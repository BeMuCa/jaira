---
id: 01M44AQ1BH7H7R09PGER2MJBJM
title: "Eine Rolle fasst einen Zeitraum aus dem Logbook fuer die Projektleitung zusammen, ohne Code"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "jaira roles install liefert jaira-role-logbook-summary: sie fasst einen Zeitraum (Standard: letzte Kalenderwoche) als Seite zusammen, die eine Projektleitung liest - was fertig ist, was laeuft, was offen ist, mit Stunden aus den Commits und kopierfertigen Texten; ohne Konfiguration ein schlichter Standard, mit einer Datei im Repo an das eigene Ziel (z. B. Jira) angepasst."
context: |-
  Was heute fehlt:
  - Eine Projektleitung fragt: wie weit seid ihr, was ist fertig, wie viele Stunden. Das Board beantwortet Fragen zum Code, nicht diese.
  - AlSa baut das jede Woche von Hand: ~/jira-erfassung-kw34.html (KW 34, 19.-23.08.2026). Inhalt: Stunden je Tag aus Commits (Sitzungen getrennt bei mehr als 90 Minuten Abstand), Tickets zu zwei Jira-Tickets gebuendelt, je Ticket Summary, Description und ein Kommentar pro Tag in Jira-Wiki-Markup mit Kopierknopf, interne Ticket-IDs als Tabelle, Stunden halb auf zwei Kunden.
  - Das meiste davon ist AlSas eigener Bedarf (Jira, zwei Kunden, Wiki-Markup). Andere brauchen anderes.

  Was schon existiert:
  - jaira logbook legt fertige Tickets unter .jaira/logbook/<initialen>-<datum>/ ab, mit gestempelten Commits; jaira logbook ohne Argument listet die letzten vier Wochen (MJ3PVV).
  - Rollen liegen in core/role/builtin/<id>/ (core/role/role.go:33).

  Entscheidungen:
  - Der Standard ist schlicht: Seite mit Zeitraum, Zahlen (fertig, in Arbeit, Stunden), Themen in der Sprache der Nutzer (was es fuer sie bedeutet, keine Dateinamen), Stand je Thema, ein kurzer Eintrag pro Tag, jeder Text als Markdown mit Kopierknopf.
  - Wer etwas anderes braucht, schreibt es in eine Datei, die die Rolle liest: .jaira/summary.md im Repo, sonst ~/.jaira/summary.md. Darin in Worten: Zielsystem, Format, Sprache, Bündelung, Stundenaufteilung. Ohne Datei gilt der Standard.
  - Stunden sind eine Schaetzung aus Commit-Zeiten und werden auf der Seite so genannt; Tage ohne Commits (z. B. Abnahme) fragt die Rolle nach, statt Stunden zu erfinden.
  - Kein Code, keine Pfade, keine Commit-Hashes im Text fuer die Projektleitung; interne IDs nur als Anhang.
definition-of-done: jaira roles install schreibt jaira-role-logbook-summary
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-04T20:47:00Z
updated-at: 2026-10-04T21:04:58Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-32325
claimed-at: 2026-10-04T21:03:09Z
outcome-what: "Neue eingebaute Rolle jaira-role-logbook-summary: collect.py (Fakten aus git, Board, Logbook, Stunden je Tag), build.py (prueft Text auf Code, baut Seite und Markdown), template.html, summary.example.md; README und NOTES."
outcome-why: "Eine Projektleitung fragt nach Stand und Stunden, nicht nach Code; AlSa baut das bisher jede Woche von Hand (jira-erfassung-kw34.html)."
outcome-resolves: "Ohne Datei ein schlichter Bericht der letzten Woche; mit .jaira/summary.md in eigener Form, z. B. Jira mit zwei Kunden."
---

# Eine Rolle fasst einen Zeitraum aus dem Logbook fuer die Projektleitung zusammen, ohne Code

## Definition of Done

- [x] jaira roles install schreibt jaira-role-logbook-summary
  proof: core/role/role_test.go roleFiles + TestRolesShipTheirFiles, TestBuiltinsAreTheShippedRoles
- [x] Ohne summary.md entsteht fuer die letzte Kalenderwoche eine Seite mit Zahlen, Themen, Stand je Thema, einem Eintrag pro Tag und Kopierknoepfen
  proof: core/role/builtin/jaira-role-logbook-summary/SKILL.md Abschnitt 1 Default; Probe KW 40 requirementsgenie: https://claude.ai/artifact/XUH5fQY3DHV4VY1pyB6P5m
- [x] Mit einer summary.md, die Jira-Wiki-Markup und zwei Kunden verlangt, entsteht eine Seite wie ~/jira-erfassung-kw34.html
  proof: core/role/builtin/jaira-role-logbook-summary/summary.example.md (Jira, zwei Kunden); Probe: https://claude.ai/artifact/3U3Bb2EUnB8wqs1d7WtbiW
- [x] Stunden je Tag stammen aus Commit-Zeiten mit der Sitzungsgrenze aus summary.md (Standard 90 Minuten) und heissen auf der Seite Schaetzung
  proof: core/role/builtin/jaira-role-logbook-summary/collect.py --gap/--lead-in, Sitzungen ueber Mitternacht zaehlen am Starttag; template note 'Schaetzung'
- [x] Der Text fuer die Projektleitung enthaelt keine Dateipfade, Funktionsnamen oder Commit-Hashes
  proof: core/role/builtin/jaira-role-logbook-summary/build.py check(): Pfad, Dateiname, Hash, snake_case, Funktionsaufruf, Backtick -> exit 1 mit Zeile
- [x] core/release/NOTES.md nennt die neue Rolle und die Datei summary.md
  proof: core/release/NOTES.md Unreleased

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-04 21:04 · Alexander Sacharov** — in-progress: Zweig auf feat/747HER gestellt, weil beide core/role/role_test.go (roleFiles) aendern; PR erst nach 747HER. Stunden: Sitzungen laufen ueber den ganzen Zeitraum, damit Arbeit nach Mitternacht eine Sitzung bleibt; sie zaehlen am Tag ihres Beginns. collect.py liest nur (git log --branches --no-merges, jaira list --json, .jaira/logbook); 'active' nimmt jedes Ticket, das sich im Zeitraum bewegt hat, auch die anderer Personen - die Rolle filtert beim Schreiben. Probe KW 40 auf requirementsgenie: 37 h geschaetzt, 28.09. ohne Commits.
