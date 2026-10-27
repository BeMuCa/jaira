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
updated-at: 2026-10-27T00:48:55Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: 4n1b66-driver
claimed-at: 2026-10-27T00:01:46Z
outcome-what: "The testing lane (core/lane/builtin/27-testing.md and the board copy .jaira/lanes/testing.md, byte-identical prompt) now says in Gates: build, format/vet, and test only the touched packages (in Go go test ./<pkg> -count=1); never run the full suite in a lane, it runs once before the pull request by the main session (in Go go test ./... -race). Pass 2 says a DoD item that asks for the full suite is left unticked and named in test-verdict as left for that run, not a finding. jaira-worker and jaira-dispatcher SKILL.md carry the rule in one line each. core/release/notes/4N1B66.md tells readers to run jaira lanes use testing --force if their board has a testing lane, jaira roles install --force with their install flag, and to fix a catalogue testing.md first."
outcome-why: "The lane demanded the full race suite in every testing pass, against the maintainer's rule that it runs once before the PR (it takes 7-8 minutes here); the rule lived only in session memory, so a worker followed the prompt. Many tickets carry a DoD item 'go test ./... -race gruen', which the new Gates rule would otherwise turn into an endless send-back. Shipped lanes and roles go into non-Go repositories, so Go is named as one case."
outcome-resolves: "DoD 1: 'run the whole suite' is gone; .jaira/lanes/testing.md:22 and core/lane/builtin/27-testing.md:23 say touched-package tests in the lane and the full -race suite once before the PR by the main session, prompts identical. DoD 2: one rule each in core/role/builtin/jaira-worker/SKILL.md:168 and core/role/builtin/jaira-dispatcher/SKILL.md:268. DoD 3: core/release/notes/4N1B66.md:1."
review-summary: "The testing lane prompt no longer says 'Build everything and run the whole suite ... with the race detector', in both copies: the shipped core/lane/builtin/27-testing.md and this board's .jaira/lanes/testing.md, whose # Prompt sections are byte-identical. Gates now says: build, run the format and vet checks, and run only the tests of the packages the change touched (in Go: go test ./<pkg> -count=1); never run the full suite in a lane - it runs once, before the pull request, by the main session (in Go: go test ./... -race). Pass 2 gets one sentence: a DoD item that asks for the full suite stays unticked, is named in test-verdict as left for that run, and is not a finding - so the common DoD item 'go test ./... -race gruen' no longer causes an endless send-back. The example verdict reads 'touched packages green'. The same rule is in core/role/builtin/jaira-worker/SKILL.md:168 and core/role/builtin/jaira-dispatcher/SKILL.md:268. core/release/notes/4N1B66.md is one '- ' line: run jaira lanes use testing --force if the board has a testing lane, jaira roles install --force with the flag you installed with, and update or delete a catalogue testing.md first. The three fixes the first review asked for are in, and each Go command sits on one line inside its backticks."
round: 3
review-gaps: "pass 2 (f7deb7fa only): trimmed the new clause in both testing lanes (.jaira/lanes/testing.md:31-33, core/lane/builtin/27-testing.md:32-34, still byte-identical) from 'An item that asks for the full suite is not run in the lane: leave it unticked...' to 'Leave an item that asks for the full suite unticked...' - 'is not run in the lane' repeated Gates' 'Never run the full suite in a lane' three lines above, and 'left for the run before the pull request' still says where it runs; same instruction. Left: 'it is not a finding' (the exception to the sentence before it); the Gates and jaira-worker/SKILL.md changes (re-wraps only, no new words); in the release note, 'role files' three times (clearer than pronouns), 'with the flag you installed with' (plain, only slightly clumsy), and the bare --force next to a backticked one (formatting, not a repeat)."
test-verdict: "pass: touched packages green (gofmt/vet/build/core/lane+role+release RC=0), DoD 1-3 verified in the tree (prompts identical, 'whole suite' gone, one-line note), release-note steps exercised on a scratch board; DoD asks no full-suite item"
review-verdict: "Send back. DoD 1-3 are met as literally written; gofmt -l empty, go vet ./..., go build ./... and go test ./core/lane ./core/role ./core/release -count=1 pass; commits use the noreply address, no trailers, no denylisted strings. But the change introduces a contradiction in the testing lane for a common DoD shape (gap 1), and both instructions in the release note fail to do what they say (gaps 2 and 3); all three reproduced, each fix a few words. Gaps 4 and 5 are questions for the maintainer, not reasons to send back."
review-round: 2
review-check: "Pending: this review sent the ticket back; the next review writes the hand-check for the fixed change."
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
- **2026-10-27 00:29 · BeMuCa** — review: sent back with three fixes, verified by the driver. (1) Both testing lanes, same bytes: pass 2 says verify every DoD item, Gates forbids the full suite, and 16 open tickets carry the DoD item 'go test ./... -race gruen' (jaira-role-brainstorm/SKILL.md:54-56 teaches it). Add a clause: such an item is not run in the lane; leave it unticked, name it in test-verdict as left for the run before the pull request, not a finding. (2) core/release/notes/4N1B66.md: 'jaira lanes use testing --force' adds the lane to a board that never had one (jaira init writes no testing lane) - say 'if your board has a testing lane'. (3) Same note: bare 'jaira roles install' exits 2 (needs --project|--global|--into) and without --force an older installed role is skipped as edited - name the flag the reader installed with plus --force, one '- ' line. Nit to fix with it: keep each Go command on one line inside its backticks (lanes Gates, worker skill). Gaps 4 (dispatcher route / optimize 'run the tests') and 5 (who runs the suite before the PR) are for the maintainer, not this ticket.
- **2026-10-27 00:35 · BeMuCa** — critique (pass 3, read only f7deb7fa against the 00:29 review note): none. Fix 1 - pass 2 in both testing lanes carries the full-suite clause, prompts byte-identical (diff of the # Prompt sections empty), worded for any repository. Fix 2 - release note says 'if your board has a testing lane'. Fix 3 - names --force plus the install flag, still one '- ' line (core/release/notes/4N1B66.md:1). Nit - no Go command breaks inside backticks in either lane or in the worker skill. Seen and let stand, not DoD-breaking: the note's '(--project or --global)' leaves out --into (internal/cli/roles.go:203), but 'the flag you installed with' still covers it; who ticks a held-back full-suite DoD item after the pre-PR run belongs to review gap 5, closed by the driver note at 00:10.
