---
id: 01M865XBRDFR99ENKNMKWV2J2W
title: "outbox: eine Schreibung, die waehrend eines Pushs fuer denselben Schluessel eingereiht wird, wird geloescht"
status: critique
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
commits: []
created-at: 2026-11-24T07:04:21Z
updated-at: 2026-11-24T07:20:30Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1216759
claimed-at: 2026-11-24T07:14:03Z
outcome-what: "outbox.Box nimmt ein Lock-File um jedes Lesen-dann-Schreiben (queueKind, RekeyMilestone, Drop*); neu Settle(seen, lease) und Replace(seen, ...) wirken nur auf den gelesenen Eintrag. Box.Flush (Sent/Rejected), rebaseEpic (Erfolg + Requeue), tidyMilestoneEntries und mergeRefused nutzen sie. Nach Sent wird ein zwischendurch eingereihter Eintrag auf die gepushte SHA umgeleast."
outcome-why: "Ein Eintrag wurde nach Push oder git-Lesen per Schluessel geloescht; eine in diesem Fenster von einer zweiten Session eingereihte Schreibung war weg (Epic: einzige Kopie)."
outcome-resolves: "DoD 1: race_test.go (alle drei Arten, gesendet + abgelehnt) und epic_race_test.go (pre-receive-Hook, direkt + Merge), rot ohne Fix. DoD 2: Release-Note. Geprueft: Nebenlaeufigkeit (Lock nie ueber Push gehalten), Fehlerpfade (Replace atomar, alter Eintrag bleibt bei Fehler), Altpfade (Settle entfernt legacy/old-milestone wie dropKind), Rechte unveraendert."
---

# outbox: eine Schreibung, die waehrend eines Pushs fuer denselben Schluessel eingereiht wird, wird geloescht

## Definition of Done

- [x] Ein Test reiht waehrend des Pushs eine zweite Schreibung fuer denselben Schluessel ein; sie ist danach noch in der Outbox oder auf dem Remote - fuer Ticket, Milestone und Epic
  proof: core/outbox/race_test.go TestAWriteQueuedDuringThePushOfItsKeySurvivesTheFlush (ticket/milestone/epic), TestALostRaceDropsOnlyTheRefusedWrite; core/refsync/epic_race_test.go TestAnEpicEditQueuedDuringThePushReachesTheRemote (pre-receive-Hook, direkt + rebaseEpic); beide rot mit sameWrite=true
- [x] Release-Note in core/release/notes/WV2J2W.md
  proof: core/release/notes/WV2J2W.md

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
