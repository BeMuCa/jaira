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
updated-at: 2026-11-22T17:11:14Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2175608
claimed-at: 2026-11-22T17:03:48Z
outcome-what: "Accept, Give back und eine gewaehlte Entscheidung im Cockpit ziehen das Ticket vorher mit jaira pull, wenn es nur auf dem Ref liegt oder die Datei hinter dem Ref ist; ohne eigenen Worktree wird die dafuer in den Haupt-Checkout gezogene Kopie danach mit git clean entfernt."
outcome-why: "Die Entscheidung lief in einem Verzeichnis, dessen Platte das Ticket nicht oder veraltet trug, und jaira lehnte jede Schreibung mit 'on its ref and not on your disk' bzw. stale_file ab."
outcome-resolves: "Signoff-Entscheidungen im Cockpit scheitern nicht mehr an einem Ticket, das nur im Worktree/auf dem Ref liegt."
---

# Eine Signoff-Entscheidung im Cockpit scheitert nicht, wenn das Ticket nur im Worktree liegt

## Definition of Done

- [x] Accept und Give back auf einem Ticket, das nur im Worktree/auf dem Ref liegt, gelingen (Test)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts:3862 describe 'a ticket not on the disk the press writes from' (accept/give back x own worktree/no worktree x ref-only)
- [x] Ohne Worktree holt das Cockpit das Ticket selbst und entscheidet, statt abzulehnen
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx ticketHere(); test 'accept, no worktree, file ref-only: pulled first, then decided, the pulled copy taken away after'
- [x] Klasse statt Symptom: jede Schreibaktion einer Entscheidung (decide-Note, Accept inkl. Dry-run, Give back) laeuft erst, wenn die Platte den Ref-Stand traegt - on-ref-only und behind-ref werden vorher mit jaira pull geholt (Tabellentest ueber beide Faelle und den Normalfall)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts table cases × state ref-only|behind|current; register.tsx decide()/accept()/giveBack() all via ticketHere
- [x] Eine nur fuer die Entscheidung in den Haupt-Checkout gezogene Kopie ist danach wieder weg, auch wenn die Entscheidung scheitert; eine Kopie im eigenen Worktree bleibt
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts:3908 'a pulled copy is taken away from the main checkout when the press is refused too'; cleans:false for own worktree rows
- [x] Ein abgelehnter pull (taken/offline) stoppt die Entscheidung mit seiner Meldung, schreibt nichts und nimmt nie --steal
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts:3921 'a refused pull stops the press with its reason, writes nothing and never steals'

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
