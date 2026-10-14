---
id: 01M4QED6YSC1DXTB8TY5F3VNVC
title: move schreibt eine veraltete Ticket-Datei von der Platte ueber einen neueren Ref
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "jaira move und jaira pull arbeiten nie mit einer Ticket-Datei, die aelter ist als der Ref des Tickets"
context: |-
  Was passiert ist (12.10.2026): 'jaira move K1W2JR --to done --force' auf master. master hing upstream hinterher; die Datei .jaira/tickets/01M263RFDQNPFBZVAYJNK1W2JR-*.md war vom 10.09. (status backlog). Der Ref refs/jaira/tickets/01M263RFDQNPFBZVAYJNK1W2JR stand auf signoff, mit Review, test-verdict, DoD-Proof und 9 Notes (Commit 3290438).
  move las die Datei von der Platte und schrieb sie als d84afdb auf den Ref und nach upstream: Review, Outcome, Proof, Notes weg, assignee zurueck auf BeMuCa. Von Hand wiederhergestellt per Forward-Commit 5ca86a2.
  Die 5 Gate-Verweigerungen (assignee, 3x outcome leer, DoD offen) kamen alle aus der alten Datei - --force hat genau das Symptom verdeckt.

  Die Branch ist nicht sicherer: critique, testing und review committen nach Board-Regel nichts, also hinkt die Ticket-Datei jeder Branch dem Ref um diese Lanes hinterher (K1W2JR: Branch critique, Ref signoff; QTHD8H: Branch in-progress, Ref signoff).
  'jaira pull <id>' hilft nicht: es ist ein No-op, sobald irgendeine Datei da ist (core/refsync/refsync.go:733), es liest den Ref dann gar nicht.
  Ohne Datei verweigert move ('on its ref and not on your disk; pull it first').
  Umgangen wurde es mit einem Skript, das die Datei mit dem Inhalt des Refs ueberschreibt und erst dann move/logbook ruft.
definition-of-done: "move auf eine Datei, die aelter ist als ihr Ref, schreibt nicht auf den Ref, sondern verweigert (auch mit --force) und nennt beide Staende; ein Test pinnt das"
tags:
  - cli
  - concurrency
blocked-by: []
related: []
commits:
  - 8a09630ee4aa14863b9c2e58ac01ee254265fe2a
  - 5bf1e64bb95e7db0b1611482af2e197fada4d1eb
  - e12c9cdb280960aeb6be02c3b7cc907f9527ad97
created-at: 2026-10-12T06:57:07Z
updated-at: 2026-10-12T08:05:00Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-37581
claimed-at: 2026-10-12T07:24:00Z
outcome-what: "Review: Diff gegen die DoD geprueft, Tests mit -race gruen"
outcome-why: "Frische Review-Lane vor der Abnahme durch eine Person"
outcome-resolves: "F3VNVC: bereit fuer signoff, mit einem benannten Vorbehalt (Uhrenversatz)"
review-summary: "Jeder Schreibweg durch den Store (Mutate: note, dod, claim, set, move; dazu Archive und Logbook) fragt vor dem Schreiben den Syncer, ob updated-at der Datei strikt vor updated-at des lokalen Refs liegt; wenn ja, verweigert er mit StaleError (beide Staende: Lane + Zeit, Hinweis auf jaira pull), Exit 3, reason stale_file. move prueft das schon vor dem Gate, damit --force nicht die Luecken der alten Datei ueberstimmt. jaira pull auf eine vorhandene Datei holt best effort, vergleicht mit dem Ref und ueberschreibt eine aeltere Datei mit den Ref-Bytes (Umbenennung bei geaendertem Titel), ohne Ref oder Claim anzufassen; liegt ein ungesendeter Schreibvorgang im Outbox, verweigert pull (reason unsent). NOTES.md-Zeile vorhanden."
review-gaps: "1) Uhrenversatz kann falsch verweigern, entgegen der critique-Note ('nie eine falsche Verweigerung'): der lokale Ref bewegt sich erst nach erfolgreichem Push (Record, refsync.go:105). Schrieb ein Teammate mit vorgehender Uhr den Ref, und ich schreibe offline danach mit meiner Uhr (Zeit < seine), liegt meine Datei 'vor' dem Ref -> der naechste Befehl verweigert stale_file, und pull verweigert unsent: festgefahren bis online. Selten, aber der Kommentar an CheckFresh/die Note behaupten das Gegenteil. 2) refresh ueberschreibt Handaenderungen an der Datei, die nicht im Outbox liegen (nur Pending wird geprueft); ok, weil 'aelter' per Definition, aber nicht erwaehnt. 3) Kein Test fuer den Umbenennungs-Zweig in refresh (Titel auf dem Ref geaendert, alte Datei geloescht) und keiner fuer Archive/Logbook-Verweigerung. Keine Luecke gegen die DoD."
test-verdict: "pass: go clean -testcache; go test ./... -race -count=1 RC=0; DoD 1-3 im Baum geprueft (Proofs-Zeilen stimmen); frisch gebautes Binary auf Scratch-Board: move --force auf alter Datei RC=3 mit beiden Staenden, --json reason stale_file, Ref unveraendert; pull meldet was/now und die Datei traegt den Ref-Stand; note danach geht; zweites pull ist wieder no-op"
review-verdict: "DoD erfuellt: move --force verweigert eine aeltere Datei und nennt beide Staende (TestMoveRefusesAFileOlderThanItsRefEvenWithForce), pull frischt auf (TestPullBringsAStaleFileUpToItsRef), NOTES-Zeile da, go test ./... -race -count=1 gruen. Annehmbar mit Vorbehalt: der Uhrenversatz-Fall (gaps 1) ist ein echter, seltener Fehlweg; entweder hinnehmen und die Note/den Kommentar korrigieren, oder spaeter eigenes Ticket."
review-check: "1. cd in einen Klon mit Remote, ein Ticket X mit Ref nehmen.  2. Datei von X auf alten Stand bringen: git show <alter-commit>:.jaira/tickets/<X>.md > .jaira/tickets/<X>.md  3. jaira move X --to done --force -> Ausgabe 'the ticket file here is older than its ref', Zeilen 'file:' und 'ref:', Exit-Code 3 (echo $?); git rev-parse refs/jaira/tickets/<ULID> unveraendert.  4. jaira note X test -> dieselbe Verweigerung.  5. jaira pull X -> 'was behind its ref and is up to date now' mit 'was:'/'now:'.  6. jaira note X test -> geht durch. Automatisch: go test ./internal/cli -run 'Stale|OlderThanItsRef' und go test ./core/refsync -run Unsent."
---

# move schreibt eine veraltete Ticket-Datei von der Platte ueber einen neueren Ref

## Definition of Done

- [x] move auf eine Datei, die aelter ist als ihr Ref, schreibt nicht auf den Ref, sondern verweigert (auch mit --force) und nennt beide Staende; ein Test pinnt das
  proof: core/refsync/refsync.go:727 CheckFresh, core/move/move.go:95 vor dem Gate, core/ticket/store.go:1300 in Mutate; TestMoveRefusesAFileOlderThanItsRefEvenWithForce
- [x] jaira pull bringt eine vorhandene, aeltere Datei auf den Stand des Refs, statt still nichts zu tun; ein Test pinnt das
  proof: core/refsync/refsync.go:787 refresh; TestPullBringsAStaleFileUpToItsRef, TestPullLeavesAStaleFileWithAnUnsentWriteAlone
- [x] Zeile in core/release/NOTES.md; go test ./... -race gruen
  proof: core/release/NOTES.md:17; go test ./... -race -count=1 exit 0

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Syncer: lokalen Ref lesen, updated_at mit der Datei vergleichen, beide Staende zurueckgeben
- [x] Store.Mutate fragt den Recorder vor fn und verweigert mit ErrStale (beide Staende, Hinweis auf jaira pull)
- [x] Syncer.Pull: vorhandene aeltere Datei auffrischen, Outbox-Pending verweigert
- [x] Tests: move --force auf alter Datei verweigert, Ref unveraendert; pull frischt auf
- [x] NOTES.md-Zeile, go test ./... -race

## Progress
- **2026-10-12 07:24 · Alexander Sacharov** — parent session: jaira-4c
- **2026-10-12 07:25 · Alexander Sacharov** — Offene Entscheidungen vor dem Plan: keine, die nur Alex treffen kann. Selbst entschieden, jeweils mit Grund:
1. 'Aelter' heisst: updated_at der Datei liegt strikt vor updated_at des lokalen Refs. Der lokale Ref ist genau der Stand, gegen den der Lease beim Push prueft (Syncer.Record), und kostet keinen Netzwerkweg. Byte-Vergleich mit der Ref-Historie verworfen: eine Datei, die nie so auf dem Ref stand (K1W2JR: Stand vom 10.09., Ref erst spaeter angelegt), faellt da durch.
2. Die Sperre sitzt in Store.Mutate, nicht nur in move: jeder Schreibweg (note, dod, claim, set, move) liest die Datei und schreibt sie auf den Ref, also verliert jeder denselben Stand. Mutate fragt den Recorder ueber ein optionales Interface; --force kommt dort gar nicht an, also greift die Sperre auch mit --force.
3. pull auf eine vorhandene Datei: fetch, dann Ref-updated_at vergleichen. Ref neuer -> derselbe Weg wie ein frischer pull (Ref uebernehmen, Datei mit den Ref-Bytes ueberschreiben; alter Dateiname weg, falls der Titel sich geaendert hat). Ref gleich alt oder aelter -> wie bisher AlreadyHere. Liegt fuer das Ticket ein ungesendeter Schreibvorgang im Outbox, verweigert pull das Ueberschreiben statt lokale Arbeit zu verwerfen.
- **2026-10-12 07:32 · Alexander Sacharov** — Umsetzung, was der Code nicht sagt:
- Pruefung auch in Store.Archive und Store.Logbook: beide schreiben per recordFiled den Dateistand auf den Ref, also derselbe Verlust. Delete bewusst nicht: Loeschen ist Absicht.
- move.Move prueft VOR dem Gate, nicht nur in Mutate: sonst listet das Gate die Luecken der alten Datei (bei K1W2JR die 5 Verweigerungen) und --force ueberstimmt sie, bevor Mutate verweigert.
- pull-Refresh schreibt den Ref nicht und aendert den Claim nicht; er laeuft auch offline gegen den lokalen Ref (fetch nur best effort), damit wiederholtes pull offline weiter geht.
- pull-Refresh nimmt keine Store-Sperre (Store.lock ist privat); bei gleichzeitigem Mutate desselben Tickets koennte es sich ueberschneiden. Fuer einen Handbefehl hingenommen.
- Ein Ref ohne updated-at gibt kein Urteil (Store.Create setzt es nicht von selbst, die CLI schon) - dann bleibt alles wie vorher.
- Test-Falle: runCLI ruft afterCommand nicht, die Ref-Schreibvorgaenge bleiben im Outbox; flushRefs() im Test aufrufen.
- **2026-10-12 07:35 · Alexander Sacharov** — critique 1: refresh soll y.Store.Load(id) statt eigener ReadFile/ParseDoc/Decode-Kette nutzen (wie der alte AlreadyHere-Zweig). Bewusst stehen gelassen: Was/Now als fertige Saetze im JSON (die Verweigerung soll Menschen beide Staende nennen; ein Agent braucht nur reason=stale_file und pull); Ueberschreiben statt Merge im Refresh (Begruendung im Kommentar an refresh); doppelte Pruefung in move.Move und Mutate (vor dem Gate noetig, in Mutate fuer alle anderen Schreibwege). Uhrenversatz zwischen Maschinen fuehrt hoechstens zu falsch-negativ, also zum alten Verhalten, nie zu einer falschen Verweigerung, solange der Ref nur Inhalte traegt, die die Datei nicht hat.
- **2026-10-12 07:40 · Alexander Sacharov** — testing: done, an die frische review-Lane uebergeben (eigener Herdr-Tab via run-lane.sh --no-worktree).
- **2026-10-12 07:45 · Alexander Sacharov** — review: DoD erfuellt, Tests gruen. Vorbehalt fuer die Abnahme: Uhrenversatz zwischen Maschinen kann offline falsch verweigern (lokaler Ref bewegt sich erst nach Push; Ref mit Zeitstempel aus der Zukunft -> eigene spaetere Schreibvorgaenge wirken 'aelter', pull verweigert dann unsent). Die critique-Note behauptet das Gegenteil. Entscheidung fuer die Person: hinnehmen oder nachbessern.
- **2026-10-12 08:02 · Alexander Sacharov** — Alex 12.10. zu review gap 1 (Uhrzeitversatz -> falsche Verweigerung): annehmen; im logbook-Commit den CheckFresh-Kommentar und die critique-Note korrigieren - 'nie eine falsche Verweigerung' stimmt nicht, Uhrzeitversatz kann eine erzeugen.
- **2026-10-12 08:04 · Alexander Sacharov** — Korrektur zur critique-Note: 'nie eine falsche Verweigerung' stimmt nicht. Uhrenversatz kann eine aktuelle Datei verweigern: Teammate mit vorgehender Uhr schreibt den Ref, hier wird danach offline mit frueherem Stempel geschrieben; der lokale Ref bewegt sich erst nach einem Push, also meldet der naechste Schreibbefehl stale_file und pull unsent, bis das Remote erreichbar ist. Von Alex so angenommen (Option 1); der Kommentar an CheckFresh (core/refsync/refsync.go) sagt es jetzt.
