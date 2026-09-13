---
id: 01M2DPRAEZ0GEMRF6E1CZ4W2BX
title: Tastenkuerzel funktionieren auch bei kyrillischem Layout
status: in-progress
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
updated-at: 2026-09-13T15:53:26Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-4206
claimed-at: 2026-09-13T15:44:03Z
---

# Tastenkuerzel funktionieren auch bei kyrillischem Layout

## Definition of Done

- [ ] Mit russischem Layout steuern j/k/h/l, q, enter, / und die uebrigen Kommandotasten das Board wie mit englischem Layout; in Such- und Editfeldern erscheinen kyrillische Zeichen weiterhin als Text

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Belegte Stellen sammeln: jeden Kommando-Switch finden, der k.String() vergleicht (model.go, home.go, board/browse/drop-Screens) und von den Text-Eingabepfaden trennen
  proof: internal/tui/keylayout.go
- [x] internal/tui/keylayout.go anlegen: cmdKey(k tea.KeyPressMsg) string - erst k.BaseCode, sonst Positions-Tabelle je Schrift, sonst k.String() unveraendert
  proof: internal/tui/keylayout.go:37 cmdKey
- [x] Kyrillische Positions-Tabelle (JZUKEN, klein und gross) schreiben, ueber unicode.Is(unicode.Cyrillic, r) ausgeloest
  proof: internal/tui/keylayout.go usPosition
- [x] Tabellentest: jede kyrillische Taste ergibt ihre QWERTY-Position; Modifier bleiben erhalten (ctrl+о -> ctrl+j); ASCII bleibt unveraendert
  proof: internal/tui/keylayout_test.go TestCmdKey*
- [x] cmdKey in den Kommando-Switches einsetzen, Text-Eingabemodi (edit.go, Filter, Create) unangetastet lassen
- [x] KeyboardEnhancements.ReportAllKeysAsEscapeCodes + ReportAlternateKeys in beiden View()-Funktionen setzen, damit BaseCode ueberhaupt ankommt
- [x] TUI-Test: kyrillische Tasten steuern das Board (j/k/h/l/q//), und im Editfeld landen sie weiterhin als Text
- [x] Zeile in core/release/NOTES.md unter ## Unreleased

## Progress
- **2026-09-13 15:44 · Alexander Sacharov** — Alexander Sacharov took this ticket over from AlSa
- **2026-09-13 15:44 · Alexander Sacharov** — Nicht nur Kyrillisch loesen. Zwei Ebenen: (1) Das Terminal nach dem physischen Tastenwert fragen - View.KeyboardEnhancements.ReportAllKeysAsEscapeCodes + ReportAlternateKeys setzen, dann liefert bubbletea v2 Key.BaseCode (key.go:326) die Taste nach PC-101-Layout, egal welche Belegung aktiv ist. Damit ist jede Sprache erledigt, ohne eine einzige Tabelle. Geht aber nur in Terminals mit Kitty-Keyboard-Protokoll (kitty, ghostty, WezTerm, foot) oder ueber die Windows Console API - Windows Terminal kann es nicht, und genau da sitzt der Melder (WSL2). (2) Fallback fuer den Rest: pro Schrift eine Positions-Tabelle, ausgeloest ueber unicode.Is(unicode.Cyrillic, r). Latein-Layouts wie AZERTY oder QWERTZ brauchen nichts, die liefern ohnehin ASCII-Buchstaben - kaputt sind nur die Nicht-Latein-Schriften. Reihenfolge: erst BaseCode pruefen, nur wenn leer die Tabelle.
- **2026-09-13 15:46 · Alexander Sacharov** — Warum die Reihenfolge BaseCode-zuerst: bubbletea v2 fuellt Key.BaseCode (key.go:326) nur, wenn das Terminal das Kitty-Keyboard-Protokoll spricht oder die Windows Console API benutzt wird. Dort ist jede Belegung der Welt erledigt, ohne eine Tabelle. Windows Terminal - wo der Melder unter WSL2 sitzt - kann es nicht, deshalb ist die Positions-Tabelle keine Alternative sondern die zweite Ebene darunter. Latein-Belegungen (AZERTY, QWERTZ, Dvorak) brauchen gar nichts: sie liefern ASCII-Buchstaben, die Switches treffen schon heute. Kaputt sind ausschliesslich die Nicht-Latein-Schriften, und davon bekommt dieses Ticket nur Kyrillisch - Griechisch, Hebraeisch und Arabisch sind dieselbe Mechanik und eine Tabelle mehr, sobald jemand sie meldet. Zweite Erkenntnis: ReportAlternateKeys allein reicht nicht. Die bubbletea-Doku (tea.go:252) sagt, dass nur als Escape-Code gemeldete Tasten davon betroffen sind - normale Buchstaben kommen als Klartext. Ohne ReportAllKeysAsEscapeCodes bleibt BaseCode bei 'j' also leer, und die erste Ebene laeuft ins Nichts.
