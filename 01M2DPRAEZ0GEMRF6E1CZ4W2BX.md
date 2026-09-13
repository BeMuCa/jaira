---
id: 01M2DPRAEZ0GEMRF6E1CZ4W2BX
title: Tastenkuerzel funktionieren auch bei kyrillischem Layout
status: critique
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
updated-at: 2026-09-13T18:35:26Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-4206
claimed-at: 2026-09-13T15:44:03Z
outcome-what: "AltGr faellt nicht mehr in die Positions-Umsetzung; Shift-Guard und AltGr sind jetzt von Tests festgenagelt"
outcome-why: "Der Windows-Decoder setzt Key.Text auch bei AltGr, und die Buchstaben-Ebene sah gar nicht auf k.Mod - AltGr+о waere als blankes j im Kommando-Handler gelandet"
outcome-resolves: "physicalRune steigt aus, sobald k.Mod etwas ueber ModShift|ModCapsLock hinaus traegt; TestCmdKeyIgnoresAltGr und der AZERTY-Fall in TestCmdKeyLeavesShiftedPunctuationAlone fallen ohne die Guards um; Kopfkommentar und NOTES-Zeile auf Buchstaben eingeschraenkt"
review-summary: |-
  Vierter Durchgang: der AltGr-Guard sitzt als eine Bedingung in physicalRune (internal/tui/keylayout.go:69), direkt neben der Bedingung, die benannte Tasten aussortiert - keine neue Ebene, keine Verzweigung bei den Aufrufern
  internal/tui/keylayout.go:69 k.Mod&^(ModShift|ModCapsLock) statt einer Aufzaehlung von ModCtrl und ModAlt: die Frage ist 'aendert der Modifier nur das gedruckte Zeichen', und das ist bei genau diesen beiden der Fall - eine Liste der verbotenen Modifier waere unvollstaendig, sobald einer dazukommt
  Kein Fund mehr offen; die Einschraenkung auf Buchstaben und die beiden Guards erzaehlen im Kommentar dieselbe Geschichte wie im Code
review-gaps: |-
  internal/tui/keylayout.go:80 AltGr faellt auf dem Windows-Console-Pfad in die Buchstaben-Ebene: der Decoder setzt Key.Text auch bei AltGr (decoder.go:2035, LEFT_CTRL|RIGHT_ALT) und Mod traegt dann ModCtrl|ModAlt. AltGr+о wuerde als blankes 'j' im Kommando-Handler landen und den Cursor bewegen, wo auf master nichts passierte. Gleiches fuer AltGr-Buchstaben lateinischer Belegungen (polnisches AltGr+a = ą, BaseCode 'a')
  internal/tui/keylayout_test.go:44 Der ModShift-Guard ist von keinem Test festgenagelt: '?', '!' und '>' stehen ohnehin nicht in usPosition, die Faelle bestehen auch ohne Guard. Der Fall, fuer den der Guard wirklich da ist, fehlt - auf AZERTY ist der Punkt shift+';', ohne Guard haette cmdKey '/' zurueckgegeben und beim Tippen eines Punktes den Filter geoeffnet
  internal/tui/keylayout.go:20 Der Kopfkommentar sagt weiterhin, BaseCode sei fuer jede Belegung richtig und brauche keine Tabelle; nach der Einschraenkung gilt das nur noch fuer Buchstaben
  core/release/NOTES.md:18 'and the rest' im ersten Punkt deckt auch ctrl+d/ctrl+u mit ab, die genau nicht mitkommen
  Bestaetigt geschlossen: die Hilfe-Taste ist auf allen drei Pfaden wieder erreichbar (auf JZUKEN ist '?' shift+7 und laeuft unveraendert durch), shift+Ziffer schaltet keine Boards mehr um, der tote ctrl/alt-Zweig ist raus, die NOTES-Zeile stimmt
test-verdict: |-
  go build, go vet, go test ./... gruen nach dem Shift-Fix
  Neu: TestCmdKeyLeavesShiftedPunctuationAlone deckt shift+/ ('?'), shift+1 ('!') und shift+. ('>') mit gesetztem BaseCode ab - genau der Pfad, der die Hilfe-Taste in den Filter geschickt haette
  Weiterhin gruen: TestBoardAnswersACyrillicLayout, TestTypingStaysCyrillicInTheFilter, TestCmdKeyIgnoresBaseCodeOnNamedKeys, TestCmdKeyMapsCyrillicToItsPhysicalKey, TestCmdKeyKeepsCase, TestCmdKeyPrefersTheTerminalsOwnBaseCode
  Von Hand bestaetigt (vor dem Shift-Fix, unveraendertem Pfad): russische Belegung steuert das Board auf Windows Terminal unter WSL2
  Offen bleibt, was hier kein Terminal hergibt: der BaseCode-Pfad eines echten Kitty-Terminals
  Binary neu gebaut unter /home/alex/.local/bin/jaira
question: ""
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
- **2026-09-13 16:58 · Alexander Sacharov** — Mensch hat es an seinem Terminal bestaetigt: mit russischer Belegung steuert das Board wie mit englischer. Windows Terminal unter WSL2, also ueber die Positions-Tabelle, nicht ueber BaseCode - die Kitty-Ebene bleibt weiterhin unbestaetigt, sie schadet hier aber nachweislich nicht.
- **2026-09-13 17:05 · Alexander Sacharov** — review: Shift ist das Loch. physicalRune darf Satzzeichen nicht ueber BaseCode ersetzen - shift+/ ist '?' und meldet trotzdem '/', also oeffnete die Hilfe-Taste den Filter, auf dem Windows-Console-Pfad sogar mit US-Belegung. Fix: nur Buchstaben ueber BaseCode/Tabelle umsetzen, Satzzeichen ausschliesslich ueber die Tabelle und nur ohne Shift. Zweitens: der ctrl/alt-Zweig in cmdKey kann nie laufen, weil der Decoder Key.Text bei jedem Modifier ueber Shift hinaus leert - raus damit, und ehrlich hinschreiben, dass ctrl-Kombinationen Sache des Terminals bleiben.
- **2026-09-13 18:29 · Alexander Sacharov** — Mensch hat nachgetestet: '?' oeffnet wieder die Hilfe, russische Belegung steuert weiterhin. Windows Terminal unter WSL2.
- **2026-09-13 18:33 · Alexander Sacharov** — review 2: AltGr ist das verbliebene Loch. Die Buchstaben-Ebene schaut gar nicht auf k.Mod, und der Windows-Decoder setzt Key.Text auch bei AltGr - also kommt AltGr+Buchstabe als blanke Kommandotaste an. Fix: vor beiden Zweigen aussteigen, sobald k.Mod etwas ueber ModShift|ModCapsLock hinaus traegt. Dazu ein Test, der den Shift-Guard wirklich festnagelt (AZERTY: shift+';' druckt '.', darf nicht '/' werden), und zwei Kommentare/NOTES-Zeilen, die nach der Einschraenkung zu viel versprechen.
