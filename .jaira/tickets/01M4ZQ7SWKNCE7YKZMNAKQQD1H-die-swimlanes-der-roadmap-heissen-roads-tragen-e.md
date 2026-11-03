---
id: 01M4ZQ7SWKNCE7YKZMNAKQQD1H
title: "Die Swimlanes der Roadmap heissen Roads, tragen einen Namen und lassen sich im TUI umbenennen"
status: in-progress
ready: true
creator: BeMuCa
assignee: berkc
goal: "Jede Swimlane der Milestone-Roadmap ist eine Road - ein eigenes Objekt auf seinem Ref mit Titel, Beschreibung und Milestone-Liste -, die ohne Titel Road-1..Road-n heisst und im TUI per Pfeiltaste gewaehlt und mit Enter benannt wird; Milestones wissen nichts von ihrer Road, Unfertiges bleibt im Milestone"
context: |-
  Was heute fehlt: Die Milestone-Roadmap im TUI (Taste L, aus D0SAHM) zeigt Swimlanes ohne Namen.
  Eine Swimlane ist heute der Baum unter einem Milestone ohne Vorgaenger; gebaut aus after: in core/milestone/layout.go (erstes after: entscheidet die Lane).
  Berk am 15.10.: Swimlanes heissen 'Road', bekommen Namen, Default Road-1..Road-n, per Pfeiltaste waehlen, mit Enter umbenennen.
  Entscheidungen Maintainer 03.11.2026 (Notizen): Road ist ein eigenes Objekt (Titel, Beschreibung, Liste ihrer Milestones); Milestones wissen nichts von ihrer Road.
  Eine Road ist die ganze Swimlane samt Verzweigungen; der Aufbau darin bleibt aus after:.
  Ein Milestone in zwei Roads ist verboten (exit 3); nach einem Merge meldet validate es.
  Road liegt nur auf ihrem Ref, wie Epics (refs/jaira/epics/ -> refs/jaira/roads/).
  Epics bleiben; das Etikett 'Epic' im TUI ist nicht Teil dieses Tickets (VD4ZDW).
  Unfertiges bleibt im Milestone - bisher rieten Hilfe und Vorlage zum Verschieben.
definition-of-done: "Jede Swimlane der Roadmap zeigt einen Road-Namen: 'road:' ihres ersten Milestones, sonst Road-1..Road-n in Anzeigereihenfolge; Test"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-15T12:05:22Z
updated-at: 2026-11-03T23:09:21Z
updated-by: berkc
claimed-by: berk-ca-130584
claimed-at: 2026-11-03T17:20:46Z
question: ""
---

# Die Swimlanes der Roadmap heissen Roads, tragen einen Namen und lassen sich im TUI umbenennen

## Definition of Done

- [x] Eine Road ist ein eigenes Objekt nur auf ihrem Ref refs/jaira/roads/<name> (wie Epics): Titel, Beschreibung, Liste ihrer Milestone-IDs; Milestones tragen kein Road-Feld; paralleles Schreiben wird gemergt wie bei Epics; Test
  proof: core/road/road.go Merge, core/road/road_test.go TestCreateRefusesAMilestoneInAnotherRoad; core/gitref/road_test.go TestRoadArrivesWithoutASharedBranch; core/refsync/road_test.go TestAConcurrentRoadListingIsMergedNotLost; no Road field in core/milestone
- [ ] Roadmap: eine Road ist eine ganze Swimlane - ihre Milestones stehen in ihrer Lane, der Aufbau darin kommt weiter aus after:; ein Milestone ohne Road steht in der Lane seines ersten Vorgaengers; der Lane-Kopf zeigt den Road-Titel, sonst Road-1..Road-n in Anzeigereihenfolge; Test
- [ ] TUI-Roadmap: vom obersten Bubble waehlt Pfeil hoch den Road-Kopf, Enter oeffnet ein Namensfeld, Enter speichert (eine unbenannte Lane wird dabei zur Road mit allen Milestones, die in ihr stehen), Esc bricht ab; Test
- [ ] CLI 'jaira road create|ls|set|add|rm' ('ls --json' traegt Titel, Beschreibung, Milestones); ein Milestone in einer zweiten Road wird mit exit 3 verweigert; liegt er nach einem Merge doch in zweien, zeichnet ihn die aeltere Road und 'jaira validate' meldet es; Test
- [ ] Unfertiges bleibt im Milestone: Hilfe von 'jaira milestone' (internal/cli/milestones.go), Vorlage im neuen Milestone-File (core/milestone/milestone.go) und Paketdoku core/epic/epic.go sagen das statt 'in den naechsten verschieben'
- [ ] README, docs/COMMANDS.md und eine Zeile in core/release/notes/KQQD1H.md

## Options

- [x] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] gitref: RoadPrefix refs/jaira/roads/, Read/ReadAt/Write/List/SHA wie Epic
- [x] outbox: KindRoad 'roads' an jeder Stelle von KindEpic, Queue/Pending/DropRoad, Sender.WriteRoad
- [x] refsync/road.go wie epic.go (RoadNames, RoadContent, RecordRoad, rebaseRoad) + Haken in refsync.go; cli/refs.go refSubject/rejectedAdvice
- [x] core/road: Road (title, description, created-at, '- <msID>'), Merge 3-Wege, LoadAll aelteste zuerst, Create/Assign/Set unter Lock 'roads', TakenError, Owner/Duplicates/Groups/NewName; Tests
- [x] milestone.Layout(all, roads): Lane = Road die listet, sonst Lane des ersten Vorgaengers, sonst eigene; Vorgaenger in fremder Lane = Extra-Pfeil + Top; mehrere Tops je Lane; roads=nil unveraendert; Tests
- [ ] CLI jaira road create|add|rm|set|ls [--json], zweite Road exit 3; Tests
- [ ] validate: Warnung road_duplicate; Test
- [ ] TUI: Roads laden, Kopfzeile je Lane (Titel / Road-n), Ziel onRoad per Pfeil hoch, Enter-Namensfeld, Enter speichert (unbenannt: Road mit allen Lane-Milestones), Esc; Hilfe; Tests
- [ ] Ratschlag-Texte 'Unfertiges bleibt im Milestone' (milestones.go Hilfe, milestone.go Vorlage, epic.go Doku)
- [ ] README, docs/COMMANDS.md, core/release/notes/KQQD1H.md

## Progress
- **2026-10-31 01:25 · berkc** — Maintainer 31.10.2026 (Punkt 11): Wenn Roads Name + Beschreibung bekommen, koennen sie als Epics dienen - eine Road ist dann ein Feature, die Milestones darauf die Runden. Im Plan pruefen, ob Epic und Road ein Objekt werden.
- **2026-11-03 17:18 · berkc** — Maintainer 03.11.2026: (1) Statt 'road:' am ersten Milestone: die Road ist ein eigenes Objekt (Name, Beschreibung, geordnete Liste ihrer Milestones) - Eltern/Kind; Milestones wissen nicht, zu welcher Road sie gehoeren. (2) Epics bleiben. Pro Board waehlbar, worauf das Etikett 'Epic' im TUI zeigt: auf einen Milestone, auf eine Road, oder auf das eigene Epic-Objekt - jedes Team entscheidet selbst. (3) Was in einem Milestone nicht fertig wird, bleibt in diesem Milestone. Klaeren im Brainstorm: Verhaeltnis Road-Liste zu after: (ersetzt die Liste after: innerhalb einer Road?), ein Milestone in zwei Roads, und dass Epic seit PG9P9Y externe Schluessel/Links traegt.
- **2026-11-03 17:22 · berkc** — Brainstorm 03.11 (Befund im Code): (a) Eine Swimlane ist heute KEINE Kette, sondern der Baum unter einem Top: core/milestone/layout.go baut Lanes aus after:, das erste after: entscheidet die Lane, Geschwister (zwei Milestones mit gleichem Vorgaenger) sind erlaubt und werden nebeneinander gezeichnet (internal/tui/layer.go placement onRight). Eine geordnete Liste kann Geschwister nicht ausdruecken. (b) Milestones sind repo-weit (.jaira/milestones/<name>-<ID>.md + Ref je ID), nicht pro Board; Epics liegen nur auf Refs (refs/jaira/epics/<name>, core/epic/epic.go). (c) Das TUI zeigt Epics heute nirgends - ein 'Epic'-Etikett, das man umbiegen koennte, existiert im TUI nicht; Epics sieht man nur in CLI und Web (Epic-Editor aus PG9P9Y). (d) 'Unfertiges wandert in den naechsten Milestone' steht nur als Rat: internal/cli/milestones.go:52 (Hilfe), core/milestone/milestone.go:174 (Kommentar im neuen Milestone-File), core/epic/epic.go:7 (Paketdoku). Punkt (3) des Maintainers dreht das um; kein Code verschiebt Tickets. (e) Pfeile in der Roadmap wandern heute zwischen Bubbles und [+]-Zielen, Enter auf einer Bubble oeffnet ihr Detail - fuer die Road-Auswahl reicht ein neues Ziel 'Road-Kopf' ueber dem Top (hoch vom Top); das entscheidet der Plan.
- **2026-11-03 17:22 · berkc** — Brainstorm: selbst entschieden (Widerspruch erwuenscht): (1) Epic und Road bleiben zwei Objekte (Maintainer 03.11 (2)), Road bekommt Name + Beschreibung + geordnete Milestone-Liste, Milestones tragen kein road:-Feld. (2) Ohne Road-Objekt bekommt eine Swimlane den Anzeigenamen Road-1..Road-n in Anzeigereihenfolge; geschrieben wird erst beim Umbenennen. (3) CLI als eigene Gruppe 'jaira road create|ls|set|add|rm|rename' wie 'epic'/'milestone', statt 'milestone set --road' - der Milestone weiss nichts von seiner Road, also gehoert der Befehl der Road. (4) Die drei Ratschlag-Texte aus (d) werden in diesem Ticket auf 'Unfertiges bleibt im Milestone' umgeschrieben. DoD 1-3 beschreiben noch das alte road:-am-ersten-Milestone-Modell und werden nach den Antworten neu geschrieben.
- **2026-11-03 17:23 · berkc** — OFFENE FRAGEN (Brainstorm 03.11, je eine Wahl, Empfehlung zuerst):
F1 Road-Liste vs after: - A (Empfohlen) Die Liste der Road IST die Reihenfolge in der Road; after: bleibt nur fuer Pfeile zwischen Roads; bestehende Swimlanes werden beim ersten Schreiben als Liste uebernommen (Baum flach nach Tiefe, dann Alter). Gibt eine Wahrheit; kostet: Geschwister in einer Lane gehen verloren (werden zwei Roads). - B Road listet nur Mitglieder, Reihenfolge/Layout kommt weiter aus after:. Behaelt Geschwister; kostet: zwei Wahrheiten (Listenordnung vs after:), die auseinanderlaufen koennen. - C Road nennt nur ihren Top-Milestone, der Rest folgt aus after:. Kleinstes Format; widerspricht 'geordnete Liste ihrer Milestones'.
F2 Ein Milestone in zwei Roads - A (Empfohlen) verboten: create/add verweigert (exit 3); kommt es per git-Merge doch zustande, zeichnet die aeltere Road ihn und validate meldet es. - B erlaubt, in beiden Roads gezeichnet (doppelte Bubble). - C erlaubt, gezeichnet in der ersten, die zweite zeigt einen Verweis.
F3 Wo liegt eine Road - A (Empfohlen) wie Milestones: Datei .jaira/roads/<name>-<ID>.md (eine Zeile pro Milestone-ID) plus Ref je ID; hand-editierbar, merge-faehig, Roadmap liest wie heute Dateien. - B nur auf Refs wie Epics (refs/jaira/roads/<name>); kein Datei-Konflikt im Branch, aber nicht im Diff sichtbar und offline nur ueber Outbox.
F4 Epic-Etikett pro Board - A (Empfohlen) eigenes Ticket: das TUI zeigt heute gar keine Epics, also braucht es erst eine Epic-Ansicht; dieses Ticket liefert nur Roads. Ablageort dort: eine committete Datei je Board neben den Lanes (.jaira/lanes/epic bzw. .jaira/boards/<n>/lanes/epic) mit milestone|road|epic. - B hier mit: Datei wie in A, und das TUI bekommt ein 'Epic'-Etikett an Karte/Roadmap. Macht das Ticket etwa doppelt so gross. - C Einstellung in ~/.jaira/settings.json (pro Person) - widerspricht 'jedes Team entscheidet'.
- **2026-11-03 22:58 · berkc** — Maintainer 03.11.2026 zu F1-F4: F1 'Eine Road ist eine ganze Swimlane' - die Road umfasst die ganze Swimlane samt Verzweigungen; Aufbau innerhalb bleibt aus after:, die Road fasst ihre Milestones zusammen (Milestones kennen ihre Road nicht). F2 verboten (exit 3, validate meldet Merge-Faelle). F3 nur Ref (wie Epics) - reicht, damit andere Klone/Boards sie sehen. F4 nicht in KQQD1H: 'Epic' wird vorerst ein Milestone sein; die Wahl gehoert in die Web-Zuordnungsseite (VD4ZDW).
- **2026-11-03 23:00 · berkc** — Brainstorm nach den Antworten 03.11/04.11 (Worker), selbst entschieden: (1) Road wie Epic per Name auf dem Ref (refs/jaira/roads/<name>, slug), Anzeige ueber title; Umbenennen im TUI aendert nur title - so braucht es kein Loeschen eines Refs (Refs werden nie geloescht). (2) Zu F1 'ganze Swimlane' + Mitgliederliste: Lane eines Milestones = die Road, die ihn listet; sonst die Lane seines ersten Vorgaengers (wie heute); sonst eigene Lane. Ein erster Vorgaenger in einer anderen Lane wird zum Pfeil zwischen Lanes, der Milestone ist dann ein Top seiner Lane; eine Lane kann so mehrere Tops nebeneinander haben. Ein neuer Milestone nach einem Road-Mitglied steht damit ohne Schreiben der Road in ihrer Lane - die Liste muss nicht bei jedem Platzieren nachgezogen werden. (3) Doppelte Mitgliedschaft nach Merge: die aeltere Road (created-at, dann Name) bekommt den Milestone. (4) F4 Epic-Etikett nicht hier (VD4ZDW). DoD 1-4 umformuliert, 5-6 neu.
- **2026-11-03 23:02 · berkc** — Plan (pre-process 04.11): Road spiegelt den Epic-Weg auf Refs (gitref -> outbox Kind -> refsync mit Merge bei verlorenem Rennen), weil F3 'nur Ref wie Epics' ist und Fetch schon refs/jaira/* holt - kein neuer Sync-Mechanismus. Eigenes Paket core/road statt Erweiterung von core/epic: Epic hat Externals, Status, Textkoerper und Ticket-Mitglieder, eine Road nichts davon. Layout bekommt die Road-Listen als [][]string, damit core/milestone core/road nicht importiert (road importiert milestone.ParseMember). Navigation im TUI ist geometrisch (layerScene.step), daher reicht ein Ziel onRoad in einer neuen Kopfzeile ueber der Lane; Pfeil hoch vom Top erreicht es ohne Sonderfall. Name des Refs fest (slug), Umbenennen im TUI aendert nur title - Refs werden nie geloescht.
- **2026-11-03 23:03 · berkc** — Step 1 gitref: RoadPrefix/RoadRefName/RoadSHA/ReadRoad/ReadRoadAt/WriteRoad/ListRoads mirror the epic ones. Test core/gitref/road_test.go mirrors milestone_test.go (no epic gitref test existed).
- **2026-11-03 23:04 · berkc** — Step 2 outbox: KindRoad (dir roads/) at queueKind ref switch, readDir kinds list, lost-race keep-queued condition, send switch; Sender gained WriteRoad (gitref.Repo has it). Test TestARoadIsQueuedSentAndKeptOnARace in outbox_test.go.
- **2026-11-03 23:07 · berkc** — Plan deviation (step 4): the plan says members are read with milestone.ParseMember, but that only accepts ticket ULIDs (ticket.ValidID); milestone IDs are six characters (milestone.IsID). core/road has its own parseMember on milestone.IsID. Plan step 3 (refsync/road.go) needs core/road to compile, so steps 3 and 4 were committed together after both tests passed. refsync has no epic rebase test; road_test.go there covers the lost race (TestAConcurrentRoadListingIsMergedNotLost).
- **2026-11-03 23:09 · berkc** — Step 5 layout: Layout(all, roads [][]string). Lane key = first listing road, else lane of first usable predecessor, else own top. A child joins its parent only when both have the same lane key; otherwise its arrow is Extra. nil roads: all 13 existing tests pass with ', nil' added; layerdraw.go passes nil until step 9 wires road.Groups.
