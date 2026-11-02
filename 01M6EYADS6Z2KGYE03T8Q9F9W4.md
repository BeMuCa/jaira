---
id: 01M6EYADS6Z2KGYE03T8Q9F9W4
title: "Der Installer erkennt den Coding-Agenten und legt die Rollen dort ab, wo er sie liest"
status: backlog
ready: false
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
updated-at: 2026-11-02T20:13:27Z
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
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

