---
id: 01M7ZYRRKJ8Q77ZZ3B8WY2CD9X
title: "Rollen, die weder Code ändern noch urteilen, laufen auf haiku"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "PR-, CI- und Test-Läufe kosten haiku-Preise, ohne dass ein Urteil schwächer wird."
context: |-
  Idee Alex 21.11.: was keinen Code ändert, auf haiku.
  Abgrenzung (Teamlead): nicht nur "ändert keinen Code", sondern auch "urteilt nicht". critique, review und bughunt ändern keinen Code, sind aber Urteile — die bleiben strong.
  Kandidaten für haiku: jaira-role-pr (push, PR öffnen, Kommentare beantworten ohne Codeänderung), jaira-role-tester (Suite laufen lassen, geerbte von neuen Fehlern trennen), das Warten auf CI, jaira-role-logbook-summary.
  Heute: Lanes haben tier cheap/strong; Rollen legen kein Modell fest (kein model: in core/role/builtin/*/SKILL.md); Chats nutzen .jaira/chat-model (sonnet). spawn.sh startet claude ohne Modellwahl pro Rolle.
  Offen für pre-process: Mechanismus (model: in der SKILL-Frontmatter, Feld pro Rolle in den Board-Settings oder spawn.sh --model) und ob ein PR-Kommentar, der Codeänderung verlangt, an einen Worker auf dem normalen Modell geht.
definition-of-done: "pr-, tester- und logbook-summary-Rolle starten auf haiku, critique/review/bughunt bleiben strong (Test oder roles list zeigt es)"
tags:
  - gates
blocked-by: []
related: []
commits:
  - 2a672cfb40c950cc9c13d61be7a94fddbfc65801
  - 7bb4c47dbf7cd5b5cd3476f2c4165d0866ad14a1
  - 6f57eee71482b17cc1c8bf81bfa64d87c440297b
  - 6822ec61b3857da5339ab2d9c7eb7aba3b892c77
  - 270f5b044686d3210f1131e87d9a657b42b60c6a
  - b5ec4070197353a63d6c813e5329b780c978a395
created-at: 2026-11-21T21:04:04Z
updated-at: 2026-11-21T22:45:50Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1274804
claimed-at: 2026-11-21T22:44:29Z
outcome-what: "tester and logbook-summary roles carry model: haiku; the pr role carries none and hands code-changing review points to an agent that inherits the session's model"
outcome-why: "pushing and running suites need no strong model, but a model: line on the pr role would be inherited by its code-fix agent"
outcome-resolves: "test and summary runs cost haiku prices; critique, review, bughunt and code fixes stay on the session model"
review-summary: "Adds 'model: haiku' to the SKILL.md headers of jaira-role-tester and jaira-role-logbook-summary, so Claude Code runs those two roles on haiku. jaira-role-pr keeps no model: line (it had one in round 1, removed in 6822ec61b after the bughunt proved a subagent inherits the skill's haiku); its comment-answering step 1 now hands a code-changing review point to an Agent with no model, which therefore runs on the session's model. TestRoleModels pins the model field of every builtin role (haiku only for tester and logbook-summary, empty for all others incl. pr, lane, bughunt). Release note core/release/notes/Y2CD9X.md has two lines; shipped.sha256 records the new hashes of the three SKILL.md files (re-recorded after the upstream/master merge)."
round: 4
back-from: testing
back-by: lane
back-count: 1
review-gaps: |-
  note: known limit, not testable from Go: that Claude Code switches to haiku on a model: line in a SKILL.md header is harness behaviour; TestRoleModels only pins the header text
  note: known limit: an account or provider without the haiku alias (or --agent copilot, which likely ignores model:) runs tester and logbook-summary on whatever it falls back to; nothing in jaira checks this
  note: the worktree holds uncommitted 0.6.0 jaira-update edits (.jaira/lanes/{critique,in-progress,optimize,testing}.md, AGENTS.md, CLAUDE.md) that are not this ticket - keep them out of the logbook commit (git add only the ticket file)
test-verdict: "pass: round 4 at b5ec40701: full go test ./... after upstream merge green except TestReapTearsDownWhatTheSetupBuilt/hangs once under load (passes alone; reap.sh untouched); vet RC=0, shipped-history --check RC=0, TestRoleModels -race ok; built binary installs model: haiku in exactly tester and logbook-summary, pr has none and carries the new step 1; DoD 1-4 checked in tree; writes nothing new beyond roles install (parallel-install race noted earlier, pre-existing)"
hunt-verdict: "pass: bughunt done at b5ec40701, 0 change, 0 pre-existing"
review-verdict: "accept: all 4 DoD items hold in the diff (tester+logbook-summary on haiku, pr/lane/bughunt without model:, pr fix agent inherits the session model, note + shipped hashes); go test ./core/role ./core/shipped pass; no defect found in the touched files"
review-round: 4
review-verdict-gaps: |-
  note: known limit, not testable from Go: that Claude Code switches to haiku on a model: line in a SKILL.md header is harness behaviour; TestRoleModels only pins the header text
  note: known limit: an account or provider without the haiku alias (or --agent copilot, which likely ignores model:) runs tester and logbook-summary on whatever it falls back to; nothing in jaira checks this
  note: the worktree holds uncommitted 0.6.0 jaira-update edits (.jaira/lanes/{critique,in-progress,optimize,testing}.md, AGENTS.md, CLAUDE.md) that are not this ticket - keep them out of the logbook commit (git add only the ticket file)
review-check: "1. cd /home/alex/projects/.worktrees/jaira-y2cd9x  2. go test ./core/role -run TestRoleModels -v  -> PASS  3. grep -rn '^model:' core/role/builtin/*/SKILL.md  -> exactly two lines: jaira-role-tester and jaira-role-logbook-summary, both 'model: haiku'  4. sed -n 480,495p core/role/builtin/jaira-role-pr/SKILL.md  -> step 1 hands a code fix to Agent(...) with no model and says why pr has no model: line  5. cat core/release/notes/Y2CD9X.md  -> two one-line notes, tester+logbook-summary on haiku, pr fix agent on the session's model  6. live (optional): in a Claude Code session run /jaira-role-tester and check the status line shows haiku while it runs"
---

# Rollen, die weder Code ändern noch urteilen, laufen auf haiku

## Definition of Done

- [x] tester- und logbook-summary-Rolle starten auf haiku; pr-Rolle, critique/review/bughunt bleiben auf dem Modell der Session (Entscheidung 21.11.: Variante 2)
  proof: core/role/role_test.go TestRoleModels; model: haiku only in core/role/builtin/jaira-role-{tester,logbook-summary}/SKILL.md
- [x] Ein PR-Kommentar, der eine Codeänderung verlangt, geht an einen Worker auf dem normalen Modell
  proof: core/role/builtin/jaira-role-pr/SKILL.md step 1 (Agent without model inherits the pr role's model = the session's; pr has no model: line)
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/Y2CD9X.md; core/shipped/shipped.sha256 (+3 lines)
- [x] pr-Rolle hat kein model:-Feld; tester und logbook-summary behalten haiku (TestRoleModels)
  proof: core/role/role_test.go TestRoleModels (haiku map: tester, logbook-summary); core/role/builtin/jaira-role-pr/SKILL.md header has no model:

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-21 21:04 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock
- **2026-11-21 21:05 · Alexander Sacharov** — decision: mechanism = model: haiku in the SKILL.md frontmatter of pr, tester, logbook-summary. Claude Code honours model: in a skill header for the turn the skill runs, so it holds however the role is started (spawn.sh tab, Skill tool, Agent). No board setting, no spawn.sh flag: smallest change, nothing new to configure.
- **2026-11-21 21:06 · Alexander Sacharov** — decision: Alex 21.11.: Grenze bestätigt — haiku nur für Rollen, die weder Code ändern noch urteilen (pr ohne Codeänderung, tester, CI-Warten, logbook-summary); critique, review, bughunt bleiben strong.
- **2026-11-21 21:07 · Alexander Sacharov** — decision: code-changing review point goes to Agent(model: opus), named explicitly because a subagent would inherit haiku from the skill. opus = what this account's workers run on; sonnet is only the cockpit chat-model. Not via spawn.sh worker: the ticket is usually in done/signoff by PR time, a full worker run would walk lanes.
- **2026-11-21 21:07 · Alexander Sacharov** — critique: release note line 1 — replace 'run jaira roles install ... remove the line from the installed SKILL.md' with 'jaira update --roles installs it'. Rest stands: model: header is the smallest mechanism (note above); opus named explicitly in pr step 1 because Agent would inherit haiku. Checked: teamlead/dispatcher start pr/tester as their own session, never inline, so haiku does not leak into a strong session.
- **2026-11-21 21:08 · Alexander Sacharov** — testing round 1: vet, shipped-history --check, check-ticket-commits vs upstream/master all exit 0 (vs origin/master it exits 1: the fork's master is stale, not this branch). roles install --into scratch: model: haiku in exactly pr/tester/logbook-summary, none in lane/bughunt; install into /proc/nope exits 1. Full go test ./... running; bughunt started.
- **2026-11-21 21:10 · Alexander Sacharov** — bughunt: done at 7bb4c47db - 0 findings in 6 modules (role package core/role read with role.go, install.go, target.go; the three SKILL.md files; spawn.sh pr lane; shipped list), first hunt; concurrency, error paths, input, permissions checked; files: core/release/notes/Y2CD9X.md core/role/builtin/jaira-role-logbook-summary/SKILL.md core/role/builtin/jaira-role-pr/SKILL.md core/role/builtin/jaira-role-tester/SKILL.md core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-21 21:10 · Alexander Sacharov** — known limits (bughunt, unverifiable here): 1) that Claude Code switches to haiku on model: in the SKILL.md header is harness behaviour, not testable from Go; 2) --agent copilot installs the same header and the pr role's Agent(model: opus) step, which Copilot likely ignores/cannot run; 3) an account without opus makes the pr role's fix-agent call fail.
- **2026-11-21 21:13 · Alexander Sacharov** — Given back from signoff: der Code-Fix-Agent der pr-Rolle darf kein festes opus nennen — auf einem Konto ohne opus scheitert er. Modell weglassen bzw. inherit (das Modell der Session), damit er auf jedem Konto läuft. Rest der Review-Notes bleibt note.
- **2026-11-21 21:13 · Alexander Sacharov** — decision (Alex give-back 21.11.): pr role step 1 names no model on the fix agent; Agent without model runs on the session's model, so it works on any account. Claim 'inherits haiku from this role' removed from the skill text. Release note line reworded to 'the session's model'.
- **2026-11-21 21:14 · Alexander Sacharov** — critique round 3: change is the person's decision verbatim (no model on the fix agent). Known limit, unverifiable from Go: if Claude Code's 'inherit' takes the model the skill switched to for the turn, the fix agent runs on haiku, not the session's default. Live check: run the pr role on a PR with a code-changing thread and look at the agent's model.
- **2026-11-21 21:16 · Alexander Sacharov** — testing: note: pre-existing, not this change: two parallel 'jaira roles install --into <same dir>' — one reports a just-written file as 'edited here', skips it and exits 3; files on disk end correct. roles install code is untouched by this branch.
- **2026-11-21 21:16 · Alexander Sacharov** — testing round 3: suite of changed files and function pass done, test-verdict written; waiting for bughunt.
- **2026-11-21 21:17 · Alexander Sacharov** — bughunt: change: blocker: core/role/builtin/jaira-role-pr/SKILL.md:488-493 - class: a subagent launched without model: from a skill whose header says model: haiku inherits haiku, not the session's model; the step that says to name no model therefore hands every code-changing review fix to haiku, which line 485-486 of the same step says is not enough to change code; the release note core/release/notes/Y2CD9X.md:2 promises the session's model and is wrong for the same reason; 6f57eee71 reversed the earlier text that said this correctly; example: session on opus with no default subagent model set, probe skill with model: haiku calls Agent(subagent_type: general-purpose) with no model, and the child reports claude-haiku-5-5; ran: claude -p /probe in a scratch repo with .claude/skills/probe/SKILL.md (model: haiku), result PARENT=claude-opus-5-5 CHILD=claude-haiku-5-5
- **2026-11-21 21:17 · Alexander Sacharov** — bughunt: done at 6f57eee71 - 1 finding (1 change, blocker) in 3 modules, follow-up hunt since 7bb4c47db (diff only, no new files); concurrency, error paths, input, permissions checked; files: core/release/notes/Y2CD9X.md core/role/builtin/jaira-role-logbook-summary/SKILL.md core/role/builtin/jaira-role-pr/SKILL.md core/role/builtin/jaira-role-tester/SKILL.md core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-21 22:01 · Alexander Sacharov** — question for Alex (give-back premise disproved by bughunt 21:17: Agent without model inherits haiku from the skill): how does the pr role's code-fix agent get a strong model on every account? 1) (Recommended) name opus, fall back to sonnet if opus is refused, never haiku — keeps pr on haiku, works on accounts without opus. 2) drop model: haiku from the pr role entirely — pr runs on the session model, haiku stays only on tester and logbook-summary; simplest, pr gets dearer. 3) start the fix as its own session (spawn.sh worker) — session model, but heavy and walks lanes on a ticket already in signoff/done.
- **2026-11-21 22:01 · Alexander Sacharov** — decision: Variante 2 (Teamlead 21.11. nachts, Alex' Regel 'Ränder wachsen → kleineres Design'): model: haiku aus der pr-Rolle entfernen; pr läuft auf dem Modell der Session, Code-Fix-Agent ohne model erbt es damit richtig. haiku bleibt bei tester und logbook-summary. Keine Opus/Sonnet-Fallback-Logik.
- **2026-11-21 22:17 · Alexander Sacharov** — in-progress round 4: variant 2 done (6822ec61b): model: haiku removed from pr role, step 1 says why, TestRoleModels pins tester+logbook-summary only, note line 1 reworded. Merged upstream/master (v0.6.0 shipped meanwhile, shipped-history went red); pr hash re-recorded after the merge. Foreign uncommitted 0.6.0 'jaira update' edits (.jaira/lanes/*, AGENTS.md, CLAUDE.md) found in this worktree: set aside for the merge and put back, not committed by me. TestReapTearsDownWhatTheSetupBuilt/hangs failed once under full-suite load (20s vs 1s bound), passes alone; reap.sh untouched by this branch.
- **2026-11-21 22:19 · Alexander Sacharov** — testing round 4: tests and function pass done, test-verdict written; waiting for bughunt.
- **2026-11-21 22:24 · Alexander Sacharov** — bughunt: done at b5ec40701 - 0 findings in 9 modules, follow-up hunt since 6f57eee71 (diff only for the 6 files read before; new and read whole: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/optimize.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md as their uncommitted jaira-update diff); concurrency, error paths, input, permissions checked; ran: go test ./core/role/ ./core/shipped/ (ok); files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/optimize.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/release/notes/Y2CD9X.md core/role/builtin/jaira-role-logbook-summary/SKILL.md core/role/builtin/jaira-role-pr/SKILL.md core/role/builtin/jaira-role-tester/SKILL.md core/role/role_test.go core/shipped/shipped.sha256
- **2026-11-21 22:40 · Alexander Sacharov** — review handed to a fresh agent via run-lane.sh (herdr tab)
- **2026-11-21 22:41 · Alexander Sacharov** — review: run-lane.sh failed twice ('claude did not come up in t19s/t19t', exit 1, tabs t19s and t19t may be left open); review handed to a fresh subagent instead
