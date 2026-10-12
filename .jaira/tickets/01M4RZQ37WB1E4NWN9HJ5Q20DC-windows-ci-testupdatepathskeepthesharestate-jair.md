---
id: 01M4RZQ37WB1E4NWN9HJ5Q20DC
title: "Windows-CI: TestUpdatePathsKeepTheShareState/jaira_self_upgrade ersetzt die laufende Test-Binary"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: go test ./... -race ist auf windows-latest wieder gruen; der self-upgrade-Pfad im Share-State-Test tauscht eine Attrappe statt der laufenden cli.test.exe
context: "Seit PR #67 (8W1R94) in master ist, ist der Windows-Job der CI rot. Gesehen am 2026-10-12 im Lauf 38308305116 (PR #68), Job 115007598537.\nFehler: internal/cli/sharestate_test.go:98, Subtests 'shared clone/jaira self upgrade' und 'private board/jaira self upgrade':\n  rename ...\\go-build...\\b351\\cli.test.exe ...\\cli.test.exe.old-9748: Access is denied.\nUrsache: der Subtest ruft runSelf(t, dir, \"upgrade\") ohne eigenes Installationsziel; self upgrade ersetzt deshalb die laufende Test-Binary. Windows erlaubt kein Umbenennen einer laufenden .exe, Linux und macOS schon - darum dort gruen.\nErster Schritt: nachsehen, wie die bestehenden self-upgrade-Tests (core/selfupdate, internal/cli/*self*_test.go) das Ziel umbiegen (Env-Variable oder Helper fuer den Binary-Pfad), und denselben Weg hier benutzen. Nicht per runtime.GOOS ueberspringen, solange ein Ziel-Override existiert.\nFolge heute: jeder PR ist auf Windows rot, auch der gemeinsame PR feat/13pqb3."
definition-of-done: "TestUpdatePathsKeepTheShareState laeuft auf windows-latest gruen; der self-upgrade-Subtest schreibt in ein Temp-Ziel, nie in die laufende Test-Binary; CI des PRs auf ubuntu, macos und windows gruen"
tags:
  - ci
blocked-by: []
related: []
commits: []
created-at: 2026-10-12T21:18:51Z
updated-at: 2026-10-12T21:24:26Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-48475
claimed-at: 2026-10-12T21:19:18Z
outcome-what: self upgrade nimmt sein Ziel aus der Naht selfExecutable (internal/cli/self.go); selfUpgradeTarget legt eine Attrappe im TempDir an und biegt die Naht dorthin; der Share-State-Subtest ruft sie
outcome-why: "Ohne Override ersetzte jeder Upgrade-Test die laufende cli.test.exe; auf Windows scheitert der zweite Replace im selben Prozess am noch gesperrten .old-<pid>"
outcome-resolves: "Windows-Job: TestUpdatePathsKeepTheShareState/*/jaira_self_upgrade Access is denied"
---

# Windows-CI: TestUpdatePathsKeepTheShareState/jaira_self_upgrade ersetzt die laufende Test-Binary

## Definition of Done

- [~] TestUpdatePathsKeepTheShareState laeuft auf windows-latest gruen; der self-upgrade-Subtest schreibt in ein Temp-Ziel, nie in die laufende Test-Binary; CI des PRs auf ubuntu, macos und windows gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-12 21:19 · Alexander Sacharov** — parent session: jaira-f5
- **2026-10-12 21:20 · Alexander Sacharov** — Ursache genauer: ein Ziel-Override gibt es nicht - 'self upgrade' nimmt os.Executable() direkt (internal/cli/self.go:58), und selfUpgradeTarget in self_test.go zeigt auf dieselbe laufende cli.test.exe. Auf Windows benennt Replace die laufende .exe in .old-<pid> um; der zweite Replace im selben Prozess (2. Board-Variante) findet .old-<pid> schon vor und gesperrt -> Access is denied. Auch TestSelfUpgradeVersionPinsInstallsExactRelease ueberschreibt die Test-Binary.
- **2026-10-12 21:20 · Alexander Sacharov** — decision: keine offene Frage an Alex. Kleinster Weg ohne GOOS-Skip: eine Test-Naht 'var selfExecutable = os.Executable' in self.go; selfUpgradeTarget legt eine Attrappe jaira(.exe) im TempDir an und biegt die Naht dorthin. Der Share-State-Subtest ruft selfUpgradeTarget. Keine Release-Note: nur Tests betroffen, nichts von aussen sichtbar.
- **2026-10-12 21:24 · Alexander Sacharov** — Lokal gruen: go test ./... -race (exit 0), GOOS=windows go vet ./internal/cli (exit 0). Kompilierte cli.test vor/nach TestUpdatePathsKeepTheShareState+TestSelfUpgrade* per sha256 verglichen: unveraendert. Offen fuer DoD 1: der Windows-Job im PR-Lauf - lokal nicht pruefbar.
