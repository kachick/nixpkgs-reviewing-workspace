# nixpkgs-reviewing-workspace

[![Report](https://github.com/kachick/nixpkgs-reviewing-workspace/actions/workflows/report.yml/badge.svg?branch=main)](https://github.com/kachick/nixpkgs-reviewing-workspace/actions/workflows/report.yml?query=branch%3Amain+)
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

As I understand that, Which github provided runner matches the Nix supported [systems](https://github.com/NixOS/nixpkgs/blob/nixos-24.11/lib/systems/flake-systems.nix).

Excluding _large_ runners here. Which cannot be used in free plan.

| Nix                                                                               | GitHub           |
| --------------------------------------------------------------------------------- | ---------------- |
| x86_64-linux                                                                      | ubuntu-24.04     |
| [aarch64-Linux](https://github.com/kachick/nixpkgs-reviewing-workspace/issues/43) | ubuntu-24.04-arm |
| x86_64-darwin                                                                     | macos-13         |
| aarch64-darwin                                                                    | macos-15         |

## Usage

Specify the [NixOS/nixpkgs PR](https://github.com/NixOS/nixpkgs/pulls) number with following CLI,

```bash
task start -- "$PR_NUMBER"
```

or input it in [workflow_dispatch](https://github.com/kachick/nixpkgs-reviewing-workspace/actions/workflows/nixpkgs-review.yml) if you can only use a smartphone.

![Image](https://github.com/user-attachments/assets/2fd03f40-7561-4c48-a35e-ed9ba309ac5f)

If you want to get results and reports in another machine, run below

```bash
task resume
```

## Traps

I have obtained a different result with the runners and my device once.\
I suspect it is caused by pnpm spec, but you should remember this case, especially if the package is using `pnpm.fetchDeps`.

- <https://github.com/NixOS/nixpkgs/pull/361460/files#r1906545103>
