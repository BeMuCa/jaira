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
updated-at: 2026-11-22T17:27:41Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2431680
claimed-at: 2026-11-22T17:10:39Z
outcome-what: "jaira reject and R on the board reject a ticket that lives only on its ref: refsync.Syncer.Reject fetches, stamps the rejection (shared ticket.StampRejection) and pushes under lease; no file is written in the checkout. Store.List hides ref tickets carrying rejected-kind; Load/show still find them. CLI fetches once when an id is unknown locally and names already_rejected / filed_away / race / not_found with exit codes."
outcome-why: "A person rejecting HS3K1E had to fetch, then pull it into an unrelated worktree, so the archive file travelled in another ticket's PR. Where a ticket lives is not the rejecter's business."
outcome-resolves: "DoD 1: TestRejectATicketThatLivesOnlyOnAStaleRef + TUI test; DoD 2: TestRejectATicketWithAFileHereStillArchivesIt; DoD 3: TestRejectOnTheRefRefusals (unknown, already rejected on ref and in archive, archived, unreachable remote, real lost push via pre-push hook) + handle/ulid; DoD 4: core/release/notes/4539MZ.md. Checked: concurrency (lease + race test, TUI off update loop), error paths (nothing written on any refusal, asserted), input (handle/ULID/unknown), permissions (--force gate unchanged, runs before any write)."
review-summary: "none"
round: 3
back-from: testing
back-by: lane
back-count: 1
review-gaps: "removed the nil-refSync branch in internal/tui/model.go rejectOnRef (New and switchBoard always set it); left Syncer.Reject's fetch-read-parse prologue, which repeats Pull/Release (pre-existing pattern, three callers now — a shared helper would touch both older functions, not this ticket's scope); kept the localPath and Pending guards in Syncer.Reject as API guards for any later caller"
test-verdict: "pass: go test ./... green, -race -count=1 on core/ticket core/refsync internal/tui internal/cli green, go vet linux+windows green, shipped-history ok. DoD 1-4 re-run by their named tests. Binary in a two-clone fixture: reject by handle from a clone that never fetched -> rc 0, no file under .jaira, creator's list drops it after fetch; repeat -> rc 3 already_rejected; two parallel rejects -> one rc 0, other rc 3 race with nothing written; unknown id -> rc 5."
hunt-verdict: "fail: 3 change findings (blocker: Syncer.Reject ignores assignee; pull revives a ref-rejected ticket with no restore path; fetch error swallowed -> exit 5 instead of 1)"
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
- **2026-11-22 17:19 · Alexander Sacharov** — critique: rejectOnRef catch-all mislabels non-network errors as no_remote; narrow to gitref.ErrOffline/ErrNoRepo, return others unchanged. Left alone: ErrUnsent guard in Syncer.Reject (protects a queued write's lease, cheap), TUI duplicate of the on-ref path (different surface: async tea.Cmd).
- **2026-11-22 17:20 · Alexander Sacharov** — critique round 2: the one finding (catch-all no_remote) is fixed in ebed05124 — only ErrOffline/ErrNoRepo get that message. Nothing else raised.
- **2026-11-22 17:24 · Alexander Sacharov** — bughunt: change: blocker: core/refsync/refsync.go:1563 (Syncer.Reject) - class: reject on the ref ignores the ref's assignee, so a ticket another clone has pulled and is working gets rejected under its owner; Pull and Release refuse that with ErrTaken (refsync.go:1422/1519), Reject has no holder check; same path from the board via internal/tui/model.go:2403 (rejectOnRef). Before this change the only route was 'jaira pull', which refused. example: ada pulls 01JATN and notes on it, bob (no file) runs 'jaira reject 01JATN --why x --force' -> exit 0; ada's board keeps the card as 'behind ref: pull it', her 'jaira note' is refused as stale (exit 3), and after 'jaira pull' the card is an ordinary active card carrying rejected-kind/why/by: bob - she can work it on to done with the rejection on it. ran: two scratch clones on one bare remote with the branch's binary
- **2026-11-22 17:24 · Alexander Sacharov** — bughunt: change: core/ticket/store.go:1016 + core/refsync/refsync.go:1394 (Pull) - class: a ticket rejected on its ref comes back onto a board through pull/claim still carrying its rejection, and nothing takes the rejection off: List hides only ref-only rejected tickets, Pull/refresh copy the ref's rejected-* fields into tickets/ and the card shows as live; 'jaira restore' only knows archive/logbook files, so there is no way to undo a ref rejection that clears the fields (unreject runs only in Restore). example: bob rejects SR2WF2 on the ref, ada runs 'jaira pull SR2WF2' (exit 0, 'is yours') and 'jaira claim SR2WF2' (exit 0) -> active card with rejected-kind: wontfix; 'jaira restore SR2WF2' -> 'not in the archive or in .jaira/logbook/'. marked change though Pull is older (RFC7GA): an archived-rejected ticket's filed ref could reach Pull before, but this change makes rejection-on-the-ref the normal state of a rejected unpulled ticket and hides it from List, so pull is now the only way back and it is silent. ran: same scratch clones
- **2026-11-22 17:24 · Alexander Sacharov** — bughunt: change: internal/cli/reject.go:134 (loadForReject) - class: any fetch failure is swallowed and the lookup then answers exit 5 'not on the board, not on a ref, not filed away', though COMMANDS.md says exit 1 when the remote cannot be reached and exit 5 only when nothing anywhere has that id; also hides auth/other fetch errors. example: a ticket created in another clone after this one last fetched, remote unreachable -> 'jaira reject 8MAYFV --why x --force' exit 5 reason not_found, while the ticket exists on the remote. Edge, not a blocker: nothing is written and a retry online works. ran: remote set-url to a missing path, reject (exit 5, --json reason not_found)
- **2026-11-22 17:24 · Alexander Sacharov** — bughunt: done at d67f1783b - 3 findings (3 change, 1 of them blocker; 0 pre-existing) in 5 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/refsync/refsync.go core/release/notes/4539MZ.md core/ticket/store.go docs/COMMANDS.md internal/cli/reject.go internal/cli/rejectref_test.go internal/tui/model.go internal/tui/rejectref_test.go
- **2026-11-22 17:27 · Alexander Sacharov** — hunt round 1 fixed: (1) Syncer.Reject refuses a ticket whose ref assignee is someone else (ErrTaken -> exit 3 taken), as Pull/Release do; (2) chose the smaller way back over a ref-aware restore: Pull of a ref-rejected ticket takes the rejection off via shared ticket.Unreject (also used by restore's unreject) and says so in output + note; (3) loadForReject: fetch failed and id unknown -> exit 1 no_remote, not 5. Sibling of the ticket's fault fixed too: 'jaira pull <handle>' in a clone that had not fetched failed with 'no ref' (resolveID fell back to the bare handle) - pull now fetches when the id is unknown locally. Not fixed (seen by hunter, not reproduced): TUI async reject vs periodic fetchRefs lock clash would fail without writing.
