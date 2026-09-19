---
id: 01M2SNEP2NN02DN5W0ENVSC1GW
title: "Die Pruefschleife gehoert ins Binary, nicht in den Katalog"
status: review
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer jaira ohne Netz benutzt, kommt an critique, optimize und testing heran und erfaehrt beim Blick auf 'jaira lanes', dass es sie gibt - ohne dass ein frisches Board dadurch eine Lane mehr bekommt"
context: |-
  Alex am 2026-09-18, nachdem 'jaira lanes' auf diesem Board dreizehn Lanes zeigte und zehn davon built-in waren.

  Was ist: 'critique', 'optimize' und 'testing' sind KEINE built-ins. Built-in sind genau zehn: backlog, brainstorm, todo, pre-process, in-progress, human, review, signoff, done, blocked (core/lane/lane.go:34, go:embed builtin/*.md). Die drei liegen im Katalog unter lanes/ und erreichen ein Board nur ueber 'jaira lanes market adopt <id>' plus 'jaira lanes add <id>'. lanes/README.md nennt das ausdruecklich Absicht: built-ins stehen auf jedem Board, Katalog-Lanes nur, wenn jemand sie absichtlich holt.

  Warum das heute stoert: ohne diese drei ist ein frisches Board ein Tracker, kein Agenten-Konveyer. Die Schleife critique -> in-progress ist das, was auf diesem Board die echten Defekte findet - an einem einzigen Tag unter anderem einen Unterscheider, der einen schreibenden Worker still zum stummen Leser gemacht haette, und einen fest verdrahteten Branch-Namen 'master' in einem ausgelieferten Prompt. Wer jaira wegen der Agenten installiert und die drei nicht kennt, bekommt nichts davon und erfaehrt auch nicht, dass es sie gibt: 'jaira lanes' zeigt nur, was installiert ist, und der Katalog meldet sich von selbst nirgends.

  Das Gegenargument steht in CLAUDE.md und ist ernst zu nehmen: 'Every feature is measured against is this smaller than paca - the project fails by growing'. Dreizehn Lanes als Voreinstellung zwingt jedem einen Konveyer auf, auch dem, der jaira als reinen Tracker aufsetzt.

  Die Entscheidung ist deshalb NICHT 'alles einbauen', sondern: was ist die Voreinstellung fuer wen. Moegliche Formen, in der Reihenfolge wachsenden Eingriffs - (1) nur sichtbar machen: ein frisches Board oder 'jaira lanes' nennt die Katalog-Lanes, die es nicht hat; (2) die drei ins Binary einbetten, aber nicht auf die Voreinstellung setzen, also ohne Netz adoptierbar; (3) ein zweites Default-Board 'agentic' neben dem schlanken, bei 'jaira init' waehlbar; (4) sie als Voreinstellung fuer jedes neue Board. Welche - das ist die eigentliche Frage dieses Tickets und gehoert in die brainstorm-Lane.
definition-of-done: "Das Ticket nennt EINE gewaehlte Form mit Begruendung, warum die anderen drei verworfen wurden - an der Scope-Regel aus CLAUDE.md gemessen, nicht am Bauchgefuehl."
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-18T07:07:21Z
updated-at: 2026-09-19T19:10:36Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-92780
claimed-at: 2026-09-19T19:09:20Z
mode: conversational
outcome-what: "internal/cli/lanes_test.go: TestLanesAddAfterRemoveAppendsAtEnd heisst jetzt TestLanesAddAfterRemoveLandsAfterItsAnchor; Kommentar und Fehlermeldung nennen den wirklichen Grund (Anker 'done'), Assertion unveraendert"
outcome-why: "Name, Kommentar und Fehlermeldung behaupteten Anhaengen ans Ende - das tut 'jaira lanes add' seit Runde 2 nicht mehr. Ein gruener Test, dessen Name das Gegenteil dessen sichert, was die Assertion schuetzt, fuehrt den naechsten Leser von insertAfterAnchor in die Irre."
outcome-resolves: "vierte und letzte Fundstelle des veralteten 'appending'-Satzes aus Fund 3; Fund der vierten Kritik"
review-summary: |-
  none
  review (Runde 1), am Diff 2fb8f17 gelesen und am gebauten Binary nachgestellt - nicht am Bericht. Der Zweig traegt genau einen Commit, das Ticket-Feld commits: ist leer, der Payload-Diff ist damit der ganze Zweig (23 Dateien, +1308/-73); nichts ist ausserhalb des Payloads passiert.

  WAS GEBAUT WURDE
  1. Ein Bit wird zwei. Lane.Builtin (core/lane/lane.go:135) heisst jetzt nur noch 'reist im Binary mit'; das neue Feld Lane.Default (lane.go:143) heisst 'steht in der Vorauswahl eines frischen Boards' und wird in parse() aus dem Frontmatter 'default-board:' gelesen, mit Vorgabe = builtin (lane.go:298, ueber das vorhandene boolOr, das 'fehlt' von 'steht auf false' unterscheidet). Die zehn vorhandenen builtin/*.md bleiben damit Zeile fuer Zeile unveraendert, und eine Lane aus einer Datei verhaelt sich wie bisher.
  2. Drei Lanes wandern in das Binary. lanes/critique.md, lanes/optimize.md und lanes/testing.md sind nach core/lane/builtin/25-, 26-, 27- VERSCHOBEN (git erkennt Renames, es gibt keine zweite Quelle) und tragen je eine neue Zeile 'default-board: false'. Die Zahlenpraefixe setzen sie zwischen 20-in-progress und 30-human, wo order() mitgelieferte Lanes nach Dateinamen haelt.
  3. Zwei Aufrufer lesen jetzt Default statt Builtin: setUp() (lane.go:639), also die Vorauswahl eines frischen Boards, und internal/tui/defaultboard.go:50, der vorgesetzte Haken im Default-Board-Bildschirm. Dazu eine dritte, die der Plan nicht vorsah: migrateLegacy() (lane.go:694) - ein altes Lane-Verzeichnis ohne order-Datei haette sonst drei Lanes bekommen, die es nie hatte. Die Begruendung steht als Kommentar an der Stelle; die uebrigen neun .Builtin-Fundstellen (share, order, TUI-Beschriftung 'built-in', tickets.go-Quelle) sind unveraendert und dort auch richtig.
  4. 'jaira lanes' bekommt einen Fuss (internal/cli/tickets.go:996-1063): die Lanes, die das Binary traegt und dieses Board nicht installiert hat, gefiltert auf 'l.Builtin && !l.Default' - eine absichtlich entfernte Vorauswahl-Lane wird also nicht wieder angeboten, eine Lane ausserhalb der Vorauswahl schon. Im Text drei Zeilen mit Beschreibung plus die Zeile, die sie holt; im JSON ein Feld 'available' mit id, name, description, agentic, model_tier und einem fertigen 'add'-Kommando. Die Quelle ist lane.Installable() - eingebettete Builtins plus ein Glob ueber ~/.jaira/lanes, kein Netz.
  5. 'jaira lanes add' haengt nicht mehr an, sondern setzt (core/lane/order.go:253-343). insertAfterAnchor() setzt die Lane hinter den Anker, den anchorIndex() findet; anchorIndex() folgt der after:-Kette DURCH nicht installierte Lanes hindurch (mit seen-Map gegen Zyklen), so dass 'jaira lanes add testing' auf einem frischen Board ueber optimize -> critique -> in-progress zwischen in-progress und human landet und nicht hinter signoff. Loest nichts auf, geht die Lane vor die erste terminale Lane, und nur dann gibt es eine Warnung - die den Anker nennt, den die Lane selbst traegt (l.After), nicht den Namen, an dem die Kette endete. Nur die neue id bewegt sich, die uebrige Reihenfolge bleibt.
  6. lane.Add() gibt jetzt vier Werte zurueck (dst, after, warnings, err). 'jaira lanes add' sagt in der Erfolgszeile 'after <id>' bzw. 'at the front' und fuehrt 'after' und 'warnings' im JSON; die TUI faltet ihre zwei Aufrufer in addedMsg() und zeigt die Warnung in der Statuszeile (die Position sieht man dort ohnehin, weil reload() die Spalten neu zeichnet).
  7. Der Katalog wird an den Tag des laufenden Binaries gebunden (core/market/market.go:51-105). apiBase() setzt '?ref=v<release.Current>' - ueber url.Parse/Query().Set und nicht per String-Verkettung, damit eine per JAIRA_MARKET_API gesetzte Adresse mit eigener Query nicht kaputtgeht. Bei 'dev' kein ref, dafuer sagt Unpinned() den Satz, den internal/cli/market.go in List und Adopt als 'note:' ausgibt und im JSON als Feld 'unpinned' fuehrt - zurueckgegeben statt gedruckt, weil apiBase() pro Anfrage laeuft.
  8. Fuenf Zeilen unter '## Unreleased' in core/release/NOTES.md, je eine pro Punkt; der Zeilenscan liest genau fuenf. docs/COMMANDS.md:153 und die Hilfetexte sagen nicht mehr 'appending'.

  WAS DAMIT DAS PROBLEM LOEST: die drei Lanes sind ohne Netz installierbar, weil Installable() sie schon vorher angeboten haette, sobald sie eingebettet sind; sie draengen sich keinem Board auf, weil setUp() jetzt Default liest; und ein frisches Board erfaehrt von ihnen, weil der Fuss sie nennt. Nachgestellt: 'jaira lanes add testing' mit JAIRA_MARKET_API auf 127.0.0.1:1 installiert sie, danach 'add optimize' und 'add critique' - Endordnung in-progress critique optimize testing human, der Fuss schweigt danach.
review-gaps: |-
  internal/tui/lanes.go: addFromCatalogue und addAvailable trugen denselben Sechszeiler samt Kommentar hinter einer sofort ueberschriebenen ls.msg-Zuweisung - in einen Helfer addedMsg(id, warnings) gefaltet, Verhalten unveraendert. Stehen gelassen mit Begruendung: terminalIDIndex neben terminalIndex (zwei Typen, in Durchgang 1 geschlossen; beide Fallbacks am Code auf gleiche Position gegengeprueft), die Unpinned()-Notiz zweimal in internal/cli/market.go (spiegelt die vorhandene Overridden()-Doppelung), die after-Nachsuche in lane.Add (ein dritter Rueckgabewert fuer eine Schleife ueber 13 Eintraege auf einem einmaligen Pfad) und der ''-Zweig in pinnedRef (Vorsicht gegen -ldflags, nicht unerreichbar). Kein toter Code; lane.Installable() auf dem 'jaira lanes'-Pfad ist netzfrei (embedded Builtins plus ein Glob), die Startzeit-Regel bleibt unberuehrt.
  review (Runde 1). Gegengeprueft, nicht geglaubt: go build ./... , go vet ./... und go test ./... -count=1 alle RC=0, keine uebersprungenen Pakete; alle 13 in den DoD-Beweisen genannten Testfunktionen existieren an den genannten Dateien; jede Datei-Fundstelle stimmt am Arbeitsbaum; die Entscheidung vom 2026-09-19 ist eingehalten - in core/market/market.go kommt weder 'Builtin' noch 'Default' vor, market filtert nichts zusaetzlich. Die DoD ist damit erfuellt. Drei Befunde, alle ausserhalb des Go-Codes oder unterhalb der Rueckweis-Schwelle:

  BEFUND 1 (der einzige ernste): lanes/README.md ist von diesem Ticket nicht angefasst worden und luegt jetzt an vier Stellen. Die Datei ist die Titelseite des Katalogs auf GitHub und das, was ein Adoptierender liest.
  - Zeile 11-13, der Schnellstart, fuehrt woertlich 'jaira lanes market adopt critique' vor. critique liegt seit diesem Commit nicht mehr in lanes/, der Aufruf schlaegt an einem Binary dieses Tags fehl.
  - Zeile 15-16: 'Aus einem Clone macht jaira lanes adopt lanes/critique.md dasselbe'. Die Datei ist weg.
  - Zeile 40-44, die Tabelle 'What is here', fuehrt critique und optimize weiter als Katalog-Lanes, und die drei Absaetze darunter erklaeren die critique/optimize-Schleife als Katalog-Sache. Im Katalog liegen nur noch changelog-writer und secrets-scan.
  - Zeile 58-63: 'jaira lanes add appends the lane as the last line of .jaira/lanes/order, so a freshly adopted lane is the rightmost column - after: is only consulted when there is no order file at all, and jaira init writes one'. Beides ist seit Fund 2 dieses Tickets falsch: Add setzt nach der Kette und konsultiert after: gerade dann, wenn eine order-Datei da ist.
  Das ist zugleich die FUENFTE Fundstelle des 'appending'-Satzes aus Fund 3, und die Notiz von Runde 4 behauptet ausdruecklich, die Suche sei erschoepfend gewesen. Sie war es nicht: gesucht wurde nach 'append(s|ing) (it) (at|to) the end' und 'end of the (column) order', und die README formuliert 'appends the lane as the last line'. Der Beweis der Planzeile 'Grep-Abdeckung' traegt also nicht.

  BEFUND 2 (Verhaltensaenderung ohne NOTES-Zeile und ohne Test): wer critique/optimize/testing frueher mit 'jaira lanes market adopt' in ~/.jaira/lanes geholt und dort angepasst hat, bekommt seine eigene Fassung nicht mehr. Installable() (core/lane/order.go:190-205) legt zuerst die Builtins und erst danach den Glob ueber UserLanesDir() in dieselbe seen-Map - der eingebettete Namensvetter gewinnt, still. Am gebauten Binary nachgestellt: mit einer ~/.jaira/lanes/critique.md, deren description auf 'MY OWN ADOPTED CRITIQUE' geaendert war, nennt der neue Fuss von 'jaira lanes' die eingebettete Beschreibung, und 'jaira lanes add critique' schreibt die eingebettete Datei aufs Board. Fuer die zehn alten Builtins galt das immer; neu ist, dass es genau die drei Lanes trifft, die bis heute NUR ueber diesen Weg zu haben waren. Die NOTES-Zeile sagt 'Adopting them from the marketplace is no longer needed' und kein Wort davon, dass eine schon adoptierte, angepasste Fassung ueberdeckt wird.

  BEFUND 3 (kosmetisch, kein Rueckweg): core/market/market.go:78 behandelt release.Current == "" wie 'dev', aber Unpinned() (Zeile 102) baut daraus den Satz 'this build reports version , which is no released tag'. release.Current ist mit "dev" vorbelegt (core/release/release.go:23), der Zweig ist also nur per -ldflags '-X ...=' erreichbar - ein Satz mit Loch statt einer Aussage, aber niemand sieht ihn.

  NICHT BEANSTANDET, damit es niemand nochmal aufmacht: die TUI-Zeile addedMsg() nennt die Nachbarlane nicht, aber reload() zeichnet die Spalten neu und setzt den Cursor auf die neue Lane - die Position ist dort sichtbar, anders als in der CLI. Und die dritte .Builtin-Umstellung in migrateLegacy(), die der Plan nicht vorsah, ist richtig und begruendet: ohne sie bekaeme ein Alt-Board ohne order-Datei drei Lanes, die es nie hatte.
test-verdict: "pass: go build/vet sauber, volle Suite mit -race und geleertem Cache gruen (RC=0), DoD 1-6 am Arbeitsbaum geprueft, Verhalten an einer frischen Testdoska und an einem Binary mit -X main.version=0.1.4 nachgestellt"
question: "Die Pruefschleife liegt im Binary, Suite und Verhalten sind geprueft - nimmst du die Arbeit an, oder soll noch etwas hinein? Zwei Punkte zum Mitentscheiden: (1) Die Arbeit ist NICHT committet - conversational-Modus, die Commit-Zeile gehoert dir. (2) 'jaira lanes market' bietet critique und optimize an einem dev-Build weiter an, weil lanes/ auf master noch existiert; das verschwindet erst mit dem Tag, der diesen Zweig enthaelt - soll das so bleiben, oder soll market eingebettete Lanes zusaetzlich herausfiltern?"
---

# Die Pruefschleife gehoert ins Binary, nicht in den Katalog

## Definition of Done

- [x] Das Ticket nennt EINE gewaehlte Form mit Begruendung, warum die anderen drei verworfen wurden - an der Scope-Regel aus CLAUDE.md gemessen, nicht am Bauchgefuehl.
  proof: Notiz 'ENTSCHIEDEN von Alex am 2026-09-18' auf diesem Ticket: Form 2 gewaehlt, Formen 1/3/4 und Weg B mit Begruendung verworfen
- [x] Ein frisch mit 'jaira init' angelegtes Board erfaehrt von der Pruefschleife, ohne dass jemand den Katalog kennt. Nachgestellt an einer leeren Testdoska: der Weg von 'jaira init' zu einer arbeitenden critique-Lane ist ohne Vorwissen gehbar.
  proof: internal/cli/tickets.go:1010 Filter 'l.Builtin && !l.Default'; Tests internal/cli/lanes_test.go TestLanesNamesTheShippedLanesThisBoardLacks, TestLanesJSONCarriesAvailable, TestLanesSaysNothingWhenTheBoardHasEverything, TestLanesDoesNotOfferBackALaneTheBoardRemoved
- [x] Wer jaira ohne Netz benutzt, kommt an die drei Lanes heran, oder die Fehlermeldung sagt, dass sie aus dem Netz kommen und wie man sie sonst bekommt. 'jaira lanes market' holt heute von GitHub.
  proof: internal/cli/lanes_test.go TestLanesAddInstallsTheReviewLoopWithoutNetwork und TestLanesAddFollowsAnAnchorThatIsItselfUninstalled - 'jaira lanes add testing' allein landet ueber die Ankerkette (core/lane/order.go anchorIndex) zwischen in-progress und human, nicht hinter signoff
- [x] Ein Board, das jaira als reinen Tracker benutzt, wird nicht mit einem Konveyer beladen, den niemand faehrt - die Wahl bleibt eine Wahl.
  proof: core/lane/defaultboard_test.go TestFreshBoardGetsOnlyTheDefaultLanes — ein frisch angelegtes Board laedt als genau die zehn Lanes von heute, obwohl das Binary jetzt dreizehn traegt
- [x] Je eine Zeile in core/release/NOTES.md fuer das, was ein Benutzer dadurch anders tut.
  proof: core/release/NOTES.md unter '## Unreleased': fuenf Zeilen, per Notes() gegengeprueft - der Zeilenscan liest genau fuenf Aenderungen
- [x] Der Katalog ist an eine Version gebunden oder sagt, dass er es nicht ist: 'jaira lanes market' holt heute von 'https://api.github.com/repos/BeMuCa/jaira/contents/lanes' ohne '?ref=' und bekommt damit den HEAD des Default-Branches, egal wie alt das laufende Binary ist. Entweder fragt der Aufruf den Tag des laufenden Binaries ab, oder er sagt dem Benutzer, dass er die Entwicklungsfassung bekommt. Nachgestellt mit einem Binary, das eine aeltere Version meldet.
  proof: core/market/market.go apiBase()+pinnedRef()+Unpinned(); Tests core/market/market_test.go TestListPinsTheCatalogueToTheRunningTag, TestDevBuildSendsNoRefAndSaysSo, TestRefIsSetOnAnAddressThatAlreadyHasAQuery; nachgestellt mit einem Binary aus -ldflags '-X main.version=0.1.4': ASKED /contents/lanes?ref=v0.1.4

## Options

- [x] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Lane.Builtin aufteilen: Feld Default im Struct (core/lane/lane.go:136), aus Frontmatter 'default-board:' in parse() (lane.go:284), Vorgabe = Builtin
- [x] setUp() (core/lane/lane.go:627) und internal/tui/defaultboard.go:50 auf Default umstellen; die uebrigen neun .Builtin-Stellen unveraendert lassen
- [x] Test: ein frisch angelegtes Board hat genau die zehn Lanes von heute (DoD 4)
- [x] lanes/critique.md, lanes/optimize.md, lanes/testing.md nach core/lane/builtin/ verschieben als 25-/26-/27-, je mit 'default-board: false'
- [x] Test: 'jaira lanes add critique' installiert die Lane bei unerreichbarem Netz (JAIRA_MARKET_API auf eine tote Adresse), und sie steht zwischen in-progress und human (DoD 3)
- [x] 'jaira lanes' um einen Fuss erweitern (internal/cli/tickets.go:977): mitgelieferte, nicht installierte Lanes namentlich plus die Zeile, die sie holt; im JSON ein Feld 'available' (DoD 2)
- [x] Test fuer diesen Fuss in Text und JSON, und dass er schweigt, wenn nichts fehlt
- [x] core/market/market.go:45: apiBase() haengt '?ref=v<release.Current>' an; bei 'dev' ohne ref plus eine Zeile auf stderr, dass die Entwicklungsfassung kommt (DoD 6)
- [x] Test: mit gesetzter Version geht der ref an den Server, mit 'dev' nicht und die Ansage erscheint
- [x] Je eine Zeile unter '## Unreleased' in core/release/NOTES.md fuer die drei Lanes im Binary, den neuen Fuss von 'jaira lanes' und den versionsgebundenen Katalog (DoD 5)
- [x] go build ./... , go vet ./... und die volle Testsuite gruen
- [x] Fund 1: Fuss von 'jaira lanes' filtert auf 'Builtin && !Default' statt Builtin - eine entfernte Vorauswahl-Lane wird nicht wieder angeboten
- [x] Fund 2: insertAfterAnchor loest die Ankerkette ueber Installable() auf, damit 'jaira lanes add testing' auf einem frischen Board nicht hinter signoff landet
  proof: core/lane/order.go anchorIndex; Test TestLanesAddFollowsAnAnchorThatIsItselfUninstalled
- [x] Tests fuer beide Funde, dann go build/vet und volle Suite
- [x] Punkt A: 'jaira lanes add' sagt in der Erfolgszeile, hinter welcher Lane sie gelandet ist - lane.Add gibt die Position zurueck, JSON bekommt ein Feld dafuer
  proof: core/lane/order.go Add gibt den Vorgaenger zurueck; internal/cli/lanes.go:128 Erfolgszeile 'after <id>' und JSON-Feld 'after'
- [x] Punkt B: die Warnung nennt den Anker, den die Lane selbst traegt (l.After), nicht den Namen, an dem die Kette endete
  proof: core/lane/order.go insertAfterAnchor: Warnung setzt l.After ein, anchorIndex gibt nur noch den Index zurueck
- [x] Punkt C: der Kommentarsatz in internal/cli/tickets.go, der behauptet, das Angebot gelte nur nie installierten Lanes - nur der Satz, nicht der Filter
  proof: internal/cli/tickets.go: der Kommentar am 'Builtin && !Default'-Filter, dazu derselbe Satz im Kopf von TestLanesDoesNotOfferBackALaneTheBoardRemoved
- [x] Tests fuer A und B, dann go build ./... , go vet ./... und die volle Suite
  proof: TestLanesAddSaysWhichLaneItLandedAfter und TestLanesAddWarnsWithTheAnchorTheLaneNames; go build/vet sauber, volle Suite gruen
- [x] Runde 4, Fund der dritten Kritik: die drei veralteten 'appending'-Saetze (internal/cli/lanes.go Long von 'lanes add', internal/tui/lanes.go Kommentar an addFromCatalogue, docs/COMMANDS.md Kommandotabelle) sagen stattdessen, dass die Lane dorthin gesetzt wird, wohin ihre after:-Kette zeigt - durch nicht installierte Lanes hindurch
  proof: internal/cli/lanes.go:100, internal/tui/lanes.go:438, docs/COMMANDS.md:153; am gebauten Binary nachgestellt: 'jaira lanes add --help' sagt jetzt 'placing it where its after: field points'
- [x] Runde 5, Fund der vierten Kritik: internal/cli/lanes_test.go TestLanesAddAfterRemoveAppendsAtEnd heisst und begruendet nicht mehr 'angehaengt', sondern nennt den wirklichen Grund; Assertion unveraendert
  proof: internal/cli/lanes_test.go:962 TestLanesAddAfterRemoveLandsAfterItsAnchor - gruen, Assertion unveraendert

## Progress
- **2026-09-18 07:10 · Alexander Sacharov** — Alex am 2026-09-18, beim Durchdenken der Katalog-Idee: 'kann man dann nicht alle Lanes in den Markt legen, und werden sie aus dem Release geladen?'

Zweite Frage nachgesehen, und die Antwort ist nein: core/market/market.go:45 baut die Adresse als 'https://api.github.com/repos/BeMuCa/jaira/contents/lanes' - ohne '?ref='. Die GitHub-Contents-API liefert ohne ref den HEAD des Default-Branches. Der Katalog ist also gar nicht versioniert: wer 0.1.4 installiert hat, bekommt heute die Lanes von heutigem master. Nennt eine Lane spaeter ein Feld oder eine Option, die sein Binary nicht kennt, kommt sie kaputt bei ihm an - und umgekehrt laesst sich der Katalog fuer ein altes Release nicht reparieren, es gibt nur einen fuer alle.

Erste Frage, warum nicht ALLES in den Katalog kann: die zehn built-ins sind der Offline-Boden. 'jaira init' funktioniert heute im Flugzeug. Mit allem im Katalog haengt das Anlegen eines Boards an GitHub - nicht erreichbar, Rate-Limit erschoepft, Repository umgezogen, und es entsteht kein Board. Dazu waere der Weg zu einem arbeitenden Board eine Kette aus zehnmal 'market adopt' plus 'lanes add' statt einem 'jaira init'.

Was das fuer dieses Ticket bedeutet: die Grenze ist nicht 'wichtig oder unwichtig', sondern 'ohne das gibt es kein Board' gegen 'das aendert, wie man arbeitet'. Und Form 2 aus dem Kontext - einbetten, aber nicht voreinstellen - ist genau deshalb ein ernster Kandidat: sie loest Offline und Versionierung in einem, weil eine eingebettete Lane mit dem Binary mitreist und damit automatisch zu ihm passt.
- **2026-09-18 07:12 · Alexander Sacharov** — Alex fragt nach, wie die Versionierung technisch aussieht: 'aendern wir dann market UND das Erstellen des Releases, damit alle Dateien als Download im Release liegen?' Drei Wege durchgesehen, damit brainstorm nicht bei null anfaengt.

A - EIN PARAMETER. Die GitHub-Contents-API nimmt '?ref='. Aus core/market/market.go:45 wird 'https://api.github.com/repos/BeMuCa/jaira/contents/lanes?ref=v<version des laufenden Binaries>'. Am Release aendert sich NICHTS - GitHub liefert den Baum jedes Tags von selbst. Ein Source-Build hat keinen Tag: dann auf den Default-Branch zurueckfallen und es laut sagen, nicht stillschweigend.

B - DATEIEN ALS RELEASE-ASSETS. Das war die Frage. Der meiste Aufwand: der Release-Workflow muss sie hochladen, und market muss statt der Verzeichnisliste die Asset-Liste lesen. Heraus kommt genau das, was A schon liefert. Der einzige Mehrwert waere Unabhaengigkeit von der Repository-Struktur - dafuer lohnt kein eigener Auslieferungsweg.

C - EINBETTEN. Fuer critique, optimize und testing der eigentliche Punkt: sie reisen mit dem Binary, die Version stimmt per Konstruktion, es braucht gar kein Netz. market bleibt dann das, was es sein sollte - der Ort fuer FREMDE Lanes, nicht fuer die eigenen.

Empfehlung des Dispatchers, nicht Entscheidung: C fuer die drei, A fuer den Rest des Katalogs, B gar nicht. Die Entscheidung gehoert der brainstorm-Lane und ist DoD 1.
- **2026-09-18 08:09 · Alexander Sacharov** — Dispatcher, beim Start am 2026-09-18: Alex ist weg und hat gesagt, dieses Ticket soll anlaufen, sobald 7XXT00 eine Menschen-Lane erreicht - das ist jetzt passiert. Er hat die Entscheidung dieses Tickets NICHT vorweggenommen, also ist sie offen und dieses Ticket besteht praktisch nur aus ihr. Nicht raten: die vier moeglichen Formen stehen im Kontext, die drei technischen Wege zur Versionierung (A ?ref=, B Release-Assets, C einbetten) in der Notiz darueber, samt der Empfehlung des vorigen Dispatchers (C fuer die drei Lanes, A fuer den Rest, B gar nicht) - eine Empfehlung, keine Entscheidung. Arbeite brainstorm aus, lege die Frage sauber vor und halte dann an. Der Zweig feat/VSC1GW haengt an master (3d12b7a, 0.3.0 gemerged).
- **2026-09-18 08:12 · Alexander Sacharov** — brainstorm, am Code nachgesehen - was WIRKLICH im Weg steht, und warum Form 2 nicht das ist, was sie im Kontext zu sein scheint.

WAS DER CODE ZEIGT

1. 'eingebettet' und 'auf jedem Board voreingestellt' sind heute EIN Bit, nicht zwei. Builtins() (core/lane/lane.go:385) liest go:embed builtin/*.md, und setUp() (core/lane/lane.go:614) schreibt beim ersten Oeffnen eines Boards genau die Lanes als Dateien, fuer die l.Builtin true ist - sofern kein Default-Board existiert. Wer critique.md einfach nach core/lane/builtin/ legt, hat damit NICHT Form 2 gebaut, sondern Form 4: dreizehn Lanes auf jedem frischen Board. Form 2 verlangt, das Feld Lane.Builtin (lane.go:136) in zwei Tatsachen zu trennen - 'reist im Binary mit' und 'steht in der Vorauswahl'. Das ist ein Frontmatter-Feld mehr (etwa 'default-board: false') oder ein zweites embed-Verzeichnis, plus die eine Zeile in setUp, die statt l.Builtin die Vorauswahl fragt.

2. Ist das getan, ist Offline geschenkt. Installable() (core/lane/order.go:175) bietet schon heute built-ins PLUS ~/.jaira/lanes an, und lane.Add() kopiert von dort auf das Board. Eine eingebettete critique-Lane ist damit ohne jedes Netz per 'jaira lanes add critique' installierbar - 'market adopt' entfaellt fuer sie ganz. Das ist DoD 3 erledigt, ohne eine Fehlermeldung schreiben zu muessen.

3. Form 1 allein erreicht Offline NICHT. Die Namen der Katalog-Lanes kennt das Binary heute nirgends; die Liste kommt ausschliesslich aus market.List() ueber die Contents-API (core/market/market.go:45). 'jaira lanes' koennte also ohne Netz gar nicht nennen, was fehlt - es sei denn, man verdrahtet die drei Namen fest, was ein halbes Einbetten ohne dessen Nutzen ist.

4. Form 3 hat die Mechanik schon, braucht aber trotzdem neuen Stoff. DefaultBoard (core/lane/defaultboard.go) ist per Benutzer und kann jede Lane-Auswahl halten - ein zweites 'agentic'-Profil waere kein neues Konzept im Kern. Neu waeren die Frage in 'jaira init' und ein zweiter mitgelieferter Auswahl-Satz; und sie hilft nur dem, der beim Anlegen schon weiss, was er waehlen soll.

5. DoD 6, Versionierung: version steht in cmd/jaira/main.go:11 (ldflags, 'dev' im Source-Build) und ist dem Paket market heute gar nicht bekannt. Weg A - '?ref=v<version>' an apiBase() - sind wenige Zeilen plus Durchreichen der Version, Fallback auf den Default-Branch mit lauter Ansage, wenn die Version 'dev' ist. Weg B (Release-Assets) liefert dasselbe und kostet den Release-Workflow dazu. Unter Form 2 braucht der Markt fuer die drei Lanes gar nicht mehr versioniert zu sein - er bleibt noetig fuer den Rest des Katalogs (changelog-writer, secrets-scan).

EMPFEHLUNG: Form 2, und Form 1 obendrauf, weil sie durch 2 erst moeglich und dann fast gratis ist. Dazu Weg A.

Form 2 einbetten, aber nicht voreinstellen: kein Board bekommt eine Lane mehr als heute (DoD 4 bleibt gewahrt), kein neues Kommando, keine Frage in 'jaira init' - es wird ein Bool in zwei Bools zerlegt und eine Netzabhaengigkeit entfernt statt eine hinzugefuegt. An der Scope-Regel gemessen ist das die einzige der vier Formen, die das Binary nicht um ein Konzept wachsen laesst.
Form 1 obendrauf: 'jaira lanes' nennt am Fuss die mitgelieferten, nicht installierten Lanes mit der Zeile, die sie holt. Ohne das ist DoD 2 nicht erfuellt - eingebettet und ungenannt ist so unsichtbar wie im Katalog. Mit 2 im Ruecken ist es eine Liste aus dem Binary statt aus dem Netz.
Form 4 faellt an der Scope-Regel: dreizehn Lanes als Voreinstellung zwingt jedem Tracker-Benutzer einen Konveyer auf, den niemand faehrt.
Form 3 faellt am Aufwand gegen den Nutzen: init-Dialog plus zweites Profil, und der Nutzen trifft nur den, der die Wahl schon versteht. Sie bleibt spaeter nachruestbar - ein 'agentic'-Default-Board ist nach Form 2 nur noch eine Datei mit dreizehn Namen darin.
Weg B faellt: eigener Auslieferungsweg fuer das Ergebnis, das A mit einem Query-Parameter liefert.

OFFEN UND MENSCHENSACHE: die Wahl selbst (DoD 1). Alex hat sie ausdruecklich nicht vorweggenommen.
- **2026-09-18 10:30 · Alexander Sacharov** — ENTSCHIEDEN von Alex am 2026-09-18, in der Sitzung mitlesend (conversational): FORM 2, mit Form 1 obendrauf und Weg A fuer den Rest des Katalogs.

Gewaehlt: die drei Lanes reisen im Binary mit, stehen aber NICHT in der Vorauswahl. Dazu nennt 'jaira lanes' die mitgelieferten, nicht installierten Lanes, damit ein frisches Board von ihnen erfaehrt. 'jaira lanes market' bekommt '?ref=v<version des laufenden Binaries>' fuer den Katalog, der weiterhin aus dem Netz kommt.

Verworfen, an der Scope-Regel aus CLAUDE.md gemessen:
- Form 1 allein: erreicht DoD 3 nicht. Das Binary kennt die Namen der Katalog-Lanes nirgends, sie kommen nur aus market.List() (core/market/market.go:45). Ohne Netz kann es nicht nennen, was fehlt - ausser man verdrahtet die drei Namen fest, was ein halbes Einbetten ohne dessen Nutzen ist.
- Form 3: die Mechanik traegt zwar (DefaultBoard, core/lane/defaultboard.go), neu waeren aber ein Dialog in 'jaira init' und ein zweiter mitgelieferter Auswahl-Satz - ein Konzept mehr im Binary, und der Nutzen trifft nur den, der die Wahl beim Anlegen schon versteht. Nach Form 2 bleibt sie jederzeit nachruestbar: ein 'agentic'-Default-Board ist dann nur noch eine Datei mit dreizehn Namen darin.
- Form 4: dreizehn Lanes als Voreinstellung zwingt jedem, der jaira als reinen Tracker aufsetzt, einen Konveyer auf, den niemand faehrt. Das ist genau das Wachstum, gegen das die Scope-Regel steht.
- Weg B (Release-Assets): ein eigener Auslieferungsweg fuer das Ergebnis, das Weg A mit einem Query-Parameter liefert.

Warum Form 2 als einzige die Scope-Regel besteht: kein Board bekommt eine Lane mehr als heute, kein neues Kommando, keine Frage in 'jaira init'. Es wird ein bestehendes Bool (Lane.Builtin, core/lane/lane.go:136) in zwei Tatsachen zerlegt - 'reist im Binary mit' und 'steht in der Vorauswahl' - und eine Netzabhaengigkeit entfernt statt eine hinzugefuegt.
- **2026-09-18 10:33 · Alexander Sacharov** — pre-process, am Code nachgesehen: wie Form 2 + Form 1 + Weg A konkret gebaut werden, und was dabei im Weg steht.

DAS EINE BIT, DAS ZWEI WERDEN MUSS
Lane.Builtin (core/lane/lane.go:136) traegt heute zwei Tatsachen: 'reist im Binary mit' und 'steht in der Vorauswahl'. Der Plan trennt sie: Builtin bleibt 'reist im Binary mit', neu kommt Default 'steht in der Vorauswahl', gelesen aus dem Frontmatter-Feld 'default-board:' mit Vorgabe = Builtin. Damit aendert sich an den zehn vorhandenen builtin/*.md-Dateien keine Zeile; die drei neuen setzen 'default-board: false'.

Warum ein Frontmatter-Feld und nicht ein zweites embed-Verzeichnis: ein zweites Verzeichnis braucht eine zweite Leseschleife in Builtins() (lane.go:385) und eine zweite Stelle, die 'ist das eine mitgelieferte Lane' beantwortet. Das Feld ist eine Zeile in parse() (lane.go:284) und eine im Lane-Struct.

DIE ELF FUNDSTELLEN VON .Builtin, JE ENTSCHIEDEN
Auf Default umzustellen sind genau zwei:
- core/lane/lane.go:627 (setUp) - die Vorauswahl eines frischen Boards. Das ist DoD 4.
- internal/tui/defaultboard.go:50 - der Haken, der im Default-Board-Bildschirm vorgesetzt ist.
Builtin bleibt richtig bei: share.go:21 und :146 (mitgelieferte Lanes publiziert man nicht), lane.go:715 (migrateLegacy), lane.go:850/857 (order), internal/tui/defaultboard.go:139, internal/tui/lanes.go:181/359/519 (Beschriftung 'built-in'), internal/cli/tickets.go:1002/1017/1114/1121 (Quelle 'built-in' und das JSON-Feld 'builtin').

ORDER IST EINE FALLE
order() (lane.go:850) haelt mitgelieferte Lanes in der Reihenfolge ihrer Dateinamen und ordnet nur die uebrigen nach 'after:'. Landen critique/optimize/testing als builtin/*.md mit falschem Zahlenpraefix, stehen sie an der falschen Stelle, sobald sie installiert sind. Deshalb 25-critique.md, 26-optimize.md, 27-testing.md - zwischen 20-in-progress und 30-human, wo ihr 'after:' sie ohnehin hinweist.

OFFLINE IST DANN GESCHENKT
Installable() (core/lane/order.go:175) bietet schon heute Builtins() plus ~/.jaira/lanes an, lane.Add() kopiert von dort. Eine eingebettete critique-Lane ist damit ohne jedes Netz per 'jaira lanes add critique' installierbar - 'market adopt' entfaellt fuer sie ganz, und es muss keine Fehlermeldung geschrieben werden (DoD 3).

DIE DREI DATEIEN AUS lanes/ NEHMEN
Vorschlag: lanes/critique.md, lanes/optimize.md, lanes/testing.md werden nach core/lane/builtin/ verschoben, nicht kopiert. Zwei Quellen fuer denselben Prompt laufen auseinander, und der Markt soll laut der Entscheidung der Ort fuer FREMDE Lanes sein. Alte Binaries verlieren dadurch nichts: mit Weg A fragen sie den Tag ihrer eigenen Version ab, und dort liegen die Dateien noch. Im Katalog bleiben changelog-writer und secrets-scan.

WEG A, UND WO ER KLEMMT
release.Current (core/release/release.go:23) wird in cli.Execute (internal/cli/root.go:100) gesetzt, und core/release importiert nichts aus jaira - core/market darf es also lesen, ohne einen Zyklus zu bauen. apiBase() (core/market/market.go:45) haengt '?ref=v<release.Current>' an. Bei 'dev' (Source-Build) gibt es keinen Tag: dann ohne ref, plus eine Zeile auf stderr, dass die Entwicklungsfassung des Katalogs kommt. Offen und in Schritt 9 zu entscheiden: ob der ref auch an ein per JAIRA_MARKET_API gesetztes Ziel gehaengt wird - dagegen spricht nichts, der Testserver ignoriert eine unbekannte Query.
- **2026-09-18 10:40 · Alexander Sacharov** — in-progress, Schritte 1-4: die Aufteilung steht, und der Plan hatte EINE Fundstelle zu wenig.

DRITTE STELLE, die auf Default muss: migrateLegacy (core/lane/lane.go, der Zweig 'kein order-File'). Der Plan nannte nur setUp und internal/tui/defaultboard.go:50. Aber der Fallback dort bedeutet 'was ein altes Verzeichnis ohne order-File implizit hatte', und das waren die zehn - haette er .Builtin weitergelesen, bekaeme jedes migrierte Altboard critique/optimize/testing dazugeschrieben. Genau der Fall, den DoD 4 verbietet, nur an einer Stelle, an der niemand ihn sucht. Die uebrigen acht .Builtin-Stellen bleiben wie geplant.

creator: BeMuCa aus den drei Dateien ENTFERNT. Sie liegen jetzt im Binary, und parse() setzt fuer built-ins ohne creator: 'jaira'. Mit der alten Zeile faellt TestBuiltinDefaultsCreatorToJaira - und die Zeile war auch inhaltlich nicht mehr wahr.

SIEBEN TESTS mussten mit, alle aus demselben Grund: sie schrieben 'built-in' und meinten 'was ein frisches Board bekommt'. Der Helfer builtinIDList filtert jetzt auf Default; neu daneben shippedIDList fuer 'alles, was im Binary steckt' (TestNoBoardReturnsTheOffer braucht das, das Angebot ist jetzt echt groesser). builtinSet in next_test.go lud Load("") - das Angebot, nicht ein Board; mit critique im Angebot war next(in-progress) plotzlich critique. Es laedt jetzt ein leeres Board.

NEBENBEFUND, ungeplant und willkommen: der Lane-Bildschirm der TUI zeigt die drei sofort in seiner 'not on board'-Spalte an - available in internal/tui/lanes.go zieht aus Installable(), und Einbetten genuegte. DoD 2 fuer die TUI ist damit schon erfuellt, ohne eine Zeile dafuer; Schritt 6 bleibt fuer 'jaira lanes' auf der Kommandozeile noetig.
- **2026-09-18 10:47 · Alexander Sacharov** — Die sechs .Builtin-Fundstellen in den Tests, vor Schritt 5 einzeln angesehen (Fund der mitlaufenden Kritik). Ergebnis: KEINE muss geaendert werden, und eine arbeitet fuer uns.

core/lane/shipped_test.go:59 verbietet, dass eine Datei unter lanes/ zu einer Lane mit Builtin=true aufloest. Haette ich critique.md KOPIERT statt verschoben, faende dieser Test es sofort - die Kopie unter lanes/ wuerde auf die eingebettete Lane aufloesen und der Test schluege an. Der Schutz gegen 'zwei Quellen fuer denselben Prompt' steht also schon im Repo; nichts nachzubauen.
core/lane/lane_test.go:207 (eine ueberschreibende Lane ist nicht Builtin), internal/tui/lanes_test.go:77 (auf einem Board ist keine Lane Builtin, es wird nichts mehr injiziert) und :665 (eine entfernte built-in taucht als verfuegbare built-in wieder auf) sprechen alle drei ueber Herkunft, nicht ueber Vorauswahl. Genau die Bedeutung, die Builtin behalten hat.

FREMDBEFUND, NICHT VON DIESEM TICKET VERURSACHT und hier nicht gefixt: wer eine built-in ueberschreibt, indem er z.B. review.md in ~/.jaira/lanes legt, verliert diese Lane auf jedem NEUEN Board. Nachgestellt: ein frisches Board laedt dann als [backlog brainstorm todo pre-process in-progress human signoff done blocked] - review fehlt. Ursache: die ueberschreibende Lane hat Builtin=false (lane_test.go:207 schreibt das ausdruecklich fest), und setUp filterte schon vorher auf Builtin. Mit der Aufteilung ist es unveraendert, weil Default fuer eine Nicht-built-in ohne 'default-board:' ebenfalls false wird. Fix waere, Default von der ueberschriebenen Lane zu erben, dort wo Load Overrides setzt - eine Zeile, aber eine Verhaltensaenderung ausserhalb dieses Tickets.
- **2026-09-18 10:49 · Alexander Sacharov** — Schritt 5: der Test fiel beim ersten Lauf, und er hatte recht - 'jaira lanes add' haengte die Lane ans ENDE der Order-Datei, ungeachtet ihres 'after:'. Die drei landeten hinter done und blocked. Eine Pruefschleife, durch die kein Ticket je laeuft, ist keine installierte Pruefschleife, also ist das kein Testproblem, sondern DoD 3.

Warum es bisher niemandem auffiel: solange jede mitgelieferte Lane ohnehin schon installiert war, kam ueber diesen Weg nur eine eigene Lane, und die haengt man sich selbst zurecht. Mit uninstalliert mitgelieferten Lanes ist Anhaengen falsch.

Gefixt in Add() (core/lane/order.go) mit insertAfterAnchor: NUR die neue id wird gesetzt, der Rest der Order-Datei bleibt Zeichen fuer Zeichen stehen. Bewusst NICHT order() ueber das ganze Board laufen lassen - die Spaltenfolge ist die Anordnung des Benutzers, und eine hinzugefuegte Lane ist kein Anlass, sie neu herzuleiten. Ohne 'after:' oder mit einem Anker, den dieses Board nicht hat, wird weiter angehaengt; das ist dieselbe Aussage, die order() aus einem fehlenden Anker liest.

Die Kette loest sich von selbst auf: critique steht nach in-progress, dann findet optimize seinen Anker critique, dann testing seinen Anker optimize. Kein Sortierdurchlauf noetig.

Das ist eine sichtbare Verhaltensaenderung von 'jaira lanes add' - gehoert in NOTES.md (Schritt 10), zusaetzlich zu den dort schon geplanten drei Zeilen.
- **2026-09-18 10:54 · Alexander Sacharov** — Schritte 6+7, drei Entscheidungen am Fuss von 'jaira lanes', die aus dem Code nicht hervorgehen.

NUR MITGELIEFERTE, NICHT DER GANZE KATALOG. Installable() liefert built-ins PLUS ~/.jaira/lanes; der Fuss filtert auf l.Builtin. Zwei Gruende: eine Lane, die der Benutzer selbst in seinen Katalog gelegt hat, ist ihm bekannt - ihn darauf hinzuweisen ist Laerm; und die Aussage 'no network needed' waere fuer sie zwar auch wahr, aber der Satz, den der Fuss zu tragen hat, ist 'das steckt in deinem Binary'. Der Markt bleibt ausdruecklich draussen: er braucht Netz, und ein Fuss, der ohne Verbindung leer bliebe oder haengt, ist schlimmer als keiner.

EINE ZEILE PRO LANE, DIE KOMMANDOZEILE EINMAL AM SCHLUSS. Erst stand 'jaira lanes add <id>' unter jeder der drei - dreimal dieselbe Zeile unter einer ohnehin zehnzeiligen Tabelle. Jetzt: id + description je Zeile, darunter einmal "Add one with 'jaira lanes add <id>'".

DER FUSS SCHWEIGT, WENN NICHTS FEHLT, und dafuer gibt es einen eigenen Test. Eine Dauerwerbung unter einem Kommando, das man staendig aufruft, wird nach drei Tagen nicht mehr gelesen - und dann fehlt sie genau bei dem einen Board, bei dem sie gezaehlt haette.

Im JSON heisst das Feld 'available' und traegt id, name, description, agentic, model_tier und 'add' mit dem fertigen Kommando. Eine Sitzung, die --json liest, soll den Prosa-Fuss nicht parsen muessen.

NACHTRAG zur frueheren Notiz: der 'not on board'-Streifen der TUI zeigt die drei bereits, ohne Zutun - er zieht aus derselben Installable().
- **2026-09-18 10:58 · Alexander Sacharov** — Funde 5-7 der mitlaufenden Kritik behoben. Was dabei gelernt wurde und aus dem Diff nicht hervorgeht:

FUND 5, creator. Ich hatte 'creator: BeMuCa' durch 'default-board: false' ERSETZT statt die Zeile dazuzuschreiben - zwei Felder ohne jeden Zusammenhang. Die Folge war still und genau die, gegen die das Feld existiert: parse() setzt Creator='jaira', wenn das Feld fehlt UND die Lane built-in ist, und built-in sind die drei jetzt. 'jaira lanes show critique' haette jaira als Autor gemeldet. Ausliefern ist nicht Urheberschaft.
Dabei fiel TestBuiltinDefaultsCreatorToJaira, und die Reparatur ist wichtiger als sie aussieht: der Test darf NICHT auf Default filtern. Default ist 'steht in der Vorauswahl' und hat mit Autorschaft nichts zu tun - das waere derselbe Fehler noch einmal, nur im Test. Er fuehrt jetzt eine benannte Liste, welche mitgelieferte Lane ihren eigenen Autor nennt, und meckert, wenn eine davon nicht mehr mitgeliefert wird.

FUND 6, der ernste. insertAfterAnchor tat bei fehlendem Anker genau das, wogegen sie geschrieben war: anhaengen, hinter done und blocked. Szenario ist real - 'jaira lanes remove in-progress', dann 'jaira lanes add critique'. Die Lane ist dann installiert und unerreichbar, was schlimmer ist als gar nicht installiert, weil 'jaira lanes' sie brav auflistet.
Der Fix ist Paritaet mit order(): vor die terminale Lane, in BEIDEN ankerlosen Faellen, und beim unaufloesbaren Anker mit derselben Warnung, die order() schon ausgibt. Wo eine Lane landet, darf nicht davon abhaengen, ob Load die Reihenfolge herleitet oder Add sie schreibt.
Das kostete eine Signaturaenderung: Add gibt jetzt (string, []string, error) zurueck. Warnungen brauchen einen Kanal, und der Rueckgabewert ist derselbe, den order() schon benutzt. Drei Aufrufer: 'jaira lanes add' druckt sie auf stderr und fuehrt sie im --json unter 'warnings'; die TUI hat nur EINE Meldungszeile, deshalb gewinnt dort die Warnung ueber das blosse 'added <id>' - dass die Lane nicht dort steht, wo ihr after: hinzeigt, ist das, was der Benutzer lesen muss.

DER NEUE TEST BEISST WIRKLICH. Gegengeprueft, indem ich 'at = len(ids)' voruebergehend zurueckgesetzt habe: der Test faellt dann mit [... done blocked critique]. Ein Regressionstest, der die alte Fassung durchlaesst, ist keiner - das kostet zwanzig Sekunden und beantwortet die Frage, statt sie zu glauben.

FUND 7: die Bedingung 'l.After == "" || at < 0' war zusammengelegt und ist jetzt geteilt, weil die beiden Faelle sich unterscheiden - ohne Anker keine Warnung, mit unaufloesbarem Anker eine. Der Einfuegepunkt ist derselbe.
- **2026-09-18 11:26 · Alexander Sacharov** — Schritte 8+9, Weg A. Drei Dinge, die der Diff nicht sagt.

DER REF WIRD GESETZT, NICHT ANGEKLEBT (Fund 2 der Kritik). '?ref=' als String an apiBase() zu haengen bricht, sobald die Adresse schon eine Query traegt - und der Testserver setzt eine, das ist also der Normalfall und keine Ecke. Heraus kaeme '?x=1?ref=v0.3.0', ein einziger unsinniger Parameter. Jetzt url.Parse + q.Set + Encode, und ein eigener Test (TestRefIsSetOnAnAddressThatAlreadyHasAQuery) haelt fest, dass die vorhandene Query dabei erhalten bleibt.

DIE DEV-ANSAGE WIRD ZURUECKGEGEBEN, NICHT GEDRUCKT (Fund 3). Sie gehoert nicht nach apiBase(): das ist eine reine Funktion, sie laeuft bei jeder Anfrage, und eine Zeile von dort erschiene mitten im --json. Vorbild ist Overridden() direkt darueber - Text zurueckgeben, drucken tut das Kommando. Heisst jetzt Unpinned() und steht im JSON unter 'unpinned' neben 'override'.

WARUM release.Current und nicht ein durchgereichter Parameter: core/release importiert nichts aus jaira, core/market darf es also lesen, ohne einen Zyklus zu bauen. cli.Execute (internal/cli/root.go) setzt die Variable ohnehin schon fuer die TUI.

Die in Schritt 8 offen gelassene Frage - geht der ref auch an ein per JAIRA_MARKET_API gesetztes Ziel - ist mit JA entschieden. Ein Codepfad statt zwei; ein selbst gehosteter Spiegel ist genauso versionierbar; und ein Testserver ignoriert eine unbekannte Query ohnehin. Ein zweiter Pfad waere genau der, den kein Test je durchlaeuft.

Am echten Binary nachgestellt, wie DoD 6 es verlangt: mit -ldflags '-X main.version=0.1.4' fragt es /contents/lanes?ref=v0.1.4 ab, der Source-Build fragt ohne ref und sagt die Zeile dazu.

WERKZEUG-DELLE, fuer die naechste Sitzung: ein 'pkill -f ...' am Anfang einer Kommandokette hat die ganze Kette mitgerissen (Exit 144), noch bevor die jaira-Schreibvorgaenge dahinter liefen. Die Haken standen danach nicht auf dem Ticket. Aufraeumbefehle nicht mit Ticket-Schreibvorgaengen in eine Zeile.
- **2026-09-18 12:04 · Alexander Sacharov** — Schritt 10, was die Zeilen tragen und warum fuenf und nicht drei. Geplant waren drei (Lanes im Binary, Fuss von 'jaira lanes', versionsgebundener Katalog). Dazugekommen sind zwei, weil unterwegs zwei weitere Dinge sichtbar wurden:
- 'jaira lanes add' setzt jetzt nach dem Anker statt anzuhaengen. Das ist eine geaenderte Wirkung eines bestehenden Aufrufs und damit client-facing nach der Regel in CLAUDE.md, auch wenn niemand danach gefragt hat.
- 'default-board:' ist ein neues Frontmatter-Feld in Lane-Dateien. Ebenfalls ausdruecklich client-facing, und wer eigene Lanes schreibt, will davon wissen.
In der Katalog-Zeile steht jetzt auch, dass der ref an ein per JAIRA_MARKET_API gesetztes Ziel geht und dessen eigene Query-Parameter erhalten bleiben - ohne das koennte jemand mit eigenem Spiegel es fuer einen Fehler halten, wenn dort ploetzlich ein ref ankommt.

Format gegengeprueft, nicht nur gelesen: Notes() aus core/release gegen die Datei laufen lassen, Ergebnis 'Unreleased hat 5 Aenderungen'. Der Parser ist ein Zeilenscan, eine umgebrochene Zeile waere still zu zwei halben Aenderungen geworden - das sieht man der Datei im Editor nicht an, dem Parser schon.

Schritt 11: go build, go vet und 'go test -count=1 ./...' (ohne Cache) alle sauber.
- **2026-09-18 12:12 · Alexander Sacharov** — critique (Lane, erster Durchgang): fuenf Funde, zwei davon am gebauten Binary reproduziert, nicht nur gelesen.

REPRODUZIERT 1 - der Fuss bietet Entferntes wieder an. Frisches Testboard, 'jaira lanes remove signoff', dann 'jaira lanes': signoff steht unter 'Shipped with this binary, not on this board'. Der Filter 'l.Builtin' (internal/cli/tickets.go:1002) meint die falsche Menge; gemeint ist 'mitgeliefert, aber nicht in der Vorauswahl', also 'l.Builtin && !l.Default'. Genau dafuer wurde Default in diesem Diff eingefuehrt - die Stelle liest es nur nicht.

REPRODUZIERT 2 - der neue Fuss laedt in eine Falle. 'jaira lanes add testing' auf einem frischen Board (der Fuss bietet die drei EINZELN an und sagt 'Add one with jaira lanes add <id>'): testing landet hinter signoff. Ein Test-Lane nach der menschlichen Abnahme. Ursache: insertAfterAnchor kennt nur zwei Faelle - Anker da, oder vor die terminale Lane. Der dritte Fall ist neu durch dieses Ticket: der Anker ist selbst eine mitgelieferte, noch nicht installierte Lane. Die Kette loest sich auf, wenn man sie ueber Installable() weiterverfolgt: optimize -> critique -> in-progress. Der Fund-6-Fix war richtig, er ist nur eine Ebene zu flach.
Der DoD-3-Test deckt das nicht, weil er die drei in der richtigen Reihenfolge installiert - dann findet jede ihren Anker schon auf dem Board.

WAS ICH STEHEN LASSE, mit Grund - damit der naechste Durchgang es nicht neu aufmacht:
- Fund 10 ('add' als fertige Kommandozeile im JSON): bewusst und in NOTES.md ausdruecklich verkauft ('so a session reading JSON does not have to parse the prose'). Eine Redundanz, die ihren Zweck nennt, ist kein Fund.
- Fund 12 (os.Stderr statt cmd.ErrOrStderr in tickets.go:1053): die Zeile ist ALT, nicht aus diesem Diff. Und die Kopplung ist nicht real - nachgemessen mit 'go test -run TestLanesJSONCarriesAvailable -v': dieser Test gibt gar keine Warnung aus, nur der Prosa-Test tut es. Eine Vereinheitlichung ist ein eigenes Ticket, keine Aenderung an dieser Lane.
- terminalIDIndex neben terminalIndex (lane.go:938): die Signaturen unterscheiden sich ([]string+Set gegen []*Lane), eine Zusammenlegung waere mehr Code als die sechs Zeilen. Absichtlich gelassen.
- Die drei verschobenen Lane-Dateien: wortgleich uebernommen, je eine Zeile dazu. Nichts zu sagen.
- **2026-09-18 12:20 · Alexander Sacharov** — in-progress, zweiter Durchgang: die beiden reproduzierten Funde der Kritik behoben. Was der Diff nicht sagt:

FUND 'REPRODUZIERT 1' war ein Ein-Bit-Fehler mit einer inhaltlichen Aussage dahinter. Der Fuss filterte auf Builtin und meinte damit 'mitgeliefert'; richtig ist 'Builtin && !Default' - mitgeliefert UND nicht in der Vorauswahl. Der Unterschied ist nicht kosmetisch: eine Lane, die jedes frische Board bekommt und die dieses Board nicht mehr hat, wurde ENTFERNT. Sie unter jedem 'jaira lanes' wieder anzubieten stellt eine Entscheidung des Benutzers dauerhaft in Frage. Das Angebot gilt dem, was ein Board nie hatte.

FUND 'REPRODUZIERT 2' war tiefer als er aussah. insertAfterAnchor kannte zwei Faelle; der dritte ist jetzt der NORMALFALL, nicht die Ecke: der Fuss bietet die drei einzeln an, also ist 'jaira lanes add testing' auf einem frischen Board der uebliche Weg, und testings Anker optimize ist selbst noch nicht installiert. Die Kette loest sich ueber Installable() auf (testing -> optimize -> critique -> in-progress). Neue Funktion anchorIndex, mit besuchter-Menge gegen eine Kette, die sich im Kreis dreht - zwei nicht installierte Lanes, die einander nennen, haetten sonst eine Endlosschleife ergeben.

DIE PARITAET MIT order() IST DABEI BEWUSST AUFGEGEBEN, und das ist die Entscheidung, die man spaeter anfechten koennte. Die frueherer Notiz sagte: eine Lane darf nicht anders landen, je nachdem ob Load die Reihenfolge herleitet oder Add sie schreibt. Das gilt weiter fuer den unaufloesbaren Anker. Fuer die Kette kann es nicht gelten: order() sieht nur die Lanes AUF dem Board und hat gar keine Kette zu verfolgen, Add haelt das Angebot in der Hand, aus dem die Lane kommt. Add ist also besser informiert, nicht abweichend. Im Kommentar an insertAfterAnchor steht das so.

FOLGE FUER EINEN BESTEHENDEN TEST, den ich umgeschrieben statt repariert habe: TestLanesAddWithoutItsAnchorStaysBeforeTheTerminalLane (in-progress entfernt, dann critique dazu) erwartete eine Warnung. Mit der Kette landet critique jetzt hinter pre-process - in-progress ist ja noch im Angebot und nennt pre-process - und es gibt nichts zu warnen. Das ist die bessere Stelle, also ist der Test veraltet und nicht der Code. Er heisst jetzt TestLanesAddFollowsTheChainPastARemovedAnchor und prueft genau das. Damit der Warnpfad nicht ungeprueft bleibt, steht daneben neu TestLanesAddWithAnAnchorNowhereStaysBeforeTheTerminalLane mit einer eigenen Katalog-Lane, deren after: auf 'nowhere' zeigt.

BEIDE NEUEN TESTS GEGENGEPRUEFT, dass sie wirklich beissen: mit zurueckgesetztem Filter meldet der eine 'available carries signoff', mit abgeschalteter Kettenverfolgung meldet der andere 'testing landed behind signoff'. Ein Regressionstest, den die alte Fassung besteht, ist keiner.

AM ECHTEN BINARY NACHGESTELLT, wie die Kritik es getan hatte: frisches Board, 'jaira lanes add testing' -> [... in-progress testing human ...] ohne Warnung; 'jaira lanes remove signoff' -> der Fuss nennt nur critique, optimize, testing.
- **2026-09-18 12:26 · Alexander Sacharov** — critique (Lane, zweiter Durchgang): drei Punkte, zwei am gebauten Binary reproduziert. Dazu die Entscheidung zu Punkt A, die der Dispatcher mir gegeben hat.

ENTSCHIEDEN - PUNKT A, die stille Kettenaufloesung BLEIBT still. Der Vorschlag war, auch bei at >= 0 zu warnen, wenn der aufgeloeste Anker nicht l.After ist. Nachgemessen, warum das nicht geht: auf einem frischen Board loest 'jaira lanes add testing' die Kette testing -> optimize -> critique -> in-progress auf und wuerde damit warnen - und genau dieser Aufruf ist der Weg, den der neue Fuss von 'jaira lanes' selbst bewirbt. Eine Warnung auf dem beworbenen Normalweg ist Laerm, und Laerm erzieht dazu, Warnungen zu ueberlesen. Add kann die beiden Faelle auch nicht trennen: 'noch nie installiert' und 'entfernt' sehen von Add aus identisch aus (dasselbe, was Punkt C am Kommentar in tickets.go aufzeigt). Also entweder immer warnen oder nie - und nie ist richtig, denn die Platzierung ist in beiden Faellen die richtige.

Was dabei aber wirklich fehlt, ist nicht eine Warnung, sondern die Auskunft WO die Lane gelandet ist: 'jaira lanes add' sagt heute nur 'added <id> to this project (<pfad>)'. Der geloeschte Test TestLanesAddWithoutItsAnchorStaysBeforeTheTerminalLane trug diese Ansage, die Kettenaufloesung hat sie ersatzlos mitgenommen. Der Nachbar gehoert in die normale Erfolgszeile, nicht auf den Warnkanal - damit erfaehrt auch der, der in-progress entfernt hat, wo critique hingekommen ist, ohne dass der Normalweg Warnungen wirft.

PUNKT B am Binary reproduziert, nicht nur gelesen: '.jaira lanes remove brainstorm', 'remove backlog', dann 'lanes add brainstorm' -> 'jaira: warning: lane brainstorm: anchor "" is not on this board'. Leere Anfuehrungszeichen, weil anchorIndex die Laufvariable 'anchor' bis ans Kettenende schiebt und insertAfterAnchor sie in die Meldung setzt. Der Benutzer hat 'backlog' geschrieben und bekommt einen Namen genannt, den er nie geschrieben hat.

PUNKT C am Binary reproduziert: 'lanes add critique', 'lanes remove critique', 'jaira lanes' -> critique steht wieder im Angebot. Das VERHALTEN ist vertretbar (eine Katalog-Lane ist eben immer nachinstallierbar), der Kommentarsatz ist es nicht: 'The offer is for the lanes a board never had' ist nachweislich falsch, und ein Kommentar, der eine Eigenschaft behauptet, die der Code nicht hat, ist schlimmer als keiner. Nur der Satz aendert sich, nicht der Filter.

WAS ICH STEHEN LASSE, damit der dritte Durchgang es nicht neu aufmacht:
- Die Warnungsweitergabe in die TUI: internal/tui/lanes.go:460 und :511 nehmen warnings[0] in die Meldungszeile. Nachgesehen, nicht angenommen - die neue Signatur verliert nichts.
- Alles aus dem ersten Durchgang, was dort mit Grund stehen blieb (Fund 10, Fund 12, terminalIDIndex, die drei verschobenen Lane-Dateien). Unveraendert gueltig, siehe Notiz vom 12:12.
- **2026-09-18 12:34 · Alexander Sacharov** — in-progress, dritter Durchgang: die drei Punkte der zweiten Kritik. Was der Diff nicht sagt:

PUNKT A als VIERTER RUECKGABEWERT von lane.Add, nicht als Struct. Ein Struct waere die sauberere Signatur, kostet aber einen neuen exportierten Typ fuer eine einzige Zeichenkette, und der Aufrufer, der sie braucht, ist genau einer. Die beiden TUI-Aufrufe (internal/tui/lanes.go:446 und :499) ignorieren sie ABSICHTLICH: dort sieht man die Position, die Liste wird neu gezeichnet und der Cursor springt auf die Lane. Nur die Kommandozeile ist blind, und nur sie sagt es jetzt.

EINE STELLE, DIE KOMISCH AUSSIEHT UND RICHTIG IST: bei unaufloesbarem Anker steht die Erfolgszeile 'added brainstorm to this project after signoff' direkt ueber der Warnung, dass backlog nicht da ist. Am Binary nachgestellt. Ueberlegt, den Nachbarn in diesem Fall wegzulassen - nicht getan: genau dann will der Benutzer wissen, WO die Lane statt dessen liegt, und die Warnung darunter erklaert das 'warum'. Wer das spaeter 'aufraeumt', nimmt die Auskunft im einzigen Fall weg, in dem sie nicht selbstverstaendlich ist.

PUNKT B hat anchorIndex den zweiten Rueckgabewert ganz genommen. Nachgesehen, ob ihn sonst jemand liest: nein, er existierte nur fuer die Meldung. Die Funktion gibt jetzt int zurueck.

PUNKT C stand ZWEIMAL da. Der falsche Satz 'The offer is for the lanes a board never had' klebte auch im Kopf von TestLanesDoesNotOfferBackALaneTheBoardRemoved in internal/cli/lanes_test.go - ein Testkommentar, der eine falsche Eigenschaft behauptet, fuehrt den naechsten Leser genauso in die Irre wie der im Produktionscode.

NOTES.md BEKAM KEINE NEUE ZEILE, sondern zwei korrigierte. Die Fuss-Zeile verkaufte dieselbe falsche Zusage wie der Kommentar ('It offers only lanes a board never had') - das ist kein Kommentarfehler mehr, sondern ein Versprechen an den Benutzer, und die Zeile ist unveroeffentlicht, also korrigierbar. Die 'lanes add'-Zeile trug die Platzierung schon, also gehoert die neue Auskunft dort hinein und nicht in eine zweite Zeile: eine Aenderung ist eine Zeile.

BEIDE NEUEN TESTS GEGENGEPRUEFT: mit zurueckgedrehtem Code meldet der eine 'the add must say which lane testing landed after', der andere 'anchor ""'. Am echten Binary nachgestellt, alle drei Faelle.
- **2026-09-18 12:37 · Alexander Sacharov** — critique (Lane, dritter Durchgang): EIN Fund, am gebauten Binary reproduziert. Die drei Punkte des zweiten Durchgangs sind erledigt und werden nicht neu aufgemacht.

NACHGEPRUEFT, nicht nur gelesen. Punkt A: frisches Testboard, 'jaira lanes add testing' -> 'added testing to this project after in-progress (<pfad>)'. Punkt B: 'remove brainstorm', 'remove backlog', 'add brainstorm' -> Warnung nennt jetzt anchor "backlog", den Namen, den der Benutzer geschrieben hat, nicht mehr die leere Zeichenkette. Punkt C: der falsche Satz ist in internal/cli/tickets.go und im Kopf von TestLanesDoesNotOfferBackALaneTheBoardRemoved durch die zutreffende Fassung ersetzt.

FUND - derselbe Fehler wie Punkt C, an drei Stellen, die Punkt C nicht angesehen hat. Punkt C war: ein Kommentar behauptet eine Eigenschaft, die der Code nicht hat. Der Fix zu Fund 6 / REPRODUZIERT 2 hat genau so eine Behauptung an drei weiteren Stellen falsch gemacht, und eine davon ist nicht ein Kommentar, sondern Hilfetext.

- internal/cli/lanes.go:100, das Long von 'lanes add': 'Adds the named lane to this board, appending it to the column order'. Am Binary nachgestellt: 'jaira lanes add --help' sagt das heute woertlich. Angehaengt wird seit dem Fix nichts mehr - die Lane landet, wohin ihre after:-Kette zeigt.
- internal/tui/lanes.go:438, Kommentar an addFromCatalogue: 'appending it at the end of the order'.
- docs/COMMANDS.md:153, dieselbe Zeile in der Kommandotabelle.

WARUM DAS IN DIESEM DURCHGANG UND NICHT SPAETER: der Hilfetext steht auf dem Weg, den DoD 2 verlangt. Wer 'jaira init' laeuft, den Fuss von 'jaira lanes' liest und dem 'Add one with jaira lanes add <id>' folgt, fragt als naechstes '--help' - und bekommt dort genau den Glauben bestaetigt, gegen den der zweite Durchgang gebaut hat: dass die Lane hinten angehaengt wird. Der Weg ohne Vorwissen ist damit nicht gehbar, sondern irrefuehrend beschriftet.

NOTES.md braucht dafuer KEINE neue Zeile: die 'lanes add'-Zeile sagt die Platzierung schon richtig. Es ist der Hilfetext, der ihr widerspricht.

WAS ICH STEHEN LASSE, damit der vierte Durchgang es nicht neu aufmacht:
- Alles aus Durchgang 1 (Fund 10, Fund 12, terminalIDIndex, die drei verschobenen Lane-Dateien) und aus Durchgang 2 (die Warnungsweitergabe in die TUI). Unveraendert gueltig, siehe Notizen 12:12 und 12:26.
- Der vierte Rueckgabewert von lane.Add statt eines Structs: Begruendung steht in der Notiz von 12:34, ein exportierter Typ fuer eine Zeichenkette mit einem Leser. Angenommen.
- Die Erfolgszeile ueber der Warnung bei unaufloesbarem Anker ('after signoff' plus Warnung zu backlog): am Binary gesehen, liest sich seltsam und ist richtig. Begruendung 12:34, geschlossen.
- **2026-09-18 18:20 · Alexander Sacharov** — in-progress, vierter Durchgang: der Fund der dritten Kritik behoben. Drei Dinge, die aus dem Diff nicht hervorgehen:
- Eine VIERTE Fundstelle desselben veralteten Satzes existiert und wurde ABSICHTLICH nicht angefasst: internal/cli/lanes_test.go:955, Test 'TestLanesAddAfterRemoveAppendsAtEnd' - Name und Kommentar behaupten beide 'appended at the end of the order'. Der Test ist gruen, aber aus dem falschen Grund: 'blocked' landet nicht durch Anhaengen am Ende, sondern weil insertAfterAnchor fuer eine Lane ohne aufloesbaren Anker vor die terminale Lane parkt und 'blocked' selbst terminal ist. Der Auftrag dieser Runde war auf die drei genannten Stellen begrenzt, also gemeldet statt stillschweigend mitgenommen. Wer sie anfasst: Name und Kommentar aendern, die Assertion nicht - sie prueft weiter das richtige Ergebnis.
- Grep-Abdeckung: gesucht wurde nach 'append', 'end of the order' und 'column order' ueber *.go und *.md. docs/AGENTS.md:279 nennt 'jaira lanes add' zwar, behauptet aber nirgends Anhaengen - keine Fundstelle. core/lane/order.go:263 sagt 'rather than at the end of the board' und beschreibt damit das NEUE Verhalten korrekt.
- Keine neue NOTES.md-Zeile geschrieben, und das ist eine Entscheidung: core/release/NOTES.md:19 sagt dem Benutzer bereits, wohin 'jaira lanes add' die Lane setzt. Der Hilfetext holt nur nach, was dort schon steht - eine zweite Zeile fuer dieselbe Aenderung wuerde der Zeilenscan als zwei Aenderungen lesen.
- **2026-09-18 18:23 · Alexander Sacharov** — critique (Lane, vierter Durchgang): EIN Fund, plus zwei Entscheidungen, die der Dispatcher mir vorgelegt hat.
- FUND: internal/cli/lanes_test.go:955, TestLanesAddAfterRemoveAppendsAtEnd. Der Worker hat die Stelle gemeldet statt sie stillschweigend mitzunehmen - richtig gemeldet, aber sie gehoert in diese Runde und nicht in ein eigenes Ticket. Sie ist kein neuer Fund, sondern der unreparierte Rest von Fund 3: derselbe Satz, vierte Stelle. Ihn liegen zu lassen hiesse, einen gruenen Test im Baum zu behalten, dessen Name und dessen Fehlermeldung beide das Gegenteil dessen behaupten, was die Assertion schuetzt - der naechste, der insertAfterAnchor anfasst, liest den Namen und glaubt, Anhaengen sei abgesichert. Kosten der Reparatur: eine Umbenennung und zwei Zeilen Prosa.
- Am Code nachgeprueft statt geglaubt: core/lane/order.go:263ff, Kommentar und Code von insertAfterAnchor. Die Diagnose des Workers stimmt - 'kein aufloesbarer Anker' parkt vor der terminalen Lane, und 'blocked' ist die terminale. Die Assertion ids[len-1]=='blocked' bleibt darum richtig.
- ENTSCHIEDEN, keine neue NOTES.md-Zeile: der Worker hat recht und die Begruendung traegt. core/release/NOTES.md:19 sagt dem Leser bereits, wohin 'jaira lanes add' die Lane setzt und dass die Kette durch nicht installierte Lanes verfolgt wird; die Hilfetexte holen nur nach, was dort schon steht. Eine zweite Zeile wuerde nach der Regel in CLAUDE.md als zweite Aenderung gelesen und dem Leser nichts nennen, das er anders tun soll - das ist der Test, den die Datei an eine Zeile stellt. Geschlossen, nicht wieder aufmachen.
- Nicht neu aufgemacht und weiterhin gueltig: Fund 10, Fund 12, terminalIDIndex und die drei verschobenen Lane-Dateien (Notiz 12:12), die Warnungsweitergabe in die TUI (Notiz 12:26), der vierte Rueckgabewert von lane.Add und die Erfolgszeile ueber der Warnung (Notiz 12:37). Die drei Stellen aus Durchgang 3 sind repariert - internal/cli/lanes.go:100, internal/tui/lanes.go:438, docs/COMMANDS.md:153 im Diff gelesen, alle drei sagen jetzt 'where its after: chain points'.
- **2026-09-18 20:54 · Alexander Sacharov** — in-progress, fuenfter Durchgang: der Fund der vierten Kritik behoben - und die Begruendung, die Kritik UND Dispatcher mitgaben, ist am Code falsch.

- Kritik und Dispatcher sagten beide: 'blocked' landet hinten, weil insertAfterAnchor bei UNAUFLOESBAREM Anker vor die terminale Lane parkt und 'blocked' selbst terminal ist. Beide Haelften stimmen nicht. core/lane/builtin/60-blocked.md:4 traegt 'after: done', und 'done' steht auf dem Board - der Anker ist also aufloesbar, anchorIndex (core/lane/order.go:315) findet ihn und gibt index(done)+1 zurueck. Und terminal ist nicht 'blocked' (60-blocked.md:7 'terminal: false'), sondern 'done' (50-done.md:7 'terminal: true').
- Gegenprobe, die es entscheidet, ohne etwas zu glauben: waere der Fallback-Weg gelaufen, haette terminalIDIndex den Index von 'done' geliefert und 'blocked' VOR 'done' eingesetzt - dann waere 'blocked' nicht mehr letzte Lane und die Assertion ids[len-1]=='blocked' waere rot. Sie ist gruen. Also lief der Anker-Weg.
- Darum heisst der Test NICHT TestLanesAddAfterRemoveLandsBeforeTheTerminalLane, wie vorgeschlagen - der Name waere die zweite falsche Behauptung an derselben Stelle. Er heisst TestLanesAddAfterRemoveLandsAfterItsAnchor, und Kommentar wie Fehlermeldung nennen 'done' als den Anker, dem die Lane folgt.
- Was der Kommentar zusaetzlich festhaelt, weil es die Stelle verwirrend macht: 'blocked' landet zufaellig an derselben Position, die Anhaengen erzeugt haette. Die Assertion bleibt trotzdem aussagekraeftig - genau weil der Fallback eine ANDERE Position ergaebe.
- Sonst nichts angefasst: keine NOTES.md-Zeile (die Entscheidung von 12:12/18:23 gilt, und eine Testumbenennung ist ohnehin nicht client-facing), keine weitere Fundstelle. go build, go vet und 'go test -count=1 ./...' gruen.
- **2026-09-18 20:58 · Alexander Sacharov** — critique (Lane, fuenfter und letzter Durchgang): KEIN Fund. Die Schleife endet hier.
- Der eine Punkt aus Durchgang 4 ist erledigt und richtig erledigt: internal/cli/lanes_test.go:955ff heisst jetzt TestLanesAddAfterRemoveLandsAfterItsAnchor, Kommentar und Fehlermeldung nennen 'done' als Anker, die Assertion ids[len-1]=='blocked' ist unveraendert.
- Die Gegenrede des Workers aus der Notiz 20:54 am Code nachgeprueft statt geglaubt, und sie stimmt in allen drei Punkten: core/lane/builtin/60-blocked.md:4 traegt 'after: done'; core/lane/builtin/50-done.md:7 traegt 'terminal: true', 60-blocked.md:7 traegt 'terminal: false'; terminalIDIndex (core/lane/order.go:349) gibt den Index der terminalen Lane selbst zurueck, nicht den danach. Der Fallback haette 'blocked' also VOR 'done' gesetzt und die Assertion rot gemacht - gelaufen ist der Anker-Weg ueber anchorIndex (order.go:315). Meine eigene Diagnose aus Durchgang 4 war an dieser Stelle falsch; der neue Name ist der richtige, der von mir vorgeschlagene waere die zweite falsche Behauptung gewesen.
- Gegengeprueft, dass Fund 3 jetzt wirklich erschoepft ist: grep ueber *.go und *.md nach 'append(s|ing) (it) (at|to) the end' und 'end of the (column) order' findet ausserhalb der Ticket-Prosa und .planning/ keine Fundstelle mehr.
- go build, go vet und 'go test -count=1 ./...' selbst gelaufen, alles gruen.
- Nicht neu aufgemacht und weiterhin gueltig: Fund 10, Fund 12, terminalIDIndex neben terminalIndex und die drei verschobenen Lane-Dateien (Notiz 12:12), die Warnungsweitergabe in die TUI (Notiz 12:26), der vierte Rueckgabewert von lane.Add und die Erfolgszeile ueber der Warnung (Notiz 12:37), die Entscheidung gegen eine weitere NOTES.md-Zeile (Notiz 18:23).
- **2026-09-18 21:04 · Alexander Sacharov** — optimize (Lane): EINE Aenderung, alles andere geprueft und bewusst stehen gelassen.

GEAENDERT - internal/tui/lanes.go: addFromCatalogue und addAvailable trugen denselben Sechszeiler samt identischem Kommentar ('One message line, and an anchor that could not be found...'), beide Male nach einer Zuweisung von ls.msg, die die naechste Zeile sofort ueberschrieb. Jetzt ein Helfer addedMsg(id, warnings) und je eine Zeile an beiden Stellen. Verhalten identisch: bei leeren warnings 'added <id>', sonst 'added <id> — <warnings[0]>'. Der Kommentar steht nur noch einmal, an der Funktion, die die Entscheidung traegt.

DUPLIKATE GESUCHT, NICHT GERATEN:
- terminalIDIndex (order.go:349) neben terminalIndex (lane.go:938): echte Doppelung ueber zwei Typen ([]string+Set gegen []*Lane). Bleibt - Durchgang 1 hat sie mit Grund geschlossen (Notiz 12:12), und eine Vereinheitlichung hiesse, ids in []*Lane aufzuloesen, um eine Schleife von vier Zeilen zu sparen.
- Gegengeprueft, dass die beiden Fallbacks wirklich dasselbe tun, statt es zu glauben: order() rechnet 'idx = terminalIndex(out) - 1' und danach 'at = idx + 1', insertAfterAnchor setzt 'at = terminalIDIndex(ids, set)'. Gleiche Position. Ebenso beim aufgeloesten Anker: order() 'indexOf + 1', anchorIndex 'i + 1'. Kein Auseinanderlaufen.
- Die Unpinned()-Notiz steht zweimal in internal/cli/market.go (list und adopt). Das spiegelt exakt die schon vorhandene Overridden()-Doppelung zwei Zeilen darueber - ein Helfer fuer beide waere eine Aenderung an unberuehrtem Code.

TOTER CODE: keiner gefunden. Alle neuen Bezeichner haben Leser - Default (setUp, migrateLegacy, tickets.go:1010, defaultboard.go:50), anchorIndex/terminalIDIndex/insertAfterAnchor (order.go), pinnedRef/Unpinned (market.go, cli/market.go). anchorIndex hat seinen zweiten Rueckgabewert schon in Durchgang 3 verloren.

KOSTEN: lane.Installable() laeuft jetzt bei jedem 'jaira lanes'. Nachgesehen statt angenommen (order.go:175ff): embedded Builtins plus ein Glob ueber UserLanesDir, kein Netz. Die Startzeit-Regel aus CLAUDE.md ist nicht beruehrt.

STEHEN GELASSEN, damit die naechste Runde es nicht neu aufmacht:
- Die Nachsuche nach 'after' in lane.Add (order.go, Schleife ueber ids nach dem Einfuegen): insertAfterAnchor kennt die Position bereits, aber sie herauszureichen kostet einen dritten Rueckgabewert an einer Funktion mit einem Aufrufer, um eine Schleife ueber dreizehn Eintraege auf einem einmaligen Pfad zu sparen. Kein Gewinn.
- 'release.Current == \"\"' in pinnedRef: release.go:23 setzt 'dev' als Vorgabe, die leere Zeichenkette entsteht nur ueber -ldflags. Der Zweig ist Vorsicht, keine unerreichbare Bedingung - entfernen wuerde ref='v' erzeugen.
- Alles aus den Durchgaengen 1-5, das dort mit Begruendung geschlossen wurde.

go build, go vet und 'go test -count=1 ./...' nach der Aenderung gruen.
- **2026-09-18 21:09 · Alexander Sacharov** — testing, alle drei Durchgaenge gruen - was genau lief und was dabei zu sehen war.

1 GATES. go build ./... und go vet ./... ohne Ausgabe. 'go clean -testcache && go test -race ./...' RC=0, alle 28 Pakete ok (core/lane 10.6s, internal/cli 55.4s, internal/tui 146.3s).

2 DIE FORDERUNG, Punkt fuer Punkt am Arbeitsbaum nachgesehen, nicht am Outcome.
DoD 1: die Notiz 'ENTSCHIEDEN von Alex am 2026-09-18' steht auf dem Ticket.
DoD 2: internal/cli/tickets.go, Filter 'l.Builtin && !l.Default' vorhanden; die vier genannten Tests laufen einzeln gruen.
DoD 3: TestLanesAddInstallsTheReviewLoopWithoutNetwork und TestLanesAddFollowsAnAnchorThatIsItselfUninstalled gruen.
DoD 4: TestFreshBoardGetsOnlyTheDefaultLanes gruen.
DoD 5: unter '## Unreleased' stehen genau fuenf '- '-Zeilen.
DoD 6: TestListPinsTheCatalogueToTheRunningTag, TestDevBuildSendsNoRefAndSaysSo, TestRefIsSetOnAnAddressThatAlreadyHasAQuery gruen.

3 FUNKTION, an einer leeren Testdoska mit dem gebauten Binary.
- 'jaira init' legt genau zehn Lane-Dateien an, ohne critique/optimize/testing.
- 'jaira lanes' schliesst mit 'Shipped with this binary, not on this board - no network needed:' und nennt die drei mit Beschreibung plus der Zeile, die sie holt.
- 'JAIRA_MARKET_API=http://127.0.0.1:1/dead jaira lanes add testing' installiert sie ohne Netz: 'added testing to this project after in-progress'. Die Reihenfolge danach: backlog brainstorm todo pre-process in-progress testing human review signoff done blocked - testing sitzt im Fluss, nicht hinter signoff.
- Versionsbindung am Binary nachgestellt: 'go build -ldflags "-X main.version=0.1.4"' gegen einen lokalen HTTP-Server, der die Anfrage mitschreibt. Der Server sah 'ASKED /contents/lanes?ref=v0.1.4', und die dev-Ansage blieb aus. Das Quell-Binary (dev) sagt umgekehrt 'note: this build reports version dev, ... comes from the development branch'.

WAS DABEI AUFFIEL, kein Fund: 'jaira lanes market' listet an diesem Quell-Binary critique und optimize weiterhin. Das ist richtig so - market liest das Verzeichnis lanes/ des Default-Branches, und dort liegen die drei noch, weil dieser Zweig nicht gemerged ist. Die NOTES-Zeile 'jaira lanes market no longer offers them' gilt ab dem Tag, der den Umzug enthaelt, und genau dorthin zeigt der neue ?ref=.

Nicht zu diesem Ticket: '?? .jaira/milestones/' liegt unverfolgt im Baum, stammt aus 0.3.0.
- **2026-09-19 15:38 · Alexander Sacharov** — Nachpruefung vor der Menschen-Lane (nichts neu gebaut, nur geprueft): go build ./... , go vet ./... und go test ./... -count=1 im Zweig feat/VSC1GW alle sauber, RC=0, keine uebersprungenen Pakete. Alle 13 in den DoD-Beweisen genannten Tests existieren und laufen; alle Datei-Beweise stimmen am Arbeitsbaum (internal/cli/tickets.go Filter 'l.Builtin && !l.Default', core/lane/order.go anchorIndex, internal/cli/lanes.go 'after' in Erfolgszeile und JSON, core/lane/builtin/25-/26-/27- mit 'default-board: false', core/market/market.go apiBase()/pinnedRef()/Unpinned(), fuenf Zeilen unter '## Unreleased'). Einzige Abweichung: der Beweis der letzten Planzeile nennt internal/cli/lanes_test.go:954, die Funktion TestLanesAddAfterRemoveLandsAfterItsAnchor steht auf 962 - Zeilendrift, der Test ist der richtige.

Zu Frage (2), Empfehlung des Dispatchers, keine Entscheidung: so lassen, market NICHT zusaetzlich filtern. Begruendung: die Bindung ?ref=v<version> aus diesem Ticket loest genau dieses Problem schon. Ein ausgeliefertes Binary fragt den Katalog SEINES Tags ab - der Tag, der diesen Zweig enthaelt, hat lanes/critique.md nicht mehr, also bietet market dort nichts doppelt an; und ein 0.3.0-Binary, das die drei Lanes nicht eingebettet traegt, bekommt sie aus dem Katalog voellig zu Recht. Doppelt angeboten wird nur auf einem dev-Build, der ohne ref den HEAD von master liest, und auch nur bis dieser Zweig auf master ist. Das Fenster ist also 'jaira-Entwickler, bis zum Merge'. Ein Filter kostet dagegen eine neue Abhaengigkeit von core/market auf lane.Builtins, Test, NOTES-Zeile - und nimmt einen Fall weg, den es wirklich gibt: eine Katalogfassung einer Lane bewusst gegen die eingebettete zu holen. An der Scope-Regel aus CLAUDE.md gemessen ist das Wachstum fuer ein Problem, das der Tag von selbst beendet.
- **2026-09-19 15:54 · Alexander Sacharov** — ENTSCHIEDEN von Alex am 2026-09-19, zur zweiten Frage der human-Lane: 'jaira lanes market' filtert eingebettete Lanes NICHT zusaetzlich heraus. Es bleibt bei dem, was dieses Ticket gebaut hat.

Begruendung, an der Scope-Regel aus CLAUDE.md gemessen: der Pin '?ref=v<version>' aus diesem Ticket loest den Fall schon. Ein Release-Binary fragt den Katalog SEINES Tags - in dem Tag, der diesen Zweig enthaelt, gibt es lanes/critique.md nicht mehr, also kein Doppelangebot. Ein Binary 0.3.0, das die drei Lanes nicht in sich traegt, bekommt sie aus dem Katalog voellig zurecht. Das Doppelangebot ist nur auf einer dev-Sicherung sichtbar (release.Current == 'dev' faehrt ohne ref gegen HEAD von master) und nur bis dieser Zweig gemerged ist - das Fenster ist 'wer an jaira selbst arbeitet, vor dem Merge'.

Der Preis des Filters waere eine neue Abhaengigkeit von core/market auf lane.Builtins(), ein Test und eine NOTES-Zeile - und er nimmt einen echten Fall weg: absichtlich die Katalog-Fassung einer Lane statt der eingebetteten holen. Wachstum fuer ein Problem, das der Tag von selbst schliesst.

Die Zeile in core/release/NOTES.md, die sagt 'jaira lanes market' biete sie nicht mehr an, bleibt richtig: ihr Leser ist der Besitzer eines Release-Binaries, und der hat den Pin.
