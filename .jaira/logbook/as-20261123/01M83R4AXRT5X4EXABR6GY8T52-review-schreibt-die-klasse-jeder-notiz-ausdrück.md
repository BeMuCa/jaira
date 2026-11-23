---
id: 01M83R4AXRT5X4EXABR6GY8T52
title: "Review schreibt die Klasse jeder Notiz ausdrücklich, und Cockpit, Abnahme und Teamlead lesen dieselbe Klasse"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Jede review-gaps-Zeile trägt ihre Klasse (blocker / edge / limit / elsewhere) als festes Präfix; das Cockpit (Verdict, Abnahme, Gesamtampel) liest nur dieses Präfix, sodass Teamlead und Cockpit dieselbe Farbe zeigen."
context: |-
  Was falsch ist (Alex 23.11.2026, Screenshot Abnahme VHGWQF): der Teamlead nannte die 4 Notizen ⚪/⚪/⚪/🔵, das Abnahme-Fenster zeigte alle 🟡 und die Gesamtampel gelb.
  Ursache: core/mod/builtin/jaira-cockpit/hooks/board.ts gapItems kennt nur BLOCKER und PRE_EXISTING per Regex; alles andere wird 'edge' (🟡). Wörter wie 'known limit', 'not this ticket', 'text only' im Fliesstext werden nicht erkannt; der Teamlead klassifiziert dagegen nach Sinn.
  Lösung: review-Lane (core/lane/builtin/40-review.md, auch critique/testing-Notizen, die in signoff landen) schreibt 'note: blocker:', 'note: edge:', 'note: limit:' oder 'note: elsewhere:' (pre-existing/not this ticket). gapItems liest nur das Präfix; ohne Präfix wie heute (edge), damit alte Tickets weiter gehen.
  Betroffen: AJ41SC-Ampel, 0GMYQB-Marker, 2V50Y5-Abnahme, jaira-role-signoff/teamlead SKILL (gleiche Wörter).
definition-of-done: "40-review.md (und jede Lane, die review-gaps/Notizen für signoff schreibt) verlangt pro Notiz genau ein Klassen-Präfix blocker|edge|limit|elsewhere (Rollen-/Lane-Test)"
tags:
  - cockpit
blocked-by: []
related: []
commits:
  - 5e9a470fe6d8c9c9b8e614c52fce9d4f19df8f8b
  - 911822f3c94787ac28b55f175d4589f134ddfee9
  - ad97b249c01d90ad524b6ec40a0ebf8e42dd4cf5
  - dd448069e77ea6f1abf818d06f38bbd0eeec414d
created-at: 2026-11-23T08:25:01Z
updated-at: 2026-11-23T09:04:26Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-27805
claimed-at: 2026-11-23T08:25:28Z
outcome-what: "review lane writes note: edge:/limit:/elsewhere: beside blocker:; cockpit classOf reads that word for the Acceptance pane and the light; signoff/teamlead skills name the same four words"
outcome-why: "the cockpit showed every non-blocker note as yellow while the teamlead classed them by meaning, so both showed different colours for the same notes"
outcome-resolves: "teamlead, signoff table and cockpit show one class per note; limit/elsewhere notes no longer turn the light yellow"
review-summary: "Delta review from ad97b249c to dd448069e: rephrasePrompt (board.ts) now drops a leading edge:/limit:/elsewhere:/fix-now: as well as blocker:/note:/pre-existing:, and both Explain prompts (cockpit texts.ts, dispatcher spawn.sh) tell the model that a note stating its class keeps it (blocker: -> [blocker], note: edge: -> [small edge], note: limit: -> [known limit], note: elsewhere: -> [not this ticket]); shipped.sha256 carries the new hashes of all three files. Both fix-now findings of the previous review are answered; nothing else changed in code."
round: 4
back-from: review
back-by: lane
back-count: 1
review-gaps: |-
  note: edge: no test pins the new class sentence in the Explain prompts or the new words in rephrasePrompt - a tests/board.test.ts expect(t.…explain…).toContain('note: edge: [small edge]') would keep the cockpit and spawn.sh wording from drifting apart
  note: elsewhere: AGENTS.md, CLAUDE.md and .jaira/lanes/{critique,in-progress,review,testing}.md are modified and uncommitted in this worktree - the 0.6.2 jaira block/lane refresh, not this ticket; leave them out of its commits
test-verdict: "pass: delta dd448069e - cockpit 'claude plugin test' 472/472 RC=0, go test core/role core/shipped RC=0, bash -n spawn.sh ok; prompt text only, writes nothing"
hunt-verdict: "pass: stop-after-fix, no hunt in round 4"
review-verdict: "accept: delta from ad97b249c answers both fix-now findings; hashes in shipped.sha256 match HEAD; go test ./shipped ./role ./lane ok; no blocker"
review-round: 4
review-verdict-gaps: |-
  note: edge: no test pins the new class sentence in the Explain prompts or the new words in rephrasePrompt - a tests/board.test.ts expect(t.…explain…).toContain('note: edge: [small edge]') would keep the cockpit and spawn.sh wording from drifting apart
  note: elsewhere: AGENTS.md, CLAUDE.md and .jaira/lanes/{critique,in-progress,review,testing}.md are modified and uncommitted in this worktree - the 0.6.2 jaira block/lane refresh, not this ticket; leave them out of its commits
review-check: "1. cd core && go test ./role -run TestReviewSignoffAndTeamleadShareTheFourNoteClasses -> ok  2. go test ./shipped -> ok (every shipped file hash is listed)  3. in the cockpit, open a signoff ticket whose review-gaps has a line 'note: limit: ...' -> its light is white (⚪), 'note: elsewhere: ...' -> blue (🔵), 'note: edge: ...' -> yellow (🟡)  4. press Explain on it -> the explain note puts that line under [known limit] / [not this ticket] / [small edge]"
stop-after-fix: "true"
---

# Review schreibt die Klasse jeder Notiz ausdrücklich, und Cockpit, Abnahme und Teamlead lesen dieselbe Klasse

## Definition of Done

- [x] 40-review.md (und jede Lane, die review-gaps/Notizen für signoff schreibt) verlangt pro Notiz genau ein Klassen-Präfix blocker|edge|limit|elsewhere (Rollen-/Lane-Test)
  proof: core/lane/builtin/40-review.md review-gaps item 2; 27-testing.md 'testing: note: edge:'; core/role/noteclass_test.go TestReviewSignoffAndTeamleadShareTheFourNoteClasses
- [x] gapItems/Gesamtampel lesen das Präfix: limit -> ⚪, elsewhere -> 🔵, edge -> 🟡, blocker -> 🔴; 'pre-existing:' zählt als elsewhere; ohne Präfix bleibt edge (Tabellentest)
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts classOf + gapItems; tests/board.test.ts 'gapItems reads the class each finding states', 'a note colours the light only when its class is edge'
- [x] jaira-role-signoff und jaira-teamlead SKILL verwenden dieselben vier Wörter
  proof: core/role/builtin/jaira-role-signoff/SKILL.md §2 table; jaira-teamlead/SKILL.md decision shape item 2; TestReviewSignoffAndTeamleadShareTheFourNoteClasses
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/GY8T52.md; core/shipped/shipped.sha256 (TestEveryInstalledFileIsAShippedVersion, TestEveryBuiltinLaneIsAShippedVersion)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 08:25 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/6958.sock
- **2026-11-23 08:26 · Alexander Sacharov** — decision (worker, nothing open for the person): format: 'blocker:' bleibt wie heute (Move-Gate und BLOCKER-Regex lesen es); Notizen schreiben 'note: edge:', 'note: limit:', 'note: elsewhere:'. 'note: pre-existing:' und 'note: fix-now:' bleiben gültig und zählen als elsewhere bzw. edge; fix-now kombiniert als 'note: edge: fix-now:'. Ein ausdrückliches Präfix schlägt die Explain-Zuordnung nach Position; ohne Präfix wie bisher (Explain, sonst pre-existing->elsewhere, sonst edge).
- **2026-11-23 08:33 · Alexander Sacharov** — in-progress: .jaira/lanes/*.md, AGENTS.md und CLAUDE.md waren nach dem Claim im Worktree verändert (jaira-Block 0.6.2 / Lanes aktualisiert, nicht von mir) - gehören nicht in diesen Commit, bleiben liegen. Explain-Notiz überstimmt eine ausdrückliche Klasse (edge/limit/elsewhere) nicht mehr; pre-existing:/fix-now: und Notizen ohne Klasse folgen ihr weiter wie bisher (Test 'a stated class outranks the Explain note').
- **2026-11-23 08:33 · Alexander Sacharov** — waiting: go test ./... (background shell b64230bz3) since 09:33
- **2026-11-23 08:35 · Alexander Sacharov** — critique: board.ts classOf - split the regexes by meaning (NOTE_CLASS = firm edge|limit|elsewhere, OLD_MARK = pre-existing|fix-now), one WORD_CLASS entry per word, gapItems decides firm by NOTE_CLASS instead of a negated OLD_MARK test
- **2026-11-23 08:37 · Alexander Sacharov** — critique: round 2 - the split of NOTE_CLASS/OLD_MARK answers round 1; nothing else raised. Fix 911822f3c under 50 lines, so optimize ran in the same pass.
- **2026-11-23 08:37 · Alexander Sacharov** — waiting: bughunt (background agent) since 09:37
- **2026-11-23 08:41 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/board.ts:937 followTitle - class: the class-word strip is repeated ((...)+), so a finding whose own text starts with edge/limit/elsewhere/pre-existing/fix-now followed by ':' loses that word in the follow-up ticket title; before this change only one old mark was stripped; example: 'note: edge: Limit: 0 means unlimited in jaira list' becomes title '0 means unlimited in jaira list' (likewise 'note: limit: Edge: lanes with one ticket' -> 'lanes with one ticket'); fix: strip one class word and at most one older mark after it, not any number; not a blocker - title only, --context keeps the whole finding; ran: node --experimental-strip-types on a scratch script calling followTitle
- **2026-11-23 08:41 · Alexander Sacharov** — hunt report: 1 finding, 1 change (edge, not a blocker), as a note on GY8T52 | read whole: board.ts (classOf, lightOf, gapItems, followTitle and their callers in register.tsx), noteclass_test.go, 27-testing.md, 40-review.md, the signoff/teamlead/worker SKILL.md, the release note, shipped.sha256, the board's uncommitted .jaira/lanes/* (the installed 0.6.2 build's text - jaira update --lanes after the release replaces it, so not a finding); ran claude plugin test (472 pass) and go test ./role ./lane ./shipped (ok) | not checked: how the cheap rephrase model treats the new 'edge:'/'limit:' word after 'note:' (rephrasePrompt only names blocker/note/pre-existing) - that is model behaviour and cannot be confirmed from the code
- **2026-11-23 08:41 · Alexander Sacharov** — bughunt: done at 911822f3c - 1 finding (1 change, 0 pre-existing) in 10 modules, first hunt; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/lane/builtin/27-testing.md core/lane/builtin/40-review.md core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/GY8T52.md core/role/builtin/jaira-role-signoff/SKILL.md core/role/builtin/jaira-teamlead/SKILL.md core/role/builtin/jaira-worker/SKILL.md core/role/noteclass_test.go core/shipped/shipped.sha256
- **2026-11-23 08:41 · Alexander Sacharov** — testing: back to in-progress for the hunt's change finding (followTitle, board.ts:937): strip one class word and at most one older mark after it
- **2026-11-23 08:43 · Alexander Sacharov** — critique: round 3 (after testing) - ad97b249c answers the hunt's followTitle finding; nothing else raised; under 50 lines, optimize in the same pass: nothing to remove
- **2026-11-23 08:44 · Alexander Sacharov** — waiting: bughunt round 2 (background agent) since 09:44
- **2026-11-23 08:45 · Alexander Sacharov** — hunt report: 0 findings, nothing written | read as a diff since 911822f3c: core/mod/builtin/jaira-cockpit/hooks/board.ts (followTitle), tests/board.test.ts, core/shipped/shipped.sha256 (hash = HEAD board.ts), uncommitted .jaira/lanes/{critique,in-progress,review,testing}.md, AGENTS.md, CLAUDE.md; no new file | could not run the vitest suite here (no runner config in the worktree); followTitle checked instead by a scratch node run over 14 inputs (class+mark, reversed case, markdown emphasis, empty rest, 'edges:'/'limited:'), all as intended
- **2026-11-23 08:45 · Alexander Sacharov** — bughunt: done at ad97b249c - 0 findings (0 change, 0 pre-existing) in 10 modules, follow-up hunt since 911822f3c; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/lane/builtin/27-testing.md core/lane/builtin/40-review.md core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/GY8T52.md core/role/builtin/jaira-role-signoff/SKILL.md core/role/builtin/jaira-teamlead/SKILL.md core/role/builtin/jaira-worker/SKILL.md core/role/noteclass_test.go core/shipped/shipped.sha256
- **2026-11-23 08:45 · Alexander Sacharov** — waiting: review lane (fresh agent via run-lane.sh, Herdr tab) since 09:45
- **2026-11-23 08:49 · Alexander Sacharov** — review: done at ad97b249c
- **2026-11-23 08:49 · Alexander Sacharov** — review: fix-now: (1) board.ts rephrasePrompt: add edge:/limit:/elsewhere: to the 'Drop a leading' line; (2) texts.ts:228 + jaira-dispatcher/scripts/spawn.sh:597 Explain prompt: a note that already starts with blocker:/note: edge|limit|elsewhere: keeps that class - one sentence each; update shipped.sha256 for changed shipped files
- **2026-11-23 08:49 · Alexander Sacharov** — worker: stop-after-fix=true set for the fix-now round (two prompt sentences, no lock, no parser). Removed core/mod/builtin/jaira-cockpit/node_modules - left untracked by my own 'claude plugin test' runs, the cause of the review's elsewhere note.
- **2026-11-23 08:51 · Alexander Sacharov** — critique: fix-now round - dd448069e adds the class words to rephrasePrompt and one sentence to both Explain prompts (texts.ts, spawn.sh), the two review fix-now notes; nothing else; optimize in the same pass: nothing to remove
- **2026-11-23 08:52 · Alexander Sacharov** — waiting: second review (new fresh agent via run-lane.sh) since 09:52
- **2026-11-23 08:53 · Alexander Sacharov** — review: done at dd448069e
- **2026-11-23 08:54 · Alexander Sacharov** — report: review notes carry blocker:/note: edge|limit|elsewhere:; cockpit classOf colours pane+light by that word; signoff/teamlead/Explain/rephrase use the same four words | signoff, context 184k of 1000k, cache warm 60 min | open for the person: 2 review notes (edge: no test pins the Explain/rephrase wording; elsewhere: uncommitted 0.6.2 jaira-block/lane refresh in worktree); branch pushed, PR not opened
- **2026-11-23 08:56 · Alexander Sacharov** — hunt report: 0 findings (0 change), nothing written as bughunt notes | read as diff since ad97b249c: core/mod/builtin/jaira-cockpit/hooks/board.ts (rephrasePrompt, checked against classOf/gapItems/followTitle), core/shipped/shipped.sha256 (new hashes match the three files at HEAD); read whole as new to the list: core/mod/builtin/jaira-cockpit/hooks/texts.ts, core/role/builtin/jaira-dispatcher/scripts/spawn.sh (change is prompt text only; bash -n ok); ran: claude plugin test . (472 pass, 0 fail) | not checked: how a model actually follows the reworded explain/rephrase prompts - needs a live run
- **2026-11-23 08:56 · Alexander Sacharov** — bughunt: done at dd448069e - 0 findings (0 change, 0 pre-existing) in 12 modules, follow-up hunt since ad97b249c; concurrency, error paths, input, permissions checked; files: core/lane/builtin/27-testing.md core/lane/builtin/40-review.md core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/GY8T52.md core/role/builtin/jaira-dispatcher/scripts/spawn.sh core/role/builtin/jaira-role-signoff/SKILL.md core/role/builtin/jaira-teamlead/SKILL.md core/role/builtin/jaira-worker/SKILL.md core/role/noteclass_test.go core/shipped/shipped.sha256
