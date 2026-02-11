package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSortPaths(t *testing.T) {
	paths := []string{
		"snapshots/nixpkgs-review-files-pr-489308-ARM64-macOS/pr-489308/report.md",
		"snapshots/nixpkgs-review-files-pr-489308-X64-Linux/pr-489308/report.md",
		"snapshots/nixpkgs-review-files-pr-489308-ARM64-Linux/pr-489308/report.md",
		"snapshots/nixpkgs-review-files-pr-489308-X64-macOS/pr-489308/report.md",
	}

	expected := []string{
		"snapshots/nixpkgs-review-files-pr-489308-X64-Linux/pr-489308/report.md",   // Tier 1, Linux
		"snapshots/nixpkgs-review-files-pr-489308-ARM64-Linux/pr-489308/report.md", // Tier 2, Linux
		"snapshots/nixpkgs-review-files-pr-489308-X64-macOS/pr-489308/report.md",   // Tier 2, macOS
		"snapshots/nixpkgs-review-files-pr-489308-ARM64-macOS/pr-489308/report.md", // Tier 3, macOS
	}

	SortPaths(paths)

	if diff := cmp.Diff(expected, paths); diff != "" {
		t.Errorf("SortPaths() mismatch (-want +got):\n%s", diff)
	}
}
