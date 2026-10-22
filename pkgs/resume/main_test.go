package main

import (
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
