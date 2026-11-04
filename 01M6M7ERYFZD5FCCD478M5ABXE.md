---
id: 01M6M7ERYFZD5FCCD478M5ABXE
title: "A lane can wait for a reply that arrives from outside the board, such as a chat message"
status: backlog
ready: false
creator: berkc
goal: "A lane declares wait-for: reply; a ticket in it stays until a reply is recorded on it, and the reply reaches the next agent as a note."
context: |-
  What is wrong today: a ticket that needs an answer from a person sits in the human lane until somebody opens the board. There is no way to say 'this lane waits for a reply that a script will bring in', e.g. a chat-bot message from the phone.
  Triggered by: ticket ENCFVZ (skill/script/receipt lanes). The maintainer listed wait-for: reply as a lane attribute on 04.11.2026; kept out of ENCFVZ's scope on purpose.
  Builds on ENCFVZ: the lane's script sends the question (with JAIRA_IDEMPOTENCY_KEY and a receipt, so it goes out once per round).
  Already decided in the ENCFVZ brainstorm note of 30.10.2026: jaira does not poll and runs no daemon. A listener the person runs (bot, cron) calls 'jaira note' and 'jaira move' on their own machine. That listener must accept messages only from the person's own chat id, or anyone who finds the bot can move tickets.
  Open: whether a reply may move the ticket forward (then a bot decides what today only a person may) or only back to rejects-to; and how the gate tells a recorded reply from an ordinary note.
definition-of-done: "A lane file may declare wait-for: reply, and a forward move out of it is refused until a reply is recorded on the ticket"
tags:
  - gates
blocked-by: []
related: []
follows: 01M67K7GEZF694985EHTENCFVZ
commits: []
created-at: 2026-11-04T21:29:19Z
updated-at: 2026-11-04T21:29:19Z
---

# A lane can wait for a reply that arrives from outside the board, such as a chat message

## Definition of Done

- [ ] A lane file may declare wait-for: reply, and a forward move out of it is refused until a reply is recorded on the ticket
- [ ] docs/ carry an example listener script that checks the sender before it writes anything; nothing in the binary talks to a chat service
- [ ] README lane format and a release note describe it

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

