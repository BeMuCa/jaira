---
id: 01M4930V33NCPD29SW0ZV3S0MW
title: Ein Ticket landet nur zusammen mit Code in einem Commit
status: done
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
commits:
  - c0e38127be11012a6af8dcda3aa804a6d7f099ad
  - aff2304a20867cc8eaa42bf7d2e4d534b5153c82
  - 9dfd359398e27dce1752727002d25bb857afaba0
  - 01cd6d3b29e3e466f61a2b840e0abb73090bd51e
created-at: 2026-10-06T17:08:45Z
updated-at: 2026-10-06T21:53:37Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-22138
claimed-at: 2026-10-06T17:08:56Z
outcome-what: "Brainstorm-Rolle und Projekt-Skill sagen jetzt wie der Agent-Block: erstellte Tickets nicht committen, außer im File-Modus; NOTES-Zeile nennt jaira roles install --force."
outcome-why: "Review fand zwei Anweisungen, die der neuen Regel widersprachen; die Brainstorm-Rolle ist die wahrscheinlichste Quelle der Backlog-Stapel auf master."
outcome-resolves: "Keine mitgelieferte Anweisung schickt mehr ein frisch erstelltes Ticket in einen Commit (git grep 'earns a commit' leer)."
question: "Regel so richtig? Danach: die 32 Altfiles auf master räumen (vorher prüfen, welche Branches sie anfassen)."
review-summary: "Der generierte Agent-Block (core/board/announce.go -> CLAUDE.md, AGENTS.md), docs/AGENTS.md, die Brainstorm-Rolle und der Projekt-Skill sagen jetzt einheitlich: ein erstelltes Ticket wird nicht committet, es reist auf seinem Ref und kommt per jaira pull in den Branch; nur ohne nutzbares Remote committen. scripts/check-ticket-commits.sh prüft jeden Nicht-Merge-Commit eines PR und lehnt einen ab, der ein Ticketfile unter .jaira/tickets/ hinzufügt oder ändert, ohne eine Datei außerhalb .jaira/ zu ändern; Löschen ist frei. Der CI-Job ticket-commits führt es auf pull_request aus."
review-gaps: "none — Rest-Risiko, nicht Lücke: ein Board im File-Modus würde in diesem Repo von CI abgelehnt; dieses Repo läuft im Ref-Modus. Der CI-Job selbst ist erst auf GitHub bewiesen, wenn der PR läuft."
test-verdict: "pass — scripts/check-ticket-commits.sh master: exit 0; git grep 'earns a commit' (ohne .planning/.jaira): leer; go test ./core/role grün."
review-verdict: "der Diff erfüllt alle sechs DoD-Punkte; keine Defekte gefunden. Unsicher nur, ob der Job auf GitHub so läuft wie lokal — das zeigt der erste PR."
review-check: "1. git switch fix/ticket-rides-with-code  2. scripts/check-ticket-commits.sh master -> keine Ausgabe, Exit-Code 0 (echo $?)  3. git grep -n 'earns a commit' -- ':!.planning' ':!.jaira' -> keine Treffer  4. sed -n 100,112p CLAUDE.md -> Absatz enthält 'A ticket you create commits nothing either'  5. sed -n 74,82p core/role/builtin/jaira-role-brainstorm/SKILL.md -> 'Do not commit them.'  6. Nach dem PR: der Job 'ticket-commits' in GitHub Actions ist grün"
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
- [x] core/role/builtin/jaira-role-brainstorm/SKILL.md sagt nicht mehr, die angelegten Tickets zu committen
  proof: core/role/builtin/jaira-role-brainstorm/SKILL.md:77 (9dfd359)
- [x] .claude/skills/jaira/SKILL.md sagt dasselbe wie der Agent-Block: ein erstelltes Ticket wird nur im File-Modus committet
  proof: .claude/skills/jaira/SKILL.md:286 (9dfd359)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-06 17:12 · Alexander Sacharov** — Die CI-Prüfung läuft nur auf pull_request über base..head. Die alte Historie verletzt die Regel dutzendfach (z.B. ec2c861, 4cf0249) — ein Check über die ganze Historie wäre für immer rot. Löschen eines Ticketfiles ist erlaubt, damit logbook/archive und das Aufräumen der 32 Altfiles auf master durchgehen. 13 Tests in internal/cli schlagen auf diesem Rechner auch ohne diese Änderung fehl; mit leerem HOME sind alle grün — Ursache ist ~/.jaira, nicht der Code.
- **2026-10-06 18:57 · Alexander Sacharov** — Alex am 2026-10-06: Regel und CI-Prüfung angenommen wie gebaut (Ticketfile nur mit Code außerhalb .jaira/). Nächster Schritt danach: die 32 Altfiles auf master räumen.
- **2026-10-06 19:03 · Alexander Sacharov** — Review 2026-10-06: zwei Anweisungen widersprechen der Regel noch — Brainstorm-Rolle (core/role/builtin/jaira-role-brainstorm/SKILL.md:77) und Projekt-Skill (.claude/skills/jaira/SKILL.md:286). Beide als DoD-Punkte 5 und 6 angehängt.
- **2026-10-06 19:05 · Alexander Sacharov** — Von testing direkt nach review, ohne human: keine offene Frage. Nach Alex' Entscheidung vom 2026-10-06 (siehe BG5QJ6) nimmt der Mensch in signoff ab.
- **2026-10-06 21:53 · Alexander Sacharov** — Abnahme 2026-10-06 (https://claude.ai/artifact/JwDPU5tVLGPXLp5uNFNGb2): alle Schritte von Alex ok. Die Maschinenprüfung 'git grep earns a commit' war rot, weil sie die Ticketfiles mitsuchte, die den alten Satz als Ursache zitieren — Schritt korrigiert (':!.jaira'), danach grün. Kein Codefehler.
