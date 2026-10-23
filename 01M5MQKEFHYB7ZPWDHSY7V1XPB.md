---
id: 01M5MQKEFHYB7ZPWDHSY7V1XPB
title: "spawn.sh: Parent-Session selbst ermitteln, ein unerreichbarer Parent bricht laut ab"
status: in-progress
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
updated-at: 2026-10-23T16:32:21Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-54254
claimed-at: 2026-10-23T15:56:29Z
outcome-what: "spawn.sh passes the calling session's messaging socket (uds:) as parent for dispatch/worker, resolves a --parent name through Claude Code's session registry and refuses an unknown one with exit 2 before opening a tab; worker/dispatcher roles forbid rerouting to another peer"
outcome-why: "a teamlead's session name changes on restart or compaction, so workers started with the old name reported to whichever peer they picked"
outcome-resolves: workers report to the session that started them; an unreachable parent is caught at spawn time
review-summary: "core/role/builtin/jaira-dispatcher/scripts/spawn.sh usage: --no-parent says the cockpit passes it only after Accept; it also does on give back (register.tsx:841) — say 'when it starts a dispatcher or worker itself'"
round: 4
review-gaps: "1) Unmentioned behaviour change in the cockpit: Accept -> next ticket (core/mod/builtin/jaira-cockpit/hooks/register.tsx:515) runs spawn.sh <id> <id> dispatch|worker without --parent. Before: that dispatcher asked in its own tab. Now spawn.sh inherits the person's CLAUDE_CODE_MESSAGING_SOCKET (the mod reads the same variable at :608, so it is in its env) and the dispatcher sends its questions via SendMessage into the person's chat. Likely wanted, but the ticket kept discuss unchanged precisely to avoid changing a cockpit button, and neither the notes nor the release note mention this one. Not verified live. 2) 'refuses one that is gone' only holds if the socket file was removed: test -e accepts a socket file a crashed session left behind, and a registry <pid>.json of a dead pid still resolves. Chosen on purpose (-e for Windows), but the release note promises more. 3) DoD 2 says 'ListAgents does not know'; the check reads Claude Code's undocumented session registry instead — the notes say so, acceptable, but it breaks silently if Claude Code changes that format (then every --parent <name> exits 2, the safe direction)."
test-verdict: "pass: go test ./core/role ./internal/cli green (RC=0), full ./... green before optimize; live: spawn.sh --parent <unknown name> and --parent uds:<missing> both exit 2 with the message and open no tab; this session's name resolves through ~/.claude/sessions to its own CLAUDE_CODE_MESSAGING_SOCKET"
review-verdict: "The diff meets all three DoD items and the tests pass; no defect found in spawn.sh. Not unconditional: the person should decide whether cockpit Accept -> next starting a dispatcher that reports to their own chat is wanted (gap 1) — if not, the cockpit needs --parent '' or spawn.sh needs a way to opt out."
review-round: 2
review-check: "1. go test ./core/role -run TestSpawn -count=1 -> ok  2. In a Claude session with Herdr: spawn.sh --parent nobody-9f x1 X1 worker; echo $? -> message '--parent nobody-9f: no running session has that name...' and 2, no new Herdr tab  3. Same without --parent on a real ticket id: the new tab's first line reads /jaira-worker <id> --parent uds:<your session's socket>  4. Cockpit: Accept a ticket so the next one starts; the new dispatcher tab's line carries no --parent and asks in its own tab"
---

# spawn.sh: Parent-Session selbst ermitteln, ein unerreichbarer Parent bricht laut ab

## Definition of Done

- [x] Ohne --parent trägt spawn.sh die aufrufende Session als Parent ein
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:92 (uds:$CLAUDE_CODE_MESSAGING_SOCKET for dispatch/worker); TestSpawnParentIsTheCallerAndAnUnknownOneStops
- [x] Ein --parent, den ListAgents nicht kennt, bricht den Spawn mit klarer Meldung ab
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:111 (exit 2 before tab create); TestSpawnParentIsTheCallerAndAnUnknownOneStops
- [x] Worker/Dispatcher leiten bei unerreichbarem Parent nicht an einen anderen Peer um, sondern fragen im eigenen Tab
  proof: core/role/builtin/jaira-worker/SKILL.md:255, core/role/builtin/jaira-dispatcher/SKILL.md:192
- [x] spawn.sh --no-parent startet Dispatcher/Worker ohne Parent; das Cockpit übergibt es beim Start nach Accept
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:518 (Accept -> next) and :841 (give back); spawn.sh --no-parent; TestSpawnParentIsTheCallerAndAnUnknownOneStops, cockpit board.test.ts 158 pass

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
- **2026-10-23 16:05 · Alexander Sacharov** — critique: teamlead SKILL.md claims the socket never goes stale; it is per process, so it survives compaction and a rename but not a restart. Reword. Left alone on purpose: name resolution through the session registry stays (DoD 2 allows names ListAgents knows; refusing all names would break existing --parent <name> callers); cockpit's own uds: for discuss in register.tsx:608 is the same pattern, not a duplicate to remove.
- **2026-10-23 16:08 · Alexander Sacharov** — verified: SendMessage accepts a uds:<socket path> address — a progress line sent that way reached the parent session's inbox (success, msg delivered). The tool's own description mentions only names, so this was the open risk of the uds: approach.
- **2026-10-23 16:11 · Alexander Sacharov** — review: DoD met, tests green. Open for the person: cockpit Accept -> next (register.tsx:515) now silently starts the dispatcher with the person's session as parent, so its questions arrive in the person's chat instead of its own tab — not mentioned in notes or release note. Decide whether that is wanted.
- **2026-10-23 16:12 · Alexander Sacharov** — open for the person (review gap 1): cockpit Accept -> next (register.tsx:515) runs spawn.sh dispatch|worker without --parent, so that dispatcher now reports into the person's own chat instead of asking in its tab. 1) keep it, add a release-note line (recommended) 2) cockpit opts out explicitly, spawn.sh gets a no-parent switch. Sent to the parent session.
- **2026-10-23 16:29 · Alexander Sacharov** — decision (person, via parent session): option 2 — a dispatcher/worker the cockpit starts stays in its own tab. The cockpit (register.tsx:515) passes an explicit no-parent, spawn.sh gets --no-parent for it.
- **2026-10-23 16:31 · Alexander Sacharov** — the give-back path (register.tsx:841, a new worker when the person sends a ticket back) also starts a worker from the person's session; the decision 'a worker the cockpit starts stays in its own tab' covers it, so it passes --no-parent too. Discuss (register.tsx:612) keeps its explicit uds: parent — unchanged. Cockpit tests: 'claude plugin test .' in core/mod/builtin/jaira-cockpit.
- **2026-10-23 16:32 · Alexander Sacharov** — critique round 2: --no-parent help names only Accept, the give-back path passes it too; reword. Left alone: an older spawn.sh would refuse --no-parent as an unknown flag, but 'jaira roles install' writes the cockpit and spawn.sh together (internal/cli/roles.go:100), so they cannot drift apart.
