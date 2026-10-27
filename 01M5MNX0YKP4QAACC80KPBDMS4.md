---
id: 01M5MNX0YKP4QAACC80KPBDMS4
title: Milestone-Dateien fuehren keine archivierten Tickets mehr als offen
status: critique
ready: true
creator: BeMuCa
goal: jaira milestone ls/show und die Kartenkante auf dem Board zaehlen ein archiviertes oder abgelegtes Ticket nicht mehr als offenen Teil des Milestones.
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkt 11): Milestone-Dateien im Haupt-Checkout listen archivierte Tickets weiter (P6Y1BG, XY4HDS, 9JGWGR) - der Milestone wirkt unfertig, und die Board-Dateien des Haupt-Checkouts hinken den Worktrees hinterher.
  Verwandt: GWK31K (archiviertes Ticket verschwindet sofort von allen Boards, signoff), N71NVG (fertiger Milestone ins Archiv), KWQCYZ (veraltete Milestone-Datei ueberschreibt beim add). Dieses Ticket ist der Rest: die Zaehlung und Anzeige, wenn die ID in der Datei bleibt.
definition-of-done: "jaira milestone ls zaehlt archivierte und abgelegte Tickets nicht als offen und zeigt sie in show als 'archiviert'/'abgelegt'; ein Test baut den Fall"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:26:04Z
updated-at: 2026-10-27T00:16:23Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: pbdms4-driver
claimed-at: 2026-10-27T00:16:23Z
round: 2
outcome-what: "jaira milestone ls prints 'N open of M ticket(s)' and --json gains 'open' (count stays the total); new 'jaira milestone show <name>' lists every member with where it is now (on the board · <lane>, archive · ..., logbook · done); the TUI M picker row counts open of total; logbook cards no longer carry the milestone colour bar. New link.Index.Open and Milestone.Open; ids stay in the milestone files."
outcome-why: "Milestone files keep the ids of archived and logbooked tickets, so ls and the picker counted them as open and a finished round looked unfinished."
outcome-resolves: "DoD 1: ls/show via Index.Open + Milestone.Open, TestLsCountsOnlyOpenTicketsAndKeepsTheTotal and TestShowSaysWhereEveryMemberIsNow. DoD 2: picker count and milestoneColors, TestMilestonePickerCountsOpenOfTotalAndLogbookCardsLoseTheBar. DoD 3 waits for the full race suite before the PR; release note core/release/notes/PBDMS4.md."
review-summary: "critique 1: ls prints 'N open of M ticket(s)', --json gains open, count stays the total. New milestone show lists each member with where it is now, built from existing helpers (milestoneKey, editRefusal, milestoneJSON, boardIndex, Whereabouts). Open = on the board and not in a terminal lane, as jaira tags and the logbook gate count it. Picker counts open of total; logbook cards lose the milestone bar. Right shape. LOOP: internal/tui/model.go:1791-1793 comment says 'the bar says what is still open', but done-lane cards keep their bar - reword. Cheap in the same round: docs/COMMANDS.md show row misses 'on its ref · jaira pull <h>' and 'nowhere here · dangling link'; .claude/skills/jaira/SKILL.md:213 ls comment lacks the open count."
---

# Milestone-Dateien fuehren keine archivierten Tickets mehr als offen

## Definition of Done

- [x] jaira milestone ls zaehlt archivierte und abgelegte Tickets nicht als offen und zeigt sie in show als 'archiviert'/'abgelegt'; ein Test baut den Fall
  proof: core/milestone/milestone.go Open + internal/cli/milestones.go ls/show; TestLsCountsOnlyOpenTicketsAndKeepsTheTotal, TestShowSaysWhereEveryMemberIsNow
- [x] Die Kartenkante und der M-Picker ignorieren sie
  proof: internal/tui/model.go milestoneColors + view.go picker count; TestMilestonePickerCountsOpenOfTotalAndLogbookCardsLoseTheBar
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/link/link.go: add (ix *Index) Open(id) bool beside Satisfied (~l.304): true only when place() is PlaceBoard or PlaceRef and not Done - 'open' as jaira tags/boards count it (tags.go:126, boards.go:47); reads no file
- [x] core/link/link_test.go: TestOpen... - backlog ticket open; done-lane ticket, s.Archive'd ticket, s.Logbook'd ticket and unknown id not open (fixtures store/make1/list, l.10-60); go test ./core/link -run TestOpen -count=1
- [x] core/milestone/milestone.go: add (m *Milestone) Open(isOpen func(id string) bool) int beside Members (l.109), counting members; file lines stay untouched (removing ids is N71NVG); small test in milestone_test.go
- [x] failing test first in internal/cli/milestones_test.go after TestLsNamesEveryMilestoneWithItsCount (l.245): 3 tickets via mkTicket, milestone create, runCLI archive h2, ticket.At(dir).Logbook(h3, "as-20260911") (pattern links_test.go:47); ls text has '1 open of' and '3 ticket(s)', ls --json open==1 and count==3
- [x] internal/cli/milestones.go ls (l.505-560): one boardIndex(s) (links.go:113); line '%3d ticket(s)' becomes '%3d open of %3d ticket(s)' (keeps milestones_test.go:239/262 green); --json adds "open" after milestoneJSON, "count" stays the total; update the Long help
- [x] NEEDS DRIVER OK (no 'milestone show' exists, milestones.go:80-87): newMilestoneShowCmd modelled on epic show (epics.go:259-309), registered in newMilestoneCmd: label, goal, 'N open of M ticket(s)', one row per member with handle, ix.Lookup(id).Whereabouts() (link.go:493: 'archive · unfinished/done', 'logbook · done', 'on the board · <lane>'), title; --json = milestoneJSON + open + tickets with place, lane, title; name lookup and refusals as the other milestone commands
- [x] extend the CLI test: show marks h2 'archive ·' and h3 'logbook · done', h1 'on the board · backlog'; show --json tickets carry place archive/logbook/board
- [x] TUI M picker: internal/tui/view.go:1646 count becomes open of total via ms.Open(m.linkIndex().Open) (model.go:1926); test in internal/tui/milestonebar_test.go with pickerRow (l.192), mkMilestone (l.13), m.store.Archive (pattern layer_test.go:751); go test ./internal/tui -run <that test> -count=1
- [x] TUI card edge: milestoneColors (model.go:1794) returns empty slots for m.isLogged(t) (logbook.go:37); do NOT filter m.mstones (model.go:431), the M filter matches logbook cards on purpose (model.go:567-569); test with fileIntoLogbook (logbook_test.go:16) that the logged card has no milestone colour
- [x] docs: docs/COMMANDS.md:142 ls row (open of total, json open) plus a show row; docs/AGENTS.md:164; README.md:862 command list
- [x] release note core/release/notes/PBDMS4.md: exactly one '- ' line, unwrapped, an instruction: ls reads open of total (json open), show lists where each ticket is now, M picker counts the same, logbook cards drop the milestone colour
- [x] per commit gofmt -l . empty, go vet ./... and go build ./... green; the full go test ./... -race runs once before the PR (driver ticks DoD 3)

## Progress
- **2026-10-27 00:04 · BeMuCa** — pre-process reasoning:
- There is no 'jaira milestone show' (milestones.go:80-87 registers create/add/rm/set/rename/ls). DoD 1 names it. Plan step 6 adds one, modelled on epic show. Driver confirms.
- Reproduced with this worktree's build on a scratch board: 3 members, 1 archived, 1 moved into a logbook folder. 'milestone ls' says '3 ticket(s)', 'jaira list' shows 1.
- 'open' = what jaira tags and boards call open (tags.go:126, boards.go:47): on the board (a file or a ref) and not in a terminal lane. So a done-lane ticket is not open either. This is the same answer the 'jaira logbook <milestone>' gate gives (logbook.go:393-404): open==0 means the milestone may be filed.
- Where 'on the board' comes from: Store.List via link.Index, not FiledAwayIDs alone. List already drops tickets filed here (refsync.go:1512-1521). Once GWK31K lands it also drops refs marked archived:, which FiledAwayIDs never sees.
- The main-checkout case in the context (archived in a worktree, its ref still up without a mark) stays counted as open until GWK31K is merged. Nothing on this branch can see it. GWK31K lives on feat/GWK31K-archive-on-ref.
- After GWK31K, such a ticket is PlaceUnknown here, so show prints 'nowhere here · dangling link', not 'archive'. Follow-up if that matters.
- Card edge: archived tickets never get a card. Only logbook cards in the done column carry the bar, so the bar goes there (milestoneColors). m.mstones stays unfiltered, because the M filter matches logbook cards on purpose (model.go:567-569).
- Left alone on purpose: roadmap L bubble and detail (layer.go:162 milestoneDone, via Satisfied; TestArchivedTicketCountsAsTheLinkIndexSays pins archived-unfinished as not done there); epic ls/show (epic.Count); the 'holds N ticket(s)' totals of create/add/rm and logbook.
- JSON: ls adds 'open'. 'count' keeps meaning the total, because changing what it means would break readers.
- **2026-10-27 00:05 · BeMuCa** — Driver decisions after pre-process: (1) step 6 approved - DoD 1 names 'milestone show', which does not exist; the smallest faithful reading is a new 'jaira milestone show <name>' modelled on 'epic show', using Whereabouts() wording (English, like the rest of the command). If the maintainer does not want a new command, drop steps 6-7 and reword DoD 1. (2) step 9 as planned - DoD 2 asks the card edge to ignore filed tickets; archived ones have no card, so only logbook cards lose the milestone bar; the M filter keeps matching logbook cards (model.go:567-569, 'narrowing the board narrows what it remembers too'). Alternative if that reads wrong at signoff: keep the bar, drop the isLogged condition in milestoneColors.
- **2026-10-27 00:09 · BeMuCa** — in-progress: show prints English Whereabouts() wording (archive · unfinished, logbook · done), not archiviert/abgelegt as DoD 1 words it. Name lookup reuses editRefusal (exit 5 for unknown). Picker row now reads 'N open of M'. Doc line for 'filed one exits 3' dropped: not verified.
- **2026-10-27 00:15 · BeMuCa** — critique 1: read the whole diff 37c3a045 and the notes. Standing on driver decisions 00:05 (new 'milestone show'; logbook cards lose the bar, M filter keeps matching) and in-progress 00:09 (English Whereabouts wording instead of archiviert/abgelegt) - not reopened.
LOOP (comment only, behaviour is right): internal/tui/model.go:1791-1793 says a logbook card loses the bar because 'that work is finished and the bar says what is still open'. False: a done-lane card still on the board is finished, is not open by Index.Open, and keeps its bar (milestoneColors only checks isLogged, model.go:1798). Fix: reword to what the code does, e.g. 'and so does a card filed in the logbook: that work has left the board. The M filter still matches such a card; only its colour goes.'
Checked and fine: Index.Open (link.go:318) = tags count (tags.go:129 default lanes, terminal skipped) and = the logbook <milestone> gate today (logbook.go:393: s.Load misses archived/logbooked because refsync Extra drops FiledAwayIDs, refsync.go:1512). JSON: ls top-level count unchanged, per-row count still total, open new; show --json overwrites milestoneJSON tickets with a superset (id, handle + place, lane, title). show reuses milestoneKey, editRefusal, milestoneJSON, milestoneLabel, boardIndex, Whereabouts - no duplicated lookup. Release note: one '- ' line, instruction, matches code. Tests: all four would fail without the change (new symbols / new text); ran them single: ok.
FOLLOW-UP 1: link.Index.terminal uses one lane set (CLI: default board, links.go:115; TUI: the open board, model.go:1930). A member on another board in its own terminal lane counts open; boards.go:52 uses per-board lanes. So link.go:316 'what the board counts call open' and the note's 'picker counts the same way' hold only on the default board. Pre-existing in Satisfied too.
FOLLOW-UP 2: after GWK31K, an archived-elsewhere ref is dropped by List (not open here) but still found by s.Load, so the logbook <milestone> gate (logbook.go:394) would refuse a milestone ls calls '0 open'.
FOLLOW-UP 3 (cheap, may ride the LOOP): docs/COMMANDS.md:143 lists the show wordings without 'on its ref · jaira pull <h>' and 'nowhere here · dangling link', which link.go:501 also prints; .claude/skills/jaira/SKILL.md:213 still lacks 'open' beside docs/AGENTS.md:164.
