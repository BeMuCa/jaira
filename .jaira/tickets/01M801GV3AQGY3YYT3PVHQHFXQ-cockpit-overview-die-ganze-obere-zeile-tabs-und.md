---
id: 01M801GV3AQGY3YYT3PVHQHFXQ
title: "Cockpit-Overview: die ganze obere Zeile (Tabs und Close) im gerahmten Stil des Close-Knopfs"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die obere Zeile der Overview-Pane ist eine ruhige Reihe gleich aussehender, gut treffbarer Knöpfe."
context: |-
  Screenshot Alex 21.11. abends (0.6.0): oben links Tabs [ Queue ] [ Log ] [ Workers ] als schmaler Text, daneben der neue gerahmte, zweizeilige Close-Knopf aus 02D92C — er steht höher/tiefer als die Tabs und sieht anders aus. Oben rechts zeigt Claude Code zusätzlich sein eigenes ×, also gibt es zwei Schliessen-Knöpfe.
  Gewünscht (Alex): die ganze obere Zeile im Stil von Close umschreiben — Tabs genauso gerahmt und gleich hoch, in einer Reihe mit Close; der aktive Tab hervorgehoben.
  Offen für pre-process: ob Close neben dem × des Hosts überhaupt bleiben soll (Vorschlag: bleibt, weil das × klein ist und auf dem Handy fehlt).
  Code: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx (tabButton, closeButton in drawOverview).
definition-of-done: "Tabs und Close stehen in einer Zeile, gleich gerahmt und gleich hoch; der aktive Tab ist hervorgehoben (Test)"
tags:
  - tui
blocked-by: []
related: []
follows: 01M7ZPP39H62XRS3R59602D92C
commits: []
created-at: 2026-11-21T21:52:10Z
updated-at: 2026-11-22T11:15:42Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-30716
claimed-at: 2026-11-22T11:07:59Z
outcome-what: "Overview-Pane: Tabs Queue/Log/Workers und Close teilen einen Helfer framed() — jeder Knopf gerahmt (round, paddingX 1), gleich hoch, in einer Zeile (overview-top-row, wrap bei schmalem Terminal); aktiver Tab cyan Rahmen + bold + primary, inaktive gedimmt. Test ersetzt den alten Close-Test; Release-Note, shipped.sha256."
outcome-why: "Screenshot 21.11.: Tabs als schmaler Text neben dem dreizeiligen gerahmten Close — ungleiche Höhe, unruhige Zeile, kleine Treffflächen."
outcome-resolves: "DoD 1 (Test) und 3; DoD 2 (Windows Terminal ansehen) bleibt für die Person. Geprüft: Nebenläufigkeit n/a (reines Rendern); Fehlerpfad: ungelesenes Board zeigt Close gerahmt (Test); ungewöhnliche Eingabe: schmales Terminal — Zeile bricht um statt abzuschneiden; Rechte n/a."
review-summary: "core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:202,204 hardcodes 'cyan' for the open tab; the file colours by theme key (GROUP, line 15-16: 'they follow the person's theme') — use the theme key 'claude' instead, and the test with it"
round: 2
back-from: critique
back-by: lane
back-count: 1
---

# Cockpit-Overview: die ganze obere Zeile (Tabs und Close) im gerahmten Stil des Close-Knopfs

## Definition of Done

- [x] Tabs und Close stehen in einer Zeile, gleich gerahmt und gleich hoch; der aktive Tab ist hervorgehoben (Test)
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:201 framed(); test 'the top row is one row of framed buttons alike' tests/board.test.ts:3402
- [ ] von einer Person angesehen: obere Zeile in Windows Terminal
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/HQHFXQ.md; core/shipped/shipped.sha256 (+1 overview-pane.tsx), shipped-history.sh --check ok

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 11:05 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2784.sock
- **2026-11-22 11:06 · Alexander Sacharov** — Worktree auf upstream/master (0.6.1) vorgespult; lokale Regeneration der Lane-Dateien durch Binary 0.6.0 verworfen (nicht Teil des Tickets). Options: planning nicht angehakt, pre-process entfällt.
- **2026-11-22 11:08 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock (neuer Worker; voriger Claim-Prozess 11624 tot, Claim übernommen)
- **2026-11-22 11:12 · Alexander Sacharov** — decision: Close bleibt neben dem × des Hosts (Vorschlag aus context übernommen, Alex wollte 'in einer Reihe mit Close'). Alle vier Knöpfe teilen einen Helfer framed() — Tabs und Close können so nicht mehr auseinanderlaufen. Aktiver Tab: cyan Rahmen + bold + variant primary; inaktive Tabs: borderDimColor + dimColor. Close bleibt key overview-close-frame.
- **2026-11-22 11:14 · Alexander Sacharov** — critique: overview-pane.tsx framed() färbt den offenen Tab mit festem 'cyan'; die Datei nutzt Theme-Keys (GROUP) — auf 'claude' umstellen, Test mit.
