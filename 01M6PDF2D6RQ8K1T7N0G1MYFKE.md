---
id: 01M6PDF2D6RQ8K1T7N0G1MYFKE
title: "Ein Milestone erreicht die Kollegen so schnell wie ein Ticket, ueber einen Ref"
status: backlog
ready: false
creator: BeMuCa
goal: "Ein angelegter oder geaenderter Milestone liegt auf einem Ref wie Tickets und Epics und erscheint bei Kollegen mit dem naechsten Fetch, ohne Commit und Merge."
context: |-
  Was schief ist: Milestones leben nur als Dateien unter .jaira/milestones/. Sie erreichen andere erst ueber Commit, Push, Merge und Pull.
  Im jaira-Repo selbst ist auf master keine einzige Milestone-Datei getrackt (git ls-files .jaira/milestones: 0, Stand 2026-11-05). Die Roadmap eines Kollegen zeigt keinen unserer Milestones.
  Tickets liegen auf refs/jaira/tickets/*, Epics auf refs/jaira/epics (grep in core/), Milestones auf keinem Ref.
  Frage von Berk (2026-11-05): ab wann sieht ein anderer meinen Milestone auf seiner Roadmap? Antwort heute: erst nach Merge.
  Vorbild: die Epics - Beschreibung plus Ticketliste auf einem eigenen Ref (jaira epic --help).
  Achtung: update never loses data - vorhandene Dateien werden uebernommen, nicht verworfen.
definition-of-done: "Milestones liegen auf einem Ref; create, add, rm, rename, set schreiben ihn"
tags:
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-11-05T17:52:49Z
updated-at: 2026-11-05T17:52:49Z
---

# Ein Milestone erreicht die Kollegen so schnell wie ein Ticket, ueber einen Ref

## Definition of Done

- [ ] Milestones liegen auf einem Ref; create, add, rm, rename, set schreiben ihn
- [ ] Ein zweiter Clone sieht einen neuen Milestone nach jaira fetch, ohne Commit
- [ ] Vorhandene Milestone-Dateien werden beim ersten Lauf auf den Ref uebernommen; ein Test belegt es
- [ ] Release-Notiz in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

