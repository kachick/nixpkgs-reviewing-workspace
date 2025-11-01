# nixpkgs-reviewing-workspace

[![Reviewing](https://github.com/kachick/nixpkgs-reviewing-workspace/actions/workflows/nixpkgs-review.yml/badge.svg?branch=main)](https://github.com/kachick/nixpkgs-reviewing-workspace/actions/workflows/nixpkgs-review.yml?query=branch%3Amain+)
[![Dev Shell](https://github.com/kachick/nixpkgs-reviewing-workspace/actions/workflows/devshell.yml/badge.svg?branch=main)](https://github.com/kachick/nixpkgs-reviewing-workspace/actions/workflows/devshell.yml?query=branch%3Amain+)

## Motivation

Checking the build result of nixpkgs in multiple platforms is an annoying task.

- I basically do not use macOS these days and only have an old and slow Intel MacBook.
- Having few devices or smartphones if I am walking.
- Several packages are very heavy to build, they take a long time and require many resources.

Solution

- GitHub Actions are free for pricing in public repositories. 😋

## Matrix

The GitHub hosted runners for Nix supported [systems](https://github.com/NixOS/nixpkgs/blob/nixos-25.05/lib/systems/flake-systems.nix) are in the list below.\
Large runners are excluded because they are not available on the free plan.

| Nix            | GitHub                                                                                                       |
| -------------- | ------------------------------------------------------------------------------------------------------------ |
| x86_64-linux   | [ubuntu-24.04](https://github.com/actions/runner-images/issues/9848)                                         |
| aarch64-linux  | [ubuntu-24.04-arm](https://github.com/actions/partner-runner-images/blob/main/images/arm-ubuntu-24-image.md) |
| x86_64-darwin  | [macos-15-intel](https://github.com/actions/runner-images/issues/13045)                                      |
| aarch64-darwin | [macos-26](https://github.com/actions/runner-images/issues/13008)                                            |

## Usage

### nixpkgs-review

Specify the [NixOS/nixpkgs PR](https://github.com/NixOS/nixpkgs/pulls) number with following CLI,

```bash
nix run .#review -- "$PR_NUMBER"
```

or input it in [workflow_dispatch](https://github.com/kachick/nixpkgs-reviewing-workspace/actions/workflows/nixpkgs-review.yml) if you can only use a smartphone.

![Image](https://github.com/user-attachments/assets/2fd03f40-7561-4c48-a35e-ed9ba309ac5f)

If you want to get results and reports in another machine, run below

```bash
nix run .#fzf-resume
```

### nix-build / nix build

Faster feedbacks much help us while early stage of adding or updating package.\
You can use [another workflow](https://github.com/kachick/nixpkgs-reviewing-workspace/actions/workflows/build.yml) for that purpose even if not yet sending a PR.

## Traps

I observed differing results between GHA runners, my local device, and Hydra.

- [Different hashes on my device and GHA runners when using `pnpm.fetchDeps`](https://github.com/NixOS/nixpkgs/pull/361460/files#r1906545103)
- [Failed post-merge, even though macOS runners succeeded](https://github.com/NixOS/nixpkgs/pull/382541#issuecomment-2670547003)
