---
id: 01M44AQ1BH7H7R09PGER2MJBJM
title: "Eine Rolle fasst einen Zeitraum aus dem Logbook fuer die Projektleitung zusammen, ohne Code"
status: done
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
commits:
  - 6555760c14d651b6fcd6eaeee877c8c6d6591311
  - 43aefa9bfadb156a9b2d0077cf1ee96cf846d7fe
  - 23704c35d168bf5534e0ef05f5c685aa4d07659d
  - 79125e6054597a1b4dc153e0d1ed5c1446f11fe4
  - 4bced36f8aa5f9c18b99c404e9ac2b09472e5fc6
created-at: 2026-10-04T20:47:00Z
updated-at: 2026-10-04T22:27:15Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-32325
claimed-at: 2026-10-04T21:03:09Z
outcome-what: "Wie Runde 4; dazu Stunden je Ticket in collect.py (hours_by_ticket je Tag und Zeitraum) und die Regel in SKILL.md, sie nur auf Wunsch und unveraendert zu zeigen."
outcome-why: "AlSa will Stunden je Ticket; eine feste Rechnung statt Schaetzung des Modells."
outcome-resolves: "Wie Runde 4."
review-summary: none
review-gaps: "none - split_hours() recomputes the day's minutes the session loop also sums, but it needs them per ticket and the duplicate is one sum; left as is"
test-verdict: "pass: go build ./... ok; go test -count=1 -race ./core/role/ ./core/release/ ./internal/wintrap/ green; collect.py requirementsgenie KW40 rc=0: total 37.5 = sum(hours_by_ticket) 37.5, each of 7 days equal; split_hours on a scratch session (lead-in, two-id commit, no-id commit) gives A 0.75/B 0.25/- 0.5 = 1.5; build.py rg-summary/rg-jira rc=0; DoD 1-9 in the tree (no Jira in role md files)"
question: "Abnahme Runde 3: Stunden je Ticket (requirementsgenie KW 40: 37,5 h, z. B. 7KSKMP 7,0 / R33C4B 3,0) - passt die Verteilung nach Zeit seit dem vorigen Commit, und sollen sie nur auf Wunsch in summary.md erscheinen? Ja -> weiter; nein -> was anders."
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
- [x] collect.py verteilt die Stunden eines Tages auf die Tickets seiner Commits (Zeit seit dem vorigen Commit, Viertelstunden nach groesstem Rest); die Tickets eines Tages ergeben genau dessen Stunden, und gezeigt werden sie nur, wenn summary.md es verlangt
  proof: core/role/builtin/jaira-role-logbook-summary/collect.py split_hours(); Probe requirementsgenie KW 40: 37,5 h, Summe je Ticket 37,5, jeder Tag gleich

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
- **2026-10-04 21:30 · Alexander Sacharov** — critique (pass 4, c7f59e2 only): comments[] is one optional list in template.html and to_md, covered by the recursive no-code check in texts(); Jira gone from shipped text without losing the example's substance. No finding.
- **2026-10-04 21:38 · Alexander Sacharov** — AlSa 04.10.: Stunden je Ticket gewuenscht. Rechnen gehoert in collect.py (sonst schaetzt das Modell jedes Mal anders), Zeigen in die summary.md der Person - der Standardbericht bleibt kurz. AlSas ~/.jaira/summary.md verlangt jetzt Tabelle je internem Ticket und Summe je Jira-Ticket.
- **2026-10-04 21:38 · Alexander Sacharov** — critique (pass 5, bbf3a1e only): split_hours() sits beside the session code that already owns the commit times, output is additive (hours_by_ticket per day and period), the role shows it only on request. No finding.
- **2026-10-04 22:27 · Alexander Sacharov** — AlSa 05.10.2026: angenommen ('сдвинь 2MJBJM и закончи'), kommt in 0.3.3.
