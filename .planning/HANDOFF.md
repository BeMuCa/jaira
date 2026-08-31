# Handoff — 2026-08-31

State before a context clear; this file is the memory. The previous handoff
(2026-08-26, the board arriving in this repository) is at `7a861b3`. Longer-
lived decisions also live in the agent memory file `jaira-design-invariants`.

## Where things are

- **Public: https://github.com/BeMuCa/jaira** — clean tree, pushed, master at
  `e9cd386`. The user's binary is `~/.local/bin/jaira`; rebuild with
  `go build -o ~/.local/bin/jaira ./cmd/jaira` after every change, never
  `go install`. Check `go version -m` — a build before committing stamps
  `+dirty`.
- `gofmt -l core internal` lists exactly `internal/cli/tickets.go`
  (pre-existing alignment group). `core/gate/gate.go` no longer appears.
- `go test ./... -race` with the cache cleared before claiming anything.
- **The board here is the work list.** `jaira next --per-lane` — not
  `NEXT-STEPS.md`, which now only carries the ground rules.
- Still unreleased: everything after `62989f1` (v0.1.0 is the only tag), now
  ~40 commits including two renames of unreleased surface (`sync`→`logbook`,
  JSON `synced`→`logged`). Ticket `7ZQ0ZN`.

## What landed 2026-08-27 … 08-31 (all pushed)

- **`jaira sync` is `jaira logbook`** (`b93a85c`): dated per-person folders,
  `restore` still reads the old `sync/` name, `archive --help` states the rule
  — finished → logbook, not being worked → archive (`7173a45`, `2MM32Y`).
- **`z` draws empty lanes thin** (`bd5b2ee`): four cells, name vertical;
  `fitWindow` budgets per lane (unit costs = the old centred row, sweep-tested).
  Measured: lipgloss v2 `Width` includes the border.
- **`s` on the home screen**: logbook entries per day, last 7 days, all boards,
  from folder names alone (`b91ecb1`, `SPDWGH`).
- **`jaira update` no longer re-privatises a shared board** (`ef8f2d9`,
  `ETR0PX`): it writes the agent note only, never .gitignore.
- **A board is its lane directory** (`743737f`, `BNZERQ`): Load injects
  nothing; first open writes the default board or the built-ins as files plus
  `order`, once, reported; `removed`/`MaterialiseWorkingSet`/`Differs` and the
  override warning are gone; legacy dirs (a `removed` file or no `order`)
  migrate in place once. This board migrated itself. **His req board migrates
  on its next jaira command** — 8 shipped files written beside its 3, wanted.
- **`jaira lanes market`** (`0782f6e`): the repo's `lanes/` on GitHub, listed
  via the contents API and adopted with `market adopt <id>`; README +
  `lanes/README.md` say a PR with one file adds yours, CI parses it.

## The user's req board (other repo, never probe with live moves)

He committed the untracking himself (`c688a93b`): `.jaira/lanes/` ignored and
untracked. Its CLAUDE.md still carries the **old** agent block — `jaira update`
there is now safe (ETR0PX) and brings the board-aware block + the lane-dir
migration. His local "How that works here" section still says `.jaira/sync/`.

## The board here (private, gitignored)

At signoff, his to accept: `4DQPMS QPJNQP 2MM32Y SPDWGH ETR0PX BNZERQ YM7QSA`.
At human: `SGPDYK` (done-cap: logbook vs archive as target, `holds:` on the
lane vs code, oldest = updated-at). Blocked: `QA3GN1` (subsumed by BNZERQ,
close when accepted). Backlog worth knowing: `88H1P4` (sessions drive lanes —
mechanism a: nudge in move output; block sentence for subagents awaits his go),
`YBC0MT` (secrets-scan + changelog-writer lanes), `FCMP17`/`NFJCTK`/`B4MGTP`
(artifacts, lane schema doc, foldable history — all folded into the schema
brainstorm), `MFD7P3`, `GEC3TK` (board.png stale), `7ZQ0ZN` (cut v0.1.1),
`TQXBY5`, `DNAEPN`, `CD9TCB` (mostly resolved by BNZERQ — verify then close).

## Known gaps, written down deliberately

- `human` is marked "never out" in the block but only `signoff` carries
  `requires-human-exit`; the gate lets an agent leave `human`. Proven with
  `--dry-run`. Not ticketed yet.
- Lane files ignore unknown frontmatter keys silently (typo in
  `output-produces` never surfaces). Goes with `NFJCTK`/the schema.
- The "no subagents" line is this Claude Code build's default prompt, not any
  file on the machine; a request in CLAUDE.md overrides it per session.

## The schema brainstorm (running, architectural path)

State in `.planning/schema-brainstorm.md`: decisions so far, the worked
example, and the one open question (frontmatter = latest + `## History` per
round, vs lane fields only in History). Spec not yet written — it goes to
`docs/superpowers/specs/` only after he approves the design. Lane inspiration:
`.planning/research/lane-ideas.md` (19 ideas; he likes `secrets-scan` and
`changelog-writer` → `YBC0MT`).
