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
