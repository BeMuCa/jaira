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

## Open — the next question to him

**Q4: build order** (approaches presented, his pick pending) — see the chat of
2026-08-31.

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
