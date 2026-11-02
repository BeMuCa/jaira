---
id: 01M6DSAAVH27P0HNBXZV182JT1
title: "Ein Epic laesst sich schliessen: danach nimmt es keine Tickets mehr an und verschwindet aus der Epic-Liste"
status: testing
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "jaira epic close <name> schliesst ein fertiges Epic. Ein geschlossenes Epic nimmt keine neuen Tickets mehr an und steht nicht mehr in der normalen Epic-Liste (CLI und web). Wieder oeffnen geht."
context: |-
  Was heute fehlt (Alex, 02.11.2026):
  - jaira epic kennt nur create/add/rm/set/show/ls (core/epic/epic.go, Ref refs/jaira/epics/<name>). Es gibt keinen Zustand 'geschlossen'.
  - 'jaira epic add' legt ein Ticket auch in ein Epic, das laengst fertig ist. Beispiel: web-0-2 steht auf '4 of 4 done' und nimmt weiter Tickets an.
  - Fertige Epics bleiben fuer immer in 'jaira epic ls' und in der Epic-Liste im web (core/web/ui/src/widgets/epic-list).
  Vorbild: Milestones haben das schon. 'jaira logbook <milestone>' legt einen fertigen Milestone ab; E70WYN markiert einen Ref mit 'status: filed', statt ihn zu loeschen. Ein Epic sollte es genauso machen: Feld im Epic-Ref, kein Loeschen. Andere Klone sehen den Zustand dann nach dem naechsten Sync.
  Offen fuer brainstorm:
  - Darf man ein Epic mit offenen Tickets schliessen? Vorschlag: verweigern (exit 3), --force schliesst trotzdem.
  - Sollen beim Schliessen die done-Tickets des Epics gleich mit 'jaira logbook' abgelegt werden (Flag)? Alex hat das noch nicht entschieden.
  - Was passiert mit einem Folge-Ticket (follows:), das spaeter auf ein Mitglied entsteht? Es zaehlt heute automatisch zum Epic.
definition-of-done: "jaira epic close <name> setzt das Epic auf geschlossen (Feld im Epic-Ref, kein Loeschen); jaira epic reopen <name> nimmt es zurueck; Test"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-11-02T09:26:46Z
updated-at: 2026-11-02T09:45:10Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-70421
claimed-at: 2026-11-02T09:33:25Z
outcome-what: "Epics lassen sich schliessen und wieder oeffnen: 'status: closed' im Epic-Ref (core/epic Close/Reopen), CLI 'epic close [--force]', 'epic reopen', 'epic add --force', 'epic ls --all'; web: POST /api/epics/{name}/close|reopen, force beim Zuordnen, Schalter 'Show closed', Knoepfe Close/Reopen, Bestaetigung bei Drop auf ein geschlossenes Epic und beim Schliessen mit offenen Tickets."
outcome-why: "Fertige Epics nahmen weiter Tickets an (web-0-2 '4 of 4 done') und blieben fuer immer in 'epic ls' und der Epic-Liste im web."
outcome-resolves: "Alle 7 DoD-Punkte: close/reopen ohne Loeschen, add verweigert mit exit 3 und nennt --force, --force laesst es geschlossen, ls/--all, web ausgeblendet/Schalter, web Knoepfe, Release-Note."
review-summary: none
round: 2
review-gaps: "zweimal '[closed]'-Markierung (epic ls, epic show) in closedMark() in internal/cli/epics.go gefaltet; IndexFunc>=0 durch ContainsFunc ersetzt (core/epic/epic.go Assign). Gelassen: dropField neben setField in core/epic/epic.go — kein vorhandener Helfer im Repo entfernt eine Frontmatter-Zeile (grep nach drop/remove/unsetField leer); written-Variable in Close/Reopen, weil edit() nicht meldet, ob geschrieben wurde."
---

# Ein Epic laesst sich schliessen: danach nimmt es keine Tickets mehr an und verschwindet aus der Epic-Liste

## Definition of Done

- [x] jaira epic close <name> setzt das Epic auf geschlossen (Feld im Epic-Ref, kein Loeschen); jaira epic reopen <name> nimmt es zurueck; Test
  proof: core/epic/epic.go:658 Close / :676 Reopen (status: closed im Ref); internal/cli/epics.go:288,327; TestCloseAndReopenWriteOneLine, TestEpicCloseAndReopen
- [x] jaira epic add in ein geschlossenes Epic verweigert ohne --force mit exit 3 und nennt --force; Test
  proof: core/epic/epic.go:632 ErrClosed; internal/cli/epics.go:64 exit 3 nennt --force; TestAssignToAClosedEpic, TestEpicCloseAndReopen
- [x] jaira epic ls zeigt geschlossene Epics nicht, jaira epic ls --all zeigt sie mit Markierung; Test
  proof: internal/cli/epics.go:406 (ls ohne --all filtert), [closed]-Markierung; TestEpicCloseAndReopen
- [x] Die Epic-Liste im web blendet geschlossene Epics aus oder zeigt sie getrennt; Test
  proof: core/web/ui/src/widgets/epic-list/ui/EpicList.svelte:35 blendet geschlossene aus, Schalter 'Show closed'; visible() in features/close-epic/model/close.ts, Test 'closed epics are hidden unless asked for'; Epic.closed core/web/web.go, TestEpicsCarryClosed
- [x] Release-Note-Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/182JT1.md
- [x] Im web laesst sich ein Epic schliessen und wieder oeffnen (Knopf in der Epic-Liste/Epic-Ansicht); Test
  proof: core/web/ui/src/widgets/epic-list/ui/EpicList.svelte:46 close() / :91 Reopen; POST /api/epics/{name}/close|reopen core/web/web.go:289,444; TestCloseReopenAndForce, TestWebClosesReopensAndForcesIntoAClosedEpic
- [x] jaira epic add in ein geschlossenes Epic geht nur mit --force, das Epic bleibt dabei geschlossen; im web nur nach ausdruecklicher Bestaetigung; Test
  proof: CLI: core/epic/epic.go Assign force, internal/cli/epics.go add --force, TestAssignToAClosedEpic, TestEpicCloseAndReopen; web: core/web/ui/src/widgets/epic-list/ui/EpicList.svelte:39 drop() fragt per confirm, internal/cli/web_source.go:352 409 ohne force, TestWebClosesReopensAndForcesIntoAClosedEpic

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/epic: Closed-Feld (status: closed), Close/Reopen, ErrClosed in Assign ohne force, ErrOpenTickets beim Schliessen, Merge uebernimmt Status; Tests
- [x] CLI: epic close/reopen [--force], epic add --force, epic ls --all mit Markierung, JSON closed; Tests
- [x] web: Epic.closed, POST /api/epics/{name}/close|reopen, Assign mit force, 409 ohne; Go-Tests
- [x] UI: Schalter 'geschlossene zeigen', Knopf schliessen/oeffnen, Bestaetigung bei Drop auf geschlossenes Epic; ui-Tests, dist neu bauen
- [x] Release-Note core/release/notes/182JT1.md

## Progress
- **2026-11-02 09:31 · Alexander Sacharov** — Entschieden von Alex (02.11.2026, ueber Teamlead): (1) Schliessen und Wieder-Oeffnen geht aus dem web UND aus der CLI. (2) Ein Ticket in ein geschlossenes Epic legen verweigert jaira epic add mit exit 3; nur mit --force geht es trotzdem, ohne das Epic wieder zu oeffnen. Im web dasselbe: Hinzufuegen zu einem geschlossenen Epic nur mit ausdruecklicher Bestaetigung. Damit ist der Vorschlag 'Hinweis auf reopen' im DoD 2 ersetzt. Offen bleibt fuer brainstorm: Schliessen mit offenen Tickets, und ob done-Tickets beim Schliessen ins Logbuch gehen.
- **2026-11-02 09:33 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/43623.sock
- **2026-11-02 09:34 · Alexander Sacharov** — Offene Entscheidungen vor dem Plan, vom Worker gesetzt (Alex kann widersprechen):
- Schliessen mit offenen Tickets: verweigert mit exit 3, --force schliesst trotzdem (Vorschlag aus dem Kontext, gleiches Muster wie Alex' Entscheidung fuer add). Web: Bestaetigungsdialog, dann mit force.
- done-Tickets beim Schliessen ins Logbuch: wird NICHT gebaut. Steht in keinem DoD-Punkt; Scope-Disziplin. Waere ein eigenes Ticket, wenn Alex es will.
- Folge-Tickets (follows:/parent:) auf Mitglieder eines geschlossenen Epics: zaehlen weiter wie heute. Sie laufen nicht ueber 'epic add', also verweigert nichts.
- Format: Frontmatter-Zeile 'status: closed' im Epic-Ref (wie 'status: filed' beim Milestone); reopen nimmt die Zeile wieder raus. Merge uebernimmt den Status von der Seite, die ihn geaendert hat.
- Web: geschlossene Epics standardmaessig ausgeblendet, ein Schalter 'geschlossene zeigen' zeigt sie markiert, mit Knopf 'wieder oeffnen'; offene haben 'schliessen'. Spiegelt 'epic ls --all'.
- **2026-11-02 09:36 · Alexander Sacharov** — core/epic: Tests brauchen fuer 'fertig' ein link.Index ueber lane.Load("") (Default-Lanes, done ist terminal); link.Build(nil, nil, ...) kennt keine terminale Lane und zaehlt nichts als fertig. Close prueft offen/fertig ueber epic.Count, also zaehlen Folge-Tickets mit.
- **2026-11-02 09:42 · Alexander Sacharov** — web: Bestaetigung ist window.confirm (features/close-epic/model/close.ts liefert den Fragetext, null = keine Frage). Der Server verweigert ohne force mit 409; die Seite fragt VOR dem Request, statt erst auf 409 zu reagieren, weil sie done/total und closed schon kennt. 'Show closed' ist lokaler Zustand in EpicList, nicht in der URL. tsc-Fehler in vite.config.ts (node:url) ist vorbestehend, nicht von diesem Ticket.
- **2026-11-02 09:43 · Alexander Sacharov** — critique Runde 1: (1) Release-Note Zeile 2 ohne '!' — niemand muss jetzt handeln. (2) reportClosing: state-Parameter weg, e.Closed lesen. Stehen gelassen: Close/Reopen als zwei Funktionen statt einer SetClosed (Close braucht den Index, Reopen nicht; zwei Namen lesen sich an den CLI-Aufrufern klarer); optionaler Body bei POST close (reopen sendet keinen, close kann ohne).
- **2026-11-02 09:44 · Alexander Sacharov** — critique Runde 2: beide Befunde aus Runde 1 umgesetzt (b76f9c10), nichts Neues.
