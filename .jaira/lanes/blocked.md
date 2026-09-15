---
id: blocked
name: Blocked
after: done
precedence: 10
agentic: false
terminal: false
requires-blocked-reason: true
description: Waiting on an external dependency. Lower precedence than active work, so a merge never reverts progress into Blocked.
---
