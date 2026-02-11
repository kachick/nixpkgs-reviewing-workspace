# Rules for AI Agents

## Reference Projects

- **General Structure**: Follow patterns in [kachick/anylang-template](https://github.com/kachick/anylang-template) and [kachick/dotfiles](https://github.com/kachick/dotfiles).
- **Go Style**: Follow patterns in [kachick/selfup](https://github.com/kachick/selfup).

## Languages

- **Go**: Use standard libraries as much as possible.
  - Formatter: Use `gofumpt`.
  - Testing: Use `github.com/google/go-cmp/cmp`.
  - Check: `go vet` and `gofumpt -d` are required.
- **Nix**: Follow patterns in `flake.nix` and `pkgs/`.

## Git

- **Verification**: Run `task check` before you commit. All tests and lints must pass.
- **Commits**: Make small commits for each change. Do not use `commit --amend`.
- **Messages**: Follow the style of recent commits. Always add `Assisted-by: Gemini <gemini@google.com>` at the end of every commit message you create.

## Coding Style

- **Comments**: Write only important comments. Do not write obvious things.
- **UI**: Keep good terminal output (e.g., connect `os.Stdin/Stdout` for `gh`).
- **Tests**: Use real data instead of mocks when you can.
