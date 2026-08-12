# Handoff — 2026-08-12 (evening)

State after the session that shipped jaira to GitHub. Written before a context
clear; this file is the memory.

## Where things are

- **Public: https://github.com/BeMuCa/jaira** — 53 commits, branch `master`,
  CI green on linux/macOS/windows at `4d3cace`.
- Module path is `github.com/BeMuCa/jaira`.
- Go is on PATH now (`~/.bashrc` and `~/.profile` both, because `.bashrc`
  returns early for non-interactive shells and a login shell reads `.profile`).
- `jaira` and `jaira-iconpreview` are built into `~/.local/bin`.
- The skill is installed globally at `~/.claude/skills/jaira/SKILL.md`, kept in
  sync with the repo copy by hand — **re-copy it when the repo copy changes.**

## Two bugs CI found that local testing could not

Both were invisible on Linux. This is the argument for keeping the CI matrix.

1. **macOS**: `/var` is a symlink to `/private/var`, so a board registered
   through one spelling and discovered through the other appeared twice on the
   launcher. `project.Canonical` is now the single reduction every caller uses.
2. **Windows**: `builtinFS.ReadFile(filepath.Join("builtin", n))` produced
   `builtin\00-backlog.md`, which does not exist in an embedded filesystem —
   those always use forward slashes. **No lanes loaded at all on Windows**, for
   the whole life of the project until now. Fixed with `path.Join`.

## What the pipeline looks like now

```
backlog → todo → pre-process → in-progress → human → review → signoff → done
                  (optional)   (Implementing) (HITL)  (agent)  (human)   + blocked
```

- **pre-process** is opt-in per ticket: a ticket's `## Options` checklist must
  tick `planning`. It produces the `## Plan` checklist and cannot be left
  without one. Skipping it is free.
- **review** is the model's step (writes `review-verdict`); **signoff** is the
  human checkpoint no agent may leave (`requires-human-exit`).
- `in-progress` and `human` keep their **ids** deliberately — those ids are
  written into existing tickets, and renaming would strand them read-only.

## Decisions not to re-open

- Tickets are markdown, not JSON. JSON storage would break the field-aware merge
  driver (a body becomes one escaped string) and make `git diff` unreadable.
  `--json` is the *machine interface*, not the storage format.
- `--signal` was removed. It accepted unchecked free text as evidence; an agent
  could close any ticket by typing a sentence.
- Progress notes live in the ticket body (`## Progress`), not a sidecar memory
  file — a separate file is one more thing to lose and would not travel in git.
- Checkbox position is load-bearing: an item only counts under its own heading.
  Anything writing items must use `ticket.AddItem`, never append to the file.
  This mistake was made three times while testing.
- The board is glanced at constantly, so bare `jaira` opening the launcher costs
  a keypress on the common path. `jaira board` is the documented direct route.

## Verified this session

120 tests green under `-race`, `gofmt`/`vet` clean, six cross-compile targets
with CGO disabled, and a real pseudo-terminal exercised: home screen, board,
detail pane, `e` field editing (write landed on disk, umlauts and newlines
intact), `E` `$EDITOR` round-trip, the sign-off view, `a` accept (a person could,
the CLI could not), `x` archive, the compact view and its serpentine layout.

## Not verified

- `go install github.com/BeMuCa/jaira/cmd/jaira@latest` — never run.
- How GitHub renders the README; the referenced files all exist remotely.
- The `GEMINI.md` / `.cursorrules` / Aider rows in `docs/AGENTS.md` are from
  knowledge, not testing.
- Windows beyond "CI passes" — nobody has run jaira there, and since no lanes
  loaded at all until today, other Windows issues may be hiding behind that.
- Whether the global skill is discovered by a session in another repository.
- The 44 requirementsgenie tickets are still unmigrated.
