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
updated-at: 2026-11-02T20:41:11Z
updated-by: Alexander Sacharov
outcome-what: "Epics tragen 'external:' (Blockliste KEY URL) im Frontmatter ihres Refs; epic create/set --external/--external-rm, epic show/ls zeigen sie, epic ls --external und list --external filtern; jaira web zeigt Links, editiert und filtert nach Key."
outcome-why: "Alex rechnet Stunden ueber JIRA-Stories ab: ein Epic je Story, und der Schluessel darf nicht in Code-Commits landen - darum nur auf dem Epic-Ref."
outcome-resolves: "Epics und ihre Tickets lassen sich nach dem Schluessel einer Aufgabe in einem anderen System finden, ohne dass der Schluessel ein Ticket-File beruehrt."
review-summary: "Epics tragen jetzt im Frontmatter ihres Refs eine Blockliste 'external:' mit je einer Zeile '<KEY> [<URL>]' (core/epic/external.go). 'jaira epic create|set --external' fuegt einen Schluessel hinzu oder verlinkt ihn neu (Vergleich ohne Gross/Klein, URL nur http/https, sonst exit 2), '--external-rm' nimmt ihn weg; der letzte entfernte Schluessel nimmt den Block mit, alte Epics bleiben byte-gleich. Von Hand geschriebene Formen (Wert auf der 'external:'-Zeile, gequotete Zeilen, Nicht-URL) werden gelesen und beim naechsten Schreiben erhalten. Merge im Ref-Rennen wendet Hinzufuegen/Entfernen je Schluessel auf die andere Seite an. 'epic show' und 'epic ls' (auch --json) zeigen die Schluessel, 'epic ls --external' filtert Epics, 'jaira list --external' zeigt Mitglieder aller Epics mit dem Schluessel samt Folge-Tickets (geschlossene Epics eingeschlossen). jaira web liefert je Schluessel key/url und ein 'link' nur fuer http(s); die Epic-Liste zeigt ihn als Link, der Epic-Editor hat KEY/URL-Zeilen und schickt add/drop als Differenz, die FilterBar hat einen Key-Chip, der Epics und Tickets (ueber die Schluessel ihrer Epics) filtert. Kein Ticket-Ref, keine Datei im Arbeitsbaum und kein Commit ausserhalb refs/jaira/epics traegt den Schluessel (Test). Release-Note in core/release/notes/PG9P9Y.md."
round: 2
review-gaps: |-
  note: EditEpicForm.svelte - das URL-Feld ist type="url"; traegt ein von Hand editiertes Epic eine Zeile wie 'C-3 not a link', blockiert die Browser-Validierung das Speichern des ganzen Formulars (auch Titel/Text), bis man die Zeile leert oder repariert
  note: der Merge-Pfad fuer Schluessel ist nur auf Ebene core/epic getestet (TestMergeKeepsBothSidesKeyChanges), nicht ueber ein echtes Rennen zweier Clones um den Ref
  note: svelte-check meldet einen Fehler in vite.config.ts (node:url-Typen) - besteht schon vor diesem Diff, nicht durch ihn
test-verdict: "pass: go build + go test -race -count=1 ./... RC=0; ui npm test 21/21; DoD 1-6 in der Arbeitskopie geprueft (Proof-Tests laufen gruen); am gebauten Binary auf einem Scratch-Board: epic create/set --external, ls zeigt [JIRA-1234, YT-7], show zeigt 'external:'-Zeilen, list --external jira-1234 zeigt Mitglied + Folge-Ticket, javascript:-URL exit 2, --external-rm entfernt; Ref traegt Blockliste; jaira web: GET /api/epics liefert key/url/link, PATCH add/drop 200, javascript: 400, das gebaute JS enthaelt Editor-Zeilen und noopener-Links. Nicht im Browser angeklickt."
claimed-by: DESKTOP-RFTCH11-47320
claimed-at: 2026-11-02T20:26:29Z
review-verdict: "Der Diff erfuellt alle sechs DoD-Punkte; keine Defekte gefunden. go build/vet, go test ./core/epic ./core/web ./internal/cli (Epic|External|Web|List), npm test (21/21) gruen; core/web/ui/dist entspricht einem frischen vite build der Quellen. Drei Notizen, kein Blocker."
review-round: 2
review-check: |-
  1. cd /home/alex/projects/.worktrees/jaira-pg9p9y && go build -o /tmp/jaira-pg ./cmd/jaira
  2. /tmp/jaira-pg epic ls  - sich ein vorhandenes Epic merken, z. B. <name>
  3. /tmp/jaira-pg epic set <name> --external "JIRA-1 https://jira.example/browse/JIRA-1"
  4. /tmp/jaira-pg epic show <name>  - Zeile 'external: JIRA-1 https://jira.example/browse/JIRA-1' ist da
  5. /tmp/jaira-pg epic ls --external jira-1  - nur dieses Epic, mit [JIRA-1] hinter dem Titel
  6. /tmp/jaira-pg list --external JIRA-1  - die Tickets des Epics samt Folge-Tickets
  7. git status  - keine Datei unter .jaira/tickets hat sich geaendert
  8. /tmp/jaira-pg epic set <name> --external "X javascript:alert(1)"  - Fehler, exit 2
  9. /tmp/jaira-pg web  - am Epic steht JIRA-1 als Link; Epic bearbeiten zeigt eine Zeile JIRA-1/URL mit x; im Filter 'Key' JIRA-1 waehlen - nur dieses Epic und seine Tickets bleiben
  10. /tmp/jaira-pg epic set <name> --external-rm JIRA-1  - epic show zeigt keine external-Zeile mehr
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
- [ ] web: eine ungueltige URL in einer external-Zeile markiert nur diese Zeile; der Rest des Epic-Formulars laesst sich speichern; Test

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
