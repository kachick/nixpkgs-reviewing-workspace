package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func assetNameFor(arch, osName string) string {
	var a, o string
	switch arch {
	case "amd64":
		a = "X64"
	case "arm64":
		a = "ARM64"
	}
	switch osName {
	case "linux":
		o = "Linux"
	case "darwin":
		o = "macOS"
	}

	if a == "" || o == "" {
		return ""
	}
	return a + "-" + o
}

func getCurrentAssetName() string {
	return assetNameFor(runtime.GOARCH, runtime.GOOS)
}

func determineDownloadArgs(runID, tempDir, currentAssetName string, all bool) ([]string, bool) {
	if all {
		return []string{"run", "download", runID, "--dir", tempDir}, false
	}

	if currentAssetName != "" {
		return []string{"run", "download", runID, "--dir", tempDir, "--pattern", "*" + currentAssetName + "*"}, false
	}

	return []string{"run", "download", runID, "--dir", tempDir}, true
}

func main() {
	var allArtifacts bool
	flag.BoolVar(&allArtifacts, "all", false, "download all architecture artifacts")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <run_id>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	runID := args[0]

	// Check the value for clarity and robustness.
	// Some other tools check only for existence, so don't define it in .envrc even as "false".
	// For example, Gemini CLI UI degraded when CI=false was set.
	inRunner := os.Getenv("GITHUB_ACTIONS") == "true"
	currentAssetName := getCurrentAssetName()

	if !inRunner {
		cmd := exec.Command("gh", "run", "watch", runID, "--interval", "10")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}

	tempDir, err := os.MkdirTemp("", "nixpkgs-reviewing-workspace.run-"+runID+".*")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}

	downloadArgs, isUnknownArch := determineDownloadArgs(runID, tempDir, currentAssetName, allArtifacts || inRunner)
	if isUnknownArch {
		fmt.Fprintf(os.Stderr, "Warning: unknown host architecture/OS (%s/%s). Falling back to downloading all artifacts.\n", runtime.GOARCH, runtime.GOOS)
	}

	fmt.Fprintf(os.Stderr, "Downloading artifacts to %s...\n", tempDir)
	downloadCmd := exec.Command("gh", downloadArgs...)
	downloadCmd.Stderr = os.Stderr
	if err := downloadCmd.Run(); err != nil {
		log.Fatalf("Failed to download artifacts: %v", err)
	}

	fmt.Fprintln(os.Stderr, "Downloaded files:")
	err = filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			absPath, err := filepath.Abs(path)
			if err != nil {
				absPath = path
			}

			displayPath := absPath
			if currentAssetName != "" && extractAssetName(absPath) == currentAssetName {
				displayPath = "\x1b[32m" + absPath + "\x1b[0m"
			}
			fmt.Fprintf(os.Stderr, "  %s\n", displayPath)
		}
		return nil
	})
	if err != nil {
		log.Printf("Error listing files: %v", err)
	}

	var reports []string
	err = filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Base(path) == "report.md" {
			reports = append(reports, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to find reports: %v", err)
	}

	if len(reports) == 0 {
		log.Fatalf("No report.md found in %s", tempDir)
	}

	SortPaths(reports)

	finalReport, err := ConcatReports(reports)
	if err != nil {
		log.Fatalf("Failed to concatenate reports: %v", err)
	}

	fmt.Println()
	fmt.Println(finalReport)

	if isUnknownArch {
		fmt.Fprintf(os.Stderr, "Error: host architecture/OS (%s/%s) is unrecognized; downloaded all artifacts as fallback.\n", runtime.GOARCH, runtime.GOOS)
		os.Exit(1)
	}
}
