---
id: 01M59MQBJHB2NW3F3WFXAW9KVT
title: "Human Review zeigt zuerst, wie der Mensch die Aenderung selbst testet"
status: in-progress
ready: true
creator: berkc
assignee: berkc
goal: "Wer ein Ticket in Human Review oeffnet, sieht oben kurze Schritte zum Ausprobieren (was starten, was druecken, was sehen) und erst danach die Review-Details."
context: |-
  Gemeldet von berkc am 19.10.2026: 'die Human Review Lane ist schlecht, frueher standen da Beschreibungen, wie man testet'.
  Die Daten gibt es: die review-Lane muss review-check liefern (.jaira/lanes/review.md, output-produces). Der Sign-off-Screen zeigt es (internal/tui/signoff.go:83), das Detail auch (view.go:1375) - aber als letzten Punkt nach summary, gaps, verdict.
  Beispiel GK1QGT: check ist ein Block mit 9 Schritten voller Pfade wie /home/berkc/code/X_General/jaira-T8A8KM und go build -ldflags ... - fuer einen Entwickler, nicht zum schnellen Ausprobieren.
  Zu klaeren beim Bearbeiten: was berkc 'frueher' gesehen hat (git log von internal/tui/signoff.go und .jaira/lanes/review.md pruefen).
  Ausserdem: der Code liegt in einem anderen Worktree/Branch. Ohne einen fertigen Befehl zum Bauen dieses Branches kann man nichts testen (siehe X3A8RZ: berkc sah den Fix nicht, weil er nur im Branch ist).
definition-of-done: "Der Sign-off-Screen und das Detail einer Karte in Human Review zeigen review-check als ersten Block, als nummerierte Schritte; Test in internal/tui"
tags:
  - tui
blocked-by: []
related: []
commits:
  - 26d4a26af9b1c7e56b09175150b0af34663ceab9
created-at: 2026-10-19T08:33:51Z
updated-at: 2026-11-04T21:42:39Z
updated-by: berkc
claimed-by: berk-ca-1280084
claimed-at: 2026-11-04T21:24:26Z
outcome-what: "Sign-off screen and detail of a human-lane ticket open with a 'try it' row of numbered steps (one per line); review-check prompt asks for at most six steps, step 1 builds and starts this change."
outcome-why: "A reviewer's check was the last row, as a developer recipe; a person waiting to sign off could not quickly try the change, and step 1 pointed at the installed binary, not this branch."
outcome-resolves: "Human review shows how to test: the check is the first thing on the sign-off screen and the ticket detail, with short steps that name what to see."
review-summary: |-
  internal/tui/view.go:1347-1354 re-implements the declared closure (value -> sourced row, empty+owed -> owedRow) that the same function defines at ~1415; hoist declared (with owed/srcs) above the identity rows and call declared(declaredField{"try it", ticket.FieldReviewCheck, oneStepPerLine(t.ReviewCheck)}) instead of a second copy
  internal/tui/owed_test.go:336 strings.Index(out,"try it") > strings.Index(out,"secrets-status") passes when the try it row is missing (-1); add an i < 0 failure as reviewcheck_test.go does, so the ordering assertion cannot pass vacuously
  internal/tui/signoff.go:19 the reflowed doc comment left one 104-char line; rewrap the paragraph to the file's ~80-column width
round: 2
---

# Human Review zeigt zuerst, wie der Mensch die Aenderung selbst testet

## Definition of Done

- [x] Der Sign-off-Screen und das Detail einer Karte in Human Review zeigen review-check als ersten Block, als nummerierte Schritte; Test in internal/tui
  proof: internal/tui/signoff.go:71 (try it row first); view.go:1359 (detail, human lanes); tests reviewcheck_test.go TestSignOffScreenShowsTheCheck, TestDetailPaneOpensWithTheCheckOnAHumanLane, TestDetailPaneKeepsTheCheckInReviewOffAHumanLane
- [x] Der Prompt der review-Lane verlangt fuer check: hoechstens ~6 Schritte, beginnend mit einem kopierbaren Befehl, der den Branch baut und startet, und je Schritt was man sehen soll
  proof: core/lane/builtin/40-review.md:63 (at most six steps, step 1 builds and starts this change); identical copy .jaira/lanes/review.md; core/shipped/shipped.sha256 last line ae89... kept, fe8fd6a9... appended; go test ./core/shipped ./core/lane ok
- [x] Release-Note-Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/AW9KVT.md (one '- ' line)

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] internal/tui: helper checkSteps(text) splits an inline '1. a  2. b' review-check into one line per step (text that already has line breaks keeps them); used by both screens
- [x] signoff.go: move the declared("check") call from last to first (before tags), label 'try it', rendered as numbered steps; nothing else in signoff.go changes (EH6Z0N merge)
- [x] view.go detailBody: for a ticket in a lane with RequiresHumanExit or RequiresQuestion, render review-check as first block after the title rule ('try it'), and leave it out of the Review group there; other lanes unchanged
- [x] tests in internal/tui: signoff and detail show check before goal/summary, one step per line; detail for a non-human lane keeps check in Review
- [x] core/lane/builtin/40-review.md + .jaira/lanes/review.md: review-check at most ~6 steps, step 1 a copyable command that builds and starts this branch (worktree path), each step names what you should see; append new hash to core/shipped/shipped.sha256, run core/shipped tests
- [x] release note core/release/notes/AW9KVT.md; go test ./internal/tui ./core/shipped ./core/lane

## Progress
- **2026-11-04 21:26 · berkc** — Plan decisions: (1) 'frueher' was not a different order: git log -S shows review-check was last on both screens since it was introduced in 0994fc07 (2026-08-20); what changed is that check texts grew into 9-step developer recipes (GK1QGT). So the fix is order + step layout + a tighter lane prompt. (2) The detail pane moves check first only for tickets in a person's lane (requires-question or requires-human-exit), not everywhere: on other lanes nobody is asked to try anything. (3) Inline '1. x  2. y' is split into lines at render time, the stored field is not rewritten — the file format stays as written. (4) EH6Z0N (branch feat/EH6Z0N-diff-view, unmerged) also edits internal/tui/signoff.go (adds a right-half diff pane). This ticket keeps its signoff.go change to moving one call up in renderSignOff's left column plus the step layout, so the merge conflict is at most one hunk.
- **2026-11-04 21:27 · berkc** — Name clash: internal/tui/reviewcheck_test.go already has package-level const checkSteps (test fixture), so the helper is named oneStepPerLine. Also TestSignOffScreenShowsTheCheck asserts check comes AFTER verdict; that is the behaviour being reversed, so the test gets inverted.
- **2026-11-04 21:34 · berkc** — Sign-off/detail: the 'try it' row reuses declared()/row() so it keeps the '(review)' provenance prefix on step 1 (same as every other row). Two owed_test assertions moved check->'try it' on the sign-off screen; the detail-pane owed assertion for lane review stays 'check' (non-human lane).
- **2026-11-04 21:39 · berkc** — critique (pass 1): design stands — render-time split (note 3), human-lane-only placement (note 2), prompt rewrite are accepted. Three fixes: (1) view.go detailBody duplicates the declared closure for the try-it row; hoist declared+owed+srcs above the identity rows and call it. (2) owed_test.go:336 ordering check passes vacuously if 'try it' is absent; guard i<0. (3) signoff.go:19 doc comment line is 104 chars after reflow; rewrap. No other findings; next pass reads only these three.
- **2026-11-04 21:42 · berkc** — critique findings fixed: (1) view.go: declared closure hoisted above the try-it row, row now declared(declaredField{"try it",...}); moved view.go:1351 -> 1359, re-ticked DoD 1 proof. (2) owed_test.go: guard fails when 'try it' index < 0. (3) signoff.go:19-22 doc comment rewrapped under 80. signoff.go:236 and :250 are >100 chars and pre-existing, left alone. go test ./internal/tui -count=1 exit 0; go vet and gofmt clean.
