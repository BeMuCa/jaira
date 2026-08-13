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
  - core/lane/defaultboard.go
  - core/lane/defaultboard_test.go
  - core/lane/builtin/60-blocked.md
  - core/ticket/store.go
  - core/ticket/body.go
  - core/ticket/body_test.go
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
  - internal/cli/init_test.go
  - internal/cli/flow.go
  - internal/tui/model.go
  - internal/tui/home.go
  - internal/tui/lanes.go
  - internal/tui/lanes_test.go
  - internal/tui/defaultboard.go
  - internal/tui/defaultboard_test.go
  - internal/tui/view.go
  - core/gate/plan_test.go
  - core/validate/validate_test.go
  - core/merge/merge.go
  - core/merge/merge_test.go
  - internal/tui/signoff_test.go
  - internal/tui/render_test.go
  - .planning/REQUIREMENTS.md
  - .planning/ROADMAP.md
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
    - "A project that changed nothing carries no lane files: jaira init writes .jaira/lanes/ only when the default board differs from the built-in set"
    - "A default board set once per user decides which lanes a new board gets and which ticket Options start ticked"
    - "An agent can find the default board file, write one from a template, and have jaira lanes report what is wrong with it"
    - "precedence decides display order; after is a checked constraint, and a lane ordered before the lane that produces an input it requires is reported at load time"
    - "The ten built-in lanes render in exactly the order they render in today"
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
    - path: "core/lane/defaultboard.go"
      provides: "the per-user default board: parse, validate, and materialise .jaira/lanes/ only when it differs from the built-ins"
    - path: "internal/tui/defaultboard.go"
      provides: "the default board screen reached from the home screen: pick lanes, pre-tick options, 'e' to $EDITOR"
    - path: "core/ticket/body.go"
      provides: "one starting ticket body shared by the CLI and the TUI, with Options pre-ticked from the default board"
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
    - from: "internal/cli/tickets.go"
      to: "core/lane.Materialise"
      via: "jaira init writes .jaira/lanes/ only when the default board differs from the built-ins"
      pattern: "lane\\.Materialise"
    - from: "core/merge/merge.go"
      to: "core/lane.Set.Progress"
      via: "merge ranks lanes by progress, not by display order, so a parking lane never wins"
      pattern: "\\.Progress\\("
    - from: "internal/tui/home.go"
      to: "internal/tui/defaultboard.go"
      via: "the default board keybind on the home screen, the only per-user surface there is"
      pattern: "defaultBoard"
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
`jaira share`, a TUI lane settings screen that publishes and adopts, a per-user
default board that decides what a fresh board starts with, and an ordering model
where `precedence` decides position and `after` becomes a checked constraint.
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
carrying tasks 2 and 3 together, split it. Tasks 8, 9 and 10 close success
criteria 10, 11 and 12, which were added to ROADMAP.md after tasks 1-7 were
written; task 10 is the riskiest in the plan because it changes where every
column appears on every board.

Three quick tasks are landing in parallel on `core/ticket/frontmatter.go`,
`core/ticket/schema.go`, `core/ticket/checklist.go`, `internal/cli/tickets.go`,
`internal/cli/checklist.go`, `internal/tui/view.go`, `internal/tui/signoff.go`,
`core/board/*` and three files under `core/lane/builtin/`. Do not trust line
numbers quoted anywhere in this plan; locate code by identifier. Tasks 1-9
deliberately edit no file under `core/lane/builtin/` (see task 3). Task 10 edits
exactly one — `60-blocked.md` — and cannot avoid it; before starting task 10,
check `git log --oneline -- core/lane/builtin/` and rebase onto whatever the
quick tasks landed rather than resolving a conflict inside a lane definition.
</execution_notes>

<scope_discipline>
The project constraint is "is this smaller than paca?". Eight places where this
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
6. **The default board screen selects lanes; it does not reorder them** (task
   8). The design note says the screen "selects and orders lanes". Ordering was
   dropped because success criterion 12, settled the same day, makes
   `precedence` the single thing that decides position. A reorder control on
   the default board screen would be a second way to position a lane — exactly
   the duplicate mechanism criterion 12 exists to remove — and it would have to
   write `precedence` back into the lane files anyway. So the screen shows
   lanes in precedence order and `e` opens the file, which is the escape hatch
   the design note already chose for everything else about a lane. If the user
   wants a reorder control despite this, it is a follow-up, not a silent drop.
7. **The default board is a selection, not a board definition** (task 8). It
   holds a list of lane ids and a list of options to pre-tick. It does not hold
   prompts, tiers, contracts, colours or per-project overrides. Anything richer
   would make it a second lane format.
8. **No `jaira default-board` command tree** (task 9). Same rule as lanes: the
   file is the API. `jaira lanes path` says where the file is, `jaira lanes
   template --board` prints its shape, and `jaira lanes` validates it. Three
   additions to a surface being built anyway, instead of a fourth command.
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

Two more from task 10, which changes what `precedence` means:

- Display order stops following the `after:` anchor and follows `precedence`.
  The ten built-ins are unaffected — that is the task's regression test — but a
  user's custom lane with no `precedence:` used to park before the terminal lane
  and now derives one from its anchor. A custom lane that declared a
  `precedence` inconsistent with its `after` will move, which is the point:
  `precedence: 12` rendering before `precedence: 5` was observed, not theorised.
- `blocked` becomes `precedence: 65` with `off-pipeline: true`, and merge ranks
  by `Progress()` rather than raw precedence. A `backlog`/`blocked` status
  collision now resolves to `backlog` instead of `blocked`; the blocking fact
  lives in `blocked-by`, which merges as a union and is unaffected.
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
  <decided date="2026-08-13">

**D-01 — `<user>` in `.jaira/shared/<user>/`: the existing `identity()`, slugified.**
`JAIRA_USER` → `git config user.name` → `$USER` → `unknown`, lowercased and reduced
to `[a-z0-9-]`. Moves to a `core/identity` package so `internal/tui` can reach it.
Two teammates with the same configured name collide into one folder; `JAIRA_USER`
is the one-line fix and is documented as such.

**D-02 — drift: warn, and check it where the user is already looking.**
Editing a lane writes through to the global catalogue immediately, so it is
*other* projects that drift, and they drift silently. Rather than checking on
every command, the lane settings screen runs the check when it opens: it compares
the lanes this project uses against their catalogue copies and shows a small
warning on the ones that differ, with a refresh action to pull the catalogue
version in. Nothing syncs by itself — the direction is always the user's.

**D-03 — `.jaira/lanes/` is authoritative.**
It is the record of which lanes this project uses; without it nothing stores that
selection at all. So the loader reads it as the project's lane list rather than as
a set of overrides.

Consequence that must be handled in Task 2, not discovered later: the export that
populates the directory writes the full working set, built-ins included. A
directory holding one hand-written file would otherwise leave a board with one
lane. The loader warns when `.jaira/lanes/` exists but holds no lane declaring
`requires-specified`, because that board can never move a ticket into work.

**D-04 (new) — a brainstorm lane ships as a built-in optional step.**
Already landed ahead of this phase: `core/lane/builtin/05-brainstorm.md`,
`requires-option: brainstorm`, `output-produces: [goal]` so it cannot be left
until the brainstorm reached an intent. The gate rework that made it possible is
in `core/gate/gate.go` (`requires-specified`).

  </decided>
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

Add the gate-weakening warning. Settled 2026-08-13 in the design note: **nothing
is off limits, including the protections.** An override may drop
`requires-human-exit` from sign-off or `requires-nonmodel-signal` from done, and
the loader allows it — the user chose freedom over a lock. So the loader's job is
not to refuse; it is to make the loss impossible to miss. When an override
replaces a built-in that had `RequiresHumanExit`, `Terminal`, `RequiresOutcome`
or `RequiresNonModelSignal` set and the replacement does not, emit a **second,
separate warning line** naming each protection that went away — not a clause
appended to the ordinary "this overrides a built-in" line, which scrolls past
unread. Word it as a consequence rather than a fact: an agent accepting its own
work is what `requires-human-exit` prevents, and that is what the warning should
say happened. This is the one place where "powerful but never silent" has to mean
more than a line in a list.

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

<task type="auto" tdd="true">
  <name>Task 8: The default board — one per-user file that decides what a fresh board starts with</name>
  <files>core/lane/defaultboard.go, core/lane/defaultboard_test.go, core/lane/lane.go, core/ticket/body.go, core/ticket/body_test.go, internal/cli/tickets.go, internal/cli/init_test.go, internal/tui/defaultboard.go, internal/tui/defaultboard_test.go, internal/tui/home.go</files>
  <behavior>
Closes success criterion 10. Depends on task 2 (`ProjectLanesDir`) and task 6
(`lane.Export`, which this reuses to write copies).

In `core/lane/defaultboard_test.go`:
- No file at all: `LoadDefaultBoard()` returns a zero board with no error. Absent
  is the normal state, not a missing-configuration error.
- A file selecting `[backlog, todo, done]` parses to those three ids in that
  order, and `options: [brainstorm]` parses to that one option.
- `Differs(board, builtins)` is false when the selection names exactly the
  built-in ids and every selected lane resolves to a built-in. It is true when a
  lane is dropped, added, or resolves to a catalogue file that overrides a
  built-in.
- `Materialise(root, set, board)` writes nothing and returns no paths when
  `Differs` is false — assert `.jaira/lanes/` does not exist afterwards. This is
  the criterion's load-bearing half.
- When `Differs` is true it writes one `.md` per selected lane into
  `ProjectLanesDir(root)`, named `<id>.md`, and the bytes of a built-in copy
  parse back to the same lane.
- An unknown lane id in `lanes:` is a warning naming the id, not an error, and
  the remaining lanes still materialise. A shared default board that names a lane
  you have not adopted must not stop `jaira init`.

In `core/ticket/body_test.go`: `NewBody` emits the Options checklist with the
named options ticked and the rest unticked, and is byte-identical to today's
output when nothing is ticked — the existing shape is the regression baseline.

In `internal/cli/init_test.go`: `jaira init` in a temp root with no default board
leaves no `.jaira/lanes/`; with a default board that drops a lane, it creates one
file per selected lane. Then `jaira create` in that root produces a body whose
Options section has the board's options ticked.
  </behavior>
  <action>
`core/lane/defaultboard.go`. `DefaultBoardPath()` returns `$JAIRA_DEFAULT_BOARD`
when set, else `~/.jaira/default-board.md`. Give it an explicit env override
rather than deriving it from `UserLanesDir()`'s parent: tests need to point at a
temp file, and inferring a sibling path from an already-overridable directory is
the kind of cleverness that breaks the first time someone sets
`JAIRA_LANES_DIR` to something that is not `<something>/lanes`.

The file is markdown with frontmatter, parsed with `ticket.ParseDoc` and
`d.List` exactly as a lane file is — the file format is the API, and a second
format would be a second thing to teach an agent:

    ---
    lanes: [backlog, brainstorm, todo, in-progress, review, signoff, done, blocked]
    options: [brainstorm]
    ---

    # Default board
    <prose explaining what this file does>

`DefaultBoard{Lanes []string; Options []string; Path string; Warnings []string}`.
Absent file, absent `lanes:`, or an empty list all mean the built-ins; say so in
the type's doc comment, because "empty means everything" is a rule a reader will
otherwise get backwards.

`Materialise(root string, set *Set, b *DefaultBoard) ([]string, error)` writes
`.jaira/lanes/` **only when the selection differs from the built-in set**. This
is the whole point of the criterion: a repo whose owner changed nothing carries
no lane files, and an absent directory means the built-ins. Compare as a set of
ids plus the origin of each resolved lane, not as a list — reordering alone must
not trigger materialisation, since order is `precedence`'s business (task 10).

Materialising a built-in needs its original bytes, and a built-in has no file on
disk. Add `Bytes(l *Lane) ([]byte, error)` to `core/lane` reading from
`builtinFS` for built-ins and from `l.Source` otherwise, and route task 6's
`Export` through it if task 6 did not already need it. Copy bytes verbatim —
never parse and re-serialise, for the reason task 6 already gives.

Wire `jaira init` (`newInitCmd` in `internal/cli/tickets.go`) to call
`LoadDefaultBoard` then `Materialise`, and report what it wrote: either "using
the built-in lanes" or "wrote N lane files from your default board". Both
outputs, human and `--json`. Surface the board's warnings on stderr the way
`jaira lanes` already surfaces the loader's.

Ticket options. Add `Set.Options() []string` returning the distinct
`RequiresOption` values across the loaded lanes, in display order — the list of
options is derived from the installed lanes rather than hardcoded, which is what
makes a user-written optional lane appear in the checklist at all. Move
`newTicketBody` out of `internal/cli/tickets.go` into `core/ticket/body.go` as
`NewBody(title, dod string, options []BodyOption) string`, where `BodyOption` is
`{Name string; Ticked bool}`. `core/ticket` must not import `core/lane` (lane
imports ticket), so the caller resolves the options and passes them in. Keep
every comment from the current function: they explain why the options start
unticked and why the Plan heading holds no checkbox, and both facts survive this
change — the default board changes the default, not the reasoning.

`internal/tui/model.go:693` creates tickets with an empty body, so a
TUI-created ticket has no Options section at all today. Pass `ticket.NewBody(...)`
there too. This is a pre-existing gap, not one this task introduces, but leaving
it means "always brainstorm" would silently not apply to half the tickets — which
would make the setting a lie rather than a limitation.

The screen. `internal/tui/defaultboard.go`, following `internal/tui/browse.go`'s
shape: its own struct, a `key(string) (…, done bool)` method, a `render(width,
height int) string`, no bubbletea imports of its own. `Home` is its own
`tea.Model` and already holds a `*lane.Set`, so this hangs off `Home` and not off
the board `Model` — the default board is per-user, and the home screen is the
only per-user surface there is. Reach it with `d` from the home screen (`a`, `r`,
`q`, `j`, `k`, `enter` and `esc` are taken; confirm against the current footer
before committing to `d`).

The screen lists every loaded lane with a checkbox, then every option with a
checkbox. `space` toggles, `e` opens the selected lane's file in `$EDITOR` via
the existing `editorCommand()` in `internal/tui/external.go`, `s` saves, `esc`
backs out. It contains **no form for a lane's prompt, tier or contract** — that is
what `e` is for, and a second way to write a lane file is a second way for the
two to disagree. It does **not** reorder lanes; see scope discipline item 6, and
say so in the screen's own help line so the absence reads as a decision rather
than a missing feature.

Saving writes the file with the frontmatter fields in a fixed order and the prose
body preserved if the file already had one. This file is small, hand-written, and
not a ticket, so a plain rewrite is acceptable here — but keep any body the user
wrote, because that is where they will have left themselves a note.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin &amp;&amp; go test ./core/lane/... ./core/ticket/... ./internal/cli/... ./internal/tui/... 2>&amp;1 | tail -20</automated>
    <automated>export PATH=$PATH:$HOME/.local/go/bin &amp;&amp; go build -o /tmp/jaira-p5 ./cmd/jaira &amp;&amp; d=$(mktemp -d) &amp;&amp; (cd "$d" &amp;&amp; git init -q . &amp;&amp; /tmp/jaira-p5 init >/dev/null &amp;&amp; test ! -d .jaira/lanes &amp;&amp; echo "no default board: no lane files")</automated>
    <automated>export PATH=$PATH:$HOME/.local/go/bin &amp;&amp; go build -o /tmp/jaira-p5 ./cmd/jaira &amp;&amp; b=$(mktemp -d)/default-board.md &amp;&amp; printf -- '---\nlanes: [backlog, todo, done]\noptions: [brainstorm]\n---\n' > "$b" &amp;&amp; d=$(mktemp -d) &amp;&amp; (cd "$d" &amp;&amp; git init -q . &amp;&amp; JAIRA_DEFAULT_BOARD="$b" /tmp/jaira-p5 init >/dev/null &amp;&amp; test -f .jaira/lanes/backlog.md &amp;&amp; test -f .jaira/lanes/done.md &amp;&amp; test ! -f .jaira/lanes/review.md &amp;&amp; echo "differs: materialised the selection only")</automated>
    <human-check>Launch the launcher with no board argument, press `d`, untick a lane and tick the brainstorm option, save. Confirm `~/.jaira/default-board.md` now holds that selection. Press `e` on a lane and confirm it opens that lane's file in your editor and that the board redraws on exit. Then run `jaira init` in a fresh repo and confirm the lane files that appear match what you picked.</human-check>
  </verify>
  <done>`~/.jaira/default-board.md` exists as a documented format, is reachable from the home screen with `d`, and decides both which lanes `jaira init` materialises and which Options a new ticket starts with; a repo whose default board matches the built-ins has no `.jaira/lanes/` directory at all, proven by a test and by the init verify above; `NewBody` is the single starting body for both the CLI and the TUI.</done>
</task>

<task type="auto" tdd="true">
  <name>Task 9: The agent-facing surface for lanes and the default board</name>
  <files>internal/cli/tickets.go, internal/cli/lanes_test.go, core/lane/defaultboard.go, core/lane/defaultboard_test.go</files>
  <behavior>
Closes success criterion 11. Depends on task 3 (the `lanes` subcommand tree) and
task 8 (the default board format).

**What already works and must not be rebuilt** — verified directly against the
current binary, not assumed:
- Writing a `.md` file into `~/.jaira/lanes/` already makes the lane appear.
- `jaira lanes` already prints each lane's source path.
- A second file reusing an id is already reported as a warning and ignored.

What is missing is discoverability and a checker for the one file that has
neither: the default board. So this task adds three things and no new verbs.

In `internal/cli/lanes_test.go`:
- `jaira lanes path` names every source of lane definitions — the catalogue, the
  project directory, and the default board file — each labelled, with a marker on
  whichever is in force. `--json` carries them as named fields plus which is
  active. Assert the default board path appears in both modes, since task 3 only
  required the two directories.
- `jaira lanes path` outside a board still names the catalogue and the default
  board and says there is no project, so an agent can write a lane before there
  is a repo to write it for.
- `jaira lanes template --board` prints a default board file that
  `LoadDefaultBoard` accepts — assert by parsing the command's own stdout, the
  same way task 3 asserts the lane template. A template the parser has drifted
  from is worse than none.
- A default board naming a lane that is not installed makes `jaira lanes` print a
  warning naming the id, on stderr in human mode and in the `warnings` array
  under `--json`, and exit 0. A warning, not a failure: a teammate's default
  board naming their own lane must not break your `jaira lanes`.
- A default board naming an option no lane requires warns the same way, naming
  the option. This is the typo case, and it is silent today because nothing
  reads the file at all.
- A default board that is not parseable warns and is treated as absent, rather
  than taking the board down.
  </behavior>
  <action>
Extend `lanes path` from task 3 with a third entry: the default board file, its
path, and whether it exists. Under `--json` add `default_board` alongside the
directory fields. This is the one call an agent makes before writing anything, so
it must answer "where does each kind of file go" completely — a path command that
names two of the three places is a trap.

Add `--board` to `lanes template`, printing a commented default board skeleton to
stdout and nothing else, same contract as the lane template: it creates no file
and asks no questions, so `jaira lanes template --board > $(jaira lanes path
--json | jq -r .default_board)` is the whole workflow. The comments must state
the two rules that are not visible from the fields: an absent or empty `lanes:`
means the built-ins, and `options:` names entries in a ticket's Options checklist
which come from the installed lanes' `requires-option`, not from a fixed list.

Add `Validate(set *Set) []string` to `core/lane/defaultboard.go` — unknown lane
id, unknown option, unparseable file — returning warnings in the same string
shape the loader already uses, and have the bare `jaira lanes` command call it
and append the results to the warnings it already prints. One warning channel,
already surfaced in `jaira lanes`, `--json`, the TUI warnings block and every
command through `loadEnv`; adding a second reporting path for the same class of
problem would be the mistake here.

Do not add `jaira lanes new/edit/move/rm` and do not add a `default-board`
command tree. Both were explicitly rejected. The loop for an agent is: `jaira
lanes path` to find where, `jaira lanes template` (or `--board`) for the shape,
write the file with the tools it already has, `jaira lanes` to check the result.
That is a complete read-write-verify loop with exactly one write path, and it is
the same loop a human uses.

Document the loop in the `lanes` parent command's `Long` help, in four lines, in
that order. An agent reads `--help` before it reads a design note, and a surface
that is discoverable only by having been told about it is not discoverable.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin &amp;&amp; go test ./internal/cli/... ./core/lane/... 2>&amp;1 | tail -20</automated>
    <automated>export PATH=$PATH:$HOME/.local/go/bin &amp;&amp; go build -o /tmp/jaira-p5 ./cmd/jaira &amp;&amp; /tmp/jaira-p5 lanes path --json | grep -q '"default_board"' &amp;&amp; /tmp/jaira-p5 lanes template --board | grep -q '^lanes:' &amp;&amp; echo "path names the board file; template emits one"</automated>
    <automated>export PATH=$PATH:$HOME/.local/go/bin &amp;&amp; go build -o /tmp/jaira-p5 ./cmd/jaira &amp;&amp; b=$(mktemp -d)/default-board.md &amp;&amp; printf -- '---\nlanes: [backlog, nosuchlane]\noptions: [nosuchoption]\n---\n' > "$b" &amp;&amp; JAIRA_DEFAULT_BOARD="$b" /tmp/jaira-p5 lanes 2>&amp;1 >/dev/null | grep -q nosuchlane &amp;&amp; JAIRA_DEFAULT_BOARD="$b" /tmp/jaira-p5 lanes 2>&amp;1 >/dev/null | grep -q nosuchoption &amp;&amp; echo "default board is validated"</automated>
  </verify>
  <done>An agent with only a shell can create a lane, change one, and write a default board: `jaira lanes path` names all three locations and which is in force, `jaira lanes template` and `jaira lanes template --board` print shapes the parsers accept (proven by parsing the commands' own output), and `jaira lanes` reports a bad lane id or a bad option in the default board as a named warning while still exiting 0.</done>
</task>

<task type="auto" tdd="true">
  <name>Task 10: `precedence` decides the order, `after` becomes a checked constraint, and contracts are checked against the order</name>
  <files>core/lane/lane.go, core/lane/lane_test.go, core/lane/builtin/60-blocked.md, core/merge/merge.go, core/merge/merge_test.go, internal/cli/flow.go, internal/cli/tickets.go</files>
  <behavior>
Closes success criterion 12. This is the riskiest task in the plan: it changes
where every column appears on every board, so its first test is a regression
test, not a feature test.

**The regression test, written first and never deleted:** load the built-ins with
an empty catalogue and assert `set.IDs()` equals exactly

    backlog brainstorm todo pre-process in-progress human review signoff done blocked

That is today's order, by id, and it must be identical after this change. Take
the baseline by running the current binary before touching anything, not from
this plan — if the quick tasks landed another built-in lane in the meantime, the
baseline is whatever the binary prints today, and this plan's list is stale.

Then, in `core/lane/lane_test.go`:
- Order follows `precedence` ascending. A custom lane with `precedence: 22`
  lands between `todo` (20) and `pre-process` (25) regardless of filename and
  regardless of what it anchors to.
- Ties are stable and deterministic: two lanes at the same precedence keep
  built-ins first, then custom lanes by sorted source path.
- A lane with `after: X` and a `precedence` at or below X's is **reported**, not
  moved. Assert the warning names both lanes and says which way to fix it.
- A lane with `after: X` and a higher precedence than X produces no warning.
- `after` naming a lane that is not installed warns that the constraint cannot be
  checked, and the lane still orders by its precedence. **This replaces task 5's
  test asserting it lands before the terminal lane** — update that test rather
  than leaving two tests asserting opposite things.
- Cyclic `after` between two lanes is reported as an unsatisfiable constraint and
  both lanes still appear, each at its own precedence. No lane ever disappears
  because a constraint is wrong.
- A custom lane with `after:` and no `precedence` derives one just above its
  anchor's instead of defaulting to 0 and jumping to the front of the board.
- A custom lane with neither still parks just before the terminal lane, as today.
- Contract check: a lane whose `input-requires` names a field that only a
  later-ordered lane's `output-produces` supplies is reported at load time,
  naming the field, the requiring lane and the producing lane.
- Base ticket fields do not trigger it: a lane requiring `goal`, `title`,
  `context`, `definition-of-done`, `assignee` or `diff` is silent even though
  `brainstorm` also declares `output-produces: [goal]`.
- A field no installed lane produces is silent — it is a ticket field this
  package does not model, not a misordering.
- **The ten built-ins produce zero warnings of either kind.** A shipped default
  that warns on every load trains the user to ignore warnings.

In `core/merge/merge_test.go`: a status collision between `done` and `blocked`
still resolves to `done`, and one between `in-progress` and `blocked` still
resolves to `in-progress`.
  </behavior>
  <action>
Read this whole action before writing code. There is a conflict in it that has to
be resolved deliberately.

**The conflict.** `Lane.Precedence`'s doc comment says today, in as many words:
"It is deliberately separate from display order: Blocked appears last on the
board but must not outrank active work in a merge." Criterion 12 says precedence
decides the order. Both cannot be true for `blocked`, which has `precedence: 10`
and renders last. Sorting by precedence today moves Blocked to third — between
Brainstorm and Todo — which fails the regression test above. And `blocked`
declares `after: done` while sorting before it, so under the new rule the shipped
board would report a constraint violation against itself on every load.

**The resolution to implement:** split the two jobs by name, not by adding a
second position field.

1. `precedence` becomes display order. Change `core/lane/builtin/60-blocked.md`
   to `precedence: 65` and add `off-pipeline: true`. Blocked then sorts last,
   satisfies its own `after: done`, and the regression test passes with exactly
   one built-in file edited.
2. Parse `off-pipeline: true` onto `Lane.OffPipeline`: a parking lane, not a
   step. Nothing is expected to progress *through* it.
3. Add `Set.Progress(id string) int` returning `-1` for an off-pipeline lane and
   for an unknown one, and the lane's precedence otherwise. Point
   `mergeStatus` in `core/merge/merge.go` and `sortByProgress` in
   `internal/cli/flow.go` at `Progress` instead of `Precedence`. Leave
   `Set.Precedence` alone: it means display position now, and callers that want
   position should keep getting it.

That keeps merge's invariant exactly as Phase 4 built it — a parking lane never
counts as further along, so a side moving a ticket to Blocked never overwrites a
side moving it to Done — while making precedence the single answer to "where does
this column go".

**One behaviour change to state in the summary and the release note:** a status
collision between `backlog` and `blocked` used to resolve to `blocked` (10 beat
0) and now resolves to `backlog` (-1 loses to 0). The blocking information is not
lost: `blocked-by` is a list and merges as a union, which is where the fact
actually lives. Add a test that pins this, so the flip is recorded rather than
discovered.

**If this resolution is rejected**, stop and say so rather than improvising. The
two alternatives both cost more: giving Blocked a precedence above Done without
`off-pipeline` silently reverts Done to Blocked on a merge, and keeping a
separate display field reintroduces the second positioning mechanism criterion 12
exists to delete.

Now the ordering itself. Rewrite `order()`:

- Sort every lane by `(Precedence, builtin-first, source-path)` with a stable
  sort over a deterministic input order. Delete the fixed-point anchor-insertion
  loop: with precedence deciding position there is nothing to place iteratively,
  and the loop's whole purpose was to be a positioning mechanism.
- Derive a precedence for lanes that declare none: anchor's precedence + 1 if
  `after` names an installed lane, otherwise the terminal lane's precedence
  minus 1, which preserves today's "park before the terminal lane" behaviour for
  a lane file that says nothing. Record the derived value on the lane so
  `jaira lanes` shows what it actually used rather than a misleading 0.
- Check `after` as a constraint and report violations. Report; do not repair.
  Repairing would make `after` a positioning mechanism again by the back door,
  and the user cannot fix what the tool silently corrected.
- Keep the cycle detection that exists today, rewritten as a constraint check:
  a cycle in `after` is a set of constraints that cannot all hold.

Then the contract check, which is the part nobody has today. `InputRequires` and
`OutputProduces` are parsed and never compared against anything (`core/lane/lane.go`,
in `parse`). Add `checkContracts(ordered []*Lane) []string`: for each lane, for
each `input-requires` field that is not a base ticket field, find every installed
lane declaring it in `output-produces`; if there is at least one producer and
every producer sorts at or after the consumer, warn — naming the field, the
consumer and the producer. Run it next to the constraint check, in the same
warnings slice.

The base ticket fields are `title`, `goal`, `context`, `definition-of-done`,
`assignee`, `diff` and `commits`: supplied by the ticket, not by any lane. The
authority on what a field means is `fieldFilled` in `core/gate/gate.go`, but
`core/gate` imports `core/lane`, so this list cannot live there without a cycle.
Declare it in `core/lane` and add a test **in `core/gate`** — which may import
lane — asserting every name in `lane.BaseTicketFields` is a case `fieldFilled`
handles. That is what stops the two lists from drifting apart, and it is three
lines.

Finally, update the surfaces that describe ordering, or the new model is
invisible: `lanes template` (task 3) must document `precedence` as the order,
`after` as a checked constraint, and `off-pipeline` as "a parking lane, never
counted as progress"; `lanes --json` gains `off_pipeline`; and `Lane.Precedence`'s
doc comment must stop saying precedence is separate from display order, because
after this task it is display order.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin &amp;&amp; go build ./... &amp;&amp; go test ./core/lane/... ./core/merge/... ./core/gate/... ./internal/... 2>&amp;1 | tail -30</automated>
    <automated>export PATH=$PATH:$HOME/.local/go/bin &amp;&amp; go build -o /tmp/jaira-p5 ./cmd/jaira &amp;&amp; JAIRA_LANES_DIR=$(mktemp -d) /tmp/jaira-p5 lanes | tail -n +2 | awk '{print $1}' | paste -sd' ' - | grep -qx "backlog brainstorm todo pre-process in-progress human review signoff done blocked" &amp;&amp; echo "built-in order unchanged"</automated>
    <automated>export PATH=$PATH:$HOME/.local/go/bin &amp;&amp; go build -o /tmp/jaira-p5 ./cmd/jaira &amp;&amp; test -z "$(JAIRA_LANES_DIR=$(mktemp -d) /tmp/jaira-p5 lanes 2>&amp;1 >/dev/null)" &amp;&amp; echo "shipped board warns about nothing"</automated>
    <automated>export PATH=$PATH:$HOME/.local/go/bin &amp;&amp; go build -o /tmp/jaira-p5 ./cmd/jaira &amp;&amp; d=$(mktemp -d) &amp;&amp; printf -- '---\nid: early-review\nname: Early\nafter: in-progress\nprecedence: 1\nagentic: false\ninput-requires: [plan]\n---\n' > "$d/early.md" &amp;&amp; JAIRA_LANES_DIR="$d" /tmp/jaira-p5 lanes 2>&amp;1 >/dev/null | grep -q "in-progress" &amp;&amp; JAIRA_LANES_DIR="$d" /tmp/jaira-p5 lanes 2>&amp;1 >/dev/null | grep -q "plan" &amp;&amp; echo "constraint and contract both reported"</automated>
  </verify>
  <done>The ten built-in lanes come out of this change in exactly the order they went in, asserted by id in a test and by the second verify above; `precedence` is the only thing that decides display order; an `after` that the precedence violates is reported and not silently repaired; a lane requiring a field only a later lane produces is reported at load time next to the cycle check; base ticket fields never trigger it, guarded by a test in `core/gate` that keeps the base-field list from drifting from `fieldFilled`; and the shipped board produces no warnings at all.</done>
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
| T-5-07 | Tampering | `lane.Materialise` / default board | mitigate | Lane ids from `~/.jaira/default-board.md` become filenames under `.jaira/lanes/`. Materialisation names files from the resolved lane's already-`validID`-constrained `ID`, never from the string in the default board, and an id matching no installed lane is a warning rather than a file. Tested in `core/lane/defaultboard_test.go`. |
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
| 10 | A per-user default board decides a new board's lanes and pre-ticked options; absent `.jaira/lanes/` means the built-ins | 8 |
| 11 | An agent can create a lane, change one and edit the default board without the TUI, including validation of the default board | 9 (path names all three, `--board` template, default board validated), 3 (the `lanes` surface it extends) |
| 12 | `precedence` decides the order; `after` is a checked constraint; a lane ordered before its input's producer is reported at load time | 10 |

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
| An override may drop a protection; the loss gets its own warning (settled 2026-08-13) | 2 |
| The default board at `~/.jaira/default-board.md`, reached from the home screen (settled 2026-08-13) | 8 |
| Directory absent means the built-ins; `init` materialises only on a difference | 8 |
| `e` opens the lane file in `$EDITOR`; no form for prompt, tier or contract | 8 |
| `precedence` decides order, `after` is a constraint (settled 2026-08-13) | 10 |
| Contracts checked against the order at load time | 10 |

## Surfaced, not in any source artifact

| Item | Task | Why it is here |
|---|---|---|
| Additive vs authoritative project lanes | 1 (D-03) | criterion 1 is only delivered under one of the two readings; not settled anywhere |
| `lane.Load()` must take a root | 2 | 16 call sites; unavoidable consequence of a per-project lane directory |
| REQUIREMENTS.md LANE-06 contradicts the design | 2 | shipping code against a live requirement is a defect either way |
| `store.go`'s "only committable content" comment becomes false | 4 | the invariant this phase breaks should not stay written down as true |
| A gate-weakening override is a privilege escalation | 2, threat model | overriding built-ins was approved; overriding `requires-human-exit` specifically was not discussed |
| `blocked` cannot satisfy criterion 12 and today's board at once | 10 | `precedence: 10` renders last today only because display order ignores precedence; making precedence the order moves Blocked to third and makes the shipped board violate its own `after: done`. Resolved with an `off-pipeline` flag and a merge-side `Progress()`; the alternatives are named in the task |
| Merge and display disagree about what precedence means | 10 | Phase 4's merge invariant ("never revert forward progress") and criterion 12 both claim the field; splitting the accessor is the only way to keep both |
| The default board screen "selects and orders"; ordering was dropped | 8, scope discipline 6 | a reorder control would be the second positioning mechanism criterion 12 removes — stated rather than silently omitted |
| The TUI creates tickets with an empty body, so Options never appear there | 8 | pre-existing, but it would make the default board's option pre-ticking a half-truth |
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
7. `jaira init` in a second throwaway repo with no default board — confirm there
   is no `.jaira/lanes/` directory at all.
8. From the launcher, press `d`, drop a lane and tick `brainstorm`, save. `jaira
   init` in a third repo — confirm `.jaira/lanes/` now holds exactly the picked
   lanes, and `jaira create "x"` produces a body with `- [x] brainstorm`.
9. `jaira lanes path --json` — confirm it names the catalogue, the project
   directory and the default board file, and says which is in force.
10. Set `precedence: 1` on a lane that declares `after: in-progress` and requires
   `plan`; `jaira lanes` reports both the violated constraint and the unsatisfiable
   contract, and still draws every column.
11. `jaira lanes` on a clean install prints the ten built-ins in their shipped
   order and no warnings at all.
</verification>

<success_criteria>
- `go build ./...` and `go test ./...` pass.
- `core/lane` has tests where it had none.
- All twelve ROADMAP Phase 5 success criteria have a named test or a recorded
  human check.
- REQUIREMENTS.md LANE-06 no longer contradicts the shipped behaviour.
- No lane, on any path, is loaded onto a board because a teammate committed it.
- `jaira share` on a real repo leaves `.jaira/lanes/` unstaged.
- The ten built-in lanes render in exactly the order they render in today,
  asserted by id in a test that runs before and after the ordering change.
- A clean install produces no lane warnings on any command.
- `jaira init` writes no lane files for a project whose default board matches the
  built-ins.
</success_criteria>

<output>
Create `.planning/phase-5-custom-and-portable-lanes/05-01-SUMMARY.md` when done.
Record in it: the answers chosen for D-01, D-02 and D-03; the two behaviour
changes from the migration section, for the release note; whether task 5's tests
passed on first run or found real gaps; and the two behaviour changes task 10
introduces — `blocked` moving from `precedence: 10` to `65` with `off-pipeline:
true`, and a `backlog`/`blocked` status collision now resolving to `backlog`.
</output>
