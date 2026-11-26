---
id: 01M86MQMBHGM68DJ476VVXQ2YR
title: Cockpit-Testlauf hinterlaesst node_modules/ - weder gitignored noch von der Mod-Installation ausgenommen
status: review
ready: true
creator: Sascha
context: |-
  Was falsch ist: ein Testlauf im Mod-Ordner (claude plugin test / vitest) legt core/mod/builtin/jaira-cockpit/node_modules/.vite/vitest/results.json an.
  Die Datei ist nicht gitignored (.gitignore im Repo-Root kennt kein node_modules ausser core/web/ui/.gitignore). Sie steht deshalb als ?? in git status und in jedem worktree_diff, den eine Lane liest.
  Zweiter Schaden, gleiche Ursache: core/mod/mod.go bettet builtin/ mit go:embed all: ein. go:embed liest die Platte, nicht git. Ein Binary, das nach einem Testlauf gebaut wird, traegt node_modules/ mit und installiert es mit dem Mod. Builtins() filtert heute nur tests/, .claude-plugin/types und tsconfig.json.
  Ausloeser: Review-Notiz auf N5ZZP6 (24.11.2026): 'the one untracked file is node_modules/.vite vitest cache'.
tags: []
blocked-by: []
related: []
follows: 01M86H5PFARVW0Z0NBMJN5ZZP6
commits: []
created-at: 2026-11-24T11:23:22Z
updated-at: 2026-11-26T08:04:36Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-38348
claimed-at: 2026-11-26T07:55:02Z
goal: "Ein Testlauf des jaira-cockpit-Mods hinterlaesst nichts, was git als neue Datei zeigt oder was ein danach lokal gebautes jaira-Binary mit dem Mod installiert."
outcome-what: "Root-.gitignore ignoriert core/mod/builtin/*/node_modules/; mod.Builtins() laesst node_modules/ eines Mods aus (Walk ueber fs.FS via builtinsIn, damit testbar)"
outcome-why: "Ein Cockpit-Testlauf hinterliess node_modules/.vite als ?? in jedem worktree_diff, und ein danach gebautes Binary haette es eingebettet und installiert"
outcome-resolves: "DoD 1-2; go test ./... gruen, go vet, shipped-history --check"
review-summary: "Die Root-.gitignore ignoriert jetzt /core/mod/builtin/*/node_modules/, sodass der vitest-Cache (node_modules/.vite) nach einem Mod-Testlauf nicht mehr in git status auftaucht. mod.Builtins() ueberspringt node_modules/ auf oberster Ebene eines Mods (neben tests/ und den generierten Typen), damit 'jaira roles install' den Cache nie mitinstalliert; dafuer wurde der Walk in builtinsIn(fs.FS) ausgelagert, den ein neuer Test mit fstest.MapFS prueft - inklusive dass ein tieferes hooks/node_modules erhalten bleibt. Kommentare sagen ehrlich, dass go:embed all: den Cache in einen lokalen Build weiter einbettet."
review-gaps: |-
  note: limit: ein lokal nach einem Testlauf gebautes Binary traegt die node_modules-Bytes weiter (go:embed all: kann kein Unterverzeichnis auslassen); installiert werden sie nicht, Release-Builds laufen auf frischem Checkout - so in mod.go:49-54 und .gitignore:41-44 dokumentiert
  note: limit: nur node_modules/ direkt im Mod-Ordner wird uebersprungen, tiefere (hooks/node_modules) bewusst nicht - der Test haelt das fest
test-verdict: "pass round 2: go test ./core/mod ok, gofmt clean; DoD 1-2 unchanged by the comment-only fix, proofs still point at .gitignore:43 (node_modules line) and TestBuiltinsLeaveATestRunsCacheOut"
hunt-verdict: "pass: stop-after-fix, no hunt in round 2"
round: 2
back-from: testing
back-by: lane
back-count: 1
stop-after-fix: "true"
review-verdict: "accept: beide DoD-Punkte erfuellt (.gitignore:44 greift laut git check-ignore, TestBuiltinsLeaveATestRunsCacheOut gruen); kein Defekt in core/mod/mod.go oder dem Install-Pfad core/role/install.go gefunden"
review-round: 2
review-verdict-gaps: |-
  note: limit: ein lokal nach einem Testlauf gebautes Binary traegt die node_modules-Bytes weiter (go:embed all: kann kein Unterverzeichnis auslassen); installiert werden sie nicht, Release-Builds laufen auf frischem Checkout - so in mod.go:49-54 und .gitignore:41-44 dokumentiert
  note: limit: nur node_modules/ direkt im Mod-Ordner wird uebersprungen, tiefere (hooks/node_modules) bewusst nicht - der Test haelt das fest
---

# core/mod/builtin/jaira-cockpit/node_modules/ is untracked and not gitig…

## Definition of Done

- [x] core/mod/builtin/*/node_modules/ ist in der Root-.gitignore: git status zeigt nach einem Testlauf im Mod-Ordner nichts unter node_modules
  proof: .gitignore:44; git check-ignore on node_modules/.vite/vitest/results.json matches it
- [x] mod.Builtins() laesst node_modules/ eines Mods aus, wie schon tests/ und tsconfig.json; Go-Test in core/mod deckt das ab
  proof: core/mod/mod.go testCache + builtinsIn; core/mod/mod_test.go TestBuiltinsLeaveATestRunsCacheOut

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-26 07:55 · Sascha** — parent session: uds:/run/user/1000/cc-socks/62036.sock
- **2026-11-26 07:56 · Sascha** — decision: Eintrag in die Root-.gitignore, nicht eine .gitignore im Mod-Ordner - eine Datei dort wuerde go:embed all: mit einbetten und mit dem Mod installieren. Muster /core/mod/builtin/*/node_modules/ fuer alle Mods, wie die tsconfig.json-Zeile daneben. Keine Release-Note: Nutzer eines Release-Binaries sehen davon nichts (Release-Build laeuft auf frischem Checkout).
- **2026-11-26 07:59 · Sascha** — waiting: bughunt (background agent) since 08:59
- **2026-11-26 08:01 · Sascha** — bughunt: change: core/mod/mod.go:23 (with mod.go:49-52, .gitignore:42, commit subject e9ba2b874) - class: a skip in builtinsIn only keeps a file off the install, while go:embed all:builtin still compiles every file on disk into the binary, so text that says 'out of the binary' is wrong for node_modules (and for the generated types/tsconfig.json beside it); example: core/mod/builtin/jaira-cockpit/node_modules/.vite/vitest/results.json present at build time is in the binary's bytes, and a full npm install there would add its whole size to every release binary; not a blocker by itself (nothing wrong is installed), unless this ticket's definition of done asks for it to be out of the binary; ran: mkdir node_modules/.vite/vitest + marker file, go build ./cmd/jaira, grep -a marker in binary (exit 0, found), install path list does not carry it
- **2026-11-26 08:01 · Sascha** — hunt report: 1 finding, 1 change (edge, not a blocker), as a note on VXQ2YR | read whole: core/mod/mod.go, core/mod/mod_test.go, .gitignore hunk, plus caller core/role/install.go (mod install loop); go test ./core/mod ./core/role pass | not checked: what 'claude plugin test' really writes into the mod folder beyond node_modules/.vite (no run of it here); a dangling symlink in node_modules/.bin was tried and does not break go build
- **2026-11-26 08:01 · Sascha** — bughunt: done at e9ba2b874 - 1 finding (1 change, 0 pre-existing) in 1 module, first hunt; concurrency, error paths, input, permissions checked; files: .gitignore core/mod/mod.go core/mod/mod_test.go
- **2026-11-26 08:02 · Sascha** — bughunt fix: wording only - go:embed all: cannot exclude a subdirectory, so a local build still carries node_modules bytes; mod.go and .gitignore now say install, not binary; goal and title reworded. stop-after-fix=true set by me (comments only, no lock, no parser). The commit subject of e9ba2b874 says 'binary' and stays: no amend after a move.
- **2026-11-26 08:03 · Sascha** — waiting: review (fresh agent via run-lane.sh) since 09:03
