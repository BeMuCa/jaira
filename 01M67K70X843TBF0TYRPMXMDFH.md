---
id: 01M67K70X843TBF0TYRPMXMDFH
title: "Die Release-Notes erscheinen als Popup statt Seite, getrennt in Features zuerst und 'Andere' darunter, luftig und leicht lesbar"
status: testing
ready: true
creator: berkc
assignee: berkc
goal: "Nach einem Upgrade zeigt ein Popup die Aenderungen: oben Features und Feature-Aenderungen, darunter unter der Ueberschrift 'Andere' der Rest, mit Abstaenden, freundlich und schnell zu erfassen."
context: |-
  Wunsch von berkc am 31.10.2026 (Punkt 6): 'I wanted the release notes to be a pop up. You made them a page. check what i wrote in the ticket again and redo. i wanted pop up only and a split between features (or feature changes) and others. Also make it more fun and easily readable. You can use space and gaps ... Features and feature changes are the most relevant, so they come first and then a header with others.'
  Heute (0.5.7, GK1QGT): eine Version pro Seite, links/rechts blaettert; 'w' im Launcher. Lies GK1QGT und seine Notizen, was berkc urspruenglich verlangt hat.
  Die Notizdateien core/release/notes/<ID>.md haben heute nur '- ' und '- ! ' (Handlungsbedarf); eine Kennzeichnung Feature vs. anderes fehlt - braucht ein Format (z.B. '- + ' fuer Feature) und die Regeln in CLAUDE.md.
definition-of-done: "Popup (Overlay ueber dem Board), keine eigene Seite; links/rechts blaettert weiter zwischen Versionen"
tags:
  - tui
blocked-by: []
related:
  - 01M5659BR7RYZ6RVC84HGK1QGT
commits: []
created-at: 2026-10-30T23:44:40Z
updated-at: 2026-10-31T00:00:15Z
updated-by: berkc
claimed-by: berk-ca-2691736
claimed-at: 2026-10-30T23:47:40Z
outcome-what: "Release notes popup and jaira update: a note line '- + text' is a feature; the popup lists features under a Features header, then the rest under Others, a blank line after each change. Parser: core/release/release.go Feature(); layout: internal/tui/releasenotes.go lines()/changeLines(); CLI: internal/cli/update.go writeNotes prints '  + text'. CLAUDE.md format bullet and note core/release/notes/MXMDFH.md added. Commit cb13a6c2 on feat/MXMDFH-release-notes-popup."
outcome-why: "Berk: the popup must be the only surface, features first, Others below with a header, gaps for readability; the feature marker had to be a format the old binaries read as a plain line, so '- + ' mirrors '- ! '."
outcome-resolves: "DoD 2 (ordering, header, gaps, test) and DoD 3 (marker and CLAUDE.md) code and tests in place; DoD 1 (popup, paging) unchanged from GK1QGT, to be verified in tmux; DoD 4 note written."
review-summary: "none"
round: 2
review-gaps: "none: duplication searched (grep CutPrefix and the '! ' marker across the repo): the only marker parser is core/release/release.go Important/Feature, the other CutPrefix hits are unrelated (epic.go, selfupdate.go, milestone.go). No dead code: every added helper has callers (Feature: release.go lines and update.go writeNotes; changeLines: lines()). No fluff left after the e39d62f4 simplification. Cost: lines() rebuilds the page on each render, as before this change; the render already did it, and a page is a few dozen lines."
---

# Die Release-Notes erscheinen als Popup statt Seite, getrennt in Features zuerst und 'Andere' darunter, luftig und leicht lesbar

## Definition of Done

- [x] Popup (Overlay ueber dem Board), keine eigene Seite; links/rechts blaettert weiter zwischen Versionen
  proof: internal/tui/releasenotes.go:41-50 key() left/h and right/l inert at the ends; popup verified in tmux: 'What is new since 0.5.0 · 9 unread', Left -> '0.5.8 · 2/9', Right -> 'Unreleased · 1/9', second Right no-op, esc closes; TestNotesViewStartsOnTheNewestAndPagesOneVersionAtATime, TestPopupAfterAnUpgradeAcrossVersionsCountsTheUnreadAndPagesThem
- [x] Je Version zuerst Features/Feature-Aenderungen, dann Ueberschrift 'Andere'; Leerzeilen und Einrueckung fuer Lesbarkeit; Test
  proof: internal/tui/releasenotes.go:73-98 lines() groups Features then Others, Others header only beside a feature; :101-119 changeLines blank line after each change; TestNotesViewSplitsFeaturesFromOthers; tmux popup showed Features header, gap, '+' bullets, Others header
- [x] Notizformat kennzeichnet Features; CLAUDE.md-Abschnitt zu Release-Notes beschreibt es
  proof: core/release/release.go:109 Feature() ('- + ' marker, '- + ! ' combined); CLAUDE.md:270 format bullet; TestFeatureMarker, TestFeatureLineStaysAPlainChange (core/release/release_test.go); internal/cli/update.go writeNotes prints ' + text', TestUpdateNotesMarkImportantLines
- [x] Release-Note-Zeile
  proof: core/release/notes/MXMDFH.md two '- + ' lines, one per change, unwrapped

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/release: a '- + ' marker (alone or as '- + ! ') read by a new Feature helper; failing test in core/release/notes_test.go first
- [x] internal/tui/releasenotes.go lines(): Features group first under a 'Features' header, plain changes under 'Others' (header only when the version has a feature line), a blank line between changes, same colours for ! lines; failing test in releasenotes_test.go first
- [x] internal/cli/update.go writeNotes: a feature line prints '  + text'; test next to the existing writeNotes test
- [x] CLAUDE.md release-notes format list: a bullet for '- + ' next to '- ! ' (old binaries read it as a plain line)
- [x] core/release/notes/MXMDFH.md: one '- ' line per change, a line that tells the user what to expect in the popup and in 'jaira update'
- [~] go test on the touched packages, then go vet ./..., then full suite with the three hanging tests skipped
- [ ] run the built binary in tmux on a scratch repo (JAIRA_HOME in a temp dir), check the popup pages and the layout, describe what was seen

## Progress
- **2026-10-30 23:48 · berkc** — Open decisions, counted once (before the plan). 1 feature marker: '- + text' for a feature or feature change, combinable with the must-act marker as '- + ! text'. Settled: mirrors '- ! ' (same parser style, ASCII, a plain line to old binaries which show '+ text'). Rejected: a word prefix such as 'Feature:' (a sentence starting with that word would be misread); a separate file per kind (breaks one file per ticket); a unicode star (old binaries and terminals render it unevenly).
- **2026-10-30 23:48 · berkc** — 2 section header labels: 'Features' above the features, 'Others' above the rest. English, because the popup title and the hint are English ('What is new since', '←/→ older/newer'). The ticket's 'Andere' in its DoD is superseded by the person's own wording 'others'. 3 the 'Others' header is printed only when the version has at least one feature line: a version with only plain notes (every version before this change) looks exactly as it did. 4 a must-act feature line sits under Features and keeps its warning colour; a must-act plain line sits under Others as before. 5 paging stays as GK1QGT decided: the popup pages only unread versions, 'w' pages all, newest first, left/h older, right/l newer. Settled, not reopened.
- **2026-10-30 23:48 · berkc** — 6 jaira update (CLI, writeNotes) prints a feature line as '  + text', like it prints '  ! text' for important lines; JSON keeps the raw string as it keeps '! text' today. That is a changed CLI output and gets its own release-note line. 7 spacing: one blank line between changes in the popup, the person asked for gaps. Settled. 8 the open point (who decides) : none needs the person; nothing is set to conversational, the person did not ask for it. Plan follows in pre-process.
- **2026-10-30 23:49 · berkc** — Plan rests on the decisions already noted above. Test order: write the failing core/release and TUI tests first, then the code. The popup layout is verified by tmux, not only by the test, because the person judges it by eye. The release note line is written last so it describes what was built.
- **2026-10-30 23:50 · berkc** — Lane change pre-process -> in-progress was refused by the send step ('was not sent', ref still at pre-process, 23:49:13, the same worker's own earlier write) while the local file already said in-progress. Re-sending with this note. If the ref still says pre-process after this, the send is broken, not the ticket, and the next worker must check it before moving.
- **2026-10-30 23:53 · berkc** — critique: internal/tui/releasenotes.go changeLines(c, feature, width) takes the feature flag from lines() and lines() strips the marker before calling it. Simplify: changeLines(c, width) parses both markers itself with release.Feature and release.Important; lines() only asks Feature to pick the group. Checked: the wrap, the bullets and the blank line are right; CLI writeNotes and the core helper are fine as they are.
- **2026-10-30 23:54 · berkc** — critique pass 2: none. Stood on the decisions note (format '- + ', English labels, Others only beside a feature, CLI '+' lines) and on the simplification note of pass 1, now in e39d62f4: changeLines reads both markers itself. Checked: Feature/Important split in core/release/release.go is the same pattern as Important; CLI writeNotes changes only the feature branch; no new file beyond core/release/notes/MXMDFH.md and the CLAUDE.md bullet.
