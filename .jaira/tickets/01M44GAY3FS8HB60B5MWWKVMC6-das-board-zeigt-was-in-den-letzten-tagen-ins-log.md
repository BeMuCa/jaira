---
id: 01M44GAY3FS8HB60B5MWWKVMC6
title: "Das Board zeigt, was in den letzten Tagen ins Logbook ging, und wie viele Tage stellt man in den Einstellungen ein"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Im TUI sieht man ohne Befehl, was zuletzt vom Board ging: die Logbook-Tickets der letzten N Tage stehen in der Done-Lane, sichtbar anders als noch nicht abgelegte; N steht in den Einstellungen und gilt auch fuer 'jaira logbook'."
context: |-
  Was fehlt:
  - Wer im Board arbeitet, sieht nicht, was zuletzt ins Logbook ging. Ein abgelegtes Ticket verschwindet einfach (Absicht: abgelegt = nicht mehr auf der Tafel).
  - Nachschauen geht nur in der CLI: 'jaira logbook' (seit MJ3PVV, 0.3.2: letzte vier Wochen, --since Nw/Nd, --since 0 = alles).
  - MJ3PVV hat das TUI ausdruecklich ausgenommen. Alex hatte aber genau das gemeint (05.10.2026): im TUI anzeigen und in den Settings einstellen, wie viele Tage.
  - Folgt auf MJ3PVV (01M43QT4X8K72505R1HDMJ3PVV). --follows liess sich nicht setzen: 'jaira create --follows' findet ein Ticket im Logbook nicht ('ticket: not found').

  Entschieden (Alex, 05.10.2026):
  - Ort: in der Done-Lane (Terminal-Lane), unter den noch nicht abgelegten Tickets, mit anderem Rahmen bzw. anderer Farbe, damit man 'fertig, liegt noch auf dem Board' von 'schon im Logbook' unterscheidet. Zuerst war eine eigene Spalte rechts vorgeschlagen; Alex' Notiz dazu: vielleicht in done, aber mit anderem Rahmen in der Farbe.
  - Zeitraum: Einstellung 'logbook-days' in ~/.jaira/settings.json, Standard 28 (wie 'jaira logbook' ohne Argument), im TUI auf dem Einstellungsbildschirm (S) aenderbar. Derselbe Wert wird Standard fuer 'jaira logbook'; --since ueberschreibt ihn weiter.

  Bekannt:
  - Das Ablagedatum steht im Ordnernamen <initialen>-<yyyymmdd>; ticket.LogbookFolderDay (core/ticket/store.go) liest es schon fuer CLI und Launcher-Grafik.
  - Die Launcher-Grafik (internal/tui/home.go, logbookDays = 7) ist etwas anderes und bleibt.
  - Logbook-Karten sind nur zum Lesen: Enter oeffnet das Ticket, kein Verschieben. Zurueckholen bleibt 'jaira restore'.
definition-of-done: "In der Done-Lane stehen unter den noch nicht abgelegten Tickets die Logbook-Tickets der letzten N Tage, neueste zuerst, sichtbar anders gerahmt bzw. gefaerbt"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-04T22:25:15Z
updated-at: 2026-10-04T22:34:45Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-734
claimed-at: 2026-10-04T22:25:33Z
outcome-what: "Terminal-Lane zeigt unter ihren Tickets die Logbook-Tickets der letzten logbook-days Tage (neueste zuerst, tuerkis, Marker '⎙ filed <Tag>', nur lesen); logbook-days in settings.json (Standard 28, 0 = aus), auf S editierbar; jaira logbook ohne --since nimmt denselben Zeitraum. Neu: settings.LogbookWindow, ticket.LoggedSince/LogbookWindowStart, internal/tui/logbook.go."
outcome-why: "Ein abgelegtes Ticket verschwand spurlos vom Board; was zuletzt fertig wurde, sah man nur in der CLI. Alex wollte es im TUI, in der Done-Lane, mit einstellbarer Tageszahl."
outcome-resolves: "Alle sechs DoD-Punkte; Tests: TestLogbookWindow, TestLoggedSinceReadsTheWindowNewestFirst, TestLogbookListingFollowsLogbookDays, TestLogbookCardsSitBelowTheTerminalLane, TestLogbookDaysZeroHidesTheCards, TestALogbookCardIsReadOnly, TestSettingsEditsLogbookDays."
---

# Das Board zeigt, was in den letzten Tagen ins Logbook ging, und wie viele Tage stellt man in den Einstellungen ein

## Definition of Done

- [x] In der Done-Lane stehen unter den noch nicht abgelegten Tickets die Logbook-Tickets der letzten N Tage, neueste zuerst, sichtbar anders gerahmt bzw. gefaerbt
  proof: internal/tui/model.go rebuild (filed cards appended to terminal column); internal/tui/logbook.go logbookShade/logbookFlag; TestLogbookCardsSitBelowTheTerminalLane
- [x] Eine Logbook-Karte laesst sich mit Enter oeffnen und nicht verschieben
  proof: internal/tui/logbook.go refuseLogged, model.go openDetail/openMove/archiveSelected/modeDetail guard; TestALogbookCardIsReadOnly
- [x] logbook-days in ~/.jaira/settings.json setzt N (Standard 28); 0 blendet die Logbook-Karten aus
  proof: core/settings/settings.go LogbookDays/LogbookWindow; TestLogbookWindow, TestLogbookDaysZeroHidesTheCards
- [x] Der Einstellungsbildschirm S zeigt und aendert logbook-days
  proof: internal/tui/settings.go 'Logbook on the board' row, editKey; TestSettingsEditsLogbookDays
- [x] jaira logbook ohne --since nimmt logbook-days als Zeitraum; --since gewinnt
  proof: internal/cli/logbook.go RunE (settings.LogbookWindow unless --since changed); TestLogbookListingFollowsLogbookDays
- [x] core/release/NOTES.md, README und docs/COMMANDS.md beschreiben die Anzeige und die Einstellung
  proof: core/release/NOTES.md Unreleased; README.md Keys + Reviewing finished work; docs/COMMANDS.md jaira logbook row + Keys

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] settings: logbook-days (Zeiger, Standard 28, 0 = aus) in core/settings
- [x] core/ticket: Fenster-Start und LoggedSince (Logbook-Tickets der letzten N Tage, neueste zuerst)
- [x] CLI: jaira logbook nimmt logbook-days, wenn --since fehlt
- [x] TUI: Logbook-Karten unten in der Terminal-Lane, eigene Farbe, nur lesen
- [x] TUI: Einstellungsbildschirm S zeigt und aendert logbook-days
- [x] Tests, NOTES.md, README, docs/COMMANDS.md

## Progress
- **2026-10-04 22:34 · Alexander Sacharov** — In-progress: Logbook-Karten haengen in rebuild() hinter den Tickets der Terminal-Lane (column.filed zaehlt sie), erkannt per Zeiger (m.loggedDay), nicht per ID - ein Ticket kann nach einem Merge zugleich auf dem Board und im Logbook liegen, die Board-Kopie muss bearbeitbar bleiben. openDetail nimmt die Logbook-Karte direkt (store.Load wuerde die Board-Kopie liefern). Ordner ohne Datum zeigt das Board nicht (die CLI-Liste schon). logbook-days ist *int, weil 0 = aus von 'nicht gesetzt' (=28) unterschieden werden muss; logbook-days 0 heisst fuer 'jaira logbook' ohne --since: alles, wie --since 0. Follow-up (n) auf einer Logbook-Karte bleibt erlaubt, er schreibt ein neues Ticket. Shade-Test: '48;5;23' ist Praefix von 234 (Lane-Grau) - deshalb mit 'm' am Ende pruefen.
