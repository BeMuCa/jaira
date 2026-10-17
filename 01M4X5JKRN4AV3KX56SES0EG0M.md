---
id: 01M4X5JKRN4AV3KX56SES0EG0M
title: "Ein Repo kann mehrere Boards haben, jedes mit eigenen Lanes fuer sein Thema"
status: signoff
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
updated-at: 2026-10-17T12:51:04Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-67696
claimed-at: 2026-10-17T12:50:10Z
outcome-what: "Mehrere Boards pro Repo: core/lane/boards.go (LoadBoard, Boards, NormalizeBoard), Feld board: im Ticket, --board an create/list/next, Einzelticket-Kommandos lesen das Board vom Ticket, set board=<name> mit Fallback auf die erste Lane, TUI-Taste F"
outcome-why: "Ein Repo braucht Boards mit eigenen Lanes neben den SWE-Lanes, z.B. ein video-Board in demo-video fuer eine Designerin"
outcome-resolves: "S0EG0M: alle sechs DoD-Punkte, Tests in internal/cli/boards_test.go und internal/tui/boards_test.go"
review-summary: "Ein Repo kann jetzt neben .jaira/lanes/ weitere Boards unter .jaira/boards/<name>/lanes/ haben (core/lane/boards.go: LoadBoard, Boards, NormalizeBoard, BoardOf, LanesOf). Ein Ticket traegt das Feld board: (core/ticket/schema.go); fehlt es, gilt das heutige Board - bestehende Boards laufen ohne Migration. create/list/next nehmen --board (unbekannter Name: Exit 2 mit Liste der Boards); show/move/dod/claim/set/logbook <id> lesen die Lanes vom Board des Tickets (envFor in internal/cli/root.go). set board=<name> wechselt das Board; fehlt die Lane dort, setzt jaira status auf die erste Lane des neuen Boards und sagt das (landed_in im JSON). validate und resume pruefen jedes Ticket gegen sein eigenes Board. TUI: F schaltet reihum durch die Boards, zeigt deren Lanes und nur deren Tickets, Kopfzeile 'board: <name>'; n und Follow-up schreiben das gezeigte Board ins Ticket. Brettweite Kommandos (logbook --all, milestone, sync, archive, tags-Liste, links, lanes, merge-driver) bleiben auf dem heutigen Board."
round: 5
review-gaps: "1) Merge-Driver (internal/cli/mergedriver.go:55) laedt nur die Lanes des heutigen Boards. core/merge/merge.go mergeStatus vergleicht Precedence: fuer ein Ticket auf einem benannten Board sind dessen Lanes dort unbekannt (-1), also 'kept ours' egal welche Seite weiter ist - oder, wenn eine Lane-ID zufaellig auch im heutigen Board existiert (z.B. human), wird nach der falschen Reihenfolge entschieden und Fortschritt kann zurueckgedreht werden. Fix klein: Board aus ours/theirs-Frontmatter lesen und lane.LoadBoard(root, board) nehmen. Verletzt die Concurrency-Vorgabe, steht aber nicht ausdruecklich im DoD. 2) Frontmatter definition-of-done steht noch auf dem alten Text 'Design als Spec unter docs/superpowers/specs/, von Berk freigegeben' - der Body-DoD und die decision-Notes sagen etwas anderes; wer das Feld liest, sieht einen nie erfuellten DoD. Fix: jaira set S0EG0M definition-of-done=... 3) Klein/kosmetisch: set board=default hinterlaesst 'board: \"\"' statt die Zeile zu entfernen; harmlos, BoardOf liest es als heutiges Board."
test-verdict: "pass (Runde 2): go clean -testcache; go test -race ./... RC=0; DoD 2-6 im Baum (boards_test.go cli+tui gruen, neu TestValidateChecksEachTicketAgainstItsBoard und createTicket-Teil von TestShiftFPicksABoard); Binaer: video-Ticket in Lane script, jaira validate RC=0 ohne unknown_lane (Review-Gap 2 behoben). TUI-Taste F und n nur ueber Tests"
review-verdict: "Alle sechs DoD-Punkte sind im Diff erfuellt und im Temp-Repo nachgestellt (create/list/next --board, Tippfehler Exit 2, set board= mit Fallback-Meldung, validate gruen, go test ./... gruen). Nicht ganz fertig wegen Gap 1: der Merge-Driver entscheidet Status-Konflikte auf benannten Boards nach den falschen Lanes - fuer ein Repo mit parallelen Sessions der Fall, fuer den das Feature gebaut ist. Empfehlung: Gap 1 und 2 vor dem Annehmen beheben (zusammen ca. 30 Minuten), oder bewusst als Folgearbeit annehmen."
review-round: 5
review-check: "1. go build -o /tmp/jaira ./cmd/jaira  2. In einem leeren Git-Repo: /tmp/jaira init, dann mkdir -p .jaira/boards/video/lanes und drei Dateien skript.md, bau.md, done.md mit Frontmatter 'id: <name>' (done.md zusaetzlich 'terminal: true') und eine Datei order mit skript, bau, done je Zeile  3. /tmp/jaira create 'Szene A' --board video -> 'In skript'  4. /tmp/jaira list -> Szene A fehlt; /tmp/jaira list --board video -> Szene A in skript  5. /tmp/jaira list --board vidoe -> 'no such board \"vidoe\"; boards here: default, video', echo $? zeigt 2  6. /tmp/jaira move <id> --to bau, dann /tmp/jaira set <id> board=default -> 'Board default has no lane \"bau\"; <id> is now in <erste Lane>'  7. /tmp/jaira validate -> keine error-Zeile  8. /tmp/jaira (TUI), F druecken -> Kopfzeile 'board: video', nur Lanes skript/bau/done; noch einmal F -> zurueck auf das heutige Board"
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
- **2026-10-17 12:50 · Alexander Sacharov** — review (frische Session, Runde 2): DoD 1-6 erfuellt und im Temp-Repo nachgestellt, go test ./... gruen. Offen: (1) internal/cli/mergedriver.go:55 merged Status-Konflikte eines Tickets auf benanntem Board mit den Lanes des heutigen Boards (core/merge/merge.go mergeStatus, Precedence -1 bzw. falsche Reihenfolge) - Fix: Board aus dem Frontmatter lesen, lane.LoadBoard. (2) Frontmatter definition-of-done traegt noch den alten Spec-Text von Berk. (3) kosmetisch: set board=default laesst 'board: ""' stehen. Nach signoff, Person entscheidet: Gap 1+2 jetzt beheben oder annehmen.
- **2026-10-17 12:51 · Alexander Sacharov** — Worker: zweiter Review (Verdikt 'nicht ganz fertig') hat nach signoff bewegt; zurueck nach in-progress fuer Gap 1 (Merge-Driver nimmt die Lanes des heutigen Boards) und Gap 2 (Frontmatter definition-of-done veraltet). Gap 3 (board: "" nach set board=default) bleibt: BoardOf liest es als heutiges Board, Doc hat kein Entfernen eines Felds.
