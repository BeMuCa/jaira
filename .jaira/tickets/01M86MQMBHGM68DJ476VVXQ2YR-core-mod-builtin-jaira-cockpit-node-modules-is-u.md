---
id: 01M86MQMBHGM68DJ476VVXQ2YR
title: Cockpit-Testlauf hinterlaesst node_modules/ - weder gitignored noch aus dem Binary gefiltert
status: critique
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
updated-at: 2026-11-26T07:58:56Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-38348
claimed-at: 2026-11-26T07:55:02Z
goal: "Ein Testlauf des jaira-cockpit-Mods hinterlaesst nichts, was git als neue Datei zeigt oder was ein danach gebautes jaira-Binary einbettet und installiert."
outcome-what: "Root-.gitignore ignoriert core/mod/builtin/*/node_modules/; mod.Builtins() laesst node_modules/ eines Mods aus (Walk ueber fs.FS via builtinsIn, damit testbar)"
outcome-why: "Ein Cockpit-Testlauf hinterliess node_modules/.vite als ?? in jedem worktree_diff, und ein danach gebautes Binary haette es eingebettet und installiert"
outcome-resolves: "DoD 1-2; go test ./... gruen, go vet, shipped-history --check"
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
