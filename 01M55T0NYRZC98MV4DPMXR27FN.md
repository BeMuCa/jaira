---
id: 01M55T0NYRZC98MV4DPMXR27FN
title: Eine Lane kann einen Skill starten statt nur einen Prompt zu tragen
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Eine Lane-Datei kann einen Skill nennen (z.B. ponytail zum Aufraeumen von Code als Testing- oder Optimize-Lane); 'jaira show --for-lane' gibt ihn aus und der Agent ruft genau diesen Skill auf."
context: |-
  Ausloeser: Berk am 17.10.: 'Wir sollten eine Lane anlegen koennen, die nur einen Skill ausfuehrt, z.B. ponytail (Skill zum Code-Aufraeumen) als Testing-Lane.'
  Heute geht das nur ueber den Prompt-Text: eine Lane ist .jaira/lanes/<id>.md mit Frontmatter (id, name, after, agentic, model-tier, rejects-to, input-requires, output-produces) und '# Prompt'. Ein Prompt 'Rufe den Skill /ponytail auf den Diff dieses Tickets auf' funktioniert schon, aber jaira weiss nichts davon: es prueft nicht, ob der Skill installiert ist, und der Agent muss es aus Freitext lesen.
  ponytail ist auf Berks Rechner nicht installiert (weder ~/.claude/skills noch Plugin-Cache; nur in einem CHANGELOG eines Plugins erwaehnt) - Name und Quelle klaeren.
  Offen fuer brainstorm: ein Frontmatter-Feld 'skill: <name>' (plus optional args), das 'show --for-lane --json' als eigenes Feld ausgibt; was mit output-produces passiert, wenn ein Skill keine jaira-Outputs kennt (der Prompt muss ihm sagen, welche Felder er zurueckschreibt); Warnung in 'jaira lanes', wenn der Skill fehlt. Kleiner als ein Plugin-System bleiben.
definition-of-done: |-
  Eine Lane-Datei kann 'skill: <name>' tragen, 'jaira lanes' validiert es
  'jaira show <id> --for-lane <lane> --json' gibt den Skill aus, und der Agent-Block sagt, dass er dann aufzurufen ist
  Ein Beispiel-Lane-Template mit Skill ist dokumentiert
  Test; Release-Note
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-17T20:49:21Z
updated-at: 2026-10-17T23:48:25Z
updated-by: berkc
---

# Eine Lane kann einen Skill starten statt nur einen Prompt zu tragen

## Definition of Done

- [ ] Eine Lane-Datei kann 'skill: <name>' tragen, 'jaira lanes' validiert es
'jaira show <id> --for-lane <lane> --json' gibt den Skill aus, und der Agent-Block sagt, dass er dann aufzurufen ist
Ein Beispiel-Lane-Template mit Skill ist dokumentiert
Test; Release-Note

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

