---
id: 01M11YYR2P6AZ2K2CERSSPDWGH
title: Der Startbildschirm zeigt fertige Tickets pro Tag der letzten sieben Tage
status: done
ready: true
creator: BeMuCa
goal: "Auf der Projektauswahl schaltet s eine Balkengrafik ein: x die letzten sieben Tage mit heute hervorgehoben, y die Tickets, die an dem Tag ins Logbuch gingen"
context: |-
  Berk am 27.08.: das Logbuch protokolliert Aktivitaet; daraus eine Statistik unten auf der Projektauswahl im TUI, per s ein- und ausblendbar.

  Datenquelle: die Ordnernamen des Logbuchs, .jaira/logbook/<initialen>-<yyyymmdd>/, ueber alle registrierten Boards (core/project.Load). Zaehlen = Verzeichnisse auflisten, kein Ticket wird geparst. Gemessen wird 'ins Logbuch gebracht', nicht 'abgenommen': ein Ticket, das in done liegen bleibt, zaehlt nicht - das steht so in der Grafik oder ihrer Hilfe.

  Annahme, zu bestaetigen: die Zahl ist ueber alle Boards summiert, weil der Bildschirm ueber den Boards steht; pro Board waere ein zweiter Modus.

  Vorsicht: die Taste s ist auf dem Browse-Bildschirm (internal/tui/browse.go:177) schon belegt - pruefen, ob Projektauswahl und Browse dieselbe Tastenbelegung teilen. Die Projektauswahl ist internal/tui/home.go.
definition-of-done: "Auf der Projektauswahl blendet s eine Grafik mit sieben Balken ein und wieder aus; heute ist farblich markiert; die Hoehe entspricht der Zahl der Logbuch-Eintraege des Tages ueber alle registrierten Boards; ein Test mit gebauten Logbuch-Ordnern prueft die Zaehlung an der Tagesgrenze; die Grafik sagt, dass sie Logbuch-Eintraege zaehlt"
blocked-by: []
commits: []
created-at: 2026-08-27T15:55:56Z
updated-at: 2026-08-31T16:21:34Z
claimed-by: EE-3NX6GL3-3495447
claimed-at: 2026-08-27T15:58:19Z
updated-by: BeMuCa
assignee: BeMuCa
outcome-what: "s auf der Projektauswahl blendet die Grafik ein/aus: 7 Spalten, Zahl oben, Balken (skaliert auf den staerksten Tag, mindestens eine Zeile bei >0), Tag des Monats unten, heute in Akzentfarbe. Daten aus den Ordnernamen des Logbuchs ueber alle Boards (Store.LoggedPerDay), alter Ordnername sync/ zaehlt mit."
outcome-why: "Der datierte Ordner ist die Statistik - kein Ticket wird gelesen, die Projektauswahl wird nicht langsamer. 'Ins Logbuch gebracht' statt 'abgenommen' steht im Titel der Grafik, weil ein in done liegendes Ticket nicht zaehlt."
outcome-resolves: "DoD: s ein/aus (Test), heute markiert (Test prueft das Datum, Dump zeigt die Farbe), Hoehe = Eintraege des Tages ueber alle Boards (Test mit zwei Boards und Grenzfaellen: 7 Tage alt draussen, ohne Datum uebersprungen), Titel nennt Logbuch-Eintraege. go test ./... -race gruen, Binary gebaut."
review-summary: "Zwei Commits. b91ecb1: Store.LoggedPerDay (core/ticket/store.go) zaehlt .md-Dateien pro Logbuch-Ordner nach dem Datum im Ordnernamen, alter Ordnername inklusive, Tagesabstand gerundet; HomeEntry.Logged wird in describe() befuellt; Home bekommt stats bool, Taste s, renderStats mit 7x4-Zellen-Zeilen und Hinweis 's stats'. Zweiter Commit: ein unerreichbarer Index-Guard weniger."
review-gaps: "Entfernt: der Index-Guard 'if i < len(days)' in renderStats - Logged hat immer logbookDays Eintraege, LoggedPerDay wird genau damit aufgerufen. Gesucht und nicht gefunden: eine zweite Stelle, die Logbuch-Ordner auflistet (logbookFolders in store.go ist die einzige, LoggedPerDay ruft sie). Kosten: LoggedPerDay laeuft pro Board beim Aufbau der Projektauswahl, auch wenn die Grafik aus ist - ein ReadDir pro Logbuch-Ordner, einstellig; nicht verlagert, ein Lazy-Pfad waere mehr Code als der Aufwand. Stehen gelassen: gofmt -l listet weiterhin internal/cli/tickets.go. Test-Helfer 'logged' (home_test) und 'mk' (logbook_test) bauen dasselbe in zwei Paketen - Testcode ueber Paketgrenzen, nicht gefaltet."
review-check: |-
  1. Im jaira-Repo: go build -o ~/.local/bin/jaira ./cmd/jaira.
  2. In irgendeinem Board mit einem abgenommenen Ticket: jaira logbook <handle> - es antwortet 'Logged ...'. Ohne ein solches Ticket: einen Ordner von Hand anlegen, .jaira/logbook/xx-20260827/, und eine leere .md hineinlegen.
  3. jaira (ohne Argumente) startet die Projektauswahl. Unten in der Zeile steht 's stats'.
  4. s druecken. Unter der Projektliste erscheint 'logbook entries per day · last 7 days · all boards', darunter sieben Spalten: oben die Zahl, dann Balken, unten der Tag des Monats. Die rechte Spalte ist heute, farblich anders, und zeigt mindestens 1.
  5. s noch einmal: die Grafik ist weg.
  6. go test ./internal/tui/ -run TestHomeStats -v und go test ./core/ticket/ -run TestLoggedPerDay -v: beide PASS.
review-verdict: "Angenommen, dieselbe Einschraenkung: eine Sitzung, kein zweites Modell. Geprueft: go test ./... -race mit leerem Cache gruen; Render-Dump 100x30 zeigt die Spalten buendig (der erste Wurf war um eine Zelle verschoben, weil TrimRight die Zeilen ungleich breit machte - gefunden, bevor es committet war). Eine Annahme steht im Ticket und im Titel der Grafik: Summe ueber alle Boards; pro Board nicht gebaut. Nicht gemacht: Lazy-Berechnung erst bei s - benannt in review-gaps."
---

# Der Startbildschirm zeigt fertige Tickets pro Tag der letzten sieben Tage

## Definition of Done

- [x] Auf der Projektauswahl blendet s eine Grafik mit sieben Balken ein und wieder aus; heute ist farblich markiert; die Hoehe entspricht der Zahl der Logbuch-Eintraege des Tages ueber alle registrierten Boards; ein Test mit gebauten Logbuch-Ordnern prueft die Zaehlung an der Tagesgrenze; die Grafik sagt, dass sie Logbuch-Eintraege zaehlt
  proof: internal/tui/home.go renderStats + case s; core/ticket/store.go LoggedPerDay; TestHomeStatsCountLogbookEntriesPerDayAcrossBoards, TestLoggedPerDayReadsTheFolderNames (Tagesgrenze, alter Ordnername, ohne Datum); Render-Dump 100x30: 7 Balken, Zahlen darueber, Tage darunter, heute rechts

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-08-27 16:02 · BeMuCa** — Annahme umgesetzt: ueber alle registrierten Boards summiert (der Schirm steht ueber den Boards). Pro Board waere ein zweiter Modus - nicht gebaut.

Gefunden: zentrierte Zeilen muessen gleich breit sein, sonst verschiebt centre.Render die Spalten gegeneinander - der erste Wurf hat Leerzeichen am Zeilenende abgeschnitten und die Zahlenzeile stand eine Zelle neben den Balken. Jede Zeile ist jetzt exakt 7x4 Zellen.

Tagesabstand gerundet statt abgeschnitten: bei Sommerzeitwechsel liegen zwei Mitternachte 23 oder 25 Stunden auseinander.

Die Taste s auf dem Browse-Bildschirm kollidiert nicht: browse ist ein Unterbildschirm, seine Tasten werden vorher abgefangen.
