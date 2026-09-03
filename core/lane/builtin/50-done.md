---
id: done
name: Done
after: signoff
precedence: 60
agentic: false
terminal: true
requires-outcome: true
requires-nonmodel-signal: true
requires-commits: true
holds: 10
description: Accepted. Every definition-of-done item must be marked done, the plan finished if there is one, and the commits that carry the change recorded. The ticket leaves the board once that work is pushed, with 'jaira archive'.
---
