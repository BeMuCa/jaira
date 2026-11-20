---
id: 01M7WWZRD7WNC5NT1PS5TW5FE2
title: "Randfaelle werden vor dem Code benannt, und wer umsetzt, prueft die Fehlerklasse statt des Beispiels"
status: in-progress
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
updated-at: 2026-11-20T17:02:11Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-42185
claimed-at: 2026-11-20T16:42:15Z
outcome-what: "Eingebaute Lanes pre-process, in-progress, testing, review und die Rolle jaira-role-bughunt: pre-process schreibt Fehlerklasse, Standard-Eingabeklassen (leer, gross, Kodierung/CRLF, Leerzeichen-Pfade, Windows, zwei Prozesse, Abbruch) und Annahmeregel als DoD-Punkte; in-progress repariert die Ursache, testet die Klasse, zieht Geschwister mit, prueft die vier Kategorien, schliesst kleine Raender in derselben Runde und fragt nach der dritten Rueckgabe derselben Lane; testing/bughunt melden Klasse + Beispiel; testing und review schicken nur fuer Blocker zurueck (falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke), 'note: pre-existing:' fuer Altes. .jaira/lanes nachgezogen, shipped.sha256, Release-Note."
outcome-why: "Review- und bughunt-Schleifen liefen 10-24 Runden, weil jede Runde den Nachbarfall eines nur am Beispiel reparierten Fehlers fand und die vier Kategorien vor dem Code nie gefragt wurden."
outcome-resolves: "DoD 1-4: pre-process 15-pre-process.md:38, in-progress 20-in-progress.md:29, bughunt SKILL.md:104 und testing 27-testing.md:102, eingebettet mit shipped.sha256 und core/release/notes/TW5FE2.md"
review-summary: none
round: 6
review-gaps: "none: Zaehlregel ersetzt die Praefixliste, nichts Doppeltes oder Verwaistes"
test-verdict: "pass: go test -race -count=1 ./... RC=0, go build ok; in-progress-Kopie byte-gleich"
hunt-verdict: "fail: bughunt done at ad6a6ffc, 1 change - 'round > 1' zaehlt auch Rueckgaben einer Person aus signoff/human als Loop-Rueckgabe"
---

# Randfaelle werden vor dem Code benannt, und wer umsetzt, prueft die Fehlerklasse statt des Beispiels

## Definition of Done

- [x] pre-process-Prompt verlangt: Fehlerklasse statt Symptom, Eingabeklassen, bei Heuristik eine Annahmeregel - als DoD-Punkte
  proof: core/lane/builtin/15-pre-process.md:38 (Name the edges before the code: class, input classes, acceptance rule, jaira dod --add)
- [x] in-progress-Prompt verlangt: Tests ueber die Klasse (parametrisiert), Suche nach Geschwister-Stellen, Selbstpruefung der vier Kategorien vor dem Move, Fix an der Ursache statt am Beispiel
  proof: core/lane/builtin/20-in-progress.md:29 (Fix the class, not the example: cause, parametrised test, siblings in scope, four categories before move)
- [x] bughunt- und testing-Prompt: ein Fund nennt die Klasse und das Beispiel; der Rueckweg nach in-progress verlangt den Fix der Klasse
  proof: core/role/builtin/jaira-role-bughunt/SKILL.md:104; core/lane/builtin/27-testing.md:112,121 (class + example; passed only when the class is fixed)
- [x] Die eingebetteten Lanes und Rollen sind geaendert, mit Release-Note
  proof: core/shipped/shipped.sha256:560-573; .jaira/lanes/{pre-process,in-progress,testing,review}.md byte-equal to builtin; core/release/notes/TW5FE2.md; go test ./core/shipped ./core/lane ./core/role/... ok
- [x] pre-process nennt die Standard-Eingabeklassen: leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg
  proof: core/lane/builtin/15-pre-process.md:48-55
- [x] Schleifenbremse in testing, review und in-progress: Blocker nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke; Folgerunde testet und huntet nur den Diff seit dem letzten Hunt; nach 3 Rueckgaben aus derselben Lane eine Frage statt Runde 4
  proof: core/lane/builtin/27-testing.md:97 (blocker = four kinds), 27-testing.md:127-130 (follow-up round diff only, then review); core/lane/builtin/40-review.md:57; core/lane/builtin/20-in-progress.md:66-87 (own 'back: <lane>, round <n>' lines, third -> human question)
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
- **2026-11-20 16:53 · Alexander Sacharov** — critique: testing.md:102 nicht-blockierende Funde als Note 'testing: note: ...' am Ticket, hunt-verdict pass mit Zaehlung; in-progress.md:68 'a review note' statt 'a review blocker:'.
- **2026-11-20 16:56 · Alexander Sacharov** — bughunt: change: core/lane/builtin/20-in-progress.md:68-71 (and .jaira/lanes/in-progress.md) - class: the return count reads note prefixes that do not mark a return one-to-one; example: in-progress counts every 'testing:' note as a return, but 27-testing.md:103 now writes 'testing: note: ...' for non-blocking findings on a passing verdict, and the hunter writes one 'bughunt: change:' note per finding (non-blockers among them stay notes on a pass), so three findings in one round or three passing testing notes read as a third return and the ticket goes to human after its first real return; the review's blocker note (40-review.md:116-118) has no prefix at all, so its returns cannot be told from in-progress's own progress notes and are never counted; fix: give each return one marker that only a send-back writes (e.g. '<lane>: return: ...' written by critique/testing/optimize/review right before 'jaira move --to in-progress') and count those; ran: read, grep -n 'testing:|critique:|bughunt:' core/lane/builtin/*.md; blame: lines from e101fb64/c20a3e06 (TW5FE2)
- **2026-11-20 16:56 · Alexander Sacharov** — bughunt: change: core/lane/builtin/20-in-progress.md:69 - class: a loop lane missing from the brake's list; example: optimize sends work back with an 'optimize:' note (26-optimize.md:61) and is a declared loop (critique/optimize/testing loops), but the counted prefixes are only critique:, testing:, bughunt: change: and the review's note, so optimize returns never trigger the third-return stop; fix: name optimize in the list, or derive it from the marker in the finding above; ran: read; blame: e101fb64 (TW5FE2)
- **2026-11-20 16:56 · Alexander Sacharov** — bughunt: done at e101fb64 - 2 findings (2 change, 0 pre-existing) in 6 modules, follow-up hunt since ab6eeb90; concurrency, error paths, input, permissions checked (third-return 'jaira move --to human --question' from in-progress passes the gate once outcomes are set, ran on a scratch board with this repo's lanes; lane copies equal builtin; manifest holds current and previous 40-review/15/20/27 hashes; go test core/shipped core/lane green); files: .jaira/lanes/in-progress.md .jaira/lanes/pre-process.md .jaira/lanes/review.md .jaira/lanes/testing.md core/lane/builtin/15-pre-process.md core/lane/builtin/20-in-progress.md core/lane/builtin/27-testing.md core/lane/builtin/40-review.md core/release/notes/TW5FE2.md core/role/builtin/jaira-role-bughunt/SKILL.md core/shipped/shipped.sha256
- **2026-11-20 16:57 · Alexander Sacharov** — back: testing, round 5
- **2026-11-20 16:57 · Alexander Sacharov** — Fix der Klasse 'Rueckgaben aus gedruckten Notes zaehlen': in-progress schreibt beim Aufnehmen einer Rueckgabe selbst eine Zeile 'back: <lane>, round <n>' und zaehlt nur diese - Ursache statt Praefixliste; optimize ist darin mit genannt. Verworfen: Feld round (zaehlt alle Lanes zusammen, haelt gesunde Tickets nach 3 verschiedenen Rueckgaben an) und neue Praefixe in vier Lanes (groesser, gleiche Mehrdeutigkeit bei mehreren Notes pro Runde).
- **2026-11-20 16:59 · Alexander Sacharov** — bughunt: change: core/lane/builtin/20-in-progress.md:72 (and .jaira/lanes/in-progress.md:72) - class: the rule reads 'round above 1' as 'a loop lane sent it back', but core/move/move.go:133 sendsBack bumps round on every move to an earlier lane, including a person's return from signoff or human; the prompt then picks the sender only among 'critique, optimize, testing or review' from the newest notes, so a signoff return (newest note: the review's, no prefix) is logged as 'back: review' and, after two real review returns, makes in-progress escalate to human instead of doing what the person asked; a human answer after the third-return escalation is likewise logged against the loop lane. Example: signoff -> in-progress. ran: scratch board with this repo's lanes, jaira move <id> --to in-progress from signoff (exit 0, round 2->3) and from human (round 3->4)
- **2026-11-20 16:59 · Alexander Sacharov** — bughunt: done at ad6a6ffc - 1 finding (1 change, 0 pre-existing) in 3 modules, follow-up hunt since e101fb64; concurrency, error paths, input, permissions checked (new block identical in builtin and board copy; manifest line b8ef26c5 matches sha256 of 20-in-progress.md; round exists in jaira show --json; 'back:' prefix collides with no other reader); files: .jaira/lanes/in-progress.md .jaira/lanes/pre-process.md .jaira/lanes/review.md .jaira/lanes/testing.md core/lane/builtin/15-pre-process.md core/lane/builtin/20-in-progress.md core/lane/builtin/27-testing.md core/lane/builtin/40-review.md core/release/notes/TW5FE2.md core/role/builtin/jaira-role-bughunt/SKILL.md core/shipped/shipped.sha256
- **2026-11-20 17:01 · Alexander Sacharov** — Frage an den Teamlead: testing hat TW5FE2 zum dritten Mal zurueckgeschickt (1: nachgetragene Interim-DoD 5-7, 2: Rueckgaben nach Notes-Praefixen gezaehlt, 3: 'round > 1' zaehlt auch Rueckgaben aus signoff/human). Runde 3 verbatim: in-progress.md:72 liest round>1 als Loop-Rueckgabe, move.go:133 erhoeht round bei jedem Rueckweg, auch von einer Person. Einschaetzung: DoD stimmt, die Schleife wird nicht tiefer - Runde 1 war Umfangserweiterung, 2 und 3 sind dieselbe Klasse (Zaehlregel), jeder Fix kleiner. Vorschlag: eine Runde noch - Rueckgabe einer Person (signoff/human) wird als 'back: person' notiert, zaehlt nicht und setzt die Zaehlung zurueck; danach review. Optionen: 1) Ja, eine Runde, dann review (empfohlen) 2) Nein, so an review/Person geben, Fund als note
- **2026-11-20 17:01 · Alexander Sacharov** — Teamlead-Antwort 17:01: Option 1, eine Runde. Rueckgabe aus einer Menschen-Lane (signoff, human, jede Lane mit requires-human-exit) und die Antwort einer Person nach einer Eskalation werden 'back: person' notiert, zaehlen nicht und setzen die Zaehlung zurueck; nicht auf round allein stuetzen. Danach Testlauf, Bughunt nur ueber den Diff, dann review - keine weitere Runde aus testing, was danach kommt ist note: fuer die Person. shipped.sha256-Fund bleibt note: pre-existing, Teamlead klaert mit Alex. Prompt-Zaehler klein halten, echte Zaehlung kommt mit DF0YD6/H6CD33.
- **2026-11-20 17:01 · Alexander Sacharov** — back: testing, round 6
