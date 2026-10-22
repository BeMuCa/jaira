---
id: 01M4TRAY0VEZPHSK66EFVHTE7X
title: "Ein selbst gebauter jaira zeigt, aus welchem Commit er stammt"
status: in-progress
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
updated-at: 2026-10-22T08:16:54Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-4187645
claimed-at: 2026-10-22T07:55:10Z
outcome-what: "A source build now shows dev-<7-char commit>, plus -dirty for uncommitted changes, in 'jaira --version' and the TUI header, via release.Display() (core/release/release.go); release.Current stays \"dev\" and a release build shows its tag unchanged."
outcome-why: "A bare 'dev' looked the same for every self-built binary, so nobody could tell which build was running."
outcome-resolves: "DoD 1: Display/displayFrom feed both cobra's Version (internal/cli/root.go:101) and versionLine (internal/tui/updatecheck.go:57); TestDisplayFrom covers release, clean, dirty and no-vcs cases; checked on real binaries built from a plain git clone. DoD 2's race suite is left to the maintainer before the PR; the release note is core/release/notes/VHTE7X.md."
review-summary: |-
  Neu: release.Display() (core/release/release.go:43-79).
  Liest per debug.ReadBuildInfo die Go-Stempel vcs.revision und vcs.modified.
  Ist release.Current nicht 'dev' (Release, -ldflags): Current kommt unveraendert zurueck.
  Ist es 'dev' und Go hat eine Revision gestempelt: 'dev-<7 Zeichen>', bei vcs.modified=true dazu '-dirty'.
  Ohne Revision bleibt es 'dev'.
  Verdrahtet an zwei Stellen:
  - internal/cli/root.go:101 gibt Display() an cobra als Version -> 'jaira --version'.
  - internal/tui/updatecheck.go:57 nutzt es in der Versionszeile oben links (Board und Startbildschirm).
  release.Current selbst bleibt 'dev'. Board-Stempel, self upgrade, update und market vergleichen weiter mit Current.
  Test: TestDisplayFrom (core/release/display_test.go), vier Faelle: Release, sauber, dirty, ohne VCS-Daten.
  Release-Notiz: core/release/notes/VHTE7X.md, eine Zeile.
review-gaps: "1. Defekt, neu durch diesen Diff: falscher Commit bei einem Worktree INNERHALB eines Checkouts.\n   Beispiel: <checkout>/.claude/worktrees/<name> - dort legt Claude Code seine Worktrees an (.gitignore:25).\n   Go findet dann das .git-Verzeichnis des aeusseren Checkouts und stempelt dessen HEAD.\n   Nachgestellt: Worktree-HEAD b0fcf07, Binary sagt 'jaira version dev-d23f19f'.\n   Mit geaenderter main.go im Worktree: trotzdem kein '-dirty'.\n   Vorher stand dort nur 'dev' - unergiebig, aber nicht falsch.\n2. Release-Notiz verspricht zu viel: 'A jaira you built yourself now names its commit'.\n   Stimmt nicht fuer: Build aus einem verlinkten Worktree wie /home/berk/git/jAIra-<ID> (.git ist eine Datei -> plain 'dev'), den Fall aus Punkt 1, und 'go run ./cmd/jaira' (plain 'dev').\n   Die Grenze steht nur in den Ticket-Notizen - weder in der Notiz noch im Kommentar an Display.\n3. '-dirty' heisst: 'git status --porcelain' ist nicht leer (GOROOT/src/cmd/go/internal/vcs/vcs.go:233).\n   Eine unversionierte Datei oder eine geaenderte Ticketdatei unter .jaira/ reicht.\n   /home/berk/git/jAIra hat gerade 12 solche Eintraege, alle unter .jaira/ -> ein Build dort zeigt '-dirty' ohne Codeaenderung.\n   Nicht falsch, aber die Notiz sollte es nicht als Code-Aenderung lesen lassen.\n4. Veralteter Kommentar: internal/tui/updatecheck.go:39 zitiert die Zeile noch als \"jaira dev\".\n5. Klein, kein Rueckgabegrund: kein Test bricht, wenn root.go:101 oder updatecheck.go:57 wieder plain Current zeigen.\n   Grund: Testbinaries tragen ohne -buildvcs=true keine VCS-Daten (GOROOT/src/cmd/go/internal/load/pkg.go:2506), Display() ist im Test immer 'dev'.\n   TestDisplayFrom deckt die Logik, nicht die Verdrahtung.\n6. DoD 2 offen: 'go test ./... -race' ist nicht gelaufen (Hausregel: erst beim Maintainer vor dem PR)."
test-verdict: "pass: gates green (gofmt/vet/build, targeted tests); real builds show dev-d23f19f, dev-d23f19f-dirty, 0.9.9; TUI header dev-d23f19f / jaira 0.9.9; release note 1 line. Full race suite (DoD 2) open for maintainer before PR."
review-verdict: |-
  Empfehlung: zurueck nach in-progress, kleine Nacharbeit ohne neue Logik.
  Kernfall stimmt, am echten Binary aus einem normalen Clone geprueft:
  - sauber: 'jaira version dev-d23f19f'
  - geaendert: 'jaira version dev-d23f19f-dirty'
  - -ldflags 0.9.9: 'jaira version 0.9.9'
  - TUI erste Zeile: 'jaira dev-d23f19f'
  release.Current bleibt 'dev'. market.go:78, update.go:39, self.go:80, updatecheck.go:56 und der Board-Stempel (release.go:221) lesen weiter Current.
  newRoot bekommt Display() nur fuer cobra.Version, sonst nichts.
  Cockpit-Regex (board.ts:516) findet in 'version dev-...' keine X.Y.Z, also keine Warnung - wie vorher. Installskripte installieren nur Releases.
  Aber: im Worktree innerhalb eines Checkouts zeigt das Binary einen falschen Commit (Gap 1), und die Release-Notiz verspricht den Commit ohne Einschraenkung (Gap 2).
  Zur Laufzeit erkennt der Code das nicht. Abhilfe ohne neues Build-Kommando:
  - Grenze in die Release-Notiz (eine Zeile) und in den Kommentar an Display schreiben.
  - Kommentar updatecheck.go:39 nachziehen.
  Ob der Commit auch in Worktree-Builds stimmen soll (Commit per -ldflags aus einem Build-Skript), ist Berks Entscheidung - eigenes Ticket.
  Nimmt Berk die Grenze so hin, reicht signoff.
  DoD 2 nicht erfuellt: 'go test ./... -race' habe ich nicht laufen lassen (Hausregel 3; der Maintainer fuehrt ihn einmal vor dem PR aus).
  Gelaufen: gofmt -l . leer; go vet ./... und go build ./... gruen; TestDisplayFrom PASS; go test ./internal/tui -run 'Dev|Version' ok; go test ./internal/cli -run 'Version|Root' ok.
review-round: 1
review-check: "Alles in einem Wegwerf-Ordner. Kein echtes Board, kein Netz.\n\n1. Leeren Pruefordner anlegen:\n   rm -rf /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck && mkdir -p /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck\n   Erwartet: keine Ausgabe.\n2. Normalen Clone holen (kein Worktree):\n   git clone -q /home/berk/git/jAIra-VHTE7X /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/clone && git -C /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/clone checkout -q d23f19fe\n   Erwartet: keine Ausgabe.\n3. Sauber bauen und Version fragen:\n   cd /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/clone && go build -o /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/jaira-clean ./cmd/jaira && /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/jaira-clean --version\n   Erwartet: jaira version dev-d23f19f\n4. Eine Datei aendern, neu bauen, fragen:\n   echo '// x' >> /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/clone/cmd/jaira/main.go && cd /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/clone && go build -o /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/jaira-dirty ./cmd/jaira && /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/jaira-dirty --version\n   Erwartet: jaira version dev-d23f19f-dirty\n5. Aenderung zuruecknehmen:\n   git -C /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/clone checkout -- cmd/jaira/main.go && git -C /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/clone status --porcelain\n   Erwartet: keine Ausgabe.\n6. Release-Build nachstellen:\n   cd /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/clone && go build -ldflags \"-X main.version=0.9.9\" -o /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/jaira-rel ./cmd/jaira && /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/jaira-rel --version\n   Erwartet: jaira version 0.9.9  (kein dev, kein Commit)\n7. Leeres Repo ohne Remote anlegen:\n   mkdir -p /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/board /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/home && cd /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/board && git init -q .\n   Erwartet: keine Ausgabe.\n8. Board darin anlegen:\n   cd /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/board && HOME=/tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/home JAIRA_HOME=/tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/home/.jaira JAIRA_NO_FETCH=1 JAIRA_NO_UPDATE_CHECK=1 /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/jaira-clean init\n   Erwartet als letzte Zeile: Run 'jaira share' when you want your team to have it.\n9. Board oeffnen:\n   cd /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/board && HOME=/tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/home JAIRA_HOME=/tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/home/.jaira JAIRA_NO_FETCH=1 JAIRA_NO_UPDATE_CHECK=1 /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/jaira-clean board\n   Erwartet: allererste Zeile oben links 'jaira dev-d23f19f'. Darunter 'board' und '0 tickets'. Mit q beenden.\n10. Bekannte Grenze - Build aus einem verlinkten Worktree:\n   git -C /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/clone worktree add -q --detach /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/wt d23f19fe && cd /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/wt && go build -o /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/jaira-wt ./cmd/jaira && /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/jaira-wt --version\n   Erwartet: jaira version dev  (ohne Commit: Go erkennt .git nur als Verzeichnis, im Worktree ist es eine Datei)\n11. Befund Gap 1 - Worktree innerhalb des Clones, mit eigenem Commit:\n   git -C /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/clone worktree add -q -b nested /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/clone/.claude/worktrees/nt d23f19fe && cd /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/clone/.claude/worktrees/nt && git -c user.name=x -c user.email=x@example.invalid commit -q --allow-empty -m nested && git rev-parse --short=7 HEAD\n   Erwartet: ein Hash, der NICHT d23f19f ist.\n12. Dort bauen und fragen:\n   cd /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/clone/.claude/worktrees/nt && go build -o /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/jaira-nt ./cmd/jaira && /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck/jaira-nt --version\n   Erwartet heute: jaira version dev-d23f19f - der Commit des aeusseren Clones, nicht der aus Schritt 11. Das ist der Fehler.\n13. Aufraeumen:\n   rm -rf /tmp/claude-1000/-home-berk-git-jAIra/bf80deba-b5c9-4820-848c-472cc90df7c4/scratchpad/vcheck\n   Erwartet: keine Ausgabe.\n14. Zum Schluss DoD 2, die ganze Suite (dauert mehrere Minuten):\n   cd /home/berk/git/jAIra-VHTE7X && go test ./... -race -count=1 2>&1 | tail -25\n   Erwartet: jede Zeile beginnt mit 'ok' oder '?'. Nirgends FAIL."
round: 2
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
- **2026-10-22 08:08 · BeMuCa** — optimize pass 1 review-gaps (kept here because review overwrites the field): none. core/gitrepo/git.go:155 shortSHA also cuts 7 chars but is unexported in a package that shells out to git, so displayFrom keeps its own two lines rather than coupling release to gitrepo. testing pass 1: pass on clone builds (dev-d23f19f, dev-d23f19f-dirty, 0.9.9) and the TUI header on a scratch board; full race suite (DoD 2) left for the maintainer before the PR.
- **2026-10-22 08:16 · BeMuCa** — review: zurueck nach in-progress empfohlen. (1) Build in einem Worktree INNERHALB eines Checkouts (z.B. <checkout>/.claude/worktrees/nt, .gitignore:25) zeigt den Commit des aeusseren Checkouts: nachgestellt Worktree-HEAD b0fcf07, Binary 'jaira version dev-d23f19f', mit geaenderter main.go kein -dirty - Go findet das aeussere .git-Verzeichnis. (2) Release-Notiz core/release/notes/VHTE7X.md verspricht den Commit ohne Einschraenkung; verlinkter Worktree und 'go run ./cmd/jaira' zeigen plain dev, Fall (1) einen falschen Commit. (3) -dirty kommt schon von einer geaenderten Ticketdatei unter .jaira/ oder einer unversionierten Datei (git status --porcelain). (4) Kommentar internal/tui/updatecheck.go:39 zitiert noch 'jaira dev'. Abhilfe ohne neue Logik: Grenze in Notiz-Zeile und Display-Kommentar nennen, Kommentar :39 nachziehen. Commit auch in Worktree-Builds (per -ldflags aus einem Build-Skript) waere ein eigenes Ticket, Berks Entscheidung.
- **2026-10-22 08:16 · BeMuCa** — review pass 1: back to in-progress for text only. (1) A worktree INSIDE a checkout (e.g. <checkout>/.claude/worktrees/<name>, .gitignore:25) gets the OUTER checkout's HEAD stamped - Go walks past the worktree's .git file to the outer .git dir; reproduced: worktree HEAD b0fcf07, binary said dev-d23f19f, no -dirty for a change in the worktree. Cannot be detected at runtime; fixing it needs a commit passed by -ldflags from a build script = Berk's decision, separate ticket. (2) Release note promises the commit without limit; say: linked worktree -> plain dev, worktree inside a checkout -> that checkout's commit, go run -> dev. (3) -dirty = git status --porcelain not empty, so untracked files and changed .jaira/ tickets count; say so. (4) internal/tui/updatecheck.go:39 comment still quotes "jaira dev". (5) no wiring test - accepted, test binaries carry no vcs data.
