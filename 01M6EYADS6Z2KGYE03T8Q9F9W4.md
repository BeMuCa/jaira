---
id: 01M6EYADS6Z2KGYE03T8Q9F9W4
title: "Der Installer erkennt den Coding-Agenten und legt die Rollen dort ab, wo er sie liest"
status: review
ready: true
creator: BeMuCa
goal: "install.sh und install.ps1 erkennen, welche Agenten auf dem Rechner sind, und jaira roles install kennt neben claude und copilot auch codex, cursor, gemini-cli, opencode, pi und hermes mit ihren Pfaden; ohne erkannten Agenten wird nichts ins Home geschrieben und ein Satz sagt, wie man es nachholt."
context: |-
  Was schief ist: jaira roles install --global schreibt immer nach ~/.claude/skills (core/role/target.go:48), ob Claude Code da ist oder nicht; install.sh ruft es fest auf (Branch V38STG, Zeile 115). Ein PM ohne Coding-Agent bekommt einen toten Ordner, ein Codex- oder Cursor-Nutzer bekommt die Rollen an einer Stelle, die sein Agent nie liest (Wunsch 2026-11-02: 'Jaira sollte ueberall funktionieren').
  Was es schon gibt: eine Tabelle je Agent in core/role/target.go (claude, copilot) mit Projekt- und Home-Pfad und dem Flag, ob der Agent die Mods laedt; --into <dir> fuer unbekannte Agenten; JAIRA_NO_ROLES=1 als Opt-out im Installer.
  Wie andere es loesen (Vercel skills CLI, docs/agents.md, gelesen 2026-11-02): der Installer erkennt installierte Agenten an ihren Home-Ordnern und fragt nur, wenn keiner da ist; -a <agent> nennt ihn explizit. Pfade dort: Claude .claude/skills bzw. ~/.claude/skills; Codex und Pi .agents/skills bzw. ~/.agents/skills; Cursor .agents/skills bzw. ~/.cursor/skills; Gemini CLI .agents/skills bzw. ~/.gemini/skills; OpenCode .agents/skills bzw. ~/.config/opencode/skills; Copilot .agents/skills bzw. ~/.copilot/skills; Hermes .hermes/skills bzw. ~/.hermes/skills. Diese Pfade vor dem Eintragen gegen die Doku des jeweiligen Agenten pruefen, die Quellen widersprechen sich teils.
  Das Cockpit bleibt Claude-only (Claude-Code-Plugin, core/mod); bei jedem anderen Agenten wird es weggelassen, wie heute bei copilot.
  Nicht gewollt: ein eigener Skill-Store oder Symlinks; die Tabelle in target.go waechst, mehr nicht.
definition-of-done: "core/role/target.go kennt die acht Agenten mit Projekt- und Home-Pfad; jaira roles install --agent <name> schreibt dorthin; ein Test je Agent"
tags:
  - cli
  - docs
blocked-by: []
related:
  - 01M55S97BCFRNZMARY8WV38STG
commits: []
created-at: 2026-11-02T20:13:27Z
updated-at: 2026-11-05T18:11:54Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: q9f9w4-driver
claimed-at: 2026-11-05T17:54:30Z
outcome-what: "target.go knows eight agents (codex, pi, cursor, gemini-cli, opencode, hermes added) with project/home/also dirs; install.sh and install.ps1 detect agents by home folder and run roles install --agent X --global each, JAIRA_AGENT overrides, no agent writes nothing and says jaira works without one; README and release note updated"
outcome-why: "roles went to ~/.claude/skills whether or not Claude Code was there, and to a place other agents never read"
outcome-resolves: "DoD 1-5; tests written but not run, install.ps1 unrun"
review-summary: none
round: 3
review-gaps: "removed: install.ps1 Install-JairaRoles '$found = @()' (both branches assign it); roles_test.go TestRolesInstallAgentSkipsWhereItAlreadyReads guard 'scope.own != also' (no table entry lists its own dir as an Also dir, so always true). left: agentDirs in roles_test.go repeats target.go's table on purpose (independent expectation, folding it would make the tests tautological); install.sh keeps two lists (agent_dir case + AGENTS order) and install.ps1 a third copy of the table, documented at target.go:39-41 because the script cannot ask an older binary; nothing hot, no cost pass findings; no Go duplicate of the detection exists"
test-verdict: "pass: gofmt/vet/build/sh -n/dash -n rc=0; 8 agents verified global+project incl. Also-skip and cockpit only for claude; install.sh none/codex/all/JAIRA_AGENT/NO_ROLES verified in scratch HOME; README+release note format ok; go test NOT run, install.ps1 NOT run (read against install.sh only)"
---

# Der Installer erkennt den Coding-Agenten und legt die Rollen dort ab, wo er sie liest

## Definition of Done

- [x] core/role/target.go kennt die acht Agenten mit Projekt- und Home-Pfad; jaira roles install --agent <name> schreibt dorthin; ein Test je Agent
  proof: core/role/target.go agents table (8 entries); internal/cli/roles_test.go TestRolesInstallAgentWritesWhereItReads (t.Run per agent); go vet compiles it, not run (no go test in this lane)
- [x] install.sh und install.ps1 erkennen die vorhandenen Agenten an ihren Home-Ordnern, installieren fuer jeden erkannten, und schreiben bei keinem erkannten nichts ins Home, sondern einen Satz mit dem Nachhol-Befehl
  proof: scripts/install.sh install_roles + agent_dir; scripts/install.ps1 Install-JairaRoles (unrun, no pwsh); harness runs: none/codex/all/override in note
- [x] JAIRA_AGENT=<name,name> uebersteuert die Erkennung; JAIRA_NO_ROLES=1 bleibt
  proof: scripts/install.sh install_roles JAIRA_AGENT branch; JAIRA_NO_ROLES first; harness runs override and nores
- [x] README-Install-Abschnitt nennt die Erkennung und die Uebersteuerung
  proof: README.md ## Install paragraph
- [x] Release-Notiz in core/release/notes/<ID>.md
  proof: core/release/notes/Q9F9W4.md

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] test first, internal/cli/roles_test.go: TestRolesInstallAgentWritesWhereItReads, one t.Run per agent (claude, codex, pi, copilot, cursor, gemini-cli, opencode, hermes): --project and --global land in the dirs of the findings note; jaira-cockpit only for claude; replaces TestRolesInstallAgentCopilotWritesWhereCopilotReads
- [x] test first: TestRolesInstallAgentSkipsWhereItAlreadyReads, one t.Run per agent with Also dirs: roles pre-installed in each Also dir (project and global) -> exit 0, one 'already reads' line, own dir not created; keeps the copilot skip test's assertions
- [x] test first, core/role/target_test.go: agents sharing a Project or Global dir agree on Mods; claude is first; table order is claude, codex, pi, copilot, cursor, gemini-cli, opencode, hermes
- [x] test first: Installed() attributes a shared dir (~/.agents/skills) once, to codex; UpdateCommand for it returns 'jaira roles install --agent codex --global'
- [x] extend TestRolesInstallHelpNamesEveryAgentsDirectories and TestRolesInstallAgentRefusesUnknownAndInto to all eight names and dirs
- [x] core/role/target.go: add the six entries in the order above with Project/Global/AlsoProject/AlsoGlobal from the findings note, Mods false; extend the agents comment with each source URL, the order rule and the claude+codex duplicate limit
- [x] internal/cli/roles.go: rewrite the --agent paragraph of the Long help generically (every agent but claude leaves the cockpit out; Also dirs skip) instead of copilot-only prose
- [x] gofmt -l, go vet ./..., go build ./... (no go test in this lane)
- [x] scripts/install.sh: agent list in table order with its home folder (~/.claude, ~/.codex, ~/.pi, ~/.copilot, ~/.cursor, ~/.gemini, ~/.config/opencode, ~/.hermes); detect by folder; JAIRA_AGENT=a,b overrides (known names in table order, unknown named and skipped); JAIRA_NO_ROLES=1 first
- [x] install.sh: per agent run 'jaira roles install --agent X --global', print its summary or 'already reads' line; failure is non-fatal with the retry line; cockpit/~/.claude lines only for claude; none found -> nothing into HOME, one sentence: jaira works without an agent (e.g. jaira web), add roles later with 'jaira roles install --agent <name> --global'
- [x] update install.sh header comment; check with sh -n and dash -n; run it with HOME=throwaway dir: no agent folders, fake ~/.codex, JAIRA_AGENT=claude,bogus, JAIRA_NO_ROLES=1
- [x] scripts/install.ps1: same detection, order, JAIRA_AGENT, messages, under $env:USERPROFILE; header comment; no exit (irm|iex); pwsh not available here, so note it is unrun
- [x] README ## Install: replace the ~/.claude/skills sentence with detection by home folder, JAIRA_AGENT override, JAIRA_NO_ROLES, and the no-agent case (jaira web, roles later)
- [x] core/release/notes/Q9F9W4.md: one-line instructions: new --agent names; installer detects agents and JAIRA_AGENT; no agent -> nothing in home; jaira update --roles now also refreshes roles in the new agents' dirs
- [x] scan the diff for employer names and hostnames; gofmt -l, go vet ./..., go build ./...; commit with Q9F9W4 in the subject, ticket file alongside
- [x] AGENTS.md table names all eight agents plus the skip rule
- [x] target.go comment names both installer scripts
- [x] installers print the add-later line for only-unknown JAIRA_AGENT
- [x] AGENTS.md states the claude+codex double listing for cursor and opencode

## Progress
- **2026-11-05 17:53 · BeMuCa** — Berk, 2026-11-05: (1) ohne erkannten Agenten sagt der Installer ausdruecklich, dass jaira trotzdem nutzbar ist - z. B. ein PM ohne Agent mit 'jaira web' (Zusammenhang H319TV). (2) Stand: V38STG ist auf master gemergt (install.sh und install.ps1 haben JAIRA_NO_ROLES), dieses Ticket baut darauf auf. (3) Windows: install.ps1 bekommt dieselbe Erkennung wie install.sh.
- **2026-11-05 17:55 · BeMuCa** — Driver (q9f9w4-driver), before the plan: no decision left that only a person can take, so no conversational mode. Settled by default, open to disagreement: (a) keep copilot as target.go has it (.github/skills, ~/.copilot/skills, Also dirs) — the ticket's .agents/skills for copilot is the Vercel table, target.go's is from Copilot's own docs (E75R4Q); (b) name the Gemini entry gemini-cli as the goal does; (c) detection by each agent's own home folder (~/.claude, ~/.codex, ~/.cursor, ~/.gemini, ~/.config/opencode, ~/.copilot, ~/.pi, ~/.hermes), not by the shared ~/.agents; (d) the installer installs claude first so an agent that also reads ~/.claude/skills finds the roles there and skips via the existing Also rule; (e) the installer keeps --global only, as today. Pre-process checks the paths against each agent's own docs.
- **2026-11-05 17:59 · BeMuCa** — Pre-process: skills dirs per agent, from each agent's own docs (read 2026-11-05). Columns: write project / write home | also-read project | also-read home | detect.
claude: .claude/skills / ~/.claude/skills | - | - | ~/.claude — https://code.claude.com/docs/en/skills (verified; no .agents/skills)
codex: .agents/skills / ~/.agents/skills | - | - (also /etc/codex/skills, admin, not ours) | ~/.codex — https://learn.chatgpt.com/docs/build-skills (was developers.openai.com/codex/skills; verified; ~/.codex only via ~/.codex/config.toml on that page, no ~/.codex/skills)
pi: .agents/skills / ~/.agents/skills | .pi/skills | .pi/agent/skills | ~/.pi — https://github.com/badlogic/pi-mono/blob/main/packages/coding-agent/docs/skills.md names only .agents/skills + ~/.agents/skills; the .pi/skills and ~/.pi/agent/skills paths are from the same doc at npm @mariozechner/pi-coding-agent@0.73.1 (cdn.jsdelivr.net), current doc says only 'user or project skills directory' — PARTLY VERIFIED
copilot: .github/skills / ~/.copilot/skills | .claude/skills, .agents/skills | .agents/skills | ~/.copilot — https://docs.github.com/en/copilot/concepts/agents/about-agent-skills (verified; matches target.go, kept as is)
cursor: .agents/skills / ~/.cursor/skills | .cursor/skills, .claude/skills, .codex/skills | .agents/skills, .claude/skills, .codex/skills | ~/.cursor — https://cursor.com/docs/context/skills (verified)
gemini-cli: .agents/skills / ~/.gemini/skills | .gemini/skills | .agents/skills | ~/.gemini — https://geminicli.com/docs/cli/skills/ (verified; same-name skills: higher tier wins, .agents beats .gemini)
opencode: .agents/skills / ~/.config/opencode/skills | .opencode/skills, .claude/skills | .claude/skills, .agents/skills | ~/.config/opencode — https://opencode.ai/docs/skills/ (verified on Linux/mac; the Windows home path is NOT verified)
hermes: .hermes/skills / ~/.hermes/skills | .agents/skills | - (skills.external_dirs is user config, not ours) | ~/.hermes — https://hermes-agent.nousresearch.com/docs/user-guide/features/skills (verified; project dirs load only after 'hermes skills trust')
All write dirs equal the Vercel table (https://github.com/vercel-labs/skills/blob/main/docs/agents.md) except copilot project, where the driver's decision (a) keeps .github/skills; Copilot's own docs list both, so nothing contradicts it. Driver decisions (a)-(e) stand.
- **2026-11-05 17:59 · BeMuCa** — Pre-process, why the plan looks like this:
1. Order matters. The agents table order is also the installers' order: claude, codex, pi, copilot, cursor, gemini-cli, opencode, hermes. Writers of a shared dir come first, so the existing Also rule makes the later ones skip. With all eight present only ~/.claude/skills, ~/.agents/skills and ~/.hermes/skills get written.
2. Known limit, not fixable without symlinks (ruled out): with claude AND codex, both ~/.claude/skills and ~/.agents/skills must exist, and cursor and opencode read both, so they see each role twice. The Also rule only protects the agent being installed. Say so in the Agent comment; do not try to fix it.
3. codex and pi write the same dir. Second run reports 'unchanged' (its own target is not in its Also list) - harmless. Installed() and UpdateCommand() name the first table entry for a shared dir (codex), which is correct because the content is identical. A test pins that agents sharing a write dir agree on Mods (all false).
4. Detection stays in the scripts, not a new binary flag: the script comes from master, the binary is the latest release, and an older release would not know such a flag. The same gap hits new names: an older binary refuses '--agent codex' (exit 2, unknown agent). So each per-agent install is non-fatal: print its output and the retry line, go on with the next.
5. JAIRA_AGENT=a,b: known names run in table order (not the given order, see 1); an unknown name gets one line naming the known ones and is skipped, not fatal. JAIRA_NO_ROLES=1 still wins over everything.
6. Nothing detected: write nothing into HOME; one sentence: no coding agent found, jaira works without one (e.g. 'jaira web'), add the roles later with 'jaira roles install --agent <name> --global'.
7. Side effect: 'jaira update --roles' loops all agents via Installed(), so roles already sitting in e.g. ~/.agents/skills or project .agents/skills now get refreshed too. Worth a release-note line.
8. The 'active now ... ~/.claude/skills' and 'Cockpit loads in your next Claude Code session' lines print only when claude was installed.
9. Scripts have no Go tests. Check them with sh -n, dash -n and a run against a throwaway HOME with fake agent folders. pwsh is not installed here, so install.ps1 cannot be run locally - say so in the outcome.
- **2026-11-05 18:03 · BeMuCa** — In-progress: tests written first but NOT run (house rule: no go test; go vet compiles them). install.sh checked with sh -n, dash -n and a harness that sources only the roles section (scratch dir q9f9w4-scratch, HOME inside it, binary rebuilt from this branch): no agent folders -> nothing written plus the 'works without an agent (jaira web), add later' sentence; only ~/.codex -> ~/.agents/skills; all eight folders -> only ~/.claude/skills (+cockpit), ~/.agents/skills, ~/.hermes/skills written, copilot/cursor/gemini-cli/opencode print 'already reads' (opencode names ~/.claude/skills, because its Also list has .claude first), pi prints 21 unchanged; JAIRA_AGENT=claude,bogus -> bogus named and skipped, claude installed; JAIRA_NO_ROLES=1 with ~/.codex -> nothing written. install.ps1 was NOT run (no pwsh here) nor syntax-checked; it mirrors install.sh by hand. Adjacent, not fixed: with JAIRA_AGENT naming only unknown agents nothing is installed and no retry sentence is printed beyond the unknown-agent line; docs/AGENTS.md (linked from README) was not updated for the six new agents.
- **2026-11-05 18:05 · BeMuCa** — critique (pass 1): docs/AGENTS.md:276-288 still lists only claude and copilot and says "for any other agent, --into names its skills directory" - now false, and README.md:812 sends readers there for the eight agents; extend that table to all eight (--project/--global) and replace the copilot-only Also paragraph with the generic skip rule plus the claude+codex double-listing limit from target.go
core/role/target.go:39 "A further one is a further entry" is made false by this change: scripts/install.sh (agent_dir + AGENTS) and scripts/install.ps1 ($known) carry the same names, order and detect folders; name both scripts in that comment so the next agent lands in all three copies
scripts/install.sh install_roles and scripts/install.ps1 Install-JairaRoles: the JAIRA_AGENT guard around the no-agent sentence leaves JAIRA_AGENT=bogus with no add-later command (goal: a sentence says how to catch up); drop the guard, one branch less, and the sentence prints whenever nothing is installed
left alone: detection in the scripts, not the binary (pre-process note 17:59 point 4); claude+codex duplicate in cursor/opencode (point 2); copilot dirs kept (driver note 17:55 a); install.ps1 unrun and the opencode Windows path (in-progress note 18:03)
- **2026-11-05 18:06 · BeMuCa** — round 2: fixed the three critique findings: AGENTS.md table now lists all 8 agents + skip rule; target.go comment names install.sh/install.ps1; install.sh/ps1 print the add-later line also when JAIRA_AGENT has only unknown names. Ran: gofmt -l, go vet, go build, sh -n/dash -n; install.sh harness (scratch HOME) with JAIRA_AGENT=foo prints unknown + add-later line, codex and none cases unchanged. install.ps1 unrun (no pwsh).
- **2026-11-05 18:07 · BeMuCa** — critique (pass 2): findings 2 and 3 of pass 1 fixed in a69f5abdf. Finding 1 half done: docs/AGENTS.md now has all eight agents and the generic skip rule, but not the claude+codex double-listing limit pass 1 asked for, and its new sentence ('a second copy would register every role twice') reads as if the skip always prevents that. Add one sentence after the skip paragraph: with both claude and codex installed, cursor and opencode read ~/.claude/skills and ~/.agents/skills and list each role twice (known limit, see core/role/target.go:47-50). Nothing else reopened.
- **2026-11-05 18:07 · BeMuCa** — in-progress round 3: docs/AGENTS.md now states the claude+codex limit (cursor and opencode list each role twice) as the last sentence of the --agent paragraph, after the Copilot example so that sentence still follows its own lead-in.
- **2026-11-05 18:08 · BeMuCa** — critique (pass 3): the one finding of pass 2 (docs/AGENTS.md lacked the claude+codex double-listing limit) is fixed in 766139a52 - docs/AGENTS.md:295-297 now says cursor and opencode list each role twice with claude and codex both installed, matching core/role/target.go:47-50. Read only that change, as the lane prompt says; nothing else reopened. Standing on pass 1 and pass 2 notes for everything left alone. review-summary=none.
- **2026-11-05 18:11 · BeMuCa** — testing: go test NOT run (house rule), install.ps1 NOT run (no pwsh), read against install.sh line by line: same names, order, folders, messages. Gates: rebuilt J rc=0; gofmt -l (4 Go files) empty; go vet ./... rc=0; go build ./... rc=0; sh -n rc=0; dash -n rc=0. Eight agents, scratch HOME and scratch git repo with init: --global lands in .claude/.agents(codex,pi)/.copilot/.cursor/.gemini/.config/opencode/.hermes skills; --project in .claude, .agents (codex,pi,cursor,gemini-cli,opencode), .github, .hermes; jaira-cockpit only under claude (both scopes). Also-dir skip: pre-installed roles in each Also dir (project and global, 19 cases) -> rc=0, 'already reads' line, skill count unchanged. install.sh roles section (sed from '# Roles:' to end, binary = rebuilt J, scratch HOME): no folders -> nothing in HOME + 'no coding agent found ... jaira works without one (e.g. jaira web) ... jaira roles install --agent <name> --global'; .codex only -> ~/.agents/skills; all eight -> only .claude/.agents/.hermes skills written, pi 21 unchanged, copilot/cursor/gemini-cli/opencode 'already reads'; JAIRA_AGENT=cursor,bogus -> bogus skipped, cursor installed; JAIRA_AGENT=bogus -> unknown line + add-later sentence, nothing written; JAIRA_NO_ROLES=1 with .codex -> skipped line, nothing written. README and release note: 4 lines, each '- ', unwrapped, instructions.
