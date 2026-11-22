---
id: 01M7SHP0X7Q27BWFQMDS1J8XZ7
title: Das Cockpit zeigt einen Milestone mit Fortschritt und seinen Tickets nach Stand sortiert
status: in-progress
ready: true
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
updated-at: 2026-11-22T22:47:05Z
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
- [ ] edge cases: no milestones, a milestone with no tickets, milestone ls failing (older jaira) shows a could-not-read line, a remembered milestone that is gone falls back to the first, a ticket jaira show cannot find stays in the list with lane '?'
- [ ] seen by a person: the Milestones tab in a real Claude Code terminal - selector paged by click and by arrow keys + Enter, rows readable, head row still on one line at 80 columns

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [ ] overview.ts: read jaira milestone ls --json into Overview.milestones (null when the read failed), keep each open ticket's lane+title from jaira list in Overview.tickets, and a pure milestoneView() that groups and sorts one milestone's rows and counts x of n
- [ ] register.tsx: run milestone ls beside lanes/list/next at open; atoms overviewMilestone (chosen id, remembered) and overviewFiled (show results for handles not in jaira list); on open and on every page run jaira show --json for the chosen milestone's absent handles
- [ ] overview-pane.tsx: fourth tab Milestones; selector [ < ] name [ > ] as buttons (click, or arrows to focus + Enter), wraps around; 'x of n done' line; rows ID | lane | title in three groups; none/unread states
- [ ] types/index.d.ts and texts.ts: new state keys and words
- [ ] tests/board.test.ts: fix the tab-label test, add milestone tests (grouping, logbook count, paging wrap, remembered after reopen, ls failure)
- [ ] release note core/release/notes/1J8XZ7.md
- [ ] gates: tsc/validate of the mod if cheap, claude plugin test once; gofmt/vet/build untouched (no Go change)

## Progress
- **2026-11-22 22:47 · BeMuCa** — pre-process, decisions settled (no person needed, disagree at acceptance): (1) data source deviates from the context: 'jaira list --milestone X' on 0.6.1 leaves out logbook tickets (release-0-6-0 has 3 tickets, list --milestone returns 0), so the tab takes the handles from 'jaira milestone ls --json' (milestones[].tickets), lane+title from the 'jaira list --json' the pane already reads, and runs 'jaira show <h> --json' only for the chosen milestone's handles list does not return - show finds logbook tickets (filed_away set) and tickets on other boards. (2) Arrow keys: the cockpit has no key handler; a Pane's arrows move the focus ring, so the selector is two Buttons [ < ] and [ > ] around the name - click, or arrows onto a button and Enter. Paging wraps around. (3) Words in English like the rest of the cockpit: 'x of n done'. (4) Groups: done = terminal lane or filed away (rejected ones too, lane shown as 'rejected', else 'logbook'); in work = agentic lane or a person's lane; everything else (ideas, backlog, todo, blocked, unknown) = not started - blocked sits there as the Queue tab already puts it under 'Put back'. Within a group board lane order. (5) Twins (two milestones of one name): selector shows name-ID; the remembered choice is the milestone ID.
