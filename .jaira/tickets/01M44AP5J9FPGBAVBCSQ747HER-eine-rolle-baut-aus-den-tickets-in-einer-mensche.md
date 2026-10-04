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
updated-at: 2026-10-04T20:53:14Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-94145
claimed-at: 2026-10-04T20:47:12Z
outcome-what: "Neue eingebaute Rolle jaira-role-acceptance (SKILL.md, template.html, build.py) baut aus Tickets in Menschen-Lanes eine Abnahme-Seite mit db-Speicher und liest die Entscheidungen als jaira-Befehle zurueck; Dispatcher und Teamlead geben Tickets einer Aenderung gemeinsam weiter; README und NOTES."
outcome-why: "Die Person soll nur pruefen, was eine Maschine nicht pruefen kann; ein Befund aus der Abnahme soll ins selbe Ticket statt ein neues Ticket zu werden."
outcome-resolves: "Abnahme pro Ticket mit wiederholter Vorbereitung und Testbefehlen fuer die Person (R33C4B: 378 Schritte fuer 55 Tickets); Befunde wurden Folgetickets."
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

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-04 20:47 · Alexander Sacharov** — AlSa 04.10.: die Rolle muss weitergegeben werden koennen - traegt ein Dispatcher mehrere Tickets, weil der Fix gross ist, geht die Abnahme fuer alle zusammen auf eine Seite.
- **2026-10-04 20:53 · Alexander Sacharov** — in-progress: Entwurf aus ~/.claude/skills uebernommen, build.py auf python3 statt uv (Rollen duerfen nichts ausser python3 voraussetzen). Neu: Sprache (de/en/ru, Labels aus der Gespraechssprache), freier Text (Kommentar je Ticket, Notiz je Block), Name einer Abnahme (--name, Meilenstein, Zweig) mit Ergaenzen statt Neubau, Weitergabe durch Dispatcher und Teamlead. Probe: Block 3 'Zweiter Lauf' aus R33C4B - 48 alte Schritte werden 10 Personenschritte plus 4 Maschinenpruefungen (alle gruen, von mir gelaufen). TestLoadEmptyProjectDirIsSetUp in core/lane faellt auch auf upstream/master, wegen ~/.jaira/default-board.md auf diesem Rechner - nicht von hier.
