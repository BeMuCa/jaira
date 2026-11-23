---
id: 01M67K6X23DJ7TDVQ6DVDVPAPE
title: "Jede Person hat pro Projekt einen festen jaira-Namen, einmal abgefragt und lokal gespeichert; er ersetzt Git-Name und Aliasse"
status: optimize
ready: true
creator: berkc
assignee: Alexander Sacharov
goal: "Ein Mensch erscheint auf dem Board immer unter genau einem jaira-Namen, egal welcher git user.name oder welche E-Mail gerade gesetzt ist. Die Namen und ihre Aliasse liegen auf einem Ref, den alle teilen und der sich konfliktfrei mergen laesst; jaira liest den Namen von dort und schreibt immer denselben."
context: |-
  Wunsch von berkc am 31.10.2026 (Punkt 2).
  Problem heute: das Board fuehrt dieselben Menschen unter mehreren Namen - BeMuCa (121) und berkc (88) fuer berkc, 'Alexander Sacharov' (48), 'alex', 'AlSa' fuer seinen Kollegen.
  Heute kommt der Name aus git user.name (internal/cli/root.go:353 identity -> core/identity); Aliasse stehen global in ~/.jaira/identity, nicht pro Projekt.
  Projekteinstellungen gibt es nur als committete Dateien (.jaira/chat-model, .jaira/release, bald .jaira/active-boards aus KVMG0K); eine lokale, gitignorierte Projekteinstellung gibt es nicht - die braucht es dafuer (z.B. .jaira/local/settings.json, gitignoriert).
  Web-Nutzer ohne Claude/CLI (PMs) werden im Web nach ihrem Namen gefragt.
  Offen: bestehende Tickets mit alten Namen umschreiben (Migration) oder nur als Alias lesen?
definition-of-done: "Erster Start in einem Projekt fragt den jaira-Namen ab und speichert ihn in einer lokalen, gitignorierten Projekteinstellung; Test"
tags:
  - cli
blocked-by: []
related:
  - 01M5MNWJF240306ZVDA41K3RA8
commits: []
created-at: 2026-10-30T23:44:36Z
updated-at: 2026-11-23T20:04:29Z
updated-by: Alexander Sacharov
outcome-what: "Neuer Ref refs/jaira/people (core/people) mit jaira people [alias|rename]; core/identity loest den Handelnden darueber auf; list/next/show/TUI zeigen und filtern unter dem jaira-Namen; snapshot sichert die Liste."
outcome-why: "Dieselben Menschen standen unter mehreren Namen auf dem Board (BeMuCa/berkc, Alexander Sacharov/alex/AlSa); Filter und Eigentum zerfielen danach."
outcome-resolves: DVPAPE
review-summary: "none"
round: 3
back-from: testing
back-by: lane
back-count: 1
review-gaps: "internal/tui/model.go: both board-open paths set me/Actor/myAliases/people by hand - now one call to refreshPeople, the same the fetch uses. Left as in round 1: identity.People one 'git show'; people.same mirrors identity.Same (import direction)."
test-verdict: "pass: go vet + GOOS=windows vet green, race run of people/identity/snapshot/gitref/cli/tui RC=0, full suite RC=0 earlier this round; DoD 4-9 verified against their tests; scratch board with two clones + bare: alias arrives, whoami on the other clone acts as berkc by email; two parallel aliases from both clones both landed (p1, p2); remote gone -> exit 3, nothing written"
hunt-verdict: "change: 2 findings - TUI keeps the people list/me from board open; gate refusal and updated-by/executed/note authors print stored names"
---

# Jede Person hat pro Projekt einen festen jaira-Namen, einmal abgefragt und lokal gespeichert; er ersetzt Git-Name und Aliasse

## Definition of Done

- [-] Erster Start in einem Projekt fragt den jaira-Namen ab und speichert ihn in einer lokalen, gitignorierten Projekteinstellung; Test
- [-] creator, assignee, claimed-by und neue Refs nutzen diesen Namen; alte Namen werden als Alias derselben Person gelesen
- [-] jaira web fragt einen Nutzer ohne Namen im Browser danach
- [x] Release-Note-Zeile
  proof: core/release/notes/DVPAPE.md
- [x] Personenliste auf einem eigenen Ref (z.B. refs/jaira/people): je Person ein jaira-Name und ihre Aliasse (git user.name, git user.email). Zwei Seiten, die gleichzeitig schreiben, werden zusammengefuehrt statt ueberschrieben (Vereinigung der Aliasse je Person); Test mit zwei Klonen
  proof: core/people/store.go Update; TestPeopleTravelBetweenClonesAndRenameOldTickets (internal/cli/people_test.go), TestAWriteThatLosesARaceIsMadeAgainOnTopOfTheWinner
- [x] core/identity loest den git-Namen/die E-Mail ueber diese Liste auf: creator, assignee, claimed-by, updated_by, Notes und neue Refs tragen immer den jaira-Namen; ein unbekannter Name bleibt wie heute der git-Name; Test
  proof: core/identity/identity.go currentIn; internal/cli/root.go s.Actor=identity(); TestCurrentIsTheJairaNameThePeopleListGives
- [x] Befehl, um einer Person einen Alias zuzuordnen (z.B. jaira people alias <name> <git-name|email>) und die Liste zu zeigen (jaira people); Test
  proof: internal/cli/people.go newPeopleAliasCmd/newPeopleCmd; TestPeopleRefusals, TestAlias
- [x] Befehl, um eine Person umzubenennen (z.B. jaira people rename <alt> <neu>): der alte Name wird Alias des neuen; Board-Filter, list --assignee und das Cockpit zeigen alte Tickets unter dem neuen Namen, ohne dass Ticketdateien umgeschrieben werden muessen; Test
  proof: internal/cli/people.go newPeopleRenameCmd; shownName in tickets.go, peopleList().Same in list/next; tui boardPeople; cockpit mine via whoami aliases; TestPeopleTravelBetweenClonesAndRenameOldTickets
- [x] jaira snapshot sichert den people-Ref mit (siehe ES625B: Backup aller refs/jaira/*)
  proof: core/snapshot/snapshot.go PeopleFile; TestASnapshotHoldsThePeopleList

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] gitref: PeopleRef + ReadPeople/WritePeople
- [x] core/people: Parse/Format/Merge/Lookup + Save mit Race-Retry; Test zwei Klone
- [x] core/identity: Current/Aliases loesen ueber Personenliste auf; Canonical-Helfer
- [x] CLI: jaira people [alias|rename]
- [x] Anzeige/Filter: list, next, show/JSON, TUI
- [x] snapshot sichert people
- [x] Release-Note, Gates

## Progress
- **2026-11-23 19:41 · Alexander Sacharov** — Alexander Sacharov took this ticket over from berkc
- **2026-11-23 19:41 · Alexander Sacharov** — decision: Alex 23.11.2026 ersetzt den Ansatz von berkc (lokale, gitignorierte Projekteinstellung, Abfrage beim ersten Start): Namen gehoeren in einen geteilten Ref, damit alle dieselbe Zuordnung sehen und sie mergebar ist ('sami imena v refs, ot tuda berutsja i pishutsja vsegda odno i to zhe'). Dazu: Alias (jeder git-Name/jede E-Mail einer Person -> ein jaira-Name) und Umbenennen. DoD 1-3 deshalb ersetzt. Bestehende Tickets nicht umschreiben: alte Namen werden ueber Aliasse gelesen. Web-Abfrage (alter DoD 3) gestrichen. Bekannte Doppelnamen fuer den Start: BeMuCa = berkc; 'Alexander Sacharov' = alex = AlSa. Vom Teamlead auf Alex' Wunsch von berkc uebernommen.
- **2026-11-23 19:42 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/48662.sock
- **2026-11-23 19:44 · Alexander Sacharov** — decisions (worker, keine offenen fuer die Person): (1) ein Ref refs/jaira/people mit einer Datei 'people', je Zeile '<jaira-name>: alias, alias' - hand-editierbar, diff-lesbar. (2) Schreiben synchron wie refWrite: fetch, mergen, push mit Lease, bei verlorenem Race neu mergen; offline/ohne Remote wird abgelehnt (kein Outbox-Kind - kleiner; nichts geht still verloren, der Befehl scheitert sichtbar). (3) Merge = Vereinigung je Person (case-insensitiv) + Falten: eine Person, deren Name Alias einer anderen ist, geht in dieser auf - so ueberlebt rename ein gleichzeitiges alias auf den alten Namen ohne Basis-Diff. (4) Keine Selbstregistrierung: unbekannter git-Name bleibt git-Name (DoD 6). (5) ~/.jaira/identity bleibt als persoenliche Aliasliste. (6) Anzeige: list/next --assignee, Ticket-JSON/show (assignee, creator) und TUI-Filter/-Anzeige zeigen den jaira-Namen; Ticketdateien bleiben unveraendert. Cockpit 'mine' laeuft ueber whoami-Aliasse, die die Personenliste mit einschliessen. (7) snapshot legt die Datei 'people' neben board/ in den Snapshot-Branch.
- **2026-11-23 19:53 · Alexander Sacharov** — in-progress: Merge-Funktion bewusst weggelassen - Update wiederholt die Aenderung auf dem neuesten Stand statt Text gegen Text zu mergen; da Fetch die lokalen Refs erzwingt und nur Push-bestaetigte Staende lokal stehen, gibt es keinen Ort, an dem zwei Listen zu mergen waeren. 'jaira people' liest den lokalen Ref ohne Fetch, wie jede andere Leseoperation. TUI: boardPeople ist eine Paket-Variable, weil matches() eine freie Funktion ist, die ~40 Tests mit nil aufrufen.
- **2026-11-23 19:53 · Alexander Sacharov** — critique: internal/cli/people.go alias - unknown person aliased to its own name returns unchanged with no person; output derefs nil p. Fix: guard nil. Sonst: Form (ein Ref, synchroner Write statt Outbox) bleibt - siehe decisions-Note (2).
- **2026-11-23 19:56 · Alexander Sacharov** — waiting: bug hunt (background subagent) since 20:56
- **2026-11-23 19:58 · Alexander Sacharov** — bughunt: change: internal/tui/model.go:385,397,398 (same at 426,433,434) - class: identity read from the people ref once per board open and never again, while the TUI's own refresh (internal/tui/refs.go:103, refsync Incoming -> Repo.Fetch of refs/jaira/*) brings in a newer people ref every round; boardPeople, m.me and m.myAliases stay at the list of the moment the board opened; example: a teammate runs 'jaira people rename berkc berk' while my board is open - my cards, detail pane and user picker keep showing berkc, and moves I make from the TUI write assignee=berkc until I restart; nothing is lost (the old name stays an alias), so not a blocker; ran: no live TUI run, read the refresh path through to the fetch refspec (gitref.go:851)
- **2026-11-23 19:58 · Alexander Sacharov** — bughunt: change: core/gate/gate.go:264 - class: person names printed straight from the ticket file instead of through the people list, so output names someone the board no longer shows; also internal/cli/tickets.go printDetail rows 'updated by' and 'executed', and note authors; example: ticket stored with assignee=BeMuCa, person renamed to berk - 'jaira list' shows berk, but 'JAIRA_USER=someoneelse jaira move <id> --to todo' refuses with 'VDTK6D belongs to BeMuCa — ask them' (exit 3); the gate line is older, the divergence comes from this change; edge, not a blocker; ran: scratch repo with bare remote, people alias + rename, set assignee=BeMuCa, move as another user
- **2026-11-23 19:58 · Alexander Sacharov** — hunt report: 2 findings, both change, neither a blocker, as notes on DVPAPE | read whole: core/people/people.go core/people/store.go internal/cli/people.go; read as the change in context: core/gitref/gitref.go core/identity/identity.go core/snapshot/snapshot.go internal/cli/root.go internal/cli/tickets.go internal/cli/flow.go internal/tui/model.go internal/tui/view.go internal/tui/users.go internal/tui/lanefocus.go; ran end to end on a scratch repo: alias, rename, exit codes 2/3/5, list --assignee across old names, ownership after rename, 18 parallel alias writes (6 got the documented 'run it again', none half-written) | not checked: the TUI finding live (no terminal session), Windows, a remote that rejects the push for reasons other than a lease
- **2026-11-23 19:58 · Alexander Sacharov** — bughunt: done at 00e1138ea - 2 findings (2 change, 0 pre-existing) in 15 modules, first hunt; concurrency, error paths, input, permissions checked; files: README.md core/gitref/gitref.go core/identity/identity.go core/identity/identity_test.go core/people/people.go core/people/people_test.go core/people/store.go core/release/notes/DVPAPE.md core/snapshot/snapshot.go core/snapshot/snapshot_test.go docs/COMMANDS.md internal/cli/flow.go internal/cli/people.go internal/cli/people_test.go internal/cli/root.go internal/cli/tickets.go internal/tui/lanefocus.go internal/tui/model.go internal/tui/users.go internal/tui/view.go
- **2026-11-23 20:03 · Alexander Sacharov** — in-progress (bughunt round 1): TUI liest Personenliste, me, Actor und Aliasse nach jedem Hintergrund-Fetch neu (refreshPeople; boardPeople jetzt atomic, weil der Syncer im Fetch-Goroutine darueber Namen aufloest). Namen in Meldungen: gate.Env.People fuer 'belongs to', refsync.Syncer.Name fuer Winner/Arrival (pull/release/claim-Verweigerungen, Ankuendigungen), updated_by in show/JSON/TUI, jaira web liest die Liste pro Abruf. Bewusst NICHT: Autorennamen in Progress-Notizen - das ist Fliesstext im Body, Umschreiben hiesse Ticketdateien aendern; Release-Note sagt es. executed_by ist ein Agent, keine Person.
