---
id: 01M4ZQ7QQB5YBTPPCCEMGWK31K
title: "Ein archiviertes Ticket verschwindet sofort von allen Boards, nicht erst wenn der Snapshot den Ref abraeumt"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Wer ein Ticket archiviert, nimmt es damit von jedem Board im Team, sofort beim naechsten fetch, und ein restore bringt es ueberall zurueck"
context: |-
  Was heute falsch ist: 'jaira archive <id>' verschiebt die Datei nach .jaira/archive/. Auf den Ref schreibt es aber den unveraenderten Inhalt.
  Code: Store.Archive (core/ticket/store.go:301-330) ruft recordFiled mit dem Dateiinhalt auf. Die Lane bleibt z.B. backlog, und nichts im Ref sagt 'archiviert'.
  Folge: Auf jedem anderen Klon steht das Ticket weiter als unberuehrtes ref-only Ticket auf dem Board. Berk hat das am 15.10. bemerkt.
  Der Ref verschwindet erst spaeter. Das macht der Snapshot-Lauf (core/snapshot, alle 72 h), und nur nachdem die Ticket-Datei in einem Branch angekommen ist, der zaehlt. Dafuer muss die Archiv-Datei committet und gemergt sein (refsync.go RecordFiled-Kommentar: 'Filing does not remove the ref, and that is the point').
  'Ungetrackt' heisst hier: .jaira/archive/ ist nicht gitignored (git check-ignore sagt nein). Die verschobene Datei ist nur neu und noch nicht committet. Das kann leicht vergessen werden.
  Vorbild, das schon funktioniert: Ein abgelegter Milestone traegt 'status: filed' auf seinem Ref, und jeder Klon blendet ihn beim naechsten fetch aus (jaira logbook --help).
  Beim Ticket-Logbook faellt es weniger auf, weil die Lane dort done ist.
definition-of-done: "Archivieren schreibt eine Markierung in den Inhalt auf dem Ref (z.B. 'status: archived', wie 'filed' bei Milestones); board, list, next und resume blenden ein ref-only Ticket mit dieser Markierung aus; Test"
tags:
  - concurrency
blocked-by: []
related:
  - 01M4TRC0ZE0GNWDWA1XDK6BGB0
commits: []
created-at: 2026-10-15T12:05:20Z
updated-at: 2026-10-15T12:37:17Z
updated-by: BeMuCa
---

# Ein archiviertes Ticket verschwindet sofort von allen Boards, nicht erst wenn der Snapshot den Ref abraeumt

## Definition of Done

- [ ] Archivieren schreibt eine Markierung in den Inhalt auf dem Ref (z.B. 'status: archived', wie 'filed' bei Milestones); board, list, next und resume blenden ein ref-only Ticket mit dieser Markierung aus; Test
- [ ] 'jaira restore' eines archivierten Tickets entfernt die Markierung auf dem Ref wieder, sodass es ueberall zurueckkommt; Test
- [ ] Die Ausgabe von 'jaira archive' sagt, dass .jaira/archive/<datei> committet werden muss (wie beim logbook), oder das Ticket begruendet, warum nicht
- [ ] Eine Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-15 12:37 · BeMuCa** — Berk am 15.10.: 'Archiv sollte auch wie Ticket funktionieren, damit jeder es schnell genug sieht.' Seine Frage, ob Tickets erst nach einem Commit fuer andere sichtbar sind: nein. 'jaira create' legt das Ticket auf seinen Ref und pusht ihn sofort; andere sehen es beim naechsten fetch, ohne Commit. Archivieren muss genauso gehen: die Markierung kommt auf den Ref und wird sofort gepusht, andere sehen das Ticket beim naechsten fetch verschwinden, ohne dass die Archiv-Datei vorher committet sein muss. Das Committen der Archiv-Datei haelt danach nur den Eintrag im Repository fest; niemandes Sicht wartet darauf.
