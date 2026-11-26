---
id: 01M86MPK2ZKG4T824GG598KFJP
title: "gitref.go:548 publish - schlaegt update-ref nach erfolgreichem Push feh…"
status: testing
ready: true
creator: Sascha
context: |-
  Was falsch ist: core/gitref/gitref.go publish (Zeile 542-553) pusht erst, dann setzt es den lokalen Ref mit update-ref. Schlaegt update-ref fehl (z.B. ein fremdes .lock auf dem Ref), gibt publish (commit, err) zurueck.
  Folge: core/outbox/outbox.go Box.Flush sieht den Fehler, meldet Failed und laesst den Eintrag mit der alten Lease stehen, obwohl die Schreibung auf dem Remote ist.
  Der naechste Flush pusht mit der alten Lease gegen die eigene Schreibung, verliert (ErrRaceLost) und meldet Rejected: Ticket-Eintraege werden verworfen, Milestone/Epic in den 'Gewinner' gemergt - beides Laerm fuer eine erfolgreiche Schreibung.
  Vorbild: Delete (gitref.go:558) ignoriert ein fehlgeschlagenes lokales update-ref -d bereits; der Remote ist die Wahrheit, der naechste Fetch zieht den lokalen Ref nach.
  Gefunden: Review zu WV2J2W am 24.11.2026 (note: elsewhere). Vorbestehend.
tags: []
blocked-by: []
related: []
follows: 01M865XBRDFR99ENKNMKWV2J2W
commits: []
created-at: 2026-11-24T11:22:48Z
updated-at: 2026-11-26T08:08:05Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-37950
claimed-at: 2026-11-26T07:54:51Z
goal: "Haelt ein anderer git-Prozess das lokale Ref-Lock kurz, waehrend eine Schreibung pusht, gilt die angenommene Schreibung als gesendet (Flush meldet Sent, der Eintrag verschwindet) und der lokale Ref steht danach auf dem Remote-Stand, sodass die naechste Schreibung nicht gegen die eigene verliert."
outcome-what: "gitref publish laesst update-ref nach angenommenem Push bis zu 1 s auf ein fremdes Ref-Lock warten (core.filesRefLockTimeout=1000 statt gits 100 ms); ein dauerhaft unbeschreibbarer Ref bleibt ein Fehler."
outcome-why: "Ein kurz gesperrter lokaler Ref liess Flush Failed melden und den Eintrag behalten; den Fehler zu schlucken (Runde 1) liess den lokalen Ref hinter dem Remote, und die naechste Schreibung verlor gegen die eigene. Warten haelt den lokalen Ref aktuell."
outcome-resolves: "DoD 1-3: TestAWriteWaitsOutABrieflyLockedLocalRef, TestALocalRefThatStaysLockedIsStillAnError, TestABrieflyLockedLocalRefNeitherFailsTheFlushNorTheNextWrite, Release-Note."
review-summary: none
review-gaps: "lockRef in gitref_test herausgezogen (Lock-Datei stand in holdRefLock und im Dauer-Lock-Test doppelt); gelassen: dieselben 6 Zeilen in outbox_test - anderes Paket, ein Testhilfspaket dafuer waere mehr als es spart"
test-verdict: "pass: r2 go test -race gitref/outbox/refsync gruen (RC=0), vet/build ok; DoD 1-3 im Baum, Tests rot gegen master-gitref.go; Scratch-Board: Lock 400 ms -> 'jaira set' exit 0, lokal = Remote, naechstes set geht durch; Dauer-Lock -> 'could not be sent' mit git-Meldung (laut), danach verliert die naechste Schreibung gemeldet (vorbestehend, Repo kaputt); zwei parallele Schreiber: einer landet, der andere 'was not sent' klar abgewiesen"
hunt-verdict: "fail: 1 change blocker (gitref.go:554 lokaler Ref bleibt hinter dem Remote), 1 pre-existing (gitref.go:998 jede Server-Ablehnung gilt als ErrRaceLost) - Note fuer die Person"
round: 2
back-from: testing
back-by: lane
back-count: 1
---

# gitref.go:548 publish - schlaegt update-ref nach erfolgreichem Push feh…

## Definition of Done

- [x] Ein Test haelt den lokalen Ref kurz (300 ms, laenger als gits 100 ms) gesperrt, waehrend Write pusht: Write liefert die SHA ohne Fehler, lokaler Ref und Remote tragen sie - fuer Ticket, Milestone, Migration, Epic; ein dauerhaft gesperrter Ref bleibt ein Fehler
  proof: core/gitref/gitref_test.go TestAWriteWaitsOutABrieflyLockedLocalRef, TestALocalRefThatStaysLockedIsStillAnError (rot gegen master-gitref.go)
- [x] Ein Outbox-Test: Flush meldet in diesem Fall Sent, der Eintrag ist weg, und die naechste Schreibung mit der Lease vom lokalen Ref geht durch statt gegen die erste zu verlieren
  proof: core/outbox/outbox_test.go TestABrieflyLockedLocalRefNeitherFailsTheFlushNorTheNextWrite (rot gegen master-gitref.go)
- [x] Release-Note in core/release/notes/98KFJP.md
  proof: core/release/notes/98KFJP.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-26 07:55 · Sascha** — parent session: uds:/run/user/1000/cc-socks/62036.sock
- **2026-11-26 07:56 · Sascha** — decision: publish gibt nach angenommenem Push nil zurueck statt den update-ref-Fehler - wie Delete es mit update-ref -d schon tut. Kein Retry: der lokale Ref ist nur ein Cache des Remote, der naechste Fetch zieht ihn nach. Offene Entscheidungen fuer die Person: keine.
- **2026-11-26 07:58 · Sascha** — Geschwister geprueft: Delete (update-ref -d) und PushBranch ignorieren das lokale update-ref schon; publish war der einzige Schreibpfad, der es als Fehler meldete. Bekannte Grenze: bleibt der lokale Ref wegen des Locks stehen, liest eine Schreibung derselben Session vor dem naechsten Fetch die alte SHA als Lease und verliert - das war vorher genauso, nur kommt jetzt die erste Schreibung nicht mehr als Failed zurueck. Testweg: ein .lock-File neben dem Ref unter .git/ laesst update-ref scheitern, ohne den Push zu stoeren.
- **2026-11-26 07:58 · Sascha** — critique: nichts. Fix sitzt im Callee (publish), alle vier Schreibpfade laufen hindurch; Muster wie Delete/PushBranch in gitref.go. Kein einfacherer Schnitt.
- **2026-11-26 08:00 · Sascha** — testing: note: edge: lokaler Ref bleibt nach Lock stehen; Beispiel: Scratch-Board, .lock auf den Ticket-Ref, 'jaira set' -> exit 0, Remote neu, lokal alt; Lock weg, naechstes 'jaira set' -> 'was not sent: ada wrote it first' (gemeldet, Datei unveraendert, kein Datenverlust). Vorbestehend (vorher landete die erste Schreibung als Failed in der Outbox und verlor genauso). Fix waere: Lease nach einem Rejected gegen den Remote-Tip pruefen, ob der eigene Inhalt schon dort liegt - neue Regel, nicht dieses Ticket.
- **2026-11-26 08:00 · Sascha** — waiting: bughunt (Subagent) und go test -race ./... since 09:00
- **2026-11-26 08:01 · Sascha** — bughunt: change: blocker: core/gitref/gitref.go:554 - class: after a push the remote accepted, a failed local update-ref is swallowed, so the local ref stays behind the remote while every lease reader treats it as 'what the remote last confirmed' (refsync.Record core/refsync/refsync.go:157, RecordDelete :176, RecordMilestone :222, core/refsync/epic.go:67, CheckFresh/isStale); the next write leased from it loses against this clone's own push, outbox.Flush (core/outbox/outbox.go:573-580) settles a ticket entry away as Rejected, and the ref keeps the older state. Same class, second example: a permanent update-ref failure (no write permission on .git/refs, full disk) is now silent on every write, where it used to surface. example: write v1 with the ref lock held -> Sent; lock released; next Record leases '' (local ref still empty); Flush fetches, pushes v2 -> 'the ticket moved since it was read', queue empty, ref carries v1. ran: go test -overlay (scratch test in core/outbox, Queue+Flush with a .lock file, then SHA/Queue/Fetch/Flush) -> second flush rejected, 0 queued, ref = v1. Marked change because line 554 is on the path and core/release/notes/98KFJP.md promises 'the next sync no longer reports it as a lost race against your own write'; the same test against the pre-change gitref.go also rejects v2, so this is the fix stopping one write short, not a regression
- **2026-11-26 08:02 · Sascha** — bughunt: pre-existing: core/gitref/gitref.go:998 - class: any server-side refusal whose stderr contains 'rejected' (pre-receive hook declined, protected ref, permission or quota refusal: '[remote rejected] ...') is classified as ErrRaceLost, and outbox.Flush (core/outbox/outbox.go:573-580) then settles a ticket entry away, so the write is dropped and reported as a lost race instead of kept as Failed; example: bare remote with a pre-receive hook that exits 1 -> outcome rejected, 'the ticket moved since it was read: remote: refs/jaira is protected', 0 entries left; blame: 03961b0ca4 / d3bc99ef64 (2026-09-10, before the ticket); ran: go test -overlay (scratch test in core/outbox with a refusing pre-receive hook)
- **2026-11-26 08:02 · Sascha** — hunt report: 2 findings (1 change blocker, 1 pre-existing), both as notes on 98KFJP | read whole: core/gitref (gitref.go and its tests), core/outbox/outbox.go Flush/send/settle and its new test, core/release/notes/98KFJP.md; callers in core/refsync read along the lease and write paths | not checked: a real concurrent fetch holding the ref lock (simulated with a .lock file); whether milestone migration (HasMilestoneRef/ReadMigration on a stale local ID ref) misbehaves before the next fetch — not reproduced, so not reported
- **2026-11-26 08:02 · Sascha** — bughunt: done at 9870084be - 2 findings (1 change, 1 pre-existing) in 3 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/gitref/gitref.go core/gitref/gitref_test.go core/outbox/outbox_test.go core/release/notes/98KFJP.md
- **2026-11-26 08:03 · Sascha** — decision: kleinerer Entwurf statt Fehler schlucken. publish laesst update-ref mit -c core.filesRefLockTimeout=1000 auf ein fremdes Lock warten (git-Default 100 ms): der realistische Fall - ein paralleler fetch haelt das Ref-Lock kurz - endet mit aktuellem lokalem Ref, also ohne Selbst-Rennen bei der naechsten Schreibung. Ein dauerhaft unbeschreibbarer lokaler Ref (Rechte, volle Platte) bleibt ein Fehler (Failed, laut), wie vor dem Ticket - ihn zu schlucken liess jede folgende Schreibung ohne sichtbare Ursache verlieren (Bughunt-Befund). Damit gilt das Ziel fuer das voruebergehende Lock; DoD 1-2 entsprechend umformuliert.
- **2026-11-26 08:05 · Sascha** — round 2: Bughunt-Blocker (lokaler Ref hinter Remote -> naechste Schreibung verliert gegen sich selbst) behoben, indem update-ref auf das Lock wartet statt den Fehler zu schlucken; Test-Edge von Runde 1 (veralteter lokaler Ref) damit fuer das kurze Lock erledigt. Tests halten das Lock 300 ms nach dem angenommenen Push - mit gits 100 ms rot (8/8), mit 1000 ms gruen. Pre-existing gitref.go:998 (jede Server-Ablehnung = ErrRaceLost) bleibt Note fuer die Person: anderes Thema als das Ziel.
- **2026-11-26 08:06 · Sascha** — critique r2: nichts. gits eigene Lock-Wartezeit (core.filesRefLockTimeout) statt einer Retry-Schleife in Go ist der kleinste Schnitt; Fehlerpfad wie vor dem Ticket.
- **2026-11-26 08:07 · Sascha** — waiting: bughunt r2 (Subagent) since 09:07
- **2026-11-26 08:08 · Sascha** — bughunt: change: blocker: core/gitref/gitref.go:555 - class: a push the remote accepted is still returned as an error whenever the local ref lock outlasts the 1 s filesRefLockTimeout (a stale .lock left by a crashed git, a fetch slower than 1 s); every caller reads that error as 'not written': outbox.Flush (core/outbox/outbox.go:589 default case) reports Failed and keeps the entry with its old lease, and its next round is refused as ErrRaceLost and, for a ticket, settled away as Rejected; refsync claim (core/refsync/refsync.go:1659) fails without writing the local file. example: lock refs/jaira/... in the clone, Write -> error, remote holds the commit, local ref absent; the next Write leased on repo.SHA() is refused 'the ticket moved since it was read' - against its own first write. ran: scratch test in core/gitref (lockRef, Write, remove lock, SHA, Write) -> second write errors.Is ErrRaceLost=true. Unclear mark: for locks over 1 s this is the same behaviour the code had before 9870084be (update-ref error returned); b70a95f7c narrowed the window but rewrote these lines, so it is marked change - a fix needs the error to carry 'remote has it' (e.g. return the sha with a distinct error the outbox settles on) rather than 'failed'
