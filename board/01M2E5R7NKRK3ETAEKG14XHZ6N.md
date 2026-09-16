---
id: 01M2E5R7NKRK3ETAEKG14XHZ6N
title: "Lanes brauchen einen Kanal an den Menschen, nicht nur an die naechste Lane"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Jede Lane kann eine Zeile hinterlassen, die fuer einen Menschen bestimmt ist; die human-Lane bekommt diese Zeilen gesammelt und in der Reihenfolge, in der sie entstanden sind"
context: |-
  Heute hat eine Lane zwei Ausgaenge, und beide zielen auf Maschinen.
  - output-produces ist der Vertrag an die naechste Lane: plan, outcome-what, review-verdict. Wird von 'jaira show --for-lane' als Eingabe weitergereicht.
  - 'jaira note' ist das Journal fuer den naechsten Agenten: Sackgassen, warum dies und nicht das.
  Was fehlt, ist der Kanal an die Person. Beobachtet am 2026-09-13 beim ersten Lauf von Ticket C9V7ZV: die pre-process-Lane hatte einen Plan mit neun Schritten fertig, der Mensch am Bildschirm erfuhr davon nichts. Die Information blieb in der Lane liegen, bis am Ende jemand zusammenfasste.
  Wer heute versucht, das ueber 'jaira note' zu loesen, verschlechtert beides: das Journal fuellt sich mit Dingen fuer den Menschen, und der Mensch liest vierzig Zeilen Agenten-Notizen, um zwei zu finden, die ihn angehen.
  Vorschlag, der nichts Neues erfindet: ein reservierter Ausgabename, den jede Lane in output-produces fuehren darf - Arbeitstitel 'human-note'. Gesammelt in Reihenfolge, als Eingabe der human-Lane ausgeliefert, im TUI am Ticket sichtbar. Der Mechanismus ist der bestehende Lane-Vertrag, es kommt kein zweiter daneben.
  Offene Frage fuer die Brainstorm-Lane, nicht vorweggenommen: ob das ein reservierter Ausgabename ist oder ein eigenes Frontmatter-Feld, und ob eine Lane hoechstens eine solche Zeile hinterlassen darf. Eine Zeile pro Lane haelt es lesbar; mehrere machen daraus ein zweites Journal.
  Nicht Teil dieses Tickets: das Berichten selbst. Dass ein Dispatcher pro Lane nach oben meldet, steht schon im dispatcher-Skill und braucht keinen Code.
definition-of-done: "Eine Lane kann eine an den Menschen gerichtete Zeile hinterlassen; 'jaira show <id> --for-lane human --json' liefert diese Zeilen gesammelt und in Entstehungsreihenfolge; das TUI zeigt sie am Ticket; 'jaira note' bleibt unveraendert das Agenten-Journal; eine Zeile in core/release/NOTES.md unter ## Unreleased; go test ./... -race gruen"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-13T20:01:18Z
updated-at: 2026-09-13T20:01:18Z
---

# Lanes brauchen einen Kanal an den Menschen, nicht nur an die naechste Lane

## Definition of Done

- [ ] Eine Lane kann eine an den Menschen gerichtete Zeile hinterlassen; 'jaira show <id> --for-lane human --json' liefert diese Zeilen gesammelt und in Entstehungsreihenfolge; das TUI zeigt sie am Ticket; 'jaira note' bleibt unveraendert das Agenten-Journal; eine Zeile in core/release/NOTES.md unter ## Unreleased; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

