---
id: 01M5WTKDTVNQ9HMQJ9NZW9AVBE
title: "Kein Lane-Prompt schreibt mehr --commits \"$(git rev-parse HEAD)\" vor, das immer nachfragt"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Ein Worker faehrt die in-progress-Lane mit der Erlaubnisliste aus .claude/settings.json ohne Rueckfrage durch; kein Prompt, kein Skill und keine Doku schreibt ihm eine Shell-Expansion vor, die Claude Code nie automatisch erlaubt."
context: |-
  Claude Code fragt bei "$(git rev-parse HEAD)" immer nach, wenn die Ersetzung ein ganzes Argument ist. Keine Allow-Regel hilft dagegen.
  Gemessen am 26.10.2026 beim Experiment fuer ZE9WZC (claude -p, --permission-prompts none, Regeln aus .claude/settings.json): denied wurden jaira move ... --commits "$(git rev-parse HEAD)" und echo "$(git rev-parse HEAD)". Durch gingen jaira note X "text $(git rev-parse HEAD)" und alle Befehle der Liste. Transkripte lagen im Scratchpad, die Zusammenfassung steht als Notiz an ZE9WZC.
  Genau diese Form schreiben vor: core/lane/builtin/20-in-progress.md:68 (und die Kopie .jaira/lanes/in-progress.md:68), .claude/skills/jaira/SKILL.md:307 und :322, docs/AGENTS.md:36, README.md:846.
  --commits ist gar nicht noetig: jaira leitet die Commit-Liste selbst aus der git-Historie ab (CLAUDE.md-Block, Punkt "jaira move").
  Auch $? und ${PIPESTATUS[0]} fragen immer nach; das betrifft nur Befehle, die ein Worker selbst zusammenbaut.
definition-of-done: "Keine der genannten Stellen schreibt mehr --commits \"$(git rev-parse HEAD)\" vor; sie lassen --commits weg oder nennen die sha woertlich"
tags:
  - docs
blocked-by: []
related: []
follows: 01M5MNWG4WFAPP5YE4T0ZE9WZC
commits: []
created-at: 2026-10-26T19:22:08Z
updated-at: 2026-10-26T19:22:08Z
---

# Kein Lane-Prompt schreibt mehr --commits "$(git rev-parse HEAD)" vor, das immer nachfragt

## Definition of Done

- [ ] Keine der genannten Stellen schreibt mehr --commits "$(git rev-parse HEAD)" vor; sie lassen --commits weg oder nennen die sha woertlich
- [ ] Ein frischer Worker auf einem Scratch-Board faehrt die in-progress-Lane mit dem neuen Prompt ohne permission_denials durch (claude -p --permission-prompts none); das Ergebnis steht als Notiz am Ticket

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

