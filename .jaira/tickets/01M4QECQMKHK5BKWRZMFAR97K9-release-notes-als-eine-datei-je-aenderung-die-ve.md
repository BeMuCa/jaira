---
id: 01M4QECQMKHK5BKWRZMFAR97K9
title: "Release-Notes als eine Datei je Aenderung, die Version setzt der Release-Build"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Zwei PRs geraten nie mehr ueber core/release/NOTES.md in Konflikt, und ein Release braucht keinen eigenen Commit mehr, um Unreleased in die Versionsnummer umzubenennen"
context: |-
  Was heute falsch ist: jede Branch haengt ihre Zeile ans Ende von '## Unreleased' in core/release/NOTES.md. Zwei PRs aendern damit dieselbe Stelle - fast jeder PR nach einem anderen hat einen Konflikt. Zuletzt PR BeMuCa/jaira#52 (K1W2JR) am 12.10.2026, vorher 8WFNQY, TECZH5, GDB5YJ (alle mit 'chore(<id>): merge master, keep every release note under Unreleased').
  Dazu kommt ein eigener Release-PR nur fuers Umbenennen (#44 'chore(release): cut 0.3.6', 2 Zeilen in NOTES.md).
  Alex' Wunsch (12.10.): das soll automatisch gehen - die Release-Pipeline traegt die Eintraege in die Version ein, die gerade getaggt wurde.

  Was schon bekannt ist:
  - NOTES.md ist per go:embed im Binary (core/release/release.go:17). jaira update waehlt per Ueberschrift: selectSince sucht die gestempelte Version als '## <version>'. Ein Binary 0.3.7 MUSS '## 0.3.7' in sich tragen - umbenennen NACH dem Bauen reicht nicht.
  - Release-Pipeline: .github/workflows/release.yaml, Tag 'v*' -> goreleaser; .goreleaser.yaml hat before.hooks (go mod tidy, go test ./...). Dort kann ein Generator laufen, bevor gebaut wird. Checkout hat fetch-depth: 0, Tags sind da.
  - Ein Commit von CI zurueck auf master verstoesst gegen 'nichts landet direkt auf master' (CLAUDE.md) - darum soll nichts zurueckgeschrieben werden.

  Vorgeschlagener Weg (im Plan pruefen, nicht gesetzt):
  1. Jede Aenderung legt core/release/notes/<ticket-id>.md an, Inhalt: eine oder mehrere Zeilen '- ...'. Eindeutige Dateinamen -> kein Konflikt.
  2. Ein Generator (go generate oder Skript im goreleaser-before-hook) ordnet jede Datei der fruehesten Version zu, deren Tag den Commit enthaelt, der sie angelegt hat (git log --diff-filter=A + git tag --contains). Ohne solchen Tag: Unreleased. Er schreibt die Datei, die eingebettet wird, neuester Release zuerst.
  3. In der Repo wird nie etwas verschoben; die Version ergibt sich aus git.
  Offen: was ein Build ohne Generator einbettet (go build, go install von Hand) - z. B. alles als Unreleased. Und ob die bestehenden Abschnitte in NOTES.md (0.3.6 und aelter) als Dateien migriert werden oder als feste Vergangenheit stehen bleiben.

  Danach muessen die Regeln nachziehen: Abschnitt 'Client-facing changes go in core/release/NOTES.md' in CLAUDE.md und AGENTS.md (pruefen, ob er aus dem Agent-Block in core/board/announce.go kommt oder handgeschrieben ist) und der Format-Kommentar oben in NOTES.md.
definition-of-done: "Eine Aenderung traegt ihre Release-Notiz als eigene Datei unter core/release/ ein; zwei Branches, die je eine anlegen, mergen ohne Konflikt - ein Test oder ein Skriptlauf mit zwei Branches belegt es"
tags:
  - release
  - ci
blocked-by: []
related: []
commits: []
created-at: 2026-10-12T06:56:51Z
updated-at: 2026-10-12T08:04:56Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-69610
claimed-at: 2026-10-12T08:04:33Z
outcome-what: "Regel ergaenzt: nach dem Tag kommt eine weitere Aenderung in <TICKET>-2.md, nie in die getaggte Datei"
outcome-why: "review gap 1, Alex Option 1: der Index ordnet nach Dateiname zu, eine angehaengte Zeile landete unter der alten Version"
outcome-resolves: "review gap 1"
review-summary: "Release notes since 0.3.6 are one file per ticket under core/release/notes/<TICKET>.md ('- ' lines, no heading); NOTES.md keeps 0.3.6 and older as closed history. Notes() = note files without index entry as Unreleased + releases from the embedded notes/versions + NOTES.md. The goreleaser before-hook runs 'go generate ./core/release' (core/release/gen -> release.WriteIndex): for every tag v* merged into HEAD, oldest first, git ls-tree of core/release/notes; a file belongs to the first tag holding it; written to the gitignored notes/versions, nothing committed. go build without generate shows all notes as Unreleased. The 5 former Unreleased lines moved into MG677R/K1W2JR/8WFNQY/GDB5YJ.md; CLAUDE.md/AGENTS.md describe the new way; .gitattributes pins notes/*.md to LF. Release CI clones with fetch-depth 0, so the tags are there."
review-gaps: "1) A line added to a note file AFTER its ticket's file was tagged is filed under the old release: WriteIndex assigns by file name (core/release/index.go, seen[name]), not by line. A board stamped with that release never sees the line in jaira update, and the old release section of newer binaries shows a change that release never had. CLAUDE.md/AGENTS.md even say 'A ticket with a second change adds a second line to its own file' without limiting it to before the tag, and the diff dropped the old rule 'a released note that grows after the fact describes a binary nobody has'. Fix: one sentence in CLAUDE.md/AGENTS.md: once a tag holds the file, a further change goes into a new file (e.g. <TICKET>-2.md), never into the old one; optionally a test. 2) Minor: a notes/versions left over from a local go generate is embedded into every later local go build and goes stale silently - only a developer machine is affected, not the release. DoD 1-5 are otherwise met by the diff; go test ./core/release -race green."
test-verdict: "pass: go build/vet + go test -count=1 -race ./... RC=0; DoD 1-5 im Baum geprueft (vier benannte Tests gruen, Proof-Zeilen CLAUDE.md:199, AGENTS.md:152, NOTES.md:1 stimmen); Funktion: Scratch-Klon mit lokalem Tag v0.3.7 -> go generate schreibt Index, git status leer, Binary mit -X main.version=0.3.7 zeigt bei Stempel 0.3.6 genau '0.3.7' mit 5 Zeilen und danach 'Nothing has changed'; go build ohne Generator zeigt dieselben 5 als Unreleased"
review-verdict: "DoD met, one gap: the docs allow appending to a released note file, which the index then files under the old version (gap 1). Small doc fix, not a blocker; the person decides whether it is done before signoff or afterwards."
review-check: "1. cd into the worktree, run: go test ./core/release -race -v -run 'TestTwoBranches|TestWriteIndex|TestComposeNotesWithoutIndex|TestNotesUnchanged' - all four PASS. 2. go build -o /tmp/jaira ./cmd/jaira, then in a scratch board run /tmp/jaira update --json on a board stamped 0.3.6 - exactly one 'Unreleased' section with 5 lines (MG677R x2, K1W2JR, 8WFNQY, GDB5YJ). 3. cd core/release && go generate . && cat notes/versions - only the comment line (no tag holds notes/ yet); git status shows nothing new (the file is gitignored). Then rm notes/versions. 4. Read CLAUDE.md section 'Client-facing changes go in core/release/notes/' and decide gap 1: the sentence 'A ticket with a second change adds a second line to its own file' should be limited to before the tag."
---

# Release-Notes als eine Datei je Aenderung, die Version setzt der Release-Build

## Definition of Done

- [x] Eine Aenderung traegt ihre Release-Notiz als eigene Datei unter core/release/ ein; zwei Branches, die je eine anlegen, mergen ohne Konflikt - ein Test oder ein Skriptlauf mit zwei Branches belegt es
  proof: core/release/index_test.go TestTwoBranchesAddingNotesMergeWithoutConflict (Gegenprobe: zwei Appends an eine Datei konfliktieren)
- [x] Ein Tag v<x.y.z> baut ein Binary, dessen 'jaira update' die neuen Eintraege unter '<x.y.z>' zeigt, ohne dass jemand NOTES.md vor dem Tag umbenannt hat - geprueft an einer Fixture-Repo mit zwei Tags
  proof: core/release/index_test.go TestWriteIndexFilesNotesUnderTheTagThatFirstHeldThem (Fixture mit v0.1.0 und annotiertem v0.2.0; Nachbau von v0.1.0 sieht 0.2.0 nicht); Generator core/release/index.go WriteIndex, Hook .goreleaser.yaml
- [x] Kein Release-Schritt schreibt auf master; ein Build ohne Pipeline (go build ./cmd/jaira) laeuft weiter und zeigt die neuen Eintraege als Unreleased
  proof: go build ./cmd/jaira + jaira update --json bei Stempel 0.3.6 -> nur 'Unreleased' mit den 5 Notiz-Dateien; core/release/index_test.go TestComposeNotesWithoutIndexShowsEveryNoteUnreleased; .goreleaser.yaml go generate schreibt nur gitignortes core/release/notes/versions, kein push
- [x] Die bestehenden Eintraege (Unreleased und alle Versionen) erscheinen nach dem Umbau in jaira update unveraendert - Test gegen den alten Stand
  proof: core/release/index_test.go TestNotesUnchangedByMovingUnreleasedIntoFiles gegen core/release/testdata/NOTES-before-AR97K9.md
- [x] CLAUDE.md, AGENTS.md, der Agent-Block (falls er die Regel enthaelt) und der Format-Kommentar in NOTES.md sagen den neuen Weg; go test ./... -race gruen
  proof: CLAUDE.md:199 und AGENTS.md:152 Abschnitt 'Client-facing changes go in core/release/notes/'; Format-Kommentar core/release/NOTES.md:1; Agent-Block (core/board/announce.go) enthaelt die Regel nicht; go test ./... -race exit 0

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/release/notes/ anlegen; die 5 Unreleased-Zeilen aus NOTES.md als <ticket>.md dorthin, '## Unreleased' aus NOTES.md entfernen; alter Stand als testdata fuer DoD 4
- [x] release.go: //go:embed NOTES.md notes; Notes() = Unreleased (Dateien ohne Indexeintrag) + Versionen aus notes/versions (Reihenfolge der Indexdatei) + feste Abschnitte aus NOTES.md
- [x] release.WriteIndex(repo): git tag --merged HEAD 'v*' sortiert, je Tag git ls-tree der notes/, aeltester Tag gewinnt; schreibt notes/versions (gitignort); go:generate-Programm core/release/gen ruft es
- [x] Tests: Fixture-Repo zwei Branches mergen konfliktfrei (DoD 1); Fixture-Repo mit zwei Tags -> Eintraege unter beiden Versionen (DoD 2); ohne Index alles Unreleased (DoD 3); neue Notes() == alter Stand (DoD 4)
- [x] .goreleaser.yaml before-hook 'go generate ./core/release'; .gitignore notes/versions
- [x] CLAUDE.md, AGENTS.md (Abschnitt Client-facing) und Format-Kommentar in NOTES.md auf den neuen Weg; go test ./... -race

## Progress
- **2026-10-12 07:23 · Alexander Sacharov** — parent session: jaira-4c
- **2026-10-12 07:25 · Alexander Sacharov** — Offene Entscheidungen vor dem Plan: keine, die nur Alex treffen kann. Selbst entschieden: (1) Dateiname core/release/notes/<ticket-id>.md, Inhalt Zeilen '- ...' - eindeutig je Ticket, also konfliktfrei. (2) Zuordnung zur Version macht ein Generator im goreleaser-before-hook: Datei gehoert zum aeltesten Tag v*, in dessen Baum (git ls-tree) sie liegt - ein Aufruf je Tag statt git log+tag --contains je Datei, gleiches Ergebnis bei linearer Tag-Folge. Er schreibt eine gitignorte Indexdatei in core/release/notes/, die mit eingebettet wird; nichts wird committet. (3) Build ohne Generator (go build): Index fehlt -> jede Notiz-Datei steht unter Unreleased (DoD 3). (4) Versionierte Abschnitte 0.3.6 und aelter bleiben als feste Vergangenheit in NOTES.md; die 5 heutigen Unreleased-Zeilen (MG677R x2, K1W2JR, 8WFNQY, GDB5YJ) wandern als Dateien nach notes/, sonst landeten sie im naechsten Release-Binary unter 'Unreleased' statt '0.3.7'. (5) Reihenfolge innerhalb einer Version: nach Dateiname - ein Build ohne git kennt keine andere. Abschnitt 'Client-facing changes' ist in CLAUDE.md und AGENTS.md handgeschrieben (hinter jaira:local), nicht aus core/board/announce.go.
- **2026-10-12 07:26 · Alexander Sacharov** — Plan-Begruendung: Indexdatei statt NOTES.md im CI umzuschreiben, weil dann ein Build ohne Generator dieselbe Quelle liest und nur die Versionszuordnung fehlt. Logik in package release (nicht in gen/main), damit die DoD-2-Fixture im selben Paket Generator und Parser zusammen pruefen kann. 'git tag --merged HEAD', damit ein Nachbau eines alten Tags keine Versionen spaeterer Tags vergibt. DoD 4 vergleicht Unreleased als Menge: die Reihenfolge innerhalb einer Version folgt jetzt dem Dateinamen.
- **2026-10-12 07:32 · Alexander Sacharov** — Gelernt: git ls-tree <tag>:<pfad> filtert aus einem Unterverzeichnis heraus trotzdem nach dem cwd-Pfad und liefert nichts - darum --full-tree in WriteIndex (go generate laeuft in core/release). Keine eigene Notiz-Datei fuer AR97K9: jaira update zeigt dasselbe wie vorher, die Aenderung betrifft nur Beitragende (CLAUDE.md/AGENTS.md). Lokal geprueft: go build + jaira update --json bei Stempel 0.3.6 zeigt genau Unreleased mit 5 Zeilen. goreleaser prueft 'dirty' vor den before-hooks, und notes/versions ist gitignort - der Hook macht den Baum nicht schmutzig.
- **2026-10-12 07:36 · Alexander Sacharov** — critique 1: (a) WriteIndex auf Gruppenliste {version, names} umbauen statt first-map + versions-slice. (b) veraltete Kommentare zum 'leeren Unreleased nach Release-Cut' in release.go (sinceEntries) und release_test.go nachziehen. Stehen gelassen mit Begruendung: eigener git()-Helfer, WriteIndex im Paket release, testdata-Kopie.
- **2026-10-12 07:42 · Alexander Sacharov** — review: DoD 1-5 erfuellt, go test ./core/release -race gruen. Luecke 1: WriteIndex ordnet nach Dateiname zu - eine Zeile, die nach dem Tag an eine bereits getaggte Notiz-Datei angehaengt wird, landet unter der alten Version, und Boards mit diesem Stempel sehen sie nie. CLAUDE.md/AGENTS.md erlauben genau das ('second line to its own file') und die alte Regel 'released note that grows after the fact' ist entfallen. Vorschlag: ein Satz - nach dem Tag eine neue Datei <TICKET>-2.md. Luecke 2 (klein): lokal liegengebliebenes notes/versions wird in spaetere go builds eingebettet.
- **2026-10-12 07:43 · Alexander Sacharov** — Frage an Alex (review gap 1): Eine Zeile, die nach dem Tag an eine schon getaggte Notiz-Datei angehaengt wird, landet unter der alten Version. Optionen: 1) (Empfohlen) Vor signoff einen Satz in CLAUDE.md/AGENTS.md: sobald ein Tag die Datei enthaelt, geht eine weitere Aenderung in eine neue Datei <TICKET>-2.md - ~5 Min, Ticket geht zurueck nach in-progress. 2) Wie es ist abnehmen, Satz spaeter als eigenes Ticket. 3) Zusaetzlich zu 1 einen Test/Pruefung, die das Anhaengen an getaggte Dateien in CI ablehnt - ~30 Min.
- **2026-10-12 08:02 · Alexander Sacharov** — Alex 12.10. zu review gap 1: Option 1 - ein Satz in CLAUDE.md/AGENTS.md: sobald ein Tag die Notiz-Datei enthaelt, kommt eine weitere Aenderung in eine neue Datei <TICKET>-2.md. Ticket zurueck nach in-progress, danach wieder signoff.
- **2026-10-12 08:04 · Alexander Sacharov** — Gap 1 umgesetzt nach Alex' Option 1: Satz in CLAUDE.md und AGENTS.md (Abschnitt Client-facing), keine Code-Aenderung, kein Test (Option 3 nicht gewaehlt). Der Abschnitt kommt nicht aus core/board/announce.go - nichts zu generieren. Review-Gap 2 (lokal liegengebliebenes notes/versions veraltet still) bewusst nicht angefasst: nur Entwicklerrechner, nicht beauftragt.
