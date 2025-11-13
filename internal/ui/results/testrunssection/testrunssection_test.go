package testrunssection

import (
	"testing"

	"github.com/adamakhtar/wing_commander/internal/testrun"
)

func TestLabel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		mode      testrun.Mode
		filepaths []string
		want      string
	}{
		{
			name: "whole suite",
			mode: testrun.ModeRunWholeSuite,
			want: "Run whole suite",
		},
		{
			name:      "single selected pattern",
			mode:      testrun.ModeRunSelectedPatterns,
			filepaths: []string{"test/example_test.rb"},
			want:      "Run 1 pattern",
		},
		{
			name:      "multiple selected patterns",
			mode:      testrun.ModeRunSelectedPatterns,
			filepaths: []string{"test/a.rb", "test/b.rb", "test/c.rb"},
			want:      "Run 3 patterns",
		},
		{
			name:      "re-run single failure",
			mode:      testrun.ModeReRunSingleFailure,
			filepaths: []string{"test/failure.rb"},
			want:      "Re-run failure",
		},
		{
			name:      "re-run all failures",
			mode:      testrun.ModeReRunAllFailures,
			filepaths: []string{"test/failure_a.rb", "test/failure_b.rb"},
			want:      "Re-run all failed",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			run := testrun.TestRun{
				Filepaths: tt.filepaths,
				Mode:      string(tt.mode),
			}

			if got := Label(run); got != tt.want {
				t.Fatalf("Label() = %q, want %q", got, tt.want)
			}
		})
	}
}
