# TUI Board
_Owner: Berk · EARS · each line one requirement · [tickets] hold the detail_

## Start and navigation
TUI-01  When `jaira` runs without a command, the CLI shall open the launcher, not a board; `jaira board` shall open the board directly. [9WE1QK]
TUI-02  The board shall move between lanes with ← →, between cards with j k ↓ ↑, and to the first/last card with g G. [9WE1QK]
TUI-03  When q is pressed on any page, the TUI shall go back one page to where esc goes; on the board q shall return to the launcher or quit. [9WE1QK]
TUI-04  While a text field has focus, the TUI shall treat q and every other letter as text. [9WE1QK]
TUI-05  The TUI shall read a command key off the physical key, so commands work on a Cyrillic layout while typed text keeps that layout. [Z4W2BX]
TUI-06  When text is pasted into a single-line field, the TUI shall insert it with line breaks folded to spaces and Unicode intact. [8XHZT5]
TUI-07  The TUI shall wrap prose to the terminal width in every view rather than cut it off at the right edge. [HREQJR, MR3CN8]
TUI-08  When ? is pressed, the board shall show every key, the meaning of every card mark and the roadmap keys. [9WE1QK]

## Cards and lanes
TUI-10  The board shall draw a card as a three-row filled band without a frame, its background alternating from card to card. [QQ3EX4]
TUI-11  The board shall show the colours of a ticket's first three tags as cells of a one-column bar at the card's left edge. [S1VM40]
TUI-12  The board shall fill the selected card in a shade that follows the terminal's light or dark background. [ANX7Y5]
TUI-13  When t is pressed, the board shall show a legend mapping each tag colour to its name. [AXTFG3]
TUI-14  While an agent, a fresh reviewer or a person holds a ticket, the board shall run a band round its card; c shall stop or start it. [C8V6DQ, 2C2TFC]
TUI-15  The board shall mark a card with glyph plus word (spec, blocked, asks, sign off, unworked, ✎ name), never by colour alone. [S1VM40]
TUI-16  When z is pressed, the board shall draw lanes holding no ticket thin instead of hiding them; pressing z again widens them. [QPJNQP]
TUI-17  While the cursor rests on a card with links, the board shall draw a coloured line to each linked card on screen. [PAP369]
TUI-18  The board shall show the running version top left, `jaira dev` on a dev build and `· up to date` when no newer release is known. [P1AE82, DNAEPN]
TUI-19  If the release cache knows a newer release, then the launcher and board head shall show a pill `^ <version>` under the version line. [T8A8KM]

## Filtering
TUI-20  When / is pressed, the board shall filter as you type over id, title, goal, context, DoD, assignee, lane, tags, body and checklists. [N172JS]
TUI-21  Where the filter holds `key:value`, the board shall narrow to that field; tag, milestone and user keys shall match exactly. [N172JS, PNC854]
TUI-22  When u or M is pressed, the board shall open one ticking picker (space ticks, enter applies, x clears) for people or milestones. [N172JS, PNC854]
TUI-23  When esc is pressed on the board, the board shall clear the filter.

## Ticket detail
TUI-30  When enter is pressed on a card, the TUI shall open the ticket; ↓ ↑ shall scroll it and j k shall jump to the neighbouring ticket. [9WE1QK]
TUI-31  The detail view shall keep the author's line breaks of a multi-line field, one wrapped line per input line. [KA9CFA]
TUI-32  The detail view shall list every field an installed lane declares; an empty one shall read `— owed by <lane>`. [A1TZ4N]
TUI-33  The detail view shall prefix a filled lane field with the lane(s) that declare it, e.g. `(critique/review)`. [D9SC5A]
TUI-34  When b is pressed on a blocked ticket, the TUI shall open the ticket it waits on; a wait on a pull request shall be said, not opened.
TUI-35  When l is pressed, the TUI shall open a window with every link in both directions, across board, refs, logbook and archive. [51RNC1]
TUI-36  When y is pressed in the detail view, the TUI shall copy the full ticket id to the clipboard.
TUI-37  When e is pressed, the TUI shall edit fields in place; E shall hand only the body and checklists to $EDITOR.
TUI-38  When n is pressed on an open ticket, the TUI shall open a follow-up beside it in a split; tab shall switch panes and esc discard.

## Changing tickets
TUI-40  When n is pressed on the board, the TUI shall create a ticket and open it for filling in.
TUI-41  When m is pressed, the TUI shall move the ticket through the same gate checks and core move as the CLI. [6BFHJ4]
TUI-42  If the gate refuses a move, then the TUI shall list the refusals and override them with f only after the person confirms with y.
TUI-43  When a person moves a ticket from a human-exit lane to an earlier lane, the TUI shall require a reason and write it as the note `sent back to <lane> from <lane>: <reason>`. [BQGM6X]
TUI-44  When a is pressed on a ticket in a human-exit lane, the TUI shall move it to the terminal lane as a person, through the gate. [6BFHJ4]
TUI-45  When f is pressed on a ticket in a human-exit lane, the TUI shall raise a follow-up instead of accepting.
TUI-46  When x is pressed, the TUI shall archive the ticket; R shall reject it with a one-line reason. [A0C081]
TUI-47  When X is pressed on an open ticket, the TUI shall delete its file only after the person types its handle back.
TUI-48  If a gesture would change a logbook ticket, then the TUI shall refuse it and say the ticket is read-only. [WKVMC6]

## Live data
TUI-50  The terminal lane shall show, below its tickets, the logbook tickets of the last N days marked `⎙ filed <date>`; N=0 hides them. [WKVMC6]
TUI-51  The board shall re-read only ticket files that changed and parse each ticket head once. [NA413Y]
TUI-52  The board shall fetch ticket refs at startup and then every `fetch-every` from settings.json (default 10m), as the CLI does. [YYR59F]
TUI-53  When a fetch brings tickets newly assigned to me, the board shall say so in the status line, never as a screen. [YYR59F]
TUI-54  If a fetch fails or a fetched milestone cannot be queued, then the board shall show the error instead of dropping it. [S1SA89]
TUI-55  When the board opens on a newer version than the user last read, the TUI shall show the unread release notes once. [GK1QGT]
TUI-56  The release-notes view shall show one version per page, newest first, ← → paging and `- !` lines highlighted. [GK1QGT]

## Boards, compact view, roadmap
TUI-60  When b is pressed, the board shall switch to the repository's next board with its own lanes and only its tickets. [S0EG0M]
TUI-61  When 1–9 or p is pressed, the board shall switch to another registered project.
TUI-62  When v is pressed, the TUI shall show the whole flow with agents counted per step; enter shall open one lane full-screen. [9WE1QK]
TUI-63  When L is pressed, the TUI shall show the Milestone Roadmap: one bubble per placed milestone, filled by its done share, arrows for `after`. [D0SAHM]
TUI-64  When a milestone is renamed in the roadmap's e dialog, the TUI shall rename it everywhere, as `jaira milestone rename` does. [E70WYN]
TUI-65  When x is pressed on a bubble and confirmed with y, the TUI shall take it off the roadmap and keep its tickets. [D0SAHM]

## Settings
TUI-70  When S is pressed, the TUI shall open settings: Lanes, Default board, Logbook days, Agents per ticket, Chat model, Notes language.
TUI-71  While the lane screen is open, enter shall add a lane not installed and remove an installed one behind a yes/no dialog that preselects no. [GDB5YJ]
TUI-72  While the lane screen is open, H and L shall move the selected lane one step, as `jaira lanes move` does. [9CQJA9]
TUI-73  The setting Agents per ticket shall toggle lanes/worker; any value but `worker` shall mean lanes. [G6T54H]
TUI-74  The logbook-days setting shall default to 28 and apply to both the board and `jaira logbook`. [WKVMC6]

## Launcher
TUI-80  The launcher shall list every registered repository once, under its main checkout, never a linked worktree or temp dir. [ZXJ72E]
TUI-81  The launcher shall count open and total tickets as `jaira list` does, ref-only tickets included, filled in after first paint. [PXYN99]
TUI-82  If projects.json cannot be parsed, then jaira shall copy it to `projects.json.broken[-<time>]` and say so before writing a new one. [BPR7C6]
TUI-83  When s is pressed in the launcher, the launcher shall toggle a seven-day chart of logbook entries per day over all boards. [SPDWGH]
TUI-84  While the launcher is open, a shall add boards from a directory browser, x remove one (default no), d edit the default board, w show notes. [9WE1QK]
