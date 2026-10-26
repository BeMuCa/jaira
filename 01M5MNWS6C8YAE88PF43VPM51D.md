---
id: 01M5MNWS6C8YAE88PF43VPM51D
title: "Ein Checkout, dessen Ticket anderswo abgelegt wurde, zeigt es als abgelegt statt als geaendert"
status: in-progress
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
updated-at: 2026-10-26T19:43:28Z
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

- [x] Failing-Test zuerst (internal/cli, neue Datei): Bare-Remote + Clone + git worktree add (Helfer wie cloneWithRemotes/runCLI); Ticket im Haupt-Checkout gepullt und committet, im Worktree move --force nach done + jaira logbook; im Haupt-Checkout zeigen list, list --json, list --actionable --json, show, show --json es als abgelegt mit Entfernen-Befehl, und die Datei ist noch da
- [x] Zweiter Fall im selben Test (D0SAHM): Logbuch-Commit des Worktrees landet auf origin/master, Worktree entfernt, Ref weg (snapshot --drop); Haupt-Checkout auf altem Commit zeigt es weiter als abgelegt, gefunden auf dem Landing-Branch
- [x] Gegenfaelle im Test: per jaira restore zurueckgeholte Kopie (gleiches updated-at) und eine nach dem Ablegen beschriebene Kopie (neueres updated-at) werden NICHT markiert; stale Kopie mindestens 1 s aelter als die abgelegte
- [x] core/ticket: Laufzeitfeld am Ticket (wo abgelegt: Worktree-Pfad oder Landing-Branch + Pfad) und optionales Interface an Store.Source; Store.List markiert Datei-Tickets (nicht ReadOnly, nicht hier abgelegt); Store-Methode fuer ein einzelnes Ticket (show)
- [x] core/gitref: Worktrees() per git worktree list --porcelain; abgelegte Dateien eines Branches per einem git ls-tree -r <rev> -- .jaira/logbook .jaira/archive (id -> Pfad + Blob); Header der Kandidaten-Blobs per einem cat-file --batch wie ReadMany
- [x] core/refsync: Syncer erfuellt das Interface - FiledAwayIDs der anderen Worktrees + Landing-Branches (settings.Landing, gesetzt in internal/cli/refs.go attachRefs und internal/tui/refs.go newSyncer); Regel: anderswo abgelegt UND updated-at der Datei hier strikt aelter; Landing-Ergebnis nach Rev-SHA cachen
- [x] Kostentest mit GIT_TRACE: ein zweites List ohne Aenderung startet hoechstens 2 zusaetzliche git-Prozesse und keinen pro Ticket
- [ ] Entfernen-Hinweis an einer Stelle fuer CLI und TUI: nicht getrackt -> rm <pfad>; auf diesem Branch getrackt -> git restore <pfad>, dann git merge <landing> (nur im Worktree abgelegt: git restore <pfad>, geht mit dem Merge des ablegenden Branches); nur beim Ausgeben berechnet, nie im Board-Tick
- [ ] gate.Actionable gibt false fuer so ein Ticket (jaira next, list --actionable)
- [ ] jaira list: Text - raus aus der Lane, eigener Block am Ende mit Ort + Befehl; JSON - eigenes Array filed_elsewhere, nicht in tickets und count
- [ ] jaira show: Zeile ueber dem Detail (Ort + Befehl); JSON-Felder filed_elsewhere und remove_stale
- [ ] Board: reload schiebt solche Tickets an den Logbuch-Fuss der Terminal-Lane (m.logged/m.loggedDay), Lane-Zaehler ohne sie; Flag 'filed elsewhere'; refuseLogged nennt den Entfernen-Befehl statt jaira restore; openDetail laedt die Datei voll; ein Render-Test, mit -run gelaufen
- [ ] Release-Notiz core/release/notes/VPM51D.md: je Aenderung eine '- '-Zeile als Anweisung (Anzeige in list/show/Board, JSON-Felder, next ueberspringt)
- [ ] gofmt -l . leer, go vet ./..., go build ./...; neue Tests der beruehrten Pakete mit -run gruen

## Progress
- **2026-10-26 19:21 · BeMuCa** — Planning option ticked by the driver: the in-progress lane is cheap tier and the ticket needs a design decision first (what the remote carries after jaira logbook, and where board/list/show detect it). Pre-process runs at strong tier.
- **2026-10-26 19:35 · BeMuCa** — pre-process, Befund + Entwurf:
- Was der Remote nach logbook/archive traegt: Der Ref bleibt stehen und bekommt die Endbytes der abgelegten Datei: RecordFiled ruft nur Record (core/refsync/refsync.go:166-182). Ein 'filed'-Marker steht nicht darin, anders als bei Meilensteinen (internal/cli/logbook.go:340ff).
- logbook/archive stempeln vorher die Commits per Mutate, und das setzt updated-at auf jetzt (core/ticket/store.go:1468, logbook.go:302, archive.go:85). Die abgelegte Kopie ist also juenger als jede alte Kopie anderswo.
- Der Ref geht erst, wenn das Ticket auf einem Landing-Branch abgelegt liegt: snapshot reap (core/snapshot/snapshot.go:239-258) mit gitref.Landed (core/gitref/gitref.go:765-790, es zaehlen nur logbook/archive-Pfade).
- Ein anderer Checkout sieht also einen von zwei Zustaenden: (S1) der Ref steht mit dem Endstand, die Ablage liegt nur im ablegenden Worktree/Branch. (S2) der Ref ist weg, die Ablage liegt auf dem Landing-Branch (origin/master). S2 ist der D0SAHM-Fall aus K6BGB0.
- Was es schon gibt: offBoardDuplicates (store.go:995-1041) gilt nur, wenn beide Kopien im selben Checkout liegen. Es warnt nur und raet bewusst nicht. Bleibt unveraendert, IDs, die hier abgelegt sind, sind keine Kandidaten. Departed (refsync.go:1597) heisst 'Ref gesehen und jetzt weg'. Es erscheint nur in fetch/validate (fetch.go:98, validate.go:84) und sagt nicht, wohin das Ticket ging.
- Gewaehlte Erkennung: (a) andere Worktrees per git worktree list, darin deren FiledAwayIDs auf der Platte, auch uncommittet. Deckt S1 fuer Worktrees ab, also den DoD-Testfall. (b) Landing-Branches (settings.Landing, core/settings/settings.go:481): ein ls-tree auf .jaira/logbook und .jaira/archive. Deckt S2 ab, mit derselben Definition wie Landed/reap/Stranded.
- Regel: markiert nur, wenn updated-at der Datei hier STRIKT aelter ist als das der abgelegten Kopie. Eine restore-Kopie hat dasselbe updated-at (Restore bewegt nur, store.go:708-772) und wird nicht markiert. Wer nach dem Ablegen hier schreibt, macht sie neuer, und sie wird auch nicht markiert.
- Verworfen: (C) ein 'filed'-Marker im Ticket/Ref wie bei Meilensteinen. Er wuerde auch fremde Clones vor dem Merge abdecken, ist aber eine Formataenderung (restore, merge, validate muessten ihn kennen). (D) alle Remote-Branches scannen: ein ls-tree pro Branch pro CLI-Aufruf, und 'sofort starten' waere dahin. (E) git log --all ueber die History: dann wuerden auch Tickets als abgelegt gelten, die spaeter zurueckgeholt wurden.
- **2026-10-26 19:35 · BeMuCa** — pre-process, Entfernen-Befehl + Luecken:
- jaira delete verworfen: es loescht auch den Ref (RecordDelete, core/ticket/store.go:419-436). In S1 nimmt das den Ref weg, bevor das Ticket gelandet ist, entgegen refsync.go:166-179. Ausserdem verweigert es Tickets, auf die andere zeigen (internal/cli/delete.go referencesTo), und ist fuer 'haette nie existieren sollen' gedacht.
- Gewaehlt: git-Befehle je Zustand, geprueft mit gitrepo.Tracks (core/gitrepo/git.go:49). Nicht getrackt: rm <pfad>. Getrackt: git restore <pfad>, dann git merge <landing>. Das ist der Handfix von D0SAHM (stash + pull master). jaira druckt schon jetzt git-Befehle aus, siehe internal/cli/share.go:78.
- git rm auf einer getrackten Datei verworfen: der naechste Merge des Landing-Branches gibt dann einen rename/delete-Konflikt.
- Luecke: ein anderer CLONE legt ab, und sein Branch ist noch nicht gemergt. Das sieht man hier nicht, ohne Marker im Ref. Der Ref zeigt den Endstand in done. Das passt zur Absicht von RecordFiled (bis zum Landen sichtbar bleiben).
- Luecke: eine stale Kopie aus derselben Sekunde wie das Ablegen wird verpasst (updated-at hat Sekundenaufloesung). Darum muss der Test 1 s Abstand halten oder updated-at in der Fixture setzen.
- K6BGB0 (PR #91, nicht auf master) aendert dieselben Stellen: Ticket-Laufzeitfeld, list-Flags, ticketJSON, Karten-Flags und Store.List 'newest'. Egal, wer zuerst merged: die Abgelegt-Pruefung laeuft auf der Datei, vor 'newest', und ein markiertes Ticket wird nicht durch den Ref-Stand ersetzt.
- Nicht im Plan, Kandidat fuer ein Folgeticket: in S2 (Ref weg) geht ein Schreiben auf eine markierte Kopie durch und legt den Ref neu an. Das Ticket lebt dann wieder auf. Ausserdem zaehlt der Blocker-Index (link.Build) eine markierte Kopie weiter als offen, und web (internal/cli/web_source.go) ist nicht angefasst.
