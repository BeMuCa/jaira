---
id: 01M897AW0X2HZVWVGHRM574VCW
title: "Cockpit Discuss: Entscheidung geht in die Arbeits-Lane des Boards, nicht fest in-progress"
status: in-progress
ready: true
creator: Sascha
goal: "Nach einer Entscheidung im Cockpit (Discuss/Explain) landet das Ticket in der Lane, die Give back seit PGXJG0 berechnet, auf jedem Board."
context: |-
  Der Discuss-Pfad im Cockpit scheitert auf Boards ohne Lane in-progress (z.B. ein Video-Board mit todo, script, script-approval, build, draft-approval, voice-approval, voice, final, acceptance, done).
  Ursache: core/mod/builtin/jaira-cockpit/hooks/texts.ts:245 sagt dem Agenten fest: jaira move <id> --to in-progress.
  PGXJG0 (24./25.11.2026) hat Give back schon auf backLaneOf (board.ts) umgestellt; der Discuss-Text blieb aussen vor, Review-Notiz von PGXJG0.
  Fix-Idee: denselben backLaneOf-Wert in den Text einsetzen.
definition-of-done: Discuss-Text nennt die Lane aus backLaneOf statt in-progress; Test in core/mod/builtin/jaira-cockpit/tests
tags:
  - cockpit
  - lanes
blocked-by: []
related: []
follows: 01M87RTTZFTY9RV9R3XPPGXJG0
commits: []
created-at: 2026-11-25T11:26:56Z
updated-at: 2026-11-25T11:46:07Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-905098
claimed-at: 2026-11-25T11:38:24Z
outcome-what: "Cockpit Discuss moves a decided ticket back to backLaneOf (w.backTo), in this chat and via spawn.sh --back <lane> in the Herdr tab"
outcome-why: boards without in-progress (video board) refused the hardcoded move; Give back already used backLaneOf since PGXJG0
outcome-resolves: Discuss decision on a video board lands in its working lane (build/script/final)
review-summary: none
review-gaps: "none - no other hardcoded Discuss move left (grep 'to in-progress' in cockpit hooks and dispatcher scripts is empty); nothing added is uncalled; openChat's back default serves the Explain callers"
test-verdict: "fail: gates green (go test ./core/... RC=0, claude plugin test 588/0), DoD 1-3 verified; bughunt change finding sends it back"
hunt-verdict: "fail: 1 change finding - register.tsx passes --back to an older spawn.sh, which refuses it (unknown flag) even on the default board"
round: 2
back-from: testing
back-by: lane
back-count: 1
---

# Cockpit Discuss: Entscheidung geht in die Arbeits-Lane des Boards, nicht fest in-progress

## Definition of Done

- [x] Discuss-Text nennt die Lane aus backLaneOf statt in-progress; Test in core/mod/builtin/jaira-cockpit/tests
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx openChat --back w.backTo; test 'Discuss on a video board moves the work back to its build lane' (tests/board.test.ts); core/role/role_test.go TestSpawnDiscussReturnsItsDecision
- [x] Default-Board: Text nennt weiter in-progress
  proof: tests/board.test.ts 'an open question shows under its ticket with Discuss' (no --back on the default board) and 'without Herdr, Discuss talks the question through in this chat' (move --to in-progress); spawn.sh default back=in-progress
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/574VCW.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-25 11:38 · Sascha** — parent session: uds:/run/user/1000/cc-socks/431929.sock
- **2026-11-25 11:39 · Sascha** — decision (worker): Herdr-Weg hat dieselbe feste Zeile in spawn.sh:644 (Kommentar dort: Aenderung an einem gehoert in den anderen). Darum neues spawn.sh-Flag --back <lane> nur fuer discuss, Default in-progress; openChat reicht w.backTo durch. Kein offener Punkt fuer die Person.
- **2026-11-25 11:42 · Sascha** — critique: none - one flag on spawn.sh mirrors the existing --copy pattern; backTo is the Waiting field Give back already uses (board.ts:192); openChat default in-progress only serves Explain, which never moves.
- **2026-11-25 11:42 · Sascha** — waiting: bughunt (background subagent) since 12:42
- **2026-11-25 11:44 · Sascha** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1194 - class: the cockpit always hands spawn.sh a flag an older spawn.sh does not know, also where the value is the default and the flag says nothing new (same pattern as --copy at :1195, which is older); example: a project whose .claude/skills/jaira-dispatcher/scripts/spawn.sh is a pre-574VCW copy (edited by hand, which jaira update --roles never replaces, and spawnScript() prefers the project copy over the personal one) - every Discuss press, also on the default board where back is in-progress, now fails with 'unknown flag: --back' instead of opening the chat as before; passing --back only when back !== 'in-progress' keeps the old path working; not a blocker: it fails loudly and the release note tells to run jaira update --roles; ran: bash <spawn.sh of fa4fe9401^> --no-worktree --no-parent --back in-progress ab12cd AB12CD discuss (exit 2, unknown flag: --back)
- **2026-11-25 11:44 · Sascha** — hunt report: 1 finding (1 change, edge, no blocker; 0 pre-existing) as a note on 574VCW | read whole: spawn.sh flag parse and discuss prompt, board.ts backLaneOf/snapshotOf, register.tsx openChat/spawnScript/parentArgs/herdr and every openChat caller, texts.ts discussInThisChat; register.tsx (3018 lines) not read end to end, only the parts the change reaches | not run: the cockpit tests (claude plugin test), left to the testing lane; checked: lane ids are validID [a-z0-9-] so --to ${back} cannot inject into the agent's shell line, jaira move needs no --board for a ticket on another board
- **2026-11-25 11:44 · Sascha** — bughunt: done at fa4fe9401 - 1 finding (1 change, 0 pre-existing) in 5 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/574VCW.md core/role/builtin/jaira-dispatcher/scripts/spawn.sh core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-25 11:46 · Sascha** — fix round 2: --back only when the lane is not in-progress (bughunt: a hand-edited older spawn.sh kept by update --roles refuses the flag; default board must keep working)
