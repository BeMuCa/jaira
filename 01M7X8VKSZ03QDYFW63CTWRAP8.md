---
id: 01M7X8VKSZ03QDYFW63CTWRAP8
title: "Was der Teamlead am 20.11. gelernt hat, steht in den ausgelieferten Rollen statt nur in einem lokalen Gedaechtnis"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Jeder Teamlead, der jaira installiert, bekommt die Arbeitsweise, die sich am 20.11.2026 bewaehrt hat: Stand vor jeder Liste neu lesen, Signoff-Bericht als Statuszeile plus eingeordnete Tabelle, eine Signoff-Rolle mit Rueckgabe-Regeln, Schlussbericht der Worker mit Feldern."
context: |-
  Heute stehen diese Regeln nur im lokalen Gedaechtnis einer Teamlead-Session auf einem Rechner; jeder andere Teamlead (Kollegen, andere Projekte) macht dieselben Fehler.
  Was an einem Tag schiefging und wie es geloest wurde:
  1. Der Teamlead listete laengst entschiedene Tickets als offen - Entscheidungen im Cockpit/Browser erreichen ihn nicht. Regel: vor jeder Liste "wartet auf dich" jaira show je Ticket und gh pr list lesen, die Liste nur bei Aenderung zeigen. (Teil im Binary: QTQW15 DoD 5.)
  2. Die Person schickte fast jedes Ticket zurueck, weil der Explain-Text jede Note wie einen Fehler klingen laesst. Die Form, die sofort verstanden wurde: Statuszeile (DoD, Tests, Blocker, bughunt auf HEAD) + Tabelle je Note mit Marke (rot Blocker, gelb kleiner Rand, weiss bekannte Grenze, blau nicht aus diesem Ticket) + eine Empfehlung + beide move-Befehle.
  3. Rueckgabe-Regeln, die Endlosrunden beenden: nach 2 Rueckgaben aus signoff gehen gelbe Raender in ein Folge-Ticket; Wortlaut allein ist nie ein Rueckgabegrund; wiederholen sich Funde einer Klasse, wird die Ursache benannt. Ausnahme: Rechte-/Sicherheitsgrenzen bleiben Blocker.
  4. Ein per Cockpit neu gestarteter Worker kennt seinen Teamlead nicht (QTQW15) - bis dahin schickt der Teamlead ihm eine Zeile.
  5. PRs scheiterten an fehlenden Zeilen in core/shipped/shipped.sha256 - vor jedem Push go test ./core/role -run Shipped (nur dieses Repo; allgemein: die Projekt-Gates vor dem Push).
  Verwandt: DF0YD6 (Explain-Form im Cockpit), MFYM1A (Haken je Note), 55CS2M (Schlussbericht mit festen Feldern), TW5FE2 (Lane-Regeln, done).
definition-of-done: "jaira-teamlead/SKILL.md: Abschnitt 'Report to the human' verlangt die Statuszeile + eingeordnete Tabelle + eine Empfehlung fuer signoff/human und das Neu-Lesen des Stands vor jeder Liste offener Entscheidungen"
tags: []
blocked-by: []
related:
  - 01M7WWZWSDHK37DGGHJFDF0YD6
  - 01M7X20YACKGBHD6P1QKMFYM1A
  - 01M7WX0C09TC63TKPXT555CS2M
  - 01M7X0CFXEAEXDHVWTD3QTQW15
commits: []
created-at: 2026-11-20T20:02:40Z
updated-at: 2026-11-20T21:43:37Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-16592
claimed-at: 2026-11-20T20:04:03Z
outcome-what: "Neue Rolle jaira-role-signoff (Statuszeile, markierte Notes-Tabelle, eine Empfehlung, Rueckgabe-Regeln); Teamlead liest den Stand vor jeder Liste neu und berichtet signoff in dieser Form; worker, dispatcher und jaira-role-pr lassen vor jedem Push die CI-Gates des Projekts lokal laufen."
outcome-why: "Die Arbeitsweise vom 20.11.2026 stand nur im lokalen Gedaechtnis einer Teamlead-Session; jeder andere Teamlead haette dieselben Fehler gemacht (veraltete Listen, Rueckgaben wegen Notes ohne Einordnung, rote PRs durch fehlende shipped.sha256-Zeilen)."
outcome-resolves: Die Regeln reisen mit den ausgelieferten Rollen; jaira update --roles bringt sie zu jedem Nutzer.
review-summary: none
round: 2
---

# Was der Teamlead am 20.11. gelernt hat, steht in den ausgelieferten Rollen statt nur in einem lokalen Gedaechtnis

## Definition of Done

- [x] jaira-teamlead/SKILL.md: Abschnitt 'Report to the human' verlangt die Statuszeile + eingeordnete Tabelle + eine Empfehlung fuer signoff/human und das Neu-Lesen des Stands vor jeder Liste offener Entscheidungen
  proof: core/role/builtin/jaira-teamlead/SKILL.md ## Report to the human
- [x] Neue Rolle jaira-role-signoff <id>...: liest Tickets in signoff, gibt die Tabelle mit Einordnung und Empfehlung nach den Rueckgabe-Regeln aus, fertige move-Befehle und ein Folge-Ticket fuer gelbe Raender nach der zweiten Rueckgabe; aendert nichts selbst
  proof: core/role/builtin/jaira-role-signoff/SKILL.md (## The return rules, ## Per ticket)
- [x] jaira-worker/jaira-dispatcher: vor jedem Push die Gates des Projekts lokal laufen lassen, die CI sonst rot machen (im Skill allgemein formuliert)
  proof: jaira-worker/SKILL.md ### Before you push; jaira-dispatcher/SKILL.md ## Before anything is pushed; jaira-role-pr/SKILL.md ## Push the branch
- [x] Release-Note und shipped.sha256; die Rollen-Tests (core/role) gruen
  proof: core/release/notes/TWRAP8.md; core/shipped/shipped.sha256; go test ./core/role ok

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-20 20:02 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-20 20:04 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/35112.sock
- **2026-11-20 20:04 · Alexander Sacharov** — decision (worker, keine offene Frage): (1) neue Rolle heisst jaira-role-signoff, nur lesend, Ausgabe in der Sprache der Person; Marken wie im Kontext (rot/gelb/weiss/blau als Emoji, weil Terminal-Markdown keine Farben hat). (2) Push-Gate steht allgemein in jaira-worker, jaira-dispatcher UND jaira-role-pr (die dritte Stelle, die pusht - Geschwisterstelle nach Interim-Regel 2). Beispiel fuer dieses Repo nur als Beispiel. (3) README-Rollentabelle und docs/how-to-work-with-jaira.html bekommen die neue Rolle, wantRoles in role_test.go auch. (4) 'Schlussbericht der Worker mit Feldern' aus dem Ziel ist 55CS2M und steht nicht in der DoD - nicht hier.
- **2026-11-20 20:10 · Alexander Sacharov** — Stand bei Usage-Limit: DoD 1-3 geschrieben, Release-Note + shipped.sha256 angehaengt, core/role gruen, committed. Offen: DoD 4 abhaken, Lanes critique/optimize/testing/review. Hashes in shipped.sha256 muessen nach jeder weiteren SKILL.md-Aenderung neu angehaengt werden.
- **2026-11-20 21:41 · Alexander Sacharov** — critique: (1) Push-Gate-Text im Dispatcher kuerzen, er pusht nicht selbst. (2) Rueckgabe-Regeln im Teamlead durch Verweis auf jaira-role-signoff ersetzen, damit es nur eine Quelle gibt. Stehen gelassen: Statuszeile+Tabelle+Empfehlung bleiben im Teamlead ausgeschrieben, weil DoD 1 genau das im Teamlead-Abschnitt verlangt.
- **2026-11-20 21:43 · Alexander Sacharov** — critique runde 2: beide Funde der Runde 1 umgesetzt (c02bb7ea). Stehen gelassen: 'jaira move <id> --to done --force' in jaira-role-signoff nennt die Ziel-Lane fest - signoff->done ist die ausgelieferte Lane-Folge, und dieselbe Annahme steht schon in teamlead/acceptance.
