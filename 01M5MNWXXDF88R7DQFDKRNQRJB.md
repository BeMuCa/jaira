---
id: 01M5MNWXXDF88R7DQFDKRNQRJB
title: "Die Warnung 'Board von aelterer Version eingerichtet' kommt einmal, nicht bei jedem Aufruf"
status: critique
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
updated-at: 2026-10-26T20:09:30Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: rnqrjb-driver
claimed-at: 2026-10-26T20:09:15Z
outcome-what: "nudgeIfStale (internal/cli/update.go) prints the 'board was set up by an older version' hint once per board and binary version: a 'nudged' marker file beside the version stamp (core/release Nudged/MarkNudged) remembers which version it was shown for. jaira update and the stamp are unchanged."
outcome-why: "The hint appeared on stderr before every command, even with --json, so agents filtered it away and people stopped reading it."
outcome-resolves: "DoD 1: the hint no longer repeats on every call; a new binary version shows it once again; afterwards jaira update's own report (previous version + what changed) is where staleness surfaces. Covered by TestNudgeIfStaleShowsOncePerBoardAndVersion and TestUpdateAfterShownHintStillReportsPreviousAndStamps. No --verbose flag exists, so none was added."
---

# Die Warnung 'Board von aelterer Version eingerichtet' kommt einmal, nicht bei jedem Aufruf

## Definition of Done

- [x] Die Warnung erscheint einmal pro Board und Binary-Version (Merker neben dem Versionsstempel) und danach nur noch bei jaira update selbst oder --verbose; ein Test deckt den zweiten Aufruf ab
  proof: core/release/release.go Nudged/MarkNudged; internal/cli/update.go nudgeIfStale; covered by TestNudgeIfStaleShowsOncePerBoardAndVersion, TestUpdateAfterShownHintStillReportsPreviousAndStamps
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/release/release.go: beside Stamped/Stamp add a marker file in the same state dir that records which binary version the stale-board hint was last shown for (read swallowing errors like Stamped; write like Stamp)
- [x] internal/cli/update.go nudgeIfStale: after the stamp check, stay silent when the marker already names release.Current; otherwise print the line and write the marker. No --verbose flag exists, none is added; jaira update keeps reporting the previous version and what changed
- [x] internal/cli/update_test.go: second nudge on the same board and version is silent; a different binary version nudges once again; jaira update after a shown hint still reports previous version and stamps
- [x] core/release/notes/RNQRJB.md: one '- ' line written as an instruction; update the nudgeIfStale doc comment

## Progress
- **2026-10-26 20:07 · BeMuCa** — Ticket opts out of brainstorm and planning (Options unchecked), so next_lane jumps todo -> in-progress; the driver wrote the plan. Interpretation taken: 'danach nur noch bei jaira update selbst' = after the one shown hint, the staleness only surfaces through jaira update's own report (previous version + what changed); update already suppresses the nudge via updateInProgress. There is no --verbose flag (global flags are -C and --json), so none is added.
- **2026-10-26 20:09 · BeMuCa** — Marker is a separate file <stateDir>/nudged holding the binary version; the version stamp itself is not touched, so jaira update still reads the true previous version. Marker write errors are swallowed (hint then just repeats).
