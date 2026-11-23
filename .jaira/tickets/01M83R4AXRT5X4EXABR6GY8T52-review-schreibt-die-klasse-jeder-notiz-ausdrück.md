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
updated-at: 2026-11-23T08:37:05Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-27805
claimed-at: 2026-11-23T08:25:28Z
outcome-what: "review lane writes note: edge:/limit:/elsewhere: beside blocker:; cockpit classOf reads that word for the Acceptance pane and the light; signoff/teamlead skills name the same four words"
outcome-why: "the cockpit showed every non-blocker note as yellow while the teamlead classed them by meaning, so both showed different colours for the same notes"
outcome-resolves: "teamlead, signoff table and cockpit show one class per note; limit/elsewhere notes no longer turn the light yellow"
review-summary: "core/mod/builtin/jaira-cockpit/hooks/board.ts: NOTE_CLASS and OLD_MARK both match pre-existing/fix-now, WORD_CLASS carries each word twice, and gapItems re-tests OLD_MARK to tell a firm class from an old mark. Simpler: NOTE_CLASS only for edge|limit|elsewhere (the firm words), OLD_MARK for the two old marks, key normalised by dropping '-'; gapItems takes the stated class when it is a blocker or NOTE_CLASS matched."
round: 2
back-from: critique
back-by: lane
back-count: 1
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
