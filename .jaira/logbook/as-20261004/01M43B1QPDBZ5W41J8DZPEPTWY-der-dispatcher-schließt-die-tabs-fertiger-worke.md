---
id: 01M43B1QPDBZ5W41J8DZPEPTWY
title: Der Dispatcher schließt die Tabs fertiger Worker nie
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "run-lane.sh liegt mit dem Fix im Binary: ein Worker, der fertig ist, wird erkannt und sein Tab geschlossen — auf jedem Rechner, nicht nur auf Alex'."
context: |-
  Was falsch ist: Ein fertiger Worker meldet in Herdr agent_status 'done', nicht 'idle'. Eine Warteschleife, die nur idle|blocked prüft, kehrt nie zurück; der Tab bleibt offen, das Board steht die Nacht über still.
  Alex hat das am 2026-10-03 lokal repariert: ~/.claude/skills/jaira-dispatcher/scripts/run-lane.sh (neu) und ein Absatz in SKILL.md.
  Im Repo fehlt beides: core/role/builtin/jaira-dispatcher/ hat nur spawn.sh. Wer die Rollen aus dem Binary installiert, bekommt den Fehler zurück.
  Weitere Fehler im lokalen run-lane.sh, gleich mit zu beheben:
  - --no-worktree fest verdrahtet; SKILL.md erlaubt das nur in drei Fällen.
  - Pfad ~/.claude/skills/... fest verdrahtet; Rollen werden auch nach <repo>/.claude/skills installiert (core/role/target.go).
  - Schickt nach spawn.sh einen zweiten Prompt; spawn.sh tippt den Lane-Befehl aber schon selbst (spawn.sh: send-text).
  - Wartet auf den Text 'auto mode' in der Pane und hat keine Zeitgrenze; ohne Auto-Modus oder bei totem Worker hängt es ewig.
  - Liest den Status aus der Ticketdatei im Repo-Root; im Worktree-Modus schreibt der Worker aber die Kopie im Worktree.
definition-of-done: "core/role/builtin/jaira-dispatcher/scripts/run-lane.sh existiert und wartet auf idle|done, schließt danach den Tab"
tags:
  - release
blocked-by: []
related: []
commits:
  - 418a9c4fd654f0c90600584b65eeb4b2d2927984
  - fe7efd174071a2e0f0a66e84c6a0538c6b8ba738
  - 86f473ebcdb53be5d53e1320acf630598f5e0e66
  - b1dba69f7c7124586696665c3f8c3cd43c1c4e3c
created-at: 2026-10-04T11:33:37Z
updated-at: 2026-10-04T15:22:51Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-12187
claimed-at: 2026-10-04T11:33:53Z
outcome-what: "Von Alex abgenommen"
outcome-why: "Alex: 'review machen und schließen'"
outcome-resolves: "signoff"
review-summary: "Liefert core/role/builtin/jaira-dispatcher/scripts/run-lane.sh mit dem Binary aus (Test prüft Ausführungsbit). Das Skript startet den Worker über spawn.sh neben sich selbst (--no-worktree nur auf Flag/Env), liest den Ticketstatus dort, wo der Worker schreibt (Worktree, sonst Repo), wartet bis das Ticket in der Lane war und sie verlassen hat, dann bis Herdr idle|done meldet, und schließt den Tab. Zeitgrenze (Default 240 min, Exit 3), Approval-Dialog (blocked) bricht mit Exit 4 ab, Tab bleibt offen. Runde 3 (86f473e): Startstatus wird so lange neu gelesen, bis er nicht leer ist (run-lane.sh:95); Timeout-Meldung nennt den Worker-Zustand statt 'left running' (run-lane.sh:81); ein schon geschlossener Tab beim Abschluss wird gemeldet statt abzustürzen (run-lane.sh:126-133)."
review-gaps: "Keine blockierenden. Die drei Runde-2-Lücken sind zu (Stub-Läufe: zwei gescheiterte Lesezugriffe, dann review→human: Exit 0, ein tab close; nur todo nach gescheitertem Lesen: kein tab close). Beobachtung, kein Fehler: schließt der Mensch den Tab, bevor der Worker idle|done gemeldet hat, liefert agent() leer und die Schleife in run-lane.sh:113 wartet bis zur Zeitgrenze (Exit 3, Meldung 'tab left open', obwohl er weg ist); der Zweig 'already gone' in run-lane.sh:131-132 greift nur im kurzen Fenster zwischen Zeile 113 und 126. Kein Absturz, kein falsch geschlossener Tab, Rückkehr per Timeout — daher nicht zurückgeschickt."
test-verdict: |-
  Bestanden. go test ./core/role/ grün (run-lane.sh mitgeliefert, mit Ausführungsbit). Stub-Läufe mit falschem herdr/jaira: done -> Exit 0 und Tab zu; blocked -> Exit 4, Tab offen; --timeout 0 wartet; gescheiterter Status-Lesezugriff -> wartet weiter, kein Tab zu. Nicht getestet: echter Herdr-Lauf.
  Runde 2: go test ./core/role/ grün; sechs Stub-Szenarien wie in der Notiz.
  Runde 3: go test ./core/role/ grün; acht Stub-Szenarien plus Tab-weg-Fall.
question: "Einmal echt in Herdr laufen lassen (eine Lane mit scripts/run-lane.sh): schließt sich der Tab, wenn der Worker fertig ist? Getestet ist nur mit Attrappen."
review-verdict: "Bestanden. Der Diff erfüllt alle fünf DoD-Punkte; die Runde-2-Lücken sind behoben und per Stub geprüft. Unsicher bleibt nur der echte Herdr-Lauf (agent_status-Werte, tab_id-Feld), der nur mit Attrappen getestet ist."
review-check: "1. cd /home/alex/projects/jaira && go test ./core/role/ — erwartet: ok.  2. bash core/role/builtin/jaira-dispatcher/scripts/run-lane.sh --help — erwartet: Hilfetext mit Exit-Codes 0/3/4.  3. In einer Herdr-Pane in einem Repo mit .jaira: run-lane.sh --no-worktree <ticket-id> <lane> im Hintergrund starten.  4. Es öffnet sich ein Tab '<ticket>/<lane>', claude tippt /jaira-role-lane selbst.  5. Wenn der Worker das Ticket aus der Lane bewegt hat und fertig ist: Skript gibt 'status: <neue Lane>' und 'closed <pane>' aus, Exit 0, der Tab ist zu.  6. Gegenprobe: denselben Lauf mit --timeout 1 starten und den Worker an einem Approval-Dialog stehen lassen — Exit 4, Tab bleibt offen."
---

# Der Dispatcher schließt die Tabs fertiger Worker nie

## Definition of Done

- [x] core/role/builtin/jaira-dispatcher/scripts/run-lane.sh existiert und wartet auf idle|done, schließt danach den Tab
  proof: core/role/builtin/jaira-dispatcher/scripts/run-lane.sh:109 (idle|done), core/role/builtin/jaira-dispatcher/scripts/run-lane.sh:122 tab close; Stub-Läufe in der Notiz
- [x] run-lane.sh nimmt --no-worktree nur, wenn es übergeben wird, und ruft spawn.sh neben sich selbst auf
  proof: core/role/builtin/jaira-dispatcher/scripts/run-lane.sh:60 (--no-worktree nur auf Flag), core/role/builtin/jaira-dispatcher/scripts/run-lane.sh:46 spawn.sh über dirname $0
- [x] Keine Warteschleife ohne Zeitgrenze; ein Approval-Dialog (blocked) bricht mit Meldung ab, ohne den Tab zu schließen
  proof: core/role/builtin/jaira-dispatcher/scripts/run-lane.sh:80 check(): Exit 3 Deadline, Exit 4 blocked ohne tab close
- [x] SKILL.md im Repo beschreibt run-lane.sh; Test prüft, dass run-lane.sh mit Ausführungsbit installiert wird
  proof: core/role/builtin/jaira-dispatcher/SKILL.md:267; TestDispatcherShipsItsScript + Install-Test in core/role/role_test.go
- [x] Zeile in core/release/NOTES.md unter Unreleased
  proof: core/release/NOTES.md:17

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-04 11:36 · Alexander Sacharov** — Bewusst weggelassen: Erkennung 'Worker idle, Ticket aber noch in der Lane'. Im conversational mode wartet ein Worker legitim auf den Menschen (SKILL.md: Pausen sind keine Stalls); das Skript würde dann fälschlich aufgeben. Nur blocked (Approval-Dialog) und Timeout brechen ab. Den zweiten Prompt ($MSG) aus Alex' lokaler Fassung gibt es nicht mehr: spawn.sh tippt den Lane-Befehl selbst.
- **2026-10-04 11:37 · Alexander Sacharov** — Kritik-Befund behoben: leerer Status zählt jetzt als 'noch in der Lane' (core/role/builtin/jaira-dispatcher/scripts/run-lane.sh:94). Stub-Test: jaira, das immer scheitert, hält das Skript im Warten, kein tab close.
- **2026-10-04 12:15 · Alexander Sacharov** — Review: zurück. run-lane.sh:94 nimmt an, das Ticket steht beim Start schon in <lane>; laut jaira-role-lane/SKILL.md:110-114 startet der Dispatcher den Worker aber vor dem Move (pre-process auf todo). Dann fällt die Status-Schleife sofort durch und run-lane.sh:95 schließt den Tab beim ersten idle. Stub-Lauf (Status todo, Lane pre-process, idle): Tab nach <1 s geschlossen, Exit 0. Fix: Startstatus merken, erst auf Eintritt in <lane> und dann auf Verlassen warten.
- **2026-10-04 12:17 · Alexander Sacharov** — Review-Befund behoben: run-lane.sh liest den Startstatus und gilt als fertig erst, wenn das Ticket in der Lane war und sie verlassen hat — oder direkt weiter zog, ohne sie zu betreten. Rückkehr in den Startstatus nach der Lane (critique schickt zurück) zählt. Stub-Läufe: nie betreten -> wartet; todo->pre-process->in-progress -> zu; todo->in-progress -> zu; in-progress->critique->in-progress -> zu; Lesefehler -> wartet.
- **2026-10-04 12:20 · Alexander Sacharov** — Review Runde 2: zurück. run-lane.sh:71 liest start einmal; scheitert dieser Lesezugriff (z.B. weil der Dispatcher das Ticket gerade schreibt), ist start leer und Z.105 nimmt den ersten Vorstatus (todo) als 'Lane übersprungen' -> Tab des gerade gestarteten Workers wird geschlossen (Stub: FAIL,todo,todo -> Exit 0 + tab close). Fix: start wiederholt lesen, bis nicht leer, check dazwischen.
- **2026-10-04 12:21 · Alexander Sacharov** — Review Runde 2 behoben: Startstatus wird wiederholt gelesen, bis er nicht leer ist (core/role/builtin/jaira-dispatcher/scripts/run-lane.sh:95); Timeout-Meldung nennt den Worker-Zustand statt 'left running'; ein schon geschlossener Tab ist kein Fehler mehr ('already gone', Exit 0). Stub: FAIL,todo,todo -> wartet; FAIL,FAIL,todo,pre-process,in-progress -> zu; Tab weg -> Exit 0.
- **2026-10-04 13:59 · Alexander Sacharov** — Review Runde 3 bestanden. Offen, nicht blockierend: schließt die Person den Tab, bevor der Worker idle|done meldet, wartet run-lane.sh bis zum Timeout (240 min) und meldet dann 'tab left open'. Bewusst gelassen — der Worker ist dann ohnehin weg, und Exit 3 weckt den Dispatcher.
