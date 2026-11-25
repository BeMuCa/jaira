# In a person's lane mid-board, Accept goes one lane on; Give back and Discuss go to the agent lane before
Status: Accepted
Date: 2026-11-25

## Context
The cockpit could accept only from a lane whose next lane is terminal (signoff -> done).
On a video board Sascha could not release a script waiting in script-approval; only a hand-typed `jaira move --to build` worked (M5JSE9).
Give back and Discuss moved work to a fixed in-progress, which a board without that lane refuses (PGXJG0, 574VCW).

## Decision
Accept moves the ticket into its next_lane, terminal or not; the Execute line names the lane ("Execute: accept → build").
An open question in a mid-board lane stops Accept until it is answered.
Only an Accept into the terminal lane opens a pull request, starts the land lane or holds the ticket in signoff.
Give back and Discuss move work to the back lane: walking back from the person's lane, the first agentic lane that is not fresh
and not a person's; a loop lane points on to its rejects-to target. Default board: in-progress; video board: draft-approval -> build.
Where `jaira lanes --json` carries no rejects_to (older jaira), the target stays in-progress.

## Rejected
- Accept only into the terminal lane (2V50Y5 to M5JSE9): a person cannot release work mid-board from the cockpit.
- A fixed in-progress for Give back and Discuss: fails on every board without that lane.
- The fresh review lane as target: it judges finished work and does none itself (PGXJG0 skips fresh lanes).

## Consequences
spawn.sh takes `--back <lane>` for discuss, so the Herdr chat and this session move to the same lane.
`jaira lanes --json` lists rejects_to for every lane.
A decision page (decision-page skill) that accepts tickets should use this same path.

## Sources
M5JSE9, PGXJG0 (decision 2026-11-25 10:33), 574VCW (decision 2026-11-25 11:39), core/mod/builtin/jaira-cockpit/hooks/board.ts backLaneOf, register.tsx accepting
