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
updated-at: 2026-11-25T10:38:49Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-439134
claimed-at: 2026-11-25T10:32:16Z
outcome-what: "Give back im Cockpit zielt auf Waiting.backTo (backLaneOf in board.ts): rückwärts von der Personen-Lane die erste agentische, nicht-fresh Lane, Loop-Lanes folgen ihrem rejects_to. jaira lanes --json liefert dafür rejects_to. Texte nennen die Ziel-Lane."
outcome-why: "Auf Boards ohne in-progress (Video-Board) scheiterte Give back mit 'no lane in-progress is installed'."
outcome-resolves: "Video-Board: draft-approval -> build, acceptance -> final (Tests); Default-Board weiter in-progress; alle 585 Cockpit-Tests grün; Release-Note geschrieben."
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
