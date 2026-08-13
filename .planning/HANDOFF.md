# Handoff — 2026-08-13 (evening)

State after the session that closed Phase 5 and reworked the board's settings,
lanes and notes. Written before a context clear; this file is the memory.

## Where things are

- **Public: https://github.com/BeMuCa/jaira** — branch `master`, everything
  pushed, CI green on linux/macOS/windows.
- ~75 commits landed today. `go test ./... -race` passes locally on all
  packages; CI runs the same on three operating systems.
- The user's binary lives at `~/.local/bin/jaira`. **`go install` does not update
  it** — it writes `~/go/bin`. Use
  `go build -o ~/.local/bin/jaira ./cmd/jaira`.

## What landed today

Data and format
- A hand-written `context: |` block scalar used to read back as the string `"|"`,
  losing the whole text. Block scalars are now read and written properly.
- `## Description` is gone from the ticket template; nothing ever read it. The
  problem statement lives in `context`, which every lane's bounded input carries.

Traceability
- The review lane must produce `review-summary` and `review-gaps`; the gate
  refuses to release a ticket from review without them.
- Every definition-of-done and plan item can carry a `proof:` sub-line, written
  with `jaira dod <id> <n> --done --proof "…"`.
- The detail pane shows the full ticket id (`y` copies it over OSC52) and a
  `when` row with created / last-touched.

Gates and lanes
- The promotion gate fires on entering the *specified zone* — a lane declaring
  `requires-specified`, which the built-in `todo` carries — rather than on
  leaving a lane literally named `backlog`. That made a thinking lane possible.
- A `brainstorm` built-in lane ships, optional per ticket
  (`requires-option: brainstorm`), and cannot be left until it produced a `goal`.
- Built-in lanes may be overridden, including dropping a protection; two distinct
  warnings fire, one naming the protection lost.
- The `signoff` lane displays as **Human Review**; its id is unchanged.
- A project's column order lives in `.jaira/lanes/order` as positions from 1;
  moving a lane shifts its neighbour. `precedence` is untouched by it.

Surfaces
- `S` opens settings, holding the lane screen and the default board. `L` is gone.
- The lane settings screen is a small board: `h`/`l` navigate, `H`/`L` move the
  selected lane, `x` removes it from the project, a `+` column adds one from the
  catalogue.
- `jaira lanes` gained `use`, `publish`, `adopt`, `add`, `remove`, `move`,
  `default`, `show`, `path`, `template`, `shared` — so an agent can do everything
  the screen can.
- Boards are listed on the board and the compact view, `1`-`9` switches between
  them from either, the one you are in is marked `▸`, and a board with a live
  agent session is marked `◆`.
- `jaira update` plus a per-command nudge when the binary is newer than the board.

Ownership and notes
- A ticket belongs to its assignee: writing to someone else's is refused with
  exit 3. The human checkpoint lanes are exempt by contract, hand-over is always
  allowed, `--force` overrides and is recorded. **A guard rail, not a lock** —
  git still merges files, so the merge rules are untouched.
- The working lane never named `jaira note` at all. It now says what a note
  contains and when one is due: **write down what the repository does not already
  say**, at every pause, never saved for the end.

## Decisions that must not be silently reversed

1. **`precedence` is a merge rank, not a column position.** `core/merge` uses it
   to decide which lane wins when two clones moved the same ticket, which is why
   `blocked` carries a deliberately low number. Renumbering it to fix display
   order would make a merge pull an in-progress ticket back into Blocked. This
   was nearly shipped as a "fix" and the plan for it was withdrawn.
2. **Lanes are written as files**, not through a CRUD form. The CLI makes them
   legible and actionable; the TUI's `n` writes a skeleton and opens `$EDITOR`.
3. **A progress note is never gated.** A note written to satisfy a gate says
   "work done, all good" and is worthless.
4. **`.jaira/lanes/` is authoritative** when present — the project's whole lane
   list, not a set of overrides. Any change to it must first materialise the full
   working set, or the board drops to one column.
5. **`jaira create` puts a ticket in the backlog even when fully specified.**
   Decided 2026-08-13: capturing and deciding to start are two acts. `--lane todo`
   is the direct route.

## Open, needs the user

- **Do the per-ticket `## Options` (brainstorm, planning) earn their keep** now
  that lanes can be added or left out per project? The user's instinct is that
  removing the lane is simpler; the counter-argument is that an agent driving
  `jaira next` needs this ticket's *route*, and route is per ticket.
- Four interactive checks nobody has run, all needing a real terminal: `S` and its
  two entries, `1`-`9` switching and the `▸`/`◆` markers, the small board's
  `H`/`L`/`x`/`+`, `y` copying over OSC52, and `$EDITOR` opening from `n`.

## Known and deliberately not fixed

- A cosmetic `anchor "" is not installed` warning can appear after materialising a
  project's lanes. Non-corrupting; deferred rather than risk a rushed change to
  `order()`'s bucketing.
- **Renaming text inside a built-in lane makes every existing project copy of it
  start warning**, because the copy now genuinely differs. Happened today when
  "sign-off" became "human review". There is no migration; re-copy or accept it.
- No way to forget a board: `jaira projects` has no removal, and the launcher has
  no key for it. Six stale `/tmp` entries had to be edited out of
  `~/.jaira/projects.json` by hand. `~/.jaira/state/` also grows without bound —
  439 directories, never pruned.
- `▸` (current board) and the header dropping its duplicate name have **no tests**.
  The suite stays green either way.
