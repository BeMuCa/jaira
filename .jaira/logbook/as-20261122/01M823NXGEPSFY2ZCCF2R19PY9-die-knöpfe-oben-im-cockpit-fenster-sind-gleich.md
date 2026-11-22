---
id: 01M823NXGEPSFY2ZCCF2R19PY9
title: "Die Knöpfe oben im Cockpit-Fenster sind gleich groß und gut klickbar: Tabs, Close und das Kreuz"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Alle Knöpfe im Kopf eines Cockpit-Fensters (Tabs wie Queue/Log/Workers, Close, das Kreuz oben rechts) haben eine Form und eine Höhe und eine Klickfläche, die man ohne Zielen trifft."
context: |-
  Was falsch ist (Screenshot Alex 22.11.2026, Fenster 'Board overview'): die Knöpfe oben sehen dreimal verschieden aus.
  - Tabs '[ Queue ] [ Log ] [ Workers ]': Text mit Klammern, eine Zeile.
  - 'Close': ein gerahmter Kasten, drei Zeilen hoch, viel größer als alles andere, steht mitten in der Tab-Zeile.
  - Das Kreuz '×' oben rechts: ein einzelnes Zeichen, schwer zu treffen.
  Vorgeschichte: Alex wollte das Kreuz sichtbarer und größer; danach wurde Close zum großen Rahmen (register.tsx:103 'three high since the framed Close stands in it').
  Code: core/mod/builtin/jaira-cockpit/hooks/register.tsx (Kopfzeile der Panes, Close um Zeile 103 und 561).
definition-of-done: "Tabs, Close und Kreuz haben dieselbe Form und Höhe (eine Zeile, '[ ... ]'-Stil oder ein gemeinsamer Stil); Close steht nicht mehr als dreizeiliger Rahmen in der Tab-Zeile"
tags:
  - cockpit
blocked-by: []
related: []
commits:
  - fce29461ae328a27d8fea01e657f2eedde5d8b39
  - b0d5faff44a67ca7eba2cb01fef8a7eac2b8c9cd
  - 05d850d81a56e2933bf363938598580a309fdcd5
created-at: 2026-11-22T17:08:22Z
updated-at: 2026-11-22T18:21:54Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2402047
claimed-at: 2026-11-22T17:08:58Z
outcome-what: "Pane heads on one line with [ Close ] at the right, also in the 'nothing recorded yet' fallback of Verdict, Give back, Explanation and the ticket pane; HQHFXQ release note rewritten to the final state"
outcome-why: "Alex's screenshot: three different button shapes; teamlead give-back: Close everywhere, no contradicting release notes"
outcome-resolves: "R19PY9"
review-summary: "Ein Helfer headRow (overview-pane.tsx) zeichnet den Kopf jedes Cockpit-Panes als eine Zeile: Tabs links als normale Buttons '[ label ]', [ Close ] (role=dismiss) rechts, nowrap, beide Teile flexShrink 0. Die dreizeiligen Rahmen von HQHFXQ sind weg (OVERVIEW_ROWS 3 -> 1, Workers-Tab bekommt 2 Zeilen zurueck). Verdict, Explanation und Message haben Close jetzt oben statt unten; Give back und das Ticket-Pane bekommen erstmals ein Close. Der Fallback 'nothing recorded yet' von Verdict, Give back, Explanation und Ticket-Pane (nothingPane) behaelt denselben Kopf. Das x des Hosts bleibt (Alex 22.11.). Release-Note R19PY9.md neu, HQHFXQ.md auf den Endstand umgeschrieben, shipped.sha256 traegt beide aktuellen Hashes."
round: 3
back-from: signoff
back-by: person
back-count: 1
review-gaps: |-
  note: der 80-Spalten-Test rechnet die Breite aus den Labels nach statt zu rendern; reicht fuer die festen englischen Labels
  note: Close-Schluessel und -Handler stehen in Verdict, Explanation und Ticket-Pane doppelt (Fallback und Kopf); eine Konstante pro Pane wuerde je eine Zeile sparen
test-verdict: "pass: round 3 on the delta since b0d5faff4: claude plugin test RC=0 (316 pass), go test release/shipped/role RC=0, shipped-history --check RC=0; give-back items verified: HQHFXQ.md now one line on the final state, 'nothing recorded yet' fallback of Verdict, Give back, Explanation, ticket pane carries the head and its Close closes the pane (test 'a pane left with nothing to show keeps its head'); writes nothing; DoD 5 left for Alex"
hunt-verdict: "pass: bughunt done at 05d850d81, 0 findings on the delta since b0d5faff4"
review-verdict: "accept: DoD 1-4 erfuellt, auch die Rueckgabe-Punkte (HQHFXQ-Note, Close im leeren Fallback); claude plugin test: 316 pass; shipped.sha256 enthaelt die aktuellen Hashes von register.tsx und overview-pane.tsx; kein Defekt im geaenderten Code gefunden"
review-round: 3
review-verdict-gaps: |-
  note: der 80-Spalten-Test rechnet die Breite aus den Labels nach statt zu rendern; reicht fuer die festen englischen Labels
  note: Close-Schluessel und -Handler stehen in Verdict, Explanation und Ticket-Pane doppelt (Fallback und Kopf); eine Konstante pro Pane wuerde je eine Zeile sparen
review-check: "1. cd core/mod/builtin/jaira-cockpit && claude plugin test .  -> 316 pass, 0 fail  2. claude --plugin-dir $PWD in einem 80 Spalten breiten Terminal starten  3. /cockpit, im Band 'Overview' druecken  4. sehen: eine Zeile '[ Queue ] [ Log ] [ Workers ]' links, '[ Close ]' rechts, keine Rahmen, offener Tab hervorgehoben  5. jeden Tab anklicken - er wird aktiv  6. seen by a person: Tabs, [ Close ] und das x oben rechts anklicken - [ Close ] und x schliessen das Pane  7. bei einem Ticket in signoff Verdict, Give back und Explain oeffnen - in jedem steht [ Close ] rechts in der ersten Zeile und schliesst es  8. ein Verdict-Pane offen lassen, das Ticket aus einer anderen Session annehmen - das Pane zeigt 'nothing recorded yet' und hat weiter [ Close ] oben rechts"
---

# Die Knöpfe oben im Cockpit-Fenster sind gleich groß und gut klickbar: Tabs, Close und das Kreuz

## Definition of Done

- [x] Tabs, Close und Kreuz haben dieselbe Form und Höhe (eine Zeile, '[ ... ]'-Stil oder ein gemeinsamer Stil); Close steht nicht mehr als dreizeiliger Rahmen in der Tab-Zeile
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:198 headRow; test 'the top row is one line of buttons alike' (core/mod/builtin/jaira-cockpit/tests/board.test.ts:3518)
- [x] Jeder Kopf-Knopf (Tabs, Close) hat eine Klickfläche über den ganzen Klammertext, mindestens 3 Zeichen; das × zeichnet Claude Code und bleibt, Close ist die grosse Schliessen-Fläche in jedem Pane mit Kopf (Test über die Layout-Breiten) - Alex 22.11.
  proof: test 'the head of every pane' (core/mod/builtin/jaira-cockpit/tests/board.test.ts:1888): every button drawn '[ label ]' >= 3 columns, Close at the right; × stays the host's (Alex 22.11.)
- [x] Gilt in jedem Pane mit Kopf (Board overview, Verdict, Give back, Report); schmale Terminals (80 Spalten) brechen die Zeile nicht kaputt (Test)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts:1914 (Overview head <= 80 columns, nowrap) and core/mod/builtin/jaira-cockpit/tests/board.test.ts:1925 (Verdict, Give back, Explanation, Message, ticket pane)
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/R19PY9.md; core/shipped/shipped.sha256 regenerated (scripts/shipped-history.sh --write)
- [x] seen by a person: Alex klickt Tabs, Close und Kreuz im laufenden Cockpit
  proof: Alex 22.11. geklickt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 17:09 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 17:11 · Alexander Sacharov** — Befund vor dem Plan: (1) Das Kreuz × oben rechts zeichnet Claude Code selbst (Pane-Rahmen des Hosts), nicht der Cockpit-Mod. $.ui.open kennt nur id, title, focus, closeOnEscape, holdToasts, backgroundColor, rows/columns, placement - keine Option, das × zu vergrößern oder auszublenden (claude-code.d.ts, PaneOpenArgs, Build 2.1.330). DoD 2 'das Kreuz wird zu einem Knopf dieser Größe' ist aus dem Mod heraus nicht machbar. (2) HQHFXQ (auf master, noch nicht released; Alex läuft 0.6.1) hat die Overview-Tabs schon zu gerahmten Knöpfen in Close-Höhe gemacht - drei Zeilen hoch, alle gleich. Der Screenshot zeigt 0.6.1, also den Stand davor. DoD 1 verlangt aber 'Close nicht mehr als dreizeiliger Rahmen' - widerspricht HQHFXQ.
- **2026-11-22 17:11 · Alexander Sacharov** — Frage an Alex (1/2): Welche Form für Tabs und Close im Kopf? 1) (Empfohlen) Eine Zeile im '[ Queue ]'-Stil für alle, Close als '[ Close ]' am rechten Rand - flach, braucht 1 statt 3 Zeilen, Klickfläche = ganzer Klammertext. 2) HQHFXQ-Stil behalten: alle gerahmt, drei Zeilen hoch - größte Klickfläche, kostet 2 Zeilen pro Pane; DoD 1-Halbsatz zu Close fällt weg. --- Frage (2/2): Das × gehört Claude Code, der Mod kann es nicht vergrößern. 1) (Empfohlen) × bleibt wie es ist; unser Close in jedem Pane mit Kopf ist der große Knopf daneben; DoD 2 wird so umformuliert. 2) Ticket für das × abgeben (Feedback an Anthropic), hier nur Tabs/Close.
- **2026-11-22 17:12 · Alexander Sacharov** — decision (Alex 22.11.): F1=1 alle Kopf-Knöpfe einzeilig im Stil '[ Queue ] [ Log ] [ Workers ]   [ Close ]', Close rechts, ganzer Klammertext klickbar (ersetzt die dreizeiligen Rahmen von HQHFXQ). F2=1 das × von Claude Code bleibt; die grosse Schliessen-Fläche ist unser Close in jedem Pane mit Kopf.
- **2026-11-22 17:19 · Alexander Sacharov** — Umsetzung: ein Helfer headRow (overview-pane.tsx) für jede Kopfzeile - Tabs links, Close rechts, Box nowrap + space-between, beide Teile flexShrink 0. Die Close-Knöpfe von Verdict/Explanation/Message standen bisher unten; sie stehen jetzt oben im Kopf, Give back und das Ticket-Pane (/cockpit ticket) bekommen erstmals ein Close. OVERVIEW_ROWS 3 -> 1, der Workers-Tab hat zwei Zeilen mehr (Test 40x100 jetzt <= 39). Offen für die Person: die Release-Note HQHFXQ.md (unreleased) sagt noch 'gerahmte Knöpfe'; R19PY9.md sagt das Gegenteil - beim Release lesen beide nacheinander, HQHFXQ.md könnte gestrichen werden.
- **2026-11-22 17:20 · Alexander Sacharov** — critique: headRow baut den Close-Knopf selbst (Schlüssel + Handler), role=dismiss nach ButtonProps; die sechs Aufrufer übergeben nur noch key und onPress.
- **2026-11-22 17:28 · Alexander Sacharov** — bughunt: pre-existing: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1944 - class: the 'nothing to show' fallback of a pane is drawn without a head, so it has no Close of its own and only the host's × closes it (small; absent on a phone). Same pattern at register.tsx:1707 (Give back), 1793 (Explanation), 2001 (ticket pane). The Overview's unread-board fallback (overview-pane.tsx:211) does have the head. Example: the Verdict pane is open, the ticket is accepted from another session and leaves the board's waiting list, the pane redraws to 'nothing' with no [ Close ]. blame: 3a1d27bc4e, 86da7c4354 and older (before the ticket); read
- **2026-11-22 17:28 · Alexander Sacharov** — bughunt: done at b0d5faff4 - 1 finding (0 change, 1 pre-existing) in 3 modules, first hunt; concurrency, error paths, input, permissions checked; ran: claude plugin test . (315 pass); files: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/R19PY9.md core/shipped/shipped.sha256
- **2026-11-22 17:35 · Alexander Sacharov** — given back by teamlead (22.11.), short round: (1) core/release/notes/HQHFXQ.md (unreleased) line about framed buttons rewrite to the final state, so the release does not carry two contradicting lines; (2) the 'nothing to show' fallback of Give back/Explanation/Verdict/Ticket pane (register.tsx ~1707/1793/1944/2001) gets the same head row with [ Close ] - Alex wants Close everywhere; test.
- **2026-11-22 17:38 · Alexander Sacharov** — Rückgaberunde: (1) HQHFXQ.md auf den Endstand umgeschrieben (eine Zeile, ohne Rahmen). (2) nothingPane (register.tsx) - der 'nothing recorded yet'-Fallback von Verdict, Give back, Explanation und Ticket-Pane trägt jetzt dieselbe Kopfzeile mit [ Close ]; Answer (nur ein Eingabefeld) und Update (hat Cancel) bleiben ohne Kopf, sie haben keinen. Test: 'a pane left with nothing to show keeps its head'. Give backs close() steht jetzt vor dem frühen Return, damit auch der Fallback die Atome leert.
- **2026-11-22 17:38 · Alexander Sacharov** — critique Runde 2: gelesen ist nur der Diff seit b0d5faff4. Liegen gelassen: Schlüssel und Handler von Close stehen in Verdict, Explanation und Ticket-Pane zweimal, einmal im Fallback und einmal im Kopf. Eine gemeinsame Konstante pro Pane spart je eine Zeile, macht aber einen weiteren Umlauf nicht wert.
- **2026-11-22 17:39 · Alexander Sacharov** — bughunt: done at 05d850d81 - 0 findings in 5 modules, follow-up hunt since b0d5faff4 (diff of register.tsx, board.test.ts, shipped.sha256; HQHFXQ.md read whole as new); concurrency, error paths, input, permissions checked; ran: claude plugin test . (316 pass), sha256sum register.tsx matches the added shipped.sha256 line; files: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/HQHFXQ.md core/release/notes/R19PY9.md core/shipped/shipped.sha256
