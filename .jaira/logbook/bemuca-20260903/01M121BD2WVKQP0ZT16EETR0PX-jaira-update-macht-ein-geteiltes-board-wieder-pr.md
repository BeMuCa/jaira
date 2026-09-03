---
id: 01M121BD2WVKQP0ZT16EETR0PX
title: jaira update macht ein geteiltes Board wieder privat
status: done
ready: true
creator: BeMuCa
goal: jaira update auf einem geteilten Board laesst .gitignore in Ruhe; nur jaira init schreibt den Ignore-Eintrag
context: |-
  Gemessen am 27.08. auf zwei Wegwerf-Boards: 'jaira init && jaira share && jaira update' - nach share steht nur /.jaira/lanes/ in .gitignore, nach update steht /.jaira/ wieder drin. Dasselbe, wenn die Zeile von Hand auskommentiert wurde (#/.jaira/), wie auf Berks req-Board. Folge: neue Tickets eines geteilten Boards werden ignoriert, Teammates sehen sie nicht mehr; schon getrackte Dateien bleiben getrackt, deshalb faellt es erst spaet auf.

  Ursache: internal/cli/update.go:82 ruft board.Prepare, und Prepare (core/board/board.go:20) ruft AddIgnore ohne Bedingung. AddIgnore prueft nur Ignored(root) (core/board/gitignore.go:29), und das erkennt eine auskommentierte Zeile nicht als Absicht.

  Was NICHT stimmt: die Hilfe von update sagt 'Re-runs the same setup init performs - the .gitignore entry and the jaira section'. Der Ignore-Eintrag ist die Haelfte, die schadet.

  Fix-Vorschlag: update ruft nur AnnounceInAgentFiles (den Block) und, wenn das Board geteilt ist, AddLanesIgnore; den /.jaira/-Eintrag schreibt nur init. Alternativ Prepare ein Flag geben. Tests in internal/cli, die 'gitignore_written' von update erwarten, anpassen.

  Warum das jetzt zaehlt: der board-beschreibende Block (0e7db74) ist auf Berks req-Board nie angekommen - dort steht noch der alte Block ohne 'This board's lanes' und ohne 'Work one lane to empty'. Ihn per jaira update zu holen wuerde das Board heute wieder privat machen.
definition-of-done: "Nach 'jaira init && jaira share && jaira update' enthaelt .gitignore keinen /.jaira/-Eintrag (Test); auf einem Board mit auskommentierter Zeile ebenso; jaira init schreibt ihn weiterhin; die update-Hilfe sagt nicht mehr, dass sie den Ignore-Eintrag schreibt"
blocked-by: []
commits:
  - ef8f2d9b946bbbbef8d29940ac818c60cb5d0dbd
  - 6c5d81b28e18049c24bafd21a1d639c263474a32
created-at: 2026-08-27T16:37:48Z
updated-at: 2026-09-03T15:48:26Z
claimed-by: EE-3NX6GL3-3987890
claimed-at: 2026-08-27T20:41:37Z
updated-by: BeMuCa
assignee: BeMuCa
outcome-what: "update ruft nur noch board.AnnounceInAgentFiles; Prepare (Ignore + Note) bleibt init und dem Browse-Screen. JSON ohne gitignore_written, Ausgabezeile 'This board is private' weg, Hilfe erklaert es."
outcome-why: "Ob ein Board geteilt ist, ist beim update laengst entschieden; den Eintrag wieder zu schreiben liess neue Tickets still im Ignore verschwinden."
outcome-resolves: "DoD: der Test schreibt die auskommentierte Zeile, laeuft update, vergleicht byteweise - unveraendert. init schreibt den Eintrag weiter (bestehende init-Tests). go test ./internal/cli gruen."
review-summary: "Ein Commit: update.go ruft AnnounceInAgentFiles statt Prepare, JSON-Feld gitignore_written entfaellt, Hilfetext. Ein Test."
review-gaps: "Nichts zu entfernen: Prepare hat weiter zwei Aufrufer (init, browse). Der update-Test fuer den Fall 'jaira share, dann update' fehlt, weil share im Test ein git-Repo braucht; der Handvariante-Test deckt denselben Pfad (AddIgnore wird nicht mehr gerufen)."
review-check: |-
  1. Wegwerf-Repo: git init && jaira init && jaira share - .gitignore enthaelt nur /.jaira/lanes/.
  2. jaira update - .gitignore ist unveraendert (vorher kam /.jaira/ zurueck).
  3. jaira update --help: der zweite Absatz sagt, warum .gitignore nicht angefasst wird.
  4. go test ./internal/cli -run TestUpdate -v: alle PASS.
review-verdict: "Angenommen; dieselbe Sitzung, kein zweites Modell. Geprueft: cli-Tests gruen, Binary gebaut. Verhaltensaenderung: update schreibt nie mehr in .gitignore - auch nicht auf einem privaten Board, wo es vorher ein No-op war (Eintrag vorhanden)."
---

# jaira update macht ein geteiltes Board wieder privat

## Definition of Done

- [x] Nach 'jaira init && jaira share && jaira update' enthaelt .gitignore keinen /.jaira/-Eintrag (Test); auf einem Board mit auskommentierter Zeile ebenso; jaira init schreibt ihn weiterhin; die update-Hilfe sagt nicht mehr, dass sie den Ignore-Eintrag schreibt
  proof: internal/cli/update.go: AnnounceInAgentFiles statt Prepare; TestUpdateLeavesASharedBoardShared (auskommentierte Zeile bleibt unveraendert); Hilfetext nennt den Grund

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

