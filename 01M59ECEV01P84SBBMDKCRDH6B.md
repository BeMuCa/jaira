---
id: 01M59ECEV01P84SBBMDKCRDH6B
title: Titel und Beschreibung eines Epics lassen sich im Web und mit jaira epic set aendern
status: todo
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer ein Epic angelegt hat, kann seinen Titel und seine Beschreibung spaeter neu schreiben - auf der Seite von 'jaira web' und in der CLI."
context: |-
  Was heute fehlt: ein Epic laesst sich nach dem Anlegen nirgends mehr aendern.
  - Web (core/web/web.go:185-191) kann Epics nur anlegen (POST /api/epics) und Tickets zuordnen/entfernen. Keine Route zum Bearbeiten.
  - CLI 'jaira epic' kennt nur create, add, rm, ls, show. Kein set.
  - core/epic/epic.go hat setField (Zeile 200) und Bytes, aber keine oeffentliche Funktion, die Titel oder Beschreibung neu schreibt.
  Anlass: Alex am 19.10.2026 - 'im Web muss man Epics bearbeiten, also neu schreiben koennen'.
  Entschieden von Alex (19.10.2026): nur Titel und Beschreibung. Der Name (= Ref refs/jaira/epics/<name>) bleibt, kein Umbenennen, kein Loeschen.
  Warum blocked-by TPCKZB: TPCKZB baut gerade Filter in dieselbe Web-UI (core/web/ui/src, core/web/web.go). Zwei Aenderungen an denselben Dateien gleichzeitig -> erst nach TPCKZB.
  Beim Schreiben bedenken: ein Epic liegt nur auf seinem Ref; zwei Leute koennen es gleichzeitig bearbeiten. Wie create/add heute mit einem zwischenzeitlich geaenderten Ref umgehen, dasselbe Muster nehmen.
definition-of-done: "core/epic: eine Funktion schreibt Titel und/oder Beschreibung eines bestehenden Epics auf seinen Ref, Mitglieder bleiben unveraendert; Test"
tags:
  - web
  - cli
blocked-by:
  - TPCKZB
related: []
commits: []
created-at: 2026-10-19T06:43:02Z
updated-at: 2026-10-19T06:45:41Z
updated-by: Alexander Sacharov
---

# Titel und Beschreibung eines Epics lassen sich im Web und mit jaira epic set aendern

## Definition of Done

- [ ] core/epic: eine Funktion schreibt Titel und/oder Beschreibung eines bestehenden Epics auf seinen Ref, Mitglieder bleiben unveraendert; Test
- [ ] 'jaira epic set <name> --title <...> --description <...>' aendert eines oder beides; unbekanntes Epic Exit 5; ohne beide Flags Exit 2
- [ ] Web: eine Route (z.B. PATCH /api/epics/{name}) und in der UI ein Bearbeiten-Knopf am Epic, Titel und Beschreibung editierbar, Speichern zeigt sofort den neuen Stand
- [ ] Hat jemand das Epic zwischenzeitlich geaendert, geht keine der beiden Aenderungen still verloren (Test)
- [ ] core/release/notes/<ID>.md; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-19 06:45 · Alexander Sacharov** — parent session: jaira-a7
- **2026-10-19 06:45 · Alexander Sacharov** — Befund (Teamlead, 19.10.2026): Epic 'anforderungen-im-chat' in ~/projects/requirementsgenie traegt ~95 Zeilen Markdown unter dem Frontmatter (## Ziel, ## Ausgangslage, ...): git -C ~/projects/requirementsgenie show refs/jaira/epics/anforderungen-im-chat:anforderungen-im-chat.md. Das Web zeigt nur die einzeilige Frontmatter-description. Ursache: core/epic/epic.go FromBytes behaelt unter dem Frontmatter nur Mitgliedszeilen; EpicList.svelte:29 rendert nur epic.description.
- **2026-10-19 06:45 · Alexander Sacharov** — decision (Alex, 19.10.2026): 'Beschreibung bearbeiten' heisst in diesem Ticket der Markdown-Text unter dem Frontmatter, nicht nur die description-Zeile. Mitgliedszeilen bleiben beim Schreiben unangetastet. Der Text wird gerendert angezeigt und ist standardmaessig eingeklappt.
- **2026-10-19 06:45 · Alexander Sacharov** — Reihenfolge (Teamlead): blocked-by TPCKZB (gleiche core/web/ui-Dateien, feat/tpckzb in-progress). Jetzt nur brainstorm/plan; kein Code, bevor TPCKZB auf master ist.
