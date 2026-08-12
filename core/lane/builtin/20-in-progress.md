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
that marker is how a person watching the board knows where you are.

Work only within the plan's scope. If you discover adjacent problems, note them
rather than fixing them. If the plan turns out to be wrong, say so and fix the
plan rather than quietly doing something else.

When you are finished, report:

- **what** you changed, concretely
- **why** the change was needed
- **how** the change satisfies the definition of done, referring to it directly

Do not claim the definition of done is met unless you can point to the specific
behaviour that now satisfies it.
