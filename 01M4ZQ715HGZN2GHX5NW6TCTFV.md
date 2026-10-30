---
id: 01M4ZQ715HGZN2GHX5NW6TCTFV
title: "Ist der letzte Ticket eines Milestones fertig, schlaegt jaira den PR mit /jaira-role-pr vor und bietet an, den Milestone abzulegen"
status: critique
ready: true
creator: BeMuCa
assignee: berkc
goal: "Sobald alle Tickets eines Milestones fertig sind, sagt jaira (und jeder Teamlead/Dispatcher) dem Menschen: jetzt PR mit /jaira-role-pr, und danach den Milestone ablegen. Geoeffnet wird der PR weiter nur, wenn ein Mensch die Rolle aufruft"
context: |-
  Was heute falsch ist: Der PR-Skill (jaira-role-pr) wird oft vergessen. Keine Lane loest ihn aus, und kein Bericht schlaegt ihn vor.
  Beispiel aus einer anderen Session (ein anderes Projekt, Ticket 6ESQWR): Das Ticket kam in signoff an, und niemand schlug '/jaira-role-pr 6ESQWR' vor. Der MR !56 musste von Hand nachgezogen werden.
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
updated-at: 2026-10-30T23:59:15Z
updated-by: berkc
question: ""
round: 2
outcome-what: "jaira move, jaira logbook <id> and jaira milestone ls print the pull-request line and the milestone filing line when a ticket or a complete milestone calls for one, and --json carries them as suggestions; the rules live in core/suggest; the roles pass the lines on word for word; jaira-role-pr takes several tickets or a milestone and merges their branches into one before the push"
outcome-why: "a milestone is one pull request over all its tickets (berkc 22.10), and the person is otherwise never told /jaira-role-pr is due; a ticket outside any milestone is a group of one"
outcome-resolves: "a finished ticket or a finished milestone now names its next step in the output; an agent never opens a pull request; the one-branch question is settled by merging the ticket branches at PR time, not by changing worktrees"
review-summary: "none: pass 1 read the diff of core/suggest, internal/cli suggest/flow/logbook/milestones and the three role paragraphs; no finding that names a file and a clearer shape; the observations are in the note of this pass"
---

# Ist der letzte Ticket eines Milestones fertig, schlaegt jaira den PR mit /jaira-role-pr vor und bietet an, den Milestone abzulegen

## Definition of Done

- [x] Milestone komplett: Hat 'jaira move <id> --to <lane>' oder 'jaira logbook <id>' jedes Ticket eines Milestones in die Endlane oder vom Board gebracht (dieselbe Regel wie 'jaira logbook <milestone>'), druckt jaira zwei Zeilen: 'Milestone <m> is complete. Next: one pull request for all its tickets. Only a person opens it: /jaira-role-pr <alle Handles>' und 'Once that pull request is merged, take milestone <m> off the board: jaira logbook <m>'; --json traegt beide unter 'suggestions' (kind pull_request mit by person, kind file_milestone mit after merge); Tests TestMoveCompletingAMilestoneSuggestsOnePullRequestAndFiling, TestLogbookOfTheLastTicketSuggestsThePullRequest
  proof: internal/cli/prsuggest_test.go TestMoveCompletingAMilestoneSuggestsOnePullRequestAndFiling, TestLogbookOfTheLastTicketSuggestsThePullRequest; core/suggest/suggest.go ForTicket
- [x] jaira-dispatcher, jaira-worker und jaira-teamlead geben jede Vorschlagszeile von jaira ('suggestions' aus move/logbook, beim Teamlead auch aus 'jaira milestone ls --json') woertlich an den Menschen weiter und rufen /jaira-role-pr nie selbst auf; Test TestRolesPassThePRSuggestionOnAndNeverRunIt
  proof: core/role/prsuggest_test.go TestRolesPassThePRSuggestionOnAndNeverRunIt; core/role/builtin/jaira-{dispatcher,worker,teamlead}/SKILL.md
- [x] jaira-role-pr nimmt mehrere Tickets bzw. einen Milestone: er legt alle angenommenen Tickets vor dem Push ab und baut bei einem Milestone EINEN Branch aus origin/master, in den er jeden Ticket-Branch mit --no-ff hineinmerged; ein genanntes Ticket ohne Commit auf seinem Branch ist ein Befund (Stopp vor dem Push), ein Merge-Konflikt ebenso. Fuer einen bestehenden PR/MR aktualisiert er Titel und Beschreibung aus den Tickets (gh pr edit / glab mr update), aber nur wenn ein Mensch ihn aufrief; Test TestPRTakesSeveralTicketsOrAMilestone
  proof: core/role/prsuggest_test.go TestPRTakesSeveralTicketsOrAMilestone; core/role/builtin/jaira-role-pr/SKILL.md ## One branch for several tickets
- [x] Je eine Zeile pro sichtbarer Aenderung in core/release/notes/6TCTFV.md: Vorschlaege bei move, bei logbook <id>, in milestone ls, Weitergabe durch die Rollen, jaira-role-pr mit mehreren Tickets
  proof: core/release/notes/6TCTFV.md (five lines); docs/COMMANDS.md move, logbook, milestone ls rows
- [-] Ein Milestone sagt mit 'pr: block' oder 'pr: single', wie seine Tickets reinkommen; gesetzt mit 'jaira milestone create|set <name> --pr <mode>', ein anderer Wert wird abgelehnt (exit 2); ohne Zeile gilt single (auch fuer jede alte Datei), ein unbekannter Wert in der Datei wird als single gelesen; 'milestone ls --json' traegt 'pr'; Tests TestCreateAndSetWriteThePRMode, TestPRRefusesAnUnknownMode, TestAMilestoneWithoutPRIsSingle
- [x] Ticket fertig: Erreicht ein Ticket per 'jaira move' oder 'jaira logbook' die Abnahme-Lane (signoff) oder die Endlane, druckt jaira den Stand. Ist es in keinem Milestone, druckt es 'Next: the pull request. Only a person opens it: /jaira-role-pr <id>' und --json 'suggestions' (kind pull_request, by person). Liegt es in einem offenen Milestone, druckt es nur 'Milestone <m>: k of n done' und KEINE PR-Zeile; Tests TestMoveToSignoffSuggestsThePullRequest, TestMoveJSONCarriesSuggestions, TestTicketInAnOpenMilestoneSuggestsNoPullRequest
  proof: internal/cli/prsuggest_test.go TestMoveToSignoffSuggestsThePullRequest, TestMoveJSONCarriesSuggestions, TestTicketInAnOpenMilestoneSuggestsNoPullRequest
- [-] Einzeln, Milestone komplett: Ist ein Milestone ohne pr: block nach move oder logbook <id> komplett, bietet jaira nur 'jaira logbook <m>' an und schlaegt KEINEN Block-PR vor; 'jaira logbook <m>' nimmt den vorgeschlagenen Befehl an; Tests TestSingleMilestoneCompleteOffersFilingAndNoBlockPullRequest, TestLogbookOfTheLastTicketOffersToFileTheMilestone, TestSuggestedFilingIsAccepted
- [x] Abnahme- und Endlane werden an ihren Eigenschaften erkannt, nicht am Namen; ein leerer oder abgelegter Milestone schlaegt nie etwas vor; ein Ticket in mehreren Milestones bekommt die Zeilen jedes kompletten Milestones; Tests TestLanesAreFoundByWhatTheyDoNotByName, TestEmptyMilestoneIsNeverComplete, TestFiledMilestoneSuggestsNothing, TestTicketInTwoMilestones
  proof: core/suggest/suggest_test.go TestLanesAreFoundByWhatTheyDoNotByName, TestEmptyMilestoneIsNeverComplete, TestFiledMilestoneSuggestsNothing, TestTicketInTwoMilestonesGetsTheLinesOfEachComplete
- [x] 'jaira milestone ls' zeigt fuer jeden kompletten Milestone dieselben Zeilen (--json: 'suggestions' am Milestone), weil das Annehmen in der TUI nichts druckt; Test TestLsCarriesTheSuggestionsOfACompleteMilestone
  proof: internal/cli/prsuggest_test.go TestLsCarriesTheSuggestionsOfACompleteMilestone

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [-] Schritt 1 (core/milestone): Feld 'pr: block|single' — erst Tests core/milestone/pr_test.go (TestPRLineIsReadAndRewrittenInPlace, TestAMilestoneWithoutPRIsSingle, TestUnknownPRValueReadsAsSingle, TestParsePRRefusesAnUnknownMode, TestAdoptTakesAPRChangeFromTheOldSide), dann neue Datei pr.go (PRBlock, PRSingle, Block(), SetPR, ParsePR) + Feld PR im Struct und 'case "pr"' in parse(); check.sh
- [-] Schritt 2 (CLI): 'milestone create/set --pr block|single' (falscher Wert exit 2), milestoneJSON traegt immer 'pr', 'milestone ls' markiert 'pr: block' — erst Tests internal/cli/milestonepr_test.go; docs/COMMANDS.md Zeilen create/set/ls; erste Zeile in core/release/notes/6TCTFV.md; check.sh
- [x] Schritt 1 (neues Paket core/suggest): AfterMove, AfterLogbook, ForMilestone, Complete, Lines, fileMilestoneCommand — erst Tests core/suggest/suggest_test.go (TestTicketInAnOpenMilestoneSuggestsNoPullRequest, TestMoveToSignoffSuggestsThePullRequest, TestMilestoneCompleteSuggestsOnePullRequestThenFiling, TestTicketInTwoMilestones, TestLanesAreFoundByWhatTheyDoNotByName, TestEmptyMilestoneIsNeverComplete, TestFiledMilestoneSuggestsNothing, TestSuggestedFilingIsAccepted, TestLinesSayOnlyAPersonOpensThePullRequest); dann Code; go test ./core/suggest/
- [x] Schritt 2 (jaira move): kleiner Hook in flow.go newMoveCmd — Zeilen nach nextStepLine, --json 'suggestions' nur wenn nicht leer — erst Tests internal/cli/prsuggest_test.go (TestMoveToSignoffSuggestsThePullRequest, TestMoveJSONCarriesSuggestions, TestMoveCompletingAMilestoneSuggestsOnePullRequestAndFiling, TestMoveIntoAnAgenticLaneSuggestsNothing); docs move-Zeile; Release-Zeilen; go test ./internal/cli/ fuer die Pakete
- [x] Schritt 3 (jaira logbook <id>): kleiner Hook in logbookOut — erst Tests (TestLogbookOfTheLastTicketSuggestsThePullRequest, TestLogbookJSONCarriesSuggestions); docs logbook-Zeile; Release-Zeile; Tests
- [x] Schritt 4 (jaira milestone ls): komplette Milestones bekommen dieselben Zeilen, --json 'suggestions' pro Milestone — erst Test TestLsCarriesTheSuggestionsOfACompleteMilestone; docs ls-Zeile; Release-Zeile; Tests
- [x] Schritt 5 (Rollen): jaira-dispatcher, jaira-worker, jaira-teamlead geben jaira-Vorschlaege woertlich weiter und rufen /jaira-role-pr nie auf; Teamlead liest 'jaira milestone ls --json' — erst Pin-Test core/role/prsuggest_test.go TestRolesPassThePRSuggestionOnAndNeverRunIt; Release-Zeile nur wenn sichtbar; Tests
- [x] Schritt 6 (Rolle jaira-role-pr): nimmt mehrere Tickets oder einen Milestone; Milestone = ein Branch aus origin/master, jeder Ticket-Branch per merge --no-ff hinein (siehe Notiz zum Ein-Branch); Ticket ohne Commit = Befund, Merge-Konflikt = Stopp vor dem Push; offener PR/MR per 'gh pr edit'/'glab mr update', nur wenn ein Mensch aufrief — erst Pin-Test TestPRTakesSeveralTicketsOrAMilestone; die Strings, die TestPRFilesTheTicketsAwayBeforeThePush (core/role/role_test.go:504) prueft, bleiben stehen; Release-Zeile; Tests
- [x] Schritt 7 (Handtest): Build aus dem Branch in einem Scratch-Repo ausserhalb von /home/berk/git (HOME und JAIRA_HOME dort): ein Milestone mit zwei Tickets, eins fertig, dann beide fertig; kein Milestone; menschlich und --json; das vorgeschlagene 'jaira logbook <m>' laeuft durch; Merge der Ticket-Branches einmal ohne, einmal mit Konflikt; Ergebnis als jaira note

## Progress
- **2026-10-15 12:20 · BeMuCa** — Berk am 15.10. (Ergaenzung): '6TCTFV soll auch umgesetzt werden! Weil dann kann ein Milestone-Block als ein PR reinkommen. Alternativ gibt es aber auch noch die Option, dass man Tickets immer noch einzeln reinmergen will und Milestones nur als eine Sammlung nutzt; fuer den Fall kann man Milestones auch wieder loeschen.'
Daraus zwei Arbeitsweisen, beide muessen gehen:
- Block: der Milestone kommt als EIN PR. Ist sein letztes Ticket fertig, schlaegt jaira '/jaira-role-pr' fuer alle seine Tickets zusammen vor und danach das Ablegen des Milestones nach dem Merge.
- Einzeln: jedes Ticket bekommt seinen eigenen PR, der Milestone ist nur eine Sammlung. Jedes Ticket, das done oder signoff erreicht, bekommt einen Vorschlag '/jaira-role-pr <id>'. Ist der Milestone komplett, bietet jaira an, ihn abzulegen oder zu loeschen, und schlaegt KEINEN Block-PR vor.
Wie ein Milestone sagt, welche Arbeitsweise gilt (z.B. Feld 'pr: block|single' per milestone create/set, mit Default), entscheidet pre-process. Ist der Default nicht offensichtlich: Empfehlung + Frage an Berk ueber die human-Lane.
Den PR oeffnet nie jaira oder ein Agent von sich aus: er wird nur vorgeschlagen; ein Mensch oeffnet ihn mit /jaira-role-pr.
Parallel: N71NVG baut das Ablegen/Loeschen eines fertigen Milestones samt Tickets. Hier nur Vorschlaege mit den Kommandos, die auf master existieren ('jaira logbook <milestone>').
- **2026-10-15 12:40 · BeMuCa** — pre-process, Grundlagen des Plans:
- Feld: 'pr: block|single' im Milestone-Frontmatter, gebaut wie 'goal' (core/milestone/milestone.go:405, :533). Gesetzt per 'milestone create|set --pr'.
- Default: keine Zeile = single. Grund 1: alte Milestones (ram-cpu, milestone-roadmap) kamen Ticket fuer Ticket rein; als block gelesen wuerden sie schon gemergte Tickets nochmal als PR vorschlagen. Grund 2: hier hat jedes Ticket seinen eigenen Branch (Dispatcher 'One worktree per ticket'); ein Block-PR braucht alle auf einem Branch, das muss man bewusst wollen. Grund 3: das TUI-Formular (internal/tui/layer.go:457) schreibt keine Zeile, CLI und TUI stimmen so ohne TUI-Aenderung ueberein.
- OFFEN fuer Berk: soll 'milestone create' ohne --pr trotzdem 'pr: block' schreiben? Sein erster Satz nennt Block. Aenderung waere eine Zeile an der create-Stelle plus layer.go:457.
- Fertig = genau die Regel von 'jaira logbook <milestone>' (internal/cli/logbook.go:393-404): Ticket in der Endlane oder nicht mehr auf dem Board. Dazu: ein leerer Milestone ist nie fertig.
- Ausloeser ist ein Zustand, kein Uebergang: 'logbook <id>' nimmt nur Tickets aus der Endlane, macht also nie etwas neu fertig. Nach move, logbook <id> und milestone ls wird gefragt: ist der Milestone jetzt komplett? Die Zeile kann sich wiederholen (z.B. bei jedem logbook in jaira-role-pr) - harmlos, sie bleibt wahr bis zum Ablegen.
- Einzel-Vorschlag bei move in die Abnahme-Lane (requires-human-exit, keine Frage-Lane, nicht agentisch - dieselbe Lesart wie core/lane/lane.go:652) oder in die Endlane. Lanes nach Eigenschaft, nie nach Namen.
- Ticket ohne Milestone bekommt den Einzel-Vorschlag. Grund: das Beispiel im Kontext (6ESQWR in signoff). single heisst 'Milestone ist nur Sammlung', Mitgliedschaft aendert also nichts. Kostet keinen Code: Bedingung ist nur 'in keinem Block-Milestone'.
- Logik in neuem Paket core/suggest: core/move/move.go:1-7 sagt, Regeln nach einem Move gehoeren in core, damit die TUI-Schreibstellen sie spaeter teilen.
- logbookMilestone wird NICHT umgebaut, dort baut N71NVG. TestSuggestedFilingIsAccepted pinnt, dass jaira nur ein Ablegen vorschlaegt, das 'jaira logbook <m>' annimmt.
- N71NVG steckt hier ein: Befehlsname nur in core/suggest fileMilestoneCommand (eine Zeile). Loeschen bei single: ein Eintrag mehr in ForMilestone plus ein Fall in Line.
- --json: 'suggestions' nur wenn nicht leer; jede bestehende Ausgabe bleibt gleich (wie filed_on_entry, internal/cli/flow.go:239).
- milestone ls bekommt die Zeilen, weil Annehmen in der TUI (internal/tui/signoff.go:187) der uebliche Weg nach done ist und nichts druckt. Der Teamlead liest 'milestone ls --json'.
- jaira-worker mitgenommen: Schwesterrolle des Dispatchers, verschiebt Tickets selbst; ohne ihn ginge die Zeile bei Worker-Tickets verloren.
- Unveraendert: core/board/announce.go (sagt schon 'jaira logbook <milestone>'), core/milestone/adopt.go (merged jedes Frontmatter-Feld generisch, adopt.go:84), nextStepLine, move.Move.
- Ausgelassen, Folge-Tickets: TUI-Schreibstellen (signoff.go:187, model.go:1949, :2004); 'jaira logbook --all'; Tickets eines Block-Milestones auf einen Branch bringen (jaira-role-pr meldet ein Ticket ohne Commit auf dem Branch und stoppt).
- **2026-10-22 00:15 · berkc** — Antwort von berkc am 22.10.2026 auf die Frage A/B: Alarm in beiden Faellen - wenn ein Ticket fertig ist UND wenn ein Milestone fertig ist. PRs pro Ticket entfallen; der PR entsteht pro Milestone ('we can skip the single ticket prs and wait for milestone'). Damit gibt es kein single/block-Attribut mehr: jeder Milestone bekommt beim Abschluss einen PR ueber alle seine Tickets. Ein Ticket ohne Milestone: offen, ob es einen eigenen PR bekommt - beim Planen pruefen und notfalls nachfragen. Folge: alle Tickets eines Milestones muessen auf einem Branch landen (siehe Plan-Notiz: heute hat jedes Ticket seinen eigenen Branch). Verwandt: 80MGXT (Worktree weg bei done, Branch bleibt bis zum Milestone-Merge).
- **2026-10-30 23:49 · berkc** — pre-process (worker, 2026-10-31): open decisions counted from the ticket and its notes. (a) single/block attribute on a milestone: CLOSED by berkc 22.10 - no attribute, one PR per milestone, alert on ticket done and on milestone done. (b) ticket without a milestone: CLOSED here - it gets the ticket-done alert with its own /jaira-role-pr <id> line, because no milestone PR will ever carry it. A group of one, not a PR per milestone ticket. Reversible if berkc says otherwise. (c) one PR over tickets that each live on their own branch: CLOSED in the plan, see the next note. Nothing is left that only the person can take, so no question and no conversational mode.
- **2026-10-30 23:49 · berkc** — pre-process: the one-branch consequence. (1) One shared worktree and branch per milestone - rejected: git allows one checkout per branch, so tickets of one milestone would serialize, and the dispatcher worktree logic would change = not smaller than paca. (2) Cherry-pick ticket commits onto a milestone branch - rejected: rewrites commit ids for no gain. (3) CHOSEN: when the milestone PR is made, jaira-role-pr creates one branch from origin/master and merges each member ticket branch into it with --no-ff; ticket commit ids stay as they are, so the ticket commit lists still match. A conflict stops before the push and is reported as a finding. Each ticket keeps its own worktree and branch as today; the branches stay until the milestone is filed (80MGXT). Risk to test in step 8: .jaira/milestones/ files that differ between branches; KWQCYZ (stale milestone files overwrite refs) must not be reproduced.
- **2026-10-30 23:49 · berkc** — pre-process: ticket without a milestone - decided. Ticket-done alert (move or logbook into signoff or done) prints one PR line for that ticket. A ticket in a milestone gets only the progress line 'Milestone m: k of n done', never a PR line. A complete milestone (every member in the terminal lane or off the board, the same rule as logbook m) gives one PR line for all its handles plus the filing line. An empty or filed milestone never suggests anything. A ticket in two milestones gets the lines of each complete one.
- **2026-10-30 23:50 · berkc** — pre-process: question field cleared - berkc answered it on 22.10 (see Progress: no single/block attribute). Plan and DoD rewritten to that decision; old items 5, 7 (DoD) and 1, 2 (plan) superseded, not deleted.
- **2026-10-30 23:58 · berkc** — handtest (plan step 9): built the branch binary to the scratchpad (handtest-6TCTFV/bin/jaira), scratch repo with HOME and JAIRA_HOME inside the scratchpad. round-one with A and B: A to signoff -> 'Milestone round-one: 0 of 2 done', no PR line; B to done -> '1 of 2 done'; A to done -> PR line for both handles + logbook line, same in milestone ls and in --json. The suggested 'jaira logbook round-one' filed it (exit 0), ls empty after. Ungrouped C to signoff -> '/jaira-role-pr <C>' line. Plain git check of the merge step: two feature branches merged --no-ff into milestone/round-one = clean; a conflicting pair = rc 1, merge --abort. NOT tested: the jaira-role-pr skill itself (an agent run that pushes to a real remote), and the milestone file differing between ticket branches (the KWQCYZ risk).
- **2026-10-30 23:59 · berkc** — critique (pass 1, worker): read the Go diff - core/suggest/suggest.go, internal/cli/suggest.go, the hooks in flow.go, logbook.go, milestones.go - and the three role SKILL paragraphs. One simplification made already: ForMilestone repeated the Filed/empty check that complete() does. Left alone on purpose: (a) groupOf reads every member of each milestone the ticket is in on every move; board sizes are tens, so no cache; (b) the same paragraph is in three role skills because the pin test reads each role on its own and the role format has no include; (c) --json suggestions also carry kind progress, which the DoD does not name; the docs say consumers branch on kind; (d) milestone ls swallows a lane-load error on purpose so the listing never fails because of suggestions.
