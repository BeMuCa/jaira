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

`jaira next` answers "what is the single furthest-along ticket"; `jaira next
--per-lane` answers "where is work waiting", one entry per lane that has any, in
pipeline order, each carrying the lane's `agentic` flag and one ticket. The two
are different questions: under the default ordering a deep queue in a late lane
hides every earlier lane, so a step inserted mid-pipeline never sees traffic
until the queue ahead of it drains.

Every ticket in a `--json` payload carries `next_lane`: where it goes when the
step it is in is finished, with the ticket's own Options applied and parking and
question lanes left out. Empty means there is nowhere left to go, and also
whenever the ticket is parked or waiting on an answer: such a ticket resumes
where it stopped, which the board does not record. It is a
convenience, not a rail — `jaira move` re-checks the gates whatever route the
caller took. The same payload carries `review` (`summary`, `gaps`, `verdict`,
`check`), so a reader can be handed the review without opening the board.

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
| `jaira next` | the next actionable ticket; `--lane`, `--assignee`, `--all`, `--per-lane` |
| `jaira lanes` | the installed lanes |
| `jaira projects` | boards you have opened |
| `jaira whoami` | the identity jaira acts as, and the other names that mean you |
| `jaira sessions` | sessions working this tree |
| `jaira resume` | everything left mid-flight, with its notes |
| `jaira validate` | check every ticket for damage; `--strict` fails on warnings |

### Being one person under several names

A ticket belongs to its `assignee`, and the gates compare that string to who is
acting. A person is not one string: a `user.name` on one machine, a work address
in a ticket a teammate's tooling assigned, a personal address in another
repository. jaira treats git's `user.email` as you as well as `user.name`, and
reads `~/.jaira/identity` (one name per line, `#` comments allowed) for the rest.
`jaira whoami` prints the resolved list.

Names not on that list are somebody else, so their tickets are refused. That is
the point of the rail — but a rail that refuses your own tickets only teaches you
to pass `--force` to everything, which protects nothing.

## Writing

| Command | What it does |
|---|---|
| `jaira init` | prepare a repository; writes a jaira section into `CLAUDE.md` |
| `jaira update` | re-apply this repository's jaira setup and print what changed since the version that last did it |
| `jaira create <title>` | create a ticket; `--goal`, `--context`, `--dod`, `--assignee`, `--lane`, `--tier`, `--blocked-by`, `--follows` (the ticket this one follows on from; must resolve) |
| `jaira set <id> k=v…` | set frontmatter fields |
| `jaira dod <id> <n> --doing\|--done\|--todo` | mark a checklist item |
| `jaira dod <id> --add "…"` | append checklist items; repeat for several |
| `jaira dod <id> --plan …` | address the Plan checklist instead of the definition of done |
| `jaira dod <id> --option <name>` | turn an optional step on for this ticket (`--todo` turns it off) |
| `jaira move <id> --to <lane>` | move lanes, applying the gates; `--what`, `--why`, `--resolves`, `--commits` (leaving the implementing lane), `--question` (entering the human lane), `--reason` (entering the blocked lane), `--from-lane` (validate piped lane output), `--force` |
| `jaira note <id> "…"` | record progress a later session would otherwise rediscover |
| `jaira claim <id>` | take a 30-minute lease so two sessions do not collide |
| `jaira archive <id>` | take a ticket off the board (nothing is deleted) |
| `jaira restore <file>` | put an archived ticket back |
| `jaira resolve <id>` | settle the fields a merge could not |
| `jaira share` | publish the board; `--undo` makes it private again |
| `jaira projects add <path>` | register a board; `--scan` searches two levels down |
| `jaira lanes use <id>` | copy a catalogue lane into this project; `--force` |
| `jaira lanes add <id>` | add a built-in or catalogue lane to this project's board, appending it to the column order; materialises the project's lane directory first if it has none yet |
| `jaira lanes remove <id>` | remove a lane from this project's board (it stays in the catalogue); refused, naming them, if any ticket sits in it |
| `jaira lanes move <id> --left\|--right` | shift a lane one column in this project's order, swapping it with its neighbour |
| `jaira lanes publish <id>` | copy a lane into `.jaira/shared/<you>/` for teammates; `--force` |
| `jaira lanes adopt <path>` | copy a teammate's shared lane (the path `lanes shared` prints) into your catalogue; `--force` |
| `jaira lanes default` | show or set which lanes and options a new board starts with; `--lanes`, `--options`, `--clear` |

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
j k ↓ ↑   card            /       filter (key:value narrows to one field)
g G       first / last    m       move ticket      ?   help
v         compact view    x       archive          r   reload
q         quit            S   settings: lanes and the default board
```

In an open ticket:

```
e   edit fields (enter newline, ctrl+s save)   a   accept (at a checkpoint)
E   edit body and checklists in $EDITOR        f   follow-up from the review
y   copy the full ticket id                    n   follow-up beside this one
b   open the ticket this one is blocked by     m   move it
jk  next / previous ticket                     tab other pane, in the split
↓↑ scroll (ctrl+d/u pages) a ticket taller than the terminal
```

`n` splits the screen: the ticket it follows on the left, the follow-up written
on the right. Nothing is written until `ctrl+s`, which creates it in the default
lane with `follows` set; `esc` discards the draft. `n` again chains from the
ticket just written. The editor keeps `tab` for its fields, so `shift+↓↑` scrolls
the left pane while you type; after saving, `tab` moves between panes. No split
below 80 columns or 20 rows.

A refused move offers `f` to override it and asks again before writing, which is
the TUI's `--force`: any refusal, reported in the output, nothing written to the
ticket.

Lane settings (`S` then lanes): the project's lanes drawn as a small board,
with a `+` column at the far right that opens the catalogue.

```
h l       select a column        x   remove the selected lane from this project
H L       move the selected      enter (on +)   catalogue: choose a lane to add
          lane one column        u   copy the selected lane into this project
tab       switch to teammates'   p   publish it to teammates
          shared lanes           n   write a new lane and open it in $EDITOR
E         edit the selected lane in $EDITOR (a built-in gets an override copy first)
                                  R   pull a drifted lane's catalogue copy in
                                  a   (shared list) adopt a teammate's lane
esc       back
```

Home screen: `enter` open · `a` add a board · `r` refresh · `q` quit.
In the picker: `a` add · `i` create a board here · `s` scan two levels · `space`
toggle a result · `enter` add what is selected.

Compact view: `a`/`d` or `←`/`→` move · `enter` open the step full width ·
`1`-`9` switch project · `v` back to the board.

In an opened step: `jk` ticket · `gG` first/last · `hl` next/previous lane
(stays in this view) · `enter` open ticket · `q`/`esc`/`v` back to the compact
view.
