---
id: 01M4WS20HSNGR8XKP3RCF914QH
title: "Epics: Funktionen als eigene Art auf eigenen Refs, nicht als Ticket"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Epic beschreibt eine Funktion für einen Menschen und sammelt die Tickets, die zu ihr gehören. Es liegt auf einer eigenen Ref, läuft durch keine Lane und zeigt, wie weit die Funktion ist - auch wenn sie aus 40 Tickets besteht."
context: |-
  Was falsch ist: Auf einem großen Board (~100 offene Tickets) sieht ein Teamlead ohne KI nicht, wie weit eine Funktion ist. Die Zahl wächst, weil beim Abschluss eines Tickets neue Folge-Tickets entstehen.
  Das Board passt ihm nicht. Er will Funktionen sehen, nicht Tickets.
  Warum kein Ticket als Epic: ein Epic hat keine DoD, keine Lanes, keinen Agenten. Es ist eine Beschreibung plus Liste.
  Warum kein Milestone: ein Milestone ist eine Runde (Zeit), ein Epic eine Funktion (Inhalt). Unfertiges wandert beim Milestone in die nächste Runde, beim Epic bleibt es drin.
  Schon da: --parent und --follows an Tickets, 'jaira links' liest den Baum. Tickets reisen schon auf Refs (Ref-Modus).
  Entscheidung Alex 14.10.: Die Mitgliedsliste steht im Epic, nicht am Ticket. So schreibt ein Teamlead, der ein Ticket einsortiert, nur die Ref des Epics - nie die Ref eines Tickets, an dem gerade ein Agent arbeitet. Keine Konflikte.
  Folge-Tickets (--follows/--parent eines Mitglieds) zählen beim Lesen automatisch mit.
  Nebenproblem: Agenten legen Folge-Tickets heute oft ohne --follows an. Dann fallen sie aus jedem Epic heraus.
  Grenze: Sync nur über git, kein Server. Synchron heißt: nach dem nächsten fetch.
definition-of-done: "Ein Epic ist eine Datei mit Frontmatter (Titel, Beschreibung für Menschen, Ticket-ids) auf einer eigenen Ref refs/jaira/epics/<id>, von Hand lesbar im Diff"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-14T08:39:27Z
updated-at: 2026-10-14T13:26:43Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-96726
claimed-at: 2026-10-14T13:23:08Z
mode: ""
outcome-what: "Epics: jaira epic create/add/rm/ls/show auf refs/jaira/epics/<name> (nur Ref, keine Datei), Fortschritt X von Y inkl. follows/parent-Folge-Tickets, ls --unassigned; jaira create setzt follows aus dem Ticket-Handle im Branch-Namen, --follows none schaltet ab"
outcome-why: "Ein Teamlead sieht auf ~100 Tickets nicht, wie weit eine Funktion ist; Folge-Tickets ohne --follows fielen aus jeder Gruppierung"
outcome-resolves: F914QH
review-summary: "Neues Paket core/epic: ein Epic ist eine Markdown-Datei (Frontmatter title/description/created-at, darunter eine Zeile '- <ticket-id>' pro Mitglied) auf refs/jaira/epics/<name>, ohne Datei im Worktree. Gelesen wird das ungesendete Outbox-Entry, sonst die lokale Ref; geschrieben über outbox.KindEpic mit Lease wie Milestones. Verliert ein Push das Rennen, merged refsync.rebaseEpic dreiseitig (Basis = Lease-Commit) die Mitglieder- und Titeländerungen und sendet bis zu 3x neu, statt den Eintrag zu verwerfen. CLI: jaira epic create/add/rm/ls/show, ls --unassigned, alles mit --json; Fortschritt = Mitglieder plus transitiv alle Tickets mit follows/parent darauf (Board + Logbook), erledigt = link.Index.Satisfied. jaira create setzt ohne --follows follows auf das Ticket, dessen Handle im Branch-Namen steht, sagt es an; --follows none schaltet ab. Dazu: refs.go meldet Epics beim Namen, milestones/links teilen memberRefs/boardIndex, docs/COMMANDS.md und core/release/notes/F914QH.md."
round: 3
review-gaps: "(1) Agent-Block nicht angepasst: Option 1 der DoD-5-Entscheidung war 'automatisch aus dem Branch + Hinweis im Agent-Block'; core/board/announce.go erwähnt weder jaira epic noch den Branch-Default noch --follows none — Agents lernen es nur aus der Ausgabe von create. (2) Nach einem verlorenen Rennen druckt 'jaira epic add' die Zahl vor dem Merge ('now lists 2'), obwohl die Ref danach 3 Mitglieder trägt; der Merge selbst bleibt still (Sent sagt nichts). Kosmetisch, im Zwei-Klon-Test nachgestellt. (3) Kein Befehl, um Titel/Beschreibung zu ändern oder ein Epic zu löschen; ohne Arbeitsdatei geht das nur mit git-Plumbing. Ein vertipptes Epic bleibt für immer in 'epic ls'. Nicht in der DoD verlangt. (4) rebaseEpic: der erste Versuch liest die lokale, noch alte Ref und verliert sicher ein zweites Mal, erst nach Fetch greift der Merge — kostet einen Push-Roundtrip, effektiv 2 statt 3 Versuche. Kein Fehler."
test-verdict: "pass: go vet ./... (RC=0), go test -race -count=1 ./... (RC=0); DoD 1-6 checked in the tree against their proofs; on a scratch remote with two clones: create/add/rm/ls/show, fetch to the second clone, parent-follow-up counted 0->1 of 2 after move to done and still counted after jaira logbook, ls --unassigned lists only the unrelated ticket, unknown epic exit 5, duplicate create exit 3"
review-verdict: "Der Diff erfüllt alle sechs DoD-Punkte; go test ./... und go vet sind grün, Ende-zu-Ende mit Bare-Remote nachgestellt (Epic nur auf Ref, X von Y inkl. Folge-Tickets, --unassigned, Branch-follows, zwei Klone fügen gleichzeitig hinzu und beide landen). Kein Defekt gefunden, der die Annahme blockiert. Offen für die Person: Gap (1) — der versprochene Hinweis im Agent-Block fehlt; entscheiden, ob er hier nachgezogen wird oder entfällt."
review-round: 2
review-check: "1. go build -o /tmp/jaira ./cmd/jaira  2. In einem Repo mit Remote: /tmp/jaira create Alpha --goal g --context c --dod d, dann /tmp/jaira create Beta --goal g --context c --dod d --follows <Alpha-Handle>  3. /tmp/jaira epic create 'Magic Login' <Alpha-Handle> --title 'Magic login' → meldet 'Filed as magic-login'  4. /tmp/jaira epic show magic-login → '0 of 2 done', Beta mit '(follow-up)'  5. git show refs/jaira/epics/magic-login:magic-login.md → Frontmatter + eine '- <id>'-Zeile; ls .jaira zeigt kein epics/  6. /tmp/jaira epic ls --unassigned → listet nur Tickets ohne Epic  7. git checkout -b feat/<alpha-handle-klein>; /tmp/jaira create Delta --goal g --context c --dod d → Zeile 'Follows <Alpha>, from branch …'; mit --follows none fehlt sie  8. Gleichzeitiges Hinzufügen aus zwei Klonen: go test ./internal/cli -run TestTwoClonesSortingIntoOneEpicBothLand"
---

# Epics: Funktionen als eigene Art auf eigenen Refs, nicht als Ticket

## Definition of Done

- [x] Ein Epic ist eine Datei mit Frontmatter (Titel, Beschreibung für Menschen, Ticket-ids) auf einer eigenen Ref refs/jaira/epics/<id>, von Hand lesbar im Diff
  proof: core/epic/epic.go:119 New + core/gitref/gitref.go:65 EpicRefName; TestAnEpicTravelsOnItsRefAndCountsFollowUps (git show refs/jaira/epics/<name>:<name>.md, no .jaira/epics on disk), TestEditsKeepEveryOtherLine
- [x] jaira epic create/add/rm/ls/show; ls und show nennen erledigt X von Y, inklusive Folge-Tickets über follows/parent
  proof: internal/cli/epics.go:18 newEpicCmd (create/add/rm/ls/show); core/epic/epic.go:378 Count + :397 Covered; TestAnEpicTravelsOnItsRefAndCountsFollowUps (0 of 2 incl. follows), TestCoveredWalksFollowsAndParent
- [x] Ein Ticket in ein Epic einsortieren schreibt nur die Ref des Epics, nie die des Tickets
  proof: core/epic/epic.go:319 Assign writes only RecordEpic; TestAnEpicTravelsOnItsRefAndCountsFollowUps checks the ticket ref SHA is unchanged; TestTwoClonesSortingIntoOneEpicBothLand (lost race merged, core/refsync/epic.go rebaseEpic)
- [x] Tickets in keinem Epic sind abfragbar (jaira epic ls --unassigned oder ähnlich), mit --json
  proof: internal/cli/epics.go:204 'jaira epic ls --unassigned [--json]', core/epic/epic.go:435 Unassigned; TestAnEpicTravelsOnItsRefAndCountsFollowUps, TestCoveredWalksFollowsAndParent
- [x] Ein Ticket, das während der Arbeit an einem anderen angelegt wird, bekommt --follows automatisch oder die Rolle setzt es; geklärt in brainstorm
  proof: internal/cli/tickets.go:153 branchTicket + --follows none; TestCreateFollowsTheTicketTheBranchIsNamedAfter
- [x] Tests, docs/COMMANDS.md, Release-Note core/release/notes/<ID>.md
  proof: core/epic/epic_test.go, internal/cli/epics_test.go, TestCreateFollowsTheTicketTheBranchIsNamedAfter (go test ./... exit 0); docs/COMMANDS.md jaira epic rows + create --follows; core/release/notes/F914QH.md
- [x] Der Agent-Block (core/board/announce.go) sagt in je einer Zeile: auf einem Ticket-Branch setzt jaira create follows auf dieses Ticket (--follows none schaltet ab); jaira epic (was ein Epic ist, ls/show/add)
  proof: core/board/announce.go:51 (branch follows) and :64 (jaira epic); same lines in CLAUDE.md/AGENTS.md; go test ./core/board ./internal/cli green; jaira validate reports no stale agent block

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] gitref: EpicPrefix, ReadEpic/WriteEpic/ListEpics/EpicSHA; outbox: KindEpic
- [x] core/epic: Format, parse, Add/Remove, Merge3, Load/LoadAll/Create/Assign über Refs-Interface
- [x] core/epic: Count (Mitglieder + follows/parent transitiv, Satisfied) und Unassigned; link.Index.All exportieren
- [x] refsync: EpicNames/EpicContent/RecordEpic; verlorenes Rennen eines Epics im Flush dreiseitig mergen und neu senden
- [x] cli: jaira epic create/add/rm/ls/show, ls --unassigned, --json
- [x] cli create: follows aus dem Branch-Handle, --follows none
- [x] Tests, docs/COMMANDS.md, core/release/notes/F914QH.md

## Progress
- **2026-10-14 08:46 · Alexander Sacharov** — parent session: jaira-34
- **2026-10-14 08:47 · Alexander Sacharov** — decision (worker, settled from code): Epics folgen dem Milestone-Muster 1:1 — Datei .jaira/epics/<name>.md auf der Platte, gespiegelt auf refs/jaira/epics/<name> über outbox/refsync, eingehend per fetch auf die Platte geschrieben. Grund: die ganze Infrastruktur (Lease, Outbox, Incoming) existiert schon für Milestones; ein reiner Ref-Speicher würde Boards ohne Remote ausschließen, die Milestones heute bedienen.
- **2026-10-14 08:47 · Alexander Sacharov** — decision (worker): Epic-id ist ein kebab-Name wie bei Milestones (tag-Regeln, Dateiname = id), 'jaira epic create <name> [id...] --title <t> --description <d>'. Frontmatter: title, description (einzeilig; längerer Text als Prosa unter der Frontmatter), created-at. Mitglieder: eine Zeile '- <ticket-id>' pro Ticket, wie Milestone. Grund: tippbar in 'jaira epic add <name> <id>', im Diff lesbar, gleiche Regeln wie Milestones.
- **2026-10-14 08:47 · Alexander Sacharov** — decision (worker): 'erledigt' = core/link Index.Satisfied (Terminal-Lane oder im Logbook abgelegt). Zählmenge eines Epics = Mitglieder plus transitiv alle Tickets, deren follows oder parent auf ein Element der Menge zeigt (Board + Logbook). 'unassigned' = offene Board-Tickets, die weder direkt noch über follows/parent in einem Epic stecken. Kein Logbook-Ablegen von Epics in diesem Ticket (Scope).
- **2026-10-14 08:47 · Alexander Sacharov** — question (DoD 5, offen für Alex): Wie bekommt ein Folge-Ticket sein --follows? 1) (Recommended) Automatisch aus dem Branch + Hinweis im Agent-Block: auf einem Branch, dessen Name einen Ticket-Handle trägt (feat/f914qh), setzt 'jaira create' ohne --follows follows auf dieses Ticket und sagt es an; '--follows none' schaltet ab. Auf master passiert nichts. Zuverlässig im Worktree-pro-Ticket-Betrieb, kein Raten über Sessions. 2) Nur Rolle/Agent-Block: der generierte CLAUDE.md-Block und die Rollen sagen 'gib --follows <aktuelles Ticket> mit'. Kleinster Code, verlässt sich aber auf genau die Disziplin, die heute fehlt. 3) Automatisch aus dem eigenen Claim: das einzige Ticket, das diese Person gerade hält. Bricht, sobald eine Person mehrere parallele Sessions hat (Claims sind pro Person, nicht pro Session) — bei Alex der Normalfall.
- **2026-10-14 08:49 · Alexander Sacharov** — decision (Alex 14.10., über jaira-34): DoD 5 = Option 1. Auf einem Branch, dessen Name einen Ticket-Handle trägt, setzt 'jaira create' ohne --follows follows auf dieses Ticket und sagt es an; '--follows none' schaltet ab; auf master passiert nichts.
- **2026-10-14 08:49 · Alexander Sacharov** — decision (Alex 14.10., über jaira-34, überstimmt die Worker-Entscheidung zum Milestone-Muster): Epics liegen NUR auf Refs. Keine Arbeitsdatei .jaira/epics/<name>.md auf der Platte; gelesen und geschrieben wird refs/jaira/epics/<name> direkt. Lease/Outbox/Fetch wiederverwenden, wo es passt.
- **2026-10-14 08:57 · Alexander Sacharov** — Interface für RFTQKP steht (an jaira-rftqkp-b4 geschickt): epic.Refs = *refsync.Syncer; Load/LoadAll/Create/Assign/Count/Unassigned. Verlorenes Rennen eines Epics: Box.Flush verwirft den Eintrag, deshalb trägt outbox.Result jetzt Entry; refsync.rebaseEpic merged dreiseitig (Basis = Lease-Commit, der lokal liegt) und schreibt sofort neu, bis 3 Runden, sonst zurück in die Queue. Mutation geprüft: ohne rebaseEpic schlägt TestTwoClonesSortingIntoOneEpicBothLand fehl.
- **2026-10-14 08:57 · Alexander Sacharov** — Pause (conversational): DoD 1-4 fertig, nicht committet; Diff und Commit-Zeile an jaira-34, warte auf Freigabe. Danach: Plan 6 (follows aus Branch-Handle in internal/cli/tickets.go um Zeile 192, '--follows none'), Plan 7 (docs/COMMANDS.md, core/release/notes/F914QH.md). jaira-rftqkp-b4 wartet auf die Commit-SHA von core/epic.
- **2026-10-14 13:12 · Alexander Sacharov** — critique: (1) internal/cli/refs.go refSubject/rejectedAdvice fehlen KindEpic — Epic-Meldungen nennen sonst den Namen per ticket.Handle verstümmelt; fixen. Bewusst stehen gelassen: (a) core/epic dupliziert frontEnd/Bytes/Add/Remove aus core/milestone (~50 Zeilen) statt einen gemeinsamen Zeilen-Datei-Typ zu ziehen — Milestone trägt status/color und Logbook-Ablage, ein gemeinsamer Typ koppelte zwei Formate, die getrennt wachsen sollen; ParseMember ist geteilt, damit die Mitgliedszeile nicht auseinanderläuft. (b) Covered liest über link.Index.All das ganze Logbook — nur in epic-Befehlen, nie beim Board-Rendern; nötig, weil abgelegte Folge-Tickets mitzählen.
- **2026-10-14 13:13 · Alexander Sacharov** — critique Runde 2: refs.go-Fund behoben (4400dddb, TestAnEpicIsReportedByItsName); sonst nichts — die zwei bewusst stehen gelassenen Punkte stehen in der critique-Note davor.
- **2026-10-14 13:23 · Alexander Sacharov** — review: alle DoD-Punkte erfüllt, Tests grün, E2E mit Bare-Remote nachgestellt. Offen für Alex: Option 1 versprach einen Hinweis im Agent-Block (core/board/announce.go) zu Branch-follows/--follows none und jaira epic — fehlt. Kleinkram: 'epic add' druckt nach verlorenem Rennen die Zahl vor dem Merge; kein Befehl zum Umbenennen/Löschen eines Epics.
- **2026-10-14 13:24 · Alexander Sacharov** — Abnahme-Seite: https://claude.ai/artifact/TzzQLMxeCQuumLtkD43Lxk (3 Person-Schritte statt 8 Review-Check-Schritten; Maschinen-Checks grün). Offen für Alex: Review-Gap 1 — Hinweis im Agent-Block (Branch-follows, jaira epic) hier nachziehen oder weglassen.
- **2026-10-14 13:25 · Alexander Sacharov** — decision (Alex 14.10., über jaira-34): Review-Gap 1 wird hier erledigt — Agent-Block in core/board/announce.go bekommt je eine kurze Zeile zu Branch-follows (--follows none) und zu jaira epic (ls/show/add). Die kleinen Review-Punkte (Zahl vor dem Merge, kein rename/delete, erste rebaseEpic-Runde) bleiben bewusst stehen.
