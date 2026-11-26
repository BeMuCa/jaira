---
id: 01M85APJ60PQG282EDVDZ9JNAV
title: "Der Katalog nimmt keine Lane mit script: an, und adopt nennt skill: und script: beim Namen"
status: backlog
ready: true
creator: BeMuCa
assignee: Berk
goal: "Eine geteilte Lane kann niemanden unbemerkt ein Skript ausfuehren lassen: der oeffentliche Katalog lehnt script: ab, und beim Uebernehmen aus jeder anderen Quelle steht es gross da."
context: |-
  ENCFVZ (Stand 24.11.2026: signoff, nicht gemerged) gibt Lanes drei neue Felder: skill:, script: und receipt:. script: nennt ein Skript im Repo, das der Agent in dieser Lane ausfuehrt.
  Eine Lane aus dem Marktplatz mit script: scripts/deploy.sh laesst also Code laufen, sobald ein Ticket in die Lane kommt. Der Pfad muss im Repo liegen (ENCFVZ, Test TestScriptRefusesAPathOutsideTheRepository) - die Lane bringt das Skript nicht mit, aber sie waehlt aus, welches laeuft.
  skill: ist schwaecher, aber dieselbe Art: die Lane bestimmt, welche fremde Anweisung der Agent laedt.
  Heute prueft das nichts, weil es die Felder auf master noch nicht gibt.
  Loesung:
  - 'jaira market lint' lehnt im Katalog jede Lane mit script: ab. Regel-Name no-script.
  - 'jaira lanes market adopt' und 'jaira boards market adopt' nennen skill: und script: in der Zusammenfassung in einer eigenen Zeile, und eine Lane mit script: braucht zusaetzlich --allow-script (fuer Firmen-Kataloge, die es bewusst erlauben).
  - 'jaira lanes adopt <pfad>' fuer die Lane eines Teammitglieds (.jaira/shared/<du>/) nennt die beiden Felder ebenso.
  Erst anfangen, wenn ENCFVZ auf master ist.
  Vorgabe: docs/MARKET.md.
definition-of-done: "jaira market lint meldet no-script fuer eine Lane mit script: (Test)"
tags:
  - cli
  - gates
blocked-by:
  - 5TED33
  - R0YRGP
  - ENCFVZ
related: []
commits: []
created-at: 2026-11-23T23:08:47Z
updated-at: 2026-11-26T10:51:31Z
updated-by: Berk
---

# Der Katalog nimmt keine Lane mit script: an, und adopt nennt skill: und script: beim Namen

## Definition of Done

- [ ] jaira market lint meldet no-script fuer eine Lane mit script: (Test)
- [ ] market adopt einer Lane mit script: ohne --allow-script endet mit Exit 3 und schreibt nichts (Test)
- [ ] Die Zusammenfassung vor dem Uebernehmen nennt skill: und script: mit ihrem Wert
- [ ] jaira lanes adopt <pfad> nennt die Felder ebenfalls
- [ ] Release-Note core/release/notes/<TICKET>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

