---
id: 01M1PPW335WNFP2AEKDNWXQ9PT
title: "Ein Ticket, das in done landet, geht sofort ins Logbuch"
status: in-progress
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Der Move/Accept nach done stempelt die Commits und legt das Ticket direkt im Logbuch ab - gemeldet mit restore-Weg; die done-Spalte ist ein Durchgang, das Logbuch die vollstaendige Chronik des Fertigen"
context: "Berk am 04.09.: 'Alle Tickets die auf Done landen gehen direkt ins Logbuch, nicht erst wenn sie aus done rausgeworfen werden. wir wollen im logbook ja alles speichern was gemacht wurde direkt.' REVIDIERT bewusst die SGPDYK-Abnahme von gestern (10 bleiben sichtbar) - der Cap-Mechanismus (holds) bleibt als Feature fuer andere Lanes/Boards bestehen, builtin done tauscht holds: 10 gegen das neue Flag. Umsetzung: Lane-Flag logbook-on-entry (nur terminal, parse-Guard wie holds); an ALLEN VIER Status-Schreibstellen (cli flow, tui applyMove/forceMove/accept) nach gelungenem Move: Commits stempeln (Ableitung wie 'jaira logbook', gate hat sie schon gefordert) und Lane komplett ins Logbuch fegen (auch Altbestand - selbstmigrierend); nie still. Startbildschirm zeigt Fertiges pro Tag aus dem Logbuch (SPDWGH), die Sicht geht nicht verloren."
definition-of-done: "Move und Accept nach done melden die Logbuch-Ablage mit restore-Datei; die abgelegte Datei traegt gestempelte Commits; done ist danach leer (Altbestand mitgenommen); logbook-on-entry auf nicht-terminaler Lane wird beim Parsen verweigert; Tests decken CLI-Move, TUI-Accept und den Parse-Guard; go test ./... -race gruen"
tags: []
blocked-by: []
commits: []
created-at: 2026-09-04T17:18:43Z
updated-at: 2026-09-04T17:43:36Z
claimed-by: EE-3NX6GL3-34378
claimed-at: 2026-09-04T17:20:09Z
updated-by: BeMuCa
outcome-what: "Lane-Flag logbook-on-entry (nur terminal): der Move/Accept nach done stempelt die Commits (env.DeriveCommits) und fegt die ganze Lane ins Logbuch - gemeldet je Ticket mit restore-Weg, --json traegt filed_on_entry+trimmed; builtin done nutzt das Flag statt holds:10; StampCommits/FileLane im core, CLI-stampCommits delegiert"
outcome-why: "Berk am 04.09.: alle Tickets, die auf done landen, gehoeren sofort ins Logbuch - das Logbuch ist die Chronik, done nur der Durchgang; revidiert bewusst die gestrige 10-bleiben-Abnahme (im Kontext dokumentiert)"
outcome-resolves: "CLI-Move, TUI-Accept und settleLane je mit Test (Commits-Erhalt am abgelegten File verifiziert, Altbestand-Sweep mit drin), Parse-Guard gegen nicht-terminale Doorways; go test ./... -race RC=0, 15 Pakete; live: dieses Board ist umgestellt, done leer, 10 Restbewohner gefegt"
executed-by: fable
review-summary: "Kritik: richtiger Ort und kein neues Muster - FileLane/StampCommits sitzen neben TrimLane/Logbook im core, die vier Schreibstellen laufen durch die zwei existierenden Sammelpunkte (flow-Switch, settleLane); der Entry-Zweig ist ein Switch-Arm neben dem Cap, keine Parallelmechanik. Ein Preis benannt: done ist auf dem Board immer leer - 'zuletzt fertig' liest man im Startbildschirm (Fertiges pro Tag, SPDWGH) und im Logbuch, nicht mehr in der Spalte. Berks explizite Richtungsentscheidung."
review-gaps: "Nichts entfernt. Gelassen: holds als Feature samt Tests (via Fixture-Patch statt builtin); die doppelte Meldezeile (filed vs. left-for) statt einer generischen - zwei verschiedene Ereignisse sollen verschieden klingen; kein Stempeln in TrimLane (Cap-Tickets behalten Gate-gepruefte Commits)."
review-check: |-
  1. Ein Ticket im signoff mit a akzeptieren: Meldung Accepted + je eine filed-Zeile mit restore-Datei; done bleibt leer.
  2. .jaira/logbook/<initialen>-<datum>/ enthaelt die Datei, Commits stehen drin.
  3. jaira restore <datei> bringt es zurueck (landet wieder in done -> beim naechsten Landen wird es erneut gefegt).
  4. jaira lanes show done -> logbook-on-entry statt Holds.
  5. move --json -> filed_on_entry: true neben trimmed.
---

# Ein Ticket, das in done landet, geht sofort ins Logbuch

## Definition of Done

- [x] Move und Accept nach done melden die Logbuch-Ablage mit restore-Datei; die abgelegte Datei traegt gestempelte Commits; done ist danach leer (Altbestand mitgenommen); logbook-on-entry auf nicht-terminaler Lane wird beim Parsen verweigert; Tests decken CLI-Move, TUI-Accept und den Parse-Guard; go test ./... -race gruen
  proof: core/lane logbook-on-entry (terminal-Guard, Template, builtin done statt holds); core/ticket FileLane+StampCommits; alle vier Schreibstellen via flow.go-Switch bzw. tui settleLane (applyMove/forceMove/accept); Tests: TestMoveIntoDoneFilesEverythingToTheLogbook (Commits gestempelt verifiziert), TestAcceptFilesTheTicketIntoTheLogbook, TestSettleLaneFilesTheDoorwayLane, Parse-Guard; go test ./... -race RC=0

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-04 17:30 · BeMuCa** — Entscheidungen: (1) holds bleibt als Feature (andere Lanes/Boards), builtin done tauscht es gegen logbook-on-entry - SGPDYK-Mechanik nicht rueckgebaut, nur nicht mehr von done genutzt (Berks Revision vom 04.09. dokumentiert im Kontext). (2) FileLane fegt die GANZE Lane (Altbestand mit) - selbstmigrierend, kein Sonderpfad. (3) Commit-Stempeln vor dem Ablegen an allen vier Schreibstellen via env.DeriveCommits; stampCommits (CLI) delegiert jetzt an core StampCommits (drei Nutzer: logbook-Cmd, archive, Entry-Filing beidseitig). (4) TUI-Cap-Branch-Test entfaellt - der Zweig ist ein Switch-Arm auf dem core-/CLI-getesteten TrimLane; CLI-Cap-Tests patchen die Lane-Kopie ihres Fixture-Boards auf holds:10 (builtin ist jetzt Doorway). (5) Board live: done.md-Kopie refreshed, die 10 Restbewohner per jaira logbook gefegt - done ist leer und bleibt Durchgang.
- **2026-09-04 17:43 · BeMuCa** — Review-Fund F1 (live reproduziert, zweimal): die Doorway verklemmt - schlaegt der Sweep fehl (z.B. PartialError durch ein unlesbares Ticket, oder eine Logbuch-Kollision beim aeltesten Bewohner), bleibt AUCH das ankommende Ticket ungefiled in done liegen, und jede weitere Landung scheitert am selben Fehler: kein Selbstheilen. Ursache: FileLane bricht beim ersten Fehler ab (List-Fehler ganz, Einzel-Fehler mitten in der oldest-first-Reihe - das ankommende ist das NEUSTE und damit als letztes dran). Fix: continue-on-error - unlesbare/kollidierende Tickets ueberspringen und benennen (PartialError), alles andere filen; das ankommende kommt damit immer durch, der Problemfall wird bei jeder Landung erneut gemeldet statt alles zu blockieren. Gleiche Philosophie wie List selbst ('one malformed ticket must not blank the whole board').
