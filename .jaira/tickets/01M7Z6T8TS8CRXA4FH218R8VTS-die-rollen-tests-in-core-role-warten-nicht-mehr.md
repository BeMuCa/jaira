---
id: 01M7Z6T8TS8CRXA4FH218R8VTS
title: "Die Rollen-Tests in core/role warten nicht mehr auf Sekunden-Pausen, die für ein echtes Terminal gedacht sind"
status: testing
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "go test ./core/role läuft lokal in etwa 10 s statt 88 s, und der CI-Lauf wird dadurch spürbar kürzer (macOS und Windows-Shard heute ~7,5 min)"
context: |-
  Was falsch ist: go test ./core/role braucht lokal 88 s, im CI auf ubuntu 134 s (mit -race). Das ist der langsamste Paket-Test und bestimmt die Dauer von macOS (7m27) und dem langsamsten Windows-Shard (7m38).
  Ursache spawn: core/role/builtin/jaira-dispatcher/scripts/spawn.sh wartet mit 'sleep 2' in einer Schleife (Zeile ~353), bis claude im Tab hochkommt, und danach noch 'sleep 1' (Zeile ~413). Im Test gibt es kein Terminal und kein claude, also wartet jeder spawn-Aufruf ~3 s umsonst.
  Gemessen am 2026-11-21 (go test -json): TestSpawnTypesTheRoleForItsLaneName 18 s, TestSpawnChatsStartOnTheChatModelWithTheirBrief 15 s, TestSpawnParentIsTheCallerAndAnUnknownOneStops 12 s, TestSpawnStartsClaudeWhereItsShellRuns 9 s, TestSpawnRecordsItsTabForReap 6 s, TestSpawnLetsTheRepositorySetUpItsWorktree 6 s.
  Ursache reap: fünf TestReap*-Tests brauchen je ~6,5 s (TestReapRunsOneAtATime, TestReapClosesAMergedTicket, TestReapComesBackToAWorktreeItCouldNotRemove, TestReapLeavesATabItDidNotOpen, TestReapNeverRemovesTheWorktreeItRunsIn). Woher die 6,5 s kommen, ist noch nicht gefunden - reap.sh hat kein sleep; zuerst messen (Stub für gh/herdr? Retry?).
  Richtung: die Wartezeit per Umgebungsvariable einstellbar machen (Default unverändert für echte Läufe), Tests setzen sie auf 0. Kein Verhalten für echte Nutzer ändern.
  Nicht gemeint: Tests umschreiben oder streichen; Sharding in ci.yaml (das ist ein anderer Hebel und braucht workflow-Rechte).
definition-of-done: go test ./core/role läuft lokal unter 20 s; vorher/nachher-Zeit steht in einer Note
tags:
  - ci
blocked-by: []
related: []
commits: []
created-at: 2026-11-21T14:05:27Z
updated-at: 2026-11-21T14:11:43Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-842456
claimed-at: 2026-11-21T14:05:56Z
outcome-what: spawn.sh liest JAIRA_SPAWN_PAUSE für seine zwei Wartezeiten (Default 2 s/1 s unverändert); die core/role-Tests setzen sie in TestMain auf 0
outcome-why: "jeder spawn.sh-Aufruf im Test schlief 3 s auf ein claude, das es nicht gibt; newReapRig ruft spawn.sh zweimal, daher die 6,5 s pro TestReap*"
outcome-resolves: "go test ./core/role 89,5 s -> 4,4 s lokal"
review-summary: none
review-gaps: "reap_test.go: r.write + os.Chmod für den sleep-Stub zu einem os.WriteFile(...,0o755) zusammengelegt, wie newReapRig seine Stubs schreibt; sonst nichts - spawn.sh-Änderung sind zwei Parameter-Expansionen und der Kopftext, kein Duplikat gefunden (grep JAIRA_SPAWN/SPAWN_PAUSE)"
---

# Die Rollen-Tests in core/role warten nicht mehr auf Sekunden-Pausen, die für ein echtes Terminal gedacht sind

## Definition of Done

- [x] go test ./core/role läuft lokal unter 20 s; vorher/nachher-Zeit steht in einer Note
  proof: go test -count=1 ./core/role: 89,5 s -> 4,4 s (Note); core/role/main_test.go:15
- [x] spawn.sh und reap.sh warten in echten Läufen genauso lange wie vorher; die Variable ist im Kopfkommentar des Skripts beschrieben
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:57-60 (Kopf/Usage), :357 und :417 Default 2/1; TestSpawnPausesAsLongAsBeforeWithoutTheVariable (core/role/reap_test.go); reap.sh hat keine Pause
- [x] Die Ursache der ~6,5 s pro TestReap* ist benannt und behoben oder als bekannte Grenze notiert
  proof: newReapRig ruft spawn.sh zweimal (core/role/reap_test.go:126-127) = 2x3 s; behoben durch JAIRA_SPAWN_PAUSE=0 in TestMain

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-21 14:05 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/4205.sock
- **2026-11-21 14:08 · Alexander Sacharov** — Ursache reap gefunden: newReapRig (core/role/reap_test.go:126-127) ruft spawn.sh zweimal auf, jeder Aufruf schläft 2+1 s = 6 s pro Rig; reap.sh selbst wartet nie. TestUncommittedFailsWhereGitDoes (6,2 s) nutzt dasselbe Rig. Also ein Hebel: die Pausen in spawn.sh. reap.sh bekommt keine Variable, es hat keine Pause. Vorher: go test -count=1 ./core/role = 89,5 s (lokal, WSL, 2026-11-21).
- **2026-11-21 14:10 · Alexander Sacharov** — Nachher: go test -count=1 ./core/role = 4,4 s (vorher 89,5 s), mit -race 5,4 s; go test ./... grün. Hebel: JAIRA_SPAWN_PAUSE in spawn.sh, TestMain (core/role/main_test.go) setzt 0 für alle Tests des Pakets statt in jedem der ~10 cmd.Env-Blöcke. Leer/ungesetzt = 2 s und 1 s wie vorher; TestSpawnPausesAsLongAsBeforeWithoutTheVariable prüft das mit einem sleep-Stub auf PATH (Mutation 2->3 macht ihn rot). Kein Release-Note: die Variable ist ein Test-Hebel, kein Nutzer muss etwas anders tun und ohne sie ändert sich nichts. core/shipped/shipped.sha256 per scripts/shipped-history.sh --write ergänzt (Test verlangt es für jede geänderte ausgelieferte Datei).
- **2026-11-21 14:11 · Alexander Sacharov** — critique: nichts. Erwogen und verworfen: ein No-op-sleep auf PATH (in built.dir von jairaPath) statt der Variable - kein Eingriff in spawn.sh, aber er schaltet jedes sleep in jedem Skripttest ab (run-lane.sh, herdr-hc.sh) und ist für den nächsten Leser überraschend; die Ticket-Richtung verlangt ausdrücklich die Variable. JAIRA_SPAWN_PAUSE passt zu JAIRA_NO_WORKTREE/JAIRA_BRANCH_PREFIX im selben Kopf. Eine Variable für beide Pausen statt zwei: Tests brauchen nur 0.
