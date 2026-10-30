---
id: 01M5GFAV1C3DV9NWZFXC4G07X9
title: "Boards wie Lanes behandeln: von Hand anlegen, Lanes waehlen, beschreiben, teilen und global speichern"
status: review
ready: true
creator: berkc
assignee: berkc
goal: "Boards wie Lanes: ein Board laesst sich im TUI und per CLI anlegen (Name, Pflicht-Beschreibung, Lanes aus Board und Katalog), landet im persoenlichen Katalog ~/.jaira/boards, aus dem jedes Repo Boards anlegen kann, und laesst sich mit dem Team teilen (publish) und uebernehmen (adopt)."
context: |-
  Wunsch von berkc am 22.10.2026: 'Regarding the Boards we need to handle them similar to the lanes, making it possible to create a board manually and pick lanes to include, write board description and also to share boards with the team or save shared boards. The boards are also saved globally, like for the tickets.'
  Heute: ein Board ist ein Ordner .jaira/boards/<name>/lanes/ mit description-Datei, nur im Repo (S0EG0M). Anlegen geht nur von Hand.
  Vorbild Lanes: persoenlicher Katalog ~/.jaira/lanes (lane.UserLanesDir, core/lane/lane.go:466), publish nach .jaira/shared/<du>/, adopt zurueck (internal/tui/lanes.go:484, :640). Fuer Boards gaebe es analog ~/.jaira/boards und .jaira/shared/<du>/boards/.
  Verwandt: Y2C9QN (Board aus Board-Typ mit einem Befehl, Beschreibung Pflicht), 57TDND (aktive Boards), 28C0QV (neue Lane auch in den Katalog), Ticket zum Board-Fenster auf b.
  Offen: 'saved globally, like for the tickets' - Tickets liegen auf Refs im Repo, nicht global; gemeint ist vermutlich der persoenliche Katalog wie bei Lanes. Bei berkc nachfragen.
definition-of-done: "Im TUI laesst sich ein Board anlegen: Name, Beschreibung (Pflicht), Lanes aus Katalog und Board ankreuzen; Test"
tags:
  - tui
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T00:14:19Z
updated-at: 2026-10-30T02:29:39Z
updated-by: berkc
claimed-by: berk-ca-1397571
claimed-at: 2026-10-30T01:57:35Z
outcome-what: "Boards get a personal catalogue ~/.jaira/boards (core/lane/boardtypes.go UserBoardsDir; boards types/create --type read it), boards create --lanes picks lanes, a created board is saved to the catalogue, boards publish/shared/adopt share boards via .jaira/shared/<you>/boards/ (core/lane/boardshare.go), and B in the TUI lane settings screen opens a create-board form (internal/tui/newboard.go)."
outcome-why: "Boards could only be made by hand in one repo; berkc wants them handled like lanes: created with picked lanes and a description, kept globally in the personal catalogue so every repo can use them, and shared/adopted with the team."
outcome-resolves: "DoD 1 TUI form (TestLaneScreenCreatesABoardFromTheForm); 2 catalogue + every repo (TestBoardInTheCatalogueIsATypeInEveryRepository, JAIRA_HOME); 3 publish/adopt (TestPublishAndAdoptBoard); 4 core/release/notes/4G07X9.md; 5 CLI --json (TestBoardsPickSaveShareAndAdoptFromTheCLI). go vet and go test ./... green."
review-summary: "Boards now work like lanes. (1) Personal catalogue ~/.jaira/boards/<name>/ (honours JAIRA_HOME): each entry is a full board-type-shaped dir (order, description, every lane file); 'jaira boards types' lists it beside the shipped types (source column, shipped wins on a name clash) and 'boards create --type <name>' builds from it in any repo (core/lane/boardtypes.go UserBoardsDir, typeLanes/dirLanes generalised over an fs.FS). (2) 'boards create --lanes a,b,c' picks lanes from today's board, catalogue lanes and built-ins (CreateBoardWithLanes/PickableLanes; unknown/duplicate lane or --lanes with --type exit 2). (3) Every create (CLI and TUI) also saves the board into the catalogue unless a catalogue board or shipped type of that name exists (SaveToCatalogue; JSON gains catalogue, catalogue_saved). (4) 'boards publish|shared|adopt' copy a board to .jaira/shared/<you>/boards/<name>/, list those, and copy one into the catalogue, refusing overwrite without --force (exit 3) and refusing shipped type names (core/lane/boardshare.go). (5) TUI: B on the lane settings screen opens a form (name, required description, lane checklist with today's lanes ticked) that calls the same core functions (internal/tui/newboard.go, key routing in model.go). Release note core/release/notes/4G07X9.md and docs/COMMANDS.md updated."
review-gaps: |-
  note: if the board is created but SaveToCatalogue then fails (e.g. unwritable ~/.jaira), boards create returns a plain error (exit 1) although the board was written - the TUI reports the same case as a partial success; the CLI could warn instead of failing
  note: with --lanes every CreateBoardWithLanes error except ErrBoardExists is reported as code unknown_lane, including an invalid board name and I/O errors
  note: publishing today's board as "default" and adopting it puts a catalogue type literally called default (boards types lists it, --type default builds it) - works, but the name says nothing
  note: the TUI form keeps lane order as listed (today's, then catalogue/built-ins); there is no way to reorder picked lanes in the form - not asked by the DoD
  note: branch is based on feat/Y2C9QN-boards-create (not merged); origin/HEAD..HEAD shows the 2 Y2C9QN commits too - the PR must land after Y2C9QN's
test-verdict: "pass: go vet + go test -count=1 ./... green (RC=0); DoD 1-5 tests run; built binary exercised on two scratch repos with scratch JAIRA_HOME: create --lanes saved to catalogue, unknown lane exit 2, other repo lists catalogue type and builds clips from it, ticket filed on it, publish/shared/adopt onto a second home all with --json. TUI form verified only through TestLaneScreenCreatesABoardFromTheForm (real key dispatch), not by eye."
review-verdict: "Approve: all five DoD items are met by the diff, each with a test; go build, go vet and go test ./... pass; manual CLI run of create --lanes/types/publish/shared/adopt behaves as documented. No defects found; only minor notes."
review-round: 1
review-check: "1. cd into a scratch git repo, run 'git init' and 'jaira init', and 'export JAIRA_HOME=/tmp/jh' so your real ~/.jaira is untouched  2. jaira boards create video --lanes backlog,review,done --description \"Video work\"  -> says board created in .jaira/boards/video/lanes and 'Saved to your catalogue in /tmp/jh/boards/video'  3. jaira boards types -> a row 'video  catalogue  Video work' under 'code  shipped'  4. in a second scratch repo (same JAIRA_HOME): jaira boards create clips --type video --description Clips -> board with backlog → review → done  5. back in the first repo: jaira boards publish video -> path .jaira/shared/<you>/boards/video; run it again -> refused, exit 3  6. jaira boards shared -> lists video with its lanes and path  7. export JAIRA_HOME=/tmp/jh2; jaira boards adopt <that path> -> 'adopted /tmp/jh2/boards/video'  8. jaira (TUI) -> open the lane settings screen, press B -> 'New board' form; type a name, tab, a description, tab, space to untick a lane, enter -> message 'board <name> created with N lanes; saved to your catalogue, ...'  9. jaira boards -> the new board is listed with its description"
---

# Boards wie Lanes behandeln: von Hand anlegen, Lanes waehlen, beschreiben, teilen und global speichern

## Definition of Done

- [x] Im TUI laesst sich ein Board anlegen: Name, Beschreibung (Pflicht), Lanes aus Katalog und Board ankreuzen; Test
  proof: internal/tui/newboard.go formKey/createBoard; TestLaneScreenCreatesABoardFromTheForm
- [x] Ein Board landet auch im persoenlichen Katalog und ist in jedem Repo waehlbar; Test mit JAIRA_HOME
  proof: core/lane/boardshare.go SaveToCatalogue, boardtypes.go BoardTypes/typeLanes read UserBoardsDir; TestBoardInTheCatalogueIsATypeInEveryRepository (JAIRA_HOME)
- [x] publish/adopt fuer Boards wie fuer Lanes; Test
  proof: core/lane/boardshare.go PublishBoard/SharedBoards/AdoptBoard; TestPublishAndAdoptBoard
- [x] Release-Note-Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/4G07X9.md
- [x] Per CLI ohne TUI: 'jaira boards create --lanes a,b' waehlt Lanes, 'boards publish|shared|adopt' teilen und uebernehmen, alles mit --json; Test
  proof: internal/cli/boards.go create --lanes, publish/shared/adopt; TestBoardsPickSaveShareAndAdoptFromTheCLI

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/lane/boardtypes.go: UserBoardsDir (~/.jaira/boards), BoardTypes/typeLanes read shipped types then the catalogue; CreateBoardWithLanes picks lanes from today's board + Installable; SaveBoard copies a repo board into a board-type dir
- [x] core/lane: PublishBoard -> .jaira/shared/<you>/boards/<name>/, SharedBoards lists, AdoptBoard validates then copies into the catalogue
- [x] internal/cli/boards.go: create --lanes, create saves to catalogue, publish/shared/adopt, all --json
- [x] TUI: B in the lane settings screen opens a create-board form (name, description, lane checklist)
- [x] tests + release note + docs/COMMANDS.md

## Progress
- **2026-10-30 01:56 · berkc** — berkc 30.10.2026 (Frage 3/13): Boards und Default-Lanes sollen an einer Stelle global liegen, damit ein 'jaira update' sie fuer alle Sessions und Repos auf einmal aktualisiert und Repo-eigene Lanes nicht ueberschrieben werden. 'global wie Tickets' meint den persoenlichen Katalog ~/.jaira, nicht die Refs. Baut auf Y2C9QN auf (boards create/types, Branch feat/Y2C9QN-boards-create, noch nicht gemergt).
- **2026-10-30 01:59 · berkc** — Blocked-by Y2C9QN (reasoning): this work is built on branch feat/Y2C9QN-boards-create (jaira boards create/types, core/lane/boardtypes.go), not yet on master. Branch feat/4G07X9-boards-catalogue is based on origin/feat/Y2C9QN-boards-create; its PR must land after Y2C9QN's. Not set as blocked_by field so the ticket can run its lanes now.
- **2026-10-30 01:59 · berkc** — Brainstorm. What the code shows: a board is .jaira/boards/<name>/lanes/ (description, order, <id>.md) in one repo only (core/lane/boards.go). Y2C9QN adds board types: embedded dirs core/lane/boardtypes/<type>/ of the same shape, 'jaira boards create --type'. Lanes have a personal catalogue ~/.jaira/lanes (lane.UserLanesDir, honours JAIRA_HOME), 'jaira lanes publish' -> .jaira/shared/<you>/<id>.md, 'jaira lanes adopt <path>' -> catalogue (core/lane/share.go). Nothing for boards. Options: (A) catalogue board = a board-type-shaped dir ~/.jaira/boards/<name>/ with full lane copies; 'boards types' and 'boards create --type' read shipped types plus the catalogue; create also saves into the catalogue; publish/adopt via .jaira/shared/<you>/boards/<name>/. Reuses typeLanes/CreateBoard; repo boards stay copies, never overwritten. (B) repo boards reference the catalogue live (no copy): one edit reaches every repo, but a clone without the person's ~/.jaira shows no lanes - breaks 'teammate clones and sees same state'. (C) jaira update syncs catalogue boards into every repo board: needs shipped-version tracking (master's core/shipped, not on this base) and overwrites repo-owned lanes - against the 30.10 note. Choose A.
- **2026-10-30 01:59 · berkc** — Decisions settled (disagree at review): 1) 'global/jaira update for all repos' per 30.10 note = one place, ~/.jaira/boards: every repo's 'boards types' and 'boards create --type' read it live, and shipped types live in the binary, so updating the binary or the catalogue reaches every repo at once; repo boards are copies and are never rewritten. 2) A shipped type's name wins over a catalogue board of the same name (create does not save a catalogue copy under it), so a binary update is never shadowed. 3) CLI picks lanes with 'boards create --lanes a,b' (resolved: today's board's file, then catalogue, then built-in). 4) TUI: a create-board form opened with B from the lane settings screen (free key there), not from the b board window - that window and board switching belong to KVMG0K, which can open the same form later. 5) publish/adopt for boards are CLI only (boards publish/shared/adopt, all with --json); the TUI lane screen's p/a stay lane-only to keep the change small.
- **2026-10-30 02:00 · berkc** — Moved ideas -> brainstorm with --force: the ideas lane exists on master (.jaira/lanes/ideas.md) but not on the Y2C9QN base this worktree is on, so jaira here reads the ticket as in an unknown lane. No gate other than that was overridden.
- **2026-10-30 02:04 · berkc** — Found out: lane.Shared globs .jaira/shared/*/*.md, so published boards go one level deeper (shared/<who>/boards/<name>/) and are never read as shared lanes (asserted in TestPublishAndAdoptBoard). Catalogue boards hold every lane file (built-ins materialised), so a published/adopted board needs nothing else. Installable lists built-ins before catalogue lanes, so for a picked id the built-in wins over a catalogue lane of the same id - same as 'jaira lanes add'.
- **2026-10-30 02:13 · berkc** — critique pass 1: read whole diff. Let stand: catalogue save on every create (DoD 2 asks for it; a board built from catalogue type X under name Y adds Y - accepted, it is a copy the person named); clearDst in boardshare.go is the directory counterpart of copyLane's file-level overwrite refusal in share.go, not a duplicate; I/O errors from CreateBoardWithLanes reported as exit 2 - accepted, minor.
- **2026-10-30 02:29 · berkc** — review (fresh context): approve, no blockers. go build/vet/test ./... green; manual CLI run of boards create --lanes, types, publish/shared/adopt in a scratch repo with JAIRA_HOME matched the docs. Notes in review-gaps (CLI exit 1 when catalogue save fails after the board was written; unknown_lane code covers non-lane errors; adoptable catalogue type named default; no reorder in TUI form; PR must land after Y2C9QN).
