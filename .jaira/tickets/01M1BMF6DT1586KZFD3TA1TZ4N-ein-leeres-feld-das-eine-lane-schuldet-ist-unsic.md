---
id: 01M1BMF6DT1586KZFD3TA1TZ4N
title: "Ein leeres Feld, das eine Lane schuldet, ist unsichtbar"
status: done
ready: true
creator: BeMuCa
goal: "Ein Ticket in review/signoff zeigt problem, what, why, resolves, summary, gaps, check - auch wenn ein Feld noch leer ist, als leere Zeile mit der Lane, die es schuldet"
context: |-
  Berk am 31.08.: Tickets in seiner review-Lane zeigen nur dod und problem, er weiss nicht einmal, wie er das Review pruefen soll. Frueher mochte er: problem, what, why, resolves, summary, gaps, check - knappe Stichpunkte mit genug Inhalt.

  Ursache, gelesen in internal/tui/view.go ~848-860: die Detailansicht rendert Outcome (what/why/resolves) und Review (summary/gaps/verdict/check) als Bloecke, aber NUR wenn mindestens ein Feld gefuellt ist; prose() unterdrueckt leere Werte generell. Ein Ticket, das review erreicht, ohne dass die Lanes davor ihre Felder geschrieben haben (sein req-Board: Agents fuellen dort nichts), zeigt also nur Basiszeilen plus DoD.

  Regel: ein Feld, das eine installierte Lane dieses Boards in produces deklariert, wird angezeigt, auch leer - als leere Zeile, die die schuldende Lane nennt (z.B. 'check    — (aus optimize, noch leer)'). Basisfelder (id, lane, assignee, creator, when, goal, context) immer. Reihenfolge fuer review/signoff: problem, what, why, resolves, summary, gaps, check.

  Gehoert zum Schema-Umbau (Schnitt 4, TUI-Falten), kann aber davor als eigener Fix landen; .planning/schema-brainstorm.md Punkte 8-10.
definition-of-done: "Ein Ticket in review ohne gefuellte Outcome/Review-Felder zeigt die sieben Zeilen in der genannten Reihenfolge, leere mit der schuldenden Lane; die Basisfelder stehen in jeder Lane darueber; ein Golden-/View-Test deckt den leeren und den gefuellten Fall ab"
blocked-by: []
commits: []
created-at: 2026-08-31T10:05:05Z
updated-at: 2026-09-03T10:26:38Z
claimed-by: EE-3NX6GL3-2715187
claimed-at: 2026-08-31T17:14:04Z
updated-by: BeMuCa
assignee: BeMuCa
outcome-what: "Neue Gruppe 'Lane fields': detailBody listet nach Outcome/Review jedes Feld, das eine installierte Lane deklariert und fuer das keine Zeile darueber zustaendig ist - leer als '— owed by <lane>', gefuellt mit seinem Wert (aus Doc().Scalar, weil nicht im Struct modelliert). renderSignOff haengt dasselbe nach check an. Reihenfolge: Lanes in Board-Reihenfolge, darin die produces-Reihenfolge; plan und diff bekommen bewusst keine Zeile. labelPad/fieldRow berechnen das Textbudget aus dem echten Label, damit ein langer Key (secrets-findings) den Wert nicht aus dem Pane schiebt."
outcome-why: "Die declared-Gruppen waren hart verdrahtet (goal, Outcome-Trio, Review-Quartett), OwedBy antwortet aber fuer jedes Feld jeder installierten Lane. Ein Board mit eigenen Katalog-Lanes hatte damit genau wieder eine unsichtbare Schuld - die das Gate durchsetzt. Gefuellte Fremdfelder muessen mit, sonst waere die Schuld sichtbar und die Antwort darauf unsichtbar."
outcome-resolves: "Ja. go test ./internal/tui ./core/gate -count=1 gruen, go test ./... -race -count=1 gruen, gofmt -l unveraendert (nur internal/cli/tickets.go). Die drei neuen Tests fallen gegen 847db93 rot ('no heading for the board's own fields'), sind also nicht vakuum. Reihenfolge im Signoff: die sieben Labels bleiben vorn, Fremdfelder danach - begruendet in den Notes."
review-summary: none
review-gaps: "Runde 2: nichts entfernt. Alle vier Review-Funde geschlossen: goal laeuft durch declaredProse (Schuld sichtbar bei Brainstorm-Opt-in), plan bekommt bewusst keine Zeile (Checkliste sagt es schon - in OwedBys Kommentar begruendet), signoff-Ueberdeckung kommentiert, Backlog-Rendering als Entscheidung testgepinnt (TestBacklogDetailCarriesTheSameDebts). Verifiziert gegen sauberen Checkout, damit parallele Edits nichts schoenfaerben. Offen fuer Berk: verdict fehlt in der Spec-Reihenfolge (Entscheidung 10) - Position im Code unveraendert gelassen."
review-check: "Terminal ~100 Spalten, jaira: 1. Ticket im Backlog oeffnen: Outcome/Review-Zeilen stehen als '— owed by <lane>' da (boardweit, per Spec-Regel). 2. jaira dod <id> --option brainstorm auf einem Ticket ohne goal, Detail oeffnen: 'goal — owed by brainstorm'. 3. jaira set <id> outcome-what=x: die Zeile zeigt x. 4. Ticket in signoff: gleiche Schulden als Sektionen; bei leerem goal aber gefuelltem context zeigt problem den context (gewollt, kommentiert)."
review-verdict: "accept (Opus-Reviews auf 8783751+847db93: OwedBy-Kontrakt, Breiten-Arithmetik, Test-Biss alle verifiziert; der Querschnitts-Fund L6/L24 ist mit 2570a2a geschlossen - Koordinator hat den Catch-all gelesen: Board-Reihenfolge statt Map-Order, gefuellte Fremdfelder zeigen ihren Wert, plan/diff begruendet ausgenommen, Label-Budget vom echten Label, Nicht-Vakuitaet gegen den Eltern-Commit belegt)"
---

# Ein leeres Feld, das eine Lane schuldet, ist unsichtbar

## Definition of Done

- [x] Ein Ticket in review ohne gefuellte Outcome/Review-Felder zeigt die sieben Zeilen in der genannten Reihenfolge, leere mit der schuldenden Lane; die Basisfelder stehen in jeder Lane darueber; ein Golden-/View-Test deckt den leeren und den gefuellten Fall ab
  proof: internal/tui/owed_test.go TestDetailShowsOwedFieldsNamingTheLane prueft die sieben Zeilen, die Reihenfolge und dass die Basisfelder darueber stehen; TestDetailShowsTheValueRatherThanTheDebt den gefuellten Fall

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] gate.OwedBy(set, ticket) map[Feld]Lane-ID: laeuft die Lanes in Reihenfolge, nutzt OutputOwed pro Lane (respektiert requires-option), erste Lane gewinnt; nil-sicher
  proof: core/gate/gate.go:608 OwedBy; TestOwedByAgreesWithOutputOwed, TestOwedByIgnoresLanesOffTheRoute, TestOwedByAttributesAFieldToTheFirstLaneThatDeclaresIt, TestOwedByIsNilSafe
- [x] detailBody: owed-Map holen; Outcome- und Review-Block auch rendern, wenn ein Feld leer aber geschuldet ist; leere Zeile als dim '— owed by <lane>' via styleLines+wrap; Blockkopf nur wenn der Block Zeilen hat
  proof: internal/tui/view.go:780 declaredField+owedRow, :866-915 detailBody; TestDetailShowsOwedFieldsNamingTheLane, TestDetailShowsTheValueRatherThanTheDebt
- [x] renderSignOff: dieselbe Regel fuer section(), damit die Signoff-Ansicht nicht ebenfalls die Schuld verschweigt
  proof: internal/tui/signoff.go:46-69; TestSignOffShowsOwedFieldsNamingTheLane
- [x] Tests: internal/tui/owed_test.go (leeres Ticket in review zeigt die sieben Zeilen in Reihenfolge mit schuldender Lane; gefuelltes Feld zeigt Wert statt Schuld-Zeile; Basisfelder stehen darueber; Zeilenbreite <= width; Signoff) und core/gate/owedby_test.go (Reihenfolge, opted-out Lane schuldet nichts, nil-sicher)
  proof: internal/tui/owed_test.go, core/gate/owedby_test.go; gegen den alten Renderer laufen sie rot (view.go/signoff.go zurueckgesetzt, 'sign-off view is missing "summary — owed by review"'), also nicht vakuum
- [x] Verifikation: go test ./internal/tui ./core/gate, go test ./... -race, gofmt -l core internal
  proof: go test ./internal/tui ./core/gate ok; go test ./... -race alle ok (tui 80.8s); gofmt -l core internal nennt nur internal/cli/tickets.go (vorher schon)
- [x] Review-Runde 1: goal-Zeile durch declared() statt prose(); plan als Body-Checkliste bewusst ohne Schuld-Zeile (in OwedBy dokumentiert); Kommentar zum problem-Row-Override im Signoff; vakuumen Off-Route-Test durch den positiven Fall ersetzt; Backlog-Fall festgenagelt
  proof: 847db93; internal/tui/view.go:838-882 declaredProse+declared vor der prose-Gruppe, core/gate/gate.go:617-624 Kommentar, internal/tui/signoff.go:59-62; TestDetailShowsTheGoalDebtWhenBrainstormIsOnTheRoute (faellt gegen 8783751 rot: gar keine goal-Zeile), TestDetailShowsNoDebtRowForABodyChecklistField, TestBacklogDetailCarriesTheSameDebts
- [x] Review-Runde 2: Catch-all 'Lane fields' fuer Felder, die nur das Board deklariert - in detailBody nach der Review-Gruppe, im Signoff nach check; gefuellte Fremdfelder zeigen ihren Wert; plan bleibt ohne Zeile; Reihenfolge = Lane-Reihenfolge, dann produces-Reihenfolge
  proof: 2570a2a; internal/tui/view.go laneFields()+fieldsWithTheirOwnRow+labelPad/fieldRow, signoff.go nach check; TestDetailShowsDebtsForFieldsOnlyTheBoardDeclares, TestDetailShowsTheValueOfAFieldOnlyTheBoardDeclares, TestSignOffAppendsFieldsOnlyTheBoardDeclares - alle drei fallen gegen 847db93 rot

## Progress
- **2026-08-31 17:21 · BeMuCa** — gate.OwedBy sitzt auf OutputOwed, nicht auf einer eigenen Producer-Map: nur so ist requires-option (skipped()) automatisch respektiert und es gibt keine zweite Wahrheit darueber, was eine Lane schuldet.
- **2026-08-31 17:21 · BeMuCa** — Bewusst KEIN Positionsfilter: auch Lanes VOR dem Ticket schulden. Ein Backlog-Ticket zeigt dadurch alle sieben Schuld-Zeilen. Das folgt der Regel aus Punkt 8-10 wortwoertlich ('ein Feld, das eine installierte Lane deklariert, wird angezeigt'), ist aber die eine Stelle, an der Berk vielleicht 'nur bis zur aktuellen Lane' will - Kandidat fuer critique.
- **2026-08-31 17:21 · BeMuCa** — Basisfelder wurden NICHT auf 'immer sichtbar auch wenn leer' umgestellt. Sie stehen schon in jeder Lane oben; nur 'when' fehlt bei einem im Speicher gebauten Ticket ohne created-at (deshalb setzt die Test-Fixture die Zeitstempel). Leere Platzhalter fuer assignee/context waeren Punkt 8 als eigener Schnitt.
- **2026-08-31 17:21 · BeMuCa** — renderSignOff hat keine Basisfelder (kein id/assignee/creator/when) - absichtlich nicht angefasst, das ist Informationsarchitektur von Punkt 8 und nicht dieser Defekt.
- **2026-08-31 17:21 · BeMuCa** — Reihenfolge musste nicht geaendert werden: detailBody rendert goal, context, ..., what, why, resolves, summary, gaps, verdict, check - die sieben stehen schon in der geforderten Folge. verdict bleibt zwischen gaps und check, es faellt nicht unter die sieben, wegwerfen wuerde aber Information kosten.
- **2026-08-31 17:21 · BeMuCa** — Testfalle: lane.Load("") im core/gate-Test laedt den Katalog aus ./lanes (critique/optimize), nicht nur core/lane/builtin - deshalb schuldet review-summary dort 'critique'. Der TUI-Test setzt JAIRA_LANES_DIR auf ein leeres Verzeichnis und bekommt die Builtins (in-progress/review). Board-spezifische Lane-Namen also nicht hart erwarten.
- **2026-08-31 17:21 · BeMuCa** — owedRow styled zeilenweise (styleLines) NACH dem wrap - ein als Block gestylter Mehrzeiler wird von lipgloss auf die breiteste Zeile aufgefuellt und schiesst hinter der Label-Spalte aus dem Pane.
- **2026-08-31 18:02 · BeMuCa** — review (reject, 4 Funde uebernommen, 1 an Berk): (L25) goal laeuft weiter durch prose/row und verschwindet leer - durch declared routen; plan-Schuld entweder als Zeile oder in OwedBy dokumentieren, dass Body-Checklisten keine Zeile haben. (L26) TestDetailOwesNothingToALaneOffTheRoute ist vakuum (der String kann wegen L25 nie erscheinen) - stattdessen Opt-in-Fall: brainstorm angehakt, goal leer -> Schuld-Zeile erscheint. (L28) firstNonEmpty(goal, context) ueberdeckt die goal-Schuld still - Kommentar dazu. (L29) Backlog-Ticket-Rendering per Test pinnen (Spec-Regel 'shown even when empty' gilt boardweit; kompakter Erstblick kommt mit B4MGTP). NICHT anfassen: die verdict-Position (Spec-Frage an Berk, Entscheidung 10 nennt verdict nicht).
- **2026-08-31 18:09 · BeMuCa** — Review-Fund, der der eigentliche Kern war: die goal-Zeile lief noch ueber prose(), das leere Werte verschluckt. Der erste Commit hat also genau den Defekt eine Zeile ueber den reparierten Zeilen stehen gelassen. declaredProse() schreibt die Leerzeile nur, wenn die Zeile auch kommt - sonst bleibt eine Luecke im Rhythmus stehen.
- **2026-08-31 18:09 · BeMuCa** — plan bekommt absichtlich KEINE Schuld-Zeile: es ist eine Body-Checkliste, die leere Checkliste sagt es schon; eine Feldzeile waere eine zweite, schlechtere Version derselben Aussage. OwedBy meldet plan trotzdem (die Frage 'welche Lane schuldet das' hat dieselbe Antwort) und sagt das jetzt im Kommentar - der Renderer entscheidet, er erbt keine Luecke. TestDetailShowsNoDebtRowForABodyChecklistField haelt beides fest.
- **2026-08-31 18:09 · BeMuCa** — Signoff: problem zeigt firstNonEmpty(goal, context), verschluckt damit die goal-Schuld, wenn Kontext da ist. Bleibt so - der Schirm fragt 'was war falsch', und darauf gibt es eine Antwort. Steht jetzt als Kommentar dort.
- **2026-08-31 18:09 · BeMuCa** — Der Off-Route-Test im TUI war vakuum (der String konnte nach Fund 1 nie erscheinen). Off-Route gehoert in core/gate (TestOwedByIgnoresLanesOffTheRoute), im TUI steht jetzt der positive Fall. Merkregel fuer diese Datei: jeden neuen Render-Test einmal gegen den alten Renderer laufen lassen.
- **2026-08-31 18:09 · BeMuCa** — Achtung fuer den naechsten Durchgang: ein anderer Agent arbeitet parallel im selben Worktree (core/validate, internal/tui/home_test.go waren fremd geaendert). Nur die eigenen vier Dateien stagen, nie git add -A.
- **2026-08-31 18:44 · BeMuCa** — End-Review Querschnitts-Fund (L6/L24): die declared-Gruppen in view.go:907/912 sind hartkodiert (goal, Outcome-Trio, Review-Quartett) - Felder adoptierter Lanes (secrets-status, secrets-findings, changelog-entry) sind unsichtbare Schulden, genau die Luecke, die dieses Ticket schliesst. Fix: Catch-all-Pass ueber owed fuer unbekannte Felder unter einer 'Lane fields'-Ueberschrift, mit Test (synthetische Lane produziert ein Fremdfeld).
- **2026-08-31 18:54 · BeMuCa** — Fund der letzten Runde: die declared-Gruppen waren hart verdrahtet, OwedBy antwortet aber fuer JEDES Feld einer installierten Lane. Ein Board mit eigenen Lanes (secrets-scan, changelog-writer) hatte damit eine Schuld, die das Gate durchsetzt und kein Schirm zeigt. Neu: laneFields() sammelt alles, was keine eigene Zeile hat, Gruppe 'Lane fields' am Ende.
- **2026-08-31 18:54 · BeMuCa** — Reihenfolge bewusst NICHT aus der owed-Map: Lanes in Board-Reihenfolge, darin die produces-Reihenfolge der Lane. Sonst wuerde sich die Ansicht zwischen zwei Renders umsortieren (Go randomisiert Map-Iteration). Der Test nagelt secrets-status < secrets-findings < changelog-entry fest.
- **2026-08-31 18:54 · BeMuCa** — Signoff: Fremdfelder haengen NACH check, nicht zwischen den sieben. Die sieben Labels sind die Lesereihenfolge der Sign-off-Entscheidung; ein boardeigenes Feld dort einzusortieren waere geraten. Angehaengt ist es trotzdem sichtbar - und das muss es sein, weil das Gate den Move daran verweigert.
- **2026-08-31 18:54 · BeMuCa** — Gefuellte Fremdfelder zeigen ihren Wert (Doc().Scalar, weil das Feld nicht im Ticket-Struct modelliert ist). Nur die leeren zu zeigen waere die schlechtere Haelfte desselben Bugs: Schuld sichtbar, Antwort darauf unsichtbar. Ein im Speicher gebautes Ticket hat keinen Doc - deshalb legt der Fuell-Test ein echtes Ticket im Store an (unbekannte Frontmatter-Keys ueberleben, geprueft).
- **2026-08-31 18:54 · BeMuCa** — labelPad(): ein Lane-Key kann laenger als die 12 Spalten sein (secrets-findings = 16). Das Textbudget wird jetzt aus dem tatsaechlich gedruckten Label berechnet statt mit 12 angenommen - sonst schiebt ein langer Key den Wert aus dem Pane. TestOwedRowsStayInsideTheWidth laeuft deshalb jetzt mit dem Katalog-Lane-Set.
- **2026-08-31 18:54 · BeMuCa** — fieldsWithTheirOwnRow enthaelt plan und diff aus einem anderen Grund als der Rest: das sind keine Frontmatter-Werte (Body-Checkliste bzw. aus den Commits gebaut), eine Label-Wert-Zeile waere dort immer die schlechtere Version der Checkliste bzw. des Commits-Blocks.
