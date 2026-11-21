---
id: 01M7ZYF17AVJAEWK40GG3DVVFZ
title: "Bei TUI-, Cockpit- oder OS-spezifischen Änderungen schreibt pre-process einen DoD-Punkt 'von einer Person angesehen'"
status: in-progress
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
updated-at: 2026-11-21T21:21:44Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3788188
claimed-at: 2026-11-21T21:19:40Z
stop-after-fix: "false"
outcome-what: "Built-in lanes: pre-process adds 'seen by a person: <what, where>' DoD items for terminal-drawn output, platform-dependent behaviour and person-only actions; in-progress adds one when no planning ran and never ticks it; testing leaves it open (agent-ticked = finding, untick); review makes it a review-check step (agent-ticked = blocker). Acceptance (build.py plan) and signoff hand the person one line 'jaira dod <id> <n> --done --proof ... && jaira move ... --force' instead of a bare --force move. shipped.sha256 + release note."
outcome-why: "02D92C, JDDCEZ, AX8K02 ticked human-only checks with stand-in proofs; and --force on the accept move passes an unticked DoD, so the item would vanish unseen."
outcome-resolves: "DoD 1-4 with proofs; categories: prompts plus one python print path — no concurrency or writes; input: dod texts with leading spaces handled (strip), a ticket without dod gets the plain move; permissions: the tick is never run by an agent."
review-summary: "Adds a definition-of-done item kind 'seen by a person:' for changes to what a terminal draws, platform-dependent behaviour, or person-only actions. pre-process now asks for one item per thing to look at; in-progress adds it itself when no planning ran and never ticks it; testing leaves it open and unticks it if an agent ticked it; review makes it a step of review-check and a blocker when agent-ticked. The acceptance role's data gains a required 'person' list (open items only); build.py check() refuses data without it, plan now runs check() and, for an accepted ticket with open person items, prints one line 'jaira dod <id> <n> --done --proof ... && jaira move --force' for the person to run instead of a bare --force move. The signoff role names open person items beside the DoD count and puts their tick in front of the accept line. Release note and shipped.sha256 updated."
round: 5
back-from: signoff
back-by: person
back-count: 1
review-gaps: |-
  note: build.py plan reads p['n'] and p['text'] without checking them — a person entry missing either crashes plan with KeyError instead of check() naming it; fix: check() validates each entry has an int n and a text
  note: this repo's own .jaira/lanes/ copies keep the old testing/review prompts until 'jaira update --lanes' is run here (decided in the notes, not a defect of the shipped builtins)
  note: build.py report does not show the person items, so the written acceptance record does not say which item the person still had to tick
test-verdict: "pass: round 3 over 2871d8f71..HEAD: go test ./core/role/... ./core/shipped/... RC=0, shipped-history --check RC=0, check-ticket-commits RC=0; DoD 1-4 opened at their proof lines; build.py on scratch data: without 'person' plan and page exit 1 naming each ticket, with it only the open item n=3 is in the person line and a ticket with [] gets the plain move; nothing written, so no two-writer case"
hunt-verdict: "pass: bughunt done at 7ccb7c724, 2 change non-blockers (malformed person element crashes plan loudly; an older acceptance set needs 'person' filled before plan/page) left as notes for review and the person, 0 pre-existing"
review-verdict: "All four DoD items are met by the diff (15-pre-process.md:77, 20-in-progress.md:164, 27-testing.md:57, 40-review.md:107, build.py:61/150, signoff SKILL.md:179, release note, shipped.sha256); go test ./core/... passes. Only prompt text and a small build.py change; no blocker, three notes."
review-round: 4
review-verdict-gaps: |-
  note: build.py plan reads p['n'] and p['text'] without checking them — a person entry missing either crashes plan with KeyError instead of check() naming it; fix: check() validates each entry has an int n and a text
  note: this repo's own .jaira/lanes/ copies keep the old testing/review prompts until 'jaira update --lanes' is run here (decided in the notes, not a defect of the shipped builtins)
  note: build.py report does not show the person items, so the written acceptance record does not say which item the person still had to tick
review-check: |-
  1. cd /home/alex/projects/.worktrees/jaira-3dvvfz
  2. Open core/lane/builtin/15-pre-process.md, section '## What only a person can check' (line 77) — it lists the three change kinds and the 'seen by a person:' prefix
  3. Write /tmp/a.json: {"blocks":[{"n":1,"tickets":["X"],"steps":[{"id":"s1","do":"d","expect":"e","covers":["X"]}]}],"tickets":[{"id":"X","title":"t","goal":"g","dod":[],"person":[{"n":3,"text":"seen by a person: emoji column in Windows Terminal"}],"round":1,"next_lane":"done","return_lane":"in-progress"}]}
  4. Write /tmp/db.json: {"results":{"s1":{"status":"ok","round":1}}}
  5. Run: python3 core/role/builtin/jaira-role-acceptance/build.py plan /tmp/a.json /tmp/db.json
  6. You should see '# X: for the person - seen by a person: ...' and one line 'jaira dod X 3 --done --proof ... && jaira move X --to done --force   # the person runs this line'
  7. Remove the "person" key from /tmp/a.json and rerun: it prints 'ticket X: no person list' and exits 1
---

# Bei TUI-, Cockpit- oder OS-spezifischen Änderungen schreibt pre-process einen DoD-Punkt 'von einer Person angesehen'

## Definition of Done

- [x] pre-process-Lane-Prompt verlangt den DoD-Punkt 'von einer Person angesehen' für die genannten Änderungsarten
  proof: core/lane/builtin/15-pre-process.md:77 '## What only a person can check'
- [x] in-progress-, testing- und review-Prompt sagen: diesen Punkt hakt nur eine Person ab
  proof: core/lane/builtin/20-in-progress.md:164, 27-testing.md:57, 40-review.md:107
- [x] shipped.sha256 und Release-Note core/release/notes/<ID>.md
  proof: core/shipped/shipped.sha256 (scripts/shipped-history.sh --write, 910 versions), core/release/notes/3DVVFZ.md
- [x] Abnahme: jaira-role-acceptance (build.py plan) und jaira-role-signoff geben für einen offenen 'seen by a person:'-Punkt eine Zeile 'dod --done && move --force' an die Person; kein Agent hakt ihn ab
  proof: core/role/builtin/jaira-role-acceptance/build.py:154 (person line: tick of the open items + move), build.py:61 (check() refuses a malformed person list in page, plan and report), SKILL.md:160,222, jaira-role-signoff/SKILL.md:179
- [x] build.py: fehlender person-Eintrag = leere Liste, falsch geformter = klare Fehlermeldung, kein KeyError (Test)
  proof: core/role/builtin/jaira-role-acceptance/build.py:61 (missing person = [], malformed = check() message); TestAcceptanceBuildPersonItems in core/role/acceptance_build_test.go
- [x] build.py report zeigt die 'seen by a person'-Punkte (Test)
  proof: core/role/builtin/jaira-role-acceptance/build.py:242 (report block lists open person items); TestAcceptanceBuildPersonItems/report_lists_open_person_items

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
- **2026-11-21 21:13 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-role-acceptance/build.py:60 (check) and :150-153 (plan) - class: check() tests only that tickets[].person is a list, never the shape of its items, so a malformed item passes 'page' and crashes 'plan' at read-back, after the person has already marked the page; example: person: ["seen by a person: it renders"] (the shape of the dod list right beside it) -> TypeError at :152, person: [{"n": 1}] -> KeyError 'text'; both exit 1 with a traceback, nothing written, db.json intact. Not a blocker in my judgement: the input breaks the documented shape, the failure is loud and loses nothing; fix is to check each item is {n: int, text: str} in check(); ran: build.py page/plan on a scratch data.json (page exit 0, plan exit 1 traceback; well-formed item exit 0)
- **2026-11-21 21:13 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-role-acceptance/build.py:60 and :135 - class: an acceptance set whose data.json predates the person key (written by the released build.py, kept in ~/.claude/acceptance/<repo>/<name>/ across rounds) is now refused at read-back (plan) and at every rebuild while extending or starting the next round, until each old ticket gets a person list; example: data.json without person -> 'ticket X: no person list ...', exit 1. Loud refusal, nothing lost, the message names what to add; SKILL.md 'If url.txt exists, you are extending' and '3. Next round' do not say that old tickets need person filled in, so an agent may write [] blindly and hand out the bare --force move the key exists to prevent. Edge, not blocker; ran: build.py page on data.json without person (exit 1); read: SKILL.md:44-57, 280-288
- **2026-11-21 21:14 · Alexander Sacharov** — bughunt: done at 7ccb7c724 - 2 findings (2 change, 0 blocker, 0 pre-existing), follow-up hunt since 2871d8f71: diff of 3 known files (build.py, acceptance SKILL.md, shipped.sha256), no new file; concurrency, error paths, input, permissions checked; go test -count=1 ./core/role/ ./core/shipped/ passes; files: core/lane/builtin/15-pre-process.md core/lane/builtin/20-in-progress.md core/lane/builtin/27-testing.md core/lane/builtin/40-review.md core/release/notes/3DVVFZ.md core/role/builtin/jaira-role-acceptance/SKILL.md core/role/builtin/jaira-role-acceptance/build.py core/role/builtin/jaira-role-signoff/SKILL.md core/shipped/shipped.sha256
- **2026-11-21 21:18 · Alexander Sacharov** — Given back from signoff (Alex 21.11.): (1) build.py: fehlender oder falsch geformter 'person'-Eintrag darf nicht mit KeyError abstürzen — fehlend = leere Liste (deckt auch Abnahme-Sets von vor dieser Version), falsch geformt = klare Fehlermeldung in check(). Ein Test für beide. (2) build.py report listet die 'seen by a person'-Punkte mit, sonst sieht die Person sie im Bericht nicht — genau der Zweck des Tickets. Note 'jaira update --lanes' bleibt known limit.
- **2026-11-21 21:21 · Alexander Sacharov** — Give-back round: check() now also runs in report, so a malformed person list never reaches the report loop. Missing key = [] means a pre-key set gets the bare --force move again — Alex's call (21.11.), SKILL.md tells the agent to fill person in when extending such a set.
