---
id: 01M7X20YACKGBHD6P1QKMFYM1A
title: "In signoff hakt die Person die Notes an, die behoben werden sollen, und Give back nimmt genau diese mit"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die Explain-/Ticket-Ansicht im Cockpit zeigt jede Note aus review-gaps mit einem Haken; die Person waehlt aus, was behoben werden soll, und Give back schreibt genau diese Notes als Auftrag an den Worker; nichts angehakt plus Accept nimmt an."
context: |-
  Heute: Give back im Cockpit bietet nur die Optionen aus der explain:-Note (1) [give back] 2) [accept]); welche der 3-5 Notes behoben werden sollen, muss die Person in Worten sagen oder der Teamlead herausfinden.
  Alex 20.11.2026: "muss man die Moeglichkeit geben, Haken zu setzen, was wir korrigieren wollen?" - nach einer Signoff-Uebersicht als Tabelle mit Marken je Note (Blocker / kleiner Rand / bekannte Grenze / nicht aus diesem Ticket).
  Alex gibt heute fast jedes Ticket zurueck, um offene Raender zu schliessen; die Auswahl ist also der Normalfall, kein Sonderfall.
  Voraussetzung: DF0YD6 (Notes werden eingeordnet und farbig gezeigt).
  Verwandt: Q63EEV (Tickets im Cockpit einzeln abarbeiten: Pane mit Accept/Give back).
definition-of-done: "Jede Note aus review-gaps steht mit Haken in der Ansicht; Blocker und kleine Raender sind vorausgewaehlt, bekannte Grenzen und fremde Notes nicht (Test)"
tags: []
blocked-by:
  - DF0YD6
related:
  - 01M7WWZWSDHK37DGGHJFDF0YD6
  - 01M5M7JTVR8G7DZ1MGHJQ63EEV
commits: []
created-at: 2026-11-20T18:03:14Z
updated-at: 2026-11-20T18:03:24Z
updated-by: Alexander Sacharov
---

# In signoff hakt die Person die Notes an, die behoben werden sollen, und Give back nimmt genau diese mit

## Definition of Done

- [ ] Jede Note aus review-gaps steht mit Haken in der Ansicht; Blocker und kleine Raender sind vorausgewaehlt, bekannte Grenzen und fremde Notes nicht (Test)
- [ ] Give back mit Auswahl schreibt eine Note 'Given back from signoff: beheben: <Notes woertlich>' und der Worker arbeitet genau diese ab; nicht angehakte Notes bleiben als akzeptiert vermerkt (Test)
- [ ] Eine nicht angehakte Note, die kein Blocker ist, haelt Accept nicht auf; ein nicht angehakter Blocker fragt einmal nach (Test)
- [ ] Eingabeklassen: keine Notes, nur ein Blocker, sehr viele Notes (Scrollen), Note ueber mehrere Zeilen, review-gaps aus frueherer Runde, zwei Personen/Sessions entscheiden gleichzeitig

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

