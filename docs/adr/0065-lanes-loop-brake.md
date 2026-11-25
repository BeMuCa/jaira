# A loop lane stops after max-rounds returns in a row, and the teamlead decides first
Status: Accepted
Date: 2026-11-22

## Context
critique, optimize and testing send work back to in-progress until they find nothing; such a loop can run forever.
The limit used to live in prompt prose, where an agent could ignore it.

## Decision
A lane declares max-rounds: n; the move records back-from, back-by and back-count, and the gate refuses the (n+1)-th return in a row (max_rounds).
The refusal hands the decision to the parent teamlead session, else to a person in the question lane; `jaira set <id> back-count=0` restores the rounds.

## Rejected
- Counting every send-back of the ticket: a person's return or another lane's is not a loop.
- Going straight to the person after round three: the teamlead can usually decide (K6PBYR).
- Letting --force reset the count: --force passes one move, it does not give rounds back.

## Consequences
A person's own send-back and a move into a person's lane are never refused by the limit.
The shipped loop lanes declare max-rounds 3; for-lane payloads carry max_rounds and rounds_left.

## Sources
H6CD33, K6PBYR (bughunt note 2026-11-22), X9GJ51
