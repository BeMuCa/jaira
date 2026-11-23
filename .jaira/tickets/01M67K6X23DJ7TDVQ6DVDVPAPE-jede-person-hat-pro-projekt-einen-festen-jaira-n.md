---
id: 01M67K6X23DJ7TDVQ6DVDVPAPE
title: "Jede Person hat pro Projekt einen festen jaira-Namen, einmal abgefragt und lokal gespeichert; er ersetzt Git-Name und Aliasse"
status: in-progress
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
updated-at: 2026-11-23T20:13:11Z
updated-by: Alexander Sacharov
outcome-what: "Neuer Ref refs/jaira/people (core/people) mit jaira people [alias|rename]; core/identity loest den Handelnden darueber auf; list/next/show/TUI zeigen und filtern unter dem jaira-Namen; snapshot sichert die Liste."
outcome-why: "Dieselben Menschen standen unter mehreren Namen auf dem Board (BeMuCa/berkc, Alexander Sacharov/alex/AlSa); Filter und Eigentum zerfielen danach."
outcome-resolves: DVPAPE
review-summary: "none"
round: 4
back-from: testing
back-by: lane
back-count: 2
review-gaps: "internal/tui/model.go: both board-open paths set me/Actor/myAliases/people by hand - now one call to refreshPeople, the same the fetch uses. Left as in round 1: identity.People one 'git show'; people.same mirrors identity.Same (import direction)."
test-verdict: "pass (round 2): go vet + windows vet green, full suite RC=0, race over cli/tui/core RC=0 before the last fold; DoD 4-9 re-verified incl. new TestOwnershipRefusalNamesTheJairaName, TestTheBoardFetchBringsInThePeopleList, pull refusal in TestPeopleTravelBetweenClonesAndRenameOldTickets; scratch board round 1: parallel aliases both land, remote gone exit 3 nothing written"
hunt-verdict: "change: blocker JAIRA_USER loses to the email lookup (identity.go currentIn); edge: TUI syncer Actor not refreshed; edge: reject heldBy/TakenFrom/web Rejected.By print stored names; pre-existing: Departed drops Repo.List error (note only)"
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
- **2026-11-23 20:05 · Alexander Sacharov** — waiting: bug hunt round 2 (background subagent) since 21:05
- **2026-11-23 20:10 · Alexander Sacharov** — bughunt: change: blocker: core/identity/identity.go:27-35 (currentIn) - class: an explicit identity override (JAIRA_USER) loses to the people-list lookup by git user.email, so whoever git's email names acts, not the person the environment names; reaches s.Actor, Syncer.Actor (refsync Pull/Release/Reject write assignee, updated-by, rejected-by, notes) and gate Request.Actor (the not_owner rail); example: people ref has berkc = berk@example.test, git user.email = berk@example.test, 'JAIRA_USER=grace jaira whoami' says 'Acting as: berkc' with grace only under 'Also me' - so 'jaira pull' under JAIRA_USER=grace writes assignee=berkc and a move on berkc's ticket is not refused; on master JAIRA_USER always won; blame: 996c3e735 + 00e1138ea (DVPAPE); ran: branch binary in a scratch repo with a bare remote, 'jaira people alias berkc berk@example.test' then whoami with and without JAIRA_USER (both print Acting as: berkc)
- **2026-11-23 20:10 · Alexander Sacharov** — bughunt: change: internal/tui/model.go:461-466 (refreshPeople) - class: who the board acts as is refreshed in some holders and not in others: m.me, store.Actor, boardPeople and myAliases follow a fetch, the ref syncer built by newSyncer(s, m.me) at model.go:391/426 keeps Actor (and the ref commit author) from board open; same pattern in 'jaira web', where s.Actor and the CLI peopleCache (internal/cli/root.go:411) are read once for the server's whole life; example: board opened as 'alex', a teammate's 'jaira people alias "Alexander Sacharov" alex' arrives with a fetch, R on a ref-only card (rejectOnRef -> y.Reject, model.go:2600) still writes rejected-by/updated-by and the outbox 'by' as 'alex' while store writes say 'Alexander Sacharov'; not a blocker: the old name stays an alias of the same person after alias/rename, so nothing changes owner; read only, live TUI not run
- **2026-11-23 20:10 · Alexander Sacharov** — bughunt: change: core/refsync/refsync.go:1798 - class: a person named in output that does not pass through the people list, while its neighbours (Winner, Incoming, gate not_owner, web cards) now do; places: Reject's heldBy -> refsync.HeldByWarning (refsync.go:1821) printed by internal/cli/reject.go:204 and internal/tui/model.go:2617; Pulled.TakenFrom and the 'took this ticket over from %s' note text (refsync.go:1613-1618); web TicketDetail Rejected.By (internal/cli/web_source.go rejection()); example: ticket assigned to BeMuCa where the list says berkc, 'jaira reject' warns 'was assigned to BeMuCa' while 'jaira list' shows berkc; the lines are older than the ticket, marked change because the gap exists only since this change introduced jaira names (unclear, so change); edge, display only; read
- **2026-11-23 20:10 · Alexander Sacharov** — bughunt: pre-existing: core/refsync/refsync.go:1993-1998 (Departed) - class: an error from Repo.List() is dropped and read as 'no refs at all'; example: a transient git failure during for-each-ref in 'jaira fetch' or 'jaira validate' reports every seen ticket with a local file as departed; report only, no file is touched; blame: 2f5713f18 (before the ticket); read
- **2026-11-23 20:10 · Alexander Sacharov** — hunt report: 4 findings (3 change, 1 blocker among them; 1 pre-existing), all as notes on DVPAPE | read whole: core/gate/gate.go core/refsync/refsync.go internal/cli/refs.go internal/cli/web_source.go internal/tui/refs.go (new to the hunt, gate/refsync via a helper agent); read as a diff since 00e1138ea: core/release/notes/DVPAPE.md internal/cli/root.go internal/cli/tickets.go internal/tui/lanefocus.go internal/tui/model.go internal/tui/users.go internal/tui/view.go (identity.go re-entered only to follow the JAIRA_USER path) | not checked: live TUI for the stale syncer Actor; a -race run of refreshPeople writing store.Actor against tea.Cmd goroutines; Windows; milestone merge paths in refsync only read
- **2026-11-23 20:10 · Alexander Sacharov** — bughunt: done at c7c77814d - 4 findings (3 change, 1 pre-existing) in 12 modules, follow-up hunt; concurrency, error paths, input, permissions checked; files: core/gate/gate.go core/gate/gate_test.go core/refsync/refsync.go core/release/notes/DVPAPE.md internal/cli/people_test.go internal/cli/refs.go internal/cli/root.go internal/cli/tickets.go internal/cli/web_source.go internal/tui/lanefocus.go internal/tui/model.go internal/tui/refs.go internal/tui/refs_test.go internal/tui/users.go internal/tui/view.go
- **2026-11-23 20:13 · Alexander Sacharov** — in-progress (bughunt round 2): blocker fixed - JAIRA_USER, wenn gesetzt, ist der einzige nachgeschlagene Name; git user.email entscheidet nur ohne ihn (identity.go currentIn, Testfall grace/someone). Kleineres Design statt weiterer Halter: Die offene TUI liest nach einem Fetch nur die Personenliste neu (Anzeige/Filter); als wen sie handelt (me, store.Actor, Syncer.Actor) bleibt vom Oeffnen - einheitlich in allen Haltern, und der Name gehoert weiterhin derselben Person. Gleiches gilt fuer jaira web. Anzeige: reject heldBy, already-rejected-by, Steal-Notiz/TakenFrom und web Rejected.By laufen ueber die Liste. pre-existing (nicht dieses Ticket): refsync Departed verschluckt Repo.List-Fehler - Notiz fuer die Person.
