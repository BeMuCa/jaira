---
id: 01M59MQ9PF25WTEXVEGH28PD7N
title: "Landet ein Ticket in einer Lane, die einen Menschen braucht, wird der Mensch aktiv benachrichtigt"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Fragt ein Agent etwas (HITL) oder wartet Arbeit auf Abnahme (Human Review), erfaehrt der Mensch es sofort, ohne aufs Board zu schauen."
context: |-
  Heute: 6TCTFV lag im HITL mit einer Frage, berkc hat es nicht gemerkt. Der Agent schreibt die Frage nur aufs Ticket (Karte zeigt '▲ asks').
  Das Cockpit-Band (core/mod/builtin/jaira-cockpit/hooks/register.tsx) zeigt wartende Tickets ueber dem Prompt. Erst seit 0.5.4 (NC323D) zaehlt es auch Lanes mit requires-question (HITL) mit - vorher nur Human Review. Nur wer das Cockpit installiert hat, sieht es.
  Was fehlt: ein aktiver Hinweis im Moment des Moves (z.B. jaira move gibt eine deutliche Zeile aus, die der Agent dem Menschen weitergeben MUSS; Desktop-Benachrichtigung / Terminal-Bell im TUI; Cockpit-Band blinkt).
  Wunsch von berkc: jede Lane, die einen Menschen braucht, traegt die Regel 'wird die Lane benutzt und nicht uebersprungen, benachrichtige den Menschen' - als Teil der ausgelieferten Lanes human.md und signoff.md und der Vorlage fuer Uebergabe-Lanes.
definition-of-done: "jaira move in eine Lane mit requires-question oder requires-human-exit gibt eine Zeile aus, die sagt, dass der Mensch jetzt gefragt ist, und was er tun soll; Test in internal/cli"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-19T08:33:49Z
updated-at: 2026-10-19T08:33:49Z
---

# Landet ein Ticket in einer Lane, die einen Menschen braucht, wird der Mensch aktiv benachrichtigt

## Definition of Done

- [ ] jaira move in eine Lane mit requires-question oder requires-human-exit gibt eine Zeile aus, die sagt, dass der Mensch jetzt gefragt ist, und was er tun soll; Test in internal/cli
- [ ] Die ausgelieferten Lanes human und signoff sagen im Prompt/Beschreibung, dass der Agent den Menschen in seiner Sitzung benachrichtigt; jaira show --for-lane enthaelt es
- [ ] Das TUI zeigt ein neu wartendes Ticket deutlich an (Meldung oder Bell), einmal pro Ticket; Test
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

