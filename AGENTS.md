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

Working a ticket:

- `jaira claim <id>` — take it first; other sessions read this board too
- `jaira show <id> --for-lane <lane> --json` — the lane's prompt, the bounded input,
  the model tier, and the outputs the lane expects back
- `jaira dod <id> <n> --doing|--done` — mark checklist items as you go
- `jaira note <id> <text>` — at every pause, write down what the repository does
  not already say: dead ends, why this and not that, what you had to find out.
  Not what the checklist and git already record. A killed session never gets a
  turn to write anything down, so do not save it for the end
- `jaira move <id> --to <lane> --what <...> --why <...> --resolves <...>` — finish
  the step. jaira works out the commit list itself from git history — the union
  of the ticket file's own history and commits naming its id — so nothing needs
  to be typed here; it is written onto the ticket once, when the ticket leaves
  the board
- `jaira resume` — work left in progress, with everything recorded about it
- on a board that has not been shared yet (`jaira init` gitignores `.jaira/`
  until `jaira share`), the ticket file is untracked, so the only thing tying a
  commit to a ticket is its handle in the commit message. Name it there —
  `fix(A3K9QP): ...` — or the derived list stays empty and the move is refused
- **the ticket rides in the same commit as the code.** Move the ticket first,
  then `git add` the changed file under `.jaira/tickets/` alongside your source
  changes and commit them together. A reviewer then sees the change and what it
  was for in one place, instead of a diff whose ticket is still in whatever
  state the last commit left it. Same for a ticket you create and hand to
  someone else: commit it, or nobody but you knows it exists — and now this is
  also what makes the commit list derivable at all: that shared commit is how
  git ties the ticket to the change
- `jaira logbook <id>` — once a ticket reaches the terminal lane, stamps its
  commits and files it under `.jaira/logbook/<you>-<date>/`, taking it off the
  board. `jaira restore <file>` brings it back

`jaira <command> --help` for everything else.

Do not edit files under `.jaira/tickets/` directly; the CLI is the write path.
The human review lane cannot be left by an agent — a person accepts the work there.

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
- `human` — **a person's, not yours** — you may move work in, never out
  Human in the loop.
- `review` — yours to work; tier strong; must produce review-summary, review-gaps, review-verdict, review-check
  A second model has judged the diff.
- `signoff` — **a person's, not yours** — you may move work in, never out
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
has answered. Told an agent should work it, hand it to a subagent that
babysits the ticket through the same route.
<!-- jaira:local -->

## Client-facing changes go in core/release/NOTES.md

Every change a user of jaira can notice gets a line in `core/release/NOTES.md`,
in the same commit as the change itself. Not a changelog generated from git
later — a line written while you still know why it matters. The file is
embedded into the binary with go:embed and is what `jaira update` reads back to
someone whose board was last touched by an older build, so an unwritten change
is a change nobody is ever told about.

Counts as client-facing: a new or renamed command or flag, changed CLI output or
exit codes, a new or changed field in the ticket frontmatter, anything that
alters the TUI, a changed default, and anything that makes an existing
invocation behave differently.

Does not count: refactors, tests, internal helpers, and anything a user cannot
observe from outside the binary.

The format is a line scan, and breaking it silently corrupts what users read:

- One change is exactly one line starting `- `. Never wrap it across two lines —
  a wrapped line is read as two separate, incomplete changes.
- A section heading is `## ` followed by the version a build reports: the git
  tag with its leading `v` stripped, so tag `v0.1.0` is `## 0.1.0`.
- Newest section first. Everything else — prose, blank lines, HTML comments — is
  ignored, so explaining yourself in the file is free.

Write each line as an instruction: what the reader must DO differently, not a
record of what a commit did.

New entries go under `## Unreleased` at the top. A version heading is written
when that version is tagged, never before: once a tag exists for a version, that
section is closed history and nothing may be added to it — a released note that
grows after the fact describes a binary nobody has. Releasing means renaming
`## Unreleased` to the version number being tagged and opening a fresh empty
`## Unreleased` above it. If no `## Unreleased` section is there when you need
one, add it.

<!-- jaira:end -->
