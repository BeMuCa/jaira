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
blocked-by: []
related: []
commits: []
created-at: 2026-10-19T06:43:02Z
updated-at: 2026-10-19T18:18:27Z
updated-by: Alexander Sacharov
mode: conversational
claimed-by: DESKTOP-RFTCH11-25649
claimed-at: 2026-10-19T18:18:27Z
---

# Titel und Beschreibung eines Epics lassen sich im Web und mit jaira epic set aendern

## Definition of Done

- [ ] core/epic: eine Funktion schreibt Titel und/oder Beschreibung eines bestehenden Epics auf seinen Ref, Mitglieder bleiben unveraendert; Test
- [ ] 'jaira epic set <name> --title <...> --description <...>' aendert eines oder beides; unbekanntes Epic Exit 5; ohne beide Flags Exit 2
- [ ] Web: eine Route (z.B. PATCH /api/epics/{name}) und in der UI ein Bearbeiten-Knopf am Epic, Titel und Beschreibung editierbar, Speichern zeigt sofort den neuen Stand
- [ ] Hat jemand das Epic zwischenzeitlich geaendert, geht keine der beiden Aenderungen still verloren (Test)
- [ ] core/release/notes/<ID>.md; go test ./... -race gruen
- [ ] Web zeigt den ganzen Markdown-Text eines Epics gerendert; 'jaira epic show' gibt ihn ebenfalls aus
- [ ] Der Text ist standardmaessig eingeklappt (die description-Zeile bleibt sichtbar) und klappt per Klick auf
- [ ] Beim Schreiben des Textes bleiben die Mitgliedszeilen des Epics unangetastet (Test)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [ ] core/epic: Body() = alle Nicht-Mitgliedszeilen unter dem Frontmatter; SetBody(text) schreibt Frontmatter + Text + Mitgliedszeilen in alter Reihenfolge; Set(store, refs, name, Patch{Title,Description,Body *string, Base}) unter dem epics-Lock; Tests
- [ ] core/epic Merge: Body wie title/description mergen; gleichzeitige Aenderung desselben Feldes nach Entscheidung (siehe Frage); Test
- [ ] CLI: 'jaira epic set <name> --title --description --body --body-file <pfad|->' (Exit 5 unbekannt, Exit 2 ohne Flag); 'epic show' gibt den Text aus, --json mit body; Tests
- [ ] Web: Epic-JSON bekommt body; PATCH /api/epics/{name} mit base (was der Editor gelesen hat) -> 409 wenn inzwischen geaendert; Source.SetEpic in internal/cli/web_source.go; Tests
- [ ] UI: EpicList zeigt description-Zeile, Text gerendert (shared/lib/markdown) und eingeklappt; Bearbeiten-Knopf -> Formular Titel/description/Text; Speichern laedt Epics neu
- [ ] core/release/notes/CRDH6B.md; go test ./... -race; UI-Build/Tests

## Progress
- **2026-10-19 06:45 · Alexander Sacharov** — parent session: jaira-a7
- **2026-10-19 06:45 · Alexander Sacharov** — Befund (Teamlead, 19.10.2026): Epic 'anforderungen-im-chat' in ~/projects/requirementsgenie traegt ~95 Zeilen Markdown unter dem Frontmatter (## Ziel, ## Ausgangslage, ...): git -C ~/projects/requirementsgenie show refs/jaira/epics/anforderungen-im-chat:anforderungen-im-chat.md. Das Web zeigt nur die einzeilige Frontmatter-description. Ursache: core/epic/epic.go FromBytes behaelt unter dem Frontmatter nur Mitgliedszeilen; EpicList.svelte:29 rendert nur epic.description.
- **2026-10-19 06:45 · Alexander Sacharov** — decision (Alex, 19.10.2026): 'Beschreibung bearbeiten' heisst in diesem Ticket der Markdown-Text unter dem Frontmatter, nicht nur die description-Zeile. Mitgliedszeilen bleiben beim Schreiben unangetastet. Der Text wird gerendert angezeigt und ist standardmaessig eingeklappt.
- **2026-10-19 06:45 · Alexander Sacharov** — Reihenfolge (Teamlead): blocked-by TPCKZB (gleiche core/web/ui-Dateien, feat/tpckzb in-progress). Jetzt nur brainstorm/plan; kein Code, bevor TPCKZB auf master ist.
- **2026-10-19 06:46 · Alexander Sacharov** — decision (Worker, aus Code/Default): 'description' bleibt die einzeilige Frontmatter-Zeile (sichtbar im eingeklappten Zustand), der Text heisst 'body'. CLI bekommt dafuer --body und --body-file (- = stdin) neben --title/--description; Web-Formular hat alle drei Felder. Grund: die Zeile und der Text sind zwei Dinge auf dem Ref; ein Flag fuer beides wuerde eins davon ueberschreiben.
- **2026-10-19 06:46 · Alexander Sacharov** — decision (Worker): Body = alle Zeilen unter dem Frontmatter ausser '- <id>'-Mitgliedszeilen. SetBody setzt die Mitgliedszeilen in alter Reihenfolge ans Ende hinter den neuen Text. Bei heutigem Layout (## Tickets + Kommentar am Ende, Mitglieder darunter) bleibt die Datei damit gleich aufgebaut. Ein '# <Titel>' als erste Ueberschrift im Body wird bei --title mit umbenannt, wenn sie dem alten Titel entsprach.
- **2026-10-19 06:46 · Alexander Sacharov** — decision (Worker): lokale Gleichzeitigkeit (Web-Formular offen, CLI aendert dasselbe Epic) -> PATCH traegt 'base' = die Werte, die der Editor gelesen hat; ein Feld, das der Editor aendern will und das sich seitdem geaendert hat -> 409 + aktueller Stand, nichts geschrieben. CLI-set schreibt unter dem Lock auf den aktuellen Stand und braucht kein base.
- **2026-10-19 06:47 · Alexander Sacharov** — Frage an Teamlead/Alex (offen): Zwei Klone aendern denselben Epic-Text gleichzeitig; refsync verliert das Rennen und ruft epic.Merge. Heute gewinnt dort 'ours' still - verletzt DoD 4. Optionen: 1) (empfohlen) Konfliktmarker: der Text wird ours, der fremde Text bleibt darunter in einem markierten Block '<!-- jaira: gleichzeitig geaendert, bitte zusammenfuehren -->' - sichtbar im Web und im Diff, kein Verlust, ~30 Zeilen. 2) Zeilenweiser 3-Wege-Merge (diff3 selbst gebaut) - nicht-ueberlappende Aenderungen fliessen still zusammen, nur echte Ueberlappung braucht Marker; deutlich mehr Code. 3) ours gewinnt, theirs geht verloren wie heute - einfach, aber DoD 4 muesste umformuliert werden. Fuer Titel/description (einzeilig) gilt bei 1) dasselbe: ours gewinnt, der fremde Wert landet im Konfliktblock.
- **2026-10-19 06:48 · Alexander Sacharov** — decision (Alex via Teamlead, 19.10.2026): Option 1. Gleichzeitige Aenderung desselben Feldes in epic.Merge: ours wird der Wert/Text, theirs bleibt darunter im Body in einem markierten Konfliktblock '<!-- jaira: gleichzeitig geaendert -->' erhalten - gilt fuer Body, Titel und description. Nichts geht still verloren.
- **2026-10-19 18:17 · Alexander Sacharov** — decision (Alex via Teamlead, 19.10.2026): nicht auf master warten. feat/crdh6b wird auf feat/tpckzb aufgebaut (Worktree steht heute auf 7f008dec = master, noch kein eigener Commit: zuerst 'git reset --keep feat/tpckzb'). Code ab jetzt erlaubt. PR von CRDH6B erst nach dem PR von TPCKZB oeffnen; im Push-Bericht sagen, dass der Branch gestapelt ist.
