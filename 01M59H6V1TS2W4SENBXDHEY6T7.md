---
id: 01M59H6V1TS2W4SENBXDHEY6T7
title: "jaira web schreibt mit, was passiert, und ein Absturz hinterlaesst seine Ursache"
status: optimize
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wenn die Seite von 'jaira web' ausfaellt, steht danach in einer Logdatei, was zuletzt passiert ist und warum."
context: |-
  Was heute fehlt: 'jaira web' loggt nichts.
  - core/web/*.go und internal/cli/web.go: kein log-Aufruf, kein recover(). Eine Panic in einem Handler oder in srv.Loop (Fetch alle 60 s, internal/cli/web.go:24) beendet den Prozess; der Stacktrace geht nur nach stderr und ist weg, wenn web nicht in einem sichtbaren Terminal laeuft.
  Anlass: Alex am 19.10.2026 - 'Web faellt nach ein paar Minuten aus'.
  Schon geprueft (Teamlead, 19.10.2026):
  - jaira 0.5.4, 'jaira web --no-open' in ~/projects/requirementsgenie 8 Minuten lang alle 20 s mit curl /api/epics abgefragt: immer 200, kein Absturz, stderr leer.
  - Alex' eigener 'jaira web' (PID 17706, cwd ~/projects/.worktrees/requirementsgenie-harness) lief um 09:32 seit 08:58 weiter - der PROZESS starb dort also nicht. Moeglich: die Seite im Browser verliert die Verbindung oder zeigt einen Fehler, waehrend der Server lebt.
  Was Alex genau sieht (Fehlermeldung, leere Seite, Prozess weg) steht als Note am Ticket, sobald bekannt.
definition-of-done: "Jede Anfrage mit Fehlerstatus (>=400), jeder fehlgeschlagene Fetch/Push und jede Panic wird mit Zeit und Ursache in eine Logdatei geschrieben; 'jaira web' nennt beim Start ihren Pfad"
tags:
  - web
blocked-by: []
related: []
commits: []
created-at: 2026-10-19T07:32:24Z
updated-at: 2026-10-19T12:56:40Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-79475
claimed-at: 2026-10-19T12:47:47Z
outcome-what: "jaira web schreibt Fehlerantworten (>=400), fehlgeschlagene Fetches/Pushes, Panics mit Stack und das beendende Signal mit Zeit nach ~/.jaira/state/<baum>/web.log und nennt den Pfad beim Start; eine Panic im Handler antwortet 500 statt die Verbindung zu kappen; die Seite zeigt nach spaetestens 10 s ein Banner, wenn der Befehl nicht mehr antwortet."
outcome-why: "Alex' 'jaira web' fiel aus, und es gab keine Spur: nichts wurde geloggt, stderr ging mit dem Terminal verloren, und die Seite stand still, ohne es zu sagen."
outcome-resolves: "Ein Ausfall von jaira web hinterlaesst jetzt seine Ursache in einer Datei, und die Seite meldet ihn sichtbar."
review-summary: none
review-gaps: "removed recorder.Unwrap (nothing uses http.ResponseController here) and the test-only syncBuffer (log.Logger serialises writes, tests are single-goroutine); left NewLogger (it fixes the timestamp flags for CLI and tests) and logf's nil check (Server.Log is optional, existing tests run without it)"
---

# jaira web schreibt mit, was passiert, und ein Absturz hinterlaesst seine Ursache

## Definition of Done

- [x] Jede Anfrage mit Fehlerstatus (>=400), jeder fehlgeschlagene Fetch/Push und jede Panic wird mit Zeit und Ursache in eine Logdatei geschrieben; 'jaira web' nennt beim Start ihren Pfad
  proof: core/web/log.go:63 logged + core/web/web.go:135,151 + internal/cli/web.go:82; TestErrorAnswersAreLogged, TestFailedFetchAndPushAreLogged
- [x] Eine Panic in einem Handler wird abgefangen (500 + Stacktrace im Log), der Server laeuft weiter; Test
  proof: core/web/log.go:63 logged; TestPanicAnswers500AndIsLogged
- [x] Die Seite zeigt sichtbar an, wenn sie den Server nicht mehr erreicht, statt still stehenzubleiben
  proof: core/web/ui/src/pages/review/ui/ReviewPage.svelte:68 heartbeat, :97 banner
- [ ] Die Ursache von Alex' Ausfall ist mit dem Log gefunden und behoben, oder als Note mit dem Logauszug festgehalten
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
- **2026-10-19 12:48 · Alexander Sacharov** — Befund 19.10. 14:48: Alex' heutiger 'jaira web' ist PID 92131 (gestartet 09:50, cwd requirementsgenie-harness, Port 58300) - nicht mehr PID 17706 von 08:58, der lief also zwischen 09:32 und 09:50 nicht weiter oder wurde neu gestartet. 92131 antwortet auf /api/sync, /api/epics, /api/unassigned und / mit 200, last_fetch 14:48. Ursache ohne Log nicht feststellbar.
- **2026-10-19 12:50 · Alexander Sacharov** — Fetch-Timeout verworfen (Note von vorhin): ein abgebrochenes 'jaira fetch' mitten im Ref-Schreiben ist riskanter als ein haengendes, und das DoD verlangt ihn nicht. Stattdessen faengt 'jaira web' jetzt SIGTERM und SIGHUP wie Ctrl-C und schreibt das Signal ins Log - ein geschlossenes Terminal ist ein Kandidat fuer Alex' Ausfall.
- **2026-10-19 12:54 · Alexander Sacharov** — Live geprueft (Scratch-Binary, dieses Repo, Port 58411): Start druckt 'jaira web: log /home/alex/.jaira/state/jaira-hey6t7-888112df/web.log'; Log hat 'started on ... (pid ...)', 'GET /api/tickets/NOPE: 404 {...not found}', 'stopped by signal hangup'. DoD 4 bleibt offen: die Ursache braucht einen Logauszug von Alex' eigenem Lauf mit dem neuen Build.
- **2026-10-19 12:55 · Alexander Sacharov** — Regel von Alex (via jaira-8b, 19.10.): keine anderen Projektnamen und keine privaten Daten (Home-Pfade, PIDs, Hostnamen) in Code, Tests, Commits, Release Notes oder Notes. Neutrale Beispiele nehmen (temporaeres Repo, example.com). Commits enden mit 'AlSa', kein Co-Authored-By.
- **2026-10-19 12:55 · Alexander Sacharov** — Frage an jaira-8b (offen): zwei fruehere Notes und der Kontext enthalten private Daten, jaira note kann nichts aendern. 1 (empfohlen) Stellen per Hand neutralisieren, mit dem naechsten Code-Commit; 2 zusaetzlich Ref neu schreiben; 3 so lassen.
- **2026-10-19 12:56 · Alexander Sacharov** — critique (Runde 1, ganzer Diff): nichts. Bewusst stehen gelassen: (a) bei Ausfall zeigt die Seite Banner UND die 'Failed to fetch'-Zeile aus load() - das Banner erklaert die Zeile, ein Zusammenlegen waere mehr Code; (b) recovered() um Fetch/Push, obwohl fetchInChild kaum panict - eine Panic in der Loop-Goroutine beendet sonst den Prozess, DoD 1 verlangt jede Panic; (c) kein bestehendes Log-Muster im Repo (grep log.New/slog leer), also kein Duplikat; (d) SIGHUP wird unter Windows nie geliefert, schadet dort nicht.
