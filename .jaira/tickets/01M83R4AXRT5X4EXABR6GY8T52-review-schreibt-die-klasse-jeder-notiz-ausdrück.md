---
id: 01M83R4AXRT5X4EXABR6GY8T52
title: "Review schreibt die Klasse jeder Notiz ausdrücklich, und Cockpit, Abnahme und Teamlead lesen dieselbe Klasse"
status: critique
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
commits: []
created-at: 2026-11-23T08:25:01Z
updated-at: 2026-11-23T08:51:41Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-27805
claimed-at: 2026-11-23T08:25:28Z
outcome-what: "review lane writes note: edge:/limit:/elsewhere: beside blocker:; cockpit classOf reads that word for the Acceptance pane and the light; signoff/teamlead skills name the same four words"
outcome-why: "the cockpit showed every non-blocker note as yellow while the teamlead classed them by meaning, so both showed different colours for the same notes"
outcome-resolves: "teamlead, signoff table and cockpit show one class per note; limit/elsewhere notes no longer turn the light yellow"
review-summary: "Full review of 5e9a470..ad97b24 (first review). The review lane (40-review.md) now asks for one of four classes per finding: 'blocker:', 'note: edge:', 'note: limit:', 'note: elsewhere:' (fix-now becomes 'note: edge: fix-now:'); testing writes 'testing: note: edge:'. In the cockpit, classOf() reads the class a finding states; gapItems uses it ahead of the Explain note's positional classes (older 'pre-existing:'/'fix-now:' marks and unmarked notes still follow Explain, then mark, else edge), and lightOf on the band counts only edge notes as yellow, so limit/elsewhere no longer turn the light yellow. followTitle strips one class word and one older mark. Signoff and teamlead SKILLs use the same four words and say a stated word is not to be re-marked. noteclass_test.go pins the words across lane and roles; board.test.ts tables cover classOf, gapItems, lightOf, followTitle. Release note and shipped.sha256 added."
round: 4
back-from: review
back-by: lane
back-count: 1
review-gaps: |-
  note: edge: fix-now: board.ts:1036 rephrasePrompt tells the model to drop a leading "blocker:", "note:" or "pre-existing:" but not the new "edge:", "limit:", "elsewhere:" - a rephrased note in the Acceptance pane can keep "edge:" in front of a mark that already says it; add the three words to that line
  note: edge: fix-now: texts.ts:228 and jaira-dispatcher/scripts/spawn.sh:597 (Explain prompt) still let the model class every note itself and do not say that a line already starting with blocker: / note: edge|limit|elsewhere: keeps that class - the pane now shows the stated class while the Explain text beside it can name another; add one sentence to both prompts
  note: elsewhere: go test ./role fails locally with TestEveryInstalledFileIsAShippedVersion because the untracked core/mod/builtin/jaira-cockpit/node_modules/.vite/... results.json is embedded as an installed file - not from this change, node_modules is not gitignored; clean CI checkout is unaffected
test-verdict: "pass: round 2 - cockpit 'claude plugin test' 472/472 RC=0 and go test core/role core/lane core/shipped RC=0 after ad97b249c; followTitle cases (Limit:/Edge:/pre-existing) in tests/board.test.ts; DoD 1-4 unchanged and verified in round 1; change writes nothing"
hunt-verdict: "pass: bughunt done at ad97b249c, 0 change, 0 pre-existing (round 1's followTitle finding fixed in ad97b249c)"
review-verdict: "accept: all four DoD items hold in the diff; no blocker. Two fix-now edges (rephrase and Explain prompts do not know the new class words) go back once."
review-round: 3
review-verdict-gaps: |-
  note: edge: fix-now: board.ts:1036 rephrasePrompt tells the model to drop a leading "blocker:", "note:" or "pre-existing:" but not the new "edge:", "limit:", "elsewhere:" - a rephrased note in the Acceptance pane can keep "edge:" in front of a mark that already says it; add the three words to that line
  note: edge: fix-now: texts.ts:228 and jaira-dispatcher/scripts/spawn.sh:597 (Explain prompt) still let the model class every note itself and do not say that a line already starting with blocker: / note: edge|limit|elsewhere: keeps that class - the pane now shows the stated class while the Explain text beside it can name another; add one sentence to both prompts
  note: elsewhere: go test ./role fails locally with TestEveryInstalledFileIsAShippedVersion because the untracked core/mod/builtin/jaira-cockpit/node_modules/.vite/... results.json is embedded as an installed file - not from this change, node_modules is not gitignored; clean CI checkout is unaffected
review-check: "1. cd core && go test ./role -run TestReviewSignoffAndTeamleadShareTheFourNoteClasses ./lane ./shipped - all ok  2. cd core/mod/builtin/jaira-cockpit && claude plugin test - 472 pass, including 'gapItems reads the class each finding states'  3. open core/lane/builtin/40-review.md, section 2 - you see blocker / note: edge / note: limit / note: elsewhere  4. in the cockpit, a signoff ticket whose review-gaps holds only 'note: limit: ...' shows a green light and a ⚪ row in the Acceptance pane"
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
