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
updated-at: 2026-09-14T19:54:49Z
assignee: "Alexander Sacharov"
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
- [ ] scripts/spawn.sh gibt jedem Worker einen eigenen TAB statt einer geteilten Vorlage: 'herdr tab create --label <ticket>/<lane> --no-focus', die Pane-Id kommt aus result.root_pane.pane_id. Ein Split teilt die Hoehe eines Bildschirms - bei vier Workern bleiben je ein paar Zeilen und niemand kann lesen, was einer tut. --cwd nimmt Herdr nur fuer das Label entgegen und loest es gegen Windows auf, ein WSL-Pfad wird ignoriert; das echte Wechseln des Verzeichnisses bleibt das 'cd' im pane run.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-14 18:23 · Alexander Sacharov** — Befund aus dem ersten echten Gebrauch (Dispatcher fuer 13VMA8, 14.09.): scripts/spawn.sh und SKILL.md widersprechen sich beim Transport. SKILL.md sagt ausdruecklich 'Not a split pane: a tab' und nennt 'herdr tab create --cwd ... --label ... --no-focus'. spawn.sh macht stattdessen 'herdr pane split --current --direction down --no-focus'. Wer dem Prompt folgt, benutzt das Skript nicht; wer das Skript benutzt, bricht den Prompt. Zweiter Befund: spawn.sh legt den Worktree mit -b feature/SLUG von HEAD des Repos an - dieses Repository benutzt feat/, und ein Dispatcher braucht oft eine andere Basis als HEAD (hier daf4312 auf feat/9ET6NC-per-board-remote); der Worktree musste deshalb von Hand vorher angelegt werden, damit spawn.sh ihn nur noch vorfindet. Dritter Befund: der .env-Block traegt COMPOSE_PROJECT_NAME=rg_SLUG - 'rg_' ist der Projektname eines fremden Repositories. Vierter Befund: spawn.sh startet immer /jaira-role-lane und kann /jaira-role-tester nicht starten, obwohl SKILL.md sagt 'Testing is not a lane'.
- **2026-09-14 18:37 · Alexander Sacharov** — Die Richtung des Vorlagen-Splits ist kein Geschmacksurteil. Bei 'down' bekommt jeder weitere Worker einen Streifen der geteilten Hoehe: bei vier Workern bleiben je ein paar Zeilen uebrig, und der Sinn der Vorlagen - zusehen koennen, was ein Worker tut - ist weg. Nach rechts geteilt behaelt jede Vorlage die volle Hoehe und verliert nur Breite, was umbrochenen Text kostet und nicht die Sicht. Alex am 2026-09-14: 'иначе очень не видно нечего'. Einstellbar bleibt es trotzdem, weil bei vielen Workern auch die Breite irgendwann aufgebraucht ist.
- **2026-09-14 19:54 · Alexander Sacharov** — Am 2026-09-14 auf Alex' Bitte sofort in der INSTALLIERTEN Kopie ~/.claude/skills/jaira-dispatcher/scripts/spawn.sh gemacht, weil er mit geteilten Vorlagen nichts mehr lesen konnte und die Worker gerade liefen. Das dreht kurzzeitig die Richtung um, die PMF635 festgelegt hat - core/role/builtin ist die Quelle, ~/.claude/skills entsteht daraus. Die Kopie im Repository ist NICHT geaendert: 13VMA8 bearbeitet dieselben Dateien im Nachbar-Worktree, und zwei Haende auf einer Datei sind ein Konflikt. Solange das so steht, wuerde 'jaira roles install --global --force' den Patch ueberschreiben. Dieses Ticket ist der Ort, an dem die Aenderung ins Repository kommt; die Sicherung liegt als spawn.sh.bak daneben. Geprueft: 'tab create' liefert result.root_pane.pane_id und result.tab.tab_id, 'bash -n' ist sauber, und --cwd mit einem Linux-Pfad kam als 'C:\Users\Alex\' zurueck.
