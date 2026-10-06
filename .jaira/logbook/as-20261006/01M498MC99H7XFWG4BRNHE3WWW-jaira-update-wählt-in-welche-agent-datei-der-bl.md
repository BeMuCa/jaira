---
id: 01M498MC99H7XFWG4BRNHE3WWW
title: "jaira update wählt, in welche Agent-Datei der Block geht"
status: done
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
commits:
  - aff2304a20867cc8eaa42bf7d2e4d534b5153c82
  - 9dfd359398e27dce1752727002d25bb857afaba0
  - 01cd6d3b29e3e466f61a2b840e0abb73090bd51e
created-at: 2026-10-06T18:46:48Z
updated-at: 2026-10-06T21:53:40Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3243
claimed-at: 2026-10-06T18:46:59Z
outcome-what: "jaira update --agent-file agents|claude|both; core/board ChooseAgentFiles, removeBlock, setClaudeImport; AnnounceInAgentFiles schreibt nur in Dateien, die den Block tragen."
outcome-why: "Wer nur AGENTS.md oder nur CLAUDE.md will, bekam die andere Datei bei jedem update zurück; Claude Code liest AGENTS.md nur ohne CLAUDE.md."
outcome-resolves: "Die Wahl ist ein Flag, hält ohne Konfiguration, und Claude Code sieht den Block in jeder Variante."
review-summary: "jaira update hat ein Flag --agent-file agents|claude|both. Es schreibt den jaira-Block nur in die gewählten Dateien und nimmt ihn aus der anderen heraus (eine Datei, die danach leer ist, wird gelöscht; der Text drumherum bleibt mit Leerzeile getrennt). Bei 'agents' bekommt eine CLAUDE.md mit eigenem Inhalt oben eine markierte Zeile @AGENTS.md, weil Claude Code AGENTS.md sonst nicht liest; 'claude' und 'both' nehmen nur diese markierte Zeile wieder heraus. Ohne Flag schreibt jaira nur in Dateien, die den Block schon tragen, und in beide, wenn keine ihn trägt — so hält die Wahl ohne Konfiguration, und alte Boards verhalten sich wie bisher. Ein Block mit jaira:local-Bereich wird nicht entfernt, sondern mit Fehler abgelehnt. Ein falscher Wert endet mit Exit 3."
review-gaps: "none — bewusst offen: wird das Entfernen wegen jaira:local abgelehnt, ist die gewählte Datei schon geschrieben und der Block steht in beiden Dateien; nichts geht verloren."
test-verdict: "pass — go test ./core/board ./internal/cli grün (leeres HOME), 9 Tests in core/board/agentfile_test.go."
question: "Soll dieses Repo selbst auf --agent-file agents umgestellt werden, oder bleibt es bei beiden?"
review-verdict: "der Diff erfüllt alle sechs DoD-Punkte; keine Defekte mehr gefunden."
review-check: "1. go build -o /tmp/j ./cmd/jaira  2. mkdir /tmp/t && cd /tmp/t && git init -q && printf '# Own\\n' > CLAUDE.md && /tmp/j init  3. /tmp/j update --agent-file agents -> Ausgabe nennt 'CLAUDE.md (block removed), CLAUDE.md (imports AGENTS.md)'  4. head -3 CLAUDE.md -> Kommentarzeile, dann '@AGENTS.md'  5. grep -c 'jaira:start' AGENTS.md CLAUDE.md -> AGENTS.md:1 CLAUDE.md:0  6. /tmp/j update -> danach grep wieder AGENTS.md:1 CLAUDE.md:0  7. /tmp/j update --agent-file cursor -> Meldung 'want agents, claude or both', echo $? zeigt 3"
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
- [x] Steht der Block mitten in der Datei, bleibt nach dem Entfernen zwischen dem Text davor und danach eine Leerzeile; ein Test pinnt das
  proof: TestRemovingABlockMidFileKeepsTheParagraphsApart

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-06 18:49 · Alexander Sacharov** — Die Wahl lebt in den Dateien: geschrieben wird nur dorthin, wo der Block schon steht; nirgends -> beide. Darum ändern sich die Aufrufer (update ohne Flag, lanes, TUI, init) nicht. Ein Block mit jaira:local wird nicht entfernt, sondern mit Fehler abgelehnt — sonst ginge handgeschriebener Text verloren. Die eigene @AGENTS.md-Zeile trägt einen Kommentar darüber, damit 'both'/'claude' nur jairas Import wieder entfernt, nie einen vom Nutzer. Dieses Repo selbst bleibt auf 'both'.
- **2026-10-06 18:57 · Alexander Sacharov** — Alex am 2026-10-06: dieses Repo bleibt auf beiden Dateien (AGENTS.md und CLAUDE.md); --agent-file ist nur verfügbar, nicht angewendet.
- **2026-10-06 19:06 · Alexander Sacharov** — Review 2026-10-06: removeBlock klebte Text vor und nach einem Block in der Dateimitte mit nur einem Zeilenumbruch zusammen; als DoD-Punkt 6 angehängt.
- **2026-10-06 19:08 · Alexander Sacharov** — Von testing direkt nach review, ohne human: keine offene Frage (Entscheidung vom 2026-10-06, BG5QJ6).
- **2026-10-06 21:53 · Alexander Sacharov** — Abnahme 2026-10-06 (https://claude.ai/artifact/JwDPU5tVLGPXLp5uNFNGb2): alle Schritte von Alex ok, alle fünf Maschinenprüfungen grün.
