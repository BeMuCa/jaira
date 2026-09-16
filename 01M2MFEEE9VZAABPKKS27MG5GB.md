---
id: 01M2MFEEE9VZAABPKKS27MG5GB
title: "spawn.sh schreibt Werte, die nur einem Projekt gehoeren"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "spawn.sh schreibt nur, was jaira selbst besitzt. Alles, was den Namen eines Projekts traegt, kommt aus einer Datei im Repository - .jaira/worktree-env, ausfuehrbar, ihre Ausgabe wird an die .env des neuen worktree gehaengt. Fehlt die Datei, verhaelt sich spawn.sh wie bisher."
context: |-
  spawn.sh legt fuer jeden Arbeiter einen worktree an und schreibt ihm eine eigene .env. Darin stehen heute HTTP_PORT und DB_PORT_HOST - Namen, die sich ein einzelnes Projekt ausgedacht hat. Ein anderes Projekt nennt seine Ports anders oder hat keine, und bekommt die Zeilen trotzdem.

  Die Folge sieht man an einem Board, das jaira benutzt: dessen spawn.sh ist von 100 auf 159 Zeilen gewachsen, weil dort IMAGE_NS, COMPOSE_PROFILES, VITE_PORT_HOST und BACKEND_PORT_HOST nachgetragen wurden. Alles vier gehoert diesem einen Projekt.

  Das Schlimme daran ist nicht das Ueberschreiben - 'roles install' laesst eine bearbeitete Datei stehen und meldet sie. Es ist das Gegenteil: die Datei ist eingefroren. Jede Verbesserung an spawn.sh wird bei diesem Anwender still uebersprungen, und er merkt es erst, wenn er die Ausgabe von 'roles install' liest. Der Fork ist auch nicht sortenrein: --no-worktree und der Lane-Name 'dispatch' stehen darin und sind allgemein, sie gehoeren hierher.

  Am 16.09.2026 aufgekommen, als in einem worktree die COMPOSE_PROFILES des Hauptverzeichnisses mitkopiert wurden und der Zweig acht Container zu viel hochfuhr.

  Noch offen und Teil der Loesung: welche Angaben der Hook braucht. Vorschlag: Slug und Port-Versatz als Argumente, damit die Datei ihre Ports genauso verteilen kann wie spawn.sh heute.
definition-of-done: "spawn.sh ruft .jaira/worktree-env auf, wenn die Datei existiert und ausfuehrbar ist, und haengt ihre Ausgabe an die .env des worktree"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-16T06:46:09Z
updated-at: 2026-09-16T06:46:09Z
---

# spawn.sh schreibt Werte, die nur einem Projekt gehoeren

## Definition of Done

- [ ] spawn.sh ruft .jaira/worktree-env auf, wenn die Datei existiert und ausfuehrbar ist, und haengt ihre Ausgabe an die .env des worktree
- [ ] Ohne diese Datei schreibt spawn.sh dieselbe .env wie vorher - ein bestehendes Board merkt nichts
- [ ] HTTP_PORT und DB_PORT_HOST stehen nicht mehr in spawn.sh; die Dokumentation zeigt, wie ein Projekt sie in seinem Hook setzt
- [ ] Die Hilfe von 'roles install' sagt, dass eine bearbeitete Datei nicht nur unveraendert bleibt, sondern auch keine Verbesserungen mehr bekommt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

