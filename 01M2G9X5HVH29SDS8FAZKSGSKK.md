---
id: 01M2G9X5HVH29SDS8FAZKSGSKK
title: "Der Dispatcher-Prompt nennt sein Transportmittel nicht, also erfindet jeder Lauf ein eigenes"
status: in-progress
ready: true
creator: Alexander Sacharov
goal: "Ein Dispatcher liest aus seinem eigenen Prompt, womit er einen Worker startet, und benutzt das mitgelieferte scripts/spawn.sh - statt sich einen Weg auszudenken, den der Berechtigungspruefer ablehnt."
context: |-
  Am 2026-09-14 liefen drei Dispatcher auf diesem Board. Zwei davon bekamen keinen einzigen Worker in eine Herdr-Vorlage und arbeiteten alle Lanes als Subagenten ab - Subagenten sterben mit der Sitzung, und genau dagegen gibt es die Vorlagen.

  Beide scheiterten gleich: sie riefen 'claude --permission-mode acceptEdits' selbst auf, und der Berechtigungspruefer lehnte das als 'Create Unsafe Agents' ab. Einer meldete danach sogar, Herdr sei nicht installiert - er hatte 'command -v herdr' gefragt, und auf diesem Rechner heisst die Datei herdr.exe und liegt unter /mnt/c.

  Der dritte Dispatcher kam durch und hatte einen Worker in einer eigenen Vorlage laufen.

  Das mitgelieferte Skript loest das Problem bereits, und zwar durch seinen Aufbau:
  core/role/builtin/jaira-dispatcher/scripts/spawn.sh nimmt den Pfad aus HERDR_BIN_PATH, spaltet eine Vorlage mit 'herdr pane split' und startet darin mit 'herdr pane run <pane> "cd <worktree> && claude"' ein blankes claude. Gestartet wird es von der Vorlage, nicht vom Bash-Werkzeug - der Pruefer sieht einen herdr-Aufruf und keine Agenten-Erzeugung, und ein --permission-mode braucht es nie. Danach schiebt es den Lane-Befehl mit 'pane send-text' und 'send-keys enter' hinein, nachdem der Zustands-Hook 'claude idle' gemeldet hat.

  Warum es trotzdem niemand benutzt: die Prompts sagen es nicht. jaira-teamlead/SKILL.md:44 und jaira-dispatcher/SKILL.md:83 verweisen beide auf 'herdr --skill' fuer die Mechanik; das Skript wird erst in jaira-dispatcher/SKILL.md:152 erwaehnt, und dort im Zusammenhang mit Projektnamen und Ports, nicht als der Weg, einen Worker zu starten. Wer den Prompt von oben liest, erfaehrt nie, dass es das Skript gibt.

  Zwei Dinge im Skript selbst passen nicht zu diesem Repository, und beide fallen erst auf, wenn es benutzt wird:
  - Zeile 'git worktree add "" -b "feature/"' legt Zweige mit dem Praefix feature/ an. Hier heissen sie feat/.
  - Der .env-Block schreibt COMPOSE_PROJECT_NAME=rg_<slug>, VITE_PORT_HOST und BACKEND_PORT_HOST - der Stapel von requirementsgenie, mit fest eingebautem Praefix rg_. In jaira gibt es keine .env, der Block wird uebersprungen und faellt darum niemandem auf.

  Nicht Teil dieses Tickets: die Berechtigungsregeln zu aendern. Es hat sich gezeigt, dass keine gebraucht werden, wenn der vorhandene Weg benutzt wird.
definition-of-done: "jaira-dispatcher/SKILL.md nennt scripts/spawn.sh dort, wo ein Worker gestartet wird, nicht erst im Abschnitt ueber Ports - und sagt dazu, dass ein selbst aufgerufenes 'claude --permission-mode ...' vom Berechtigungspruefer abgelehnt wird, damit niemand es noch einmal versucht."
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-14T15:52:22Z
updated-at: 2026-09-14T20:31:20Z
assignee: "Alexander Sacharov"
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-772869
claimed-at: 2026-09-14T20:17:54Z
outcome-what: an einen Worker uebergeben
outcome-why: "Transport-Passagen aus den installierten Kopien nach core/role/builtin portieren"
---

# Der Dispatcher-Prompt nennt sein Transportmittel nicht, also erfindet jeder Lauf ein eigenes

## Definition of Done

- [ ] jaira-dispatcher/SKILL.md nennt scripts/spawn.sh dort, wo ein Worker gestartet wird, nicht erst im Abschnitt ueber Ports - und sagt dazu, dass ein selbst aufgerufenes 'claude --permission-mode ...' vom Berechtigungspruefer abgelehnt wird, damit niemand es noch einmal versucht.
- [ ] jaira-teamlead/SKILL.md nennt dasselbe an der Stelle, an der es einen Dispatcher in eine Vorlage schickt - heute steht dort nur 'Run herdr --skill for the mechanics'.
- [ ] Beide Prompts sagen, dass 'herdr' auf einem WSL-Rechner nicht unter diesem Namen im PATH stehen muss und HERDR_BIN_PATH die verlaessliche Antwort ist - 'command -v herdr' beantwortet die Frage falsch.
- [ ] scripts/spawn.sh legt Zweige mit dem Praefix an, den dieses Repository benutzt, nicht mit feature/.
- [ ] Der .env-Block in scripts/spawn.sh traegt keinen fest eingebauten Projektnamen eines fremden Repositories mehr - entweder abgeleitet oder aus dem Skript heraus.
- [ ] Nachgestellt: ein Dispatcher, der nur seinen eigenen Prompt liest, startet einen Worker in einer eigenen Vorlage, ohne 'claude --permission-mode' selbst aufzurufen.
- [ ] Eine Zeile in core/release/NOTES.md unter ## Unreleased, weil die ausgelieferten Prompts sich aendern.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-14 18:23 · Alexander Sacharov** — Befund aus dem ersten echten Gebrauch (Dispatcher fuer 13VMA8, 14.09.): scripts/spawn.sh und SKILL.md widersprechen sich beim Transport. SKILL.md sagt ausdruecklich 'Not a split pane: a tab' und nennt 'herdr tab create --cwd ... --label ... --no-focus'. spawn.sh macht stattdessen 'herdr pane split --current --direction down --no-focus'. Wer dem Prompt folgt, benutzt das Skript nicht; wer das Skript benutzt, bricht den Prompt. Zweiter Befund: spawn.sh legt den Worktree mit -b feature/SLUG von HEAD des Repos an - dieses Repository benutzt feat/, und ein Dispatcher braucht oft eine andere Basis als HEAD (hier daf4312 auf feat/9ET6NC-per-board-remote); der Worktree musste deshalb von Hand vorher angelegt werden, damit spawn.sh ihn nur noch vorfindet. Dritter Befund: der .env-Block traegt COMPOSE_PROJECT_NAME=rg_SLUG - 'rg_' ist der Projektname eines fremden Repositories. Vierter Befund: spawn.sh startet immer /jaira-role-lane und kann /jaira-role-tester nicht starten, obwohl SKILL.md sagt 'Testing is not a lane'.
- **2026-09-14 20:18 · Alexander Sacharov** — Dispatcher-Lauf 2026-09-14 (2. Anlauf, Worktree jaira-13VMA8 auf feat/13VMA8-pr-is-the-humans). Transport: Herdr ist da, HERDR_ENV=1, HERDR_BIN_PATH=/mnt/c/.../herdr.exe - 'command -v herdr' waere hier wieder falsch gewesen. Ich benutze scripts/spawn.sh wie mein eigener Prompt es jetzt vorschreibt. Ausgangslage: die INSTALLIERTEN Kopien unter ~/.claude/skills sind heute von Hand gepatcht und dienen als Vorlage fuer den Transport-Abschnitt; sie tragen aber noch die ALTE Pull-Request-Regel, die 13VMA8 gerade umgedreht hat. Es werden nur die Transport-Passagen portiert, nicht die Dateien.
- **2026-09-14 20:19 · Alexander Sacharov** — Arbeitsanweisung fuer die in-progress-Lane (vom Dispatcher, damit sie auf dem Board steht und nicht in einer Sitzung stirbt):

Gearbeitet wird im Worktree /home/alex/projects/.worktrees/jaira-13VMA8 auf dem Zweig feat/13VMA8-pr-is-the-humans. KEINEN neuen Worktree und KEINEN neuen Zweig anlegen. /home/alex/projects/jaira und .worktrees/jaira-9ET6NC nicht anfassen.

Die Arbeit ist auf diesem Rechner schon getan, aber am falschen Ort: die INSTALLIERTEN Kopien unter ~/.claude/skills/jaira-dispatcher/SKILL.md, ~/.claude/skills/jaira-teamlead/SKILL.md und ~/.claude/skills/jaira-dispatcher/scripts/spawn.sh sind heute von Hand gepatcht worden. 'diff -u core/role/builtin/<rolle>/SKILL.md ~/.claude/skills/<rolle>/SKILL.md' zeigt genau, was zu portieren ist.

FALLE: die installierten Kopien tragen noch die ALTE Pull-Request-Regel ('ein Worker oeffnet den Pull Request', 'Close it once the pull request is open'). Ticket 13VMA8 hat diese Regel in core/role/builtin gerade umgedreht: ein Agent schiebt seinen Zweig und hoert auf, der Mensch oeffnet den Pull Request. Nur die Transport-Absaetze portieren, niemals die Dateien im Ganzen kopieren, sonst wird 13VMA8 rueckgaengig gemacht.

Zusaetzlich, von den installierten Kopien nicht abgedeckt:
- spawn.sh Zeile 18: '-b feature/$slug' -> dieses Repository benutzt 'feat/'.
- spawn.sh Zeilen 30/34/35: COMPOSE_PROJECT_NAME=rg_$slug, VITE_PORT_HOST, BACKEND_PORT_HOST sind der Stapel eines fremden Projekts (requirementsgenie). Kein fest eingebauter fremder Projektname mehr - ableiten oder herausnehmen.
- core/role/builtin/jaira-dispatcher/SKILL.md:157 'not when the pull request merges' begrenzt die Lebensdauer eines Worktrees. Unter der neuen Regel sieht ein Dispatcher nie einen offenen Pull Request; der Satz muss im selben Durchgang umgeschrieben werden.
- eine Zeile in core/release/NOTES.md unter ## Unreleased.

KEINEN Pull Request oeffnen, aktualisieren oder mergen. Zweig schieben und aufhoeren.
- **2026-09-14 20:31 · Alexander Sacharov** — Nachgestellt (DoD-Punkt 6), und zwar von diesem Dispatcher selbst: mein eigener Prompt sagt mir inzwischen, scripts/spawn.sh zu benutzen. Aufruf war 'bash ~/.claude/skills/jaira-dispatcher/scripts/spawn.sh 13VMA8 KSGSKK in-progress /home/alex/projects/jaira'. Es lief beim ersten Versuch durch und gab die Pane-Id 'w2:p1W' zurueck - ein eigener Tab, kein Split. Ich habe 'claude --permission-mode' nirgends selbst aufgerufen, also gab es auch keine Ablehnung als 'Create Unsafe Agents'. Das Skript hat mir nicht widersprochen; es gibt aus diesem Lauf keinen neuen Befund gegen das Skript.

Ein Detail fuer den, der spawn.sh anfasst: weil der Worktree schon existierte (.worktrees/jaira-13VMA8), lief der ganze Block mit 'git worktree add -b feature/$slug' und dem .env-Teil gar nicht erst an. Genau deshalb sind diese beiden Fehler bisher niemandem aufgefallen - sie liegen in dem Zweig, der nur beim allerersten Aufruf fuer einen Slug ausgefuehrt wird.

Kleiner Stolperstein ohne Bezug zum Ticket: 'jaira show <id> --json' hat kein Feld 'lane', die Lane steht in 'status'.
