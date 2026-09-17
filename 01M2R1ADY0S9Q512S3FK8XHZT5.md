---
id: 01M2R1ADY0S9Q512S3FK8XHZT5
title: Einfuegen aus der Zwischenablage kommt in keinem Eingabefeld an
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer etwas in die Suche der Tafel einfuegt, sieht es im Feld stehen und die Tafel filtert danach - egal auf welcher Tastaturbelegung und egal ob der Text Kyrillisch, Umlaute oder Emoji enthaelt"
context: |-
  Alex am 2026-09-17: 'in den Suchfilter laesst sich nichts einfuegen'. Strg+C/Strg+V in der Suche tut nichts - der Text aus der Zwischenablage kommt nicht an.

  Die Ursache steht fest, nachgesehen im Code:
  - internal/tui/model.go:897 - Model.Update hat genau einen Zweig fuer Tasten: 'case tea.KeyPressMsg'.
  - charm.land/bubbletea/v2@v2.0.8 input.go:37 - ein Einfuegen kommt NICHT als Taste, sondern als eigenes 'tea.PasteMsg'.
  - Niemand faengt dieses PasteMsg. Es faellt durch das Ende von Update und ist weg.

  Das trifft nicht nur die Suche. Jeder Modus, der Text sammelt, liest ausschliesslich k.Text: modeFilter (model.go:932), modeCreate (957), modeDelete (975) und modeEdit (edit.go:136). In allen vieren wird ein Einfuegen still verschluckt.

  Schon geklaert, nicht nochmal untersuchen: Strg+V laesst sich NICHT ueber die Tastaturbelegung abbilden. internal/tui/keylayout.go:38-42 sagt warum - der Dekoder loescht Key.Text, sobald ein Modifikator jenseits von Shift gedrueckt ist, es bleibt kein Zeichen uebrig, das cmdKey verschieben koennte. Das ist auch nicht noetig: ein Einfuegen meldet das Terminal als Ereignis, nicht als Tastenkombination, also ist die Behebung von der Belegung unabhaengig - kyrillische und deutsche Belegung bekommen sie geschenkt, sobald PasteMsg behandelt wird.

  Zweiter Punkt, der beim Anfassen nicht kaputtgehen darf: der Puffer ist rune-basiert. m.input += k.Text und das Backspace ueber []rune (model.go:928) tragen heute Umlaute und Kyrillisch heil durch - ein eingefuegter Text muss genauso ankommen, byteweise abgeschnitten waere ein Rueckschritt.

  Strg+C selbst bleibt, was es ist: ein Beenden (model.go:999, 1282, 1297). Kopieren AUS der Tafel heraus ist nicht Teil dieses Tickets - der Terminal-Emulator macht das mit seiner eigenen Auswahl.
definition-of-done: "Ein Einfuegen in die Suche landet im Feld: Model.Update behandelt tea.PasteMsg, der Text haengt an m.input, und die Tafel filtert sofort danach wie beim Tippen. Mit Test nachgewiesen, der ein PasteMsg direkt einspeist."
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-09-17T15:56:15Z
updated-at: 2026-09-17T16:04:32Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-45551
claimed-at: 2026-09-17T15:57:28Z
---

# Einfuegen aus der Zwischenablage kommt in keinem Eingabefeld an

## Definition of Done

- [x] Ein Einfuegen in die Suche landet im Feld: Model.Update behandelt tea.PasteMsg, der Text haengt an m.input, und die Tafel filtert sofort danach wie beim Tippen. Mit Test nachgewiesen, der ein PasteMsg direkt einspeist.
  proof: internal/tui/model.go:900 case tea.PasteMsg -> internal/tui/paste.go:24 Model.paste; TestPasteIntoFilterLandsAndFilters
- [x] Kein Eingabefeld verschluckt ein Einfuegen mehr: modeFilter, modeCreate, modeDelete und modeEdit nehmen den eingefuegten Text genauso an wie getippten. Je ein Test pro Modus.
  proof: internal/tui/paste.go:25-41 covers modeEdit, modeFilter, modeCreate, modeDelete; TestPasteReachesEveryInputMode (create/delete/edit subtests) + TestPasteIntoFilterLandsAndFilters
- [x] Mehrbyte-Text ueberlebt das Einfuegen unveraendert: kyrillischer Text, Umlaute und ein Emoji stehen nach dem Einfuegen Zeichen fuer Zeichen im Puffer, und ein Backspace danach entfernt genau ein Zeichen, nicht ein Byte. Mit Test nachgewiesen.
  proof: TestPasteKeepsMultiByteTextWhole in internal/tui/paste_test.go — 'Gruesse Privet <emoji>' arrives whole and one backspace removes one rune
- [x] Ein mehrzeiliger eingefuegter Text zerlegt die Suche nicht: das Feld ist eine Zeile, also entscheidet das Ticket bewusst, was mit Zeilenumbruechen passiert (verwerfen oder zu Leerzeichen falten), und ein Test haelt diese Entscheidung fest.
  proof: internal/tui/paste.go:47-63 sanitisePaste folds newline runs to one space; TestMultiLinePasteIsFoldedToSpaces (5 cases) and TestPasteKeepsLinesInTheFieldEditor
- [x] Die Tastaturbelegung spielt keine Rolle: eine Notiz am Code haelt fest, warum Strg+V nicht ueber cmdKey abgebildet wird (keylayout.go:38-42, der Dekoder loescht Key.Text bei Modifikatoren) und warum das Behandeln von PasteMsg die Belegung ueberfluessig macht.
  proof: internal/tui/paste.go:6-23 doc comment on Model.paste, and internal/tui/model.go:899 at the branch
- [x] core/release/NOTES.md traegt unter '## Unreleased' eine Zeile: was der Benutzer jetzt TUN kann - in die Suche einfuegen.
  proof: core/release/NOTES.md:17 under '## Unreleased'

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] sanitisePaste in model.go: strip CR, fold newline runs to one space, trim ends - shared by the single-line fields
- [x] Update: case tea.PasteMsg -> m.paste(msg.Content), routed per mode like the k.Text default branch
- [x] modeFilter/modeCreate/modeDelete take the folded text; modeEdit keeps newlines (enter inserts one there)
- [x] comment at the PasteMsg branch: why ctrl+v cannot go through cmdKey (keylayout.go:38-42)
- [x] paste_test.go: one test per mode, one for multi-byte + backspace, one for multiline folding
- [x] NOTES.md line under ## Unreleased

## Progress
- **2026-09-17 15:57 · Alexander Sacharov** — Entscheidung zu DoD 4, vom Dispatcher getroffen nachdem Alex sie ihm ueberlassen hat: mehrzeiliger eingefuegter Text wird zu Leerzeichen gefaltet, nicht verworfen. Grund: wer eine Zeile aus einem Terminal oder einer Datei kopiert, nimmt den abschliessenden Umbruch versehentlich mit - verwerfen wuerde dann den ganzen Text schlucken und wie der Fehler aussehen, den dieses Ticket behebt. Aufeinanderfolgende Umbrueche werden zu einem Leerzeichen, fuehrende und abschliessende fallen weg.
- **2026-09-17 16:04 · Alexander Sacharov** — Die Behebung sitzt bewusst NICHT im Tastenpfad. Ein Einfuegen kommt als tea.PasteMsg an, nie als Taste - deshalb ein eigener Zweig in Model.Update neben tea.KeyPressMsg statt einer Erweiterung von m.key(). Ueber cmdKey waere es nicht gegangen (keylayout.go: Key.Text ist bei Modifikatoren jenseits von Shift leer), und es ist auch nicht noetig - die Belegung spielt fuer ein Terminal-Ereignis keine Rolle.

modeEdit faellt aus der Faltung heraus: dort fuegt Enter eine Zeile ein (edit.go:83), der Puffer ist also mehrzeilig, und ein eingefuegter Absatz behaelt dort seine Form. Nur die einzeiligen Felder (Filter, Titel, Loeschbestaetigung) laufen durch sanitisePaste.

Erste Fassung von sanitisePaste hat jede Zeile einzeln getrimmt - verworfen, weil ein einzeiliges Einfuegen mit fuehrendem Leerzeichen dann veraendert ankommt. Jetzt nur noch: CR normalisieren, Umbrueche an den Enden weg, Umbruchketten zu einem Leerzeichen. Innenliegende Leerzeichen bleiben unangetastet.

Gegenprobe gemacht: mit auskommentiertem PasteMsg-Zweig fallen 14 der neuen Faelle um, mit dem Zweig laeuft die ganze Suite gruen.
