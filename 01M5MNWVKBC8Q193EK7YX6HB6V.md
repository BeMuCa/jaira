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
updated-at: 2026-10-26T19:59:12Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: x6hb6v-driver
claimed-at: 2026-10-26T19:57:50Z
outcome-what: "Der Cockpit-Einzeiler (core/mod/builtin/jaira-cockpit/hooks/digest.ts) liest keine Lane mehr aus dem Berichtstext; register.tsx liest beim Bericht `jaira list --json` aller Boards frisch (lanesNow) und zeigt 'vorher → jetzt' aus der Map lastSeen, nur wenn beide bekannt sind und sich unterscheiden."
outcome-why: "Der Text nannte beilaeufig beliebige Lanes ('list -> human', 'done -> in-progress', 'moved to review'), sodass der Einzeiler falsche oder nie betretene Lanes zeigte; refreshBoard ergaenzt lastSeen nur beim ersten Auftauchen, damit der Vergleich nicht schon die neue Lane enthaelt, wenn ein minutenlang arbeitender Subagent berichtet."
outcome-resolves: "Der Einzeiler stimmt mit dem Board ueberein: der echte Wechsel oder keiner. Vier neue Tests in core/mod/builtin/jaira-cockpit/tests/cockpit.test.ts decken die drei beobachteten Faelle und den positiven Fall ab; Release-Notiz core/release/notes/X6HB6V.md."
review-summary: |-
  core/mod/builtin/jaira-cockpit/hooks/register.tsx:199 reads the board (lanesNow) before digest knows there is a ticket: every failure or ticket-less report spawns jaira list + jaira boards + one list per side board, and a failure session.receive passes on whole reaches prompt.submit and does it all again. Instead make digest's laneChange parameter async (digest.ts:79-86, (ticket) => Promise<string | undefined>, awaited after the isFailure/ticketOf checks) so the board is read only for a report that will be shortened.
  core/mod/builtin/jaira-cockpit/hooks/register.tsx:213-230 lanesNow repeats refreshBoard's list+boards+per-board loop (register.tsx:313-337) with a different failure rule (undefined vs continue). Replace it with one `jaira show <ticket> --json` and its `status`, as noteOn already does (register.tsx:388): Store.Load resolves an id on every board (core/ticket/store.go:1174); checked on a scratch board, a ticket on board "ops" is absent from `jaira list --json` while `jaira show <id> --json` returns status=backlog, board=ops. lanesNow, its board loop and the side-board failure rule then go.
  core/mod/builtin/jaira-cockpit/hooks/board.ts:146-150 laneChangeOf formats a piece of the one-liner (`before → now`) and has one caller; board.ts reads jaira's JSON, the line format lives in digest.ts (the removed laneOf built the same `a → b` there). Move it into digest.ts beside digest().
  core/release/notes/X6HB6V.md:1 records what changed ("The teamlead chat now names ...") instead of telling the reader what to do, as CLAUDE.md asks and the neighbours do (7V1XPB "Start ...", A0C081 "Read why ...", TPCKZB "Tell lanes apart ..."); "as of the last report" is also wrong for the first report, which compares with the lane the cockpit first saw on the board. Rewrite as one instruction, e.g. "- Trust the lane change in an agent's one-line summary in the cockpit: it is what the board shows (`before → now` since the cockpit last saw that ticket), never a lane the report merely names; a ticket that did not move shows no lane, and `/cockpit report <id>` still gives the full text."
  core/mod/builtin/jaira-cockpit/tests/cockpit.test.ts:304-306 first() takes lanes and to only to write lanes[ticket] = to, and all four callers (313, 322, 331, 342) pass the lane the ticket already stands in, so it is a no-op. Drop both parameters and the assignment: first($, ticket).
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
