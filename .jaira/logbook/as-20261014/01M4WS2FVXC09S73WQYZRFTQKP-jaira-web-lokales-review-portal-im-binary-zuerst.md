---
id: 01M4WS2FVXC09S73WQYZRFTQKP
title: "jaira web: lokales Review-Portal im Binary, zuerst Epics"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "jaira web startet eine lokale Seite im Browser: rechts die Epics mit Beschreibung und Fortschritt, links die Tickets ohne Epic. Ein Mensch ohne KI und ohne Board sortiert dort Tickets per Ziehen in Epics, und die Änderung reist über git zu allen."
context: |-
  Was falsch ist: Ein Teamlead will den Fortschritt sehen und steuern, ohne KI und ohne das TUI-Board. Siehe F914QH (Epics), auf dem diese Seite aufbaut.
  Gedacht als kleines Review-Portal, das wachsen darf. Später kann das Board selbst dorthin wandern - das ist NICHT Teil dieses Tickets, aber die Struktur muss es zulassen.
  Entscheidungen Alex 14.10.:
  - Kein Server: das Binary startet die Seite lokal (nur 127.0.0.1). Sync über git: fetch etwa jede Minute, push nach jeder Änderung.
  - Frontend: Svelte 5 + Vite, gebaut zu statischen Dateien, per go:embed ins Binary (wie core/role, core/mod). Kein SvelteKit-SSR; wenn SvelteKit, dann adapter-static.
  - Backend: net/http aus der Standardbibliothek (Go 1.22+ Routing), kein Go-Framework. Handler rufen dieselben core/-Pakete wie die CLI.
  - Frontend-Struktur: Feature-Sliced Design (Schichten app/pages/widgets/features/entities/shared, Imports nur nach unten, öffentliche API je Slice über index.ts).
  Falle: README bietet 'go install ...@latest'. Das baut ohne Node, go:embed liest die Platte. Also muss das gebaute dist/ im Repo liegen, und CI prüft, dass es zu den Quellen passt (neu bauen, git diff --exit-code).
  Grenze: 'kleiner als paca?' - erste Version nur zwei Spalten, Ziehen, Epic-Beschreibung, Fortschritt.
definition-of-done: jaira web startet einen HTTP-Server nur auf 127.0.0.1 und öffnet den Browser; JSON-Endpunkte für Epics und Tickets ohne Epic
tags:
  - cli
blocked-by: []
related: []
follows: 01M4WS20HSNGR8XKP3RCF914QH
commits:
  - e34498d014b8f16e1364c503d7de28fbbc4edc02
  - 063d6b278adf9d3e998bf33cb9a8e911d0681d38
  - 7222c881d6bd03733c08a30150ff9bc1d0d4fabc
  - 667d96160d57a4598a44451d3b0457408da8fbe8
  - b7d536048cbaeac329df7372ce3bea0d0a05b562
  - 9be21b4eef05009a0f745a6e1ffa6bc74f468850
  - b6952b9a465d530d8ea55a393f15003f12f7ccf5
  - 46df5ae2564b65697b7466e1192a649058a1afbc
  - 3425c9074674dce20a697a3a1ae6a6660b9c13e2
  - 2db62fa44aeb0448962b9eadf819e375bf2ae86f
  - d6fee7b86a48f0acfbf0114e1852955517a14c10
  - f3b808fec15de0381a9ef4d19ff7645b8e870de6
created-at: 2026-10-14T08:39:42Z
updated-at: 2026-10-14T15:29:33Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-82065
claimed-at: 2026-10-14T14:58:29Z
outcome-what: "jaira web: a local review portal on 127.0.0.1 — open tickets in no epic on the left, epics with description and done X of Y on the right, drag a ticket onto an epic to write its ref and push; fetches every minute and shows the last sync. Svelte 5 + Vite in FSD, built into core/web/ui/dist and embedded; CI job web-dist checks the build matches the sources."
outcome-why: A teamlead wants to see and sort the work by feature without an agent and without the terminal board (F914QH added the epics this page shows).
outcome-resolves: RFTQKP
review-summary: "Round 2 (Alex's acceptance on 14.10, DoD 7-10). A card opens read-only details: GET /api/tickets/{id} reads the ticket the way 'jaira show' does (Store.LoadAnywhere) — title, lane, goal, context, DoD items with done/open, and the notes. The notes come from the new core/ticket.Notes: it splits '## Progress' at dated headings and gives date, author and text. ticketBrief now uses the same splitter, so there is one parser instead of two. The page renders each note as markdown (marked + DOMPurify, raw HTML escaped beforehand) with its date and author. Only the newest 5 are shown, the rest behind 'Show all'. New writes: POST /api/epics -> epic.NormalizeName + epic.Create (409 if it exists, 400 for a bad or empty name). Dragging between epics is POST /api/epics/{to}/tickets {id, from}: it adds first, then removes, and the error names the step that failed. Dragging back to the left is DELETE /api/epics/{name}/tickets/{id} -> epic.Assign add=false. A follow-up card is not draggable and its card says 'via follows/parent'. Every write pushes once, a failed one too, so a half move still goes out. Docs, help text and the release note gained a line each."
round: 6
review-gaps: "1. Minor: DELETE on a ticket that is simply not in the epic answers 409 'it counts through follows or parent'. That is wrong for that case (reproduced: ticket B, epic 'two'). The page reaches it only from a stale view, e.g. a teammate removed the ticket in the last minute. The status is right, only the message is wrong. 2. Dragging, the dialog and the markdown were checked only in headless Chrome (testing lane). No person has tried them in a real browser yet. 3. A note whose text has a line starting '## ' ends the Progress section. AppendNote already behaved that way and Notes() follows it, so this is not new here. 4. The branch carries the unmerged F914QH, so the F914QH PR must go first."
test-verdict: "pass: go test -race -count=1 ./... RC=0, GOOS=windows build RC=0, npm ci + build leaves dist byte-identical; DoD 7-10 verified in the tree; headless Chrome on a scratch copy of the acceptance board: detail dialog with 5 of 7 notes + Show all, markdown bold rendered, <script>/javascript: inert (window.pwned unset), + New epic creates and a duplicate says 'exists already', drag member to another epic moves it (1/2), drag back left removes it, follow-up has draggable=false and 'via follows/parent', no console errors; API: half-failed move says 'is in pay-later now, but taking it out of no-such failed', foreign Host/Origin 403, remove with remote down reports push_error and goes out with the next change (second clone sees it)"
review-verdict: "The diff meets DoD 7-10, and DoD 1-6 are unchanged. go vet and the tests for core/web, core/ticket and internal/cli pass. The dist rebuilt from npm ci + build is byte-identical. All proof references point at existing lines. Live check on a scratch board with a bare remote: create, 409 on a duplicate, assign, move (both refs on the remote), remove, and detail with author/date notes all behave as specified. No defects block acceptance. One wrong error message (gap 1) is cosmetic."
review-round: 6
review-check: "1. In a scratch repo with a remote, run: go run ./cmd/jaira web   2. The browser opens: tickets in no epic on the left, epics on the right.   3. Click '+ New epic', enter name 'test', title 'Test', a description, submit. A 'Test' card appears on the right.   4. Drag a ticket from the left onto 'Test'. It disappears on the left, 'Test' shows '1 tickets'.   5. Create a second epic, open the 'tickets' fold of 'Test', drag the ticket onto the second epic. It moves. git ls-remote origin 'refs/jaira/epics/*' lists both refs.   6. Drag it back to the left column. It shows there again.   7. Click a ticket title. A dialog shows lane, goal, context, DoD with ☑/☐ and the notes with date and author. A note '<img src=x onerror=alert(1)> **bold**' shows the tag as text and 'bold' in bold, and no alert pops up.   8. A ticket with more than 5 notes shows 'Show all N notes'. Clicking it shows all of them."
---

# jaira web: lokales Review-Portal im Binary, zuerst Epics

## Definition of Done

- [x] jaira web startet einen HTTP-Server nur auf 127.0.0.1 und öffnet den Browser; JSON-Endpunkte für Epics und Tickets ohne Epic
  proof: internal/cli/web.go newWebCmd (web.Listen 127.0.0.1, web.OpenBrowser); core/web/web.go Handler /api/epics /api/unassigned; TestListenIsLoopback, TestEpicsAndUnassigned, TestForeignHostAndOriginRefused
- [x] Seite: rechts Epics mit Beschreibung und erledigt X von Y, links Tickets ohne Epic; Ziehen eines Tickets in ein Epic schreibt die Epic-Ref und pusht
  proof: core/web/ui/src/pages/review/ui/ReviewPage.svelte (UnassignedList left, EpicList right with EpicProgress 'X of Y done'), features/assign-to-epic DropTarget; internal/cli/web_source.go Assign -> epic.Assign + refs.Flush; TestWebShowsEpicsAndADropTravels (drop reaches second clone). Drag in a real browser not yet tried by a person
- [x] Frontend Svelte 5 + Vite in Feature-Sliced Design; dist per go:embed im Binary, 'go install' ohne Node liefert die Seite mit
  proof: core/web/ui/src (app/pages/widgets/features/entities/shared, index.ts per slice); core/web/web.go //go:embed all:ui/dist; TestServesEmbeddedPage
- [x] CI prüft, dass das eingecheckte dist zu den Svelte-Quellen passt
  proof: .github/workflows/ci.yaml job web-dist: npm ci, npm run build, git diff --exit-code core/web/ui/dist
- [x] Hintergrund-fetch etwa jede Minute; die Seite zeigt den Stand der letzten Synchronisierung
  proof: internal/cli/web.go webFetchEvery=1m, Server.Loop; SyncStatus widget shows last_fetch; TestLoopFetchesAtOnceAndStops, TestSyncRecordsFetch; live: /api/sync returned last_fetch on this board
- [x] Tests (Go-Handler), docs/COMMANDS.md, Release-Note core/release/notes/<ID>.md
  proof: core/web/web_test.go, internal/cli/web_test.go; docs/COMMANDS.md row 'jaira web'; core/release/notes/RFTQKP.md
- [x] Klick auf ein Ticket öffnet eine Detailansicht (Titel, Lane, Ziel, Kontext, DoD mit Häkchen), nur lesend; Endpunkt GET /api/tickets/{id} aus denselben core-Daten wie 'jaira show --json'
  proof: core/web/web.go:280 handleTicket (GET /api/tickets/{id}); internal/cli/web_source.go:79 Ticket via Store.LoadAnywhere like 'jaira show'; widgets/ticket-details TicketDetails.svelte in a <dialog> from ReviewPage; TestTicketDetail, TestWebCreatesMovesRemovesAndOpens; headless Chrome: click opens title/lane/goal/context/DoD
- [x] In der Detailansicht alle Notes lesbar: Markdown gerendert, je Note Datum und Autor, chronologisch, lange Verläufe eingeklappt mit 'alle anzeigen'; als FSD-Slice (entities/ticket, widgets/ticket-details)
  proof: core/ticket/body.go:138 Notes, shared with ticketBrief (TestNotesReadsBackWhatAppendNoteWrote, TestTicketBriefCarriesTheNewestNotes); entities/ticket/ui/NoteItem.svelte date+author, shared/lib/markdown.ts marked + raw HTML escaped + DOMPurify; TicketDetails folds all but the newest 5 behind 'Show all'; headless Chrome: <script>, onerror and javascript: links render as text/no href, window.pwned stays unset
- [x] Neues Epic auf der Seite anlegen (Name, Titel, Beschreibung), schreibt refs/jaira/epics/<name> und pusht wie 'jaira epic create'
  proof: core/web/web.go:295 handleCreateEpic (POST /api/epics, 409 exists, 400 no name); internal/cli/web_source.go:101 CreateEpic -> epic.NormalizeName + epic.Create; features/create-epic CreateEpicForm; TestWebCreatesMovesRemovesAndOpens (second clone sees pay-later after fetch)
- [x] Ticket aus einem Epic zurück nach links ziehen nimmt es heraus ('jaira epic rm'); auf ein anderes Epic ziehen verschiebt es (raus aus dem alten, rein ins neue, beide Refs gepusht; scheitert einer, sagt die Seite welcher). Folge-Tickets ohne eigene Zeile sind nicht ziehbar und sagen warum
  proof: internal/cli/web_source.go:121 Assign with from (add first, then remove; error names the step), :142 Unassign, :165 assign -> epic.Assign; core/web/web.go:320 DELETE route, :328 write pushes once; UnassignedList is a DropTarget, TicketCard of a follow-up not draggable with 'via follows/parent'; TestWebCreatesMovesRemovesAndOpens, TestMoveRemoveCreate; headless Chrome drag move and drag back

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/web: Server with EpicSource interface, JSON endpoints /api/epics /api/unassigned /api/assign /api/sync, embed dist, 127.0.0.1 only
- [x] Go handler tests
- [x] internal/cli web command: listen, open browser, background fetch loop, adapter over epic package (stub until F914QH)
- [x] Svelte 5 + Vite app in FSD: shared/api, entities/epic+ticket, features/assign-to-epic, widgets/epic-list+unassigned-list, pages/review, app
- [x] build dist, commit; CI job dist-check
- [x] docs/COMMANDS.md, release note RFTQKP.md
- [x] wire the real core/epic once F914QH lands

## Progress
- **2026-10-14 08:46 · Alexander Sacharov** — parent session: jaira-34
- **2026-10-14 08:47 · Alexander Sacharov** — decision (worker, settled from defaults, no open decision for the person): Go code in core/web (net/http, handlers, embed of ui/dist); Svelte 5 + plain Vite (no SvelteKit) in core/web/ui, FSD layers under src/; dist committed at core/web/ui/dist with package-lock.json. Listens on 127.0.0.1, port 0 (random free) unless --port; --no-open skips the browser. Browser: xdg-open / open / rundll32 (Windows) / wslview or cmd.exe start on WSL. Drag = native HTML5 DnD, no library. Background sync: goroutine every 60s calls refs.Incoming() (fetches all refs/jaira/*, epics included); /api/sync reports last fetch time and error. Assign = epic add + flush (push). Epic data comes from F914QH (worker jaira-f914qh-93), parallel; web defines an EpicSource interface and stubs it until F914QH lands — no second epic model. CI: setup-node, npm ci, npm run build, git diff --exit-code core/web/ui/dist.
- **2026-10-14 08:48 · Alexander Sacharov** — F914QH agreed interface (jaira-f914qh-93, 14.10.): core/epic, epics keyed by Name; LoadAll(root), Count(e, ix *link.Index) Progress{Done,Total,Tickets[]{ID,Title,Lane,Done,Direct}}, Unassigned(epics, ix) []*ticket.Ticket, Assign(s, rec Recorder, name, ids, add) — one call takes lock+saves+queues ref; push = refs.Flush afterwards. done = link.Index.Satisfied. Web stubs against this in internal/cli/web_epics.go until feat/f914qh lands.
- **2026-10-14 08:50 · Alexander Sacharov** — F914QH interface update: epics live on refs only (Alex). epic.LoadAll(r Refs), Assign(s, r Refs, name, ids, add), Unassigned(epics, ix, board); *refsync.Syncer implements Refs; ErrNoRefs on a board without remote. Then refs.Flush(). So the web's background fetch = refs.Incoming() (one fetch of refs/jaira/*) is enough; no disk write needed for epics.
- **2026-10-14 08:56 · Alexander Sacharov** — Pause: core/web + jaira web + Svelte UI + CI + docs + release note committed. Epics side is a stub in internal/cli/web_source.go (Epics empty, Assign errors, Unassigned = all open tickets). Waiting for F914QH (feat/f914qh) to wire core/epic; then tick DoD 1,2,5,6. Added .gitattributes core/web/ui/** eol=lf (wintrap test demands it for embedded files). .jaira/milestones/ untracked in this worktree came from the child 'jaira fetch' — not mine to commit.
- **2026-10-14 08:57 · Alexander Sacharov** — Wiring plan: draft adapter ready (uses F914QH's epicRefs() and epicIndex() from internal/cli/epics.go, so no second copy). F914QH is conversational — a person commits it; waiting for its sha on feat/f914qh, then merge it into feat/rftqkp and replace internal/cli/web_source.go.
- **2026-10-14 13:11 · Alexander Sacharov** — critique round 1: core/web/browser.go:51 isWSL is a copy of core/notify/notify.go:92 — export it once (notify.IsWSL) and call that from browserCommand
core/web/web.go Server.Now is a clock nobody needs: the only test using it asserts LastPush.IsZero(); drop the field and now(), call time.Now()
internal/cli/web.go:119 lastLine (12 lines) for a child that prints one error line — use strings.TrimSpace(stderr.String())
core/web/web.go POST /api/sync + SyncStatus "Sync now" button are outside the DoD (show last sync, not trigger one) — remove the route, handleSyncNow, syncNow() in entities/sync and the button; the minute loop and the push after a drop cover it
- **2026-10-14 13:11 · Alexander Sacharov** — critique let stand: fetch via child 'jaira fetch' instead of in-process refs.Incoming — in-process would update the seen record and swallow the 'ticket is yours' notification the CLI fetch raises; separate /api/epics and /api/unassigned (DoD names both); vite dev proxy (needed to develop the page at all).
- **2026-10-14 13:13 · Alexander Sacharov** — in-progress round 2: all four critique findings applied (notify.IsWSL exported and reused; Server.Now dropped; lastLine -> strings.TrimSpace; POST /api/sync, syncNow and the button removed, dist rebuilt). Proofs re-checked: still name existing tests/lines.
- **2026-10-14 13:14 · Alexander Sacharov** — critique round 2: all four round-1 findings addressed in 667d9616; the repair introduced nothing that breaks the DoD. none.
- **2026-10-14 13:22 · Alexander Sacharov** — review (fresh context): one defect, reproduced. Pushing from 'jaira web' hides an unreachable remote: refs.Flush returns (nil, nil) on ErrOffline, and per-entry Failed/Rejected come back only as reports, which srv.Push in internal/cli/web.go throws away. So push_error stays empty and the page never shows 'push failed — queued'. Repro: origin -> https://nonexistent.invalid/x.git, POST /api/epics/auth/tickets -> 200 without push_error, nothing on the remote. Fix: in srv.Push, return an error when a report is not Sent. Also: pushNow runs outside webSource.mu (two drops = two concurrent Flushes, not reproduced); the branch carries the unmerged F914QH.
- **2026-10-14 13:24 · Alexander Sacharov** — review round 1 sent back from signoff (gap 1, reproduced by the reviewer): a drop with an unreachable remote showed no push_error, because refs.Flush returns nil on ErrOffline and reports Unsent/Failed as results. Fix: internal/cli/web.go webPush turns any non-Sent report and a non-empty outbox into the error; gap 2: core/web Server.pushing mutex serialises flushes. Test TestWebPushReportsAQueuedDrop (remote http://127.0.0.1:1). Found on the way: refSubject in internal/cli/refs.go (F914QH code) names an epic as a ticket handle — reported to jaira-f914qh-93, not edited here to avoid a conflict. Gap 4 (branch carries F914QH) is for the PR order, not code.
- **2026-10-14 13:24 · Alexander Sacharov** — critique (after review fix): read only 9be21b4e — webPush and the push mutex answer gaps 1 and 2; nothing to change. none.
- **2026-10-14 13:31 · Alexander Sacharov** — merged feat/f914qh up to 340e0825 (refSubject names epics; epicIndex became boardIndex in links.go) — web_source.go follows the rename.
- **2026-10-14 13:36 · Alexander Sacharov** — review (fresh context): no defects; checked by hand in headless Chrome against a scratch board with a bare remote - drag moves the ticket and pushes the epic ref, unreachable remote shows push_error and the next drop sends the queue, foreign Host/Origin 403, teammate's epic arrives with the minute fetch, dist rebuild byte-identical. Open: PR order - F914QH is not on upstream/master, so its PR goes first.
- **2026-10-14 13:38 · Alexander Sacharov** — after review round 2 (signoff, no defects): merged feat/f914qh d5b28056 (agent block text only, none of RFTQKP's files); suite green. PR order: F914QH first.
- **2026-10-14 14:56 · Alexander Sacharov** — Abnahme 14.10. (Alex): die Seite läuft, aber man muss ein Ticket öffnen können. Klick auf eine Ticket-Karte (links oder in einem Epic) öffnet eine Detailansicht mit Titel, Lane, Ziel, Kontext und DoD (angehakt/offen). Darunter alle Notes/Kommentare, schön lesbar: Markdown gerendert, Datum und Autor je Note, neueste unten, lange Listen eingeklappt mit 'alle anzeigen'. Nur lesen, kein Editieren. Daten wie 'jaira show <id> --json'.
- **2026-10-14 14:57 · Alexander Sacharov** — Abnahme 14.10. (Alex), Nachtrag: Epics auch auf der Seite anlegen; Tickets aus einem Epic herausnehmen und zwischen Epics verschieben (DoD 9, 10).
- **2026-10-14 14:58 · Alexander Sacharov** — decision (worker, round 2, DoD 7-10, settled from defaults): GET /api/tickets/{id} returns title, lane, goal, context, dod items (text, done), notes parsed from '## Progress' (at, author, text) via a new ticket.Notes(body) beside AppendNote. Markdown rendered in the browser with marked 18.0.14 and the output passed through DOMPurify 3.4.16 (raw HTML in a note is escaped before parsing as well, so a <script> shows as text). Notes beyond the newest 5 fold behind 'show all'. POST /api/epics {name,title,description} -> epic.Create (name normalised like the CLI, 409 if it exists). DELETE /api/epics/{name}/tickets/{id} -> epic.Assign add=false. Move = POST /api/epics/{to}/tickets {id, from}: add to the new epic first, then remove from the old one, so a failure leaves the ticket in both rather than in none; the error names which step failed; one Flush pushes both refs and webPush names a rejected ref. A follow-up (direct=false) is not draggable and its card says 'counted via follows/parent'; the server answers 409 for removing a non-member.
- **2026-10-14 15:09 · Alexander Sacharov** — critique round 1 (round 2 of the ticket): one finding - ticket.Notes duplicates ticketBrief's note splitting in internal/cli/tickets.go; unify on ticket.Notes with a Raw field, splitting at the brief's dated heading rule. Let stand: write() pushes after a failed write too (a half move must still go out, comment says so); from==name guard in webSource.Assign is needed (add no-op then remove would drop the ticket); escaping raw HTML before DOMPurify is the ticket's 'no raw HTML' rule, not redundant.
- **2026-10-14 15:11 · Alexander Sacharov** — critique round 2: read only d6fee7b8 - ticketBrief now takes ticket.Notes' Raw, briefNote regex gone, the brief test passes unchanged; nothing else in the round-2 diff changed. none.
- **2026-10-14 15:12 · Alexander Sacharov** — optimize (round 2): removed three type exports nobody imports across a slice; bundle byte-identical. Looked at and kept: TicketDetail.id (the page's identity of the opened ticket), the from==name guard, escapeHTML before DOMPurify (see critique note).
- **2026-10-14 15:16 · Alexander Sacharov** — testing round 2: pass (see test-verdict). Scratch copy of the acceptance fixture under the session scratchpad, the acceptance fixture itself untouched.
- **2026-10-14 15:18 · Alexander Sacharov** — review round 3 (fresh context, DoD 7-10): no blocking defects. Live check on a scratch board with a bare remote: create/409/assign/move/remove/detail all as specified, dist byte-identical, tests green. Minor: the DELETE 409 message says 'counts through follows or parent' even when the ticket is simply not in that epic. Open: no person has tried it in a real browser; PR order F914QH first.
- **2026-10-14 15:20 · Alexander Sacharov** — acceptance round 2: block 2 on https://claude.ai/artifact/TzzQLMxeCQuumLtkD43Lxk (v4) gained b2-s06..s13 (detail, notes markdown + XSS note, show all, + New epic and duplicate, follow-up not draggable, move between epics, second clone sees both refs, drag back left); fixture rebuilt from scratch by ~/.claude/acceptance/jaira/feat-f914qh/rftqkp-fixture.sh with the f3b808fe binary, so b2-s01..s05 start from the round-1 state again.
