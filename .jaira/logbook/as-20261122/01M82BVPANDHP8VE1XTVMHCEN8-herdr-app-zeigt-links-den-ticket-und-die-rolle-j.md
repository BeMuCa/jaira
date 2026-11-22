---
id: 01M82BVPANDHP8VE1XTVMHCEN8
title: "Herdr App zeigt links den Ticket und die Rolle jeder Agent-Tab (worker-8tt890, pr-watjdq) statt 'Claude Code'"
status: done
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
commits:
  - 81cccfc1138ec175b57988efefeefa17e7617d4f
  - 56f111c29e6807ea845915271ec2e5c2580ca117
  - b0ffc7062608550e29a55a7035cb20e8b7b3ece0
created-at: 2026-11-22T19:31:20Z
updated-at: 2026-11-22T19:43:20Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3820322
claimed-at: 2026-11-22T19:31:51Z
outcome-what: "herdr-hc.sh 'tab create' gibt hc run --name <rolle>-<id klein> aus dem Label <ID>/<rolle>; bei name_taken (aus hc's stderr gelesen) -2..-9, dann ohne Namen; Test mit stderr-Stub, shipped.sha256, Release-Note"
outcome-why: "Herdr App zeigt eine Kind-Tab nur mit ihrem hc-Namen; ohne Namen stand links zehnmal 'Claude Code'"
outcome-resolves: "Alex sieht in der Agent-Leiste, welche Rolle an welchem Ticket arbeitet"
review-summary: "herdr-hc.sh 'tab create' builds an hc name from spawn.sh's label <ID>/<lane> as <lane>-<id lowercase>, forced into a-z0-9_- and cut to 32, and passes it as hc run --name; Herdr App's agent list then shows worker-8tt890 instead of 'Claude Code'. When hc refuses with name_taken (read from hc's stderr, where hc 0.1.20 writes it) it retries with -2..-9 (shortening the base so the total stays 32), then without a name, so a tab always opens; any other hc refusal is printed and exits 1 as before. No label: no --name, unchanged. Test TestHerdrHcNamesTheTabAfterItsLabel with a fake hc covers names, odd characters, length, taken once/twice/all, other refusal; shipped.sha256 and release note added."
review-gaps: |-
  note: on the success path hc's stderr now goes into the temp file and is dropped - a warning hc prints while still starting the terminal is no longer visible; cat "$err" >&2 after the loop would keep it
  note: a label whose lane part maps to a leading dash would pass --name -x to hc, which may read it as a flag; unreachable today (labels are <ID>/<lane>), a ${name#-} strip would close it
test-verdict: "pass: go test -race ./role ./shipped green (RC=0; whole suite RC=0 in round 1); DoD 1-3 verified in the tree, 4 (seen by a person) left for the person; real hc 0.1.20: two 'tab create --label ZZ99QQ/worker' opened worker-zz99qq and worker-zz99qq-2; missing --cwd: hc's spawn_failed on stderr, rc 1, no terminal, no temp file left; two writers: names are claimed by hc itself, the loser of a race gets name_taken and takes the next suffix"
hunt-verdict: "pass: bughunt done at b0ffc7062, 0 change, 0 pre-existing"
round: 2
back-from: testing
back-by: lane
back-count: 1
review-verdict: "accept: DoD 1-3 met in code and tests; live check against hc 0.1.20: two tab create --label RV99QQ/probe opened probe-rv99qq and probe-rv99qq-2; go test ./role -run TestHerdrHc and ./shipped pass, sha256 matches shipped.sha256; DoD 4 is the person's"
review-round: 2
review-verdict-gaps: |-
  note: on the success path hc's stderr now goes into the temp file and is dropped - a warning hc prints while still starting the terminal is no longer visible; cat "$err" >&2 after the loop would keep it
  note: a label whose lane part maps to a leading dash would pass --name -x to hc, which may read it as a flag; unreachable today (labels are <ID>/<lane>), a ${name#-} strip would close it
review-check: "1. In Herdr App, in a jaira repo, start a worker the usual way (cockpit 'start' on a ticket, or: core/role/builtin/jaira-dispatcher/scripts/spawn.sh <ID> worker)  2. Look at the agent list on the left: the new tab reads worker-<id in lowercase>, e.g. worker-8tt890, not 'Claude Code'  3. Start a second worker for the same ticket while the first is still open: it opens and reads worker-<id>-2  4. Close both tabs. Automated: cd core && go test ./role -run TestHerdrHcNamesTheTabAfterItsLabel"
---

# Herdr App zeigt links den Ticket und die Rolle jeder Agent-Tab (worker-8tt890, pr-watjdq) statt 'Claude Code'

## Definition of Done

- [x] herdr-hc.sh 'tab create' gibt hc run --name '<rolle>-<id klein>' (aus dem Label <ID>/<rolle>, z.B. worker-8tt890, pr-watjdq, review-…), nur erlaubte Zeichen, höchstens 32 (Test)
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh tab create; TestHerdrHcNamesTheTabAfterItsLabel (worker/pr/review/odd characters/too long)
- [x] Ist der Name schon vergeben (alte Tab desselben Tickets lebt noch), startet die Tab trotzdem: mit Suffix -2, -3 … oder ohne Namen (Test)
  proof: herdr-hc.sh reads hc's stderr for name_taken; TestHerdrHcNamesTheTabAfterItsLabel (taken once/twice/all taken, stub refuses on stderr like hc); real hc: second ZZ99QQ/worker opened as worker-zz99qq-2
- [x] Ohne Label bleibt alles wie heute; shipped.sha256 aktualisiert, Release-Note core/release/notes/<ID>.md
  proof: TestHerdrHcNamesTheTabAfterItsLabel/no_label; core/shipped/shipped.sha256; core/release/notes/MHCEN8.md
- [-] seen by a person: in Herdr App eine Tab über spawn.sh starten (z.B. Cockpit-Worker) — links steht worker-<id> statt 'Claude Code'

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
- **2026-11-22 19:38 · Alexander Sacharov** — critique round 2: nichts. Temp-Datei für stderr statt 2>&1: auf dem Erfolgsweg bleibt stdout reines JSON für 'jaira json get address', auch wenn hc dabei etwas auf stderr schreibt.
- **2026-11-22 19:39 · Alexander Sacharov** — bughunt: done at b0ffc7062 - 0 findings in 4 modules, follow-up hunt since 56f111c29 (diff of herdr-hc.sh, role_test.go, shipped.sha256; release note unchanged); concurrency, error paths, input, permissions checked; ran: hc run --cwd /nonexistent (exit 1, refusal JSON on stderr, stdout empty), go test ./role -run TestHerdrHc and ./shipped (pass), sha256 of herdr-hc.sh matches shipped.sha256; files: core/release/notes/MHCEN8.md core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-22 19:43 · Alexander Sacharov** — DoD 4 (live in Herdr App) nach --force offen geblieben: wird nach dem Merge geprüft - erster spawn nach 'jaira roles install' muss links worker-<id> zeigen; sonst Folge-Ticket.
