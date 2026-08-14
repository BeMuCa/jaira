# Handoff — 2026-08-15

State after the session that made commits and blocked-reasons mandatory lane
output, made `follows:` real, and gave the open ticket a viewport. Written
before a context clear; this file is the memory. The previous handoff
(2026-08-13, Phase 5 close + settings/lanes/notes rework) is in git history at
`d0a1b40` — its decisions are restated below where they still bind.

## Where things are

- **Public: https://github.com/BeMuCa/jaira** — branch `master`, in sync with
  origin; the session's last code commit is `333e70f`, this handoff sits on top.
- `go test ./... -race` passes on all 11 packages; `go vet` clean.
- The user's binary lives at `~/.local/bin/jaira`. **`go install` does not update
  it** — it writes `~/go/bin`. Use
  `go build -o ~/.local/bin/jaira ./cmd/jaira`. Rebuilt this session.

## What landed this session (4 commits)

**`e426e73` — commits are mandatory implementing-lane output.**
Trigger: real feedback — three accepted tickets (DAHC06, YDACKQ, PJFVD1 on the
user's German-language "requirement" board) carried `commits: (keine)`, so at
sign-off there was no way to see what was being accepted. The mechanism already
existed end to end (`FieldCommits`, `fieldFilled`, the `--from-lane` decoder);
the only change was adding `commits` to `20-in-progress.md`'s `output-produces`
plus a prompt paragraph telling the agent to pass
`--commits "$(git rev-parse HEAD)"`. Both move paths verified on the running
binary: flags path and JSON-on-stdin path each refuse with exit 3 without
commits and pass with them.

**`c605edb` — `follows:` is visible, settable and durable.**
It was written by exactly one code path (the sign-off screen's follow-up
action) and rendered nowhere — a write-only field. Now: `jaira create
--follows <handle>` (resolves the prefix, refuses a dead link with exit 5,
normalises to the full id), rendered in `jaira show`, the TUI detail pane and
`--json` (`follows`), placed canonically after `blocked-by`. The sign-off
follow-up also writes the predecessor's commits into the context prose
("That work shipped in <sha>.") so the answer to "what was already done"
survives the predecessor being archived.

**`5a47359` — the open ticket is clipped to the terminal and scrolls.**
`renderDetail` never read `m.height`; a 41-line ticket rendered whole into a
16-line terminal, pushing handle/title/goal off the top with no key to bring
them back. Now: `detailScroll` offset, clamped in the renderer (only it knows
the content length), `ctrl+d`/`ctrl+u`/`pgdown`/`pgup` to scroll, a
`+N more ·` prefix on the footer when clipped, reset on open. `j`/`k` still
switch tickets — untouched. Board-column scrolling was *not* the bug: verified
correct across a 9×8 size matrix before touching anything.

**`333e70f` — the blocked lane refuses a ticket that cannot say what it is
waiting on.** New field `blocked-reason`, new lane flag
`requires-blocked-reason` (on `60-blocked.md`), new violation
`needs_blocked_reason` (exit 3), `jaira move --reason "…"`. `blocked-by` is
accepted in place of a typed reason. Rendered as `waiting on` in show/TUI —
**only while the ticket sits in a parking lane**; the field stays on the
ticket as history after unparking. Editable in the TUI field editor.

## Decisions made this session — do not silently reverse

1. **Parking is exempt from the dependency check and from the leaving lane's
   output contract** (`core/gate/gate.go`, the `parking` flag = target lane has
   `requires-blocked-reason`). Both were found the hard way: the dependency
   gate locked tickets with open blockers *out of Blocked* (exit 4 — the very
   tickets the lane is for), and the new commits contract made a ticket stopped
   mid-work with no commits yet unparkable. Tests:
   `TestUnresolvedBlockerDoesNotBarTheBlockedLane`,
   `TestParkingIsExemptFromTheLeavingLanesContract`.
2. **`blocked-reason` is its own field, not a reuse of `question`.** Chosen by
   the user over the one-line `requires-question: true` alternative, because a
   ticket wandering HITL→Blocked would carry its stale question as a fake
   reason, and `requires-question` also waives the owner guard.
3. **A stale reason is not rendered on an active ticket.** `waiting on` shows
   only while the current lane requires it; deleting the field on unpark was
   rejected in favour of keeping history.
4. **`create --follows` validates; `jaira set follows=` does not.** Deliberate:
   `set` validates no field, and making `follows` its lone special case is
   inconsistency in the other direction. Known consequence: `jaira set <id>
   follows=NOPE99` writes a dead link, exit 0. Reopening this is a user
   decision.
5. **Archives are cheap; deletion buys nothing.** Measured: 1000 archived
   tickets ≈ 320 KB in git after gc; the board scan never reads `archive/`
   (sibling of `tickets/`, `store.go`). Deleting a committed file does not
   shrink history. Any future "cleanup" feature should be a listing/filter
   problem, not a deletion cascade.

The five standing decisions from 2026-08-13 (precedence is a merge rank; lanes
are files; a progress note is never gated; `.jaira/lanes/` is authoritative;
`create` lands in backlog) still bind — see `d0a1b40` for their full text.

## The lane-copy trap, again

`20-in-progress.md` and `60-blocked.md` both changed. **Every project copy of
those lanes now warns as modified and does NOT have the new gates** until
re-copied. This bit once before (sign-off → human review rename). If the
user's "requirement" board has its own lane files, the commits gate is not
active there — that board is where the original complaint came from.

## Open, needs the user

- **DAHC06, YDACKQ, PJFVD1 still record no commits.** The gate only works
  forward. The 07.08 commits need manual attribution:
  `jaira set <handle> commits=<sha1>,<sha2>` (commits is a list field).
- **Follow-up chains: archive protection / cascade.** The user sketched:
  warn (but allow) archiving a predecessor with an open follow-up; cascading
  archive along the `follows` chain. Deferred until the links are visible —
  which they now are. Note: it is *archive*, reversible, not delete.
- **The TUI move dialog does not prompt for a reason** when moving to Blocked;
  it shows the gate message and points at the field editor — same contract as
  `question` for HITL. Acceptable parity or a papercut worth a prompt?
- Per-ticket `## Options` earning their keep — unchanged from last handoff.

## Known and deliberately not touched

- `gofmt` flags `core/gate/gate.go`, `internal/cli/tickets.go`,
  `internal/tui/lanes.go`, `internal/tui/view.go` — **pre-existing at
  `d0a1b40`**, verified per file against HEAD before this session's edits; the
  complaints sit in untouched map literals/alignment groups. Session edits are
  clean. Reformatting them would pollute blame for no behaviour change.
- Board render has a floor of 10 lines: terminals ≤ 9 rows overflow by ~2
  lines regardless of selection. Found while chasing the scroll bug; cosmetic;
  not fixed.
- `jaira archive` (no args) lists everything unfiltered — the real "archive
  gets full" cost is this listing, not storage.
- No way to forget a board; `~/.jaira/state/` grows unbounded — unchanged from
  last handoff.
