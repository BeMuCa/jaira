---
id: 01M55S97BCFRNZMARY8WV38STG
title: "Die Installation per curl richtet Rollen und Cockpit gleich mit ein und zeigt dabei, was sie tut"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Wer jaira mit dem curl-Einzeiler (scripts/install.sh) oder install.ps1 installiert, hat danach Binary, Rollen-Skills und Cockpit eingerichtet, ohne einen zweiten Befehl, und sieht dabei Fortschritt: Download-Balken, 'Setting up the skills', 'Setting up the Cockpit'."
context: |-
  Was heute ist: scripts/install.sh (Zeile ~86-94) und scripts/install.ps1 hoeren auf, sobald das Binary liegt. Rollen und Cockpit kommen erst mit 'jaira roles install --global', und das sagt niemandem jemand.
  Ausloeser: Berk am 17.10.: hat per curl installiert, im Repo 'jaira' getippt, Board war da - 'jaira init' brauchte er nie. Rollen und Cockpit fehlten aber. Er will KEIN Onboarding-Dialog, sondern: der curl-Install installiert und aktiviert Rollen und Cockpit gleich mit.
  Bekannt: 'jaira roles install --global' schreibt 34 Dateien nach ~/.claude/skills, inklusive Cockpit (am 17.10. mit 0.5.2 laufen lassen: zweiter Lauf '0 written, 34 unchanged', exit 0). Das Cockpit laedt erst in der naechsten Claude-Code-Session.
  Gewuenscht sind sichtbare Schritte: Download mit Fortschrittsbalken (curl --progress-bar o.ae.), dann Zeilen wie 'Setting up the skills', 'Setting up the Cockpit', am Ende was jetzt aktiv ist.
  Offen: ein Opt-out (Umgebungsvariable, z.B. JAIRA_NO_ROLES=1) fuer CI und Leute ohne Claude Code; nicht-interaktiv lauffaehig bleiben (curl | sh hat kein TTY auf stdin).
definition-of-done: |-
  install.sh installiert nach dem Binary die Rollen und das Cockpit global, mit Fortschrittsbalken beim Download und je einer Statuszeile fuer Skills und Cockpit
  install.ps1 tut dasselbe
  Ein Opt-out per Umgebungsvariable laesst nur das Binary installieren
  Ein zweiter Lauf ist idempotent und meldet unveraenderte Dateien statt Fehler
  Release-Note in core/release/notes/<ID>.md
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-17T20:36:33Z
updated-at: 2026-10-17T20:37:24Z
updated-by: berkc
---

# Die Installation per curl richtet Rollen und Cockpit gleich mit ein und zeigt dabei, was sie tut

## Definition of Done

- [ ] install.sh installiert nach dem Binary die Rollen und das Cockpit global, mit Fortschrittsbalken beim Download und je einer Statuszeile fuer Skills und Cockpit
install.ps1 tut dasselbe
Ein Opt-out per Umgebungsvariable laesst nur das Binary installieren
Ein zweiter Lauf ist idempotent und meldet unveraenderte Dateien statt Fehler
Release-Note in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-17 20:37 · berkc** — berkc took this ticket over from BeMuCa
