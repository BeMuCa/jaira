---
id: 01M8488Z540RWJ8F6J8CDKMG91
title: "Cockpit-Abnahme: offene 'seen by a person'-Punkte stehen als eigener Abschnitt da und werden mit Accept abgehakt"
status: testing
ready: true
creator: Alexander Sacharov
goal: "Eine Person nimmt ein Ticket mit offenem 'seen by a person'-DoD in einem Schritt an: die Ansicht zeigt diese Punkte als eigenen Abschnitt 'Selbst prüfen', und Accept hakt sie ab und nimmt dann an - statt mit 'criterion N is still open' abzulehnen."
context: |-
  Was heute falsch ist: Ein Ticket mit einem DoD-Punkt 'seen by a person: ...' kommt in signoff mit diesem Punkt offen an - das ist Absicht, nur die Person kann ihn abhaken. Drückt der Maintainer im jaira-cockpit (Abnahme-Ansicht, Приём) Accept, lehnt jaira ab: 'the definition of done is not met: criterion 7 (...) is still open. Satisfy it, mark it with jaira dod <ID> 7 --done, then try this move again'.
  Passiert am 23.11.2026 zweimal hintereinander (FMWBA8 DoD 7, MK8HD2 DoD 5). Jedes Mal musste der Teamlead den Punkt per CLI abhaken und der Maintainer ein zweites Mal Accept drücken.
  Die Ansicht zeigt heute nirgends, dass noch ein Punkt auf die Person wartet - man erfährt es erst aus der Ablehnung.
  Wunsch des Maintainers: ein eigener Absatz in der Ansicht, der zeigt, was er noch selbst tun/ansehen muss; das tun, dann annehmen - in einem Durchgang.
  Wo: core/mod/builtin/jaira-cockpit/hooks/register.tsx (Abnahme-Ansicht, Execute/Accept, executionOf in hooks/board.ts), DoD-Daten aus 'jaira show --json' (dod_items: n, text, done).
  Erkennung: DoD-Text beginnt mit 'seen by a person' (so schreiben die Lanes ihn heute).
definition-of-done: "Hat das Ticket offene DoD-Punkte, die mit 'seen by a person' beginnen, zeigt die Abnahme-Ansicht sie oberhalb der Review-Notizen als eigenen, abgesetzten Abschnitt (Titel z.B. 'Selbst prüfen'), jeden Punkt mit seinem Text"
tags:
  - cockpit
  - acceptance
blocked-by: []
related: []
commits: []
created-at: 2026-11-23T13:07:10Z
updated-at: 2026-11-23T13:17:15Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-97038
claimed-at: 2026-11-23T13:07:36Z
outcome-what: "Acceptance pane lists open 'seen by a person' DoD items under 'Check yourself' above the review notes; Execute needs each ticked, then runs jaira dod <id> <n> --done --proof 'seen by <whoami> in the cockpit' before the gates and accepts in one press"
outcome-why: "Accept was refused with 'criterion N is still open' for items only a person can tick; teamlead had to tick by CLI and the maintainer press twice (FMWBA8, MK8HD2 on 23.11.)"
outcome-resolves: "one-press acceptance of tickets with 'seen by a person' items"
review-summary: none
review-gaps: "folded the open-item test into one isOpen in board.ts (openOf and seenItemsOf shared it); hoisted the unticked list in execute (was filtered twice). Left: storyOf's own done/superseded count (pre-existing, counts done, not open)."
---

# Cockpit-Abnahme: offene 'seen by a person'-Punkte stehen als eigener Abschnitt da und werden mit Accept abgehakt

## Definition of Done

- [x] Hat das Ticket offene DoD-Punkte, die mit 'seen by a person' beginnen, zeigt die Abnahme-Ansicht sie oberhalb der Review-Notizen als eigenen, abgesetzten Abschnitt (Titel z.B. 'Selbst prüfen'), jeden Punkt mit seinem Text
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2376; test 'core/mod/builtin/jaira-cockpit/tests/board.test.ts › what the person checks themselves › open items stand above the review notes...'
- [x] Jeder dieser Punkte hat einen Haken, den die Person setzt; Accept ist erst möglich, wenn alle gesetzt sind, und sagt sonst in der Zeile, was noch fehlt
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1496 (needSeen) + texts.ts summary 'needSeen'; test 'Execute with an item not ticked accepts nothing and names the item'
- [x] Accept mit gesetzten Haken ruft für jeden Punkt 'jaira dod <id> <n> --done --proof "seen by <Person> in the cockpit"' auf und dann die Annahme - ein Druck, keine Ablehnung wegen offener Punkte
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1637 acceptWith before-hook; test 'ticked, one Execute ticks it on the ticket with the person named, then accepts'
- [x] Ein Ticket ohne solche Punkte sieht und verhält sich wie bisher
  proof: test 'a ticket without such items shows no section, and Execute accepts as before'; full suite 486 pass
- [x] Tests in core/mod/builtin/jaira-cockpit/tests: Abschnitt erscheint nur bei offenen Punkten; Accept ohne Haken gesperrt; mit Haken werden die dod-Aufrufe vor dem move gemacht
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts:5569 describe 'what the person checks themselves' (6 tests)
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/DKMG91.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 13:07 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/6958.sock
- **2026-11-23 13:08 · Alexander Sacharov** — decisions (worker, none open for the person): 1) section title in English ('Check yourself') - texts.ts keeps every cockpit text English; 'Selbst prüfen' in the DoD was 'z.B.'. 2) person name in the proof = 'jaira whoami' identity, 'a person' where it says none. 3) ticks are written right before the dry run of the move: a later gate refusal leaves them ticked - true anyway, the person did look. 4) only the Acceptance pane's Execute ticks; Accept from an Explain choice keeps today's refusal (outside this DoD). 5) items come from 'jaira show --json' read when the pane opens; path: todo -> in-progress (no planning option).
- **2026-11-23 13:14 · Alexander Sacharov** — in-progress: tests via 'claude plugin test .' in core/mod/builtin/jaira-cockpit (486 pass, ~60s). No typecheck in CI for the mod. Unrelated dirty files in worktree (.jaira/lanes/*, AGENTS.md, CLAUDE.md - jaira's own regen) left out of commits.
- **2026-11-23 13:15 · Alexander Sacharov** — critique: none. Checked: before-hook in accept() is needed because the ticket's checkout (dir) is only known inside accept/ticketHere; jaira dod --done on an already ticked item exits 0, so a second press after a gate refusal re-ticks harmlessly; needWhy precedes needSeen on purpose (reason first, then ticks).
