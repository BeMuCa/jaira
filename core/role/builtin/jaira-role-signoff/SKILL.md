---
name: jaira-role-signoff
description: "Read tickets waiting in signoff and put each to the person as one status line, one table that sorts every open note into blocker, small edge, known limit or not from this ticket, and one recommendation with both move lines ready to paste. Changes nothing itself. Invoked as /jaira-role-signoff [ticket-id...] by a person or a teamlead. Use when work sits in signoff and someone has to decide whether to accept it or give it back."
---

# Signoff: one decision per ticket, made easy

A ticket in signoff is finished work a model has already reviewed. What is left
is one decision a person takes: accept it or give it back. Your job is to make
that decision take a minute, not to take it.

Measured on 2026-11-20: a person sent back almost every ticket from signoff,
because the reviewer's notes read one after another like a list of faults —
known limits, a ten-minute edge and a real blocker all in the same voice. The
same notes, sorted into a table with a mark per line and one recommendation on
top, were understood at once and accepted. This role is that shape.

**You change nothing.** No move, no note, no `set`, no commit, no ticket. You
read and you report; the lines you hand back are for the person to run. A
person accepts work in signoff, never an agent — and the person deciding is
also the person who writes the reason down.

## Which tickets

The ids you were given. None given: every ticket in signoff, on every board:

```bash
jaira boards
jaira list --lane signoff --json                    # today's board
jaira list --lane signoff --board <name> --json     # every other board
```

## Read the state now, never from memory

Before you put a ticket in front of the person, read it fresh — the person may
have accepted or given it back a minute ago in the board or the browser, and
nothing tells you:

```bash
jaira show <id> --json      # status, handle, dod_items, review, body
gh pr list --state open      # or glab mr list: what is waiting as a pull request
```

The keys of that JSON, as this role reads them: `.status` the lane,
`.handle` the six characters, `.dod_items[]` with `.done` per item,
`.review.verdict`, `.review.gaps` and `.review.summary` the review lane's
outputs. There is no `notes` key: the notes are the `- **<date> · <who>** —
<text>` lines under `## Progress` in `.body`.

A ticket no longer in signoff drops out of the report without a word. Nothing
left: say so in one line and stop.

## Per ticket

### 1. The status line

One line, each check with ✅, ❌ or `?`. Every check is about the commit that would
ship — the head of the ticket's branch — and nothing older:

The ticket's branch is the one whose name carries its six-character handle (`handle` in `jaira show <id> --json`) as
a segment of its own — `feat/twrap8`, `fix/TWRAP8-signoff`, `alsa/twrap8` —
the rule `jaira create` uses to read a ticket off a branch. Never assume a
`feat/` prefix:

```bash
seg='(^|[^[:alnum:]])<handle>([^[:alnum:]]|$)'   # the six characters, never the ULID
refs=$(git for-each-ref --format='%(refname:short)' refs/heads refs/remotes | grep -iE "$seg")
head=
for c in $refs; do                 # the copy every other copy is an ancestor of
  for r in $refs; do git merge-base --is-ancestor "$r" "$c" || continue 2; done
  head=$(git rev-parse "$c"); break
done
```

Local and remote copies are compared by ancestry, not by where they live: the
head is the commit every other copy leads up to — a local branch ahead of its
push, or a remote one a teammate pushed after your last fetch. No copy at all,
or copies that diverged with none of them leading on from all the others:
`head` stays empty, the two checks bound to it — tests and bug hunt — are
`?` (DoD, blocker, review and returns are read off the ticket as usual), and the
status line names `$refs` — which one ships is the person's to say, not yours
to guess.

A CI result that belongs to another commit, or no result at all, is `?`,
never ✅ — nobody knows yet. The bug hunt is different: it is the board's own
work, so a hunt that does not name the head was not done there, and that is ❌.

- **DoD** `x/y` ticked, from `.dod_items[].done`
- **tests** — the CI run for exactly that commit, of the workflow that runs the
  tests (named in `.github/workflows/`): `gh run list --commit "$head"
  --workflow <ci workflow> --json conclusion,headSha`; on GitLab `glab ci list
  --sha "$head"`. One commit can have several runs (a push and a pull request
  event, a re-run): the newest decides — `gh` lists newest first. Its
  conclusion `success` is ✅, a failure ❌, and a run still queued or in
  progress (empty conclusion) is `?`. A run for an older commit says nothing
  about this one: `?` — and if the branch is ahead of what was pushed, say
  that is why
- **bug hunt on HEAD** — the last `— bughunt: done at <sha> ` note names
  `$head` — the note's sha is a prefix of the full `$head`, whatever length
either clone abbreviates to. An older sha does not cover what a later
  round changed: ❌
- **blocker** — none, or how many 🔴 the table below holds
- **review** — `.review.verdict`
- **returns from signoff** — how many are noted (see return rule 2)

### 2. The table

Every open point the ticket carries — each line of `.review.gaps`, each
`note:` line in it, each finding of the hunt that is still open, each `decision:`
note the person has not seen — one row:

| # | the point in one plain sentence | mark |
|---|---|---|

The sentence says what a user of the software would notice, not what the
reviewer wrote. The mark is exactly one of:

- 🔴 **blocker** — a wrong result, lost data, a crash, or a security or
  permission boundary that does not hold. Nothing else is a blocker
- 🟡 **small edge** — a real gap in this ticket's own change that one round of
  about ten minutes closes
- ⚪ **known limit** — a boundary the ticket draws on purpose, written down; no
  work follows from it
- 🔵 **not from this ticket** — pre-existing, or another ticket's subject; say
  where it gets fixed (the ticket id, or "new ticket")

Wording alone — a comment, a name, a sentence that could read better — is never
🔴 and never 🟡. Leave it out of the table, or ⚪ if the reviewer insisted.

### 3. One recommendation

Bold, one line, with its reason: **accept** or **give back**, never both, never
"up to you". Then both lines, ready to paste:

```bash
jaira move <id> --to done --force      # accept: signoff only lets a person move on
jaira note <id> "Given back from signoff: <the rows to fix, by number>" && jaira move <id> --to in-progress
```

`--force` on the accept line is the person's to pass: the gate refuses an
agent's move out of signoff, and the line is handed to the person precisely so
that they, not you, run it.

## The return rules

These decide the recommendation. They exist because without them a ticket goes
round signoff forever, each round finding the next smaller edge. The status
line comes first: the table only sorts notes, and a red check is no note.

0. **A ❌ on the status line → give back, whatever the table says.** Tests
   failed at the head: give back with the failing run's link. No hunt names
   the head: give back so the hunt runs on it — the change is not done
   without one. A `?` is not a ❌, but it is not a ✅ either: never
   recommend accept over it without saying which check is unknown and what
   settles it (push and wait for CI, run the hunt), and put that in the
   recommendation's reason.
1. **Any 🔴 → give back.** Only the 🔴 rows, and the 🟡 rows with them.
2. **Only 🟡, ⚪ and 🔵:** count the ticket's earlier returns from signoff —
   the notes that start `Given back from signoff:`, which the cockpit and the
   give-back line above write. A ticket moved back by hand in the board
   leaves no such note, so the count is a floor: put it on the status line
   (`returns from signoff: 1 noted`), where the person sees it and can say it
   was more. `round` is no substitute — it rises with every critique, optimize
   and testing loop too.
   - fewer than two: **give back**, the 🟡 rows by number — one small round
     costs less than a follow-up ticket nobody picks up
   - two or more: **accept**, and hand back a ready follow-up ticket for the
     🟡 rows instead (below). A third round for edges is what this rule ends
3. **Only ⚪ and 🔵 → accept.** A known limit and someone else's subject are
   no reason to hold this ticket.
4. **The same class keeps coming back** — the third finding about, say, CRLF
   input, or error paths that lose the message: name the cause in one line under
   the table, not a fourth example. Giving back the examples one at a time is
   what made the rounds endless.
5. **Exception: a permission or security boundary stays 🔴**, whatever the
   round. Rule 2's limit never accepts one.

The follow-up for rule 2, ready to paste and naming the ticket it follows:

```bash
jaira create "<what is still missing, as a fact>" --follows <id> --goal "<...>" --context "<the 🟡 rows, each in one line, and that <id> was accepted with them after two returns>" --dod "<...>" --assignee "<who owns <id>>"
```

## The report

Per ticket: the status line, the table, the recommendation in bold, the two (or
three) lines. In the language the person speaks to you — the marks and the
commands stay as they are. Then one line on what comes after the decision: a
pull request to open, a worker to give back to, nothing.

Several tickets: one block each, the ticket with a 🔴 first, then the ones you
recommend giving back, then the ones you recommend accepting. A list of open
decisions is shown again only when it changed since you last showed it.
