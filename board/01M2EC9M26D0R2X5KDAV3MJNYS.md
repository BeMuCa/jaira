---
id: 01M2EC9M26D0R2X5KDAV3MJNYS
title: Das Hook-Beispiel liest die Lane-Rollen aus dem Board statt sie hart zu kennen
status: backlog
ready: true
creator: Alexander Sacharov
assignee: alex
goal: "Das von 'jaira hook example' gedruckte Skript entscheidet anhand der Board-Lanes, welcher Lane-Wechsel einen Ton verdient, damit es auch auf einem Board mit eigenen Lane-Namen klingelt."
context: |-
  Das Beispielskript aus core/hook/example/notify.sh nennt die Lanes 'human', 'signoff' und 'done' im Klartext.
  Auf einem Board, dessen Lanes anders heissen, ist es dadurch stumm - und zwar ohne Fehler, also merkt es niemand.
  Die Information gibt es im Binary bereits: core/lane kennt pro Lane, ob sie einem Menschen gehoert (requires-human-exit) und ob sie terminal ist.
  Sie kommt nur nicht heraus: 'jaira lanes --json' gibt requires-question und requires-human-exit nicht aus (geprueft am 2026-09-13). Ein Skript kann die Rolle einer Lane also nicht erfragen.
  Zwei Wege sind offen: entweder 'jaira lanes --json' um die Rollen-Felder erweitern und das Skript sie lesen lassen, oder jaira uebergibt dem Hook gleich eine Variable wie JAIRA_STATUS_ROLE. Der zweite Weg aendert den Hook-Vertrag in core/hook/hook.go und war in Ticket Y0A7DT ausdruecklich ausgeschlossen.
  Nicht Teil davon: an der Ton-Regel selbst etwas aendern. Ein Ton gehoert weiter nur dem Zustand, in dem sich ohne den Menschen nichts bewegt.
definition-of-done: "Das gedruckte Beispielskript klingelt auf einem Board mit umbenannten Lanes fuer die Lane, die einem Menschen gehoert, und bleibt fuer die Agenten-Lanes stumm; ein Test belegt das mit umbenannten Lanes; go test ./... -race gruen"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-13T21:55:39Z
updated-at: 2026-09-13T21:55:39Z
---

# Das Hook-Beispiel liest die Lane-Rollen aus dem Board statt sie hart zu kennen

## Definition of Done

- [ ] Das gedruckte Beispielskript klingelt auf einem Board mit umbenannten Lanes fuer die Lane, die einem Menschen gehoert, und bleibt fuer die Agenten-Lanes stumm; ein Test belegt das mit umbenannten Lanes; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

