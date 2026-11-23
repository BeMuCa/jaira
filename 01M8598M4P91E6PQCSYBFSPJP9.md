---
id: 01M8598M4P91E6PQCSYBFSPJP9
title: "Eine Lane-Datei kann mit agents: N sagen, dass N Agenten die Lane laufen; show --for-lane reicht es weiter"
status: pre-process
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Eine lesende Lane laesst sich in ihrer Lane-Datei mit agents: N besetzen; jaira prueft den Wert, reicht ihn an den Agenten weiter, und die Rollen-Skills sagen, wie N Ergebnisse zu einem werden."
context: |-
  Heute laeuft jede Lane mit genau einem Agenten; es gibt keinen Schluessel dafuer (ATDG74, Idee von Berk: 5x Haiku in critique).
  jaira startet selbst keine Agenten. Der Schwarm ist deshalb: ein Schluessel in der Lane-Datei, den 'jaira show <id> --for-lane <lane> --json' ausgibt, und ein Absatz in den Skills jaira-worker und jaira-role-lane (core/role/builtin/).
  Entschieden im Brainstorm zu ATDG74 (Teamlead, 2026-11-23):
  - Schluessel 'agents: N' im Frontmatter, ganze Zahl 1 bis 9, fehlt er, gilt 1 und nichts aendert sich.
  - Erlaubt nur auf lesenden Lanes. agents > 1 auf einer Lane, deren output-produces outcome-what oder plan enthaelt, ist ein Fehler beim Laden - N Agenten im selben Worktree kollidieren.
  - Zusammenfuehren: Vereinigung der Befunde; der tragende Agent prueft jeden Befund am Code, streicht Doppelte und falsche, und schreibt EIN Ergebnis und EIN Verdict. Keine Mehrheit (verliert den Befund, den nur einer fand), kein eigener Zusammenfasser (prueft nichts).
  - Speicher: vor jedem Start free -h; unter 2 GiB verfuegbar nacheinander statt parallel; nie mehr als 4 zugleich. Das steht im Skill-Text, das Binary startet nichts.
  Code: core/lane/lane.go (Lane-Struct, Parser um Zeile 333-440, max-rounds ist das Muster fuer einen Zahlen-Schluessel; Gleichheit um Zeile 960), core/lane/facts.go, die Ausgabe von show --for-lane in internal/cli (dort wo model_tier gesetzt wird), 'jaira lanes show'.
  Achtung: der Zweig origin/feat/ENCFVZ-skill-lane (nicht gemergt) aendert dieselben Stellen in lane.go; klein halten, damit der Konflikt klein bleibt. Aendert sich eine Datei unter core/role/, verlangt ein Test eine neue core/shipped/shipped.sha256 - nachlesen, wie die erzeugt wird.
definition-of-done: "agents: N in einer Lane-Datei wird gelesen; 0, negative Werte, mehr als 9 und Nicht-Zahlen werden beim Laden mit einer Meldung abgelehnt, die Datei und Schluessel nennt; Test in core/lane"
tags:
  - concurrency
  - cli
blocked-by: []
related: []
follows: 01M67K72RDJY2N2R2KQGATDG74
commits: []
created-at: 2026-11-23T22:43:41Z
updated-at: 2026-11-23T22:56:58Z
updated-by: BeMuCa
claimed-by: worker-fspjp9
claimed-at: 2026-11-23T22:56:02Z
---

# Eine Lane-Datei kann mit agents: N sagen, dass N Agenten die Lane laufen; show --for-lane reicht es weiter

## Definition of Done

- [ ] agents: N in einer Lane-Datei wird gelesen; 0, negative Werte, mehr als 9 und Nicht-Zahlen werden beim Laden mit einer Meldung abgelehnt, die Datei und Schluessel nennt; Test in core/lane
- [ ] agents > 1 auf einer Lane, die outcome-what oder plan liefert, wird beim Laden abgelehnt; Test
- [ ] jaira show <id> --for-lane <lane> --json traegt agents (1 ohne Schluessel); die Textausgabe und jaira lanes show nennen es nur, wenn es groesser als 1 ist; Test in internal/cli
- [ ] jaira-worker und jaira-role-lane sagen, was bei agents > 1 zu tun ist: N Agenten, Speicherregel, Befunde vereinigen und am Code pruefen, ein Verdict
- [ ] docs und Release-Notiz in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [ ] read the max-rounds path end to end: lane.go parser 422-431, flow.go 750 and 815, tickets.go lanes show 1523 and 1536
- [ ] failing tests in core/lane: default 1, 1 to 9 accepted, 0 -1 10 x 2.5 refused with file and key named, agents>1 refused on an outcome-what or plan lane
- [ ] failing test in internal/cli: show --for-lane json agents 1 without key and 3 with key; text line only above 1; lanes show
- [ ] implement agents in core/lane/lane.go: field next to MaxRounds, parser, equality; no change to facts.go (the board note never shows lane agents)

## Progress

