# Cockpit: band and acceptance
_Owner: Sascha · EARS · each line one requirement · [tickets] hold the detail_

## Install and scope
COCK-01  The jaira-cockpit mod shall ship inside the jaira binary and be written as <skills>/jaira-cockpit/ by `jaira roles install`, for Claude Code only. [NZZDRH, E75R4Q]
COCK-02  While a session carries a ticket (its prompt started a worker, dispatcher or lane role), the cockpit shall draw no band and shorten no report. [NZZDRH, N7JW2Q]
COCK-03  Every text the cockpit shows a person shall be English; only the ticket's own texts keep the language they were written in. [F80QZT, FMWBA8]
COCK-04  When the board is found and the installed jaira is older than 0.3.7, the cockpit shall say both versions once and warn that Give back is refused. [WP1CHK]

## The band above the prompt
COCK-10  The band shall show in a frame titled jaira-cockpit with a head line "N worker | N review | N waiting for you" and the session's usage. [NZZDRH, QRHB35]
COCK-11  The band shall reread the board every 60 seconds and after every turn, so a ticket moved by another session appears without a command. [NZZDRH, 7QDQ4W]
COCK-12  The band shall list as waiting every ticket in a person's lane (human_exit or requires_question) on every board of the repository. [NC323D]
COCK-13  Where `jaira whoami` names the person, the band shall give buttons to own and unassigned tickets and fold teammates' into one line, Discuss only. [NZZDRH]
COCK-14  The band shall show at most 3 waiting tickets, one line each, and a "+N more" button that opens all of them. [PDZ444]
COCK-15  Each band row shall show id, lane and marks: DoD tick, bughunt tick, blocker count, note count, an old-review round and "?" for an open question. [PDZ444, AJ41SC]
COCK-16  If the bughunt fails only because the checkout lacks the hunt's commit, then the band shall show "?" for it, not a red cross. [AJ41SC]
COCK-17  Each band row shall carry one Acceptance button; Read explanation appears once an explain note exists, Discuss once a question is set. [2V50Y5, 13PQB3]
COCK-18  When a press on a row runs, the cockpit shall dim that row's buttons and ignore further presses on it until the action ends or 20 seconds pass. [0EW5XK, N5ZZP6]
COCK-19  The band shall show the last action's notice on top and teammates' tickets last. [TWB0R6]
COCK-20  While the context reaches 100k tokens the band shall advise /compact, and from 300k tokens to start a new session. [QRHB35]
COCK-21  While jaira's update cache names a newer release not dismissed, the band shall offer Update and Dismiss; JAIRA_NO_UPDATE_CHECK=1 hides it. [EDSQZP]
COCK-22  When Update is pressed, the cockpit shall show the `jaira update --roles --lanes --dry-run` preview with Cancel focused before it writes anything. [EDSQZP, PWAF22]

## The Acceptance pane
COCK-30  When Acceptance is pressed, the cockpit shall open one focused pane for that ticket that replaces separate Accept, Give back and Verdict buttons. [2V50Y5]
COCK-31  The pane shall open with a light: red to give back, yellow to look first, green to accept, naming its first reason and how many follow. [AJ41SC]
COCK-32  The light shall be red for a refusing verdict, a blocker, a bughunt not on HEAD or an open non-person DoD item, all read off the ticket. [AJ41SC, GY8T52]
COCK-33  The pane shall show goal, what was done and why, round and returns, DoD count, bughunt, verdict, commits, an existing pull request and the release note. [2V50Y5]
COCK-34  The pane shall separate its blocks by a coloured title and a dim rule, so it reads by skimming. [MK8HD2]
COCK-35  The pane shall list each review note with its class mark and the choices ok, fix and ticket; a blocker starts as fix, every other note as ok. [MFYM1A, GY8T52]
COCK-36  A note's class shall be the one the note states (blocker:, note: edge|limit|elsewhere), else the Explain note's, else edge. [GY8T52]
COCK-37  When a note was decided at an earlier acceptance, the pane shall preselect that decision and keep its follow-up ticket. [MFYM1A]
COCK-38  When "?" beside a note is pressed, the cockpit shall show one haiku sentence on why it matters, from the ticket's own diff of the named files. [2V50Y5]
COCK-39  Where settings.json names a language, Translate (key t) shall show head and notes in it via one haiku call; a second press returns the original. [FMWBA8, 0GMYQB]
COCK-40  The Translate button shall never write to the ticket and every opening of the pane shall start in the original. [FMWBA8]
COCK-41  The pane shall list open "seen by a person" DoD items under "Check yourself" for the person to tick. [DKMG91, 3DVVFZ]
COCK-42  The pane's last line shall say what Execute will do before it is pressed. [2V50Y5]
COCK-78  Where the next lane does not end the work, the pane's last line shall name it, as "Execute: accept → <lane>". [M5JSE9]

## Execute
COCK-50  When a note is marked fix or words are typed, Execute shall give the ticket back to its back lane (COCK-81) with those notes and words as the reason. [2V50Y5, MFYM1A, PGXJG0]
COCK-51  If a blocker is not marked fix and nothing is typed, then Execute shall refuse and ask why the blocker is accepted. [MFYM1A]
COCK-52  If "seen by a person" items remain unticked, then Execute shall refuse to accept and list them. [DKMG91]
COCK-53  If jaira names no next lane for the ticket, then Execute shall accept nothing and ask for a note marked fix or typed words. [2V50Y5, NC323D, M5JSE9]
COCK-54  When Execute accepts, the cockpit shall tick each seen item with proof "seen by <name> in the cockpit", unless its text changed since opening. [DKMG91]
COCK-55  When Execute runs, the cockpit shall create one follow-up ticket, following this one, per note marked ticket and say each that failed. [2V50Y5]
COCK-56  If the ticket left its lane since the pane opened, then Execute shall write nothing and say who moved it. [MFYM1A]
COCK-57  The pane shall stay open with the person's marks on any refusal and close once the ticket moved. [2V50Y5]
COCK-79  If a ticket in a lane that does not end the work still has an open question, then Execute shall refuse to accept until it is answered. [M5JSE9]

## Accept
COCK-60  The Accept step shall move into the ticket's next lane, terminal or not, run `move --dry-run` first and pass with --force only the needs_human refusal. [NZZDRH, GPFA7B, M5JSE9]
COCK-61  When a press is about to write, the cockpit shall first pull a ticket that is ref-only or stale on disk and remove a copy it pulled into the main checkout after. [76ERHH, DF46B9]
COCK-62  When the ticket has no own worktree, Accept shall move it in the worktree of the one branch holding its newest commit beyond HEAD, refusing if several do. [VHGWQF]
COCK-63  When Accept succeeds, the cockpit shall note "accepted by the person" on the ticket and submit one line to this session naming the move. [409XCV]
COCK-64  When Accept moved a ticket whose branch is ahead, has no pull request and all its tickets are done, the cockpit shall start the pr lane via spawn.sh. [GPFA7B, CRASZS]
COCK-65  Where `.jaira/pr` on the base says none, Accept shall start the lane that lands the branch on the default branch instead of a pull request. [XB97NJ]
COCK-66  Where `.jaira/pr-when` says signoff, Accept shall open the pull request first, keep the ticket in signoff and accept only a settled review round. [JG3REW]
COCK-67  The Accept step shall never start the next ticket; which ticket runs next is the teamlead's call. [X6JHQX]
COCK-68  If the pull request cannot be opened, then the ticket shall stay in done and the cockpit shall say why. [GPFA7B]
COCK-80  Only an Accept into the terminal lane shall start the pr or land lane or hold the ticket in signoff; an Accept mid-board moves it one lane on. [M5JSE9]

## Give back
COCK-70  The Give back step shall move to the ticket's back lane first, then note "Given back from <lane>: <reason>"; a refused move writes nothing. [NZZDRH, WP1CHK, PGXJG0]
COCK-71  When given back, the work shall return to the old worker while its tab lives, cache is warm and context is under workerContextShare (default 50%). [QRHB35, P122GE]
COCK-72  While the old worker is working, Give back shall type into its recorded pane, or else name its tab, and never start a second worker. [QRHB35, QTQW15]
COCK-73  If the old worker waits at an approval dialog, then Give back shall type nothing, start nobody and say where to answer. [M31BYV]
COCK-74  If Herdr's tab list cannot be read, then Give back shall start no worker and type nothing. [W80HNH, 8NWE4P]
COCK-75  When a new worker is needed, Give back shall close the old tab first and start it on the branch carrying the ticket's commits. [QTQW15]
COCK-76  A worker Give back starts shall report to the ticket's noted parent session, else this session, else none. [T5QN9P, 7V1XPB]
COCK-77  When Give back has moved a ticket, the cockpit shall note who got the work and why and submit one line to this session. [409XCV]
COCK-81  The back lane shall be, walking back from the person's lane, the first agentic lane that is not fresh and not a person's, or its rejects-to target; else in-progress. [PGXJG0]
