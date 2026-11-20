---
id: 01M7WWZRD7WNC5NT1PS5TW5FE2
title: "Randfaelle werden vor dem Code benannt, und wer umsetzt, prueft die Fehlerklasse statt des Beispiels"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "pre-process schreibt Fehlerklasse, Eingabeklassen und bei Heuristiken eine Annahmeregel als DoD-Punkte; in-progress testet die Klasse, sucht Geschwister-Stellen und prueft selbst die vier Kategorien, bevor es abgibt; bughunt meldet einen Fund als Klasse."
context: |-
  Problem: Review und bughunt finden in jedem Ticket etwas, Schleifen laufen 10-24 Runden. Die Funde sind meist richtig.
  Analyse 20.11.2026 an 34 Tickets eines Projekts: bughunt 47 Funde aus der Aenderung selbst, 5 vorbestehend; critique fast immer "nichts"; Review am Ende 0 Blocker, 2-6 Notes.
  Ursache 1: die vier Kategorien (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte) stehen nur in testing/bughunt/review (.jaira/lanes/review.md:31). pre-process, in-progress (Tier cheap, "work only within the plan's scope") und critique ("Do not check whether it works") fragen nie danach.
  Ursache 2: ein Fund kommt als ein Beispiel, der Fix testet das Beispiel, die naechste Runde findet den Nachbarfall (wiederholt beobachtet: Seitenzahl "3 / 12" dann "3/12"; Zahlen "0,05", dann "2,50", dann "1.234,56").
  Ursache 3: Heuristik-Tickets (PDF, Zahlen) haben keine Annahmeregel; die Schleife konvergierte erst, als der Teamlead eine setzte ("Textverlust ist schlimmer als ein uebrig gebliebener Kopf/Fuss").
  Ursache 4: Tickets beschreiben das Symptom ("NUL in Werkzeug A"), nicht die Klasse ("alles, was die DB nicht speichert, in jedem Schreibwerkzeug").
  Ursache 5: an der Wurzel reparieren statt am Beispiel hat den Nachbarfall-Zyklus beendet (strukturierte Daten statt String parsen).
  Alex: pre-process baut die DoD so, dass in-progress weiss, was zu tun ist; aber wer umsetzt, muss selbst denken statt Haken setzen, sonst pruefen ihn danach zehn andere.
  Eine Explain-Taste nur fuer blocker: gehoert zum Ticket ueber die Schleifenbremse.
definition-of-done: "pre-process-Prompt verlangt: Fehlerklasse statt Symptom, Eingabeklassen, bei Heuristik eine Annahmeregel - als DoD-Punkte"
tags:
  - gates
blocked-by: []
related:
  - 01M7WWZWSDHK37DGGHJFDF0YD6
  - 01M67RWZDWBS8E44Z3KGXKS048
  - 01M67K7JAF6R9N2VSPEW3CBCN6
  - 01M5ZJ4NNCYN6YBWVAPWAHG254
commits: []
created-at: 2026-11-20T16:35:12Z
updated-at: 2026-11-20T16:45:57Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-42185
claimed-at: 2026-11-20T16:42:15Z
---

# Randfaelle werden vor dem Code benannt, und wer umsetzt, prueft die Fehlerklasse statt des Beispiels

## Definition of Done

- [x] pre-process-Prompt verlangt: Fehlerklasse statt Symptom, Eingabeklassen, bei Heuristik eine Annahmeregel - als DoD-Punkte
  proof: core/lane/builtin/15-pre-process.md:38 (Name the edges before the code: class, input classes, acceptance rule, jaira dod --add)
- [x] in-progress-Prompt verlangt: Tests ueber die Klasse (parametrisiert), Suche nach Geschwister-Stellen, Selbstpruefung der vier Kategorien vor dem Move, Fix an der Ursache statt am Beispiel
  proof: core/lane/builtin/20-in-progress.md:29 (Fix the class, not the example: cause, parametrised test, siblings in scope, four categories before move)
- [ ] bughunt- und testing-Prompt: ein Fund nennt die Klasse und das Beispiel; der Rueckweg nach in-progress verlangt den Fix der Klasse
- [ ] Die eingebetteten Lanes und Rollen sind geaendert, mit Release-Note

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-20 16:35 · Alexander Sacharov** — Alex 20.11.2026: er schickt inzwischen jedes Ticket wegen der Review-Funde zurueck, bis die Raender zu sind ('Raender offen lassen ist schlecht'). Jede Note ist also faktisch eine weitere Runde mit dem Menschen - die Raender billiger vor dem Code zu benennen spart genau diese Runden.
- **2026-11-20 16:42 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/35112.sock
- **2026-11-20 16:43 · Alexander Sacharov** — decision: keine offene Entscheidung fuer die Person. Selbst gesetzt: (1) pre-process laeuft nur mit Option planning (requires-option) - deshalb benennt in-progress Klasse/Eingabeklassen/Annahmeregel selbst, wenn kein Plan-Lane lief; (2) 'Work only within the plan's scope' bleibt, aber Geschwister-Stellen derselben Fehlerklasse gelten ausdruecklich als im Scope, andere Nachbarprobleme weiter nur als Note; (3) pre-process schreibt die Punkte mit 'jaira dod <id> --add' in die DoD, nicht nur in den Plan, weil nur die DoD das Gate liest; (4) critique, review und Tier bleiben unveraendert - nicht in der DoD; (5) jede geaenderte eingebaute Datei bekommt ihren sha256 in core/shipped/shipped.sha256, .jaira/lanes bekommt dieselbe Fassung.
