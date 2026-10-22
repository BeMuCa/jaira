---
id: 01M4TRAVWQFDY2W0MPEA339DNB
title: "Ein Testbuild von jaira laesst sich starten, ohne dass er Refs ins Remote sendet"
status: todo
ready: true
creator: BeMuCa
goal: "Wer einen ungeprueften Build gegen einen echten Checkout startet, kann das Senden an das Remote abschalten."
context: |-
  Das TUI sendet beim Start und bei jedem Ref-Tick alle gesammelten Refs ins Remote (internal/tui/refs.go:78, y.Flush()). Die CLI sendet ebenfalls (internal/cli/refs.go:88 und :248).
  Ein Testbuild auf einem echten Checkout veroeffentlicht also alles, was er schreibt - auch das, was ein fehlerhafter Build falsch schreibt.
  Am 2026-10-13 (D0SAHM, Milestone Roadmap) lief deshalb jeder manuelle Test auf einer Kopie: cp -a des Boards, git remote remove origin, HOME und JAIRA_HOME in ein Scratch-Verzeichnis.
  Einen Schalter gibt es heute nicht. Die Umgebungsvariablen, die der Code liest: JAIRA_DEFAULT_BOARD, JAIRA_HOME, JAIRA_LANES_DIR, JAIRA_NO_UPDATE_CHECK, JAIRA_USER (grep auf Getenv, 2026-10-13).
  Vorbild fuer den Namen: JAIRA_NO_UPDATE_CHECK.
definition-of-done: "Mit gesetztem Schalter sendet weder das TUI noch die CLI Refs; beide sagen sichtbar, dass nicht gesendet wird; ein Test belegt, dass Flush dann nicht laeuft"
tags:
  - concurrency
  - cli
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-13T13:48:19Z
updated-at: 2026-10-22T07:55:19Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-4187593
claimed-at: 2026-10-22T07:55:10Z
---

# Ein Testbuild von jaira laesst sich starten, ohne dass er Refs ins Remote sendet

## Definition of Done

- [ ] Mit gesetztem Schalter sendet weder das TUI noch die CLI Refs; beide sagen sichtbar, dass nicht gesendet wird; ein Test belegt, dass Flush dann nicht laeuft
- [ ] Der Schalter steht in docs/COMMANDS.md; go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

