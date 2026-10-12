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
updated-at: 2026-10-12T15:59:07Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-57410
claimed-at: 2026-10-12T15:43:47Z
outcome-what: ".gitattributes pinnt core/release/notes/**/*.md; release-notes.yaml testet den geschnittenen Stand vor dem Push; TestMatchPattern deckt notes/<version>/"
outcome-why: "v0.3.7 brach im Tag-Build ab, weil die neuen notes/0.3.7/-Dateien von keiner .gitattributes-Zeile gepinnt waren und der Release-PR keine CI startet"
outcome-resolves: "Release v0.3.7 kann neu getaggt werden; ein kuenftiger notes/<version>/-Ordner scheitert nicht mehr, und ein aehnlicher Fehler stoppt den Knopf statt den Tag-Build"
review-summary: "Drei Aenderungen: (1) .gitattributes pinnt jetzt core/release/notes/**/*.md statt notes/*.md - git check-attr liefert eol=lf fuer notes/X.md, notes/0.3.7/ und jeden kuenftigen notes/<version>/-Ordner, damit faellt TestRepositoryIsClean (wintrap) im Tag-Build nicht mehr. (2) release-notes.yaml fuehrt nach 'go run ./core/release/cut' und vor dem Push 'go test ./...' auf dem geschnittenen Arbeitsbaum aus - derselbe Befehl wie der goreleaser-before-Hook, also stoppt ein solcher Fehler kuenftig den Knopf statt den Tag-Build. (3) TestMatchPattern deckt notes/X.md und notes/9.9.9/X.md gegen das neue Muster ab."
review-gaps: "Klein: Der neue Schritt 'Test the cut tree' setzt keine globale git-Identitaet, waehrend ci.yaml und release.yaml das ausdruecklich 'fuer die Store-Tests' tun. Geprueft: go test ./... laeuft mit GIT_CONFIG_GLOBAL=/dev/null und user.useConfigOnly=true komplett gruen (die Tests setzen ihre Identitaet selbst), heute also kein Defekt - nur eine Abweichung, die ein kuenftiger Test ohne eigene Identitaet erst im Knopf zeigt. Nicht pruefbar ausser auf GitHub: der Workflow-Schritt selbst ist nie gelaufen. Offen und ausserhalb des Diffs: Tag v0.3.7 neu setzen und release.yaml starten (Mensch)."
test-verdict: "pass: go test ./... RC=0, go vet RC=0; git check-attr eol=lf fuer notes/X.md, notes/0.3.7/, notes/9.9.9/; Knopf simuliert (cut 9.9.9 in temp-Worktree, danach go test ./... RC=0); alte Zeile macht TestRepositoryIsClean rot. Workflow selbst nicht auf GitHub gelaufen."
review-verdict: "Diff erfuellt die DoD: mit neuer Zeile wintrap gruen mit notes/0.3.7/ und einem angelegten notes/9.9.9/, mit alter Zeile rot (selbst nachgeprueft). Kein Defekt gefunden; einziger Vorbehalt ist der ungelaufene Workflow-Schritt auf GitHub."
review-check: "1. git check-attr eol -- core/release/notes/X.md core/release/notes/9.9.9/X.md -> beide Zeilen 'eol: lf'  2. mkdir core/release/notes/9.9.9 && echo '- probe' > core/release/notes/9.9.9/ZZ.md  3. go test -count=1 ./internal/wintrap -> ok  4. rm -r core/release/notes/9.9.9  5. Nach dem Merge: Tag v0.3.7 auf den Merge-Commit neu setzen, gh workflow run release.yaml --ref v0.3.7 -> Lauf gruen, GitHub-Release v0.3.7 existiert  6. Beim naechsten Release-Knopf: im Actions-Log erscheint der Schritt 'Test the cut tree' vor 'Push release/<version>' und ist gruen"
---

# Release v0.3.7 scheitert: .gitattributes pinnt die Versionsordner unter notes/ nicht

## Definition of Done

- [x] .gitattributes pinnt 'core/release/notes/**/*.md text eol=lf' (alle Versionsordner); go test ./internal/wintrap gruen mit einem notes/0.3.7/-Ordner und mit einem kuenstlichen neuen Ordner (z. B. notes/9.9.9/)
  proof: .gitattributes:26; TestRepositoryIsClean gruen mit notes/0.3.7/ und mit kuenstlichem notes/9.9.9/ (mit alter Zeile rot); TestMatchPattern internal/wintrap/wintrap_test.go:119-120
- [x] release-notes.yaml faehrt vor dem Oeffnen des Release-PR 'go test ./internal/wintrap ./core/release/...' auf dem geschnittenen Stand, so dass ein solcher Fehler den Knopf stoppt statt den Tag-Build
  proof: .github/workflows/release-notes.yaml:48-49 'Test the cut tree' (go test ./..., wie .goreleaser.yaml:11) nach cut und git-Identitaet (:38), vor Push/PR
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
- **2026-10-12 15:52 · Alexander Sacharov** — Review: DoD selbst nachgeprueft (alte Zeile -> TestRepositoryIsClean rot, neue -> gruen mit notes/9.9.9/). go test ./... ohne globale git-Identitaet und mit user.useConfigOnly=true gruen - der fehlende 'git config --global' im neuen Workflow-Schritt ist heute harmlos, nur eine Abweichung zu ci.yaml/release.yaml. Basis fuer den Diff war f212ed3, nicht origin/master (der Fork-Remote ist veraltet).
- **2026-10-12 15:54 · Alexander Sacharov** — Alex 2026-10-12: send back for review gap (git identity): copy the 'Configure git for fixture repositories' step from release.yaml into release-notes.yaml before the test step.
- **2026-10-12 15:59 · Alexander Sacharov** — Git-Identitaet: Schritt 'Configure git for fixture repositories' aus release.yaml uebernommen, vor 'Test the cut tree'. Der Push-Schritt setzt danach 'git config user.name github-actions' repo-lokal - das schlaegt --global, der Release-Commit bleibt also von github-actions.
