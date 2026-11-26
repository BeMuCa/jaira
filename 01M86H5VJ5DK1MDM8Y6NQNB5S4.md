---
id: 01M86H5VJ5DK1MDM8Y6NQNB5S4
title: "herdr-hc.sh 'hc get' has no timeout of its own - a hanging hc get hangs…"
status: backlog
ready: true
creator: Sascha
context: "note: elsewhere: herdr-hc.sh 'hc get' has no timeout of its own - a hanging hc get hangs tab list and every caller; before this ticket"
tags: []
blocked-by: []
related: []
follows: 01M866HXTHQS29K2ZMNQ8NWE4P
commits: []
created-at: 2026-11-24T10:21:14Z
updated-at: 2026-11-26T07:55:54Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-36559
claimed-at: 2026-11-26T07:54:46Z
goal: "herdr-hc.sh 'tab list' endet in begrenzter Zeit, auch wenn hc get (oder hc ls) haengt: ein Terminal, dessen hc get nicht in der Frist antwortet, zaehlt als stumm (Exit 3 + stderr, wie seit 8NWE4P)"
---

# herdr-hc.sh 'hc get' has no timeout of its own - a hanging hc get hangs…

## Definition of Done

- [ ] herdr-hc.sh 'tab list' bricht jedes hc get nach einer Frist ab (Standard 3 s, JAIRA_HC_TIMEOUT); ein haengendes Terminal steht danach als stumm auf stderr und die Liste endet mit Exit 3; Test mit einem hc-Stub, dessen get schlaeft
- [ ] Ein haengendes hc ls laesst 'tab list' nach derselben Frist mit Exit ungleich 0 enden; Test
- [ ] Ohne timeout(1): die Frist laeuft in reinem bash (bash 3.2 tauglich), ohne liegengebliebene Prozesse
- [ ] Release-Note in core/release/notes/QNB5S4.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-26 07:54 · Sascha** — parent session: uds:/run/user/1000/cc-socks/62036.sock
