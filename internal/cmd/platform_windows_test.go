//go:build windows
// +build windows

package cmd

import (
	"os/exec"
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
			name:    "Run dir command",
			timeout: 5 * time.Second,
			args:    "dir",
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
			args:    "timeout 10",
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
			args:    []string{"cmd.exe", "/C", "echo", "hello"},
			wantErr: false,
			checkFunc: func(out string) bool {
				return strings.Contains(out, "hello")
			},
		},
		{
			name:    "Run with dir command",
			timeout: 5 * time.Second,
			args:    []string{"cmd.exe", "/C", "dir"},
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
			args:    []string{"cmd.exe", "/C", "timeout", "10"},
			wantErr: true,
			checkFunc: func(out string) bool {
				return true
			},
		},
		{
			name:    "Single arg",
			timeout: 5 * time.Second,
			args:    []string{"cmd.exe"},
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

func TestChangeToUTF8String(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		charset charset
		want    string
	}{
		{
			name:    "Convert UTF-8 string",
			input:   "hello",
			charset: UTF8,
			want:    "hello",
		},
		{
			name:    "Convert empty UTF-8 string",
			input:   "",
			charset: UTF8,
			want:    "",
		},
		{
			name:    "Convert UTF-8 with special chars",
			input:   "测试",
			charset: UTF8,
			want:    "测试",
		},
		{
			name:    "Convert GB18030 string",
			input:   "hello",
			charset: GB18030,
			want:    "hello",
		},
		{
			name:    "Convert empty GB18030 string",
			input:   "",
			charset: GB18030,
			want:    "",
		},
		{
			name:    "Default charset",
			input:   "hello",
			charset: "unknown",
			want:    "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := changeToUTF8String(tt.input, tt.charset)
			if got != tt.want {
				t.Errorf("changeToUTF8String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHideWindowCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "Hide window for simple command",
			args: []string{"echo", "test"},
		},
		{
			name: "Hide window for complex command",
			args: []string{"cmd.exe", "/C", "dir"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &exec.Cmd{}
			hideWindowCmd(cmd)
			if cmd.SysProcAttr == nil {
				t.Error("hideWindowCmd() did not set SysProcAttr")
			}
			if !cmd.SysProcAttr.HideWindow {
				t.Error("hideWindowCmd() did not set HideWindow to true")
			}
		})
	}
}
