---
id: 01M43QT4X8K72505R1HDMJ3PVV
title: "jaira logbook listet nur, was in den letzten vier Wochen vom Board ging"
status: critique
ready: true
creator: Alexander Sacharov
goal: "'jaira logbook' ohne Argument zeigt nur Tickets, die im gewählten Zeitraum ins Logbook gingen — Standard vier Wochen —, damit man sieht, was zuletzt vom Board ging, statt der ganzen Geschichte."
context: |-
  Was falsch ist: 'jaira logbook' ohne Argument listet das ganze Logbook. Auf diesem Board sind das am 2026-10-04 schon 74 Zeilen, die ältesten vom 2026-09-11. Was zuletzt vom Board ging, sucht man darin von Hand.
  Auslöser: Alex am 2026-10-04: im Logbook sollen nur Tickets der letzten N Zeit erscheinen, Standard vier Wochen.
  Bekannt: Jeder Ordner heißt <initialen>-<yyyymmdd> (.jaira/logbook/as-20260911/) — das Datum des Ablegens steht also schon im Pfad, kein neues Feld nötig.
  Bekannt: '--all' ist bei 'jaira logbook' schon belegt und heißt 'alles aus der Terminal-Lane ablegen' (internal/cli/archive.go). Der neue Schalter braucht einen anderen Namen, z. B. '--since 8w' und '--since 0' für alles.
  Gilt auch für '--json'. Die Aktivitätsgrafik im Launcher (internal/tui/home.go, logbookDays = 7) ist ein anderes Ding und bleibt.
  Kein Löschen: ältere Einträge bleiben auf der Platte und in git, sie werden nur nicht gelistet.
definition-of-done: "'jaira logbook' ohne Argument listet nur Einträge, deren Ordnerdatum höchstens vier Wochen zurückliegt"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-04T15:16:39Z
updated-at: 2026-10-04T15:26:06Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-66852
claimed-at: 2026-10-04T15:23:16Z
outcome-what: "jaira logbook listet ohne Argument die letzten vier Wochen; --since setzt das Fenster, 0 = alles; Fußzeile nennt die ausgelassenen; --json mit hidden und since"
outcome-why: "Das ganze Logbook (74 Einträge) verdeckte, was zuletzt vom Board ging"
outcome-resolves: Logbook-Liste ohne Zeitfenster
---

# jaira logbook listet nur, was in den letzten vier Wochen vom Board ging

## Definition of Done

- [x] 'jaira logbook' ohne Argument listet nur Einträge, deren Ordnerdatum höchstens vier Wochen zurückliegt
  proof: internal/cli/logbook.go:205; TestLogbookListsTheLastFourWeeksByDefault
- [x] Ein Schalter setzt den Zeitraum, und ein Wert zeigt wieder alles; '--all' behält seine Bedeutung
  proof: internal/cli/logbook.go:102 --since, 0 = alles (internal/cli/logbook.go:174); --all unverändert, --since mit --all/ID verweigert (internal/cli/logbook.go:83)
- [x] Die Liste sagt am Ende, wie viele ältere Einträge ausgeblendet sind und wie man sie sieht
  proof: internal/cli/logbook.go:244
- [x] --json folgt demselben Zeitraum; Test in internal/cli deckt Standard, Schalter und 'alles' ab
  proof: internal/cli/logbook.go:225; TestLogbookListsTheLastFourWeeksByDefault, TestLogbookSinceRefusesWhatItCannotMean in internal/cli/logbook_test.go
- [x] Zeile in core/release/NOTES.md unter Unreleased
  proof: core/release/NOTES.md:17

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-04 15:24 · Alexander Sacharov** — Entscheidungen beim Bau: Schalter heißt --since (Werte wie 4w, 10d, 48h; 0 = alles), weil --all schon 'ablegen' heißt. Datum kommt aus dem Ordnernamen <initialen>-<yyyymmdd>; ein Ordner ohne lesbares Datum wird immer gezeigt — was man nicht datieren kann, versteckt man nicht. Grenze auf ganze Tage (lokale Mitternacht), weil der Ordner nur einen Tag kennt.
