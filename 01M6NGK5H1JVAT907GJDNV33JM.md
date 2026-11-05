---
id: 01M6NGK5H1JVAT907GJDNV33JM
title: "Nur noch jaira-worker: der Dispatcher faellt weg, seine Skripte bleiben"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Es gibt eine Rolle, die ein Ticket durch die Lanes traegt: jaira-worker. Welche Lane einen frischen Agenten, welches Modell und wie viele Agenten bekommt, steht an der Lane, nicht in der Wahl worker/lanes."
context: |-
  Entscheidung berkc 05.11.2026: 'es macht kein Sinn mehr worker und dispatcher zu haben, weil auf den Lanes festgelegt ist welche Agenten genutzt werden und wieviele und welches Model'.
  Heute: Setting dispatch = worker|lanes (~/.jaira/settings.json; G6T54H macht worker zum Default). jaira-dispatcher startet je Lane einen Agenten, jaira-worker traegt selbst und gibt nur fresh-context-Lanes ab (XKS048 macht critique/optimize fresh).
  spawn.sh, reap.sh und run-lane.sh liegen unter core/role/builtin/jaira-dispatcher/scripts/ und werden von Teamlead und Worker benutzt - sie muessen umziehen (z.B. nach jaira-worker/scripts), alte Pfade in Skills und Tests nachziehen.
  Erst nach G6T54H, XKS048 und CXH98W (signoff) bauen.
  role-acceptance bleibt (berkc 05.11.2026).
definition-of-done: "jaira-dispatcher ist aus core/role/builtin entfernt; das Setting dispatch verschwindet oder wird ignoriert; Release-Note '- !'"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-11-05T09:28:14Z
updated-at: 2026-11-05T09:28:14Z
---

# Nur noch jaira-worker: der Dispatcher faellt weg, seine Skripte bleiben

## Definition of Done

- [ ] jaira-dispatcher ist aus core/role/builtin entfernt; das Setting dispatch verschwindet oder wird ignoriert; Release-Note '- !'
- [ ] spawn.sh, reap.sh und run-lane.sh liegen an einem neuen Ort; teamlead, worker und Tests nutzen ihn; shipped-Hashes ergaenzt
- [ ] jaira update entfernt eine installierte jaira-dispatcher-Rolle oder sagt, dass sie veraltet ist; Test

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

