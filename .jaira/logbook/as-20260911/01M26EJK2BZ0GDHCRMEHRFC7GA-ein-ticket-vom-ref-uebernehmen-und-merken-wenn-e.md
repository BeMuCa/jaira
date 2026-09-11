---
id: 01M26EJK2BZ0GDHCRMEHRFC7GA
title: "Ein Ticket vom Ref uebernehmen, und merken wenn es die Tafel verlassen hat"
status: done
ready: true
creator: Alexander Sacharov
goal: "Ein Ticket lebt bis zur Uebernahme nur auf seinem Ref, und erst 'jaira pull' legt es hier als Datei hin - damit existiert es zu jeder Zeit in genau einem Klon und ein Merge kann es nicht doppeln"
context: |-
  Die Konflikte, die beim Durchsprechen von 8566KF gefunden wurden, sollen nicht aufgeloest, sondern unmoeglich gemacht werden. Beides ist gemessen, nicht vermutet.

  Was heute passiert:

  1. Ein Ticket, das nur auf einem Ref liegt, ist in 'jaira list' und auf dem Board keine Karte, und es gibt keinen Befehl, es auf die eigene Platte zu holen - man muesste 'git show refs/jaira/tickets/<id>:<id>.md' von Hand hinschreiben.
  2. Legt ein Klon das Ticket ins Logbuch, waehrend ein anderer es noch unter tickets/ hat, existiert es nach dem Merge zweimal (fuer git ist das rename/modify). Erkannt wird das nicht: Store.idIndex in core/ticket/store.go vergleicht Ids nur innerhalb von tickets/, weil Paths() nur dieses Verzeichnis liest.

  Der Entwurf, auf den wir gekommen sind: solange niemand das Ticket bearbeitet, ist das Ref der Speicherort und die Datei existiert nirgends. Das Uebernehmen legt sie hin - und weil ein Uebernehmen ein Schreibvorgang aufs Ref und damit ein Compare-and-Swap ist, gewinnt genau einer. Es bleibt kein Ehrenwort ('jeder zieht es nur einmal'), sondern die Mechanik erzwingt es. Ein zweiter Puller wird abgelehnt und erfaehrt, wer schneller war.

  Damit gibt es zu jeder Zeit genau eine Ticketdatei. add/add und die Logbuch-Dublette treten nicht mehr auf, statt behandelt zu werden.

  Drei Preise, ausdruecklich abgewogen und nicht uebersehen:

  1. 'Klonen und die Tafel sehen' gilt nicht mehr wortwoertlich. Der Default-Refspec zieht die Refs nicht, also wird aus 'clone, dann jaira' ein 'clone, dann jaira fetch'. Die README-Zeile muss mitgeaendert werden statt still falsch zu werden.
  2. Der unbearbeitete Backlog liegt nur auf dem Remote. Abgefedert durch die Regel 'die Datei wird bei pull und beim Verlassen der Tafel committet' - alles, woran gearbeitet wurde, und alles Abgeschlossene steht damit in der Historie. Fuer den unberuehrten Backlog ist der Snapshot-Branch aus PTQ3XT das Backup.
  3. Zwei Modi: ein Board mit Remote hat das Ref als Speicherort, ein Board ohne Remote (jaira init gitignored .jaira/) die Datei. Diese Verzweigung muss an einer Stelle stehen und benannt sein, sonst wird sie ueberall halb nachgebaut.
definition-of-done: "'jaira pull <id>' holt ein Ticket, das nur auf seinem Ref liegt, als Datei nach .jaira/tickets/ und setzt in derselben Operation den Uebernehmer als assignee, mit CAS - verliert der Aufrufer das Rennen, wird nichts hingelegt und die Meldung nennt, wer schneller war; 'jaira fetch' bleibt reines Lesen und legt nie eine Datei an; ein Ticket, dessen Ref verschwunden ist, waehrend die Datei noch unter tickets/ liegt, wird von 'jaira fetch' und 'jaira validate' als 'hat die Tafel verlassen' gemeldet, mit dem Befehl der es hier nachzieht, und nie still verschoben; core/ticket erkennt zwei Dateien mit derselben Id auch ueber tickets/, logbook/ und archive/ hinweg und meldet sie wie die bestehende Dublette"
tags:
  - concurrency
  - cli
blocked-by: []
commits:
  - 7cd9cd046d6488f72821e060081bebf3a0982197
  - cc94119280855fc2230411228dcbdc7aaae894d6
  - 8bfff13b6395a732eb27efc29c2ef9c04f6e3165
  - 30555d9dc50ad33208a30787967bd56803cba146
  - 19419ee2305171a728b526c4852c032d01112ab0
  - ef4b2df0fbde5449a58c40fe08502155ef345e72
  - bf7cde877747ef8973e54d589f2c5d9eead1b525
  - 2f5713f18c2fb6bec663c0329c29b27c9be563db
  - f2364213c06c1d2bbaa66ee0e8df8674523393fe
  - 3dbf02839cb83e92d9cb8f0e38eb225b454782ea
  - 675d8b7f1a47050f576ffa2974eb5c3e26a5256f
  - 231606dcce8401caa2480d855b8c62ce4f4074ba
created-at: 2026-09-10T20:01:34Z
updated-at: 2026-09-11T13:31:43Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-149957
claimed-at: 2026-09-11T11:25:14Z
assignee: Alexander Sacharov
question: "Alles aus der DoD steht und ist gegen echtes git belegt. Zwei Dinge fuer dich: (1) Nimmst du an, dass nichts still verschoben wird - ein woanders abgeraeumtes Ticket und eine Dublette zwischen Tafel und Logbuch werden nur gemeldet, mit dem Befehl der es aufloest? (2) PTQ3XT (Snapshot-Branch) ist jetzt der einzige offene Punkt der ganzen Konstruktion: soll ich direkt weitermachen, oder willst du erst diese vier Tickets abnehmen?"
outcome-what: "Review-Pass gegen die DoD von RFC7GA abgeschlossen"
outcome-why: "vier verbliebene Luecken aus dem Entwurf (stiller Steal, Dubletten ueber Verzeichnisse, verschwindendes Ref, Board-Karte statt Zaehler) code-seitig verifiziert, kein Defekt gefunden"
outcome-resolves: "review-summary/review-gaps/review-verdict/review-check gesetzt, Ticket wartet auf menschliches Signoff"
review-summary: "Vier Luecken aus dem Entwurf 8566KF/PP5SCQ geschlossen: (1) --steal ist nicht mehr still — core/refsync.Pull setzt beim Uebernehmen von jemand anderem eine Notiz auf das Ticket ('X took this ticket over from Y', ueber die neu extrahierte ticket.AppendNote) und core/refsync.Pulled traegt TakenFrom, das internal/cli/pull.go laut ausgibt. (2) core/ticket.Store.offBoardDuplicates vergleicht Ids jetzt auch gegen logbook/ und archive/, nicht nur innerhalb von tickets/, und meldet eine Dublette ueber PartialError mit beiden Fundorten und der Ansage 'nur du kannst sagen welches'. (3) Ein verschwundenes Ref, waehrend die Datei noch da ist, wird gemeldet statt verloren zu gehen: core/refsync.Departed vergleicht das seen-Register gegen die aktuellen Refs, und Incoming haelt einen verschwundenen Eintrag im seen-Register fest, solange die lokale Datei noch existiert — sonst haette der Fetch, der den Abgang melden soll, die einzige Evidenz vorher ueberschrieben. internal/cli/fetch.go und validate.go rufen Departed auf und drucken pro Ticket den Befehl, der es aufloest (logbook/archive/pull). (4) internal/tui zeigt ref-only Tickets jetzt als normale Karte mit einem 'pull'-Marker statt als blosse Zahl in der Statuszeile."
review-gaps: "none — go build ./... und go test ./core/... ./internal/... laufen gruen (inkl. der 11 core/refsync- und 2 neuen core/ticket-Tests aus diesem Diff). Code-Lesung bestaetigt: fileOnRefOnly (internal/cli/refs.go) ist tatsaechlich die einzige Stelle, die zwischen den zwei Speichermodi entscheidet; die README-Zeile 'clone and see the same board' ist auf 'clone, jaira fetch' korrigiert (README.md:13); ForgetDeparted loescht den seen-Eintrag, sodass ein einmal behandelter Abgang nicht erneut gemeldet wird."
review-verdict: "Der Diff erfuellt die Definition of Done Punkt fuer Punkt: pull/create/fetch-Trennung war bereits Stand von PP5SCQ, dieser Diff schliesst die drei verbliebenen Luecken (stiller Steal, Dubletten ueber Verzeichnisse, verschwindendes Ref) plus einen vierten, im Vorgaenger-Ticket zugesagten Punkt (Board-Karte statt Zaehler). Keine Defekte gefunden, alle Tests gruen."
review-check: |-
  1. cd /home/alex/projects/jaira und go test ./core/refsync/... ./core/ticket/... -run 'Steal|Departed|OffBoard|Duplicate' -v — zeigt TestStealingRecordsWhoItWasTakenFrom, TestATicketTakenOffTheBoardElsewhereIsReportedHere, TestATicketOnTheBoardAndInTheLogbookIsReported, TestATicketOnTheBoardAndInTheArchiveIsReported, alle PASS.
  2. Fuer die Handprobe: zwei Klone eines Test-Repos mit Remote anlegen, in Klon A 'jaira create' ein Ticket, in Klon B 'jaira fetch' dann 'jaira pull <id>' — Ausgabe '<id> is yours'.
  3. In Klon A danach 'jaira pull <id> --steal' — Ausgabe '<id> taken over from <B>' und die Ticketdatei enthaelt unter '## Progress' die Zeile '... took this ticket over from ...'.
  4. In Klon B 'jaira fetch' — die Meldung nennt den Uebernehmer.
  5. Fuer Dubletten: eine Ticketdatei sowohl unter .jaira/tickets/ als auch unter .jaira/logbook/<ordner>/ ablegen, dann 'jaira validate' oder 'jaira list' — Fehlermeldung nennt beide Fundorte und 'nur du kannst sagen welches'.
---

# Ein Ticket vom Ref uebernehmen, und merken wenn es die Tafel verlassen hat

## Definition of Done

- [x] 'jaira pull <id>' holt ein Ticket vom Ref als Datei nach .jaira/tickets/ und setzt in derselben Operation den Uebernehmer als assignee, mit CAS - verliert der Aufrufer das Rennen, wird keine Datei angelegt und die Meldung nennt, wer schneller war; 'jaira create' legt auf einem Board mit Remote keine lokale Datei mehr an, sondern nur das Ref, und sagt das; 'jaira fetch' bleibt reines Lesen; ein Board ohne Remote verhaelt sich unveraendert wie heute, und diese Verzweigung steht an genau einer Stelle im Code; ein Ticket, dessen Ref verschwunden ist, waehrend die Datei noch unter tickets/ liegt, wird von fetch und validate als 'hat die Tafel verlassen' gemeldet und nie still verschoben; core/ticket erkennt zwei Dateien mit derselben Id auch ueber tickets/, logbook/ und archive/ hinweg; die README-Zeile ueber 'clone and see the same board' ist auf den neuen Ablauf korrigiert
  proof: jaira pull/release/fetch, list+next+show sehen Ref-Tickets ([pull it], on-ref-only), Schreiben verweigert mit exit 3; core/ticket offBoardDuplicates fuer tickets+logbook+archive; refsync.Departed in fetch und validate; README mit Sequenzdiagramm; Board ohne Remote unveraendert (fileOnRefOnly ist die einzige Verzweigung). Tests: core/refsync (11), core/ticket (2 neue), Smoke ueber alle Kommandos

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/refsync.Pull(id): fetch, Ref lesen, assignee auf mich setzen, synchron mit CAS pushen - und NUR bei akzeptiertem Push die lokale Datei schreiben
  proof: core/refsync/refsync.go Pull + internal/cli/pull.go; Tests TestOnlyOneCloneCanPullATicket, TestPullingATicketYouAlreadyHaveIsANoOp; Smoke gegen zwei Klone
- [x] Reihenfolge begruenden und festhalten: Push zuerst, Datei danach. Umgekehrt haette der Verlierer des Rennens eine Datei, die niemandem gehoert
  proof: core/refsync/refsync.go Pull + internal/cli/pull.go; Tests TestOnlyOneCloneCanPullATicket, TestPullingATicketYouAlreadyHaveIsANoOp; Smoke gegen zwei Klone
- [x] Pull ist die eine Operation, die synchron schreibt und nicht ueber die Outbox geht - ohne Netz gibt es kein Ref zu lesen, also auch nichts zu uebernehmen
  proof: core/refsync/refsync.go Pull + internal/cli/pull.go; Tests TestOnlyOneCloneCanPullATicket, TestPullingATicketYouAlreadyHaveIsANoOp; Smoke gegen zwei Klone
- [x] Idempotenz: liegt die Datei schon hier, ist pull ein erfolgreicher No-op (Projektregel fuer wiederholbare Kommandos)
  proof: core/refsync/refsync.go Pull + internal/cli/pull.go; Tests TestOnlyOneCloneCanPullATicket, TestPullingATicketYouAlreadyHaveIsANoOp; Smoke gegen zwei Klone
- [x] internal/cli/pull.go: jaira pull <id>, --json, Exit 3 wenn das Rennen verloren ist, mit Namen des Gewinners
  proof: core/refsync/refsync.go Pull + internal/cli/pull.go; Tests TestOnlyOneCloneCanPullATicket, TestPullingATicketYouAlreadyHaveIsANoOp; Smoke gegen zwei Klone
- [x] Tests mit zwei Klonen: berk pullt und hat die Datei, ada verliert und hat KEINE Datei, doppeltes pull ist No-op
  proof: core/refsync/refsync.go Pull + internal/cli/pull.go; Tests TestOnlyOneCloneCanPullATicket, TestPullingATicketYouAlreadyHaveIsANoOp; Smoke gegen zwei Klone
- [x] README: eigener Abschnitt zum Ablauf mit einem mermaid-Sequenzdiagramm, das zeigt, wie zwei Leute ein Ticket austauschen (create nur aufs Ref, fetch liest, pull gewinnt per CAS, Verlierer bekommt die Meldung) - mermaid, weil GitHub es selbst zeichnet und es im Gegensatz zu einem png von Hand aenderbar bleibt
  proof: core/refsync Release + internal/cli/release.go; internal/cli/refs.go fileOnRefOnly (die einzige Stelle, an der die zwei Modi entschieden werden); Tests TestReleaseHandsTheTicketBack, TestAnAssignedTicketIsReservedForItsAssignee; Smoke: create=0 Dateien, pull, release, Fremd-pull
- [x] README: die Zeile 'clone and see the same board' auf den neuen Ablauf korrigieren, statt sie still falsch werden zu lassen
  proof: core/refsync Release + internal/cli/release.go; internal/cli/refs.go fileOnRefOnly (die einzige Stelle, an der die zwei Modi entschieden werden); Tests TestReleaseHandsTheTicketBack, TestAnAssignedTicketIsReservedForItsAssignee; Smoke: create=0 Dateien, pull, release, Fremd-pull
- [x] jaira release <id>: den assignee auf dem Ref per CAS loeschen und die lokale Datei entfernen - danach kann es jeder holen. Ohne das ist eine Zuweisung eine Reservierung, die niemand zurueckgeben kann
  proof: README.md Abschnitt 'How two people hand work over' mit mermaid-Sequenzdiagramm plus Tabelle 'What each step buys'
- [x] Reservierung dokumentieren: 'create --assignee berk' materialisiert bei berk NICHTS, es reserviert nur; berk sieht 'neu fuer dich' und holt es selbst. In welchem Branch er arbeitet, ist nicht Sache des Zuweisenden
  proof: README.md: 'a teammate clones, runs jaira fetch, and sees the same board' - an beiden Stellen (Einleitung und Start-Abschnitt) korrigiert; pull und release in der Kommandoliste
- [x] Board: ein Ticket, das mir zugewiesen ist und noch nicht hier liegt, ist eine eigene Karte mit 'pull' statt nur eine Zahl in der Hinweiszeile - ein zugewiesenes Ticket ist der Fall, fuer den das alles existiert
  proof: internal/cli/tickets.go: Zeilenmarker [pull it]; ticketJSON traegt on-ref-only; alle Lesekommandos gegen ein Board mit ausschliesslich Ref-Tickets durchgeprueft
- [x] --steal sagt laut, von wem genommen wurde, und schreibt es in eine Notiz am Ticket
  proof: core/refsync Pulled.TakenFrom + ticket.AppendNote (aus internal/cli/resume.go herausgezogen, jetzt eine Stelle fuer beide Aufrufer); internal/cli/pull.go sagt es laut; TestStealingRecordsWhoItWasTakenFrom

## Progress
- **2026-09-10 20:14 · Alexander Sacharov** — Reihenfolge, in der das gebaut werden muss: erst pull (dieses Ticket), dann darf create aufhoeren, lokal zu schreiben. Umgekehrt gaebe es einen Zustand, in dem ein Ticket auf einem Ref liegt und niemand es holen kann.

PTQ3XT (Snapshot-Branch) ist die Abfederung von Preis 2 und kann parallel laufen - es haengt nicht an diesem Ticket, aber dieses Ticket sollte nicht ohne es in Produktion gehen, sonst liegt der Backlog eine Zeit lang ohne Backup nur auf dem Remote.
- **2026-09-10 20:23 · Alexander Sacharov** — pull steht. Zwei Dinge, die der Test gefunden hat und die im Entwurf falsch waren:

1. CAS allein sichert das Uebernehmen NICHT. Pull liest das Ref unmittelbar vor dem Schreiben, also ist der Lease immer aktuell und der Push gelingt immer - im ersten Wurf hat ada berks Uebernahme einfach ueberschrieben, und der Test hat es sofort gezeigt. Der eigentliche Waechter ist der assignee AUF DEM REF: ein Ticket, das jemand schon geholt hat, wird mit ErrTaken abgelehnt, mit Namen. Das CAS bleibt trotzdem noetig - es deckt den Fall ab, den die Pruefung nicht sehen kann: zwei Klone, die im selben Moment ein noch unbesetztes Ticket lesen.

2. Ein aufgenommenes Ticket gehoert niemandem. Meine Testfixture setzte assignee beim Anlegen, und damit war jeder pull sofort 'taken'. Das ist nicht nur ein Fixture-Fehler, es widerspricht der Regel, auf der der ganze Ablauf steht (README: 'a captured ticket belongs to nobody; whoever pulls it out of the backlog becomes its assignee'). assignee heisst 'das ist meins zu bearbeiten', nicht 'ich habe es geschrieben' - dafuer ist creator da.

Reihenfolge im Code ist die halbe Mechanik: erst der Push, dann die Datei. Andersherum haelt der Verlierer eines Rennens eine Ticketdatei, die jemand anderem gehoert - genau die Dublette, die dieser Entwurf unmoeglich machen soll, statt sie aufzuloesen. Belegt: nach der Ablehnung hat ada 0 Dateien.

pull ist die eine Operation, die synchron schreibt und nicht ueber die Outbox geht. Das ist keine Ausnahme von der Offline-Regel, sondern deren Folge: ohne Route zum Remote gibt es kein Ref zu lesen, also nichts zu uebernehmen.

Handprobe mit dem Binary: berk pullt und hat die Datei; ada wird abgelehnt ('berk has it, it is in backlog', exit 3) und hat keine; ein zweiter pull bei berk ist ein No-op; --steal nimmt es trotzdem.
- **2026-09-10 20:38 · Alexander Sacharov** — Preis 1 entschieden (Alexander): 'Klonen und die Tafel sehen' wird zu zwei Schritten, 'git clone && jaira fetch', und danach ist es fuer alle gleich. Das ist die ganze Loesung - keine Sonderbehandlung, keine Materialisierung beim Clone, nur eine Zeile mehr in der Anleitung. Die README-Zeile wird entsprechend geaendert, nicht weggelassen.

Damit bleibt von den drei Preisen im Kontext nur noch einer, der Arbeit macht: die zwei Modi (mit Remote das Ref, ohne Remote die Datei). Preis 2 (Backlog nur auf dem Remote) traegt PTQ3XT, alle drei Tage im Hintergrund.
- **2026-09-10 20:40 · Alexander Sacharov** — create und release stehen. Die Verzweigung der zwei Modi liegt an genau einer Stelle: cli.fileOnRefOnly. Ein Board mit Remote schickt das frisch angelegte Ticket aufs Ref und nimmt die lokale Datei wieder weg; ein Board ohne Remote behaelt sie, unveraendert wie bisher. Jedes andere Kommando liest, was da ist, und muss nicht wissen, in welchem Modus es laeuft.

Warum die Datei weggeht, waehrend niemand daran arbeitet: eine Datei in irgendeinem Checkout ist eine Kopie, die ein Merge doppeln kann, und sie verdeckt, wem das Ticket gehoert. 'jaira pull' holt sie zurueck, fuer genau einen Klon.

Eine Ausnahme, die bleiben muss: konnte das Ticket nicht gesendet werden, BLEIBT die Datei hier. Sie ist dann hier korrekt, traegt den unsent-Marker und geht mit dem naechsten Kommando raus. Die Datei fuer eine Schreibung wegzunehmen, die die Maschine nie verlassen hat, waere der eine Weg, ein Ticket wirklich zu verlieren.

release loescht den assignee auf dem Ref per CAS und entfernt danach die Datei - in dieser Reihenfolge, damit ein Fehlschlag das Ticket noch dir laesst und nicht 'niemandem, aber noch auf deiner Platte'. Ein fremdes Ticket wird abgelehnt und nennt den Halter; --force fuer den Fall, dass die Person nicht zurueckkommt.

Handprobe, Ablauf komplett: ada legt zwei Tickets an (0 Dateien bei ihr), berk fetcht und sieht beide als 'ref-only', zieht sein zugewiesenes, gibt es zurueck (Datei weg), danach kann ada es ziehen. Ein Board ohne Remote verhaelt sich unveraendert.
- **2026-09-10 20:48 · Alexander Sacharov** — Der Sammelpunkt ist jetzt global, und der Audit dazu hat drei Fehler gefunden - alle vom selben Typ.

Wie es gemacht ist: der Store bekommt neben Recorder (nach draussen) eine Source (nach innen), und core/refsync erfuellt beide. List haengt die Tickets an, die das Board sieht ohne eine Datei zu haben; Load faellt darauf zurueck, wenn kein Pfad passt. Damit sehen list, next, show, tasks und das Board alles, ohne von Refs zu wissen - genau wie beim Schreiben. Ein Aufrufer, der eine zweite Stelle fragen muesste, waere ein Aufrufer, der die halbe Tafel zeigt.

Gefunden beim Durchgehen ALLER Kommandos gegen ein Board, dessen Tickets nur auf Refs liegen:

1. archive, delete und logbook nahmen t.Path, der bei einem Ticket ohne Datei leer ist. filepath.Base("") ist ".", also war das Ziel das Archivverzeichnis selbst und archive antwortete mit Unsinn: 'archive already exists in the archive', exit 1. Jetzt eine Pruefung onlyOnRef in genau diesen drei Funktionen - jedes Kommando, das eine Ticketdatei bewegt oder entfernt, laeuft durch sie.

2. Meine eigene Source hat pull kaputt gemacht: Load antwortet jetzt auch fuer ein Ticket auf dem Ref, also hielt pull es fuer 'schon hier' und tat nichts (sichtbar als 'is already here:' mit leerem Pfad). Die Frage ist, ob eine DATEI da ist, nicht ob das Board das Ticket sieht - es sieht jedes Ref. Dafuer gibt es jetzt localPath/localIDs, die nur Dateinamen lesen. Dieselbe Falle steckte in Release und Reconcile.

3. validate meldete einen leeren Pfad. Steht jetzt als '(on its ref; pull it to work on it)' - der Unterschied zwischen 'geh und reparier die Datei' und 'das ist hier nicht deins zu reparieren'.

Achtung fuer spaeter, das ist die Lehre: Load darf nicht mehr als 'ist es hier' gelesen werden. Wer einen Pfad braucht, muss ihn ueber die Dateinamen holen. Der Typ traegt das jetzt: Ticket.ReadOnly.

Ergebnis des Audits: alle Lesekommandos sehen die Tickets, alle Schreibkommandos verweigern mit exit 3 und einem Satz, der den naechsten Schritt nennt.
- **2026-09-10 20:56 · Alexander Sacharov** — Alles aus der Definition of Done steht. Was beim Bauen gefunden wurde und im Entwurf nicht stand, in der Reihenfolge des Gewichts:

1. CAS sichert das Uebernehmen nicht. Der Waechter ist der assignee auf dem Ref; das CAS deckt nur den Fall 'zwei lesen im selben Moment ein unbesetztes Ticket'. Der Test hat das sofort gezeigt, der Entwurf hatte es falsch.

2. Der Sammelpunkt musste global werden, sonst waren list und next auf einem frisch geklonten Board blind - 'No tickets match' bei voller Tafel. Loesung symmetrisch zum Schreiben: Store.Source neben Store.Recorder, erfuellt von core/refsync.

3. Genau daraus folgte die naechste Falle: Store.Load heisst seither nicht mehr 'ist es hier'. Es hat pull kaputt gemacht (hielt jedes Ref fuer 'schon hier'), und dieselbe Falle steckte in Release und Reconcile. Wer einen Pfad braucht, holt ihn ueber die Dateinamen (localPath). Der Typ traegt es jetzt: Ticket.ReadOnly.

4. Der Audit ueber ALLE Kommandos hat drei Fehler vom Typ 'leerer Pfad' gefunden: archive, delete und logbook nahmen t.Path, filepath.Base("") ist ".", und archive antwortete mit 'archive already exists in the archive'. Jetzt eine Pruefung in genau diesen drei Funktionen.

5. Der Fetch, der einen Abgang melden soll, hat die einzige Evidenz vorher ueberschrieben. Ein verschwundenes Ref bleibt jetzt im seen-Register, solange die Datei da ist - unerledigt, nicht Geschichte.

Was absichtlich NICHT gemacht wurde: nichts wird still verschoben. Weder ein Ticket, das woanders abgeraeumt wurde, noch eine Dublette zwischen Tafel und Logbuch. Beides wird gemeldet mit dem Befehl, der es aufloest, weil ein Werkzeug, das hier raet, irgendwann fertige Arbeit wieder oeffnet.

Offen fuer PTQ3XT und nicht hier: der Snapshot-Branch. Solange er fehlt, liegt ein unberuehrter Backlog nur auf dem Remote - das ist Preis 2 aus dem Kontext, bewusst so, aber es sollte nicht lange so bleiben.
- **2026-09-11 11:27 · Alexander Sacharov** — Review gegen DoD: alle vier Punkte im Code verifiziert (nicht nur im Diff-Text) — fileOnRefOnly als einzige Verzweigung, README-Zeile korrigiert, Departed-Persistenz in Incoming, offBoardDuplicates ueber logbook+archive. go build und volle Testsuite gruen. Keine Gaps gefunden.
