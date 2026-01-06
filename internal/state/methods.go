package state

import (
	"time"
)

// RecordAttestation adds or updates an attestation for a stage.
// It returns true if the attestation was updated (existed previously).
func (s *State) RecordAttestation(stageID string, attestation Attestation) bool {
	if s.Attestations == nil {
		s.Attestations = make(map[string][]Attestation)
	}

	updated := false
	attestations := s.Attestations[stageID]
	for i, existing := range attestations {
		if existing.Role == attestation.Role {
			attestations[i] = attestation
			updated = true
			break
		}
	}
	if !updated {
		s.Attestations[stageID] = append(attestations, attestation)
	} else {
		s.Attestations[stageID] = attestations // Re-assign if we updated in place (map slice)
	}
	return updated
}

// UpdateStageOutput records the output for a completed stage.
func (s *State) UpdateStageOutput(stageID string, files []string, promptHash string) {
	if s.StageOutputs == nil {
		s.StageOutputs = make(map[string]StageOutput)
	}
	s.StageOutputs[stageID] = StageOutput{
		CompletedAt: time.Now().UTC(),
		Files:       files,
		PromptHash:  promptHash,
		// Assuming OutputRelPath logic is handled by caller or unnecessary here
	}

	// Mark stage as completed if not already
	alreadyCompleted := false
	for _, id := range s.CompletedStages {
		if id == stageID {
			alreadyCompleted = true
			break
		}
	}
	if !alreadyCompleted {
		s.CompletedStages = append(s.CompletedStages, stageID)
	}
}

// HasAttestation checks if a specific role has an attestation with a specific status.
// If status is empty, checks for any attestation.
func (s *State) HasAttestation(stageID, role, status string) bool {
	attestations, ok := s.Attestations[stageID]
	if !ok {
		return false
	}
	for _, a := range attestations {
		if a.Role == role {
			if status == "" || a.Status == status {
				return true
			}
		}
	}
	return false
}
