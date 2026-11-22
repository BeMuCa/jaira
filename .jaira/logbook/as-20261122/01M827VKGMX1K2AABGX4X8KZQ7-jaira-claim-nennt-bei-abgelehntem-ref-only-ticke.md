---
id: 01M827VKGMX1K2AABGX4X8KZQ7
title: "jaira claim nennt bei abgelehntem Ref-only-Ticket den richtigen Weg, und die 4539MZ-Release-Note markiert die geänderten Exit-Codes"
status: done
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
commits:
  - 6cfec1b75fc975973dfe139136b76ef62ef57644
  - 38ab5cb00fbe4be08a6f929085ab6350dc6bd4dd
created-at: 2026-11-22T18:21:23Z
updated-at: 2026-11-22T19:14:02Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-412670
claimed-at: 2026-11-22T18:43:32Z
outcome-what: "claim nennt bei auf dem Ref abgelehntem Ticket ohne Datei nur noch 'jaira pull <handle>'; die Exit-Code-Zeile der 4539MZ-Note ist mit '- ! ' markiert"
outcome-why: "Die alte Meldung riet, eine Datei zu löschen, die es nicht gibt; die Exit-Code-Änderung bricht Skripte und muss in jaira update hervorgehoben werden"
outcome-resolves: Die zwei Ränder aus der 4539MZ-Abnahme
review-summary: "claim.go: the refusal for a rejected ticket now picks its advice by t.ReadOnly - a ticket rejected only on its ref (no file here) is told just 'jaira pull <handle>', one with a file keeps 'delete its file here and jaira pull'. rejectref_test.go asserts both texts (ref-only: pull, no 'delete its file'; with file: 'delete its file here'). 4539MZ.md: the exit-code line now starts '- ! ' so jaira update highlights it, and the claim line names the ref-only path too."
review-gaps: none
test-verdict: "pass: go build/vet/test ./... green (exit 0); DoD 1 by TestATicketRejectedOnItsRefStillExistsForLinksAndCannotBeClaimed (ref-only claim names only 'jaira pull <h>', file case still names delete+pull), DoD 2 seen at 4539MZ.md:2; failure path is the refusal itself, exit 3 checked; claim refusal writes nothing, no two-writer case"
hunt-verdict: "pass: 0 findings in 3 modules"
review-verdict: "accept: both DoD items met in 6cfec1b75; ReadOnly is set exactly when the store has no file (core/ticket/store.go:1247); test passes; no defects in claim.go along the four lines"
review-round: 1
review-check: "1. go test ./internal/cli -run TestATicketRejectedOnItsRefStillExistsForLinksAndCannotBeClaimed -count=1  -> ok  2. sed -n 2p core/release/notes/4539MZ.md  -> line starts with '- ! '  3. sed -n 100,113p internal/cli/claim.go  -> ReadOnly branch gives only 'jaira pull <handle>'"
review-verdict-gaps: none
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
- **2026-11-22 18:49 · Alexander Sacharov** — bughunt: done at 6cfec1b75 - 0 findings in 3 modules, follow-up hunt (diff since b51c7b87f = commit 6cfec1b75); concurrency, error paths, input, permissions checked; files: core/release/notes/4539MZ.md internal/cli/claim.go internal/cli/rejectref_test.go
