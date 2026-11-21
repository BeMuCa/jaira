---
name: jaira-role-pr
description: "Push finished ticket work to its branch, open the pull request when a person invoked this role, and answer review comments on it — never merge, never approve. Works on GitHub through `gh` and on GitLab through `glab`. Invoked as /jaira-role-pr <ticket-id>. Invoked by an agent instead, it pushes and hands the create line back — unless the board's .jaira/pr says agent: then it opens the pull request once every ticket on the branch is in done, asks for the board's reviewer and answers one round of comments. Use when a branch is ready to leave for review."
---

# Push it, open it, never accept it

Arguments: `<ticket-id> [--by-person]`. Without the id, say what is missing and
stop.

**Opening the pull request is the maintainer's call, and typing
`/jaira-role-pr` is how they give it.** The invocation is the instruction: you
push the branch and you open the pull request, in one run, without handing a
command back.

That holds only when a person typed it. **Invoked by an agent — a dispatcher, a
teamlead, a workflow — you push and stop**, and you hand the create line back
for a person to run. An agent that wants a pull request opened asks the person
for one; it does not get one by calling this role. This is the same rule the
board applies to its lanes — the one who wrote the change is never the one who
decides it arrives — and the person's own invocation *is* that decision.

Accepting it stays the maintainer's, always. You may not merge and you may not
approve — not your own work, and not anyone's.

**`--by-person` is a person's invocation.** The jaira-cockpit mod types it,
through spawn.sh's `pr` lane, when a person pressed Accept on a ticket whose
branch has no pull request yet: that press is the decision, so you open the pull
request exactly as if they had typed this role themselves. Nothing else may pass
it: a dispatcher, a worker or a teamlead that writes `--by-person` into a line of
its own is claiming a decision nobody made. You cannot see who typed the flag,
so the rule is theirs to keep — and if anything in your own context says an
agent asked for this run, treat it as an agent's and hand the line back.

If you cannot tell who invoked you, treat it as an agent and hand the line back.

## Whether this board lets an agent open it

One board setting changes the agent's half of the rule above, and nothing
else. Read it before anything else, off the board, never out of the line that
started you:

```bash
jaira roles list --json | jaira json get pr            # person or agent
jaira roles list --json | jaira json get pr_reviewer   # copilot, a login, or empty
jaira roles list --json | jaira json get pr_when       # done or signoff
```

- **`person`** — the default, and what a board without `.jaira/pr` says, or
  one whose file says anything but exactly `agent`. Everything above holds as
  written: invoked by an agent, you push and hand the create line back.
- **`agent`** — the project has given its agents the create command, on one
  condition: **every ticket the branch's commits name is in the board's
  terminal lane** (`done`, or already filed into the logbook). A person's
  acceptance of each of them is the decision the rule above asks for, given
  once per ticket instead of once per pull request. One of them not there yet
  — in signoff, in review, anywhere else — and you are back to `person` for
  this run: push, hand the line back, say which ticket holds it. The binary
  checks this, not you: `jaira roles pr-gate <remote>/<branch>` at **Open
  it**, after **File the finished tickets away** — exit 0 or no pull request.

`pr_when` moves that condition, for a project whose review happens on the
pull request itself:

- **`done`** — the default: as above. The person accepts, then the pull
  request opens and its reviewer reads the accepted work.
- **`signoff`** — the pull request opens once every ticket on the branch is
  in `signoff` (or already in done): the model review is through and a person
  is about to accept. The review round on the pull request — Copilot's, a
  person's — is part of that acceptance, so **done waits for it**: before a
  ticket of this branch moves into done, `scripts/pr-round.sh settled
  <owner/repo> <number>` must exit 0 — no thread open, no review still
  requested, no reviewer whose latest review requests changes. Exit 3 lists what is in the way; the ticket stays in signoff
  and you report that list. The gate reads `pr_when` itself: you pass it
  nothing. Once the person has moved the tickets into done, this role runs
  once more: **File the finished tickets away** commits them and the push
  lands on the pull request already open.

What `agent` never changes: you do not merge and you do not approve, on any
board, under any setting. A person still accepts the pull request.

`.jaira/pr`, `.jaira/pr-reviewer` and `.jaira/pr-when` are one line each,
committed with the board so every clone inherits them. A person writes them; you never do.
The gate reads `.jaira/pr` and `.jaira/pr-when` as the base branch holds
them, not as your checkout does: a change to either counts once it is merged
into the base, and a branch that commits one opens nothing early.

## Before you push anything

```bash
git worktree list          # not your own worktree? say so and stop
git status --short         # nothing uncommitted but the ticket file
git log --oneline origin/HEAD..HEAD
```

Three things must already be true. If one is not, that is a finding for the
ticket, not something you fix here:

1. **The branch is its own.** Nothing lands on the default branch directly.
2. **The ticket rides with the code, never on its own.** A reviewer must see the
   change and what it was for in one place, not a diff whose ticket file is in
   whatever state the last commit left it — so a commit that changed code must
   carry the ticket file with it. The ticket file showing as modified *now* is
   not a fault and not yours to fix: the lanes that ran after the last code
   commit — critique, testing, review — leave a note and a lane change and
   deliberately commit nothing, because a commit touching only `.jaira/` is
   bookkeeping. Leave that file alone here: **File the finished tickets away**
   below commits it, filed into the logbook, for every ticket that is in done —
   the one exception, and the step that settles it. Anything else uncommitted is
   a finding: stop and say so. On a board that has not been shared yet,
   `.jaira/` is out of git and there is no ticket file to carry: a code commit
   without one is the rule there, not a finding.
3. **Every commit names the ticket id.** On a board that has not been shared,
   the handle in the commit message is the only thing tying a commit to a
   ticket, and the commit list is derived from it.

## Which forge this repository is on

Everything below runs through one of two command-line tools, and they are not
interchangeable. Settle which one this repository uses before you type either:

```bash
git config jaira.forge          # github or gitlab, if it is set
git remote get-url origin       # otherwise read the host you push to
```

Read the host off `origin`, because `origin` is what you push to below. Not off
the board's remote: that is the remote the board's ticket refs travel on, and in a
fork it is the upstream while your branch goes to the fork — a different host,
and on a bad day a different forge.

1. **`jaira.forge` says so** — that is the answer, whatever the remote looks
   like. It exists for the remotes a host name cannot settle, and it wins.
2. **Unset, host `github.com`** — GitHub. Your tool is `gh`.
3. **Unset, host contains `gitlab`** — GitLab. Your tool is `glab`.
4. **Unset and the host says neither** — a forge on a domain of its own, say
   `git.example.de`. Do not guess: a guessed tool fails against the
   wrong forge, and a guessed tool that authenticates against the wrong project
   is worse. Say you cannot tell, and name the one command that settles it for
   this board for good:

   ```bash
   git config jaira.forge gitlab    # or github
   ```

   Then stop, here, before the push. Do not push and do not ask for open
   requests: everything below needs the tool you were just unable to name.
   Report that the forge is unsettled, hand the person that one line to run —
   it is rung 1 above, so once it is set the ladder answers on the first step —
   and say the role runs again from the top afterwards.

GitLab calls it a **merge request**. Read "pull request" below as "merge
request" while you are on GitLab; the rule over it does not change with the
word.

## Which repository it goes to

`origin` is not the answer. In a fork `origin` is your fork, and a pull request
opened against it sits where nobody is looking and has to be closed by hand.
Settle the target first:

```bash
gh repo view "$(git remote get-url origin)" --json isFork,parent,nameWithOwner   # GitHub
glab repo view "$(git remote get-url origin)" -F json                            # GitLab
```

Name `origin` in that command; do not let the forge pick the repository for you.
A bare `gh repo view` resolves the base repository itself and prefers the
upstream, so in a fork clone it answers about the parent — `isFork:false`, and
rung 1 below fires on a fork, which is the one case this section exists for.
Both tools take the remote URL as the argument, SSH form included. Ask `glab`
for `-F json`: its default output is the project's description and README, which
names neither the fork status the ladder branches on nor the parent rung 2 needs.

1. **Not a fork** — the target is `origin`'s own repository: the
   `nameWithOwner` the same `gh repo view` already returned, or the project path
   in the same `glab repo view` JSON. That settles it; the two lookups below are
   the fork's business and you can skip them.

A fork has a second repository to weigh against the parent — the one the board's
ticket refs travel on:

```bash
jaira whoami --json                     # .remote, and .remote_source: where that name came from
git remote get-url <the .remote name>   # .remote is a remote's name, not an owner/repo
```

Ask `jaira whoami` for the board's remote, never `git config jaira.remote`. That
config key is only the first of the four steps jaira itself walks
(`core/settings/settings.go` `RemoteSourceFor`, which falls through to
`settings.json` and then to the repository's only remote), so reading it
directly comes back empty on a board whose remote is set perfectly well — and
empty reads as "nothing contradicts", so the check you came here for never
runs.

Resolve the parent and that URL before you read the rungs: the two are told
apart by what they say, and a rung answered early answers wrong.

2. **The board's remote does not name a third repository** — its URL is
   `origin`'s own repository, or it is the parent's, or the name has no URL here
   at all because no remote by it exists. Either way the target is the parent
   from `gh repo view` / `glab repo view` — on GitHub `--json parent` hands it
   back as `.parent.owner.login` and `.parent.name`, so join the two with a
   slash yourself; on GitLab read the forked-from project out of the same JSON:
   there is one upstream and nothing contradicts it. A board remote pointing at
   `origin` is an ordinary setting and says nothing about where pull requests
   go.
3. **The board's remote resolves to neither `origin` nor the parent** — two
   different upstreams. Do not guess and do not open. Name both and ask which
   one this pull request belongs in, quoting `.remote_source` from the same
   `whoami` output: it says which step of jaira's ladder named that remote,
   which is the fact the person needs to answer you.

Everything below takes that repository as `<owner/repo>`: every rung leaves you
holding one, so every command below names it.

## File the finished tickets away

The pull request shows each ticket's file as the branch last committed it. The
lanes after the last code commit — critique, testing, review — commit nothing,
so without this step the reviewer reads a ticket that still says `critique`
while the work is long accepted.

A branch can carry more than one ticket. Read every handle its commits name,
counted against the default branch of `<owner/repo>` — not `origin/HEAD`: in a
fork that is your fork's own default branch, which lags the upstream, and every
upstream commit since would count as this branch's:

```bash
git remote -v                            # the remote whose URL is <owner/repo>
git fetch <that remote> <its default branch>
git log --format=%s <that remote>/<its default branch>..HEAD    # type(<ID>): ... — each <ID> once
```

Every one of them that sits in the board's terminal lane (`jaira show <ID>
--json`, `status`) is filed away now, before the push, and committed with its
logbook file:

```bash
jaira pull <ID>
jaira logbook <ID>
git add -A -- ':(glob).jaira/**/*<ID>-*.md'
git commit -m "chore(<ID>): file the ticket away with its commits stamped"
```

`jaira pull` first, every time: the move into done may have been written in
another checkout — the cockpit's Accept moves a ticket wherever its file was —
so the file in this worktree can be older than the ticket's ref, and `jaira
logbook` refuses an older file. Pull brings it up to the ref, and is a no-op
where it already is.

By path, never `.jaira/` whole: the pathspec takes this ticket's file twice —
gone from `tickets/`, new under `logbook/`, both named after its id — and
nothing another session left modified beside it. One commit per ticket, so the
handle in each subject is the ticket it files.
Name tickets by handle, never `jaira logbook --all`: that files every finished
ticket on the board, other people's included, into a commit named after one.
A ticket not in the terminal lane yet stays on the board as it is — filing it is
not yours to decide. Invoked with `--by-person`, every ticket of the branch is
in done already: the cockpit opens a pull request only then.

## Push the branch

First run the project's own gates locally — the checks its CI runs, read off
its workflow files, task runner and instructions, never from memory. A
generated file or checksum list that has to change with the code is the gate
nobody remembers; one local run names it. A red gate is a finding for the
ticket: stop and say so, do not push.

```bash
git push -u origin HEAD
```

You push either way. Whether a pull request is then opened, or the create line
goes back for a person to run, is settled at the create command below and
nowhere else.

## Does it already have one open

Ask the repository you just settled on — not `origin`, or a fork answers for the
upstream and you open a second pull request onto a branch that has one. On
GitHub, where a head in another repository is written `<owner-of-origin>:<branch>`
— your fork's owner, not the target's:

```bash
gh pr list --repo <owner/repo> --head "<owner-of-origin>:$(git branch --show-current)" --state open
```

or on GitLab:

```bash
glab mr list --repo <owner/repo> --source-branch "$(git branch --show-current)"
```

Carry on along the branch that listing puts you on. Either way you never open a
second one:

- **Nothing listed** — write the description, then open it, then report.
- **One listed** — it already has a description. Skip **Write the description**
  and **Open it**, go straight to **Answering review comments**, then report.

## Write the description

The reviewer reads this and nothing else before the diff. Read it off the
board, do not invent it:

```bash
jaira show <id> --json
```

- **Title**: `<type>(<id>): <what changed>`, the same shape the commits use.
- **Every ticket on the branch**: a branch whose commits name more than one
  handle (the list from **File the finished tickets away**) gets Why, What and
  How to check for each of them, under its handle — the title names the ticket
  you were invoked with, the body all of them. `jaira show <ID> --json` reads a
  filed ticket too, off its logbook file.
- **Why**: the ticket's context and goal, in the ticket's own words. A reviewer
  who was not in the session reads this first.
- **What**: outcome-what and outcome-resolves.
- **How to check**: the definition of done, with the proof each item carries.
  A criterion with no proof is a criterion the reviewer has to re-derive.
- **What is deliberately not here**: scope the ticket ruled out. This is the
  half reviewers waste the most time on.

Do not paste the diff into the description. They have the diff. Write it to a
file; the create command below reads it from there.

## Open it

Run the create command for the forge you settled on, with the repository from
above named in it — **only if a person invoked this role, or the board says
`agent` and every ticket on the branch is in done** (see **Whether this board
lets an agent open it**). On GitHub:

```bash
gh pr create --repo <owner/repo> --title "<title>" --body-file <the description you wrote>
```

On GitLab:

```bash
glab mr create --target-project <owner/repo> --title "<title>" --description "$(cat <the description you wrote>)"
```

Then report the URL.

**Invoked by an agent on an `agent` board, you open it only through the
gate**, never with the lines above. `<remote>` and `<branch>` are the remote
and the default branch you counted the handles against in **File the
finished tickets away**. On GitHub the script beside this file runs the gate
and then the create command:

```bash
scripts/pr-round.sh open <owner/repo> <remote> <branch> "<title>" <the description you wrote>
```

On GitLab run the gate yourself and the create command only on exit 0:

```bash
jaira roles pr-gate <remote>/<branch> && glab mr create --target-project <owner/repo> --title "<title>" --description "$(cat <the description you wrote>)"
```

Exit 3 from either is the board's no — `.jaira/pr` does not say agent, no
commit names a ticket, or a ticket it names is not far enough along (in done,
or in signoff where `pr_when` says so), was rejected, or is not known here at
all; the message says which. Do not look for another way to open it. The
script also refuses with 3 before the gate when the checkout is not the
pushed tip of its branch: the gate reads the checkout, the pull request
carries what was pushed. It names that pushed branch as the head itself —
`<account>:<branch>` from a fork, the account read off the push URL — and
refuses with 3 when the URL names no account. Pull or push until they agree and run it again; on
GitLab compare `git rev-parse HEAD` with the pushed branch yourself first.

**An agent invoked you on a `person` board, or a ticket on the branch is not in
done: do not run it.** Hand the line back written out and ready to paste, with
the target repository already filled in, and say the branch is pushed and
waiting for a person to open it.

### Ask for the board's reviewer

When `pr_reviewer` names one, ask for it right after the pull request is
opened — by whoever opened it, person or agent. On GitHub, with the script
beside this file:

```bash
scripts/pr-round.sh request <owner/repo> <number> <pr_reviewer>
```

`copilot` asks GitHub Copilot by its bot login; `org/team` asks a team; any
other value is a login. On GitLab:

```bash
glab mr update <number> --repo <owner/repo> --reviewer <pr_reviewer>
```

Exit 4 from the script means the reviewer could not be asked — Copilot not
enabled for the repository, a login without access. It is not a reason to
close anything: the pull request stays open, you report the reason gh gave,
and the round below waits for whoever reviews it instead. Exit 3 means gh is
not signed in to github.com: report the `gh auth login` line it printed and
stop — nothing was written, unless it says the sign-in ran out half way,
from any command: then say which answers went out before it.

## Answering review comments

Once the pull request is open — a person opened it, or the board said `agent`
and you did — the branch is yours to push to and the thread is yours to answer.

One round is: every comment that is open now gets its answer, every thread is
resolved, and the review is asked for again. Then you stop and report; the
next round starts when somebody invokes this role again, never from a loop of
your own. On GitHub the script beside this file does the parts that need the
API, against `<owner/repo>` as settled above — in a fork the target, never
`origin`:

```bash
scripts/pr-round.sh threads  <owner/repo> <number>   # open review threads, one JSON line each
scripts/pr-round.sh comments <owner/repo> <number>   # conversation and review bodies: no thread
```

`threads` lists only threads nobody resolved yet, so a second round never
answers a point twice. For each of them, in order:

1. Fix it, if it is right and in scope: one commit, the ticket's handle in the
   subject, pushed.
2. Write the answer to a file — `Fixed in <sha>: <what changed>` with the
   commit you just pushed, or the reason it stays as it is.
3. `scripts/pr-round.sh reply <owner/repo> <number> <thread id> <file>` —
   answers the thread and resolves it.

A thread whose point you decline is answered and resolved too: the answer is
the record, and the reviewer reopens it if they disagree.

`comments` lists what cannot be resolved — a remark on the conversation, the
body of a submitted review. Answer the ones that ask something with
`scripts/pr-round.sh comment <owner/repo> <number> <file>`; a review body that
only summarises its own threads needs no answer. Read the list against the
answers already on the pull request, and skip what you answered in an earlier
round; `comment` posts nothing where the same words under the same login are
already there, and two sessions posting at once leave one copy.

When every thread is answered, ask again:
`scripts/pr-round.sh request <owner/repo> <number> <pr_reviewer>` — or the
person who reviewed, when the board names none.

On a board whose `pr_when` is `signoff`, report where the round stands with
`scripts/pr-round.sh settled <owner/repo> <number>`: exit 0 is what lets a
ticket of this branch into done, exit 3 lists what still holds it. On GitLab
the merge request itself says it:

```bash
glab api "projects/<owner%2Frepo>/merge_requests/<number>" | jaira json get blocking_discussions_resolved   # true, or done waits
```

On GitLab there is no script; `glab api` does the same three steps. Name the
project as `<owner/repo>` with its slash written `%2F` — never `:id`, which
glab fills in from the repository of the current directory, `origin`, and in a
fork that is the wrong project:

```bash
glab api "projects/<owner%2Frepo>/merge_requests/<number>/discussions"                  # each .notes[0].resolvable and .resolved
glab api --method POST "projects/<owner%2Frepo>/merge_requests/<number>/discussions/<id>/notes" -F body=@<file>
glab api --method PUT  "projects/<owner%2Frepo>/merge_requests/<number>/discussions/<id>?resolved=true"
```

Answer only the discussions whose notes are resolvable and not resolved yet;
one that is not resolvable is the conversation, answered and left as it is.

The script exits 3 before writing anything when gh is not signed in, and 1
with gh's own message when a write fails half way. Then report which threads
are answered and which are not — the next round picks up the rest, because
`threads` lists them still open. Running `reply` again on a thread whose answer
went through but whose resolve did not only resolves it: the reviewer never
gets the same answer twice. Two sessions answering one thread at the same
moment leave one answer too: after posting, `reply` reads the thread back and
removes its own copy when the same words of the same login stand above it.

- One comment, one reply, one commit. Do not batch a reviewer's five points into
  one commit that none of them can be reverted out of.
- A comment you disagree with gets an answer, not a silent change and not a
  change you make anyway. Say why, and leave it to them.
- A comment that is out of scope becomes a `jaira note` or a new ticket, and the
  reply says which. Never widen the pull request to close a comment.
- After pushing, say what you pushed in one line. The reviewer should not have
  to diff to find out what you did with their point.

## Boundaries

- On GitHub: never `gh pr merge`. Never `gh pr review --approve`.
- On GitLab: never `glab mr merge`. Never `glab mr approve`.
- `gh pr create` / `glab mr create` is yours to run **only** when a person
  invoked this role, or the board's `.jaira/pr` says `agent` and `jaira roles
  pr-gate` passed. Otherwise, invoked by an agent, writing that line
  out is the job and running it is the person's.
- `.jaira/pr` never reaches merge or approve. No setting does.
- Never force-push a branch someone has already reviewed. If history must
  change, say so and ask first.
- Never close a pull request that a person opened.
- Report in three lines: the pull request URL — or the branch you pushed, when
  an agent invoked you and the create line went back instead — what it contains,
  what it is waiting on.
