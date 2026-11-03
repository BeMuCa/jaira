---
id: 01M67RYCQZ131NH1BATCCZWKNY
title: "Release-Notes beginnen mit einer Hook-Zeile, und wer mehrere Versionen ueberspringt, sieht die Notes aller dazwischen"
status: pre-process
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
updated-at: 2026-11-03T17:23:22Z
updated-by: berkc
claimed-by: berk-ca-130447
claimed-at: 2026-11-03T17:20:45Z
---

# Release-Notes beginnen mit einer Hook-Zeile, und wer mehrere Versionen ueberspringt, sieht die Notes aller dazwischen

## Definition of Done

- [ ] Ein Praefix fuer die Hook-Zeile ist festgelegt, dokumentiert in CLAUDE.md und core/release, und alte Leser zeigen sie als normale Zeile
- [ ] jaira update und das Popup zeigen die Hook-Zeilen zuerst; Test
- [ ] Ein Board, das 3 Versionen zurueckliegt, sieht die Notes aller 3; Test
- [ ] Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-03 17:20 · berkc** — teamlead 03.11.2026: baut auf MXMDFH (Popup, Praefix '- + ') auf, das noch nicht auf master ist. Branch feat/CZWKNY-release-hook startet von origin/feat/MXMDFH-release-notes-popup; der PR geht nach MXMDFH.
- **2026-11-03 17:23 · berkc** — pre-process findings: skipped versions ALREADY work. release.Since(stamp) (core/release/release.go selectSince) returns every entry newer than the stamp; jaira update (internal/cli/update.go writeNotes) prints all of them; the popup (internal/tui/releasenotes.go openReleaseNotes, GK1QGT) pages one version per page through all of them. Only missing: a deterministic test for '3 versions behind' (the existing popup test reads real notes and t.Skips under 3). DoD 3 = tests only, no code.
- **2026-11-03 17:23 · berkc** — decisions settled without asking (disagree at review): (1) hook prefix is '- > ' — any '- X ' line is still one plain change to changeLines, so older readers show '> text'; (2) a hook does not combine with + or !; (3) hooks go first PER VERSION, done once in release.entryOf (stable partition), so update text, update --json and popup all get the order from one place; the popup additionally puts them above the Features header; (4) 'every note file opens with a hook' is documented in CLAUDE.md/AGENTS.md, NOT enforced by a test — a test would fail every in-flight branch's note file; (5) MXMDFH.md (another ticket's unreleased note) is left alone; (6) NOTES.md is closed history, so 'core/release docs' = the package doc of core/release/release.go.
