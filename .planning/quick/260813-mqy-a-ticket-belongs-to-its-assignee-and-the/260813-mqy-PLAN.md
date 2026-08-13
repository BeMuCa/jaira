---
phase: quick/260813-mqy
plan: 01
type: execute
wave: 1
depends_on: []
files_modified:
  - core/identity/identity.go
  - core/identity/identity_test.go
  - core/gate/gate.go
  - core/gate/gate_test.go
  - internal/cli/root.go
  - internal/cli/flow.go
  - docs/AGENTS.md
autonomous: true
requirements: [QUICK-260813-mqy]

must_haves:
  truths:
    - "Writing to a ticket assigned to someone else is refused, and the refusal names the owner and the two ways forward"
    - "The human checkpoint lanes are exempt: reviewing and signing off someone else's work is the point of them"
    - "Taking a ticket over by setting its assignee is always allowed, so a ticket is never frozen by an absent owner"
    - "--force overrides the refusal and the override is recorded, exactly as it already is for gate refusals"
    - "A ticket with no assignee belongs to nobody and is writable by anyone"
    - "Nothing about the merge driver changes: ownership is a client-side guard rail, not a lock, because git merges files whether or not the CLI approved the write"
  artifacts:
    - path: "core/identity/identity.go"
      provides: "Who am I — moved out of internal/cli so core and the TUI can ask"
      exports: ["Current", "Slug"]
    - path: "core/gate/gate.go"
      provides: "the ownership check, alongside the existing gates"
  key_links:
    - from: "core/gate/gate.go"
      to: "core/identity.Current"
      via: "compare the actor against the ticket's assignee"
      pattern: "identity\\."
---

<objective>
Nothing enforces ownership today. `assignee` is a field the promotion gate
requires to be non-empty (`core/gate/gate.go:83`) and a filter for `list`
(`internal/cli/flow.go:254`); no write is ever refused because the ticket is
someone else's.

Purpose: two teammates working the same board should almost never edit the same
ticket, so the merge path stops being a thing anyone meets in practice. The CLI
refuses to write to a ticket that belongs to someone else and says who owns it.

Explicitly NOT the purpose: removing the merge driver. jaira has no server, and
the file format is meant to stay hand-editable, so ownership can only be a
refusal in this binary. Git will still merge two versions after an offline pull,
a rebase, or a hand edit. The merge rules stay exactly as they are.
</objective>

<context>
@CLAUDE.md
@core/gate/gate.go
@internal/cli/flow.go
@internal/cli/root.go

<interfaces>
`identity()` lives at internal/cli/root.go:278-292 today: `JAIRA_USER`, then
`git config user.name`, then `$USER`/`$USERNAME`/`$LOGNAME`, then a fallback.
`internal/cli` imports `internal/tui`, so anything both need must live in core.
Phase 5's plan also wants this function in a core package for the shared-lane
folder name — moving it here serves both, and that plan should be left to find
it already moved.

The `--force` mechanism already exists for gate refusals: `internal/cli/flow.go:22,
158, 186, 190, 204`. Reuse it; do not invent a second override flag.

Lane contracts already mark the human checkpoints: `RequiresQuestion` on `human`
and `RequiresHumanExit` on `signoff`. Use those flags rather than hardcoding lane
ids, the way the gate rework for `requires-specified` did.

Exit code 3 is "a gate refused the operation" and is already documented.

Go is not on the default PATH:
  export PATH=$PATH:$HOME/.local/go/bin
</interfaces>
</context>

<tasks>

<task type="auto">
  <name>Task 1: Move identity into core</name>
  <files>core/identity/identity.go, core/identity/identity_test.go, internal/cli/root.go</files>
  <action>
Create `core/identity` with `Current(dir string) string`, holding the body of
`identity()` from internal/cli/root.go:278-292 verbatim — same order of sources,
same fallback. Add `Slug(name string) string`, lowercasing and reducing to
`[a-z0-9-]`, collapsing runs and trimming leading/trailing separators, because
the shared-lane folder in phase 5 needs a name that is safe as a path component.

`internal/cli`'s `identity()` becomes a one-line call through to it, so ticket
attribution keeps behaving exactly as it does now. Do not change any caller.

Tests: `Current` honours `JAIRA_USER` above everything else; falls through when
it is empty. `Slug` on "BeMuCa" → "bemuca", on "Anna Müller" → something safe as
a path component with no separators at either end, and on a string of only
punctuation → a non-empty fallback rather than "", since an empty path component
would silently write to the parent directory.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./core/identity/... ./internal/cli/... -count=1 2>&1 | tail -5</automated>
  </verify>
  <done>core/identity holds Current and Slug with tests; internal/cli calls through; every existing test still passes.</done>
</task>

<task type="auto">
  <name>Task 2: A ticket belongs to its assignee</name>
  <files>core/gate/gate.go, core/gate/gate_test.go, internal/cli/flow.go</files>
  <action>
Add an ownership check to the gate, next to the existing ones. It refuses when
all of these hold:

- the ticket has a non-empty `assignee`
- the actor is not that assignee (compare case-insensitively, trimmed)
- the ticket is NOT currently in a lane whose contract marks it a human
  checkpoint (`RequiresQuestion` or `RequiresHumanExit`)
- the move is not a hand-over — a change of `assignee` is always allowed

Add the actor to whatever the gate already receives rather than reading the
environment from inside `core/gate`; the gate must stay a pure function of its
inputs, which is what makes it testable. Follow how `req.Interactive` is already
threaded through.

The message must name the owner and both ways forward, in the repo's plain
register, e.g.:

    ABC123 belongs to anna — ask her, take it over with
    'jaira set ABC123 assignee=<you>', or override with --force

Give it its own violation code alongside the existing ones so `--json` consumers
can branch on it. Exit code stays 3.

`--force` already overrides gate violations and records the override
(internal/cli/flow.go:158,190) — the ownership violation must flow through that
same path, with no second flag.

A ticket with an empty `assignee` belongs to nobody and is writable by anyone;
that is what makes a fresh scrap usable before it has been assigned.

Tests in core/gate/gate_test.go, matching the existing style:
- another owner, ordinary lane → refused, message names the owner
- another owner, but the ticket sits in a `RequiresHumanExit` lane → allowed
- another owner, but the ticket sits in a `RequiresQuestion` lane → allowed
- another owner, and the request sets a new assignee → allowed
- no assignee at all → allowed
- the actor IS the assignee, differing only in case → allowed
- the refusal carries the ownership code, not a generic one
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./... -count=1 2>&1 | grep -Ev "no test files" | tail -12</automated>
  </verify>
  <done>The gate refuses a write to someone else's ticket outside the human lanes, --force overrides it, hand-over is always allowed, and every case above is covered by a test.</done>
</task>

<task type="auto">
  <name>Task 3: Say that a ticket has an owner</name>
  <files>docs/AGENTS.md</files>
  <action>
docs/AGENTS.md tells whoever drives the board what it will refuse. Add a short
paragraph in the surrounding register: a ticket belongs to its assignee, another
person's ticket is refused with exit code 3, the human checkpoint lanes are
exempt because reviewing someone else's work is their purpose, taking a ticket
over is always allowed, and `--force` overrides and is recorded.

Say plainly that this is a guard rail and not a lock: the files are still plain
markdown in git, so a hand edit or an offline merge can still produce two
versions of a ticket, and the merge rules are what handle that. A reader who
believes ownership is enforced everywhere will be surprised exactly once, badly.

A few sentences. This file is a reference, not a tutorial.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && grep -ci "belongs to" docs/AGENTS.md</automated>
  </verify>
  <done>docs/AGENTS.md describes ownership, its exemptions, and its limits.</done>
</task>

</tasks>

<verification>
export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./... -count=1

Manual, against a real binary in a throwaway repo with JAIRA_HOME redirected:
create a ticket, set its assignee to someone else, then try to move it with
JAIRA_USER set to a different name. Expect exit 3 and a message naming the owner.
Repeat with --force and expect it to go through.
</verification>

<success_criteria>
- Writing to someone else's ticket is refused, naming the owner and the ways forward.
- Review and sign-off still work on other people's tickets.
- A ticket is never frozen: taking it over is always allowed.
- --force overrides, through the existing mechanism, and is recorded.
- The merge driver and its rules are untouched.
- `go build ./... && go vet ./... && go test ./... -count=1` passes.
</success_criteria>

<output>
Create `.planning/quick/260813-mqy-a-ticket-belongs-to-its-assignee-and-the/260813-mqy-SUMMARY.md` when done
</output>
