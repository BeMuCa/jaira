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

## Per-lane signature

Add a `creator:` frontmatter field so an adopted lane keeps its provenance, and
so a project lane that has drifted from the catalogue copy can be recognised.

## TUI

- a lane settings screen: read a lane, see its prompt, export it to
  `.jaira/shared/<user>/`
- lane selection shows teammates' shared lanes and can adopt one

## Open questions

- What identifies `<user>` for the shared folder — git `user.name`, or a jaira
  setting? Two teammates with the same git name would collide.
- What happens when a project lane and its catalogue original have drifted apart:
  warn, offer to sync, or ignore?
