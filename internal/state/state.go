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
	Epistemics      Epistemics             `json:"epistemics,omitempty"`
}

type Epistemics struct {
	Assumptions   []Assumption   `json:"assumptions,omitempty"`
	OpenQuestions []OpenQuestion `json:"open_questions,omitempty"`
	Decisions     []Decision     `json:"decisions,omitempty"`
	Risks         []Risk         `json:"risks,omitempty"`
	Disputes      []Dispute      `json:"disputes,omitempty"`
	Confidence    Confidence     `json:"confidence,omitempty"`
}

type Assumption struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Status    string    `json:"status"` // open, validated, invalidated
	Owner     string    `json:"owner,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type OpenQuestion struct {
	ID      string   `json:"id"`
	Text    string   `json:"text"`
	Tags    []string `json:"tags,omitempty"`
	Status  string   `json:"status"` // open, resolved, deferred
	Answer  string   `json:"answer,omitempty"`
	Context string   `json:"context,omitempty"` // file or section reference
}

type Decision struct {
	ID           string    `json:"id"`
	Text         string    `json:"text"`
	Rationale    string    `json:"rationale"`
	Alternatives []string  `json:"alternatives,omitempty"`
	Status       string    `json:"status"` // proposed, accepted, reversed
	CreatedAt    time.Time `json:"created_at"`
}

type Risk struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	Severity   string `json:"severity"` // low, medium, high
	Mitigation string `json:"mitigation,omitempty"`
	Status     string `json:"status"` // open, mitigated, accepted
}

type Dispute struct {
	ID        string     `json:"id"`
	Topic     string     `json:"topic"`
	Positions []Position `json:"positions,omitempty"`
	Status    string     `json:"status"` // open, resolved
}

type Position struct {
	Owner string `json:"owner"`
	Claim string `json:"claim"`
}

type Confidence struct {
	Overall string            `json:"overall"` // low, medium, high
	ByStage map[string]string `json:"by_stage,omitempty"`
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

func NewState(protocol string) State {
	return State{
		Protocol:        protocol,
		StartedAt:       time.Now(),
		CompletedStages: []string{},
		StageOutputs:    make(map[string]StageOutput),
		Approvals:       make(map[string][]Approval),
		Epistemics: Epistemics{
			Assumptions:   []Assumption{},
			OpenQuestions: []OpenQuestion{},
			Decisions:     []Decision{},
			Risks:         []Risk{},
			Disputes:      []Dispute{},
			Confidence: Confidence{
				ByStage: make(map[string]string),
			},
		},
	}
}

func Load(path string) (State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewState(""), nil
		}
		return State{}, err
	}

	if len(data) == 0 {
		return NewState(""), nil
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return State{}, err
	}

	// Ensure maps/slices are not nil even after unmarshal
	if s.CompletedStages == nil {
		s.CompletedStages = []string{}
	}
	if s.StageOutputs == nil {
		s.StageOutputs = make(map[string]StageOutput)
	}
	if s.Approvals == nil {
		s.Approvals = make(map[string][]Approval)
	}
	if s.Epistemics.Assumptions == nil {
		s.Epistemics.Assumptions = []Assumption{}
	}
	if s.Epistemics.OpenQuestions == nil {
		s.Epistemics.OpenQuestions = []OpenQuestion{}
	}
	if s.Epistemics.Decisions == nil {
		s.Epistemics.Decisions = []Decision{}
	}
	if s.Epistemics.Risks == nil {
		s.Epistemics.Risks = []Risk{}
	}
	if s.Epistemics.Disputes == nil {
		s.Epistemics.Disputes = []Dispute{}
	}
	if s.Epistemics.Confidence.ByStage == nil {
		s.Epistemics.Confidence.ByStage = make(map[string]string)
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
