---
id: 01M865XBRDFR99ENKNMKWV2J2W
title: "outbox: eine Schreibung, die waehrend eines Pushs fuer denselben Schluessel eingereiht wird, wird geloescht"
status: done
ready: true
creator: Alexander Sacharov
goal: "Ein Outbox-Eintrag wird nach dem Senden nur geloescht, wenn er noch derselbe ist, der gelesen wurde; eine inzwischen neu eingereihte Schreibung bleibt erhalten."
context: |-
  Was falsch ist: Ein Outbox-Eintrag wird nach einem langsamen Schritt (Push, git-Lesen) per Schluessel geloescht, ohne zu pruefen, ob er noch der gelesene ist. Kein Lock wird gehalten.
  Folge: eine Schreibung, die eine zweite Session in diesem Fenster fuer denselben Schluessel einreiht, wird geloescht.
  Stellen: core/refsync/epic.go rebaseEpic (DropEpic nach WriteEpic); core/outbox/outbox.go:471 Box.Flush (dropKind nach dem Senden, alle Arten); core/refsync/refsync.go tidyMilestoneEntries (DropMilestone nach ReadMilestone, ausserhalb des Milestone-Locks, den mergeRefused nimmt).
  Schaden: Epic - der Eintrag ist die einzige Kopie der Aenderung, sie ist weg. Milestone - die Datei gilt nicht mehr als ausstehend, der naechste Fetch ueberschreibt sie. Ticket - die Unsent-Markierung verschwindet, die Aenderung wartet bis zur naechsten Schreibung.
  Beispiel: 'jaira epic add' in einer zweiten Session, waehrend der Flush der ersten dasselbe Epic pusht.
  Nachgewiesen: Bughunt zu MMKENA am 24.11.2026 mit einem pre-receive-Hook, der mitten im Push eine zweite Epic-Schreibung einreiht - sie war danach weder in der Outbox noch auf dem Remote. Gab es vor MMKENA (d3bc99ef64, 9cb410a730, 90151d5877).
definition-of-done: "Ein Test reiht waehrend des Pushs eine zweite Schreibung fuer denselben Schluessel ein; sie ist danach noch in der Outbox oder auf dem Remote - fuer Ticket, Milestone und Epic"
tags:
  - concurrency
blocked-by: []
related: []
follows: 01M8543DKBZA0P1CHFGKMMKENA
commits:
  - b5220edeb4568e308c440409958797f5764b0b6b
  - cb71b1c7a9ce85b2d0a21a6afe3acf597730924c
  - 83a9db6d4a9cfd88c401aaa41a194c401e99bb7e
  - 6ee1875e9efad2c73eef16f9a3fa7e229c757bb6
  - 28a06ab4dd2a349e4f3136de146d1a0ed765106b
  - 92f14119a8bc2416e32a4f762b75775a478e3cd0
created-at: 2026-11-24T07:04:21Z
updated-at: 2026-11-24T11:28:27Z
assignee: Alexander Sacharov
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-2599564
claimed-at: 2026-11-24T10:20:56Z
outcome-what: "outbox.Box nimmt ein Lock-File um jedes Lesen-dann-Schreiben; Settle(seen) und Replace(seen, ...) wirken nur auf den gelesenen Eintrag; Flush leased einen waehrend des Sendens eingereihten neueren Eintrag auf die gesendete SHA um - nach einem Delete auf 'kein Ref' (leere Lease)."
outcome-why: "Nach Signoff-Rueckgabe: ein waehrend eines Deletes eingereihter Eintrag behielt die Lease des geloeschten Refs und wurde beim naechsten Flush als Rejected verworfen. Settle(seen, \"\") trug zwei Bedeutungen; jetzt Settle(seen) = nie umleasen, settle(seen, &sha) nur in Flush."
outcome-resolves: "DoD 1-3: race_test.go (alle Arten, gesendet, abgelehnt, Delete) und epic_race_test.go; DoD 2 Release-Note."
review-summary: "Delta-Review ab cb71b1c7a (ohne den upstream-Merge 83a9db6d4): Box.Settle hat keine Lease mehr und laesst eine inzwischen neu eingereihte Schreibung samt Lease stehen (fuer Ablehnung, Merge, tidy, rebaseEpic). Nur Flush ruft nach erfolgreichem Senden das private settle(seen, &sha) auf; nach einem gesendeten Delete ist sha leer, und die waehrend des Pushs eingereihte Schreibung wird auf 'kein Ref' umgeleast - force-with-lease=<ref>: verlangt genau, dass der Ref fehlt. Damit ist der Edge der letzten Review (alte Lease nach Delete -> naechster Flush Rejected) geschlossen; neuer Test TestAWriteQueuedDuringItsDeleteIsLeasedOnNoRef."
review-gaps: |-
  note: elsewhere: refsync.go:157-161, :222-226 und epic.go:67-71 - Lease wird vor dem Outbox-Lock gelesen (vorbestehend; Folgeticket AEZ3HD)
  note: elsewhere: core/identity TestCurrentHonoursJairaUserAboveEverything liest den echten people-Ref und ist rot, unabhaengig von dieser Aenderung (Folgeticket X2M486)
  note: elsewhere: gitref.go:548 publish - schlaegt update-ref nach erfolgreichem Push fehl, meldet Flush Failed und der Eintrag bleibt, obwohl er auf dem Remote ist; der naechste Flush verliert gegen die eigene Schreibung (vorbestehend, new ticket)
test-verdict: "pass: go test -count=1 ./core/... gruen bis auf core/identity TestCurrentHonoursJairaUserAboveEverything (vorbestehend, X2M486); neu TestAWriteQueuedDuringItsDeleteIsLeasedOnNoRef gruen, rot bei 6ee1875e9 ohne Fix"
hunt-verdict: "pass: stop-after-fix, no hunt in round 2 - settle aendert keine Sperre (Lock-Aufnahme unveraendert), nur welche Lease ein neuerer Eintrag bekommt"
review-verdict: "accept: DoD 1-3 erfuellt; Delete-Lease-Edge geschlossen, eine erfolgreiche Schreibung liefert nie eine leere SHA (publish), ein zweites Delete auf 'kein Ref' laeuft durch (git: exit 0); keine neuen Defekte"
review-round: 3
review-verdict-gaps: |-
  note: elsewhere: refsync.go:157-161, :222-226 und epic.go:67-71 - Lease wird vor dem Outbox-Lock gelesen (vorbestehend; Folgeticket AEZ3HD)
  note: elsewhere: core/identity TestCurrentHonoursJairaUserAboveEverything liest den echten people-Ref und ist rot, unabhaengig von dieser Aenderung (Folgeticket X2M486)
  note: elsewhere: gitref.go:548 publish - schlaegt update-ref nach erfolgreichem Push fehl, meldet Flush Failed und der Eintrag bleibt, obwohl er auf dem Remote ist; der naechste Flush verliert gegen die eigene Schreibung (vorbestehend, new ticket)
review-check: "1. cd /home/alex/projects/.worktrees/jaira-wv2j2w  2. go test ./core/outbox -run 'TestAWriteQueuedDuringItsDeleteIsLeasedOnNoRef|TestAWriteQueuedDuringThePushOfItsKeySurvivesTheFlush|TestALostRaceDropsOnlyTheRefusedWrite|TestSettleAndReplaceLeaveANewerEntryAlone' -v  3. alle vier PASS  4. go test ./core/refsync -run TestAnEpicEditQueuedDuringThePushReachesTheRemote -v  5. PASS. Kein UI-Pfad: das Rennen braucht zwei Sessions im selben Moment, die Tests stellen es mit einem Sender nach, der waehrend des Pushs einreiht."
round: 3
back-from: critique
back-by: lane
back-count: 1
stop-after-fix: "true"
---

# outbox: eine Schreibung, die waehrend eines Pushs fuer denselben Schluessel eingereiht wird, wird geloescht

## Definition of Done

- [x] Ein Test reiht waehrend des Pushs eine zweite Schreibung fuer denselben Schluessel ein; sie ist danach noch in der Outbox oder auf dem Remote - fuer Ticket, Milestone und Epic
  proof: core/outbox/race_test.go TestAWriteQueuedDuringThePushOfItsKeySurvivesTheFlush (ticket/milestone/epic), TestALostRaceDropsOnlyTheRefusedWrite; core/refsync/epic_race_test.go TestAnEpicEditQueuedDuringThePushReachesTheRemote (pre-receive-Hook, direkt + rebaseEpic); beide rot mit sameWrite=true
- [x] Release-Note in core/release/notes/WV2J2W.md
  proof: core/release/notes/WV2J2W.md
- [x] Nach einem gesendeten Delete wird eine waehrend des Pushs eingereihte Schreibung auf 'kein Ref' (leere Lease) umgeleast
  proof: core/outbox/race_test.go TestAWriteQueuedDuringItsDeleteIsLeasedOnNoRef; core/outbox/outbox.go settle (Flush leased auch auf leere SHA); rot bei 6ee1875e9 mit altem Settle(e, sha)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-24 07:13 · Alexander Sacharov** — teamlead 24.11.2026: Alex: starten. MMKENA (refsync.Flush, PR #246) wird gerade gelandet; vor dem PR upstream/master einmergen, Konflikt in core/refsync/refsync.go moeglich.
- **2026-11-24 07:14 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/3417.sock
- **2026-11-24 07:15 · Alexander Sacharov** — decision: Box bekommt ein Lock-File (<outbox>/.lock, ticket.LockFile) um jedes Lesen-Vergleichen-Schreiben: queueKind, RekeyMilestone und ein neues bedingtes Drop. Gleich heisst: Op, Content, File, Lease, Migration gleich - UpdatedAt nicht, denn derselbe Inhalt neu eingereiht ist nach dem Senden schon auf dem Remote. Nach erfolgreichem Senden und geaendertem Eintrag wird der neue Eintrag auf die gesendete SHA umgeleast (nicht geloescht): er enthaelt die gesendete Datei schon (ganze Datei, ersetzt auf derselben Basis); mit der alten Lease wuerde ihn der naechste Flush als Rejected verwerfen (Ticket). Ausnahme rebaseEpic: gesendet wurde ein Merge, den der neue Eintrag nicht kennt - dort bleibt seine alte Lease, der naechste Flush merged ihn 3-Wege.
- **2026-11-24 07:18 · Alexander Sacharov** — Gefunden: Box.Flush loeschte nach Sent per Schluessel - auch die Ticket-Unsent-Markierung. Ein neuer Eintrag mit alter Lease waere beim naechsten Flush als Rejected verworfen worden (Ticket) - darum Settle(seen, sha) mit Umleasen statt nur Stehenlassen. mergeRefused-Drop (Zeile ~956) laeuft unter dem Milestone-Lock, ist trotzdem auf Settle umgestellt (billig). mergeRefused Drop+RecordMilestone (~971) unveraendert: unter Milestone-Lock, und Box-Lock reicht dort nicht ueber zwei Aufrufe. Kein Nebenlaeufigkeits-Stresstest mit Goroutinen: Lock ist ticket.LockFile, dort schon getestet.
- **2026-11-24 07:19 · Alexander Sacharov** — waiting: go test ./... (Hintergrund, eigene Session) since 08:19
- **2026-11-24 07:21 · Alexander Sacharov** — critique: nichts. Gelassen: refsync.go mergeRefused Drop+RecordMilestone (laeuft unter dem Milestone-Lock, RecordMilestone holt die Lease vom Repo - Replace passt dort nicht ohne zweite Lease-Quelle). Replace schreibt auch, wenn der Eintrag inzwischen fehlt: wie vorher DropEpic+QueueEpic; im schlimmsten Fall verliert der Merge den naechsten Push und wird erneut gemerged, nichts geht verloren.
- **2026-11-24 07:23 · Sascha** — testing: note: pre-existing: core/identity TestCurrentHonoursJairaUserAboveEverything liest den echten people-Ref des Repos (refs/jaira/people, heute 08:22 mit 'Berk: BeMuCa, berkc' aktualisiert) und wird damit rot (Current()="Berk", want "berk"). Nicht diese Aenderung (identity unberuehrt); Test braucht ein isoliertes Repo.
- **2026-11-24 07:24 · Sascha** — waiting: bughunt (Hintergrund-Agent dieser Session) since 08:24
- **2026-11-24 07:27 · Sascha** — bughunt: pre-existing: core/refsync/refsync.go:157-161 - class: the lease is read from the local ref (Repo.SHA / MilestoneSHA / EpicSHA) before Box.Queue takes the outbox lock, so a flush in another session that pushes and settles the same key in between leaves the new entry leased on a SHA the remote no longer holds; also refsync.go:222-226 (RecordMilestone) and epic.go:67-71 (RecordEpic); example: session B reads lease L0, session A's flush pushes the pending write (ref -> S1) and Settle drops it, B queues its newer write on L0 -> the next flush loses the race to our own write, and a ticket write is reported Rejected and dropped from the outbox (milestone/epic go through merge instead); blame: d66fcad478, 7556b46a3c, 9cb410a730 (all before the ticket); ran: scratch stress test via go test -overlay (one flusher + a queuer reading the lease outside the lock): 2-8 own-write rejections per run, last write neither on the remote nor queued in 7/10 runs; with lease read and Queue made atomic: 0 rejections, 10/10 pass
- **2026-11-24 07:27 · Sascha** — hunt report: 1 finding, 0 change, 1 pre-existing, as a note on WV2J2W | read whole: core/outbox (outbox.go, race_test.go and package tests), core/refsync/epic.go, epic_race_test.go, the outbox callers in refsync.go (Record*, Flush, mergeRefused, tidyMilestoneEntries), plus ticket.LockFile and gitref Write*/WriteMigration return values | not settled: two flushes in one worktree at once still report our own newer ticket write as Rejected and drop it when it is queued after the first push lands and sent by the second flush before the first settles (seen in a stress run, the last write still landed); what the board does with that local file afterwards is outside these modules, so I did not file it as a finding
- **2026-11-24 07:28 · Sascha** — bughunt: done at cb71b1c7a - 1 finding (0 change, 1 pre-existing) in 3 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/outbox/outbox.go core/outbox/race_test.go core/refsync/epic.go core/refsync/epic_race_test.go core/refsync/refsync.go core/release/notes/WV2J2W.md
- **2026-11-24 07:28 · Sascha** — waiting: review (frischer Agent via run-lane.sh) since 08:28
- **2026-11-24 07:30 · Sascha** — review: done at cb71b1c7a
- **2026-11-24 07:32 · Sascha** — report: Outbox-Lock + Settle/Replace: ein Eintrag wird nach dem Push nur geloescht, wenn er noch der gesendete ist; neuer Eintrag wird umgeleast (Ticket/Milestone/Epic, rebaseEpic, tidy) | signoff, context unknown | 3 Review-Notes: Delete-Lease-Edge, Lease vor Box-Lock gelesen (vorbestehend, Folgeticket?), identity-Test liest echten people-Ref; Branch gepusht, upstream/master eingemergt, PR ist Alex' Sache
- **2026-11-24 10:20 · Sascha** — Given back from signoff: beheben:
- note: edge: outbox.go Settle - nach einem gesendeten Delete (sha leer) behaelt eine waehrend des Pushs eingereihte neue Schreibung ihre alte Lease auf einen Ref, den es nicht mehr gibt, und wird beim naechsten Flush als Rejected verworfen; selten (Loeschen und Neuschreiben desselben Tickets in zwei Sessions gleichzeitig), braucht eine Regel fuer die Lease eines geloeschten Refs
Folge-Ticket:
- AEZ3HD: note: elsewhere: refsync.go:157-161, :222-226 und epic.go:67-71 - die Lease wird vor dem Outbox-Lock gelesen; ein Flush einer anderen Session dazwischen laesst den neuen Eintrag auf einer alten SHA stehen (Bughunt-Notiz, vor diesem Ticket; new ticket)
- X2M486: note: elsewhere: core/identity TestCurrentHonoursJairaUserAboveEverything liest den echten people-Ref und ist rot, unabhaengig von dieser Aenderung
- **2026-11-24 10:20 · Sascha** — cockpit: given back to a new worker - the old worker has no Herdr tab open any more; branch feat/wv2j2w is the ticket's worktree's, worktree /home/alex/projects/.worktrees/jaira-wv2j2w; reports to uds:/run/user/1000/cc-socks/3417.sock, the ticket's parent session
- **2026-11-24 10:21 · Sascha** — parent session: uds:/run/user/1000/cc-socks/3417.sock
- **2026-11-24 10:22 · Sascha** — fix (Rueckgabe signoff): Box.settle(seen, lease, release) - Flush uebergibt release=true, auch bei leerer SHA nach Delete; Settle(seen, lease) bleibt wie bisher (leer = Merge, alte Lease behalten). Leere Lease heisst bei gitref.Write 'kein Ref', also genau der Zustand nach dem Delete. Abgelehntes Delete unveraendert (Folgeticket AEZ3HD-Klasse).
- **2026-11-24 10:23 · Sascha** — stop-after-fix=true gesetzt: Fix 11 Zeilen + Test, kein neuer Lock, kein Parser
- **2026-11-24 10:24 · Sascha** — critique: Settle(seen, "") hat nur Aufrufer mit leerer Lease ausser Flush - Signatur auf Settle(seen) verkleinern, settle(seen, lease) fuer Flush; bool faellt weg
- **2026-11-24 10:25 · Sascha** — critique+optimize (lean, stop-after-fix): nichts mehr; Zeiger statt bool, weil nil = nicht umleasen und "" = kein Ref zwei echte Zustaende sind. Cockpit-Dateien im Diff kommen aus dem upstream-Merge (W80HNH), nicht von hier.
- **2026-11-24 10:26 · Sascha** — testing: done at 6ee1875e9 (lean: nur outbox/refsync/core neu gelaufen)
- **2026-11-24 10:27 · Sascha** — testing: done at 28a06ab4d
- **2026-11-24 10:27 · Sascha** — waiting: review (frischer Agent via run-lane.sh) since 11:27
- **2026-11-24 10:29 · Sascha** — review: done at 28a06ab4d
- **2026-11-24 10:29 · Sascha** — report: Rueckgabe behoben - nach gesendetem Delete wird eine waehrend des Pushs eingereihte Schreibung auf 'kein Ref' umgeleast; Settle(seen) leased nie um, nur Flush via settle(seen, &sha) | signoff, context unknown | 3 elsewhere-Notes (AEZ3HD, X2M486, neu: gitref publish update-ref-Fehler nach Push -> Folgeticket?); Ticketdatei uncommittet bis Logbook; Branch gepusht
- **2026-11-24 10:32 · Sascha** — hunt report: 0 findings, nothing noted | read as a diff since cb71b1c7a: core/outbox/outbox.go (Settle/settle split, Flush lease after a sent Delete), core/outbox/race_test.go, core/refsync/epic.go, core/refsync/refsync.go; epic_race_test.go and release note unchanged; cockpit/W80HNH hunks came by merge and are not this ticket; checked gitref.Repo Write/WriteMilestone/WriteEpic/WriteMigration return a non-empty SHA whenever err is nil, so an empty lease reaches a newer entry only after a Delete; a second Delete queued during the first has the same lease and is dropped by sameWrite; ran: go test -race -count=3 ./core/outbox/ ./core/refsync/ (ok) | not checked: a live push against a real remote of a write leased on no ref after a delete - covered only by the fake sender
- **2026-11-24 10:32 · Sascha** — bughunt: done at 28a06ab4d - 0 findings in 4 modules, follow-up hunt since cb71b1c7a; concurrency, error paths, input, permissions checked; files: core/outbox/outbox.go core/outbox/race_test.go core/refsync/epic.go core/refsync/epic_race_test.go core/refsync/refsync.go core/release/notes/WV2J2W.md
- **2026-11-24 11:22 · Sascha** — Accepted from signoff: Folge-Ticket:
- EPZ78W: note: elsewhere: refsync.go:157-161, :222-226 und epic.go:67-71 - Lease wird vor dem Outbox-Lock gelesen (vorbestehend; Folgeticket AEZ3HD)
- DZ30ZC: note: elsewhere: core/identity TestCurrentHonoursJairaUserAboveEverything liest den echten people-Ref und ist rot, unabhaengig von dieser Aenderung (Folgeticket X2M486)
- 98KFJP: note: elsewhere: gitref.go:548 publish - schlaegt update-ref nach erfolgreichem Push fehl, meldet Flush Failed und der Eintrag bleibt, obwohl er auf dem Remote ist; der naechste Flush verliert gegen die eigene Schreibung (vorbestehend, new ticket)
- **2026-11-24 11:23 · Sascha** — cockpit: accepted by the person, signoff -> done
- **2026-11-24 11:28 · Sascha** — pr: push holds at step 3 - go test ./... -race -count=1 exit 1, only core/identity TestCurrentHonoursJairaUserAboveEverything (Current()="Berk", want "berk"; reads this machine's people ref, X2M486 in backlog, branch does not touch core/identity). 36 packages ok, vet/build/check-ticket-commits/shipped-history ok. Question: 1) push and open the PR anyway, the failure is X2M486 (recommended) 2) stop until X2M486 is fixed
- **2026-11-24 11:28 · Sascha** — pr: parent answered 1 - push and open; the identity test is fixed by AA2C3G
