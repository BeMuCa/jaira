---
id: 01M7ZYRRKJ8Q77ZZ3B8WY2CD9X
title: "Rollen, die weder Code ändern noch urteilen, laufen auf haiku"
status: critique
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
commits: []
created-at: 2026-11-21T21:04:04Z
updated-at: 2026-11-21T22:17:59Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-3504002
claimed-at: 2026-11-21T21:13:37Z
outcome-what: "tester and logbook-summary roles carry model: haiku; the pr role carries none and hands code-changing review points to an agent that inherits the session's model"
outcome-why: "pushing and running suites need no strong model, but a model: line on the pr role would be inherited by its code-fix agent"
outcome-resolves: "test and summary runs cost haiku prices; critique, review, bughunt and code fixes stay on the session model"
review-summary: "none"
round: 4
back-from: testing
back-by: lane
back-count: 1
review-gaps: none
test-verdict: "pass: round 3 tests of changed files green with -race (core/role, core/shipped, core/release RC=0), shipped-history --check RC=0, vet RC=0; DoD 1-3 checked in tree (SKILL.md:484-491 names no model, note line 2 says session's model, new pr hash in shipped.sha256); built binary, roles install --into scratch writes the new pr text and model: haiku in exactly pr/tester/logbook-summary; install into /proc/nope exits 1; two parallel installs: files correct, one exits 3 on a race (pre-existing, note)"
hunt-verdict: "fail: bughunt done at 6f57eee71, 1 change blocker: Agent without model inside a model: haiku skill runs on haiku (measured: parent opus, child haiku), so pr step 1 and release note line 2 are wrong"
review-verdict: "Diff meets all three DoD items: exactly pr/tester/logbook-summary carry model: haiku and TestRoleModels fails if any judging role gains a model line; a code-changing PR point goes to an opus worker; release note and shipped.sha256 are present. Checked the other paths in jaira-role-pr that could change code (none-mode merge, conflicts) - it aborts conflicts and stops, so haiku never resolves code there. tester only reads logs. No blocker; three notes about robustness of the opus hand-off."
review-round: 2
review-verdict-gaps: |-
  note: jaira-role-pr step 1 hard-codes model "opus"; an account or agent host without opus fails that call - name it as "the strongest model you have, never haiku" or fall back to sonnet if opus is refused
  note: jaira-role-pr step 1 pushes the subagent's commit without checking it - add one line: confirm the printed sha exists on HEAD and its subject names <id> before pushing
  note: model: haiku is installed verbatim for --agent copilot too, where the header and the Agent(model) step likely mean nothing; harmless, but untested (bughunt known limit)
review-check: "1. cd into this worktree  2. go test ./core/role -run TestRoleModels -count=1  -> ok  3. head -5 core/role/builtin/jaira-role-pr/SKILL.md -> line 'model: haiku' under description; same for jaira-role-tester and jaira-role-logbook-summary  4. grep -l '^model:' core/role/builtin/*/SKILL.md -> exactly those three files, no jaira-role-lane, no jaira-role-bughunt  5. sed -n 484,492p core/role/builtin/jaira-role-pr/SKILL.md -> step 1 sends code-changing points to Agent with model opus  6. Live check (optional): in a fresh claude session run /jaira-role-tester and look at the model shown in the status line -> haiku for that turn"
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
