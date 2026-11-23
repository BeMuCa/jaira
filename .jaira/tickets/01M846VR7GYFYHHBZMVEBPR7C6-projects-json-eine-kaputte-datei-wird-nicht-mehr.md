---
id: 01M846VR7GYFYHHBZMVEBPR7C6
title: "projects.json: eine kaputte Datei wird nicht mehr still mit einem Eintrag überschrieben"
status: critique
ready: true
creator: Alexander Sacharov
goal: "Ist ~/.jaira/projects.json nicht lesbar, geht die Projektliste nicht verloren: jaira legt eine Kopie an, sagt es auf stderr und überschreibt das Original nicht still."
context: |-
  Was heute falsch ist: core/project/project.go:56 (Load) liest eine projects.json, die sich nicht parsen lässt, als leere Liste. Der nächste Remember() schreibt dann eine Liste mit nur dem aktuellen Projekt - alle anderen Einträge sind weg, ohne Meldung.
  Wann das passiert: Datei von Hand bearbeitet und dabei ein Komma vergessen, oder eine alte zerrissene Datei aus der Zeit vor ZXJ72E (das parallele Schreiben ist seit ZXJ72E atomar und gesperrt).
  Gefunden: Review von ZXJ72E als 'elsewhere'-Notiz, vom Maintainer am 23.11.2026 zur Reparatur freigegeben.
  Bekannt: Remember() schreibt seit ZXJ72E unter ticket.LockFile + ticket.WriteAtomic; core/settings macht für settings.json schon eine .bak-Kopie (settings.json.bak liegt in ~/.jaira).
definition-of-done: "Eine projects.json, die sich nicht parsen lässt, wird vor dem nächsten Schreiben nach projects.json.broken kopiert (eine bestehende .broken wird nicht still überschrieben, z.B. Zeitstempel im Namen)"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-11-23T12:42:28Z
updated-at: 2026-11-23T13:32:17Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-74262
claimed-at: 2026-11-23T13:30:23Z
outcome-what: "Remember copies an unparsable projects.json to projects.json.broken (or .broken-<time>) and says so on stderr before writing a new list; a failed copy leaves the original and returns an error"
outcome-why: "Load reads a broken file as an empty list, so the next Remember silently kept one board and dropped the rest"
outcome-resolves: "BPR7C6: a broken project list is no longer lost without a word"
---

# projects.json: eine kaputte Datei wird nicht mehr still mit einem Eintrag überschrieben

## Definition of Done

- [x] Eine projects.json, die sich nicht parsen lässt, wird vor dem nächsten Schreiben nach projects.json.broken kopiert (eine bestehende .broken wird nicht still überschrieben, z.B. Zeitstempel im Namen)
  proof: core/project/project.go:keepBroken, TestRememberKeepsBrokenList
- [x] jaira meldet auf stderr einmal, dass die Liste nicht lesbar war und wo die Kopie liegt
  proof: core/project/project.go:keepBroken Fprintf to Stderr; TestRememberKeepsBrokenList (one line, then silent)
- [x] Eine leere oder fehlende projects.json verhält sich wie bisher (keine Kopie, keine Meldung)
  proof: TestRememberMissingOrEmptyListIsQuiet
- [x] Tests: kaputte Datei -> Kopie + Meldung + neue Liste; fehlende Datei -> nichts davon
  proof: core/project/worktree_test.go: TestRememberKeepsBrokenList, TestRememberMissingOrEmptyListIsQuiet
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/BPR7C6.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 13:30 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/6958.sock
- **2026-11-23 13:30 · Alexander Sacharov** — decision: no open decisions for the person. Copy goes to projects.json.broken; if that exists, projects.json.broken-<UTC yyyymmddThhmmssZ> (never overwritten). Copy+message only in Remember (the write path), not in Load: a read alone loses nothing. Empty/whitespace-only file counts as missing. Message goes to a package writer (os.Stderr by default, swappable in tests); a TUI board switch mid-session could print over the screen, but the first Remember runs before the TUI starts (internal/cli/home.go openBoardAt) and rewrites the file valid, so that only hits a file broken by hand during a session.
