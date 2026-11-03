---
id: 01M67RXNWJMJRW2EG721MG80DF
title: "Das Cockpit zeigt nur die Tickets, die mir gehoeren"
status: in-progress
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
updated-at: 2026-11-03T22:49:27Z
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

- [x] plugin.json (core/mod/builtin/jaira-cockpit/.claude-plugin/plugin.json): add userConfig key showTeammates after workerContextShare: {"type": "boolean", "title": "Teammates' waiting tickets", "description": "Show the band's line with the tickets waiting for your teammates. Off: the band shows only your own tickets and nobody's.", "default": false}. Keep version 0.2.1.
- [x] register.tsx (hooks/register.tsx): in register(), after the line 'const share = ...workerContextShare...', add a comment line '// Whether the band shows the line with the teammates' waiting tickets; off unless set.' and 'const showTeammates = options.showTeammates === true'. In the band render handler replace 'const theirs = view?.theirs ?? []' with 'const theirs = showTeammates ? view?.theirs ?? [] : []'. Touch nothing else (board.ts unchanged).
- [x] tests/board.test.ts: give the four tests that use the teammates' line the option, as the 'the share is a setting' test does - test(name, { options: { showTeammates: true } }, async ($, on) => ...): 'only the person's own and nobody's tickets get buttons and count; ...' (~l.1044), 'a Discuss chat on a teammate's ticket stays open ...' (~l.1072), 'the band reads top down ...' (~l.1121), 'opened by +N more, ...' (~l.1149).
- [x] tests/board.test.ts: new test right after 'without a word from jaira whoami ...' (~l.1090): 'by default a teammate's waiting ticket is not in the band; the person's own, an alias's and nobody's are'. w.me = { identity: 'Alexander Sacharov', aliases: ['Alexander Sacharov', 'alex@example.org'] }; w.tickets = [{ ...SIGNOFF, assignee: 'Alexander Sacharov' }, { ...SIGNOFF, handle: 'AB12CD', assignee: 'alex@EXAMPLE.org' }, { ...SIGNOFF, handle: 'CD34EF', assignee: '' }, { ...SIGNOFF, handle: 'EF34GH', assignee: 'Berk' }, WORKING]; startSession; headLine = '0 worker | 0 review | 3 waiting for you | 5h 34%'; band has accept-MG677R, accept-AB12CD, accept-CD34EF; find({ key: 'theirs' }) undefined; await drawn($) contains neither 'EF34GH' nor 'Berk'.
- [x] run: cd core/mod/builtin/jaira-cockpit && claude plugin test . -> 176 pass, 0 fail (check exit code 0).
- [~] release note: create core/release/notes/MG80DF.md with exactly one line: '- ! Expect your teammates' waiting tickets to be gone from the cockpit band; turn on the jaira-cockpit setting showTeammates to get their line back.'
- [ ] shipped hashes: from repo root run 'sha256sum core/mod/builtin/jaira-cockpit/.claude-plugin/plugin.json core/mod/builtin/jaira-cockpit/hooks/register.tsx >> core/shipped/shipped.sha256' (append, never delete a line; tests/ are not shipped). Then: export PATH=$HOME/.local/go/bin:$PATH; go test ./core/shipped ./core/mod -count=1 -> ok, exit 0.
- [ ] tick DoD: 1 and 2 with --proof naming the new test, 3 with --proof core/release/notes/MG80DF.md; then jaira move MG80DF --to critique --what --why --resolves; then commit code + ticket together: git add the 5 changed files plus .jaira/tickets/01M67RXNWJMJRW2EG721MG80DF-*.md; git commit -m 'feat(MG80DF): the cockpit band hides teammates' waiting tickets unless showTeammates is on'.

## Progress
- **2026-11-03 17:21 · berkc** — Before plan (worker 2026-11-03): most of this already exists since NZZDRH (4b2c3701, 2026-10-12): board.ts snapshotOf splits waiting (assignee = whoami identity/alias, or empty) from theirs; counter 'N waiting for you' counts only own+unassigned; teammates' tickets sit in ONE collapsed band line 'N waiting for Berk, Carla' (key 'theirs', register.tsx ~1340) that opens to rows with Discuss only. Tests: board.test.ts ~1044-1165. whoami on this machine resolves berkc + aliases berk.calabakan@gmail.com, BeMuCa; all waiting tickets on this board today are berkc/BeMuCa. Unassigned stays shown (settled: belongs to nobody yet, as context proposes). OPEN for the person: what happens to the collapsed theirs line? 1) (Recommended) remove it and the teammate Discuss from the band entirely - literal DoD, smallest; 2) keep it, hidden by default behind a settings switch; 3) move teammates' waiting tickets into the Overview pane instead of the band; 4) nothing to build - the trigger was BeMuCa missing from ~/.jaira/identity, close as already done by NZZDRH. Stopped before pre-process until answered.
- **2026-11-03 22:45 · berkc** — Maintainer 03.11.2026 (teamlead): Teamkollegen-Zeile bleibt, aber standardmaessig versteckt; ein Settings-Schalter blendet sie ein (Option 2). Unzugewiesene bleiben sichtbar.
- **2026-11-03 22:47 · berkc** — Plan reasoning: 'settings switch' = a plugin userConfig boolean showTeammates in core/mod/builtin/jaira-cockpit/.claude-plugin/plugin.json, read in register(on, options) like cacheGuardTokens/workerContextShare - not ~/.jaira/settings.json, which holds jaira's own dispatch choice. Filter at render only (register.tsx band handler, 'const theirs = view?.theirs ?? []'); board.ts snapshotOf stays, its split and counter already do DoD 2. No version bump: QRHB35 (98718866) added workerContextShare without one; release role here is 'worker', repo has no per-change bump rule. Shipped hashes: append sha256sum lines; do NOT rerun scripts/shipped-history.sh - it also adds 5 unrelated historic lines (40-review.md, board.ts) from other branches. Baseline before change: claude plugin test . = 175 pass, 0 fail.
- **2026-11-03 22:47 · berkc** — Plan step 1: the --plan --done marker was written in the same batch as the plugin.json Edit, and the Edit was refused (file not read yet). Re-did the Edit after a Read; step 1 final state verified (userConfig now has showTeammates, JSON parses).
