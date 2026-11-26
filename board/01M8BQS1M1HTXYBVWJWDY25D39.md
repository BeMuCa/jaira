---
id: 01M8BQS1M1HTXYBVWJWDY25D39
title: Ein uebernommenes Board traegt seinen Ersteller und liegt getrennt von den eigenen
status: backlog
ready: true
creator: Berk
assignee: Berk
goal: "Im persoenlichen Katalog sieht man, welche Boards (und Lanes) man selbst gebaut und welche man von wem uebernommen hat - in der TUI mit Ersteller, und auf der Platte getrennt."
context: |-
  Wunsch von Berk, 26.11.2026, beim Ansehen von 4G07X9 (Boards wie Lanes, ~/.jaira/boards, publish/shared/adopt; Branch origin/feat/4G07X9-boards-catalogue, signoff).
  Heute: eine Lane traegt creator: (core/lane/lane.go:168), gestempelt beim Veroeffentlichen und Uebernehmen; mitgelieferte Lanes gelten als 'jaira'. Ein Board auf dem 4G07X9-Branch hat kein solches Feld (git grep creator in core/board und internal/cli/boards*: kein Treffer).
  Berks Idee: unter ~/.jaira/boards nach eigenen und uebernommenen trennen, uebernommene in einem Ordner je Ersteller. Nur Ordnerstruktur - die TUI zeigt weiter alles, liest aber den Ersteller mit an.
  Zu klaeren: Ersteller nur im Frontmatter oder auch als Ordner (Ordner und Feld koennen auseinanderlaufen); was ein Umbenennen des Erstellers bedeutet; Zusammenhang mit MMY2CS (Herkunft merken).
definition-of-done: "Ein Board traegt creator: wie eine Lane; publish stempelt ihn, adopt behaelt ihn"
tags: []
blocked-by:
  - 01M8BQFBMS778FE9WK6YDAW7WP
related: []
follows: 01M5GFAV1C3DV9NWZFXC4G07X9
commits: []
created-at: 2026-11-26T10:52:46Z
updated-at: 2026-11-26T10:52:59Z
updated-by: Berk
---

# Ein uebernommenes Board traegt seinen Ersteller und liegt getrennt von den eigenen

## Definition of Done

- [ ] Ein Board traegt creator: wie eine Lane; publish stempelt ihn, adopt behaelt ihn
- [ ] Eigene und uebernommene Boards liegen in ~/.jaira/boards getrennt; uebernommene je Ersteller
- [ ] Die TUI zeigt den Ersteller eines Boards und einer Lane im Katalog
- [ ] Release-Notiz in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

