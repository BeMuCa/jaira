---
id: 01M851Y62KHD9DVD2GDAM31BYV
title: "Give back tippt die Entscheidung in einen Worker-Tab, der an einem Freigabe-Dialog haengt"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Give back im Cockpit schreibt nie in einen Worker-Tab, in dem gerade ein Freigabe-Dialog offen ist (Herdr agent_status=blocked); die Person erfaehrt stattdessen, dass der Worker auf ihre Freigabe wartet."
context: |-
  Was falsch ist: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1798 (Give back, alter Worker wird wiederverwendet) macht 'herdr pane send-text' und danach 'send-keys enter' in die Pane des alten Workers.
  Haengt der Worker dort an einem Freigabe-Dialog (z.B. Bash-Befehl), landet der Text im Dialog und das Enter beantwortet ihn - im Namen der Person, ohne dass sie es sieht.
  Herdr meldet diesen Zustand als agent_status=blocked. spawn.sh prueft genau das schon vor seinem send-keys und bricht ab (spawn.sh: 'claude blocked' -> exit 1).
  register.tsx:1790 kennt nur working ('nie ersetzen'), nicht blocked.
  Gefunden vom Bughunt zu RJFA19 am 23.11.2026 (pre-existing, nicht durch RJFA19 entstanden). RJFA19 zeigt den Zustand jetzt im Worker-Tab an.
  Offen fuer den Plan: bei blocked nichts tippen und die Person auf den Dialog hinweisen, oder einen neuen Worker starten - der alte haengt dann aber weiter.
definition-of-done: Give back tippt nichts und drueckt kein Enter in einen Tab mit agent_status=blocked; Test in core/mod/builtin/jaira-cockpit/tests
tags:
  - cockpit
blocked-by: []
related: []
commits: []
created-at: 2026-11-23T20:35:39Z
updated-at: 2026-11-23T20:41:03Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-41227
claimed-at: 2026-11-23T20:36:42Z
outcome-what: "Give back stops before the reuse/replace choice when the carrier tab's agent_status is blocked: nothing typed, no Enter, no tab closed, no new worker; ticket note, teamlead line and toast name the waiting tab"
outcome-why: "send-text + Enter into a pane at an approval dialog answered the dialog for the person, unseen"
outcome-resolves: "M31BYV: decision no longer lands in an approval dialog"
---

# Give back tippt die Entscheidung in einen Worker-Tab, der an einem Freigabe-Dialog haengt

## Definition of Done

- [x] Give back tippt nichts und drueckt kein Enter in einen Tab mit agent_status=blocked; Test in core/mod/builtin/jaira-cockpit/tests
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1794; test 'give back > a worker at an approval dialog gets nothing typed and no Enter...' tests/board.test.ts:2413
- [x] Die Person sieht in der Meldung von Give back, dass der Worker an einem Freigabe-Dialog wartet und in welchem Tab
  proof: hooks/texts.ts workerBlocked (toast names the tab), asserted in tests/board.test.ts:2413
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/M31BYV.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 20:36 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/48662.sock
- **2026-11-23 20:36 · Alexander Sacharov** — decision: blocked tab = like working: never replaced, nothing typed; handedTo none with the tab label, no second worker (the old one still hangs at its dialog). Settled by DoD 2, no question open.
