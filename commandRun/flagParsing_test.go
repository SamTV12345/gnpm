package commandRun

import "testing"

func TestFilterArgs(t *testing.T) {
	var args = []string{"--pm", "npm", "--global", "typescript", "--save-dev", "--runtimeVersion"}
	var filteredArgs = FilterArgs(args)
	if len(filteredArgs) != 5 {
		t.Errorf("Expected 5 arguments, got %d", len(filteredArgs))
	}
}

func TestFilterArgsWithValue(t *testing.T) {
	var args = []string{"--pm", "npm", "--global", "typescript", "--save-dev", "--runtimeVersion='>=18.0.0'"}
	var filteredArgs = FilterArgs(args)
	if len(filteredArgs) != 5 {
		t.Errorf("Expected 5 arguments, got %d", len(filteredArgs))
	}
}
