---
id: 01M8543DKBZA0P1CHFGKMMKENA
title: "refsync: der Fehler von Box.List beim Lesen des Outbox wird verschluckt"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Scheitert das Auflisten des Outbox, meldet refsync das als Fehler statt so zu tun, als waere der Outbox leer."
context: |-
  Was falsch ist: core/refsync/refsync.go:847 'entries, _ = y.Box.List()' wirft den Fehler weg und arbeitet mit einer leeren Liste weiter. Ein voruebergehender Lesefehler sieht dann aus wie 'nichts zu senden' - ausstehende Ticket-Aenderungen werden in diesem Lauf still nicht gepusht.
  Gleiche Fehlerart wie RGFJ0Q (Departed verschluckte den Fehler von Repo.List); gefunden im Review zu RGFJ0Q am 23.11.2026, gab es schon vorher.
  Erst nach RGFJ0Q starten: beide aendern core/refsync/refsync.go.
definition-of-done: "Ein Fehler von Box.List wird zurueckgegeben bzw. auf stderr gemeldet, statt als leerer Outbox zu gelten; Test mit einem scheiternden Box.List"
tags:
  - concurrency
blocked-by: []
related: []
follows: 01M851Y9VNEV64RE6FAXRGFJ0Q
commits:
  - 188b7c2a37134c8ecb043c94ad298707ad9aed32
  - 8d09d48a25d98e270c249a5594185da482e30632
created-at: 2026-11-23T21:13:28Z
updated-at: 2026-11-24T07:09:17Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-45309
claimed-at: 2026-11-24T06:57:55Z
outcome-what: Flush in core/refsync returns the error of its second outbox listing instead of treating the outbox as empty; a test makes exactly that listing fail
outcome-why: "an empty listing skipped the merge of a refused milestone write, so the next fetch would overwrite that edit"
outcome-resolves: "MMKENA: Box.List error swallowed in refsync.Flush"
review-summary: "refsync.Flush listet den Outbox nach tidyMilestoneEntries ein zweites Mal; ein Fehler dabei wurde bisher mit _ verworfen und als leere Queue gelesen, sodass die Merge-Logik fuer einen abgelehnten Milestone-Write uebersprungen wurde. Jetzt gibt Flush den Fehler zurueck, bevor irgendetwas gesendet wird - die Eintraege bleiben in der Queue. Alle drei Listings laufen ueber eine Paketvariable listBox, damit ein Test genau den zweiten Aufruf scheitern lassen kann (Hook in export_test.go). Alle vier Aufrufer (cli/refs.go x2, cli/web.go, tui/refs.go) melden den Fehler bereits weiter. Release-Note core/release/notes/MMKENA.md dabei."
review-gaps: "none"
test-verdict: "pass: go vet + go test -race -count=1 ./... green (RC=0), check-ticket-commits ok; DoD 1 verified: TestFlushReportsAFailedSecondOutboxListing passes and fails with the old 'entries, _ =' line restored; DoD 2: core/release/notes/MMKENA.md one line. The second listing fails on its own only via the listBox seam, so the failure path was exercised through that test (error returned, write stays queued, next Flush sends it); the change writes nothing new, no two-writer case"
hunt-verdict: "pass: bughunt done at 188b7c2a3, 0 change, 1 pre-existing filed as WV2J2W"
review-verdict: "accept: beide DoD-Punkte erfuellt; Fehler wird vor dem Senden zurueckgegeben, Queue bleibt erhalten; keine Defekte in Flush gefunden"
review-round: 1
review-verdict-gaps: none
review-check: "1. cd core && go test ./refsync -run TestFlushReportsAFailedSecondOutboxListing -v  2. Erwartet: PASS - der Test laesst das zweite Outbox-Listing scheitern, prueft dass Flush genau diesen Fehler zurueckgibt, der Write in der Queue bleibt und ein zweiter Flush ihn dann sendet  3. cat core/release/notes/MMKENA.md - eine Zeile, beginnt mit '- '. Kein UI-Pfad: das Outbox-Verzeichnis unlesbar zu machen trifft schon das erste Listing, nicht das zweite."
---

# refsync: der Fehler von Box.List beim Lesen des Outbox wird verschluckt

## Definition of Done

- [x] Ein Fehler von Box.List wird zurueckgegeben bzw. auf stderr gemeldet, statt als leerer Outbox zu gelten; Test mit einem scheiternden Box.List
  proof: core/refsync/refsync.go:860 returns the error; core/refsync/refsync_test.go:1024 TestFlushReportsAFailedSecondOutboxListing (fails without the fix)
- [x] Release-Note in core/release/notes/<ID>.md, falls Nutzer eine neue Meldung sehen
  proof: core/release/notes/MMKENA.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 21:59 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/48662.sock
- **2026-11-23 21:59 · Alexander Sacharov** — decision: no open decisions - second Box.List returns its error like the first one in Flush does (refsync.go:829); queue stays untouched, next command retries
- **2026-11-24 06:58 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/3417.sock (restarted worker)
- **2026-11-24 06:58 · Alexander Sacharov** — decision: the second Box.List in Flush (refsync.go:853) can only fail on its own in the window after tidyMilestoneEntries' own List, so no on-disk trick reaches it in a test; a package var listBox in refsync.go, swapped through export_test.go, makes that one call fail. Real effect of the swallowed error was not 'nothing pushed' (Box.Flush lists by itself) but an empty queued map: a refused milestone write skipped mergeRefused and its edit would be lost on the next fetch.
- **2026-11-24 07:00 · Alexander Sacharov** — same pattern left alone, pre-existing and outside refsync: internal/tui/refs.go:67 (unsent marks, redrawn every 2s, a miss is cosmetic) and internal/cli/web.go:173 (runs after Flush, which now already failed on the same listing)
- **2026-11-24 07:01 · Alexander Sacharov** — critique: none. listBox follows the package-var seam already used in internal/cli/flow.go:980 (atATerminal) and internal/cli/self.go:34 (runUpgraded); dropping the second listing instead would mean tidyMilestoneEntries rebuilding the queue after rekeying, more code for the same result
- **2026-11-24 07:03 · Alexander Sacharov** — waiting: bughunt (background subagent) since 08:03
- **2026-11-24 07:03 · Alexander Sacharov** — bughunt: pre-existing: core/refsync/epic.go rebaseEpic (case err == nil: DropEpic after WriteEpic) - class: an outbox entry is dropped by key after a slow step (a push, a git read) without checking it is still the entry that was read, and no lock is held, so a write another session queued for the same key in that window is deleted; also in core/outbox/outbox.go:471 Box.Flush (dropKind after send, every kind) and core/refsync/refsync.go:1039-1041 tidyMilestoneEntries (drop after ReadMilestone, outside the milestone lock mergeRefused takes); for an epic the entry is the only copy of the edit, for a milestone the file is then not pending and the next fetch overwrites it, for a ticket the unsent mark goes and the edit waits for the next write; example: 'jaira epic add' in a second session while the first one's Flush is pushing the same epic; blame: d3bc99ef64, 9cb410a730, 90151d5877 (before the ticket); ran: scratch test in core/refsync with a pre-receive hook that queues a second epic write mid-push, then Flush - the second write was neither queued nor on the remote (test failed as predicted, file removed)
- **2026-11-24 07:03 · Alexander Sacharov** — hunt report: 1 finding, 0 change, 1 pre-existing, as a note on MMKENA; the change itself (listBox seam, error on the second listing, release note) is correct and every Flush caller (internal/cli/refs.go flushRefs and fileOnRefOnly, internal/cli/web.go webPush, internal/tui/refs.go fetchRefs) leaves the queue intact on the new error | read whole: core/refsync (refsync.go, epic.go, schedule.go, export_test.go, the new test), plus outbox.List/readDir/Flush for the call path | not checked: the milestone and ticket variants of the drop race were read, not run; only the epic variant was reproduced
- **2026-11-24 07:03 · Alexander Sacharov** — bughunt: done at 188b7c2a3 - 1 finding (0 change, 1 pre-existing) in 1 module, first hunt; concurrency, error paths, input, permissions checked; files: core/refsync/export_test.go core/refsync/refsync.go core/refsync/refsync_test.go core/release/notes/MMKENA.md
- **2026-11-24 07:04 · Alexander Sacharov** — testing: bughunt pre-existing finding (outbox entry dropped after a slow step, race) is another subject than this goal; filed as WV2J2W
- **2026-11-24 07:04 · Alexander Sacharov** — waiting: review (fresh agent, Herdr tab via run-lane.sh) since 08:04
- **2026-11-24 07:05 · Alexander Sacharov** — review: done at 188b7c2a3
- **2026-11-24 07:06 · Alexander Sacharov** — report: refsync.Flush returns the error of its second outbox listing instead of an empty queue; test via listBox seam; release note | signoff, context unknown | nothing open for this ticket; pre-existing outbox drop race filed as WV2J2W (backlog, needs assignee)
- **2026-11-24 07:07 · Alexander Sacharov** — cockpit: accepted by the person, signoff -> done
