---
id: 01M2GRVQSSSN75K3AY94NYW4M7
title: Die Definition of Done steht an zwei Stellen und die CLI haelt sie nicht zusammen
status: backlog
ready: false
creator: Alexander Sacharov
goal: "Ein Kriterium der Definition of Done hat eine Wahrheit, nicht zwei: was in der Checkliste steht, steht auch im Frontmatter, ohne dass jemand daran denkt."
context: |-
  Ein Ticket traegt seine Abnahmekriterien doppelt: als Kaestchen im Rumpf und als Feld 'definition-of-done' im Frontmatter. Die beiden koennen auseinanderlaufen, und niemand merkt es.

  Der gefaehrliche Teil ist, welche Fassung die Arbeit steuert: 'jaira show --for-lane' reicht das FRONTMATTER an jede Lane weiter. Ein Worker liest also die Fassung, die niemand mehr pflegt, und baut danach.

  Wie es auseinanderlaeuft: 'jaira create --dod' fuellt beide. 'jaira dod <id> <n> --text' aendert danach nur das Kaestchen. 'jaira dod --add' haengt ein Kaestchen an, das im Frontmatter nie auftaucht. Ein Weg, beide zusammen zu aendern, existiert nicht - 'jaira set definition-of-done=...' trifft nur das Feld, 'jaira dod --text' nur die Liste.

  Dreimal an einem Tag gesehen, am 2026-09-14:
  - S1VM40: critique fand das Frontmatter mit der alten Zusage ueber den dritten Farbplatz, waehrend der Rumpf schon die neue trug. Ohne diesen Fund haetten testing und review gegen eine ueberholte Bedingung geprueft.
  - 13VMA8: das Frontmatter sagte 'ein Agent macht ihn auf und merged ihn nie' - die Regel, die das Ticket gerade umgedreht hatte. Der Rumpf sagte das Gegenteil. Gefunden hat es die testing-Lane, nicht das Werkzeug.
  - APABM4, 9ET6NC und VM0A76 tragen denselben Rest aus einer verwandten Ecke: '--dod' nahm frueher nur einen Wert, alle weiteren Absaetze blieben als lose Prosa im Rumpf liegen und behaupten dort bis heute Dinge, die nicht mehr gelten.

  Nicht Teil dieses Tickets: dass '--dod' nur einen Wert nahm - das ist mit VM0A76 behoben.

  Offen und Sache der Plan-Lane: ob das Frontmatter-Feld ueberhaupt bleiben muss. Zwei Orte fuer eine Wahrheit sind der Fehler; ein abgeleitetes Feld waere ein Weg, ein einziger Ort ein anderer.
definition-of-done: "Ein Kriterium, das ueber 'jaira dod --text' umformuliert wird, steht danach in beiden Fassungen gleich - nachgestellt an einem Ticket, dessen Frontmatter vorher abwich."
tags:
  - cli
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-09-14T20:13:44Z
updated-at: 2026-09-14T20:14:06Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
---

# Die Definition of Done steht an zwei Stellen und die CLI haelt sie nicht zusammen

## Definition of Done

- [ ] Ein Kriterium, das ueber 'jaira dod --text' umformuliert wird, steht danach in beiden Fassungen gleich - nachgestellt an einem Ticket, dessen Frontmatter vorher abwich.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

