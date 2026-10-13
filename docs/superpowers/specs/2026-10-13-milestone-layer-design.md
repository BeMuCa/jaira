# Milestone layer — bubbles, swimlanes, arrows

2026-10-13 · ticket D0SAHM · outcome of the brainstorm of 2026-10-12/13 (decisions
are in the ticket's notes) · status: **approved by Berk to build** ("Ich vertrau
all deinen Entscheidungen").

## Problem

The board answers "where does each ticket stand". Nothing answers "what comes
after this round of work, and how far along is each round". Milestones already
group tickets (`core/milestone`), but they are a flat list: no goal text, no
order, no picture.

## What the layer is

A second view of the TUI, opened with `L` from the board. Every *placed*
milestone is a bubble; arrows mean "first A, then B". The arrows are a picture
only: nothing is blocked, nothing is hinted.

```
  ╭──────────────╮                  ┊  ╭──────────────╮  [+]
  │ api-umbau    │                  ┊  │ ui-runde     │
  │          6/6 │                  ┊  │          2/5 │
  ╰──────┬───────╯                  ┊  ╰──────┬───────╯
         ▼                          ┊         ▼
  ┏━━━━━━━━━━━━━━┓                  ┊  ╭──────────────╮
  ┃ release-0.4  ┃  [+]             ┊  │ hilfe-texte  │
  ┃          5/7 ┃                  ┊  │          3/3 │
  ┗━━━━━━┯━━━━━━━┛                  ┊  ╰──────────────╯
       ┌─┴───────────────┐          ┊
       ▼                 ▼          ┊
╭──────────────╮  ╭──────────────╮  ┊
│ doku         │  │ tests        │  ┊
│          0/2 │  │          1/4 │  ┊
╰──────────────╯  ╰──────────────╯  ┊
```

## Data — two new frontmatter keys on the milestone file

```
---
color: 209
created-at: 2026-10-12T…
goal: Version 0.4 is out          ← new, free text, one line
after: [api-umbau, ui-runde]      ← new, names of predecessor milestones
---
- 01M…
```

- No `after:` line: the milestone is not in the layer (not placed).
- `after: []`: top of a swimlane of its own.
- `after: [b]`: placed under `b`, in `b`'s swimlane.
- The **first** predecessor decides the swimlane and the parent; every further
  name is an extra arrow.
- An arrow edit touches only the successor's file, so two people adding arrows
  to different milestones never write the same file.
- Older binaries keep both lines verbatim (the parser keeps lines it does not
  know; checked with the 2026-10-05 binary).

Names in `after:` are normalized like milestone names. Unknown names and arrows
that would close a cycle are ignored when drawing (never an error, never a
crash); the layer itself refuses to create them.

## Layout — a pure function

`core/milestone` gets a pure layout function: placed milestones in, a grid out.

1. **Parent** of a placed milestone = its first predecessor, if that name is a
   placed milestone and does not close a cycle; otherwise it has no parent.
2. Milestones without a parent are **swimlane tops**. Each top with all its
   descendants (through the parent relation) is one swimlane.
3. Swimlanes are ordered left to right by the `created-at` of their top
   (ties: name).
4. Inside a swimlane, a tidy tree: row = depth; children ordered by
   `created-at` (ties: name); a subtree is as wide as the sum of its children
   (at least 1); the parent sits in the column of its first child. Swimlane
   width = width of its tree.
5. **Extra arrows** = every further valid predecessor. Same-swimlane ones are
   joins (c→h, g→h); others are cross-lane arrows. All are drawn dim.
6. Output: bubbles with (lane, column, row), parent arrows, extra arrows,
   swimlane column spans. Deterministic for equal input.

## Look

- Bubble: 16 columns × 4 rows; rounded border in the milestone colour; name
  (truncated with `…` past 12 characters) and `done/total` right-aligned.
- Fill (L9, variant A): the two inner rows fill left to right in proportion
  done/total over 14 cells; 0 < done < total never shows empty or full. Fill
  colour = milestone colour mixed 55 % toward white; on truecolor terminals
  as RGB, otherwise the nearest of the 256 standard colours. Text on the fill
  is black when the fill's relative luminance is above 0.35, else white;
  outside the fill the terminal default.
- Selected bubble: thick border (`┏━┓┃┗┛`), bold.
- Grid cell: 21 columns wide (bubble 16 + 5 gap), 6 rows high (bubble 4 + 2
  connector rows). Swimlanes are separated by a dim `┊` in the gap.
- Parent arrows: `┬` in the parent's bottom border, `│`, splits as
  `┌─┴──┐`, `▼` above the child. Extra arrows dim, routed: down one row into
  the gap below the predecessor, horizontally to the gap column left of the
  target, vertically to the gap row above the target, horizontally to its
  centre, `▼`.
- Done = what core/link's relation index answers for the gate and the link
  window (`Satisfied`): the ticket is in a terminal lane, filed in the logbook
  on any day (not only the board's logbook window), or archived from a
  terminal lane. Unknown ids count for nothing. (Critique, 2026-10-13: the
  layer first kept its own answer beside the index; one answer is the point.)

### The `[+]` targets (L8)

Drawn only around the selected bubble, in reserved space, so the layout never
moves while the cursor moves:

- **right** of the selected bubble (in its 5-column gap): a new sibling — same
  first predecessor. Not shown on a swimlane top (a sibling of a top would be a
  new swimlane, which only the far-right `[+]` makes).
- **below** the selected bubble (in its connector rows): a new successor. Shown
  only when the bubble has no child; further successors are siblings of its
  children.
- **far right, top row**: always present; a new swimlane.
- Empty layer: only that `[+]`, with `enter: place the first milestone`.

## Interaction

| key | where | does |
|---|---|---|
| `L` | board | open the layer |
| `L` `esc` `q` | layer | back to the board |
| arrows | layer | move to the nearest bubble or `[+]` in that direction |
| `enter` | `[+]` | choose: `n` new milestone · `c` connect an existing one |
| `enter` | bubble | milestone detail |
| `e` | bubble / detail | edit goal and predecessors |
| `c` | bubble | next palette colour |
| `x` | bubble | take it out of the layer (y/n) |
| `↑↓` `enter` | detail | pick a ticket, open it |
| `b` | detail | board filtered to this milestone (`milestone:<name>`) |
| `esc` | detail | back to the layer |

- **New milestone**: form with name and goal; colour assigned like
  `jaira milestone create` (random free palette colour). Placed with
  `after: [<pred>]` (right/below `[+]`) or `after: []` (far-right `[+]`).
- **Connect existing**: a picker of milestones, without the anchor itself and
  without any that would close a cycle. The far-right `[+]` lists only unplaced
  ones. An unplaced pick is placed where the `[+]` was; a placed pick only gains
  the extra arrow (it does not move).
- **Edit** (`e`): goal (text) and predecessors (comma-separated names; empty =
  swimlane top). Unknown names and cycles are refused with a message. The name
  itself is not editable: jaira has no rename anywhere today.
- **Remove** (`x`): drops the `after:` line, so the milestone leaves the layer
  and keeps its tickets. Its successors are re-attached to its predecessors
  (the chain stays connected). Deleting milestones is N71NVG's.
- **Detail**: goal, `after: …`, then one row per member ticket: id, title,
  lane (or `logbook`). `enter` opens the ticket detail; `esc` there returns to
  the milestone detail.

### Keys that move (Berk, 2026-10-12/13)

- Board: `l` opens the links window (was `L`); `h` and `l` no longer change
  lanes — the arrow keys do. `L` opens the layer.
- Ticket detail: `l` opens the links window (was `L`).
- Other views (move dialog, lane editor, browse) keep their `h`/`l`.
- PAP369's link lines are already on master; their help line and tests follow
  the move from `L` to `l`.

## CLI (every agent must be able to do what the TUI does)

- `jaira milestone create <name> --goal <text> --after <a,b>` (`--after ""`
  = top of a new swimlane).
- `jaira milestone set <name> [--goal <text>] [--after <a,b>] [--no-layer]`.
- `jaira milestone ls --json` carries `goal` and `after` (`after` absent =
  not placed).

## Writes

Every write from the TUI goes through the same path as the CLI (milestone lock,
atomic save, ref record), extracted into one shared function so the two cannot
drift.

## Testing

- `core/milestone`: parse/write round trip of `goal`/`after`, unknown lines
  kept; layout: single lane, split, join, cross-lane, lane order, cycle,
  unknown name, empty.
- `internal/cli`: create/set flags, `ls --json`.
- `internal/tui`: `L` opens, empty layer `[+]`, navigation, new via form,
  connect, edit, colour, remove with re-attach, detail and `b`, fill cells and
  colour codes (truecolor and 256), key moves (`l` links, `h` inert).
- The built binary in tmux on a copy of this repo's board.

## Not in this change

Renaming milestones, deleting or archiving them (N71NVG), persons or assignees
on swimlanes, any rule that blocks work because of an arrow.
