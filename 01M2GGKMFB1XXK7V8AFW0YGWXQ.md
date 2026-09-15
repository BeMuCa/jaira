---
id: 01M2GGKMFB1XXK7V8AFW0YGWXQ
title: "Ein Sprint ist eine eigene Datei, keine Markierung am einzelnen Ticket"
status: backlog
ready: true
creator: Alexander Sacharov
goal: "Wer plant, legt einen Sprint als eigene Datei an, sieht auf jeder Karte am rechten Rand welche Tickets zusammengehoeren, und filtert das Board mit einem Griff darauf."
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
updated-at: 2026-09-15T14:56:43Z
assignee: "Alexander Sacharov"
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-69299
claimed-at: 2026-09-15T14:56:26Z
---

# Ein Sprint ist eine eigene Datei, keine Markierung am einzelnen Ticket

## Definition of Done

- [ ] Ein Sprint laesst sich anlegen und benennen, und seine Datei fuehrt die Tickets auf, die zu ihm gehoeren. Sie ist von Hand editierbar und im Diff lesbar, wie eine Ticket-Datei.
- [ ] Die Sprint-Datei reist auf einem Ref: wer sie zieht, sieht denselben Sprint, ohne auf das Mergen eines Zweiges zu warten.
- [ ] Jeder Sprint bekommt seine Farbe, ohne dass jemand eine aussucht.
- [ ] Eine Karte zeigt die Farbe ihres Sprints am RECHTEN Rand, deutlich getrennt von den Tag-Farben am linken; ein Ticket ohne Sprint zeigt dort nichts und die Karte wird dadurch nicht breiter.
- [ ] Das Board zieht sich mit einer Geste auf einen Sprint zusammen, und 'jaira list' hat den entsprechenden Schalter.
- [ ] Unerledigte Arbeit wandert in den naechsten Sprint, indem eine Datei bearbeitet wird - nicht indem jedes Ticket einzeln angefasst wird. Nachgestellt mit mindestens drei Tickets, von denen zwei weiterwandern.
- [ ] Je eine Zeile in core/release/NOTES.md unter ## Unreleased fuer das, was ein Benutzer davon merkt.

## Options

- [x] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-15 14:55 · Alexander Sacharov** — Alex hat am 2026-09-15 aus dem Sprint einen Milestone gemacht. Das ist keine Umbenennung, es aendert die Mechanik - wer dieses Ticket arbeitet, liest ab hier und nicht den Entwurf vom 14.09.

Was ein Milestone ist: eine frei benannte Gruppe, kein Zeitkasten. Er endet nicht an einem Datum, sondern wenn seine Tickets fertig sind. Er ist an keine Release-Version gebunden und kann alles sein, was jemand gruppieren will.

Mehrere laufen gleichzeitig und nebeneinander. Es gibt keinen 'aktuellen' Milestone, auf den sich das Board von allein zusammenzieht.

Damit faellt der teuerste Punkt des alten Entwurfs weg: 'Unerledigte Arbeit wandert in den naechsten Sprint'. Es gibt kein Weiterwandern mehr - unfertige Arbeit bleibt in ihrem Milestone stehen, bis sie fertig ist. Der DoD-Punkt dazu gehoert gestrichen.

Was bleibt und der eigentliche Grund fuer die Datei ist: die Mitgliederliste steht IM Milestone, nicht als Feld an jedem Ticket. Alex' Begruendung am 15.09.: bei einer Aenderung soll man nicht die richtigen Tickets zusammensuchen muessen. Eine Datei bearbeiten statt zwanzig Tickets anfassen - und Tickets liegen auf ihren Refs, also je Ticket ein pull, eine Aenderung, ein release. Genau das macht es sonst niemand.

Die Datei reist auf einem Ref, wie im alten Entwurf. Farbe zufaellig, Markierung rechts auf der Karte, Filter auf den Milestone - alles unveraendert uebernommen.

Offen und Sache der Brainstorm-Lane: ob ein Ticket in mehr als einem Milestone stehen darf, und wenn ja, was die rechte Kante der Karte dann zeigt. Bei parallelen, frei benannten Gruppen ist Mehrfachzugehoerigkeit wahrscheinlich, beim Sprint war sie es nicht.

Was der Milestone NICHT loest, und das war die Ausgangsnot vom 14.09.: 'zeig mir nur die zwanzig, die ich in dieser Sitzung durchgehe'. Ein Milestone mit vierzig Tickets filtert auf vierzig Tickets. Wenn diese Not bleibt, ist sie ein eigenes Ticket und nicht dieses.
