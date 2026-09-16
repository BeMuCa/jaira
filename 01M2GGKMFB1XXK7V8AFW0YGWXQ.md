---
id: 01M2GGKMFB1XXK7V8AFW0YGWXQ
title: "Ein Sprint ist eine eigene Datei, keine Markierung am einzelnen Ticket"
status: in-progress
ready: true
creator: Alexander Sacharov
goal: "Wer plant, legt einen Milestone als eigene Datei an, die die zugehoerigen Tickets aufzaehlt, sieht deren Farbe am rechten Rand jeder Karte und zieht das Board mit einem Griff auf diesen Milestone zusammen - eine Datei bearbeiten statt zwanzig Tickets einzeln anzufassen."
context: |-
  Beim Planen am 2026-09-14 lagen 80-90 alte Tickets auf dem Board und 20 neue sollten dazukommen. Gebraucht wurde ein Weg, die neuen in der Sitzung schnell durchzugehen, waehrend die alten daneben liegen bleiben.

  Der erste Versuch war ein Tag, sprint-0914, und er funktioniert fuer genau eine Sitzung. Danach nicht mehr, und das ist der Grund fuer dieses Ticket.

  Was an einem Tag nicht geht: unerledigte Arbeit muss im naechsten Sprint wieder auftauchen. Ein Tag haengt am einzelnen Ticket, also heisst Weitertragen, jedes Ticket einzeln anzufassen - und Tickets liegen auf ihren Refs, das ist je Ticket ein pull, eine Aenderung und ein release. Bei zwanzig Tickets ist das der Grund, warum es niemand macht.

  Eine Datei wird stattdessen kopiert: die Zeilen, die nicht fertig sind, wandern in einem Griff in die Datei des naechsten Sprints. Aus zwanzig Vorgaengen wird eine Bearbeitung.

  Alex' Entwurf, am 2026-09-14 entschieden:
  - Der Sprint ist eine eigene Datei. Sie fuehrt die Tickets auf, die zu ihm gehoeren, und wird wie ein Ticket von Hand editiert und im Diff gelesen.
  - Sie reist auf einem Ref, damit sie schnell bei allen ankommt und nicht auf einen Merge wartet. Der Snapshot-Zweig jaira/board ist ein moeglicher Ort.
  - Die Farbe wird zufaellig vergeben. Gleichzeitige Sprints sind nie viele, ein Zusammenfall ist also verschmerzbar und niemand soll eine Farbe aussuchen muessen.
  - Auf der Karte steht der Sprint RECHTS. Links liegen die Tag-Farben (S1VM40), rechts eine Farbe, die sagt welche Tickets zusammengehoeren.
  - Es braucht einen Filter auf den Sprint, damit sich das Board in einer Geste auf die wichtigen Tickets zusammenzieht.

  Verhaeltnis zu S1VM40: dort bekam die linke Leiste drei Plaetze, und der dritte war fuer die Sprint-Markierung reserviert. Weil der Sprint jetzt nach rechts wandert, ist der dritte Platz links frei und gehoert dem dritten Tag. Diese Aenderung gehoert zu S1VM40 und nicht hierher.

  Nicht entschieden und Sache der Plan-Lane: das genaue Format der Datei, wo sie liegt, und ob ein Ticket in mehr als einem Sprint stehen darf.
definition-of-done: "Ein Sprint laesst sich anlegen und benennen, und seine Datei fuehrt die Tickets auf, die zu ihm gehoeren. Sie ist von Hand editierbar und im Diff lesbar, wie eine Ticket-Datei."
tags:
  - tui
  - cli
blocked-by: []
related:
  - 01M2FQEEQN61ZE9AJ4Y4S1VM40
commits:
  - bb903c8d77f82fae3840895ccd46a9ad569ef84f
  - 85cfde30113a1c24da62ca7136bdcaacdbb6d8b6
  - 5b0d17ce5a712900708342b600c0e59dcef88267
  - c359211dbc52d7cca6286f32e06eb3e8c736a242
  - ade63fe0eac8077144f48ef491da073ea7176087
  - 29afd307dee1524f4d96da72e094c13015c125f8
  - 4078d9774653ab9b785d5a84d0b5caf0009529c9
  - c08ecb911b1d5a686c213bc7e717f6dcb0b954b0
  - 2ff06a626737804dcdc2ff5f05b36efa898c0e37
created-at: 2026-09-14T17:49:30Z
updated-at: 2026-09-16T06:54:04Z
assignee: "Alexander Sacharov"
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-35903
claimed-at: 2026-09-16T06:51:55Z
outcome-what: "Reworded the two refusals a user meets when a milestone was filed somewhere else: 'jaira logbook <name>' on an already-filed milestone now names the marked file and says the restore has to run in the tree that filed it, instead of pointing at a 'jaira restore' that fails here; 'jaira milestone create' on a filed name no longer offers hand-deleting 'status: filed' as the way back, and core/release/NOTES.md and docs/COMMANDS.md say the same."
outcome-why: "Both messages sent the reader down a path that does not work. Every state reaching the logbook refusal is a tree without a logbook copy, so its 'jaira restore <name>.md' answers 'is not in the archive'. And deleting the mark by hand in a second clone rides back out on the ref, pulls the milestone onto that board alone and strands the logbook copy in the filer's tree, where 'jaira restore' then hits 'is already on the board'."
outcome-resolves: "critique round 6, findings 1 and 2. No behaviour changed; go build, go vet and go test ./core/... ./internal/cli/... are green."
review-summary: "internal/cli/milestones.go:237 editMembers loads the milestone with milestone.Load, which does not filter Filed(), and never asks ms.Filed() — so 'jaira milestone add/rm <name>' on a marked file still on disk succeeds and prints '<handle> added to <name>' and '<name> now holds N ticket(s)' for a milestone milestone.LoadAll hides, 'jaira milestone ls' does not name and no card paints. That file is a normal state, not a corner: core/refsync/refsync.go:IncomingMilestones writes a filed ref over a file this tree already has, and a filing that marked but did not move leaves one too — the two states core/milestone/milestone.go:LoadAll's doc names. Both other doors into that file already have a refusal (internal/cli/milestones.go:121 for create, internal/cli/logbook.go:283 for logbook); add/rm is the third and has none. Refuse there the same way: name the file, its mark, and that 'jaira restore <name>.md' in the tree that filed it is what puts it back — the wording at milestones.go:122 is the model. No test covers this; internal/cli/milestones_test.go has no add/rm case on a marked file."
review-gaps: "Entfernt: outbox.QueueKind/PendingKind/DropKind sind unexportiert (queueKind/pendingKind/dropKind) - kein Aufrufer ausserhalb core/outbox, auch kein Test; die drei kind.or(KindTicket)-Zeilen darin und die in Box.path sind weg, weil jeder Aufrufer den Kind selbst benennt oder ihn normalisiert von der Platte bekommt (Kind.or bleibt dort, wo Kind aus JSON kommt: readEntry-Pfad, readDir, Flush). milestoneJSON ruft ms.Members() einmal statt zweimal - jeder Aufruf kopierte die ganze Slice. Stehengelassen und warum: milestone.parse duplziert die Frontmatter-Lesung von ticket.ParseDoc nur scheinbar - ParseDoc lehnt eine kaputte Datei ab und kann keine Body-Zeilen editieren, milestone muss beides koennen, ein Umbau waere eine Verhaltensaenderung; cardColors/milestoneColors teilen die Form, nicht die Quelle (Registry vs Index), ein gemeinsamer Helfer waere ein Callback und laenger; Index.Matches normalisiert je Ticket, genau wie das vorhandene tag.Matches daneben in tickets.go:507 - dieselbe Kosten, gleiche Stelle, kein Grund nur die eine Haelfte zu aendern; gitref.Root/MilestonePrefix und milestone.Subdir sind exportiert ohne externen Aufrufer, benennen aber das Ref- bzw. Platten-Layout wie das vorhandene gitref.Prefix und ticket.DirName. Vorhandener toter Code nicht angefasst (staticcheck U1000, alle drei aelter als dieser Branch): internal/cli/share.go:17 isShared, internal/tui/model.go:256 laneStart, internal/tui/model.go:609 currentLane."
test-verdict: "pass: Suite gruen (build/vet/go test ./... -race, Cache geleert, RC=0), DoD 1-7 im Baum nachgeprueft, Verhalten mit dem echten Binary auf einem Scratch-Board und zwei Clones ausgefuehrt"
question: "Testing ist durch: build/vet/test -race gruen, DoD 1-7 nachgeprueft, Milestone-Anlegen, Hand-Edit-Weitertragen, Ref-Transport und Board-Filter am echten Binary vorgefuehrt. Nimmst du die Arbeit an, oder soll noch etwas geprueft werden, bevor sie in review geht?"
---

# Ein Sprint ist eine eigene Datei, keine Markierung am einzelnen Ticket

## Definition of Done

- [x] Ein Milestone laesst sich anlegen und benennen, und seine Datei fuehrt die Tickets auf, die zu ihm gehoeren. Sie ist von Hand editierbar und im Diff lesbar, wie eine Ticket-Datei.
  proof: core/milestone/milestone_test.go:TestSaveKeepsHandEditsVerbatim
- [x] Die Milestone-Datei reist auf einem Ref: wer sie zieht, sieht denselben Milestone, ohne auf das Mergen eines Zweiges zu warten.
  proof: core/gitref/milestone_test.go:TestMilestoneArrivesWithoutASharedBranch
- [x] Jeder Milestone bekommt seine Farbe, ohne dass jemand eine aussucht.
  proof: internal/cli/milestones_test.go:TestCreateAssignsAColourIntoTheMilestoneFile
- [x] Eine Karte zeigt die Farbe ihres Milestones am RECHTEN Rand, deutlich getrennt von den Tag-Farben am linken; ein Ticket ohne Milestone zeigt dort nichts und die Karte wird dadurch nicht breiter. Mehrfachzugehoerigkeit ist erlaubt: rechts stehen bis zu drei Plaetze, in Dateireihenfolge.
  proof: internal/tui/milestonebar_test.go:TestRightBarWidthIsTheSameWithAndWithoutMilestones
- [x] Das Board zieht sich mit einer Geste auf einen Milestone zusammen, und 'jaira list' hat den entsprechenden Schalter.
  proof: internal/tui/milestonebar_test.go:TestPickerNarrowsTheBoardAndReleasesIt
- [x] Unerledigte Arbeit wandert in den naechsten Sprint, indem eine Datei bearbeitet wird - nicht indem jedes Ticket einzeln angefasst wird. Nachgestellt mit mindestens drei Tickets, von denen zwei weiterwandern.
  proof: testing lane, reenacted on a scratch board: sprint-1 with 3 tickets, 2 lines moved by hand into sprint-2.md, one edit — jaira list --milestone then shows 1 and 2
- [x] Je eine Zeile in core/release/NOTES.md unter ## Unreleased fuer das, was ein Benutzer davon merkt.
  proof: core/release/NOTES.md:17
- [x] Ein Milestone verschwindet nur auf Kommando, nie von allein: ein leer geraeumter Milestone bleibt stehen, bis jemand ihn ablegt. Nachgestellt, indem das letzte Ticket herausgenommen wird - danach nennt 'jaira milestone ls' ihn unveraendert.
  proof: internal/cli/milestones_test.go:TestRmDropsTheLineAndKeepsTheMilestone, TestCreateWithNoTicketsLeavesTheFileLyingThere
- [x] 'jaira logbook' legt einen Milestone ab wie ein Ticket: er wandert unter .jaira/logbook/, 'jaira milestone ls' nennt ihn nicht mehr, das Board zeigt ihn nicht und der M-Filter kennt ihn nicht. 'jaira restore' holt ihn zurueck, samt Ticket-Liste und Farbe.
  proof: internal/cli/milestones_test.go:TestFilingAMilestoneTakesItOffTheBoardAndRestoreBringsItBack
- [x] Der Ref eines abgelegten Milestones wird NICHT geraeumt: er bleibt stehen und traegt im Frontmatter den Status 'abgelegt'. Nachgestellt an zwei Arbeitsbaeumen - der zweite zieht die Refs und schreibt den abgelegten Milestone NICHT wieder aufs Board.
  proof: internal/cli/milestoneref_test.go:TestAFiledMilestoneStaysOffTheOtherCloneAndKeepsItsRef
- [x] Der Name eines abgelegten Milestones ist belegt: 'jaira milestone create' mit demselben Namen wird abgelehnt und sagt, dass dieser Milestone abgelegt ist und mit 'jaira restore' zurueckkommt.
  proof: internal/cli/milestones_test.go:TestFilingAMilestoneTakesItOffTheBoardAndRestoreBringsItBack (milestone create refused, names 'jaira restore')
- [x] Je eine Zeile in core/release/NOTES.md unter ## Unreleased fuer das Ablegen eines Milestones und fuer den belegten Namen.
  proof: core/release/NOTES.md:24 (filing) and :25 (the taken name)

## Options

- [x] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Format festlegen und als Paket-Doku in core/milestone hinschreiben: .jaira/milestones/<name>.md, Frontmatter name/colour/created-at, Body eine Zeile je Mitglied - wie eine Ticket-Datei, und eine Zeile ist das Kleinste, was git mergen kann
- [x] core/milestone: Load/LoadAll/Save nach dem Vorbild von core/tag (core/tag/tag.go:189 Load, :371 Save) - Zeilen verbatim erhalten, WriteAtomic, Mitgliederliste parsen
- [x] Farbvergabe: zufaellig aus einer Palette, die keine Farbe eines schon vorhandenen Milestones doppelt; die Farbe steht in der Milestone-Datei selbst, kein zweites zentrales Registry
- [x] Test: eine Milestone-Datei von Hand editieren (Kommentar, Leerzeile, eigene Reihenfolge) und nach Load/Save unveraendert wiederfinden
- [x] Index ID -> Milestones einmal beim Laden bauen, neben link.Build in loadEnv (internal/cli/root.go:275) und im TUI (internal/tui/model.go:352); CLI und TUI lesen denselben Index
- [x] CLI: jaira milestone create/add/rm/ls - je Aufruf genau ein Dateischreibvorgang, damit zwanzig Tickets gruppieren eine Bearbeitung bleibt
- [x] CLI: jaira list --milestone <name> (neben --tag, internal/cli/tickets.go:528) und der Schluessel milestone:<name> in beiden matches() (internal/cli/tickets.go:534, internal/tui/model.go:598) - exakt wie tag, nicht als Teilstring
- [x] TUI: rechte Kartenkante in renderCardBlock (internal/tui/view.go:533) - zweite Balkenzelle rechts, inner = w-2, die Spalte IMMER reserviert, auch ohne Milestone, sonst flattern die Titel in einer Lane
- [x] TUI: milestoneColors gespiegelt zu cardColors (internal/tui/model.go:1366) - cardSlots Plaetze, in Dateireihenfolge, ein vierter Milestone faerbt nichts
- [x] Test: Karte mit 0, 1 und 4 Milestones - gleiche Kartenbreite, gleiche Textbreite, hoechstens drei gefaerbte Plaetze rechts
- [x] TUI: die Geste - Picker wie die Tag-Box auf 't' (internal/tui/model.go:1010), die Auswahl setzt m.filter auf milestone:<name> und nutzt damit den vorhandenen Filterweg
- [x] gitref: Namensraum von 'Ticket-ID' auf '(Art, Name)' verallgemeinern - Prefix (core/gitref/gitref.go:39), RefName (:144), Fetch-Refspec (:575), List/ListRemote/idsFrom (:587-625); refs/jaira/milestones/<name> neben refs/jaira/tickets/<id>
- [x] snapshot: pruefen, dass reap (core/snapshot/snapshot.go:234) nur Ticket-Refs loescht und den zweiten Namensraum nicht anfasst; der Snapshot-Zweig bleibt Backup und wird NICHT der Ablageort
- [x] outbox und refsync auf die zweite Art ausdehnen: b.path (core/outbox/outbox.go:77) kollidiert sonst zwischen einem Milestone-Namen und einem Ticket-Handle - je Art ein Unterordner
- [x] Milestone-Schreibweg an den Ref haengen, wie attachRefs es fuer Tickets tut (internal/cli/root.go:255)
- [x] Test: zwei Klone, in einem ein Milestone angelegt, im anderen nach jaira fetch sichtbar - ohne dass ein Zweig gemergt wurde
- [x] core/release/NOTES.md unter ## Unreleased: je eine Zeile fuer den Befehl, den Listen-Schalter, die Board-Geste, die rechte Kartenkante und das Dateiformat
- [-] core/gitref: DeleteMilestone als Gegenstueck zu Delete - refDelete auf MilestoneRefName(name), ErrNoRef heisst nichts zu tun
- [-] core/outbox: OpDelete fuer KindMilestone zulassen - send() (core/outbox/outbox.go:361) weist ihn heute ab, samt Kommentar, der das Gegenteil behauptet
- [-] core/refsync: RecordMilestoneDelete nach dem Vorbild von RecordDelete (core/refsync/refsync.go:133) - Lease aus MilestoneSHA, Weg ueber die Outbox, damit Loeschen offline funktioniert
- [-] core/milestone: Delete(root, name) entfernt die Datei; os.ErrNotExist durchreichen, wie Load es tut
- [-] internal/cli: 'jaira milestone rm' loescht Datei und Ref, wenn die letzte Zeile herausgenommen wird - ein Aufruf, ein Schreibvorgang; Meldung sagt, dass der Milestone weg ist
- [-] Test DoD 8: letztes Ticket herausnehmen - Datei weg, 'milestone ls' nennt ihn nicht, Index leer, also keine Karte traegt seine Farbe
- [-] Test: 'jaira milestone create' ohne Tickets legt eine leere Datei an und sie bleibt liegen - die Regel haengt am Herausnehmen, nicht am Leersein
- [-] Test mit zwei Klonen: nach dem Loeschen holt 'jaira fetch' den Milestone NICHT zurueck (sonst schreibt IncomingMilestones ihn wieder hin)
- [-] internal/cli logbook: 'jaira logbook <name>' erkennt einen Milestone, nachdem die Ticket-Aufloesung nicht greift; er wandert nach .jaira/logbook/<initials>-<datum>/milestones/<name>.md, und sein Ref wird geraeumt
- [-] Weigerung festlegen und bauen: ein Milestone geht nur ins Logbuch, wenn jede seiner Ticket-Zeilen in der Terminal-Lane steht oder schon abgelegt ist - Gegenstueck zu der Regel, die 'jaira logbook <id>' fuer ein Ticket hat
- [-] listLogbook zeigt abgelegte Milestones mit an, sonst ist die Datei da und die Liste sagt es nicht
- [-] Restore: eine Datei aus einem milestones/-Unterordner landet in .jaira/milestones/ statt in TicketsDir (core/ticket/store.go:464) und kommt ueber RecordMilestone wieder auf ihren Ref; Farbe und Mitgliederliste unveraendert
- [-] Test DoD 9: ablegen - vom Board weg und 'milestone ls' nennt ihn nicht; restore - Liste und Farbe zurueck, Karte wieder gefaerbt
- [-] Hilfetexte und docs/COMMANDS.md: das Verschwinden bei 'jaira milestone rm' und das Ablegen bei 'jaira logbook' - beide Stellen liest man VOR dem Aufruf
- [-] DoD 10: je eine Zeile in core/release/NOTES.md unter ## Unreleased fuer das Ablegen und fuer das Verschwinden eines leer geraeumten Milestones
- [x] Arbeitsbaum zurueckdrehen: der komplette Loeschweg aus Runde 1 faellt weg - gitref.DeleteMilestone+refDelete, milestone.Delete, outbox OpDelete/PendingMilestone, refsync.RecordMilestoneDelete und das automatische Wegnehmen in internal/cli/milestones.go; 'git checkout --' auf die fuenf Dateien, denn DoD 8 will das Gegenteil und DoD 10 braucht keinen Ref-Abbau
- [x] Test DoD 8: letztes Ticket mit 'jaira milestone rm' herausnehmen - die Datei bleibt, 'milestone ls' nennt ihn unveraendert; dazu 'create' ohne Tickets, die leere Datei bleibt ebenfalls liegen
- [x] core/milestone: Status aus der Frontmatter lesen ('status: filed'), Filed() dazu, und SetStatus, das die Zeile in der Frontmatter setzt, einfuegt oder entfernt, ohne eine andere Zeile anzufassen - dieselbe verbatim-Regel wie Add/Remove; die Paket-Doku sagt, was 'filed' bedeutet
- [x] Test: SetStatus auf eine von Hand editierte Datei (Kommentar, eigene Reihenfolge, fehlende Frontmatter) - jede andere Zeile unveraendert, Load liest den Status zurueck
- [x] core/ticket/store.go: Ablegen und Zurueckholen fuer den Unterordner milestones/ - eine LogbookMilestone-Haelfte schiebt .jaira/milestones/<name>.md nach .jaira/logbook/<ordner>/milestones/, und Restore (store.go:464) findet eine Datei auch dort und legt sie nach milestone.Dir statt TicketsDir; Mehrdeutigkeit bleibt ein Fehler wie heute
- [x] internal/cli logbook: 'jaira logbook <name>' erkennt einen Milestone, nachdem die Ticket-Aufloesung nicht greift - nur ausdruecklich benannt, nie von '--all' mitgenommen; die Meldung nennt den Ablageort und 'jaira restore'
- [x] Weigerung: ein Milestone geht nur ins Logbuch, wenn jede seiner Ticket-Zeilen in der Terminal-Lane steht oder nicht mehr auf dem Board liegt - Gegenstueck zu der Regel, die 'jaira logbook <id>' fuer ein Ticket hat
- [x] Der Ref bleibt stehen und traegt den Status: beim Ablegen SetStatus('filed') auf den Inhalt, der ueber recordMilestone/RecordMilestone an den Ref geht. Kein Ref-Abbau - ein geraeumter Ref gibt den Namen wieder frei (DoD 11) und der zweite Klon erfaehrt nichts
- [x] core/refsync IncomingMilestones (refsync.go:194): einen Ref, dessen Inhalt 'filed' sagt, NICHT auf die Platte schreiben. Eine schon vorhandene lokale Datei wird dabei nicht geloescht - jaira loescht keine Datei, die es nur gelesen hat
- [x] Test DoD 10 mit zwei Arbeitsbaeumen: im ersten ablegen, im zweiten 'jaira fetch' - der Milestone erscheint dort nicht auf dem Board, und der Ref steht weiter und traegt 'filed'
- [x] internal/cli milestone create: nach der lokalen Pruefung (milestones.go:107) den Ref lesen - traegt er 'filed', wird abgelehnt mit dem Hinweis auf 'jaira restore'. Ohne brauchbare Refs greift derselbe Blick ins lokale Logbuch, sonst legt ein ungeteiltes Board denselben Namen zweimal an
- [x] Restore eines Milestones: die Status-Zeile wieder entfernen und ueber recordMilestone erneut an den Ref - sonst liegt er auf dem Board, waehrend sein Ref weiter 'filed' sagt und jeder andere Klon ihn ausblendet
- [x] listLogbook/logbookNames (internal/cli/logbook.go:224) listen den Unterordner milestones/ mit, sonst liegt die Datei da und die Liste sagt es nicht
- [x] Test DoD 9 und 11: ablegen - vom Board weg, 'milestone ls' schweigt, keine Karte traegt die Farbe, das Logbuch nennt ihn; 'create' mit demselben Namen wird abgelehnt; 'restore' - Liste, Farbe und Karte zurueck
- [x] Hilfetexte und docs/COMMANDS.md: 'jaira logbook' nennt den Milestone-Fall, 'jaira milestone rm' bleibt bei 'der Milestone bleibt stehen', 'create' nennt den belegten Namen - alle drei liest man VOR dem Aufruf
- [x] DoD 12: je eine Zeile in core/release/NOTES.md unter ## Unreleased fuer das Ablegen eines Milestones und fuer den belegten Namen
- [x] editMembers: nach erfolgreichem Load auf ms.Filed() pruefen und mit fail(ExitValidation, "milestone_filed", ...) abweisen - Wortlaut nach milestones.go:122 (Datei, Markierung, 'jaira restore <name>.md' in dem Baum, der abgelegt hat)
- [x] editMembers ErrNotExist-Zweig: vor dem 'create'-Rat milestoneFiled(s, name) fragen, wie create es bei milestones.go:130 tut, und auf 'jaira restore' zeigen statt auf 'create'
- [~] Tests in internal/cli/milestones_test.go: add/rm auf einer markierten Datei auf der Platte, und add/rm auf einem in DIESEN Baum abgelegten Milestone
- [ ] core/release/NOTES.md unter ## Unreleased: eine Zeile fuer die Weigerung von 'jaira milestone add/rm' bei einem abgelegten Milestone

## Progress
- **2026-09-15 14:55 · Alexander Sacharov** — Alex hat am 2026-09-15 aus dem Sprint einen Milestone gemacht. Das ist keine Umbenennung, es aendert die Mechanik - wer dieses Ticket arbeitet, liest ab hier und nicht den Entwurf vom 14.09.

Was ein Milestone ist: eine frei benannte Gruppe, kein Zeitkasten. Er endet nicht an einem Datum, sondern wenn seine Tickets fertig sind. Er ist an keine Release-Version gebunden und kann alles sein, was jemand gruppieren will.

Mehrere laufen gleichzeitig und nebeneinander. Es gibt keinen 'aktuellen' Milestone, auf den sich das Board von allein zusammenzieht.

Damit faellt der teuerste Punkt des alten Entwurfs weg: 'Unerledigte Arbeit wandert in den naechsten Sprint'. Es gibt kein Weiterwandern mehr - unfertige Arbeit bleibt in ihrem Milestone stehen, bis sie fertig ist. Der DoD-Punkt dazu gehoert gestrichen.

Was bleibt und der eigentliche Grund fuer die Datei ist: die Mitgliederliste steht IM Milestone, nicht als Feld an jedem Ticket. Alex' Begruendung am 15.09.: bei einer Aenderung soll man nicht die richtigen Tickets zusammensuchen muessen. Eine Datei bearbeiten statt zwanzig Tickets anfassen - und Tickets liegen auf ihren Refs, also je Ticket ein pull, eine Aenderung, ein release. Genau das macht es sonst niemand.

Die Datei reist auf einem Ref, wie im alten Entwurf. Farbe zufaellig, Markierung rechts auf der Karte, Filter auf den Milestone - alles unveraendert uebernommen.

Offen und Sache der Brainstorm-Lane: ob ein Ticket in mehr als einem Milestone stehen darf, und wenn ja, was die rechte Kante der Karte dann zeigt. Bei parallelen, frei benannten Gruppen ist Mehrfachzugehoerigkeit wahrscheinlich, beim Sprint war sie es nicht.

Was der Milestone NICHT loest, und das war die Ausgangsnot vom 14.09.: 'zeig mir nur die zwanzig, die ich in dieser Sitzung durchgehe'. Ein Milestone mit vierzig Tickets filtert auf vierzig Tickets. Wenn diese Not bleibt, ist sie ein eigenes Ticket und nicht dieses.
- **2026-09-15 14:59 · Alexander Sacharov** — Brainstorm-Lane, gelesen im Code am 2026-09-15. Grundlage ist Alex' Milestone-Entscheidung vom 15.09., nicht der Sprint-Entwurf vom 14.09.

WAS DER CODE ZEIGT

1. Die Kostenrechnung stimmt. Ein Ticket wird geschrieben, indem sein eigener Ref mit compare-and-swap gepusht wird (core/gitref/gitref.go:39 Prefix, :654 pushRefspec mit --force-with-lease). Zwanzig Tickets zu gruppieren heisst also zwanzig Push-Vorgaenge, jeder mit Rennen und Outbox. Eine Datei ist ein Vorgang. Das ist der Grund fuer dieses Ticket, und er haelt.

2. Der Snapshot-Zweig jaira/board ist der FALSCHE Ort, und zwar laut seiner eigenen Dokumentation. core/snapshot/snapshot.go:1-27 sagt woertlich: "It is a backup, not the storage: the working state is always on the refs." Der Zweig wird aus den Refs neu gebaut, im Hintergrund, hoechstens alle 72 Stunden (DefaultEvery, :52), und er loescht dabei Refs gelandeter Tickets. Wer die Mitgliederliste dort ablegt, legt die Wahrheit in einen Cache, der sie beim naechsten Lauf ueberschreibt. Der Punkt aus dem Entwurf vom 14.09. ("Der Snapshot-Zweig ist ein moeglicher Ort") faellt damit weg.

3. Ein Ref traegt heute ausschliesslich Tickets. refs/jaira/tickets/<id>, ein Fetch-Refspec (gitref.go:575), Incoming/Departed, Outbox, die Read-only-Karte mit "pull" im Board (internal/tui/view.go:664). Nichts davon kennt etwas anderes als eine Ticket-ID. Ein zweiter Namensraum ist machbar, aber er ist nicht gratis.

4. Board-weite Dateien gibt es schon, und sie reisen im Zweig, nicht auf einem Ref: .jaira/tags (core/tag/tag.go:32, eine Zeile je Tag, handgeschrieben) und .jaira/lanes. .jaira/lanes bleibt sogar nach 'jaira share' absichtlich gitignoriert (core/board/gitignore.go:18). Es gibt also Praezedenz fuer "kleine Tatsache ueber das Board als Textdatei" - aber keine fuer "board-weite Datei, die auf einem Ref reist".

5. Das Filtern ist heute ein Feld AM Ticket: internal/cli/tickets.go:499 prueft tag.Matches(t.Tags, filter). Steht die Mitgliederliste in der Milestone-Datei, dreht sich die Richtung um: beim Laden muss aus den Dateien ein Index ID -> Milestones gebaut werden, den TUI und CLI gemeinsam nutzen. Das ist eine kleine, aber neue Schicht.

6. Der rechte Kartenrand ist frei, aber nicht umsonst. renderCardBlock (internal/tui/view.go:533) malt je Kartenzeile eine Balkenzelle links und gibt dem Text inner = w-1. Rechts eine Zelle heisst inner = w-2, also eine Spalte weniger Titel. Und: die Spalte muss IMMER reserviert werden, auch bei einem Ticket ohne Milestone - sonst stehen in einer Lane Karten mit unterschiedlich breitem Text nebeneinander und die Titel flattern. "Die Karte wird nicht breiter" ist erfuellbar, "der Titel wird nicht kuerzer" nicht.

DREI WEGE

A) Datei im Baum, unter .jaira/milestones/<name>.md, reist wie .jaira/tags im Zweig.
   Gibt: nichts Neues zu bauen ausser Lesen, Index und Anzeige. Handeditierbar, zeilenweise mergebar, im Diff lesbar, ein Reviewer sieht die Gruppierung im selben Commit wie die Arbeit.
   Kostet: sie kommt erst an, wenn ein Zweig gemergt ist. Genau das hat Alex abgelehnt. Auf einem ungeteilten Board (.jaira/ gitignoriert bis 'jaira share') kommt sie ueberhaupt nicht an.

B) Datei im Baum UND auf einem eigenen Ref, refs/jaira/milestones/<name> - also exakt das Muster, das Tickets schon haben (gitref.go:19-24: dieselben Bytes, einmal im Commit, einmal auf einem Kanal, der ohne Zweig sichtbar ist).
   Gibt: sie ist sofort bei allen, ohne Merge, und sie steht trotzdem im Diff. Kein neues Konzept fuer den Benutzer, es ist das Ticket-Muster ein zweites Mal.
   Kostet: zweiter Ref-Namensraum. Fetch, Prune, Rennen, Outbox, "unsent"-Anzeige und Konfliktbehandlung muessen auf etwas verallgemeinert werden, das keine Ticket-ID ist. Das ist die Haelfte der Arbeit dieses Tickets.

C) Kein neues Objekt: ein Feld milestone am Ticket, wie tags.
   Gibt: fast nichts zu bauen. Filter, Farbregistry und Kartenplaetze existieren; es waere ein zweiter Tag-Typ.
   Kostet: die zwanzig Ref-Schreibvorgaenge, wegen derer das Ticket existiert. Und beim Aendern einer Gruppe muss man die richtigen Tickets erst zusammensuchen - Alex' Begruendung vom 15.09.

WAS ICH TUN WUERDE

B, aber in zwei Schritten, und Schritt eins ist A. Erst die Datei, das Format, der Index, die rechte Kartenkante und der Filter - alles, was den Milestone ueberhaupt erst sichtbar macht, und alles im Baum. Dann der Ref als zweiter Schritt, wenn sich zeigt, dass der Merge-Weg im Alltag zu langsam ist. Grund: der Ref-Namensraum ist die Haelfte des Aufwands und traegt null Nutzen, solange es noch nichts anzuzeigen gibt; und ein Ticket, das beides auf einmal will, ist das Ticket, das nicht fertig wird. C waere billig, aber es loest die Not nicht, die im Ticket steht.

DIE OFFENE FRAGE: MEHRFACHZUGEHOERIGKEIT

Ja, und zwar nicht als Wahl, sondern als Folge. Steht die Mitgliederliste IN der Milestone-Datei, kann niemand verhindern, dass zwei Dateien dieselbe ID nennen: das zu verbieten braeuchte eine Pruefung ueber alle Dateien hinweg, und genau die kann eine handgeschriebene Einzeldatei nicht leisten. Mehrfachzugehoerigkeit ist also erlaubt, weil sie nicht verhinderbar ist.

Die rechte Kante zeigt sie wie die linke die Tags: die Karte hat schon drei Zeilen (cardSlots, internal/tui/view.go:494), rechts entstehen damit drei Plaetze, oben der erste Milestone, darunter der zweite und dritte, in Dateireihenfolge. Derselbe Mechanismus wie cardColors (internal/tui/model.go:1366), gespiegelt. Ein vierter Milestone bleibt am Ticket und faerbt nichts - reine Anzeigegrenze, wie bei S1VM40 entschieden.

ZWEI SACHEN, DIE AM TICKET NOCH FAUL SIND

- DoD-Punkt 6 ("Unerledigte Arbeit wandert in den naechsten Sprint") beschreibt Mechanik, die es nach Alex' Entscheidung nicht mehr gibt. Er gehoert gestrichen, bevor das Ticket in todo liegt - ein Gate am Ende liest ihn und weist die Arbeit sonst zurueck, wenn das Herausfinden am teuersten ist.
- Titel und DoD sagen weiter "Sprint". Das Ziel sagt ab jetzt Milestone. Wer das Ticket arbeitet, liest sonst zwei verschiedene Entwuerfe.
- **2026-09-15 15:02 · Alexander Sacharov** — Dispatcher, 2026-09-15 nach der Brainstorm-Lane: DoD 1-5 von 'Sprint' auf 'Milestone' umgeschrieben, DoD 4 um die Mehrfachzugehoerigkeit (bis zu drei Plaetze rechts, Dateireihenfolge) ergaenzt, DoD 6 (Weiterwandern in den naechsten Sprint) als superseded markiert - Alex' Entscheidung vom 15.09. sagt woertlich, dass dieser Punkt gestrichen gehoert. Der TITEL sagt weiterhin 'Sprint'; jaira hat keinen Befehl zum Umbenennen, das muss von Hand oder per eigenem Ticket passieren. DoD 2 ('reist auf einem Ref') bleibt stehen, obwohl die Brainstorm-Lane B in zwei Schritten empfiehlt und den Ref-Namensraum zurueckstellen will - das ist eine Empfehlung an die Plan-Lane, nicht eine Aenderung an Alex' Entscheidung.
- **2026-09-15 15:05 · Alexander Sacharov** — Plan-Lane, 2026-09-15, im Code nachgelesen. Warum der Plan so aussieht und nicht anders.

FORMAT: .jaira/milestones/<name>.md, Frontmatter fuer name und colour, im Body eine Zeile je Mitglied.
Der DoD sagt woertlich 'von Hand editierbar und im Diff lesbar, wie eine Ticket-Datei' - also Markdown mit Frontmatter und nicht das Zeilenformat von .jaira/tags. Eine Zeile je Mitglied trotzdem, weil eine Zeile das Kleinste ist, was git mergen kann; zwei Leute, die gleichzeitig ein Ticket in denselben Milestone legen, sollen sich nicht in die Quere kommen. Im Body ein Handle je Zeile, nicht die volle ID: wer die Datei von Hand aufmacht, liest Handles, und das Board zeigt ueberall Handles.

FARBE STEHT IN DER MILESTONE-DATEI, nicht in einem zentralen Register. Bei Tags war das Register noetig, weil eine Farbe eine Eigenschaft des TAGS ist und viele Tickets denselben Tag tragen (core/tag/tag.go:1-14). Ein Milestone hat aber schon eine eigene Datei - eine zweite board-weite Datei dafuer waere ein zweites Ding zum Mergen ohne Gegenwert. Vergabe nach dem Muster von Registry.Assign: zufaellig aus der Palette, was noch keiner hat.

DER INDEX DREHT DIE RICHTUNG UM. Heute ist Filtern ein Feld AM Ticket (internal/cli/tickets.go:499, tag.Matches(t.Tags, filter)). Die Mitgliedschaft steht kuenftig in der Milestone-Datei, also muss beim Laden ein Index ID -> Milestones gebaut werden, sonst braeuchte jede Karte einen Scan ueber alle Dateien. Das ist die eine wirklich neue Schicht in diesem Ticket. Vorbild ist link.Build in loadEnv (internal/cli/root.go:275): einmal bauen, dann nur lesen.

DER SNAPSHOT-ZWEIG IST NICHT DER ORT. Der Entwurf vom 14.09. nannte jaira/board als moeglichen Ablageort. Seine eigene Doku widerspricht (core/snapshot/snapshot.go:1-27): 'It is a backup, not the storage: the working state is always on the refs.' Er wird alle 72 Stunden aus den Refs neu gebaut und loescht dabei Refs gelandeter Tickets. Wer die Mitgliederliste dort ablegt, legt sie in einen Cache, der sie ueberschreibt. Darum Schritt 13: nachpruefen, dass reap den neuen Namensraum nicht mitnimmt.

WARUM DER REF TROTZDEM IM PLAN STEHT. Die Brainstorm-Lane hat empfohlen, in zwei Schritten zu bauen und den Ref-Namensraum zurueckzustellen (Weg B, Schritt eins ist A). Das bleibt eine gute Reihenfolge und ist sie hier auch: Schritte 1-11 sind die Datei im Baum, Schritte 12-16 der Ref. Aber weglassen kann der Plan ihn nicht - DoD-Punkt 2 verlangt ihn woertlich, und ein Gate am Ende liest ihn. Wenn Alex den Ref doch zuruecknehmen will, ist das eine Aenderung am DoD und nicht am Plan.

WAS DER REF KOSTET: er ist die Haelfte der Arbeit. Heute traegt ein Ref ausschliesslich Tickets - Prefix, RefName, Fetch-Refspec, List, idsFrom, die Outbox mit gitref.RefName(id), die 'unsent'-Anzeige im Board. Nichts davon kennt etwas anderes als eine Ticket-ID. Die stillste Falle ist die Outbox: b.path ist <id>.json, also kollidiert ein Milestone namens wie ein Handle mit einem Ticket. Darum je Art ein Unterordner (Schritt 14).

RECHTE KARTENKANTE: 'die Karte wird nicht breiter' ist erfuellbar, 'der Titel wird nicht kuerzer' nicht. renderCardBlock gibt dem Text heute inner = w-1 (internal/tui/view.go:533); eine zweite Balkenzelle macht daraus w-2. Und die Spalte muss IMMER reserviert werden, auch bei einem Ticket ohne Milestone - sonst stehen in einer Lane Karten mit unterschiedlich breitem Text nebeneinander. Darum Schritt 10 als eigener Test.

DIE GESTE: kein neuer Filtermechanismus. Das Board hat schon m.filter mit key:value (internal/tui/model.go:598) und die Tag-Box auf 't' (:1010). Ein Milestone-Picker, dessen Auswahl m.filter auf milestone:<name> setzt, ist eine Geste und kein zweiter Weg zum selben Ziel. Der Schluessel wird exakt verglichen, nicht als Teilstring - genau wie tag, sonst widerspricht das Board 'jaira list --milestone'.

NOCH FAUL AM TICKET, ausserhalb dieser Lane: der TITEL sagt weiter 'Sprint'. jaira hat keinen Umbenennen-Befehl. Wer das Ticket aufmacht, liest im Titel den Entwurf vom 14.09. und im Ziel den vom 15.09.
- **2026-09-15 15:10 · Alexander Sacharov** — core/milestone, 2026-09-15. Was das Paket nicht tut und warum: keine eigene Palette - Palette = tag.Palette, weil Tag links und Milestone rechts auf der Karte stehen und nie verwechselt werden koennen; sechzehn weitere Werte wuerden nur naeher an die Statusfarben (39/214/203/78/141) ruecken. Mitglied ist die VOLLE ULID, nicht das Handle: die Datei ist auch der Merge-Gegenstand, und ein Handle ist nicht garantiert eindeutig. parseMember akzeptiert nur eine gueltige ULID, damit ein gewoehnlicher Markdown-Bullet in der Prosa nicht als Mitglied gelesen wird - das ist der Grund, warum die Datei ueberhaupt Prosa enthalten darf. Load/Save haelt die Zeilen verbatim wie core/tag, Frontmatter wird NICHT ueber core/ticket geparst: die Ticket-Frontmatter kennt Schema und Pflichtfelder, ein Milestone hat drei Zeilen.
- **2026-09-15 15:15 · Alexander Sacharov** — TUI, 2026-09-15. Taste ist M, nicht m - m ist 'move' und muss das ueberall bleiben. Die Geste schreibt in m.filter ('milestone:<name>') statt eine zweite Verengung daneben zu halten: damit raeumt esc auf dem Board sie genauso weg wie einen getippten Filter, und / zeigt, worauf das Board verengt ist. Im Picker loest x den Filter, weil sich niemand merkt, dass esc auf dem BOARD das tut. matches() hat jetzt einen dritten Parameter (milestone.Index) statt einer zweiten Funktion matchesIn - zwei Namen fuer eine Frage driften. inner in renderCardBlock ist w-2: die rechte Zelle ist IMMER reserviert, auch ohne Milestone, sonst wandern die Titel einer Lane um eine Spalte, wenn ein Ticket einer Gruppe beitritt. Getestet in internal/tui/milestonebar_test.go.
- **2026-09-15 15:27 · Alexander Sacharov** — Refs, 2026-09-15. gitref wurde NICHT auf ein generisches (Art, Name) umgebaut, wie Plan-Schritt 12 es woertlich sagt: das haette jede Signatur und jede Aufrufstelle in refsync, outbox, internal/cli/refs.go und internal/tui/refs.go angefasst. Stattdessen Root = 'refs/jaira/', Prefix = Root+'tickets/', MilestonePrefix = Root+'milestones/', und die gemeinsamen Teile als private refSHA/refRead/refWrite/refDelete/listNames. Die Ticket-API ist unveraendert, die Milestone-API steht daneben (MilestoneSHA/ReadMilestone/WriteMilestone/DeleteMilestone/ListMilestones). Fetch holt Root/* statt Prefix/* - ein Fetch bringt beide Arten, sonst zeigt ein Board eine Gruppe, deren Datei niemand hat. snapshot.reap braucht keine Aenderung: Run() nimmt Repo.List(), und das liest nur Prefix - abgesichert durch TestTheTwoNamespacesStayApart, das absichtlich Ticket 01TEST und Milestone 01TEST nebeneinanderlegt. outbox: Unterordner je Art (tickets/, milestones/); der ALTE flache Pfad wird weiter GELESEN, damit ein von einem aelteren Build eingestellter Schreibvorgang beim Upgrade nicht verlorengeht. MilestoneSender ist ein zweites Interface, nicht eine Erweiterung von Sender - ein alter Sender kompiliert weiter. Nicht gemacht und bewusst: 'jaira show' zeigt die Milestones eines Tickets nicht; das gehoert in ein eigenes Ticket.
- **2026-09-15 15:32 · Alexander Sacharov** — critique: Der Bau ist im Kern richtig - zwei Ref-Namensraeume statt eines generischen (Art, Name), die Datei als Format, der Index als einzige neue Schicht, die rechte Kante immer reserviert. Nichts davon muss anders. Neun Findings, alle im Kleinen, und sie fallen in drei Gruppen.

1) SPEKULATIVES: der ganze Loeschweg fuer Milestones ist gebaut und wird von nichts aufgerufen. refsync.go:197 RecordMilestoneDelete hat keinen Aufrufer, outbox.go:280 PendingMilestone keinen, gitref.DeleteMilestone nur einen Test - und die OpDelete-Verzweigung in outbox.send erreicht niemand, weil kein Befehl einen Milestone loescht (die Hilfe sagt woertlich: das ist 'rm' auf einer Datei). Entweder 'jaira milestone delete' anbieten und verdrahten, oder die drei weg. Dasselbe fuer gitref ListRemoteMilestones.

2) FEHLERBEHANDLUNG FUER EINEN UNMOEGLICHEN ZUSTAND: MilestoneSender (outbox.go:282) plus der Typ-Assert bei :364 schuetzt vor einem Sender, der keine Milestones traegt. Box.Flush hat GENAU EINE Aufrufstelle, refsync.go:344, und die uebergibt *gitref.Repo, das beide Haelften erfuellt. Kein anderer Typ im Baum erfuellt Sender. Die zwei Methoden gehoeren an Sender selbst; dann fallen der Assert, die Fehlermeldung und die untypische (error, bool)-Rueckgabe von send() gemeinsam weg.

3) ZWEITE NAMEN FUER VORHANDENE DINGE: validColour (milestones.go:360) ist tag.ValidColour abgeschrieben; itoa/atoi (milestone.go:358) sind strconv, das core/tag fuer genau dieses Feld schon benutzt; Milestone.path wird geschrieben und nie gelesen, weil alle milestone.Path() rechnen; activeMilestones() ist m.milestones unter zweitem Namen, und der Picker indiziert beide Listen mit einem Cursor.

EINS DAVON IST KEIN AUFRAEUMEN, SONDERN FORMAT: milestone.go:129 schreibt 'colour:' in die Frontmatter. internal/cli/tags.go:233 schreibt die Regel dieses Boards woertlich hin - 'color, matching --color: one spelling on the machine surface' - und milestones.go:347 sowie der Flag halten sie schon. Die Datei IST die API (CLAUDE.md); nach dem Release ist das ein Bruch, jetzt ist es eine Zeile.

NICHT AUFGEMACHT, absichtlich: dass IncomingMilestones die lokale Datei ueberschreibt, dass die Palette mit tag geteilt wird, dass der Index in newListCmd statt in loadEnv gebaut wird, und dass gitref nicht auf (Art, Name) verallgemeinert wurde statt Plan-Schritt 12 - alle vier sind in den Notizen begruendet und die Begruendung traegt.
- **2026-09-15 15:39 · Alexander Sacharov** — critique-Runde 1 abgearbeitet, 2026-09-15. Alle neun Findings behoben; zwei Entscheidungen, die der Code nicht erklaert.

LOESCHWEG: entfernt statt verdrahtet (refsync.RecordMilestoneDelete, outbox.PendingMilestone, gitref.DeleteMilestone, gitref.ListRemoteMilestones, der OpDelete-Zweig fuer Milestones in outbox.send und die zwei Tests darauf). Grund: kein Befehl loescht einen Milestone, der Plan sah keinen vor, und Scope-Disziplin schlaegt Symmetrie. FOLGE, bewusst offen gelassen: wer .jaira/milestones/<name>.md von Hand loescht, bekommt die Datei beim naechsten 'jaira fetch' zurueck, weil der Ref stehen bleibt. Das ist eine echte Falle und braucht 'jaira milestone delete <name>', das Datei, Ref und Outbox-Eintrag zusammen wegnimmt - eigenes Ticket, nicht hier.

SENDER: WriteMilestone sitzt jetzt direkt auf outbox.Sender, MilestoneSender ist weg. Nachgesehen: Box.Flush hat genau eine Aufrufstelle (refsync.go:344) und die uebergibt *gitref.Repo. Damit fallen Typ-Assert, Fehlermeldung und die (error, bool)-Rueckgabe von send() weg; ein unbekannter Op kommt jetzt als gewoehnlicher Fehler zurueck und wird wie vorher zu Failed. fakeSender im outbox-Test hat WriteMilestone dazubekommen.

FORMAT: die Frontmatter-Zeile heisst 'color:', nicht 'colour:' - internal/cli/tags.go:233 schreibt diese Regel fuer die Maschinenoberflaeche hin. Das Go-Feld heisst weiter Colour, wie in core/tag; die Regel gilt fuer das, was in der Datei steht. Mitgeaendert: core/milestone/milestone_test.go, core/gitref/milestone_test.go und die NOTES-Zeile zum Dateiformat. Jetzt eine Zeile, nach dem Release ein Bruch.

Plan-Schritt 5 nachgetragen abgehakt: der Index wird gebaut (internal/cli/milestones.go:29 milestoneIndex, internal/tui/model.go:374), nur in newListCmd statt in loadEnv - die critique hat genau das als begruendet durchgehen lassen.
- **2026-09-15 15:44 · Alexander Sacharov** — critique-Runde 2, 2026-09-15. Die Findings aus Runde 1 sind alle abgearbeitet und werden nicht wieder aufgemacht: der Loeschweg ist weg, MilestoneSender ist in Sender aufgeloest, die Frontmatter sagt color:, validColour/itoa/atoi/path/activeMilestones sind ersetzt, der fetch-Hilfetext stimmt. Der Bau bleibt im Kern richtig.

Sechs neue Findings, zwei davon nicht kosmetisch.

1) IDENTITAET DES MILESTONE. core/milestone/milestone.go:145: Load setzt m.Name nur dann aus dem Dateinamen, wenn die Frontmatter keinen hat; Save schreibt nach Path(root, m.Name). Die Datei ist laut CLAUDE.md und laut dem eigenen Hilfetext von Hand editierbar - wer also die Zeile name: aendert, legt beim naechsten jaira milestone add eine zweite Datei an, waehrend die alte mitsamt ihrem Ref liegen bleibt. Zwei Wahrheiten fuer einen Namen. Entweder der Dateiname gewinnt immer (Load setzt m.Name = name), oder die Zeile name: verschwindet aus New - sie wird sonst nirgends gebraucht.

2) DIE ALTE OUTBOX WIRD GELESEN, NICHT MIGRIERT. core/outbox/outbox.go:196: List liest sowohl das flache Verzeichnis als auch tickets/, und QueueKind loescht die flache Datei nach dem Superseden nicht. Nachgemessen mit einem Wegwerf-Test: nach einem Queue auf ein altes Ticket enthaelt List zwei Eintraege derselben ID, der veraltete zuerst. Erst wird also der alte Inhalt gesendet, danach traegt der neue Eintrag ein Lease, das der Remote nicht mehr hat. Ein paralleler Lesepfad loest das nicht - QueueKind soll legacyPath nach dem Schreiben entfernen (den Key hat es), oder List auf (Kind, ID) deduplizieren und den Unterordner gewinnen lassen.

3) RESTE VON RUNDE 1. core/gitref/gitref.go:364 refDelete und :659 listRemoteNames haben je einen Aufrufer, dem sie eine Konstante uebergeben - die zweiten Aufrufer sind mit dem Loeschweg gegangen. Beide inlinen.

4) DREIMAL DIESELBE REGEL. Farbe 0 bedeutet keine Farbe steht in internal/cli/milestones.go:292, internal/tui/model.go milestoneColors und internal/tui/view.go:1531 - und milestones.go:119 laesst --color 0 durch, was dann ueberall als farblos erscheint. Ein HasColour() am Milestone, an allen drei Stellen benutzt, und am Flag eine Entscheidung.

5) core/outbox/outbox.go:253 DropKind vergleicht Pfad-Strings, um zu wissen, in welchem Durchlauf es steht, und ist die einzige Methode, die kind nicht normalisiert.

6) internal/tui/view.go:1501 und :1532 sind derselbe Swatch-Ausdruck.

NICHT AUFGEMACHT: dass IncomingMilestones die lokale Datei ueberschreibt, die geteilte Palette, der Index in newListCmd, gitref ohne generisches (Art, Name) - alle vier stehen seit Runde 1 mit ihrer Begruendung. Ebenso, dass milestoneIndex() genau einen Aufrufer hat: das ist die Naht, an der die naechsten Leser haengen, und sie ist als solche beschrieben.
- **2026-09-15 15:45 · Alexander Sacharov** — in-progress, 2026-09-15, Runde 2 der critique-Findings. Wahl bei Finding 1 (zwei Wahrheiten fuer den Namen): der DATEINAME gewinnt immer, und die Zeile name: verschwindet aus New. Die Alternative - name: gewinnt und Save benennt die Datei um - waere ein Rename samt Ref-Umzug bei jedem Save gewesen, also genau der Vorgang, den dieses Ticket abschaffen will.
- **2026-09-15 15:48 · Alexander Sacharov** — in-progress, Runde 2 abgearbeitet, 2026-09-15. Was der Code nicht erklaert:
- Finding 2 war real und ist jetzt durch TestQueueSupersedesTheEntryAnOlderBuildLeft (core/outbox/outbox_test.go) abgesichert - nachgemessen: ohne das os.Remove in QueueKind liefert List zwei Eintraege derselben ID, der veraltete zuerst. Der PARALLELE LESEPFAD bleibt, nur der Schreibvorgang raeumt auf: einen Eintrag, den dieser Build nie wieder anfasst, darf er auch nicht loeschen.
- Finding 1: 'name:' ist aus New UND aus parse() raus, nicht nur aus New. Wuerde parse den Schluessel weiterlesen, saehe ein Leser wieder zwei Wahrheiten, obwohl nur eine gilt. Alte Dateien mit der Zeile bleiben lesbar - die Zeile ist dann Prosa in der Frontmatter und ueberlebt Save verbatim, wie jede andere unbekannte Zeile.
- Finding 4: HasColour() sitzt in core/milestone, nicht im TUI, weil die Regel 'Farbe 0 heisst keine Farbe' eine Eigenschaft des Dateiformats ist und nicht der Darstellung. --color 0 wird jetzt am Flag abgewiesen (internal/cli/milestones.go), damit die Regel nicht erst auf der Karte sichtbar wird.
- Nicht angefasst, absichtlich: tag.ValidColour bleibt 0-255. Das Tag-Registry hat kein 'keine Farbe' in dieser Form; eine gemeinsame Verschaerfung waere eine Aenderung an S1VM40s Gegenstand.
- **2026-09-15 15:53 · Alexander Sacharov** — critique-Runde 3, 2026-09-15. Die sechs Findings aus Runde 2 sind alle abgearbeitet und werden nicht wieder aufgemacht: der Dateiname ist die einzige Identität (Load setzt m.Name immer, name: ist aus New und parse raus), QueueKind räumt den flachen Alt-Eintrag weg und TestQueueSupersedesTheEntryAnOlderBuildLeft misst es nach, HasColour() steht an allen drei Stellen und --color 0 wird abgewiesen, refDelete/listRemoteNames sind inline, DropKind normalisiert kind, swatch() ist ein Helfer. Der Bau selbst bleibt richtig.

Drei Findings in Runde 3, und alle drei sind dieselbe Sorte: die Oberfläche, die jemand VOR dem Aufruf liest, beschreibt noch den Stand vor Runde 2.

1) internal/cli/milestones.go:51 verspricht ein name: in der Frontmatter, das es nicht mehr gibt. Das ist nicht kosmetisch: die Datei IST die API (CLAUDE.md), der Hilfetext ist die einzige Stelle, an der ein Mensch das Format erklärt bekommt, und wer nach dieser Erklärung ein name: hinschreibt, bekommt keine Fehlermeldung, sondern eine Zeile, die nichts tut.

2) internal/cli/milestones.go:78 und der Flag-Text bei :146 sagen 0-255, die Prüfung bei :123 sagt 1-255. Die Fehlermeldung ist die einzige richtige der drei, und sie sieht man erst, nachdem man dem Hilfetext geglaubt hat. NOTES.md:18 sagt bereits 1-255 - die Hilfe ist jetzt die letzte Stelle, die widerspricht.

3) docs/COMMANDS.md kennt milestone überhaupt nicht - weder die vier Unterbefehle noch --milestone in der list-Zeile 57. README.md:676 nennt diese Datei die vollständige Referenz, und jede andere Befehlsfamilie steht dort, jaira tag und jaira tags eingeschlossen (Zeilen 130-131). Eine Befehlsfamilie, die nur die eingebaute Hilfe kennt, findet niemand, der nicht schon weiß, dass es sie gibt.

NICHT aufgemacht, weil begründet und die Begründung trägt: die geteilte Palette mit tag, der Index in newListCmd statt loadEnv, dass IncomingMilestones die lokale Datei überschreibt, der fehlende Löschweg (eigenes Ticket), und dass gitref zwei Namensräume nebeneinander hat statt eines generischen (Art, Name). Auch nicht aufgemacht: dass das Ticket im TITEL weiter 'Sprint' sagt - dafür fehlt jaira ein Umbenennen-Befehl, das ist kein Finding an diesem Diff.
- **2026-09-15 15:54 · Alexander Sacharov** — Dispatcher, 2026-09-15: ANGEHALTEN nach der dritten Ruecksendung aus der critique-Lane. Die Regel des Dispatchers ist, nach drei Ruecksendungen desselben Lanes an einen Menschen zu uebergeben, und ich fuehre keine vierte Runde.

Das Ticket liegt in in-progress. Die drei offenen Findings aus Runde 3 sind alle Dokumentations-Nachzug aus den Runden 1-3, nicht neue Mechanik: der Hilfetext von 'jaira milestone' nennt noch 'name:' in der Frontmatter, das es seit Runde 2 nicht mehr gibt; --color sagt an zwei Stellen 0-255 und an einer 1-255; docs/COMMANDS.md hat noch keine Zeile fuer 'jaira milestone' und kein --milestone bei 'jaira list'.

Was ich als Beobachtung dazuschreibe und worauf ich ausdruecklich NICHT handle: die Findings werden kleiner (9, dann 6, dann 3), keines wurde wiederholt, und Runde 3 findet nur, was Runde 1 und 2 verursacht haben. Von innen sieht das nach einer Schleife aus, die konvergiert. Genau so sieht aber auch eine Schleife aus, die nie endet, weil jede Runde tiefer liest. Welche von beiden das hier ist, kann nur ein Mensch entscheiden.

Zum Weitermachen: eine vierte in-progress-Runde mit genau diesen drei Findings ist ein kleiner, klar begrenzter Auftrag. Danach stehen noch optimize, testing und review aus.
- **2026-09-15 18:09 · Alexander Sacharov** — Alex hat am 2026-09-15 die vierte in-progress-Runde freigegeben, nachdem der Dispatcher nach drei critique-Ruecksendungen angehalten hatte. Danach laufen optimize, testing und review wie vorgesehen.

Ausserdem entschieden, und damit ist der Punkt zu: die rechte Kartenkante bleibt, wie sie gebaut ist - eine Zelle je Milestone, von oben, hoechstens drei. KEINE Volleinfaerbung aller drei Zellen bei nur einem Milestone. Grund: die Karte soll zeigen, in wie vielen Gruppen ein Ticket steht, und das geht nur, wenn eine Zelle eine Gruppe ist. Wer das spaeter wieder aufmachen will, liest hier, dass es nicht vergessen, sondern entschieden wurde.

Auf dem Board liegen drei Milestones demo-ui, demo-board-dateien und demo-naechste-version, absichtlich mit ueberlappenden Tickets (0YGWXQ in allen dreien, DBJTKQ in zweien). Sie sind eine Vorfuehrung, keine echte Gruppierung, und sie liegen auch auf upstream. Wer sie wegraeumt: es gibt noch keinen Befehl dafuer, das ist Ticket N71NVG.
- **2026-09-15 18:11 · Alexander Sacharov** — in-progress, Runde 4 (critique-Runde 3), 2026-09-15. Alle drei Findings waren Text, kein Verhalten - der Code selbst ist unveraendert, go vet und go test ./... bleiben gruen.

Was der Code nicht erklaert:
- Der lange Hilfetext von 'jaira milestone' (internal/cli/milestones.go:42-48) sagt weiterhin, dass unerledigte Arbeit in die Datei des naechsten Milestones wandert. Das ist KEIN Rueckfall auf den Sprint-Entwurf vom 14.09.: gemeint ist das Editieren von Hand, das der Nutzer selbst macht, nicht eine Mechanik, die jaira ausfuehrt. Dieselbe Formulierung steht in core/release/NOTES.md:18. Wer sie streichen will, streicht sie an beiden Stellen und nicht nur an einer.
- --color steht jetzt an vier Stellen als 1-255: Langtext (:80), Flag-Beschreibung (:148), die Abweisung (:125) und COMMANDS.md. Drei davon sind Text, und genau die liest man VOR dem Aufruf - deshalb war das ein Finding und keine Kleinigkeit.
- docs/COMMANDS.md: die vier milestone-Zeilen stehen in der Writing-Tabelle, auch 'ls'. In der Looking-Tabelle waere es sauberer sortiert, aber dann stuenden die vier Befehle einer Familie in zwei Tabellen - 'jaira tags' steht auch bei Looking und 'jaira tag' bei Writing, und genau das sucht man beim Lesen zweimal.

Nicht angefasst, ausserhalb der Lane: .jaira/milestones/demo-board-dateien.md, demo-naechste-version.md und demo-ui.md liegen als untracked Dateien im Worktree - Rueckstand eines Demolaufs von Hand. Sie sind NICHT committet worden; wer aufraeumt, loescht sie einfach.
- **2026-09-15 18:15 · Alexander Sacharov** — critique-Runde 4, 2026-09-15: KEIN Finding. Die drei Findings aus Runde 3 sind nachgemessen behoben - internal/cli/milestones.go:50-54 sagt jetzt, dass der Dateiname der Name ist und die Frontmatter color und created-at traegt; grep ueber 0-255/1-255 findet --color nur noch als 1-255 an :80, :125, :148 und in docs/COMMANDS.md:132, die einzige verbliebene 0-255-Stelle ist 'jaira tag', das seit Runde 2 begruendet unangetastet bleibt; docs/COMMANDS.md:57 hat --milestone und :132-135 die vier Zeilen fuer create/add/rm/ls.

Was ich darueber hinaus gelesen habe, statt nur den Zuwachs: core/milestone/milestone.go ganz, internal/cli/milestones.go ganz, der TUI-Diff (model.go/view.go) und der Ref-Weg (outbox.go, refsync.go, fetch.go, tickets.go). Nichts Neues. Jedes exportierte Stueck von core/milestone hat einen Aufrufer - nachgesehen fuer Names, For, Matches, HasColour, Members, NormalizeName, Build, LoadAll, AssignColour, Add, Remove - also keine Abstraktion ohne Nutzer mehr, was in Runde 1 und 2 die haeufigste Sorte war.

Ausdruecklich NICHT aufgemacht, weil in Runde 1 bis 3 stehengelassen und die Begruendung traegt: dass IncomingMilestones die lokale Datei ueberschreibt, die mit tag geteilte Palette, der Index in newListCmd statt loadEnv, gitref mit zwei Namensraeumen statt einem generischen (Art, Name), der fehlende Loeschweg (Ticket N71NVG), und der Titel, der weiter 'Sprint' sagt.

Die Schleife hat konvergiert: 9, 6, 3, 0. Das ist der vorgesehene Ausgang der Lane und keine Nachsicht - eine vierte Runde haette nur wiederholt, was schon beantwortet ist.
- **2026-09-15 18:21 · Alexander Sacharov** — optimize: staticcheck U1000 ueber core/milestone, core/outbox, core/gitref, internal/cli und internal/tui findet nichts, was dieser Branch verwaist hat - die drei Treffer (share.go isShared, model.go laneStart, model.go currentLane) stehen schon auf master. staticcheck ist im Repo nicht installiert, lief ueber 'go run honnef.co/go/tools/cmd/staticcheck@latest'.
- **2026-09-15 18:21 · Alexander Sacharov** — optimize: .jaira/milestones/demo-*.md sind Handtest-Dateien aus einer frueheren Lane und bleiben untracked - sie gehoeren nicht in den Commit, aber jemand sollte sie am Ende loeschen.
- **2026-09-15 18:28 · Alexander Sacharov** — testing: go build/vet ./... RC=0, go clean -testcache && go test ./... -race green in all 28 Pakete (internal/tui 123s). DoD 1-5 und 7: jede Proof-Zeile im Baum nachgeschlagen, die Tests laufen einzeln gruen. DoD 6 war ungehakt und ohne Proof - auf einem Scratch-Board mit dem echten Binary nachgestellt: sprint-1 mit drei Tickets, zwei Zeilen von Hand in sprint-2.md verschoben, danach 'jaira list --milestone' 1 bzw. 2 Tickets; Kommentar und Leerzeilen der Datei ueberleben spaetere 'milestone add/rm'. Ref-Weg mit zwei echten Clones geprueft: refs/jaira/milestones/sprint-x steht auf dem Remote, im zweiten Clone die Datei geloescht, 'jaira fetch' holt sie zurueck und nennt 'Milestones updated from the remote: sprint-x'. Nebenbei: --color 0 und 256 werden mit RC=2 abgelehnt, Umbenennen per mv der Datei wirkt sofort, --milestone nope gibt 'No tickets match.' RC=0. Nicht mitcommittet: .jaira/milestones/demo-*.md im Worktree sind Demo-Dateien aus einer frueheren Lane und bleiben untracked.
- **2026-09-15 20:14 · Alexander Sacharov** — Alex hat die Arbeit am 2026-09-15 in der human-Lane geprueft und drei Dinge verlangt. Nachgesehen, was davon schon da ist, damit die naechste Runde nicht zweimal baut:

FERTIG - der Ref. core/refsync/refsync.go:170 RecordMilestone stellt die Bytes der Milestone-Datei in die Outbox, :194 IncomingMilestones schreibt jede Milestone-Datei, die die Refs tragen, auf die Platte und legt .jaira/milestones/ dabei an. Das ist DoD 2 und war schon abgehakt. Hier ist nichts zu tun.

FEHLT - das Ablegen. Es gibt keine Verbindung zwischen Milestone und Logbuch: kein Treffer auf 'milestone' in core/logbook oder core/archive.

FEHLT - das Verschwinden. 'jaira milestone' hat add, create, ls, rm. 'rm' nimmt TICKETS AUS einem Milestone heraus; es gibt kein Kommando und keinen Pfad, der den Milestone selbst entfernt. Ein leer geraeumter Milestone bleibt also als Datei, als Ref und als Farbe auf dem Board stehen.

Daraus sind DoD 8-10 geworden.

Meine Lesart von Alex' Satz, bevor jemand anders sie anders liest: ein Milestone soll denselben Lebenslauf haben wie ein Ticket - das Logbuch legt ihn ab, restore holt ihn zurueck. Wenn Alex etwas anderes meinte (etwa: das Logbuch nimmt beim Ablegen eines Tickets dessen Zeile aus dem Milestone), ist DoD 9 falsch formuliert und gehoert korrigiert statt gebaut.

Offen und Sache der Plan-Lane, weil es den Bau entscheidet: verschwindet der leere Milestone VON ALLEIN, sobald die letzte Zeile herausgenommen wird, oder braucht es dafuer ein Kommando? Automatisch ist bequemer und laesst sich nicht vergessen; es loescht aber eine frisch mit 'jaira milestone create' angelegte, noch leere Datei sofort wieder - und genau so legt man einen Milestone an, bevor man weiss, was hineinkommt. Wer das automatisch baut, braucht eine Antwort darauf.

Ref-Transport beim Loeschen nicht vergessen: eine Datei von der Platte zu nehmen raeumt refs/jaira/ nicht. Wer den Milestone entfernt, muss auch seinen Ref raeumen, sonst schreibt IncomingMilestones ihn beim naechsten Zug wieder hin.
- **2026-09-15 20:19 · Alexander Sacharov** — pre-process, Runde 2 (DoD 8-10), 2026-09-15. Die offene Frage aus der human-Lane ist entschieden, und zwar so:

VERSCHWINDEN IST EIN EREIGNIS, KEIN ZUSTAND. Ein Milestone verschwindet in dem Moment, in dem seine LETZTE Ticket-Zeile herausgenommen wird - nicht deshalb, weil seine Liste leer ist. Damit faellt der Einwand aus der human-Lane weg: 'jaira milestone create <name>' ohne Tickets legt eine leere Datei an, und die bleibt liegen, weil niemand etwas herausgenommen hat. Genau so legt man einen Milestone an, bevor man weiss, was hineinkommt.

Die Alternative waere ein eigenes Kommando ('jaira milestone delete'). Dagegen spricht, dass DoD 8 woertlich das Herausnehmen des letzten Tickets als Nachstellung nennt, und dass ein Milestone, den man vergisst zu loeschen, als Farbe auf Karten weiterlebt. Ein automatisches Loeschen laeuft NUR im Schreibweg (milestone rm), nie beim Lesen: eine von Hand leer editierte Datei wird nicht beim naechsten Laden geloescht - jaira loescht keine Datei, die es nur gelesen hat.

ZWEI AUSGAENGE, ZWEI BEDEUTUNGEN, das ist kein Widerspruch zwischen DoD 8 und 9: leer geraeumt heisst aufgegeben und die Gruppe ist spurlos weg; fertig heisst abgerechnet und die Gruppe wandert samt Liste ins Logbuch, wo sie nachlesbar bleibt.

WAS DER CODE DAFUER NOCH NICHT HAT:
- core/outbox/outbox.go:361 send() lehnt OpDelete fuer KindMilestone ausdruecklich ab ('A milestone is only ever written, never deleted'). Dieser Satz ist ab jetzt falsch und der Weg muss gebaut werden.
- core/gitref hat Delete nur fuer Tickets; MilestoneRefName ist da, das Gegenstueck zu Delete fehlt.
- core/ticket/store.go:464 Restore legt JEDE Datei nach TicketsDir. Ein abgelegter Milestone braucht ein eigenes Ziel.

REF RAEUMEN IST PFLICHT, nicht Kosmetik: refsync.IncomingMilestones (refsync.go:194) schreibt jede Milestone-Datei, die die Refs tragen, auf die Platte. Wer nur die Datei loescht, bekommt sie beim naechsten 'jaira fetch' zurueck. Das gilt fuer das Verschwinden UND fuer das Ablegen ins Logbuch - anders als bei einem Ticket, wo RecordFiled den Ref absichtlich stehen laesst und erst der Snapshot-Lauf ihn raeumt; fuer Milestones raeumt der Snapshot nichts.

ABLAGEORT: .jaira/logbook/<initials>-<datum>/milestones/<name>.md, ein Unterordner. Grund: ein Milestone-Dateiname (frei gewaehlt) kann mit einem Ticket-Dateinamen kollidieren, und Restore muss am Fundort erkennen, wohin die Datei zurueckgehoert - Ordner statt Namensraten.

NOCH IMMER NICHT BESTAETIGT, und wer es anders liest, korrigiert DoD 9 statt es zu bauen: dass mit Alex' Forderung der Lebenslauf des Milestones selbst gemeint ist (Logbuch legt den Milestone ab) und nicht, dass das Ablegen eines TICKETS dessen Zeile aus dem Milestone nimmt. Gebaut wird die erste Lesart, weil DoD 9 sie woertlich so sagt.
- **2026-09-15 20:22 · Alexander Sacharov** — Alex hat am 2026-09-15 die offene Frage aus der letzten Notiz beantwortet. DoD 8-12 tragen das Ergebnis; hier steht die Begruendung, damit spaeter niemand zurueckbaut.

Kein Automatismus. Ein Milestone geht nur auf Kommando weg. Die automatische Loeschung eines leer geraeumten Milestones ist ausdruecklich NICHT gewollt - sie wuerde eine frisch angelegte, noch leere Datei sofort wieder wegnehmen, und genau so legt man einen Milestone an. DoD 8 ist deshalb umgedreht: es prueft jetzt, dass der leere Milestone STEHEN BLEIBT.

Ablegen wie ein Ticket, bestaetigt. 'jaira logbook' nimmt den Milestone vom Board, 'jaira restore' holt ihn zurueck.

Der eigentliche Fund dieser Runde, und er kommt von Alex: der Ref wird beim Ablegen NICHT geraeumt. Stattdessen traegt der Milestone einen Status.

Warum das die richtige Loesung ist - es raeumt genau die Falle weg, die in der Notiz vom 15.09. schon stand: eine Datei von der Platte zu nehmen raeumt refs/jaira/ nicht, und IncomingMilestones (core/refsync/refsync.go:194) schreibt sie beim naechsten Zug wieder hin. Ein abgelegter Milestone kaeme also von allein zurueck, und zwar auf jedem Rechner, der zieht. Mit einem Status im Frontmatter ist der Ref die Wahrheit statt der Gegner: er sagt 'abgelegt', und der Client schreibt die Datei nicht mehr aufs Board.

Der Ref haelt ausserdem den Namen belegt. Das ist DoD 11 und Alex' zweite Sorge in einem Satz ('sonst kann es sich wiederholen'): wer nach dem Ablegen denselben Namen noch einmal anlegt, bekommt zwei Milestones mit einer Identitaet, die auf verschiedenen Rechnern verschieden aussehen. Statt dessen wird das Anlegen abgelehnt, mit dem Hinweis auf 'jaira restore'.

Damit ist die Plan-Lane nicht mehr blockiert. Was sie noch selbst entscheidet: wie der Status heisst und wo er steht (Frontmatter-Feld der Milestone-Datei ist der naheliegende Ort), und ob 'jaira logbook' den Milestone von sich aus anfasst oder ob er ausdruecklich benannt werden muss.
- **2026-09-15 20:24 · Alexander Sacharov** — Dispatcher, 2026-09-15 20:2x: die in-progress-Runde 1 wurde ABGEBROCHEN, nicht fertig. Grund: Alex hat mitten in der Runde DoD 8 umgedreht und DoD 10-12 dazugelegt, und der Worker baute genau das Gegenteil.

Was er gebaut hatte, unfertig und NICHT committet, im Arbeitsbaum liegend (core/gitref/gitref.go, core/milestone/milestone.go, core/outbox/outbox.go, core/refsync/refsync.go, internal/cli/milestones.go, zusammen ~120 Zeilen): den kompletten LOESCHWEG. gitref.DeleteMilestone plus ein herausgezogenes refDelete, milestone.Delete auf der Datei, outbox.PendingMilestone, OpDelete fuer KindMilestone im Sendeweg, und in milestones.go das automatische Wegnehmen eines leer geraeumten Milestones.

Dieser Weg ist ab jetzt FALSCH, in beiden Haelften:
- DoD 8 sagt jetzt, ein leer geraeumter Milestone BLEIBT STEHEN. Das automatische Loeschen in milestones.go ist damit nicht mehr gewollt.
- DoD 10 sagt, der Ref eines abgelegten Milestones wird NICHT geraeumt, sondern traegt den Status. DeleteMilestone/refDelete/OpDelete loesen also ein Problem, das es nicht mehr gibt.

Der Diff liegt als Patch unter dem Scratchpad dieser Dispatcher-Sitzung (inprogress-round1-deletepath.patch), falls jemand eine Zeile daraus doch braucht. Der Arbeitsbaum wurde ABSICHTLICH nicht zurueckgesetzt: was davon stehenbleibt, entscheidet die Plan-Lane und nicht der Dispatcher.

Der Plan-Schritt 'core/gitref: DeleteMilestone als Gegenstueck zu Delete' steht auf [~] und ist der erste, den die neue Plan-Runde anfassen muss.
- **2026-09-15 20:29 · Alexander Sacharov** — pre-process, Runde 3 (DoD 8-12 nach Alex' Umkehrung), 2026-09-15. Die Plan-Schritte 18-32 sind SUPERSEDED, nicht erledigt: sie beschrieben den Loeschweg, und den gibt es nicht mehr. Die Schritte 33-48 sind der Ersatz und stehen in der Reihenfolge, in der sie gebaut werden.

DER GANZE DIFF IM ARBEITSBAUM GEHT WEG, alle fuenf Dateien. Nachgesehen, ob eine Zeile davon ueberlebt: keine. gitref.DeleteMilestone, milestone.Delete, outbox OpDelete fuer KindMilestone, outbox.PendingMilestone, refsync.RecordMilestoneDelete und das automatische Wegnehmen in milestones.go loesen alle dasselbe Problem - einen Milestone-Ref abbauen - und DoD 10 sagt ausdruecklich, dass er stehen bleibt. Auch das herausgezogene refDelete geht zurueck: es haette dann wieder nur einen Aufrufer, und genau das war Finding 3 der critique-Runde 2. Der Patch liegt im Scratchpad der Dispatcher-Sitzung, falls doch jemand eine Zeile sucht.

WIE DER STATUS AUSSIEHT, das war die offene Frage der Plan-Lane: eine Zeile 'status: filed' in der Frontmatter der Milestone-Datei, neben color und created-at. Kein eigenes Registry und keine zweite Datei, weil die Datei laut CLAUDE.md die API ist und ein Mensch in einem Diff lesen koennen muss, warum ein Milestone nicht mehr auf dem Board steht. Englisch, weil die Dateien und der Code englisch sind - im Ticket heisst es 'abgelegt', in der Datei 'filed'. Fehlt die Zeile, steht der Milestone auf dem Board; das ist der Zustand jeder heute existierenden Datei, also braucht es keine Migration.

Setzen und Entfernen brauchen ein SetStatus in core/milestone, kein Neuschreiben der Datei: Save schreibt m.lines verbatim zurueck, und wer die Frontmatter neu baut, verliert Kommentare und Reihenfolge - genau das, was DoD 1 abgesichert hat.

WARUM DER REF DIE WAHRHEIT IST UND NICHT DER GEGNER: IncomingMilestones (core/refsync/refsync.go:194) schreibt jede Milestone-Datei, die die Refs tragen, auf die Platte. Ein abgelegter Milestone mit abgebautem Ref kaeme nicht zurueck - aber sein Name waere auch wieder frei, und dann legen zwei Rechner zwei Milestones mit einer Identitaet an. Mit dem Status im Ref-Inhalt ist beides in einem Zug geloest: der Leser sieht 'filed' und schreibt ihn nicht aufs Board (Schritt 41), und 'milestone create' sieht denselben Ref und lehnt den Namen ab (Schritt 43).

WAS ICH ENTSCHIEDEN HABE, WEIL DIE DoD ES OFFEN LAESST:
- 'jaira logbook <name>' fasst einen Milestone nur an, wenn er BENANNT ist. '--all' sweept die Terminal-Lane und bleibt bei Tickets. Grund: --all ist der Schnitt, den jemand macht, wenn er Stunden eintraegt, und dieser Schnitt darf keine Gruppe mitnehmen, an die niemand gedacht hat - dieselbe Sorge, die im Hilfetext von logbook schon mit den neunundvierzig fremden Tickets steht.
- Reihenfolge der Aufloesung: erst Ticket, dann Milestone. Ein Milestone-Name ist frei gewaehlt und koennte wie ein Handle aussehen; das Ticket gewinnt, weil es die haeufigere und die aeltere Bedeutung ist.
- Ein Milestone, der auf einem anderen Rechner abgelegt wird, verschwindet dort NICHT automatisch vom Board: der Fetch schreibt ihn nur nicht mehr hin. Eine lokale Datei zu loeschen, weil ein Ref das sagt, waere das erste Mal, dass jaira eine Datei entfernt, die es nur gelesen hat. DoD 10 verlangt es auch nicht - dort steht 'schreibt ihn NICHT wieder aufs Board'. Wer dort aufraeumt, legt selber ab.
- Der Ref-Blick in 'create' reicht auf einem ungeteilten Board nicht (dort gibt es keine Refs), deshalb zusaetzlich der Blick ins lokale Logbuch. Zwei Quellen, aber die Suchfunktion braucht Restore ohnehin.

NICHT ENTSCHIEDEN UND AUCH NICHT NOETIG: ob ein abgelegter Milestone jemals wieder aus dem Logbuch verschwindet. 'jaira restore' holt ihn zurueck, sonst liegt er da - genau wie ein abgelegtes Ticket.
- **2026-09-15 20:36 · Alexander Sacharov** — in-progress Runde 3, 2026-09-15. Was beim Bauen anders kam als im Plan.

REIHENFOLGE BEIM ABLEGEN ist keine Geschmacksfrage: SetStatus -> Save -> recordMilestone -> LogbookMilestone. recordMilestone (internal/cli/milestones.go) liest die Datei von der PLATTE, nicht aus dem Objekt. Wer zuerst verschiebt, schickt eine leere Datei an den Ref oder gar nichts - und dann steht der Ref ohne Status und jeder andere Klon schreibt den Milestone zurueck aufs Board.

KEIN IMPORTZYKLUS MOEGLICH: core/milestone importiert core/ticket, also kann core/ticket/store.go den Unterordnernamen nicht von dort holen. Loesung: ticket.MilestonesSubdir ist die Quelle, milestone.Subdir = ticket.MilestonesSubdir. Eine Schreibweise, in der Richtung, die der Compiler erlaubt.

RESTORE MERKT SICH DEN FUNDORT, nicht den Namen: store.Restore sammelt jetzt neben logMatches auch logDirs und leitet das Ziel daraus ab (milestones/ -> MilestonesDir, sonst TicketsDir). Am Namen zu erkennen, ob eine Datei ein Milestone ist, geht nicht - ein Milestone-Name ist frei gewaehlt.

STOLPERSTELLE im Hilfetext: cobra Long ist ein Backtick-String. Ein Backtick-Zitat wie 'status: filed' darin beendet das Literal und der Compiler zeigt auf eine Stelle 200 Zeilen weiter unten. Anfuehrungszeichen benutzen.
- **2026-09-15 20:43 · Alexander Sacharov** — in-progress Runde 3 fertig, 2026-09-15. Was beim Bauen entschieden wurde und im Plan nicht stand.

DAS ABLEGEN BRAUCHT EIN GATE, ANALOG ZUM TICKET: logbookMilestone laedt jedes Mitglied und verlangt Terminal-Lane ODER 'nicht mehr auf dem Board' (s.Load schlaegt fehl = abgelegt/archiviert). Ein Ticket, das nicht geladen werden kann, gilt als fertig genug - sonst kann ein Milestone, dessen Tickets schon abgelegt sind, nie ins Logbuch.

RESTORE ERKENNT DEN MILESTONE AM ZIELPFAD, nicht am Argument: unfileMilestone (internal/cli/archive.go) prueft filepath.Base(filepath.Dir(dst)) == 'milestones'. Der Aufrufer nennt nur den Dateinamen, und ob das ein Ticket oder ein Milestone ist, weiss erst store.Restore, nachdem es die Datei gefunden hat.

MILESTONE-NAME BELEGT: zwei Quellen, Ref zuerst, dann lokales Logbuch (milestoneFiled in internal/cli/milestones.go). Der Ref allein reicht nicht - ein ungeteiltes Board hat keine Refs, und dann wuerde derselbe Name zweimal vergeben. Das Logbuch allein reicht auch nicht - wer nie gefetcht hat, sieht die fremde Ablage nicht.

NICHT GEBAUT UND ABSICHTLICH: ein anderswo abgelegter Milestone verschwindet beim Fetch NICHT vom eigenen Board. IncomingMilestones ueberspringt den Ref nur. Eine lokale Datei zu loeschen, weil ein Ref das sagt, waere das erste Mal, dass jaira eine Datei entfernt, die es nur gelesen hat; DoD 10 verlangt es nicht.

ADJACENT, NICHT ANGEFASST: 'jaira restore' hat keinen Hinweis darauf, dass der Name eines Milestones ohne .md-Endung nicht funktioniert - man muss 'round-one.md' schreiben. Das ist bei Tickets genauso und waere eine eigene Aenderung.
- **2026-09-15 20:46 · Alexander Sacharov** — critique-Runde 5 (DoD 8-12), 2026-09-15. Der Bau ist im Kern richtig: Milestone verlaesst das Board nur auf Befehl, der Ref bleibt stehen und traegt den Status, restore holt die Datei dorthin zurueck, wo sie gefunden wurde. Fuenf Findings, und die ersten drei haengen zusammen.

1. 'status: filed' ist eine zweite Wahrheit, die das Board nicht liest. core/milestone/milestone.go:LoadAll filtert Filed() NICHT - nachgesehen, Filed() wird nur an drei Stellen gelesen (refsync.go:224, archive.go:143, milestones.go:342). Dass ein abgelegter Milestone vom Board verschwindet, haengt allein daran, dass die Datei verschoben wurde. Die Zeile soll das entscheiden, nicht der Ort der Datei.

2. Daraus folgt die Ordnung in internal/cli/logbook.go:logbookMilestone - Save, recordMilestone, dann Rename - deren Kommentar ausdruecklich sagt, sie sei keine Vorliebe. Schlaegt das Rename fehl (dst existiert schon), steht die Datei markiert weiter auf dem Board, waehrend jeder andere Clone sie versteckt. Mit einem Filter in LoadAll ist dieser Halbzustand harmlos.

3. core/refsync/refsync.go:224 ueberspringt einen abgelegten Ref ganz. Ein Clone, der die Datei schon hat, erfaehrt damit NIE, dass abgelegt wurde - genau das, wofuer der Ref da ist. Die Begruendung im Kommentar ('jaira loescht keine Datei, die es nur gelesen hat') traegt die Alternative nicht: den markierten Inhalt normal schreiben loescht nichts und braucht den Sonderfall gar nicht.

4. internal/cli/archive.go:143 schluckt den Save-Fehler und meldet trotzdem 'Restored'.

5. internal/cli/milestones.go:milestoneFiled laeuft die Logbook-Ordner von Hand ab und sieht den Alt-Ordner nicht, den core/ticket/store.go:logbookFolders abdeckt - und aus dem Restore einen Milestone sehr wohl zurueckholen kann. Dritter handgeschriebener Walk neben logbookNames und Restore.

Kein Finding aus den Runden 1-4 wird wieder aufgemacht; die betrafen den Bau vor DoD 8-12 und sind nachgemessen erledigt.
- **2026-09-15 20:49 · Alexander Sacharov** — in-progress Runde 4 (critique-Runde 5), 2026-09-15, angefangen. Die fuenf Findings werden in dieser Reihenfolge gebaut: LoadAll filtert Filed() (1), damit faellt die Begruendung der Reihenfolge in logbookMilestone weg (2), IncomingMilestones schreibt den markierten Inhalt normal (3), unfileMilestone gibt den Save-Fehler zurueck (4), milestoneFiled ruft eine Store-Suche auf logbookFolders statt eines dritten Handwalks (5).
- **2026-09-15 20:58 · Alexander Sacharov** — in-progress Runde 4 fertig (critique-Runde 5), 2026-09-15. Alle fuenf Findings behoben. Was der Code nicht erklaert:

FINDING 3 HAT EINE FALLE, die im Finding nicht steht: 'den markierten Inhalt normal schreiben' darf NICHT heissen, ihn auch dort hinzuschreiben, wo es die Datei gar nicht gibt. Sonst legt der Fetch auf dem ablegenden Klon selbst die Datei wieder unter .jaira/milestones/ an, waehrend die Kopie im Logbuch liegt - und 'jaira restore' verweigert danach mit 'is already on the board'. IncomingMilestones schreibt einen 'filed'-Ref deshalb nur ueber eine Datei, die schon da ist (haveErr == nil). Genau das ist auch der Zweck des Findings: ein Klon, DER DIE DATEI HAT, soll erfahren, dass abgelegt wurde.

ZWEI LISTEN AUS IncomingMilestones (wrote, filed) statt einer: 'Milestones updated from the remote' und daneben 'Milestones filed elsewhere, now off this board'. Einen abgelegten Milestone unter 'updated' zu melden schickt den Leser zu 'jaira milestone ls', das ihn nicht mehr nennt. JSON hat dafuer milestones_filed.

FALLOUT VON FINDING 1, den das Finding nicht nennt: sobald LoadAll Filed() filtert, kann eine markierte Datei legitim auf der Platte liegen. Zwei Wege brauchen deshalb eine eigene Abweisung - 'jaira logbook <name>' legt einen schon markierten Milestone nicht ein zweites Mal ab, und 'jaira milestone create' sagt bei einer markierten lokalen Datei nicht mehr 'already exists' (das schickt zu einem Listing, das ihn nicht nennt), sondern nennt die Zeile und den Weg zurueck: Zeile von Hand entfernen oder der Ableger macht 'jaira restore'.

NICHT GETESTET, absichtlich: der Fehlerweg von Finding 4 (Save schlaegt fehl, nachdem Restore die Datei zurueckgeschoben hat). Ihn herbeizufuehren hiesse das Verzeichnis schreibgeschuetzt zu machen, und daran scheitert schon das Rename davor - der Test wuerde etwas anderes messen als er behauptet.

FINDING 5: die Suche heisst jetzt ticket.Store.FiledMilestone und liegt auf logbookFolders, deckt also den Alt-Ordner .jaira/sync/ mit ab - aus dem Restore einen Milestone sehr wohl zurueckholen kann. TestFiledMilestoneFindsBothLogbookFolders misst beide Ordner nach.
- **2026-09-15 21:01 · Alexander Sacharov** — critique round 6, two findings, both about what the user is told when a milestone is filed somewhere else.

1. internal/cli/logbook.go:283 — the 'already filed' refusal points at 'jaira restore <name>.md'. Every state that reaches this refusal is a tree WITHOUT a logbook copy: a fetch wrote the marked file back (refsync writes it only where the file already exists), or a filing marked and did not move, or a restore moved the file back and failed to unmark it. In all three ticket.Store.Restore answers 'is not in the archive or in .jaira/logbook/'. Point at the tree that filed it instead — milestones.go:120 already words this correctly.

2. internal/cli/milestones.go:121 and core/release/NOTES.md:26 recommend deleting the 'status: filed' line by hand. That is the duplicate the refsync doc comment (core/refsync/refsync.go:206-213) refuses to create: the hand edit rides out on the ref with the next milestone command, the filer's tree fetches the unmarked file onto its board, and its own 'jaira restore' then hits the 'is already on the board' guard in core/ticket/store.go:544 with the logbook copy left stranded. The file staying hand-editable is not the same as the hand edit being the documented way back.

Checked and left alone: LoadAll's filed filter and its single-caller effects (tui/model.go:369, milestones.go:30/136/302) are consistent; Store.FiledMilestone swallowing logbookFolders' error matches the os.ReadDir it replaced; IncomingMilestones has one caller. The modify/delete conflict between a clone that commits the marked file and the filer who git-mv'd it is real but is not new — the previous skip behaviour produced the same conflict shape — so it is not re-raised here.
- **2026-09-15 21:04 · Alexander Sacharov** — in-progress round 5 (critique round 6), 2026-09-15. Both findings were message text; no behaviour changed and the suite stays green.
- Finding 1 (logbook.go): the refusal now names the file and its mark and says the restore has to run in the tree that filed it. Checked before rewording: every route here — a fetch writing a marked file back, a filing that marked and did not move, a restore that moved back and failed to unmark — leaves this tree WITHOUT a logbook copy, so the old advice hit 'is not in the archive or in .jaira/logbook/'.
- Finding 2 (milestones.go, NOTES.md, COMMANDS.md): the hand edit of 'status: filed' is no longer offered as the way back. It is still possible — the file stays hand-editable, that is the format's promise — but it is not documented, because doing it in a second clone puts the milestone back on that board alone while the logbook copy stays stranded in the filer's tree, and that tree's own 'jaira restore' then hits 'is already on the board' (core/ticket/store.go:544). Hand-editable is not the same as the documented route.
- Not changed: internal/cli/fetch.go's long text and docs/COMMANDS.md:146/147 describe the mark without recommending removing it, so they were already right.
- **2026-09-15 21:07 · Alexander Sacharov** — critique round 7, one finding. Both findings from round 6 are text and are measured fixed — logbook.go:283 now names the file, its mark and the tree that holds the copy, and neither milestones.go:122 nor core/release/NOTES.md:26 offers the hand edit any more; neither is re-opened.

FINDING: internal/cli/milestones.go:237 — 'jaira milestone add/rm' is the third door into a milestone file and the only one without a Filed() refusal. editMembers calls milestone.Load, which reads a marked file happily, edits it, saves it, records it to the ref and reports '<name> now holds N ticket(s)' — for a milestone LoadAll hides, 'milestone ls' will not name and no card paints. Reachable through the two states LoadAll's own doc names: a fetch writing a filed ref over a file this tree has (core/refsync/refsync.go:IncomingMilestones), and a filing that marked and did not move. Refuse it the way create does at milestones.go:121, wording modelled on :122.

Nothing is lost when it happens — the status line survives the write, so the ref stays filed and no other clone is misled. What is wrong is only what the user is told, which is why this is one finding and not a redesign.

CHECKED AND LEFT ALONE: unfileMilestone recording an unmarked milestone that was never marked is a no-op write, not a bug; AssignColour no longer sees a filed milestone's colour since LoadAll filters, but a filed milestone paints nothing, so a reused colour cannot collide on screen; milestoneFiled's ref-before-logbook order is right for an unshared board; fetch.go's two lists and their JSON keys match what IncomingMilestones returns.
- **2026-09-15 21:08 · Alexander Sacharov** — Dispatcher, 2026-09-15 21:1x: HALT NACH DREI RUECKLAEUFEN. Die critique-Lane hat die Arbeit dreimal zurueckgeschickt (Runde 5: fuenf Findings, Runde 6: zwei, Runde 7: eines). Die Regel des Dispatchers sagt, dass hier ein Mensch entscheidet, und der Dispatcher startet deshalb KEINE vierte Runde.

Was ein Mensch zum Entscheiden braucht - die Beobachtung, nicht die Empfehlung:
- Die Findings werden kleiner und keines wird wieder aufgemacht. Runde 6 und 7 sagen das ausdruecklich und nennen jedes Mal, was nachgemessen und liegengelassen wurde.
- Runde 5 betraf den Bau ('status: filed' entschied nicht, was auf dem Board steht). Runde 6 und 7 betreffen nur noch, was der Benutzer zu lesen bekommt, wenn er auf eine markierte Datei trifft.
- Das offene Finding aus Runde 7: 'jaira milestone add/rm' ist die dritte Tuer in eine markierte Milestone-Datei und die einzige ohne Weigerung; create (milestones.go:121) und logbook (logbook.go:283) haben je eine. Es geht dabei nichts verloren - die Status-Zeile ueberlebt den Schreibvorgang, der Ref bleibt 'filed' - falsch ist nur die Meldung an den Benutzer.

Die zwei Lesarten, zwischen denen nur ein Mensch entscheiden kann, und sie sehen von hier gleich aus: entweder ist die Definition of Done unvollstaendig (jede Tuer in eine abgelegte Datei braucht eine Weigerung, und das gehoert hingeschrieben), oder die Schleife konvergiert, ohne je zu enden, weil jede Runde eine Ebene tiefer liest und tiefer immer geht.

Das Ticket steht in in-progress. Der Zweig feat/0YGWXQ traegt alles bis 65aff0a; DoD 8-12 sind gebaut, die Suite mit -race ist gruen.
- **2026-09-16 06:51 · Alexander Sacharov** — Alex hat am 2026-09-16 entschieden: der Befund aus critique-Runde 7 wird behoben, nicht mitgeliefert. Damit ist die Schleife wieder offen - sie war nach drei Ruecklaeufen angehalten worden, weil die Regel das zur Entscheidung eines Menschen macht, und der Mensch hat sie getroffen.

Zu bauen sind DREI Dinge. Punkt 2 stand in keiner critique-Runde; er ist beim Nachlesen des Codes fuer Alex' Frage 'wie willst du das fixen' aufgefallen.

PUNKT 1 - der fehlende Waechter, der urspruengliche Befund.
internal/cli/milestones.go:237 editMembers laedt mit milestone.Load und behandelt danach nur os.ErrNotExist. ms.Filed() wird nie gefragt, obwohl es die Methode gibt (core/milestone/milestone.go:309). Also schreibt 'jaira milestone add/rm' in eine Datei, die 'status: filed' traegt, und meldet freundlich, wie viele Tickets sie jetzt haelt.

Was den Befund schwerer macht als er im review-Feld steht: unmittelbar nach ms.Save steht recordMilestone(s, ms). Die Aenderung bleibt nicht lokal - sie geht auf den Ref. Ein abgelegter Milestone wird damit bei JEDEM wieder lebendig, der zieht. Das ist kein stiller Erfolg, das ist eine Veroeffentlichung.

Zu bauen: nach dem erfolgreichen Load, vor der Schleife, auf ms.Filed() pruefen und mit fail(ExitValidation, "milestone_filed", ...) abweisen. Den Wortlaut nicht neu erfinden - milestones.go:122 ist das Vorbild und nennt drei Dinge: die Datei, ihre Markierung, und dass 'jaira restore <name>.md' in dem Baum, der abgelegt hat, sie zurueckholt.

PUNKT 2 - der falsche Rat im ErrNotExist-Zweig. NEU, in keiner Runde genannt.
Ist der Milestone in DIESEM Baum abgelegt worden, liegt seine Datei nicht mehr unter .jaira/milestones/. Load gibt ErrNotExist, und milestones.go:242 antwortet 'no milestone %q on this board - jaira milestone create %s starts it'. Wer diesem Rat folgt, laeuft in den naechsten Abweisungstext ('filed into the logbook'). Zwei Schritte fuer eine Auskunft, und der erste schickt in die falsche Richtung.

create macht es bereits richtig und hat ZWEI Abweisungen: die markierte Datei auf der Platte (milestones.go:116) UND milestoneFiled(s, name), das im Logbuch-Ordner nachsieht (milestones.go:130). add/rm muss dieselbe zweite Frage stellen und sofort auf 'jaira restore' zeigen statt auf 'create'.

PUNKT 3 - Tests. internal/cli/milestones_test.go hat keinen einzigen add/rm-Fall auf einer markierten Datei. Zwei Faelle: markierte Datei auf der Platte, und in das Logbuch dieses Baumes abgelegt.

Warum das kein Randfall ist, fuer den Fall dass jemand kuerzen will: core/refsync/refsync.go:194 IncomingMilestones schreibt einen abgelegten Ref ueber eine Datei, die dieser Baum schon hat. Die markierte Datei auf der Platte entsteht also dadurch, dass man das Ablegen eines anderen zieht - der normale Weg, nicht ein Sonderfall.

Keine neue DoD-Zeile: Punkt 1 und 2 sind beide das, was DoD 9 und 11 schon verlangen ('ist vom Board weg', 'der Name ist belegt'). Eine NOTES.md-Zeile nur, wenn der Abweisungstext etwas ist, das ein Benutzer von aussen merkt - das tut er hier, also wahrscheinlich ja.
