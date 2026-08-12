# Home screen and multi-project — specified 2026-08-12, not built

Decided in conversation, recorded here so it survives a context clear.

- **Icon.** `icon/jAIra.png`, 311×286. It has **no alpha channel** — 88,946
  pixels, all opaque; the dark ground is a baked-in gradient. "No background"
  therefore means keying out by luminance, not honouring transparency.
  `go run ./scripts/iconpreview` renders six candidate styles; **a style has not
  been chosen yet**.
- **A terminal tab icon is not achievable.** `OSC 0/1/2` set title text only.
  In-terminal images need Sixel / Kitty / iTerm2 protocols, none universal, and
  Windows Terminal under WSL2 is partial at best. The icon can only be drawn in
  the viewport.
- **Adding a project: all three ways.** A directory browser in the TUI, a
  `jaira projects add <path>` command, and auto-discovery that scans **at most
  two levels deep** below a chosen root — the user's code all lives under one
  `code/` directory with boards two levels down, and an unbounded scan of a home
  directory is slow enough to be a bug.
- **Multiple projects open at once**, not just switching between them.
- **Active agent sessions collected across projects** and marked in the project
  list, so one board shows which projects have agents working and which are
  waiting on a human. The `checkpoint` / `sessions` commands and
  `~/.jaira/state/<worktree>/` are the existing foundation.
- **Claude should be able to launch jaira** with a given project open — implies
  something like `jaira --project <path>` or `jaira board -C <path>`.
- **Open question, unanswered:** whether bare `jaira` shows the home screen
  always, or only when the current directory is not itself a board. The board is
  glanced at constantly, so a mandatory keypress on the common path has a real
  cost.

# Next steps — open tasks

Ordered. Each has the command to start with and how you know it is done.
Background in `HANDOFF.md`; that file also lists decisions not to re-open.

Always first:

```bash
export PATH=$HOME/.local/go/bin:$PATH
```

---

## 1. Install the skill globally — blocks everything else

**Why first:** the skill only exists at `/home/berk/git/jAIra/.claude/skills/jaira/`.
Outside this repo Claude does not know jaira exists, so it will not use the board
in requirementsgenie or anywhere else. Every other task assumes this is done.

**Waiting on:** the user. Asked, not yet answered.

```bash
mkdir -p ~/.claude/skills
cp -r /home/berk/git/jAIra/.claude/skills/jaira ~/.claude/skills/jaira
```

**Done when:** `ls ~/.claude/skills/jaira/SKILL.md` exists, and a Claude session in
an unrelated repo containing `.jaira/` reaches for `jaira` commands unprompted.

**Note:** a copy goes stale. If the skill changes, re-copy — or symlink instead
(`ln -s /home/berk/git/jAIra/.claude/skills/jaira ~/.claude/skills/jaira`) and
accept that it breaks if the repo moves.

---

## 2. Two TUI layout overflow bugs

**Status:** an agent reproduced these twice at 20×20, including with an ASCII-only
fixture, so they are not a wide-character artifact. **Never reproduced here.** Its
findings report stayed in its own context; only its audit arrived.

Either resume the agent:

> resume agent `ae2c189a2446fce69` and ask for the two layout findings with repro
> and verbatim output

or re-derive, which is quick — add to `internal/tui/render_test.go`:

```go
func TestNoLineOverflowsAtAnySize(t *testing.T) {
    for _, size := range [][2]int{{20,5},{20,20},{30,10},{40,20},{80,24},{120,40},{200,60}} {
        m := newTestModel(t, size[0], size[1])
        for _, mode := range []mode{modeBoard, modeDetail, modeHelp, modeProjects, modeMessage} {
            m.mode = mode
            for _, line := range strings.Split(stripANSI(m.render()), "\n") {
                if lipgloss.Width(line) > size[0] {
                    t.Errorf("%dx%d mode %v: line is %d cols: %q",
                        size[0], size[1], mode, lipgloss.Width(line), line)
                }
            }
        }
    }
}
```

Measure with `lipgloss.Width`, never `len()` — display width is the thing that
wraps a terminal.

**Also never stressed for overflow:** the diff view, the projects switcher, and the
message modal. The fuzz test only proved they do not panic.

**Done when:** that test passes at all seven sizes across every mode.

---

## 3. Migrate the 44 requirementsgenie tickets

**Approved by the user.** Script at `scripts/migrate-tickets.py` — written by the
sync agent, since read, parameterized (it had hardcoded `/tmp` paths and took no
arguments), and dry-run against copies of all 44. Results below are verified.

Target: `~/git/requirementsgenie-feature-requirements-coverage-elicitation/tickets/`
— 44 German markdown files on branch `feature/requirements-coverage-elicitation`.

What jaira already handles natively (verified on a real ticket):
- the title as an H1 in the body — falls back to it when `title:` is absent
- the `## Definition of Done` checklist — parsed, and ticking all boxes satisfies
  the Done gate
- Jira frontmatter (`jira`, `type`, `component`, `priority`, `labels`, `epic_link`)
  survives writes byte-for-byte

What the migration must add:
- `id:` — a ULID per ticket
- `status:` — a starting lane
- rename each file to `<ulid>-<slug>.md`

**Already verified on copies** (`/tmp/migcheck3`, may have been reaped — re-run to
reproduce):

- 44/44 byte-identical to the originals apart from the two inserted lines
- all Jira fields survive: `jira` `type` `component` `priority` `labels`
  `epic_link` 44/44, `assignee` 43/44 (one original genuinely lacks it)
- umlauts intact in all 44; jaira reads all 44 and shows a title for every one
- ULIDs derive from each file's mtime, so `ls` reflects real authoring order
  (Aug 6 → Aug 11) rather than batch order

**Invocation** — it takes two arguments and never writes to the source:

```bash
D=~/git/requirementsgenie-feature-requirements-coverage-elicitation
mkdir -p /tmp/mig/src
cp -p $D/tickets/*.md /tmp/mig/src/          # -p matters: plain cp resets mtimes,
                                             # which collapses all the ULIDs together
python3 scripts/migrate-tickets.py /tmp/mig/src /tmp/mig/out
```

**To apply for real**, in the requirementsgenie repo (ask first — `init` writes a
`.gitignore` entry):

```bash
cd $D && jaira init
python3 ~/git/jAIra/scripts/migrate-tickets.py tickets .jaira/tickets
jaira list                                   # expect 44
# then decide what happens to the old tickets/ dir — the script does not remove it
```

**Open question for the user:** after migrating, does the original `tickets/`
directory stay (as the Jira export staging area) or go? The script copies rather
than moves, so both exist until someone decides.

**Done when:** 44 in `jaira list` with titles, and a fully-ticked ticket reaches
`done` without `--signal` while a partially-ticked one does not.

**Note on slugs:** filenames keep German characters (`klären-ob-die-knöpfe…`),
because jaira's own `Slug()` treats any unicode letter as a word character. Consistent
with tickets jaira creates itself; change both together if you want ASCII-only names.

**Do not** run `jaira init` inside the real requirementsgenie repo without asking —
it writes a `.gitignore` entry.

---

## 4. Collect the agent findings that never arrived

Three of four agents ended with an audit instead of their report. Titles only:

**Store/merge** — resume `adb7f3391f65e3fac`:
- cherry-pick reverts a list deletion
- octopus merge fails
- an empty `conflict-theirs-*` key is left behind after `resolve`

**Sync/share** — resume `a6a66be715bf53031`: findings beyond the two already fixed.
Areas covered were `.gitignore` variants (`**/.jaira/`, `!.jaira/keep`, missing
trailing newline), `share` in a non-git directory, and tasks↔sync round-trip
convergence.

Reproduce each before acting. Two of the four already-fixed bugs came with a wrong
diagnosis attached, so treat reports as leads.

**Done when:** each is reproduced or shown not to reproduce, and either fixed or
recorded here with a reason for deferring.

---

## 5. TUI-11 — the one unmet v1 requirement

*"If a ticket is changed elsewhere while the user is editing it, the conflict is
surfaced rather than silently overwritten."*

Deliberately deferred: the board has no in-place field editor, so there is no edit
buffer to clobber. It only becomes real alongside an editor. The mechanism it would
use already exists — `updated-at` staleness detection, the same field the merge
driver keys on.

**Done when:** either an editor exists and surfaces the conflict, or the
requirement is formally moved to v2 in `REQUIREMENTS.md`.

---

## 6. Smaller things

- **Run the board under a real TTY.** Never done — no TTY in this environment. The
  11 render tests call `View()` directly; nobody has watched it run. `jaira` in a
  terminal is the only way to close this.
- **Verify fsnotify live refresh for real.** The watcher is wired in but was only
  exercised through the 2-second timer path. Touch a ticket file while the board is
  open and confirm it repaints.
- **Test the `sync-tasks.sh` hook against a live session.** The CLI is verified; the
  shell wrapper that digs the task list out of a real hook envelope is not, because
  the envelope shape is unconfirmed. It tries several `jq` paths and falls back to
  handing the whole payload to jaira's tolerant parser.
- **Add `~/.local/bin` and `~/.local/go/bin` to `~/.bashrc`.** Not done; the user
  hit "command go not found" twice because of it.
- **Consider dogfooding.** jAIra has no `.jaira/` of its own. Tracking this list as
  real tickets would exercise the tool and keep the list in git — but persisting it
  across a context clear needs `jaira share`, which commits `.jaira/` into this
  repo. That is a structural choice for the user to make, not one to assume.

---

## Guardrails

- Set `JAIRA_HOME` and `JAIRA_LANES_DIR` under `/tmp` in every test invocation, or
  runs write into the real `~/.jaira`.
- `$?` after a pipe reports the last command in the pipe. Measure exit codes with
  `cmd >/dev/null 2>&1; echo $?` — several early readings were wrong because of this.
- Never edit files under `.jaira/tickets/` directly; the CLI is the write path.
