---
id: 01M84KMZD7A1JXCHKJTGJ9H38M
title: "Git-Mittel für den Ref-Sync prüfen: atomic push, git notes, Hooks, Schutz von refs/jaira"
status: ideas
ready: true
creator: Alexander Sacharov
goal: "Entscheiden, welche eingebauten Git-Mittel der Ref-Sync aus ES625B nutzt, damit Teilen über refs/jaira/* sicher, automatisch und ohne eigenen Server läuft."
context: |-
  Anlass: ES625B will alles Geteilte (Team-Version, Tags, Lanes, Milestones, Logbuch) auf refs/jaira/* legen. Der Maintainer fragte am 23.11.2026, welche Git-Funktionen dafür noch helfen.
  Was jaira heute schon nutzt: eigene Refs refs/jaira/{tickets,epics,milestones}; Compare-and-swap beim Push über --force-with-lease (core/gitref/gitref.go); einen Merge-Treiber für Ticket-Dateien und union für .jaira/tags (.jaira/.gitattributes).
  Was es noch nicht nutzt, Kandidaten:
  - git push --atomic: mehrere Refs (z.B. Ticket + Logbuch + config) entweder alle oder keiner - kein halber Zustand.
  - git notes (refs/notes/jaira): Daten an einen Commit hängen ohne ihn zu ändern - z.B. Ticket-ID, bughunt-Ergebnis, Review-Urteil pro Commit. GitHub zeigt Notes nicht mehr an.
  - Hooks: pre-push oder post-merge, die refs/jaira mitschieben/holen, damit niemand 'jaira fetch' vergisst; reference-transaction-Hook lokal.
  - remote.origin.fetch-Refspec +refs/jaira/*:refs/jaira/* in der Repo-Config, damit ein normales git fetch die Refs mitbringt.
  - git bundle: Board offline als eine Datei weitergeben.
  Risiko: GitHub-Rulesets schützen nur Branches und Tags, keine eigenen Refs - jeder mit Push-Recht kann refs/jaira löschen oder zurücksetzen. Ein Schutz (Backup-Ref, Prüfsumme, Warnung bei Rücksprung) gehört dazu.
  Die Ref-Pfade in diesem Ticket sind Vorschläge, endgültig in ES625B.
definition-of-done: "Für jeden Kandidaten (atomic push, notes, Hooks, Refspec, bundle) steht im Ticket: nutzen oder nicht, mit Grund"
tags:
  - release
  - cli
blocked-by: []
related: []
follows: 01M84KJD20SF7R9EAADAES625B
commits: []
created-at: 2026-11-23T16:25:57Z
updated-at: 2026-11-23T16:26:13Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
round: 2
back-from: backlog
back-by: person
back-count: 1
---

# Git-Mittel für den Ref-Sync prüfen: atomic push, git notes, Hooks, Schutz von refs/jaira

## Definition of Done

- [ ] Für jeden Kandidaten (atomic push, notes, Hooks, Refspec, bundle) steht im Ticket: nutzen oder nicht, mit Grund
- [ ] Antwort auf das Risiko 'eigene Refs sind auf GitHub ungeschützt' steht fest
- [ ] Die Entscheidungen sind als DoD-Punkte in ES625B übernommen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

