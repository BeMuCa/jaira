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
updated-at: 2026-10-05T13:38:33Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-323571
claimed-at: 2026-10-05T12:48:54Z
executed-by: opus
outcome-what: "Critique Runde 1: Schluessel mit Leerzeichen vor dem Wert bleibt eine Bedingung, leere Alternativen zaehlen nicht, Quotes schuetzen Kommas, Milestone-Picker besitzt nur milestone:, Hinweis im L-Fenster ohne %q, Doku (COMMANDS.md, README) nennt u und die Syntax"
outcome-why: "Kritik: 'tag: ui' fand nichts mehr, 'title:zzz,' liess alles durch, 'Doe, John' traf niemanden, M loeschte den User-Filter"
outcome-resolves: "Je Befund ein Test (SpaceAfterAColon, EmptyAlternative, QuotesKeepACommaInAValue, SplitTermKeepsQuotedCommas, MilestonePickerKeepsTheUserFilter), jeder per Mutation als scharf bewiesen; go test ./... -race RC=0"
review-summary: "Form passt (users.go folgt der Milestone-Auswahl, ein /-Filter, Tests gruen). Falsch: (1) Zerlegen an Leerzeichen macht 'tag: ui', 'milestone: q4-cleanup', 'tag:needs review' leer, die vorher trafen - 'verbreitert nur' stimmt nicht (Note, Testname, Kommentar). (2) Leerer Komma-Teil wie 'title:zzz,' laesst jedes Ticket durch. (3) Anfuehrungszeichen schuetzen kein Komma - 'Nachname, Vorname' aus dem Picker trifft niemanden. Optional: links.go:147 druckt den Filter noch mit %q; M ersetzt den ganzen Filter statt nur milestone:; Picker scrollt nicht; Namen nur bei ' \\t' gequotet; docs/COMMANDS.md + README nennen user/Kombination nicht."
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
- **2026-10-05 13:38 · BeMuCa** — Critique Runde 1 umgesetzt: (1) 'tag: ui' - ein Schluessel ohne Wert nimmt das naechste Wort (filterTerms klebt an ':'). 'tag:needs review' braucht jetzt Quotes - 'verbreitert nur' gilt nur fuer reine Phrasen; Kommentar, Testname und NOTES korrigiert. (2) leere Komma-Teile zaehlen nicht (filterValues), ein Schluessel ganz ohne Wert laesst das Board wie es ist. (3) ein splitOutsideQuotes fuer Leerzeichen UND Komma - Quotes schuetzen jetzt auch Kommas ('Doe, John'); Picker quotet bei Leerraum oder Komma. (4) links.go-Hinweis druckt den Filter wie getippt. (5) M und x im Milestone-Fenster besitzen nur noch milestone: (splitTerm, gemeinsam mit dem User-Picker) - sonst loeschten sich die beiden Picker gegenseitig. Nicht gemacht: Scrollen im Picker (echte Boards haben 2-3 Namen; Milestone-Fenster hat dieselbe Grenze).
