#!/usr/bin/env bash

set -euxo pipefail

pr_number="$1"
output_dir="$(mktemp --directory)"

gh workflow run 'nixpkgs-review.yml' --field pr-number="$pr_number"
sleep 10 # TODO: Ensure to get correct ID
run_id="$(gh run list --workflow=nixpkgs-review.yml --branch main --limit 1 --json databaseId --jq '.[].databaseId')"

gh run watch "$run_id" # Don't use --exit-status to make sure the downloading
gh run download "$run_id" --dir "$output_dir"

echo "Downloaded the files in $output_dir"

tree "$output_dir"

fd --absolute-path . "$output_dir"

# TODO: Sort the results with nixpkgs Tier - https://github.com/NixOS/nixpkgs/blob/nixos-24.11/lib/systems/flake-systems.nix
# TODO: Trim excess headers
fd --absolute-path report.md "$output_dir" | xargs cat
