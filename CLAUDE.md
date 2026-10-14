## Project

**jAIra**

A lightweight, git-synced kanban board for the tasks Claude Code generates while working in a repo. Tickets are markdown files committed inside the repo, so the board travels with the code — a teammate clones, runs `jaira`, and sees the same state with no server, no accounts, and no setup. Each lane is a configurable agent step (prompt + model tier + input/output contract), which turns the board from a tracker into a visible, resumable agent pipeline.

**Core Value:** **You never lose track of what a task was for or where it stands** — across session boundaries, agent runs, and teammates.

Everything else is negotiable. If the board can't answer "what was this ticket supposed to do, why, and what state is it in?" months later, it has failed.

### Constraints

- **Tech stack**: Single static binary, Go or Rust — no runtime dependency, because every teammate must be able to install and run it trivially
- **Storage**: Markdown with frontmatter, inside the repo — the file format *is* the API, and must stay hand-editable and diff-readable
- **Sync**: Git only — no server, so any feature requiring central coordination is off the table by construction
- **Performance**: Instant startup — the board is glanced at constantly, and a slow board will not be opened
- **Compatibility**: CLI must be usable by any bash-capable agent, not only Claude Code
- **Concurrency**: Multiple parallel Claude sessions and multiple teammates write the same store, so the format must minimize and survive conflicts
- **Scope discipline**: Every feature is measured against "is this smaller than paca?" — the project fails by growing



## Technology Stack

The stack decision and its research live in [docs/STACK.md](docs/STACK.md); the unabridged version is `.planning/research/STACK.md`.


<!-- jaira:start -->
## Task tracking: jaira

This repository has a jaira board (`.jaira/`). Multi-step work is tracked
there as markdown tickets so it survives session boundaries.

Capturing and picking work:

- `jaira create <title> --goal <...> --context <...> --dod <...>` — one call files a
  complete ticket; without a goal, a definition of done, the context it came from
  and an assignee it cannot leave the backlog
- the context is the only record of why a ticket exists. Write it for someone who
  was not in this conversation and reads it weeks from now: what is wrong today,
  what triggered it, what is already known or ruled out. Write it as if that
  reader has mild ADHD and knows none of what you know — lead with what is wrong,
  short concrete lines with one point each, names and paths rather than
  adjectives, no jargon and no preamble. It may span several lines, but someone
  should be able to act after the first two. If acting on it would need a
  question answered first, it is not finished
- a ticket you find while working another one follows on from it: on that
  ticket's branch (`feat/a3k9qp`) `jaira create` sets `follows` to it by itself and
  says so; `--follows none` when the new ticket is unrelated
- `jaira list --actionable --json` — everything that could be started right now
- `jaira next --json` — the single next actionable ticket
- `jaira tags` — the tags this board already uses, with how many open tickets
  carry each. **Read it before you tag anything** and reuse the name that is
  already there for that subject; never invent a synonym — "ui", "frontend" and
  "gui" on one board are three names for one thing and filter to nothing.
  `jaira tag <id> <name>...` adds tags, and `jaira create --tag <name>` sets them
  at capture. A name jaira has not seen is new and gets a colour in the
  hand-editable `.jaira/tags`; the board filter and `jaira list --tag <name>`
  read them back
- `jaira epic ls` — the features this board is built out of, each with done X of Y
  (follow-ups counted); an epic is a description plus a ticket list on its own
  ref, not a ticket. `jaira epic show <name>` lists its tickets, `jaira epic add
  <name> <id>...` sorts tickets in, `jaira epic ls --unassigned` lists the rest

Working a ticket:

- `jaira claim <id>` — take it first; other sessions read this board too
- `jaira show <id> --for-lane <lane> --json` — the lane's prompt, the bounded input,
  the model tier, the outputs the lane expects back, and `mode`. Empty — the
  usual case — means the lanes run autonomously. `conversational` means a person
  found open design decisions on this ticket and is reading along: show the diff
  after every definition-of-done item instead of at the end, and, where the lane
  changed code, hand back a ready-to-paste `git add`/`git commit` line with the
  ticket id in its subject rather than committing yourself. Read it off the
  ticket, never off the line that started you — that is what makes it survive a
  killed session
- `jaira dod <id> <n> --done --proof "<file:line or test name>"` — tick an
  acceptance criterion and say what makes it true. These are what the terminal
  lane's gate reads: one left unticked refuses the move however finished the
  work is, and it refuses it at the end, when the cost of finding out is
  highest
- `jaira dod <id> <n> --doing|--done --plan` — the method, a second and
  separate list. Ticking the plan does not tick the definition of done
- `jaira note <id> <text>` — at every pause, write down what the repository does
  not already say: dead ends, why this and not that, what you had to find out.
  Not what the checklist and git already record. A killed session never gets a
  turn to write anything down, so do not save it for the end
- `jaira move <id> --to <lane> --what <...> --why <...> --resolves <...>` — finish
  the step. The three outcome flags belong to the lane that produces the outcome
  (in-progress); every later lane reads it, moves without them, and is refused
  when it would replace it. jaira works out the commit list itself from git
  history — the union of the ticket file's own history and commits naming its
  id — so nothing needs to be typed here; it is written onto the ticket once,
  when the ticket leaves the board
- `jaira resume` — work left in progress, with everything recorded about it
- **every commit names the ticket id** — `fix(A3K9QP): ...`. jaira derives a
  ticket's commit list from two sources: the history of the ticket file itself,
  and the commits that name its id. The first source is thin by design (see the
  point on lanes that change no code), and on a board that has not been shared
  yet (`jaira init` gitignores `.jaira/` until `jaira share`) it does not exist
  at all — so the handle in the subject line is what the list is actually built
  from. Leave it out and the list stays empty and the move into the last lane
  is refused
- **the ticket rides with the code, never on its own.** Move the ticket first,
  then `git add` the changed file under `.jaira/tickets/` alongside your source
  changes and commit them together. A reviewer then sees the change and what it
  was for in one place, instead of a diff whose ticket is still in whatever
  state the last commit left it
- **a lane that changed no code commits nothing.** critique, testing and review
  leave a note and a lane change and no source change; a commit carrying only
  that is bookkeeping, and a branch with one of them per lane hides the work
  inside it. Leave the ticket file modified in the worktree — the next commit
  that carries code takes it along; if no further code commit follows, the
  commit that files the ticket away with `jaira logbook <id>` carries its
  final state. Nothing is lost by waiting: the lane's writes are already on
  the ticket. A ticket you create commits nothing either: it travels on its
  ref, everybody sees it there, and `jaira pull <id>` brings it into the
  branch that works it. Only when `jaira create` reports it as a file on your
  disk — the board has no usable remote — commit it, or nobody but you knows
  it exists
- `jaira logbook <id>` — once a ticket reaches the terminal lane, stamps its
  commits and files it under `.jaira/logbook/<you>-<date>/`, taking it off the
  board. `jaira restore <file>` brings it back

`jaira <command> --help` for everything else.

Do not edit files under `.jaira/tickets/` directly; the CLI is the write path.
An agent cannot move work on out of the human review lane — a person accepts it there.
Sending it back to an earlier lane is not an acceptance and needs no --force.

A `jaira:local` marker — an HTML comment of that name — added by hand anywhere
inside this block makes everything between it and the end marker survive the
next regeneration. Nothing writes it for you, and there is none here until
somebody adds one; project-specific rules belong behind it rather than fighting
this note from outside the block.

## This board's lanes

Order: backlog → brainstorm → todo → pre-process → in-progress → critique → optimize → testing → human → review → signoff → done → blocked
Loop: critique sends work back to in-progress, and that repeats until critique has nothing left to say.
Loop: optimize sends work back to in-progress, and that repeats until optimize has nothing left to say.
Loop: testing sends work back to in-progress, and that repeats until testing has nothing left to say.

- `backlog` — no agent step; move through it
  Captured but not yet specified enough to work on.
- `brainstorm` — yours to work; tier strong; must produce goal
  Working out what the ticket should even be.
- `todo` — no agent step; move through it
  Specified and ready to be picked up.
- `pre-process` — yours to work; tier strong; must produce plan
  Working out how the change will be made.
- `in-progress` — yours to work; tier cheap; must produce outcome-what, outcome-why, outcome-resolves
  Carrying out the plan.
- `critique` — yours to work; tier strong; must produce review-summary
  Judges whether this is the right implementation, not whether it works.
- `optimize` — yours to work; tier strong; must produce review-gaps
  Removes what the change does not need — code that already exists elsewhere, code nobody calls, and code that carries its weight in nothing.
- `testing` — yours to work; tier cheap; must produce test-verdict
  Runs the change and checks it against the ticket - does the demanded thing exist, and does it work.
- `human` — **a person's, not yours** — you may move work in, never out; only with an open question (`--question`); without one, work passes it by
  Human in the loop.
- `review` — **fresh context** — an agent's, but not the one that carried the ticket here: hand it to a new agent; tier strong; must produce review-summary, review-gaps, review-verdict, review-check
  A second model has judged the diff.
- `signoff` — **a person's, not yours** — you may move work in, or send it back to an earlier lane, never on
  Reviewed by a model, waiting for a person to accept it or send it back.
- `done` — no agent step; move through it; terminal
  Accepted.
- `blocked` — no agent step; move through it; parking: work returns to the lane it left
  Waiting on an external dependency.

Nothing moves a ticket for you. There is no daemon and no runner: a lane's
prompt runs because a session ran it. Drive the board from the session you are
in — `jaira next --per-lane --json` says which lanes have work waiting and
which of them you are allowed to work, `jaira show <id> --for-lane <lane> --json`
hands you that lane's prompt and its bounded input, and `jaira move` puts the
result in the next lane. Work one lane to empty before starting the next, or
the lane nobody drives is the one that fills up.

Told to start or work a ticket, drive it this way yourself — lane by lane,
loops included — until it sits in a human lane, then continue once the human
has answered. The next lane is `next_lane` in `jaira show <id> --json`: it
passes a lane that only takes a question by, so a ticket with nothing to
decide goes on to the model review and reaches a person once, at the end.
Told an agent should work it, hand it to a subagent that babysits the ticket
through the same route.

<!-- jaira:local -->

## Work rides on a branch and arrives through a PR

Nothing lands on `master` directly. A change is made on a branch of its own, the
ticket rides in the same commits as the code it belongs to, and the branch
reaches `master` through a pull request — which is what makes the diff readable:
the change and what it was for in one place. It rides with the code and never
alone: a lane that changed no code — critique, testing, review — commits
nothing at all and leaves the ticket file for the next commit that carries
code, or, when no further code commit follows, for the commit that files the
ticket away with `jaira logbook`.

**The pull request belongs to the maintainer from the moment it exists.** An
agent pushes its branch and stops there: it does not open the pull request, does
not merge one, and does not approve one. Opening it is the maintainer's call,
given as an instruction — so the person who wrote the change is never the person
who decides it arrives. This is the same rule the board already applies to lanes
— a review agent cannot certify its own work — said once more where it also has
to hold.
<!-- jaira:end -->

## Client-facing changes go in core/release/notes/

Every change a user of jaira can notice gets a note in
`core/release/notes/<TICKET>.md` — the ticket's six-character handle, so
`core/release/notes/AR97K9.md` — in the same commit as the change itself. Not a
changelog generated from git later — a line written while you still know why it
matters. The notes are embedded into the binary with go:embed and are what
`jaira update` reads back to someone whose board was last touched by an older
build, so an unwritten change is a change nobody is ever told about.

One file per ticket is what keeps two pull requests from ever conflicting over
release notes: each adds a file nobody else touches. A ticket with a second
change adds a second line to its own file. A file directly in
`core/release/notes/` is Unreleased; once a release has moved it away, a
further change for the same ticket creates `core/release/notes/<TICKET>.md`
again.

Counts as client-facing: a new or renamed command or flag, changed CLI output or
exit codes, a new or changed field in the ticket frontmatter, anything that
alters the TUI, a changed default, and anything that makes an existing
invocation behave differently.

Does not count: refactors, tests, internal helpers, and anything a user cannot
observe from outside the binary.

The format is a line scan, and breaking it silently corrupts what users read:

- One change is exactly one line starting `- `. Never wrap it across two lines —
  a wrapped line is read as two separate, incomplete changes.
- Everything else in the file — prose, blank lines, HTML comments — is ignored.
- No version heading. A note file never says which release it belongs to.

Write each line as an instruction: what the reader must DO differently, not a
record of what a commit did.

Releasing is a button, and nothing in it lands on master directly:

1. Run the `release-notes` workflow from the Actions tab on master with the
   version, `X.Y.Z`. It moves every file directly in `core/release/notes/` into
   `core/release/notes/<version>/` on a branch `release/<version>` and opens a
   pull request.
2. Merging that pull request runs `release-tag`: it tags the merge commit
   `v<version>` and starts `release.yaml`, which builds the binaries as before.

A release is tested once, in step 1: the button runs `go test ./...` on the cut
tree before it opens the pull request. The pull request and its merge change
only `core/release/notes/**`, which `ci.yaml` skips, and the tag build runs no
tests. A release pull request that changes anything else gets the full CI run.

The binary reads a release's version from its folder name. Never move a note
into a version folder by hand, and never edit a note inside one: that release
has shipped.

`core/release/NOTES.md` holds the releases up to 0.3.6 and is closed history:
nothing is added to it.
