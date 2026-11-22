---
id: 01M827VCK42HAEPYFCK3DF46B9
title: "Explain und Discuss im Cockpit holen ein Ticket ohne Worktree oder mit veralteter Kopie zuerst, wie Accept seit 76ERHH"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Explain und Discuss funktionieren für ein Ticket, das nur auf seinem Ref liegt oder dessen Kopie veraltet ist, genauso wie Accept/Give back nach 76ERHH."
context: |-
  Was falsch ist: 76ERHH (22.11.2026) brachte Accept, Give back und Entscheidung im Cockpit bei: ein Ref-only- oder veraltetes Ticket erst pullen, eine Kopie im Haupt-Checkout danach entfernen.
  Explain und Discuss machen das nicht und scheitern bei einem Ticket ohne Worktree genauso wie früher Accept (Review-Notiz von 76ERHH, pre-existing).
  Code: core/mod/builtin/jaira-cockpit/hooks/register.tsx (Explain/Discuss), die Pull-Hilfe aus 76ERHH wiederverwenden statt neu bauen.
definition-of-done: "Explain und Discuss auf ein Ref-only-Ticket und auf eines mit veralteter Kopie: holen es zuerst über dieselbe Hilfe wie Accept (76ERHH), eine Kopie im Haupt-Checkout wird danach entfernt (Test)"
tags:
  - cockpit
blocked-by: []
related: []
commits: []
created-at: 2026-11-22T18:21:16Z
updated-at: 2026-11-22T19:09:03Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1159000
claimed-at: 2026-11-22T18:52:19Z
outcome-what: "Explain and Discuss in the cockpit pull a ref-only or stale ticket through ticketHere before the chat opens; a copy pulled into the main checkout is removed once the chat has done its part (refreshBoard's chats loop), or at once when no chat could be opened"
outcome-why: "The chat writes its note with jaira note, which refuses a ticket that is only on its ref; Explain/Discuss failed on any ticket without a worktree, as Accept did before 76ERHH"
outcome-resolves: "DF46B9: Explain/Discuss on ref-only or stale tickets"
review-summary: "none"
round: 3
back-from: testing
back-by: lane
back-count: 1
review-gaps: "none removed: ticketDir stays (still read by two callers at register.tsx:1210/1232); left MAIN/ALONE repeated in the new describe in board.test.ts beside the 76ERHH one - describe-scoped, hoisting them would touch the 76ERHH block for no behaviour"
test-verdict: "pass: claude plugin test (2.1.x) 368/368 RC=0 at 31f8c9e0d, the 8 new DF46B9 tests fail with the fix stashed; go build ./... RC=0, shipped-history --check RC=0, go test core/shipped core/role core/release core/mod RC=0. DoD 1-3 verified in the tree (register.tsx:936, board.test.ts:4316/4360, notes/DF46B9.md). Failure paths covered by the engine tests: refused pull (no spawn, no clean, message pane with --steal text), spawn.sh missing (copy cleaned at once). Not run live in Claude Code against a real ref-only ticket: the plugin runs only inside Claude Code; the jaira pull/note/git clean sequence is the one 76ERHH ships. Two writers: the change writes only through jaira pull and git clean of a path it pulled itself; jaira serialises the pull"
hunt-verdict: "fail: bughunt done at 31f8c9e0d, 2 change blockers - the copy pulled for a chat is left forever when the chat ends without a note or the session ends, and is removed under a non-Herdr Discuss between its note and its set/move"
---

# Explain und Discuss im Cockpit holen ein Ticket ohne Worktree oder mit veralteter Kopie zuerst, wie Accept seit 76ERHH

## Definition of Done

- [x] Explain und Discuss auf ein Ref-only-Ticket und auf eines mit veralteter Kopie: holen es zuerst über dieselbe Hilfe wie Accept (76ERHH), eine Kopie im Haupt-Checkout wird danach entfernt (Test)
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:949 openChat pulls via ticketHere; core/mod/builtin/jaira-cockpit/tests/board.test.ts:4321 'Explain and Discuss on a ticket not on the disk the chat writes from' (ref-only/behind x explain/discuss, with and without Herdr)
- [x] Ein verweigerter pull (fremder assignee) stoppt mit Meldung ohne --steal, wie bei Accept (Test)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts:4396 '<kind>: a refused pull stops the press with its reason, opens no chat and never steals'
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/DF46B9.md:1-2; core/shipped/shipped.sha256 regenerated, shipped-history --check RC=0
- [x] Die Kopie im Haupt-Checkout entfernt der Chat selbst als letzten Schritt (Pfad im Prompt von spawn.sh und texts.ts); das Cockpit entfernt sie nur, wenn kein Chat startete; bleibt sie nach dem Chat liegen, nennt das Band den Befehl (Test)
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:91 --copy + last step; core/mod/builtin/jaira-cockpit/hooks/texts.ts:25 removeCopy; tests core/mod/builtin/jaira-cockpit/tests/board.test.ts:4379 'a copy the chat left behind is named beside the band', 'a chat that cannot be opened takes the pulled copy away at once', core/role/role_test.go TestSpawnChatRemovesTheCopyItWasHanded

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 18:52 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 18:53 · Alexander Sacharov** — decision: Explain/Discuss pull through ticketHere (76ERHH) before the chat starts. The copy in the main checkout cannot go when the press ends, as with Accept: the chat writes its note (jaira note = Store.Mutate) later and that refuses a ref-only ticket. So leave() runs when the chats loop in refreshBoard finds the chat done (its fresh note, or for Discuss the question gone). Without Herdr the chat is registered in the same map, so the same loop takes the copy away. Known limit: the leave lives in memory - a cockpit restarted, or a chat closed without a note, leaves the copy.
- **2026-11-22 18:58 · Alexander Sacharov** — critique: chats and leaves in register.tsx are two maps over the same label, always written together - fold leave into the chats entry ({ since, leave }), keeping the earlier leave on a second press, and inline opened(). The non-Herdr path now registering a chat (so a Discuss decision there is offered to read too) is intended: it is what lets the same loop take the copy away.
- **2026-11-22 18:59 · Alexander Sacharov** — critique round 2: nothing left. Reuses ticketHere unchanged (register.tsx:550) and the existing chats loop; the one-map finding from round 1 is fixed in 31f8c9e0d.
- **2026-11-22 19:01 · Alexander Sacharov** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/register.tsx:944 (with :416-428, :1659) - class: a copy pulled into the main checkout for a chat is removed only when that chat leaves a fresh note, so every chat that ends without one keeps it there for good; example: an Explain tab closed (or its subagent failing) before it writes explain:, a Discuss chat left while the question stays, or the session ending with a chat open - chats is in-memory, session.end (:1659) does not run the leaves, and the next press finds the leftover 'behind' and keeps it (ticketHere stays for anything but ref-only, :563). The copy left there outranks the ref and shows the ticket in a lane it left (ticketHere's own comment, :546-547); read
- **2026-11-22 19:02 · Alexander Sacharov** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/register.tsx:418-428 (and :944, :960-963) - class: leave() deletes the copy while a chat is still writing through it; examples: (a) Discuss without Herdr has no tab, so the agent_status==='working' guard (:418) never holds - the 60s tick (:1581, void refreshBoard) that lands between the session's 'note decision:' and its 'set question=' / 'move --to in-progress' (texts.ts:168-170) sees the fresh note, runs git clean, and those writes fail with 'on its ref and not on your disk', the very refusal this ticket removes; (b) a second press on the same label while the first chat finishes: ticketHere sees 'current' and gets a no-op leave, refreshBoard deletes the label and removes the copy during spawn.sh's 120s, then opened() stores the no-op leave for a chat running without its file; (c) spawn.sh timing out (run -> undefined, :960) after the tab came up runs here.leave() under the running chat; read
- **2026-11-22 19:02 · Alexander Sacharov** — bughunt: done at 31f8c9e0d - 2 findings (2 change, 0 pre-existing) in 4 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/DF46B9.md core/shipped/shipped.sha256
- **2026-11-22 19:03 · Alexander Sacharov** — question: Who removes the copy a chat ran from? The cockpit cannot see when a chat wrote its last line (non-Herdr Discuss has no tab; a closed tab or ended session loses the in-memory leave) - both hunt findings come from that. Options: 1) (Recommended) The chat removes it itself as its last step: the cockpit passes the pulled path into the chat prompt (spawn.sh explain/discuss + texts.ts explainBySubagent/discussInThisChat) and only removes it when no chat could be opened; the leave tracking in refreshBoard goes again. Copy stays only if the chat dies. 2) The cockpit waits for the runner to be idle: Herdr tab not working or gone, non-Herdr subagent finished / this session's turn over, plus all leaves on session.end. More state, more timing edges. 3) The copy stays for chats: drop 'removed afterwards' from DoD 1 for Explain/Discuss; once the chat is done the band names the git clean line. Smallest, no race, but a copy stays in the main checkout.
- **2026-11-22 19:03 · Alexander Sacharov** — decision (teamlead 22.11.): Option 1 - der Chat (Explain/Discuss) entfernt die Kopie selbst als letzten Schritt, Pfad steht in seinem Prompt; das Cockpit entfernt sie nur, wenn der Chat gar nicht startet. Option 3 nicht: eine liegengebliebene Kopie im Haupt-Checkout verdeckt dem Cockpit die Lane des Worktrees. Bleibt eine Kopie nach Absturz liegen, zeigt das Band den Befehl zum Entfernen.
- **2026-11-22 19:04 · Alexander Sacharov** — decision (parent, answer 1): the chat removes the copy itself as its last step; the cockpit removes it only when the chat never started. Not 3: a copy left in the main checkout hides the worktree's lane. If a crash leaves a copy behind, the band shows the remove command.
- **2026-11-22 19:08 · Alexander Sacharov** — round 2 (answer 1): leave tracking in refreshBoard removed. The chat gets the copy's path (spawn.sh --copy, texts.ts removeCopy) and runs git clean itself as its last step. The cockpit removes the copy only when spawn.sh is missing (no chat). A failed spawn.sh (may have opened a tab on a timeout) does not remove it. A minute after a chat is done (COPY_CHECK_MS), a copy still on disk gets a band line with the git clean command, once. Known limit: a Herdr tab closed without a note, or a cockpit session that ends mid-chat, never counts as done, so no band line. An older installed spawn.sh refuses --copy - release note line '- !' says to run jaira update --roles.
