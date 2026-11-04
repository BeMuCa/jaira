---
id: 01M6M72CGAVYB8WR1YZBN5GSWM
title: Die Ticket-ID auf der Karte ist gut lesbar statt fast unsichtbar grau
status: optimize
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
updated-at: 2026-11-04T21:36:26Z
updated-by: berkc
claimed-by: berk-ca-1280093
claimed-at: 2026-11-04T21:24:26Z
outcome-what: "The ticket handle on board cards and ticket headers is toned by terminal background: 248 on dark, 240 on light (was colFaint 240 everywhere)."
outcome-why: "On the dark selection fill (239) colFaint 240 is ten grey levels away, so the handle vanished on the selected card. The lane counter keeps colFaint."
outcome-resolves: "Ticket ID hard to read on dark terminals and on the selected card (berkc, 04.11.2026)."
review-summary: none
review-gaps: "handle_test.go: replaced the inline absolute-value branch with the package's existing abs() (links_test.go:543). Left: grey() in handle_test.go (takes a color.Color, paletteRGB takes an index - they differ); handleStyle() building a fresh lipgloss style per call (render path, cheap, mirrors selectionFill/laneShade); no dead code - styHandle removed with all 7 callers switched."
---

# Die Ticket-ID auf der Karte ist gut lesbar statt fast unsichtbar grau

## Definition of Done

- [x] styHandle nutzt eine hellere Farbe (z.B. 245-250 oder eine adaptive Farbe hell/dunkel), die auf der ausgewaehlten und der normalen Karte lesbar ist; Test in internal/tui
  proof: internal/tui/view.go handleStyle(); internal/tui/handle_test.go TestHandleToneReadsOnEveryCardFill
- [x] Release-Note-Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/N5GSWM.md:1

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] view.go: replace package var styHandle with method (m *Model) handleStyle() that picks handleDark=248 / handleLight=240 off m.darkBG, like laneShade does; keep colFaint for lane counter
- [x] switch all 7 styHandle call sites (view.go x3, lanefocus.go, edit.go, signoff.go, layerdraw.go) to m.handleStyle(); all of them show a ticket handle, so they should all follow
- [x] test in internal/tui: for dark and light, the handle tone is at least 60 grey levels (paletteRGB) away from both lane shades and the selection fill
- [x] release note core/release/notes/N5GSWM.md, one line
- [x] go test ./internal/tui -count=1

## Progress
- **2026-11-04 21:26 · berkc** — plan: no single grey reads on both selection fills (dark 239=78, light 250=188 grey). 240 on dark selection is 10 levels apart = invisible, which is the bug. Fixed 245 gives only ~50 on either; adaptive gives 90 (dark 248 on 239) and keeps light at 240 (100 on 250), which already read fine. 248 not 250 so the ID stays calmer than the title and only a step above the name (244). Every styHandle use is a ticket handle beside a title, so all follow; lane counter keeps colFaint via styLaneCount.
- **2026-11-04 21:31 · berkc** — Claim: jaira claim said berk-ca-1280093 holds it (3m); did not --steal, worked in the worktree. Mutation check: dark handle back at 240 fails TestHandleToneReadsOnEveryCardFill (88 is 40 from shade B 48; 240 on selection 239 is 10 off, so the selection fails too). Lane-counter stays colFaint; only handle sites switched. Test helper grey() and the RGBA reading live in handle_test.go; abs() already exists in links_test.go, so none added.
- **2026-11-04 21:34 · berkc** — critique pass 1: no findings. handleStyle() (view.go) mirrors selectionFill (palette.go:43) and laneShade (view.go) - a darkBG switch over two consts - so no new pattern; a package var mutated on BackgroundColorMsg (model.go:1101) would be the only alternative and is worse. All 7 sites show a ticket handle beside a title; lane counter correctly keeps colFaint via styLaneCount. Checked logbook card fills too (logbook.go:27-30: 23/24 dark, 195/153 light), not covered by the test but far from 248/240, no regression. Let stand: test comment says every fill is grey, true for the three fills it checks.
- **2026-11-04 21:36 · berkc** — optimize: only change was abs() reuse in handle_test.go. Considered precomputing two package-level handle styles instead of NewStyle per call - not worth it, per-card cost is negligible and laneShade/selectionFill set the per-call pattern.
