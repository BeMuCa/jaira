---
id: 01M5MNWVKBC8Q193EK7YX6HB6V
title: "Die Einzeiler des Cockpits nennen die Lane, in der das Ticket wirklich steht"
status: critique
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
updated-at: 2026-10-26T20:01:38Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: x6hb6v-driver
claimed-at: 2026-10-26T19:57:50Z
outcome-what: "Der Cockpit-Einzeiler zeigt 'vorher → jetzt' nur noch aus dem Board: shorten() in core/mod/builtin/jaira-cockpit/hooks/register.tsx liest die Lane des einen Tickets per `jaira show <id> --json`, erst nachdem digest() ein Ticket und keinen Fehlschlag gefunden hat (digest ist jetzt async). Die Lane, die der Berichtstext nennt, wird nie gelesen; laneChangeOf liegt in digest.ts."
outcome-why: "Der Berichtstext nannte beliebige Lanes ('list -> human', 'done -> in-progress', 'moved to review'), die das Ticket nie hatte; das Board ist die einzige Quelle, die stimmt. lastSeen wird nur beim ersten Auftauchen und bei Berichten gesetzt, weil refreshBoard jede Minute liest und ein minutenlang arbeitender Subagent sonst nie einen Wechsel zeigte."
outcome-resolves: "Der Einzeiler stimmt mit dem Board ueberein: der echte Wechsel oder keiner. Vier Tests in core/mod/builtin/jaira-cockpit/tests/cockpit.test.ts ('the lane in the line comes from the board') decken die drei beobachteten Faelle und den positiven Fall ab; Release-Notiz core/release/notes/X6HB6V.md."
review-summary: none
round: 2
---

# Die Einzeiler des Cockpits nennen die Lane, in der das Ticket wirklich steht

## Definition of Done

- [x] Der Einzeiler leitet den Lanewechsel aus dem Board ab, nicht aus dem Berichtstext; nennt der Text eine Lane, die das Ticket nicht hat, erscheint kein Wechsel
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx laneOf()/shorten(), digest.ts async digest(text, laneChange)
- [x] Ein Test deckt die drei beobachteten Fehlfaelle ab
  proof: tests/cockpit.test.ts 'the lane in the line comes from the board' (4 tests, fake board answers jaira show)
- [x] Release-Notiz unter core/release/notes/
  proof: core/release/notes/X6HB6V.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-26 19:48 · BeMuCa** — driver: pre-process is not ticked for this ticket (Options), so the design is written here. Why not 'previous poll vs now': refreshBoard reads jaira list --json every 60 s and after every turn; a subagent works for minutes, so the poll sees its moves long before its report arrives - 'last poll vs now' would show no change almost always. Design instead: register.tsx keeps one map handle -> lane 'as of the last report about this ticket'. refreshBoard only ADDS tickets it has not seen yet (first sight), never overwrites. On a report, shorten() reads jaira list --json fresh (default board plus the other boards refreshBoard reads), emits '<before> -> <now>' when both are known and differ, then stores <now>. The report text is no longer read for lanes at all (laneOf/ARROW/INTO/LANES go), so a lane only the text names can never show up. Verdict stays text-derived: not in the DoD, not one of the observed failures.
- **2026-10-26 19:51 · BeMuCa** — lanesNow() liest bei jedem Bericht frisch; in shorten wird board nicht gecacht (sonst entfiele checkVersion in refreshBoard). Bei Lesefehler einer Nebenboard-Liste: gar kein Lanewechsel.
- **2026-10-26 19:57 · BeMuCa** — critique (first pass): five findings, all in review-summary. (1+2, one fix) read the lane only once digest has found a ticket and no failure — make digest's laneChange callback async — and read it with one 'jaira show <ticket> --json' (status) instead of lanesNow: show resolves the id on every board (core/ticket/store.go:1174; checked on a scratch board: an 'ops' ticket is missing from 'jaira list --json' but 'jaira show <id> --json' gives status and board), so the copy of refreshBoard's board loop and its separate failure rule go, and failure reports (which shorten sees twice, session.receive then prompt.submit) spawn no jaira at all. (3) laneChangeOf to digest.ts — it formats the line, board.ts parses jaira. (4) release note as an instruction, and not 'as of the last report' (the first report compares with first sight). (5) cockpit.test.ts first(): drop the no-op lanes/to parameters. Let stand: lastSeen seeded at first sight and updated only on reports (driver note 19:48) — found no case where it shows a lane the ticket did not take; no board cache in shorten (note 19:51); verdict stays text-derived (19:48); the edits to the existing tests drop only the text-derived lane, every verdict/question/failure/usage/marker assertion is unchanged.
- **2026-10-26 20:01 · BeMuCa** — critique (second pass): none left - checked f04f549e against the five first-pass findings, all addressed: (1) digest is async and awaits laneChange only after isFailure/ticketOf (digest.ts:85-92); (2) lanesNow and its board loop are gone, laneOf reads one 'jaira show <id> --json' status (register.tsx:218-223; checked: show --json carries top-level status); (3) laneChangeOf now in digest.ts:76; (4) release note is one instruction line; (5) first($, ticket) without the no-op parameters. Repair introduced nothing that breaks the DoD; first-pass let-stands untouched.
