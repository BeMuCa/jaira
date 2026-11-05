---
id: 01M6NHNRNKP678V30QDJAN1XR5
title: jaira plugin ls/install/rm holt Plugins aus dem Katalog panda-bobr/jaira-releases
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Eine Person installiert, aktualisiert und entfernt ein Plugin aus dem Katalog mit einem Befehl, und ein Binary mit falscher Pruefsumme wird nicht installiert."
context: |-
  Heute gibt es keinen Befehl dafuer. Katalog: plugins.json in https://github.com/panda-bobr/jaira-releases (Format: R6X6FH).
  Ablauf wie Obsidian: Eintrag -> repo -> jaira-plugin.json vom Default-Branch -> Release v<version> -> Asset fuer os/arch + checksums.txt -> sha256 pruefen -> ~/.jaira/plugins/<id>/.
  Vorbild im Code: core/market/market.go liest lanes/ ueber die GitHub-API, gleiches Muster.
  minJairaVersion hoeher als die laufende Version: verweigern mit Hinweis auf jaira self upgrade.
definition-of-done: jaira plugin ls zeigt den Katalog mit installierter und neuester Version; Test gegen einen Fake-Server
tags:
  - cli
blocked-by:
  - R6X6FH
related: []
commits: []
created-at: 2026-11-05T09:47:08Z
updated-at: 2026-11-05T09:47:08Z
---

# jaira plugin ls/install/rm holt Plugins aus dem Katalog panda-bobr/jaira-releases

## Definition of Done

- [ ] jaira plugin ls zeigt den Katalog mit installierter und neuester Version; Test gegen einen Fake-Server
- [ ] jaira plugin install <id> prueft sha256 und verweigert bei Abweichung (exit 3); Test
- [ ] jaira plugin rm <id> entfernt das Plugin; Test
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

