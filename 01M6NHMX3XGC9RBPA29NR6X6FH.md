---
id: 01M6NHMX3XGC9RBPA29NR6X6FH
title: "Plugins: Entscheidung und Standard stehen in docs/PLUGINS.md"
status: review
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
updated-at: 2026-11-05T09:54:19Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-21000
claimed-at: 2026-11-05T09:47:52Z
outcome-what: docs/PLUGINS.md schreibt den Plugin-Standard und die Entscheidung vom 05.11.2026 fest; docs/schema/ enthaelt JSON-Schemas fuer jaira-plugin.json und einen plugins.json-Eintrag
outcome-why: "Wer ein Plugin schreiben will, braucht einen Ort mit Aufruf, Umgebung, Schreibweg, external, Geheimnissen, Exit-Codes, Manifest und Katalogweg - bisher stand das nur in der Teamlead-Sitzung"
outcome-resolves: Alle drei DoD-Punkte; jaira selbst startet noch keine Plugins (spaeteres Ticket)
review-summary: "Zwei Commits, nur Doku: docs/PLUGINS.md legt den Plugin-Standard fest (Aufruf jaira <id> -> ~/.jaira/plugins/<id>/jaira-<id>, Umgebung JAIRA_ROOT/JAIRA_BOARD/JAIRA_VERSION, Schreiben nur ueber die CLI, external als '<id>:<key>', Geheimnisse in ~/.jaira/secrets.json (wartet auf 3Q92JY), Exit-Codes wie jaira, Manifest jaira-plugin.json + 'manifest --json', Releases per goreleaser, Weg in den Katalog per Fork+PR) und haelt die Entscheidung vom 05.11.2026 samt Absagen (Go-plugin, Jira eingebaut) fest. docs/schema/ bringt JSON-Schemas (Draft 2020-12) fuer das Manifest und einen plugins.json-Eintrag. jaira selbst aendert sich nicht."
round: 2
review-gaps: |-
  note: jaira show --json und jaira list --json geben das Feld external nicht aus - ein Plugin, das wie empfohlen ueber die Lesebefehle liest, findet seinen Schluessel (jira:PROJ-123) nicht und muss doch das Markdown parsen; entweder die JSON-Ausgabe ergaenzt external im Dispatch-Ticket, oder der Abschnitt Writing nennt das
  note: core/merge/merge_test.go fuehrt external als YAML-Block (provider/key), PLUGINS.md als String '<id>:<key>' - der Merge behandelt beides opak, aber der Test zeigt eine Form, die der Standard nicht mehr vorsieht
  note: Exit-Code 5 heisst in PLUGINS.md 'not found - ticket, issue or key', in COMMANDS.md 'no such ticket, or an ambiguous id prefix' - fuer Plugins sinnvoll erweitert, aber nicht wortgleich mit dem behaupteten gleichen Vertrag
test-verdict: "pass: go build ./... und go test ./... gruen (RC=0); DoD 1-3 im Baum geprueft (Abschnitte in docs/PLUGINS.md, beide Schemas Draft-2020-12-gueltig, Doku-Beispiele validieren, Bad_ID abgelehnt); 'jaira set <id> external=jira:PROJ-1' auf Scratch-Board ausgefuehrt, Feld steht im Frontmatter"
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
- **2026-11-05 09:51 · Alexander Sacharov** — critique Runde 2: nichts mehr. Bewusst stehen gelassen: Exit-Code-Tabelle uebernimmt auch 4 (dependencies) - gleiche Tabelle wie jaira ist der Punkt; JAIRA_BOARD-Herkunft (--board vor der id) legt das Dispatch-Ticket genau fest.
