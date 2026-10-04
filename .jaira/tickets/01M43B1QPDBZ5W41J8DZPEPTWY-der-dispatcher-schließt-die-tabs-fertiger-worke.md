---
id: 01M43B1QPDBZ5W41J8DZPEPTWY
title: Der Dispatcher schließt die Tabs fertiger Worker nie
status: review
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
commits: []
created-at: 2026-10-04T11:33:37Z
updated-at: 2026-10-04T12:17:34Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-12187
claimed-at: 2026-10-04T11:33:53Z
outcome-what: "Getestet"
outcome-why: Befund eingearbeitet
outcome-resolves: "testing"
review-summary: "Neues scripts/run-lane.sh neben spawn.sh im eingebetteten jaira-dispatcher: startet den Worker über spawn.sh (--no-worktree nur auf Flag/JAIRA_NO_WORKTREE), liest den Ticketstatus im Worktree (oder im Repo, wenn dort kein .jaira liegt), wartet erst bis der Status die Lane verlassen hat, dann bis Herdr idle|done meldet, druckt status/outcome/review/question und schließt den Tab. blocked -> Exit 4 ohne Schließen, --timeout (Default 240 min) -> Exit 3. SKILL.md bekommt einen Absatz dazu, role_test.go prüft beide Skripte samt Ausführungsbit, NOTES.md eine Zeile unter Unreleased. go test ./core/role/ grün."
review-gaps: "Defekt: run-lane.sh:94 setzt voraus, dass das Ticket beim Start schon in <lane> steht. jaira-role-lane/SKILL.md:110-114 sagt das Gegenteil: der Dispatcher startet jeden Worker, BEVOR er das Ticket in dessen Lane schiebt ('pre-process auf einem Ticket in todo ist der Normalfall'). Dann ist status != lane sofort, die erste Schleife fällt durch, und run-lane.sh:95 schließt den Tab beim ersten idle — direkt nach spawn.sh (Herdr-Hook hat 'working' evtl. noch nicht gemeldet) oder bei jeder Pause im conversational mode. Stub-Lauf belegt: Status todo, Lane pre-process, agent idle -> 'closed P1', Exit 0, nach 2 pane-get-Aufrufen. Genau der Fall, den die Notiz vom 11:36 ausschließen wollte. Vorschlag: Startstatus merken und warten, bis der Status weder <lane> noch der Startstatus ist (bzw. erst <lane> erreicht und dann verlassen wurde); SKILL.md-Absatz (core/role/builtin/jaira-dispatcher/SKILL.md:286) entsprechend präzisieren. Sonst nichts: set -e-Fallen (Zeile 60 &&-Liste, ), Worktree-Ableitung und Timeout geprüft, in Ordnung."
test-verdict: |-
  Bestanden. go test ./core/role/ grün (run-lane.sh mitgeliefert, mit Ausführungsbit). Stub-Läufe mit falschem herdr/jaira: done -> Exit 0 und Tab zu; blocked -> Exit 4, Tab offen; --timeout 0 wartet; gescheiterter Status-Lesezugriff -> wartet weiter, kein Tab zu. Nicht getestet: echter Herdr-Lauf.
  Runde 2: go test ./core/role/ grün; sechs Stub-Szenarien wie in der Notiz.
question: "Einmal echt in Herdr laufen lassen (eine Lane mit scripts/run-lane.sh): schließt sich der Tab, wenn der Worker fertig ist? Getestet ist nur mit Attrappen."
review-verdict: "Zurück nach in-progress: DoD formal erfüllt, aber run-lane.sh schließt im Normalablauf des Dispatchers (Worker vor dem Move gestartet) den Tab eines gerade gestarteten Workers."
review-check: "1. mkdir -p /tmp/rl/bin /tmp/rl/s /tmp/rl/repo/.jaira  2. git show 418a9c4:core/role/builtin/jaira-dispatcher/scripts/run-lane.sh > /tmp/rl/s/run-lane.sh; printf '#!/bin/sh\\necho P1\\n' > /tmp/rl/s/spawn.sh  3. Attrappe /tmp/rl/bin/jaira, die {\"status\":\"todo\"} ausgibt; Attrappe /tmp/rl/bin/herdr, die bei 'pane get' {\"result\":{\"pane\":{\"agent_status\":\"idle\",\"tab_id\":\"T1\"}}} ausgibt und jeden Aufruf in eine Logdatei schreibt; chmod +x alles  4. HERDR_BIN_PATH=/tmp/rl/bin/herdr PATH=/tmp/rl/bin:$PATH bash /tmp/rl/s/run-lane.sh --no-worktree ABC pre-process /tmp/rl/repo  5. Heute: sofort 'closed P1' und 'herdr tab close T1' im Log. Nach dem Fix: das Skript wartet weiter (Strg-C), bis die jaira-Attrappe erst pre-process und dann eine andere Lane meldet. 6. go test ./core/role/ muss grün bleiben."
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
