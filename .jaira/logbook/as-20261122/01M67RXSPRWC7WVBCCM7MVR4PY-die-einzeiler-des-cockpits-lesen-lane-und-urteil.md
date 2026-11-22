---
id: 01M67RXSPRWC7WVBCCM7MVR4PY
title: Die Einzeiler des Cockpits lesen Lane und Urteil aus einer festen Schlusszeile statt aus dem ersten passenden Wort
status: done
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
commits:
  - 37403d24eb2298d0b436c06334dd004377843836
  - 7b31318f743d06d640b8475303762b7f299a231b
  - f7d381741929fa423acbd836dfdc98ce734194be
created-at: 2026-10-31T01:24:29Z
updated-at: 2026-11-22T17:03:13Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1092384
claimed-at: 2026-11-22T16:25:37Z
outcome-what: "digest.ts no longer guesses: without the closing status line the one-liner shows the ticket id (plus an open question or needs-you line); jaira-role-lane now ends its report with the same status line as worker and dispatcher"
outcome-why: "laneOf/verdictOf picked the first arrow or verdict word anywhere in a report, so 'reject? reason', 'nothing rejected' or 'all the way to signoff' became a wrong lane or verdict"
outcome-resolves: "the cockpit one-liner reads lane and verdict only from the fixed closing line; checked: concurrency n/a (pure function), error paths (failure detection unchanged, keeps the report whole), unusual input (empty/no ticket passes whole; CRLF and frame tags covered by the status-line tests), permissions n/a (fewer regexes over agent text)"
review-summary: "Round 3 on top of the accepted rounds 1-2: the cockpit's status-line pattern (STATUS_HEAD in hooks/digest.ts) and the start-line pattern (START in hooks/sessions.ts) now accept lane ids with digits (qa2, 2-review, qa2->3-ship). Before, such a status line was not recognised at all and - since round 1 removed the free-text guessing - the one-liner fell back to the bare id. Overall the branch: digest() reads lane and verdict only from the closing 'jaira: <ID> <from> -> <to> · ...' line, shows the id alone without it (questions and needs-you still shown), jaira-role-lane ends its report with that line, jaira-worker's text agrees, one release note per change."
review-gaps: |-
  note: the comments at hooks/digest.ts:84 and hooks/sessions.ts:38 say the pattern is what core/lane accepts; core/lane/lane.go validID also allows a leading hyphen (-qa), which the cockpit pattern refuses - fix: say "except a leading hyphen" or make validID refuse it (testing note, teamlead ruled no round for wording)
  note: pre-existing: core/lane/lane.go validID accepts a lane id starting with a hyphen, which the CLI then reads as a flag (jaira move --to -qa)
test-verdict: "pass: round 3 delta since 37403d24e; claude plugin test 319 pass 0 fail (RC=0) incl. 'status line: lane id with *' (4 cases) and the qa2 start line asserting the stored lane; go vet core + go test core/role core/mod core/shipped core/release RC=0; shipped-history --check RC=0; DoD 1-6 verified in the tree; writes nothing new, two-writer pass n/a"
hunt-verdict: "pass: bughunt done at f7d381741, 1 change (not a blocker: leading-hyphen lane id, unreachable from the CLI; note for the person), 0 pre-existing"
review-verdict: "accept: DoD 1-6 met; digit lane ids read in both patterns with table tests; manifest hashes match the four changed files; 319 cockpit tests green (claude plugin test)"
review-round: 3
review-verdict-gaps: |-
  note: the comments at hooks/digest.ts:84 and hooks/sessions.ts:38 say the pattern is what core/lane accepts; core/lane/lane.go validID also allows a leading hyphen (-qa), which the cockpit pattern refuses - fix: say "except a leading hyphen" or make validID refuse it (testing note, teamlead ruled no round for wording)
  note: pre-existing: core/lane/lane.go validID accepts a lane id starting with a hyphen, which the CLI then reads as a flag (jaira move --to -qa)
review-check: "1. cd core/mod/builtin/jaira-cockpit && claude plugin test .  -> 319 pass, 0 fail  2. look for the tests 'status line: lane id with a digit at the end' and 'recognises its start line a lane whose id holds digits' in that output -> both (pass)  3. in the cockpit, a worker report ending 'jaira: <ID> qa2 -> signoff · verdict: approve' shows '<ID> · qa2 → signoff · approve'  4. a report without a 'jaira: ...' last line shows only '[shortened by jaira-cockpit] <ID>', never a lane or 'rejected'"
round: 3
back-from: signoff
back-by: person
back-count: 1
stop-after-fix: "false"
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
- [x] STATUS_HEAD reads lane ids with digits (qa2, 2-review) like core/lane validID; table test
  proof: core/mod/builtin/jaira-cockpit/hooks/digest.ts LANE/STATUS_HEAD; sibling hooks/sessions.ts START; tests 'status line: lane id with *' (cockpit.test.ts) and 'recognises its start line a lane whose id holds digits' (board.test.ts)

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
- **2026-11-22 16:47 · Alexander Sacharov** — review (round 2): accept - both doc contradictions closed, manifest hashes match, 314 cockpit tests green; open note: digit lane ids in STATUS_HEAD (bughunt)
- **2026-11-22 16:48 · Alexander Sacharov** — given back by teamlead (22.11.), round 3: close the last note - STATUS_HEAD in digest.ts also matches lane ids with digits ([a-z0-9][a-z0-9-]*), table test (qa2, 2-review). Testing hunts the delta since 37403d24e so the hunt covers HEAD.
- **2026-11-22 16:51 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/digest.ts:86 - class: lane ids that core/lane validID (core/lane/lane.go:477) accepts but the cockpit's lane pattern [a-z0-9][a-z0-9-]* rejects, i.e. a lane id starting with a hyphen; same pattern in core/mod/builtin/jaira-cockpit/hooks/sessions.ts:40 (START); example: status line 'jaira: 55CS2M -qa -> signoff · verdict: ok' is not recognised at all, so the one-liner drops lanes and verdict and shows the id alone, and '/jaira-role-lane MG677R -qa' records lane '' - the comments on both lines claim parity with core/lane, which is not true for this input; not a blocker: such a lane id is close to unusable from the CLI anyway (jaira move --to -qa reads as a flag), so either validID should refuse a leading hyphen or the comments should say the cockpit narrows it; blame: f7d3817 (names MVR4PY); ran: node --experimental-strip-types on statusLineOf/startOf (undefined, lane '')
- **2026-11-22 16:51 · Alexander Sacharov** — bughunt: done at f7d381741 - 1 finding (1 change, 0 pre-existing) in 6 modules, follow-up hunt since 37403d24e; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/digest.ts core/mod/builtin/jaira-cockpit/hooks/sessions.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/tests/cockpit.test.ts core/release/notes/MVR4PY.md core/role/builtin/jaira-role-lane/SKILL.md core/role/builtin/jaira-worker/SKILL.md core/shipped/shipped.sha256
- **2026-11-22 16:51 · Alexander Sacharov** — testing: note: bughunt round 3 (leading-hyphen lane id '-qa' passes core/lane validID but not LANE/START) is not a blocker - such a lane cannot be reached from the CLI ('--to -qa' reads as a flag), and [a-z0-9][a-z0-9-]* is the pattern the teamlead asked for. Left: the comments at digest.ts:83 and sessions.ts:36 say 'what core/lane accepts' - one word off for that one case. Teamlead ruled no further round for wording.
- **2026-11-22 16:53 · Alexander Sacharov** — review (round 3): accept - digit lane ids in STATUS_HEAD and START, table tests, manifest ok, 319 cockpit tests green; open note: comment wording vs validID leading hyphen
