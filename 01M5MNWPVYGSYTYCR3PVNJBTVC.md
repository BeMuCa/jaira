---
id: 01M5MNWPVYGSYTYCR3PVNJBTVC
title: "Ein offener Plan-Punkt sperrt die Ziellane nicht mehr, wenn die Definition of Done erfuellt ist"
status: in-progress
ready: true
creator: BeMuCa
goal: "Ob ein Ticket in die Ziellane darf, entscheidet allein die Definition of Done; ein offener Plan-Schritt sperrt nicht mehr, sondern wird beim gelungenen Move als Hinweis genannt - der Plan ist die Methode, die DoD das Ergebnis."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkt 8): GH0GKD mit DoD 9/9 wurde an done abgelehnt, weil Plan-Schritt 30 offen war.
  Nachgeprueft am 23.10.: core/gate/gate.go:411-424 verweigert jede Ziellane mit CodePlanIncomplete, solange ein Plan-Punkt offen ist. Der Kommentar dort sagt, Plan und DoD seien nicht unabhaengig.
  CLAUDE.md dieses Repos sagt das Gegenteil: 'Ticking the plan does not tick the definition of done' - der Plan ist eine zweite, separate Liste.
  Nicht entschieden: Gate weg (nur Warnung), oder der Worker hakt beim Abschluss den letzten Plan-Schritt selbst ab. Brainstorm entscheidet; die einfachere Variante ist das Gate zu einer Warnung zu machen.
definition-of-done: Ein Move in eine Ziellane mit vollstaendiger DoD und offenem Plan-Punkt wird nicht mehr abgelehnt; die Ausgabe nennt den offenen Plan-Punkt als Hinweis; ein Test deckt den Fall
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:25:53Z
updated-at: 2026-10-26T19:28:03Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-264760
claimed-at: 2026-10-26T19:12:35Z
---

# Ein offener Plan-Punkt sperrt die Ziellane nicht mehr, wenn die Definition of Done erfuellt ist

- [x] Gate: Plan-Sperre in gate.go entfernen, Hinweis-Helfer OpenPlanHints
  proof: core/move/move_test.go TestMoveIntoDoneWithOpenPlanStepCarriesAHint; core/gate/plan_test.go TestUnfinishedPlanDoesNotBlockTerminalLane
- [ ] move.Result.Hints fuellen bei Move in Ziellane
- [x] CLI Text+JSON+dry-run zeigen den Hinweis
  proof: core/gate/plan_test.go TestUnfinishedPlanAddsNothingToDoDRefusal
- [ ] TUI: model.go zwei Pfade + signoff Accept zeigen den Hinweis
- [ ] Tests: plan_test.go drehen, move-Test fuer Hinweis
- [ ] Docs (SKILL, README, AGENTS, 50-done) + Release-Notiz
- [ ] gofmt/vet/build + Paket-Tests, DoD abhaken

## Definition of Done

- [ ] Ein Move in eine Ziellane mit vollstaendiger DoD und offenem Plan-Punkt wird nicht mehr abgelehnt; die Ausgabe nennt den offenen Plan-Punkt als Hinweis; ein Test deckt den Fall
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/
- [ ] Bei unvollstaendiger DoD nennt die Ablehnung nur die offenen DoD-Punkte, keinen Plan-Schritt mehr
- [ ] Der Hinweis erscheint in der Textausgabe, im JSON und beim Move aus der TUI, das Annehmen in signoff eingeschlossen
- [ ] Die Texte, die das Plan-Gate beschreiben (.claude/skills/jaira/SKILL.md, README.md, docs/AGENTS.md, core/lane/builtin/50-done.md), sagen nichts Falsches mehr

## Options

- [x] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-26 19:17 · BeMuCa** — BRAINSTORM 1/3 - Befund im Code.
- core/gate/gate.go:417-431: jede Ziellane (target.Terminal) verweigert mit CodePlanIncomplete, eine Meldung pro offenem Plan-Schritt. Die DoD spielt dabei keine Rolle.
- Die DoD-Sperre steht davor, gate.go:358-385, haengt aber an target.RequiresNonModelSignal, nicht an Terminal. Das ist per Default dasselbe (core/lane/lane.go:344), eine Lane kann es abwaehlen.
- Folge: DoD 9/9 + ein offener Plan-Schritt = Ablehnung. So bei GH0GKD (requirementsgenie): offen war Schritt 30 "R6-13 Abnahme-Vorbereitung", die Arbeit war laut Notiz 23.10. 07:56 erledigt, nur nicht abgehakt.
- Der TUI-Accept in signoff (internal/tui/signoff.go:177-183) hat kein Override: die Person kommt nur weiter, indem sie den Plan-Schritt abhakt.
- Der Kommentar gate.go:404-416 begruendet die Sperre mit "erst Design fertig, nicht gebaut". Aber Plan-Haken und DoD-Haken setzt derselbe Worker selbst. Der Plan-Haken ist kein unabhaengiger Beleg, er faengt nur das Vergessen eines zweiten Hakens.
- Die Doku ist gespalten. Agent-Block core/board/announce.go:98-104 (-> CLAUDE.md/AGENTS.md): die DoD ist das, was das Gate der Ziellane liest; der Plan ist eine zweite, getrennte Liste. Dagegen sagen die Plan-Sperre ausdruecklich an: .claude/skills/jaira/SKILL.md:105-107, README.md:247-249, docs/AGENTS.md:227-228, core/lane/builtin/50-done.md:11.
- Einen Kanal fuer "Move gelungen, aber mit Hinweis" gibt es heute nicht: gate.Violation ist nur Ablehnung (gate.go:41-51), move.Result hat Refused/Overrode/Trimmed/SettleErr (core/move/move.go:58-70).
- **2026-10-26 19:18 · BeMuCa** — BRAINSTORM 2/3 - Optionen und Entscheidung.
(a) Plan-Sperre weg, offene Plan-Schritte als Hinweiszeile beim gelungenen Move. Kosten: ein neuer Hinweis-Kanal in move.Result, Ausgabe in CLI (Text + JSON) und an drei TUI-Stellen. Gibt auf: die harte Sperre bei vergessenem Haken - die faengt aber keine echte Luecke, die DoD-Sperre und signoff/review bleiben.
(b) Beim Abschluss den letzten Plan-Schritt automatisch abhaken. Kosten: jaira schreibt einen Haken, den niemand gesetzt hat - die Datei behauptet Arbeit ohne Beleg. Trifft nur "der letzte" Schritt; offen kann jeder sein (Schritt 3 von 5). Alle abhaken waere Faelschen. Im Worker-Skill statt im Code hilft es der Person in der TUI nicht. Verworfen.
(c) Plan-Sperre ersatzlos streichen, kein Hinweis. Am billigsten (Block loeschen, Test drehen). Gibt auf: den sichtbaren Hinweis, den die DoD des Tickets ausdruecklich verlangt. Verworfen.
Entscheidung: (a). Die DoD ist das Ergebnis und allein das Gate; der Plan ist die Methode und wird nur gemeldet. Das passt zu dem, was der Agent-Block heute schon sagt.
Plan-Sperre bei UNVOLLSTAENDIGER DoD: entfaellt auch. Begruendung:
- Die DoD-Sperre lehnt dann sowieso ab (CodeNeedsSignal je offenem Kriterium, gate.go:362-375; ohne DoD-Liste gate.go:376-383).
- Behalten wuerde nur zusaetzliche Bullets und einen zweiten Code liefern, und eine Regel "Plan sperrt nur, wenn die DoD auch sperrt" - ein Gate, das von einem anderen Gate abhaengt.
- Einziger Fall, in dem sie allein sperren wuerde: eine Ziellane mit requires-nonmodel-signal: false. Die hat die Beleg-Pflicht bewusst abgewaehlt; dort per Plan weiter zu sperren widerspricht dieser Wahl.
- Vorschlag: der Hinweis erscheint nur bei einem Move, der landet. Eine Ablehnung sagt schon, was zu tun ist.
- **2026-10-26 19:18 · BeMuCa** — BRAINSTORM 3/3 - Beruehrungspunkte fuer pre-process (nur gefunden, nichts geaendert).
Code:
- core/gate/gate.go:404-432: Block und Kommentar weg oder umbauen; der Hinweis darf kein Violation sein, sonst zeigt die TUI ihn als Ablehnung (Vorpruefung model.go:1982, signoff.go:177).
- core/gate/gate.go:35 CodePlanIncomplete: wird nicht mehr ausgegeben; pre-process entscheidet: loeschen (plan_test.go benutzt ihn) oder behalten. Kein anderer Verbraucher: core/mod (Cockpit) verzweigt nicht darauf, board.ts:213-225; refusalCode flow.go:345-352 kennt ihn nicht.
- core/ticket/schema.go:579 PlanComplete() liefert die offenen Texte, aber keine Nummern; der Hinweis braucht die Nummer fuer 'jaira dod <id> <n> --done --plan'.
- core/move/move.go:58-70 Result + move.go:139-146: Hinweis hier eintragen, dann bekommen alle vier Schreibstellen ihn.
- internal/cli/flow.go:246-252 Text: Hinweis VOR der Zeile "X -> done" ausgeben (Regel BDV0HM dort: die letzte Zeile ist der Move).
- internal/cli/flow.go:232-244 JSON: neuer Schluessel. Vgl. offenes Ticket 01M26PW18D... "Bei --json steht eine Warnung im JSON, nicht auf einem zweiten Stream".
- internal/cli/flow.go:320-341 reportDryRun: Hinweis optional auch hier.
- TUI: internal/tui/model.go:2012-2027 (normaler Move, notify(...,false)), model.go:2072-2090 (forceMove), signoff.go:196-214 (Accept-Meldung).
- internal/cli/tickets.go:1615-1619 plan_complete im JSON: bleibt wahr, keine Aenderung.
Tests:
- core/gate/plan_test.go:36-63 TestUnfinishedPlanBlocksTerminalLane: dreht sich um (kein Refuse, Hinweis nennt "implement"); Kommentar 36-40 wird falsch. plan_test.go:67-104 pruefen nur "kein CodePlanIncomplete" - an den Konstanten-Entscheid anpassen.
- core/ticket/superseded_test.go:80 nutzt PlanComplete: unberuehrt.
Text, der falsch wird:
- .claude/skills/jaira/SKILL.md:105-107; README.md:247-249; docs/AGENTS.md:227-228; core/lane/builtin/50-done.md:11 (description).
- Bestehende Boards tragen die alte done-Beschreibung in .jaira/lanes/done.md:11. Ob das eine Correction braucht (core/lane/corrections.go) oder stehen bleiben darf: pre-process entscheidet.
- NICHT anfassen: core/lane/corrections.go:94 doneDoorway (woertliche Kopie einer ausgelieferten Lane, Geschichte), core/move/oldboard_test.go:29 (Fixture), core/release/NOTES.md:117 (geschlossen).
- core/board/announce.go:98-104 bleibt wahr; hoechstens ein Halbsatz "ein offener Plan-Schritt wird nur gemeldet".
- Release-Notiz core/release/notes/NJBTVC.md (geaenderte Move-Ausgabe, Move wird nicht mehr abgelehnt).
- **2026-10-26 19:24 · BeMuCa** — Driver: DoD um 3 Punkte ergaenzt, die der Brainstorm vorgeschlagen hat (nur-DoD-Ablehnung, Hinweis in Text/JSON/TUI, Doku). DoD-Punkt 2 (go test ./... -race) laeuft einmal vor dem PR, nicht in einer Lane.
