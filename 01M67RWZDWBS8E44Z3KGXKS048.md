---
id: 01M67RWZDWBS8E44Z3KGXKS048
title: "critique und optimize laufen je in einem frischen Agenten, nicht im Agenten, der gebaut hat"
status: optimize
ready: true
creator: berkc
assignee: berkc
goal: "Die Lanes critique und optimize tragen context: fresh in der ausgelieferten Lane und auf diesem Board, so dass Worker und Dispatcher sie je einem neuen Agenten geben - zwei zusaetzliche Augenpaare statt Selbstkontrolle."
context: |-
  Was falsch ist: critique und optimize macht heute derselbe Agent, der in in-progress gebaut hat. Er findet seine eigenen Luecken nicht.
  Ausloeser: Experiment 3CBCN6 (31.10.2026): der Nur-Haiku-Lauf liess eine veraltete Hilfezeile und einen falschen Kommentar stehen; seine eigene critique/optimize fand das nicht. Der Maintainer: 'setz das bitte auch um das 2 agenten starten sollen fuer critique und optimize.'
  Bekannt: review hat schon context: fresh (core/lane/lane.go:60-76, ContextFresh). jaira-worker gibt eine fresh-Lane an einen neuen Agenten (core/role/builtin/jaira-worker/SKILL.md).
  Offen fuer den Plan: kostet zwei Agentenstarts mehr pro Ticket; beide Lanes sind tier strong.
definition-of-done: "core/lane/builtin critique und optimize und .jaira/lanes/critique.md, optimize.md tragen context: fresh; jaira show <id> --for-lane critique --json zeigt das"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-31T01:24:02Z
updated-at: 2026-10-31T02:54:10Z
updated-by: berkc
claimed-by: berk-ca-3153991
claimed-at: 2026-10-31T02:32:33Z
outcome-what: "docs/how-to-work-with-jaira.html:387 trennt jetzt Schluss: shipped-Board = critique, optimize, review; ohne Angabe = letzte Agenten-Lane vor der ersten Menschen-Lane; Commit 065df0d9"
outcome-why: "critique r2 (02:50): beide Klauseln waren im Satz verschmolzen, so dass der Default critique/optimize/review las"
outcome-resolves: "docs/how-to-work-with-jaira.html:387 nennt die drei Lanes fuer das shipped-Board und beschreibt den Default getrennt, wie README.md:840-843"
review-summary: |-
  core/role/builtin/jaira-worker/SKILL.md:16 still says "On the shipped board that is `review`." — change to critique, optimize and review (the change made all three fresh; step 4 and line 157 already say so)
  core/role/builtin/jaira-teamlead/SKILL.md:143 "only a lane marked fresh context — `review` on the shipped board —" — name critique, optimize and review, as README.md:78 now does; append the new teamlead SKILL hash to core/shipped/shipped.sha256
  docs/how-to-work-with-jaira.html:387 (German worker row) "also <code>review</code>" — name critique, optimize und review on the shipped board, keep the no-declaration default sentence as README.md:840 does
  Left alone: lane files, AGENTS.md/CLAUDE.md lane block, impliedFresh/review decision (note 02:35), setUpOldBoard (note 02:38), release note incl. "- !" line for hand-edited lane files, worker step 4 fresh->fresh line, dispatcher unchanged (one worker per lane already) — all sound
  round 2: worker SKILL.md:16 fixed; teamlead SKILL.md:143 fixed + hash appended (core/shipped ok)
  docs/how-to-work-with-jaira.html:387 still wrong: 'ohne Angabe die letzte Agenten-Lane vor der ersten Menschen-Lane, auf dem ausgelieferten Board also critique, optimize und review' says the no-declaration default IS those three. Split it like README.md:840-843: 'Nur eine Lane mit <code>context: fresh</code> – auf dem ausgelieferten Board <code>critique</code>, <code>optimize</code> und <code>review</code>; ohne Angabe die letzte Agenten-Lane vor der ersten Menschen-Lane – bekommt ...'
  round 3 (fresh agent): docs/how-to-work-with-jaira.html:387 now split like README.md:840-843 (shipped board: critique, optimize, review; no-declaration default as its own clause) - fixed. Whole diff re-read: lane files, AGENTS/CLAUDE lane block, README:78/840, worker+teamlead SKILL, shipped hashes, release note, tests consistent; grep finds no stale 'only review' mention; go test core/lane core/shipped ok. Nothing left to say.
round: 3
review-gaps: "none - considered: stripLaneContexts exists twice (core/lane/context_test.go:51, internal/cli/context_test.go:177); left, test helpers cannot cross packages and the repo has no shared testutil package, exporting it from core/lane would put test code in the binary. setUpOldBoard (7 callers) kept. No dead code, no fluff, no hot path. go test core/lane core/shipped, internal/cli -run Context ok."
---

# critique und optimize laufen je in einem frischen Agenten, nicht im Agenten, der gebaut hat

## Definition of Done

- [x] core/lane/builtin critique und optimize und .jaira/lanes/critique.md, optimize.md tragen context: fresh; jaira show <id> --for-lane critique --json zeigt das
  proof: core/lane/builtin/25-critique.md:9, 26-optimize.md:9, 40-review.md:9 (context: fresh); jaira show XKS048 --for-lane critique|optimize|review --json shows context fresh
- [x] jaira update --lanes bringt es auf ein bestehendes Board (oder die Release-Notiz sagt, was zu tun ist)
  proof: .jaira/lanes/critique.md:9, optimize.md:9, review.md:9 changed by 'jaira update --lanes' (updated); AGENTS.md/CLAUDE.md lane block only; core/lane TestRefreshShipped* ok
- [x] Release-Notiz unter core/release/notes/
  proof: core/release/notes/XKS048.md (one '- ' line, one '- ! ' line)

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Lane files: add the line 'context: fresh' directly after 'model-tier: strong' in core/lane/builtin/25-critique.md, 26-optimize.md AND 40-review.md (review must become explicit: impliedFresh in core/lane/lane.go:648 switches the implied review default off as soon as any lane declares a valid context). Change nothing else in them.
- [x] core/shipped/shipped.sha256: keep every old hash, append the new sha256 of each changed file (lanes + SKILL.md below) under its own path; run 'go test ./core/shipped -count=1' and follow its message.
- [x] core/lane/context_test.go: setUpBoard stays as it was (the shipped board as installed — refresh_test.go boardWithTesting relies on that). Add setUpOldBoard(t) = setUpBoard(t) + stripLaneContexts(t, root) (strips every 'context: fresh' line from ProjectLanesDir(root)/*.md; the helper already in the worktree is fine) and switch all 7 setUpBoard calls in context_test.go to setUpOldBoard. TestContextDefaultsToReviewOnTheShippedBoard: shipped-offer half wants 'critique,optimize,review' (no '*'), old-board half keeps 'review*'; rename TestContextShippedLanesAreFreshAndAnOldBoardDefaultsToReview, fix its comment. refresh_test.go untouched. Run 'go test ./core/lane -count=1': all pass.
- [x] internal/cli/context_test.go: add helper stripLaneContexts(t, dir) (same stripping, lane.ProjectLanesDir(dir)); in TestLanesAndForLaneReportContext, TestAnInvalidContextIsReportedNotDropped and TestAnInvalidContextLeavesTheDefaultFresh call runLanes(t, dir) once so the files exist, then stripLaneContexts, before the rest. Run 'go test ./internal/cli -run Context -count=1' and 'go test ./core/board -count=1'.
- [x] core/role/builtin/jaira-worker/SKILL.md 'The fresh lane' step 4: add that when the fresh agent moved the ticket into another fresh lane (critique -> optimize), start a NEW agent on that lane (step 2) — never the previous one, never yourself. 'Judging your own work' (line ~157): 'In critique, optimize and testing' -> 'In a lane not marked fresh (testing on the shipped board)'. Dispatcher needs nothing (it gives every lane its own worker).
- [x] README.md:78 and :840-842: the fresh lanes on the shipped board are critique, optimize and review, not only review (keep the 'board where no lane says' sentence for older boards).
- [x] Build: export PATH=$HOME/.local/go/bin:$PATH; go build -o /tmp/claude-1000/jaira-XKS048 ./cmd/jaira. In this worktree run '/tmp/claude-1000/jaira-XKS048 update --lanes': critique.md, optimize.md, review.md must report updated; then '/tmp/claude-1000/jaira-XKS048 show XKS048 --for-lane critique --json' shows context fresh; git diff CLAUDE.md AGENTS.md shows critique/optimize marked **fresh context** (if update rewrote anything unrelated in those files, keep only the lane-block lines). Commit .jaira/lanes + agent files with the code.
- [x] core/release/notes/XKS048.md, one unwrapped line each: '- critique and optimize now run in a new agent each (context: fresh), and review declares context: fresh in its file; run jaira update --lanes to get them.' and '- ! A lane file update --lanes skipped (you edited it) gets no context line: add context: fresh to critique.md, optimize.md and review.md by hand — once any lane declares context, review is no longer fresh by default.'
- [x] Commit feat(XKS048): ... with the ticket file; jaira dod the 3 DoD items with proofs.

## Progress
- **2026-10-31 02:35 · berkc** — Plan decision: review/40-review.md gets context: fresh explicitly too. Reason: impliedFresh (core/lane/lane.go:648) gives review its fresh only while NO lane declares a valid context; no shipped lane declares one today, so adding fresh to critique alone would silently make review non-fresh. Rejected alternative: change impliedFresh to add the default even when fresh is declared elsewhere — that reverses TestContextDeclaredReplacesTheDefault and the 0.3.7 release note. Consequence: the implied-default tests (core/lane/context_test.go setUpBoard, internal/cli/context_test.go 3 tests) must strip context lines to simulate an old board; impliedFresh itself needs no code change and handles several fresh lanes fine (it returns early). Residual risk: a board whose review.md was hand-edited is skipped by update --lanes and loses fresh review once critique is refreshed -> '- !' release note line. This worktree's .jaira/lanes critique/optimize/review are byte-identical to the builtins (known hashes), so update --lanes refreshes them here; the master tree's critique.md may differ. Worker SKILL step 4 only handled fresh->person or fresh->back; fresh critique moving into fresh optimize needed a line. Dispatcher already uses one worker per lane: no change.
- **2026-10-31 02:37 · berkc** — Stopped at plan step 3. Plan step 3 makes setUpBoard strip 'context: fresh' from the lane files (simulating an old board). That breaks core/lane/refresh_test.go TestRefreshShippedUpdatesAnOlderShippedLane, which the plan did not list: boardWithTesting uses setUpBoard, and after the strip review.md equals the old shipped bytes (known in shipped.sha256), so RefreshShipped reports review as 'updated' where the test expects every lane except testing 'unchanged'. Only that one test fails (go test ./core/lane). Options for the planner: (a) exempt review in that test's loop, (b) strip only in the context tests (not in setUpBoard), (c) have boardWithTesting skip the strip. Not decided here. Done so far: plan 1 (three lane files get 'context: fresh'), plan 2 partly (25-critique, 26-optimize, 40-review hashes appended; go test ./core/shipped passes; SKILL.md hash still to append). Plan 3 edits in core/lane/context_test.go are in the worktree, uncommitted.
- **2026-10-31 02:38 · berkc** — Planner on the 02:37 stop: option (b). setUpBoard keeps meaning 'the shipped board as installed' — refresh_test.go (boardWithTesting) tests RefreshShipped against exactly that, and a stripped review.md is a genuinely older shipped version, so 'updated' would be correct behaviour, not a bug; exempting review in that loop (a) would weaken the test. (c) would add a flag to a helper. Instead context_test.go gets its own setUpOldBoard (setUpBoard + strip), used by all 7 context tests, which are about the implied default of a pre-field board. Plan step 3 reworded accordingly.
- **2026-10-31 02:39 · berkc** — Planner amendment applied: setUpBoard back to installed shipped board; setUpOldBoard (setUpBoard + stripLaneContexts) used by all 7 context_test.go calls. go test ./core/lane -count=1: ok.
- **2026-10-31 02:40 · berkc** — Step 4 done: internal/cli stripLaneContexts helper added; TestLanesAndForLaneReportContext, TestAnInvalidContextIsReportedNotDropped, TestAnInvalidContextLeavesTheDefaultFresh run runLanes once, then strip. go test ./internal/cli -run Context and ./core/board: ok.
- **2026-10-31 02:45 · berkc** — critique (fresh): three stale 'only review is fresh' statements the diff missed — worker SKILL.md:16, teamlead SKILL.md:143 (+ shipped hash), docs/how-to-work-with-jaira.html:387. Everything else stands (see review-summary). Fix those and return to critique.
- **2026-10-31 02:50 · berkc** — critique r2: worker:16 and teamlead:143(+hash) fixed. docs/how-to-work-with-jaira.html:387 fix merged the two clauses so the German text now says the no-declaration default is critique/optimize/review ('also'); separate them as README.md:840-843 does: shipped board = critique, optimize, review; no declaration = last agentic lane before the first human lane. Only remaining finding.
