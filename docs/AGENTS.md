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
jaira dod <id> 3 --text "the wording it should have had"    # reword it; state and proof stay
jaira dod <id> 4 --superseded                               # [-] it will not happen; never tick a stale item instead

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

`jaira next` hands back the furthest-along ticket, which is the right answer when
you are finishing work and the wrong one when you are looking for it: a deep
queue in a late lane hides every earlier lane. `jaira next --per-lane --json`
reports every lane that has work, in pipeline order, with its `agentic` flag —
use it to find the front line, then `--lane <id>` to work one.

Carry a started ticket forward rather than stopping at each lane: after finishing
a step, move it to its `next_lane` and keep going until the next lane is a human
step (`requires-human-exit` or `requires-question` in `jaira lanes show <id>`),
a gate refuses, or the work is blocked. Those are the only three places an agent
should hand a ticket back.

Do not work out the route from the column order. Every `--json` ticket names
`next_lane` — the lane this ticket goes to next, with its own Options applied and
the parking and question lanes left out. It is empty for a ticket that is itself
parked or waiting on an answer, because such a ticket goes back to whatever
raised the question rather than onward. A lane that sends work back also names
where: `rejects-to` in `jaira lanes show <id>`.

A review hands the reader something they can act on. `review-check` is the steps
a person follows to check the change themselves, as a flow, and the review lane
declares it as an output — so a review that skips it cannot leave the lane.

**The ticket rides in the same commit as the code.** Move it first, then stage the
changed file under `.jaira/tickets/` next to your source changes and commit them
together. A reviewer reading that commit sees the change and what it was for at
once; split across two commits they get a diff whose ticket is still in the state
the previous commit left it. The same goes for a ticket created and handed to
someone else: commit it, or nobody but its author knows it exists. jaira never
commits for you — it reads git and writes only files.

The loop ends past the last lane, not at it: an accepted ticket comes off the
board with `jaira archive <id>` once the work is pushed, and in that order. A
board archived ahead of its push has forgotten a ticket whose code has not
arrived at the teammate who pulls it. An agent does not push on its own
initiative, so it does not archive on its own initiative either; when the user
has pushed, or asks for it, take the ticket off the board. A follow-up keeps its
`follows` link to a logged or archived predecessor, so nothing is lost by clearing them.

When something outside the repository stops the work, park the ticket rather
than abandoning it: `jaira move <id> --to blocked --reason "<what it is waiting
on>"`. The blocked lane refuses a ticket that cannot say what it is waiting on
(`blocked-by` naming another ticket counts as the answer), and parking is exempt
from the leaving lane's output contract — a ticket stopped mid-work has not
produced its output yet, and that is not held against it.

Lanes before the first lane that declares itself the start of the specified
zone (`requires-specified: true` — the built-in `todo` lane, by default) are for
tickets that are not specified yet: a scrap, an idea, something being thought
through. The promotion gate fires when a ticket crosses into that zone, not
when it leaves a lane with a particular name.

## Before creating a ticket

Before `jaira create`, search the board for what it already decided.
`jaira list -q <term> --json` matches over title, goal, context, definition of
done, assignee and status — not just the title — and on a small board a single
`jaira list --json` returns everything in one call. Read anything that looks
close with `jaira show <id> --json` before judging it.

```bash
jaira list -q "session cookie" --json
jaira show <id> --json
```

Related tickets are normal; a contradiction is not — an existing ticket whose
goal, definition of done or context already decided the same question a
different way. When one turns up, stop and put it to the user: name both ids,
quote the line that contradicts, and offer the ways forward — adjust the new
ticket to match the existing decision, create it anyway as a deliberate
supersession with `--follows <id>`, or drop it. A pure duplicate is the same
check's easy case: point at the existing ticket instead of writing a twin.

Nothing in the binary enforces this. It is a judgment call, so it is on the
agent to make it before writing.

## Leaving a trail

Write `jaira note <id> <text>` before you stop, or when something did not
work — anything not written down is lost once the session ends. Three parts:
what you were doing, what you found — especially what did not work — and what
the next step is. Write it as if the reader has mild ADHD and knows none of what
you know: what you were doing first, short concrete lines with one point each,
names and paths rather than adjectives, no jargon and no preamble.

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

- **Leave the human review lane.** It is a human checkpoint. A review agent writes
  `review-summary`, `review-gaps` and `review-verdict`, then moves the ticket to
  human review and stops; a person accepts the work in the board or raises a
  follow-up. Attempting the move exits 3. When the person has reviewed the work
  and says so, the agent may finish the acceptance on their behalf with
  `--force`, which is reported in that command's output rather than written to
  the ticket — the decision stays human, only the keystroke is delegated.
- **Close a ticket without meeting its definition of done.** Every criterion must
  be marked done, the plan finished if the ticket has one, and the commits that
  carry the change recorded on the ticket. There is no evidence flag to pass.
  `--force` exists, is reported in the command's output, and is the user's call
  rather than the agent's.
- **Write to someone else's ticket.** A ticket belongs to its `assignee`; a move
  by anyone else exits 3, naming the owner. The human checkpoint lanes are
  exempt — reviewing and signing off someone else's work is the whole point of
  them — and taking a ticket over with `jaira set <id> assignee=<you>` is always
  allowed, so a ticket is never frozen by an owner who has gone quiet. `--force`
  overrides and is reported in the output, same as any other gate refusal. In the
  board the same override is `f`, which asks once more before it writes.

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

## A note on lane order

A board's column order follows each lane's `after:` anchor, never `precedence`
— unless the project has its own order file (see below), in which case the
order file decides the position of every lane it names, and `after:` is no
longer consulted for those. `precedence` is the rank a merge uses to decide
which lane wins when two clones moved the same ticket — `jaira lanes` and
`jaira lanes show` label it accordingly, not as a position.

A project can also add, remove or reorder the lanes it uses, independently of
its column order file's day-to-day up-keep: `jaira lanes add <id>`,
`jaira lanes remove <id>` (refused, naming them, if a ticket is in it), and
`jaira lanes move <id> --left|--right`. All three take `--json`, and all three
are the exact calls the settings screen's board makes — a lane added,
removed or moved from either place is the same fact from the other. If the
project has no lane directory of its own yet, the first such change writes
one holding every lane the board currently shows, so it never drops to a
single column because one lane was touched.

## Building and sharing a lane without the TUI

`jaira lanes path` finds where a lane file belongs, `jaira lanes template`
prints its skeleton, and `jaira lanes` (or `jaira lanes show <id>`) confirms it
loaded cleanly once written. `jaira lanes use <id>` puts it to work in this
project, and `jaira lanes publish <id>` hands it to teammates, who pick it up
with `jaira lanes adopt <path>` — the path `jaira lanes shared` prints.
