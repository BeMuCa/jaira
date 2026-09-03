---
id: 01M1KAFM6C20J8H7Z3178W1R94
title: "Repo ziehen und jaira updaten ist ein einfacher, durchgespielter Prozess - und der Share-Zustand bleibt dabei immer erhalten"
status: todo
ready: false
creator: BeMuCa
goal: "Der komplette Weg 'jaira installieren/updaten + Board-Repo ziehen/aktualisieren' ist einmal Schritt fuer Schritt durchgespielt, so einfach wie moeglich gemacht und dokumentiert; bei JEDEM Update-Pfad bleibt der Share-Zustand des Boards erhalten (shared bleibt shared, privat bleibt privat)"
context: "Berk am 03.09. vor einem Context-Clear: wenn wir ein Update von jaira ziehen, muss der Share-Zustand im Repo beibehalten werden. Teilweise schon gebaut: ETR0PX (accepted) - jaira update schreibt nur noch den Agent-Block, nie .gitignore; nur jaira init gitignoriert. ZU PRUEFEN: gilt das auf ALLEN Update-Pfaden (jaira self upgrade, go build/install eines neuen Stands, jaira update auf geteiltem vs. privatem Board, erster Befehl nach Binary-Wechsel inkl. Lane-Dir-Migration und Agent-Block-Refresh)? Zweiter Teil: den Prozess einmal als Nutzer durchgehen - frisches Clone eines geteilten Boards, jaira installieren (welcher Befehl? go install trifft NICHT ~/.local/bin - bekannte Falle), Board oeffnen, Update ziehen, jaira update ausfuehren - jede Reibung notieren und glaetten (READMEs, Fehlermeldungen, evtl. ein 'jaira doctor'-artiger Check). DoD-Messlatte: ein Teammate ohne Vorwissen kommt mit README allein durch."
definition-of-done: "Ein dokumentierter Durchlauf clone->install->open->update existiert (Skript oder README-Abschnitt); ein Test oder Durchlauf belegt je Update-Pfad, dass .gitignore/Share-Zustand unangetastet bleibt; gefundene Reibungspunkte sind als Tickets eingefangen oder gefixt"
tags: []
blocked-by: []
commits: []
created-at: 2026-09-03T09:44:29Z
updated-at: 2026-09-03T09:44:29Z
---

# Repo ziehen und jaira updaten ist ein einfacher, durchgespielter Prozess - und der Share-Zustand bleibt dabei immer erhalten

## Definition of Done

- [ ] Ein dokumentierter Durchlauf clone->install->open->update existiert (Skript oder README-Abschnitt); ein Test oder Durchlauf belegt je Update-Pfad, dass .gitignore/Share-Zustand unangetastet bleibt; gefundene Reibungspunkte sind als Tickets eingefangen oder gefixt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

