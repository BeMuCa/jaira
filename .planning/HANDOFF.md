# Handoff — 2026-08-18

State before a context clear; this file is the memory. The previous handoff
(2026-08-15, gates for commits and blocked reasons) is at `364a260`. Everything
below landed after it, is pushed, and CI is green on linux/macos/windows at
`bf2c16d`. Longer-lived decisions also live in the agent memory file
`jaira-design-invariants` (twelve entries).

## Where things are

- **Public: https://github.com/BeMuCa/jaira** — `master` = origin, clean tree.
- **v0.1.0 is released**: tag + GitHub release with 6 binaries + checksums,
  built by `.github/workflows/release.yaml` (tag `v*` → goreleaser).
  `scripts/install.sh` is the curl-able installer (checksum-verified, tested
  end to end against the real release). **Everything after `62989f1` is not in
  any release** — cutting v0.1.1 means: add a `## 0.1.1` block to
  `core/release/NOTES.md`, tag, push the tag.
- The user's binary: `~/.local/bin/jaira`, built from `bf2c16d`.
  `go build -o ~/.local/bin/jaira ./cmd/jaira` — never `go install`.
- CI was red for ~12 runs on Windows only; fixed by `.gitattributes` pinning
  `core/lane/builtin/*.md` to LF (autocrlf handed go:embed CRLF bytes).

## What landed since the last handoff

Gates and flow
- **Commits are demanded at done, not on leaving implementing** (`a349673`,
  user directive): `requires-commits: true` on 50-done.md, CodeNeedsCommits,
  `--commits` still encouraged in the implementing prompt. Review accepts
  uncommitted work; nothing is accepted that cannot be checked.
- **Capture belongs to nobody; the pull claims** (`6b77a3b`, user's team
  flow): create leaves assignee empty (CLI and TUI; template/--assignee win),
  moving an unassigned ticket assigns the mover — staged before the gate so
  the promotion gate is satisfied by the pull. Existing assignees are never
  overwritten. Team ritual: pull → drag into todo (= claim) → push.
  Teammates' claims render as a yellow `@name` on the card.
- An agent may finish a human's acceptance with `--force` when the human said
  so in conversation — documented in SKILL.md/AGENTS.md, gate unchanged.
- DoD completeness is deliberately NOT required to enter review (review judges
  whether the DoD is met; DoD items like "reviewed and merged" would deadlock).
  Enforced at done only.

TUI
- **Board switching swaps the store in place** (`c310cdb`): no more program
  restart per switch, so the terminal no longer flashes through. board.go/
  home.go lost their restart loops; Model.SwitchTo is gone.
- **Nothing off-screen goes unnamed**: key hints wrap onto more lines instead
  of truncating/dropping (every footer + board status bar, which the board
  measures to size its columns); lane-settings columns wrap into rows;
  catalogue lanes not on the board show as dimmed "not on board" columns
  (enter installs); `E` edits a lane in $EDITOR (a built-in gets a catalogue
  override copy first). Footers name actions only — movement keys
  (hjkl/jk/arrows/wasd) live in `?` help.
- **Every screen uses the full terminal width** (`d8e026c`): the 78-column
  caps and the 34-column board-column cap are gone; visible columns stretch.
- Arrows scroll an open ticket (detail + sign-off, both clipped to the
  window); j/k switch tickets; `b` walks to the blocker (footer offers it only
  when one exists). Sign-off renders label-left/text-right like the detail
  pane (labels: problem/what/why/resolves/summary/gaps/verdict).
- **Body sections render as fields** (`6b77a3b`): options/progress/whatever
  headings remain render label-left/content-right in the detail pane; raw
  `##` no longer leaks; empty sections are skipped. Storage unchanged.
  The CLI's `show` still prints the body raw — TUI only, by scope.
- The wordmark's a+i (the AI in jAIra) render in sunset pink (256-colour 211),
  `internal/tui/wordmarkstyle.go`; slicing is by runes, not bytes.
- Filter takes `key:value` (id/ticket, title, goal, context, assignee,
  lane/status, body); unknown keys fall back to full text; a known key with an
  empty field matches nothing.
- New tickets seed `## Progress` (where `jaira note` always wrote); the dead
  seeded `## Notes` is gone. `## Description` has been dead since 13.08.

CLI and plumbing
- `lanes use` resolves from the catalogue, so `--force` actually refreshes a
  stale project copy instead of copying it onto itself (`4f198f3`). This is
  the ONLY way to refresh a project copy of a built-in (TUI `R` compares
  against the user catalogue only, never embedded built-ins).
- `follows` is settable (`create --follows`, validates, exit 5 on dead ref),
  visible in show/TUI/--json; sign-off follow-ups carry the predecessor's
  commits into the context prose.
- Blocked lane requires a reason (`--reason`, or `blocked-by` counts);
  parking is exempt from the dependency check and the leaving lane's output
  contract; the reason renders as "waiting on" only while parked.
- shotgen finds the sign-off lane by probing instead of counting lanes.

Docs and assets
- README: the launcher screenshot IS the header (wordmark + pink AI + moon +
  example projects); all four screenshots regenerated from a scripted demo
  world; zero em dashes and zero Unicode ellipses (user: AI tells);
  GitHub repo description likewise de-dashed. Team flow documented in the
  Concurrency section.
- Screenshot recipe: build demo boards under a scratch `JAIRA_HOME` +
  `JAIRA_LANES_DIR`, then `go run ./scripts/shotgen <board> <view> <cols>
  <rows> | python3 scripts/termshot.py /dev/stdin docs/img/<view>.png --cols
  <cols>` (views: home, board, pipeline, signoff, edit).

## Open, needs the user

- **The req board** (separate repo, unreachable from here): run `jaira update`
  (agent-note block), `jaira lanes use in-progress --force` / `done --force` /
  `blocked --force` if it has project lane copies, and attribute the 07.08
  commits to DAHC06/YDACKQ/PJFVD1 via `jaira set <h> commits=<sha>,<sha>`.
- Cut v0.1.1? Everything user-facing since the release is unreleased.
- Follow-up chains: archive warning/cascade along `follows` — deferred until
  the links were visible; they are now.
- `jaira set follows=XXX` writes dead links (set validates nothing, decided);
  reopening is the user's call.
- Per-ticket `## Options` earning their keep — still open from 13.08.

## Known and deliberately not touched

- `gofmt` flags gate.go/tickets.go/lanes.go/view.go — pre-existing alignment
  groups, verified pre-existing at `d0a1b40`; session code is clean.
- `jaira archive` (no args) lists everything unfiltered; storage itself is
  cheap (measured: 1000 archived tickets ≈ 320 KB in git).
- Board render has a ~10-line floor; terminals ≤ 9 rows overflow slightly.
- No way to forget a board; `~/.jaira/state/` grows unbounded.
