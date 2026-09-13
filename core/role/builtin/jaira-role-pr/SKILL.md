---
name: jaira-role-pr
description: "Open a pull request for finished ticket work, push to it, and answer review comments — never merge and never approve. Invoked as /jaira-role-pr <ticket-id>, normally by a dispatcher. Use when a branch is ready to leave for review."
---

# Open it, answer it, never accept it

Arguments: `<ticket-id>`. Without it, say what is missing and stop.

Opening a pull request is a contributor's job. Accepting one is the
maintainer's. You are the contributor: you may open, push, and answer. You may
not merge, and you may not approve — not your own work, and not anyone's.

## Before you open anything

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

## The description is the ticket, not a summary of the diff

Read it off the board, do not invent it:

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

## Answering review comments

- One comment, one reply, one commit. Do not batch a reviewer's five points into
  one commit that none of them can be reverted out of.
- A comment you disagree with gets an answer, not a silent change and not a
  change you make anyway. Say why, and leave it to them.
- A comment that is out of scope becomes a `jaira note` or a new ticket, and the
  reply says which. Never widen the pull request to close a comment.
- After pushing, say what you pushed in one line. The reviewer should not have
  to diff to find out what you did with their point.

## Boundaries

- **Never `gh pr merge`.** Never `gh pr review --approve`.
- Never force-push a branch someone has already reviewed. If history must
  change, say so and ask first.
- Never close a pull request that a person opened.
- Report in three lines: the pull request URL, what it contains, what it is
  waiting on.
