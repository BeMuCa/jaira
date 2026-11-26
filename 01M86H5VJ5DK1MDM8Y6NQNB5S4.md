---
id: 01M86H5VJ5DK1MDM8Y6NQNB5S4
title: "herdr-hc.sh 'hc get' has no timeout of its own - a hanging hc get hangs…"
status: done
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
commits:
  - 21d6e7e89f4107c2dd795739df33fdd32f8dd46a
  - ff4b72ae3fd900b4401433f263178b4afa16f2d2
  - 5e9dfcf8281b07421c566533b6068b8c3364cba8
created-at: 2026-11-24T10:21:14Z
updated-at: 2026-11-26T09:24:12Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-36559
claimed-at: 2026-11-26T07:54:46Z
goal: "herdr-hc.sh 'tab list' endet in begrenzter Zeit, auch wenn hc get (oder hc ls) haengt: ein Terminal, dessen hc get nicht in der Frist antwortet, zaehlt als stumm (Exit 3 + stderr, wie seit 8NWE4P)"
outcome-what: "herdr-hc.sh: bounded() gibt hc ls und jedem hc get in 'tab list' (und dem Space-Lookup in 'tab create') JAIRA_HC_TIMEOUT Sekunden (3), dann kill; Ausgabe ueber Dateien. Test um haengende get/ls erweitert, Manifest und Release-Note."
outcome-why: "Ein Terminal, dessen hc get nie antwortete, liess 'tab list' und damit reap.sh (unter reap.lock), Give back und die Worker-Liste des Cockpits haengen."
outcome-resolves: "DoD 1-4: haengendes get -> stumm, Exit 3; haengendes ls -> Exit !=0; reines bash ohne Prozessreste; Release-Note."
review-summary: "Full review (no earlier review) of 21d6e7e89..5e9dfcf82. herdr-hc.sh gets bounded(): every hc call in 'tab list' (hc ls, each hc get) and the herdr.space lookup in 'tab create' runs in the background with its output in temp files, and a plain-bash sleep+kill watchdog kills it once the call's single budget is spent. JAIRA_HC_TIMEOUT (default 5 s, read as decimal) is that budget for the whole herdr-hc.sh call, at most secs+1 s; a terminal not heard from in time, or not asked any more because the budget was spent, is named on stderr and the list exits 3 as since 8NWE4P; a hung hc ls fails the list. An invalid JAIRA_HC_TIMEOUT (empty, non-digits, 0) stops herdr-hc.sh with exit 2 before hc is asked. Test: hc stub with hanging get/ls, a 2 s budget, five hung terminals under 5 s; separate test for bad values. Release note and shipped.sha256 updated. Checked by hand: 3 hung terminals at 5 s -> 6.0 s, exit 3, no processes left; go test ./core/role ./core/shipped green."
review-gaps: |-
  note: edge: JAIRA_HC_TIMEOUT longer than bash's 64-bit range passes the digit check and wraps in $((10#$secs)) - 99999999999999999999 becomes ~7.7e18 s, so a hung hc get waits without limit again; fix: refuse more than 6 digits in the case at herdr-hc.sh:40-43, plus one test value
  note: limit: the budget is the whole call's - once it is spent the remaining terminals are not asked and are named as mute even when they would answer (checked: t3 answers, still named); every caller drops an exit-3 list, so nothing reads them (ticket note round 2)
  note: limit: an invalid JAIRA_HC_TIMEOUT now stops every herdr-hc.sh command with exit 2, tab create and pane run included, not only tab list; the release note says so
  note: limit: the watchdog kills hc alone, not processes hc started; a real hc is one binary, and whether kill -KILL stops a native hc.exe under Git Bash was not checked (no Windows here)
  note: elsewhere: tab get, pane get, tab close, pane run and send-keys in herdr-hc.sh:157-188 still call hc without a limit and are reached from reap.sh:418-419 under its lock; deliberately left until 0TREV1 lands (ticket note 08:02)
test-verdict: "pass: Runde 3 - go vet + go test -race ./... gruen (RC=0) auf 5e9dfcf82, TestHerdrHc* 3x hintereinander gruen; shipped-history --check ok. DoD 1-5 im Baum: 08 -> Exit 0 (8 s), 00/0/abc/-1/1.5/' 3' -> Exit 2 ohne hc-Aufruf; 1 haengendes Terminal bei secs=1/5 -> 1,06 s/5,06 s, Exit 3, stderr nennt es; Nachbar der Runde-2-Note (010) geprueft: 10 s, nicht 8."
hunt-verdict: "pass: bughunt done at 5e9dfcf82, 0 Blocker, 1 change-edge (Ueberlauf bei >18 Ziffern in JAIRA_HC_TIMEOUT, Note fuer Review/Person)"
round: 3
back-from: testing
back-by: lane
back-count: 2
review-verdict: "accept: all five DoD items hold in the code and tests; no blocker, one edge (oversized JAIRA_HC_TIMEOUT wraps) and four limits"
review-round: 3
review-verdict-gaps: |-
  note: edge: JAIRA_HC_TIMEOUT longer than bash's 64-bit range passes the digit check and wraps in $((10#$secs)) - 99999999999999999999 becomes ~7.7e18 s, so a hung hc get waits without limit again; fix: refuse more than 6 digits in the case at herdr-hc.sh:40-43, plus one test value
  note: limit: the budget is the whole call's - once it is spent the remaining terminals are not asked and are named as mute even when they would answer (checked: t3 answers, still named); every caller drops an exit-3 list, so nothing reads them (ticket note round 2)
  note: limit: an invalid JAIRA_HC_TIMEOUT now stops every herdr-hc.sh command with exit 2, tab create and pane run included, not only tab list; the release note says so
  note: limit: the watchdog kills hc alone, not processes hc started; a real hc is one binary, and whether kill -KILL stops a native hc.exe under Git Bash was not checked (no Windows here)
  note: elsewhere: tab get, pane get, tab close, pane run and send-keys in herdr-hc.sh:157-188 still call hc without a limit and are reached from reap.sh:418-419 under its lock; deliberately left until 0TREV1 lands (ticket note 08:02)
review-check: "1. cd into this worktree  2. go test ./core/role -run TestHerdrHc -count=1 -v  3. you see one_hangs, all_hang, hangs_and_mute, five_hang and ls_hangs PASS, each in about 3 s or less, and TestHerdrHcRefusesATimeoutThatIsNoSeconds PASS for abc, 0, 00, -1, 1.5 and ' 3'  4. in Herdr App, with one worker tab frozen: open the cockpit worker list - it comes back within about 6 s instead of hanging 15 s, and the frozen terminal is reported as not answering  5. JAIRA_HC_TIMEOUT=abc bash core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh tab list prints 'JAIRA_HC_TIMEOUT must be seconds above 0' and exit code 2"
---

# herdr-hc.sh 'hc get' has no timeout of its own - a hanging hc get hangs…

## Definition of Done

- [x] herdr-hc.sh 'tab list' endet innerhalb eines Budgets fuer den ganzen Aufruf (Standard 5 s, JAIRA_HC_TIMEOUT), auch wenn mehrere hc get haengen; jedes nicht gehoerte Terminal steht als stumm auf stderr und die Liste endet mit Exit 3; Test mit einem hc-Stub, dessen get schlaeft
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:39-68 (Budget end=SECONDS+secs+1, ein hc hoechstens secs), :143; TestHerdrHcTabListSaysWhichTerminalDidNotAnswer/one_hangs, all_hang, hangs_and_mute, five_hang (<5 s bei 2 s Budget)
- [x] Ein haengendes hc ls laesst 'tab list' nach derselben Frist mit Exit ungleich 0 enden; Test
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:138; TestHerdrHcTabListSaysWhichTerminalDidNotAnswer/ls_hangs
- [x] Ohne timeout(1): die Frist laeuft in reinem bash (bash 3.2 tauglich), ohne liegengebliebene Prozesse
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:48-68 sleep+kill-Watchdog, kein timeout(1), Ausgabe ueber Dateien; Stub ohne exec in role_test.go, Test misst die Dauer
- [x] Release-Note in core/release/notes/QNB5S4.md
  proof: core/release/notes/QNB5S4.md
- [x] Ein JAIRA_HC_TIMEOUT, das keine ganzen Sekunden ueber 0 sind, wird mit Exit 2 abgelehnt, bevor hc gefragt wird; Test
  proof: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:39-46 (10#, 0 abgelehnt); TestHerdrHcRefusesATimeoutThatIsNoSeconds (abc, 0, 00, -1, 1.5, ' 3'); 08 -> Exit 0, von Hand

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
- **2026-11-26 08:09 · Sascha** — waiting: bughunt Runde 2 (background agent) since 09:09
- **2026-11-26 08:10 · Sascha** — bughunt: change: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:37-40 - class: JAIRA_HC_TIMEOUT with leading zeros passes the digits-only guard, but $((SECONDS + secs)) reads it as octal; examples: 08/09 -> bash 'value too great for base' + 'end: unbound variable', exit 1 instead of the promised exit 2 (release note '- !' line says any other value is refused with exit 2); 00 -> passes the guard, budget 0, 'hc ls not asked', exit 1; 010 -> silently 8 s instead of 10; ran: JAIRA_HC_TIMEOUT=<v> bash herdr-hc.sh tab list with a stub hc (rc 1/1/1/0, stub asked nothing for 08/09/00); not a blocker - every caller still reads the failure as tabs unknown; fix: strip leading zeros or use 10#$secs and refuse a result of 0
- **2026-11-26 08:10 · Sascha** — bughunt: change: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:40-43,52 - class: the budget is counted in whole SECONDS ticks, so it is secs minus up to 1 s depending on where in the wall-clock second the script starts; the watchdog sleeps the truncated 'left' too; example: JAIRA_HC_TIMEOUT=1 with an hc that answers at once reports a live terminal as mute and exits 3 when started just before a second tick; ran: 10 runs of tab list started ~20 ms before a tick with a stub hc (2/10 exit 3 'hc get gave no answer for t1'); with the default 5 the budget is 4-5 s, not a blocker; also the source of possible flakes in TestHerdrHcTabListSaysWhichTerminalDidNotAnswer at JAIRA_HC_TIMEOUT=2; fix: end=$((SECONDS + secs + 1)) or document the -1 s, or the guard minimum as 2
- **2026-11-26 08:10 · Sascha** — hunt report: 2 findings, both change, neither a blocker, as notes on QNB5S4 (octal JAIRA_HC_TIMEOUT; whole-second budget) | read as a diff since 21d6e7e89: herdr-hc.sh (plus the whole script for context), role_test.go, release note, shipped.sha256 (new hash matches the file); no new files | not checked: behaviour on macOS bash 3.2 and Git Bash, no such machine here
- **2026-11-26 08:10 · Sascha** — bughunt: done at ff4b72ae3 - 2 findings (2 change, 0 pre-existing) in 4 modules, follow-up hunt since 21d6e7e89; concurrency, error paths, input, permissions checked; files: core/release/notes/QNB5S4.md core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-26 08:11 · Sascha** — testing (Runde 2): class: Zahl aus der Umgebung wird von $(( )) als Oktal gelesen; example: JAIRA_HC_TIMEOUT=08 -> Exit 1, 00 -> Budget 0; fix: nach der Ziffern-Pruefung 10#$secs und 0 ablehnen, Test um 08/00/010 erweitern. Dazu: end um 1 s erhoehen, damit das Budget nie unter secs faellt (SECONDS zaehlt ganze Sekunden).
- **2026-11-26 08:15 · Sascha** — Runde 3: 10#$secs gegen Oktal (08 gilt jetzt als 8 s, 00 -> Exit 2); end=SECONDS+secs+1, ein einzelnes hc wartet hoechstens secs - gemessen 1 haengendes bei secs=1/5: 1,06 s/5,06 s. Folge im Test: ob Terminals nach einem Haenger noch gefragt werden, haengt nun vom Rest-Budget ab; der Test verlangt fuer Haenger-Faelle nur noch 'gelistet oder als stumm genannt', 3x hintereinander gruen.
- **2026-11-26 08:15 · Sascha** — critique+optimize (Runde 3, ein Durchgang, 39 Zeilen): nichts. 10# und die Kappung sind je eine Zeile an der einen Stelle, durch die jeder hc-Aufruf geht; keine Duplikate, nichts verwaist.
- **2026-11-26 08:15 · Sascha** — waiting: bughunt Runde 3 (background agent) since 09:15
- **2026-11-26 08:16 · Sascha** — bughunt: change: core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh:41-47 - class: JAIRA_HC_TIMEOUT digits beyond bash's 64-bit range pass the [!0-9] check and wrap silently in $((10#$secs)) and in end=$((SECONDS + secs + 1)); example: 9223372036854775807 makes end negative, so 'tab list' exits 1 with 'hc ls not asked: the 9223372036854775807s were up' without asking hc; 18446744073709551617 wraps to a 1 s budget, not refused; ran: JAIRA_HC_TIMEOUT=<v> herdr-hc.sh tab list against a stub hc (exit 1 / exit 0). Edge, not a blocker: nobody sets such a value, but a digit-count cap (e.g. more than 6 digits refused) would close the class
- **2026-11-26 08:16 · Sascha** — hunt report: 1 finding (1 change, edge, not a blocker) as a note on QNB5S4 | read as a diff since ff4b72ae3: herdr-hc.sh, role_test.go, shipped.sha256 (new hash matches the file); QNB5S4.md unchanged | not checked: real hc on Windows/Git Bash; SECONDS rounding was checked only by reading
- **2026-11-26 08:16 · Sascha** — bughunt: done at 5e9dfcf82 - 1 finding (1 change, 0 pre-existing) in 4 modules, follow-up hunt since ff4b72ae3; concurrency, error paths, input, permissions checked; files: core/release/notes/QNB5S4.md core/role/builtin/jaira-dispatcher/scripts/herdr-hc.sh core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-26 08:17 · Sascha** — testing: note: edge: JAIRA_HC_TIMEOUT mit mehr Ziffern als 64 Bit (9223372036854775807) laeuft in $(( )) ueber - Exit 1 ohne hc-Aufruf oder 1 s Budget; Fix waere eine Ziffernkappe (>6 Ziffern -> Exit 2); kein Blocker, niemand setzt solche Werte, Ergebnis bleibt 'Tabs unbekannt'.
- **2026-11-26 08:17 · Sascha** — waiting: review (fresh agent, run-lane.sh tab) since 09:17
- **2026-11-26 08:20 · Sascha** — review: done at 5e9dfcf82
- **2026-11-26 08:21 · Sascha** — report: herdr-hc.sh: ein Budget (JAIRA_HC_TIMEOUT, 5 s, dezimal) fuer jeden Aufruf, bash-Watchdog um hc ls/get in 'tab list' und den Space-Lookup; ungueltiger Wert -> Exit 2; Tests, Release-Note, Manifest | signoff, context unknown | offen: Abnahme durch Alex; Review: 1 edge (Ueberlauf >18 Ziffern), 4 limits; tab get/pane get erst nach 0TREV1; Branch nicht gepusht
