---
id: 01M5M7JTVR8G7DZ1MGHJQ63EEV
title: "Tickets im Cockpit einzeln abarbeiten: Ticket aufklappen, Details rechts, Accept/Deny darunter"
status: backlog
ready: false
creator: BeMuCa
goal: "A person works off waiting tickets in the cockpit one by one: expanding a ticket opens a pane with id, one-liner, how it was implemented, information and open questions, with Accept and Deny beneath."
context: |-
  Berk, 2026-10-22, after walking through five finished tickets in chat (id + one-liner + information + how it was implemented, then the open question, or accept / explain / fix):
  This is also a better way to go about working off the tickets in the jaira-cockpit I think, we should add it or at least think about how we could improve the cockpit to make it possible. For example making it possible to expand the tickets, where then a mod window opens on the right with all of these things on it (ticket id + oneliner + ... everything I just mentioned.) and the buttons beneath it to accept, or deny. If you have further questions you can just ask about it in the chat.
  What the cockpit has today: a band with Accept / Give back / Explain / Discuss per waiting ticket (core/mod/builtin/jaira-cockpit/hooks/board.ts), and Explain writes a note with options onto the ticket. What is missing is one place per ticket that shows the whole picture before the button.
definition-of-done: "Brainstorm with Berk has decided the pane's content, where the information comes from (ticket fields vs. a generated summary) and how Accept/Deny map onto signoff moves"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T11:15:50Z
updated-at: 2026-10-26T19:08:32Z
assignee: BeMuCa
updated-by: BeMuCa
---

# Tickets im Cockpit einzeln abarbeiten: Ticket aufklappen, Details rechts, Accept/Deny darunter

## Definition of Done

- [ ] Brainstorm with Berk has decided the pane's content, where the information comes from (ticket fields vs. a generated summary) and how Accept/Deny map onto signoff moves
- [ ] Follow-up tickets for the implementation are filed

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

