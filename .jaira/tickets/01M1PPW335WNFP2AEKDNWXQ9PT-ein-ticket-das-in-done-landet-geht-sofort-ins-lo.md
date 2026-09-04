---
id: 01M1PPW335WNFP2AEKDNWXQ9PT
title: "Ein Ticket, das in done landet, geht sofort ins Logbuch"
status: signoff
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
updated-at: 2026-09-04T17:56:09Z
claimed-by: EE-3NX6GL3-34378
claimed-at: 2026-09-04T17:20:09Z
updated-by: BeMuCa
outcome-what: "Jam-Fix: FileLane ueberspringt Unfilebares (unlesbar, Stempel-Fehler, Kollision), benennt es als PartialError und filed den Rest - der Ankoemmling kommt immer durch, das Problem wird bei jeder Landung erneut gemeldet"
outcome-why: "Das Zweitmodell-Review reproduzierte live einen Dauerstau: ein Fehler mitten im oldest-first-Sweep liess das ankommende Ticket (als neustes zuletzt dran) ungefiled liegen, jede Folge-Landung scheiterte identisch"
outcome-resolves: "Core-Test deckt broken.md UND Kollision in einem Sweep (genau die zwei fileable gehen), CLI-Test beweist den Ankoemmling trotz Problem; go test ./... -race RC=0, Binary neu"
executed-by: fable
review-summary: "Runde 2 (Jam-Fund): continue-on-error sitzt IM FileLane statt in vier Aufrufern - eine Stelle, gleiche Philosophie wie List selbst ('one malformed ticket must not blank the whole board'); PartialError wiederverwendet statt neuem Fehlertyp. Der Cap-Pfad (TrimLane) behaelt bewusst den Abbruch: dort haengt die Auswahl des Aeltesten an vollstaendigem Wissen, bei der Doorway gibt es nichts zu waehlen."
review-gaps: "Nichts entfernt. Gelassen: TrimLane-Abbruchverhalten (begruendet verschieden); die Meldezeile heisst weiter 'trimming the lane failed' auch fuer Doorway-Probleme (ein Wortlaut, von Tests gepinnt - Umbenennen waere eigener Mini-Schnitt); Probe c des Reviewers (selectByID nach Filing) bleibt code-verifiziert."
review-check: |-
  1. Ein Ticket im signoff mit a akzeptieren: Meldung Accepted + filed-Zeile(n) mit restore-Datei; done bleibt leer.
  2. .jaira/logbook/<initialen>-<heute>/ enthaelt die Datei, Commits stehen drin (grep commits:).
  3. jaira restore <datei> holt es zurueck nach done; das naechste Landen fegt es erneut - kein Stau.
  4. Stau-Probe: eine Kauderwelsch-Datei nach .jaira/tickets/broken.md legen, ein Ticket nach done moven -> es wird GEFILED und stderr nennt broken.md; broken.md loeschen, Meldung verschwindet.
  5. jaira lanes show done -> logbook-on-entry; move --json -> filed_on_entry: true.
test-verdict: "pass: Gates gruen (go test ./... -race, Cache geleert, RC=0, 15 Pakete; Build + Binary), DoD am Baum verifiziert (Flag+Guard in lane.go, FileLane/StampCommits in trim.go, vier Schreibstellen via flow-Switch/settleLane, Tests je Pfad), Verhalten exerziert: CLI-Move filed 4/4 mit gestempelten Commits (Testlauf), Jam-Fall filed den Ankoemmling trotz broken.md, dieses Board lief live durch den Sweep (done leer, 22 Dateien im Logbuch)"
review-verdict: "accept (Zweitmodell Sonnet, zwei Durchgaenge). Runde 1 fand den Doorway-Jam live (Kollisions-Probe, zweifach reproduziert); nach dem Fix hat der Reviewer BEIDE Jam-Proben am selbst gebauten HEAD-Binary (b82f8d3) erneut gefahren: beide Ankoemmlinge werden gefiled, das Problem wird je Landung benannt, EXIT=0 - Ausgaben woertlich kopiert, Fix-Diff selbst gelesen (errors.As/PartialError, skip+continue statt Abbruch); alle vier Pakete -race RC=0. Offengelegt: die Unlesbar-Variante stuetzt sich auf den gepinnten Test (vom Reviewer -race bestaetigt), nicht auf eigene Live-Repro; sein Berichts-VOLLTEXT ging im Notification-Kanal wiederholt verloren - Substanz (Proben, Verdikt) kam an und steht hier."
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
- **2026-09-04 17:50 · BeMuCa** — F1 gefixt: FileLane ueberspringt und benennt (PartialError) statt abzubrechen - das ankommende Ticket kommt immer durch, Problemfaelle werden bei jeder Landung erneut gemeldet. Gepinnt: core (unlesbares Board-File + Logbuch-Kollision in EINEM Sweep: genau die zwei fileable gehen, der Rest bleibt benannt liegen) und CLI (TestDoorwayKeepsFilingPastAProblem: Ankoemmling gefiled trotz broken.md). Reviewer-Restnotizen: Probe c (selectByID nach Filing) nur code-verifiziert - Fallback-Logik model.go greift; restore->resweep am selben Tag live gruen (restore MOVED die Datei zurueck, keine Kollision).
