---
id: 01M4WP3E1FWG8JR3KVRRG2P50F
title: "Ein Formatwechsel erreicht alte jaira-Versionen und Agenten, bevor er Daten kostet"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Aendert ein Release das Dateiformat des Boards, schreibt eine aeltere jaira-Version nicht mehr still im alten Format weiter, und ein Agent erfaehrt davon und legt Board-Dateien im neuen Format an"
context: |-
  Was heute falsch ist: Ein Release, das das Board-Format aendert, erreicht nur, wer das neue Binary hat. Aeltere jaira-Versionen im Team schreiben im alten Format weiter, und nichts sagt es ihnen.
  Anlass: EBTX6X. Milestone-Dateien heissen dort <name>-<ID>.md statt <name>.md, der Ref refs/jaira/milestones/<ID> statt .../<name>. Die Review fand am 14.10.: Aenderungen einer 0.4.x-Version an einer schon umgestellten Milestone gingen still verloren.
  EBTX6X uebernimmt solche Aenderungen jetzt nachtraeglich (Variante C). Das repariert im Nachhinein und beugt nicht vor.
  Berks Grundsatz, 14.10.: 'ein Update von uns sollte nie dafuer sorgen dass bestehende Sachen verschwinden, man sollte immer die Sachen einfach an das neue anpassen'.
  Luecke 1, alte Builds: Es gibt keine ueber git geteilte Format-Markierung. release.Stamped (internal/cli/update.go:50) merkt sich nur pro Board-Statusordner, welche Version zuletzt lief, und ein aelteres Binary prueft gar nichts.
  Luecke 2, Agenten: Milestone-Dateien sind von Hand editierbar ('the file IS the interface', core/refsync/refsync.go Kommentar zu IncomingMilestones). Ein Agent mit altem Skill oder altem CLAUDE.md-Block legt von Hand <name>.md an. Wichtige Release-Notes erreichen ihn nicht.
  Grenze: Das hilft erst ab dem naechsten Formatwechsel. 0.4.x kann eine Markierung nicht mehr nachtraeglich lernen.
  Verwandt: T8A8KM (Popup und Markierung wichtiger Release-Notes, DoD 3-5), PWAF22 (Upgrade zieht Skills und update nach).
definition-of-done: "Das Board traegt eine ueber git geteilte Format-Version (Ort im Ticket begruendet); ein Build, der eine neuere Format-Version sieht als er kennt, verweigert jedes Schreiben mit exit 3 und nennt 'jaira self upgrade'; Lesen bleibt erlaubt; Test"
tags:
  - release
blocked-by: []
related:
  - 01M21PJTWPK4DZ89PRY6T8A8KM
  - 01M2HX6SYCFB27V26R5APWAF22
follows: 01M4VCBB5FCBME5SPR6REBTX6X
commits: []
created-at: 2026-10-14T07:47:47Z
updated-at: 2026-10-14T07:47:47Z
---

# Ein Formatwechsel erreicht alte jaira-Versionen und Agenten, bevor er Daten kostet

## Definition of Done

- [ ] Das Board traegt eine ueber git geteilte Format-Version (Ort im Ticket begruendet); ein Build, der eine neuere Format-Version sieht als er kennt, verweigert jedes Schreiben mit exit 3 und nennt 'jaira self upgrade'; Lesen bleibt erlaubt; Test
- [ ] Ein Release, das die Format-Version erhoeht, stellt beim ersten Lauf alle betroffenen Daten des Boards um und setzt die Version erst danach; Test an einem Board im alten Format
- [ ] Eine als wichtig markierte Release-Note (Markierung aus T8A8KM DoD 4) erreicht einen Agenten: jaira nennt sie in einer Ausgabe, die ein Agent liest (z.B. resume/next --json), bis 'jaira update' gelaufen ist; Test
- [ ] Agent-Block und jaira-Skill sagen: nach einer gemeldeten Formataenderung Board-Dateien nur im neuen Format anlegen und den Nutzer darauf hinweisen
- [ ] Eine Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

