---
name: jaira-role-brainstorm
description: "Turn a raw idea into tickets somebody else can act on — or say it is not a ticket. Invoked as /jaira-role-brainstorm <idea>. Use before anything is on the board, when what is wanted is still a sentence rather than a task."
---

# Shape the idea into tickets, or refuse it

You work upstream of the board. Nothing exists yet: there is an idea, and your
job is to end with tickets a stranger could pick up — or with a clear statement
that there is no ticket here.

This is not the `brainstorm` lane. That lane works a ticket that already exists,
and its prompt lives on the board (`jaira show <id> --for-lane brainstorm
--json`). If a ticket id was handed to you, use the lane, not this.

## Understand before you shape

1. **What is wrong today?** Not what is missing — what actually hurts, and to
   whom. An idea whose owner cannot name a current pain is a preference, and it
   should be said so.
2. **What triggered it now?** Ideas have dates. The trigger usually names the
   real requirement better than the idea does.
3. **What has already been tried or ruled out?** Look: `CLAUDE.md` / `AGENTS.md`,
   related tickets, `jaira note` entries, git history. Re-proposing something
   this project rejected on purpose wastes everyone's afternoon.
4. **What would have to be true for this to be a bad idea?** Ask it out loud. An
   idea nobody argued against has not been shaped, only accepted.

If the answer would change what gets built and you do not have it, hand out
`/jaira-role-research` before going further. Shaping on top of a guessed fact bakes
the guess into the definition of done.

## Then measure it against the project

Every project has a line it will not cross — read it before proposing anything.
For this repository it is in `CLAUDE.md`: scope discipline, and the constraints
that are there by construction. An idea that violates a constraint is not a
small ticket, it is a different project. Say which constraint, and stop.

## Splitting

Cut along **what a user can observe**, never along the layers of the code. Three
tickets called model, handler and UI are one ticket wearing a costume: none of
them can be finished, reviewed or reverted alone.

Each resulting ticket must be able to answer, on its own:

- **goal** — what it is for, one sentence
- **context** — what is wrong today, what triggered it, what is ruled out.
  Written for someone who was not here and reads it weeks from now: what is
  wrong first, short concrete lines with one point each, names and paths rather
  than adjectives. If acting on it needs a question answered first, it is not
  finished.
- **definition of done** — checkable. "Works well" is not a criterion; "`t`
  draws the list over the board, the board stays visible, `go test ./... -race`
  green" is.

A ticket you cannot give a definition of done is a ticket you have not
understood yet. Go back to question 1.

## Refuse cleanly

Three outcomes are all legitimate, and the last two are not failures:

- **tickets** — file them
- **one ticket, smaller than the idea** — say what you cut and why
- **no ticket** — the pain is not real, it is already solved, or it crosses a
  project constraint. Say which, in one line, and stop

## File them

```bash
jaira tags                       # reuse a name, never invent a synonym
jaira create "<title>" --goal … --context … --dod … --tag <existing>
```

Then commit the ticket files. A ticket nobody but you knows about is a note in a
dead session.

## Report

Per ticket: id, title, one line on why it exists. Then what you deliberately did
not file, and why. That second list is the one that stops the idea coming back
next month unchanged.
