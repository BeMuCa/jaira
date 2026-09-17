---
id: 01M2R1ADY0S9Q512S3FK8XHZT5
title: Einfuegen aus der Zwischenablage kommt in keinem Eingabefeld an
status: testing
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
updated-at: 2026-09-17T17:23:24Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-6497
claimed-at: 2026-09-17T17:01:52Z
outcome-what: "Removed the forwarding wrapper Model.paste; the PasteMsg branch calls insertText directly and carries the rationale"
outcome-why: "A wrapper with one caller that only forwards is one indirection between the event and the code that handles it, and it kept the explanation in a different file from the branch it explains"
outcome-resolves: "optimize pass: no duplication, no dead code, no behaviour change"
review-summary: "none"
review-gaps: "removed Model.paste (internal/tui/paste.go) — a wrapper with one caller that only forwarded to insertText; its rationale now sits at the 'case tea.PasteMsg' branch in model.go and the tea import went with it. Left alone: foldToOneLine has no duplicate in the repo (view.go wrap* folds the other way, edit.go:178 is display-only, core/lane/corrections.go:219 is file reading in another package), and the two ReplaceAll on the per-keystroke path allocate nothing when there is no match. No dead code and no behaviour change; tests green, go vet clean."
test-verdict: "pass: suite green (go build/vet/test RC=0, -race on internal/tui RC=0, Windows vet+build RC=0), DoD 1-6 verified in the working tree, and a real bracketed paste (ESC[200~ … ESC[201~) fed to the binary in a pty lands in the filter with Cyrillic/umlauts/emoji intact and folds a multi-line paste to 'paste bug', narrowing the board"
---

# Einfuegen aus der Zwischenablage kommt in keinem Eingabefeld an

## Definition of Done

- [x] Ein Einfuegen in die Suche landet im Feld: Model.Update behandelt tea.PasteMsg, der Text haengt an m.input, und die Tafel filtert sofort danach wie beim Tippen. Mit Test nachgewiesen, der ein PasteMsg direkt einspeist.
  proof: internal/tui/model.go:911 case tea.PasteMsg -> internal/tui/model.go:924 Model.insertText; TestPasteIntoFilterLandsAndFilters (paste_test.go feeds a real tea.PasteMsg)
- [x] Kein Eingabefeld verschluckt ein Einfuegen mehr: modeFilter, modeCreate, modeDelete und modeEdit nehmen den eingefuegten Text genauso an wie getippten. Je ein Test pro Modus.
  proof: internal/tui/model.go:914-939 insertText covers modeEdit, modeFilter, modeCreate, modeDelete and is the single path for typed and pasted text alike; TestPasteReachesEveryInputMode (create/delete/edit subtests) + TestPasteIntoFilterLandsAndFilters
- [x] Mehrbyte-Text ueberlebt das Einfuegen unveraendert: kyrillischer Text, Umlaute und ein Emoji stehen nach dem Einfuegen Zeichen fuer Zeichen im Puffer, und ein Backspace danach entfernt genau ein Zeichen, nicht ein Byte. Mit Test nachgewiesen.
  proof: TestPasteKeepsMultiByteTextWhole in internal/tui/paste_test.go — 'Gruesse Privet <emoji>' arrives whole and one backspace removes one rune
- [x] Ein mehrzeiliger eingefuegter Text zerlegt die Suche nicht: das Feld ist eine Zeile, also entscheidet das Ticket bewusst, was mit Zeilenumbruechen passiert (verwerfen oder zu Leerzeichen falten), und ein Test haelt diese Entscheidung fest.
  proof: internal/tui/model.go:921-927 insertText normalises CRLF and lone CR to \n once; internal/tui/paste.go:41 foldToOneLine folds newline runs to one space in a single pass; TestMultiLinePasteIsFoldedToSpaces (6 cases), TestPasteKeepsLinesInTheFieldEditor, TestPasteNormalisesALoneCarriageReturnInTheFieldEditor
- [x] Die Tastaturbelegung spielt keine Rolle: eine Notiz am Code haelt fest, warum Strg+V nicht ueber cmdKey abgebildet wird (keylayout.go:38-42, der Dekoder loescht Key.Text bei Modifikatoren) und warum das Behandeln von PasteMsg die Belegung ueberfluessig macht.
  proof: internal/tui/model.go:900-910 comment at the 'case tea.PasteMsg' branch: why ctrl+v cannot go through cmdKey (keylayout.go clears Key.Text under modifiers) and why the event makes the layout irrelevant
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
- [x] critique gap (2. Durchlauf): \r\n und einzelnes \r einmal oben in insertText zu \n normalisieren; modeEdit haengt dann nur noch an, foldToOneLine verliert seine ReplaceAll-Zeilen

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
- **2026-09-17 17:07 · Alexander Sacharov** — critique (2. Durchlauf): Die drei Punkte des ersten Durchlaufs sind sauber abgearbeitet - insertText ist der eine Weg in die Puffer (model.go:914), key()s drei default-Zweige und edit.go:136 rufen dieselbe Funktion, foldToOneLine faltet in einem Durchlauf. Geprueft und nicht beanstandet: Model.paste ist zwar ein Einzeiler mit einem Aufrufer, spiegelt aber genau m.key(msg) im selben switch - das ist das vorhandene Muster, kein neues. Die Leerpruefung in modeFilter steht dort zu Recht: ohne sie laeuft rebuild() ueber die ganze Tafel fuer einen Einfuegevorgang, der nichts aendert. Offen bleibt genau ein Punkt: in insertText gibt es zwei Definitionen von 'Zeilenumbruch'. modeEdit (model.go:922) kennt nur \r\n, foldToOneLine (paste.go:40) kennt \r\n und \r. Terminals schicken in einer Klammer-Einfuegung durchaus einzelne CR; die landen dann roh im Editorpuffer und damit in der Ticketdatei. paste_test.go deckt den Fall nicht ab - TestPasteKeepsLinesInTheFieldEditor testet nur \r\n. Fix: einmal oben in insertText normalisieren, dann haengt modeEdit nur noch an und foldToOneLine kommt ohne seine ReplaceAll-Zeilen aus.
- **2026-09-17 17:10 · Alexander Sacharov** — Zweiter Critique-Punkt eingearbeitet: "Zeilenumbruch" ist jetzt einmal definiert. insertText (model.go:921-927) normalisiert \r\n und einzelnes \r zu \n, bevor irgendein Puffer den Text sieht; modeEdit haengt danach nur noch an, und foldToOneLine kommt ohne seine beiden ReplaceAll-Zeilen aus.

Warum oben und nicht in jedem Zweig: die beiden Definitionen sind genau so auseinandergelaufen - modeEdit kannte nur \r\n, foldToOneLine auch das einzelne \r. Ein Zweig mehr haette denselben Fehler wieder eingebaut.

Nachweis, dass es vorher wirklich kaputt war: TestPasteNormalisesALoneCarriageReturnInTheFieldEditor faellt mit dem alten ReplaceAll(text, "\r\n", "\n") im modeEdit-Zweig um - das einzelne \r landete roh im editBuf und damit in der Ticketdatei.

Keine NOTES.md-Zeile dafuer: die vorhandene Zeile unter ## Unreleased beschreibt das Einfuegen als Ganzes, und diese Runde aendert nichts, was ein Benutzer ausserhalb davon bemerken koennte.
- **2026-09-17 17:13 · Alexander Sacharov** — critique (3. Durchlauf): nichts mehr zu beanstanden, review-summary=none. Der Punkt des zweiten Durchlaufs ist abgearbeitet - insertText (model.go:921-927) normalisiert \r\n und einzelnes \r einmal zu \n, bevor irgendein Puffer den Text sieht; modeEdit haengt nur noch an, foldToOneLine (paste.go:41) traegt keine ReplaceAll-Zeile mehr und dokumentiert die Abhaengigkeit ('Carriage returns are already gone'). Eine Definition von Zeilenumbruch, an einer Stelle.

Geprueft und ohne Befund, damit es kein vierter Durchlauf nochmal aufmacht:
- Vollstaendigkeit: 'grep k.Text' ueber internal/tui findet genau vier Aufrufstellen (model.go:975, 996, 1012, edit.go:136), alle rufen insertText. Kein Puffer liest k.Text noch direkt, kein Feld fehlt.
- Keine Verhaltensaenderung im Tastenpfad: 686daa5:model.go gegen heute verglichen - modeFilter setzte schon vorher m.filter und rief rebuild() beim Tippen, modeCreate/modeDelete haengten nur an. insertText bildet beides unveraendert ab, nur die 'if k.Text != ""' sind in die eine Leerpruefung oben gewandert.
- Die Behebung ist nicht tot: charm.land/bubbletea/v2@v2.0.8 cursed_renderer.go:115 schaltet Bracketed Paste von sich aus ein (DisableBracketedPasteMode ist der Ausschalter, niemand setzt ihn). tea.NewProgram ohne Option genuegt also, es fehlt kein Aktivierungsschritt.
- Nicht wieder aufgemacht: Model.paste als Einzeiler mit einem Aufrufer (spiegelt m.key(msg) im selben switch) und die Leerpruefung im modeFilter-Zweig (spart ein rebuild() ueber die ganze Tafel) - beide im 2. Durchlauf bewusst stehen gelassen.
- editKey 'enter' haengt \n direkt an statt ueber insertText: identisches Ergebnis, der Umweg waere reine Indirektion.
- **2026-09-17 17:15 · Alexander Sacharov** — optimize: entfernt wurde der Wrapper Model.paste in internal/tui/paste.go - er leitete nur an insertText weiter und hatte genau einen Aufrufer. Die Begruendung, warum ein Einfuegen keine Tastenbindung sein kann, steht jetzt direkt am 'case tea.PasteMsg' in model.go, also dort, wo jemand sie liest. paste.go haelt nur noch foldToOneLine; der tea-Import ist damit weg. Geprueft und stehen gelassen: foldToOneLine hat kein Gegenstueck im Repo - view.go wrapLines/wrap gehen in die andere Richtung, edit.go:178 ersetzt Umbrueche fuer die Anzeige durch ' ⏎ ', core/lane/corrections.go:219 normalisiert CRLF beim Dateilesen und gehoert einem anderen Paket. Die zwei ReplaceAll in insertText laufen bei jedem Tastendruck, kosten aber nichts: strings.Replace gibt bei null Treffern denselben String ohne Allokation zurueck. Tests gruen, go vet sauber.
- **2026-09-17 17:23 · Alexander Sacharov** — testing: Bestanden, drei Durchgaenge.

1. Tore: go build ./... RC=0, go vet ./... RC=0, go test ./... RC=0 (ganze Suite), go test -race -count=1 ./internal/tui RC=0 (114s, kein Race). Windows-Gegenprobe wie in README 'Development': GOOS=windows GOARCH=amd64 go vet ./... und go build ./cmd/jaira, beide RC=0.

2. Die Forderung, Punkt fuer Punkt am Arbeitsbaum geprueft, nicht am Outcome-Text:
- DoD 1: model.go:911 'case tea.PasteMsg' -> model.go:924 insertText; TestPasteIntoFilterLandsAndFilters gruen. paste_test.go:14 pasteInto speist wirklich ein tea.PasteMsg ein, kein nachgebautes KeyPressMsg - der Nachweis, den DoD 1 verlangt.
- DoD 2: grep k.Text ueber internal/tui findet vier Stellen (model.go:975, 996, 1012, edit.go:136), alle rufen insertText. TestPasteReachesEveryInputMode (create/delete/edit) + TestPasteIntoFilterLandsAndFilters, alle gruen.
- DoD 3: TestPasteKeepsMultiByteTextWhole gruen, 'Gruesse Privet Emoji' kommt heil an, ein Backspace entfernt ein Zeichen.
- DoD 4: TestMultiLinePasteIsFoldedToSpaces, 6 Faelle, alle gruen; dazu TestPasteKeepsLinesInTheFieldEditor und TestPasteNormalisesALoneCarriageReturnInTheFieldEditor.
- DoD 5: Begruendung steht an model.go:900-910 am Zweig selbst.
- DoD 6: core/release/NOTES.md:17 unter '## Unreleased', eine Zeile, als Anweisung geschrieben.

3. Verhalten am echten Programm geprueft, nicht nur am Test - die Tests umgehen den Dekoder, indem sie das PasteMsg direkt einspeisen, also war offen, ob das Terminal ueberhaupt eine Klammer-Einfuegung schickt und ob bubbletea sie ohne Zusatzoption einschaltet. Binary gebaut, in einem echten pty gestartet, Testtafel mit zwei Tickets, '/' gedrueckt und die rohe Bytefolge ESC[200~ ... ESC[201~ hineingeschrieben:
- 'Privet ueoe Emoji' (Kyrillisch + Umlaute + Emoji): steht Zeichen fuer Zeichen im Feld, Tafel filtert sofort (0 Tickets).
- 'paste\nbug\r\n' (mehrzeilig, mit CRLF am Ende): im Feld steht 'paste bug', Tafel zeigt 1 Ticket. Faltung und CR-Normalisierung greifen auch auf dem echten Weg durch den Dekoder.
Damit ist bestaetigt, was der 3. Critique-Durchlauf nur aus dem Quelltext geschlossen hatte: Bracketed Paste ist ohne Zusatzoption aktiv, die Behebung ist nicht tot.

Ein Befund, ohne Ruecklauf, weil er keinen Code betrifft: die Proof-Zeilen von DoD 1 und DoD 5 zeigten noch auf 'internal/tui/paste.go:24 Model.paste' und 'paste.go:9-23', die der optimize-Durchgang entfernt hat. Zeigten also ins Leere. Beide hier auf die heutigen Stellen gesetzt (model.go:911/924 bzw. model.go:900-910). paste.go traegt nur noch foldToOneLine.
