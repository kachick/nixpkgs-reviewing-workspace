# nixpkgs-reviewing-workspace

[![Dev Shell Status](https://github.com/kachick/nixpkgs-reviewing-workspace/actions/workflows/devshell.yml/badge.svg?branch=main)](https://github.com/kachick/nixpkgs-reviewing-workspace/actions/workflows/devshell.yml?query=branch%3Amain+)

## Motivation

Checking the build result of nixpkgs in multiple platforms is an annoying task.

- I basically do not use macOS these days and only have an old and slow Intel MacBook.
- Having few devices or smartphones if I am walking.
- Several packages are very heavy to build, they take a long time and require many resources.

Solution

- GitHub Actions are free for pricing in public repositories. 😋

## Matrix

As I understand that, Which github provided runner matches the Nix supported [systems](https://github.com/NixOS/nixpkgs/blob/nixos-24.11/lib/systems/flake-systems.nix).

Excluding _large_ runners here. Which cannot be used in free plan.

| Nix            | GitHub           |
| -------------- | ---------------- |
| x86_64-linux   | ubuntu-24.04     |
| aarch64-Linux  | ubuntu-24.04-arm |
| x86_64-darwin  | macos-13         |
| aarch64-darwin | macos-15         |

## Usage

```bash
./review.bash "$PR_NUMBER"
```
