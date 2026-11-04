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
updated-at: 2026-11-04T00:42:22Z
updated-by: berkc
claimed-by: berk-ca-512097
claimed-at: 2026-11-04T00:37:05Z
outcome-what: "Release notes take a '- > text' hook line; core/release sorts each release's hooks first, and jaira update and the release-notes popup print them first, with '>'."
outcome-why: "Readers skip long notes; a one-line pitch per release makes them open it. Older binaries read the same line as an ordinary change, so no reader breaks."
outcome-resolves: "Hook marker, hook-first ordering, and three-versions-behind paging are covered by tests in core/release, internal/cli and internal/tui; the format is documented in CLAUDE.md and AGENTS.md; the note is core/release/notes/CZWKNY.md."
review-summary: "CLAUDE.md:274 + AGENTS.md:227 make a hook mandatory in every note file; with one file per ticket (0.5.7 = 9 files) that puts N hooks above Features and brings back the drowning the ticket is for. Make it optional ('may open with one') unless the person wants it mandatory. Open question on the ticket. No other finding."
question: ""
round: 2
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
- [ ] Hook-Format '- > <Headline>: <ein kurzer Satz>'; das Popup zeigt die Headline fett, den Satz gedimmt darunter, eine Leerzeile zwischen Hooks; Test
- [ ] Popup: Hooks, Features und '!'-Zeilen immer sichtbar, einfache Zeilen unter 'Others' standardmaessig eingeklappt, Taste o klappt sie aus; jaira update zeigt weiter alles; Test
- [ ] Ein Test schlaegt fehl bei einer Hook-Zeile ueber 80 Zeichen oder mehr als einem Hook pro Note-Datei

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
- [x] R2 popup: notesView.others + key o; lines() buckets hooks/features/important/plain, plain folded to 'n more · o shows them'; hookLines (headline bold, sentence dim); changeLines drops hook case; hint 'o others'
- [x] R2 tui tests: v.key("o") in 4 existing tests; TestNotesViewHookIsAHeadlineWithItsSentenceBelow; TestNotesViewFoldsPlainChangesUnderOthers
- [x] R2 core/release/notes_test.go TestHookLinesAreShortAndOnePerFile over embedded notesFS (<=80 runes, <=1 per file); prove it fails on a long hook
- [x] R2 docs: release.go package doc + Hook doc; CLAUDE.md + AGENTS.md hook bullet (optional headline, max one, <=80, popup folds plain); CZWKNY.md rewritten with a headline hook
- [ ] R2 gofmt/vet/tests; tick DoD 5-7, re-proof DoD 1; move to critique; one commit code+ticket

## Progress
- **2026-11-03 17:20 · berkc** — teamlead 03.11.2026: baut auf MXMDFH (Popup, Praefix '- + ') auf, das noch nicht auf master ist. Branch feat/CZWKNY-release-hook startet von origin/feat/MXMDFH-release-notes-popup; der PR geht nach MXMDFH.
- **2026-11-03 17:23 · berkc** — pre-process findings: skipped versions ALREADY work. release.Since(stamp) (core/release/release.go selectSince) returns every entry newer than the stamp; jaira update (internal/cli/update.go writeNotes) prints all of them; the popup (internal/tui/releasenotes.go openReleaseNotes, GK1QGT) pages one version per page through all of them. Only missing: a deterministic test for '3 versions behind' (the existing popup test reads real notes and t.Skips under 3). DoD 3 = tests only, no code.
- **2026-11-03 17:23 · berkc** — decisions settled without asking (disagree at review): (1) hook prefix is '- > ' — any '- X ' line is still one plain change to changeLines, so older readers show '> text'; (2) a hook does not combine with + or !; (3) hooks go first PER VERSION, done once in release.entryOf (stable partition), so update text, update --json and popup all get the order from one place; the popup additionally puts them above the Features header; (4) 'every note file opens with a hook' is documented in CLAUDE.md/AGENTS.md, NOT enforced by a test — a test would fail every in-flight branch's note file; (5) MXMDFH.md (another ticket's unreleased note) is left alone; (6) NOTES.md is closed history, so 'core/release docs' = the package doc of core/release/release.go.
- **2026-11-03 17:27 · berkc** — DEVIATION (test only): plan's TestPopupThreeVersionsBehindPagesAllThree failed: the popup title reads 'What is new since 0.5.5 · 4 unread', so strings.Contains("0.5.5 · ") matched the title, not a page. Also Since(0.5.5) returns 4 entries here (Unreleased(2) + 0.5.8 + 0.5.7 + 0.5.6), since notes/ holds unreleased files; the test already counts n dynamically. Fix: match the page header in its box, '│ <v> · ' (rendered as '│   │ 0.5.8 · 2/4'). Plan code otherwise kept.
- **2026-11-03 17:28 · berkc** — claim: 'jaira claim CZWKNY' refused (exit 3): claimed by berk-ca-130447 (4m before start). Did not pass --steal; ticket writes (dod, note, move) went through under assignee berkc. Teamlead should confirm the holder is the dispatching session.
- **2026-11-03 17:30 · berkc** — critique (pass 1): one design question, no code finding. CLAUDE.md:274 / AGENTS.md:227 make the hook MANDATORY ('Every note file opens with one'), which the goal asks for, but it works against the context ('big changes drown among many small lines'). Releases hold one file per ticket: 0.5.7 = 9 files, 0.5.5 = 8 files, so the popup and 'jaira update' would print 8-9 bold hooks above Features, the same pile one level up. For a file with a single change (most of them) the hook either repeats that change or is the only line, so most changes end up in the hook block and MXMDFH's Features section empties out. Question: is a hook required in every note file? Options: (1) RECOMMENDED: hook optional, written only for a change worth pitching; change only the two doc bullets to 'A note file may open with one'; code and tests stay as they are. Cost: relies on the author's judgement, nothing enforces it. (2) Keep it mandatory, as the goal says. Cost: N hooks per release, hooks repeat one-line files, and the drowning comes back. (3) One hook per release, written when the release is cut (release-notes workflow / Cut), not per ticket. Cost: new machinery in Cut plus a per-version file, more than this ticket's scope. Left standing, settled by the 17:23 note: prefix '- > ' (1), no combining (2), ordering in entryOf plus the popup's hooks bucket (3), not enforced by test (4). Rest of diff: no finding.
- **2026-11-03 23:05 · berkc** — Maintainer 04.11.2026: 'Hook' heisst eine Ueberschrift neben '- > ', z.B. 'Neuer Skill: /interreview' plus kurze Beschreibung 'Nuetzlich zum Durchgehen von Tickets in review'. Release-News sollen kurz und knackig sein; Zeilen wie die zu 'jaira roles install --project ... --force' (0.5.7/PWAF22) gehoeren nicht in die News. Heute: 0.5.7 hatte 32 Zeilen mit im Schnitt 214 Zeichen - eine Textwand. Format-Entscheidung steht aus.
- **2026-11-04 00:37 · berkc** — Maintainer decision 04.11.2026 (via coordinator): hook = HEADLINE. (1) format '- > <Headline>: <one short sentence>'; popup shows headline bold, sentence dimmed on the line below, blank line between hooks. (2) optional, at most one hook per ticket note file, only for changes a user really notices (new command, skill, view); docs say so. (3) popup: hooks and '!' lines always visible; plain lines under 'Others' collapsed by default, one key expands; jaira update terminal output keeps showing everything. (4) a test fails when a hook line is longer than ~80 chars.
- **2026-11-04 00:39 · berkc** — round 2 plan choices: plain = no marker; '!' non-feature lines stay visible under Others; 'Others' header only below a hook or feature (keeps a plain-only version a plain list); fold state is one bool on the view that survives paging, 'o' toggles; length cap counts the whole line incl. '- > ' in runes; test also enforces at most one hook per file since the maintainer asked for that; jaira update text output untouched.
