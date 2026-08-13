# Lane system: catalogue, project lanes, and sharing

**Captured:** 2026-08-13
**Status:** design agreed with the user, not yet planned as a phase

This is the input for a future phase. It records decisions, not implementation.

## Why

Lanes are the pipeline. Today they cannot be inspected, edited, reordered or
shared: `jaira lanes` is read-only and does not even print the prompt, custom
lanes live only in `~/.jaira/lanes`, and a custom lane whose id collides with a
built-in is silently discarded (`core/lane/lane.go:309-316`). A user who wants
"swap the review lane with HITL and add a second review after it" has no way to
do it, and no way to hand the result to a teammate.

## Storage layout

| Location | Holds | Git |
|---|---|---|
| `~/.jaira/lanes/` | the global catalogue — every lane this user built or adopted | outside any repo |
| `.jaira/lanes/` | only the lanes this project actually uses, as full copies exported from the catalogue | **gitignored, including after `jaira share`** |
| `.jaira/shared/<user>/` | lanes deliberately published to teammates | committed |

Decisions behind this:

- The project directory holds copies, not name references. A reference would
  resolve against a catalogue the teammate does not have.
- It holds only the lanes in use, not the whole catalogue. Export from the
  catalogue is a deliberate act, done through the lane settings screen.
- Project lanes stay private by default. Sharing is opt-in through `shared/`,
  not a side effect of `jaira share`.
- Adopting someone's shared lane copies it into the adopter's global catalogue.

## Gitignore consequence

`jaira share` today removes `/.jaira/` from `.gitignore` entirely and tells the
user to `git add .jaira` (`internal/cli/share.go:110`), which would commit the
lanes. Under this design `share` must additionally write a plain `/.jaira/lanes/`
entry. This is an ordinary ignore rule, not a negation inside an ignored tree, so
no nested-negation trickery is needed.

## Built-in lanes become overridable

The current refusal (`core/lane/lane.go:309-316`) is reversed: a custom lane may
shadow a built-in, including its prompt. It must produce a warning, never a
silent swap. The warning channel already exists and is surfaced everywhere
(`jaira lanes`, `--json`, the TUI warnings block, `internal/cli/root.go:240`).

## CLI surface

Lanes become legible, and stay writable as files — the file format is the API.

- `jaira lanes show <id>` — the full lane including its prompt
- the prompt joins `jaira lanes --json` (today it is absent)
- `jaira lanes path` — where lane files belong, so an agent can write one
- a lane template to write against

No `lanes new/edit/move/rm`. Reordering is `after:` plus `precedence:` in the
file, and `order()` already detects cycles and reports them as warnings
(`core/lane/lane.go:396`) — so "edit the file, then run `jaira lanes` to check"
is a complete loop for an agent, without a second write path.

That loop already works today for a custom lane in `~/.jaira/lanes/`: writing a
file there makes the lane appear, `jaira lanes` names its source path, and a
second file reusing an id is reported as a warning and ignored — all verified
directly. What an agent is missing is discoverability (where files belong, what
shape they take) and the same loop for the default board, which has no validator
at all.

The same check turned up an ordering problem: display order follows the `after:`
anchor, not `precedence`, so two lanes anchored to the same lane are ordered by
whichever file was read first. A lane with `precedence: 12` rendering before one
with `precedence: 5` was observed, not theorised. Either the column or the
ordering has to give.

## Per-lane signature

Add a `creator:` frontmatter field so an adopted lane keeps its provenance, and
so a project lane that has drifted from the catalogue copy can be recognised.

## TUI

- a lane settings screen: read a lane, see its prompt, export it to
  `.jaira/shared/<user>/`
- lane selection shows teammates' shared lanes and can adopt one

## Settled 2026-08-13

- `<user>` is the existing `identity()`, slugified (`JAIRA_USER` → git
  `user.name` → `$USER`). `JAIRA_USER` is the escape hatch when two teammates
  share a configured name.
- Drift is checked when the lane settings screen opens, not on every command,
  and shown as a per-lane warning with a refresh action. Editing a lane writes
  through to the catalogue, so it is other projects that drift.
- `.jaira/lanes/` is authoritative — it is the only record of which lanes a
  project uses. The export writes the full working set, built-ins included.
- A `brainstorm` lane ships as a built-in optional step (`requires-option:
  brainstorm`, `output-produces: [goal]`).

## The default board (settled 2026-08-13)

`.jaira/lanes/` being authoritative leaves a hole: something has to fill it. The
default board is that something, and it is a per-user setting, so it lives
globally at `~/.jaira/default-board.md` and is reachable from the home screen —
the only per-user surface there is.

It sets two things:

- **which lanes** a newly initialised board gets
- **which ticket Options are pre-ticked** in new tickets, so "always brainstorm"
  is a setting rather than a habit

A repo where nothing was changed still gets no lane files: **directory absent
means the built-ins**, directory present means it is the project's lane list.
Nobody ends up with ten lane files in `.jaira/` for having changed nothing.

The screen selects and orders lanes and pre-ticks options; it does not contain a
form for a lane's prompt, tier or contract. `e` opens the lane file in `$EDITOR`
instead, which is full editing power without a second way to write the same file
— and it is the pattern the TUI already uses for a ticket body
(`internal/tui/external.go`, `VISUAL` then `EDITOR`).
