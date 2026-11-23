---
id: 01M851Y62KHD9DVD2GDAM31BYV
title: "Give back tippt die Entscheidung in einen Worker-Tab, der an einem Freigabe-Dialog haengt"
status: done
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
commits:
  - 3c78621803b6cf25094dd42b88336d1ef202f97c
  - 9986c069459e854ef1adbd1fa26c82e1233c25be
  - cdf607c973ec696915cbd60cd2ba36e93376c08b
  - 191070c3679617c508fc83baa99ec7aeee5a3ae8
  - 14682393b12f2cc69fcb39bed3716cf01d3a08f8
  - a2d8aa5d3f2de2626336f97dceb8b5888668fce9
created-at: 2026-11-23T20:35:39Z
updated-at: 2026-11-23T21:21:30Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-93876
claimed-at: 2026-11-23T21:06:44Z
outcome-what: "Give back never types into or presses Enter in a worker waiting at an approval dialog: the carrier tab's status and the recorded pane's own status (herdr pane get, before the text and before the Enter) are checked; a dialog before the text: nothing sent, the tab or pane named; a dialog after it: no Enter, recorded as handed to the old worker, the person told to answer the dialog and press Enter only if the decision still waits"
outcome-why: "send-text + Enter into a pane at an approval dialog answered the dialog for the person, unseen"
outcome-resolves: "M31BYV: decision no longer lands in an approval dialog"
review-summary: "full review of 3c7862180..14682393b (no earlier review). Give back now checks the carrier tab (agent_status=blocked -> nothing typed, nothing closed, no spawn) and asks Herdr about the recorded pane itself (paneBlocked via 'pane get', mapped in herdr-hc.sh too) before send-text and again before Enter; a dialog after the text reports the decision as the old worker's and presses no Enter. DoD 1-3 met: three tests in tests/board.test.ts cover blocked tab, blocked pane behind an idle tab, dialog between text and Enter; texts workerBlocked/dialogBeforeEnter name tab or pane; release note is one line. Ran: claude plugin test (525 pass), go test ./core/shipped ./core/role ./core/release ok."
round: 4
back-from: testing
back-by: lane
back-count: 2
review-gaps: |-
  note: edge: carrierTab (sessions.ts:105) returns only the first <id>/worker or <id>/dispatch tab - with both open and the second at a dialog, a cold-cache give back closes the first and spawns a new worker beside the blocked one; fix: treat the ticket as blocked when any carrier tab in the list is blocked
  note: limit: a dialog opening in the milliseconds between the last pane get and send-text/Enter is not caught - check-then-act, said in the in-progress note
  note: elsewhere: register.tsx:1784 - a failed herdr tab list reads as no tab open and can start a second worker beside a live one; there before this ticket (82aa5654a3), new ticket
test-verdict: "pass: round 3 - cockpit plugin test rc=0 (525 pass), go test ./core/role ./core/shipped rc=0 (manifest regenerated); DoD 1-3 re-verified; dialog after the text: no Enter, note/toast/teamlead line say the old worker has the decision"
hunt-verdict: "pass: bughunt done at 14682393b, 0 findings (round 3); round 1 pre-existing (failed tab list -> second worker) left as a note for the person"
review-verdict: "accept: DoD 1-3 met with tests; no blocker - one edge with two carrier tabs, one known limit, one pre-existing fault elsewhere"
review-round: 4
review-verdict-gaps: |-
  note: edge: carrierTab (sessions.ts:105) returns only the first <id>/worker or <id>/dispatch tab - with both open and the second at a dialog, a cold-cache give back closes the first and spawns a new worker beside the blocked one; fix: treat the ticket as blocked when any carrier tab in the list is blocked
  note: limit: a dialog opening in the milliseconds between the last pane get and send-text/Enter is not caught - check-then-act, said in the in-progress note
  note: elsewhere: register.tsx:1784 - a failed herdr tab list reads as no tab open and can start a second worker beside a live one; there before this ticket (82aa5654a3), new ticket
review-check: "1. cd core/mod/builtin/jaira-cockpit && claude plugin test .  2. see 525 pass, including the three give back tests named 'approval dialog' / 'dialog that opens in the second'  3. live (Herdr): start a worker for a ticket, let it raise a permission dialog, press Give back on that ticket in the cockpit  4. nothing appears in the worker tab, no new tab opens, and the toast says the worker waits at an approval dialog in tab <ID>/worker"
---

# Give back tippt die Entscheidung in einen Worker-Tab, der an einem Freigabe-Dialog haengt

## Definition of Done

- [x] Give back tippt nichts und drueckt kein Enter in einen Tab mit agent_status=blocked; Test in core/mod/builtin/jaira-cockpit/tests
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1794 (tab) and core/mod/builtin/jaira-cockpit/hooks/register.tsx:1850 paneBlocked before send-text and before Enter; tests core/mod/builtin/jaira-cockpit/tests/board.test.ts:2418, core/mod/builtin/jaira-cockpit/tests/board.test.ts:2433, core/mod/builtin/jaira-cockpit/tests/board.test.ts:2451
- [x] Die Person sieht in der Meldung von Give back, dass der Worker an einem Freigabe-Dialog wartet und in welchem Tab
  proof: hooks/texts.ts workerBlocked / dialogBeforeEnter name the tab or pane; asserted in core/mod/builtin/jaira-cockpit/tests/board.test.ts:2418, core/mod/builtin/jaira-cockpit/tests/board.test.ts:2433, core/mod/builtin/jaira-cockpit/tests/board.test.ts:2451
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
- **2026-11-23 20:51 · Alexander Sacharov** — optimize: core/role TestEveryInstalledFileIsAShippedVersion was red - changed cockpit files need their version in core/shipped/shipped.sha256; regenerated with scripts/shipped-history.sh --write. Gate for this ticket = plugin test AND go test ./...
- **2026-11-23 20:54 · Alexander Sacharov** — waiting: bughunt round 2 (background subagent) since 21:54
- **2026-11-23 20:54 · Alexander Sacharov** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1815 - class: the second paneBlocked check runs after send-text already delivered the decision, yet reports it as not handed over (atDialog -> handedTo('none') note, toast, teamlead line 'worker ... type the decision into it'); example: under Herdr App (HC_TERMINAL_ID set, bin = herdr-hc.sh) 'pane send-text' is 'hc send', which pastes AND presses Enter (herdr-hc.sh:116-119), so the old worker has already received and submitted the decision; it then opens an approval dialog within the 1s sleep (likely: it starts acting on the decision), and the ticket gets 'given back to nobody (no worker started) ... no Enter pressed', the person is told to type the decision again (worker gets it twice) and the teamlead is told nobody carries it. On plain Herdr the text sits pasted in the input under the dialog and the same 'type the decision into it' makes a duplicate too; test 'a dialog that opens in the second between the text and its Enter gets no Enter' (board.test.ts) asserts this note; blame: cdf607c97 (names M31BYV); read
- **2026-11-23 20:55 · Alexander Sacharov** — hunt report: 1 finding, 1 change (blocker), note on M31BYV | read as a diff since 9986c0694: register.tsx, texts.ts, board.test.ts, release note M31BYV.md; read whole (new to the list): core/shipped/shipped.sha256 (current register.tsx/texts.ts hashes present; go test ./core/shipped ./core/role green); claude plugin test: 525 pass | not checkable here: the ms window between paneBlocked and send-text (a dialog opening there still receives the pasted keys) - needs live Herdr, inherent to check-then-act
- **2026-11-23 20:55 · Alexander Sacharov** — bughunt: done at 191070c36 - 1 finding (1 change, 0 pre-existing) in 5 modules, follow-up hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/M31BYV.md core/shipped/shipped.sha256
- **2026-11-23 20:55 · Alexander Sacharov** — testing: class: after send-text the decision has arrived, whatever happens next; any report after that point that says 'nobody / type it again' is wrong; example: bughunt round 2 (register.tsx:1815, herdr-hc.sh:116-119); fix: keep the check before the Enter (plain Herdr: an Enter into the dialog would answer it) but report handedTo old - text is in its pane, Enter not pressed, answer the dialog, then press Enter there only if the decision still waits in its input. Kept rather than dropped: on the working path a dialog within that 1s is likely.
- **2026-11-23 20:59 · Alexander Sacharov** — waiting: bughunt round 3 (background subagent) since 21:59
- **2026-11-23 21:02 · Alexander Sacharov** — hunt report: 0 findings, nothing filed | read as a diff since 191070c36: register.tsx (givingBack whole function around the hunk, paneBlocked), texts.ts, board.test.ts, shipped.sha256; release note unchanged since 191070c36 | ran: claude plugin test . in core/mod/builtin/jaira-cockpit (525 pass, 0 fail); new shipped.sha256 hashes match HEAD files; real Herdr/hc timing of a dialog between text and Enter not reproducible here
- **2026-11-23 21:02 · Alexander Sacharov** — bughunt: done at 14682393b - 0 findings in 5 modules, follow-up hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/M31BYV.md core/shipped/shipped.sha256
- **2026-11-23 21:02 · Alexander Sacharov** — waiting: review lane (fresh agent, run-lane.sh tab) since 22:02
- **2026-11-23 21:07 · Alexander Sacharov** — review: done at 14682393b
- **2026-11-23 21:07 · Alexander Sacharov** — report: Give back never types into or presses Enter in a worker at an approval dialog (tab status + herdr pane get before text and before Enter); toast/note/teamlead line name the tab or pane | signoff, context 158k of 1000k, cache warm 60 min after last answer | 3 review notes (edge: two carrier tabs; limit: ms window; elsewhere: failed tab list -> second worker, pre-existing) - person decides; live check in Herdr per review-check; branch not pushed
- **2026-11-23 21:12 · Alexander Sacharov** — Accepted from signoff: angenommen:
- note: edge: carrierTab (sessions.ts:105) returns only the first <id>/worker or <id>/dispatch tab - with both open and the second at a dialog, a cold-cache give back closes the first and spawns a new worker beside the blocked one; fix: treat the ticket as blocked when any carrier tab in the list is blocked
- note: limit: a dialog opening in the milliseconds between the last pane get and send-text/Enter is not caught - check-then-act, said in the in-progress note
- note: elsewhere: register.tsx:1784 - a failed herdr tab list reads as no tab open and can start a second worker beside a live one; there before this ticket (82aa5654a3), new ticket
- **2026-11-23 21:13 · Alexander Sacharov** — cockpit: accepted by the person, signoff -> done
