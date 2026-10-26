---
id: review
name: Review
after: human
precedence: 50
agentic: true
terminal: false
model-tier: strong
input-requires: [goal, definition-of-done, outcome-what, outcome-resolves, diff]
output-produces: [review-summary, review-gaps, review-verdict, review-check]
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

2. `review-gaps` — what is missing or unconvincing, one finding per line, and
   every line starts with what kind of finding it is:

   - `blocker:` — the diff breaks a definition-of-done item, or introduces a
     defect. Nothing else is a blocker.
   - `note:` — everything else: a simpler form, a missing test of something
     that works, a neighbouring fault, a "would be nice".

   If you found nothing, write `none` explicitly — an empty field means nobody
   looked, and "none" is you saying you looked and found nothing.

       jaira set <handle> review-gaps="none"
       jaira set <handle> review-gaps="note: the export test covers 1MB, not 40MB"

   Ask yourself specifically:
   - Does the diff actually satisfy the definition of done? If not, what is missing?
   - Are there defects introduced by this change?
   - Is anything claimed in the outcome not supported by the diff?

   A note is one line in this field and nothing else. Never create a ticket
   from it, and never ask for one to be created: whether a note becomes a
   ticket is the person's decision when they accept the work, not yours and
   not the next agent's.

3. `review-verdict` — your conclusion.

       jaira set <handle> review-verdict="the diff matches the criteria; no defects found"

4. `review-check` — how a person checks this themselves, as a flow: do this,
   then that, then look at this. Every field above is an account of what
   happened; this is the only one the reader can act on, and without it signing
   off means trusting the account or rebuilding the steps from the diff.

       jaira set <handle> review-check="1. go run ./cmd/app  2. open /export and pick a 40MB range  3. the download starts within a second instead of hanging"

   Write it as if the reader has mild ADHD and knows none of what you know —
   numbered steps, one action each, exact commands and exact paths rather than
   "run the tests", no jargon and no preamble. Name what they should see, not
   only what to press: a step whose outcome is unstated cannot be failed.

   If it cannot be checked by hand at all, say that and say why — a test name
   is a legitimate answer ("go test ./auth -run TestSameSite covers it; there is
   no UI path"). "Review the diff" is not.

Judge against the definition of done. `none` is a normal, expected answer: a
change that meets its criteria and breaks nothing has no blocker, however much
you could add to it.

Then move the ticket on. A blocker sends it back to the lane the work was done
in (`in-progress` on the default board), with the reason in a note
(`jaira note <id> <text>`), so the implementer reading the ticket later sees it
in the trail rather than only in a field the next review overwrites. Without a
blocker it goes on to the person, notes and all. You are not the last word: a
person reads your verdict alongside the implementer's account and decides. Say
plainly when you are unsure rather than rounding up to approval.
