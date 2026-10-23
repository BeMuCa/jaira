---
id: 01M5MNX64Z9DZE1KTFZ9SZFTQD
title: "Ergebnisse fuer eine Person liegen auf einem stabilen Pfad, nie unter /tmp"
status: backlog
ready: false
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
updated-at: 2026-10-23T15:26:09Z
---

# Ergebnisse fuer eine Person liegen auf einem stabilen Pfad, nie unter /tmp

## Definition of Done

- [ ] jaira-worker, jaira-dispatcher, jaira-role-lane und jaira-role-acceptance sagen: ein Ergebnis fuer eine Person wird unter einem stabilen Pfad abgelegt (Board-Setting 'deliverables' oder, wenn leer, <repo>/.jaira/deliverables/<ticket>/), nie unter /tmp, und der Pfad steht als Notiz am Ticket
- [ ] Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

