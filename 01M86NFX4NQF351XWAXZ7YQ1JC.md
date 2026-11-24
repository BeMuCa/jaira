---
id: 01M86NFX4NQF351XWAXZ7YQ1JC
title: "SDD-Basis: Specs in EARS (max. 100 Zeilen je Datei) und ADRs aus allen Tickets erzeugen"
status: backlog
ready: true
creator: Sascha
assignee: Sascha
goal: "docs/specs/ beschreibt in EARS, was jaira heute tut (je Datei hoechstens 100 Zeilen, je Zeile ein Requirement mit ID und Tickets), docs/adr/ haelt die Entscheidungen (je ADR hoechstens 40 Zeilen); dazu ein Bericht ueber doppelte Tickets und widerspruechliche Entscheidungen."
context: |-
  Warum: Sascha und Berk arbeiten viel parallel, reden wenig; sie machen Dinge doppelt (z.B. X2M486 = AA2C3G) und widersprechen sich (Road vs Track fuer Swimlanes). Es fehlt eine gemeinsame Startbasis.
  Entscheidung Sascha 24.11.2026: Spec Driven Development + ADR. Menschen schreiben keine Specs, sie bestaetigen sie nur. Format EARS (The ... shall / When / While / If ... then / Where), kurz: ein Mensch liest hoechstens 100 Zeilen je Datei.
  Quellen: .jaira/logbook/** (~250 Tickets mit goal, DoD, outcome, 'decision:'-Notizen), .jaira/archive (abgelehnt -> ADR Rejected), offene Tickets (-> ADR Proposed), CLAUDE.md, docs/, Code als Gegenprobe.
  Gebaut einmalig per Claude-Workflow (mehrere Agenten je Bereich), dann PR; Sascha und Berk lesen und bestaetigen.
definition-of-done: "docs/specs/README.md indexiert alle Spec-Dateien (eine Zeile je Datei, Owner); jede Spec-Datei hat hoechstens 100 Zeilen, jede Requirement-Zeile ist EARS mit ID und [Tickets]"
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-11-24T11:36:37Z
updated-at: 2026-11-24T11:36:56Z
updated-by: Sascha
---

# SDD-Basis: Specs in EARS (max. 100 Zeilen je Datei) und ADRs aus allen Tickets erzeugen

## Definition of Done

- [ ] docs/specs/README.md indexiert alle Spec-Dateien (eine Zeile je Datei, Owner); jede Spec-Datei hat hoechstens 100 Zeilen, jede Requirement-Zeile ist EARS mit ID und [Tickets]
- [ ] docs/adr/NNNN-*.md: je ADR hoechstens 40 Zeilen mit Status, Context, Decision, Rejected, Sources (Ticket-IDs)
- [ ] Bericht ueber doppelte Tickets und widerspruechliche Entscheidungen als Notiz am Ticket; jeder Widerspruch ist ein ADR mit Status Proposed
- [ ] Gegenprobe: jede Requirement-Zeile stimmt mit dem heutigen Code ueberein (stichprobenhaft je Bereich, im Workflow geprueft)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

