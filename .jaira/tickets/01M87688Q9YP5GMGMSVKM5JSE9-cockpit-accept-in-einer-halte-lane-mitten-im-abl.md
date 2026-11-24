---
id: 01M87688Q9YP5GMGMSVKM5JSE9
title: "Cockpit: Accept in einer Halte-Lane mitten im Ablauf schiebt das Ticket in die nächste Lane"
status: in-progress
ready: true
creator: Sascha
assignee: Sascha
goal: "Eine Person kann ein Ticket in jeder Personen-Lane mit Accept freigeben: in einer Lane vor dem Ende geht es in next_lane weiter, nicht nur von signoff nach done."
context: |-
  Was falsch ist: Accept im Cockpit-Pane 'Wartet auf dich' funktioniert nur, wenn die nächste Lane terminal ist (done).
  In jeder anderen Personen-Lane kommt nur: '<ID> waits in <lane>, which does not end the work: nothing to accept there'.
  Auslöser 24.11.2026: Board 'video' im Repo demo-video, Ticket in Lane script-approval (Halt Skript), next_lane = build. Sascha wollte das Skript freigeben und konnte es nicht; nur 'jaira move --to build' von Hand ging.
  Betrifft auch die Lane 'human' auf dem default-Board, wenn die Frage beantwortet ist.
  Code: core/mod/builtin/jaira-cockpit/hooks/board.ts:158 (acceptTo nur bei next.terminal), :1182 (step notTerminal), hooks/register.tsx:1502, hooks/texts.ts:80. Test, der das heutige Verhalten festschreibt: tests/board.test.ts:854.
  Offen zu klären im Plan: ob jaira move aus einer Personen-Lane ohne --force geht, wenn eine Person es auslöst (Cockpit-Button = Person).
  Geplant danach: die Entscheidungsseite (Skill decision-page, Markierungen 'accept' zweier Personen) soll denselben Weg nutzen, um Tickets aus Personen-Lanes weiterzuschieben.
definition-of-done: In einer Personen-Lane mit nicht-terminaler next_lane bietet das Cockpit Accept an; Accept verschiebt das Ticket nach next_lane
tags:
  - cockpit
  - lanes
blocked-by: []
related: []
commits: []
created-at: 2026-11-24T16:29:33Z
updated-at: 2026-11-24T16:39:03Z
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-51795
claimed-at: 2026-11-24T16:30:00Z
outcome-what: "Cockpit Accept moves a ticket from any person's lane one lane on (script-approval -> build, human -> review), not only signoff -> done; the Execute line names the target mid-board"
outcome-why: "Sascha could not release a script in the video board's script-approval lane from the cockpit; only a hand-typed jaira move went"
outcome-resolves: M5JSE9
review-summary: |-
  core/mod/builtin/jaira-cockpit/hooks/board.ts executionOf: an open question in a lane that does not end the work (human) is now accepted past - jaira move keeps the question (core/gate/gate.go:325 only checks entering), so review gets a ticket whose question nobody answered; the ticket's context says human only 'wenn die Frage beantwortet ist'. Add a step needAnswer (question open, lane does not end) that refuses with 'answer it first', words typed still give back
  core/mod/builtin/jaira-cockpit/hooks/texts.ts summary: the default ends = true has one caller that always passes it; make it required
round: 2
back-from: critique
back-by: lane
back-count: 1
---

# Cockpit: Accept in einer Halte-Lane mitten im Ablauf schiebt das Ticket in die nächste Lane

## Definition of Done

- [x] In einer Personen-Lane mit nicht-terminaler next_lane bietet das Cockpit Accept an; Accept verschiebt das Ticket nach next_lane
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts:162; test 'Execute moves it one lane on and says where'
- [x] Der Text des Buttons/der Zusammenfassung nennt das Ziel (z. B. 'Accept → build'), damit klar ist, dass das Ticket nicht nach done geht
  proof: core/mod/builtin/jaira-cockpit/hooks/texts.ts:73 ('Execute: accept → review' asserted in that test)
- [x] signoff → done verhält sich unverändert; Give back (note fix / getippter Grund) unverändert
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:731,760 gate PR/hold on acceptEnds; signoff tests ('Execute: done', move --to done --force) and give-back tests unchanged and green
- [x] Tests: board.test.ts deckt beide Fälle ab (terminal und nicht-terminal); der alte notTerminal-Test ist angepasst
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'Execute moves it one lane on', 'names no lane after it', executionOf table rows 'in a lane that does not end' / 'with no next lane'; claude plugin test 525 pass
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/M5JSE9.md
- [x] Eine offene Frage in einer Lane, die die Arbeit nicht beendet (human), hält Accept an: Execute verschiebt nichts und verlangt die Antwort; getippter Text gibt weiterhin zurück
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts:1191 waitsForAnswer + step needAnswer; test 'a question still open in a lane that does not end the work: Execute moves nothing and asks for the answer'

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-24 16:30 · Sascha** — parent session: uds:/run/user/1000/cc-socks/3417.sock
- **2026-11-24 16:30 · Sascha** — no open decisions: --force question settled from code - accepting() (register.tsx) dry-runs jaira move, otherRefusal drops only needs_human, then --force; the same path serves a mid-flow person lane. Plan: Waiting gets acceptEnds; acceptTo = next_lane whenever it exists; afterAccept (PR/land) and heldInSignoff only when it ends; summary 'Execute: accept -> <lane>' mid-flow, 'Execute: done' unchanged
- **2026-11-24 16:36 · Sascha** — critique: (1) an open question in a mid-board lane must stop Accept - jaira does not clear the question on leaving human, so review would get it unanswered; new step needAnswer. (2) summary's ends default is speculative, make it required
