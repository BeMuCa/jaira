---
id: 01M84KJD20SF7R9EAADAES625B
title: "Alles Geteilte auf Refs: Team-Version, Tags, Lanes, Milestones, Logbuch - im Repo nur der Ticket-Schnappschuss im PR"
status: ideas
ready: true
creator: Alexander Sacharov
goal: "Zwischen Teammitgliedern wird nur über refs/jaira/* geteilt; Dateien unter .jaira/ im Repo gibt es nur noch als Schnappschuss des Tickets im PR. Ein Versions-Wächter auf einem Ref hält alle auf derselben jaira-Version."
context: |-
  Was heute falsch ist: jaira teilt Zustand auf zwei Wegen. Tickets, Epics und Milestones liegen auf refs/jaira/*; Lanes (.jaira/lanes/*.md), Tags (.jaira/tags), Logbuch (.jaira/logbook/, über 200 Dateien) und Archiv liegen als Dateien im Repo, Milestones zusätzlich als lokale Datei.
  Folge 1: 'jaira update' einer neueren Version schreibt .jaira/lanes, AGENTS.md und CLAUDE.md neu. Am 23.11.2026 blockierte dieser Rest in fast jedem Worktree den signoff-Gate, und jeder Worker musste fragen.
  Folge 2: hat einer jaira aktualisiert und der andere nicht, arbeiten beide mit verschiedenen Lanes und Regeln, ohne es zu merken.
  Folge 3: das Logbuch lässt das Repo wachsen und erzeugt Merge-Konflikte.
  Entscheidung des Maintainers (23.11.2026): alles Geteilte auf Refs, auch Tags und Lanes; im Repo nur der Ticket-Text, der im PR neben dem Code stehen soll.
  Versions-Wächter (Idee des Maintainers): refs/jaira/config hält die Team-Version. Beim Start liest jaira sie: Ref neuer -> jaira aktualisiert sich selbst und zeigt die Release-Notes; eigenes Binary neuer und ein echtes Release (kein Dev-Build) -> jaira hebt den Ref an; Ref nicht erreichbar -> letzte bekannte Version. Notbremse 'jaira hold <version>' für ein kaputtes Release. Aktualisieren nur beim Start, nie während Worker laufen.
  Bekanntes Risiko: Lanes auf einem Ref ändern sich ohne PR-Review und erreichen alle sofort; der Versions-Wächter sorgt dafür, dass alle Binaries die neuen Felder verstehen. GitHub zeigt refs/jaira nicht an - Logbuch nur noch mit jaira lesbar.
  Konflikt mit der heutigen CLAUDE.md-Regel 'vor self upgrade den Menschen fragen' - wird zu einer Einstellung (auto-upgrade).
  Diskussionsstand auf dem Whiteboard 'jaira: Begriffe ordnen' (Abschnitt 'Wo liegt was?').
definition-of-done: "Konzept entschieden: welche Refs (config, tags, lanes, logbook, archive), Format je Ref, Migration bestehender Dateien"
tags:
  - release
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-11-23T16:24:33Z
updated-at: 2026-11-23T16:24:56Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
round: 2
back-from: backlog
back-by: person
back-count: 1
---

# Alles Geteilte auf Refs: Team-Version, Tags, Lanes, Milestones, Logbuch - im Repo nur der Ticket-Schnappschuss im PR

## Definition of Done

- [ ] Konzept entschieden: welche Refs (config, tags, lanes, logbook, archive), Format je Ref, Migration bestehender Dateien
- [ ] Versions-Wächter: Lesen beim Start, Selbst-Update, Anheben nur durch Releases, hold, Offline-Fallback
- [ ] Tags, Lanes (nur Abweichungen vom eingebauten Stand), Milestones, Logbuch und Archiv lesen und schreiben über Refs; .jaira/ im Repo enthält nur noch Ticket-Schnappschüsse
- [ ] jaira update schreibt keine Dateien mehr ins Repo
- [ ] Migration: ein Befehl hebt ein bestehendes Board auf Refs, ohne Verlust
- [ ] Release-Note in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

