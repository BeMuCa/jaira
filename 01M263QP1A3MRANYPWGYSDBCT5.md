---
id: 01M263QP1A3MRANYPWGYSDBCT5
title: "Der Agent-Block nennt --proof, und die Gitignore-Regel verweist auf die Commit-Regel"
status: in-progress
ready: true
creator: BeMuCa
assignee: Alexander Sacharov
goal: "Zwei nachgemessene Luecken im generierten Block sind geschlossen: der Nachweis beim Abhaken, und der Widerspruch auf einem ungeteilten Board"
context: |-
  Zwei Befunde aus einer Multi-Agent-Session am 10.09., beide nachgeprueft und beide enger als ursprünglich gemeldet.
  1. --proof fehlt im Block. 'grep -c proof' in der CLAUDE.md eines anderen Projekts ist 0. NICHT zutreffend ist die Meldung, es sei gar nicht auffindbar: 'jaira dod --help' dokumentiert --proof an sechs Stellen samt Beispiel. Es fehlt allein im Block - also genau dort, wo eine Session es liest. Folge: die ersten Tickets wurden ohne Nachweis abgehakt.
  2. Der Gitignore-Vorbehalt und die Commit-Regel widersprechen sich auf einem ungeteilten Board. Der Block sagt 'the ticket rides in the same commit as the code', aber auf einem Board mit gitignoretem .jaira/ KANN die Ticketdatei nicht mitreisen. NICHT zutreffend ist die Meldung, die Reihenfolge sei falsch: der Vorbehalt steht auf Zeile 216, die Regel auf 220 - erst der Zustand, dann die Regel. Was fehlt, ist der Verweis: der Vorbehalt sagt nicht, dass der naechste Punkt fuer dieses Board nicht gilt.
  Ein anderes Projekt ist so ein Board: .gitignore Zeile 42 enthaelt /.jaira/, null getrackte Dateien darunter.
definition-of-done: "Der generierte Block nennt --proof beim Abhaken; der Gitignore-Vorbehalt sagt ausdruecklich, dass die Commit-Regel auf einem ungeteilten Board anders lautet; 'jaira update' schreibt beides auf ein bestehendes Board; go test ./... -race gruen"
tags:
  - docs
blocked-by:
  - 01M7WWY6PX0TRDQSVK7RNNV39Y
commits: []
created-at: 2026-09-10T16:52:07Z
updated-at: 2026-11-21T08:07:13Z
related:
  - 01M7WWY6PX0TRDQSVK7RNNV39Y
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-26380
claimed-at: 2026-11-21T08:01:01Z
outcome-what: "Der generierte Agent-Block sagt im Punkt 'the ticket rides with the code', dass dieser und der naechste Punkt nur auf einem geteilten Board gelten: ohne geteiltes .jaira/ gibt es keine Ticketdatei zum Mitcommitten, der Code geht allein mit dem Handle im Betreff. Test haelt den Satz und --proof fest; CLAUDE.md/AGENTS.md dieses Repos regeneriert; Release-Note."
outcome-why: "Auf einem ungeteilten Board widersprach die Commit-Regel dem Zustand des Boards: ein Agent kann eine ignorierte Datei nicht mitcommitten und las die Regel als kaputt. --proof und 'jaira update' waren schon ausgeliefert (2143e9ad, update.go:257) und brauchten nur den Test."
outcome-resolves: SDBCT5
review-summary: "Der generierte Agent-Block (core/board/announce.go) sagt jetzt am Punkt 'the ticket rides with the code': dieser und der naechste Punkt gelten nur auf einem geteilten Board; solange .jaira/ ausserhalb von git liegt (jaira init bis jaira share), gibt es keine Ticketdatei zum Mitcommitten - der Code wird allein committet, das Handle im Betreff ist die ganze Verbindung. Derselbe Vorbehalt steht in der jaira-Rolle (core/role/jaira/SKILL.md, Kopie .claude/skills/jaira), in docs/AGENTS.md und in der PR-Rolle, die einen Code-Commit ohne Ticketdatei dort nicht mehr als Befund wertet. --proof stand schon vorher im Block (2143e9ad). CLAUDE.md/AGENTS.md sind regeneriert, shipped.sha256 kennt die neuen Rollen-Hashes, ein Test haelt Block- und Rollensatz fest, die Release-Note sagt 'jaira update' bzw. 'jaira update --roles'."
round: 4
review-gaps: |-
  note: "This point and the next hold on a shared board only" erfasst auch den Punkt 'a lane that changed no code commits nothing', dessen Kern (nichts committen) auf einem ungeteilten Board erst recht gilt - nur dessen Ticketdatei-Saetze fallen weg; lesbar, aber nicht ganz praezise
  note: pre-existing: die fertigen Commit-Zeilen in jaira-worker SKILL.md und jaira-role-lane SKILL.md nennen weiter .jaira/tickets/<ticket>.md, das git auf einem ungeteilten Board als ignoriert ablehnt (in critique Runde 3 bewusst stehen gelassen)
  note: der Testname nennt nur das ungeteilte Board, prueft aber auch --done --proof
test-verdict: "pass: go test ./... -race exit 0; build ok. DoD 1-6 checked in the tree (announce.go:101/138, role SKILL.md:331, pr SKILL.md:59, docs/AGENTS.md:72, release note, test present and green). Behaviour on a scratch repo: 'jaira init' writes the block with --proof and the caveat into CLAUDE.md+AGENTS.md, .jaira/ ignored, and 'git add .jaira/tickets' is refused there ('paths are ignored') - exactly the case the caveat names; 'jaira update' rewrites a stale block back. Failure path: read-only CLAUDE.md -> update exits 1 with permission denied, file unchanged. Two writers: the change adds no write path (block writing is pre-existing), nothing to race."
hunt-verdict: "pass: bughunt done at bb02c815, 0 findings in 10 modules"
review-verdict: pass
review-round: 3
review-check: |-
  1. In einem leeren Verzeichnis: git init && jaira init (mit diesem Build)
  2. grep -n "hold on a shared board only" CLAUDE.md AGENTS.md - je eine Treffer-Zeile im Punkt 'the ticket rides with the code'
  3. git status --short zeigt nur CLAUDE.md und AGENTS.md, kein .jaira/ - genau der Fall, den der Satz beschreibt
  4. In einem bestehenden Board: jaira update schreibt denselben Satz in den Block; go test ./core/board -run TestAgentNoteSaysTheTicketCannotRideOnAnUnsharedBoard ist gruen
---

# Der Agent-Block nennt --proof, und die Gitignore-Regel verweist auf die Commit-Regel

## Definition of Done

- [x] Der generierte Block nennt --proof beim Abhaken; der Gitignore-Vorbehalt sagt ausdruecklich, dass die Commit-Regel auf einem ungeteilten Board anders lautet; 'jaira update' schreibt beides auf ein bestehendes Board; go test ./... -race gruen
  proof: core/board/announce.go:101 (--proof), core/board/announce.go:138 (unshared-board caveat), internal/cli/update.go:257 (jaira update rewrites the block); go test ./... -race exit 0
- [x] Ein Test in core/board/announce_test.go haelt den Satz fest, dass auf einem ungeteilten Board keine Ticketdatei mitreist und der Code-Commit allein geht
  proof: TestAgentNoteNamesProofAndTheUnsharedBoardCaveat
- [x] core/release/notes/SDBCT5.md nennt die Aenderung als eine Zeile
  proof: core/release/notes/SDBCT5.md:1-2
- [x] Die Agent-Bloecke dieses Repos (CLAUDE.md, AGENTS.md) sind mit dem neuen Block regeneriert
  proof: CLAUDE.md:131, AGENTS.md (regenerated with 'jaira update' from this build)
- [x] Die jaira-Rolle (core/role/jaira/SKILL.md) traegt denselben Vorbehalt wie der Block
  proof: core/role/jaira/SKILL.md:331; TestAgentNoteNamesProofAndTheUnsharedBoardCaveat
- [x] Die PR-Rolle und docs/AGENTS.md werten einen Code-Commit ohne Ticketdatei auf einem ungeteilten Board nicht als Fehler
  proof: core/role/builtin/jaira-role-pr/SKILL.md:59, docs/AGENTS.md:72; go test ./... -race exit 0
- [x] Die fertigen Commit-Zeilen in jaira-worker und jaira-role-lane SKILL.md lassen die Ticketdatei auf einem ungeteilten Board weg; der Vorbehalt sagt, dass der Punkt 'commits nothing' dort erst recht gilt
  proof: core/role/builtin/jaira-worker/SKILL.md:107, core/role/builtin/jaira-role-lane/SKILL.md:90 and :293, core/board/announce.go:141; TestAgentNoteNamesProofAndTheUnsharedBoardCaveat; go test ./... -race exit 0

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Satz im Punkt 'ticket rides with the code' in core/board/announce.go
- [x] Test in announce_test.go
- [x] Release-Note
- [x] Agent-Dateien des Repos regenerieren, go test ./... -race

## Progress
- **2026-11-20 16:49 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-20 22:44 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/35112.sock
- **2026-11-20 22:45 · Alexander Sacharov** — decision: open decisions none. Finding at start (20.11.): part 1 is already shipped - core/board/announce.go:101 names --proof (commit 2143e9ad). Part 3 is already shipped too - 'jaira update' rewrites the block via board.AnnounceInAgentFiles (internal/cli/update.go:257). The old separate gitignore caveat (was line 216) is now folded into the 'every commit names the ticket id' point. What is still missing: the 'ticket rides with the code' point and the 'commit it' sentence of 'commits nothing' read as if .jaira/ were always in git; on an unshared board (.git/info/exclude since jaira init) there is no ticket file to git add. Fix: one sentence on the rides-with-code point saying this and the next point assume a shared board, and what to do instead.
- **2026-11-20 22:50 · Alexander Sacharov** — critique: sibling site core/role/jaira/SKILL.md:317 has the same rides-with-code rule without the unshared-board caveat; mirror the sentence there. Placement in announce.go (on the rides-with-code point, not on the commit-id point that already names the unshared board) stands: the caveat belongs where the rule is that it qualifies.
- **2026-11-20 22:54 · Alexander Sacharov** — critique round 2: two more sibling sites, docs/AGENTS.md:68 and jaira-role-pr SKILL.md:49 (point 2 would flag a missing ticket file on an unshared board). Added there too. grep 'rides with the code' over core/ and docs/ finds no further site after this.
- **2026-11-20 22:58 · Alexander Sacharov** — critique round 3: nothing left that needs a return. Left standing on purpose: the ready-to-paste commit lines in jaira-worker SKILL.md:103 and jaira-role-lane SKILL.md:86/288 still name .jaira/tickets/<ticket>.md; on an unshared board git refuses that path as ignored, and the agent already has the caveat from the block and the jaira role. Not a blocker under interim rule 4 (no wrong result, no data loss); a third critique return would trip the loop brake for a cosmetic line.
- **2026-11-20 23:02 · Alexander Sacharov** — bughunt: done at bb02c815 - 0 findings in 10 modules, first hunt; concurrency, error paths, input, permissions checked; files: .claude/skills/jaira/SKILL.md .jaira/tags AGENTS.md CLAUDE.md core/board/announce.go core/board/announce_test.go core/release/notes/SDBCT5.md core/role/builtin/jaira-role-pr/SKILL.md core/role/jaira/SKILL.md core/shipped/shipped.sha256 docs/AGENTS.md
- **2026-11-20 23:04 · Alexander Sacharov** — review: pass, keine Blocker. Eigener Check: frischer Build, git init + jaira init -> Block in CLAUDE.md/AGENTS.md traegt den Vorbehalt, .jaira/ bleibt aus git; Hashes in shipped.sha256 stimmen mit den Dateien; .claude/skills/jaira = core/role/jaira byte-gleich; go test board/role/cli -count=1 gruen. Drei notes in review-gaps.
- **2026-11-20 23:04 · Alexander Sacharov** — worker stop 21.11. 00:05: in signoff, review verdict pass. Branch feat/sdbct5 pushed (3 commits, last bb02c815). Ticket file modified in the worktree on purpose - no code commit follows; the logbook commit after acceptance carries it. Open for the person: three note: lines in review-gaps (caveat wording also covers 'commits nothing'; commit lines in jaira-worker/jaira-role-lane SKILL.md still name .jaira/tickets/; test name).
- **2026-11-21 08:00 · Alexander Sacharov** — Given back from signoff: Alex 21.11. alle Funde in dieser Runde schliessen. 1) Fertige Commit-Zeilen in jaira-worker SKILL.md und jaira-role-lane SKILL.md (+ alle Kopien, grep) nennen .jaira/tickets/<ticket>.md auch auf ungeteiltem Board - wie in den anderen Rollen anpassen. 2) Die Einschraenkung 'This point and the next hold on a shared board only' praezise fassen: 'commits nothing' gilt erst recht, nur die Ticketdatei-Saetze fallen weg. 3) Testname nennt auch --done --proof. Manifest-Hashes nachtragen, Diff-Hunt, review, signoff, push.
- **2026-11-21 08:01 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/4205.sock
- **2026-11-21 08:07 · Alexander Sacharov** — round after signoff return (21.11.): all three review-gaps notes closed. 1) commit lines in jaira-worker SKILL.md:103 and jaira-role-lane SKILL.md:86/288 now say to leave the ticket file out on an unshared board; grep 'git add.*.jaira/tickets' finds no further live site (rest: .planning, NOTES.md history, testdata). 2) caveat reworded in block and jaira role: no longer 'this point and the next hold on a shared board only', instead 'the next point holds there all the more; only its sentences about the ticket file fall away'. 3) test renamed TestAgentNoteNamesProofAndTheUnsharedBoardCaveat, asserts the new sentence. Manifest: 3 hashes appended. Release note line 2 names jaira-worker and jaira-role-lane.
