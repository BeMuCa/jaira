---
id: 01M2DPRAEZ0GEMRF6E1CZ4W2BX
title: Tastenkuerzel funktionieren auch bei kyrillischem Layout
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: Die Board-Tasten reagieren unabhaengig vom aktiven Tastaturlayout
context: "Wer die TUI mit russischer Tastaturbelegung bedient, kann sie gar nicht steuern: keine einzige Taste tut etwas. Grund: internal/tui/model.go vergleicht die Taste als Zeichen ('j', 'k', 'q', ...). Bei kyrillischem Layout liefert das Terminal 'о', 'л', 'й' - kein case trifft, die Taste faellt durch. Betroffen sind alle Kommando-Switches: model.go ab Zeile 871, home.go 215/230/237/258, edit.go 82. Gesucht ist eine Umsetzung der JZUKEN-Positionen auf ihre QWERTY-Entsprechung, bevor ein Kommando gesucht wird. Wichtig: in den Text-Eingabemodi (Suche, Bearbeiten) muss das kyrillische Zeichen unveraendert durchgehen, sonst kann niemand mehr russisch tippen."
definition-of-done: "Mit russischem Layout steuern j/k/h/l, q, enter, / und die uebrigen Kommandotasten das Board wie mit englischem Layout; in Such- und Editfeldern erscheinen kyrillische Zeichen weiterhin als Text"
tags:
  - tui
blocked-by: []
commits:
  - 401dbb12dd6c6cccc54ec365da9bbee39ee691fe
  - 5f741efaf92470683f4184e7343c3637f9e7ad73
  - 959e28e90bbf4ae47e6a1d5d2cb553c4fb3ee95f
  - fd3071416eb0af64f95a23b0268c769209b111f8
  - f6bc9d3d86d98a4c6f5950f802cda9943e350df2
  - 2131091d083a4440c721923dac74728c484c835b
  - 673d01004a815ebed38ceb74452aba1b5fc99df9
  - 8a0ce8d4626c37ba2558186d7b51a37362c3031b
  - 7729eb4c4bb128d9ad0f318c94d4b6f8f3b3570c
  - 61d73da4a8dcb1a69f41835b0fccee40f8b39f63
  - c1e4af3f4c501177cb96bf17e289797876317b9d
  - 52fd9acc456525a239d6de87c65d218d678c83a0
  - f5504c417d45a83995b2d78d529c91e7945d5737
  - d80c8e7a0e704dc7148236c15c0af6415c6ad975
created-at: 2026-09-13T15:39:12Z
updated-at: 2026-09-13T19:06:20Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-4206
claimed-at: 2026-09-13T15:44:03Z
outcome-what: "Lock-Zustaende zaehlen nicht mehr als gehaltener Modifier"
outcome-why: "Kitty-Terminals melden NumLock bei jedem Tastendruck - der Guard haette die Umsetzung dort ganz abgeschaltet"
outcome-resolves: "physicalRune erlaubt Shift, CapsLock, NumLock und ScrollLock; TestCmdKeyIgnoresTheLockStates haelt es fest"
review-summary: |-
  Kommandotasten werden nicht mehr als Zeichen gelesen, sondern als Taste. internal/tui/keylayout.go bringt cmdKey: es beantwortet, welche Taste auf einer US-PC-101-Tastatur gedrueckt worden waere, und alle vier Stellen, die vorher k.String() in einen Kommando-Switch gaben (model.go:877, home.go 215/230/237/258, edit.go:82), fragen jetzt cmdKey
  Die Antwort kommt aus zwei Quellen: Key.BaseCode, das ein Terminal mit Kitty-Protokoll selbst meldet und das fuer Buchstaben jede Belegung abdeckt, und sonst der Tabelle usPosition mit der kyrillischen JZUKEN-Belegung. Windows Terminal - wo der Melder sitzt - meldet kein BaseCode, dort traegt die Tabelle
  Vier Guards halten die Umsetzung von allem fern, wo sie schaden wuerde: benannte Tasten (enter, space, Pfeile), Satzzeichen mit Shift (shift+/ druckt '?' und meldet trotzdem '/'), gehaltene Modifier inklusive AltGr, und - als Ausnahme davon - die Lock-Zustaende, die Kitty bei jedem Tastendruck mitschickt
  Texteingabe laeuft nicht durch cmdKey: Filter, Titel-Eingabe und Editor lesen weiter k.Text, russischer Text bleibt russisch
  Dazu in beiden View()-Funktionen die KeyboardEnhancements, ohne die BaseCode gar nicht erst ankommt, und eine Zeile in core/release/NOTES.md
review-gaps: |-
  none - vier Durchgaenge eines zweiten Modells, drei davon mit Fund (Shift-Satzzeichen, AltGr, NumLock), alle behoben und je mit einem Test festgenagelt, der ohne seinen Guard umfaellt. Der letzte Durchgang blockiert nichts mehr
  Ehrlich offen, aber kein Mangel am Diff: der echte Kitty-Pfad und der Windows-Console-Pfad sind nur ueber konstruierte KeyPressMsg gedeckt, weil hier kein solches Terminal laeuft. Beide gefundenen Fehler in diesen Pfaden kamen aus dem Lesen des Decoders, nicht aus einem Test
test-verdict: |-
  go build, go vet, go test ./... gruen nach dem Lock-Fix
  Acht Tabellentests und zwei Board-Tests; vier Guards, vier Faelle, die ohne ihren Guard umfallen
  Dreimal von Hand bestaetigt: russische Belegung steuert das Board, '?' oeffnet die Hilfe (Windows Terminal, WSL2)
  Unveraendert offen und nur ueber konstruierte KeyPressMsg gedeckt: der echte Kitty-Pfad und der Windows-Console-Pfad. Der NumLock-Fehler kam genau von dort, gefunden hat ihn das Lesen des Decoders, nicht ein Test - wer an einem kitty, ghostty, WezTerm oder foot sitzt, sollte es einmal mit eingeschaltetem NumLock ausprobieren
  Binary neu gebaut unter /home/alex/.local/bin/jaira
question: ""
review-verdict: "Erfuellt die Definition of Done und ist merge-reif: der Melder hat es viermal an seinem Terminal bestaetigt, das zweite Modell blockiert nichts mehr"
review-check: |-
  1. go build -o /tmp/jaira-check ./cmd/jaira
  2. /tmp/jaira-check board - das Board geht auf
  3. Tastatur auf Russisch umstellen
  4. о und л druecken - der Cursor geht im Lane runter und rauf (das sind die Tasten j und k)
  5. р und д druecken - der Cursor springt in die Lane links und rechts
  6. Punkt-Taste druecken - unten geht der Filter auf
  7. русский ins Filterfeld tippen - es erscheint russisch, nicht als Kommando
  8. esc, dann ? druecken - die Hilfe geht auf, nicht der Filter
  9. esc, dann й - das Board schliesst
  10. go test ./internal/tui/ -run TestCmdKey -v - acht Tests, alle PASS
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
- **2026-09-13 18:40 · Alexander Sacharov** — Mensch hat den dritten Stand bestaetigt: Board reagiert mit russischer Belegung, Windows Terminal unter WSL2.
- **2026-09-13 18:43 · Alexander Sacharov** — review 3: NumLock. Der Guard war eine Zeile zu eng - ModNumLock kommt auf Kitty-Terminals bei jedem Tastendruck mit und wird von der Bibliothek selbst als 'beeinflusst den Text nicht' behandelt (decoder.go:1475). Mit NumLock an haette der Guard die ganze Funktion auf genau den Terminals ausgeschaltet, fuer die die BaseCode-Ebene gebaut wurde, und niemand haette es hier gemerkt. Lehre fuer den Kommentar: Lock-Zustaende sind keine Modifier.
- **2026-09-13 18:51 · Alexander Sacharov** — Mensch hat den vierten Stand gegengeprueft.
- **2026-09-13 19:06 · Alexander Sacharov** — Pull Request offen: https://github.com/BeMuCa/jaira/pull/17 (sashasoft90:feat/Z4W2BX-keys-layout-independent -> BeMuCa/jaira master, 14 Commits). Beschreibung steht aus dem Ticket, nicht aus dem Diff: Kontext, Ergebnis, Pruefanleitung und der bewusst nicht erledigte Teil (keine Tabellen fuer Griechisch/Hebraeisch/Arabisch, ctrl-Kombinationen bleiben Sache des Terminals, Satzzeichen lesen kein BaseCode).
