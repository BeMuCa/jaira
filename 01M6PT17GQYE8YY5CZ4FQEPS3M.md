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
updated-at: 2026-11-05T21:38:41Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-54459
claimed-at: 2026-11-05T21:33:03Z
---

# jaira-cockpit: session.append wird seit Claude Code 2.1.318 nicht mehr erreicht

## Definition of Done

- [ ] Eine Peer-Meldung erscheint in Claude Code 2.1.318 wieder gekuerzt mit anklickbarer Zeile
- [ ] Kein Hinweis 'session.append hook was skipped' mehr fuer jaira-cockpit
- [ ] Die Knoepfe im Cockpit-Panel ([ Overview ], [ N waiting for ... ]) reagieren in 2.1.318 wieder auf Klick, und das Panel scrollt nicht von selbst

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-05 21:33 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/47752.sock
- **2026-11-05 21:34 · Alexander Sacharov** — Alex (05.11.2026, Screenshot): das Cockpit-Panel wird gezeichnet (1 worker | 0 review | ..., [ Overview ], [ 48 waiting for BeMuCa, berkc ]), aber kein Knopf ist anklickbar - nirgends; ausserdem springt die Scroll-Position des Menues beliebig. Also nicht nur session.append: Maus/Klick und Scroll im Panel sind in 2.1.318 ebenfalls kaputt. Auch /feedback zeigte in dieser Sitzung '(no content)'.
- **2026-11-05 21:38 · Alexander Sacharov** — Ursache session.append-Hinweis: cc-plugin-sec-default (eingebaut in 2.1.318, sitzt bei Team/Enterprise ganz aussen, Tier prepend) ruft next.to an den user-Plugins vorbei fuer session.append{door=attachment|hook-context|compaction} (PASSED_OVER_DOORS im Binary). Absicht: Prompt-Inhalt ausser Reichweite persoenlicher Plugins. Unser Hook (register.tsx:1179) hat keinen Matcher, haengt also auch an diesen Tueren -> Hinweis. rewriteRow (filter.ts:80) wirkt aber nur auf door prompt|delivery, die sec-default nicht ueberspringt. Fix: Matcher { door: ['prompt','delivery'] } am Hook. Kein Prioritaets-Hebel fuer user-Plugins: next.to ueberspringt nie eine Tier mit mehr Autoritaet, user kann nicht vor prepend.
