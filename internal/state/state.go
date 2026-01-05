package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type State struct {
	Protocol        string                 `json:"protocol"`
	CurrentStage    string                 `json:"current_stage"`
	CompletedStages []string               `json:"completed_stages"`
	StartedAt       time.Time              `json:"started_at"`
	SpecVersion     string                 `json:"spec_version"`
	StageOutputs    map[string]StageOutput `json:"stage_outputs"`
	Approvals       map[string][]Approval  `json:"approvals"`
}

type StageOutput struct {
	CompletedAt time.Time `json:"completed_at"`
	Files       []string  `json:"files"`
	PromptHash  string    `json:"prompt_hash"`
}

type Approval struct {
	Role       string    `json:"role"`
	ApprovedBy string    `json:"approved_by"`
	ApprovedAt time.Time `json:"approved_at"`
	Notes      string    `json:"notes,omitempty"`
}

func Load(path string) (State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return State{}, nil
		}
		return State{}, err
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return State{}, err
	}

	return s, nil
}

// Save writes the state to disk atomically using a temp file + rename pattern.
// This prevents corruption from partial writes or concurrent CLI invocations.
func Save(path string, s State) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	// Write to temp file in same directory (required for atomic rename)
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".state.*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()

	// Clean up temp file on any error
	success := false
	defer func() {
		if !success {
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	// Atomic rename
	if runtime.GOOS == "windows" {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}
	success = true
	return nil
}

func (s State) IsStageCompleted(id string) bool {
	for _, stage := range s.CompletedStages {
		if stage == id {
			return true
		}
	}
	return false
}
