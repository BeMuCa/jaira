---
id: 01M801GV3AQGY3YYT3PVHQHFXQ
title: "Cockpit-Overview: die ganze obere Zeile (Tabs und Close) im gerahmten Stil des Close-Knopfs"
status: done
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
commits:
  - fab81a6cc13cf6beb70e51396ad638ae32418798
  - 81e626549955fd0c033ed88f4f16895c9e3a4277
  - 9d963a97cbf8a8684551b19d2931c5035de687ef
created-at: 2026-11-21T21:52:10Z
updated-at: 2026-11-22T12:20:52Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2591038
claimed-at: 2026-11-22T11:57:16Z
outcome-what: "Overview-Pane: Tabs Queue/Log/Workers und Close teilen einen Helfer framed() — jeder Knopf gerahmt (round, paddingX 1), gleich hoch, in einer Zeile (overview-top-row, kein Umbruch: OVERVIEW_ROWS=3); offener Tab Rahmen+Text im Theme-Akzent 'claude' + bold + primary, andere gedimmt. Test ersetzt den alten Close-Test; Release-Note, shipped.sha256."
outcome-why: "Screenshot 21.11.: Tabs als schmaler Text neben dem dreizeiligen gerahmten Close — ungleiche Höhe, unruhige Zeile, kleine Treffflächen."
outcome-resolves: "DoD 1 (Test) und 3; DoD 2 (Windows Terminal ansehen) bleibt für die Person. Geprüft: Nebenläufigkeit n/a (reines Rendern); Fehlerpfad: ungelesenes Board zeigt Close gerahmt (Test); ungewöhnliche Eingabe: schmale Pane — Zeile bricht bewusst nicht um (Workers-Tab rechnet mit 3 Reihen), wird abgeschnitten; Rechte n/a."
review-summary: "overview-pane.tsx: neuer Helfer framed() zeichnet alle vier Knöpfe der oberen Zeile (Queue, Log, Workers, Close) gleich - runder Rahmen, paddingX 1, Label mit zwei Leerzeichen gepolstert, also gleich hoch. Der offene Tab hat Rahmen und fette Schrift in der Theme-Farbe 'claude' plus variant primary; die anderen Tabs sind gedimmt; Close ist neutral fett. Die Zeile bricht nicht um und hat keinen alignItems mehr, damit sie die von register.tsx eingeplanten 3 Reihen (OVERVIEW_ROWS) hält. Test ersetzt den alten Close-Test: vier Rahmen, gleiche Props, eine Zeile ohne flexWrap, nur der offene Tab hervorgehoben. Release-Note und shipped.sha256 ergänzt."
round: 4
back-from: signoff
back-by: person
back-count: 1
review-gaps: |-
  note: schmale Pane (< ~55 Spalten): die Zeile bricht absichtlich nicht um; ob der Host die Knöpfe abschneidet oder ihren Text innerhalb des Rahmens umbricht (dann wird die Zeile doch höher als OVERVIEW_ROWS=3), sieht kein Test - nur der Baum wird geprüft; bei DoD 2 an einer schmalen Pane mitansehen
  note: DoD 2 (Person sieht die Zeile in Windows Terminal an) ist offen - gehört der Person in signoff, kein Code-Mangel
test-verdict: "pass (round 2, changed since last hunt: overview-pane.tsx top row without flexWrap + test asserts no wrap): claude plugin test 296/296 RC=0, shipped-history.sh --check RC=0; DoD 1 test board.test.ts:3402 green; DoD 2 left for the person; writes nothing"
hunt-verdict: "pass: 0 findings at 9d963a97c (round 2 over diff since 81e626549); open for the person: narrow pane (<55 cols) not seen — no wrap by design, host clips"
review-verdict: "accept: DoD 1 und 3 im Diff belegt, 296/296 Cockpit-Tests grün (claude plugin test .), shipped.sha256 enthält den aktuellen Hash von overview-pane.tsx; DoD 2 ist die Sichtprüfung der Person"
review-round: 3
review-verdict-gaps: |-
  note: schmale Pane (< ~55 Spalten): die Zeile bricht absichtlich nicht um; ob der Host die Knöpfe abschneidet oder ihren Text innerhalb des Rahmens umbricht (dann wird die Zeile doch höher als OVERVIEW_ROWS=3), sieht kein Test - nur der Baum wird geprüft; bei DoD 2 an einer schmalen Pane mitansehen
  note: DoD 2 (Person sieht die Zeile in Windows Terminal an) ist offen - gehört der Person in signoff, kein Code-Mangel
review-check: "1. In einem Repo mit jaira-Board Claude Code in Windows Terminal starten (Cockpit dieses Branches: claude --plugin-dir /home/alex/projects/.worktrees/jaira-hqhfxq/core/mod/builtin/jaira-cockpit)  2. Im Cockpit-Band auf Overview drücken  3. Oben stehen Queue, Log, Workers, Close in EINER Reihe, alle im runden Rahmen, alle gleich hoch  4. Queue ist offen: nur sein Rahmen und Text sind in der Akzentfarbe, Log und Workers gedimmt  5. Auf Log, dann Workers klicken: die Hervorhebung wandert mit, Close bleibt neutral  6. Fenster schmal ziehen (< 55 Spalten): die Reihe bricht nicht auf zwei Zeilen um; im Workers-Tab bleibt die Bodenlinie sichtbar  7. Close klicken: die Pane schließt"
question: ""
---

# Cockpit-Overview: die ganze obere Zeile (Tabs und Close) im gerahmten Stil des Close-Knopfs

## Definition of Done

- [x] Tabs und Close stehen in einer Zeile, gleich gerahmt und gleich hoch; der aktive Tab ist hervorgehoben (Test)
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:201 framed(); test 'the top row is one row of framed buttons alike' tests/board.test.ts:3402
- [x] von einer Person angesehen: obere Zeile in Windows Terminal
  proof: Alex 22.11. angesehen
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
- **2026-11-22 11:22 · Alexander Sacharov** — bughunt: done at 9d963a97c - 0 findings in 2 modules (diff since 81e626549: overview-pane.tsx top row lost flexWrap, board.test.ts asserts it; shipped.sha256 carries the new hash; lanes/AGENTS/CLAUDE diffs are jaira update text, no code), follow-up hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/HQHFXQ.md core/shipped/shipped.sha256 .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/optimize.md .jaira/lanes/pre-process.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md
- **2026-11-22 11:22 · Alexander Sacharov** — für die Person bei DoD 2: auch eine schmale Pane (<55 Spalten, Handy) ansehen — die Zeile bricht absichtlich nicht um (OVERVIEW_ROWS=3), der Host schneidet; Tests sehen nur den Baum, nicht die Breite.
- **2026-11-22 11:25 · Alexander Sacharov** — review: verdict accept (296/296 cockpit tests). Move nach signoff verweigert: Gate "bug hunt covers HEAD" sieht uncommittete Änderungen in .jaira/lanes/*.md, AGENTS.md, CLAUDE.md (Regeneration durch jaira update, neuerer Lane-Text, kein Code dieses Tickets). Frage: was mit diesen Dateien?
1. (empfohlen) Person verschiebt mit --force nach signoff: ! jaira move HQHFXQ --to signoff --force — die Regeneration bleibt im Worktree liegen und kommt nicht in HQHFXQ-Commits
2. Regeneration verwerfen (git checkout -- .jaira/lanes AGENTS.md CLAUDE.md), dann normaler move — wie am 22.11. 11:06 schon einmal gemacht; verliert den neueren Lane-Text, bis jaira update wieder läuft
3. Regeneration als eigenen chore-Commit auf diesem Branch — vermischt jaira update mit HQHFXQ im PR
- **2026-11-22 11:26 · Alexander Sacharov** — decision (Worker): Option 2 — jaira-update-Regeneration in .jaira/lanes, AGENTS.md, CLAUDE.md verworfen wie um 11:06; kein Code dieses Tickets, kommt mit dem nächsten jaira update auf master. Kein --force nötig.
- **2026-11-22 11:56 · Alexander Sacharov** — Given back from signoff: beheben:
dazu: SKILL пишет «при красном CI ничего не запушено», но merge-коммит к этому моменту уже │ неверная строка в │ │ запуше
angenommen:
- note: schmale Pane (< ~55 Spalten): die Zeile bricht absichtlich nicht um; ob der Host die Knöpfe abschneidet oder ihren Text innerhalb des Rahmens umbricht (dann wird die Zeile doch höher als OVERVIEW_ROWS=3), sieht kein Test - nur der Baum wird geprüft; bei DoD 2 an einer schmalen Pane mitansehen
- note: DoD 2 (Person sieht die Zeile in Windows Terminal an) ist offen - gehört der Person in signoff, kein Code-Mangel
- **2026-11-22 11:57 · Alexander Sacharov** — cockpit: given back to a new worker - nothing recorded about the old worker session (cache and context unknown)
- **2026-11-22 11:57 · Alexander Sacharov** — Worker: Give-back-Grund 11:56 (SKILL schreibt 'bei rotem CI nichts gepusht', Merge-Commit aber schon gepusht) passt nicht zu diesem Branch — HQHFXQ ändert keine SKILL.md, nur overview-pane.tsx/Test/Notes. Vermutlich Text eines anderen Tickets. Frage an die Person: 1. (empfohlen) Irrtum — HQHFXQ unverändert zurück nach signoff 2. Fund gehört doch hierher — welche SKILL-Datei/Zeile?
- **2026-11-22 12:01 · Alexander Sacharov** — decision (Person, 22.11.): Give-back 11:56 war ein Irrtum (Text eines anderen Tickets). HQHFXQ unverändert zurück nach signoff; Review-Urteil accept von 11:25 gilt weiter.
