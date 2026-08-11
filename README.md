# jaira

A kanban board for the work you do with coding agents, stored as markdown files
inside your repository.

Hand a coding agent one task and it reliably becomes five. Within a session that
is fine. Across sessions it is not: what each sub-task was *for* and where it got
to both evaporate, and there is no artifact left to reconstruct them from.

jaira makes that state durable and visible. Tickets are files in the repo, so the
board travels with the code — a teammate clones, runs `jaira`, and sees the same
board. No server, no accounts, no setup.

```
repo/                        localhost is not involved
├── .jaira/
│   ├── tickets/             one markdown file per ticket, committed
│   └── .gitignore           session + lock state stays local
└── src/
```

## Install

```bash
go install github.com/berk/jaira/cmd/jaira@latest
```

Or download a binary from the releases page. It is a single static executable with
no runtime dependency; `git` must be on `PATH`, which it already is if you have a
repository for jaira to work in.

## Start

```bash
cd your-repo
jaira init      # creates .jaira/, registers the merge driver for this clone
jaira           # opens the board
```

Commit `.jaira/` and your team has the board.

## A ticket

```markdown
---
id: 01KZSAMZMANMNFNHS05J0VAG6T
title: Fix session cookie dropped on 302
status: in-progress
ready: true
creator: berk
assignee: berk
executed-by: haiku
goal: Session must survive the OAuth round-trip
context: Reported in chat while debugging Safari logouts
definition-of-done: "session survives OAuth round-trip, covered by a test"
blocked-by: []
commits:
  - 4f2a1c9
model-tier: cheap
outcome-what: Set SameSite=Lax and re-issued the cookie on 302
outcome-why: The cookie was dropped cross-site, silently logging users out
outcome-resolves: >-
  The DoD asked for survival across the OAuth round-trip; re-issuing on
  redirect closes the gap, verified by session_test.go
created-at: 2026-08-11T21:11:27Z
updated-at: 2026-08-11T21:14:03Z
---

Reproduced on Safari only.
```

The file format *is* the API. It is hand-editable, and writing one field rewrites
only that field's bytes — so a lane change shows up in git as a one-line diff, not
a reformatted file.

Three details are load-bearing:

- **`assignee` is always a human**, even when an agent does the work. The model is
  recorded separately as `executed-by`. Ownership of an outcome does not transfer
  to a language model.
- **`outcome-resolves`** is not a restatement of what changed. It is the argument
  that the change satisfies the definition of done — enough to review without
  opening the code.
- **`external:`** is reserved for a future Jira/YouTrack adapter. jaira never
  interprets it and never rewrites it.

## The lanes

```
BACKLOG → TODO → IN PROGRESS → HUMAN → REVIEW → DONE        + BLOCKED
                                 │        ▲
                                 └ needs ─┘
                                   you
```

Anything can be thrown into the backlog. Nothing *leaves* it without a goal, a
definition of done, the context it came from, and an assignee.

That gate is the point of the tool. Capture should be frictionless and execution
should not be: an agent that starts work with no checkable target produces work
you cannot review. Acquiring a definition of done is the price of admission to a
run.

`DONE` additionally requires a signal that did not come from a model:

```bash
jaira move JJN9KH --to done --signal "go test ./... passed; berk reviewed the diff"
```

A review agent cannot certify its own work. This is not politeness about AI — it
is that LLM-as-judge is measurably poor at catching real defects, so a terminal
state gated on a model's own assessment would mean nothing.

## Lanes are pipeline stages

A lane can bind a prompt, a model tier, and an input/output contract. That turns
the board from a progress display into a pipeline you can watch: a cheap model
implements, an expensive one reviews.

Custom lanes are single files in `~/.jaira/lanes/`, so sharing one is sending
someone the file.

```markdown
---
id: critique
name: Critique
after: review
precedence: 55
agentic: true
model-tier: strong
input-requires: [goal, definition-of-done, outcome-what, diff]
output-produces: [outcome-resolves]
---
# Prompt

Look for defects the implementer would not have noticed...
```

Ask the tool to assemble the input rather than letting the agent decide what
context matters:

```bash
jaira show JJN9KH --for-lane critique --json
```

That returns the prompt, only the declared fields, and the diff of the ticket's
own commits. If the agent chose its own context, the contract would be a
suggestion.

A teammate without your `critique` lane still sees those tickets — in a read-only
passthrough column. Hiding them would be the worse failure.

## Concurrency

Two people moving the same ticket both rewrite the same `status:` line. Line-based
merging calls that a conflict, on the single most common operation the whole system
performs.

So `jaira init` registers a **field-aware merge driver** for the clone:

| Field | Resolution |
|---|---|
| `status` | whichever lane is further along — never revert progress |
| `blocked-by`, `commits` | union; neither side's addition is lost |
| other scalars | the more recent `updated-at` wins |
| prose (`goal`, `outcome-*`, body) | a real conflict, scoped to that field |

A conflicted ticket stays valid YAML — the contested value is parked in a
`conflict-theirs-<field>` key and listed under `merge-conflicts`, rather than the
file being filled with markers. Conflict markers would make the frontmatter
unparseable and blank the ticket on everyone's board until someone resolved it.
`jaira resolve` shows both sides and can take either.

Registration writes to `.git/config`, because git deliberately does not let a
clone configure an executable on your behalf. jaira says so when it does it rather
than acting quietly.

Known limit: merge drivers only run for merges git performs locally. A conflict
resolved through a hosting provider's web UI falls back to line-level markers.

Within one machine, concurrent CLI calls take a per-ticket lock and write
atomically, because a single write path prevents schema drift but does nothing
about two processes interleaving.

## With Claude Code

Install the skill from `.claude/skills/jaira/` and Claude will use the CLI
directly. Optionally install `hooks/sync-tasks.sh` as a `PostToolUse` hook on the
task tools to mirror Claude's task list into the backlog.

Mirrored tickets land in the backlog behind the gate. Lane movement is
deliberately *not* mirrored: letting an external status push a ticket into the
pipeline would route around the gate that makes the pipeline worth having.

The sync is idempotent — a task already mapped to a ticket updates it rather than
creating a second one, and an unchanged list writes nothing — so a
board→tasks→board round trip settles instead of oscillating.

## Commands

```
jaira                      open the board
jaira init                 prepare a repository
jaira create <title>       create a ticket
jaira list                 list tickets
jaira show <id>            show one ticket
jaira set <id> k=v…        set fields
jaira move <id> --to …     move lanes, applying the gates
jaira next                 the next actionable ticket
jaira claim <id>           take a 30-minute lease on a ticket
jaira lanes                installed lanes
jaira checkpoint           record what this session is doing
jaira sessions             sessions working this tree
jaira sync-tasks           mirror an agent task list into the backlog
jaira tasks                emit the board as a task list
jaira resolve <id>         settle the fields a merge could not resolve
jaira projects             boards you have opened
```

Every read command takes `--json`. Exit codes are a stable contract:

| Code | Meaning |
|---|---|
| 0 | success |
| 1 | unexpected error |
| 2 | usage error |
| 3 | a gate refused the operation |
| 4 | unresolved dependencies |
| 5 | no such ticket, or an ambiguous id |

Under `--json`, refusals are structured on stderr with a `field` naming what to
supply — so an agent can fix and retry without parsing prose.

## Keys

```
h l ← →   lane            enter   open ticket      n   new ticket
j k ↓ ↑   card            d       diff             m   move lane
g G       first / last    /       filter           r   reload
                          ?       help             q   quit
```

## What this deliberately is not

`paca`, Jira, and Linear already exist. jaira is smaller on purpose, and the
following are excluded with reasons rather than left as future work:

| Not built | Why |
|---|---|
| Sprints, custom fields, roles, saved views, dashboards | Full-PM-tool features. The project fails by growing. |
| Server, accounts, authentication | Git is the sync layer. Repository access already *is* the permission model. |
| Web UI, VS Code extension | Sharing is solved by git. Splitting across UI surfaces is the bloat trap. |
| Branch per ticket | Too heavy for small tasks; breaks down under parallel agents. |
| Board-spawned background agents | Orchestration stays in your session. Process lifecycle is a large surface for little gain. |
| Multi-type dependency graphs | A flat `blocked-by` covers the actual need. |
| A CRDT ticket engine | Solves conflicts completely, but abandons hand-editable plain files — a worse trade than tolerating rare prose conflicts. |
| Jira / YouTrack sync | Deferred. The `external:` block reserves room so no migration is needed. |

Worth knowing honestly: file-based issue trackers have a poor track record.
ticgit, Bugs Everywhere and others are dead or niche, and Fossil explicitly
rejected mutable working-tree ticket files for exactly the merge-conflict reason
above. jaira's bet is that a field-aware merge driver plus a deliberately small
surface makes the plain-file approach work where those attempts did not. That bet
is not yet proven by adoption.

## Development

```bash
go test ./...
go build ./cmd/jaira
```

Layering is enforced by the module graph: `core/` imports nothing from `cmd/` or
`internal/`. The CLI and the TUI are peers over the same core, which is the only
reason "both interfaces enforce the same rules" is true rather than aspirational.
