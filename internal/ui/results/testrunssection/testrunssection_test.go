package testrunssection

import (
	"testing"

	"github.com/adamakhtar/wing_commander/internal/testrun"
)

func TestLabel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		mode     testrun.Mode
		patterns []testrun.TestPattern
		want     string
	}{
		{
			name:     "whole suite",
			mode:     testrun.ModeRunWholeSuite,
			patterns: []testrun.TestPattern{},
			want:     "Run whole suite",
		},
		{
			name: "single selected pattern",
			mode: testrun.ModeRunSelectedPatterns,
			patterns: func() []testrun.TestPattern {
				p, _ := testrun.PatternsFromStrings([]string{"test/example_test.rb"})
				return p
			}(),
			want: "Run 1 pattern",
		},
		{
			name: "multiple selected patterns",
			mode: testrun.ModeRunSelectedPatterns,
			patterns: func() []testrun.TestPattern {
				p, _ := testrun.PatternsFromStrings([]string{"test/a.rb", "test/b.rb", "test/c.rb"})
				return p
			}(),
			want: "Run 3 patterns",
		},
		{
			name: "re-run single failure",
			mode: testrun.ModeReRunSingleFailure,
			patterns: func() []testrun.TestPattern {
				p, _ := testrun.PatternsFromStrings([]string{"test/failure.rb"})
				return p
			}(),
			want: "Re-run failure",
		},
		{
			name: "re-run all failures",
			mode: testrun.ModeReRunAllFailures,
			patterns: func() []testrun.TestPattern {
				p, _ := testrun.PatternsFromStrings([]string{"test/failure_a.rb", "test/failure_b.rb"})
				return p
			}(),
			want: "Re-run all failed",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			run := testrun.TestRun{
				Patterns: tt.patterns,
				Mode:     string(tt.mode),
			}

			if got := Label(run); got != tt.want {
				t.Fatalf("Label() = %q, want %q", got, tt.want)
			}
		})
	}
}
