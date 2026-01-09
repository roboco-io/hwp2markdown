package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestConvertCommand(t *testing.T) {
	// Find the sample HWPX file
	fixtureDir := "fixtures"
	sampleFile := filepath.Join(fixtureDir, "sample.hwpx")

	if _, err := os.Stat(sampleFile); os.IsNotExist(err) {
		t.Skipf("sample file not found: %s", sampleFile)
	}

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "hwp2markdown_test", "../cmd/hwp2markdown")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}
	defer os.Remove("hwp2markdown_test")

	tests := []struct {
		name       string
		args       []string
		wantErr    bool
		wantOutput []string
	}{
		{
			name:    "basic convert",
			args:    []string{"convert", sampleFile},
			wantErr: false,
		},
		{
			name:    "convert with verbose",
			args:    []string{"convert", sampleFile, "-v"},
			wantErr: false,
		},
		{
			name:    "convert non-existent file",
			args:    []string{"convert", "nonexistent.hwpx"},
			wantErr: true,
		},
		{
			name:    "convert unsupported format",
			args:    []string{"convert", "test.txt"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("./hwp2markdown_test", tc.args...)
			output, err := cmd.CombinedOutput()

			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v\noutput: %s", err, output)
				}
			}

			for _, want := range tc.wantOutput {
				if !strings.Contains(string(output), want) {
					t.Errorf("output should contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestExtractCommand(t *testing.T) {
	fixtureDir := "fixtures"
	sampleFile := filepath.Join(fixtureDir, "sample.hwpx")

	if _, err := os.Stat(sampleFile); os.IsNotExist(err) {
		t.Skipf("sample file not found: %s", sampleFile)
	}

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "hwp2markdown_test", "../cmd/hwp2markdown")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}
	defer os.Remove("hwp2markdown_test")

	tests := []struct {
		name       string
		args       []string
		wantErr    bool
		wantFormat string
	}{
		{
			name:       "extract as json",
			args:       []string{"extract", sampleFile},
			wantErr:    false,
			wantFormat: "json",
		},
		{
			name:       "extract as text",
			args:       []string{"extract", sampleFile, "--format", "text"},
			wantErr:    false,
			wantFormat: "text",
		},
		{
			name:    "extract non-existent file",
			args:    []string{"extract", "nonexistent.hwpx"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("./hwp2markdown_test", tc.args...)
			output, err := cmd.CombinedOutput()

			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v\noutput: %s", err, output)
				}

				if tc.wantFormat == "json" && !strings.Contains(string(output), "{") {
					t.Errorf("expected JSON output, got: %s", output)
				}
			}
		})
	}
}

func TestProvidersCommand(t *testing.T) {
	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "hwp2markdown_test", "../cmd/hwp2markdown")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}
	defer os.Remove("hwp2markdown_test")

	cmd := exec.Command("./hwp2markdown_test", "providers")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("unexpected error: %v\noutput: %s", err, output)
	}

	// Check that all providers are listed
	providers := []string{"anthropic", "openai", "gemini", "ollama"}
	for _, p := range providers {
		if !strings.Contains(string(output), p) {
			t.Errorf("output should contain provider %q, got: %s", p, output)
		}
	}
}

func TestVersionCommand(t *testing.T) {
	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "hwp2markdown_test", "../cmd/hwp2markdown")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}
	defer os.Remove("hwp2markdown_test")

	cmd := exec.Command("./hwp2markdown_test", "version")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("unexpected error: %v\noutput: %s", err, output)
	}

	if !strings.Contains(string(output), "hwp2markdown") {
		t.Errorf("output should contain 'hwp2markdown', got: %s", output)
	}
}

func TestConfigCommand(t *testing.T) {
	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "hwp2markdown_test", "../cmd/hwp2markdown")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}
	defer os.Remove("hwp2markdown_test")

	t.Run("config show", func(t *testing.T) {
		cmd := exec.Command("./hwp2markdown_test", "config", "show")
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Errorf("unexpected error: %v\noutput: %s", err, output)
		}

		if !strings.Contains(string(output), "default_provider") {
			t.Errorf("output should contain 'default_provider', got: %s", output)
		}
	})

	t.Run("config path", func(t *testing.T) {
		cmd := exec.Command("./hwp2markdown_test", "config", "path")
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Errorf("unexpected error: %v\noutput: %s", err, output)
		}

		if !strings.Contains(string(output), "config.yaml") {
			t.Errorf("output should contain 'config.yaml', got: %s", output)
		}
	})
}

func TestHelpCommand(t *testing.T) {
	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "hwp2markdown_test", "../cmd/hwp2markdown")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}
	defer os.Remove("hwp2markdown_test")

	cmd := exec.Command("./hwp2markdown_test", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("unexpected error: %v\noutput: %s", err, output)
	}

	expectedStrings := []string{"hwp2markdown", "convert", "extract", "providers", "config"}
	for _, s := range expectedStrings {
		if !strings.Contains(string(output), s) {
			t.Errorf("output should contain %q, got: %s", s, output)
		}
	}
}
