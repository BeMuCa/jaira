---
id: 01M2K275PSXDASRN6Y1VN71NVG
title: "Ein fertiger Milestone geht ins Archiv, und ein falscher laesst sich loeschen"
status: backlog
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer eine Runde Arbeit abgeschlossen hat, nimmt ihren Milestone in einem Griff vom Board - Datei und Ref weg, Inhalt weiter lesbar - und wer sich vertan hat, loescht ihn ganz."
context: |-
  Milestones entstehen in 0YGWXQ (.jaira/milestones/<name>.md plus refs/jaira/milestones/<name>). Es gibt keinen Weg, einen wieder wegzubekommen.

  Was heute fehlt, beides:
  - ARCHIVIEREN. Ein Milestone ist fertig, wenn seine Tickets fertig sind. Danach steht er weiter im M-Picker und faerbt weiter Kartenkanten. Nach zehn Runden Arbeit sind zehn tote Milestones im Picker und der Picker ist unbrauchbar - genau das Problem, das jaira logbook fuer Tickets schon geloest hat (.jaira/logbook/<wer>-<datum>/, 'jaira restore' bringt zurueck).
  - LOESCHEN. 'jaira milestone delete' existiert nicht. Die critique-Lane von 0YGWXQ hat den Loeschweg in Runde 1 als toten Code entfernt (RecordMilestoneDelete, PendingMilestone, gitref.DeleteMilestone hatten keinen Aufrufer) - richtig entfernt, aber der Befehl, der sie haette aufrufen sollen, wurde nie geschrieben.

  Ausgeloest am 2026-09-15: Alex hat einen Milestone 'test' angelegt und wollte ihn wieder weghaben. Es ging nur von Hand: rm der Datei, 'git update-ref -d refs/jaira/milestones/test' und ein Push mit Doppelpunkt gegen upstream. Drei Schritte, von denen zwei git-Wissen brauchen, fuer etwas, das ein Befehl sein muss.

  Der Ref ist der Teil, den man dabei vergisst. Loescht jemand nur die Datei, kommt der Milestone beim naechsten 'jaira fetch' vom Ref zurueck.

  Vorbild ist da und soll benutzt werden, nicht neu erfunden: jaira logbook fuer Tickets. Offen und Sache der Plan-Lane: ob das Archiv derselbe Ordner ist oder ein eigener, und was passiert, wenn ein Milestone archiviert wird, dessen Tickets noch offen sind.
definition-of-done: "Ein fertiger Milestone laesst sich in einem Griff archivieren: Datei und Ref verschwinden vom Board, der Inhalt bleibt lesbar - nach dem Vorbild von jaira logbook"
tags:
  - cli
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-09-15T17:35:45Z
updated-at: 2026-09-15T17:36:13Z
updated-by: Alexander Sacharov
---

# Ein fertiger Milestone geht ins Archiv, und ein falscher laesst sich loeschen

## Definition of Done

- [ ] Ein fertiger Milestone laesst sich in einem Griff archivieren: Datei und Ref verschwinden vom Board, der Inhalt bleibt lesbar - nach dem Vorbild von jaira logbook
- [ ] Ein archivierter Milestone faerbt keine Kartenkante mehr und steht nicht im M-Picker
- [ ] Es gibt einen Weg zurueck, wie 'jaira restore' ihn fuer Tickets hat
- [ ] Archivieren sagt es, wenn der Milestone noch unerledigte Tickets fuehrt, und archiviert die Tickets nicht mit
- [ ] 'jaira milestone delete <name>' entfernt Datei UND Ref, lokal und auf dem Remote - sonst kommt der Milestone beim naechsten fetch zurueck
- [ ] Je eine Zeile in core/release/NOTES.md unter ## Unreleased fuer das, was ein Benutzer davon merkt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

