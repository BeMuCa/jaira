---
id: in-progress
name: In Progress
after: todo
precedence: 30
agentic: true
terminal: false
model-tier: cheap
input-requires: [goal, definition-of-done, context]
output-produces: [outcome-what, outcome-why, outcome-resolves]
description: Being implemented.
---
# Prompt

Implement the change this ticket asks for.

You are given the ticket's goal, its context, and its definition of done. Work
only within that scope — if you discover adjacent problems, note them rather
than fixing them.

When you are finished, report:

- **what** you changed, concretely
- **why** the change was needed
- **how** the change satisfies the definition of done, referring to it directly

Do not claim the definition of done is met unless you can point to the specific
behaviour that now satisfies it.
