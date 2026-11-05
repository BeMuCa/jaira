---
id: 01M6NNAXQ276DMK3A1GH3XXF5B
title: "jaira web: ein Ticket laesst sich per Knopf aus einem Epic nehmen, auch auf dem Handy"
status: review
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer in jaira web ein Epic ansieht, nimmt jedes Ticket mit einem Klick oder Tipp heraus, ohne auf die CLI ('jaira epic rm') auszuweichen."
context: |-
  Heute nimmt man ein Ticket nur per Drag-and-Drop aus einem Epic: Karte aus der Epic-Liste auf 'Tickets in no epic' ziehen (core/web/ui/src/widgets/epic-list/ui/EpicList.svelte:74, Ziel UnassignedList.svelte).
  Gemeldet von Alex am 05.11.2026: 'manchmal kann ich nicht entfernen, dann muss ich die CLI nehmen'.
  Bekannte Ursachen:
  - Drag-and-Drop geht auf dem Telefon nicht (Touch hat kein HTML5-dragstart). Alex arbeitet gerade vom Telefon aus.
  - Karten mit direct=false (ueber follows/parent mitgezaehlt) sind nicht ziehbar, TicketCard.svelte:23. Sie stehen nicht in der Ticketliste des Epics, also gibt es dort nichts zu entfernen - die Karte muss das sagen und zeigen, welches Ticket sie hineinzieht.
  - Archivierte Tickets (place != board) sind in der Unassigned-Liste nicht ziehbar.
  Die API gibt es schon: DELETE /api/epics/{name}/tickets/{id} (core/web/web.go:280, handleUnassign) - CLI-Gegenstueck 'jaira epic rm'.
  Pruefen, ob ein archiviertes Ticket aus einem Epic genommen werden kann (CLI vs. Web gleich).
definition-of-done: "Jede direkte Ticket-Karte in einem Epic hat einen Entfernen-Knopf, der DELETE /api/epics/{name}/tickets/{id} ruft; funktioniert per Tipp ohne Drag; Test"
tags:
  - web
blocked-by: []
related: []
commits: []
created-at: 2026-11-05T10:51:07Z
updated-at: 2026-11-05T11:07:37Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-85123
claimed-at: 2026-11-05T10:51:26Z
outcome-what: "Epic cards in jaira web get a Remove button (tap, no drag); follow-up cards name the member they are counted through (new 'via' field from core/epic); the web resolves filed tickets with LoadAnywhere like 'jaira epic rm'."
outcome-why: "Removing only worked by HTML5 drag, which phones lack; dashed cards gave no hint which line to remove; logbook/archive tickets answered 404 in the web while the CLI removed them."
outcome-resolves: "Taking any ticket out of an epic in jaira web, on a phone too, without falling back to 'jaira epic rm'."
review-summary: "Epic cards in jaira web now carry a 'Remove' button: a plain click/tap button (TicketCard.svelte) shown only on direct member cards (EpicList.svelte), wired through ReviewPage to the existing removeFromEpic, i.e. DELETE /api/epics/{name}/tickets/{id}, so taking a ticket out no longer needs a drag and works on a phone. core/epic covered() now also returns, for each follow-up, the member it was reached from (transitively, cmp.Or keeps the top member); it is passed out as a new 'via' field on epic.Ticket and web.Ticket, and a dashed follow-up card says 'counted through <handle> (follows/parent) - take that one out to remove this', the handle opening that member. webSource.load switched from Store.Load to LoadAnywhere, so DELETE (and POST assign) find logbook and archived tickets the way 'jaira epic rm/add' do instead of answering 404. Help text of 'jaira web', release note core/release/notes/3XXF5B.md and the rebuilt ui/dist come along."
round: 3
review-gaps: |-
  note: the Remove button itself has no automated test (core/web/ui has no Svelte component runner); the server side is covered, a real tap on a phone was not exercised
  note: TicketCard grid is now four columns on every card, also in the unassigned list where no button appears - the empty fourth column only adds a gap
  note: LoadAnywhere also applies to Assign, so the API now accepts adding a filed ticket to an epic - matches 'jaira epic add', the page itself never drags a filed ticket
  note: Remove has no confirm and no undo message - deliberate per the 10:52 decision note, re-adding is a drag
test-verdict: "pass: go build rc=0, go test -count=1 ./... rc=0, -race on core/epic core/web internal/cli rc=0, npm test rc=0; DoD 1-4 verified in the tree; scratch board via jaira web: DELETE member 200, logbook member 200, follow-up 409, /api/epics carries via; EpicList server-rendered: Remove on both direct cards (board+logbook), none on the follow-up, which names its member; web --help updated. Not run: a real tap in a phone browser"
review-verdict: "The diff meets all four DoD items: Remove button on direct cards calling DELETE (tap, no drag), follow-up cards name the member via the new 'via' field, filed tickets removable as with 'jaira epic rm' (TestWebRemovesAFiledTicketFromAnEpic), release note present. go test core/epic core/web internal/cli and npm test (13 pass) green here; built bundle matches src. No blocker; the only unverified part is the tap on a real phone."
review-round: 3
review-check: "1. cd /home/alex/projects/.worktrees/jaira-3xxf5b && go build -o /tmp/jaira-3xxf5b ./cmd/jaira  2. in a repo with an epic run /tmp/jaira-3xxf5b web --no-open and open the printed URL on your phone (same network) or with the browser's phone emulation  3. unfold an epic: every solid (member) card shows a 'Remove' button; dashed cards show none  4. tap 'Remove' on a member: the card leaves the epic, the done-count drops, and 'jaira epic show <epic>' no longer lists it  5. on a dashed card read 'counted through XXXXXX (follows/parent)'; tap XXXXXX - the member opens  6. for a logbook ticket in an epic, tap 'Remove': it goes, with no 'not found' error"
---

# jaira web: ein Ticket laesst sich per Knopf aus einem Epic nehmen, auch auf dem Handy

## Definition of Done

- [x] Jede direkte Ticket-Karte in einem Epic hat einen Entfernen-Knopf, der DELETE /api/epics/{name}/tickets/{id} ruft; funktioniert per Tipp ohne Drag; Test
  proof: EpicList.svelte:79 onremove on every direct card -> ReviewPage removeFromEpic (DELETE /api/epics/{name}/tickets/{id}); button TicketCard.svelte:34 is a plain click/tap button; tests TestWebCreatesMovesRemovesAndOpens (DELETE 200 for a member, 409 for a follow-up). No Svelte component test runner in core/web/ui - the button itself is checked by hand.
- [x] Eine indirekte Karte (via follows/parent) nennt das Ticket, ueber das sie mitgezaehlt wird, statt still nicht ziehbar zu sein; Test
  proof: core/epic/epic.go Ticket.Via from covered(); TicketCard.svelte:54 'counted through <handle>'; tests TestCoveredWalksFollowsAndParent (via two steps down), TestWebCreatesMovesRemovesAndOpens (via in /api/epics), handle.test.ts
- [x] Ein archiviertes Ticket laesst sich im Web genauso aus dem Epic nehmen wie mit 'jaira epic rm', oder die Seite sagt warum nicht; Test
  proof: internal/cli/web_source.go load -> LoadAnywhere, same as resolveRef for 'jaira epic rm'; TestWebRemovesAFiledTicketFromAnEpic (404 before, 200 after)
- [x] Release-Note-Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/3XXF5B.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/epic: Ticket.Via = the member a follow-up is counted through; Covered/Count carry it; Go test
- [x] core/web + internal/cli: web.Ticket.Via passed through; webSource.load uses LoadAnywhere so archived/logbook tickets can be removed like 'jaira epic rm'; Go test
- [x] UI: Remove button on direct cards in EpicList calling DELETE via onunassign; indirect card text names the via ticket; vitest
- [x] release note core/release/notes/3XXF5B.md; build ui dist; go test ./...

## Progress
- **2026-11-05 10:51 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/1770.sock
- **2026-11-05 10:52 · Alexander Sacharov** — decisions (none open for the person, settled from code): (1) Remove button = small text button 'Remove' on every direct card in an epic, no confirm dialog - drag has none either and 'add' undoes it; aria-label names ticket and epic. (2) An indirect card names the epic member it is counted through (the line you would remove), found by walking follows/parent up in epic.Covered; new field 'via' on epic.Ticket and web.Ticket. (3) Archived: root cause found - webSource.load uses Store.Load (board only), 'jaira epic rm' resolves via LoadAnywhere (board+logbook+archive). Web switches to LoadAnywhere, so web and CLI behave the same; no 'why not' message needed.
- **2026-11-05 10:54 · Alexander Sacharov** — archived/logbook: reproduced - DELETE on a filed ticket answered 404 'not found' before load switched to LoadAnywhere (TestWebRemovesAFiledTicketFromAnEpic fails on the old line).
- **2026-11-05 10:57 · Alexander Sacharov** — critique (round 1): drop epicCardExit - one caller, restates ticket.direct, via half unused; EpicList back to ticket.direct; handleOf alone into model/handle.ts, not exported. Let stand: LoadAnywhere also for Assign (CLI 'epic add' resolves the same way), via computed in core/epic not the UI (the UI cannot see follows/parent), no confirm dialog (decision note 10:52).
- **2026-11-05 10:58 · Alexander Sacharov** — in-progress round 2: critique fix done - epic-card.ts gone, EpicList keys on ticket.direct, handleOf in model/handle.ts; proofs of DoD 1/2 rewritten (they named the deleted test). The button has no component test: core/web/ui only runs node --test on .ts, no Svelte test runner - the review-check should ask for a tap on a phone.
- **2026-11-05 10:58 · Alexander Sacharov** — critique round 2: round-1 finding addressed in 1ff410d2 (epicCardExit gone, EpicList on ticket.direct, handleOf unexported in model/handle.ts). Nothing further.
- **2026-11-05 11:02 · Alexander Sacharov** — testing round 1: gates green - go build ./... rc=0, go test -race -count=1 core/epic core/web internal/cli rc=0, go test -count=1 ./... rc=0 (33 ok), npm test rc=0. DoD 1-4 verified in tree. Function on a scratch board (bare remote, jaira web --no-open): /api/epics gives the follow-up via=<member id>, members via empty; DELETE member 200, DELETE logbook member 200, DELETE follow-up 409 'counts through follows or parent'; epic show afterwards 0 of 0. Finding: internal/cli/web.go:36-42 help text says 'drag ... back to the left to take it out' and 'a ticket that counts only because it follows a member ... cannot be dragged' - add that every member card has a Remove button (a tap, works on a phone, filed tickets too) and that a follow-up's card names the member that brings it in.
- **2026-11-05 11:03 · Alexander Sacharov** — in-progress round 3: web --help (internal/cli/web.go) now names the Remove button, filed tickets and the follow-up's member line; internal/cli tests rc=0.
- **2026-11-05 11:03 · Alexander Sacharov** — critique round 3: read only f370fa68 (help text answering the testing finding) - says what the page does, no new shape. Nothing.
- **2026-11-05 11:05 · Alexander Sacharov** — testing round 2: help-text finding fixed (f370fa68, seen in the rebuilt binary's web --help). EpicList rendered with svelte/server through vite ssrLoadModule (DOMPurify stubbed, script kept in the scratchpad): aria-label 'Take MEMBER out of Demo' and 'Take LOGGED out of Demo', follow-up shows 'counted through MEMBER (follows/parent)' and no button. A real tap on a phone was not exercised - for the review-check.
- **2026-11-05 11:07 · Alexander Sacharov** — review: no blocker. Diff meets DoD 1-4; go test core/epic core/web internal/cli + npm test green; dist bundle matches src. Unverified only: a real tap on a phone (in review-check).
