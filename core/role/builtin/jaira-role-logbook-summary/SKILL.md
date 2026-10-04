---
name: jaira-role-logbook-summary
description: "Summarise a stretch of work for a project lead - what is finished, what is under way, what is open, hours estimated from commit times - as one page with every text ready to copy, and no code in it. Plain by default; a summary.md in the repository or in ~/.jaira says what the reader's own system needs. Invoked as /jaira-role-logbook-summary [period]. Use when someone asks how far the work has come, for a weekly report, or for hours to book."
---

# What happened, for someone who does not read code

The reader is a project lead. They want to know how far the work has come,
what is finished, what is still open and how many hours it took — not which
files changed. Everything here serves that reader.

## 1. What is wanted

Read, in this order, and stop at the first that exists:

1. `.jaira/summary.md` in the repository
2. `~/.jaira/summary.md`

It is prose, written by a person, describing what they need: the target system
(a ticket system, an email, a wiki), the markup (Markdown, plain text, a
system's own markup),
the language, how to bundle the work (one entry per theme, two tickets for two
customers), how to split hours, the session rule for hours. Follow it where it
speaks; where it is silent, the default below applies. `summary.example.md`
next to this file shows what such a file can say.

**Without a file, the default:**

- period: the last full calendar week, Monday to Sunday, unless the request
  names another
- language: the language of the conversation
- markup: Markdown
- bundling: by theme — a milestone, a `follows` chain, a shared tag; 2 to 6
  themes, a ticket in exactly one
- hours: from commit times, a new session after 90 minutes without a commit,
  30 minutes counted before a session's first commit

Do not ask what the file would have answered. The default is meant to be
good enough to send.

## 2. Collect

```bash
python3 <this skill's folder>/collect.py --from <YYYY-MM-DD> --to <YYYY-MM-DD> [--gap 90] [--lead-in 30] [--author <email>]...
```

Run it in the repository whose work is being summarised. It reads git history,
`jaira list --json` and `.jaira/logbook/`, and prints facts: commits and
sessions per day with an hour estimate, tickets filed into the logbook in the
period (`done`), and tickets on the board that moved in it (`active`). It
also shares each day's hours out over the tickets its commits name
(`hours_by_ticket`, per day and for the period): the time between two commits
belongs to the ticket of the later one, and a day's tickets add up to that
day's hours exactly. One
person's work across several repositories: run it in each and add up.

For each ticket you will write about, read what it was for: `jaira show <id>`
reads it on the board and in the logbook alike, and `done` already carries it.
The `goal` and the `outcome` (`what`, `why`) are your raw material.

**Days with no commits** inside a working week: say so on the page. If the
person may have worked on them anyway — an acceptance in the browser, a
meeting — ask once, listing those days, and take their hours as they give
them, marked as their own figure. Never invent hours.

## 3. Write

For each theme:

- **title**: what changed for the people who use the product, in their words
  ("Two runs no longer collide"), not the ticket's title if that is technical
- **status**: `done`, `active` or `open`, with a label in the page's language
- **body**: three to six short lines — what it is about, the result, what is
  still open. Plain words. No file names, function names, paths, commit
  hashes, or code formatting: `build.py` refuses a page that has them
- **copy**: the same content in the target markup, if that differs from the
  body; otherwise leave it out
- **comments**: when the target keeps a running record per entry — a ticket
  with comments, a log — one comment per day that entry moved, saying what
  was finished that day, each in the target markup. The facts are in
  `days[].by_ticket`; a day the entry did not move gets no comment

Then:

- **figures**: finished, under way, hours (estimate) — three to five numbers
- **hours**: one row per day — sessions, commits, hours; a total row last,
  and a note that says the hours are estimated from commit times and how
- **a short entry per working day**, if the target wants one: what got done
  that day, one or two sentences — a section of its own
- **appendix**: internal ticket ids with their titles, for whoever needs to
  trace a line back. Ids belong only here
- **hours per ticket**, only when `summary.md` asks for them: take them from
  `hours_by_ticket` as they are — never estimate or round them yourself, so
  that they add up to the hours per day. A theme's hours are the sum of its
  tickets'. `-` is time on commits that name no ticket; show it as such

Write `summary.json` beside nothing in the repository — in
`~/.claude/summaries/<repo>/<from>_<to>/`:

```json
{
  "eyebrow": "Weekly report · CW 40", "title": "...", "subtitle": "...", "lang": "en",
  "figures": [{"label": "finished", "value": "8"}, {"label": "hours", "value": "38.5", "note": "estimate"}],
  "hours": {"title": "Hours per day", "columns": ["Day", "Sessions", "Commits", "Hours"],
            "rows": [{"cells": ["Mon 28.09.", "09:10-12:40", "14", "4.0"]}], "note": "..."},
  "sections": [{"title": "What changed", "intro": "", "items": [
      {"title": "...", "status": "done", "status_label": "finished", "meta": "3 tickets", "body": "...", "copy": "...",
       "comments": [{"date": "30.09.", "text": "..."}]}]}],
  "appendix": {"title": "Internal ticket ids", "lines": ["7P6J0C  A refused start says the run did not start"]}
}
```

## 4. Build and hand over

```bash
python3 <this skill's folder>/build.py page summary.json summary.html summary.md
```

A refusal names the line with code in it: rewrite that line, do not work
around the check.

Publish `summary.html` with the Artifact tool (no capabilities, `icon`
`report`, a one-line `description`). Reply with the link and three lines: the
period, the figures, and anything the reader should know before sending it on
— a day with no commits, hours the person gave you, a theme still open.

## What you never do

- invent hours, or round a day with no commits up to a working day
- put code, paths, hashes or function names in the text for the reader
- leave a ticket out because it is unfinished: open work is part of how far
  the work has come
- write into the repository: the summary is about the work, not part of it
