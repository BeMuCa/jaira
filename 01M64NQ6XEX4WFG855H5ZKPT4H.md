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
updated-at: 2026-10-29T21:05:14Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1236953
claimed-at: 2026-10-29T20:37:14Z
outcome-what: "jaira web zeigt den Ref-Stand bei älterer Datei, Erstelldatum und 'untouched' auf jeder Karte, und 'archive' als Lane für archivierte/abgelegte Tickets."
outcome-why: Lanes ohne Code schreiben nur den Ref; abgelegte Tickets behalten ihre alte Lane und lasen sich wie offene Arbeit.
outcome-resolves: Seite und Ref widersprechen sich nicht mehr; Liegengebliebenes und Abgelegtes ist auf der Karte erkennbar.
review-summary: "board() in internal/cli/web_source.go liest jetzt alle Ticketdateien, fragt mit dem neuen refsync.Syncer.Newer in einem einzigen git-Durchgang (Repo.ReadMany) alle Refs ab und ersetzt jede Datei, deren Ref ein späteres updated-at hat (isStale), durch den Ref-Stand – bevor der Index gebaut wird. Damit stehen Karten in Epics und 'unassigned' sowie die Detailansicht (Ticket() ruft Newer für das eine Ticket) in der Lane des Refs; eine neuere Datei gewinnt weiter. web.Ticket bekommt created_at/updated_at, die Karte zeigt 'created <Datum> · untouched 3d/5h/<1h' (time.ts age/date). Archivierte bzw. Logbook-Tickets zeigen 'archive' statt ihrer alten Lane (Epics über ix.Lookup-Place, Detail über den Pfad außerhalb tickets/). Epics zählen signoff getrennt (Epic.Signoff) und der Fortschrittsbalken zeigt dafür ein bernsteinfarbenes zweites Segment mit Text '4 done · 3 in signoff of 10'. Release-Note core/release/notes/ZKPT4H.md mit vier Zeilen, dist-Bundle neu gebaut und passt zum Quelltext."
round: 5
review-gaps: "Keine Defekte gefunden. Kleinigkeiten: (1) board() dupliziert boardIndex() aus internal/cli/links.go:113 bis auf den Newer-Schritt – zwei Kopien derselben Ladelogik. (2) Ein ticket, das der Web-Server als Ref-Stand zeigt, bleibt im Checkout veraltet; die Seite verschweigt das (kein Hinweis 'Datei älter als Ref') – gewollt laut Kontext, aber unsichtbar. (3) Fortschrittsbalken: pct(done)+pct(signoff) wird einzeln gerundet und kann 101% ergeben; durch overflow:hidden unsichtbar. (4) Der Signoff-Zähler und 'archive' sind nur über Go-Tests abgedeckt; die Svelte-Anzeige (Karte, Balken) hat keinen Test."
test-verdict: "pass: go test -race -count=1 ./... RC=0 nach DoD 5/6; dist == Neubau aus src; DoD 5 per TestWebShowsAnArchivedTicketAsArchive, DoD 6 per TestWebCountsSignoffApartFromDone (beide Gegenproben rot), DoD 1-4 wie Runde 1; jaira web auf diesem Board liefert signoff je Epic (hier 0, keine signoff-Tickets im Epic)"
review-verdict: "Der Diff erfüllt alle sechs DoD-Punkte; TestWebShowsTheNewerOfFileAndRef deckt genau Datei=in-progress/Ref=signoff und die Gegenrichtung ab. go vet sauber, go test ./internal/cli ./core/refsync ./core/web ./core/ticket ./core/link ./core/epic grün. Keine Defekte; die Lücken sind kosmetisch bzw. Wartung."
review-round: 4
review-check: "1. go test ./internal/cli -run 'TestWebShowsTheNewerOfFileAndRef|TestWebShowsAnArchivedTicketAsArchive|TestWebCountsSignoffApartFromDone' -count=1 → ok  2. go build -o /tmp/jaira ./cmd/jaira und in einem Repo mit geteiltem Board /tmp/jaira web starten  3. ein Ticket nehmen, das per 'jaira move' (ohne Code-Commit) in signoff steht, dessen Datei im Checkout aber noch in-progress sagt: die Karte zeigt 'signoff'  4. jede Karte trägt darunter 'created <Datum> · untouched <N>d|h|<1h'  5. ein Epic mit einem signoff-Ticket: der Balken hat ein bernsteinfarbenes Segment, Text '<x> done · <y> in signoff of <z>'  6. ein mit 'jaira archive' abgelegtes Epic-Mitglied: Karte und geöffnetes Detail zeigen Lane 'archive'"
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
- [~] board() in internal/cli/web_source.go und boardIndex() in internal/cli/links.go sind eine Ladefunktion
- [ ] Karte und Detailansicht im Web sagen, wenn die Datei im Checkout älter ist als der Ref und die Seite den Ref zeigt

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
- **2026-10-29 20:54 · Alexander Sacharov** — critique (DoD 5/6): stehen gelassen – signoff per Lane-Name statt per Lane-Merkmal (nur RequiresNonModelSignal && !RequiresQuestion trennt es von human, mehr Umweg für dieselbe Antwort; DoD nennt signoff wörtlich). Ticket(): 'Pfad außerhalb tickets/' statt link.Index für die Detailansicht, weil ein Index für ein einzelnes Ticket das ganze Board liest.
- **2026-10-29 20:54 · Alexander Sacharov** — testing: auf diesem Board zeigt ein Epic 3 Karten mit leerer Lane – Mitglieder, die hier unbekannt sind (Lookup → PlaceUnknown). War vor diesem Ticket genauso, nicht Teil davon.
- **2026-10-29 21:02 · Alexander Sacharov** — explain: Der Review fand keinen Defekt, nur Kleinigkeiten: board() in web_source.go kopiert die Ladelogik von boardIndex() (links.go:113), die Seite verschweigt, dass die Datei im Checkout älter als der Ref ist, der Balken kann durch Rundung 101% ergeben (unsichtbar), und die Svelte-Anzeige hat keinen Test.
Warum: Alle sechs DoD-Punkte sind erfüllt und die Tests sind grün, die Lücken betreffen Wartung und Kosmetik. Zwei Kopien derselben Ladelogik können später auseinanderlaufen.
1) Zurückgeben: board() und boardIndex() zusammenführen und einen Hinweis 'Datei älter als Ref' ergänzen [give back]
2) So annehmen und die Kleinigkeiten bewusst stehen lassen [accept]
Empfehlung: 2 - kein Defekt, und jede Änderung vergrößert das Ticket, obwohl das Projekt klein bleiben soll.
- **2026-10-29 21:04 · Alexander Sacharov** — decision: 1) Zurückgeben: board() und boardIndex() zusammenführen und einen Hinweis 'Datei älter als Ref' ergänzen
- **2026-10-29 21:04 · Alexander Sacharov** — Given back from signoff: 1) Zurückgeben: board() und boardIndex() zusammenführen und einen Hinweis 'Datei älter als Ref' ergänzen
- **2026-10-29 21:04 · Alexander Sacharov** — cockpit: given back to the old worker - tab alive, cache warm for 54 more min, context 170k
- **2026-10-29 21:04 · Alexander Sacharov** — Runde 3, Entscheidungen: (a) Zusammenführen heißt: boardIndex (links.go) macht den Newer-Schritt selbst und meldet, welche Tickets ersetzt wurden; webSource.board ruft es. Damit zeigen auch 'jaira epic ls/show' und 'jaira links' den Ref-Stand – derselbe Fehler, gleiche Regel; Release-Note dafür. (b) Newer setzt den ersetzten Tickets nicht mehr ReadOnly und behält den Dateipfad: sonst ordnet link.Build sie als PlaceRef ('nicht gepullt') ein, was in 'jaira links' sichtbar falsch wäre (Review-Lücke 3 der ersten Runde). Geschrieben wird über diese Objekte nie – Mutate liest per id neu. (c) Hinweis: Feld behind_ref an Karte und Detail; Karte zeigt 'file behind ref' in der Datumszeile, Tooltip nennt 'jaira pull <handle>'; Detail zeigt dieselbe Zeile unter dem Kopf.
