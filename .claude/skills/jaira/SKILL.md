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

That lands in the backlog. It cannot *start* until it has a goal, a definition of
done, the context it came from, and an assignee. Supply them at creation to save a
round trip:

```bash
jaira create "Fix session cookie dropped on 302" \
  --goal "Session must survive the OAuth round-trip" \
  --dod "session survives OAuth round-trip, covered by a test" \
  --context "Reported while debugging Safari logouts"
```

Write the `--context` yourself from the conversation that produced the ticket.
That field is the whole point: it is what makes the ticket comprehensible later.
A good context says what problem prompted this and what had already been ruled
out. "User asked for it" is useless.

The definition of done lives in the body as a **checklist**, and `create` writes
the heading for you. Each item must be checkable by someone who was not here.
"Works properly" is not a definition of done; "`go test ./auth` passes and the
cookie survives a 302" is.

Mark items as you go, with the command rather than by editing the file:

```bash
jaira dod <handle> 2 --doing    # you are working on item 2 now
jaira dod <handle> 2 --done     # you have satisfied it
jaira dod <handle> 1 --todo     # you were wrong, put it back
```

Only one item can be `--doing` at a time; marking a second moves the marker. That
marker is how a person watching the board knows which criterion you are on.

Do not mark an item done you have not verified.

There is a second, optional checklist under a `## Plan` heading: the method you
are following — write the spec, design it, implement it — as opposed to the
criteria for acceptance. Address it with `--plan`:

```bash
jaira dod <handle> 3 --doing --plan
```

The Plan does not gate anything on its own, but a terminal lane refuses a ticket
whose plan is unfinished: the criteria cannot have been met while the work that
meets them is still in progress.

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

`backlog → todo → in-progress → human → review → done`, plus `blocked`.

- **human** — you need a decision only the user can make. Move the ticket here
  with the question attached rather than guessing:
  `jaira move <handle> --to human --question "Should expired sessions redirect or 401?"`
- **review** — implemented, awaiting sign-off. **You cannot move a ticket out of
  this lane.** It is a human checkpoint: write your verdict with
  `jaira set <handle> review-verdict="..."` and stop. The user accepts it in the
  board, or raises a follow-up ticket. Attempting the move exits 3.
- **done** — requires the definition-of-done checklist to be fully marked, and
  the plan too if the ticket has one. There is no evidence flag to pass; you
  cannot certify your own work complete, and you should not try to route around
  it. `--force` exists, is recorded on the ticket, and is the user's call.

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
