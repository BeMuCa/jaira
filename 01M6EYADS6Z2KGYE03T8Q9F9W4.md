---
id: 01M6EYADS6Z2KGYE03T8Q9F9W4
title: "Der Installer erkennt den Coding-Agenten und legt die Rollen dort ab, wo er sie liest"
status: pre-process
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
updated-at: 2026-11-05T17:59:03Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: q9f9w4-driver
claimed-at: 2026-11-05T17:54:30Z
---

# Der Installer erkennt den Coding-Agenten und legt die Rollen dort ab, wo er sie liest

## Definition of Done

- [ ] core/role/target.go kennt die acht Agenten mit Projekt- und Home-Pfad; jaira roles install --agent <name> schreibt dorthin; ein Test je Agent
- [ ] install.sh und install.ps1 erkennen die vorhandenen Agenten an ihren Home-Ordnern, installieren fuer jeden erkannten, und schreiben bei keinem erkannten nichts ins Home, sondern einen Satz mit dem Nachhol-Befehl
- [ ] JAIRA_AGENT=<name,name> uebersteuert die Erkennung; JAIRA_NO_ROLES=1 bleibt
- [ ] README-Install-Abschnitt nennt die Erkennung und die Uebersteuerung
- [ ] Release-Notiz in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

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
