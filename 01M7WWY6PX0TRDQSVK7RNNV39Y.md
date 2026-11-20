---
id: 01M7WWY6PX0TRDQSVK7RNNV39Y
title: Ein privates Board steht in .git/info/exclude statt in der .gitignore des Arbeitsbaums
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: Ein privates Board (vor jaira share) wird ueber .git/info/exclude ausgeschlossen; jaira schreibt fuer ein privates Board nie in die getrackte .gitignore.
context: |-
  Heute: board.AddIgnore (core/board/gitignore.go:44, aufgerufen aus core/board/board.go:22) haengt "# jaira board ... /.jaira/" an die .gitignore des Arbeitsbaums.
  Folge: in jedem neuen Worktree ist .gitignore geaendert. Ein Push-Guard im Projekt lehnt das ab, jeder Push scheitert, der Teamlead raeumt es bei jedem Ticket von Hand weg.
  Gemeldet vom Teamlead eines Projekts mit privatem Board am 20.11.2026; Alex: "das muss unterstuetzt werden, nicht umgangen".
  .git/info/exclude liegt im gemeinsamen git-dir (git rev-parse --git-common-dir) und gilt fuer alle Worktrees, landet nie in einem Diff.
  Offen zu klaeren im Plan: Migration bestehender Boards (Zeile aus .gitignore entfernen, in exclude setzen) und was board.Ignored/jaira share/RemoveIgnore dann lesen.
definition-of-done: "Ein neues privates Board und jeder neue Worktree davon lassen git status sauber; der Ausschluss steht in <git-common-dir>/info/exclude (Test)"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-11-20T16:34:22Z
updated-at: 2026-11-20T16:34:22Z
---

# Ein privates Board steht in .git/info/exclude statt in der .gitignore des Arbeitsbaums

## Definition of Done

- [ ] Ein neues privates Board und jeder neue Worktree davon lassen git status sauber; der Ausschluss steht in <git-common-dir>/info/exclude (Test)
- [ ] jaira share hebt den Ausschluss in info/exclude auf und laesst /.jaira/lanes/ privat wie heute (Test)
- [ ] Ein Board, das die alte .gitignore-Zeile hat, wird bei jaira update umgestellt und die Zeile entfernt (Test)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

