---
id: 01M44GAY3FS8HB60B5MWWKVMC6
title: "Das Board zeigt, was in den letzten Tagen ins Logbook ging, und wie viele Tage stellt man in den Einstellungen ein"
status: testing
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
updated-at: 2026-10-04T22:38:05Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-734
claimed-at: 2026-10-04T22:25:33Z
outcome-what: "Lese-Schutz der Logbook-Karten aus der Tastenliste in modeDetail in die Aufgerufenen verlegt (startEdit, openInEditor, X); Test deckt E ab."
outcome-why: "Critique: eine zweite Tastenliste wird beim naechsten schreibenden Key vergessen; openMove/archiveSelected schuetzen sich schon selbst."
outcome-resolves: "Critique-Befund aus Durchgang 1; DoD unveraendert erfuellt."
review-summary: "none"
review-gaps: "Model.logbookDays gestrichen (nur in loadLogbook gelesen, und gleichnamig mit der Konstante logbookDays der Launcher-Grafik in home.go) - loadLogbook liest settings selbst; nil-Pruefung in isLogged gestrichen (Map-Lookup mit nil ist false). Stehen gelassen: LoggedPerDay (store.go) neben LoggedSince - zaehlt nur Dateinamen ohne Tickets zu lesen und schneidet anders (Tage bis heute), Zusammenlegen wuerde die Launcher-Grafik teurer machen; settingsActionLogbookDays als Zeilenmarker; vorhandenes gateEnv() pro Karte in renderCard ist alt."
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
- **2026-10-04 22:35 · Alexander Sacharov** — critique (1. Durchgang): ein Befund - der Lese-Schutz im Detail haengt an einer eigenen Tastenliste vor dem switch in model.go; er gehoert in die Aufgerufenen (startEdit in edit.go, openInEditor in external.go, Fall X), so wie openMove/archiveSelected es schon tun. Bewusst stehen gelassen: settingsActionLogbookDays als Marker der Zeile (nie zurueckgegeben) - passt zur Eintragstabelle; logbook-days 0 = CLI listet alles (wie --since 0, dokumentiert); Ordner ohne Datum nicht auf dem Board (Notiz in-progress).
- **2026-10-04 22:36 · Alexander Sacharov** — in-progress (Runde 2): Critique-Befund umgesetzt - refuseLogged sitzt jetzt in startEdit (edit.go), openInEditor (external.go) und im Fall X; der vorgeschaltete switch in modeDetail ist weg. Test prueft zusaetzlich E.
- **2026-10-04 22:36 · Alexander Sacharov** — critique (2. Durchgang): nur der Befund aus Durchgang 1 und cb4bb84 gelesen - Schutz sitzt in startEdit/openInEditor/X, Tastenliste entfernt. Erledigt, keine neuen Befunde.
- **2026-10-04 22:38 · Alexander Sacharov** — AlSa 05.10.2026: kommt in 0.3.4, nicht in 0.3.3. Der Zweig ging von der alten release/0.3.3 ab und wird vor dem PR auf master gesetzt.
