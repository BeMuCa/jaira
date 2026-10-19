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
updated-at: 2026-10-19T12:50:40Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-79475
claimed-at: 2026-10-19T12:47:47Z
---

# jaira web schreibt mit, was passiert, und ein Absturz hinterlaesst seine Ursache

## Definition of Done

- [ ] Jede Anfrage mit Fehlerstatus (>=400), jeder fehlgeschlagene Fetch/Push und jede Panic wird mit Zeit und Ursache in eine Logdatei geschrieben; 'jaira web' nennt beim Start ihren Pfad
- [ ] Eine Panic in einem Handler wird abgefangen (500 + Stacktrace im Log), der Server laeuft weiter; Test
- [ ] Die Seite zeigt sichtbar an, wenn sie den Server nicht mehr erreicht, statt still stehenzubleiben
- [ ] Die Ursache von Alex' Ausfall ist mit dem Log gefunden und behoben, oder als Note mit dem Logauszug festgehalten
- [ ] core/release/notes/<ID>.md; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/web/log.go: Logger (Datei, Rotation), recover-Middleware, Fehlerstatus-Logging, Loop/Fetch/Push-Fehler + recover
- [x] internal/cli/web.go: Log oeffnen, Pfad beim Start nennen, Fetch-Timeout
- [ ] Tests: Panic -> 500 + Stack im Log, Server laeuft weiter; >=400 und Fetch-Fehler im Log
- [ ] UI: Herzschlag + Banner 'antwortet nicht mehr', dist neu bauen
- [ ] release note, go test ./... -race

## Progress
- **2026-10-19 12:47 · Alexander Sacharov** — parent session: jaira-8b
- **2026-10-19 12:48 · Alexander Sacharov** — decision (worker, keine offenen Entscheidungen fuer Alex): Log nach <StateDir>/web.log (pro Arbeitsbaum unter ~/.jaira/state/..., ausserhalb des Repos - landet nie in einem Commit); beim Start rotiert ab 1 MB nach web.log.1, damit es nicht waechst. Panic: eigenes recover-Middleware (net/http faengt Handler-Panics zwar selbst, schreibt aber nur nach stderr und kappt die Verbindung ohne 500) plus recover in Loop/FetchNow, denn eine Panic dort beendet den Prozess. Seite: Herzschlag auf /api/sync alle 10 s; scheitert er, ein Banner 'jaira web antwortet nicht mehr seit HH:MM' oben, verschwindet beim naechsten Erfolg und laedt dann neu. Fetch-Kind bekommt einen Timeout (2 min), damit ein haengender git fetch nicht still stehenbleibt.
- **2026-10-19 12:48 · Alexander Sacharov** — Befund 19.10. 14:48: Alex' heutiger 'jaira web' ist PID 92131 (gestartet 09:50, cwd requirementsgenie-harness, Port 58300) - nicht mehr PID 17706 von 08:58, der lief also zwischen 09:32 und 09:50 nicht weiter oder wurde neu gestartet. 92131 antwortet auf /api/sync, /api/epics, /api/unassigned und / mit 200, last_fetch 14:48. Ursache ohne Log nicht feststellbar.
