---
id: 01M5ZJ4NNCYN6YBWVAPWAHG254
title: Die Felder einer Schleifen-Lane ueberschreiben nicht die Runde davor
status: backlog
ready: false
creator: BeMuCa
goal: "review-summary und review-gaps, die critique, optimize und review schreiben, bleiben je Runde und je Lane erhalten statt sich gegenseitig zu ersetzen; jaira show zeigt die Historie."
context: |-
  Was schief ist: die Review-Lane ueberschrieb beim Move den review-summary-Text, den critique geschrieben hatte - zum zweiten Mal in zwei Teamlead-Durchlaeufen (27.10.2026, Feedback L314).
  Ursache: critique und review liefern beide review-summary, optimize und review beide review-gaps (Lane-Dateien dieses Boards). jaira move schreibt das Feld und ersetzt den alten Wert.
  jaira set --append sagt in seiner Hilfe, dass Review-Felder damit ihre Historie ueber Schleifenrunden behalten - move benutzt diesen Weg nicht.
  Folge: die Begruendung, warum critique zurueckgeschickt hat, ist nach dem Review weg; der Mensch in signoff liest nur die letzte Stimme.
  Zusammenhang: NFJCTK (Schema fuer Lane-Felder), FCMP17 (Lane-Artefakte an die naechste Lane).
definition-of-done: "Ein Move aus critique mit --review-summary und danach ein Move aus review mit --review-summary hinterlassen beide Texte im Ticket, mit Lane und Runde erkennbar"
tags:
  - gates
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:51:59Z
updated-at: 2026-10-27T20:51:59Z
---

# Die Felder einer Schleifen-Lane ueberschreiben nicht die Runde davor

## Definition of Done

- [ ] Ein Move aus critique mit --review-summary und danach ein Move aus review mit --review-summary hinterlassen beide Texte im Ticket, mit Lane und Runde erkennbar
- [ ] jaira show <id> zeigt beide Eintraege; --json liefert sie als Liste
- [ ] Ein Test in core/ belegt, dass der zweite Schreiber den ersten nicht ersetzt
- [ ] Release-Notiz in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

