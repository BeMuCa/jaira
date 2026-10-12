---
id: 01M4RCGVBXRPSJ3GA8P25HDAJZ
title: "Release v0.3.7 scheitert: .gitattributes pinnt die Versionsordner unter notes/ nicht"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "v0.3.7 wird gebaut und veroeffentlicht, und kein kuenftiger Release scheitert mehr an einem neuen notes/<version>/-Ordner."
context: |-
  Was falsch ist: der Release-Lauf https://github.com/BeMuCa/jaira/actions/runs/38271972105 (Tag v0.3.7, 2026-10-12) bricht in goreleaser ab: TestRepositoryIsClean (internal/wintrap/wintrap_test.go:83) meldet 'go:embed pulls in 7 file(s) that no line in .gitattributes pins, starting at core/release/notes/0.3.7/8WFNQY.md'.
  Ursache: .gitattributes:25 pinnt nur 'core/release/notes/*.md text eol=lf', nicht die Unterordner notes/<version>/, die der Release-Knopf (AR97K9) anlegt. Das Review von AR97K9 hatte das als 'kosmetisch' notiert (gap 1) - es ist es nicht, der Wintrap-Test verlangt es.
  Warum CI es nicht fing: der vom Knopf geoeffnete Release-PR startet ci.yaml nicht (GITHUB_TOKEN, von Alex so akzeptiert). Der Fehler zeigte sich erst im Tag-Build.
  Stand: Tag v0.3.7 steht auf 44c6fec (Merge des Release-PR), ein GitHub-Release v0.3.7 gibt es NICHT (gh release view: not found). Niemand hat 0.3.7 heruntergeladen.
  Nach dem Fix: Tag v0.3.7 loeschen und auf den Merge-Commit des Fixes neu setzen, dann 'gh workflow run release.yaml --ref v0.3.7' - das tut ein Mensch (Alex), nicht der Agent.
definition-of-done: ".gitattributes pinnt 'core/release/notes/**/*.md text eol=lf' (alle Versionsordner); go test ./internal/wintrap gruen mit einem notes/0.3.7/-Ordner und mit einem kuenstlichen neuen Ordner (z. B. notes/9.9.9/)"
tags:
  - release
  - ci
blocked-by: []
related:
  - 01M4QECQMKHK5BKWRZMFAR97K9
commits: []
created-at: 2026-10-12T15:43:24Z
updated-at: 2026-10-12T15:47:43Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-57410
claimed-at: 2026-10-12T15:43:47Z
outcome-what: ".gitattributes pinnt core/release/notes/**/*.md; release-notes.yaml testet den geschnittenen Stand vor dem Push; TestMatchPattern deckt notes/<version>/"
outcome-why: "v0.3.7 brach im Tag-Build ab, weil die neuen notes/0.3.7/-Dateien von keiner .gitattributes-Zeile gepinnt waren und der Release-PR keine CI startet"
outcome-resolves: "Release v0.3.7 kann neu getaggt werden; ein kuenftiger notes/<version>/-Ordner scheitert nicht mehr, und ein aehnlicher Fehler stoppt den Knopf statt den Tag-Build"
review-summary: ".github/workflows/release-notes.yaml:41 testet nur ./internal/wintrap und ./core/release/..., der Tag-Build (.goreleaser.yaml:11 before.hooks) aber 'go test ./...' - jeder andere Test, den ein Schnitt bricht (internal/cli/update_test.go liest release.Current), faellt weiter erst im Tag-Build. Den Schritt auf 'go test ./...' stellen, wie goreleaser; lokal 43 s."
---

# Release v0.3.7 scheitert: .gitattributes pinnt die Versionsordner unter notes/ nicht

## Definition of Done

- [x] .gitattributes pinnt 'core/release/notes/**/*.md text eol=lf' (alle Versionsordner); go test ./internal/wintrap gruen mit einem notes/0.3.7/-Ordner und mit einem kuenstlichen neuen Ordner (z. B. notes/9.9.9/)
  proof: .gitattributes:26; TestRepositoryIsClean gruen mit notes/0.3.7/ und mit kuenstlichem notes/9.9.9/ (mit alter Zeile rot); TestMatchPattern internal/wintrap/wintrap_test.go:119-120
- [x] release-notes.yaml faehrt vor dem Oeffnen des Release-PR 'go test ./internal/wintrap ./core/release/...' auf dem geschnittenen Stand, so dass ein solcher Fehler den Knopf stoppt statt den Tag-Build
  proof: .github/workflows/release-notes.yaml:42-43 'Test the cut tree' (go test ./..., wie .goreleaser.yaml:11) nach cut, vor Push/PR
- [x] Kein Eintrag in core/release/notes noetig, falls fuer Nutzer unsichtbar; sonst eine Zeile in core/release/notes/<ID>.md (Hinweis: notes/0.3.7/ ist geschlossen)
  proof: keine Note: nur .gitattributes und CI-Schritt, fuer Nutzer des Binaries unsichtbar

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-12 15:43 · Alexander Sacharov** — parent session: jaira-4c
- **2026-10-12 15:44 · Alexander Sacharov** — Offene Entscheidungen: keine. (1) Muster 'core/release/notes/**/*.md' ersetzt 'notes/*.md' - '**/' deckt null oder mehr Ordner, also bleiben die Unreleased-Dateien direkt in notes/ gepinnt (git und wintrap.matchPattern verhalten sich gleich, TestMatchPattern hat den Fall SKILL.md). (2) Der Testschritt laeuft in release-notes.yaml NACH 'cut' und VOR dem Push, damit er den geschnittenen Stand prueft und ein Fehler nichts pusht. (3) Keine Release-Note: .gitattributes und ein CI-Schritt sind fuer Nutzer des Binaries unsichtbar.
- **2026-10-12 15:45 · Alexander Sacharov** — Geprueft: mit der alten Zeile 'notes/*.md' faellt TestRepositoryIsClean (exit 1), mit '**/*.md' gruen. Der Workflow-Schritt testet den Arbeitsbaum nach 'go run ./core/release/cut' - cut verschiebt die Dateien auf der Platte, git add kommt erst danach, also sieht go test genau den geschnittenen Stand. Lokal nicht ausfuehrbar: der Workflow selbst (nur YAML-Syntax geprueft).
- **2026-10-12 15:46 · Alexander Sacharov** — critique: Testschritt in release-notes.yaml auf 'go test ./...' stellen - dasselbe, was .goreleaser.yaml:11 im Tag-Build faehrt, sonst bleibt jeder Test ausserhalb der zwei Pakete eine Luecke. Obermenge der in DoD 2 genannten Pakete; Kosten ~43 s lokal.
- **2026-10-12 15:47 · Alexander Sacharov** — Simuliert in einem temporaeren Worktree: Note ZZZZZZ.md angelegt, 'go run ./core/release/cut 9.9.9' -> notes/9.9.9/ZZZZZZ.md, danach 'go test ./...' exit 0. Worktree wieder entfernt.
