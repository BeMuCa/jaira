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
updated-at: 2026-11-23T20:42:25Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-49023
claimed-at: 2026-11-23T20:36:54Z
outcome-what: "refsync.Departed returns ([]Departure, error); a failed Repo.List is an error, not an empty ref list. fetch and validate go through new departedHere, which warns on stderr and reports no departures. whoami's file-only count no longer reads a List error as no refs. Release note added."
outcome-why: "A passing git error (lock, broken ref) made every seen ticket look taken off the board elsewhere, so fetch/validate told the person to file or archive live tickets."
outcome-resolves: "DoD 1-4: Departed errors on List failure (refsync test with a broken packed-refs), both callers report no departure and a warning (CLI test), ForgetDeparted untouched as it never reads Departed, whoami sibling fixed. Checked: concurrency — no new shared state, seen record untouched on error; error paths — List error now surfaces; unusual input — broken packed-refs; permissions — n/a."
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
  proof: internal/cli/whoami.go:188 returns 0 like the Paths error above it

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 20:36 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/48662.sock
- **2026-11-23 20:37 · Alexander Sacharov** — decision: Departed returns ([]Departure, error). On a Repo.List failure it reports no departures and the error; fetch and validate print a warning on stderr and carry on with an empty list (JSON: departed stays []). Why not silent nil: validate is the command run when the board looks wrong, a hidden git error there costs the person the clue. Why not fail the command: departures are an extra report, the rest of fetch/validate is still right. ForgetDeparted (snapshot --drop) never reads Departed, so nothing to change there.
- **2026-11-23 20:42 · Alexander Sacharov** — found out: a packed-refs line git cannot parse makes 'git for-each-ref' exit 128 — the test fault. Broken from the start it also fails 'git remote get-url', so Usable() is false and Departed never runs; the bug bites only once Usable is cached (long-lived syncer, or a lock/ref going bad mid-command). The CLI test therefore drives departedHere on the syncer the last command left in refs. Left alone: refsync.go:847 'entries, _ = y.Box.List()' drops an outbox listing error the same way — outbox, not refs, other subject.
