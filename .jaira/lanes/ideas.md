---
id: ideas
name: Feature Ideas
precedence: -5
agentic: false
terminal: false
description: Ideas that are not ripe for the backlog yet. Collect everything said about an idea on its ticket; shape it into backlog tickets with /jaira-role-brainstorm when it is ready.
---

# Feature Ideas

A thought, a wish, a "what if" — not yet a task anyone could pick up. Nothing
here needs a goal, a definition of done or an assignee.

Everything the person says about an idea goes onto its ticket, word for word
where it matters: `jaira note <id> <text>`. The idea ticket is the record of
how the work came about.

When an idea is ripe, shape it with `/jaira-role-brainstorm`: the tickets that
come out of it land in the backlog and point back here with
`jaira create ... --follows <idea-id>`.
