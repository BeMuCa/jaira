---
id: 01M5MNX3NGYBHKCQ2YBQ4N1B66
title: Die Testing-Lane verlangt keinen vollen Suite-Lauf mehr; die Suite laeuft einmal vor dem PR
status: in-progress
ready: true
creator: BeMuCa
goal: "Der Prompt der Testing-Lane sagt, was die Regel ist: Paket-Tests der beruehrten Pakete in der Lane, die volle Race-Suite einmal vor dem PR durch die Hauptsession - und kein Worker startet sie mehr aus eigenem Antrieb."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkt 15): ein Worker lief die volle Suite, obwohl die Regel 'einmal pro Milestone' lautet. Die Regel stand nur im Gedaechtnis der Session, nicht im Lane-Prompt.
  Nachgeprueft am 20.10.: .jaira/lanes/testing.md:22 dieses Boards verlangt ausdruecklich 'run the whole suite ... with the race detector'. Prompt und Regel widersprechen sich; die Regel stammt vom 14.10. ('Test nur final bevor wir pr preparieren'), weil die Suite 7-8 Minuten dauert.
  Betrifft auch die Katalog-Lane lanes/testing.md im Repo, falls sie denselben Satz traegt.
definition-of-done: ".jaira/lanes/testing.md und die Katalog-Lane sagen: in der Lane die Tests der beruehrten Pakete (go test ./<pkg> -count=1), die volle Suite mit -race einmal vor dem PR durch die Hauptsession; der Satz 'run the whole suite' ist weg"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:26:07Z
updated-at: 2026-10-27T00:10:45Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: 4n1b66-driver
claimed-at: 2026-10-27T00:01:46Z
outcome-what: "Rewrote the Gates paragraph of the testing lane in .jaira/lanes/testing.md and core/lane/builtin/27-testing.md (byte-identical prompt): a lane builds, runs format/vet and tests only the touched packages (Go: go test ./<pkg> -count=1 plus gofmt/vet/build); the example verdict now says 'touched packages green'. Added the rule as one line to the jaira-worker and jaira-dispatcher role skills, and a release note in core/release/notes/4N1B66.md."
outcome-why: "The testing lane told agents to run the whole suite with the race detector, but the full suite takes 7-8 minutes and the maintainer's rule is that no lane runs it: the main session runs go test ./... -race once, before the pull request. The rule lived only in session memory, so a worker followed the prompt instead."
outcome-resolves: "DoD 1: 'run the whole suite' is gone and both lane copies carry the new Gates rule with an identical prompt (.jaira/lanes/testing.md:22, core/lane/builtin/27-testing.md:23). DoD 2: the rule is one line in core/role/builtin/jaira-worker/SKILL.md:168 and core/role/builtin/jaira-dispatcher/SKILL.md:268. DoD 3: core/release/notes/4N1B66.md tells readers to run jaira lanes use testing --force and jaira roles install."
review-summary: |-
  core/role/builtin/jaira-role-pr/SKILL.md:37 ("Before you push anything") - the diff forbids the full suite in every lane and says "the main session" runs it once before the PR, but no shipped file tells any session to run it: the old Gates paragraph was the only shipped instruction to run a full suite, and "main session" appears in no role file. Either add a fourth precondition there (the repository's full suite green on the head being pushed; Go: go test ./... -race) or keep it as the maintainer's own practice - a maintainer decision, see the human-lane question.
  core/lane/builtin/27-testing.md:25-27 and .jaira/lanes/testing.md:24-26 - "with the race detector (go test ./... -race)" is stated for every repository, though the progress note of 2026-10-27 00:02 says the Go command is phrased as the Go case because the builtin ships to non-Go repositories; move it into the Go case, e.g. "Never run the full suite in a lane: it runs once, before the pull request, by the main session (in Go `go test ./... -race`)." - same bytes in both files (lanesEquivalent compares Prompt).
  core/role/builtin/jaira-dispatcher/SKILL.md:268 and core/role/builtin/jaira-worker/SKILL.md:168-169 - "it takes 7-8 minutes" is this repository's suite, written into role files that jaira roles install puts into every repository (--global: all of them), and go test ./... -race is again unconditional; drop the number ("the full suite is the slowest check, so it runs once ...") and name go test ./... -race as the Go case, as the lane does.
  core/release/notes/4N1B66.md:1 - jaira lanes use testing --force copies the catalogue's testing.md when there is one (the catalogue lays over the built-ins, core/lane/lane.go:682), and the maintainer's catalogue ~/.jaira/lanes/testing.md:22 still says "Build everything and run the whole suite"; on a scratch board with that file, use --force wrote the old Gates back (board testing.md:22). Add to the note: if jaira lanes path shows a catalogue holding testing.md, update or delete that copy first. The catalogue file itself is outside the repository and needs the maintainer's own hand - it is also what "die Katalog-Lane" in DoD 1 means in jaira's own vocabulary (jaira lanes path: "Catalogue").
round: 2
---

# Die Testing-Lane verlangt keinen vollen Suite-Lauf mehr; die Suite laeuft einmal vor dem PR

## Definition of Done

- [x] .jaira/lanes/testing.md und die Katalog-Lane sagen: in der Lane die Tests der beruehrten Pakete (go test ./<pkg> -count=1), die volle Suite mit -race einmal vor dem PR durch die Hauptsession; der Satz 'run the whole suite' ist weg
  proof: .jaira/lanes/testing.md:22, core/lane/builtin/27-testing.md:23
- [x] jaira-worker/SKILL.md und jaira-dispatcher/SKILL.md tragen dieselbe Regel in einer Zeile
  proof: core/role/builtin/jaira-worker/SKILL.md:168, core/role/builtin/jaira-dispatcher/SKILL.md:268
- [x] Release-Notiz unter core/release/notes/
  proof: core/release/notes/4N1B66.md:1

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-27 00:02 · BeMuCa** — Gates phrases the Go command as the Go case because core/lane/builtin ships to every repository, not only Go ones; the lane prompt must stay byte-identical between .jaira/lanes/testing.md and the builtin (lanesEquivalent). Example verdict now says 'touched packages green'.
- **2026-10-27 00:09 · BeMuCa** — critique: four findings, the full text is in review-summary. (1) Open question for the maintainer: the "runs once before the PR" half of the rule has no home - no shipped role tells any session to run the full suite now (the old Gates was the only place), and "main session" is defined nowhere; the candidate home is the "Before you push anything" list in core/role/builtin/jaira-role-pr/SKILL.md:37. Wait for the answer before touching it. (2) Lane Gates, both files, same bytes: put "with the race detector (go test ./... -race)" inside the Go case, as the 00:02 note intends. (3) Worker SKILL.md:168 and dispatcher SKILL.md:268: drop "7-8 minutes" (true only for this repository; the roles are installed everywhere) and name go test ./... -race as the Go case. (4) Release note: lanes use --force copies a catalogue testing.md over the builtin (core/lane/lane.go:682); the maintainer's ~/.jaira/lanes/testing.md:22 still has the old Gates - reproduced on a scratch board, use --force wrote "Build everything and run the whole suite" back. Tell readers to update or delete a catalogue copy first; the maintainer refreshes his own copy by hand (outside the repo). Left alone: the rule appearing in three files and the dispatcher line being unwrapped - both are what DoD 2 asks for ("in einer Zeile").
- **2026-10-27 00:10 · BeMuCa** — driver: critique finding 1 (nothing shipped tells a session to run the full suite; jaira-role-pr) is not taken into this ticket and does not block it: it is outside DoD 1-3, and .github/workflows/ci.yaml:70 already runs go test ./... -race on every PR. It goes to the maintainer as a question in the driver's report. Findings 2-4 go back to in-progress. The maintainer's own catalogue ~/.jaira/lanes/testing.md:22 still has the old Gates text and lays over the shipped lane; it is outside the repository and is left for his hand.
