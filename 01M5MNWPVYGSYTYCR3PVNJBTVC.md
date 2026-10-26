---
id: 01M5MNWPVYGSYTYCR3PVNJBTVC
title: "Ein offener Plan-Punkt sperrt die Ziellane nicht mehr, wenn die Definition of Done erfuellt ist"
status: backlog
ready: false
creator: BeMuCa
goal: "Ein Ticket mit erfuellter Definition of Done kommt nach done, auch wenn ein Plan-Schritt nicht abgehakt ist; der Plan ist die Methode, die DoD das Ergebnis."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkt 8): GH0GKD mit DoD 9/9 wurde an done abgelehnt, weil Plan-Schritt 30 offen war.
  Nachgeprueft am 23.10.: core/gate/gate.go:411-424 verweigert jede Ziellane mit CodePlanIncomplete, solange ein Plan-Punkt offen ist. Der Kommentar dort sagt, Plan und DoD seien nicht unabhaengig.
  CLAUDE.md dieses Repos sagt das Gegenteil: 'Ticking the plan does not tick the definition of done' - der Plan ist eine zweite, separate Liste.
  Nicht entschieden: Gate weg (nur Warnung), oder der Worker hakt beim Abschluss den letzten Plan-Schritt selbst ab. Brainstorm entscheidet; die einfachere Variante ist das Gate zu einer Warnung zu machen.
definition-of-done: Ein Move in eine Ziellane mit vollstaendiger DoD und offenem Plan-Punkt wird nicht mehr abgelehnt; die Ausgabe nennt den offenen Plan-Punkt als Hinweis; ein Test deckt den Fall
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:25:53Z
updated-at: 2026-10-26T19:10:34Z
assignee: BeMuCa
updated-by: BeMuCa
---

# Ein offener Plan-Punkt sperrt die Ziellane nicht mehr, wenn die Definition of Done erfuellt ist

## Definition of Done

- [ ] Ein Move in eine Ziellane mit vollstaendiger DoD und offenem Plan-Punkt wird nicht mehr abgelehnt; die Ausgabe nennt den offenen Plan-Punkt als Hinweis; ein Test deckt den Fall
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [x] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

