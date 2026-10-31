---
id: 01M67RWX6BPVT0QKN5ZAX1M1AF
title: "Neue Tickets haben die Option planning standardmaessig angehakt, damit jedes Ticket durch den Planer geht"
status: todo
ready: true
creator: berkc
assignee: berkc
goal: "Ein neu angelegtes Ticket hat unter Options 'planning' angehakt, solange die Person keine eigene Default-Board-Datei hat, die etwas anderes sagt; die pre-process-Lane (Opus-Planer) laeuft dadurch fuer jedes Ticket."
context: |-
  Was falsch ist: Neue Tickets haben alle Options unangehakt (core/lane/defaultboard.go:185 ResolveOptions tickt nur, was die Default-Board-Datei nennt; ohne Datei nichts). Die pre-process-Lane (requires-option: planning) wird deshalb fast immer uebersprungen.
  Ausloeser: Experiment 3CBCN6 auf NDE933 (31.10.2026). Nur-Haiku vergass, Hilfezeile und Kommentar mitzuziehen (Urteil 8/7/7/6); Opus-Plan + Haiku-Bau lag vorne (9/9/8/9). Der Maintainer: 'should be on per default der planning tag. damit jeder durch den opus planer geht.'
  Bekannt: Optionen kommen aus requires-option der installierten Lanes (core/lane/lane.go:233 Options). Die Default-Board-Datei pro Person bleibt der Weg, es abzuschalten.
  Haengt zusammen mit CXH98W: dass pre-process wirklich mit Opus laeuft, erzwingt erst CXH98W (Modell pro Lane fuer Worker).
definition-of-done: "jaira create ohne Default-Board-Datei legt ein Ticket an, dessen Options 'planning' angehakt hat; Test dafuer"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-31T01:24:00Z
updated-at: 2026-10-31T02:31:33Z
updated-by: berkc
claimed-by: berk-ca-3129208
claimed-at: 2026-10-31T01:27:35Z
---

# Neue Tickets haben die Option planning standardmaessig angehakt, damit jedes Ticket durch den Planer geht

## Definition of Done

- [ ] jaira create ohne Default-Board-Datei legt ein Ticket an, dessen Options 'planning' angehakt hat; Test dafuer
- [ ] Eine Default-Board-Datei, die options ohne planning nennt, laesst planning unangehakt; Test dafuer
- [ ] Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

