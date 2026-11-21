---
id: 01M7ZYF17AVJAEWK40GG3DVVFZ
title: "Bei TUI-, Cockpit- oder OS-spezifischen Änderungen schreibt pre-process einen DoD-Punkt 'von einer Person angesehen'"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Was nur ein Mensch an echtem Terminal oder Betriebssystem prüfen kann, steht als eigener DoD-Punkt auf dem Ticket und wird bei der Abnahme gesehen, nicht vergessen."
context: |-
  Was heute fehlte (21.11.):
  - 02D92C: Emoji-Breite in Windows Terminal nur über die Unicode-Breitenklasse belegt, kein Mensch hat hingesehen; der DoD-Punkt galt trotzdem als erledigt.
  - JDDCEZ: Test lief lokal und auf Linux grün, fiel erst auf macOS-CI (kein timeout(1)).
  - AX8K02: der Fall "Person setzt done ohne --force" liess sich vom Agenten nicht spielen, belegt nur über den Refusal-Text.
  Die Lanes behandeln solche Punkte wie alle anderen und haken sie mit einem Ersatzbeweis ab.
  Gewünscht: der shipped pre-process-Lane (core/lane/builtin, Datei für pre-process) schreibt bei Änderungen an TUI, Cockpit, Terminal-Ausgabe oder OS-spezifischem Verhalten einen eigenen DoD-Punkt "von einer Person angesehen: <was, wo>" mit jaira dod --add. Agenten-Lanes haken diesen Punkt nie ab (Lane-Prompts von in-progress/testing/review sagen das), die Abnahme zeigt ihn als offenen Schritt.
  Offen für pre-process: woran der Lane erkennt, dass ein Punkt nur menschlich prüfbar ist (Liste von Pfaden wie core/mod/, internal/tui/ ist der Vorschlag).
definition-of-done: "pre-process-Lane-Prompt verlangt den DoD-Punkt 'von einer Person angesehen' für die genannten Änderungsarten"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-11-21T20:58:45Z
updated-at: 2026-11-21T21:12:14Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2704396
claimed-at: 2026-11-21T20:59:22Z
stop-after-fix: "false"
outcome-what: "Built-in lanes: pre-process adds 'seen by a person: <what, where>' DoD items for terminal-drawn output, platform-dependent behaviour and person-only actions; in-progress adds one when no planning ran and never ticks it; testing leaves it open (agent-ticked = finding, untick); review makes it a review-check step (agent-ticked = blocker). Acceptance (build.py plan) and signoff hand the person one line 'jaira dod <id> <n> --done --proof ... && jaira move ... --force' instead of a bare --force move. shipped.sha256 + release note."
outcome-why: "02D92C, JDDCEZ, AX8K02 ticked human-only checks with stand-in proofs; and --force on the accept move passes an unticked DoD, so the item would vanish unseen."
outcome-resolves: "DoD 1-4 with proofs; categories: prompts plus one python print path — no concurrency or writes; input: dod texts with leading spaces handled (strip), a ticket without dod gets the plain move; permissions: the tick is never run by an agent."
review-summary: none
round: 4
back-from: testing
back-by: lane
back-count: 2
review-gaps: none
test-verdict: "fail: bughunt change finding at build.py:145 - plan ticks done/superseded 'seen by a person:' items too"
hunt-verdict: "fail: bughunt done at 2871d8f71, 2 change (1 blocker)"
---

# Bei TUI-, Cockpit- oder OS-spezifischen Änderungen schreibt pre-process einen DoD-Punkt 'von einer Person angesehen'

## Definition of Done

- [x] pre-process-Lane-Prompt verlangt den DoD-Punkt 'von einer Person angesehen' für die genannten Änderungsarten
  proof: core/lane/builtin/15-pre-process.md:77 '## What only a person can check'
- [x] in-progress-, testing- und review-Prompt sagen: diesen Punkt hakt nur eine Person ab
  proof: core/lane/builtin/20-in-progress.md:164, 27-testing.md:57, 40-review.md:107
- [x] shipped.sha256 und Release-Note core/release/notes/<ID>.md
  proof: core/shipped/shipped.sha256 (scripts/shipped-history.sh --write, 903 versions), core/release/notes/3DVVFZ.md
- [x] Abnahme: jaira-role-acceptance (build.py plan) und jaira-role-signoff geben für einen offenen 'seen by a person:'-Punkt eine Zeile 'dod --done && move --force' an die Person; kein Agent hakt ihn ab
  proof: core/role/builtin/jaira-role-acceptance/build.py:150 (person line: tick of the open items + move), core/role/builtin/jaira-role-acceptance/build.py:61 (check() and plan refuse data without 'person'), core/role/builtin/jaira-role-acceptance/SKILL.md:160,220, jaira-role-signoff/SKILL.md:179; scratch: only the open n=3 is ticked, the ticked n=2 is not, missing key exits 1

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-21 20:59 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock
- **2026-11-21 20:59 · Alexander Sacharov** — Branch feat/3dvvfz was cut from the unpushed F1V59P tip (8ae150955); reset onto upstream/master (012caa9f7) so the PR carries only this ticket.
- **2026-11-21 21:00 · Alexander Sacharov** — decision: no path list. The built-in lanes ship into every repository that runs jaira, so core/mod/ or internal/tui/ would mean nothing elsewhere. pre-process recognises the kind of change instead: what a terminal draws (TUI, colours, widths, glyphs, a plugin/cockpit panel), behaviour that differs by OS/shell/terminal, and anything that needs a real person at a real keyboard. The item text starts with the fixed English prefix 'seen by a person:' because the prompts are English and later lanes match on it.
- **2026-11-21 21:00 · Alexander Sacharov** — decision: in-progress adds the item itself when no planning step ran (this very ticket had no 'planning' option); otherwise the most common path would never get one.
- **2026-11-21 21:01 · Alexander Sacharov** — Left .jaira/lanes/ (this repo's own board) alone: testing.md and review.md are byte-identical to the previous built-in and pre-process/in-progress carry local edits; bringing them along is 'jaira update --lanes', not this ticket.
- **2026-11-21 21:02 · Alexander Sacharov** — critique: the lanes leave 'seen by a person:' open, but nothing ticks it at acceptance — jaira-role-acceptance moves accepted tickets with --force, so the item stays [ ] into the logbook. Fix: acceptance ticks it from the passed step (proof = the person's mark); signoff names open such items beside DoD x/y.
- **2026-11-21 21:03 · Alexander Sacharov** — critique round 2: the criteria list appears in pre-process and again, shorter, in in-progress — left on purpose, each lane sees only its own prompt and in-progress is the path without planning.
- **2026-11-21 21:06 · Alexander Sacharov** — bughunt: change: blocker: core/role/builtin/jaira-role-acceptance/build.py:138-140 - class: the accept paths hand out 'jaira move <id> --to <terminal> --force', and --force overrides every gate refusal (core/move/move.go:126), the done lane's DoD check included, while nothing ties a 'seen by a person:' item's tick to that move; also in jaira-role-acceptance/SKILL.md:216-221 (ticking relies on 'the item's step', but data.json links steps to tickets, not to DoD items, and build.py check/state/plan never read the dod list) and jaira-role-signoff/SKILL.md:108-110+170 (an open item is reported as 'not missing work' beside an accept line with --force and no 'jaira dod <id> <n> --done' line); example: signoff ticket with dod ['code done','seen by a person: emoji column lines up in Windows Terminal'] whose review-check has no step for the second item - plan prints 'accepted - every step ok' and 'jaira move AAAAAA --to done --force', the ticket lands in done with the item open; fix direction: build.py refuses or flags acceptance while a dod entry starting 'seen by a person:' has no step naming it, and the signoff accept line carries the dod --done line(s) first; ran: build.py page+plan on a scratch data.json/db.json (page exit 0, plan exit 0 with the --force move)
- **2026-11-21 21:06 · Alexander Sacharov** — bughunt: done at 88246fa20 - 1 finding (1 change, 0 pre-existing) in 6 modules, first hunt; concurrency, error paths, input, permissions checked; manifest hashes match the working tree, go test ./core/lane/... ./core/role/... ./core/shipped/... ./core/release/... passes; files: core/lane/builtin/15-pre-process.md core/lane/builtin/20-in-progress.md core/lane/builtin/27-testing.md core/lane/builtin/40-review.md core/release/notes/3DVVFZ.md core/role/builtin/jaira-role-acceptance/SKILL.md core/role/builtin/jaira-role-signoff/SKILL.md core/shipped/shipped.sha256
- **2026-11-21 21:08 · Alexander Sacharov** — Fix for bughunt change finding: the agent no longer ticks 'seen by a person:' items at all (round-1 design had acceptance tick them from the passed step, but steps map to tickets, not DoD items, so build.py could not tell which step proved which item). Smaller design: build.py plan/signoff hand the person one line 'jaira dod <id> <n> --done --proof ... && jaira move ... --force'. Dead end: adding a step->item mapping to data.json — more schema for the same result.
- **2026-11-21 21:08 · Alexander Sacharov** — stop-after-fix turned off again: it was set for a prompts-only change, but the fix for the bughunt finding put code into build.py, which deserves its own hunt.
- **2026-11-21 21:10 · Alexander Sacharov** — bughunt: change: blocker: core/role/builtin/jaira-role-acceptance/build.py:145 - class: plan picks 'seen by a person:' items by text alone and ignores each item's state, so the line it hands the person ticks items that are already done or superseded; acceptance SKILL.md:160 also only passes the texts, so the state never reaches plan. example: round-2 acceptance of a ticket whose person item was ticked in round 1 with 'real: saw it in WT' prints 'jaira dod <id> 1 --done --proof <what I saw, where> && ... move --force' - running it replaces the person's real proof (data loss), and a [-] superseded item becomes [x] beside work nobody did (wrong result). Only open items (state todo/doing) belong in the line, as signoff SKILL.md:178 already says ('every open'). ran: build.py plan on scratch data with a done and a superseded person item (both printed), then jaira dod --done on a scratch board (exit 0, proof overwritten, [-] -> [x]); blame: 2871d8f71
- **2026-11-21 21:10 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-role-acceptance/build.py:145 - class: a missing input field read as 'nothing to do'; example: a data.json ticket without 'dod' (check() at build.py:52 does not require it) gets t.get('dod', []) = [] and plan prints the bare 'jaira move <id> --to <next> --force', i.e. exactly the move past an unticked 'seen by a person:' item this change exists to prevent, with no warning. check() should require 'dod' on every ticket like it requires next_lane/return_lane (build.py:58). read; blame: 2871d8f71
- **2026-11-21 21:10 · Alexander Sacharov** — bughunt: done at 2871d8f71 - 2 findings (2 change, 1 of them blocker, 0 pre-existing), follow-up hunt since 88246fa20: diff of 4 known files, 1 module read whole (build.py); concurrency, error paths, input, permissions checked; go test -count=1 ./core/role/ ./core/shipped/ passes; files: core/lane/builtin/15-pre-process.md core/lane/builtin/20-in-progress.md core/lane/builtin/27-testing.md core/lane/builtin/40-review.md core/release/notes/3DVVFZ.md core/role/builtin/jaira-role-acceptance/SKILL.md core/role/builtin/jaira-role-acceptance/build.py core/role/builtin/jaira-role-signoff/SKILL.md core/shipped/shipped.sha256
- **2026-11-21 21:12 · Alexander Sacharov** — Round 3 fix (bughunt 2871d8f71): build.py no longer parses dod texts — it could not see an item's state, so it re-ticked done and superseded items. The agent now passes 'person': the open seen-by-a-person items with their .dod_items[].n, read off jaira; the key is required (check() and plan exit 1 without it), which also closes the non-blocker 'missing dod field gives a bare --force move'. The __pycache__ the hunter saw came from my py_compile run; removed, never committed.
