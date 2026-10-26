---
id: 01M5WGN0K4DS6QD6RKBB167HY8
title: "spawn.sh: eigenes .jaira/worktree-setup des Repos aufrufen, wenn vorhanden"
status: done
ready: true
creator: Alexander Sacharov
goal: "Ein Worker-worktree bekommt genau die Umgebung, die das Repo für worktrees vorsieht: hat das Repo ein ausführbares .jaira/worktree-setup, ruft spawn.sh es auf (worktree, slug, offset, root) statt nur drei Ports selbst anzuhängen."
context: |-
  Was falsch ist: spawn.sh kopiert .env und hängt nur COMPOSE_PROJECT_NAME, HTTP_PORT, DB_PORT_HOST, DB_PORT_TEST_HOST an. Repos brauchen mehr (Speicher-Ports, Vite/Backend-Ports, IMAGE_NS, COMPOSE_PROFILES leer …). Folge: Worker-Stacks kollidieren mit dem Hauptstack (z. B. Port 8333 belegt), Worker ergänzen die .env jedes Mal von Hand — in einer Woche in fast jedem Bericht.
  Ein Repo kann das schon selbst beschreiben: ein Skript .jaira/worktree-setup <worktree> <slug> <n> <root>, das die .env kopiert und alles Nötige anhängt (die Repo-Doku verweist darauf, als würde spawn.sh es aufrufen — tut es aber nicht).
  Ziel: liegt ein ausführbares .jaira/worktree-setup im Repo, ruft spawn.sh es nach 'git worktree add' auf (offset wie heute berechnet); sonst wie heute.
definition-of-done: Mit .jaira/worktree-setup im Repo ruft spawn.sh es auf und hängt selbst nichts an (Test)
tags:
  - cli
blocked-by: []
related: []
commits:
  - f40230edcf4a5f007d5d739570479773a7681f98
created-at: 2026-10-26T16:28:14Z
updated-at: 2026-10-26T16:41:00Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-90346
claimed-at: 2026-10-26T16:28:33Z
outcome-what: "spawn.sh ruft ein ausfuehrbares .jaira/worktree-setup des Repos mit worktree, slug, offset, root auf und schreibt dann selbst keine .env; ohne Datei wie bisher; Fehler oder fehlendes chmod +x bricht vor dem Tab ab und raeumt worktree und Branch weg"
outcome-why: Die vier fest eingebauten Port-Zeilen reichen Repos mit mehr Diensten nicht; Worker-Stacks kollidierten mit dem Hauptstack und Worker ergaenzten die .env von Hand
outcome-resolves: "Worker-worktrees bekommen die Umgebung, die das Repo selbst beschreibt; kein Hand-Nachtragen der .env mehr"
review-summary: "spawn.sh prüft beim Anlegen eines neuen worktrees, ob im Repo ein .jaira/worktree-setup liegt. Ist es ausführbar, wird es nach 'git worktree add' als '<worktree> <slug> <offset 1-40> <root>' aufgerufen (Ausgabe nach stderr, stdout bleibt die Pane-ID), und spawn.sh schreibt selbst keine .env. Scheitert es, entfernt spawn.sh worktree und Branch wieder und bricht vor 'tab create' ab ('failed (exit N)'). Liegt die Datei da, ist aber nicht ausführbar, bricht spawn.sh schon vor 'git worktree add' mit 'chmod +x'-Hinweis ab. Ohne Datei unverändert: .env kopieren + COMPOSE_PROJECT_NAME/HTTP_PORT/DB_PORT_HOST/DB_PORT_TEST_HOST. Der Offset wird jetzt immer berechnet. Neuer Test TestSpawnLetsTheRepositorySetUpItsWorktree deckt alle vier Fälle, Release-Note core/release/notes/167HY8.md."
review-gaps: "Keine blockierenden. Klein: (1) Was ein gescheitertes Setup außerhalb des worktrees angerichtet hat (z. B. schon gestartete Container), räumt spawn.sh nicht weg – das bleibt Sache des Skripts, steht aber nirgends. (2) Der Ohne-Datei-Test prüft die .env nur auf enthaltene Zeilen, nicht auf exakte Gleichheit mit dem alten Ergebnis; der unveränderte Code-Zweig macht das unbedenklich. (3) Im Dispatcher-worktree (wt=root) und mit --no-worktree läuft das Setup nicht – wie bisher bei der .env, so gewollt laut Note. Geprüft: 'git worktree remove --force' löscht auch gitignorierte Dateien wie eine vom Setup geschriebene .env (von Hand nachgestellt)."
test-verdict: "pass: go build ./... + go test -race -count=1 ./... RC=0; DoD 1-3 per TestSpawnLetsTheRepositorySetUpItsWorktree (schlaegt gegen altes spawn.sh fehl); Handlauf an Scratch-Repo mit herdr-Stub: Setup exit 7 -> kein TAB-CREATED, worktree und feat/abc1 entfernt, exit 1 mit Meldung; ohne chmod +x -> exit 1 vor worktree add"
review-verdict: "Diff erfüllt alle drei DoD-Punkte; go test ./core/role grün, go vet sauber. Keine Defekte gefunden. Annahme."
review-round: 1
review-check: "1. go test ./core/role -run TestSpawnLetsTheRepositorySetUpItsWorktree -v  → PASS mit den Subtests setup, failing_setup, setup_not_executable  2. Echt: im Repo eine Datei .jaira/worktree-setup anlegen mit Inhalt '#!/usr/bin/env bash' + 'echo \"$@\" >&2; cp \"$4/.env\" \"$1/.env\"' und chmod +x  3. In Herdr einen Arbeiter starten (jaira-dispatcher spawn.sh <slug> <ID> review)  4. Im Pane-Log stehen worktree-Pfad, Slug, eine Zahl 1-40 und der Repo-Pfad; die .env im neuen worktree hat KEINE '# worker stack'-Zeilen  5. chmod -x auf die Datei, erneut starten mit anderem Slug → Abbruch mit 'chmod +x', kein Tab, kein neuer worktree  6. Datei löschen → .env bekommt wieder die vier Port-Zeilen"
---

# spawn.sh: eigenes .jaira/worktree-setup des Repos aufrufen, wenn vorhanden

## Definition of Done

- [x] Mit .jaira/worktree-setup im Repo ruft spawn.sh es auf und hängt selbst nichts an (Test)
  proof: TestSpawnLetsTheRepositorySetUpItsWorktree/setup (core/role/role_test.go:806): Argumente wt/slug/offset/root, .env bleibt genau wie vom Skript geschrieben; core/role/builtin/jaira-dispatcher/scripts/spawn.sh:129
- [x] Ohne die Datei verhält sich spawn.sh wie bisher (Test)
  proof: TestSpawnLetsTheRepositorySetUpItsWorktree (core/role/role_test.go:792): ohne Skript .env-Kopie + COMPOSE_PROJECT_NAME/HTTP_PORT/DB_PORT_HOST/DB_PORT_TEST_HOST
- [x] Fehler im Setup-Skript bricht den Spawn mit klarer Meldung ab, ohne halben Tab
  proof: TestSpawnLetsTheRepositorySetUpItsWorktree/failing_setup (core/role/role_test.go:823) und /setup_not_executable: Exit != 0, kein 'tab create', worktree und Branch weg, Meldung 'failed (exit 3)' bzw. 'chmod +x'; core/role/builtin/jaira-dispatcher/scripts/spawn.sh:108,core/role/builtin/jaira-dispatcher/scripts/spawn.sh:129

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-26 16:28 · Alexander Sacharov** — parent session: requirementsgenie-22
- **2026-10-26 16:29 · Alexander Sacharov** — Offene Entscheidungen gezaehlt, keine braucht die Person - alle aus Ticket und Code entschieden:
1) Fehler im Setup-Skript: Abbruch (DoD 3 dieses Tickets), nicht 'melden und weiter' wie die Notiz an 7MG5GB vorschlug. Das Skript laeuft vor 'tab create', also entsteht kein Tab.
2) Nach dem Fehlschlag raeumt spawn.sh den eben angelegten worktree und Branch wieder ab. Grund: bliebe der worktree stehen, ueberspringt der naechste Aufruf das Setup ('if [ ! -d $wt ]') und der Arbeiter laeuft ohne .env.
3) Datei da, aber nicht ausfuehrbar: Abbruch mit 'chmod +x'-Hinweis, VOR 'git worktree add'. Still auf die alten drei Ports zurueckfallen waere genau die Verwechslung, die 7MG5GB beschreibt.
4) Ohne Datei: alter Weg unveraendert (.env kopieren + 4 Zeilen), wie DoD 2 verlangt. 7MG5GB will diesen Weg ganz streichen - das bleibt dort.
5) stdout des Setup-Skripts geht nach stderr: spawn.sh gibt auf stdout nur die Pane-ID aus, run-lane.sh liest sie.
6) --no-worktree: Skript laeuft nicht (wt existiert schon, wie heute fuer die .env).
Verhaeltnis zu 7MG5GB: dieses Ticket ist dessen erster Schritt (Hook aufrufen). Offen bleibt dort: alten .env-Weg entfernen, roles-install-Hilfe, Meldung fuer Boards ohne Skript.
- **2026-10-26 16:33 · Alexander Sacharov** — in-progress fertig. Der neue Test schlaegt gegen das alte spawn.sh fehl (geprueft: alle drei Subtests FAIL), gegen das neue gruen. Offset wird jetzt immer berechnet, auch ohne .env (vorher nur im .env-Zweig) - harmlos, nur eine Zahl. Hilfe von spawn.sh nennt den Aufruf; Release-Note core/release/notes/167HY8.md.
- **2026-10-26 16:34 · Alexander Sacharov** — critique: nichts. Geprueft: einfachere Form (Pruefung vor 'worktree add' + Aufraeumen nach Fehlschlag ist noetig, weil nur ein fehlendes Verzeichnis eingerichtet wird), Passung (gleiches Muster wie die bestehende .env-Verzweigung, if/elif statt zweitem Block), run-lane.sh:69 reicht den Exit von spawn.sh unter set -e durch. Die Entscheidungen aus der ersten Note stehen.
