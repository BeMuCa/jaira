---
id: 01M5MQKEFHYB7ZPWDHSY7V1XPB
title: "spawn.sh: Parent-Session selbst ermitteln, ein unerreichbarer Parent bricht laut ab"
status: done
ready: true
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
commits:
  - 5b028ed15562bcac3721dca870e5f4098b749b97
  - fc282ed8dc14df798c2d82d6eefa8aaf555e2259
  - 97a86381aec403cfc78372860ffa4be46c7762e7
  - 3c7d46000179021c3c223a20414a063d4321aa47
  - 3785c18ef0e6b9eaa4cece8a8e7108a430920568
created-at: 2026-10-23T15:55:47Z
updated-at: 2026-10-23T16:38:16Z
outcome-what: "spawn.sh passes the calling session's messaging socket (uds:) as parent for dispatch/worker, resolves a --parent name through Claude Code's session registry and refuses an unknown one with exit 2 before opening a tab; worker/dispatcher roles forbid rerouting to another peer"
outcome-why: "a teamlead's session name changes on restart or compaction, so workers started with the old name reported to whichever peer they picked"
outcome-resolves: workers report to the session that started them; an unreachable parent is caught at spawn time
review-summary: "spawn.sh started a worker/dispatcher with whatever --parent name it was given; a stale name (teamlead restarted/compacted) made workers improvise and report to the oldest peer. Now: without --parent, spawn.sh passes the calling session's own messaging socket (uds:$CLAUDE_CODE_MESSAGING_SOCKET) as parent for the dispatch and worker lanes. An explicit --parent is checked before any tab is made: uds:<path> must exist, a name is looked up in Claude Code's session registry (~/.claude/sessions/<pid>.json) and replaced by that session's uds: address; unknown -> exit 2 with a message naming it. New --no-parent switch; the cockpit passes it on both of its own starts (next ticket after Accept, register.tsx:518; give-back worker, :841) so those agents keep asking in their own tab. worker/dispatcher SKILL.md forbid rerouting to another peer when the parent is unreachable (note on ticket, ask in own tab); teamlead SKILL.md now says: pass no --parent. Four release-note lines in core/release/notes/7V1XPB.md."
round: 4
review-gaps: "No DoD gap. Minor, not blocking: (1) two running sessions with the same name -> the name lookup silently takes the first registry file in glob order; ListAgents would show both. (2) the 'session is running' test is only [ -e socket ]; on this machine all 88 dead registry entries had their socket removed, so normal exit is covered, but a SIGKILLed session leaving its socket file behind would pass (unverified). (3) DoD 3 (no rerouting to another peer) is SKILL.md prose only, no test can hold it - inherent to a prompt rule."
test-verdict: "pass: go test ./... green (RC=0), cockpit 'claude plugin test' 158 pass; DoD 1-4 checked in the tree; --no-parent covered by TestSpawnParentIsTheCallerAndAnUnknownOneStops (worker typed without --parent while the caller's socket is in the env), not run live because it would open a real tab"
review-verdict: "The diff meets all four DoD items; go test ./core/role and cockpit 'claude plugin test .' (158 pass) green; spawn.sh refuses a stale name and a dead uds: with exit 2 before a tab, checked by hand. Approve; the two minor gaps can stay."
review-round: 4
review-check: "1. In a Claude session in this worktree run: bash core/role/builtin/jaira-dispatcher/scripts/spawn.sh --no-worktree --parent lead-9f x7 X7 worker ; echo $?  -> message '--parent lead-9f: no running session has that name ...' and 2, no Herdr tab opens.  2. Same with --parent uds:/nope.sock -> 'no session listens there' and 2.  3. go test ./core/role -run TestSpawn -> ok (covers: no --parent passes uds:<caller socket> for worker and dispatch; review lane gets none; --no-parent gets none).  4. In the cockpit (claude --plugin-dir core/mod/builtin/jaira-cockpit, inside Herdr) accept a ticket so the next one starts -> the new tab's first line is '/jaira-dispatcher <id>' with no --parent.  5. As teamlead, run spawn.sh <slug> <id> dispatch without --parent -> the tab's line ends in '--parent uds:/run/user/.../<pid>.sock'."
---

# spawn.sh: Parent-Session selbst ermitteln, ein unerreichbarer Parent bricht laut ab

## Definition of Done

- [x] Ohne --parent trägt spawn.sh die aufrufende Session als Parent ein
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:97 (uds:$CLAUDE_CODE_MESSAGING_SOCKET for dispatch/worker); TestSpawnParentIsTheCallerAndAnUnknownOneStops
- [x] Ein --parent, den ListAgents nicht kennt, bricht den Spawn mit klarer Meldung ab
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:116 (exit 2 before tab create); TestSpawnParentIsTheCallerAndAnUnknownOneStops
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
- **2026-10-23 15:57** — decision: spawn.sh fills the parent in itself as uds:$CLAUDE_CODE_MESSAGING_SOCKET, not as a session name. Why: the name is what goes stale (it changes on restart/compaction); the socket belongs to the process and survives a rename. Claude Code sets this variable for every Bash call of a session.
decision: an explicit --parent is checked before a tab is made. uds:<path> must be a live socket (test -S). A plain name is resolved through Claude Code's session registry (~/.claude/sessions/<pid>.json: name + messagingSocketPath) to its uds: address, so the worker gets the stable address, not the name; a name not found there, or no registry, aborts with exit 2. The registry is internal to Claude Code and undocumented — the only way a script can see what ListAgents sees.
decision: the automatic parent applies only to the lanes dispatch and worker (what the ticket is about). discuss keeps its parent only when one is passed: the cockpit starts it, and sending its decision somewhere new would change that button's behaviour.
decision: worker/dispatcher SKILL.md get an explicit rule: a parent SendMessage cannot reach is never replaced by another peer; ask in the own tab (AskUserQuestion) and note that the parent was unreachable.
- **2026-10-23 16:04** — found: Claude Code keeps a registry of running sessions in ~/.claude/sessions/<pid>.json (name, messagingSocketPath, nameSince); the name in it changes, the socket path is per process. The 'oldest peer' fallback was never written in any SKILL.md — the workers improvised it when SendMessage failed, so the fix is an explicit prohibition, not removing a rule. -e instead of -S on the socket: unverified how a native Windows claude's messaging address looks to Git Bash; -S could stop every cockpit Discuss there.
- **2026-10-23 16:05** — critique: teamlead SKILL.md claims the socket never goes stale; it is per process, so it survives compaction and a rename but not a restart. Reword. Left alone on purpose: name resolution through the session registry stays (DoD 2 allows names ListAgents knows; refusing all names would break existing --parent <name> callers); cockpit's own uds: for discuss in register.tsx:608 is the same pattern, not a duplicate to remove.
- **2026-10-23 16:08** — verified: SendMessage accepts a uds:<socket path> address — a progress line sent that way reached the parent session's inbox (success, msg delivered). The tool's own description mentions only names, so this was the open risk of the uds: approach.
- **2026-10-23 16:11** — review: DoD met, tests green. Open for the person: cockpit Accept -> next (register.tsx:515) now silently starts the dispatcher with the person's session as parent, so its questions arrive in the person's chat instead of its own tab — not mentioned in notes or release note. Decide whether that is wanted.
- **2026-10-23 16:12** — open for the person (review gap 1): cockpit Accept -> next (register.tsx:515) runs spawn.sh dispatch|worker without --parent, so that dispatcher now reports into the person's own chat instead of asking in its tab. 1) keep it, add a release-note line (recommended) 2) cockpit opts out explicitly, spawn.sh gets a no-parent switch. Sent to the parent session.
- **2026-10-23 16:29** — decision (person, via parent session): option 2 — a dispatcher/worker the cockpit starts stays in its own tab. The cockpit (register.tsx:515) passes an explicit no-parent, spawn.sh gets --no-parent for it.
- **2026-10-23 16:31** — the give-back path (register.tsx:841, a new worker when the person sends a ticket back) also starts a worker from the person's session; the decision 'a worker the cockpit starts stays in its own tab' covers it, so it passes --no-parent too. Discuss (register.tsx:612) keeps its explicit uds: parent — unchanged. Cockpit tests: 'claude plugin test .' in core/mod/builtin/jaira-cockpit.
- **2026-10-23 16:32** — critique round 2: --no-parent help names only Accept, the give-back path passes it too; reword. Left alone: an older spawn.sh would refuse --no-parent as an unknown flag, but 'jaira roles install' writes the cockpit and spawn.sh together (internal/cli/roles.go:100), so they cannot drift apart.
- **2026-10-23 16:36** — review (fresh context): DoD met, tests green, stale-name and dead-uds refusal checked by hand. Minor, left open: same-name sessions take the first registry hit; a SIGKILLed session's leftover socket file would pass the -e check (unverified).
