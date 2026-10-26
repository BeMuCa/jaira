---
id: 01M5MNWVKBC8Q193EK7YX6HB6V
title: "Die Einzeiler des Cockpits nennen die Lane, in der das Ticket wirklich steht"
status: in-progress
ready: true
creator: BeMuCa
goal: "Die Kurzfassung eines Agentenberichts im Cockpit stimmt mit dem Board ueberein: sie nennt den tatsaechlichen Lanewechsel, oder keinen."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkt 10): die Einzeiler waren zweimal falsch ('GH0GKD · list → human', 'done → in-progress') und nannten einmal eine Lane, die das Ticket nie betreten hat. Der volle Bericht musste jedes Mal gelesen werden - dann ist die Kurzfassung wertlos.
  Ursache vermutlich: core/mod/builtin/jaira-cockpit/hooks/digest.ts liest den Lanewechsel aus dem Text des Berichts (Regex auf 'a -> b'), nicht aus dem Board.
  Vorschlag: den Wechsel aus jaira list --json vor und nach dem Bericht ableiten (der Cockpit liest das Board ohnehin, board.ts:3) und den Text nur fuer die Frage nutzen.
definition-of-done: "Der Einzeiler leitet den Lanewechsel aus dem Board ab, nicht aus dem Berichtstext; nennt der Text eine Lane, die das Ticket nicht hat, erscheint kein Wechsel"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:25:58Z
updated-at: 2026-10-26T19:51:49Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: x6hb6v-driver
claimed-at: 2026-10-26T19:50:01Z
---

# Die Einzeiler des Cockpits nennen die Lane, in der das Ticket wirklich steht

## Definition of Done

- [x] Der Einzeiler leitet den Lanewechsel aus dem Board ab, nicht aus dem Berichtstext; nennt der Text eine Lane, die das Ticket nicht hat, erscheint kein Wechsel
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx lanesNow()/shorten(), digest.ts digest(text, laneChange)
- [ ] Ein Test deckt die drei beobachteten Fehlfaelle ab
- [ ] Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-26 19:48 · BeMuCa** — driver: pre-process is not ticked for this ticket (Options), so the design is written here. Why not 'previous poll vs now': refreshBoard reads jaira list --json every 60 s and after every turn; a subagent works for minutes, so the poll sees its moves long before its report arrives - 'last poll vs now' would show no change almost always. Design instead: register.tsx keeps one map handle -> lane 'as of the last report about this ticket'. refreshBoard only ADDS tickets it has not seen yet (first sight), never overwrites. On a report, shorten() reads jaira list --json fresh (default board plus the other boards refreshBoard reads), emits '<before> -> <now>' when both are known and differ, then stores <now>. The report text is no longer read for lanes at all (laneOf/ARROW/INTO/LANES go), so a lane only the text names can never show up. Verdict stays text-derived: not in the DoD, not one of the observed failures.
