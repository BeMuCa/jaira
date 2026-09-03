---
id: 01M0YT5A5QR240Z1W2JKCD9TCB
title: Ein absichtlicher Lane-Override soll nicht mehr warnen
status: human
ready: true
creator: BeMuCa
goal: "Eine Lane-Datei kann erklaeren, dass sie eine eingebaute Lane ersetzen will"
context: "Berks Board meldet bei jedem Befehl: 'lane review.md: id \"review\" overrides the built-in lane of the same name'. Die Warnung stimmt - er hat review.md am 20.08. bewusst zu einer Menschen-Lane gemacht - aber sie ist damit nutzlos: eine Warnung, die immer kommt, wird nicht mehr gelesen. Vorschlag: ein Feld wie 'overrides: review' im Lane-Kopf laesst sie verstummen, alles andere warnt weiter."
definition-of-done: Mit dem Feld ist die Warnung weg; ohne das Feld kommt sie unveraendert; jaira lanes show zeigt das Feld
blocked-by: []
commits: []
created-at: 2026-08-26T10:34:25Z
updated-at: 2026-08-31T16:38:20Z
claimed-by: EE-3NX6GL3-2641173
claimed-at: 2026-08-31T16:33:58Z
updated-by: BeMuCa
assignee: BeMuCa
question: "Die Warnung, die dieses Ticket stumm schalten soll ('lane X: id \"review\" overrides the built-in lane of the same name'), gibt es im aktuellen Code nicht mehr. Commit 743737f (27.08., nach diesem Ticket vom 26.08.) hat sie bewusst entfernt: eine vom Nutzer geschriebene Lane-Datei ist die Lane, keine Ueberschreibung — core/lane/lane.go warnt seit da nur noch bei einer Protection, die dabei verloren geht (droppedProtections, Invariante 4), nie mehr bei einer reinen Verhaltensaenderung. Test core/lane/lane_test.go:378 (TestLoadOverrideChangedPromptIsMarkedNotWarned) belegt das explizit. Damit ist Berks review.md-Fall aus dem Kontext vermutlich schon still, sofern sie keine Protection dropt. Wie soll ich weitermachen: (a) Ticket schliessen/archivieren, weil 743737f es bereits geloest hat, (b) ein overrides:-Feld trotzdem einfuehren, das rein dokumentarisch ist und nichts stummschaltet, oder (c) die alte Immer-Warnung bewusst wieder einfuehren (Reversal einer Design-Entscheidung), nur damit es etwas gibt, das das Feld stummschalten kann?"
outcome-what: "Recherche statt Code-Aenderung: geprueft, ob die Ticket-Praemisse (immer feuernde Override-Warnung) im aktuellen Code noch zutrifft"
outcome-why: "Vor der Implementierung muss die im Ticket beschriebene Warnung tatsaechlich existieren, sonst gibt es nichts, das ein neues Feld stummschalten koennte"
outcome-resolves: "Ungeloest: die Warnung existiert seit Commit 743737f nicht mehr, DoD wie geschrieben nicht erfuellbar, Frage an Berk gestellt"
---

# Ein absichtlicher Lane-Override soll nicht mehr warnen

## Definition of Done

- [ ] Mit dem Feld ist die Warnung weg; ohne das Feld kommt sie unveraendert; jaira lanes show zeigt das Feld

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-08-31 16:37 · BeMuCa** — Die im Ticket beschriebene Warnung existiert im aktuellen Code nicht mehr. Commit 743737f (27.08., 'a board is its lane directory') hat die einfache Meldung 'id %q overrides the built-in lane of the same name' bewusst entfernt (core/lane/lane.go, damals Zeile ~426) mit der Begruendung: eine vom Nutzer geschriebene Lane-Datei IST die Lane, keine Ueberschreibung von irgendwas. Seitdem warnt readLaneDir() (core/lane/lane.go:526-536) nur noch bei droppedProtections (Invariante 4, IMMER, unveraendert von diesem Ticket). TestLoadOverrideChangedPromptIsMarkedNotWarned (core/lane/lane_test.go:378) verifiziert explizit: ein geaenderter Prompt/Override OHNE Protection-Verlust erzeugt KEINE Warnung mehr. Berks review.md-Fall (Kontext des Tickets, 20.08.) waere also mit heutigem Code schon still - vorausgesetzt sie dropt keine Protection. Das DoD ('ohne Feld kommt sie unveraendert') setzt eine Warnung voraus, die es nicht mehr gibt: nicht umsetzbar wie geschrieben. Ticket an human mit Frage.
