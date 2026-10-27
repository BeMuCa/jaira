---
id: 01M5ZJ4WRN0FRAKQCCB805CSWX
title: "spawn.sh startet keinen Worker, wenn der freie Speicher nicht fuer einen weiteren Build reicht"
status: backlog
ready: false
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
updated-at: 2026-10-27T20:52:06Z
---

# spawn.sh startet keinen Worker, wenn der freie Speicher nicht fuer einen weiteren Build reicht

## Definition of Done

- [ ] spawn.sh lehnt mit Exit 3 und Meldung ab, wenn MemAvailable unter der Grenze liegt; --force startet trotzdem
- [ ] Die Grenze ist eine Variable am Kopf von spawn.sh mit Kommentar, warum der Default so ist
- [ ] jaira-teamlead/SKILL.md sagt, was der Teamlead tut, wenn spawn.sh ablehnt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

