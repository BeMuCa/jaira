---
id: 01M4TRASJQKBMA0RQ5MPQNQ0RP
title: "Der Agent-Block sagt, dass ein Agent ein Ticket aus human weiterschiebt, sobald die Frage beantwortet ist"
status: backlog
ready: true
creator: BeMuCa
goal: "Was der erzeugte Agent-Block ueber die human-Lane sagt, stimmt mit ihrem Gate ueberein."
context: |-
  Der erzeugte Agent-Block (CLAUDE.md, AGENTS.md) nennt human 'a person's, not yours — you may move work in, never out' (core/board/announce.go:210-211).
  Das Gate sagt etwas anderes: human hat nur requires-question, kein requires-human-exit (core/lane/builtin/30-human.md:8). Ein Agent darf human also verlassen.
  Derselbe Block sagt weiter unten 'until it sits in a human lane, then continue once the human has answered' - also doch raus. Der Block widerspricht sich selbst.
  Ursache: announce.go:209 legt f.HumanExit und f.Question in einen gemeinsamen case.
  Folge am 2026-10-13: Nach Berks Antwort auf PNC854 wurde der erste Move aus human abgelehnt (von Claude Codes Klassifizierer; dass dieser Satz der Grund war, ist eine Vermutung).
  Verwandt: 64GNZT (Namen von human und signoff).
definition-of-done: "Fuer eine Lane mit nur requires-question sagt der Block: eine Person beantwortet die Frage, danach schiebt der Agent weiter; fuer requires-human-exit bleibt 'never out'; ein Test in core/board deckt beide Faelle"
tags:
  - gates
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-10-13T13:48:17Z
updated-at: 2026-10-13T14:16:38Z
assignee: BeMuCa
updated-by: BeMuCa
---

# Der Agent-Block sagt, dass ein Agent ein Ticket aus human weiterschiebt, sobald die Frage beantwortet ist

## Definition of Done

- [ ] Fuer eine Lane mit nur requires-question sagt der Block: eine Person beantwortet die Frage, danach schiebt der Agent weiter; fuer requires-human-exit bleibt 'never out'; ein Test in core/board deckt beide Faelle
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/; der CLAUDE.md-Block dieses Repos ist neu erzeugt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

