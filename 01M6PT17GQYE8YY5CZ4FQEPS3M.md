---
id: 01M6PT17GQYE8YY5CZ4FQEPS3M
title: "jaira-cockpit: session.append wird seit Claude Code 2.1.318 nicht mehr erreicht"
status: review
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
updated-at: 2026-11-05T22:00:53Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-54459
claimed-at: 2026-11-05T21:33:03Z
outcome-what: "session.append hook of jaira-cockpit now matches only door prompt|delivery; release note; test for the delivery and attachment doors"
outcome-why: "the built-in cc-plugin-sec-default in 2.1.318 passes person plugins over on attachment/hook-context/compaction rows, and a hook without a matcher reported itself as skipped there"
outcome-resolves: "the 'session.append hook was skipped' notice; dead clicks and scroll trace to the fullscreen renderer Claude Code turned off on this machine, not to the mod"
review-summary: "none"
round: 3
review-gaps: "none - the diff is one matcher, one comment block, one removed check, one test and the release note; left alone the tsc error at register.tsx:381 that was there before this change (null is not void)"
test-verdict: "pass for what a machine can check: go test ./... rc=0, claude plugin test 176/176, DoD 2 verified live on 2.1.318 (notice with the master mod, none with the branch). DoD 1: rewriting on door prompt|delivery is proven by a test, the clickable line needs a person in the fullscreen renderer. DoD 3: no code change possible - fullscreen is switched off on this machine (fullscreenAutoDisabled)"
question: ""
---

# jaira-cockpit: session.append wird seit Claude Code 2.1.318 nicht mehr erreicht

## Definition of Done

- [x] Eine Peer-Meldung erscheint in Claude Code 2.1.318 wieder gekuerzt mit anklickbarer Zeile
  proof: Alex live on 2.1.318 after /tui fullscreen: peer messages arrive as '[shortened by jaira-cockpit] ... - full report: /cockpit report <ID>' and are clickable; path session.receive/prompt.submit (register.tsx:1146/1155), test 'a shortened message is rewritten on the delivery door too, and an attachment row is left alone'
- [x] Kein Hinweis 'session.append hook was skipped' mehr fuer jaira-cockpit
  proof: register.tsx:1184 matcher door prompt|delivery; live claude -p 2.1.318 --plugin-dir: master mod logs 'session.append bypassed by cc-plugin-sec-default' + 'hook was skipped', this branch logs neither
- [x] Die Knoepfe im Cockpit-Panel ([ Overview ], [ N waiting for ... ]) reagieren in 2.1.318 wieder auf Klick, und das Panel scrollt nicht von selbst
  proof: Alex live on 2.1.318: clicks and scroll work after /tui fullscreen; cause ~/.claude.json fullscreenAutoDisabled {2.1.318, strikes 3}, no code change

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
- **2026-11-05 21:47 · Alexander Sacharov** — testing: live run claude -p --plugin-dir (2.1.318): with master's register.tsx/filter.ts the debug log shows 'jaira-cockpit: session.append bypassed by cc-plugin-sec-default' and the skipped notice; with the branch there is neither, and session.append settles in 2.2ms. DoD 1 and 3 cannot be proven without a person.
- **2026-11-05 21:47 · Alexander Sacharov** — question: DoD 3 (clicks/scroll) is not a mod bug: Claude Code 2.1.318 switched off fullscreen on this machine on 2026-11-04 after 3 failed starts (~/.claude.json fullscreenAutoDisabled). 1) (Recommended) Alex runs /tui fullscreen, then checks DoD 1+3 live in the cockpit (claude --plugin-dir <worktree>/core/mod/builtin/jaira-cockpit) - an agent ticks them on that word. 2) Remove DoD 3 from the ticket as not a code matter, check only DoD 1. 3) Something else.
- **2026-11-05 21:50 · Alexander Sacharov** — Alex (via teamlead, 05.11.2026): the cockpit works on Windows - fullscreen was auto-disabled only on this WSL machine, DoD 3 is no mod bug. Answer to the question: option 1 - Alex runs /tui fullscreen here and checks DoD 1 and DoD 3 live with claude --plugin-dir <worktree>/core/mod/builtin/jaira-cockpit. The worker waits for his word, which the teamlead passes on; then tick DoD 1+3 and hand the ticket to a fresh review.
- **2026-11-05 21:50 · Alexander Sacharov** — Correction to the previous note: Claude Code on Windows is v2.1.317, not 2.1.318. That it works there does not prove the fullscreen diagnosis - it only shows that 2.1.317 works. What decides DoD 1 and DoD 3 is still Alex's live check on 2.1.318 after /tui fullscreen (option 1 stands).
- **2026-11-05 22:00 · Alexander Sacharov** — DoD 1 path: the '[shortened by jaira-cockpit] ... - full report: /cockpit report <ID>' lines the teamlead saw on master's mod under 2.1.318 come from session.receive (register.tsx:1146, relayed peer messages) and prompt.submit (register.tsx:1155) - neither is in sec-default's scan list (PASSED_OVER_DOORS only touches session.append on attachment|hook-context|compaction). The session.append rewrite on door prompt|delivery was never skipped either. So shortening never broke: what failed was the click, through fullscreenAutoDisabled, and the notice came from the hook having no matcher on the attachment doors. DoD 1 holds on master as well as on the branch.
