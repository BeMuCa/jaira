---
id: 01M86NG0VN5W1KWA5B3ZMYMCDF
title: jaira spec find und Spec-Gates in jaira validate
status: backlog
ready: true
creator: Sascha
assignee: Sascha
goal: "Agenten finden ein Requirement ohne Lesen ganzer Dateien (jaira spec find <woerter>), jaira create warnt bei einem aehnlichen Requirement oder Ticket, und jaira validate prueft die Spec-Regeln ohne Modell."
context: |-
  Teil von SDD (siehe 7YQ1JC). Ziel: Token sparen und Doppelte vermeiden.
  jaira spec find <woerter>: Volltext ueber die EARS-Zeilen in docs/specs/*.md, Ausgabe ID, Zeile, Tickets; --json.
  jaira create ruft es mit dem Titel auf und zeigt die besten Treffer ('aehnlich: COCK-03 [W80HNH]').
  jaira validate: Spec-Datei <=100 Zeilen, jede Requirement-Zeile EARS mit ID, IDs eindeutig, genannte Tickets existieren, ADR <=40 Zeilen mit Status, kein ADR Proposed aelter als 14 Tage (Warnung).
  Kein Modell, reines Go.
definition-of-done: "jaira spec find <woerter> [--json] liefert passende Requirements mit ID, Zeile, Datei, Tickets; Test"
tags:
  - cli
blocked-by: []
related: []
follows: 01M86NFX4NQF351XWAXZ7YQ1JC
commits: []
created-at: 2026-11-24T11:36:41Z
updated-at: 2026-11-24T11:37:56Z
updated-by: Sascha
---

# jaira spec find und Spec-Gates in jaira validate

## Definition of Done

- [ ] jaira spec find <woerter> [--json] liefert passende Requirements mit ID, Zeile, Datei, Tickets; Test
- [ ] jaira create zeigt bis zu 3 aehnliche Requirements/Tickets als Hinweis (nicht blockierend); Test
- [ ] jaira validate prueft die Spec- und ADR-Regeln (Zeilenlimits, EARS-Form, ID eindeutig, Tickets existieren, alte Proposed als Warnung); Test
- [ ] Release-Note in core/release/notes/<ID>.md
- [ ] Nennt eine 'decision:'-Notiz eines Tickets eine Requirement-ID (z.B. COCK-03), prueft jaira validate (bzw. das signoff-Gate), dass diese ID im Diff der Ticket-Branch in docs/specs/ geaendert wurde; sonst Fehler mit ID und Datei; Test

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

