# Deferred items — 260823-khj

Out-of-scope findings noticed during execution, not fixed (per the scope
boundary rule: only auto-fix issues directly caused by this task's changes).

- `internal/cli/tickets.go` has a pre-existing `gofmt` drift (two map-literal
  alignment blocks are not gofmt-clean) unrelated to any file this quick task
  touched. `git log -1 -- internal/cli/tickets.go` shows it was already in
  that state before this task started. Left as-is; a future pass touching
  that file should run `gofmt -w` on it.
