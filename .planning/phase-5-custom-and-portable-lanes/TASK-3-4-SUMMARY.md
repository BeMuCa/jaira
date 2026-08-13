# Task 3 & 4: `creator:` provenance, a legible `jaira lanes`, and private project lanes through `jaira share`

**Commits:**
- `abbe883` — feat: give every lane a creator and make jaira lanes legible (Task 3)
- `2b34fa7` — feat: keep a project's lane files private through jaira share (Task 4)

## Task 3 — what changed

- `Lane.Creator` is parsed from `creator:`. Defaults to `jaira` when the lane
  is built-in and the field is absent (no edits to the nine built-in files, as
  the plan directed). A custom lane with no `creator:` reports empty.
- `Lane.Overrides` is new: set by `Load` (not `parse`) to the id of the
  built-in a custom lane displaced, empty otherwise. Needed for `lanes show
  --json`'s `overrides` field.
- `lane.ProjectLanesActive(root)` is a new exported helper extracting the
  D-03 authority check (`ProjectLanesDir(root)` exists and holds ≥1 `.md`
  file) that `Load` already computed inline. `Load` now calls it too, so
  there is one place this boolean is computed rather than two that could
  drift apart — used by both `Load` and the new `lanes path` command.
- `jaira lanes` is now a parent command (`cmd.AddCommand`, following
  `projects.go`'s shape) with three subcommands, all working outside a board
  via `bestEffortRoot()` exactly as the bare command already did:
  - `lanes show <id>` — human output is a labelled block (ID, Name, Anchor,
    Precedence, Tier, Input, Output, Creator, Source, Overrides if any,
    Description if any) then the full prompt body. `--json` carries id, name,
    precedence, agentic, terminal, model_tier, builtin, input_requires,
    output_produces, source, prompt, creator, after, description, overrides.
    Unknown id → `fail(ExitUsage, "no_such_lane", ...)`, same pattern `move`
    already uses.
  - `lanes path` — prints the catalogue dir and (if inside a project) the
    project dir, marking whichever `ProjectLanesActive` says is in force.
    Outside a project it prints only the catalogue and says so in words
    ("not in a project directory"), matching the plan's explicit requirement
    that this command work outside a board. `--json`: `{catalogue, project,
    in_project, active}`.
  - `lanes template` — a `const` string, not generated: every field `parse`
    reads, one-line comments, explicitly states `model-tier` is a local alias
    and not a model name. Proven parseable by a test that writes the command's
    own stdout into a temp catalogue and asserts `lane.Load("")` picks up its
    `id: my-lane` with no warning naming the file.
- `lanes --json` (the existing list) now carries `prompt` and `creator` per
  lane, per the plan's explicit instruction.
- The human table (`ID NAME PREC AGENTIC TIER SOURCE`) is unchanged — no
  `CREATOR` column, per the plan's "if it does not fit, leave the table
  alone" — `lanes show` carries it instead.
- `internal/cli/lanes_test.go` is new, driving the real cobra tree via
  `newRoot("test")` + `SetArgs` (the pattern `checklist_test.go` /
  `update_test.go` use) with `JAIRA_LANES_DIR` set to an isolated temp dir per
  test. Covers: prompt+creator in the list JSON, `show`'s full human output,
  `show <unknown>`'s exit code and reason, `show --json`'s field set, `path`
  in and out of a project (both output modes, both active/inactive), the
  template round-trip, and `lanes typo` falling through to the parent's
  `noArgs()` usage error.
- `core/lane/lane_test.go` gained three tests: `creator:` parses onto
  `Lane.Creator`; every built-in defaults to `jaira`; a custom lane with no
  `creator:` stays empty.

### Deviations (Task 3)

**1. [Rule 3 — avoid duplicated logic] Extracted `ProjectLanesActive` from
`Load`'s inline check rather than duplicating the same `os.Stat` +
`filepath.Glob` test inside the new `lanes path` command.** Not explicitly
asked for, but the alternative was two copies of the same D-03 boundary
condition that could silently disagree after a future edit to one of them.
`Load`'s behaviour is unchanged — same test, same call sites, now shared.

No other deviations. The `<decided>` block, `lane-system-design.md`, and
Task 2's summary were all consistent with what Task 3 needed.

## Task 4 — what changed

- `core/board/gitignore.go`: `LanesIgnoreLine = "/.jaira/lanes/"`, plus
  `AddLanesIgnore(root)` / `RemoveLanesIgnore(root)` following `AddIgnore` /
  `RemoveIgnore`'s shape exactly (idempotent add, comment-plus-blank-line-aware
  remove). `Ignored(root)` is untouched — it only matches the whole-board
  line, never the lanes-only one, confirmed by test.
- `internal/cli/share.go`: after `RemoveIgnore` succeeds, calls
  `AddLanesIgnore` and reports it as `lanes_ignored` in `--json` and as a
  human line ("This project's lane files stay private: /.jaira/lanes/ is
  still ignored."). In `--undo`, calls `RemoveLanesIgnore` after `AddIgnore`
  and reports `lanes_ignore_removed`. The closing instruction keeps `git add
  .jaira .gitignore` as the actual command (gitignore already makes git skip
  `.jaira/lanes/` recursively under `git add .jaira`) but now says explicitly
  that the project's lanes, if any, stay out of it — per the plan's "name
  what is actually being published" instruction, this was a wording fix, not
  a command change, since the git command already worked correctly.
- `core/ticket/store.go`: `SharedSubdir = "shared"` and `SharedDir()`
  (matching `ArchiveDir()`'s shape) added; nothing creates the directory yet
  (Task 6's job, on first publish). The comment above `SessionsSubdir` no
  longer claims `.jaira/` holds only committed content — it now names which
  three subdirectories are committed, that `lanes/` is this machine's and
  gitignored even on a shared board, and that sessions/locks never enter the
  repo.
- `core/board/board_test.go`: four new tests — fresh-board add (whole-board
  line absent, lanes line present), idempotent add, undo-path removal, and
  the `Ignored()` non-fooling case (this is the one the plan calls out as
  mattering most: a false positive here would silently disable the merge
  driver on a shared board).

### Deviations (Task 4)

**1. [Adaptation — the plan's literal verify command would fetch the
published package, not local code] The second `<automated>` verify block in
the plan does `cd $(mktemp -d) && ... && go run github.com/BeMuCa/jaira/cmd/jaira
init`.** Run exactly as written, this `cd`s out of the module entirely and
then asks `go run` to resolve a full import path with no local module in
scope — `go` would either fail to find the module (as it did when I tried it
literally) or, in an environment with network access and no replace
directive, silently fetch and run the **public, previously-released**
`github.com/BeMuCa/jaira` package instead of the code just written in this
session. Either outcome defeats the point of the check. I ran the equivalent
proof from inside the repository root using `go run github.com/BeMuCa/jaira/cmd/jaira
-C "$TMPD" ...` instead — same net effect (a fresh repo, `init`, a hand-placed
`.jaira/lanes/x.md`, `share`, then `git add -A`), but resolves the module
locally. Verbatim output below. Flagging this so nobody re-runs the plan's
exact string expecting it to test local changes.

No other deviations.

## Verify output (verbatim)

Task 3:

```
$ export PATH=$PATH:$HOME/.local/go/bin && go test ./internal/cli/... ./core/lane/... 2>&1 | tail -20
ok  	github.com/BeMuCa/jaira/internal/cli	(cached)
ok  	github.com/BeMuCa/jaira/core/lane	(cached)

$ export PATH=$PATH:$HOME/.local/go/bin && go run ./cmd/jaira lanes template | go run ./cmd/jaira lanes show review >/dev/null; go run ./cmd/jaira lanes --json | grep -q '"prompt"' && echo "prompt in json"
prompt in json
```

`jaira lanes show nope` / `jaira lanes typo` exit codes (compiled binary, not
`go run`, since `go run` does not propagate the child's exact exit code):

```
$ /tmp/jaira-bin lanes show nope; echo "exit=$?"
jaira: no lane "nope" is installed; available: backlog, brainstorm, todo, pre-process, in-progress, human, review, signoff, done, blocked
exit=2
$ /tmp/jaira-bin lanes typo; echo "exit=$?"
jaira: lanes takes no arguments, received 1

Usage: jaira lanes [flags]
exit=2
```

Zero-warnings proof, clean install outside this checkout:

```
$ JAIRA_LANES_DIR="$FAKEHOME/lanes" go run github.com/BeMuCa/jaira/cmd/jaira -C "$TMPD" lanes
ID             NAME             PREC   AGENTIC  TIER     SOURCE
backlog        Backlog          0      false    —        built-in
brainstorm     Brainstorm       5      true     strong   built-in
todo           Todo             20     false    —        built-in
pre-process    Pre-process      25     true     strong   built-in
in-progress    Implementing     30     true     cheap    built-in
human          HITL             40     false    —        built-in
review         Review           50     true     strong   built-in
signoff        Sign-off         55     false    —        built-in
done           Done             60     false    —        built-in
blocked        Blocked          10     false    —        built-in
(nothing on stderr)
```

Task 4:

```
$ export PATH=$PATH:$HOME/.local/go/bin && go test ./core/board/... ./core/ticket/... ./internal/cli/... 2>&1 | tail -20
ok  	github.com/BeMuCa/jaira/core/board	0.026s
ok  	github.com/BeMuCa/jaira/core/ticket	0.053s
ok  	github.com/BeMuCa/jaira/internal/cli	0.206s
```

Real-repo proof (see Deviation 1 for why this replaces the plan's literal
command):

```
$ go run github.com/BeMuCa/jaira/cmd/jaira -C "$TMPD" init >/dev/null 2>&1
$ mkdir -p "$TMPD/.jaira/lanes" && echo x > "$TMPD/.jaira/lanes/x.md"
$ go run github.com/BeMuCa/jaira/cmd/jaira -C "$TMPD" share
The board is now part of the repository.
Removed /.jaira/ from .gitignore.
This project's lane files stay private: /.jaira/lanes/ is still ignored.
Wrote .jaira/.gitattributes so git merges tickets field by field.
Registered the merge driver in .git/config for this clone.

Commit to publish 0 ticket(s) (this project's lanes, if any, stay out of it):
  git add .jaira .gitignore && git commit -m "share jaira board"

Teammates then clone, and jaira binds the merge driver on their first command.

$ (cd "$TMPD" && git add -A && git status --porcelain)
A  .gitignore
A  .jaira/.gitattributes
A  AGENTS.md
A  CLAUDE.md
$ (cd "$TMPD" && git status --porcelain | grep -qv "\.jaira/lanes" && echo "lanes not staged")
lanes not staged
```

With an actual ticket present, `.jaira/tickets/<file>.md` is staged
alongside the above and `.jaira/lanes/` still is not (verified separately,
same run pattern, ticket file appeared as `A  .jaira/tickets/...` with no
`lanes` entry anywhere in `git status --porcelain`).

`share --undo` round trip on the same repo — confirms no orphan lanes line:

```
$ go run github.com/BeMuCa/jaira/cmd/jaira -C "$TMPD" share --undo
The board is private again.
Added /.jaira/ to .gitignore.

Tickets already committed stay in history. To remove them:
  git rm -r --cached .jaira && git commit -m "make jaira board private"

$ cat "$TMPD/.gitignore"

# jaira board — private to this machine. Run 'jaira share' to publish it.
/.jaira/
```

Full suite after both commits:

```
$ export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./... -count=1
?   	github.com/BeMuCa/jaira/cmd/jaira	[no test files]
ok  	github.com/BeMuCa/jaira/core/board	0.026s
ok  	github.com/BeMuCa/jaira/core/gate	0.269s
?   	github.com/BeMuCa/jaira/core/gitrepo	[no test files]
ok  	github.com/BeMuCa/jaira/core/lane	0.270s
ok  	github.com/BeMuCa/jaira/core/merge	0.153s
ok  	github.com/BeMuCa/jaira/core/project	0.045s
ok  	github.com/BeMuCa/jaira/core/release	0.011s
?   	github.com/BeMuCa/jaira/core/session	[no test files]
ok  	github.com/BeMuCa/jaira/core/ticket	0.053s
ok  	github.com/BeMuCa/jaira/core/validate	0.126s
ok  	github.com/BeMuCa/jaira/internal/cli	0.206s
ok  	github.com/BeMuCa/jaira/internal/tui	1.797s
?   	github.com/BeMuCa/jaira/scripts/iconpreview	[no test files]
?   	github.com/BeMuCa/jaira/scripts/shotgen	[no test files]
```

## What the next task's author needs to know

- **Task 5 is next in sequence but is explicitly held** per the environment
  note in this run: it depends on nothing this run touched, but the
  `do_not_touch` block covering `precedence`/ordering still applies —
  nothing here changed ordering, `precedence` values, or `order()`.
- `lane.ProjectLanesActive(root)` is now the one place D-03's "is the project
  directory authoritative" test lives. If Task 8's `Materialise` or Task 6's
  `Export` need to know whether writing into `ProjectLanesDir` will flip a
  project from catalogue-mode to project-mode, this is the function to call
  rather than re-deriving it.
- `lane.Overrides` exists on every `*Lane` now (empty string when not
  overriding). Task 6's lane settings screen ("marks overrides") can read
  this directly instead of comparing against `Builtins()` itself.
- `ticket.Store.SharedDir()` exists and returns `<root>/.jaira/shared`, but
  **nothing creates it yet** — Task 6's `Publish` is expected to
  `MkdirAll` it on first use, per the plan.
- The D-02 drift check still does not exist anywhere (confirmed unchanged
  from Task 2's summary) — still entirely Task 6's responsibility, in the
  lane settings screen.
- `board.LanesIgnoreLine` / `AddLanesIgnore` / `RemoveLanesIgnore` are
  exported and idempotent; nothing else in this task's scope calls them
  besides `jaira share`.

## Known stubs

None — Task 3 only added CLI surface backed by real loader data; Task 4 only
touched gitignore plumbing and a doc comment.

## Threat Flags

None. Task 3 only exposes read-only lane data already loaded and warned about
by the existing loader (no new trust boundary). Task 4 only changes what a
`.gitignore` file excludes; it does not change what is readable, writable, or
executable, and does not touch the merge driver's `Ignored()` semantics.

## Self-Check

```
$ [ -f core/lane/lane.go ] && echo FOUND || echo MISSING
FOUND
$ [ -f internal/cli/lanes_test.go ] && echo FOUND || echo MISSING
FOUND
$ [ -f core/board/gitignore.go ] && echo FOUND || echo MISSING
FOUND
$ [ -f core/ticket/store.go ] && echo FOUND || echo MISSING
FOUND
$ git log --oneline --all | grep -q abbe883 && echo FOUND || echo MISSING
FOUND
$ git log --oneline --all | grep -q 2b34fa7 && echo FOUND || echo MISSING
FOUND
```

## Self-Check: PASSED
