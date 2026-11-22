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
updated-at: 2026-11-22T14:54:12Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3877259
claimed-at: 2026-11-22T14:47:51Z
outcome-what: "Review marks notes fix-now or note and sends fix-now back once; small fix rounds run lean (stop-after-fix, critique+optimize in one pass, delta review from 'review: done at <commit>'); testing still hunts when a lock or parser changed"
outcome-why: Small fixable notes reached signoff and each give-back cost a full lane run (9 times on 2026-11-21)
outcome-resolves: "Tickets arrive at signoff without fixable edges, and a few-line fix no longer costs a full pass"
review-summary: none
round: 3
back-from: testing
back-by: lane
back-count: 1
review-gaps: "none: the 50-line / lock-or-parser rule stands in review, critique, testing and the worker skill on purpose - each lane prompt is read alone by its own agent; nothing orphaned"
test-verdict: "pass: go build + go test ./core/... green (rc=0), shipped-history --check ok; DoD 1-3 checked in the tree; fresh board from the built binary gets review.md with fix-now; prompt-only change writes nothing, so no two-writer case"
hunt-verdict: "fail: bughunt done at fe7428b63, 1 change - fix-round size via git diff --shortstat misses untracked files"
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
