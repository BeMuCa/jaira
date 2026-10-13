---
id: 01M4THCXMAGS2JQN133DE70WYN
title: Der Name einer Milestone laesst sich auf der Milestone Roadmap und per CLI aendern
status: critique
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
updated-at: 2026-10-13T12:34:14Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-1775362
claimed-at: 2026-10-13T12:29:45Z
executed-by: opus
outcome-what: "Critique 1 answered: milestone.CheckFree (core/milestone/write.go) is the one 'is this name free on disk' check, called by Rename and by both create doors (CLI create, TUI createMilestone). The e-form keeps an unchanged hand-written name (Q4.md) instead of renaming it to its normalized form."
outcome-why: "The same three-way switch was written three times and a nil,nil bug had to be fixed twice; normalizing an unchanged name renamed hand-named files on every goal edit."
outcome-resolves: "DoD1-3 unchanged and still covered by their tests; TestEditKeepsAHandWrittenName pins the hand-name case. vet + go test ./... -race green (28 ok)."
review-summary: "core/milestone/write.go:150-157, internal/cli/milestones.go:127-141, internal/tui/layer.go:398-405 write the same 'is this name free on disk' switch three times (on the board -> exists, filed and HoldsName -> refused, rename mark -> free, other error -> returned); the nil,nil panic of the 12:07 note had to be fixed in two copies. Export it once as milestone.CheckFree(root, name) error returning *ExistsError / *FiledError / nil, call it from Rename and from both create build funcs, each door wording the two errors itself."
---

# Der Name einer Milestone laesst sich auf der Milestone Roadmap und per CLI aendern

## Definition of Done

- [x] jaira milestone rename <alt> <neu> benennt Datei und Ref um, schreibt die after:-Listen aller Nachfolger um und lehnt einen vergebenen oder abgelegten Namen ab; Tests in internal/cli
  proof: core/milestone/write.go Rename; internal/cli/milestones.go newMilestoneRenameCmd; TestRenameMovesTheMilestoneAndItsSuccessorsFollow, TestRenameRefusalsAndTheMarkOfARename, TestRenameMovesTheFileRewritesSuccessorsAndMarksTheOldName, TestRenameRefusesATakenNameButNotTheMarkOfARename
- [x] Im e-Dialog der Milestone Roadmap ist der Name das erste Feld; ein geaenderter Name benennt ueber denselben Weg um; TUI-Test
  proof: internal/tui/layer.go startLayerEdit (name first), submitLayerForm -> renameMilestone -> milestone.Rename; TestEditRenamesAMilestone, TestEditGoalAndPredecessors
- [x] Ein zweiter Klon, der nach dem Umbenennen fetcht, hat danach nur die Milestone unter dem neuen Namen und nicht zusaetzlich die alte; Test mit zwei Repos wie in den refsync-Tests
  proof: internal/cli/milestoneref_test.go TestARenamedMilestoneArrivesUnderItsNewNameOnly (two clones, bare remote; fails without the mark on the old ref)
- [ ] go vet und go test ./... -race gruen; im echten Terminal auf einer Board-Kopie ohne Remote durchgespielt; Release-Notiz unter core/release/notes/

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
- [ ] Binary bauen, in tmux auf Board-Kopie ohne Remote durchspielen; Zwei-Klon-Probe mit lokalem bare-Repo

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
