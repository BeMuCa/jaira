---
quick_id: 260824-pgj
slug: fix-the-critique-prompt-self-contradicti
date: 2026-08-24
mode: quick
---

# Quick task — the critique prompt contradicts itself

`.planning/NEXT-STEPS.md` item 1. His objection, in his words: "finde
standardmäßig etwas" against "loope bis die Kritisier-Lane nix mehr zu meckern
hat". If the lane always finds something, the loop never ends.

## The defect

`lanes/critique.md:61-65` ends the prompt with:

    Two rules for this lane:
    **Default to finding something.** A critique that approves everything is
    worth nothing. If you genuinely find nothing on the fourth pass, say so
    plainly rather than inventing a finding to look useful.

The escape clause is there, and the headline overshadows it. A reader following
the bold sentence produces a finding on every pass, and `rejects-to: in-progress`
turns that into an unterminating loop.

## Decision

Keep the bias against lazy approval, make termination explicit. Rewrite the
closing rules so that:

1. a finding names a file and a concrete alternative, or it is not a finding;
2. "nothing left" is the expected way the lane finishes, not a failure of nerve;
3. the loop's end condition is stated in the prompt itself.

Do not touch the frontmatter (`rejects-to`, the contract, the description already
says "until there is nothing left to say"), the four questions, or the three
outcome branches — they are not the defect. The README snippet at
`README.md:187-200` is a shape example with different prose; leave it.

## Tasks

1. Rewrite the closing rules of `lanes/critique.md`. Verify: the file holds no
   sentence that reads as "always find something"; the end condition is written
   in the prompt.
2. `go test ./core/lane/... -race` — `shipped_test.go` parses everything under
   `lanes/`. Verify: green, no load warnings.
3. `go build -o ~/.local/bin/jaira ./cmd/jaira`, then refresh the copies:
   `jaira lanes adopt lanes/critique.md --force` (catalogue) and, on his board,
   `jaira lanes use critique --force`. Verify: `jaira lanes show critique`
   prints the new text in both places.
4. `go test ./... -race` and `gofmt -l core internal` (must list exactly
   `core/gate/gate.go` and `internal/cli/tickets.go`).

## Done when

The prompt cannot be read as "always find something", `core/lane/shipped_test.go`
still passes, and the catalogue and his board carry the new text.
