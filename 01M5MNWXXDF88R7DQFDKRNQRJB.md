---
id: 01M5MNWXXDF88R7DQFDKRNQRJB
title: "Die Warnung 'Board von aelterer Version eingerichtet' kommt einmal, nicht bei jedem Aufruf"
status: todo
ready: true
creator: BeMuCa
goal: "Wer jaira update noch nicht laufen liess, hoert es einmal pro Board und Version - nicht vor jeder Ausgabe, wo Agenten es mitparsen und Menschen es ueberlesen."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkt 11): 'jaira: this board was set up by an older version of jaira — run jaira update' steht vor jeder Ausgabe, auch bei --json auf stderr. Agenten filtern es weg (grep -v), Menschen sehen es nicht mehr.
  Es gibt schon einen Versionsstempel pro Board (core/release/release.go: Stamped/Stamp). Der Hinweis kann sich daran merken, dass er fuer diese Version gezeigt wurde.
  Verwandt: T8A8KM (Release-Pille), PWAF22 (update bringt das Setup nach).
definition-of-done: Die Warnung erscheint einmal pro Board und Binary-Version (Merker neben dem Versionsstempel) und danach nur noch bei jaira update selbst oder --verbose; ein Test deckt den zweiten Aufruf ab
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:26:01Z
updated-at: 2026-10-26T20:06:48Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: rnqrjb-driver
claimed-at: 2026-10-26T20:06:42Z
---

# Die Warnung 'Board von aelterer Version eingerichtet' kommt einmal, nicht bei jedem Aufruf

## Definition of Done

- [ ] Die Warnung erscheint einmal pro Board und Binary-Version (Merker neben dem Versionsstempel) und danach nur noch bei jaira update selbst oder --verbose; ein Test deckt den zweiten Aufruf ab
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

