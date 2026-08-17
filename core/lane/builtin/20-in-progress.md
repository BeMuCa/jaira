---
id: in-progress
name: Implementing
after: pre-process
precedence: 30
agentic: true
terminal: false
model-tier: cheap
input-requires: [goal, definition-of-done, context, plan]
output-produces: [outcome-what, outcome-why, outcome-resolves]
description: Carrying out the plan.
---
# Prompt

Carry out the plan on this ticket.

You are given the goal, the context, the definition of done, and the plan worked
out in the previous step. Work through the plan in order. Mark the step you are
on with `jaira dod <id> <n> --doing --plan`, and mark it done when it is done —
that marker is how a person watching the board knows where you are. When you
mark a definition-of-done item done, record the evidence in the same call:
`jaira dod <id> <n> --done --proof "<file:line or test name>"`, naming the
specific code or test that makes it true.

Work only within the plan's scope. If you discover adjacent problems, note them
rather than fixing them. If the plan turns out to be wrong, say so and fix the
plan rather than quietly doing something else.

## Progress notes

Write a note with `jaira note <id> <text>`. You are writing to yourself, in a
session that has lost this conversation and can read only the repository and this
ticket.

**The rule: write down what the repository does not already say.**

That test decides both what to write and what to leave out. Applied:

- **a dead end** — "tried X, it does not work because Y". Nothing in the code
  records a thing that was not done. Without it the next session spends its first
  hour rediscovering it. This is the single most valuable line you can leave and
  the one most often skipped.
- **why this and not that** — the code shows the approach you took, never the one
  you rejected or the reason. Without it someone later "fixes" it back.
- **what you had to find out** — "the exporter buffers everything in writeAll()".
  It cost you time to learn; it costs the next reader the same time again.
- **work that is half done** — "callers in Y updated, Z still on the old shape".
- **what you would do next**, when it is not simply the next line of the plan.

And what to skip, because the board already carries it: which step you are on
(the `--doing` marker says so), that a step is finished (the checklist says so),
that a commit exists (git says so), that the definition of done is met (the
`--proof` on the item says so). A note repeating those is a transcript, and a
transcript is not read.

**When to ask yourself the question:** at every natural pause — a plan step
finished, a piece implemented, a test going green, before starting something
long, on handover. Not on every turn.

Do not save it for when you are about to stop. A session killed by a usage limit
or a crash never gets that turn, which is exactly the case the note exists for.

When you are finished, report:

- **what** you changed, concretely
- **why** the change was needed
- **how** the change satisfies the definition of done, referring to it directly
- **which commits** carry it: `--commits "$(git rev-parse HEAD)"` on the move,
  naming every commit you made. Record them now, while you still know which they
  are — at review and at sign-off the diff shown is the diff of exactly these
  commits, and the done lane refuses a ticket that names none, because a change
  with no commits is a change nobody can check, then or months later.

Do not claim the definition of done is met unless you can point to the specific
behaviour that now satisfies it — the `--proof` on each item is where that
pointer goes.
