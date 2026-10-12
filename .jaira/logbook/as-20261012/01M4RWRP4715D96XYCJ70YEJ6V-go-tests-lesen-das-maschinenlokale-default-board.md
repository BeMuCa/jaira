---
id: 01M4RWRP4715D96XYCJ70YEJ6V
title: Go-Tests lesen das maschinenlokale Default-Board aus ~/.jaira
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "go test ./... -race ist auf jeder Maschine gruen, egal was in ~/.jaira liegt"
context: |-
  16 Tests in 'go test ./... -race' sind rot auf dieser Maschine (Alex, WSL).
  Sie sind identisch rot auf der Basis 43080ed7 (master) - kein Ticket hat sie eingefuehrt.
  Gefunden vom Worker von 8W1R94 am 2026-10-12 im testing-Lane.
  Ursache laut Worker: das maschinenlokale Default-Board in ~/.jaira (settings.json bzw. Default-Board) leckt in die Tests - die Tests lesen das echte HOME statt eines isolierten.
  Erster Schritt: go test ./... -race 2>&1 | grep -- '--- FAIL' fuer die Liste der 16; dann pruefen, welcher Code ~/.jaira liest (os.UserHomeDir / settings-Pfad) und warum der Test-Helper HOME nicht auf ein t.TempDir() setzt.
  Folge heute: jede Testing-Lane muss 16 rote Tests von Hand als 'vorbestehend' wegerklaeren.
definition-of-done: "Alle Tests isolieren HOME (bzw. den ~/.jaira-Pfad) auf ein Temp-Verzeichnis; go test ./... -race ist gruen mit einem befuellten ~/.jaira auf der Maschine; ein Test belegt, dass ein Default-Board in ~/.jaira das Ergebnis nicht aendert"
tags:
  - ci
blocked-by: []
related: []
commits:
  - d2577cb0dddd88ea6749b5cc986905e99e246cc1
created-at: 2026-10-12T20:27:18Z
updated-at: 2026-10-12T21:03:49Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-92214
claimed-at: 2026-10-12T20:27:43Z
outcome-what: "Auf upstream/master neu aufgesetzt: DefaultBoardPath und UserLanesDir folgen JAIRA_HOME ueber inUserHome (core/lane/lane.go); Beweistest core/lane/jairahome_test.go; Release-Note core/release/notes/0YEJ6V.md. Testisolation kommt von AK3Z6D."
outcome-why: "Upstream isoliert die Tests schon (AK3Z6D); offen war nur, dass JAIRA_HOME Default-Board und Lane-Katalog nicht verlegte - die einzigen Per-User-Pfade ohne es."
outcome-resolves: Default-Board und Lane-Katalog unter verlegtem JAIRA_HOME; Branch ohne Konflikte zu upstream
review-summary: "Ein Commit (d2577cb0) auf upstream/master, 0 Commits dahinter. core/lane/lane.go bekommt inUserHome(name): $JAIRA_HOME/<name> wenn gesetzt, sonst ~/.jaira/<name>. DefaultBoardPath und UserLanesDir nutzen es nach ihren eigenen Overrides (JAIRA_DEFAULT_BOARD, JAIRA_LANES_DIR), damit folgen Default-Board und Lane-Katalog JAIRA_HOME wie settings, identity, projects, store und update-cache schon. Dadurch greift die upstream-Testisolation (internal/testhome.Run setzt JAIRA_HOME) jetzt auch fuer diese zwei Pfade. Neu: Beweistest core/lane/jairahome_test.go und Release-Note core/release/notes/0YEJ6V.md (eine Zeile, Format korrekt)."
review-gaps: "none: Diff deckt das DoD; Isolation aller Test-Pakete kommt von AK3Z6D, das auf diesem Branch enthalten ist. Geprueft und unkritisch: (1) SaveDefaultBoard schreibt unter gesetztem JAIRA_HOME jetzt dorthin statt nach ~/.jaira - beabsichtigt und in der Release-Note genannt; (2) letzte UserHomeDir-Nutzer ohne JAIRA_HOME sind role/target.go (~/.claude/skills, kein jaira-Zustand) und selfupdate/install.go - ausserhalb des Tickets. Worktree-Rest .jaira/milestones/ ist nicht von diesem Ticket."
test-verdict: "pass: go vet + go test ./... -race -count=1 RC=0 auf feat/0yej6v-2 (29 Pakete ok, 0 FAIL) mit befuelltem ~/.jaira; TestJairaHomeHidesTheRealDefaultBoard rot ohne inUserHome, gruen mit; Binary: jaira init mit leerem JAIRA_HOME -> Built-ins inkl. backlog, ohne JAIRA_HOME -> Lanes aus ~/.jaira/default-board.md wie bisher"
review-verdict: "Erfuellt das DoD, keine Defekte gefunden, mergebar ohne Konflikte. go test ./... -race -count=1 RC=0 auf feat/0yej6v-2 mit befuelltem ~/.jaira (selbst nachgeprueft). Nutzersichtbare Aenderung: wer JAIRA_HOME setzt, bekommt Default-Board und Lanes von dort - das ist der einzige Punkt, den eine Person bewusst akzeptieren sollte."
review-check: "1. cd /home/alex/projects/.worktrees/jaira-0yej6v  2. go test ./core/lane -run TestJairaHomeHidesTheRealDefaultBoard -count=1 -v  -> PASS  3. go test ./... -race -count=1 2>&1 | grep -c FAIL  -> 0 (mit deinem echten ~/.jaira/default-board.md)  4. go build -o /tmp/j ./cmd/jaira && d=$(mktemp -d) && git -C $d init -q && JAIRA_HOME=$(mktemp -d) /tmp/j -C $d init && /tmp/j -C $d lanes  -> Built-in-Lanes inkl. backlog, nicht die aus ~/.jaira/default-board.md  5. git log --oneline HEAD..upstream/master | wc -l  -> 0 (keine Konflikte)"
---

# Go-Tests lesen das maschinenlokale Default-Board aus ~/.jaira

## Definition of Done

- [x] Alle Tests isolieren HOME (bzw. den ~/.jaira-Pfad) auf ein Temp-Verzeichnis; go test ./... -race ist gruen mit einem befuellten ~/.jaira auf der Maschine; ein Test belegt, dass ein Default-Board in ~/.jaira das Ergebnis nicht aendert
  proof: core/lane/jairahome_test.go:TestJairaHomeHidesTheRealDefaultBoard (rot ohne inUserHome, gruen mit); Isolation aller Test-Pakete: internal/testhome.Run (AK3Z6D, upstream); go test ./... -race -count=1 RC=0 auf feat/0yej6v-2 mit befuelltem ~/.jaira

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-12 20:27 · Alexander Sacharov** — parent session: jaira-f5
- **2026-10-12 20:30 · Alexander Sacharov** — Keine offenen Entscheidungen. Ursache gefunden: core/lane/defaultboard.go DefaultBoardPath und core/lane/lane.go UserLanesDir ignorieren JAIRA_HOME (settings, project, identity, ticket/store, selfupdate/cache beachten es). Deshalb rot sogar in internal/cli, dessen TestMain JAIRA_HOME schon auf ein Temp-Verzeichnis setzt: ~/.jaira/default-board.md (ohne backlog) gibt neuen Test-Boards 9 Lanes ohne backlog. Entscheidung: (1) beide Funktionen fallen nach ihrem eigenen Override auf $JAIRA_HOME zurueck, vor ~/.jaira - eine Stellschraube fuer allen Per-User-Zustand; (2) TestMain mit JAIRA_HOME=Temp in jedem Test-Paket, das Per-User-Zustand lesen kann (statt HOME umzubiegen: HOME umbiegen nimmt git die ~/.gitconfig und damit user.name). Verworfen: HOME in Tests umbiegen; nur JAIRA_DEFAULT_BOARD in Tests setzen (UserLanesDir leckt dann weiter).
- **2026-10-12 20:35 · Alexander Sacharov** — critique: nichts. Fix im Callee (DefaultBoardPath/UserLanesDir) statt in Tests ist richtig - identity.go:92 behauptete schon, JAIRA_HOME verlege allen Per-User-Zustand. testhome.Main verallgemeinert das bestehende internal/cli/main_test.go statt ein neues Muster daneben zu stellen. Bewusst stehen gelassen: TestMain auch in Paketen, die heute nicht rot waren (gitrepo, tag, validate ...) - die DoD verlangt alle Tests isoliert, und sie erreichen core/lane bzw. settings transitiv.
- **2026-10-12 20:46 · Alexander Sacharov** — review: Duplikat von upstream AK3Z6D (151643cd). Entscheidung fuer den Menschen in signoff: 1) (empfohlen) zurueck nach in-progress: Branch auf upstream/master neu aufsetzen, nur inUserHome + jairahome_test.go + core/release/notes/0YEJ6V.md behalten, testhome-Aenderungen verwerfen — kleiner Diff, keine Konflikte. 2) Schliessen als Duplikat — kostet nichts, aber JAIRA_HOME verlegt dann weiter nicht Default-Board und Lane-Katalog. 3) Branch so mergen und Konflikte beim Merge loesen — nicht empfohlen, zwei testhome-Varianten muessten von Hand zusammengefuehrt werden.
- **2026-10-12 20:47 · Alexander Sacharov** — worker: Review hat in signoff abgeliefert, Befund: upstream/master hat den Testteil schon mit AK3Z6D (151643cd, internal/testhome.Run) gruendlicher geloest; Branch basiert auf 43080ed7, 137 Commits hinter upstream, 21 Konflikte. Nur neu: DefaultBoardPath/UserLanesDir folgen JAIRA_HOME (Nutzerverhalten) + Beweistest. Entscheidung fuer die Person in signoff: 1) (empfohlen) zurueck nach in-progress, auf upstream/master neu aufsetzen, nur Produktionsteil + Beweistest + core/release/notes/0YEJ6V.md behalten; 2) Ticket schliessen, wenn JAIRA_HOME fuer Default-Board/Lanes nicht gewollt ist.
- **2026-10-12 20:50 · Alexander Sacharov** — decision: Alex 2026-10-12 (signoff): Option 1 - zurueck nach in-progress, auf upstream/master neu aufsetzen; nur JAIRA_HOME fuer DefaultBoardPath/UserLanesDir, den Beweistest und core/release/notes/0YEJ6V.md behalten, den Testteil von AK3Z6D (internal/testhome) nutzen statt eigenem.
- **2026-10-12 20:50 · Alexander Sacharov** — worker: Neuaufbau begonnen. Patches gesichert (prod.patch: defaultboard.go/lane.go, jairahome_test.go). Schritt: Branch feat/0yej6v auf upstream/master (9f6e9877) setzen, Ticket per jaira pull zurueckholen, Patch anwenden, Notes-Datei core/release/notes/0YEJ6V.md.
- **2026-10-12 20:51 · Alexander Sacharov** — worker: 'git reset --hard upstream/master' vom Permission-Classifier abgelehnt (Irreversible Local Destruction). Nicht umgangen; Frage an jaira-f5. Patches liegen im Scratchpad des Workers; Branch unveraendert auf 4883c757, Ticket-Datei modifiziert, in-progress.
- **2026-10-12 20:52 · Alexander Sacharov** — worker: neu aufgebaut auf Branch feat/0yej6v-2 (von upstream/master 9f6e9877, Alex Option 2). feat/0yej6v bleibt liegen (letzter Commit parkt nur den Ticketstand) und wird spaeter geloescht. Auf dem neuen Branch nur: inUserHome in core/lane, core/lane/jairahome_test.go, core/release/notes/0YEJ6V.md; Testisolation kommt von AK3Z6D internal/testhome.Run.
- **2026-10-12 20:56 · Alexander Sacharov** — critique (Runde 2): nur die Funde der Review gelesen. Duplikat testhome -> verworfen, upstream internal/testhome.Run genutzt; fehlendes Unset von JAIRA_DEFAULT_BOARD/JAIRA_LANES_DIR -> erledigt durch upstream Run; Release-Note -> jetzt core/release/notes/0YEJ6V.md. Nichts Neues.
