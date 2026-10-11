# Lane eval

Runs the same tickets in two ways of working and compares them:

- **lane-per-agent**: a fresh agent per lane, each one started with the Skill
  `jaira-role-lane <id> <lane>`;
- **worker**: one agent carries the ticket through every lane before `review`.

In both arms a fresh agent then works `review`. A blind judge compares the two
diffs of each ticket and run, twice, with X and Y swapped, because a model
judge leans towards whichever diff it reads first.

This is a developer tool for this repository, not a jaira feature. It is not in
the binary.

## Parts

| File | Does |
|------|------|
| `.claude/workflows/lane-eval.js` | the Claude Code workflow: setup, both arms, review, judge, record |
| `scripts/lane-eval/clone.sh` | one clone per arm, ticket and run; no remote, fixed base, never removed |
| `scripts/lane-eval/judge-prep.sh` | the two judge folders for one ticket and run, plus a line in `KEY.txt` |
| `scripts/lane-eval/metrics.sh` | agent transcripts + `results.json` → `metrics.csv` and `summary.csv` |

## Start it

Only a person's own session can start a workflow. Ask Claude in this
repository to run the workflow `lane-eval` with these args (paths absolute,
`out` a directory that does not exist yet and does not lie under `$TMPDIR` or
`/tmp`, where the judge folders go — `clone.sh` refuses that layout):

```json
{
  "repo": "/home/alex/projects/jaira",
  "out": "/home/alex/lane-eval/2026-10-11",
  "tickets": ["BG5QJ6", {"id": "GDB5YJ", "base": "43080ed"}],
  "n": 3
}
```

`repo` is where the tickets and the helpers come from: the workflow runs
`<repo>/scripts/lane-eval/*.sh`, so it must be a checkout that has them. `arms`
defaults to both, `n` to 3, `maxSteps` (lane agents per run before
lane-per-agent gives up) to 20. Try it on one ticket with `"n": 1` first.

The clones land at `<out>/<arm>/<ticket>/r<run>` with their metadata beside
them in `r<run>.json`. Nothing is pushed and nothing is cleaned up: every
clone can be opened afterwards to see what an agent did.

## Read the result

When the workflow is done, take the transcript directory from its result and run:

```bash
bash scripts/lane-eval/metrics.sh <out> <transcript dir>
```

It writes `<out>/metrics.csv`, one row per agent:

```
arm,ticket,run,lane,input_tokens,cache_write_tokens,cache_read_tokens,output_tokens,cost_eq,peak_context,duration_ms,tool_uses,subagents,loops,verdict
```

- All numbers per agent come from its transcript, because a workflow's
  `agent()` returns only the agent's answer. Every arm agent's prompt starts
  with a `[lane-eval id=… arm=… …]` line; that line is how the script finds it.
  The transcript does not start with it: the workflow harness puts the request
  that started the workflow first and the prompt, indented, in the second user
  message, so the script looks for the marker at the start of any line of the
  first two. The transcripts lie in
  `~/.claude/projects/<project>/<session>/subagents/workflows/<run id>/`, the
  `Transcript dir` the Workflow tool prints when it starts.
  Setup, judges and the recorder carry none and are not counted.
- An agent an arm agent starts itself with the Agent tool carries no marker
  either, but it is part of that arm's cost. Claude Code writes
  `agent-<id>.meta.json` beside every subagent transcript, with
  `parentAgentId`; the script follows it up to the nearest marked agent and adds
  the nested one's tokens, `cost_eq` and tool uses to that row (`peak_context`
  is the larger of the two). `subagents` says how many were added. A session
  started in its own terminal tab (Herdr's `spawn.sh`) is not a subagent, has
  no `parentAgentId` and stays uncounted; the workflow's prompts start none.
- `input_tokens`, `cache_write_tokens`, `cache_read_tokens`, `output_tokens` are
  what the agent consumed: each assistant message counted once (the transcript
  repeats it per content block, under one message id), summed over all turns.
- `cost_eq` is the column to compare arms by: `input + cache_write × 1.25 +
  cache_read × 0.1`, the consumption in input-token prices. Output is left out;
  it is small and listed on its own.
- `peak_context` is the context of the agent's last turn, what the Agent tool
  reports as `total_tokens`. It is a size, not a cost: a long single session has
  one large peak but pays for its whole context again on every turn, so summing
  peaks favours the worker arm. On the hand-run A/B test of 2026-10-11 the
  summed peaks were 1.10 M for lane-per-agent against 0.45 M for the worker
  (2.5 times), `cost_eq` 2.32 M against 1.50 M (1.5 times).
- `lane` is the lane the agent worked. The worker arm has one row,
  `in-progress..testing`, plus its `review` row.
- `loops` is per run: how often critique, optimize or testing sent the ticket
  back to in-progress. The worker reports its own.
- `verdict` is per run: `win` or `lose` when both judge passes picked the same
  arm, `tie`, `split` when the two passes disagree (position decided, not
  quality), `unjudged` when a judge or the other arm is missing.

`<out>/summary.csv` has one line per arm: runs, agents, the subagents added to
them, the token columns and
`cost_eq` summed (and `cost_eq` per run), the largest `peak_context`, time, tool
uses, loops per run, and the verdict counts.

## What keeps the judge blind and the arms in their clones

- The judge folders are fresh temp directories outside `<out>`
  (`$TMPDIR/judge.XXXXXXXX`, one per pass), holding `ticket.md`, `X.diff`,
  `Y.diff`, `X.ticket.md`, `Y.ticket.md` and, for a golden ticket,
  `reference.diff`. The judge's prompt names
  only its folder; neither the folder's path nor anything beside it names an
  arm, a clone or `<out>`. Which diff was X and which Y, and where each folder
  is, stands only in `<out>/KEY.txt`.
- `X.ticket.md` and `Y.ticket.md` are what each arm wrote onto its ticket: the
  definition of done with each item's proof, the outcome, and the commit list,
  every entry with the subject git in the clone gives it, or `(not a commit in
  this repository)`. The judge needs them because the typical mistakes of an
  arm sit there, not in the diff: on the hand-run A/B test of 2026-10-11 the
  lane-per-agent arm left a `--proof` pointing at a line optimize had moved and
  put two SHAs into one commits entry, and a judge reading diffs alone saw
  neither. Notes, author, lane history and timestamps stay out: how an agent
  writes its notes tells the arms apart.
- An arm agent's prompt names no path but its clone, and its first instruction
  is `cd <clone>`. The clone has no remote.

What is left, because a workflow cannot do more:

- `agent()` has no working-directory option. Every agent starts in the session's
  directory — the real repository — and is in its clone only after that first
  `cd`. A forgotten `cd` still writes to the real board.
- Agents use the real `~/.jaira` (roles, notify hook) and the real `~/.claude`.
- Nothing stops a judge from listing `/tmp`; the path just no longer leads it to
  `<out>` or a name that gives the arm away, since an `<out>` under the temp
  root is refused.
- The judge folders are not removed; old passes stay in `$TMPDIR` until it is
  cleaned.

## Golden tickets from the logbook

A ticket filed in `.jaira/logbook/` is a golden ticket: someone accepted the
work, so there is a reference to compare with. Name it by its handle and give no
base:

- the base is the parent of the oldest commit whose subject names the ticket;
- the ticket is filed afresh in the clone, in `todo`, from its title, goal,
  context and definition of done. Its plan, notes, outcome and review stay out,
  so the run cannot copy them. It gets a new handle, which `r<run>.json` records;
- the accepted change (the commits listed on the logbook ticket, without
  `.jaira/`) goes to `<out>/reference/<ticket>.diff` and into every judge
  folder.

Pick tickets whose context can be acted on without the conversation that filed
them, and whose change is code rather than board bookkeeping.

A ticket that is still on the board works too, but needs a `base` in which it
already exists; it starts from the lane it sits in, walked out of backlog.

## Known trap

On this machine about 16 Go tests fail on an untouched master, because they read
`~/.jaira/default-board.md`. Every arm agent is told to set `HOME`, `USERPROFILE`
and `JAIRA_HOME` to one absolute empty temp directory before `go test`; a
relative one breaks the Go module cache.

## Checking the script

A workflow body may `return` and `await` at top level, so `node --check` needs
it wrapped first:

```bash
{ echo 'async function wf(args, agent, parallel, pipeline, phase, log, budget) {'
  sed 's/^export const meta/const meta/' .claude/workflows/lane-eval.js
  echo '}'; } > /tmp/lane-eval.check.js && node --check /tmp/lane-eval.check.js
```
