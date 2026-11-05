---
id: 01M6NHMX3XGC9RBPA29NR6X6FH
title: "Plugins: Entscheidung und Standard stehen in docs/PLUGINS.md"
status: backlog
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
updated-at: 2026-11-05T09:47:52Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-21000
claimed-at: 2026-11-05T09:47:52Z
---

# Plugins: Entscheidung und Standard stehen in docs/PLUGINS.md

## Definition of Done

- [ ] docs/PLUGINS.md beschreibt Aufruf, Umgebung, Schreibweg, external-Feld, Geheimnisse, Exit-Codes und Manifest
- [ ] docs/PLUGINS.md beschreibt den Weg in den Katalog (Fork + PR an panda-bobr/jaira-releases) und warum kein Go-plugin
- [ ] Ein JSON-Schema fuer jaira-plugin.json und fuer einen plugins.json-Eintrag liegt im Repo

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

