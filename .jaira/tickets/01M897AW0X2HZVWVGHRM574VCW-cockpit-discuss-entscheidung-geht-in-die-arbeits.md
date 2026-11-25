---
id: 01M897AW0X2HZVWVGHRM574VCW
title: "Cockpit Discuss: Entscheidung geht in die Arbeits-Lane des Boards, nicht fest in-progress"
status: critique
ready: true
creator: Sascha
goal: "Nach einer Entscheidung im Cockpit (Discuss/Explain) landet das Ticket in der Lane, die Give back seit PGXJG0 berechnet, auf jedem Board."
context: |-
  Der Discuss-Pfad im Cockpit scheitert auf Boards ohne Lane in-progress (z.B. ein Video-Board mit todo, script, script-approval, build, draft-approval, voice-approval, voice, final, acceptance, done).
  Ursache: core/mod/builtin/jaira-cockpit/hooks/texts.ts:245 sagt dem Agenten fest: jaira move <id> --to in-progress.
  PGXJG0 (24./25.11.2026) hat Give back schon auf backLaneOf (board.ts) umgestellt; der Discuss-Text blieb aussen vor, Review-Notiz von PGXJG0.
  Fix-Idee: denselben backLaneOf-Wert in den Text einsetzen.
definition-of-done: Discuss-Text nennt die Lane aus backLaneOf statt in-progress; Test in core/mod/builtin/jaira-cockpit/tests
tags:
  - cockpit
  - lanes
blocked-by: []
related: []
follows: 01M87RTTZFTY9RV9R3XPPGXJG0
commits: []
created-at: 2026-11-25T11:26:56Z
updated-at: 2026-11-25T11:42:01Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-905098
claimed-at: 2026-11-25T11:38:24Z
outcome-what: "Cockpit Discuss moves a decided ticket back to backLaneOf (w.backTo), in this chat and via spawn.sh --back <lane> in the Herdr tab"
outcome-why: boards without in-progress (video board) refused the hardcoded move; Give back already used backLaneOf since PGXJG0
outcome-resolves: Discuss decision on a video board lands in its working lane (build/script/final)
---

# Cockpit Discuss: Entscheidung geht in die Arbeits-Lane des Boards, nicht fest in-progress

## Definition of Done

- [x] Discuss-Text nennt die Lane aus backLaneOf statt in-progress; Test in core/mod/builtin/jaira-cockpit/tests
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx openChat --back w.backTo; test 'Discuss on a video board moves the work back to its build lane' (tests/board.test.ts); core/role/role_test.go TestSpawnDiscussReturnsItsDecision
- [x] Default-Board: Text nennt weiter in-progress
  proof: tests/board.test.ts 'an open question shows under its ticket with Discuss' (--back in-progress) and 'without Herdr, Discuss talks the question through in this chat'; spawn.sh default back=in-progress
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/574VCW.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-25 11:38 · Sascha** — parent session: uds:/run/user/1000/cc-socks/431929.sock
- **2026-11-25 11:39 · Sascha** — decision (worker): Herdr-Weg hat dieselbe feste Zeile in spawn.sh:644 (Kommentar dort: Aenderung an einem gehoert in den anderen). Darum neues spawn.sh-Flag --back <lane> nur fuer discuss, Default in-progress; openChat reicht w.backTo durch. Kein offener Punkt fuer die Person.
