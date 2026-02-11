package main

import (
	"regexp"
	"sort"
	"strings"
)

type NixSystem struct {
	Platform string
	Tier     int
}

type Runner struct {
	Arch string
	OS   string
}

func (r Runner) AssetName() string {
	return r.Arch + "-" + r.OS
}

var NixByAssetName = map[string]NixSystem{
	"X64-Linux":   {Platform: "x86_64-linux", Tier: 1},
	"ARM64-Linux": {Platform: "aarch64-linux", Tier: 2},
	"X64-macOS":   {Platform: "x86_64-darwin", Tier: 2},
	"ARM64-macOS": {Platform: "aarch64-darwin", Tier: 3},
}

var assetRegex = regexp.MustCompile(`nixpkgs-review-files-.*`)

func extractAssetName(path string) string {
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if assetRegex.MatchString(part) {
			for assetName := range NixByAssetName {
				if strings.Contains(part, assetName) {
					return assetName
				}
			}
		}
	}
	return ""
}

func getPriority(path string) (int, int) {
	assetName := extractAssetName(path)
	nix, ok := NixByAssetName[assetName]
	if !ok {
		return 999, 999 // Unknown
	}

	tier := nix.Tier
	myFavor := 0
	if !strings.HasSuffix(assetName, "Linux") {
		myFavor = 1 // Simplified: 0 for Linux, 1 for others (macOS)
	}

	return tier, myFavor
}

func SortPaths(paths []string) {
	sort.Slice(paths, func(i, j int) bool {
		tI, fI := getPriority(paths[i])
		tJ, fJ := getPriority(paths[j])

		if tI != tJ {
			return tI < tJ
		}
		return fI < fJ
	})
}
