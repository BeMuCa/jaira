---
id: 01M5M7JTVR8G7DZ1MGHJQ63EEV
title: "Tickets im Cockpit einzeln abarbeiten: Ticket aufklappen, Details rechts, Accept/Deny darunter"
status: human
ready: true
creator: BeMuCa
goal: "A person works off the tickets waiting for them one at a time in the cockpit: pressing a ticket in the band opens its pane with title, goal, how it was done, what the review found, the DoD and the open question, read from the ticket's existing fields, with Accept, Give back (Deny), Explain or Discuss beneath, each doing what the band's button does today."
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
updated-at: 2026-10-26T19:20:30Z
assignee: BeMuCa
updated-by: BeMuCa
question: |-
  1. Pane content: (a) existing fields only, long ones folded, mixed languages; (b) an agent writes a short card per ticket in your language, as in the chat walkthrough, at one model call per ticket; (c) (a) now and (b) later only if (a) reads badly. Recommended (c).
  2. The open question of a signoff ticket: (a) you read it out of the verdict and newest notes in the pane; (b) the review lane writes it into the existing question field, which the band already shows with Discuss (overlaps ZWK0PF); (c) a new field of its own. Recommended (b), as its own follow-up.
  3. Opening a ticket from the band: (a) the ticket's ID in the row becomes a button; (b) (a) plus digits 1-3 from an empty prompt, at the cost that a prompt starting with a digit presses the row instead. Recommended (a).
  4. After Accept or Deny in the pane: (a) the pane closes; (b) the pane moves on to the next waiting ticket. Recommended (b).
---

# Tickets im Cockpit einzeln abarbeiten: Ticket aufklappen, Details rechts, Accept/Deny darunter

## Definition of Done

- [ ] Brainstorm with Berk has decided the pane's content, where the information comes from (ticket fields vs. a generated summary) and how Accept/Deny map onto signoff moves
- [ ] Follow-up tickets for the implementation are filed

## Options

- [x] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-26 19:20 · BeMuCa** — brainstorm: what is missing (read at 13aaa578, paths under core/mod/builtin/jaira-cockpit/):
- The band row draws only ID | lane | verdict. The title is in Waiting but is never drawn (hooks/register.tsx:1203).
- The cockpit reads none of goal, outcome what/why/resolves, review summary/gaps/check or the DoD proofs. Waiting carries 8 fields only (hooks/board.ts:9-24).
- A ticket pane exists (`cockpit-ticket`, hooks/register.tsx:70). It opens only by typing `/cockpit ticket <id>` (hooks/register.tsx:1137) and shows reports and usage, no content and no buttons.
- A signoff ticket has no field for its open question. Today it sits in the verdict's prose or a driver note. `question` is empty on all 11 signoff tickets.
- The choice part exists already: an Explain chat's `explain:` note becomes digit options in the explanation pane (hooks/register.tsx:1318-1368). It is offered only when the review names a reservation or a gap (hooks/board.ts:192-197).
Proposal:
- The row's ID becomes a Button. It opens the existing `cockpit-ticket` pane with focus.
- Where the pane sits is the surface's choice: docked beside the transcript in fullscreen from 110 columns, otherwise inline above the prompt.
- Content: one `jaira show --json`, run in the ticket's directory. Lines in order: title, goal, outcome what/why/resolves, verdict, DoD with proofs, open question.
- The open question is the human ticket's `question`, else the newest unanswered `explain:` note with its options inline.
- Folded: review summary, gaps, check (up to 7118 characters), newest notes, reports and usage.
- No new field and no agent call are needed for this.
- Accept = existing accept(): signoff -> done, dry run then --force, then afterAccept.
- Deny = existing Give back: reason pane, move to in-progress, note "Given back from <lane>: ...", hand to a worker.
- "explain more" = Explain. "fix something" = Give back. "other" = Own answer.
- A human ticket gets no Accept (acceptTo is null). It gets Discuss and Give back.
- Every button also gets a letter hotkey (a/g/e/d), so a scrolled card still works.
Options:
- A: fields only, mod-only, one ticket.
- B: A, plus the review lane writes the signoff open point into `question` (overlaps ZWK0PF).
- C: an agent writes a card in the person's language, one model call per ticket.
Recommended: A now (follow-up F1), B as its own follow-up (F2), C only if A reads badly.
Follow-ups:
- F1, one ticket, mod only: row button, the card from `jaira show --json` with long parts folded, Accept/Give back/Explain/Discuss in the pane with letter hotkeys, explain options inline, next ticket or close after a press, `/cockpit ticket` kept, tests and a release note.
- F2, own ticket, not mod-only: where a signoff ticket's open question lives (option B). Touches the review lane prompt and role texts; coordinate with ZWK0PF.
Not verified: that the docked pane sits on the right; the plugin API only says "beside".
