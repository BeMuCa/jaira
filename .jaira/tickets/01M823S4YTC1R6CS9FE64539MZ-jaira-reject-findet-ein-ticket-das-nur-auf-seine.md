---
id: 01M823S4YTC1R6CS9FE64539MZ
title: "jaira reject findet ein Ticket, das nur auf seinem Ref liegt, und lehnt es ohne vorheriges pull ab"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "'jaira reject <id> --why ... --force' funktioniert in jedem Checkout, auch wenn das Ticket nur auf seinem Ref liegt oder die lokale Ref-Kopie veraltet ist: es holt den Ref, lehnt ab und schreibt das Ergebnis so, dass alle es sehen - ohne dass die Person einen Branch oder ein pull wählen muss."
context: |-
  Was falsch ist (Alex 22.11.2026, HS3K1E):
  1. 'jaira reject HS3K1E --force' im Haupt-Checkout: 'ticket: not found: HS3K1E' - die lokale Ref-Kopie war veraltet; erst 'jaira fetch' half.
  2. Danach: 'ticket: not found: 01M7X84R87N1Z4K96A39HS3K1E' - reject verlangt die Datei auf der Platte; das Ticket lag nur auf refs/jaira/tickets/<ULID>.
  3. Behelf: in einem fremden Worktree (jaira-bwt53t) 'jaira pull HS3K1E && jaira reject ...' - das Archiv-File reist nun im PR eines anderen Tickets mit.
  Eine Person will nur sagen 'nicht machen, weil ...'. Wo das Ticket liegt, ist nicht ihre Sache.
  Code: internal/cli (reject/archive), core/refsync (Ref lesen/schreiben). Ähnlich: 'jaira dod/note' auf Ref-only-Tickets verlangen heute auch erst ein pull.
definition-of-done: "reject auf ein Ticket, das nur auf dem Ref liegt: holt den Ref (auch wenn die lokale Kopie veraltet ist), lehnt ab und schreibt die Ablehnung auf den Ref, ohne Datei in einem Branch zu hinterlassen; alle anderen sehen es nach fetch (Test)"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-11-22T17:10:08Z
updated-at: 2026-11-22T17:19:12Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2431680
claimed-at: 2026-11-22T17:10:39Z
outcome-what: "jaira reject and R on the board reject a ticket that lives only on its ref: refsync.Syncer.Reject fetches, stamps the rejection (shared ticket.StampRejection) and pushes under lease; no file is written in the checkout. Store.List hides ref tickets carrying rejected-kind; Load/show still find them. CLI fetches once when an id is unknown locally and names already_rejected / filed_away / race / not_found with exit codes."
outcome-why: "A person rejecting HS3K1E had to fetch, then pull it into an unrelated worktree, so the archive file travelled in another ticket's PR. Where a ticket lives is not the rejecter's business."
outcome-resolves: "DoD 1: TestRejectATicketThatLivesOnlyOnAStaleRef + TUI test; DoD 2: TestRejectATicketWithAFileHereStillArchivesIt; DoD 3: TestRejectOnTheRefRefusals (unknown, already rejected on ref and in archive, archived, unreachable remote, real lost push via pre-push hook) + handle/ulid; DoD 4: core/release/notes/4539MZ.md. Checked: concurrency (lease + race test, TUI off update loop), error paths (nothing written on any refusal, asserted), input (handle/ULID/unknown), permissions (--force gate unchanged, runs before any write)."
---

# jaira reject findet ein Ticket, das nur auf seinem Ref liegt, und lehnt es ohne vorheriges pull ab

## Definition of Done

- [x] reject auf ein Ticket, das nur auf dem Ref liegt: holt den Ref (auch wenn die lokale Kopie veraltet ist), lehnt ab und schreibt die Ablehnung auf den Ref, ohne Datei in einem Branch zu hinterlassen; alle anderen sehen es nach fetch (Test)
  proof: internal/cli/rejectref_test.go TestRejectATicketThatLivesOnlyOnAStaleRef; internal/tui/rejectref_test.go TestRejectACardThatLivesOnlyOnItsRef; core/refsync/refsync.go Syncer.Reject
- [x] reject auf ein Ticket mit Datei im aktuellen Branch verhält sich wie heute (Archiv im Branch) (Test)
  proof: internal/cli/rejectref_test.go TestRejectATicketWithAFileHereStillArchivesIt; TestRejectWritesTheReasonAndArchives
- [x] Eingabeklassen: kurze ID und ULID, kein Netz/kein Remote, Ref beim Push überholt (Rennen), Ticket schon abgelehnt, Ticket unbekannt - jede mit klarer Meldung und Exit-Code (Tabellentest)
  proof: internal/cli/rejectref_test.go TestRejectOnTheRefRefusals (+ handle/ulid subtests of TestRejectATicketThatLivesOnlyOnAStaleRef)
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/4539MZ.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 17:10 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 17:12 · Alexander Sacharov** — decision (worker, settled, no open question): ref-only reject = synchronous fetch+read+stamp+write-with-lease on the ref, like pull/release (core/refsync). The rejected ticket stays on its ref with rejected-* fields; Store.List hides ref tickets carrying rejected-kind, Load/show still find them. Race lost = exit 3 reason race, no auto-retry. Not found locally + usable remote = fetch once and look again (fixes stale-copy case). Already filed away rejected = exit 3 already_rejected. Out of scope: 'jaira restore' of a ref-only rejection (restore takes a file), pull of a rejected ref.
- **2026-11-22 17:19 · Alexander Sacharov** — found out: TUI forbids network in the update loop (model.go New comment), so R on a ref-only card runs Syncer.Reject as a tea.Cmd (refRejectedMsg). Race test reproduces a real lost push with a pre-push hook that pushes another commit onto the ref first. Siblings left out on purpose (different commands, not this goal): archive, delete, dod, note on ref-only tickets still say 'pull it first'. A ticket pulled by someone else and then rejected on the ref keeps their file; their next write sees the ref ahead.
