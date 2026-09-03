---
id: 01M0YT4RYWWAAYV8J8VKSGPDYK
title: Die done-Lane haelt hoechstens zehn Tickets
status: signoff
ready: true
creator: BeMuCa
goal: "Aeltere Tickets in done wandern selbsttaetig ins Archiv, die zehn neuesten bleiben"
context: "Berks Vorschlag am 26.08.: done sammelt sich an und verstopft jede Liste. Regel: die zehn neuesten bleiben, alles aeltere geht nach .jaira/archive/. Offen und vor dem Bauen zu klaeren: wann das laeuft (bei jedem Board-Laden? bei einem eigenen Befehl? beim Move nach done?) und ob es still passieren darf - der Rest des Werkzeugs verschiebt nie eine Datei ohne Aufforderung. Archivieren ist umkehrbar (jaira restore), das spricht dafuer."
definition-of-done: "Bei mehr als zehn Tickets in done wird das aelteste archiviert; was passiert ist, steht sichtbar irgendwo; ein Test deckt die Grenze ab"
blocked-by: []
commits: []
created-at: 2026-08-26T10:34:07Z
updated-at: 2026-09-03T12:11:47Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-2180407
claimed-at: 2026-09-03T11:14:22Z
assignee: BeMuCa
question: |-
  Drei Entscheidungen, bevor gebaut wird:

  1. WOHIN gehen die aelteren Tickets? Das Ticket sagt .jaira/archive/. Seit heute heisst der Ort fuer fertige Arbeit .jaira/logbook/<initialen>-<datum>/ - und ein Ticket in done IST fertige Arbeit. Ins Logbuch getrimmt zaehlen sie spaeter in der Statistik pro Tag mit; im Archiv nicht. Aber: das Logbuch traegt die Initialen dessen, der den Move ausloest, nicht dessen, der das Ticket abgenommen hat, und das Datum des Trims, nicht der Abnahme. Empfehlung: Logbuch, weil 'fertig' dort hingehoert und Archiv 'nicht mehr verfolgt' heisst.

  2. WO steht die Zehn? (a) Im Lane-File der done-Lane als 'holds: 10' - sichtbar in 'jaira lanes show', pro Board aenderbar, 0 = unbegrenzt; die eingebaute done.md deklariert 10, Projektkopien bekommen es beim Drift-Refresh. (b) Fest im Code fuer die Terminal-Lane. Empfehlung (a): eine Regel, die Dateien bewegt, soll dort stehen, wo man sie lesen und abschalten kann.

  3. 'Aeltestes' = kleinstes updated-at. Einverstanden? (Alternative gibt es keine ohne neues Feld.)
outcome-what: "Review-Funde umgesetzt: accept() als vierte Schreibstelle trimmt, Tie-Break zur neueren ULID plus newest-Pin fuer den ausloesenden Move, Partial-Trim wird immer gemeldet, CLI-Trim-Fehler faellt den gelungenen Move nicht mehr (stderr-Warnung + trim_error im JSON), holds nur noch auf terminalen Lanes"
outcome-why: "Das unabhaengige Opus-Review fand 5 reproduzierte Funde - darunter die vierte Move-Schreibstelle (menschlicher Accept-Weg), auf der der Cap gar nicht feuerte"
outcome-resolves: "Alle 5 Funde geschlossen und je mit Test gepinnt (accept-Pfad, Ties, Pin, trim_error, terminal-Guard); go test ./... -race RC=0, Cache geleert; Binary neu gebaut"
executed-by: fable
review-summary: "Zweite Runde nach 5 Review-Funden. Kritik am Fix-Diff: jeder Fund minimal geschlossen, kein Beifang - accept() ruft denselben trimHolds wie die Move-Pfade (kein neues Muster), der newest-Pin ist ein Parameter statt einer zweiten Sortierlogik, der CLI-Fehlerpfad aendert nur Meldeweg (stderr + trim_error), nicht Semantik des Moves. Terminal-Guard sitzt im parse, wo der Lane-Autor die Begruendung liest. Der Meldungssatz existiert weiterhin doppelt (flow.go/model.go) - unveraendert bewusster Zustand."
review-gaps: "Nichts entfernt. Gelassen: Negativ-Check auf holds; 'trimmed' im --json immer praesent (stabiles Schema), trim_error nur bei Fehler; Abbruch bei PartialError von List (nicht raten, welches Ticket das aelteste ist); Fixture-Helfer je Package. Bekannte Restluecke: Logbook-Namenskollision MITTEN im Trim (Partial-Fall) ohne deterministischen Test - das Melde-Verhalten ist implementiert, die Kollision aber nicht stabil zu inszenieren ohne logbookFolder()-Interna im Test nachzubauen."
review-check: "1. Neu bauen; frisches Board: 'jaira lanes show done' -> 'Holds: 10'. 2. Elf Tickets nach done zwingen -> der elfte Move meldet die Logbuch-Zeile mit restore-Befehl; genauso beim Accept mit 'a' im signoff (vierte Schreibstelle, jetzt gedeckt). 3. Drei Moves in EINER Sekunde in eine holds-Lane -> es geht das aelteste, nie der Neuankoemmling. 4. move --json -> 'trimmed'; bei kaputtem Ticket auf dem Board zusaetzlich 'trim_error', Exit bleibt 0, der Move steht. 5. holds: 3 auf einer nicht-terminalen Lane -> 'jaira lanes' verweigert mit Begruendung. 6. Auf DIESEM Board greift der Cap erst nach Drift-Refresh von done.md - done haelt 21: der naechste done-Move/Accept danach schickt 11 ins Logbuch, gemeldet."
review-verdict: "accept (Opus-End-Review auf b1c4ecb, zweiter Durchgang nach 5 Funden aus dem ersten - alle mit eigenen Proben verifiziert gefixt: accept-Pfad per Overlay-Test ausgefuehrt, Tie/Pin ueber 550 Kombinationen gemessen, trim_error/stderr-Verhalten am echten Binary probiert, 'grep trimHolds' bestaetigt alle vier Schreibstellen). Restnotizen des Reviewers, nicht blockierend: der Partial-Trim-Fall (Kollision mitten im Trim) ist nur per Probe gepinnt, nicht als Test; Board-Anzeige tie't auf ID< waehrend der Trim auf ID> tie't - kosmetisch. Nicht geprueft: echtes Tastendruecken im Terminal (nur ueber Go-Tests) und dieses Board selbst (holds erst nach Drift-Refresh)."
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
  proof: core/ticket/trim.go TrimLane (mit newest-Pin); VIER Aufrufstellen: cli/flow.go nach Mutate, tui/model.go applyMove + forceMove, tui/signoff.go accept() - die vierte fand das Opus-Review
- [x] Sagen, was passiert ist: CLI eine Zeile pro Ticket mit dem restore-Befehl, --json ein 'trimmed'-Feld, TUI ein notify; nie still
  proof: CLI: Zeile je Ticket mit 'jaira restore'; --json: trimmed-Feld (trimmedJSON, flow.go); TUI: notify via trimHolds-Message; nie still
- [x] Grenztests: genau 10 in done -> nichts bewegt; 11 -> das aelteste geht; Ziel und restore-Rueckweg geprueft
  proof: trim_test.go (Grenze, Restore, Ties, Pin), holdcap_test.go (CLI 5 Faelle inkl. trim_error), trimholds_test.go + signoffaccept_test.go (TUI inkl. accept-Pfad), holds_test.go (parse inkl. terminal-Guard); go test ./... -race RC=0 nach Fix-Runde

## Progress
- **2026-08-26 17:07 · BeMuCa** — Entschieden am 26.08.: der Schnitt laeuft beim Zug nach done UND sagt es. Beide offenen Fragen des Tickets sind damit beantwortet: wann (beim Zug, nicht beim Laden, nicht auf Befehl) und ob still (nein). Begruendung: automatisch, aber ohne die Regel zu brechen, dass in diesem Werkzeug nichts eine Datei unangekuendigt bewegt.
- **2026-08-27 15:52 · BeMuCa** — Entschieden (Berk, 26./27.08.): laeuft beim Move nach done, und sagt es - nicht beim Laden, nicht auf eigenen Befehl, nie still.

Gefunden beim Planen: es gibt keine gemeinsame Move-Funktion. Der CLI-Move schreibt in internal/cli/flow.go:207 (s.Mutate + SetReady), der TUI-Move in internal/tui/model.go:1291 (moveMutation) und noch einmal :1343 fuer den forcierten Pfad. Der Trim muss also an drei Stellen aufgerufen werden oder die drei bekommen erst eine gemeinsame Funktion - letzteres ist ein eigener Schnitt, nicht dieser.

'Aeltestes' = kleinstes updated-at. Ein Ticket in done wird normalerweise nicht mehr angefasst, also ist updated-at die Annahme-Zeit; eine spaete Notiz macht es juenger. Es gibt kein Feld fuer den Zeitpunkt des Lane-Wechsels. Verworfen: die ULID (Erstellzeit), sie misst das Falsche.
- **2026-09-03 07:13 · BeMuCa** — Empfehlung des Agenten (02.09., wartet auf Berks Ja): (1) Ziel = LOGBUCH (done ist fertige Arbeit; Archiv heisst 'wird nicht verfolgt'); (2) die Zehn als 'holds: 10' im Lane-File der done-Lane (sichtbar, pro Board aenderbar, 0 = unbegrenzt); (3) aeltestes = kleinstes updated-at; (4) Trigger = beim Move nach done, GEMELDET in der Move-Ausgabe (welche Tickets ins Logbuch gingen) - nie still, kein Board-Load-Nebeneffekt. Berks Frage vom 01.09. nach 'sync und nach 1 Monat archive' wurde abgeraten: wuerde die gerade vereinheitlichte Bedeutung (fertig->logbook, zurueckgezogen->archive) zerstoeren und braeuchte einen Beweger, den jaira nicht hat.
- **2026-09-03 11:34 · BeMuCa** — Gebaut wie empfohlen (Ja vom 03.09.): Ziel=Logbuch, holds:10 im Lane-File, aeltestes=kleinstes updated-at, Trigger=Move, immer gemeldet. Entscheidungen beim Bauen: (1) Overflow bricht ab, wenn List() Fehler meldet (auch PartialError) - auf einem Board mit unlesbaren Tickets wird nicht geraten, welches das aelteste ist. (2) Kein eigener Lock um den Trim; Rename schlaegt bei parallelem Trim einfach fehl und wird gemeldet - gleiche Exponierung wie 'jaira logbook'. (3) Kein Commit-Stempeln beim Trim: Tickets in done haben ihre Commits schon (requires-commits-Gate); 'jaira logbook' stempelt weiterhin. WICHTIG fuer dieses Board: die Projektkopien unter .jaira/lanes/ haben holds NICHT - der Cap wird hier erst aktiv, wenn done.md per Drift-Refresh nachgezogen wird, und done hat aktuell 21 Tickets: der erste done-Move danach schickt 11 ins Logbuch (gemeldet).
- **2026-09-03 11:48 · BeMuCa** — Opus-Review (2. Modell, eigene Proben): 5 Funde, alle reproduziert. (1) internal/tui/signoff.go accept() ist die VIERTE Status-Schreibstelle und trimmt nicht - ausgerechnet der menschliche Abnahme-Weg nach done. (2) trim.go Tie-Break invertiert: bei gleichem updated-at (Sekundenaufloesung!) fliegt das NEUSTE Ticket ins Logbuch (ID< statt ID>); 3x reproduziert. (3) Teilfehler im Trim: schon bewegte Tickets werden verschluckt statt gemeldet (Callers werfen das partial-Ergebnis weg) - verletzt 'nie still'. (4) CLI meldet nach GELUNGENEM Move exit 1 und verliert Erfolgszeile/JSON, Retry ist already_in_lane und trimmt nie - Lane bleibt fuer immer ueber dem Cap. (5) holds auf nicht-terminaler Lane schickt UNFERTIGE Arbeit ins Logbuch (umgeht Terminal-Guard von 'jaira logbook', zaehlt in LoggedPerDay als fertig). Gaps: kein Tie-Test, kein accept()-Test, Entscheidungen 3/4 ungepinnt.
- **2026-09-03 11:59 · BeMuCa** — Alle 5 Review-Funde gefixt, Suite gruen (-race, Cache geleert): (1) signoff.go accept() trimmt jetzt (vierte Schreibstelle; Test TestAcceptEnforcesTheCap, adaptiert von der Review-Probe). (2) Tie-Break gedreht (ID> = neuere ULID bleibt) UND der ausloesende Move ist als 'newest' gepinnt - Overflow/TrimLane tragen den Parameter; Tests TestOverflowBreaksTiesTowardTheNewerULID + TestOverflowNeverPicksTheJustMovedTicket. (3) Partial-Trim wird ueberall gemeldet (Callers zeigen trimmed auch bei err). (4) CLI: Trim-Fehler faellt den Move nicht mehr - Erfolgszeile/JSON bleiben, Warnung auf stderr, --json bekommt trim_error; Tests TestMoveSurvivesATrimFailureAndSaysSo + TestMoveJSONCarriesTheTrimError. (5) parse verweigert holds>0 auf nicht-terminaler Lane (Logbuch = fertige Arbeit); TestHoldsRefusesANonTerminalLane. Restluecke ehrlich benannt: Logbook-Namenskollision mitten im Trim (F3-Szenario) hat keinen deterministischen Test - Verhalten ist aber jetzt: alles bereits Bewegte wird genannt.
