---
id: 01M2E85S75MEF7YJJRJ6C9QS8F
title: "Was die Tafel treibt, reist mit dem Binary"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer jaira installiert, bekommt nicht nur die Lanes, sondern auch das, was sie faehrt: die Rollen-Prompts und einen Weg, ueber Lane-Wechsel benachrichtigt zu werden"
context: |-
  Die Tafel reist mit dem Repository, die Dinge, die sie ANTREIBEN, reisen gar nicht.
  Stand 2026-09-13 existieren drei Stuecke ausschliesslich auf einem Rechner und in keinem git:
  - ~/.claude/skills/ mit sieben Rollen-Prompts (teamlead, dispatcher, role-lane, role-tester, role-pr, role-brainstorm, role-research)
  - ~/.jaira/settings.json mit remote und hook
  - ~/.jaira/notify.sh, das jeden Lane-Wechsel in eine Herdr-Benachrichtigung verwandelt
  Folge: ein Teamkollege klont, startet jaira, sieht dieselbe Tafel - und hat niemanden, der sie faehrt, und erfaehrt nichts, wenn sich etwas bewegt.
  Dieses Ticket ist die Klammer, es wird selbst nicht gebaut. Darunter haengen die Kinder, die je ein Stueck erledigen.
  Nicht dazu gehoert settings.json selbst: das sind die Vorlieben eines Rechners (welcher Remote, welcher Hook) und gehoert genau dorthin, wo es liegt. Ausgeliefert wird, was man einstellt, nicht die Einstellung.
definition-of-done: Alle Kinder dieses Tickets sind abgeschlossen; ein frischer Klon auf einem zweiten Rechner kommt ohne Handarbeit an Rollen und Benachrichtigung
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-13T20:43:39Z
updated-at: 2026-09-13T20:43:39Z
---

# Was die Tafel treibt, reist mit dem Binary

## Definition of Done

- [ ] Alle Kinder dieses Tickets sind abgeschlossen; ein frischer Klon auf einem zweiten Rechner kommt ohne Handarbeit an Rollen und Benachrichtigung

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

