---
id: 01M83Q9ETRH50TZBG2X1PZJW3D
title: "TUI lässt sich mit der Maus bedienen: Karten anklicken, mit dem Rad scrollen, Knöpfe und Reiter klicken"
status: ideas
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer jaira im Terminal öffnet, kann Karten, Reiter, Knöpfe und Filter mit der Maus anklicken und mit dem Rad scrollen; per Einstellung abschaltbar, damit Textmarkieren weiter geht."
context: |-
  Was heute fehlt (Alex 23.11.2026): die TUI kennt keine Maus, nur Tasten.
  Technik: Bubble Tea v2 (charm.land/bubbletea/v2) schaltet Maus mit v.MouseMode = tea.MouseModeCellMotion in View() (internal/tui/view.go:108 neben AltScreen) ein; Klicks kommen als tea.MouseClickMsg mit x,y.
  Die eigentliche Arbeit: jede Ansicht muss wissen, was an Zelle x,y liegt (Karte, Reiter, Knopf, Listenzeile). internal/tui hat ~100 Dateien, ~27k Zeilen.
  Schätzung: Board-Karten + Rad ~halber Tag; alles (Reiter, Knöpfe, Filter, Fenster) 2-3 Tage, besser in 3-4 Tickets je Ansicht.
  Nebenwirkung: mit Maus-Modus markieren die meisten Terminals Text nur noch mit Shift - daher Einstellung im Settings-Fenster (S).
  Vorbild für klickbare Knöpfe: das Cockpit (core/mod/builtin/jaira-cockpit), R19PY9 Kopfzeile.
definition-of-done: "Erster Schnitt: Klick auf eine Karte wählt/öffnet sie, Mausrad scrollt Board und Listen (Test über MouseMsg)"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-11-23T08:10:20Z
updated-at: 2026-11-23T08:10:32Z
updated-by: Alexander Sacharov
round: 2
back-from: backlog
back-by: person
back-count: 1
---

# TUI lässt sich mit der Maus bedienen: Karten anklicken, mit dem Rad scrollen, Knöpfe und Reiter klicken

## Definition of Done

- [ ] Erster Schnitt: Klick auf eine Karte wählt/öffnet sie, Mausrad scrollt Board und Listen (Test über MouseMsg)
- [ ] Einstellung 'Maus' im Settings-Fenster (S) und in ~/.jaira/settings.json; aus = Verhalten wie heute (Test)
- [ ] Weitere Ansichten (Reiter, Knöpfe, Filter, Fenster) als Folge-Tickets aufgeteilt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

