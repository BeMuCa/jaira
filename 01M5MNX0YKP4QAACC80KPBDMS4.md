---
id: 01M5MNX0YKP4QAACC80KPBDMS4
title: Milestone-Dateien fuehren keine archivierten Tickets mehr als offen
status: pre-process
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
updated-at: 2026-10-27T00:04:03Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: pbdms4-driver
claimed-at: 2026-10-26T20:22:44Z
round: 2
---

# Milestone-Dateien fuehren keine archivierten Tickets mehr als offen

## Definition of Done

- [ ] jaira milestone ls zaehlt archivierte und abgelegte Tickets nicht als offen und zeigt sie in show als 'archiviert'/'abgelegt'; ein Test baut den Fall
- [ ] Die Kartenkante und der M-Picker ignorieren sie
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [ ] core/link/link.go: add (ix *Index) Open(id) bool beside Satisfied (~l.304): true only when place() is PlaceBoard or PlaceRef and not Done - 'open' as jaira tags/boards count it (tags.go:126, boards.go:47); reads no file
- [ ] core/link/link_test.go: TestOpen... - backlog ticket open; done-lane ticket, s.Archive'd ticket, s.Logbook'd ticket and unknown id not open (fixtures store/make1/list, l.10-60); go test ./core/link -run TestOpen -count=1
- [ ] core/milestone/milestone.go: add (m *Milestone) Open(isOpen func(id string) bool) int beside Members (l.109), counting members; file lines stay untouched (removing ids is N71NVG); small test in milestone_test.go
- [ ] failing test first in internal/cli/milestones_test.go after TestLsNamesEveryMilestoneWithItsCount (l.245): 3 tickets via mkTicket, milestone create, runCLI archive h2, ticket.At(dir).Logbook(h3, "as-20260911") (pattern links_test.go:47); ls text has '1 open of' and '3 ticket(s)', ls --json open==1 and count==3
- [ ] internal/cli/milestones.go ls (l.505-560): one boardIndex(s) (links.go:113); line '%3d ticket(s)' becomes '%3d open of %3d ticket(s)' (keeps milestones_test.go:239/262 green); --json adds "open" after milestoneJSON, "count" stays the total; update the Long help
- [ ] NEEDS DRIVER OK (no 'milestone show' exists, milestones.go:80-87): newMilestoneShowCmd modelled on epic show (epics.go:259-309), registered in newMilestoneCmd: label, goal, 'N open of M ticket(s)', one row per member with handle, ix.Lookup(id).Whereabouts() (link.go:493: 'archive · unfinished/done', 'logbook · done', 'on the board · <lane>'), title; --json = milestoneJSON + open + tickets with place, lane, title; name lookup and refusals as the other milestone commands
- [ ] extend the CLI test: show marks h2 'archive ·' and h3 'logbook · done', h1 'on the board · backlog'; show --json tickets carry place archive/logbook/board
- [ ] TUI M picker: internal/tui/view.go:1646 count becomes open of total via ms.Open(m.linkIndex().Open) (model.go:1926); test in internal/tui/milestonebar_test.go with pickerRow (l.192), mkMilestone (l.13), m.store.Archive (pattern layer_test.go:751); go test ./internal/tui -run <that test> -count=1
- [ ] TUI card edge: milestoneColors (model.go:1794) returns empty slots for m.isLogged(t) (logbook.go:37); do NOT filter m.mstones (model.go:431), the M filter matches logbook cards on purpose (model.go:567-569); test with fileIntoLogbook (logbook_test.go:16) that the logged card has no milestone colour
- [ ] docs: docs/COMMANDS.md:142 ls row (open of total, json open) plus a show row; docs/AGENTS.md:164; README.md:862 command list
- [ ] release note core/release/notes/PBDMS4.md: exactly one '- ' line, unwrapped, an instruction: ls reads open of total (json open), show lists where each ticket is now, M picker counts the same, logbook cards drop the milestone colour
- [ ] per commit gofmt -l . empty, go vet ./... and go build ./... green; the full go test ./... -race runs once before the PR (driver ticks DoD 3)

## Progress

