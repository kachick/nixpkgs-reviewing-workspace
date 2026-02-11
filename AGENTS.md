# Project Rules for AI Agents

## Languages & Frameworks
- **Go**: Use standard libraries as much as possible.
  - Formatter: `gofumpt`
  - Testing: Use `github.com/google/go-cmp/cmp` for deep equality.
  - CI/Lint: `go vet` and `gofumpt -d` are mandatory.
- **Nix**: Adhere to existing patterns in `flake.nix` and `pkgs/`.

## Git & Workflow
- **Commits**: Create new, small commits for each logical change. Do not `amend` unless specifically asked. This helps the maintainer review the history before squash-merging.
- **Messages**: Follow the style of recent commits. Use `Assisted-by: Gemini <gemini@google.com>` for significant AI-generated changes.

## Coding Style
- **Comments**: Minimal. Focus on *why*, not *what*. Delete obvious comments like `// Verify files exist`.
- **UI/UX**: Maintain rich terminal output (e.g., connecting `Stdin/Stdout` for interactive tools like `gh`).
- **Tests**: Prefer real data snapshots over mocks where feasible.
