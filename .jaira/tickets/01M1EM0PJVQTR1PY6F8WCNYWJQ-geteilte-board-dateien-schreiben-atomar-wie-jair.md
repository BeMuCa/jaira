---
id: 01M1EM0PJVQTR1PY6F8WCNYWJQ
title: Geteilte Board-Dateien schreiben atomar wie .jaira/tags
status: backlog
ready: false
creator: BeMuCa
goal: Die kleinen geteilten Store-Dateien (lanes/order voran) schreiben ueber ticket.WriteAtomic statt ueber ein truncierendes os.WriteFile
context: "Fund des Tag-Reviews am 01.09. (L19), bewusst aus 79GEPW herausgehalten: seit ticket.WriteAtomic exportiert ist, faellt auf, dass mehrere geteilte Dateien noch plain os.WriteFile nutzen - ein paralleler Leser kann eine halb geschriebene Datei sehen. Genannte Stellen: core/lane/order.go:62 (pikant: core/tag nennt order als sein Format-Vorbild, ist selbst aber atomar+gelockt), core/session/session.go:98, internal/cli/sync.go:79, core/board/announce.go:314/329/341, core/lane/share.go:54/116/172, core/board/gitignore.go (4 Stellen), core/lane/defaultboard.go:101/143, core/project/project.go:91. Je Stelle pruefen, ob Atomicitaet reicht oder (wie bei tags) auch ein Lock um Load-Modify-Save gehoert - nicht pauschal umstellen, sondern je Datei begruenden."
definition-of-done: "Jede genannte Stelle schreibt atomar oder traegt einen Kommentar, warum nicht noetig; Lock nur wo ein Load-Modify-Save-Fenster existiert; go test ./... -race gruen"
blocked-by: []
commits: []
created-at: 2026-09-01T13:54:53Z
updated-at: 2026-09-01T13:54:53Z
---

# Geteilte Board-Dateien schreiben atomar wie .jaira/tags

## Definition of Done

- [ ] Jede genannte Stelle schreibt atomar oder traegt einen Kommentar, warum nicht noetig; Lock nur wo ein Load-Modify-Save-Fenster existiert; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

