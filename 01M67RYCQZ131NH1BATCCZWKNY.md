---
id: 01M67RYCQZ131NH1BATCCZWKNY
title: "Release-Notes beginnen mit einer Hook-Zeile, und wer mehrere Versionen ueberspringt, sieht die Notes aller dazwischen"
status: in-progress
ready: true
creator: berkc
assignee: berkc
goal: "Jede Release-Note-Datei hat als erste Zeile einen Einzeiler, der wie ein Instagram-Hook neugierig macht; jaira update und das Release-Notes-Popup zeigen ihn zuoberst. Wer von 0.5.4 auf 0.5.8 springt, sieht die Notes von 0.5.5, 0.5.6, 0.5.7 und 0.5.8, nicht nur die der neuesten."
context: |-
  Was fehlt: grosse Aenderungen (Modell pro Lane, Default-Worker, Teamlead-Default-Modelle) gehen zwischen vielen kleinen Zeilen unter.
  Ausloeser: Maintainer 31.10.2026, Punkt 13: 'die erste Linie eine Oneliner beschreibung wie eine instagram hook ... Und ausserdem wenn man nicht eine version, sondern mehrere ueberspringt sollte man trotzdem sehen welche notes bei den anderen relevant war.'
  Bekannt: Format ist ein Zeilenscan (CLAUDE.md, Abschnitt Release-Notes): '- ', '- ! ', MXMDFH bringt '- + ' fuer Features. Ein neues Praefix fuer die Hook-Zeile muss alte Leser nicht brechen. Ob update heute schon alle uebersprungenen Versionen zeigt: pruefen (GK1QGT blaettert pro Version).
definition-of-done: "Ein Praefix fuer die Hook-Zeile ist festgelegt, dokumentiert in CLAUDE.md und core/release, und alte Leser zeigen sie als normale Zeile"
tags:
  - release
blocked-by: []
related:
  - 01M67K70X843TBF0TYRPMXMDFH
commits: []
created-at: 2026-10-31T01:24:49Z
updated-at: 2026-11-03T17:28:57Z
updated-by: berkc
claimed-by: berk-ca-130447
claimed-at: 2026-11-03T17:20:45Z
---

# Release-Notes beginnen mit einer Hook-Zeile, und wer mehrere Versionen ueberspringt, sieht die Notes aller dazwischen

## Definition of Done

- [x] Ein Praefix fuer die Hook-Zeile ist festgelegt, dokumentiert in CLAUDE.md und core/release, und alte Leser zeigen sie als normale Zeile
  proof: CLAUDE.md + AGENTS.md '- > ' bullet; core/release/release.go Hook + package doc; TestHookLineStaysAPlainChange
- [x] jaira update und das Popup zeigen die Hook-Zeilen zuerst; Test
  proof: release.entryOf hooks first: TestEntryPutsEveryHookOfAReleaseFirst; TestUpdateNotesShowHooks; TestNotesViewPutsHooksFirst
- [x] Ein Board, das 3 Versionen zurueckliegt, sieht die Notes aller 3; Test
  proof: TestSinceThreeVersionsBehindReturnsAllThree; TestUpdateThreeVersionsBehindPrintsAllThree; TestPopupThreeVersionsBehindPagesAllThree
- [x] Release-Notiz unter core/release/notes/
  proof: core/release/notes/CZWKNY.md

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/release/release.go: Hook() (CutPrefix '> ') after Feature; entryOf stable-partitions hooks first; package doc paragraph on the three markers and old readers
- [x] core/release/release_test.go: TestHookMarker, TestHookLineStaysAPlainChange, TestEntryPutsEveryHookOfAReleaseFirst, TestSinceThreeVersionsBehindReturnsAllThree (fstest 0.5.4..0.5.7, stamp 0.5.4)
- [x] internal/cli/update.go writeNotes: hook line prints '  > text' first branch; tests TestUpdateNotesShowHooks, TestUpdateThreeVersionsBehindPrintsAllThree (stamp 0.5.5, current 0.5.8, real notes)
- [x] internal/tui/releasenotes.go: styHook bold; lines() emits hooks before Features header; changeLines bullet '> '; tests TestNotesViewPutsHooksFirst, TestPopupThreeVersionsBehindPagesAllThree (popupBoard 0.5.8/0.5.5)
- [x] CLAUDE.md and AGENTS.md: '- > ' hook bullet after the '- + ' bullet in the release-notes section (never NOTES.md, never notes/<version>/, never MXMDFH.md)
- [x] core/release/notes/CZWKNY.md: '- > ' hook line + json-order line
- [x] gofmt -l, go vet, go test ./core/release ./core/shipped, ./internal/cli -run TestUpdate, ./internal/tui -run 'NotesView|Popup|NotesBase'; tick DoD with proofs; move to critique; one commit code+ticket

## Progress
- **2026-11-03 17:20 · berkc** — teamlead 03.11.2026: baut auf MXMDFH (Popup, Praefix '- + ') auf, das noch nicht auf master ist. Branch feat/CZWKNY-release-hook startet von origin/feat/MXMDFH-release-notes-popup; der PR geht nach MXMDFH.
- **2026-11-03 17:23 · berkc** — pre-process findings: skipped versions ALREADY work. release.Since(stamp) (core/release/release.go selectSince) returns every entry newer than the stamp; jaira update (internal/cli/update.go writeNotes) prints all of them; the popup (internal/tui/releasenotes.go openReleaseNotes, GK1QGT) pages one version per page through all of them. Only missing: a deterministic test for '3 versions behind' (the existing popup test reads real notes and t.Skips under 3). DoD 3 = tests only, no code.
- **2026-11-03 17:23 · berkc** — decisions settled without asking (disagree at review): (1) hook prefix is '- > ' — any '- X ' line is still one plain change to changeLines, so older readers show '> text'; (2) a hook does not combine with + or !; (3) hooks go first PER VERSION, done once in release.entryOf (stable partition), so update text, update --json and popup all get the order from one place; the popup additionally puts them above the Features header; (4) 'every note file opens with a hook' is documented in CLAUDE.md/AGENTS.md, NOT enforced by a test — a test would fail every in-flight branch's note file; (5) MXMDFH.md (another ticket's unreleased note) is left alone; (6) NOTES.md is closed history, so 'core/release docs' = the package doc of core/release/release.go.
- **2026-11-03 17:27 · berkc** — DEVIATION (test only): plan's TestPopupThreeVersionsBehindPagesAllThree failed: the popup title reads 'What is new since 0.5.5 · 4 unread', so strings.Contains("0.5.5 · ") matched the title, not a page. Also Since(0.5.5) returns 4 entries here (Unreleased(2) + 0.5.8 + 0.5.7 + 0.5.6), since notes/ holds unreleased files; the test already counts n dynamically. Fix: match the page header in its box, '│ <v> · ' (rendered as '│   │ 0.5.8 · 2/4'). Plan code otherwise kept.
- **2026-11-03 17:28 · berkc** — claim: 'jaira claim CZWKNY' refused (exit 3): claimed by berk-ca-130447 (4m before start). Did not pass --steal; ticket writes (dod, note, move) went through under assignee berkc. Teamlead should confirm the holder is the dispatching session.
