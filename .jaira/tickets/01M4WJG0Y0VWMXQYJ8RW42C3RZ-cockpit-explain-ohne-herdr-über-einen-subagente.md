---
id: 01M4WJG0Y0VWMXQYJ8RW42C3RZ
title: "Cockpit: Explain ohne Herdr über einen Subagenten statt Fehler"
status: testing
ready: true
creator: Alexander Sacharov
goal: "Explain im Cockpit funktioniert auch ohne Herdr: die Sitzung startet einen Subagenten, der den Review-Befund erklärt und als explain-Note ans Ticket schreibt. Discuss ohne Herdr läuft im aktuellen Chat. Die Fehlermeldung nennt die tatsächliche Ursache."
context: |-
  Was falsch ist: Auf Windows ohne Herdr zeigt [Explain] im Cockpit nur 'Explain 7MX3EJ failed - A chat of its own needs Herdr (HERDR_ENV=1).' Es passiert sonst nichts.
  Wo: core/mod/builtin/jaira-cockpit/hooks/register.tsx, Funktion openChat (Zeile ~551). Ohne HERDR_ENV=1 oder ohne spawn.sh bricht sie mit t.needsHerdr ab. Einen Weg über einen Subagenten gibt es nicht.
  Erwartet hat Alex (14.10.2026): ohne Herdr soll ein Subagent laufen.
  Was es schon gibt: $.prompt.submit schickt Text an die laufende Sitzung (benutzt in decide, register.tsx ~627). Der Explain-Prompt steht in core/role/builtin/jaira-dispatcher/scripts/spawn.sh (lane explain, Zeile ~247). Er schreibt 'jaira note <id> "explain: ..."', und das Cockpit liest diese Note schon unter [Read explanation].
  Discuss: ein Gespräch mit einer Person, das kann ein Subagent nicht führen. Ohne Herdr gehört es in den aktuellen Chat (per $.prompt.submit).
  Zweiter Fehler: derselbe Text t.needsHerdr erscheint bei drei Ursachen (kein HERDR_ENV, kein spawn.sh, keine Board). Man sieht nicht, welche.
definition-of-done: "Explain ohne HERDR_ENV=1 schickt per $.prompt.submit einen Auftrag an die Sitzung, einen Subagenten mit dem Explain-Prompt aus spawn.sh zu starten; kein Fehler-Toast"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-14T06:44:46Z
updated-at: 2026-10-14T06:55:38Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-27701
claimed-at: 2026-10-14T06:45:25Z
outcome-what: "Cockpit: Explain ohne Herdr bittet die Sitzung um einen Subagenten, Discuss läuft im aktuellen Chat; eigene Meldungen für fehlendes spawn.sh und fehlendes Board; auf nativem Windows findet das Cockpit spawn.sh unter USERPROFILE und startet es mit Git Bash"
outcome-why: "Ohne Herdr und auf nativem Windows endete Explain nur in 'needs Herdr' - auch dann, wenn Herdr lief"
outcome-resolves: "Explain/Discuss ohne Herdr, irreführende needsHerdr-Meldung, Windows-in-Herdr (DoD 1-6)"
review-summary: none
review-gaps: "spawnScript: let found/try/catch zu fs.exists(...).catch(() => false) gefaltet; keine Duplikate (USERPROFILE/Git-Bash-Suche gibt es sonst nirgends im Cockpit oder in den Rollen); prNeedsHerdr/nextNeedsHerdr unverändert gelassen (andere Pfade, nicht Teil des Tickets); Prompt-Duplikat zu spawn.sh bewusst (siehe Note)"
---

# Cockpit: Explain ohne Herdr über einen Subagenten statt Fehler

## Definition of Done

- [x] Explain ohne HERDR_ENV=1 schickt per $.prompt.submit einen Auftrag an die Sitzung, einen Subagenten mit dem Explain-Prompt aus spawn.sh zu starten; kein Fehler-Toast
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:568-571 openChat ohne HERDR_ENV -> $.prompt.submit(t.explainBySubagent) (texts.ts, Prompt wie spawn.sh lane explain); Test 'without Herdr, Explain asks this session for a subagent that writes the explain note'
- [x] Discuss ohne HERDR_ENV=1 führt das Gespräch per $.prompt.submit im aktuellen Chat; die Entscheidung landet als Note am Ticket
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:568-571 -> $.prompt.submit(t.discussInThisChat): decision-Note + question= leeren; Test 'without Herdr, Discuss talks the question through in this chat'
- [x] Kein spawn.sh bei HERDR_ENV=1 und keine Board haben je eine eigene Meldung in texts.ts
  proof: core/mod/builtin/jaira-cockpit/hooks/texts.ts noSpawn und noBoard statt needsHerdr; register.tsx:566 (noBoard) und :574 (noSpawn); Test 'with Herdr but no spawn.sh, Explain says spawn.sh is missing'
- [x] Tests in core/mod/builtin/jaira-cockpit/tests decken Explain und Discuss ohne Herdr ab
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 5 neue Tests (ohne Herdr Explain/Discuss, kein spawn.sh, Windows USERPROFILE+Git Bash, CLAUDE_CODE_GIT_BASH_PATH); claude plugin test: 142 pass 0 fail, exit 0
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/42C3RZ.md
- [x] Natives Windows in Herdr: home() fällt auf USERPROFILE zurück, und spawn.sh wird über bash gestartet (Git Bash), damit Explain/Discuss dort einen Herdr-Tab öffnen
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx home() HOME||USERPROFILE mit /; spawnScript() gibt [bash, spawn.sh]; bash() nimmt auf OS=Windows_NT CLAUDE_CODE_GIT_BASH_PATH, sonst <git --exec-path>/../../../bin/bash.exe; Tests 'on native Windows, ...'

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-14 06:45 · Alexander Sacharov** — parent session: jaira-34
- **2026-10-14 06:46 · Alexander Sacharov** — Offene Entscheidungen vor der Arbeit: keine, die nur Alex treffen kann. Selbst entschieden:
1) Prompt-Quelle: der Cockpit-Weg ohne Herdr trägt Explain- und Discuss-Auftrag selbst in texts.ts (englisch, wie alle Cockpit-Texte), inhaltlich gleich zu spawn.sh lane explain/discuss. Warum: ohne Herdr ist spawn.sh oft gar nicht installiert (jaira roles install), und DoD 3 nennt 'kein spawn.sh' nur bei HERDR_ENV=1. Kommentar an beiden Stellen verweist auf die andere.
2) Explain ohne Herdr: $.prompt.submit an die laufende Sitzung: starte einen Subagenten mit diesem Prompt; der Subagent liest selbst jaira show <id> --brief und schreibt die explain:-Note. Das Chat-Modell (.jaira/chat-model) wird dort nicht durchgereicht - Agent kennt nur sonnet/opus/haiku-Namen.
3) Discuss ohne Herdr: $.prompt.submit mit dem Discuss-Auftrag an die Sitzung selbst (kein SendMessage, kein Tab).
4) Meldungen: needsHerdr fällt weg; neu noSpawn (HERDR_ENV=1, spawn.sh fehlt) und noBoard.
- **2026-10-14 06:47 · Alexander Sacharov** — Nachtrag 14.10.: Alex sagt, auf Windows läuft die Sitzung IN Herdr und trotzdem kommt needsHerdr. Verdacht 1: home() in register.tsx:147 liest nur HOME; natives Windows hat oft nur USERPROFILE -> spawnScript sucht /.claude/skills/... und findet nichts -> needsHerdr. Verdacht 2: auch gefunden würde spawn.sh per $.process.run direkt gestartet; Windows führt .sh ohne bash nicht aus. Im Cockpit gibt es keinerlei Windows-Behandlung (kein USERPROFILE, kein bash).
- **2026-10-14 06:53 · Alexander Sacharov** — Umsetzung: spawn.sh startet jetzt überall als 'bash spawn.sh ...' (auch Linux/WSL) - ein Weg, ein Test. Auf Windows NICHT bare 'bash': CreateProcess sucht System32 vor PATH, dort ist bash.exe die WSL-bash. Daher CLAUDE_CODE_GIT_BASH_PATH (setzt Claude Code für Windows), sonst aus git --exec-path abgeleitet. Testmotor nimmt nur absolute POSIX-Pfade (C:/... wird vor dem Mock verworfen) und fällt bei fehlendem HOME offenbar auf das echte HOME zurück - daher HOME:'' und USERPROFILE '\\Users\\a' im Test. Echt auf Windows ungeprüft.
- **2026-10-14 06:54 · Alexander Sacharov** — critique: nichts. Geprüft: Duplikat der Prompts zu spawn.sh steht auf der Note 'Offene Entscheidungen' (bewusst, Querverweis in texts.ts und spawn.sh); 'bash spawn.sh' statt Shebang ist auf Linux/macOS gleichwertig (Shebang ist /usr/bin/env bash, also dieselbe PATH-Suche); bash() außerhalb des try in spawnScript, damit ein Fehler dort nicht als 'kein spawn.sh' erscheint.
