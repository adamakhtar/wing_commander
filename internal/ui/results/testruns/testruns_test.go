package testruns

import "testing"

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
				Filepaths: tt.filepaths,
				Mode:      tt.mode,
			}

			if got := run.Label(); got != tt.want {
				t.Fatalf("Label() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDefaultMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		filepaths []string
		want      Mode
	}{
		{
			name:      "no filepaths defaults to whole suite",
			filepaths: nil,
			want:      ModeRunWholeSuite,
		},
		{
			name:      "blank filepath treated as whole suite",
			filepaths: []string{""},
			want:      ModeRunWholeSuite,
		},
		{
			name:      "multiple filepaths default to selected patterns",
			filepaths: []string{"test/a.rb", "test/b.rb"},
			want:      ModeRunSelectedPatterns,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := defaultMode(tt.filepaths); got != tt.want {
				t.Fatalf("defaultMode() = %q, want %q", got, tt.want)
			}
		})
	}
}
