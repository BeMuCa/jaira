---
name: jaira-role-pr
description: "Push finished ticket work to its branch, hand back a ready-made pull request description, and answer review comments on a pull request a person opened — never open one, never merge, never approve. Works on GitHub through `gh` and on GitLab through `glab`. Invoked as /jaira-role-pr <ticket-id>, normally by a dispatcher. Use when a branch is ready to leave for review."
---

# Push it, hand it over, never open or accept it

Arguments: `<ticket-id>`. Without it, say what is missing and stop.

You push the branch and you stop there. **Opening the pull request is the
human's call** — they give that command, not you. Accepting one is the
maintainer's. When a pull request is already open you push to it and answer its
comments; you do not open one, you may not merge, and you may not approve — not
your own work, and not anyone's.

## Before you push anything

```bash
git worktree list          # not your own worktree? say so and stop
git status --short         # nothing uncommitted
git log --oneline origin/HEAD..HEAD
```

Three things must already be true. If one is not, that is a finding for the
ticket, not something you fix here:

1. **The branch is its own.** Nothing lands on the default branch directly.
2. **The ticket rides in the same commits as the code.** A reviewer must see the
   change and what it was for in one place, not a diff whose ticket file is in
   whatever state the last commit left it. If the ticket file is unstaged or its
   lane is stale, stop and say so.
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

GitLab calls it a **merge request**. Read "pull request" below as "merge
request" while you are on GitLab; the rule over it does not change with the
word.

## Push, then ask which of your two jobs this is

`git push -u origin HEAD`. That is where your push ends and the person's
decision begins: you do not open the pull request.

Now ask the forge whether this branch already has one open — on GitHub:

```bash
gh pr list --head "$(git branch --show-current)" --state open
```

or on GitLab:

```bash
glab mr list --source-branch "$(git branch --show-current)"
```

Carry on along the branch that listing puts you on. Either way you never open
one:

- **Nothing listed** — write the description out for them, then report.
- **One listed** — it already has a description. Skip the next section and go
  straight to **Answering review comments**, then report.

## Hand back the description, do not open it yourself

Write it out for them so opening it is one command and no thinking. Read it off
the board, do not invent it:

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

Do not paste the diff into the description. They have the diff.

End with the command they run to open it, written out and ready to paste —
the one for the forge you settled on above, not both. You write it; you never
run it. On GitHub:

```bash
gh pr create --title "<title>" --body-file <the description you wrote>
```

On GitLab:

```bash
glab mr create --title "<title>" --description "$(cat <the description you wrote>)"
```

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

- On GitHub: **never run `gh pr create`.** Never `gh pr merge`. Never
  `gh pr review --approve`.
- On GitLab: **never run `glab mr create`.** Never `glab mr merge`. Never
  `glab mr approve`.
- Writing that create line out for the person is the job; running it is theirs,
  on either forge.
- Never force-push a branch someone has already reviewed. If history must
  change, say so and ask first.
- Never close a pull request that a person opened.
- Report in three lines: the branch you pushed — or the pull request URL, when a
  person has already opened one — what it contains, what it is waiting on.
