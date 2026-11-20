---
id: 01M7WWYG213BWB06A7D1CRASZS
title: "Ein offener PR wird begleitet: jaira sieht neue Kommentare, CI-Stand und wartende Re-Requests und meldet sie"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Pro Ticket mit PR zeigt jaira den Stand (Kommentare offen/beantwortet, CI, Review angefordert) und meldet einer Session, wenn sich daran etwas aendert, damit die naechste Runde startet."
context: |-
  Heute: den Stand eines PR (z.B. "3 offen, 2 beantwortet, Re-Request ausstehend") haelt der Teamlead im Kopf; jaira weiss von PRs nichts ausser in jaira-role-pr.
  Alex 20.11.2026: Automatik, um einen PR zu begleiten, und Moeglichkeiten, ihn zu verfolgen.
  Kein Daemon (Grundsatz von jaira): Abfrage beim Board-Refresh, im Cockpit oder per Befehl; gh/glab als Quelle.
  Verwandt: QFWRXH (PR-Spalte auf dem Board), das Ticket mit der optionalen PR-Runde.
definition-of-done: "Ein Befehl oder eine Cockpit-Ansicht zeigt je Ticket mit PR: offene/beantwortete Kommentare, CI-Stand, ob ein Review aussteht"
tags: []
blocked-by:
  - 01M5ZJ5DKQZ83GYBGPJHAX8K02
related:
  - 01M7WWYB5XZMD01DZY3MJG3REW
  - 01M6HZQPZZFAF3MXY3AWQFWRXH
  - 01M7WX040JC9Z202EH852EF48S
commits: []
created-at: 2026-11-20T16:34:31Z
updated-at: 2026-11-20T17:58:18Z
updated-by: Alexander Sacharov
---

# Ein offener PR wird begleitet: jaira sieht neue Kommentare, CI-Stand und wartende Re-Requests und meldet sie

## Definition of Done

- [ ] Ein Befehl oder eine Cockpit-Ansicht zeigt je Ticket mit PR: offene/beantwortete Kommentare, CI-Stand, ob ein Review aussteht
- [ ] Kommt ein neuer Kommentar oder wird CI rot, erfaehrt es die Session des Teamleads ohne Nachfragen

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
