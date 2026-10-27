---
id: 01M5MNX3NGYBHKCQ2YBQ4N1B66
title: Die Testing-Lane verlangt keinen vollen Suite-Lauf mehr; die Suite laeuft einmal vor dem PR
status: review
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
updated-at: 2026-10-27T00:29:40Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: 4n1b66-driver
claimed-at: 2026-10-27T00:01:46Z
outcome-what: "The testing lane (core/lane/builtin/27-testing.md and the board copy .jaira/lanes/testing.md, byte-identical prompt), the jaira-worker and jaira-dispatcher role skills and a release note now say: a lane runs only the tests of the packages its change touched; the full suite runs once, before the pull request, by the main session. The Go commands (go test ./<pkg> -count=1, go test ./... -race) are given as the Go case only. The release note tells readers to run jaira lanes use testing --force and jaira roles install, and to update or delete a catalogue copy of testing.md first."
outcome-why: "The lane demanded the full race suite in every testing pass, which takes 7-8 minutes here, against the maintainer's rule that it runs once before the PR; the rule lived only in session memory. The shipped lane and role files also go into non-Go repositories, so the rule is stated generally with Go as one case and no timing number. A catalogue testing.md lays over the built-in on lanes use --force, so an old catalogue copy would bring the old text back."
outcome-resolves: "DoD 1: 'run the whole suite' is gone; .jaira/lanes/testing.md:22 and core/lane/builtin/27-testing.md:23 say touched-package tests in the lane and the full -race suite once before the PR by the main session. DoD 2: one rule each in core/role/builtin/jaira-worker/SKILL.md:168 and core/role/builtin/jaira-dispatcher/SKILL.md:268. DoD 3: core/release/notes/4N1B66.md:1."
review-summary: "The branch has no Go code changes. It is three docs commits that rewrite the Gates step of the testing lane, in the shipped built-in core/lane/builtin/27-testing.md and in this board's copy .jaira/lanes/testing.md. The two prompt bodies are byte-identical, and go test ./core/lane passes. Before, Gates said: build everything and run the whole suite with the race detector. Now a testing pass builds, runs the format and vet checks, and runs only the tests of the packages the change touched (Go: go test ./<pkg> -count=1). It is told never to run the full suite in a lane, because the main session runs it once before the pull request (Go: go test ./... -race). The example verdict changes from 'suite green' to 'touched packages green'. The same rule appears as one bullet in core/role/builtin/jaira-worker/SKILL.md:168-170 and as one line in core/role/builtin/jaira-dispatcher/SKILL.md:268. A new one-line release note, core/release/notes/4N1B66.md, tells readers to run 'jaira lanes use testing --force' and 'jaira roles install', and to update or delete a catalogue testing.md first. The repository's lanes/ catalogue has no testing.md, so the built-in is the only other copy that DoD 1 calls 'Katalog-Lane'."
round: 2
review-gaps: "(1) Defect in the lane: 16 open tickets in .jaira/tickets/ besides this one have a DoD with 'go test ./... -race gruen', and core/role/builtin/jaira-role-brainstorm/SKILL.md:54-56 teaches that as the model DoD. Pass 2 of the same lane says to verify every DoD item ('a criterion you could not verify is a finding, not a pass'), while Gates now forbids the full suite: a testing worker on such a ticket must break the new rule or fail the ticket back over and over. (2) Release note: 'jaira lanes use testing --force' on an existing board adds the testing lane to a board that never had it (scratch board from jaira init: it wrote .jaira/lanes/testing.md and appended testing to .jaira/lanes/order); testing is not among the lanes jaira init writes. (3) Release note: bare 'jaira roles install' exits 2 ('choose exactly one of --project, --global or --into'); with --global but without --force an installed older jaira-worker/SKILL.md is skipped as 'edited here' (exit 3) and the rule never arrives (core/role/install.go:115-116 compares bytes only). (4) Not in the DoD: on the dispatcher route each lane runs as /jaira-role-lane without a brief, so the dispatcher line at :268 reaches no worker; core/lane/builtin/26-optimize.md:43-44 still says 'After each edit, run the tests'. (5) 'runs once before the PR by the main session' is written nowhere a session would act on it (jaira-role-pr, jaira-teamlead say nothing about tests) - critique raised it, the driver deferred it to the maintainer. (6) Nits: Go commands break across lines inside backticks (both lanes, worker skill); 'vet' is Go vocabulary; the dispatcher rule is one 158-character line; the dispatcher line leaves out the touched-packages half."
test-verdict: "pass: gofmt/vet/build rc=0, go test ./core/lane ./core/role ./core/release rc=0, DoD 1-3 verified in the tree (prompts identical, 'whole suite' gone, rule lines in worker+dispatcher, note is one '- ' line), scratch board: lanes show testing has new text, lanes use testing --force replaced old Gates text, roles install wrote both rule lines"
review-verdict: "Send back. DoD 1-3 are met as literally written; gofmt -l empty, go vet ./..., go build ./... and go test ./core/lane ./core/role ./core/release -count=1 pass; commits use the noreply address, no trailers, no denylisted strings. But the change introduces a contradiction in the testing lane for a common DoD shape (gap 1), and both instructions in the release note fail to do what they say (gaps 2 and 3); all three reproduced, each fix a few words. Gaps 4 and 5 are questions for the maintainer, not reasons to send back."
review-round: 2
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
