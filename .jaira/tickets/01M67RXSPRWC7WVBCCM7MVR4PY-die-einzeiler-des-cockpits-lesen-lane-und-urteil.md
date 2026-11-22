---
id: 01M67RXSPRWC7WVBCCM7MVR4PY
title: Die Einzeiler des Cockpits lesen Lane und Urteil aus einer festen Schlusszeile statt aus dem ersten passenden Wort
status: critique
ready: true
creator: berkc
assignee: Alexander Sacharov
goal: "Ein Worker-Bericht endet mit einer festen Zeile (z.B. 'jaira: <ID> <von> -> <nach> · verdict: <...>'), und digest.ts liest Lane und Urteil nur aus ihr; fehlt sie, zeigt der Einzeiler nur die ID statt zu raten."
context: |-
  Was falsch ist: core/mod/builtin/jaira-cockpit/hooks/digest.ts raet. laneOf nimmt den ersten Pfeil irgendwo im Bericht (Zeile 67), verdictOf das erste Wort wie 'reject' irgendwo (Zeile 77-84). Ein Bericht, der 'reject? reason' oder 'nichts zurueckgewiesen' erwaehnt, wird zu '-> signoff · rejected'.
  Ausloeser: Maintainer 31.10.2026: 'why is our cockpit wrong this often? is it because haiku?' - Nein: kein Modell, reine Regex. Haiku klassifiziert nur die Threads.
  Bekannt: X6HB6V (signoff) korrigiert die Lane. Das Urteil bleibt falsch. Die Schlusszeile muss in jaira-worker, jaira-dispatcher und jaira-role-lane stehen.
definition-of-done: "Worker-, Dispatcher- und Lane-Rolle schreiben die feste Schlusszeile"
tags:
  - tui
blocked-by: []
related:
  - 01M5MNWVKBC8Q193EK7YX6HB6V
commits: []
created-at: 2026-10-31T01:24:29Z
updated-at: 2026-11-22T16:37:44Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1092384
claimed-at: 2026-11-22T16:25:37Z
outcome-what: "digest.ts no longer guesses: without the closing status line the one-liner shows the ticket id (plus an open question or needs-you line); jaira-role-lane now ends its report with the same status line as worker and dispatcher"
outcome-why: "laneOf/verdictOf picked the first arrow or verdict word anywhere in a report, so 'reject? reason', 'nothing rejected' or 'all the way to signoff' became a wrong lane or verdict"
outcome-resolves: "the cockpit one-liner reads lane and verdict only from the fixed closing line; checked: concurrency n/a (pure function), error paths (failure detection unchanged, keeps the report whole), unusual input (empty/no ticket passes whole; CRLF and frame tags covered by the status-line tests), permissions n/a (fewer regexes over agent text)"
---

# Die Einzeiler des Cockpits lesen Lane und Urteil aus einer festen Schlusszeile statt aus dem ersten passenden Wort

## Definition of Done

- [x] jaira-role-lane schreibt die feste Schlusszeile wie Worker und Dispatcher (55CS2M); diese zwei haben sie schon
  proof: core/role/builtin/jaira-role-lane/SKILL.md: Boundaries, last item (status line)
- [x] digest.ts: ohne Schlusszeile zeigt der Einzeiler nur die ID - laneOf/verdictOf raten nicht mehr aus dem Fliesstext; ein Bericht mit 'reject' im Text ergibt kein 'rejected' (Test)
  proof: core/mod/builtin/jaira-cockpit/hooks/digest.ts digest(): ': [ticket]'; test 'no status line, only the ticket: *' (9 cases) in tests/cockpit.test.ts
- [x] Release-Notiz unter core/release/notes/
  proof: core/release/notes/MVR4PY.md
- [x] acceptance rule: a one-liner with only the ID is better than one with a wrong lane or verdict
  proof: tests/cockpit.test.ts guessCases: every free-text case expects the bare id

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 16:25 · Alexander Sacharov** — Alexander Sacharov took this ticket over from berkc
- **2026-11-22 16:25 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 16:32 · Alexander Sacharov** — decision (Alex 22.11.): Ticket auf den Rest verkleinert. Die Schlusszeile in worker/dispatcher und das Lesen in digest.ts brachte 55CS2M (e7046a407). Offen: (1) Schlusszeile in jaira-role-lane, (2) ohne Schlusszeile nur die ID, kein Raten. DoD 1-2 umformuliert.
- **2026-11-22 16:33 · Alexander Sacharov** — decision (worker): without a status line, question:/needs you: are still read from the report text. They are not lane or verdict, and dropping them would hide what waits for a person. Only laneOf/verdictOf go. Planning lane skipped: the Options checklist does not tick it, and the method is in the DoD.
