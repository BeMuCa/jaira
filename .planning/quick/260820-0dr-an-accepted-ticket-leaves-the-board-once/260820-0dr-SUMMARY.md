---
quick_id: 260820-0dr
status: complete
---

# Summary: an accepted ticket leaves the board once its work is pushed

## The gap that was found

Nothing said when a ticket should be archived. `archive` was documented as a
command (README, COMMANDS.md, SKILL.md) and as a key (`x`), and `push` appeared in
no lane file and nowhere in `core/lane` or `core/gate` — only twice in README as
team-flow prose. So a done ticket sat on the board indefinitely unless someone
decided on their own that it should not.

## What was written down

The rule, and the reason for its order: an accepted ticket leaves the board once
its work is pushed. Archiving ahead of the push hands a teammate a board that has
forgotten a ticket while the code it describes has not arrived, which is the exact
state this board exists to prevent. A follow-up keeps its `follows` link to an
archived predecessor, so clearing finished tickets loses nothing.

Four places, chosen because they are where the rule gets read:

- `core/lane/builtin/50-done.md` — one clause on the `description`. Verified it
  actually surfaces: `jaira lanes show done` prints it (`internal/cli/tickets.go:824`).
- `.claude/skills/jaira/SKILL.md` — "Taking a ticket off the board" now says when,
  why the order matters, and that an agent does not push on its own initiative, so
  the archive waits on the user.
- `docs/AGENTS.md` — the loop no longer ends at the hand-off; the prose after it
  says the ticket comes off the board after the push, and not on the agent's
  initiative. Deliberately not in the bash block: that block is what an agent runs
  on one ticket in one session, and pushing is not part of it.
- `README.md` — one paragraph closing "Reviewing finished work" for a person.

## What was deliberately not done

No enforcement. There is no push check anywhere in the binary and this task did
not add one: archiving is the step after the last lane, not a lane transition, so
a gate has nothing to fire on. `core/gate` untouched.

The done lane keeps its shape — frontmatter only, no body. A body there would be
dead weight: `core/lane/lane.go:244` turns a lane body into its `Prompt`, and the
lane is `agentic: false`, so no prompt is ever dispatched.

## Worth knowing

The edit changes the **embedded** built-in only. A project that has its own copy of
the `done` lane keeps the old description until `jaira lanes use done --force`,
which is the only refresh path (the TUI's `R` compares against the user catalogue,
never the embedded built-ins).

## Verification

`go build ./...`, `go test ./... -race` all green. `jaira lanes show done` prints
the new sentence. README still has 0 em dashes and 0 Unicode ellipses.
