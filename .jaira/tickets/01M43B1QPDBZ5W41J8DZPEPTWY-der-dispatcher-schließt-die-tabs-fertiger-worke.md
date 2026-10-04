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
updated-at: 2026-10-04T12:21:36Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-12187
claimed-at: 2026-10-04T11:33:53Z
outcome-what: "Getestet"
outcome-why: "Befunde eingearbeitet"
outcome-resolves: "testing"
review-summary: "Runde 2 (fe7efd1): run-lane.sh liest nach spawn.sh einmal den Ausgangsstatus (start, Z.71) und zählt die Lane erst als fertig, wenn das Ticket in der Lane war und sie verlassen hat, oder wenn es in einen Status wechselt, der weder start noch die Lane ist (Worker hat die Lane übersprungen). Leerer Status zählt nie als Ausgang. Danach wie gehabt: warten auf idle|done, Ausgabe, Tab schließen. SKILL.md sagt jetzt, dass der Start vor dem Move in die Lane in Ordnung ist. Der Runde-1-Defekt (Tab schließt beim ersten idle, weil das Ticket noch nicht in der Lane war) ist im Normalfall weg: Stub-Lauf todo->pre-process->in-progress Exit 0 erst nach dem Wechsel; Lane übersprungen Exit 0; critique zurück nach in-progress Exit 0."
review-gaps: "1) Defekt, run-lane.sh:71: start=\"$(status)\" wird nicht wiederholt, wenn das Lesen scheitert. Dann ist start leer, und Z.105 ([ \"$s\" != \"$start\" ]) wertet schon den ersten gelesenen Vorstatus (z.B. todo) als 'Worker hat die Lane übersprungen' -> Schleife bricht ab, agent() ist direkt nach dem Start idle, Tab wird geschlossen. Genau der Runde-1-Fehler auf anderem Weg. Stub belegt: Statusfolge FAIL,todo,todo -> Exit 0, tab close nach 3 Lesezugriffen. Wahrscheinlich gerade dort, wo es zählt: der Dispatcher schreibt das Ticket (claim/move) genau in dem Moment, in dem run-lane.sh startet. Fix: start lesen, bis es nicht leer ist (mit check dazwischen). 2) Klein, kein Rückgabegrund: Kehrt ein Worker ohne Statuswechsel zurück (z.B. critique auf einem Ticket in in-progress, das nie nach critique kam), wartet das Skript die vollen 240 min und meldet dann 'timed out ... left running', obwohl der Worker idle ist — bewusste Entscheidung laut Note, aber die Meldung von Exit 3 führt in die Irre. 3) Klein, vorbestehend: Z.121 herdr pane get ohne Schutz — ist der Tab schon von Hand zu, endet das Skript mit Python-Fehler und Exit 1, nicht dokumentiert."
test-verdict: |-
  Bestanden. go test ./core/role/ grün (run-lane.sh mitgeliefert, mit Ausführungsbit). Stub-Läufe mit falschem herdr/jaira: done -> Exit 0 und Tab zu; blocked -> Exit 4, Tab offen; --timeout 0 wartet; gescheiterter Status-Lesezugriff -> wartet weiter, kein Tab zu. Nicht getestet: echter Herdr-Lauf.
  Runde 2: go test ./core/role/ grün; sechs Stub-Szenarien wie in der Notiz.
  Runde 3: go test ./core/role/ grün; acht Stub-Szenarien plus Tab-weg-Fall.
question: "Einmal echt in Herdr laufen lassen (eine Lane mit scripts/run-lane.sh): schließt sich der Tab, wenn der Worker fertig ist? Getestet ist nur mit Attrappen."
review-verdict: "Zurück nach in-progress. Der Runde-1-Defekt ist im Normalfall behoben, aber ein gescheiterter erster Statuslesezugriff (run-lane.sh:71) bringt ihn zurück: Tab eines gerade gestarteten Workers wird geschlossen. Einzeiliger Fix, danach sollte es passen. go test ./core/role/ grün, bash -n ok."
review-check: "1. go test ./core/role/ -run 'TestDispatcherShipsItsScript|TestInstallWritesEveryFile' -> ok. 2. Stub-Verzeichnis anlegen mit falschem herdr (pane get liefert agent_status idle, tab close wird geloggt), falschem jaira (gibt Status aus einer Liste zurück, Zeile FAIL = Exit 1) und sleep, das sofort endet. 3. export HERDR_BIN_PATH=<stub>/herdr HERDR_ENV=1 PATH=<stub>:$PATH. 4. Liste todo,todo,pre-process,pre-process,in-progress; bash core/role/builtin/jaira-dispatcher/scripts/run-lane.sh --no-worktree TICK1 pre-process <leeres Verzeichnis mit .jaira/> -> Exit 0, tab close erst nach dem Wechsel auf in-progress. 5. Liste FAIL,todo,todo,todo -> heute: Exit 0 und tab close, obwohl das Ticket nie in pre-process war (das ist der Defekt). Nach dem Fix muss das Skript hier weiter warten. 6. Echt in Herdr einmal eine Lane mit scripts/run-lane.sh laufen lassen: Tab schließt sich erst, wenn der Worker fertig ist."
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
