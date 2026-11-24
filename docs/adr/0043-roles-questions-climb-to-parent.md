# Questions climb to the parent session, addressed by its socket
Status: Accepted
Date: 2026-10-23

## Context
Workers asked the person in their own tabs; questions sat unseen across many tabs.
A parent named by session name went stale after a restart or compaction, and reports went nowhere.

## Decision
A worker or dispatcher with --parent notes its question on the ticket and sends it to the parent; only the top session asks the person.
spawn.sh passes the parent as uds:$CLAUDE_CODE_MESSAGING_SOCKET and checks it before a tab exists.
The teamlead decides what ticket, code or a rule settle (`decision (teamlead <date>)`), and asks the person only otherwise.

## Rejected
- Passing the session name (changes on restart).
- Replacing an unreachable parent with another peer (another teamlead would get the reports).
- Sending every third send-back straight to the human lane (K6PBYR: teamlead first).

## Consequences
An unreachable parent means the worker asks in its own tab and notes it; nothing is lost silently.

## Sources
TECZH5, 7V1XPB (decisions 2026-10-23), K6PBYR, T5QN9P, TWRAP8
