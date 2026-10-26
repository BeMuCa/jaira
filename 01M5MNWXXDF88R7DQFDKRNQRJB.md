---
id: 01M5MNWXXDF88R7DQFDKRNQRJB
title: "Die Warnung 'Board von aelterer Version eingerichtet' kommt einmal, nicht bei jedem Aufruf"
status: human
ready: true
creator: BeMuCa
goal: "Wer jaira update noch nicht laufen liess, hoert es einmal pro Board und Version - nicht vor jeder Ausgabe, wo Agenten es mitparsen und Menschen es ueberlesen."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkt 11): 'jaira: this board was set up by an older version of jaira — run jaira update' steht vor jeder Ausgabe, auch bei --json auf stderr. Agenten filtern es weg (grep -v), Menschen sehen es nicht mehr.
  Es gibt schon einen Versionsstempel pro Board (core/release/release.go: Stamped/Stamp). Der Hinweis kann sich daran merken, dass er fuer diese Version gezeigt wurde.
  Verwandt: T8A8KM (Release-Pille), PWAF22 (update bringt das Setup nach).
definition-of-done: Die Warnung erscheint einmal pro Board und Binary-Version (Merker neben dem Versionsstempel) und danach nur noch bei jaira update selbst oder --verbose; ein Test deckt den zweiten Aufruf ab
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:26:01Z
updated-at: 2026-10-26T20:26:19Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: ""
claimed-at: ""
outcome-what: "nudgeIfStale (internal/cli/update.go) prints the 'board was set up by an older version' hint once per board and binary version: a 'nudged' marker file beside the version stamp (core/release Nudged/MarkNudged) remembers which version it was shown for. jaira update and the stamp are unchanged."
outcome-why: "The hint appeared on stderr before every command, even with --json, so agents filtered it away and people stopped reading it."
outcome-resolves: "DoD 1: the hint no longer repeats on every call; a new binary version shows it once again; afterwards jaira update's own report (previous version + what changed) is where staleness surfaces. Covered by TestNudgeIfStaleShowsOncePerBoardAndVersion and TestUpdateAfterShownHintStillReportsPreviousAndStamps. No --verbose flag exists, so none was added."
review-summary: |-
  critique 1: README.md:97 still says a newer binary on an older board 'says so on every command' - false after this change; reword to once per board and jaira version, on the first command after the upgrade. Optional: reflow the nudgeIfStale doc comment (internal/cli/update.go:20-26, line 24 is 88 chars); make core/release/notes/RNQRJB.md tell the reader to act on the hint because it will not repeat, rather than 'run jaira update to see the previous version' (update rewrites CLAUDE.md/AGENTS.md, it is not a status query).
  critique 2: round-1 findings handled (README.md:97-98, update.go comment reflow, release note). Open: callers that run jaira with --json and discard stderr use up the one hint unseen - the jaira-cockpit poll (core/mod/builtin/jaira-cockpit/hooks/register.tsx:281-285, every 60s and after each turn) and the dispatcher scripts run-lane.sh/reap.sh. Reproduced on a scratch board: a 'list --json 2>/dev/null' wrote nudged, the next visible 'list' printed nothing. Fix needs a design choice (cockpit forwards the line vs. --json neither prints nor consumes it) - sent to human.
round: 2
question: "Kritik 2: Aufrufer mit --json, die stderr wegwerfen (Cockpit-Poll alle 60s: core/mod/builtin/jaira-cockpit/hooks/register.tsx:281; Dispatcher run-lane.sh, reap.sh), verbrauchen den einen Hinweis - niemand sieht ihn (auf Scratch-Board nachgestellt). Welcher Weg? A) Das Cockpit zeigt die Zeile als Notice per say(), wie checkVersion - deckt nur das Cockpit ab. B) Bei --json wird der Hinweis weder gedruckt noch als gezeigt gemerkt; er kommt beim ersten Aufruf ohne --json - deckt alle --json-Aufrufer ab, kehrt aber 'regardless of --json' (Kommentar in internal/cli/update.go) um und nimmt PW18DT (Warnungen bei --json) ein Stueck vorweg. Empfehlung: B."
---

# Die Warnung 'Board von aelterer Version eingerichtet' kommt einmal, nicht bei jedem Aufruf

## Definition of Done

- [x] Die Warnung erscheint einmal pro Board und Binary-Version (Merker neben dem Versionsstempel) und danach nur noch bei jaira update selbst oder --verbose; ein Test deckt den zweiten Aufruf ab
  proof: core/release/release.go Nudged/MarkNudged; internal/cli/update.go nudgeIfStale; covered by TestNudgeIfStaleShowsOncePerBoardAndVersion, TestUpdateAfterShownHintStillReportsPreviousAndStamps
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/release/release.go: beside Stamped/Stamp add a marker file in the same state dir that records which binary version the stale-board hint was last shown for (read swallowing errors like Stamped; write like Stamp)
- [x] internal/cli/update.go nudgeIfStale: after the stamp check, stay silent when the marker already names release.Current; otherwise print the line and write the marker. No --verbose flag exists, none is added; jaira update keeps reporting the previous version and what changed
- [x] internal/cli/update_test.go: second nudge on the same board and version is silent; a different binary version nudges once again; jaira update after a shown hint still reports previous version and stamps
- [x] core/release/notes/RNQRJB.md: one '- ' line written as an instruction; update the nudgeIfStale doc comment

## Progress
- **2026-10-26 20:07 · BeMuCa** — Ticket opts out of brainstorm and planning (Options unchecked), so next_lane jumps todo -> in-progress; the driver wrote the plan. Interpretation taken: 'danach nur noch bei jaira update selbst' = after the one shown hint, the staleness only surfaces through jaira update's own report (previous version + what changed); update already suppresses the nudge via updateInProgress. There is no --verbose flag (global flags are -C and --json), so none is added.
- **2026-10-26 20:09 · BeMuCa** — Marker is a separate file <stateDir>/nudged holding the binary version; the version stamp itself is not touched, so jaira update still reads the true previous version. Marker write errors are swallowed (hint then just repeats).
- **2026-10-26 20:12 · BeMuCa** — critique 1: one must-fix. README.md:97 says the newer binary 'says so on every command' - this change makes that false; reword to once per board and jaira version, on the first command after the upgrade. Optional: internal/cli/update.go:24 is 88 chars in a paragraph wrapped at ~78 (reflow); core/release/notes/RNQRJB.md could say 'run jaira update when you see it, it will not remind you again' instead of 'to see the previous version' (update rewrites CLAUDE.md/AGENTS.md, it is not a status query). Checked and let stand, do not reopen: separate nudged file beside version (DoD names the place; stamp untouched keeps update's previous true, note 20:09); scope is per working tree because StateDir is per tree (core/ticket/store.go:822) - a fresh git worktree shows the hint once; RepoStateDir would make it once per clone but would sit apart from the stamp the DoD ties it to; Nudged/MarkNudged copy Stamped/Stamp except the file name - folding them into one helper would refactor existing code; no --verbose (note 20:07); swallowed write error (note 20:09); TUI (jaira, jaira board) never goes through openStore so it cannot use up the one hint unseen; scripts/lane-eval/clone.sh:55 and the cockpit OLDER_BUILD filter stay needed because the line still appears once.
- **2026-10-26 20:18 · BeMuCa** — critique 2: round-1 items all handled (README.md:97-98 reworded, update.go:24 reflowed, release note now says run jaira update when you see it). New finding, not cleared by critique 1 (it checked unseen consumption only for the TUI): the jaira-cockpit mod polls 'jaira lanes/list/boards --json' in the main checkout at plugin load, every 60s and after every turn (core/mod/builtin/jaira-cockpit/hooks/register.tsx:281-285, :1049-1053) and drops stderr on success. That poll is the first command after a binary upgrade, so it uses up the one hint and nobody sees it for that version. Reproduced on a scratch board: binaries built with -X main.version=0.5.0/0.5.9, first 'list --json 2>/dev/null' wrote nudged=0.5.9, the next visible 'list' printed nothing. Before this change the agent's own next command still showed it. Proposed fix: the cockpit shows the line it would drop - board.ts exports a non-global match for the OLDER_BUILD line, refreshBoard passes the first match from lanes/list/boards stderr to say() (toast + notice, same channel checkVersion uses at register.tsx:349-353), plus a case in tests/board.test.ts. Alternative: an env opt-out the cockpit sets so its polls never write the marker.
- **2026-10-26 20:25 · BeMuCa** — Driver, paused in human: under option B the remaining non-JSON silent consumer is core/role/builtin/jaira-dispatcher/scripts/spawn.sh:241 ('show --brief 2>/dev/null'); it would still use up the hint. Option A covers only the cockpit; run-lane.sh:76/137 and reap.sh:141 ('show --json 2>/dev/null') keep consuming it. The cockpit already strips the line via OLDER_BUILD in hooks/board.ts:202. Commits so far: 14b1c046 (code+tests+note), 2b0041d7 (docs after critique 1). DoD 2 (go test ./... -race) not run - house rule: full race suite only before the PR.
- **2026-10-26 20:26 · BeMuCa** — Lead (an agent, not a person) recommends option B for the human question: under --json the stale-board hint is neither printed nor marked as shown, so it appears once on the next plain call. Reason given: the hint is for people, and a --json caller is a program. Effect on PW18DT if adopted: it settles only the stale-board hint (a hint meant for people does not appear under --json at all, neither on stderr nor as a JSON field). It does not settle warnStaleClaims (internal/cli/flow.go), which is what PW18DT is about. Waiting for a person's answer in human; the driver did not move the ticket out of human.
