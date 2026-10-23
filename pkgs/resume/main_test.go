package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAssetNameFor(t *testing.T) {
	tests := []struct {
		arch string
		os   string
		want string
	}{
		{"amd64", "linux", "X64-Linux"},
		{"arm64", "linux", "ARM64-Linux"},
		{"amd64", "darwin", "X64-macOS"},
		{"arm64", "darwin", "ARM64-macOS"},
		{"386", "linux", ""},
		{"amd64", "windows", ""},
		{"", "", ""},
	}

	for _, tt := range tests {
		got := assetNameFor(tt.arch, tt.os)
		if got != tt.want {
			t.Errorf("assetNameFor(%q, %q) = %q, want %q", tt.arch, tt.os, got, tt.want)
		}
	}
}

func TestDetermineDownloadArgs(t *testing.T) {
	tests := []struct {
		name             string
		runID            string
		tempDir          string
		currentAssetName string
		all              bool
		wantArgs         []string
		wantUnknown      bool
	}{
		{
			name:             "download matched arch only",
			runID:            "12345",
			tempDir:          "/tmp/dir",
			currentAssetName: "X64-Linux",
			all:              false,
			wantArgs:         []string{"run", "download", "12345", "--dir", "/tmp/dir", "--pattern", "*X64-Linux*"},
			wantUnknown:      false,
		},
		{
			name:             "download all when all flag is true",
			runID:            "12345",
			tempDir:          "/tmp/dir",
			currentAssetName: "X64-Linux",
			all:              true,
			wantArgs:         []string{"run", "download", "12345", "--dir", "/tmp/dir"},
			wantUnknown:      false,
		},
		{
			name:             "fallback to all with unknown arch flag when asset name is empty",
			runID:            "12345",
			tempDir:          "/tmp/dir",
			currentAssetName: "",
			all:              false,
			wantArgs:         []string{"run", "download", "12345", "--dir", "/tmp/dir"},
			wantUnknown:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArgs, gotUnknown := determineDownloadArgs(tt.runID, tt.tempDir, tt.currentAssetName, tt.all)
			if diff := cmp.Diff(tt.wantArgs, gotArgs); diff != "" {
				t.Errorf("determineDownloadArgs() args mismatch (-want +got):\n%s", diff)
			}
			if gotUnknown != tt.wantUnknown {
				t.Errorf("determineDownloadArgs() isUnknownArch = %v, want %v", gotUnknown, tt.wantUnknown)
			}
		})
	}
}

func TestBuildRunURL(t *testing.T) {
	tests := []struct {
		name       string
		serverURL  string
		repository string
		runID      string
		want       string
	}{
		{
			name:       "default fallback when env vars are empty",
			serverURL:  "",
			repository: "",
			runID:      "123456",
			want:       "https://github.com/kachick/nixpkgs-reviewing-workspace/actions/runs/123456",
		},
		{
			name:       "custom server and repository",
			serverURL:  "https://custom.github.com",
			repository: "custom-owner/custom-repo",
			runID:      "789012",
			want:       "https://custom.github.com/custom-owner/custom-repo/actions/runs/789012",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildRunURL(tt.serverURL, tt.repository, tt.runID)
			if got != tt.want {
				t.Errorf("buildRunURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAppendFooter(t *testing.T) {
	tests := []struct {
		name    string
		report  string
		runURL  string
		version string
		want    string
	}{
		{
			name:    "report without trailing newline and without version",
			report:  "hello world",
			runURL:  "https://github.com/kachick/nixpkgs-reviewing-workspace/actions/runs/123",
			version: "",
			want:    "hello world\n\n---\n\n[Built in this run](https://github.com/kachick/nixpkgs-reviewing-workspace/actions/runs/123)\n",
		},
		{
			name:    "report with trailing newlines and without version",
			report:  "hello world\n\n",
			runURL:  "https://github.com/kachick/nixpkgs-reviewing-workspace/actions/runs/123",
			version: "",
			want:    "hello world\n\n---\n\n[Built in this run](https://github.com/kachick/nixpkgs-reviewing-workspace/actions/runs/123)\n",
		},
		{
			name:    "report with version",
			report:  "hello world",
			runURL:  "https://github.com/kachick/nixpkgs-reviewing-workspace/actions/runs/123",
			version: "nixpkgs-review 4.0.0",
			want:    "hello world\n\n---\n\n[Built in this run](https://github.com/kachick/nixpkgs-reviewing-workspace/actions/runs/123) (nixpkgs-review 4.0.0)\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := appendFooter(tt.report, tt.runURL, tt.version)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("appendFooter() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFindReviewVersion(t *testing.T) {
	t.Run("returns empty string when file does not exist", func(t *testing.T) {
		tempDir := t.TempDir()
		got, err := findReviewVersion(tempDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "" {
			t.Errorf("findReviewVersion() = %q, want empty string", got)
		}
	})

	t.Run("finds version file in nested directory", func(t *testing.T) {
		tempDir := t.TempDir()
		subDir := filepath.Join(tempDir, "pr-123")
		if err := os.MkdirAll(subDir, 0o755); err != nil {
			t.Fatalf("failed to create sub dir: %v", err)
		}
		versionFile := filepath.Join(subDir, reviewVersionFileName)
		if err := os.WriteFile(versionFile, []byte("nixpkgs-review 4.0.0\n"), 0o644); err != nil {
			t.Fatalf("failed to write version file: %v", err)
		}

		got, err := findReviewVersion(tempDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := "nixpkgs-review 4.0.0"
		if got != want {
			t.Errorf("findReviewVersion() = %q, want %q", got, want)
		}
	})

	t.Run("skips empty version files and picks first valid", func(t *testing.T) {
		tempDir := t.TempDir()
		dirA := filepath.Join(tempDir, "a")
		dirB := filepath.Join(tempDir, "b")
		if err := os.MkdirAll(dirA, 0o755); err != nil {
			t.Fatalf("failed to create dirA: %v", err)
		}
		if err := os.MkdirAll(dirB, 0o755); err != nil {
			t.Fatalf("failed to create dirB: %v", err)
		}
		if err := os.WriteFile(filepath.Join(dirA, reviewVersionFileName), []byte("  \n"), 0o644); err != nil {
			t.Fatalf("failed to write empty version file: %v", err)
		}
		if err := os.WriteFile(filepath.Join(dirB, reviewVersionFileName), []byte("nixpkgs-review 4.0.0\n"), 0o644); err != nil {
			t.Fatalf("failed to write version file: %v", err)
		}

		got, err := findReviewVersion(tempDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := "nixpkgs-review 4.0.0"
		if got != want {
			t.Errorf("findReviewVersion() = %q, want %q", got, want)
		}
	})
}
