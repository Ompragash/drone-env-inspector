// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    Args
		wantErr bool
	}{
		{
			name: "missing env_name",
			args: Args{
				EnvName: "",
			},
			wantErr: true,
		},
		{
			name: "valid single env_name",
			args: Args{
				EnvName: "MY_VAR",
			},
			wantErr: false,
		},
		{
			name: "valid multiple env_names",
			args: Args{
				EnvName: "VAR1,VAR2,VAR3",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExec(t *testing.T) {
	// Create a temporary directory for output files
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "drone_output")
	secretOutputFile := filepath.Join(tempDir, "harness_secret_output")

	// Set up environment variables for the test
	os.Setenv("DRONE_OUTPUT", outputFile)
	os.Setenv("HARNESS_OUTPUT_SECRET_FILE", secretOutputFile)
	os.Setenv("TEST_VAR_1", "value1")
	os.Setenv("TEST_VAR_2", "value2")
	defer func() {
		os.Unsetenv("DRONE_OUTPUT")
		os.Unsetenv("HARNESS_OUTPUT_SECRET_FILE")
		os.Unsetenv("TEST_VAR_1")
		os.Unsetenv("TEST_VAR_2")
	}()

	tests := []struct {
		name           string
		args           Args
		wantErr        bool
		expectedOutput string
		checkSecret    bool
	}{
		{
			name: "single existing env var",
			args: Args{
				EnvName: "TEST_VAR_1",
				Secret:  false,
			},
			wantErr:        false,
			expectedOutput: "TEST_VAR_1=value1\n",
			checkSecret:    false,
		},
		{
			name: "multiple existing env vars",
			args: Args{
				EnvName: "TEST_VAR_1,TEST_VAR_2",
				Secret:  false,
			},
			wantErr:        false,
			expectedOutput: "TEST_VAR_1=value1\nTEST_VAR_2=value2\n",
			checkSecret:    false,
		},
		{
			name: "non-existing env var",
			args: Args{
				EnvName: "NON_EXISTING_VAR",
				Secret:  false,
			},
			wantErr:        false,
			expectedOutput: "NON_EXISTING_VAR=\n",
			checkSecret:    false,
		},
		{
			name: "secret output",
			args: Args{
				EnvName: "TEST_VAR_1",
				Secret:  true,
			},
			wantErr:        false,
			expectedOutput: "TEST_VAR_1=value1\n",
			checkSecret:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear the output files before each test
			os.WriteFile(outputFile, []byte{}, 0644)
			os.WriteFile(secretOutputFile, []byte{}, 0644)

			err := Exec(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var fileToCheck string
			if tt.checkSecret {
				fileToCheck = secretOutputFile
			} else {
				fileToCheck = outputFile
			}

			content, err := os.ReadFile(fileToCheck)
			if err != nil {
				t.Errorf("failed to read output file: %v", err)
				return
			}

			if string(content) != tt.expectedOutput {
				t.Errorf("output mismatch: got %q, want %q", string(content), tt.expectedOutput)
			}
		})
	}
}

func TestExecWithSpacesInEnvNames(t *testing.T) {
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "drone_output")

	os.Setenv("DRONE_OUTPUT", outputFile)
	os.Setenv("SPACED_VAR", "spaced_value")
	defer func() {
		os.Unsetenv("DRONE_OUTPUT")
		os.Unsetenv("SPACED_VAR")
	}()

	args := Args{
		EnvName: "  SPACED_VAR  ",
		Secret:  false,
	}

	err := Exec(context.Background(), args)
	if err != nil {
		t.Errorf("Exec() unexpected error: %v", err)
		return
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Errorf("failed to read output file: %v", err)
		return
	}

	expected := "SPACED_VAR=spaced_value\n"
	if string(content) != expected {
		t.Errorf("output mismatch: got %q, want %q", string(content), expected)
	}
}

func TestExecEmptyEnvNameInList(t *testing.T) {
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "drone_output")

	os.Setenv("DRONE_OUTPUT", outputFile)
	os.Setenv("VAR_A", "a")
	os.Setenv("VAR_B", "b")
	defer func() {
		os.Unsetenv("DRONE_OUTPUT")
		os.Unsetenv("VAR_A")
		os.Unsetenv("VAR_B")
	}()

	// Test with empty entries in comma-separated list
	args := Args{
		EnvName: "VAR_A,,VAR_B",
		Secret:  false,
	}

	err := Exec(context.Background(), args)
	if err != nil {
		t.Errorf("Exec() unexpected error: %v", err)
		return
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Errorf("failed to read output file: %v", err)
		return
	}

	// Should only have VAR_A and VAR_B, skipping the empty entry
	if !strings.Contains(string(content), "VAR_A=a") || !strings.Contains(string(content), "VAR_B=b") {
		t.Errorf("output missing expected variables: got %q", string(content))
	}
}

