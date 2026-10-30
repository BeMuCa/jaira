---
id: 01M658C9JPRZQPQA5JMJER0PDG
title: "Ein Ticket in human review laeuft im Hauptarbeitsbaum statt im eigenen Worktree, damit UI-Tests nur eine App auf einem Port brauchen"
status: ideas
ready: true
creator: berkc
assignee: berkc
goal: "UI-Abnahme laeuft gegen eine einzige laufende App: was in human review ankommt, landet im Hauptarbeitsbaum."
context: |-
  Idee von berkc am 30.10.2026: 'merging worktrees into the main work when the ticket reaches human review ... That way when we have UI tests, we dont have to call the same frontend on two different ports.'
  Heute: jeder Worker hat einen eigenen Worktree; reap.sh entfernt ihn erst bei done + gemergtem PR (80MGXT: Worktree weg bei done).
  Problem dabei: mehrere Tickets gleichzeitig in human review kollidieren im Hauptbaum; 'nichts landet direkt auf master'.
  Alternative: ein Integrationszweig (z.B. review/integration), in den jedes Ticket bei human review gemergt wird, und die App laeuft von dort - ein Port, master bleibt sauber.
definition-of-done: "Entscheidung: Hauptbaum vs. Integrationszweig, mit Begruendung"
tags:
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-10-30T01:56:49Z
updated-at: 2026-10-30T23:45:40Z
updated-by: berkc
---

# Ein Ticket in human review laeuft im Hauptarbeitsbaum statt im eigenen Worktree, damit UI-Tests nur eine App auf einem Port brauchen

## Definition of Done

- [ ] Entscheidung: Hauptbaum vs. Integrationszweig, mit Begruendung

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-30 23:45 · berkc** — berkc 31.10.2026: 'wouldnt it be better to just pull it into the main worktree because we also have the advantage that the ticket updates land in the mainboard aswell? integration worktree seems useless.' Gegenargumente zur Entscheidung: (1) der Hauptworktree steht auf master, und auf master landet nichts direkt - ein Merge dort waere ein lokaler Commit auf master; (2) zwei Tickets gleichzeitig in human review mischen sich im Hauptbaum, und ein zurueckgeschicktes Ticket muss wieder heraus (Revert); (3) die eigenen ungesicherten Aenderungen im Hauptbaum geraten dazwischen. Ticket-Updates landen ohnehin auf den Refs, das Board sieht sie ohne Merge. Moeglicher Mittelweg: Checkout des Ticket-Branches im Hauptbaum (git switch) fuer die Abnahme statt Merge.
