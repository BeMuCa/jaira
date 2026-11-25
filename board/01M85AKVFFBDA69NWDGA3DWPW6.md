---
id: 01M85AKVFFBDA69NWDGA3DWPW6
title: jaira market index schreibt index.json fuer ein Katalog-Verzeichnis
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Ein Befehl erzeugt aus einem Katalog-Verzeichnis die Datei index.json, die Client und Website lesen, sodass niemand sie von Hand pflegt."
context: |-
  Heute gibt es keine Index-Datei. 'jaira lanes market' fragt die GitHub-Contents-API nach dem Verzeichnis und holt dann jede Lane einzeln: N+1 Anfragen, und ohne Anmeldung erlaubt die API 60 Anfragen pro Stunde (core/market/market.go, List).
  Eine Website kann so auch nichts anzeigen, ohne selbst die API zu fragen.
  Loesung: 'jaira market index <verzeichnis>' schreibt index.json auf stdout. Das Katalog-Repo ruft es in CI nach jedem Merge auf.
  Je Eintrag: kind (lane oder board), id, name, description, path, sha256 des Inhalts, min-jaira. Bei einer Lane zusaetzlich agentic, model-tier, after, rejects-to. Bei einem Board die Lane-Reihenfolge und die Dateien mit je eigenem sha256.
  min-jaira: die kleinste jaira-Version, die alle Felder der Datei kennt. Das ersetzt den Tag-Pin von heute (apiBase/pinnedRef in core/market/market.go), der nur geht, solange der Katalog im Code-Repo liegt.
  Kopfzeile der Datei: format-Version und der Commit, aus dem sie erzeugt wurde.
  Ausgabe ist stabil sortiert, damit ein Diff nur zeigt, was sich geaendert hat.
  Vorgabe: docs/MARKET.md.
definition-of-done: "jaira market index <dir> gibt gueltiges JSON mit allen genannten Feldern aus, sortiert nach kind und id"
tags:
  - cli
blocked-by:
  - 5TED33
related: []
commits: []
created-at: 2026-11-23T23:07:18Z
updated-at: 2026-11-23T23:07:18Z
---

# jaira market index schreibt index.json fuer ein Katalog-Verzeichnis

## Definition of Done

- [ ] jaira market index <dir> gibt gueltiges JSON mit allen genannten Feldern aus, sortiert nach kind und id
- [ ] Zwei Laeufe ueber dasselbe Verzeichnis liefern dieselben Bytes
- [ ] Ein Verzeichnis, das jaira market lint ablehnt, erzeugt keinen Index, sondern Exit 1
- [ ] Test: sha256 im Index stimmt mit dem Dateiinhalt ueberein
- [ ] Release-Note core/release/notes/<TICKET>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

