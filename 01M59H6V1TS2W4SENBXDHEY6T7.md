---
id: 01M59H6V1TS2W4SENBXDHEY6T7
title: "jaira web schreibt mit, was passiert, und ein Absturz hinterlaesst seine Ursache"
status: in-progress
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
updated-at: 2026-10-19T14:29:51Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-27212
claimed-at: 2026-10-19T13:47:43Z
outcome-what: "jaira web schreibt Fehlerantworten (>=400), fehlgeschlagene Fetches/Pushes, Panics mit Stack und das beendende Signal mit Zeit nach ~/.jaira/state/<baum>/web.log und nennt den Pfad beim Start; eine Panic im Handler antwortet 500 statt die Verbindung zu kappen; die Seite zeigt nach spaetestens 10 s ein Banner, wenn der Befehl nicht mehr antwortet."
outcome-why: "Alex' 'jaira web' fiel aus, und es gab keine Spur: nichts wurde geloggt, stderr ging mit dem Terminal verloren, und die Seite stand still, ohne es zu sagen."
outcome-resolves: "Ein Ausfall von jaira web hinterlaesst jetzt seine Ursache in einer Datei, und die Seite meldet ihn sichtbar."
review-summary: "jaira web oeffnet beim Start <StateDir>/web.log (ausserhalb des Repos, ab 1 MB beim naechsten Start nach web.log.1 verschoben) und druckt 'jaira web: log <pfad>'. Eine Middleware um alle Routen (core/web/log.go logged) schreibt jede Antwort >=400 mit Methode, Pfad, Status und den ersten 512 Byte des Fehlerbodys ins Log, faengt Handler-Panics ab, loggt sie mit Stack und antwortet 500 statt die Verbindung zu kappen. Fetch und Push laufen durch recovered(): eine Panic wird zum Fehler mit Stack, wird geloggt und landet in PushError, die Loop-Goroutine laeuft weiter. http.Server.ErrorLog zeigt auf dasselbe Log; SIGTERM/SIGHUP beenden wie Ctrl-C und das Signal steht im Log. Die Seite fragt alle 10 s /api/sync ab; scheitert das, erscheint ein rotes Banner 'has not answered since HH:MM', das beim naechsten Erfolg verschwindet und alles neu laedt. dist ist neu gebaut, Release Note mit zwei Zeilen vorhanden."
review-gaps: "1 (blockiert Push, nicht Code): die in d89bb20c committete Ticketdatei enthaelt private Daten gegen Alex' Regel vom 19.10. - fremder Projektname im Kontext und in einer Note, PIDs, ein Home-Pfad. 2 (klein): der logKeep-Kommentar verspricht Rotation waehrend eines Laufs, die es nicht gibt. 3 (klein): ein nicht oeffnbares web.log macht jaira web fatal."
test-verdict: "pass: go test -race -count=1 ./... RC=0; DoD 1,2,3,5 verified in the tree and live (scratch binary: log path printed, 404 and SIGHUP logged; headless Chrome: banner appears ~10 s after the command is killed, disappears after restart); DoD 4 needs a log from Alex' own run"
question: ""
review-verdict: "Der Code erfuellt die Definition of Done; keine Defekte gefunden, nur zwei kleine Unschaerfen (Kommentar zu Rotation, Log-Fehler ist fatal). Nicht pushbar, solange die privaten Daten in der committeten Ticketdatei stehen - das entscheidet Alex (Note 13:46)."
review-round: 1
review-check: "1. cd in den Worktree und: go build -o /tmp/jw ./cmd/jaira   2. /tmp/jw web --no-open   3. Die Ausgabe hat eine zweite Zeile 'jaira web: log <pfad>/web.log' - Pfad merken   4. In zweitem Terminal: curl -s http://127.0.0.1:<port>/api/tickets/NOPE   5. tail <pfad>/web.log zeigt 'started on ...' und 'GET /api/tickets/NOPE: 404 {...not found}' mit Zeitstempel   6. Die Seite im Browser oeffnen (Adresse aus Schritt 3), dann im ersten Terminal Ctrl-C   7. Nach spaetestens 10 s steht oben ein rotes Banner 'jaira web has not answered since HH:MM'   8. Schritt 2 erneut starten (gleicher Port ist nicht garantiert - dann neuen Tab oeffnen); bei gleichem Port verschwindet das Banner binnen 10 s   9. tail web.log zeigt 'stopped by signal interrupt'   Panic-Pfad hat keinen Handweg: go test ./core/web -run 'TestPanic|TestFailedFetch' -v"
round: 2
---

# jaira web schreibt mit, was passiert, und ein Absturz hinterlaesst seine Ursache

## Definition of Done

- [x] Jede Anfrage mit Fehlerstatus (>=400), jeder fehlgeschlagene Fetch/Push und jede Panic wird mit Zeit und Ursache in eine Logdatei geschrieben; 'jaira web' nennt beim Start ihren Pfad
  proof: core/web/log.go:63 logged + core/web/web.go:135,151 + internal/cli/web.go:82; TestErrorAnswersAreLogged, TestFailedFetchAndPushAreLogged
- [x] Eine Panic in einem Handler wird abgefangen (500 + Stacktrace im Log), der Server laeuft weiter; Test
  proof: core/web/log.go:63 logged; TestPanicAnswers500AndIsLogged
- [x] Die Seite zeigt sichtbar an, wenn sie den Server nicht mehr erreicht, statt still stehenzubleiben
  proof: core/web/ui/src/pages/review/ui/ReviewPage.svelte:68 heartbeat, :97 banner
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
