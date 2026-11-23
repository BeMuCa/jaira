---
id: 01M8455PXYM7Y2A36J7YZXJ72E
title: "Projektliste: ein verknüpfter Worktree erscheint nicht als eigenes Projekt, sein Haupt-Checkout steht dafür"
status: in-progress
ready: true
creator: Alexander Sacharov
goal: "Die Projektliste des Launchers (Projects:) zeigt jedes Repository einmal, unter seinem Haupt-Checkout, statt jeden Worktree und jede Wegwerf-Kopie, in der jaira einmal lief."
context: |-
  Was heute falsch ist: jeder Start von jaira trägt seinen Ordner in ~/.jaira/projects.json ein (core/project/project.go, Remember()).
  Folge: auf der Maschine des Maintainers stehen dort 10 Einträge für eigentlich 2 Repositories - 6 davon sind verknüpfte git-Worktrees unter ~/projects/.worktrees/<repo>-<slug>, 3 sind Wegwerf-Repos in Scratch-Ordnern von Agenten (/tmp/claude-*/.../scratchpad/...) und Abnahme-Sandboxen (~/.claude/acceptance/.../sandbox).
  Der Launcher (internal/tui/home.go, NewHome/boardRoots) listet alle, die Projektauswahl (internal/tui/view.go renderProjects) ebenso.
  Wunsch des Maintainers (23.11.2026): automatisch erkennen - ist jaira im Haupt-Checkout offen und in einem Worktree desselben Repos auch, soll der Worktree nicht als eigenes Board erscheinen.
  Bekannt: core/board/gitignore.go:332 liest 'git worktree list' schon; 'git rev-parse --git-common-dir' sagt für einen verknüpften Worktree, zu welchem Haupt-Repo er gehört.
  Ausgeschlossen: eine Option zum Abschalten - es soll ohne Einstellung stimmen.
definition-of-done: "Ein Start in einem verknüpften Worktree trägt nicht den Worktree, sondern den Haupt-Checkout des Repos in projects.json ein (oder gar nichts, wenn der schon drinsteht)"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-11-23T12:12:57Z
updated-at: 2026-11-23T12:18:38Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-27553
claimed-at: 2026-11-23T12:13:22Z
---

# Projektliste: ein verknüpfter Worktree erscheint nicht als eigenes Projekt, sein Haupt-Checkout steht dafür

## Definition of Done

- [x] Ein Start in einem verknüpften Worktree trägt nicht den Worktree, sondern den Haupt-Checkout des Repos in projects.json ein (oder gar nichts, wenn der schon drinsteht)
  proof: core/project/project.go:89 Remember -> Board -> worktree.go:25 mainCheckout; TestRememberRecordsMainCheckout
- [x] Schon vorhandene Worktree-Einträge in projects.json zeigt der Launcher nicht mehr als eigene Zeile, wenn ihr Haupt-Checkout bekannt ist; ein Worktree, dessen Ordner nicht mehr existiert, verschwindet ebenfalls
  proof: core/project/project.go:65 Load folds via Board, drops missing .jaira; internal/tui/home.go:128; TestLoadFoldsOldWorktreeEntries
- [x] Repos unter dem Temp-Verzeichnis (os.TempDir) werden nicht eingetragen
  proof: core/project/worktree.go:93 throwaway, used in Remember and Load; TestRememberSkipsTemp
- [x] Ein Haupt-Checkout bleibt eine Zeile, auch wenn gerade mehrere seiner Worktrees jaira offen haben
  proof: core/project/project.go:65 dedupe by Board root; TestRememberRecordsMainCheckout (two worktrees + main = one line)
- [x] Tests für: Worktree -> Haupt-Checkout, Temp-Repo ignoriert, alter Worktree-Eintrag ausgeblendet
  proof: core/project/worktree_test.go: TestBoardIsMainCheckout, TestRememberSkipsTemp, TestLoadFoldsOldWorktreeEntries
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/ZXJ72E.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 12:13 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/6958.sock
- **2026-11-23 12:13 · Alexander Sacharov** — decisions (worker, none open for a person): 1) worktree->main detection reads files only (.git file 'gitdir:' + its 'commondir'), no git exec - Remember runs on every start and startup must stay instant. 2) the main checkout replaces a worktree only when it has its own .jaira dir; a worktree whose main checkout has no board (private board only in the worktree) stays as it is, else the board would vanish from the list. 3) Load() does the mapping/filtering, so every reader (launcher, switcher, 'jaira projects') sees the same list; Remember writes the cleaned list back. 4) temp = paths under os.TempDir() after canonical(); ~/.claude/acceptance sandboxes are not covered (DoD names only os.TempDir).
- **2026-11-23 12:18 · Alexander Sacharov** — found out: the shipped tests register t.TempDir() boards (internal/tui/projects_test.go, qback_test.go), so a plain os.TempDir filter would empty their registries; throwaway() stands aside when the registry (JAIRA_HOME) itself lives in temp - also true for any sandbox that points JAIRA_HOME into temp. Launcher extras discovered below cwd are mapped too (boardRoots), so 'jaira' started inside a worktree lists the main checkout; 'jaira board' still opens the worktree's own board. Not touched: projects.json has no lock, two starts at once can lose an entry - pre-existing.
