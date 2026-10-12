---
id: 01M4QECQMKHK5BKWRZMFAR97K9
title: "Release-Notes als eine Datei je Aenderung, die Version setzt der Release-Build"
status: testing
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
updated-at: 2026-10-12T07:37:30Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-37358
claimed-at: 2026-10-12T07:23:55Z
outcome-what: "wie zuvor; WriteIndex baut jetzt eine Gruppenliste je Tag, Kommentare zum alten Release-Cut nachgezogen"
outcome-why: "critique 1"
outcome-resolves: "critique-1-Befunde (a) und (b)"
review-summary: "Runde 2: keine Befunde. critique-1 (a) und (b) erledigt in dffae51. Stehen gelassen wie in Note 'critique 1' begruendet. Bewusst nicht behandelt: Vorab-Tags (v0.4.0-rc1) - das Repo hat keine, und die Ordnung nach version:refname waere dann die Frage des ersten solchen Tags."
review-gaps: "Entfernt: Fehlerzweig beim Lesen einer eingebetteten Notiz, die fs.Glob gerade gefunden hat (kann nicht fehlschlagen); Duplikat-Pruefung assigned[name] im Index (WriteIndex schreibt jeden Namen genau einmal). Gelassen: changeLines neben parseNotes - drei Zeilen, ein gemeinsamer Helfer waere laenger; eigener git()-Helfer (gitref.run ist an Ticket-Refs gebunden); Test-Helfer run() wie git(t,...) in jedem anderen Testpaket; testdata/NOTES-before-AR97K9.md (DoD 4); Pruefung changes[name] im Index (eine nach dem Release geloeschte Notiz kann dort stehen). Kosten: Notes() laeuft einmal je jaira update, nichts zu heben."
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
