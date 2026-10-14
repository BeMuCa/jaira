---
id: 01M4WS2FVXC09S73WQYZRFTQKP
title: "jaira web: lokales Review-Portal im Binary, zuerst Epics"
status: in-progress
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
commits: []
created-at: 2026-10-14T08:39:42Z
updated-at: 2026-10-14T13:28:17Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-93522
claimed-at: 2026-10-14T13:22:11Z
outcome-what: "jaira web: a local review portal on 127.0.0.1 — open tickets in no epic on the left, epics with description and done X of Y on the right, drag a ticket onto an epic to write its ref and push; fetches every minute and shows the last sync. Svelte 5 + Vite in FSD, built into core/web/ui/dist and embedded; CI job web-dist checks the build matches the sources."
outcome-why: A teamlead wants to see and sort the work by feature without an agent and without the terminal board (F914QH added the epics this page shows).
outcome-resolves: RFTQKP
review-summary: "none"
round: 4
review-gaps: "none: 9be21b4e adds webPush (one caller, no copy of flushRefs' report loop worth sharing: that one prints per-kind advice, this one returns the first failure) and one mutex"
test-verdict: "pass: go vet + go test -race -count=1 ./... RC=0, windows vet/build RC=0; npm ci + build reproduces committed dist (git diff --exit-code RC=0); clean clone builds without node_modules and serves the page; DoD 1-6 verified in the tree; on a scratch bare remote with two clones: epic 0 of 2 with description, unassigned [Unrelated], POST drop -> unassigned 0 and grace's 'jaira epic show' 0 of 3 after fetch, last_fetch/last_push reported, foreign Host 403, cross-origin POST 403, listens on 127.0.0.1 only, Ctrl-C (SIGINT) stops it. Not tried: dragging in a real browser"
review-verdict: "Mostly right, one defect. DoD 1, 3, 4, 5 and 6 match the diff. go test ./core/web ./internal/cli is green, and live on a scratch remote epics, unassigned, sync, the drop and 403 for a foreign Host all work. DoD 2 holds when the remote is reachable. When it is not, the drop is silently only queued and the page says nothing, although help, docs and the release note promise otherwise (gap 1, reproduced). Small fix in internal/cli/web.go srv.Push. My advice: send it back to in-progress for gap 1, or accept it knowingly with a follow-up. Unsure about gap 2 (concurrent flush), not reproduced."
review-round: 2
review-check: "1. cd /home/alex/projects/.worktrees/jaira-rftqkp && go build -o /tmp/jw ./cmd/jaira  2. mkdir /tmp/jwt && cd /tmp/jwt && git init -q --bare r.git && git clone -q r.git a && cd a && git commit -q --allow-empty -m init && git push -q -u origin HEAD  3. /tmp/jw init && /tmp/jw create 'Loose one' --goal g --context c --dod d && /tmp/jw epic create auth --description 'Users sign in'  4. /tmp/jw web -> a browser tab opens on http://127.0.0.1:<port>/; on the left 'Loose one', on the right 'auth' with 'Users sign in' and 0 of 0  5. drag 'Loose one' onto 'auth' -> it disappears on the left, 'auth' shows it and 0 of 1; 'Last sync' at the top shows a time  6. in a second terminal: git -C /tmp/jwt/r.git show refs/jaira/epics/auth:auth.md -> the ticket id is in it (the push arrived)  7. gap 1: in /tmp/jwt/a run git remote set-url origin https://nonexistent.invalid/x.git, create and drag a second ticket -> today NO 'push failed' appears at the top; after the fix it should  8. Ctrl-C in the jaira web terminal -> the command ends"
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
