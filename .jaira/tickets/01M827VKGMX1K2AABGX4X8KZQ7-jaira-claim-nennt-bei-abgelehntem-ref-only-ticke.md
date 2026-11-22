---
id: 01M827VKGMX1K2AABGX4X8KZQ7
title: "jaira claim nennt bei abgelehntem Ref-only-Ticket den richtigen Weg, und die 4539MZ-Release-Note markiert die geänderten Exit-Codes"
status: review
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: Die zwei kleinen Ränder aus der 4539MZ-Abnahme sind geschlossen.
context: |-
  Was falsch ist (Review-Notizen 4539MZ, 22.11.2026, angenommen ohne weitere Runde):
  1. internal/cli/claim.go:102 - bei einem auf dem Ref abgelehnten Ticket ohne Datei sagt die Meldung 'delete its file here and jaira pull'; es gibt keine Datei. Richtig: nur 'jaira pull <handle>', wenn t.ReadOnly.
  2. core/release/notes/4539MZ.md Zeile 2 ändert Exit-Codes, die Skripte lesen (1 -> 3/5), ist aber nicht mit '- ! ' markiert; jaira update hebt sie so nicht hervor. Noch unreleased: Zeile in derselben Datei ändern.
definition-of-done: "claim auf ein abgelehntes Ref-only-Ticket nennt nur 'jaira pull <handle>' (Test)"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-11-22T18:21:23Z
updated-at: 2026-11-22T18:48:31Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-412670
claimed-at: 2026-11-22T18:43:32Z
outcome-what: "claim nennt bei auf dem Ref abgelehntem Ticket ohne Datei nur noch 'jaira pull <handle>'; die Exit-Code-Zeile der 4539MZ-Note ist mit '- ! ' markiert"
outcome-why: "Die alte Meldung riet, eine Datei zu löschen, die es nicht gibt; die Exit-Code-Änderung bricht Skripte und muss in jaira update hervorgehoben werden"
outcome-resolves: Die zwei Ränder aus der 4539MZ-Abnahme
review-summary: none
review-gaps: none
test-verdict: "pass: go build/vet/test ./... green (exit 0); DoD 1 by TestATicketRejectedOnItsRefStillExistsForLinksAndCannotBeClaimed (ref-only claim names only 'jaira pull <h>', file case still names delete+pull), DoD 2 seen at 4539MZ.md:2; failure path is the refusal itself, exit 3 checked; claim refusal writes nothing, no two-writer case"
hunt-verdict: "pass: 0 findings in 3 modules"
---

# jaira claim nennt bei abgelehntem Ref-only-Ticket den richtigen Weg, und die 4539MZ-Release-Note markiert die geänderten Exit-Codes

## Definition of Done

- [x] claim auf ein abgelehntes Ref-only-Ticket nennt nur 'jaira pull <handle>' (Test)
  proof: internal/cli/claim.go:103; TestATicketRejectedOnItsRefStillExistsForLinksAndCannotBeClaimed
- [x] core/release/notes/4539MZ.md Zeile mit den Exit-Codes beginnt mit '- ! '
  proof: core/release/notes/4539MZ.md:2

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 18:43 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 18:46 · Alexander Sacharov** — Die claim-Zeile in 4539MZ.md gleich mitgeändert (noch unreleased): sie nannte nur den Weg mit Datei. Kein eigenes X8KZQ7.md, weil beide Änderungen 4539MZs eigene, noch nicht veröffentlichte Zeilen korrigieren.
- **2026-11-22 18:47 · Alexander Sacharov** — bughunt: done at b51c7b87f - 0 findings in 3 modules, first hunt (uncommitted change); concurrency, error paths, input, permissions checked; files: core/release/notes/4539MZ.md internal/cli/claim.go internal/cli/rejectref_test.go
