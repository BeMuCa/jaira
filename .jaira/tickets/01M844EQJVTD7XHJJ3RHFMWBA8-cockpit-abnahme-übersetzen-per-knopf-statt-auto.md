---
id: 01M844EQJVTD7XHJJ3RHFMWBA8
title: "Cockpit-Abnahme: Übersetzen per Knopf statt automatisch, alles in einem Aufruf"
status: critique
ready: true
creator: Alexander Sacharov
goal: Die Abnahme-Ansicht übersetzt nur noch auf Knopfdruck; ein Druck übersetzt Kopf und alle Notizen auf einmal.
context: |-
  Was heute falsch ist: Die Abnahme-Ansicht (Приём) im jaira-cockpit übersetzt jedes Öffnen automatisch, sobald settings.json 'language' nennt.
  Wo: core/mod/builtin/jaira-cockpit/hooks/register.tsx, readStory() ruft rephrase() zweimal auf - einmal für die Review-Notizen (rephrasePrompt), einmal für den Kopf mit Titel, Ziel, Was/Warum und Ampelgrund (headPrompt, Schlüssel id+':head').
  Folge: zwei haiku-Aufrufe bei jedem neuen Ticket, und der Text springt vom Original zur Übersetzung, während man liest.
  Wunsch des Maintainers (23.11.2026): ein Knopf 'Übersetzen', damit es nicht immer automatisch passiert - schnell und bequem - und ein Druck übersetzt alles auf einmal.
  Bekannt: Übersetzungen werden im Store unter REPHRASED_PREFIX gehalten (rephrased(), storeNewest()); ein schon übersetztes Ticket kann aus dem Store sofort gezeigt werden, ohne Aufruf.
  Offen gelassen: ob 'alles' auch die Explain-Begründungen (reasonPrompt) einschließt - Explain wird ohnehin pro Notiz gedrückt, bleibt also wie es ist.
definition-of-done: "Ohne Druck zeigt die Abnahme-Ansicht den Originaltext und ruft kein Modell zur Übersetzung auf, auch wenn settings.json 'language' setzt"
tags:
  - cockpit
  - acceptance
blocked-by: []
related: []
commits: []
created-at: 2026-11-23T12:00:24Z
updated-at: 2026-11-23T12:13:25Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-68987
claimed-at: 2026-11-23T12:01:00Z
outcome-what: "Acceptance pane translates only on Translate (key t): head and all review notes in one haiku call, kept in the store; a second press or a new opening shows the original. Removed the two automatic calls in readStory, rephrasePrompt/headPrompt replaced by translatePrompt."
outcome-why: "Every opening with 'language' set made two model calls and the text jumped from original to translation while the person read."
outcome-resolves: "DoD 1-6: no call without a press, one call per press, store hit without call, no button without language, tests, release note. Checked: concurrency (stale answers after second press/reopen show nothing, mutation-tested), error paths (failed/misfitting answer keeps original, next press asks again), unusual input (no notes -> prompt without notes part; empty goal filtered), permissions (nothing written to the ticket)."
review-summary: |-
  core/mod/builtin/jaira-cockpit/hooks/register.tsx:1534 translate() builds the head from s.title/s.goal while the pane draws story_?.title || w.title and story_?.goal || w.goal - a ticket with an empty story goal shows w.goal untranslated; build head from the same expressions the render uses
  core/mod/builtin/jaira-cockpit/hooks/register.tsx:1531 translate() re-reads settings.json although openAcceptance already put the language into storyAtom - read s.language: the button and the call then agree, one read less
  core/mod/builtin/jaira-cockpit/hooks/register.tsx:1551 untranslate() fires its atom updates with void, then translate() awaits update(translatingAtom, id) - the null may land after the id; make untranslate async and await both updates
round: 2
back-from: critique
back-by: lane
back-count: 1
---

# Cockpit-Abnahme: Übersetzen per Knopf statt automatisch, alles in einem Aufruf

## Definition of Done

- [x] Ohne Druck zeigt die Abnahme-Ansicht den Originaltext und ruft kein Modell zur Übersetzung auf, auch wenn settings.json 'language' setzt
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1359 (openAcceptance resets, readStory asks no model); test 'with a language set the pane opens on the original and asks nothing until Translate is pressed'
- [x] Ein Knopf 'Übersetzen' (mit Taste) in der Abnahme-Ansicht übersetzt Kopf und alle Notizen in einem einzigen Modellaufruf
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2409 Button hotkey t; register.tsx:1544 + board.ts:1052 translatePrompt; tests 'one press translates the head and every note in one call; a second press goes back to the original', 'a head the pane draws from the board row - no goal in show - is translated as drawn'
- [x] Ein Ticket, dessen Übersetzung schon im Store liegt, wird nach dem Druck ohne neuen Aufruf übersetzt gezeigt; ein zweiter Druck schaltet zurück zum Original
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1562 store hit; test 'kept in the store per ticket and texts: a press on a ticket translated before asks nothing, ...' and 'one press translates ... a second press goes back'
- [x] Ohne 'language' in settings.json ist der Knopf nicht zu sehen
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2408 s?.language guard; test 'each class has its mark in the Acceptance pane; without a language there is no Translate and nothing is sent'
- [x] Tests in core/mod/builtin/jaira-cockpit/tests decken: kein Aufruf ohne Druck, ein Aufruf mit Druck, Store-Treffer ohne Aufruf
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts describe 'review notes as a person reads them' (15 tests; claude plugin test . 481 pass)
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/FMWBA8.md
- [ ] seen by a person: in the cockpit's Acceptance pane with 'language' set, the head and notes show the original; 't' or the Translate button turns them into the language at once, a second press back

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 12:01 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/6958.sock
- **2026-11-23 12:02 · Alexander Sacharov** — decisions (worker, none open for a person): 1) key 't' as Button hotkey (letter allowed by API; digits 1-4 are the Explanation pane's). 2) each opening of the pane starts on the original - the toggle is reset in openAcceptance, so DoD 1 holds on every open. 3) one prompt (head texts first, then notes) -> one JSON array, stored under one key rephrased:<id>; old '<id>:head' entries age out through storeNewest's 50-ticket cap. 4) a failed or misfitting call leaves the original and the toggle off; a new press asks again. 5) while the call runs the button reads 'Translating…'.
- **2026-11-23 12:07 · Alexander Sacharov** — in-progress: the old auto-translation ran two calls (notes, '<id>:head'); now translate() builds head+notes and asks once via rephrased() with translatePrompt. Stale-call guard (newest map) checked by mutation: removing it fails 'a call answered after a second press...'. Open risk for the person: whether the 't' hotkey fires while the reason Input has focus is engine-defined (API doc: hotkey works while the pane holds focus, says nothing of Input) - check at acceptance. .jaira/lanes, AGENTS.md, CLAUDE.md were rewritten by jaira on claim (0.6.4 block), not by this ticket: left out of the commit.
- **2026-11-23 12:08 · Alexander Sacharov** — critique: (1) head texts must use the render's fallbacks (story title/goal || w.title/w.goal); (2) translate takes the language from storyAtom.language, not a second settings read; (3) untranslate awaits its updates so 'translating' cannot be cleared after it was set.
- **2026-11-23 12:13 · Alexander Sacharov** — in-progress round 2: all three critique findings fixed. Learned: a superseded call's answer is still kept in the store (rephrased() stores before the newest check) - the next press is then a store hit with no call; the race test now asserts that instead of expecting a second call.
