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
updated-at: 2026-09-13T15:55:03Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-4206
claimed-at: 2026-09-13T15:44:03Z
outcome-what: "Tastenkommandos werden ueber die physische Tastenposition gelesen statt ueber das Zeichen, das die Belegung druckt"
outcome-why: "Mit kyrillischer Belegung traf kein einziger case in den TUI-Switches, das Board reagierte auf gar nichts mehr"
outcome-resolves: "cmdKey in internal/tui/keylayout.go: erst Key.BaseCode vom Terminal, sonst JZUKEN-Positionstabelle; eingesetzt in model.go, home.go, edit.go; Texteingabe liest weiter k.Text und bleibt kyrillisch"
review-summary: |-
  internal/tui/keylayout.go:70 physicalRune trusts Key.BaseCode for every keypress, also fuer benannte Tasten: spricht das Terminal das Kitty-Protokoll und meldet ein BaseCode fuer enter, space oder einen Pfeil, gibt cmdKey dessen Rune roh zurueck statt 'enter'/'space' - Abbruch nur auf gedruckte Zeichen einschraenken, indem k.Text genau eine graphische Nicht-Leerzeichen-Rune sein muss
  internal/tui/home.go:215 vier Aufrufstellen bekamen cmdKey einzeln; das ist die vorhandene Form (jeder Screen nimmt einen string entgegen) und bleibt richtig - keine Aenderung
  internal/tui/keylayout.go:85 usPosition traegt ausser Kyrillisch nur '.' - kein spekulativer Ausbau auf Griechisch/Hebraeisch, das ist so richtig
---

# Tastenkuerzel funktionieren auch bei kyrillischem Layout

## Definition of Done

- [x] Mit russischem Layout steuern j/k/h/l, q, enter, / und die uebrigen Kommandotasten das Board wie mit englischem Layout; in Such- und Editfeldern erscheinen kyrillische Zeichen weiterhin als Text
  proof: internal/tui/keylayout_test.go: TestBoardAnswersACyrillicLayout und TestTypingStaysCyrillicInTheFilter, go test ./... gruen

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
  proof: internal/tui/model.go:877, home.go, edit.go:82
- [x] KeyboardEnhancements.ReportAllKeysAsEscapeCodes + ReportAlternateKeys in beiden View()-Funktionen setzen, damit BaseCode ueberhaupt ankommt
  proof: internal/tui/view.go, internal/tui/home.go View()
- [x] TUI-Test: kyrillische Tasten steuern das Board (j/k/h/l/q//), und im Editfeld landen sie weiterhin als Text
  proof: internal/tui/keylayout_test.go TestBoardAnswersACyrillicLayout, TestTypingStaysCyrillicInTheFilter
- [x] Zeile in core/release/NOTES.md unter ## Unreleased
  proof: core/release/NOTES.md

## Progress
- **2026-09-13 15:44 · Alexander Sacharov** — Alexander Sacharov took this ticket over from AlSa
- **2026-09-13 15:44 · Alexander Sacharov** — Nicht nur Kyrillisch loesen. Zwei Ebenen: (1) Das Terminal nach dem physischen Tastenwert fragen - View.KeyboardEnhancements.ReportAllKeysAsEscapeCodes + ReportAlternateKeys setzen, dann liefert bubbletea v2 Key.BaseCode (key.go:326) die Taste nach PC-101-Layout, egal welche Belegung aktiv ist. Damit ist jede Sprache erledigt, ohne eine einzige Tabelle. Geht aber nur in Terminals mit Kitty-Keyboard-Protokoll (kitty, ghostty, WezTerm, foot) oder ueber die Windows Console API - Windows Terminal kann es nicht, und genau da sitzt der Melder (WSL2). (2) Fallback fuer den Rest: pro Schrift eine Positions-Tabelle, ausgeloest ueber unicode.Is(unicode.Cyrillic, r). Latein-Layouts wie AZERTY oder QWERTZ brauchen nichts, die liefern ohnehin ASCII-Buchstaben - kaputt sind nur die Nicht-Latein-Schriften. Reihenfolge: erst BaseCode pruefen, nur wenn leer die Tabelle.
- **2026-09-13 15:46 · Alexander Sacharov** — Warum die Reihenfolge BaseCode-zuerst: bubbletea v2 fuellt Key.BaseCode (key.go:326) nur, wenn das Terminal das Kitty-Keyboard-Protokoll spricht oder die Windows Console API benutzt wird. Dort ist jede Belegung der Welt erledigt, ohne eine Tabelle. Windows Terminal - wo der Melder unter WSL2 sitzt - kann es nicht, deshalb ist die Positions-Tabelle keine Alternative sondern die zweite Ebene darunter. Latein-Belegungen (AZERTY, QWERTZ, Dvorak) brauchen gar nichts: sie liefern ASCII-Buchstaben, die Switches treffen schon heute. Kaputt sind ausschliesslich die Nicht-Latein-Schriften, und davon bekommt dieses Ticket nur Kyrillisch - Griechisch, Hebraeisch und Arabisch sind dieselbe Mechanik und eine Tabelle mehr, sobald jemand sie meldet. Zweite Erkenntnis: ReportAlternateKeys allein reicht nicht. Die bubbletea-Doku (tea.go:252) sagt, dass nur als Escape-Code gemeldete Tasten davon betroffen sind - normale Buchstaben kommen als Klartext. Ohne ReportAllKeysAsEscapeCodes bleibt BaseCode bei 'j' also leer, und die erste Ebene laeuft ins Nichts.
- **2026-09-13 15:53 · Alexander Sacharov** — Zwei Entscheidungen beim Bauen, die man dem Code nicht ansieht. (1) Key.String() ist als Quelle unbrauchbar: ultraviolet/key.go:391 gibt Key.Text zurueck, sobald die Taste Text erzeugt hat, und ignoriert BaseCode genau dort, wo es gebraucht wird - bei Buchstaben. Keystroke() liest BaseCode zwar, schreibt Shift aber als 'shift+g' statt 'G', und das Board bindet 'G', 'E', 'X', 'S' als Grossbuchstaben. Deshalb baut cmdKey den String selbst zusammen. (2) '.' ist als einzige Latein-Taste in der Tabelle: auf JZUKEN sitzt der Punkt auf der Taste, die auf US-Layout '/' ist, und ohne diese Zeile ist der Filter mit russischer Belegung gar nicht erreichbar. Geprueft, dass '.' und ',' in keinem Screen als Kommando gebunden sind - ein US-Nutzer verliert also nichts, er bekommt einen zweiten Weg in den Filter. NICHT lokal verifizierbar: die KeyboardEnhancements-Ebene. Hier laeuft Windows Terminal unter WSL2, das kein Kitty-Protokoll spricht, also kommt BaseCode nie an und der Zweig bleibt ungetestet ausser im Unit-Test mit gesetztem BaseCode. Wer an einem kitty, ghostty, WezTerm oder foot sitzt, sollte einmal Tippen und Steuern pruefen: ReportAllKeysAsEscapeCodes schaltet Klartext ab, ReportAssociatedText muss den Text zurueckbringen - stimmt das in einem echten Terminal nicht, bleiben Filter und Editfeld leer.
- **2026-09-13 15:54 · Alexander Sacharov** — critique: cmdKey darf nur bei gedruckten Zeichen eingreifen. Heute reicht ein gesetztes Key.BaseCode, und bubbletea setzt das auf einem Kitty-Terminal auch fuer benannte Tasten - dann liefert cmdKey die nackte Rune von enter oder space statt 'enter'/'space', und die Switches treffen nichts mehr. Auf Windows Terminal faellt das nie auf, weil BaseCode dort immer 0 ist: der Fehler waere genau auf den Terminals aufgetreten, fuer die die erste Ebene ueberhaupt gebaut wurde. Fix: in physicalRune zuerst k.Text pruefen - genau eine Rune, unicode.IsGraphic, kein Leerzeichen -, sonst k.String() unveraendert.
