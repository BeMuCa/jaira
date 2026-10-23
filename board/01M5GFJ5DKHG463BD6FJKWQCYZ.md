---
id: 01M5GFJ5DKHG463BD6FJKWQCYZ
title: Eine veraltete Milestone-Datei ueberschreibt beim naechsten milestone add still den neueren Stand auf dem Ref
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Ein Schreiben auf einen Milestone aus einer Datei, die aelter als ihr Ref ist, wird verweigert oder zuerst nachgezogen - es gehen nie Tickets aus einem Milestone verloren."
context: |-
  Passiert am 22.10.2026: im Worktree ../jaira-T8A8KM lagen alte Kopien unter .jaira/milestones/ (untracked, vom 14.10.). jaira milestone add jaira-board-fan-out KVMG0K 4G07X9 7BWHCN schrieb die alte Datei plus die drei neuen auf den Ref und pushte. Verloren: 57TDND und Y2C9QN (board-fan-out), PP24H1 und XR27FN (erweiterungen-abschauen). Keine Warnung.
  Wiederhergestellt aus der Ref-Historie (git show refs/jaira/milestones/<ID>~1), indem aus dem Haupt-Checkout erneut hinzugefuegt wurde.
  Ursache (aus dem Code gelesen, nicht per Test bestaetigt): core/refsync/refsync.go:204 RecordMilestone nimmt als Lease den SHA des lokalen Refs (y.Repo.MilestoneSHA). Der lokale Ref ist nach dem Fetch aktuell, also passt der Lease immer - ob die Datei auf der Platte zu diesem Ref passt, prueft niemand. Fuer Tickets gibt es genau dafuer CheckFresh/ErrStale (core/ticket/store.go:218).
  Gleiches Risiko in jedem Worktree mit alten Milestone-Kopien: ../jaira-K6BGB0, ../jaira-V38STG, ../jaira-X3A8RZ.
definition-of-done: "Ein milestone add/rm/set/rename aus einer Datei, die aelter als ihr Ref ist, schreibt nicht auf den Ref (Fehler mit Hinweis) oder liest vorher den Ref-Stand ein; Test in core/refsync oder internal/cli"
tags:
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T00:18:19Z
updated-at: 2026-10-22T00:18:19Z
---

# Eine veraltete Milestone-Datei ueberschreibt beim naechsten milestone add still den neueren Stand auf dem Ref

## Definition of Done

- [ ] Ein milestone add/rm/set/rename aus einer Datei, die aelter als ihr Ref ist, schreibt nicht auf den Ref (Fehler mit Hinweis) oder liest vorher den Ref-Stand ein; Test in core/refsync oder internal/cli
- [ ] Der Test spielt den Fall von oben nach: alte Datei, neuerer Ref, add -> kein Ticket geht verloren
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

