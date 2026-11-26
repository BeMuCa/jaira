---
id: 01M85AM09PYFEBF9Q4E6BTHBEX
title: jaira lanes market liest index.json aus dem Katalog-Repo und prueft die Pruefsumme
status: backlog
ready: true
creator: BeMuCa
assignee: Berk
goal: "jaira lanes market zeigt und holt Lanes aus dem Katalog-Repo mit einer Anfrage fuer die Liste, und eine Datei, deren Pruefsumme nicht zum Index passt, wird nicht uebernommen."
context: |-
  Heute liest core/market/market.go das Verzeichnis lanes/ im Repo BeMuCa/jaira ueber die GitHub-Contents-API (Konstanten repo und dir, Funktion apiBase), eine Anfrage je Lane dazu.
  Folgen: 60 Anfragen pro Stunde ohne Anmeldung; neue Lanes erst nach einem Release (Pin auf den Tag in pinnedRef); nichts prueft, ob die geholte Datei die ist, die gereviewt wurde.
  Loesung: der Client holt <basis>/index.json und danach nur die gewaehlte Datei <basis>/<path>, beides als statische Datei per HTTPS. Basis ist die Pages-Adresse des Katalog-Repos.
  - sha256 der geholten Datei muss zum Index passen, sonst Abbruch mit Exit 1 und beiden Werten.
  - Eintraege mit min-jaira groesser als die laufende Version werden gelistet, aber als 'braucht jaira X' markiert, und adopt lehnt sie ab mit Hinweis auf 'jaira self upgrade'.
  - JAIRA_MARKET_API bleibt die Umstellung fuer Tests und Spiegel, bedeutet jetzt aber die Basis-Adresse; gleiche Regel wie bisher (https, http nur auf loopback - selfupdate.OverrideOf).
  - maxBody bleibt die Grenze fuer jede Antwort.
  - Ausgabe und --json von 'jaira lanes market' bleiben, plus Feld sha256.
  Kein Clone, kein git im Client: nur zwei HTTP-GETs.
  Vorgabe: docs/MARKET.md. Braucht das Katalog-Repo mit veroeffentlichtem index.json.
definition-of-done: jaira lanes market listet aus index.json mit genau einer Anfrage (Test mit httptest-Server zaehlt)
tags:
  - cli
blocked-by:
  - TGYQGJ
related: []
commits: []
created-at: 2026-11-23T23:07:23Z
updated-at: 2026-11-26T10:51:01Z
updated-by: Berk
---

# jaira lanes market liest index.json aus dem Katalog-Repo und prueft die Pruefsumme

## Definition of Done

- [ ] jaira lanes market listet aus index.json mit genau einer Anfrage (Test mit httptest-Server zaehlt)
- [ ] adopt bricht ab, wenn der Inhalt nicht zum sha256 des Index passt (Test)
- [ ] Ein Eintrag mit zu hohem min-jaira wird markiert und von adopt abgelehnt (Test)
- [ ] JAIRA_MARKET_API zeigt auf eine Basis-Adresse; Test in internal/cli/market_test.go angepasst
- [ ] Release-Note core/release/notes/<TICKET>.md, mit '- ! ' fuer die geaenderte Bedeutung von JAIRA_MARKET_API

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

