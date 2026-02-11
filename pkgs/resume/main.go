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

func getCurrentAssetName() string {
	var a, o string
	switch runtime.GOARCH {
	case "amd64":
		a = "X64"
	case "arm64":
		a = "ARM64"
	}
	switch runtime.GOOS {
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

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <run_id>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	runID := args[0]

	isCI := os.Getenv("CI") == "true"
	currentAssetName := getCurrentAssetName()

	if !isCI {
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

	fmt.Fprintf(os.Stderr, "Downloading artifacts to %s...\n", tempDir)
	downloadCmd := exec.Command("gh", "run", "download", runID, "--dir", tempDir)
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
}
