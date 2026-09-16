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
`jaira.remote`: that is the remote the board's ticket refs travel on, and in a
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

## Push, then ask which of your two jobs this is

`git push -u origin HEAD`. Then open the pull request, if a person invoked
you; if an agent did, the push is where you stop.

Now ask the forge whether this branch already has one open — on GitHub:

```bash
gh pr list --head "$(git branch --show-current)" --state open
```

or on GitLab:

```bash
glab mr list --source-branch "$(git branch --show-current)"
```

Carry on along the branch that listing puts you on. Either way you never open a
second one:

- **Nothing listed** — write the description, then open it, then report.
- **One listed** — it already has a description. Skip the next section and go
  straight to **Answering review comments**, then report.

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
file — you do not open anything until the next section has settled where it goes.

## Which repository it goes to

`origin` is not the answer. In a fork `origin` is your fork, and a pull request
opened against the fork instead of the upstream sits where nobody is looking and
has to be closed by hand. Settle the target **before** you type a create command:

```bash
gh repo view --json isFork,parent,nameWithOwner   # GitHub
glab repo view                                    # GitLab
git config jaira.remote                           # the remote the board travels on
```

1. **Not a fork** — the target is `origin`'s own repository, and the create
   command below needs no repository flag.
2. **A fork** — the target is the parent, not your fork. `jaira.remote` is the
   remote the board's ticket refs travel on and it names that same upstream;
   read the owner/repo off it and pass it to the create command explicitly.
3. **A fork whose parent and `jaira.remote` are not the same repository** — do
   not guess and do not open. Name both and ask which one this pull request
   belongs in.

## Open it

Run the create command for the forge you settled on above — not both, and with
the repository from the step above named in it. On GitHub:

```bash
gh pr create --repo <owner/repo> --title "<title>" --body-file <the description you wrote>
```

On GitLab:

```bash
glab mr create --target-project <path> --title "<title>" --description "$(cat <the description you wrote>)"
```

Leave `--repo` / `--target-project` off only in case 1, where there is a single
repository the command can mean.

Report the URL. An agent invoked you instead: do not run the command — hand it
back, written out and ready to paste, with the target repository already filled
in, and say the branch is pushed and waiting for a person to open it.

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
