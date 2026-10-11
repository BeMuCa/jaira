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
`out` a directory that does not exist yet):

```json
{
  "repo": "/home/alex/projects/jaira",
  "out": "/tmp/lane-eval-2026-10-11",
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
arm,ticket,run,lane,tokens,duration_ms,tool_uses,loops,verdict
```

- `tokens`, `duration_ms`, `tool_uses` come from the agent's transcript, because
  a workflow's `agent()` returns only the agent's answer. Every arm agent's
  prompt starts with a `[lane-eval id=… arm=… …]` line; that line is how the
  script finds it. Setup, judges and the recorder carry none and are not counted.
  `tokens` is the context of the agent's last turn, which is what the Agent tool
  reports as `total_tokens`, to within about one percent.
- `lane` is the lane the agent worked. The worker arm has one row,
  `in-progress..testing`, plus its `review` row.
- `loops` is per run: how often critique, optimize or testing sent the ticket
  back to in-progress. The worker reports its own.
- `verdict` is per run: `win` or `lose` when both judge passes picked the same
  arm, `tie`, `split` when the two passes disagree (position decided, not
  quality), `unjudged` when a judge or the other arm is missing.

`<out>/summary.csv` has one line per arm: runs, agents, tokens (total and per
run), time, tool uses, loops per run, and the verdict counts.

Which diff was X and which Y stands only in `<out>/KEY.txt`. The judge folders
`<out>/judge/<ticket>-r<run>-p1|p2/` hold `ticket.md`, `X.diff`, `Y.diff` and,
for a golden ticket, `reference.diff`. The judge is told to read nothing else.

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
