package testruns

import (
	"testing"

	"github.com/adamakhtar/wing_commander/internal/types"
)

func TestTestRunLabel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		mode      Mode
		filepaths []string
		want      string
	}{
		{
			name: "whole suite",
			mode: ModeRunWholeSuite,
			want: "Whole suite",
		},
		{
			name:      "single selected pattern",
			mode:      ModeRunSelectedPatterns,
			filepaths: []string{"test/example_test.rb"},
			want:      "1 test pattern",
		},
		{
			name:      "multiple selected patterns",
			mode:      ModeRunSelectedPatterns,
			filepaths: []string{"test/a.rb", "test/b.rb", "test/c.rb"},
			want:      "3 test patterns",
		},
		{
			name:      "re-run single failure",
			mode:      ModeReRunSingleFailure,
			filepaths: []string{"test/failure.rb"},
			want:      "Re-run single failure",
		},
		{
			name:      "re-run all failures",
			mode:      ModeReRunAllFailures,
			filepaths: []string{"test/failure_a.rb", "test/failure_b.rb"},
			want:      "Re-run failed",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			run := TestRun{
				TestRun: types.TestRun{
					Filepaths: tt.filepaths,
					Mode:      string(tt.mode),
				},
			}

			if got := run.Label(); got != tt.want {
				t.Fatalf("Label() = %q, want %q", got, tt.want)
			}
		})
	}
}
