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
updated-at: 2026-10-31T00:03:20Z
updated-by: berkc
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
- **2026-10-31 00:03 · berkc** — Ergebnis 31.10.2026, Ticket NDE933 (z ueberspringt leere Lanes), beide von origin/master, Ticket nicht beruehrt. Lauf A (exp/NDE933-a, ed602e79): nur Haiku high, ohne Plan - 108k Tokens (120k inkl. dreifach wiederholtem Bericht), ca. 7,5 min Agentenzeit, 70 Tool-Aufrufe, 3 Dateien +94/-9, 5 Tests + Mutationsprobe. Lauf B (exp/NDE933-b, 31871f63): Plan durch Opus 5.5 (74,5k Tokens, 105 s, 17 Aufrufe), Umsetzung Haiku high nach Plan (85k Tokens, 11,5 min, 37 Aufrufe) - zusammen 160k Tokens, davon 74,5k Opus (teurer), 4 Dateien +142/-11, 6 Tests. Ein Planfehler (Testhelfer key('right') liefert Code 0) liess 5 Tests zuerst scheitern; Haiku hat ihn korrigiert. Beide: go vet und go test ./internal/tui gruen (von mir nachgeprueft). Blinde Bewertung durch frischen Opus-Agenten (Variante 1 = B, 2 = A): B 9/9/8/9 (Korrektheit/Tests/Stil/Vollstaendigkeit) gegen A 8/7/7/6; Empfehlung B. A liess den Hilfetext fuer z und einen nun falschen Kommentar stehen, keine Tests fuer Wrap-around und 'alles leer', und friert den Cursor ein, wenn alle Lanes leer sind. Fazit: Opus-Plan + Haiku ergab hier spuerbar vollstaendigere Arbeit bei etwa halber Zahl Tool-Aufrufe in der Umsetzung, kostet aber mehr Tokens und Geld (Opus-Anteil) und dauerte laenger. Eine Stichprobe von 1 - fuer eine Regel weitere Laeufe noetig. Hinweis: der Bewerter hatte kein Go im PATH und urteilte nur aus dem Code.
