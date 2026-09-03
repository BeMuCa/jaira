---
id: 01M1KFMETY9MKH4V38TVD0SAHM
title: Ein Epic-Layer zeigt Storys als verbundene Bubbles ueber dem Board
status: backlog
ready: true
creator: BeMuCa
goal: "Mit L wechselt das TUI auf einen Epic/Story-Layer: Bubbles (eine je Epic, in Tag-Farbe) bilden einen Flow von 'erst A, dann B'; jede Bubble sammelt die Tickets ihres Tags, fuellt sich mit deren Fertigstellungsgrad, und Enter springt vom Bubble-Detail zum Ticket aufs Board"
context: "Berk am 03.09., grosser Feature-Wunsch fuer ein naechstes Update - die Tag-Farben (79GEPW/AXTFG3) waren der erste Schritt dahin. Kern: eine zweite Ansicht (Taste L) mit Epics/Storys als Kreisen/Bubbles. Verbindungen als Pfeile = Reihenfolge (erst Thema A, dann B). Zuordnung Epic->Tickets laeuft ueber den Tag: Board-Farben zeigen, zu welchem Epic ein Ticket gehoert. Bedienung, wie von Berk beschrieben: initial eine Plus-Bubble; Enter oeffnet ein TUI-Template (wie Ticket-Anlage) und erstellt die Bubble (Default leer, zufaellige Farbe). Danach: rechts neben jeder Bubble eine Plus-Bubble (parallele/ungeordnete Epics, erste Reihe = Sammelreihe), darunter eine Plus-Bubble mit Pfeil (Nachfolger). Nach unten erstellte Bubbles bieten wieder rechts/unten an. y = Split/Abzweigung: von der Vorgaenger-Bubble eine zweite Kante zu einer neuen Bubble auf derselben Ebene wie die aktuelle (Berk ist offen fuer intuitivere Split-Gesten). Loeschen oeffnet einen Dialog: nur Bubble, oder Bubble samt Tickets ihres Tags - Default ohne Tickets. Fuellstand: Bubble fuellt sich mit Farbe nach Anteil fertiger Tickets ihres Tags. Bubble oeffnen: obere Haelfte die ausgefuellten Felder (z.B. Ziel des Epics), untere Haelfte die Ticket-Liste, dort auswaehlbar und mit Enter Sprung zum Ticket im Board. Vor dem Bauen zu brainstormen: Datenmodell (Bubble=Datei? Kanten wo?), Konfliktverhalten bei parallelen Sessions, ob Epic==Tag oder Epic hat einen Tag, Layout-Algorithmus im Terminal, und die Split-Geste."
definition-of-done: L wechselt zwischen Board und Epic-Layer; Bubbles anlegen/verbinden/splitten/loeschen wie im Kontext; Fuellstand nach fertigen Tickets des Tags; Bubble-Detail springt per Enter zum Ticket; das Datenmodell uebersteht parallele Sessions und git-Merge
tags: []
blocked-by: []
commits: []
created-at: 2026-09-03T11:14:30Z
updated-at: 2026-09-03T11:15:04Z
assignee: BeMuCa
updated-by: BeMuCa
---

# Ein Epic-Layer zeigt Storys als verbundene Bubbles ueber dem Board

## Definition of Done

- [ ] L wechselt zwischen Board und Epic-Layer; Bubbles anlegen/verbinden/splitten/loeschen wie im Kontext; Fuellstand nach fertigen Tickets des Tags; Bubble-Detail springt per Enter zum Ticket; das Datenmodell uebersteht parallele Sessions und git-Merge

## Options

- [x] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

