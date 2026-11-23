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
updated-at: 2026-11-23T08:43:39Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-27805
claimed-at: 2026-11-23T08:25:28Z
outcome-what: "review lane writes note: edge:/limit:/elsewhere: beside blocker:; cockpit classOf reads that word for the Acceptance pane and the light; signoff/teamlead skills name the same four words"
outcome-why: "the cockpit showed every non-blocker note as yellow while the teamlead classed them by meaning, so both showed different colours for the same notes"
outcome-resolves: "teamlead, signoff table and cockpit show one class per note; limit/elsewhere notes no longer turn the light yellow"
review-summary: "none"
round: 3
back-from: testing
back-by: lane
back-count: 1
review-gaps: "none - PRE_EXISTING replaced by OLD_MARK (one caller set: classOf, gapItems); nothing duplicated or uncalled left; the old-word check in noteclass_test.go was dropped in-progress as prose-only"
test-verdict: "pass: go vet + go test ./... RC=0 at 5e9a470; after 911822f core/role, core/lane RC=0 and cockpit 'claude plugin test' 472/472; DoD 1-4 checked in the tree (noteclass_test, gapItems/light table tests, notes file, shipped.sha256); built binary: fresh init review lane asks for note: edge/limit/elsewhere, 'update --lanes --dry-run' on an older board updates review.md + testing.md; unknown command exits 2; change writes nothing, so no two-writer case"
hunt-verdict: "fail: bughunt change - followTitle strips any number of class words, so a finding starting 'Limit:' loses that word in the follow-up title"
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
