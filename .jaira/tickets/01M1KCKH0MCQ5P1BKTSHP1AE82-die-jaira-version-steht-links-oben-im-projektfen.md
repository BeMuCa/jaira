---
id: 01M1KCKH0MCQ5P1BKTSHP1AE82
title: Die jaira-Version steht links oben im Projektfenster
status: backlog
ready: true
creator: BeMuCa
goal: "Das TUI zeigt die laufende Binary-Version sichtbar im Board, links oben in der Ecke des Projektfensters"
context: "Berk am 03.09.: beim Arbeiten mit mehreren Binary-Staenden (self upgrade, go build, ~/.local/bin) ist nicht sichtbar, welche Version laeuft. Wunsch: Version ganz links oben in der Ecke des Projektfensters. ACHTUNG Kollision: DNAEPN (done, 03.09.) hat gerade entschieden, dass die bestehende Versionszeile (internal/tui/updatecheck.go, gezeichnet in home.go und view.go/Fusszeile) auf dev-Builds SCHWEIGT, weil 'jaira dev' nichts beantwortet; auf Releases zeigt sie 'jaira 0.1.1 - up to date'. Vor dem Bauen klaeren: Platzierung links oben zusaetzlich oder statt der Fusszeile, und ob ein dev-Build 'dev' zeigen soll (widersprueche DNAEPN) oder weiter nichts."
definition-of-done: "Das Board zeigt links oben die Version; ein Dev-Build zeigt 'dev'; ein Test deckt die Platzierung ab"
tags: []
blocked-by: []
commits: []
created-at: 2026-09-03T10:21:34Z
updated-at: 2026-09-03T10:28:37Z
assignee: BeMuCa
updated-by: BeMuCa
---

# Die jaira-Version steht links oben im Projektfenster

## Definition of Done

- [ ] Das Board zeigt links oben die Version; ein Dev-Build zeigt 'dev'; ein Test deckt die Platzierung ab

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

