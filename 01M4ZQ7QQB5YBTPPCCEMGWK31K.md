---
id: 01M4ZQ7QQB5YBTPPCCEMGWK31K
title: "Ein archiviertes Ticket verschwindet sofort von allen Boards, nicht erst wenn der Snapshot den Ref abraeumt"
status: in-progress
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
updated-at: 2026-10-15T13:19:25Z
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

- [x] Test zuerst, core/ticket/archiveref_test.go (neu): Fake-Recorder refRec (Record/RecordFiled/RecordDelete merken den letzten Inhalt). TestArchivePutsTheArchivedLineOnTheRefOnly: Archive(id) -> der RecordFiled-Inhalt dekodiert (ParseDoc+Decode) mit Archived == '<yyyy-mm-dd heute> <Actor>' und gleichem updated-at wie die Datei; die Datei in .jaira/archive/ ist byte-gleich zur Datei vor dem Verschieben und hat keinen archived-Key (DoD 1)
- [x] Test zuerst, gleiche Datei: TestAnArchivedRefOnlyTicketIsOffTheListButStillLoads: s.Source = refSource (core/ticket/trim_test.go:218) mit einem Ticket mit Archived und einem ohne -> List() nennt nur das ohne; Load(Handle des archivierten) liefert es mit ReadOnly und Archived (DoD 1; List ist, was board, list, next und resume lesen)
- [x] Test zuerst, gleiche Datei: TestRestoreTakesTheArchivedLineOffTheRef (Archive, dann Restore(datei) -> letzter Record-Inhalt == Bytes der wiederhergestellten Datei, ohne archived-Key) und TestRestoreDoesNotWriteOverANewerRef (refRec implementiert CheckFresh und meldet ErrStale -> Restore gibt dst UND einen Fehler mit errors.Is(err, ErrStale) zurueck, Datei liegt in tickets/, kein Record) (DoD 2)
- [x] core/ticket, Feld: FieldArchived = "archived" in schema.go neben FieldTags; Ticket.Archived string (Kommentar: Wert '<yyyy-mm-dd> <wer>', leer = nicht archiviert, steht nur im Inhalt auf dem Ref); decodeFields liest es (schema.go ~640). Nicht in canonicalOrder
- [x] core/ticket, Archive: Store.Archive (store.go:330) gibt RecordFiled die Dateibytes plus Zeile archived (Doc.SetScalar, Wert heute + ' ' + s.Actor, getrimmt); die Datei in .jaira/archive/ und updated-at/updated-by bleiben unberuehrt. Logbook (store.go:402) bleibt bei recordFiled. Doc-Kommentar von Archive (store.go:297-300) um den Ref ergaenzen
- [x] core/ticket, List und Restore: Store.List (store.go:867) haengt nur Extras ohne Archived an. Store.Restore (store.go:634): nur fuer eine Datei, die aus ArchiveDir kam, nach dem Rename Load(id), s.CheckFresh, dann s.record(id, Dateibytes); Rueckgabe dst plus dieser Fehler (ErrStale/ErrNotRecorded). Logbook- und Milestone-Restore unveraendert. Tests aus Schritt 1-3 muessen kompilieren (go vet ./core/ticket/)
- [x] Test zuerst, internal/cli/archiveref_test.go (neu), mit twoBoards, runAndSend, aTicket, handleFromList, gitOut (internal/cli/milestoneref_test.go) und jsonCLI (internal/cli/tags_test.go:64): TestAnArchivedTicketLeavesTheOtherBoardAtItsNextFetch: ada aTicket + pull; grace fetch -> list --json nennt den Handle (Vorbedingung); ada runAndSend archive <h> -> Ausgabe enthaelt 'only keeps the record in the repository'; gitref.Repo{Dir: ada, Remote: origin}.Read(id) traegt 'archived:'; git ls-files .jaira/archive in ada leer (nichts committet); grace fetch -> list --json nennt ihn nicht mehr; grace show <h> enthaelt 'archived' (DoD 1, 3)
- [x] Test zuerst, gleiche Datei: TestARestoredTicketComesBackOnTheOtherBoard: wie oben archivieren, dann ada runAndSend restore <datei> -> ada-Ref ohne 'archived:'; grace fetch -> list --json nennt den Handle wieder (DoD 2)
- [x] internal/cli umsetzen: archive.go Textausgabe (:86-88) bekommt, nur wenn refs.Usable() == nil, die Zeile 'Every other board drops it at its next fetch; committing .jaira/archive/<datei> only keeps the record in the repository.' --json bleibt gleich. restore RunE (archive.go:127): dst leer -> not_archived wie heute; dst gesetzt + Fehler -> warnRef auf stderr, restore gilt als gelungen. show (tickets.go:787-794): bei t.Archived != '' die Zeile 'off the board: archived <wert>, only on its ref' und im --json der Key 'archived'. Long-Hilfe von archive und restore um je einen Satz zum Ref ergaenzen
- [ ] Doku, die dadurch falsch wird, nachziehen (kein Code): RecordFiled-Kommentar core/refsync/refsync.go:166-179, README.md:579-584 (archive nimmt das Ticket sofort von allen Boards, logbook wie bisher), docs/COMMANDS.md:155 (archive) und :158 (restore nimmt die Zeile vom Ref)
- [ ] core/release/notes/GWK31K.md: genau eine '- '-Zeile, nie umbrochen, als Anweisung: wer mit 'jaira archive <id>' archiviert, nimmt das Ticket beim naechsten fetch von jedem Board, ohne Commit; 'jaira restore <datei>' bringt es ueberall zurueck, 'jaira show <id>' sagt archived; das Committen von .jaira/archive/<datei> haelt nur den Eintrag im Repository fest
- [ ] Pruefen ohne go test (Berk verbietet es): gofmt -l core internal leer, go vet ./... (kompiliert die neuen Tests mit), go build -race -o /tmp/claude-1000/-home-berk-git-jAIra/a751ce77-35ca-405d-9c2e-d6ae252f56ea/scratchpad/gwk31k/jaira ./cmd/jaira
- [ ] fetch-Bericht: refsync.Incoming liest auch archived (Arrival.Archived, json 'archived'); printArrivals (internal/cli/fetch.go ~272) zeigt bei einem markierten Ref 'archived' statt der Lane, announceArrivals meldet ihn nicht als 'a ticket is yours'; Test in internal/cli/archiveref_test.go (grace fetch nennt den Handle mit 'archived')

## Progress
- **2026-10-15 12:37 · BeMuCa** — Berk am 15.10.: 'Archiv sollte auch wie Ticket funktionieren, damit jeder es schnell genug sieht.' Seine Frage, ob Tickets erst nach einem Commit fuer andere sichtbar sind: nein. 'jaira create' legt das Ticket auf seinen Ref und pusht ihn sofort; andere sehen es beim naechsten fetch, ohne Commit. Archivieren muss genauso gehen: die Markierung kommt auf den Ref und wird sofort gepusht, andere sehen das Ticket beim naechsten fetch verschwinden, ohne dass die Archiv-Datei vorher committet sein muss. Das Committen der Archiv-Datei haelt danach nur den Eintrag im Repository fest; niemandes Sicht wartet darauf.
- **2026-10-15 12:56 · BeMuCa** — pre-process, Teil 1 - die Marke und wo sie steht:
- Marke = neue Frontmatter-Zeile 'archived: <yyyy-mm-dd> <wer>', z.B. 'archived: 2026-10-15 BeMuCa'. Leer oder fehlend = nicht archiviert.
- Nicht 'status: archived': status ist bei Tickets die Lane (core/ticket/schema.go:19), und restore muss die Lane unveraendert zurueckgeben. Wertform wie N71NVGs 'deleted: <datum> <wer>'.
- Unbekannte Keys ueberleben jedes Schreiben: geschrieben wird ueber Doc, nie ueber das Struct (core/ticket/frontmatter.go:1-13, schema.go:216-218). Lesen braucht trotzdem ein Feld: Ticket.Archived (string) + FieldArchived, dekodiert in decodeFields (schema.go ~640).
- Die Marke steht NUR im Inhalt auf dem Ref. Die Datei in .jaira/archive/ bleibt byte-gleich, updated-at/updated-by bleiben gleich.
- Warum nur auf dem Ref: Datei und Ref haben denselben updated-at. CheckFresh vergleicht nur updated-at (core/refsync/refsync.go:1205-1230), also wird ein restore nicht als stale verweigert. Und keine Datei bekommt je eine Restzeile: Doc kann keinen Key loeschen (frontmatter.go hat kein Delete; merge.go:386-392 schreibt sogar bewusst leer statt zu loeschen).
- Archive (core/ticket/store.go:330): RecordFiled bekommt Dateibytes + Zeile statt nur Dateibytes.
- Restore (store.go:634): schreibt heute GAR NICHTS auf den Ref. Neu, nur fuer ein Ticket aus .jaira/archive/: nach dem Zurueckschieben CheckFresh, dann record(Dateibytes) -> Ref wieder ohne Zeile, kommt ueberall zurueck. Ist der Ref neuer (jemand schrieb seitdem), wird nichts geschrieben: Datei ist zurueck, CLI warnt mit dem ErrStale-Text ('jaira pull').
- Restore-CLI (internal/cli/archive.go:127) macht heute aus jedem Fehler not_archived. Darum gibt Restore dst + Fehler zurueck, und die CLI warnt nur, wenn dst gesetzt ist.
- Kein Ref mehr (Snapshot hat ihn abgeraeumt): record legt ihn neu an (Lease ''), das Ticket ist dann wie gepullt. Ohne Remote macht Record nichts (refsync.go:124-131).
- Snapshot kopiert den Ref-Inhalt nach board/ im Branch jaira/board (core/snapshot/snapshot.go:17-27): die Sicherung sagt dann auch archiviert, richtig so. Der Reaper prueft nur Pfade .jaira/archive/<id>* (core/gitref/gitref.go:770-786), unveraendert.
- **2026-10-15 12:56 · BeMuCa** — pre-process, Teil 2 - wer was sieht, Push, Ausgabe:
- Ausblenden an EINER Stelle: Store.List (core/ticket/store.go:860-867) haengt nur Extras ohne Archived an. Durch List lesen: list und next (loadEnv, internal/cli/root.go:288), resume (internal/cli/resume.go:88), Board-TUI (internal/tui/model.go:420, home.go:167).
- Nicht in refsync.Extra: Extra bedient auch Load (fromSource, store.go:1160). Dann saehe 'jaira show <h>' nur 'not found'.
- show: Load findet das Ticket weiter (ReadOnly; Schreiben wird mit 'pull it first' verweigert). show sagt eine Zeile 'off the board: archived <wert>', --json bekommt 'archived'. Grund: wer es gestern sah, erfaehrt, was passiert ist.
- pull: unveraendert. Archivieren braucht eine Datei (onlyOnRef, store.go:166-179), also meist ein pull -> assignee = Archivierer -> andere bekommen ErrTaken (refsync.go:1364-1372). Mit --steal kommt die Datei samt Zeile; dort ist sie eine Datei und wird gezeigt, fuer alle anderen bleibt sie verborgen. Randfall, kein Code.
- Klon mit Datei in tickets/: fetch schreibt nie Ticket-Dateien (Incoming refsync.go:936-999; Kommentar internal/cli/fetch.go:132-137). Extra ueberspringt lokale IDs (refsync.go:1536-1541). Er sieht weiter seine Datei, bis der Branch mit dem Verschieben gemergt ist. pull schreibt nur eine aeltere Datei um (refresh refsync.go:1284-1323, isStale nach updated-at); die Marke aendert updated-at nicht. Nichts wird ueberschrieben, nichts geht verloren. DoD 1 verlangt nur ref-only.
- Push (DoD 3): passiert schon heute vor Prozessende. recordFiled -> refsync.Record setzt dirty (refsync.go:124-141) -> Execute ruft afterCommand -> flushRefs (root.go:103-106, refs.go:204-207, 246-252). create schickt synchron in fileOnRefOnly (refs.go:83-110), nur weil es danach die Datei loescht. Kein neuer Push-Code. TUI-Taste x (internal/tui/model.go:2192) geht durch dieselbe Store.Archive. Im Test ersetzt runAndSend afterCommand (internal/cli/milestoneref_test.go:17-25).
- Ausgabe archive (internal/cli/archive.go:86-88): neue Zeile nur wenn refs.Usable() == nil. Ohne Remote waere 'jedes Board laesst es beim naechsten fetch fallen' falsch.
- --json von archive bleibt gleich: der Satz ist Rat fuer Menschen. Ob gesendet wurde, steht wie bei jedem Schreiben auf stderr (flushRefs/warnRef, refs.go:246-358).
- **2026-10-15 12:56 · BeMuCa** — pre-process, Teil 3 - Gleichzeitigkeit, Logbook, N71NVG, bekannte Luecken:
- A archiviert, B schreibt dasselbe Ticket. B hat dafuer eine Datei: geteilter Branch oder --steal. Beide leasen denselben SHA.
- Der zweite Push wird abgewiesen und VERWORFEN (core/outbox/outbox.go:424-427, 470-479). Tickets werden dabei nicht gemergt, nur Milestones und Epics (refsync.go:680-691).
- B zuerst: As Marke kommt nie an. A sieht '<h> was not sent: <B ...>' (refs.go:262-272). Ticket bleibt ueberall sichtbar - die sichere Richtung. Neu archivieren: restore (schreibt nicht, Ref ist neuer) -> pull -> archive.
- A zuerst: Bs naechster Schreibvorgang nach seinem fetch ist nicht stale (gleicher updated-at) und ersetzt den Ref-Inhalt. Marke weg, Ticket wieder sichtbar. Dasselbe last-writer-wins wie fuer jedes Feld auf dem Ref.
- core/merge: Die Marke steht fast nie in einer Datei in tickets/. Wenn doch (nach pull eines markierten Refs), nimmt mergeScalar die einseitige Aenderung (core/merge/merge.go:155-158, 250-262). Kein Code, kein Test geplant.
- Ausser Scope: die Marke gegen Bs spaeteren Push schuetzen. Braucht Merge-on-reject fuer Tickets (wie mergeRefused, refsync.go:705ff) oder eine Sperre in Record wie N71NVGs ErrMilestoneDeleted. Folgeticket, falls gewuenscht.
- Logbook (store.go:402) ruft recordFiled weiter unmarkiert, ausser Scope. Die Zeile wuerde dort genauso passen. Aber das kehrt die Begruendung in RecordFiled (refsync.go:166-179, README.md:579-584: Ref sichtbar lassen, damit niemand fertige Arbeit doppelt aufschreibt) auch fuer fertige Arbeit um. Berks Entscheidung.
- Bekannt, nicht in der DoD, kein Code: 'jaira fetch' listet markierte Refs weiter, und 'a ticket is yours' kann fuer ein markiertes Ticket feuern (internal/cli/fetch.go:211-229, 272-291). Ein Link auf ein anderswo archiviertes Ticket heisst dort 'unknown' statt 'archive'; Satisfied bleibt false wie beim Archivierer (core/link/link.go:304-313).
- N71NVG: nichts im Code geteilt. Milestones haben ein eigenes Format (core/milestone: 'status: filed' + 'deleted:'), Tickets gehen ueber core/ticket.Doc, und status heisst dort Lane. Ein gemeinsamer Helper muesste nach Art verzweigen - kein Gewinn, nur Kopplung.
- Geteilt ist die Regel 'der Ref bleibt stehen und bekommt eine Zeile' und die Wertform '<yyyy-mm-dd> <wer>'. GWK31K braucht keinen Code von N71NVG.
- Merge-Risiko mit N71NVG nur in der Doku: beide aendern docs/COMMANDS.md (N71NVG :138/:157, GWK31K :155/:158) und README.md.
- **2026-10-15 13:14 · BeMuCa** — Treiber nach Handtest des Core-Commits 4deebea5 (zwei Klone ueber ein lokales bare repo): list/next/resume auf B blenden beide archivierten Tickets aus, restore bringt sie zurueck. ABER B's 'jaira fetch' meldete danach weiter 'W9SM3S todo Probe todo @you new ref-only' - der Bericht nennt die alte Lane eines archivierten Tickets, und announceArrivals kann dafuer 'a ticket is yours' ausloesen. Darum Plan-Schritt 13: fetch sagt 'archived'. Nicht gemacht (Nachfolger): ein Klon, der das Ticket als Datei hat, wird nicht ueber Departed() informiert; Departed() reagiert nur auf einen verschwundenen Ref.
