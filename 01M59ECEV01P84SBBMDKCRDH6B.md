---
id: 01M59ECEV01P84SBBMDKCRDH6B
title: Titel und Beschreibung eines Epics lassen sich im Web und mit jaira epic set aendern
status: signoff
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
updated-at: 2026-10-19T20:36:02Z
updated-by: Alexander Sacharov
mode: conversational
claimed-by: DESKTOP-RFTCH11-52791
claimed-at: 2026-10-19T20:32:26Z
outcome-what: "core/epic: Body/Text/SetBody/Set with Patch+Base, Merge keeps both sides of a field in a marked block; CLI 'jaira epic set' and body in 'epic show'/--json; web PATCH /api/epics/{name} (409 on stale base) with an Edit form and the text rendered folded under the description"
outcome-why: "An epic could not be changed after it was created, and its markdown text below the frontmatter was invisible in the web (requirementsgenie epic anforderungen-im-chat)"
outcome-resolves: "DoD 1-8: Set/SetBody keep member lines (TestSetKeepsTheMembers), epic set exits 5/2 (TestEpicSetRewritesTheTextNotTheMembers), web edit + 409 (TestWebRewritesAnEpicAndRefusesAStaleEditor), merge clash block (TestMergeKeepsBothSidesOfOneField), text rendered and folded (headless Chromium), release note, race suite green"
review-summary: "core/epic gets Body() (text under the frontmatter without member lines), Text() (Body without the title heading and HTML comments, for display), SetBody() (new text, member lines appended verbatim in order) and Set(store, refs, name, Patch) - one locked write through a new edit() helper that Assign now shares too; Patch.Base refuses with ErrChanged when a field being changed differs from what the editor read, an empty title gives ErrNoTitle, a heading '# <old title>' is renamed along. Merge now also merges the text; where both clones changed the same field (title, description, text) ours wins and theirs is appended to the text in a '<!-- jaira: changed at the same time -->' block with a visible bold line; a retitle is not counted as a text edit (titleMark). CLI: 'jaira epic set <name> --title/--description/--body/--body-file <path|->' (exit 5 unknown, 2 without flags or with empty title), 'epic show' prints Text(), --json carries body. Web: PATCH /api/epics/{name} with base (409 via web.ErrConflict, 400 on empty title), Epic JSON has body+text, EpicList shows the text folded under the description rendered via renderMarkdown, an Edit button opens features/edit-epic/EditEpicForm which sends only changed fields plus base. dist rebuilt, release note CRDH6B.md with 4 lines."
review-gaps: "Not blocking, both narrow: (1) a line in the text that looks like a member line ('- <full 26-char ULID> ...') is taken by Body()/FromBytes as a member - writing such a list item via 'epic set --body' or the web form silently adds that ticket to the epic and the line vanishes from the text (core/epic/epic.go Body + FromBytes:176). Short handles are not affected (ParseMember needs a full ULID). (2) After a 409 the form stays open with the old base, so every further Save answers 409 again; the only way out is Cancel + Edit, which drops what was typed - the person has to copy their text out first (EditEpicForm.svelte base is fixed at open). Claims in outcome-resolves are backed by the diff; go test ./core/epic ./core/web ./internal/cli re-run green here. Nothing committed yet: commits=[], all work is in the worktree on top of feat/tpckzb."
test-verdict: "pass: go vet ./... RC=0, go test -race -count=1 ./... RC=0 (31 packages ok), npm ci && npm run build leaves dist identical; DoD 1-8 checked against their tests; scratch board with bare remote: epic set exits 5 (unknown) / 2 (no flag, empty title), member line kept on the ref; headless Chromium: text folded under the visible description, opens rendered (h2, bold, list), Edit->Save shows new title+paragraph at once, save after a change behind the open form shows the 409 message in the form"
review-verdict: "The diff meets all 8 DoD items; no defect that breaks a DoD item found. Two small gaps above, neither needs to block acceptance - (1) could become a follow-up if epic texts ever list full ticket ids. Still to do outside this lane: commit the worktree (conversational mode, Alex commits) and open the PR only after TPCKZB's."
review-round: 1
review-check: "1. cd /home/alex/projects/.worktrees/jaira-crdh6b && go build -o /tmp/jaira-crdh6b ./cmd/jaira  2. /tmp/jaira-crdh6b epic ls - pick any epic name, call it E  3. /tmp/jaira-crdh6b epic set E --body $'## Goal\\nfirst **bold** line' - prints 'Epic \"E\" rewritten on refs/jaira/epics/E'  4. /tmp/jaira-crdh6b epic show E - you see '## Goal' and the bold line under the description, and the same tickets as before  5. /tmp/jaira-crdh6b epic set E - exits 2 with 'nothing to set'; /tmp/jaira-crdh6b epic set nosuch --title x - exits 5  6. /tmp/jaira-crdh6b web, open the page - under E's description a folded 'Text'; click it - '## Goal' renders as a heading, the line in bold  7. click 'Edit', change the title, Save - the new title shows at once  8. click 'Edit' again, leave the form open; in a terminal run /tmp/jaira-crdh6b epic set E --title other; back in the form change the title, Save - a red message 'was changed since you opened it' appears and nothing is overwritten  9. merge conflict block: go test ./core/epic -run TestMergeKeepsBothSidesOfOneField -v (no UI path, needs two clones)"
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
- **2026-10-19 20:31 · Alexander Sacharov** — testing: alles gruen, siehe test-verdict. Einziger Konsolen-Eintrag im Browser ist der erwartete 409.
- **2026-10-19 20:36 · Alexander Sacharov** — Abnahme-Seite: https://claude.ai/artifact/FTJQKnj5LLAFS1RzzGfMbm (4 Personen-Schritte statt 9 review-check-Schritten, 3 Maschinen-Checks gruen). Wegwerf-Board dafuer unter ~/.claude/acceptance/jaira/feat-crdh6b/board, Binary daneben. Review-Luecken offen fuer Alex: (1) eine Textzeile '- <volle 26-stellige ULID>' wird als Mitglied gelesen; (2) nach einem 409 bleibt das Formular auf dem alten base, nur Cancel+Edit hilft, Getipptes geht dabei verloren.
