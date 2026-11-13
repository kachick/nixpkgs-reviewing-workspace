package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type runJob struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	URL        string `json:"url"`
}

// parseIncompleteJobs returns finished jobs that did not succeed.
// Jobs still running are ignored because the report job itself is running when called inside the workflow.
// Matching on the "review" job name is avoided to keep this working after the workflow renames jobs.
func parseIncompleteJobs(data []byte) ([]runJob, error) {
	var run struct {
		Jobs []runJob `json:"jobs"`
	}
	if err := json.Unmarshal(data, &run); err != nil {
		return nil, err
	}

	var incomplete []runJob
	for _, job := range run.Jobs {
		if job.Status != "completed" {
			continue
		}
		switch job.Conclusion {
		case "success", "skipped":
			continue
		}
		incomplete = append(incomplete, job)
	}
	return incomplete, nil
}

func fetchIncompleteJobs(runID string) ([]runJob, error) {
	out, err := exec.Command("gh", "run", "view", runID, "--json", "jobs").Output()
	if err != nil {
		return nil, err
	}
	return parseIncompleteJobs(out)
}

func appendIncompleteJobs(report string, jobs []runJob) string {
	if len(jobs) == 0 {
		return report
	}

	var builder strings.Builder
	builder.WriteString(strings.TrimRight(report, "\n"))
	builder.WriteString("\n\nMissing results from these jobs:\n\n")
	for _, job := range jobs {
		fmt.Fprintf(&builder, "- [%s](%s): %s\n", job.Name, job.URL, job.Conclusion)
	}
	return builder.String()
}
