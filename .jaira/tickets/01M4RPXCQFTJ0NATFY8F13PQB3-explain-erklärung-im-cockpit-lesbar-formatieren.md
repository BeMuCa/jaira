---
id: 01M4RPXCQFTJ0NATFY8F13PQB3
title: Explain-Erklärung im Cockpit lesbar formatieren und Optionen direkt beantwortbar machen
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die Erklärung eines Review-Befunds ist auf einen Blick lesbar, und die Person kann eine der Optionen (a/b/c) mit einem Druck wählen"
context: "Die Explain-Erklärung im Cockpit ist ein einziger Textblock und schwer zu lesen; antworten kann man darauf nicht.\n- Screenshot von Alex 2026-10-12: Pane 'Explanation of 776466' zeigt alles als einen Absatz, Optionen '(a) ... (b) ... (c) ...' mitten im Fliesstext.\n- Woher der Text kommt: spawn.sh (core/role/builtin/jaira-dispatcher/scripts/spawn.sh ~Zeile 171) lässt den Explain-Chat 'jaira note <id> \"explain: <Erklärung>\"' schreiben, 'in wenigen Sätzen' - kein Format vorgegeben.\n- Wie er gezeigt wird: Knopf 'Read explanation' -> show() in core/mod/builtin/jaira-cockpit/hooks/register.tsx (~Zeile 306), reiner Text in einem Pane.\n- Zwei Teile:\n  1) Format: Explain-Prompt verlangt Struktur (Befund / warum es zählt / Optionen je eine Zeile / Empfehlung), und das Pane rendert Zeilenumbrüche bzw. Listen.\n  2) Antworten: Optionen als Knöpfe im Pane ('a', 'b', 'c'); die Wahl wird als 'decision:'-Note aufs Ticket geschrieben, wie es der Discuss-Chat heute tut (noteOf(..., 'decision') in hooks/board.ts).\n- Ungeklärt, für brainstorm: was nach der Wahl passiert - nur notieren, oder gleich Accept / Give back mit der Wahl als Grund auslösen."
definition-of-done: "Die Explain-Note hat eine feste Struktur (Befund, Begründung, Optionen je Zeile, Empfehlung), vorgegeben im Explain-Prompt von spawn.sh"
tags: []
blocked-by: []
related: []
commits: []
created-at: 2026-10-12T18:45:00Z
updated-at: 2026-10-12T19:02:26Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-65657
claimed-at: 2026-10-12T18:56:47Z
outcome-what: "Explain-Pane (cockpit-explanation) zeigt die Note zeilenweise mit einer Zeile je Option und einem Select (a/b/c + Own answer…, Empfehlung vorgewählt); die Wahl oder eine eigene Antwort (eigenes Pane mit focus) wird decision:-Note; Verdict-Knopf zeigt das ganze Verdikt; Explain erscheint auch, wenn review.gaps eine Lücke nennt; spawn.sh verlangt die Zeilenform."
outcome-why: "Die Erklärung war ein Textblock ohne Antwortmöglichkeit, das Verdikt im Band abgeschnitten und nirgends ganz lesbar, und Explain fehlte bei Verdikten ohne Vorbehaltswort (WP1CHK)."
outcome-resolves: "DoD 1-9: spawn.sh:171, board.ts explanationOf/namesGap, register.tsx Panes EXPLANATION/ANSWER und Verdict-Knopf, sechs Tests in tests/board.test.ts, core/release/notes/13PQB3.md"
---

# Explain-Erklärung im Cockpit lesbar formatieren und Optionen direkt beantwortbar machen

## Definition of Done

- [x] Die Explain-Note hat eine feste Struktur (Befund, Begründung, Optionen je Zeile, Empfehlung), vorgegeben im Explain-Prompt von spawn.sh
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:171 (explain prompt: Befund / Warum: / a) b) ... / Empfehlung:, one line each)
- [x] Das Pane 'Explanation of <id>' zeigt diese Struktur mit Absätzen und einer Zeile je Option; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts:223 explanationOf, core/mod/builtin/jaira-cockpit/hooks/register.tsx:964 pane; test 'a shaped explanation reads as paragraphs and one line per option, the recommended one chosen' (core/mod/builtin/jaira-cockpit/tests/board.test.ts:421)
- [x] Jede Option ist im Pane wählbar; die Wahl landet als 'decision:'-Note auf dem Ticket und das Band zeigt sie; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:413 decide; test 'an option picked goes onto the ticket as a decision, and the band offers it to read' (core/mod/builtin/jaira-cockpit/tests/board.test.ts:442)
- [x] Eine alte Explain-Note ohne Struktur wird weiter angezeigt wie bisher; Test
  proof: test 'what the Explain chat wrote onto the ticket is read in the cockpit' (core/mod/builtin/jaira-cockpit/tests/board.test.ts:386): unshaped note shown as one Text, whole
- [x] Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/13PQB3.md
- [x] Neben den Optionen gibt es immer 'Eigene Antwort': ein Textfeld im Pane, der Text landet als 'decision:'-Note wie eine gewählte Option; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1005 answer pane; test 'an own answer is typed in a pane that takes the keyboard and goes onto the ticket as a decision' (core/mod/builtin/jaira-cockpit/tests/board.test.ts:458)
- [x] Das Pane ist interaktiv: Optionen und 'Eigene Antwort' per Pfeiltasten/Enter wählbar, Fokus liegt beim Öffnen auf der Liste (wie das Give-back-Pane mit focus: true); Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:984 Select autoFocus in a pane opened with focus: true; tests at core/mod/builtin/jaira-cockpit/tests/board.test.ts:421 (autoFocus, options a,b,c,own) and core/mod/builtin/jaira-cockpit/tests/board.test.ts:458
- [x] Der volle Text des Review-Verdikts ist aus dem Band lesbar: ein Druck auf die (abgeschnittene) Verdikt-Zeile oder ein Knopf öffnet ihn ganz im Pane, auch ohne Explain-Note; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:868 Verdict button; test 'the whole review verdict is read from the band, with no explanation on the ticket' (core/mod/builtin/jaira-cockpit/tests/board.test.ts:493)
- [x] Explain erscheint, sobald das Review eine Lücke nennt (review.gaps nicht leer bzw. nicht 'none'), nicht nur wenn das Verdikt ein Wort aus RESERVATION trifft; Test mit dem WP1CHK-Verdikt
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts:160 namesGap, core/mod/builtin/jaira-cockpit/hooks/register.tsx:861; test 'Explain is offered once the review names a gap, whatever words the verdict uses' (core/mod/builtin/jaira-cockpit/tests/board.test.ts:477) with WP1CHK's verdict

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-12 18:47 · Alexander Sacharov** — Alex 2026-10-12: Hauptsache, es gibt immer die Möglichkeit, eine eigene Antwort zu schreiben, und die Auswahl ist interaktiv (wie AskUserQuestion mit 'Other'). Vorbild im Code: askWhatIsMissing/GIVE_BACK-Pane in hooks/register.tsx öffnet ein Pane mit focus: true und autoFocus-Feld.
- **2026-10-12 18:56 · Alexander Sacharov** — Teamlead 2026-10-12: Branch feat/13pqb3 ist von feat/twb0r6 abgezweigt (das von feat/qrhb35): TWB0R6 hat die Reihenfolge im Band-Render von hooks/register.tsx neu gebaut, dieses Ticket setzt darauf auf. PR stapelt sich auf TWB0R6 und QRHB35 (PR #63). Parallel laufen GPFA7B (accept) und WP1CHK (giveBack/Version) in derselben Datei - deren Funktionen nicht anfassen.
- **2026-10-12 18:56 · Alexander Sacharov** — parent session: jaira-f5
- **2026-10-12 18:58 · Alexander Sacharov** — Alex 2026-10-12 (Screenshot): im Band ist das Verdikt hinter '…' abgeschnitten ('The diff meets all five DoD items … the one weak sp…') und nirgends ganz lesbar; Explain gibt es nur bei einem Vorbehalt. Daher DoD 8: voller Verdikt-Text per Druck im Pane.
- **2026-10-12 18:58 · Alexander Sacharov** — Offene Entscheidungen vor dem Plan: keine. 'Was nach der Wahl passiert' entscheidet DoD 3: nur die decision:-Note und das Band zeigt sie - kein Accept/Give back ausgelöst; das bleiben die Knöpfe im Band. Form (aus Code gesetzt): eigenes Pane 'cockpit-explanation' statt des Meldungs-Panes; Struktur = Zeilen in der Note, Optionen als Zeilen 'a) ...', 'Empfehlung: <Buchstabe>' wählt die Option im Select vor. Auswahl per Select (Engine-Element, Pfeile+Enter, autoFocus). 'Eigene Antwort' ist die letzte Option im Select; gewählt öffnet sie ein eigenes Pane mit focus:true und autoFocus-Feld, genau wie Give back - ein Input, das erst nach einem Redraw erscheint, bekommt den Fokus nicht sicher. Bandtext 'decided in Discuss' wird 'decided', weil die Entscheidung jetzt auch aus dem Explain-Pane kommt.
- **2026-10-12 19:00 · Alexander Sacharov** — Alex 2026-10-12 (Screenshot): WP1CHK in signoff hatte keinen Explain-Knopf, obwohl das Verdikt einen Vorbehalt hat: 'No defect that blocks: the one weak spot is that a failed note after a successful move is not reported (gap 1) ... Accept, with gap 1 either fixed now or knowingly left.' Ursache: hasReservation() in hooks/board.ts (~Zeile 136) sucht nur Wörter (reservation, but, aber, caveat ...); 'weak spot', 'gap', 'either ... or' trifft keins. Eine Wortliste rät; review.gaps sagt es direkt. DoD 9.
- **2026-10-12 19:01 · Alexander Sacharov** — Gefunden: Select ist ein Engine-Element (Pfeile bewegen, Enter wählt, autoFocus wie Input); Tests wählen mit pane.select({key, value}). Typecheck: kein tsc lokal installiert, nur 'claude plugin validate' (passed) und 'claude plugin test' (88 pass). DoD 9: review.gaps kommt bei Tickets in 'human' (vor review) noch vom optimize-Lane (dort steht, was entfernt wurde) - auch dann erscheint Explain; das ist die wörtliche Forderung, aber ein möglicher Fehlalarm, den der Reviewer sehen soll.
