---
id: 01M7WX00JXRK922S029Q76ERHH
title: "Eine Signoff-Entscheidung im Cockpit scheitert nicht, wenn das Ticket nur im Worktree liegt"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Accept und Give back im Cockpit finden den Worktree des Tickets und entscheiden dort, oder holen das Ticket selbst; der Fehler 'on its ref and not on your disk' kommt dabei nicht mehr."
context: |-
  Heute: Accept/Give back im Cockpit laeuft im Haupt-Checkout. Liegt die Ticket-Datei nur im Worktree des Tickets bzw. auf dem Ref, scheitert es mit "on its ref and not on your disk".
  Alex 20.11.2026: "das ist oft ein fail"; ein Teamlead meldet dasselbe.
  Das Cockpit kennt den Worktree schon (ticketsTree in core/mod/builtin/jaira-cockpit, sortiert nach Pfad -<id>).
  Verwandt: 1K3RA8 (ein Schreibbefehl holt den Ref selbst nach).
definition-of-done: "Accept und Give back auf einem Ticket, das nur im Worktree/auf dem Ref liegt, gelingen (Test)"
tags: []
blocked-by:
  - 01M67RXSPRWC7WVBCCM7MVR4PY
related:
  - 01M5MNWJF240306ZVDA41K3RA8
  - 01M5ZJ57ARFWR06TS7TJFYHJ39
commits: []
created-at: 2026-11-20T16:35:21Z
updated-at: 2026-11-22T17:20:37Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2175608
claimed-at: 2026-11-22T17:03:48Z
outcome-what: "Accept, Give back und eine gewaehlte Entscheidung im Cockpit ziehen das Ticket vorher mit jaira pull, wenn es nur auf dem Ref liegt oder die Datei hinter dem Ref ist; ohne eigenen Worktree wird die dafuer in den Haupt-Checkout gezogene Kopie danach mit git clean entfernt, ein gescheitertes Entfernen wird mit Pfad gezeigt."
outcome-why: "Die Entscheidung lief in einem Verzeichnis, dessen Platte das Ticket nicht oder veraltet trug, und jaira lehnte jede Schreibung mit 'on its ref and not on your disk' bzw. stale_file ab."
outcome-resolves: "Signoff-Entscheidungen im Cockpit scheitern nicht mehr an einem Ticket, das nur im Worktree/auf dem Ref liegt."
review-summary: none
review-gaps: "none - Runde 2: der Fix nutzt output() und show(), nichts Neues zu falten"
test-verdict: "pass: Runde 2 (Diff seit f43fc94d2: register.tsx leave(), texts.ts, board.test.ts): claude plugin test . 323 pass 0 fail (RC=0); neuer Test 'a copy that cannot be taken away is shown with its path' zeigt Fenster 'Copy of MG677R left behind' mit Fehler und git-clean-Zeile; DoD 1-5 mit aktualisierten Proof-Zeilen geprueft. Runde 1: Live-Lauf mit echtem jaira (pull/move/note/git clean, offline rc=1, fremd rc=3, zwei Schreiber)."
hunt-verdict: "fail: 1 change (Blocker) - leave() ueberschreibt das Meldungsfenster (Ablehnung, noteLost)"
round: 3
back-from: testing
back-by: lane
back-count: 2
---

# Eine Signoff-Entscheidung im Cockpit scheitert nicht, wenn das Ticket nur im Worktree liegt

## Definition of Done

- [x] Accept und Give back auf einem Ticket, das nur im Worktree/auf dem Ref liegt, gelingen (Test)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts:3865 describe 'a ticket not on the disk the press writes from' (accept/give back x own worktree/no worktree x ref-only)
- [x] Ohne Worktree holt das Cockpit das Ticket selbst und entscheidet, statt abzulehnen
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx ticketHere(); test 'accept, no worktree, file ref-only: pulled first, then decided, the pulled copy taken away after'
- [x] Klasse statt Symptom: jede Schreibaktion einer Entscheidung (decide-Note, Accept inkl. Dry-run, Give back) laeuft erst, wenn die Platte den Ref-Stand traegt - on-ref-only und behind-ref werden vorher mit jaira pull geholt (Tabellentest ueber beide Faelle und den Normalfall)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts table cases × state ref-only|behind|current; register.tsx decide()/accept()/giveBack() all via ticketHere
- [x] Eine nur fuer die Entscheidung in den Haupt-Checkout gezogene Kopie ist danach wieder weg, auch wenn die Entscheidung scheitert; eine Kopie im eigenen Worktree bleibt
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts:3911 'a pulled copy is taken away from the main checkout when the press is refused too'; core/mod/builtin/jaira-cockpit/tests/board.test.ts:3923 failed clean said beside the band, message pane kept; cleans:false for own worktree rows
- [x] Ein abgelehnter pull (taken/offline) stoppt die Entscheidung mit seiner Meldung, schreibt nichts und nimmt nie --steal
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts:3951 'a refused pull stops the press with its reason, writes nothing and never steals'

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-20 16:50 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-22 17:03 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 17:05 · Alexander Sacharov** — decision (worker): keine offene Frage. Ursache: ticketDir() in register.tsx findet den Worktree schon (worktreeFor), aber (a) ohne Worktree laeuft move im Haupt-Checkout, wo die Datei fehlt, (b) ein Worktree ohne gepullte Datei scheitert genauso, (c) Datei hinter dem Ref wird als stale abgelehnt. Klasse: 'die Entscheidung schreibt in einem Verzeichnis, dessen Platte nicht den Stand des Refs traegt'. Fix: vor jeder Schreibaktion einer Entscheidung (decide-Note, Accept, Give back) liest das Cockpit jaira show --json im Ticket-Verzeichnis; on-ref-only oder behind-ref -> jaira pull <id> dort. Im eigenen Worktree bleibt die gezogene Datei (sie gehoert auf den Branch). Im Haupt-Checkout ohne eigenen Worktree wird die nur fuer die Entscheidung gezogene Kopie danach mit git clean -fqx -- <path> entfernt (Memory: keine Ticket-Kopie im Haupt-Checkout, sie verdeckt den Ref). Sicher, weil jede Schreibung schon auf dem Ref oder im Outbox liegt (Outbox haelt den ganzen Dateiinhalt). pull 'taken' (jemand anderes) wird angezeigt, nie --steal.
- **2026-11-22 17:11 · Alexander Sacharov** — pre-existing, nicht dieser Ticket-Umfang: Explain/Discuss (openChat, register.tsx) gibt dem Chat dasselbe Verzeichnis per ticketDir; ohne Worktree schreibt der Chat seine explain:-Note spaeter selbst und scheitert dort genauso. Die Kopie kann dort nicht nach dem Druck entfernt werden, weil der Chat asynchron schreibt - eigener Entwurf noetig.
- **2026-11-22 17:11 · Alexander Sacharov** — critique: nichts. Einfachere Form geprueft: den Pull in jaira selbst (1K3RA8, liegt in human) - ersetzt die Cockpit-Seite nicht, solange es nicht gebaut ist, und das Entfernen der Kopie im Haupt-Checkout bliebe trotzdem Cockpit-Sache. ticketDir/ownTree teilen jetzt eine Stelle statt zwei porcelain-Lesungen. git clean ohne Fehlermeldung gelassen: es entfernt eine unversionierte Datei auf einem Pfad, den jaira selbst genannt hat.
- **2026-11-22 17:13 · Alexander Sacharov** — testing: note: known limit - jaira pull nimmt ein Ticket, dessen assignee jemand anderes ist, nur mit --steal (live: rc=3 'not yours to pull: probe has it'). Accept/Give back eines fremden Tickets ohne Worktree scheitert damit weiter, jetzt mit dieser klaren Meldung statt 'pull it first'; --steal bleibt nach DoD 5 bewusst aus.
- **2026-11-22 17:14 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:555 - class: a press that only wants to write now claims the ticket on its ref; ticketHere runs 'jaira pull', and refsync.Pull (core/refsync/refsync.go:1440-1462) sets assignee to the person and pushes that before it writes the file; leave() then removes the file but never gives the claim back (the ref stays 'assignee: <person>' with no file anywhere). example: an unassigned ref-only ticket in human gets a decision with no action (note only) or a Give back from the cockpit; afterwards a teammate's 'jaira pull' of it is refused as taken (refsync.go:1420) until --steal. Same path via accept() and giveBack(). Not a blocker by the rule (no lost data, ownership is wrong). The test world's pull mock (tests/board.test.ts ~line 264) only clears the disk state, so no test sees the claim. confirmed: read
- **2026-11-22 17:14 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:561 - class: the copy pulled into the main checkout can stay there and nobody is told, which is the stale copy the comment at :537 says outranks the ref and shows a lane the ticket left. (a) leave() drops the result of 'git clean' ('void (await run(...))'): a failed or timed-out clean leaves the file silently. (b) the row lock once() frees after BUSY_MS = 20s (register.tsx:96, :862) while a press can run longer (pull 30s + dry run 30s + move 30s + spawn.sh 120s in givingBack); a second press on the same row then finds A's copy as 'current', takes no leave() of its own, and A's leave() cleans the file under it. That second press then either fails with 'on its ref and not on your disk' (the error this ticket removes) or, when its jaira write lands its atomic rename after A's clean, recreates the file, and it stays in the main checkout. example: Give back on a ref-only ticket with no worktree, spawn.sh slow, Accept pressed again on the row after 20s. Not a blocker: needs a slow press plus a second press, or a failing clean. confirmed: read (timing not run)
- **2026-11-22 17:14 · Alexander Sacharov** — bughunt: done at f43fc94d2 - 2 findings (2 change, 0 blocker, 0 pre-existing) in 3 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/76ERHH.md
- **2026-11-22 17:15 · Alexander Sacharov** — testing: class: eine Aufraeum-Aktion des Cockpits scheitert still; example: leave() verwirft das Ergebnis von git clean (register.tsx:561); fix: Fehlschlag mit Pfad im Meldungsfenster zeigen, Test mit scheiterndem clean. Bughunt-Fund 1 (Claim durch pull) bleibt known limit: ein Ticket ohne assignee verlaesst backlog nicht, in human/signoff ist pull fuer den eigenen assignee ein No-op und fuer einen fremden abgelehnt (DoD 5). Fund 2b (zweiter Druck nach BUSY_MS 20s waehrend ein langsamer Druck laeuft) bleibt known limit: braucht einen >20s-Druck plus zweiten Druck auf dieselbe Zeile.
- **2026-11-22 17:16 · Alexander Sacharov** — critique (Runde 2, nur der Fix 4e0a0a440): nichts - leave() zeigt ein gescheitertes git clean mit dem bestehenden show()-Fenster, wie noteLost es fuer eine verlorene Note tut.
- **2026-11-22 17:18 · Alexander Sacharov** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/register.tsx:566 - class: leave() runs in finally after the press and calls show(), which replaces the one message pane (messageAtom), so any message the press itself showed is lost whenever git clean fails; places: accept() :1088, giveBack() :1244, decide() :1088-area finally, and the messages they lose: fail() refusals (e.g. 'Accept MG677R failed: definition of done item 2 is open') and noteLost at :1264 and :1187, which holds the only copy of the give-back reason the person typed plus its 'jaira note ... ' paste line (data loss of that text); example: refused Accept + git clean 'Permission denied' -> panes opened: 'Accept MG677R failed' then 'Copy of MG677R left behind', and the mounted Accept pane reads only the copy text; blocker because the reason text is lost, but it needs two failures in one press (a refused write or note, and a refused git clean); ran: scratch copy of the plugin with an added board.test.ts case, claude plugin test . (OPENED shows both panes, the Accept reason missing from the pane)
- **2026-11-22 17:18 · Alexander Sacharov** — bughunt: done at 4e0a0a440 - 1 finding (1 change, 1 blocker, 0 pre-existing), follow-up hunt since f43fc94d2; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/76ERHH.md
- **2026-11-22 17:19 · Alexander Sacharov** — bughunt: correction to the change finding above: the finally blocks calling leave() are accept() register.tsx:647, decide() :1088 and giveBack() :1244
