---
id: 01M4TRA5NAZ2BWSH0XWD64GNZT
title: "Die Lanes human und signoff sagen mit ihrem Namen, wofuer sie da sind"
status: backlog
ready: false
creator: BeMuCa
goal: "Ein Agent und ein Mensch erkennen am Namen, dass human eine offene Frage ist und signoff die Abnahme."
context: |-
  Die Lane human heisst auf dem Board 'HITL', die Lane signoff heisst 'Human Review' (core/lane/builtin/30-human.md:3, core/lane/builtin/45-signoff.md:3).
  Der Agent-Block sagt 'The human review lane cannot be left by an agent — a person accepts the work there.' (core/board/announce.go:120). Gemeint ist signoff, ein Agent liest 'human'.
  Folge am 2026-10-13: Claude hat Berk erklaert, human sei die Abnahme. Falsch - human verlangt nur eine Frage (requires-question), die Abnahme ist signoff (requires-human-exit).
  Berk sah PNC854 in human und fragte 'da ist doch gar keine Frage?'.
  Vorschlag aus der Session: human heisst 'Question', signoff heisst 'Sign-off', und announce.go:120 nennt die Lane-ID statt 'the human review lane'. Die Namen darf brainstorm noch aendern.
definition-of-done: "Die eingebauten Lanes zeigen human als 'Question' und signoff als 'Sign-off' (core/lane/builtin/30-human.md, 45-signoff.md); was ein Board mit eigener Kopie in .jaira/lanes/ sieht, ist entschieden und notiert"
tags:
  - gates
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-10-13T13:47:56Z
updated-at: 2026-10-13T14:17:04Z
assignee: BeMuCa
updated-by: BeMuCa
---

# Die Lanes human und signoff sagen mit ihrem Namen, wofuer sie da sind

## Definition of Done

- [ ] Die eingebauten Lanes zeigen human als 'Question' und signoff als 'Sign-off' (core/lane/builtin/30-human.md, 45-signoff.md); was ein Board mit eigener Kopie in .jaira/lanes/ sieht, ist entschieden und notiert
- [ ] core/board/announce.go:120 nennt die Lane per ID (signoff) statt 'the human review lane'; ein Test in core/board prueft den Satz; der CLAUDE.md-Block dieses Repos ist neu erzeugt
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

