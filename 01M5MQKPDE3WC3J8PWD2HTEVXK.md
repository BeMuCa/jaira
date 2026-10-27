---
id: 01M5MQKPDE3WC3J8PWD2HTEVXK
title: "Teamlead: Release-Weg je Repo als Werkzeuge im Skill"
status: backlog
ready: false
creator: Alexander Sacharov
goal: "Eine frische Teamlead-Session released ein fertiges Ticket ohne Rückfragen zum Weg: prüfen, Changelog/Version, MR+Merge+Tag mit Wiederholung, Abnahme-Block, Aufräumen, Stand neu bauen — als Werkzeuge, je Repo konfiguriert."
context: |-
  Gefunden am 23.10.2026 in einer langen Teamlead-Session für ein anderes Projekt (Notizen: ~/.claude/acceptance/<projekt>/teamlead-staerken.md).
  Heute steht der Release-Weg nur im Kopf einer Session. In einem anderen Projekt wurden dafür sechs Skripte gebaut (rel_check.sh, cl.py, rel_ship.sh, addblock.py, cleanup.sh, warm.sh unter ~/.claude/acceptance/<projekt>/bin) — vorher lagen sie im /tmp-Scratchpad und waren nach einem Neustart weg.
  Gelernte Stolperstellen, die in die Werkzeuge gehören: MR-Anlage direkt nach dem Push schlägt manchmal fehl (Wiederholung); Merge direkt nach Anlage ebenso; Changelog-Text mit Anführungszeichen bricht den Build (escapen + Build nach dem Bump); Abnahme-Block-Nummern müssen über Step-IDs zählen; nach Rebuild das Einbettungsmodell vorwärmen.
  Teamlead-Skill beschreibt Lanes und Dispatch, aber nicht den Weg 'fertig → Stand → MR → Tag → Abnahme → aufräumen'.
  Vorschlag: Abschnitt im jaira-teamlead-Skill + Werkzeuge im Skill-Ordner, repo-spezifische Teile (Stand-worktree, Prüfbefehle, Changelog-Datei, Tag-Schema) in einer Repo-Einstellung.
definition-of-done: "Teamlead-Skill enthält den Release-Weg mit Werkzeugen, repo-spezifisch konfigurierbar"
tags:
  - release
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:55:55Z
updated-at: 2026-10-27T20:53:32Z
assignee: BeMuCa
updated-by: BeMuCa
---

# Teamlead: Release-Weg je Repo als Werkzeuge im Skill

## Definition of Done

- [ ] Teamlead-Skill enthält den Release-Weg mit Werkzeugen, repo-spezifisch konfigurierbar
- [ ] Werkzeuge liegen dauerhaft (nicht /tmp) und haben Wiederholung für MR/Merge
- [ ] Probe: frische Session released ein Ticket in einem anderen Projekt ohne Rückfragen zum Weg

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

