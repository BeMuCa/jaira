---
id: 01M4VGHHEKBXMV9120AD3SJMV8
title: "Der jaira-Skill sagt, wie ein PR mit seinen Tickets entsteht, und verweist auf die PR-Rolle"
status: critique
ready: true
creator: BeMuCa
goal: "Ein Agent, der fuer ein jaira-Repo einen PR vorbereitet, legt die angenommenen Tickets auf dem Branch ins Logbuch, bevor er pusht, weil der jaira-Skill es ihm sagt."
context: |-
  Am 2026-10-13 hat Claude PR #74 geoeffnet, waehrend dessen Tickets noch in signoff standen. Alex (sashasoft90) setzte 'Changes requested': 'Hi berk schliss mal tickets ab'.
  Die Regel steht in der PR-Rolle core/role/builtin/jaira-role-pr/SKILL.md ('File the finished tickets away': pull, logbook, ein chore-Commit pro Ticket, vor dem Push).
  Der jaira-Skill (.claude/skills/jaira/SKILL.md) erwaehnt die PR-Rolle nirgends und sagt das Gegenteil: 'An accepted ticket leaves the board once its work is pushed ... the logbook waits on it' (SKILL.md:427 auf feat/milestone-roadmap). Der Agent-Block in CLAUDE.md nennt die Rolle auch nicht.
  Optionale Schritte erklaert der Skill schon ('## Optional steps': ein nicht angehakter Schritt wird beim Hineinschieben abgelehnt) - das muss nicht dazu.
  Verwandt: PWAF22 (Skills nach einem Upgrade nachziehen) - ein Skill, der nicht neu installiert wird, hilft nicht.
definition-of-done: "SKILL.md hat einen Abschnitt zum PR: angenommene Tickets werden auf dem PR-Branch mit jaira pull, jaira logbook und einem chore-Commit pro Ticket abgelegt, bevor gepusht wird; er verweist auf /jaira-role-pr"
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-10-13T20:51:24Z
updated-at: 2026-10-22T08:09:00Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-4188296
claimed-at: 2026-10-22T07:55:16Z
outcome-what: "The jaira skill (.claude/skills/jaira/SKILL.md) gets a section 'Making a pull request': accepted tickets are filed away on the PR branch with jaira pull, jaira logbook and one chore commit each before the push, and it points to /jaira-role-pr. The 'logbook waits on the push' paragraph now separates a direct push from a PR. The agent block (core/board/announce.go) gets one line naming /jaira-role-pr, guarded by TestAgentNoteNamesThePullRequestRole; CLAUDE.md and AGENTS.md regenerated with a build of this branch; release note core/release/notes/3SJMV8.md."
outcome-why: "PR #74 was opened while its tickets still sat in signoff and got 'changes requested'. The filing rule existed only in the PR role; the jaira skill said the opposite (log only after the push) and neither the skill nor the agent block named the role."
outcome-resolves: "DoD 1: SKILL.md section 'Making a pull request'. DoD 2: the paragraph in 'Taking a ticket off the board' distinguishes direct push and PR. DoD 3: announce.go bullet + TestAgentNoteNamesThePullRequestRole + regenerated CLAUDE.md/AGENTS.md. DoD 4: a fresh subagent told only to make a PR described filing before the push; result in a note on the ticket."
executed-by: claude-sonnet
review-summary: |-
  .claude/skills/jaira/SKILL.md:435-441 says the rule three times and calls the PR case "the exception, and only apparently", which contradicts itself; "the logbook waits until the code is pushed" (436) repeats "so the logbook waits on it" one sentence later. Instead: make 435-437 read "this board exists to prevent. That is the order for a direct push to the default branch. Do not push on your own initiative; the push is the user's call, so the logbook waits on it." and replace the paragraph 438-441 with "On a pull request the order turns round, for the same reason: the code reaches the default branch only when the pull request merges, so file the accepted tickets on its branch *before* pushing it — they arrive in that same merge. See **Making a pull request**."
  .claude/skills/jaira/SKILL.md:476-477 "The pull request itself is the maintainer's to open, not yours." is this repository's own house rule (CLAUDE.md:207-224, behind jaira:local) copied into a skill every jaira repository reads, and it contradicts the role the same paragraph points to, which opens the PR when a person invoked it (core/role/builtin/jaira-role-pr/SKILL.md:10-13, 269-288). Replace that sentence with the role's rule: "It opens the pull request only when a person invoked it; started by an agent, it pushes and hands the create line back for a person to run."
  .claude/skills/jaira/SKILL.md:469, 470 and 473 write `<ID>` for the value 467-468 call `<handle>`; the skill uses `<handle>` 40 times, and `<ID>` elsewhere (207, 252, 254) means a milestone's ID. Use `<handle>` in all three: it matches the pathspec (a ticket file ends in <handle>-<slug>.md) and is what commit subjects carry (fix(A3K9QP)).
  core/board/announce.go:145 (and its copies in CLAUDE.md and AGENTS.md) says what a Claude Code skill does rather than what the reader does, is the only entry in the block that names a slash command (an AGENTS.md reader without `jaira roles install` has no /jaira-role-pr), and leaves out "on the branch". Replace with "- a pull request: file accepted tickets on its branch before the push (`/jaira-role-pr`)" and regenerate CLAUDE.md/AGENTS.md as before (scratch-HOME `jaira update` with a build of this branch); TestAgentNoteNamesThePullRequestRole keeps passing because the literal stays.
round: 2
---

# Der jaira-Skill sagt, wie ein PR mit seinen Tickets entsteht, und verweist auf die PR-Rolle

## Definition of Done

- [x] SKILL.md hat einen Abschnitt zum PR: angenommene Tickets werden auf dem PR-Branch mit jaira pull, jaira logbook und einem chore-Commit pro Ticket abgelegt, bevor gepusht wird; er verweist auf /jaira-role-pr
  proof: .claude/skills/jaira/SKILL.md:459 (section Making a pull request)
- [x] Der Satz 'the logbook waits on [the push]' (SKILL.md:427) widerspricht dem nicht mehr: er unterscheidet direkten Push und PR
  proof: .claude/skills/jaira/SKILL.md:439 (direct push vs PR)
- [x] Der Agent-Block (core/board/announce.go) nennt /jaira-role-pr fuer PRs in einer Zeile; CLAUDE.md und AGENTS.md neu erzeugt; ein Test prueft die Zeile
  proof: core/board/announce.go:145; TestAgentNoteNamesThePullRequestRole; CLAUDE.md/AGENTS.md regenerated
- [x] Ein frischer Subagent mit dem Auftrag 'mach einen PR fuer Ticket X' beschreibt das Ablegen vor dem Push, ohne dass man es ihm sagt; das Ergebnis steht als Notiz am Ticket
  proof: jaira note on 3SJMV8, 2026-10-22: fresh-subagent PR experiment

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-22 07:59 · BeMuCa** — Items 1-3 done in worktree, uncommitted. 'jaira update' with scratch HOME touched only CLAUDE.md and AGENTS.md (one bullet each). Test needs the literal with closing backtick: `/jaira-role-pr`.
- **2026-10-22 08:00 · BeMuCa** — DoD 4 experiment, 2026-10-22: a fresh sonnet subagent was told only 'make a PR for ticket 3SJMV8 in this worktree; read CLAUDE.md and the jaira skill first; describe your steps, do not run git push, gh, or jaira move/logbook'. Its plan, unprompted: commit code with the ticket -> list the handles the branch's commits name -> for each ticket in done: jaira pull, jaira logbook, git add -A -- ':(glob).jaira/**/*<ID>-*.md', one chore commit -> only then git push -u origin HEAD -> hand the gh pr create line back. It named SKILL.md 'Making a pull request', the new CLAUDE.md bullet and the PR role (which the skill points to) as its sources. Caveat: the maintainer's auto-memory index, outside the repo and visible to every subagent, also mentions the PR filing order, so the run was not fully blind; the agent cited the repo files, not memory, for that step. It also noted 3SJMV8 itself was not in done yet, so filing would be skipped for it - correct.
- **2026-10-22 08:06 · BeMuCa** — critique (pass 1, whole diff): four findings, each with its fix, also in review-summary.
.claude/skills/jaira/SKILL.md:435-441 says the rule three times and calls the PR case "the exception, and only apparently", which contradicts itself; "the logbook waits until the code is pushed" (436) repeats "so the logbook waits on it" one sentence later. Instead: make 435-437 read "this board exists to prevent. That is the order for a direct push to the default branch. Do not push on your own initiative; the push is the user's call, so the logbook waits on it." and replace the paragraph 438-441 with "On a pull request the order turns round, for the same reason: the code reaches the default branch only when the pull request merges, so file the accepted tickets on its branch *before* pushing it — they arrive in that same merge. See **Making a pull request**."
.claude/skills/jaira/SKILL.md:476-477 "The pull request itself is the maintainer's to open, not yours." is this repository's own house rule (CLAUDE.md:207-224, behind jaira:local) copied into a skill every jaira repository reads, and it contradicts the role the same paragraph points to, which opens the PR when a person invoked it (core/role/builtin/jaira-role-pr/SKILL.md:10-13, 269-288). Replace that sentence with the role's rule: "It opens the pull request only when a person invoked it; started by an agent, it pushes and hands the create line back for a person to run."
.claude/skills/jaira/SKILL.md:469, 470 and 473 write `<ID>` for the value 467-468 call `<handle>`; the skill uses `<handle>` 40 times, and `<ID>` elsewhere (207, 252, 254) means a milestone's ID. Use `<handle>` in all three: it matches the pathspec (a ticket file ends in <handle>-<slug>.md) and is what commit subjects carry (fix(A3K9QP)).
core/board/announce.go:145 (and its copies in CLAUDE.md and AGENTS.md) says what a Claude Code skill does rather than what the reader does, is the only entry in the block that names a slash command (an AGENTS.md reader without `jaira roles install` has no /jaira-role-pr), and leaves out "on the branch". Replace with "- a pull request: file accepted tickets on its branch before the push (`/jaira-role-pr`)" and regenerate CLAUDE.md/AGENTS.md as before (scratch-HOME `jaira update` with a build of this branch); TestAgentNoteNamesThePullRequestRole keeps passing because the literal stays.
Left standing: the four-line command block in SKILL.md 'Making a pull request' duplicates core/role/builtin/jaira-role-pr/SKILL.md:186-191, but DoD 1 demands exactly those steps and the pathspec is what keeps an agent from adding .jaira/ whole; the details (which handles, no --all, pull for freshness) rightly stay only in the role. TestAgentNoteNamesThePullRequestRole matches the existing TestAgentNoteNames* pattern in core/board/announce_test.go. Checked and accurate: 'the tickets filed on the PR branch arrive in that same merge' - filing leaves the ref up (core/refsync/refsync.go:166-179), so the rewrite in finding 1 keeps that claim. Standing on the DoD 4 note (fresh-subagent experiment) for discoverability; these findings are wording and consistency only.
