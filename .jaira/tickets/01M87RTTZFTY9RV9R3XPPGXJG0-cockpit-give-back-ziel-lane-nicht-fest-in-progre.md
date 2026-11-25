---
id: 01M87RTTZFTY9RV9R3XPPGXJG0
title: "Cockpit Give back: Ziel-Lane nicht fest in-progress"
status: critique
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
commits: []
created-at: 2026-11-24T21:54:16Z
updated-at: 2026-11-25T10:47:19Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-439134
claimed-at: 2026-11-25T10:32:16Z
outcome-what: "Give back im Cockpit zielt auf Waiting.backTo (backLaneOf in board.ts): rückwärts von der Personen-Lane die erste agentische, nicht-fresh Lane, Loop-Lanes folgen ihrem rejects_to. jaira lanes --json liefert dafür rejects_to. Texte nennen die Ziel-Lane."
outcome-why: "Auf Boards ohne in-progress (Video-Board) scheiterte Give back mit 'no lane in-progress is installed'."
outcome-resolves: "Video-Board: draft-approval -> build, acceptance -> final (Tests); Default-Board weiter in-progress; alle 585 Cockpit-Tests grün; Release-Note geschrieben."
review-summary: "Das Cockpit rechnet pro Board ein backTo für jedes wartende Ticket (backLaneOf in hooks/board.ts): von der Personen-Lane rückwärts die erste agentische, nicht-fresh, nicht-Personen-Lane, bei einer Loop-Lane deren rejects_to. Give back, alle Texte (texts.ts, sessions.ts backMessage) und die Zeile an die Teamlead-Session nutzen dieses Ziel statt fest in-progress. jaira lanes --json liefert dafür neu rejects_to; ohne das Feld (älteres jaira) bleibt in-progress. Default-Board bleibt in-progress, Video-Board landet in build/final. 585 Cockpit-Tests und go test ./internal/cli grün."
review-gaps: |-
  note: edge: fix-now: backLaneOf nimmt rejects_to[0] ungeprüft (hooks/board.ts:121) - ein Ziel, das nicht auf dem Board steht oder eine Personen-Lane ist (rejects_to [human, in-progress]), lässt den Move scheitern bzw. landet bei der Person; fix: erstes rejects_to-Ziel nehmen, das im lanes-Array steht, agentisch, nicht fresh und nicht isPersons ist, sonst l.id; plus ein Test dafür
  note: elsewhere: discussInThisChat (hooks/texts.ts) und die Discuss-Lane in spawn.sh sagen weiter move --to in-progress - auf einem Board ohne in-progress scheitert dieser Schritt nach einer beantworteten Frage; neues Ticket
  note: limit: backTo wird aus dem Board-Snapshot gelesen; ändert jemand zwischen Anzeige und Druck die Lanes, geht der Move an die alte Ziel-Lane oder wird von jaira abgelehnt
test-verdict: "pass: claude plugin test . 585 pass RC=0; go test internal/cli core/shipped core/lane RC=0; DoD 1-4 checked in the tree (named tests, texts.ts/sessions.ts, release note); jaira lanes --json lists rejects_to on the repo board and a fresh init board (fresh board without loop lanes still resolves to in-progress); unknown board exits 2; writes only via jaira move, path unchanged"
hunt-verdict: "pass: bughunt done at 876a1f3e0, 1 change (edge, not a blocker: rejects_to[0] unchecked), noted for review and person, 0 pre-existing"
review-verdict: "accept: alle vier DoD-Punkte erfüllt, Tests grün; ein fix-now-Rand bei rejects_to[0]"
review-round: 1
review-verdict-gaps: |-
  note: edge: fix-now: backLaneOf nimmt rejects_to[0] ungeprüft (hooks/board.ts:121) - ein Ziel, das nicht auf dem Board steht oder eine Personen-Lane ist (rejects_to [human, in-progress]), lässt den Move scheitern bzw. landet bei der Person; fix: erstes rejects_to-Ziel nehmen, das im lanes-Array steht, agentisch, nicht fresh und nicht isPersons ist, sonst l.id; plus ein Test dafür
  note: elsewhere: discussInThisChat (hooks/texts.ts) und die Discuss-Lane in spawn.sh sagen weiter move --to in-progress - auf einem Board ohne in-progress scheitert dieser Schritt nach einer beantworteten Frage; neues Ticket
  note: limit: backTo wird aus dem Board-Snapshot gelesen; ändert jemand zwischen Anzeige und Druck die Lanes, geht der Move an die alte Ziel-Lane oder wird von jaira abgelehnt
review-check: "1. cd core/mod/builtin/jaira-cockpit && claude plugin test .  2. du siehst 585 pass, 0 fail, darunter 'Give back on a video board moves the ticket into its build lane and says so'  3. jaira lanes --json | jq '.[] | {id, rejects_to}'  4. du siehst critique/optimize/testing mit [\"in-progress\"], alle anderen mit null  5. Live (optional): im Cockpit ein Ticket in signoff per Give back zurückgeben - es landet in in-progress, die Meldung nennt in-progress"
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
