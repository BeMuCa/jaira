---
id: 01M59H6V1TS2W4SENBXDHEY6T7
title: "jaira web schreibt mit, was passiert, und ein Absturz hinterlaesst seine Ursache"
status: review
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wenn die Seite von 'jaira web' ausfaellt, steht danach in einer Logdatei, was zuletzt passiert ist und warum."
context: |-
  Was heute fehlt: 'jaira web' loggt nichts.
  - core/web/*.go und internal/cli/web.go: kein log-Aufruf, kein recover(). Eine Panic in einem Handler oder in srv.Loop (Fetch alle 60 s, internal/cli/web.go:24) beendet den Prozess; der Stacktrace geht nur nach stderr und ist weg, wenn web nicht in einem sichtbaren Terminal laeuft.
  Anlass: Alex am 19.10.2026 - 'Web faellt nach ein paar Minuten aus'.
  Schon geprueft (Teamlead, 19.10.2026):
  - jaira 0.5.4, 'jaira web --no-open' in einem anderen Repo 8 Minuten lang alle 20 s mit curl /api/epics abgefragt: immer 200, kein Absturz, stderr leer.
  - Alex' eigener 'jaira web' im selben anderen Repo lief beim Nachsehen seit ueber einer halben Stunde weiter - der PROZESS starb dort also nicht. Moeglich: die Seite im Browser verliert die Verbindung oder zeigt einen Fehler, waehrend der Server lebt.
  Was Alex genau sieht (Fehlermeldung, leere Seite, Prozess weg) steht als Note am Ticket, sobald bekannt.
definition-of-done: "Jede Anfrage mit Fehlerstatus (>=400), jeder fehlgeschlagene Fetch/Push und jede Panic wird mit Zeit und Ursache in eine Logdatei geschrieben; 'jaira web' nennt beim Start ihren Pfad"
tags:
  - web
blocked-by: []
related: []
commits: []
created-at: 2026-10-19T07:32:24Z
updated-at: 2026-10-19T20:51:09Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-11846
claimed-at: 2026-10-19T20:45:53Z
outcome-what: "jaira web schreibt Fehlerantworten (>=400), fehlgeschlagene Fetches/Pushes, Panics mit Stack und das beendende Signal mit Zeit nach ~/.jaira/state/<baum>/web.log und nennt den Pfad beim Start; eine Panic im Handler antwortet 500 statt die Verbindung zu kappen; die Seite zeigt nach spaetestens 10 s ein Banner, wenn der Befehl nicht mehr antwortet."
outcome-why: "Alex' 'jaira web' fiel aus, und es gab keine Spur: nichts wurde geloggt, stderr ging mit dem Terminal verloren, und die Seite stand still, ohne es zu sagen."
outcome-resolves: "Ein Ausfall von jaira web hinterlaesst jetzt seine Ursache in einer Datei, und die Seite meldet ihn sichtbar."
review-summary: "jaira web schreibt jetzt ein Log nach ~/.jaira/state/<arbeitsbaum>/web.log (beim Start ab 1 MB nach web.log.1 rotiert) und nennt den Pfad beim Start. Eine Middleware um alle Handler loggt jede Antwort >=400 mit dem Anfang des Fehlerkoerpers und faengt Panics ab: Stack ins Log, Antwort 500 statt gekappter Verbindung. Fetch und Push laufen hinter recover(), Fehlschlaege und Panics landen im Log und im Sync-Status, die Fetch-Schleife laeuft weiter. SIGTERM/SIGHUP beenden wie Ctrl-C, das Signal steht im Log. Laesst sich das Log nicht oeffnen, warnt der Befehl und loggt (inkl. http.Server-ErrorLog) nach stderr weiter. Die Seite fragt alle 10 s /api/sync; scheitert das, zeigt sie oben ein rotes Banner 'has not answered since HH:MM' und laedt beim naechsten Erfolg alles neu. Release Note vorhanden."
review-gaps: "none - Runde 3 aendert nur den Fallback-Logger (internal/cli/web.go:71); der io-Import fiel mit io.Discard weg"
test-verdict: "pass: go test -race -count=1 ./... RC=0; Proof-Tests (Panic/ErrorAnswers/FailedFetchAndPush) PASS; DoD-Proofs auf aktuelle Zeilen korrigiert (log.go:64, web.go:137,153, ReviewPage.svelte:96); live: schreibgeschuetztes State-Verzeichnis -> Warnung, Log-Zeilen auf stderr, /api/sync 200"
question: ""
review-verdict: "Der Diff erfuellt DoD 1-5, keine Defekte; die drei Rueckgabepunkte aus signoff (Kommentar, Warnen statt Abbrechen, Ticketdatei an HEAD) sind erledigt. Push erst nach Alex' History-Bereinigung (Note 14:34)."
review-round: 2
review-check: "1. go build -o ./jw ./cmd/jaira   2. ./jw web --no-open   3. Zweite Ausgabezeile lautet 'jaira web: log <pfad>/web.log' - Pfad und Port merken   4. Zweites Terminal: curl -s http://127.0.0.1:<port>/api/tickets/NOPE - Antwort ist 404 'not found'   5. tail <pfad>/web.log zeigt 'started on ...' und 'GET /api/tickets/NOPE: 404 {...not found}' mit Uhrzeit   6. Adresse im Browser oeffnen, dann im ersten Terminal Ctrl-C   7. Binnen 10 s erscheint oben ein rotes Banner 'jaira web has not answered since HH:MM'   8. tail web.log zeigt 'stopped by signal interrupt'   9. Panic-Pfad hat keinen Handweg: go test ./core/web -run 'TestPanic|TestFailedFetch' -v - beide PASS"
round: 3
---

# jaira web schreibt mit, was passiert, und ein Absturz hinterlaesst seine Ursache

## Definition of Done

- [x] Jede Anfrage mit Fehlerstatus (>=400), jeder fehlgeschlagene Fetch/Push und jede Panic wird mit Zeit und Ursache in eine Logdatei geschrieben; 'jaira web' nennt beim Start ihren Pfad
  proof: core/web/log.go:64 logged + core/web/web.go:137,153 + internal/cli/web.go:88; TestErrorAnswersAreLogged, TestFailedFetchAndPushAreLogged
- [x] Eine Panic in einem Handler wird abgefangen (500 + Stacktrace im Log), der Server laeuft weiter; Test
  proof: core/web/log.go:64 logged; TestPanicAnswers500AndIsLogged
- [x] Die Seite zeigt sichtbar an, wenn sie den Server nicht mehr erreicht, statt still stehenzubleiben
  proof: core/web/ui/src/pages/review/ui/ReviewPage.svelte:68 heartbeat, :96 banner
- [x] Die Ursache von Alex' Ausfall ist mit dem Log gefunden und behoben, oder als Note mit dem Logauszug festgehalten
  proof: decision note 19.10. (Alex, Option 2): nicht reproduziert, ab diesem Build im Log
- [x] core/release/notes/<ID>.md; go test ./... -race gruen
  proof: core/release/notes/HEY6T7.md; go test -race ./... exit 0

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/web/log.go: Logger (Datei, Rotation), recover-Middleware, Fehlerstatus-Logging, Loop/Fetch/Push-Fehler + recover
- [x] internal/cli/web.go: Log oeffnen, Pfad beim Start nennen, Fetch-Timeout
- [x] Tests: Panic -> 500 + Stack im Log, Server laeuft weiter; >=400 und Fetch-Fehler im Log
- [x] UI: Herzschlag + Banner 'antwortet nicht mehr', dist neu bauen
- [x] release note, go test ./... -race

## Progress
- **2026-10-19 12:47 · Alexander Sacharov** — parent session: jaira-8b
- **2026-10-19 12:48 · Alexander Sacharov** — decision (worker, keine offenen Entscheidungen fuer Alex): Log nach <StateDir>/web.log (pro Arbeitsbaum unter ~/.jaira/state/..., ausserhalb des Repos - landet nie in einem Commit); beim Start rotiert ab 1 MB nach web.log.1, damit es nicht waechst. Panic: eigenes recover-Middleware (net/http faengt Handler-Panics zwar selbst, schreibt aber nur nach stderr und kappt die Verbindung ohne 500) plus recover in Loop/FetchNow, denn eine Panic dort beendet den Prozess. Seite: Herzschlag auf /api/sync alle 10 s; scheitert er, ein Banner 'jaira web antwortet nicht mehr seit HH:MM' oben, verschwindet beim naechsten Erfolg und laedt dann neu. Fetch-Kind bekommt einen Timeout (2 min), damit ein haengender git fetch nicht still stehenbleibt.
- **2026-10-19 12:48 · Alexander Sacharov** — Befund 19.10. nachmittags: Alex' 'jaira web' im anderen Repo ist inzwischen ein neu gestarteter Prozess - der vom Morgen lief nicht weiter oder wurde neu gestartet. Der neue antwortet auf /api/sync, /api/epics, /api/unassigned und / mit 200, der letzte Fetch ist aktuell. Ursache ohne Log nicht feststellbar.
- **2026-10-19 12:50 · Alexander Sacharov** — Fetch-Timeout verworfen (Note von vorhin): ein abgebrochenes 'jaira fetch' mitten im Ref-Schreiben ist riskanter als ein haengendes, und das DoD verlangt ihn nicht. Stattdessen faengt 'jaira web' jetzt SIGTERM und SIGHUP wie Ctrl-C und schreibt das Signal ins Log - ein geschlossenes Terminal ist ein Kandidat fuer Alex' Ausfall.
- **2026-10-19 12:54 · Alexander Sacharov** — Live geprueft (Scratch-Binary, dieses Repo, freier Port): Start druckt 'jaira web: log ~/.jaira/state/<arbeitsbaum>/web.log'; Log hat 'started on ... (pid ...)', 'GET /api/tickets/NOPE: 404 {...not found}', 'stopped by signal hangup'. DoD 4 bleibt offen: die Ursache braucht einen Logauszug von Alex' eigenem Lauf mit dem neuen Build.
- **2026-10-19 12:55 · Alexander Sacharov** — Regel von Alex (via jaira-8b, 19.10.): keine anderen Projektnamen und keine privaten Daten (Home-Pfade, PIDs, Hostnamen) in Code, Tests, Commits, Release Notes oder Notes. Neutrale Beispiele nehmen (temporaeres Repo, example.com). Commits enden mit 'AlSa', kein Co-Authored-By.
- **2026-10-19 12:55 · Alexander Sacharov** — Frage an jaira-8b (offen): zwei fruehere Notes und der Kontext enthalten private Daten, jaira note kann nichts aendern. 1 (empfohlen) Stellen per Hand neutralisieren, mit dem naechsten Code-Commit; 2 zusaetzlich Ref neu schreiben; 3 so lassen.
- **2026-10-19 12:56 · Alexander Sacharov** — critique (Runde 1, ganzer Diff): nichts. Bewusst stehen gelassen: (a) bei Ausfall zeigt die Seite Banner UND die 'Failed to fetch'-Zeile aus load() - das Banner erklaert die Zeile, ein Zusammenlegen waere mehr Code; (b) recovered() um Fetch/Push, obwohl fetchInChild kaum panict - eine Panic in der Loop-Goroutine beendet sonst den Prozess, DoD 1 verlangt jede Panic; (c) kein bestehendes Log-Muster im Repo (grep log.New/slog leer), also kein Duplikat; (d) SIGHUP wird unter Windows nie geliefert, schadet dort nicht.
- **2026-10-19 13:00 · Alexander Sacharov** — testing: go test -race -count=1 ./... RC=0. Headless Chrome gegen ein Scratch-Binary: Seite ohne Banner; Befehl beendet -> nach ~10 s Banner 'jaira web has not answered since <Zeit> - is the command still running? ...'; Befehl neu gestartet -> Banner weg. Bereinigung der Notes/Kontext (Option 2 von Alex via jaira-8b): Hand-Edit der Ticketdatei wurde in dieser Session vom Berechtigungs-Klassifikator abgelehnt - nicht ausgefuehrt, liegt bei jaira-8b/Alex.
- **2026-10-19 13:46 · Alexander Sacharov** — decision (Alex via jaira-8b, 19.10.): DoD 4 Option 2 - mit Note schliessen. Nicht reproduziert; ab diesem Build steht jede Ursache eines Ausfalls (Fehlerantwort, Fetch/Push-Fehler, Panic mit Stack, beendendes Signal) in der beim Start genannten Logdatei. Datenbereinigung: Alex gibt sie selbst im Worker-Tab frei; bis dahin weder Branch noch Ref pushen.
- **2026-10-19 13:48 · Alexander Sacharov** — review (frisch): Code deckt DoD 1-3,5; keine Defekte. Blockierend fuer Push: private Daten (fremder Projektname, PIDs, Home-Pfad) in der committeten Ticketdatei (d89bb20c, Kontext + Notes 12:48/12:54) - Bereinigung liegt laut 13:46 bei Alex. Klein: logKeep-Kommentar ueberverspricht (Rotation nur beim Start); nicht oeffnbares web.log macht jaira web jetzt fatal.
- **2026-10-19 14:23 · Alexander Sacharov** — explain: Der Code erfuellt die Definition of Done, es gibt keine Defekte. Aber die in d89bb20c committete Ticketdatei enthaelt private Daten (fremder Projektname, PIDs, Home-Pfad), dazu zwei kleine Unschaerfen: der logKeep-Kommentar verspricht mehr Rotation, als es gibt, und ein nicht oeffnbares web.log macht jaira web fatal.
Warum: Mit den privaten Daten darf weder Branch noch Ref gepusht werden, Alex' Regel vom 19.10. Die kleinen Punkte sind kein Risiko, aber der Kommentar luegt und ein nicht beschreibbares State-Verzeichnis verhindert jetzt den Start.
1) Zurueckgeben: Ticketdatei bereinigen, Kommentar korrigieren, Log-Fehler nur warnen statt abbrechen [give back]
2) Zurueckgeben: nur die Ticketdatei bereinigen, die zwei kleinen Punkte bleiben [give back]
3) So annehmen; Bereinigung gibt Alex selbst im Worker-Tab frei, kleine Punkte bleiben stehen [accept]
Empfehlung: 1 - die Bereinigung ist ohnehin vor dem Push noetig, und die beiden kleinen Punkte sind in derselben Runde billig zu beheben.
- **2026-10-19 14:27 · Alexander Sacharov** — decision: 1) Zurueckgeben: Ticketdatei bereinigen, Kommentar korrigieren, Log-Fehler nur warnen statt abbrechen
- **2026-10-19 14:28 · Alexander Sacharov** — Given back from signoff: 1) Zurueckgeben: Ticketdatei bereinigen, Kommentar korrigieren, Log-Fehler nur warnen statt abbrechen
- **2026-10-19 14:28 · Alexander Sacharov** — cockpit: given back to the old worker - tab alive, cache warm for 21 more min, context 132k
- **2026-10-19 14:28 · Alexander Sacharov** — Runde 2 (Rueckgabe aus signoff, Alex Option 1): Bereinigung der Ticketdatei freigegeben durch Alex' Entscheidung im Cockpit - Hand-Edit als Ausnahme von 'nur CLI', weil jaira note nichts aendern kann; die Branch-Commits werden vor jedem Push per filter-branch neu geschrieben. Dazu logKeep-Kommentar korrigieren und ein nicht oeffnbares web.log nur warnen statt abbrechen.
- **2026-10-19 14:30 · Alexander Sacharov** — Runde 2 umgesetzt: logKeep-Kommentar sagt jetzt, dass nur beim Start rotiert wird (core/web/log.go:14). Ein nicht oeffnbares web.log warnt auf stderr ('cannot open the log ...; running without one') und der Server laeuft ohne Log weiter (internal/cli/web.go:70) - live geprueft mit schreibgeschuetztem State-Verzeichnis: Warnung, /api/sync 200. Ticketdatei bereinigt (Kontext + zwei Notes + review-gaps), auch der aktuelle Stand auf dem Ticket-Ref ist sauber; die aeltere Ref-History enthaelt die Daten weiter (siehe Report).
- **2026-10-19 14:34 · Alexander Sacharov** — decision (Alex via jaira-8b, 19.10.): auch die History des Ticket-Refs wird bereinigt und nach upstream force-gepusht. Ein Agent fuehrt weder Umschreiben noch Force-Push aus (Klassifikator lehnt beides ab, Alex' Regel): Alex startet das vorbereitete Skript (Backup-Refs, filter-branch fuer Branch und Ticket-Ref, Trockenpruefung) selbst und pusht danach mit --force-with-lease. Skript in einem Wegwerf-Klon getestet: alle Versionen frei von den Strings, Branch ausserhalb der Ticketdatei unveraendert.
- **2026-10-19 14:35 · Alexander Sacharov** — critique Runde 2 (nur die drei Review-Punkte): alle beantwortet - Kommentar stimmt (log.go:16), Log-Fehler warnt (cli/web.go:72), Ticketdatei sauber. Nichts Neues.
- **2026-10-19 14:38 · Alexander Sacharov** — review Runde 2 (frisch): DoD erfuellt, keine Defekte. Offen nur: History-Bereinigung vor Push (Alex, 14:34); klein: ErrorLog=io.Discard, wenn das Log nicht aufgeht.
- **2026-10-19 20:43 · Alexander Sacharov** — explain: Der Review fand keinen Code-Defekt, nur zwei Kleinigkeiten: Laesst sich web.log nicht oeffnen, bekommt http.Server ErrorLog=io.Discard, und dessen eigene Meldungen (z.B. Accept-Fehler) sind ganz weg (internal/cli/web.go:71,99). Ausserdem stehen in den Commits d89bb20c und b2389dea in der Ticketdatei noch ein fremder Projektname, PIDs und ein Home-Pfad.
Warum: Ohne Log sollen wenigstens die Server-Meldungen auf stderr bleiben, wie vor der Aenderung. Die History darf mit fremden Daten nicht gepusht werden.
1) Zurueckgeben: ErrorLog auf einen stderr-Logger statt io.Discard setzen, dann erneut pruefen [give back]
2) Annehmen wie es ist; Alex bereinigt die History (Skript, force-with-lease) vor dem Push, ErrorLog bleibt als bekannte Kleinigkeit [accept]
Empfehlung: 2 - DoD 1-5 sind erfuellt, der io.Discard-Fall tritt nur bei unoeffnbarem Log auf und das Ticket soll nicht weiter wachsen.
- **2026-10-19 20:45 · Alexander Sacharov** — decision: 1) Zurueckgeben: ErrorLog auf einen stderr-Logger statt io.Discard setzen, dann erneut pruefen
- **2026-10-19 20:45 · Alexander Sacharov** — Given back from signoff: 1) Zurueckgeben: ErrorLog auf einen stderr-Logger statt io.Discard setzen, dann erneut pruefen
- **2026-10-19 20:45 · Alexander Sacharov** — cockpit: given back to a new worker - the old worker has no Herdr tab open any more
- **2026-10-19 20:47 · Alexander Sacharov** — Runde 3 (Rueckgabe aus signoff, ErrorLog): geht web.log nicht auf, schreibt der Logger nach stderr statt io.Discard (internal/cli/web.go:71) - http.Server-Meldungen, Fehlerantworten, Panics und das Signal bleiben sichtbar wie vor dem Log. Live geprueft mit schreibgeschuetztem State-Verzeichnis: Warnung, dann 'started on', '404 ...', 'stopped by signal terminated' auf stderr; /api/sync 200. go test -race ./... RC=0.
- **2026-10-19 20:48 · Alexander Sacharov** — critique Runde 3 (nur der ErrorLog-Befund aus signoff 20:43): beantwortet - Fallback-Logger schreibt nach cmd.ErrOrStderr() (internal/cli/web.go:71), io-Import entfaellt. Nichts Neues.
- **2026-10-19 20:49 · Alexander Sacharov** — testing Runde 3: Suite gruen, Proof-Tests PASS. DoD-Proofs 1-3 zeigten seit Runde 2 eine Zeile zu hoch (log.go:63 ist Kommentar, logged steht auf :64; web.go 135/151 -> 137/153; Banner :97 -> :96) - korrigiert. review-check enthielt einen Home-Pfad (aus Review Runde 2) - neutralisiert, bevor er in einem Commit blieb.
