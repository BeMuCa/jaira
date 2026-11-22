---
id: 01M82BVPANDHP8VE1XTVMHCEN8
title: "Herdr App zeigt links den Ticket und die Rolle jeder Agent-Tab (worker-8tt890, pr-watjdq) statt 'Claude Code'"
status: in-progress
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
updated-at: 2026-11-22T19:37:52Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3820322
claimed-at: 2026-11-22T19:31:51Z
outcome-what: "herdr-hc.sh 'tab create' gibt hc run --name <rolle>-<id klein> aus dem Label <ID>/<rolle>, auf [a-z0-9_-] und 32 Zeichen gebracht; bei name_taken -2..-9, dann ohne Namen; Test TestHerdrHcNamesTheTabAfterItsLabel, shipped.sha256, Release-Note"
outcome-why: "Herdr App zeigt eine Kind-Tab nur mit ihrem hc-Namen; ohne Namen stand links zehnmal 'Claude Code'"
outcome-resolves: "Alex sieht in der Agent-Leiste, welche Rolle an welchem Ticket arbeitet"
review-summary: none
review-gaps: "removed the per-name regexp check in TestHerdrHcNamesTheTabAfterItsLabel: the exact comparison of every --name already pins them"
test-verdict: "fail: with real hc a taken name opens no tab — hc writes its name_taken error to stderr, herdr-hc.sh looks for it in stdout"
hunt-verdict: "fail: bughunt done at 56f111c, 1 change (name_taken read from stdout, hc writes it to stderr — same as the testing finding)"
round: 2
back-from: testing
back-by: lane
back-count: 1
---

# Herdr App zeigt links den Ticket und die Rolle jeder Agent-Tab (worker-8tt890, pr-watjdq) statt 'Claude Code'

## Definition of Done

- [x] herdr-hc.sh 'tab create' gibt hc run --name '<rolle>-<id klein>' (aus dem Label <ID>/<rolle>, z.B. worker-8tt890, pr-watjdq, review-…), nur erlaubte Zeichen, höchstens 32 (Test)
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh tab create; TestHerdrHcNamesTheTabAfterItsLabel (worker/pr/review/odd characters/too long)
- [x] Ist der Name schon vergeben (alte Tab desselben Tickets lebt noch), startet die Tab trotzdem: mit Suffix -2, -3 … oder ohne Namen (Test)
  proof: herdr-hc.sh reads hc's stderr for name_taken; TestHerdrHcNamesTheTabAfterItsLabel (taken once/twice/all taken, stub refuses on stderr like hc); real hc: second ZZ99QQ/worker opened as worker-zz99qq-2
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
- **2026-11-22 19:36 · Alexander Sacharov** — testing: class: herdr-hc.sh reads hc's refusal from the wrong stream — hc run --json prints errors to stderr (checked: 2>/dev/null leaves stdout empty, rc 1), the script greps $out (stdout) for name_taken; the test stub echoes its error to stdout and so passes. example: two 'tab create --label ZZ99QQ/worker' against real hc: second prints name_taken, returns nothing, no tab. fix: capture hc's stderr separately (temp file) and grep that; make the stub write errors to stderr like hc does, so the test covers the real stream.
- **2026-11-22 19:36 · Alexander Sacharov** — bughunt: change: blocker: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:66-68 - class: hc's refusal is looked for on stdout, but hc (0.1.20) writes every --json error to stderr and leaves stdout empty, so the name_taken check never matches; the -2..-9 suffixes and the unnamed fallback are never tried, and a taken name means no tab at all (before this change the tab was always created). The fake hc in core/role/role_test.go TestHerdrHcNamesTheTabAfterItsLabel prints its errors on stdout, which is why the test passes; same pattern for both of its error branches (name_taken and bad_cwd). example: a live terminal named probe-ab12cd, then herdr-hc.sh tab create --cwd /tmp --label AB12CD/probe -> prints the name_taken JSON from hc's stderr, exits 1, no terminal created; ran: that against real hc 0.1.20 (exit 1), and hc run --bg --json --name <taken> sleep 1 >o 2>e (stdout empty, name_taken JSON on stderr, exit 1)
- **2026-11-22 19:36 · Alexander Sacharov** — bughunt: done at 56f111c29 - 1 finding (1 change, 0 pre-existing) in 4 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/release/notes/MHCEN8.md core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh core/role/role_test.go core/shipped/shipped.sha256
