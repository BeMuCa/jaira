---
id: 01M7Z6T8TS8CRXA4FH218R8VTS
title: "Die Rollen-Tests in core/role warten nicht mehr auf Sekunden-Pausen, die für ein echtes Terminal gedacht sind"
status: done
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
commits:
  - ca2928da558275d99b28583aced79c8da1538bd3
  - 6c243af0da564c85f106ce75d38bf2ea85caf16a
created-at: 2026-11-21T14:05:27Z
updated-at: 2026-11-21T18:37:33Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-842456
claimed-at: 2026-11-21T14:05:56Z
outcome-what: spawn.sh liest JAIRA_SPAWN_PAUSE für seine zwei Wartezeiten (Default 2 s/1 s unverändert); die core/role-Tests setzen sie in TestMain auf 0
outcome-why: "jeder spawn.sh-Aufruf im Test schlief 3 s auf ein claude, das es nicht gibt; newReapRig ruft spawn.sh zweimal, daher die 6,5 s pro TestReap*"
outcome-resolves: "go test ./core/role 89,5 s -> 4,4 s lokal"
review-summary: "spawn.sh liest seine zwei Wartezeiten (2 s zwischen den bis zu 20 Blicken auf den Tab, 1 s vor dem Enter) jetzt aus JAIRA_SPAWN_PAUSE; ungesetzt oder leer bleiben es 2 und 1. TestMain in core/role setzt die Variable für das ganze Paket auf 0, und weil jeder Skripttest os.Environ() bzw. r.env weiterreicht, schläft kein spawn.sh-Aufruf im Test mehr - newReapRig rief spawn.sh zweimal (2x3 s), daher die 6,5 s je TestReap*. Ein neuer Test legt einen sleep-Stub auf PATH, leert die Variable und prüft, dass spawn.sh genau 2 und 1 Sekunde verlangt. shipped.sha256 bekommt den neuen spawn.sh-Hash. Gemessen: go test -count=1 ./core/role 4,4 s (vorher laut Note 89,5 s)."
review-gaps: |-
  note: JAIRA_SPAWN_PAUSE=abc (kein Zahlwert) lässt sleep unter set -e scheitern, nachdem der Tab schon offen ist - wie andere Fehlerpfade dort; Fix wäre ein Fallback auf den Default bei nicht-numerischem Wert, nur nötig falls die Variable je für Nutzer gedacht ist
  note: keine Release-Note für JAIRA_SPAWN_PAUSE, bewusst (Test-Hebel, Default unverändert) - vertretbar, weil sich ohne die Variable nichts ändert; die Variable steht aber in der Usage von spawn.sh, also sichtbar
  note: die TestReap*-Tests skippen auf Windows (reap_test.go:35), der Gewinn dort kommt nur aus den role_test.go-Spawn-Tests; CI-Zeit vorher/nachher je Shard ist nicht gemessen
  note: pre-existing: spawn.sh:186-190 Port-Offset cksum(slug) mod 40 + 1 kollidiert häufig (8r8vts/32ttss/twb0r6 alle Offset 5) - aus dem Bughunt, nicht dieses Ticket
test-verdict: "pass: go build ./... rc=0, go test -race -count=1 ./... rc=0 (core/role 6,3 s); DoD 1-3 im Baum geprüft (Zeiten, spawn.sh:57-60/357/417, TestSpawnPausesAsLongAsBeforeWithoutTheVariable, newReapRig:126-127); spawn.sh auf Stub: default 3,0 s, =0 0,03 s; nie-idle-Pfad rc=1 mit 40 s default; schreibt nichts Neues"
hunt-verdict: "pass: bughunt done at 6c243af0, 0 change, 1 pre-existing (Port-Offset-Kollision in spawn.sh:186-190, Note für die Person)"
review-verdict: "Der Diff erfüllt alle drei DoD-Punkte: Paket läuft in 4,2 s (selbst nachgemessen), echte Läufe warten wie vorher (Default-Expansion + Stub-Test), Ursache benannt und behoben. Kein Blocker; shipped-history --check und core/shipped-Tests grün."
review-round: 1
review-check: |-
  1. cd /home/alex/projects/.worktrees/jaira-8r8vts
  2. time go test -count=1 ./core/role  -> ok nach etwa 4-5 s (vorher etwa 90 s)
  3. go test -count=1 ./core/role -run TestSpawnPausesAsLongAsBeforeWithoutTheVariable -v  -> PASS: ohne Variable schläft spawn.sh 2 s und 1 s wie bisher
  4. grep -n JAIRA_SPAWN_PAUSE core/role/builtin/jaira-dispatcher/scripts/spawn.sh  -> Beschreibung im Usage-Kopf (Zeile 57) und die zwei sleep-Zeilen mit Default 2 und 1
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
- **2026-11-21 14:13 · Alexander Sacharov** — bughunt: pre-existing: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:186-190 - the port offset is cksum(slug) mod 40 + 1, so two slugs often get the same offset: their worktrees get identical HTTP_PORT/DB_PORT_HOST/DB_PORT_TEST_HOST (and .jaira/worktree-setup the same offset), and the second docker stack fails on busy ports. The comment's claim 'two workers never share a stack' is false. Triggers whenever two colliding tickets run at once in a repo with a .env or worktree-setup; blame: 8561b2cc (C9V7ZV) / 3859cf8 (167HY8), before the ticket; ran: cksum over this repo's feat/* slugs - 8r8vts, 32ttss and twb0r6 all get offset 5, and 10 offsets are shared by two or more slugs
- **2026-11-21 14:13 · Alexander Sacharov** — bughunt: done at 6c243af0 - 1 finding (0 change, 1 pre-existing) in 2 modules (spawn.sh; package core/role tests), first hunt; concurrency, error paths, input, permissions checked; files: core/role/builtin/jaira-dispatcher/scripts/spawn.sh core/role/main_test.go core/role/reap_test.go core/shipped/shipped.sha256
- **2026-11-21 14:14 · Alexander Sacharov** — testing: Funktion auf Herdr-Stub (Scratch-Repo): ungesetzt 3,04 s, leer 3,04 s, =0 0,03 s, alle rc=0 und Pane p1 gedruckt. Fehlerpfad: Stub nie idle -> ungesetzt 40,3 s (20x2 s wie vorher), =0 sofort, beide rc=1 'claude did not come up'. JAIRA_SPAWN_PAUSE=abc -> sleep-Fehler, rc=1, Tab bleibt offen (wie andere Fehlerpfade, siehe Hunt). Zwei Schreiber: die Änderung schreibt nichts Neues. Hunt-Fund pre-existing (Port-Offset-Kollision spawn.sh:186-190) bleibt als Note für die Person - anderes Thema, kein Satellit-Ticket angelegt (Regel: keine Satelliten-Tickets).
- **2026-11-21 17:50 · Alexander Sacharov** — explain: Der Review findet keinen Blocker: go test ./core/role läuft jetzt in ca. 4 s statt 90 s, weil spawn.sh seine zwei Pausen aus JAIRA_SPAWN_PAUSE liest (Default 2 s und 1 s) und die Tests die Variable auf 0 setzen. Offen sind nur kleine Kanten: ein nicht-numerischer Wert bricht spawn.sh ab, es gibt keine Release-Note, und die CI-Zeit ist nicht gemessen.
Warum: Die Variable steht in der Usage von spawn.sh und ist damit für Nutzer sichtbar, aber ohne sie ändert sich nichts am echten Verhalten. Die Kanten treffen nur jemanden, der die Variable falsch setzt.
1) Zurückgeben: bei nicht-numerischem JAIRA_SPAWN_PAUSE auf den Default zurückfallen und eine Release-Note schreiben [give back]
2) So annehmen: Variable bleibt ein Test-Hebel, die Kanten stehen als Notes am Ticket [accept]
Empfehlung: 2 - das Ziel (Tests ~10 s statt 88 s) ist erreicht, die Kanten sind klein und betreffen nur Fehlbedienung einer Test-Variable.
