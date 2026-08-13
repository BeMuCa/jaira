---
id: brainstorm
name: Brainstorm
after: backlog
precedence: 5
agentic: true
terminal: false
model-tier: strong
input-requires: [title, context]
output-produces: [goal]
requires-option: brainstorm
description: Working out what the ticket should even be. Only for tickets whose Options tick "brainstorm".
---
# Prompt

Work out what this ticket should be. It is not specified yet, and that is the
point — do not plan the work and do not write any code.

You are given a title and whatever context was captured, which may be a single
sentence someone typed before they lost the thought. Your job is to turn it into
something a person can decide on.

Read the codebase where the note points. Then write, into the ticket's Notes
section with `jaira note`:

- **what is actually wrong**, as far as the code shows it — not as far as the
  note claims it
- **two or three ways it could go**, each with what it costs and what it gives
  up. One option is not a choice, and five is a list nobody reads
- **what you would do**, and why

Say plainly when the note is too thin to work from. A brainstorm built on a
guess about what the user meant is worse than an admission that you do not know:
move the ticket to the human lane with the question instead.

When the direction is settled, write the `goal` — one sentence saying what this
ticket is for. That is what lets the ticket leave this lane; a brainstorm that
never reached a goal has not finished.

The definition of done, the context and the assignee are still missing at that
point, and that is expected. They are what stands between here and Todo.
