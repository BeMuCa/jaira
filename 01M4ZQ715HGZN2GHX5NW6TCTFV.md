---
id: 01M4ZQ715HGZN2GHX5NW6TCTFV
title: "Ist der letzte Ticket eines Milestones fertig, schlaegt jaira den PR mit /jaira-role-pr vor und bietet an, den Milestone abzulegen"
status: pre-process
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Sobald alle Tickets eines Milestones fertig sind, sagt jaira (und jeder Teamlead/Dispatcher) dem Menschen: jetzt PR mit /jaira-role-pr, und danach den Milestone ablegen. Geoeffnet wird der PR weiter nur, wenn ein Mensch die Rolle aufruft"
context: |-
  Was heute falsch ist: Der PR-Skill (jaira-role-pr) wird oft vergessen. Keine Lane loest ihn aus, und kein Bericht schlaegt ihn vor.
  Beispiel aus einer anderen Session (requirementsgenie, Ticket 6ESQWR): Das Ticket kam in signoff an, und niemand schlug '/jaira-role-pr 6ESQWR' vor. Der MR !56 musste von Hand nachgezogen werden.
  Absicht von jaira, und die bleibt: Den PR oeffnet nur ein Mensch. jaira-role-pr SKILL.md:11 'Opening the pull request is the maintainer's call, and typing /jaira-role-pr is how they give it'; :16 von einem Agenten aufgerufen pusht er nur. Dispatcher :409 und Teamlead :167 oeffnen nie einen PR.
  Berk am 15.10.: 'wenn wir die Tickets, die unter einem Milestone stehen, alle fertiggestellt haben, ein PR vorschlagen und anbieten, das mit dem jaira pr skill zu machen.'
  Berk fragte auch, warum fertige Milestones (milestone-roadmap, ram-cpu) noch auf der Roadmap stehen. Ablegen geht schon mit 'jaira logbook <milestone>' (Ref traegt dann status: filed), aber niemand schlaegt es vor.
  Andere Ideen, die verworfen bzw. eingeordnet wurden:
  - nur eine Memory-Regel fuer eine Session: zu schmal
  - CLAUDE.md jaira:local: gilt nur fuer ein Repo
  - PostToolUse-Hook: zerbrechlich, nur Claude Code
  Gewaehlt: in jaira selbst. Die CLI-Ausgabe erreicht jeden Agenten, nicht nur Claude (Projekt-Constraint 'CLI must be usable by any bash-capable agent').
  Verwandt: 3SJMV8 (PR-Abschnitt im jaira-Skill), N71NVG (fertigen Milestone archivieren).
definition-of-done: "Macht ein 'jaira move <id> --to done' oder 'jaira logbook <id>' den letzten offenen Ticket eines Milestones fertig, druckt jaira eine Zeile: PR mit '/jaira-role-pr <ids>' vorschlagen und danach 'jaira logbook <milestone>'; --json traegt es als Feld; Test"
tags: []
blocked-by: []
related:
  - 01M4VGHHEKBXMV9120AD3SJMV8
  - 01M2K275PSXDASRN6Y1VN71NVG
commits: []
created-at: 2026-10-15T12:04:57Z
updated-at: 2026-10-15T12:39:03Z
updated-by: BeMuCa
---

# Ist der letzte Ticket eines Milestones fertig, schlaegt jaira den PR mit /jaira-role-pr vor und bietet an, den Milestone abzulegen

## Definition of Done

- [ ] Macht ein 'jaira move <id> --to done' oder 'jaira logbook <id>' den letzten offenen Ticket eines Milestones fertig, druckt jaira eine Zeile: PR mit '/jaira-role-pr <ids>' vorschlagen und danach 'jaira logbook <milestone>'; --json traegt es als Feld; Test
- [ ] jaira-dispatcher und jaira-teamlead nennen in ihrem Bericht den PR-Vorschlag, sobald ein Milestone komplett ist, und rufen /jaira-role-pr nie selbst auf
- [ ] jaira-role-pr nimmt mehrere Tickets bzw. einen Milestone und legt alle angenommenen Tickets vor dem Push ab; fuer einen bestehenden PR/MR aktualisiert er Titel und Beschreibung aus den Tickets
- [ ] Eine Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [ ] Schritt 1 (core/milestone): Feld 'pr: block|single' — erst Tests core/milestone/pr_test.go (TestPRLineIsReadAndRewrittenInPlace, TestAMilestoneWithoutPRIsSingle, TestUnknownPRValueReadsAsSingle, TestParsePRRefusesAnUnknownMode, TestAdoptTakesAPRChangeFromTheOldSide), dann neue Datei pr.go (PRBlock, PRSingle, Block(), SetPR, ParsePR) + Feld PR im Struct und 'case "pr"' in parse(); check.sh
- [ ] Schritt 2 (CLI): 'milestone create/set --pr block|single' (falscher Wert exit 2), milestoneJSON traegt immer 'pr', 'milestone ls' markiert 'pr: block' — erst Tests internal/cli/milestonepr_test.go; docs/COMMANDS.md Zeilen create/set/ls; erste Zeile in core/release/notes/6TCTFV.md; check.sh
- [ ] Schritt 3 (neues Paket core/suggest): AfterMove, AfterLogbook, ForMilestone, Complete, Lines, fileMilestoneCommand — erst Tests core/suggest/suggest_test.go (12 Tests laut design.md, u.a. Block komplett, Einzeln komplett, leer, abgelegt, zwei Milestones, Lanes nach Eigenschaft); check.sh
- [ ] Schritt 4 (jaira move): kleiner Hook in flow.go newMoveCmd — Zeilen nach nextStepLine, --json 'suggestions' nur wenn nicht leer — erst Tests internal/cli/prsuggest_test.go (TestMoveToSignoffSuggestsThePullRequest, TestMoveJSONCarriesSuggestions, TestMoveCompletingABlockMilestoneSuggestsOnePullRequestAndFiling, TestMoveIntoAnAgenticLaneSuggestsNothing); docs move-Zeile; Release-Zeile; check.sh
- [ ] Schritt 5 (jaira logbook <id>): kleiner Hook in logbookOut — erst Tests (TestLogbookOfTheLastTicketOffersToFileTheMilestone, TestLogbookJSONCarriesSuggestions, TestSuggestedFilingIsAccepted, TestLogbookOfABlockTicketSuggestsTheBlockPullRequest); docs logbook-Zeile; Release-Zeile; check.sh
- [ ] Schritt 6 (jaira milestone ls): komplette Milestones bekommen dieselben Zeilen, --json 'suggestions' pro Milestone — erst Test TestLsCarriesTheSuggestionsOfACompleteMilestone; docs ls-Zeile; Release-Zeile; check.sh
- [ ] Schritt 7 (Rollen): jaira-dispatcher, jaira-worker, jaira-teamlead geben jaira-Vorschlaege woertlich weiter und rufen /jaira-role-pr nie auf; Teamlead liest 'jaira milestone ls --json' — erst Pin-Test core/role/prsuggest_test.go TestRolesPassThePRSuggestionOnAndNeverRunIt; Release-Zeile; check.sh
- [ ] Schritt 8 (Rolle jaira-role-pr): nimmt mehrere Tickets oder einen Milestone; Ticket ohne Commit auf dem Branch = Befund, Stopp vor dem Push; offener PR/MR per 'gh pr edit'/'glab mr update' aktualisiert, nur wenn ein Mensch aufrief — erst Pin-Test TestPRTakesSeveralTicketsOrAMilestone; die Strings, die TestPRFilesTheTicketsAwayBeforeThePush (core/role/role_test.go:504) prueft, bleiben stehen; Release-Zeile; check.sh
- [ ] Schritt 9: Handtest mit einem Build aus dem Branch in einem Scratch-Repo ausserhalb von /home/berk/git (HOME und JAIRA_HOME dort): kein Milestone, single, block; menschlich und --json; das vorgeschlagene 'jaira logbook <m>' laeuft durch — Ergebnis als jaira note

## Progress
- **2026-10-15 12:20 · BeMuCa** — Berk am 15.10. (Ergaenzung): '6TCTFV soll auch umgesetzt werden! Weil dann kann ein Milestone-Block als ein PR reinkommen. Alternativ gibt es aber auch noch die Option, dass man Tickets immer noch einzeln reinmergen will und Milestones nur als eine Sammlung nutzt; fuer den Fall kann man Milestones auch wieder loeschen.'
Daraus zwei Arbeitsweisen, beide muessen gehen:
- Block: der Milestone kommt als EIN PR. Ist sein letztes Ticket fertig, schlaegt jaira '/jaira-role-pr' fuer alle seine Tickets zusammen vor und danach das Ablegen des Milestones nach dem Merge.
- Einzeln: jedes Ticket bekommt seinen eigenen PR, der Milestone ist nur eine Sammlung. Jedes Ticket, das done oder signoff erreicht, bekommt einen Vorschlag '/jaira-role-pr <id>'. Ist der Milestone komplett, bietet jaira an, ihn abzulegen oder zu loeschen, und schlaegt KEINEN Block-PR vor.
Wie ein Milestone sagt, welche Arbeitsweise gilt (z.B. Feld 'pr: block|single' per milestone create/set, mit Default), entscheidet pre-process. Ist der Default nicht offensichtlich: Empfehlung + Frage an Berk ueber die human-Lane.
Den PR oeffnet nie jaira oder ein Agent von sich aus: er wird nur vorgeschlagen; ein Mensch oeffnet ihn mit /jaira-role-pr.
Parallel: N71NVG baut das Ablegen/Loeschen eines fertigen Milestones samt Tickets. Hier nur Vorschlaege mit den Kommandos, die auf master existieren ('jaira logbook <milestone>').
