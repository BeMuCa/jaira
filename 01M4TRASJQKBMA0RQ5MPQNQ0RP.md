---
id: 01M4TRASJQKBMA0RQ5MPQNQ0RP
title: "Der Agent-Block sagt, dass ein Agent ein Ticket aus human weiterschiebt, sobald die Frage beantwortet ist"
status: review
ready: true
creator: BeMuCa
goal: "Was der erzeugte Agent-Block ueber die human-Lane sagt, stimmt mit ihrem Gate ueberein."
context: |-
  Der erzeugte Agent-Block (CLAUDE.md, AGENTS.md) nennt human 'a person's, not yours — you may move work in, never out' (core/board/announce.go:210-211).
  Das Gate sagt etwas anderes: human hat nur requires-question, kein requires-human-exit (core/lane/builtin/30-human.md:8). Ein Agent darf human also verlassen.
  Derselbe Block sagt weiter unten 'until it sits in a human lane, then continue once the human has answered' - also doch raus. Der Block widerspricht sich selbst.
  Ursache: announce.go:209 legt f.HumanExit und f.Question in einen gemeinsamen case.
  Folge am 2026-10-13: Nach Berks Antwort auf PNC854 wurde der erste Move aus human abgelehnt (von Claude Codes Klassifizierer; dass dieser Satz der Grund war, ist eine Vermutung).
  Verwandt: 64GNZT (Namen von human und signoff).
definition-of-done: "Fuer eine Lane mit nur requires-question sagt der Block: eine Person beantwortet die Frage, danach schiebt der Agent weiter; fuer requires-human-exit bleibt 'never out'; ein Test in core/board deckt beide Faelle"
tags:
  - gates
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-10-13T13:48:17Z
updated-at: 2026-10-22T08:14:18Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-4189248
claimed-at: 2026-10-22T07:55:22Z
outcome-what: "The generated agent block now says a question-only lane (human) is a person's to answer, after which the agent moves the ticket on itself; a requires-human-exit lane (signoff) keeps 'never on'. The static line says 'a lane that waits for a person's acceptance' instead of 'the human review lane'. CLAUDE.md and AGENTS.md are regenerated."
outcome-why: "The block said an agent may never move work out of human, although that lane has no human-exit gate and the same block says to continue once the human has answered; agents refused to move on after a person's answer."
outcome-resolves: "The block now matches the human lane's gate and its own later paragraph; TestLaneSectionQuestionLaneIsLeftByTheAgentAcceptanceLaneIsNot covers the question-only and the requires-human-exit lane."
review-summary: |-
  critique 1: none. The human mark still says how work enters (the appended 'only with an open question (--question)'); 'move the ticket on yourself' matches the decision and the later 'continue once the human has answered'. The static line 'a lane that waits for a person's acceptance' is accurate: on a stock board only signoff has requires-human-exit (core/lane/builtin/45-signoff.md:8). The new test overlaps lanesection_test.go:46 and :173 and stays because the DoD asks for one test over both cases. Release note is one instruction line.
  review 1: Two text changes in core/board/announce.go, both rendered into the jaira block of CLAUDE.md and AGENTS.md. (1) A lane with only requires-question (stock board: `human`) used to share signoff's wording 'a person's, not yours — you may move work in, never out'; it now reads '**a person's to answer** — once they have answered, you move the ticket on yourself' (announce.go:249), still followed by 'only with an open question (`--question`); without one, work passes it by'. A requires-human-exit lane (stock board: `signoff`) keeps 'a person's, not yours — you may move work in, or send it back to an earlier lane, never on' (announce.go:247, unchanged, still checked first). (2) The static line 'An agent cannot move work on out of the human review lane' now says 'out of a lane that waits for a person's acceptance' (announce.go:149), so it no longer reads as lane id `human`. Both now match the gate - only requires-human-exit refuses a forward move (core/gate/gate.go:395), and human has only requires-question (core/lane/builtin/30-human.md:8) - and the block's own later paragraph 'continue once the human has answered'. New test TestLaneSectionQuestionLaneIsLeftByTheAgentAcceptanceLaneIsNot (core/board/lanesection_test.go:197) pins both lane lines; the older test at :173 carries the new human text. Repo CLAUDE.md and AGENTS.md regenerated; release note core/release/notes/QNQ0RP.md is one '- ' line.
review-gaps: |-
  optimize 1: simplified TestLaneSectionQuestionLaneIsLeftByTheAgentAcceptanceLaneIsNot (core/board/lanesection_test.go:197) - the line-split loop, its two once-used vars and the 'never out' negative check gave way to two Contains checks anchored on the line prefix, the style of the neighbour test at :173; the anchored check already rules out the old mark (mutation-checked: old human text and a cut signoff mark both fail it). Both cases stay in one test, as the DoD asks. Left: overlap with :46 and :173 (DoD wants this test; critique let it stand); docs/AGENTS.md:218 and the role skills that say no agent leaves human are hand-written, outside this ticket's generated block, already noted by critique as a follow-up. No duplication, dead code or cost in announce.go - two string literals changed, switch order unchanged.
  review 1: No defect found in the diff. Checked: gofmt -l empty, go vet and go build green; go test -race -count=1 ./core/board ok. Mutations: the old human text put back fails the new test at lanesection_test.go:203; merging the HumanExit and Question cases fails it at :206. The branch binary's `jaira update` on a remote-less copy carrying master's CLAUDE.md and AGENTS.md changes exactly two lines in each, and the result is byte-identical to the committed files. On that copy a human->review dry-run is allowed without --force, and signoff->done is refused with 'is a human checkpoint'. Open, not a defect of the diff: DoD 2's go test ./... -race (maintainer, before the PR); the tui suite was not run. Minor: the new test repeats what lanesection_test.go:46 and :173 already pin; it stays because the DoD asks for one test over both cases. Outside this diff: the shipped role skills still say no agent leaves human (core/role/builtin/jaira-worker/SKILL.md:291, jaira-dispatcher/SKILL.md:429, jaira-teamlead/SKILL.md:207, jaira-role-acceptance/SKILL.md:200) - confirmed, already a recorded follow-up. Also pre-existing and out of scope: the gate does not check that the question was answered (the human->review dry-run passed while the question was still open), so 'once they have answered' is an instruction only.
test-verdict: "testing 1: pass: gofmt/vet/build clean, go test -race ./core/board ok (TestLaneSection* incl. new QuestionLane test pass); release note is one '- ' line; CLAUDE.md/AGENTS.md byte-identical to branch binary's update output on a remote-less scratch repo; human line says person answers then agent moves on, signoff keeps 'never on'; scratch gate: human->review dry-run allowed w/o --force, signoff->done dry-run refused (human checkpoint). Full race suite (DoD 2) left to maintainer before PR, so DoD 2 not fully tickable."
---

# Der Agent-Block sagt, dass ein Agent ein Ticket aus human weiterschiebt, sobald die Frage beantwortet ist

## Definition of Done

- [x] Fuer eine Lane mit nur requires-question sagt der Block: eine Person beantwortet die Frage, danach schiebt der Agent weiter; fuer requires-human-exit bleibt 'never out'; ein Test in core/board deckt beide Faelle
  proof: core/board/announce.go:249; TestLaneSectionQuestionLaneIsLeftByTheAgentAcceptanceLaneIsNot (core/board/lanesection_test.go:197)
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/; der CLAUDE.md-Block dieses Repos ist neu erzeugt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Test first in core/board/announce_test.go: a question-only lane's line says a person answers and the agent moves on after the answer (no 'never out'); a requires-human-exit lane keeps 'you may move work in, or send it back to an earlier lane, never on'
- [x] Change the f.Question case text in core/board/announce.go laneSection; check the static line 'An agent cannot move work on out of the human review lane' (announce.go:149) for consistency
- [x] Release note core/release/notes/QNQ0RP.md, one '- ' line written as an instruction
- [x] Build the branch binary, run its update in the worktree (HOME/JAIRA_HOME in scratchpad, JAIRA_NO_FETCH=1), keep the regenerated CLAUDE.md/AGENTS.md block

## Progress
- **2026-10-22 07:57 · BeMuCa** — Test put in lanesection_test.go beside the other laneSection tests (announce_test.go has none); the old assertion there ('never out') was updated. Static line 'human review lane' reworded to 'a lane that waits for a person's acceptance' because the static half is board-agnostic and 'human' is also a lane id. Update on the worktree touched only AGENTS.md and CLAUDE.md.
- **2026-10-22 08:03 · BeMuCa** — critique 1: nothing to send back. (1) human mark: dropping 'you may move work in' loses nothing - the appended mark 'only with an open question (`--question`)' still says how work enters; 'move the ticket on yourself' mirrors the decision and the later paragraph 'continue once the human has answered'. (2) static line 'a lane that waits for a person's acceptance' is accurate: on a stock board only signoff has requires-human-exit (core/lane/builtin/45-signoff.md:8, description 'waiting for a person to accept it'); the old 'human review lane' meant signoff by its display name 'Human Review' but reads as lane id human inside the block, so the rewording is right. README.md:820, docs/AGENTS.md:218 and .claude/skills/jaira/SKILL.md:376 say 'human review lane' = signoff and do not contradict. (3) test: TestLaneSectionQuestionLaneIsLeftByTheAgentAcceptanceLaneIsNot overlaps lanesection_test.go:46 and :173, which pin the same marks; let stand - the DoD asks for one test over both cases and its failure message names the rule. (4) release note: one '- ' line, instruction form, same shape as earlier 'Run `jaira update` so ...' notes. Outside this ticket's goal (the generated block), for a follow-up: shipped role skills still say no agent leaves human - core/role/builtin/jaira-worker/SKILL.md:290-291 and jaira-dispatcher/SKILL.md:428-429 'you may deliver into one, never out of it' for (human, signoff); jaira-teamlead/SKILL.md:207 'You never move a ticket out of a human lane'; jaira-role-acceptance/SKILL.md:199-200 'an agent cannot leave a human lane otherwise' (without --force). The 2026-10-12 critique on 4R72WD let worker/dispatcher stand as describing the stop, not the gate; under this decision they now also forbid the move after the answer.
- **2026-10-22 08:09 · BeMuCa** — testing 1 passed on a remote-less scratch copy: human->review dry-run allowed without --force, signoff->done refused. DoD 2 stays open only for the full go test ./... -race, which the maintainer runs once before the PR. Follow-up not filed yet (outside this ticket's generated block): core/role/builtin/jaira-worker/SKILL.md:291, jaira-dispatcher/SKILL.md:429, jaira-teamlead/SKILL.md:207 and jaira-role-acceptance/SKILL.md:200 still say no agent moves a ticket out of human.
