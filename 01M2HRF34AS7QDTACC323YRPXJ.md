---
id: 01M2HRF34AS7QDTACC323YRPXJ
title: "spawn.sh kann nur eine Rolle starten, obwohl der Prompt zwei verlangt"
status: backlog
ready: true
creator: Alexander Sacharov
goal: "Ein Dispatcher kann mit dem mitgelieferten Skript jede Lane besetzen, die sein eigener Prompt von ihm verlangt - nicht nur die, deren Name fest eingebaut ist."
context: "scripts/spawn.sh:71 schickt einen fest verdrahteten Befehl in die Vorlage:\n\n  \"$herdr\" pane send-text \"$pane\" \"/jaira-role-lane $ticket $lane\"\n\nDer Prompt des Dispatchers verlangt aber zwei verschiedene Rollen. jaira-dispatcher/SKILL.md:48 sagt woertlich: '/jaira-role-lane <id> <lane>. Testing is not a lane: /jaira-role-tester <id>.'\n\nDamit kann ein Dispatcher, der seinem eigenen Prompt folgt UND das mitgelieferte Skript benutzt, nicht beides tun. Genau diese Kombination verlangt KSGSKK von ihm.\n\nGefunden am 2026-09-15 von dem Dispatcher, der 13VMA8 fuehrte, beim Starten der testing-Lane - also von dem ersten, dessen Prompt ihn ueberhaupt auf spawn.sh verwies. Er hat mit lane=testing weitergemacht, was auf diesem Board vertretbar ist, weil testing hier eine eingerichtete Lane mit eigenem Prompt und einem verpflichtenden test-verdict ist. Auf einem Board ohne diese Lane haette er nichts starten koennen.\n\nZweiter Befund aus demselben Lauf, gleiche Datei: spawn.sh vertraut der Pane-Id, die 'tab create' zurueckgibt. Sie wurde stale, als ein zweiter Dispatcher im selben Moment seinen eigenen Tab aufmachte. Der betroffene Lauf hat die Pane ueber die Tab-Id neu aufgeloest und kam durch; das Skript kann das nicht.\n\nNicht Teil dieses Tickets: der WSL-Start (wsl.exe --cd), das ist in KSGSKK behoben und seit dem 2026-09-15 auch in der installierten Kopie."
definition-of-done: "scripts/spawn.sh startet die Rolle, die zur Lane gehoert: fuer testing '/jaira-role-tester <id>', sonst '/jaira-role-lane <id> <lane>'. Nachgestellt, indem beide Lanes einmal wirklich gestartet werden."
tags:
  - cli
blocked-by: []
related:
  - 01M2G9X5HVH29SDS8FAZKSGSKK
commits: []
created-at: 2026-09-15T05:26:04Z
updated-at: 2026-09-16T06:48:19Z
assignee: "Alexander Sacharov"
updated-by: Alexander Sacharov
---

# spawn.sh kann nur eine Rolle starten, obwohl der Prompt zwei verlangt

## Definition of Done

- [ ] scripts/spawn.sh startet die Rolle, die zur Lane gehoert: fuer testing '/jaira-role-tester <id>', sonst '/jaira-role-lane <id> <lane>'. Nachgestellt, indem beide Lanes einmal wirklich gestartet werden.
- [ ] Die Pane, in die getippt wird, wird nicht aus einem Wert von vorhin geglaubt: ein zweiter Dispatcher, der im selben Moment einen Tab aufmacht, laesst den Lauf nicht mehr in die falsche Vorlage schreiben. Nachgestellt mit zwei gleichzeitigen Starts.
- [ ] Faellt der Start doch um, sagt das Skript in welcher Vorlage und woran - statt einer Zeile, die nur 'claude did not come up' meldet.
- [ ] Eine Zeile in core/release/NOTES.md unter ## Unreleased, weil das Skript mit den Rollen ausgeliefert wird.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-16 06:48 · Alexander Sacharov** — Am 2026-09-15 und 16 dreimal live bestaetigt, von der Teamlead-Sitzung selbst - der Befund ist keine Lesefrucht mehr.

Ein Teamlead soll laut ~/.claude/skills/jaira-teamlead/SKILL.md den Dispatcher ueber scripts/spawn.sh starten. spawn.sh:96 tippt fest '/jaira-role-lane $ticket $lane'. Es gibt keinen Weg, damit '/jaira-dispatcher <id>' zu starten - also genau die Rolle, fuer die das Skript im Dispatcher-Skill liegt.

Was die Sitzung stattdessen getan hat, dreimal (9ZZSFT, 0YGWXQ zweimal): spawn.sh in ein Scratchpad kopiert und die eine Zeile per sed durch ${JAIRA_SPAWN_PROMPT:-...} ersetzt. Das lief, hinterlaesst aber eine zweite Kopie des Skripts, die niemand pflegt - dieselbe Krankheit, die 7MG5GB an einem fremden Board beschreibt.

Damit ist der Befund breiter als der Titel sagt. Es fehlt nicht nur '/jaira-role-tester': es fehlt jede Rolle ausser einer. Drei Aufrufer wollen drei verschiedene Zeilen - der Dispatcher '/jaira-role-lane' und '/jaira-role-tester' (SKILL.md:48), der Teamlead '/jaira-dispatcher'. Eine Loesung, die nur testing nachtraegt, laesst den Teamlead wieder mit sed arbeiten.

Vorschlag fuer die Plan-Lane, nicht entschieden: den Prompt zum Argument machen statt ihn aus der Lane abzuleiten. Ein Aufrufer, der nichts sagt, bekommt weiter '/jaira-role-lane <id> <lane>'; wer etwas anderes braucht, uebergibt es. Das deckt alle drei Aufrufer mit einer Zeile und braucht keine Liste von Sonderfaellen im Skript.

DoD 1 ist damit zu eng formuliert ('fuer testing ..., sonst ...'). Wer ihn woertlich baut, hat den Teamlead-Fall nicht geloest.
