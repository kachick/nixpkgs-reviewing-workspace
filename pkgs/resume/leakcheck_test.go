package main

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestBetterleaksChecker(t *testing.T) {
	if _, err := exec.LookPath("betterleaks"); err != nil {
		t.Skip("betterleaks executable not found in PATH")
	}

	ctx := context.Background()
	checker := &BetterleaksChecker{}

	t.Run("clean text passes", func(t *testing.T) {
		err := checker.Check(ctx, "Generated review report for package")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("detects secret pattern and redacts value from error", func(t *testing.T) {
		// Use dummy GitHub PAT pattern documented in upstream docs
		rawSecret := "ghp_aB3dE5fG7hI9jK1mN3pQ5rS7tU9vW1xY3zA5"
		dummyToken := "GITHUB_TOKEN=" + rawSecret
		err := checker.Check(ctx, dummyToken)
		if err == nil {
			t.Fatal("expected error for detected secret, got nil")
		}
		if !errors.Is(err, ErrLeakDetected) {
			t.Fatalf("expected ErrLeakDetected, got: %v", err)
		}
		if strings.Contains(err.Error(), rawSecret) {
			t.Errorf("error message must redact raw secret, got: %v", err)
		}
	})
}
