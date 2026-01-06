package state

import (
	"time"
)

// RecordApproval adds or updates an approval for a stage.
// It returns true if the approval was updated (existed previously).
func (s *State) RecordApproval(stageID, role, user, notes string) bool {
	if s.Approvals == nil {
		s.Approvals = make(map[string][]Approval)
	}

	newApproval := Approval{
		Role:       role,
		ApprovedBy: user,
		ApprovedAt: time.Now().UTC(),
		Notes:      notes,
	}

	updated := false
	approvals := s.Approvals[stageID]
	for i, existing := range approvals {
		if existing.Role == role {
			approvals[i] = newApproval
			updated = true
			break
		}
	}
	if !updated {
		s.Approvals[stageID] = append(approvals, newApproval)
	} else {
		s.Approvals[stageID] = approvals // Re-assign if we updated in place (map slice)
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

// HasApproval checks if a specific role has approved a stage.
func (s *State) HasApproval(stageID, role string) bool {
	approvals, ok := s.Approvals[stageID]
	if !ok {
		return false
	}
	for _, a := range approvals {
		if a.Role == role {
			return true
		}
	}
	return false
}
