---
id: 01M257R8QBBG3613E8V0GDB5YJ
title: "Enter entfernt eine installierte Lane, so wie es sie hinzufuegt"
status: backlog
ready: true
creator: BeMuCa
assignee: Alexander Sacharov
goal: "Enter ist im Lane-Screen der eine Schalter: auf einer noch nicht installierten Lane fuegt es sie hinzu, auf einer installierten nimmt es sie wieder weg"
context: |-
  Berk am 09.09.: 'man sollte eine lane die hinzugefuegt wurde mit enter auch wieder entfernen'.
  WICHTIG - Entfernen gibt es schon, auf x. internal/tui/lanes.go:284-287 ruft startRemove, das einen Ja/Nein-Dialog aufmacht (lanes.go:347); die Fusszeile nennt 'x remove'. Es fehlt also keine Faehigkeit, es fehlt die Taste, die Berk erwartet.
  Der Konflikt, der vor dem Bauen zu entscheiden ist: enter ist heute schon belegt (lanes.go:288-293). Auf der Plus-Spalte oeffnet es den Katalog, auf einer noch nicht installierten Lane fuegt es sie hinzu. Beides ist 'enter fuegt hinzu'.
  Berks Wunsch macht daraus einen Umschalter: enter auf einer installierten Lane = entfernen. Das ist in sich stimmig, aber enter wird damit auf einer Spalte destruktiv, und enter ist die Taste, die man blind drueckt.
  Was dagegen absichert: der Ja/Nein-Dialog existiert schon und zeichnet 'no' zuerst, damit blindes Enter-Haemmern nichts loescht (lanes.go:810-820). Diese Vorsicht muss erhalten bleiben - genau sie ist der Grund, warum enter hier ueberhaupt vertretbar ist.
  x soll bleiben. Zwei Wege zur selben Handlung sind hier kein Widerspruch, sondern der Schutz fuer den, der x schon kennt.
definition-of-done: "Enter auf einer installierten Lane startet dieselbe Entfernung wie x, samt Ja/Nein-Dialog mit 'no' als Vorauswahl; enter auf der Plus-Spalte und auf einer nicht installierten Lane fuegt weiter hinzu; x funktioniert unveraendert; die Fusszeile nennt beide; Tests decken beide Enter-Faelle; go test ./... -race gruen"
tags:
  - tui
blocked-by: []
commits: []
created-at: 2026-09-10T08:43:06Z
updated-at: 2026-10-11T19:44:56Z
updated-by: Alexander Sacharov
---

# Enter entfernt eine installierte Lane, so wie es sie hinzufuegt

## Definition of Done

- [ ] Enter auf einer installierten Lane startet dieselbe Entfernung wie x, samt Ja/Nein-Dialog mit 'no' als Vorauswahl; enter auf der Plus-Spalte und auf einer nicht installierten Lane fuegt weiter hinzu; x funktioniert unveraendert; die Fusszeile nennt beide; Tests decken beide Enter-Faelle; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-11 19:44 · Alexander Sacharov** — Code aus dem A/B-Versuch vom 11.10.2026 uebernommen (Arm A, Lane-Pipeline; blinder Richter knapp fuer diese Variante; Review: approve). Offen laut Review: Doc-Kommentare in internal/tui/lanes.go um startRemove nennen teils nur x. Weiter im Worker-Modus; Alex hat das Ticket von BeMuCa uebernommen, um es zum Release zu bringen.
