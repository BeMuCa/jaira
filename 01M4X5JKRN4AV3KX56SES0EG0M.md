---
id: 01M4X5JKRN4AV3KX56SES0EG0M
title: "Ein Repo kann mehrere Boards haben, jedes mit eigenen Lanes fuer sein Thema"
status: critique
ready: true
creator: BeMuCa
assignee: Alexander Sacharov
goal: "In einem Repo laufen mehrere Boards nebeneinander (z.B. backend, ui, video), jedes mit eigenen Lanes; ein Ticket gehoert genau einem Board und laeuft durch dessen Lanes"
context: |-
  Was heute fehlt: Ein Repo hat genau ein Board. Ein Board ist sein Lane-Ordner (.jaira/lanes/, siehe README 'A board is its lane directory'). Alle Tickets laufen durch dieselben Lanes.
  Unsere Lanes passen fuer Software-Entwicklung (SWE). Fuer andere Arbeit passen sie nicht.
  Berks Idee vom 14.10., seine Beispiele:
  - UI-Board: Lanes speziell fuer UI, z.B. verlangt testing dort immer Playwright-Tests.
  - Video-Board: Lanes fuer Ideen-Generierung, Skript, Voice-over und Schnitt.
  - Backend-Board: z.B. fuer 'in der Datenbank fuer user eine Spalte hinzufuegen'.
  Jedes ist ein eigenes Board. Der Teamlead verteilt Tickets darauf (eigenes Ticket), und ein jaira-worker bekommt Ticket plus Board und treibt es durch dessen Lanes.
  Offen fuer brainstorm:
  - wo die Boards liegen (z.B. .jaira/boards/<name>/lanes/) und wie ein Ticket sein Board traegt (Feld board:)
  - wie TUI und Launcher mehrere Boards eines Repos zeigen
  - wie es mit dem heutigen Einzel-Board vertraeglich bleibt (Berks Regel: ein Update darf nichts verschwinden lassen)
  Haengt am Ticket-Schema (BGZ0SW, NFJCTK): eine Lane eines solchen Boards muss eigene Felder erklaeren koennen.
definition-of-done: "Design als Spec unter docs/superpowers/specs/, von Berk freigegeben: Ablage der Boards, Board-Feld am Ticket, Anzeige in TUI und Launcher, Migration des heutigen Boards"
tags: []
blocked-by: []
related:
  - 01M4X5J1QRR70AM2B2JPBGZ0SW
  - 01M11YYR0SMVCXNQPZ68NFJCTK
commits: []
created-at: 2026-10-14T12:18:13Z
updated-at: 2026-10-17T12:44:35Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-73711
claimed-at: 2026-10-17T12:17:54Z
outcome-what: "Mehrere Boards pro Repo: core/lane/boards.go (LoadBoard, Boards, NormalizeBoard), Feld board: im Ticket, --board an create/list/next, Einzelticket-Kommandos lesen das Board vom Ticket, set board=<name> mit Fallback auf die erste Lane, TUI-Taste F"
outcome-why: "Ein Repo braucht Boards mit eigenen Lanes neben den SWE-Lanes, z.B. ein video-Board in demo-video fuer eine Designerin"
outcome-resolves: "S0EG0M: alle sechs DoD-Punkte, Tests in internal/cli/boards_test.go und internal/tui/boards_test.go"
review-summary: "Runde 5: keine Befunde. lane.LanesOf ersetzt beide Caches; TicketsOnBoards haelt Tickets() als Einzelboard-Huelle, damit die bestehenden Tests unveraendert bleiben. Bekannte Grenze (Note Stand): ein Blocker auf einem anderen Board wird im Link-Index mit den Lanes des gezeigten Boards beurteilt."
round: 5
review-gaps: "1) Fehler TUI: n (createTicket, internal/tui/model.go:2122) und Follow-up (followUpFields, internal/tui/followup.go:173) schreiben auf einem mit F gewaehlten Board status=erste Lane dieses Boards, aber kein board:. Das Ticket gehoert damit zum heutigen Board, verschwindet sofort aus der Ansicht, und hat das heutige Board diese Lane nicht (z.B. 'script'), ist es read-only. Fix: fields[ticket.FieldBoard]=m.boardName, wenn nicht leer; Test in internal/tui/boards_test.go. 2) Fehler: jaira validate laedt nur das heutige Board (internal/cli/validate.go:35) und meldet jedes Ticket in einer Lane, die nur sein Board hat, als error 'lane is not installed ... read-only', Exit 3. Nachgestellt im Temp-Repo: Board video mit Lane script, Ticket dorthin bewegt -> validate exit 3. Die Note 'Stand' nennt validate bewusst brettweit, die Folge (falscher Fehler auf einem gueltigen Repo) steht dort nicht. 3) Kleiner: Die TUI baut den Link-Index aus den gefilterten Tickets (model.go linkIndex mit m.tickets), die CLI aus allen Boards; ein Blocker auf einem anderen Board erscheint in der TUI als unbekannt. 4) Bewusst ausgelassen laut Notes, kein Befund: logbook ohne id, milestone, sync, archive, lanes, links, merge-driver bleiben auf dem heutigen Board. DoD 1-6 sind durch den Diff gedeckt; DoD 1 weicht vom urspruenglichen DoD-Text (Spec von Berk freigegeben) ab, per decision-Note von Alexander am 17.10."
test-verdict: "pass: go test -race ./... nach go clean -testcache RC=0; DoD 2-6 im Baum geprueft (boards_test.go in cli und tui laufen gruen, release note da); Binaer auf Scratch-Repo: create/list/next --board video, unbekanntes Board Exit 2 mit Liste, move/claim/dod/note/tag auf video-Ticket, set board=video landet in backlog mit Meldung, board= zurueck. TUI-Taste F nur ueber Tests (echte Tastendispatch), nicht interaktiv"
review-verdict: "Nicht fertig: Die sechs DoD-Punkte sind erfuellt und getestet, aber zwei Fehler gehoeren in diesen Umfang - ein in der TUI auf einem weiteren Board angelegtes Ticket landet ohne board: auf dem falschen Board (Gap 1), und jaira validate meldet gueltige Tickets anderer Boards als Fehler mit Exit 3 (Gap 2). Empfehlung: zurueck nach in-progress fuer beide, je ein Test; Gap 3 kann mit."
review-round: 3
review-check: "1. cd /home/alex/projects/.worktrees/jaira-s0eg0m && go build -o /tmp/jb ./cmd/jaira  2. mkdir /tmp/mb && cd /tmp/mb && git init -q && git commit -q --allow-empty -m i && /tmp/jb init  3. mkdir -p .jaira/boards/video/lanes && cp .jaira/lanes/backlog.md .jaira/lanes/done.md .jaira/boards/video/lanes/ && sed 's/^id: todo/id: script/' .jaira/lanes/todo.md > .jaira/boards/video/lanes/script.md && printf 'backlog\\nscript\\ndone\\n' > .jaira/boards/video/lanes/order  4. /tmp/jb create 'vid' --board video --goal g --context c --dod d -> 'Created <ID>'  5. /tmp/jb list -> vid fehlt; /tmp/jb list --board video -> vid steht in Backlog  6. /tmp/jb list --board vidoe -> Exit 2, nennt 'default, video'  7. /tmp/jb set <ID> assignee=T && /tmp/jb move <ID> --to script && /tmp/jb show <ID> -> Zeilen 'board video', 'lane script'  8. /tmp/jb set <ID> board=default -> 'Board default has no lane \"script\"; <ID> is now in backlog'  9. /tmp/jb set <ID> board=video && /tmp/jb move <ID> --to script && /tmp/jb validate -> heute: error 'lane \"script\" is not installed', Exit 3 (Gap 2)  10. /tmp/jb (TUI) -> F druecken: Kopf zeigt 'board: video', Spalten Backlog/Script/Done, nur vid; nochmal F -> zurueck  11. auf video n druecken, Titel tippen, Enter -> heute verschwindet das neue Ticket aus der Ansicht (Gap 1)"
---

# Ein Repo kann mehrere Boards haben, jedes mit eigenen Lanes fuer sein Thema

## Definition of Done

- [x] Entscheidungen stehen als decision-Notes am Ticket (Ablage .jaira/boards/<name>/lanes/, Feld board:, Boardwechsel, Migration); keine eigene Spec, Alexander hat die Richtung am 17.10. vorgegeben
  proof: decision-Notes vom 17.10. am Ticket (Ablage, Name/Existenz, Feld board:/Migration, list/next-Default, Boardwechsel)
- [x] TUI: Shift+F waehlt ein Board; das Board zeigt dann dessen Lanes und nur dessen Tickets; ohne Wahl das heutige Board; Test. Mehr TUI-Darstellung nicht
  proof: internal/tui/boards_test.go TestShiftFPicksABoard, TestShiftFOnASingleBoardExplains; internal/tui/model.go nextBoard und Filter in reload()
- [x] jaira create, list, next und show nehmen ein Board (z.B. --board ui); ohne Angabe gilt das heutige Board, nichts aendert sich still; Test
  proof: internal/cli/boards_test.go TestCreateListNextTakeABoard, TestUnknownBoardIsRefused; --board in internal/cli/tickets.go (create, list) und flow.go (next); show liest das Board vom Ticket (envFor, internal/cli/root.go)
- [x] Ein bestehendes Einzel-Board laeuft nach dem Update unveraendert weiter; Test
  proof: internal/cli/boards_test.go TestSingleBoardRunsUnchanged; ohne board: gilt lane.Load wie bisher (core/lane/boards.go LoadBoard)
- [x] Eine Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/S0EG0M.md
- [x] Ein Ticket wechselt das Board mit jaira set <id> board=<name>; existiert seine Lane dort nicht, landet es in der ersten Lane des neuen Boards und jaira sagt das; Test
  proof: internal/cli/boards_test.go TestSetBoardMovesToFirstLaneWhenItsLaneIsMissing; internal/cli/tickets.go set: FieldBoard-Zweig

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/lane: LoadBoard(root,name), BoardLanesDir, Boards(root), ValidBoardName; Load teilt den Rest mit LoadBoard
- [x] core/ticket: Feld board (FieldBoard, Ticket.Board)
- [x] cli: lanesFor/envFor per Board; create/list/next --board; show/move/dod/claim lesen das Board vom Ticket
- [x] cli: set board=<name> mit Fallback auf erste Lane und Meldung
- [x] tui: Shift+F waehlt Board, Lanes und Tickets gefiltert
- [x] Tests + Release-Note

## Progress
- **2026-10-17 12:17 · Alexander Sacharov** — Alexander Sacharov took this ticket over from BeMuCa
- **2026-10-17 12:17 · Alexander Sacharov** — decision 17.10. (Alexander, Teamlead-Session): Alexander uebernimmt das Ticket von Berk und grenzt den Umfang ein: kein von Berk freigegebener Spec, TUI nur Filter per Shift+F. Tickets bekommen KEINE Ordner pro Board — sie leben bis zum Logbook auf refs/jaira/tickets/<ID>, der Pfad traegt keine Information. Zugehoerigkeit ist das Frontmatter-Feld board:; fehlt es, gilt das heutige Board (.jaira/lanes/), so laeuft jedes bestehende Board unveraendert. Ordner nur fuer die Lanes eines weiteren Boards: .jaira/boards/<name>/lanes/ mit eigener order. Boardwechsel = eine Zeile im Frontmatter statt Rename (kein rename/modify-Konflikt). Ausloeser: ~/projects/demo-video braucht ein Board 'video' (Skript -> human -> Bau, Szenen parallel -> Probe-Render -> human -> Stimme -> Finalrender), damit eine Designerin ohne Programmierkenntnis Videos nach Vorlage anlegt, waehrend Code und Animationen (shared/, scripts/) weiter durch die SWE-Lanes laufen.
- **2026-10-17 12:18 · Alexander Sacharov** — parent session: jaira-79
- **2026-10-17 12:19 · Alexander Sacharov** — decision: Ablage. Ein weiteres Board ist .jaira/boards/<name>/lanes/ mit Lane-Dateien und eigener order-Datei. Das heutige Board (.jaira/lanes/) hat keinen Namen; auf der Kommandozeile steht 'default' fuer es (--board default, set board=default oder board= leer). Ein Board-Name ist klein, Ziffern und Bindestriche, 'default' ist reserviert.
- **2026-10-17 12:19 · Alexander Sacharov** — decision: Ein benanntes Board existiert genau dann, wenn sein lanes/-Ordner Lane-Dateien enthaelt. Ein unbekannter Name ist ein Fehler (Exit 2) und nennt die vorhandenen Boards. Kein automatisches Anlegen wie beim heutigen Board: ein Tippfehler in --board darf kein neues Board erzeugen. Ein Board legt man an, indem man den Ordner mit Lane-Dateien fuellt (z.B. aus .jaira/lanes/ kopiert); kein eigenes Kommando - Scope.
- **2026-10-17 12:19 · Alexander Sacharov** — decision: Feld board: im Frontmatter. Fehlt es, gehoert das Ticket zum heutigen Board - das ist die ganze Migration: bestehende Tickets und .jaira/lanes/ bleiben unberuehrt, nichts wird umgeschrieben.
- **2026-10-17 12:19 · Alexander Sacharov** — decision: jaira list und next ohne --board zeigen nur Tickets des heutigen Boards, wie die TUI ohne Wahl. Kommandos auf ein einzelnes Ticket (show, move, dod, claim, set, note) brauchen kein --board: sie lesen das Board vom Ticket und pruefen es gegen dessen Lanes.
- **2026-10-17 12:19 · Alexander Sacharov** — decision: Boardwechsel ist jaira set <id> board=<name>, eine Zeile im Frontmatter. Gibt es die Lane des Tickets im neuen Board nicht, setzt jaira status auf die erste Lane in der order des neuen Boards und sagt das auf stdout.
- **2026-10-17 12:23 · Alexander Sacharov** — Stand: core/lane/boards.go (LoadBoard, Boards, NormalizeBoard) und Feld board in core/ticket/schema.go fertig; CLI-Helfer loadBoardEnv/envFor/onBoard in internal/cli/root.go. Brettweite Kommandos (logbook ohne id, milestone, tags-Liste, sync, archive, validate, links, lanes) bleiben bewusst auf dem heutigen Board - nicht im DoD. Bekannte Grenze: link.Build beurteilt einen Blocker auf einem anderen Board mit den Lanes des gerade geladenen Boards.
- **2026-10-17 12:29 · Alexander Sacharov** — decision: TUI zeigt neben dem Repo-Namen 'board: <name>', solange nicht das heutige Board offen ist - ohne die Zeile sieht man nach F nicht, welches Board man vor sich hat. F schaltet reihum (heute, dann die anderen nach Name), statt eines Auswahlfensters: kleiner, und bei zwei, drei Boards gleich schnell. Die Hilfe fuehrt F in der Zeile von p, weil eine eigene Zeile '@name' aus dem 60-Zeilen-Fenster von TestHelpExplainsEveryCardMark schob.
- **2026-10-17 12:29 · Alexander Sacharov** — critique: internal/cli/tickets.go set: der board-Zweig steht zwischen dem mode-Kommentar und if k == FieldMode - der Kommentar beschreibt jetzt den falschen Block; board-Zweig vor den mode-Kommentar ziehen
internal/cli/tickets.go create und set: loadBoardEnv/envFor lesen alle Tickets (s.List) nur fuer die Lanes - ein Helfer boardLanes(root, board) in root.go, der ErrNoSuchBoard auf Exit 2 abbildet, statt dessen; loadBoardEnv benutzt ihn selbst
internal/cli/checklist.go:137, tags.go:320: lane.LoadBoard direkt - unbekanntes Board gibt dort Exit 1 statt 2; denselben Helfer nehmen
internal/cli/root.go onBoard und internal/tui/model.go reload filtern beide mit NormalizeBoard(t.Board) == board - eine Funktion lane.BoardOf(t) dafuer
core/lane/boards.go BoardsDir/BoardLanesDir haben keinen Aufrufer ausserhalb des Pakets - klein schreiben
- **2026-10-17 12:31 · Alexander Sacharov** — critique (Runde 2): internal/cli/root.go onBoard: der Fehlerzweig ist unerreichbar - beide Aufrufer (list, next) haben den Namen schon durch loadBoardEnv geprueft; onBoard gibt nur die Liste zurueck. Sonst nichts: boardLanes, BoardOf und die set-Reihenfolge aus Runde 1 sind erledigt.
- **2026-10-17 12:39 · Alexander Sacharov** — review (frische Session): zwei Fehler im Umfang. (1) TUI-Anlegen per n und Follow-up auf einem mit F gewaehlten Board schreibt kein board: - Ticket landet auf dem heutigen Board, ggf. in einer Lane, die es dort nicht gibt (internal/tui/model.go createTicket, internal/tui/followup.go followUpFields). (2) jaira validate kennt nur .jaira/lanes/ und meldet Tickets in board-eigenen Lanes als error, Exit 3 (internal/cli/validate.go:35); im Temp-Repo nachgestellt. Klein: TUI-Link-Index nur aus dem gefilterten Board. Empfehlung zurueck nach in-progress.
- **2026-10-17 12:39 · Alexander Sacharov** — Worker: Review-Agent hat mit Verdikt 'Nicht fertig' (Gaps 1-3) nach signoff bewegt; zurueck nach in-progress fuer Gap 1 (TUI n/Follow-up ohne board:), Gap 2 (validate brettweit), Gap 3 (TUI-Link-Index aus gefilterten Tickets).
- **2026-10-17 12:42 · Alexander Sacharov** — Review-Gaps 1-3 behoben: TUI n und Follow-up schreiben board: des mit F gewaehlten Boards (internal/tui/model.go createTicket, followup.go followUpFields; nur wenn nicht leer, weil Store.Create leere Felder als 'board: ""' schreibt). validate prueft per validate.TicketsOnBoards jedes Ticket gegen die Lanes seines Boards; Board ohne Lane-Dateien = Fehler unknown_lane am Feld board. TUI: m.allTickets fuer Link-Index und gateEnv.All, m.tickets bleibt das gezeigte Board. Tests: TestValidateChecksEachTicketAgainstItsBoard, TestShiftFPicksABoard (createTicket-Teil).
- **2026-10-17 12:43 · Alexander Sacharov** — critique (Runde 4): internal/cli/resume.go und internal/cli/validate.go bauen je denselben Cache Lanes-pro-Board; eine Funktion lane.LanesOf(root) in core/lane/boards.go, beide rufen sie
