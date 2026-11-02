---
id: 01M6EY4FPSNN0693MDECPG9P9Y
title: "Ein Epic traegt einen externen Schluessel mit Link, und Epics und Tickets lassen sich danach suchen und filtern"
status: signoff
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: Ein Epic hat ein eigenes Feld fuer den Schluessel einer Aufgabe in einem anderen System (z. B. JIRA-1234) mit URL. Danach kann man in CLI und web Epics und ihre Tickets suchen und filtern.
context: |-
  Warum (Alex, 02.11.2026):
  - Alex will Stunden ueber echte JIRA-Stories abrechnen: pro JIRA-Story ein Epic, die Tickets darin sind die Arbeit, die darueber bezahlt wird.
  - Die Abrechnung soll NICHT in Code-Commits landen. Deshalb das Feld auf dem Epic (eigener Ref refs/jaira/epics/<name>), nicht auf dem Ticket: die Ticket-Datei wird mit dem Code committet, der Epic-Ref nicht.
  - JIRA-Stories und Epics ueberschneiden sich bei der Bezahlung. Ein Ticket darf also in zwei Epics liegen (geht heute schon mit jaira epic add).
  Was heute fehlt:
  - Epic-Frontmatter hat nur title, description, created-at (core/epic/epic.go). Einen Schluessel kann man nur in description schreiben: nicht filterbar, kein Link.
  Entwurf (brainstorm darf aendern):
  - Feld 'external:' im Epic-Frontmatter, eine Liste, je Zeile '<SCHLUESSEL> <URL>'. Allgemein, nicht 'jira' - es gibt auch YouTrack.
  - jaira epic create|set --external '<KEY> <URL>' (wiederholbar), --external-rm <KEY>.
  - jaira epic ls --external <KEY> und jaira list --external <KEY> (Tickets aller Epics mit diesem Schluessel).
  - web: Schluessel als Link am Epic, Filter nach Schluessel.
  Achtung: der Epic-Ref wird zum Remote gepusht - alle mit Zugriff aufs Repo sehen die Schluessel. Fuer Alex ok (bewusst gewaehlt statt einer lokalen Datei).
  Haengt an 182JT1 (PR #120, epic close), das dieselben Dateien in core/epic und web aendert: erst nach dessen Merge anfangen.
  Spaeter, nicht hier: Stundensumme je Schluessel in /jaira-role-logbook-summary.
definition-of-done: "Epic-Frontmatter hat 'external:' (Liste '<KEY> <URL>'); jaira epic create/set --external fuegt hinzu, --external-rm nimmt weg; Datei bleibt von Hand editierbar; Test"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-11-02T20:10:12Z
updated-at: 2026-11-02T22:06:58Z
updated-by: Alexander Sacharov
outcome-what: "Epics tragen 'external:' (Blockliste KEY URL) im Frontmatter ihres Refs; epic create/set --external/--external-rm, epic show/ls zeigen sie, epic ls --external und list --external filtern; jaira web zeigt Links, editiert und filtert nach Key."
outcome-why: "Alex rechnet Stunden ueber JIRA-Stories ab: ein Epic je Story, und der Schluessel darf nicht in Code-Commits landen - darum nur auf dem Epic-Ref."
outcome-resolves: "Epics und ihre Tickets lassen sich nach dem Schluessel einer Aufgabe in einem anderen System finden, ohne dass der Schluessel ein Ticket-File beruehrt."
review-summary: "Epics carry an 'external:' block list (one '<KEY> <URL>' line each) in the frontmatter on their own ref; 'epic create/set --external/--external-rm' add and take keys (case-insensitive, http(s) only), 'epic show/ls' (+--json) print them, 'epic ls --external' and 'list --external' find epics and their tickets including follow-ups, and no ticket file or code commit carries a key. jaira web shows keys as links, edits them as KEY/URL rows and filters by a Key chip. A bad link marks only its own row and the rest saves; round 3 makes each read row remember its key (from), so renaming a row whose hand-written link is bad keeps the old line instead of dropping it (heldAs in edit.ts is the single rule, used by keysEdit and the hint). Every place an epic is named, 'epic add/rm' summary line included, shows the title with the name beside it."
round: 5
review-gaps: |-
  note: edit.test.ts does not cover clearing the key of a bad-link row (heldAs keeps the line, which is right, but no test pins it)
  note: svelte-check reports one error in core/web/ui/vite.config.ts (missing node types); the file is not touched by this branch
test-verdict: "pass: go build/vet ./..., go test ./... -race RC=0 (33 Pakete ok), npm test 24/24 RC=0, dist = frischer vite build; DoD 9 per edit.test.ts + Quelle (Zeilen tragen from), DoD 10 live auf Scratch-Board: epic add/rm drucken '\"Web 0.2\" (web-0-2) now lists N ticket(s).'"
claimed-by: DESKTOP-RFTCH11-91590
claimed-at: 2026-11-02T21:59:27Z
review-verdict: "the diff matches all ten DoD items; no blocker. go build/vet, go test core/epic core/web internal/cli, npm test 24/24 green, and core/web/ui/dist is byte-identical to a fresh vite build of the current sources"
review-round: 5
review-check: |-
  1. go build -o /tmp/jaira ./cmd/jaira, then cd into a scratch repo that has a jaira board
  2. /tmp/jaira epic create web-0-2 --title "Web 0.2" --external "JIRA-12 https://jira.example/12" - you see: Epic "Web 0.2" (web-0-2) created
  3. /tmp/jaira epic add web-0-2 <ticket-id> - the last line reads: "Web 0.2" (web-0-2) now lists 1 ticket(s).
  4. /tmp/jaira list --external jira-12 - the ticket is listed (lower case matches); git status shows no changed ticket file
  5. /tmp/jaira web, open the epic: JIRA-12 is a link; in the editor type "not a link" in its URL and change the title, save - the title is saved, only that row is marked, and jaira epic show web-0-2 still has the old URL
  6. the rename case needs a hand-written bad link on the ref and has no quick UI path: cd core/web/ui && npm test - "renaming a row with a bad link keeps the line it was read as" passes
---

# Ein Epic traegt einen externen Schluessel mit Link, und Epics und Tickets lassen sich danach suchen und filtern

## Definition of Done

- [x] Epic-Frontmatter hat 'external:' (Liste '<KEY> <URL>'); jaira epic create/set --external fuegt hinzu, --external-rm nimmt weg; Datei bleibt von Hand editierbar; Test
  proof: core/epic/external.go writeExternal; tests TestExternalKeysAreABlockListInTheFrontmatter, TestExternalKeysReadAHandEditedFile, TestMergeKeepsBothSidesKeyChanges; CLI TestEpicExternalKeysFindTheWork
- [x] jaira epic show und epic ls (auch --json) zeigen die Schluessel; jaira epic ls --external <KEY> filtert; Test
  proof: internal/cli/epics.go keysMark/epicJSON external, ls --external; test TestEpicExternalKeysFindTheWork
- [x] jaira list --external <KEY> zeigt die Tickets aller Epics mit diesem Schluessel, Folge-Tickets eingeschlossen; Test
  proof: internal/cli/epics.go ticketsForKey + list --external in internal/cli/tickets.go; test TestEpicExternalKeysFindTheWork (member + follow-up)
- [x] web: Schluessel am Epic als klickbarer Link, setzen/entfernen im Epic-Editor, Filter nach Schluessel; Test
  proof: core/web/ui/src/widgets/epic-list/ui/EpicList.svelte (links), features/edit-epic/ui/EditEpicForm.svelte (rows), features/filter-tickets/model/filter.ts epicPasses/key chip; tests TestWebEditsAnEpicsExternalKeys, edit.test.ts 'the keys go as...', filter.test.ts 'the key filter...'
- [x] Kein Ticket-File und kein Code-Commit traegt den Schluessel - er steht nur im Epic-Ref; Test
  proof: test TestAnExternalKeyStaysOnTheEpicRef (ticket ref unchanged, no working-tree file, no commit outside refs/jaira/epics)
- [x] Release-Note-Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/PG9P9Y.md
- [x] web: eine ungueltige URL in einer external-Zeile markiert nur diese Zeile; der Rest des Epic-Formulars laesst sich speichern; Test
  proof: core/web/ui/src/features/edit-epic/ui/EditEpicForm.svelte (Link-Feld ohne type=url, Zeile .bad + Hinweis), model/edit.ts linkable/keysEdit; test edit.test.ts 'a bad link holds back its own row only'
- [x] Wo ein Epic genannt wird (web und CLI), steht sein Titel, der Name klein daneben; Test
  proof: core/epic/epic.go Label/LabelOf, internal/cli/epics.go titleAndName, web_source.go messages, core/web/ui/src/entities/epic/model/label.ts + EpicList.svelte; tests TestTheCLINamesAnEpicByItsTitle, TestLabel, label.test.ts
- [x] Umbenennen des KEY einer Zeile mit ungueltigem Link im Epic-Editor verliert die Zeile nicht; Test
  proof: core/web/ui/src/features/edit-epic/model/edit.ts heldAs + keysEdit (held), EditEpicForm.svelte rows carry from; test edit.test.ts 'renaming a row with a bad link keeps the line it was read as'
- [x] epic add/rm nennen in der Schlusszeile Titel und Namen; Test
  proof: internal/cli/epics.go editEpic summary line uses e.Label(); test TestTheCLINamesAnEpicByItsTitle (add/rm 'now lists')

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/epic: External-Typ, Parse/Schreiben der Blockliste, AddExternal/DropExternal in Patch, Merge je KEY, Tests
- [x] CLI: epic create/set --external/--external-rm, epic show/ls (+json) zeigen, epic ls --external, list --external, Tests
- [x] web: Go-Typen (Epic.External, EpicEdit add/drop), webSource, Tests
- [x] web UI: Links am Epic, Editor-Zeilen, Key-Chip im Filter, vitest; dist bauen
- [x] Test: Ticket-Datei bleibt unveraendert; Release-Note core/release/notes/PG9P9Y.md

## Progress
- **2026-11-02 20:11 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/43623.sock
- **2026-11-02 20:11 · Alexander Sacharov** — Branch auf upstream/master 783c6009 gesetzt (enthaelt 182JT1 PR #120 und DF70XX PR #119) - die Abhaengigkeit aus dem Kontext ist erfuellt.
- **2026-11-02 20:12 · Alexander Sacharov** — Offene Entscheidungen gezaehlt: keine fuer Alex. Selbst entschieden (Reviewer darf widersprechen):
1) Format: YAML-Blockliste im Frontmatter, 'external:' und darunter je Schluessel eine Zeile '  - <KEY> <URL>'. Grund: Kontext sagt 'je Zeile', eine Zeile pro Schluessel ist im Diff lesbar, gueltiges YAML (wiederholte 'external:'-Zeilen waeren doppelte Keys). Fehlt 'external:', schreibt jaira nichts (alte Epics bleiben byte-gleich).
2) KEY = erstes Wort (kein Leerzeichen), URL optional; ohne URL nur Text, kein Link. URL nur http(s) - CLI lehnt anderes ab (exit 2), web rendert href nur fuer http(s) (kein javascript:-Link).
3) Vergleich der Schluessel ohne Gross/Klein (jira-12 == JIRA-12), gespeichert wie geschrieben. --external mit vorhandenem KEY ersetzt dessen URL an seiner Stelle.
4) Patch arbeitet mit add/drop pro KEY statt Liste ersetzen - CLI und Web-Editor (der schickt die Differenz zu dem, was er gelesen hat). Merge im Rennen um den Ref: wie Mitglieder, Hinzufuegen/Entfernen je KEY wird auf theirs angewandt.
5) Web: Chip-Filter 'Key' in der FilterBar wie Tag; ein Ticket traegt die Schluessel aller Epics, die es zaehlen (Folge-Tickets eingeschlossen) - so filtert er Epic-Liste und lose Tickets gleich wie 'jaira list --external'. Epic-Editor: Zeilen KEY/URL mit Entfernen-Knopf und eine leere Zeile zum Hinzufuegen.
- **2026-11-02 20:21 · Alexander Sacharov** — critique 1: zwei kleine Funde - keysEdit doppelter Key-Vergleich (edit.ts), Short von 'epic set' nennt die Keys nicht. Sonst stehen gelassen: Inline-'external: KEY URL' lesen bleibt (verhindert Datenverlust beim naechsten writeExternal); web.External traegt URL und Link, damit LinkableURL nur in Go entschieden wird; ticketsForKey liegt bei den Epic-Helfern in epics.go.
- **2026-11-02 20:22 · Alexander Sacharov** — critique 2: beide Funde aus Runde 1 behoben (edit.ts keysEdit, epic set Short) - nichts weiter.
- **2026-11-02 20:24 · Alexander Sacharov** — review: an einen frischen Agenten uebergeben (Herdr-Tab via run-lane.sh --no-worktree).
- **2026-11-02 20:26 · Alexander Sacharov** — review: frischer Agent. Gelesen: Go- und UI-Diff 783c6009..HEAD (2 Commits), worktree_diff enthaelt nur fremde .jaira/milestones/ - nicht Teil des Tickets. Geprueft: go build/vet, go test core/epic core/web internal/cli, npm test, dist = frischer vite build. Kein Blocker; drei Notizen in review-gaps (type=url blockiert Speichern bei handeditierter Nicht-URL).
- **2026-11-02 20:28 · Alexander Sacharov** — Abnahmeseite: https://claude.ai/artifact/N5BogNLe6e1kXmNtufAqn2 (8 Maschinenpruefungen gruen auf einem Test-Board unter ~/.claude/acceptance/jaira/feat-pg9p9y/fixture, 5 Personen-Schritte in jaira web).
- **2026-11-02 20:41 · Alexander Sacharov** — decision: Alex (02.11., ueber Teamlead): der type=url-Fund aus dem Review wird in diesem Ticket behoben, kein eigenes Ticket. Die anderen zwei Review-Notizen bleiben Notizen.
- **2026-11-02 20:41 · Alexander Sacharov** — decision: Alex (02.11., ueber Teamlead): wo ein Epic genannt wird (web und CLI), steht sein Titel, der Name klein daneben (den tippt man in Befehle). Gleiches Ticket, kein eigenes.
- **2026-11-02 20:45 · Alexander Sacharov** — critique 3 (Runde fuer DoD 7+8): kein Fund. Stehen gelassen: epic.Label ("Web 0.2" (web-0-2), fuer Saetze) und cli titleAndName (Web 0.2  (web-0-2), fuer Tabellen-/Kopfzeilen) sind zwei Formen mit Absicht - Anfuehrungszeichen gehoeren in einen Satz, nicht in eine Spalte. Fehlermeldungen fuer ein unbekanntes Epic bleiben beim Namen (es gibt keinen Titel). web: TicketCard/Details/UnassignedList nennen kein Epic, nur EpicList-Kopf geaendert.
- **2026-11-02 20:46 · Alexander Sacharov** — review Runde 2: an einen neuen frischen Agenten uebergeben.
- **2026-11-02 20:47 · Alexander Sacharov** — review Runde 2: frischer Agent. Gelesen: 62869144 (DoD 7+8) gegen 783c6009..HEAD; worktree_diff nur fremde .jaira/milestones/. go build/vet, go test core/epic core/web internal/cli, npm test 23/23, dist = frischer vite build. Kein Blocker, drei Notizen in review-gaps.
- **2026-11-02 21:59 · Alexander Sacharov** — Entschieden von Alex (02.11.2026, ueber Teamlead), Runde 3: beide Review-Notizen werden in diesem Ticket behoben. (1) Datenverlust: KEY einer Zeile mit ungueltigem, von Hand geschriebenem Link im Epic-Editor umbenennen -> alter KEY weg, neuer nicht gespeichert. Die Zeile darf nie verloren gehen. (2) epic add/rm nennen in der Schlusszeile den Titel, den Namen klein daneben, wie ueberall sonst (DoD 8). Neuer Worker, weil der Cache des alten kalt war.
- **2026-11-02 22:01 · Alexander Sacharov** — Runde 3 umgesetzt: (1) Zeilen im Epic-Editor merken sich ihren gelesenen KEY (from). keysEdit haelt den from-KEY einer Zeile mit ungueltigem Link vom drop zurueck - Umbenennen ohne gueltigen Link wird nicht gespeichert, die alte Zeile bleibt; der Hinweis sagt 'stays as <from>'. Verworfen: Umbenennung trotzdem mit kaputtem Link senden - der Server lehnt nicht-http(s) ab (ParseExternal), das wuerde das ganze Speichern scheitern lassen. (2) epic add/rm Schlusszeile 'now lists' nimmt e.Label() statt name.
- **2026-11-02 22:01 · Alexander Sacharov** — critique 4 (Runde 3): core/web/ui/src/features/edit-epic/ui/EditEpicForm.svelte:78 compares row.key with row.from in its own case rule inside the template; keysEdit (model/edit.ts) owns 'keys compare in any case' - export heldAs(row) from edit.ts (the read key a bad-link renamed row stays as, else '') and use it in both keysEdit's held list and the hint
- **2026-11-02 22:02 · Alexander Sacharov** — critique 5: Fund aus critique 4 behoben (heldAs in edit.ts, Template und keysEdit nutzen es) - nichts weiter.
- **2026-11-02 22:04 · Alexander Sacharov** — review Runde 3: an einen neuen frischen Agenten uebergeben.
- **2026-11-02 22:06 · Alexander Sacharov** — review Runde 3: frischer Agent. Gelesen: 62869144..143e2ae8 (heldAs/from in edit.ts, EditEpicForm-Hinweis, epic add/rm Schlusszeile) gegen den ganzen Diff seit 783c6009; worktree_diff nur fremde .jaira/milestones/. go build/vet, go test core/epic core/web internal/cli, npm test 24/24, dist = frischer vite build (byte-gleich). Randfaelle geprueft: KEY leeren, Gross/Klein-Umbenennung, Umbenennung mit gutem Link, Entfernen-Knopf - keiner verliert eine Zeile. Kein Blocker, zwei Notizen.
