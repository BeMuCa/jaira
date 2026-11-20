---
id: 01M6M7EA5EG5TWJSTVNQH6CD33
title: A lane can say how many times it may send work back before a person decides
status: backlog
ready: true
creator: berkc
goal: "A loop lane declares max-rounds: <n>, and once a ticket has been sent back from it n times the next send-back goes to a person instead of to in-progress."
context: |-
  What is wrong today: the limit of three send-backs per loop lane (critique, optimize, testing) lives only in the prose of the worker role (core/role/builtin/jaira-worker/SKILL.md, 'the same lane sent work back three times'). The binary does not know it.
  So an agent that does not read that skill loops forever, and a board cannot choose a different limit for one lane.
  Triggered by: ticket ENCFVZ (skill/script/receipt lanes), decision of the maintainer on 04.11.2026 listing max-rounds as a lane attribute; kept out of ENCFVZ's scope on purpose.
  Already there: the ticket field 'round' counts send-backs (core/move/move.go, sendsBack); review-round stamps the round a verdict was written in.
  Open: whether the gate refuses the n+1-th send-back (exit 3, name the human lane) or only 'jaira show --for-lane' reports rounds left; and whether round counts per lane or per ticket (today it is per ticket).
definition-of-done: "A lane file may declare max-rounds: <n>; the loader refuses a non-positive or non-numeric value"
tags:
  - gates
blocked-by:
  - 01M7WWZWSDHK37DGGHJFDF0YD6
related:
  - 01M7WWZWSDHK37DGGHJFDF0YD6
follows: 01M67K7GEZF694985EHTENCFVZ
commits: []
created-at: 2026-11-04T21:29:03Z
updated-at: 2026-11-20T16:48:51Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
---

# A lane can say how many times it may send work back before a person decides

## Definition of Done

- [ ] A lane file may declare max-rounds: <n>; the loader refuses a non-positive or non-numeric value
- [ ] Sending work back from that lane for the (n+1)-th time is refused with a message naming the human lane, and --for-lane shows the rounds left
- [ ] README lane format and a release note describe it

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-20 16:48 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
