---
quick_id: 260824-pgj
slug: fix-the-critique-prompt-self-contradicti
date: 2026-08-24
status: complete
---

# The critique prompt no longer contradicts itself

`.planning/NEXT-STEPS.md` item 1, closed.

## What changed

`lanes/critique.md`, the closing rules only — 5 lines out for 14. Frontmatter,
the four questions and the three outcome branches are untouched; they were not
the defect.

Gone: **"Default to finding something."** A reader following that bold sentence
produces a finding on every pass, and `rejects-to: in-progress` turns that into a
loop with no exit. The escape clause under it ("if you genuinely find nothing on
the fourth pass") never stood a chance against the headline.

In its place, two rules:

- **A finding names a file and a concrete alternative.** "Could be cleaner",
  "consider extracting this", "this may not scale" are named as non-findings. The
  bias against lazy approval is kept, and put on the same footing as inventing
  something: "do not approve a diff you did not read, and do not manufacture a
  finding to look thorough — both produce a critique nobody can act on."
- **The loop ends when a pass finds nothing.** The end condition is now in the
  prompt: a pass with no finding *by that definition* is where the lane is done —
  `review-summary="none"` and move on — and that is "the expected way out, not a
  failure of nerve". Plus the two ways a pass fakes progress: re-raising a
  finding already addressed, and re-opening a trade-off let stand earlier.

## Verified

- `go test ./core/lane/... -race` green — `shipped_test.go` loads every file
  under `lanes/` and fails on warnings.
- `go test ./... -race` green; `gofmt -l core internal` lists exactly
  `core/gate/gate.go` and `internal/cli/tickets.go` (the documented baseline).
- Binary rebuilt to `~/.local/bin/jaira`; `jaira lanes show critique` prints the
  new rules.
- Copies refreshed: `~/.jaira/lanes/critique.md` (catalogue) and, on his board,
  `.jaira/lanes/critique.md` — which was byte-identical to the old shipped file,
  so no customisation of his was overwritten. **That copy is an uncommitted
  change in his repo, for him to commit.**

## Note

`README.md:187-200` also shows a `critique` lane, but as a shape example with
different prose and a different contract. Left alone.
