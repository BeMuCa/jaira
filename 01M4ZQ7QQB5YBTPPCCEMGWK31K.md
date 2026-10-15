---
id: 01M4ZQ7QQB5YBTPPCCEMGWK31K
title: "Ein archiviertes Ticket verschwindet sofort von allen Boards, nicht erst wenn der Snapshot den Ref abraeumt"
status: pre-process
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
updated-at: 2026-10-15T12:56:06Z
updated-by: BeMuCa
---

# Ein archiviertes Ticket verschwindet sofort von allen Boards, nicht erst wenn der Snapshot den Ref abraeumt

## Definition of Done

- [ ] Archivieren schreibt eine Markierung in den Inhalt auf dem Ref (z.B. 'status: archived', wie 'filed' bei Milestones); board, list, next und resume blenden ein ref-only Ticket mit dieser Markierung aus; Test
- [ ] 'jaira restore' eines archivierten Tickets entfernt die Markierung auf dem Ref wieder, sodass es ueberall zurueckkommt; Test
- [ ] Archivieren pusht die Markierung sofort auf den Ref, wie 'jaira create' ein neues Ticket pusht: ein anderer Klon sieht das Ticket nach seinem naechsten fetch nicht mehr, ohne dass die Archiv-Datei committet ist. Die Ausgabe von 'jaira archive' sagt, dass das Committen von .jaira/archive/<datei> nur den Eintrag im Repository festhaelt; Test
- [ ] Eine Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [ ] Test zuerst, core/ticket/archiveref_test.go (neu): Fake-Recorder refRec (Record/RecordFiled/RecordDelete merken den letzten Inhalt). TestArchivePutsTheArchivedLineOnTheRefOnly: Archive(id) -> der RecordFiled-Inhalt dekodiert (ParseDoc+Decode) mit Archived == '<yyyy-mm-dd heute> <Actor>' und gleichem updated-at wie die Datei; die Datei in .jaira/archive/ ist byte-gleich zur Datei vor dem Verschieben und hat keinen archived-Key (DoD 1)
- [ ] Test zuerst, gleiche Datei: TestAnArchivedRefOnlyTicketIsOffTheListButStillLoads: s.Source = refSource (core/ticket/trim_test.go:218) mit einem Ticket mit Archived und einem ohne -> List() nennt nur das ohne; Load(Handle des archivierten) liefert es mit ReadOnly und Archived (DoD 1; List ist, was board, list, next und resume lesen)
- [ ] Test zuerst, gleiche Datei: TestRestoreTakesTheArchivedLineOffTheRef (Archive, dann Restore(datei) -> letzter Record-Inhalt == Bytes der wiederhergestellten Datei, ohne archived-Key) und TestRestoreDoesNotWriteOverANewerRef (refRec implementiert CheckFresh und meldet ErrStale -> Restore gibt dst UND einen Fehler mit errors.Is(err, ErrStale) zurueck, Datei liegt in tickets/, kein Record) (DoD 2)
- [ ] core/ticket, Feld: FieldArchived = "archived" in schema.go neben FieldTags; Ticket.Archived string (Kommentar: Wert '<yyyy-mm-dd> <wer>', leer = nicht archiviert, steht nur im Inhalt auf dem Ref); decodeFields liest es (schema.go ~640). Nicht in canonicalOrder
- [ ] core/ticket, Archive: Store.Archive (store.go:330) gibt RecordFiled die Dateibytes plus Zeile archived (Doc.SetScalar, Wert heute + ' ' + s.Actor, getrimmt); die Datei in .jaira/archive/ und updated-at/updated-by bleiben unberuehrt. Logbook (store.go:402) bleibt bei recordFiled. Doc-Kommentar von Archive (store.go:297-300) um den Ref ergaenzen
- [ ] core/ticket, List und Restore: Store.List (store.go:867) haengt nur Extras ohne Archived an. Store.Restore (store.go:634): nur fuer eine Datei, die aus ArchiveDir kam, nach dem Rename Load(id), s.CheckFresh, dann s.record(id, Dateibytes); Rueckgabe dst plus dieser Fehler (ErrStale/ErrNotRecorded). Logbook- und Milestone-Restore unveraendert. Tests aus Schritt 1-3 muessen kompilieren (go vet ./core/ticket/)
- [ ] Test zuerst, internal/cli/archiveref_test.go (neu), mit twoBoards, runAndSend, aTicket, handleFromList, gitOut (internal/cli/milestoneref_test.go) und jsonCLI (internal/cli/tags_test.go:64): TestAnArchivedTicketLeavesTheOtherBoardAtItsNextFetch: ada aTicket + pull; grace fetch -> list --json nennt den Handle (Vorbedingung); ada runAndSend archive <h> -> Ausgabe enthaelt 'only keeps the record in the repository'; gitref.Repo{Dir: ada, Remote: origin}.Read(id) traegt 'archived:'; git ls-files .jaira/archive in ada leer (nichts committet); grace fetch -> list --json nennt ihn nicht mehr; grace show <h> enthaelt 'archived' (DoD 1, 3)
- [ ] Test zuerst, gleiche Datei: TestARestoredTicketComesBackOnTheOtherBoard: wie oben archivieren, dann ada runAndSend restore <datei> -> ada-Ref ohne 'archived:'; grace fetch -> list --json nennt den Handle wieder (DoD 2)
- [ ] internal/cli umsetzen: archive.go Textausgabe (:86-88) bekommt, nur wenn refs.Usable() == nil, die Zeile 'Every other board drops it at its next fetch; committing .jaira/archive/<datei> only keeps the record in the repository.' --json bleibt gleich. restore RunE (archive.go:127): dst leer -> not_archived wie heute; dst gesetzt + Fehler -> warnRef auf stderr, restore gilt als gelungen. show (tickets.go:787-794): bei t.Archived != '' die Zeile 'off the board: archived <wert>, only on its ref' und im --json der Key 'archived'. Long-Hilfe von archive und restore um je einen Satz zum Ref ergaenzen
- [ ] Doku, die dadurch falsch wird, nachziehen (kein Code): RecordFiled-Kommentar core/refsync/refsync.go:166-179, README.md:579-584 (archive nimmt das Ticket sofort von allen Boards, logbook wie bisher), docs/COMMANDS.md:155 (archive) und :158 (restore nimmt die Zeile vom Ref)
- [ ] core/release/notes/GWK31K.md: genau eine '- '-Zeile, nie umbrochen, als Anweisung: wer mit 'jaira archive <id>' archiviert, nimmt das Ticket beim naechsten fetch von jedem Board, ohne Commit; 'jaira restore <datei>' bringt es ueberall zurueck, 'jaira show <id>' sagt archived; das Committen von .jaira/archive/<datei> haelt nur den Eintrag im Repository fest
- [ ] Pruefen ohne go test (Berk verbietet es): gofmt -l core internal leer, go vet ./... (kompiliert die neuen Tests mit), go build -race -o /tmp/claude-1000/-home-berk-git-jAIra/a751ce77-35ca-405d-9c2e-d6ae252f56ea/scratchpad/gwk31k/jaira ./cmd/jaira

## Progress
- **2026-10-15 12:37 · BeMuCa** — Berk am 15.10.: 'Archiv sollte auch wie Ticket funktionieren, damit jeder es schnell genug sieht.' Seine Frage, ob Tickets erst nach einem Commit fuer andere sichtbar sind: nein. 'jaira create' legt das Ticket auf seinen Ref und pusht ihn sofort; andere sehen es beim naechsten fetch, ohne Commit. Archivieren muss genauso gehen: die Markierung kommt auf den Ref und wird sofort gepusht, andere sehen das Ticket beim naechsten fetch verschwinden, ohne dass die Archiv-Datei vorher committet sein muss. Das Committen der Archiv-Datei haelt danach nur den Eintrag im Repository fest; niemandes Sicht wartet darauf.
