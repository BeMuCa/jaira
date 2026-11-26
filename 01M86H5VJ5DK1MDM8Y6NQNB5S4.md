---
id: 01M86H5VJ5DK1MDM8Y6NQNB5S4
title: "herdr-hc.sh 'hc get' has no timeout of its own - a hanging hc get hangs…"
status: optimize
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
updated-at: 2026-11-26T08:02:16Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-36559
claimed-at: 2026-11-26T07:54:46Z
goal: "herdr-hc.sh 'tab list' endet in begrenzter Zeit, auch wenn hc get (oder hc ls) haengt: ein Terminal, dessen hc get nicht in der Frist antwortet, zaehlt als stumm (Exit 3 + stderr, wie seit 8NWE4P)"
outcome-what: "herdr-hc.sh: bounded() gibt hc ls und jedem hc get in 'tab list' (und dem Space-Lookup in 'tab create') JAIRA_HC_TIMEOUT Sekunden (3), dann kill; Ausgabe ueber Dateien. Test um haengende get/ls erweitert, Manifest und Release-Note."
outcome-why: "Ein Terminal, dessen hc get nie antwortete, liess 'tab list' und damit reap.sh (unter reap.lock), Give back und die Worker-Liste des Cockpits haengen."
outcome-resolves: "DoD 1-4: haengendes get -> stumm, Exit 3; haengendes ls -> Exit !=0; reines bash ohne Prozessreste; Release-Note."
review-summary: none
---

# herdr-hc.sh 'hc get' has no timeout of its own - a hanging hc get hangs…

## Definition of Done

- [x] herdr-hc.sh 'tab list' bricht jedes hc get nach einer Frist ab (Standard 3 s, JAIRA_HC_TIMEOUT); ein haengendes Terminal steht danach als stumm auf stderr und die Liste endet mit Exit 3; Test mit einem hc-Stub, dessen get schlaeft
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:25-50 (bounded), :123; TestHerdrHcTabListSaysWhichTerminalDidNotAnswer/one_hangs, all_hang, hangs_and_mute (core/role/role_test.go)
- [x] Ein haengendes hc ls laesst 'tab list' nach derselben Frist mit Exit ungleich 0 enden; Test
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:118; TestHerdrHcTabListSaysWhichTerminalDidNotAnswer/ls_hangs
- [x] Ohne timeout(1): die Frist laeuft in reinem bash (bash 3.2 tauglich), ohne liegengebliebene Prozesse
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:34-50 sleep+kill watchdog, kein timeout(1); Stub ohne exec in role_test.go: Test misst <10 s, also haelt nichts die Pipes
- [x] Release-Note in core/release/notes/QNB5S4.md
  proof: core/release/notes/QNB5S4.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-26 07:54 · Sascha** — parent session: uds:/run/user/1000/cc-socks/62036.sock
- **2026-11-26 07:55 · Sascha** — decision (worker): Umfang nur 'tab list' (hc ls + hc get darin). 'tab get'/'pane get' bleiben ohne Frist: reap.sh:417 liest ein gescheitertes pane get als 'Tab zu' und raeumte dann den Worktree eines offenen Tabs ab - ein Haenger dort ist sicherer als ein Timeout. Frist 3 s: Cockpit gibt 'tab list' 15 s gesamt (register.tsx:212), so passen mehrere stumme Terminals. Watchdog in bash statt timeout(1): macOS hat keins, Git Bash findet timeout.exe (siehe reap.sh:389).
- **2026-11-26 08:01 · Sascha** — Dead end: hc-Ausgabe direkt in die Pipe des Aufrufers - ein Kind des getoeteten hc (Stub ohne exec: sleep) hielt stderr offen, Go/Cockpit warteten trotzdem 20 s. Darum schreibt bounded() stdout UND stderr in Dateien und gibt sie danach aus. Vier Kategorien: Nebenlaeufigkeit - jeder Aufruf eigene mktemp-Dateien; PID-Reuse ausgeschlossen, solange hc ungewartet Zombie ist, Restfenster zwischen wait und kill des Watchdogs ist Mikrosekunden. Fehlerpfade - Absage von hc behaelt Exit und stderr (not_found-Erkennung unveraendert), Ablauf -> Exit 137 + 'gave no answer within Ns'. Eingabe - JAIRA_HC_TIMEOUT nicht numerisch/0 -> sleep endet sofort, alles gilt als stumm, Exit 3 = sicherer Fehler. Rechte - nichts Neues. Manifest core/shipped/shipped.sha256 neu geschrieben (TestEveryInstalledFileIsAShippedVersion).
- **2026-11-26 08:02 · Sascha** — critique: nichts. Ein Muster neben reap.sh:389 (timeout(1)-Probe) ist hier gewollt: reap.sh braucht den Gruppen-Kill fuer docker compose und verweigert ohne coreutils; hier hiesse Verweigern, dass 'tab list' auf einem Mac ohne coreutils immer scheitert. JAIRA_HC_TIMEOUT hat einen Aufrufer (der Test) und steht in der Release-Note. Ort stimmt: im Aufgerufenen, kein Cockpit-Code.
