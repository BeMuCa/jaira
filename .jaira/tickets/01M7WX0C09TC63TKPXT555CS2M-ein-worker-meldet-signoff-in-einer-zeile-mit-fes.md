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
updated-at: 2026-11-21T19:18:09Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2279502
claimed-at: 2026-11-21T19:00:51Z
outcome-what: "jaira-worker and jaira-dispatcher end their report with a fixed status line 'jaira: <ID> <from> -> <to> · branch · commits · gates · notes · verdict', documented in their Report sections; the jaira-cockpit's digest.ts reads that line (statusLineOf) and shortens the report to its fields, falling back to the old guessing without it. Frame tags sharing the line are stripped; a question is appended only when the shown line does not already hold it, and ends at a '·'"
outcome-why: "The teamlead read the ticket after every report to find branch, commits and whether it could push; a fixed line answers that from the report alone"
outcome-resolves: "55CS2M; shares the line prefix with MVR4PY. Checked: concurrency n/a (pure function); error path: no line -> old path, failure text keeps the report whole; unusual input: CRLF, lowercase, unicode arrow, '-' values, unknown keys, overlong fields clipped, verdict holding '·', '->' or a question, quoted duplicate line, frame tags on either side; permissions: peer text only displayed, every field clipped"
review-summary: none
round: 4
back-from: testing
back-by: lane
back-count: 2
review-gaps: "none (round 2: the repair adds two frame regexes and a raw field, each used once; nothing to fold). Round 1: dropped the redundant verdict check and its cast, clip() reused; laneOf/verdictOf fallback kept by decision (MVR4PY removes it)."
test-verdict: "pass (round 2, diff since da4ea383: digest.ts, cockpit.test.ts, manifest): claude plugin test 241 pass RC=0, go test core/shipped core/mod RC=0, shipped-history --check ok. Both hunter findings fixed as a class: frame tags on either side of the status line (closing </result>, one-line <agent-message ...>…</agent-message>) and a question inside the verdict shown once while a question above the line still shows - each a lineCases test. DoD 2-3 proofs re-pointed (digest.ts:131/:161). Pure function, writes nothing."
hunt-verdict: "fail: bughunt done at 23ce7aaa, 2 change: removing the raw status line before questionOf drops questions the line does not show, and removes the first identical line instead of the last"
---

# Ein Worker meldet signoff in einer Zeile mit festen Feldern: Branch, Commits, Guard, offene Notes

## Definition of Done

- [x] jaira-worker und jaira-dispatcher schreiben die Schlussmeldung im festen Format, dokumentiert im Skill
  proof: core/role/builtin/jaira-worker/SKILL.md:409 (Report: fixed line + field meanings), core/role/builtin/jaira-dispatcher/SKILL.md:556 (same line, points to jaira-worker's Report)
- [x] Das Cockpit liest die Felder aus dieser Zeile
  proof: core/mod/builtin/jaira-cockpit/hooks/digest.ts:131 statusLineOf (frame tags stripped per line), used in digest() at :161
- [x] Eingabeklassen der Zeile im Cockpit-Test: vollständige Zeile, leere Werte '-', fehlender Schlüssel, CRLF, Zeile nicht am Ende, zwei Zeilen (die letzte zählt), Fehlertext bleibt trotz Zeile ganz
  proof: core/mod/builtin/jaira-cockpit/tests/cockpit.test.ts:70 lineCases ('status line: ...' tests incl. frame tags, questions once, 244 pass)

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
- **2026-11-21 19:07 · Alexander Sacharov** — critique: (1) digest.ts hardcodes the labels commits/gates/notes - move them into texts.ts filterWords like marker/question/fullReport. (2) jaira-dispatcher SKILL.md copies the status-line block verbatim from jaira-worker - replace it with a pointer to jaira-worker's Report section, the pattern the dispatcher already uses for 'Before you push'.
- **2026-11-21 19:11 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/digest.ts:131 - class: verdict takes everything to the end of the physical line, so an engine frame closed on the same line leaks into it (shell is digest(report) in register.tsx:230 on the whole notification/agent-message, not on body()); example: '...· verdict: approve</result>' or '...· verdict: approve</agent-message>' is shown as 'approve</result>' / 'approve</agent-message>' in the chat line; the test fixtures always put the closing tag on its own line, the inline form '<agent-message from="x">text</agent-message>' is already used at cockpit.test.ts:259; ran: node --experimental-strip-types scratch calling digest() on both framings; not a blocker (report kept whole, only the shown verdict is wrong)
- **2026-11-21 19:12 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/digest.ts:168 - class: questionOf() still scans the whole text, including the status line, so a verdict that carries a question (the worker skill allows any text there) is shown twice; example: 'verdict: approve · question: shall we?' gives '... · approve · question: shall we? · feat/55cs2m · ... · question: shall we?'; ran: node --experimental-strip-types scratch calling digest(); edge, not a blocker
- **2026-11-21 19:12 · Alexander Sacharov** — bughunt: done at da4ea383 - 2 findings (2 change, 0 pre-existing) in 4 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/digest.ts core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/cockpit.test.ts core/release/notes/55CS2M.md core/role/builtin/jaira-dispatcher/SKILL.md core/role/builtin/jaira-teamlead/SKILL.md core/role/builtin/jaira-worker/SKILL.md core/shipped/shipped.sha256
- **2026-11-21 19:12 · Alexander Sacharov** — testing: class: the status line is read from the raw physical line, so engine frame tags on the same line (<agent-message ...>…</agent-message>, </result>) either leak into the last field or keep the line from matching at all; example: hunter's 'verdict: approve</result>'; fix: strip leading opening and trailing closing frame tags per line before matching, test over both sides and the single-line inline frame. Second: questionOf must not read the status line itself.
- **2026-11-21 19:15 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/digest.ts:176 - class: question text that sits inside the status line but is not shown by the digest is now dropped, because the whole raw line is removed before questionOf, not only what the line displays; places: (a) a verdict longer than 60 chars is clipped, so the question tail is gone from the chat; (b) a question under a key outside FIELDS/verdict (before verdict) vanishes entirely; example: 'jaira: 55CS2M critique -> human · notes: 1 · question: keep cache? · verdict: approve' digests to '... approve · notes 1' with no question (at da4ea383 it showed 'question: keep cache?'); 'verdict: question: shall we keep the on-disk cache when the ref moves underneath it, or drop it?' shows only '...when the ref move…'; not a blocker (format puts questions in verdict, start of question still visible); ran: node --experimental-strip-types scratch script against HEAD and da4ea383 digest.ts
- **2026-11-21 19:15 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/digest.ts:176 - class: text.replace(status.raw, '') removes the FIRST occurrence of the line text, while statusLineOf picks the LAST matching line; any earlier line containing the same text (a quoted worker status line, '> jaira: ...', a repeated line) is removed instead and the real line stays; example: a report quoting '> jaira: 55CS2M critique -> human · verdict: question: keep cache?' and ending with the same line digests to '... question: keep cache? · question: keep cache?'; edge, not a blocker; fix direction: drop the line by index (lines i) rather than by string match; ran: node --experimental-strip-types scratch script against HEAD digest.ts
- **2026-11-21 19:15 · Alexander Sacharov** — bughunt: done at 23ce7aaa - 2 findings (2 change, 0 pre-existing) in 2 modules, follow-up hunt since da4ea383 (diff only; no new files); concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/digest.ts core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/cockpit.test.ts core/release/notes/55CS2M.md core/role/builtin/jaira-dispatcher/SKILL.md core/role/builtin/jaira-teamlead/SKILL.md core/role/builtin/jaira-worker/SKILL.md core/shipped/shipped.sha256
- **2026-11-21 19:16 · Alexander Sacharov** — testing: class: round-2 rule 'read the question from the report minus the raw status line' loses or doubles questions (unknown key, clipped verdict, quoted duplicate line). decision: delete the rule instead of patching it - questionOf reads the whole text as before, and the question is only appended when the shown line does not already contain it. No raw field.
