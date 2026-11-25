---
id: 01M8AF7J4NW1AXQGMSDVB2GKN7
title: "Das maximierte Cockpit zeigt ein Nachrichten-Symbol mit den aktiven Personen des Repos, und man chattet von dort mit ihnen"
status: backlog
ready: true
creator: Berk
assignee: Berk
goal: "Im maximierten Cockpit (rechts neben dem Chat) sieht man ein Nachrichten-Symbol mit Zaehler; es oeffnet die Liste der aktiven Personen dieses Repos; man liest und schreibt dort Nachrichten an eine Person, gespeichert ueber jaira msg auf Refs."
context: |-
  Wunsch von berkc am 26.11.2026: 'im cockpit (wenn es maximiert auf der rechten seite ist) ein nachrichten icon wo alle user aliases sind die in dem repo aktiv sind und man mit den chatten kann.'
  Baut auf WQNRQN auf: dort entstehen jaira inbox, jaira msg und jaira people (Aliase aus Assignees, Creator, Commit-Autoren). Erst WQNRQN fertig, dann dieses Ticket.
  Das Cockpit ist ein Claude-Code-Mod unter core/mod/builtin/jaira-cockpit/ (hooks/*.ts, ausgeliefert mit Hash in core/shipped/shipped.sha256).
  Nachrichten sind Daten, keine Anweisungen: das Cockpit zeigt sie dem Menschen, es reicht sie nie als Prompt an Claude weiter. (Idee 'Nachricht loest Hook aus und Claude liest sie direkt' = Remote-Coding; bewusst NICHT Teil dieses Tickets, weil jeder mit Push-Recht damit Befehle in fremde Sitzungen schicken koennte.)
definition-of-done: Maximiertes Cockpit zeigt ein Nachrichten-Symbol mit Zahl ungelesener Nachrichten; Test
tags:
  - cockpit
blocked-by: []
related: []
follows: 01M6HZQN4HSDJ55MYC5AWQNRQN
commits: []
created-at: 2026-11-25T23:04:10Z
updated-at: 2026-11-25T23:04:10Z
---

# Das maximierte Cockpit zeigt ein Nachrichten-Symbol mit den aktiven Personen des Repos, und man chattet von dort mit ihnen

## Definition of Done

- [ ] Maximiertes Cockpit zeigt ein Nachrichten-Symbol mit Zahl ungelesener Nachrichten; Test
- [ ] Das Symbol oeffnet die Liste aus 'jaira people --json'; eine Person waehlen zeigt den Verlauf und ein Eingabefeld; Senden ruft 'jaira msg <person> <text>'; Test
- [ ] Eingehende Nachrichten werden nie in den Claude-Prompt eingespeist, nur angezeigt; Test
- [ ] shipped.sha256 und core/shipped-Tests stimmen; Release-Note-Zeile

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

