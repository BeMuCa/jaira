---
id: 01M5GFB0CEC5NVDXHQWSAMJB8G
title: "Pruefen, ob Tickets ganz auf Refs leben koennen statt als Dateien im Repo (Archiv inklusive)"
status: ideas
ready: true
creator: berkc
assignee: berkc
goal: "Begruendet entschieden, ob Ticket-, Archiv- und Logbook-Dateien im Repo bleiben oder nur noch auf refs/jaira/* liegen."
context: |-
  Frage von berkc am 19.10.2026: warum muss das Archivieren in git gepusht werden? Wenn wir mit Refs arbeiten, koennen die Tickets nicht ganz aus dem Repo raus, und jeder neue Nutzer liest nur die Refs?
  Stand (geprueft): jaira archive verschiebt die Datei nach .jaira/archive/ im Repo - das braucht einen Commit. Der Ref bleibt und bekommt den Endstand (core/refsync/refsync.go:191 RecordFiled = Record). Nur jaira delete loescht den Ref (recordDelete). GWK31K (Human Review) schreibt beim Archivieren zusaetzlich eine Zeile 'archived:' auf den Ref, damit jedes Board das Ticket sofort ausblendet - ohne Commit.
  Was Dateien im Repo heute leisten:
  - Das Ticket faehrt im PR mit (Regel in CLAUDE.md: 'the ticket rides with the code'): Reviewer sehen Aenderung und Zweck an einer Stelle.
  - GitHub zeigt Dateien, Refs zeigt es nicht.
  - Von Hand lesbar und editierbar (Constraint 'Storage: Markdown ... inside the repo').
  - git clone holt refs/jaira/* nicht von selbst; jaira holt sie mit eigener Refspec.
  Was nur-Refs braechte: kein 'pull it', kein Commit fuer Archiv/Logbook, keine Merge-Konflikte in Ticket-Dateien.
definition-of-done: "Entscheidungsnotiz: was Dateien im Repo leisten, was nur-Refs leistet, was dabei verloren ginge, mit Empfehlung"
tags:
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T00:14:24Z
updated-at: 2026-10-22T00:14:24Z
---

# Pruefen, ob Tickets ganz auf Refs leben koennen statt als Dateien im Repo (Archiv inklusive)

## Definition of Done

- [ ] Entscheidungsnotiz: was Dateien im Repo leisten, was nur-Refs leistet, was dabei verloren ginge, mit Empfehlung
- [ ] Falls nur-Refs: Folge-Tickets fuer Migration, PR-Ansicht des Tickets und Lesbarkeit auf GitHub angelegt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

