---
id: 01M7SHP0X7Q27BWFQMDS1J8XZ7
title: Das Cockpit zeigt einen Milestone mit Fortschritt und seinen Tickets nach Stand sortiert
status: backlog
ready: false
creator: BeMuCa
goal: "Das Overview-Pane des Cockpits hat einen Tab Milestones: oben ein Waehler '< milestone-name >' zum Durchblaettern, darunter 'x von n fertig' und die Tickets des Milestones mit Lane, in Arbeit ganz oben, nicht begonnene in der Mitte, fertige ganz unten."
context: |-
  Wunsch von Berk (2026-11-19): im Cockpit sieht man nicht, wie weit ein Milestone ist. Dafuer muss man heute ins Board (M) oder jaira list --milestone <name> tippen.
  Heute hat das Overview-Pane (core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx) die Tabs Queue, Log und Workers. Es liest nur das Board ueber jaira lanes, jaira list und jaira next --all (overview.ts Kopf). Milestones kommen dort nicht vor.
  Datenquelle: jaira milestone ls fuer die Namen, jaira list --milestone <name> --json fuer die Tickets mit Lane. Kein neuer Lesepfad noetig.
  Sortierung: in Arbeit = jede agentische Lane plus human, review, signoff; nicht begonnen = ideas, backlog, todo; fertig = done und die im Logbook. Innerhalb einer Gruppe nach Lane-Reihenfolge.
  'x von n' zaehlt done plus Logbook gegen alle Tickets des Milestones.
  Der Waehler merkt sich den zuletzt gewaehlten Milestone, wie das Pane heute Tab und Faltung behaelt (register.tsx).
  Zusammenhang: 35Y6N5 (dasselbe in jaira web), XGBVRN (Ersteller im Picker), 1MYFKE (Milestones auf einem Ref - bis dahin zeigt der Tab nur die lokalen Milestone-Dateien).
  Nicht gewollt: Milestones im Cockpit anlegen oder aendern. Nur ansehen.
definition-of-done: "Das Overview-Pane hat einen Tab Milestones mit Waehler '< name >', der mit Pfeiltasten oder Klick durch alle Milestones blaettert"
tags:
  - tui
blocked-by: []
related:
  - 01M6H733PCXM02SCZJ6T35Y6N5
  - 01M6PDF0D9HS2SE2SCPQXGBVRN
  - 01M6PDF2D6RQ8K1T7N0G1MYFKE
commits: []
created-at: 2026-11-19T09:19:56Z
updated-at: 2026-11-22T22:43:54Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-385148
claimed-at: 2026-11-22T22:43:54Z
---

# Das Cockpit zeigt einen Milestone mit Fortschritt und seinen Tickets nach Stand sortiert

## Definition of Done

- [ ] Das Overview-Pane hat einen Tab Milestones mit Waehler '< name >', der mit Pfeiltasten oder Klick durch alle Milestones blaettert
- [ ] Unter dem Waehler steht 'x von n fertig' und die Ticketliste mit ID, Titel und Lane: in Arbeit oben, nicht begonnen in der Mitte, fertig unten
- [ ] Der zuletzt gewaehlte Milestone bleibt nach Schliessen und Oeffnen des Panes gewaehlt
- [ ] Release-Notiz in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

