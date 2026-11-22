---
id: 01M67RXSPRWC7WVBCCM7MVR4PY
title: Die Einzeiler des Cockpits lesen Lane und Urteil aus einer festen Schlusszeile statt aus dem ersten passenden Wort
status: critique
ready: true
creator: berkc
assignee: Alexander Sacharov
goal: "Ein Worker-Bericht endet mit einer festen Zeile (z.B. 'jaira: <ID> <von> -> <nach> · verdict: <...>'), und digest.ts liest Lane und Urteil nur aus ihr; fehlt sie, zeigt der Einzeiler nur die ID statt zu raten."
context: |-
  Was falsch ist: core/mod/builtin/jaira-cockpit/hooks/digest.ts raet. laneOf nimmt den ersten Pfeil irgendwo im Bericht (Zeile 67), verdictOf das erste Wort wie 'reject' irgendwo (Zeile 77-84). Ein Bericht, der 'reject? reason' oder 'nichts zurueckgewiesen' erwaehnt, wird zu '-> signoff · rejected'.
  Ausloeser: Maintainer 31.10.2026: 'why is our cockpit wrong this often? is it because haiku?' - Nein: kein Modell, reine Regex. Haiku klassifiziert nur die Threads.
  Bekannt: X6HB6V (signoff) korrigiert die Lane. Das Urteil bleibt falsch. Die Schlusszeile muss in jaira-worker, jaira-dispatcher und jaira-role-lane stehen.
definition-of-done: "Worker-, Dispatcher- und Lane-Rolle schreiben die feste Schlusszeile"
tags:
  - tui
blocked-by: []
related:
  - 01M5MNWVKBC8Q193EK7YX6HB6V
commits: []
created-at: 2026-10-31T01:24:29Z
updated-at: 2026-11-22T16:45:06Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1092384
claimed-at: 2026-11-22T16:25:37Z
outcome-what: "digest.ts no longer guesses: without the closing status line the one-liner shows the ticket id (plus an open question or needs-you line); jaira-role-lane now ends its report with the same status line as worker and dispatcher"
outcome-why: "laneOf/verdictOf picked the first arrow or verdict word anywhere in a report, so 'reject? reason', 'nothing rejected' or 'all the way to signoff' became a wrong lane or verdict"
outcome-resolves: "the cockpit one-liner reads lane and verdict only from the fixed closing line; checked: concurrency n/a (pure function), error paths (failure detection unchanged, keeps the report whole), unusual input (empty/no ticket passes whole; CRLF and frame tags covered by the status-line tests), permissions n/a (fewer regexes over agent text)"
review-summary: "The cockpit's digest (core/mod/builtin/jaira-cockpit/hooks/digest.ts) no longer guesses: laneOf/verdictOf and the LANES list are deleted, so a report without the closing 'jaira: <ID> <from> -> <to> · … · verdict: …' line is shortened to the ticket id (plus an open question / needs-you line if there is one) instead of a lane or verdict picked out of free text. jaira-role-lane SKILL.md now ends every report with that same status line, so lane workers get the full one-liner too. Tests replace the old guessing cases with nine free-text reports (reject inside a question, negated verdict words, a quoted status line, arrows) that must yield only the id; release note and shipped.sha256 updated."
review-gaps: |-
  note: core/role/builtin/jaira-worker/SKILL.md:465 still says 'a report without it is guessed at, field by field' - now false; should say the cockpit then shows only the ticket id
  note: core/role/builtin/jaira-role-lane/SKILL.md:323 (question case: 'numbered list as the last thing in your pane') and :332 (role question 'as the last line') now contradict the new 'always the last line' status-line rule - say the status line comes after the list, or that a role question writes none
test-verdict: "pass: claude plugin test (2.1.330) 314 pass 0 fail; go vet + go test ./... RC=0; shipped-history --check RC=0. DoD 1-4 verified in the tree: role-lane Boundaries carries the status line, digest() ': [ticket]', 9 free-text cases incl. 'reject? reason' and 'nothing rejected' give the bare id, release note present. Function through the real session.receive/prompt.submit hooks: status line read, no line → id only, failed report stays whole. Writes nothing new: two-writer pass n/a"
hunt-verdict: "pass: bughunt done at 37403d24e, 1 change (not a blocker: a lane id with digits on a custom board gives the bare id; note for the person), 0 pre-existing"
review-verdict: "accept: all four DoD items met in the diff; free-text guessing is gone, 314 cockpit tests green; two doc-drift notes only"
review-round: 1
review-verdict-gaps: |-
  note: core/role/builtin/jaira-worker/SKILL.md:465 still says 'a report without it is guessed at, field by field' - now false; should say the cockpit then shows only the ticket id
  note: core/role/builtin/jaira-role-lane/SKILL.md:323 (question case: 'numbered list as the last thing in your pane') and :332 (role question 'as the last line') now contradict the new 'always the last line' status-line rule - say the status line comes after the list, or that a role question writes none
review-check: "1. cd core/mod/builtin/jaira-cockpit && claude plugin test .  2. expect '314 pass, 0 fail', including the 9 tests named 'no status line, only the ticket: …'  3. open core/mod/builtin/jaira-cockpit/hooks/digest.ts and search laneOf / verdictOf - no hits  4. live: claude --plugin-dir <worktree>/core/mod/builtin/jaira-cockpit, let a subagent send 'T8QS1Y: moved to human with --question \"reject? reason: dup\"' - the chat shows '[shortened by jaira-cockpit] T8QS1Y · question: reject? reason: dup — …' with no lane and no 'rejected'"
round: 2
back-from: signoff
back-by: person
back-count: 1
stop-after-fix: "true"
---

# Die Einzeiler des Cockpits lesen Lane und Urteil aus einer festen Schlusszeile statt aus dem ersten passenden Wort

## Definition of Done

- [x] jaira-role-lane schreibt die feste Schlusszeile wie Worker und Dispatcher (55CS2M); diese zwei haben sie schon
  proof: core/role/builtin/jaira-role-lane/SKILL.md: Boundaries, last item (status line)
- [x] digest.ts: ohne Schlusszeile zeigt der Einzeiler nur die ID - laneOf/verdictOf raten nicht mehr aus dem Fliesstext; ein Bericht mit 'reject' im Text ergibt kein 'rejected' (Test)
  proof: core/mod/builtin/jaira-cockpit/hooks/digest.ts digest(): ': [ticket]'; test 'no status line, only the ticket: *' (9 cases) in tests/cockpit.test.ts
- [x] Release-Notiz unter core/release/notes/
  proof: core/release/notes/MVR4PY.md
- [x] acceptance rule: a one-liner with only the ID is better than one with a wrong lane or verdict
  proof: tests/cockpit.test.ts guessCases: every free-text case expects the bare id
- [x] review notes closed: jaira-worker SKILL no longer says a report without the status line is guessed at; jaira-role-lane question cases put the status line after the list
  proof: core/role/builtin/jaira-worker/SKILL.md:465; core/role/builtin/jaira-role-lane/SKILL.md:324, :333, :364

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 16:25 · Alexander Sacharov** — Alexander Sacharov took this ticket over from berkc
- **2026-11-22 16:25 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 16:32 · Alexander Sacharov** — decision (Alex 22.11.): Ticket auf den Rest verkleinert. Die Schlusszeile in worker/dispatcher und das Lesen in digest.ts brachte 55CS2M (e7046a407). Offen: (1) Schlusszeile in jaira-role-lane, (2) ohne Schlusszeile nur die ID, kein Raten. DoD 1-2 umformuliert.
- **2026-11-22 16:33 · Alexander Sacharov** — decision (worker): without a status line, question:/needs you: are still read from the report text. They are not lane or verdict, and dropping them would hide what waits for a person. Only laneOf/verdictOf go. Planning lane skipped: the Options checklist does not tick it, and the method is in the DoD.
- **2026-11-22 16:38 · Alexander Sacharov** — critique: none. Checked: the lane role points at the worker's Report section, which carries the format (core/role/builtin/jaira-worker/SKILL.md:412); digest() keeps one branch per case; no other module guessed lanes from text (laneOf/verdictOf had no callers outside digest.ts).
- **2026-11-22 16:40 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/digest.ts:85 (STATUS_HEAD) + :162 - class: a status line the head regex does not recognise now loses lane and verdict entirely, because the free-text fallback that used to catch them is gone; the regex only accepts lane ids [a-z][a-z-]*, while core/lane/lane.go:477 (validID) allows digits anywhere in a lane id, and the new jaira-role-lane SKILL.md:354 tells every lane on every board to write this line; example: 'jaira: AB12CD qa2 -> signoff · verdict: approve' digested to '[shortened by jaira-cockpit] AB12CD — full report…' (before 37403d24e: 'AB12CD · approve'); likewise '2-review -> signoff · verdict: request changes' lost 'request changes'; not a blocker: nothing wrong is shown, only less; the regex line is from e7046a4 (55CS2M), but losing the fallback is this change; ran: node --experimental-strip-types scratch script against digest() at HEAD and at 37403d24e^
- **2026-11-22 16:40 · Alexander Sacharov** — bughunt: done at 37403d24e - 1 finding (1 change, 0 pre-existing) in 3 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/digest.ts core/mod/builtin/jaira-cockpit/tests/cockpit.test.ts core/release/notes/MVR4PY.md core/role/builtin/jaira-role-lane/SKILL.md core/shipped/shipped.sha256
- **2026-11-22 16:40 · Alexander Sacharov** — testing: note: bughunt change finding confirmed (core/lane/lane.go:477 validID allows digits; STATUS_HEAD lane groups are [a-z][a-z-]*). Not a blocker under the default rule: the line shows less (bare id), never a wrong lane or verdict — the acceptance rule of DoD 4. Fix if wanted: lane groups [a-z0-9][a-z0-9-]* in STATUS_HEAD plus a test case 'qa2 -> signoff'; one line, no other reader of that regex.
- **2026-11-22 16:42 · Alexander Sacharov** — review: accept - DoD 1-4 met, 314 cockpit tests green; two doc-drift notes (jaira-worker SKILL.md:465 'guessed at', jaira-role-lane question section vs. 'always last line') left for the person
- **2026-11-22 16:43 · Alexander Sacharov** — given back by teamlead (22.11.): two small text contradictions from review close in this round - worker SKILL:465 'guessed at, field by field' -> 'shows only the ticket id'; role-lane SKILL:323/332 question cases vs 'always the last line' -> the status line comes after the list/question. Then delta review.
- **2026-11-22 16:44 · Alexander Sacharov** — stop-after-fix set: this round changes only two skill texts and the manifest, no code - a further hunt finds nothing worth a round (teamlead asked for a light round).
