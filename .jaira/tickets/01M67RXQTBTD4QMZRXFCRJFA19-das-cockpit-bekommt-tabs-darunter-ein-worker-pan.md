---
id: 01M67RXQTBTD4QMZRXFCRJFA19
title: "Das Cockpit bekommt Tabs, darunter ein Worker-Panel: je Agent ein Icon, sein Ticket, seine Lane und ob er festhaengt"
status: critique
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
updated-at: 2026-11-23T19:44:12Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-40650
claimed-at: 2026-11-23T19:41:00Z
outcome-what: "Worker-Tab und workers-Tool zeigen 'waits for approval' ('!' in warning), wenn Herdr fuer den Agenten-Tab agent_status=blocked meldet; Kopfzeile zaehlt solche Agenten weiter"
outcome-why: Ein Agent am Freigabe-Dialog sah aus wie ein ruhender und fiel aus der Zaehlung
outcome-resolves: "DoD 2-4: vierter Zustand im Worker-Tab, im workers-Tool, Release-Notiz"
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
