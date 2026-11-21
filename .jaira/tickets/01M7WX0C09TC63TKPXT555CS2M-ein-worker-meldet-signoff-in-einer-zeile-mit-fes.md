---
id: 01M7WX0C09TC63TKPXT555CS2M
title: "Ein Worker meldet signoff in einer Zeile mit festen Feldern: Branch, Commits, Guard, offene Notes"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die Schlussmeldung eines Workers/Dispatchers hat ein festes Format, das ein Teamlead ohne Nachfragen und ohne das Ticket zu lesen auswerten kann."
context: |-
  Heute: Schlussmeldungen sind freier Text in drei Zeilen; der Teamlead liest das Ticket nach, um Branch, Commits und ob gepusht werden kann zu finden. Alex: "sonst liest er viel".
  Vorschlag des Teamleads: eine Zeile mit Branch, Commits (Spanne), Guard/Lint ok, offene Notes (Anzahl), Lane.
  Verwandt: MVR4PY (Cockpit liest Lane und Urteil aus einer festen Schlusszeile).
definition-of-done: "jaira-worker und jaira-dispatcher schreiben die Schlussmeldung im festen Format, dokumentiert im Skill"
tags: []
blocked-by:
  - 01M7WX07G4E8ERZ30PH10J88M5
related:
  - 01M67RXSPRWC7WVBCCM7MVR4PY
commits: []
created-at: 2026-11-20T16:35:33Z
updated-at: 2026-11-21T19:07:02Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2279502
claimed-at: 2026-11-21T19:00:51Z
outcome-what: "jaira-worker and jaira-dispatcher end their report with a fixed status line 'jaira: <ID> <from> -> <to> · branch · commits · gates · notes · verdict', documented in their Report sections; the jaira-cockpit's digest.ts reads that line (statusLineOf) and shortens the report to its fields, falling back to the old guessing without it"
outcome-why: "The teamlead read the ticket after every report to find branch, commits and whether it could push; a fixed line answers that from the report alone"
outcome-resolves: "55CS2M; shares the line prefix with MVR4PY, which can now remove the guessing. Checked: concurrency n/a (pure function); error path: no line -> old path, failure text still keeps the report whole; unusual input: CRLF, lowercase handle/lanes, unicode arrow, '-' values, unknown keys, overlong fields clipped, verdict holding '·' or '->'; permissions: peer text is only displayed, every field clipped to 60"
---

# Ein Worker meldet signoff in einer Zeile mit festen Feldern: Branch, Commits, Guard, offene Notes

## Definition of Done

- [x] jaira-worker und jaira-dispatcher schreiben die Schlussmeldung im festen Format, dokumentiert im Skill
  proof: core/role/builtin/jaira-worker/SKILL.md:409, core/role/builtin/jaira-dispatcher/SKILL.md:556 (Report section, fixed line + field meanings)
- [x] Das Cockpit liest die Felder aus dieser Zeile
  proof: core/mod/builtin/jaira-cockpit/hooks/digest.ts:124 statusLineOf, used in digest() at :154
- [x] Eingabeklassen der Zeile im Cockpit-Test: vollständige Zeile, leere Werte '-', fehlender Schlüssel, CRLF, Zeile nicht am Ende, zwei Zeilen (die letzte zählt), Fehlertext bleibt trotz Zeile ganz
  proof: core/mod/builtin/jaira-cockpit/tests/cockpit.test.ts:70 lineCases ('status line: ...' tests, 237 pass)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] digest.ts: statusLineOf() parses the fixed 'jaira: <ID> <from> -> <to> · key: value ...' line; digest() prefers it over guessing
- [x] cockpit tests: table test over input classes (full line, '-' values, missing keys, CRLF, line not last, two lines -> last wins, failure text still stays whole)
- [x] jaira-worker + jaira-dispatcher SKILL.md: Report section documents the line, field meanings and order
- [x] release note core/release/notes/55CS2M.md; regenerate shipped manifest if required; run cockpit tests + go test for roles

## Progress
- **2026-11-20 16:49 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-21 18:58 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/4205.sock
- **2026-11-21 18:58 · Alexander Sacharov** — decision (worker, no open question left): 1) Format, last line of the report: 'jaira: <ID> <from> -> <to> · branch: <branch> · commits: <first>..<last> (<n>) · gates: ok|red|none · notes: <n> · verdict: <text>'; fixed keys, fixed order, '-' for an empty value. Teamlead's proposal plus MVR4PY's 'jaira: <ID> <from> -> <to>' prefix, so both tickets share one line. 2) notes = count of 'note:' lines in review-gaps. gates = the project's local gates before push: ok/red, none where nothing was run (no code). 3) Cockpit: digest.ts reads lane, verdict and the four fields from that line when it is there; without it the old guessing stays — removing it is MVR4PY's job. 4) Scope: jaira-worker and jaira-dispatcher as the DoD names; jaira-role-lane reports to its dispatcher, which writes the line.
- **2026-11-21 19:00 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock (restarted worker; previous claimer process gone)
- **2026-11-21 19:03 · Alexander Sacharov** — statusLineOf in digest.ts: last matching line wins (a report quoting an earlier one ends with its own); verdict is last and takes the rest of the line so it may hold '·' or '->'; ticket from the line outranks the handle count. isFailure still runs first, so 'gates: red' alone shortens (the line shows 'gates red'), a '3 failed' in the text keeps it whole.
