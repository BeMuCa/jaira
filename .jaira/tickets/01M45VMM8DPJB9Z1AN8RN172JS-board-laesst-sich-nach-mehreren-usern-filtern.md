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
updated-at: 2026-10-07T09:27:51Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-3594284
claimed-at: 2026-10-07T09:18:49Z
executed-by: opus
outcome-what: "Die Tastenleiste unten auf dem Board nennt jetzt 'u users' hinter 't tags'"
outcome-why: "Berk fand die Taste u beim Ausprobieren nicht, weil sie nur in der Hilfe (?) stand"
outcome-resolves: "DoD 6: statusBar() in internal/tui/view.go fuehrt 'u users' in derselben Liste wie alle anderen Tasten; die Leiste bricht um statt zu kuerzen, also geht der Hinweis auch auf schmalen Terminals nicht verloren - TestNarrowBoardShowsAllKeysAndFits prueft das bei 30x24"
review-summary: "Runde 3: Befund aus Runde 2 behoben - 'fix: crash' und 'todo: write' finden ihre Tickets wieder wie auf master; 'tag: ui' und 'milestone: q4-cleanup' bleiben eine Bedingung; alle Proben wie beabsichtigt; go test ./internal/tui -race gruen. Picker-Scrollen bewusst offen. Nichts mehr offen."
review-gaps: "Drei Kopien von 'eine Bedingung im Filter ersetzen, input nachziehen, rebuild' (Milestone enter, Milestone x, applyUsers) in setFilterTerm(key, values...) gefaltet. In matchField die beiden strings.TrimSpace entfernt: filterValues trimmt jeden Wert, und um einen Schluessel steht kein Leerzeichen mehr - beide waren wirkungslos. Gelassen: renderUsers/keyUsers folgen dem bestehenden Picker-Muster der Milestones (ein gemeinsamer Picker waere Umbau); Neu-Zerlegen des Filters pro Ticket ist vernachlaessigbar; der known-Check per matchField(t, key, \"\", ms) vermeidet eine zweite Schluesselliste. Kein bestehender quote-bewusster Splitter im Repo."
test-verdict: "pass: go test ./... -race RC=0 (Cache geleert), gofmt/vet sauber; DoD 1-5 per Tests und live in tmux (160x30): u-Liste mit Anzahl, zwei angekreuzt inkl. 'Alexander Sacharov' und 'Doe, John' -> user:\"Alexander Sacharov\",\"Doe, John\", 2 Tickets; Haken beim Wiederoeffnen da; x/esc wie spezifiziert; /-Filter UND, ODER, Quotes, 'assignee: sam' geprueft. M-Picker nur per Test (Fixture ohne Milestones)."
question: "User-Filter so abnehmen? Taste u, mehrere ankreuzen, schreibt user:a,b in den /-Filter; /-Filter: Leerzeichen = und, Komma = oder, Anfuehrungszeichen halten zusammen. Verhaltensaenderung: ein Feldwert mit Leerzeichen braucht jetzt Quotes (tag:\"needs review\"); M ersetzt nur noch milestone:. Zum Ausprobieren: Branch feat/N172JS-user-filter bauen (dein ~/.local/bin/jaira ist noch PAP369)."
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
- [x] Die Tastenleiste unten auf dem Board nennt 'u users', damit man die Taste ohne '?' findet; Test
  proof: internal/tui/view.go:982 keys enthaelt 'u users'; TestNarrowBoardShowsAllKeysAndFits (internal/tui/hints_test.go:39) war rot, jetzt gruen; live tmux 160x30: Zeile 'enter open · … · t tags · u users · n new · …'; go test ./... -race RC=0

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
- [x] 'u users' in die Tastenleiste (statusBar in internal/tui/view.go) neben 't tags'; TestNarrowBoardShowsAllKeysAndFits pinnt es

## Progress
- **2026-10-05 12:49 · BeMuCa** — Pre-process: Vorbild ist die Milestone-Auswahl (model.go ~1110, view.go renderMilestones): Geste schreibt den normalen Filter, damit esc und / weiter genau einen Filter kennen. u ist auf dem Board frei (grep 'case "u"' leer). Kein Eintrag in der Statusleiste - M steht dort auch nicht, und jede Taste mehr laesst die Leiste frueher umbrechen; nur in der Hilfe. Phrase->Woerter verbreitert nur: was 'export csv' als Phrase traf, enthaelt beide Woerter.
- **2026-10-05 13:14 · BeMuCa** — In-progress: Namen mit Leerzeichen sind echt - 'Alexander Sacharov' ist auf dem jAIra-Board 49x creator. 'Leerzeichen trennt Bedingungen' haette ihn zerlegt; deshalb filterTerms: doppelte Anfuehrungszeichen halten Woerter zusammen (user:"Alexander Sacharov"), bringen nebenbei die exakte Phrasensuche zurueck. Milestone- und Tag-Namen sind kebab (core/milestone NormalizeName -> tag.Normalize), brauchen keine Quotes.
Kopfzeile zeigte den Filter mit %q - mit Quotes im Filter kam '\"' heraus; jetzt wie getippt (view.go header).
Komma in einem Namen wird nicht unterstuetzt (Komma = oder).
- **2026-10-05 13:38 · BeMuCa** — Critique Runde 1 umgesetzt: (1) 'tag: ui' - ein Schluessel ohne Wert nimmt das naechste Wort (filterTerms klebt an ':'). 'tag:needs review' braucht jetzt Quotes - 'verbreitert nur' gilt nur fuer reine Phrasen; Kommentar, Testname und NOTES korrigiert. (2) leere Komma-Teile zaehlen nicht (filterValues), ein Schluessel ganz ohne Wert laesst das Board wie es ist. (3) ein splitOutsideQuotes fuer Leerzeichen UND Komma - Quotes schuetzen jetzt auch Kommas ('Doe, John'); Picker quotet bei Leerraum oder Komma. (4) links.go-Hinweis druckt den Filter wie getippt. (5) M und x im Milestone-Fenster besitzen nur noch milestone: (splitTerm, gemeinsam mit dem User-Picker) - sonst loeschten sich die beiden Picker gegenseitig. Nicht gemacht: Scrollen im Picker (echte Boards haben 2-3 Namen; Milestone-Fenster hat dieselbe Grenze).
- **2026-10-05 13:54 · BeMuCa** — Critique Runde 2: Kleben an ':' verschluckte das Leerzeichen - 'fix: crash' (unbekannter Schluessel -> Volltext) suchte 'fix:crash'. Jetzt mit Leerzeichen geklebt; bekannte Schluessel trimmen ihren Wert ohnehin. Testfall in TestFilterSpaceAfterAColonStaysOneCondition.
- **2026-10-07 09:18 · BeMuCa** — Berk am 07.10. beim Ausprobieren: 'u' funktioniert, aber unten in der Tastenleiste steht kein u - man weiss nicht, dass man es druecken kann. Die Entscheidung aus pre-process (u nur in der Hilfe, wie M) ist damit von ihm zurueckgenommen, fuer u. M ist ebenfalls nicht in der Leiste; danach hat er nicht gefragt, also bleibt es so. Zurueck aus human nach in-progress auf seine Ansage.
- **2026-10-07 09:27 · BeMuCa** — u-Hinweis: steht direkt hinter 't tags', weil beide ein Fenster ueber dem Board oeffnen. Kein neuer NOTES.md-Eintrag: 'u' ist selbst noch unreleased, die bestehende Unreleased-Zeile 'Press u on the board…' beschreibt die Taste schon; ein Hinweis in der Leiste aendert daran nichts fuer jemanden, der von 0.3.5 kommt. M bleibt ohne Hinweis in der Leiste - nicht gefragt.
