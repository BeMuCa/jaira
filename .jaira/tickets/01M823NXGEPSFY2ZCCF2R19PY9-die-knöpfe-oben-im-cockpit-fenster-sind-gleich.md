---
id: 01M823NXGEPSFY2ZCCF2R19PY9
title: "Die Knöpfe oben im Cockpit-Fenster sind gleich groß und gut klickbar: Tabs, Close und das Kreuz"
status: critique
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
commits: []
created-at: 2026-11-22T17:08:22Z
updated-at: 2026-11-22T17:22:18Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2402047
claimed-at: 2026-11-22T17:08:58Z
outcome-what: "Every cockpit pane head is one line: tabs as [ label ] on the left, [ Close ] at the right (headRow); framed 3-line buttons removed; Give back and the ticket pane got a Close"
outcome-why: "Alex's screenshot: tabs, a 3-line framed Close and a tiny × looked three different ways and were hard to hit"
outcome-resolves: "R19PY9: one shape, one height, Close as the large close target in every headed pane"
review-summary: "core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx headRow: all six callers build the same Close Button by hand (key, t.close, onPress); let headRow build it from a close key and handler, with role=\"dismiss\" (ButtonProps: the host's documented mark for a site's close button, native close on desktop)"
round: 2
back-from: critique
back-by: lane
back-count: 1
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
- [ ] seen by a person: Alex klickt Tabs, Close und Kreuz im laufenden Cockpit

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
