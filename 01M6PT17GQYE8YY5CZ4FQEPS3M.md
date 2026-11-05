---
id: 01M6PT17GQYE8YY5CZ4FQEPS3M
title: "jaira-cockpit: session.append wird seit Claude Code 2.1.318 nicht mehr erreicht"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: Cockpit kuerzt und verlinkt Peer-Meldungen wieder in Claude Code 2.1.318 und spaeter
context: |-
  Was kaputt ist: seit dem Update auf Claude Code 2.1.318 (05.11.2026) kuerzt jaira-cockpit keine Meldungen mehr, und die Zeilen sind nirgends mehr anklickbar.
  Ausloeser: Claude Code meldet 'jaira-cockpit: its session.append hook was skipped. "cc-plugin-sec-default" runs first and sent the event on past it.'
  cc-plugin-sec-default ist ein eingebautes Modul von Claude Code, nicht von uns; es laeuft vor uns auf session.append und reicht das Event nicht an jaira-cockpit weiter.
  Betroffen: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1179 (on('session.append', ... rewriteRow)) und hooks/filter.ts:10.
  Bekannt: mcp__jaira-cockpit__report funktioniert noch - die Berichte selbst gehen nicht verloren, nur Kuerzen und Anklicken fehlen.
  Offen: ob es eine Reihenfolge-/Prioritaets-Option fuer Hooks gibt, oder ob ein anderes Event (z. B. vor dem Anhaengen) dieselbe Umschreibung erlaubt. Ein Bericht an Anthropic ist parallel abgeschickt (/feedback).
  Live-Test: claude --plugin-dir <worktree>/core/mod/builtin/jaira-cockpit
definition-of-done: Eine Peer-Meldung erscheint in Claude Code 2.1.318 wieder gekuerzt mit anklickbarer Zeile
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-11-05T21:32:27Z
updated-at: 2026-11-05T21:33:00Z
updated-by: Alexander Sacharov
---

# jaira-cockpit: session.append wird seit Claude Code 2.1.318 nicht mehr erreicht

## Definition of Done

- [ ] Eine Peer-Meldung erscheint in Claude Code 2.1.318 wieder gekuerzt mit anklickbarer Zeile
- [ ] Kein Hinweis 'session.append hook was skipped' mehr fuer jaira-cockpit

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

