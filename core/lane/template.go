package lane

// Template is the skeleton a fresh lane file starts from: every field parse()
// reads, commented, plus a prompt body — a file this parses without error, so
// the template cannot rot away from the parser it feeds.
//
// Shared by 'jaira lanes template' and the lane settings screen's 'n' key, so
// there is one skeleton to keep in sync with parse(), not two.
const Template = `---
id: my-lane                      # required. lowercase letters, digits and dashes.
name: My Lane                    # optional; defaults to id.
description: One line saying what this lane is for.
after: human                     # anchors ordering: this lane sits after the named lane.
precedence: 42                   # optional. Merge order only — never display order (see 'after').
agentic: true                    # true if a subagent works this lane; false for a human-only step.
model-tier: strong               # a local alias (e.g. cheap, strong) — NOT a model name. This
                                  # indirection is what lets a shared lane file survive a model rename.
terminal: false                  # true marks the lane where signed-off work lands.
holds: 0                         # optional, terminal lane only: a move landing here files the oldest beyond the newest N into the logbook. 0 = unlimited.
logbook-on-entry: false          # optional, terminal lane only: a move landing here files the ticket straight into the logbook, commits stamped.
requires-question: false         # true means a ticket needs an open question before entering.
requires-specified: false        # true marks the first lane a ticket may not skip its promotion fields at.
requires-outcome: false          # defaults to the value of terminal if this key is absent.
requires-nonmodel-signal: false  # defaults to the value of terminal if this key is absent.
requires-human-exit: false       # true means no agent may move a ticket out of this lane.
requires-option: ""              # names a ticket Options entry; the lane only applies if it is ticked.
rejects-to: ""                   # the lane work goes back to when this one finds it wanting: the loop's
                                  # back edge. Must name an installed lane, and not this one. A list —
                                  # [in-progress, human] — declares two back edges: a flaw goes back to
                                  # be fixed, a decision goes to a person.
input-requires: [goal]           # ticket fields assembled into a subagent's bounded input.
output-produces: [outcome]       # ticket fields the subagent must return before it can move on.
creator: you                     # optional provenance; left empty if you'd rather not say.
---

# Prompt

Write the instruction given to the subagent here.
`
