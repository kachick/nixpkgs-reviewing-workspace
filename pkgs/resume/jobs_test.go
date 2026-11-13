package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseIncompleteJobs(t *testing.T) {
	data := []byte(`{
  "jobs": [
    {"name": "setup", "status": "completed", "conclusion": "success", "url": "https://example.com/1"},
    {"name": "review (x86_64-linux)", "status": "completed", "conclusion": "failure", "url": "https://example.com/2"},
    {"name": "review (aarch64-darwin (sandbox=true))", "status": "completed", "conclusion": "cancelled", "url": "https://example.com/3"},
    {"name": "review (aarch64-linux)", "status": "completed", "conclusion": "skipped", "url": "https://example.com/4"},
    {"name": "report / report", "status": "in_progress", "conclusion": "", "url": "https://example.com/5"}
  ]
}`)

	got, err := parseIncompleteJobs(data)
	if err != nil {
		t.Fatalf("parseIncompleteJobs() error = %v", err)
	}

	want := []runJob{
		{Name: "review (x86_64-linux)", Status: "completed", Conclusion: "failure", URL: "https://example.com/2"},
		{Name: "review (aarch64-darwin (sandbox=true))", Status: "completed", Conclusion: "cancelled", URL: "https://example.com/3"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("parseIncompleteJobs() mismatch (-want +got):\n%s", diff)
	}
}

func TestAppendIncompleteJobs(t *testing.T) {
	report := "hello world\n\n---\n\n[Built in this run](https://example.com/run)\n"

	t.Run("keeps report as is without incomplete jobs", func(t *testing.T) {
		if diff := cmp.Diff(report, appendIncompleteJobs(report, nil)); diff != "" {
			t.Errorf("appendIncompleteJobs() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("lists incomplete jobs", func(t *testing.T) {
		jobs := []runJob{
			{Name: "review (aarch64-darwin (sandbox=true))", Status: "completed", Conclusion: "cancelled", URL: "https://example.com/3"},
		}
		want := "hello world\n\n---\n\n[Built in this run](https://example.com/run)\n\nMissing results from these jobs:\n\n- [review (aarch64-darwin (sandbox=true))](https://example.com/3): cancelled\n"
		if diff := cmp.Diff(want, appendIncompleteJobs(report, jobs)); diff != "" {
			t.Errorf("appendIncompleteJobs() mismatch (-want +got):\n%s", diff)
		}
	})
}
