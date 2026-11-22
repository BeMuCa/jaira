---
id: 01M7X84R87N1Z4K96A39HS3K1E
title: "Zwei Wortlaut-Raender aus TW5FE2: Release-Note-Zeile teilen, testing-Folgerunde bei unbekanntem Commit"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die zwei kleinen Raender, die bei der Abnahme von TW5FE2 offen blieben, sind geschlossen."
context: |-
  Bei der Abnahme von TW5FE2 (20.11.2026, dritte Signoff-Runde) als kleine Raender stehen gelassen, statt einer vierten Runde:
  1. core/release/notes/TW5FE2.md Zeile 7 traegt zwei Aenderungen (review:-Prefix und testing-Folgerunde) in einer Zeile; der Release-Notes-Leser sieht sie als eine.
  2. core/lane/builtin/27-testing.md: die Folgerunde nimmt den Commit aus der letzten Zeile "bughunt: done at <sha>" fuer git diff --name-only; nach einem Rebase existiert er nicht mehr und der Befehl scheitert.
  TW5FE2 ist dann schon ausgeliefert: die Note geht als neue Datei core/release/notes/<dieses Ticket>.md, nicht in die TW5FE2-Datei, falls die bereits in einem Versionsordner liegt.
definition-of-done: "Die zwei Aenderungen stehen als zwei '- '-Zeilen in den Release-Notes (bzw. eine neue Note-Datei, falls TW5FE2 schon released ist)"
tags: []
blocked-by: []
related: []
follows: 01M7WWZRD7WNC5NT1PS5TW5FE2
commits: []
created-at: 2026-11-20T19:50:10Z
updated-at: 2026-11-22T17:09:22Z
updated-by: Alexander Sacharov
rejected-kind: wontfix
rejected-why: erledigt in TW5FE2
rejected-by: Alexander Sacharov
rejected-at: "2026-11-22T17:09:22Z"
---

# Zwei Wortlaut-Raender aus TW5FE2: Release-Note-Zeile teilen, testing-Folgerunde bei unbekanntem Commit

## Definition of Done

- [ ] Die zwei Aenderungen stehen als zwei '- '-Zeilen in den Release-Notes (bzw. eine neue Note-Datei, falls TW5FE2 schon released ist)
- [ ] 27-testing.md (builtin und .jaira/lanes byte-gleich) sagt: ist der Commit aus 'bughunt: done at' unbekannt (git cat-file -e scheitert), laeuft die ganze Suite wie in Runde 1; shipped.sha256 aktualisiert

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

