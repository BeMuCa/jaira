---
id: 01M6PDF2D6RQ8K1T7N0G1MYFKE
title: "Ein Milestone erreicht die Kollegen so schnell wie ein Ticket, ueber einen Ref"
status: in-progress
ready: true
creator: BeMuCa
goal: "Ein angelegter oder geaenderter Milestone liegt auf einem Ref wie Tickets und Epics und erscheint bei Kollegen mit dem naechsten Fetch, ohne Commit und Merge."
context: |-
  Was schief ist: Milestones leben nur als Dateien unter .jaira/milestones/. Sie erreichen andere erst ueber Commit, Push, Merge und Pull.
  Im jaira-Repo selbst ist auf master keine einzige Milestone-Datei getrackt (git ls-files .jaira/milestones: 0, Stand 2026-11-05). Die Roadmap eines Kollegen zeigt keinen unserer Milestones.
  Tickets liegen auf refs/jaira/tickets/*, Epics auf refs/jaira/epics (grep in core/), Milestones auf keinem Ref.
  Frage von Berk (2026-11-05): ab wann sieht ein anderer meinen Milestone auf seiner Roadmap? Antwort heute: erst nach Merge.
  Vorbild: die Epics - Beschreibung plus Ticketliste auf einem eigenen Ref (jaira epic --help).
  Achtung: update never loses data - vorhandene Dateien werden uebernommen, nicht verworfen.
definition-of-done: "Milestones liegen auf einem Ref; create, add, rm, rename, set schreiben ihn"
tags:
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-11-05T17:52:49Z
updated-at: 2026-11-21T20:45:35Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2112118
claimed-at: 2026-11-21T20:33:18Z
outcome-what: "jaira fetch (und der Board-Fetch) schickt eine ID-benannte Milestone-Datei ohne Ref und ohne wartenden Outbox-Eintrag einmal auf refs/jaira/milestones/<ID> und nennt sie (Text, JSON milestones_published, Board-Zeile)."
outcome-why: "Milestones liegen seit 0YGWXQ/EBTX6X schon auf Refs; die einzige Luecke war eine Datei, die vor dem Remote angelegt wurde - Migrate uebergeht ID-benannte Dateien, so blieb sie fuer immer auf einer Platte."
outcome-resolves: "Jeder Milestone dieses Boards erreicht die Kollegen mit dem naechsten Fetch, auch einer aus der Zeit ohne Remote."
review-summary: "core/refsync/refsync.go fetch(): der Merker war prozessweit und blieb nach einem guten Fetch true - im Board (Fetch per Timer) haetten spaetere Offline-Ticks wieder gesendet. Behoben: er heisst jetzt 'der letzte Fetch erreichte den Remote' (Store(err == nil)), Test TestAFailedFetchAfterAGoodOneSendsNothing. Sonst keine Befunde; Ort (FetchMilestones statt Migrate) und Wiederverwendung wie in Runde 1."
review-gaps: "nichts zu entfernen: ~40 Zeilen Kern in core/refsync/refsync.go publishUnsent, je eine Ausgabezeile in internal/cli/fetch.go und internal/tui/refs.go; der Zwillings-Ausschluss ist ein map-Zaehler und hat einen eigenen Test."
test-verdict: "pass: go vet ./... und GOOS=windows vet gruen, go test ./... RC=0, -race core/refsync RC=0; DoD 1-4 im Baum geprueft; Normalfall: Milestone vor dem Remote angelegt, erster jaira fetch meldet 'now sent to the remote', Ref auf dem Remote, zweiter fetch sendet nichts; Fehlerpfad: unerreichbarer Remote -> fetch exit 1, kein Ref, nichts gemeldet, nach Reparatur kommt er an; zwei fetch parallel: genau einer sendet m1,m2, je ein Commit pro Ref, Outbox leer"
hunt-verdict: "change: 2 Befunde (1 blocker) - publishUnsent sendet auch ohne erfolgreichen Fetch (ErrOffline), und ein nicht lesbarer Milestone-File bricht den ganzen fetch ab"
round: 3
back-from: critique
back-by: lane
back-count: 1
---

# Ein Milestone erreicht die Kollegen so schnell wie ein Ticket, ueber einen Ref

## Definition of Done

- [x] Milestones liegen auf einem Ref; create, add, rm, rename, set schreiben ihn
  proof: core/milestone/write.go:82 Edit/Rename/Create -> saveAndRecord -> refsync.RecordMilestone (seit 0YGWXQ/EBTX6X); E2E 2026-11-21: create/rename/set erschienen auf refs/jaira/milestones/<ID>
- [x] Ein zweiter Clone sieht einen neuen Milestone nach jaira fetch, ohne Commit
  proof: core/gitref/milestone_test.go TestMilestoneArrivesWithoutASharedBranch; core/refsync/publish_test.go TestAMilestoneOnThisDiskOnlyReachesTheOtherCloneOnTheNextFetch; E2E zwei Clones: jaira fetch in b zeigt q1/q3 ohne Commit
- [x] Vorhandene Milestone-Dateien werden beim ersten Lauf auf den Ref uebernommen; ein Test belegt es
  proof: core/refsync/refsync.go publishUnsent (aus FetchMilestones), Test core/refsync/publish_test.go TestAMilestoneOnThisDiskOnlyReachesTheOtherCloneOnTheNextFetch; alte Dateien ohne ID: core/milestone/migrate_test.go
- [x] Release-Notiz in core/release/notes/<ID>.md
  proof: core/release/notes/1MYFKE.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] refsync.FetchMilestones: nach Migrate jede ID-benannte Datei ohne lokalen Ref und ohne wartenden Outbox-Eintrag einmal mit RecordMilestone einreihen; eine ID auf zwei Dateien wird ausgelassen
- [x] MilestonesFetched.Published; jaira fetch (Text und JSON milestones_published) und Board-Zeile nennen sie
- [x] Test in core/refsync: Datei vor dem Remote angelegt, fetch, zweiter Clone sieht sie
- [x] Release-Notiz core/release/notes/1MYFKE.md

## Progress
- **2026-11-21 20:33 · Alexander Sacharov** — decision: Alex 21.11.: Milestones ganz auf Refs, keine Datei mehr, die man committen könnte — damit sie nie wieder hin und her committet werden. Auslöser: 21.11. hat ein logbook-Commit (git add -A .jaira) 12 untracked Milestone-Dateien nach master gebracht, PR #167 nimmt sie wieder raus und setzt /.jaira/milestones/ in .gitignore. Verwandt: KWQCYZ (Berk) — eine veraltete Datei überschreibt den Ref; ohne Datei als Quelle fällt diese Ursache weg, im Ticket erwähnen.
- **2026-11-21 20:33 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock
- **2026-11-21 20:35 · Alexander Sacharov** — Befund 2026-11-21: Die Praemisse im Kontext stimmt nicht mehr. Milestones liegen seit 0YGWXQ (Sept.) und EBTX6X (Okt.) auf refs/jaira/milestones/<ID>; create/add/rm/rename/set schreiben ihn, jaira fetch schreibt sie in einen zweiten Clone (E2E mit zwei Clones am 2026-11-21 nachgeprueft, ohne Commit). Alle 24 Milestones im jaira-Repo haben einen Ref auf upstream. Echte Luecke: eine ID-benannte Datei ohne Ref (angelegt bevor das Board einen Remote hatte) kommt nie auf den Ref - milestone.Migrate ueberspringt ID-benannte Dateien, nur alte Dateien ohne ID werden uebernommen. E2E: offline-XTTYQF.md blieb nach fetch ohne Ref, legacy.md kam an. Das ist der Rest von DoD 3.
- **2026-11-21 20:40 · Alexander Sacharov** — Nebenbefund, nicht dieses Ticket: ein mit 'jaira share' geteiltes Board ohne Tickets kommt im Clone nur als .jaira/.gitattributes an (leere Ordner trackt git nicht); jaira meldet dort 'no .jaira directory found; run jaira init'. Gesehen beim E2E am 2026-11-21, mkdir .jaira/tickets half.
- **2026-11-21 20:41 · Alexander Sacharov** — bughunt: change: blocker: core/refsync/refsync.go:520 (publishUnsent, called from FetchMilestones:487) - class: 'no local ref' is taken to mean 'no ref on the remote', but FetchMilestones also runs when no fetch happened: Incoming() returns nil,nil on ErrOffline (refsync.go:1009-1012), and both 'jaira fetch' (internal/cli/fetch.go:93) and the board's background fetch (internal/tui/refs.go:111) go on to FetchMilestones anyway. Any ID-named file whose ref this clone has not fetched yet is then queued with lease "". On the next flush the push is refused, and mergeRefused merges it against an empty base (no Base recorded, no migration commit), so every ticket in the stale copy counts as 'added' and every field the teammate changed is only a conflict comment. example: grace's milestone alpha-K3F9QZ has A,B on its ref; ada gets the same file through a merged branch; grace removes B; ada runs fetch while the remote is unreachable -> 'now sent to the remote: alpha'; back online, flush reports Merged:true Conflicts:[] and the ref has A,B again. The teammate's removal is undone and nobody is told. ran: scratch test in core/refsync (twoSides, FetchMilestones without a prior Repo.Fetch, then Flush): fails, B is back on the ref. Fix direction: publish only after a fetch that actually reached the remote (e.g. FetchMilestones takes or checks a 'fetched' flag, as migrateMilestones already does with refs.FetchRefs())
- **2026-11-21 20:41 · Alexander Sacharov** — bughunt: change: core/refsync/refsync.go:524-527 (publishUnsent) - class: an error reading one milestone file returns at once with a plain error, not a RecordError. That contradicts its own doc ('a file that could not be queued does not stop the others'). FetchMilestones then returns it, 'jaira fetch' exits with that error and prints nothing, not even arrivals (internal/cli/fetch.go:94-96), and the board's fetch shows the error on every tick. While old-format files are pending, migrateMilestones (internal/cli/milestones.go:663) also makes every writing milestone command fail this way. The files that sort after the bad one are never sent. milestone.LoadAll tolerates the same file, so before this change the fetch did not fail. example: a dangling symlink gone-K3F9QZ.md beside real-M7P2RT.md -> published [] and err 'open .../gone-K3F9QZ.md: no such file or directory' on every fetch; the same happens with a file without read permission, or one that git checkout removes between ReadDir and ReadFile (git takes no jaira lock). ran: scratch test TestScratchDanglingLink in core/refsync (LoadAll err nil, FetchMilestones err non-nil, Published empty)
- **2026-11-21 20:41 · Alexander Sacharov** — bughunt: done at fd04ce378 - 2 findings (2 change, 1 of them a blocker, 0 pre-existing) in 5 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/refsync/publish_test.go core/refsync/refsync.go core/release/notes/1MYFKE.md internal/cli/fetch.go internal/cli/milestones.go internal/tui/refs.go
- **2026-11-21 20:44 · Alexander Sacharov** — Runde 2 (Bughunt-Befunde behoben): (1) publishUnsent laeuft nur, wenn ein Fetch in diesem Prozess den Remote erreicht hat - Syncer.fetch() setzt dafuer ein atomic.Bool; alle Repo.Fetch-Aufrufe in refsync gehen jetzt durch ihn. Offline heisst 'kein lokaler Ref' nicht 'kein Ref auf dem Remote'. Test TestNothingIsSentWithoutAFetchThatReachedTheRemote. (2) eine nicht lesbare Datei wird uebersprungen und als RecordError (Warnung) gemeldet, die anderen gehen raus. Test TestAnUnreadableFileDoesNotStopTheOthers.
- **2026-11-21 20:45 · Alexander Sacharov** — Critique Runde 2: fetched-Merker auf 'letzter Fetch' umgestellt (vorher: irgendein Fetch im Prozess), sonst haette der Board-Timer nach einem guten Fetch offline wieder gesendet.
