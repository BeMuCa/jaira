---
id: 01M5MNWPVYGSYTYCR3PVNJBTVC
title: "Ein offener Plan-Punkt sperrt die Ziellane nicht mehr, wenn die Definition of Done erfuellt ist"
status: backlog
ready: false
creator: BeMuCa
goal: "Ein Ticket mit erfuellter Definition of Done kommt nach done, auch wenn ein Plan-Schritt nicht abgehakt ist; der Plan ist die Methode, die DoD das Ergebnis."
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
updated-at: 2026-10-26T19:17:57Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-264760
claimed-at: 2026-10-26T19:12:35Z
---

# Ein offener Plan-Punkt sperrt die Ziellane nicht mehr, wenn die Definition of Done erfuellt ist

## Definition of Done

- [ ] Ein Move in eine Ziellane mit vollstaendiger DoD und offenem Plan-Punkt wird nicht mehr abgelehnt; die Ausgabe nennt den offenen Plan-Punkt als Hinweis; ein Test deckt den Fall
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

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
