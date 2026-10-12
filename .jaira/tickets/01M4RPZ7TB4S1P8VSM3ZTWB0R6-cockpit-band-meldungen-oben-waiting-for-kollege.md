---
id: 01M4RPZ7TB4S1P8VSM3ZTWB0R6
title: "Cockpit-Band: Meldungen oben, 'waiting for <Kollege>' immer unten"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Das Band hat eine feste Reihenfolge: Meldungen oben, die eigene Arbeit in der Mitte, die Tickets der Kollegen unten"
context: |-
  Im Cockpit-Band steht die letzte Meldung unter dem Block der Kollegen. Dadurch springt die Liste.
  - Screenshot von Alex 2026-10-12: '[ 1 waiting for BeMuCa ]', darunter '776466 accepted: done', darunter die Zeile 'context 73k | cache ...'.
  - Wunsch von Alex: Tickets von Kollegen ('waiting for <name>', t.theirs in hooks/texts.ts) stehen immer ganz unten, die Meldungen (Verlauf, noticeAtom, z. B. '<id> accepted: done') oben.
  - Code: Render des Bands in core/mod/builtin/jaira-cockpit/hooks/register.tsx, ab ~Zeile 800. Reihenfolge heute: Kopfzeile -> 'Waiting for you' -> Kollegen (theirs) -> decisions -> notice -> Session-Zeile (context/cache/threads).
  - Soll: Kopfzeile -> notice -> 'Waiting for you' -> decisions -> Kollegen -> Session-Zeile.
  - Die Session-Zeile bleibt die letzte Zeile, sie ist der Status dieser Sitzung, keine Liste.
definition-of-done: "Das Band rendert in der Reihenfolge Kopfzeile, Meldung, 'Waiting for you', Entscheidungen, Kollegen, Session-Zeile; Test mit allen Teilen gleichzeitig"
tags: []
blocked-by: []
related: []
commits: []
created-at: 2026-10-12T18:46:01Z
updated-at: 2026-10-12T18:53:30Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-48411
claimed-at: 2026-10-12T18:49:44Z
outcome-what: "AbovePrompt-Render in core/mod/builtin/jaira-cockpit/hooks/register.tsx umgeordnet: Meldung (notice) direkt unter der Kopfzeile, dann 'Waiting for you' mit +N more, dann Entscheidungen, dann der Kollegen-Block ('waiting for <name>' und seine Zeilen), zuletzt die Session-Zeile. Zwei Tests in tests/board.test.ts, Release-Note core/release/notes/TWB0R6.md."
outcome-why: "Die Meldung stand unter dem Kollegen-Block; kam oder ging sie, sprang die Liste, und der Blick auf die eigene Arbeit fand sie an wechselnder Stelle."
outcome-resolves: "Band hat feste Reihenfolge: Meldungen oben, eigene Arbeit in der Mitte, Kollegen unten (DoD 1, 2, 3)."
---

# Cockpit-Band: Meldungen oben, 'waiting for <Kollege>' immer unten

## Definition of Done

- [x] Das Band rendert in der Reihenfolge Kopfzeile, Meldung, 'Waiting for you', Entscheidungen, Kollegen, Session-Zeile; Test mit allen Teilen gleichzeitig
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:809; test 'waiting for you > the band reads top down: head, notice, waiting for you, decisions, the teammates' tickets, the session line' (core/mod/builtin/jaira-cockpit/tests/board.test.ts:512), fails on the old order
- [x] Der Kollegen-Block steht auch aufgeklappt ('+N more') unter allen eigenen Tickets; Test
  proof: test 'waiting for you > opened by +N more, the person's own tickets still sit above the teammates'' (core/mod/builtin/jaira-cockpit/tests/board.test.ts:536)
- [x] Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/TWB0R6.md:1

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Render in register.tsx umordnen: notice direkt unter Kopfzeile, decisions vor den Kollegen-Block
- [x] Test: alle Teile gleichzeitig, Reihenfolge per findAll(Text)
- [x] Test: aufgeklappt (+N more und theirs offen) stehen Kollegen unter allen eigenen
- [x] Release-Note core/release/notes/TWB0R6.md

## Progress
- **2026-10-12 18:49 · Alexander Sacharov** — Teamlead 2026-10-12: Branch feat/twb0r6 (Worktree /home/alex/projects/.worktrees/jaira-twb0r6) ist von feat/qrhb35 (82aa5654) abgezweigt, nicht von master: der Cockpit-Code mit workers()/giveBack ohne --force existiert nur dort. feat/qrhb35 ist noch nicht gepusht/gemergt; der PR dieses Tickets stapelt sich darauf. Parallel laufen WP1CHK (giveBack/Version), GPFA7B (accept), TWB0R6 (Band-Render) - alle in hooks/register.tsx; nur den eigenen Bereich anfassen.
- **2026-10-12 18:49 · Alexander Sacharov** — parent session: jaira-f5
- **2026-10-12 18:50 · Alexander Sacharov** — Offene Entscheidungen: keine. Reihenfolge steht im Kontext wörtlich; Session-Zeile bleibt letzte Zeile.
- **2026-10-12 18:53 · Alexander Sacharov** — Testweg: der Wurzel-Box-Text aus ui.find({type:'Box'}) ist das ganze Band in Zeichenreihenfolge; indexOf darauf prüft die Reihenfolge über Text und Button hinweg (findAll nur über Text sieht die Kollegen-Zeile nicht, sie ist ein Button). Der +N-more-Test war schon mit dem alten Code grün: theirs stand dort schon unter dem more-Knopf - er schützt nur vor Rückfall. Der Reihenfolge-Test schlägt mit dem alten Render fehl (geprüft durch Rücksetzen von register.tsx). Lauf: claude plugin test in core/mod/builtin/jaira-cockpit, 83 pass.
