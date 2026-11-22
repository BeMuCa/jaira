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
updated-at: 2026-11-22T14:51:31Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3877259
claimed-at: 2026-11-22T14:47:51Z
outcome-what: "Review marks notes fix-now or note and sends fix-now back once; small fix rounds run lean (stop-after-fix, critique+optimize in one pass, delta review from 'review: done at <commit>'); testing still hunts when a lock or parser changed"
outcome-why: Small fixable notes reached signoff and each give-back cost a full lane run (9 times on 2026-11-21)
outcome-resolves: "Tickets arrive at signoff without fixable edges, and a few-line fix no longer costs a full pass"
review-summary: |-
  core/lane/builtin/25-critique.md: the one-pass rule measures the payload diff, which on a fix-now round is the whole change, so a 5-line fix on a 300-line ticket never gets the one pass the worker skill promises - measure from review: done at <commit> when the notes carry one
  core/lane/builtin/40-review.md: delta size via git diff --stat prints per-file bars; --shortstat gives the one line to count
round: 2
back-from: critique
back-by: lane
back-count: 1
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
