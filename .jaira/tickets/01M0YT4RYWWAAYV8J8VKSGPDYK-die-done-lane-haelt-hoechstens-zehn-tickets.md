---
id: 01M0YT4RYWWAAYV8J8VKSGPDYK
title: Die done-Lane haelt hoechstens zehn Tickets
status: human
ready: true
creator: BeMuCa
goal: "Aeltere Tickets in done wandern selbsttaetig ins Archiv, die zehn neuesten bleiben"
context: "Berks Vorschlag am 26.08.: done sammelt sich an und verstopft jede Liste. Regel: die zehn neuesten bleiben, alles aeltere geht nach .jaira/archive/. Offen und vor dem Bauen zu klaeren: wann das laeuft (bei jedem Board-Laden? bei einem eigenen Befehl? beim Move nach done?) und ob es still passieren darf - der Rest des Werkzeugs verschiebt nie eine Datei ohne Aufforderung. Archivieren ist umkehrbar (jaira restore), das spricht dafuer."
definition-of-done: "Bei mehr als zehn Tickets in done wird das aelteste archiviert; was passiert ist, steht sichtbar irgendwo; ein Test deckt die Grenze ab"
blocked-by: []
commits: []
created-at: 2026-08-26T10:34:07Z
updated-at: 2026-09-03T07:13:59Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-3483472
claimed-at: 2026-08-27T15:52:02Z
assignee: BeMuCa
question: |-
  Drei Entscheidungen, bevor gebaut wird:

  1. WOHIN gehen die aelteren Tickets? Das Ticket sagt .jaira/archive/. Seit heute heisst der Ort fuer fertige Arbeit .jaira/logbook/<initialen>-<datum>/ - und ein Ticket in done IST fertige Arbeit. Ins Logbuch getrimmt zaehlen sie spaeter in der Statistik pro Tag mit; im Archiv nicht. Aber: das Logbuch traegt die Initialen dessen, der den Move ausloest, nicht dessen, der das Ticket abgenommen hat, und das Datum des Trims, nicht der Abnahme. Empfehlung: Logbuch, weil 'fertig' dort hingehoert und Archiv 'nicht mehr verfolgt' heisst.

  2. WO steht die Zehn? (a) Im Lane-File der done-Lane als 'holds: 10' - sichtbar in 'jaira lanes show', pro Board aenderbar, 0 = unbegrenzt; die eingebaute done.md deklariert 10, Projektkopien bekommen es beim Drift-Refresh. (b) Fest im Code fuer die Terminal-Lane. Empfehlung (a): eine Regel, die Dateien bewegt, soll dort stehen, wo man sie lesen und abschalten kann.

  3. 'Aeltestes' = kleinstes updated-at. Einverstanden? (Alternative gibt es keine ohne neues Feld.)
outcome-what: "Plan geschrieben (5 Schritte), keine Codezeile."
outcome-why: "Drei Punkte sind Entscheidungen, keine Ableitungen: Ziel-Ordner, Ort der Grenze, Mass fuer 'aeltestes'."
outcome-resolves: "Nichts erfuellt; die DoD nennt 'Archiv' als Ziel, und genau das steht zur Frage."
---

# Die done-Lane haelt hoechstens zehn Tickets

## Definition of Done

- [ ] Bei mehr als zehn Tickets in done wird das aelteste archiviert; was passiert ist, steht sichtbar irgendwo; ein Test deckt die Grenze ab

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [ ] Grenze deklarieren: 'holds: 10' im Frontmatter der Terminal-Lane (core/lane parse + validate, 0 oder fehlend = unbegrenzt); builtin/50-done.md deklariert 10, Projektkopien holen es ueber den Drift-Refresh
- [ ] core/ticket: Tickets einer Lane nach updated-at sortieren und alles jenseits der neuesten N zurueckgeben (Store-Methode, Test mit 12 Tickets: genau die 2 aeltesten)
- [ ] Eine Trim-Funktion im core, aufgerufen nach einem gelungenen Move in eine Lane mit holds - drei Aufrufstellen: internal/cli/flow.go nach dem Mutate (~207), internal/tui/model.go applyMove (~1291) und der forcierte Pfad (~1343)
- [ ] Sagen, was passiert ist: CLI eine Zeile pro Ticket mit dem restore-Befehl, --json ein 'trimmed'-Feld, TUI ein notify; nie still
- [ ] Grenztests: genau 10 in done -> nichts bewegt; 11 -> das aelteste geht; Ziel und restore-Rueckweg geprueft

## Progress
- **2026-08-26 17:07 · BeMuCa** — Entschieden am 26.08.: der Schnitt laeuft beim Zug nach done UND sagt es. Beide offenen Fragen des Tickets sind damit beantwortet: wann (beim Zug, nicht beim Laden, nicht auf Befehl) und ob still (nein). Begruendung: automatisch, aber ohne die Regel zu brechen, dass in diesem Werkzeug nichts eine Datei unangekuendigt bewegt.
- **2026-08-27 15:52 · BeMuCa** — Entschieden (Berk, 26./27.08.): laeuft beim Move nach done, und sagt es - nicht beim Laden, nicht auf eigenen Befehl, nie still.

Gefunden beim Planen: es gibt keine gemeinsame Move-Funktion. Der CLI-Move schreibt in internal/cli/flow.go:207 (s.Mutate + SetReady), der TUI-Move in internal/tui/model.go:1291 (moveMutation) und noch einmal :1343 fuer den forcierten Pfad. Der Trim muss also an drei Stellen aufgerufen werden oder die drei bekommen erst eine gemeinsame Funktion - letzteres ist ein eigener Schnitt, nicht dieser.

'Aeltestes' = kleinstes updated-at. Ein Ticket in done wird normalerweise nicht mehr angefasst, also ist updated-at die Annahme-Zeit; eine spaete Notiz macht es juenger. Es gibt kein Feld fuer den Zeitpunkt des Lane-Wechsels. Verworfen: die ULID (Erstellzeit), sie misst das Falsche.
- **2026-09-03 07:13 · BeMuCa** — Empfehlung des Agenten (02.09., wartet auf Berks Ja): (1) Ziel = LOGBUCH (done ist fertige Arbeit; Archiv heisst 'wird nicht verfolgt'); (2) die Zehn als 'holds: 10' im Lane-File der done-Lane (sichtbar, pro Board aenderbar, 0 = unbegrenzt); (3) aeltestes = kleinstes updated-at; (4) Trigger = beim Move nach done, GEMELDET in der Move-Ausgabe (welche Tickets ins Logbuch gingen) - nie still, kein Board-Load-Nebeneffekt. Berks Frage vom 01.09. nach 'sync und nach 1 Monat archive' wurde abgeraten: wuerde die gerade vereinheitlichte Bedeutung (fertig->logbook, zurueckgezogen->archive) zerstoeren und braeuchte einen Beweger, den jaira nicht hat.
