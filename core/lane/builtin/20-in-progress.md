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

Write a note (`jaira note <id> <text>`) when you finish a plan step. That is the
one trigger you can actually rely on: a session does not get a turn to write
something down when a usage limit or a crash ends it, so "write one before you
stop" fires exactly never in the case that matters most. Finishing a step
happens on its own schedule and leaves the trail at a point someone else could
pick up from.

The note says what the checklist marker cannot: what you learned doing the step,
and what the next one actually needs. If the plan alone is enough for the next
reader, one line is enough — do not pad it. The value is a readable trail, not a
transcript.

Write one outside that rhythm too when: something you tried did not work, which
is the single most valuable line a later session can read and the one most often
skipped; you learn something that contradicts the ticket's context or its plan —
fix the plan too, but say what changed and why; you are about to start something
long, so an interruption in the middle is recoverable; or you are handing the
ticket over or stopping deliberately.

When you are finished, report:

- **what** you changed, concretely
- **why** the change was needed
- **how** the change satisfies the definition of done, referring to it directly

Do not claim the definition of done is met unless you can point to the specific
behaviour that now satisfies it — the `--proof` on each item is where that
pointer goes.
