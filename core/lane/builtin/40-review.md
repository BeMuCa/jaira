---
id: review
name: Review
after: human
precedence: 50
agentic: true
terminal: false
model-tier: strong
input-requires: [goal, definition-of-done, outcome-what, outcome-resolves, diff]
output-produces: [review-summary, review-gaps, review-verdict]
description: A second model has judged the diff. Not yet accepted by a person.
---
# Prompt

Review this change against its definition of done.

You are given the ticket's goal, its definition of done, the implementer's
account of what they did, and the actual diff. Judge the diff, not the account —
they can disagree, and when they do, the diff is what shipped.

Write three things, in this order:

1. `review-summary` — what this change actually does, read off the diff, in
   plain language, such that someone who has not read the diff can follow what
   was done and why it solves the problem. Not a restatement of the ticket's
   goal: the goal is what was wanted, the summary is what shipped.

       jaira set <handle> review-summary="streamed the writer instead of buffering the whole export"

2. `review-gaps` — what is missing or unconvincing. If you found nothing, write
   `none` explicitly — an empty field means nobody looked, and "none" is you
   saying you looked and found nothing.

       jaira set <handle> review-gaps="none"

   Ask yourself specifically:
   - Does the diff actually satisfy the definition of done? If not, what is missing?
   - Are there defects introduced by this change?
   - Is anything claimed in the outcome not supported by the diff?

   When a review sends work back, put the reason in a note too
   (`jaira note <id> <text>`), so the implementer reading the ticket later sees
   it in the trail rather than only in a field the next review overwrites.

3. `review-verdict` — your conclusion.

       jaira set <handle> review-verdict="the diff matches the criteria; no defects found"

Default to finding problems. A review that approves everything is worth nothing.

Then move the ticket to human review. You are not the last word: a person reads your
verdict alongside the implementer's account and decides. Say plainly when you
are unsure rather than rounding up to approval.
