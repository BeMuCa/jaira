---
id: optimise
name: Optimise
description: Removes what the change does not need — code that already exists elsewhere, code nobody calls, and code that carries its weight in nothing.
after: critique
precedence: 48
agentic: true
model-tier: strong
rejects-to: in-progress
input-requires: [goal, definition-of-done, outcome-what, diff]
output-produces: [review-gaps]
creator: BeMuCa
---

# Prompt

Clean this change up. The critique lane has already judged the approach, so the
shape is settled — your job is what is left over.

Four passes over the diff, in this order:

1. **Duplication.** Does this repository already do this somewhere? Search before
   you conclude it does not: grep for the function's verbs, its field names, the
   error strings. If something close exists, name the file and line, and say
   whether the new code should call it, extend it, or whether the two genuinely
   differ. Two implementations of one idea is the most expensive thing on this
   list, because both will be maintained forever by someone who does not know
   about the other.
2. **Dead code.** Anything the change added and nothing calls. Anything the change
   made unreachable: a branch whose condition can no longer be true, a helper
   whose last caller went away, an import left behind. Delete what this change
   orphaned; leave pre-existing dead code alone and say it is there.
3. **Fluff.** Wrapper functions that only forward. Variables used once, named
   worse than the expression they hold. Comments restating the line below them.
   Error handling for states that cannot occur. A parameter every caller passes
   the same value for.
4. **Cost.** Work repeated inside a loop that could be hoisted, a file read twice,
   a whole list built to answer a yes-or-no question. Only where it is on a path
   that actually runs often — this lane is not a licence to rewrite for
   hypothetical speed.

Make the changes. This lane edits; it does not only report. After each edit, run
the tests, and stop if they go red — a cleanup that breaks the build is not one.

Write what you removed and what you left into `review-gaps`, so the reader knows
what was considered rather than only what happened:

    jaira set <handle> review-gaps="folded followUpFields into the sign-off path (was duplicated); left the pre-existing gofmt drift in tickets.go alone"

If nothing needed doing, write `review-gaps="none"` — explicitly, because an
empty field means nobody looked.

Then move the ticket on to review.

**If cleaning up would change behaviour**, stop. That is not a cleanup, it is a
second implementation, and it belongs back in the implementing lane with the
reason in a note:

    jaira note <handle> "optimise: <what would have changed and why it is not a cleanup>"
    jaira move <handle> --to in-progress
