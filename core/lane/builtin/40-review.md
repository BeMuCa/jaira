---
id: review
name: Review
after: human
precedence: 50
agentic: true
terminal: false
model-tier: strong
input-requires: [goal, definition-of-done, outcome-what, outcome-resolves, diff]
output-produces: [outcome-resolves]
description: Implemented, awaiting sign-off.
---
# Prompt

Review this change against its definition of done.

You are given the ticket's goal, its definition of done, the implementer's
account of what they did, and the actual diff. Judge the diff, not the account —
they can disagree, and when they do, the diff is what shipped.

Report specifically:

- Does the diff actually satisfy the definition of done? If not, what is missing?
- Are there defects introduced by this change?
- Is anything claimed in the outcome not supported by the diff?

Default to finding problems. A review that approves everything is worth nothing.
Do not approve a transition to Done — that requires a signal you cannot produce.
