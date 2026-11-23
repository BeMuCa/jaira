---
id: 01M67K6X23DJ7TDVQ6DVDVPAPE
title: "Jede Person hat pro Projekt einen festen jaira-Namen, einmal abgefragt und lokal gespeichert; er ersetzt Git-Name und Aliasse"
status: critique
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
updated-at: 2026-11-23T20:50:36Z
updated-by: Alexander Sacharov
outcome-what: "jaira web schreibt unter dem aktuellen jaira-Namen (actAsNow vor jeder Schreibaktion) und die Seite liest nach jedem Fetch neu"
outcome-why: "Rueckgabe aus signoff: Namen im Web sollen reaktiv sein"
outcome-resolves: "Limit 'jaira web handelt unter dem Namen vom Start' fuer web aufgehoben; TUI bleibt wie angenommen"
review-summary: "internal/cli/web_source.go actAsNow sets the global peopleCache=nil under w.mu, but Server.pushNow -> webPush -> refs.Flush -> winner -> refs.Name (shownName -> peopleList) reads and lazily writes it outside w.mu: data race. Instead: in jaira web give refs.Name its own resolver that reads coreidentity.People(root) per call, as cards() already does, and leave peopleCache alone"
round: 7
back-from: critique
back-by: lane
back-count: 1
review-gaps: |-
  note: elsewhere: internal/cli/sync.go:182 stores a Claude Code task owner as assignee without Canonical - an agent task owner, not a typed person name; new ticket if it should resolve too
  note: limit: open TUI/jaira web act under the name read at open, note authors keep their written name (accepted at signoff)
test-verdict: "pass: go vet (linux+windows), go test ./... green at 4c2ddaf8b; TestAssigneeIsStoredUnderTheJairaName covers create --assignee GH and set assignee=gh -> file stores 'grace'"
hunt-verdict: "pass: stop-after-fix, no hunt in round 5"
claimed-by: DESKTOP-RFTCH11-64758
claimed-at: 2026-11-23T20:42:04Z
review-verdict: "accept: give-back edge closed in tickets.go create/set; test green; no blocker in the delta"
review-round: 5
review-verdict-gaps: |-
  note: elsewhere: internal/cli/sync.go:182 stores a Claude Code task owner as assignee without Canonical - an agent task owner, not a typed person name; new ticket if it should resolve too
  note: limit: open TUI/jaira web act under the name read at open, note authors keep their written name (accepted at signoff)
review-check: "1. cd into a jaira checkout with a remote  2. jaira people alias <you> TestAlias  3. jaira create 'probe' --assignee TestAlias  4. jaira show <new id> --json | grep assignee -> your jaira name, not TestAlias  5. jaira set <id> assignee=testalias, then git show refs/jaira/tickets/<ULID>:<ULID>.md | grep assignee -> your jaira name  6. or: go test ./internal/cli -run TestAssigneeIsStoredUnderTheJairaName"
stop-after-fix: ""
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
  proof: core/identity/identity.go currentIn; internal/cli/root.go s.Actor=identity(); internal/cli/tickets.go create/set Canonical(assignee); TestCurrentIsTheJairaNameThePeopleListGives, TestAssigneeIsStoredUnderTheJairaName
- [x] Befehl, um einer Person einen Alias zuzuordnen (z.B. jaira people alias <name> <git-name|email>) und die Liste zu zeigen (jaira people); Test
  proof: internal/cli/people.go newPeopleAliasCmd/newPeopleCmd; TestPeopleRefusals, TestAlias
- [x] Befehl, um eine Person umzubenennen (z.B. jaira people rename <alt> <neu>): der alte Name wird Alias des neuen; Board-Filter, list --assignee und das Cockpit zeigen alte Tickets unter dem neuen Namen, ohne dass Ticketdateien umgeschrieben werden muessen; Test
  proof: internal/cli/people.go newPeopleRenameCmd; shownName in tickets.go, peopleList().Same in list/next; tui boardPeople; cockpit mine via whoami aliases; TestPeopleTravelBetweenClonesAndRenameOldTickets
- [x] jaira snapshot sichert den people-Ref mit (siehe ES625B: Backup aller refs/jaira/*)
  proof: core/snapshot/snapshot.go PeopleFile; TestASnapshotHoldsThePeopleList
- [x] jaira web reagiert auf Namensaenderungen: Schreiben unter dem aktuellen jaira-Namen, Seite liest nach jedem Fetch neu; Test
  proof: internal/cli/web_source.go actAsNow; core/web/ui/src/pages/review/ui/ReviewPage.svelte heartbeat; TestWebWritesUnderTheNameThePeopleListGivesNow

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
- **2026-11-23 20:15 · Alexander Sacharov** — waiting: bug hunt round 3 (background subagent) since 21:15
- **2026-11-23 20:16 · Alexander Sacharov** — hunt report: 0 findings, nothing written to the board but this note and the closing line | read whole: core/people/people.go, core/people/store.go, core/identity/identity.go, internal/cli/people.go, core/snapshot/snapshot.go; core/gitref/gitref.go and internal/cli/flow.go read in the DVPAPE hunks plus the code on their paths (Usable, Fetch, refRead/refWrite/publish, Tree/Nest, LsTree; next --assignee), not every line; read as diff since c7c77814d: core/refsync/refsync.go, internal/cli/web_source.go, internal/tui/model.go, core/identity/identity.go | ran against a bare remote: 8 parallel 'people alias' from one clone (4 arrive, 4 exit 1 with ErrBusy, local ref = remote ref, nothing lost), exit codes for unknown rename (5), bad names (2), no-op alias --json (person null, 0), rename chain round-trip; not checked: a second clone writing, because the scratch board was unshared and the clone had no .jaira
- **2026-11-23 20:16 · Alexander Sacharov** — bughunt: done at 37a51cd47 - 0 findings in 13 modules, follow-up hunt; concurrency, error paths, input, permissions checked; files: README.md core/gate/gate.go core/gate/gate_test.go core/gitref/gitref.go core/identity/identity.go core/identity/identity_test.go core/people/people.go core/people/people_test.go core/people/store.go core/refsync/refsync.go core/release/notes/DVPAPE.md core/snapshot/snapshot.go core/snapshot/snapshot_test.go docs/COMMANDS.md internal/cli/flow.go internal/cli/people.go internal/cli/people_test.go internal/cli/refs.go internal/cli/root.go internal/cli/tickets.go internal/cli/web_source.go internal/tui/lanefocus.go internal/tui/model.go internal/tui/refs.go internal/tui/refs_test.go internal/tui/users.go internal/tui/view.go
- **2026-11-23 20:17 · Alexander Sacharov** — waiting: fresh review lane (run-lane.sh herdr tab) since 21:17
- **2026-11-23 20:20 · Alexander Sacharov** — review: done at 37a51cd47
- **2026-11-23 20:21 · Alexander Sacharov** — report: refs/jaira/people + jaira people [alias|rename]; identity resolves the actor via the list (JAIRA_USER first); list/next/show/JSON/web/TUI/refusals show jaira names, ticket files untouched; snapshot carries 'people' | signoff, context 253k of 1000k, cache warm 60 min after last answer | 5 review notes (set --assignee stores alias literally; open TUI/web act under the name from open; note authors stay; email-in-aliases under JAIRA_USER pre-existing; Departed drops List error pre-existing); branch not pushed; acceptance by person via review-check steps
- **2026-11-23 20:30 · Alexander Sacharov** — Given back from signoff: beheben:
- note: edge: 'jaira set <id> assignee=BeMuCa' stores the alias literally (checked: file says BeMuCa, list shows berk) - display is right, but the file does not carry the jaira name; fix: run the assignee value through peopleList().Canonical in set/create --assignee
angenommen:
- note: limit: an open TUI and a running 'jaira web' keep acting under the name read at open; only the displayed list refreshes after a fetch (model.go refreshPeople note)
- note: limit: progress-note authors and executed_by keep the name they were written under; the release note says so
- note: elsewhere: core/identity Aliases adds git user.email even when JAIRA_USER names someone else, so IsMe under JAIRA_USER=grace counts the checkout owner's email-assigned tickets as grace's; same on master, before this change
- note: elsewhere: core/refsync/refsync.go Departed drops the Repo.List error and reads it as no refs (found by the bug hunt, before this change)
- **2026-11-23 20:30 · Alexander Sacharov** — cockpit: given back to a new worker - nothing recorded about the old worker session (cache and context unknown); branch feat/dvpape is the ticket's worktree's, worktree /home/alex/projects/.worktrees/jaira-dvpape; reports to uds:/run/user/1000/cc-socks/48662.sock, the ticket's parent session
- **2026-11-23 20:31 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/48662.sock
- **2026-11-23 20:34 · Alexander Sacharov** — in-progress (give-back signoff): create --assignee und set assignee= speichern den jaira-Namen (peopleList().Canonical); Test TestAssigneeIsStoredUnderTheJairaName. stop-after-fix=true gesetzt: Fix < 20 Zeilen, kein Lock, kein Parser.
- **2026-11-23 20:35 · Alexander Sacharov** — waiting: fresh review lane (run-lane.sh herdr tab) since 21:35
- **2026-11-23 20:37 · Alexander Sacharov** — review: done at 4c2ddaf8b
- **2026-11-23 20:37 · Alexander Sacharov** — report: create --assignee / set assignee= store the jaira name (4c2ddaf8b) | signoff, context unknown | 2 review notes (sync.go task owner not canonicalised - elsewhere; TUI/web act under name from open - limit); branch not pushed
- **2026-11-23 20:41 · Alexander Sacharov** — Given back from signoff: beheben:
- note: limit: open TUI/jaira web act under the name read at open, note authors keep their written name (accepted at signoff)
dazu: надо чтобы web обновлял имена реактивно. такое возможно?
angenommen:
- note: elsewhere: internal/cli/sync.go:182 stores a Claude Code task owner as assignee without Canonical - an agent task owner, not a typed person name; new ticket if it should resolve too
- **2026-11-23 20:41 · Alexander Sacharov** — cockpit: given back to a new worker - nothing recorded about the old worker session (cache and context unknown); branch feat/dvpape is the ticket's worktree's, worktree /home/alex/projects/.worktrees/jaira-dvpape; reports to uds:/run/user/1000/cc-socks/48662.sock, the ticket's parent session
- **2026-11-23 20:42 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/48662.sock
- **2026-11-23 20:43 · Alexander Sacharov** — in-progress (give-back 2, web reaktiv): Anzeige in jaira web las die Liste schon pro Abruf; offen war (a) Schreiben unter dem Namen vom Start und (b) Seite liest erst im 60s-Takt neu. Plan: webSource loest vor jedem Schreiben neu auf, wer handelt (s.Actor, refs.Actor, peopleCache), hinter einem Mutex; Seite laedt neu, sobald /api/sync einen neuen last_fetch meldet (Heartbeat 10s). TUI bleibt wie angenommen.
- **2026-11-23 20:47 · Alexander Sacharov** — in-progress: stop-after-fix=true gesetzt - Fix ~50 Zeilen (Go 18 + Svelte 6 + Test), kein neuer Lock (vorhandener webSource-Mutex), kein Parser. Gates: go vet, go test ./..., -race TestWeb*, npm build/test, shipped-history gruen.
- **2026-11-23 20:48 · Alexander Sacharov** — critique: actAsNow -> peopleCache-Reset racet mit refs.Flush (laeuft in pushNow ausserhalb w.mu, ruft refs.Name=shownName). Fix: web setzt refs.Name auf People(root) pro Aufruf wie cards(); actAsNow fasst peopleCache nicht an. Ausserdem stop-after-fix wieder geloescht: Diff 56 Zeilen + dist, ueber der 50er-Grenze.
