---
id: 01M2QFVJ5M1PSCX422VPD8CSAA
title: 0.3.0 schneiden
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer 'jaira update' nach 0.3.0 laufen laesst, liest vier Aenderungen und nicht zweiundzwanzig Zeilen Entwicklungsweg - und die Version ist geschnitten, sobald die vier Tickets in signoff angenommen sind"
context: |-
  Stand am 2026-09-17: master ist gleich upstream/master, letzter Tag ist v0.2.1. Naechste Version ist 0.3.0, weil Milestones eine Funktion sind und kein Patch.

  Was drin ist - vier Aenderungen: Milestones ganz (0YGWXQ), eine Lane ohne Code-Aenderung committet nicht mehr und der Handle im Betreff wird zur Bedingung (9ZZSFT), spawn.sh kennt die Lane 'dispatch' und --no-worktree (7KX89C), und jaira-role-pr macht den Pull Request auf, wenn ein Mensch die Rolle ruft.

  Das Problem, das vor dem Schnitt behoben gehoert: unter '## Unreleased' stehen 22 Zeilen, davon 18 zu Milestones. Mehrere davon beschreiben DASSELBE Verhalten in verschiedenen Entwicklungsstaenden - erst 'create weist den Namen ab', dann 'folge dem Ref in den Baum, der abgelegt hat', dann 'nein, lies den Abweisungstext als Hinweis auf dein eigenes Logbuch'. Ein Benutzer hat die Zwischenstaende nie gesehen; fuer ihn ist das ein Verhalten, nicht drei.

  Warum das zaehlt und nicht Kosmetik ist: core/release/NOTES.md ist mit go:embed im Binary und ist genau das, was 'jaira update' jemandem vorliest, dessen Board ein aelterer Build angefasst hat. Zweiundzwanzig Zeilen fuer vier Aenderungen ist das, was dieser Mensch bekommt.

  Die Sektion ist nicht getaggt, also ist sie offen und darf umgeschrieben werden - CLAUDE.md verbietet nur, eine bereits getaggte Sektion nachtraeglich wachsen zu lassen.

  Blockiert vom Menschen: 9ZZSFT, 0YGWXQ, 7KX89C und GTQHNH stehen in signoff. Ohne ihre Annahme wird nichts geschnitten. Das Falten der Notizen haengt nicht daran und kann sofort laufen.
definition-of-done: "Unter '## Unreleased' steht je eine Zeile je Aenderung, die ein Benutzer von aussen merkt - kein Zwischenstand, den nie jemand ausgeliefert bekommen hat, und keine zwei Zeilen fuer ein Verhalten. Die 18 Milestone-Zeilen sind zusammengefasst; was jede Zeile den Leser TUN laesst, bleibt erhalten."
tags:
  - release
blocked-by: []
related: []
commits: []
created-at: 2026-09-17T10:51:02Z
updated-at: 2026-09-17T11:18:33Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-33921
claimed-at: 2026-09-17T11:14:50Z
mode: conversational
---

# 0.3.0 schneiden

## Definition of Done

- [x] Unter '## Unreleased' steht je eine Zeile je Aenderung, die ein Benutzer von aussen merkt - kein Zwischenstand, den nie jemand ausgeliefert bekommen hat, und keine zwei Zeilen fuer ein Verhalten. Die 18 Milestone-Zeilen sind zusammengefasst; was jede Zeile den Leser TUN laesst, bleibt erhalten.
  proof: core/release/NOTES.md:21-25 — 18 Milestone-Zeilen zu 5 gefaltet, '## Unreleased' traegt jetzt 9 statt 22 Zeilen
- [ ] Das Zeilenformat ueberlebt das Falten: jede Aenderung ist genau eine Zeile, die mit '- ' beginnt, keine ist ueber zwei Zeilen umgebrochen. Nachgewiesen mit dem Test, der core/release/NOTES.md scannt.
- [ ] Keine Aenderung geht beim Falten verloren: jede der vier Aenderungen ist im gefalteten Text nachweisbar, und die Kommandos und Tasten, die ein Benutzer tippt ('jaira milestone', 'M', '--milestone', '--no-worktree', 'dispatch', 'jaira restore'), stehen weiter darin.
- [ ] Die Umbenennung von '## Unreleased' nach '## 0.3.0' und die frische leere '## Unreleased' darueber faehren in dem Commit, den der Mensch taggt - nicht frueher. Der Tag selbst wird von einem Menschen gesetzt, nicht von einem Agenten.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Die 18 Milestone-Zeilen (NOTES.md 21-38) auf 5 Zeilen falten, je eine pro Faehigkeit nach der Entscheidung vom 2026-09-17 11:14
- [x] Zeilenformat pruefen: jede Aenderung genau eine '- '-Zeile, kein Umbruch - go test ./core/release/
- [ ] Nachweisen, dass kein Kommando und keine Taste verloren ging: jaira milestone, M, --milestone, --no-worktree, dispatch, jaira restore
- [ ] Umbenennung nach '## 0.3.0' NICHT hier - sie gehoert in den Commit, den der Mensch taggt, nach Annahme der vier signoff-Tickets

## Progress
- **2026-09-17 11:14 · Alexander Sacharov** — Entscheidung vom Menschen, vor der Plan-Lane: die 18 Milestone-Zeilen werden auf ca. 5 Zeilen gefaltet, je eine pro Faehigkeit, die ein Benutzer anders benutzt - (1) Milestone anlegen/aendern per CLI und per Hand in .jaira/milestones/<name>.md, (2) Board und 'jaira list' auf eine Runde einschraenken (--milestone, Taste M), (3) Milestone-Farben am RECHTEN Kartenrand lesen, (4) Milestone reist auf seinem eigenen Ref und kommt mit 'jaira fetch', (5) Milestone mit 'jaira logbook <name>' ablegen und mit 'jaira restore <name>.md' zurueckholen. Alle Zwischenstaende und Corner-Case-Abweisungen (drei Fassungen der create-Abweisung, Locking, verlorene Rennen) fallen weg - ein Benutzer hat sie nie gesehen. Nicht 4 Zeilen (eine Milestone-Zeile wuerde zu lang und Kommandos ertrinken darin), nicht 8-10. Die anderen drei Aenderungen bleiben je eine Zeile. DoD 3 bleibt gewahrt: jaira milestone, M, --milestone, --no-worktree, dispatch, jaira restore stehen weiter im Text.
