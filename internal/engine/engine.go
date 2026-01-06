package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"strings"

	"specfirst/internal/assets"
	"specfirst/internal/config"
	"specfirst/internal/protocol"
	"specfirst/internal/state"
	"specfirst/internal/store"
)

// Engine is the main application coordinator.
type Engine struct {
	Config   config.Config
	Protocol protocol.Protocol
	State    state.State
}

// NewEngine creates a new Engine instance.
func NewEngine(cfg config.Config, proto protocol.Protocol, s state.State) *Engine {
	return &Engine{
		Config:   cfg,
		Protocol: proto,
		State:    s,
	}
}

// Load loads the engine from the filesystem.
// It handles the full configuration loading flow: Config -> Protocol -> State.
func Load(protocolOverride string) (*Engine, error) {
	// 1. Load Config
	cfg, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	// 2. Resolve Active Protocol
	activeProto := cfg.Protocol
	if protocolOverride != "" {
		activeProto = protocolOverride
	} else if cfg.Protocol == "" {
		activeProto = assets.DefaultProtocolName
	}

	// 3. Load Protocol
	proto, err := loadProtocol(activeProto)
	if err != nil {
		return nil, fmt.Errorf("loading protocol %s: %w", activeProto, err)
	}

	// 4. Load State
	s, err := loadState()
	if err != nil {
		return nil, fmt.Errorf("loading state: %w", err)
	}

	// 5. Initialize State if needed
	s = ensureStateInitialized(s, proto)

	return NewEngine(cfg, proto, s), nil
}

// AttestStage records an attestation for a stage.
// Returns a list of warnings and any error that occurred.
func (e *Engine) AttestStage(stageID, role, user, status, notes string, conditions []string) ([]string, error) {
	var warnings []string

	// 1. Verify approval is declared in protocol
	declared := false
	for _, approval := range e.Protocol.Approvals {
		if approval.Stage == stageID && approval.Role == role {
			declared = true
			break
		}
	}
	if !declared {
		return nil, fmt.Errorf("approval not declared in protocol: stage=%s role=%s", stageID, role)
	}

	// 2. Warn if stage is not completed (but allow it)
	if !e.State.IsStageCompleted(stageID) {
		warnings = append(warnings, fmt.Sprintf("stage %s is not yet completed; attestation recorded preemptively", stageID))
	}

	// 3. Record approval in state (as Attestation)
	attestation := state.Attestation{
		Role:       role,
		AttestedBy: user,
		Status:     status,
		Rationale:  notes,
		Conditions: conditions,
		Date:       time.Now().UTC(),
	}
	updated := e.State.RecordAttestation(stageID, attestation)
	if updated {
		warnings = append(warnings, fmt.Sprintf("updating existing attestation for role %s", role))
	}

	// 4. Save state
	return warnings, e.SaveState()
}

// SaveState saves the current state to disk.
func (e *Engine) SaveState() error {
	return state.Save(store.StatePath(), e.State)
}

// Private helpers (migrated from cmd/helpers.go basically)

func loadConfig() (config.Config, error) {
	cfg, err := config.Load(store.ConfigPath())
	if err != nil {
		return config.Config{}, err
	}
	if cfg.Protocol == "" {
		cfg.Protocol = assets.DefaultProtocolName
	}
	if cfg.ProjectName == "" {
		if wd, err := os.Getwd(); err == nil {
			cfg.ProjectName = filepath.Base(wd)
		} else {
			cfg.ProjectName = "project" // Fallback when working directory is unavailable
		}
	}
	if cfg.CustomVars == nil {
		cfg.CustomVars = map[string]string{}
	}
	if cfg.Constraints == nil {
		cfg.Constraints = map[string]string{}
	}
	return cfg, nil
}

func loadProtocol(name string) (protocol.Protocol, error) {
	// If name looks like a path or has .yaml extension, load it directly
	if filepath.IsAbs(name) || strings.Contains(name, string(os.PathSeparator)) || strings.HasSuffix(name, ".yaml") {
		return protocol.Load(filepath.Clean(name))
	}
	// Otherwise treat as a protocol name in the protocols directory
	path := store.ProtocolsPath(name + ".yaml")
	return protocol.Load(path)
}

func loadState() (state.State, error) {
	s, err := state.Load(store.StatePath())
	if err != nil {
		return state.State{}, err
	}
	// Note: The initialization logic (ensuring non-nil maps) is already in state.Load
	// allowing us to trust the returned struct is safe to use.
	return s, nil
}

func ensureStateInitialized(s state.State, proto protocol.Protocol) state.State {
	if s.Protocol == "" {
		s.Protocol = proto.Name
	}
	if s.SpecVersion == "" {
		s.SpecVersion = proto.Version
	}
	// Other fields are initialized in state.Load/NewState
	return s
}
