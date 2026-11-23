---
id: 01M5ZJ5BFTK2EF9RNRJP35PV7F
title: "jaira milestone create ohne --after sagt, dass der Milestone nicht auf der Roadmap steht"
status: testing
ready: true
creator: BeMuCa
goal: "Legt jemand einen Milestone ohne --after an, sagt die Ausgabe, dass er im Picker steht, nicht auf der Roadmap, und wie er dorthin kommt (jaira milestone set <name> --after ...)."
context: |-
  Was schief ist: ein Milestone ohne --after landet im M-Picker, nicht auf der Roadmap; der Mensch schaute auf die Roadmap und hielt ihn fuer verloren (O14, 27.10.2026).
  Heute steht das nur in jaira milestone create --help ('Without --after the milestone stays off the roadmap'), nicht in der Ausgabe des Befehls.
  Ein Satz in der Ausgabe reicht.
definition-of-done: "jaira milestone create <name> ohne --after druckt einen Satz mit dem set --after Befehl"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:52:21Z
updated-at: 2026-11-23T01:22:55Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: worker-35pv7f
claimed-at: 2026-11-23T01:19:53Z
outcome-what: "milestone create without --after now prints a second line: the milestone is in the picker (M), not on the roadmap (L), and 'jaira milestone set <name> --after <a,b>' puts it there. --after given (also --after \"\") and --json print nothing new. Test assertions added to TestCreateWritesGoalAndAfter; release note added."
outcome-why: A milestone made without --after went to the M picker only; the person looked on the roadmap and thought it lost (O14). The hint was only in --help.
outcome-resolves: "DoD 1 (hint with the set --after command) and DoD 2 (release note). Checked: concurrency (no new write, the line prints after the locked write returned), error paths (an error returns before the line), unusual input (normalised name like 'Round Two' -> round-two is the one printed; --after \"\" counts as placed), permissions n/a. go test not run per worker rules."
review-summary: "none"
round: 2
back-from: critique
back-by: lane
back-count: 1
review-gaps: "none — the change is one if-block (internal/cli/milestones.go:206-208) reusing the command's own writer and name; nothing duplicated, nothing unused. The test was folded into a table in TestCreateWritesGoalAndAfter rather than a new test, so no new helper."
test-verdict: "pass: gofmt clean, go vet ./internal/cli RC=0, go build ./... RC=0; go test NOT run (house rule tonight: one suite before the PR, run by the main session) — TestCreateWritesGoalAndAfter is unverified until then. DoD 1 verified on a scratch board: 'milestone create q3' prints the hint line, '--after q1', '--after=' and --json do not, and the hinted 'milestone set q3 --after q1' succeeds (after: [XM44T3] in the file); 'Round Two' hints round-two. DoD 2: core/release/notes/35PV7F.md is one unwrapped '- ' line. Failure paths: duplicate name exit 3, bad name exit 2, unknown --after exit 5 with no file written, none print the hint. Two parallel 'create par': one exit 3, one created + hint, one file."
---

# jaira milestone create ohne --after sagt, dass der Milestone nicht auf der Roadmap steht

## Definition of Done

- [x] jaira milestone create <name> ohne --after druckt einen Satz mit dem set --after Befehl
  proof: internal/cli/milestones.go:207 (printed only when --after was not given); TestCreateWritesGoalAndAfter in internal/cli/milestones_test.go; manual run: 'milestone create q3' prints the hint, '--after q1', '--after=' and --json do not, and the suggested 'milestone set q3 --after q1' succeeds
- [x] Release-Notiz in core/release/notes/<ID>.md
  proof: core/release/notes/35PV7F.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 01:20 · BeMuCa** — Checked branch build 37e27ef69: 'milestone create q1' without --after prints only 'Milestone "q1" created with 0 ticket(s): <path>' — no roadmap hint, so the change is needed (not a reject). Decisions settled without the person: hint goes to stdout as a second line after the created line, only when --after was not given at all (--after "" does put it on the roadmap as a swimlane top); --json output unchanged (no message field to add to); the hint names the milestone by its name — create refuses a name already on the board, so the name is unique there.
- **2026-11-23 01:21 · BeMuCa** — critique: core/release/notes/35PV7F.md:1 says 'read the second line', but with a normalised name ('Round Two') a 'Filed ... as ...' line comes first and the hint is the third; say 'the line under the created line' instead. Rest let stand: one if-block after the existing created line is the smallest shape, stdout matches the other create output, the TUI create path (internal/tui/layer.go createMilestone) always places on the roadmap so it has no sibling.
- **2026-11-23 01:22 · BeMuCa** — bughunt: done at 46f276515 - 0 findings (0 change, 0 pre-existing) in 1 module (internal/cli milestone create), first hunt; concurrency, error paths, input, permissions checked; files: core/release/notes/35PV7F.md internal/cli/milestones.go internal/cli/milestones_test.go
