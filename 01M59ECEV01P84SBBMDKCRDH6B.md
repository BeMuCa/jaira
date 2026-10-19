---
id: 01M59ECEV01P84SBBMDKCRDH6B
title: Titel und Beschreibung eines Epics lassen sich im Web und mit jaira epic set aendern
status: optimize
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
updated-at: 2026-10-19T20:29:53Z
updated-by: Alexander Sacharov
mode: conversational
claimed-by: DESKTOP-RFTCH11-25649
claimed-at: 2026-10-19T18:18:27Z
outcome-what: "core/epic: Body/Text/SetBody/Set with Patch+Base, Merge keeps both sides of a field in a marked block; CLI 'jaira epic set' and body in 'epic show'/--json; web PATCH /api/epics/{name} (409 on stale base) with an Edit form and the text rendered folded under the description"
outcome-why: "An epic could not be changed after it was created, and its markdown text below the frontmatter was invisible in the web (requirementsgenie epic anforderungen-im-chat)"
outcome-resolves: "DoD 1-8: Set/SetBody keep member lines (TestSetKeepsTheMembers), epic set exits 5/2 (TestEpicSetRewritesTheTextNotTheMembers), web edit + 409 (TestWebRewritesAnEpicAndRefusesAStaleEditor), merge clash block (TestMergeKeepsBothSidesOfOneField), text rendered and folded (headless Chromium), release note, race suite green"
review-summary: "none - the critique beside the work raised 7 findings over 4 pieces (retitle with unchanged body, shared edit() helper, one titleLine rule, Text() display form, ErrNoTitle, Fields json tags, help wrap); all applied, see notes 'Kritik Stueck 1-6'. Shape fits: edit() as Assign's twin, web keeps its own EpicFields/EpicEdit like web.Epic, features/edit-epic mirrors features/create-epic."
review-gaps: "removed: edit()/Set returned the epic as read alongside an error, no caller read it (CLI and web discard it) - now nil, doc trimmed, TestSetRefusesAStaleBase reads the ref instead. Searched for an existing comment stripper, blank-line trim and file-or-stdin helper: none to reuse (stdin read is the 2-copy pattern of internal/cli/sync.go:117, left). Left: Epics() computes Body()+Text() per epic on each minute reload - a few string ops, not worth caching."
---

# Titel und Beschreibung eines Epics lassen sich im Web und mit jaira epic set aendern

## Definition of Done

- [x] core/epic: eine Funktion schreibt Titel und/oder Beschreibung eines bestehenden Epics auf seinen Ref, Mitglieder bleiben unveraendert; Test
  proof: core/epic/epic.go Set(); TestSetKeepsTheMembers, TestSetRetitlesAndWritesOnlyWhatChanged
- [x] 'jaira epic set <name> --title <...> --description <...>' aendert eines oder beides; unbekanntes Epic Exit 5; ohne beide Flags Exit 2
  proof: internal/cli/epics.go newEpicSetCmd; TestEpicSetRewritesTheTextNotTheMembers
- [x] Web: eine Route (z.B. PATCH /api/epics/{name}) und in der UI ein Bearbeiten-Knopf am Epic, Titel und Beschreibung editierbar, Speichern zeigt sofort den neuen Stand
  proof: core/web/web.go handleSetEpic (PATCH /api/epics/{name}); internal/cli/web_source.go SetEpic; core/web/ui/src/features/edit-epic/ui/EditEpicForm.svelte; TestWebRewritesAnEpicAndRefusesAStaleEditor; headless Chromium: Edit->Save shows new title and text at once
- [x] Hat jemand das Epic zwischenzeitlich geaendert, geht keine der beiden Aenderungen still verloren (Test)
  proof: core/epic Merge()+clash(): TestMergeKeepsBothSidesOfOneField, TestMergeTellsARetitleFromATextEdit; Set Base: TestSetRefusesAStaleBase; web 409: TestWebRewritesAnEpicAndRefusesAStaleEditor
- [x] core/release/notes/<ID>.md; go test ./... -race gruen
  proof: core/release/notes/CRDH6B.md; go test -race -count=1 ./... RC=0; go vet ./... RC=0
- [x] Web zeigt den ganzen Markdown-Text eines Epics gerendert; 'jaira epic show' gibt ihn ebenfalls aus
  proof: EpicList.svelte {@html renderMarkdown(epic.text)}; internal/cli/epics.go epic show prints e.Text(); TestEpicSetRewritesTheTextNotTheMembers
- [x] Der Text ist standardmaessig eingeklappt (die description-Zeile bleibt sichtbar) und klappt per Klick auf
  proof: EpicList.svelte details.text without open, .desc outside it; headless Chromium: folded on load, opens on click
- [x] Beim Schreiben des Textes bleiben die Mitgliedszeilen des Epics unangetastet (Test)
  proof: core/epic/epic.go SetBody(); TestSetKeepsTheMembers

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/epic: Body() = alle Nicht-Mitgliedszeilen unter dem Frontmatter; SetBody(text) schreibt Frontmatter + Text + Mitgliedszeilen in alter Reihenfolge; Set(store, refs, name, Patch{Title,Description,Body *string, Base}) unter dem epics-Lock; Tests
- [x] core/epic Merge: Body wie title/description mergen; gleichzeitige Aenderung desselben Feldes nach Entscheidung (siehe Frage); Test
- [x] CLI: 'jaira epic set <name> --title --description --body --body-file <pfad|->' (Exit 5 unbekannt, Exit 2 ohne Flag); 'epic show' gibt den Text aus, --json mit body; Tests
- [x] Web: Epic-JSON bekommt body; PATCH /api/epics/{name} mit base (was der Editor gelesen hat) -> 409 wenn inzwischen geaendert; Source.SetEpic in internal/cli/web_source.go; Tests
- [x] UI: EpicList zeigt description-Zeile, Text gerendert (shared/lib/markdown) und eingeklappt; Bearbeiten-Knopf -> Formular Titel/description/Text; Speichern laedt Epics neu
- [x] core/release/notes/CRDH6B.md; go test ./... -race; UI-Build/Tests

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
- **2026-10-19 18:18 · Alexander Sacharov** — parent session: jaira-20
- **2026-10-19 18:19 · Alexander Sacharov** — Branch auf feat/tpckzb gestapelt per 'git merge --ff-only feat/tpckzb' (reset --keep wurde vom Auto-Mode verweigert; ff ist gleichwertig, weil HEAD=master Vorfahre von feat/tpckzb war). Stand 0f1b2e4d.
- **2026-10-19 18:22 · Alexander Sacharov** — Schritt 1 fertig (Set/Body/SetBody). Body = Text unter Frontmatter ohne Mitgliedszeilen, Leerzeilen am Rand getrimmt; SetBody schreibt Front + Leerzeile + Text + Leerzeile + Mitgliedszeilen -> New()-Layout bleibt byte-gleich. Set schreibt nichts, wenn Fields() gleich bleibt (auch wenn SetBody Mitgliedszeilen nach unten sortiert haette). Leerer --title wird ignoriert (Titel kann nicht leer werden). UI-Hinweis fuer spaeter: renderMarkdown escaped rohes HTML -> der '<!-- One ticket id per line -->'-Kommentar jedes Epics wuerde als Text erscheinen; fuer den Epic-Text HTML-Kommentare vor dem Rendern entfernen. Deshalb wird der Konfliktblock (Schritt 2) neben dem Kommentar-Marker auch eine sichtbare fette Zeile tragen.
- **2026-10-19 18:24 · Alexander Sacharov** — Kritik (daneben) Stueck 1 umgesetzt: (1) Set behandelt einen Body, der dem aktuellen gleicht (nach cleanBody), als unveraendert -> retitle laeuft auch, wenn das Web-Formular alle drei Felder schickt; Test in TestSetRetitlesAndWritesOnlyWhatChanged. (2) Set und Assign teilen jetzt edit(s, r, name, change) (Lock, Load, RecordEpic nur bei Aenderung); Create bleibt eigen, weil es 'nicht vorhanden' prueft. Hinweis fuer Web: epic.ErrChanged auf das bestehende web.ErrConflict (409) abbilden, keinen zweiten 409-Pfad.
- **2026-10-19 18:33 · Alexander Sacharov** — decision (Alex via Teamlead, 19.10.2026): Schritt 1 (DoD 1+8) angenommen. Konfliktmarker englisch: '<!-- jaira: changed at the same time -->' plus sichtbare fette Zeile - ersetzt den deutschen Wortlaut der Entscheidung von 06:48.
- **2026-10-19 18:34 · Alexander Sacharov** — Merge: Texte werden ueber plain() verglichen - eine Ueberschrift '# <eigener Titel>' am Textanfang wird durch titleMark ersetzt. Sonst wuerde ein Umbenennen (das die Ueberschrift mit aendert) auf einem Klon und ein neuer Absatz auf dem anderen als Konflikt gelten. Erster Entwurf mit Sonderfaellen fuer 'ours hat nur umbenannt' war unlesbar, deshalb verworfen. Konfliktblock: '<!-- jaira: changed at the same time -->' + fette sichtbare Zeile + fremder Wert + End-Kommentar; leerer fremder Wert erscheint als '(empty)'.
- **2026-10-19 18:36 · Alexander Sacharov** — Kritik Stueck 2: retitle und plain entschieden 'ist die erste Ueberschrift der Titel?' verschieden (oneLine nur auf einer Seite). Jetzt ein Helfer titleLine(), beide nutzen ihn; TestTheTitleHeadingIsFoundOneWay.
- **2026-10-19 19:29 · Alexander Sacharov** — decision (Alex via Teamlead, 19.10.2026): Schritt 2 (DoD 4, Merge-Konfliktblock inkl. titleLine-Fix) angenommen; weiter mit Schritt 3 (CLI).
- **2026-10-19 19:31 · Alexander Sacharov** — Kritik Stueck 3: (1) Epic.Text() = Body ohne die Titel-Ueberschrift (titleLine) und ohne <!-- -->-Kommentare; 'epic show' und spaeter die gerenderte Web-Ansicht nutzen Text, --json/Formular/Base bleiben bei Body (exakter Round-Trip). Ein frisch angelegtes Epic hat damit keinen Text. (2) leerer --title: epic.ErrNoTitle aus Set (gilt auch fuer Web), epicFail -> Exit 2. Vorher meldete die CLI faelschlich 'Nothing to do'.
- **2026-10-19 20:20 · Alexander Sacharov** — decision (Alex via Teamlead, 19.10.2026): Schritt 3 (CLI, inkl. Text()/ErrNoTitle-Fixes) angenommen. Schritte 4-6 ohne Zwischenstopp: am Ende EIN Gesamt-Diff und EINE Commit-Zeile fuer alles, was dann noch uncommittet ist. critique/optimize/testing/review laufen wie vorgesehen; mode bleibt conversational (Alex committet selbst).
- **2026-10-19 20:22 · Alexander Sacharov** — Set: Text wird vor dem Titel geschrieben und retitle laeuft immer (titleLine prueft gegen den alten Titel). Grund: das Web-Formular schickt Titel + geaenderten Text; vorher blieb dann '# Alter Titel' im Text stehen. Web: EpicEdit mit Zeigern + base; ErrChanged -> web.ErrConflict (409, 'reload'), ErrNoTitle -> ErrInvalid (400). Epic-JSON hat body (roh, fuers Formular) und text (Epic.Text, zum Rendern).
- **2026-10-19 20:26 · Alexander Sacharov** — UI live geprueft (headless Chromium via uv run --with playwright, Scratch-Board mit bare remote unter dem Scratchpad): Text eingeklappt, description sichtbar, Klick rendert ## Goal/**link**/Liste; Edit->Save zeigt neuen Titel und Absatz sofort; Formular offen + PATCH von aussen -> Save zeigt die 409-Meldung im Formular. Formular schickt nur geaenderte Felder + base, damit eine fremde Aenderung an einem nicht angefassten Feld keinen 409 ausloest. svelte-check: nur der bekannte vite.config.ts-Fehler (node:url). dist neu gebaut mit npm ci && npm run build. 'jaira web --help' nennt jetzt Edit.
- **2026-10-19 20:27 · Alexander Sacharov** — Kritik Stueck 4-6: json-Tags an epic.Fields entfernt (kein Leser, core/web hat eigene EpicFields); 'jaira web --help'-Absatz auf ~80 Spalten neu umbrochen und nennt das eingeklappte Text-Feld. Daneben-Kritik damit fertig, alle Funde umgesetzt.
- **2026-10-19 20:29 · Alexander Sacharov** — critique-Lane (Worker, Diff gelesen ohne dist): kein offener Fund. Losgelassen, weil es dem optimize-Lane gehoert: edit()/Set geben bei einem Fehler das gelesene Epic zurueck, kein Aufrufer liest es.
