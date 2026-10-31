---
id: 01M67RWZDWBS8E44Z3KGXKS048
title: "critique und optimize laufen je in einem frischen Agenten, nicht im Agenten, der gebaut hat"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Die Lanes critique und optimize tragen context: fresh in der ausgelieferten Lane und auf diesem Board, so dass Worker und Dispatcher sie je einem neuen Agenten geben - zwei zusaetzliche Augenpaare statt Selbstkontrolle."
context: |-
  Was falsch ist: critique und optimize macht heute derselbe Agent, der in in-progress gebaut hat. Er findet seine eigenen Luecken nicht.
  Ausloeser: Experiment 3CBCN6 (31.10.2026): der Nur-Haiku-Lauf liess eine veraltete Hilfezeile und einen falschen Kommentar stehen; seine eigene critique/optimize fand das nicht. Der Maintainer: 'setz das bitte auch um das 2 agenten starten sollen fuer critique und optimize.'
  Bekannt: review hat schon context: fresh (core/lane/lane.go:60-76, ContextFresh). jaira-worker gibt eine fresh-Lane an einen neuen Agenten (core/role/builtin/jaira-worker/SKILL.md).
  Offen fuer den Plan: kostet zwei Agentenstarts mehr pro Ticket; beide Lanes sind tier strong.
definition-of-done: "core/lane/builtin critique und optimize und .jaira/lanes/critique.md, optimize.md tragen context: fresh; jaira show <id> --for-lane critique --json zeigt das"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-31T01:24:02Z
updated-at: 2026-10-31T01:24:02Z
---

# critique und optimize laufen je in einem frischen Agenten, nicht im Agenten, der gebaut hat

## Definition of Done

- [ ] core/lane/builtin critique und optimize und .jaira/lanes/critique.md, optimize.md tragen context: fresh; jaira show <id> --for-lane critique --json zeigt das
- [ ] jaira update --lanes bringt es auf ein bestehendes Board (oder die Release-Notiz sagt, was zu tun ist)
- [ ] Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

