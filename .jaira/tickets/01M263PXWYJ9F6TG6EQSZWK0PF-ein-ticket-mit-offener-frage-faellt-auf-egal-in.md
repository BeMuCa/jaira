---
id: 01M263PXWYJ9F6TG6EQSZWK0PF
title: "Ein Ticket mit offener Frage faellt auf, egal in welcher Lane es liegt"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Eine gestellte, unbeantwortete Frage ist auf dem Board sichtbar, auch wenn das Ticket nicht in der human-Lane sitzt"
context: |-
  Auf dem Board eines anderen Projekts tragen 9GEGTB, NEEY9A und SXDBX6 ein gefuelltes question-Feld (931, 1059 und 19 Zeichen) und liegen in review, nicht in human.
  9GEGTB fragt woertlich 'Do I delete the stray {clarify_response} from prompts/v1_powerpoint.txt?' - eine Frage, die vor der Abnahme beantwortet gehoert.
  In review ist sie unsichtbar: die Lane heisst 'Review', nicht 'wartet auf deine Antwort', und die Karte zeigt sie nicht.
  jaira kann das schon anzeigen - internal/tui/view.go:1026 rendert bei einem question-Feld das Kennzeichen 'waiting on your answer'. Vor dem Bauen pruefen, ob das nur in der Detailansicht haengt oder auch auf der Karte, und warum es hier nicht greift.
definition-of-done: "Ein Ticket mit gefuelltem question-Feld ist auf der Board-Karte als solches erkennbar, unabhaengig von seiner Lane; ein Test deckt eine Frage ausserhalb der human-Lane ab; go test ./... -race gruen"
tags:
  - cli
blocked-by: []
commits: []
created-at: 2026-09-10T16:51:42Z
updated-at: 2026-09-10T16:51:42Z
---

# Ein Ticket mit offener Frage faellt auf, egal in welcher Lane es liegt

## Definition of Done

- [ ] Ein Ticket mit gefuelltem question-Feld ist auf der Board-Karte als solches erkennbar, unabhaengig von seiner Lane; ein Test deckt eine Frage ausserhalb der human-Lane ab; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

