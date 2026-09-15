---
id: 01M2GPAV3W200QAWWGNQ1K9KZS
title: "Ein Board, das es schon gibt, bekommt eine geaenderte Lane nie zu sehen"
status: todo
ready: true
creator: Alexander Sacharov
goal: "Eine Korrektur an einer ausgelieferten Lane erreicht auch die Boards, die es schon gibt - ohne dass jemand auf jedem Rechner eine Zeile von Hand loescht."
context: |-
  74VM40 hat 'logbook-on-entry' aus den mitgelieferten Lanes entfernt, weil ein Move nach done 49 fremde fertige Tickets ins Logbuch fegte (Issue #6). Auf jedem Board, das vor dieser Aenderung entstanden ist, passiert das weiterhin.

  Die Ursache steht in core/lane/lane.go:479. Load() liest bei ProjectLanesActive NUR das Lane-Verzeichnis des Projekts und legt die eingebauten Lanes nie darunter. Eine Lane-Datei, die einmal geschrieben wurde, ist damit fuer immer die Wahrheit - auch wenn das Binary sie laengst korrigiert hat.

  Nachgestellt, zweimal auf diesem Board: am 11.09. und wieder am 14.09. hat ein Move nach done fremde fertige Tickets mitgenommen. Beide Male hat ein Mensch die Zeile von Hand aus .jaira/lanes/done.md geloescht - und weil .jaira/lanes/ in gitignore steht, pro Checkout einzeln. Am 14.09. waren es zwei Verzeichnisse.

  Dazu kommt, dass die Release-Notiz das Gegenteil behauptet: sie sagt glatt 'Finishing a ticket no longer files anything' und erklaert danach nur, wie man das alte Verhalten zurueckholt. Wer sie liest und ein aelteres Board hat, glaubt etwas Falsches ueber sein eigenes Werkzeug.

  Der Befund stammt aus der review-Lane von 74VM40 am 2026-09-14, die ihn ausdruecklich als Grund genannt hat, 74VM40 NICHT abzunehmen.

  Die Spannung, die die Plan-Lane aufloesen muss: 'ein Board ist sein Lane-Verzeichnis' ist eine bewusste Entscheidung (Commit 743737f, 27.08.) - eine vom Nutzer geschriebene Lane-Datei IST die Lane, keine Ueberschreibung. Eine Migration darf das nicht stillschweigend umdrehen und darf niemandem eine selbst geaenderte Lane unter den Fuessen wegziehen. Gesucht ist der Weg dazwischen: eine Korrektur erreicht das Board, eine Anpassung des Nutzers ueberlebt sie.
definition-of-done: "Ein Board, dessen done.md noch 'logbook-on-entry: true' traegt, fegt beim naechsten Move nach done keine fremden fertigen Tickets mehr ins Logbuch. Nachgestellt an einem Board-Fixture, das mit der alten Lane-Datei angelegt wurde."
tags:
  - gates
  - cli
blocked-by: []
related:
  - 01M28MHSDBABYVD8785A74VM40
commits: []
created-at: 2026-09-14T19:29:33Z
updated-at: 2026-09-15T06:57:32Z
assignee: "Alexander Sacharov"
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-40747
claimed-at: 2026-09-15T06:57:22Z
---

# Ein Board, das es schon gibt, bekommt eine geaenderte Lane nie zu sehen

## Definition of Done

- [ ] Ein Board, dessen done.md noch 'logbook-on-entry: true' traegt, fegt beim naechsten Move nach done keine fremden fertigen Tickets mehr ins Logbuch. Nachgestellt an einem Board-Fixture, das mit der alten Lane-Datei angelegt wurde.
- [ ] Eine vom Nutzer selbst geaenderte Lane ueberlebt die Migration unveraendert: was er geschrieben hat, wird nicht zurueckgesetzt. Die Entscheidung von 743737f - ein Board ist sein Lane-Verzeichnis - bleibt gueltig.
- [ ] Der Nutzer erfaehrt, dass eine seiner Lanes von einer Korrektur betroffen ist, statt es an seinem Verhalten zu merken. Nachgestellt an dem Board-Fixture aus Punkt 1.
- [ ] Die Zeile in core/release/NOTES.md, die 'Finishing a ticket no longer files anything' behauptet, stimmt danach fuer alle Boards - oder sie sagt, fuer welche sie nicht gilt und was zu tun ist.
- [ ] Eine Zeile in core/release/NOTES.md unter ## Unreleased.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

