---
id: 01M86H5PFARVW0Z0NBMJN5ZZP6
title: "Cockpit: Execute in der Abnahme wird nach dem Klick grau und nimmt keinen zweiten Klick, bis der Schritt fertig ist"
status: critique
ready: true
creator: Sascha
assignee: Sascha
goal: "Nach einem Klick auf Execute (bzw. Accept) im Abnahme-Pane sieht die Person sofort, dass der Klick angekommen ist: der Knopf wird grau mit 'laeuft…', nimmt keinen weiteren Klick, und wird wieder normal, wenn der Schritt fertig ist (oder spaetestens nach einer Zeitgrenze)."
context: |-
  Was falsch ist: Im Cockpit-Abnahme-Pane (core/mod/builtin/jaira-cockpit/hooks/register.tsx, Knopf Execute, Texte in hooks/texts.ts:63-70) sieht der Knopf nach dem Klick genauso aus wie davor. Die Person weiss nicht, ob der Klick ankam, und klickt nochmal.
  Ein zweiter Klick wird heute erst abgewiesen, wenn er schon laeuft: texts.ts:64 busyRow ('an earlier press ... is still at work on it') - als Textzeile, nicht am Knopf.
  Ausloeser: Sascha (Alex) 24.11.2026: 'die Knopf zum Annehmen unten ... besser wenn ich klicke, dass sie nicht mehr klickbar ist, damit klar ist, dass ich gedrueckt habe - z.B. 5 Sekunden grau'.
  Gilt fuer Execute in allen Modi (accept, give back) und den Accept-Knopf des Panes.
definition-of-done: "Nach dem Klick ist Execute sofort grau/dim mit einem Text wie 'laeuft…' und reagiert nicht auf weitere Klicks oder Enter; Test"
tags:
  - cockpit
blocked-by: []
related: []
commits: []
created-at: 2026-11-24T10:21:08Z
updated-at: 2026-11-24T10:52:41Z
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-2605896
claimed-at: 2026-11-24T10:21:37Z
outcome-what: "Execute in the Acceptance pane reads 'running…' and is drawn dim while its row's lock is held; a press on it or Enter in the field does nothing then (the busyRow toast is gone). When BUSY_MS (20 s) frees Execute's lock while its step is still at work, one toast says so; other rows (Explain, release) stay silent. The lock is keyed per press, not per millisecond."
outcome-why: "A press gave no visible answer, so the person pressed again (Alex 24.11.2026)."
outcome-resolves: "DoD 1-3; checked: concurrency (two presses same ms - fixed key collision), error paths (refused step frees at once, tested), unusual input (Enter with text while locked), permissions (n/a, no outside input)"
review-summary: "Full review of ccf90e1d4..603483615. The Acceptance pane's Execute button and the field's submit label now read 'running…' (dim) whenever the row lock busyAtom holds the ticket; a press while locked goes into once(), which returns without doing anything - the old busyRow toast is gone. once() marks each lock with clock time plus a per-press counter, so two presses in the same millisecond cannot free each other's lock. When BUSY_MS (20 s) frees a lock whose action has not ended, once() says the optional 'late' text - only Execute's two callers (button, field Enter) pass one, so Explain, Discuss, Acceptance and the release row stay silent. A finished or refused step frees the lock in finally and Execute reads normal again. Tests cover running…/no second press/no Enter, refusal resets the label, the 20 s timeout toast said once, and Explain's held lock silent. Release note added. claude plugin test: 546 pass, 0 fail."
round: 4
back-from: review
back-by: lane
back-count: 1
review-gaps: |-
  note: edge: fix-now: Enter in the field with nothing typed while Execute reads running… still toasts nothingToFix (register.tsx:2583 runs before once()); DoD 1 says Enter does nothing while busy - return early when busy before the nothingToFix check, plus a test line in 'Execute reads running…' pressing Enter on an empty field
  note: limit: time limit is BUSY_MS=20 s, not the 30 s the DoD names as an example - decided in the ticket's notes
  note: limit: Execute also reads running… while an Explain chat holds the row lock - one lock per row by design, decided in critique
  note: elsewhere: core/mod/builtin/jaira-cockpit/node_modules/ is untracked and not gitignored - a test run leaves .vite/vitest/results.json, which lands in every lane's worktree_diff; .gitignore entry, new ticket
test-verdict: "pass: round 2 on 603483615 - claude plugin test . exit 0 (546 pass, 0 fail); DoD 1-3 re-verified by their tests; bughunt finding (stillRunning on non-Execute rows) closed and tested: Explain's held chat past 20 s toasts nothing. Change writes only the in-memory row lock. DoD 4 left for the person."
hunt-verdict: "pass: follow-up hunt at 603483615 - 0 findings; round 1's change finding fixed in 603483615"
review-verdict: "accept: DoD 1-3 met by the diff and its tests, no blocker; one fix-now edge (empty Enter while busy still toasts)"
review-round: 3
review-verdict-gaps: |-
  note: edge: fix-now: Enter in the field with nothing typed while Execute reads running… still toasts nothingToFix (register.tsx:2583 runs before once()); DoD 1 says Enter does nothing while busy - return early when busy before the nothingToFix check, plus a test line in 'Execute reads running…' pressing Enter on an empty field
  note: limit: time limit is BUSY_MS=20 s, not the 30 s the DoD names as an example - decided in the ticket's notes
  note: limit: Execute also reads running… while an Explain chat holds the row lock - one lock per row by design, decided in critique
  note: elsewhere: core/mod/builtin/jaira-cockpit/node_modules/ is untracked and not gitignored - a test run leaves .vite/vitest/results.json, which lands in every lane's worktree_diff; .gitignore entry, new ticket
review-check: "1. Build and install this branch's cockpit (claude --plugin-dir core/mod/builtin/jaira-cockpit)  2. Open the cockpit, find a ticket waiting in signoff, press Acceptance on its row  3. Press Execute: it turns dim and reads 'running…' at once  4. Press Execute again and press Enter in the field: nothing happens, no toast  5. The pane closes when the ticket is accepted, or the refusal shows in red and Execute reads 'Execute' again (this is the 'seen by a person' DoD item - tick it yourself)  6. Automated: cd core/mod/builtin/jaira-cockpit && claude plugin test . - 'Execute reads running…', 'a step refused reads Execute again', 'a step that has not ended after 20 seconds' pass"
---

# Cockpit: Execute in der Abnahme wird nach dem Klick grau und nimmt keinen zweiten Klick, bis der Schritt fertig ist

## Definition of Done

- [x] Nach dem Klick ist Execute sofort grau/dim mit einem Text wie 'laeuft…' und reagiert nicht auf weitere Klicks oder Enter; Test
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'Execute reads running… and takes no press and no Enter until its step is done'; hooks/register.tsx Execute button label/dim from busy
- [x] Ist der Schritt fertig (Erfolg oder Fehler), wird der Knopf wieder normal; endet nichts, spaetestens nach einer Zeitgrenze (z.B. 30 s) mit einem Hinweis; Test
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'a step refused reads Execute again at once, the refusal beside it', 'a step that has not ended after 20 seconds frees Execute and says so once'; hooks/register.tsx once() stillRunning
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/N5ZZP6.md
- [ ] seen by a person: in the Acceptance pane, Execute pressed on a ticket whose move takes a moment reads 'running…' dim until the pane closes or the refusal shows

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-24 10:21 · Sascha** — parent session: uds:/run/user/1000/cc-socks/3417.sock
- **2026-11-24 10:22 · Sascha** — open decisions: none. Settled: (1) the existing row lock (once/busyAtom, BUSY_MS=20s, register.tsx:1063) already blocks a second press and draws the button dim - the gap is only the label and the timeout hint; (2) busy Execute reads 'running…' and its press does nothing (was: busyRow toast), Input submitLabel follows; (3) time limit stays BUSY_MS=20s (DoD says 'e.g. 30 s'; 20 s is already tested and shipped); (4) when the limit frees a row whose action is still running, once() says so in one line - not for 'hold' actions, whose lock is the limit by design.
- **2026-11-24 10:35 · Sascha** — found out: (1) reading an atom ($.state.get) inside a $.clock.after callback is refused by the host ('no hooks module ... no scan to allow it') - update() works there, so once() learns whether its lock was still there from inside its own update; (2) the lock was keyed by clock.now(), so two presses in the same millisecond (band Acceptance, then Execute) freed each other's lock at BUSY_MS - real life almost never, the mock clock always; now a per-press counter is added. (3) busyRow text removed: a locked Execute now says 'running…' on itself; the old 'said, not dropped' test was changed to check the label instead.
- **2026-11-24 10:36 · Sascha** — critique: (1) move 'let presses' above once()'s doc comment. Let stand: Execute reads 'running…' also while an Explain chat holds the row lock (hold, 20 s) - the lock is one per row by design (comment at register.tsx 'One lock for the row and its pane'), and a label naming the holder would need a second state for one rare case; stillRunning applies to every once() caller, worded for 'its buttons', which is right for all of them.
- **2026-11-24 10:39 · Sascha** — waiting: bug hunt (background subagent) since 11:39
- **2026-11-24 10:41 · Sascha** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1087 - class: once() now speaks t.stillRunning(id) for every caller, but its text (texts.ts:69) assumes id is a ticket, while once() also locks the non-ticket release row (RELEASE='release', register.tsx:2362 Update, :2823 upgrade-confirm; jaira self upgrade runs with a 300 s timeout at :545); example: an upgrade slower than 20 s toasts 'release: still at work after 20 s - its buttons take a press again; look at the ticket before you press once more' - no ticket exists, and it invites a second Update while the binary is being replaced; ran: scratch copy of tests/board.test.ts holding 'jaira self upgrade', press release-update then upgrade-confirm, clock.advance(20_000) -> that exact toast (claude plugin test, exit 0); not a blocker: a misleading notice, nothing written wrong by this change
- **2026-11-24 10:41 · Sascha** — hunt report: 1 finding, 1 change (edge, no blocker), as a note on N5ZZP6 | read whole: once() and every caller of it, busyAtom and its reset on session.start, openChat, openAcceptance, previewUpgrade/upgradeJaira, the Acceptance pane render with Execute and the field's onSubmit in register.tsx; texts.ts; the lock tests in board.test.ts; release note N5ZZP6.md - the rest of the 2859-line register.tsx only by search for the lock's users | not checked: whether the host's update() from 'claude-code' may run its callback more than once (the freed flag is set inside it) - the engine's source is not in this repo
- **2026-11-24 10:42 · Sascha** — bughunt: done at 0127e7b20 - 1 finding (1 change, 0 pre-existing) in 4 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/N5ZZP6.md
- **2026-11-24 10:44 · Sascha** — fixed bughunt change finding: once() takes an optional 'late' text and says it only where given - Execute's two presses (button, field Enter) pass t.stillRunning; Explain, Discuss, Acceptance and the release row stay silent. Text now names Execute. Test: Explain's held chat past 20 s says nothing.
- **2026-11-24 10:45 · Sascha** — waiting: bug hunt round 2 (background subagent) since 11:45
- **2026-11-24 10:47 · Sascha** — hunt report: 0 findings, nothing written | read as a diff since 0127e7b20: hooks/register.tsx (once() whole and all its 8 callers; both execute() callers pass late, the rest stay silent), hooks/texts.ts, tests/board.test.ts; none read whole, no new files | nothing left unchecked; ran: claude plugin test . (546 pass, 0 fail)
- **2026-11-24 10:47 · Sascha** — bughunt: done at 603483615 - 0 findings (0 change, 0 pre-existing) in 4 modules, follow-up hunt since 0127e7b20; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/N5ZZP6.md
- **2026-11-24 10:47 · Sascha** — waiting: fresh review lane (run-lane.sh, Herdr tab) since 11:47
- **2026-11-24 10:50 · Sascha** — review: done at 603483615
- **2026-11-24 10:50 · Sascha** — review: fix-now: Enter on an empty field while Execute reads running… still toasts nothingToFix (register.tsx:2583) - return early when the row is busy, before the nothingToFix check, and assert no toast in 'Execute reads running… and takes no press and no Enter'
- **2026-11-24 10:52 · Sascha** — review fix-now closed: onSubmit returns before the nothing-to-fix check while busy (register.tsx onSubmit), test 'Execute reads running…' now also presses Enter on an empty field. Other review notes left as they are: 20 s limit and running… during Explain are decided; node_modules/.gitignore is elsewhere, for the person.
