---
id: 01M7ZYF17AVJAEWK40GG3DVVFZ
title: "Bei TUI-, Cockpit- oder OS-spezifischen Änderungen schreibt pre-process einen DoD-Punkt 'von einer Person angesehen'"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Was nur ein Mensch an echtem Terminal oder Betriebssystem prüfen kann, steht als eigener DoD-Punkt auf dem Ticket und wird bei der Abnahme gesehen, nicht vergessen."
context: |-
  Was heute fehlte (21.11.):
  - 02D92C: Emoji-Breite in Windows Terminal nur über die Unicode-Breitenklasse belegt, kein Mensch hat hingesehen; der DoD-Punkt galt trotzdem als erledigt.
  - JDDCEZ: Test lief lokal und auf Linux grün, fiel erst auf macOS-CI (kein timeout(1)).
  - AX8K02: der Fall "Person setzt done ohne --force" liess sich vom Agenten nicht spielen, belegt nur über den Refusal-Text.
  Die Lanes behandeln solche Punkte wie alle anderen und haken sie mit einem Ersatzbeweis ab.
  Gewünscht: der shipped pre-process-Lane (core/lane/builtin, Datei für pre-process) schreibt bei Änderungen an TUI, Cockpit, Terminal-Ausgabe oder OS-spezifischem Verhalten einen eigenen DoD-Punkt "von einer Person angesehen: <was, wo>" mit jaira dod --add. Agenten-Lanes haken diesen Punkt nie ab (Lane-Prompts von in-progress/testing/review sagen das), die Abnahme zeigt ihn als offenen Schritt.
  Offen für pre-process: woran der Lane erkennt, dass ein Punkt nur menschlich prüfbar ist (Liste von Pfaden wie core/mod/, internal/tui/ ist der Vorschlag).
definition-of-done: "pre-process-Lane-Prompt verlangt den DoD-Punkt 'von einer Person angesehen' für die genannten Änderungsarten"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-11-21T20:58:45Z
updated-at: 2026-11-21T21:01:56Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2704396
claimed-at: 2026-11-21T20:59:22Z
stop-after-fix: "true"
outcome-what: "Built-in lanes: pre-process adds a DoD item 'seen by a person: <what, where>' for changes to what a terminal draws, platform-dependent behaviour, or person-only actions; in-progress adds it when no planning ran and never ticks it; testing leaves it open (agent-ticked = finding, untick); review makes it a review-check step (agent-ticked = blocker). shipped.sha256 + release note."
outcome-why: "02D92C, JDDCEZ, AX8K02 ticked human-only checks with stand-in proofs (width table, Linux run, refusal text)."
outcome-resolves: "DoD 1-3: 15-pre-process.md:77, 20-in-progress.md:164, 27-testing.md:57, 40-review.md:107, shipped.sha256, core/release/notes/3DVVFZ.md. Categories: prompts only, no code — concurrency/error paths/permissions do not apply; input class = item text, matched by fixed prefix."
---

# Bei TUI-, Cockpit- oder OS-spezifischen Änderungen schreibt pre-process einen DoD-Punkt 'von einer Person angesehen'

## Definition of Done

- [x] pre-process-Lane-Prompt verlangt den DoD-Punkt 'von einer Person angesehen' für die genannten Änderungsarten
  proof: core/lane/builtin/15-pre-process.md:77 '## What only a person can check'
- [x] in-progress-, testing- und review-Prompt sagen: diesen Punkt hakt nur eine Person ab
  proof: core/lane/builtin/20-in-progress.md:164, 27-testing.md:57, 40-review.md:107
- [x] shipped.sha256 und Release-Note core/release/notes/<ID>.md
  proof: core/shipped/shipped.sha256 (+4 versions, scripts/shipped-history.sh --write), core/release/notes/3DVVFZ.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-21 20:59 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock
- **2026-11-21 20:59 · Alexander Sacharov** — Branch feat/3dvvfz was cut from the unpushed F1V59P tip (8ae150955); reset onto upstream/master (012caa9f7) so the PR carries only this ticket.
- **2026-11-21 21:00 · Alexander Sacharov** — decision: no path list. The built-in lanes ship into every repository that runs jaira, so core/mod/ or internal/tui/ would mean nothing elsewhere. pre-process recognises the kind of change instead: what a terminal draws (TUI, colours, widths, glyphs, a plugin/cockpit panel), behaviour that differs by OS/shell/terminal, and anything that needs a real person at a real keyboard. The item text starts with the fixed English prefix 'seen by a person:' because the prompts are English and later lanes match on it.
- **2026-11-21 21:00 · Alexander Sacharov** — decision: in-progress adds the item itself when no planning step ran (this very ticket had no 'planning' option); otherwise the most common path would never get one.
- **2026-11-21 21:01 · Alexander Sacharov** — Left .jaira/lanes/ (this repo's own board) alone: testing.md and review.md are byte-identical to the previous built-in and pre-process/in-progress carry local edits; bringing them along is 'jaira update --lanes', not this ticket.
