---
id: 01M82ZBBANXWEKXWDTM1ZZTCWS
title: "Die Abnahme-Ansicht zeigt oben, wie man prueft, was gemacht wurde, die offene Frage, den Code-Diff und die DoD - und sonst nichts"
status: in-progress
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
updated-at: 2026-11-23T21:13:41Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-79231
claimed-at: 2026-11-23T20:54:57Z
---

# Die Abnahme-Ansicht zeigt oben, wie man prueft, was gemacht wurde, die offene Frage, den Code-Diff und die DoD - und sonst nichts

## Definition of Done

- [x] Die Ansicht in human, review und signoff zeigt oben in dieser Reihenfolge: So pruefst du, Was gemacht wurde, Offene Frage (wenn vorhanden), Code-Aenderungen, DoD mit Proofs
  proof: internal/tui/signoff.go:81-131 signOffBody order; view.go render() via atPersonsLane; tests TestSignOffViewShowsTheDecisionFirst, TestHumanLaneOpensTheDecisionScreen
- [ ] Felder ohne Nutzen fuer die Entscheidung sind von dieser Ansicht entfernt; welche, steht nach dem Brainstorm als Notiz am Ticket
- [ ] Ein leeres Pflichtfeld zeigt, welche Lane es haette schreiben muessen
- [ ] Release-Notiz in core/release/notes/<ID>.md
- [ ] Eine 'blocker:'-Zeile aus review-gaps bleibt oben sichtbar; problem, resolves, summary, verdict und die uebrigen gaps sind hinter einer Taste erreichbar; Test

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
