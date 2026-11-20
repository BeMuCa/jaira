---
id: 01M7WX0C09TC63TKPXT555CS2M
title: "Ein Worker meldet signoff in einer Zeile mit festen Feldern: Branch, Commits, Guard, offene Notes"
status: backlog
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
updated-at: 2026-11-20T16:49:56Z
updated-by: Alexander Sacharov
---

# Ein Worker meldet signoff in einer Zeile mit festen Feldern: Branch, Commits, Guard, offene Notes

## Definition of Done

- [ ] jaira-worker und jaira-dispatcher schreiben die Schlussmeldung im festen Format, dokumentiert im Skill
- [ ] Das Cockpit liest die Felder aus dieser Zeile

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
