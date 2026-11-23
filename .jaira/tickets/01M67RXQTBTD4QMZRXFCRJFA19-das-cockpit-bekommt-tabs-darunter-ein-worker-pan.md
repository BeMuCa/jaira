---
id: 01M67RXQTBTD4QMZRXFCRJFA19
title: "Das Cockpit bekommt Tabs, darunter ein Worker-Panel: je Agent ein Icon, sein Ticket, seine Lane und ob er festhaengt"
status: testing
ready: true
creator: berkc
assignee: Alexander Sacharov
goal: "Im jaira-cockpit gibt es Tabs; ein Tab 'Worker' zeigt je laufendem Agenten ein kleines Icon, darunter das Ticket und die Lane, in der er gerade ist, und das Icon wechselt, wenn der Agent blockiert ist oder Hilfe braucht (z.B. wartet auf Freigabe eines Bash-Befehls)."
context: |-
  Was fehlt: wer gerade was macht, sieht man nur in den Agenten-Meldungen im Chat. Ob ein Worker an einer Berechtigungsfrage haengt, sieht man gar nicht.
  Ausloeser: Maintainer 31.10.2026, Punkt 5: 'Workerpanel zb. Da sieht man dann die agenten als kleine icons und das ticket was bearbeitet wird und darunter die lane wo er grad ist. Er zeigt aber auch ob er blockiert ist, oder assistenz braucht ... das koennten zb bash befehle oder aehnliches sein.'
  Bekannt: das Cockpit hat schon ein Tool workers (mcp__jaira-cockpit__workers) und /cockpit workers (hooks/register.tsx:1112). Ob die Plugin-API Tabs und den Wartezustand eines Subagenten (Freigabe offen) liefert, ist nicht geprueft - zuerst klaeren.
  Welche anderen Tabs: im Brainstorm festlegen (Kandidaten: Wartet auf mich, Worker, Threads).
definition-of-done: "Brainstorm hat Tabs, Icon-Zustaende und die Datenquelle fuer 'blockiert/braucht Hilfe' festgelegt"
tags:
  - tui
blocked-by: []
related:
  - 01M5M7JTVR8G7DZ1MGHJQ63EEV
commits: []
created-at: 2026-10-31T01:24:27Z
updated-at: 2026-11-23T19:51:56Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-40650
claimed-at: 2026-11-23T19:41:00Z
outcome-what: "Worker-Tab und workers-Tool zeigen 'waits for approval' ('!' in warning), wenn Herdr fuer den Agenten-Tab agent_status=blocked meldet; Kopfzeile zaehlt solche Agenten weiter"
outcome-why: Ein Agent am Freigabe-Dialog sah aus wie ein ruhender und fiel aus der Zaehlung
outcome-resolves: "DoD 2-4: vierter Zustand im Worker-Tab, im workers-Tool, Release-Notiz"
review-summary: "none"
review-gaps: "none: one-line type fix, nothing to remove"
test-verdict: "pass: claude plugin test 522 pass/0 fail (RC=0), go test core/mod core/release RC=0, shipped-history --check ok; DoD 2-4 in tree verified (overview.ts:154, overview-pane.tsx:132, sessions.ts:151, notes/RJFA19.md) and by 4 new tests; live: hc get liefert status.state (working in this tab), 'blocked' is in hc's vocabulary and herdr-hc.sh passes it through tab list unchanged - a real approval dialog was not provoked (would need a refused permission prompt); failure path: herdr tab list failing leaves tabs [] (unchanged); change writes nothing"
hunt-verdict: "pass: bughunt done at 7584a4bde, 0 blocker; 1 change edge (types/index.d.ts:20 Worker.state ohne 'blocked', note), 1 pre-existing same subject (register.tsx:1790 Give back in blocked tab, note for the person)"
review-verdict: "accept: alle drei DoD-Punkte erfuellt und getestet (claude plugin test 522/522); ein fix-now am Typvertrag, kein Blocker"
review-round: 1
review-verdict-gaps: |-
  note: edge: fix-now: core/mod/builtin/jaira-cockpit/types/index.d.ts:20 - der Worker-Typ, der laut Kommentar hooks/overview.ts spiegelt, erlaubt state nur working|idle|question; 'blocked' ergaenzen (eine Zeile, kein neues Verhalten)
  note: elsewhere: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1790 - Give back behandelt nur 'working' als Agent im Zug; bei einem blocked Worker-Tab mit warmem Cache tippt es backMessage plus Enter in den offenen Freigabe-Dialog, bei kaltem schliesst es den Tab mitten im Zug; vor diesem Ticket da (05a8d99fc, QTQW15), braucht eine eigene Entscheidung - neues Ticket
review-check: "1. cd core/mod/builtin/jaira-cockpit && claude plugin test  2. erwartet: 522 pass, 0 fail, darunter 'an agent held at an approval dialog waits for approval', 'a worker held at an approval dialog says so instead of a next round', 'an agent held at an approval dialog is still in its turn and counts'  3. live (Herdr noetig): einen Worker per spawn.sh starten und ihn einen Bash-Befehl ausfuehren lassen, der eine Freigabe verlangt; Dialog NICHT beantworten  4. im Cockpit den Overview-Tab Workers oeffnen: die Figur des Tickets traegt ein gelbes '!', die Zeile sagt 'waits for approval'  5. /cockpit workers eingeben: die Zeile des Tickets endet auf 'waits for approval: answer the dialog in its tab'  6. Dialog beantworten: '!' wird zum Spinner, Text 'working'"
round: 2
back-from: review
back-by: lane
back-count: 1
stop-after-fix: "true"
---

# Das Cockpit bekommt Tabs, darunter ein Worker-Panel: je Agent ein Icon, sein Ticket, seine Lane und ob er festhaengt

## Definition of Done

- [-] Brainstorm hat Tabs, Icon-Zustaende und die Datenquelle fuer 'blockiert/braucht Hilfe' festgelegt
- [x] Worker-Tab im Overview zeigt einen vierten Zustand 'wartet auf Freigabe' (eigenes Zeichen, z.B. '!' in warning-Farbe), wenn Herdr fuer den Tab des Agenten agent_status=blocked meldet (offener Approval-Dialog); Test in tests/
  proof: core/mod/builtin/jaira-cockpit/hooks/overview.ts:154, overview-pane.tsx:132; test 'an agent held at an approval dialog waits for approval, its mark a !'
- [x] Release-Notiz unter core/release/notes/
  proof: core/release/notes/RJFA19.md
- [x] Auch das workers-Tool (mcp__jaira-cockpit__workers, /cockpit workers) nennt den Zustand 'wartet auf Freigabe' fuer diesen Agenten
  proof: core/mod/builtin/jaira-cockpit/hooks/sessions.ts:151; test 'a worker held at an approval dialog says so instead of a next round'

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-03 17:18 · berkc** — Gepruefft 03.11.2026: der Worker-Tab existiert schon (F80QZT, Alexander Sacharov, 16.10.): Overview hat Tabs queue/log/workers, je Agent eine Figur mit Ticket, Lane, naechster Lane, Fortschritt; Zustaende working/idle/question (core/mod/builtin/jaira-cockpit/hooks/overview.ts:122-176). Es fehlt nur: Zustand 'wartet auf Freigabe' (z.B. Bash-Befehl). Ticket auf diesen Rest verkleinern.
- **2026-11-23 19:40 · Alexander Sacharov** — Alexander Sacharov took this ticket over from berkc
- **2026-11-23 19:40 · Alexander Sacharov** — decision: Alex 23.11.2026: Ticket eingeengt auf den Rest, den berkc am 03.11. benannt hat. Worker-Tab mit Ticket/Lane/working/idle/question existiert (F80QZT, overview.ts:121-176). Neu nur: Zustand 'wartet auf Freigabe'. Datenquelle: Herdr agent_status=blocked je Tab (spawn.sh liest denselben Wert, overview.ts:161 liest agent_status schon, kennt blocked aber nicht). DoD 1 (Brainstorm) dadurch erledigt. berkc: vom Teamlead auf Alex' Wunsch uebernommen.
- **2026-11-23 19:41 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/48662.sock
- **2026-11-23 19:41 · Alexander Sacharov** — decision (worker): keine offenen Entscheidungen. Zustand heisst 'blocked', Text 'waits for approval', Zeichen '!' in warning; blocked schlaegt working und question (ein offener Dialog haelt alles an). Quelle nur Herdr-Tab agent_status=blocked; Subagenten liefern keinen solchen Status (WORKING=pending/running/waiting). workers-Tool: Carrier-Tab blocked -> 'waits for approval: answer the dialog in its tab' statt next round.
- **2026-11-23 19:43 · Alexander Sacharov** — Geschwister mitgezogen: Kopfzeile zaehlt blocked-Tabs als laufend (board.ts IN_TURN), Discuss/Explain-Tab wird bei blocked nicht abgeraeumt (register.tsx:441). Offen, nicht hier: Give back (register.tsx:1790) tippt in einen blocked Worker-Tab und beantwortet damit evtl. den Dialog - braucht eigene Entscheidung. Subagenten haben keinen Freigabe-Status in der API; workers-Tool listet nur worker/dispatch-Traeger, keine Lane-Agenten.
- **2026-11-23 19:44 · Alexander Sacharov** — critique: none. Stehen gelassen: cache-Angabe im workers-Tool fuer blocked bleibt die aus lastAnswerAt (ein Agent am Dialog macht keinen API-Call, 'warm (working)' waere falsch); Strings in sessions.ts inline wie die bestehenden.
- **2026-11-23 19:45 · Alexander Sacharov** — waiting: bughunt (background subagent) since 20:45
- **2026-11-23 19:47 · Alexander Sacharov** — bughunt: pre-existing: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1790 - class: Give back treats only Herdr 'working' as an agent in its turn, so a tab held at an approval dialog ('blocked') is handled like a resting worker; this change put IN_TURN on board.ts:175 and register.tsx:441 but not on this one; example: give back a ticket whose <id>/worker tab is 'blocked' - with a warm cache the cockpit types backMessage plus Enter into that pane, so the Enter lands on the approval dialog the person has not read; with a cold cache it runs 'herdr tab close' on the tab mid-turn and spawns a new worker, which the comment at 1786-1789 says must never happen; blame: 05a8d99fc (QTQW15, before the ticket); read
- **2026-11-23 19:47 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/types/index.d.ts:20 - class: the plugin's public state contract no longer mirrors hooks/overview.ts; WorkerState gained 'blocked' (overview.ts:122) but the Worker type in index.d.ts, which says it mirrors it, still allows only 'working' | 'idle' | 'question', and PluginState.workers is typed by it; example: a consumer that switches over Worker.state by this contract never handles 'blocked'; nothing typechecks it (claude plugin test passes 522/522); not a blocker; read
- **2026-11-23 19:47 · Alexander Sacharov** — hunt report: 2 findings (1 change edge, 1 pre-existing), notes on RJFA19 | first hunt: read the cockpit hooks unit whole where it touches agent status (board.ts, overview.ts, overview-pane.tsx, sessions.ts, texts.ts, register.tsx at every agent_status/state use), plus types/index.d.ts; claude plugin test: 522 pass | could not run Herdr here, so what Enter does on a live approval dialog in the Give-back finding is read from code, not reproduced
- **2026-11-23 19:47 · Alexander Sacharov** — bughunt: done at 7584a4bde - 2 findings (1 change, 1 pre-existing) in 7 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx core/mod/builtin/jaira-cockpit/hooks/overview.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/sessions.ts core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/RJFA19.md
- **2026-11-23 19:47 · Alexander Sacharov** — waiting: review (fresh agent via run-lane.sh) since 20:47
- **2026-11-23 19:49 · Alexander Sacharov** — review: done at 7584a4bde
- **2026-11-23 19:49 · Alexander Sacharov** — review: fix-now: types/index.d.ts:20 Worker.state um 'blocked' ergaenzen (spiegelt hooks/overview.ts:122)
- **2026-11-23 19:50 · Alexander Sacharov** — fix-now geschlossen: types/index.d.ts:20 Worker.state um 'blocked' ergaenzt (1 Zeile). stop-after-fix=true von mir gesetzt (Fix < 50 Zeilen, kein Lock/Parser). 'elsewhere'-Notiz (Give back register.tsx:1790) bleibt Zeile fuer die Person.
