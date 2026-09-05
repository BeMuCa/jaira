---
id: 01M1Q59ATPCREV9W8FNW6BFHJ4
title: CLI und TUI teilen EINE Settle-Funktion im core
status: critique
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Die Doorway/Cap-Entscheidung nach einem Move (logbook-on-entry -> FileLane mit Stempeln, holds -> TrimLane mit Pin) lebt einmal im core; flow.go und settleLane rufen sie statt je einen eigenen Switch zu tragen"
context: "Berk am 04.09. (Feedback L122): vier Status-Schreibstellen bedeuten, dass jede Nach-dem-Move-Regel vierfach verdrahtet wird - das Review fand die vierte (accept) nur durch Suchen. Erster Schnitt: die SETTLE-Semantik (welche Regel, welcher Ordner, Stempeln) zentral als core-Funktion; die UIs behalten nur ihre Meldungs-Formatierung. Die VOLLE Move-Vereinheitlichung (Gate+Mutate+Claim-Regeln aller vier Stellen) bleibt ein eigener groesserer Schnitt - separat eingefangen, hier bewusst nicht."
definition-of-done: Eine core-Funktion entscheidet Doorway/Cap und liefert Trimmed+filed+Fehler; flow.go und tui settleLane rufen sie (kein doppelter Switch mehr); Verhalten byte-gleich (bestehende Tests gruen ohne Anpassung)
tags: []
blocked-by: []
commits: []
created-at: 2026-09-04T21:30:37Z
updated-at: 2026-09-05T02:33:02Z
claimed-by: EE-3NX6GL3-404958
claimed-at: 2026-09-05T02:29:45Z
updated-by: BeMuCa
outcome-what: "lane.Settle im core: die Doorway/Cap-Entscheidung samt Vorrang lebt einmal; CLI-flow und TUI-settleLane delegieren und behalten nur ihre Meldungs-Formatierung - die zwei parallelen Switches sind weg"
outcome-why: "Feedback L122, erster Schnitt: vier Schreibstellen hiessen bisher zwei identische Regel-Switches, die auseinanderdriften koennen - das Review fand accept() nur durch Suchen; die Regel gehoert dahin, wo lane und ticket sich treffen (lane importiert ticket, nicht umgekehrt - daher core/lane)"
outcome-resolves: "Bestehende Doorway/Cap/Jam-Tests alle gruen ohne eine Zeile Anpassung (byte-gleiches Verhalten), plus Routing-Unit-Test inkl. Doorway-schlaegt-Cap; die volle Move-Vereinheitlichung ist als JJ32B4 im Backlog eingefangen, bewusst nicht hier"
executed-by: fable
---

# CLI und TUI teilen EINE Settle-Funktion im core

## Definition of Done

- [x] Eine core-Funktion entscheidet Doorway/Cap und liefert Trimmed+filed+Fehler; flow.go und tui settleLane rufen sie (kein doppelter Switch mehr); Verhalten byte-gleich (bestehende Tests gruen ohne Anpassung)
  proof: core/lane/settle.go: Settle entscheidet Doorway/Cap einmal (Doorway gewinnt); flow.go und tui settleLane rufen sie, beide Switches geloescht; TestSettleRoutesDoorwayThenCap (nil/plain/Cap-mit-Pin/Doorway-gewinnt); Verhalten byte-gleich: voller Lauf -race RC=0 OHNE Anpassung bestehender Tests (test18)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

