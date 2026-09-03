---
id: 01M0ZCFA03ESJ6RJVH8EQA3GN1
title: jaira init schreibt Lane-Dateien ohne Reihenfolge
status: done
ready: true
creator: BeMuCa
goal: "Ein mit jaira init angelegtes Board zeigt seine Lanes in der Reihenfolge, die das Default-Board nennt"
context: |-
  Gefunden am 26.08. beim Verifizieren von 4DQPMS, nicht vom Ticket selbst behauptet.

  Was falsch ist: 'jaira init' mit einer eigenen Default-Board-Auswahl schreibt die Lane-Dateien, aber keine 'order'-Datei. Beim naechsten Laden ist keine Lane mehr Builtin, also muss core/lane/order() die Reihenfolge aus den 'after:'-Feldern rekonstruieren - und das kann es nicht eindeutig.

  Warum nicht: zwei Lanes koennen denselben Anker haben. brainstorm und todo sind beide 'after: backlog'. order() setzt eine Lane direkt hinter ihren Anker, also schiebt die spaeter eingefuegte die frueher eingefuegte nach rechts. Gemessen in einem Test: erwartet [backlog brainstorm todo pre-process in-progress human review signoff done blocked], geladen [backlog todo pre-process in-progress human review signoff done blocked brainstorm] - brainstorm landet am Ende.

  Was NICHT die Ursache ist: 'direkt hinter den Anker' ist Absicht. Genau so landet critique zwischen in-progress und human, vor dem eingebauten Nachfolger, der denselben Anker hat. An order() darf man dafuer nicht ruehren.

  Wo es zu reparieren ist: core/lane/defaultboard.go:154 Materialise() schreibt b.Lanes als Dateien und ruft SaveOrder nie. lane.Add, lane.Remove und lane.MoveLane (core/lane/order.go:286,353,383) rufen es. b.Lanes ist schon eine geordnete Liste, die Information ist also vorhanden und wird nur weggeworfen. Aufrufstelle: internal/cli/tickets.go:55.

  Warum es hier nicht auffaellt: dieses Board hat eine order-Datei, weil es durch 'lanes remove' entstanden ist, und das ruft SaveOrder.
definition-of-done: "Nach 'jaira init' mit einer Auswahl, die nicht dem Builtin-Board entspricht, laedt das Board dieselbe Lane-Reihenfolge, die das Default-Board nennt - abgesichert durch einen Test, der init-Materialisierung und erneutes Laden vergleicht"
blocked-by:
  - BNZERQ
commits: []
created-at: 2026-08-26T15:54:27Z
updated-at: 2026-09-03T10:28:55Z
claimed-by: EE-3NX6GL3-2116865
claimed-at: 2026-08-26T16:00:35Z
updated-by: BeMuCa
assignee: BeMuCa
question: |-
  Wie soll 'jaira init' eine verkleinerte Lane-Auswahl umsetzen? Zwei verteidigbare Wege, beide Fehler haengen daran.

  A) Tombstone statt Kopien. Materialise schreibt fuer die NICHT gewaehlten Builtins eine 'removed'-Datei und legt fuer gewaehlte, unveraenderte Builtins keine Datei an. Dann bleiben sie Builtin, behalten ihre gelieferte Reihenfolge, und die Umsortierung verschwindet von selbst, ohne dass order() angefasst wird. Passt zu 'order is precedence's business' und zu 'ein Repo, dessen Besitzer nichts geaendert hat, traegt keine Lane-Dateien'. Kostet: ein neu angelegtes Board sieht anders aus als heute, eine removed-Datei statt zehn Kopien. Bestehende Boards laufen unveraendert weiter, ihre Dateien sind autoritativ.

  B) Reihenfolge mitschreiben. Materialise ruft SaveOrder mit der Auswahl, wie lane.Add, lane.Remove und lane.MoveLane es tun. Kleiner Eingriff, repariert die Reihenfolge - aber nicht, dass nicht gewaehlte Builtins zurueckkommen; dafuer braucht es trotzdem den Tombstone. Und es macht das Default-Board zur Instanz fuer Reihenfolge, was defaultboard.go:117 ausschliesst.

  Empfehlung A: repariert beide Symptome mit einer Aenderung und keines gegen einen Satz, der im Code steht. B repariert die Haelfte und muss den Satz umschreiben.

  Nicht zur Wahl: order() so aendern, dass Lanes mit gleichem Anker ihre Ladereihenfolge behalten. Das wuerde critique hinter human schieben statt davor, also dein echtes Board umrouten.
outcome-what: "Kein Code geaendert. Diese Runde hat gemessen statt reparariert: Materialise mit vier gewaehlten Lanes, danach Load - es kommen zehn Lanes zurueck, [pre-process in-progress human review signoff done blocked backlog todo brainstorm]. Testentwuerfe, die den im Ticket vermuteten Fix vorausgesetzt haben, sind verworfen."
outcome-why: "Der im Ticket vorgeschlagene Fix (SaveOrder in Materialise) widerspricht core/lane/defaultboard.go:117, wo steht dass Reihenfolge ausdruecklich nicht Sache des Default-Boards ist. Und er wuerde nur die Haelfte reparieren: dass nicht gewaehlte Builtins zurueckkommen, liegt an der fehlenden removed-Tombstone, nicht an der Reihenfolge."
outcome-resolves: "Die DoD ist nicht erfuellt und soll es hier auch nicht sein. Sie setzt einen Weg voraus, der sich beim Nachmessen als der schlechtere von zwei erwiesen hat. Welcher Weg gilt, ist eine Entscheidung - deshalb liegt die Frage jetzt bei einem Menschen und nicht ein halber Fix im Repo."
blocked-reason: "Symptom von BNZERQ; wird mit dem Redesign erledigt, nicht einzeln"
review-summary: "Kein eigener Diff - der Bug ist durch BNZERQ (743737f) mitgeloest: init materialisiert das Board als Dateien inklusive order"
review-gaps: none
review-verdict: "durch BNZERQ subsumiert, per Code-Trace, Tests und Live-Check bestaetigt"
review-check: "In einem leeren Verzeichnis: git init; jaira init mit einer Auswahl, die NICHT dem Builtin-Board entspricht; dann cat .jaira/lanes/order - die Datei existiert und nennt die gewaehlte Reihenfolge. jaira (TUI) zeigt dieselbe Reihenfolge."
---

# jaira init schreibt Lane-Dateien ohne Reihenfolge

## Definition of Done

- [ ] Nach 'jaira init' mit einer Auswahl, die nicht dem Builtin-Board entspricht, laedt das Board dieselbe Lane-Reihenfolge, die das Default-Board nennt - abgesichert durch einen Test, der init-Materialisierung und erneutes Laden vergleicht

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-08-26 16:02 · BeMuCa** — Gemessen, bevor irgendwas geaendert wurde - und es steht schlechter als das Ticket behauptet.

Zweite, groessere Tatsache: 'jaira init' mit einer verkleinerten Auswahl verkleinert das Board gar nicht. Materialise(root, set, {Lanes: [backlog brainstorm todo done]}) schreibt vier Dateien; das darauffolgende Load(root) liefert ZEHN Lanes: [pre-process in-progress human review signoff done blocked backlog todo brainstorm]. Grund: core/lane/lane.go:398 spritzt die Builtins immer ein, und nur die 'removed'-Tombstone-Datei (lane.go:470) schliesst eine Lane aus. Materialise schreibt die nie - SaveRemoved wird nur von lane.Remove aufgerufen. Der Doc-Kommentar von DefaultBoard sagt aber 'which lanes a freshly initialised board gets'.

Und der im Ticket vorgeschlagene Fix widerspricht dem Code: core/lane/defaultboard.go:117 sagt woertlich 'reordering alone must not trigger materialisation, since order is precedence's business, not the default board's'. Materialise soll also gerade NICHT ueber Reihenfolge entscheiden. SaveOrder dort einzubauen kaempft gegen diesen Satz.

Deshalb liegt das Ticket jetzt bei einer Entscheidung und nicht bei einer Reparatur. Testentwuerfe, die die Ticket-Antwort schon vorausgesetzt haben, sind wieder verworfen - nichts committed.
- **2026-08-27 16:37 · BeMuCa** — Antwort von Berk am 27.08.: weder A noch B. Das Board soll sein Lane-Verzeichnis sein, nichts wird eingespritzt, 'removed' faellt weg. Damit ist dieses Ticket ein Symptom des Redesign-Tickets 'Das Board ist sein Lane-Verzeichnis' und wird mit ihm erledigt: wenn init die Auswahl als Dateien plus order schreibt und Load nichts dazulegt, gibt es weder die Umsortierung noch die zurueckkehrenden Builtins.
- **2026-08-31 16:29 · BeMuCa** — Verifiziert am 31.08. (read-only, Wegwerf-Verzeichnis mit isoliertem JAIRA_HOME): jaira init schreibt die order-Datei jetzt selbst - BNZERQ (743737f) hat das subsumiert. Code-Pfad: internal/cli/tickets.go:22-50 -> core/lane/lane.go:403-482 -> defaultboard.go:117-153 -> order.go. Tests: TestFirstLoadUsesTheDefaultBoard, TestInitWithNoDefaultBoardWritesTheBuiltinsAsFiles (init_test.go/lane_test.go, Bodies gelesen, nicht nur Namen). Live: init mit eigener Auswahl -> .jaira/lanes/order existiert und traegt die gewaehlte Reihenfolge. Kein eigener Code noetig; Commits liegen auf BNZERQ.
