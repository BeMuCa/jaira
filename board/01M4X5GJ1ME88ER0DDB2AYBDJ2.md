---
id: 01M4X5GJ1ME88ER0DDB2AYBDJ2
title: Nach jaira logbook bleibt keine Arbeitskopie des Tickets und kein Launcher-Eintrag davon liegen
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Ist ein Ticket abgelegt und seine Arbeit gepusht, bietet jaira an, die Arbeitskopie (git worktree) des Tickets und ihren Launcher-Eintrag zu entfernen, und entfernt nur, was sicher weg kann"
context: |-
  Was heute falsch ist: Agenten legen pro Ticket eine Arbeitskopie an (git worktree add -b feat/<ID>-... /home/berk/git/jAIra-<ID>). Nach dem Merge bleibt sie liegen, und niemand raeumt sie weg.
  Stand am 14.10.: 7 solche Ordner lagen herum (jAIra-D0SAHM, -E70WYN, -PNC854, -9WE1QK, -J1ZVT9, -EBTX6X, -PR). Ihre Commits waren alle schon in PR #74.
  Jeder Ordner haelt veraltete Ticket-Kopien und eine eigene .jaira/milestones/. Wer dort jaira startet, sieht einen alten Stand (siehe K6BGB0).
  Jeder Ordner, in dem jaira einmal lief, steht als eigenes Board in ~/.jaira/projects.json und damit im Launcher (am 14.10. z.B. jAIra-E70WYN).
  Berk am 14.10.: 'Erstell Tickets bzgl der Worktree aufraeum Thematik. Wahrscheinlich wird es wie dein L12 ablaufen'. L12 war: nach dem Merge aufraeumen, und 'jaira logbook' bietet an, Arbeitskopie und Launcher-Eintrag zu entfernen.
  Grenzen: nie die Haupt-Arbeitskopie entfernen. Nur entfernen, wenn die Arbeitskopie sauber ist (untracked .jaira/milestones/ zaehlt nicht) und ihr Branch gepusht ist.
  Verwandt: JDXV60 (ein Ticket in zwei Arbeitsbaeumen), TS6YCJ (parallele Agenten im selben Arbeitsbaum), PWAF22 (Wegwerf-Boards in 'jaira projects').
definition-of-done: "Nach 'jaira logbook <id>' in einer Arbeitskopie (git worktree, nicht die Haupt-Arbeitskopie), die sauber ist und deren Branch gepusht ist, fragt jaira, ob die Arbeitskopie weg soll; ja entfernt sie mit 'git worktree remove' und nimmt ihren Eintrag aus ~/.jaira/projects.json; Test"
tags:
  - cli
blocked-by: []
related:
  - 01M2HW4788Z5G0M05G3PJDXV60
  - 01M263S8VMXMVQDFD3BTTS6YCJ
  - 01M2HX6SYCFB27V26R5APWAF22
commits: []
created-at: 2026-10-14T12:17:06Z
updated-at: 2026-10-14T12:17:06Z
---

# Nach jaira logbook bleibt keine Arbeitskopie des Tickets und kein Launcher-Eintrag davon liegen

## Definition of Done

- [ ] Nach 'jaira logbook <id>' in einer Arbeitskopie (git worktree, nicht die Haupt-Arbeitskopie), die sauber ist und deren Branch gepusht ist, fragt jaira, ob die Arbeitskopie weg soll; ja entfernt sie mit 'git worktree remove' und nimmt ihren Eintrag aus ~/.jaira/projects.json; Test
- [ ] Ist die Arbeitskopie nicht sauber oder ihr Branch nicht gepusht, entfernt jaira nichts und sagt, was fehlt; Test
- [ ] Ohne Terminal (Agent, --json) fragt jaira nicht; ein ausdruecklicher Schalter entfernt, und der jaira-role-pr-Skill nennt ihn fuer die Zeit nach dem Merge; Test
- [ ] Eine Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

