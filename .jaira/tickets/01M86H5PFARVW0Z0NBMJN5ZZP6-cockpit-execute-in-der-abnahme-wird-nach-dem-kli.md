---
id: 01M86H5PFARVW0Z0NBMJN5ZZP6
title: "Cockpit: Execute in der Abnahme wird nach dem Klick grau und nimmt keinen zweiten Klick, bis der Schritt fertig ist"
status: critique
ready: true
creator: Sascha
assignee: Sascha
goal: "Nach einem Klick auf Execute (bzw. Accept) im Abnahme-Pane sieht die Person sofort, dass der Klick angekommen ist: der Knopf wird grau mit 'laeuft…', nimmt keinen weiteren Klick, und wird wieder normal, wenn der Schritt fertig ist (oder spaetestens nach einer Zeitgrenze)."
context: |-
  Was falsch ist: Im Cockpit-Abnahme-Pane (core/mod/builtin/jaira-cockpit/hooks/register.tsx, Knopf Execute, Texte in hooks/texts.ts:63-70) sieht der Knopf nach dem Klick genauso aus wie davor. Die Person weiss nicht, ob der Klick ankam, und klickt nochmal.
  Ein zweiter Klick wird heute erst abgewiesen, wenn er schon laeuft: texts.ts:64 busyRow ('an earlier press ... is still at work on it') - als Textzeile, nicht am Knopf.
  Ausloeser: Sascha (Alex) 24.11.2026: 'die Knopf zum Annehmen unten ... besser wenn ich klicke, dass sie nicht mehr klickbar ist, damit klar ist, dass ich gedrueckt habe - z.B. 5 Sekunden grau'.
  Gilt fuer Execute in allen Modi (accept, give back) und den Accept-Knopf des Panes.
definition-of-done: "Nach dem Klick ist Execute sofort grau/dim mit einem Text wie 'laeuft…' und reagiert nicht auf weitere Klicks oder Enter; Test"
tags:
  - cockpit
blocked-by: []
related: []
commits: []
created-at: 2026-11-24T10:21:08Z
updated-at: 2026-11-24T10:36:19Z
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-2605896
claimed-at: 2026-11-24T10:21:37Z
outcome-what: "Execute in the Acceptance pane reads 'running…' and is drawn dim while its row's lock is held; a press on it or Enter in the field does nothing then (the busyRow toast is gone). When BUSY_MS (20 s) frees a row whose step is still at work, once() says so in one toast; a held request stays silent. The lock is keyed per press, not per millisecond."
outcome-why: "A press gave no visible answer, so the person pressed again (Alex 24.11.2026)."
outcome-resolves: "DoD 1-3; checked: concurrency (two presses same ms - fixed key collision), error paths (refused step frees at once, tested), unusual input (Enter with text while locked), permissions (n/a, no outside input)"
---

# Cockpit: Execute in der Abnahme wird nach dem Klick grau und nimmt keinen zweiten Klick, bis der Schritt fertig ist

## Definition of Done

- [x] Nach dem Klick ist Execute sofort grau/dim mit einem Text wie 'laeuft…' und reagiert nicht auf weitere Klicks oder Enter; Test
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'Execute reads running… and takes no press and no Enter until its step is done'; hooks/register.tsx Execute button label/dim from busy
- [x] Ist der Schritt fertig (Erfolg oder Fehler), wird der Knopf wieder normal; endet nichts, spaetestens nach einer Zeitgrenze (z.B. 30 s) mit einem Hinweis; Test
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'a step refused reads Execute again at once, the refusal beside it', 'a step that has not ended after 20 seconds frees Execute and says so once'; hooks/register.tsx once() stillRunning
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/N5ZZP6.md
- [ ] seen by a person: in the Acceptance pane, Execute pressed on a ticket whose move takes a moment reads 'running…' dim until the pane closes or the refusal shows

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-24 10:21 · Sascha** — parent session: uds:/run/user/1000/cc-socks/3417.sock
- **2026-11-24 10:22 · Sascha** — open decisions: none. Settled: (1) the existing row lock (once/busyAtom, BUSY_MS=20s, register.tsx:1063) already blocks a second press and draws the button dim - the gap is only the label and the timeout hint; (2) busy Execute reads 'running…' and its press does nothing (was: busyRow toast), Input submitLabel follows; (3) time limit stays BUSY_MS=20s (DoD says 'e.g. 30 s'; 20 s is already tested and shipped); (4) when the limit frees a row whose action is still running, once() says so in one line - not for 'hold' actions, whose lock is the limit by design.
- **2026-11-24 10:35 · Sascha** — found out: (1) reading an atom ($.state.get) inside a $.clock.after callback is refused by the host ('no hooks module ... no scan to allow it') - update() works there, so once() learns whether its lock was still there from inside its own update; (2) the lock was keyed by clock.now(), so two presses in the same millisecond (band Acceptance, then Execute) freed each other's lock at BUSY_MS - real life almost never, the mock clock always; now a per-press counter is added. (3) busyRow text removed: a locked Execute now says 'running…' on itself; the old 'said, not dropped' test was changed to check the label instead.
