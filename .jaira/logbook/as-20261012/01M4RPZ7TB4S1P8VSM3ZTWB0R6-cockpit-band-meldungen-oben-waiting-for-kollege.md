---
id: 01M4RPZ7TB4S1P8VSM3ZTWB0R6
title: "Cockpit-Band: Meldungen oben, 'waiting for <Kollege>' immer unten"
status: done
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
commits:
  - 8a074857b0e696bdbe50056a9ea115d5e7797b40
  - 5af87faaf4ad0059502ecc9315642350dcc0dc3b
created-at: 2026-10-12T18:46:01Z
updated-at: 2026-10-12T19:09:15Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-48411
claimed-at: 2026-10-12T18:49:44Z
outcome-what: "AbovePrompt-Render in core/mod/builtin/jaira-cockpit/hooks/register.tsx umgeordnet: Meldung (notice) direkt unter der Kopfzeile, dann 'Waiting for you' mit +N more, dann Entscheidungen, dann der Kollegen-Block ('waiting for <name>' und seine Zeilen), zuletzt die Session-Zeile. Zwei Tests in tests/board.test.ts, Release-Note core/release/notes/TWB0R6.md."
outcome-why: "Die Meldung stand unter dem Kollegen-Block; kam oder ging sie, sprang die Liste, und der Blick auf die eigene Arbeit fand sie an wechselnder Stelle."
outcome-resolves: "Band hat feste Reihenfolge: Meldungen oben, eigene Arbeit in der Mitte, Kollegen unten (DoD 1, 2, 3)."
review-summary: "In register.tsx (AbovePrompt-Render des Cockpit-Bands) sind zwei Blöcke verschoben, keine Logik neu: die Meldung (notice, z. B. 'MG677R accepted: done') steht jetzt direkt unter der Kopfzeile statt ganz unten, und die Entscheidungen ('<id> | decided in Discuss') stehen vor dem Kollegen-Block statt danach. Damit lautet die Reihenfolge Kopfzeile -> Meldung -> Waiting for you (+N more) -> Entscheidungen -> 'N waiting for <name>' mit seinen Zeilen -> Session-Zeile. Kommt oder geht die Meldung, verschiebt sich nur alles darunter gleichmäßig, die Kollegen bleiben immer unten. Zwei neue Tests in tests/board.test.ts lesen das ganze Band als einen String (Helfer drawn() neben headLine/sessionLine) und prüfen die Reihenfolge per indexOf: einmal mit allen Teilen gleichzeitig, einmal aufgeklappt (+N more und Kollegen offen). Release-Note core/release/notes/TWB0R6.md, eine Zeile."
review-gaps: "Keine Fehler. Drei Randpunkte, keiner blockiert: (1) Entscheidungen haben weder Überschrift noch Abstand; ohne eigene wartende Tickets kleben sie direkt unter der Meldung - war vorher unter dem Kollegen-Block genauso, nur kosmetisch. (2) Der +N-more-Test war laut Notiz schon mit dem alten Code grün, er schützt nur vor Rückfall; DoD 2 ist damit belegt, aber nicht durch eine Änderung erzeugt. (3) Der Branch feat/twb0r6 ist von feat/qrhb35 (82aa5654) abgezweigt, das noch nicht gemergt ist: der PR stapelt sich darauf und kann erst nach QRHB35 sauber nach master. Geprüft: Diff 82aa5654..HEAD deckt sich mit outcome-what; claude plugin test in core/mod/builtin/jaira-cockpit: 83 pass, 0 fail; Critique-Befund (drawn() top-level) in 5af87faa umgesetzt."
test-verdict: "pass: claude plugin test RC=0 (83 pass), claude plugin validate RC=0, go build ./... RC=0, go test ./core/mod/... RC=0; DoD 1-3 verified in the tree (register.tsx:809, board.test.ts:514 and :538 pass, TWB0R6.md one '- ' line); the order test fails on the old render; the live band in Claude Code was not looked at"
review-verdict: "Der Diff erfüllt alle drei DoD-Punkte, die kleinste mögliche Form (nur Blöcke verschoben). Keine Defekte gefunden. Empfehlung: annehmen; beim PR beachten, dass er auf feat/qrhb35 aufsetzt."
review-check: "1. cd /home/alex/projects/.worktrees/jaira-twb0r6/core/mod/builtin/jaira-cockpit && claude plugin test  -> '83 pass, 0 fail'. 2. Im Claude-Code-Cockpit dieses Worktrees ein eigenes Ticket in signoff mit Accept annehmen. 3. Erwartet: die Zeile '<id> accepted: done' steht direkt unter der Kopfzeile 'jaira-cockpit', ÜBER 'Waiting for you'. 4. Hat ein Kollege wartende Tickets: die Zeile '[ N waiting for <name> ]' steht als letzte Liste, direkt über der Zeile 'context ...k | cache ...'. 5. Diese Kollegen-Zeile anklicken und, falls mehr als 3 eigene Tickets, '+N more' -> alle eigenen Tickets und 'show fewer' stehen weiter über 'N waiting for <name>'."
---

# Cockpit-Band: Meldungen oben, 'waiting for <Kollege>' immer unten

## Definition of Done

- [x] Das Band rendert in der Reihenfolge Kopfzeile, Meldung, 'Waiting for you', Entscheidungen, Kollegen, Session-Zeile; Test mit allen Teilen gleichzeitig
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:809; test 'waiting for you > the band reads top down: head, notice, waiting for you, decisions, the teammates' tickets, the session line' (core/mod/builtin/jaira-cockpit/tests/board.test.ts:514), fails on the old order
- [x] Der Kollegen-Block steht auch aufgeklappt ('+N more') unter allen eigenen Tickets; Test
  proof: test 'waiting for you > opened by +N more, the person's own tickets still sit above the teammates'' (core/mod/builtin/jaira-cockpit/tests/board.test.ts:538)
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
- **2026-10-12 18:53 · Alexander Sacharov** — critique: drawn() in tests/board.test.ts gehört als top-level async function neben headLine/sessionLine (:190-201), nicht als const im describe. Sonst nichts: Render-Umbau ist die kleinste Form (nur Blöcke verschoben, keine neue Logik); Plan-Punkt 2 sagt findAll(Text), gebaut ist Wurzel-Box-Text - die Testweg-Notiz erklärt das, kein Befund.
