---
phase: quick-260819-hco
plan: 01
subsystem: docs
tags: [jaira-cli, skill, agents-md, ticket-lifecycle]

requires: []
provides:
  - "Pre-create conflict check documented in .claude/skills/jaira/SKILL.md"
  - "Same rule documented in docs/AGENTS.md for non-Claude agents"
affects: [jaira-skill, jaira-docs]

tech-stack:
  added: []
  patterns:
    - "Search-before-write: jaira list -q / jaira show --json before jaira create, ask-the-user gate on genuine contradiction"

key-files:
  created: []
  modified:
    - .claude/skills/jaira/SKILL.md
    - docs/AGENTS.md

key-decisions:
  - "Described jaira list -q's real match fields (title, goal, context, DoD, assignee, status) per tickets.go matches(), not the narrower --help one-liner"
  - "Reused --follows for deliberate supersession rather than inventing a new flag"
  - "Kept 'What an agent cannot do' in AGENTS.md untouched — this rule is a judgment call, not a gate-enforced refusal"

patterns-established:
  - "Pattern: pre-write conflict check — search, read close hits, ask on contradiction, name handles/ids and quote the line"

requirements-completed: [QUICK-260819-hco]

duration: 12min
completed: 2026-08-19
---

# Quick Task 260819-hco: Pre-create conflict check Summary

**Taught agents (Claude and any shell-capable agent) to search the board with `jaira list -q` / `jaira show --json` before `jaira create`, and to stop and ask the user — naming both handles and quoting the contradicting line — when an existing ticket already decided the same question differently.**

## Performance

- **Duration:** 12 min
- **Started:** 2026-08-19T10:25:00Z (approx)
- **Completed:** 2026-08-19T10:37:14Z
- **Tasks:** 2 (1 verification-only, 1 doc edit + commit)
- **Files modified:** 2

## Accomplishments
- New `## Before you create, check what the board already decided` section in `.claude/skills/jaira/SKILL.md`, placed between `## When to create tickets` and `## Writing a good ticket from a request`, in the file's second-person/terse/example-driven voice.
- New `## Before creating a ticket` section in `docs/AGENTS.md`, placed between `## The loop` and `## Leaving a trail`, in the file's more explanatory register, explicitly noting nothing in the binary enforces this.
- Every command quoted (`jaira list -q`, `jaira show <handle/id> --json`, `--follows`) verified against the real CLI (`--help` output and `internal/cli/tickets.go`) before being written.

## Task Commits

1. **Task 1: Confirm the CLI surface** - verification-only, no commit (ran `go run ./cmd/jaira list --help`, `go run ./cmd/jaira show --help`, read `internal/cli/tickets.go:420-428` `matches()`)
2. **Task 2: Land the pre-create conflict check in both files** - `d280620` (docs)

**Plan metadata:** commit deferred to orchestrator per constraints (SUMMARY.md/STATE.md not committed by this executor).

## Files Created/Modified
- `.claude/skills/jaira/SKILL.md` - Added pre-create conflict-check section (27 lines) with a `jaira list -q` / `jaira show --json` command block, the three ways forward (honor / `--follows` supersede / drop), and the duplicate case.
- `docs/AGENTS.md` - Added the same rule (24 lines) in a more explanatory register for any shell-capable agent, with an explicit "nothing in the binary enforces this" note; `## What an agent cannot do` list left untouched.

## Decisions Made
- Described `-q`'s real search fields (title, goal, context, DoD, assignee, status — from `matches()` at `internal/cli/tickets.go:422`) rather than the `--help` string's narrower "title, goal and id" summary, per the plan's explicit instruction.
- No new CLI flags, no Go changes, no lane-file changes — reused `--follows` exactly as it already exists.

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

Both files carry the rule; no further work implied by this quick task. Verification gates (help output, diff scope, section placement, command existence) all passed — see Self-Check below.

## Self-Check

- `go run ./cmd/jaira list --help 2>&1 | grep -q -- '-q, --query'` → PASSED
- `go run ./cmd/jaira show --help 2>&1 | grep -q -- '--json'` → PASSED
- `grep -q 'jaira list -q' .claude/skills/jaira/SKILL.md` → PASSED
- `grep -q 'jaira show' .claude/skills/jaira/SKILL.md` → PASSED
- `grep -q -- '--follows' .claude/skills/jaira/SKILL.md` → PASSED
- `grep -q 'jaira list -q' docs/AGENTS.md` → PASSED
- `grep -q -- '--follows' docs/AGENTS.md` → PASSED
- `git diff --stat HEAD~1` → only `.claude/skills/jaira/SKILL.md` and `docs/AGENTS.md` changed (51 insertions, 0 deletions), confirmed
- New sections sit in the required locations in both files (confirmed via `grep -n '^## '`)
- `## What an agent cannot do` in `docs/AGENTS.md` unchanged (confirmed no diff touched that block)
- Commit `d280620` exists in `git log --oneline -1`: confirmed

## Self-Check: PASSED

---
*Phase: quick-260819-hco*
*Completed: 2026-08-19*
