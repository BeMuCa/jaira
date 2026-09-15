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
created-at: 2026-09-14T17:49:30Z
updated-at: 2026-09-15T15:54:05Z
assignee: "Alexander Sacharov"
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-28259
claimed-at: 2026-09-15T15:33:55Z
outcome-what: "Alle sechs Findings der critique-Runde 2 behoben: Dateiname ist der einzige Milestone-Name (name: raus aus New und parse), QueueKind loescht den superseded flachen Outbox-Eintrag, HasColour() als einzige Stelle fuer 'Farbe 0 heisst keine Farbe' plus Abweisung von --color 0, refDelete/listRemoteNames inline, DropKind normalisiert kind statt Pfade zu vergleichen, ein swatch-Helfer statt zwei Ausdruecken."
outcome-why: "Finding 1 und 2 waren echte Fehler: ein von Hand geaenderter name: legte beim naechsten add eine zweite Datei an, und ein von einem aelteren Build hinterlassener Outbox-Eintrag wurde neben dem neuen gesendet - erst der veraltete Inhalt, dann ein Lease, das der Remote nicht mehr hat. Nachgemessen mit TestQueueSupersedesTheEntryAnOlderBuildLeft, der ohne den Fix zwei Eintraege derselben ID sieht."
outcome-resolves: "Kein DoD-Punkt aendert sich - die sechs Findings waren Korrektheit und Doppelung innerhalb der schon gebauten Mechanik. go vet und go test ./... sind gruen."
review-summary: |-
  internal/cli/milestones.go:51 der Hilfetext von 'jaira milestone' sagt weiter "Frontmatter carries the name, the colour and when it was created" - seit Runde 2 schreibt New() kein name: mehr und parse() liest keines; wer das liest und von Hand ein name: einträgt, ändert nichts und merkt es nie. Ersetzen durch: der Dateiname IST der Name, die Frontmatter trägt color und created-at - genau wie core/milestone/milestone.go:154 und core/release/NOTES.md:18 es schon sagen.
  internal/cli/milestones.go:78 und :146 nennen den Bereich weiter "--color <0-255>" bzw. "ANSI-256 colour (0-255)", während :123 jetzt "1-255" fordert und 0 zurückweist; drei Stellen, zwei davon falsch, und die falschen sind die, die man vor dem Aufruf liest. Beide auf 1-255 ändern, mit dem Grund in einem Halbsatz (0 färbt keine Zelle).
  docs/COMMANDS.md:57 listet für 'jaira list' weiter nur --lane/--assignee/--tag/--query/--actionable, und die Befehlstabelle hat keine Zeile für milestone create/add/rm/ls - obwohl README.md:676 diese Datei als vollständige Referenz ausweist und jede andere Befehlsfamilie (jaira tag, jaira tags, jaira lanes ...) dort steht. Vier Zeilen in die Writing-Tabelle neben 'jaira tag' und --milestone in Zeile 57 nachtragen.
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
- [-] Unerledigte Arbeit wandert in den naechsten Sprint, indem eine Datei bearbeitet wird - nicht indem jedes Ticket einzeln angefasst wird. Nachgestellt mit mindestens drei Tickets, von denen zwei weiterwandern.
- [x] Je eine Zeile in core/release/NOTES.md unter ## Unreleased fuer das, was ein Benutzer davon merkt.
  proof: core/release/NOTES.md:17

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
