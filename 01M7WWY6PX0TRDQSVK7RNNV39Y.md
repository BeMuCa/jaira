---
id: 01M7WWY6PX0TRDQSVK7RNNV39Y
title: Ein privates Board steht in .git/info/exclude statt in der .gitignore des Arbeitsbaums
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: Ein privates Board (vor jaira share) wird ueber .git/info/exclude ausgeschlossen; jaira schreibt fuer ein privates Board nie in die getrackte .gitignore.
context: |-
  Heute: board.AddIgnore (core/board/gitignore.go:44, aufgerufen aus core/board/board.go:22) haengt "# jaira board ... /.jaira/" an die .gitignore des Arbeitsbaums.
  Folge: in jedem neuen Worktree ist .gitignore geaendert. Ein Push-Guard im Projekt lehnt das ab, jeder Push scheitert, der Teamlead raeumt es bei jedem Ticket von Hand weg.
  Gemeldet vom Teamlead eines Projekts mit privatem Board am 20.11.2026; Alex: "das muss unterstuetzt werden, nicht umgangen".
  .git/info/exclude liegt im gemeinsamen git-dir (git rev-parse --git-common-dir) und gilt fuer alle Worktrees, landet nie in einem Diff.
  Offen zu klaeren im Plan: Migration bestehender Boards (Zeile aus .gitignore entfernen, in exclude setzen) und was board.Ignored/jaira share/RemoveIgnore dann lesen.
definition-of-done: "Ein neues privates Board und jeder neue Worktree davon lassen git status sauber; der Ausschluss steht in <git-common-dir>/info/exclude (Test)"
tags:
  - cli
blocked-by: []
related:
  - 01M263QP1A3MRANYPWGYSDBCT5
  - 01M4TRC37SNBER05X090DPSFGE
commits: []
created-at: 2026-11-20T16:34:22Z
updated-at: 2026-11-20T17:14:34Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-46134
claimed-at: 2026-11-20T16:42:22Z
outcome-what: "Ein privates Board steht jetzt in <git-common-dir>/info/exclude statt in der getrackten .gitignore (core/board/gitignore.go); jaira update verschiebt eine alte .gitignore-Zeile dorthin und schliesst ein privates Board ohne Ausschluss dort aus; share/--undo arbeiten auf beiden Dateien; ticket.LockFile aus Store.lock herausgezogen."
outcome-why: "Eine von jaira geaenderte .gitignore erscheint in jedem neuen Worktree als Aenderung, ein Push-Guard lehnt dann jeden Push ab; info/exclude gilt fuer alle Worktrees und taucht in keinem Diff auf."
outcome-resolves: Neue Worktrees eines Repos mit privatem Board bleiben sauber; kein Wegraeumen der .gitignore-Zeile von Hand mehr.
review-summary: "Ein privates Board wird jetzt per 'git rev-parse --path-format=absolute --git-common-dir --show-prefix' in <git-common-dir>/info/exclude ausgeschlossen (Muster /<prefix>.jaira/), nicht mehr in der getrackten .gitignore; ohne git bleibt .gitignore der Ort. Schreiben unter Sperrdatei info/exclude.jaira.lock (ticket.LockFile, aus Store.lock herausgezogen) und per WriteAtomic, CRLF wird erhalten. share/share --undo heben den Ausschluss in beiden Dateien auf bzw. setzen ihn in info/exclude; /.jaira/lanes/ bleibt in .gitignore. jaira update verschiebt eine alte .gitignore-Zeile nach info/exclude bzw. setzt den Ausschluss fuer ein privates Board ohne Ausschluss - nie bei einem Board mit .jaira/.gitattributes oder getrackten .jaira-Dateien. Hilfetexte, README, COMMANDS.md, Release-Note angepasst."
round: 4
review-gaps: |-
  blocker: unusual input / git < 2.31 - gitPlace (core/board/gitignore.go) uebergibt --path-format=absolute; aelteres git kennt die Option nicht und gibt sie als erste Zeile aus (git rev-parse reicht unbekannte --Optionen durch, exit 0). parsePlace nimmt "--path-format=absolute" als common dir: jaira init legt ./--path-format=absolute/info/exclude im Repo an, info/exclude bleibt leer, git status zeigt ?? .jaira/ und ?? --path-format=absolute/ - das Board ist NICHT privat (nachgestellt mit einem git-Wrapper, der das alte Verhalten emuliert). Fix: erste Zeile nur annehmen, wenn filepath.IsAbs(filepath.FromSlash(common)) gilt, sonst "" (Fallback .gitignore) - plus Tabellenfall in TestParsePlaceReadsWhatGitPrints
  note: pre-existing: commonGitDir in core/ticket/store.go:878 hat dieselbe Schwaeche (liefert "--path-format=absolute\n.git" als Schluessel fuer stateKey) - gleiche IsAbs-Pruefung waere billig
  note: dropEntry entfernt in info/exclude JEDE Kommentarzeile mit privateMarker - bei zwei Boards (Wurzel + /sub/.jaira/) nimmt share im einen den Kommentar des anderen mit, dessen Musterzeile bleibt kommentarlos stehen (kosmetisch)
  note: ein Unterverzeichnis mit Glob-Zeichen im Namen (z.B. a[1]) ergibt ein Muster /a[1]/.jaira/, das git als Zeichenklasse liest - das Board waere nicht ausgeschlossen; selten, Escapen von []*? waere die Abhilfe
  note: Ignored() kostet einen git-Aufruf pro Kommando (bindDriverIfShared) - auf Windows eher 20-50 ms als die gemessenen 2,4 ms unter Linux; bewusst in critique gelassen
test-verdict: "pass: go test -race ./... gruen (RC=0) auf 394a544a; DoD 1-7 an den Tests geprueft; Binary: Board in git mit committeter /.jaira/-Zeile -> update rc 0, git status leer; vorige Runde: init/Schreibschutz/3x parallel wie notiert"
hunt-verdict: "pass: bughunt done at 394a544a, 0 findings (Folgerunde ueber den Fix-Diff); Runde 1 (44007130): 1 change, behoben in 394a544a"
review-verdict: "Zurueck nach in-progress: DoD 1-7 sind erfuellt und getestet (go test ./core/board ./core/ticket ./internal/cli gruen, Paralleltest 30/30), aber auf git < 2.31 (Ubuntu 20.04, Debian 11) macht init das Board nicht privat und legt ein Verzeichnis '--path-format=absolute/' im Repo an - falsches Ergebnis mit Datenpreisgabe, Fix ist eine Zeile."
review-round: 3
review-check: "1. cd $(mktemp -d) && git init -q && jaira init  2. git status --short zeigt nur AGENTS.md/CLAUDE.md, NICHT .jaira/  3. tail -2 .git/info/exclude zeigt den jaira-Kommentar und /.jaira/  4. git worktree add ../wt && cd ../wt && git status --short ist leer  5. jaira share -> info/exclude ohne /.jaira/, .gitignore hat /.jaira/lanes/  6. jaira share --undo -> /.jaira/ wieder in .git/info/exclude  7. altes git: go test ./core/board -run TestParsePlaceReadsWhatGitPrints mit dem neuen Fall \"--path-format=absolute\\n.git\\n\" -> exclude muss \"\" sein"
---

# Ein privates Board steht in .git/info/exclude statt in der .gitignore des Arbeitsbaums

## Definition of Done

- [x] Ein neues privates Board und jeder neue Worktree davon lassen git status sauber; der Ausschluss steht in <git-common-dir>/info/exclude (Test)
  proof: internal/cli/privateexclude_test.go TestANewPrivateBoardAndItsWorktreesLeaveGitStatusClean
- [x] jaira share hebt den Ausschluss in info/exclude auf und laesst /.jaira/lanes/ privat wie heute (Test)
  proof: internal/cli/privateexclude_test.go TestShareLiftsTheExcludeAndKeepsLanesPrivate
- [x] Ein Board, das die alte .gitignore-Zeile hat, wird bei jaira update umgestellt und die Zeile entfernt (Test)
  proof: internal/cli/privateexclude_test.go TestUpdateMovesTheOldGitignoreLineToExclude, TestUpdateLeavesABoardInGitAsItIs
- [x] Eingabeklassen als Tabellentest: info/exclude fehlt (auch info/ fehlt), leer, ohne Schluss-Newline, mit CRLF; Board-Wurzel mit Leerzeichen im Pfad; Board in einem Unterverzeichnis des Repos (Muster /sub/.jaira/) (Test)
  proof: core/board/exclude_test.go TestAddAndRemoveIgnoreOverEveryShapeOfExclude, TestABoardInASubdirectoryIsExcludedByItsOwnPath
- [x] Worktree vs. Haupt-Checkout: init/share/update aus einem Worktree schreiben in dasselbe gemeinsame info/exclude und nie in eine worktree-eigene Datei (Test)
  proof: internal/cli/privateexclude_test.go TestANewPrivateBoardAndItsWorktreesLeaveGitStatusClean (git-dir des Worktrees bleibt leer; share/--undo im Worktree)
- [x] Zwei jaira-Aufrufe gleichzeitig: AddIgnore/RemoveIgnore/MoveIgnore schreiben info/exclude unter einer Sperrdatei daneben (ticket.LockFile) und per Temp-Datei + Rename; nach parallelen AddIgnore steht der Eintrag genau einmal darin (Test)
  proof: core/board/exclude_test.go TestParallelAddIgnoreWritesTheEntryOnce (-count=30 gruen)
- [x] Windows: der Pfad von git (C:/... mit Schraegstrichen) wird zu einem gueltigen Dateipfad; CRLF-Zeilen werden erkannt und beim Entfernen bleibt der Rest der Datei byteweise erhalten (Test)
  proof: core/board/exclude_test.go TestParsePlaceReadsWhatGitPrints (Windows-Laufwerkspfad auf windows-latest, git<2.31-Ausgabe), TestAddAndRemoveIgnoreOverEveryShapeOfExclude/CRLF

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/board/gitignore.go: excludePath(root) ueber git rev-parse --git-common-dir; AddIgnore/RemoveIgnore/Ignored auf info/exclude, Fallback .gitignore ohne git
- [x] board-Tests: neues Board + neuer Worktree lassen git status sauber; Eintrag steht in info/exclude
- [x] share/share --undo: Ausgaben nennen die Datei; Test dass lanes-Zeile bleibt und exclude leer wird
- [x] jaira update: MigrateIgnore (alte .gitignore-Zeile -> exclude; privates Board ohne Ausschluss -> exclude) + Test
- [x] init-/share-/update-Hilfetexte, docs/COMMANDS.md, README, Release-Note core/release/notes/NNV39Y.md

## Progress
- **2026-11-20 16:42 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/35112.sock
- **2026-11-20 16:43 · Alexander Sacharov** — decision: keine offenen Entscheidungen fuer die Person; selbst entschieden: (1) Ziel ist <git-common-dir>/info/exclude, ermittelt per 'git rev-parse --path-format=absolute --git-common-dir' wie commonGitDir in core/ticket/store.go:878 — kein Nachbau der .git-Datei-Aufloesung. (2) Kein git-Repo (init ausserhalb von git): Fallback auf .gitignore wie bisher, dort ist sie nicht getrackt und haelt das Board privat, falls spaeter git init kommt. (3) Ignored() liest .gitignore UND info/exclude: eine von Hand in .gitignore geschriebene .jaira/-Zeile zaehlt weiter als privat. (4) share entfernt die Zeile aus beiden Dateien; /.jaira/lanes/ bleibt in der .gitignore, denn nach share ist das Board geteilt und die .gitignore wird ohnehin mit committet. (5) JSON-Key gitignore_written von init bleibt (Vertrag), bedeutet jetzt 'Ausschluss geschrieben'.
- **2026-11-20 16:43 · Alexander Sacharov** — decision Migration (jaira update): (a) .gitignore hat die /.jaira/-Zeile -> Zeile samt jaira-Kommentar raus, Eintrag in info/exclude; die geaenderte .gitignore ist eine einmalige Aenderung zum Committen, update sagt das. (b) Teamkollege, der diesen Commit pullt: seine .gitignore hat die Zeile nicht mehr, sein privates Board waere ungetrackt sichtbar. Deshalb setzt update den exclude-Eintrag auch, wenn .jaira/ nirgends ignoriert ist, keine .jaira/.gitattributes hat und git keine Datei unter .jaira/ trackt — das ist genau ein privates Board ohne Ausschluss. Ein geteiltes Board (Attribute oder getrackte Dateien) wird nie angefasst. Nur in 'jaira update', nicht im automatischen Refresh jedes Kommandos: der schreibt nichts an der Repo-Konfiguration.
- **2026-11-20 16:50 · Alexander Sacharov** — Interim-Regeln vom Teamlead (Alex 20.11.2026, bis TW5FE2 ausgeliefert): DoD nennt Fehlerklassen (leer, gross, Kodierung/CRLF, Leerzeichen-Pfade, Windows, zwei Prozesse, Fehler halbwegs); in-progress prueft Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte, Tabellentests, Geschwister-Stellen per grep; Funde nennen Klasse+Beispiel; Blocker nur falsches Ergebnis/Datenverlust/Absturz/Sicherheit; nach Fix-Durchgang Test+Hunt nur ueber den Diff seit letztem Hunt; 3 Rueckgaben -> Frage an Teamlead; kleine review-gaps sofort schliessen, Vorbestehendes als 'note: pre-existing:'.
- **2026-11-20 16:50 · Alexander Sacharov** — decision Nebenlaeufigkeit: keine Sperrdatei. info/exclude ist pro Repo gemeinsam, die Store-Sperren liegen pro Arbeitsbaum und wuerden zwei Worktrees nicht serialisieren. Stattdessen: lesen, Inhalt bauen, Temp-Datei im selben Verzeichnis, Rename. Zwei parallele AddIgnore bauen beide 'alt + ein Eintrag', der letzte Rename gewinnt -> genau ein Eintrag, nie eine halbe Datei. Verlorenes Update nur, wenn ein Dritter zeitgleich eine andere Zeile schreibt; das nehmen wir hin.
- **2026-11-20 16:55 · Alexander Sacharov** — Korrektur zur decision Nebenlaeufigkeit: Rename allein reicht NICHT — TestParallelAddIgnoreWritesTheEntryOnce fiel in 2 von 15 Laeufen mit doppeltem Eintrag (B liest nach As Rename, prueft aber vor As Schreiben). Jetzt Sperrdatei <info/exclude>.jaira.lock; die Logik von Store.lock ist dafuer als ticket.LockFile(path) herausgezogen (gleiche Stale-/Windows-Behandlung), nicht nachgebaut. 30 von 30 Laeufen gruen. Ausserdem gefunden: CRLF-exclude bekam nach Add+Remove ein einzelnes \r am Ende — gefixt in dropEntry.
- **2026-11-20 16:58 · Alexander Sacharov** — critique: (1) IgnoreFile in IgnoreFileName falten, exportierte Funktion mit einem Aufrufer. (2) share.go Long-Text neu umbrechen. Stehen gelassen: Ignored() kostet auf geteilten Boards einen git-Aufruf mehr pro Kommando (bindDriverIfShared) — gemessen 2,4 ms unter Linux; die Alternative, .git-Datei/commondir selbst aufzuloesen, lehnt core/ticket/store.go:870 aus gutem Grund ab. Der zweite git-Aufruf neben commonGitDir (store.go:878) ist gewollt: gitPlace braucht zusaetzlich --show-prefix, in einem Aufruf.
- **2026-11-20 17:00 · Alexander Sacharov** — critique Runde 2: nichts mehr. Steht auf der vorigen critique-Note (Kosten von Ignored, zweiter git-Aufruf). Bewusst gelassen: MoveIgnore nimmt auch eine von Hand geschriebene .jaira/-Zeile aus .gitignore — gleiche Wirkung, Ziel des Tickets.
- **2026-11-20 17:02 · Alexander Sacharov** — bughunt: change: internal/cli/update.go:344 - moveIgnore moves a .gitignore entry before it checks whether the board is shared (the hasAttributes/Tracks guard is only on line 352), so 'jaira update' on a board whose .jaira/ files are committed but whose tracked .gitignore lists /.jaira/ (the 'clone of one whose ignore line was committed' case tickets.go names) rewrites the team's tracked .gitignore and tells the user to commit it, although update's help (update.go:185) says such a board is never touched; once committed, teammates lose the ignore and their private tickets show up as untracked; ran: git init, jaira init, jaira share --undo, append /.jaira/ to .gitignore, git add -f .jaira .gitignore && commit, jaira update (exit 0, prints 'Moved this private board's entry ... commit the changed .gitignore once', git status shows ' M .gitignore')
- **2026-11-20 17:03 · Alexander Sacharov** — bughunt: done at 44007130 - 1 finding (1 change, 0 pre-existing) in 6 modules, first hunt; concurrency, error paths, input, permissions checked; files: README.md core/board/announce.go core/board/board.go core/board/exclude_test.go core/board/gitignore.go core/release/notes/NNV39Y.md core/ticket/store.go docs/COMMANDS.md internal/cli/privateexclude_test.go internal/cli/share.go internal/cli/sharestate_test.go internal/cli/tickets.go internal/cli/update.go
- **2026-11-20 17:04 · Alexander Sacharov** — testing: Klasse des bughunt-Funds: jeder Zweig von moveIgnore (update.go) muss hinter derselben Privat-Pruefung stehen (keine .jaira/.gitattributes, nichts unter .jaira/ getrackt). Heute prueft nur der 'added'-Zweig sie; der 'moved'-Zweig aendert die committete .gitignore auch auf einem Board, das in git ist — nach dem Commit waere es fuer Teamkollegen nicht mehr ignoriert. Fix: Pruefung an den Anfang von moveIgnore; Test fuer beide Zweige ueber die Klasse (Attribute / getrackte Datei).
- **2026-11-20 17:04 · Alexander Sacharov** — bughunt-Fund behoben an der Klasse: moveIgnore prueft 'in git' (Attribute oder getrackte .jaira-Datei) jetzt einmal am Anfang fuer beide Zweige. TestUpdateLeavesABoardInGitAsItIs deckt shared clone, shared clone mit committeter Zeile und getrackte Tickets ohne Attribute ab, je mit und ohne --dry-run; faellt auf dem alten Code in 2 von 3 Faellen.
- **2026-11-20 17:05 · Alexander Sacharov** — critique Runde 3 (Diff 44007130..394a544a): nichts — die Privat-Pruefung steht jetzt einmal vor beiden Zweigen, kein neues Muster.
- **2026-11-20 17:06 · Alexander Sacharov** — bughunt: done at 394a544a - 0 findings in 2 modules, follow-up hunt (diff since 44007130: internal/cli/update.go, internal/cli/privateexclude_test.go); concurrency, error paths, input, permissions checked; files: README.md core/board/announce.go core/board/board.go core/board/exclude_test.go core/board/gitignore.go core/release/notes/NNV39Y.md core/ticket/store.go docs/COMMANDS.md internal/cli/privateexclude_test.go internal/cli/share.go internal/cli/sharestate_test.go internal/cli/tickets.go internal/cli/update.go
- **2026-11-20 17:11 · Alexander Sacharov** — review: blocker - auf git < 2.31 reicht rev-parse die unbekannte Option --path-format=absolute als erste Ausgabezeile durch; parsePlace nimmt sie als git-common-dir, init schreibt ./--path-format=absolute/info/exclude ins Repo und das Board bleibt sichtbar (nachgestellt mit git-Wrapper). Fix: in parsePlace nur einen absoluten Pfad als common dir annehmen, sonst Fallback; Tabellenfall dazu. Rest der review-gaps sind notes.
- **2026-11-20 17:11 · Alexander Sacharov** — review Runde 1 zurueck: (1) Blocker git<2.31: parsePlace nimmt eine relative erste Zeile ('--path-format=absolute') als common dir -> Fix an der Klasse: nur absolute Pfade annehmen, sonst kein git-Ort (Fallback .gitignore); dieselbe Pruefung in commonGitDir (core/ticket/store.go) mitziehen, gleiche Ursache. (2) dropEntry nimmt jeden jaira-Kommentar mit -> nur den Kommentar direkt ueber einer entfernten Zeile. (3) Glob-Zeichen im Unterverzeichnis-Praefix escapen. Gelassen: Kosten von Ignored auf Windows (critique-Note).
