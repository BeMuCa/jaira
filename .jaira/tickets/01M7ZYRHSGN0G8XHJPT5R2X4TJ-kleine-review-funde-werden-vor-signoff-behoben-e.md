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
updated-at: 2026-11-22T16:17:29Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-738753
claimed-at: 2026-11-22T16:16:51Z
outcome-what: "Review marks notes fix-now or note and sends fix-now back once; small fix rounds run lean (stop-after-fix, critique+optimize in one pass, delta review from 'review: done at <commit>'); testing still hunts when a lock or parser changed"
outcome-why: Small fixable notes reached signoff and each give-back cost a full lane run (9 times on 2026-11-21)
outcome-resolves: "Tickets arrive at signoff without fixable edges, and a few-line fix no longer costs a full pass"
review-summary: "Full review (the board's installed review prompt is still the old one, so no delta), commits 1d181f428..96f28e718. 40-review.md: every note of this change is marked 'note: fix-now:' (few lines, no new behaviour) or plain 'note:'; fix-now notes send the ticket back to in-progress once with verdict accept, bounded by an earlier 'review: fix-now:' note; every review closes with 'review: done at <commit>', and a later review whose fix since that commit is under 50 lines, has no untracked file and touches no lock/parser reads only that delta. 25-critique.md: a change under 50 lines (or the fix since the last review) gets critique and optimize in one pass. 27-testing.md: stop-after-fix no longer skips the hunt when locking or a parser is touched; fix-round test selection also counts untracked files. jaira-worker SKILL.md step 5: the worker closes exactly the fix-now notes and sets stop-after-fix for a small fix round. Cockpit followTitle strips 'fix-now:' like 'pre-existing:', with a test. Release note and shipped.sha256 updated."
round: 5
back-from: signoff
back-by: person
back-count: 1
review-gaps: |-
  note: core/lane/builtin/25-critique.md:70 - the one-pass rule has neither the lock/parser exclusion nor the merge-base ancestor check that 40-review.md:133-144 and the ticket's decision 5 attach to 'small', so a 26-line cockpit parser fix (this ticket's last round) got one-pass critique; add the same two clauses (a few lines)
  note: core/lane/builtin/25-critique.md:73 and 40-review.md:136,139 - the size pathspec `-- . ':(exclude).jaira'` is cwd-relative, so an agent outside the repo root measures only its subtree; use `-- :/ ':(exclude,top).jaira'` as jaira-role-lane already does
  note: pre-existing: core/mod/builtin/jaira-cockpit/hooks/board.ts:667-668 - PREFIXED/PRE_EXISTING leave the closing `**` of `**note:**`, so followTitle and the pre-existing class miss a mark after it (bughunt round 3, blame 36a787f03 MFYM1A)
test-verdict: "pass: round 3 - claude plugin test (cockpit) 303 pass rc=0, new test checks follow-up title without fix-now:; go test shipped+lane ok, shipped-history --check ok; DoD 1-3 hold"
hunt-verdict: "pass: bughunt done at 96f28e718, 0 change, 1 pre-existing (cockpit '**note:**' emphasis, MFYM1A) left as a note for the person - no satellite ticket"
review-verdict: "accept: all three DoD items are met by the diff; go test ./core/shipped ./core/lane/... ok; three notes, none a blocker. The earlier review's fix-now round already ran (worker note 15:00), so these stay plain notes and the ticket goes to the person"
review-round: 4
review-verdict-gaps: |-
  note: core/lane/builtin/25-critique.md:70 - the one-pass rule has neither the lock/parser exclusion nor the merge-base ancestor check that 40-review.md:133-144 and the ticket's decision 5 attach to 'small', so a 26-line cockpit parser fix (this ticket's last round) got one-pass critique; add the same two clauses (a few lines)
  note: core/lane/builtin/25-critique.md:73 and 40-review.md:136,139 - the size pathspec `-- . ':(exclude).jaira'` is cwd-relative, so an agent outside the repo root measures only its subtree; use `-- :/ ':(exclude,top).jaira'` as jaira-role-lane already does
  note: pre-existing: core/mod/builtin/jaira-cockpit/hooks/board.ts:667-668 - PREFIXED/PRE_EXISTING leave the closing `**` of `**note:**`, so followTitle and the pre-existing class miss a mark after it (bughunt round 3, blame 36a787f03 MFYM1A)
review-check: "1. go test ./core/shipped/ ./core/lane/...  -> both ok  2. cd core/mod/builtin/jaira-cockpit && claude plugin test .  -> 303 pass, incl. \"a follow-up's title drops the review's fix-now mark\"  3. grep -n 'fix-now' core/lane/builtin/40-review.md  -> the marking paragraph (~line 71) and the send-back sentence at the end  4. grep -n 'one pass' core/lane/builtin/25-critique.md  -> the under-50-lines rule  5. jaira update --lanes --dry-run  -> lists 25-critique, 27-testing, 40-review as changed  6. cat core/release/notes/R2X4TJ.md  -> three one-line notes"
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
- **2026-11-22 15:01 · Alexander Sacharov** — critique+optimize in one pass (fix since review: 26 lines): both fix-now answered, nothing to clean
- **2026-11-22 15:01 · Alexander Sacharov** — waiting: bughunt round 3 since 16:01
- **2026-11-22 15:03 · Alexander Sacharov** — bughunt: pre-existing: core/mod/builtin/jaira-cockpit/hooks/board.ts:667 - class: a review line whose emphasis closes after the colon (`**note:** ...`) - PREFIXED strips `**note:` but leaves the closing `**`, so every mark after it is missed; also board.ts:668 PRE_EXISTING (such a pre-existing note is classed small edge, not "not this ticket") and board.ts:782 followTitle (the new fix-now strip inherits it); example: followTitle("**note:** **fix-now:** a") gives "** **fix-now:** a", and followTitle("**note:** the guard") gives "** the guard"; blame: 36a787f03 (MFYM1A); ran: node --experimental-strip-types on followTitle/gapItems with these inputs
- **2026-11-22 15:03 · Alexander Sacharov** — bughunt: done at 96f28e718 - 1 finding (0 change, 1 pre-existing) in 8 modules, follow-up hunt since b6777c680 (diff of 40-review.md, R2X4TJ.md, shipped.sha256; 25-critique.md, 27-testing.md, jaira-worker/SKILL.md unchanged; board.ts read whole, board.test.ts new test and its fixtures read); concurrency, error paths, input, permissions checked; ran: claude plugin test . (303 pass), go test ./core/shipped/ ./core/lane/... ok, sha256 of current 40-review.md and board.ts found in shipped.sha256; files: core/lane/builtin/25-critique.md core/lane/builtin/27-testing.md core/lane/builtin/40-review.md core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/R2X4TJ.md core/role/builtin/jaira-worker/SKILL.md core/shipped/shipped.sha256
- **2026-11-22 15:03 · Alexander Sacharov** — waiting: fresh review agent round 2 (run-lane.sh tab) since 16:03
- **2026-11-22 15:05 · Alexander Sacharov** — review: done at 96f28e718
- **2026-11-22 16:16 · Alexander Sacharov** — give-back (Teamlead 22.11.): zwei kleine Ränder schließen: (1) 25-critique.md:70 die one-pass-Regel bekommt dieselbe Lock/Parser-Ausnahme und den merge-base-Check wie 40-review.md:133-144; (2) Größen-pathspec in 25-critique.md:73 und 40-review.md:136,139 auf '-- :/ :(exclude,top).jaira' wie jaira-role-lane.
- **2026-11-22 16:17 · Alexander Sacharov** — give-back closed: 25-critique.md one-pass rule now has the merge-base ancestor check and the lock/parser exception of 40-review.md; size pathspec in both prompts is '-- :/ :(exclude,top).jaira' (checked from a subdirectory in a scratch repo: .jaira/ excluded, untracked new file listed); shipped.sha256 regenerated, go test ./core/shipped/ ./core/lane/... ok
