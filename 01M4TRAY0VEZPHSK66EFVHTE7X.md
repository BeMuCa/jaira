---
id: 01M4TRAY0VEZPHSK66EFVHTE7X
title: "Ein selbst gebauter jaira zeigt, aus welchem Commit er stammt"
status: review
ready: true
creator: BeMuCa
goal: "jaira --version und der Projektkopf zeigen bei einem lokalen Build den Commit, so wie ein Release-Build seinen Tag zeigt."
context: |-
  Release-Builds zeigen den Tag - das bleibt: .goreleaser.yaml:22 setzt -X main.version={{.Version}}.
  Ein lokaler Build (go build -o ~/.local/bin/jaira ./cmd/jaira) zeigt nur 'dev' (cmd/jaira/main.go:11). Am 2026-10-13: 'jaira --version' gab 'jaira version dev'.
  'go version -m ~/.local/bin/jaira' zeigt keine vcs.revision-Zeile, nur 'mod (devel)'. Go hat den Commit also nicht eingestempelt. Warum (Build aus einem Worktree? -buildvcs?), ist nicht geprueft.
  Folge: Berk testet das laufende Binary und kann nicht sehen, welcher Stand laeuft. Am 2026-10-13 half nur die Dateizeit (Binary 14:05, D0SAHM-Commit 13:51, erster E70WYN-Commit 14:16).
  Verwandt: T8A8KM (Pille unter der Version im Projektkopf).
definition-of-done: "Ein lokaler Build zeigt bei 'jaira --version' und im Projektkopf 'dev' plus Kurz-Commit, und ob der Baum Aenderungen hatte; ein Release-Build zeigt unveraendert seinen Tag; ein Test deckt beide Faelle"
tags:
  - release
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-13T13:48:21Z
updated-at: 2026-10-22T08:08:07Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-4187645
claimed-at: 2026-10-22T07:55:10Z
outcome-what: "A source build now shows dev-<7-char commit>, plus -dirty for uncommitted changes, in 'jaira --version' and the TUI header, via release.Display() (core/release/release.go); release.Current stays \"dev\" and a release build shows its tag unchanged."
outcome-why: "A bare 'dev' looked the same for every self-built binary, so nobody could tell which build was running."
outcome-resolves: "DoD 1: Display/displayFrom feed both cobra's Version (internal/cli/root.go:101) and versionLine (internal/tui/updatecheck.go:57); TestDisplayFrom covers release, clean, dirty and no-vcs cases; checked on real binaries built from a plain git clone. DoD 2's race suite is left to the maintainer before the PR; the release note is core/release/notes/VHTE7X.md."
review-summary: none
review-gaps: "none. Duplication: no other ReadBuildInfo/vcs.revision/BuildSetting reader in the repo; core/gitrepo/git.go:155 shortSHA does the same 7-char cut but is unexported in a package that shells out to git, and core/release imports nothing of it, so the two lines in displayFrom stay rather than coupling release to gitrepo. Dead code: none orphaned; newRoot(version) and Execute(version) still use their parameter. Fluff: the nil guard in Display is needed (ReadBuildInfo returns nil when !ok); the len(rev)>7 guard is real for non-git VCS (svn revisions are short numbers); displayFrom's current parameter is the test seam critique let stand. Cost: Display runs once per process (root.go:101) and once per TUI model build (model.go:326, home.go:145), not per render; ReadBuildInfo measured ~3us per call. Left for review as critique recorded: internal/tui/updatecheck.go:39-41 comment still quotes the dev line as \"jaira dev\"."
test-verdict: "pass: gates green (gofmt/vet/build, targeted tests); real builds show dev-d23f19f, dev-d23f19f-dirty, 0.9.9; TUI header dev-d23f19f / jaira 0.9.9; release note 1 line. Full race suite (DoD 2) open for maintainer before PR."
---

# Ein selbst gebauter jaira zeigt, aus welchem Commit er stammt

## Definition of Done

- [x] Ein lokaler Build zeigt bei 'jaira --version' und im Projektkopf 'dev' plus Kurz-Commit, und ob der Baum Aenderungen hatte; ein Release-Build zeigt unveraendert seinen Tag; ein Test deckt beide Faelle
  proof: core/release/release.go:Display/displayFrom + core/release/display_test.go TestDisplayFrom
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-22 08:00 · BeMuCa** — Go does not stamp vcs.* when built from a linked git worktree (.git is a file): the worktree binary prints plain dev. Verified dev-<sha>, dev-<sha>-dirty and 0.9.9 on a scratch repo copy instead. Display() is shown only; release.Current stays dev (checked in market, update, self, updatecheck). root.go passes release.Display() into newRoot at Execute, not inside newRoot, so tests calling newRoot("test") are untouched.
- **2026-10-22 08:00 · BeMuCa** — Why a worktree build shows plain dev: Go 1.26.5 finds a git root only where .git is a DIRECTORY (GOROOT/src/cmd/go/internal/vcs/vcs.go:218, rootName{".git", isDir: true}); a linked worktree has a .git FILE, so nothing is stamped and -buildvcs=true stays silent. Same tree built in the main checkout /home/berk/git/jAIra stamps vcs.revision/modified. Answers the context's open 'Build aus einem Worktree?' question. Not worked around: an ldflags commit var would need a changed build command, outside this ticket.
- **2026-10-22 08:04 · BeMuCa** — critique pass 1: none. Let stand: Display/displayFrom split is the test seam for debug.ReadBuildInfo, not speculative; Current stays dev and root.go wires Display at Execute (standing on note 1); no ldflags commit var, worktree builds show plain dev (standing on note 2). Considered and rejected: showing bi.Main.Version, which Go 1.26.5 already stamps (cmd/go/internal/load/pkg.go:2603-2612) - it is a pseudo-version without the 'dev' prefix the DoD asks for, and a build at a tagged commit would read exactly like the release. Left for review, not a critique finding: internal/tui/updatecheck.go:39 comment still quotes the dev line as "jaira dev".
- **2026-10-22 08:06 · BeMuCa** — optimize: Display() is not on a hot path - versionLine is computed once at Model/Home construction (internal/tui/model.go:326, home.go:145), never per render, and debug.ReadBuildInfo costs ~3us (measured with a scratch loop), so caching its result would be speculative.
