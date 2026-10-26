---
id: 01M5MNWS6C8YAE88PF43VPM51D
title: "Ein Checkout, dessen Ticket anderswo abgelegt wurde, zeigt es als abgelegt statt als geaendert"
status: pre-process
ready: true
creator: BeMuCa
goal: "Nach jaira logbook in einem Worktree verwirrt die alte Ticketdatei in den anderen Worktrees und im Haupt-Checkout niemanden mehr: sie gilt als abgelegt, und jaira bietet an, sie zu entfernen."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkt 9): nach einem --force-Move und jaira logbook in einem Worktree zeigten die anderen Worktrees die Ticketdatei weiter als geaendert; der naechste Worker hielt das Ticket fuer offen.
  K6BGB0 hat den Fall 'Datei hinter Ref' geloest; dieser hier ist 'Datei ohne Ref, Ticket upstream im Logbuch' - K6BGB0s Notiz vom 14.10. nennt ihn als dritte Variante, die offen bleibt.
  Verwandt: DPSFGE (git status im Haupt-Checkout ohne Ticket-Rauschen), AYBDJ2 (Arbeitskopie nach logbook aufraeumen).
definition-of-done: "Board, jaira list und jaira show erkennen eine Ticketdatei, deren Ticket im Logbuch eines anderen Checkouts oder auf dem Remote abgelegt ist, und zeigen sie als 'abgelegt'; ein Test baut den Fall mit zwei Worktrees"
tags:
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:25:56Z
updated-at: 2026-10-26T19:35:18Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-374545
claimed-at: 2026-10-26T19:21:26Z
---

# Ein Checkout, dessen Ticket anderswo abgelegt wurde, zeigt es als abgelegt statt als geaendert

## Definition of Done

- [ ] Board, jaira list und jaira show erkennen eine Ticketdatei, deren Ticket im Logbuch eines anderen Checkouts oder auf dem Remote abgelegt ist, und zeigen sie als 'abgelegt'; ein Test baut den Fall mit zwei Worktrees
- [ ] jaira bietet an, die veraltete Kopie zu entfernen (Befehl in der Meldung); nichts wird still geloescht
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [ ] Failing-Test zuerst (internal/cli, neue Datei): Bare-Remote + Clone + git worktree add (Helfer wie cloneWithRemotes/runCLI); Ticket im Haupt-Checkout gepullt und committet, im Worktree move --force nach done + jaira logbook; im Haupt-Checkout zeigen list, list --json, list --actionable --json, show, show --json es als abgelegt mit Entfernen-Befehl, und die Datei ist noch da
- [ ] Zweiter Fall im selben Test (D0SAHM): Logbuch-Commit des Worktrees landet auf origin/master, Worktree entfernt, Ref weg (snapshot --drop); Haupt-Checkout auf altem Commit zeigt es weiter als abgelegt, gefunden auf dem Landing-Branch
- [ ] Gegenfaelle im Test: per jaira restore zurueckgeholte Kopie (gleiches updated-at) und eine nach dem Ablegen beschriebene Kopie (neueres updated-at) werden NICHT markiert; stale Kopie mindestens 1 s aelter als die abgelegte
- [ ] core/ticket: Laufzeitfeld am Ticket (wo abgelegt: Worktree-Pfad oder Landing-Branch + Pfad) und optionales Interface an Store.Source; Store.List markiert Datei-Tickets (nicht ReadOnly, nicht hier abgelegt); Store-Methode fuer ein einzelnes Ticket (show)
- [ ] core/gitref: Worktrees() per git worktree list --porcelain; abgelegte Dateien eines Branches per einem git ls-tree -r <rev> -- .jaira/logbook .jaira/archive (id -> Pfad + Blob); Header der Kandidaten-Blobs per einem cat-file --batch wie ReadMany
- [ ] core/refsync: Syncer erfuellt das Interface - FiledAwayIDs der anderen Worktrees + Landing-Branches (settings.Landing, gesetzt in internal/cli/refs.go attachRefs und internal/tui/refs.go newSyncer); Regel: anderswo abgelegt UND updated-at der Datei hier strikt aelter; Landing-Ergebnis nach Rev-SHA cachen
- [ ] Kostentest mit GIT_TRACE: ein zweites List ohne Aenderung startet hoechstens 2 zusaetzliche git-Prozesse und keinen pro Ticket
- [ ] Entfernen-Hinweis an einer Stelle fuer CLI und TUI: nicht getrackt -> rm <pfad>; auf diesem Branch getrackt -> git restore <pfad>, dann git merge <landing> (nur im Worktree abgelegt: git restore <pfad>, geht mit dem Merge des ablegenden Branches); nur beim Ausgeben berechnet, nie im Board-Tick
- [ ] gate.Actionable gibt false fuer so ein Ticket (jaira next, list --actionable)
- [ ] jaira list: Text - raus aus der Lane, eigener Block am Ende mit Ort + Befehl; JSON - eigenes Array filed_elsewhere, nicht in tickets und count
- [ ] jaira show: Zeile ueber dem Detail (Ort + Befehl); JSON-Felder filed_elsewhere und remove_stale
- [ ] Board: reload schiebt solche Tickets an den Logbuch-Fuss der Terminal-Lane (m.logged/m.loggedDay), Lane-Zaehler ohne sie; Flag 'filed elsewhere'; refuseLogged nennt den Entfernen-Befehl statt jaira restore; openDetail laedt die Datei voll; ein Render-Test, mit -run gelaufen
- [ ] Release-Notiz core/release/notes/VPM51D.md: je Aenderung eine '- '-Zeile als Anweisung (Anzeige in list/show/Board, JSON-Felder, next ueberspringt)
- [ ] gofmt -l . leer, go vet ./..., go build ./...; neue Tests der beruehrten Pakete mit -run gruen

## Progress
- **2026-10-26 19:21 · BeMuCa** — Planning option ticked by the driver: the in-progress lane is cheap tier and the ticket needs a design decision first (what the remote carries after jaira logbook, and where board/list/show detect it). Pre-process runs at strong tier.
