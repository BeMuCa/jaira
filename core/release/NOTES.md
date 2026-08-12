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

- A board created from the TUI browse screen now gets the same setup `jaira init` gives it — `.jaira/` gitignored and the jaira section written into `AGENTS.md` and `CLAUDE.md`. Boards created that way before this release have neither; run `jaira update` in each of them.
- The jaira section in `AGENTS.md` and `CLAUDE.md` now names the full working loop rather than a handful of commands. Re-read it — it is the contract for how to drive the board.
- The ticket detail pane shows the full ticket id and `y` copies it. Use the full id when handing a ticket to another tool; a prefix can turn ambiguous as the board grows.
