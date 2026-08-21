# Deferred Items — quick task 260821-t7s

Out-of-scope discoveries found while executing this task. Not fixed here per
the scope boundary rule (only auto-fix issues directly caused by this task's
own changes).

| File | Issue | Notes |
|------|-------|-------|
| `core/gate/gate.go` | `gofmt -l` reports this file as not gofmt-clean | Pre-existing before this task (confirmed via `git stash` against the branch tip); unrelated to `core/selfupdate` or `internal/cli/self.go`. Task 1's `<verify>` runs `gofmt -l core internal` repo-wide, which surfaces this even though it was not introduced here. |
| `internal/cli/tickets.go` | `gofmt -l` reports this file as not gofmt-clean | Same as above — pre-existing, unrelated to this task's files. |
