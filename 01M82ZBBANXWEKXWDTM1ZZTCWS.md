---
id: 01M82ZBBANXWEKXWDTM1ZZTCWS
title: "Die Abnahme-Ansicht zeigt oben, wie man prueft, was gemacht wurde, die offene Frage, den Code-Diff und die DoD - und sonst nichts"
status: todo
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
updated-at: 2026-11-23T20:54:31Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
---

# Die Abnahme-Ansicht zeigt oben, wie man prueft, was gemacht wurde, die offene Frage, den Code-Diff und die DoD - und sonst nichts

## Definition of Done

- [ ] Die Ansicht in human, review und signoff zeigt oben in dieser Reihenfolge: So pruefst du, Was gemacht wurde, Offene Frage (wenn vorhanden), Code-Aenderungen, DoD mit Proofs
- [ ] Felder ohne Nutzen fuer die Entscheidung sind von dieser Ansicht entfernt; welche, steht nach dem Brainstorm als Notiz am Ticket
- [ ] Ein leeres Pflichtfeld zeigt, welche Lane es haette schreiben muessen
- [ ] Release-Notiz in core/release/notes/<ID>.md
- [ ] Eine 'blocker:'-Zeile aus review-gaps bleibt oben sichtbar; problem, resolves, summary, verdict und die uebrigen gaps sind hinter einer Taste erreichbar; Test

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 20:54 · Alexander Sacharov** — Alexander Sacharov took this ticket over from BeMuCa
- **2026-11-23 20:54 · Alexander Sacharov** — decision: Alex 23.11.2026: vom Teamlead uebernommen, sofort umsetzen. Gestapelt auf EH6Z0N (berkc, signoff, Branch upstream/feat/EH6Z0N-diff-view, noch nicht auf master): beide aendern internal/tui/signoff.go, der Diff-Teil kommt von dort. Der Branch feat/zztcws enthaelt deshalb die EH6Z0N-Commits; gemergt wird nach EH6Z0N. Brainstorm-Vorgabe vom Teamlead: problem, resolves, summary und verdict fallen von der Ansicht und bleiben hinter einer Taste (z.B. 'i' Details) erreichbar; eine 'blocker:'-Zeile aus review-gaps bleibt oben sichtbar, weil sie die Entscheidung aendert; die uebrigen gaps hinter derselben Taste. Die genaue Liste schreibt die Brainstorm-Lane als Notiz (DoD 2).
