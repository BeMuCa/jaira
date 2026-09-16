---
id: 01M2MFEEE9VZAABPKKS27MG5GB
title: "spawn.sh schreibt Werte, die nur einem Projekt gehoeren"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "spawn.sh legt den worktree an, oeffnet die Kachel und startet den Arbeiter - mehr nicht. Alles, was ein Projekt zum Loslaufen braucht, macht ein Skript im Repository, das spawn.sh aufruft und mit worktree-Pfad, Slug, Port-Versatz und Wurzel versorgt. Ohne dieses Skript legt spawn.sh nur den worktree an."
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
updated-at: 2026-09-16T06:47:16Z
updated-by: Alexander Sacharov
---

# spawn.sh schreibt Werte, die nur einem Projekt gehoeren

## Definition of Done

- [ ] spawn.sh ruft ein Einrichtungsskript des Repositorys auf und uebergibt ihm worktree-Pfad, Slug, Port-Versatz und Repository-Wurzel
- [ ] Ohne diese Datei schreibt spawn.sh dieselbe .env wie vorher - ein bestehendes Board merkt nichts
- [ ] HTTP_PORT und DB_PORT_HOST stehen nicht mehr in spawn.sh; die Dokumentation zeigt, wie ein Projekt sie in seinem Hook setzt
- [ ] Die Hilfe von 'roles install' sagt, dass eine bearbeitete Datei nicht nur unveraendert bleibt, sondern auch keine Verbesserungen mehr bekommt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-16 06:46 · Alexander Sacharov** — Praezisiert am 16.09.2026 auf Ansage: jaira soll den worktree anlegen und sonst nichts. Die erste Fassung dieses Tickets hat noch von einer Datei gesprochen, die Zeilen fuer die .env ausgibt - das ist zu eng gedacht.

Warum zu eng: das Kopieren der .env ist selbst schon projektspezifisch. Ein Projekt ohne .env braucht es nicht, ein anderes will daneben noch 'npm install' oder eine Datenbank hochfahren. Wenn jaira die .env kopiert und ein Hook nur Zeilen anhaengen darf, bleibt die eine Entscheidung, die am haeufigsten falsch ist, weiter in jaira.

Also: ein Einrichtungsskript, kein Zeilenlieferant. jaira ruft es auf, gibt ihm worktree-Pfad, Slug, Port-Versatz und Repository-Wurzel, und was danach in dem Verzeichnis steht, ist Sache des Projekts. Rueckwaerts vertraeglich bleibt es dadurch, dass ein Board ohne dieses Skript einen worktree ohne .env bekommt - wer die alte Bequemlichkeit will, schreibt sich das Skript einmal hin, und jaira zeigt eins als Vorlage.
