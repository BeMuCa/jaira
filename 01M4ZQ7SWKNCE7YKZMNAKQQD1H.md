---
id: 01M4ZQ7SWKNCE7YKZMNAKQQD1H
title: "Die Swimlanes der Roadmap heissen Roads, tragen einen Namen und lassen sich im TUI umbenennen"
status: brainstorm
ready: true
creator: BeMuCa
assignee: berkc
goal: "Eine Road ist ein eigenes Objekt mit Name, Beschreibung und geordneter Milestone-Liste, das die Roadmap als benannte Swimlane zeichnet (ohne Road Road-1..Road-n) und das man im TUI per Pfeiltaste waehlt und mit Enter umbenennt; Milestones wissen nichts von ihrer Road, Unfertiges bleibt im Milestone"
context: |-
  Was heute fehlt: Die Roadmap (L, aus D0SAHM) zeigt Swimlanes, also Ketten von Milestones, die ueber after: verbunden sind. Eine Swimlane hat keinen Namen.
  Der Anfang einer Swimlane ist ein Milestone ohne Vorgaenger ('jaira milestone create --after ""' macht ihn zum Anfang einer eigenen Swimlane).
  Berk am 15.10.: 'die einzelnen Swimlanes als Milestone Road bezeichnen, und die sollen jeweils auch Namen bekommen. Per Default zaehlen sie einfach hoch, Road-1 .. Road-n. Aber man kann auf sie wechseln per Pfeiltaste und mit Enter deren Bezeichnung dann aendern. Dafuer muss man den Milestones dann wahrscheinlich noch einen Road-Namen hinzufuegen?'
  Vorschlag: Den Namen traegt nur der erste Milestone der Road, als Feld 'road: <name>'. Die anderen erben ihn ueber die Kette, damit es keine doppelte Wahrheit gibt. Ein Milestone wird wie heute ueber seinen Ref geteilt (seit EBTX6X ueber die ID).
  Offen fuer brainstorm: Was passiert, wenn zwei Roads zusammenlaufen (ein Milestone mit after auf zwei Roads), und welche Pfeiltasten (in der Roadmap bewegen sich Pfeile heute zwischen Bubbles).
definition-of-done: "Jede Swimlane der Roadmap zeigt einen Road-Namen: 'road:' ihres ersten Milestones, sonst Road-1..Road-n in Anzeigereihenfolge; Test"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-15T12:05:22Z
updated-at: 2026-11-03T23:00:24Z
updated-by: berkc
claimed-by: berk-ca-130584
claimed-at: 2026-11-03T17:20:46Z
question: "Brainstorm: 4 Designfragen offen (Road-Liste vs after:, Milestone in zwei Roads, Ablage der Road, Epic-Etikett) - siehe Notiz OFFENE FRAGEN, je mit Empfehlung"
---

# Die Swimlanes der Roadmap heissen Roads, tragen einen Namen und lassen sich im TUI umbenennen

## Definition of Done

- [ ] Eine Road ist ein eigenes Objekt nur auf ihrem Ref refs/jaira/roads/<name> (wie Epics): Titel, Beschreibung, Liste ihrer Milestone-IDs; Milestones tragen kein Road-Feld; paralleles Schreiben wird gemergt wie bei Epics; Test
- [ ] In der Roadmap waehlt eine Pfeiltaste die Road, Enter oeffnet ein Namensfeld, Enter speichert 'road:' am ersten Milestone (geteilt wie jede Milestone-Aenderung), Esc bricht ab; Test
- [ ] 'jaira milestone set <name> --road <road>' setzt den Namen auch ueber die CLI; 'milestone ls --json' traegt road; Test
- [ ] README, docs/COMMANDS.md und eine Zeile in core/release/notes/<ID>.md

## Options

- [x] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

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
