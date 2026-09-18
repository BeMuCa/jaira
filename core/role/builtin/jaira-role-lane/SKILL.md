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
running beside it, which writes nothing — but only when you were started on
`critique`. Started on any other lane you are that lane and you write, whatever
the `status` says. Only that JSON has the field — the `--for-lane` one does
not, so it cannot tell you. Read it here, before the `jaira claim` below, and
read it once: the section below says why a second read gives the wrong answer.

The second gives you the lane prompt, the bounded input, the outputs the lane
owes back, and the `mode` key the section below turns on. Are you the lane that
writes? Then follow it, produce exactly those outputs, and nothing beyond them.
Are you the critique running beside the work instead? Then that prompt is
something to read and not an order: the critique lane's own prompt tells you to
`jaira set <handle> review-summary=…` and to `jaira move`, and those are the
writes the section below takes away from you. Read it for what the lane is
meant to judge, and leave its outputs to the dispatcher that started you.

Does your lane judge that payload's diff? The ordinary `critique`, `testing`
and `review` lanes do — and the payload names what it was built from. Its
`commits` field lists the SHAs the diff covers and `commits_source` says where
they came from: the diff is the ticket's whole committed history, not a slice
of it, and where `commits_source` ends in `+worktree` the uncommitted work is
appended to it too — so a lane that changed no code and a conversational ticket
that has not been committed yet are both in front of you without your going to
look. The one commit that escapes is a commit that neither names the ticket id
nor touches its file; that is what the rule "every commit names the ticket id"
is for, and you can see it by counting `commits` against
`git log origin/HEAD..HEAD --oneline`.

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

Both halves have to hold: your lane argument is `critique`, **and** the
ticket's `status` is something else. Anything else is the lane the ticket is
actually in — work it as written above and below, including the ordinary
`critique` lane on a ticket whose `status` already says `critique`. The
dispatcher starts every worker before it moves the ticket into that worker's
lane, so a lane argument that merely differs from the `status` is the normal
case and says nothing; `pre-process` on a ticket still in `todo` is a
`pre-process` worker that claims, works and moves, not a silent reader. Only
`critique` beside a `status` that is not `critique` is this critique. Then
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

**What you judge is in the payload, and it is the second half of the diff.**
The `--for-lane critique` payload you read at the top already carries the
uncommitted work: `showForLane` appends the working tree to the commit diff, so
the `diff` has two halves separated by the line

```
uncommitted work in the working tree
```

Everything above that line is the earlier rounds — finished, already judged,
not yours. Everything below it is the work running beside you, and that is your
judgment object. On the first round there is no line at all, because there are
no commits yet: the whole diff is the worktree.

You do not have to hunt for the marker. `commits_source` says which halves are
there: it ends in `+worktree` when the uncommitted half is appended, and on the
first round it is plain `worktree`.

`complete: false` on that first round is still normal, but not because the diff
is missing — it is not. The critique lane also requires `outcome-what` and
`outcome-resolves`, and those are written by the worker beside you when it ends
its lane. So: do not wait for the payload to fill, and do not report that you
had nothing to read.

One case leaves you without that half, and the payload names it: a
`worktree_error` key means git could not read the working tree, so nothing
uncommitted is in the diff however clean `commits_source` looks. That is when
you fall back to reading the disk yourself:

```bash
git status --short -- :/ ':(exclude,top).jaira/tickets'
git diff -- :/ ':(exclude,top).jaira/tickets'
git diff --cached -- :/ ':(exclude,top).jaira/tickets'
```

`git diff` is the work in progress, `--cached` anything already staged, and
`git status --short` the new files neither of them shows — read a `??` line's
file with `cat`. The pathspec is what keeps that readable: the worker beside
you rewrites the ticket file with every `jaira dod` and every `jaira note`, so
without it you get hundreds of lines of the ticket's own prose around the few
lines of code you came to judge. `:/` is the repository root and `,top` anchors
the exclusion there, so all three commands say the same thing from any
directory.

The goal, the definition of done and the notes in the same payload are the
ticket as it stands right now — the critique prompt asks for the notes before
anything else, and they are what says which findings are already closed. They
stay the measure. The one thing neither they nor the diff carry is what the
implementer has written onto the ticket since the payload was built: read that
with `jaira show <id> --json`, not out of a diff.

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
git status --short -- :/ ':(exclude,top).jaira/tickets'
git diff -- :/ ':(exclude,top).jaira/tickets'
```

The pathspec is not decoration either. From your first `jaira dod` onwards the
ticket file under `.jaira/tickets/` is modified and stays modified for the rest
of the lane, so a bare `git status --short` is never empty again: every later
item would pause, and what you would put in front of the person is the ticket's
own diff. Excluding that one directory is what keeps "no code, no pause" true
after the first item. `:/` is the repository root, so both commands say the same
thing from any directory.

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
