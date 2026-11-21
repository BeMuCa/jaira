---
id: 01M7ZPP39H62XRS3R59602D92C
title: "Cockpit-Log: Filter nach Lane, größerer Close-Knopf, ein Zeichen pro Ereignis"
status: critique
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
commits: []
created-at: 2026-11-21T18:42:48Z
updated-at: 2026-11-21T19:11:40Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2284465
claimed-at: 2026-11-21T19:01:26Z
outcome-what: "Overview Log: lane filter row (board order, several, AND with people), Close framed in the tab row, 🆕/🔙 sign first on every Log line; manifest + release note"
outcome-why: "Alex wants to see only e.g. blocked or signoff, hit Close without scrolling, and read an event's kind at a glance"
outcome-resolves: Log filters only by person; Close is a small word under long lists; created/put back differ only by colour
review-summary: "core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:180 filterRow picks its button-key prefix by testing key === 'people' — hidden coupling on a string; pass the prefix ('overview-person' / 'overview-lane') as the parameter instead of deriving it"
round: 2
back-from: critique
back-by: lane
back-count: 1
---

# Cockpit-Log: Filter nach Lane, größerer Close-Knopf, ein Zeichen pro Ereignis

## Definition of Done

- [x] Log-Tab hat eine Reihe Lane-Buttons; Auswahl mehrerer Lanes zeigt nur deren Ereignisse, zusammen mit dem Personenfilter (Test)
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:192; test 'the Log filters to the lanes chosen, in board order, several at once, together with the people filter' (core/mod/builtin/jaira-cockpit/tests/board.test.ts:2754)
- [x] Close steht oben neben den Tabs und ist sichtbar größer als ein Tab-Button
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:131 (round frame, paddingX 1, in the tab row); test 'Close stands in the tab row, framed and wider than a tab' (core/mod/builtin/jaira-cockpit/tests/board.test.ts:2790)
- [x] Jede Log-Zeile beginnt mit dem Zeichen ihrer Ereignisart; Spalten verrutschen in Windows Terminal nicht (geprüft)
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:200 (sign first; only single Emoji_Presentation code points, East Asian Wide, no VS16/ZWJ); test 'every Log line starts with its kind's sign' (core/mod/builtin/jaira-cockpit/tests/board.test.ts:2813) holds the width class and the clock column. Visual check in a real Windows Terminal is the acceptance step
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
