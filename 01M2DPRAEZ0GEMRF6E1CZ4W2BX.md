---
id: 01M2DPRAEZ0GEMRF6E1CZ4W2BX
title: Tastenkuerzel funktionieren auch bei kyrillischem Layout
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: Die Board-Tasten reagieren unabhaengig vom aktiven Tastaturlayout
context: "Wer die TUI mit russischer Tastaturbelegung bedient, kann sie gar nicht steuern: keine einzige Taste tut etwas. Grund: internal/tui/model.go vergleicht die Taste als Zeichen ('j', 'k', 'q', ...). Bei kyrillischem Layout liefert das Terminal 'о', 'л', 'й' - kein case trifft, die Taste faellt durch. Betroffen sind alle Kommando-Switches: model.go ab Zeile 871, home.go 215/230/237/258, edit.go 82. Gesucht ist eine Umsetzung der JZUKEN-Positionen auf ihre QWERTY-Entsprechung, bevor ein Kommando gesucht wird. Wichtig: in den Text-Eingabemodi (Suche, Bearbeiten) muss das kyrillische Zeichen unveraendert durchgehen, sonst kann niemand mehr russisch tippen."
definition-of-done: "Mit russischem Layout steuern j/k/h/l, q, enter, / und die uebrigen Kommandotasten das Board wie mit englischem Layout; in Such- und Editfeldern erscheinen kyrillische Zeichen weiterhin als Text"
tags:
  - tui
blocked-by: []
commits: []
created-at: 2026-09-13T15:39:12Z
updated-at: 2026-09-13T15:44:03Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-4206
claimed-at: 2026-09-13T15:44:03Z
---

# Tastenkuerzel funktionieren auch bei kyrillischem Layout

## Definition of Done

- [ ] Mit russischem Layout steuern j/k/h/l, q, enter, / und die uebrigen Kommandotasten das Board wie mit englischem Layout; in Such- und Editfeldern erscheinen kyrillische Zeichen weiterhin als Text

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-13 15:44 · Alexander Sacharov** — Alexander Sacharov took this ticket over from AlSa
