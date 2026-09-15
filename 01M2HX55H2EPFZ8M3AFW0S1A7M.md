---
id: 01M2HX55H2EPFZ8M3AFW0S1A7M
title: jaira lanes use --all bringt alle Lanes eines Boards auf den Stand des Binaries
status: backlog
ready: true
creator: Alexander Sacharov
goal: Nach einem Upgrade des Binaries bringt ein Aufruf alle Lanes des Boards auf den mitgelieferten Stand - ohne dass jemand sie einzeln aufzaehlt und ohne dass eine handgeaenderte Lane dabei still verschwindet.
definition-of-done: "'jaira lanes use --all' bringt jede Lane des Boards auf den mitgelieferten oder Katalog-Stand. Nachgestellt an einem Board mit dreizehn Lanes, von denen mehrere veraltet sind: ein Aufruf genuegt, und die Ausgabe sagt je Lane, was passiert ist."
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-15T06:48:02Z
updated-at: 2026-09-15T06:48:28Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
context: |-
  Nach 'jaira self update' haengen die Lanes eines Boards zurueck. Sie aktualisieren sich nicht mit dem Binary: eine Lane-Datei ist die Lane, und das Board behaelt die Fassung, die einmal geschrieben wurde.

  Heute geht das nur einzeln. 'jaira lanes use <id> --force' nimmt genau eine Lane. Dieses Board hat dreizehn, also dreizehn Aufrufe - und man muss vorher wissen, welche ueberhaupt installiert sind.

  Der Fall ist nicht theoretisch: die Rollen haben mit 'jaira roles install --global --force' laengst einen Sammelbefehl, die Lanes nicht.

  Die Gefahr, die der Schalter nicht mitbringen darf: eine handgeaenderte Lane darf nicht still verschwinden. Am 2026-09-14 wurde auf diesem Board 'logbook-on-entry: true' von Hand aus .jaira/lanes/done.md entfernt, weil ein Move nach done drei fremde fertige Tickets ins Logbuch fegte. Ein blindes --force ueber alle Lanes haette das zurueckgeschrieben. 'jaira roles install' macht es an dieser Stelle richtig vor: was der Mensch geaendert hat, wird gemeldet und nicht angefasst, der Befehl endet mit Code 3, und erst --force ueberschreibt.

  Seit dem 2026-09-15 liegen die Lane-Dateien dieses Boards unter git (Commit 42d8253) statt in .gitignore. Damit reisen sie mit dem Zweig, und ein Sammelbefehl aendert etwas, das andere sehen - ein weiterer Grund, dass er sagt was er tut, statt es zu tun.

  Nicht Teil dieses Tickets, aber die groessere Frage dahinter: 'jaira update' verspricht in seiner eigenen Hilfe, das Setup eines Boards auf den Stand des installierten Binaries zu bringen, regeneriert aber nur den Textblock in CLAUDE.md und AGENTS.md. Rollen und Lanes gehoeren genauso dazu.
---

# jaira lanes use --all bringt alle Lanes eines Boards auf den Stand des Binaries

## Definition of Done

- [ ] 'jaira lanes use --all' bringt jede Lane des Boards auf den mitgelieferten oder Katalog-Stand. Nachgestellt an einem Board mit dreizehn Lanes, von denen mehrere veraltet sind: ein Aufruf genuegt, und die Ausgabe sagt je Lane, was passiert ist.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

