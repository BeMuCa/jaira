---
id: 01M87RTTZFTY9RV9R3XPPGXJG0
title: "Cockpit Give back: Ziel-Lane nicht fest in-progress"
status: done
ready: true
creator: Sascha
goal: "Give back im Cockpit schickt ein Ticket in die Arbeits-Lane seines eigenen Boards zurück, auch wenn das Board keine Lane in-progress hat."
context: "Give back im Cockpit scheitert auf jedem Board ohne Lane in-progress.\nGesehen 24.11.2026 auf einem Video-Board (Lanes: todo, script, script-approval, build, draft-approval, voice-approval, voice, final, acceptance, done, blocked):\n  jaira: cannot move <ID> to in-progress: no lane \"in-progress\" is installed\nUrsache: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1757 ruft fest jaira move <id> --to in-progress. Auch texts.ts:106/120/151/245 und sessions.ts:170 nennen in-progress fest.\nM5JSE9 (0.6.7) hat Accept schon auf next_lane umgestellt; Give back ist die Gegenrichtung und fehlt noch.\nOffen: welche Lane ist das Ziel? Vorschlag: die letzte Agenten-Lane vor der Personen-Lane, in der das Ticket steht (Video: build vor draft-approval, final vor acceptance); auf dem Default-Board bleibt das in-progress."
definition-of-done: "Give back auf einem Board ohne in-progress bewegt das Ticket in die Agenten-Lane vor der aktuellen Personen-Lane, Test in tests/"
tags:
  - cockpit
  - lanes
blocked-by: []
related: []
commits:
  - 876a1f3e0135451b3a8067529ab8abf5a35222aa
  - 62df7c16fd1bdbe2f09301fe1786273ea10ea0c1
created-at: 2026-11-24T21:54:16Z
updated-at: 2026-11-25T11:29:17Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-439134
claimed-at: 2026-11-25T10:32:16Z
outcome-what: "Give back im Cockpit zielt auf Waiting.backTo (backLaneOf in board.ts): rückwärts von der Personen-Lane die erste agentische, nicht-fresh Lane, Loop-Lanes folgen ihrem rejects_to. jaira lanes --json liefert dafür rejects_to. Texte nennen die Ziel-Lane."
outcome-why: "Auf Boards ohne in-progress (Video-Board) scheiterte Give back mit 'no lane in-progress is installed'."
outcome-resolves: "Video-Board: draft-approval -> build, acceptance -> final (Tests); Default-Board weiter in-progress; alle 585 Cockpit-Tests grün; Release-Note geschrieben."
review-summary: "Delta review from 876a1f3e0 (commit 62df7c16f, 15+/3-): backLaneOf in hooks/board.ts now picks the first rejects_to entry that is on the board, agent-worked, not fresh and not a person's lane, and falls back to the loop lane itself when none qualifies. This closes the round-1 fix-now: a rejects_to naming a missing lane (video board without in-progress) or a person's lane (human first) no longer becomes the Give back target. A new test covers both cases; shipped.sha256 carries the new board.ts hash (verified with sha256sum). claude plugin test: 586 pass, 0 fail."
review-gaps: "none"
test-verdict: "pass: round 2 fix only - claude plugin test . 586 pass RC=0 incl. 'a rejects_to naming a lane not on the board, or a person's lane, is passed by'; go test core/shipped RC=0; DoD 1-4 proofs unchanged and still true"
hunt-verdict: "pass: stop-after-fix, no hunt in round 2"
review-verdict: "accept: round-1 fix-now closed in 62df7c16f, DoD 1-4 met, 586 cockpit tests green"
review-round: 2
review-verdict-gaps: none
review-check: "1. cd core/mod/builtin/jaira-cockpit  2. claude plugin test .  3. expect 586 pass, 0 fail, including 'a rejects_to naming a lane not on the board, or a person's lane, is passed by' and 'Give back on a video board moves the ticket into its build lane and says so'  4. in the cockpit on the default board, press Give back on a ticket in signoff: it lands in in-progress and the toast names in-progress  5. on a board without in-progress (e.g. a video board), Give back from draft-approval lands in build and the toast names build"
round: 2
back-from: review
back-by: lane
back-count: 1
stop-after-fix: "true"
---

# Cockpit Give back: Ziel-Lane nicht fest in-progress

## Definition of Done

- [x] Give back auf einem Board ohne in-progress bewegt das Ticket in die Agenten-Lane vor der aktuellen Personen-Lane, Test in tests/
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'Give back on a video board moves the ticket into its build lane and says so'; backLaneOf in hooks/board.ts
- [x] Default-Board: Give back landet weiter in in-progress (bestehende Tests grün)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'the shipped board gives back to in-progress, past the review, the human lane and the loop lanes'; claude plugin test . 585 pass
- [x] Texte in texts.ts und sessions.ts nennen die echte Ziel-Lane statt in-progress
  proof: hooks/texts.ts returnedTo/tabsUnreadable/workerBlocked/toSessionGivenBack/noteLost and hooks/sessions.ts backMessage take the target lane
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/PGXJG0.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-25 10:32 · Sascha** — parent session: uds:/run/user/1000/cc-socks/431929.sock
- **2026-11-25 10:33 · Sascha** — decision (worker): Ziel von Give back = rückwärts von der Personen-Lane in Board-Reihenfolge die erste agentische, nicht-fresh, nicht-Personen-Lane; eine Loop-Lane (rejects_to gesetzt, z.B. testing/optimize/critique) verweist auf ihr rejects_to-Ziel. Default-Board: signoff <- review(fresh) <- human <- testing -> in-progress. Video: draft-approval -> build, acceptance -> final. rejects_to kommt dafür in 'jaira lanes --json' (bisher nur in 'lanes show'); ohne das Feld (älteres jaira) bleibt in-progress wie bisher. Keine offene Frage für die Person: der Vorschlag im Kontext deckt es.
- **2026-11-25 10:38 · Sascha** — Nicht angefasst: discussInThisChat (texts.ts) und die Discuss-Lane von spawn.sh sagen weiter 'move --to in-progress' - das ist die Antwort auf eine Frage im human-Lane, nicht Give back. Auf einem Board ohne in-progress scheitert dieser Schritt noch; gehört zu einem eigenen Ticket, falls gewünscht. shipped.sha256 um die neuen Hashes der 4 Cockpit-Dateien ergänzt (wie M5JSE9).
- **2026-11-25 10:39 · Sascha** — critique: Form geprüft. Alternative 'jaira rechnet back_lane wie next_lane' verworfen: TUI hat kein festes in-progress (grep internal/tui, core/move), einziger Nutzer ist das Cockpit, das lanes --json schon liest; ein Feld mehr in jedem Ticket-JSON und Go-Logik für einen Aufrufer. backLaneOf liegt neben isPersons in board.ts, wo die Lane-Lesung schon lebt. noteLost sagt jetzt 'moved to <lane>' statt 'went back to in-progress', weil auch Accept es benutzt (dort war 'went back' falsch).
- **2026-11-25 10:41 · Sascha** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/board.ts:121 - class: backLaneOf takes rejects_to[0] without checking that the target is installed on the ticket's board or is a lane an agent works; jaira only warns about a rejects-to naming an uninstalled lane (core/lane/lane.go:1154-1163), and a list may name a person's lane first (template.go:31 allows [in-progress, human] in any order); example 1: a video board [script, build, testing(rejects_to in-progress), draft-approval] - Give back from draft-approval targets in-progress, which that board lacks, so the move is refused; example 2: critique with rejects_to [human, in-progress] - Give back from signoff targets human, a person's lane, not the worker's; ran: node --experimental-strip-types scratch script calling backLaneOf on both lane lists (printed in-progress and human); not a blocker: needs a board jaira already warns about, or a reversed list; fix: take the first rejects_to entry that is on the board, agentic, not fresh and not isPersons, else fall back to the lane itself
- **2026-11-25 10:41 · Sascha** — hunt report: 1 finding, 1 change (edge, not a blocker), note on PGXJG0 | read whole: board.ts, sessions.ts, texts.ts (diffed parts), register.tsx give-back/accept paths (giveBack, givingBack, giveBackWith, acceptWith, accept/accepting, deciding, execute, the board snapshot loader), tickets.go lanes --json, lane.go RejectsTo + Validate | not checked: a live Give back in the cockpit on a real second board (needs Herdr + a configured video board); the reproduction ran against backLaneOf directly
- **2026-11-25 10:41 · Sascha** — bughunt: done at 876a1f3e0 - 1 finding (1 change, 0 pre-existing) in 8 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/sessions.ts core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/PGXJG0.md core/shipped/shipped.sha256 internal/cli/tickets.go
- **2026-11-25 10:42 · Sascha** — testing: note: edge: backLaneOf nimmt rejects_to[0] ungeprüft (Hunter-Fund board.ts:121); trifft nur ein Board, auf dem jaira schon vor rejects-to auf eine fehlende Lane warnt, oder eine Liste mit Personen-Lane zuerst; fix: erstes rejects_to-Ziel, das auf dem Board, agentisch, nicht fresh und nicht isPersons ist, sonst die Lane selbst
- **2026-11-25 10:42 · Sascha** — waiting: review lane (fresh agent via run-lane.sh) since 11:42
- **2026-11-25 10:45 · Sascha** — review: done at 876a1f3e0
- **2026-11-25 10:45 · Sascha** — review: fix-now: backLaneOf (hooks/board.ts:121) soll das erste rejects_to-Ziel nehmen, das auf dem Board steht, agentisch, nicht fresh und nicht isPersons ist, sonst l.id; Test mit rejects_to [human, in-progress] und mit einem Ziel, das auf dem Board fehlt
- **2026-11-25 10:46 · Sascha** — stop-after-fix=true gesetzt (worker): fix-now unter 50 Zeilen, kein Lock, kein Parser
- **2026-11-25 10:47 · Sascha** — critique+optimize (ein Pass, stop-after-fix): Fix 62df7c16f, 15+/3- - works() bündelt die drei Bedingungen, die die Schleife schon prüfte, und prüft damit auch das rejects_to-Ziel; nichts doppelt, nichts zu entfernen.
- **2026-11-25 10:48 · Sascha** — waiting: review lane round 2 (new fresh agent via run-lane.sh) since 11:48
- **2026-11-25 10:49 · Sascha** — review: done at 62df7c16f
- **2026-11-25 10:50 · Sascha** — report: Give back zielt auf backLaneOf (Agenten-Lane vor der Personen-Lane, Loop-Lanes über geprüftes rejects_to), jaira lanes --json liefert rejects_to, Texte nennen die Ziel-Lane | signoff, context unknown | offen: Discuss-Pfad (discussInThisChat, spawn.sh discuss) sagt weiter in-progress - eigenes Ticket, falls gewünscht; Ticket-Datei uncommitted bis Logbook
- **2026-11-25 11:26 · Sascha** — cockpit: accepted by the person, signoff -> done
