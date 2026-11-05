---
id: 01M6NHMX3XGC9RBPA29NR6X6FH
title: "Plugins: Entscheidung und Standard stehen in docs/PLUGINS.md"
status: critique
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
updated-at: 2026-11-05T10:10:45Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-21000
claimed-at: 2026-11-05T09:47:52Z
outcome-what: "docs/PLUGINS.md: Plugin-Standard samt Entscheidung vom 05.11.2026, dazu Events (jaira-<id> event <name>) und Web (Seiten unter /plugins/<id>/, api CGI-artig); docs/schema/ mit Schemas fuer jaira-plugin.json (inkl. events, web) und einen plugins.json-Eintrag"
outcome-why: "Wer ein Plugin schreiben will, braucht einen Ort mit allen drei Kanaelen (Befehl, Event, Web), Schreibweg, external, Geheimnissen, Exit-Codes, Manifest und Katalogweg - bisher stand das nur in der Teamlead-Sitzung"
outcome-resolves: Alle vier DoD-Punkte; jaira selbst startet noch keine Plugins (8Z3VRF und Folgetickets im Epic plugins)
review-summary: |-
  docs/PLUGINS.md Events: a plugin that calls jaira set/note from its event handler fires the next event - an endless loop with two plugins; say that jaira passes JAIRA_IN_EVENT=1 and writes made with it set fire no events
  docs/PLUGINS.md Web: 'pages or an api' - but the manifest has no way to declare an api without pages; say the capability web is what lets jaira web forward /plugins/<id>/api/*
  docs/PLUGINS.md Web: a file under web/api/ collides with the forwarded path; say api/ is reserved and never served as a file
round: 4
review-gaps: |-
  note: jaira show --json und jaira list --json geben das Feld external nicht aus - ein Plugin, das wie empfohlen ueber die Lesebefehle liest, findet seinen Schluessel (jira:PROJ-123) nicht und muss doch das Markdown parsen; entweder die JSON-Ausgabe ergaenzt external im Dispatch-Ticket, oder der Abschnitt Writing nennt das
  note: core/merge/merge_test.go fuehrt external als YAML-Block (provider/key), PLUGINS.md als String '<id>:<key>' - der Merge behandelt beides opak, aber der Test zeigt eine Form, die der Standard nicht mehr vorsieht
  note: Exit-Code 5 heisst in PLUGINS.md 'not found - ticket, issue or key', in COMMANDS.md 'no such ticket, or an ambiguous id prefix' - fuer Plugins sinnvoll erweitert, aber nicht wortgleich mit dem behaupteten gleichen Vertrag
test-verdict: "pass: go build ./... und go test ./... gruen (RC=0); DoD 1-3 im Baum geprueft (Abschnitte in docs/PLUGINS.md, beide Schemas Draft-2020-12-gueltig, Doku-Beispiele validieren, Bad_ID abgelehnt); 'jaira set <id> external=jira:PROJ-1' auf Scratch-Board ausgefuehrt, Feld steht im Frontmatter"
review-verdict: "Alle drei DoD-Punkte erfuellt: PLUGINS.md deckt Aufruf, Umgebung, Schreibweg, external, Geheimnisse, Exit-Codes, Manifest, Katalogweg und die Go-plugin-Absage ab; beide Schemas liegen im Repo. Geprueft: jaira set external=jira:PROJ-123 und note --json funktionieren wie beschrieben, JAIRA_ROOT gibt es fuer Hooks, external wird beim Merge opak behandelt. Kein Blocker; der wichtigste Hinweis ist, dass external in der JSON-Lesesicht fehlt."
review-round: 2
review-check: "1. less docs/PLUGINS.md - Abschnitte Invocation, Environment, Writing, The external field, Secrets, Exit codes, The manifest, Releases, Into the catalogue sind da  2. uvx check-jsonschema --check-metaschema docs/schema/jaira-plugin.schema.json docs/schema/plugins-entry.schema.json - beide 'ok'  3. Manifest-Beispiel aus PLUGINS.md in /tmp/m.json kopieren, dann uvx check-jsonschema --schemafile docs/schema/jaira-plugin.schema.json /tmp/m.json - 'ok'  4. in einem Testboard: jaira set <id> external=jira:PROJ-123, dann grep external .jaira/tickets/*.md - zeigt external: \"jira:PROJ-123\""
---

# Plugins: Entscheidung und Standard stehen in docs/PLUGINS.md

## Definition of Done

- [x] docs/PLUGINS.md beschreibt Aufruf, Umgebung, Schreibweg, external-Feld, Geheimnisse, Exit-Codes und Manifest
  proof: docs/PLUGINS.md: sections Invocation/Environment, Writing, external field, Secrets, Exit codes, The manifest
- [x] docs/PLUGINS.md beschreibt den Weg in den Katalog (Fork + PR an panda-bobr/jaira-releases) und warum kein Go-plugin
  proof: docs/PLUGINS.md: 'Into the catalogue' (fork + PR to panda-bobr/jaira-releases) and 'Why plugins' (Not Go's plugin package)
- [x] Ein JSON-Schema fuer jaira-plugin.json und fuer einen plugins.json-Eintrag liegt im Repo
  proof: docs/schema/jaira-plugin.schema.json, docs/schema/plugins-entry.schema.json - Draft 2020-12, alle JSON-Beispiele der Doku validieren (jsonschema check_schema + validate)
- [x] docs/PLUGINS.md beschreibt Events (Feld events, Aufruf, stdin, Timeout, Fehler) und Web (web.pages, /plugins/<id>/, api-Weiterleitung, capability web) sowie dass external in show/list --json fehlt (8Z3VRF); beide Schemas nachgezogen
  proof: docs/PLUGINS.md: '## Events', '## Web' (Pages, Its own data), 8Z3VRF-Satz unter Writing; docs/schema/jaira-plugin.schema.json: events, web, capability web - Beispiele validieren; plugins-entry.schema.json unveraendert (Katalog nennt nur das Repo)

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
- **2026-11-05 10:07 · Alexander Sacharov** — Alex (05.11.2026, Teamlead-Sitzung): zurueck nach in-progress. Es fehlt der dritte Kanal - jaira benachrichtigt ein Plugin von sich aus. Abschnitt 'Events' in docs/PLUGINS.md + Feld 'events' (geschlossene Liste: create, move, claim, set) in docs/schema/jaira-plugin.schema.json. Vertrag: nach erfolgreichem Schreiben ruft jaira 'jaira-<id> event <name>' fuer jedes abonnierte Plugin auf; stdin = JSON wie 'jaira show --json' nach der Aenderung plus actor und from_lane; Timeout 5 s wie core/hook; Fehler des Plugins lassen den Befehl nicht scheitern; ein Plugin, das sein Ziel nicht erreicht, puffert selbst. core/hook bleibt fuer Nutzerskripte. Code dafuer gehoert zu 8Z3VRF. Ausserdem in PLUGINS.md vermerken: external fehlt heute in show/list --json, kommt mit 8Z3VRF.
- **2026-11-05 10:08 · Alexander Sacharov** — Alex (05.11.2026), zweiter Zusatz: ein Plugin kann jaira web um eigene Seiten erweitern (z.B. Jira-Integration anders dargestellt). Vertrag fuer docs/PLUGINS.md, Abschnitt 'Web': Manifest-Feld 'web': {pages: [{path, title}]}; die Dateien liegen im Release-Asset unter web/ und landen in ~/.jaira/plugins/<id>/web/; jaira web liefert sie unter /plugins/<id>/ aus und zeigt jede page als Tab/Link. Daten: die Seite liest jaira webs bestehende JSON-API; fuer eigene Daten leitet jaira web /plugins/<id>/api/<pfad> CGI-artig weiter an 'jaira-<id> web <METHOD> /<pfad>' (Body auf stdin, Antwort JSON auf stdout, Status per Exit-Code/Feld) - kein dauerhaft laufender Plugin-Prozess, kein eigener Port. Neue capability 'web'. Warum kein eigener Server des Plugins: Ports, Lebensdauer und Windows-Prozesse waeren jaira webs Problem. Code dafuer kommt in ein eigenes Ticket im Epic plugins.
- **2026-11-05 10:09 · Alexander Sacharov** — Runde 2, selbst entschieden (Vertrag aus den Notes von 10:07/10:08 sonst 1:1): (1) Event-stdin = show-JSON-Objekt mit zwei zusaetzlichen Top-Level-Keys actor und from_lane (from_lane leer ausser bei move). (2) Plugins werden nacheinander, in id-Reihenfolge, aufgerufen; stdout/stderr des Plugins verworfen wie bei core/hook. (3) Web-API: Query-String in JAIRA_QUERY, Body auf stdin, JSON auf stdout; HTTP-Status aus dem Exit-Code (0->200, 2->400, 3->403, 5->404, sonst 500) - kein Statusfeld, ein Weg reicht. (4) web/-Dateien kommen als eigenes Release-Asset jaira-<id>_web.tar.gz (die Binaries sind nackte Dateien, kein Archiv). (5) Events brauchen keine eigene capability; neue capability nur 'web'. (6) plugins-entry.schema.json bleibt unveraendert: der Katalog nennt nur das Repo, events/web stehen im Manifest.
- **2026-11-05 10:10 · Alexander Sacharov** — critique Runde 3 (nach Events/Web): Event-Schleife ueber eigene Schreibbefehle -> JAIRA_IN_EVENT=1, Schreiben damit loest keine Events aus (selbst entschieden, sonst kann ein Plugin sich selbst endlos triggern); api ohne pages nicht deklarierbar -> capability web schaltet die Weiterleitung frei; web/api/ als Datei reserviert.
