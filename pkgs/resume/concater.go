package main

import (
	"os"
	"strings"
)

func ConcatReports(paths []string) (string, error) {
	var builder strings.Builder
	for i, path := range paths {
		b, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		content := string(b)

		// nixpkgs-review report.md is separated by "---" into common metadata (header) and per-platform results (body).
		parts := strings.SplitN(content, "---", 2)
		header := parts[0]

		if i == 0 {
			// Keep the header (metadata) from the first file only.
			builder.WriteString(header)
		}

		if len(parts) > 1 {
			body := parts[1]
			if i > 0 {
				// Ensure there's a clear separation between platform sections.
				builder.WriteByte('\n')
			}
			// Append the separator and the platform-specific body.
			builder.WriteString("---")
			builder.WriteString(body)
		}
	}
	return builder.String(), nil
}
