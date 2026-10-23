---
id: 01M5MQKEFHYB7ZPWDHSY7V1XPB
title: "spawn.sh: Parent-Session selbst ermitteln, ein unerreichbarer Parent bricht laut ab"
status: critique
ready: true
creator: Alexander Sacharov
goal: "Worker und Dispatcher melden sich immer bei der Session, die sie gestartet hat — auch nach Neustart oder Komprimierung des Teamleads."
context: |-
  Gefunden am 23.10.2026 in einer langen Teamlead-Session in einem anderen Projekt.
  Was falsch ist: Der Teamlead startete Worker mit --parent <projekt>-9f. So hieß seine Session früher; nach Neustart/Komprimierung heißt sie anders (z. B. <projekt>-92) und der Name aus der Zusammenfassung ist veraltet.
  Jeder Worker meldete 'parent <projekt>-9f nicht erreichbar' und schickte seinen Bericht an den ältesten Peer. Das war zufällig der richtige Teamlead — ein zweiter Teamlead bekäme die Berichte seiner Worker nie.
  Wo: core/role/builtin/jaira-dispatcher/scripts/spawn.sh (--parent), Rollen jaira-worker/jaira-dispatcher (Fallback 'ältester Peer'), jaira-teamlead SKILL.md (sagt: --parent <your-name> wie ListAgents es zeigt).
  Idee: spawn.sh ermittelt den Namen der aufrufenden Session selbst (Umgebung/Socket der Claude-Session), --parent nur noch als Überschreibung; ein angegebener, aber nicht existierender Parent bricht beim Spawn ab statt später still umzuleiten. Der Worker leitet nie an einen fremden Peer um.
definition-of-done: Ohne --parent trägt spawn.sh die aufrufende Session als Parent ein
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:55:47Z
updated-at: 2026-10-23T16:04:33Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-54254
claimed-at: 2026-10-23T15:56:29Z
outcome-what: "spawn.sh passes the calling session's messaging socket (uds:) as parent for dispatch/worker, resolves a --parent name through Claude Code's session registry and refuses an unknown one with exit 2 before opening a tab; worker/dispatcher roles forbid rerouting to another peer"
outcome-why: "a teamlead's session name changes on restart or compaction, so workers started with the old name reported to whichever peer they picked"
outcome-resolves: workers report to the session that started them; an unreachable parent is caught at spawn time
---

# spawn.sh: Parent-Session selbst ermitteln, ein unerreichbarer Parent bricht laut ab

## Definition of Done

- [x] Ohne --parent trägt spawn.sh die aufrufende Session als Parent ein
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:89 (uds:$CLAUDE_CODE_MESSAGING_SOCKET for dispatch/worker); TestSpawnParentIsTheCallerAndAnUnknownOneStops
- [x] Ein --parent, den ListAgents nicht kennt, bricht den Spawn mit klarer Meldung ab
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:94 (exit 2 before tab create); TestSpawnParentIsTheCallerAndAnUnknownOneStops
- [x] Worker/Dispatcher leiten bei unerreichbarem Parent nicht an einen anderen Peer um, sondern fragen im eigenen Tab
  proof: core/role/builtin/jaira-worker/SKILL.md:254, core/role/builtin/jaira-dispatcher/SKILL.md:192

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-23 15:57 · Alexander Sacharov** — decision: spawn.sh fills the parent in itself as uds:$CLAUDE_CODE_MESSAGING_SOCKET, not as a session name. Why: the name is what goes stale (it changes on restart/compaction); the socket belongs to the process and survives a rename. Claude Code sets this variable for every Bash call of a session.
decision: an explicit --parent is checked before a tab is made. uds:<path> must be a live socket (test -S). A plain name is resolved through Claude Code's session registry (~/.claude/sessions/<pid>.json: name + messagingSocketPath) to its uds: address, so the worker gets the stable address, not the name; a name not found there, or no registry, aborts with exit 2. The registry is internal to Claude Code and undocumented — the only way a script can see what ListAgents sees.
decision: the automatic parent applies only to the lanes dispatch and worker (what the ticket is about). discuss keeps its parent only when one is passed: the cockpit starts it, and sending its decision somewhere new would change that button's behaviour.
decision: worker/dispatcher SKILL.md get an explicit rule: a parent SendMessage cannot reach is never replaced by another peer; ask in the own tab (AskUserQuestion) and note that the parent was unreachable.
- **2026-10-23 16:04 · Alexander Sacharov** — found: Claude Code keeps a registry of running sessions in ~/.claude/sessions/<pid>.json (name, messagingSocketPath, nameSince); the name in it changes, the socket path is per process. The 'oldest peer' fallback was never written in any SKILL.md — the workers improvised it when SendMessage failed, so the fix is an explicit prohibition, not removing a rule. -e instead of -S on the socket: unverified how a native Windows claude's messaging address looks to Git Bash; -S could stop every cockpit Discuss there.
