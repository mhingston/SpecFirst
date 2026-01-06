package state

import (
	"encoding/json"
	"testing"
)

func TestEpistemicHelpers(t *testing.T) {
	s := State{}

	// Test AddAssumption
	id := s.AddAssumption("Assumption 1", "user")
	if id != "A1" {
		t.Errorf("expected ID A1, got %s", id)
	}
	if len(s.Epistemics.Assumptions) != 1 {
		t.Errorf("expected 1 assumption, got %d", len(s.Epistemics.Assumptions))
	}

	// Test AddOpenQuestion
	id = s.AddOpenQuestion("Question 1", []string{"tag1"}, "section1")
	if id != "Q1" {
		t.Errorf("expected ID Q1, got %s", id)
	}
	if len(s.Epistemics.OpenQuestions) != 1 {
		t.Errorf("expected 1 question, got %d", len(s.Epistemics.OpenQuestions))
	}

	// Test AddDecision
	id = s.AddDecision("Decision 1", "Rationale 1", []string{"Alt 1"})
	if id != "D1" {
		t.Errorf("expected ID D1, got %s", id)
	}
	if len(s.Epistemics.Decisions) != 1 {
		t.Errorf("expected 1 decision, got %d", len(s.Epistemics.Decisions))
	}
}

func TestJSONMarshaling(t *testing.T) {
	s := State{}
	s.AddAssumption("A1", "user")

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("marshaling failed: %v", err)
	}

	var s2 State
	if err := json.Unmarshal(data, &s2); err != nil {
		t.Fatalf("unmarshaling failed: %v", err)
	}

	if len(s2.Epistemics.Assumptions) != 1 {
		t.Errorf("expected 1 assumption after unmarshal, got %d", len(s2.Epistemics.Assumptions))
	}
}
