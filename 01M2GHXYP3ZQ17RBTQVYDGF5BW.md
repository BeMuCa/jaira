---
id: 01M2GHXYP3ZQ17RBTQVYDGF5BW
title: "Ein Remote, den es gibt und der nicht funktioniert, wird wie ein fehlender gemeldet"
status: backlog
ready: false
creator: Alexander Sacharov
goal: "Ein Remote, der da ist aber nicht antwortet, wird als das gemeldet was er ist - nicht als fehlend, und nicht als funktionierend."
context: |-
  Die Review-Lane von VM0A76 hat es am 2026-09-14 von Hand nachgestellt und bewusst nicht behoben, weil es eine Formaenderung ist und kein Textfehler.

  Fall: der Remote ist konfiguriert und existiert, antwortet aber nicht - tote URL, kein Netz, keine Berechtigung.

  Was dann passiert, zwei Meldungen, beide falsch:
  - 'jaira create' schreibt 'no remote "origin" — this repository has origin'. Der Satz widerspricht sich in sich selbst und laesst den Leser ratlos zurueck.
  - 'jaira whoami' sagt 'Ref mode: yes', waehrend das create, das gerade lief, eine Datei geschrieben hat. Die Auskunft und das Verhalten widersprechen sich.

  Die Ursache ist, dass zwei verschiedene Zustaende in einen Fehler zusammenfallen: 'kein Remote konfiguriert' und 'Remote konfiguriert, aber unerreichbar'. Sie brauchen verschiedene Antworten - der erste ist mit 'git remote add' behoben, der zweite nicht.

  Kein Rueckschritt durch VM0A76: unerreichbare Remotes verhielten sich vorher auch schon schief. Die neuen Meldungen machen es nur zum ersten Mal sichtbar, weil sie ueberhaupt etwas ueber den Modus sagen.

  Warum es nicht in den Zweig von VM0A76 kam: der trug zu dem Zeitpunkt zwei Tickets und war fertig und gruen. Eine Formaenderung dort einzuziehen haette 0.2.1 aufgehalten.
definition-of-done: "Ein konfigurierter, aber unerreichbarer Remote erzeugt keine Meldung mehr, die sich selbst widerspricht: 'no remote \"X\" — this repository has X' kommt nicht mehr vor. Nachgestellt an einem Repository, dessen Remote auf eine tote URL zeigt."
tags:
  - cli
blocked-by: []
related:
  - 01M2G02EHTV6RJPPQKR5VM0A76
commits: []
created-at: 2026-09-14T18:12:37Z
updated-at: 2026-09-14T18:12:51Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
---

# Ein Remote, den es gibt und der nicht funktioniert, wird wie ein fehlender gemeldet

## Definition of Done

- [ ] Ein konfigurierter, aber unerreichbarer Remote erzeugt keine Meldung mehr, die sich selbst widerspricht: 'no remote "X" — this repository has X' kommt nicht mehr vor. Nachgestellt an einem Repository, dessen Remote auf eine tote URL zeigt.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

