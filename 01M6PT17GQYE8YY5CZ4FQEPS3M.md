---
id: 01M6PT17GQYE8YY5CZ4FQEPS3M
title: "jaira-cockpit: session.append wird seit Claude Code 2.1.318 nicht mehr erreicht"
status: testing
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
updated-at: 2026-11-05T21:47:38Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-54459
claimed-at: 2026-11-05T21:33:03Z
outcome-what: "session.append hook of jaira-cockpit now matches only door prompt|delivery; release note; test for the delivery and attachment doors"
outcome-why: "the built-in cc-plugin-sec-default in 2.1.318 passes person plugins over on attachment/hook-context/compaction rows, and a hook without a matcher reported itself as skipped there"
outcome-resolves: "the 'session.append hook was skipped' notice; dead clicks and scroll trace to the fullscreen renderer Claude Code turned off on this machine, not to the mod"
review-summary: "none"
round: 3
review-gaps: "none - the diff is one matcher, one comment block, one removed check, one test and the release note; left alone the tsc error at register.tsx:381 that was there before this change (null is not void)"
test-verdict: "fail: TestEveryInstalledFileIsAShippedVersion - shipped.sha256 lacks the new versions of filter.ts and register.tsx"
---

# jaira-cockpit: session.append wird seit Claude Code 2.1.318 nicht mehr erreicht

## Definition of Done

- [~] Eine Peer-Meldung erscheint in Claude Code 2.1.318 wieder gekuerzt mit anklickbarer Zeile
- [x] Kein Hinweis 'session.append hook was skipped' mehr fuer jaira-cockpit
  proof: register.tsx:1184 matcher door prompt|delivery; live claude -p 2.1.318 --plugin-dir: master mod logs 'session.append bypassed by cc-plugin-sec-default' + 'hook was skipped', this branch logs neither
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
- **2026-11-05 21:38 · Alexander Sacharov** — Ursache DoD 3 (Klicks tot, Scroll springt): Claude Code hat am 2026-11-04 21:41 den Fullscreen-Renderer auf diesem Rechner abgeschaltet, nach 3 fehlgeschlagenen Starts: ~/.claude.json fullscreenAutoDisabled {version 2.1.318, strikes 3}. Der klassische Renderer hat keine Mausklicks. Kein Fehler der Mod: claude plugin test . laeuft auf 2.1.318 mit 175/175 gruen, inkl. Press-Tests. Abhilfe liegt bei Alex: /tui fullscreen oder CLAUDE_CODE_NO_FLICKER=1 (laut Meldung im Binary). Gilt nur fuer 2.1.318; die naechste Version setzt es zurueck.
- **2026-11-05 21:41 · Alexander Sacharov** — Code fertig: Matcher als Literal { door: ['prompt','delivery'] } in register.tsx - mit einer Konstante zeigt claude plugin validate door=?, die Engine liest den Matcher aus dem Quelltext. Test 'a shortened message is rewritten on the delivery door too, and an attachment row is left alone' (176/176 gruen). DoD 1 braucht noch den Live-Test in 2.1.318 (claude --plugin-dir ...): ob eine Peer-Meldung, die waehrend eines laufenden Zuges ankommt, ueber door attachment kommt - die kann seit 2.1.318 KEIN Personen-Plugin mehr umschreiben (sec-default-Politik), session.receive kuerzt sie dann trotzdem, wenn es feuert. tsc meldet einen alten Fehler register.tsx:381 (null nicht void), nicht aus dieser Aenderung.
- **2026-11-05 21:41 · Alexander Sacharov** — critique: DELIVERY_DOORS in filter.ts duplicates the matcher literal in register.tsx (the only caller of rewriteRow). Fix: drop the constant and the door check in rewriteRow, the matcher is the only door gate.
- **2026-11-05 21:42 · Alexander Sacharov** — critique round 2: nothing left - the doors now live only in the matcher literal (register.tsx:1184), rewriteRow no longer checks them.
- **2026-11-05 21:44 · Alexander Sacharov** — testing: go test ./... rc=1, core/role TestEveryInstalledFileIsAShippedVersion: core/shipped/shipped.sha256 lacks 85cf6274... filter.ts and 364f06f7... register.tsx. Fix: append the two lines the test prints (append rather than rerunning the script, so the file only grows).
- **2026-11-05 21:46 · Alexander Sacharov** — critique round 3 (after testing): the only change is two appended lines in core/shipped/shipped.sha256, the way the test asks for it - nothing to say. optimize: none.
