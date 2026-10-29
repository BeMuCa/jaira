---
id: 01M64NQ6XEX4WFG855H5ZKPT4H
title: "Web zeigt den Stand vom Ref, Erstelldatum und wie lange ein Ticket unberührt ist"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Das Board im Browser zeigt für jedes Ticket die Lane, die auf seinem Ref steht, und dazu, wann es angelegt wurde und wie lange es schon niemand angefasst hat."
context: |-
  Gefunden am 29.10.2026 auf einem Board mit geteilten Refs:
  - jaira web zeigt ein Ticket als in-progress, auf seinem Ref steht es seit Tagen in signoff.
  - Ursache: web liest die Ticketdatei im Arbeitsverzeichnis (master). Die Datei kam mit dem letzten Code-Commit dorthin; spätere Lanes ohne Code (critique, review, signoff) schreiben nur den Ref. Die Datei bleibt stehen.
  - Betroffen waren dort über 20 Tickets: Datei sagt critique/testing/in-progress, Ref sagt signoff.
  - Gewünscht (Alex): web nimmt den Ref und führt ihn mit der Datei zusammen – der neuere Stand gewinnt (updated-at), wie 'jaira pull' es schon beim Vergleich tut ('the ticket file here is older than its ref').
  - Außerdem fehlen auf der Karte: Erstelldatum (created) und 'unberührt seit' (Zeit seit updated-at, z.B. '3 Tage').
  - Regel für diese Arbeit: kein Name, Pfad und keine Ticket-ID eines anderen Projekts in Code, Tests, Commits, Notes oder Release-Notes. Testdaten sind erfunden.
definition-of-done: "Ein Ticket, dessen Datei älter ist als sein Ref, erscheint im Web in der Lane des Refs (Test mit Datei=in-progress, Ref=signoff)"
tags:
  - web
blocked-by: []
related: []
commits: []
created-at: 2026-10-29T20:30:44Z
updated-at: 2026-10-29T20:52:33Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1236953
claimed-at: 2026-10-29T20:37:14Z
outcome-what: "jaira web zeigt den Ref-Stand bei älterer Datei, Erstelldatum und 'untouched' auf jeder Karte, und 'archive' als Lane für archivierte/abgelegte Tickets."
outcome-why: Lanes ohne Code schreiben nur den Ref; abgelegte Tickets behalten ihre alte Lane und lasen sich wie offene Arbeit.
outcome-resolves: Seite und Ref widersprechen sich nicht mehr; Liegengebliebenes und Abgelegtes ist auf der Karte erkennbar.
review-summary: "jaira web liest jetzt alle Ticket-Refs in einem git-Aufruf (neues refsync.Syncer.Newer über gitref.ReadMany) und ersetzt jede Ticketdatei, deren updated-at älter ist als der Ref, durch den Ref-Stand, bevor der Link-Index gebaut wird — dadurch zeigen Karten, Epic-Zählung und Detailansicht die Lane vom Ref; eine neuere Datei bleibt, weil derselbe isStale-Vergleich wie bei pull/CheckFresh entscheidet. Zusätzlich trägt web.Ticket created_at/updated_at, und TicketCard zeigt 'created <Datum>' und 'untouched 3d/5h/<1h' (im Browser aus updated_at berechnet). dist ist aus dem Quellstand gebaut (Neubau ergibt dieselben Hashes)."
round: 4
review-gaps: "Kein Blocker. (1) Der Test prüft nur lose Karten (/api/unassigned) und die Detailansicht; der Epic-Pfad (Epics(): Zählung und Daten, inkl. LoadAnywhere für abgelegte Mitglieder) ist ungetestet, obwohl das Outcome 'Epic-Zählung' nennt. (2) age()/date() haben keinen Test — es gibt im Repo kein UI-Testframework, das räumt der Implementierer selbst ein. (3) Der ersetzte Ticketstand ist ReadOnly mit leerem Path; im Link-Index gilt er damit als 'nur auf dem Ref' (PlaceRef), obwohl eine Datei existiert — die Seite zeigt Place heute nicht, fällt also nicht auf, ist aber eine stille Bedeutungsverschiebung. (4) Scheitert ReadMany, zeigt die Seite ohne Hinweis wieder den Dateistand. (5) TUI und 'jaira list' lesen weiter nur die Datei — bewusst außerhalb des Tickets (Note vom 29.10.)."
test-verdict: "pass: go build ./... und go test -race -count=1 ./... RC=0; dist == npm ci && npm run build; DoD 1+2 per TestWebShowsTheNewerOfFileAndRef (Gegenprobe ohne Newer rot), DoD 3+4 per Test (created_at/updated_at) und age()/date() unter node geprüft (<1h, 5h, 1d, 3d, leer); jaira web --no-open auf diesem Board liefert Daten auf jeder Karte"
review-verdict: "Der Diff erfüllt alle vier DoD-Punkte; go build/vet und die Tests in internal/cli (TestWeb*), core/refsync, core/web, core/link laufen grün, das ausgelieferte Bundle entspricht dem Quellstand. Keine Defekte gefunden; die Lücken oben sind Testabdeckung und eine Randbedeutung von ReadOnly, nichts, das die Annahme aufhalten muss."
review-round: 2
review-check: "1. cd in das Repo dieses Branches, go run ./cmd/jaira web  2. ein Ticket nehmen, das per 'jaira move' auf dem Ref in signoff liegt, dessen Datei unter .jaira/tickets/ aber noch 'status: in-progress' und ein älteres updated-at hat  3. Seite neu laden: die Karte steht in signoff, nicht in in-progress; Karte anklicken — Detail zeigt ebenfalls signoff  4. in der Datei updated-at auf 2099-01-01T00:00:00Z setzen, neu laden: Karte zeigt jetzt in-progress  5. unter jedem Kartentitel steht eine graue Zeile 'created <Datum> · untouched <Nd|Nh|<1h>'  6. Datei danach mit 'git checkout -- .jaira/tickets/<datei>' zurücksetzen. Automatisch: go test ./internal/cli -run TestWebShowsTheNewerOfFileAndRef"
---

# Web zeigt den Stand vom Ref, Erstelldatum und wie lange ein Ticket unberührt ist

## Definition of Done

- [x] Ein Ticket, dessen Datei älter ist als sein Ref, erscheint im Web in der Lane des Refs (Test mit Datei=in-progress, Ref=signoff)
  proof: internal/cli/web_source.go board(): refs.Newer(all) ersetzt veraltete Dateien vor link.Build; TestWebShowsTheNewerOfFileAndRef (Datei=in-progress/2026-09-10, Ref=signoff → Karte und Detail signoff)
- [x] Ist die Datei neuer als der Ref, gewinnt die Datei
  proof: core/refsync/refsync.go Syncer.Newer nimmt ein Ticket nur auf, wenn isStale(Datei, Ref); TestWebShowsTheNewerOfFileAndRef (Datei updated-at 2099 → in-progress)
- [x] Jede Karte zeigt das Erstelldatum
  proof: core/web/ui/src/entities/ticket/ui/TicketCard.svelte 'created {date(created_at)}'; created_at gesetzt in internal/cli/web_source.go Unassigned() und Epics(); TestWebShowsTheNewerOfFileAndRef prüft created_at
- [x] Jede Karte zeigt, wie lange das Ticket unberührt ist (seit updated-at), in Tagen/Stunden
  proof: core/web/ui/src/shared/lib/time.ts age() (d ab 24h, sonst h, <1h); TicketCard.svelte 'untouched {age(updated_at)}'; TestWebShowsTheNewerOfFileAndRef prüft updated_at
- [x] Ein Ticket, das im Archiv bzw. im Logbook liegt, zeigt auf seiner Karte im Web den Status 'archive' statt der Lane, in der es abgelegt wurde (z.B. nicht 'backlog') – sonst sieht es aus, als müsse man es noch machen
  proof: internal/cli/web_source.go Epics(): ix.Lookup Place archive/logbook → Lane archived; Ticket(): Pfad außerhalb tickets/ → archived; TestWebShowsAnArchivedTicketAsArchive (Gegenprobe mit backlog rot)
- [x] Der Fortschrittsbalken eines Epics im Web (EpicProgress) zeigt neben 'done' auch die Tickets in 'signoff' als eigenen, farblich abgesetzten Teil, mit Zahl im Tooltip/Text (z.B. '4 done · 3 in signoff of 10') – Abnahme dauert, und heute sieht fertig-wartende Arbeit aus wie nicht angefangen
  proof: core/web/ui/src/entities/epic/ui/EpicProgress.svelte zweites Segment .signoff (--waiting) + Text '{done} done · {signoff} in signoff of {total}'; internal/cli/web_source.go Epics() zählt Lane signoff; TestWebCountsSignoffApartFromDone (Gegenprobe rot)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] refsync: exportierte Methode Newer(t) – liefert das Ref-Ticket, wenn isStale(t, ref), sonst nil
- [x] web_source: board() ersetzt veraltete Dateien durch den Ref-Stand, bevor link.Build den Index baut; Ticket(id) ebenso
- [x] web.Ticket: created_at/updated_at; webSource füllt sie für Epic- und lose Karten
- [x] UI: TicketCard zeigt created und untouched; dist neu bauen
- [x] Tests: web_source mit Datei=in-progress/Ref=signoff und umgekehrt; Go-Test für die Felder; Release-Note

## Progress
- **2026-10-29 20:35 · Alexander Sacharov** — decision (Alex, 29.10.): keine Erwähnung anderer Projekte – weder in Code, Tests, Fixtures, Commits, Notes noch Release-Notes. Kontext dafür umgeschrieben.
- **2026-10-29 20:37 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/1226588.sock
- **2026-10-29 20:38 · Alexander Sacharov** — Entscheidungen (selbst getroffen, keine offen für Alex): (1) Regel 'neuer gewinnt' ganz, nicht der Feld-Merge aus refsync.Reconcile – der hält status nie rückwärts und würde DoD 2 (Datei neuer, frühere Lane) brechen. Vergleich über refsync.isStale, dieselbe Definition wie pull/CheckFresh. (2) Nur jaira web; TUI/list lesen weiter die Datei – gleicher Fehler möglich, nicht Teil dieses Tickets. (3) Karte englisch wie der Rest der Seite: 'created 14.10.26' (locale) und 'untouched 3d' / '5h' / '<1h'; Tage ab 24h, sonst Stunden, berechnet im Browser aus updated_at. (4) Auch die Detailansicht nimmt den neueren Stand. (5) Epic-Karten eines abgelegten Tickets: Daten per LoadAnywhere.
- **2026-10-29 20:41 · Alexander Sacharov** — Gegenprobe: mit Newer→nil fällt TestWebShowsTheNewerOfFileAndRef ('card in-progress; want signoff'). Live-Probe: jaira web --no-open auf diesem Board, /api/unassigned liefert created_at/updated_at, Seite lädt das neue Bundle. Kein UI-Testframework im Repo (kein vitest) – age()/date() nur per Augenschein im Build.
- **2026-10-29 20:43 · Alexander Sacharov** — critique (Runde 1): card() in web_source.go hat einen Aufrufer – inlinen. Stehen gelassen: board() dupliziert ~8 Zeilen von boardIndex (links.go:113) statt boardIndex zu ändern, weil boardIndex auch epic ls/show bedient und dessen Verhalten nicht Teil dieses Tickets ist.
- **2026-10-29 20:45 · Alexander Sacharov** — optimize: Newer war pro Ticket ein git-Prozess – jede Seitenanfrage 5x langsamer (0,1 → 0,55 s bei 130 Tickets). Jetzt gebündelt über gitref.ReadMany wie Syncer.Extra; zurück auf 0,08–0,12 s. Verhalten unverändert, Test grün.
- **2026-10-29 20:48 · Alexander Sacharov** — review: grün, keine Defekte. Lücken nur Testabdeckung (Epic-Pfad, age()/date()) und ReadOnly-Ersatz = PlaceRef im Index. Nicht zu diesem Ticket gehörig: .jaira/milestones/ liegt untracked im Worktree und taucht deshalb in worktree_diff auf.
- **2026-10-29 20:49 · Alexander Sacharov** — decision (Alex, 29.10.): Karten von archivierten/abgelegten Tickets (z.B. in einem Epic, 'via follows/parent') zeigten das Lane-Badge 'backlog' plus die Zeile 'in the archive'. Das Badge soll 'archive' heißen, wenn das Ticket im Archiv liegt – 'backlog' liest sich wie offene Arbeit. Als DoD 5 aufgenommen. Testdaten erfunden, keine fremden Tickets.
- **2026-10-29 20:50 · Alexander Sacharov** — DoD 5: Entscheidung – Archiv und Logbook zeigen beide 'archive' (DoD-Text 'Archiv bzw. Logbook'); auch die Detailansicht, damit Karte und Detail nicht auseinanderlaufen. Die Zeile 'in the archive' aus Alex' Beschreibung gibt es in core/web weder hier noch auf origin/master – nichts zu entfernen.
- **2026-10-29 20:50 · Alexander Sacharov** — decision (Alex, 29.10.): DoD 6 – signoff gehört in den Epic-Fortschrittsbalken von jaira web (core/web/ui/src/entities/epic/ui/EpicProgress.svelte, heute nur done/total). Grund: Tickets warten lange in signoff und zählen sonst wie offene Arbeit. done bleibt die einzige Zahl, die als fertig zählt; signoff ist ein zweites Segment.
- **2026-10-29 20:52 · Alexander Sacharov** — DoD 6: Lane per Name 'signoff' (const in web_source.go) – die Lane-Konfiguration hat kein Merkmal, das signoff vom human-Lane unterscheidet; ein Board ohne signoff zählt 0. Signoff-Zahl nach der Ref-Ersetzung, zählt also auch Tickets, deren Datei noch in-progress sagt. Text ohne signoff unverändert '{done} of {total} done'. Farbe neues Token --waiting (hell #c98a1b / dunkel #e8b14f).
