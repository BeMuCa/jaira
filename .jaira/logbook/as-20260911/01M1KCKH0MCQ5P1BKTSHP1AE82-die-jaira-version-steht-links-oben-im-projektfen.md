---
id: 01M1KCKH0MCQ5P1BKTSHP1AE82
title: Die jaira-Version steht links oben im Projektfenster
status: done
ready: true
creator: BeMuCa
goal: "Das TUI zeigt die laufende Binary-Version sichtbar im Board, links oben in der Ecke des Projektfensters"
context: "Berk am 03.09.: beim Arbeiten mit mehreren Binary-Staenden (self upgrade, go build, ~/.local/bin) ist nicht sichtbar, welche Version laeuft. Wunsch: Version ganz links oben in der Ecke des Projektfensters. ACHTUNG Kollision: DNAEPN (done, 03.09.) hat gerade entschieden, dass die bestehende Versionszeile (internal/tui/updatecheck.go, gezeichnet in home.go und view.go/Fusszeile) auf dev-Builds SCHWEIGT, weil 'jaira dev' nichts beantwortet; auf Releases zeigt sie 'jaira 0.1.1 - up to date'. Vor dem Bauen klaeren: Platzierung links oben zusaetzlich oder statt der Fusszeile, und ob ein dev-Build 'dev' zeigen soll (widersprueche DNAEPN) oder weiter nichts."
definition-of-done: "Das Board zeigt links oben die Version; ein Dev-Build zeigt 'dev'; ein Test deckt die Platzierung ab"
tags: []
blocked-by: []
commits:
  - 6c5d81b28e18049c24bafd21a1d639c263474a32
  - fac61b4e92969a4b28a2584efb4ffd6504633e0c
  - 8277d9585c33844fe8afa061e44e59fe99086486
  - 596f1a58dcdfb536e67774ab1c2485b0a2a4b730
  - ec2c861207e60b6c93894a637ef40b9ed60ee27e
  - d952f0e3d27917fa6f5037b01d0aa158135245ae
  - 6bd21c73399d119c3edd146759e47737c2828ef0
created-at: 2026-09-03T10:21:34Z
updated-at: 2026-09-11T14:51:51Z
assignee: BeMuCa
updated-by: Alexander Sacharov
claimed-by: EE-3NX6GL3-4099823
claimed-at: 2026-09-08T23:46:46Z
outcome-what: "Runde 2 nach Critique: versionHead() aufgeloest (renderBoard schreibt truncate(m.versionLine, m.width) direkt), die drei toten Leer-Wachen in view.go und home.go entfernt, h.render() im Home-Test einmal statt dreimal aufgerufen."
outcome-why: "versionLine() liefert seit Runde 1 auch auf dev-Builds eine Zeile, damit waren die Wachen und der Wrapper Reste der entfernten Verhaltensweise."
outcome-resolves: "Gleiche Ausgabe, drei Verzweigungen und eine Methode weniger; go test ./... -race gruen."
review-summary: |-
  internal/tui/view.go:771 versionHead() hat einen Aufrufer und tut nach dem Wegfall des dev-Leerfalls nur noch truncate; direkt in renderBoard einsetzen, der erklaerende Kommentar steht dort schon
  internal/tui/view.go:183 und internal/tui/home.go:342 pruefen auf versionLine != "" - dieser Zustand kann nicht mehr eintreten, weil versionLine() jetzt auch auf einem dev-Build eine Zeile liefert (updatecheck.go:59, alle vier returns nicht leer). Das ist ein Rest der Verhaltensweise, die dieses Ticket gerade entfernt hat; Bedingung raus, Zeile immer schreiben
  internal/tui/home.go:342 'if v := h.versionLine; v != ""' passt nicht zum Stil der Datei, die sonst 'if h.versionLine != ""' schreibt (alte Fusszeilenstelle) - mit dem Wegfall der Bedingung erledigt
  internal/tui/updatecheck_test.go:99-105 ruft h.render() dreimal auf; einmal in eine Variable
  Runde 2: none - die drei Findings aus Runde 1 sind umgesetzt, ein Durchgang ueber den Diff findet nichts, was eine Datei und eine konkrete Alternative benennen koennte.
review-gaps: |-
  Entfernt: dieselbe Begruendung stand dreifach - der Satz 'eigene Zeile, damit eine zweite darunter passt' in internal/tui/updatecheck.go (Doc-Kommentar) UND internal/tui/view.go:177 (Aufrufstelle), und der Fusszeilen-Umzug zusaetzlich in internal/tui/home.go:339. Jetzt steht die Layout-Begruendung nur an der Aufrufstelle, die Funktions-Begruendung nur am Doc-Kommentar, und home.go zeigt auf renderBoard.
  Nicht entfernt und warum: truncate(X.versionLine, X.width) steht in view.go:181 und home.go:341 zweimal - eine gemeinsame Hilfsfunktion braeuchte ein Interface oder eine freie Funktion mit zwei Parametern fuer je einen Aufrufer, das ist teurer als die Wiederholung. headLines = strings.Count(head,'\\n')+1 (view.go:206) sieht heute nach Overkill aus, weil head genau eine Zeile ist - es ist die vom Ticket verlangte Luecke fuer T8A8KM und benutzt dieselbe Redewendung wie sbLines zwei Zeilen darunter.
  Nichts verwaist: kein Import, kein Helfer und kein Zweig ohne Aufrufer; styMeta, truncate und centre haben alle weitere Nutzer. Kein Verhaltensunterschied in diesem Durchgang - nur Kommentare.
test-verdict: "pass: go test ./... -race -count=1 gruen auf dem wiederhergestellten Baum, RC=0, alle 16 Pakete ok (internal/tui 269.467s, internal/cli 78.027s); go build -o /tmp/jaira-p1ae82 ./cmd/jaira ok, gofmt -l internal/tui leer. DoD 1 dreiteilig einzeln geprueft: Board zeigt links oben die Version (am laufenden Binary im pty: erste Zeile 'jaira dev', darunter die Kopfzeile mit '21 tickets', darunter die Board-Leiste), dev-Build zeigt 'dev' (jaira --version sagt 'version dev', dieselbe Aufnahme), und der Platzierungstest deckt sie wirklich ab - mutiert geprueft, 'head := \"\"' macht TestBoardVersionSitsInTheTopLeftCorner, TestBoardHeadNamesADevBuild und TestBoardHeadCarriesTheVersionIndicator rot. Fusszeile in beiden Screens ohne Version, 'jaira dev' kommt je Aufnahme genau einmal vor."
question: |-
  Zwei Dinge zum Ansehen, beide nur mit Augen entscheidbar.
  1. Der Launcher (die Projektliste, 'jaira' ohne Argument) hatte die Versionszeile auch in seiner Fusszeile. Ich habe sie dort ebenfalls nach links oben gezogen, ueber das Wordmark - sonst haette der Launcher die Version ganz verloren, denn 'in der Fusszeile steht keine Version mehr' trifft auch ihn. Soll sie dort oben bleiben, oder soll der Launcher gar keine Version zeigen? Letzteres ist eine geloeschte Zeile in internal/tui/home.go:341, sonst nichts.
  2. Die Zeile im Board sieht so aus: 'jaira dev' linksbuendig in Zeile 1, darunter die Kopfzeile mit '21 tickets' rechts, darunter die Board-Leiste. Sie kostet dem Board eine Zeile Hoehe. Auf einem Release steht dort die lange Form 'jaira 0.1.1 - 0.1.2 available - run: jaira self upgrade'; die wird kurz, sobald T8A8KM die zweite Haelfte in die Pille darunter zieht. Passt das so, oder soll die Zeile noch etwas anderes sagen?
  Gepruefte Fakten dazu stehen in den Notes; go test ./... -race ist gruen.
---

# Die jaira-Version steht links oben im Projektfenster

## Definition of Done

- [x] Das Board zeigt links oben die Version; ein Dev-Build zeigt 'dev'; ein Test deckt die Platzierung ab
  proof: internal/tui/view.go:183 (eigene erste Zeile) + :771 versionHead; internal/tui/updatecheck.go:59 dev zeigt 'jaira dev'; internal/tui/home.go:342; Tests TestBoardVersionSitsInTheTopLeftCorner, TestBoardHeadNamesADevBuild, TestBoardHeadCarriesTheVersionIndicator, TestHomeHeadNamesADevBuild; go test ./... -race gruen

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] versionLine(): dev-Zweig gibt 'jaira dev' statt "" zurueck, weiter vor PollCache
- [x] Doc-Kommentar von versionLine() nachziehen: Fusszeile -> links oben, dev-Absatz umschreiben, DNAEPN-Teilumkehr benennen
- [x] Board: eigene erste Zeile links oben in renderBoard, Kopfblock-Hoehe in bodyHeight einrechnen
- [x] Board: versionLine aus statusBar entfernen
- [x] Launcher: versionLine aus der Fusszeile in Home.render entfernen und links oben setzen
- [x] Tests nachziehen: die vier Fusszeilen-Zusicherungen in devfooter_test.go und updatecheck_test.go
- [x] Neuer Test deckt die Platzierung ab: Version in der ersten Zeile links, Fusszeile ohne Version, dev zeigt 'dev'
- [x] go build -o /tmp/jaira-p1ae82 ./cmd/jaira und go test ./... -race gruen

## Progress
- **2026-09-08 23:45 · BeMuCa** — Berk hat am 08.09. die beiden offenen Fragen aus dem Kontext entschieden.
1. Platzierung: ERSETZEN, nicht zusaetzlich. Die Versionszeile verschwindet aus der Fusszeile und steht nur noch links oben im Projektkopf. Damit loest sich die notierte Kollision mit DNAEPN auf - es gibt danach nur noch eine Stelle, die eine Version behauptet.
2. Dev-Build: zeigt 'dev'. Das steht so schon im DoD und geht bewusst gegen DNAEPN, das die Fusszeile auf dev-Builds schweigen liess. Begruendung: links oben ist die Zeile eine Identitaetsangabe ('welche Binary laeuft hier'), nicht mehr eine Upgrade-Empfehlung - und genau die Frage hatte Berk beim Wechsel zwischen self upgrade, go build und ~/.local/bin.
Folgeticket T8A8KM haengt daran und setzt spaeter eine Pille '^ <version>' unter diese Zeile. Die Platzierung hier muss also Platz fuer eine zweite Zeile darunter lassen.
- **2026-09-09 07:36 · BeMuCa** — pre-process, was ich verifiziert und entschieden habe.

Zwei Zeichenstellen, nicht eine. versionLine() (internal/tui/updatecheck.go:35) hat heute zwei Aufrufer: Model.versionLine, gezeichnet in statusBar (internal/tui/view.go:829) = Board-Fusszeile, und Home.versionLine, gezeichnet in Home.render (internal/tui/home.go:394) = Launcher-Fusszeile. 'Projektkopf' meint nur das Board. Die Regel 'in der Fusszeile steht keine Version mehr' trifft aber beide Fusszeilen. Entscheidung: die Zeile wandert in BEIDEN Screens von unten nach links oben, damit der Launcher die Version nicht verliert und beide Screens dieselbe Stelle benutzen. Wenn Berk den Launcher lieber ganz ohne Version will, ist das eine Zeile weniger - bewusst so gebaut, nicht uebersehen.

Platzierung: eigene Zeile GANZ oben, nicht links in die Kopfzeile mit '13 tickets'. Grund ist T8A8KM: die Pille steht 'unter der Versionszeile'. Saesse die Version links in der Kopfzeile, waere die Zeile darunter schon die Board-Leiste ('1 requirementsgenie | 2 jAIra') - kein Platz. Als eigene erste Zeile kann T8A8KM eine zweite Zeile einfach anhaengen. Deshalb zaehle ich die Hoehe des Kopfblocks (strings.Count+1) statt eine 1 zu addieren: eine zweite Zeile rechnet sich dann von selbst mit, ohne dass T8A8KM die Hoehenrechnung anfassen muss.

Der dev-Fall bricht DNAEPN nur halb. DNAEPN hatte zwei Gruende: (a) 'jaira dev' sagt nichts Neues, (b) release.Current='dev' mit einer veroeffentlichten Version verglichen ergibt immer 'unterschiedlich', die Zeile wuerde also ein Upgrade auf AELTEREN Code bewerben und auf ein Kommando zeigen, das mit dev_build ablehnt. P1AE82 kippt nur (a). (b) bleibt: der dev-Zweig kehrt weiter VOR selfupdate.PollCache um und gibt nur 'jaira dev' zurueck - kein 'up to date', keine Pille, und kein abgekoppelter Check-Kindprozess fuer Contributor.

Keine Golden-/Snapshot-Tests in internal/tui - kein testdata-Verzeichnis, kein Treffer auf 'golden'. Was bricht, sind vier gezielte Zusicherungen in internal/tui/devfooter_test.go und internal/tui/updatecheck_test.go, die die Version genau in der Fusszeile erwarten.
- **2026-09-09 07:47 · BeMuCa** — in-progress, umgesetzt und am gerenderten Bild geprueft.

So sieht der Board-Kopf jetzt aus (stripANSI-Dump aus einem Wegwerf-Test, 120x24):
  jaira 0.1.1 - 0.1.2 available - run: jaira self upgrade
  001                                                        6 tickets
  +--- Backlog ---+ +--- Brainstorm ---+ ...
Fusszeile: '! 1 enter open . v compact . ... . q quit' - keine Version mehr.
Der Launcher hat die Zeile ueber dem Wordmark, ebenfalls links oben, Fusszeile dort ebenfalls ohne Version.

Was ich pruefen musste, statt es anzunehmen: internal/tui hat KEINE Golden-Tests (kein testdata/, kein 'golden'-Treffer). Was brach, waren vier handgeschriebene Zusicherungen, die die Version genau in der Fusszeile erwarteten. Statt sie stumpf umzustellen habe ich sie umgedreht: jeder Test behauptet jetzt beides - oben steht sie, unten nicht - damit ein Rueckfall nach unten auffliegt und nicht nur ein Verschwinden.

devfooter_test.go heisst jetzt devbuild_test.go (git mv). Der alte Name behauptete eine Fusszeile, die es fuer diese Zeile nicht mehr gibt.

Die Hoehenrechnung in renderBoard zaehlt den Kopfblock (strings.Count(head,'\\n')+1) statt eine 1 abzuziehen. Absicht: T8A8KM haengt eine zweite Zeile an, und dann stimmt bodyHeight von selbst. Wer dort '- 1' hinschreibt, schiebt die Spalten unter den Bildschirmrand.

Sackgasse, die ich verworfen habe: die Version links in die Kopfzeile selbst (dort, wo rechts '6 tickets' steht). Kostet keine Zeile und ist auch 'links oben' - aber die Zeile DARUNTER ist dann schon die Board-Leiste, und T8A8KM will die Pille darunter. Haette T8A8KM gezwungen, erst eine Zeile einzuschieben.

Ausserhalb des Scope aufgefallen, NICHT angefasst: internal/cli/update.go nudgeIfStale haelt die CLI bewusst still ueber Releases - unberuehrt, die Aenderung betrifft nur das TUI.
- **2026-09-09 07:54 · BeMuCa** — Beim Commit fast Fremdgut mitgenommen: 'git add -A .jaira/tickets' hat fuenf ungetrackte Tickets von Berk (81XRXX, VS5DFW, DBJTKQ, T8A8KM, ZP3R48) in meinen Commit gezogen. Rueckgaengig mit 'git reset --soft HEAD~1' und 'git restore --staged <datei>' - sie liegen wieder ungetrackt da, wie vorher. Merke fuer die naechste Session: beim Ticket-Commit die EINE Ticketdatei mit Pfad adden, nie '-A' auf .jaira/tickets.
- **2026-09-09 07:55 · BeMuCa** — critique, Runde 1. Drei Findings, alle vom gleichen Ursprung: der dev-Zweig gibt jetzt IMMER eine Zeile zurueck, also sind die drei '!= ""'-Wachen aus der alten Welt tot. Damit schrumpft versionHead() auf ein truncate und gehoert in renderBoard hinein.

Was ich geprueft und ausdruecklich STEHEN lasse, damit die naechste Runde es nicht wieder aufmacht:
- Model und Home werden nur in New() (model.go:230) bzw. NewHome() gebaut, nirgends als Literal. Es gibt also keinen Nullwert-Pfad, der ein leeres versionLine in renderBoard traegt - deshalb ist das Entfernen der Wachen sicher und nicht bloss huebscher.
- headLines = strings.Count(head,'\\n')+1 bleibt, obwohl head heute genau eine Zeile ist. Das ist keine Spekulation, sondern die vom Ticket verlangte Luecke fuer T8A8KMs Pille.
- Die Zeile als eigene Zeile statt links in die Kopfzeile: begruendet, nicht wieder aufmachen (siehe Note vom in-progress-Schritt).
- view.go zerlegt seinen Kopf ohnehin in Methoden (header, boardProjectLine, renderSessions). versionHead() haette also zum Muster gepasst - der Grund, sie trotzdem aufzuloesen, ist allein, dass nach Finding 2 nichts mehr drin steht.
- **2026-09-09 08:00 · BeMuCa** — in-progress, Runde 2: die drei Critique-Findings sind umgesetzt. versionHead() ist weg, renderBoard schreibt die Zeile direkt, die drei toten '!= ""'-Wachen sind raus, der Test haelt h.render() in einer Variablen. Kein Verhaltensunterschied - dieselbe Ausgabe, drei Verzweigungen weniger.
- **2026-09-09 08:01 · BeMuCa** — critique, Runde 2: keine Findings mehr, Lane fertig. Geprueft und bewusst nicht als Finding gefuehrt:
- renderBoard schreibt die Zeile jetzt bedingungslos, also erscheint sie auch in modeMove/modeFilter/modeCreate, wo statusBar frueher frueh zurueckkehrte und die Version dort verschwand. Das ist konsistenter, nicht schlechter: eine Identitaetsangabe soll nicht davon abhaengen, ob gerade ein Prompt offen ist.
- Andere Screens (renderDetail, renderHelp, renderPipeline...) zeigen die Version weiterhin nicht. Das war vorher genauso - statusBar haengt nur an renderBoard -, also keine Regression.
- internal/tui/fitwindow_test.go:120 misst die Board-Hoehe hinter m.View(), und View() endet in clampBlock, das jede Zeile hart abschneidet. Der Test kann ein Ueberlaufen von bodyHeight also gar nicht sehen. Das ist bestehender Zustand, nicht von diesem Ticket gemacht - aber der Grund, warum ich headLines lieber messe als annehme. Fuer die testing-Lane: ein dev-Build auf einem kurzen Terminal ist der Fall, in dem der Kopf eine Zeile mehr kostet als vorher (unten fiel dort nichts weg, weil die Fusszeile auf dev schwieg).
- **2026-09-09 08:06 · BeMuCa** — optimize: nur Kommentar-Deduplikation, kein Code. Der einzige echte Fund war dreifach erzaehlte Begruendung. Vor dem Aufraeumen gesucht: kein bestehendes Muster fuer eine Kopfzeile-links-oben im Repo, kein zweiter Ort, der release.Current ins TUI zeichnet (grep versionLine, release.Current) - internal/cli/update.go nudgeIfStale ist die CLI-Seite und bleibt bewusst stumm.
- **2026-09-09 08:08 · BeMuCa** — testing, Pass 3 (Funktion) am laufenden Binary, nicht nur im Test:
Board mit einem pty aufgenommen - '(sleep 4; printf q) | script -qec "stty rows 30 cols 140; /tmp/jaira-p1ae82 board" /dev/null' - und die Escape-Sequenzen weggeschnitten. Erste Zeile: 'jaira dev', danach die Kopfzeile mit '21 tickets' rechts, danach die Board-Leiste '1 jAIra | 2 requirementsgenie | 3 LHChecker | 4 localhtml', danach die Spalten. Letzte Zeile: 'enter open . v compact . ... . q quit' - keine Version. 'jaira dev' kommt in der ganzen Aufnahme genau einmal vor.
Launcher genauso aufgenommen (bare '/tmp/jaira-p1ae82' aus einem leeren Verzeichnis): 'jaira dev' in der ersten Zeile ueber dem Wordmark, Fusszeile endet mit 'q quit', ebenfalls genau ein Vorkommen.
/tmp/jaira-p1ae82 --version sagt 'jaira version dev' - es ist wirklich der dev-Fall, den DNAEPN stumm gestellt hatte.
Wie das reproduzierbar ist, falls es nochmal gebraucht wird: der pty-Umweg ist noetig, weil bubbletea ohne Terminal nichts zeichnet; ohne 'stty rows/cols' im script-Kommando nimmt es 80x24.
- **2026-09-09 08:09 · BeMuCa** — testing, Pass 2 (die Forderung): den Platzierungs-Test nicht nur laufen lassen, sondern mutiert. 'head := truncate(...)' plus WriteString in renderBoard durch 'head := ""' ersetzt und die vier Tests laufen lassen - alle rot, wortwoertlich:
  TestBoardVersionSitsInTheTopLeftCorner: board first line = "001 ... 6 tickets", want it to start flush left with the version
  TestBoardHeadNamesADevBuild: board first line = "001 ... 6 tickets", want it to name the dev build
  TestBoardHeadCarriesTheVersionIndicator: board first line = "001 ... 6 tickets", want the version indicator naming the available release
Danach view.go aus der Kopie zurueckgeholt, dieselben Tests gruen. Der Test deckt die Platzierung also wirklich ab und nicht nur die Anwesenheit des Strings irgendwo.
Wichtig fuer die naechste Session: die Mutation lief, waehrend im Hintergrund schon ein 'go test ./... -race -count=1' unterwegs war. Dessen Ergebnis war damit wertlos - go test kompiliert das Paket beim Erreichen des Pakets, nicht am Anfang. Lauf abgebrochen und nach dem Zurueckholen neu gestartet.
- **2026-09-09 08:19 · BeMuCa** — Nachtrag zur Verifikation: bestehende gofmt-Drift in internal/cli/tickets.go. Nicht von mir - 'git show 823fc81:internal/cli/tickets.go | gofmt -l' meldet sie schon vor meinen Commits, und die Datei steht nicht in 'git diff --name-only 823fc81..HEAD'. Liegen gelassen, wie es die Regel fuer fremden toten/schiefen Code verlangt. Wer sie anfasst, sollte es in einem eigenen Commit tun, sonst rauscht ein reiner Formatlauf in einen Feature-Diff.
- **2026-09-11 14:51 · Alexander Sacharov** — Abgenommen von Alexander. Bestaetigt an zwei Stellen: internal/tui/home.go:81 setzt die Versionszeile, und auf den heute aufgenommenen Vorfuehrungen steht 'jaira dev' links oben im Board - ein Dev-Build, der sich als solcher nennt. Die Aenderung ist ausserdem seit 0.1.2 veroeffentlicht.
