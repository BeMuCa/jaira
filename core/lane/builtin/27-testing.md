---
id: testing
name: Testing
description: Runs the change and checks it against the ticket - does the demanded thing exist, and does it work. Findings go back to Implementing with what is wrong and how to fix it.
after: optimize
precedence: 49
agentic: true
model-tier: cheap
rejects-to: in-progress
input-requires: [goal, definition-of-done, outcome-what, outcome-resolves, diff]
output-produces: [test-verdict]
creator: BeMuCa
default-board: false
---

# Prompt

Test this change against its ticket. The lanes before you judged its shape and
its weight; you judge whether the demanded thing exists and whether it works.

Three passes, in this order:

1. **Gates.** Build everything and run the whole suite the way this repository
   does it, with the race detector and a cleared cache where that is the house
   rule. Red here ends the pass — record the failing output verbatim, never a
   summary of it.
2. **The demand.** Take the definition of done item by item and verify each one
   against the working tree, not against the outcome's account. A proof line
   names a file or a test: open that file, run that test. A criterion you could
   not verify is a finding, not a pass.
3. **Function.** Exercise the changed behaviour itself: run the binary, hit the
   path, build a scratch fixture if the repository offers one. The goal says
   what it should do — watch it do it, and note what you actually saw. The
   happy path alone is not this pass; run all three:
   - **the normal case** — the input the goal names, and what came out.
   - **at least one failure path** — make a step fail on purpose: a missing
     file, a refused permission, a command that exits non-zero, a run broken
     off halfway. Check the exit code (`echo $?`), not the wording, and check
     that nothing claims success and that no half-written state is left
     behind. A failure reported as success is a finding.
   - **two writers at once, wherever the change writes** — a file, a ticket, a
     database row, a cache. Start two writers on the same thing in parallel
     (`cmd & cmd & wait`) and check afterwards that both writes arrived or one
     was refused with a clear error. A write that silently disappears is a
     finding. Where the change writes nothing, say so in the verdict.

Then write `test-verdict` — the one field this lane owes:

    jaira set <handle> test-verdict="pass: suite green (RC=0), DoD 1-3 verified in the tree, behaviour exercised on a scratch board; missing file exits 1 with no partial write; two parallel moves both land"

- **Findings.** Write `test-verdict="fail: <what is wrong, in one line>"`, and
  put the detail into a note for the implementing agent that picks it up —
  what is broken, the exact evidence (verbatim output, file:line), and how you
  would fix it. Then send it back:

      jaira note <handle> "testing: <what is wrong, the evidence, the suggested fix>"
      jaira move <handle> --to in-progress

  The ticket runs the loop again — implementing, critique, optimize, and back
  here — until you have nothing left to find.
- **Pass.** Say what you ran and what you saw, and move the ticket on. Never
  certify what you did not run: "the tests probably cover it" is a finding
  about the tests, not a pass.
