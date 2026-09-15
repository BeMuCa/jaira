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
commits: []
created-at: 2026-09-14T17:49:30Z
updated-at: 2026-09-15T15:13:13Z
assignee: "Alexander Sacharov"
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-79843
claimed-at: 2026-09-15T15:02:47Z
outcome-what: "Der Plan steht: 17 Schritte, Datei im Baum zuerst (1-11), Ref-Namensraum danach (12-16), NOTES.md zuletzt."
outcome-why: "Die Plan-Lane hat im Code nachgesehen, was der Milestone kostet, und die offenen Punkte des Entwurfs entschieden: Format, Ablageort, Farbvergabe, Indexrichtung."
outcome-resolves: "Format .jaira/milestones/<name>.md festgelegt, Snapshot-Zweig als Ablageort verworfen, Ref bleibt im Plan weil DoD 2 ihn verlangt."
---

# Ein Sprint ist eine eigene Datei, keine Markierung am einzelnen Ticket

## Definition of Done

- [ ] Ein Milestone laesst sich anlegen und benennen, und seine Datei fuehrt die Tickets auf, die zu ihm gehoeren. Sie ist von Hand editierbar und im Diff lesbar, wie eine Ticket-Datei.
- [ ] Die Milestone-Datei reist auf einem Ref: wer sie zieht, sieht denselben Milestone, ohne auf das Mergen eines Zweiges zu warten.
- [ ] Jeder Milestone bekommt seine Farbe, ohne dass jemand eine aussucht.
- [ ] Eine Karte zeigt die Farbe ihres Milestones am RECHTEN Rand, deutlich getrennt von den Tag-Farben am linken; ein Ticket ohne Milestone zeigt dort nichts und die Karte wird dadurch nicht breiter. Mehrfachzugehoerigkeit ist erlaubt: rechts stehen bis zu drei Plaetze, in Dateireihenfolge.
- [ ] Das Board zieht sich mit einer Geste auf einen Milestone zusammen, und 'jaira list' hat den entsprechenden Schalter.
- [-] Unerledigte Arbeit wandert in den naechsten Sprint, indem eine Datei bearbeitet wird - nicht indem jedes Ticket einzeln angefasst wird. Nachgestellt mit mindestens drei Tickets, von denen zwei weiterwandern.
- [ ] Je eine Zeile in core/release/NOTES.md unter ## Unreleased fuer das, was ein Benutzer davon merkt.

## Options

- [x] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Format festlegen und als Paket-Doku in core/milestone hinschreiben: .jaira/milestones/<name>.md, Frontmatter name/colour/created-at, Body eine Zeile je Mitglied - wie eine Ticket-Datei, und eine Zeile ist das Kleinste, was git mergen kann
- [x] core/milestone: Load/LoadAll/Save nach dem Vorbild von core/tag (core/tag/tag.go:189 Load, :371 Save) - Zeilen verbatim erhalten, WriteAtomic, Mitgliederliste parsen
- [x] Farbvergabe: zufaellig aus einer Palette, die keine Farbe eines schon vorhandenen Milestones doppelt; die Farbe steht in der Milestone-Datei selbst, kein zweites zentrales Registry
- [x] Test: eine Milestone-Datei von Hand editieren (Kommentar, Leerzeile, eigene Reihenfolge) und nach Load/Save unveraendert wiederfinden
- [ ] Index ID -> Milestones einmal beim Laden bauen, neben link.Build in loadEnv (internal/cli/root.go:275) und im TUI (internal/tui/model.go:352); CLI und TUI lesen denselben Index
- [ ] CLI: jaira milestone create/add/rm/ls - je Aufruf genau ein Dateischreibvorgang, damit zwanzig Tickets gruppieren eine Bearbeitung bleibt
- [ ] CLI: jaira list --milestone <name> (neben --tag, internal/cli/tickets.go:528) und der Schluessel milestone:<name> in beiden matches() (internal/cli/tickets.go:534, internal/tui/model.go:598) - exakt wie tag, nicht als Teilstring
- [x] TUI: rechte Kartenkante in renderCardBlock (internal/tui/view.go:533) - zweite Balkenzelle rechts, inner = w-2, die Spalte IMMER reserviert, auch ohne Milestone, sonst flattern die Titel in einer Lane
- [x] TUI: milestoneColors gespiegelt zu cardColors (internal/tui/model.go:1366) - cardSlots Plaetze, in Dateireihenfolge, ein vierter Milestone faerbt nichts
- [~] Test: Karte mit 0, 1 und 4 Milestones - gleiche Kartenbreite, gleiche Textbreite, hoechstens drei gefaerbte Plaetze rechts
- [ ] TUI: die Geste - Picker wie die Tag-Box auf 't' (internal/tui/model.go:1010), die Auswahl setzt m.filter auf milestone:<name> und nutzt damit den vorhandenen Filterweg
- [ ] gitref: Namensraum von 'Ticket-ID' auf '(Art, Name)' verallgemeinern - Prefix (core/gitref/gitref.go:39), RefName (:144), Fetch-Refspec (:575), List/ListRemote/idsFrom (:587-625); refs/jaira/milestones/<name> neben refs/jaira/tickets/<id>
- [ ] snapshot: pruefen, dass reap (core/snapshot/snapshot.go:234) nur Ticket-Refs loescht und den zweiten Namensraum nicht anfasst; der Snapshot-Zweig bleibt Backup und wird NICHT der Ablageort
- [ ] outbox und refsync auf die zweite Art ausdehnen: b.path (core/outbox/outbox.go:77) kollidiert sonst zwischen einem Milestone-Namen und einem Ticket-Handle - je Art ein Unterordner
- [ ] Milestone-Schreibweg an den Ref haengen, wie attachRefs es fuer Tickets tut (internal/cli/root.go:255)
- [ ] Test: zwei Klone, in einem ein Milestone angelegt, im anderen nach jaira fetch sichtbar - ohne dass ein Zweig gemergt wurde
- [ ] core/release/NOTES.md unter ## Unreleased: je eine Zeile fuer den Befehl, den Listen-Schalter, die Board-Geste, die rechte Kartenkante und das Dateiformat

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
