---
id: 01M86H5VJ5DK1MDM8Y6NQNB5S4
title: "herdr-hc.sh 'hc get' has no timeout of its own - a hanging hc get hangs…"
status: optimize
ready: true
creator: Sascha
context: |-
  herdr-hc.sh 'tab list' (core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:82) ruft hc ls und je Terminal hc get ohne Frist auf.
  Haengt eins davon, haengt 'tab list' - und damit jeder Aufrufer.
  reap.sh:133 haengt dann unter reap.lock: jede spaetere Runde wird uebersprungen, nichts wird mehr abgeraeumt.
  Das Cockpit (register.tsx:212) bricht nach 15 s ab und sieht dann gar keine Tabs, nicht nur das eine stumme Terminal.
  Gefunden bei 8NWE4P (Exit 3 fuer stumme Terminals); dort blieb der Haenger selbst offen.
tags: []
blocked-by: []
related: []
follows: 01M866HXTHQS29K2ZMNQ8NWE4P
commits: []
created-at: 2026-11-24T10:21:14Z
updated-at: 2026-11-26T08:09:23Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-36559
claimed-at: 2026-11-26T07:54:46Z
goal: "herdr-hc.sh 'tab list' endet in begrenzter Zeit, auch wenn hc get (oder hc ls) haengt: ein Terminal, dessen hc get nicht in der Frist antwortet, zaehlt als stumm (Exit 3 + stderr, wie seit 8NWE4P)"
outcome-what: "herdr-hc.sh: bounded() gibt hc ls und jedem hc get in 'tab list' (und dem Space-Lookup in 'tab create') JAIRA_HC_TIMEOUT Sekunden (3), dann kill; Ausgabe ueber Dateien. Test um haengende get/ls erweitert, Manifest und Release-Note."
outcome-why: "Ein Terminal, dessen hc get nie antwortete, liess 'tab list' und damit reap.sh (unter reap.lock), Give back und die Worker-Liste des Cockpits haengen."
outcome-resolves: "DoD 1-4: haengendes get -> stumm, Exit 3; haengendes ls -> Exit !=0; reines bash ohne Prozessreste; Release-Note."
review-summary: none
review-gaps: "none: nichts entfernt. Geprueft: Duplikat - reap.sh:389 timeout(1)-Probe tut Aehnliches, bleibt getrennt (siehe critique-Note); tote Pfade - keine, alle drei bounded-Aufrufer laufen; Kosten - je hc-Aufruf 2 mktemp + Watchdog-Subshell + sleep, bei N Terminals ~4N Forks mehr pro 'tab list', gemessen 1 s Lauf mit 2 Terminals inkl. 1 s Frist, also unkritisch; gelassen."
test-verdict: "pass: go vet + go test -race ./... nach go clean -testcache gruen (RC=0); shipped-history --check ok; DoD 1-4 im Baum geprueft (TestHerdrHcTabListSaysWhichTerminalDidNotAnswer inkl. one_hangs/all_hang/hangs_and_mute/ls_hangs). Funktion mit hc-Stub: t2 haengt -> nach 3,08 s Liste mit t1, stderr nennt t2, Exit 3; zwei 'tab list' parallel -> beide Exit 3, keine tmp-Reste; JAIRA_HC_TIMEOUT=x -> sofort Exit 137 (sicherer Fehler). Schreibt nichts ausser eigenen mktemp-Dateien."
hunt-verdict: "fail: bughunt done at 21d6e7e89, 2 change (keine Gesamtfrist fuer 'tab list'; JAIRA_HC_TIMEOUT ungeprueft), 1 pre-existing (tab get/pane get ohne Frist - gleiches Thema, Note fuer die Person, deckt sich mit der 0TREV1-Note)"
round: 2
back-from: testing
back-by: lane
back-count: 1
---

# herdr-hc.sh 'hc get' has no timeout of its own - a hanging hc get hangs…

## Definition of Done

- [x] herdr-hc.sh 'tab list' endet innerhalb eines Budgets fuer den ganzen Aufruf (Standard 5 s, JAIRA_HC_TIMEOUT), auch wenn mehrere hc get haengen; jedes nicht gehoerte Terminal steht als stumm auf stderr und die Liste endet mit Exit 3; Test mit einem hc-Stub, dessen get schlaeft
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:36-63 (Budget end=SECONDS+secs, bounded), :133; TestHerdrHcTabListSaysWhichTerminalDidNotAnswer/one_hangs, all_hang, hangs_and_mute, five_hang (<5 s bei 2 s Budget)
- [x] Ein haengendes hc ls laesst 'tab list' nach derselben Frist mit Exit ungleich 0 enden; Test
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:130; TestHerdrHcTabListSaysWhichTerminalDidNotAnswer/ls_hangs
- [x] Ohne timeout(1): die Frist laeuft in reinem bash (bash 3.2 tauglich), ohne liegengebliebene Prozesse
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:41-63 sleep+kill-Watchdog, kein timeout(1), Ausgabe ueber Dateien; Stub ohne exec in role_test.go, Test misst die Dauer
- [x] Release-Note in core/release/notes/QNB5S4.md
  proof: core/release/notes/QNB5S4.md
- [x] Ein JAIRA_HC_TIMEOUT, das keine ganzen Sekunden ueber 0 sind, wird mit Exit 2 abgelehnt, bevor hc gefragt wird; Test
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:37-39; TestHerdrHcRefusesATimeoutThatIsNoSeconds (abc, 0, -1, 1.5, ' 3')

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-26 07:54 · Sascha** — parent session: uds:/run/user/1000/cc-socks/62036.sock
- **2026-11-26 07:55 · Sascha** — decision (worker): Umfang nur 'tab list' (hc ls + hc get darin). 'tab get'/'pane get' bleiben ohne Frist: reap.sh:417 liest ein gescheitertes pane get als 'Tab zu' und raeumte dann den Worktree eines offenen Tabs ab - ein Haenger dort ist sicherer als ein Timeout. Frist 3 s: Cockpit gibt 'tab list' 15 s gesamt (register.tsx:212), so passen mehrere stumme Terminals. Watchdog in bash statt timeout(1): macOS hat keins, Git Bash findet timeout.exe (siehe reap.sh:389).
- **2026-11-26 08:01 · Sascha** — Dead end: hc-Ausgabe direkt in die Pipe des Aufrufers - ein Kind des getoeteten hc (Stub ohne exec: sleep) hielt stderr offen, Go/Cockpit warteten trotzdem 20 s. Darum schreibt bounded() stdout UND stderr in Dateien und gibt sie danach aus. Vier Kategorien: Nebenlaeufigkeit - jeder Aufruf eigene mktemp-Dateien; PID-Reuse ausgeschlossen, solange hc ungewartet Zombie ist, Restfenster zwischen wait und kill des Watchdogs ist Mikrosekunden. Fehlerpfade - Absage von hc behaelt Exit und stderr (not_found-Erkennung unveraendert), Ablauf -> Exit 137 + 'gave no answer within Ns'. Eingabe - JAIRA_HC_TIMEOUT nicht numerisch/0 -> sleep endet sofort, alles gilt als stumm, Exit 3 = sicherer Fehler. Rechte - nichts Neues. Manifest core/shipped/shipped.sha256 neu geschrieben (TestEveryInstalledFileIsAShippedVersion).
- **2026-11-26 08:02 · Sascha** — critique: nichts. Ein Muster neben reap.sh:389 (timeout(1)-Probe) ist hier gewollt: reap.sh braucht den Gruppen-Kill fuer docker compose und verweigert ohne coreutils; hier hiesse Verweigern, dass 'tab list' auf einem Mac ohne coreutils immer scheitert. JAIRA_HC_TIMEOUT hat einen Aufrufer (der Test) und steht in der Release-Note. Ort stimmt: im Aufgerufenen, kein Cockpit-Code.
- **2026-11-26 08:02 · Sascha** — not this ticket (teamlead 26.11.): 0TREV1 (feat/0trev1, 2fc5359c9) laesst reap.sh und close-tab.sh nur not_found als geschlossenen Tab lesen, ein Timeout behaelt den Tab. Sind beide gelandet, kann auch 'tab get'/'pane get' in herdr-hc.sh eine Frist bekommen - hier bewusst nicht, nichts auf 0TREV1s ungemergtem Code gebaut.
- **2026-11-26 08:02 · Sascha** — waiting: bughunt (background agent) since 09:02
- **2026-11-26 08:04 · Sascha** — bughunt: change: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:118-130 - class: the deadline applies to each hc call, not to the whole tab list; terminals are asked one after another, so the list takes up to 3 s (ls) + 3 s for each mute terminal and has no total limit; example: 6 terminals, 5 of them hung, take 15.15 s, which is past the cockpit's 15 s limit (register.tsx:212, run() timeoutMs=15_000); the cockpit then gets 'gave no answer' (safe, read as tabs unknown) but still blocks for 15 s on every refresh, which is what this change was meant to stop; reap.sh has no outer limit and holds its lock for N x 3 s; ran: a stub hc with 5 hung terminals, herdr-hc.sh tab list -> rc=3, took 15.15s. Not a blocker: the result is the same safe 'unknown' as before, only the time is wrong
- **2026-11-26 08:04 · Sascha** — bughunt: change: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:33 - class: JAIRA_HC_TIMEOUT is not checked; any value sleep refuses kills every hc call at once; example: JAIRA_HC_TIMEOUT=0 (a person may read it as 'no limit'), abc or -1 -> sleep fails at once, the watchdog kills hc, tab list exits 137 with 'hc ls gave no answer within 0s', and tab create quietly loses the herdr.space label so the new tab lands outside the space; reap.sh:154 already checks JAIRA_TEARDOWN_TIMEOUT with case ''|*[!0-9]*|0 -> exit 2; ran: JAIRA_HC_TIMEOUT=abc|0|-1 herdr-hc.sh tab list with an hc that answers in 0.3 s -> rc=137 each, 0.5 -> rc=0
- **2026-11-26 08:05 · Sascha** — bughunt: pre-existing: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:138,151 (tab get, pane get), also tab close 143-148, pane run 161, send-keys 167 - class: hc calls with no time limit that are still reached while reap.sh holds its lock; example: a terminal answers tab list, then stops answering while reap.sh waits on gh or worktree-teardown (up to 300 s), and reap.sh:418-419 'pane get'/'tab get' hang forever, the lock stays held and no later reap runs - the same hang this ticket fixed for tab list; close-tab.sh:159-162 is reached the same way from reap.sh:142; blame: 655874a6b9 (2026-11-11, before the ticket); ran: herdr-hc.sh pane get on a hung stub terminal under timeout 8 -> rc=124
- **2026-11-26 08:05 · Sascha** — hunt report: 3 findings (2 change - neither a blocker, 1 pre-existing), all as notes on QNB5S4 | read whole: herdr-hc.sh, the herdr-hc part of role_test.go, release note QNB5S4.md, shipped.sha256 line; callers read along: reap.sh, close-tab.sh, cockpit register.tsx run()/tab list | not checked: whether bash kill -KILL stops a native Windows hc.exe under Git Bash (no Windows here); whether the test's JAIRA_HC_TIMEOUT=1 kills a healthy stub hc on a slow CI runner and makes 'all answer' fail (not reproduced)
- **2026-11-26 08:05 · Sascha** — bughunt: done at 21d6e7e89 - 3 findings (2 change, 1 pre-existing) in 4 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/release/notes/QNB5S4.md core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-26 08:05 · Sascha** — testing: class: Frist gilt je hc-Aufruf, nicht fuer die ganze Liste; example: 5 haengende Terminals -> 15,15 s, das Cockpit (register.tsx:212, 15 s) laeuft trotzdem in seinen Timeout - genau das sollte das Ticket beenden; fix: eine Gesamtfrist fuer 'tab list' (Rest der Frist an jedes get, abgelaufen -> restliche Terminals stumm ohne zu fragen), Test mit mehreren haengenden Terminals unter der Cockpit-Grenze. Mit dabei: JAIRA_HC_TIMEOUT wie reap.sh:154 pruefen (leer/nicht numerisch/0 -> Standard oder Exit 2), Test ueber abc/0/-1.
- **2026-11-26 08:08 · Sascha** — Runde 2: Budget gilt jetzt fuer den ganzen Aufruf (5 s), nicht je hc. Folge: nach einem Haenger werden die restlichen Terminals nicht mehr gefragt und zaehlen als stumm. Bewusst so statt paralleler hc get: alle Aufrufer verwerfen stdout bei Exit != 0 (register.tsx:358 agentsOf, :1992 refreshBoard, :1825 Give back, reap.sh:133), also liest niemand, was sie gesagt haetten - parallel waere mehr Code fuer nichts. Gemessen: 5 haengende von 6 -> 5,1 s, keine Prozessreste. Ungueltiges JAIRA_HC_TIMEOUT -> Exit 2 wie reap.sh:154 (statt still alles stumm und tab create ohne herdr.space).
- **2026-11-26 08:09 · Sascha** — critique (Runde 2): nichts. Ganzes Budget statt Parallel-get ist die kleinere Form (siehe Runde-2-Note); Pruefung von JAIRA_HC_TIMEOUT folgt dem Muster reap.sh:154.
