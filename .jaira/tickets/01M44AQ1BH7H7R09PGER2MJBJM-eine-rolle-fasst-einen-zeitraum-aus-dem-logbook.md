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
updated-at: 2026-10-04T21:30:30Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-32325
claimed-at: 2026-10-04T21:03:09Z
outcome-what: "Wie Runde 2; dazu Kommentare je Tag pro Eintrag (Seite und Markdown) und eine Rolle ohne festes Zielsystem."
outcome-why: "Abnahme AlSa: Tageskommentare fehlten, und das Zielsystem gehoert in die Datei der Person, nicht in die Rolle."
outcome-resolves: "Wie Runde 2."
review-summary: none
review-gaps: "none - pass 2 after rebase: no new diff beyond f6e0805, which optimize itself wrote; earlier gaps note stands"
test-verdict: "pass: go build ./... ok; go test -count=1 -race ./core/role/ ./core/release/ ./internal/wintrap/ green; DoD 1-6 checked in the tree (role in wantRoles/roleFiles, SKILL default, summary.example.md, collect --gap/--lead-in, build check(), NOTES Unreleased); collect.py on requirementsgenie KW40 rc=0, 37.0 h, 28.09. 0 h; build.py page rg-summary (8 items) and rg-jira (4 items) rc=0"
question: "Abnahme: passen die zwei Probeseiten (Standard KW 40 https://claude.ai/artifact/XUH5fQY3DHV4VY1pyB6P5m, Jira https://claude.ai/artifact/3U3Bb2EUnB8wqs1d7WtbiW) zu dem, was die Projektleitung lesen soll? Ja -> weiter; nein -> was fehlt oder stoert."
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
- [x] Ein Eintrag kann Kommentare je Tag tragen, jeder mit eigenem Kopierknopf und im Markdown unter dem Eintrag
  proof: core/role/builtin/jaira-role-logbook-summary/template.html .comments, core/role/builtin/jaira-role-logbook-summary/build.py to_md(); Probe https://claude.ai/artifact/3U3Bb2EUnB8wqs1d7WtbiW Version 2
- [x] Die Rolle nennt kein bestimmtes Zielsystem (kein Jira in SKILL.md, README, NOTES, summary.example.md); was ein System braucht, steht nur in der summary.md der Person
  proof: grep -i jira in core/role/builtin/jaira-role-logbook-summary/*.md und der neuen README-Zeile: 0 Treffer; AlSas Jira-Regeln stehen in ~/.jaira/summary.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-04 21:04 · Alexander Sacharov** — in-progress: Zweig auf feat/747HER gestellt, weil beide core/role/role_test.go (roleFiles) aendern; PR erst nach 747HER. Stunden: Sitzungen laufen ueber den ganzen Zeitraum, damit Arbeit nach Mitternacht eine Sitzung bleibt; sie zaehlen am Tag ihres Beginns. collect.py liest nur (git log --branches --no-merges, jaira list --json, .jaira/logbook); 'active' nimmt jedes Ticket, das sich im Zeitraum bewegt hat, auch die anderer Personen - die Rolle filtert beim Schreiben. Probe KW 40 auf requirementsgenie: 37 h geschaetzt, 28.09. ohne Commits.
- **2026-10-04 21:05 · Alexander Sacharov** — critique (pass 1): core/role/builtin/jaira-role-logbook-summary/collect.py:40 frontmatter() is a hand-written YAML-subset parser for logbook files; jaira show <handle> --json already reads a filed ticket (fields filed_away, path, goal, outcome, tags, follows, commits) - list the files with jaira logbook --since 0 --json, keep the folder-date filter, and read each with jaira show instead, then delete frontmatter()
core/role/builtin/jaira-role-logbook-summary/build.py:30 texts() lists fields by hand and so misses title, eyebrow, item titles, meta, status_label, figure labels and hours.note, while checking hours rows day/note that the summary.json schema in SKILL.md does not have - walk every string in the data recursively and skip only appendix (and hours.rows cells), so the no-code check (DoD 5) covers whatever the page shows
- **2026-10-04 21:06 · Alexander Sacharov** — in-progress Runde 2: Critique umgesetzt - collect.py liest das Logbook ueber 'jaira logbook --since 0 --json' und 'jaira show <handle> --json' statt eigenem YAML-Parser; build.py prueft jeden sichtbaren Text rekursiv (Titel, Eyebrow, Meta, Labels, Notizen), ausgenommen appendix, lang, status und die Stunden-Zellen. Probe: KW 40 und Jira-Beispiel bauen weiter, Titel mit build.py und Meta mit Hash werden abgewiesen.
- **2026-10-04 21:07 · Alexander Sacharov** — critique (pass 2): both pass-1 findings addressed in bbb907d - collect.py reads the logbook via jaira logbook --json + jaira show, frontmatter() gone; build.py texts() walks every string except appendix/lang/status/hours.rows. Nothing new that breaks the DoD.
- **2026-10-04 21:09 · Alexander Sacharov** — testing: go test -count=1 ./internal/wintrap/ FAIL - "core/role/role.go:33: windows trap 2: //go:embed pulls in 6 file(s) that no line in .gitattributes pins, starting at core/role/builtin/jaira-role-acceptance/build.py". Three of the six are this ticket (jaira-role-logbook-summary/collect.py, build.py, template.html), three come from 747HER (acceptance build.py, template.html, report.html). Fix: in .gitattributes beside core/role/builtin/**/*.md add "core/role/builtin/**/*.py text eol=lf" and "core/role/builtin/**/*.html text eol=lf" - covers both roles. Otherwise green: go build ./..., go test -race ./core/role/ ./core/release/ ok; collect.py on requirementsgenie-R33C4B KW40 rc=0, 37.0 h, 28.09. without commits, 1 done (753FQ5), sessions over midnight count on the start day; build.py page rg-summary (8 items) and rg-jira (4 items) rc=0; check() catches path, file name, hash, snake_case, function call in title/meta/hours.note and skips appendix.
- **2026-10-04 21:09 · Alexander Sacharov** — in-progress Runde 3: Testing-Befund (wintrap, .py/.html nicht auf LF) war derselbe wie bei 747HER und ist dort in 933ad03 behoben; feat/2MJBJM auf feat/747HER neu aufgesetzt, kein eigener Code. internal/wintrap gruen.
- **2026-10-04 21:09 · Alexander Sacharov** — critique (pass 3): the testing finding (unpinned .py/.html) is answered by the rebase onto 933ad03 (747HER pins core/role/builtin/**/*.py and *.html to LF); no new diff in this ticket.
- **2026-10-04 21:30 · Alexander Sacharov** — Abnahme AlSa 04.10.: im Jira-Beispiel fehlten die Kommentare je Tag - es gab nur Summary und Description, weil ein Eintrag keine Kommentare tragen konnte. Und: Jira gehoert nicht in die Rolle, jeder schreibt sein System in seine summary.md. Umgesetzt: comments[] je Eintrag, Jira aus SKILL.md/README/NOTES/Beispiel entfernt, AlSas Regeln (nach KW34) in ~/.jaira/summary.md.
