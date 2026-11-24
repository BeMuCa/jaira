# Moving a ticket: gates, definition of done, claim and filing
_Owner: Sascha · EARS · each line one requirement · [tickets] hold the detail_

## One move, one gate
GATE-01  Every lane change from the CLI, the board and the sign-off accept shall run through one move function: stage, gate, write, settle. [JJ32B4, 6BFHJ4]
GATE-02  `jaira move` shall check the gates at the moment of writing, whatever `jaira next` said before.
GATE-03  When a ticket is already in the target lane and the move brings no field, `jaira move` shall succeed and write nothing.
GATE-04  When `jaira move --dry-run` runs, jaira shall run the same gates on the staged ticket, write nothing, and exit as the real move would.
GATE-05  If a gate refuses, then `jaira move` shall exit 3, or 4 when any refusal is a dependency, and name each refusal with a stable code.
GATE-06  The fields a move brings (outcome, question, reason, commits, executed-by) shall be written before the gate and kept if the gate refuses.
GATE-07  If the ticket file on disk is older than the ticket's ref, then `jaira move` shall refuse before the gate runs. [F3VNVC]
GATE-08  When `--force` overrides refusals, jaira shall print them first and the successful move as the last line. [HVHZ7Q]
GATE-09  The gate shall judge only the lane being left and the lane being entered; lanes skipped in between are not checked. [D28H7V, VW1WB1]
GATE-10  If a ticket's assignee is someone else, then a move shall be refused (not_owner), except out of a question or signoff lane.
GATE-11  When an unassigned ticket is moved, the move shall make the mover its assignee.

## Specification and dependencies
GATE-20  If a ticket crosses into the first requires-specified lane (todo) without goal, definition of done, context and assignee, then the move shall be refused naming each missing field.
GATE-21  If any forward move (parking excluded) starts work while a blocked-by entry is unresolved, then it shall be refused (exit 4). [2EF48S]
GATE-22  While a blocked-by entry `pr#<n>` names an unmerged pull request, the ticket shall be refused as waits_on_pr. [2EF48S]
GATE-23  If the target lane requires-option and the ticket's Options do not tick it, then the move shall be refused (step_not_chosen). [8WFNQY]
GATE-24  If the target is the blocked lane and neither --reason, blocked-reason nor blocked-by is given, then the move shall be refused.

## Outputs and outcome
GATE-30  If a forward move leaves a lane whose output-produces field is still empty, then the move shall be refused (missing_lane_output). [A1TZ4N]
GATE-31  If a move from a lane that does not produce the outcome would replace a filled outcome field, then it shall exit 3 (outcome_not_yours). [QTHD8H]
GATE-32  A move from a later lane shall be allowed to fill an empty outcome field. [QTHD8H]
GATE-33  If the target lane requires-outcome and outcome-what, -why or -resolves is empty, then the move shall be refused naming each.

## Definition of done and plan
GATE-40  `jaira dod <id> <n>` shall set one of four states: --todo [ ], --doing [~], --done [x], --superseded [-].
GATE-41  A done or superseded item shall count as settled; a doing item shall not.
GATE-42  When `--proof` is given, jaira shall record the evidence under the item, replacing an earlier proof line. [SDBCT5]
GATE-43  The `--plan` flag shall address the Plan checklist, which is separate: ticking the plan shall not tick the definition of done.
GATE-44  If the target lane requires-nonmodel-signal and a definition-of-done item is unsettled, then the move shall be refused, naming each item.
GATE-45  If the target lane requires-nonmodel-signal and the ticket has no definition-of-done checklist, then the move shall be refused.
GATE-46  If the target lane is terminal and a plan item is unsettled, then the move shall be refused (plan_incomplete). [NJBTVC]
GATE-47  When a forward move leaves the implementing lane with no DoD item settled, jaira shall warn on stderr and still move. [77B4GA]

## Commits
GATE-50  If the target lane requires-commits and the ticket records none, then jaira shall derive them from git (ticket file history and ids in subjects).
GATE-51  If git finds no commits either, then the move shall be refused and name branches that carry commits for the ticket.
GATE-52  When a move enters a requires-commits lane, jaira shall write the derived commits onto the ticket in the same write. [8FE2EG]

## Claim and release
GATE-60  `jaira claim <id>` shall take a 30-minute lease for this session; a live claim of another session is refused without --steal. [K1W2JR]
GATE-61  When a claim has expired, claiming shall take it over and report whose claim it was and how old. [K1W2JR]
GATE-62  `jaira claim --release` shall drop one's own claim at any time, and another's only once it has expired. [7WY0YT, K1W2JR]
GATE-63  If a ticket is rejected, then `jaira claim` shall refuse it and name `jaira pull` as the way back. [X8KZQ7]
GATE-64  `jaira release <id>` shall clear the assignee on the ticket's ref, then remove the local file; another's ticket needs --force.

## Leaving the board
GATE-70  `jaira reject <id> --why` shall refuse without --force (exit 3) and name the question lane move an agent makes instead. [A0C081]
GATE-71  When a person rejects with --force, jaira shall write rejected-kind (wontfix, obsolete, duplicate), -why, -by, -at and archive it. [A0C081]
GATE-72  If the rejected ticket lives only on its ref, then jaira shall reject it on the ref and write no file in the checkout. [4539MZ]
GATE-73  When a ticket reaches done on the shipped board, jaira shall not file it; only a person's `jaira logbook <id>` or `--all` files it. [74VM40]
GATE-74  If `jaira logbook <id>` names a ticket outside the terminal lane, then it shall refuse (exit 3, not_terminal). [74VM40]
GATE-75  When jaira files a ticket, it shall stamp its commits first and move it to `.jaira/logbook/<initials>-<yyyymmdd>/`. [74VM40]
GATE-76  Where a lane declares logbook-on-entry or holds, a move into it shall file the lane or trim its oldest tickets, and report each. [WXQ9PT, SGPDYK]
GATE-77  `jaira restore <file>` shall bring a filed or archived ticket back and take a rejection off it. [A0C081]
