# Ticket schema — grammar, lane fields, History, routes-on

2026-08-31 · outcome of the schema brainstorm (state: `.planning/schema-brainstorm.md`,
decisions 1–10 there are binding) · status: **awaiting Berk's review** — nothing
below is implemented.

## Problem

Three failures, all live on this board today:

- A lane file with a typo in `output-produces` is silently a lane with no
  contract — unknown frontmatter keys are ignored without a word
  (`core/lane/lane.go` `parse`, ~224).
- A ticket that loops through critique twice keeps only the last
  `review-summary`; the first round is gone. Nothing records how a ticket got
  to where it stands.
- A ticket sitting in review shows a reviewer only the fields that are filled
  (`internal/tui/view.go` ~848: blocks render only when non-empty), so an
  unworked ticket in review looks like a finished one minus the reading. The
  reviewer cannot see what is owed. (Ticketed: `A1TZ4N`.)

Underneath all three: the ticket file has no declared schema. Fields exist
because code happens to know their names.

## Design

### 1. The file is the truth; the grammar is fixed

The markdown ticket file is the single source of truth; the TUI and CLI render
and edit it, never extend it. Every field a form shows is a field in the file,
and a form writes the same bytes the CLI writes.

The *schema* varies per board. The *grammar* never does:

- frontmatter: `key: value`, label left, content right, key = label
  (kebab-case → words; no separate label field, ever)
- body: `## Heading` + prose or a checklist

### 2. Six shapes, fixed forever

A declaration is a `<…>` value; anything else is a default.

| shape | meaning | filled when |
|---|---|---|
| `<text>` | required free text | non-empty |
| `<text?>` | optional free text | always (may stay empty) |
| `<one-of: a \| b>` | exactly one listed value | value ∈ list |
| `<list>` | list of strings | ≥ 1 item |
| `<date>` | ISO date | parses as a date |
| `<yes/no>` | boolean | `yes` or `no` |

Body sections declared in the template (`## Risks` etc.) are prose or
checklist; they have no shape beyond that. There is no seventh shape and no
extension mechanism — a board that needs more does not get it from jaira.

### 3. Two declaration sources, no third file

- **`.jaira/TEMPLATE.md`** declares the fields the **creator** fills in. This
  file already exists as a mechanism (`internal/cli/tickets.go` `templateBody`,
  ~289: `.jaira/template.md|TEMPLATE.md`, frontmatter copied as defaults).
  Change: a `<…>` value is a declaration — the create path asks for it and the
  create form renders it — not a literal default to copy. Non-`<…>` values keep
  today's behaviour: board defaults, copied verbatim.
- **Each lane file** declares the fields **it** produces. `output-produces`
  keeps its list form (`[review-summary]`, shape defaults to `<text>`) and
  gains a mapping form:

      output-produces:
        review-summary: <text>
        review-verdict: <one-of: accept | reject>

The ticket's field set is the union: base fields + template declarations +
every installed lane's produces. Nothing else is a schema source.

`input-requires` stays a plain name list — it references fields, it does not
declare them. Enforcement is unchanged: bounded input via `show --for-lane`
with `missing` (`internal/cli/flow.go` ~496), and the gate refuses a forward
move while a declared output is empty (`core/gate/gate.go` ~423,
`OutputOwed` ~433). Both gain shape-awareness: "filled" means "matches the
shape" (today's `fieldFilled`, gate.go ~487, already falls through to the raw
doc for unmodelled names; the fallthrough learns shapes).

### 4. Frontmatter holds the latest; `## History` holds every round

- The ticket's frontmatter keeps the **newest** value of each lane field.
  Gates, `--json`, and the merge driver keep reading one flat key, unchanged.
- Every lane pass appends one **immutable** block to a `## History` body
  section on the move out of the lane:

      ### critique · 1 · 2026-08-31T14:02Z · strong
      - review-summary: view.go:309 hardcodes the lane id; use RequiresHumanExit
      → in-progress

  Header: lane · round number (per lane) · timestamp · model tier. Body: that
  lane's produced fields as `- key: value` lines. Footer: the lane it moved to.
- Frontmatter answers "where does it stand?"; History answers "how did it get
  here?". History is append-only; jaira never edits an existing block.
- The human lane's `question`/`answer` exchange becomes a History block like
  any other lane pass; the `question` frontmatter field itself is unchanged.
- Concurrency: two sessions appending to the same section can conflict; the
  merge driver resolves per frontmatter key today, and its treatment of body
  sections must be tested with History before cut 2 ships (test to write, not
  a behaviour to assume).

### 5. `routes-on` makes the loop edge deterministic

New lane key: `routes-on: <field>`. On leaving the lane:

- field value is the sentinel `none` → forward (normal `after:` order)
- anything else → the lane's `rejects-to`

This drives `jaira next`/`next_lane` and the destination `jaira move` picks
when `--to` is not given. **It is routing, not a gate**: an explicit `--to`
still goes where it says, and the backwards-move freedom stays
(invariant 14 — no state machine). Validation at load: `routes-on` names a
field in this lane's produces, and the lane has a `rejects-to`
(`checkContracts`, `core/lane/lane.go` ~871, already warns on bad
`rejects-to`; this extends it).

### 6. Views

- **Base fields, always visible, every view**: id, lane, assignee, creator,
  when (created/updated), goal, context — regardless of lane or what else is
  set.
- **Display style**: concise bullet points that carry the point — not prose
  walls, not bare keywords.
- **The reviewer's view**: a ticket in review/signoff shows, in order:
  problem (goal), what, why, resolves, summary, gaps, check.
- **A declared field is shown even when empty** — the empty row names the lane
  that owes it (`gate.OutputOwed` already computes this). This is `A1TZ4N` and
  lands with or before cut 4.
- First glance shows the base form; lane fields and History are folded,
  expandable (`B4MGTP`). History renders as rounds.

### 7. Validation — where it refuses, where it warns

- **Refuses** (gate, on a forward move): a declared output that is empty or
  does not match its shape. Same exit code and JSON error shape as today's
  `missing_lane_output`.
- **Warns** (`jaira lanes` / load / `jaira validate`): unknown frontmatter
  keys in a lane file (the typo problem — `NFJCTK`), a `routes-on` naming a
  field the lane does not produce, a shape that does not parse, a `<…>`
  declaration jaira does not recognise.
- **`jaira set` keeps validating nothing** — invariant 9 stands. Write
  anything, exit 0; the gate is where a bad value stops mattering. Reopening
  invariant 9 is Berk's call and not part of this spec.

### 8. Migration and compatibility

- History starts at adoption: existing tickets get a `## History` section on
  their first lane pass after cut 2 ships. Nothing is derived retroactively.
- Existing lane files (list-form produces, no shapes) parse exactly as today:
  list form means every field is `<text>`.
- `show --for-lane --json` keeps `produces` as a name list (agents parse it —
  stable schema promise in CLAUDE.md); shapes ride in a new `produces-shapes`
  map. No existing JSON key changes meaning.
- `## Options` and `## Plan` are untouched.
- Unknown ticket frontmatter (`external:` etc.) survives untouched — writes
  stay AST-level single-field patches (`core/ticket/frontmatter.go`
  `SetScalar`).

### 9. Non-goals

- Lane-produced **artifacts** (named files a lane saves for the next) —
  `FCMP17`, deliberately its own cut, not in this spec.
- New shapes, per-board grammar extensions, a label field.
- Enforcing loops with a state machine (invariant 14).
- Per-field validation on `jaira set` (invariant 9).

## Where this lands in the code

| concern | today | change |
|---|---|---|
| lane frontmatter parse | `core/lane/lane.go` `parse` ~224, unknown keys silent | known-key set + warning; mapping-form produces |
| shape type | — | one small type in `core/ticket` (both lane and template readers use it) |
| template read | `internal/cli/tickets.go` `templateBody` ~289 copies scalars as defaults | `<…>` values become declarations, not defaults |
| "is it filled" | `gate.go` `fieldFilled` ~487, doc fallthrough | shape-aware fallthrough |
| bounded input | `flow.go` `showForLane` ~496 / `fieldValue` ~598 | shape-aware; `produces-shapes` in JSON |
| owed fields | `gate.OutputOwed` ~433 | feeds empty-row rendering (A1TZ4N) |
| move write path | **three sites**: `internal/cli/flow.go` ~207, `internal/tui/model.go` `applyMove` ~1291 + forced ~1343 | History append must run at all three — unify them or touch all three, decided in cut 2's plan |
| load-time warnings | `checkContracts` ~871 | + routes-on checks, + unknown-key warnings |
| TUI detail pane | `view.go` ~848 renders non-empty blocks only | base-first + folds + empty declared rows |

## Build order — six cuts, each shippable alone, each through the board's lanes

1. **Grammar**: shape type, shape-aware parse of template + lane produces,
   gate/bounded-input shape awareness, unknown-lane-key warning.
   Done when: a typo'd lane key warns; a `<one-of>` field refuses a stray
   value at the gate; every existing board parses unchanged.
2. **`## History`**: append one block per lane pass at every move write-site;
   merge behaviour for concurrent appends tested.
   Done when: two critique rounds are both readable on the ticket.
3. **`routes-on`**: lane key + validation; drives `next_lane`, `jaira next`,
   and `move` without `--to`.
   Done when: critique with `review-summary: none` forwards, anything else
   routes to `rejects-to`, with no `--to` typed.
4. **TUI**: base-first view, folds, History-as-rounds, reviewer's view order,
   empty declared rows (`A1TZ4N`, `B4MGTP`).
5. **Create form**: generated from TEMPLATE.md declarations; writes the same
   bytes `jaira create` writes.
6. **Docs**: `docs/TICKET.md`, `docs/LANES.md` (`NFJCTK`), `jaira ticket
   template` (prints a starter TEMPLATE.md, sibling of `jaira lanes template`).

Each cut is a ticket worked through this board's own lanes. Cut 1 starts only
after this spec is approved.

## Tickets this spec folds in or touches

- `NFJCTK` — lane schema doc → cut 6 (+ the cut-1 warning)
- `B4MGTP` — foldable history → cut 4
- `A1TZ4N` — reviewer's view / empty declared fields → cut 4 (or earlier)
- `FCMP17` — artifacts → explicitly out (own cut, later)
- `TQXBY5` — multi-target rejects-to: untouched; routes-on composes with it
  later if it lands
