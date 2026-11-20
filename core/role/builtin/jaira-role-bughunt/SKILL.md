---
name: jaira-role-bughunt
description: "Hunt for defects in whole modules with fresh eyes - concurrency, error paths, unusual input, permissions - verify every finding against the code, and leave each one on the board. Invoked as /jaira-role-bughunt [ticket-id | path...]. Use when the testing lane starts it beside its tests, or on any module that nobody has read as a whole."
---

# Find the bugs nobody was looking for

A ticket that passed critique, testing and review can still ship defects:
every one of those lanes reads the change through the ticket — its goal, its
diff, often the implementer's own account of it. Code that was copied in or
moved along with the change is outside that frame, and so is every fault that
lives in the lines around a hunk rather than in it. An external reviewer
without the frame finds them. You are that reviewer.

So you do not read the ticket's goal, context, outcome or notes, and nobody
briefs you. What you need is the code.

You change nothing: no source file, no commit, no `jaira move`, no `dod`, no
`set`. You write findings, and only in the two places below.

## What you read

**Given a ticket id:** the files its change touches, and nothing of its story.

```bash
# the files of every commit that names the ticket, and the uncommitted ones;
# ticket files — on the board, filed, rejected — stay out: they are the story
x=(':/' ':(exclude,top).jaira/tickets' ':(exclude,top).jaira/logbook' ':(exclude,top).jaira/archive')
{ git -c core.quotePath=false log -i --grep=<id> --name-only --format= -- "${x[@]}"
  git -c core.quotePath=false diff --name-only HEAD -- "${x[@]}"
  git -c core.quotePath=false ls-files --others --exclude-standard -- "${x[@]}"; } | sort -u
```

The commits are found by the ticket id anywhere in their message, the way jaira itself
derives a ticket's commits — not by a merge base: `origin/HEAD` can name a
fork that is hundreds of commits behind, and the files of all those commits
are not this change. Read nothing else of the ticket: `jaira show` only ever
runs piped into the one `grep` below, for your own closing line of the last
round.

**Given paths:** those files or directories.

**Has this ticket been hunted before?** Your own closing line from the last
round is the one thing you read off the ticket — that line, nothing else:

```bash
jaira show <id> | grep -E -- "— bughunt: done at [0-9a-f]{7,} " | tail -n 1
```

The pattern is anchored on the dash `jaira show` puts before every note's
text and on the commit after it, so it matches a note that *is* a closing
line and not one that only quotes the words.

- **No such line, or its commit is not an ancestor of `HEAD`**
  (`git merge-base --is-ancestor <commit> HEAD` exits non-zero — the branch
  was rebased): a first hunt. Read every one of those **modules whole** — the
  file, and where a file is one part of a unit (a Go package, a class split
  over files, a service with its store), the unit. A lost update or a
  swallowed error is rarely visible in the lines that changed; it is visible
  in how they meet the rest.
- **A line whose commit is an ancestor of `HEAD`:** a follow-up hunt. The
  files that line lists were read whole already, and reading them whole again
  only finds the next old defect, round after round. Of those files read only
  what changed since: `git diff <commit> -- <file>`, with enough of the
  surrounding function to follow it. A file in today's list that the line
  does not name is new to the change: read it whole. A follow-up hunt with no
  diff and no new file has nothing to read — say so in the closing line.

Given paths, there is no earlier round: always read them whole.

## What you look for

Go through all four, in every module, and say afterwards that you did:

1. **Concurrency** — two writers at once; read-modify-write without a lock, a
   transaction or a version check; a handler that runs while the state it
   reads is replaced; events that arrive while nobody listens and are lost.
2. **Error paths** — an error dropped, or logged and then reported as success;
   a step that fails halfway and leaves partial state; a cleanup or rollback
   that never runs; a retry that repeats a non-idempotent write.
3. **Unusual input** — empty, huge, duplicated, unordered, out of range,
   missing, another encoding, locale or time zone than the happy path assumes;
   data that lands under the wrong key, heading or owner.
4. **Permissions and boundaries** — who may call this, what it may touch, what
   it trusts that came from outside, what it exposes.

Correctness first. Style, naming and "would be nicer" are not findings here.

## Verify before you report

A finding you have not checked against the code is a guess, and a guess
written onto a ticket costs someone an hour to disprove. For each one:

- open the lines again and follow the call path that triggers it;
- where you can, make it happen — a scratch test, two commands in parallel
  (`cmd & cmd & wait`), a missing file — and check the exit code (`echo $?`),
  not the wording;
- drop it if the path cannot be reached, or if a guard elsewhere already
  holds. Say nothing about the ones you dropped.

What is left is a finding: `path:line`, what goes wrong, the input or timing
that makes it go wrong, and how you confirmed it (`read`, or `ran: <command>`).

**Name the class and the example.** The input that made it go wrong is one
example; the fault is a class of inputs or a pattern in the code. Say both —
`class: amounts with a thousands separator; example: "1.234,56" is read as
1.234`. A finding that names only the example gets a fix for only the example,
and the next round finds its neighbour. Where the same pattern sits in other
places, list them in the same finding: one class is one finding, however many
places carry it.

## Change or pre-existing

Mark every finding as one of two, and check it — do not guess:

- **`change`** — this change brought the defect in: a line on the path that
  triggers it was written by a commit that names the ticket id, or is not
  committed yet.
- **`pre-existing`** — the defect was there before: every line on that path,
  and the call that reaches it, is older than the ticket's commits, and the
  input that triggers it reached it before this change, too.

```bash
git blame -L <from>,<to> -- <file>               # who wrote the lines the defect rests on
git log -1 --format=%B <commit> | grep -i <id>  # does that commit name the ticket id?
```

A line shown as `Not Committed Yet` belongs to the change. Look at the caller
too: old code that a new call now reaches with input it never got before is a
`change`. Where the evidence does not settle it, mark it `change` and say in
the finding why it was unclear — a `pre-existing` mark is a promise that the
change did not do this.

Only a `change` finding is work on this ticket. A `pre-existing` one never
sends it back: whoever carries the ticket sorts it into a note for the person
or a ticket of its own, because they know what the ticket is about and you do
not.

## Where findings go

**Given a ticket id** — one note per finding on that ticket, and nothing else
on the board. The mark comes first:

```bash
jaira note <id> "bughunt: change: internal/store/save.go:88 - class: read-modify-write without a version check, also in archive.go:52; example: two parallel saves both read v3 and the second overwrites the first; ran: go test -run TestParallelSave -count=20 (fails 7/20)"
jaira note <id> "bughunt: pre-existing: internal/store/load.go:41 - class: a file without frontmatter is accepted as a ticket; example: an empty file is read as a ticket with no id; blame: 2f1c0a7 (before the ticket); ran: touch empty.md && jaira list (exit 0, ghost row)"
```

A finding never becomes a ticket of its own from you.

**Given paths only** — one ticket per finding, filed so someone who never met
you can act on it (`jaira create` with `--goal`, `--context`, `--dod`; the
context leads with what goes wrong and where). On a ticket's branch `follows`
is set by itself. There is no change to compare with, so no mark.

**Always last** — one closing line, findings or not, so whoever started you
knows you have finished and looked. It names the commit you hunted — a later
round changes the code, and a closing line from before it must not read as a
hunt of what ships now — and, after `files:`, every file of your list,
separated by spaces: the next round reads it back to know what it has seen.

```bash
jaira note <id> "bughunt: done at $(git rev-parse --short HEAD) - 2 findings (1 change, 1 pre-existing) in 5 modules, first hunt; concurrency, error paths, input, permissions checked; files: internal/store/save.go internal/store/load.go"
```

With paths only, the same line is your report instead of a note.

## Report

Three lines, to whoever started you:

- how many findings, how many of them `change`, and where they went (notes on
  `<id>`, or the new ticket ids)
- the modules you read whole, and the ones you read only as a diff since the
  last round
- anything you could not check, and why — a finding you could not confirm is
  this line, never a note
