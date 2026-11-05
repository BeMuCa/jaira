---
id: 01M6NHMX3XGC9RBPA29NR6X6FH
title: "Plugins: Entscheidung und Standard stehen in docs/PLUGINS.md"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer ein jaira-Plugin schreiben will, findet in einem Dokument, wie es aufgerufen wird, was es darf und wie es in den Katalog kommt."
context: |-
  Entschieden von Alex am 05.11.2026 (Teamlead-Sitzung): jaira bekommt Plugins nach dem Obsidian-Muster.
  - Ein Plugin ist ein eigenes Programm jaira-<id> in einem eigenen Repo seines Autors, mit eigenen Releases. jaira speichert keinen Plugin-Code.
  - Abgelehnt: Go-Paket plugin (kein Windows, gleiche Toolchain, cgo - bricht das statische Binary); Jira fest eingebaut (jede Firma nutzt Jira anders, jaira waechst).
  - Katalog = Whitelist: plugins.json in https://github.com/panda-bobr/jaira-releases. Ein Autor forkt das Repo und oeffnet einen PR mit seinem Eintrag {id, name, author, description, repo}; CI prueft, Alex/Berk mergen. Updates brauchen keinen PR.
  - Vorlage: https://github.com/panda-bobr/jaira-plugin-template (GitHub-Template, "Use this template", kein Fork).
  - Spaeter: eine Website mit DB und Registrierung ersetzt den PR und spiegelt per Bot in plugins.json (wie Obsidian heute). Das Format von plugins.json bleibt dabei gleich.
  - Beide Repos sind am 05.11.2026 leer angelegt; jaira zieht mit Berk in die Organisation panda-bobr um. Umbenennen: entschieden, jaira bleibt.
  Was dieses Ticket schreibt: docs/PLUGINS.md (Standard) und die Entscheidung oben als Abschnitt darin.
  Inhalt des Standards (Vorschlag aus der Sitzung):
  - Aufruf: 'jaira <id> ...' startet ~/.jaira/plugins/<id>/jaira-<id>; Umgebung JAIRA_ROOT, JAIRA_BOARD, JAIRA_VERSION.
  - Schreiben nur ueber die CLI (jaira create/set/note --json), nie direkt in .jaira/tickets/.
  - Verweise ins Fremdsystem im Feld external als '<id>:<schluessel>' (z.B. jira:PROJ-123). Das Feld ist dafuer reserviert: core/ticket/schema.go:77.
  - Geheimnisse nur in ~/.jaira/secrets.json (siehe 3Q92JY).
  - Exit-Codes wie jaira (0 ok, 3 verweigert); Unterbefehl 'manifest --json' fuer jaira web.
  - Manifest jaira-plugin.json: id, name, version, minJairaVersion, commands, capabilities. Release-Tag = v<version>, Assets jaira-<id>_<os>_<arch> + checksums.txt (goreleaser).
definition-of-done: "docs/PLUGINS.md beschreibt Aufruf, Umgebung, Schreibweg, external-Feld, Geheimnisse, Exit-Codes und Manifest"
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-11-05T09:46:40Z
updated-at: 2026-11-05T09:50:54Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-21000
claimed-at: 2026-11-05T09:47:52Z
outcome-what: docs/PLUGINS.md schreibt den Plugin-Standard und die Entscheidung vom 05.11.2026 fest; docs/schema/ enthaelt JSON-Schemas fuer jaira-plugin.json und einen plugins.json-Eintrag
outcome-why: "Wer ein Plugin schreiben will, braucht einen Ort mit Aufruf, Umgebung, Schreibweg, external, Geheimnissen, Exit-Codes, Manifest und Katalogweg - bisher stand das nur in der Teamlead-Sitzung"
outcome-resolves: Alle drei DoD-Punkte; jaira selbst startet noch keine Plugins (spaeteres Ticket)
review-summary: |-
  docs/schema/*.schema.json: $id points at BeMuCa/jaira, the repo moves to panda-bobr - the URL goes stale and nothing resolves it; drop $id
  docs/PLUGINS.md Secrets: 'Until that lands, a plugin reads its token from where 3Q92JY puts it' is circular; say the file does not exist until 3Q92JY ships and its layout is defined there
  docs/PLUGINS.md manifest: 'and a later jaira plugins' promises a command nobody decided; drop it
round: 2
---

# Plugins: Entscheidung und Standard stehen in docs/PLUGINS.md

## Definition of Done

- [x] docs/PLUGINS.md beschreibt Aufruf, Umgebung, Schreibweg, external-Feld, Geheimnisse, Exit-Codes und Manifest
  proof: docs/PLUGINS.md: sections Invocation/Environment, Writing, external field, Secrets, Exit codes, The manifest
- [x] docs/PLUGINS.md beschreibt den Weg in den Katalog (Fork + PR an panda-bobr/jaira-releases) und warum kein Go-plugin
  proof: docs/PLUGINS.md: 'Into the catalogue' (fork + PR to panda-bobr/jaira-releases) and 'Why plugins' (Not Go's plugin package)
- [x] Ein JSON-Schema fuer jaira-plugin.json und fuer einen plugins.json-Eintrag liegt im Repo
  proof: docs/schema/jaira-plugin.schema.json, docs/schema/plugins-entry.schema.json - Draft 2020-12, doc examples validate (jsonschema check_schema + validate)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-05 09:47 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/1770.sock
- **2026-11-05 09:48 · Alexander Sacharov** — Offene Entscheidungen: keine fuer die Person. Selbst entschieden: (1) docs/PLUGINS.md auf Englisch (Projektsprache). (2) Schemas unter docs/schema/jaira-plugin.schema.json und docs/schema/plugins-entry.schema.json, JSON Schema 2020-12. (3) capabilities = geschlossene Liste read/write/network/secrets; nur angezeigt, nicht erzwungen (kein Sandbox). (4) commands = Liste {name, summary}. (5) id: lowercase-kebab, eingebaute jaira-Befehle gewinnen bei Namensgleichheit. (6) Das Dokument sagt offen: Standard, jaira startet Plugins noch nicht - Umsetzung ist ein spaeteres Ticket. (7) external als einzelner String '<id>:<key>' - 'jaira set <id> external=jira:PROJ-1' funktioniert heute schon (auf Scratch-Board geprueft). (8) JAIRA_ROOT ist schon die Hook-Variable (core/hook/hook.go:64) - gleiche Bedeutung. (9) Keine Release-Note: reine Doku, fuer Nutzer des Binaries nicht sichtbar.
- **2026-11-05 09:50 · Alexander Sacharov** — Geprueft: beide Schemas sind gueltiges Draft 2020-12 (python jsonschema check_schema), die zwei JSON-Beispiele in docs/PLUGINS.md validieren dagegen, eine id 'Bad_ID' wird abgelehnt. Kein Test im Repo dafuer - reine Doku-Dateien, nichts im Binary liest sie.
- **2026-11-05 09:50 · Alexander Sacharov** — critique: $id in beiden Schemas raus (Repo zieht nach panda-bobr um); Secrets-Satz zu 3Q92JY ist zirkulaer - neu formulieren; 'a later jaira plugins' streichen. docs/AGENTS.md:8 ('no plugin API') bleibt stehen: spricht ueber Agenten, und jaira startet Plugins noch nicht.
