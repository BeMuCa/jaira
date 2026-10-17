---
id: 01M4X5JKRN4AV3KX56SES0EG0M
title: "Ein Repo kann mehrere Boards haben, jedes mit eigenen Lanes fuer sein Thema"
status: done
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
definition-of-done: "Mehrere Boards pro Repo laufen (Lanes unter .jaira/boards/<name>/lanes/, Feld board:), create/list/next --board, set board=, TUI Shift+F, Einzel-Board unveraendert; Entscheidungen als decision-Notes statt Spec (Alexander, 17.10.)"
tags: []
blocked-by: []
related:
  - 01M4X5J1QRR70AM2B2JPBGZ0SW
  - 01M11YYR0SMVCXNQPZ68NFJCTK
commits: []
created-at: 2026-10-14T12:18:13Z
updated-at: 2026-10-17T14:02:52Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-67696
claimed-at: 2026-10-17T12:50:10Z
outcome-what: "Mehrere Boards pro Repo: core/lane/boards.go (LoadBoard, Boards, NormalizeBoard), Feld board: im Ticket, --board an create/list/next, Einzelticket-Kommandos lesen das Board vom Ticket, set board=<name> mit Fallback auf die erste Lane, TUI-Taste F"
outcome-why: "Ein Repo braucht Boards mit eigenen Lanes neben den SWE-Lanes, z.B. ein video-Board in demo-video fuer eine Designerin"
outcome-resolves: "S0EG0M: alle sechs DoD-Punkte, Tests in internal/cli/boards_test.go und internal/tui/boards_test.go"
review-summary: "Runde 7 (Taste b): keine Befunde. b auf der Board-Ebene war frei; die Detailansicht (model.go modeDetail, b = Blocker folgen) und die Milestone-Ebene (layer.go keyLayer) kehren vor der Board-Ebene zurueck, die Bedeutungen ueberschneiden sich nicht. Bewusst stehen gelassen (decision Alexander): derselbe Buchstabe heisst auf dem Board 'naechstes Board' und im offenen Ticket 'Blocker folgen'."
round: 7
review-gaps: "Runde 4 (Taste b): nichts zu entfernen - eine Zeile case, Hilfe, Note und Test umbenannt. Frueheres siehe Notes."
test-verdict: "pass (Runde 4, Taste b): go clean -testcache; go test -race ./... RC=0; TestBPicksABoard und TestBOnASingleBoardExplains druecken b durch die echte Tastendispatch; TestHelpExplainsEveryCardMark gruen (Hilfe 'p b' haelt das 60-Zeilen-Fenster)"
review-verdict: "Die Diff erfuellt alle sechs DoD-Punkte; go test ./... gruen; create/list/next --board, Tippfehler-Board (Exit 2, kein Ordner angelegt), set board= hin und zurueck mit Landen in der ersten Lane und validate im Temp-Repo mit dem gebauten Binary nachgestellt. Die Gaps der beiden frueheren Reviews (TUI n/Follow-up, validate, Link-Index, Merge-Driver, Frontmatter-DoD) sind behoben. Offen nur Grenzen fuer brettweites Logbook-Ablegen und kosmetische board:-Schreibweisen - nicht blockierend. Annehmbar."
review-round: 6
review-check: "1. go build -o /tmp/jaira ./cmd/jaira  2. mkdir /tmp/r && cd /tmp/r && git init && git commit --allow-empty -m i && /tmp/jaira init  3. mkdir -p .jaira/boards/video/lanes && cp .jaira/lanes/*.md .jaira/boards/video/lanes/ (oder eigene Lane-Dateien)  4. /tmp/jaira create V1 --board video --goal g --context c --dod d  -> das Ticket hat 'board: video'  5. /tmp/jaira list  -> V1 fehlt; /tmp/jaira list --board video -> V1 steht da  6. /tmp/jaira list --board vidoe -> 'no such board \"vidoe\"; boards here: default, video', Exit 2, kein Ordner vidoe entsteht  7. /tmp/jaira set <V1-id> board=default -> 'Updated', ggf. 'Board default has no lane ...; ... is now in backlog'  8. /tmp/jaira tui (bzw. jaira), F druecken -> Kopfzeile zeigt 'board: video', nur dessen Lanes und Tickets; nochmal F -> zurueck zum heutigen Board"
---

# Ein Repo kann mehrere Boards haben, jedes mit eigenen Lanes fuer sein Thema

## Definition of Done

- [x] Entscheidungen stehen als decision-Notes am Ticket (Ablage .jaira/boards/<name>/lanes/, Feld board:, Boardwechsel, Migration); keine eigene Spec, Alexander hat die Richtung am 17.10. vorgegeben
  proof: decision-Notes vom 17.10. am Ticket (Ablage, Name/Existenz, Feld board:/Migration, list/next-Default, Boardwechsel)
- [x] TUI: b waehlt ein Board (urspruenglich Shift+F, per decision 17.10. b); das Board zeigt dann dessen Lanes und nur dessen Tickets; ohne Wahl das heutige Board; Test. Mehr TUI-Darstellung nicht
  proof: internal/tui/boards_test.go TestBPicksABoard, TestBOnASingleBoardExplains; internal/tui/model.go nextBoard und Filter in reload()
- [x] jaira create, list, next und show nehmen ein Board (z.B. --board ui); ohne Angabe gilt das heutige Board, nichts aendert sich still; Test
  proof: internal/cli/boards_test.go TestCreateListNextTakeABoard, TestUnknownBoardIsRefused; --board in internal/cli/tickets.go (create, list) und flow.go (next); show liest das Board vom Ticket (envFor, internal/cli/root.go)
- [x] Ein bestehendes Einzel-Board laeuft nach dem Update unveraendert weiter; Test
  proof: internal/cli/boards_test.go TestSingleBoardRunsUnchanged; ohne board: gilt lane.Load wie bisher (core/lane/boards.go LoadBoard)
- [x] Eine Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/S0EG0M.md
- [x] Ein Ticket wechselt das Board mit jaira set <id> board=<name>; existiert seine Lane dort nicht, landet es in der ersten Lane des neuen Boards und jaira sagt das; Test
  proof: internal/cli/boards_test.go TestSetBoardMovesToFirstLaneWhenItsLaneIsMissing; internal/cli/tickets.go set: FieldBoard-Zweig
- [x] Board wechselt mit b statt F; Hilfe, Release-Note und Test nennen b; Test
  proof: internal/tui/model.go case "b" (Board-Ebene) -> nextBoard; Hilfe internal/tui/view.go Zeile 'p b'; core/release/notes/S0EG0M.md Zeile 3; Test internal/tui/boards_test.go TestBPicksABoard (press b)

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
- **2026-10-17 12:53 · Alexander Sacharov** — Review-2-Gaps: Gap 1 behoben - merge-driver laedt die Lanes des Boards aus ours (internal/cli/mergedriver.go mergeBoard), Test TestMergeDriverUsesTheTicketsBoard (gegengeprueft: rot ohne Fix; erste Fassung mit 'done' war nicht trennscharf, weil done auch auf dem heutigen Board liegt). Wechselt eine Seite das Board, gilt das Board von ours. Gap 2: definition-of-done im Frontmatter neu gesetzt.
- **2026-10-17 12:56 · Alexander Sacharov** — review (frische Session, Runde 3): DoD 1-6 erfuellt, go test ./... gruen, im Temp-Repo mit gebautem Binary nachgestellt. Fruehere Gaps behoben. Nicht blockierend offen: brettweites Logbook-Ablegen (logbook.go:136, lane/settle.go:18) arbeitet nach Lane-ID ueber alle Boards; Ticket eines verschwundenen Boards nur per set board=default rettbar; board:-Feld in drei Schreibweisen (leer-quoted, quoted, unquoted). Verdikt: annehmbar.
- **2026-10-17 12:57 · Alexander Sacharov** — explain: Der Review fand keine Luecke im DoD, nur drei kleine Grenzen: Board-weites Ablegen ins Logbook arbeitet nach Lane-ID ueber alle Boards, ein Ticket eines verschwundenen Boards ist nur per set board=default rettbar, und das Feld board: steht in drei Schreibweisen im Frontmatter.
Warum: Das Logbook kann ein Ticket auf einem anderen Board ablegen, wenn dessen Lane gleich heisst, aber nicht terminal ist; die anderen beiden Punkte sind nur unschoen, nichts geht verloren.
1) Zurueckgeben: Logbook-Ablegen auf das Board des Tickets begrenzen und die board:-Schreibweise vereinheitlichen [give back]
2) So annehmen und die Grenzen als bekannt stehen lassen, sie sind nicht blockierend [accept]
Empfehlung: 2 - alle sechs DoD-Punkte sind erfuellt, go test ./... ist gruen, und der Reviewer nennt das Ergebnis annehmbar.
- **2026-10-17 12:57 · Alexander Sacharov** — Abnahmeseite: https://claude.ai/artifact/HdF2uR51Qm7NkPg23nMiX9 (5 Personenschritte TUI, 2 Maschinenchecks gruen; Demo-Repo ~/.claude/acceptance/jaira/feat-s0eg0m/demo, Binaer daneben)
- **2026-10-17 13:02 · Alexander Sacharov** — decision 17.10.: Taste b statt F, Alexander (F/f haben schon andere Bedeutungen, f = Follow-up in der Detailansicht). Weitergegeben von der Teamlead-Session jaira-79.
- **2026-10-17 14:02 · Alexander Sacharov** — decision: 2) So annehmen und die Grenzen als bekannt stehen lassen, sie sind nicht blockierend
