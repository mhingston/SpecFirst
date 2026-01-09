package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"specfirst/internal/domain"
)

// LoadState reads state from disk.
func LoadState(path string) (domain.State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return domain.NewState(""), nil
		}
		return domain.State{}, err
	}

	if len(data) == 0 {
		return domain.NewState(""), nil
	}

	var s domain.State
	if err := json.Unmarshal(data, &s); err != nil {
		return domain.State{}, err
	}

	// Ensure maps/slices are not nil
	if s.CompletedStages == nil {
		s.CompletedStages = []string{}
	}
	if s.StageOutputs == nil {
		s.StageOutputs = make(map[string]domain.StageOutput)
	}
	if s.Attestations == nil {
		s.Attestations = make(map[string][]domain.Attestation)
	}
	if s.Epistemics.Assumptions == nil {
		s.Epistemics.Assumptions = []domain.Assumption{}
	}
	if s.Epistemics.OpenQuestions == nil {
		s.Epistemics.OpenQuestions = []domain.OpenQuestion{}
	}
	if s.Epistemics.Decisions == nil {
		s.Epistemics.Decisions = []domain.Decision{}
	}
	if s.Epistemics.Risks == nil {
		s.Epistemics.Risks = []domain.Risk{}
	}
	if s.Epistemics.Disputes == nil {
		s.Epistemics.Disputes = []domain.Dispute{}
	}
	if s.Epistemics.Confidence.ByStage == nil {
		s.Epistemics.Confidence.ByStage = make(map[string]string)
	}

	return s, nil
}

// SaveState writes the state to disk atomically using a temp file + rename pattern.
func SaveState(path string, s domain.State) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".state.*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()

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

	// Robust Rename with Retry
	// Windows often holds locks briefly (AV, Indexing, Backup). Retry helps avoid flakes.
	var renameErr error
	for attempt := 0; attempt < 5; attempt++ {
		// On Windows, Rename fails if dest exists. We try to remove it first.
		if runtime.GOOS == "windows" {
			_ = os.Remove(path)
		}

		renameErr = os.Rename(tmpPath, path)
		if renameErr == nil {
			success = true
			return nil
		}

		// If it's not a link/exist/permission error, fail immediately
		if !os.IsExist(renameErr) && !os.IsPermission(renameErr) {
			break
		}

		// Backoff before retry
		if attempt < 4 {
			time.Sleep(50 * time.Millisecond)
		}
	}

	return fmt.Errorf("failed to save state after retries: %w", renameErr)
}
