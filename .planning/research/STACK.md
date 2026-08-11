# Stack Research

**Domain:** Single-binary terminal TUI kanban board over git-committed markdown tickets, driven by a CLI that AI agents call from bash
**Researched:** 2026-08-11
**Confidence:** HIGH on language choice, git integration, CLI framework, distribution. MEDIUM on TUI framework specifics. MEDIUM-LOW on YAML round-trip fidelity (explicitly flagged — see Q3).

## 1. Language: Go — not Rust

**Recommendation: Go 1.26.x. Confidence: HIGH.**

Both languages satisfy "single static binary, no runtime dependency, instant startup" — this alone doesn't decide it. Three concrete, verified differentiators do:

1. **Cross-compilation from the stated dev environment (Linux/WSL2) to macOS and Windows is a solved, zero-toolchain problem in Go and an unsolved one in Rust.** `GOOS=darwin GOARCH=arm64 go build` produces a working macOS binary from a Linux box with nothing extra installed, as long as the codebase is pure Go (`CGO_ENABLED=0`). Rust cross-compiling to macOS from Linux requires either a macOS SDK/osxcross setup or `cross` (Docker-based). Verified via GitHub API: `cross-rs/cross`'s last release is **v0.2.5, Feb 2023** — the standard Rust cross-compilation helper is effectively unmaintained. This directly matters because the team's primary (and by the constraints, only stated) dev environment is Linux/WSL2 and the binary must reach Windows and macOS teammates.
2. **YAML round-trip fidelity — a hard project constraint — is meaningfully better supported in the Go ecosystem** (see Q3 in detail). Rust's canonical YAML+serde library, `serde_yaml`, was archived/deprecated by its maintainer in March 2024 with no official successor; the closest lossless-editing crates (`yaml-edit`, `rust-yaml`) are small, low-adoption, early-stage projects. Go has an actively maintained library (`goccy/go-yaml`, v1.19.2, Jan 2026) with explicit AST/path-based node editing and a comment-preservation API built for exactly this use case.
3. **The Go TUI ecosystem (Charm's Bubble Tea/Lip Gloss/Bubbles stack) is "batteries-included" for this exact shape of app** — bordered multi-column layouts, list/table components, expandable panes — and is proven at scale in comparable real tools (`gh dash`, `k9s`-style dashboards, `soft-serve`, Charm's own `gum`, `glow`). Rust's `ratatui` is equally mature as a *rendering primitive* but leaves more of the higher-level widget/layout composition to the app, which is more work for a kanban-with-detail-panes UI specifically.

Binary size (Go ~8–15MB static, Rust ~2–5MB stripped) and raw runtime performance are not decisive: neither is a stated constraint, and both are irrelevant at "instant startup, single download" scale — a 12MB binary over a GitHub Release or Homebrew is still instant to fetch and instant to start (no GC pause is observable in a TUI that redraws on keypress/file-change, not in a hot loop).

**Do not treat this as "it depends."** Given this project's specific constraints (WSL2 dev box, cross-platform teammates, hand-edited/diffable YAML frontmatter as the API), Go is the lower-risk choice on every axis that actually matters here.

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

```bash
go mod init github.com/<org>/jaira

# Core
go get github.com/charmbracelet/bubbletea@v2.0.8
go get github.com/charmbracelet/lipgloss@v2.0.6
go get github.com/charmbracelet/bubbles@v2.1.1
go get github.com/spf13/cobra@v1.10.2
go get github.com/goccy/go-yaml@v1.19.2
go get github.com/fsnotify/fsnotify@v1.10.1
go get github.com/alecthomas/chroma/v2@v2.27.0

# Testing
go get github.com/charmbracelet/x/exp/teatest
```

## 2. TUI Framework Detail

**Bubble Tea (v2.0.8) + Lip Gloss (v2.0.6) + Bubbles (v2.1.1). Confidence: MEDIUM-HIGH** (versions verified via GitHub releases API; maturity claims are well-established community consensus, not independently benchmarked by this research).

| Concern | Assessment |
|---|---|
| Maturity/maintenance | Active — v2.0.8 released within the research window; Charm (the maintaining company) ships multiple production tools on this stack (`gh dash`, `glow`, `soft-serve`, `gum`, `crush`) |
| Multi-column kanban + keyboard nav | `bubbles/list` per lane column + a custom `Model` that tracks focused column/row; Lip Gloss `JoinHorizontal` composes N lane columns side by side with bordered boxes — this is a well-trodden pattern in the Charm ecosystem, not a novel one |
| Inline expandable detail pane | `bubbles/viewport` for a scrollable detail pane that toggles visibility in the `Update`/`View` cycle on a keypress — standard Bubble Tea pattern (collapse/expand is just conditional rendering in `View()`) |
| Live file-watch-driven re-render | Bubble Tea supports sending custom messages into the update loop from a goroutine via `Program.Send(msg)` (or a `tea.Cmd` that reads from a channel) — wire `fsnotify` events into this channel; this is the documented pattern for external-event-driven updates, not a workaround |
| Syntax-highlighted diff rendering | Chroma emits ANSI-escaped output directly (or styled tokens you feed into Lip Gloss); render the result inside a `viewport` — no special integration needed, both operate on plain strings with ANSI codes |

**Rust alternative for reference (ratatui v0.30.2):** Equally capable as a rendering primitive (immediate-mode redraw, `Layout`/`Block`/`List`/`Table` widgets exist and are actively maintained — verified current release June 2026). If Rust had been chosen, ratatui would be the correct pick over `cursive` (retained-mode, smaller community, less active) or `tui-rs` (unmaintained predecessor to ratatui — do not use). Not recommended here only because Go wins on the language-level decision above, not because ratatui is deficient.

## 3. Markdown + YAML Frontmatter Parsing — Round-Trip Fidelity

**Recommendation: `goccy/go-yaml` (v1.19.2), used at the AST level via its `parser`/`ast` packages with path-based node editing — not full struct marshal/unmarshal. Confidence: MEDIUM (capability verified as real and purpose-built; perfect byte-for-byte fidelity in every edge case is NOT guaranteed by any current Go or Rust library — flagged explicitly below).**

**What was checked and why it matters:**

- `serde_yaml` (Rust, the canonical YAML+serde crate) was **archived by its maintainer in March 2024** with an explicit "deprecated" note and no designated successor. This is disqualifying for a hard round-trip-fidelity requirement — building on a dead library for the project's core file format is a real risk, not a style preference.
- Rust's round-trip-focused alternatives (`yaml-edit`, built on `rowan` lossless syntax trees like rust-analyzer; `rust-yaml`, explicitly inspired by Python's `ruamel.yaml`) exist but are small, low-adoption projects — this is an immature corner of the Rust ecosystem as of this research.
- `gopkg.in/yaml.v3` (Go, the original niemeyer-maintained package used by most existing Go tools) was **marked unmaintained by its own author in April 2025**. It has since been forked under community governance as `go.yaml.in/yaml/v3` (v3.0.1), which is safe for simple marshal/unmarshal but — like its predecessor — does **not** guarantee comment/order preservation through a full unmarshal→marshal round trip (documented as a known, unresolved limitation in `go-yaml/yaml` issue #709).
- `goccy/go-yaml` is a separately-authored, actively maintained library (independent of the niemeyer/`go.yaml.in` lineage) that explicitly targets this problem: it provides `yaml.CommentToMap()`/`yaml.WithComment()` for comment round-tripping, preserves key order (order is positional in its Node/AST representation, not alphabetized), and exposes a `parser`/`ast` API for **path-based node editing** — i.e., parse the frontmatter to an AST, locate one field (e.g. `status:`) by YAMLPath, mutate only that node's scalar value, and re-print the AST. Untouched nodes keep their original tokens, so unrelated fields, comments, and blank lines are not touched.

**Practical recommendation for this project, given the hard constraint:**

1. Treat the ticket file as three zones: YAML frontmatter delimiters, the frontmatter body, and the markdown body. Read the whole file as bytes; only the frontmatter body goes through YAML processing.
2. For **reads** (rendering the board), unmarshal normally — full fidelity doesn't matter for a read-only view.
3. For **writes** (CLI mutating one or two fields, e.g. `status`), use `goccy/go-yaml`'s AST path-edit workflow rather than "unmarshal into a struct, mutate the struct, marshal the whole struct back out." Full struct remarshal is the actual danger here: it silently drops fields not in the Go struct (a real risk given the `external:` reserved block, which must survive round trips even though this tool doesn't understand its contents) and can reformat scalars (quoting style, block vs. flow) even when order is preserved.
4. **Be honest about the ceiling:** no current Go or Rust library — including `goccy/go-yaml`'s AST mode — is a general-purpose, provably lossless YAML editor for every legal YAML construct (anchors/aliases, exotic flow styles, tag directives). For a hand-authored, schema-constrained frontmatter format like this project's (flat keys, scalars, one array field `commits[]`, no anchors expected), AST-level single-field patching is realistic and low-risk. If a future ticket file is manually hand-edited into something exotic (e.g., a YAML anchor), a round-trip write could still reformat that specific construct. Mitigation: validate/normalize frontmatter shape on first read (reject or flag anchors/aliases) rather than relying on the editor to survive them silently.

This is the single highest-uncertainty area of this research (MEDIUM confidence, not HIGH) — flag it for a spike/prototype early in implementation: write a small AST-edit-and-reprint test against a handful of real-looking ticket files with comments and verify diff output is a single-line change before building the rest of the write path on top of it.

## 4. Git Integration

**Recommendation: shell out to the system `git` binary via `os/exec`. Do not use go-git, libgit2/git2-rs, or gix as the primary path. Confidence: HIGH.**

| Approach | Correctness | Binary size / static-binary impact | Assumes `git` on PATH? |
|---|---|---|---|
| **Shell out to `git`** (recommended) | Byte-identical to what a human running `git diff`/`git log` sees — same diff algorithm, same rename detection, same `.gitattributes` handling, no reimplementation drift | Zero added weight; pure Go `os/exec` | Yes — but this is free: if there's no `git` on the machine, there's no repo for this tool to operate on in the first place. This is not an added dependency risk. |
| go-git (pure Go, v5.19.2 stable / v6.0.0-alpha.5 in progress) | Reimplements git in Go; historically has had gaps vs. real git behavior (partial `.gitattributes`/hook support, diff output formatting that does not match `git diff` exactly) — a real risk when the whole point is showing a reviewer a diff they'll recognize | Pure Go, keeps static binary — the only real advantage | No `git` binary needed, but not needed here since git is already required |
| libgit2 (via CGO) / git2-rs (Rust) | Correct, battle-tested (used by GitHub itself historically) | **Disqualifying**: libgit2 is a C library. Binding to it requires CGO in Go or an FFI boundary in Rust, which either statically links a large C library (bloats and complicates the build matrix) or dynamically links it (reintroduces the exact "runtime dependency" this project is constitutionally against) | N/A |
| gix / gitoxide (pure Rust, actively released — verified v0.86.0-era crates, Aug 2026) | Fast-moving, well-regarded, but still explicitly documents partial feature coverage vs. canonical git for some plumbing | Pure Rust, keeps static binary | Moot — Rust isn't the chosen language |

Use `os/exec` for: `git log --format=...` (commit metadata for the `commits[]` field), `git show <sha>` (per-commit diff), and `git diff <sha1> <sha2>` (combined diff across a ticket's commits) — then pipe the raw text output through Chroma's `diff` lexer for the TUI's syntax-highlighted view. This mirrors what tools like `lazygit` and `gh` do and is the well-established pattern for Go CLIs that need git output a human will recognize.

## 5. File Watching

**Recommendation: `fsnotify` v1.10.1. Confidence: HIGH on the library, MEDIUM on WSL2-specific behavior (reasoned from well-documented WSL2 filesystem architecture, not independently reproduced in this research session).**

- `fsnotify` wraps the native OS mechanism: inotify (Linux, including WSL2's own Linux filesystem), kqueue (macOS/BSD), `ReadDirectoryChangesW` (Windows). Current release requires Go 1.23+, which is satisfied by the recommended Go 1.26.x.
- **WSL2 caveat (the stated primary dev environment):** inotify works correctly for files inside the WSL2 Linux filesystem (e.g. a repo cloned under `/home/...`). It is well-documented as **unreliable for files under `/mnt/c/...`** (Windows drives mounted into WSL2 via the 9P/DrvFs protocol) when those files are modified by a *Windows-side* process — DrvFs does not reliably surface inotify events for cross-boundary writes. Microsoft's own WSL guidance already recommends keeping git repos on the Linux filesystem for performance reasons unrelated to this tool; this reinforces the same recommendation for `jaira`'s file-watch feature to work reliably.
- Practical mitigation: document that `.jaira/tickets/` should live inside the WSL2 filesystem (not `/mnt/c`) for live-refresh to be reliable, and add a debounced (e.g. 500ms–1s coalesced) fallback poll of directory mtimes as a safety net for the cross-boundary case rather than depending on inotify alone.
- Multiple external writers (another Claude session, `git pull`, manual edit) means events can arrive in bursts (e.g. a `git pull` touching many files at once) — debounce/coalesce fsnotify events before triggering a re-render rather than re-rendering per individual event.

## 6. CLI Framework

**Recommendation: Cobra v1.10.2. Confidence: HIGH on the library, MEDIUM on the specific "agent-friendly" conventions (these are design choices this project makes, not something a library provides out of the box).**

Cobra is the standard for this shape of tool (`kubectl`, `gh`, `hugo`, `docker`, `helm`), gives built-in shell completion generation, and its command tree is easy for a "skill" document to describe programmatically (Cobra commands are self-describing via `Use`/`Short`/`Long`/flag metadata, which can be dumped to generate agent-facing docs). Lighter alternatives exist (`alecthomas/kong` — struct-tag driven, less boilerplate, smaller ecosystem) but Cobra's ubiquity means any agent (Claude Code or otherwise) has almost certainly seen extensive Cobra-based CLI examples in training data, which slightly helps agent tool-use reliability — a soft but real benefit for a tool explicitly designed to be "drivable by any bash-capable agent."

**Agent-friendly conventions to build on top of Cobra (project-level design, not a library):**

- **`--json` on every read command** (`jaira list --json`, `jaira show <id> --json`): emit a stable, versioned JSON schema instead of the human table/TUI-styled text. Never mix human-readable text and JSON on the same stream.
- **Stable, documented exit codes**: `0` success, `1` generic/unexpected error, `2` usage error (Cobra's own default for bad flags — keep it, don't override), and reserve a small set of project-specific codes (e.g. `3` = validation/schema error such as a missing `definition-of-done` on promotion, `4` = dependency-blocked). Document the table once and keep it stable across releases — agents branch on exit codes, so changing them is a breaking change.
- **Machine-readable errors on `--json`**: when `--json` is set, errors should also be JSON on stderr (e.g. `{"error":{"code":"blocked","message":"...","ticket":"..."}}`) rather than a free-text sentence, so an agent can parse failure reasons without regex.
- **Idempotent/safe re-invocation**: since multiple parallel Claude sessions may call the CLI concurrently against the same file store, commands should be safe to retry (e.g. moving a ticket to a lane it's already in is a no-op success, not an error) — this is a design implication of the "concurrency" constraint in PROJECT.md, not a library feature.

## 7. Testing

**TUI snapshot/golden testing. Confidence: MEDIUM (the recommended tool is explicitly experimental).**

- `charmbracelet/x/exp/teatest` provides `teatest.NewTestModel(...)`, `WaitFor`, and `RequireEqualOutput`, integrating with `charmbracelet/x/exp/golden` for `-update`-flag-driven golden file snapshots of a Bubble Tea program's rendered output. This is Charm's own answer to "how do I test a TUI" and is what their example repos use.
- **Caveat, stated plainly:** this package lives in the `x/exp` namespace — Charm's own convention for "no API stability guarantee." Expect it to shift under Bubble Tea version bumps. Mitigation: wrap it behind a small internal test-helper package (one file) so a `teatest` API change is a one-place fix, not a scattered rewrite across every test file.
- Alternative if `teatest` churn becomes a problem: hand-roll golden testing by calling `Model.View() string` directly in a plain Go test, stripping or keeping ANSI codes as needed, and diffing against a committed `testdata/*.golden` file with a manual `-update` flag (`flag.Bool("update", false, ...)` + `os.WriteFile` when set). This is more code but zero dependency risk — reasonable to start here and adopt `teatest` only if its `WaitFor`/async-message-handling helpers are needed for the file-watch-driven re-render tests specifically.

**Git-dependent code testing:**

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

**If a future v2 needs a persistent index/cache (e.g. for fast search across hundreds of tickets):**
- Use a pure-Go embedded store (`modernc.org/sqlite`, or even simpler, an in-memory index rebuilt from the markdown files on startup)
- Because any CGO-based store (real `sqlite3`, `libgit2`) reintroduces the exact runtime-dependency problem this stack is built to avoid

**If the `external:` Jira/YouTrack adapter (deferred to v2) is eventually built:**
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

---
*Stack research for: single-binary terminal TUI kanban board over git-committed markdown tickets, CLI-driven by AI agents*
*Researched: 2026-08-11*
