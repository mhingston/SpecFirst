package template

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"

	"specfirst/internal/protocol"
)

type Input struct {
	Name    string
	Content string
}

type Data struct {
	StageName   string
	ProjectName string
	Inputs      []Input
	Outputs     []string
	Intent      string
	Language    string
	Framework   string
	CustomVars  map[string]string
	Constraints map[string]string

	StageType      string
	Prompt         *protocol.PromptConfig
	OutputContract *protocol.OutputContract
}

type cachedTemplate struct {
	tmpl    *template.Template
	modTime time.Time
}

// templateCache stores parsed templates keyed by file path.
// This cache is designed for CLI usage (short-lived processes) and does not
// implement eviction. If used in a long-running service, templates would
// accumulate indefinitely. For such use cases, consider using a bounded cache
// with LRU eviction.
var (
	templateCache   = make(map[string]cachedTemplate)
	templateCacheMu sync.RWMutex
)

func Render(path string, data Data) (string, error) {
	// Check file modification time
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("reading template %s: %w", path, err)
	}
	currentModTime := info.ModTime()

	// Check cache with read lock
	templateCacheMu.RLock()
	cached, found := templateCache[path]
	templateCacheMu.RUnlock()

	// Use cached template if found and not stale
	if found && !cached.modTime.Before(currentModTime) {
		var buf bytes.Buffer
		if err := cached.tmpl.Execute(&buf, data); err != nil {
			return "", fmt.Errorf("executing template %s: %w", path, err)
		}
		return buf.String(), nil
	}

	// Need to load/reload template - acquire write lock
	templateCacheMu.Lock()
	defer templateCacheMu.Unlock()

	// Re-check after acquiring write lock (another goroutine may have updated)
	if cached, found := templateCache[path]; found && !cached.modTime.Before(currentModTime) {
		var buf bytes.Buffer
		if err := cached.tmpl.Execute(&buf, data); err != nil {
			return "", fmt.Errorf("executing template %s: %w", path, err)
		}
		return buf.String(), nil
	}

	// Load and parse template
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading template %s: %w", path, err)
	}

	tmpl, err := template.New("stage").Funcs(template.FuncMap{
		"join": strings.Join,
	}).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("parsing template %s: %w", path, err)
	}

	templateCache[path] = cachedTemplate{
		tmpl:    tmpl,
		modTime: currentModTime,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template %s: %w", path, err)
	}

	return buf.String(), nil
}
