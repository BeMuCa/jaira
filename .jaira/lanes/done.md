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
description: Accepted. Every definition-of-done item must be marked done, the plan finished if there is one, and the commits that carry the change recorded. The move that lands here stamps the commits. Filing is a separate act — 'jaira logbook' takes finished tickets off the board when somebody says so, and 'jaira restore' brings one back.
---
