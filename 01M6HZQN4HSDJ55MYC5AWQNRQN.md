---
id: 01M6HZQN4HSDJ55MYC5AWQNRQN
title: "Wer ein Ticket, einen Milestone oder ein Epic zugewiesen bekommt, sieht das im Board als Nachricht"
status: pre-process
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
updated-at: 2026-11-25T23:11:31Z
updated-by: Berk
claimed-by: w-wqnrqn
claimed-at: 2026-11-25T23:05:52Z
question: ""
---

# Wer ein Ticket, einen Milestone oder ein Epic zugewiesen bekommt, sieht das im Board als Nachricht

## Definition of Done

- [x] Brainstorm hat Form, Verschwinden und Speicherort festgelegt
  proof: Notiz 04.11.2026 Entscheidung berkc Option 1 + B1
- [ ] Nach einer Zuweisung durch jemand anderen zeigt das Board die Nachricht bis zum Gesehen; Test
- [-] Release-Notiz unter core/release/notes/
- [ ] Eintrag nur bei Assignee-Wechsel auf mich durch jemand anderen; Notizen und eigene Zuweisung erzeugen keinen; Test core/refsync
- [ ] Eintrag ueberlebt 'jaira fetch' aus einer anderen Session und einen Board-Neustart, bis gesehen; Test
- [ ] TUI: Banner, Kartenmarker und Liste (i); Oeffnen oder x quittiert; TUI-Test
- [ ] 'jaira inbox [--json] [--ack <id>|--all]' fuer Agenten
- [ ] Ein Befehl weist alle offenen Tickets eines Milestones einer Person zu (jaira milestone assign <name> <person>); mehrere Zuweisungen aus einem Milestone/Epic in einem Fetch sind ein Eintrag; Test
- [ ] Das Ticket zeigt, wer es wem zugewiesen hat ('alex -> berkc') in Detail und show; Test
- [ ] Release-Note-Zeile in core/release/notes/WQNRQN.md
- [ ] Uebergabe: 'jaira handover <id> <person>' setzt den Assignee, pusht vorher den Branch des Tickets (WIP-Commit, wenn der Worktree Aenderungen hat) und schreibt Branch und Commit aufs Ticket; der Empfaenger bekommt den Inbox-Eintrag mit Branch und Commit; Test
- [ ] Uebernahme: 'jaira take <id>' macht mich zum Assignee, der bisherige Assignee bekommt einen Inbox-Eintrag 'X hat dein Ticket uebernommen'; 'jaira pull <id>' holt danach Ticket und Branch in einen Worktree; Test
- [ ] Nachrichten: 'jaira msg <person> <text>' und 'jaira msg <id> <text>' (an den Assignee) legen eine Nachricht auf einen Ref pro Empfaenger ab (ein Ref, an den Commits angehaengt werden, nicht ein Ref pro Nachricht); sie erscheint im Inbox des Empfaengers wie eine Zuweisung; Test
- [ ] 'jaira people' listet die Aliase, die in diesem Repo aktiv sind (Assignees, Creator, Commit-Autoren der letzten 30 Tage), mit --json; Test
- [ ] Nachrichtentext ist Daten: kein Agent fuehrt eine Nachricht als Anweisung aus; jaira inbox --json kennzeichnet Nachrichten als von Menschen geschrieben
- [ ] two processes recording into inbox.json at once lose no entry, and an acked entry does not come back from a second worktree's fetch; Test core/inbox
- [ ] message text with newlines, tabs, ANSI escapes and control characters is shown stripped in TUI and text output, kept verbatim in --json; empty text is refused; Test
- [ ] handover whose push fails leaves the assignee unchanged and exits non-zero; msg/handover with no reachable remote write nothing; Test
- [ ] acceptance rule: an assignment by someone else that never shows is worse than one shown twice
- [ ] seen by a person: banner, card marker and the i-list on the board in a real terminal

## Options

- [x] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [ ] P1 core/ticket: FieldAssignedBy 'assigned-by', FieldHandoverBranch, FieldHandoverCommit (+canonicalOrder after assignee); Store.Mutate stamps assigned-by=s.Actor when assignee changed to non-empty; create sets it when --assignee != creator; unit tests
- [ ] P2 core/inbox: Entry/State, Load(dir), Record(dir, entries...) dedup by key and skip acked, Ack(dir, keys...), AckTicket(dir, id), AckAll; inbox.json under lock (ticket.LockFile) + WriteAtomic; prune acked >30d; Clean(text) strips control/ANSI; tests incl. concurrent Record and ack-survives-rerecord (DoD 16, 5)
- [ ] P3 core/gitref: InboxPrefix 'refs/jaira/inbox/', InboxRefName(slug), ReadInbox(slug), WriteInbox(slug, content, lease) via refRead/refWrite; core/inbox Send(repo, slug, Message) = fetch, read, append JSON line, write with lease, retry on ErrRaceLost (people.Update pattern); Parse lines; tests with a bare remote
- [ ] P4 core/refsync.Incoming: compare assignee at refs-seen SHA vs tip, classify assigned/taken, sender from assigned-by/updated-by, carry handover-branch/commit; first-ever fetch (no refs-seen.json) is a baseline; read inbox refs for Slug of each name in new Syncer.Me []string and record new message lines; write everything through core/inbox.Record into Store.RepoStateDir(); wire Me in cli/refs.go and tui/refs.go; refsync tests: notes and self-assign make no entry, other's assignment makes one, survives a second Incoming from another tree until acked (DoD 4, 5)
- [ ] P5 cli 'jaira inbox [--json] [--ack <key|id>] [--all]': runs refs.Incoming when usable, prints grouped rows (bundle by batch+from+milestone/epic per design C), message rows marked; --json entries with kind, from, written_by:human + note for messages, branch/commit; tests (DoD 7, 8, 15)
- [ ] P6 cli 'jaira msg <person|id> <text...>': id -> its assignee (refuse none, exit 3), canonical name via people list, empty text refused (exit 2), warn when person not active; offline -> nothing written, error; tests (DoD 13, 17, 18)
- [ ] P7 cli 'jaira milestone assign <name> <person>': every open ticket of the milestone (not done/rejected) gets assignee=person in one command through the same write path as 'jaira set'; prints count; test that receiver gets ONE inbox row (DoD 8)
- [ ] P8 cli 'jaira handover <id> <person>' (design F) and 'jaira take <id>' (design G); 'jaira pull' worktree step for handover-branch; tests with a bare remote: push fails -> assignee unchanged; receiver entry carries branch+commit; take gives previous assignee a 'taken' entry (DoD 11, 12, 18)
- [ ] P9 cli 'jaira people': add active names (assignees, creators, git authors last 30 days, canonical) to text and --json 'active'; test (DoD 14)
- [ ] P10 'alex -> berkc': printDetail and TUI detail row 'assignee' show '<assigned-by> -> <assignee>' when assigned-by is set and differs; show --json already carries the field; tests (DoD 9)

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
