---
id: 01M5ZJ4SWDVEZ3C6BV0VXZBF0V
title: "Ein Worker schreibt nur in seinem Worktree, startet nichts im Hintergrund und laesst am Lane-Ende nichts liegen"
status: backlog
ready: false
creator: BeMuCa
goal: "Die Rolle jaira-worker verlangt: Schreibpfade nur im eigenen Worktree, docker nur mit eigenem Projektnamen, Laeufe nur im Vordergrund mit Timeout, commit+push nach jedem Lane-Schritt, am Lane-Ende git status leer und origin aktuell, und keine eigene Antwort auf eine Frage in human."
context: |-
  Was schief ist (Teamlead-Session 25.-27.10.2026, 16 Worker, Feedback L310/L316/L325/O4/O16):
  - Worker schrieben ins Hauptverzeichnis statt ins eigene Worktree und stoppten vermutlich dessen laufenden Stack.
  - Worker starteten Hintergrundlaeufe und konnten sie nicht beenden (pkill eigener Prozesse wird vom Classifier abgelehnt).
  - Zwei Worker bauten denselben Stack (gleicher Compose-Projektname) fast gleichzeitig.
  - Ein Driver (ZE9WZC) ging aus human mit der Antwort des Teamlead-Agenten weiter, ein anderer (RNQRJB) lehnte dasselbe ab: nirgends steht, dass nur eine Person in human antwortet.
  Heute: core/role/builtin/jaira-worker/SKILL.md ab Zeile 265 beschreibt den Worktree, aber keine Schreibsperre, keinen Projektnamen, kein Vordergrund-Gebot, keinen Endzustand je Lane.
  Das Muster 'Worktree ab master, neutraler Claim, frischer Subagent je Lane, Stopp in human/signoff' lief fuer 19 Tickets durch (O16) und gehoert in die Rolle.
  Nicht gewollt: ein Hook, der Pfade blockiert - erst die Regel, dann sehen, ob sie reicht.
definition-of-done: jaira-worker/SKILL.md enthaelt die sechs Regeln als Liste mit je einem Satz Grund
tags:
  - docs
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:52:03Z
updated-at: 2026-10-27T20:52:03Z
---

# Ein Worker schreibt nur in seinem Worktree, startet nichts im Hintergrund und laesst am Lane-Ende nichts liegen

## Definition of Done

- [ ] jaira-worker/SKILL.md enthaelt die sechs Regeln als Liste mit je einem Satz Grund
- [ ] Die Regel 'in human antwortet nur eine Person; eine Antwort eines Agenten gilt nicht' steht in jaira-worker und jaira-dispatcher
- [ ] jaira-dispatcher gibt den Compose-Projektnamen und das Timeout im Brief vor
- [ ] Auf einem Scratch-Board nachgestellt: ein Worker beendet eine Lane mit leerem git status und gepushtem Branch

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

