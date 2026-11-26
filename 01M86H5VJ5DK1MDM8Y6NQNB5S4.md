---
id: 01M86H5VJ5DK1MDM8Y6NQNB5S4
title: "herdr-hc.sh 'hc get' has no timeout of its own - a hanging hc get hangs…"
status: in-progress
ready: true
creator: Sascha
context: |-
  herdr-hc.sh 'tab list' (core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:82) ruft hc ls und je Terminal hc get ohne Frist auf.
  Haengt eins davon, haengt 'tab list' - und damit jeder Aufrufer.
  reap.sh:133 haengt dann unter reap.lock: jede spaetere Runde wird uebersprungen, nichts wird mehr abgeraeumt.
  Das Cockpit (register.tsx:212) bricht nach 15 s ab und sieht dann gar keine Tabs, nicht nur das eine stumme Terminal.
  Gefunden bei 8NWE4P (Exit 3 fuer stumme Terminals); dort blieb der Haenger selbst offen.
tags: []
blocked-by: []
related: []
follows: 01M866HXTHQS29K2ZMNQ8NWE4P
commits: []
created-at: 2026-11-24T10:21:14Z
updated-at: 2026-11-26T08:01:12Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-36559
claimed-at: 2026-11-26T07:54:46Z
goal: "herdr-hc.sh 'tab list' endet in begrenzter Zeit, auch wenn hc get (oder hc ls) haengt: ein Terminal, dessen hc get nicht in der Frist antwortet, zaehlt als stumm (Exit 3 + stderr, wie seit 8NWE4P)"
---

# herdr-hc.sh 'hc get' has no timeout of its own - a hanging hc get hangs…

## Definition of Done

- [x] herdr-hc.sh 'tab list' bricht jedes hc get nach einer Frist ab (Standard 3 s, JAIRA_HC_TIMEOUT); ein haengendes Terminal steht danach als stumm auf stderr und die Liste endet mit Exit 3; Test mit einem hc-Stub, dessen get schlaeft
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:25-50 (bounded), :123; TestHerdrHcTabListSaysWhichTerminalDidNotAnswer/one_hangs, all_hang, hangs_and_mute (core/role/role_test.go)
- [x] Ein haengendes hc ls laesst 'tab list' nach derselben Frist mit Exit ungleich 0 enden; Test
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:118; TestHerdrHcTabListSaysWhichTerminalDidNotAnswer/ls_hangs
- [ ] Ohne timeout(1): die Frist laeuft in reinem bash (bash 3.2 tauglich), ohne liegengebliebene Prozesse
- [ ] Release-Note in core/release/notes/QNB5S4.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-26 07:54 · Sascha** — parent session: uds:/run/user/1000/cc-socks/62036.sock
- **2026-11-26 07:55 · Sascha** — decision (worker): Umfang nur 'tab list' (hc ls + hc get darin). 'tab get'/'pane get' bleiben ohne Frist: reap.sh:417 liest ein gescheitertes pane get als 'Tab zu' und raeumte dann den Worktree eines offenen Tabs ab - ein Haenger dort ist sicherer als ein Timeout. Frist 3 s: Cockpit gibt 'tab list' 15 s gesamt (register.tsx:212), so passen mehrere stumme Terminals. Watchdog in bash statt timeout(1): macOS hat keins, Git Bash findet timeout.exe (siehe reap.sh:389).
