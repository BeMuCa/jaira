---
id: 01M82BVPANDHP8VE1XTVMHCEN8
title: "Herdr App zeigt links den Ticket und die Rolle jeder Agent-Tab (worker-8tt890, pr-watjdq) statt 'Claude Code'"
status: testing
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "In der linken Agent-Leiste von Herdr App steht bei jeder Tab, die spawn.sh oder das Cockpit startet, Rolle und Ticket, damit Alex sieht, welcher Claude woran arbeitet."
context: |-
  Was falsch ist (Alex 22.11.2026, Screenshot): links in Herdr App stehen zehn Zeilen 'Claude Code', ohne Ticket oder Rolle.
  Herausgefunden: Herdr App zeigt für eine Kind-Tab ihren hc-Namen (getestet: hc run --name worker-test01 erscheint so in der Leiste). Der Titel (z.B. 'Jaira-worker 8TT890') und das Label jaira.tab=8TT890/worker werden dort nicht gezeigt.
  Code: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh, 'tab create' - kennt das Label '<ID>/<rolle>' schon, ruft hc run ohne --name.
  hc --name: 1-32 Zeichen a-z 0-9 - _, eindeutig unter laufenden Terminals.
definition-of-done: "herdr-hc.sh 'tab create' gibt hc run --name '<rolle>-<id klein>' (aus dem Label <ID>/<rolle>, z.B. worker-8tt890, pr-watjdq, review-…), nur erlaubte Zeichen, höchstens 32 (Test)"
tags:
  - cockpit
blocked-by: []
related: []
commits: []
created-at: 2026-11-22T19:31:20Z
updated-at: 2026-11-22T19:35:27Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3820322
claimed-at: 2026-11-22T19:31:51Z
outcome-what: "herdr-hc.sh 'tab create' gibt hc run --name <rolle>-<id klein> aus dem Label <ID>/<rolle>, auf [a-z0-9_-] und 32 Zeichen gebracht; bei name_taken -2..-9, dann ohne Namen; Test TestHerdrHcNamesTheTabAfterItsLabel, shipped.sha256, Release-Note"
outcome-why: "Herdr App zeigt eine Kind-Tab nur mit ihrem hc-Namen; ohne Namen stand links zehnmal 'Claude Code'"
outcome-resolves: "Alex sieht in der Agent-Leiste, welche Rolle an welchem Ticket arbeitet"
review-summary: none
review-gaps: "removed the per-name regexp check in TestHerdrHcNamesTheTabAfterItsLabel: the exact comparison of every --name already pins them"
---

# Herdr App zeigt links den Ticket und die Rolle jeder Agent-Tab (worker-8tt890, pr-watjdq) statt 'Claude Code'

## Definition of Done

- [x] herdr-hc.sh 'tab create' gibt hc run --name '<rolle>-<id klein>' (aus dem Label <ID>/<rolle>, z.B. worker-8tt890, pr-watjdq, review-…), nur erlaubte Zeichen, höchstens 32 (Test)
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh tab create; TestHerdrHcNamesTheTabAfterItsLabel (worker/pr/review/odd characters/too long)
- [x] Ist der Name schon vergeben (alte Tab desselben Tickets lebt noch), startet die Tab trotzdem: mit Suffix -2, -3 … oder ohne Namen (Test)
  proof: TestHerdrHcNamesTheTabAfterItsLabel (taken once/taken twice/long and taken/all taken/other refusal)
- [x] Ohne Label bleibt alles wie heute; shipped.sha256 aktualisiert, Release-Note core/release/notes/<ID>.md
  proof: TestHerdrHcNamesTheTabAfterItsLabel/no_label; core/shipped/shipped.sha256; core/release/notes/MHCEN8.md
- [ ] seen by a person: in Herdr App eine Tab über spawn.sh starten (z.B. Cockpit-Worker) — links steht worker-<id> statt 'Claude Code'

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 19:31 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 19:32 · Alexander Sacharov** — decision: hc antwortet auf doppelten --name mit kind name_taken, exit 1 (probiert). Daher: Name aus Label <ID>/<rolle> -> <rolle>-<id klein>, sonst [a-z0-9_-] erzwungen, auf 32 gekürzt; bei name_taken nacheinander -2 bis -9, danach ohne --name. Andere Fehler von hc bleiben Fehler.
- **2026-11-22 19:34 · Alexander Sacharov** — tr arbeitet auf Bytes: ein ü wird zu zwei '-'. Bewusst so gelassen, Namen bleiben gültig; Label kommen ohnehin aus Ticket-ID und Lane-Namen. Cockpit startet Tabs nur über spawn.sh (kein eigenes 'tab create' in core/mod), daher reicht die eine Stelle. Vier Kategorien: Nebenläufigkeit — zwei Spawns gleichzeitig mit gleichem Namen: der zweite bekommt name_taken und nimmt -2; Fehlerpfade — andere hc-Fehler gehen nach stderr, exit 1 wie vorher; Eingabe — leer/lang/Sonderzeichen getestet; Rechte — nichts Neues von außen.
- **2026-11-22 19:35 · Alexander Sacharov** — critique: nichts. Wiederholungsschleife in tab create folgt dem Muster von 'tab close' im selben Skript; Name wird nur dort gebaut, wo das Label schon da ist (herdr-hc.sh), spawn.sh bleibt unverändert, weil der alte Herdr kein --name kennt.
