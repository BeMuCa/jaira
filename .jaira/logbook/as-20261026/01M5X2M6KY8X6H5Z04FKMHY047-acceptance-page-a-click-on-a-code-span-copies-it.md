---
id: 01M5X2M6KY8X6H5Z04FKMHY047
title: "Acceptance page: a click on a code span copies it; line breaks in step text show"
status: done
ready: true
goal: "A person running an acceptance page copies the exact text a step tells them to send with one click, and a step written as numbered sub-steps shows one sub-step per line."
context: |-
  Steps that tell a person to send a sentence put it in backticks, so it renders as <code>. Copying it means selecting by hand, and a selection easily misses the first or last character.
  Steps with several sub-steps (1) … 2) … 3) …) render as one paragraph: newlines in the step text are dropped by fmt(). Long steps become hard to follow.
  Found while a person worked through a 24-step acceptance page.
definition-of-done: "Clicking a <code> span copies its text to the clipboard and outlines it briefly in the ok colour"
tags:
  - acceptance
blocked-by: []
related: []
commits:
  - 5dd939eb961806df0eee523b7e8cfc9d02045512
created-at: 2026-10-26T21:42:22Z
updated-at: 2026-10-26T21:44:02Z
outcome-what: Click on a code span copies it; newlines in step text become line breaks
outcome-why: Copying sentences to send by hand missed characters; numbered sub-steps ran together
outcome-resolves: "template.html: cursor:copy, click handler, fmt newline"
---

# Acceptance page: a click on a code span copies it; line breaks in step text show

## Definition of Done

- [x] Clicking a <code> span copies its text to the clipboard and outlines it briefly in the ok colour
  proof: core/role/builtin/jaira-role-acceptance/template.html click handler on code
- [x] A newline in a step's text renders as a line break
  proof: template.html fmt(): \n -> <br>

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

