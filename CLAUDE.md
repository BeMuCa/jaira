<!-- GSD:project-start source:PROJECT.md -->

## Project

**jAIra**

A lightweight, git-synced kanban board for the tasks Claude Code generates while working in a repo. Tickets are markdown files committed inside the repo, so the board travels with the code — a teammate clones, runs `jaira`, and sees the same state with no server, no accounts, and no setup. Each lane is a configurable agent step (prompt + model tier + input/output contract), which turns the board from a tracker into a visible, resumable agent pipeline.

**Core Value:** **You never lose track of what a task was for or where it stands** — across session boundaries, agent runs, and teammates.

Everything else is negotiable. If the board can't answer "what was this ticket supposed to do, why, and what state is it in?" months later, it has failed.

### Constraints

- **Tech stack**: Single static binary, Go or Rust — no runtime dependency, because every teammate must be able to install and run it trivially
- **Storage**: Markdown with frontmatter, inside the repo — the file format *is* the API, and must stay hand-editable and diff-readable
- **Sync**: Git only — no server, so any feature requiring central coordination is off the table by construction
- **Performance**: Instant startup — the board is glanced at constantly, and a slow board will not be opened
- **Compatibility**: CLI must be usable by any bash-capable agent, not only Claude Code
- **Concurrency**: Multiple parallel Claude sessions and multiple teammates write the same store, so the format must minimize and survive conflicts
- **Scope discipline**: Every feature is measured against "is this smaller than paca?" — the project fails by growing

<!-- GSD:project-end -->

<!-- GSD:stack-start source:research/STACK.md -->

## Technology Stack

## 1. Language: Go — not Rust

## Recommended Stack

### Core Technologies

| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| Go | 1.26.5 (current stable, verified via go.dev/dl) | Language/toolchain | Native cross-compilation, static binaries by default, fast startup, mature stdlib `os/exec` for git shell-out |
| Bubble Tea | v2.0.8 (verified via GitHub releases API) | TUI runtime (Elm-architecture: Model/Update/View) | De facto standard Go TUI framework; v2 is GA (not beta) as of this research; used by Charm's own production tools |
| Lip Gloss | v2.0.6 | Terminal styling (borders, columns, colors, layout) | Pairs with Bubble Tea; gives you the bordered multi-column kanban look with minimal code; handles ANSI width/wrapping correctly (grapheme-aware) |
| Bubbles | v2.1.1 | Pre-built TUI components (list, table, viewport, textinput, textarea) | `list`/`table` components map directly onto "lane of cards" and "card detail pane"; avoids hand-rolling scroll/selection logic |
| Cobra | v1.10.2 | CLI command/flag framework for the `jaira` binary | Industry-standard for exactly this shape of tool (`kubectl`, `gh`, `hugo`, `docker` all use it); built-in shell completion generation and command tree, which the "skill that teaches Claude the CLI" can introspect |
| goccy/go-yaml | v1.19.2 | Frontmatter parsing with AST-level, order/comment-aware editing | Only actively maintained Go (or Rust) YAML library with an explicit workflow for reversible transforms; see Q3 |

### Supporting Libraries

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| fsnotify | v1.10.1 | Cross-platform file watching (inotify/kqueue/ReadDirectoryChangesW) | Watching `.jaira/tickets/` for external edits to trigger board re-render |
| alecthomas/chroma | v2.27.0 (stable; v3 is alpha as of this research — stay on v2) | Syntax highlighting for the git diff view rendered inside the TUI | Chroma ports Pygments' lexer set including a `diff` lexer; pipe `git show`/`git diff` output through it, then through Lip Gloss for terminal color codes |
| glamour | v2.0.1 | Markdown rendering (optional) | Only if ticket bodies/detail panes should render markdown (headers, lists) rather than raw text — evaluate against complexity budget before adding |
| go.yaml.in/yaml/v3 | v3.0.1 (community-governed fork) | Fallback/simple-case YAML (struct marshal only, no round-trip needs) | Use only for structures where round-trip fidelity does not matter (e.g., internal config, not tickets) — see "What NOT to Use" for why the *original* `gopkg.in/yaml.v3` is off the table |
| goreleaser | latest (actively released; nightly builds as of this research) | Build/release automation: cross-platform matrix, GitHub Releases, Homebrew tap formula | The standard distribution pipeline for Go CLIs; one YAML config produces linux/darwin/windows × amd64/arm64 binaries, checksums, and a Homebrew formula in one CI step |
| charmbracelet/x/exp/teatest | experimental (`x/exp` namespace — no API stability guarantee) | Golden/snapshot testing of Bubble Tea `View()` output | See Q7 — use with awareness it can break on Bubble Tea upgrades |

### Development Tools

| Tool | Purpose | Notes |
|------|---------|-------|
| goreleaser | Cross-compilation + packaging + release | Set `CGO_ENABLED=0` in the build matrix to guarantee static binaries on every target |
| golangci-lint | Static analysis | Standard for Go CLI projects; catches unchecked errors, which matter a lot around file I/O on tickets |
| GitHub Actions | CI: test matrix + goreleaser release on tag | Run tests on linux/macOS/windows runners; WSL2-specific file-watch behavior should additionally be manually smoke-tested (CI runners are not WSL2) |

## Installation

# Core

# Testing

## 2. TUI Framework Detail

| Concern | Assessment |
|---|---|
| Maturity/maintenance | Active — v2.0.8 released within the research window; Charm (the maintaining company) ships multiple production tools on this stack (`gh dash`, `glow`, `soft-serve`, `gum`, `crush`) |
| Multi-column kanban + keyboard nav | `bubbles/list` per lane column + a custom `Model` that tracks focused column/row; Lip Gloss `JoinHorizontal` composes N lane columns side by side with bordered boxes — this is a well-trodden pattern in the Charm ecosystem, not a novel one |
| Inline expandable detail pane | `bubbles/viewport` for a scrollable detail pane that toggles visibility in the `Update`/`View` cycle on a keypress — standard Bubble Tea pattern (collapse/expand is just conditional rendering in `View()`) |
| Live file-watch-driven re-render | Bubble Tea supports sending custom messages into the update loop from a goroutine via `Program.Send(msg)` (or a `tea.Cmd` that reads from a channel) — wire `fsnotify` events into this channel; this is the documented pattern for external-event-driven updates, not a workaround |
| Syntax-highlighted diff rendering | Chroma emits ANSI-escaped output directly (or styled tokens you feed into Lip Gloss); render the result inside a `viewport` — no special integration needed, both operate on plain strings with ANSI codes |

## 3. Markdown + YAML Frontmatter Parsing — Round-Trip Fidelity

- `serde_yaml` (Rust, the canonical YAML+serde crate) was **archived by its maintainer in March 2024** with an explicit "deprecated" note and no designated successor. This is disqualifying for a hard round-trip-fidelity requirement — building on a dead library for the project's core file format is a real risk, not a style preference.
- Rust's round-trip-focused alternatives (`yaml-edit`, built on `rowan` lossless syntax trees like rust-analyzer; `rust-yaml`, explicitly inspired by Python's `ruamel.yaml`) exist but are small, low-adoption projects — this is an immature corner of the Rust ecosystem as of this research.
- `gopkg.in/yaml.v3` (Go, the original niemeyer-maintained package used by most existing Go tools) was **marked unmaintained by its own author in April 2025**. It has since been forked under community governance as `go.yaml.in/yaml/v3` (v3.0.1), which is safe for simple marshal/unmarshal but — like its predecessor — does **not** guarantee comment/order preservation through a full unmarshal→marshal round trip (documented as a known, unresolved limitation in `go-yaml/yaml` issue #709).
- `goccy/go-yaml` is a separately-authored, actively maintained library (independent of the niemeyer/`go.yaml.in` lineage) that explicitly targets this problem: it provides `yaml.CommentToMap()`/`yaml.WithComment()` for comment round-tripping, preserves key order (order is positional in its Node/AST representation, not alphabetized), and exposes a `parser`/`ast` API for **path-based node editing** — i.e., parse the frontmatter to an AST, locate one field (e.g. `status:`) by YAMLPath, mutate only that node's scalar value, and re-print the AST. Untouched nodes keep their original tokens, so unrelated fields, comments, and blank lines are not touched.

## 4. Git Integration

| Approach | Correctness | Binary size / static-binary impact | Assumes `git` on PATH? |
|---|---|---|---|
| **Shell out to `git`** (recommended) | Byte-identical to what a human running `git diff`/`git log` sees — same diff algorithm, same rename detection, same `.gitattributes` handling, no reimplementation drift | Zero added weight; pure Go `os/exec` | Yes — but this is free: if there's no `git` on the machine, there's no repo for this tool to operate on in the first place. This is not an added dependency risk. |
| go-git (pure Go, v5.19.2 stable / v6.0.0-alpha.5 in progress) | Reimplements git in Go; historically has had gaps vs. real git behavior (partial `.gitattributes`/hook support, diff output formatting that does not match `git diff` exactly) — a real risk when the whole point is showing a reviewer a diff they'll recognize | Pure Go, keeps static binary — the only real advantage | No `git` binary needed, but not needed here since git is already required |
| libgit2 (via CGO) / git2-rs (Rust) | Correct, battle-tested (used by GitHub itself historically) | **Disqualifying**: libgit2 is a C library. Binding to it requires CGO in Go or an FFI boundary in Rust, which either statically links a large C library (bloats and complicates the build matrix) or dynamically links it (reintroduces the exact "runtime dependency" this project is constitutionally against) | N/A |
| gix / gitoxide (pure Rust, actively released — verified v0.86.0-era crates, Aug 2026) | Fast-moving, well-regarded, but still explicitly documents partial feature coverage vs. canonical git for some plumbing | Pure Rust, keeps static binary | Moot — Rust isn't the chosen language |

## 5. File Watching

- `fsnotify` wraps the native OS mechanism: inotify (Linux, including WSL2's own Linux filesystem), kqueue (macOS/BSD), `ReadDirectoryChangesW` (Windows). Current release requires Go 1.23+, which is satisfied by the recommended Go 1.26.x.
- **WSL2 caveat (the stated primary dev environment):** inotify works correctly for files inside the WSL2 Linux filesystem (e.g. a repo cloned under `/home/...`). It is well-documented as **unreliable for files under `/mnt/c/...`** (Windows drives mounted into WSL2 via the 9P/DrvFs protocol) when those files are modified by a *Windows-side* process — DrvFs does not reliably surface inotify events for cross-boundary writes. Microsoft's own WSL guidance already recommends keeping git repos on the Linux filesystem for performance reasons unrelated to this tool; this reinforces the same recommendation for `jaira`'s file-watch feature to work reliably.
- Practical mitigation: document that `.jaira/tickets/` should live inside the WSL2 filesystem (not `/mnt/c`) for live-refresh to be reliable, and add a debounced (e.g. 500ms–1s coalesced) fallback poll of directory mtimes as a safety net for the cross-boundary case rather than depending on inotify alone.
- Multiple external writers (another Claude session, `git pull`, manual edit) means events can arrive in bursts (e.g. a `git pull` touching many files at once) — debounce/coalesce fsnotify events before triggering a re-render rather than re-rendering per individual event.

## 6. CLI Framework

- **`--json` on every read command** (`jaira list --json`, `jaira show <id> --json`): emit a stable, versioned JSON schema instead of the human table/TUI-styled text. Never mix human-readable text and JSON on the same stream.
- **Stable, documented exit codes**: `0` success, `1` generic/unexpected error, `2` usage error (Cobra's own default for bad flags — keep it, don't override), and reserve a small set of project-specific codes (e.g. `3` = validation/schema error such as a missing `definition-of-done` on promotion, `4` = dependency-blocked). Document the table once and keep it stable across releases — agents branch on exit codes, so changing them is a breaking change.
- **Machine-readable errors on `--json`**: when `--json` is set, errors should also be JSON on stderr (e.g. `{"error":{"code":"blocked","message":"...","ticket":"..."}}`) rather than a free-text sentence, so an agent can parse failure reasons without regex.
- **Idempotent/safe re-invocation**: since multiple parallel Claude sessions may call the CLI concurrently against the same file store, commands should be safe to retry (e.g. moving a ticket to a lane it's already in is a no-op success, not an error) — this is a design implication of the "concurrency" constraint in PROJECT.md, not a library feature.

## 7. Testing

- `charmbracelet/x/exp/teatest` provides `teatest.NewTestModel(...)`, `WaitFor`, and `RequireEqualOutput`, integrating with `charmbracelet/x/exp/golden` for `-update`-flag-driven golden file snapshots of a Bubble Tea program's rendered output. This is Charm's own answer to "how do I test a TUI" and is what their example repos use.
- **Caveat, stated plainly:** this package lives in the `x/exp` namespace — Charm's own convention for "no API stability guarantee." Expect it to shift under Bubble Tea version bumps. Mitigation: wrap it behind a small internal test-helper package (one file) so a `teatest` API change is a one-place fix, not a scattered rewrite across every test file.
- Alternative if `teatest` churn becomes a problem: hand-roll golden testing by calling `Model.View() string` directly in a plain Go test, stripping or keeping ANSI codes as needed, and diffing against a committed `testdata/*.golden` file with a manual `-update` flag (`flag.Bool("update", false, ...)` + `os.WriteFile` when set). This is more code but zero dependency risk — reasonable to start here and adopt `teatest` only if its `WaitFor`/async-message-handling helpers are needed for the file-watch-driven re-render tests specifically.
- Do not commit nested `.git` directories as test fixtures (this nests a git repo inside the project's own repo, which git itself handles awkwardly, and most tooling/CI treats nested `.git` as a submodule boundary). Instead, **build fixture repos at test runtime**: `t.TempDir()` + `os/exec` calls to `git init`, `git config user.email/user.name` (required in CI, which has no global git config), `git add`, `git commit` to construct a realistic repo with real commit SHAs, then run the code under test against that path. This is slower than an in-memory fake but is the only approach that produces byte-identical output to what `git show`/`git diff` produce in production, which matters directly given the recommendation in Q4 to shell out to real `git`.
- Skip these tests (or mark them build-tag-gated) in any environment without `git` on PATH, though in practice every dev/CI environment for this project will have it.

## 8. What NOT to Use

| Avoid | Why | Use Instead |
|-------|-----|-------------|
| `serde_yaml` (Rust) | Archived/deprecated by its own maintainer, March 2024; no official successor exists for new projects | N/A — moot once Go is chosen; if evaluating Rust independently, use `saphyr`/`yaml-rust2` for parsing but accept no comment/round-trip support |
| `gopkg.in/yaml.v3` (original niemeyer package) | Marked "unmaintained" by its author, April 2025; superseded | `go.yaml.in/yaml/v3` (community fork, safe for simple marshal/unmarshal) or `goccy/go-yaml` (for round-trip editing) |
| Full struct marshal/unmarshal as the *write* strategy for ticket frontmatter | Silently reorders/reformats fields and can drop unrecognized keys (a real risk for the reserved, currently-unused `external:` block) — violates the explicit "must not reorder or reformat unrelated fields" constraint | AST/path-based single-field patching via `goccy/go-yaml`'s `parser`/`ast` packages (see Q3) |
| `libgit2`/`git2-rs` (CGO/FFI-bound git library) | Requires linking a C library — either bundled statically (bloats and complicates the cross-compile matrix) or dynamically (reintroduces the "runtime dependency" the project explicitly rejects) | Shell out to the system `git` binary via `os/exec` |
| `go-git` as the primary diff/log source for the human-facing diff viewer | Reimplements git in pure Go; output formatting and edge-case behavior (renames, `.gitattributes`) can diverge from what `git diff` actually shows, which matters because the reviewer's mental model is "the real git diff" | Shell out to system `git` |
| `cross-rs/cross` (Rust cross-compilation helper) | Last released Feb 2023 — effectively unmaintained; Rust cross-compilation to macOS from a Linux dev box remains genuinely painful without it | N/A — moot once Go is chosen |
| Any CGO dependency anywhere in the binary (e.g. `mattn/go-sqlite3` if a future feature wants embedded SQLite) | CGO defeats trivial cross-compilation and can reintroduce dynamic linking against a C runtime, undermining "single static binary, no runtime dependency" | If persistence beyond markdown files is ever needed, use a pure-Go implementation (e.g. `modernc.org/sqlite`) — but note the project's own storage design (markdown files, no DB) makes this moot for now |
| Building the TUI in Node/Ink despite Node 24 being available in the dev environment | Node is a runtime dependency by definition — directly violates the "no runtime dependency" constraint and the explicit "Go or Rust" constraint in PROJECT.md; tempting only because Node tooling is already on the dev machine | Go + Bubble Tea, per the recommendation above |
| `charmbracelet/x/exp/teatest` treated as a stable, permanent API | It's explicitly experimental (`x/exp` namespace) | Wrap it behind a thin internal helper so upgrades are contained, or fall back to hand-rolled golden files (see Q7) |
| Committing nested `.git` directories as test fixtures | Git and most CI tooling treat a nested `.git` as ambiguous/submodule-like; fragile and confusing | Generate fixture repos at test runtime via `os/exec` + `t.TempDir()` |

## Stack Patterns by Variant

- Use a pure-Go embedded store (`modernc.org/sqlite`, or even simpler, an in-memory index rebuilt from the markdown files on startup)
- Because any CGO-based store (real `sqlite3`, `libgit2`) reintroduces the exact runtime-dependency problem this stack is built to avoid
- Keep it entirely additive to the frontmatter schema — never make the core CLI's read/write path aware of `external:`'s internal shape
- Because the AST-level editing strategy in Q3 only stays safe if unknown blocks are never touched by the core tool's serializer

## Version Compatibility

| Package A | Compatible With | Notes |
|-----------|-----------------|-------|
| Go 1.26.5 | fsnotify v1.10.1 | fsnotify v1.10.0+ requires Go 1.23+; satisfied |
| Bubble Tea v2.0.8 | Lip Gloss v2.0.6, Bubbles v2.1.1 | These three are versioned and released together by Charm and are the intended combination — do not mix Bubble Tea v2 with Bubbles v1 (different API surface, v1 targeted Bubble Tea v1) |
| Cobra v1.10.2 | go.yaml.in/yaml/v3 | Cobra itself migrated its internal YAML dependency from `gopkg.in/yaml.v3` to `go.yaml.in/yaml/v3` in this version specifically because of the former's unmaintained status — corroborates the Q3 finding independently |
| goccy/go-yaml v1.19.2 | Go 1.26.5 | No known compatibility issues found; independent of the yaml.v3/go.yaml.in lineage |

## Sources

- GitHub Releases API (`api.github.com/repos/.../releases`) — used to verify current version numbers directly rather than trusting summarized web pages, for: bubbletea, lipgloss, bubbles, cobra, fsnotify, ratatui, go-git, goccy/go-yaml, chroma, glamour, cross-rs/cross, gitoxide. HIGH confidence (primary source, machine-readable).
- crates.io API (`crates.io/api/v1/crates/...`) — verified gix, insta, ratatui, clap, notify crate current versions. HIGH confidence.
- go.dev/dl (JSON endpoint) — verified current stable Go release (1.26.5). HIGH confidence.
- [go-yaml/yaml issue #709](https://github.com/go-yaml/yaml/issues/709) — documents the known limitation that no existing Go YAML library fully round-trips comment positions. MEDIUM-HIGH confidence (primary source, direct maintainer discussion).
- [zerokspot.com — "A maintained YAML library for Go again!"](https://zerokspot.com/weblog/2026/02/07/maintained-golang-yaml-library/) and [go-yaml/yaml.in fork announcement search results] — corroborate the April 2025 unmaintained status of `gopkg.in/yaml.v3` and the `go.yaml.in` fork. MEDIUM confidence (community source, cross-checked against Cobra's own dependency migration as independent corroboration).
- [github.com/goccy/go-yaml](https://github.com/goccy/go-yaml) and its issue tracker (#297, #651, #225) — describes AST-based reversible transformation and comment-map API. MEDIUM confidence (README + issues describe the capability; exact byte-for-byte fidelity was not independently hand-tested in this research session — flagged as a spike recommendation in Q3).
- [users.rust-lang.org — "Serde-yaml deprecation. alternatives?"](https://users.rust-lang.org/t/serde-yaml-deprecation-alternatives/108868) — community discussion confirming `serde_yaml` archival and the fragmented state of Rust YAML alternatives. MEDIUM confidence.
- [Charm blog — "Writing Bubble Tea Tests"](https://charm.land/blog/teatest/) and `pkg.go.dev/github.com/charmbracelet/x/exp/teatest` — describes teatest/golden testing workflow and its experimental status. MEDIUM confidence.
- General Go/Rust ecosystem knowledge (Cobra's ubiquity in `kubectl`/`gh`/`hugo`/`docker`; WSL2 DrvFs/inotify interaction; standard `os/exec` git-shelling pattern used by `lazygit`) — MEDIUM confidence, well-established community knowledge not independently re-verified line-by-line in this session but consistent across multiple independent tools/projects the researcher is aware of.
- Note: several early WebSearch results for "Go vs Rust 2026 benchmarks" (e.g. tech-insider.org, ztabs.co) had signs of low-quality/AI-generated SEO content (implausible specific claims like "12x benchmark gap," "$25K salary divide") and were **explicitly discarded** rather than cited — flagging this per the honesty requirement rather than silently omitting the search.

<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->

## Conventions

Conventions not yet established. Will populate as patterns emerge during development.
<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->

## Architecture

Architecture not yet mapped. Follow existing patterns found in the codebase.
<!-- GSD:architecture-end -->

<!-- GSD:skills-start source:skills/ -->

## Project Skills

No project skills found. Add skills to any of: `.claude/skills/`, `.agents/skills/`, `.cursor/skills/`, `.github/skills/`, or `.codex/skills/` with a `SKILL.md` index file.
<!-- GSD:skills-end -->

<!-- GSD:workflow-start source:GSD defaults -->

## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:

- `/gsd:quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd:debug` for investigation and bug fixing
- `/gsd:execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->

<!-- GSD:profile-start -->

## Developer Profile

> Profile not yet configured. Run `/gsd:profile-user` to generate your developer profile.
> This section is managed by `generate-claude-profile` -- do not edit manually.
<!-- GSD:profile-end -->

<!-- jaira:start -->
## Task tracking: jaira

This repository has a jaira board (`.jaira/`). Multi-step work is tracked
there as markdown tickets so it survives session boundaries.

Capturing and picking work:

- `jaira create <title> --goal <...> --context <...> --dod <...>` — one call files a
  complete ticket; without a goal, a definition of done, the context it came from
  and an assignee it cannot leave the backlog
- the context is the only record of why a ticket exists. Write it for someone who
  was not in this conversation and reads it weeks from now: what is wrong today,
  what triggered it, what is already known or ruled out. Write it as if that
  reader has mild ADHD and knows none of what you know — lead with what is wrong,
  short concrete lines with one point each, names and paths rather than
  adjectives, no jargon and no preamble. It may span several lines, but someone
  should be able to act after the first two. If acting on it would need a
  question answered first, it is not finished
- `jaira list --actionable --json` — everything that could be started right now
- `jaira next --json` — the single next actionable ticket
- `jaira tags` — the tags this board already uses, with how many open tickets
  carry each. **Read it before you tag anything** and reuse the name that is
  already there for that subject; never invent a synonym — "ui", "frontend" and
  "gui" on one board are three names for one thing and filter to nothing.
  `jaira tag <id> <name>...` adds tags, and `jaira create --tag <name>` sets them
  at capture. A name jaira has not seen is new and gets a colour in the
  hand-editable `.jaira/tags`; the board filter and `jaira list --tag <name>`
  read them back

Working a ticket:

- `jaira claim <id>` — take it first; other sessions read this board too
- `jaira show <id> --for-lane <lane> --json` — the lane's prompt, the bounded input,
  the model tier, and the outputs the lane expects back
- `jaira dod <id> <n> --doing|--done` — mark checklist items as you go
- `jaira note <id> <text>` — at every pause, write down what the repository does
  not already say: dead ends, why this and not that, what you had to find out.
  Not what the checklist and git already record. A killed session never gets a
  turn to write anything down, so do not save it for the end
- `jaira move <id> --to <lane> --what <...> --why <...> --resolves <...>` — finish
  the step. jaira works out the commit list itself from git history — the union
  of the ticket file's own history and commits naming its id — so nothing needs
  to be typed here; it is written onto the ticket once, when the ticket leaves
  the board
- `jaira resume` — work left in progress, with everything recorded about it
- on a board that has not been shared yet (`jaira init` gitignores `.jaira/`
  until `jaira share`), the ticket file is untracked, so the only thing tying a
  commit to a ticket is its handle in the commit message. Name it there —
  `fix(A3K9QP): ...` — or the derived list stays empty and the move is refused
- **the ticket rides in the same commit as the code.** Move the ticket first,
  then `git add` the changed file under `.jaira/tickets/` alongside your source
  changes and commit them together. A reviewer then sees the change and what it
  was for in one place, instead of a diff whose ticket is still in whatever
  state the last commit left it. Same for a ticket you create and hand to
  someone else: commit it, or nobody but you knows it exists — and now this is
  also what makes the commit list derivable at all: that shared commit is how
  git ties the ticket to the change
- `jaira logbook <id>` — once a ticket reaches the terminal lane, stamps its
  commits and files it under `.jaira/logbook/<you>-<date>/`, taking it off the
  board. `jaira restore <file>` brings it back

`jaira <command> --help` for everything else.

Do not edit files under `.jaira/tickets/` directly; the CLI is the write path.
The human review lane cannot be left by an agent — a person accepts the work there.

A `jaira:local` marker — an HTML comment of that name — added by hand anywhere
inside this block makes everything between it and the end marker survive the
next regeneration. Nothing writes it for you, and there is none here until
somebody adds one; project-specific rules belong behind it rather than fighting
this note from outside the block.

## This board's lanes

Order: backlog → brainstorm → todo → pre-process → in-progress → critique → optimize → testing → human → review → signoff → done → blocked
Loop: critique sends work back to in-progress, and that repeats until critique has nothing left to say.
Loop: optimize sends work back to in-progress, and that repeats until optimize has nothing left to say.
Loop: testing sends work back to in-progress, and that repeats until testing has nothing left to say.

- `backlog` — no agent step; move through it
  Captured but not yet specified enough to work on.
- `brainstorm` — yours to work; tier strong; must produce goal
  Working out what the ticket should even be.
- `todo` — no agent step; move through it
  Specified and ready to be picked up.
- `pre-process` — yours to work; tier strong; must produce plan
  Working out how the change will be made.
- `in-progress` — yours to work; tier cheap; must produce outcome-what, outcome-why, outcome-resolves
  Carrying out the plan.
- `critique` — yours to work; tier strong; must produce review-summary
  Judges whether this is the right implementation, not whether it works.
- `optimize` — yours to work; tier strong; must produce review-gaps
  Removes what the change does not need — code that already exists elsewhere, code nobody calls, and code that carries its weight in nothing.
- `testing` — yours to work; tier cheap; must produce test-verdict
  Runs the change and checks it against the ticket - does the demanded thing exist, and does it work.
- `human` — **a person's, not yours** — you may move work in, never out
  Human in the loop.
- `review` — yours to work; tier strong; must produce review-summary, review-gaps, review-verdict, review-check
  A second model has judged the diff.
- `signoff` — **a person's, not yours** — you may move work in, never out
  Reviewed by a model, waiting for a person to accept it or send it back.
- `done` — no agent step; move through it; terminal
  Accepted.
- `blocked` — no agent step; move through it; parking: work returns to the lane it left
  Waiting on an external dependency.

Nothing moves a ticket for you. There is no daemon and no runner: a lane's
prompt runs because a session ran it. Drive the board from the session you are
in — `jaira next --per-lane --json` says which lanes have work waiting and
which of them you are allowed to work, `jaira show <id> --for-lane <lane> --json`
hands you that lane's prompt and its bounded input, and `jaira move` puts the
result in the next lane. Work one lane to empty before starting the next, or
the lane nobody drives is the one that fills up.

Told to start or work a ticket, drive it this way yourself — lane by lane,
loops included — until it sits in a human lane, then continue once the human
has answered. Told an agent should work it, hand it to a subagent that
babysits the ticket through the same route.
<!-- jaira:end -->
