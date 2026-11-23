---
id: 01M5ZJ4WRN0FRAKQCCB805CSWX
title: "spawn.sh startet keinen Worker, wenn der freie Speicher nicht fuer einen weiteren Build reicht"
status: testing
ready: true
creator: BeMuCa
goal: "Vor dem Start eines Workers liest spawn.sh den freien Speicher und lehnt ab oder warnt, wenn ein weiterer paralleler Build den Rechner zum Absturz braechte; die Grenze steht an einer Stelle und laesst sich setzen."
context: |-
  Was schief ist: drei parallele Go-Builds auf einer WSL mit 7 GB fuehrten zu vier Abstuerzen und einmal zu 29 leeren Git-Objekten (Teamlead-Session 27.10.2026, Feedback L317). Nichts in jaira warnt vor Speicher.
  Heute: core/role/builtin/jaira-dispatcher/scripts/spawn.sh startet Worker ohne Blick auf die Maschine (grep MemAvailable findet nichts).
  Bekannt: /proc/meminfo MemAvailable reicht unter Linux/WSL; ein Go-Build mit -race braucht grob 1-2 GB.
  Nicht gewollt: ein Scheduler. Eine Zahl (freier Speicher in MB je Worker, Default) und ein Abbruch mit klarer Meldung reichen; der Teamlead entscheidet dann, ob er wartet.
definition-of-done: "spawn.sh lehnt mit Exit 3 und Meldung ab, wenn MemAvailable unter der Grenze liegt; --force startet trotzdem"
tags:
  - concurrency
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:52:06Z
updated-at: 2026-11-23T01:40:12Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: worker-05cswx
claimed-at: 2026-11-23T01:32:21Z
outcome-what: "spawn.sh prueft MemAvailable vor Worktree und Tab: unter JAIRA_MIN_FREE_MB (Default 2048) Exit 3, --force startet trotzdem, ohne /proc/meminfo eine Zeile und Start. run-lane.sh meldet das als Exit 6. Teamlead-, Dispatcher- und Worker-SKILL sagen, was dann zu tun ist."
outcome-why: Drei parallele Go-Builds auf einer 7-GB-WSL stuerzten am 27.10. viermal ab; nichts warnte vor Speicher.
outcome-resolves: "Kein Worker startet mehr still in eine Maschine, die fuer seinen Build zu wenig freien Speicher hat."
review-summary: none
review-gaps: "none: round 2 adds a trap and two mappings, nothing duplicated or orphaned; own unreleased hashes from 12cccd2 replaced in shipped.sha256 rather than piling up"
test-verdict: "pass: round 2 over 12cccd2..6bc8af7: go vet + go build ./... green (go test not run, house rule). Extracted check of the committed spawn.sh (trap included): 976 MB -> 3, force -> 0, missing meminfo / no MemAvailable -> one line + 0, limit 0 -> 0, 20-digit limit -> 2, 9-digit limit -> 3, abc -> 2, discuss -> 0, a later command failing with 3 -> 1, with 5 -> 5. run-lane.sh copy with stub spawn: 3->6, 4->1, 5->1, 2->1. spawn.sh not run end to end (opens tabs). DoD 5 left for the person."
hunt-verdict: "fail: bughunt done at 12cccd2ee, 2 change: run-lane read any spawn.sh exit 3 (e.g. herdr after the tab exists) as 'nothing started' -> second tab; huge JAIRA_MIN_FREE_MB skipped the check"
round: 2
back-from: testing
back-by: lane
back-count: 1
---

# spawn.sh startet keinen Worker, wenn der freie Speicher nicht fuer einen weiteren Build reicht

## Definition of Done

- [x] spawn.sh lehnt mit Exit 3 und Meldung ab, wenn MemAvailable unter der Grenze liegt; --force startet trotzdem
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:140 exit 3 (EXIT trap keeps 3 for this refusal only), core/role/builtin/jaira-dispatcher/scripts/spawn.sh:110 --force; extracted check vs fake meminfo 976 MB -> rc=3, force=1 -> started
- [x] Die Grenze ist eine Variable am Kopf von spawn.sh mit Kommentar, warum der Default so ist
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:7-12 min_free_mb with comment (1-2 GB per -race build, 2026-10-27 crashes), JAIRA_MIN_FREE_MB
- [x] jaira-teamlead/SKILL.md sagt, was der Teamlead tut, wenn spawn.sh ablehnt
  proof: core/role/builtin/jaira-teamlead/SKILL.md:237 wait, tell the person, --force only on their word
- [x] Ohne /proc/meminfo oder ohne MemAvailable (macOS, Git Bash) ueberspringt spawn.sh die Pruefung mit einer Zeile auf stderr und startet; nie ein Abbruch
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:134; extracted check with missing file and with no MemAvailable line -> one stderr line, rc=0
- [ ] seen by a person: spawn.sh auf macOS oder unter Git Bash gestartet zeigt die Zeile, dass der Speicher nicht geprueft wurde, und oeffnet den Tab

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 01:33 · BeMuCa** — Entscheidungen (keine offen fuer die Person): Grenze min_free_mb am Kopf von spawn.sh, setzbar mit JAIRA_MIN_FREE_MB, Default 2048 MB (ein -race-Build 1-2 GB). 0 schaltet ab, Nicht-Zahl -> Exit 2. Exit 3 bei zu wenig MemAvailable, vor Worktree und Tab. --force startet trotzdem. Ohne /proc/meminfo oder ohne MemAvailable-Zeile (macOS, Git Bash, Kernel <3.14): eine Zeile auf stderr, Start geht weiter. Ausgenommen sind die Lanes explain und discuss: das sind Chats fuer eine Person aus den Cockpit-Knoepfen, kein Build. Bekannte Grenze: gemessen wird der Stand jetzt; zwei Spawns direkt hintereinander bestehen beide, weil der erste noch nicht baut - kein Scheduler, wie der Kontext verlangt. TestMain in core/role setzt JAIRA_MIN_FREE_MB=0, damit die Skript-Tests auf knappen Maschinen nicht an der Pruefung scheitern.
- **2026-11-23 01:35 · BeMuCa** — run-lane.sh: spawn.sh lief dort unter set -e in einer Zuweisung, sein Exit 3 waere als run-lane-Exit 3 (= Timeout, Worker laeuft noch, keinen zweiten starten) gelesen worden. Daher bildet run-lane.sh spawns 3 auf 6 ab (nichts gestartet); Dispatcher- und Worker-SKILL nennen 6. Geprueft mit Stub-spawn.sh: 3->6, 1->1. spawn.sh selbst wurde nicht echt ausgefuehrt (oeffnet Tabs; vom Klassifizierer verweigert) - das --force-Parsing ist nur gelesen, die Pruefung als ausgeschnittener Block getestet.
- **2026-11-23 01:36 · BeMuCa** — critique: nichts. Erwogen und stehen gelassen: (a) Pruefung in run-lane statt spawn.sh - nein, Teamlead und Cockpit rufen spawn.sh direkt; (b) JAIRA_MIN_FREE_MB=0 neben --force doppelt - nein, 0 gilt fuer eine ganze Sitzung/die Tests, --force fuer einen Aufruf; (c) Zahlenpruefung des Env-Werts - bleibt, sonst wuerde [ -gt ] mit Text die Pruefung still ueberspringen.
- **2026-11-23 01:36 · BeMuCa** — testing: pause - bughunt laeuft im Hintergrund auf 12cccd2, test-verdict steht; hunt-verdict fehlt noch.
- **2026-11-23 01:38 · BeMuCa** — bughunt: change: core/role/builtin/jaira-dispatcher/scripts/run-lane.sh:79-83 - class: spawn.sh runs under set -e, so any failing command inside it (herdr/hc via herdr-hc.sh, jaira json get) leaves with its own exit code, and run-lane maps 3 to 6 and passes every other code through unmapped, where 3/4/5/6 already mean worker states; the new 6 mapping only covers the one code spawn.sh uses on purpose. example 1: a herdr 'pane run' that exits 3 after 'tab create' made a tab -> run-lane exits 6 'nothing was started', the dispatcher waits and starts a second tab beside the first; example 2 (older, same class): tab create answering JSON without result.root_pane.pane_id -> jaira json get exits 5 -> run-lane exits 5 'the worker stopped on a question' with no question on the ticket. Fix the class: give spawn.sh's deliberate memory refusal a code nothing else produces, or map every non-zero spawn rc other than the memory refusal to one 'spawn failed' code. Not a blocker: whether the real herdr/hc ever exit 3 or 4 was not checked. ran: stub herdr + copies of spawn.sh/run-lane.sh in a scratch repo, run-lane.sh --no-worktree TST001 critique (pane run exit 3 -> rc=6, tab create logged; pane_id missing -> rc=5)
- **2026-11-23 01:38 · BeMuCa** — bughunt: change: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:123-127 - class: JAIRA_MIN_FREE_MB is checked for digits only, not for size; a value bash cannot hold as an integer makes [ -gt 0 ] fail with 'integer expression expected', the if is false and the memory check is skipped instead of refused (exit 2 like a non-number). example: JAIRA_MIN_FREE_MB=99999999999999999999 -> spawn.sh goes on to create the tab and start claude with 2694 MB free. Edge, not a blocker. ran: JAIRA_MIN_FREE_MB=99999999999999999999 spawn.sh --no-worktree tst001 TST001 critique <scratch repo> with a stub herdr (reached 'pane run')
- **2026-11-23 01:38 · BeMuCa** — bughunt: done at 12cccd2ee - 2 findings (2 change, 0 pre-existing) in 2 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/release/notes/05CSWX.md core/role/builtin/jaira-dispatcher/SKILL.md core/role/builtin/jaira-dispatcher/scripts/run-lane.sh core/role/builtin/jaira-dispatcher/scripts/spawn.sh core/role/builtin/jaira-teamlead/SKILL.md core/role/builtin/jaira-worker/SKILL.md core/role/main_test.go core/shipped/shipped.sha256
- **2026-11-23 01:39 · BeMuCa** — testing: class: exit-code collision zwischen spawn.sh und run-lane.sh; example: bughunt note 1 (herdr exit 3 -> run-lane 6, zweiter Tab) und note 2 (20-stellige Grenze -> Pruefung uebersprungen); fix: spawn.sh EXIT-Trap macht jede nicht absichtliche 3 zu 1, run-lane bildet 3->6 und jede andere Nicht-Null auf 1 ab (deckt auch das aeltere 5-Beispiel); Grenze mit mehr als 9 Ziffern -> Exit 2.
- **2026-11-23 01:39 · BeMuCa** — critique runde 2: nichts. EXIT-Trap statt eigenem Exit-Code, weil DoD 1 Exit 3 vorschreibt; zwei Zeilen, kein Umbau der uebrigen Exits.
