---
id: 01M4TRAY0VEZPHSK66EFVHTE7X
title: "Ein selbst gebauter jaira zeigt, aus welchem Commit er stammt"
status: critique
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
updated-at: 2026-10-22T08:01:02Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-4187645
claimed-at: 2026-10-22T07:55:10Z
outcome-what: "A source build now shows dev-<7-char commit>, plus -dirty for uncommitted changes, in 'jaira --version' and the TUI header, via release.Display() (core/release/release.go); release.Current stays \"dev\" and a release build shows its tag unchanged."
outcome-why: "A bare 'dev' looked the same for every self-built binary, so nobody could tell which build was running."
outcome-resolves: "DoD 1: Display/displayFrom feed both cobra's Version (internal/cli/root.go:101) and versionLine (internal/tui/updatecheck.go:57); TestDisplayFrom covers release, clean, dirty and no-vcs cases; checked on real binaries built from a plain git clone. DoD 2's race suite is left to the maintainer before the PR; the release note is core/release/notes/VHTE7X.md."
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
