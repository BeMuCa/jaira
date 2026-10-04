---
id: 01M44AP5J9FPGBAVBCSQ747HER
title: Eine Rolle baut aus den Tickets in einer Menschen-Lane eine Abnahme-Seite und liest die Entscheidungen zurueck
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "jaira roles install liefert jaira-role-acceptance: die Rolle macht aus allen Tickets einer Menschen-Lane eine Abnahme-Seite (Artifact), auf der die Person nur noch prueft, was eine Maschine nicht pruefen kann, und schickt ihre Entscheidungen zurueck aufs Board - ein Befund geht in dasselbe Ticket."
context: |-
  Was heute schiefgeht:
  - Eine Abnahme heute: die Person geht jedes Ticket einzeln nach seinem review-check durch.
  - Beispiel requirementsgenie feat/R33C4B, 04.10.2026: 55 Tickets, 378 Abnahme-Schritte.
  - Davon 51 Schritte Testbefehle (task test:one, npx vitest), die ein Agent selbst laufen lassen kann.
  - 52 Schritte 'localhost:8109 oeffnen', 25 'anmelden', 24 'Projekt 93' - dieselbe Vorbereitung je Ticket neu.
  - 9 bis 11 Schritte starten je einen Lauf mit echtem Modell (kostet Tokens), obwohl ein Lauf pro Thema alle Tickets abdeckt.
  - Ein Befund aus der Abnahme wurde meist ein neues Ticket (TUI: nur 'a accept' und 'f follow-up', internal/tui/signoff.go:151). 25 der 56 Tickets dieser Woche sind so entstanden.

  Was schon existiert:
  - AlSa hat die Seite von Hand bauen lassen: https://claude.ai/artifact/NhU4pMRcgx1Tu4uAWemjQY (Bloecke nach Thema, ok/nein/Uebersprungen je Schritt, Annehmen/Zurueck je Ticket, Speicher im Artifact-db).
  - Entwurf der Rolle auf AlSas Rechner: ~/.claude/skills/jaira-role-acceptance/ (SKILL.md, template.html, build.py).
  - Rollen liegen in core/role/builtin/<id>/ und werden mit allen Dateien eingebettet (core/role/role.go:33, go:embed all:builtin).

  Entscheidungen:
  - Pro Block EIN Szenario statt Schritte pro Ticket; jeder Schritt nennt die Tickets, die er belegt (covers).
  - Maschinen-Schritte laufen vorher und stehen fertig auf der Seite.
  - Ein Zurueck wird DoD-Punkt desselben Tickets (jaira dod --add) und geht nach in-progress; ein neues Ticket nur, wenn der Befund ausserhalb des Ziels liegt.
  - Annehmen auf der Seite ist die Entscheidung der Person; die Rolle fuehrt sie mit --force aus und sagt das.
  - build.py nur mit der Python-Standardbibliothek, aufgerufen mit python3 - jaira darf keine weitere Laufzeit verlangen als die, die ein Agent ohnehin hat.

  Verwandt: PAAAJH (Rueckweg aus einer Menschen-Lane).
definition-of-done: "jaira roles install schreibt jaira-role-acceptance mit SKILL.md, template.html und build.py"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-04T20:46:32Z
updated-at: 2026-10-04T21:05:20Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-94145
claimed-at: 2026-10-04T20:47:12Z
outcome-what: "Wie Runde 3, dazu .gitattributes fuer die neuen .py- und .html-Dateien der Rollen."
outcome-why: "Ein Windows-Checkout haette sie mit CRLF eingebettet."
outcome-resolves: "Wie Runde 3."
review-summary: "none"
review-gaps: "build.py: sechsmal wiederholtes ((coll.get(k) or {}).get('note') or '').strip() in einen Helfer said() gezogen; state() gibt nur noch die Urteile zurueck statt (out, results) - plan/report lesen results selbst aus db. template.html: stepState() nutzt stepStale() statt die Rundenpruefung zu wiederholen. Ausgaben von page/plan/report auf dem Demo-Datensatz vorher/nachher identisch. Gelassen: REPORT-Labels in build.py und STR in template.html sind getrennte Woerterbuecher (Python vs. Seite, verschiedene Schluessel) - kein gemeinsamer Ort ohne Build-Schritt; state()/verdict() bewusst doppelt in Python und JS, weil Seite und Rueckleser dieselbe Regel brauchen; roleFiles-Tabelle in role_test.go bereits in Runde 1 entdoppelt; Dispatcher-/Teamlead-Texte ohne Ueberschuss."
test-verdict: "fail: go test ./internal/wintrap rot - build.py, template.html, report.html der neuen Rolle sind in .gitattributes nicht auf eol=lf gepinnt; DoD 1-13 sonst im Baum bestaetigt"
---

# Eine Rolle baut aus den Tickets in einer Menschen-Lane eine Abnahme-Seite und liest die Entscheidungen zurueck

## Definition of Done

- [x] jaira roles install schreibt jaira-role-acceptance mit SKILL.md, template.html und build.py
  proof: core/role/role_test.go TestAcceptanceShipsItsBuilder, TestInstallWritesEveryFile
- [x] build.py page verweigert eine Seite, auf der ein Ticket von keinem Schritt und keiner Maschinenpruefung belegt ist, oder ein Schritt ein unbekanntes Ticket nennt
  proof: core/role/builtin/jaira-role-acceptance/build.py:33 check(); von Hand: fehlende Abdeckung, unbekanntes Ticket und lang fr ergeben exit 1 mit drei Zeilen
- [x] build.py plan macht aus den gespeicherten Markierungen pro Ticket: zurueck (dod --add, note, move in-progress), angenommen, oder Rueckfrage an die Person
  proof: core/role/builtin/jaira-role-acceptance/build.py:83 plan(); von Hand mit Block 3 aus R33C4B: return, accept trotz Fehler, return ohne Grund, Blocknotiz
- [x] Eine Markierung aus einer frueheren Runde zaehlt auf der Seite als offen, sobald das Ticket eine neue Runde hat
  proof: core/role/builtin/jaira-role-acceptance/template.html:171 stepState/verdict vergleichen round mit der Runde des Tickets
- [x] Die Rolle laesst keinen Testbefehl fuer die Person stehen und fuehrt nichts aus, was schreibt, Tokens kostet oder den Stack der Person umstellt
  proof: core/role/builtin/jaira-role-acceptance/SKILL.md:83 Abschnitt 'Split every review-check step' mit 'Never run'
- [x] Ein Test in core/role prueft, dass die Rolle mit allen drei Dateien eingebettet ist
  proof: core/role/role_test.go TestAcceptanceShipsItsBuilder
- [x] core/release/NOTES.md nennt die neue Rolle
  proof: core/release/NOTES.md:17
- [x] Die Rolle nimmt beliebig viele Ticket-IDs auf einmal und fuehrt sie unter einem Namen (Meilenstein, sonst Zweig); ein zweiter Aufruf mit weiteren IDs ergaenzt dieselbe Seite statt eine neue zu bauen
  proof: core/role/builtin/jaira-role-acceptance/SKILL.md:30 Abschnitt 'Which tickets, under which name' und 'you are extending'
- [x] Dispatcher und Teamlead geben Tickets, die zu einer Aenderung gehoeren und in einer Menschen-Lane stehen, gemeinsam an jaira-role-acceptance weiter und nennen den Link in ihrem Bericht; aufgerufen von einem Agenten antwortet die Rolle in drei Zeilen
  proof: core/role/builtin/jaira-dispatcher/SKILL.md:417, core/role/builtin/jaira-teamlead/SKILL.md:33, core/role/builtin/jaira-role-acceptance/SKILL.md 'Who called you'
- [x] Auf der Seite laesst sich freier Text schreiben: ein Kommentar je Ticket (auch beim Annehmen) und eine Bemerkung je Block; build.py plan traegt beides ins Board
  proof: core/role/builtin/jaira-role-acceptance/template.html:135 Ticket-Kommentar immer sichtbar, Blocknotiz in notes/b<n>; build.py plan druckt beides
- [x] Die Sprache der Seite folgt der Sprache des Gespraechs mit der Person (de, en, ru; sonst en); Schritte bleiben in der Sprache des Tickets, und was die Person schreibt, geht woertlich in die Notiz und uebersetzt in den DoD-Punkt
  proof: core/role/builtin/jaira-role-acceptance/SKILL.md:60 Abschnitt Language; template.html STR ru/de/en
- [x] Ein Ticket gilt als angenommen, wenn alle seine Schritte ok und seine Maschinenpruefungen gruen sind; einen eigenen Annehmen-Knopf gibt es nicht, und jedes 'nein' klaert die Rolle selbst: zurueck ins Ticket, neuer Fehler mit --follows, falscher Schritt, oder Rueckfrage
  proof: core/role/builtin/jaira-role-acceptance/template.html verdict() und core/role/builtin/jaira-role-acceptance/build.py state(); core/role/builtin/jaira-role-acceptance/SKILL.md Abschnitt 'Read back' mit return/new-bug/step/asked
- [x] build.py report schreibt nach dem Zuruecklesen ein festes Abnahmeprotokoll (HTML und Markdown): je Block angenommen, nicht angenommen, offen, und je fehlgeschlagenem Schritt Erwartung, Gesehenes und Entscheidung
  proof: core/role/builtin/jaira-role-acceptance/build.py report(); von Hand mit Block 3: 4 angenommen, 3 nicht, 1 offen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-04 20:47 · Alexander Sacharov** — AlSa 04.10.: die Rolle muss weitergegeben werden koennen - traegt ein Dispatcher mehrere Tickets, weil der Fix gross ist, geht die Abnahme fuer alle zusammen auf eine Seite.
- **2026-10-04 20:53 · Alexander Sacharov** — in-progress: Entwurf aus ~/.claude/skills uebernommen, build.py auf python3 statt uv (Rollen duerfen nichts ausser python3 voraussetzen). Neu: Sprache (de/en/ru, Labels aus der Gespraechssprache), freier Text (Kommentar je Ticket, Notiz je Block), Name einer Abnahme (--name, Meilenstein, Zweig) mit Ergaenzen statt Neubau, Weitergabe durch Dispatcher und Teamlead. Probe: Block 3 'Zweiter Lauf' aus R33C4B - 48 alte Schritte werden 10 Personenschritte plus 4 Maschinenpruefungen (alle gruen, von mir gelaufen). TestLoadEmptyProjectDirIsSetUp in core/lane faellt auch auf upstream/master, wegen ~/.jaira/default-board.md auf diesem Rechner - nicht von hier.
- **2026-10-04 20:54 · Alexander Sacharov** — critique (Runde 1): (1) build.py:345 move nach in-progress ohne --force - gate.go:389 verweigert jeden Agenten-Move aus einer Menschen-Lane; ausserdem Lane-ID hart kodiert -> --force drucken, Rueckweg-Lane aus data.json. (2) build.py:334,344 deutsches Praefix 'Abnahme (Runde n)' fest in einer ausgelieferten Rolle -> englisch/neutral. (3) SKILL.md:188 Annehmen -> next_lane, TUI-Taste a -> Terminal-Lane (signoff.go:168) -> angleichen. (4) role_test.go: TestAcceptanceShipsItsBuilder dupliziert TestDispatcherShipsItsScript -> eine Tabelle Rolle->Pflichtdateien. Worktree-Diff (.jaira/milestones/demo-*) gehoert nicht zu diesem Ticket.
- **2026-10-04 20:59 · Alexander Sacharov** — in-progress Runde 2: Critique-Befunde umgesetzt - plan druckt --force und nimmt next_lane/return_lane aus data.json (build.py check verlangt beide), Praefix englisch ('acceptance (round n)'), Annehmen geht zu next_lane und SKILL.md sagt, warum nicht ins Terminal wie die Taste a (die Seite bestaetigt Verhalten, Review/Signoff laufen danach weiter), Tests als Tabelle roleFiles. Dazu AlSa 04.10.: kein eigener Annehmen/Zurueck-Knopf mehr - die Schritte entscheiden, und jedes 'nein' klaert der Agent selbst (Ticket zurueck, neuer Fehler, falscher Schritt, Rueckfrage); am Ende ein festes Protokoll (build.py report), weil die Seite sich weiter aendert. Beispiel: https://claude.ai/artifact/DHsAgCGbByU3A1KxWNGBKp und https://claude.ai/artifact/DLztPxZgKdtGbEd7HNYUJn
- **2026-10-04 20:59 · Alexander Sacharov** — critique (Runde 2): Befunde 1, 2, 4 aus Runde 1 behoben; 3 (next_lane statt Terminal) in SKILL.md begruendet - geschlossen. Neu durch die Umstellung auf 'Schritte entscheiden': ein Ticket, das nur von Maschinenpruefungen belegt ist (check() erlaubt das), wird nie angenommen - build.py:113 'own and all(...)', template.html:176 'steps.length && ...'. Entweder ohne Personenschritte nach gruenen Maschinenpruefungen annehmen (beide Stellen), oder check() verlangt mindestens einen Personenschritt je Ticket.
- **2026-10-04 21:00 · Alexander Sacharov** — in-progress Runde 3: Critique-Befund 2 umgesetzt - ein Ticket ohne Personenschritt gilt als angenommen, wenn alle seine Maschinenpruefungen gruen sind (build.py state, template.html verdict); eine nicht gelaufene Pruefung nimmt nichts an. AlSa 04.10.: 'да так мне нравится' zu Schritte-entscheiden + Protokoll.
- **2026-10-04 21:00 · Alexander Sacharov** — critique (Runde 3): Befund aus Runde 2 behoben (acf4f02) - build.py state() und template.html verdict() nehmen ein nur maschinell belegtes Ticket bei gruenen Pruefungen an, beide Stellen gleich. Nichts offen.
- **2026-10-04 21:04 · Alexander Sacharov** — testing: Gate rot. go test -count=1 ./internal/wintrap: --- FAIL: TestRepositoryIsClean / wintrap_test.go:83: core/role/role.go:33: windows trap 2: //go:embed pulls in 3 file(s) that no line in .gitattributes pins, starting at core/role/builtin/jaira-role-acceptance/build.py. Fix: in .gitattributes neben core/role/builtin/**/*.md und *.sh auch core/role/builtin/**/*.py text eol=lf und core/role/builtin/**/*.html text eol=lf eintragen (Kommentar dort ergaenzen). Sonst gruen: go test -race ./core/role ./core/release ok; tui TestArchiveFromTheBoard und core/lane TestLoadEmptyProjectDirIsSetUp fallen nur mit dem echten HOME (~/.jaira/default-board.md), mit leerem HOME gruen - nicht von hier. Funktion: jaira roles install --into <tmp> schreibt SKILL.md, build.py, template.html, report.html; build.py page/plan/report auf acc-demo laufen (8 Tickets, 10 Schritte, 4 Maschinenpruefungen; plan: 4 angenommen mit move --force, Fall b3-s02, Blocknotiz; report: 4/3/1); page auf bad.json exit 1 mit fr, unbekanntem Ticket und fehlender Abdeckung. Alle DoD-Proofs im Baum gefunden.
- **2026-10-04 21:05 · Alexander Sacharov** — in-progress Runde 4: Testing-Befund - .gitattributes pinnt jetzt core/role/builtin/**/*.py und *.html auf LF (internal/wintrap TestRepositoryIsClean).
