---
id: 01M82ZBBANXWEKXWDTM1ZZTCWS
title: "Die Abnahme-Ansicht zeigt oben, wie man prueft, was gemacht wurde, die offene Frage, den Code-Diff und die DoD - und sonst nichts"
status: testing
ready: true
creator: BeMuCa
goal: "Oeffnet eine Person ein Ticket in human, review oder signoff in der TUI, sieht sie oben in dieser Reihenfolge: So pruefst du (review-check), Was gemacht wurde (outcome what/why), Offene Frage (question, nur wenn eine da ist), die Code-Aenderungen und die DoD mit Proofs. Felder, die der Person bei der Entscheidung nicht helfen, stehen nicht mehr auf dieser Ansicht."
context: |-
  Was heute schief ist: internal/tui/signoff.go renderSignOff (ab Zeile 22) zeigt Problem, What/Why/Resolves, Summary, Gaps, Verdict und als LETZTE Zeile 'check' (review-check, signoff.go:83). Wie man prueft, steht also ganz unten; davor viel Text, der beim Entscheiden nicht hilft.
  Ausloeser: Berk 2026-11-23 - 'ja so soll es sein' zu: oben So pruefst du / Was gemacht wurde / Offene Frage, dazu Code-Aenderungen und DoD; 'rest muss nicht angezeigt werden'; 'da ist viel sinnloses sonst was dem user nix bringt'.
  Code-Aenderungen: EH6Z0N (berkc, signoff) baut den Diff gegen master je Datei mit Sprung nach VS Code (internal/tui/diffpane.go, Commit 2a73eae23). Dieses Ticket ordnet ihn nur ein und baut ihn nicht neu.
  Die Felder muessen in der richtigen Lane gefuellt werden: review-check schreibt die review-Lane (TYTVBZ), outcome die in-progress-Lane, question jede Lane, die nach human schickt. Fehlt ein Feld, zeigt die Ansicht 'fehlt - <Lane> haette es schreiben muessen' statt einer leeren Stelle.
  Zu klaeren im Brainstorm: welche Felder genau wegfallen (summary, gaps, verdict, problem?) und ob sie hinter einer Taste erreichbar bleiben.
  Zusammenhang: Q63EEV (Cockpit-Pane je Ticket, berkc), 2V50Y5 (Abnahme-Knopf im Cockpit), 5QR36R (Abnahme im Chat).
definition-of-done: "Die Ansicht in human, review und signoff zeigt oben in dieser Reihenfolge: So pruefst du, Was gemacht wurde, Offene Frage (wenn vorhanden), Code-Aenderungen, DoD mit Proofs"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-11-23T01:11:56Z
updated-at: 2026-11-23T21:19:26Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-79231
claimed-at: 2026-11-23T20:54:57Z
outcome-what: "Sign-off screen (signoff and now also human) shows blocker lines from review-gaps, check, what, why, question, commits, DoD with proofs; problem, resolves, summary, verdict, remaining gaps, tags and board lane fields moved behind 'i' (details block, per ticket via Model.detailsFor). atPersonsLane() is the one predicate for the screen and EH6Z0N's diff pane; accept/follow-up stay signoff-only."
outcome-why: The check sat last under a page of account that does not help the decision (Berk 2026-11-23); a ticket in human opened to the generic detail dump.
outcome-resolves: "DoD 1-5; tests TestSignOffViewShowsTheDecisionFirst, TestSignOffKeepsBlockersOnTop, TestHumanLaneOpensTheDecisionScreen; checked: concurrency n/a (render only), error paths n/a, unusual input - gaps with list markers/case/indent and only-blocker gaps, empty fields owed rows; permissions n/a"
review-summary: none
review-gaps: "duplication: searched core/ and internal/ for a review-gaps prefix parser (blocker:/fix-now:/note:) - none, gapLines stays; atPersonsLane removed the copy of the lane test in diffPaneShown. dead code: none orphaned - firstNonEmpty still used 3 times, laneFields/owedRow/fieldRow still called from the details block. fluff: fixed view.go:134 comment that still said 'human checkpoint' for the wider predicate; left the banner/hint locals. cost: render-only, gapLines splits one short string per render - left."
---

# Die Abnahme-Ansicht zeigt oben, wie man prueft, was gemacht wurde, die offene Frage, den Code-Diff und die DoD - und sonst nichts

## Definition of Done

- [x] Die Ansicht in human, review und signoff zeigt oben in dieser Reihenfolge: So pruefst du, Was gemacht wurde, Offene Frage (wenn vorhanden), Code-Aenderungen, DoD mit Proofs
  proof: internal/tui/signoff.go:81-131 signOffBody order; view.go render() via atPersonsLane; tests TestSignOffViewShowsTheDecisionFirst, TestHumanLaneOpensTheDecisionScreen
- [x] Felder ohne Nutzen fuer die Entscheidung sind von dieser Ansicht entfernt; welche, steht nach dem Brainstorm als Notiz am Ticket
  proof: decision note 'decisions settled (worker...)' lists them; signoff.go:136-160 details block behind detailsFor; TestSignOffViewShowsTheDecisionFirst hidden list
- [x] Ein leeres Pflichtfeld zeigt, welche Lane es haette schreiben muessen
  proof: signoff.go declared() -> owedRow '— owed by <lane>' for check/what/why; TestSignOffShowsOwedFieldsNamingTheLane
- [x] Release-Notiz in core/release/notes/<ID>.md
  proof: core/release/notes/ZZTCWS.md
- [x] Eine 'blocker:'-Zeile aus review-gaps bleibt oben sichtbar; problem, resolves, summary, verdict und die uebrigen gaps sind hinter einer Taste erreichbar; Test
  proof: signoff.go:81 gapLines/isBlockerLine; model.go case "i" toggles detailsFor; tests TestSignOffKeepsBlockersOnTop, TestSignOffViewShowsTheDecisionFirst

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] renderSignOff on person lanes (human+signoff); a/f stay signoff-only
- [x] signOffBody reordered: blocker, check, what, why, question, commits, DoD
- [x] 'i' toggles a details block: problem/resolves/summary/verdict/gaps/tags/lane fields
- [x] update existing signoff tests, add order/blocker/details tests
- [x] release note core/release/notes/ZZTCWS.md

## Progress
- **2026-11-23 20:54 · Alexander Sacharov** — Alexander Sacharov took this ticket over from BeMuCa
- **2026-11-23 20:54 · Alexander Sacharov** — decision: Alex 23.11.2026: vom Teamlead uebernommen, sofort umsetzen. Gestapelt auf EH6Z0N (berkc, signoff, Branch upstream/feat/EH6Z0N-diff-view, noch nicht auf master): beide aendern internal/tui/signoff.go, der Diff-Teil kommt von dort. Der Branch feat/zztcws enthaelt deshalb die EH6Z0N-Commits; gemergt wird nach EH6Z0N. Brainstorm-Vorgabe vom Teamlead: problem, resolves, summary und verdict fallen von der Ansicht und bleiben hinter einer Taste (z.B. 'i' Details) erreichbar; eine 'blocker:'-Zeile aus review-gaps bleibt oben sichtbar, weil sie die Entscheidung aendert; die uebrigen gaps hinter derselben Taste. Die genaue Liste schreibt die Brainstorm-Lane als Notiz (DoD 2).
- **2026-11-23 20:55 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/48662.sock
- **2026-11-23 20:59 · Alexander Sacharov** — stacked on EH6Z0N (merge commit f03f8598b): signoff.go conflict resolved keeping EH6Z0N's diff layout; master's 'm send back' hint carried into renderSignOff's hint (and diffpane_test.go's narrow-screen want)
- **2026-11-23 21:00 · Alexander Sacharov** — decisions settled (worker, no person needed): (1) 'human, review und signoff' = the lanes a person owns, human and signoff, the same predicate EH6Z0N's diffPaneShown uses (RequiresHumanExit||RequiresQuestion); the agent lane 'review' keeps the plain detail view - EH6Z0N's title 'human review und signoff' names the human lane, not the review lane. Today a ticket in human shows renderDetail; it now gets this view too, but a/f (accept, follow-up) stay signoff-only. (2) Top of the view: blocker lines from review-gaps, check, what, why, question (only when set), Commits/diff, DoD with proofs. (3) DoD 2 list - moved behind 'i' details: problem (goal/context), resolves, summary, verdict, all of review-gaps, tags, and the board's own lane fields (laneFields). Labels stay English like the rest of the TUI. (4) An empty check/what/why keeps its row as '— owed by <lane>' (existing owedRow wording, same on every screen).
- **2026-11-23 21:09 · Alexander Sacharov** — waiting: full go vet + go test ./internal/tui ./internal/cli ./core/... (bg task bqbdbt649) since 22:09
- **2026-11-23 21:14 · Alexander Sacharov** — critique (pass 1, whole diff f03f8598b..0cfce40a2): simpler shape - none; the order lives in one function and atPersonsLane replaces the copy of the lane test that was in diffPaneShown. Fit - gapLines is the first parser of review-gaps prefixes in Go (grep for blocker:/fix-now:/note: in core and internal found none), so no existing helper to reuse. Speculative - none. Let stand: hint for human lane reuses detailHints (same keys as the detail view there, a/f stay signoff-only); the gaps row behind 'i' is skipped when review-gaps holds only blocker lines (they are on top already). review-summary=none.
- **2026-11-23 21:17 · Alexander Sacharov** — testing: function pass on a scratch board (jaira board, tmux 120x45, ticket in signoff with review-gaps 'blocker: resets on restart / note: wording'): screen shows blocker, check, what, why, DoD; 'i' adds Details with problem, resolves, summary, gaps (only the note line), verdict. testing: note: edge: the blocker row reads 'blocker      blocker: resets on restart' - label and line prefix repeat the word; fix: strip the prefix in the top row. waiting: race suite (bg) and bug hunter (subagent) since 22:17
- **2026-11-23 21:19 · Alexander Sacharov** — bughunt: change: blocker: internal/tui/signoff.go:189-192 (isBlockerLine) - class: a review-gaps blocker line written in any form other than a bare/dashed/starred 'blocker:' prefix is not recognised, so the sign-off screen moves it behind 'i' as an ordinary gap and the person can press 'a' without ever seeing it (before this change every gap was on the screen). Forms missed: numbered list '1. blocker: x', '**Blocker**: x', '`blocker`: x', 'blocker : x' - all of which the cockpit's BLOCKER/MARKER regexes in core/mod/builtin/jaira-cockpit/hooks/board.ts:340-345 count as blockers, and DF0YD6's release note promises; example: review-gaps='1. blocker: DoD 2 unmet' renders no blocker row; ran: scratch test calling isBlockerLine on 8 forms (false for the 4 above, true for 'blocker:', '- blocker:', '**blocker:**', tab+'BLOCKER:')
- **2026-11-23 21:19 · Alexander Sacharov** — bughunt: change: internal/tui/view.go:136 (render routes atPersonsLane to renderSignOff) - class: everything detailBody showed for a ticket in the human lane that signOffBody does not carry is now unreachable in the TUI, not even behind 'i': the Progress notes (body, where the question's background is usually written), context when a goal is set, plan, blocked-by / waiting-on, follows, assignee; only 'E' (editor) still shows them. Example: a human-lane ticket with goal, context and a note 'option A vs B' - the person asked the question sees neither the context nor the note; ran: scratch render test, Status human: note=false context=false with and without 'i'; same ticket in critique: note=true context=true. Not a blocker under the rule (no wrong result in data), but the 'i' details list was designed for signoff tickets and now also serves the lane where a person must answer from that background
- **2026-11-23 21:19 · Alexander Sacharov** — bughunt: pre-existing: internal/tui/diffpane.go:202-210 (splitDiffGit) and :106-112 ('+++' path) - class: a path git prints quoted (non-ASCII, or containing a quote, backslash or tab) is kept with its quotes and octal escapes; the fallback diff core/gitrepo/git.go:149 (Repo.Diff, 'git show --patch' without -c core.quotePath=false, used when no branch carries the handle) quotes every non-ASCII name. Example: file 'ä.go' becomes header '"b/\303\244.go"' and 'o' opens a file that does not exist; blame: 766032e54 (EH6Z0N) and 29a6fbb83 - no ZZTCWS commit on the path, diffPaneShown's predicate unchanged; ran: git show --patch on a scratch repo with ä.go prints 'diff --git "a/\303\244.go" "b/\303\244.go"'
