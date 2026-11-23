---
id: 01M851Y9VNEV64RE6FAXRGFJ0Q
title: "refsync.Departed verschluckt den Fehler von Repo.List und liest ihn als 'keine Refs'"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Scheitert das Auflisten der Ticket-Refs, meldet Departed das als Fehler statt so zu tun, als haetten alle gesehenen Tickets das Board verlassen."
context: |-
  Was falsch ist: core/refsync/refsync.go:1969 (func Departed) ignoriert den Fehler von Repo.List und arbeitet mit einer leeren Liste weiter.
  Folge: ein voruebergehender git-Fehler (Lock, kaputter Ref, Timeout) laesst jedes zuvor gesehene Ticket als 'von jemand anderem vom Board genommen' erscheinen. Wer Departed liest (TUI-Hinweis, ForgetDeparted) zeigt dann falsche Meldungen und vergisst Tickets.
  Gefunden vom Bughunt zu DVPAPE am 23.11.2026, gab es schon vorher auf master.
definition-of-done: Departed gibt bei einem Fehler von Repo.List einen Fehler zurueck (oder meldet nichts) statt alle gesehenen Tickets als gegangen; Test mit einem scheiternden Repo.List
tags:
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-11-23T20:35:43Z
updated-at: 2026-11-23T20:50:14Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-49023
claimed-at: 2026-11-23T20:36:54Z
outcome-what: "refsync.Departed returns ([]Departure, error); a failed Repo.List is an error, not an empty ref list. fetch and validate go through new departedHere, which warns on stderr and reports no departures. whoami's file-only count no longer reads a List error as no refs. Release note added."
outcome-why: "A passing git error (lock, broken ref) made every seen ticket look taken off the board elsewhere, so fetch/validate told the person to file or archive live tickets."
outcome-resolves: "DoD 1-4: Departed errors on List failure (refsync test with a broken packed-refs), both callers report no departure and a warning (CLI test), ForgetDeparted untouched as it never reads Departed, whoami sibling fixed. Checked: concurrency — no new shared state, seen record untouched on error; error paths — List error now surfaces; unusual input — broken packed-refs; permissions — n/a."
review-summary: none
review-gaps: "none — checked: departedHere is the one wrapper with two callers (fetch, validate), not a forwarder: it turns the error into the warning; no code orphaned (Departed's old present-map guard replaced, not left); the packed-refs breakage is written in both test packages because refsync_test and cli tests share no helper package — two 4-line copies left as they are"
test-verdict: "pass: go test ./... RC=0, -race on refsync and cli RC=0, DoD 1-4 verified in the tree (both new tests fail with the old Departed); failure path = broken packed-refs after Usable, exercised by the tests; the change writes nothing new"
hunt-verdict: "fail: bughunt change: whoami fileOnlyCount turns a failed listing into a silent 0 — say it on stderr like departedHere; pre-existing seen-record truncated read folded into this ticket (same feature: it loses departures)"
round: 2
back-from: testing
back-by: lane
back-count: 1
---

# refsync.Departed verschluckt den Fehler von Repo.List und liest ihn als 'keine Refs'

## Definition of Done

- [x] Departed gibt bei einem Fehler von Repo.List einen Fehler zurueck (oder meldet nichts) statt alle gesehenen Tickets als gegangen; Test mit einem scheiternden Repo.List
  proof: core/refsync/refsync.go:1983; TestDepartedReportsAFailedRefListingAsAnErrorNotAsDepartures (broken packed-refs)
- [x] Alle Aufrufer von Departed behandeln den Fehler: keine falsche Meldung, kein ForgetDeparted; Test
  proof: internal/cli/fetch.go:249 departedHere used by fetch and validate; TestAFailedRefListingReportsNoDepartures; ForgetDeparted (snapshot.go:61) never reads Departed
- [x] Release-Note in core/release/notes/<ID>.md, falls die TUI dadurch anders meldet
  proof: core/release/notes/RGFJ0Q.md (CLI warning on stderr; TUI does not call Departed)
- [x] whoami: a failed Repo.List no longer counts every ticket file as file-only (same class: list error read as no refs)
  proof: internal/cli/whoami.go:191 returns 0 and warns
- [x] whoami: a failed ref listing is said on stderr, not only counted as 0 file-only
  proof: internal/cli/whoami.go:191; TestAFailedRefListingReportsNoDepartures (fileOnlyCount part)
- [x] refs-seen.json is written atomically (ticket.WriteAtomic), so a fetch reading it while a background fetch writes never sees an empty record and drops a departure
  proof: core/refsync/refsync.go:1255; TestTheSeenRecordIsNeverReadHalfWritten (fails with os.WriteFile)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 20:36 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/48662.sock
- **2026-11-23 20:37 · Alexander Sacharov** — decision: Departed returns ([]Departure, error). On a Repo.List failure it reports no departures and the error; fetch and validate print a warning on stderr and carry on with an empty list (JSON: departed stays []). Why not silent nil: validate is the command run when the board looks wrong, a hidden git error there costs the person the clue. Why not fail the command: departures are an extra report, the rest of fetch/validate is still right. ForgetDeparted (snapshot --drop) never reads Departed, so nothing to change there.
- **2026-11-23 20:42 · Alexander Sacharov** — found out: a packed-refs line git cannot parse makes 'git for-each-ref' exit 128 — the test fault. Broken from the start it also fails 'git remote get-url', so Usable() is false and Departed never runs; the bug bites only once Usable is cached (long-lived syncer, or a lock/ref going bad mid-command). The CLI test therefore drives departedHere on the syncer the last command left in refs. Left alone: refsync.go:847 'entries, _ = y.Box.List()' drops an outbox listing error the same way — outbox, not refs, other subject.
- **2026-11-23 20:42 · Alexander Sacharov** — critique: none. Considered: (1) nil-on-error without a signature change was simpler but hides the git error from validate, the command run when the board looks wrong — kept the error (see decision note). (2) departedHere sits beside strandedHere in fetch.go, same shape; warning text follows share.go's 'jaira: warning:'. (3) whoami returning 0 on a List error states '0 file-only' — accepted: same answer fileOnlyCount already gives on a Paths error, and an inflated count was the louder false claim.
- **2026-11-23 20:45 · Alexander Sacharov** — bughunt: change: internal/cli/whoami.go:187-190 - class: a failed ref listing reaches the reader as a definite answer with no signal (the departedHere rule 'never stands in for nothing silently' is not applied here); example: one ticket file-only, the ref listing fails after refs.Usable() passed (broken packed-refs line) -> 'File only: 0 tickets' and JSON file_only: 0, nothing on stderr, where the true count is 1. Not a blocker: only reachable when the listing breaks between Usable() and List() in one run, and the old code was also wrong there (it overcounted). ran: scratch test in internal/cli (fileOnlyCount before=1, after=0, List err='unexpected line in .git/packed-refs'), removed afterwards
- **2026-11-23 20:45 · Alexander Sacharov** — bughunt: pre-existing: core/refsync/refsync.go:1226-1256 (loadSeen/saveSeen), used by Incoming:1125-1170 and ForgetDeparted:2010 - class: read-modify-write of refs-seen.json with no lock and an in-place os.WriteFile (truncate then write), while a detached 'jaira fetch --json' runs after every command; example: a foreground fetch reads the record mid-write, gets {} (json error -> seen{}), and Incoming saves next without the 'ref gone, file still here' ids it is meant to keep -> that departure is never reported again (lost evidence, the thing Departed exists for). blame: 98009909c (8566KF, 2026-09-10), 2f5713f18 - before this ticket. ran: scratch test, saveSeen in a loop beside loadSeen: 19798 of 20000 reads empty
- **2026-11-23 20:45 · Alexander Sacharov** — hunt report: 2 findings (1 change edge, not a blocker; 1 pre-existing), notes on RGFJ0Q | read whole: core/refsync/refsync.go (the departed/seen/incoming/stranded paths in full, rest scanned), internal/cli/fetch.go, internal/cli/validate.go, internal/cli/whoami.go, core/refsync/refsync_test.go, internal/cli/departed_test.go, core/release/notes/RGFJ0Q.md; also core/gitref List/listNames | not checked: the race in finding 2 end to end through two real jaira processes (no remote in this session), only the record's truncated-read window
- **2026-11-23 20:45 · Alexander Sacharov** — bughunt: done at 23ded2de1 - 2 findings (1 change, 1 pre-existing) in 7 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/refsync/refsync.go core/refsync/refsync_test.go core/release/notes/RGFJ0Q.md internal/cli/departed_test.go internal/cli/fetch.go internal/cli/validate.go internal/cli/whoami.go
- **2026-11-23 20:48 · Alexander Sacharov** — testing: back to in-progress for the bughunt change finding (whoami silent 0). Folded in the pre-existing seen-record finding (fold-into-ticket rule; it is the same feature — a torn read loses the departures Departed reports): fix is saveSeen -> ticket.WriteAtomic. Not a lock: the remaining read-modify-write race only re-reports or re-learns an id, which the next Incoming recomputes from the refs anyway.
