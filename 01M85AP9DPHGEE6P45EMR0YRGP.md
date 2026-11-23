---
id: 01M85AP9DPHGEE6P45EMR0YRGP
title: "Vor dem Uebernehmen zeigt jaira lanes market, was die Lane tut - und ein Agent uebernimmt nicht allein"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Wer eine Lane aus dem Marktplatz holt, sieht vorher Vertrag und Prompt, und ohne eine Person am Terminal wird keine fremde Lane uebernommen."
context: |-
  Heute schreibt 'jaira lanes market adopt <id>' die Datei sofort nach ~/.jaira/lanes (internal/cli/market.go, newLanesMarketAdoptCmd). Vorher zeigt es nichts.
  Die Liste zeigt nur id, name und description. Den Prompt sieht man erst nach dem Uebernehmen mit 'jaira lanes show <id>' oder im Browser.
  Warum das zaehlt: ein Lane-Prompt ist eine Anweisung an einen Agenten mit den Rechten des Nutzers. Eine fremde Lane ist damit so heikel wie fremder Code. Ein Agent, der 'adopt' selbst aufruft, holt sich seine naechsten Anweisungen aus dem Netz.
  Loesung:
  - 'jaira lanes market show <id>' druckt Vertrag und ganzen Prompt, ohne etwas zu schreiben.
  - 'adopt' druckt vor dem Schreiben: id, Quelle, sha256, model-tier, wohin die Lane zurueckschickt, Groesse des Prompts, und jede Adresse (URL), die im Prompt steht.
  - An einem Terminal fragt adopt einmal nach. Ohne Terminal lehnt es mit Exit 3 ab, solange --yes fehlt - dieselbe Form wie 'jaira reject' ohne --force.
  - Die erzeugte Agenten-Notiz (CLAUDE.md/AGENTS.md-Block, 'Task tracking: jaira') sagt: eine Lane aus dem Marktplatz zu uebernehmen entscheidet eine Person; --yes ist nicht Sache des Agenten.
  Ob die Ablehnung ohne Terminal so gewollt ist, steht als Frage im Frage-Ticket dieses Milestones.
  Vorgabe: docs/MARKET.md.
definition-of-done: "jaira lanes market show <id> druckt Vertrag und Prompt und schreibt nichts (Test)"
tags:
  - cli
  - gates
blocked-by:
  - BTHBEX
related: []
commits: []
created-at: 2026-11-23T23:08:38Z
updated-at: 2026-11-23T23:08:38Z
---

# Vor dem Uebernehmen zeigt jaira lanes market, was die Lane tut - und ein Agent uebernimmt nicht allein

## Definition of Done

- [ ] jaira lanes market show <id> druckt Vertrag und Prompt und schreibt nichts (Test)
- [ ] adopt druckt die Zusammenfassung mit sha256 und den URLs aus dem Prompt vor dem Schreiben
- [ ] adopt ohne Terminal und ohne --yes endet mit Exit 3 und schreibt nichts (Test)
- [ ] Die erzeugte Agenten-Notiz nennt die Regel; jaira update --roles schreibt sie nach
- [ ] Release-Note core/release/notes/<TICKET>.md mit '- ! ' fuer Skripte, die adopt aufrufen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

