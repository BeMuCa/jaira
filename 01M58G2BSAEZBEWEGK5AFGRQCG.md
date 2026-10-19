---
id: 01M58G2BSAEZBEWEGK5AFGRQCG
title: "Eine kaputte Lane-Datei verschwindet still aus den Settings, und das Lane-Formular fragt Felder ab, die niemand versteht"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Wer eine Lane anlegt, sieht sie danach auf dem Board oder bekommt gesagt, was an der Datei falsch ist; das Formular fragt nur, was man wirklich entscheiden muss."
context: |-
  Heute: berkc hat im Lanes-Screen zwei Lanes angelegt. .jaira/lanes/new-lane.md ist die unveraenderte Vorlage (id my-lane) - die sieht er als 'default lane'.
  .jaira/lanes/new-lane-2.md ist kaputt: die description ging ueber zwei Zeilen ('The Agent that works' / '.') und die Zeile 'after: human' fehlt. jaira list warnt nur: 'missing required field "id"', der TUI zeigt nichts - die Lane ist einfach weg.
  Das Skelett kommt aus writeLaneSkeleton (internal/tui/lanes.go:578) und hat 20+ Felder mit Kommentaren.
  Unklar fuer berkc: precedence (nur Merge-Rang, core/merge/merge.go:377, nicht Anzeigereihenfolge). Siehe auch die Diskussion, precedence ganz abzuschaffen.
  Pro Feld pruefen: braucht man es beim Anlegen? Kandidaten zum Weglassen: precedence, holds, logbook-on-entry, requires-outcome, requires-nonmodel-signal, creator.
  Nebenbei: die Vorlage unveraendert speichern legt eine Lane 'My Lane' an - nicht gewollt.
definition-of-done: "Eine Lane-Datei, die nicht geparst werden kann, zeigt im Lanes-Screen eine Fehlerzeile mit Dateiname und Grund statt zu verschwinden; Test"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-18T21:53:14Z
updated-at: 2026-10-19T08:34:30Z
updated-by: berkc
---

# Eine kaputte Lane-Datei verschwindet still aus den Settings, und das Lane-Formular fragt Felder ab, die niemand versteht

## Definition of Done

- [ ] Eine Lane-Datei, die nicht geparst werden kann, zeigt im Lanes-Screen eine Fehlerzeile mit Dateiname und Grund statt zu verschwinden; Test
- [ ] Die unveraendert gespeicherte Vorlage legt keine Lane an (oder fragt nach); Test
- [ ] Liste in der Notiz: jedes Feld der Vorlage mit behalten/weglassen/Default und Grund; die Vorlage folgt der Liste
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

