---
id: 01M2FYEHZYFCW7DSNRHY9ET6NC
title: "Der Remote-Name gilt pro Rechner, gebraucht wird er pro Board"
status: backlog
ready: false
creator: Alexander Sacharov
goal: "Auf einem Board, dessen Repository den eingestellten Remote nicht hat, funktionieren die ref-Befehle wieder - ohne dass ein Ticket dadurch im falschen Repository landet."
context: |-
  'jaira release' ist auf jedem Board ausser diesem tot. Es bricht ab mit: no remote "upstream".

  Damit fehlt genau die Haelfte des Mechanismus, die ein Mensch braucht: ein Ticket nehmen und, wenn er es doch nicht macht, wieder hergeben. Alex ist am 2026-09-14 auf dem requirementsgenie-Board bei FAQ3MW darauf gelaufen und hat den assignee von Hand aus der Datei geloescht.

  Die Ursache ist eine Einstellung, die auf der falschen Ebene liegt:
  - ~/.jaira/settings.json enthaelt {"remote": "upstream"}. Diese Datei gilt pro Rechner und damit fuer JEDES Board, das auf ihm geoeffnet wird.
  - core/settings/settings.go:130 RemoteName() gibt den eingestellten Namen unbesehen zurueck. Auf den Vorgabewert origin faellt es nur zurueck, wenn der String leer ist.
  - core/gitref/gitref.go:112-117 Repo.remote() macht dasselbe noch einmal.
  - core/gitref/gitref.go:100-105 fuehrt dann 'git remote get-url upstream' aus und liefert ErrNoRepo: no remote %q.

  Welches Repository welche Remotes hat:
  - /home/alex/projects/jaira: origin = sashasoft90/jaira (Fork), upstream = BeMuCa/jaira. Hier stimmt die Einstellung.
  - /home/alex/projects/requirementsgenie: nur origin = git.esprit-engineering.de/.../requirementsgenie. Hier gibt es kein upstream, und jeder ref-Befehl faellt um.

  Die Einstellung kam am 2026-09-13 dazu, als das jaira-Board auf upstream umzog. Sie war fuer dieses eine Repository richtig und hat alle anderen mitgerissen.

  Wichtig fuer die Plan-Lane, sonst wird die naheliegende Loesung die falsche: 'wenn der eingestellte Remote fehlt, nimm origin' ist NICHT sicher. In genau diesem Repository ist origin der Fork. Ein Ticket-Ref, das still nach origin geht statt nach upstream, landet im Fork - und das ist der Fehler, gegen den die Einstellung ueberhaupt eingefuehrt wurde. Ein stiller Rueckfall tauscht einen lauten Abbruch gegen einen leisen Datenverlust.

  Nicht Teil dieses Tickets: wohin die Einstellung genau wandert (Board-Datei, Eintrag je Pfad in ~/.jaira/settings.json, oder etwas Drittes). Das entscheidet die Plan-Lane; die beiden Bedingungen unten gelten fuer jeden dieser Wege.
definition-of-done: |-
  In einem Repository, dessen einziger Remote origin heisst, laeuft 'jaira release <id>' durch, waehrend ~/.jaira/settings.json weiterhin remote: upstream sagt. Nachgestellt an einem Fixture-Repository mit genau einem Remote.

  Im jaira-Repository selbst (origin = Fork, upstream = BeMuCa) geht der Ticket-Ref weiterhin nach upstream. Ein Test haelt das fest, damit die Loesung nicht darin bestehen kann, ueberall still auf origin auszuweichen.

  Der Remote laesst sich pro Board festlegen, nicht nur pro Rechner - auf welchem Weg auch immer die Plan-Lane das loest.

  Bricht eine ref-Operation doch am Remote ab, nennt die Meldung drei Dinge: den eingestellten Namen, die Remotes, die dieses Repository tatsaechlich hat, und den Befehl, der es geradezieht. Nicht nur 'no remote "upstream"'.

  Eine Zeile in core/release/NOTES.md unter ## Unreleased.
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-14T12:32:09Z
updated-at: 2026-09-14T12:36:43Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-93721
claimed-at: 2026-09-14T12:36:43Z
---

# Der Remote-Name gilt pro Rechner, gebraucht wird er pro Board

## Definition of Done

- [ ] In einem Repository, dessen einziger Remote origin heisst, laeuft 'jaira release <id>' durch, waehrend ~/.jaira/settings.json weiterhin remote: upstream sagt. Nachgestellt an einem Fixture-Repository mit genau einem Remote.

Im jaira-Repository selbst (origin = Fork, upstream = BeMuCa) geht der Ticket-Ref weiterhin nach upstream. Ein Test haelt das fest, damit die Loesung nicht darin bestehen kann, ueberall still auf origin auszuweichen.

Der Remote laesst sich pro Board festlegen, nicht nur pro Rechner - auf welchem Weg auch immer die Plan-Lane das loest.

Bricht eine ref-Operation doch am Remote ab, nennt die Meldung drei Dinge: den eingestellten Namen, die Remotes, die dieses Repository tatsaechlich hat, und den Befehl, der es geradezieht. Nicht nur 'no remote "upstream"'.

Eine Zeile in core/release/NOTES.md unter ## Unreleased.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

