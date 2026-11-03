---
id: 01M67RXNWJMJRW2EG721MG80DF
title: "Das Cockpit zeigt nur die Tickets, die mir gehoeren"
status: todo
ready: true
creator: berkc
assignee: berkc
goal: "Band und Zaehler des jaira-cockpit zeigen nur Tickets, deren assignee die Person ist, die die Session fuehrt (git user.name plus Aliase aus ~/.jaira/identity); fremde Tickets erscheinen nicht."
context: |-
  Was falsch ist: das Cockpit-Band (core/mod/builtin/jaira-cockpit/hooks/board.ts) listet alle wartenden Tickets des Boards, auch die von Kollegen.
  Ausloeser: Maintainer 31.10.2026: 'Erstens will ich nur die tickets von mir und nicht von anderen angezeigt kriegen.'
  Bekannt: Identitaet = git user.name plus Aliase in ~/.jaira/identity. Tickets ohne assignee: im Plan entscheiden (Vorschlag: zeigen, sie gehoeren noch niemandem).
definition-of-done: Ein wartendes Ticket mit fremdem assignee erscheint nicht im Band; Test in hooks
tags:
  - tui
blocked-by: []
related:
  - 01M5M7JTVR8G7DZ1MGHJQ63EEV
commits: []
created-at: 2026-10-31T01:24:25Z
updated-at: 2026-11-03T17:21:52Z
updated-by: berkc
claimed-by: berk-ca-130899
claimed-at: 2026-11-03T17:20:49Z
---

# Das Cockpit zeigt nur die Tickets, die mir gehoeren

## Definition of Done

- [ ] Ein wartendes Ticket mit fremdem assignee erscheint nicht im Band; Test in hooks
- [ ] Ein Ticket mit eigenem assignee oder einem Alias erscheint; Test
- [ ] Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-03 17:21 · berkc** — Before plan (worker 2026-11-03): most of this already exists since NZZDRH (4b2c3701, 2026-10-12): board.ts snapshotOf splits waiting (assignee = whoami identity/alias, or empty) from theirs; counter 'N waiting for you' counts only own+unassigned; teammates' tickets sit in ONE collapsed band line 'N waiting for Berk, Carla' (key 'theirs', register.tsx ~1340) that opens to rows with Discuss only. Tests: board.test.ts ~1044-1165. whoami on this machine resolves berkc + aliases berk.calabakan@gmail.com, BeMuCa; all waiting tickets on this board today are berkc/BeMuCa. Unassigned stays shown (settled: belongs to nobody yet, as context proposes). OPEN for the person: what happens to the collapsed theirs line? 1) (Recommended) remove it and the teammate Discuss from the band entirely - literal DoD, smallest; 2) keep it, hidden by default behind a settings switch; 3) move teammates' waiting tickets into the Overview pane instead of the band; 4) nothing to build - the trigger was BeMuCa missing from ~/.jaira/identity, close as already done by NZZDRH. Stopped before pre-process until answered.
