---
id: 01M59MQ7PT0S3NYN5J7AWFZ4EY
title: Neue Lanes starten aus einer Vorlage pro Lane-Art statt aus einer Datei mit allen Feldern
status: ideas
ready: true
creator: berkc
assignee: berkc
goal: "Wer eine Lane anlegt, waehlt zuerst die Art (Agent-Schritt, Skill ausfuehren, Uebergabe an einen Menschen, Endstation) und bekommt nur die Felder, die diese Art braucht, vorbelegt."
context: |-
  Heute: n im Lanes-Screen schreibt writeLaneSkeleton (internal/tui/lanes.go:578), eine Vorlage mit 20+ Feldern und Kommentaren. berkc hat damit eine kaputte Lane gebaut (siehe FGRQCG).
  Wunsch von berkc am 19.10.2026, zwei Beispiele:
  - Skill-Lane: nur Name, Beschreibung und welcher Skill laeuft. Haengt an XR27FN (Lane, die nur einen Skill ausfuehrt) - das Feld dafuer gibt es erst mit XR27FN.
  - Agent->Mensch-Uebergabe: requires-question: true, agentic: false, und im Prompt die Zeile, dass der Mensch benachrichtigt wird, sobald ein Ticket hier landet (siehe Ticket zur Benachrichtigung bei HITL-Lanes).
  Weitere naheliegende Arten: normaler Agent-Schritt (model-tier, input/output), Review mit frischem Kontext (context: fresh), Endstation (terminal).
  Die Vorlagen gehoeren ins Binary (wie die ausgelieferten Lanes), nicht in jedes Repo.
definition-of-done: n im Lanes-Screen fragt zuerst nach der Art und oeffnet eine Vorlage nur mit den Feldern dieser Art; Test pro Art
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-19T08:33:47Z
updated-at: 2026-10-22T00:15:37Z
updated-by: berkc
round: 2
---

# Neue Lanes starten aus einer Vorlage pro Lane-Art statt aus einer Datei mit allen Feldern

## Definition of Done

- [ ] n im Lanes-Screen fragt zuerst nach der Art und oeffnet eine Vorlage nur mit den Feldern dieser Art; Test pro Art
- [ ] Die Uebergabe-an-Mensch-Vorlage setzt requires-question und enthaelt die Benachrichtigungs-Zeile; Test
- [ ] jaira lanes new --kind <art> (oder gleichwertig) macht dasselbe auf der CLI; Test
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-22 00:15 · berkc** — berkc am 22.10.2026: gehoert in die Lane Feature Ideas, erst durchdenken (siehe W0JX10).
