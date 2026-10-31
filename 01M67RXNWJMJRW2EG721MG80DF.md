---
id: 01M67RXNWJMJRW2EG721MG80DF
title: "Das Cockpit zeigt nur die Tickets, die mir gehoeren"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Band und Zaehler des jaira-cockpit zeigen nur Tickets, deren assignee die Person ist, die die Session fuehrt (git user.name plus Aliase aus ~/.jaira/identity); fremde Tickets erscheinen nicht."
context: |-
  Was falsch ist: das Cockpit-Band (core/mod/builtin/jaira-cockpit/hooks/board.ts) listet alle wartenden Tickets des Boards, auch die von Kollegen.
  Ausloeser: Maintainer 31.10.2026: 'Erstens will ich nur die tickets von mir und nicht von anderen angezeigt kriegen.'
  Bekannt: Identitaet = git user.name plus Aliase in ~/.jaira/identity. Tickets ohne assignee: im Plan entscheiden (Vorschlag: zeigen, sie gehoeren noch niemandem).
definition-of-done: Ein wartendes Ticket mit fremdem assignee erscheint nicht im Band; Test in hooks
tags:
  - tui
blocked-by: []
related:
  - 01M5M7JTVR8G7DZ1MGHJQ63EEV
commits: []
created-at: 2026-10-31T01:24:25Z
updated-at: 2026-10-31T01:24:25Z
---

# Das Cockpit zeigt nur die Tickets, die mir gehoeren

## Definition of Done

- [ ] Ein wartendes Ticket mit fremdem assignee erscheint nicht im Band; Test in hooks
- [ ] Ein Ticket mit eigenem assignee oder einem Alias erscheint; Test
- [ ] Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

