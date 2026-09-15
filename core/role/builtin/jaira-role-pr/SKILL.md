---
name: jaira-role-pr
description: "Push finished ticket work to its branch, hand back a ready-made pull request description, and answer review comments on a pull request a person opened — never open one, never merge, never approve. Invoked as /jaira-role-pr <ticket-id>, normally by a dispatcher. Use when a branch is ready to leave for review."
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
gh pr list --head "$(git branch --show-current)" --state open
```

That last line decides which of your two jobs this is; you branch on it after
the push. Either way you never open one.

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

Then `git push -u origin HEAD`. That is where your push ends and the person's
decision begins: you do not open the pull request. Carry on below, along the
branch the `gh pr list` above put you on:

- **Nothing listed** — write the description out for them, then report.
- **A pull request listed** — it already has a description. Skip the next
  section and go straight to **Answering review comments**, then report.

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

End with the command they run to open it, written out and ready to paste. You
write it; you never run it:

```bash
gh pr create --title "<title>" --body-file <the description you wrote>
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

- **Never run `gh pr create`.** Never `gh pr merge`. Never
  `gh pr review --approve`. Writing the `gh pr create` line out for the person
  is the job; running it is theirs.
- Never force-push a branch someone has already reviewed. If history must
  change, say so and ask first.
- Never close a pull request that a person opened.
- Report in three lines: the branch you pushed — or the pull request URL, when a
  person has already opened one — what it contains, what it is waiting on.
