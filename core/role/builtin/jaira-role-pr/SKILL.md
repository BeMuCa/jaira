---
name: jaira-role-pr
description: "Push finished ticket work to its branch, open the pull request when a person invoked this role, and answer review comments on it — never merge, never approve. Works on GitHub through `gh` and on GitLab through `glab`. Invoked as /jaira-role-pr <ticket-id>. Invoked by an agent instead, it pushes and hands the create line back. Use when a branch is ready to leave for review."
---

# Push it, open it, never accept it

Arguments: `<ticket-id>`. Without it, say what is missing and stop.

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

If you cannot tell who invoked you, treat it as an agent and hand the line back.

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
   bookkeeping. Push what is committed and leave that file alone. Anything else
   uncommitted is a finding: stop and say so.
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
   `git.esprit-engineering.de`. Do not guess: a guessed tool fails against the
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

## Push the branch

```bash
git push -u origin HEAD
```

That is the whole of this step: whether a pull request then gets opened, or the
create line goes back for a person to run, is settled at the create command
below and nowhere else.

## Which repository it goes to

`origin` is not the answer, and you need the answer before the next command, not
before the last one. In a fork `origin` is your fork, and a pull request opened
against it sits where nobody is looking and has to be closed by hand. Settle the
target first:

```bash
gh repo view "$(git remote get-url origin)" --json isFork,parent,nameWithOwner   # GitHub
glab repo view "$(git remote get-url origin)"                                    # GitLab
jaira whoami --json                              # the board's remote, and where it came from
```

Name `origin` in that command; do not let the forge pick the repository for you.
A bare `gh repo view` resolves the base repository itself and prefers the
upstream, so in a fork clone it answers about the parent — `isFork:false`, and
rung 1 below fires on a fork, which is the one case this section exists for.
Both tools take the remote URL as the argument, SSH form included.

Ask `jaira whoami` for the board's remote, never `git config jaira.remote`. That
config key is only the first of the four steps jaira itself walks
(`core/settings/settings.go` `RemoteSourceFor`, which falls through to
`settings.json` and then to the repository's only remote), so reading it
directly comes back empty on a board whose remote is set perfectly well — and
empty reads as "nothing contradicts", so the check you came here for never
runs.

`.remote` in that JSON is the *name* of a remote — `upstream`, say — not an
`owner/repo`. Turn it into one before you use it:

```bash
git remote get-url <the .remote name>
```

Resolve both the parent and that URL before you read the ladder: the rungs are
told apart by what the two say, and a rung answered early answers wrong.

1. **Not a fork** — the target is `origin`'s own repository: the
   `nameWithOwner` the same `gh repo view` already returned, or the path
   `glab repo view` printed.
2. **A fork, and the board's remote does not name a third repository** — its URL
   is `origin`'s own repository, or it is the parent's, or the name has no URL
   here at all because no remote by it exists. Either way the target is
   the parent from `gh repo view` / `glab repo view` — `--json parent` hands it
   back as `.parent.owner.login` and `.parent.name`, so join the two with a
   slash yourself: there is one upstream and nothing contradicts it. A board
   remote pointing at `origin` is an ordinary setting and says nothing about
   where pull requests go.
3. **A fork whose board remote resolves to neither `origin` nor the parent** —
   two different upstreams. Do not guess and do not open. Name both and ask
   which one this pull request belongs in, quoting `.remote_source` from the
   same `whoami` output: it says which step of jaira's ladder named that remote,
   which is the fact the person needs to answer you.

Everything below takes that repository as `<owner/repo>`. Every rung leaves you
holding one, so every command below names it — there is no branch where the flag
is left off.

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

This is where the two jobs part, and the only place they do. Run the create
command for the forge you settled on, with the repository from above named in
it — **only if a person invoked this role.** On GitHub:

```bash
gh pr create --repo <owner/repo> --title "<title>" --body-file <the description you wrote>
```

On GitLab:

```bash
glab mr create --target-project <owner/repo> --title "<title>" --description "$(cat <the description you wrote>)"
```

Then report the URL.

**An agent invoked you instead: do not run it.** Hand the line back written out
and ready to paste, with the target repository already filled in, and say the
branch is pushed and waiting for a person to open it.

## Answering review comments

Once a person has opened the pull request, the branch is yours to push to and
the thread is yours to answer.

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
  invoked this role. Invoked by an agent, writing that line out is the job and
  running it is the person's.
- Never force-push a branch someone has already reviewed. If history must
  change, say so and ask first.
- Never close a pull request that a person opened.
- Report in three lines: the pull request URL — or the branch you pushed, when
  an agent invoked you and the create line went back instead — what it contains,
  what it is waiting on.
