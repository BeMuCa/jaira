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
