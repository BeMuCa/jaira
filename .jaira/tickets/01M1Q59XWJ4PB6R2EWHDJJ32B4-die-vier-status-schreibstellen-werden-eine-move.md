---
id: 01M1Q59XWJ4PB6R2EWHDJJ32B4
title: Die vier Status-Schreibstellen werden eine Move-Funktion
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Gate-Check, Claim-Regel, Mutation und Settle leben in EINER core-Funktion; CLI move, TUI applyMove/forceMove und accept rufen sie und behalten nur UI-Belange (Formatierung, pending/confirm, --from-lane-Validierung)"
context: "Berk am 04.09. (Feedback L122, zweiter Teil): der bekannte 'eigene Schnitt' aus den Learnings - vier Schreibstellen (cli/flow.go, tui applyMove, forceMove, signoff accept) teilen heute moveMutation nur paarweise; jede neue Nach-Move-Regel braucht alle vier, und grep findet accept nicht (schreibt Status direkt). Der Settle-Teil wurde bereits zentralisiert (Vorgaenger-Ticket); dieser Schnitt vereinheitlicht den Rest. Gross und invariantennah (CLI und TUI muessen identisch gaten, Interactive-Flag, Claim-beim-Pull) - bewusst NICHT im Schnellverfahren gebaut."
definition-of-done: "Eine core-Move-Funktion traegt Gate+Mutate+Settle; alle vier Aufrufer delegieren; CLI- und TUI-Verhalten byte-gleich (bestehende Tests gruen); ein Test beweist, dass eine neue Nach-Move-Regel nur noch eine Stelle braucht"
tags: []
blocked-by: []
commits: []
created-at: 2026-09-04T21:30:57Z
updated-at: 2026-09-04T21:30:57Z
---

# Die vier Status-Schreibstellen werden eine Move-Funktion

## Definition of Done

- [ ] Eine core-Move-Funktion traegt Gate+Mutate+Settle; alle vier Aufrufer delegieren; CLI- und TUI-Verhalten byte-gleich (bestehende Tests gruen); ein Test beweist, dass eine neue Nach-Move-Regel nur noch eine Stelle braucht

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

