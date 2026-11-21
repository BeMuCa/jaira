---
id: 01M6PDF2D6RQ8K1T7N0G1MYFKE
title: "Ein Milestone erreicht die Kollegen so schnell wie ein Ticket, ueber einen Ref"
status: critique
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
updated-at: 2026-11-21T20:38:15Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2112118
claimed-at: 2026-11-21T20:33:18Z
outcome-what: "jaira fetch (und der Board-Fetch) schickt eine ID-benannte Milestone-Datei ohne Ref und ohne wartenden Outbox-Eintrag einmal auf refs/jaira/milestones/<ID> und nennt sie (Text, JSON milestones_published, Board-Zeile)."
outcome-why: "Milestones liegen seit 0YGWXQ/EBTX6X schon auf Refs; die einzige Luecke war eine Datei, die vor dem Remote angelegt wurde - Migrate uebergeht ID-benannte Dateien, so blieb sie fuer immer auf einer Platte."
outcome-resolves: "Jeder Milestone dieses Boards erreicht die Kollegen mit dem naechsten Fetch, auch einer aus der Zeit ohne Remote."
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
