---
id: 01M5MNX0YKP4QAACC80KPBDMS4
title: Milestone-Dateien fuehren keine archivierten Tickets mehr als offen
status: backlog
ready: false
creator: BeMuCa
goal: jaira milestone ls/show und die Kartenkante auf dem Board zaehlen ein archiviertes oder abgelegtes Ticket nicht mehr als offenen Teil des Milestones.
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkt 11): Milestone-Dateien im Haupt-Checkout listen archivierte Tickets weiter (P6Y1BG, XY4HDS, 9JGWGR) - der Milestone wirkt unfertig, und die Board-Dateien des Haupt-Checkouts hinken den Worktrees hinterher.
  Verwandt: GWK31K (archiviertes Ticket verschwindet sofort von allen Boards, signoff), N71NVG (fertiger Milestone ins Archiv), KWQCYZ (veraltete Milestone-Datei ueberschreibt beim add). Dieses Ticket ist der Rest: die Zaehlung und Anzeige, wenn die ID in der Datei bleibt.
definition-of-done: "jaira milestone ls zaehlt archivierte und abgelegte Tickets nicht als offen und zeigt sie in show als 'archiviert'/'abgelegt'; ein Test baut den Fall"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:26:04Z
updated-at: 2026-10-26T20:22:44Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: pbdms4-driver
claimed-at: 2026-10-26T20:22:44Z
---

# Milestone-Dateien fuehren keine archivierten Tickets mehr als offen

## Definition of Done

- [ ] jaira milestone ls zaehlt archivierte und abgelegte Tickets nicht als offen und zeigt sie in show als 'archiviert'/'abgelegt'; ein Test baut den Fall
- [ ] Die Kartenkante und der M-Picker ignorieren sie
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

