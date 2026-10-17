---
id: 01M4TRC0ZE0GNWDWA1XDK6BGB0
title: "Das Board zeigt den Stand des Refs, wenn die Ticket-Datei dahinter zurueckliegt"
status: testing
ready: true
creator: BeMuCa
goal: "Board, 'jaira list' und 'jaira show' zeigen den neuesten Stand eines Tickets: liegt die Datei im Checkout hinter ihrem Ref, den des Refs, mit einem Kennzeichen, dass 'jaira pull' faellig ist - und ein Schreiben auf die veraltete Datei bleibt verweigert."
context: |-
  Am 2026-10-13 stand D0SAHM auf Berks Board (Haupt-Checkout /home/berk/git/jAIra) in in-progress. Sein Ref refs/jaira/tickets/01M1KFMETY9MKH4V38TVD0SAHM sagte schon signoff.
  Grund: Board und 'jaira list' lesen die Datei im Checkout, und die lag hinter dem Ref. Gearbeitet wurde das Ticket in einem anderen Worktree (/home/berk/git/jAIra-D0SAHM).
  'jaira pull D0SAHM' wusste es: 'D0SAHM was behind its ref and is up to date now - was: in-progress ... now: signoff'. jaira erkennt den Fall also schon.
  Es war das zweite Mal in derselben Session. Berk haette fast ein Ticket fuer haengend gehalten, das auf seine Abnahme wartete.
  JDXV60 beschreibt dasselbe Problem zwischen zwei Worktrees (2026-09-15), vor den Refs. Das hier ist der Rest, der mit Refs bleibt.
  Offen fuer brainstorm: Board liest dann den Ref, oder zeigt nur ein Kennzeichen 'Ref ist weiter' an der Karte.
definition-of-done: "Liegt die Datei eines Tickets hinter seinem Ref, zeigt das Board auf der Karte den Stand des Refs oder ein Kennzeichen dafuer; 'jaira list' und 'jaira show' sagen es ebenso; ein Test baut genau diesen Fall"
tags:
  - concurrency
  - tui
blocked-by: []
related: []
follows: 01M2HW4788Z5G0M05G3PJDXV60
commits: []
created-at: 2026-10-13T13:48:57Z
updated-at: 2026-10-17T20:59:37Z
assignee: berkc
updated-by: berkc
claimed-by: berk-ca-297639
claimed-at: 2026-10-17T20:39:34Z
outcome-what: "Board, jaira list und jaira show zeigen bei einer Ticket-Datei, die hinter ihrem Ref liegt, den Stand des Refs und markieren sie (Karte '⇣ behind', Zeile '[behind ref: pull it]', JSON behind-ref, show-Zeile 'behind')."
outcome-why: "Ein in einem anderen Worktree oder von einem Teammitglied bewegtes Ticket stand hier in seiner alten Lane; die Erkennung (isStale, wie bei pull) gab es schon, sie wurde nur beim Lesen nicht genutzt."
outcome-resolves: K6BGB0; deckt JDXV60 DoD 1 fuer den Ref-Fall ab
review-summary: "none"
round: 3
review-gaps: "Runde 3: MilestoneHad und das neue Had teilen sich had(ref, content) statt die rev-list-Schleife zu verdoppeln. Unveraendert stehen gelassen: zweites for-each-ref pro List (~2 ms bei 60 Refs). Kein toter Code: das alte freie isStale ist ersetzt, alle drei Aufrufer nutzen die Methode."
test-verdict: "fail: liegen Datei und Ref in derselben Sekunde (updated-at gleich), gilt die Datei nicht als hinter dem Ref - Board zeigt den alten Stand und ein note schreibt ihn still auf den Ref zurueck (todo -> backlog)"
---

# Das Board zeigt den Stand des Refs, wenn die Ticket-Datei dahinter zurueckliegt

## Definition of Done

- [x] Liegt die Datei eines Tickets hinter seinem Ref, zeigt das Board auf der Karte den Stand des Refs oder ein Kennzeichen dafuer; 'jaira list' und 'jaira show' sagen es ebenso; ein Test baut genau diesen Fall
  proof: internal/cli/stalefile_test.go TestListAndShowGiveTheRefStateOfAFileBehindIt + TestAFileBehindItsRefWithinTheSameSecondIsStillBehind; internal/tui/checklist_test.go TestCardMarksAFileBehindItsRef; core/ticket/store.go Store.newest/Newest; core/refsync/refsync.go Syncer.Ahead/isStale
- [x] go test ./... -race gruen; Release-Notiz unter core/release/notes/
  proof: go vet ./... && go test -race ./... all ok 2026-10-17; core/release/notes/K6BGB0.md

## Options

- [x] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] ticket: Feld BehindRef am Ticket (Laufzeit, kein Frontmatter) und optionales Interface AheadSource; Store.List ersetzt Datei-Tickets, deren Ref weiter ist, durch das Ref-Ticket (Path bleibt die Datei, BehindRef=true); Store.Newest(t) fuer show
- [x] refsync: Syncer.Ahead(files) - ListHashes + extraCache nach Ref-Hash + ReadMany nur fuer neue Hashes, Urteil per isStale (dieselbe Definition wie pull/CheckFresh)
- [x] Ausgabe: TUI-Karte Kennzeichen 'pull due', jaira list Zeile [behind ref - pull], JSON behind-ref, show-Text Zeile; Load bleibt Datei
- [x] Tests: refsync (Datei hinter Ref -> List zeigt Ref-Lane+BehindRef; Mutate verweigert mit ErrStale, Ref unveraendert), TUI-Karte zeigt Kennzeichen
- [x] go vet + go test -race; manuell mit zwei Worktrees eines Temp-Repos; Release-Notiz core/release/notes/K6BGB0.md

## Progress
- **2026-10-14 14:08 · BeMuCa** — Dritte Variante am 14.10.: D0SAHM stand auf Berks Board in signoff, obwohl es laengst abgenommen und abgelegt war (logbook bemuca-20261013, per PR #74 in master). Grund: Der Haupt-Checkout /home/berk/git/jAIra stand noch auf dem alten master fa222c5d, und dort lag eine veraltete, geaenderte Kopie .jaira/tickets/...D0SAHM...md mit status: signoff. Den Ref refs/jaira/tickets/<ULID> gibt es nach dem Ablegen nicht mehr, also kann 'zeig den Stand des Refs' hier nicht helfen. Das Board muss auch erkennen: Datei ohne Ref, deren Ticket upstream im Logbuch liegt -> als abgelegt bzw. veraltet zeigen statt als offen. Behoben fuer heute durch git stash der veralteten Kopien + git pull master.
- **2026-10-17 20:39 · berkc** — berkc took this ticket over from BeMuCa
- **2026-10-17 20:39 · berkc** — Brainstorm (2026-10-17, entschieden vom Nutzer): Was falsch ist laut Code: Store.List zeigt fuer jede Ticket-Datei die Datei; Syncer.Extra (core/refsync/refsync.go) ueberspringt jeden Ref, fuer den eine Datei existiert ('if local[id] || filed[id] continue'). Liegt die Datei hinter dem Ref, zeigen Board, list und show die alte Lane. 'jaira pull' erkennt den Fall schon mit isStale (refsync.go: ref.updated_at strikt nach file.updated_at) - dieselbe Definition wie CheckFresh. Schreibseite: Mutate/Archive/Logbook rufen schon CheckFresh und verweigern mit StaleError ('run jaira pull'), ein Schreiben auf eine veraltete Datei kann den Ref also heute schon nicht still zurueckdrehen.
Optionen: (A) nur Kennzeichen an der Karte, Datei-Stand bleibt sichtbar - billig, aber Board zeigt weiter die falsche Lane. (B) Board/list/show zeigen den Stand des Refs plus Kennzeichen 'Datei veraltet' - zeigt die echte Lane, Schreiben bleibt an der Datei und wird von CheckFresh verweigert. (C) automatisch pullen beim Lesen - schreibt Dateien ohne Befehl, verletzt 'Lesen schreibt nie'.
Gewaehlt: B (Nutzerentscheidung 2026-10-17): newest wins in List und show, Marker an Karte/Zeile/JSON. Load bleibt bei der Datei, weil Archive/Logbook/move den geladenen Ticket gegen CheckFresh pruefen - ein Ref-Ticket dort wuerde die Pruefung bestehen und die alte Datei ablegen.
- **2026-10-17 20:40 · berkc** — Plan-Begruendung: Load wird NICHT umgestellt. Archive, Logbook und der move-Gate laden per Load und pruefen danach CheckFresh(t) - wuerde Load das Ref-Ticket liefern, bestuende es die Pruefung (gleiches updated_at) und Archive wuerde die alte Datei ablegen und deren Stand per RecordFiled auf den Ref schreiben. Deshalb nur List (Board, list, next, web) und show ueber Store.Newest. Mutate liest unter Lock die Datei neu und ruft CheckFresh -> StaleError; das bleibt die Schreib-Sperre. Separates optionales Interface statt Extra-Signatur zu aendern: ein zusaetzliches 'git for-each-ref' pro List, Decodes nur bei neuem Ref-Hash (extraCache geteilt). Nicht abgedeckt: D0SAHM-Variante (Datei ohne Ref, Ticket upstream abgelegt) - dort gibt es keinen Ref zum Vergleichen.
- **2026-10-17 20:49 · berkc** — In-progress: Store.List ersetzt Datei-Tickets, deren Ref nach isStale weiter ist, durch das Ref-Ticket (BehindRef, Path=Datei); show nutzt Store.Newest nur fuer Tickets auf dem Board (nicht Logbuch). Schreibseite war schon sicher (Mutate/Archive/Logbook/move rufen CheckFresh -> StaleError, exit 3); der Test prueft das am selben Fall, Ref-SHA unveraendert. Manuell nachgestellt mit zwei Worktrees eines Temp-Repos: main zeigt todo + [behind ref: pull it], note verweigert, Ref-SHA gleich; nach pull Marker weg. Gleiches-Sekunde-Grenzfall: isStale vergleicht updated_at mit Sekundenaufloesung - zwei Schreibvorgaenge in derselben Sekunde gelten nicht als 'weiter' (bestehende Grenze von pull/CheckFresh, nicht neu).
- **2026-10-17 20:49 · berkc** — critique: store.go newest vereinfachen - ReadOnly-Filter und Kopie entfernen, ts direkt an Ahead (Begruendung im review-summary). Ahead/Extra teilen sich eine fast gleiche Lese-und-Decode-Schleife: das ist Sache der optimize-Lane, hier stehen gelassen. Load bleibt bewusst bei der Datei (siehe Plan-Notiz).
- **2026-10-17 20:50 · berkc** — in-progress (Runde 2): store.go newest vereinfacht wie von critique verlangt; Tests gruen.
- **2026-10-17 20:50 · berkc** — critique Runde 2: nichts mehr - newest ist jetzt drei Zeilen; Load bleibt bewusst Datei (Plan-Notiz); die doppelte Decode-Schleife Ahead/Extra geht an optimize.
- **2026-10-17 20:52 · berkc** — testing: Zwei-Worktree-Nachstellung ohne Pause (repro.sh): pull in main 20:51:52, move --to todo in wt2 20:51:52, danach 'jaira note' in main: Noted, exit 0, Ref danach 'status: backlog updated-at 20:51:53' - der todo-Stand ist still weg. Mit dem installierten 0.5.2 genauso, also vorbestehend: isStale vergleicht updated-at mit Sekundenaufloesung und urteilt bei Gleichstand 'nicht hinter'. Vorschlag: bei gleichem updated-at und anderen Bytes als die Ref-Spitze entscheidet die Herkunft - die Datei ist hinter dem Ref, wenn ihr Blob von dem Ref aus erreichbar ist (frueherer Stand). Das Muster gibt es schon: gitref.Repo.MilestoneHad (rev-list --objects) fuer Meilensteine; isStale ist die eine Definition fuer pull, CheckFresh und Ahead, dort einbauen.
- **2026-10-17 20:54 · berkc** — in-progress (Runde 3, nach testing): isStale ist jetzt eine Syncer-Methode isStale(file, ref, refContent). Ungleiche Stempel entscheiden wie bisher. Bei gleichem updated-at: die Datei ist hinter dem Ref, wenn ihre Bytes andere sind als die der Ref-Spitze und ihr Blob von dem Ref aus erreichbar ist (gitref.Repo.Had, gleiche Logik wie MilestoneHad, beide jetzt ueber had(ref, content)). Wirkt fuer pull, CheckFresh und Ahead zugleich. Kosten auf dem Tick: der Bytevergleich mit dem gecachten Ref-Inhalt (extraEntry.content) kostet nichts; hash-object + rev-list laufen nur bei Gleichstand UND anderen Bytes, also nur fuer genau so eine Datei bis zum pull. Test: TestAFileBehindItsRefWithinTheSameSecondIsStillBehind (ohne Fix rot: show status backlog, note nicht verweigert). repro.sh ohne Pause: note jetzt exit 3, Ref bleibt todo.
- **2026-10-17 20:54 · berkc** — critique Runde 3: nichts. Der Gleichstand wird dort entschieden, wo 'hinter dem Ref' schon definiert ist (isStale), und mit dem vorhandenen Muster (MilestoneHad -> had); keine zweite Definition. Der Doc()==nil-Schutz in isStale bleibt: guenstig, und ein Ticket ohne Doc darf nicht in eine Panik laufen.
