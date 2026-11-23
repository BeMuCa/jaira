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
updated-at: 2026-11-23T20:57:47Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-49023
claimed-at: 2026-11-23T20:36:54Z
outcome-what: "refsync.Departed returns ([]Departure, error); a failed Repo.List is an error, not an empty ref list. fetch and validate go through new departedHere, which warns on stderr and reports no departures. whoami's file-only count no longer reads a List error as no refs. Release note added."
outcome-why: "A passing git error (lock, broken ref) made every seen ticket look taken off the board elsewhere, so fetch/validate told the person to file or archive live tickets."
outcome-resolves: "DoD 1-4: Departed errors on List failure (refsync test with a broken packed-refs), both callers report no departure and a warning (CLI test), ForgetDeparted untouched as it never reads Departed, whoami sibling fixed. Checked: concurrency — no new shared state, seen record untouched on error; error paths — List error now surfaces; unusual input — broken packed-refs; permissions — n/a."
review-summary: "Full review at 860bb94dc. Syncer.Departed now returns ([]Departure, error): when git cannot list the ticket refs it returns the error and no departures, instead of treating 'no refs' as 'every seen ticket left the board'. fetch and validate call it through a new departedHere helper (internal/cli/fetch.go:249) that prints 'jaira: warning: could not check ...' on stderr and reports an empty list, so the rest of the command still runs. whoami's file-only count got the same treatment: a failed listing warns on stderr and counts 0 instead of every file. saveSeen now writes refs-seen.json via ticket.WriteAtomic, so a fetch reading it while a background fetch writes never sees a torn, empty record that would forget departures. Tests: a broken packed-refs line drives the failure in refsync and CLI; a concurrent save/load test covers the torn read. Release note added."
review-gaps: |-
  note: edge: fix-now: internal/cli/fetch.go:242-248 - departedHere was inserted between strandedHere's doc comment and strandedHere, so godoc now gives departedHere both comments and strandedHere none; move the strandedHere comment back above func strandedHere (line 258)
  note: limit: whoami prints file_only 0 (text and JSON) when the listing fails, with only a stderr warning to say the number is not real; chosen on purpose in testing round 2 over an 'unknown' value in the output contract
  note: limit: refs-seen.json is still read-modify-write without a lock (Incoming vs ForgetDeparted); a lost update only re-reports or re-learns one id, accepted in the testing note
  note: elsewhere: core/refsync/refsync.go:847 'entries, _ = y.Box.List()' drops an outbox listing error the same way; there before this change, new ticket
test-verdict: "pass: vet ok, go test -race refsync+cli RC=0 (round 2), cli+release tests RC=0 after the comment/note change, DoD 1-6 verified in the tree"
hunt-verdict: "pass: round 3 hunt at 860bb94dc, 0 findings"
round: 4
back-from: review
back-by: lane
back-count: 1
review-verdict: "accept: all six DoD items are met and proven by tests that fail on the old code; one fix-now (a misplaced doc comment in fetch.go)"
review-round: 3
review-verdict-gaps: |-
  note: edge: fix-now: internal/cli/fetch.go:242-248 - departedHere was inserted between strandedHere's doc comment and strandedHere, so godoc now gives departedHere both comments and strandedHere none; move the strandedHere comment back above func strandedHere (line 258)
  note: limit: whoami prints file_only 0 (text and JSON) when the listing fails, with only a stderr warning to say the number is not real; chosen on purpose in testing round 2 over an 'unknown' value in the output contract
  note: limit: refs-seen.json is still read-modify-write without a lock (Incoming vs ForgetDeparted); a lost update only re-reports or re-learns one id, accepted in the testing note
  note: elsewhere: core/refsync/refsync.go:847 'entries, _ = y.Box.List()' drops an outbox listing error the same way; there before this change, new ticket
review-check: "1. cd into this worktree  2. go test ./core/refsync -run 'TestDepartedReportsAFailedRefListingAsAnErrorNotAsDepartures|TestTheSeenRecordIsNeverReadHalfWritten' -count=1 -v  3. both PASS  4. go test ./internal/cli -run TestAFailedRefListingReportsNoDepartures -count=1 -v  5. PASS: no departures and the warning 'could not check for tickets taken off the board elsewhere' on a broken packed-refs  6. no hand path: the failure only happens when git breaks between two calls in one command"
stop-after-fix: "true"
---

# refsync.Departed verschluckt den Fehler von Repo.List und liest ihn als 'keine Refs'

## Definition of Done

- [x] Departed gibt bei einem Fehler von Repo.List einen Fehler zurueck (oder meldet nichts) statt alle gesehenen Tickets als gegangen; Test mit einem scheiternden Repo.List
  proof: core/refsync/refsync.go:1986; TestDepartedReportsAFailedRefListingAsAnErrorNotAsDepartures (broken packed-refs)
- [x] Alle Aufrufer von Departed behandeln den Fehler: keine falsche Meldung, kein ForgetDeparted; Test
  proof: internal/cli/fetch.go:246 departedHere used by fetch and validate; TestAFailedRefListingReportsNoDepartures; ForgetDeparted (snapshot.go:61) never reads Departed
- [x] Release-Note in core/release/notes/<ID>.md, falls die TUI dadurch anders meldet
  proof: core/release/notes/RGFJ0Q.md (CLI warning on stderr; TUI does not call Departed)
- [x] whoami: a failed Repo.List no longer counts every ticket file as file-only (same class: list error read as no refs)
  proof: internal/cli/whoami.go:192 returns 0 and warns
- [x] whoami: a failed ref listing is said on stderr, not only counted as 0 file-only
  proof: internal/cli/whoami.go:192; TestAFailedRefListingReportsNoDepartures (fileOnlyCount part)
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
- **2026-11-23 20:50 · Alexander Sacharov** — critique round 2: none — the whoami warning threads cmd.ErrOrStderr() through boardState like departedHere takes it, rather than writing to os.Stderr, so the test can read it; WriteAtomic is the store's existing helper, no second one.
- **2026-11-23 20:51 · Alexander Sacharov** — bughunt: change: internal/cli/whoami.go:192 - class: a count that could not be worked out is still printed as a real number on stdout; the warning goes only to stderr and nothing on stdout marks the value as unknown; example: Repo.List fails between refs.Usable() and List() (a transient git lock) -> fileOnlyCount returns 0, whoami prints 'File only: 0 tickets on this disk and on no ref' (whoami.go:106) and --json emits "file_only": 0 (whoami.go:74) with no field saying it is unknown, so a JSON consumer that ignores stderr reads 'none'. This contradicts the line added to core/release/notes/RGFJ0Q.md ('instead of reporting a file-only count it could not work out') and the code comment at whoami.go:189-190 ('counted as 0 it would claim none is'). Not a blocker: the value is warned about, the window is narrow. Fix direction: carry 'unknown' through boardGit (e.g. -1 / omit file_only in JSON, print 'unknown' in text) or reword the release note; ran: scratch test with broken packed-refs before whoami (board then reads as ref mode no, so the path is only reachable transiently); read: departed_test.go:61 confirms fileOnlyCount returns 0 on a failed listing
- **2026-11-23 20:52 · Alexander Sacharov** — hunt report: 1 finding, 1 change (edge, not blocker), as a note on RGFJ0Q | read whole: core/refsync/seen_test.go (new); read as a diff since 23ded2de1: core/refsync/refsync.go (saveSeen/loadSeen, Incoming, Departed, ForgetDeparted followed), internal/cli/whoami.go, internal/cli/departed_test.go, core/release/notes/RGFJ0Q.md; also read ticket.WriteAtomic | could not reproduce the whoami path end to end through the CLI: a ref listing broken from the start makes the board read as ref mode no, so only a transient failure between Usable and List reaches it; the lost-update race on refs-seen.json (Incoming vs ForgetDeparted) was checked and only re-raises a forgotten departure, which the code accepts as a repeated notification
- **2026-11-23 20:52 · Alexander Sacharov** — bughunt: done at a0a66de74 - 1 finding (1 change, 0 pre-existing) in 8 modules, follow-up hunt; concurrency, error paths, input, permissions checked; files: core/refsync/refsync.go core/refsync/refsync_test.go core/refsync/seen_test.go core/release/notes/RGFJ0Q.md internal/cli/departed_test.go internal/cli/fetch.go internal/cli/validate.go internal/cli/whoami.go
- **2026-11-23 20:53 · Alexander Sacharov** — testing round 2: back for the hunt's whoami edge. Chose the smaller fix the hunter offered — reword the release note and the comment to say what happens (warning on stderr, count reads 0) — over an 'unknown' value in boardGit and the JSON: that would change whoami's output contract for a window only a git error between Usable() and List() opens; edges growing on a sibling ask for a smaller design, not a bigger fix.
- **2026-11-23 20:54 · Alexander Sacharov** — waiting: bughunt round 3 (subagent) since 21:54
- **2026-11-23 20:54 · Alexander Sacharov** — hunt report: 0 findings, nothing written beyond this report | read as a diff since a0a66de74: internal/cli/whoami.go (comment only, fileOnlyCount read in full with its caller), core/release/notes/RGFJ0Q.md (note now matches the 0-plus-warning behaviour, checked against departed_test.go:61); the other 6 files unchanged since the last round | nothing left unchecked
- **2026-11-23 20:54 · Alexander Sacharov** — bughunt: done at 860bb94dc - 0 findings in 8 modules, follow-up hunt; concurrency, error paths, input, permissions checked; files: core/refsync/refsync.go core/refsync/refsync_test.go core/refsync/seen_test.go core/release/notes/RGFJ0Q.md internal/cli/departed_test.go internal/cli/fetch.go internal/cli/validate.go internal/cli/whoami.go
- **2026-11-23 20:55 · Alexander Sacharov** — waiting: review lane (fresh agent via run-lane.sh) since 21:55
- **2026-11-23 20:56 · Alexander Sacharov** — review: done at 860bb94dc
- **2026-11-23 20:56 · Alexander Sacharov** — review: fix-now: internal/cli/fetch.go:242-248 - strandedHere's doc comment now sits above departedHere's; move it back above func strandedHere
- **2026-11-23 20:57 · Alexander Sacharov** — review fix-now closed: strandedHere's doc comment moved back above it (fetch.go). Set stop-after-fix=true: a comment move, no lock, no parser. Proofs on DoD 1,2,4,5 refreshed to current lines. The review's other notes (whoami 0 limit, seen record lock limit, Box.List elsewhere) stay as notes for the person.
