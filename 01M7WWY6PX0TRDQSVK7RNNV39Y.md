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
updated-at: 2026-11-20T16:50:02Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-46134
claimed-at: 2026-11-20T16:42:22Z
---

# Ein privates Board steht in .git/info/exclude statt in der .gitignore des Arbeitsbaums

## Definition of Done

- [ ] Ein neues privates Board und jeder neue Worktree davon lassen git status sauber; der Ausschluss steht in <git-common-dir>/info/exclude (Test)
- [ ] jaira share hebt den Ausschluss in info/exclude auf und laesst /.jaira/lanes/ privat wie heute (Test)
- [ ] Ein Board, das die alte .gitignore-Zeile hat, wird bei jaira update umgestellt und die Zeile entfernt (Test)

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [~] core/board/gitignore.go: excludePath(root) ueber git rev-parse --git-common-dir; AddIgnore/RemoveIgnore/Ignored auf info/exclude, Fallback .gitignore ohne git
- [ ] board-Tests: neues Board + neuer Worktree lassen git status sauber; Eintrag steht in info/exclude
- [ ] share/share --undo: Ausgaben nennen die Datei; Test dass lanes-Zeile bleibt und exclude leer wird
- [ ] jaira update: MigrateIgnore (alte .gitignore-Zeile -> exclude; privates Board ohne Ausschluss -> exclude) + Test
- [ ] init-/share-/update-Hilfetexte, docs/COMMANDS.md, README, Release-Note core/release/notes/NNV39Y.md

## Progress
- **2026-11-20 16:42 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/35112.sock
- **2026-11-20 16:43 · Alexander Sacharov** — decision: keine offenen Entscheidungen fuer die Person; selbst entschieden: (1) Ziel ist <git-common-dir>/info/exclude, ermittelt per 'git rev-parse --path-format=absolute --git-common-dir' wie commonGitDir in core/ticket/store.go:878 — kein Nachbau der .git-Datei-Aufloesung. (2) Kein git-Repo (init ausserhalb von git): Fallback auf .gitignore wie bisher, dort ist sie nicht getrackt und haelt das Board privat, falls spaeter git init kommt. (3) Ignored() liest .gitignore UND info/exclude: eine von Hand in .gitignore geschriebene .jaira/-Zeile zaehlt weiter als privat. (4) share entfernt die Zeile aus beiden Dateien; /.jaira/lanes/ bleibt in der .gitignore, denn nach share ist das Board geteilt und die .gitignore wird ohnehin mit committet. (5) JSON-Key gitignore_written von init bleibt (Vertrag), bedeutet jetzt 'Ausschluss geschrieben'.
- **2026-11-20 16:43 · Alexander Sacharov** — decision Migration (jaira update): (a) .gitignore hat die /.jaira/-Zeile -> Zeile samt jaira-Kommentar raus, Eintrag in info/exclude; die geaenderte .gitignore ist eine einmalige Aenderung zum Committen, update sagt das. (b) Teamkollege, der diesen Commit pullt: seine .gitignore hat die Zeile nicht mehr, sein privates Board waere ungetrackt sichtbar. Deshalb setzt update den exclude-Eintrag auch, wenn .jaira/ nirgends ignoriert ist, keine .jaira/.gitattributes hat und git keine Datei unter .jaira/ trackt — das ist genau ein privates Board ohne Ausschluss. Ein geteiltes Board (Attribute oder getrackte Dateien) wird nie angefasst. Nur in 'jaira update', nicht im automatischen Refresh jedes Kommandos: der schreibt nichts an der Repo-Konfiguration.
- **2026-11-20 16:50 · Alexander Sacharov** — Interim-Regeln vom Teamlead (Alex 20.11.2026, bis TW5FE2 ausgeliefert): DoD nennt Fehlerklassen (leer, gross, Kodierung/CRLF, Leerzeichen-Pfade, Windows, zwei Prozesse, Fehler halbwegs); in-progress prueft Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte, Tabellentests, Geschwister-Stellen per grep; Funde nennen Klasse+Beispiel; Blocker nur falsches Ergebnis/Datenverlust/Absturz/Sicherheit; nach Fix-Durchgang Test+Hunt nur ueber den Diff seit letztem Hunt; 3 Rueckgaben -> Frage an Teamlead; kleine review-gaps sofort schliessen, Vorbestehendes als 'note: pre-existing:'.
