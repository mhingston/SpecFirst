package state

import (
	"fmt"
	"time"
)

// Helper methods for Epistemics

func (s *State) AddAssumption(text string, owner string) string {
	id := fmt.Sprintf("A%d", len(s.Epistemics.Assumptions)+1)
	s.Epistemics.Assumptions = append(s.Epistemics.Assumptions, Assumption{
		ID:        id,
		Text:      text,
		Status:    "open",
		Owner:     owner,
		CreatedAt: time.Now(),
	})
	return id
}

func (s *State) AddOpenQuestion(text string, tags []string, context string) string {
	id := fmt.Sprintf("Q%d", len(s.Epistemics.OpenQuestions)+1)
	s.Epistemics.OpenQuestions = append(s.Epistemics.OpenQuestions, OpenQuestion{
		ID:      id,
		Text:    text,
		Tags:    tags,
		Status:  "open",
		Context: context,
	})
	return id
}

func (s *State) AddDecision(text string, rationale string, alternatives []string) string {
	id := fmt.Sprintf("D%d", len(s.Epistemics.Decisions)+1)
	s.Epistemics.Decisions = append(s.Epistemics.Decisions, Decision{
		ID:           id,
		Text:         text,
		Rationale:    rationale,
		Alternatives: alternatives,
		Status:       "proposed",
		CreatedAt:    time.Now(),
	})
	return id
}

func (s *State) AddRisk(text string, severity string) string {
	id := fmt.Sprintf("R%d", len(s.Epistemics.Risks)+1)
	s.Epistemics.Risks = append(s.Epistemics.Risks, Risk{
		ID:       id,
		Text:     text,
		Severity: severity,
		Status:   "open",
	})
	return id
}

func (s *State) AddDispute(topic string) string {
	id := fmt.Sprintf("X%d", len(s.Epistemics.Disputes)+1)
	s.Epistemics.Disputes = append(s.Epistemics.Disputes, Dispute{
		ID:     id,
		Topic:  topic,
		Status: "open",
	})
	return id
}

func (s *State) CloseAssumption(id string, status string) bool {
	for i, a := range s.Epistemics.Assumptions {
		if a.ID == id {
			s.Epistemics.Assumptions[i].Status = status
			return true
		}
	}
	return false
}

func (s *State) ResolveOpenQuestion(id string, answer string) bool {
	for i, q := range s.Epistemics.OpenQuestions {
		if q.ID == id {
			s.Epistemics.OpenQuestions[i].Status = "resolved"
			s.Epistemics.OpenQuestions[i].Answer = answer
			return true
		}
	}
	return false
}

func (s *State) UpdateDecision(id string, status string) bool {
	for i, d := range s.Epistemics.Decisions {
		if d.ID == id {
			s.Epistemics.Decisions[i].Status = status
			return true
		}
	}
	return false
}

func (s *State) MitigateRisk(id string, mitigation string, status string) bool {
	for i, r := range s.Epistemics.Risks {
		if r.ID == id {
			s.Epistemics.Risks[i].Mitigation = mitigation
			s.Epistemics.Risks[i].Status = status
			return true
		}
	}
	return false
}

func (s *State) ResolveDispute(id string) bool {
	for i, d := range s.Epistemics.Disputes {
		if d.ID == id {
			s.Epistemics.Disputes[i].Status = "resolved"
			return true
		}
	}
	return false
}
