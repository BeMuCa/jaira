---
id: 01M6NHNN9AF2W8R6MDQK8Z3VRF
title: "jaira <id> startet ein installiertes Plugin"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Befehl, den jaira nicht kennt, startet das gleichnamige installierte Plugin, mit derselben Umgebung auf Linux, macOS und Windows."
context: |-
  Heute: 'jaira jira ...' endet mit 'unknown command'. Es gibt keinen Weg, jaira zu erweitern, ohne den Kern zu aendern.
  Gewollt (Standard in R6X6FH, docs/PLUGINS.md): unbekannter Unterbefehl <id> -> ~/.jaira/plugins/<id>/jaira-<id>(.exe) starten, Argumente, stdin/stdout/stderr und Exit-Code durchreichen.
  Umgebung: JAIRA_ROOT, JAIRA_BOARD, JAIRA_VERSION.
  Eingebaute Befehle gehen immer vor; ein Plugin kann keinen Kernbefehl ueberdecken.
  Kein PATH-Suchen: nur was 'jaira plugin install' abgelegt hat, laeuft (Whitelist).
definition-of-done: "Unbekannter Unterbefehl startet ~/.jaira/plugins/<id>/jaira-<id> mit Argumenten, Umgebung und Exit-Code; Test"
tags:
  - cli
blocked-by:
  - R6X6FH
related: []
commits: []
created-at: 2026-11-05T09:47:04Z
updated-at: 2026-11-05T21:41:05Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-75959
claimed-at: 2026-11-05T21:39:09Z
---

# jaira <id> startet ein installiertes Plugin

## Definition of Done

- [ ] Unbekannter Unterbefehl startet ~/.jaira/plugins/<id>/jaira-<id> mit Argumenten, Umgebung und Exit-Code; Test
- [ ] Ein eingebauter Befehl gleichen Namens gewinnt; Test
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md
- [ ] Events: create, move, claim und set aus der CLI und Moves aus der TUI rufen jedes abonnierte Plugin als 'jaira-<id> event <name>' auf (stdin = show --json + actor + from_lane, 5 s Timeout, Fehler aendert nichts, JAIRA_IN_EVENT=1 unterdrueckt Events); Test

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-05 21:39 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/47752.sock
