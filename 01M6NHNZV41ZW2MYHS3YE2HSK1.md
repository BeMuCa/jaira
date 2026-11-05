---
id: 01M6NHNZV41ZW2MYHS3YE2HSK1
title: "panda-bobr/jaira-plugin-template: lauffaehige Vorlage mit Beispielbefehl und Release-Pipeline"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer 'Use this template' drueckt und einen Tag setzt, hat ein Release, das die Katalog-CI ohne Aenderung annimmt."
context: |-
  Das Repo https://github.com/panda-bobr/jaira-plugin-template ist seit 05.11.2026 leer und als Template markiert.
  Inhalt: Go-Modul, Beispielbefehl 'jaira-hello' (liest 'jaira list --json', schreibt eine note), jaira-plugin.json, goreleaser-Konfig (Assets jaira-<id>_<os>_<arch>, checksums.txt), Action fuer Test und Release auf Tag, README mit den drei Schritten bis zum Katalog-PR.
  Format und Regeln: R6X6FH (docs/PLUGINS.md).
definition-of-done: "Ein Repo aus der Vorlage baut mit 'git tag v0.1.0 && git push --tags' ein Release mit allen Assets und checksums.txt; einmal vorgefuehrt"
tags:
  - ci
blocked-by:
  - R6X6FH
related: []
commits: []
created-at: 2026-11-05T09:47:15Z
updated-at: 2026-11-05T10:23:54Z
updated-by: Alexander Sacharov
---

# panda-bobr/jaira-plugin-template: lauffaehige Vorlage mit Beispielbefehl und Release-Pipeline

## Definition of Done

- [ ] Ein Repo aus der Vorlage baut mit 'git tag v0.1.0 && git push --tags' ein Release mit allen Assets und checksums.txt; einmal vorgefuehrt
- [ ] Das Release besteht die CI von panda-bobr/jaira-releases

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-05 10:23 · Alexander Sacharov** — Teamlead (05.11.2026): der Standard (R6X6FH, abgenommen) liegt noch nicht auf master, sondern auf origin/feat/r6x6fh: docs/PLUGINS.md und docs/schema/*.json - lesen mit 'git show origin/feat/r6x6fh:docs/PLUGINS.md'. Die Arbeit selbst gehoert in das Repo panda-bobr/e2hsk1; Alex ist dort admin. Ein Agent pusht seinen Branch und oeffnet keinen PR.
