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
updated-at: 2026-09-17T16:00:52Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-45551
claimed-at: 2026-09-17T15:57:28Z
---

# Einfuegen aus der Zwischenablage kommt in keinem Eingabefeld an

## Definition of Done

- [ ] Ein Einfuegen in die Suche landet im Feld: Model.Update behandelt tea.PasteMsg, der Text haengt an m.input, und die Tafel filtert sofort danach wie beim Tippen. Mit Test nachgewiesen, der ein PasteMsg direkt einspeist.
- [ ] Kein Eingabefeld verschluckt ein Einfuegen mehr: modeFilter, modeCreate, modeDelete und modeEdit nehmen den eingefuegten Text genauso an wie getippten. Je ein Test pro Modus.
- [ ] Mehrbyte-Text ueberlebt das Einfuegen unveraendert: kyrillischer Text, Umlaute und ein Emoji stehen nach dem Einfuegen Zeichen fuer Zeichen im Puffer, und ein Backspace danach entfernt genau ein Zeichen, nicht ein Byte. Mit Test nachgewiesen.
- [ ] Ein mehrzeiliger eingefuegter Text zerlegt die Suche nicht: das Feld ist eine Zeile, also entscheidet das Ticket bewusst, was mit Zeilenumbruechen passiert (verwerfen oder zu Leerzeichen falten), und ein Test haelt diese Entscheidung fest.
- [ ] Die Tastaturbelegung spielt keine Rolle: eine Notiz am Code haelt fest, warum Strg+V nicht ueber cmdKey abgebildet wird (keylayout.go:38-42, der Dekoder loescht Key.Text bei Modifikatoren) und warum das Behandeln von PasteMsg die Belegung ueberfluessig macht.
- [ ] core/release/NOTES.md traegt unter '## Unreleased' eine Zeile: was der Benutzer jetzt TUN kann - in die Suche einfuegen.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-17 15:57 · Alexander Sacharov** — Entscheidung zu DoD 4, vom Dispatcher getroffen nachdem Alex sie ihm ueberlassen hat: mehrzeiliger eingefuegter Text wird zu Leerzeichen gefaltet, nicht verworfen. Grund: wer eine Zeile aus einem Terminal oder einer Datei kopiert, nimmt den abschliessenden Umbruch versehentlich mit - verwerfen wuerde dann den ganzen Text schlucken und wie der Fehler aussehen, den dieses Ticket behebt. Aufeinanderfolgende Umbrueche werden zu einem Leerzeichen, fuehrende und abschliessende fallen weg.
