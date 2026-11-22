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
updated-at: 2026-11-22T11:21:05Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-30716
claimed-at: 2026-11-22T11:07:59Z
outcome-what: "Overview-Pane: Tabs Queue/Log/Workers und Close teilen einen Helfer framed() — jeder Knopf gerahmt (round, paddingX 1), gleich hoch, in einer Zeile (overview-top-row, kein Umbruch: OVERVIEW_ROWS=3); offener Tab Rahmen+Text im Theme-Akzent 'claude' + bold + primary, andere gedimmt. Test ersetzt den alten Close-Test; Release-Note, shipped.sha256."
outcome-why: "Screenshot 21.11.: Tabs als schmaler Text neben dem dreizeiligen gerahmten Close — ungleiche Höhe, unruhige Zeile, kleine Treffflächen."
outcome-resolves: "DoD 1 (Test) und 3; DoD 2 (Windows Terminal ansehen) bleibt für die Person. Geprüft: Nebenläufigkeit n/a (reines Rendern); Fehlerpfad: ungelesenes Board zeigt Close gerahmt (Test); ungewöhnliche Eingabe: schmale Pane — Zeile bricht bewusst nicht um (Workers-Tab rechnet mit 3 Reihen), wird abgeschnitten; Rechte n/a."
review-summary: "none — round 2: the one finding (fixed 'cyan') is answered in 81e626549 with the theme key 'claude', test follows"
round: 3
back-from: testing
back-by: lane
back-count: 1
review-gaps: "none — framed() replaces both the old tabButton and closeButton bodies (no second copy left); the framed Box in register.tsx:1587 is a column panel, not a button, so not the same idea; closeButton kept as a name: used twice (unread board and top row)"
test-verdict: "pass: claude plugin test (cockpit) 296/296 RC=0, go build ./... RC=0, go test core/shipped core/release core/mod RC=0, shipped-history.sh --check RC=0; DoD 1 verified by test 'the top row is one row of framed buttons alike' (board.test.ts:3402: four frames same border/padding, one row in order, open tab only one with accent frame + primary, on every tab); DoD 3 note file + manifest line present; failure path: unreadable board still shows Close framed (same test); writes nothing; DoD 2 (Windows Terminal look) left for the person"
hunt-verdict: "change: overview-pane.tsx:309 flexWrap lets the top row wrap on a narrow pane while register.tsx:104 OVERVIEW_ROWS=3 budgets 3 rows — Workers tab spills; fix: no wrap"
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
- **2026-11-22 11:19 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:309 - class: the top row's height is assumed to be a constant (OVERVIEW_ROWS = 3 at hooks/register.tsx:104, subtracted from bodyRows at register.tsx:1914 to size the Workers tab), but the row now wraps (flexWrap="wrap", added by fab81a6cc) and is about 55 cells wide (four framed buttons of 13+11+15+13 plus 3 gaps); example: a pane narrower than ~55 columns (a phone, a narrow dock) wraps Close, or Close and tabs, onto a second or third 3-row line, so the top row is 6-9 rows high while the Workers picture plus one line per agent is still sized for 3 - the ground line and the bottom agent lines fall below the pane and have to be scrolled to, which the Workers tab promises they never do (overview-pane.tsx:283-285). Before the change the row had no wrap and stayed 3 high. Fix either the wrap or the subtraction (measure, or derive OVERVIEW_ROWS from whether the row fits bodyColumns). Not a blocker: layout only, content reachable by scrolling. ran: nothing - the test harness checks the element tree, not the layout; confirmed by read and width arithmetic
- **2026-11-22 11:19 · Alexander Sacharov** — bughunt: done at 81e626549 - 1 finding (1 change, 0 pre-existing) in 2 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/HQHFXQ.md core/shipped/shipped.sha256 .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/optimize.md .jaira/lanes/pre-process.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md
- **2026-11-22 11:19 · Alexander Sacharov** — testing: Bug-Hunt-Fund übernommen — flexWrap auf der Tab-Zeile widerspricht OVERVIEW_ROWS=3 (register.tsx:104). Ohne wrap bleibt die Zeile drei Reihen hoch, wie vor diesem Ticket; schmale Pane schneidet ab statt umzubrechen. Theme-Key 'claude' für borderColor ist gültig: Color = ThemeKey | string im Host-Typ.
