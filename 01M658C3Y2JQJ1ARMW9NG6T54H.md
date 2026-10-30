---
id: 01M658C3Y2JQJ1ARMW9NG6T54H
title: Agents per ticket steht per Default auf worker statt lanes
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Wer nichts eingestellt hat, bekommt fuer ein ganzes Ticket jaira-worker (ein Agent durch alle Lanes, nur fresh-context-Lanes neu) statt jaira-dispatcher."
context: |-
  Wunsch von berkc am 30.10.2026: 'Agents per ticket ... das soll default auf worker sein'.
  Heute: core/settings/settings.go:278 DispatchMode gibt 'lanes' zurueck, ausser settings.json sagt ausdruecklich 'worker'. Hilfetext internal/cli/roles.go:40 nennt 'lanes' den Default; Settings-Zeile internal/tui/settings.go:56.
  berkc selbst hat schon dispatch: worker in ~/.jaira/settings.json - betrifft also neue Nutzer.
  Wer 'lanes' ausdruecklich gesetzt hat, behaelt es.
definition-of-done: "DispatchMode liefert worker, wenn dispatch leer oder unbekannt ist; 'lanes' nur ausdruecklich; Test"
tags:
  - cli
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-30T01:56:44Z
updated-at: 2026-10-30T01:56:44Z
---

# Agents per ticket steht per Default auf worker statt lanes

## Definition of Done

- [ ] DispatchMode liefert worker, wenn dispatch leer oder unbekannt ist; 'lanes' nur ausdruecklich; Test
- [ ] Hilfetexte (roles, Settings) nennen worker als Default
- [ ] Release-Note-Zeile

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

