---
id: 01M45VMM8DPJB9Z1AN8RN172JS
title: Board laesst sich nach mehreren Usern filtern
status: done
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
commits:
  - 84c69ec675df35ab5e2fc0bfeb155548ea1ee4b1
  - d37b92bcbb230b5bc39446a7d1715fe5830368fc
  - b62d02832b166bd7701301180a5ba95a54884b64
  - 518d4ce94dd20bfff3639a527188f12d26c09fe6
  - 0666af540148979c8f0a43bed28632bdec260012
  - a99798907b217e1f3b2a71fde9d11c74d803ca6b
  - 31c5cea19545c8a50d17281ea2bcd02db6ff21cc
  - 72f2a14b095e970a6ec9a6deb4788fcd3d1385af
created-at: 2026-10-05T11:02:02Z
updated-at: 2026-10-07T20:02:04Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-3594284
claimed-at: 2026-10-07T09:18:49Z
executed-by: claude-opus-5-5
outcome-what: "second-model review of the whole branch: summary, gaps, verdict and a 13-step check written"
outcome-why: "diff meets DoD 1-3, 5, 6 as written and 4 via the header; no defects; suite green under -race; live tmux check passed"
outcome-resolves: "review lane; a person accepts in signoff"
review-summary: "The board filter (/) now takes several conditions at once: words split at spaces and all must hold, a comma inside one value lists alternatives (tag:ui,cli), and double quotes keep a phrase or a name with a space or comma whole (user:\"Alexander Sacharov\"). A new key user:<name> matches a ticket whose assignee OR creator equals the name (exact, case-insensitive). A new key u on the board opens a \"Users\" window listing every assignee and creator of the loaded tickets (logbook cards included) with a count, most tickets first; space ticks, enter writes user:a,b into the ordinary filter and keeps every other condition, x removes only the user: part, esc closes without changing anything; reopening shows the names already in the filter ticked. The milestone picker (M, and x inside it) now owns only the milestone: part of the filter instead of overwriting the whole filter. The header prints the filter as typed instead of %q (no escaped quotes). The bottom key bar names \"u users\" next to \"t tags\"; help (?), README and docs/COMMANDS.md describe u and the new / syntax; three lines under ## Unreleased in core/release/NOTES.md. Tests: 9 new filter tests (internal/tui/filter_test.go), 8 picker tests (internal/tui/userpicker_test.go), key-bar pin in hints_test.go."
review-gaps: "No defect found. Three things a person should know before signing off: (1) DoD item 4 says the status bar shows the active user filter; what shipped shows it in the header line at the top (filter: user:...), where the filter has always been, and the bottom key bar does not show it. The implementer says so in the proof; it is a wording deviation, not a bug, but it is yours to accept. (2) The Users window does not scroll: a board with more names than screen rows overflows (same limit as the milestone picker; this board has 2-3 names). (3) Behaviour change for filters typed before this change: \"export csv\" now matches the two words anywhere, and \"tag:needs review\" now means tag needs AND word review; quotes restore the old reading. NOTES.md says so. Not part of the ticket: untracked .jaira/milestones/demo-*.md files sit in the worktree; they are not in any commit of this branch."
test-verdict: "pass: go build/vet/go test ./... -race -count=1 RC=0; DoD 6 (u users in statusBar view.go:982, TestNarrowBoardShowsAllKeysAndFits) verified; live 160x30 race binary shows 'u users' after 't tags' and u opens the picker"
question: "User-Filter so abnehmen? Taste u, mehrere ankreuzen, schreibt user:a,b in den /-Filter; /-Filter: Leerzeichen = und, Komma = oder, Anfuehrungszeichen halten zusammen. Verhaltensaenderung: ein Feldwert mit Leerzeichen braucht jetzt Quotes (tag:\"needs review\"); M ersetzt nur noch milestone:. Zum Ausprobieren: Branch feat/N172JS-user-filter bauen (dein ~/.local/bin/jaira ist noch PAP369)."
review-verdict: "The diff satisfies DoD 1, 2, 3, 5 and 6 as written, and DoD 4 with the filter shown in the header instead of the status bar; no defects found. Verified independently: read the whole net diff (11 files, 55af6d1..a997989), go test ./... -race -count=1 green (28 packages, RC=0), gofmt and go vet clean, and a scratch build driven in tmux 160x30: key bar shows \"u users\", u opens the Users window with names and counts, space ticks, enter writes user:\"Alexander Sacharov\" into the header and shrinks the lanes (Backlog 136 -> 70), esc clears it. Commit list whole (6 commits = git log origin/master..HEAD). Ready for a person to accept."
review-check: "1. Build: cd /home/berk/git/jAIra-N172JS && go build -race -o ~/.local/bin/jaira ./cmd/jaira - ends with no output.  2. Start: cd /home/berk/git/jAIra-N172JS && jaira - the start screen lists Projects; move to jAIra-N172JS with j/k and press enter - the board opens.  3. Read the bottom line of the board: it contains \"t tags · u users · n new\".  4. Press u - a window titled Users opens over the board, one name per line with a count after it (e.g. BeMuCa first), and its last line reads \"space tick · enter show only the ticked · x show everyone again · esc close\".  5. Press space - the circle before the first name turns from empty to filled; press j, then space - a second name is ticked.  6. Press enter - the window closes; the top-left of the board now reads filter: user:<name1>,<name2> (a name with a space stands in double quotes); the counts in the lane headers are smaller than before, and every card left names one of the two as @assignee or was written by them (enter opens a card and shows assignee and creator; esc closes it).  7. Press / and type a space and the 6-character id of one visible card (e.g. \" 7MG5GB\"), then enter - only that one card is left and the header shows the user: part followed by the id.  8. Press esc - the filter line disappears and the lane counts are the full board again.  9. Press u - no name is ticked; press space on the first name and enter, then u again - that name shows a filled circle; press x - the window closes and the filter line is gone.  10. Press / and type user: and nothing else - the board stays whole; add a name nobody has (user:nobody) - every lane shows 0; press esc.  11. Press / and type two words that appear in one card title, in the wrong order, enter - the card is still shown; put the two words in double quotes, in the wrong order - the card disappears; esc.  12. Press ? - the help lists u and the / line mentions spaces, comma and quotes; esc.  13. Tests: cd /home/berk/git/jAIra-N172JS && go test ./internal/tui -race -count=1 -run \"TestFilter|TestUserPicker|TestUOpens|TestEscOnTheBoard|TestHeaderShows|TestSplitTerm|TestMilestonePickerKeeps|TestNarrowBoard\" - prints a line starting with ok."
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
- **2026-10-07 09:29 · BeMuCa** — critique Runde 4 (nur a997989, Rueckkehr aus human): keine Befunde. 'u users' steht in derselben keys-Liste in statusBar (internal/tui/view.go:982) wie 't tags'; kein einfacherer Ort, nichts Spekulatives; TestNarrowBoardShowsAllKeysAndFits pinnt es. Stehe auf der Notiz 2026-10-07 09:27: kein neuer NOTES.md-Eintrag (die Unreleased-Zeile 'Press u on the board' deckt die Taste ab) und M bleibt ohne Hinweis - beides nicht neu aufgemacht. Veralteter Kommentar 'Hints are dropped from the right' ueber zHint (view.go ~966) widerspricht 'Wrapped, never dropped' darunter - stammt nicht aus diesem Ticket, nicht angefasst.
- **2026-10-07 09:36 · BeMuCa** — testing 2026-10-07: round on DoD 6 / a997989. Full suite green -race -count=1; live tmux check: key bar shows 'u users', u opens picker window. Items 1-5 not re-verified individually; covered by suite.
- **2026-10-07 09:45 · BeMuCa** — review 2026-10-07 (second model, whole branch 84c69ec..a997989): no defect. Read the net diff, suite green under -race (28 pkgs, RC=0), gofmt/vet clean, scratch build driven in tmux: key bar, picker, tick, apply, esc all as the DoD says. esc after enter looked unchanged in the first tmux run - that was tmux escape-time latency (clears with a longer wait), not the code. Three things left for the person in review-gaps: DoD 4 wording (header, not status bar), picker does not scroll, old unquoted multi-word filters read differently now. The user picker sees logbook cards too (openUsers walks m.logged) - consistent, since rebuild filters logged cards with the same matches().
- **2026-10-07 10:14 · BeMuCa** — Abnahme 07.10.: Berk hat im Chat alles angenommen (woertlich: 'ich nehme alles an!'). Die Abnahmeseite https://claude.ai/artifact/UBzQJqGtTpF8B8oDnHuGPM war dabei nicht markiert. Er hat entschieden, Claude hat ausgefuehrt: jaira move --to done --force.
- **2026-10-07 19:55 · BeMuCa** — 07.10.: PR #41 ist gemerged (14:40). Der Ticket-Ref stand noch auf critique, weil die Moves danach nur in diesem Worktree landeten; diese Notiz zieht den Ref auf den Stand done nach.
