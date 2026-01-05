package task

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    TaskList
		wantErr bool
	}{
		{
			name:    "simple yaml",
			content: `tasks: [{id: T1, title: Task 1, goal: Goal 1}]`,
			want:    TaskList{Tasks: []Task{{ID: "T1", Title: "Task 1", Goal: "Goal 1"}}},
			wantErr: false,
		},
		{
			name:    "markdown wrap",
			content: "Here are the tasks:\n```yaml\ntasks:\n  - id: T1\n    title: Task 1\n    goal: Goal 1\n```",
			want:    TaskList{Tasks: []Task{{ID: "T1", Title: "Task 1", Goal: "Goal 1"}}},
			wantErr: false,
		},
		{
			name:    "invalid content",
			content: "not yaml at all",
			want:    TaskList{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskList_Validate(t *testing.T) {
	tests := []struct {
		name string
		tl   TaskList
		want int // number of warnings
	}{
		{
			name: "valid",
			tl: TaskList{Tasks: []Task{
				{ID: "T1", Title: "T1", Goal: "G1"},
				{ID: "T2", Title: "T2", Goal: "G2", Dependencies: []string{"T1"}},
			}},
			want: 0,
		},
		{
			name: "duplicate id",
			tl: TaskList{Tasks: []Task{
				{ID: "T1", Title: "T1", Goal: "G1"},
				{ID: "T1", Title: "T1", Goal: "G1"},
			}},
			want: 1,
		},
		{
			name: "missing dependency",
			tl: TaskList{Tasks: []Task{
				{ID: "T1", Title: "T1", Goal: "G1", Dependencies: []string{"T2"}},
			}},
			want: 1,
		},
		{
			name: "circular dependency",
			tl: TaskList{Tasks: []Task{
				{ID: "T1", Title: "T1", Goal: "G1", Dependencies: []string{"T2"}},
				{ID: "T2", Title: "T2", Goal: "G2", Dependencies: []string{"T1"}},
			}},
			want: 2, // T1 -> T2 -> T1 and T2 -> T1 -> T2
		},
		{
			name: "empty id/title/goal",
			tl: TaskList{Tasks: []Task{
				{ID: "", Title: "", Goal: ""},
			}},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tl.Validate()
			if len(got) != tt.want {
				t.Errorf("Validate() = %v warnings, want %v", len(got), tt.want)
			}
		})
	}
}
