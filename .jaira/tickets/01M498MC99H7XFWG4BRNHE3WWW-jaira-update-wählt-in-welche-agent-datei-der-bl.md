---
id: 01M498MC99H7XFWG4BRNHE3WWW
title: "jaira update wählt, in welche Agent-Datei der Block geht"
status: human
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "jaira update --agent-file agents|claude|both legt fest, ob der jaira-Block in AGENTS.md, CLAUDE.md oder beiden steht, und die Wahl hält ohne neue Konfiguration."
context: |-
  Was falsch ist: jaira schreibt den Block immer in AGENTS.md UND CLAUDE.md (core/board/announce.go:318, agentFiles). Wer nur eine der beiden Dateien will, bekommt die andere bei jedem update und jeder Lane-Änderung zurück.
  Auslöser: Alex will wählen können, weil Claude Code inzwischen AGENTS.md liest.
  Bekannt (code.claude.com/docs/en/memory): Claude Code liest AGENTS.md nur, wenn es KEINE CLAUDE.md gibt. Gibt es eine, zählt nur sie; AGENTS.md kommt dann nur über eine Zeile '@AGENTS.md' in CLAUDE.md hinein.
  Folge für 'agents': Block aus CLAUDE.md entfernen; bleibt CLAUDE.md leer, wird sie gelöscht (dann liest Claude AGENTS.md); bleibt eigener Inhalt, kommt '@AGENTS.md' hinein, sonst sieht Claude den Block nicht.
  Folge für 'claude': Block aus AGENTS.md entfernen, leere AGENTS.md löschen.
  Wo die Wahl lebt: in den Dateien selbst — geschrieben wird in die Dateien, die den Block schon tragen; trägt keine ihn, in beide (heutiges Verhalten). Kein neues Feld in settings.json und keine Datei in .jaira/.
  Aufrufer, die mitziehen müssen: internal/cli/update.go:86, internal/cli/lanes.go:471, internal/tui/lanes.go:142, core/board/board.go:23 (init), NoteIsCurrent für validate.
definition-of-done: "jaira update --agent-file agents schreibt den Block nur nach AGENTS.md, nimmt ihn aus CLAUDE.md und sorgt dafür, dass Claude Code ihn trotzdem liest (leere CLAUDE.md gelöscht, sonst @AGENTS.md darin)"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-06T18:46:48Z
updated-at: 2026-10-06T18:57:17Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3243
claimed-at: 2026-10-06T18:46:59Z
outcome-what: "jaira update --agent-file agents|claude|both; core/board ChooseAgentFiles, removeBlock, setClaudeImport; AnnounceInAgentFiles schreibt nur in Dateien, die den Block tragen."
outcome-why: "Wer nur AGENTS.md oder nur CLAUDE.md will, bekam die andere Datei bei jedem update zurück; Claude Code liest AGENTS.md nur ohne CLAUDE.md."
outcome-resolves: "Die Wahl ist ein Flag, hält ohne Konfiguration, und Claude Code sieht den Block in jeder Variante."
review-summary: "none — die Wahl lebt in den Dateien, Aufrufer bleiben unverändert, Claude Code sieht den Block in allen drei Varianten. Bewusst stehen gelassen: wird das Entfernen wegen eines jaira:local-Bereichs abgelehnt, ist die gewählte Datei schon geschrieben; der Block steht dann in beiden Dateien wie vorher, nichts geht verloren."
review-gaps: "none — removeBlock teilt sich die Markersuche nicht mit managedBlock, weil es entfernt statt ersetzt; hasLine ist der einzige Zeilenvergleich im Paket."
test-verdict: "pass — go test ./core/board ./internal/cli ./internal/tui grün (leeres HOME). E2E mit frischem Build in leerem Repo: update --agent-file agents -> Block nur in AGENTS.md, CLAUDE.md mit eigenem Inhalt bekommt @AGENTS.md; update ohne Flag danach -> AGENTS.md:1 CLAUDE.md:0; --agent-file cursor -> exit 3 mit Meldung."
question: "Soll dieses Repo selbst auf --agent-file agents umgestellt werden, oder bleibt es bei beiden?"
---

# jaira update wählt, in welche Agent-Datei der Block geht

## Definition of Done

- [x] jaira update --agent-file agents schreibt den Block nur nach AGENTS.md, nimmt ihn aus CLAUDE.md und sorgt dafür, dass Claude Code ihn trotzdem liest (leere CLAUDE.md gelöscht, sonst @AGENTS.md darin)
  proof: TestChooseAgentsDeletesAClaudeFileHoldingOnlyTheBlock, TestChooseAgentsImportsIntoAClaudeFileWithOwnContent
- [x] jaira update --agent-file claude schreibt nur nach CLAUDE.md und nimmt den Block aus AGENTS.md (leere Datei gelöscht)
  proof: TestChooseClaudeRemovesTheBlockFromAgentsAndKeepsTheRest
- [x] Ein späteres jaira update ohne Flag und eine Lane-Änderung schreiben nur in die gewählte Datei
  proof: TestTheChoiceHoldsForTheNextWrite; core/board/announce.go chosenAgentFiles
- [x] Tests in core/board decken agents, claude, both und das Behalten der Wahl ab
  proof: core/board/agentfile_test.go
- [x] core/release/NOTES.md hat eine Zeile unter Unreleased
  proof: core/release/NOTES.md:17

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-06 18:49 · Alexander Sacharov** — Die Wahl lebt in den Dateien: geschrieben wird nur dorthin, wo der Block schon steht; nirgends -> beide. Darum ändern sich die Aufrufer (update ohne Flag, lanes, TUI, init) nicht. Ein Block mit jaira:local wird nicht entfernt, sondern mit Fehler abgelehnt — sonst ginge handgeschriebener Text verloren. Die eigene @AGENTS.md-Zeile trägt einen Kommentar darüber, damit 'both'/'claude' nur jairas Import wieder entfernt, nie einen vom Nutzer. Dieses Repo selbst bleibt auf 'both'.
- **2026-10-06 18:57 · Alexander Sacharov** — Alex am 2026-10-06: dieses Repo bleibt auf beiden Dateien (AGENTS.md und CLAUDE.md); --agent-file ist nur verfügbar, nicht angewendet.
