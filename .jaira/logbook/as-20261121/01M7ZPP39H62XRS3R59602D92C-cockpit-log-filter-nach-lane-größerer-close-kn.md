---
id: 01M7ZPP39H62XRS3R59602D92C
title: "Cockpit-Log: Filter nach Lane, größerer Close-Knopf, ein Zeichen pro Ereignis"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Im Log-Tab des Cockpits lässt sich nach Lane filtern, Close ist leicht zu treffen, und jede Log-Zeile zeigt ihr Ereignis auf einen Blick."
context: |-
  Was heute fehlt, alles in core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx (drawOverview):
  1. Der Log-Tab filtert nur nach Personen (state.people, Buttons über der Historie). Nach Lane filtern geht nicht — Alex will z. B. nur blocked oder nur signoff sehen.
  2. Der Close-Knopf (closeButton, Label w.close aus hooks/texts.ts) ist ein kleines Wort ganz unten, unter langen Listen; schwer zu treffen.
  3. Log-Zeilen sind nur Text; created und put back unterscheiden sich nur durch Farbe.
  Entschieden von Alex am 2026-11-21:
  - Filter = die Lane, in der der Ticket jetzt steht (ev.lane), mehrere wählbar, kombiniert mit dem Personenfilter wie der bestehende (leer = alle).
  - Close: größer (Rahmen/Padding) und nach oben neben die Tabs.
  - Ein Emoji pro Ereignisart in der Log-Zeile (z. B. 🆕 created, ↩️ put back), nicht im ganzen Cockpit.
  Risiko: Emoji-Breite im Terminal ist uneinheitlich (Windows Terminal, WSL); Zeilen dürfen nicht verrutschen — prüfen, sonst ein schmales Zeichen statt Emoji.
  Tests: tests/ im Cockpit-Ordner; live prüfen mit claude --plugin-dir <worktree>/core/mod/builtin/jaira-cockpit.
definition-of-done: "Log-Tab hat eine Reihe Lane-Buttons; Auswahl mehrerer Lanes zeigt nur deren Ereignisse, zusammen mit dem Personenfilter (Test)"
tags:
  - tui
blocked-by: []
related: []
commits:
  - 87ada520e40eba11e9f84cd487edb6dfdbe32525
  - 82cce911b9fb7502666d52a9c385b50f0fd0f251
  - 3c905a44a0c4d04ffa2203fe17a6d2f9ed8eaa80
  - 1b3db6c523106cc1de6e034aaea04bf2d6abd9c8
created-at: 2026-11-21T18:42:48Z
updated-at: 2026-11-21T19:45:30Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3317363
claimed-at: 2026-11-21T19:18:40Z
outcome-what: "Overview Log: lane filter row (board order, several, AND with people), Close framed in the tab row, 🆕/🔙 sign first on every Log line; a chosen lane or person without events keeps its button; manifest + release note"
outcome-why: "Alex wants to see only e.g. blocked or signoff, hit Close without scrolling, and read an event's kind at a glance; a filter that hides its own button cannot be undone"
outcome-resolves: "Log filters only by person; Close is a small word under long lists; created/put back differ only by colour; bughunt 19:15 stale filter value without a button"
review-summary: "Im Overview-Log des Cockpits gibt es unter der Personenreihe jetzt eine zweite Knopfreihe mit Lanes (eigenes 'All', Reihenfolge wie auf dem Board, unbekannte Lanes alphabetisch hinten); beide Filter gelten zusammen und filtern nach der Lane, in der ein Ticket jetzt steht. Ein gewählter Wert ohne Ereignis behält seinen Knopf, damit der Filter aufhebbar bleibt. Overview liefert dafür neu 'lanes' (Board-Reihenfolge), der Zustand liegt im neuen Atom overviewLanes, Personen und Lanes teilen sich eine Toggle-Funktion und eine filterRow. Close steht jetzt in der Tab-Reihe oben, in einem runden Rahmen mit gepolstertem Label (Button kennt keine Größe), OVERVIEW_ROWS 2→3 passt die Workers-Höhe an. Jede Log-Zeile beginnt mit 🆕 (angelegt) oder 🔙 (zurückgelegt) — einzelne Emoji_Presentation-Codepoints ohne VS16, damit die Spalten nicht verrutschen. Release-Note und shipped.sha256 sind dabei."
round: 3
back-from: signoff
back-by: person
back-count: 1
review-gaps: |-
  note: DoD 3 'geprüft' in Windows Terminal ist nur über die Breitenklasse (East_Asian_Width=W) belegt, nicht mit dem Auge — ein Blick in WT gehört in die Abnahme
  note: der Lane-Filter filtert nach der aktuellen Lane eines Tickets, nicht nach der Lane zum Zeitpunkt des Ereignisses — so gewollt (Release-Note sagt es nicht ausdrücklich), aber ein 'created' von gestern erscheint heute unter z.B. review
test-verdict: "pass: round after fix, diff since 3c905a44 (overview-pane.tsx, board.test.ts, shipped.sha256): claude plugin test 231/231 RC=0 at 1b3db6c5; go test role/shipped/release RC=0; shipped-history --check RC=0. DoD 1-4 proofs re-pointed and opened. Function: new test chooses signoff+Alex, ticket moves to done and to Berk, reopened Log shows none with both buttons pressed, pressing them brings the event back and the buttons go; same test fails on the pre-fix code (hunk reverted). Writes: only session atoms, no file — no concurrent-writer case."
hunt-verdict: "pass: bughunt done at 1b3db6c5, 0 change, 0 pre-existing (follow-up hunt over the diff since 3c905a44)"
review-verdict: "Der Diff erfüllt alle vier DoD-Punkte; claude plugin test . läuft grün (231 pass), shipped.sha256 enthält die Hashes aller vier ausgelieferten Dateien, Schlüssel der Knopfreihen kollidieren nicht. Keine Defekte gefunden; einzige Unsicherheit ist die WT-Darstellung der Emojis, die nur ein Mensch am echten Terminal bestätigen kann."
review-round: 3
review-check: "1. claude --plugin-dir /home/alex/projects/.worktrees/jaira-02d92c/core/mod/builtin/jaira-cockpit in einem Repo mit .jaira starten (am besten in Windows Terminal)  2. Cockpit-Overview öffnen: oben stehen Queue, Log, Workers und daneben ein gerahmter, breiterer 'Close'  3. Tab Log öffnen: unter der Personenreihe eine Reihe Lane-Knöpfe in Board-Reihenfolge, beginnend mit 'All'  4. Eine Lane drücken, dann eine zweite: nur Ereignisse dieser Lanes bleiben; zusätzlich eine Person drücken: beide Filter gelten zusammen  5. Jede Zeile beginnt mit 🆕 oder 🔙, und die Uhrzeitspalte steht in allen Zeilen bündig  6. Close drücken: der Pane geht zu"
review-verdict-gaps: |-
  note: DoD 3 'geprüft' in Windows Terminal ist nur über die Breitenklasse (East_Asian_Width=W) belegt, nicht mit dem Auge — ein Blick in WT gehört in die Abnahme
  note: der Lane-Filter filtert nach der aktuellen Lane eines Tickets, nicht nach der Lane zum Zeitpunkt des Ereignisses — so gewollt (Release-Note sagt es nicht ausdrücklich), aber ein 'created' von gestern erscheint heute unter z.B. review
---

# Cockpit-Log: Filter nach Lane, größerer Close-Knopf, ein Zeichen pro Ereignis

## Definition of Done

- [x] Log-Tab hat eine Reihe Lane-Buttons; Auswahl mehrerer Lanes zeigt nur deren Ereignisse, zusammen mit dem Personenfilter (Test)
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:174 (buttons = values present ∪ values chosen), :180 (filter); tests 'the Log filters to the lanes chosen, in board order, several at once, together with the people filter' (tests/board.test.ts:2754) and 'a lane or person chosen with no event left keeps its button, and pressing it clears the filter' (tests/board.test.ts:2790)
- [x] Close steht oben neben den Tabs und ist sichtbar größer als ein Tab-Button
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:131 (round frame, paddingX 1, in the tab row); test 'Close stands in the tab row, framed and wider than a tab' (core/mod/builtin/jaira-cockpit/tests/board.test.ts:2823)
- [x] Jede Log-Zeile beginnt mit dem Zeichen ihrer Ereignisart; Spalten verrutschen in Windows Terminal nicht (geprüft)
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:203 (sign first; only single Emoji_Presentation code points, East Asian Wide, no VS16/ZWJ); test 'every Log line starts with its kind's sign' (core/mod/builtin/jaira-cockpit/tests/board.test.ts:2846) holds the width class and the clock column. Visual check in a real Windows Terminal is the acceptance step
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/02D92C.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-21 18:58 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/4205.sock
- **2026-11-21 19:01 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock (replaces the earlier one; previous claimant process was dead)
- **2026-11-21 19:04 · Alexander Sacharov** — decision: Close — Button has no size prop in the engine (ButtonProps, 2.1.330), so Close sits in a round-bordered Box beside the tabs with a padded label; key overview-close and label 'Close' stay. Lane filter: buttons in board order (Overview gets the board's lane order), only lanes with events shown, own 'All', keys overview-lane-<lane>/overview-lanes-all, kept in an atom like people. Glyphs: only single code points with Emoji_Presentation (East Asian Wide, no VS16) — 🆕 created, 🔙 put back — because ↩️ (U+21A9+FE0F) is exactly the ambiguous-width class that shifts columns in Windows Terminal; a test holds every kind's glyph to that class.
- **2026-11-21 19:09 · Alexander Sacharov** — Windows Terminal: no WT screen reachable from this WSL session (WT_SESSION unset), so the 'geprüft' is by width class, not by eye: 🆕 U+1F195 and 🔙 U+1F519 are East_Asian_Width=W with default emoji presentation — WT, conhost-free WSL terminals and Claude Code's string-width all count them 2 cells. The ambiguous class (↩️ = U+21A9+FE0F, ⏪ etc. with VS16) was avoided on purpose. A person should glance at the Log tab in WT at acceptance.
- **2026-11-21 19:10 · Alexander Sacharov** — critique round 1: overview-pane.tsx filterRow derives the per-button key prefix from key === 'people'; pass the prefix in. Let stand: two 'All' buttons (one per row, distinct keys), view.lanes ?? [] fallback for an Overview kept in $.state from before this change (hot reload keeps state), lanes not on the board sorted last.
- **2026-11-21 19:11 · Alexander Sacharov** — critique round 2: the round-1 finding (key prefix) is answered in 82cce911; nothing else raised.
- **2026-11-21 19:15 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:176 - class: a Log filter keeps a chosen value that has no button any more. The buttons are built only from the values present in today's events (laneIds :176, names :171), but state.lanes / state.people (overviewLanesAtom / overviewPeopleAtom, kept across reopenings) are applied unchanged at :178. So the Log shows 'none' with no pressed button to undo, only a dim 'All'. For lanes it happens all the time, because the filter uses the lane a ticket stands in now. The people half of the same pattern (:171, a chosen person whose events age out of the 7 days) is older (9cab6144). Example: choose lane signoff, the ticket moves on to done, open Overview again: the Log shows 'none' though it has 1 event, lane buttons are [done, backlog], and 'All' is not primary. ran: scratch test in a copy of the mod, claude plugin test (handles [], NONE true, signoff button gone). Not a blocker: 'All' clears it
- **2026-11-21 19:15 · Alexander Sacharov** — bughunt: done at 3c905a44 - 1 finding (1 change, 0 pre-existing) in 4 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx core/mod/builtin/jaira-cockpit/hooks/overview.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/types/index.d.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/02D92C.md core/shipped/shipped.sha256
- **2026-11-21 19:18 · Alexander Sacharov** — Given back from signoff: Filterknöpfe = vorhandene Werte ∪ gewählte Werte, für Lanes und Personen, damit ein gewählter Wert ohne Ereignis abwählbar bleibt (Bughunt 19:15). Ein Test: signoff wählen, Ticket geht nach done, Knopf signoff ist noch da und hebt den Filter auf.
- **2026-11-21 19:22 · Alexander Sacharov** — give-back fixed: names/present now take state.people/state.lanes in too (overview-pane.tsx:174-175), so a chosen value with no event keeps a pressed button; once unchosen with no event it goes. New test fails on the old code (checked by reverting the hunk), 232/232 pass; shipped.sha256 regenerated. Release note unchanged: the existing line already promises the lane filter, this only makes it undoable.
- **2026-11-21 19:22 · Alexander Sacharov** — critique round 3 (after give-back): read only the bughunt 19:15 finding and 1b3db6c5; the union of present and chosen values at overview-pane.tsx:174-175 answers it in place, same pattern for both rows; nothing raised.
- **2026-11-21 19:24 · Alexander Sacharov** — bughunt: done at 1b3db6c5 - 0 findings in 3 changed files, follow-up hunt since 3c905a44 (diff only: overview-pane.tsx filter buttons, board.test.ts, shipped.sha256); concurrency, error paths, input, permissions checked; ran: claude plugin test . (231 pass); files: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx core/mod/builtin/jaira-cockpit/hooks/overview.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/types/index.d.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/02D92C.md core/shipped/shipped.sha256
