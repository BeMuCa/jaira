---
id: 01M7WWZWSDHK37DGGHJFDF0YD6
title: "Die Schleifenbremse steht am Ticket: was ein Blocker ist und dass nach einem Fix kein neuer Bughunt kommt"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Ticket traegt ein Blocker-Kriterium und die Option stop-after-fix; testing, bughunt und review lesen sie, und das Cockpit bietet Explain nur bei einem blocker an."
context: |-
  Heute: ein Teamlead schreibt mehrmals am Tag von Hand "fix -> test -> kein Bughunt -> review, Blocker nur bei falschem Ergebnis oder Absturz". Das hat jedes Mal sofort konvergiert, steht aber nur in Notizen und gilt fuer das naechste Ticket nicht.
  Ohne Schwelle ist jeder Randfall ein Blocker; jeder Fix ist neuer Code, den der naechste Bughunt jagt (das Ziel wandert mit).
  review-gaps hat immer 2-6 note:-Zeilen; das Cockpit zeigt Explain schon bei note: allein (core/mod/builtin/jaira-cockpit/hooks/board.ts:196 explainable/namesGap), deshalb hat jedes Ticket in signoff ein Explain.
  Die Schwelle gehoert in DoD/Plan, gesetzt von pre-process (siehe Ticket zu Randfaellen vor dem Code); der Mensch kann sie aendern.
  Verwandt: H6CD33 (max-rounds einer Schleifen-Lane).
definition-of-done: "Ticket-Felder fuer Blocker-Kriterium und stop-after-fix; jaira show --for-lane gibt sie testing, bughunt und review mit (Test)"
tags:
  - gates
blocked-by: []
related:
  - 01M6M7EA5EG5TWJSTVNQH6CD33
commits: []
created-at: 2026-11-20T16:35:17Z
updated-at: 2026-11-20T16:39:15Z
updated-by: Alexander Sacharov
---

# Die Schleifenbremse steht am Ticket: was ein Blocker ist und dass nach einem Fix kein neuer Bughunt kommt

## Definition of Done

- [ ] Ticket-Felder fuer Blocker-Kriterium und stop-after-fix; jaira show --for-lane gibt sie testing, bughunt und review mit (Test)
- [ ] Mit stop-after-fix laeuft nach einem Fix-Durchgang nur der Testlauf, kein neuer Bughunt
- [ ] Das Cockpit zeigt Explain nur, wenn review-gaps eine blocker:-Zeile hat oder das Urteil einen Vorbehalt nennt; reine note:-Zeilen stehen als Liste ohne Knopf (Test in board.test.ts)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

