---
id: 01M64NQ6XEX4WFG855H5ZKPT4H
title: "Web zeigt den Stand vom Ref, Erstelldatum und wie lange ein Ticket unberührt ist"
status: review
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
updated-at: 2026-10-29T20:48:40Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1236953
claimed-at: 2026-10-29T20:37:14Z
outcome-what: "jaira web zeigt bei einer Ticketdatei, die älter ist als ihr Ref, den Stand vom Ref (Karte, Epic-Zählung, Detail); jede Karte zeigt Erstelldatum und 'untouched Nd/Nh'."
outcome-why: Lanes ohne Code schreiben nur den Ref; die Datei im Checkout blieb stehen und die Seite zeigte Tickets in signoff als in-progress.
outcome-resolves: Seite und Ref widersprechen sich nicht mehr; liegengebliebene Tickets sind auf der Karte erkennbar.
review-summary: "jaira web liest jetzt alle Ticket-Refs in einem git-Aufruf (neues refsync.Syncer.Newer über gitref.ReadMany) und ersetzt jede Ticketdatei, deren updated-at älter ist als der Ref, durch den Ref-Stand, bevor der Link-Index gebaut wird — dadurch zeigen Karten, Epic-Zählung und Detailansicht die Lane vom Ref; eine neuere Datei bleibt, weil derselbe isStale-Vergleich wie bei pull/CheckFresh entscheidet. Zusätzlich trägt web.Ticket created_at/updated_at, und TicketCard zeigt 'created <Datum>' und 'untouched 3d/5h/<1h' (im Browser aus updated_at berechnet). dist ist aus dem Quellstand gebaut (Neubau ergibt dieselben Hashes)."
round: 2
review-gaps: "Kein Blocker. (1) Der Test prüft nur lose Karten (/api/unassigned) und die Detailansicht; der Epic-Pfad (Epics(): Zählung und Daten, inkl. LoadAnywhere für abgelegte Mitglieder) ist ungetestet, obwohl das Outcome 'Epic-Zählung' nennt. (2) age()/date() haben keinen Test — es gibt im Repo kein UI-Testframework, das räumt der Implementierer selbst ein. (3) Der ersetzte Ticketstand ist ReadOnly mit leerem Path; im Link-Index gilt er damit als 'nur auf dem Ref' (PlaceRef), obwohl eine Datei existiert — die Seite zeigt Place heute nicht, fällt also nicht auf, ist aber eine stille Bedeutungsverschiebung. (4) Scheitert ReadMany, zeigt die Seite ohne Hinweis wieder den Dateistand. (5) TUI und 'jaira list' lesen weiter nur die Datei — bewusst außerhalb des Tickets (Note vom 29.10.)."
test-verdict: "pass: go build ./... und go test -race -count=1 ./... RC=0; dist == npm ci && npm run build; DoD 1+2 per TestWebShowsTheNewerOfFileAndRef (Gegenprobe ohne Newer rot), DoD 3+4 per Test (created_at/updated_at) und age()/date() unter node geprüft (<1h, 5h, 1d, 3d, leer); jaira web --no-open auf diesem Board liefert Daten auf jeder Karte"
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
