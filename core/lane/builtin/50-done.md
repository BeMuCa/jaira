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
description: Accepted. Every definition-of-done item must be marked done, the plan finished if there is one, and the commits that carry the change recorded. Finished tickets stay here until somebody files them — 'jaira logbook --all' puts the lot into today's folder, 'jaira logbook <id>' files one.
---
