---
id: 01M4930V33NCPD29SW0ZV3S0MW
title: Ein Ticket landet nur zusammen mit Code in einem Commit
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: Ein Ticketfile unter .jaira/tickets/ kommt nie ohne Codeänderung in einen Commit; Backlog-Tickets leben nur auf ihrem Ref.
context: |-
  Was falsch ist: auf master liegen 32 Ticketfiles unter .jaira/tickets/, fast alle im Status backlog. Sie gammeln dort, obwohl seit 19419ee (RFC7GA) jaira create ein Ticket nur auf refs/jaira/tickets/<id> ablegt und das File löscht.
  Ursache 1: der generierte Agent-Block sagt noch 'The one ticket that still earns a commit of its own is a ticket you create and hand to someone else — commit it' (core/board/announce.go:110, kopiert nach CLAUDE.md:107, AGENTS.md:80, docs/AGENTS.md:81). Das stammt aus der Zeit vor den Refs; Agenten folgen ihm und committen frisch erstellte Tickets.
  Ursache 2: nichts prüft die Regel. Ein Commit, der nur Ticketfiles trägt, geht durch.
  Regel (mit Alex abgestimmt): ein Ticketfile darf nur in einen Commit, der auch eine Datei außerhalb von .jaira/ ändert. Ausnahmen: jaira logbook / archive (Ticket verlässt tickets/), reines Löschen eines Ticketfiles.
  Ausnahme im Text: im File-Modus (kein nutzbares Remote, create meldet 'as a file on your disk') muss das Ticket weiter committet werden, sonst sieht es niemand.
  Nicht Teil davon: die 32 alten Files von master räumen (eigener Schritt, erst prüfen, welche Branches sie noch anfassen); ein eigenes jaira-Kommando für fremde Repos.
definition-of-done: "core/board/announce.go sagt: im Ref-Modus wird ein erstelltes Ticket nicht committet; nur im File-Modus"
tags:
  - ci
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-10-06T17:08:45Z
updated-at: 2026-10-06T17:12:45Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-22138
claimed-at: 2026-10-06T17:08:56Z
outcome-what: "Der Agent-Block sagt nicht mehr, ein erstelltes Ticket zu committen; CI lehnt Commits ab, die Ticketfiles ohne Code tragen (scripts/check-ticket-commits.sh)."
outcome-why: "Agenten folgten der alten Regel und committeten frische Tickets nach master, obwohl sie längst auf ihrem Ref leben; ohne Prüfung kam nichts davon auf."
outcome-resolves: Backlog-Tickets landen nicht mehr auf master; ein Ticket kommt nur noch über jaira pull und mit Code in einen Branch.
---

# Ein Ticket landet nur zusammen mit Code in einem Commit

## Definition of Done

- [x] core/board/announce.go sagt: im Ref-Modus wird ein erstelltes Ticket nicht committet; nur im File-Modus
  proof: core/board/announce.go:107
- [x] CLAUDE.md, AGENTS.md und docs/AGENTS.md sind aus dem neuen Text regeneriert
  proof: CLAUDE.md:107, AGENTS.md:80, docs/AGENTS.md:81
- [x] CI schlägt fehl, wenn ein Commit im PR ein Ticketfile unter .jaira/tickets/ hinzufügt oder ändert, ohne eine Datei außerhalb von .jaira/ zu ändern
  proof: .github/workflows/ci.yaml:11, scripts/check-ticket-commits.sh
- [x] core/release/NOTES.md hat eine Zeile unter Unreleased
  proof: core/release/NOTES.md:17

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-06 17:12 · Alexander Sacharov** — Die CI-Prüfung läuft nur auf pull_request über base..head. Die alte Historie verletzt die Regel dutzendfach (z.B. ec2c861, 4cf0249) — ein Check über die ganze Historie wäre für immer rot. Löschen eines Ticketfiles ist erlaubt, damit logbook/archive und das Aufräumen der 32 Altfiles auf master durchgehen. 13 Tests in internal/cli schlagen auf diesem Rechner auch ohne diese Änderung fehl; mit leerem HOME sind alle grün — Ursache ist ~/.jaira, nicht der Code.
