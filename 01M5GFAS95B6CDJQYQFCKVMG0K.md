---
id: 01M5GFAS95B6CDJQYQFCKVMG0K
title: "Boards wechseln mit Shift+Pfeil und ein Board-Fenster auf b mit aktiven Boards, Ersteller und Lanes"
status: in-progress
ready: true
creator: berkc
assignee: berkc
goal: "Zwischen Boards springt man mit Shift+hoch/runter, und b zeigt alle Boards mit Name, Ersteller und Lanes; dort schaltet die Leertaste ein Board fuer den Teamlead aktiv oder inaktiv."
context: |-
  Wunsch von berkc am 22.10.2026, woertlich: 'pressing shift + Arrow keys (up/down) Rotates through the boards that are active. We said that the active boards are selected in the Settings, but we should also add the hotkey b, which opens up a pop up with all the available boards as long tabs, with the name of the board as a header, creator of the board and the lanes listed below. With the active boards (the one that the team lead can choose from) having a tick symbol. On that view you can select and unselect the boards with the space key. The order of the boards is: active boards, most recently active, all the other boards.'
  Heute: b auf dem Board springt zum naechsten Board (nextBoard(), internal/tui/model.go); im geoeffneten Ticket springt b zum Ticket, von dem es blockiert wird. Das Ticket zum Link-Oeffnen-Fehler will b im Ticket ausserdem fuer 'zeig mir dieses Ticket auf dem Board' - die Belegung von b muss einmal fuer alle drei entschieden werden.
  Verwandt: 57TDND (Settings: aktive Boards fuer den Teamlead waehlen) - das Fenster hier ist eine zweite Stelle fuer dieselbe Auswahl, muss also dieselbe Einstellung schreiben. Y2C9QN (Boards aus Board-Typen). Ein Board hat heute einen description-Datei unter .jaira/boards/<name>/; einen Ersteller speichert es noch nicht.
definition-of-done: Shift+hoch/runter wechselt reihum durch die aktiven Boards; Test in internal/tui
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T00:14:17Z
updated-at: 2026-10-30T02:11:56Z
updated-by: berkc
claimed-by: berk-ca-1397863
claimed-at: 2026-10-30T01:57:38Z
---

# Boards wechseln mit Shift+Pfeil und ein Board-Fenster auf b mit aktiven Boards, Ersteller und Lanes

## Definition of Done

- [x] Shift+hoch/runter wechselt reihum durch die aktiven Boards; Test in internal/tui
  proof: internal/tui/model.go:1718 stepBoard; TestShiftArrowsStepThroughActiveBoardsOnly, TestShiftDownPicksABoard
- [x] b oeffnet ein Fenster: pro Board Name als Ueberschrift, Ersteller, Lanes; aktive mit Haken; Reihenfolge aktiv, zuletzt benutzt, Rest; Test
  proof: internal/tui/boardwindow.go:120 openBoards/renderBoards; TestBOpensTheBoardWindow
- [x] Leertaste schaltet aktiv/inaktiv und schreibt dieselbe Einstellung wie 57TDND; Test
  proof: internal/tui/boardwindow.go:179 -> core/lane/active.go:69 SetBoardActive (.jaira/active-boards, the setting 57TDND is to read); TestSpaceTogglesABoardActive, TestActiveBoards
- [x] Die Belegung von b (Board, Blocked-by-Sprung, Ticket auf dem Board zeigen) ist entschieden und in der Hilfe beschrieben
  proof: internal/tui/view.go:1767 help line; b on the board = board window, b in an open ticket left to P3VS02 (decision berkc 30.10.)
- [x] Release-Note-Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/KVMG0K.md

## Options

- [x] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/lane/active.go: ActiveBoardsPath, ActiveBoards(root) (default zuerst, dann Name; keine Datei = alle), SetBoardActive(root,name,on) - das letzte aktive bleibt; Tests in core/lane
- [x] core/gitrepo: FirstAuthor(path) - Autor des ersten Commits, der path anlegte; Test
- [x] internal/tui/boardwindow.go: modeBoards, openBoards/keyBoards/renderBoards (Name als Ueberschrift, Ersteller, Lanes, Haken, Reihenfolge aktiv/zuletzt benutzt/Rest), space schaltet, enter zeigt das Board; zuletzt-benutzt-Stempel unter Store.RepoStateDir()
- [x] model.go: b auf dem Board oeffnet das Fenster; shift+down/up wechselt reihum durch die aktiven Boards (nextBoard ersetzt); shift im offenen Ticket unveraendert
- [x] Hilfe (view.go) beschreibt b und shift+↓↑ auf dem Board; b im Ticket bleibt P3VS02
- [x] Tests in internal/tui: TestBPicksABoard umstellen, Fenster-Reihenfolge/Inhalt, space schreibt .jaira/active-boards, shift cycelt nur aktive
- [ ] Release-Note core/release/notes/KVMG0K.md; go vet ./... und go test ./...

## Progress
- **2026-10-23 21:04 · berkc** — Entscheidung berkc 23.10.2026: Taste b bleibt fuer Boards (Board-Popup). Die Sprung-Funktion aus P3VS02 braucht eine andere Taste.
- **2026-10-30 01:59 · berkc** — brainstorm: Was heute ist (geprueft): b auf dem Board ruft nextBoard() (internal/tui/model.go:1705, 2353) und springt reihum durch .jaira/lanes + .jaira/boards/*. Shift+hoch/runter ist auf dem Board frei; im offenen Ticket scrollt es (model.go:1575) - bleibt so. Eine Einstellung 'aktive Boards' gibt es nicht (57TDND ist noch backlog). Ein Board speichert keinen Ersteller (auch Y2C9QN/4G07X9 nicht).
- **2026-10-30 01:59 · berkc** — brainstorm Optionen fuer den Ort der aktiven Boards: (A) im Repo, Datei .jaira/active-boards, ein Boardname pro Zeile, fehlt die Datei sind alle Boards aktiv - gibt: das ganze Team und jeder Teamlead sieht dieselbe Auswahl, handeditierbar, merge-freundlich; kostet: ein Leertaste-Druck macht eine Datei im Repo geaendert, die man committen muss. (B) pro Rechner in ~/.jaira/settings.json, nach Repo-Pfad geschluesselt - gibt: kein Commit; kostet: zwei Teamleads im Team legen dasselbe Ticket auf verschiedene Boards, und jede Worktree/jeder Klon hat einen anderen Pfad. (C) eine Marker-Datei pro Board in dessen Lanes-Ordner - gibt: zwei Umschaltungen kollidieren nie; kostet: eine Einstellung ueber N Dateien verteilt, today's board liegt woanders (.jaira/lanes), und es greift in die Board-Ordner ein, an denen 4G07X9 arbeitet.
- **2026-10-30 01:59 · berkc** — brainstorm Entscheidung (Worker, Begruendung): (A) .jaira/active-boards im Repo. Grund: 'aktiv' sagt, wohin der Teamlead neue Tickets legen darf - das ist eine Regel des Teams, nicht eine Vorliebe einer Person; mit (B) wuerden zwei Personen dasselbe Ticket auf verschiedene Boards legen. Vorbild: .jaira/chat-model ist ebenfalls eine Board-Einstellung im Repo. Keine Datei = alle aktiv, damit aendert sich fuer bestehende Repos nichts. Das letzte aktive Board laesst sich nicht abschalten (sonst haette der Teamlead nichts zur Wahl). 57TDND soll dieselbe Datei lesen/schreiben (Funktionen in core/lane).
- **2026-10-30 01:59 · berkc** — brainstorm weitere Festlegungen (Worker): 'zuletzt benutzt' = wann das Board zuletzt in der TUI gezeigt wurde, pro Rechner und Repo unter Store.RepoStateDir() (~/.jaira/repo/<key>/), weil das Navigation einer Person ist und kein Teamzustand. Ersteller = Autor des ersten Commits, der den Lanes-Ordner des Boards anlegte (git log), 'noch nicht committet' wenn es keinen gibt - so braucht es kein neues Feld und kein Eingriff in boards create (4G07X9). Reihenfolge im Fenster: aktive in fester Reihenfolge (default zuerst, dann Name - dieselbe, durch die Shift+Pfeil laeuft), dann inaktive zuletzt benutzte (neueste zuerst), dann der Rest nach Name. enter im Fenster zeigt das gewaehlte Board (sonst kaeme man an ein inaktives Board nicht mehr heran, weil b nicht mehr reihum springt). Shift+Pfeil auf einem inaktiven Board geht zum naechsten aktiven. b im offenen Ticket bleibt hier unberuehrt (P3VS02).
- **2026-10-30 02:00 · berkc** — pre-process: jaira boards --json und der Teamlead-Skill lesen active-boards hier NICHT - das ist 57TDNDs DoD und internal/cli/boards.go ist gerade Baustelle von 4G07X9/Y2C9QN. Dieses Ticket legt Speicherort und core/lane-Funktionen fest, die 57TDND dann nur noch aufruft.
