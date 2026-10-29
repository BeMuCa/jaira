---
id: 01M5GF4QA5BPN7PPT2DWW0JX10
title: Eine Lane Feature Ideas sammelt unreife Ideen vor dem Backlog; Brainstorm macht daraus verlinkte Backlog-Tickets
status: ideas
ready: true
creator: berkc
assignee: berkc
goal: "Jede Idee hat einen Platz auf dem Board, bevor sie ein Ticket ist; das Backlog enthaelt nur durchdachte Tickets, und jedes davon zeigt, aus welcher Idee es kam."
context: |-
  Wunsch von berkc am 22.10.2026, woertlich: 'Lane idea: Feature Ideas; Everything that pops into your mind, but is not ripe enough to land in the backlog. As the backlog contains thought out tickets that can be worked on. Then we could also put the brainstorming lane after the feature idea lane. the feature idea lane contains thoughts and ideas, if you want to turn them into tickets, you can brainstorm them (in the brainstorming lane) and afterwards turn them into the finalized backlog tickets (we might even link the finalized backlog tickets to the feature idea ticket, so that we can keep in mind how the ticket came into being). the tickets in the feature ideas lane should contain a collection of things that the' - der Satz bricht dort ab; berkc fragen, wie er weitergeht.
  Ausserdem: alles, was berkc zu einer Idee sagt, wird auf dem Ideen-Ticket notiert (jaira note).
  Schon erledigt (Commit mit dieser ID auf feat/T8A8KM-update-flow): .jaira/lanes/ideas.md, erste Zeile in .jaira/lanes/order, nicht agentisch, ohne requires-specified - jaira create --lane ideas legt Tickets ohne Ziel/DoD dort ab (der Spezifiziert-Gate greift erst ab todo, core/gate/gate.go:236).
  Offen:
  - Reihenfolge ideas -> brainstorm -> backlog -> todo statt heute backlog -> brainstorm -> todo. Heute gilt brainstorm nur mit Option 'brainstorm' (requires-option) und erzeugt das goal DESSELBEN Tickets; gewuenscht ist: aus einer Idee entstehen NEUE Backlog-Tickets.
  - Es gibt schon den Skill /jaira-role-brainstorm (core/role/builtin/jaira-role-brainstorm/SKILL.md): 'Turn a raw idea into tickets ... or say it is not a ticket'. Er passt genau auf den Schritt Idee -> Tickets.
  - Verlinkung: jaira create --follows <idee> existiert; zu klaeren, ob das Ideen-Ticket danach archiviert wird oder in einer Lane 'shaped' bleibt.
  - Die Lane nur auf diesem Board oder als ausgelieferte Lane fuer jedes neue Board?
definition-of-done: "Lane ideas steht vor backlog, nimmt Tickets ohne goal/DoD auf; jaira create --lane ideas funktioniert; Test"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T00:10:58Z
updated-at: 2026-10-22T00:11:21Z
updated-by: berkc
---

# Eine Lane Feature Ideas sammelt unreife Ideen vor dem Backlog; Brainstorm macht daraus verlinkte Backlog-Tickets

## Definition of Done

- [ ] Lane ideas steht vor backlog, nimmt Tickets ohne goal/DoD auf; jaira create --lane ideas funktioniert; Test
- [ ] Der Weg Idee -> Backlog-Tickets ist entschieden (Lane-Reihenfolge, Rolle des brainstorm-Lanes, was mit dem Ideen-Ticket danach passiert) und umgesetzt
- [ ] Ein aus einer Idee entstandenes Ticket zeigt die Idee (follows) und die Idee zeigt ihre Tickets (jaira links); Test
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md, falls die Lane ausgeliefert wird

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

