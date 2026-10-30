---
id: 01M67K7JAF6R9N2VSPEW3CBCN6
title: "Experiment: macht ein Pre-Process mit Opus 5.5 die Umsetzung durch Haiku besser, billiger oder schneller?"
status: ideas
ready: true
creator: berkc
assignee: berkc
goal: "Mit Zahlen entscheiden, ob pre-process immer laufen soll und mit welchem Modell: Codequalitaet, Tokenkosten und Dauer fuer dasselbe Ticket, einmal nur Haiku, einmal Opus-Plan + Haiku."
context: |-
  Wunsch von berkc am 31.10.2026 (Punkt 23): 'it seems like if we used pre process always + with a strong model it could make haiku better? we should test them - you can use one of the jaira web tickets for example - code quality change and token cost change and speed change - running it without preprocessing fully on haiku only vs with preprocessing using opus 5.5 for only the preprocessing and the other lanes default model so probably haiku.'
  Beobachtung aus der Session 30.10.: ein Haiku-Worker (GBFM32) hat seinen eigenen Worktree samt ungesicherter Arbeit geloescht; ein anderer (G6T54H) hat vier Doku-Stellen uebersehen.
  Methode: dasselbe Ticket zweimal von origin/master auf getrennten Branches; Lauf A alle Lanes Haiku high ohne pre-process; Lauf B pre-process Opus, Rest Haiku high; ein dritter, frischer Opus-Agent bewertet beide Diffs blind.
definition-of-done: "Beide Laeufe mit Tokens, Dauer, Testergebnis und Blindbewertung als Notiz; Empfehlung"
tags:
  - gates
blocked-by: []
related:
  - 01M67C18P3AB3GAPJ103DAQKZX
  - 01M658C7SP4V29E95574CXH98W
commits: []
created-at: 2026-10-30T23:44:58Z
updated-at: 2026-10-30T23:44:58Z
---

# Experiment: macht ein Pre-Process mit Opus 5.5 die Umsetzung durch Haiku besser, billiger oder schneller?

## Definition of Done

- [ ] Beide Laeufe mit Tokens, Dauer, Testergebnis und Blindbewertung als Notiz; Empfehlung

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

