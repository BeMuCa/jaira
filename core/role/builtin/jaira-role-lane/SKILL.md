---
name: jaira-role-lane
description: "Work exactly one jaira lane of one ticket, then stop. Invoked as /jaira-role-lane <ticket-id> <lane>, normally by a teamlead session. Use when handed a single lane to carry out and nothing else."
---

# One lane, then stop

Arguments: `<ticket-id> <lane>`. Without both, say what is missing and stop.

The lane's instructions are not here — the board carries them. Read first, and
write nothing at all until you have — `jaira claim` included, writing as it does
your name onto the ticket. Two reads, in this order:

```bash
jaira show <ticket-id> --json
jaira show <ticket-id> --for-lane <lane> --json
```

The first carries the ticket's `status`, and on a ticket in conversational mode
that is what decides whether you are the worker that writes or the critique
running beside it, which writes nothing. Only that JSON has the field — the
`--for-lane` one does not, so it cannot tell you. Read it here, before the
`jaira claim` below, and read it once: the section below says why a second read
gives the wrong answer.

The second gives you the lane prompt, the bounded input, the outputs the lane
owes back, and the `mode` key the section below turns on. Follow it, produce
exactly those outputs, and nothing beyond them.

Are you the lane that writes? Then take the ticket and finish the step
yourself. The critique running beside the work does none of the following — it
reads, and reports what it finds to the dispatcher that started it:

- `jaira claim <ticket-id>` — before you work it; other sessions read this board
  too. Not if the read above made you that critique: claiming is the first
  thing it must not do, because the name it writes is the one the implementing
  worker needs. `2>/dev/null || true`: a ticket already claimed by you is not
  an error worth stopping for
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
is reading along. Three things change, and only these three:

**1. Are you the critique that runs beside the work? Then write nothing.**

In this mode the dispatcher starts a second worker on the same ticket at the
same time — `/jaira-role-lane <id> critique`, while the implementing lane is
still running — and that worker reads only. Tell which of the two you are from
the ticket, not from the line that started you, which is gone the moment your
session is: the `status` in the `jaira show <id> --json` you already read at the
top of this file, before you wrote anything.

Your lane argument equal to the ticket's `status` means you are the lane the
ticket is actually in: work it as written above and below. Your lane argument
different from it means you are the critique running beside that lane. Then
read, and hand every finding to the dispatcher in the moment you have it —
and write nothing at all: no `jaira claim`, no `jaira dod`, no `jaira note`, no
`jaira set`, no `jaira move`, no `review-summary`, no commit line. Everything
this file tells a worker to write is off, claiming the ticket included. Exactly
one place writes those, and it is the dispatcher that started you.

**Read it once, and do not read it again.** The status moves under you: the
implementing worker ends its lane with `jaira move --to critique`, and from that
moment the ticket's status equals your lane argument. Ask a second time and the
answer flips — the running critique takes itself for the ordinary critique lane
and writes the `review-summary` the dispatcher is about to write, which is the
two writers on one field this whole point exists to prevent. So the first read
decides, and nothing later un-decides it.

That leaves one gap, and it has a rule of its own: if your session is restarted
or compacted and you can no longer tell from your own context which of the two
you were started as, **you are not the one that writes.** Say so to the
dispatcher and let it answer — it knows, it started you. Asking costs one line;
guessing wrong costs a second `review-summary` written over the first.

Report and stop. A running critique does not sit waiting for more work after it
has handed over its findings — it is done, and the dispatcher starts whatever
comes next.

This does not touch the ordinary critique lane: there the ticket's status is
`critique` at the moment that worker first reads it, and its lane argument is
`critique` too — the two match, and it writes its `review-summary` and moves the
ticket as always.

**2. Show the code after every definition-of-done item, not at the end.**

Work one item, then run **both** of these:

```bash
git status --short
git diff
```

`git diff` alone is blind to a file that does not exist in the index yet, and a
definition-of-done item made of one new file — a new test, a new package — is
the common case, not the corner. Take the empty diff as "nothing happened" and
that item passes the person by in silence, which is the one thing this mode is
here to prevent. `git status --short` is what sees it: an untracked file shows
there as `??`.

Both empty? No pause — a documentation item or one that only describes a check
produces no code, and pausing on it makes the mode slower than typing the change
by hand.

Otherwise pause. Show the diff itself, not a summary of it, and for every `??`
line show the new file's contents as well (`cat <path>`, or
`git diff --no-index /dev/null <path>` if you want it as a diff). Put it beside
the `--proof` you just recorded for the item, and wait for the person before
building anything on top of it.

Do not reach for `git add` to make the file visible — not even `-A -N`. Another
session may hold the same worktree, and this mode in particular runs the worker
in the checked-out directory rather than one of its own.

Late is the expensive time to disagree with a shape. That is the whole reason
the mode exists.

**3. Do not commit. Hand back the commit line instead.**

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
