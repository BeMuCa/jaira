# Cockpit: chats, reports, threads and workers
_Owner: Sascha · EARS · each line one requirement · [tickets] hold the detail_

## Explain and Discuss
COCKW-01  Where Herdr runs, Explain all and Discuss shall open a chat in a Herdr tab `<id>/explain` or `<id>/discuss` via spawn.sh, in the ticket's worktree. [NZZDRH, 58DCYM]
COCKW-02  Where Herdr is absent, Explain shall ask this session to start one subagent that writes the explain note, and Discuss shall talk here. [42C3RZ]
COCKW-03  An Explain chat shall write one `explain:` note: a status line, one classed line per review note, 2-4 numbered options and a recommendation. [13PQB3, GY8T52]
COCKW-04  A Discuss chat shall note `decision:`, clear the question, and move the ticket back to in-progress only when the decision means more work. [NZZDRH]
COCKW-05  When a chat's note is on the ticket from its opening minute on and its agent is out of its turn, the cockpit shall close its tab. [NZZDRH]
COCKW-06  If a chat already holds a pulled copy of the ticket, then a second Explain or Discuss on it shall be refused until that chat ends. [DF46B9]
COCKW-07  When Read explanation is pressed, the cockpit shall show the options one row each, picked by digit, with "Own answer…" after the last. [13PQB3]
COCKW-08  When an option is picked, the cockpit shall note `decision:` and then give back for [give back], accept for [accept], else tell this session. [13PQB3, MFYM1A]
COCKW-09  When an own answer is submitted, the cockpit shall note it as a decision and give the ticket back with it as the reason. [13PQB3]
COCKW-10  The cockpit shall read whether an explanation is answered off the ticket (a decision note after it), not off session state. [13PQB3]
COCKW-11  When a Discuss chat noted a decision, the band shall show "<id> | decided" with Read decision until it is read. [NZZDRH]

## Agent reports in the chat
COCKW-20  When an agent report or task-notification names a ticket, the cockpit shall replace it in the chat with one line led by "[shortened by jaira-cockpit]". [T8QS1Y, QEPS3M]
COCKW-21  If a report is a failure or names no ticket, then the cockpit shall pass it on whole. [T8QS1Y]
COCKW-22  Where a report ends in a status line, the one-liner shall read lane move, verdict, branch, commits, gates and notes from that line alone. [MVR4PY, 55CS2M, TT3Z11]
COCKW-23  The one-liner shall add a question or "needs you" the report contains and end with "full report: /cockpit report <id>". [T8QS1Y]
COCKW-24  The cockpit shall keep the full text per project: 10 reports per ticket, 20 tickets, 12,000 characters per report. [T8QS1Y]
COCKW-25  The report tool and `/cockpit report <id>` shall return only reports not read before, newest first; `all` adds the older ones. [0J88M5]
COCKW-26  When it is called without a ticket, the report tool shall list the tickets that have kept reports, most recent first. [T8QS1Y]
COCKW-27  The cockpit shall rewrite only rows arriving by the prompt and delivery doors of session.append. [QEPS3M]
COCKW-28  When `/cockpit ticket <id>` runs, the cockpit shall open a pane with the ticket's reports, marked read, and its usage rows from metrics.csv. [NZZDRH]

## Threads and "Waiting for you"
COCKW-30  When the person finishes a turn, the cockpit shall file it into a conversation thread with one haiku call, never per agent report. [NZZDRH, 1M3X77]
COCKW-31  A thread shall be marked "?" while a question waits for the person, "->" once a ticket holds it and "~" after 5 turns without either. [1M3X77]
COCKW-32  The band shall show "threads: Q open question(s), S left lying"; while Q > 0 it is a button, also on key q. [1M3X77]
COCKW-33  The "Waiting for you" pane shall list unanswered questions, last asked first, each with its ticket, title and options, the first nine picked by digit. [1M3X77]
COCKW-34  When an option is picked, the cockpit shall send it as the person's message and drop the question at once. [1M3X77]
COCKW-35  When "Answer in the chat…" is pressed, the cockpit shall close the pane and fill the prompt with what the answer is for. [1M3X77]
COCKW-36  The cockpit shall keep threads per session, for the 20 sessions heard from last. [1M3X77]

## Cache guard
COCKW-40  If the prompt cache went cold and the context is at least cacheGuardTokens (default 30k), then a person's prompt shall be held. [QRHB35]
COCKW-41  A held prompt shall offer Send, Compact first, Start over (/clear) and Cancel, with the cost at 1.25 × usdPerMTok; other than Send keeps the words. [QRHB35]
COCKW-42  The guard shall hold no slash command, no text typed by the cockpit ([jaira-cockpit]) and no agent message; with nobody to ask it sends. [QRHB35]

## Workers
COCKW-50  The workers tool and `/cockpit workers` shall list this session's context and cache, then each agent carrying a ticket with lane, context and cache. [QRHB35]
COCKW-51  The list shall say for each worker whether the ticket's next round goes to it or to a new one, and why, by the rule Give back uses. [QRHB35]
COCKW-52  Every session with a Herdr pane that carries a ticket shall record its pane, ticket, lane, context and cache; a record older than 24 h is dropped. [QRHB35]
COCKW-53  The band shall count as worker or review only a Herdr tab in its turn or a running subagent that names a ticket of this board. [NZZDRH, F80QZT]

## Overview
COCKW-60  The band's Overview button shall open a pane with the tabs Queue, Log and Workers, read afresh at each press. [DSKZ70, RJFA19]
COCKW-61  The Queue tab shall list In work, Queue and Put back beneath each other, each folding at its heading. [4ZM13S]
COCKW-62  The Log tab shall show what was created and put back in the last 7 days, by day, with who did it. [DSKZ70, 02D92C]
COCKW-63  The Log shall filter by people and by lanes with chips that cycle off, chosen and left out; All clears its row. [02D92C, 96E1Z4]
COCKW-64  The Workers tab shall draw each agent as a figure at one stone whose size follows the open DoD items of tickets with an agent. [F80QZT, RNZYJR]
COCKW-65  Each Workers line shall show ticket, lane, kind, state (working, resting, waits for an answer, waits for approval), time seen and last report. [F80QZT, RJFA19]
COCKW-66  A ticket shall keep its colour and hat while its agents run, and two running tickets never share a colour while colours suffice. [F80QZT]
COCKW-67  The Workers tab shall animate every 250 ms only while it is open and shown. [F80QZT]
COCKW-68  Every cockpit pane's head shall be one row of equal framed buttons, tabs left and Close right. [R19PY9, HQHFXQ]

## Pull requests and cleanup
COCKW-70  When the pr lane Accept started opens a pull request, the cockpit shall find it within 20 one-minute looks and tell this session once. [CRASZS]
COCKW-71  When a report names a pull request of this repository that is open, the cockpit shall tell this session to follow it, once per pull request. [CRASZS]
COCKW-72  The line handed over shall say to run `jaira pr --watch --number N` under a Monitor; the cockpit itself polls no pull request. [CRASZS]
COCKW-73  The cockpit shall read GitHub Enterprise SSH remotes and call gh with host/owner/repo. [XSGWPX]
COCKW-74  Where Herdr runs, a person's session shall run reap.sh from its first minute and then at most every 10 minutes. [N7JW2Q]
COCKW-75  When a session ends with a board, the cockpit shall append its token totals as a row to ~/.jaira/metrics.csv. [NZZDRH]
