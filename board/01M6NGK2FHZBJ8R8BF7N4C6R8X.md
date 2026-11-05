---
id: 01M6NGK2FHZBJ8R8BF7N4C6R8X
title: Ein Hinweis auf eine neue jaira-Version geht im Chat nicht unter
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Meldet jaira oder der Agent im Chat eine neue Version, sieht die Person das als farbige Info-Zeile, nicht als einen Satz im Fliesstext."
context: |-
  Wunsch berkc 05.11.2026: 'wenn eine neue jaira version existiert und der chat es uns sagen will, das ganze in einer Infofarbe schreiben, damit es nicht untergeht'.
  Heute: jaira next/resume/claim/move schreiben den Hinweis auf stderr; der Agent gibt ihn als normalen Text weiter (Agent-Block in CLAUDE.md: 'when jaira reports a newer release ...').
  Das Cockpit hat schon eine Band-Zeile 'jaira X is released (this is Y)' mit Update-Knopf, Farbe warning (core/mod/builtin/jaira-cockpit/hooks/register.tsx:1284, texts.ts:33). berkc hat sie trotzdem nicht bemerkt.
  Ein Agent kann im Chat keine Farbe setzen; Farbe gibt es nur ueber das Cockpit (Band, Toast, Statuszeile).
  Zu klaeren: ob die Band-Zeile auffaelliger wird (info-Farbe, eigene Zeile oben) oder ein Toast beim ersten Erkennen reicht.
definition-of-done: "Eine neue Version erscheint im Cockpit in info-Farbe an einer Stelle, die beim naechsten Prompt sichtbar ist (Band oben oder Toast); Test in den Cockpit-Tests"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-11-05T09:28:11Z
updated-at: 2026-11-05T09:28:11Z
---

# Ein Hinweis auf eine neue jaira-Version geht im Chat nicht unter

## Definition of Done

- [ ] Eine neue Version erscheint im Cockpit in info-Farbe an einer Stelle, die beim naechsten Prompt sichtbar ist (Band oben oder Toast); Test in den Cockpit-Tests
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

