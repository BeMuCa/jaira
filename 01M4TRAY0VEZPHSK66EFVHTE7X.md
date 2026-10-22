---
id: 01M4TRAY0VEZPHSK66EFVHTE7X
title: "Ein selbst gebauter jaira zeigt, aus welchem Commit er stammt"
status: backlog
ready: true
creator: BeMuCa
goal: "jaira --version und der Projektkopf zeigen bei einem lokalen Build den Commit, so wie ein Release-Build seinen Tag zeigt."
context: |-
  Release-Builds zeigen den Tag - das bleibt: .goreleaser.yaml:22 setzt -X main.version={{.Version}}.
  Ein lokaler Build (go build -o ~/.local/bin/jaira ./cmd/jaira) zeigt nur 'dev' (cmd/jaira/main.go:11). Am 2026-10-13: 'jaira --version' gab 'jaira version dev'.
  'go version -m ~/.local/bin/jaira' zeigt keine vcs.revision-Zeile, nur 'mod (devel)'. Go hat den Commit also nicht eingestempelt. Warum (Build aus einem Worktree? -buildvcs?), ist nicht geprueft.
  Folge: Berk testet das laufende Binary und kann nicht sehen, welcher Stand laeuft. Am 2026-10-13 half nur die Dateizeit (Binary 14:05, D0SAHM-Commit 13:51, erster E70WYN-Commit 14:16).
  Verwandt: T8A8KM (Pille unter der Version im Projektkopf).
definition-of-done: "Ein lokaler Build zeigt bei 'jaira --version' und im Projektkopf 'dev' plus Kurz-Commit, und ob der Baum Aenderungen hatte; ein Release-Build zeigt unveraendert seinen Tag; ein Test deckt beide Faelle"
tags:
  - release
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-13T13:48:21Z
updated-at: 2026-10-22T07:52:35Z
assignee: BeMuCa
updated-by: BeMuCa
---

# Ein selbst gebauter jaira zeigt, aus welchem Commit er stammt

## Definition of Done

- [ ] Ein lokaler Build zeigt bei 'jaira --version' und im Projektkopf 'dev' plus Kurz-Commit, und ob der Baum Aenderungen hatte; ein Release-Build zeigt unveraendert seinen Tag; ein Test deckt beide Faelle
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

