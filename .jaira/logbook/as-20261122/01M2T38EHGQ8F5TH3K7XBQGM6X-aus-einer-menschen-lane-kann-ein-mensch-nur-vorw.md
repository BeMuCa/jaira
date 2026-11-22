---
id: 01M2T38EHGQ8F5TH3K7XBQGM6X
title: "Aus einer Menschen-Lane kann ein Mensch nur vorwaerts oder zur Seite, nicht zurueck"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Mensch, der in human oder signoff hinsieht und Nacharbeit will, schickt das Ticket selbst in eine frueher liegende Lane zurueck - ohne --force, ohne Folgeticket, und ohne dass derselbe Weg nach done aufgeht"
context: |-
  Aufgefallen Alex am 2026-09-18 an 7XXT00: Ticket stand in signoff, die review-Lane hatte drei Befunde geschrieben, einer davon ausdruecklich 'gehoert vor die Annahme'. Es gibt keinen Weg, es zur Nacharbeit zurueckzuschicken.

  Was die Doska heute kann. Auf dem signoff-Schirm stehen genau zwei Aktionen (internal/tui/signoff.go:151): 'a' nimmt an und geht direkt in m.lanes.Terminal(), also done (signoff.go:170), und 'f' legt ein Folgeticket an. Der move-Dialog 'm' hilft nicht: der Gate laesst einen Ausgang aus einer requires-human-exit-Lane nur bei req.Interactive durch (core/gate/gate.go:389), und dieses Flag wird ausschliesslich in signoff.go:179 und :189 gesetzt - das sind accept und followUp. Sonst nirgends. Ein Mensch am Gerät bekommt deshalb dieselbe Absage wie ein Agent: 'An agent cannot move a ticket out of it'.

  Warum es so ist, und warum die Begruendung nicht mehr traegt. Absichtlich, siehe fe5a979 und den Kommentar an signoff.go:239: 'Rejecting work without recording what is left undone is how the reason for a ticket gets lost'. Die Entscheidung unterscheidet 'angenommen' von 'braucht eigene Arbeit'. Sie kennt den haeufigsten Fall nicht: 'nicht fertig, die Arbeit geht auf DIESEM Ticket weiter'. Genau das koennen alle Agenten-Lanes - critique, optimize und testing schicken nach in-progress zurueck und wiederholen bis zur Stille. Die Schleife existiert fuer Agenten und fehlt dem Menschen, dessen Ablehnung die einzige ist, die niemand umgehen darf.

  Ein Folgeticket ist hier der falsche Zug, und es widerspricht der Regel dieses Boards, Befunde in das Ticket zu falten, das in der Hand liegt, statt Satellitentickets zu erzeugen. Ein unfertiges Ergebnis ist kein neues Thema.

  Heute bleibt nur 'jaira move <id> --to in-progress --force' aus einem Terminal ohne Agenten-Erkennung - also der Gate mit Gewalt an der Stelle, an der er am wenigsten umgangen werden sollte.
definition-of-done: "Ein Mensch kann ein Ticket aus human und signoff in jede Lane VOR der aktuellen schicken: der move-Dialog der Doska setzt Interactive, und der Gate an core/gate/gate.go:389 laesst diesen Ausgang durch. Ohne --force."
tags:
  - gates
  - tui
blocked-by: []
related: []
commits:
  - 8252caf76b52960f85291d904eda18adbd2ee6b1
  - 810f63cf5866372833c7e16456b96f5aeb56239b
  - 029d1e78625e3dafd48a309fabfe5033b1b230f7
  - 3b4dd01ac601c39ee28e2c8b80d3acb6a5db01f6
  - 871f0068f3b3c7505aff795e060ef610ee6e06a0
created-at: 2026-09-18T11:08:36Z
updated-at: 2026-11-22T19:32:22Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-910660
claimed-at: 2026-11-22T18:49:11Z
outcome-what: "Die Doska schickt ein Ticket aus signoff/human mit m in eine fruehere Lane zurueck und fragt dabei, was offen ist; die Antwort landet als Notiz im selben Schreibvorgang (move.Request.Note). Fusszeile nennt 'm send back'. Die human-Lane-Prosa sagt nicht mehr 'never out'."
outcome-why: "Der Gate erlaubt den Rueckweg schon seit Q484YS, aber ohne Begruendung ging der Grund verloren (fe5a979), und der signoff-Schirm nannte den Weg nicht."
outcome-resolves: "Mensch muss nicht mehr --force oder Folgeticket nehmen, um Nacharbeit auf demselben Ticket anzustossen."
review-summary: "Ein Mensch kann ein Ticket aus signoff oder human in der Doska mit 'm' in eine fruehere Lane zurueckschicken. Ist das Ziel frueher als die aktuelle Lane, oeffnet sich ein Prompt 'what is still open'; leer bleibt er offen, esc bricht ab, enter verschiebt und schreibt im selben Schreibvorgang die Notiz 'sent back to <lane> from <lane>: <Grund>' (neues Feld move.Request.Note). Der Prompt merkt sich Ticket-ID, Ausgangslane und Rueckweg (struct sendBack), damit ein Reload alle 2 s weder ein anderes Ticket verschiebt noch einen inzwischen weitergewanderten Ticket bewegt; eine Meldung ueber dem Prompt kehrt zum Prompt zurueck, und das Schliessen setzt mode und returnTo zurueck. Der Gate selbst ist unveraendert (seit Q484YS lehnt gate.go:423 nur vorwaerts ab); ein Test pinnt beide Richtungen. Die Fusszeile des signoff-Schirms nennt 'm send back', die Lane-Prosa fuer human sagt nicht mehr 'never out', und core/release/notes/BQGM6X.md hat zwei Zeilen."
round: 4
back-from: testing
back-by: lane
back-count: 2
review-gaps: |-
  note: pre-existing: openMove (internal/tui/model.go:2021) und applyMove (model.go:2077) nehmen m.selected(), den Board-Cursor, nicht m.detail und keine beim Oeffnen gemerkte ID - ein Reload zwischen 'm' und enter im Lane-Picker kann ein anderes Ticket treffen; der neue Prompt pinnt erst danach. Fix: Ticket-ID beim openMove festhalten wie sendingBack.
  note: DoD 4 gilt nur in der Doska - ein 'jaira move' aus signoff zurueck ueber die CLI schreibt keinen Grund (bewusste Entscheidung, Notiz 18:50); bekannte Grenze.
  note: pre-existing: notify ueber modeReject/modeCreate faellt weiter auf das Board mit halb getipptem m.input zurueck (model.go:1976) - modeSendBack ist jetzt behandelt, die beiden anderen Prompts nicht.
test-verdict: "pass: go vet + go test ./... gruen (RC=0) bei 3b4dd01ac, -race auf internal/tui gruen; DoD 1-6 im Baum geprueft; Verhalten per TUI-Tests: Rueckgabe mit Grund, leerer Grund haelt Prompt, esc bewegt nichts, m nach done abgelehnt, zwischendurch verschobenes Ticket wird nicht bewegt, Meldung ueber dem Prompt kehrt zu ihm zurueck und endet sauber"
hunt-verdict: "pass: 3 Runden; Runde 1 zwei change-Befunde (Prompt haelt Ticket nicht, notify schliesst Prompt), Runde 2 ein change-Blocker (nil-Panic nach Meldung), alle behoben; Runde 3 bei 3b4dd01ac ohne Befund"
review-verdict: "accept: alle sechs DoD-Punkte durch Code und Tests belegt; go test ./internal/tui ./core/move ./core/board ./core/gate gruen; nur pre-existing Befunde und eine dokumentierte Grenze"
review-round: 4
review-verdict-gaps: |-
  note: pre-existing: openMove (internal/tui/model.go:2021) und applyMove (model.go:2077) nehmen m.selected(), den Board-Cursor, nicht m.detail und keine beim Oeffnen gemerkte ID - ein Reload zwischen 'm' und enter im Lane-Picker kann ein anderes Ticket treffen; der neue Prompt pinnt erst danach. Fix: Ticket-ID beim openMove festhalten wie sendingBack.
  note: DoD 4 gilt nur in der Doska - ein 'jaira move' aus signoff zurueck ueber die CLI schreibt keinen Grund (bewusste Entscheidung, Notiz 18:50); bekannte Grenze.
  note: pre-existing: notify ueber modeReject/modeCreate faellt weiter auf das Board mit halb getipptem m.input zurueck (model.go:1976) - modeSendBack ist jetzt behandelt, die beiden anderen Prompts nicht.
review-check: "1. go build -o /tmp/jaira ./cmd/jaira  2. in einem Testrepo ein Ticket nach signoff bringen, /tmp/jaira starten, das Ticket mit enter oeffnen  3. Fusszeile zeigt 'm send back'  4. m druecken, in-progress waehlen, enter  5. Statuszeile zeigt 'send back <ID> to in-progress — what is still open:'; enter auf leerem Text tut nichts  6. 'Tests fehlen' tippen, enter  7. Ticket steht in in-progress, unter Progress steht 'sent back to in-progress from signoff: Tests fehlen'  8. nochmal nach signoff, m, done waehlen, enter -> Ablehnung 'human checkpoint' (nur a nimmt an). Automatisch: go test ./internal/tui -run 'SendBack|HumanExitGate'"
---

# Aus einer Menschen-Lane kann ein Mensch nur vorwaerts oder zur Seite, nicht zurueck

## Definition of Done

- [x] Der Gate an core/gate/gate.go:389 greift nur noch fuer einen Ausgang VORWAERTS: ein Ziel, das in der Lane-Reihenfolge nach der aktuellen liegt, bleibt abgelehnt. Ein Ziel VOR der aktuellen wird durchgelassen - ohne --force und ohne Interactive, also auch fuer einen Agenten. Die Richtung entscheidet, nicht der Actor.
  proof: core/gate/gate.go:423 (seit Q484YS 8cea1b42e); TestHumanExitGate
- [x] Vorwaerts bleibt allein die accept-Taste: ein move in die terminale Lane oder in irgendeine Lane nach der aktuellen wird weiter abgelehnt, aus der Doska wie aus der CLI. Mit Test fuer beide Richtungen.
  proof: TestHumanExitGate (CLI ohne Interactive und ByPerson nach done abgelehnt); TestMoveDialogCannotAcceptOutOfSignoff
- [x] Ein Agent gewinnt nichts, was er nicht schon hat: er darf zurueckschicken, aber nicht annehmen. Ein Test pinnt, dass ein Vorwaerts-move aus der CLI abgelehnt bleibt, und dass der Rueckwaerts-move ohne Interactive durchgeht.
  proof: TestHumanExitGate: To done ohne Interactive -> CodeNeedsHuman, To in-progress ohne Interactive -> frei
- [x] Der Grund fuer die Rueckgabe wird festgehalten, damit die Begruendung von fe5a979 gewahrt bleibt: die Rueckgabe schreibt eine Notiz mit der Ziel-Lane und dem, was der Mensch offen sieht, statt stillschweigend die Lane zu wechseln.
  proof: internal/tui/model.go moveTo + confirmSendBack + sendsBackFromPerson; core/move/move.go Request.Note; TestSendBackFromSignoffWritesTheReason, TestSendBackDoesNotMoveATicketThatMovedMeanwhile
- [x] Die Fusszeile des signoff-Schirms nennt die neue Aktion, und core/release/NOTES.md hat eine Zeile unter ## Unreleased.
  proof: internal/tui/signoff.go footer 'm send back'; core/release/notes/BQGM6X.md (NOTES.md ist geschlossen)
- [x] core/board/announce.go:210 beschreibt eine Lane mit requires-question nicht mehr als eine, aus der niemand Arbeit herausbewegen darf: human traegt kein requires-human-exit, ein Agent darf sie verlassen, und der erzeugte Lane-Block muss sagen, was der Gate tut. Heute sagt die Prosa 'you may move work in, never out' und der Gate laesst den Ausgang zu.
  proof: core/board/announce.go case f.Question; core/board/lanesection_test.go; CLAUDE.md/AGENTS.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 18:49 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 18:50 · Alexander Sacharov** — decision: DoD 1-3 sind seit Q484YS (8cea1b42e, 37208bf84) in master: gate.go:423 lehnt nur 'to > from' ab, ohne Interactive. Hier bleibt fuer sie nur ein Test-Check, kein Code.
- **2026-11-22 18:50 · Alexander Sacharov** — decision: DoD 4 nur in der Doska. Der move-Dialog 'm' fragt nach einer Begruendung, wenn ein Mensch ein Ticket aus einer requires-human-exit-Lane in eine fruehere Lane schickt; leer laesst den Prompt offen (wie reject). Die Notiz lautet 'sent back to <lane>: <text>'. Die CLI bekommt keine Pflicht: Agenten schreiben heute note+move, und eine neue Pflicht wuerde bestehende Aufrufe aendern. Keine eigene Taste: 'b' ist in der Detailansicht schon belegt (blocked-Kette), und 'm' kann es bereits.
- **2026-11-22 18:50 · Alexander Sacharov** — decision: DoD 5 Release-Zeile in core/release/notes/BQGM6X.md statt core/release/NOTES.md - NOTES.md ist seit AR97K9 geschlossene Geschichte.
- **2026-11-22 18:55 · Alexander Sacharov** — critique: applyMove bekommt den Grund als Parameter statt ueber das Uebergangsfeld m.sendBackWhy; sonst nichts - sendsBackFromPerson dupliziert move.isPersons nur teilweise (ohne !Agentic) und bleibt, weil isPersons unexportiert ist und ein agentischer Lauf hier nie gefragt werden soll.
- **2026-11-22 18:56 · Alexander Sacharov** — critique pass 2: der einzige Befund (sendBackWhy-Feld) ist durch applyMove(why) erledigt; nichts weiter.
- **2026-11-22 18:58 · Alexander Sacharov** — bughunt: change: blocker: internal/tui/model.go:2059 (applyMove takes m.selected() at enter) with internal/tui/model.go:1377-1384 (modeSendBack) - class: the send-back prompt pins nothing it was opened for — neither the ticket nor the lane it left; the tick reload every 2s (model.go:1110) re-resolves both while the person types, unlike reject, which pins m.rejecting. example 1: two tickets A and B in signoff, prompt open on A, another session deletes/files/moves A to another board, the reload clamps the cursor onto B, enter sends B (never opened, the detail pane still shows A) back to in-progress with A's reason as its note. example 2: A is moved signoff->review by someone else while the reason is typed; at enter sendsBackFromPerson(review,in-progress) is false, so A is sent back with no note and the typed reason is dropped silently. ran: scratch test in internal/tui (Store.Delete A + m.reload, then enter: B status=in-progress, note carries A's reason; Mutate A status=review + reload, then enter: A in-progress, no note, no message)
- **2026-11-22 18:58 · Alexander Sacharov** — bughunt: change: internal/tui/model.go:1937-1950 (notify) reached from the new modeSendBack - class: a background message raised while a text prompt is open takes the prompt down; notify files modeSendBack under default, so returnTo becomes modeBoard, not the ticket screen; the half-typed reason stays in m.input and the next keys act on the board (q dismisses, then letters are board commands). The same holds pre-existing for modeReject and modeCreate. example: prompt open, typed 'half a re', refFetchedMsg{err} arrives (model.go:1135-1142), then q: mode=board, input='half a re', ticket unmoved. ran: scratch test in internal/tui (Update(refFetchedMsg{err}) then q)
- **2026-11-22 18:58 · Alexander Sacharov** — bughunt: done at 810f63cf5 - 2 findings (2 change, 1 of them blocker) in 4 modules, first hunt; concurrency, error paths, input, permissions checked; files: AGENTS.md CLAUDE.md core/board/announce.go core/board/lanesection_test.go core/move/move.go core/release/notes/BQGM6X.md internal/tui/model.go internal/tui/move_test.go internal/tui/render_test.go internal/tui/returnto_test.go internal/tui/sendback_test.go internal/tui/signoff.go internal/tui/signoff_test.go internal/tui/team_test.go internal/tui/view.go
- **2026-11-22 18:59 · Alexander Sacharov** — testing: bughunt change-Befunde (Prompt haelt Ticket/Lane nicht; notify schliesst Prompt) -> zurueck nach in-progress. Fix: sendBack-Zustand mit id, from, target und eigenem back-Modus.
- **2026-11-22 19:01 · Alexander Sacharov** — critique pass 3: nur der Fix 029d1e786 gelesen - sendBack-Struktur folgt dem rejecting-Muster (model.go), nichts zu sagen.
- **2026-11-22 19:03 · Alexander Sacharov** — bughunt: change: blocker: internal/tui/model.go:1972 (notify sets returnTo=modeSendBack) + model.go:2306 (finishMove: m.mode = m.returnTo) - class: returnTo still says modeSendBack after the prompt has closed and sendingBack is nil, so any later 'm.mode = m.returnTo' reopens a prompt with no send-back behind it, and esc (model.go:1398 m.sendingBack.back) or enter (confirmSendBack model.go:2087 sb.back) dereferences nil and crashes the TUI; the esc-from-prompt path (model.go:1397) leaves the same stale returnTo behind; example: open a signoff ticket, m -> in-progress, type a reason, any message arrives over the prompt (fetch/reload error) and is dismissed, enter -> the ticket is moved, finishMove puts the board back into modeSendBack, next esc or enter panics with nil pointer dereference; ran: scratch test in internal/tui (signoffModel, pickLane, typeInto, notify, esc, enter, then esc / typeInto+enter) - status=in-progress, mode==modeSendBack, sendingBack=nil, both keys panic; blame: 029d1e786 fix(BQGM6X)
- **2026-11-22 19:03 · Alexander Sacharov** — bughunt: done at 029d1e786 - 1 finding (1 change, blocker) in 2 modules, follow-up hunt since 810f63cf5 (diff of internal/tui/model.go and internal/tui/view.go plus their tests; no new files); concurrency, error paths, input, permissions checked; files: AGENTS.md CLAUDE.md core/board/announce.go core/board/lanesection_test.go core/move/move.go core/release/notes/BQGM6X.md internal/tui/model.go internal/tui/move_test.go internal/tui/render_test.go internal/tui/returnto_test.go internal/tui/sendback_test.go internal/tui/signoff.go internal/tui/signoff_test.go internal/tui/team_test.go internal/tui/view.go
- **2026-11-22 19:04 · Alexander Sacharov** — testing round 2: bughunt change/blocker - notify ueber dem Prompt liess returnTo=modeSendBack stehen, finishMove oeffnete den Prompt leer, naechste Taste nil-Panic. Fix: leaveSendBack setzt mode und returnTo zurueck; nil-Guard im modeSendBack-Fall; TestSendBackAfterAMessageFinishesCleanly.
- **2026-11-22 19:06 · Alexander Sacharov** — bughunt: done at 3b4dd01ac - 0 findings in 1 module, follow-up hunt since 029d1e786 (diff of internal/tui/model.go and internal/tui/sendback_test.go: leaveSendBack, the nil-sendingBack guard, confirmSendBack; callers of moveTo, notify, finishMove, reload and the background message paths followed; no new files); concurrency, error paths, input, permissions checked; ran: go test ./internal/tui/ -run 'SendBack|Signoff|Message' -count=3 -race (ok); files: AGENTS.md CLAUDE.md core/board/announce.go core/board/lanesection_test.go core/move/move.go core/release/notes/BQGM6X.md internal/tui/model.go internal/tui/move_test.go internal/tui/render_test.go internal/tui/returnto_test.go internal/tui/sendback_test.go internal/tui/signoff.go internal/tui/signoff_test.go internal/tui/team_test.go internal/tui/view.go
