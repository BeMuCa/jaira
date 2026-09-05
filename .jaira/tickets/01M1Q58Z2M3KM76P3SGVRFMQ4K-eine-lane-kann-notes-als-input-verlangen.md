---
id: 01M1Q58Z2M3KM76P3SGVRFMQ4K
title: Eine Lane kann notes als Input verlangen
status: critique
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "input-requires darf 'notes' nennen: jaira show --for-lane liefert dann die Progress-Notes des Tickets im bounded input - Findings aus testing/critique erreichen den Implementing-Agenten ohne Voll-Lektuere; builtin in-progress deklariert notes"
context: "Berk am 04.09. (Feedback L120): Notes sind der Uebergabekanal zwischen Lanes (testing schreibt Befund+Loesungsvorschlag als Note), aber show --for-lane reicht nur deklarierte Felder - der Implementing-Agent muss zufaellig das volle Ticket lesen. Fix an der Quelle: 'notes' als erlaubter input-requires-Wert, der die Progress-Sektion mitliefert. Builtin in-progress bekommt notes in input-requires - bekannter Preis: bestehende Board-Kopien driften bis zum Refresh (gleicher Weg wie done.md gestern)."
definition-of-done: "show <id> --for-lane in-progress --json enthaelt die Notes des Tickets; eine Lane ohne notes-Deklaration bekommt sie nicht; builtin in-progress deklariert notes; ein Test deckt beide Faelle"
tags: []
blocked-by: []
commits: []
created-at: 2026-09-04T21:30:25Z
updated-at: 2026-09-05T02:41:27Z
claimed-by: EE-3NX6GL3-413400
claimed-at: 2026-09-05T02:34:06Z
updated-by: BeMuCa
outcome-what: "input-requires kennt 'notes': show --for-lane liefert die Progress-Eintraege des Tickets als bounded input (leeres Journal wird weggelassen, nie als missing gemeldet); builtin in-progress deklariert notes; SuppliedFields haelt den Lane-Validator ruhig"
outcome-why: "Feedback L120: testing/critique geben Findings als Note zurueck, aber der Implementing-Agent bekam sie im bounded input nie zu sehen - der Uebergabekanal endete vor dem Empfaenger"
outcome-resolves: "Drei Faelle getestet (geliefert wenn deklariert+vorhanden; weggelassen wenn leer; nie ungefragt an andere Lanes); Preis dokumentiert: bestehende Board-Kopien von in-progress driften bis zum Refresh (dieses Board ist refreshed)"
executed-by: fable
---

# Eine Lane kann notes als Input verlangen

## Definition of Done

- [x] show <id> --for-lane in-progress --json enthaelt die Notes des Tickets; eine Lane ohne notes-Deklaration bekommt sie nicht; builtin in-progress deklariert notes; ein Test deckt beide Faelle
  proof: flow.go showForLane: notes-Fall (Progress via progressNotes, leer=weggelassen statt missing); ticket.SuppliedFields kennt notes (Lane-Validator warnt nicht); builtin 20-in-progress deklariert notes, Board-Kopie refreshed; TestForLaneDeliversNotesWhenDeclared deckt liefern/weglassen/nicht-deklariert; voller Lauf -race RC=0 (test20)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

