---
id: 01M43QT4X8K72505R1HDMJ3PVV
title: "jaira logbook listet nur, was in den letzten vier Wochen vom Board ging"
status: done
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
commits:
  - 0d758023f589726bda2a4808c207d3fcafe822ae
  - 075f1632bab9c2500bf56e7054439700ecc3ba34
created-at: 2026-10-04T15:16:39Z
updated-at: 2026-10-04T18:19:34Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-66852
claimed-at: 2026-10-04T15:23:16Z
outcome-what: "Review bestanden, Felder gesetzt"
outcome-why: "Alex hat MJ3PVV aus human angenommen"
outcome-resolves: MJ3PVV wartet auf die Abnahme in signoff
review-summary: |-
  internal/cli/logbook.go:250 logbookDay dupliziert das Lesen des Ordnerdatums aus core/ticket/store.go:463-470 (LoggedPerDay) — Helfer einmal in core/ticket anlegen (z. B. ticket.LogbookFolderDay(name, loc)) und von LoggedPerDay und listLogbook aufrufen, sonst laufen zwei Parser für dasselbe Format auseinander.
  internal/cli/logbook.go:174-195 parseSince nimmt über time.ParseDuration auch Stunden/Minuten an (48h, 1h, 90m), die Grenze wird aber auf Mitternacht gerundet (logbook.go:211) — "--since 1h" listet den ganzen heutigen Tag, die "feinere" Grenze gibt es nicht. Nur Nw/Nd (und 0) annehmen, Fenster als Tageszahl führen und mit now.AddDate(0,0,-n) statt now.Add(-window) rechnen (auch DST-sicher); 48h aus Hilfetext (logbook.go:38,102), Fehlermeldung und core/release/NOTES.md streichen.
  critique Runde 2: none — Runde-1-Befund (doppelter Ordnerdatum-Parser, Stunden feiner als ein Tag) ist laut Note 15:28 behoben und im Diff bestätigt: ticket.LogbookFolderDay (core/ticket/store.go:452) wird von LoggedPerDay und listLogbook genutzt, parseSince nimmt nur w/d. Keine einfachere Form, kein neues Muster neben einem bestehenden, nichts Spekulatives.
  review: 'jaira logbook' ohne Argument filtert jetzt nach dem Datum im Ordnernamen (<initialen>-<yyyymmdd>, gelesen von ticket.LogbookFolderDay, den auch die Launcher-Grafik nutzt). Fenster: --since Nw/Nd, Standard 4w, Grenze ist lokale Mitternacht von heute minus N Tagen; --since 0 zeigt alles; undatierte Ordner werden immer gezeigt. Die Fußzeile nennt Anzahl und Startdatum und 'N older not listed — --since 0'. --json liefert logbook/count gefiltert plus hidden und since. --since mit --all oder ID, Stunden ('48h'), negative oder einheitslose Werte sind Exit 2. README (Logbook-Satz, zwei Command-Zeilen, neuer Abschnitt Roles), docs/COMMANDS.md und eine NOTES.md-Zeile beschreiben das.
review-gaps: "Nichts zu entfernen: der Datumshelfer ist geteilt (ticket.LogbookFolderDay), sonst keine Doppelung."
review-verdict: "Diff erfüllt die Definition of Done; keine Defekte gefunden. Einzige Unschärfe für Nutzer: '--since 0d' heißt 'alles'. Doku-Aussagen in README/COMMANDS gegen Code und 'roles install --help' geprüft und korrekt."
review-check: "1. cd /home/alex/projects/jaira  2. go run ./cmd/jaira logbook | tail -3 — letzte Zeilen: 'N in the logbook since <heute-28 Tage> (--since 4w)' und 'M older not listed — jaira logbook --since 0 lists everything'  3. go run ./cmd/jaira logbook --since 0 | tail -1 — '74 in the logbook' (alles, keine 'older'-Zeile)  4. go run ./cmd/jaira logbook --since 1d --json — Felder count, hidden, since (= gestriges Datum)  5. go run ./cmd/jaira logbook --since 48h — Fehler 'is not a window', Exit 2  6. go test ./internal/cli -run 'TestLogbookListsTheLastFourWeeksByDefault|TestLogbookSinceRefusesWhatItCannotMean' — ok"
test-verdict: "Bestanden. go test ./... grün; auf diesem Board: ohne Argument 47 gelistet, 27 ausgelassen; --since 0 alle 74; --since 4, 48h, -1w, mit --all und mit ID: Usage-Fehler (TestLogbookSinceRefusesWhatItCannotMean)."
question: Abnehmen? Review ist schon gelaufen und bestanden (Felder gesetzt); offen sind nur die drei nicht blockierenden Punkte in der letzten Notiz.
---

# jaira logbook listet nur, was in den letzten vier Wochen vom Board ging

## Definition of Done

- [x] 'jaira logbook' ohne Argument listet nur Einträge, deren Ordnerdatum höchstens vier Wochen zurückliegt
  proof: internal/cli/logbook.go:198; TestLogbookListsTheLastFourWeeksByDefault
- [x] Ein Schalter setzt den Zeitraum, und ein Wert zeigt wieder alles; '--all' behält seine Bedeutung
  proof: internal/cli/logbook.go:102 --since, 0 = alles (internal/cli/logbook.go:175); --since mit --all/ID verweigert (internal/cli/logbook.go:83)
- [x] Die Liste sagt am Ende, wie viele ältere Einträge ausgeblendet sind und wie man sie sieht
  proof: internal/cli/logbook.go:238
- [x] --json folgt demselben Zeitraum; Test in internal/cli deckt Standard, Schalter und 'alles' ab
  proof: internal/cli/logbook.go:219; TestLogbookListsTheLastFourWeeksByDefault, TestLogbookSinceRefusesWhatItCannotMean
- [x] Zeile in core/release/NOTES.md unter Unreleased
  proof: core/release/NOTES.md:17
- [x] README und docs/COMMANDS.md beschreiben das Zeitfenster; README hat einen Abschnitt zu den Rollen (Alex, 2026-10-04: 'und aktualisiere die README')
  proof: README.md:615 (Rollen), README.md:673, README.md:691, docs/COMMANDS.md:145

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-04 15:24 · Alexander Sacharov** — Entscheidungen beim Bau: Schalter heißt --since (Werte wie 4w, 10d, 48h; 0 = alles), weil --all schon 'ablegen' heißt. Datum kommt aus dem Ordnernamen <initialen>-<yyyymmdd>; ein Ordner ohne lesbares Datum wird immer gezeigt — was man nicht datieren kann, versteckt man nicht. Grenze auf ganze Tage (lokale Mitternacht), weil der Ordner nur einen Tag kennt.
- **2026-10-04 15:27 · Alexander Sacharov** — critique: internal/cli/logbook.go:250 logbookDay dupliziert das Lesen des Ordnerdatums aus core/ticket/store.go:463-470 (LoggedPerDay) — Helfer einmal in core/ticket anlegen (z. B. ticket.LogbookFolderDay(name, loc)) und von LoggedPerDay und listLogbook aufrufen, sonst laufen zwei Parser für dasselbe Format auseinander.
internal/cli/logbook.go:174-195 parseSince nimmt über time.ParseDuration auch Stunden/Minuten an (48h, 1h, 90m), die Grenze wird aber auf Mitternacht gerundet (logbook.go:211) — "--since 1h" listet den ganzen heutigen Tag, die "feinere" Grenze gibt es nicht. Nur Nw/Nd (und 0) annehmen, Fenster als Tageszahl führen und mit now.AddDate(0,0,-n) statt now.Add(-window) rechnen (auch DST-sicher); 48h aus Hilfetext (logbook.go:38,102), Fehlermeldung und core/release/NOTES.md streichen.
- **2026-10-04 15:28 · Alexander Sacharov** — Kritik Runde 1 behoben: Ordnerdatum liest jetzt ein gemeinsamer Helfer ticket.LogbookFolderDay (core/ticket/store.go), den auch LoggedPerDay nutzt. --since nimmt nur noch Wochen und Tage (4w, 10d, 0); Stunden hätten eine Grenze versprochen, die ein Ordner mit Tagesdatum nicht ziehen kann. Grenze per time.Date(..., Day()-days) statt Add, damit Sommerzeit nicht verschiebt.
- **2026-10-04 15:30 · Alexander Sacharov** — critique Runde 2 (unabhängig): nichts mehr zu sagen. Stehe auf Note 15:28 — Parser-Duplikat und Stunden-Fenster sind behoben. Nicht blockierend, für review: internal/cli/logbook.go:181 nimmt '--since 0d'/'0w' an und listet dann alles (wie '--since 0'), obwohl 0d als 'nur heute' gelesen werden kann; internal/cli/logbook.go:205 Kommentar sagt 'AddDate', der Code nutzt time.Date(..., Day()-days).
- **2026-10-04 15:30 · Alexander Sacharov** — review (unabhängig): bestanden. Nicht blockierend: '--since 0d' listet alles (internal/cli/logbook.go:181); veralteter Kommentar 'AddDate' (internal/cli/logbook.go:205); README-Abschnitt Roles (README.md:615) ist fachfremd, aber korrekt.
- **2026-10-04 15:31 · Alexander Sacharov** — Kritik Runde 2 und Review bestanden. Offen, nicht blockierend: (a) --since 0d und 0w zeigen alles wie --since 0 (internal/cli/logbook.go:181) — wer 'nur heute' meint, bekommt alles; (b) Kommentar logbook.go:205 nennt AddDate, der Code nutzt time.Date(..., Day()-days); (c) der Rollen-Abschnitt im README gehört thematisch nicht zu diesem Ticket, kam auf Alex' Bitte mit.
