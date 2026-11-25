# Storage: ticket files, refs, filing, boards and tags
_Owner: Sascha · EARS · each line one requirement · [tickets] hold the detail_

## Layout and init
STOR-01  The store shall keep all repository state under .jaira/: tickets/, archive/, logbook/, milestones/, lanes/, boards/, tags. [BNZERQ, S0EG0M]
STOR-02  The store shall keep session state and lock files outside the repository, under ~/.jaira/state/<tree> or $JAIRA_HOME. [JNYE3T, AK3Z6D]
STOR-03  When jaira init runs, jaira shall create .jaira/tickets/ and the board's lanes, and running it again shall change nothing already there. [QA3GN1]
STOR-04  When jaira init makes a board private, jaira shall exclude .jaira/ in <git-common-dir>/info/exclude, or in .gitignore outside git. [NNV39Y]
STOR-05  When a fresh clone of a shared board runs jaira init or jaira update, jaira shall leave .gitignore unchanged so the board stays shared. [ETR0PX, 8W1R94]
STOR-06  When jaira share runs, jaira shall drop its exclude, ignore /.jaira/lanes/ and /.jaira/milestones/, and set up the merge driver. [NNV39Y, 1MYFKE]
STOR-07  If jaira share runs outside a git repository, then jaira shall refuse with exit 3 and reason not_a_repo.
STOR-08  If a rule jaira did not write still ignores .jaira/ after jaira share, then jaira shall warn on stderr and name that rule. [NNV39Y]
STOR-09  When jaira share --undo runs, jaira shall exclude .jaira/ again and print the git rm --cached line, never removing committed tickets itself.
STOR-10  The share and init commands shall never create a commit; committing is left to the person. [V3S0MW]

## Ticket file format
STOR-11  The store shall keep one ticket per file, .jaira/tickets/<ULID>-<slug>.md, YAML frontmatter then markdown body, with no counter file.
STOR-12  The store shall derive a ticket's identity from the frontmatter id, never from the filename; the slug is not renamed on title edits.
STOR-13  The store shall cut a slug at 48 bytes on a UTF-8 rune boundary, so a filename is never invalid UTF-8. [BWT53T]
STOR-14  The CLI shall resolve an exact id, unique prefix or unique suffix (the handle is the last six characters), else exit 5, reason not_found.
STOR-15  When jaira writes one field, the store shall splice only that value and leave every other byte, key order and comment unchanged. [KA9CFA]
STOR-16  The store shall keep frontmatter keys it does not know, so a field written by a newer version or a plugin survives every write.
STOR-17  If frontmatter uses YAML anchors or aliases, then the store shall refuse to rewrite the file, exiting 3 with reason unsafe_yaml.
STOR-18  The store shall read Definition of Done, Plan and Options as checkboxes [ ], [~], [x], [-] under their body headings; an unknown marker counts as open.
STOR-19  If one ticket file cannot be read, then listing shall show every other ticket and report the broken file by name.
STOR-20  If two files carry one id, or a ticket lies both on the board and in logbook/archive, then the store shall report it and move nothing.
STOR-21  While the board runs, the store shall re-read only ticket files whose size or modification time changed since the last listing. [NA413Y]

## Writes and concurrency
STOR-22  When a ticket is written, the store shall take a per-ticket lock file, re-read the file under it and replace the file atomically. [8X9SHE]
STOR-23  If a lock file is older than 30 seconds, then the store shall break it; if a live lock stays taken for 5 seconds, the write shall fail.
STOR-24  When the store writes a ticket, it shall stamp updated-at and updated-by with the writer's name. [DVPAPE]
STOR-25  If the ticket file on disk is older than its ref, then a write shall be refused with exit 3, reason stale_file, naming jaira pull. [F3VNVC]
STOR-26  While a ticket file lies behind its ref, jaira list, show and the board shall present the ref's state and mark it behind-ref. [K6BGB0]

## Tickets on refs
STOR-27  Where the board has a usable remote, jaira shall carry each ticket's full file on refs/jaira/tickets/<id>, outside refs/heads. [8566KF]
STOR-28  When a ticket is written, jaira shall queue the bytes in an outbox and send them after the command, never waiting on the network mid-write. [8566KF, WV2J2W]
STOR-29  When queued writes are sent, jaira shall push each ref as a compare-and-swap and, if it loses, name who wrote the winning state. [8566KF]
STOR-30  When jaira create succeeds on a board with a usable remote, jaira shall send the ticket to its ref and remove the local file. [V3S0MW, VM0A76]
STOR-31  If the remote is missing or unusable, then jaira create shall keep the file and print why, and --json shall carry file-only-reason. [VM0A76]
STOR-32  The store shall list tickets that exist only on a ref as read-only cards beside the files, in list, next, show and the board. [PXYN99]
STOR-33  If a command must change a ticket that exists only on its ref, then it shall refuse with exit 3, reason on_ref_only, naming jaira pull. [4539MZ, X8KZQ7]
STOR-34  When jaira pull <id> runs, jaira shall record the caller as assignee on the ref first and only then write the file. [RFC7GA]
STOR-35  If someone else already holds the ticket on its ref, then jaira pull shall refuse and name them, unless --steal is given. [RFC7GA]
STOR-36  When jaira pull runs on a file behind its ref with no unsent write, jaira shall bring the file up to the ref and print both states. [F3VNVC]
STOR-37  When jaira release <id> runs, jaira shall clear the caller as assignee on the ref first, then remove the local file. [VM0A76]
STOR-38  The CLI shall pick the remote as: git config jaira.remote, then settings remote if it exists here, then the only remote, else that name. [9ET6NC]
STOR-39  When a command finishes and the last fetch is older than fetch-every (default 10m), jaira shall start a detached fetch. [PP5SCQ, YYR59F]
STOR-40  When a fetch brings a ticket newly assigned to the user, jaira shall raise a desktop notification unless notify-off is set. [8566KF]
STOR-41  When jaira fetch finds a ref it saw before gone while the file is still here, it shall report the ticket departed and move nothing. [RFC7GA, RGFJ0Q]
STOR-42  When jaira fetch runs, jaira shall report tickets filed here that reached no landing branch within landing-grace (default 72h). [RA7PFE]

## Snapshot and reaping
STOR-43  When the last snapshot is older than snapshot-every (default 72h), jaira shall start a detached jaira snapshot after a command. [PTQ3XT, YYR59F]
STOR-44  The snapshot shall write every ticket on a ref under board/ on branch jaira/board via plumbing, never checked out, plus the people file. [PTQ3XT]
STOR-45  If the set of refs is unchanged since the last snapshot, then jaira snapshot shall write no commit and push nothing. [PTQ3XT]
STOR-46  The snapshot clock shall be kept per clone, so a new worktree does not snapshot on its first command. [NWFE2B]
STOR-47  When a snapshot is on the remote, jaira shall delete the refs of tickets filed in logbook/ or archive/ on a landing branch. [RA7PFE]
STOR-48  If no landing branch resolves from landing-branches or the remote HEAD, then jaira snapshot shall delete no ref. [RA7PFE]
STOR-49  When jaira snapshot --drop <id> runs, jaira shall delete that one ref without waiting for its ticket to land.

## Filing: logbook, archive, reject, restore, delete
STOR-50  When jaira logbook <id> runs on a terminal-lane ticket, jaira shall stamp its commits and move it to logbook/<initials>-<yyyymmdd>/. [74VM40]
STOR-51  If the ticket is not in the terminal lane, then jaira logbook shall refuse with exit 3, reason not_terminal, naming the move. [74VM40]
STOR-52  The board shall never file a ticket into the logbook on its own; reaching the terminal lane files nothing. [74VM40]
STOR-53  When jaira logbook --all runs, jaira shall file every terminal-lane ticket into today's folder and leave milestones in place. [74VM40]
STOR-54  When jaira logbook runs without an id, jaira shall list what was filed within logbook-days (default 28); --since overrides it. [MJ3PVV, WKVMC6]
STOR-55  When jaira archive <id> runs, jaira shall stamp commits best-effort and move the file to .jaira/archive/ from any lane. [2MM32Y]
STOR-56  When a ticket is filed or archived, jaira shall write its final state to its ref and keep the ref until the ticket lands. [RA7PFE]
STOR-57  When jaira reject <id> --why runs with --force, jaira shall write rejected-kind, -why, -by, -at and archive the ticket. [A0C081]
STOR-58  If jaira reject runs without --force, then it shall refuse with exit 3 and point to the human lane with a question. [A0C081]
STOR-59  If the rejected ticket exists only on its ref, then jaira shall write the rejection on the ref, clear the assignee and create no file. [4539MZ]
STOR-60  The board and jaira list shall hide every ticket carrying rejected-kind, whether it is a file or on a ref. [4539MZ]
STOR-61  When jaira restore <file> runs, jaira shall move the file back from archive/ or logbook/ by basename and remove the rejection fields. [A0C081]
STOR-62  When jaira delete <id> runs, jaira shall remove the file and its ref only after the handle is typed back or --force is given.
STOR-63  If other tickets still point at the ticket, then jaira delete shall refuse unless --force is given.
STOR-64  If the archive or logbook folder already holds a file of that name, then filing shall refuse rather than overwrite.

## Merging
STOR-65  When git merges a ticket file on a shared board, the merge driver shall union blocked-by, commits, tags and related. [NQCKF5]
STOR-66  When both sides moved the ticket, the merge driver shall keep the status further along the lane order; other scalars follow recency.
STOR-67  If both sides changed one prose field differently, then the merge driver shall keep ours and park theirs in conflict-theirs-<field>.
STOR-68  If the merge base is empty (add/add), then the merge driver shall merge as if the base had no fields. [NQCKF5]

## Boards
STOR-69  The store shall treat .jaira/lanes/ as the board named default and each .jaira/boards/<name>/lanes/ holding a lane file as a named board. [S0EG0M]
STOR-70  The store shall record a ticket's board in its board: field; an absent field means the default board. [S0EG0M]
STOR-71  If a command names a board with no lane files, then it shall exit 2 with reason no_such_board. [S0EG0M]
STOR-72  When jaira list or next runs without --board, jaira shall show only the default board's tickets. [S0EG0M]
STOR-73  When jaira boards runs, jaira shall list each board with its open-ticket count and its lanes/description text, writing nothing. [8NHD8Q]

## Tags
STOR-74  The store shall keep tags as a frontmatter list and tag colours in .jaira/tags, one "name: <ansi256>" line per tag. [79GEPW]
STOR-75  When a tag is written, jaira shall store it as lowercase kebab and say so; a name with other characters shall be refused, not trimmed. [79GEPW]
STOR-76  When a tag new to the board is added, jaira shall append one line with a free palette colour and leave every other line untouched. [79GEPW]
STOR-77  When jaira tags runs, jaira shall list every known tag with its colour and the count of open tickets carrying it. [79GEPW]
