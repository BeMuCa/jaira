# Lanes and their contracts
_Owner: Sascha · EARS · each line one requirement · [tickets] hold the detail_

## The board is its lane files
LANE-01  The board shall consist of exactly the lane files in `.jaira/lanes/`, ordered by the `order` file beside them. [BNZERQ, QA3GN1]
LANE-02  The binary shall ship 13 lanes; critique, optimize and testing carry `default-board: false` and join a board only on `jaira lanes add`. [VSC1GW]
LANE-03  When a board's lane file shares a shipped lane's id, jaira shall use the file as that lane. [BNZERQ, CD9TCB]
LANE-04  If such a file drops requires-human-exit, terminal, requires-outcome or requires-nonmodel-signal of the shipped lane, then jaira shall warn and name the dropped protection. [CD9TCB]
LANE-05  When the first command after an upgrade runs, jaira shall replace every shipped lane file nobody edited by hand and say so in one line. [NH788Q, 1K9KZS]
LANE-06  If a lane's rejects-to names itself or a lane that is not installed, then jaira shall warn about that target. [TQXBY5]
LANE-07  If a lane requires an input that no earlier lane produces, then jaira shall warn and name the lane that produces it, if any. [RFMQ4K]
LANE-08  If a ticket's status names a lane that is not installed, then jaira shall show it in a read-only column and refuse every move of it. [BNZERQ]
LANE-58  `jaira lanes --json` shall list each lane's rejects_to, so a reader finds where a loop lane sends work back. [PGXJG0]

## The contract of one lane
LANE-10  Each lane shall declare its contract in frontmatter: agentic, model-tier, input-requires, output-produces, and the requires-* rails.
LANE-11  When `jaira show <id> --for-lane <lane> --json` runs, jaira shall return the lane's prompt, model tier, bounded input, produces and missing.
LANE-12  The for-lane payload shall carry mode, context, round, back, blocker_rule and stop_after_fix beside the input. [GTQHNH, MG677R, DF0YD6]
LANE-13  Where a lane declares max-rounds, the payload shall carry max_rounds and rounds_left; where it declares verdict words, verdicts. [H6CD33, F1V59P]
LANE-14  Where a lane requires `notes`, the payload shall carry the ticket's progress notes and never report an empty journal as missing. [RFMQ4K]
LANE-15  Where a lane requires `diff`, the payload shall union the recorded commits with git's derivation and add uncommitted work as worktree_diff. [7XXT00]
LANE-16  If neither commits nor worktree yield a diff, then the payload shall list `diff` as missing and say which half was empty. [7XXT00]
LANE-17  When `jaira move --from-lane <lane>` reads JSON on stdin, jaira shall refuse output lacking a field the lane produces (exit 3).
LANE-18  Where a lane declares requires-option, the lane shall be on a ticket's path only when the ticket's Options tick that option. [8WFNQY]

## Path through the board
LANE-20  The next lane of a ticket shall be the first later lane that is not a parking lane, not a question lane, and not an unticked option lane. [BG5QJ6]
LANE-21  While a ticket sits in a terminal, parking or question lane, jaira shall report no next lane. [BG5QJ6]
LANE-22  When a move lands a ticket in an agentic lane, jaira shall print the `jaira show --for-lane` line that works it. [88H1P4]
LANE-23  `jaira next --per-lane --json` shall list every lane with actionable work, its waiting count, whether it is agentic, and its first ticket.
LANE-24  The blocked lane shall be a parking lane: entering it needs no dependency check, counts no round and owes none of the left lane's output.

## Loops and send-backs
LANE-30  Where a lane declares rejects-to, those lanes shall be its documented back edges; any move to an earlier lane is allowed. [TQXBY5]
LANE-31  When a move goes to an earlier lane (parking excluded), jaira shall raise the ticket's round by one. [TW5FE2]
LANE-32  When a move goes back, jaira shall record back-from, back-by (lane or person) and back-count, the run of returns in a row from that lane. [H6CD33]
LANE-33  If a person sends work back, or another lane does, then jaira shall restart back-count at 1. [H6CD33, K6PBYR]
LANE-34  If a lane with max-rounds n would send a ticket back an (n+1)-th time in a row, then jaira shall refuse with code max_rounds. [H6CD33]
LANE-35  The max_rounds refusal shall say the teamlead decides first, or a person in the question lane, and how `back-count=0` restores the rounds. [K6PBYR]
LANE-36  The max_rounds limit shall not apply to a move by a person (board key or interactive terminal), nor into a lane a person answers in. [H6CD33, K6PBYR]
LANE-37  The shipped critique, optimize and testing lanes shall send back to in-progress with max-rounds 3. [H6CD33, ZEFXXM]

## Verdicts and the bug hunt
LANE-40  Where a lane declares verdict-send-back words, a review-verdict starting with one shall refuse every agent move past that lane, parking excluded (exit 3, verdict_sends_back). [F1V59P]
LANE-41  When a review-verdict is anything other than a declared send-back word, free text included, the gate shall pass it. [F1V59P]
LANE-42  The verdict refusal shall name the lane to go back to: the rejects-to lanes, else the last earlier lane producing outcome-what. [F1V59P]
LANE-43  While a board has a lane that produces hunt-verdict, an agent move of a ticket that changed code into a requires-human-exit lane shall be refused unless its newest bug hunt covers HEAD (needs_bughunt, hunt_uncommitted). [G8HNHT, Z66Z8E]
LANE-44  The bug hunt gate shall pass a hunt state of ok or skipped, and shall not guard the question lane. [G8HNHT, DF0YD6]

## Human lanes and context
LANE-50  The question lane (`human`) shall refuse entry without a question, passed with --question or already on the ticket. [BG5QJ6]
LANE-51  An agent shall be allowed to move a ticket on out of the question lane once its question is answered. [BG5QJ6, QNQ0RP]
LANE-52  The signoff lane shall refuse any forward move not made by a key press in the board, unless --force is passed. [Q484YS]
LANE-53  A move from a requires-human-exit lane to an earlier lane shall need neither --force nor a key press. [Q484YS, BQGM6X]
LANE-54  When a person sends a ticket back from signoff in the board, jaira shall ask why and write the answer as a note in the same write. [BQGM6X]
LANE-55  Where no lane declares `context:`, the last agentic lane before the first requires-human-exit lane shall be fresh (review on the shipped board). [MG677R]
LANE-56  Where a lane is fresh, the payload's context shall say so, and that lane is handed to an agent that did not carry the ticket. [MG677R]
LANE-57  `jaira set mode=` shall accept only empty or `conversational`; the mode lives on the ticket so a fresh session reads it. [GTQHNH]
