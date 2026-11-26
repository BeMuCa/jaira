---
id: 01M86MQMBHGM68DJ476VVXQ2YR
title: Cockpit-Testlauf hinterlaesst node_modules/ - weder gitignored noch aus dem Binary gefiltert
status: testing
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
updated-at: 2026-11-26T07:59:54Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-38348
claimed-at: 2026-11-26T07:55:02Z
goal: "Ein Testlauf des jaira-cockpit-Mods hinterlaesst nichts, was git als neue Datei zeigt oder was ein danach gebautes jaira-Binary einbettet und installiert."
outcome-what: "Root-.gitignore ignoriert core/mod/builtin/*/node_modules/; mod.Builtins() laesst node_modules/ eines Mods aus (Walk ueber fs.FS via builtinsIn, damit testbar)"
outcome-why: "Ein Cockpit-Testlauf hinterliess node_modules/.vite als ?? in jedem worktree_diff, und ein danach gebautes Binary haette es eingebettet und installiert"
outcome-resolves: "DoD 1-2; go test ./... gruen, go vet, shipped-history --check"
review-summary: "no findings: builtinsIn is the one indirection, needed because go:embed is fixed at compile time and the filter can only be tested over a stand-in fs.FS; testCache sits beside generatedTypes/generatedConfig in the same walk (core/mod/mod.go), same pattern; gitignore line sits beside the tsconfig.json line for the same folder pattern"
review-gaps: "none: 3 files, nothing duplicated elsewhere (core/web/ui/.gitignore covers only core/web/ui), every new symbol has a caller"
test-verdict: "pass: go test ./... green (exit 0), go test -race -count=1 ./core/mod ok, go vet, shipped-history --check ok; DoD 1: git check-ignore matches .gitignore:43 for node_modules/.vite/vitest/results.json; DoD 2: TestBuiltinsLeaveATestRunsCacheOut ok; function: in a git-archive copy with a real node_modules/.vite/vitest/results.json, the file is in the embedded FS but in no Mod.Files, so it is not installed"
---

# core/mod/builtin/jaira-cockpit/node_modules/ is untracked and not gitig…

## Definition of Done

- [x] core/mod/builtin/*/node_modules/ ist in der Root-.gitignore: git status zeigt nach einem Testlauf im Mod-Ordner nichts unter node_modules
  proof: .gitignore:43; git check-ignore on node_modules/.vite/vitest/results.json matches it
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
