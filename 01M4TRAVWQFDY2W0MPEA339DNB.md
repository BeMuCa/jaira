---
id: 01M4TRAVWQFDY2W0MPEA339DNB
title: "Ein Testbuild von jaira laesst sich starten, ohne dass er Refs ins Remote sendet"
status: in-progress
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
updated-at: 2026-10-22T08:01:50Z
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

- [x] Tests zuerst: core/gitref - mit JAIRA_NO_PUSH=1 verweigern Write und PushBranch mit ErrPushOff, das Remote bleibt unveraendert; core/refsync - Flush mit Schalter gibt ErrPushOff, sendet nichts, der Eintrag bleibt in der Outbox und geht ohne Schalter danach raus
- [x] core/gitref: PushDisableEnv = JAIRA_NO_PUSH, ErrPushOff, PushDisabled(); pushRefspec und PushBranch (die einzigen zwei git push im Binary) brechen davor ab - deckt pull, release, snapshot und Epic-Rebase ohne weitere Stellen
- [x] core/refsync Syncer.Flush: bei PushDisabled vor dem Fetch nil, gitref.ErrPushOff zurueck - die Queue bleibt wie offline, nichts geht verloren
- [ ] CLI: flushRefs (internal/cli/refs.go:250) gibt die eine stderr-Zeile aus; fileOnRefOnly/putOnRef reichen ErrPushOff als why durch; webPush zeigt es
- [ ] TUI internal/tui/refs.go:85: ErrPushOff aus Flush ist kein Fehler, der Fetch laeuft weiter; Hinweis in der Versionszeile (versionLine, Home und Model)
- [ ] docs/COMMANDS.md neben JAIRA_NO_UPDATE_CHECK, core/release/notes/339DNB.md

## Progress
- **2026-10-22 07:58 · BeMuCa** — Plan vom Treiber geschrieben (pre-process wurde von next_lane uebersprungen, in-progress meldete 'missing: plan'). Warum der Schalter in core/gitref und nicht nur an den drei Flush-Aufrufen: der Kontext nennt nur internal/cli/refs.go:88/:248 und internal/tui/refs.go:78, aber es gibt mehr Sendewege - internal/cli/web.go:109 (webPush), Syncer.Pull/Release (core/refsync/refsync.go:1395/:1478 schreiben den Ref direkt), der Snapshot (core/snapshot/snapshot.go:185 PushBranch, :252/:275 Delete) und rebaseEpic (core/refsync/epic.go:113). Alle laufen durch genau zwei git-push-Aufrufe: core/gitref/gitref.go PushBranch und pushRefspec. Flush prueft zusaetzlich frueh, damit die Outbox wie offline erhalten bleibt und kein Fetch/keine Warnung pro Eintrag entsteht.
