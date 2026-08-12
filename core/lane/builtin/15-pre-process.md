---
id: pre-process
name: Pre-process
after: todo
precedence: 25
agentic: true
terminal: false
model-tier: strong
input-requires: [goal, definition-of-done, context]
output-produces: [plan]
requires-option: planning
description: Working out how the change will be made. Only for tickets whose Options tick "planning".
---
# Prompt

Work out how this change should be made. Do not make it yet.

Read the goal, the context and the definition of done. Then decide the method:
what has to be understood first, what has to be designed, what has to be written,
and in what order. Where the shape of the change is not obvious, say what the
options are and which you would take.

Write the result as a `## Plan` checklist in the ticket body — one unticked item
per step, in the order they should happen:

```markdown
## Plan

- [ ] read how the existing exporter streams rows
- [ ] design the chunking boundary
- [ ] failing test for a 40MB export
- [ ] implement
```

Keep the steps at the size of a thing you can finish and mark done. "Implement
the feature" is not a step; it is the whole ticket restated.

If the definition of done cannot be met as written — it is ambiguous, or it asks
for something the codebase makes impossible — say so and move the ticket to the
HITL lane with that question rather than planning around it.
