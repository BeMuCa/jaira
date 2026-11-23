---
id: 01M851Y62KHD9DVD2GDAM31BYV
title: "Give back tippt die Entscheidung in einen Worker-Tab, der an einem Freigabe-Dialog haengt"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Give back im Cockpit schreibt nie in einen Worker-Tab, in dem gerade ein Freigabe-Dialog offen ist (Herdr agent_status=blocked); die Person erfaehrt stattdessen, dass der Worker auf ihre Freigabe wartet."
context: |-
  Was falsch ist: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1798 (Give back, alter Worker wird wiederverwendet) macht 'herdr pane send-text' und danach 'send-keys enter' in die Pane des alten Workers.
  Haengt der Worker dort an einem Freigabe-Dialog (z.B. Bash-Befehl), landet der Text im Dialog und das Enter beantwortet ihn - im Namen der Person, ohne dass sie es sieht.
  Herdr meldet diesen Zustand als agent_status=blocked. spawn.sh prueft genau das schon vor seinem send-keys und bricht ab (spawn.sh: 'claude blocked' -> exit 1).
  register.tsx:1790 kennt nur working ('nie ersetzen'), nicht blocked.
  Gefunden vom Bughunt zu RJFA19 am 23.11.2026 (pre-existing, nicht durch RJFA19 entstanden). RJFA19 zeigt den Zustand jetzt im Worker-Tab an.
  Offen fuer den Plan: bei blocked nichts tippen und die Person auf den Dialog hinweisen, oder einen neuen Worker starten - der alte haengt dann aber weiter.
definition-of-done: Give back tippt nichts und drueckt kein Enter in einen Tab mit agent_status=blocked; Test in core/mod/builtin/jaira-cockpit/tests
tags:
  - cockpit
blocked-by: []
related: []
commits: []
created-at: 2026-11-23T20:35:39Z
updated-at: 2026-11-23T20:49:56Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-41227
claimed-at: 2026-11-23T20:36:42Z
outcome-what: "Give back never types into or presses Enter in a worker waiting at an approval dialog: the carrier tab's status and the recorded pane's own status (herdr pane get, before the text and before the Enter) are checked; ticket note, teamlead line and toast name the waiting tab or pane; nothing replaced"
outcome-why: "send-text + Enter into a pane at an approval dialog answered the dialog for the person, unseen"
outcome-resolves: "M31BYV: decision no longer lands in an approval dialog"
review-summary: "none"
round: 3
back-from: testing
back-by: lane
back-count: 1
review-gaps: "none: no duplicate - sessions.ts:151 'waits for approval: answer the dialog in its tab' is the /cockpit workers line, a different message; label local used three times, kept; nothing orphaned"
test-verdict: "pass: plugin suite green (rc=0, 523 pass), DoD 1-3 verified in the tree; Herdr not runnable here, behaviour exercised through the mod's test world"
hunt-verdict: "fail: guard reads the carrier tab's status but types into record.pane, which may be another (blocked) pane; and nothing is re-read before the Enter"
---

# Give back tippt die Entscheidung in einen Worker-Tab, der an einem Freigabe-Dialog haengt

## Definition of Done

- [x] Give back tippt nichts und drueckt kein Enter in einen Tab mit agent_status=blocked; Test in core/mod/builtin/jaira-cockpit/tests
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1794 (tab) and core/mod/builtin/jaira-cockpit/hooks/register.tsx:1844 paneBlocked before send-text and before Enter; tests core/mod/builtin/jaira-cockpit/tests/board.test.ts:2418, core/mod/builtin/jaira-cockpit/tests/board.test.ts:2433, core/mod/builtin/jaira-cockpit/tests/board.test.ts:2451
- [x] Die Person sieht in der Meldung von Give back, dass der Worker an einem Freigabe-Dialog wartet und in welchem Tab
  proof: hooks/texts.ts workerBlocked names the tab or pane; asserted in core/mod/builtin/jaira-cockpit/tests/board.test.ts:2418 and core/mod/builtin/jaira-cockpit/tests/board.test.ts:2433
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/M31BYV.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 20:36 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/48662.sock
- **2026-11-23 20:36 · Alexander Sacharov** — decision: blocked tab = like working: never replaced, nothing typed; handedTo none with the tab label, no second worker (the old one still hangs at its dialog). Settled by DoD 2, no question open.
- **2026-11-23 20:41 · Alexander Sacharov** — critique: blocked block split the working comment from its code (register.tsx:1786-1799); moved above it. Shape otherwise kept: early return like place.stop, reuses handedTo + say, no new helper.
- **2026-11-23 20:43 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1790 - class: the blocked guard reads one 'herdr tab list' snapshot and is never re-read before the keystrokes that follow; the working path (1799-1810) still types send-text into a worker at work and presses Enter 1s later (1809-1810); example: worker is 'working' at the tab list, raises a permission dialog during the send-text or the 1s sleep, the Enter lands in the dialog and answers it with its default for the person, unseen; marked change because the new check sits on this path and was meant to stop exactly this, though the typing lines are older (05a8d99fc3, 8fc5821f63); an edge, a timing window not a deterministic fault; read
- **2026-11-23 20:43 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1790 - class: the guard checks one object and the text goes to another: it reads agent_status of carrierTab (first tab labelled <id>/worker or <id>/dispatch, sessions.ts:105) but types into record.pane of carrierOf (newest worker/dispatch record of any session for the ticket, sessions.ts:51, not matched to that tab); example: spawned tab MG677R/worker idle, and a newer '/jaira-worker MG677R' session in another pane sits at an approval dialog - the check sees idle, chooseWorker may reuse, send-text and Enter go to the blocked pane; also two carrier tabs (worker and dispatch) where the blocked one is listed second; edge; read
- **2026-11-23 20:43 · Alexander Sacharov** — bughunt: pre-existing: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1784 - class: a failed 'herdr tab list' (non-zero exit or unparsable JSON, tabsOf returns []) is treated as 'no tab open'; example: tab list fails while MG677R/worker is alive (blocked or working) - the blocked/working guards are skipped, chooseWorker says tabAlive false, the old tab is not closed (tab undefined) and spawn.sh starts a second worker beside it; blame: 82aa5654a3 / 86da7c4354 (before the ticket); read
- **2026-11-23 20:43 · Alexander Sacharov** — hunt report: 3 findings, 2 change (both edges, no blocker), 1 pre-existing, all as notes on M31BYV | read whole: givingBack/giveBack and helpers in hooks/register.tsx, hooks/sessions.ts, the changed texts.ts lines, the give-back tests in tests/board.test.ts, release note; rest of register.tsx (2709 lines) only around the give-back path | could not run Herdr here (Windows binary), so the race and pane-mismatch findings are confirmed by reading only
- **2026-11-23 20:43 · Alexander Sacharov** — bughunt: done at 9986c0694 - 3 findings (2 change, 1 pre-existing) in 4 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/M31BYV.md
- **2026-11-23 20:44 · Alexander Sacharov** — testing: class: the blocked check looks at a tab, the keys go to a pane - status of one object, keystrokes to another, read once; example: bughunt notes 1+2 (record.pane in another session's pane at a dialog, carrier tab idle; dialog raised during the 1s before Enter); fix: ask 'herdr pane get <record.pane>' (result.pane.agent_status, as spawn.sh:549) right before send-text and again right before send-keys enter; blocked at either -> no keys, handedTo none naming the pane; keep the tab check (it stops closing a blocked tab and spawning beside it). Test both points.
- **2026-11-23 20:49 · Alexander Sacharov** — in-progress round 2: blocked is now asked of the recorded pane itself (herdr pane get, as spawn.sh) before send-text and before Enter, besides the tab check. Remaining window: a dialog opening between the last pane get and the Enter (milliseconds) - no check closes that, accepted. Pre-existing bughunt finding (failed tab list -> second worker) left as a note for the person: other subject (tab list error path), not this goal.
