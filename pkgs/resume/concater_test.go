package main

import (
	"os"
	"strings"
	"testing"
)

func TestConcatReports_RealData(t *testing.T) {
	paths := []string{
		"snapshots/nixpkgs-review-files-pr-489308-X64-Linux/pr-489308/report.md",
		"snapshots/nixpkgs-review-files-pr-489308-ARM64-Linux/pr-489308/report.md",
		"snapshots/nixpkgs-review-files-pr-489308-X64-macOS/pr-489308/report.md",
		"snapshots/nixpkgs-review-files-pr-489308-ARM64-macOS/pr-489308/report.md",
	}

	for _, p := range paths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Skipf("Snapshot file %s not found, skipping integration test", p)
		}
	}

	got, err := ConcatReports(paths)
	if err != nil {
		t.Fatalf("ConcatReports failed: %v", err)
	}

	if !strings.Contains(got, "## `nixpkgs-review` result") {
		t.Error("Result should contain the primary header from the first file")
	}

	count := strings.Count(got, "---")
	if count != 4 {
		t.Errorf("Expected 4 '---' markers, got %d", count)
	}

	secondContent, err := os.ReadFile(paths[1])
	if err != nil {
		t.Fatalf("Failed to read second file: %v", err)
	}

	secondParts := strings.SplitN(string(secondContent), "---", 2)
	if len(secondParts) < 2 {
		t.Fatalf("Second file %s does not contain '---' separator", paths[1])
	}
	metadata := secondParts[0]

	if metadata != "" {
		count := strings.Count(got, metadata)
		if count != 1 {
			t.Errorf("Metadata from second file should have been trimmed. Expected 1 occurrence (from first file), got %d.\nMetadata: %q", count, metadata)
		}
	}

	body := "---" + secondParts[1]
	if !strings.Contains(got, body) {
		t.Error("Body from the second file should be present in the result")
	}
}
