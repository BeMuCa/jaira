---
id: 01M2G02EHTV6RJPPQKR5VM0A76
title: "Ein Board im Datei-Modus sagt es nicht, kommt nicht zurueck und laesst sich nicht pruefen"
status: in-progress
ready: true
creator: Alexander Sacharov
goal: "Wer ein Ticket anlegt, sieht in derselben Zeile, ob es auf einem Ref liegt oder als Datei; ein Datei-Ticket kommt mit einem Befehl auf seinen Ref; und ein Befehl sagt, in welchem Modus dieses Board laeuft und warum."
context: |-
  Alex hat am 2026-09-14 eine Stunde verloren, weil siebzehn Tickets auf dem requirementsgenie-Board still als Datei angelegt wurden statt auf ihren Refs. Kein Befehl hat das gesagt, und zurueck ging es auch nicht.

  Die Ursache selbst - der Remote-Name gilt pro Rechner statt pro Board - ist 9ET6NC und wird dort behoben. Dieses Ticket ist alles andere: dass derselbe Fehler wieder eine Stunde kosten wuerde, auch wenn er aus einem anderen Grund auftritt.

  Fuenf Stellen, an denen das Werkzeug es haette sagen koennen:

  1. 'jaira create' nennt den gewaehlten Modus nicht. internal/cli/refs.go:69 fileOnRefOnly() gibt bei fehlendem Remote einfach false zurueck, und create schweigt darueber. Im Ref-Modus steht 'On its ref, not on your disk' da; im Datei-Modus steht nichts.

  2. Der Diagnosetext existiert, steht aber an der falschen Stelle. Nur 'jaira release' nennt die wirkliche Ursache: 'this board does not carry tickets on refs: gitref: no repository or no such remote: no remote "upstream"'. release laeuft am Ende, der Fehler entsteht bei create am Anfang.

  3. Es gibt keinen Weg zurueck. Ein Datei-Ticket laesst sich nicht auf seinen Ref legen. 'jaira release' darauf antwortet 'gitref: no ref for this ticket' - eine Feststellung, kein Hinweis. Der einzige Weg war, siebzehn Tickets mit neuen Ids neu anzulegen.

  Das Stueck dafuer ist schon da und ist universell: Record() (core/refsync/refsync.go:108) stellt jede Schreibung in die Outbox, und ein Ticket ohne Ref least den leeren String, also 'ich erwarte, dass dieser Ref nicht existiert'. fileOnRefOnly() (internal/cli/refs.go:69) schickt und entfernt danach die Datei. An create gebunden ist davon nur der Funktionsname.

  4. Kein Befehl zeigt den Zustand des Boards. 'jaira whoami' zeigt die Identitaet, aber nichts von der git-Seite: welcher Remote erwartet wird, ob er da ist, in welchem Modus das Board laeuft, wie viele Tickets nur als Datei liegen.

  5. 'jaira create --dod' nimmt genau einen Eintrag. Mehrzeilig uebergeben wird alles ab dem zweiten Absatz zu loser Prosa im Rumpf, ohne Kaestchen. Weitere Punkte gehen nur ueber 'jaira dod --add', und das verlangt die Datei - also 'jaira pull', und das schreibt den Aufrufer als assignee ein. Ein Ticket mit sieben Kriterien laesst sich im Ref-Modus also nicht anlegen, ohne es sich selbst zuzuweisen.

  Nachgestellt am 2026-09-14: APABM4, S1VM40 und 9ET6NC wurden mit fuenf bis sechs Kriterien angelegt und trugen danach je genau eines. Bei 9ET6NC lief der Gate der Endlane dadurch gegen ein Fuenftel der Bedingungen, waehrend ein Worker den Rest als Kontext las.
definition-of-done: "'jaira create' nennt den Modus in beiden Faellen. Im Datei-Modus sagt es, dass das Ticket als Datei und nicht auf einem Ref liegt, nennt den Remote-Namen, nach dem gesucht wurde, und den Grund - nicht nur das Schweigen von heute."
tags:
  - cli
blocked-by: []
related:
  - 01M2FYEHZYFCW7DSNRHY9ET6NC
commits: []
created-at: 2026-09-14T13:00:30Z
updated-at: 2026-09-14T13:38:27Z
assignee: "Alexander Sacharov"
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-155812
claimed-at: 2026-09-14T13:29:29Z
---

# Ein Board im Datei-Modus sagt es nicht, kommt nicht zurueck und laesst sich nicht pruefen

## Definition of Done

- [ ] 'jaira create' nennt den Modus in beiden Faellen. Im Datei-Modus sagt es, dass das Ticket als Datei und nicht auf einem Ref liegt, nennt den Remote-Namen, nach dem gesucht wurde, und den Grund - nicht nur das Schweigen von heute.
- [ ] Ein Datei-Ticket kommt mit einem Befehl auf seinen Ref und die lokale Datei verschwindet dabei. Die Logik dafuer ist die vorhandene fileOnRefOnly (internal/cli/refs.go:69), nicht eine zweite Kopie davon. Nachgestellt auf einem Fixture-Board, dessen Ticket im Datei-Modus entstanden ist.
- [ ] Ein Befehl zeigt den git-Zustand des Boards in einem Aufruf: den eingestellten Remote-Namen, die Remotes die dieses Repository hat, ob der Ref-Modus laeuft, und wie viele Tickets nur als Datei liegen. Nachgestellt auf einem Board mit passendem und auf einem mit fehlendem Remote.
- [ ] 'jaira create --dod' nimmt den Schalter mehrfach, wie --tag es tut. Nachgestellt: ein Ticket mit drei Kriterien wird im Ref-Modus mit einem Aufruf angelegt und traegt danach drei Kaestchen, ohne dass es dafuer gepullt wurde.
- [ ] Der Diagnosetext, der heute nur aus 'jaira release' kommt, erscheint dort wo der Zustand entsteht. Nachgestellt: auf einem Board ohne passenden Remote nennt schon der erste 'jaira create' den Grund, nicht erst ein Befehl am Ende der Kette.
- [ ] Je eine Zeile in core/release/NOTES.md unter ## Unreleased fuer jede von aussen sichtbare Aenderung: die Modus-Zeile in create, der neue Befehl fuer den Ref-Nachtrag, der neue Zustandsbefehl, und der wiederholbare --dod.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] read how create picks the mode: fileOnRefOnly (internal/cli/refs.go:69) and its caller internal/cli/tickets.go:293-311, plus gitref Repo.Usable/noRemote — the diagnostic text already exists since 9ET6NC, only a caller at create is missing
- [x] make fileOnRefOnly answer with the reason, not a bare bool: return the mode and the refs.Usable() error, so create can say why it chose the file
- [ ] failing test in internal/cli: 'create' on a board whose remote is absent prints a line naming the remote looked for and the reason; --json carries the same as a field
- [x] implement that create line and the json field (DoD 1 and DoD 5 are the same change)
- [x] decide the way back: extend 'jaira release' to a ticket with no ref, or add a new command — write the decision and its reason as a jaira note
- [x] generalise fileOnRefOnly so it also takes an existing ticket: Record() the current bytes with the empty lease, flush, then drop the file — one function, not a second copy
- [ ] failing test: a fixture board whose ticket was created in file mode gets a remote; one command puts it on its ref and the local file is gone
- [x] implement the way back
- [x] design the state command (DoD 3): extend 'jaira whoami' with the git side, or add a separate command — and fix the four facts it prints: remote name and where it came from, remotes this repository has, ref mode yes/no with the reason, how many tickets lie as files only
- [ ] failing test for the state command on a board with a matching remote and on one without
- [x] implement the state command, text and --json
- [x] make --dod repeatable (DoD 4): StringArray like --tag, ticket.NewBody takes several items, the frontmatter definition-of-done keeps the first and the body checklist carries all of them
- [ ] failing test: create with three --dod in ref mode leaves three boxes without anybody pulling the ticket
- [x] implement the repeatable --dod, and check gate.Ready/missingFields still read it
- [ ] fix the misleading comment on TestTheRefGoesToTheConfiguredRemoteWhenTheRepositoryHasIt (internal/cli/boardremote_test.go): both remotes point at the same bare repository, so the comment promises a check the test does not make
- [ ] one line per externally visible change under ## Unreleased in core/release/NOTES.md (DoD 6)
- [ ] go test ./... and tick each DoD box with its proof

## Progress
- **2026-09-14 13:27 · Alexander Sacharov** — Plan-Lane, Begruendung. DoD 1 und DoD 5 sind eine einzige Aenderung: 9ET6NC hat den Diagnosetext schon gebaut (gitref Repo.noRemote nennt gesuchten Remote, vorhandene Remotes und 'git config jaira.remote'). Er fehlt nur bei create, weil fileOnRefOnly (internal/cli/refs.go:75) bei refs.Usable() != nil ein nacktes false zurueckgibt und den Grund wegwirft. Also: Rueckgabe um den Grund erweitern, nicht einen zweiten Text schreiben.
- **2026-09-14 13:27 · Alexander Sacharov** — Plan-Lane, offene Entwurfsfragen, die in Schritt 5 und 9 entschieden werden. (a) Weg zurueck: 'jaira release' erweitern statt neuem Befehl - release heisst schon 'zurueck aufs Board', entfernt schon die Datei und laeuft heute nur deshalb ins Leere, weil es refs.Release() ohne Ref fragt. Gegen release spricht, dass es zusaetzlich den assignee loescht; fuer ein ungearbeitetes Datei-Ticket ist genau das richtig, fuer ein gerade bearbeitetes nicht. (b) Zustandsbefehl: 'jaira whoami' erweitern statt neuem Befehl - whoami ist schon 'was denkt jaira ueber diese Umgebung', und die Projektregel misst jedes Feature an 'kleiner als paca'. Beides bewusst als Option im Plan gelassen, weil die Kritik-Lane das umdrehen darf.
- **2026-09-14 13:27 · Alexander Sacharov** — Plan-Lane, zu DoD 4: das Frontmatter-Feld definition-of-done ist ein einzelner String, die Kaestchen im Rumpf sind die Wahrheit fuer den Gate (core/gate/gate.go:504-510 liest DoDItems, wenn es welche gibt, sonst DoD). Beide duerfen also auseinanderlaufen - dieses Ticket hier tut es bereits. Darum: --dod mehrfach, erstes Vorkommen ins Frontmatter, alle in ticket.NewBody (core/ticket/body.go:25, Signatur von string auf []string). Kein neues Feld, kein Listen-Frontmatter.
- **2026-09-14 13:32 · Alexander Sacharov** — In-progress lane, decision on plan step 5 (the way back). 'jaira release' is extended; no new command. Reasons: release already means 'back onto the board', already removes the local file, and already clears the assignee — which is exactly the right end state for the case this ticket came from, seventeen unworked tickets sitting as files. The only thing missing was that refs.Release() asks the ref first and a file ticket has none, so it died on gitref.ErrNoRef. Rejected: a separate 'jaira publish/push'. It would need its own name, help, tests and NOTES line for a path a board takes once, and it would leave two commands whose difference is invisible from the outside ('the file goes away and the ticket is on its ref' in both). Cost of the choice, so nobody 'fixes' it back: release on a file ticket that somebody is actually working clears their assignee too. That is release's meaning, not a bug; 'jaira pull' takes it straight back.
- **2026-09-14 13:32 · Alexander Sacharov** — In-progress lane, decision on plan step 9 (the state command). 'jaira whoami' is extended with a Board block; no new command. Reasons: whoami is already 'what does jaira think about this environment before I start', and the four git facts are the other half of that same question — identity decides whose tickets these are, the remote decides where they go. A second command ('jaira status') would be a second place to look for one answer, would read as a git-status lookalike it is not, and the project rule measures every feature against 'is this smaller than paca'. Cost: whoami now opens the store, which it did not. It is written so that a failure there prints the identity exactly as before — whoami outside a board must not start erroring.
- **2026-09-14 13:38 · Alexander Sacharov** — In-progress lane, two things found on the way. (1) create now prints the file-mode line on a board with no git repository at all, not only on one whose remote is missing — the DoD asks for the mode 'in beiden Faellen' and refs.Usable() returns the same gitref.ErrNoRepo for both. Somebody who uses jaira without git therefore sees two extra lines per create. Left in deliberately; if the critique lane wants it gone, the place is fileModeReason/the else-branch in internal/cli/tickets.go, and the distinguishing test is whether the error string carries the noRemote tail. (2) releaseFromFile writes through s.Mutate rather than touching the file, because Mutate already takes the per-ticket lock, stamps updated-at/updated-by and calls Recorder.Record inside the lock — writing the file by hand here would have been a second, weaker write path for the one case where the ticket is about to leave the disk.
