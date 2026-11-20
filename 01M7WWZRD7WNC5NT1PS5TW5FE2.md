---
id: 01M7WWZRD7WNC5NT1PS5TW5FE2
title: "Randfaelle werden vor dem Code benannt, und wer umsetzt, prueft die Fehlerklasse statt des Beispiels"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "pre-process schreibt Fehlerklasse, Eingabeklassen und bei Heuristiken eine Annahmeregel als DoD-Punkte; in-progress testet die Klasse, sucht Geschwister-Stellen und prueft selbst die vier Kategorien, bevor es abgibt; bughunt meldet einen Fund als Klasse."
context: |-
  Problem: Review und bughunt finden in jedem Ticket etwas, Schleifen laufen 10-24 Runden. Die Funde sind meist richtig.
  Analyse 20.11.2026 an 34 Tickets eines Projekts: bughunt 47 Funde aus der Aenderung selbst, 5 vorbestehend; critique fast immer "nichts"; Review am Ende 0 Blocker, 2-6 Notes.
  Ursache 1: die vier Kategorien (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte) stehen nur in testing/bughunt/review (.jaira/lanes/review.md:31). pre-process, in-progress (Tier cheap, "work only within the plan's scope") und critique ("Do not check whether it works") fragen nie danach.
  Ursache 2: ein Fund kommt als ein Beispiel, der Fix testet das Beispiel, die naechste Runde findet den Nachbarfall (wiederholt beobachtet: Seitenzahl "3 / 12" dann "3/12"; Zahlen "0,05", dann "2,50", dann "1.234,56").
  Ursache 3: Heuristik-Tickets (PDF, Zahlen) haben keine Annahmeregel; die Schleife konvergierte erst, als der Teamlead eine setzte ("Textverlust ist schlimmer als ein uebrig gebliebener Kopf/Fuss").
  Ursache 4: Tickets beschreiben das Symptom ("NUL in Werkzeug A"), nicht die Klasse ("alles, was die DB nicht speichert, in jedem Schreibwerkzeug").
  Ursache 5: an der Wurzel reparieren statt am Beispiel hat den Nachbarfall-Zyklus beendet (strukturierte Daten statt String parsen).
  Alex: pre-process baut die DoD so, dass in-progress weiss, was zu tun ist; aber wer umsetzt, muss selbst denken statt Haken setzen, sonst pruefen ihn danach zehn andere.
  Eine Explain-Taste nur fuer blocker: gehoert zum Ticket ueber die Schleifenbremse.
definition-of-done: "pre-process-Prompt verlangt: Fehlerklasse statt Symptom, Eingabeklassen, bei Heuristik eine Annahmeregel - als DoD-Punkte"
tags:
  - gates
blocked-by: []
related:
  - 01M7WWZWSDHK37DGGHJFDF0YD6
  - 01M67RWZDWBS8E44Z3KGXKS048
  - 01M67K7JAF6R9N2VSPEW3CBCN6
  - 01M5ZJ4NNCYN6YBWVAPWAHG254
commits: []
created-at: 2026-11-20T16:35:12Z
updated-at: 2026-11-20T16:53:41Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-42185
claimed-at: 2026-11-20T16:42:15Z
outcome-what: "Eingebaute Lanes pre-process, in-progress, testing, review und die Rolle jaira-role-bughunt: pre-process schreibt Fehlerklasse, Standard-Eingabeklassen (leer, gross, Kodierung/CRLF, Leerzeichen-Pfade, Windows, zwei Prozesse, Abbruch) und Annahmeregel als DoD-Punkte; in-progress repariert die Ursache, testet die Klasse, zieht Geschwister mit, prueft die vier Kategorien, schliesst kleine Raender in derselben Runde und fragt nach der dritten Rueckgabe derselben Lane; testing/bughunt melden Klasse + Beispiel; testing und review schicken nur fuer Blocker zurueck (falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke), 'note: pre-existing:' fuer Altes. .jaira/lanes nachgezogen, shipped.sha256, Release-Note."
outcome-why: "Review- und bughunt-Schleifen liefen 10-24 Runden, weil jede Runde den Nachbarfall eines nur am Beispiel reparierten Fehlers fand und die vier Kategorien vor dem Code nie gefragt wurden."
outcome-resolves: "DoD 1-4: pre-process 15-pre-process.md:38, in-progress 20-in-progress.md:29, bughunt SKILL.md:104 und testing 27-testing.md:102, eingebettet mit shipped.sha256 und core/release/notes/TW5FE2.md"
review-summary: |-
  core/lane/builtin/27-testing.md:102 'a note: line for the review' - testing schreibt kein review-gaps; sagen: eine Note testing: note: ... am Ticket, hunt-verdict bleibt pass und zaehlt sie
  core/lane/builtin/20-in-progress.md:68 'a review blocker:' - review schreibt seine Rueckgabe als Note ohne festes Praefix; 'a review note' schreiben
round: 3
review-gaps: "none: kein Code, nur Prompts; die vier Kategorien stehen jetzt auch in pre-process und in-progress - bewusst, jede Lane liest nur ihren eigenen Prompt; keine toten Stellen"
test-verdict: "fail: DoD 5-7 (Interim-Regeln des Teamleads, nachgetragen) sind noch nicht in den Lanes; Suite gruen (go test -race ./... RC=0)"
---

# Randfaelle werden vor dem Code benannt, und wer umsetzt, prueft die Fehlerklasse statt des Beispiels

## Definition of Done

- [x] pre-process-Prompt verlangt: Fehlerklasse statt Symptom, Eingabeklassen, bei Heuristik eine Annahmeregel - als DoD-Punkte
  proof: core/lane/builtin/15-pre-process.md:38 (Name the edges before the code: class, input classes, acceptance rule, jaira dod --add)
- [x] in-progress-Prompt verlangt: Tests ueber die Klasse (parametrisiert), Suche nach Geschwister-Stellen, Selbstpruefung der vier Kategorien vor dem Move, Fix an der Ursache statt am Beispiel
  proof: core/lane/builtin/20-in-progress.md:29 (Fix the class, not the example: cause, parametrised test, siblings in scope, four categories before move)
- [x] bughunt- und testing-Prompt: ein Fund nennt die Klasse und das Beispiel; der Rueckweg nach in-progress verlangt den Fix der Klasse
  proof: core/role/builtin/jaira-role-bughunt/SKILL.md:104; core/lane/builtin/27-testing.md:110,119 (class + example; passed only when the class is fixed)
- [x] Die eingebetteten Lanes und Rollen sind geaendert, mit Release-Note
  proof: core/shipped/shipped.sha256:560-569; .jaira/lanes/{pre-process,in-progress,testing,review}.md byte-equal to builtin; core/release/notes/TW5FE2.md; go test ./core/shipped ./core/lane ./core/role/... ok
- [x] pre-process nennt die Standard-Eingabeklassen: leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg
  proof: core/lane/builtin/15-pre-process.md:48-55
- [x] Schleifenbremse in testing, review und in-progress: Blocker nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke; Folgerunde testet und huntet nur den Diff seit dem letzten Hunt; nach 3 Rueckgaben aus derselben Lane eine Frage statt Runde 4
  proof: core/lane/builtin/27-testing.md:97 (blocker = four kinds), 27-testing.md:125-128 (follow-up round diff only, then review); core/lane/builtin/40-review.md:57; core/lane/builtin/20-in-progress.md:66 (third return -> human question)
- [x] Kleine Raender der Aenderung aus review-gaps werden in derselben Runde geschlossen statt als note: zur Person zu gehen; Vorbestehendes ausserhalb des Umfangs steht als 'note: pre-existing: ...'
  proof: core/lane/builtin/20-in-progress.md:58 (close small edges in the same round); core/lane/builtin/40-review.md:61-66 (note names class and fix; note: pre-existing:)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-20 16:35 · Alexander Sacharov** — Alex 20.11.2026: er schickt inzwischen jedes Ticket wegen der Review-Funde zurueck, bis die Raender zu sind ('Raender offen lassen ist schlecht'). Jede Note ist also faktisch eine weitere Runde mit dem Menschen - die Raender billiger vor dem Code zu benennen spart genau diese Runden.
- **2026-11-20 16:42 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/35112.sock
- **2026-11-20 16:43 · Alexander Sacharov** — decision: keine offene Entscheidung fuer die Person. Selbst gesetzt: (1) pre-process laeuft nur mit Option planning (requires-option) - deshalb benennt in-progress Klasse/Eingabeklassen/Annahmeregel selbst, wenn kein Plan-Lane lief; (2) 'Work only within the plan's scope' bleibt, aber Geschwister-Stellen derselben Fehlerklasse gelten ausdruecklich als im Scope, andere Nachbarprobleme weiter nur als Note; (3) pre-process schreibt die Punkte mit 'jaira dod <id> --add' in die DoD, nicht nur in den Plan, weil nur die DoD das Gate liest; (4) critique, review und Tier bleiben unveraendert - nicht in der DoD; (5) jede geaenderte eingebaute Datei bekommt ihren sha256 in core/shipped/shipped.sha256, .jaira/lanes bekommt dieselbe Fassung.
- **2026-11-20 16:46 · Alexander Sacharov** — critique: testing.md:110 als Pass-Regel formulieren (pass nur, wenn auch ein Nachbar geht); in-progress.md:35 'it' durch 'the definition of done' ersetzen. Sonst nichts: vier Kategorien bewusst in jeder Lane wiederholt, weil jede Lane nur ihren eigenen Prompt liest.
- **2026-11-20 16:50 · Alexander Sacharov** — Interim-Regeln vom Teamlead (Alex 20.11.2026, gelten bis TW5FE2 ausgeliefert): 1 pre-process: DoD nennt Fehlerklasse statt Symptom als DoD-Punkte; Eingabeklassen leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg; bei Heuristik Annahmeregel. 2 in-progress: vor dem Move vier Kategorien selbst pruefen; Tabellentest ueber die Klasse; Geschwister per grep mitziehen; Fix an der Ursache; Raender selbst denken statt Haken setzen. 3 bughunt/testing: Fund nennt Klasse und Beispiel, Rueckweg verlangt Fix der Klasse. 4 Schleifenbremse: Blocker nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke; nach Fix-Durchgang Testlauf + Bughunt nur ueber Diff seit letztem Hunt, dann review; nach 3 Rueckgaben aus derselben Lane Frage an den Teamlead. 5 kleine Raender aus review-gaps in derselben Runde schliessen statt note: stehen lassen; Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'. Fuer TW5FE2: diese fuenf Punkte sind der Kern dessen, was in die eingebetteten Lanes kommt.
- **2026-11-20 16:50 · Alexander Sacharov** — testing: Suite gruen (go test -race -count=1 ./... RC=0), DoD 1-4 im Baum; zurueck nach in-progress, weil die Interim-Regeln als DoD 5-7 nachgetragen sind. Laufender bughunt deckt nur ab6eeb90.
- **2026-11-20 16:50 · Alexander Sacharov** — bughunt: pre-existing: core/shipped/shipped.sha256 (read by core/role/install.go:147 writeFile) - class: versions that went out in a release tag are missing from the hand-merged manifest, so 'jaira update --roles' takes an untouched install of them for a person's edit and skips it (Skipped, needs --force); example: core/mod/builtin/jaira-cockpit/hooks/board.ts as shipped in v0.5.2/v0.5.3 (sha256 a06a3c1c...), and 7 more tagged versions of register.tsx, texts.ts, types/index.d.ts, board.ts and jaira-role-acceptance/template.html; the files TW5FE2 touches are complete; never in the manifest (git log -S<hash> on shipped.sha256 is empty); ran: scripts/shipped-history.sh vs the committed manifest (59 lines only in the generated output, none only in the committed one), then checked each against git ls-tree of the v* tags (8 shipped)
- **2026-11-20 16:50 · Alexander Sacharov** — bughunt: done at ab6eeb90 - 1 finding (0 change, 1 pre-existing) in 5 modules, first hunt; concurrency, error paths, input, permissions checked (dod --add on pre-process and two parallel dod --add on a scratch board: both land; CLI flags named by the prompts exist; lane copies equal builtin; manifest has every version of the touched files; go test core/shipped core/lane core/role core/release green); files: .jaira/lanes/in-progress.md .jaira/lanes/pre-process.md .jaira/lanes/testing.md core/lane/builtin/15-pre-process.md core/lane/builtin/20-in-progress.md core/lane/builtin/27-testing.md core/release/notes/TW5FE2.md core/role/builtin/jaira-role-bughunt/SKILL.md core/shipped/shipped.sha256
