---
id: 01M6M72CGAVYB8WR1YZBN5GSWM
title: Die Ticket-ID auf der Karte ist gut lesbar statt fast unsichtbar grau
status: in-progress
ready: true
creator: berkc
assignee: berkc
goal: Die sechsstellige Ticket-ID auf jeder Board-Karte ist auf dunklem und hellem Terminal auf einen Blick lesbar.
context: |-
  Was falsch ist: die ID (z.B. 6TCTFV) auf der Karte ist kaum zu erkennen. berkc am 04.11.2026 mit Screenshot: 'die ist schwer zu erkennen'.
  Ursache: internal/tui/view.go:63 styHandle = Foreground(colFaint), colFaint = Color("240") (view.go:45) - dasselbe Grau wie Rahmen und Lane-Zaehler.
  Der Name daneben (berkc) ist heller, die ID wirkt dadurch wie ausgegraut.
  Nicht gewollt: eine Farbe pro Ticket oder eine neue Einstellung - nur ein lesbarer, ruhiger Ton, der nicht mit Tag- und Lane-Farben konkurriert.
definition-of-done: "styHandle nutzt eine hellere Farbe (z.B. 245-250 oder eine adaptive Farbe hell/dunkel), die auf der ausgewaehlten und der normalen Karte lesbar ist; Test in internal/tui"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-11-04T21:22:33Z
updated-at: 2026-11-04T21:31:24Z
updated-by: berkc
claimed-by: berk-ca-1280093
claimed-at: 2026-11-04T21:24:26Z
---

# Die Ticket-ID auf der Karte ist gut lesbar statt fast unsichtbar grau

## Definition of Done

- [ ] styHandle nutzt eine hellere Farbe (z.B. 245-250 oder eine adaptive Farbe hell/dunkel), die auf der ausgewaehlten und der normalen Karte lesbar ist; Test in internal/tui
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] view.go: replace package var styHandle with method (m *Model) handleStyle() that picks handleDark=248 / handleLight=240 off m.darkBG, like laneShade does; keep colFaint for lane counter
- [x] switch all 7 styHandle call sites (view.go x3, lanefocus.go, edit.go, signoff.go, layerdraw.go) to m.handleStyle(); all of them show a ticket handle, so they should all follow
- [ ] test in internal/tui: for dark and light, the handle tone is at least 60 grey levels (paletteRGB) away from both lane shades and the selection fill
- [ ] release note core/release/notes/N5GSWM.md, one line
- [ ] go test ./internal/tui -count=1

## Progress
- **2026-11-04 21:26 · berkc** — plan: no single grey reads on both selection fills (dark 239=78, light 250=188 grey). 240 on dark selection is 10 levels apart = invisible, which is the bug. Fixed 245 gives only ~50 on either; adaptive gives 90 (dark 248 on 239) and keeps light at 240 (100 on 250), which already read fine. 248 not 250 so the ID stays calmer than the title and only a step above the name (244). Every styHandle use is a ticket handle beside a title, so all follow; lane counter keeps colFaint via styLaneCount.
