---
id: 01M4V3RJKTPTBT43QAWCE75R4Q
title: "jaira roles install --agent copilot legt die Rollen dort ab, wo Copilot sie liest"
status: in-progress
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
updated-at: 2026-10-13T18:54:02Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-342252
claimed-at: 2026-10-13T18:53:03Z
outcome-what: "'jaira roles install' hat --agent (Default claude, neu copilot): eine Agent-Tabelle in core/role/target.go nennt je Agent das Projekt- und das globale Skills-Verzeichnis und die weiteren, die er liest; copilot schreibt nach .github/skills bzw. ~/.copilot/skills ohne das Cockpit-Mod, schreibt nichts wenn die Rollen schon in .claude/skills/.agents/skills liegen, unbekannter Agent und --agent mit --into enden mit Exit 2"
outcome-why: "Wer mit GitHub Copilot arbeitet, bekam die Rollen nur ueber --into und musste das Verzeichnis selbst kennen"
outcome-resolves: "Alle fuenf DoD-Punkte: Recherche-Note mit Quellen, Zielverzeichnisse, keine Doppelkopie, Usage-Fehler, Doku und Release-Note"
review-summary: "'jaira roles install' hat ein Flag --agent (Default claude, neu copilot). Eine Tabelle in core/role/target.go nennt je Agent das Projekt- und das globale Zielverzeichnis, die weiteren Verzeichnisse, die er liest, und ob er das Cockpit-Mod laedt. copilot schreibt nach .github/skills bzw. ~/.copilot/skills, ohne jaira-cockpit. Liegt schon eine Rolle in einem Verzeichnis, das Copilot sonst liest (.claude/skills, .agents/skills; global ~/.agents/skills), schreibt der Befehl nichts, Exit 0, und nennt in einer Zeile den Befehl, der die Rollen dort aktualisiert (role.UpdateCommand: --project/--global des Agenten, der dorthin schreibt, sonst --into <dir>). Unbekannter Agent und --agent mit --into enden mit Exit 2. Die Flag-Hilfe von --project/--global nennt die Verzeichnisse aller Agenten aus derselben Tabelle. Das exportierte role.Install ist durch Agent.Install ersetzt; Tests in core/role rufen install(dir, force, true). README, docs/AGENTS.md (Tabelle) und core/release/notes/E75R4Q.md (zwei Zeilen) sind nachgezogen."
round: 4
review-gaps: "none - drei Randbemerkungen, keine verhindert die Abnahme: (1) Der Long-Text von 'roles install --help' sagt im Absatz vor --agent weiterhin '--project writes into .claude/skills' / '--global writes into ~/.claude/skills'; der folgende --agent-Absatz stellt es richtig, die Flag-Zeilen stimmen. (2) Im Skip-Fall wird --force stillschweigend ignoriert (gewollt: auch --force soll keine Doppelkopie anlegen, aber die Zeile sagt es nicht). (3) Bekannte, bewusst offene Luecke aus der testing-Note: erst --agent copilot --project, dann claude --project ergibt die Rollen doppelt fuer Copilot."
test-verdict: "pass: go test -race -count=1 ./... RC=0; DoD 1-6 im Baum geprueft; Binary: Skip-Zeile endet mit 'update them there with jaira roles install --project' (.claude/skills) bzw. --into (~/.agents/skills), JSON hat update; --help nennt claude- und copilot-Pfade fuer --project/--global"
review-verdict: "Der Diff erfuellt alle sechs DoD-Punkte, auch den nachgereichten DoD 6 (Update-Befehl in der Skip-Zeile, Flag-Hilfe mit allen Agenten-Pfaden). go build, go vet und go test ./core/role ./internal/cli gruen; manueller Lauf mit gebautem Binary (Projekt- und Global-Skip, JSON, unbekannter Agent, --into) verhaelt sich wie beschrieben. Commit-Liste vollstaendig: 5 Commits auf master, alle mit E75R4Q. Kein Defekt gefunden."
review-round: 3
review-check: "1. go build -o /tmp/jaira ./cmd/jaira  2. mkdir /tmp/r && cd /tmp/r && git init && /tmp/jaira init  3. /tmp/jaira roles install --agent copilot --project -> letzte Zeile 'roles in /tmp/r/.github/skills'; ls .github/skills zeigt jaira-teamlead, aber kein jaira-cockpit  4. rm -rf .github && /tmp/jaira roles install --project && /tmp/jaira roles install --agent copilot --project -> genau eine Zeile 'copilot already reads the roles in /tmp/r/.claude/skills; nothing written to ... — update them there with 'jaira roles install --project''; echo $? zeigt 0; ls .github meldet 'No such file'  5. /tmp/jaira roles install --agent vim --project -> 'unknown agent \"vim\"; known agents: claude, copilot', echo $? zeigt 2  6. /tmp/jaira roles install --agent copilot --into /tmp/x -> Fehler, echo $? zeigt 2  7. HOME=/tmp/h /tmp/jaira roles install --agent copilot --global -> 'roles in /tmp/h/.copilot/skills'  8. /tmp/jaira roles install --help -> --project-Zeile nennt 'claude: .claude/skills, copilot: .github/skills', --global-Zeile '~/.claude/skills' und '~/.copilot/skills'"
---

# jaira roles install --agent copilot legt die Rollen dort ab, wo Copilot sie liest

## Definition of Done

- [x] Eine Note mit Quellen sagt, welche Skills-Verzeichnisse Copilot (VS Code, CLI, Coding Agent) projektweit und global liest und ob es SKILL.md mit name/description versteht
  proof: jaira note (Recherche, 13.10.2026) mit den vier Quell-URLs
- [x] 'jaira roles install --agent copilot --project' und '--global' schreiben in das Verzeichnis aus dieser Note; ohne --agent bleibt alles wie heute (claude)
  proof: core/role/target.go:50; TestRolesInstallAgentCopilotWritesWhereCopilotReads (internal/cli/roles_test.go:294); bestehende TestRolesInstall*Project/Global unveraendert gruen
- [x] Liest Copilot ein Verzeichnis, in dem die Rollen schon liegen, legt --agent copilot keine zweite Kopie an und sagt das in einer Zeile
  proof: internal/cli/roles.go:171; TestRolesInstallAgentCopilotSkipsWhereRolesAlreadyAre (internal/cli/roles_test.go:323)
- [x] Ein unbekannter --agent-Wert bricht mit Exit-Code fuer Usage ab und nennt die bekannten Agenten; --agent zusammen mit --into wird abgelehnt
  proof: internal/cli/roles.go:121 und :130; TestRolesInstallAgentRefusesUnknownAndInto (internal/cli/roles_test.go:341)
- [x] jaira roles install --help, README und docs/AGENTS.md nennen --agent; Release-Note als Datei core/release/notes/<Handle>.md
  proof: internal/cli/roles.go Long-Text + --agent-Flag; README.md:661; docs/AGENTS.md Tabelle Agent/--project/--global; core/release/notes/E75R4Q.md
- [x] Die 'already reads'-Zeile nennt 'jaira roles install --project' zum Aktualisieren; die Hilfe von --project/--global stimmt auch fuer --agent copilot
  proof: core/role/target.go:130 Dirs + :146 UpdateCommand; internal/cli/roles.go:171 und :183-184; TestRolesInstallAgentCopilotSkipsWhereRolesAlreadyAre, TestRolesInstallHelpNamesEveryAgentsDirectories (internal/cli/roles_test.go:359)
- [x] Der Long-Text von 'roles install --help' nennt die Verzeichnisse von --project/--global nicht mehr fest als Claude-Pfade, sondern sagt, dass --agent sie bestimmt
  proof: internal/cli/roles.go:99-107 (Long-Text); go test ./internal/cli gruen

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
- **2026-10-13 17:19 · Alexander Sacharov** — testing: Bekannte Luecke ausserhalb der DoD - wer zuerst --agent copilot --project (.github/skills) und danach claude --project (.claude/skills) installiert, hat die Rollen fuer Copilot doppelt; ein weiterer --agent copilot-Lauf aktualisiert .github/skills dann nicht mehr, sondern meldet .claude/skills. claude kennt Copilots Verzeichnisse nicht, das bleibt bewusst so (eine Zeile pro Agent, kein Querwissen).
- **2026-10-13 17:21 · Alexander Sacharov** — review: DoD erfuellt, Tests gruen, manuell geprueft. Offen als kleine Nachbesserung: Skip-Meldung nennt nicht, wie man die Rollen in .claude/skills aktualisiert; Flag-Hilfe --project/--global nennt immer Claude-Pfade. Details in review-gaps.
- **2026-10-13 17:39 · Alexander Sacharov** — Alex am 13.10. (aus signoff zurueck): review-gaps 2+3 nachbessern. (a) Die 'already reads'-Zeile nennt 'jaira roles install --project' als Weg zum Aktualisieren. (b) Die Flag-Hilfe von --project/--global nennt nicht nur die Claude-Pfade, sondern sagt, dass --agent den Pfad bestimmt (oder nennt die Copilot-Pfade mit). An denselben Worker (jaira-e75r4q-9b), Cache warm.
- **2026-10-13 17:40 · Alexander Sacharov** — explain: Der Review hat keinen Fehler gefunden, alle fuenf DoD-Punkte sind erfuellt und die Tests sind gruen. Offen sind zwei kleine Textluecken: Die Meldung "copilot already reads the roles in .claude/skills" sagt nicht, wie man die Rollen dort aktualisiert, und die Flag-Hilfe von --project/--global nennt immer nur die Claude-Pfade, auch bei --agent copilot.
Warum: Nach einem Upgrade laeuft "jaira roles install --agent copilot" sonst ins Leere, ohne dass der Nutzer erfaehrt, dass er "jaira roles install --project" ohne --agent braucht. Die Flag-Zeile schickt Copilot-Nutzer zum falschen Verzeichnis.
1) Beide Texte nachbessern: Skip-Meldung nennt den Aktualisierungsbefehl, Flag-Hilfe sagt, dass --agent den Pfad bestimmt [give back]
2) So annehmen, die Texte sind klein und die Langhilfe und docs/AGENTS.md stimmen schon [accept]
3) Texte nachbessern und zusaetzlich die doppelte Installation (.github/skills und .claude/skills je nach Reihenfolge) abfangen [give back]
Empfehlung: 1 - zwei Textzeilen, wenige Minuten Arbeit, und genau das hat Alex am 13.10. beim Zuruecksenden aus signoff schon verlangt; die Doppel-Installation ist bewusst ausserhalb der DoD.
- **2026-10-13 17:41 · Alexander Sacharov** — DoD 6: Hinweis-Befehl kommt aus role.UpdateCommand - das --project/--global des Agenten, der genau dorthin schreibt (.claude/skills -> 'jaira roles install --project'), sonst --into <dir> (.agents/skills schreibt kein Agent). Flag-Hilfe aus role.Dirs, also aus derselben Tabelle - ein neuer Agent erscheint dort von selbst. Proof von DoD 3 auf roles.go:171 nachgezogen.
- **2026-10-13 17:41 · Alexander Sacharov** — critique Runde 3 (DoD 6): UpdateCommand und Dirs lesen dieselbe Agent-Tabelle wie die Ziele - kein zweites Wissen ueber Pfade; nichts weiter.
- **2026-10-13 17:47 · Alexander Sacharov** — review Runde 3: alle sechs DoD-Punkte erfuellt, Tests gruen, manuell geprueft inkl. DoD 6. Kein Defekt; drei Randbemerkungen in review-gaps (Long-Text-Absatz vor --agent, --force im Skip-Fall still, bekannte Doppel-Installation).
- **2026-10-13 18:52 · Alexander Sacharov** — decision: 1) Beide Texte nachbessern: Skip-Meldung nennt den Aktualisierungsbefehl, Flag-Hilfe sagt, dass --agent den Pfad bestimmt
- **2026-10-13 18:52 · Alexander Sacharov** — Given back from signoff: 1) Beide Texte nachbessern: Skip-Meldung nennt den Aktualisierungsbefehl, Flag-Hilfe sagt, dass --agent den Pfad bestimmt
- **2026-10-13 18:52 · Alexander Sacharov** — cockpit: given back to a new worker - nothing recorded about the old worker session (cache and context unknown)
- **2026-10-13 18:53 · Alexander Sacharov** — Give-back 18:52 waehlte Option 1 aus der explain-Note von 17:40 - die ist aelter als DoD 6 (17:41, Commit 28594e54): Skip-Zeile mit Update-Befehl und Flag-Zeilen mit allen Agent-Pfaden sind schon drin. Was von Option 1 noch nicht stimmt, ist review-gap (1): der Long-Text von 'roles install --help' sagt im Absatz vor --agent weiterhin fest '--project writes into .claude/skills' / '--global writes into ~/.claude/skills'. Diese Runde macht genau diesen Absatz agent-abhaengig, sonst nichts.
