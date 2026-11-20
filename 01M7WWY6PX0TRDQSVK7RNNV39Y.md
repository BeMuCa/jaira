---
id: 01M7WWY6PX0TRDQSVK7RNNV39Y
title: Ein privates Board steht in .git/info/exclude statt in der .gitignore des Arbeitsbaums
status: optimize
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
updated-at: 2026-11-20T17:00:42Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-46134
claimed-at: 2026-11-20T16:42:22Z
outcome-what: "Ein privates Board steht jetzt in <git-common-dir>/info/exclude statt in der getrackten .gitignore (core/board/gitignore.go); jaira update verschiebt eine alte .gitignore-Zeile dorthin und schliesst ein privates Board ohne Ausschluss dort aus; share/--undo arbeiten auf beiden Dateien; ticket.LockFile aus Store.lock herausgezogen."
outcome-why: "Eine von jaira geaenderte .gitignore erscheint in jedem neuen Worktree als Aenderung, ein Push-Guard lehnt dann jeden Push ab; info/exclude gilt fuer alle Worktrees und taucht in keinem Diff auf."
outcome-resolves: Neue Worktrees eines Repos mit privatem Board bleiben sauber; kein Wegraeumen der .gitignore-Zeile von Hand mehr.
review-summary: "none"
round: 2
---

# Ein privates Board steht in .git/info/exclude statt in der .gitignore des Arbeitsbaums

## Definition of Done

- [x] Ein neues privates Board und jeder neue Worktree davon lassen git status sauber; der Ausschluss steht in <git-common-dir>/info/exclude (Test)
  proof: internal/cli/privateexclude_test.go TestANewPrivateBoardAndItsWorktreesLeaveGitStatusClean
- [x] jaira share hebt den Ausschluss in info/exclude auf und laesst /.jaira/lanes/ privat wie heute (Test)
  proof: internal/cli/privateexclude_test.go TestShareLiftsTheExcludeAndKeepsLanesPrivate
- [x] Ein Board, das die alte .gitignore-Zeile hat, wird bei jaira update umgestellt und die Zeile entfernt (Test)
  proof: internal/cli/privateexclude_test.go TestUpdateMovesTheOldGitignoreLineToExclude, TestUpdateLeavesASharedBoardUnexcluded
- [x] Eingabeklassen als Tabellentest: info/exclude fehlt (auch info/ fehlt), leer, ohne Schluss-Newline, mit CRLF; Board-Wurzel mit Leerzeichen im Pfad; Board in einem Unterverzeichnis des Repos (Muster /sub/.jaira/) (Test)
  proof: core/board/exclude_test.go TestAddAndRemoveIgnoreOverEveryShapeOfExclude, TestABoardInASubdirectoryIsExcludedByItsOwnPath
- [x] Worktree vs. Haupt-Checkout: init/share/update aus einem Worktree schreiben in dasselbe gemeinsame info/exclude und nie in eine worktree-eigene Datei (Test)
  proof: internal/cli/privateexclude_test.go TestANewPrivateBoardAndItsWorktreesLeaveGitStatusClean (git-dir des Worktrees bleibt leer; share/--undo im Worktree)
- [x] Zwei jaira-Aufrufe gleichzeitig: AddIgnore/RemoveIgnore/MoveIgnore schreiben info/exclude unter einer Sperrdatei daneben (ticket.LockFile) und per Temp-Datei + Rename; nach parallelen AddIgnore steht der Eintrag genau einmal darin (Test)
  proof: core/board/exclude_test.go TestParallelAddIgnoreWritesTheEntryOnce (-count=30 gruen)
- [x] Windows: der Pfad von git (C:/... mit Schraegstrichen) wird zu einem gueltigen Dateipfad; CRLF-Zeilen werden erkannt und beim Entfernen bleibt der Rest der Datei byteweise erhalten (Test)
  proof: core/board/exclude_test.go TestParsePlaceReadsWhatGitPrints, TestAddAndRemoveIgnoreOverEveryShapeOfExclude/CRLF

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
