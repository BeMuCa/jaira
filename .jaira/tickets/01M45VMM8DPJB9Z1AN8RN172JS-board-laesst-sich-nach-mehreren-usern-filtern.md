---
id: 01M45VMM8DPJB9Z1AN8RN172JS
title: Board laesst sich nach mehreren Usern filtern
status: critique
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Per Taste eine Liste der User oeffnen, mehrere ankreuzen, und das Board zeigt nur deren Tickets - kombinierbar mit dem /-Filter"
context: |-
  Heute nimmt der /-Filter genau eine Bedingung (matches, internal/tui/model.go:642).
  Beispiel: 'assignee:berk 7MG5GB' wird als 'assignee enthaelt berk 7mg5gb' gelesen und trifft nichts.
  Berk will (05.10.2026): Taste -> Liste der User -> mehrere ankreuzen -> Board zeigt nur deren Tickets.
  Entschieden von Berk am 05.10.: User = Assignee ODER Creator. Mehrfachauswahl, nicht nur einer.
  Berks Frage dazu: sollte / nicht zwei Dinge gleichzeitig filtern koennen, z.B. Assignee und Ticket-Nr? -> gehoert in dieses Ticket.
  Vorbild fuer 'Geste schreibt in den Filter': die Milestone-Auswahl setzt filter='milestone:<name>' (internal/tui/model.go:1129).
  Offen: Syntax fuer mehrere Bedingungen im /-Filter; schreibt die Auswahl in den /-Filter oder ist sie ein eigener Filter daneben; welche Taste.
definition-of-done: "Eine Taste oeffnet eine Liste aller Assignees und Creators der Tickets auf dem Board; space kreuzt an, enter wendet an; Test"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-05T11:02:02Z
updated-at: 2026-10-05T13:14:47Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-323571
claimed-at: 2026-10-05T12:48:54Z
executed-by: opus
outcome-what: "Taste u oeffnet ein Users-Fenster (internal/tui/users.go): Namen aus assignee+creator aller geladenen Tickets mit Anzahl, space kreuzt an, enter schreibt user:a,b in den /-Filter, x entfernt, esc verwirft. Der /-Filter (matches in model.go) trennt Bedingungen an Leerzeichen (UND), Komma = oder je Feld, Anfuehrungszeichen halten Woerter zusammen; neuer Schluessel user: = assignee oder creator, exakt"
outcome-why: Berk will das Board per Taste auf mehrere Personen einschraenken und im /-Filter zwei Dinge zugleich filtern koennen (z.B. Person und Ticket-Nr)
outcome-resolves: "DoD 1-2: Picker-Tests und live in tmux; DoD 3: vier Filter-Tests fuer UND/ODER/Quotes/Phrase; DoD 4: Filter in der Kopfzeile wie typed, esc auf dem Board leert ihn (Tests); DoD 5: NOTES.md, -race RC=0. Jeder neue Test per Mutation scharf (Quoting, %q, Quote-Entfernung)"
---

# Board laesst sich nach mehreren Usern filtern

## Definition of Done

- [x] Eine Taste oeffnet eine Liste aller Assignees und Creators der Tickets auf dem Board; space kreuzt an, enter wendet an; Test
  proof: TestUOpensTheUserPicker, TestUserPickerNarrowsTheBoardToTheTickedPeople (internal/tui/userpicker_test.go); internal/tui/users.go openUsers/keyUsers; live tmux: u -> Users-Fenster mit tester 6, Alexander Sacharov 1
- [x] Board zeigt nur Tickets, deren Assignee oder Creator angekreuzt ist; Test mit zwei angekreuzten Usern
  proof: TestUserPickerNarrowsTheBoardToTheTickedPeople (zwei angekreuzt), TestFilterUserIsAssigneeOrCreator (internal/tui/filter_test.go); live: nur 'Alexanders card' sichtbar
- [x] Der /-Filter kombiniert mehrere Bedingungen, z.B. User und Ticket-Nr gleichzeitig; Test
  proof: TestFilterTermsAllHaveToMatch, TestFilterCommaMeansOr, TestFilterQuotesKeepWordsTogether, TestFilterPhraseNeverFindsLess (internal/tui/filter_test.go); matches/filterTerms/matchField internal/tui/model.go
- [x] Statusleiste zeigt den aktiven User-Filter; esc hebt ihn auf; Test
  proof: Filter steht wie bisher in der Kopfzeile des Boards (nicht in der Statusleiste - dort stand er nie): TestHeaderShowsTheFilterAsWritten; esc: TestEscOnTheBoardClearsTheUserFilter (internal/tui/userpicker_test.go); live bestaetigt
- [x] core/release/NOTES.md hat eine Zeile unter ## Unreleased; go test ./... -race gruen
  proof: core/release/NOTES.md 2 Zeilen unter ## Unreleased; go test ./... -race RC=0

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] matches(): Leerzeichen trennt Bedingungen (alle muessen passen), Komma im Wert = oder; neuer Schluessel user: = assignee ODER creator, exakter Name ohne Gross/Klein
- [x] failing tests fuer matches: UND, ODER, user trifft creator, Phrase findet nie weniger als vorher, bestehende filter_test gruen
- [x] Taste u: Fenster 'Users' (modeUsers) wie die Milestone-Auswahl - Namen aus assignee+creator aller geladenen Karten mit Anzahl; space kreuzt an, enter wendet an, x entfernt, esc schliesst
- [x] Auswahl schreibt nur den user:-Teil in den /-Filter, andere Bedingungen bleiben; beim Oeffnen sind die Namen aus dem Filter angekreuzt
- [x] Hilfe (?) nennt u und die neue /-Syntax; NOTES.md zwei Zeilen
- [x] go test ./... -race, gofmt, Binary bauen, live im tmux pruefen

## Progress
- **2026-10-05 12:49 · BeMuCa** — Pre-process: Vorbild ist die Milestone-Auswahl (model.go ~1110, view.go renderMilestones): Geste schreibt den normalen Filter, damit esc und / weiter genau einen Filter kennen. u ist auf dem Board frei (grep 'case "u"' leer). Kein Eintrag in der Statusleiste - M steht dort auch nicht, und jede Taste mehr laesst die Leiste frueher umbrechen; nur in der Hilfe. Phrase->Woerter verbreitert nur: was 'export csv' als Phrase traf, enthaelt beide Woerter.
- **2026-10-05 13:14 · BeMuCa** — In-progress: Namen mit Leerzeichen sind echt - 'Alexander Sacharov' ist auf dem jAIra-Board 49x creator. 'Leerzeichen trennt Bedingungen' haette ihn zerlegt; deshalb filterTerms: doppelte Anfuehrungszeichen halten Woerter zusammen (user:"Alexander Sacharov"), bringen nebenbei die exakte Phrasensuche zurueck. Milestone- und Tag-Namen sind kebab (core/milestone NormalizeName -> tag.Normalize), brauchen keine Quotes.
Kopfzeile zeigte den Filter mit %q - mit Quotes im Filter kam '\"' heraus; jetzt wie getippt (view.go header).
Komma in einem Namen wird nicht unterstuetzt (Komma = oder).
