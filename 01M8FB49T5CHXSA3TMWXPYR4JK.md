---
id: 01M8FB49T5CHXSA3TMWXPYR4JK
title: "Cockpit: Ticket mit offener Frage in human zeigt Frage statt Abnahme-Checkliste"
status: backlog
ready: true
creator: Sascha
goal: "Steht ein Ticket in einer Frage-Lane (human) mit offener Frage, zeigt der cockpit die Frage und ein Antwortfeld – keine 'Check yourself'-Liste und kein Accept, das danach abgelehnt wird."
context: |-
  Was falsch ist: ein Worker stellt eine Frage und schiebt das Ticket nach human. Der cockpit zeigt dafür denselben Abnahme-Bildschirm wie in signoff: 'Check yourself - tick what you have seen' mit den 'seen by a person'-Punkten der DoD.
  Die Person hakt einen Punkt ab, drückt Execute und bekommt erst dann: 'nothing to accept in this lane - mark a note fix, or type what is missing' (texts.ts:80/89, step 'noNext').
  Folge: die Person glaubt, sie solle abnehmen, und hakt DoD-Punkte ab, die erst am Ende prüfbar sind (die Arbeit ist noch gar nicht fertig).
  Wo: core/mod/builtin/jaira-cockpit/hooks/board.ts:1231 wählt den step; noNext greift erst beim Execute. register.tsx:2665 zeigt die seen-Liste unabhängig davon.
  Richtig wäre: Lane ohne next_lane + Ticket hat question → Frage oben zeigen, Antwort eintippen = Antwort (note + zurück an den Worker), seen-Liste ausblenden.
  Gesehen: 27.11.2026 in einem Projekt mit Video-Board.
definition-of-done: "Ticket in human mit question: cockpit zeigt die Frage, keine seen-Checkliste (Test in tests/board.test.ts)"
tags:
  - cockpit
blocked-by: []
related: []
commits: []
created-at: 2026-11-27T20:28:41Z
updated-at: 2026-11-27T20:28:51Z
assignee: Sascha
updated-by: Sascha
---

# Cockpit: Ticket mit offener Frage in human zeigt Frage statt Abnahme-Checkliste

## Definition of Done

- [ ] Ticket in human mit question: cockpit zeigt die Frage, keine seen-Checkliste (Test in tests/board.test.ts)
- [ ] Eingetippte Antwort landet als note am Ticket und geht an den Worker zurück; kein 'nothing to accept'-Fehler (Test)
- [ ] signoff-Abnahme unverändert (bestehende Tests grün)
- [ ] Von einer Person im cockpit gesehen
- [ ] Release-Note core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

