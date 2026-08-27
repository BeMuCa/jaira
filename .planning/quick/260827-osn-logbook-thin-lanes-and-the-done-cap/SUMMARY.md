---
quick_id: 260827-osn
slug: logbook-thin-lanes-and-the-done-cap
date: 2026-08-27
status: complete
---

# The logbook, thin lanes, and where the done cap stopped

Four of his items from one message, in the order he gave them. Two tickets
were driven through the board's own lanes rather than fixed in place, as the
test of whether a session can run the pipeline — see below.

## `jaira sync` is `jaira logbook` (`b93a85c`, pushed)

Unreleased (PR #5, after v0.1.0), so a rename with no migration owed. `sync`
implied a server and collided with `sync-tasks`. `.jaira/sync/` →
`.jaira/logbook/`, `Store.Sync` → `Store.Logbook`, JSON `synced` → `logged` /
`logbook`. `Restore` still reads the old folder name — one exists in the
wild, `alsa-20260823` on his req board — and a test builds one by hand.

`archive --help` now says it stamps commits (it has since PR #5) and says in
one sentence which command is for what: archive for a ticket not being
worked, from any lane; logbook for finished work, terminal lane only.

## `QPJNQP` — z draws an empty lane thin (`bd5b2ee`, `f629024`, `0c5423c`)

His words: hiding can confuse. An empty lane stays on the board at four
cells, its name down the column one letter per line, no count. Lanes with
tickets take the room. At 200 columns eight full lanes fit; thin, all ten.

`laneWindow(idx, n, perScreen)` became `fitWindow(idx, costs, budget)` —
grows from the focused lane outward, left then right, while the next lane
fits. Unit costs reproduce the old centred, clamped row exactly; the sweep in
`centre_test.go` proves it. Four mechanisms that existed only because lanes
became invisible are gone: the keep-rule, cursor relocation on toggle, the
h/l skip, the "hidden" notice.

**Measured on the way:** lipgloss v2's `Width` is the whole box, border
included (`Width(22)` renders 22 cells). The full columns have always been
budgeted two cells wider than they render. Left alone; the comment at the
one place this touches now says what the two cells are. `Width(3)` wraps
`" B"` and doubles every line — which is why thin is four.

## The pipeline test: it ran, but it was me

Both tickets went in-progress → critique → in-progress → critique → optimize
→ review → signoff, with every gate honoured (`critique` refused nothing;
`human` refused a move without `--question`; `signoff` refused an agent
exit on `--dry-run`). Critique found one thing each time and the loop closed
on the second pass with `review-summary=none`.

But every step was this session typing the commands because it was told to
work the board. No second model judged anything — subagents were off. What
this proves is that the mechanism holds; what it does not prove is that a
fresh session picks the board up unprompted. That test is still: open a new
session on a board with work in `critique`, say "go", watch.

## `SGPDYK` — planned, at `human`

The done-lane cap. His two answers recorded: on the move into done, and it
says so. Planned in five steps. Found while planning: **there is no shared
move function** — CLI (`flow.go:207`), TUI (`model.go:1291`) and the forced
TUI path (`:1343`) each write the move themselves, so a trim has three call
sites or the three get a common function first, which is its own cut.

Three decisions handed to him with a recommendation each: trim into the
logbook rather than the archive (finished is finished, and it then counts in
the per-day stat); declare the cap as `holds: 10` on the done lane file
rather than in code; oldest = smallest `updated-at`, there being no field
for when a ticket entered its lane.

## Also this session

- `human` is marked "a person's, never out" in the managed block but only
  `signoff` carries `requires-human-exit`; the gate lets an agent leave
  `human`. Proved with `--dry-run` on this board. Not yet a ticket.
- His req board tracks `.jaira/lanes/{critique,optimize,review}.md`, `order`
  and `removed` in git — his personal board structure is in the shared repo.
  `jaira share` was never run there (the `/.jaira/` ignore line is commented
  out by hand), so `AddLanesIgnore` never fired. His call.
- `adopt` takes any lane file, not only a shared one — verified against a
  temp catalogue. Only its help text says otherwise.

## Verified

- `go test ./... -race`, cache cleared, after every commit: all packages ok
- `gofmt -l core internal`: `internal/cli/tickets.go` only
- `~/.local/bin/jaira` rebuilt at each step; `jaira sync` exits 2, `jaira
  logbook --help` reads as intended; a 200×24 render dump with z shows all
  ten lanes with vertical names
