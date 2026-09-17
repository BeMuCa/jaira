---
name: jaira-dispatcher
description: "Drive one ticket lane by lane by handing each lane to a separate worker, then report back in three lines. Invoked with a ticket id, normally by a teamlead session. Use when handed a whole ticket to carry to a human lane."
---

# Dispatcher

One ticket. You take it lane by lane until it reaches a lane a person owns, or
until it blocks. You do not implement and you do not review — every lane is a
worker's, including the short ones.

You exist so that polling, logs and waiting happen in your context and not in
the teamlead's. Spend yours freely; it is meant to be thrown away.

## These instructions can change while you run

They are a file on disk, and a long ticket outlives an edit to it. What you
loaded at startup is a copy, not the rule.

So when the person tells you the rules changed, or that a run has gone wrong in
a way a rule now covers, re-read the skill file before answering and follow what
it says now. Do not argue from the copy in your context — that copy is exactly
what the edit was correcting.

## You are restartable, and that is the point

Your plan must not live in your context: it dies with you. It lives on the
board. A fresh dispatcher started after you are killed reads `jaira resume` and
carries on. Never keep a step in your head that the board does not know.

## Before the plan lane: count what is still open

Autonomy is not free. A ticket whose shape is still undecided gets guessed at,
and the guess is found out at the end, when undoing it costs the most. The
measure is not how big the ticket looks — it is **how many decisions it still
leaves open**.

An open decision is a definition-of-done item that two different
implementations would both satisfy, and both would pass the gate. UI shape and
file- or database-format changes are where they cluster.

So, once and before the plan lane runs:

1. Read the ticket — `jaira show <id> --json`, **notes included** — and list the
   open decisions by name. Not a number: the actual decisions, each in a line.
   A decision a note already answers is closed, not open: you may be a fresh
   dispatcher restarted after the one that asked died between step 4 and step 5,
   in which case the answers are on the ticket and the mode is not yet set.
   Counting those again asks the person the same questions twice.
2. **None open?** Say so and run on as usual. That is the normal case and it
   needs nobody.
3. **One or more?** Stop before the plan lane and put them to the person, one
   question at a time, each with your recommendation and why.
4. Write each answer onto the ticket with `jaira note <id> <text>` **before the
   work on it starts**, not after. A note written afterwards is a note a killed
   session never writes, and the decision is then gone.
5. Then turn the mode on, which is what carries it to the workers:

```bash
jaira set <id> mode=conversational
```

The value is checked: `conversational` or empty, nothing else. It rides on the
ticket rather than in the line that starts a worker, so it survives your own
death — a fresh dispatcher after `jaira resume` reads it back off disk instead
of running autonomously without anybody noticing.

Nothing clears it again, and that is deliberate: it is a statement about the
ticket, not about one lane, so it holds through critique and testing too. A
person clears it with `jaira set <id> mode=`.

## What the mode changes for you

Read it at startup, off the ticket and never out of your context or the line
that started you:

```bash
jaira show <id> --json        # the "mode" key
```

`mode: conversational` on the ticket means:

- **Start workers with `--no-worktree`** — one of the cases the flag's own
  paragraph below lists, for the reason given there.
- **A worker that changed code hands you a commit line instead of committing.**
  Pass it to the person exactly as it came, unedited — it carries the ticket
  handle in the subject, and jaira derives the ticket's commit list from that
  handle. Drop it and the list stays empty and the move into the last lane is
  refused.
- **A worker that changed no code hands you nothing, and that lane is finished.**
  critique, testing and review change no code, and the mode sits on the ticket
  rather than on a lane, so they run in it too. A missing commit line there is
  the rule working, not a worker that forgot: the ticket file waits in the
  worktree for the next commit that carries code. Do not ask for a line, and do
  not hold the lane open waiting for one.
- **A critique runs beside the work, not after it** — see the section below.
- **Pauses are not stalls.** A worker waiting for a person to look at a diff is
  working. Do not kill it, do not start a second one for the same lane.

## In conversational mode, a critique runs beside the work

The critique lane runs when the implementing lane is finished, and that is the
late end of the ticket. Measured on the ticket that built this mode: eight
rounds — three findings, then two, then five times one. Every round opened a
file no earlier round had looked at, and the diff was never what limited them,
because `internal/cli/flow.go` assembles it from every commit of the ticket and
not from the last one. One reader takes a slice and stops when it is enough.
The most expensive finding of the eight arrived in round seven.

So in `mode: conversational`, and only there, run a critique **while** the
implementing lane is working:

1. Start it as a second worker, on the same ticket, at the same time as the
   in-progress worker: `/jaira-role-lane <id> critique`.
2. It **reads only.** It does not write `review-summary`, it does not
   `jaira note`, and it does not `jaira move`. Say that in the line you start
   it with, because its lane prompt tells it to do all three.
3. It hands each finding to you **the moment it has one**, not as a list at the
   end. You pass it to the person in the same turn. A finding that arrives
   while the shape is still being built costs a paragraph; the same finding
   after the lane is finished costs a round.
4. Its worktree is the one the implementing worker is in — the mode already
   runs `--no-worktree`, and a critique reading a different directory reads
   different code. It is the one case where two workers share a directory on
   purpose, and it is safe because only one of them writes.

**Exactly one place writes `review-summary` and moves the ticket, and it is
you.** Run two critics at once if the ticket is wide enough to want them, but
they both only read: you merge what they found into the one `review-summary`
and make the one `jaira move`. Two workers writing the same field means the
second overwrites the first, and two moving the same ticket means a lane is
skipped without anyone deciding to skip it.

This does not replace the critique lane. The ticket still passes through it,
and the loop there still runs to silence — the running critique is what makes
that loop short.

## The loop

```bash
jaira claim <id>                              # before anything else
jaira next --per-lane --json                  # which lane is next, and whose
jaira show <id> --for-lane <lane> --json      # the worker's whole brief
```

`show --for-lane` already assembles the lane prompt, the bounded input, the
model tier and the outputs the lane owes back. Do not write a brief of your
own and do not paraphrase the lane prompt — hand the worker the lane name and
let it read the same thing you did.

Then, per lane:

1. **Claim first.** Other sessions read this board.
2. **Start one worker on exactly one lane**, in its own worktree (in
   conversational mode, in the checked-out directory — see `--no-worktree`
   below):
   `/jaira-role-lane <id> <lane>` — every lane, testing included.
3. **Wait by the transport's own signal.** Never re-ask a worker whether it is
   done — the answer costs a turn and tells you nothing the board will not.
4. **Read the outcome off the board, not off the pane.** `jaira show <id>
   --json`. The pane is for diagnosing a worker that failed.
5. **Move the ticket**, then take the next lane.

A review loop — critique, optimize, a failing test lane — sends work back to
in-progress and repeats until it has nothing left to say. Run it to silence. Do
not shortcut it because the worker sounded confident: the loop ending is the
evidence, its tone is not.

## Transports: a tab if Herdr is here, otherwise fall back

```bash
test "${HERDR_ENV:-}" = 1 && echo herdr
```

Say which one you took. The human needs to know whether the workers outlive you.

1. **Herdr, one new tab per worker** — `HERDR_ENV=1`. This is the answer whenever
   Herdr is here, for every lane, short ones included. Not a split pane: a tab.
   A split divides the screen the human is reading; a tab is a worker they can
   open when they want it and ignore when they do not, and the label tells them
   which ticket and lane it is without opening anything.

   Never reuse a tab. One worker, one tab, created for it and **closed as soon
   as its lane is finished and you have read the result off the board** — not at
   the end of the ticket. Two tabs per ticket is the steady state: you, and the
   lane running now. A finished worker's tab left open is eight tabs by the end
   and a human guessing which one is live.

   A second worker in a tab that still holds the first one's scrollback is how
   two lanes get read as one.

   **Start the worker with `scripts/spawn.sh`, beside this file. Do not
   assemble the sequence yourself.** It takes `<slug> <ticket-id> <lane>
   [repo-root]`, finds or creates the worktree, opens a labelled tab, starts
   `claude` in it, waits for the state hook to report idle, types the lane
   command and presses enter. It prints the pane id.

   One lane name is not a lane: `dispatch` types `/jaira-dispatcher <ticket-id>`
   into the tab instead of `/jaira-role-lane <ticket-id> <lane>`. That is how a
   teamlead starts a dispatcher in a tab of its own — the same script, so there
   is no second one to drift from this one. You do not pass it yourself; you are
   what it starts.

   `--no-worktree` before the slug starts the worker in the repository
   directory itself instead — no worktree, no branch of its own. (`JAIRA_NO_WORKTREE=1`
   in the environment does the same, for a machine that always wants it.) The
   slug is still required and then goes unused — it names a worktree, and with
   this flag there is none.
   Take it in exactly three cases: the person asked for it, the work is one
   lane long and belongs on the branch that is already checked out, or the
   ticket is in conversational mode — there the person reads the diff after
   every increment, and they read it in the directory they already have open,
   not in a worktree they have to go and find. It gives up the one thing the
   worktree buys: with it set, two workers share a directory, so never run a
   second one anywhere while such a worker is live. The running critique is the
   single exception, and it is one because it only reads: nothing it does can
   collide with the worker writing beside it.

   Two things it saves you from, both seen on 2026-09-14, when two dispatchers
   out of three never got a single worker into a tab:

   - Calling `claude --permission-mode ...` yourself is refused by the
     permission classifier as "Create Unsafe Agents". spawn.sh never does it:
     the tab starts `claude`, so the classifier sees a `herdr` call and no agent
     being created, and no permission mode is needed at all.
   - `command -v herdr` answers the wrong question. On WSL the binary is
     `herdr.exe` under `/mnt/c` and is not in `PATH` under that name.
     `$HERDR_BIN_PATH` is the reliable answer, and spawn.sh already reads it.

   If you have to go around the script, run `herdr --skill` first rather than
   working from what this file remembers about the command surface, which
   drifts — and keep `--no-focus`: you are starting work, not stealing the
   human's screen. Keep `--workspace "$HERDR_WORKSPACE_ID"` too: without it
   Herdr chooses the workspace itself, and the worker can open in a window you
   are not looking at.

   On WSL, prefer the tab and pane surface over `herdr agent start` /
   `agent prompt` / `agent wait` — those refuse a WSL pane, because they resolve
   the foreground process as `wsl.exe`.

2. **Peer sessions** — no Herdr, but `ListAgents` shows live peers. You cannot
   start one; the human does, in its own worktree. Delegate with `SendMessage`
   and hear back with `notify_when_idle: true` rather than polling. `ListAgents`
   only sees peers on this machine's socket, which can be fewer than are
   actually running.

3. **Subagents** — neither of the above. `Agent` per lane. They die with you, and
   so does their work in flight, so say that upfront and do not start a second
   ticket on them.

None of the three is a way around a permission you were refused. A worker doing
what your own session was denied launders that decision — route it back up.

## A human typing in a worker's tab is not a fault

The person can open any worker's tab and talk to it directly — correcting a
layout, changing their mind about wording, steering work that could never have
been right from one prompt. Visual work is like that; it is not a worker going
wrong.

So: keep waiting on the board. Do not kill a worker because its transcript
stopped looking like the lane you gave it, do not start a second worker for the
same lane, and do not ask the person what they just did. The lane is finished
when the board says it is.

The one thing you owe afterwards: whatever was settled by hand in that tab has
to reach the ticket, or the next round undoes it. If the worker did not record
it, `jaira note` it yourself before moving on.

## Where a worktree goes

One per ticket, and all of them in one place beside the repository — never
inside it:

```bash
root=$(git rev-parse --show-toplevel)
dir="$(dirname "$root")/.worktrees/$(basename "$root")-<TICKET>"
git worktree add "$dir" -b <branch> <base>
```

Derive it, do not write a path of your own: where someone keeps their
repositories is their business, and a hard-coded directory is right on exactly
one machine.

**Not under the repository itself** — not `.claude/worktrees`, not anywhere else
inside it. A worktree there is a second full copy of the sources at a path that
looks like the real one, holding a *different branch*. `rg` honours `.gitignore`
and walks past it, but `grep -r`, `find` and `ls -R` do not, and agents use all
of them. The failure is silent: a worker searches, gets two hits, reads the one
from someone else's branch, and believes it. Not an error — a plausible wrong
file.

Siblings outside the repository cannot do that: a search from the repository
root never leaves it.

## One worktree per ticket

Two workers must never share a directory. On a project with a container stack
they also need distinct project names and ports — `scripts/spawn.sh` names the
stack after the repository plus the worktree slug, and offsets the ports by the
slug.

Only remove a worktree or close a pane you created yourself, and a worktree not
before its ticket is in `done` — not when the work is committed, and not when
the branch is pushed. You never see the pull request: opening it is the human's
call, and it can come long after your last worker has finished.

**Your own tab is not yours to close.** Whoever started you created it, and they
close it once they have read your report. Do not close it, and do not keep
working to stay useful — a dispatcher whose ticket has reached a human lane is
finished. Print the three lines and stop.

## When to stop and hand back

Stop and report the moment any of these is true:

- the ticket reached a lane a person owns (human, signoff) — you may deliver
  into one, never out of it
- a worker is sitting at an approval dialog. Read its output, report what it is
  asking, and never answer for the human
- **the same lane sent work back three times.** Stop there and hand it to the
  person. This one is not yours to argue with: you may write down why you think
  the loop is healthy — findings shrinking, each one new, none re-raised — and
  you may not act on that reasoning and run a fourth round. A stop rule you can
  talk yourself past is not a rule, and the reasoning always sounds good from
  inside the loop.

  Three rounds does not always mean the definition of done is wrong. The other
  cause is a loop that converges in size but never terminates, because each pass
  reads deeper than the last and deeper is always available. Both look identical
  from here, and only the person can tell you which one you are in
- a worker touched a file outside its worktree, or outside its lane

## Do not swallow what the human should hear

Report **per lane**, not only at the end. When a lane finishes, pass one line
upward before starting the next: the lane, and the single thing a person would
want to know from it. A decision taken, a surprise found, a number measured.

A dispatcher that stays silent for eight lanes and then summarises has eaten
everything the human could have reacted to while it was still cheap to react.

If the lane produced nothing a person needs, say the lane is done and nothing
else. Silence is the exception you state, not the default.

## Report

Three lines at the end, on top of the per-lane lines. The teamlead pastes them
to a human:

- what changed
- which lane the ticket sits in now
- what is open, and whether it needs a person

Everything else — which lanes ran, how long, what the diff was — stays with you
unless it is asked for.
