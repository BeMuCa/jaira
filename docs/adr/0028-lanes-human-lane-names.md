# The two person lanes need names that say which is the question and which the acceptance
Status: Proposed
Date: 2026-10-13

## Context
Lane id `human` is displayed "HITL" (a question to a person); lane id `signoff` is displayed "Human Review" (the acceptance).
Agent text says "the human review lane cannot be left by an agent" meaning signoff; agents read it as `human`.
On 2026-10-13 a session explained `human` as the acceptance lane. The model review lane `review` adds a third "review".

## Decision
Open. Side A (today's files): ids human/signoff, names HITL/Human Review.
Side B (64GNZT): names that say "question" and "acceptance", and agent text that uses the ids only.

## Rejected
- None yet.

## Consequences
Until decided, docs and prompts should name lanes by id (human, signoff), never by display name.

## Sources
64GNZT, core/lane/builtin/30-human.md, core/lane/builtin/45-signoff.md, CLAUDE.md agent block
