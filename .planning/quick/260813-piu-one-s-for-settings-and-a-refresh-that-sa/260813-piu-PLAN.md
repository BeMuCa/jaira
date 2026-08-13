---
phase: quick/260813-piu
plan: 01
type: execute
wave: 1
depends_on: []
files_modified:
  - internal/tui/model.go
  - internal/tui/settings.go
  - internal/tui/settings_test.go
  - internal/tui/lanes.go
  - internal/tui/lanes_test.go
  - internal/tui/view.go
  - internal/tui/home.go
  - README.md
  - docs/COMMANDS.md
autonomous: true
requirements: [QUICK-260813-piu]

must_haves:
  truths:
    - "S on the board opens settings, and both the lane screen and the default board are reachable from there"
    - "The default board is reachable from inside a board, not only from the launcher"
    - "L no longer opens the lane screen: one door, not two"
    - "The lane screen's footer says refresh, not refresh drift"
    - "The shared-lane section is visible even when nobody has published one, and says how it fills"
  artifacts:
    - path: "internal/tui/settings.go"
      provides: "the settings screen: a menu over the lane screen and the default board"
  key_links:
    - from: "internal/tui/model.go"
      to: "modeSettings"
      via: "S on the board"
      pattern: "modeSettings"
---

<objective>
Three things the user hit while looking for settings.

`L` opened the lane screen and nothing else, so "settings" had no single door.
The default board sat on the launcher's `d` (`internal/tui/home.go:227`) and was
unreachable from inside a board — which is why it was never found. The lane
screen's footer said "refresh drift", which names an implementation detail rather
than the action. And the shared-lane list, with its `tab` and `a` keys, is hidden
whenever nobody has published a lane (`internal/tui/lanes.go:327-330`), so a user
with an empty `.jaira/shared/` concludes the feature does not exist.

Output: one `S`, a settings screen holding both existing screens, and two smaller
corrections to what the lane screen says.
</objective>

<context>
@CLAUDE.md
@internal/tui/lanes.go
@internal/tui/home.go
@internal/tui/model.go

<interfaces>
`L` is bound at internal/tui/model.go:612-614 and sets modeLanes.
The launcher's `d` is at internal/tui/home.go:227, building a defaultBoardScreen
via lane.LoadDefaultBoard(). Its footer is at home.go:330.
The lane screen's footer strings are at internal/tui/lanes.go:327-330.
`S` and `s` are both unbound on the board — verified.

Go is not on the default PATH:
  export PATH=$PATH:$HOME/.local/go/bin
</interfaces>
</context>

<tasks>

<task type="auto">
  <name>Task 1: One S for settings</name>
  <files>internal/tui/model.go, internal/tui/settings.go, internal/tui/settings_test.go, internal/tui/view.go, internal/tui/home.go</files>
  <action>
Add `modeSettings` and a `settings.go` holding a small menu screen. `S` on the
board opens it. Remove the `L` binding — two doors into the same room, one of
which contains the other, is how a UI accretes.

The menu has two entries, each opening the screen that already exists:
- **Lanes** — the existing lane screen (modeLanes, `internal/tui/lanes.go`)
- **Default board** — the existing default board screen
  (`newDefaultBoardScreen`, currently only reachable from the launcher)

Give each entry a one-line description saying what it is for, since "Lanes" alone
does not tell a first-time reader that this is where a lane is taken from the
catalogue into the project.

Keys: `jk`/arrows to move, `enter` to open, `esc`/`q` back to the board. Closing
an inner screen returns to the settings menu, not to the board — otherwise
changing two settings means opening the menu twice.

The launcher keeps its own `d`: it is the only screen available when no board is
open, so removing it would make the default board unreachable from there.

Tests in `internal/tui/settings_test.go`, in the style of the existing TUI tests:
- `S` on the board opens modeSettings; `L` no longer opens the lane screen
- `enter` on the Lanes entry reaches the lane screen
- `enter` on the Default board entry reaches the default board screen
- `esc` from an inner screen returns to the settings menu, and `esc` again to the board
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./internal/tui/... -count=1 2>&1 | tail -5</automated>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go test ./... -count=1 2>&1 | grep -c FAIL</automated>
  </verify>
  <done>S opens a settings menu holding both screens; L is gone; esc steps back one level at a time; every listed test passes.</done>
</task>

<task type="auto">
  <name>Task 2: Say refresh, and show the shared list even when it is empty</name>
  <files>internal/tui/lanes.go, internal/tui/lanes_test.go</files>
  <action>
Two corrections to the lane screen.

**The footer** says "R refresh drift". Drift is the reason the key exists, not
what pressing it does. Make it read `R refresh` in both footer variants
(lanes.go:327-330).

**The shared list** is omitted entirely when `len(ls.shared) == 0`, which is the
common case on a board where nobody has published anything — so `tab` and `a`
never appear and the whole sharing feature reads as absent. Render the section
regardless, with an empty state that says how it fills: a teammate publishes a
lane with `p` and commits `.jaira/shared/`. Keep `tab` and `a` in the footer only
when there is something to switch to or adopt — a key that does nothing is worse
than an absent one — but the section itself, and its one-line explanation, are
always visible.

Tests: the rendered screen names the shared section and its empty-state line when
no shared lanes exist; the footer says "refresh" and never "refresh drift".
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go test ./internal/tui/... -count=1 2>&1 | tail -3</automated>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && grep -c "refresh drift" internal/tui/lanes.go</automated>
  </verify>
  <done>The footer says refresh; the shared section and its explanation are visible on a board with no shared lanes; grep for "refresh drift" returns 0.</done>
</task>

<task type="auto">
  <name>Task 3: Update the key references</name>
  <files>README.md, docs/COMMANDS.md, internal/tui/view.go</files>
  <action>
Both key blocks and the `?` help screen list the board's keys. Replace `L` with
`S`, and say what settings holds — lanes and the default board — so a reader
knows where to look for either. Keep the blocks in their existing terse shape.
  </action>
  <verify>
    <automated>export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go test ./internal/tui/... -count=1 2>&1 | tail -3</automated>
    <automated>grep -c "settings" README.md docs/COMMANDS.md</automated>
  </verify>
  <done>Both key tables and the help screen name S and say what it contains; no reference to L as the lane screen remains.</done>
</task>

</tasks>

<verification>
export PATH=$PATH:$HOME/.local/go/bin && go build ./... && go vet ./... && go test ./... -count=1

Human check, unverifiable here: open a board, press S, confirm both entries open
and that esc steps back one level at a time.
</verification>

<success_criteria>
- S on the board opens settings; both screens are reachable from it; L is gone.
- The default board can be reached from inside a board.
- The lane screen says "refresh" and shows the shared section even when empty.
- `go build ./... && go vet ./... && go test ./... -count=1` passes.
</success_criteria>

<output>
Create `.planning/quick/260813-piu-one-s-for-settings-and-a-refresh-that-sa/260813-piu-SUMMARY.md` when done
</output>
