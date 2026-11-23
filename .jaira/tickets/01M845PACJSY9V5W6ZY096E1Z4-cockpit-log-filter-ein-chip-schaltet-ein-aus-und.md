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
updated-at: 2026-11-23T12:37:42Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-75968
claimed-at: 2026-11-23T12:22:34Z
outcome-what: "Log filter chips cycle off → chosen → left out (✕, red) → off; per-row state {only, not} under new atom keys; All clears both"
outcome-why: "'everything except done' needed one click per other lane; the maintainer chose a per-chip cycle on 23.11."
outcome-resolves: 96E1Z4 DoD 1-6; bughunt blocker on old-build atom shape
review-summary: none
review-gaps: "hoisted f.only.includes(id) into 'on' in filterRow (was computed twice); considered: no other filter cycle in the mod to reuse, toggled() removed with its last callers, no dead imports; left the local 'set' helper (two callers, names the union)"
test-verdict: "fail: go test ./... RC=1 — core/role TestEveryInstalledFileIsAShippedVersion: shipped.sha256 lacks the 4 changed cockpit files; cockpit suite 484/484 green, DoD 1-6 verified in tree"
hunt-verdict: "fail: bughunt done at f6394636d — 1 change blocker: an old-build string[] in overviewPeople/overviewLanes crashes the Overview after a hot reload"
round: 2
back-from: testing
back-by: lane
back-count: 1
---

# Cockpit-Log-Filter: ein Chip schaltet ein, aus und auf 'außer'

## Definition of Done

- [x] Ein Chip durchläuft per Klick: aus -> eingeschlossen -> ausgeschlossen -> aus; der ausgeschlossene Zustand ist am Chip sichtbar (✕ vor dem Namen, andere Farbe)
  proof: core/mod/builtin/jaira-cockpit/hooks/overview.ts:208 cycled; hooks/overview-pane.tsx:266 filterRow; test 'a chip cycles off, chosen, left out and off again'
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
- **2026-11-23 12:29 · Alexander Sacharov** — critique: none — pure cycled/passes beside pickOf in overview.ts (same home as other pure Log helpers); register.tsx only swaps toggled for cycled; no existing exclude pattern in the mod (grep exclud/negat: only board.ts verdict wording, unrelated). AGENTS.md, CLAUDE.md, .jaira/lanes/*.md were modified in the worktree by a jaira regeneration, not by this ticket — left uncommitted on purpose.
- **2026-11-23 12:33 · Alexander Sacharov** — testing: class: every file the build installs must have its hash in core/shipped/shipped.sha256 — any edit to an installed file (core/mod/builtin/**, roles, lanes) needs it; example: shipped_test.go:74 lacks overview-pane.tsx, overview.ts, register.tsx, types/index.d.ts; fix: run scripts/shipped-history.sh --write and commit the manifest with the code, then go test ./core/role
- **2026-11-23 12:33 · Alexander Sacharov** — waiting: bughunt agent (background subagent of this worker) since 13:33
- **2026-11-23 12:35 · Alexander Sacharov** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:254 - class: overviewPeople/overviewLanes atom values in the shape an older build stored (string[]) are read as Filter without a shape check; the same happens in passes() and cycled() (overview.ts:208,216), reached through register.tsx:2494 and 2502/2504; example: a session that pressed a Log chip or All under the previous build keeps ['alex'] or [] in the host store (types/index.d.ts says the state survives a hot reload; overview-pane.tsx:259 already guards view.lanes for the same reason). After the mod reloads, set(state.people) throws 'f.only is not iterable', and because the log rows are built before the tab is chosen the whole Overview pane fails on every tab (Queue and Workers too). A chip press throws in cycled(). Fix idea: normalise on read (Array.isArray(v) ? {only: v, not: []} : v) or give the atoms a new key; ran: node --experimental-strip-types with passes(['alex'],'bob') / cycled(['alex'],'bob') / [...['alex'].only]: TypeError each time; blame: cfb151663 (names 96E1Z4)
- **2026-11-23 12:35 · Alexander Sacharov** — hunt report: 1 finding, 1 change (blocker), note on 96E1Z4 | read whole: hooks/overview.ts, hooks/overview-pane.tsx, types/index.d.ts; register.tsx read in the overview atoms and the Overview render/update hook and not end to end (2599 lines); tests/board.test.ts and the uncommitted lane/AGENTS/CLAUDE regenerations read as a diff only (prompt text, no code paths) | not checked: whether the host really hands an old-build atom value to a reloaded mod (taken from the types comment and the existing view.lanes guard), not run inside Claude Code; concurrency of update() is the host's job
- **2026-11-23 12:35 · Alexander Sacharov** — bughunt: done at f6394636d - 1 finding (1 change, 0 pre-existing) in 6 modules, first hunt; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx core/mod/builtin/jaira-cockpit/hooks/overview.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/types/index.d.ts core/release/notes/96E1Z4.md
- **2026-11-23 12:37 · Alexander Sacharov** — fix for bughunt blocker: new atom keys overviewPeopleFilter/overviewLanesFilter instead of a shape guard — an old string[] under the old key is never read, no normalise code to maintain; manifest regenerated with scripts/shipped-history.sh --write. No test seeds an old-build atom: the harness offers no way to preload the host store.
