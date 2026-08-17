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
