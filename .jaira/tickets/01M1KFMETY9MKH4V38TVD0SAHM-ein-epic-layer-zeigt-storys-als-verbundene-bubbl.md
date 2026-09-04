---
id: 01M1KFMETY9MKH4V38TVD0SAHM
title: Ein Epic-Layer zeigt Storys als verbundene Bubbles ueber dem Board
status: brainstorm
ready: true
creator: BeMuCa
goal: "Mit L wechselt das TUI auf einen Epic/Story-Layer: Bubbles (eine je Epic, in Tag-Farbe) bilden einen Flow von 'erst A, dann B'; jede Bubble sammelt die Tickets ihres Tags, fuellt sich mit deren Fertigstellungsgrad, und Enter springt vom Bubble-Detail zum Ticket aufs Board"
context: "Berk am 03.09., grosser Feature-Wunsch fuer ein naechstes Update - die Tag-Farben (79GEPW/AXTFG3) waren der erste Schritt dahin. Kern: eine zweite Ansicht (Taste L) mit Epics/Storys als Kreisen/Bubbles. Verbindungen als Pfeile = Reihenfolge (erst Thema A, dann B). Zuordnung Epic->Tickets laeuft ueber den Tag: Board-Farben zeigen, zu welchem Epic ein Ticket gehoert. Bedienung, wie von Berk beschrieben: initial eine Plus-Bubble; Enter oeffnet ein TUI-Template (wie Ticket-Anlage) und erstellt die Bubble (Default leer, zufaellige Farbe). Danach: rechts neben jeder Bubble eine Plus-Bubble (parallele/ungeordnete Epics, erste Reihe = Sammelreihe), darunter eine Plus-Bubble mit Pfeil (Nachfolger). Nach unten erstellte Bubbles bieten wieder rechts/unten an. y = Split/Abzweigung: von der Vorgaenger-Bubble eine zweite Kante zu einer neuen Bubble auf derselben Ebene wie die aktuelle (Berk ist offen fuer intuitivere Split-Gesten). Loeschen oeffnet einen Dialog: nur Bubble, oder Bubble samt Tickets ihres Tags - Default ohne Tickets. Fuellstand: Bubble fuellt sich mit Farbe nach Anteil fertiger Tickets ihres Tags. Bubble oeffnen: obere Haelfte die ausgefuellten Felder (z.B. Ziel des Epics), untere Haelfte die Ticket-Liste, dort auswaehlbar und mit Enter Sprung zum Ticket im Board. Vor dem Bauen zu brainstormen: Datenmodell (Bubble=Datei? Kanten wo?), Konfliktverhalten bei parallelen Sessions, ob Epic==Tag oder Epic hat einen Tag, Layout-Algorithmus im Terminal, und die Split-Geste."
definition-of-done: L wechselt zwischen Board und Epic-Layer; Bubbles anlegen/verbinden/splitten/loeschen wie im Kontext; Fuellstand nach fertigen Tickets des Tags; Bubble-Detail springt per Enter zum Ticket; das Datenmodell uebersteht parallele Sessions und git-Merge
tags: []
blocked-by: []
commits: []
created-at: 2026-09-03T11:14:30Z
updated-at: 2026-09-04T16:50:53Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-4179675
claimed-at: 2026-09-04T16:44:21Z
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
- **2026-09-04 16:50 · BeMuCa** — Brainstorm-Entscheidungspunkte (an Berk gestellt, 04.09.): K1 Speicherung: eine Markdown-Datei je Epic unter .jaira/epics/ (Frontmatter: name, tag, goal, after: [Vorgaenger-Epics] als Kanten) - hand-editierbar, diff-bar, Kantenliste merge=union. K2 Bindung: Epic HAT einen Tag (Feld) statt Epic==Tag - umbenennbar, Farbe bleibt allein in .jaira/tags; Fuellstand = fertige/alle Tickets des Tags. K3 Parallel-Sessions: WriteAtomic + Store-Lock wie core/tag. K4 Layout: Ebene = laengster Pfad von den Wurzeln; erste Reihe = Sammelreihe (parallel), darunter Nachfolger; Kanten als einfache Linien. K5 Split-Geste: Vorschlag 'v'=Kante-ziehen-Modus (deckt Split UND Merge) statt nur y; Berks y als Alternative. K6 Loeschen: Dialog Bubble-only (Default) vs. 'mit Tickets' = Tickets ARCHIVIEREN (umkehrbar), nie hart loeschen. K7 Anlage: L toggelt Layer, Plus-Bubbles rechts (parallel) und unten (Nachfolger), Enter -> TUI-Formular (name, ziel, tag, farbe optional, Default random).
