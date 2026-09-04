---
id: 01M1PMZDDEQFFYGSR2WFTYTVBZ
title: Der review-check gehoert der review-Lane allein
status: critique
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: Nur die review-Lane deklariert und schreibt review-check; optimize produziert nur noch review-gaps - das (optimize)-Provenienz-Label vor dem check verschwindet
context: "Berk am 04.09.: 'Check sollte nur der Review Lane gehoeren.' Historie: optimize war als letzter Agent vor dem Menschen gedacht (Prompt nannte check 'the handover'), review kam nach human dazu und deklariert ihn ebenfalls - seither schulden ZWEI Lanes dasselbe Feld und die Detailansicht zeigt '(optimize/review)' davor. Aenderung an ZWEI Kopien: lanes/optimize.md (Katalog im Repo, committet) und .jaira/lanes/optimize.md (Board-Kopie, gitignored) - output-produces verliert review-check, der Prompt verliert den check-Abschnitt ('Then write review-check ... Review the diff is not.'). Der verwaltete CLAUDE.md-Block nennt optimizes Output erst nach 'jaira update' neu (steht ohnehin aus wegen 76WCCW)."
definition-of-done: lanes/optimize.md und .jaira/lanes/optimize.md deklarieren nur review-gaps; der Prompt verlangt keinen check mehr; shipped-Lane-Parsing bleibt gruen; ein Ticket in optimize kann ohne check nach human/review weiter
tags: []
blocked-by: []
commits: []
created-at: 2026-09-04T16:45:35Z
updated-at: 2026-09-04T16:49:40Z
claimed-by: EE-3NX6GL3-4183114
claimed-at: 2026-09-04T16:46:02Z
updated-by: BeMuCa
outcome-what: "review-check aus der optimize-Lane entfernt (Katalog lanes/optimize.md + Board-Kopie): output-produces nur noch review-gaps, der Prompt verweist die Hand-Pruefung an die review-Lane"
outcome-why: "Berk am 04.09.: der check gehoert der review-Lane allein - zwei deklarierende Lanes erzeugten das verwirrende (optimize/review)-Provenienz-Label und doppelte Schreibpflicht"
outcome-resolves: "lanes show optimize zeigt Output: review-gaps; shipped-Parsing gruen; go test ./... -race RC=0; der Gate-Beweis (optimize ohne check verlassen) ist der Weg dieses Tickets selbst"
executed-by: fable
---

# Der review-check gehoert der review-Lane allein

## Definition of Done

- [x] lanes/optimize.md und .jaira/lanes/optimize.md deklarieren nur review-gaps; der Prompt verlangt keinen check mehr; shipped-Lane-Parsing bleibt gruen; ein Ticket in optimize kann ohne check nach human/review weiter
  proof: beide optimize.md: output-produces [review-gaps], check-Abschnitt aus dem Prompt raus; 'jaira lanes show optimize' zeigt Output: review-gaps; go test ./... -race RC=0 inkl. shipped-Lane-Parsing; Gate-Beweis: dieses Ticket passiert optimize ohne check

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

