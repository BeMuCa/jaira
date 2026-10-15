---
id: 01M2K275PSXDASRN6Y1VN71NVG
title: "Ein fertiger Milestone geht ins Archiv, und ein falscher laesst sich loeschen"
status: pre-process
ready: true
creator: Alexander Sacharov
assignee: BeMuCa
goal: "Wer eine Runde Arbeit abgeschlossen hat, nimmt ihren Milestone samt seinen fertigen Tickets in einem Griff vom Board, und wer einen Milestone falsch angelegt hat, loescht ihn so, dass er auf keinem Klon zurueckkommt."
context: |-
  Milestones entstehen in 0YGWXQ (.jaira/milestones/<name>.md plus refs/jaira/milestones/<name>). Es gibt keinen Weg, einen wieder wegzubekommen.

  Was heute fehlt, beides:
  - ARCHIVIEREN. Ein Milestone ist fertig, wenn seine Tickets fertig sind. Danach steht er weiter im M-Picker und faerbt weiter Kartenkanten. Nach zehn Runden Arbeit sind zehn tote Milestones im Picker und der Picker ist unbrauchbar - genau das Problem, das jaira logbook fuer Tickets schon geloest hat (.jaira/logbook/<wer>-<datum>/, 'jaira restore' bringt zurueck).
  - LOESCHEN. 'jaira milestone delete' existiert nicht. Die critique-Lane von 0YGWXQ hat den Loeschweg in Runde 1 als toten Code entfernt (RecordMilestoneDelete, PendingMilestone, gitref.DeleteMilestone hatten keinen Aufrufer) - richtig entfernt, aber der Befehl, der sie haette aufrufen sollen, wurde nie geschrieben.

  Ausgeloest am 2026-09-15: Alex hat einen Milestone 'test' angelegt und wollte ihn wieder weghaben. Es ging nur von Hand: rm der Datei, 'git update-ref -d refs/jaira/milestones/test' und ein Push mit Doppelpunkt gegen upstream. Drei Schritte, von denen zwei git-Wissen brauchen, fuer etwas, das ein Befehl sein muss.

  Der Ref ist der Teil, den man dabei vergisst. Loescht jemand nur die Datei, kommt der Milestone beim naechsten 'jaira fetch' vom Ref zurueck.

  Vorbild ist da und soll benutzt werden, nicht neu erfunden: jaira logbook fuer Tickets. Offen und Sache der Plan-Lane: ob das Archiv derselbe Ordner ist oder ein eigener, und was passiert, wenn ein Milestone archiviert wird, dessen Tickets noch offen sind.
definition-of-done: "Ein fertiger Milestone laesst sich in einem Griff archivieren: Datei und Ref verschwinden vom Board, der Inhalt bleibt lesbar - nach dem Vorbild von jaira logbook"
tags:
  - cli
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-09-15T17:35:45Z
updated-at: 2026-10-15T12:38:25Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-3783128
claimed-at: 2026-10-15T12:19:23Z
---

# Ein fertiger Milestone geht ins Archiv, und ein falscher laesst sich loeschen

## Definition of Done

- [ ] 'jaira logbook <milestone>' nimmt einen fertigen Milestone in einem Griff vom Board: erst jedes seiner Tickets, das in diesem Klon als Datei auf dem Board liegt (Commits gestempelt wie bei 'jaira logbook <id>'), dann den Milestone selbst - alles in denselben Ordner .jaira/logbook/<wer>-<datum>/ (der Milestone unter milestones/). Danach steht keins seiner Tickets mehr auf dem Board ausser einem, dessen Datei nur auf seinem Ref liegt (es liegt auf einem fremden Branch); die Ausgabe nennt es. --json nennt die abgelegten Tickets
- [x] Ein abgelegter Milestone faerbt keine Kartenkante mehr, steht nicht im M-Picker und nicht auf der Milestone Roadmap
  proof: core/milestone/milestone.go:359 LoadAll laesst 'status: filed' aus; internal/tui/model.go:404-409 baut Kartenkanten/M-Picker (mstones) und Roadmap (milestones, internal/tui/layer.go:149) nur daraus; TestLoadAllLeavesOutAFiledMilestone (core/milestone/milestone_test.go:306), TestFilingAMilestoneTakesItOffTheBoardAndRestoreBringsItBack (internal/cli/milestones_test.go:333)
- [x] Es gibt einen Weg zurueck wie fuer Tickets: 'jaira restore <name>-<ID>.md' holt den abgelegten Milestone mit Ticketliste und Farbe zurueck, und jedes mit abgelegte Ticket kommt mit 'jaira restore <ticketdatei>' zurueck
  proof: internal/cli/archive.go:160-182 unfileMilestone; TestFilingAMilestoneTakesItOffTheBoardAndRestoreBringsItBack (internal/cli/milestones_test.go:333, Restore ab :407); Ticket: core/ticket/store.go:634 Restore, TestRestoreFindsLoggedTicket (core/ticket/logbook_test.go:108)
- [ ] Fuehrt der Milestone noch ein Ticket, das die Terminal-Lane nicht erreicht hat, legt 'jaira logbook <milestone>' gar nichts ab - weder ein Ticket noch den Milestone -, nennt jedes solche Ticket mit seiner Lane und endet mit Exit 3 (milestone_unfinished)
- [ ] 'jaira milestone delete <milestone>' loescht einen falschen Milestone fuer alle Klone: der Ref wird NICHT entfernt, er bekommt eine Marke ('status: filed' plus 'deleted: <datum> <wer>'), und die Datei hier ist weg. Ein zweiter Klon, der die Datei hat - als <name>-<ID>.md oder im alten Format <name>.md -, verliert sie beim naechsten 'jaira fetch' und schiebt sie nie wieder hoch, auch nicht aus einem noch wartenden Schreibvorgang. Belegt durch einen Zwei-Klon-Test fuer beide Dateiformen
- [ ] core/release/notes/N71NVG.md hat je eine '- '-Zeile (eine Zeile, nie umbrochen) fuer jede Aenderung, die ein Benutzer merkt: logbook <milestone> legt die Tickets mit ab, 'jaira milestone delete', und 'jaira fetch' meldet geloeschte Milestones
- [ ] Loeschen laesst die Tickets des Milestones unberuehrt (die Ausgabe nennt, wie viele ihre Gruppe verlieren), haengt Nachfolger auf der Roadmap an seine Vorgaenger um wie 'milestone set --no-roadmap', gibt den Namen frei und verweigert einen abgelegten Milestone (Exit 3)

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-13 12:02 · BeMuCa** — Ueberschneidung mit E70WYN (Milestone umbenennen, 2026-10-13): E70WYN nimmt einen alten Namen von anderen Klonen, indem sein Ref STEHEN bleibt und markiert wird - 'status: filed' plus neue Zeile 'renamed-to: <neu>' - derselbe Weg wie jaira logbook. Es wurde KEIN Ref-Loeschen gebaut; die delete-Frage hier bleibt offen.
Fuer delete wichtig: ein geloeschter Ref erreicht keinen Klon, der die Datei schon hat (refsync IncomingMilestones loescht nie). Eine Marke erreicht ihn.
Neu durch E70WYN: Milestone.HoldsName() - eine Marke mit renamed-to haelt den Namen NICHT (nichts zum Zurueckholen), eine reine filed-Marke weiter schon. Falls delete eine Marke nutzt, ist das die Stelle fuer 'Name wieder frei'.
- **2026-10-15 12:19 · BeMuCa** — Berk, 15.10., in seinen Worten: "Wir sollten auch umsetzen, dass die Milestones als Ganzes vom Board kommen, damit wir, wenn die fertig sind, keine Reste da uebrig haben." Gemeint: ist ein Milestone fertig, nimmt EIN Befehl den Milestone und seine fertigen Tickets zusammen vom Board - kein Milestone mehr im Picker oder auf der Roadmap, keins seiner fertigen Tickets mehr auf dem Board. Unfertige Tickets werden nicht abgelegt: jaira sagt welche und laesst sie stehen. Dazu, Lehre vom 14.10.: wir haben die Refs von sechs Test-Milestones auf origin geloescht; Alex' Klon hatte die Dateien noch, sein jaira hat sie um 20:20 neu migriert und als neue ID-Refs wieder hochgeschoben. Ein delete, das nur den Ref wegnimmt, macht jeder Klon mit der Datei rueckgaengig - das Loeschen braucht eine Marke, die jeder Klon befolgt. Ticket von Alex auf BeMuCa uebernommen, weil Berk es fuer das naechste Release will.
- **2026-10-15 12:38 · BeMuCa** — pre-process, Teil 1 - Milestone samt Tickets ablegen. Entscheidung und warum:
- Kein neuer Befehl. 'jaira logbook <milestone>' (internal/cli/logbook.go logbookMilestone ~346) legt kuenftig auch die Tickets ab. Der Befehl heisst schon 'vom Board nehmen'; ein zweiter Befehl wuerde zwei Wege fuer dasselbe schaffen.
- Die Sperre bleibt wie heute: ein Mitglied auf dem Board, das die Terminal-Lane nicht erreicht hat -> Exit 3 milestone_unfinished, alle genannt, NICHTS abgelegt. Ablegen trotz offener Tickets wuerde Tickets ohne Gruppe zuruecklassen - genau die Reste, die Berk nicht will.
- Reihenfolge: erst die Tickets (StampCommits + s.Logbook, wie logbookOut), dann der Milestone. Bricht ein Lauf mittendrin ab, steht der Milestone noch da; der zweite Lauf zaehlt schon abgelegte Tickets als fertig (s.Load findet sie nicht mehr, Extra blendet FiledAwayIDs aus: core/refsync/refsync.go:1500ff).
- Erster Fehler bricht ab (kein Weiterlaufen wie FileLane). Grund: einfacher, und der Neustart ist sicher.
- Ein Mitglied, das nur auf seinem Ref liegt (t.ReadOnly), wird NICHT abgelegt und nicht als Sperre gezaehlt. Das ist meist ein Ticket, das auf einem PR-Branch schon abgelegt ist und dessen Ref erst beim Merge verschwindet. Es ist nicht Sache dieses Klons (wie filesInLane in core/ticket/trim.go:27). Die Ausgabe nennt es (--json: ref_only).
- Ein Ticket, das AUCH in einem zweiten, offenen Milestone steht: wird trotzdem abgelegt (es ist fertig). Die Zeile im zweiten Milestone bleibt stehen, dessen Sperre zaehlt es als fertig (logbook.go: s.Load-Fehler -> continue). Das ist heute schon so bei 'jaira logbook <id>'. Kein Extra-Code.
- Gleicher Ordner fuer alles: logbookFolder() einmal vor der Schleife ausrechnen.
- --json bekommt 'tickets' (gleiche Form wie 'filed' bei logbook --all, trimmedJSON in internal/cli/flow.go:857) und 'ref_only' (Handles). 'count' bleibt die Zahl der Mitglieder.
- Bestehender Test TestFilingAMilestoneTakesItOffTheBoardAndRestoreBringsItBack (internal/cli/milestones_test.go:333) bricht dann: sein Farb-Check verlangt, dass das done-Ticket noch auf dem Board steht. Umbauen: Ticket mit 'jaira restore <ticketdatei>' zurueckholen und dann pruefen, dass es keine Milestone-Farbe traegt.
