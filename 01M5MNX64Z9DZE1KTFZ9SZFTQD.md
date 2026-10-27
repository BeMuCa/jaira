---
id: 01M5MNX64Z9DZE1KTFZ9SZFTQD
title: "Ergebnisse fuer eine Person liegen auf einem stabilen Pfad, nie unter /tmp"
status: in-progress
ready: true
creator: BeMuCa
goal: "Was ein Worker fuer einen Menschen erzeugt (Vergleichsseite, Bericht, Screenshot), ueberlebt einen Crash: die Rollen sagen, wohin es gehoert, und der Pfad steht am Ticket."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkt 16): die Vergleichsseite eines Workers lag im Scratchpad unter /tmp und war nach dem WSL-Crash weg. Dasselbe passierte dieser Session mit Backup-Bundle und SHA-Map.
  Das Scratchpad ist fuer Zwischenergebnisse richtig; ein Ergebnis fuer eine Person ist keins. Die Rollen (jaira-worker, jaira-role-lane, jaira-role-acceptance) nennen heute keinen Ort dafuer.
definition-of-done: "jaira-worker, jaira-dispatcher, jaira-role-lane und jaira-role-acceptance sagen: ein Ergebnis fuer eine Person wird unter einem stabilen Pfad abgelegt (Board-Setting 'deliverables' oder, wenn leer, <repo>/.jaira/deliverables/<ticket>/), nie unter /tmp, und der Pfad steht als Notiz am Ticket"
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:26:09Z
updated-at: 2026-10-27T00:05:00Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: szftqd-driver
claimed-at: 2026-10-27T00:01:42Z
---

# Ergebnisse fuer eine Person liegen auf einem stabilen Pfad, nie unter /tmp

## Definition of Done

- [x] jaira-worker, jaira-dispatcher, jaira-role-lane und jaira-role-acceptance sagen: ein Ergebnis fuer eine Person wird unter einem stabilen Pfad abgelegt (Board-Setting 'deliverables' oder, wenn leer, <repo>/.jaira/deliverables/<ticket>/), nie unter /tmp, und der Pfad steht als Notiz am Ticket
  proof: core/role/builtin/jaira-worker/SKILL.md:331, jaira-dispatcher/SKILL.md:474, jaira-role-lane/SKILL.md:342, jaira-role-acceptance/SKILL.md:52
- [ ] Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-27 00:03 · BeMuCa** — Driver: brainstorm und pre-process uebersprungen (Optionen nicht gesetzt, next_lane lt. lane.OnPath). in-progress an Subagent (tier cheap) gegeben; Vorgabe: kein neues deliverables-Setting, nur Fallback .jaira/deliverables/<handle>/, ausser es ist trivial.
- **2026-10-27 00:03 · BeMuCa** — (a) No 'deliverables' setting exists: grep -rni deliverables over *.go finds nothing; settings are the machine JSON in core/settings/settings.go plus single files per board. The skills say 'the board setting deliverables if the board has one' but nothing reads it, so the fallback .jaira/deliverables/<handle>/ is what actually applies. Not added (not text-only). (b) .jaira/deliverables/ is safe: .jaira subdirs are read by name only (boards/ in core/lane/boards.go:60, tickets/logbook in core/ticket/store.go:519); board discovery core/project/discover.go:45 checks only .jaira/tickets. No scan of .jaira itself. (c) 'jaira init' writes /.jaira/ to .gitignore (core/board/gitignore.go:10,44-56); 'jaira share' removes it so deliverables are then untracked-but-visible and get committed if someone git adds them; only /.jaira/lanes/ stays ignored (gitignore.go:19). Large files there would reach the repo.
