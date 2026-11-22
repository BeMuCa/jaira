---
id: 01M7ZYRHSGN0G8XHJPT5R2X4TJ
title: Kleine Review-Funde werden vor signoff behoben; ein kleiner Fix bekommt nur Delta-Review
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Ticket kommt nach signoff ohne behebbare Kleinigkeiten, und ein Fix von wenigen Zeilen kostet keinen vollen Lane-Durchlauf mehr."
context: |-
  Was heute teuer war (21.11.): Review schrieb "kein Blocker" plus behebbare Notes; der Teamlead gab sie dann aus signoff zurück (9 Mal) — jedes Mal ein kompletter Durchlauf critique/optimize/testing/bughunt/review.
  Seit 0.5.15 schicken testing und review nur bei Blockern zurück, kleine Ränder "fahren mit". Ergebnis: die Ränder kommen in signoff an und gehen dann trotzdem zurück.
  Gewünscht, drei Teile (Lane-Prompts in core/lane/builtin und Worker-Skill core/role/builtin/jaira-worker):
  1. Review markiert jede Note als fix-now (behebbar in wenigen Zeilen, kein neues Verhalten) oder note; der Worker behebt fix-now noch vor signoff, ohne neue Regel einzuführen (Alex' Regel: Ränder wachsen → kleineres Design).
  2. Eine Fix-Runde mit kleinem Diff: stop-after-fix setzen, Tests nur für Geändertes, Review nur über das Delta seit der letzten Review.
  3. Diff unter ~50 Zeilen: critique und optimize in einem Durchgang.
  Vorsicht: heute fand der Bughunt nach "kleinen" Fixes zweimal einen Blocker (DBJHNE, AX8K02) — Grenze für "klein" konservativ, Bughunt bei Änderung an Locking oder Parsing nie überspringen.
  Start erst nach 3DVVFZ (dieselben Lane-Prompts).
definition-of-done: Review-Prompt verlangt fix-now/note pro Fund; Worker-Skill behebt fix-now vor signoff
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-11-21T21:03:57Z
updated-at: 2026-11-22T15:00:52Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3877259
claimed-at: 2026-11-22T14:47:51Z
outcome-what: "Review marks notes fix-now or note and sends fix-now back once; small fix rounds run lean (stop-after-fix, critique+optimize in one pass, delta review from 'review: done at <commit>'); testing still hunts when a lock or parser changed"
outcome-why: Small fixable notes reached signoff and each give-back cost a full lane run (9 times on 2026-11-21)
outcome-resolves: "Tickets arrive at signoff without fixable edges, and a few-line fix no longer costs a full pass"
review-summary: "Full review (first review, no earlier 'review: done at'). Prompt-only change, four lane/role files plus release note and shipped hashes. 40-review.md: every note is marked 'note: fix-now:' (few lines, no new behaviour) or plain 'note:'; with fix-now and no blocker the verdict stays accept but the reviewer moves the ticket back to in-progress itself with a 'review: fix-now:' note; every review ends with 'review: done at <HEAD>'; a later review whose delta since that commit is an ancestor, under 50 changed lines (git diff --shortstat, .jaira excluded), has no untracked file and touches no lock/parser reads only that delta and writes plain notes. 25-critique.md: under 50 changed lines (whole diff, or the fix since 'review: done at') critique runs optimize's passes itself and moves through optimize. 27-testing.md: stop-after-fix skips the hunt except when the fix touches locking or a parser; fix-round test selection also counts untracked files. jaira-worker SKILL step 5: close exactly the fix-now notes, set stop-after-fix for a small lean round, clear it after a blocker. All four prompt hashes are in shipped.sha256; go test ./core/shipped ./core/lane/... ./core/release/... pass."
round: 4
back-from: signoff
back-by: person
back-count: 1
review-gaps: |-
  note: fix-now: 40-review.md says the fix-now round 'happens at most once', but that only holds when the next review is a delta review - a fix round that touches a lock/parser or reaches 50 lines gets a full review, which may mark fix-now again, and nothing bounds the repeat; fix: let a review mark fix-now only when the ticket carries no earlier 'review: fix-now:' note
  note: fix-now: cockpit followTitle (core/mod/builtin/jaira-cockpit/hooks/board.ts:780) strips 'note:' and 'pre-existing:' but not 'fix-now:', so a fix-now line that does reach signoff becomes a follow-up titled 'fix-now: ...'; fix: strip it in the same regex, plus one test in board.test.ts
  note: the pathspec ':(exclude).jaira' with '-- .' in 25-critique.md/40-review.md is relative to the cwd; from a subdirectory it counts only that subtree - ':/' and ':(exclude,top).jaira' as the lane-role skill already uses would be cwd-proof
test-verdict: "pass: round 2 - shipped + lane tests green (rc=0), shipped-history --check ok; untracked-file rule checked on a scratch repo (new.go listed, .jaira/ excluded); DoD 1-3 still hold"
hunt-verdict: "pass: bughunt done at b6777c680, 0 change, 0 pre-existing"
review-verdict: "accept: all three DoD items are met by the diff and the hashes ship; no blocker - the at-most-once fix-now bound has a gap on lock/parser rounds (note), which costs an extra round, not a wrong result"
review-round: 3
review-verdict-gaps: |-
  note: fix-now: 40-review.md says the fix-now round 'happens at most once', but that only holds when the next review is a delta review - a fix round that touches a lock/parser or reaches 50 lines gets a full review, which may mark fix-now again, and nothing bounds the repeat; fix: let a review mark fix-now only when the ticket carries no earlier 'review: fix-now:' note
  note: fix-now: cockpit followTitle (core/mod/builtin/jaira-cockpit/hooks/board.ts:780) strips 'note:' and 'pre-existing:' but not 'fix-now:', so a fix-now line that does reach signoff becomes a follow-up titled 'fix-now: ...'; fix: strip it in the same regex, plus one test in board.test.ts
  note: the pathspec ':(exclude).jaira' with '-- .' in 25-critique.md/40-review.md is relative to the cwd; from a subdirectory it counts only that subtree - ':/' and ':(exclude,top).jaira' as the lane-role skill already uses would be cwd-proof
review-check: "1. git diff 180daf772..HEAD -- core/lane/builtin/40-review.md  2. see the 'Mark every note ... fix-now' paragraph and, at the end, 'No blocker but fix-now notes: the verdict stays accept, and the ticket goes back'  3. git diff 180daf772..HEAD -- core/lane/builtin/25-critique.md core/lane/builtin/27-testing.md core/role/builtin/jaira-worker/SKILL.md  4. see 'A small change is critiqued and optimized in one pass', the lock/parser exception in testing, and worker step 5  5. go test ./core/shipped/ ./core/lane/...  -> ok (hashes of the new prompts are listed)  6. live proof comes on the next ticket after 'jaira update --lanes': a review with only fix-now notes should send the ticket back to in-progress, not to signoff"
---

# Kleine Review-Funde werden vor signoff behoben; ein kleiner Fix bekommt nur Delta-Review

## Definition of Done

- [x] Review-Prompt verlangt fix-now/note pro Fund; Worker-Skill behebt fix-now vor signoff
  proof: core/lane/builtin/40-review.md fix-now paragraph + end-of-prompt send-back; core/role/builtin/jaira-worker/SKILL.md 'The fresh lane' step 5
- [x] Fix-Runde mit kleinem Diff nutzt stop-after-fix und Delta-Review; kleine Diffs haben critique+optimize in einem Durchgang
  proof: core/role/builtin/jaira-worker/SKILL.md fresh lane step 5 (stop-after-fix); 40-review.md 'A delta review'; 25-critique.md 'A small change is critiqued and optimized in one pass'; 27-testing.md lock/parser exception
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/R2X4TJ.md; core/shipped/shipped.sha256 (scripts/shipped-history.sh --write)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 14:47 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 14:49 · Alexander Sacharov** — decisions (settled, none open for the person): 1) fix-now is a sub-mark of note ('note: fix-now: ...') so cockpit/signoff parsers keep working; 2) review with fix-now lines but no blocker moves the ticket back to in-progress itself (verdict stays accept) - works for worker and dispatcher alike and the ticket never touches signoff; 3) only a full review may write fix-now - a delta review writes plain notes, so the fix-now round happens at most once (edges grow -> smaller design); 4) delta base = the review's own closing note 'review: done at <commit>', like bughunt; 5) 'small' = under 50 changed lines outside .jaira/ AND no lock/parser touched; 6) critique+optimize in one pass lives in the critique prompt (the critique agent moves through optimize itself), so it holds for per-lane dispatch too; 7) the worker sets stop-after-fix for a small fix-now round and clears it again if a later round is not small
- **2026-11-22 14:51 · Alexander Sacharov** — critique: 25-critique.md one-pass rule must measure the fix since 'review: done at <commit>' when there is one, not the whole payload diff; 40-review.md use --shortstat
- **2026-11-22 14:51 · Alexander Sacharov** — critique pass 2: both findings answered in fe7428b63 (one-pass rule measures from review: done at; --shortstat)
- **2026-11-22 14:52 · Alexander Sacharov** — waiting: bughunt subagent since 15:52
- **2026-11-22 14:53 · Alexander Sacharov** — bughunt: change: blocker: core/lane/builtin/40-review.md:134-136 and core/lane/builtin/25-critique.md:73-74 - class: the fix-round size is measured with 'git diff --shortstat <commit> -- . ':(exclude).jaira'', which never counts untracked files, so a fix that adds a new file not yet committed is counted as almost nothing. The delta review then 'reads only that diff' and never sees the new file at all. The same blind spot sits in the older line core/lane/builtin/27-testing.md:147 ('git diff --name-only <commit>' for a fix round), which predates this ticket. Example: a fix round that edits one line and adds a 500-line new.go that is not committed yet (the worker skill says review judges uncommitted work) is counted as '1 file changed, 1 insertion', so it gets a delta review and a one-pass critique, and new.go reaches signoff unreviewed. Fix the class: count untracked files too (e.g. add 'git ls-files --others --exclude-standard' to the count and to the files the delta review reads), or require the fix to be committed before the size is measured; ran: scratch repo, seq 500 > new.go untracked + 1-line edit, git diff --shortstat HEAD -- . ':(exclude).jaira' prints '1 file changed, 1 insertion(+)' (rc=0); blame: fe7428b63 (R2X4TJ)
- **2026-11-22 14:53 · Alexander Sacharov** — bughunt: done at fe7428b63 - 1 finding (1 change, 0 pre-existing) in 6 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/lane/builtin/25-critique.md core/lane/builtin/27-testing.md core/lane/builtin/40-review.md core/release/notes/R2X4TJ.md core/role/builtin/jaira-worker/SKILL.md core/shipped/shipped.sha256
- **2026-11-22 14:54 · Alexander Sacharov** — fixed bughunt finding by class: review/critique size rule now also requires no untracked file outside .jaira/; the sibling in 27-testing.md fix-round test selection now adds git ls-files --others; checked on a scratch repo: untracked new.go listed, .jaira/ file excluded
- **2026-11-22 14:54 · Alexander Sacharov** — critique pass 3: the untracked-file finding is answered in all three prompts; nothing else reopened
- **2026-11-22 14:54 · Alexander Sacharov** — waiting: bughunt round 2 subagent since 15:54
- **2026-11-22 14:55 · Alexander Sacharov** — bughunt: done at b6777c680 - 0 findings in 6 modules, follow-up hunt since fe7428b63 (diff of 25-critique.md, 27-testing.md, 40-review.md, shipped.sha256; R2X4TJ.md and jaira-worker/SKILL.md unchanged); concurrency, error paths, input, permissions checked; ran: sha256 of every current and earlier R2X4TJ version found in shipped.sha256, go test ./core/shipped/ ./core/lane/... ok; files: core/lane/builtin/25-critique.md core/lane/builtin/27-testing.md core/lane/builtin/40-review.md core/release/notes/R2X4TJ.md core/role/builtin/jaira-worker/SKILL.md core/shipped/shipped.sha256
- **2026-11-22 14:55 · Alexander Sacharov** — waiting: fresh review agent (run-lane.sh tab) since 15:55
- **2026-11-22 14:57 · Alexander Sacharov** — review: done at b6777c680
- **2026-11-22 14:58 · Alexander Sacharov** — worker: review left 2 fix-now notes (installed board lanes still run the old review prompt, so it moved to signoff); pulling back to close them before the person sees it - dogfooding this ticket's rule. Plain note (cwd-relative pathspec) stays for the person. Fix is small, but touches the cockpit parser regex -> no stop-after-fix, hunt runs
- **2026-11-22 15:00 · Alexander Sacharov** — closed review's 2 fix-now: 40-review.md now marks fix-now only while no earlier 'review: fix-now:' note exists (bounds it to one round for full reviews too); cockpit followTitle strips 'fix-now:' + test in board.test.ts (claude plugin test: 303 pass). Left the plain note (cwd-relative pathspec) for the person
