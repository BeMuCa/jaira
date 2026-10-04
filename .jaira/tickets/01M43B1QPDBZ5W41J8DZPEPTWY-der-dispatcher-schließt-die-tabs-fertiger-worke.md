---
id: 01M43B1QPDBZ5W41J8DZPEPTWY
title: Der Dispatcher schließt die Tabs fertiger Worker nie
status: human
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
updated-at: 2026-10-04T11:38:21Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-12187
claimed-at: 2026-10-04T11:33:53Z
outcome-what: "Getestet"
outcome-why: Alle DoD-Punkte belegt
outcome-resolves: "testing"
review-summary: |-
  Runde 1: leerer Status (jaira show scheitert) wurde als 'Lane verlassen' gelesen und hätte den Tab eines arbeitenden Workers geschlossen. Zurück nach in-progress.
  Runde 2: nichts mehr. Die Worktree-Ableitung steht doppelt (spawn.sh und run-lane.sh); bewusst, mit Kommentar, statt spawn.sh eine zweite Ausgabe zu geben.
review-gaps: "Nichts zu entfernen: alle Flags werden genutzt, keine zweite Kopie einer Hilfsfunktion außer der kommentierten Worktree-Ableitung."
test-verdict: "Bestanden. go test ./core/role/ grün (run-lane.sh mitgeliefert, mit Ausführungsbit). Stub-Läufe mit falschem herdr/jaira: done -> Exit 0 und Tab zu; blocked -> Exit 4, Tab offen; --timeout 0 wartet; gescheiterter Status-Lesezugriff -> wartet weiter, kein Tab zu. Nicht getestet: echter Herdr-Lauf."
question: "Einmal echt in Herdr laufen lassen (eine Lane mit scripts/run-lane.sh): schließt sich der Tab, wenn der Worker fertig ist? Getestet ist nur mit Attrappen."
---

# Der Dispatcher schließt die Tabs fertiger Worker nie

## Definition of Done

- [x] core/role/builtin/jaira-dispatcher/scripts/run-lane.sh existiert und wartet auf idle|done, schließt danach den Tab
  proof: core/role/builtin/jaira-dispatcher/scripts/run-lane.sh:95 (wartet auf idle|done), core/role/builtin/jaira-dispatcher/scripts/run-lane.sh:108 tab close; Stub-Lauf: Exit 0, Tab geschlossen
- [x] run-lane.sh nimmt --no-worktree nur, wenn es übergeben wird, und ruft spawn.sh neben sich selbst auf
  proof: core/role/builtin/jaira-dispatcher/scripts/run-lane.sh:60 (--no-worktree nur auf Flag), core/role/builtin/jaira-dispatcher/scripts/run-lane.sh:46 spawn.sh über dirname $0
- [x] Keine Warteschleife ohne Zeitgrenze; ein Approval-Dialog (blocked) bricht mit Meldung ab, ohne den Tab zu schließen
  proof: core/role/builtin/jaira-dispatcher/scripts/run-lane.sh:79 check(): Deadline Exit 3 (core/role/builtin/jaira-dispatcher/scripts/run-lane.sh:82), blocked Exit 4 ohne tab close (core/role/builtin/jaira-dispatcher/scripts/run-lane.sh:88); Stub-Lauf blocked: Exit 4
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
