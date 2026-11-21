---
id: 01M7ZR6VCT7AQD88YMWARNZYJR
title: "Cockpit-Workers: ein Steinbruch statt einzelner Bahnen, und nur Tickets mit eigenem Agenten"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Der Workers-Tab zeigt auf einen Blick, welche Agenten von Alex gerade woran arbeiten, als gemeinsamen Steinbruch an einem großen Stein."
context: |-
  Was heute verwirrt (core/mod/builtin/jaira-cockpit/hooks/overview.ts workersOf, overview-pane.tsx trackRows):
  1. Jede Figur läuft auf einer eigenen Bahn von Lane zu Lane; man sieht nicht, welche Rolle (worker/review/bughunt) sie ist, wie lange sie dran ist und was sie gerade tut.
  2. Tickets mit offener Frage erscheinen als Figur, auch ohne Agent und auch fremde: am 21.11. standen dort RNQRJB, SZFTQD, T22Z2C (human) und ZE9WZC (signoff) — alle von BeMuCa, seit Oktober liegengeblieben. Ursache: workersOf schiebt eine Figur, wenn die Frage offen ist ("if (walkers.length === before && asks)").
  Entschieden von Alex am 21.11.:
  - Bild: ein Steinbruch. Alle Agenten hauen gemeinsam an einem großen Stein (die offene Arbeit); Stein wird kleiner, wenn DoD-Punkte abgehakt werden, ein Brocken fällt, wenn ein Ticket in signoff/done kommt.
  - Unter dem Bild pro Agent eine Zeile: ID, Lane, Rolle (aus dem Tab-Label <ID>/<rolle>), arbeitet/wartet, seit wann, erste Zeile des letzten Berichts (/cockpit report hat sie schon).
  - Figur nur für einen echten Agenten (Herdr-Tab oder Subagent). Ein Ticket mit offener Frage ohne Agent ist keine Figur.
  Offen für pre-process: was "der Stein" genau misst (offene DoD-Punkte aller Tickets mit Agent ist der Vorschlag), und wo Fragen ohne Agent stattdessen stehen (Queue-Tab hat schon Gruppen).
  Gleiche Datei wie 02D92C (overview-pane.tsx) — erst nach 02D92C starten.
definition-of-done: Ein Ticket mit offener Frage und ohne Agent erscheint nicht im Workers-Tab (Test in tests/)
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-11-21T19:09:25Z
updated-at: 2026-11-21T20:08:07Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-243922
claimed-at: 2026-11-21T19:47:52Z
outcome-what: "Workers-Tab neu: workersOf liefert nur echte Agenten (Herdr-Tab <ID>/<rolle> oder Subagent) mit Rolle und offenen DoD-Punkten; eine Szene mit einem Stein (stoneSize, wächst mit der Wurzel der offenen Punkte), Krabben links/rechts am Boden, kleine Krabben wenn zu viele; ein Brocken fällt, wenn die offenen Punkte weniger werden; darunter eine Zeile pro Agent (ID, Lane, Rolle, Zustand, seit wann, erste Zeile des letzten Berichts). Walk-/Bahn-Logik (progressOf, stepOf, Spot) entfernt."
outcome-why: "Bahnen pro Figur brauchten ab 5 Agenten Scrollen, zeigten weder Rolle noch Dauer noch Tätigkeit, und Tickets mit offener Frage ohne Agent (auch fremde, seit Oktober liegende) standen als Figur da."
outcome-resolves: "DoD 1-5: nur Agenten (Test), ein Stein mit Größenberechnung (Test), eine Zeile pro Agent (Test), live geprüft + Release-Note, 1-30 Agenten passen in 40×100 (parametrisierter Test)."
review-summary: none
round: 3
back-from: testing
back-by: lane
back-count: 1
review-gaps: "sceneRows: placesOf for the big crab is computed once and reused instead of twice (fits lambda gone). Searched for an existing duration formatter (cache.ts, sessions.ts only print whole minutes inline) — none to reuse, durationOf stays. Nothing the change orphaned: progressOf/stepOf/Spot/WALK_STEP, the tests' plan helper and Ticket.next_lane/plan_items went with it. Left the per-read store lookup of reports (one key per ticket every 5 s) as is."
test-verdict: "fail: go test ./core/role RC=1 — TestEveryInstalledFileIsAShippedVersion: core/shipped/shipped.sha256 lacks the versions of 6 changed cockpit files; cockpit suite itself green (237 pass, RC=0), DoD 1-5 verified in the tree and live"
hunt-verdict: "fail: bughunt done at 330eb71c4, 1 change (reportLine strips any <…> span from report content), 1 pre-existing (subagent naming two handles counted at both — same subject, note for the person)"
---

# Cockpit-Workers: ein Steinbruch statt einzelner Bahnen, und nur Tickets mit eigenem Agenten

## Definition of Done

- [x] Ein Ticket mit offener Frage und ohne Agent erscheint nicht im Workers-Tab (Test in tests/)
  proof: core/mod/builtin/jaira-cockpit/hooks/overview.ts:149 workersOf; test 'the workers tab > a ticket with an open question and no agent is no figure…' and '…nobody on the board'
- [x] Alle Agenten stehen in einer Szene an einem Stein; seine Größe folgt den offenen DoD-Punkten (Test für die Größenberechnung)
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:94 sceneRows, core/mod/builtin/jaira-cockpit/hooks/overview.ts:195 stoneSize; tests 'the stone: a pebble for nothing open…', 'every agent stands at one stone…', 'a definition-of-done item ticked makes the stone smaller…'
- [x] Pro Agent eine Zeile mit ID, Lane, Rolle, Zustand, Dauer und erster Zeile des letzten Berichts
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:153 workerLines; test 'one line per agent: ticket, lane, kind, state, how long…'
- [x] Live geprüft mit claude --plugin-dir; Release-Note core/release/notes/<ID>.md
  proof: live 2026-11-21: tmux 210×46, claude --plugin-dir core/mod/builtin/jaira-cockpit, Overview → Workers: 5 Agenten an einem Stein ('5 open'), 5 Zeilen darunter; core/release/notes/RNZYJR.md
- [x] Alle Agenten und ihre Zeilen passen ohne Scrollen auf einen Bildschirm (Pane etwa 40 Zeilen × 100 Spalten, wie am 21.11.); bei mehr Agenten wird die Szene kleiner, nicht länger
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:94 sceneRows (placesOf: big crab, else small); tests '1|4|8|16|30 agents and all their lines fit 40 rows × 100 columns without scrolling'

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] overview.ts: workersOf nur Agenten (Rolle, offene DoD-Punkte), stoneOf, reportLine, Walk-Logik raus
- [x] overview-pane.tsx: Szene (Stein + Krabben links/rechts, klein bei vielen) + eine Zeile pro Agent
- [x] register.tsx: seit-wann-Atom, letzter Bericht je Ticket, Brocken-Frame; types/texts
- [x] Tests im Block 'the workers tab' neu; claude plugin test grün
- [x] Release-Note, Live-Check mit --plugin-dir

## Progress
- **2026-11-21 19:10 · Alexander Sacharov** — decision: Alex 21.11.: heute läuft jede Figur auf einer eigenen Bahn von ~6 Zeilen, ab 5 Agenten muss man scrollen. Ziel: eine Szene, ein Bildschirm, kein Scrollen.
- **2026-11-21 19:48 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock
- **2026-11-21 19:49 · Alexander Sacharov** — decision (worker, offene Punkte aus dem Kontext selbst entschieden): (1) Der Stein misst die offenen DoD-Punkte (superseded zählt nicht) aller Tickets, an denen gerade ein Agent ist — der Vorschlag aus dem Kontext. Ohne offene Punkte bleibt ein Kieselstein, damit die Szene nie leer ist. (2) Fragen ohne Agent bekommen keinen neuen Ort: das Band zeigt Tickets in human/signoff schon als 'waiting for you'; der Workers-Tab zeigt nur Agenten. (3) 'seit wann' = seit das Cockpit den Agenten zum ersten Mal sieht (pro <ID>/<rolle> gemerkt, solange er lebt): Herdr-Tabs und Subagenten nennen keine Startzeit, die Session-Records auch nicht. (4) Rolle = Teil nach '/' im Tab-Label; Subagent = 'subagent'. Zustand: working / resting / waits (Ticket hat offene Frage). (5) Brocken: verschwindet ein Ticket mit offenen Punkten aus der Szene, fällt einige Frames lang ein Brocken vom Stein.
- **2026-11-21 19:59 · Alexander Sacharov** — gefunden: der Pane-Test zählt Zeilen über findAll — eine Bildzeile aus nur einem Stil erscheint dort doppelt (äußerer Text + einzelner Run); pictureOf in tests/board.test.ts entdoppelt deshalb. Live-Check geht ohne Person: tmux new-session ... 'claude --plugin-dir <mod>', Klick per SGR-Maus-Sequenz (tmux send-keys -l ESC[<0;x;yM/m); unter ~200 Spalten sagt das Cockpit 'No room for a pane'. 6-Minuten-Clock-Advance mit laufenden Frames läuft in den 5-s-Testtimeout — Test wechselt den Tab, damit keine Frames laufen. Vier Kategorien: Nebenläufigkeit (readWorkers hat 'reading'-Sperre, Berichte unter dem reports-Lock), Fehlerpfade (Store-Fehler → Zeilen ohne Bericht statt leerer Tab), ungewöhnliche Eingaben (leerer/reiner Tag-Bericht → Titel; 0/negativ offene Punkte → Kiesel; zu viele Agenten → '+N more' im Bild, Zeilen bleiben), Rechte: nichts Neues von außen.
- **2026-11-21 20:00 · Alexander Sacharov** — critique (Runde 1): (1) register.tsx workersStone {open, fellAt} → nur fellAt; vorheriges open aus workersAtom (stoneOf) vor dem Überschreiben. (2) reportLine gehört nach hooks/digest.ts zu den anderen Report-Parsern. (3) overview-pane.tsx: export von WorkersScene und FALL_FRAMES weg, keiner außerhalb nutzt sie.
- **2026-11-21 20:01 · Alexander Sacharov** — critique (Runde 2): alle drei Befunde aus Runde 1 erledigt (52bdc0f09) — workersFell statt workersStone, reportLine in digest.ts, Exporte weg. Nichts Neues.
- **2026-11-21 20:04 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/digest.ts:118 (reportLine) - class: report text containing a literal '<' before a later '>' on the same line (generics, comparisons, placeholders like <id>) has everything between them cut out, because /<[^>]*>/g strips any angle-bracket span, not only the notification's own tags; example: 'keyed by Map<string, number> and a < b check > fixed' shows as 'keyed by Map and a fixed', and 'if n<3 then x>2' as 'if n2' on the Workers tab line; the stored report is intact, only the shown line is wrong; not a blocker; ran: node --experimental-strip-types on reportLine with both strings
- **2026-11-21 20:04 · Alexander Sacharov** — bughunt: pre-existing: core/mod/builtin/jaira-cockpit/hooks/overview.ts:171 (workersOf) - class: a subagent whose description or name mentions more than one ticket handle is counted as an agent at every one of them (plain substring includes); example: description 'worker AAA111 (follows BBB222)' puts a working subagent at AAA111 and a phantom one at BBB222, and since this change the stone also adds BBB222's open items (stone 5 instead of 2); blame: the same includes() test is in fdccb30e3^ overview.ts:174, only the stone count is new; ran: node --experimental-strip-types workersOf/stoneOf with that agent
- **2026-11-21 20:04 · Alexander Sacharov** — bughunt: done at 330eb71c4 - 2 findings (1 change, 1 pre-existing) in 7 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/digest.ts core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx core/mod/builtin/jaira-cockpit/hooks/overview.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/types/index.d.ts core/release/notes/RNZYJR.md
- **2026-11-21 20:05 · Alexander Sacharov** — testing: class: a generated manifest not regenerated with the change; example: go test ./core/role → TestEveryInstalledFileIsAShippedVersion lists digest.ts, overview-pane.tsx, overview.ts, register.tsx, texts.ts, types/index.d.ts; fix: scripts/shipped-history.sh --write after the last code edit, and run go test ./core/role before every push. Plus bughunt change: reportLine — strip only lines that are nothing but tags, never <…> inside content; test with 'Map<string, number>' and 'n<3 then x>2'.
- **2026-11-21 20:08 · Alexander Sacharov** — Runde 2 in-progress: reportLine überspringt jetzt nur Zeilen, die ganz aus Tags/ganzen Elementen bestehen (TAGS_ONLY in digest.ts), Inhalt behält seine spitzen Klammern; shipped.sha256 per scripts/shipped-history.sh --write neu; go test ./... RC=0. Der pre-existing Befund (Subagent mit zwei Handles) bleibt als Notiz für die Person: gleiche Stelle wie vor diesem Ticket, ein Fix bräuchte eine Regel, welche Handle in einer Beschreibung 'die eigene' ist.
