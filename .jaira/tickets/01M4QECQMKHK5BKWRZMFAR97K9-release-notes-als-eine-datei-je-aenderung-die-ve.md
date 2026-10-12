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
commits:
  - 06094d3bb92243c16d3732f3e3a33795c1f1109e
  - dffae51a85b55d4290a471d32ab6e476057d9f90
  - 1fcf5a7051cc561ef6c6f64f3d799e49d666baf8
  - 1abf2eb227e43be9caf4c0b1a607b74a4bf7f63f
created-at: 2026-10-12T06:56:51Z
updated-at: 2026-10-12T08:33:38Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-89604
claimed-at: 2026-10-12T08:27:32Z
outcome-what: "Variante A: Release-Version aus dem Ordnernamen core/release/notes/<version>/; Knopf release-notes.yaml verschiebt per release.Cut und oeffnet PR release/<version>; release-tag.yaml taggt den Merge und startet release.yaml; Versionsindex und Regel <TICKET>-2.md entfernt"
outcome-why: "Alex hat nach dem Merge-Vorschlag zu #55 Variante A gewaehlt: einfacher als ein Index beim Build, Version steht sichtbar im Repo"
outcome-resolves: "DoD 6-9"
review-summary: "Release-Notizen liegen jetzt je Ticket als eigene Datei core/release/notes/<TICKET>.md (Zeilen '- ...'); die 5 bisherigen Unreleased-Zeilen sind dorthin umgezogen, NOTES.md ist geschlossene Historie bis 0.3.6. Welche Version eine Datei bekommt, steht nirgends im Repo: der goreleaser-before-hook 'go generate ./core/release' ruft release.WriteIndex, das je Tag v* (nur --merged HEAD, aeltester zuerst) per git ls-tree die notes/ liest und die Datei dem aeltesten Tag zuordnet, der sie enthaelt; das Ergebnis geht in das gitignorte notes/versions, das mit eingebettet wird. Ohne Index (go build) zeigt composeNotes jede Datei unter Unreleased. Letzter Commit 1abf2eb ergaenzt in CLAUDE.md und AGENTS.md die von Alex gewaehlte Regel: sobald ein Tag die Datei haelt, geht eine weitere Aenderung in <TICKET>-2.md."
review-gaps: "Gap 1 der vorigen Review ist mit 1abf2eb geschlossen, wie Alex es gewaehlt hat (Option 1, nur Regeltext, keine Pruefung). Offen und bewusst stehen gelassen: (a) die Regel '<TICKET>-2.md' erzwingt nichts - eine Zeile, die doch an eine getaggte Datei angehaengt wird, landet still unter der alten Version; (b) ein lokal per go generate erzeugtes notes/versions wird in spaetere go builds eingebettet und veraltet still (nur Entwicklerrechner). Nicht zum Ticket gehoerend: das untracked .jaira/milestones/ im worktree_diff (drei demo-Milestones vom 15.09.) - nicht mit committen. Sonst keine Defekte gefunden; release.yaml checkt mit fetch-depth: 0 aus, also sieht WriteIndex die Tags."
test-verdict: "pass: Runde 3 nur Doku (1abf2eb); go test -count=1 -race ./core/release ./internal/wintrap RC=0; der Satz steht in CLAUDE.md und AGENTS.md; DoD-Proofs CLAUDE.md:199/AGENTS.md:152 zeigen weiter auf die Abschnittsueberschrift; Funktion seit Runde 2 unveraendert (Tag-Build v0.3.7 im Scratch-Klon gepruef)"
review-verdict: "Der Diff erfuellt DoD 1-5; go test ./core/release -race und go vet gruen. Keine neuen Defekte. Die Restluecken (a) und (b) sind bekannt und von Alex nicht beauftragt - bereit fuer die Abnahme."
review-check: "1. go test ./core/release -race -run 'TestTwoBranches|TestWriteIndex|TestComposeNotesWithoutIndex|TestNotesUnchanged' -v  -> vier PASS. 2. go build -o /tmp/jaira ./cmd/jaira  3. ls core/release/notes  -> 8WFNQY.md GDB5YJ.md K1W2JR.md MG677R.md, kein versions. 4. /tmp/jaira update --json auf einem Board mit Stempel 0.3.6  -> genau ein Eintrag 'Unreleased' mit 5 Zeilen. 5. Lies CLAUDE.md Abschnitt 'Client-facing changes go in core/release/notes/'  -> der Satz zu <TICKET>-2.md steht dort, wortgleich in AGENTS.md."
---

# Release-Notes als eine Datei je Aenderung, die Version setzt der Release-Build

## Definition of Done

- [x] Eine Aenderung traegt ihre Release-Notiz als eigene Datei unter core/release/ ein; zwei Branches, die je eine anlegen, mergen ohne Konflikt - ein Test oder ein Skriptlauf mit zwei Branches belegt es
  proof: core/release/notes_test.go TestTwoBranchesAddingNotesMergeWithoutConflict (Gegenprobe: zwei Appends an eine Datei konfliktieren)
- [x] Ein Tag v<x.y.z> baut ein Binary, dessen 'jaira update' die neuen Eintraege unter '<x.y.z>' zeigt, ohne dass jemand NOTES.md vor dem Tag umbenannt hat - geprueft an einer Fixture-Repo mit zwei Tags
  proof: Variante A: Version aus Ordnername notes/<version>/, kein Umbenennen in NOTES.md - core/release/notes_test.go TestCutMovesEveryNoteIntoTheVersionFolder; lokal: cut 0.3.7 + go build, jaira update --json bei Stempel 0.3.6 -> genau '0.3.7' mit 7 Zeilen. Der echte Tag-Build auf GitHub ist bis nach dem Merge ungeprueft
- [x] Kein Release-Schritt schreibt auf master; ein Build ohne Pipeline (go build ./cmd/jaira) laeuft weiter und zeigt die neuen Eintraege als Unreleased
  proof: Kein Workflow pusht auf master: release-notes.yaml pusht nur release/<version>, release-tag.yaml nur den Tag; go build ./cmd/jaira zeigt Dateien direkt in notes/ als Unreleased - core/release/notes_test.go TestComposeNotesReadsVersionFolders, TestNotesUnchangedByMovingUnreleasedIntoFiles
- [x] Die bestehenden Eintraege (Unreleased und alle Versionen) erscheinen nach dem Umbau in jaira update unveraendert - Test gegen den alten Stand
  proof: core/release/notes_test.go TestNotesUnchangedByMovingUnreleasedIntoFiles gegen core/release/testdata/NOTES-before-AR97K9.md
- [x] CLAUDE.md, AGENTS.md, der Agent-Block (falls er die Regel enthaelt) und der Format-Kommentar in NOTES.md sagen den neuen Weg; go test ./... -race gruen
  proof: CLAUDE.md:211 und :236, AGENTS.md:164 und :189 wortgleich (Abschnitt 'Client-facing changes go in core/release/notes/'); core/release/NOTES.md:4 Format-Kommentar; Agent-Block (core/board/announce.go) enthaelt die Regel nicht
- [x] Eine GitHub Action mit workflow_dispatch (Eingabe version) verschiebt alle Dateien direkt unter core/release/notes/ nach core/release/notes/<version>/, legt Branch release/<version> an und oeffnet den PR; nichts geht direkt auf master
  proof: .github/workflows/release-notes.yaml (workflow_dispatch input version, nur master; go run ./core/release/cut, Branch release/<version>, gh pr create); lokal geprueft: cut + commit in Kopie verschiebt 6 Dateien nach notes/0.3.7/, zweiter Lauf exit 1. UNGEPRUEFT bis nach Merge: Lauf auf GitHub, Repo-Einstellung 'Allow GitHub Actions to create pull requests'
- [x] Der Merge eines release/<version>-PR setzt per Action den Tag v<version> auf den Merge-Commit, und release.yaml baut daraus wie bisher
  proof: .github/workflows/release-tag.yaml (pull_request closed+merged, head release/X.Y.Z aus eigener Repo -> git tag -a v<version> merge_commit_sha, push, gh workflow run release.yaml --ref v<version>); release.yaml hat dafuer workflow_dispatch. UNGEPRUEFT bis nach Merge
- [x] Das Binary nimmt die Version aus dem Ordnernamen notes/<version>/, Dateien direkt in notes/ zeigt es als Unreleased; go generate, notes/versions und die Regel <TICKET>-2.md sind entfernt; Tests decken Ordner, Unreleased und eine Probe des Verschiebeskripts ab
  proof: core/release/release.go composeNotes/entryOf/newer; core/release/cut.go Cut; Tests core/release/notes_test.go TestComposeNotesReadsVersionFolders, TestComposeNotesWithoutNotesShowsOnlyHistory, TestCutMovesEveryNoteIntoTheVersionFolder, TestCutRefuses; index.go, gen/, notes/versions (.gitignore, .goreleaser.yaml) und Regel <TICKET>-2.md entfernt
- [x] CLAUDE.md, AGENTS.md und der Format-Kommentar in NOTES.md beschreiben den Release-Ablauf (Knopf, PR, Merge setzt Tag); go test ./... -race gruen
  proof: CLAUDE.md:236 / AGENTS.md:189 'Releasing is a button' (Knopf release-notes, PR release/<version>, Merge -> release-tag setzt v<version>, startet release.yaml); NOTES.md:4; go vet ./... und go test ./... -race exit 0

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
- **2026-10-12 08:25 · Alexander Sacharov** — Alex 12.10. nach dem Merge-Vorschlag zu #55: Variante A statt Versionsindex beim Build. Neue Notizen liegen in core/release/notes/<ID>.md; ein Release verschiebt alle Dateien aus notes/ nach notes/<version>/. Ausloeser: GitHub Action workflow_dispatch mit Eingabe version -> Branch release/<version> mit dem Verschieben, PR oeffnen. Merge eines release/*-PR -> zweite Action setzt den Tag v<version> auf den Merge-Commit -> release.yaml baut wie bisher. Nichts landet direkt auf master. Der Versionsindex (go generate, notes/versions, git tag --contains) entfaellt; das Binary liest die Version aus dem Ordnernamen; Dateien direkt in notes/ sind Unreleased. Damit entfaellt auch die Regel <TICKET>-2.md. PR #55 bleibt offen und waechst um diese Aenderung.
- **2026-10-12 08:27 · Alexander Sacharov** — Variante A, selbst entschieden (keine offene Frage an Alex): (1) Verschieben als Go-Funktion release.Cut + core/release/cut/main.go statt Shell - so prueft go test das Skript direkt, plattformunabhaengig; ersetzt core/release/gen. (2) Eingabe version: fuehrendes v wird abgeschnitten, sonst genau X.Y.Z aus Ziffern; Cut verweigert, wenn notes/<version>/ schon existiert oder keine Datei direkt in notes/ liegt (leerer Release-PR waere ohne Commit). (3) Reihenfolge: Versionsordner absteigend nach X.Y.Z numerisch, darunter NOTES.md (bis 0.3.6); innerhalb eines Ordners nach Dateiname. (4) Tag-Action: ein Tag-Push mit GITHUB_TOKEN loest release.yaml NICHT aus (GitHub startet auf GITHUB_TOKEN-Events keine Workflows, ausser workflow_dispatch) - darum bekommt release.yaml zusaetzlich workflow_dispatch und die Tag-Action ruft 'gh workflow run release.yaml --ref v<version>'. (5) Release-Action laeuft nur auf master (if github.ref == refs/heads/master). Ungeprueft bis nach dem Merge: Repo-Einstellung 'Allow GitHub Actions to create and approve pull requests' muss an sein; ein von GITHUB_TOKEN geoeffneter PR startet ci.yaml nicht - verlangt Branch-Protection gruene CI, haengt der Release-PR.
- **2026-10-12 08:33 · Alexander Sacharov** — in-progress Variante A: Bug selbst gefunden und behoben - echo "version=$(go run ...)" verschluckt den Exit-Code von cut, der naechste Schritt liefe mit leerer Version; jetzt Zuweisung in eigener Zeile, unter bash -eo pipefail (GitHub-Standard) lokal geprueft: Ablehnung stoppt den Schritt. Lokale Probe: Kopie des Baums, cut 0.3.7 + Commit = 6 Renames nach notes/0.3.7/; jaira update bei Stempel 0.3.6 (HOME im scratchpad, der Stempel liegt unter ~/.jaira/state/) zeigt genau 0.3.7 mit 7 Zeilen. Keine eigene Notiz-Datei fuer AR97K9: Nutzer des Binaries merken nichts.
