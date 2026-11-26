---
id: 01M6HZQN4HSDJ55MYC5AWQNRQN
title: "Wer ein Ticket, einen Milestone oder ein Epic zugewiesen bekommt, sieht das im Board als Nachricht"
status: testing
ready: true
creator: berkc
assignee: Berk
goal: "Weist mir jemand anderes ein Ticket zu (einzeln oder alle Tickets eines Milestones/Epics mit einem Befehl), zeigt das Board eine bleibende Nachricht mit Absender, bis ich das Ticket oeffne oder quittiere; das Ticket zeigt 'Von -> An'; gesehen liegt lokal unter ~/.jaira/state, nie im Ticket."
context: |-
  Wunsch berkc 31.10.2026, Punkt 5, im Original abgebrochen: 'after you get assigned a ticket or milestone (or epic) from a colleague or the web user (PM) you see it in the board as a ...'. Am 04.11. ergaenzt: 'as a message or something - weiss aber nicht wie man das umsetzen kann, muss man brainstormen'.
  Heute: Zuweisung landet per Ref beim naechsten jaira fetch; es gibt eine Desktop-Benachrichtigung (8566KF), im Board selbst nichts.
  Offen fuer brainstorm: Form (Banner oben, Markierung auf der Karte, Posteingang-Fenster), wann sie verschwindet (erstes Oeffnen, Quittieren), wo 'gesehen' gespeichert wird (lokal, nicht im Ticket - sonst Commit-Laerm), Gegenstueck im Web. Haengt an DVPAPE (fester Name pro Person).
definition-of-done: "Brainstorm hat Form, Verschwinden und Speicherort festgelegt"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-11-04T00:35:52Z
updated-at: 2026-11-26T10:05:06Z
updated-by: Berk
claimed-by: w-wqnrqn
claimed-at: 2026-11-26T09:58:12Z
question: ""
outcome-what: "The board reads what waits for you (core/inbox): a banner above the lanes, a ✉ mark on cards, an i list where enter opens or reads and x acks; a plain pull makes a worktree only from a handover, and take makes the guessed one."
outcome-why: "An assignment or message from a teammate must stay visible until the person looks at it, and a plain pull must not create worktrees nobody asked for."
outcome-resolves: "DoD 2, 6, 15, 17 and 19 are ticked with the tests that prove them (internal/tui/inbox_test.go, internal/cli/msg_test.go, core/inbox); the pull and take worktree rules are covered by TestPlainPullMakesNoWorktreeFromAHandleNamedBranch and TestTakeMakesTheWorktreeFromAHandleNamedBranch; DoD 20 stays open for a person."
review-summary: none
round: 2
back-from: critique
back-by: lane
back-count: 1
review-gaps: "removed internal/cli/handover.go shortSHA (copy of core/gitrepo shortSHA, now exported as gitrepo.ShortSHA); folded inbox.Record (forwarder with only test callers) and RecordNew into one Record returning what it added. Left: refsync decodeTicket (small, mirrors inline ParseDoc+Decode already repeated in refsync), inbox --all loading the file before AckAll and building groups it does not need (rare command, not a hot path), Arrival.Mine/Changed (still in fetch --json output), gitrepo CheckedOut vs core/board/gitignore.go:334 worktree list (different question), people active, milestone assign (peopleList is cached)"
---

# Wer ein Ticket, einen Milestone oder ein Epic zugewiesen bekommt, sieht das im Board als Nachricht

## Definition of Done

- [x] Brainstorm hat Form, Verschwinden und Speicherort festgelegt
  proof: Notiz 04.11.2026 Entscheidung berkc Option 1 + B1
- [x] Nach einer Zuweisung durch jemand anderen zeigt das Board die Nachricht bis zum Gesehen; Test
  proof: internal/tui/inbox_test.go TestInboxBannerAndCardMarkerShowWhatWaits (banner and card mark after a reload; nothing marked before); TestOpeningATicketsDetailAcksItsEntries; the entry survives a reload from disk in core/inbox TestRecordThenLoadKeepsTheEntry
- [-] Release-Notiz unter core/release/notes/
- [x] Eintrag nur bei Assignee-Wechsel auf mich durch jemand anderen; Notizen und eigene Zuweisung erzeugen keinen; Test core/refsync
  proof: core/refsync/inbox_test.go TestAnAssignmentByAnotherRecordsAnEntry, TestAReassignmentAwayRecordsATakenEntry, TestNotesAndSelfAssignmentsRecordNothing, TestTheFirstFetchIsABaseline
- [x] Eintrag ueberlebt 'jaira fetch' aus einer anderen Session und einen Board-Neustart, bis gesehen; Test
  proof: core/refsync/inbox_test.go TestAnAssignmentByAnotherRecordsAnEntry (a second syncer on the same clone still sees the entry) and TestAnAckedAssignmentDoesNotComeBack; entries are in inbox.json in the repo state dir
- [x] TUI: Banner, Kartenmarker und Liste (i); Oeffnen oder x quittiert; TUI-Test
  proof: internal/tui/inbox_test.go TestInboxBannerAndCardMarkerShowWhatWaits (banner, card mark); TestIOpensTheInboxAndXAcksTheRow (i list, x acks, esc closes); TestEnterOnATicketRowOpensItAndAcksIt; TestOpeningATicketsDetailAcksItsEntries; TestEnterOnAMessageShowsItsTextCleanedAndAcksIt
- [x] 'jaira inbox [--json] [--ack <id>|--all]' fuer Agenten
  proof: internal/cli/inbox_test.go TestInboxShowsAnAssignmentUntilItIsAcknowledged
- [x] Ein Befehl weist alle offenen Tickets eines Milestones einer Person zu (jaira milestone assign <name> <person>); mehrere Zuweisungen aus einem Milestone/Epic in einem Fetch sind ein Eintrag; Test
  proof: internal/cli/milestone_assign_test.go TestMilestoneAssignGivesTheReceiverOneInboxRow
- [x] Das Ticket zeigt, wer es wem zugewiesen hat ('alex -> berkc') in Detail und show; Test
  proof: internal/cli/handedshow_test.go TestShowSaysWhoHandedTheTicketOver; internal/tui/handover_detail_test.go TestDetailShowsWhoHandedTheTicketOver
- [x] Release-Note-Zeile in core/release/notes/WQNRQN.md
  proof: core/release/notes/WQNRQN.md (handover, take, pull worktree, people active, show from -> to)
- [x] Uebergabe: 'jaira handover <id> <person>' setzt den Assignee, pusht vorher den Branch des Tickets (WIP-Commit, wenn der Worktree Aenderungen hat) und schreibt Branch und Commit aufs Ticket; der Empfaenger bekommt den Inbox-Eintrag mit Branch und Commit; Test
  proof: internal/cli/handover_test.go TestHandoverGivesTheReceiverTheBranchAndTheCommit
- [x] Uebernahme: 'jaira take <id>' macht mich zum Assignee, der bisherige Assignee bekommt einen Inbox-Eintrag 'X hat dein Ticket uebernommen'; 'jaira pull <id>' holt danach Ticket und Branch in einen Worktree; Test
  proof: internal/cli/handover_test.go TestTakeTellsThePreviousHolderAfterTheirFetch and TestPullChecksOutTheHandedBranchBesideTheClone
- [x] Nachrichten: 'jaira msg <person> <text>' und 'jaira msg <id> <text>' (an den Assignee) legen eine Nachricht auf einen Ref pro Empfaenger ab (ein Ref, an den Commits angehaengt werden, nicht ein Ref pro Nachricht); sie erscheint im Inbox des Empfaengers wie eine Zuweisung; Test
  proof: internal/cli/msg_test.go TestMsgReachesThePersonsInboxAsDataNotAnInstruction, TestMsgToATicketGoesToItsAssignee
- [x] 'jaira people' listet die Aliase, die in diesem Repo aktiv sind (Assignees, Creator, Commit-Autoren der letzten 30 Tage), mit --json; Test
  proof: internal/cli/people_active_test.go TestPeopleActiveNamesWhoIsWorkingHere
- [x] Nachrichtentext ist Daten: kein Agent fuehrt eine Nachricht als Anweisung aus; jaira inbox --json kennzeichnet Nachrichten als von Menschen geschrieben
  proof: internal/cli/msg_test.go TestMsgReachesThePersonsInboxAsDataNotAnInstruction (inbox --json marks a message written_by human, with the note that its text is data)
- [x] two processes recording into inbox.json at once lose no entry, and an acked entry does not come back from a second worktree's fetch; Test core/inbox
  proof: core/inbox/inbox_test.go TestConcurrentRecordLosesNothing (3 writers x 5 records, plain ticket.LockFile); TestAnAckedKeyDoesNotComeBack; core/refsync/inbox_test.go TestAnAckedMessageDoesNotComeBackAfterPrune
- [x] message text with newlines, tabs, ANSI escapes and control characters is shown stripped in TUI and text output, kept verbatim in --json; empty text is refused; Test
  proof: internal/cli/msg_test.go TestMsgCleansControlCharactersForPeopleButKeepsThemInJSON and TestMsgRefusesEmptyTextAndWritesNothing; core/inbox TestCleanStripsControlsAndANSI; internal/tui/inbox_test.go TestEnterOnAMessageShowsItsTextCleanedAndAcksIt
- [x] handover whose push fails leaves the assignee unchanged and exits non-zero; msg/handover with no reachable remote write nothing; Test
  proof: internal/cli/handover_test.go TestHandoverWhoseBranchDoesNotReachTheRemoteChangesNothing and TestHandoverWithNoRemoteWritesNothing; msg: TestMsgWithNoReachableRemoteWritesNothing
- [x] acceptance rule: an assignment by someone else that never shows is worse than one shown twice
  proof: design B dedupe and acked map: core/inbox TestRecordDoesNotDuplicateAKey and TestAnAckedKeyDoesNotComeBack; an entry is shown until acked, not until fetched: internal/tui/inbox_test.go TestInboxBannerAndCardMarkerShowWhatWaits
- [ ] seen by a person: banner, card marker and the i-list on the board in a real terminal

## Options

- [x] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] P1 core/ticket: FieldAssignedBy 'assigned-by', FieldHandoverBranch, FieldHandoverCommit (+canonicalOrder after assignee); Store.Mutate stamps assigned-by=s.Actor when assignee changed to non-empty; create sets it when --assignee != creator; unit tests
- [x] P2 core/inbox: Entry/State, Load(dir), Record(dir, entries...) dedup by key and skip acked, Ack(dir, keys...), AckTicket(dir, id), AckAll; inbox.json under lock (ticket.LockFile) + WriteAtomic; prune acked >30d; Clean(text) strips control/ANSI; tests incl. concurrent Record and ack-survives-rerecord (DoD 16, 5)
- [x] P3 core/gitref: InboxPrefix 'refs/jaira/inbox/', InboxRefName(slug), ReadInbox(slug), WriteInbox(slug, content, lease) via refRead/refWrite; core/inbox Send(repo, slug, Message) = fetch, read, append JSON line, write with lease, retry on ErrRaceLost (people.Update pattern); Parse lines; tests with a bare remote
- [x] P4 core/refsync.Incoming: compare assignee at refs-seen SHA vs tip, classify assigned/taken, sender from assigned-by/updated-by, carry handover-branch/commit; first-ever fetch (no refs-seen.json) is a baseline; read inbox refs for Slug of each name in new Syncer.Me []string and record new message lines; write everything through core/inbox.Record into Store.RepoStateDir(); wire Me in cli/refs.go and tui/refs.go; refsync tests: notes and self-assign make no entry, other's assignment makes one, survives a second Incoming from another tree until acked (DoD 4, 5)
- [x] P5 cli 'jaira inbox [--json] [--ack <key|id>] [--all]': runs refs.Incoming when usable, prints grouped rows (bundle by batch+from+milestone/epic per design C), message rows marked; --json entries with kind, from, written_by:human + note for messages, branch/commit; tests (DoD 7, 8, 15)
- [x] P6 cli 'jaira msg <person|id> <text...>': id -> its assignee (refuse none, exit 3), canonical name via people list, empty text refused (exit 2), warn when person not active; offline -> nothing written, error; tests (DoD 13, 17, 18)
- [x] P7 cli 'jaira milestone assign <name> <person>': every open ticket of the milestone (not done/rejected) gets assignee=person in one command through the same write path as 'jaira set'; prints count; test that receiver gets ONE inbox row (DoD 8)
- [x] P8 cli 'jaira handover <id> <person>' (design F) and 'jaira take <id>' (design G); 'jaira pull' worktree step for handover-branch; tests with a bare remote: push fails -> assignee unchanged; receiver entry carries branch+commit; take gives previous assignee a 'taken' entry (DoD 11, 12, 18)
- [x] P9 cli 'jaira people': add active names (assignees, creators, git authors last 30 days, canonical) to text and --json 'active'; test (DoD 14)
- [x] P10 'alex -> berkc': printDetail and TUI detail row 'assignee' show '<assigned-by> -> <assignee>' when assigned-by is set and differs; show --json already carries the field; tests (DoD 9)
- [x] P11 TUI: load inbox on reload; banner line above lanes when entries exist ('✉ N new for you: ...  i open'); ✉ marker on cards with an entry; modeInbox list on 'i' (j/k, enter opens ticket or shows message stripped, x acks, esc); openDetail acks the ticket's entries; help text; TUI tests (DoD 2, 6)
- [x] P12 README '## Commands' + docs/COMMANDS.md entries for inbox, msg, handover, take, milestone assign; core/release/notes/WQNRQN.md one line per client-facing change (DoD 10)
- [x] P13 go vet + go test for touched packages (core/ticket core/inbox core/gitref core/refsync internal/cli internal/tui) and core/shipped; tick DoD with proofs; commit feat(WQNRQN)

## Progress
- **2026-11-04 21:27 · berkc** — BRAINSTORM (Code-Befund, dann Optionen, Empfehlung zuerst)

Was heute wirklich passiert (gelesen, nicht vermutet):
- Zuweisung kommt per refs/jaira/tickets/* an; refsync.Incoming (core/refsync/refsync.go:947) meldet Arrival{Mine, Changed}.
- TUI zeigt dann 30 s lang eine Zeile in der Fusszeile: "● CW2X0F is yours: <Titel>" (internal/tui/refs.go:150, view.go:1104, flashFor model.go:1001). Danach weg.
- "Gesehen" wird beim FETCH gespeichert, nicht beim Hinschauen: ~/.jaira/state/<tree>/refs-seen.json (refsync.go:123, saveSeen). Holt ein Agent oder 'jaira fetch' die Refs, ist die Nachricht verbraucht, bevor das Board sie je zeigt. Das ist der eigentliche Fehler.
- Changed = Ref-SHA hat sich bewegt, nicht "Assignee hat gewechselt": jede Notiz eines Kollegen auf einem Ticket, das schon mir gehoert, meldet erneut "is yours". Absender wird nicht gezeigt, obwohl updated-by im Ticket steht (core/ticket/schema.go:36).
- Milestones und Epics haben kein Assignee-/Owner-Feld (.jaira/milestones/*.md: nur color, created-at, Ticketliste). "Milestone zugewiesen" gibt es im Datenmodell nicht.
- Wer bin ich: JAIRA_USER > git user.name > $USER, Aliasse ~/.jaira/identity (core/identity/identity.go:16). DVPAPE (fester Name) ist noch backlog.
- jaira web ist read-only, kein Begriff "aktueller Nutzer".

Optionen:

1. (EMPFOHLEN) Bleibender Posteingang lokal + Banner + Kartenmarker
   Speicher: ~/.jaira/state/<repo>/inbox.json (RepoStateDir, alle Worktrees eines Clones teilen ihn), nie im Ticket.
   Eintrag entsteht, wenn der Assignee im Ref von "nicht ich" auf "ich" wechselt UND updated-by != ich. Verschwindet erst, wenn ich das Ticket oeffne oder 'x' druecke. 'jaira inbox [--json] [--ack]' fuer Agenten.
     jaira · default                                          berkc
     ✉ 2 neu fuer dich: CW2X0F von alex · pm-web (5 Tickets) von anna   i oeffnen  x gesehen
     backlog            todo               in-progress
     ┌──────────────┐   ┌──────────────┐
     │✉ CW2X0F alex │   │ ...          │
     │jaira web legt│
     └──────────────┘
     [i] ┌ Fuer dich ─────────────────────────────────┐
         │▸ CW2X0F  jaira web legt Tickets an  alex 2h │
         │  pm-web  Milestone, 5 Tickets        anna 1d│
         │ enter oeffnen · x gesehen · esc             │
         └─────────────────────────────────────────────┘
   Kostet: neue Datei + Diff alter/neuer Assignee in Incoming, ein Panel im TUI. Gibt auf: gesehen gilt pro Rechner, nicht pro Person.

2. Nur Markierung auf der Karte, gesehen = einmal geoeffnet
     todo
     ┌──────────────────┐
     │● NEU CW2X0F alex→│
     │jaira web legt ...│
     └──────────────────┘
   Kostet am wenigsten (ein Set von IDs lokal, ein Marker im Card-Render). Gibt auf: Karte in eingeklappter/gefilterter Lane oder auf anderem Board sieht man nicht; Milestone/Epic hat keine Karte; Agenten bekommen nichts.

3. Posteingang als Git-Ref pro Person (refs/jaira/inbox/<name>), quittieren pusht
   Gleiches Bild wie 1, zusaetzlich im Web und auf jedem Rechner derselbe Stand; ein PM sieht "noch nicht gesehen".
   Kostet: neue Ref-Art, Push bei jedem Quittieren, CAS-Konflikte, braucht DVPAPE zwingend (Name als Ref-Schluessel). Gibt auf: "lokal, kein Commit-Laerm" aus dem Kontext.

Warum 1: behebt genau den Fehler (gesehen beim Fetch statt beim Menschen), bleibt git-only ohne neuen Ref und ohne Laerm, funktioniert vor DVPAPE ueber identity.IsMe, und ein Web-Gegenstueck kann spaeter dieselbe Datei lesen (jaira web laeuft lokal). 3 bleibt als Ausbau, falls PMs "gesehen" sehen wollen.

Zweite Entscheidung (Milestone/Epic):
 B1 (empfohlen) Kein neues Feld: werden mehrere Tickets eines Milestones/Epics in einem Fetch mir zugewiesen, wird daraus EIN Eintrag "pm-web (5 Tickets) von anna".
 B2 Neues Feld owner: in Milestone- und Epic-Datei, eigene Nachricht "pm-web gehoert jetzt dir". Formataenderung, mehr Arbeit, eigenes Ticket wert.

Vorschlag geschaerftes Ziel: Weist mir jemand anderes ein Ticket zu (einzeln oder gebuendelt aus einem Milestone/Epic), zeigt das TUI-Board eine bleibende Nachricht mit Absender, bis ich das Ticket oeffne oder quittiere; "gesehen" liegt lokal unter ~/.jaira/state, nie im Ticket.
Vorschlag DoD:
 1 Eintrag nur bei Assignee-Wechsel auf mich durch jemand anderen; Notizen/eigene Zuweisung erzeugen keinen; Test core/refsync
 2 Eintrag ueberlebt 'jaira fetch' aus einer anderen Session und Board-Neustart, bis gesehen; Test
 3 TUI: Banner + Kartenmarker + Liste (i), Oeffnen oder x quittiert; TUI-Test
 4 'jaira inbox [--json] [--ack <id>|--all]'
 5 Mehrere Zuweisungen aus einem Milestone/Epic in einem Fetch = ein Eintrag
 6 Release-Notiz core/release/notes/WQNRQN.md
Nicht hier: Web-Gegenstueck (nach DVPAPE/CW2X0F eigenes Ticket).
- **2026-11-04 21:27 · berkc** — Nachtrag Brainstorm: es gibt kein 'jaira assign' - zugewiesen wird per 'set assignee=', 'tickets create --assignee' (internal/cli/tickets.go:460), pull/move (core/move/move.go:100). Erkennung muss daher in refsync.Incoming alten gegen neuen Assignee vergleichen, nicht an einem Befehl haengen. Das Cockpit kennt 'me' schon ueber 'jaira whoami --json' (core/mod/builtin/jaira-cockpit/hooks/register.tsx:105) und koennte 'jaira inbox --json' als Zaehler zeigen - eigenes Folgeticket. Ungeprueft: ob ein leeres refs-seen.json beim ersten Fetch alle meine Tickets als neu meldet.
- **2026-11-04 21:40 · berkc** — Entscheidung berkc 04.11.2026: Option 1 (lokaler Posteingang + Banner + Kartenmarker + Liste i) und B1 (gebuendelt). Zusatz wortgetreu: 'Es soll reichen zu sagen "milestone x zu alsa" um alle tickets im milestone zuzuweisen. auf dem ticket ein "Alsa -> Berk" waere aber zusaetzlich gut'. Also: ein Befehl weist alle Tickets eines Milestones zu (z.B. 'jaira milestone assign <name> <person>'), und das Ticket zeigt, wer an wen zugewiesen hat (Von -> An).
- **2026-11-25 23:04 · Berk** — berkc 26.11.2026, woertlich: 'kann man so dann auch die arbeit einem Kollegen sharen zb? Ich glaub wir hatten die idee dass wir anderen Usern tickets überreichen kann indem man ihn zum assignee macht. Wenn wir da noch adden dass man den aktuellsten commit mitgibt, dann könnte der Kolegge das einfach übernehmen und weiterführen oder? bzw. ein anderer Kollege kann sich das ticket nach absprache einfach ziehen (also sich selber assignee machen) und dann weiterführen. ... Vielleicht könnte man sogar für User ein Postfach anlegen wo man sich per ref nachrichten schreiben könnte ... so ein chat ticket wäre funny. ich sehe grad WQNRQN ist genau das glaub ich oder? wenn ja erweiter es um unsere neuen ideen und dann setz das ticket um.' Entscheidung Main-Session: Uebergabe/Uebernahme/Nachrichten/people als DoD 11-15 in WQNRQN; Cockpit-Icon als Folgeticket; 28PD7N nutzt danach denselben Inbox.
- **2026-11-25 23:06 · Berk** — worker w-wqnrqn: /home/berkc/code/X_General/jaira-WQNRQN existed as a stale detached worktree (0.5.10, uncommitted lane edits, not mine) - left untouched; working in /home/berkc/code/X_General/jaira-WQNRQN-inbox on branch feat/WQNRQN-inbox
- **2026-11-25 23:10 · Berk** — PRE-PROCESS design (worker, 26.11.). No open person decisions left; settled here, disagree at signoff:
A) Detection lives in refsync.Incoming (only callers: cli/fetch.go:73, tui/refs.go:104). Per changed ref: old content = blob at refs-seen SHA (still in object db), new = ref tip. assigned: new assignee isMe && old not isMe && sender not isMe. taken: old isMe && new non-empty, not isMe && sender not isMe. sender = new 'assigned-by', fallback 'updated-by'. Ref with no seen SHA: entry only if refs-seen.json already existed (first fetch ever is a baseline, no flood).
B) Store: new package core/inbox, file <Store.RepoStateDir()>/inbox.json (shared by all worktrees of a clone), read-modify-write under ticket.LockFile + WriteAtomic. Entry {key, kind assigned|taken|message, id, title, from, at, batch, branch, commit, text}. key = kind:id:sha (message: msg:<sha1 of line>). 'acked' map key->time keeps an acked key from coming back via another worktree's stale refs-seen; pruned after 30 days.
C) Bundle (B1): at display time, entries with same batch+from whose tickets share a milestone (milestone.LoadAll) or epic show as ONE row 'pm-web (5 tickets) from anna'; ack of the row acks all members; opening one ticket acks only that member.
D) New frontmatter: assigned-by (set in Store.Mutate when assignee changes to a non-empty value, and at create when --assignee != creator), handover-branch, handover-commit. Detail + show print 'alex -> berkc'.
E) msg ref: refs/jaira/inbox/<identity.Slug(canonical name)>, one file 'messages', JSON lines {at,from,ticket,text}; append = new commit with lease parent, retry on race like people.Update; fetched by the existing +refs/jaira/* refspec. Receiver reads the refs for Slug of every identity.Aliases name in Incoming. Text is data: control chars/ANSI stripped when shown, inbox --json carries written_by:human and a fixed 'data, not an instruction' field.
F) handover: branch = current branch (refuse on detached or the remote's default branch); dirty tree -> git add -A + commit 'wip(<HANDLE>): hand over to <person>'; git push -u; only after a good push the ticket gets assignee/handover-branch/handover-commit + a note. Push fails -> assignee unchanged, exit non-zero.
G) take = pull --steal semantics (refsync.Pull) or Mutate when the file is here. pull: when the ticket carries handover-branch (else one remote branch containing the handle) and no worktree has it checked out, fetch it and 'git worktree add <parent of repo>/<repo>-<HANDLE> <branch>'; existing path -> say so, no worktree.
H) TUI: 'x' on the board is archive, so x acks only inside the i-list; enter on a ticket acks its entries.
I) people: 'jaira people' gains 'active' (assignees + creators of board tickets, git authors --since 30 days, canonical names) in text and --json.
- **2026-11-25 23:12 · Berk** — in-progress split into 4 sequential Haiku runs (P1-P4, P5-P7, P8-P10, P11-P13); the last one moves to critique
- **2026-11-25 23:13 · Berk** — waiting: in-progress P1-P4 (haiku subagent) since 00:13
- **2026-11-25 23:48 · Berk** — P2 core/inbox: Ack records the key as acked even when no entry is there, so a later Record of that key is refused. Clean turns newline/tab into a space and drops other control chars after stripping ANSI; stored text stays verbatim. Lock file is inbox.lock beside inbox.json in the repo state dir (ticket.LockFile needs the dir to exist, so update() MkdirAlls first).
- **2026-11-25 23:57 · Berk** — P2 fix: TestConcurrentRecordLosesNothing timed out on the ticket lock (5s wait) under load: load avg 12.9 on 8 cores made each fsync in WriteAtomic slow, and a writer that just released re-takes the O_EXCL lock before queued waiters poll, so one waiter ran out of its wait. Decision: update() retries LockFile up to lockAttempts=5 instead of failing the write; the test is not shrunk. Lock release is already on every path (defer unlock). Commit 294aa5dea landed with this test failing: my pipe hid the go test exit code. Fixed here in a new commit.
- **2026-11-25 23:59 · Berk** — worker: inbox lock fix b12eb5542 retries ticket.LockFile 5x. Cause was starvation, not a leak: 8 in-process hot writers x 10 records with fsync'd WriteAtomic on WSL, and LockFile's backoff grows to 100ms so a waiter loses to fresh grabbers past its 5s. Real contention is 2-3 processes. Verified go test ./core/inbox -count=3 green. Critique may prefer a smaller test over the retry loop.
- **2026-11-26 00:27 · Berk** — P4 refsync.Incoming: decisions. (1) The baseline is refs-seen.json existing before the fetch; a clone's first fetch records no assignment entries but still records messages. (2) Entry.Batch is not in the design: set to updated-at truncated to the minute (RFC3339, UTC), so P5 bundles by batch+from, but a milestone assign that crosses a minute boundary splits. P7 may want an explicit batch. (3) Inbox entries are recorded before refs-seen is saved; a failed record returns an error so the next fetch repeats it. (4) Me is set from identity.Aliases in cli/refs.go attachRefs and tui/refs.go newSyncer. (5) The old ticket state is read with the new gitref.Repo.ReadAt(id, sha) from the last seen commit.
- **2026-11-26 00:32 · Berk** — worker: P1-P4 verified by worker (ticket, inbox -count=3, gitref, refsync all ok; build+vet clean). Decision on Entry.Batch: one value per Incoming run (the fetch time), not updated-at per minute — DoD 8 says 'in one fetch', and a minute boundary would split a milestone assign. P5 agent changes refsync.go:1181 accordingly.
- **2026-11-26 00:33 · Berk** — waiting: in-progress P5-P7 (haiku subagent #2) since 01:33
- **2026-11-26 00:39 · Berk** — P5-P7 worker: step 0 done. Batch is one time.Now() per Incoming run (RFC3339 UTC), passed to assignmentEntries; the old updated-at-per-minute batch split a milestone assign crossing a minute. Test: TestOneFetchGivesItsEntriesOneBatch. Design A-I read from ticket lines 142-150; C = display-time bundling by batch+from+milestone/epic, so core/inbox gets a Rows grouping.
- **2026-11-26 00:41 · Berk** — worker: TestAFinishedTicketThatNeverLandedIsReported failed once (refsync_test.go:831 'a finished ticket that never landed was not reported: []') on a full-package run at cc5f791ad while other agents loaded the machine; then green 5/5 alone and 2/2 full package. Treated as flaky under load, not caused by this ticket; testing lane should re-run it.
- **2026-11-26 00:42 · Berk** — P5 inbox: rows come from inbox.Rows (core/inbox/rows.go) bundling by batch+sender+milestone/epic (design C); milestone name before epic. --all acks everything (assumption: --all with no --ack clears the inbox, --ack with --all is exit 2). --ack by row key acks the whole row, by ticket id only that ticket's entries. Message rows are not bundled. Release note line in core/release/notes/WQNRQN.md.
- **2026-11-26 00:46 · Berk** — P6 msg: target is a ticket id when s.Load finds one (its assignee, message carries the ticket id), else a person name; canonical via people list, slug via core/identity.Slug. Empty text refused exit 2. Ticket without assignee exit 3. Offline/unreachable remote: inbox.Send returns error, nothing written (TestMsgWithNoReachableRemoteWritesNothing). Not ticked: DoD 15 (agent half is behaviour, not code), 17 (TUI half), 18 (handover half).
- **2026-11-26 00:49 · Berk** — P7 milestone assign: loads the milestone, skips done (terminal lane of the ticket's board) and already-assigned tickets, refuses ref-only members before writing anything (exit 3, pull first), then s.Mutate SetScalar(assignee) per ticket. Receiver test: one inbox row, three tickets (TestMilestoneAssignGivesTheReceiverOneInboxRow). Unticked: DoD 9 and onward. Note: a stray dod call overwrote DoD 7 proof with 'x' for a moment; restored to the inbox test.
- **2026-11-26 00:51 · Berk** — worker: P5-P7 verified (go test ./internal/cli -run 'Milestone|Msg|Inbox|TestEveryCommandIsDocumented' ok, build+vet clean). Kept: 'inbox --all' alone acks everything, with --ack it is exit 2 (matches DoD wording '--ack <id>|--all'); a CLI test for --all is still missing -> P8-P10 agent adds it.
- **2026-11-26 09:00 · Berk** — waiting: P8-P10 agent resumed after API rate-limit cutoff (uncommitted work in tree) since 10:00
- **2026-11-26 09:16 · Berk** — P8-P10 worker (26.11.): P8 done in 771106263. Decisions: (1) take on a ticket held on disk changes it in place via Store.Mutate; a ref-only ticket goes through refs.Pull(steal). (2) Found: a takeover was invisible to the previous holder whose own pull was never recorded in refs-seen, since only Incoming wrote it; Pull now calls noteSeen after its write, and only when a refs-seen record already exists so a never-fetched clone keeps its baseline. A clone that has never fetched still gets no taken entry (design A, deliberate). (3) handover commits dirty work with git add -A (design F, as written). (4) pull now also creates a worktree for every pull where exactly one remote branch contains the handle: a behaviour change for all pulls, not only handed-over tickets. (5) Dead end: git log --since stops at the first commit past the cutoff, so a backdated commit at HEAD hides newer ones; the people test orders commits accordingly.
- **2026-11-26 09:16 · Berk** — P9 worker: people active = assignees and creators of s.List() plus RecentAuthors 30.days (name, email only when the list knows the email), canonical via the people list; commit da42dfea9. P10 worker: people.List.Handed returns from -> to, used by printDetail and the TUI detail row; commit 812d15034. Unverified: show --json carrying assigned-by was not checked by a test; the field exists in core/ticket schema (AssignedBy json key not re-read).
- **2026-11-26 09:17 · Berk** — Open for the human: DoD 2 6 17 20 and the TUI inbox (P11-P13) are not touched here; DoD 15 and 19 not ticked.
- **2026-11-26 09:23 · Berk** — worker: P8-P10 verified (build+vet ok; go test ./internal/cli -run 'Handover|Take|Pull|Show|Inbox|People|Msg|Milestone|TestEveryCommandIsDocumented' ok 164s). Fix ordered for P11 agent: plain 'jaira pull' must NOT guess a worktree from remote branch names — agents pull constantly and would get surprise sibling worktrees (also under --json). Worktree only when the ticket carries handover-branch; the handle-named-branch guess only inside 'jaira take'.
- **2026-11-26 09:29 · Berk** — worker (w-wqnrqn, 26.11.) step 0: pullWorktree(s, id, guess bool). Plain pull passes false (handover-branch only). take passes true: take had no worktree call before, so it now makes the guessed one itself, and prints it unless --json. Dead end: the guess test needs 'git fetch' in the receiving clone first; jaira fetch only brings refs/jaira/*, so RemoteBranches is empty without it.
- **2026-11-26 09:36 · Berk** — worker (w-wqnrqn, 26.11.) P11 TUI decisions: enter on a message acks and shows its cleaned text; enter on a bundle row acks the whole row and lists its tickets (no single ticket to open); enter on a single ticket opens detail via showDetail, which acks. ackTicket writes only when the last reload showed an entry for that id. Epics are read (local for-each-ref) only while something waits. Help line added; TestHelpExplainsEveryCardMark needed 61 rows instead of 60 (the help grew one row, the card marks sit at row 59). Footer hint not added: the banner says 'i open'.
- **2026-11-26 09:43 · Berk** — critique round 1: core/refsync/refsync.go:1200 messageEntries re-reads the whole inbox ref on every fetch and core/inbox/inbox.go:29 forgets acked keys after 30 days, so every acked message comes back a month later (breaks DoD 16) - record messages only when the inbox ref tip moved since the last fetch (keep its sha in refs-seen like ticket refs), as Incoming already does for tickets
internal/tui/refs.go:181-188 announceArrivals still raises the desktop notification and the "is yours" flash on any ref move of a ticket of mine (the bug the brainstorm named: a colleague note re-announces) beside the new banner - drive both from the Assigned entries this fetch recorded, or drop them for the banner
internal/tui/inbox.go inboxGroups and inboxLine copy internal/cli/inbox.go inboxGroups, printInboxRow and verbFor by hand (the comment says so) - move the milestone/epic group map and the one-line row text into core/inbox (e.g. inbox.Groups(milestones, epics), inbox.Line(row, name)) and call it from both
core/inbox/inbox.go:155-163 retries ticket.LockFile 5x on top of its own timeout (up to 5x the wait) only to pass an 8-writer in-process test (note 2026-11-25 23:59) - drop the loop and shrink TestConcurrentRecordLosesNothing to the 2-3 processes that really contend
Left standing: handover git add -A is DoD 11 (WIP commit) and design F; pull/take worktree rule per note 2026-11-26 09:23; first-fetch baseline per note 00:27.
- **2026-11-26 09:52 · Berk** — in-progress send-back fixes (worker w-wqnrqn): (1) messages: refsync.messageEntries now keeps the last-seen inbox commit per slug in refs-seen as inbox:<slug>; a fetch returns only lines absent from the commit it last saw (keyed by line content), so a pruned ack does not come back; a first sight still records all lines. Old commit unreadable => every line counts as new (rare; history gc'd). (2) notifications: inbox.RecordNew returns what a fetch added; Arrival.Recorded marks an assigned entry that this fetch recorded. TUI announceArrivals and the CLI jaira fetch announceArrivals both notify only on Recorded; the TUI flash is gone (the banner replaces it). The CLI copy had the same bug as the TUI copy, so it is fixed in the same change. (3) grouping and row text: inbox.Groups and inbox.Line in core/inbox/rows.go; TUI and CLI call them. (4) lock: update takes ticket.LockFile once; TestConcurrentRecordLosesNothing is 3x5. Mutation check: disabling both tip guards makes TestAnAckedMessageDoesNotComeBackAfterPrune fail.
- **2026-11-26 09:54 · Berk** — critique round 2 (w-wqnrqn, 26.11.): read only round-1 findings and their fixes 2a6ade45f, b24172b1a, 1d784d287. All four answered: (1) inbox tip kept in refs-seen as inbox:<slug>, lines diffed against last-seen commit; checked that Departed and the vanished-ref loop skip inbox: keys (both need a local ticket file), harmless. (2) flash gone, notify only on Arrival.Recorded in TUI and CLI fetch. (3) inbox.Groups/inbox.Line in core/inbox/rows.go used by both. (4) lock retry loop dropped, test 3x5. build+vet ok; go test core/inbox core/gitref core/refsync, tui -run Inbox|Arriv|Help, cli -run Inbox|Fetch|Msg all ok. review-summary=none.
- **2026-11-26 09:57 · Berk** — optimize (w-wqnrqn, 26.11.): two behaviour-neutral cleanups, see review-gaps. Checked and kept: TUI flash removal left m.flash with other callers (refs.go:145, model.go:1124/1220), so nothing orphaned; peopleList() in the milestone-assign loop is cached (root.go:415). Tests: core/inbox gitrepo refsync full, cli -run Handover|Take|Pull|Inbox|Msg|Milestone|People|Show|Fetch|TestEveryCommandIsDocumented, tui -run Inbox|Handover|Arriv|Help|Detail all ok; build+vet clean.
- **2026-11-26 09:58 · Berk** — testing: round 2 started at 58e157861; bug hunt started in background (haiku); running gates, then DoD demand and function passes
- **2026-11-26 10:02 · Berk** — testing r2: go test -race -count=1 on core/ticket core/inbox core/gitref core/refsync core/people core/gitrepo core/shipped internal/tui internal/cli (3 skips): all 9 ok, exit 0; go build ./... and go vet on touched pkgs exit 0
- **2026-11-26 10:04 · Berk** — bughunt: change: blocker: core/identity/identity.go:117 - class: two different people reduce to one inbox slug, so mail to one lands in the other's inbox; example: 'Zoë' and 'Zo' both slug to 'zo', and names with no a-z or 0-9 (e.g. '张伟' and '李娜') both slug to 'unnamed'; the reader opens the inbox of Slug of each of its own names (core/refsync/refsync.go:1217, names from internal/cli/refs.go:42), and the sender writes to the same slug (internal/cli/msg.go:55), so each reads the other's mail; blame: Slug is 018021475 (before the ticket), the route at msg.go:55 is 5886c6b14 (names WQNRQN); read only, no go toolchain here so not run
- **2026-11-26 10:04 · Berk** — bughunt: change: blocker: core/inbox/send.go:75-81 - class: a send that reached the remote is reported as failed when the local ref update after the push fails, so a retry delivers the message twice; example: publish pushes the inbox commit (core/gitref/gitref.go:579) and then runs update-ref (gitref.go:584); if update-ref fails, for instance on a ref lock held by a concurrent jaira fetch, Send returns that error (send.go:80) although the inbox already holds the line, and a second run appends it again with a new at time; blame: gitref.go:578-588 is 03961b0ca, 69b246649, bb903c8d7 (none name the ticket), send.go:63-83 is 1fdc32973 (names WQNRQN); read only, not run, and the timing window is narrow
- **2026-11-26 10:04 · Berk** — bughunt: change: blocker: internal/cli/msg.go:47-55 - class: mail addressed to a name that is not on the people list is stored in an inbox its reader never opens, and the sender is told it was sent; example: a ticket assigned to 'Berk Canli' while the reader's own names are 'berkc' and the git email: the sender writes slug 'berk-canli' (msg.go:55), the reader opens only 'berkc' and the email slug (core/identity/identity.go Aliases), so the message is never shown; the only signal is the stderr line 'not on this board's people list; the message is sent anyway' (msg.go:49) before the 'Sent to' reply (msg.go:68); blame: msg.go:47-56 is 5886c6b14 (names WQNRQN); read only, not run
- **2026-11-26 10:04 · Berk** — bughunt: change: edge: core/gitref/gitref.go:468 - class: inbox refs are fetched by every clone, so anyone with the remote can read every person's messages and the ticket titles in them; example: Fetch uses +refs/jaira/*:refs/jaira/* (gitref.go:887) and the inbox namespace sits under refs/jaira/inbox/ (gitref.go:468), so a message to A is on the disk of every teammate who fetches; blame: the InboxPrefix line is 1fdc32973 (names WQNRQN), the Root refspec predates it; read only; whether messages are meant to be private to the recipient is a decision for a person, the code does not state it
- **2026-11-26 10:04 · Berk** — bughunt: done at 58e157861 - 4 findings (3 blocker, 1 edge; all change, 0 pre-existing), partial coverage, first hunt; files read whole: core/inbox/inbox.go core/inbox/send.go core/inbox/rows.go core/gitref/gitref.go internal/cli/msg.go internal/cli/handover.go
- **2026-11-26 10:05 · Berk** — testing: class: take decides on the local file, not the ref, so a take can report success and not take the ticket; example: alice's clone still had 6FM3JP as hers after bob's pull and alice's handover, 'jaira take 6FM3JP' printed 'is already yours' and exited 0, while refs/jaira/tickets/01M8BMY85V1ERGYSGVW96FM3JP on the remote said assignee bob (remote commit 7dddf77); fix at the cause: take reads the holder from the ref's ticket content (after fetch) and not from the working-tree file, and exits non-zero when the ref does not show the take; test over the class: take on a ticket whose local file is stale against the ref must take it (ref shows the new assignee) or refuse non-zero, never print 'already yours' for a holder the ref does not name (internal/cli/handover.go:134-147 is the check)
- **2026-11-26 10:05 · Berk** — testing: class: a receiver whose last-seen state already names them gets no inbox entry for a new assignment, because fetch records an entry only when the assignee changes from the last-seen commit (core/refsync/refsync.go assignmentEntries: mine(now.Assignee) && !mine(oldAssignee)); example: bob fetched 5NZJWE while assignee was bob, alice then set herself and handed it back to bob before bob fetched again, and bob's 'jaira inbox' said 'Nothing waiting for you' (the handover branch and commit never reached him); the same loss hits any re-handover to a holder who has not fetched since; fix at the cause: record an entry when the sender changed the assignee OR the handover-branch/handover-commit moved since the last-seen commit; test over the class: receiver is already the last-seen assignee, the ref moves away and back, the next fetch must record one entry carrying branch and commit
- **2026-11-26 10:05 · Berk** — testing: class: a receiver's pull rewrites assigned-by, so the record of who handed the ticket over is lost; example: remote history for 01M8BMY85V1ERGYSGVW96FM3JP: acabbe4 (alice) has assignee bob, assigned-by alice; 7dddf77 (bob, his pull) has assignee bob, assigned-by bob, updated-by bob; the assignee did not change, so store.go:1606-1610 should not have restamped it; cause not traced; fix at the cause: the pull path must not change assigned-by unless the assignee changes; test over the class: pull and any other write that leaves the assignee as it was keeps assigned-by unchanged, checked on the ref after the write
