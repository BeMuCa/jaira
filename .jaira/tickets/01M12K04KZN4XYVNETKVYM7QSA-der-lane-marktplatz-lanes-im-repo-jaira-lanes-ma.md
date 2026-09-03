---
id: 01M12K04KZN4XYVNETKVYM7QSA
title: "Der Lane-Marktplatz: lanes/ im Repo, jaira lanes market, ein Absatz fuer Beitraege"
status: done
ready: true
creator: BeMuCa
goal: "Wer jaira auf GitHub findet, sieht die Lanes anderer, holt sich eine mit einem Befehl und weiss, wie er seine eigene beisteuert"
context: |-
  Berk am 27.08.: ein Marktplatz fuer Nutzer, die das jaira-Repo auf GitHub finden. Grundlage ist schon da: lanes/ im Repo haelt critique.md und optimize.md (ein Katalog, KEINE Builtins - die liegen unter core/lane/builtin/ und sind in die Binary kompiliert), und core/lane/shipped_test.go prueft in CI jede Datei darin. Ein PR, der eine Lane hinzufuegt, wird also schon jetzt geprueft.

  Was fehlt: (1) 'jaira lanes market' - listet lanes/ auf dem GitHub-Hauptzweig (Contents-API, https://api.github.com/repos/BeMuCa/jaira/contents/lanes; core/selfupdate spricht schon mit api.github.com fuer Releases, gleicher Trust-Root, gleiches Client-Muster); 'jaira lanes market adopt <id>' laedt die Datei und uebernimmt sie ueber lane.Adopt in den Katalog (parst vor dem Kopieren; Adopt heisst zustimmen, den Prompt auszufuehren - die Hilfe sagt das schon). (2) Ein Absatz im README, dort wo Lanes erklaert werden ('Lanes are pipeline stages'), und lanes/README.md: was eine Datei braucht, wie der PR aussieht, dass CI sie prueft. (3) docs/COMMANDS.md um die neuen Befehle ergaenzen.

  Kein Server, keine Registry, kein eigenes Format: das Verzeichnis IST der Marktplatz.
definition-of-done: "jaira lanes market listet die Lanes aus lanes/ des GitHub-Repos mit id, Name, Beschreibung (Test gegen einen httptest-Server); jaira lanes market adopt <id> legt die Datei im Katalog ab und lanes add findet sie danach; ohne Netz eine klare Fehlermeldung, kein Absturz; README-Absatz und lanes/README.md nennen den PR-Weg und die CI-Pruefung; docs/COMMANDS.md hat beide Befehle"
blocked-by: []
commits: []
created-at: 2026-08-27T21:46:13Z
updated-at: 2026-08-31T16:21:33Z
claimed-by: EE-3NX6GL3-4119427
claimed-at: 2026-08-27T21:46:32Z
updated-by: BeMuCa
assignee: BeMuCa
outcome-what: "Entry ohne Raw; market.Fetch(entry) holt die eine Datei; adopt ruft List (fuer die id-Aufloesung) und dann Fetch. Test fuer Fetch dazu."
outcome-why: "Listing und Adoption teilen keinen Zustand mehr; das Listing traegt nur, was die Tabelle braucht."
outcome-resolves: "Kritikpunkt umgesetzt; market- und cli-Tests gruen."
review-summary: "Zwei Commits. 0782f6e: Paket core/market (List ueber die GitHub-Contents-API, jede Datei geparst, README und Kaputtes ausgelassen; Host-Override nach der Regel der Release-Hosts), CLI 'lanes market' und 'lanes market adopt <id>' (Temp-Datei -> lane.Adopt), selfupdate.OverrideOf und lane.Parse exportiert, README/lanes-README/COMMANDS. Zweiter Commit: Entry ohne Raw, Fetch fuer adopt."
review-gaps: "Gesucht: eine zweite HTTP-Hilfsfunktion im Repo - selfupdate hat ihr eigenes fetch mit maxDownload 128 MiB fuer Archive; market.fetch ist 20 Zeilen mit 4 MiB. Nicht gefaltet: die Grenzen unterscheiden sich absichtlich, und ein gemeinsamer Fetcher haette einen Parameter, den jede Seite anders setzt. Entfernt in der Kritik-Runde: Entry.Raw. Toter Code: errNotFound in market wird definiert und in fetch zurueckgegeben, aber kein Aufrufer unterscheidet ihn - stehen gelassen als benannter Fehler statt Zahl, weil List die Meldung durchreicht. Stehen gelassen: gofmt-Eintrag internal/cli/tickets.go. Kosten: List laedt jede Datei, um Beschreibung zu zeigen - die Contents-API kennt nur Dateinamen; bei zwei Lanes zwei Requests, ein Index waere Infrastruktur fuer ein Problem, das es noch nicht gibt."
review-check: |-
  1. jaira lanes market - Tabelle mit ID, NAME, DESCRIPTION; heute zwei Zeilen: critique und optimize, darunter der Hinweis auf 'market adopt'.
  2. JAIRA_LANES_DIR=$(mktemp -d) jaira lanes market adopt optimize - Ausgabe 'adopted .../optimize.md' und der Hinweis auf 'lanes show' und 'lanes add'; ls des Verzeichnisses zeigt optimize.md.
  3. Dieselbe Zeile noch einmal - Ablehnung mit 'already exists ... --force'.
  4. jaira lanes market adopt nope - Fehler nennt die vorhandenen ids, Exit 2.
  5. JAIRA_MARKET_API=http://127.0.0.1:1/x jaira lanes market - eine Fehlerzeile, Exit 1, kein Absturz.
  6. README: Abschnitt 'Lanes are pipeline stages', fetter Satz 'Lanes other people made' mit dem PR-Hinweis; lanes/README.md hat 'Adding yours'; docs/COMMANDS.md hat zwei market-Zeilen.
  7. go test ./core/market/ ./internal/cli/ -run 'Market' -v - alle PASS.
review-verdict: "Angenommen; dieselbe Sitzung, kein zweites Modell. Geprueft: go test ./... -race nach dem ersten Commit gruen, market/cli nach dem zweiten; live gegen GitHub gelistet und adoptiert (Ausgabe im Protokoll). Sicherheitsseite unveraendert zur geteilten Lane: Adopt parst vor dem Kopieren, nichts landet auf einem Board ohne 'lanes add', die Hilfe sagt 'lies es vorher'. Netzabhaengig, sagt es, schreibt bei Fehler nichts."
---

# Der Lane-Marktplatz: lanes/ im Repo, jaira lanes market, ein Absatz fuer Beitraege

## Definition of Done

- [x] jaira lanes market listet die Lanes aus lanes/ des GitHub-Repos mit id, Name, Beschreibung (Test gegen einen httptest-Server); jaira lanes market adopt <id> legt die Datei im Katalog ab und lanes add findet sie danach; ohne Netz eine klare Fehlermeldung, kein Absturz; README-Absatz und lanes/README.md nennen den PR-Weg und die CI-Pruefung; docs/COMMANDS.md hat beide Befehle
  proof: core/market (3 Tests gegen httptest), internal/cli/market.go (4 Tests: Listing, adopt+add, unbekannte id, ohne Netz); live: 'jaira lanes market' gegen GitHub listet critique und optimize, 'market adopt optimize' schreibt in den Katalog; README 'Lanes other people made', lanes/README.md 'Adding yours', COMMANDS.md zwei Zeilen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-08-27 21:51 · BeMuCa** — critique: List soll nur listen (Metadaten braucht es fuer die Tabelle - dafuer muss die Datei gelesen werden, weil die Contents-API nur Dateinamen kennt; also bleibt ein Download pro Lane fuer die Tabelle unvermeidlich, solange kein Index existiert). Damit faellt Finding 1 halb: adopt braucht sowieso genau die eine Datei, List braucht alle fuer die Beschreibung. Umsetzen: Entry ohne Raw, adopt holt die gewaehlte Datei mit Fetch(url) erneut - ein Request mehr, aber List bleibt reines Lesen von Metadaten und adopt haengt nicht am Zustand der Liste.
