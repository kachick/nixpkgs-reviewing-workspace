#!/usr/bin/env bash

set -euo pipefail

if [[ "$CI" != 'true' ]]; then
  set -x
fi

run_id="$1"
output_dir="$(mktemp --tmpdir --directory "nixpkgs-reviewing-workspace.run-${run_id}.XXXX")"

# https://github.blog/changelog/2020-04-15-github-actions-sets-the-ci-environment-variable-to-true/
if [[ "$CI" != 'true' ]]; then
  # Replace with like a `gh workflow run --watch` if https://github.com/cli/cli/issues/3559 is resolved
  gh run watch "$run_id" --interval 10 # Don't use --exit-status to make sure the downloading
fi

{
  gh run download "$run_id" --dir "$output_dir"

  echo "Downloaded the files in $output_dir"

  tree "$output_dir"

  fd --absolute-path . "$output_dir"
} >&2

fd --absolute-path report.md "$output_dir" | ruby ./sort.rb | ruby ./concat.rb
