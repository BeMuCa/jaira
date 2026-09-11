---
id: 01M26F9333EWZS2TN8EHPTQ3XT
title: "Ein Snapshot-Branch traegt die Tafel als Dateien, ohne dass jemand ihn auscheckt"
status: done
ready: true
creator: Alexander Sacharov
goal: "Ein elternloser Branch jaira/board haelt zu jedem Zeitpunkt genau die Tickets, die gerade auf Refs liegen, wird per Plumbing ohne Checkout geschrieben und kann nie mit den Arbeitsdateien kollidieren"
context: |-
  Sobald ein Ticket nur noch auf seinem Ref lebt und erst beim Uebernehmen als Datei landet (siehe RFC7GA), liegt der unbearbeitete Backlog nur auf dem Remote. Geht das Remote verloren, ist er weg - heute steht er in der git-Historie.

  Der Snapshot-Branch schliesst das, und zwar als Opt-in: er ist ein Backup und nicht der Speicherort.

  Gemessen, nicht vermutet - der ganze Ablauf wurde von Hand mit Plumbing durchgespielt:

  - Das Baumobjekt wird mit hash-object/mktree/commit-tree gebaut, genau wie core/gitref es fuer ein Ticket-Ref schon tut. Kein Checkout, kein Anfassen des Arbeitsbaums: der Test lief auf master, git status blieb leer.
  - Der erste Snapshot-Commit hat keinen Parent (elternlos), jeder weitere den vorigen. Damit ist 'git log jaira/board' die Geschichte der Tafel und 'git diff jaira/board~1 jaira/board' sagt, was sich geaendert hat. Belegt: 'board/01AAA.md | 6 ------ / board/01CCC.md | 6 ++++++'.
  - Hinzufuegen und Entfernen ergeben sich von selbst, ohne Vergleich: der Baum wird jedes Mal aus dem AKTUELLEN Satz Refs neu gebaut. Ein Ticket, dessen Ref beim Logbuch geloescht wurde, ist im naechsten Snapshot einfach nicht mehr drin - und bleibt in den vorigen Commits auffindbar.

  Zwei Entscheidungen, die nicht offensichtlich sind und beide Gruende haben:

  1. Die Dateien liegen unter board/<id>.md und NICHT unter .jaira/tickets/. Wuerden sie am selben Pfad liegen, bekaeme der erste, der diesen Branch versehentlich merged, einen add/add-Konflikt ueber alle Tickets auf einmal - die schlimmere Version des Problems, von dem dieser Entwurf weggeht. Ein anderer Pfad macht die Kollision konstruktiv unmoeglich.
  2. Angehaengt statt force-gepusht. Ein Snapshot ist ein Commit mit dem vorigen als Parent, damit es Historie gibt und zwei gleichzeitige Snapshots mit demselben --force-with-lease aufloesen, das core/gitref schon benutzt.

  Bekannte Schwaeche, ausdruecklich: ein Snapshot ist ein Stand von damals. Er braucht einen Befehl und optional einen Timer, sonst veraltet er. Das Arbeitsstand liegt immer auf den Refs; der Snapshot ist das Backup.
definition-of-done: "'jaira snapshot' baut den Branch aus dem aktuellen Satz refs/jaira/tickets/* per hash-object, mktree und commit-tree, ohne Checkout und ohne den Arbeitsbaum zu beruehren, und laesst sich in jedem Repo-Zustand aufrufen; der erste Commit ist elternlos, jeder weitere haengt am vorigen, sodass git log die Geschichte der Tafel ist; ein Ticket, dessen Ref verschwunden ist, fehlt im neuen Snapshot und bleibt in den vorigen Commits lesbar; die Dateien liegen unter board/<id>.md, damit ein versehentlicher Merge dieses Branches nie mit .jaira/tickets/ kollidiert; zwei gleichzeitige Snapshots loesen per --force-with-lease auf statt sich zu ueberschreiben; der Branchname ist konfigurierbar und der Befehl sagt, was er hinzugefuegt und entfernt hat"
tags:
  - concurrency
blocked-by: []
commits:
  - cc94119280855fc2230411228dcbdc7aaae894d6
  - 8bfff13b6395a732eb27efc29c2ef9c04f6e3165
  - 30555d9dc50ad33208a30787967bd56803cba146
  - b8e1cb1a0313c41a18fabe933a33b6fc8f3cb81a
  - fe864ec705086aafcd4b16fab0556d9dd2c3ceb7
  - 2a35f28cbb411e419354d97b6f6e8ef95db9c714
  - 3dbf02839cb83e92d9cb8f0e38eb225b454782ea
  - 44c13986ecafa96cb9760608649b13c97f9acb7a
created-at: 2026-09-10T20:13:52Z
updated-at: 2026-09-11T13:31:47Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-155213
claimed-at: 2026-09-11T11:27:20Z
assignee: Alexander Sacharov
question: "Snapshot steht: elternloser Branch jaira/board, alle drei Tage im Hintergrund, board/<id>.md, kein Commit wenn sich nichts aendert. Zwei Fragen: (1) Reicht dir der Default von drei Tagen, oder soll er kuerzer sein? (2) Soll 'jaira snapshot' auch etwas sagen, wenn nichts zu tun war - im Moment sagt es 'already current'."
outcome-what: "Review-Durchlauf gegen die DoD abgeschlossen"
outcome-why: "Diff erfuellt beide DoD-Punkte, ein kleiner Test-Gap notiert"
outcome-resolves: "PTQ3XT"
review-summary: "core/snapshot baut jaira/board per hash-object/mktree/commit-tree neu aus dem aktuellen Satz Ticket-Refs, ohne Checkout: Dateien unter board/<id>.md, erster Commit elternlos, jeder weitere haengt am vorigen (Rev von refs/heads/<branch> als Parent/Lease), unveraendert wird gar nichts geschrieben (Tree-Hash-Vergleich), zwei Laeufe loesen per --force-with-lease auf, Branchname ueber settings.json konfigurierbar. Ausgeloest als abgekoppelter Hintergrundprozess alle drei Tage nach dem Muster der Update-Pruefung, nie auf dem Kommandopfad. Reaping (Ref loeschen) passiert direkt danach, nur fuer Tickets die in einer Landing-Branch (settings.json, Globs) gefunden werden. Zweiter Teil im selben Commit: logbook/archive loeschen den Ref nicht mehr sofort, sondern schreiben den Endzustand (RecordFiled), damit ein Ticket waehrend eines laufenden Reviews fuer alle sichtbar bleibt; fetch/validate melden gestrandete Tickets (nie gelandet) nach der Gnadenfrist."
review-gaps: "Keine dedizierte Test-Abdeckung fuer den Konfliktfall zweier gleichzeitiger 'jaira snapshot'-Laeufe auf denselben Branch (--force-with-lease-Pfad in PushBranch) - der Mechanismus teilt sich Code mit dem laengst getesteten Ticket-Ref-Push, aber ein eigener Test dafuer fehlt. Sonst nichts: DoD-Punkte sind einzeln im Code und in Tests nachweisbar, NOTES.md wurde bereits in einem fruehen Commit dieser Serie geschrieben, build und alle betroffenen Pakete sind gruen."
review-verdict: "Diff erfuellt beide DoD-Punkte des Tickets; ein kleiner Test-Gap beim Branch-Force-Push-Konflikt, kein Grund zur Ablehnung."
review-check: "1. go build ./... im Repo-Root - muss ohne Fehler durchlaufen  2. go test ./core/snapshot/... ./core/gitref/... ./core/refsync/... ./core/settings/... ./internal/cli/... - alle ok  3. go test ./core/snapshot/ -run TestTheFirstSnapshotIsParentless -v - zeigt, dass der erste Commit ohne Parent ist  4. go test ./core/snapshot/ -run TestRemovingARefRemovesItFromTheNextSnapshotAndKeepsItInHistory -v - zeigt Added/Removed und dass ein entfernter Ref in alten Commits lesbar bleibt  5. go test ./core/refsync/ -run TestAFinishedTicketThatNeverLandedIsReported -v - zeigt, dass Filing den Ref nicht mehr sofort loescht"
---

# Ein Snapshot-Branch traegt die Tafel als Dateien, ohne dass jemand ihn auscheckt

## Definition of Done

- [x] 'jaira snapshot' baut den Branch aus dem aktuellen Satz refs/jaira/tickets/* per hash-object, mktree und commit-tree, ohne Checkout und ohne den Arbeitsbaum zu beruehren, und laesst sich in jedem Repo-Zustand aufrufen; der erste Commit ist elternlos, jeder weitere haengt am vorigen, sodass git log die Geschichte der Tafel ist; ein Ticket, dessen Ref verschwunden ist, fehlt im neuen Snapshot und bleibt in den vorigen Commits lesbar; die Dateien liegen unter board/<id>.md, damit ein versehentlicher Merge dieses Branches nie mit .jaira/tickets/ kollidiert; zwei gleichzeitige Snapshots loesen per --force-with-lease auf statt sich zu ueberschreiben; der Branchname ist konfigurierbar und der Befehl sagt, was er hinzugefuegt und entfernt hat
  proof: core/snapshot: Baum aus den Ref-Blobs per Plumbing, elternloser erster Commit, unveraenderter Baum-Hash schreibt nichts; Smoke: Arbeitsbaum unberuehrt, git log jaira/board ist die Geschichte
- [x] der Tree wird billig neu berechnet und nur gepusht wenn sich sein Hash geaendert hat; ausgeloest wird das im Hintergrund, wenn der letzte Snapshot aelter als drei Tage ist (Intervall konfigurierbar), nach dem Muster der Update-Pruefung - nie auf dem Kommandopfad, nie als Timer, den jemand von Hand einrichten muss, und niemals bei jedem Schreibvorgang
  proof: Due/SpawnRun: Stempel vor dem Lauf, Kind mit JAIRA_NO_SNAPSHOT=1, kein Wait, stdio auf DevNull; TestDueOnlyWhenTheLastSnapshotIsOldEnough

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/snapshot: Baum aus dem aktuellen Satz Refs bauen (board/<id>.md), Blobs per cat-file --batch-check aus den Refs nehmen statt neu zu hashen - sie liegen schon im Repository
  proof: core/snapshot/snapshot.go (Runner.Run: BlobsOf/Tree/Nest/Commit, Tree-Hash-Vergleich statt Commit, PushBranch mit --force-with-lease) und core/snapshot/schedule.go (Stamp/Due/SpawnRun nach dem Muster der Update-Pruefung); internal/cli/snapshot.go; 6 Tests in core/snapshot
- [x] Commit anhaengen: erster elternlos, jeder weitere am vorigen; ist der Baum-Hash unveraendert, gar kein Commit und kein Push
  proof: core/snapshot/snapshot.go (Runner.Run: BlobsOf/Tree/Nest/Commit, Tree-Hash-Vergleich statt Commit, PushBranch mit --force-with-lease) und core/snapshot/schedule.go (Stamp/Due/SpawnRun nach dem Muster der Update-Pruefung); internal/cli/snapshot.go; 6 Tests in core/snapshot
- [x] Push mit --force-with-lease auf den Snapshot-Branch, damit zwei gleichzeitige Snapshots sich nicht ueberschreiben
  proof: core/snapshot/snapshot.go (Runner.Run: BlobsOf/Tree/Nest/Commit, Tree-Hash-Vergleich statt Commit, PushBranch mit --force-with-lease) und core/snapshot/schedule.go (Stamp/Due/SpawnRun nach dem Muster der Update-Pruefung); internal/cli/snapshot.go; 6 Tests in core/snapshot
- [x] Ausloesen im Hintergrund, wenn der letzte Snapshot aelter als drei Tage ist (Intervall in settings.json), nach dem Muster der Update-Pruefung - nie auf dem Kommandopfad
  proof: core/snapshot/snapshot.go (Runner.Run: BlobsOf/Tree/Nest/Commit, Tree-Hash-Vergleich statt Commit, PushBranch mit --force-with-lease) und core/snapshot/schedule.go (Stamp/Due/SpawnRun nach dem Muster der Update-Pruefung); internal/cli/snapshot.go; 6 Tests in core/snapshot
- [x] jaira snapshot: den Lauf von Hand ausloesen, und sagen was hinzugekommen und was entfernt wurde
  proof: core/snapshot/snapshot.go (Runner.Run: BlobsOf/Tree/Nest/Commit, Tree-Hash-Vergleich statt Commit, PushBranch mit --force-with-lease) und core/snapshot/schedule.go (Stamp/Due/SpawnRun nach dem Muster der Update-Pruefung); internal/cli/snapshot.go; 6 Tests in core/snapshot
- [x] Tests: zwei Laeufe, dazwischen ein Ref weg und eins neu; git log auf dem Branch ist die Geschichte, git diff nennt die Aenderung; unveraenderter Satz Refs erzeugt keinen zweiten Commit
  proof: core/snapshot/snapshot.go (Runner.Run: BlobsOf/Tree/Nest/Commit, Tree-Hash-Vergleich statt Commit, PushBranch mit --force-with-lease) und core/snapshot/schedule.go (Stamp/Due/SpawnRun nach dem Muster der Update-Pruefung); internal/cli/snapshot.go; 6 Tests in core/snapshot

## Progress
- **2026-09-10 20:22 · Alexander Sacharov** — Frage geklaert (Alexander): muss der Snapshot bei jeder Aenderung laufen? Nein - periodisch, und das ist nicht Faulheit, sondern die richtige Antwort. Drei Gruende, in der Reihenfolge ihres Gewichts:

1. Die Historie gibt es schon, und zwar genauer: jedes Ticket hat sein eigenes Ref mit eigener Historie ('git log refs/jaira/tickets/<id>'). Der Snapshot ist NICHT die Quelle der Geschichte, er ist das Backup. Ein Commit pro Schreibvorgang kauft also etwas, das bereits da ist, und zahlt mit Rauschen: hunderte Commits am Tag auf dem Branch.

2. Lokal kostet er fast nichts. Die Blobs liegen schon im Repository - es sind dieselben Objekte, die die Refs tragen -, also kommen pro Snapshot ein Tree und ein Commit hinzu. Und hat sich der Satz Refs nicht geaendert, hat der neu gebaute Tree denselben Hash wie der auf dem Branch: dann ist gar kein Commit noetig, null Objekte.

3. Teuer ist nur der Push, ein Netz-Roundtrip. Der darf aus demselben Grund nicht an jeder Kommandozeile haengen wie der Grund, aus dem 'jaira list' nie aufs Remote wartet.

Daraus die Regel fuer die Umsetzung: den Tree immer billig neu berechnen, aber nur pushen, wenn sich sein Hash geaendert hat; den Push an etwas haengen, das ohnehin ins Netz geht (nach einem schreibenden Kommando) oder an einen Timer im Bereich 10-15 Minuten. Ein zehn Minuten alter Backup-Stand ist kein Problem: der Arbeitsstand liegt immer auf den Refs, der Snapshot ist fuer den Fall 'Remote verloren'.
- **2026-09-10 20:37 · Alexander Sacharov** — Kadenz entschieden (Alexander): alle 3 Tage, nicht alle 10-15 Minuten. Sein Argument ist das bessere: die Refs liegen bei jedem Teilnehmer lokal, also ist die Tafel aus JEDEM Klon rekonstruierbar - der Snapshot ist nicht das Mittel gegen 'das Remote ist gerade ausgefallen'.

Wem er wirklich dient, sind zwei andere: wer zum ersten Mal klont (der hat keine Refs) und wer nie gefetcht hat. Beides sind langsame Faelle, keine Minutenfaelle. Und faellt der Server wirklich aus, muss sowieso jemand ueberlegen, wie die Tafel wieder erreichbar wird - dabei hilft ein zehn Minuten alter Stand nicht mehr als ein drei Tage alter.

Wichtiger als das Intervall: kein Timer, den jemand von Hand einrichten muss. Das Muster gibt es hier schon und es ist das richtige - die Update-Pruefung laeuft hoechstens einmal am Tag, im Hintergrund und NIE auf dem Kommandopfad (README: 'never on the command path: nothing you run ever waits on the network for this'). Ein Kommando merkt also, dass der letzte Snapshot aelter als drei Tage ist, und erledigt ihn im Hintergrund.
- **2026-09-10 21:05 · Alexander Sacharov** — Der Snapshot bekommt eine zweite Aufgabe, entschieden zusammen mit RA7PFE: er jaetet die Refs, deren Ticket in einer Landebranch angekommen ist.

Das ist kein Anhaengsel, sondern der Grund, warum es hierhin gehoert: der Snapshot hat das Ticket unmittelbar davor in den Snapshot-Branch geschrieben, also liegt es im Moment des Loeschens in zwei Ablagen - im Snapshot und in der gemergten Branch. Ein Loeschen an jeder anderen Stelle haette diese Garantie nicht.

Reihenfolge im Lauf, nicht beliebig: erst schreiben, dann jaeten. Und was sich nicht aufloest (keine Landebranch bestimmbar), wird nicht angefasst.
- **2026-09-10 21:20 · Alexander Sacharov** — Gebaut. Drei Dinge, die beim Bauen entschieden oder gefunden wurden:

1. Die Blobs werden NICHT neu gehasht. Sie liegen schon im Repository, hingelegt von den Ref-Schreibungen, also holt der Lauf sie mit einem cat-file --batch-check und baut nur Tree und Commit dazu. Ein Snapshot kostet damit zwei Objekte, egal wie viele Tickets.

2. Kein Commit, wenn sich nichts geaendert hat, und zwar an der einzigen Stelle gepruefft, an der es billig ist: der neu gebaute Tree-Hash gegen den Tree des Branches. Gleicher Hash heisst gleicher Inhalt, ohne Vergleich der Dateien.

3. Der Lauf ist die einzige Stelle, die den Remote nach seinem HEAD fragen darf ('git remote set-head origin -a'), wenn die Landebranch nicht aufloesbar ist. Beim Klon eines leeren Repos gibt es origin/HEAD nicht, und genau so ein Board ist der Normalfall bei jaira init - im Smoke-Test war es zuerst nicht aufgeloest und es wurde korrekt NICHTS gejaetet, danach hat der Lauf es selbst gesetzt.

Hintergrundlauf ist eine Kopie des Musters aus der Update-Pruefung, mit denselben Gruenden: Stempel VOR dem Lauf (ein gescheiterter Lauf darf nicht beim naechsten Kommando erneut ausloesen), Kind mit JAIRA_NO_SNAPSHOT=1 als Rekursionsschutz, kein Wait, stdio auf DevNull, und nie aus einem Testbinary.
- **2026-09-11 11:29 · Alexander Sacharov** — Review-Durchlauf: DoD gegen den Diff (fe864ec) geprueft, Code gelesen (core/snapshot, core/gitref, core/refsync, internal/cli/snapshot.go, internal/cli/fetch.go, core/ticket/store.go), build und relevante Testpakete gruen. Einziger gefundener Gap: --force-with-lease-Konflikt auf dem Snapshot-Branch selbst hat keinen dedizierten Test.
