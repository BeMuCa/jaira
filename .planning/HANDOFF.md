# Handoff — 2026-08-12

State of jAIra after the v1 build and one adversarial test round. Written before a
context clear; this file is the memory.

## Environment (matters — nothing here is on PATH by default)

- Go 1.26.5 installed to `~/.local/go` (no sudo used). **Not on PATH.**
  Every command needs `export PATH=$HOME/.local/go/bin:$PATH` first.
- No system Go, no Rust, no passwordless sudo.
- GSD CLI is `node ~/.claude/get-shit-done/bin/gsd-tools.cjs <cmd>` — the docs say
  `gsd-sdk query <cmd>`, which does not exist here. Subcommands take no `query`.
- Build: `go build -o ~/.local/bin/jaira ./cmd/jaira`
- When testing, ALWAYS set `JAIRA_HOME` and `JAIRA_LANES_DIR` under `/tmp`.
  Without them, test runs write into the real `~/.jaira`.

## Where things stand

HEAD `f46110c`. 8 roadmap phases implemented; 83 of 84 v1 requirements met
(TUI-11 deferred — no in-place field editor exists, so there is no edit buffer to
clobber). `gofmt`, `go vet`, `go test ./... -race` all clean.

Four adversarial agents ran. Everything below was reproduced first-hand before
being acted on.

## Fixed and verified this round

| # | Problem | Fix |
|---|---|---|
| 1 | A self-written lane with `terminal: true` let an agent mark work done with no evidence, and counted as resolving blockers | `terminal` now implies `requires-outcome` + `requires-nonmodel-signal` unless explicitly opted out |
| 2 | `set` ignored the read-only rule for tickets in unrecognized lanes (`move` honoured it) | `set` refuses too |
| 3 | Missing positional arg exited 1, not the documented 2 | own `exactArgs`/`minArgs`/`noArgs` validators |
| 4 | Frontmatter >16KB: ticket visible in `list` but unresolvable by `show`/`set` — two read paths disagreed | single `readHeader` with full-read fallback |
| 5 | Ticket resolution used the *filename*, so a renamed file became unreachable by any reference | resolve by frontmatter `id` |
| 6 | Merge driver only bound in `jaira init`, so a teammate who cloned and never ran init silently got line-based merging | binds on any command when a committed `.gitattributes` is present |
| 7 | `sync-tasks` overwrote a human-written `goal`; a task could hijack any ticket via `metadata.jaira_id` | sync may only touch tickets it created (`source: agent-task`), and may only seed an *empty* goal |
| 8 | CRLF files lost the CR on the opening `---` | opening line's exact ending preserved |
| 9 | Writing a symlinked ticket replaced the link, forking content | `EvalSymlinks` before the atomic write |
| 10 | Two files declaring one id: one silently shadowed | first wins deterministically, duplicate reported |
| 11 | `~/.jaira` (global config) could be mistaken for a board by `Discover` | a board is identified by `.jaira/tickets/`, not the directory name |
| 12 | Test suite wrote into the real `~/.jaira` | tests set `JAIRA_HOME` before touching the filesystem |

Design changes the user asked for, implemented:

- **Boards start private.** `init` gitignores `.jaira/`; `jaira share` publishes,
  `--undo` reverses. Publishing is a decision, not a default.
- **Ephemeral state left the repo.** Sessions and locks now live in
  `~/.jaira/state/<worktree>/`, so `.jaira/` holds only committed content.
- **Definition of done is a checklist in the body**, not a frontmatter scalar.
  Ticking every box satisfies the Done gate with no `--signal` — a ticked box is a
  human file edit, which is evidence a model cannot manufacture. Verified:
  unticked → exit 3, ticked → exit 0.

## Open — not yet done

1. **TUI layout overflow, 2 bugs.** The TUI agent reproduced these twice, at 20×20,
   including with an ASCII-only fixture (so not a wide-char artifact). **I have not
   reproduced or fixed them.** Its findings report never reached me — only its
   audit did. Either re-derive by rendering at 20×20 and asserting no line exceeds
   the width, or resume agent `ae2c189a2446fce69` and ask for the two findings.
   Not stressed for overflow at all: diff view, projects switcher, message modal.

2. **Migrate the 44 requirementsgenie tickets.** User approved ("you can migrate
   them yes"). Script rescued to `scripts/migrate-tickets.py` (106 lines, written
   by the sync agent, validated against copies — **I have not read or run it**).
   Target: `~/git/requirementsgenie-feature-requirements-coverage-elicitation/tickets/`
   (44 files, branch `feature/requirements-coverage-elicitation`, German, Jira
   frontmatter, no `id`, no `status`, title in an H1, DoD already a checklist).
   Migration must: add `id` (ULID) and `status`, rename to `<ulid>-<slug>.md`.
   jaira already handles the H1-title fallback and the DoD checklist natively.
   Jira fields (`jira`, `type`, `component`, `priority`, `labels`, `epic_link`)
   survive writes byte-for-byte — verified on a real ticket.
   **Work on copies first. Never run `jaira init` in that repo without asking.**

3. **The skill is not installed globally.** It exists only at
   `/home/berk/git/jAIra/.claude/skills/jaira/`, so Claude does not know jaira
   exists in any other repo — including requirementsgenie. Copy to
   `~/.claude/skills/jaira/` to fix. User was asked, has not answered.

4. **Store-agent findings I never verified.** Titles only; the report did not reach
   me: cherry-pick reverting a list deletion; octopus merge failure; an empty
   `conflict-theirs-*` key left behind after resolve. Resume agent
   `adb7f3391f65e3fac` for details, or re-derive.

5. **Sync/share findings I never received** beyond the two I fixed. Areas covered
   were `.gitignore` variants (`**/.jaira/`, `!.jaira/keep`, no trailing newline),
   `share` in a non-git directory, and the tasks↔sync round-trip convergence.
   Resume agent `a6a66be715bf53031`.

## Things not to re-litigate

- Tickets are committed *when shared* — that is the whole team-sync mechanism the
  user specified. Private-by-default changes the timing, not the model.
- The merge driver only matters for shared boards with more than one writer. For a
  solo private board it never runs.
- `assignee` is always a human; `executed-by` records the model. Never reassign a
  ticket to a model.
- ULID: first 10 chars are the timestamp, last 16 are random. The 6-char handle
  comes from the **tail** — an earlier head-based handle collided for tickets
  created in the same millisecond, which is the normal case for agent bursts.
- Adoption risk is real and unproven: file-based trackers have a poor record, and
  Fossil explicitly rejected mutable working-tree ticket files for the same
  conflict reason. The bet is the field-aware merge driver plus a small surface.

## Next command

```bash
export PATH=$HOME/.local/go/bin:$PATH
cd /home/berk/git/jAIra && go test ./... && git log --oneline | head -5
```
