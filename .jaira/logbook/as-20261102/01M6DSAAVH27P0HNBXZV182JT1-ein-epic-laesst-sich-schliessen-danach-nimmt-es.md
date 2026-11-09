---
id: 01M6DSAAVH27P0HNBXZV182JT1
title: "Ein Epic laesst sich schliessen: danach nimmt es keine Tickets mehr an und verschwindet aus der Epic-Liste"
status: done
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
commits:
  - e60357b786a076788144e840024a11f93391a48c
  - b76f9c108b815ffd4f7cc01b901454b83a3359d3
  - 23d42f61d8b230082499dced7d3ecc24007d9290
created-at: 2026-11-02T09:26:46Z
updated-at: 2026-11-02T10:24:49Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-70421
claimed-at: 2026-11-02T09:33:25Z
outcome-what: "Epics lassen sich schliessen und wieder oeffnen: 'status: closed' im Epic-Ref (core/epic Close/Reopen), CLI 'epic close [--force]', 'epic reopen', 'epic add --force', 'epic ls --all'; web: POST /api/epics/{name}/close|reopen, force beim Zuordnen, Schalter 'Show closed', Knoepfe Close/Reopen, Bestaetigung bei Drop auf ein geschlossenes Epic und beim Schliessen mit offenen Tickets."
outcome-why: "Fertige Epics nahmen weiter Tickets an (web-0-2 '4 of 4 done') und blieben fuer immer in 'epic ls' und der Epic-Liste im web."
outcome-resolves: "Alle 7 DoD-Punkte: close/reopen ohne Loeschen, add verweigert mit exit 3 und nennt --force, --force laesst es geschlossen, ls/--all, web ausgeblendet/Schalter, web Knoepfe, Release-Note."
review-summary: "Ein Epic bekommt einen Zustand 'geschlossen': die Frontmatter-Zeile 'status: closed' im Epic-Ref (core/epic/epic.go), reopen nimmt die Zeile wieder heraus, geloescht wird nichts. Neu: 'jaira epic close [--force]' (verweigert mit exit 3, solange Mitglieder oder Folge-Tickets offen sind), 'jaira epic reopen', 'jaira epic add --force' (sonst exit 3 bei geschlossenem Epic, Epic bleibt dabei geschlossen), 'jaira epic ls' laesst geschlossene weg, '--all' zeigt sie mit '[closed]', JSON traegt 'closed'. epic.Merge uebernimmt den Status von der Seite, die ihn geaendert hat. Web: Feld closed, POST /api/epics/{name}/close|reopen, force beim Zuordnen, 409 ohne force; die Epic-Liste versteckt geschlossene hinter 'Show closed (n)', hat Knoepfe Close/Reopen und fragt per confirm vor Drop auf ein geschlossenes Epic und vor dem Schliessen mit offenen Tickets. Release-Note mit zwei Zeilen, dist neu gebaut."
round: 2
review-gaps: |-
  note: Rennen zwischen zwei Klonen: Assign prueft "geschlossen" am lokalen Ref; verliert der Push gegen ein inzwischen entfernt geschlossenes Epic, nimmt epic.Merge (core/refsync/epic.go:112) das neue Mitglied ohne --force auf. Selten und ohne Datenverlust, aber die Sperre ist nur lokal
  note: Der Web-Dialog fragt anhand des geladenen Seitenstands; wurde das Epic seit dem Laden anderswo geschlossen, kommt statt der Frage ein 409-Fehler - die Meldung nennt die Bestaetigung, bietet sie aber nicht an
  note: Kein Test fuer epic.Merge mit status (eine Seite schliesst, die andere fuegt hinzu)
test-verdict: "pass: go test -race -count=1 ./... RC=0, UI npm test 15/15, svelte-check nur der vorbestehende vite.config.ts-Fehler; DoD 1-7 im Baum geprueft; Binary auf Scratch-Board: close ohne force exit 3 mit --force-Hinweis, force schliesst, add exit 3 / --force bleibt closed, ls blendet aus, ls --all markiert [closed], Ref traegt status: closed, reopen nimmt es zurueck; jaira web: close/drop 409 ohne force, 200 mit, reopen 200, Bundle enthaelt 'Show closed'. Nicht geprueft: Klick im Browser (confirm-Dialog)."
review-verdict: "Der Diff erfuellt alle 7 DoD-Punkte; keine Fehler gefunden. go build/vet sauber, go test ./core/epic ./core/web und internal/cli -run 'Epic|Web' gruen, UI-Tests 3/3 gruen, vite build reproduziert dist bit-gleich (index-Cn-OkEy4.js, index-CjHLen2W.css)."
review-round: 2
review-check: "1. go build -o /tmp/jaira ./cmd/jaira (im Worktree)  2. In einem Repo mit Board und Remote: /tmp/jaira epic create probe, dann /tmp/jaira epic add probe <offene-ticket-id>  3. /tmp/jaira epic close probe -> exit 3, Meldung nennt --force  4. /tmp/jaira epic close probe --force -> 'Epic \"probe\" closed on refs/jaira/epics/probe'  5. git show refs/jaira/epics/probe:probe.md -> Zeile 'status: closed', Mitglied noch da  6. /tmp/jaira epic ls -> probe fehlt; /tmp/jaira epic ls --all -> probe mit '[closed]'  7. /tmp/jaira epic add probe <andere-id> -> exit 3; mit --force -> geht, epic ls zeigt probe weiter nicht  8. /tmp/jaira web oeffnen: probe fehlt in der Epic-Liste, Haken 'Show closed (1)' zeigt es grau mit Badge 'closed' und Knopf Reopen  9. Ticket auf probe ziehen -> Rueckfrage 'is closed ... stays closed'; Abbrechen schreibt nichts  10. Reopen klicken -> probe erscheint ohne Haken wieder, Knopf heisst Close"
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
- **2026-11-02 09:47 · Alexander Sacharov** — review an frischen Agenten uebergeben (run-lane.sh --no-worktree, Herdr-Tab).
- **2026-11-02 09:49 · Alexander Sacharov** — review (frischer Agent): alle DoD-Punkte am Diff bestaetigt, keine Blocker. Drei notes in review-gaps (lokale Sperre vs. Merge-Rennen, 409 bei veraltetem Seitenstand, kein Merge-Test mit status). Geprueft: Tests gruen, dist reproduzierbar.
- **2026-11-02 10:06 · Alexander Sacharov** — explain: Der Review fand keinen Fehler, alle 7 Punkte der Definition of Done sind am Diff bestätigt und die Tests laufen grün. Offen sind drei kleine Lücken: Die Sperre eines geschlossenen Epics gilt nur lokal (Rennen zwischen zwei Klonen), der Web-Dialog zeigt bei veraltetem Seitenstand einen 409-Fehler statt der Frage, und es fehlt ein Test für epic.Merge mit status.
Warum: Keine der Lücken verliert Daten und sie treten selten auf; im schlimmsten Fall landet ein Ticket in einem Epic, das anderswo gerade geschlossen wurde.
1) Zurückgeben: Merge-Rennen schließen (Mitglied bei geschlossenem Epic nicht ohne force übernehmen), 409-Dialog mit Bestätigung anbieten, Merge-Test mit status ergänzen [give back]
2) Ticket so annehmen und die drei Lücken als bekannte Grenzen stehen lassen [accept]
Empfehlung: 2 - die Funktion ist vollständig und geprüft, die Lücken sind selten und ohne Datenverlust, und ein Nachbessern vergrößert den Diff eines schon abgenommenen Features.
- **2026-11-02 10:24 · Alexander Sacharov** — decision: 2) Ticket so annehmen und die drei Lücken als bekannte Grenzen stehen lassen
