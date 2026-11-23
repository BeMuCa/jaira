---
id: 01M5ZJ4WRN0FRAKQCCB805CSWX
title: "spawn.sh startet keinen Worker, wenn der freie Speicher nicht fuer einen weiteren Build reicht"
status: in-progress
ready: true
creator: BeMuCa
goal: "Vor dem Start eines Workers liest spawn.sh den freien Speicher und lehnt ab oder warnt, wenn ein weiterer paralleler Build den Rechner zum Absturz braechte; die Grenze steht an einer Stelle und laesst sich setzen."
context: |-
  Was schief ist: drei parallele Go-Builds auf einer WSL mit 7 GB fuehrten zu vier Abstuerzen und einmal zu 29 leeren Git-Objekten (Teamlead-Session 27.10.2026, Feedback L317). Nichts in jaira warnt vor Speicher.
  Heute: core/role/builtin/jaira-dispatcher/scripts/spawn.sh startet Worker ohne Blick auf die Maschine (grep MemAvailable findet nichts).
  Bekannt: /proc/meminfo MemAvailable reicht unter Linux/WSL; ein Go-Build mit -race braucht grob 1-2 GB.
  Nicht gewollt: ein Scheduler. Eine Zahl (freier Speicher in MB je Worker, Default) und ein Abbruch mit klarer Meldung reichen; der Teamlead entscheidet dann, ob er wartet.
definition-of-done: "spawn.sh lehnt mit Exit 3 und Meldung ab, wenn MemAvailable unter der Grenze liegt; --force startet trotzdem"
tags:
  - concurrency
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:52:06Z
updated-at: 2026-11-23T01:35:32Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: worker-05cswx
claimed-at: 2026-11-23T01:32:21Z
---

# spawn.sh startet keinen Worker, wenn der freie Speicher nicht fuer einen weiteren Build reicht

## Definition of Done

- [x] spawn.sh lehnt mit Exit 3 und Meldung ab, wenn MemAvailable unter der Grenze liegt; --force startet trotzdem
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:135 exit 3, core/role/builtin/jaira-dispatcher/scripts/spawn.sh:110 --force; extracted check against fake meminfo 976 MB -> rc=3, with force=1 -> started
- [x] Die Grenze ist eine Variable am Kopf von spawn.sh mit Kommentar, warum der Default so ist
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:7-12 min_free_mb with comment (1-2 GB per -race build, 2026-10-27 crashes), JAIRA_MIN_FREE_MB
- [x] jaira-teamlead/SKILL.md sagt, was der Teamlead tut, wenn spawn.sh ablehnt
  proof: core/role/builtin/jaira-teamlead/SKILL.md:237 wait, tell the person, --force only on their word
- [x] Ohne /proc/meminfo oder ohne MemAvailable (macOS, Git Bash) ueberspringt spawn.sh die Pruefung mit einer Zeile auf stderr und startet; nie ein Abbruch
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:130; extracted check with missing file and with no MemAvailable line -> one stderr line, rc=0
- [ ] seen by a person: spawn.sh auf macOS oder unter Git Bash gestartet zeigt die Zeile, dass der Speicher nicht geprueft wurde, und oeffnet den Tab

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 01:33 · BeMuCa** — Entscheidungen (keine offen fuer die Person): Grenze min_free_mb am Kopf von spawn.sh, setzbar mit JAIRA_MIN_FREE_MB, Default 2048 MB (ein -race-Build 1-2 GB). 0 schaltet ab, Nicht-Zahl -> Exit 2. Exit 3 bei zu wenig MemAvailable, vor Worktree und Tab. --force startet trotzdem. Ohne /proc/meminfo oder ohne MemAvailable-Zeile (macOS, Git Bash, Kernel <3.14): eine Zeile auf stderr, Start geht weiter. Ausgenommen sind die Lanes explain und discuss: das sind Chats fuer eine Person aus den Cockpit-Knoepfen, kein Build. Bekannte Grenze: gemessen wird der Stand jetzt; zwei Spawns direkt hintereinander bestehen beide, weil der erste noch nicht baut - kein Scheduler, wie der Kontext verlangt. TestMain in core/role setzt JAIRA_MIN_FREE_MB=0, damit die Skript-Tests auf knappen Maschinen nicht an der Pruefung scheitern.
