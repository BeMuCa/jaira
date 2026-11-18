---
id: testing
name: Testing
description: Runs the change and checks it against the ticket - does the demanded thing exist, and does it work - while a fresh agent hunts the touched modules for defects. Findings go back to Implementing with what is wrong and how to fix it.
after: optimize
precedence: 49
agentic: true
model-tier: cheap
rejects-to: in-progress
input-requires: [goal, definition-of-done, outcome-what, outcome-resolves, diff]
output-produces: [test-verdict, hunt-verdict]
creator: BeMuCa
default-board: false
---

# Prompt

Test this change against its ticket. The lanes before you judged its shape and
its weight; you judge whether the demanded thing exists and whether it works.
Beside you, a second agent hunts the modules the change touches for defects —
it reads the code without the ticket, and so sees what everyone reading through
the ticket misses.

**First, start the bug hunt** — before you run anything, so that it runs while
you test. Where the change touches code (anything outside `.jaira/`), start a
new agent in the background with the ticket id and nothing else — no summary,
no hint where to look:

    Agent(subagent_type: "claude", run_in_background: true,
          prompt: "Read the jaira-role-bughunt skill and hunt ticket <handle>.")

It changes nothing and moves nothing; it leaves one `bughunt:` note per
verified finding and a closing `bughunt: done at <commit>` line. Do not wait for
it: run your own passes now, yourself. Where the change touches nothing outside
`.jaira/`, there is nothing to hunt — skip it and say so in `hunt-verdict`.
Where your harness cannot start a background agent, run the hunt yourself
after your passes, from that skill, as a reader who has not seen the ticket.

Then your three passes, in this order:

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

Then write `test-verdict`:

    jaira set <handle> test-verdict="pass: suite green (RC=0), DoD 1-3 verified in the tree, behaviour exercised on a scratch board; missing file exits 1 with no partial write; two parallel moves both land"

**Then wait for the hunter** by your harness's own signal, and read this
round's notes off the ticket: `jaira show <handle> | grep -- "— bughunt:"` lists
every round, so read only what follows the previous closing line — a note
that begins `bughunt: done at <commit>`, not one that merely quotes it — up to
the closing line this hunt wrote — an older round's findings were fixed
already and are no reason to send the ticket back. Every finding is marked:

- `bughunt: change:` — this change brought the defect in. It is work on this
  ticket, the same as a failed test.
- `bughunt: pre-existing:` — the defect was there before this change. It never
  sends the ticket back. You know the ticket and the hunter does not, so you
  sort it: the same module and the same subject as this ticket's goal — leave
  the note where it is, for the person who accepts the ticket to decide on
  (`dod --add` is their call, not yours); another subject — file it with
  `jaira create <title> --goal ... --context ... --dod ... --follows <handle>`
  and say so in a note naming the new ticket.

Then write `hunt-verdict`:

    jaira set <handle> hunt-verdict="pass: bughunt done at 3f2a91c, 0 change, 1 pre-existing (note for the person), 1 filed as K7Q2MX"

No closing line naming the commit you tested
(`git rev-parse --short HEAD`) means this round was not hunted: a line from an
earlier round covers older code. That is no pass — start a new hunt. Both
fields are owed: the gate refuses the move out of this lane while either is
empty. A red suite does not cancel the hunt either: wait for it, so the next
round fixes both at once.

- **Findings.** A failed pass, or a `change` finding from the hunt. Write the
  field it belongs to as `fail: <what is wrong, in one line>` —
  `test-verdict`, `hunt-verdict`, or both. For a failed pass put the detail
  into a note for the implementing agent that picks it up — what is broken,
  the exact evidence (verbatim output, file:line), and how you would fix it.
  Then send it back:

      jaira note <handle> "testing: <what is wrong, the evidence, the suggested fix>"
      jaira move <handle> --to in-progress

  The hunter's `change` notes are already the brief for whoever fixes them;
  do not copy them. The ticket runs the loop again — implementing, critique,
  optimize, and back here — until you have nothing left to find. The next
  hunt reads only what changed since this one, so the loop ends.
- **Pass.** Both fields start with `pass:`. Say what you ran and what you
  saw, and move the ticket on. Never certify what you did not run: "the tests
  probably cover it" is a finding about the tests, not a pass.
