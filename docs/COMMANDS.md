# Command reference

Every read command takes `--json`. Every command takes `-C <dir>` to run as if
started somewhere else.

## Exit codes

Agents branch on these, so they are a contract rather than an accident.

| Code | Meaning |
|---|---|
| 0 | success |
| 1 | unexpected error |
| 2 | usage error — bad flags or arguments |
| 3 | a gate refused the operation |
| 4 | unresolved dependencies |
| 5 | no such ticket, or an ambiguous id prefix |

Under `--json`, a refusal is structured on stderr with a `code` and often a
`field`, so nothing has to be parsed out of a sentence:

```json
{"error":{"reason":"gate_refused","code":3,"violations":[
  {"code":"needs_nonmodel_signal","field":"definition-of-done",
   "message":"the definition of done is not met: criterion 2 (\"documented\") is still open…"}]}}
```

## Looking

| Command | What it does |
|---|---|
| `jaira` | the home screen: every board, and what each needs |
| `jaira board` | open the board here directly |
| `jaira list` | list tickets; `--lane`, `--assignee`, `--query`, `--actionable` |
| `jaira show <id>` | one ticket in full |
| `jaira show <id> --for-lane <lane>` | the prompt and bounded input a lane's agent should get |
| `jaira next` | the next actionable ticket |
| `jaira lanes` | the installed lanes |
| `jaira projects` | boards you have opened |
| `jaira sessions` | sessions working this tree |
| `jaira resume` | everything left mid-flight, with its notes |
| `jaira validate` | check every ticket for damage; `--strict` fails on warnings |

## Writing

| Command | What it does |
|---|---|
| `jaira init` | prepare a repository; writes a jaira section into `CLAUDE.md` |
| `jaira create <title>` | create a ticket; `--goal`, `--context`, `--dod`, `--assignee`, `--lane`, `--tier` |
| `jaira set <id> k=v…` | set frontmatter fields |
| `jaira dod <id> <n> --doing\|--done\|--todo` | mark a checklist item |
| `jaira dod <id> --add "…"` | append checklist items; repeat for several |
| `jaira dod <id> --plan …` | address the Plan checklist instead of the definition of done |
| `jaira dod <id> --option <name>` | turn an optional step on for this ticket (`--todo` turns it off) |
| `jaira move <id> --to <lane>` | move lanes, applying the gates |
| `jaira note <id> "…"` | record progress a later session would otherwise rediscover |
| `jaira claim <id>` | take a 30-minute lease so two sessions do not collide |
| `jaira archive <id>` | take a ticket off the board (nothing is deleted) |
| `jaira restore <file>` | put an archived ticket back |
| `jaira resolve <id>` | settle the fields a merge could not |
| `jaira share` | publish the board; `--undo` makes it private again |
| `jaira projects add <path>` | register a board; `--scan` searches two levels down |

## Agent plumbing

| Command | What it does |
|---|---|
| `jaira tasks` | emit the board as a task list an agent can adopt |
| `jaira sync-tasks` | mirror an agent's task list into the backlog |
| `jaira checkpoint` | record what this session is working on |
| `jaira merge-driver` | called by git; not run by hand |

## Keys

Board:

```
h l ← →   lane            enter   open ticket      n   new ticket
j k ↓ ↑   card            /       filter           m   move lane
g G       first / last    ?       help             r   reload
v         compact view    x       archive          q   quit
```

In an open ticket:

```
e   edit fields (enter newline, ctrl+s save)   a   accept (at a checkpoint)
E   edit body and checklists in $EDITOR        f   raise a follow-up
```

Home screen: `enter` open · `jk` move · `a` add a board · `r` refresh · `q` quit.
In the picker: `a` add · `i` create a board here · `s` scan two levels · `space`
toggle a result · `enter` add what is selected.

Compact view: `a`/`d` or `←`/`→` move · `enter` open the step · `1`-`9` switch
project · `v` back to the board.
