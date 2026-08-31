# The schema brainstorm — state as of 2026-08-31

Running under the brainstorming skill, architectural path: questions one at a
time → approaches → sectioned design → spec in `docs/superpowers/specs/` →
implementation plan. **Nothing is implemented.** This file is the state between
sessions; the spec supersedes it once written.

## Decided

1. **The markdown file is the truth; the TUI is a convenience** (his "a").
   Every field a form shows is a field in the file; a form writes the same
   bytes the CLI writes. Constraint he added: everything in the file follows a
   **fixed grammar** the TUI can fully parse and render — label left, content
   right — for display and for the create form. The *schema* varies per board;
   the *grammar* never does.
2. **Two declaration sources, no third file.** The board template
   (`.jaira/TEMPLATE.md`) declares the fields the **creator** fills in; each
   lane declares the fields **it** produces. jaira derives the ticket from the
   union. (He asked whether the TUI can build fields dynamically from keywords
   in the md — yes, provided the keyword carries the *shape*.)
3. **Six shapes, fixed forever**: `<text>` `<text?>` `<one-of: a | b>`
   `<list>` `<date>` `<yes/no>` in frontmatter; `## Heading` + text or a
   checklist in the body. A value is a default; a `<…>` is a declaration.
4. **The key is the label**: frontmatter key / body heading, kebab→words. No
   separate label field.
5. `requires`/`produces` keep today's enforcement (bounded input via
   `show --for-lane` + `missing`; the gate refuses leaving with an empty
   declared output) and gain shapes. `routes-on: <field>` makes the loop edge
   deterministic: sentinel value (`none`) → forward, else → `rejects-to`.

## The worked example he has seen

Template with `severity: <one-of: low | medium | high>` and a `## Risks`
section; critique lane with `produces: review-summary: <text>` and
`routes-on`; the resulting ticket; the TUI mock with folds and a `## History`
section holding one appended block per lane pass
(`### critique · 1 · <ts> · <model>` … `→ in-progress`). In the chat transcript
of 2026-08-30 and reproducible from the decisions above.

## Decided (continued)

6. **Q3, decided 2026-08-31: latest + History.** The ticket's frontmatter
   keeps the newest value of every lane field — gates, `--json` and the merge
   driver keep reading one field, unchanged — and every lane pass appends one
   immutable block to `## History`. Frontmatter answers "where does it
   stand?", History answers "how did it get here?".

## Decided (continued)

7. **Q4, decided 2026-08-31: build order A** — grammar-first, six board-sized
   cuts, each shippable alone: (1) shapes + validation + unknown-key warning,
   (2) `## History`, (3) `routes-on`, (4) TUI folds + round view, (5) the
   generated create form, (6) `docs/TICKET.md`, `docs/LANES.md`,
   `jaira ticket template`. **Execution starts after a context clear, on his
   "go" — not before.** The spec is written first
   (`docs/superpowers/specs/…-schema-design.md`), he reviews it, then cut 1.

8. **Base fields, always visible** (his requirement, 2026-08-31): every ticket
   view leads with `id, lane, assignee, creator, when (created/updated), goal,
   context` — regardless of lane, regardless of what else is set.

9. **Display style**: concise bullet points that still carry enough to
   understand the point — not prose walls, not bare keywords.

10. **The reviewer's view**: a ticket in review/signoff shows, in this order:
    problem (goal), what, why, resolves, summary, gaps, check — concise
    bullets. The detail pane already renders Outcome and Review blocks
    (view.go ~848-860) but ONLY when non-empty, so a ticket that reached
    review unworked shows just base + DoD and the reviewer cannot see what is
    owed. Rule for the spec: a field a lane on this board declares is shown
    even when empty (an empty row names the lane that owes it) — ticketed as
    a defect to land with or before cut 4.

## Open

Nothing — the design questions are closed. Next artefact: the spec.

## Parked for the spec

- What existing tickets do on first contact (migration: seed `## History`
  empty? derive nothing retroactively — history starts at adoption).
- `## Options`/`## Plan` unchanged; `question`/`answer` become a History
  block from the human lane.
- Docs owed: `docs/TICKET.md`, `docs/LANES.md`, `jaira ticket template`,
  a `validate` warning for unknown lane keys (ties `NFJCTK`), and the
  marketplace README already points at `lanes template`.
- Related tickets folded in: `NFJCTK` (lane schema doc), `B4MGTP` (foldable
  history), `FCMP17` (artifacts — deliberately NOT in this spec; own cut).
