<!--
Format rules — read before editing:
  - Newest release first.
  - A release starts with a line reading exactly "## " followed by the exact
    version string a build of that release reports (the git tag with its
    leading v stripped, e.g. tag v0.1.0 reports "0.1.0").
  - Inside a release, every change is exactly one line starting "- ".
    Never wrap a change across two lines — the parser is a line scan, and a
    wrapped line is read as two separate, incomplete changes.
  - Anything else (this comment, blank lines, prose) is ignored, so it is
    free to use for context.
  - Write each change as an instruction: what the reader must DO differently,
    not a record of what commit did what.
-->

## Unreleased
- Expect `jaira logbook --all` and a lane's `holds` cap alike to leave tickets that live only on a ref where they are, instead of naming each one as a problem they could not file: a ticket with no file on this disk is not yours to file, so a lane full of other people's ref-only tickets no longer trims your own off the board, and the count the board reminds you with now matches the set the command takes.
- Run the command the board now names: once ten or more finished tickets sit in the terminal lane, the hint bar reads `⌸ N to file: jaira logbook --all` instead of only the count, so you can act on the reminder without looking the command up.
- Read the mode off `jaira create` instead of guessing: a ticket that stays a file on your disk now says so on the same line, names the remote it looked for and why it is not usable, and `--json` carries that same sentence as `file-only-reason` beside `on-ref-only`. Outside a git repository it says that instead, and names no remote and no way back, because neither exists there.
- Put a file-only ticket back on its ref with `jaira release <id>`: give the repository the remote it was missing, then release, and the ticket travels on its ref while the local file goes — no need to recreate it under a new id. Where `create` advises it on a ticket that already has an assignee, it now also says that releasing clears them.
- Ask `jaira whoami` where your tickets go: inside a board it now also prints the remote this board uses and where that name came from, the remotes the repository actually has, whether tickets travel on refs (with the reason when they do not), and how many tickets exist on this disk and on no ref. `--json` adds `board`, `remote`, `remote_source`, `remotes` (always a list, empty rather than null), `ref_mode`, `ref_mode_reason` and `file_only`. The remote it names is the one the ref commands actually use, and the reason is a sentence rather than a raw `gitref:` error. Outside a git repository it names no remote at all, in the text and in `--json` alike, because there is none to name.
- Pass `--dod` to `jaira create` once per acceptance criterion, the way `--tag` already works: every one becomes its own checkbox in the ticket body and the first also fills the `definition-of-done` frontmatter, so a ticket with six criteria no longer arrives carrying one and no longer has to be pulled to gain the rest.
- Set the remote for one board with `git config jaira.remote <name>` instead of letting `~/.jaira/settings.json` decide for every checkout on the machine: the board setting wins and never falls back, `"remote"` in settings.json now applies only where the repository really has a remote by that name, a repository with exactly one remote uses it without any configuration, and a ref command that still cannot resolve a remote names the remotes the repository does have and the command that fixes it.
- Read a card's left bar as three tags instead of one: its three cells now show the colours of the ticket's first, second and third tag, top to bottom in the order the tags stand on the ticket, so a board can be read along three axes at once; a cell with no tag of its own or whose tag has no line in `.jaira/tags` stays in the lane's shade, and tags past the third stay on the ticket and in `jaira show` but colour nothing.
- Expect a multi-line ticket title to read as one line on its board card: the break becomes a space instead of adding a fourth row to a three-row card, which used to push the cards under it out of the lane's row budget.

## 0.2.0
- Stop inventing your first notification hook: `jaira hook example` prints a working script for the `"hook"` setting that stays silent for every agent lane and rings the terminal bell only when a ticket reaches a person (`human`, `signoff`) or finishes (`done`), telling the two apart by the line it prints, and says which line to replace to deliver through ntfy.sh or a desktop notifier instead.
- The critique lane now reads the whole diff on its first pass; later passes only re-check what it already found, so a review loop cannot run for ever by reading deeper each round.
- Keep a ticket's worktree until the ticket reaches `done`, not until its pull request merges: both shipped prompts now name the same moment, where `jaira-teamlead` used to name the earlier one. Re-run `jaira roles install --global --force` to pick it up.
- Put every worktree in `.worktrees/` beside the repository, never inside it: the `jaira-dispatcher` prompt now derives the path and `scripts/spawn.sh` creates it there, so `grep -r`, `find` and `ls -R` from the repository root can no longer walk into a second copy of the sources on someone else's branch and read the wrong file. Re-run `jaira roles install --global --force` to pick it up.
- Re-run `jaira roles install --global --force` to pick up the sharpened `jaira-dispatcher` and `jaira-teamlead` prompts: a dispatcher now stops after three round-trips of one lane without arguing itself into a fourth, reports one line per finished lane instead of only at the end, closes each worker's tab as soon as its lane is read off the board, and leaves a worker alone when a human types in its tab; a teamlead now closes the dispatcher tab it opened once the pull request is open, and starts a fresh session rather than resuming a finished one. From here these prompts are edited in `core/role/builtin` and installed from there, never the other way round.
- Install the agent role prompts this binary now ships: `jaira roles install --project` writes them into `.claude/skills` so they arrive with a clone, `--global` into `~/.claude/skills`, `--into <dir>` anywhere else; `jaira roles list` names them. A file you have edited is reported and left alone, and the command exits 3; `--force` replaces it.
- Invoke the shipped roles under their `jaira-` prefixed names — `/jaira-teamlead`, `/jaira-role-lane` and the rest. If you wrote these prompts by hand before under their bare names, those copies still answer to the old name; delete them yourself once you have switched.
- Steer the board with a Cyrillic keyboard: letter commands are read by the key's position, so `j`, `k`, `g`, `q` and the rest of them work through a Russian layout, where the board used to answer nothing at all. The key that is `/` on a US keyboard opens the filter under its Cyrillic-layout character too. Typed text is untouched — Russian still goes into filters and edit fields as Russian.
- On a terminal that speaks the Kitty keyboard protocol (kitty, ghostty, WezTerm, foot), letter commands work under any layout, not only Cyrillic: the board asks the terminal for the physical key behind each press. Punctuation and shifted keys are left as the layout printed them, so `?` still opens the help; keys held with ctrl or alt are the terminal's business as before.

## 0.1.4

- Expect one snapshot and one fetch per clone rather than per checkout: their clocks now live with the repository, so adding a git worktree no longer starts a fresh interval and no longer writes a snapshot commit on its first command.
- Read whole titles: cards have no frame any more. Each one is a filled band across the lane with its tag's colour in a solid cell down the left, and the shade alternates so stacked cards stay apart. Titles gain two columns and every card gives two rows back to the lane, so more of the lane fits on screen.
- Press `c` to fill the selected card in its tag's colour instead of neutral grey, and `c` again for plain. It ships on; where your terminal shows 24-bit colour the tint is a quiet one, and on a palette-only terminal it is mixed further so the tag's colour survives rather than snapping to grey.
- Look for the filled card to find the cursor: the selected ticket is painted a clear step above the lane's own shades, instead of only marking itself with a bold title that vanished on a column of coloured cards. On a light terminal the fill goes the other way, darker rather than lighter.
- Finishing a ticket no longer files anything: it stays in `done` with everybody else's until you cut. Run `jaira logbook --all` when you account for your hours and the whole lane goes into today's folder; `jaira logbook <id>` still files one, and `jaira logbook` alone still only lists. A board holding more than ten finished tickets says so in its hint bar and files nothing on its own.
- If you want the old doorway back, set `logbook-on-entry: true` on your terminal lane yourself — it ships off, because filing on entry means finishing one ticket files everybody's.
- Work in the board reaches the team by itself now: it sends what it queued on its own background run, so a day spent in the board no longer leaves everybody else looking at yesterday's tickets.
- A ticket you file into the logbook stays off the board — it used to come back as a card, because its ref outlives the filing by design and nothing skipped it.
- Re-read the `jaira dod` lines in the jaira section of `AGENTS.md` and `CLAUDE.md` (run `jaira update` to refresh it): they now say that the numbered criteria are what the terminal lane's gate reads, and that `--plan` is a second, separate list whose ticks do not count towards it.

- Dialogs now float over the board instead of replacing it: the link window, the tag legend and every refusal or note are drawn as a centred box with the board still visible behind, so you keep the card and the lane you were looking at.
- Press `L` on a card, or on an open ticket, to see every ticket linked to it — what it waits on, what waits on it, what it is part of, what it contains, what it relates to and what follows it — with the logbook and the archive searched too, so a link no longer dies when the work behind it finishes. `enter` jumps to a linked card — and opens it, when you pressed `L` while reading a ticket — and `esc` puts back exactly the screen you came from.
- Record containment with the new `parent` field: a ticket names the one it is part of (`jaira set <id> parent=<id>`, or `jaira create --parent <id>`), and children — and their children, to any depth — are read back from that. There is no `children` field to keep in step.
- Record a loose connection with the new `related` field (`jaira set <id> related=<id>,<id>`, or `jaira create --related <id>`). Write it on either side; both sides show it.
- Run `jaira links <id>` for the same picture on the command line, with `--json` for a machine.
- `jaira show <id>` now prints what a ticket is part of, what it contains (children at every depth, the filed ones included) and what it relates to; `--json` carries the same as `parent`, `related` and `filed_away`.
- Write a link with a handle: `jaira set <id> blocked-by=<handle>`, `parent=` and `related=` now resolve whatever reference you can read off the board into the full ticket id, and refuse a reference that names no ticket. A handle used to be stored verbatim and the link then resolved to nothing forever after.
- Stop working around a blocker that finished: a `blocked-by` whose ticket has been filed into the logbook, or archived from a terminal lane, now counts as cleared instead of blocking forever and dropping the ticket out of `jaira list --actionable`.
- `jaira validate` now reports a dependency or a parent as dangling only when the id exists nowhere at all, and reports a ticket that is its own parent or sits in a parent ring as an error.
- `jaira show <id>` now finds a ticket that has left the board, printing where it is filed instead of "not found".

## 0.1.3

- Tickets now travel on a git ref of their own (`refs/jaira/tickets/<id>`), so a ticket reaches whoever it is assigned to without anybody sharing a branch: run `jaira fetch` after cloning, and again whenever you want to see what the team has moved.
- Run `jaira pull <id>` to take a ticket over: it records you as its assignee on the ref and only then writes the file here, so exactly one clone holds a ticket at a time — a ticket somebody else has pulled is refused, naming them, and `--steal` takes it anyway and writes a note on the ticket saying so.
- Run `jaira release <id>` to hand a ticket back: it clears you as assignee on the ref and removes the local file, and until you do, nobody else can pull it — an assignment is a reservation.
- On a board with a remote, `jaira create` no longer leaves a ticket file on your disk: the ticket lives on its ref until somebody pulls it, so commit nothing for a ticket you only filed. A board with no remote behaves exactly as before.
- `jaira list`, `jaira next` and `jaira show` include tickets that are still on their refs and mark them `[pull it]`; every write refuses them with exit 3 and names the command that changes that.
- The refs are now brought up to date by themselves every ten minutes — in a detached process beside whatever command you ran, and on the open board — so a ticket assigned to you reaches you even if you never open the board and only ever read; no command waits for it, `"fetch-every": "10m"` in `~/.jaira/settings.json` changes the interval for both, and `JAIRA_NO_FETCH=1` turns it off.
- A ticket that arrives while the board is open says so on a line of the board, not on a screen you have to dismiss — it was the other way round and it interrupted whatever you were doing.
- `jaira logbook` and `jaira archive` no longer remove the ticket's ref — they write its final state to it, so the ticket stays visible to everybody while your branch waits to be merged; the ref is cleared later, once the ticket has arrived in a landing branch.
- Run `jaira snapshot` to write the board as files to the `jaira/board` branch and clear the refs of landed tickets; it otherwise happens by itself in the background every three days, so a first-time cloner has the board even before fetching refs.
- Set `landing-branches` in `~/.jaira/settings.json` (a list, globs allowed) to say which branches mean a ticket has arrived — without it the remote's own HEAD is used, and if even that cannot be resolved no ref is ever removed.
- A finished ticket that never reached a landing branch is named by `jaira fetch` and `jaira validate` after three days (`"landing-grace": "72h"` in `~/.jaira/settings.json`); push the branch that holds it, or give up on it with `jaira snapshot --drop <id>`.
- A ticket newly assigned to you raises a desktop notification once; turn it off with `"notify-off": true` in `~/.jaira/settings.json`.
- Point `hook` in `~/.jaira/settings.json` at a script to be told about `move` and `claim` immediately: it gets `JAIRA_EVENT`, `JAIRA_TICKET`, `JAIRA_TITLE`, `JAIRA_STATUS`, `JAIRA_ASSIGNEE`, `JAIRA_ACTOR` and `JAIRA_ROOT`, and what it does with them — Slack, ntfy, Telegram — is yours.
- A write made with no route to the remote is kept and sent with your next command; the card says `⇅ unsent` until it goes, so being offline no longer looks like a lost write.
- Merging two branches that each created the same ticket file no longer loses one side silently: an add/add merge is resolved field by field like any other, keeping the further lane and both sides' lists.
- `jaira validate` and `jaira list` now report a ticket that is on the board and filed away in the logbook or archive at once, and a ticket somebody else has taken off the board while a file for it is still here — nothing is moved for you, because which copy is right is yours to say.
- Set `remote`, `snapshot-branch`, `snapshot-every`, `fetch-every` and `landing-grace` in `~/.jaira/settings.json` if the defaults (`origin`, `jaira/board`, `72h`, `10m`, `72h`) do not suit — the three intervals are durations, so `3s` and `72h` are both said the same way, and a missing, damaged or unreadable value means the default rather than a refusal to start.
- Read the new README sections before using any of this: what git is actually doing underneath, how two people hand a ticket over, and three recordings of it happening between two machines.

## 0.1.2

- Look for the running version in the top-left corner of the board and the launcher; the footer no longer carries it, so a script or a screenshot that read the last line for it must read the first instead.
- A source build now names itself `jaira dev` in that corner instead of staying silent, so you can tell which binary is running when several are installed — this deliberately reverses the 0.1.1 rule that the line only speaks when it can name a release.

## 0.1.1

- Tickets carry `tags`: run `jaira tags` FIRST to see the board's vocabulary and reuse a name for that subject rather than inventing a synonym, then `jaira tag <id> <name>...` or `jaira create --tag <name>`; filter with `jaira list --tag <name>` or `tag:<name>` in the board's `/` filter, and hand-edit the shared colours in `.jaira/tags`.
- Every board's generated agent block is now out of date, so `jaira validate` reports `AGENTS.md`/`CLAUDE.md` as stale until you run `jaira update` once per board and commit the result.
- A hand-written `[-]` in any checklist now reads as withdrawn, not open: it stops blocking completion and nothing reports it as done — retick it to `[ ]` if the item is still wanted.
- `jaira sync` is now `jaira logbook`, and finished tickets file under `.jaira/logbook/<you>-<date>/`; `restore` still reads the old `sync/` folder, and the JSON field `synced` is now `logged`.
- A board is its lane directory: the first command on a board writes the default board or the built-ins as lane files plus an `order` file, once, and says so; a legacy board (a `removed` file, or no `order`) migrates in place on its next command — expect that one-time write and commit it on shared boards.
- `z` now draws an empty lane four cells thin with its name vertical instead of hiding it; press `z` again to widen.
- `s` on the project screen charts logbook entries per day over the last seven days, across all boards.
- `jaira update` on a shared board leaves `.gitignore` alone; only `jaira init` gitignores `.jaira/`.
- `jaira lanes market` lists the lanes published in the project's GitHub `lanes/` catalogue and `market adopt <id>` copies one into yours; `secrets-scan` and `changelog-writer` ship there — a freshly added lane lands as the rightmost column until you move its line in `.jaira/lanes/order`.
- Every screen wraps long text to the terminal width instead of cutting it off at the right edge: paths break mid-word, checklist items and proofs wrap on the sign-off screen, and the key-hint footers wrap so no key disappears on a narrow terminal.
- A source build shows no version line in the footers; the line only speaks when it can name a release.
- A field an installed lane declares and nobody has filled shows as `— owed by <lane>` in the detail and sign-off panes — an unworked ticket in review no longer looks like a finished one.
- `rejects-to:` may name more than one lane (`rejects-to: [in-progress, human]`), and everything that renders back edges says "in-progress or human".
- `jaira move` into an agentic lane names the command that works it, and `jaira hook print` emits a Claude Code Stop hook that refuses to end a session while an agentic lane still holds waiting work.
- `jaira validate` warns when a ticket's context or a note names another ticket's handle that is not in `blocked-by`, suggesting the exact command to declare it (`--strict` turns the warnings into an error).

## 0.1.0

- The done lane refuses a ticket that records no commits. Record them on the move out of implementing — `jaira move <id> --to review --commits "$(git rev-parse HEAD)"` — so the diff at review and sign-off is the diff of exactly those commits.
- The blocked lane refuses a ticket that cannot say what it is waiting on: pass `--reason "…"` on the move, or record the blocking ticket in `blocked-by`. Parking is exempt from the dependency check and the leaving lane's output contract.
- `follows:` links a follow-up to the ticket whose review produced it — settable with `jaira create --follows <id>`, visible in `show`, the board and `--json`.
- The open ticket clips to the terminal and scrolls: arrows line by line, `ctrl+d`/`ctrl+u` by page. The sign-off screen lays out label left, text right, like the detail pane, and `b` jumps to the ticket this one is blocked by.
- The board filter understands `key:value` — `assignee:berk`, `lane:review`, `ticket:<id>` — and plain text still searches everything.
- New tickets seed a `## Progress` section, which is where `jaira note` has always written; the seeded-but-never-used `## Notes` heading is gone.
- A board created from the TUI browse screen now gets the same setup `jaira init` gives it — `.jaira/` gitignored and the jaira section written into `AGENTS.md` and `CLAUDE.md`. Boards created that way before this release have neither; run `jaira update` in each of them.
- The jaira section in `AGENTS.md` and `CLAUDE.md` now names the full working loop rather than a handful of commands. Re-read it — it is the contract for how to drive the board.
- The ticket detail pane shows the full ticket id and `y` copies it. Use the full id when handing a ticket to another tool; a prefix can turn ambiguous as the board grows.
