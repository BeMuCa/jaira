---
id: 01M81V2D6K5ZVMY4E9JBS1SA89
title: "Drei Altfehler aus den 1MYFKE-Hunts: Epic-Lease nach ErrRaceLost, CRLF-AGENTS.md driftet immer, TUI verschluckt Fetch-Fehler"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die drei Altfehler sind behoben, jeder mit einem Test, der ohne den Fix fehlschlägt."
context: |-
  Drei Altfehler, die die Bughunts an 1MYFKE am 22.11.2026 gefunden haben. Sie stammen nicht von 1MYFKE; Alex will sie behoben haben.
  - core/refsync/epic.go:121-137: rebaseEpic hält nach ErrRaceLost den Lease nicht richtig (Rennen zwischen zwei Schreibern auf den Epic-Ref).
  - core/board/announce.go: eine AGENTS.md mit CRLF-Zeilenenden gilt immer als abgewichen (drifted), auch direkt nach jaira update.
  - internal/tui/refs.go:107: die TUI verschluckt den RecordError aus dem Hintergrund-FetchMilestones; ein fehlgeschlagener Fetch wird nirgends gezeigt.
  - Details stehen in den Notizen von 1MYFKE (review-gaps, Notiz "pre-existing, from the hunts").
definition-of-done: "rebaseEpic hält nach ErrRaceLost den Lease korrekt (core/refsync/epic.go, Test)"
tags:
  - concurrency
blocked-by: []
related: []
commits:
  - 844ade861d6a8b6d2b28cb41126a27ec7c369f8a
  - dfff60c36adbe8c32e23629f4fd894929eb05033
  - 9659685475da23e74cb521f540049f338cfe4d50
  - e3d736ef62af155e7023adf9ec7e1b834154b37c
  - dc514243cf9b6166e99f118f3942cd1fd30b8e93
created-at: 2026-11-22T14:37:54Z
updated-at: 2026-11-22T16:19:25Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3557981
claimed-at: 2026-11-22T14:38:23Z
outcome-what: "rebaseEpic queues an unsent merge on the lease it was merged onto; the agent-file block follows the file's CRLF line endings; the board warns when a fetched milestone could not be queued"
outcome-why: "three pre-existing defects from the 1MYFKE hunts: a teammate's epic change overwritten on the next flush, CRLF AGENTS.md always drifted, a milestone silently never reaching the remote"
outcome-resolves: S1SA89
review-summary: "Drei Altfehler behoben. (1) core/refsync/epic.go rebaseEpic: ein Epic-Edit, der alle Runden gegen fremde Schreiber verliert, wird jetzt mit dem Lease des Standes wieder eingereiht, auf den er gemergt wurde (oursLease), statt mit EpicSHA nach dem Fetch - der nächste Flush verliert dann gegen den unbekannten Schreiber und merged erneut, statt dessen Mitglied zu überschreiben. (2) core/board: neuer Helper lineEnding (CRLF, wenn die Mehrheit der Zeilen CRLF ist); managedBlock baut den jaira-Block in den Zeilenenden der Zeilen außerhalb des Blocks, der Anhängepfad in announceInAgentFile ebenso, appendEntry nutzt denselben Helper - eine CRLF-AGENTS.md gilt nach jaira update nicht mehr als drifted. (3) TUI: ein milestone.RecordError aus dem Hintergrund-FetchMilestones kommt als refFetchedMsg.warn an und erscheint in der Statuszeile (flash, kein Moduswechsel) mit demselben Wortlaut wie jaira fetch (RecordError.Warning, von CLI und TUI geteilt). Release-Note S1SA89.md mit drei Zeilen."
round: 4
back-from: testing
back-by: lane
back-count: 2
review-gaps: |-
  note: Text hinter <!-- jaira:local --> behält seine eigenen Zeilenenden (announce.go:448), ein CRLF-Block mit LF-Lokalregeln von einem alten Build bleibt gemischt; Drift-Check bleibt konsistent (dieselbe Funktion), Fix: preserved per ReplaceAll auf nl normalisieren
  note: der CRLF-Pfad mit Versionszeile (release.Current != dev) ist ungetestet; versionLineRe nimmt \r?\n, funktioniert also, Fix: ein Test mit gesetztem release.Current
  note: pre-existing: removeBlock und setClaudeImport in core/board/announce.go nehmen weiter LF an (CRLF-CLAUDE.md: gemischte Leerzeilen bzw. @AGENTS.md-Import bleibt stehen)
  note: pre-existing: Flush hält keinen epics-Lock, rebaseEpic DropEpic+QueueEpic kann ein paralleles jaira epic add im selben Klon verschlucken (core/refsync/epic.go:119/140, core/outbox/outbox.go:471/481)
  note: pre-existing: der msg.err-Zweig von refFetchedMsg (internal/tui/model.go:1134) nutzt noch m.notify und nimmt im Hintergrund dem Tippenden die Tasten
test-verdict: "pass: round 3 - go vet ./core/board, go test -race ./core/board exit 0, internal/cli agent-file tests exit 0; DoD 2 re-verified: TestAnLFBlockInACRLFFileBecomesCRLF fails without the fix and passes with it"
hunt-verdict: "pass: bughunt done at dc514243c, 0 blocker; 1 change non-blocker (jaira:local text keeps its own line endings inside a converted block) left as the hunter's note for the person; earlier pre-existing: other agent-file writers LF-only, epic outbox drop vs concurrent add - notes for the person"
review-verdict: "accept: alle vier DoD-Punkte erfüllt, jeder Fix mit Test, der Epic-Test fällt nachweislich ohne den Fix (gegen 180daf772 geprüft); keine Blocker, nur Randnotizen"
review-round: 4
review-verdict-gaps: |-
  note: Text hinter <!-- jaira:local --> behält seine eigenen Zeilenenden (announce.go:448), ein CRLF-Block mit LF-Lokalregeln von einem alten Build bleibt gemischt; Drift-Check bleibt konsistent (dieselbe Funktion), Fix: preserved per ReplaceAll auf nl normalisieren
  note: der CRLF-Pfad mit Versionszeile (release.Current != dev) ist ungetestet; versionLineRe nimmt \r?\n, funktioniert also, Fix: ein Test mit gesetztem release.Current
  note: pre-existing: removeBlock und setClaudeImport in core/board/announce.go nehmen weiter LF an (CRLF-CLAUDE.md: gemischte Leerzeilen bzw. @AGENTS.md-Import bleibt stehen)
  note: pre-existing: Flush hält keinen epics-Lock, rebaseEpic DropEpic+QueueEpic kann ein paralleles jaira epic add im selben Klon verschlucken (core/refsync/epic.go:119/140, core/outbox/outbox.go:471/481)
  note: pre-existing: der msg.err-Zweig von refFetchedMsg (internal/tui/model.go:1134) nutzt noch m.notify und nimmt im Hintergrund dem Tippenden die Tasten
review-check: "1. go test ./core/refsync -run TestAnEpicThatKeepsLosing  -> ok (Remote-Hook lässt jeden Push verlieren; am Ende hält das Remote-Epic alle Mitglieder)  2. go test ./core/board -run 'CRLF|StrayCRLF|LFBlockInACRLF' -v  -> drei Tests PASS  3. Handprobe: go build -o /tmp/jaira ./cmd/jaira, in einem Scratch-Repo mit Board: sed -i 's/$/\\r/' AGENTS.md, dann /tmp/jaira update und /tmp/jaira validate  -> validate meldet keine Drift, file AGENTS.md zeigt 'with CRLF line terminators'  4. go test ./internal/tui -run TestTheBoardSaysWhenAMilestoneCouldNotBeQueued  -> PASS (eine TUI-Handprobe bräuchte eine read-only Outbox; der Test baut genau das)"
---

# Drei Altfehler aus den 1MYFKE-Hunts: Epic-Lease nach ErrRaceLost, CRLF-AGENTS.md driftet immer, TUI verschluckt Fetch-Fehler

## Definition of Done

- [x] rebaseEpic hält nach ErrRaceLost den Lease korrekt (core/refsync/epic.go, Test)
  proof: core/refsync/epic.go rebaseEpic oursLease; TestAnEpicThatKeepsLosingIsRequeuedOnTheStateItWasMergedOnto (core/refsync/epic_test.go)
- [x] AGENTS.md mit CRLF gilt nach jaira update nicht als drifted (core/board/announce.go, Test)
  proof: core/board/announce.go managedBlock (line ending of the lines outside the block) + announceInAgentFile, core/board/gitignore.go lineEnding (majority); TestACRLFAgentFileIsNotDrifted, TestAStrayCRLFLineLeavesTheBlockInLF, TestAnLFBlockInACRLFFileBecomesCRLF (core/board/crlf_test.go)
- [x] Ein fehlgeschlagener Hintergrund-FetchMilestones wird in der TUI angezeigt statt verschluckt (internal/tui/refs.go, Test)
  proof: internal/tui/refs.go refFetchedMsg.warn + internal/tui/model.go refFetchedMsg flash (status line, no mode change); TestTheBoardSaysWhenAMilestoneCouldNotBeQueued (internal/tui/recorderror_test.go)
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256 falls nötig
  proof: core/release/notes/S1SA89.md; shipped.sha256 unverändert (keine installierte Datei geändert, scripts/shipped-history.sh --check ok)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 14:38 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 14:40 · Alexander Sacharov** — Epic-Lease: Fix = der Text 'ours' wird mit dem Lease des Standes eingereiht, auf den er gemergt wurde (oursLease), nicht mit EpicSHA nach dem Fetch. Test: pre-receive-Hook, der den Ref bei jedem Push weiterschiebt und ablehnt; Hook muss GIT_OBJECT_DIRECTORY/QUARANTINE unsetten, sonst verschwinden seine Objekte.
- **2026-11-22 14:44 · Alexander Sacharov** — CRLF: Regel = Datei enthält irgendwo CRLF -> Block in CRLF (wie gitignore.go appendEntry). Auch der appended-Pfad schreibt CRLF, sonst driftet ein frisch angehängter Block sofort. TUI: RecordError wird als refFetchedMsg.warn durchgereicht und mit dem Wortlaut von 'jaira fetch' per notify gezeigt; Test baut den echten Fehler (alte backend.md + Outbox read-only).
- **2026-11-22 14:44 · Alexander Sacharov** — critique: (1) CRLF-Erkennung dreimal (announce.go x2, gitignore.go:442) -> ein Helper lineEnding(s) in core/board. (2) Warnsatz in TUI und CLI doppelt -> (*milestone.RecordError).Warning() und beide nutzen. epic.go: direktes Box.QueueEpic statt RecordEpic bleibt (RecordEpic ist Teil des epic.Refs-Interfaces, eine Lease-Variante dort wäre mehr Fläche als zwei Zeilen).
- **2026-11-22 14:45 · Alexander Sacharov** — critique Runde 2: beide Befunde erledigt (lineEnding in gitignore.go, RecordError.Warning in write.go, CLI-Ausgabe unverändert); nichts Neues.
- **2026-11-22 14:48 · Alexander Sacharov** — testing: Funktion geprüft mit gebautem Binary auf Scratch-Repo: AGENTS.md per sed auf CRLF -> validate meldet keine Drift, update lässt 182/182 Zeilen CRLF, zwei parallele update -> weiterhin 182/182 CRLF, validate sauber. Epic- und TUI-Pfad über die echten Tests (Hook-Remote bzw. read-only Outbox).
- **2026-11-22 14:49 · Alexander Sacharov** — bughunt: change: internal/tui/model.go:1144 - class: a background fetch result shown with m.notify, which switches the board into modeMessage and takes the screen from whatever the person is doing (announceArrivals in refs.go says such news must be a flash line, never a screen); same pattern in the pre-existing msg.err branch at model.go:1134; example: typing a new ticket title in modeCreate when a refFetchedMsg with warn arrives - mode becomes modeMessage, returnTo modeBoard (notify's default branch, model.go:1924), and the Enter meant to submit the ticket only dismisses the warning and lands on the board, the ticket is not created; the same hits modeEdit, modeFilter, modeDelete, modeReject, modeMilestones; in modeMessage with a pending refused-move offer the warning replaces its text while m.pending stays armed for f/y; not a blocker: only when queueing a milestone fails; ran: scratch test in internal/tui (mode create=5 -> 6 after the warn, Enter -> 0, removed afterwards)
- **2026-11-22 14:49 · Alexander Sacharov** — bughunt: change: core/board/gitignore.go:432 (lineEnding) used by core/board/announce.go:437 and :733 - class: the whole managed block takes CRLF when the file holds a single CRLF anywhere, inside or outside the block; example: an LF CLAUDE.md with one pasted CRLF line, or a local area saved by a Windows editor, has its whole block rewritten to CRLF by 'jaira update' - a full-block diff, and a teammate on an older build sees it drifted and rewrites it to LF, so the block flips with every update between builds; not a blocker: the result is consistent with the drift check (no update loop on one build); ran: scratch test in core/board (mixed file -> block in CRLF, NoteIsCurrent true, removed afterwards)
- **2026-11-22 14:49 · Alexander Sacharov** — bughunt: pre-existing: core/board/announce.go:600-613 (removeBlock) and :638-661 (setClaudeImport) - class: the other agent-file writers still assume LF, the same class this change fixed in managedBlock and the append path; example 1: removeBlock on a CRLF CLAUDE.md leaves '# P\r\n\r\n\n\r\n\r\nAfter.' (TrimRight/TrimLeft on "\n" stop at the \r, a bare \n is inserted); example 2: setClaudeImport(want=false) on a CRLF CLAUDE.md does nothing (strings.Contains with the LF claudeImport fails), so the @AGENTS.md import stays after AGENTS.md is dropped; example 3: setClaudeImport(want=true) prepends the import in LF to a CRLF file; blame: aff2304a20 (before the ticket); ran: scratch test in core/board, removed afterwards
- **2026-11-22 14:49 · Alexander Sacharov** — bughunt: pre-existing: core/refsync/epic.go:119 and :140-142, core/outbox/outbox.go:471/481 - class: an outbox entry is dropped or replaced after its send without checking it is still the entry that was read, and Flush holds no 'epics' lock while every epic edit does (core/epic/epic.go:768); example: the board's background Flush is in rebaseEpic having read the pending epic, a CLI 'jaira epic add' in the same clone queues a new member on that entry, then rebaseEpic's DropEpic + QueueEpic (or the success-path DropEpic) removes it - the member is lost with no copy anywhere, since an epic has no local file; blame: 9cb410a730/ae8ef91a50 (F914QH, before the ticket); this change only altered the lease of the requeue, not the window; confirmed by read
- **2026-11-22 14:49 · Alexander Sacharov** — bughunt: done at 965968547 - 4 findings (2 change, 2 pre-existing) in 6 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/board/announce.go core/board/crlf_test.go core/board/gitignore.go core/milestone/write.go core/refsync/epic.go core/refsync/epic_test.go core/release/notes/S1SA89.md internal/cli/milestones.go internal/tui/model.go internal/tui/recorderror_test.go internal/tui/refs.go
- **2026-11-22 14:50 · Alexander Sacharov** — testing: class: Warnung als modale Meldung in einem Hintergrundpfad; example: m.notify im refFetchedMsg-Zweig schaltet in modeMessage, Enter beim Ticket-Anlegen schließt nur die Meldung; fix: Warnung in die Statuszeile (flash) wie announceMilestones, ohne Moduswechsel. class: Zeilenende-Erkennung per 'irgendein CRLF'; example: eine CRLF-Zeile in LF-AGENTS.md schaltet den Block auf CRLF; fix: lineEnding nach Mehrheit der Zeilen, Test mit einer verirrten CRLF-Zeile.
- **2026-11-22 14:51 · Alexander Sacharov** — Runde 2: Warnung jetzt in der Statuszeile (flash, hinter den News desselben Fetch), kein Moduswechsel. lineEnding nach Mehrheit der Zeilen. Beide Tests fallen ohne Fix. Altfunde 3 (removeBlock/setClaudeImport nur LF) und 4 (Epic-Outbox-Drop vs. paralleles epic add) bleiben Notizen für die Person.
- **2026-11-22 14:53 · Alexander Sacharov** — bughunt: change: core/board/gitignore.go:434 (lineEnding, majority rule from e3d736ef6) via core/board/announce.go:437 managedBlock - class: the line ending is voted on by the whole file, including the jaira block jaira itself wrote, so a CRLF file whose block an older build appended in LF keeps that LF block forever: the update rewrites it in LF, and NoteIsCurrent then calls the mixed file current, so nothing repairs it (the 'any CRLF' rule at 965968547 did). Same vote in appendEntry gitignore.go:453 for .gitignore/info/exclude, smaller since an entry is 2 lines; example: CRLF '# Project / Some notes' (4 lines) + LF block of 9 lines -> action updated, LF lines remain, current=true. Possible fix: count only the lines outside the managed block. Not a blocker: no data loss, the result is stable, a fresh checkout with autocrlf makes the file CRLF again; ran: go test -overlay with a scratch test in core/board (FAIL: mixed endings remain after update)
- **2026-11-22 14:53 · Alexander Sacharov** — bughunt: done at e3d736ef6 - 1 finding (1 change, 0 pre-existing) in 3 modules, follow-up hunt since 965968547 (diff only: core/board/gitignore.go core/board/crlf_test.go internal/tui/model.go internal/tui/recorderror_test.go); concurrency, error paths, input, permissions checked; files: core/board/announce.go core/board/crlf_test.go core/board/gitignore.go core/milestone/write.go core/refsync/epic.go core/refsync/epic_test.go core/release/notes/S1SA89.md internal/cli/milestones.go internal/tui/model.go internal/tui/recorderror_test.go internal/tui/refs.go
- **2026-11-22 14:53 · Alexander Sacharov** — testing: class: Zeilenende wird über den eigenen Block mitgezählt; example: 4 CRLF-Nutzerzeilen + 9 LF-Blockzeilen -> Block bleibt LF, current=true; fix: Mehrheit nur über die Zeilen außerhalb des Blocks, ohne solche Zeilen über die ganze Datei; Test mit genau diesem Fall.
- **2026-11-22 14:55 · Alexander Sacharov** — bughunt: change: core/board/announce.go:448 - class: the text a person keeps behind <!-- jaira:local --> is copied into the block in its own line endings, while the rest of the block now takes the line ending of the lines outside it (announce.go:443), so the two can differ and the block ends up mixed - the thing managedBlock's comment says it prevents; example: a CRLF AGENTS.md whose block an older build wrote in LF with two local rules after the marker -> after jaira update the block has 142 CRLF lines and the local rules stay '\nmy rule 1\nmy rule 2' in LF; NoteIsCurrent says current and a second update is byte-stable, so no drift loop and no lost text, an edge, not a blocker; the reverse (LF file, CRLF local text) is the same class; ran: scratch test in core/board calling announceInAgentFile on that file (removed afterwards)
- **2026-11-22 14:56 · Alexander Sacharov** — bughunt: done at dc514243c - 1 finding (1 change, 0 pre-existing) in 1 module, follow-up hunt since e3d736ef6 (diff only: core/board/announce.go core/board/crlf_test.go); concurrency, error paths, input, permissions checked; files: core/board/announce.go core/board/crlf_test.go core/board/gitignore.go core/milestone/write.go core/refsync/epic.go core/refsync/epic_test.go core/release/notes/S1SA89.md internal/cli/milestones.go internal/tui/model.go internal/tui/recorderror_test.go internal/tui/refs.go
