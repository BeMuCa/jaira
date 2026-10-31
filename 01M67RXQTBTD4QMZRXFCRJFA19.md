---
id: 01M67RXQTBTD4QMZRXFCRJFA19
title: "Das Cockpit bekommt Tabs, darunter ein Worker-Panel: je Agent ein Icon, sein Ticket, seine Lane und ob er festhaengt"
status: backlog
ready: true
creator: berkc
assignee: berkc
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
updated-at: 2026-10-31T01:24:27Z
---

# Das Cockpit bekommt Tabs, darunter ein Worker-Panel: je Agent ein Icon, sein Ticket, seine Lane und ob er festhaengt

## Definition of Done

- [ ] Brainstorm hat Tabs, Icon-Zustaende und die Datenquelle fuer 'blockiert/braucht Hilfe' festgelegt
- [ ] Tab Worker zeigt je Agent Icon, Ticket und Lane; Test
- [ ] Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

