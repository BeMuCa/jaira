---
name: jaira
description: Use when working in a repository that has a .jaira/ directory — tracking multi-step work as tickets, picking up the next task, recording outcomes, or working the board through its lanes. Also use when the user asks to "work the board", mentions tickets/lanes/backlog, or when a task decomposes into several steps that should survive this session.
---

# jaira

A kanban board stored as markdown files in the repository. Tickets outlive this
session, so work you record here is still legible weeks later — to you, to the
user, and to their teammates who have the same files.

`jaira` is a CLI. Drive it with Bash. Every read command takes `--json`.

## The one rule that matters

**Never edit files under `.jaira/tickets/` directly.** Use the CLI. It enforces
the schema, writes atomically, takes a lock so parallel sessions cannot clobber
each other, and rewrites exactly one field's bytes so the git diff stays one line.
Hand-editing loses all four properties.

## When to create tickets

When a request decomposes into more than about three steps, or when work will
outlive this session. Capture is cheap:

```bash
jaira create "Fix session cookie dropped on 302"
```

That lands in the backlog, assigned to nobody: whoever pulls the ticket out of
the backlog claims it in the same move. It cannot *start* until it has a goal, a
definition of done, and the context it came from. Supply them at creation to
save a round trip:

```bash
jaira create "Fix session cookie dropped on 302" \
  --goal "Session must survive the OAuth round-trip" \
  --dod "session survives OAuth round-trip, covered by a test" \
  --context "Reported while debugging Safari logouts"
```

When the ticket exists *because of* another one — a review found a gap, the work
turned out to need a second pass — say so with `--follows`, which takes a handle:

```bash
jaira create "Retry the export on a dropped connection" \
  --follows DAHC06 \
  --goal "..." --dod "..." --context "..."
```

Without it a follow-up looks like an unrelated ticket, and the chain of why is
gone. It must resolve to a real ticket, so a typo exits 5 rather than writing a
dead link. `jaira show` and the board's detail pane display it, and it rides in
`--json` as `follows`.

Write the `--context` yourself from the conversation that produced the ticket.
That field is the whole point: it is what makes the ticket comprehensible later.
A good context says what problem prompted this and what had already been ruled
out. "User asked for it" is useless. Write it so someone who was not in the
conversation can act on it weeks later, without asking anyone.

The definition of done lives in the body as a **checklist**, and `create` writes
the heading for you. Each item must be checkable by someone who was not here.
"Works properly" is not a definition of done; "`go test ./auth` passes and the
cookie survives a 302" is.

Mark items as you go, with the command rather than by editing the file:

```bash
jaira dod <handle> 2 --doing    # you are working on item 2 now
jaira dod <handle> 2 --done --proof "internal/x.go:12; TestX"   # satisfied, and here is why
jaira dod <handle> 1 --todo     # you were wrong, put it back
```

Only one item can be `--doing` at a time; marking a second moves the marker. That
marker is how a person watching the board knows which criterion you are on.

Do not mark an item done you have not verified. `--proof` records the evidence
in the same call that ticks the item — a file:line or a test name, so the claim
is checkable rather than taken on trust. Setting it again on the same item
replaces the line rather than stacking a second one, and it can be passed on its
own to record evidence without changing the marker.

There is a second, optional checklist under a `## Plan` heading: the method you
are following — write the spec, design it, implement it — as opposed to the
criteria for acceptance. Address it with `--plan`:

```bash
jaira dod <handle> 3 --doing --plan
```

The Plan does not gate anything on its own, but a terminal lane refuses a ticket
whose plan is unfinished: the criteria cannot have been met while the work that
meets them is still in progress.

## Writing a good ticket from a request

When the user describes work, do not invent the fields. Extract what they gave
you and ask for what is missing — the gate will refuse the ticket later anyway,
and asking now costs one message instead of a stalled run.

From a request, take:

- **title** — what is different afterwards, not what you will do
- **goal** — the outcome, one sentence
- **context** — where this came from, so it still makes sense in a month
- **definition of done** — checkable statements a person who was not here could
  verify. This is the one people skip and the one the gate blocks on.

Ask when the definition of done is missing or unfalsifiable. "Works properly"
and "is improved" are not criteria. A single question is enough:

> Before I create this: how will we know it is done? I have "the export finishes
> without the timeout" — is there anything else that has to be true?

If the user is mid-flow and does not want to stop, create the ticket anyway with
what you have. It sits in the backlog, which is exactly what the backlog is for,
and `jaira validate` will list it as still needing a definition of done.

```bash
jaira create "Export survives a slow network" \
  --goal "The CSV export completes on a 3G connection instead of timing out" \
  --context "Reported while debugging the customer's failed month-end export" \
  --dod "a 40MB export completes on a throttled connection"
```

A board may define its own ticket shape in `.jaira/TEMPLATE.md`; `create` uses it
when it exists, so follow the headings you find rather than imposing these.

## Before you stop, and when you start

A session that ends abruptly — a limit, a crash — leaves nothing behind except
what was written down. Record what a later session would otherwise rediscover:

```bash
jaira note <handle> "writeAll buffers the whole file; flushing per 5k rows works
  but the header assumes a single pass"
```

Starting a session, ask what was left mid-flight:

```bash
jaira resume --json
```

That returns every ticket with a step still marked in progress, an expired claim,
or parked in a working lane — with its notes and the step it was on.

## Optional steps

A ticket's `## Options` checklist decides which steps it uses. Planning is
opt-in, so most tickets go straight from todo to implementing:

```bash
jaira dod <handle> --option planning          # this ticket gets a planning pass
jaira dod <handle> --option planning --todo   # it does not
```

Moving into a step the ticket did not opt into is refused; skipping past it is
free.

## Working the board

```bash
jaira next --json            # the next actionable ticket
jaira list --json            # everything
jaira show <handle> --json   # one ticket in full
```

A handle is the short id printed on the board, e.g. `JJN9KH`. Full ids work too.

Advance a ticket by recording what you did:

```bash
jaira move <handle> --to review \
  --what "Set SameSite=Lax and re-issued the cookie on 302" \
  --why "The cookie was dropped cross-site, silently logging users out" \
  --resolves "The DoD asked for survival across the OAuth round-trip; re-issuing on redirect closes the gap, covered by session_test.go" \
  --commits "$(git rev-parse HEAD)" \
  --executed-by haiku
```

`--resolves` is not a restatement of `--what`. It is the causal argument that the
change satisfies the definition of done. A reviewer reads it instead of the code.

`--executed-by` records which model did the work. The `assignee` stays the human
who owns the outcome — never reassign a ticket to a model.

## The lanes

`backlog → todo → pre-process → in-progress → human → review → signoff → done`,
plus `blocked`.

- **pre-process** — work out *how*, before writing anything. Its output is a
  `## Plan` checklist, and it cannot be left without one:
  `jaira dod <handle> --plan --add "read the exporter" --add "implement"`
- **in-progress** (shown as Implementing) — carry out that plan. Mark the step
  you are on with `--doing --plan` so a watcher can see where you are.

- **human** — you need a decision only the user can make. Move the ticket here
  with the question attached rather than guessing:
  `jaira move <handle> --to human --question "Should expired sessions redirect or 401?"`
- **review** — a second model judges the diff. Write, in this order:
  `jaira set <handle> review-summary="..."` (what the change does, from the
  diff), `jaira set <handle> review-gaps="..."` (what is missing — write
  "none" if nothing is), then `jaira set <handle> review-verdict="..."` (your
  conclusion). All three are required before the ticket can leave this lane.
- **done** — requires the definition-of-done checklist to be fully marked, the
  plan too if the ticket has one, and the commits that carry the change
  recorded on the ticket. There is no evidence flag to pass; you cannot certify
  your own work complete, and you should not try to route around it.

  The human review lane cannot be left by an agent on its own judgement — but
  when the user has reviewed the work in this conversation and tells you it
  passed, finishing the acceptance for them is legitimate:
  `jaira move <handle> --to done --force`. The `--force` is recorded on the
  ticket, so the trail shows a person decided and you executed. Never use it
  because *you* judge the work done; only because they said so.

Some lanes are agentic: they carry a prompt and a model tier. To work one, ask the
tool to assemble the bounded input rather than deciding for yourself what context
is relevant:

```bash
jaira show <handle> --for-lane review --json
```

That returns the lane's prompt, only the fields the lane declared it needs, and
the diff of the ticket's own commits. Spawn a subagent at the lane's `model_tier`
with exactly that.

## Working an agentic lane

When several sessions may be running, take a claim first so two agents do not pick
up the same ticket. Claims are 30-minute leases and expire on their own, so a
crashed session never wedges the board:

```bash
jaira claim <handle>
```

Then let the tool assemble the input and hand back structured output:

```bash
jaira show <handle> --for-lane review --json     # prompt + declared fields + diff
# … spawn a subagent at the lane's model_tier with exactly that …
echo '{"outcome":{"what":"…","why":"…","resolves":"…"},"executed_by":"opus"}' \
  | jaira move <handle> --to review --from-lane in-progress
```

`--from-lane` validates your output against what the lane declared it produces and
refuses it with the missing field names if incomplete. Fix and retry rather than
working around it.

## Taking a ticket off the board

```bash
jaira archive <handle>      # moves it to .jaira/archive/, nothing is deleted
jaira archive               # list what has been archived
jaira restore <file>        # put it back
```

## Checking the board is intact

```bash
jaira validate              # exit 3 if any ticket is damaged
jaira validate --json
```

Reports unparseable ids, lanes that are not installed, missing timestamps, and
dependencies on tickets that do not exist. An unspecified backlog ticket is a
warning, not an error.

## Dependencies

```bash
jaira set <handle> blocked-by=<other-handle>
```

A ticket with unresolved blockers cannot start. `jaira next` skips them.

Parking a ticket in the blocked lane requires saying what it is waiting on —
a blocked ticket with no recorded reason looks the same on day one and day
ninety. Either the blocker is a ticket (`blocked-by`, accepted as the reason),
or you say it on the move:

```bash
jaira move <handle> --to blocked --reason "vendor API returns 500s, ticket open with them"
```

Without either, the move exits 3 (`needs_blocked_reason`). Parking is exempt
from the leaving lane's output contract and from the dependency check — a
ticket stopped mid-work has not produced its output yet, and its open blockers
are the reason it is here, not grounds to refuse entry.

## Keeping the board honest

Record what you are doing when the topic changes, so the board reflects reality
even after this session ends:

```bash
jaira checkpoint --focus "auth refactor" --why "session cookies leak on 302" --ticket <handle>
```

## If a merge conflicted

Concurrent lane moves, blockers and commit lists merge on their own. Only
competing prose rewrites conflict:

```bash
jaira resolve <handle>                  # show both sides, per field
jaira resolve <handle> --take-theirs    # or --take-ours
```

The ticket stays readable while a conflict is outstanding, so do not panic-edit
the file.

## Exit codes

Branch on these rather than parsing messages:

| Code | Meaning |
|---|---|
| 0 | success |
| 1 | unexpected error |
| 2 | usage error |
| 3 | a gate refused the operation |
| 4 | unresolved dependencies |
| 5 | no such ticket, or an ambiguous id |

With `--json`, a refusal prints `{"error":{"reason":…,"violations":[{"code","field","message"}]}}`
on stderr. Read `field` to learn what to supply, fix it, and retry.

## When a gate refuses

It is telling you something is genuinely missing. Supply it. Do not reach for
`--force` — it exists for the user, not for you, and every use is recorded.
