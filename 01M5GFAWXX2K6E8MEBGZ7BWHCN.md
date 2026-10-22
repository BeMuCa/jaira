---
id: 01M5GFAWXX2K6E8MEBGZ7BWHCN
title: "jaira lernt aus Sitzungen und Commits, wie jemand arbeitet, und schlaegt ein Board mit passenden Lanes vor"
status: ideas
ready: true
creator: berkc
assignee: berkc
goal: "Ein Skill erkennt wiederkehrende Arbeitsmuster (z.B. implementieren, dann Kritik-, Optimierungs- und Test-Agent) und schlaegt daraus Lanes und ein Board vor."
context: |-
  Idee von berkc am 22.10.2026: 'Would it make sense to add a feature which makes jaira learn how the user was programming? so Checking if in the session there are reoccurring schemas. for me i regularly implemented then called a critique agent, optimize and test agent - so I created the 3 lanes we currently have. The skill when called should check for these schemas and create a board with lanes for that. and not sure if the commit history might have some things that could help this aswell'
  Moegliche Quellen:
  - Claude-Code-Sitzungsprotokolle (~/.claude/projects/<repo>/*.jsonl): welche Subagenten/Skills in welcher Reihenfolge aufgerufen wurden.
  - Commit-Historie: Muster wie feat -> fix(review) -> test, Commit-Praefixe, Tickets pro Lane.
  - jaira selbst: ~/.jaira/metrics.csv und die Lane-Wechsel der Tickets.
  Datenschutz: Sitzungsprotokolle enthalten alles - nur lokal lesen, nichts senden, Vorschlag vor dem Schreiben zeigen.
  Verwandt: PP24H1 (Interview baut Board), Y2C9QN (Board-Typen), WFZ4EY (Lane-Vorlagen).
definition-of-done: Ein Skill liest lokal Sitzungsprotokolle und Commit-Historie und nennt gefundene Muster mit Haeufigkeit
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T00:14:21Z
updated-at: 2026-10-22T00:14:21Z
---

# jaira lernt aus Sitzungen und Commits, wie jemand arbeitet, und schlaegt ein Board mit passenden Lanes vor

## Definition of Done

- [ ] Ein Skill liest lokal Sitzungsprotokolle und Commit-Historie und nennt gefundene Muster mit Haeufigkeit
- [ ] Er schlaegt Lanes in Reihenfolge vor und legt sie erst nach Zustimmung als Board an
- [ ] Ausprobiert an diesem Repo: findet critique/optimize/testing nach in-progress

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

