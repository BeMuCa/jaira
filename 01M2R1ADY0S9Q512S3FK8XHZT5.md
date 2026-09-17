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
updated-at: 2026-09-17T17:05:17Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-6497
claimed-at: 2026-09-17T17:01:52Z
outcome-what: "Ein Einfuegen aus der Zwischenablage kommt jetzt in jedem Eingabefeld der Tafel an. Model.Update hat einen eigenen Zweig 'case tea.PasteMsg' neben tea.KeyPressMsg (internal/tui/model.go:900), der nach internal/tui/paste.go:24 Model.paste fuehrt; der bedient modeFilter, modeCreate, modeDelete und modeEdit. Die drei einzeiligen Felder laufen vorher durch sanitisePaste (paste.go:47-63): CR normalisiert, Umbrueche an den Enden weg, Umbruchketten zu einem Leerzeichen; modeEdit bleibt aussen vor, weil sein Puffer mehrzeilig ist. Neu: internal/tui/paste.go und paste_test.go mit 160 Zeilen Tests, eine Zeile in core/release/NOTES.md unter '## Unreleased'."
outcome-why: "Ein Einfuegen kommt in bubbletea v2 als tea.PasteMsg an, nie als Taste (input.go:37). Update hatte genau einen Tastenzweig, also fiel das PasteMsg durch das Ende der Funktion und war weg - still, in allen vier Feldern, nicht nur in der Suche. Ueber die Tastaturbelegung war das nicht zu loesen: cmdKey kann Strg+V nicht verschieben, weil der Dekoder Key.Text loescht, sobald ein Modifikator jenseits von Shift gedrueckt ist (keylayout.go:38-42). Es ist auch nicht noetig - ein Terminal-Ereignis hat keine Belegung, also bekommen kyrillische und deutsche Tastaturen die Behebung geschenkt."
outcome-resolves: "DoD 1 bis 6. Gegenprobe gemacht: mit auskommentiertem PasteMsg-Zweig fallen 14 der neuen Faelle um."
review-summary: |-
  internal/tui/paste.go:25-41 wiederholt die Zuordnung Modus->Puffer, die schon in internal/tui/model.go:936-940, 962-964, 980-982 und internal/tui/edit.go:136-138 steht - samt der Regel, dass modeFilter zusaetzlich m.filter setzt und rebuild() ruft; stattdessen eine Methode m.insertText(s string) in model.go neben key() anlegen, die den Modus-Switch und die Faltung einmal haelt, und sie aus den drei default-Zweigen von key(), aus dem Ende von editKey() und aus paste() rufen - getippter Text enthaelt nie einen Umbruch, die Faltung ist dort also ein No-op.
  internal/tui/paste.go:57-62: die Schleife 'for strings.Contains(text, "\n\n")' laeuft wiederholt ueber den ganzen Text; stattdessen einmal strings.Split(text, "\n"), leere Teile weglassen, mit " " joinen - das ersetzt Trim und Schleife durch einen Durchlauf.
  internal/tui/paste.go:38-41: das 'if s := sanitisePaste(text); s != ""' in modeCreate/modeDelete bewacht nichts, ein leeres Anhaengen aendert den Puffer nicht; dort direkt 'm.input += sanitisePaste(text)'.
---

# Einfuegen aus der Zwischenablage kommt in keinem Eingabefeld an

## Definition of Done

- [x] Ein Einfuegen in die Suche landet im Feld: Model.Update behandelt tea.PasteMsg, der Text haengt an m.input, und die Tafel filtert sofort danach wie beim Tippen. Mit Test nachgewiesen, der ein PasteMsg direkt einspeist.
  proof: internal/tui/model.go:902 case tea.PasteMsg -> internal/tui/paste.go:24 Model.paste -> internal/tui/model.go:914 Model.insertText; TestPasteIntoFilterLandsAndFilters
- [x] Kein Eingabefeld verschluckt ein Einfuegen mehr: modeFilter, modeCreate, modeDelete und modeEdit nehmen den eingefuegten Text genauso an wie getippten. Je ein Test pro Modus.
  proof: internal/tui/model.go:914-939 insertText covers modeEdit, modeFilter, modeCreate, modeDelete and is the single path for typed and pasted text alike; TestPasteReachesEveryInputMode (create/delete/edit subtests) + TestPasteIntoFilterLandsAndFilters
- [x] Mehrbyte-Text ueberlebt das Einfuegen unveraendert: kyrillischer Text, Umlaute und ein Emoji stehen nach dem Einfuegen Zeichen fuer Zeichen im Puffer, und ein Backspace danach entfernt genau ein Zeichen, nicht ein Byte. Mit Test nachgewiesen.
  proof: TestPasteKeepsMultiByteTextWhole in internal/tui/paste_test.go — 'Gruesse Privet <emoji>' arrives whole and one backspace removes one rune
- [x] Ein mehrzeiliger eingefuegter Text zerlegt die Suche nicht: das Feld ist eine Zeile, also entscheidet das Ticket bewusst, was mit Zeilenumbruechen passiert (verwerfen oder zu Leerzeichen falten), und ein Test haelt diese Entscheidung fest.
  proof: internal/tui/paste.go:39-49 foldToOneLine folds newline runs to one space in a single pass; TestMultiLinePasteIsFoldedToSpaces (5 cases) and TestPasteKeepsLinesInTheFieldEditor
- [x] Die Tastaturbelegung spielt keine Rolle: eine Notiz am Code haelt fest, warum Strg+V nicht ueber cmdKey abgebildet wird (keylayout.go:38-42, der Dekoder loescht Key.Text bei Modifikatoren) und warum das Behandeln von PasteMsg die Belegung ueberfluessig macht.
  proof: internal/tui/paste.go:9-23 doc comment on Model.paste, and internal/tui/model.go:900-901 at the branch
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
- [x] critique gap 1: insertText(text) in model.go holds the mode->buffer switch once; key()'s three default branches, editKey()'s tail and paste() all call it
- [x] critique gap 2: fold newlines in one pass (Split/drop empty/Join) instead of the repeated Contains loop
- [x] critique gap 3: drop the s != "" guards in modeCreate/modeDelete; one empty-check at the top of insertText
- [x] run the tui tests, unchanged, to show typed and pasted text both still land

## Progress
- **2026-09-17 15:57 · Alexander Sacharov** — Entscheidung zu DoD 4, vom Dispatcher getroffen nachdem Alex sie ihm ueberlassen hat: mehrzeiliger eingefuegter Text wird zu Leerzeichen gefaltet, nicht verworfen. Grund: wer eine Zeile aus einem Terminal oder einer Datei kopiert, nimmt den abschliessenden Umbruch versehentlich mit - verwerfen wuerde dann den ganzen Text schlucken und wie der Fehler aussehen, den dieses Ticket behebt. Aufeinanderfolgende Umbrueche werden zu einem Leerzeichen, fuehrende und abschliessende fallen weg.
- **2026-09-17 16:04 · Alexander Sacharov** — Die Behebung sitzt bewusst NICHT im Tastenpfad. Ein Einfuegen kommt als tea.PasteMsg an, nie als Taste - deshalb ein eigener Zweig in Model.Update neben tea.KeyPressMsg statt einer Erweiterung von m.key(). Ueber cmdKey waere es nicht gegangen (keylayout.go: Key.Text ist bei Modifikatoren jenseits von Shift leer), und es ist auch nicht noetig - die Belegung spielt fuer ein Terminal-Ereignis keine Rolle.

modeEdit faellt aus der Faltung heraus: dort fuegt Enter eine Zeile ein (edit.go:83), der Puffer ist also mehrzeilig, und ein eingefuegter Absatz behaelt dort seine Form. Nur die einzeiligen Felder (Filter, Titel, Loeschbestaetigung) laufen durch sanitisePaste.

Erste Fassung von sanitisePaste hat jede Zeile einzeln getrimmt - verworfen, weil ein einzeiliges Einfuegen mit fuehrendem Leerzeichen dann veraendert ankommt. Jetzt nur noch: CR normalisieren, Umbrueche an den Enden weg, Umbruchketten zu einem Leerzeichen. Innenliegende Leerzeichen bleiben unangetastet.

Gegenprobe gemacht: mit auskommentiertem PasteMsg-Zweig fallen 14 der neuen Faelle um, mit dem Zweig laeuft die ganze Suite gruen.
- **2026-09-17 16:23 · Alexander Sacharov** — critique: Die Form ist im Kern richtig - ein eigener Zweig fuer tea.PasteMsg statt eines Umwegs ueber m.key() ist die richtige Entscheidung und bleibt so. Gegenprobe gemacht: den PasteMsg in ein KeyPressMsg{Text: ...} umzubauen und durch m.key() zu schicken waere kuerzer, wuerde aber in den Nicht-Eingabe-Modi gefaehrlich - ein eingefuegtes 'q' liefe dann durch cmdKey in das Beenden. Nicht nochmal aufmachen.

Ebenfalls geprueft und in Ordnung: die vier Modi sind vollstaendig. Ausser modeFilter, modeCreate, modeDelete und modeEdit sammelt kein Bildschirm Text - browse.go, defaultboard.go, dropboard.go, lanes.go und settings.go bekommen von Home nur cmdKey(msg) als String und haben gar keinen Puffer (internal/tui/home.go:213ff). Es fehlt also kein Eingabefeld.

Was zurueckgeht, sind drei Punkte zur Form, alle in internal/tui/paste.go:
1. Die Zuordnung Modus->Puffer steht jetzt zweimal. paste.go:25-41 sagt noch einmal, was model.go:936-940, 962-964, 980-982 und edit.go:136-138 schon sagen - inklusive der Regel, dass modeFilter zusaetzlich m.filter setzt und rebuild() ruft. Wer spaeter ein fuenftes Textfeld ergaenzt oder diese Regel aendert, muss an paste.go denken. Eine Methode m.insertText(s string) in model.go neben key(), die den Modus-Switch samt Faltung einmal haelt und aus key(), editKey() und paste() gerufen wird, legt die Regel an eine Stelle; getippter Text enthaelt nie einen Umbruch, die Faltung stoert dort nicht.
2. sanitisePaste (paste.go:57-62) faltet die Umbruchketten mit einer Schleife aus wiederholtem ReplaceAll ueber den ganzen Text. Split auf '\n', leere Teile weglassen, mit Leerzeichen joinen macht dasselbe in einem Durchlauf und macht Trim und Schleife zusammen ueberfluessig.
3. Das 's != ""' in modeCreate/modeDelete (paste.go:38-41) bewacht nichts - ein leeres Anhaengen aendert den Puffer nicht. Im modeFilter-Zweig darf es bleiben, dort spart es ein rebuild().

Das Verhalten ist von allen drei Punkten nicht betroffen, die Tests in paste_test.go muessen unveraendert gruen bleiben.
- **2026-09-17 17:05 · Alexander Sacharov** — Kritik eingearbeitet: die Modus->Puffer-Zuordnung steht jetzt einmal in Model.insertText (model.go:914) statt viermal. Getippter und eingefuegter Text laufen durch dieselbe Funktion - key()s drei default-Zweige, das Ende von editKey() und paste() rufen sie alle. Damit kann ein Feld nicht mehr getippten Text nehmen und eingefuegten verschlucken, was genau der Fehler dieses Tickets war.

Warum insertText in model.go und nicht in paste.go: getippter Text ist der haeufigere Weg dorthin, paste.go traegt nur noch das PasteMsg-Eingangstor und die Faltung. sanitisePaste heisst jetzt foldToOneLine - der Name 'Paste' war falsch, seit auch getippter Text durchlaeuft (dort ein No-op, eine Taste liefert nie einen Umbruch).

Die leere Eingabe wird einmal oben in insertText abgefangen. Das ersetzt sowohl die drei 'if k.Text != ""' im Tastenpfad als auch die zwei wirkungslosen 'if s != ""' im Einfuegepfad. In modeFilter bleibt eine zweite Pruefung stehen, weil dort erst die Faltung leer werden kann ('\n' allein) und ein rebuild() fuer nichts sonst laufen wuerde.

Faltung jetzt in einem Durchlauf: Split auf \n, leere Teile weg, mit Leerzeichen joinen. Ersetzt Trim plus die Schleife 'for strings.Contains(text, "\n\n")', die den Text wiederholt abging. Verhalten identisch, die fuenf Faelle in TestMultiLinePasteIsFoldedToSpaces laufen unveraendert durch.

Keine Testaenderung noetig und bewusst keine gemacht: die 160 Zeilen paste_test.go pruefen Verhalten, nicht Struktur. Dass sie nach dem Umbau unveraendert gruen sind, ist der Nachweis, dass der Umbau nichts verschoben hat. go test ./... komplett gruen.
