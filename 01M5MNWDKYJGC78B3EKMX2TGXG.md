---
id: 01M5MNWDKYJGC78B3EKMX2TGXG
title: "Ein Worker findet seinen Teamlead nach einem Neustart wieder, und eine Frage an die Person kommt als fertige Zeile"
status: backlog
ready: false
creator: BeMuCa
goal: "Nach einem Crash oder Neustart des Teamleads erreichen die Nachrichten der Worker ihn weiter, ohne dass jemand den neuen Namen von Hand nachreicht; und eine Frage, die nur die Person per Kommando loesen kann, kommt mit der Zeile, die sie tippen muss."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkte 3 und 4).
  Nach einem Neustart bekommt jede Session einen neuen Namen. Worker behalten das beim Start uebergebene --parent und schicken ihre Berichte ins Leere, bis der Teamlead jedem den neuen Namen nachsendet.
  Zweiter Fall: ein Worker brauchte Zugangsdaten in einem Container; Agenten duerfen das zu Recht nicht kopieren. Die Frage ging an die Person, aber ohne die Zeile, die sie haette tippen koennen - eine Runde verloren. Claude Code kennt '! <befehl>' im Prompt der Person.
  Wo das Elternfeld gesetzt wird: spawn.sh (--parent), jaira-worker/SKILL.md und jaira-dispatcher/SKILL.md (Berichtsweg).
definition-of-done: Der Worker schlaegt die aktuelle Adresse seines Teamleads beim Senden nach (Session-Record oder Board) statt den Startwert einzufrieren; ein Test deckt den Namenswechsel ab
tags:
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:25:44Z
updated-at: 2026-10-26T19:07:55Z
assignee: BeMuCa
updated-by: BeMuCa
---

# Ein Worker findet seinen Teamlead nach einem Neustart wieder, und eine Frage an die Person kommt als fertige Zeile

## Definition of Done

- [ ] Der Worker schlaegt die aktuelle Adresse seines Teamleads beim Senden nach (Session-Record oder Board) statt den Startwert einzufrieren; ein Test deckt den Namenswechsel ab
- [ ] jaira-worker/SKILL.md und jaira-dispatcher/SKILL.md sagen: eine Frage, die nur die Person per Kommando loesen kann, enthaelt die fertige '! <befehl>'-Zeile
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

