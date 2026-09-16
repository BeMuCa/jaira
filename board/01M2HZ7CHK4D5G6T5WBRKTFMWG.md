---
id: 01M2HZ7CHK4D5G6T5WBRKTFMWG
title: "Die Zeile, die ein Mensch in der Sackgasse kopiert, setzt einen Wert fuer ihn ein"
status: backlog
ready: true
creator: Alexander Sacharov
goal: "Wer in der Forge-Sackgasse landet, bekommt eine Zeile, die nichts fuer ihn entscheidet - er waehlt, und die Rolle raet auch im Hilfetext nicht."
definition-of-done: "Die Zeile, die jaira-role-pr an Rung 4 herausgibt, setzt keinen Wert vor: ein Mensch, der sie unbesehen kopiert, bekommt entweder einen Fehler oder muss waehlen, aber nie stillschweigend 'gitlab'. Nachgestellt, indem die Ausgabe woertlich in eine Shell kopiert wird."
tags:
  - cli
blocked-by: []
related:
  - 01M28FMQGC4CNQT8Z9WY13VMA8
commits: []
created-at: 2026-09-15T07:24:12Z
updated-at: 2026-09-15T07:24:45Z
assignee: ""
updated-by: Alexander Sacharov
context: "Gefunden von der review-Lane von 13VMA8 am 2026-09-15, als nicht blockierend eingestuft und von Alex bei der Abnahme bewusst durchgelassen, um 0.2.1 nicht aufzuhalten.\n\ncore/role/builtin/jaira-role-pr/SKILL.md:62 gibt dem Menschen an Rung 4 der Forge-Leiter diese Zeile:\n\n  git config jaira.forge gitlab    # or github\n\nWer sie unbesehen kopiert - und in einer Sackgasse kopiert man unbesehen -, setzt gitlab. Auf einer GitHub-Enterprise-Instanz unter eigener Domain ist das die falsche Haelfte.\n\nWarum es mehr ist als ein Schoenheitsfehler: Rung 4 existiert, WEIL die Rolle den Forge nicht erraten darf. Sie haelt sich im Verhalten daran und raet dann im Hilfetext doch - nur eben an der Stelle, an der der Mensch nicht hinsieht. Alex hat den Grundsatz am 2026-09-15 so formuliert: die Zeile vorschlagen, und wenn der Mensch sie nicht ausfuehrt, haben wir keine Information - erfinden ist nicht noetig.\n\nNicht schlimm genug fuer einen Rueckruf: der falsche Wert faellt beim naechsten Lauf laut um, er verdirbt nichts still.\n\nZwei Zeilen statt einer mit Kommentar waeren der naheliegende Weg; welcher es wird, entscheidet die Arbeit."
---

# Die Zeile, die ein Mensch in der Sackgasse kopiert, setzt einen Wert fuer ihn ein

## Definition of Done

- [ ] Die Zeile, die jaira-role-pr an Rung 4 herausgibt, setzt keinen Wert vor: ein Mensch, der sie unbesehen kopiert, bekommt entweder einen Fehler oder muss waehlen, aber nie stillschweigend 'gitlab'. Nachgestellt, indem die Ausgabe woertlich in eine Shell kopiert wird.
- [ ] Derselbe Grundsatz gilt ueberall, wo ein Prompt eine Befehlszeile zum Kopieren herausgibt und einen von mehreren gueltigen Werten einsetzen muesste: nachgesehen in core/role/builtin, und wo es noch vorkommt, mitgezogen.
- [ ] Eine Zeile in core/release/NOTES.md unter ## Unreleased, weil die ausgelieferten Prompts sich aendern.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

