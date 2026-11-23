---
id: 01M8488Z540RWJ8F6J8CDKMG91
title: "Cockpit-Abnahme: offene 'seen by a person'-Punkte stehen als eigener Abschnitt da und werden mit Accept abgehakt"
status: done
ready: true
creator: Alexander Sacharov
goal: "Eine Person nimmt ein Ticket mit offenem 'seen by a person'-DoD in einem Schritt an: die Ansicht zeigt diese Punkte als eigenen Abschnitt 'Selbst prüfen', und Accept hakt sie ab und nimmt dann an - statt mit 'criterion N is still open' abzulehnen."
context: |-
  Was heute falsch ist: Ein Ticket mit einem DoD-Punkt 'seen by a person: ...' kommt in signoff mit diesem Punkt offen an - das ist Absicht, nur die Person kann ihn abhaken. Drückt der Maintainer im jaira-cockpit (Abnahme-Ansicht, Приём) Accept, lehnt jaira ab: 'the definition of done is not met: criterion 7 (...) is still open. Satisfy it, mark it with jaira dod <ID> 7 --done, then try this move again'.
  Passiert am 23.11.2026 zweimal hintereinander (FMWBA8 DoD 7, MK8HD2 DoD 5). Jedes Mal musste der Teamlead den Punkt per CLI abhaken und der Maintainer ein zweites Mal Accept drücken.
  Die Ansicht zeigt heute nirgends, dass noch ein Punkt auf die Person wartet - man erfährt es erst aus der Ablehnung.
  Wunsch des Maintainers: ein eigener Absatz in der Ansicht, der zeigt, was er noch selbst tun/ansehen muss; das tun, dann annehmen - in einem Durchgang.
  Wo: core/mod/builtin/jaira-cockpit/hooks/register.tsx (Abnahme-Ansicht, Execute/Accept, executionOf in hooks/board.ts), DoD-Daten aus 'jaira show --json' (dod_items: n, text, done).
  Erkennung: DoD-Text beginnt mit 'seen by a person' (so schreiben die Lanes ihn heute).
definition-of-done: "Hat das Ticket offene DoD-Punkte, die mit 'seen by a person' beginnen, zeigt die Abnahme-Ansicht sie oberhalb der Review-Notizen als eigenen, abgesetzten Abschnitt (Titel z.B. 'Selbst prüfen'), jeden Punkt mit seinem Text"
tags:
  - cockpit
  - acceptance
blocked-by: []
related: []
commits:
  - 2291675a0a58d01f36218d9be1922396e3faa07d
  - 0bba5eb40d0b02aa4d1bd8e69961a0acc727d97e
  - 40153c99d91f666dc89b835af64f0ce619cc8ec7
created-at: 2026-11-23T13:07:10Z
updated-at: 2026-11-23T16:20:02Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-97038
claimed-at: 2026-11-23T13:07:36Z
outcome-what: "Acceptance pane lists open 'seen by a person' DoD items under 'Check yourself' above the review notes; Execute needs each ticked, re-reads the ticket, ticks still-unchanged items with proof 'seen by <whoami> in the cockpit' before the gates and accepts in one press; a later refusal names what was ticked"
outcome-why: "Accept was refused with 'criterion N is still open' for items only a person can tick (FMWBA8, MK8HD2 on 23.11.)"
outcome-resolves: "one-press acceptance of tickets with 'seen by a person' items"
review-summary: "Full review (no earlier review). The cockpit's Acceptance pane now reads a ticket's open DoD items starting 'seen by a person:' from 'jaira show --json' when it opens and lists them under 'Check yourself' above the review notes, each a toggle [ ]/[✓] (seenAtom). Execute refuses with the unticked items named until all are ticked. With all ticked, Accept runs a new 'before' hook inside the ticket's checkout: it re-reads dod_items, stops if an item was superseded, reworded or removed since the pane opened (seenToTick), skips one ticked meanwhile, and runs 'jaira dod <id> <n> --done --proof \"seen by <whoami> in the cockpit\"' for the rest - then the usual dry run, gates and move. A refusal after some ticks names them (refusedTicked) instead of 'nothing was changed'. A ticket without such items behaves as before. Release note core/release/notes/DKMG91.md."
review-gaps: |-
  note: edge: register.tsx:705 heldInSignoff (pr-when=signoff boards) and the failed --force move report through fail() after the before-hook already ticked items, so that message does not say items were ticked - pass the ticked list into those messages the way refusedTicked does
  note: limit: ticks written before the gates stay ticked when a later gate refuses (worker decision 3, the person did look)
  note: limit: only the Acceptance pane's Execute ticks; Accept from an Explain choice still refuses on an open 'seen by a person' item (worker decision 4)
test-verdict: "pass: round 2 - cockpit suite 492 pass (RC=0), shipped-history --check ok; both hunt findings retested by tests: dod failing on item 5 -> refusal names item 2 as ticked, no move; gate refusal after ticks names 2, 5; reworded/superseded item -> no dod, no move, seenChanged shown; item ticked meanwhile -> only the other ticked, proof kept; DoD 1-6 verified"
hunt-verdict: "pass: bughunt done at 40153c99d, 1 change edge (signoff-hold/forced-move-failure messages after ticks do not name them; ticks are true, re-press skips them), 0 pre-existing, note for the person"
round: 2
back-from: testing
back-by: lane
back-count: 1
review-verdict: "accept: all six DoD items met in the diff, 492 cockpit tests pass; the re-read before ticking closes the stale-pane case; no blocker"
review-round: 2
review-verdict-gaps: |-
  note: edge: register.tsx:705 heldInSignoff (pr-when=signoff boards) and the failed --force move report through fail() after the before-hook already ticked items, so that message does not say items were ticked - pass the ticked list into those messages the way refusedTicked does
  note: limit: ticks written before the gates stay ticked when a later gate refuses (worker decision 3, the person did look)
  note: limit: only the Acceptance pane's Execute ticks; Accept from an Explain choice still refuses on an open 'seen by a person' item (worker decision 4)
review-check: "1. cd core/mod/builtin/jaira-cockpit && claude plugin test .  -> 492 pass, 0 fail; describe 'what the person checks themselves' covers the pane  2. live: on a scratch ticket in signoff add 'jaira dod <id> --add \"seen by a person: the key names the branch\"'  3. claude --plugin-dir core/mod/builtin/jaira-cockpit, open the ticket's Acceptance pane  4. see 'Check yourself - tick what you have seen:' above the review notes with '[ ] <n>. seen by a person: ...'  5. press Execute unticked -> the line names the item, nothing moves  6. tick it ([✓]), press Execute -> ticket goes to done; 'jaira show <id>' shows the item [x] with proof 'seen by <you> in the cockpit'"
---

# Cockpit-Abnahme: offene 'seen by a person'-Punkte stehen als eigener Abschnitt da und werden mit Accept abgehakt

## Definition of Done

- [x] Hat das Ticket offene DoD-Punkte, die mit 'seen by a person' beginnen, zeigt die Abnahme-Ansicht sie oberhalb der Review-Notizen als eigenen, abgesetzten Abschnitt (Titel z.B. 'Selbst prüfen'), jeden Punkt mit seinem Text
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2383; test 'what the person checks themselves › open items stand above the review notes...'
- [x] Jeder dieser Punkte hat einen Haken, den die Person setzt; Accept ist erst möglich, wenn alle gesetzt sind, und sagt sonst in der Zeile, was noch fehlt
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1496 (needSeen) + texts.ts summary 'needSeen'; test 'Execute with an item not ticked accepts nothing and names the item'
- [x] Accept mit gesetzten Haken ruft für jeden Punkt 'jaira dod <id> <n> --done --proof "seen by <Person> in the cockpit"' auf und dann die Annahme - ein Druck, keine Ablehnung wegen offener Punkte
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1638 acceptWith before-hook (re-reads dod_items, board.ts seenToTick); tests 'ticked, one Execute ticks it...', 'an item ticked meanwhile keeps its own proof'
- [x] Ein Ticket ohne solche Punkte sieht und verhält sich wie bisher
  proof: test 'a ticket without such items shows no section, and Execute accepts as before'; full suite 486 pass
- [x] Tests in core/mod/builtin/jaira-cockpit/tests: Abschnitt erscheint nur bei offenen Punkten; Accept ohne Haken gesperrt; mit Haken werden die dod-Aufrufe vor dem move gemacht
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts describe 'what the person checks themselves' (12 tests)
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/DKMG91.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 13:07 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/6958.sock
- **2026-11-23 13:08 · Alexander Sacharov** — decisions (worker, none open for the person): 1) section title in English ('Check yourself') - texts.ts keeps every cockpit text English; 'Selbst prüfen' in the DoD was 'z.B.'. 2) person name in the proof = 'jaira whoami' identity, 'a person' where it says none. 3) ticks are written right before the dry run of the move: a later gate refusal leaves them ticked - true anyway, the person did look. 4) only the Acceptance pane's Execute ticks; Accept from an Explain choice keeps today's refusal (outside this DoD). 5) items come from 'jaira show --json' read when the pane opens; path: todo -> in-progress (no planning option).
- **2026-11-23 13:14 · Alexander Sacharov** — in-progress: tests via 'claude plugin test .' in core/mod/builtin/jaira-cockpit (486 pass, ~60s). No typecheck in CI for the mod. Unrelated dirty files in worktree (.jaira/lanes/*, AGENTS.md, CLAUDE.md - jaira's own regen) left out of commits.
- **2026-11-23 13:15 · Alexander Sacharov** — critique: none. Checked: before-hook in accept() is needed because the ticket's checkout (dir) is only known inside accept/ticketHere; jaira dod --done on an already ticked item exits 0, so a second press after a gate refusal re-ticks harmlessly; needWhy precedes needSeen on purpose (reason first, then ticks).
- **2026-11-23 13:19 · Alexander Sacharov** — testing: note: edge: a jaira dod that fails on the 2nd of several ticked items leaves the 1st ticked while the pane says 'nothing was changed'; harmless (re-tick exits 0, the person did see it), fix would be a wording change in t.refused for this path.
- **2026-11-23 13:19 · Alexander Sacharov** — waiting: bughunt (background subagent) since 14:19
- **2026-11-23 13:20 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1647 - class: Execute writes 'jaira dod <id> <n> --done --proof ...' from the seen list snapshotted at pane open (openAcceptance:1369, read in ticketDir without a pull, so possibly from a stale disk copy) and never re-checks that item n is still that open 'seen by a person' item; jaira dod --done overwrites any state and proof unconditionally. example: while the pane is open someone supersedes item 3 ([-]), rewords it with 'dod 3 --text', or ticks it with a real proof - Execute turns [-] back into [x], ticks wording the person never saw, or replaces the proof with 'seen by <me> in the cockpit'. fix: re-read dod_items in here.dir inside before() and tick only items whose n, text and open state still match the snapshot, else refuse. ran: scratch board, 'jaira dod A8SEGH 1 --superseded' then '--done --proof p' -> rc 0, item [x], proof replaced. Not a blocker by itself: needs a DoD edit on a signoff ticket while its pane is open.
- **2026-11-23 13:20 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1503 (with accept():695 and acceptWith():1647) - class: the seen items are ticked on the ticket BEFORE the dry-run gates and heldInSignoff, but every refusal after that is shown as t.refused = '<id> was not accepted, nothing was changed:' - partial writes reported as none. example: two 'seen by a person' items ticked in the pane plus an unrelated open item 'tests pass' -> Execute ticks 1 and 2 on the ticket, the dry run refuses on criterion 3, and the pane says nothing was changed; same when the 2nd of two dod calls fails (item 1 already ticked). fix: say which items were ticked in that refusal (or word t.refused without 'nothing was changed' when before() wrote). ran: read of accept()/accepting()/execute(); blame 2291675a0.
- **2026-11-23 13:20 · Alexander Sacharov** — hunt report: 2 findings, both change, neither a blocker, as bughunt notes on DKMG91 | read whole the accept path of register.tsx (accept, accepting, ticketHere, heldInSignoff, openAcceptance, execute, acceptWith, once, Acceptance pane render), board.ts DoD/execution helpers, texts.ts, types/index.d.ts, board.test.ts, release note; uncommitted .jaira/lanes/*.md, AGENTS.md, CLAUDE.md are 'jaira update' prompt text, read as a diff, nothing to find; cockpit suite ran green (claude plugin test: 486 pass) | not checked live in the TUI: the stale-seen-list race was reproduced only at the jaira dod level on a scratch board, not through a real pane
- **2026-11-23 13:20 · Alexander Sacharov** — bughunt: done at 0bba5eb40 - 2 findings (2 change, 0 pre-existing) in 6 modules, first hunt; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/types/index.d.ts core/release/notes/DKMG91.md
- **2026-11-23 13:21 · Alexander Sacharov** — testing: class: Accept writes DoD ticks before gates and reports later refusals as 'nothing was changed', and ticks from the pane-open snapshot without re-reading; example: bughunt notes 13:20 (register.tsx:1503, :1647); fix: in before() re-read dod_items in the ticket's dir, tick only items still open with the same text (skip ones done meanwhile, refuse on superseded/reworded), and name the ticked items in any refusal after them; tests over both.
- **2026-11-23 13:24 · Alexander Sacharov** — in-progress round 2: both bughunt findings fixed. before() re-reads 'jaira show --json' in the ticket's dir: an item done meanwhile is skipped (keeps its proof), one superseded/reworded/gone ticks nothing and refuses (seenChanged). A refusal after ticks says 'only item(s) N were ticked' instead of 'nothing was changed' (texts.refusedTicked). 492 pass.
- **2026-11-23 13:24 · Alexander Sacharov** — critique round 2: read only the delta 0bba5eb40..40153c99d (the hunt's two fixes); seenToTick compares oneLine text on both sides, the extra show runs in the same dir as the ticks; nothing to say.
- **2026-11-23 13:26 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:715-736 - class: Accept exits after the seen items were already ticked (before-hook, register.tsx:695, from 2291675a0) that do not go through refused, so the message never names the ticks: heldInSignoff hold (fail + staysInLane), movedBefore (movedFirst), and jaira move --force failing (fail with done.why); only the dry-run refusal and a failing dod name them via refusedTicked; example: board with pr-when signoff, person ticks item 5 and presses Accept, the PR lane starts or the PR is not settled, item 5 is now done on the ticket but the message only says the ticket stays in signoff; not a blocker (the ticks are true, nothing lost, a re-press skips them via seenToTick); read
- **2026-11-23 13:26 · Alexander Sacharov** — hunt report: 1 finding, 1 change (edge, not a blocker), noted on DKMG91 | diff since 0bba5eb40 read: hooks/board.ts (seenToTick), hooks/register.tsx (acceptWith/accept/accepting/heldInSignoff), hooks/texts.ts, tests/board.test.ts, .jaira/lanes/*.md, AGENTS.md, CLAUDE.md (prose prompts, uncommitted) | not run: the cockpit test suite (claude-code/testing harness, the testing lane runs it); show-then-dod window (TOCTOU) left open by design, as dod --done on an already-done item is harmless
- **2026-11-23 13:26 · Alexander Sacharov** — bughunt: done at 40153c99d - 1 finding (1 change, 0 pre-existing) in 6 modules, follow-up hunt since 0bba5eb40; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/types/index.d.ts core/release/notes/DKMG91.md
- **2026-11-23 13:26 · Alexander Sacharov** — waiting: bughunt round 2 (background subagent) since 14:26
- **2026-11-23 13:26 · Alexander Sacharov** — waiting: review (fresh agent, Herdr tab via run-lane.sh) since 14:26
- **2026-11-23 13:29 · Alexander Sacharov** — review: done at 40153c99d
- **2026-11-23 13:29 · Alexander Sacharov** — review: verdict accept, but move to signoff refused: the hunt gate sees uncommitted files the hunt at 40153c99d read - the only dirty files are jaira's own regen (.jaira/lanes/{critique,in-progress,review,testing}.md, AGENTS.md, CLAUDE.md), not this ticket's code. Dispatcher: commit or restore those (or --force as the person), then 'jaira move DKMG91 --to signoff'.
- **2026-11-23 16:20 · Alexander Sacharov** — report: Acceptance pane lists open 'seen by a person' items under 'Check yourself', Execute re-reads, ticks them (proof 'seen by <whoami> in the cockpit') and accepts in one press; 3 commits, 492 cockpit tests | signoff, context 197k of 1000k, cache warm 60 min | open for the person: 1 edge (signoff-hold/forced-move-failure messages do not name ticks), 2 limits (ticks stay on a later refusal; Explain-choice Accept unchanged); branch not pushed
