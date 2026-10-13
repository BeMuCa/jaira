---
id: 01M4V3RJKTPTBT43QAWCE75R4Q
title: "jaira roles install --agent copilot legt die Rollen dort ab, wo Copilot sie liest"
status: testing
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "'jaira roles install --agent <name>' installiert die Rollen in das Skills-Verzeichnis des genannten Agenten; Copilot ist der erste Agent neben claude (Default), weitere kommen als eine Zeile in einer Tabelle dazu"
context: |-
  Was heute falsch ist: 'jaira roles install' schreibt nur nach .claude/skills (--project) oder ~/.claude/skills (--global). Wer mit GitHub Copilot arbeitet, bekommt die Rollen (jaira-teamlead, jaira-worker, jaira-role-* ...) nicht, ausser er kennt --into und das richtige Verzeichnis.
  Ausloeser: Alex am 13.10.: 'wir stellen alle Skills fuer Claude auf, sie muessen auch fuer andere Agenten installiert werden', Copilot zuerst.
  Entschieden (Alex, 13.10.): ein Flag --agent mit Default claude, Copilot als erster weiterer Eintrag; Codex/Cursor spaeter je eine Tabellenzeile, nicht in diesem Ticket.
  Code: core/role/target.go (skillsDir = .claude/skills, ProjectTarget, GlobalTarget); internal/cli/roles.go:87-117 (--project/--global/--into).
  Wichtige Falle, steht in target.go:8-16: liest ein Agent mehrere Skills-Verzeichnisse, ist dieselbe Rolle doppelt registriert. Falls Copilot auch .claude/skills liest, darf --agent copilot keine zweite Kopie danebenlegen - dann sagt der Befehl, dass schon installiert ist.
  ZUERST klaeren (pre-process, mit /jaira-role-research): welche Verzeichnisse liest Copilot heute fuer Skills - VS Code Agent Mode, Copilot CLI, Copilot Coding Agent - projektweit und persoenlich (Kandidaten: .github/skills, ~/.copilot/skills, .claude/skills); und versteht Copilot das SKILL.md-Format mit Frontmatter name/description so, wie die Rollen es schreiben. Quellen zitieren, Versionsstand nennen.
  Nicht Teil: Copilot-Instruktionsdatei (.github/copilot-instructions.md) - der Agent-Block geht schon nach AGENTS.md, das Copilot liest; erst pruefen, ob das reicht.
definition-of-done: "Eine Note mit Quellen sagt, welche Skills-Verzeichnisse Copilot (VS Code, CLI, Coding Agent) projektweit und global liest und ob es SKILL.md mit name/description versteht"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-13T17:08:03Z
updated-at: 2026-10-13T17:15:56Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-201739
claimed-at: 2026-10-13T17:08:32Z
outcome-what: "'jaira roles install' hat --agent (Default claude, neu copilot): eine Agent-Tabelle in core/role/target.go nennt je Agent das Projekt- und das globale Skills-Verzeichnis und die weiteren, die er liest; copilot schreibt nach .github/skills bzw. ~/.copilot/skills ohne das Cockpit-Mod, schreibt nichts wenn die Rollen schon in .claude/skills/.agents/skills liegen, unbekannter Agent und --agent mit --into enden mit Exit 2"
outcome-why: "Wer mit GitHub Copilot arbeitet, bekam die Rollen nur ueber --into und musste das Verzeichnis selbst kennen"
outcome-resolves: "Alle fuenf DoD-Punkte: Recherche-Note mit Quellen, Zielverzeichnisse, keine Doppelkopie, Usage-Fehler, Doku und Release-Note"
review-summary: "none"
round: 2
review-gaps: "core/role/target.go: Tabelle Agents -> agents (nur paketintern gelesen; CLI geht ueber AgentNamed/AgentNames). Gesucht nach vorhandenem Helfer fuer 'unbekannter Wert, bekannte sind ...' - keiner, die strings.Join-Stellen in gitref/store sind Fehlertexte eigener Art. Kosten: InstalledIn liest Builtins einmal pro Aufruf, nur auf dem einmaligen install-Pfad - gelassen."
---

# jaira roles install --agent copilot legt die Rollen dort ab, wo Copilot sie liest

## Definition of Done

- [x] Eine Note mit Quellen sagt, welche Skills-Verzeichnisse Copilot (VS Code, CLI, Coding Agent) projektweit und global liest und ob es SKILL.md mit name/description versteht
  proof: jaira note (Recherche, 13.10.2026) mit den vier Quell-URLs
- [x] 'jaira roles install --agent copilot --project' und '--global' schreiben in das Verzeichnis aus dieser Note; ohne --agent bleibt alles wie heute (claude)
  proof: core/role/target.go:50; TestRolesInstallAgentCopilotWritesWhereCopilotReads (internal/cli/roles_test.go:294); bestehende TestRolesInstall*Project/Global unveraendert gruen
- [x] Liest Copilot ein Verzeichnis, in dem die Rollen schon liegen, legt --agent copilot keine zweite Kopie an und sagt das in einer Zeile
  proof: internal/cli/roles.go:164; TestRolesInstallAgentCopilotSkipsWhereRolesAlreadyAre (internal/cli/roles_test.go:323)
- [x] Ein unbekannter --agent-Wert bricht mit Exit-Code fuer Usage ab und nennt die bekannten Agenten; --agent zusammen mit --into wird abgelehnt
  proof: internal/cli/roles.go:121 und :130; TestRolesInstallAgentRefusesUnknownAndInto (internal/cli/roles_test.go:341)
- [x] jaira roles install --help, README und docs/AGENTS.md nennen --agent; Release-Note als Datei core/release/notes/<Handle>.md
  proof: internal/cli/roles.go Long-Text + --agent-Flag; README.md:661; docs/AGENTS.md Tabelle Agent/--project/--global; core/release/notes/E75R4Q.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-13 17:08 · Alexander Sacharov** — parent session: jaira-6a
- **2026-10-13 17:09 · Alexander Sacharov** — Recherche (Stand 13.10.2026, Quellen): Copilot liest Skills projektweit aus .github/skills, .claude/skills und .agents/skills - gleich fuer VS Code, Copilot CLI und Coding/Cloud Agent. Persoenlich: CLI und Cloud Agent lesen ~/.copilot/skills und ~/.agents/skills (NICHT ~/.claude/skills); VS Code liest zusaetzlich ~/.claude/skills. SKILL.md mit YAML-Frontmatter name (lowercase, Bindestriche, = Verzeichnisname, max 64) und description (max 1024) ist genau das Format, das die Rollen schreiben. Doppelte Namen: VS Code zeigt beide Skills an, GitHub-Doku sagt nichts. Quellen: https://docs.github.com/en/copilot/how-tos/copilot-cli/customize-copilot/create-skills ; https://docs.github.com/en/copilot/how-tos/use-copilot-agents/cloud-agent/create-skills ; https://code.visualstudio.com/docs/agent-customization/agent-skills ; https://learn.microsoft.com/en-us/visualstudio/ide/copilot-agent-skills
- **2026-10-13 17:09 · Alexander Sacharov** — decision: Ziel fuer --agent copilot: --project -> .github/skills, --global -> ~/.copilot/skills (die Verzeichnisse, die alle drei Copilot-Oberflaechen lesen und die Copilot selbst empfiehlt). Keine zweite Kopie: liegt schon eine jaira-Rolle in einem anderen Verzeichnis, das ALLE Copilot-Oberflaechen lesen (projekt: .claude/skills, .agents/skills; global: ~/.agents/skills), schreibt der Befehl nichts und sagt in einer Zeile wo sie liegen, Exit 0. ~/.claude/skills zaehlt global NICHT: Copilot CLI und Cloud Agent sehen es nicht - sonst haette die CLI gar keine Rollen; Preis: VS Code zeigt sie dann doppelt. Das jaira-cockpit-Mod (Claude-Code-Plugin, keine SKILL.md) wird fuer copilot nicht geschrieben.
- **2026-10-13 17:13 · Alexander Sacharov** — Install bekam eine Agent-Methode statt eines dritten Parameters: die zehn Install(dir, force)-Aufrufe in core/role/role_test.go bleiben so unveraendert. Geprueft wird 'schon installiert' ueber irgendeine Rollen-SKILL.md, nicht ueber alle: schon eine einzige waere doppelt registriert.
- **2026-10-13 17:14 · Alexander Sacharov** — critique: core/role/install.go:46 behaelt das exportierte Install(dir, force) nur noch fuer core/role/role_test.go; das CLI ruft Agent.Install. Exportiertes Install streichen, Tests rufen das unexportierte install(dir, force, true) - ein Einstieg statt zwei. Alles andere bleibt stehen: Agent-Tabelle in target.go ist die vom Ticket verlangte 'eine Zeile pro Agent'; ~/.claude/skills nicht in AlsoGlobal ist in der decision-Note begruendet.
- **2026-10-13 17:15 · Alexander Sacharov** — Nach critique: exportiertes role.Install entfernt; einziger Einstieg ist Agent.Install, die Tests in core/role rufen install(dir, force, true). Die fruehere Note zur Agent-Methode gilt nur noch insoweit.
- **2026-10-13 17:15 · Alexander Sacharov** — critique Runde 2: der Befund aus Runde 1 ist behoben (0655cb10); nichts weiter.
