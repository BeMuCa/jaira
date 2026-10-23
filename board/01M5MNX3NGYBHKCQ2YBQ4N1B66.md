---
id: 01M5MNX3NGYBHKCQ2YBQ4N1B66
title: Die Testing-Lane verlangt keinen vollen Suite-Lauf mehr; die Suite laeuft einmal vor dem PR
status: backlog
ready: false
creator: BeMuCa
goal: "Der Prompt der Testing-Lane sagt, was die Regel ist: Paket-Tests der beruehrten Pakete in der Lane, die volle Race-Suite einmal vor dem PR durch die Hauptsession - und kein Worker startet sie mehr aus eigenem Antrieb."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkt 15): ein Worker lief die volle Suite, obwohl die Regel 'einmal pro Milestone' lautet. Die Regel stand nur im Gedaechtnis der Session, nicht im Lane-Prompt.
  Nachgeprueft am 20.10.: .jaira/lanes/testing.md:22 dieses Boards verlangt ausdruecklich 'run the whole suite ... with the race detector'. Prompt und Regel widersprechen sich; die Regel stammt vom 14.10. ('Test nur final bevor wir pr preparieren'), weil die Suite 7-8 Minuten dauert.
  Betrifft auch die Katalog-Lane lanes/testing.md im Repo, falls sie denselben Satz traegt.
definition-of-done: ".jaira/lanes/testing.md und die Katalog-Lane sagen: in der Lane die Tests der beruehrten Pakete (go test ./<pkg> -count=1), die volle Suite mit -race einmal vor dem PR durch die Hauptsession; der Satz 'run the whole suite' ist weg"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:26:07Z
updated-at: 2026-10-23T15:26:07Z
---

# Die Testing-Lane verlangt keinen vollen Suite-Lauf mehr; die Suite laeuft einmal vor dem PR

## Definition of Done

- [ ] .jaira/lanes/testing.md und die Katalog-Lane sagen: in der Lane die Tests der beruehrten Pakete (go test ./<pkg> -count=1), die volle Suite mit -race einmal vor dem PR durch die Hauptsession; der Satz 'run the whole suite' ist weg
- [ ] jaira-worker/SKILL.md und jaira-dispatcher/SKILL.md tragen dieselbe Regel in einer Zeile
- [ ] Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

