# How to develop

## Supported host environment

The goal of this project is to test builds on multiple platforms that differ from the host. However, the local tools currently assume `x86_64-linux` as the host OS to match the author's personal setup.

Rather than hiding unexpected patterns or trying to support every host architecture, the tools should report clear errors when running on unexpected host platforms.

## Pull requests

This repository does not accept pull requests from outside for now.

- **Personal use**: As stated in [`9efbf524d95d2c03d245c0864e528232718578a6`](https://github.com/kachick/nixpkgs-reviewing-workspace/commit/9efbf524d95d2c03d245c0864e528232718578a6), this repository is kept only for personal use and quick customization.
- **Security reasons**: The author temporarily pauses external pull requests on repositories that provide GitHub Actions. A small mistake in GitHub Actions settings can cause security issues, and handling spam PRs creates [noise](https://github.com/kachick/wait-other-jobs/issues/1402).

## Note

- nixpkgs-review supports shallow clone since <https://github.com/Mic92/nixpkgs-review/issues/267> => <https://github.com/Mic92/nixpkgs-review/commit/055465e55d131ffb1e1617f46d3bade0b87bbe69>
