---
id: 01M2SNEP2NN02DN5W0ENVSC1GW
title: "Die Pruefschleife gehoert ins Binary, nicht in den Katalog"
status: in-progress
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
updated-at: 2026-09-18T12:02:19Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-56920
claimed-at: 2026-09-18T10:31:58Z
mode: conversational
outcome-what: "Plan in elf Schritten, am Code gemessen"
outcome-why: "Form 2 haengt an einer einzigen Aufteilung: Lane.Builtin traegt zwei Tatsachen, und nur zwei der elf Fundstellen gehoeren auf die neue"
outcome-resolves: "Der Weg von der Entscheidung zur Umsetzung steht: welches Feld, welche Datei, welcher Test, in welcher Reihenfolge"
---

# Die Pruefschleife gehoert ins Binary, nicht in den Katalog

## Definition of Done

- [x] Das Ticket nennt EINE gewaehlte Form mit Begruendung, warum die anderen drei verworfen wurden - an der Scope-Regel aus CLAUDE.md gemessen, nicht am Bauchgefuehl.
  proof: Notiz 'ENTSCHIEDEN von Alex am 2026-09-18' auf diesem Ticket: Form 2 gewaehlt, Formen 1/3/4 und Weg B mit Begruendung verworfen
- [x] Ein frisch mit 'jaira init' angelegtes Board erfaehrt von der Pruefschleife, ohne dass jemand den Katalog kennt. Nachgestellt an einer leeren Testdoska: der Weg von 'jaira init' zu einer arbeitenden critique-Lane ist ohne Vorwissen gehbar.
  proof: internal/cli/tickets.go newLanesCmd — Fuss 'Shipped with this binary, not on this board'; Tests internal/cli/lanes_test.go TestLanesNamesTheShippedLanesThisBoardLacks, TestLanesJSONCarriesAvailable, TestLanesSaysNothingWhenTheBoardHasEverything
- [x] Wer jaira ohne Netz benutzt, kommt an die drei Lanes heran, oder die Fehlermeldung sagt, dass sie aus dem Netz kommen und wie man sie sonst bekommt. 'jaira lanes market' holt heute von GitHub.
  proof: internal/cli/lanes_test.go TestLanesAddInstallsTheReviewLoopWithoutNetwork — mit JAIRA_MARKET_API auf 127.0.0.1:1 installiert 'jaira lanes add critique|optimize|testing' alle drei und setzt sie zwischen in-progress und human
- [x] Ein Board, das jaira als reinen Tracker benutzt, wird nicht mit einem Konveyer beladen, den niemand faehrt - die Wahl bleibt eine Wahl.
  proof: core/lane/defaultboard_test.go TestFreshBoardGetsOnlyTheDefaultLanes — ein frisch angelegtes Board laedt als genau die zehn Lanes von heute, obwohl das Binary jetzt dreizehn traegt
- [ ] Je eine Zeile in core/release/NOTES.md fuer das, was ein Benutzer dadurch anders tut.
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
- [~] Je eine Zeile unter '## Unreleased' in core/release/NOTES.md fuer die drei Lanes im Binary, den neuen Fuss von 'jaira lanes' und den versionsgebundenen Katalog (DoD 5)
- [ ] go build ./... , go vet ./... und die volle Testsuite gruen

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
