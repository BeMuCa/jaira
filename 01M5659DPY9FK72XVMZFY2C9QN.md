---
id: 01M5659DPY9FK72XVMZFY2C9QN
title: Ein zweites Board laesst sich mit einem Befehl aus einem Board-Typ anlegen statt Lane-Dateien von Hand zu kopieren
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "'jaira boards create <name> --type <typ>' legt ein Board mit den Lanes und der Beschreibung eines Board-Typs an (z.B. code, research, content), und jaira liefert ein paar Typen mit."
context: |-
  Berk am 18.10.: 'Im Moment haben wir nur 1 Board. Wir wollen mehrere Board-Typen ermoeglichen.'
  Heute: mehrere Boards pro Repo gibt es seit 0.5.3 (Ticket S0EG0M, done): Lane-Dateien und order-Datei nach .jaira/boards/<name>/lanes/, 'jaira create --board <name>', 'jaira boards' listet sie, 'b' im Board wechselt. Berk hat 0.5.2 installiert und sieht das deshalb noch nicht.
  Was fehlt: kein Befehl legt ein Board an ('jaira boards' hat keine Unterbefehle, nur Auflisten). Man muss die Dateien von Hand anlegen. Es gibt keine Board-Typen: nur die eingebauten Code-Lanes, 'jaira lanes default' und 'jaira lanes template --board'.
  Beispiel, dass es gebraucht wird: Berks Repo videoeditor hat ein Board mit eigenen Lanes (Reel Research, Creative Direction, Shoot, Edit ...), von Hand gebaut.
  Verwandt: PP24H1 (Interview-Skill soll Boards passend anlegen - braucht diesen Befehl), YA7KW7 (eine Bitte ueber zwei Boards), BGZ0SW (Ticket-Schema aus Lanes).
  Offen fuer brainstorm: welche Typen mitgeliefert werden; ob ein Typ nur Lanes + Beschreibung ist oder auch Tags/Milestone-Vorlagen; ob eigene Typen im Katalog (~/.jaira) liegen koennen. Kleiner als paca bleiben.
definition-of-done: "'jaira boards create <name> --type <typ>' schreibt .jaira/boards/<name>/lanes/ mit Lanes, order und description des Typs; ein vorhandenes Board wird nicht ueberschrieben (exit 3)"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-18T00:06:22Z
updated-at: 2026-10-18T00:32:01Z
updated-by: berkc
---

# Ein zweites Board laesst sich mit einem Befehl aus einem Board-Typ anlegen statt Lane-Dateien von Hand zu kopieren

## Definition of Done

- [ ] 'jaira boards create <name> --type <typ>' schreibt .jaira/boards/<name>/lanes/ mit Lanes, order und description des Typs; ein vorhandenes Board wird nicht ueberschrieben (exit 3)
- [ ] 'jaira boards types' listet die mitgelieferten Typen mit einer Zeile Beschreibung
- [ ] Ein Board ohne --type bekommt die eingebauten Code-Lanes
- [ ] Tests; Release-Note

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

