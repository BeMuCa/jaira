# Where an agent's model is chosen is undecided
Status: Proposed
Date: 2026-11-21

## Context
Side A: Y2CD9X (decision 2026-11-21) puts `model: haiku` in the SKILL.md frontmatter of tester and logbook-summary; no board setting, no spawn.sh flag.
Side B: N2QT94 has spawn.sh pass --model from the lane's model-tier, per machine and board; CXH98W reads orchestrator and worker models from settings.
Today spawn.sh passes --model only for explain/discuss (the chat model).

## Decision
Open: role frontmatter, lane model-tier via spawn.sh, or a settings entry — and which wins when they disagree.

## Rejected
- `model: haiku` on jaira-role-pr (Y2CD9X decision 2026-11-21 22:01: removed; pr runs on the session's model).

## Consequences
Three mechanisms for one choice; a subagent inherits a skill's model, so the order of precedence matters.

## Sources
Y2CD9X, N2QT94, CXH98W
