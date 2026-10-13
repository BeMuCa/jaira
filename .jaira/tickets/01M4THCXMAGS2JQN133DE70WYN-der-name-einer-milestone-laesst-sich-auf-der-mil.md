---
id: 01M4THCXMAGS2JQN133DE70WYN
title: Der Name einer Milestone laesst sich auf der Milestone Roadmap und per CLI aendern
status: testing
ready: true
creator: BeMuCa
goal: "Wer eine Milestone umbenennen will, tut das im Bearbeiten-Dialog der Milestone Roadmap (e) oder mit jaira milestone rename <alt> <neu>; danach heisst sie ueberall neu - Datei, Ref und die after:-Verweise ihrer Nachfolger - und kein anderer Klon behaelt eine verwaiste Kopie unter dem alten Namen."
context: |-
  Heute kann jaira Milestones nirgends umbenennen: jaira milestone kennt nur create, add, rm, ls und set; der e-Dialog der Milestone Roadmap aendert nur Ziel und Vorgaenger (internal/tui/layer.go, startLayerEdit).
  Berk am 2026-10-13: 'Editing a Milestone block in the Milestone Roadmap the name of the Milestone should also be editable.' In D0SAHM war genau das als eigenes Ticket entschieden (Entscheidung C2) - das ist es.
  Schwer ist nicht die Datei, sondern die Verteilung: der Name ist Dateiname (.jaira/milestones/<name>.md) UND Ref (refs/jaira/milestones/<name>). core/refsync IncomingMilestones loescht nie, ein anderer Klon behaelt also die alte Datei, wenn nur ein neuer Ref entsteht.
  Nachfolger nennen den Namen in after: [..] (D0SAHM); ein Board-Filter milestone:<alt> passt danach auf nichts.
  Verwandt: N71NVG (Milestone archivieren/loeschen: Datei und Ref weg) - dieselbe Frage, wie ein Ref auf anderen Klonen verschwindet. Mit N71NVG abstimmen, dessen Entscheidungen nicht vorwegnehmen.
definition-of-done: "jaira milestone rename <alt> <neu> benennt Datei und Ref um, schreibt die after:-Listen aller Nachfolger um und lehnt einen vergebenen oder abgelegten Namen ab; Tests in internal/cli"
tags:
  - tui
  - cli
blocked-by: []
related: []
follows: 01M1KFMETY9MKH4V38TVD0SAHM
commits: []
created-at: 2026-10-13T11:47:07Z
updated-at: 2026-10-13T13:39:30Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-1970307
claimed-at: 2026-10-13T13:33:56Z
executed-by: opus
outcome-what: "Review 1 answered: refsync.MilestoneOnRef + CLI refuseIfTaken + TUI nameTaken refuse a name a teammate's milestone has on this tree's ref but not yet on disk (exit 3, 'jaira fetch brings it here') for rename and create; rename help/COMMANDS exit code for a filed old name; TUI reloads when goal/after is refused after a rename; release note line for the new refusal."
outcome-why: "Every write's flush fetches refs without writing files, so rename and create, asking only the disk and filings, silently overwrote a teammate's milestone on the remote (reproduced by review) - DoD 1 asks for a taken name to be refused."
outcome-resolves: "DoD 1: TestANameOnlyOnARefIsTakenForRenameAndCreate (two clones) - rename b zeta and create zeta exit 3 while zeta is only on ada's ref, grace's zeta keeps its goal; fails without the check. DoD 2-4 unchanged; vet + go test ./... -race green (28 ok)."
review-summary: "none"
review-gaps: "Optimize after review 1 (commit 3c08f8a only): folded the doubled return in submitLayerForm's refused-Edit branch into one if (internal/tui/layer.go). Left: refsync MilestoneOnRef repeats the ref read of MilestoneFiled - folding it in as a new MilestoneFiledAt would reach refuseIfFiled's other callers (logbook, editRefusal) and reshape what critique settled; create/rename read the target's ref twice (once per command, not hot); CLI refuseIfTaken vs TUI nameTaken follows the existing CLI/TUI split. Earlier optimize round (f840f0e): one const for rename's refusal clause, one slices.Index in renameMilestone; Rename's successor loop vs ReattachSuccessors (layout.go) left as a follow-up candidate; pre-existing: SetAfter lowercases hand-written predecessor names, refuseIfFiled/editRefusal comments do not list rename."
test-verdict: "pass: vet clean, go test ./... -race -count=1 green (28 ok, exit 0); DoD 1-3 proof tests re-run by name and pass; binary (f840f0e) played in tmux on a board copy without remote - e-form name first, taken name refused, rename from layer and detail, successors and board filter follow; two clones + third clone on a local bare repo: after fetch the second clone lists only the new name, the third never gets the old file, renaming back reaches both"
review-verdict: "not yet: DoD 1 is met for names on disk, in the logbook or filed on a ref, but a teammate's milestone that is only on the ref is silently overwritten (reproduced) - back to in-progress for the shared name check; the minors go along where cheap, the rest is listed for Berk."
---

# Der Name einer Milestone laesst sich auf der Milestone Roadmap und per CLI aendern

## Definition of Done

- [x] jaira milestone rename <alt> <neu> benennt Datei und Ref um, schreibt die after:-Listen aller Nachfolger um und lehnt einen vergebenen oder abgelegten Namen ab; Tests in internal/cli
  proof: core/milestone/write.go Rename; internal/cli/milestones.go newMilestoneRenameCmd; TestRenameMovesTheMilestoneAndItsSuccessorsFollow, TestRenameRefusalsAndTheMarkOfARename, TestRenameMovesTheFileRewritesSuccessorsAndMarksTheOldName, TestRenameRefusesATakenNameButNotTheMarkOfARename
- [x] Im e-Dialog der Milestone Roadmap ist der Name das erste Feld; ein geaenderter Name benennt ueber denselben Weg um; TUI-Test
  proof: internal/tui/layer.go startLayerEdit (name first), submitLayerForm -> renameMilestone -> milestone.Rename; TestEditRenamesAMilestone, TestEditGoalAndPredecessors
- [x] Ein zweiter Klon, der nach dem Umbenennen fetcht, hat danach nur die Milestone unter dem neuen Namen und nicht zusaetzlich die alte; Test mit zwei Repos wie in den refsync-Tests
  proof: internal/cli/milestoneref_test.go TestARenamedMilestoneArrivesUnderItsNewNameOnly (two clones, bare remote; fails without the mark on the old ref)
- [x] go vet und go test ./... -race gruen; im echten Terminal auf einer Board-Kopie ohne Remote durchgespielt; Release-Notiz unter core/release/notes/
  proof: go vet ./... clean + go test ./... -race -count=1: 28 ok, exit 0 (at f840f0e); tmux play of the scratch build on a copy of fix2/board (no remote): e shows name first, taken name refused in the form, rename from layer and from detail, board filter followed; core/release/notes/E70WYN.md

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/milestone: renamed-to lesen (RenamedTo), HoldsName(), FiledError.RenamedTo, ExistsError, Rename() unter dem Lock (neue Datei, Nachfolger umschreiben, alte Datei weg, Marke auf den alten Ref); Speicherschleife mit Edit teilen; Tests in write_test.go
- [x] core/refsync MilestoneFiled: eine Umbenennungs-Marke haelt den Namen nicht (HoldsName)
- [x] CLI: jaira milestone rename <alt> <neu> mit Exit-Codes; editRefusal und logbook sagen 'umbenannt in b' statt restore; create-Tuer fragt HoldsName; Hilfe 'rename the file' ersetzen; Tests in milestones_test.go
- [x] CLI fetch: Umbenennungen eigen listen (Text und JSON milestones_renamed)
- [x] Zwei-Klon-Test in internal/cli/milestoneref_test.go: zweiter Klon hat nach fetch nur b, Nachfolger nennt b, Zuruecktaufen geht
- [x] TUI: e-Dialog mit name als erstem Feld (drei Felder), Speichern = pruefen, Rename, Edit; Board-Filter umschreiben; createMilestone fragt HoldsName; Hilfe; TUI-Tests
- [x] Doku: docs/COMMANDS.md, README Keys, Spec-Zeile, core/release/notes/E70WYN.md
- [x] go vet und go test ./... -race gruen, Commit mit Ticket
- [x] Binary bauen, in tmux auf Board-Kopie ohne Remote durchspielen; Zwei-Klon-Probe mit lokalem bare-Repo

## Progress
- **2026-10-13 12:01 · BeMuCa** — pre-process, Entscheidungen zum Umbenennen (a -> b):
1. Der alte Name verschwindet von anderen Klonen ueber eine MARKE auf seinem Ref, nicht ueber Loeschen. Der Ref a bleibt stehen und traegt den alten Inhalt plus 'status: filed' und eine neue Zeile 'renamed-to: b'.
   Warum: refsync IncomingMilestones schreibt einen filed-Ref ueber eine Datei, die ein Klon schon hat, schreibt ihn nie neu hin und loescht nie. Ein geloeschter Ref erreicht keinen Klon (so begruendet in internal/cli/logbook.go logbookMilestone). Auch ein alter jaira-Build nimmt a dann vom Board, weil er 'status: filed' kennt.
2. Auf dem zweiten Klon bleibt a.md auf der Platte liegen, markiert - nicht auf dem Board, nicht in ls, nicht im M-Picker. Genau das hinterlaesst heute ein logbook eines anderen Klons (TestAClonePlanningTheMilestoneLearnsItWasFiled: jaira loescht keine Datei, die es nur gelesen hat). DoD 3 heisst hier: Board, ls und Picker kennen nur b.
3. Der alte Name ist danach FREI, anders als bei logbook. Ein abgelegter Name bleibt belegt, weil restore ihn zurueckbraucht; nach einem Umbenennen gibt es nichts zurueckzuholen. Damit geht Zuruecktaufen (b -> a, z.B. nach einem Tippfehler) und ein neuer Milestone a. Regel an einer Stelle: Milestone.HoldsName() = filed und kein renamed-to; refsync.MilestoneFiled und die create-Tueren fragen das.
4. Im umbenennenden Klon wird a.md geloescht (kein Ablegen, nichts ins logbook).
5. Kein Ref-Loeschen gebaut - N71NVG (delete) wird nichts vorweggenommen; Ueberschneidung auf N71NVG notiert.
6. e-Dialog: drei Felder name, goal, after. Beim Speichern erst goal/after pruefen, dann umbenennen (milestone.Rename, derselbe Weg wie die CLI), dann goal/after schreiben. Ein abgelehnter Name schreibt also nichts. Der Board-Filter milestone:a wird auf b umgeschrieben.
7. jaira fetch listet Umbenennungen eigen ('Milestones renamed elsewhere: a -> b', JSON milestones_renamed), nicht unter 'filed' - dort stuende 'restore bringt es zurueck', und das stimmt fuer eine Umbenennung nicht.
8. Die Hilfe von jaira milestone sagt heute 'rename the file to rename the milestone' - auf einem geteilten Board falsch (der alte Ref schreibt a.md beim naechsten fetch zurueck). Wird durch den Befehl ersetzt.
- **2026-10-13 12:07 · BeMuCa** — in-progress: core Rename, refsync, CLI rename/fetch und Zwei-Klon-Test stehen. Gegenprobe: ohne die Marke auf dem alten Ref faellt TestARenamedMilestoneArrivesUnderItsNewNameOnly (grace hat dann first-round UND round-one). Gefunden beim Testen: die create-Tuer hatte nach 'else if err == nil && !have.Filed()' ein 'else if !errors.Is(err, os.ErrNotExist)', das bei err == nil (Marke auf der Platte) nil,nil zurueckgab -> Panic in Create; jetzt 'err != nil &&'. Dieselbe Kette steht in internal/tui/layer.go createMilestone - dort mitziehen.
- **2026-10-13 12:11 · BeMuCa** — TUI: der e-Dialog hat jetzt drei Felder (name, goal, after) als Slices statt [2]string; tab/↓ vor, shift+tab/↑ zurueck (vorher schalteten alle vier zwischen zwei Feldern hin und her). Achtung Tests: key("shift+tab") im Test-Helper erzeugt Code 0, also gar keine Taste - echtes shift+tab ist tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift}. Nach dem Umbenennen folgen form.edit, detailName, Cursor (sel) und ein Board-Filter milestone:<alt> dem neuen Namen. Docs: COMMANDS.md hat eine rename-Zeile; README Keys und COMMANDS Keys sagen 'e name, goal, predecessors'; die Spec-Zeile 'name not editable' ist ersetzt. Die D0SAHM-Release-Notiz bleibt unangetastet (fremdes Ticket), E70WYN.md beschreibt das Neue.
- **2026-10-13 12:17 · BeMuCa** — Nebenbefund (nicht in diesem Ticket behoben): jaira show --for-lane <lane> --json gibt unter input.notes von einer mehrzeiligen Notiz nur die erste Zeile weiter - die pre-process-Notiz mit den acht Entscheidungen kam bei critique nur als Ueberschrift an. In der Ticketdatei steht sie vollstaendig. Ab hier schreibe ich Notizen einzeilig. Der critique-Agent bekam die volle Progress-Sektion aus der Datei.
- **2026-10-13 12:28 · BeMuCa** — critique pass 1 (opus subagent, whole diff 6fa3639, full notes from the ticket file): one finding, clear fix - the free-name-on-disk switch exists three times (Rename, CLI create, TUI createMilestone); export it as milestone.CheckFree and call it from all three. All eight pre-process decisions left standing with their recorded reasons. Not a critique finding but a defect in this change, fixed in the same round: internal/tui/layer.go submitLayerForm normalizes the name field and compares with the raw f.edit, so a hand-named file like Q4.md would be renamed to q4 on any e save - compare the unchanged field text with f.edit first.
- **2026-10-13 12:29 · BeMuCa** — in-progress after critique 1: milestone.CheckFree(root, name) in core/milestone/write.go answers 'is this name free on disk' once (*ExistsError, *FiledError, nil for no file or a rename mark); Rename, the CLI create build and the layer's createMilestone call it and word the errors themselves - the three copies are gone, and the TUI lost its os import. submitLayerForm now normalizes the name field only when it differs from f.edit, so an unchanged hand name (Q4.md) is never renamed; TestEditKeepsAHandWrittenName fails with the old compare (checked by mutation).
- **2026-10-13 12:36 · BeMuCa** — critique pass 2 (same opus subagent, read only finding L1 and b1592e7): none. L1 addressed - CheckFree answers the four disk states (no file free, on board ExistsError, holding filing FiledError, rename mark free) exactly as the three removed copies did for all callers; the CLI rename's fe.Name == to check still matches. vet + go test of core/milestone, internal/cli, internal/tui, core/refsync ok (my full -race run before the commit: 28 ok).
- **2026-10-13 12:44 · BeMuCa** — optimize (opus subagent): two cleanups, behaviour unchanged - one const for rename's refusal clause, one slices.Index in renameMilestone. Everything considered and left is in review-gaps with its reason. Its remark that the filter does not follow a hand-named Q4 is not a defect of this change: core/milestone Index.Matches normalizes the filter value, so milestone:q4 never matched a file named Q4.md before either.
- **2026-10-13 12:52 · BeMuCa** — testing: tmux (scratch build of f840f0e, HOME/JAIRA_HOME in scratch, copy of fix2/board, no remote, 1 ticket in runde-1): L, enter, b -> filter milestone:runde-1; L, e -> form 'name runde-1 / goal / after'; name parallel -> 'milestone "parallel" already exists', form stays; name Fundament -> bubble fundament, runde-2/runde-2b still below it, files fundament.md, runde-2 after [fundament, parallel], runde-2b after [fundament], heading '# fundament'; L -> filter milestone:fundament, card still shown. From the detail: e, Tab x3 cycles name-goal-after-name, shift+Tab back; rename to basis -> detail shows basis, esc -> roadmap cursor on basis, board filter milestone:basis. CLI on the same copy: rename runde-2 'Zweite Runde' -> 'Filed as zweite-runde', runde-3 after rewritten; onto parallel exit 3; nope exit 5; same name 'Nothing to do'; bad/name exit 2; --json carries renamed_from. Two clones (ada, grace) + carol on a local bare repo: ada rename round-one first-round -> bare repo ref round-one = status: filed + renamed-to: first-round; grace fetch prints 'Milestones renamed elsewhere: round-one → first-round', ls lists only first-round and round-two, round-one.md lies marked, round-two after [first-round]; fresh carol gets first-round and round-two only; ada renames back -> grace lists round-one again, round-two after [round-one].
- **2026-10-13 13:16 · BeMuCa** — review pass 1 -> in-progress. Fix: refsync gets MilestoneOnRef(name) (ref carries a milestone on the board, not filed) and both name-giving doors ask it (CLI create + rename via one refuseIfTaken, TUI nameTaken) -> exit 3 'exists on its ref but not on this disk yet, jaira fetch brings it here'; this also closes the same pre-existing hole in create. Also cheap: help/COMMANDS exit codes for a filed old name (d), TUI reload when goal/after fails after a rename (e). Left for Berk, recorded in review-gaps: lost race sends the mark anyway (a), chained renames (b), no hint on the renaming tree (c), mark queue failure (f).
- **2026-10-13 13:18 · BeMuCa** — in-progress after review 1: refsync.MilestoneOnRef(name) - this tree's ref under name carries a milestone not filed; CLI refuseIfTaken (refuseIfFiled + that check, exit 3 'exists on its ref but not on this disk yet - jaira fetch brings it here') replaces refuseIfFiled in create and in rename's taken; TUI nameFiled renamed nameTaken with the same check. This also closes the hole create had before this ticket. TestANameOnlyOnARefIsTakenForRenameAndCreate (two clones) fails without the check - verified by mutation. Cheap minors done: rename help and COMMANDS give exit 3 for a filed old name; the TUI reloads when goal/after is refused after a rename went through (no test: needs a concurrent writer between the pre-check and the locked write). Release note line 5 added for the new refusal.
- **2026-10-13 13:29 · BeMuCa** — critique after review 1 (opus subagent, read only the review defect and 3c08f8a): none. MilestoneOnRef sits beside MilestoneFiled without overlap (filed ref / live ref / rename mark = neither); create and rename share refuseIfTaken (CLI) and nameTaken (TUI); the check cannot live in core CheckFree because core/milestone cannot import refsync; rename back and unshared boards unaffected.
- **2026-10-13 13:34 · BeMuCa** — optimize after review 1 (opus subagent, 3c08f8a only): one control-flow fold in submitLayerForm, behaviour unchanged; rest left with reasons in review-gaps.
