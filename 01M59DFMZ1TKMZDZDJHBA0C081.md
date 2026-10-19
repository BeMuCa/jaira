---
id: 01M59DFMZ1TKMZDZDJHBA0C081
title: "Ein Ticket laesst sich mit Grund ablehnen: jaira reject, und das Archiv sagt warum"
status: todo
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer ein Ticket nicht macht, sagt mit einem Befehl warum - und Monate spaeter steht der Grund am Ticket, getrennt von bloss weggeraeumten."
context: |-
  Was heute fehlt: ein abgelehntes Ticket sieht genauso aus wie ein weggeraeumtes.
  - 'jaira archive <id>' (core/ticket/store.go:301, Taste x auf dem Board) verschiebt nach .jaira/archive/ und schreibt KEINEN Grund.
  - 'jaira delete <id>' loescht endgueltig - nur fuer versehentlich angelegte Tickets, keine Ablehnung.
  - Folge: 'jaira archive' (Liste) und das Web-Portal koennen nicht sagen, ob ein Ticket ein Duplikat, verworfen oder ueberholt war. Das verletzt den Kernwert 'was war das Ticket und wo steht es'.
  Anlass: Alex am 19.10.2026 - 'wir haben keine richtige Moeglichkeit, ein Ticket zu loeschen oder abzulehnen'.
  Entschieden von Alex (19.10.2026):
  - Neuer Befehl 'jaira reject <id> --why <text> [--duplicate-of <id>]'. Art: wontfix | duplicate | obsolete (--duplicate-of setzt duplicate).
  - Schreibt ins Frontmatter einen Block 'rejected:' mit kind, why, by, at; danach archiviert er wie 'jaira archive' (gleicher Code, kein zweiter Weg).
  - Nur ein Mensch lehnt ab. Ein Agent schiebt das Ticket stattdessen nach human mit --question 'ablehnen? Grund: ...'. Wie ein Agent erkannt wird: denselben Mechanismus nehmen, der heute den Ausgang aus human sperrt (requires-human-exit, siehe VW1WB1) - nichts Neues erfinden.
  Verworfen: Grund nur an archive haengen (Aufraeumen und Ablehnen blieben ununterscheidbar); eigene Lane 'rejected' (aendert die Lanes jedes Boards).
  Offen fuer die Plan-Lane: ob x auf dem Board kuenftig nach dem Grund fragt oder eine eigene Taste ablehnt (x bleibt archivieren, siehe model.go:1532-Kommentar).
definition-of-done: "'jaira reject <id> --why <text>' schreibt rejected.kind/why/by/at ins Frontmatter und verschiebt das Ticket nach .jaira/archive/; ohne --why Exit 2"
tags:
  - cli
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-19T06:27:18Z
updated-at: 2026-10-19T06:27:44Z
updated-by: Alexander Sacharov
---

# Ein Ticket laesst sich mit Grund ablehnen: jaira reject, und das Archiv sagt warum

## Definition of Done

- [ ] 'jaira reject <id> --why <text>' schreibt rejected.kind/why/by/at ins Frontmatter und verschiebt das Ticket nach .jaira/archive/; ohne --why Exit 2
- [ ] --duplicate-of <id> setzt kind=duplicate und verlinkt das Original; eine unbekannte id wird mit Exit 5 abgelehnt
- [ ] Ein Agent wird mit Exit 3 abgewiesen und die Meldung nennt den Weg ueber human --question
- [ ] Auf dem Board laesst sich ein Ticket mit einer einzeiligen Begruendung ablehnen
- [ ] 'jaira archive' (Liste), 'jaira show' und das Web-Portal zeigen kind und why; 'jaira restore' nimmt den rejected-Block wieder ab
- [ ] Release-Note core/release/notes/<ID>.md, Agent-Block/jaira-Skill nennen reject; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

