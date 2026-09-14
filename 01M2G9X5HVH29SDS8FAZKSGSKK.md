---
id: 01M2G9X5HVH29SDS8FAZKSGSKK
title: "Der Dispatcher-Prompt nennt sein Transportmittel nicht, also erfindet jeder Lauf ein eigenes"
status: backlog
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
updated-at: 2026-09-14T15:52:42Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
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

