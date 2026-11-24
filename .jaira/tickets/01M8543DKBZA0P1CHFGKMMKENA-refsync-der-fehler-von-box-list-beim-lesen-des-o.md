---
id: 01M8543DKBZA0P1CHFGKMMKENA
title: "refsync: der Fehler von Box.List beim Lesen des Outbox wird verschluckt"
status: critique
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
commits: []
created-at: 2026-11-23T21:13:28Z
updated-at: 2026-11-24T07:01:06Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-45309
claimed-at: 2026-11-24T06:57:55Z
outcome-what: Flush in core/refsync returns the error of its second outbox listing instead of treating the outbox as empty; a test makes exactly that listing fail
outcome-why: "an empty listing skipped the merge of a refused milestone write, so the next fetch would overwrite that edit"
outcome-resolves: "MMKENA: Box.List error swallowed in refsync.Flush"
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
