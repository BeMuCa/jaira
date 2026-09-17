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
- **did you change code?** Then commit it and the ticket file together:
  `git add` the files you changed and `.jaira/tickets/<this ticket>.md` by path
  — never `-A`, another session may be working in the same worktree — and name
  the ticket id in the commit subject, `fix(A3K9QP): …`. That handle is what
  jaira derives the ticket's commit list from; without it the move into the last
  lane is refused.
- **did you change no code? Then commit nothing.** A lane that only left a note
  and a lane change has nothing to show a reviewer, and a commit for that alone
  turns the branch history into one entry per lane with the actual work buried
  in it. Leave the ticket file modified in the worktree; the next commit that
  carries code takes it along; if no further code commit follows, the commit
  that files the ticket away with `jaira logbook <id>` carries its final state.

## Is this ticket in conversational mode?

The same `show --for-lane --json` you already read carries a `mode` key beside
`model_tier`. Read it there — not out of the line that started you, which is
gone the moment your session is.

Empty is the normal case and everything above stands unchanged.

`conversational` means a person found open design decisions on this ticket and
is reading along. Two things change, and only these two:

**1. Show the code after every definition-of-done item, not at the end.**

Work one item, then:

```bash
git diff
```

Empty output? No pause — a documentation item or one that only describes a
check produces no code, and pausing on it makes the mode slower than typing the
change by hand. Otherwise show that diff — the diff itself, not a summary of it
— beside the `--proof` you just recorded for the item, and wait for the person
before building anything on top of it.

Late is the expensive time to disagree with a shape. That is the whole reason
the mode exists.

**2. Do not commit. Hand back the commit line instead.**

You still `jaira move` the ticket, and you still leave the ticket file changed
in the worktree. What you do not do is run `git commit`.

Changed no code? Then hand back no line either — the rule above still holds, and
a mode does not suspend it. critique, testing and review run in this mode too,
because the mode sits on the ticket and not on a lane, and a commit line handed
to a person there produces a commit carrying the ticket file alone. In this mode
that lands harder than elsewhere: a line written out ready to paste reads as
already decided, and the person pastes it. Leave the ticket file in the worktree
for the next commit that carries code, exactly as you would outside the mode.

Changed code? Write the command out ready to paste, with the paths already
filled in:

```bash
git add <the files you changed> .jaira/tickets/<this ticket>.md
git commit -m "fix(A3K9QP): <what changed>"
```

**The handle in the subject is not decoration.** jaira derives the ticket's
commit list from two places — the ticket file's own history, which is thin on
purpose, and the commits naming its id. A commit written by hand without the
handle leaves that list empty, and the move into the last lane is then refused
however finished the work is. So hand back the whole line, handle included and
the ticket file already in the `git add`, rather than telling the person to
commit.

(The pattern is jaira-role-pr's: invoked by an agent it hands back the create
line instead of running it.)

## Boundaries

- **This lane only.** A problem you spot outside it is a `jaira note` or a new
  ticket, never a fix you slip in. The lane you were given is the deliverable.
- **Your worktree only.** `git worktree list` first; if you are not in your own
  worktree, say so and stop rather than touching what another session holds.
- `review` and `human` are a person's lanes. Deliver into them and stop there.
- When you are done, report in three lines: what changed, what the next lane is,
  what is still open. A teamlead reads your pane, so the last thing on screen
  should be that summary and not a wall of diff.
