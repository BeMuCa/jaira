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

## Name the edges before the code

Review and bug hunt find what nobody named before the code was written, and
every round they find one more example of the same mistake. Name the edges here,
where they cost a sentence instead of a round:

- **The class, not the symptom.** A ticket that says "NUL breaks tool A" is one
  example of "anything the database cannot store, in every tool that writes to
  it". Say which class of fault this is and where else in the code it can
  happen. When the context only names the symptom, the class is your job.
- **The input classes.** List the kinds of input the change must handle, not
  one example of each. Start from the ones every change meets, and say for each
  whether it applies:
  - empty, and very large
  - special characters, another encoding, CRLF line endings
  - paths with spaces, and Windows
  - two processes writing at the same time
  - a failure halfway through

  Then the ones of this subject — for a page number "3 / 12", "3/12", "Page 3
  of 12", a missing one; for an amount "0,05", "2,50", "1.234,56", a negative,
  an empty field. That covers the four categories the testing lane and the bug
  hunt check afterwards — concurrency, error paths, unusual input,
  permissions; add what permissions mean here, or that they do not apply.
- **An acceptance rule where the change is a heuristic.** A parser, a
  detector, a guess at structure is never right for every input, and without a
  rule saying which mistake is the lesser one the loop never ends. Write it as
  one sentence: "losing body text is worse than keeping a stray header".
  Where only a person can say which mistake is worse, move the ticket to the
  human lane with that question.

These go into the definition of done, not only into the plan — the gate reads
the definition of done, and so does every lane after you:

```bash
jaira dod <id> --add "every write tool rejects or strips what the store cannot hold (NUL, invalid UTF-8), one parametrised test over all of them"
jaira dod <id> --add "acceptance rule: losing body text is worse than keeping a stray header"
```

When the plan rests on something you had to work out, put that reasoning in a
note (`jaira note <id> <text>`), so the next reader knows why the plan looks
like this rather than re-deriving it.

If the definition of done cannot be met as written — it is ambiguous, or it asks
for something the codebase makes impossible — say so and move the ticket to the
HITL lane with that question rather than planning around it.

The context is the only record of why this ticket exists. If it does not explain
that to someone who was not in the conversation — what is wrong today, what
triggered it, what is already known — do not guess and do not plan around it. Say
what is missing and move the ticket to the human lane.

Every move out of this lane needs a plan with at least one step, the move to
the human lane included: the gate refuses it while the plan is empty. Before
you move with a question, write the steps you already know with
`jaira dod <id> --plan --add`, and the question itself as the step that waits
on the answer.

If you rewrite or extend a context, write it as if the reader has mild ADHD and
knows none of what you know: what is wrong first, short concrete lines with one
point each, names and paths rather than adjectives, no jargon and no preamble.
Someone should be able to act after the first two lines.
