---
id: 01M4THDAZW15SJ4Q1R959WE1QK
title: q fuehrt auf jeder Seite der TUI genau eine Seite zurueck
status: critique
ready: true
creator: BeMuCa
goal: "Auf jeder Seite der TUI fuehrt q eine Seite zurueck - dorthin, wohin auch esc fuehrt -, und nur auf der obersten Seite beendet q jaira."
context: |-
  Heute ist q nicht einheitlich (Karte aller Seiten vom 2026-10-13, Pfade unter internal/tui):
  - Kompaktansicht (v): q beendet jaira, statt zurueck zum Board zu gehen (model.go:1208); esc tut dort gar nichts. Wer aus der Lane-Ansicht (enter) zweimal q drueckt, ist ungewollt raus.
  - Lane-Katalog im Lane-Editor: nur esc schliesst, q wird ignoriert (lanes.go:205).
  - Meldungen aus Einstellungen, Lanes, Standard-Board oder Projekte schliessen aufs Board statt zurueck auf die Seite, die sie gezeigt hat (notify, model.go:1727-1733).
  - Ueberall sonst gehen q und esc schon gleich eine Seite zurueck: Detail, Hilfe, Legende, Pickers u und M, Links, Projekte, Einstellungen, Lane-Editor, Standard-Board, Milestone Roadmap mit allen Dialogen.
  - In Textfeldern ist q ein Buchstabe und kann nicht zurueck heissen: Filter /, neues Ticket n, Loeschen, Feld-Editor e, Formular der Milestone Roadmap. Dort bleibt esc der Rueckweg.
  - Board: q beendet jaira (model.go:1489); kam das Board aus dem Launcher, landet man wieder im Launcher (cli/home.go). Im Launcher beenden q und esc.
  Berk am 2026-10-13: 'Man soll konsistent mit q eine Seite zurück kommen. Bei manchen muss man esc drücken.' Er will den Rueckweg-Fluss jeder Seite sehen, bevor gebaut wird.
definition-of-done: "Kompaktansicht: q und esc tun dasselbe wie q auf dem Board (zurueck zum Launcher bzw. beenden); Lane-Katalog schliesst mit q; Meldungen schliessen zurueck auf ihre Seite; ein Test prueft q=esc auf jeder Seite"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-13T11:47:20Z
updated-at: 2026-10-13T14:17:04Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-2027046
claimed-at: 2026-10-13T13:50:27Z
outcome-what: "Kompaktansicht: esc beendet wie q (model.go:1211) - zurueck zum Launcher bzw. raus bei 'jaira board'. Lane-Katalog: q schliesst wie esc (lanes.go:207). notify: Einstellungen, Lane-Editor, Standard-Board und Projektliste sind ihr eigener Rueckweg, und modal.go underlying zeichnet sie unter der Meldung. Projektliste: ein Board, das nicht aufgeht, zeigt seinen Fehler statt ihn zu verschlucken. Hilfe-Zeile fuer q (view.go:1672), README.md und docs/COMMANDS.md unter Keys, Release-Notiz core/release/notes/9WE1QK.md. Neuer Test internal/tui/qback_test.go (27 Board-Seiten, 5 Launcher-Seiten, Board als oberste Seite, Hilfe)."
outcome-why: "Berk 2026-10-13: 'Man soll konsistent mit q eine Seite zurueck kommen. Bei manchen muss man esc druecken.' Drei Stellen brachen die Regel: esc tat in der Kompaktansicht nichts, q schloss den Lane-Katalog nicht, und Meldungen auf vier Seiten warfen aufs Board statt zurueck."
outcome-resolves: "DoD 1: compact_view-Fall im Tabellentest (q und esc beenden beide, wie q auf dem Board). DoD 2: lane_catalogue-Fall. DoD 3: TestMessageDismissesToThePageItWasRaisedOn mit settings/lane editor/default board/boards plus TestABoardThatWillNotOpenSaysSoOverTheList. DoD 4: TestQGoesWhereEscGoesOnEveryPage + TestQGoesWhereEscGoesOnTheLauncher + TestHelpNamesQAsOnePageBack; Textfelder bewusst ausgenommen. DoD 5: go vet + go test ./... -race gruen, Release-Notiz geschrieben. Gegenprobe: ohne die Fixes scheitern genau compact_view, lane_catalogue und die vier Meldungsfaelle."
executed-by: claude-opus-5-5
---

# q fuehrt auf jeder Seite der TUI genau eine Seite zurueck

## Definition of Done

- [x] Kompaktansicht: q und esc tun dasselbe wie q auf dem Board - aus dem Launcher geoeffnet geht es zurueck zum Launcher, und erst dort beenden q/esc jaira
  proof: internal/tui/model.go:1211 esc neben q/ctrl+c in modePipeline -> tea.Quit wie das Board (model.go:1495), runHome internal/cli/home.go:22 zeigt danach den Launcher; TestQGoesWhereEscGoesOnEveryPage/compact_view
- [x] Lane-Katalog: q schliesst ihn wie esc
  proof: internal/tui/lanes.go:207 case "esc", "q" im catalogueOpen-Zweig; TestQGoesWhereEscGoesOnEveryPage/lane_catalogue
- [x] Eine Meldung schliesst zurueck auf die Seite, die sie gezeigt hat (Einstellungen, Lanes, Standard-Board, Projekte)
  proof: internal/tui/model.go notify: modeSettings, modeLanes, modeDefaultBoard, modeProjects merken sich als returnTo; modal.go underlying zeichnet sie unter der Meldung; model.go modeProjects enter laesst eine Fehlermeldung stehen; TestMessageDismissesToThePageItWasRaisedOn, TestABoardThatWillNotOpenSaysSoOverTheList
- [x] Ein Test geht jede Seite durch und prueft, dass q dorthin fuehrt, wohin esc fuehrt (Textfelder ausgenommen); die Hilfe nennt q als 'eine Seite zurueck'
  proof: internal/tui/qback_test.go TestQGoesWhereEscGoesOnEveryPage (27 Seiten), TestQOnTheBoardLeavesIt, TestQGoesWhereEscGoesOnTheLauncher (5 Seiten), TestHelpNamesQAsOnePageBack; Hilfe internal/tui/view.go:1672
- [x] go test ./... -race gruen; Release-Notiz unter core/release/notes/
  proof: go vet ./... && go test ./... -race -count=1: alle Pakete ok (internal/tui 290s); core/release/notes/9WE1QK.md

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Kompaktansicht: esc neben q und ctrl+c in modePipeline (internal/tui/model.go:1208) - beendet das Programm genau wie q auf dem Board (model.go:1489); runHome (internal/cli/home.go:22) zeigt danach wieder den Launcher
- [x] Lane-Katalog: q neben esc im catalogueOpen-Zweig (internal/tui/lanes.go:205)
- [x] notify (model.go:1727): Einstellungen, Lane-Editor, Standard-Board und Projekte merken sich selbst als Rueckweg; Kommentar anpassen; TestMessageDismissesToThePageItWasRaisedOn (returnto_test.go) erwartet bei settings jetzt modeSettings, dazu lanes, default board, projects
- [x] Neuer Test internal/tui/qback_test.go: jede Seite ohne Textfeld ueber echte Tasten oeffnen, auf zwei gleich aufgebauten Modellen einmal q und einmal esc - beide landen auf derselben Seite und beenden gleich, und esc verlaesst die Seite wirklich; Board: q beendet; Launcher und seine Unterseiten ebenso
- [x] Hilfe (?, view.go renderHelp): Zeile 'q eine Seite zurueck' und Zeile zur Kompaktansicht; README.md und docs/COMMANDS.md im Abschnitt Keys
- [x] Release-Notiz core/release/notes/9WE1QK.md, je Aenderung eine Zeile '- ' als Anweisung
- [x] go vet ./... und go test ./... -race gruen, dann Commit mit der Ticketdatei

## Progress
- **2026-10-13 13:28 · BeMuCa** — Berk 2026-10-13 zur q-Karte: (1) die roten Stellen auf 'eine Seite zurueck' ziehen - ja, aber die Kompaktansicht landet wie das Board beim Launcher, nicht beim Board; erst im Launcher beendet q/esc. (2) Board bleibt oberste Seite: q zurueck zum Launcher bzw. beenden, wenn ohne Launcher gestartet - 'sieht gut aus'. (3) Launcher: q und esc beenden beide - bleibt.
- **2026-10-13 13:56 · BeMuCa** — pre-process, Befund im Code (Stand bc05c9f): q in der Kompaktansicht tut heute schon genau das, was q auf dem Board tut - beide rufen m.Close() und tea.Quit (model.go:1208 und model.go:1489). Aus dem Launcher gestartet faengt runHome (internal/cli/home.go:22) das Programmende ab und zeigt den Launcher wieder. Fuer DoD 1 aendert sich also nur esc: es tat dort nichts (pipelineKey kennt kein esc) und tut jetzt dasselbe wie q. Das Board selbst bleibt wie es ist: esc leert dort den Filter, q geht zurueck/beendet.
- **2026-10-13 13:56 · BeMuCa** — Entscheidung (Claude, aus Berks Antwort 2): ohne Launcher gestartet (jaira board) beenden q und esc in der Kompaktansicht jaira - dieselbe Regel wie auf dem Board. Grund: die Kompaktansicht ist eine zweite Ansicht derselben obersten Seite (v schaltet hin und her), keine Seite darunter; Berk hat fuer das Board 'q zurueck zum Launcher bzw. beenden' bestaetigt. Das TUI-Modell weiss nicht, ob ein Launcher davor lief (tui.New hat kein solches Feld) - und muss es auch nicht, weil runHome das Programmende so oder so richtig deutet.
- **2026-10-13 13:56 · BeMuCa** — Entscheidung: bei den Meldungen nur die vier genannten Seiten (Einstellungen, Lane-Editor, Standard-Board, Projekte) auf 'zurueck zu sich selbst' stellen, nicht pauschal jede Seite. Grund: Legende und die Picker u/M merken sich ihren eigenen Rueckweg in m.returnTo (model.go:1525, ticks.go openTicks) - eine Meldung darueber wuerde ihn ueberschreiben und die Seite haette danach keinen Rueckweg mehr. Die vier Seiten halten ihren Zustand in eigenen Feldern (settingsScreen, laneScreen, board, projects) und lesen returnTo nie. Ausloeser im Alltag: ein Fehler beim Abholen der Refs (refFetchedMsg, model.go:1005) oder beim Nachladen nach 'n' im Lane-Editor (newLaneDoneMsg) kommt mitten auf diesen Seiten an.
- **2026-10-13 13:56 · BeMuCa** — Entscheidung: das frontmatter-Feld definition-of-done stand noch auf der alten Fassung 'Kompaktansicht: q und esc gehen zurueck zum Board' - das widerspricht Berks Antwort (1). Auf den Stand der Checkliste gezogen; massgeblich fuer das Gate ist die Checkliste. Textfelder bleiben ausgenommen: Filter /, neues Ticket n, Loeschen X, Feld-Editor e, Folgeticket n im Detail, Formular der Milestone Roadmap, und die Zahl 'logbook days' in den Einstellungen (dort tippt man nur Ziffern, esc bricht ab).
- **2026-10-13 13:57 · BeMuCa** — in-progress: dazu gefunden - Enter auf einem Board in der Projektliste (p) rief switchBoard und setzte danach stur modeBoard (model.go ~1351). Scheiterte das Oeffnen (ticket.Discover-Fehler), hatte notify schon modeMessage gesetzt, und die Meldung wurde sofort ueberschrieben: der Fehler war nie zu sehen. Das ist die einzige Meldung, die die Projektliste selbst ausloest - ohne den Fix waere DoD 3 fuer 'Projekte' nur fuer Hintergrund-Meldungen wahr. Jetzt: modeBoard nur, wenn keine Meldung offen ist; die Meldung schliesst zurueck auf die Liste.
- **2026-10-13 14:00 · BeMuCa** — in-progress: Test-Aufbau internal/tui/qback_test.go - jede Seite wird auf ZWEI frischen Modellen gleich geoeffnet (einmal q, einmal esc), weil sich ein Model nicht klonen laesst. Verglichen wird page(m) (Modus plus Unterseite: Katalog/Rueckfrage im Lane-Editor, Stufe der Roadmap, detailFrom), nicht der gerenderte Bildschirm - zwei Fixture-Boards liegen in verschiedenen Temp-Ordnern, ein Render-Vergleich waere daran gescheitert. Jeder Fall nennt die Seite, die er erreichen muss (at), sonst haette z.B. ein Picker, der statt sich zu oeffnen eine Meldung zeigt, still bestanden. Gegenprobe: ohne die Fixes in model.go/lanes.go schlagen genau 'compact view' und 'lane catalogue' fehl.
- **2026-10-13 14:05 · BeMuCa** — in-progress: die Hilfe bekommt EINE Zeile fuer q (view.go renderHelp, Abschnitt 'Move around'), nicht drei. Grund: TestHelpExplainsEveryCardMark (150x60) und TestHelpNamesZ (40x40) pruefen, ob Eintraege im sichtbaren Fenster stehen - schon drei Zeilen schoben '@name' und 'z' unter den Rand. Bei 40 Spalten bricht auch eine Zeile mehrfach um, darum TestHelpNamesZ von Hoehe 40 auf 50 (44 war das Minimum; die Absicht des Tests ist die Breite, nicht die Hoehe). Die ausfuehrliche Regel (Board ist oberste Seite, Kompaktansicht wie Board, Textfelder nur esc) steht in README.md und docs/COMMANDS.md unter Keys.
- **2026-10-13 14:06 · BeMuCa** — in-progress: Folge der notify-Aenderung - eine Meldung wird ueber underlying(returnTo) gezeichnet (modal.go:82), und das kannte Einstellungen, Lane-Editor, Standard-Board und Projektliste nicht: die Meldung stand ueber dem Board, und erst beim Schliessen kam die Einstellungsseite zum Vorschein. underlying zeichnet diese vier Seiten jetzt selbst, so dass Schliessen nur den Kasten wegnimmt (die Regel, die der Kommentar dort schon nennt). TestMessageDismissesToThePageItWasRaisedOn vergleicht dafuer die oberste Bildschirmzeile (sie liegt ueber dem Kasten) vor und nach dem Schliessen.
