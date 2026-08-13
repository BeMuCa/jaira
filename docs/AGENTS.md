# Using jaira with a coding agent

jaira is not tied to Claude Code. Anything that can run a shell command and read
JSON can work the board — Claude Code, Codex, Aider, Cursor, a local model behind
Ollama, or a shell script.

The whole integration surface is: **run a command, read the JSON, branch on the
exit code.** There is no SDK, no plugin API, and nothing to keep in sync.

## The loop

```bash
# 1. What should I work on?
jaira next --json

# 2. What exactly does this step want from me?
jaira show <id> --for-lane in-progress --json
#    → { "prompt": "...", "input": { "goal": "...", "plan": "..." },
#        "produces": ["outcome-what", ...], "missing": [], "model_tier": "cheap" }

# 3. Say where you are, as you go.
jaira dod <id> 2 --doing --plan
jaira dod <id> 2 --done --plan
jaira dod <id> 3 --done --proof "internal/x.go:12; TestX"   # tick it and record why in one call

# 4. Record what a later session would have to rediscover.
jaira note <id> "writeAll buffers the whole file; flushing per 5k rows works"

# 5. Hand the work on.
jaira move <id> --to review \
  --what "streamed the writer" \
  --why "it buffered the whole export" \
  --resolves "a 40MB export now completes on a throttled link" \
  --commits "$(git rev-parse HEAD)" \
  --executed-by <your-model-name>
```

If a gate refuses, the message says what to do and exit code 3 says it was a
refusal rather than a crash. An agent can act on it without a human translating.

Lanes before the first lane that declares itself the start of the specified
zone (`requires-specified: true` — the built-in `todo` lane, by default) are for
tickets that are not specified yet: a scrap, an idea, something being thought
through. The promotion gate fires when a ticket crosses into that zone, not
when it leaves a lane with a particular name.

## Starting a session

```bash
jaira resume --json
```

Returns every ticket left mid-flight — a step still marked in progress, an
expired claim, a ticket parked in a working lane — with the notes written against
it and the step it was on. A session that died to a usage limit leaves nothing
behind except what was written down; this is how the next one picks it up.

## Working in parallel

```bash
jaira claim <id>        # a 30-minute lease, so two sessions do not collide
jaira checkpoint --focus "chunking the CSV writer" --ticket <id>
```

Claims are advisory leases, not locks: an abandoned one expires rather than
needing an unlock. `jaira sessions` shows what every session in the tree is on,
and the board's compact view (`v`) counts them per step.

## What an agent cannot do

Three things, deliberately:

- **Leave the sign-off lane.** It is a human checkpoint. A review agent writes
  `review-summary`, `review-gaps` and `review-verdict`, then moves the ticket to
  sign-off and stops; a person accepts the work in the board or raises a
  follow-up. Attempting the move exits 3.
- **Close a ticket without meeting its definition of done.** Every criterion must
  be marked done, and the plan finished if the ticket has one. There is no
  evidence flag to pass. `--force` exists, is recorded on the ticket, and is the
  user's call rather than the agent's.
- **Write to someone else's ticket.** A ticket belongs to its `assignee`; a move
  by anyone else exits 3, naming the owner. The human checkpoint lanes are
  exempt — reviewing and signing off someone else's work is the whole point of
  them — and taking a ticket over with `jaira set <id> assignee=<you>` is always
  allowed, so a ticket is never frozen by an owner who has gone quiet. `--force`
  overrides and is recorded, same as any other gate refusal.

Ownership here is a guard rail, not a lock: tickets are still plain markdown
files in git, so a hand edit or an offline merge can produce two versions of a
ticket regardless of what this binary would have refused. The merge rules
handle that case and are unchanged by any of this.

## Telling an agent the board exists

There is no single convention for this. `jaira init` writes the same section
into both **`AGENTS.md`** — the closest thing to a cross-tool standard, and what
Codex reads — and **`CLAUDE.md`**, between markers:

```
<!-- jaira:start -->
## Task tracking: jaira
This repository has a jaira board (`.jaira/`) …
<!-- jaira:end -->
```

For anything else, copy that block into whatever your tool reads:

| Tool | File |
|---|---|
| Codex, and increasingly others | `AGENTS.md` — written for you |
| Claude Code | `CLAUDE.md` — written for you, plus a skill in `.claude/skills/jaira/` |
| Gemini CLI | `GEMINI.md` |
| Cursor | `.cursorrules` or `.cursor/rules/` |
| Aider | the conventions file named in its config |

Re-running `jaira init` updates the block in place and leaves the rest of each
file alone.

For Claude Code there is also a skill in `.claude/skills/jaira/`; copy it to
`~/.claude/skills/jaira/` to have it available in every repository.

## A note on model tiers

A lane can declare `model-tier: cheap` or `strong`. jaira does not run models and
does not know what those names mean on your setup — it passes the tier through in
`--for-lane` output, and the thing driving the agent decides what to launch. That
is the only place model choice appears.
