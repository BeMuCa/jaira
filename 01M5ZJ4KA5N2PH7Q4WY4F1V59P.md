---
id: 01M5ZJ4KA5N2PH7Q4WY4F1V59P
title: "Ein Review-Verdikt 'zurueckschicken' laesst das Ticket nicht nach signoff"
status: backlog
ready: true
creator: BeMuCa
goal: "Steht im review-verdict der Review-Lane 'send back' (oder das Gegenstueck des Boards), nimmt jaira move aus review nur noch in-progress als Ziel; signoff wird mit Exit 3 abgelehnt und nennt das Verdikt."
context: |-
  Was schief ist: die Review-Lane schrieb als Verdikt 'zurueckschicken', das Ticket landete trotzdem in signoff (Ticket 25DZW3 auf einem anderen Board, Teamlead-Session 25.-27.10.2026).
  Heute prueft kein Gate, ob Verdikt und Ziellane zusammenpassen: review muss review-verdict liefern, aber move liest den Inhalt nicht.
  Der Mensch sah ein Ticket zur Abnahme, das der Reviewer selbst nicht freigegeben hatte.
  Bekannt: die Gates leben in core/gate/gate.go; welche Werte review-verdict annimmt, steht in der Lane-Datei review der Boards (grep review-verdict core/lane/).
  Nicht gewollt: ein Freitext-Parser. Ein fester kleiner Wertevorrat (z. B. accept, send-back) reicht; alles andere ist kein Verdikt.
  Zusammenhang: L312/L327 der Feedback-Liste vom 27.10.2026.
definition-of-done: "Ein Move aus review nach signoff mit review-verdict 'send back' wird mit Exit 3 abgelehnt; die Meldung nennt das Verdikt und die erlaubte Ziellane"
tags:
  - gates
blocked-by:
  - 01M6M7EA5EG5TWJSTVNQH6CD33
related:
  - 01M7WWZWSDHK37DGGHJFDF0YD6
commits: []
created-at: 2026-10-27T20:51:56Z
updated-at: 2026-11-20T16:49:00Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
---

# Ein Review-Verdikt 'zurueckschicken' laesst das Ticket nicht nach signoff

## Definition of Done

- [ ] Ein Move aus review nach signoff mit review-verdict 'send back' wird mit Exit 3 abgelehnt; die Meldung nennt das Verdikt und die erlaubte Ziellane
- [ ] Ein Move aus review nach in-progress mit 'send back' geht ohne --force
- [ ] Die Werte, die review-verdict annimmt, stehen in der Lane-Datei und in jaira show --for-lane review
- [ ] Release-Notiz in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-20 16:49 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
