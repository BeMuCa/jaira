---
id: 01M0YT5A6P63902V39HNTQXBY5
title: Eine Lane soll zwei Rueckwege nennen koennen
status: done
ready: true
creator: BeMuCa
goal: rejects-to kann mehr als eine Ziel-Lane benennen
context: "Berk am 25.08.: die Kritik-Lane schickt heute nur zurueck nach in-progress. Sie soll auch an die HITL-Lane abgeben koennen, wenn die Entscheidung dem Menschen gehoert. rejects-to nimmt heute genau einen Wert (core/lane/lane.go:98, validiert: muss installiert sein, nicht man selbst) - 'rejects-to: human' geht also schon, nur nicht beides. Heute traegt der Prompt das: lanes/critique.md sagt, Entscheidungen gehen an die HITL-Lane. Zwei Kanten waeren eine Schema-Aenderung."
definition-of-done: "Eine Lane kann zwei Rueckwege deklarieren, beide werden validiert und im generierten Block als Loop genannt; bestehende Lane-Dateien mit einem Wert funktionieren unveraendert"
blocked-by: []
commits: []
created-at: 2026-08-26T10:34:25Z
updated-at: 2026-09-03T10:28:35Z
claimed-by: EE-3NX6GL3-2836152
claimed-at: 2026-08-31T18:08:54Z
updated-by: BeMuCa
assignee: BeMuCa
outcome-what: "Zwei Review-Findings behoben. 1) 'jaira lanes show' verband die Rueckwege mit Komma - genau die Darstellung, die der Helfer des generierten Blocks selbst ablehnt. orList heisst jetzt board.OrList (exportiert), beide Aufrufer benutzen ihn: 'Rejects to:  in-progress or human'. 2) rejects-to: \"\" - die Form, die core/lane/template.go in jede handgeschriebene Lane schreibt - ist jetzt gepinnt (kein Rueckweg, keine Warnung). Dazu ein Regressionstest fuer die CLI-Darstellung."
outcome-why: "Eine Tatsache, eine Formulierung: 'in-progress, human' liest sich wie eine Reihenfolge (erst zurueck, dann zum Menschen), und zwei Schreibweisen derselben Sache zwingen den Leser, sie erst gleichzusetzen. Exportiert statt kopiert, weil internal/cli core/board schon importiert und eine zweite Kopie des Satzes genau der Grund fuer das Auseinanderlaufen war. Und die leere Form ist der haeufigste Zustand des Feldes ueberhaupt, weil das Template sie ueberall hinschreibt - ungetestet war sie die stillste Bruchstelle."
outcome-resolves: "Findings 1 und 2 aus dem Review. Verifiziert: go test ./core/lane ./core/board ./internal/cli -count=1 gruen; go test ./... -race -count=1 komplett gruen (keine FAIL-Zeile); gofmt -l core internal unveraendert (nur internal/cli/tickets.go, vorbestehend); Ende-zu-Ende auf einem Temp-Board: Liste ergibt 'Rejects to:  in-progress or human', Skalar weiter 'Rejects to:  in-progress'. Achtung Commit-Zuordnung: der Diff liegt in 20ffa1f (fremde Commit-Message, paralleler Agent hat den geteilten Index committet), a5923da ist der leere Commit mit der TQXBY5-Handle und der Erklaerung - siehe Notiz."
review-summary: none
review-gaps: "Nichts entfernt. Gelassen und geprueft: rejects_to im --json ist jetzt immer ein Array - die Flaeche existiert am Release-Tag 62989f1 nicht (git grep leer), also unreleased Surface, gedeckt durch die sync->logbook-Praezedenz; ein Slice-Feld statt Scalar+Extra, damit kein Leser das zweite Ziel verlieren kann; Built-ins und lanes/*.md byte-identisch (Drift-Warn-Falle vermieden); kein Gate-Code angefasst (Invariante 14: rejects-to deklariert, erzwingt nichts)."
review-check: "1. In einem Wegwerf-Board eine Lane-Datei mit rejects-to: [in-progress, human] anlegen: jaira lanes zeigt 'Rejects to: in-progress, human', --json ein Array, der verwaltete Block sagt 'sends work back to in-progress or human'. 2. lanes/critique.md (Scalar-Form, unveraendert) parst ohne neue Warnung, Block-Text wie bisher. 3. rejects-to: [critique] auf critique selbst -> Warnung 'names itself'; ein nicht installiertes Ziel -> Warnung je Ziel. 4. go test ./core/lane ./core/board -count=1 gruen."
review-verdict: "accept (Opus-End-Review auf 62decc0: alle RejectsTo-Leser konvertiert, Back-Compat byte-identisch, JSON-Widening unreleased; Delta vom Koordinator gelesen: board.OrList als eine Stimme fuer beide Ausgaben - Export statt Kopie, weil internal/cli core/board schon importiert -, Leerform rejects-to: \"\" testgepinnt, lanes-show-Rendering testgepinnt)"
---

# Eine Lane soll zwei Rueckwege nennen koennen

## Definition of Done

- [x] Eine Lane kann zwei Rueckwege deklarieren, beide werden validiert und im generierten Block als Loop genannt; bestehende Lane-Dateien mit einem Wert funktionieren unveraendert
  proof: rejects-to nimmt Skalar und Liste (core/lane/lane.go:328), jedes Ziel wird einzeln validiert (:931), der generierte Block nennt beide ('critique sends work back to in-progress or human', core/board/announce.go:120/163); lanes/critique.md mit einem Wert unveraendert gruen (temp-Board: 'Loop: critique sends work back to in-progress')

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] RejectsTo wird []string; Parser akzeptiert Skalar UND Liste (d.Scalar/d.List wie input-requires)
  proof: core/lane/lane.go:328 rejectsTo(): d.List, sonst d.Scalar; TestRejectsToIsParsed + TestRejectsToAcceptsAList
- [x] checkContracts warnt je Ziel und nennt es (selbst / nicht installiert)
  proof: core/lane/lane.go:931 Schleife ueber alle Ziele; TestRejectsToWarnsPerTarget prueft beide Warnungen
- [x] announce.go: LaneFact.RejectsTo []string, Loop-Zeile nennt beide Ziele lesbar
  proof: core/board/announce.go:120 orList + :163 Loop-Zeile; TestLaneSectionNamesEveryBackEdge
- [x] Konsumenten mitziehen: lanes show Text+JSON, market JSON, laneFacts CLI+TUI, lanesEquivalent, template-Kommentar
  proof: internal/cli/tickets.go:912 (Join), market.go/browse.go/laneFacts ueber []string; template.go + lanes/README.md dokumentieren die Liste
- [x] Built-ins und lanes/*.md bleiben byte-identisch (Drift-Falle); Tests fuer Liste, Validierung, Block
  proof: git status: kein core/lane/builtin und kein lanes/*.md geaendert; shipped_test.go laedt lanes/ warnungsfrei
- [x] Verifizieren: go test ./core/lane ./core/board, go test ./... -race, gofmt
  proof: go test ./core/lane ./core/board -count=1 ok; go test ./... -race -count=1 komplett ok; gofmt -l nur internal/cli/tickets.go (vorbestehend)

## Progress
- **2026-08-31 18:17 · BeMuCa** — RejectsTo ist jetzt []string statt eines zweiten Feldes: mit 'RejectsTo string + RejectsToAlso []string' haette jeder Leser, der nur das alte Feld kennt, die zweite Kante still verschluckt. Preis dafuer: 'rejects_to' in 'lanes show --json' und 'market list --json' ist jetzt immer ein Array (Bruch fuer alles, was einen String erwartet hat) - im Repo und in den Skills gibt es keinen solchen Leser, gefunden per grep nur Go-Code.
Deshalb mussten vier bestehende Test-Assertions den Typ mitmachen (rejectsto_test.go 3x, core/board/lanesection_test.go 1x). Verhalten unveraendert, nur slices.Equal statt ==.
Der Parser braucht kein Format-Sniffing: Doc.Scalar weist Sequenzen ab, Doc.List weist Skalare ab - also beides fragen und nehmen, was antwortet (core/lane/lane.go:328).
Falle in core/lane-Tests: catalogueWith() ruft Load() mit leerem Root, das liefert 'Built-ins + Katalog'. Jedes so gebaute Set warnt schon von sich aus ueber die review-Lane ('requires outcome-what'). Ein Test darf dort nie auf 'keine Warnungen' pruefen, nur auf 'keine rejects-to-Warnung'.
Built-ins und lanes/*.md bewusst byte-identisch gelassen: eine Textaenderung in einer Built-in-Datei laesst jede bestehende Projektkopie als Drift/Override warnen. Die DoD verlangt keine Demo der Listenform in einer Lane-Datei, nur dass sie geht - dokumentiert ist sie in core/lane/template.go und lanes/README.md.
Kein Gate angefasst (Invariante 14): rejects-to deklariert weiter nur.
core/lane/share.go Bytes() kopiert Lane-Dateien verbatim, darum ueberlebt die Listenform Export/Publish/Materialise ohne Zutun - hier war nichts zu tun.
Nebenbefund: 'gofmt -l core internal' meldet internal/cli/tickets.go bereits vor dieser Aenderung (zwei map-Literal-Alignments, ~Zeile 72 und ~1121). Nicht angefasst.
Der Working Tree hatte parallel fremde Aenderungen (core/validate/*, internal/cli/flow.go, internal/cli/root.go) - nur meine sieben Dateien wurden committed.
- **2026-08-31 18:44 · BeMuCa** — End-Review accept mit 2 Funden: tickets.go:912 joint mit ', ' was orLists eigener Kommentar als 'liest sich wie eine Sequenz' ablehnt - orList exportieren oder die 8 Zeilen duplizieren; rejects-to: "" (die Form aus template.go:25) ist ungetestet - Probe-Lane in catalogueWith + Assertion len==0 ohne Warnung.
- **2026-08-31 18:53 · BeMuCa** — Review-Findings abgearbeitet.
1) orList ist jetzt board.OrList (exportiert) statt einer Kopie in internal/cli. Begruendung: das Repo hat keinen gemeinsamen Text-Helfer-Ort (dash() liegt lokal in internal/cli/tickets.go:1043, in core/* gibt es keinen exportierten Prosa-Helfer) - also setzt jede Wahl hier einen Praezedenzfall. Export gewinnt, weil genau die Doppel-Formulierung das Finding war: 'lanes show' und der generierte Block beschreiben dieselbe Tatsache und duerfen nicht auseinanderlaufen. internal/cli importiert core/board schon (board.LaneFact, board.Prepare), also kostet der Export keine neue Abhaengigkeit. Eine Kopie waere genau das, was die optimize-Lane als 'Code, der schon woanders existiert' anstreicht.
2) rejects-to: "" ist jetzt gepinnt (TestRejectsToEmptyStringDeclaresNothing): das ist die haeufigste Form des Feldes ueberhaupt, weil core/lane/template.go sie in jede handgeschriebene Lane schreibt. Sie muss 'kein Rueckweg' bedeuten und darf nichts warnen - der Pfad dahin ist unauffaellig: d.List liefert (nil, nil), dann d.Scalar liefert "", daraus wird items = [""] und erst das TrimSpace-Filter macht daraus nil.
3) Zusatz-Test in internal/cli/lanes_test.go (TestLanesShowSpeaksBothBackEdges): die Darstellung 'in-progress or human' war sonst nur im Temp-Board von Hand geprueft, und genau dort war der Fehler.
Ehrlich zum Testlauf: ein 'go test ./... -race -count=1' schlug in internal/cli fehl ('FAIL github.com/BeMuCa/jaira/internal/cli 37.323s'), ohne dass der Testname im Tail sichtbar war. Danach isoliert ('go test ./internal/cli -race' ok, 27.234s) und der komplette Lauf nochmal gruen, ohne eine einzige FAIL-Zeile. Vermutliche Ursache: parallel arbeitende Agenten schreiben im selben Worktree an internal/cli/flow.go und root.go, waehrend der Lauf kompiliert - kein Befund in diesem Diff, aber wer denselben Fehler wieder sieht, sollte dort zuerst schauen.
- **2026-08-31 18:54 · BeMuCa** — Commit-Unfall, wichtig fuer den Reviewer: meine vier Dateien lagen im Index (git add), als ein parallel arbeitender Agent committet hat - sein Commit 20ffa1f ('test(88H1P4): the hook is executed, not just quoted') traegt deshalb meinen Diff mit. Inhalt geprueft und unveraendert (git show HEAD -- core/board/announce.go internal/cli/tickets.go). Nicht zurueckgesetzt: nichts ist gepusht, aber die Historie eines noch arbeitenden Agenten umzuschreiben waere schlimmer als die falsche Commit-Message, und sein Ticket kann die SHA schon notiert haben. Statt dessen a5923da: ein absichtlich leerer Commit mit der TQXBY5-Handle und der Erklaerung, damit die aus git abgeleitete Commit-Liste des Tickets ueberhaupt etwas findet. Lehre fuer den naechsten parallelen Lauf: in einem gemeinsamen Worktree 'git add' und 'git commit' in einem Aufruf, nicht in zwei Schritten.
