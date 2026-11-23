---
id: 01M2HRF34AS7QDTACC323YRPXJ
title: "spawn.sh kann nur eine Rolle starten, obwohl der Prompt zwei verlangt"
status: in-progress
ready: true
creator: Alexander Sacharov
goal: "Ein Dispatcher kann mit dem mitgelieferten Skript jede Lane besetzen, die sein eigener Prompt von ihm verlangt - nicht nur die, deren Name fest eingebaut ist."
context: "scripts/spawn.sh:71 schickt einen fest verdrahteten Befehl in die Vorlage:\n\n  \"$herdr\" pane send-text \"$pane\" \"/jaira-role-lane $ticket $lane\"\n\nDer Prompt des Dispatchers verlangt aber zwei verschiedene Rollen. jaira-dispatcher/SKILL.md:48 sagt woertlich: '/jaira-role-lane <id> <lane>. Testing is not a lane: /jaira-role-tester <id>.'\n\nDamit kann ein Dispatcher, der seinem eigenen Prompt folgt UND das mitgelieferte Skript benutzt, nicht beides tun. Genau diese Kombination verlangt KSGSKK von ihm.\n\nGefunden am 2026-09-15 von dem Dispatcher, der 13VMA8 fuehrte, beim Starten der testing-Lane - also von dem ersten, dessen Prompt ihn ueberhaupt auf spawn.sh verwies. Er hat mit lane=testing weitergemacht, was auf diesem Board vertretbar ist, weil testing hier eine eingerichtete Lane mit eigenem Prompt und einem verpflichtenden test-verdict ist. Auf einem Board ohne diese Lane haette er nichts starten koennen.\n\nZweiter Befund aus demselben Lauf, gleiche Datei: spawn.sh vertraut der Pane-Id, die 'tab create' zurueckgibt. Sie wurde stale, als ein zweiter Dispatcher im selben Moment seinen eigenen Tab aufmachte. Der betroffene Lauf hat die Pane ueber die Tab-Id neu aufgeloest und kam durch; das Skript kann das nicht.\n\nNicht Teil dieses Tickets: der WSL-Start (wsl.exe --cd), das ist in KSGSKK behoben und seit dem 2026-09-15 auch in der installierten Kopie."
definition-of-done: "scripts/spawn.sh startet die Rolle, die zur Lane gehoert: fuer testing '/jaira-role-tester <id>', sonst '/jaira-role-lane <id> <lane>'. Nachgestellt, indem beide Lanes einmal wirklich gestartet werden."
tags:
  - cli
blocked-by: []
related:
  - 01M2G9X5HVH29SDS8FAZKSGSKK
commits: []
created-at: 2026-09-15T05:26:04Z
updated-at: 2026-11-23T18:08:43Z
assignee: "Alexander Sacharov"
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-16863
claimed-at: 2026-11-23T17:56:21Z
outcome-what: "spawn.sh nimmt tab_id und pane_id aus 'tab create' und prueft per 'pane get' vor 'pane run' und vor 'send-text', dass die Pane noch im eigenen Tab steht; sonst tippt es nichts, schliesst den eigenen leeren Tab und nennt Pane, Tab und Label. Die come-up-Fehlerzeile nennt jetzt Pane, Tab, Label und Zustand."
outcome-why: Bei zwei gleichzeitigen Tab-Starts landete die Rolle eines Dispatchers in der Pane des anderen; das Skript glaubte der Pane-Id von vorhin.
outcome-resolves: "DoD 2, 3, 4, 5 (DoD 1 erledigt/superseded)"
review-summary: "spawn.sh keeps the tab id from 'tab create' beside the pane id and, through a new pane_is_ours, asks Herdr ('pane get' -> result.pane.tab_id) whether the pane still sits in that tab - once before the pane record and 'pane run', once more right before send-text. A pane found in another tab stops the spawn: nothing typed, its own tab closed, a message naming pane, tab and <ticket>/<lane>. The 'claude did not come up' line now also names tab and label. Two Go tests with a fake herdr reproduce two simultaneous starts and the come-up failure; release note and shipped hash added. Full review (no earlier review note), at b17384e62."
round: 3
back-from: review
back-by: lane
back-count: 1
review-gaps: |-
  note: edge: fix-now: spawn.sh pane_is_ours - a failed 'herdr pane get' (exit !=0, no tab_id) gives pt='' and is reported as a tab race; at the late check before send-text it closes our own tab with claude already running and leaves the pane record behind (bughunt 18:04 reproduced it). Fix: on empty pt say 'could not check pane X', exit 1 without tab close; the race message and tab close only when pt names another tab
  note: limit: when Herdr names no tab on create the check is skipped, and a swapped create answer that carries the other tab_id too passes the check - only a fake herdr was available, no live Herdr
  note: elsewhere: worktree still holds uncommitted 'jaira update' output (.jaira/lanes/*.md, AGENTS.md, CLAUDE.md) and the branch sits on old master 5da1f4471 - teamlead note 17:56 asked for 'git merge upstream/master'; keep those files out of this ticket's commits
test-verdict: "pass: go build/vet ./... and go test -race -count=1 ./... green (RC=0); DoD 2-5 checked in the tree - TestSpawnTypesOnlyIntoItsOwnTab and TestSpawnNamesPaneAndTabWhenClaudeDoesNotComeUp pass and both fail on the old spawn.sh; normal case run through herdr-hc.sh with a stub hc: exit 0, claude started and role typed into t9, one record t9; failure path and two writers at once are the race test (exit 1, nothing typed into the foreign pane, one record); real Herdr not available on this machine (no herdr server), so the race is reproduced with a stub only"
hunt-verdict: "pass: bughunt done at b17384e62, 1 change (edge, not a blocker: a failed 'pane get' is read as a race and closes the own tab; note for review), 0 pre-existing"
review-verdict: "accept: DoD 2-5 met in code and tests (go test ./role -run TestSpawn passes); one fix-now edge on a failed pane get"
review-round: 2
review-verdict-gaps: |-
  note: edge: fix-now: spawn.sh pane_is_ours - a failed 'herdr pane get' (exit !=0, no tab_id) gives pt='' and is reported as a tab race; at the late check before send-text it closes our own tab with claude already running and leaves the pane record behind (bughunt 18:04 reproduced it). Fix: on empty pt say 'could not check pane X', exit 1 without tab close; the race message and tab close only when pt names another tab
  note: limit: when Herdr names no tab on create the check is skipped, and a swapped create answer that carries the other tab_id too passes the check - only a fake herdr was available, no live Herdr
  note: elsewhere: worktree still holds uncommitted 'jaira update' output (.jaira/lanes/*.md, AGENTS.md, CLAUDE.md) and the branch sits on old master 5da1f4471 - teamlead note 17:56 asked for 'git merge upstream/master'; keep those files out of this ticket's commits
review-check: "1. cd core && go test ./role -run TestSpawn -count=1 -v  2. see TestSpawnTypesOnlyIntoItsOwnTab and TestSpawnNamesPaneAndTabWhenClaudeDoesNotComeUp end in PASS  3. git diff master..HEAD -- core/role/builtin/jaira-dispatcher/scripts/spawn.sh: pane_is_ours is called twice, before the pane record and right before send-text  4. the race itself cannot be forced by hand in a live Herdr; the fake-herdr test is the proof"
---

# spawn.sh kann nur eine Rolle starten, obwohl der Prompt zwei verlangt

## Definition of Done

- [-] scripts/spawn.sh startet die Rolle, die zur Lane gehoert: fuer testing '/jaira-role-tester <id>', sonst '/jaira-role-lane <id> <lane>'. Nachgestellt, indem beide Lanes einmal wirklich gestartet werden.
- [x] Die Pane, in die getippt wird, wird nicht aus einem Wert von vorhin geglaubt: ein zweiter Dispatcher, der im selben Moment einen Tab aufmacht, laesst den Lauf nicht mehr in die falsche Vorlage schreiben. Nachgestellt mit zwei gleichzeitigen Starts.
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh pane_is_ours (vor dem Pane-Record, vor pane run und vor send-text); Test core/role/spawnpane_test.go TestSpawnTypesOnlyIntoItsOwnTab (zwei gleichzeitige Starts, Fake-herdr gibt dem ersten die Pane des zweiten; nur p2 wird gestartet/betippt, nur ein Record; schlaegt auf dem alten spawn.sh fehl)
- [x] Faellt der Start doch um, sagt das Skript in welcher Vorlage und woran - statt einer Zeile, die nur 'claude did not come up' meldet.
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh pane_is_ours (Race-Meldung; 'could not check pane' wenn herdr nicht antwortet) + come-up-Zeile; Tests TestSpawnNamesPaneAndTabWhenClaudeDoesNotComeUp, TestSpawnKeepsItsTabWhenHerdrDoesNotAnswer
- [x] Release-Note in core/release/notes/3YRPXJ.md, weil das Skript mit den Rollen ausgeliefert wird.
  proof: core/release/notes/3YRPXJ.md
- [x] Fehlerzeile 'claude did not come up in <pane>: <state>' gibt es schon; pruefen, dass sie bei der Race (falsche Pane) die richtige Pane und den Tab nennt, nicht die stale Pane.
  proof: Bei der Race bricht spawn.sh schon vor 'pane run' ab und nennt die fremde Pane UND den eigenen Tab (TestSpawnTypesOnlyIntoItsOwnTab prueft 'pane p2 is not in tab t1'); die come-up-Zeile nennt Pane, Tab und Label (TestSpawnNamesPaneAndTabWhenClaudeDoesNotComeUp)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-16 06:48 · Alexander Sacharov** — Am 2026-09-15 und 16 dreimal live bestaetigt, von der Teamlead-Sitzung selbst - der Befund ist keine Lesefrucht mehr.

Ein Teamlead soll laut ~/.claude/skills/jaira-teamlead/SKILL.md den Dispatcher ueber scripts/spawn.sh starten. spawn.sh:96 tippt fest '/jaira-role-lane $ticket $lane'. Es gibt keinen Weg, damit '/jaira-dispatcher <id>' zu starten - also genau die Rolle, fuer die das Skript im Dispatcher-Skill liegt.

Was die Sitzung stattdessen getan hat, dreimal (9ZZSFT, 0YGWXQ zweimal): spawn.sh in ein Scratchpad kopiert und die eine Zeile per sed durch ${JAIRA_SPAWN_PROMPT:-...} ersetzt. Das lief, hinterlaesst aber eine zweite Kopie des Skripts, die niemand pflegt - dieselbe Krankheit, die 7MG5GB an einem fremden Board beschreibt.

Damit ist der Befund breiter als der Titel sagt. Es fehlt nicht nur '/jaira-role-tester': es fehlt jede Rolle ausser einer. Drei Aufrufer wollen drei verschiedene Zeilen - der Dispatcher '/jaira-role-lane' und '/jaira-role-tester' (SKILL.md:48), der Teamlead '/jaira-dispatcher'. Eine Loesung, die nur testing nachtraegt, laesst den Teamlead wieder mit sed arbeiten.

Vorschlag fuer die Plan-Lane, nicht entschieden: den Prompt zum Argument machen statt ihn aus der Lane abzuleiten. Ein Aufrufer, der nichts sagt, bekommt weiter '/jaira-role-lane <id> <lane>'; wer etwas anderes braucht, uebergibt es. Das deckt alle drei Aufrufer mit einer Zeile und braucht keine Liste von Sonderfaellen im Skript.

DoD 1 ist damit zu eng formuliert ('fuer testing ..., sonst ...'). Wer ihn woertlich baut, hat den Teamlead-Fall nicht geloest.
- **2026-11-23 17:55 · Alexander Sacharov** — decision: Alex 23.11.2026: Ticket eingeengt auf DoD 2 (Pane-Race bei gleichzeitigen Starts) + 3 + 4. DoD 1 ist durch den heutigen spawn.sh erledigt (Rollen per Lane-Name). Sofort starten.
- **2026-11-23 17:56 · Alexander Sacharov** — teamlead: Worktree steht auf altem master 5da1f4471 (vor #224). Vor dem ersten Commit 'git merge upstream/master', sonst blockiert Update-Rest in .jaira/lanes, AGENTS.md, CLAUDE.md die Gates.
- **2026-11-23 17:56 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/48662.sock
- **2026-11-23 17:57 · Alexander Sacharov** — decision (worker): keine offenen Entscheidungen. Herdr-IDs sind laut 'herdr --skill' stabil, aber eine Pane, die in einen anderen Workspace wandert, bekommt eine neue ID. spawn.sh nimmt deshalb tab_id UND pane_id aus 'tab create' und prueft per 'pane get' (result.pane.tab_id) vor 'pane run' und vor 'send-text', dass die Pane noch im eigenen Tab steht. Passt sie nicht: Abbruch mit Tab, Label und Pane in der Meldung, nichts getippt. Kein Neu-Aufloesen ueber 'pane list': dessen JSON-Form ist hier nicht pruefbar (kein laufender herdr-Server auf dieser Maschine, nur Herdr App/hc), und ein Raten an der Form waere ein neuer Fehler. Unter hc (herdr-hc.sh) ist Pane=Tab, die Pruefung laeuft immer durch. Nachstellung der Race mit Fake-herdr im Go-Test, nicht live.
- **2026-11-23 18:00 · Alexander Sacharov** — critique: Pane-Check vor den Pane-Record in spawn.sh ziehen - sonst bleibt bei der Race ein Record mit der fremden Pane stehen; Abbruch schliesst den eigenen Tab schon, also braucht er keinen Record. Sonst nichts: zweiter Check vor send-text bleibt (claude laeuft 40 s, so lange kann die Pane wandern); Muster 'pane get | result.pane.tab_id' ist das aus reap.sh/close-tab.sh.
- **2026-11-23 18:01 · Alexander Sacharov** — critique pass 2: der einzige Befund (Check vor den Pane-Record) ist in 148c249d9 behoben, Test prueft nur einen Record p2 mit der Lane des Gewinners. Nichts weiter.
- **2026-11-23 18:04 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:468 - class: any failure of 'herdr pane get' inside pane_is_ours (non-zero exit, timeout, output without result.pane.tab_id) is read as 'the pane sits in another tab': the spawn blames a tab race and closes $tab; at the second check (line 643) that tab is our own with claude already started in it by 'pane run' (line 565), and the pane record written at line 492 is left behind for a closed tab - the opposite of what the release note promises ('closes its own empty tab'); example: herdr answers the third 'pane get' with exit 1 -> message 'pane p1 is not in tab t1 ... but in no tab - another tab was opened at the same moment', tab t1 closed, record p1 left; fix: tell 'pane get failed' apart from 'pane in another tab' (exit 1 without closing on a failed get), and on the late check remove the record it wrote; ran: Herdr stub failing the Nth pane get, FAILN=1 and FAILN=3 (both exit 1, closed: t1, FAILN=3 run: p1, records: p1); not a blocker - the spawn still fails visibly and the lane can be started again
- **2026-11-23 18:04 · Alexander Sacharov** — hunt report: 1 finding (1 change, edge, not a blocker), as a note on 3YRPXJ | read whole: core/role/builtin/jaira-dispatcher/scripts/spawn.sh, core/role/spawnpane_test.go, herdr-hc.sh (as the other side of pane get/tab create); read as diff: core/release/notes/3YRPXJ.md, core/shipped/shipped.sha256 (current spawn.sh hash present), uncommitted .jaira/lanes/*.md, AGENTS.md, CLAUDE.md (regenerated prompt text, no code) | could not check against a real Herdr whether a swapped 'tab create' answer can also carry the other tab's tab_id (then the check passes and tab close would hit the other tab) - only the stub was available; go test ./role -run TestSpawn -count=3 passes
- **2026-11-23 18:04 · Alexander Sacharov** — bughunt: done at b17384e62 - 1 finding (1 change, 0 pre-existing) in 4 modules, first hunt; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/release/notes/3YRPXJ.md core/role/builtin/jaira-dispatcher/scripts/spawn.sh core/role/spawnpane_test.go core/shipped/shipped.sha256
- **2026-11-23 18:05 · Alexander Sacharov** — testing: note: edge: class: herdr 'pane get' failure gleich Race behandelt; example: dritter pane get exit 1 -> Tab mit laufendem claude geschlossen, Record bleibt; fix: get-Fehler getrennt melden, ohne tab close (Bughunt-Note 18:04)
- **2026-11-23 18:05 · Alexander Sacharov** — waiting: review lane (fresh agent via run-lane.sh) since 19:05
- **2026-11-23 18:06 · Alexander Sacharov** — review: done at b17384e62
- **2026-11-23 18:06 · Alexander Sacharov** — review: fix-now: spawn.sh pane_is_ours - empty pt (pane get failed) must not count as race: message 'could not check pane', exit 1, no tab close; tab close + race message only when pt names another tab. Add a test case with a failing pane get.
- **2026-11-23 18:08 · Alexander Sacharov** — fix-now aus review: leeres pt (herdr pane get scheitert) ist keine Race mehr - eigene Meldung, kein tab close; Test TestSpawnKeepsItsTabWhenHerdrDoesNotAnswer. merge upstream/master nicht gemacht: git merge bricht an den uncommitteten 'jaira update'-Dateien ab (.jaira/lanes/*.md gleich upstream, AGENTS.md/CLAUDE.md nur Versionsstempel 0.6.5), und das Verwerfen per git checkout wurde mir verweigert - an den Teamlead.
