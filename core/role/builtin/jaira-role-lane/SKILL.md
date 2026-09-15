---
name: jaira-role-lane
description: "Work exactly one jaira lane of one ticket, then stop. Invoked as /jaira-role-lane <ticket-id> <lane>, normally by a teamlead session. Use when handed a single lane to carry out and nothing else."
---

# One lane, then stop

Arguments: `<ticket-id> <lane>`. Without both, say what is missing and stop.

The lane's instructions are not here — the board carries them:

```bash
jaira claim <ticket-id> 2>/dev/null || true
jaira show <ticket-id> --for-lane <lane> --json
```

That gives you the lane prompt, the bounded input, and the outputs the lane owes
back. Follow it, produce exactly those outputs, and nothing beyond them.

Then finish the step yourself:

- `jaira dod <id> <n> --doing|--done` as you go, not at the end
- `jaira note <id> <text>` at every pause — dead ends, why this and not that.
  A killed session never gets a turn to write anything down.
- `jaira move <id> --to <next-lane> --what … --why … --resolves …`
- **did you change code?** Then commit it and the ticket file together: `git add`
  the files you changed and `.jaira/tickets/<this ticket>.md` by path — never
  `-A`, another session may be working in the same worktree — and name the ticket
  id in the commit subject, `fix(A3K9QP): …`. That handle is what jaira derives
  the ticket's commit list from; without it the move into the last lane is
  refused.
- **did you change no code? Then commit nothing.** A lane that only left a note
  and a lane change has nothing to show a reviewer, and a commit for that alone
  turns the branch history into one entry per lane with the actual work buried in
  it. Leave the ticket file modified in the worktree; the next commit that
  carries code takes it along; if no further code commit follows, the commit
  that files the ticket away with `jaira logbook <id>` carries its final state.

## Boundaries

- **This lane only.** A problem you spot outside it is a `jaira note` or a new
  ticket, never a fix you slip in. The lane you were given is the deliverable.
- **Your worktree only.** `git worktree list` first; if you are not in your own
  worktree, say so and stop rather than touching what another session holds.
- `review` and `human` are a person's lanes. Deliver into them and stop there.
- When you are done, report in three lines: what changed, what the next lane is,
  what is still open. A teamlead reads your pane, so the last thing on screen
  should be that summary and not a wall of diff.
