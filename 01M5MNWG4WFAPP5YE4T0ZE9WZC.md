---
id: 01M5MNWG4WFAPP5YE4T0ZE9WZC
title: "Eine Erlaubnisliste im Repo laesst Worker arbeiten, ohne bei jedem Befehl zu fragen"
status: critique
ready: true
creator: BeMuCa
goal: "Worker und Subagenten fragen nur bei wirklich kritischen Befehlen (rm, push, reset --hard, kill-server) nach; alles, was eine Lane braucht, ist im Repo erlaubt, und die Doku sagt, wo Verbots- und Erlaubnislisten liegen und welche Settings ein Worktree liest."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkt 5) und am 23.10. nachgeprueft: das Repo hat keine .claude/settings.json mit permissions; die User-Settings erlauben im Auto-Modus nur sechs jaira-Unterbefehle (create, list, next, show, claim, note). go build/vet/test, gofmt, git add/commit/worktree, jaira move/dod/set/pull/logbook fragen deshalb nach.
  Eine Verbotsliste, die ein Teamlead in einen Worktree legte, wurde nie gelesen: Worktrees lesen die Settings des Haupt-Checkouts - das steht nirgends.
  Verwandt: spawn.sh-Ticket zum Berechtigungsmodus (gleiche Runde). Der Skill 'fewer-permission-prompts' kann die Nur-Lese-Haelfte aus Transkripten vorschlagen.
definition-of-done: "Eine committete .claude/settings.json erlaubt die Befehle der Lanes (go build/vet/test/gofmt, git add/commit/status/diff/log/show/worktree/branch/fetch/stash, alle jaira-Unterbefehle, mkdir/cp/ls/cat/grep/find, tmux ohne kill-server) und verbietet rm -rf ausserhalb des Scratchpads, git push, git reset --hard, git filter-branch, gh pr create/merge, tmux kill-server"
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:25:47Z
updated-at: 2026-10-26T19:33:23Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-253301
claimed-at: 2026-10-26T19:10:55Z
outcome-what: "Added a committed .claude/settings.json. permissions.allow covers the lane commands: go build/vet/test, gofmt, the git commands lanes use, every jaira subcommand, tmux, the read/copy tools, rm -r /tmp/*, Edit(/**). permissions.deny: rm -rf/-fr, git push, git reset --hard, git filter-branch, gh pr create/merge, tmux kill-server. README: '### Permissions in a worktree' under With Claude Code (which settings a worktree reads, trust, autoMode, /permissions), plus a paragraph under Development on what this repo's file allows and denies and why rm -rf is denied even under /tmp. No release note: nothing a jaira user runs changes."
outcome-why: "Workers stopped to ask for routine lane commands (go build, git commit, jaira move/dod), because the repo carried no permission rules, and the user settings' autoMode.allow takes prose for the classifier, not patterns. A deny list a teamlead put into a worktree was never read, and nothing said which file a worktree reads."
outcome-resolves: "DoD 1: .claude/settings.json:3. Deny is checked before allow, so rm -rf is denied everywhere and rm -r /tmp/* is allowed. DoD 2: README.md:672, the worktree reads its own settings.json, but settings.local.json and trust come from the main checkout's root. DoD 3: a fresh headless sonnet worker on a scratch board ran the in-progress lane with 0 permission denials (note on ZE9WZC). Shell expansions such as the lane prompt's --commits \"$(git rev-parse HEAD)\" always ask; follow-up W9AVBE. Open: critique finding 1, deny or ask for git push / gh pr create, for the maintainer."
review-summary: ".claude/settings.json:41,44 deny `git push` and `gh pr create`, and a deny in this file beats every allow, settings.local.json included: /jaira-role-pr (core/role/builtin/jaira-role-pr/SKILL.md:212 \"git push -u origin HEAD\", :308 gh pr create when a person invoked it) and README.md:1121 / CLAUDE.md:223 (\"An agent pushes its branch\") can no longer run in this repo, not even from the maintainer's own session. The goal says these commands should ask (\"fragen nach\"), DoD 1 says forbid (\"verbietet\"). Maintainer decides: keep deny, or move both to `permissions.ask` and keep `gh pr merge` denied (question on the human move)."
round: 2
---

# Eine Erlaubnisliste im Repo laesst Worker arbeiten, ohne bei jedem Befehl zu fragen

## Definition of Done

- [x] Eine committete .claude/settings.json erlaubt die Befehle der Lanes (go build/vet/test/gofmt, git add/commit/status/diff/log/show/worktree/branch/fetch/stash, alle jaira-Unterbefehle, mkdir/cp/ls/cat/grep/find, tmux ohne kill-server) und verbietet rm -rf ausserhalb des Scratchpads, git push, git reset --hard, git filter-branch, gh pr create/merge, tmux kill-server
  proof: .claude/settings.json:3
- [x] docs/COMMANDS.md oder README sagt: Worktrees lesen die Settings des Haupt-Checkouts; eine Verbotsliste gehoert in dessen settings.local.json oder in die User-Settings
  proof: README.md:672
- [x] Ein frischer Worker auf einem Scratch-Board faehrt eine Lane ohne eine einzige Rueckfrage durch; das Ergebnis steht als Notiz am Ticket
  proof: jaira note on ZE9WZC: run 2 on scratch ticket FFDE20, 0 permission denials

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-26 19:12 · BeMuCa** — Driver: pre-process was passed by (next_lane went todo -> in-progress, the plan is empty). The lane worker gets its plan from the DoD plus the driver brief. Permission semantics are being checked against the Claude Code docs before writing; the user settings put their jaira rules under autoMode.allow, not permissions.allow. The dispatcher's spawn.sh drives herdr panes, not tmux, so the tmux rules only cover what a session types by hand.
- **2026-10-26 19:17 · BeMuCa** — Driver: the DoD as written cannot be met literally. Claude Code checks deny before allow and ignores how narrow a rule is, so 'deny rm -rf except in the scratchpad' cannot be written. Chosen instead: deny rm -rf / rm -fr everywhere and allow 'rm -r /tmp/*' for scratch boards. Edit(/**) was added beyond the DoD list, because default mode prompts on every file edit and DoD 3 could not pass without it. The user settings' autoMode.allow entries such as 'Bash(jaira create:*)' are prose for the auto-mode classifier, not tool patterns, and autoMode is never read from a project file. That is the likely reason the six jaira commands looked like the only ones allowed. Sources: code.claude.com/docs/en/permissions.md, settings.md, auto-mode-config.md.
- **2026-10-26 19:17 · BeMuCa** — Beyond the ticket list I added Bash(git ls-files *): core/role/builtin/jaira-role-lane/SKILL.md:108 tells a lane to run 'git ls-files --others --exclude-standard'. Also Bash(git config user.email) (exact, no wildcard): house rule 1 (maintainer's global instructions) requires it before the first commit; not in a lane prompt. No other lane prompt in .jaira/lanes/*.md names a shell command outside the rules.
- **2026-10-26 19:22 · BeMuCa** — DoD 3 experiment, scratch board without remote, sonnet via claude -p --permission-prompts none, rules passed with --settings from a copy at the repo root.
Run 1 (ticket R78J1W): 2 denials, both shell expansions, not missing rules.
 - jaira move ... --commits "$(git rev-parse HEAD)" : a substitution that is the whole argument. The in-progress lane prompt itself tells workers to type this.
 - a compound command with $? (echo rc=$?).
Probes: every listed command passes; jaira note "text $(...)" passes; echo "$(...)", whole-argument $(...), $? and ${PIPESTATUS[0]} are denied. No allow rule can approve those.
Run 2 (ticket FFDE20, prompt forbade $(...), backticks, $?, ${...}): 0 permission_denials, subtype success, 10 turns, committed 359106eb, ticket moved to critique.
Proven: the rule content (allow/deny list) works for a lane run via --settings.
NOT proven: the trust step of a committed file (-p ignores project allow rules in an untrusted folder, stderr said so); prompts an interactive session would show.
Follow-up: the lane text must stop prescribing $(git rev-parse HEAD) for --commits; the driver files it.
- **2026-10-26 19:28 · BeMuCa** — critique (pass 1, whole diff): five findings, the first needs the maintainer.
1. .claude/settings.json:41,44 deny git push and gh pr create. A deny in the committed file beats every allow, settings.local.json included, so /jaira-role-pr (SKILL.md:212 git push -u origin HEAD; :308 gh pr create when a person invoked it) and the rule at README.md:1138 / CLAUDE.md:223 that an agent pushes its branch stop working in this repo, the maintainer's own session included. The goal says ask ("fragen nach"), the DoD says forbid ("verbietet"). Decide: keep deny, or move push/create (and maybe rm -rf, reset --hard, filter-branch, kill-server) to permissions.ask; a headless worker with --permission-prompts none is refused either way (facts F14).
2. core/release/notes/ZE9WZC.md: delete. Nothing in the binary changed; no go:embed reads .claude/, and jaira update would tell every user about a file in a repo that is not theirs.
3. README.md:674-678: drop the five bullets that copy the JSON (they already miss git config user.email and rm -fr); point at .claude/settings.json and /permissions instead.
4. README.md:670-702: rewrap as paragraphs at 80 columns like README.md:661-668.
5. README.md:670: the part about this repo's own allow/deny lists belongs under ## Development (README.md:1100) next to the push rule; the worktree-settings part stays under With Claude Code.
Left alone, standing on the notes: rm -rf denied everywhere and Edit(/**) (19:17 driver note), git ls-files and git config user.email (19:17 note), the $(...) lane text split off to W9AVBE (19:22 note). Also left alone: sed/cp/find and git branch/worktree/stash wildcards reach past the project or drop work (sed -i, find -delete, branch -D, worktree remove --force); the README already says the deny list is not a sandbox, and the DoD names these commands.
- **2026-10-26 19:29 · BeMuCa** — Driver: critique findings 2-5 (delete the release note, README points to the JSON instead of repeating it, rewrap at 80 columns, repo-specific part under Development) go back to in-progress now. They do not depend on the deny-or-ask decision in finding 1. Finding 1 is held for the maintainer: the next critique pass sends it to human with the critic's question, so he is asked once, with only that decision left. Until he answers, the deny list stays as the DoD says.
- **2026-10-26 19:30 · BeMuCa** — Second pass: permissions docs split by audience. General part (worktree reads, trust, autoMode) stays under With Claude Code as 'Permissions in a worktree'; what this repo's settings.json allows/denies (rm -rf under /tmp, not a sandbox) moved next to the PR paragraph under Development. Allow/deny bullets dropped, the JSON is the list. Release note deleted: nothing the binary embeds reads .claude/. Finding 1 (deny vs ask for push/pr create) left for the maintainer; settings.json untouched.
