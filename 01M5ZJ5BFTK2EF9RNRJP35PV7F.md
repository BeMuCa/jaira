---
id: 01M5ZJ5BFTK2EF9RNRJP35PV7F
title: "jaira milestone create ohne --after sagt, dass der Milestone nicht auf der Roadmap steht"
status: in-progress
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
updated-at: 2026-11-23T01:21:16Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: worker-35pv7f
claimed-at: 2026-11-23T01:19:53Z
---

# jaira milestone create ohne --after sagt, dass der Milestone nicht auf der Roadmap steht

## Definition of Done

- [x] jaira milestone create <name> ohne --after druckt einen Satz mit dem set --after Befehl
  proof: internal/cli/milestones.go:207 (printed only when --after was not given); TestCreateWritesGoalAndAfter in internal/cli/milestones_test.go; manual run: 'milestone create q3' prints the hint, '--after q1', '--after=' and --json do not, and the suggested 'milestone set q3 --after q1' succeeds
- [ ] Release-Notiz in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 01:20 · BeMuCa** — Checked branch build 37e27ef69: 'milestone create q1' without --after prints only 'Milestone "q1" created with 0 ticket(s): <path>' — no roadmap hint, so the change is needed (not a reject). Decisions settled without the person: hint goes to stdout as a second line after the created line, only when --after was not given at all (--after "" does put it on the roadmap as a swimlane top); --json output unchanged (no message field to add to); the hint names the milestone by its name — create refuses a name already on the board, so the name is unique there.
