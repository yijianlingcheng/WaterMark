//go:build darwin
// +build darwin

package cmd

import (
	"strings"
	"testing"
	"time"

	"WaterMark/pkg"
)

func TestCommandRun(t *testing.T) {
	tests := []struct {
		name      string
		timeout   time.Duration
		args      string
		wantErr   bool
		checkFunc func(string) bool
	}{
		{
			name:    "Run simple echo command",
			timeout: 5 * time.Second,
			args:    "echo hello",
			wantErr: false,
			checkFunc: func(out string) bool {
				return strings.Contains(out, "hello")
			},
		},
		{
			name:    "Run ls command",
			timeout: 5 * time.Second,
			args:    "ls -la",
			wantErr: false,
			checkFunc: func(out string) bool {
				return len(out) > 0
			},
		},
		{
			name:    "Run invalid command",
			timeout: 5 * time.Second,
			args:    "invalid_command_xyz",
			wantErr: true,
			checkFunc: func(out string) bool {
				return true
			},
		},
		{
			name:    "Command with timeout",
			timeout: 100 * time.Millisecond,
			args:    "sleep 10",
			wantErr: true,
			checkFunc: func(out string) bool {
				return true
			},
		},
		{
			name:    "Empty command",
			timeout: 5 * time.Second,
			args:    "",
			wantErr: false,
			checkFunc: func(out string) bool {
				return true
			},
		},
		{
			name:    "Command with spaces",
			timeout: 5 * time.Second,
			args:    "echo test with spaces",
			wantErr: false,
			checkFunc: func(out string) bool {
				return strings.Contains(out, "test with spaces")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CommandRun(tt.timeout, tt.args)

			if tt.wantErr {
				if err == pkg.NoError {
					t.Errorf("CommandRun() expected error but got none")
				}
			} else {
				if err != pkg.NoError {
					t.Errorf("CommandRun() unexpected error: %v", err)
				}
			}

			if !tt.checkFunc(got) {
				t.Errorf("CommandRun() output check failed for: %s", got)
			}
		})
	}
}

func TestCommandRunWithArgs(t *testing.T) {
	tests := []struct {
		name      string
		timeout   time.Duration
		args      []string
		wantErr   bool
		checkFunc func(string) bool
	}{
		{
			name:    "Run with echo command",
			timeout: 5 * time.Second,
			args:    []string{"echo", "hello"},
			wantErr: false,
			checkFunc: func(out string) bool {
				return strings.Contains(out, "hello")
			},
		},
		{
			name:    "Run with ls command",
			timeout: 5 * time.Second,
			args:    []string{"ls", "-la"},
			wantErr: false,
			checkFunc: func(out string) bool {
				return len(out) > 0
			},
		},
		{
			name:    "Run with invalid command",
			timeout: 5 * time.Second,
			args:    []string{"invalid_command_xyz"},
			wantErr: true,
			checkFunc: func(out string) bool {
				return true
			},
		},
		{
			name:    "Command with timeout",
			timeout: 100 * time.Millisecond,
			args:    []string{"sleep", "10"},
			wantErr: true,
			checkFunc: func(out string) bool {
				return true
			},
		},
		{
			name:    "Single arg",
			timeout: 5 * time.Second,
			args:    []string{"echo"},
			wantErr: false,
			checkFunc: func(out string) bool {
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CommandRunWithArgs(tt.timeout, tt.args)

			if tt.wantErr {
				if err == pkg.NoError {
					t.Errorf("CommandRunWithArgs() expected error but got none")
				}
			} else {
				if err != pkg.NoError {
					t.Errorf("CommandRunWithArgs() unexpected error: %v", err)
				}
			}

			if !tt.checkFunc(got) {
				t.Errorf("CommandRunWithArgs() output check failed for: %s", got)
			}
		})
	}
}
