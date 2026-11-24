# People and planning: identity, people list, whoami, milestones, roadmap, epics
_Owner: Berk · EARS · each line one requirement · [tickets] hold the detail_

## Identity
PPL-01  The jaira CLI shall act as the person the board's people list gives for git user.name, else for git user.email, writing that person's jaira name. [DVPAPE]
PPL-02  If the people list knows neither name, then jaira shall act under git user.name, else $USER/$USERNAME/$LOGNAME, else "unknown". [DVPAPE]
PPL-03  Where JAIRA_USER is set, jaira shall act as that name, resolved through the people list, and shall not look the person up by git user.email. [Q030WZ]
PPL-04  The jaira CLI shall count as "me" the acting name, that person's people-list names, git user.name, git user.email and each line of ~/.jaira/identity. [DVPAPE]
PPL-05  While JAIRA_USER is set, jaira shall count as "me" only that name and its people-list names, never the checkout's git user.email or alias file. [Q030WZ]
PPL-06  The jaira CLI shall compare person names case-insensitively, ignoring surrounding space, and shall treat two empty names as different people. [DVPAPE]
PPL-07  If a ticket's assignee is not "me", then move shall refuse it with exit 3 naming the owner's jaira name, unless it sits in a human lane, the move hands it over, or --force is given. [Q030WZ]
PPL-08  The jaira CLI shall treat a ticket with an empty assignee as belonging to nobody and writable by anyone. [DVPAPE]
PPL-09  When a ticket is filed into the logbook, jaira shall use the folder <initials>-<yyyymmdd> of the acting name, the whole slug for a one-word name. [DVPAPE]

## whoami
PPL-10  When `jaira whoami` runs, the jaira CLI shall print the acting name, every other name it counts as "me", and the alias file path.
PPL-11  While JAIRA_USER is set, whoami shall mark the alias file "(not read under JAIRA_USER)", and --json shall carry alias_file_read false. [Q030WZ]
PPL-12  When whoami runs inside a board, it shall also print the board root, the remote used and its source, the remotes, ref mode and the file-only ticket count.
PPL-13  If whoami runs outside a board, then it shall still print the identity part and exit 0.
PPL-14  If the board is not in a git repository, then whoami shall omit the Remote and Remotes rows and the remote fields of --json.

## People list
PPL-15  The people list shall be the file `people` on refs/jaira/people, one person a line: `<jaira name>: <alias>, <alias>`. [DVPAPE]
PPL-16  When the people list is read, jaira shall skip blank lines, # comments and lines it cannot parse instead of failing. [DVPAPE]
PPL-17  When jaira writes the people list, it shall write it sorted case-insensitively under a fixed comment header. [DVPAPE]
PPL-18  When `jaira people` runs, the jaira CLI shall list each person's jaira name and aliases (--json: people, ref). [DVPAPE]
PPL-19  When `jaira people alias <name> <other>...` runs, jaira shall add each name to the person <name> means, starting a new person if nobody has <name>. [DVPAPE]
PPL-20  If an alias already names another person, then `people alias` shall refuse with exit 3 and write nothing. [DVPAPE]
PPL-21  When `jaira people rename <old> <new>` runs, jaira shall make <new> the person's jaira name and keep <old> as one of their aliases. [DVPAPE]
PPL-22  If rename's <new> is a name of another person, then it shall exit 3; if <old> names nobody, then it shall exit 5. [DVPAPE]
PPL-23  If a person name is empty, holds a colon, comma or line break, or starts with #, then the people commands shall refuse it with exit 2. [DVPAPE]
PPL-24  When a people write runs, jaira shall fetch, apply the change to the newest list and push with a lease, redoing it up to four times on a lost race. [DVPAPE]
PPL-25  If the remote is unreachable, then a people write shall write nothing and exit 1; without a repository or remote it shall exit 3. [DVPAPE]
PPL-26  The jaira CLI shall never rewrite a ticket file for an alias or a rename; readers show the person's current jaira name instead. [DVPAPE]
PPL-27  When `create --assignee` or `set <id> assignee=` names any alias, jaira shall store that person's jaira name. [DVPAPE]
PPL-28  When `list` or `next` runs with --assignee, jaira shall match tickets stored under any name of that person. [DVPAPE]
PPL-29  When `jaira snapshot` runs, it shall back up the people list as `people` beside board/. [DVPAPE]
PPL-30  When `u` is pressed on the board, the TUI shall open a picker that narrows the board to the ticked people. [N172JS]

## Milestones
PPL-31  A milestone shall be the file .jaira/milestones/<name>-<ID>.md: frontmatter, then one `- <ticket id>` line per member. [0YGWXQ, EBTX6X]
PPL-32  A milestone ID shall be six random characters of the ticket-handle alphabet with at least one letter, kept for the milestone's life. [EBTX6X]
PPL-33  Each milestone shall travel on refs/jaira/milestones/<ID> and reach teammates with their next fetch, without a commit. [EBTX6X, 1MYFKE]
PPL-34  When `share` runs, or `init`/`update` on a shared board, jaira shall add /.jaira/milestones/ to .gitignore. [1MYFKE]
PPL-35  If milestone files are tracked in git, then `update` and `validate` shall warn (milestones_tracked) with the `git rm -r --cached` line. [1MYFKE]
PPL-36  When a milestone is about to be written, jaira shall first bring a file that is older than its ref up to the ref. [1MYFKE, KWQCYZ]
PPL-37  Where the board has no remote, jaira shall still write milestones and queue them for their ref once one exists (fetch: milestones_published). [1MYFKE]
PPL-38  The jaira CLI shall store a milestone name as lowercase kebab by the tag rules and say so when it changed the name typed. [0YGWXQ]
PPL-39  When a milestone is created, jaira shall give it a random colour no other milestone uses, unless --color 1-255 is given (else exit 2). [0YGWXQ]
PPL-40  If a milestone on the board already has the name, then `milestone create` shall exit 3; a filed milestone does not hold its name. [EBTX6X]
PPL-41  Every milestone command shall take a milestone as name, ID or <name>-<ID>, and one that writes shall exit 3 on a name two milestones carry. [EBTX6X]
PPL-42  When `milestone add` or `rm` runs, jaira shall change all given tickets in one write and keep every other line verbatim. [0YGWXQ]
PPL-43  When `milestone add` names a ticket already in it, jaira shall leave it at its place; `rm` of the last ticket shall leave the milestone standing.
PPL-44  When `milestone rename` runs, jaira shall change only the name; ID, members, goal, colour and roadmap place stay. Taken name: exit 3. [E70WYN, EBTX6X]
PPL-45  When `milestone ls` runs, jaira shall print each milestone's colour, label, ticket count and goal (--json: milestones, count, dir, milestones_twins).
PPL-46  If two milestones share a name, then fetch and every milestone command shall warn on stderr with the rename that ends it, and refuse nothing. [EBTX6X]
PPL-47  If one ID lies in two files, then every milestone write, logbook and restore included, shall exit 3 (milestone_twin_files) naming both. [8TT890]
PPL-48  When `jaira logbook <milestone>` runs, jaira shall move its file under logbook/<folder>/milestones/ and write status: filed onto its ref.
PPL-49  If a milestone holds a ticket still on the board outside the terminal lane, then `jaira logbook <milestone>` shall exit 3 (milestone_unfinished).
PPL-50  When `jaira logbook --all` runs, jaira shall file no milestone.
PPL-51  If a milestone is filed, then milestone edits shall exit 3 and point to `jaira restore <file>`.
PPL-52  When `list --milestone <m>` runs, jaira shall show only that milestone's tickets, of both milestones where two carry the name. [EBTX6X]
PPL-53  The TUI shall paint a milestone's colour on the right edge of its cards; `M` shall tick several milestones to filter the board. [0YGWXQ, PNC854]

## Milestone Roadmap
PPL-54  Where a milestone has an `after:` line, the TUI shall show it on the Milestone Roadmap; without one it stays off. [D0SAHM]
PPL-55  When `--after a,b` is given, jaira shall store the predecessors as IDs; the first decides the swimlane, `--after ""` makes a swimlane top. [D0SAHM, EBTX6X]
PPL-56  If a predecessor is not a milestone on the board, then jaira shall exit 5; if the arrow closes a loop, then exit 3 naming the loop. [D0SAHM]
PPL-57  When `milestone set --no-roadmap` or `x` takes a milestone off, jaira shall hang its successors on its own predecessors. [D0SAHM]
PPL-58  If `milestone set` gets none of --goal/--after/--no-roadmap, or both --after and --no-roadmap, then it shall exit 2. [D0SAHM]
PPL-59  When `L` is pressed, the TUI shall open the roadmap: a bubble per placed milestone filled by its done share, arrows from first to next. [D0SAHM]
PPL-60  While on the roadmap, the TUI shall create (n) or connect (c) at a [+], edit name, goal and predecessors (e), recolour (c) and unplace (x). [D0SAHM, E70WYN]

## Epics
PPL-61  An epic shall live only on refs/jaira/epics/<name> as <name>.md, with no file in the working tree. [F914QH]
PPL-62  If the board has no repository with a remote, then every epic command shall exit 3 (no_refs). [F914QH]
PPL-63  When `epic create` runs, jaira shall store a kebab name, --title (default the name), a one-line --description and any --external keys. [F914QH, PG9P9Y]
PPL-64  When tickets are put into or taken out of an epic, jaira shall write the epic's ref only, never a ticket's. [F914QH]
PPL-65  The jaira CLI shall count an epic as done X of Y over its members and every follows:/parent: descendant; done is terminal lane or logbook. [F914QH]
PPL-66  When `epic ls --unassigned` runs, jaira shall list the open tickets in no epic, neither as member nor follow-up. [F914QH]
PPL-67  When `epic close` runs, jaira shall write status: closed onto the ref; unfinished work exits 3 unless --force. [182JT1]
PPL-68  While an epic is closed, `epic add` shall exit 3 without --force, `epic ls` shall hide it, and `epic reopen` shall take it back. [182JT1]
PPL-69  When `epic set` runs, jaira shall rewrite title, description, body or external keys and leave member lines untouched. [CRDH6B, PG9P9Y]
PPL-70  If `epic set` names nothing to set or an empty title, then it shall exit 2; an unknown epic exits 5. [CRDH6B]
PPL-71  When `epic ls --external <key>` or `list --external <key>` runs, jaira shall show only the epics, or their tickets, carrying that key. [PG9P9Y]
PPL-72  If two clones change one epic field at once, then jaira shall keep one value and put the other below `<!-- jaira: changed at the same time -->`. [CRDH6B]
PPL-73  Where an epic is named in output, jaira shall print its title with the name to type beside it. [PG9P9Y]
