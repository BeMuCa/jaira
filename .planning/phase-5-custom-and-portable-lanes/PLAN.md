---
phase: 5-custom-and-portable-lanes
plan: 01
type: execute
wave: 1
depends_on: []
mode: mvp
files_modified:
  - core/lane/lane.go
  - core/lane/lane_test.go
  - core/lane/share.go
  - core/lane/share_test.go
  - core/ticket/store.go
  - core/board/gitignore.go
  - core/board/board_test.go
  - core/identity/identity.go
  - core/identity/identity_test.go
  - internal/cli/root.go
  - internal/cli/tickets.go
  - internal/cli/lanes_test.go
  - internal/cli/share.go
  - internal/cli/checklist.go
  - internal/cli/resume.go
  - internal/cli/validate.go
  - internal/cli/mergedriver.go
  - internal/cli/json_test.go
  - internal/tui/model.go
  - internal/tui/home.go
  - internal/tui/lanes.go
  - internal/tui/lanes_test.go
  - internal/tui/view.go
  - core/gate/plan_test.go
  - core/validate/validate_test.go
  - core/merge/merge_test.go
  - internal/tui/signoff_test.go
  - internal/tui/render_test.go
  - .planning/REQUIREMENTS.md
autonomous: false
requirements: [LANE-02, LANE-03, LANE-04, LANE-05, LANE-06, LANE-07, LANE-08]
user_setup: []

must_haves:
  truths:
    - "A lane file dropped in .jaira/lanes/ appears on this project's board and nowhere else"
    - "A custom lane may override a built-in, and the override is always named in a warning"
    - "jaira lanes show <id> prints a lane's full prompt without a ticket in hand"
    - "jaira lanes path names the directory to write a lane into, and jaira lanes template prints one to write against"
    - "After jaira share, .jaira/tickets/ and .jaira/shared/ are committable and .jaira/lanes/ is still ignored"
    - "Every lane reports a creator, defaulting to jaira for built-ins"
    - "A lane can be published to .jaira/shared/<user>/ from the TUI, and a teammate can adopt it into their catalogue"
    - "A ticket in an unrecognised lane still renders, and cannot be advanced"
  artifacts:
    - path: "core/lane/lane_test.go"
      provides: "the lane package's first tests: loading, precedence between the three sources, override warnings, ordering"
      min_lines: 120
    - path: "core/lane/share.go"
      provides: "Export/Publish/Adopt — copying a lane file between catalogue, project and shared, as pure testable functions"
    - path: "core/identity/identity.go"
      provides: "who <user> is, in a package both internal/cli and internal/tui can import"
    - path: "internal/tui/lanes.go"
      provides: "the lane settings screen: read a lane, see its prompt, use it here, publish it, adopt a teammate's"
    - path: "internal/cli/lanes_test.go"
      provides: "jaira lanes / lanes show / lanes path / lanes template output, human and --json"
  key_links:
    - from: "core/lane.Load"
      to: "<root>/.jaira/lanes"
      via: "project lane directory, authoritative when present"
      pattern: "ProjectLanesDir"
    - from: "internal/cli/share.go"
      to: "core/board.AddLanesIgnore"
      via: "share writes /.jaira/lanes/ so publishing the board does not publish the prompts"
      pattern: "LanesIgnoreLine"
    - from: "internal/tui/lanes.go"
      to: "core/lane.Publish"
      via: "the export keybind on the lane settings screen"
      pattern: "lane\\.(Publish|Adopt|Export)"
---

<objective>
Lanes stop being a fixed thing baked into the binary. After this phase a person
or an agent can read a lane in full, write a new one from a template, scope a
project to the lanes it actually uses, override a built-in deliberately, and
hand a lane to a teammate through the repository.

Purpose: the lane is the unit of pipeline design. Today it can only be listed,
never read, never changed, never shared — which makes the pipeline something
the tool has an opinion about rather than something the user owns.

Output: a root-aware lane loader with a three-source layout, a legible `jaira
lanes` surface, a gitignore rule that keeps project prompts private through
`jaira share`, and a TUI lane settings screen that publishes and adopts.
</objective>

<execution_context>
@$HOME/.claude/get-shit-done/workflows/execute-plan.md
@$HOME/.claude/get-shit-done/templates/summary.md
</execution_context>

<context>
@.planning/lane-system-design.md
@.planning/ROADMAP.md
@.planning/REQUIREMENTS.md
@./CLAUDE.md

@core/lane/lane.go
@internal/cli/share.go
@core/board/gitignore.go
</context>

<phase_goal>
**As a** person whose pipeline does not match the one that shipped, **I want to**
read, write, reorder and hand over lane definitions as plain files, **so that**
the board's steps are mine and my teammates can run them too.
</phase_goal>

<execution_notes>
Go is not on the default PATH in this environment. Every command in this plan is
written as `export PATH=$PATH:$HOME/.local/go/bin && go ...` — keep the prefix.

Each task below is a commit boundary. Stop, run the task's verify, and commit
before starting the next one. Task 2 is the largest (a signature change across
16 call sites) and should be executed in its own session; if a single session is
carrying tasks 2 and 3 together, split it.

Three quick tasks are landing in parallel on `core/ticket/frontmatter.go`,
`core/ticket/schema.go`, `core/ticket/checklist.go`, `internal/cli/tickets.go`,
`internal/cli/checklist.go`, `internal/tui/view.go`, `internal/tui/signoff.go`,
`core/board/*` and three files under `core/lane/builtin/`. Do not trust line
numbers quoted anywhere in this plan; locate code by identifier. This plan
deliberately does not edit any file under `core/lane/builtin/` (see task 3).
</execution_notes>

<scope_discipline>
The project constraint is "is this smaller than paca?". Five places where this
phase is deliberately built smaller than the design note could be read to imply:

1. **`creator:` is not written into the nine built-in lane files.** `parse`
   defaults it to `jaira` for built-ins instead. Same observable result, nine
   fewer file edits, and no collision with the quick tasks editing three of
   those files right now.
2. **`jaira lanes template` prints to stdout.** It does not create a file, ask
   where to put it, or scaffold a directory. `jaira lanes template > $(jaira
   lanes path)/my-lane.md` is the whole workflow, and it composes.
3. **Drift between a project lane and its catalogue original is a warning, not
   a sync feature** (subject to the decision in task 1). A sync would have to
   answer "which direction", and neither answer is right for both cases.
4. **No `jaira lanes new/edit/move/rm`.** Explicitly rejected by the user. The
   file is the API; `jaira lanes` is how you check the file.
5. **"Use this lane in this project" is one keybind on a screen being built
   anyway.** It is `cp` with a confirmation. It is not worth its own command,
   and would not be worth the keybind either if the settings screen did not
   already have the lane in hand.
</scope_discipline>

<migration>
Anyone with lanes already in `~/.jaira/lanes` keeps them, unmoved and unrenamed.
That directory becomes "the catalogue" by definition rather than by migration,
and `JAIRA_LANES_DIR` keeps overriding its location exactly as today.

Two behaviour changes to state in the summary and the release note:

- A custom lane whose id collides with a built-in used to be **ignored** with a
  warning. It now **overrides** with a warning. Anyone who had such a file was
  getting the built-in; after this phase they get theirs. The warning names the
  file, the id and the built-in it displaced, so this is visible on the first
  command rather than discovered later.
- If a `.jaira/lanes/` directory exists in a project, the catalogue is not
  loaded for that project (subject to decision D-03 in task 1). No existing
  install has such a directory, so nothing changes until someone creates one.
</migration>

<tasks>

<task type="checkpoint:decision" gate="blocking">
  <name>Task 1: Settle the three open design questions</name>
  <decision>Three questions the design note leaves open, all of which change observable behaviour and none of which should be picked silently by an implementer.</decision>
  <context>
Two are recorded at the bottom of `.planning/lane-system-design.md`. The third
surfaced while reading the code and is the most consequential of the three.
  </context>
  <options>
    <option id="D-01">
      <name>D-01 — What identifies `&lt;user&gt;` in `.jaira/shared/&lt;user&gt;/`</name>
      <pros>
**Recommended: reuse the existing `identity()` function, slugified.**

`internal/cli/root.go` already has `identity()`: `JAIRA_USER` env var, then
`git config user.name`, then `$USER`/`$USERNAME`/`$LOGNAME`, then `"unknown"`.
It is already how tickets are attributed. Reusing it means the shared folder and
the ticket's assignee say the same thing about the same person, there is no new
setting to explain, and `JAIRA_USER` is already the documented escape hatch for
anyone the default gets wrong.

Implementation: move `identity()` into a new `core/identity` package (both
`internal/cli` and `internal/tui` need it, and `internal/tui` cannot import
`internal/cli` without a cycle), add `Slug()` that lowercases and reduces to
`[a-z0-9-]` so it is safe as a path component, and leave `internal/cli` calling
through to it so ticket attribution is unchanged.
      </pros>
      <cons>
The trade-off, stated plainly: two teammates both configured as "Alex" collide
into one folder, and the second one to publish overwrites the first. The
mitigations are that the collision is visible (the folder is committed, so it
shows up in a diff and in review) and that `JAIRA_USER` fixes it in one line.

Alternatives if that is not good enough:
- **Slug of `git config user.email`** — effectively unique, since git already
  requires it to commit. Costs: an email address becomes a committed directory
  name, which not everyone wants published in a repo they do not control.
- **`&lt;name&gt;-&lt;6 hex of email hash&gt;`** — readable and unique. Costs: an
  uglier path, and a folder name that changes if someone changes their git email,
  orphaning their previous folder.
      </cons>
    </option>
    <option id="D-02">
      <name>D-02 — What happens when a project lane and its catalogue original drift apart</name>
      <pros>
**Recommended: warn at load, do nothing else.**

When a lane id exists in both `.jaira/lanes/` and the catalogue and the file
bytes differ, emit `lane &lt;id&gt;: this project's copy differs from your catalogue
copy`. The warning channel already exists and already reaches `jaira lanes`,
`--json`, the TUI warnings block and every command through `loadEnv`. Cost is
roughly ten lines and one test.
      </pros>
      <cons>
The trade-off: the user is told about the drift and then has to resolve it with
`cp`, which the tool will not do for them.

Alternatives:
- **Offer to sync from the settings screen** — but "sync" has a direction, and
  the tool cannot know which copy is the intended one. A project copy edited for
  this repo and a catalogue copy improved since the export are both legitimate,
  and pointing the arrow the wrong way destroys work.
- **Ignore it** — cheapest, and defensible since the project copy is
  authoritative anyway. Rejected because silent divergence between two files
  with the same id is exactly the confusion `creator:` was added to help with.
      </cons>
    </option>
    <option id="D-03">
      <name>D-03 — Is `.jaira/lanes/` additive to the catalogue, or authoritative over it? (not in the design note)</name>
      <pros>
**Recommended: authoritative when present.** If `.jaira/lanes/` exists and
contains at least one `.md` file, load built-ins plus the project's lanes and do
not read the catalogue for that project. Otherwise load built-ins plus the
catalogue, exactly as today.

Why this is not a free choice: success criterion 1 says the project directory
holds "copies of only the lanes that project actually uses". Under an additive
reading, a user with twenty catalogue lanes still sees all twenty in every
project, and the project directory changes nothing except which copy wins — the
criterion is not delivered. Authoritative-when-present is the only reading under
which exporting a lane into a project means anything.

It also removes the reason to think of `.jaira/lanes/` as a sharing mechanism.
It is gitignored, so a teammate never receives it; its only job is scoping this
machine's catalogue down to this project. Authoritative is that job.
      </pros>
      <cons>
The trade-off: creating `.jaira/lanes/` and putting one file in it silently
removes nineteen lanes from that project's board. That is a sharp edge.

Mitigations if this option is chosen: `jaira lanes` prints which source is
active, `jaira lanes path` names both directories and says which one is
currently in force, and the loader warns when a project directory exists but is
empty (which would otherwise look like "my lanes vanished").

The alternative — **additive, project overrides catalogue** — is gentler and
never hides a lane, at the cost of not delivering criterion 1. If this is
chosen, say so, because criterion 1's wording then needs to change with it.
      </cons>
    </option>
  </options>
  <resume-signal>Answer all three: D-01 (identity / email / name+hash), D-02 (warn / sync / ignore), D-03 (authoritative / additive). Tasks 2, 6 and 7 are blocked until these are settled.</resume-signal>
  <files>.planning/phase-5-custom-and-portable-lanes/PLAN.md</files>
  <action>
Put D-01, D-02 and D-03 to the user with the recommendation and the trade-off for
each, then write the chosen answers into this task as a `&lt;decided&gt;` block naming
the option id, the choice, and the date.

This task writes no code. It exists as a task rather than as a note because the
three answers change observable behaviour in tasks 2, 6 and 7, and an implementer
who picks one silently would bury a product decision in a commit. D-03 in
particular decides whether success criterion 1 in ROADMAP.md is deliverable as
written or has to be reworded, so it must be settled before task 2 starts.
  </action>
  <verify>
    <automated>grep -q "&lt;decided&gt;" .planning/phase-5-custom-and-portable-lanes/PLAN.md</automated>
  </verify>
  <done>All three decisions are recorded in this task as a `&lt;decided&gt;` block, and if D-03 went to "additive" then ROADMAP.md criterion 1 has been reworded to match.</done>
</task>

<task type="auto" tdd="true">
  <name>Task 2: Root-aware lane loading, the three-source layout, and overridable built-ins</name>
  <files>core/lane/lane.go, core/lane/lane_test.go, internal/cli/root.go, internal/cli/tickets.go, internal/cli/checklist.go, internal/cli/resume.go, internal/cli/validate.go, internal/cli/mergedriver.go, internal/cli/json_test.go, internal/tui/model.go, internal/tui/home.go, core/gate/plan_test.go, core/validate/validate_test.go, core/merge/merge_test.go, internal/tui/signoff_test.go, .planning/REQUIREMENTS.md</files>
  <behavior>
The `core/lane` package has no tests at all today. `core/lane/lane_test.go` is
created in this task and is the acceptance surface for it:

- Built-ins alone: `Load("")` with an empty catalogue returns the nine shipped
  lanes in filename order.
- Catalogue: a lane file in the catalogue dir loads and orders by its `after:`.
- Project directory, per D-03: a lane in `&lt;root&gt;/.jaira/lanes/` loads; and if
  D-03 is authoritative, a catalogue lane absent from the project directory does
  not appear for that root, while the same catalogue lane does appear for a root
  with no project directory.
- Override: a custom lane with `id: review` replaces the built-in Review lane,
  including its prompt, and produces a warning naming the file, the id and the
  built-in displaced. Assert the warning text, not just its presence.
- Gate-weakening override: a lane that overrides a built-in and drops
  `requires-human-exit` or `terminal` evidence produces a second, distinct
  warning saying which protection was removed (see T-5-01 in the threat model).
- Duplicate ids within one directory still resolve deterministically to the
  first by sorted filename, with a warning, as today.
- Empty project directory (per D-03 mitigation): warns rather than silently
  yielding a built-ins-only board.
- Drift, per D-02: same id in project and catalogue with different bytes warns.
  </behavior>
  <action>
Change `Load()` to `Load(root string)` and add `ProjectLanesDir(root string)`
returning `filepath.Join(root, ticket.DirName, "lanes")` — `core/lane` already
imports `core/ticket`, so use `ticket.DirName` rather than a second literal.
`root == ""` means "no project in hand" and loads built-ins plus catalogue.

Resolution order is built-ins, then the source chosen by D-03, with later
definitions of the same id replacing earlier ones. Replace the collision refusal
in `Load` (the branch that appends "collides with a built-in lane and was
ignored" and `continue`s) with a replacement that keeps the custom lane, records
the displaced lane's origin, and warns. Update the package doc comment and the
comment above that branch: both currently assert that a custom lane must never
shadow a built-in, and after this change those comments are false. Rewrite them
to say why overriding is allowed and why it is loud — the reasoning in the
design note is what belongs there.

Add the gate-weakening warning: when an override replaces a built-in that had
`RequiresHumanExit`, `Terminal`, `RequiresOutcome` or `RequiresNonModelSignal`
set and the replacement does not, name the removed protection in the warning.
This is cheap and it is the one place where "powerful but never silent" has to
mean more than "a line scrolled past".

Then update all sixteen `lane.Load()` call sites. Pass a root wherever one is in
hand: `s.Root` in `internal/cli` (`loadEnv`, `tickets.go`, `checklist.go`,
`resume.go`, `validate.go`, `mergedriver.go`) and `m.store.Root` in
`internal/tui/model.go`. `internal/tui/home.go` is the launcher and spans many
projects with no single root — pass `""` there, and add a comment saying the
launcher can only show catalogue lanes because it is not inside one board.
Test call sites pass `""` or the test's temp root as appropriate.

Finally, correct `.planning/REQUIREMENTS.md` LANE-06. It currently reads "A lane
whose `id` collides with a base lane is rejected with a clear error, never
silently overriding it", which this task deliberately reverses. ROADMAP.md
already records that the design note supersedes the original criterion; make
REQUIREMENTS.md say the same thing: a lane whose id collides with a base lane
overrides it and the override is always reported as a warning, never silent.
Shipping code that contradicts a live requirement is worse than either choice.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go test ./core/lane/... ./core/gate/... ./core/validate/... ./core/merge/... ./internal/... 2>&amp;1 | tail -20</automated>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && grep -rn "lane.Load()" --include=*.go . | grep -v "^./.git" | wc -l | grep -qx 0 && echo "no zero-arg Load remains"</automated>
  </verify>
  <done>`core/lane/lane_test.go` exists and covers every case in the behavior block; `go build ./...` passes with no zero-argument `lane.Load()` anywhere; a lane file overriding a built-in changes the board and produces a named warning; REQUIREMENTS.md LANE-06 no longer contradicts the shipped behaviour.</done>
</task>

<task type="auto" tdd="true">
  <name>Task 3: `creator:` provenance and a legible `jaira lanes`</name>
  <files>core/lane/lane.go, core/lane/lane_test.go, internal/cli/tickets.go, internal/cli/lanes_test.go</files>
  <behavior>
In `core/lane/lane_test.go`:
- `creator:` in a lane file's frontmatter is parsed onto `Lane.Creator`.
- A built-in with no `creator:` field reports `jaira`.
- A custom lane with no `creator:` reports empty, not `jaira` — absent
  provenance and "shipped by the tool" are different facts.

In a new `internal/cli/lanes_test.go`, driving the cobra command with a
`JAIRA_LANES_DIR` and a temp root, as `internal/cli/json_test.go` already does:
- `jaira lanes --json` includes `prompt` and `creator` for every lane.
- `jaira lanes show review` prints the id, name, anchor, precedence, tier,
  input/output contract, creator, source, and the full prompt body.
- `jaira lanes show nope` exits 2 with reason `no_such_lane`.
- `jaira lanes show review --json` is valid JSON carrying the same fields.
- `jaira lanes path` prints the catalogue directory and the project directory,
  and marks which is currently in force.
- `jaira lanes template` prints a lane file that `lane.parse` accepts — assert
  this by parsing the command's own stdout, so the template cannot rot away from
  the parser.
  </behavior>
  <action>
Add `Creator string` to `Lane`, parsed from `creator:` in `parse`. Default it to
`jaira` when `builtin` is true and the field is absent, rather than editing the
nine files under `core/lane/builtin/` — three of those files are being edited by
quick tasks in parallel, and the default is the same observable result for no
merge risk. Note the default in a comment so the next reader does not go looking
for a `creator:` line in the built-in files and conclude it is missing.

Turn `newLanesCmd` into a parent command that keeps its current bare-`jaira
lanes` table as its own `RunE`, and add three subcommands:

- `lanes show <id>` — the whole lane including the prompt body. Human output is
  a labelled block then the prompt; `--json` carries the same fields as the list
  plus `prompt`, `creator`, `after`, `description`, and `overrides` (the id of a
  built-in this lane displaced, empty otherwise). Unknown id fails with
  `ExitUsage` and reason `no_such_lane`, matching how `move` already reports it.
- `lanes path` — prints `lane.UserLanesDir()` and `lane.ProjectLanesDir(root)`,
  each labelled, with a marker on whichever is in force for this directory under
  D-03. Under `--json`, emit both as named fields plus which is active. This is
  what an agent calls before writing a file, so it must work outside a board too
  (no project root: print the catalogue and say there is no project).
- `lanes template` — writes a commented lane skeleton to stdout and nothing
  else. Every field `parse` reads, each with a one-line comment, and a prompt
  body section. It must state that `model-tier` is a local alias (`cheap`,
  `strong` — the two the built-ins use) and not a model name, because that
  contract is invisible from the field alone and is what makes a shared lane
  file survive a model rename.

Add `prompt` and `creator` to the existing `lanes --json` object. Add `CREATOR`
to the human table only if it fits the existing column budget; if it does not,
leave the table alone and let `lanes show` carry it — the table is already six
columns wide.

Keep `Args: noArgs()` on the parent. Cobra dispatches a matching subcommand name
before argument validation runs, so `jaira lanes show x` works and `jaira lanes
typo` falls through to the parent and exits 2 as a usage error. Add a test for
the `typo` case so this stays true across cobra upgrades.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go test ./internal/cli/... ./core/lane/... 2>&amp;1 | tail -20</automated>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go run ./cmd/jaira lanes template | go run ./cmd/jaira lanes show review >/dev/null; go run ./cmd/jaira lanes --json | grep -q '"prompt"' && echo "prompt in json"</automated>
  </verify>
  <done>`jaira lanes show review` prints the Review lane's full prompt with no ticket involved; `jaira lanes --json` carries prompt and creator; `jaira lanes path` names both directories; `jaira lanes template` emits something `parse` accepts, proven by a test that parses the command's own output.</done>
</task>

<task type="auto" tdd="true">
  <name>Task 4: Project lanes stay private through `jaira share`</name>
  <files>core/board/gitignore.go, core/board/board_test.go, internal/cli/share.go, core/ticket/store.go</files>
  <behavior>
In `core/board/board_test.go`, against a temp root:
- On a fresh `.gitignore`, `RemoveIgnore` followed by `AddLanesIgnore` leaves a
  file that does not ignore `/.jaira/` but does ignore `/.jaira/lanes/`.
- `AddLanesIgnore` is idempotent: calling it twice reports `changed=false` the
  second time and writes no duplicate line.
- `AddIgnore` (the `share --undo` path) removes the `/.jaira/lanes/` line, since
  `/.jaira/` already covers it and a leftover line is a puzzle for the next
  reader.
- `Ignored(root)` still reports true only for the whole-board entry and is not
  fooled into true by the lanes-only line. This is the one that matters:
  `isShared` and `bindDriverIfShared` both branch on it, so a false positive
  would silently disable the merge driver on a shared board.
  </behavior>
  <action>
Add `LanesIgnoreLine = "/.jaira/lanes/"` next to `IgnoreLine` in
`core/board/gitignore.go`, with `AddLanesIgnore(root)` and
`RemoveLanesIgnore(root)` following the shape of the existing pair. This is an
ordinary ignore rule on a path outside any ignored tree, so no negation pattern
and no nested-gitignore trickery is involved — say so in a comment, because the
first instinct on reading "ignore one child of a directory we just un-ignored"
is to reach for `!` patterns that do not work inside an ignored parent.

In `internal/cli/share.go`: after `RemoveIgnore` succeeds, call
`AddLanesIgnore`, and report it in both output modes — a new `lanes_ignored`
field in the `--json` object, and a line in the human output saying the
project's lanes stay private. Change the closing instruction from `git add
.jaira .gitignore` to name what is actually being published, since `git add
.jaira` now stages a directory whose `lanes/` child is ignored and that is worth
being explicit about rather than leaving to the user's model of gitignore. In
the `--undo` branch, call `RemoveLanesIgnore` after `AddIgnore`.

Add `SharedSubdir = "shared"` to `core/ticket/store.go` beside `TicketsSubdir`
and `ArchiveSubdir`, plus a `SharedDir()` accessor matching `ArchiveDir()`.
Nothing creates that directory yet — task 6 does, on first publish. Creating it
empty at `init` would commit an empty directory's worth of confusion to every
board that never shares a lane.

While in `store.go`: the comment above `SessionsSubdir` currently claims
`.jaira/` "contains only content that is meant to be committed, so there is no
mixed directory to explain and no per-repo gitignore needed to protect scratch
state". This phase makes that false — `.jaira/lanes/` is deliberately never
committed. Update the comment to state the new rule: `tickets/`, `archive/` and
`shared/` are committed, `lanes/` is this machine's, and sessions and locks stay
out of the repo entirely.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go test ./core/board/... ./core/ticket/... ./internal/cli/... 2>&amp;1 | tail -20</automated>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && cd $(mktemp -d) && git init -q . && git config user.email t@t && git config user.name t && go run github.com/BeMuCa/jaira/cmd/jaira init >/dev/null 2>&1; mkdir -p .jaira/lanes && echo x > .jaira/lanes/x.md && go run github.com/BeMuCa/jaira/cmd/jaira share >/dev/null && git add -A && git status --porcelain | grep -qv "\.jaira/lanes" && echo "lanes not staged"</automated>
  </verify>
  <done>After `jaira share` on a real repo, `git add -A` stages `.jaira/tickets/` and leaves `.jaira/lanes/` unstaged; `share --undo` restores a single whole-board ignore with no orphan lanes line; `Ignored()` is unchanged in meaning.</done>
</task>

<task type="auto" tdd="true">
  <name>Task 5: Prove ordering, tier aliasing and unknown-lane handling — and close whatever the proof breaks</name>
  <files>core/lane/lane_test.go, core/lane/lane.go, core/gate/plan_test.go</files>
  <behavior>
Three success criteria (3, 4 and 6) describe behaviour that reading the code
suggests is already implemented. This task's job is to prove that with tests
rather than assert it in a summary, and to fix whatever the tests find. Write
the tests first; if they all pass on the first run, that is the deliverable and
the task is a small one — say so honestly in the summary rather than
manufacturing a change.

Ordering (criterion 4, LANE-05):
- A custom lane with `after: human` lands immediately after Human.
- A chain — B anchored to A, A anchored to a built-in — resolves regardless of
  filename order, which is what the fixed-point loop in `order()` is for.
- A lane anchored to an id that is not installed lands before the terminal lane
  and warns, naming the missing anchor. This is the shared-lane-file case and is
  the one most likely to regress.
- Two lanes anchored to each other produce the cycle warning and both still
  appear, so a bad file never removes a column.
- A lane with no `after:` at all lands before the terminal lane.

Model tier (criterion 3, LANE-04):
- An agentic lane with no `model-tier` is a parse error (already true).
- A lane's `ModelTier` round-trips as the exact alias string given, and is not
  resolved, rewritten or validated against a list of model names anywhere.
  Assert this with a grep-style guard test: `ModelTier` must never be compared
  against a string containing `claude`, `gpt`, `sonnet`, `opus` or `haiku` in
  non-test source. A shared lane file surviving a model rename depends entirely
  on nothing in this repo knowing what a model is called.

Unknown lanes (criterion 6, LANE-07, LANE-08):
- `Set.Columns` appends a passthrough column for a status no lane claims, and
  sorts multiple unknowns deterministically.
- `Set.Precedence` returns -1 for an unknown lane, so a merge never promotes a
  ticket into one.
- In `core/gate/plan_test.go`: moving a ticket *out of* an unrecognised lane is
  refused with `CodeUnknownLane`, and moving one *into* an id no lane claims is
  refused with `CodeNoSuchLane`. Both paths exist in `gate.Decide`; neither
  appears to be covered.
  </behavior>
  <action>
Write the tests above. Fix only what they break, and prefer the smallest fix
that makes the criterion true.

If the tier guard test is awkward to express as a Go test, express it as a
`go vet`-adjacent shell check inside the test using `os/exec` over `go list`
output, or as a plain `filepath.WalkDir` over `.go` files skipping `_test.go`
and `builtin/`. Either is fine; a test that reads the source tree is the honest
way to assert an absence.

Do not add tier validation against a fixed alias list. `cheap` and `strong` are
what the built-ins use, but a user inventing `local` or `nano` for their own
runner is exactly the point of an alias, and a whitelist would break them. The
template documents the convention (task 3); the code does not enforce it.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go test ./core/lane/... ./core/gate/... -run 'Order|Anchor|Tier|Unknown|Passthrough|Cycle' -v 2>&amp;1 | tail -30</automated>
  </verify>
  <done>Every one of criteria 3, 4 and 6 has a named test asserting it; a lane anchored to a missing lane orders and warns; an unrecognised lane renders as a column and refuses movement in both directions; no non-test source compares a model tier against a model name.</done>
</task>

<task type="auto">
  <name>Task 6: Lane settings screen — read a lane, use it here, publish it</name>
  <files>core/identity/identity.go, core/identity/identity_test.go, core/lane/share.go, core/lane/share_test.go, internal/cli/root.go, internal/tui/lanes.go, internal/tui/lanes_test.go, internal/tui/model.go, internal/tui/view.go</files>
  <action>
Implements D-01. Do not start until task 1 is answered.

First, the core, so the TUI stays thin and the behaviour is unit-testable
without a terminal:

`core/identity/identity.go` — move `identity()` out of `internal/cli/root.go`
(the TUI cannot import `internal/cli` without a cycle) as `identity.Who(dir
string) string`, preserving the existing precedence exactly: `JAIRA_USER`, then
`git config user.name` in `dir`, then `USER`/`USERNAME`/`LOGNAME`, then
`unknown`. Add `Slug(string) string` reducing to `[a-z0-9-]` with runs collapsed
and edges trimmed, falling back to `unknown` if nothing survives. Leave
`internal/cli`'s `identity()` in place as a one-line call through, so ticket
attribution is provably unchanged. Test `Slug` against names with spaces, dots,
accents, a leading dot, path separators and an empty string — this value becomes
a directory name, so `../` and `/` must not survive it (T-5-03).

`core/lane/share.go` — three pure functions over file paths, each copying a lane
file and returning the destination path:
- `Export(l *Lane, dstDir string) (string, error)` — catalogue to project.
- `Publish(l *Lane, dstDir string) (string, error)` — to `.jaira/shared/<slug>/`.
- `Adopt(srcPath, dstDir string) (*Lane, string, error)` — parse first, then
  copy, so an unparseable file never lands in the catalogue.

All three derive the filename from the lane's validated `ID` (already
constrained to `[a-z0-9-]` by `validID`), never from the source filename, and
all three `MkdirAll` the destination. Refuse to overwrite an existing file
unless an explicit `overwrite bool` is passed; the caller asks. Copy bytes
verbatim — do not parse and re-serialise, because the file format is the API and
round-tripping through a struct would reorder fields and drop unknown keys.
Stamp `creator:` on publish only if the file has no `creator:` line, and do it
as a line insert after the opening `---`, not a YAML rewrite.

Then the screen. `internal/tui/lanes.go`, following `internal/tui/browse.go`'s
shape exactly — a struct holding its own state, a `key(string) (…, done bool)`
method, and a `render(width, height int) string`; no bubbletea imports of its
own. Add `modeLanes` to the `mode` enum in `model.go`, a case in `key()`, a case
in the view dispatch, and a keybind from the board (`L` is free; confirm against
the current help screen and pick another if not).

The screen lists every loaded lane with its source (built-in, catalogue,
project) and marks overrides. Selecting one shows its full prompt in a scrolling
pane — this is `jaira lanes show` inside the TUI and should read from the same
`Lane` fields so the two cannot disagree. Keybinds: `u` use this lane in this
project (`Export` into `ProjectLanesDir`), `p` publish (`Publish` into
`SharedDir()/Slug(Who(root))`), `esc` back. Both actions report the written path
in the message line, and both refuse rather than overwrite, saying so.

Test `internal/tui/lanes_test.go` at the model level, as the existing TUI tests
do: drive `key()` and assert on state and on the files that appear on disk, not
on rendered pixels. `core/lane/share_test.go` covers the copy semantics.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go test ./core/identity/... ./core/lane/... ./internal/tui/... 2>&amp;1 | tail -20</automated>
    <human-check>Open the board, press the lane settings key, select the Review lane, confirm its full prompt is readable and scrolls. Press `p`. Confirm the message names a path under `.jaira/shared/`, that the file is there, and that `git status` shows it as untracked-but-not-ignored. Press `p` again and confirm it refuses rather than overwriting.</human-check>
  </verify>
  <done>`core/identity` exists and `internal/cli`'s attribution is unchanged; `Export`/`Publish`/`Adopt` are unit-tested including the refuse-to-overwrite and path-safety cases; the lane settings screen reads a lane's prompt and writes it to both destinations, reporting the path.</done>
</task>

<task type="auto">
  <name>Task 7: Adopt a teammate's shared lane</name>
  <files>internal/tui/lanes.go, internal/tui/lanes_test.go, core/lane/share.go, core/lane/share_test.go, internal/cli/tickets.go</files>
  <action>
Implements the second half of criterion 8. Depends on task 6.

Add `Shared(root string) ([]SharedLane, error)` to `core/lane/share.go`: walk
`.jaira/shared/*/*.md`, parse each, and return the lane alongside the folder
name it came from and its path. A file that fails to parse is skipped with a
warning rather than failing the walk — one teammate's broken file must not hide
everyone else's lanes. Shared lanes are **not** loaded by `Load`: they are
visible, adoptable, and inert until adopted. That separation is the whole
security story for T-5-02 and must be stated in the function's comment, because
the obvious "helpful" change later is to load them automatically.

On the lane settings screen, add a second section listing shared lanes grouped
by folder, each labelled with its `creator:` where present. Selecting one shows
its full prompt in the same pane the local lanes use — the prompt must be read
before adopting, since adopting means agreeing to run someone else's
instructions at whatever tier the file declares. Press `a` to adopt: `Adopt`
copies it into `UserLanesDir()`, refusing if that id already exists there unless
confirmed, and the message names the catalogue path written.

Make the same list available without the TUI: `jaira lanes --shared` (or a
`lanes shared` subcommand, whichever reads better against the surface built in
task 3) listing shared lanes with folder, id, creator and path, honouring
`--json`. An agent has no way to press `a`, and "read what teammates published"
is a read operation the CLI already promises for everything else.

Test: build a temp root with two `.jaira/shared/<name>/` folders including one
unparseable file, assert `Shared` returns the good ones and warns about the bad
one, assert adoption writes into a temp `JAIRA_LANES_DIR` and that the adopted
lane then appears in `Load`.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go test ./core/lane/... ./internal/tui/... ./internal/cli/... 2>&amp;1 | tail -20</automated>
    <human-check>With a lane published under `.jaira/shared/someone-else/`, open the lane settings screen and confirm it appears in the shared section with its creator. Select it, read the prompt, press `a`, and confirm the message names a path in `~/.jaira/lanes/` and that the lane then appears in `jaira lanes` for a different project.</human-check>
  </verify>
  <done>Shared lanes are listed in the TUI and by the CLI, are never loaded onto a board until adopted, show their prompt before adoption, and adopting copies the file into the catalogue where `Load` then finds it.</done>
</task>

</tasks>

<threat_model>
No packages are installed in this phase — Go standard library and existing
module dependencies only — so the package legitimacy gate does not apply.

## Trust Boundaries

| Boundary | Description |
|----------|-------------|
| repository → lane loader | `.jaira/shared/**` is committed markdown authored by anyone who can push to the repo, read by a tool that turns lane bodies into agent prompts |
| lane file → subagent | a lane's markdown body becomes an instruction executed at a self-declared `model-tier` |
| lane `id` / `<user>` → filesystem | both become path components under `.jaira/` and `~/.jaira/` |
| override → gate | a custom lane replacing a built-in also replaces the gate conditions that built-in enforced |

## STRIDE Threat Register

| Threat ID | Category | Component | Disposition | Mitigation Plan |
|-----------|----------|-----------|-------------|-----------------|
| T-5-01 | Elevation of Privilege | `core/lane.Load` override path | mitigate | A lane overriding a built-in can drop `requires-human-exit`, `terminal`, `requires-outcome` or `requires-nonmodel-signal` and thereby let an agent sign off its own work — the exact failure `gate.Decide` exists to prevent. Task 2 emits a distinct warning naming the removed protection, and task 3's `lanes show` marks overridden built-ins. Tested in `core/lane/lane_test.go`. |
| T-5-02 | Tampering | `core/lane.Shared` / adopt path | mitigate | A committed shared lane is untrusted prompt content that would otherwise run at `model-tier: strong`. Shared lanes are never loaded by `Load` (task 7), adoption is an explicit keypress, and the full prompt is displayed before adoption. |
| T-5-03 | Tampering | `identity.Slug`, lane filename derivation | mitigate | `<user>` and lane ids become path components. `Slug` reduces to `[a-z0-9-]` and cannot yield `..` or a separator; `share.go` derives filenames from the already-`validID`-constrained `Lane.ID`, never from a source filename. Tested with traversal inputs. |
| T-5-04 | Spoofing | `creator:` frontmatter | accept | `creator:` is self-declared and unverifiable without signing. It is provenance, not authentication. Document it as such in `lanes template` so nobody trusts it as an authorisation signal. |
| T-5-05 | Information Disclosure | `jaira share` / `.gitignore` | mitigate | Project lane prompts may contain repo-specific detail the author did not intend to publish. Task 4 writes `/.jaira/lanes/` on share and proves with a real-repo test that `git add -A` does not stage it. |
| T-5-06 | Denial of Service | `core/lane.Load` ordering | accept | A malformed or cyclic set of lane files could in principle produce a confusing board. `order()` already terminates on cycles and warns, and task 5 tests it. The blast radius is one user's own board. |
</threat_model>

<coverage_audit>
Every success criterion, requirement and design-note decision maps to a task.

## GOAL

| Item | Covered by |
|---|---|
| Read, write, reorder and share lane definitions | Tasks 2–7 collectively |

## Success criteria (ROADMAP Phase 5)

| # | Criterion | Task |
|---|---|---|
| 1 | One md file; `~/.jaira/lanes/` catalogue; `.jaira/lanes/` project copies | 2 (D-03 decides authoritative vs additive) |
| 2 | `.jaira/lanes/` gitignored even after `share` | 4 |
| 3 | `model-tier` is a local alias, not a model name | 5 (proof + guard test), 3 (documented in the template) |
| 4 | Anchor ordering, including an anchor absent locally | 5 |
| 5 | A custom lane MAY override a built-in, always warned | 2 |
| 6 | Unknown lane renders as passthrough; movement refused | 5 |
| 7 | `lanes show` / prompt in `--json` / `lanes path` / a template | 3 |
| 8 | Export to `.jaira/shared/<user>/`; teammates' lanes visible; adopt | 6 (publish), 7 (visible + adopt) |
| 9 | Every lane records a `creator:` | 3 |

## REQ (REQUIREMENTS.md)

| ID | Task | Note |
|---|---|---|
| LANE-02 | 2 | portable single file; the catalogue is unchanged, the project dir is new |
| LANE-03 | 3 | prompt, tier and contract already parsed; this makes them readable |
| LANE-04 | 5 | proven, plus a guard against a future model-name comparison |
| LANE-05 | 5 | |
| LANE-06 | 2 | **reversed by the design note.** Task 2 rewords REQUIREMENTS.md so the requirement and the code agree |
| LANE-07 | 5 | |
| LANE-08 | 5 | |

## DESIGN (`.planning/lane-system-design.md`)

| Decision | Task |
|---|---|
| Three-location storage layout | 2, 4 |
| Project holds copies, not references | 6 (`Export`) |
| Project lanes private by default | 4 |
| Adopting copies into the adopter's catalogue | 7 |
| `share` writes `/.jaira/lanes/` | 4 |
| Built-ins become overridable, loudly | 2 |
| `lanes show`, prompt in `--json`, `lanes path`, a template | 3 |
| No `lanes new/edit/move/rm` | not built, by instruction |
| `creator:` signature | 3 |
| TUI lane settings screen | 6 |
| Lane selection shows and adopts shared lanes | 7 |
| Open question: `<user>` identity | 1 (D-01) |
| Open question: drift | 1 (D-02) |

## Surfaced, not in any source artifact

| Item | Task | Why it is here |
|---|---|---|
| Additive vs authoritative project lanes | 1 (D-03) | criterion 1 is only delivered under one of the two readings; not settled anywhere |
| `lane.Load()` must take a root | 2 | 16 call sites; unavoidable consequence of a per-project lane directory |
| REQUIREMENTS.md LANE-06 contradicts the design | 2 | shipping code against a live requirement is a defect either way |
| `store.go`'s "only committable content" comment becomes false | 4 | the invariant this phase breaks should not stay written down as true |
| A gate-weakening override is a privilege escalation | 2, threat model | overriding built-ins was approved; overriding `requires-human-exit` specifically was not discussed |
</coverage_audit>

<verification>
Full suite, from a clean tree:

    export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go test ./... 2>&1 | tail -30

End-to-end, in a throwaway repo:

1. `jaira init`, then `jaira lanes` — nine built-ins, no warnings.
2. `jaira lanes template > $(jaira lanes path --json | jq -r .catalogue)/hitl.md`,
   edit `id`/`after`, then `jaira lanes` — the new lane appears, ordered.
3. Copy it to `id: review` — `jaira lanes` warns that a built-in was overridden
   and names it; `jaira lanes show review` prints the new prompt.
4. `jaira share`, then `git add -A && git status --porcelain` — tickets staged,
   `.jaira/lanes/` absent.
5. From the TUI lane settings screen, publish a lane; confirm the file lands in
   `.jaira/shared/<you>/` and is stageable.
6. Rename that folder to simulate a teammate, restart, adopt it, and confirm it
   appears in `jaira lanes` from a different project.
</verification>

<success_criteria>
- `go build ./...` and `go test ./...` pass.
- `core/lane` has tests where it had none.
- All nine ROADMAP Phase 5 success criteria have a named test or a recorded
  human check.
- REQUIREMENTS.md LANE-06 no longer contradicts the shipped behaviour.
- No lane, on any path, is loaded onto a board because a teammate committed it.
- `jaira share` on a real repo leaves `.jaira/lanes/` unstaged.
</success_criteria>

<output>
Create `.planning/phase-5-custom-and-portable-lanes/05-01-SUMMARY.md` when done.
Record in it: the answers chosen for D-01, D-02 and D-03; the two behaviour
changes from the migration section, for the release note; and whether task 5's
tests passed on first run or found real gaps.
</output>
