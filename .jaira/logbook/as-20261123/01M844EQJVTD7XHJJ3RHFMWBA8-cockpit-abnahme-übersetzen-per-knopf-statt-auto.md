---
id: 01M844EQJVTD7XHJJ3RHFMWBA8
title: "Cockpit-Abnahme: Übersetzen per Knopf statt automatisch, alles in einem Aufruf"
status: done
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
commits:
  - eadc689614cb0d1c538230adc8753997145da5b1
  - 46e4bb718ceef830903881a1c8028df6d033db01
  - d36de86f81c06f891621a73823680c51cfa66b91
created-at: 2026-11-23T12:00:24Z
updated-at: 2026-11-23T12:39:54Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-68987
claimed-at: 2026-11-23T12:01:00Z
outcome-what: "Acceptance pane translates only on Translate (key t): head and all review notes in one haiku call, kept in the store; a second press or a new opening shows the original. Removed the two automatic calls in readStory, rephrasePrompt/headPrompt replaced by translatePrompt."
outcome-why: "Every opening with 'language' set made two model calls and the text jumped from original to translation while the person read."
outcome-resolves: "DoD 1-6: no call without a press, one call per press, store hit without call, no button without language, tests, release note. Checked: concurrency (stale answers after second press/reopen show nothing, mutation-tested), error paths (failed/misfitting answer keeps original, next press asks again), unusual input (no notes -> prompt without notes part; empty goal filtered), permissions (nothing written to the ticket)."
review-summary: "The Acceptance pane no longer calls the model when it opens: openAcceptance resets the ticket to the original (untranslate) and only records settings.json 'language' on the pane state. A new Translate button (hotkey t, shown only when a language is set) builds the head exactly as the pane draws it (title/goal with board-row fallbacks, what, why, the light's reason) plus every review note, and sends them in one haiku call via the new translatePrompt; the answer is kept in the store under rephrased:<id>, so a ticket translated before shows at once. A second press, or reopening the pane, goes back to the original; a call answered after that is dropped through the newest-press counter. The old two-call auto-translation (rephrasePrompt, headPrompt, ':head' key) is removed. Full review at d36de86f8, 481 plugin tests pass."
round: 2
back-from: critique
back-by: lane
back-count: 1
review-gaps: |-
  note: edge: whether hotkey t fires while the reason Input has focus (typing a t toggles translation instead of writing it) is engine-defined - check it at acceptance; if it does, the hotkey needs a guard or another key
  note: limit: head and notes share one answer under maxTokens 4000 (register.tsx rephrased); a ticket with very long texts and many notes may get a cut answer and stays in the original
  note: limit: two presses inside the few awaits before translating is set start two calls instead of toggling - the newer wins, per the 22.11.2026 decision not to hold calls back
test-verdict: "pass: claude plugin test . 481 pass 0 fail, go test ./core/release ok, scripts/shipped-history.sh --check RC=0; DoD 1-6 verified by their named tests in the engine's pane harness (no call on open, one call per press with head+notes, store hit without call, no button without language, release note present); failure path: model down / misfitting answer keeps the original and the next press asks again; two writers: concurrent presses and a reopen while a call waits - the stale answer shows nothing (mutation-checked), the store write itself is unlocked as before this change. DoD 7 (seen by a person: live cockpit, and whether 't' fires while the reason field has focus) left for the person."
hunt-verdict: "pass: bughunt done at d36de86f8, 0 change, 0 pre-existing; two suspicions not confirmed (stuck 'translating…' label on a reopen race; 4000-token cap with head+notes in one answer) - in its hunt report note"
review-verdict: "accept: DoD 1-6 met by the diff and tests (481 pass); no defect found along concurrency, error paths, input or boundaries"
review-round: 2
review-verdict-gaps: |-
  note: edge: whether hotkey t fires while the reason Input has focus (typing a t toggles translation instead of writing it) is engine-defined - check it at acceptance; if it does, the hotkey needs a guard or another key
  note: limit: head and notes share one answer under maxTokens 4000 (register.tsx rephrased); a ticket with very long texts and many notes may get a cut answer and stays in the original
  note: limit: two presses inside the few awaits before translating is set start two calls instead of toggling - the newer wins, per the 22.11.2026 decision not to hold calls back
review-check: "1. In ~/.claude/settings.json set \"language\": \"Russian\"  2. Start the cockpit: claude --plugin-dir /home/alex/projects/.worktrees/jaira-fmwba8/core/mod/builtin/jaira-cockpit  3. Open the Acceptance pane of a ticket in signoff - title, goal and notes are in the original language, button reads Translate  4. Press t (or the Translate button) - label shows translating…, then head and all notes appear in Russian, label reads Original  5. Press t again - the original returns  6. Press t once more - Russian appears at once, no wait (store hit)  7. Click into the reason input and type a word containing t - check the letter is typed and translation does not toggle  8. Remove language from settings.json, reopen the pane - no Translate button"
question: ""
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
- [x] seen by a person: in the cockpit's Acceptance pane with 'language' set, the head and notes show the original; 't' or the Translate button turns them into the language at once, a second press back
  proof: Alex, cockpit Acceptance pane check 23.11.2026: t translates head and notes, second press back, t types in the reason field
- [-] Der Übersetzen-Druck übersetzt auch die Mitte der Ansicht: den Review-Satz im History-Block und die Release-Note-Zeile - im selben einen Aufruf

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
- **2026-11-23 12:13 · Alexander Sacharov** — critique pass 2: the three findings of pass 1 are addressed (head uses the render's fallbacks, language from storyAtom, untranslate awaited); nothing else reopened.
- **2026-11-23 12:15 · Alexander Sacharov** — waiting: bughunt (background subagent) since 13:15
- **2026-11-23 12:20 · Alexander Sacharov** — hunt report: 0 findings, nothing written to the board | read whole: hooks/register.tsx (acceptance pane, translate/untranslate/rephrased/storeNewest, render), hooks/board.ts (languageOf, translatePrompt, rephrasedOf), hooks/texts.ts, types/index.d.ts, tests/board.test.ts (translate tests), release note; the uncommitted .jaira/lanes/*.md, AGENTS.md and CLAUDE.md changes are lane-prompt prose with no code path, read for content only | not confirmed, so not a finding: (1) translate() takes 'mine' and sets translatingAtom only after several awaits, so a reopen landing in that window could leave the button on 'translating…' until one press; a scratch test with the reopen racing the press in both orders did not hit it (ran: claude plugin test on a scratch copy, button back to 'Translate', next press translated). (2) head and notes now share one call with maxTokens 4000, so a ticket with long what/why and many notes may get a cut-off answer and never translate; not checkable without the real model
- **2026-11-23 12:21 · Alexander Sacharov** — bughunt: done at d36de86f8 - 0 findings (0 change, 0 pre-existing) in 5 modules, first hunt; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/types/index.d.ts core/release/notes/FMWBA8.md
- **2026-11-23 12:21 · Alexander Sacharov** — testing: note: limit: head and notes now share one answer under maxTokens 4000 (register.tsx rephrased); a ticket with very long what/why and many notes could get a cut answer, which shows the original - not reproducible without the real model.
- **2026-11-23 12:21 · Alexander Sacharov** — waiting: review lane (fresh agent via run-lane.sh) since 13:21
- **2026-11-23 12:24 · Alexander Sacharov** — review: done at d36de86f8
- **2026-11-23 12:24 · Alexander Sacharov** — review: accept written, but move to signoff refused - the signoff hunt gate sees uncommitted changes in .jaira/lanes/{critique,in-progress,review,testing}.md, AGENTS.md, CLAUDE.md (jaira 0.6.4 regeneration on claim, not this ticket's code; the hunt read them). Options:
1. (recommended) commit those regenerated files on this branch as chore(FMWBA8): lane prompts regenerated by jaira 0.6.4 together with the ticket file, then move - they are prose, no code path, so no new hunt
2. git restore those files to HEAD, then move - the regeneration comes back on the next claim
3. a person moves with --force
- **2026-11-23 12:25 · Alexander Sacharov** — question (from review agent, carried to parent): review accepts (DoD 1-6, 3 notes), but the move to signoff is refused: uncommitted .jaira/lanes/*.md, AGENTS.md, CLAUDE.md, regenerated by jaira 0.6.4 on claim - not this ticket's change. Options: 1) restore them with git checkout (Recommended: they are regenerated anyway, and a version bump of the agent block belongs to its own update commit, not to FMWBA8); 2) commit them on this branch as chore(FMWBA8); 3) move with --force, files left as they are.
- **2026-11-23 12:25 · Alexander Sacharov** — decision (teamlead 2026-11-23): option 1 - regenerated board files restored; master carries the 0.6.4 files (PR #224).
- **2026-11-23 12:25 · Alexander Sacharov** — report: Acceptance pane translates only on Translate (key t), head+notes in one haiku call, store hit without call, second press/reopen back to original; 3 commits, 481 tests | signoff, context 196k of 1000k, cache warm 60 min after last answer | open for the person: DoD 7 seen-by-a-person, incl. whether 't' fires while the reason field has focus; review accept with 3 notes; branch not pushed
- **2026-11-23 12:38 · Alexander Sacharov** — decision (Alex 23.11.2026, Acceptance check): 't' works, and typing 't' in the reason field is fine. Gap: the middle of the pane stays untranslated - the History block (review verdict sentence) and the Release note line. The translate press must cover them too, in the same single call.
- **2026-11-23 12:39 · Alexander Sacharov** — decision (Alex 23.11.2026): accept as is - review sentence and release note stay English; DoD 8 superseded. Pane readability (sections/stripes) goes to a new ticket.
