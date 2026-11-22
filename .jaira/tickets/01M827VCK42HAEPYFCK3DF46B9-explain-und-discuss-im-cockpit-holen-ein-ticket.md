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
updated-at: 2026-11-22T18:59:36Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1159000
claimed-at: 2026-11-22T18:52:19Z
outcome-what: "Explain and Discuss in the cockpit pull a ref-only or stale ticket through ticketHere before the chat opens; a copy pulled into the main checkout is removed once the chat has done its part (refreshBoard's chats loop), or at once when no chat could be opened"
outcome-why: "The chat writes its note with jaira note, which refuses a ticket that is only on its ref; Explain/Discuss failed on any ticket without a worktree, as Accept did before 76ERHH"
outcome-resolves: "DF46B9: Explain/Discuss on ref-only or stale tickets"
review-summary: "core/mod/builtin/jaira-cockpit/hooks/register.tsx:136-140 keeps two maps under one label (chats: since, leaves: the copy's leave) that are set and deleted together; make chats one Map<label, { since, leave }> and drop leaves and the opened() closure at :946"
round: 2
back-from: critique
back-by: lane
back-count: 1
---

# Explain und Discuss im Cockpit holen ein Ticket ohne Worktree oder mit veralteter Kopie zuerst, wie Accept seit 76ERHH

## Definition of Done

- [x] Explain und Discuss auf ein Ref-only-Ticket und auf eines mit veralteter Kopie: holen es zuerst über dieselbe Hilfe wie Accept (76ERHH), eine Kopie im Haupt-Checkout wird danach entfernt (Test)
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx: openChat pulls via ticketHere; test core/mod/builtin/jaira-cockpit/tests/board.test.ts:4316 'Explain and Discuss on a ticket not on the disk the chat writes from' (ref-only/behind × explain/discuss, copy cleaned once the chat is done; without Herdr too)
- [x] Ein verweigerter pull (fremder assignee) stoppt mit Meldung ohne --steal, wie bei Accept (Test)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts:4360 '<kind>: a refused pull stops the press with its reason, opens no chat and never steals'
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/DF46B9.md:1; core/shipped/shipped.sha256 regenerated (scripts/shipped-history.sh --write)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 18:52 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 18:53 · Alexander Sacharov** — decision: Explain/Discuss pull through ticketHere (76ERHH) before the chat starts. The copy in the main checkout cannot go when the press ends, as with Accept: the chat writes its note (jaira note = Store.Mutate) later and that refuses a ref-only ticket. So leave() runs when the chats loop in refreshBoard finds the chat done (its fresh note, or for Discuss the question gone). Without Herdr the chat is registered in the same map, so the same loop takes the copy away. Known limit: the leave lives in memory - a cockpit restarted, or a chat closed without a note, leaves the copy.
- **2026-11-22 18:58 · Alexander Sacharov** — critique: chats and leaves in register.tsx are two maps over the same label, always written together - fold leave into the chats entry ({ since, leave }), keeping the earlier leave on a second press, and inline opened(). The non-Herdr path now registering a chat (so a Discuss decision there is offered to read too) is intended: it is what lets the same loop take the copy away.
