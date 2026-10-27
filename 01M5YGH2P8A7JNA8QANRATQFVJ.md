---
id: 01M5YGH2P8A7JNA8QANRATQFVJ
title: "Ein Subagent heisst nach seinem Ticket, nicht 'general purpose'"
status: backlog
ready: false
creator: BeMuCa
goal: "In der Agentenliste von Claude Code und im Cockpit steht bei jedem Worker-Subagenten das Ticket (ID + Lane), an dem er arbeitet - nicht der generische Typ."
context: |-
  Berk, 2026-10-27: 'When the teamlead starts subagents (worker agents) for tickets, they are named general purpose most of the times. Can we change them to hold the ticket name that they are working on?'
  Heute: das Agent-Tool kennt nur subagent_type (z.B. general-purpose), description (Kurzlabel) und prompt. Claude Code zeigt den Typ; das Cockpit erkennt Worker nur, wenn description oder name die Ticket-ID traegt (core/mod/builtin/jaira-cockpit/hooks/board.ts, 'said.includes(id)'). Per spawn.sh gestartete Herdr-Tabs heissen schon '<id>/<lane>'.
  Zwei Wege, Brainstorm entscheidet: (a) die Rollen (jaira-teamlead, jaira-dispatcher, jaira-worker) schreiben vor, dass jede Agent-description mit '<ID> <lane>' beginnt; (b) eigene Agententypen unter .claude/agents/ (z.B. jaira-worker, jaira-lane), die als Typname erscheinen und die Rolle gleich mitbringen.
definition-of-done: Brainstorm hat (a) oder (b) oder beides entschieden und begruendet
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T11:04:34Z
updated-at: 2026-10-27T11:04:34Z
---

# Ein Subagent heisst nach seinem Ticket, nicht 'general purpose'

## Definition of Done

- [ ] Brainstorm hat (a) oder (b) oder beides entschieden und begruendet
- [ ] Folge-Tickets fuer die Umsetzung sind angelegt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

