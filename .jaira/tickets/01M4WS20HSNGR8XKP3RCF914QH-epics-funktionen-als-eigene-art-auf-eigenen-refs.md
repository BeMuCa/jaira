---
id: 01M4WS20HSNGR8XKP3RCF914QH
title: "Epics: Funktionen als eigene Art auf eigenen Refs, nicht als Ticket"
status: testing
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
updated-at: 2026-10-14T13:16:18Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-81098
claimed-at: 2026-10-14T08:46:03Z
mode: ""
outcome-what: "Epics: jaira epic create/add/rm/ls/show auf refs/jaira/epics/<name> (nur Ref, keine Datei), Fortschritt X von Y inkl. follows/parent-Folge-Tickets, ls --unassigned; jaira create setzt follows aus dem Ticket-Handle im Branch-Namen, --follows none schaltet ab"
outcome-why: "Ein Teamlead sieht auf ~100 Tickets nicht, wie weit eine Funktion ist; Folge-Tickets ohne --follows fielen aus jeder Gruppierung"
outcome-resolves: F914QH
review-summary: "none"
round: 2
review-gaps: "folded the member-list JSON of milestoneJSON and epicJSON into memberRefs (internal/cli/milestones.go); links.go now calls boardIndex (internal/cli/links.go) instead of its own copy of lane.Load + partial s.List + link.Build; epic.LockName unexported (one caller). Left: core/epic frontEnd/Bytes/Add/Remove mirror core/milestone by decision (critique note); Count rebuilds the follows/parent map per epic in 'epic ls' — O(epics x tickets), only on epic commands, not worth a second API"
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
