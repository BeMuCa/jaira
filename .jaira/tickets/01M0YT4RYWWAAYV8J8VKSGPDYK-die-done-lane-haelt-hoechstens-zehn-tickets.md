---
id: 01M0YT4RYWWAAYV8J8VKSGPDYK
title: Die done-Lane haelt hoechstens zehn Tickets
status: critique
ready: true
creator: BeMuCa
goal: "Aeltere Tickets in done wandern selbsttaetig ins Archiv, die zehn neuesten bleiben"
context: "Berks Vorschlag am 26.08.: done sammelt sich an und verstopft jede Liste. Regel: die zehn neuesten bleiben, alles aeltere geht nach .jaira/archive/. Offen und vor dem Bauen zu klaeren: wann das laeuft (bei jedem Board-Laden? bei einem eigenen Befehl? beim Move nach done?) und ob es still passieren darf - der Rest des Werkzeugs verschiebt nie eine Datei ohne Aufforderung. Archivieren ist umkehrbar (jaira restore), das spricht dafuer."
definition-of-done: "Bei mehr als zehn Tickets in done wird das aelteste archiviert; was passiert ist, steht sichtbar irgendwo; ein Test deckt die Grenze ab"
blocked-by: []
commits: []
created-at: 2026-08-26T10:34:07Z
updated-at: 2026-09-03T11:35:08Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-2180407
claimed-at: 2026-09-03T11:14:22Z
assignee: BeMuCa
question: |-
  Drei Entscheidungen, bevor gebaut wird:

  1. WOHIN gehen die aelteren Tickets? Das Ticket sagt .jaira/archive/. Seit heute heisst der Ort fuer fertige Arbeit .jaira/logbook/<initialen>-<datum>/ - und ein Ticket in done IST fertige Arbeit. Ins Logbuch getrimmt zaehlen sie spaeter in der Statistik pro Tag mit; im Archiv nicht. Aber: das Logbuch traegt die Initialen dessen, der den Move ausloest, nicht dessen, der das Ticket abgenommen hat, und das Datum des Trims, nicht der Abnahme. Empfehlung: Logbuch, weil 'fertig' dort hingehoert und Archiv 'nicht mehr verfolgt' heisst.

  2. WO steht die Zehn? (a) Im Lane-File der done-Lane als 'holds: 10' - sichtbar in 'jaira lanes show', pro Board aenderbar, 0 = unbegrenzt; die eingebaute done.md deklariert 10, Projektkopien bekommen es beim Drift-Refresh. (b) Fest im Code fuer die Terminal-Lane. Empfehlung (a): eine Regel, die Dateien bewegt, soll dort stehen, wo man sie lesen und abschalten kann.

  3. 'Aeltestes' = kleinstes updated-at. Einverstanden? (Alternative gibt es keine ohne neues Feld.)
outcome-what: "done-Lane-Cap gebaut: 'holds: N' im Lane-Frontmatter (builtin done: 10, Template synchron), Store.Overflow/TrimLane in core/ticket/trim.go, Trim nach gelungenem Move an allen drei Schreibstellen (CLI flow.go, TUI applyMove, TUI forceMove), gemeldet als CLI-Zeile mit restore-Befehl, --json-Feld 'trimmed' und TUI-notify; 'jaira lanes show' zeigt den Cap"
outcome-why: "done sammelt sich an und verstopft jede Liste (Berks Vorschlag 26.08., Empfehlung am 03.09. angenommen: Logbuch, holds im Lane-File, updated-at, Trigger beim Move, nie still)"
outcome-resolves: "DoD verlangt: bei mehr als zehn in done geht das aelteste (erfuellt, Ziel Logbuch statt Archiv per Entscheidung), sichtbar (drei Meldewege), Test deckt Grenze (trim_test/holdcap_test/trimholds_test: 10->nichts, 11->aeltestes, restore-Rueckweg); go test ./... -race RC=0, Smoke auf frischem Board zeigt Holds: 10"
executed-by: fable
---

# Die done-Lane haelt hoechstens zehn Tickets

## Definition of Done

- [x] Bei mehr als zehn Tickets in done wird das aelteste archiviert; was passiert ist, steht sichtbar irgendwo; ein Test deckt die Grenze ab
  proof: Ziel=Logbuch (Berks Ja, 03.09.); TestMoveIntoACappedLaneTrimsTheOldest + TestTrimLaneFilesTheOverflowAndRestoreBringsItBack decken Grenze und Rueckweg; sichtbar: CLI-Zeile, --json trimmed, TUI-notify

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Grenze deklarieren: 'holds: 10' im Frontmatter der Terminal-Lane (core/lane parse + validate, 0 oder fehlend = unbegrenzt); builtin/50-done.md deklariert 10, Projektkopien holen es ueber den Drift-Refresh
  proof: core/lane/lane.go parse+Holds-Feld; core/lane/builtin/50-done.md holds:10; TestBuiltinDoneDeclaresTheCap, TestHoldsParsesAndDefaultsToUnlimited
- [x] core/ticket: Tickets einer Lane nach updated-at sortieren und alles jenseits der neuesten N zurueckgeben (Store-Methode, Test mit 12 Tickets: genau die 2 aeltesten)
  proof: core/ticket/trim.go Overflow (updated-at desc, tie=ID); TestOverflowReturnsExactlyTheOldestBeyondTheCap (12->genau die 2 aeltesten)
- [x] Eine Trim-Funktion im core, aufgerufen nach einem gelungenen Move in eine Lane mit holds - drei Aufrufstellen: internal/cli/flow.go nach dem Mutate (~207), internal/tui/model.go applyMove (~1291) und der forcierte Pfad (~1343)
  proof: core/ticket/trim.go TrimLane; Aufrufe: internal/cli/flow.go nach Mutate, internal/tui/model.go applyMove + forceMove via trimHolds
- [x] Sagen, was passiert ist: CLI eine Zeile pro Ticket mit dem restore-Befehl, --json ein 'trimmed'-Feld, TUI ein notify; nie still
  proof: CLI: Zeile je Ticket mit 'jaira restore'; --json: trimmed-Feld (trimmedJSON, flow.go); TUI: notify via trimHolds-Message; nie still
- [x] Grenztests: genau 10 in done -> nichts bewegt; 11 -> das aelteste geht; Ziel und restore-Rueckweg geprueft
  proof: trim_test.go (10->nichts, 11->aeltestes, Restore-Rueckweg), holdcap_test.go (CLI 3 Faelle), trimholds_test.go (TUI 2 Faelle); go test ./... -race RC=0

## Progress
- **2026-08-26 17:07 · BeMuCa** — Entschieden am 26.08.: der Schnitt laeuft beim Zug nach done UND sagt es. Beide offenen Fragen des Tickets sind damit beantwortet: wann (beim Zug, nicht beim Laden, nicht auf Befehl) und ob still (nein). Begruendung: automatisch, aber ohne die Regel zu brechen, dass in diesem Werkzeug nichts eine Datei unangekuendigt bewegt.
- **2026-08-27 15:52 · BeMuCa** — Entschieden (Berk, 26./27.08.): laeuft beim Move nach done, und sagt es - nicht beim Laden, nicht auf eigenen Befehl, nie still.

Gefunden beim Planen: es gibt keine gemeinsame Move-Funktion. Der CLI-Move schreibt in internal/cli/flow.go:207 (s.Mutate + SetReady), der TUI-Move in internal/tui/model.go:1291 (moveMutation) und noch einmal :1343 fuer den forcierten Pfad. Der Trim muss also an drei Stellen aufgerufen werden oder die drei bekommen erst eine gemeinsame Funktion - letzteres ist ein eigener Schnitt, nicht dieser.

'Aeltestes' = kleinstes updated-at. Ein Ticket in done wird normalerweise nicht mehr angefasst, also ist updated-at die Annahme-Zeit; eine spaete Notiz macht es juenger. Es gibt kein Feld fuer den Zeitpunkt des Lane-Wechsels. Verworfen: die ULID (Erstellzeit), sie misst das Falsche.
- **2026-09-03 07:13 · BeMuCa** — Empfehlung des Agenten (02.09., wartet auf Berks Ja): (1) Ziel = LOGBUCH (done ist fertige Arbeit; Archiv heisst 'wird nicht verfolgt'); (2) die Zehn als 'holds: 10' im Lane-File der done-Lane (sichtbar, pro Board aenderbar, 0 = unbegrenzt); (3) aeltestes = kleinstes updated-at; (4) Trigger = beim Move nach done, GEMELDET in der Move-Ausgabe (welche Tickets ins Logbuch gingen) - nie still, kein Board-Load-Nebeneffekt. Berks Frage vom 01.09. nach 'sync und nach 1 Monat archive' wurde abgeraten: wuerde die gerade vereinheitlichte Bedeutung (fertig->logbook, zurueckgezogen->archive) zerstoeren und braeuchte einen Beweger, den jaira nicht hat.
- **2026-09-03 11:34 · BeMuCa** — Gebaut wie empfohlen (Ja vom 03.09.): Ziel=Logbuch, holds:10 im Lane-File, aeltestes=kleinstes updated-at, Trigger=Move, immer gemeldet. Entscheidungen beim Bauen: (1) Overflow bricht ab, wenn List() Fehler meldet (auch PartialError) - auf einem Board mit unlesbaren Tickets wird nicht geraten, welches das aelteste ist. (2) Kein eigener Lock um den Trim; Rename schlaegt bei parallelem Trim einfach fehl und wird gemeldet - gleiche Exponierung wie 'jaira logbook'. (3) Kein Commit-Stempeln beim Trim: Tickets in done haben ihre Commits schon (requires-commits-Gate); 'jaira logbook' stempelt weiterhin. WICHTIG fuer dieses Board: die Projektkopien unter .jaira/lanes/ haben holds NICHT - der Cap wird hier erst aktiv, wenn done.md per Drift-Refresh nachgezogen wird, und done hat aktuell 21 Tickets: der erste done-Move danach schickt 11 ins Logbuch (gemeldet).
