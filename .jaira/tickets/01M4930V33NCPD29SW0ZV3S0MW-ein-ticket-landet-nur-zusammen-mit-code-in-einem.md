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
updated-at: 2026-10-06T19:05:07Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-22138
claimed-at: 2026-10-06T17:08:56Z
outcome-what: "Brainstorm-Rolle und Projekt-Skill sagen jetzt wie der Agent-Block: erstellte Tickets nicht committen, außer im File-Modus; NOTES-Zeile nennt jaira roles install --force."
outcome-why: "Review fand zwei Anweisungen, die der neuen Regel widersprachen; die Brainstorm-Rolle ist die wahrscheinlichste Quelle der Backlog-Stapel auf master."
outcome-resolves: "Keine mitgelieferte Anweisung schickt mehr ein frisch erstelltes Ticket in einen Commit (git grep 'earns a commit' leer)."
question: "Regel so richtig? Danach: die 32 Altfiles auf master räumen (vorher prüfen, welche Branches sie anfassen)."
review-summary: "Der generierte Agent-Block (core/board/announce.go, daraus CLAUDE.md/AGENTS.md) und docs/AGENTS.md sagen jetzt: ein erstelltes Ticket wird nicht committet, es reist auf seinem Ref und kommt per jaira pull in den Branch; nur im File-Modus committen. Neu ist scripts/check-ticket-commits.sh: es geht jeden Nicht-Merge-Commit von base..head durch und schlägt fehl, wenn einer ein Ticketfile unter .jaira/tickets/ hinzufügt oder ändert, ohne eine Datei außerhalb von .jaira/ zu ändern; Löschen ist frei. Der CI-Job ticket-commits läuft es auf jedem pull_request."
review-gaps: "Zwei Anweisungen sagen weiter das Gegenteil und werden jetzt von CI abgelehnt: (1) core/role/builtin/jaira-role-brainstorm/SKILL.md:77 'Then commit the ticket files' — die Brainstorm-Rolle legt Tickets an und committet sie; genau so landen Backlog-Tickets stapelweise auf master. (2) .claude/skills/jaira/SKILL.md:286 'The one ticket that still earns a commit of its own is one you create and hand to someone else: commit it'. Der DoD-Punkt 1 ist damit nur für den generierten Block erfüllt, das Ziel ('nie ohne Codeänderung') nicht."
test-verdict: "pass — scripts/check-ticket-commits.sh master auf dem Branch: OK. Fixture-Repo: Commit nur mit .jaira/tickets/X.md -> exit 1 mit Commit-Namen; Ticket+main.go -> durch; git mv nach .jaira/logbook/ -> durch. go test ./internal/cli ./core/... mit leerem HOME grün; die 13 Fehler mit echtem HOME bestehen identisch ohne diese Änderung (Ursache ~/.jaira)."
review-verdict: "nicht fertig: der Diff erfüllt die vier DoD-Punkte, aber zwei mitgelieferte Anweisungen widersprechen der Regel weiter, und eine davon (Brainstorm-Rolle) ist die wahrscheinlichste Quelle der Altfiles. Zurück nach in-progress."
review-check: "1. git switch fix/ticket-rides-with-code  2. scripts/check-ticket-commits.sh master -> keine Ausgabe, exit 0  3. git grep -n 'earns a commit' -> keine Treffer  4. grep -n 'commit the ticket files' core/role/builtin/jaira-role-brainstorm/SKILL.md -> keine Treffer  5. sed -n 100,115p CLAUDE.md -> der Absatz sagt 'A ticket you create commits nothing either'"
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
- [ ] core/role/builtin/jaira-role-brainstorm/SKILL.md sagt nicht mehr, die angelegten Tickets zu committen
- [ ] .claude/skills/jaira/SKILL.md sagt dasselbe wie der Agent-Block: ein erstelltes Ticket wird nur im File-Modus committet

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-06 17:12 · Alexander Sacharov** — Die CI-Prüfung läuft nur auf pull_request über base..head. Die alte Historie verletzt die Regel dutzendfach (z.B. ec2c861, 4cf0249) — ein Check über die ganze Historie wäre für immer rot. Löschen eines Ticketfiles ist erlaubt, damit logbook/archive und das Aufräumen der 32 Altfiles auf master durchgehen. 13 Tests in internal/cli schlagen auf diesem Rechner auch ohne diese Änderung fehl; mit leerem HOME sind alle grün — Ursache ist ~/.jaira, nicht der Code.
- **2026-10-06 18:57 · Alexander Sacharov** — Alex am 2026-10-06: Regel und CI-Prüfung angenommen wie gebaut (Ticketfile nur mit Code außerhalb .jaira/). Nächster Schritt danach: die 32 Altfiles auf master räumen.
- **2026-10-06 19:03 · Alexander Sacharov** — Review 2026-10-06: zwei Anweisungen widersprechen der Regel noch — Brainstorm-Rolle (core/role/builtin/jaira-role-brainstorm/SKILL.md:77) und Projekt-Skill (.claude/skills/jaira/SKILL.md:286). Beide als DoD-Punkte 5 und 6 angehängt.
