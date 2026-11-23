---
id: 01M845PACJSY9V5W6ZY096E1Z4
title: "Cockpit-Log-Filter: ein Chip schaltet ein, aus und auf 'außer'"
status: critique
ready: true
creator: Alexander Sacharov
goal: "In den Filterzeilen des Log-Fensters (Personen und Lanes) lässt sich jeder Chip auch ausschließen, damit 'alles außer X', 'alles außer X und Y' und 'nur A, aber nicht B' mit Klicks gehen."
context: |-
  Was heute fehlt: die Chips im Log-Fenster des jaira-cockpit (Zeile Personen: All / Namen; Zeile Lanes: All / ideas / backlog / ...) können nur einschließen. 'Alles außer done' heißt heute: jede andere Lane einzeln anklicken.
  Wo: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx filterRow() (Zeilen ~187-273), Zustand in overviewPeopleAtom/overviewLanesAtom, Aktionen in hooks/register.tsx ~2505 (everybody, allLanes, person, lane).
  Entscheidung des Maintainers (23.11.2026): Variante 'Zyklus pro Chip' - erster Klick wählt ein, zweiter Klick schließt aus (Chip zeigt '✕ name', rot/durchgestrichen), dritter Klick setzt zurück.
  Er will kombinieren können: mehrere ausgeschlossene Chips ('zwei nicht diese') und eingeschlossene neben ausgeschlossenen ('nicht diesen, aber diesen').
  Gleiche Regel für beide Zeilen. 'All' setzt die Zeile ganz zurück (nichts ein-, nichts ausgeschlossen).
  Verworfen: Klick bei aktivem All schließt aus (gleicher Klick, zwei Bedeutungen); Shift+Klick (im Terminal unzuverlässig).
definition-of-done: "Ein Chip durchläuft per Klick: aus -> eingeschlossen -> ausgeschlossen -> aus; der ausgeschlossene Zustand ist am Chip sichtbar (✕ vor dem Namen, andere Farbe)"
tags:
  - cockpit
blocked-by: []
related: []
commits: []
created-at: 2026-11-23T12:22:01Z
updated-at: 2026-11-23T12:29:31Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-75968
claimed-at: 2026-11-23T12:22:34Z
outcome-what: "Log filter chips in the cockpit cycle off → chosen → left out (✕, red) → off; each row's state is {only, not}; All clears both"
outcome-why: "'everything except done' needed one click per other lane; the maintainer chose a per-chip cycle on 23.11."
outcome-resolves: 96E1Z4 DoD 1-6
---

# Cockpit-Log-Filter: ein Chip schaltet ein, aus und auf 'außer'

## Definition of Done

- [x] Ein Chip durchläuft per Klick: aus -> eingeschlossen -> ausgeschlossen -> aus; der ausgeschlossene Zustand ist am Chip sichtbar (✕ vor dem Namen, andere Farbe)
  proof: core/mod/builtin/jaira-cockpit/hooks/overview.ts:208 cycled; core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:266; test 'a chip cycles off, chosen, left out and off again'
- [x] Nur Ausschlüsse: gezeigt wird alles außer den ausgeschlossenen; mehrere Ausschlüsse gehen
  proof: core/mod/builtin/jaira-cockpit/hooks/overview.ts:216 passes; test 'left out alone, the Log shows everything else, several left out at once'
- [x] Ein- und Ausschlüsse gemischt: gezeigt wird, was eingeschlossen und nicht ausgeschlossen ist; ein Eintrag wird über Person und Lane zugleich gefiltert
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx chosen = passes(people) && passes(lanes); test 'chosen and left out together, across both rows'
- [x] 'All' setzt seine Zeile zurück und ist nur dann hervorgehoben, wenn die Zeile weder Ein- noch Ausschlüsse hat
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:267 idle; test 'All clears its row of chosen and left out alike and is lit only then'
- [x] Tests in core/mod/builtin/jaira-cockpit/tests für Zyklus, nur-Ausschluss, gemischt, All-Reset
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts:3935-4040, claude plugin test . 484 pass
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/96E1Z4.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 12:22 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/6958.sock
- **2026-11-23 12:23 · Alexander Sacharov** — open decisions: none — variant 'cycle per chip' decided by maintainer in context. settled myself: excluded chip = '✕ name' in theme color 'error' with strikethrough; state per row becomes {only, not} (session atom, no migration needed)
- **2026-11-23 12:23 · Alexander Sacharov** — edges: an excluded chip with no event left must keep its button (like an included one) so it can be cycled back; filter state is session-only, no on-disk format change
- **2026-11-23 12:29 · Alexander Sacharov** — done differently than noted: no strikethrough — the mod's Text has no such prop in use anywhere; '✕ ' prefix + theme colour 'error' carry it. Three older Log tests changed: a second press now leaves out instead of clearing (test 3 presses twice more).
