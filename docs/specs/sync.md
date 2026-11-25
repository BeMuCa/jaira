# Sync: refs, outbox, fetch and concurrency
_Owner: Sascha · EARS · each line one requirement · [tickets] hold the detail_

## Refs
SYNC-01  The jaira CLI shall carry each ticket on its own ref refs/jaira/tickets/<id>, whose commit tree holds only <id>.md, the whole file. [8566KF]
SYNC-02  The jaira CLI shall carry milestones on refs/jaira/milestones/<ID>, epics on refs/jaira/epics/<name> and people on refs/jaira/people. [1MYFKE, F914QH, DVPAPE]
SYNC-03  When it writes a ref, the jaira CLI shall push with --force-with-lease on the SHA the remote last confirmed, committing on that SHA as parent. [8566KF]
SYNC-04  The jaira CLI shall move its local ref only after the remote accepted the push, so the local ref is always the next lease. [8566KF]
SYNC-05  If a push is refused as stale or non-fast-forward, then the jaira CLI shall treat it as a lost race, never as an unreachable remote. [8566KF]
SYNC-06  The jaira CLI shall keep its refs outside refs/heads, so the default refspec, branches and pull requests never carry them. [8566KF]

## Remote
SYNC-07  The jaira CLI shall pick the remote in this order: git config --local jaira.remote, settings.json "remote" if this repo has it, the repository's only remote, else the configured name (default origin) unchanged. [9ET6NC]
SYNC-08  While a board has no git, no repository or no usable remote, the jaira CLI shall queue nothing and keep every ticket as a file. [8566KF, VM0A76]
SYNC-09  If 'jaira fetch', 'jaira pull' or 'jaira release' runs without a usable remote, then it shall fail with exit 1 and reason no_remote. [RFC7GA]

## Outbox
SYNC-10  When a command writes a ticket, milestone or epic, the jaira CLI shall queue the whole content in the outbox and return without the network. [8566KF]
SYNC-11  The jaira CLI shall keep the outbox in the per-working-tree state directory outside the repository, never in a commit. [8566KF]
SYNC-12  When a write is queued for a key that already has an unsent one, the outbox shall replace the content and keep the first lease. [8566KF]
SYNC-13  When a command that queued a write finishes, also with an error, the jaira CLI shall flush the outbox once. [8566KF]
SYNC-14  While a command has queued nothing, the jaira CLI shall not flush or touch the network on its path, so read commands stay instant. [8566KF, PP5SCQ]
SYNC-15  When it flushes, the jaira CLI shall fetch refs/jaira/* first, then send the entries oldest first. [8566KF]
SYNC-16  If the remote is unreachable during a flush, then the jaira CLI shall stop sending and keep every entry for the next command. [8566KF]
SYNC-17  If the remote refuses a queued ticket write as a lost race, then the jaira CLI shall drop the entry and name who wrote the ticket first. [8566KF]
SYNC-18  If the remote refuses a milestone write, then the jaira CLI shall merge the edit into the winner, write the file and queue it again. [1MYFKE]
SYNC-19  If the remote refuses an epic write, then the jaira CLI shall merge both sides' members and resend up to three times, else keep it queued. [F914QH, S1SA89]
SYNC-20  If git refuses a queued write for any other reason, then the jaira CLI shall keep the entry and report it as failed. [8566KF]
SYNC-21  If the outbox cannot be listed during a flush, then the jaira CLI shall return that error and leave the queue untouched. [MMKENA]
SYNC-65  If a session queues a write for a key while another session is sending that key, then the newer write shall stay queued for the next command. [WV2J2W]
SYNC-22  While the board TUI runs, it shall flush its outbox before each background fetch, so board edits reach the team without a CLI call. [17TJ1H]
SYNC-23  While a ticket has an unsent write, the board shall mark its card as unsent. [8566KF]

## Ref-only tickets
SYNC-24  When 'jaira create' runs on a board with a usable remote, it shall send the ticket to its ref and delete the file once the remote accepted it. [RFC7GA]
SYNC-25  If a new ticket could not be sent, then 'jaira create' shall keep its file and say why it stayed. [RFC7GA, VM0A76]
SYNC-26  The jaira CLI shall show tickets that live only on their ref in list, next and the board, as read-only tickets without a file. [RFC7GA, PXYN99]
SYNC-27  If a write command names a ticket that lives only on its ref, then it shall refuse with exit 3 (on_ref_only) and point to 'jaira pull'. [RFC7GA]
SYNC-28  When 'jaira pull <id>' takes a ref-only ticket, it shall set the caller as assignee on the ref first and write the file only after that push. [RFC7GA]
SYNC-29  If the ref names another assignee, then 'jaira pull' shall refuse naming them, unless --steal, which notes the takeover on the ticket. [RFC7GA]
SYNC-30  If two clones pull the same ticket at once, then the loser shall get no file and a message naming the winner. [RFC7GA]
SYNC-31  When 'jaira pull' finds the file here behind its ref, it shall overwrite it with the ref and print both states, unless a write is unsent. [F3VNVC]
SYNC-32  When 'jaira pull' takes a ticket rejected on its ref, it shall lift the rejection and say so in a note on the ticket. [4539MZ]
SYNC-33  When 'jaira release <id>' runs, it shall clear the assignee on the ref first and then delete the file; another person's ticket needs --force. [RFC7GA]
SYNC-34  If 'jaira release' finds a file but no ref, then it shall clear the assignee, put the ticket on its ref and delete the file. [VM0A76]
SYNC-35  When 'jaira reject' targets a ref-only ticket, it shall fetch, stamp the rejection and clear the assignee on the ref under lease, writing no file. [4539MZ]
SYNC-36  While a ref carries a rejection, list and the board shall hide that ticket, and 'jaira show' shall still find it. [4539MZ]
SYNC-37  When 'jaira pull' or 'jaira reject' gets an id this clone does not know, it shall fetch once before resolving it. [4539MZ]

## Stale files
SYNC-38  If a ticket file is older than its ref, then every store write, move (also --force), archive and logbook shall refuse with exit 3 (stale_file), naming both states and 'jaira pull'. [F3VNVC]
SYNC-39  The jaira CLI shall judge "older than its ref" by updated-at, and on a tie by whether the file's bytes are a state the ref carried before. [F3VNVC]
SYNC-40  While a ticket file is behind its ref, list, show and the board shall show the ref's state marked behind-ref, leaving the file as it is. [K6BGB0]
SYNC-41  While a ticket has an unsent write here, the jaira CLI shall not mark it as behind its ref. [K6BGB0]

## Fetch
SYNC-42  When 'jaira fetch' runs, it shall fetch refs/jaira/* with --prune in one round trip and list the tickets, newly assigned to you first. [8566KF]
SYNC-43  The 'jaira fetch' command shall write no ticket file; milestone files are the only files it writes or removes. [RFC7GA, 1MYFKE]
SYNC-44  When a fetched ticket assigned to you changed since last seen, 'jaira fetch' shall raise one desktop notification, unless notify-off or --quiet. [8566KF]
SYNC-45  When any command finishes, the jaira CLI shall spawn a detached 'jaira fetch --json' when the last fetch is older than fetch-every (default 10 min). [PP5SCQ]
SYNC-46  The jaira CLI shall keep the background fetch stamp per repository, shared by every worktree of one clone. [NWFE2B]
SYNC-47  Where JAIRA_NO_FETCH is set, the jaira CLI shall start no background fetch. [PP5SCQ]
SYNC-48  While the board TUI runs, it shall fetch every fetch-every; while 'jaira web' runs, it shall fetch every minute. [17TJ1H, DVPAPE]

## Leaving the board
SYNC-49  When a ticket is filed with logbook or archive, the jaira CLI shall write its final state on the ref and leave the ref standing. [RA7PFE]
SYNC-50  When a filed ticket is in the logbook or archive of a landing branch on the remote, the snapshot run shall delete its ref. [RA7PFE]
SYNC-51  The jaira CLI shall take landing branches from settings, else from the remote's HEAD, always as revisions on the remote. [RA7PFE]
SYNC-52  While a ticket is filed away here, the board shall not show its still-standing ref as a card. [17TJ1H]
SYNC-53  If a filed ticket has not landed after landing-grace (default 3 days), then fetch shall list it as stranded and name 'jaira snapshot --drop'. [RA7PFE]
SYNC-54  If a ref seen before is gone while its file is still here, then fetch and validate shall report it as departed, move nothing, and name logbook, archive and pull. [RFC7GA]
SYNC-55  If the ticket refs cannot be listed, then fetch and validate shall warn on stderr and report no departures. [RGFJ0Q]

## Milestones
SYNC-56  The jaira CLI shall keep .jaira/milestones/ out of git through .gitignore and share milestone files only through their refs. [1MYFKE]
SYNC-57  When a milestone is about to be written, the jaira CLI shall first bring a file that fell behind its ref up to the ref, keeping hand edits the ref never carried. [1MYFKE]

## Sessions and claims
SYNC-58  When 'jaira claim' runs, it shall record the session on the ticket; a claim not renewed within 30 minutes counts as abandoned. [7WY0YT, K1W2JR]
SYNC-59  If another session holds a live claim, then 'jaira claim' shall refuse naming the holder, unless --steal is given. [7WY0YT]
SYNC-60  When 'jaira claim --release' runs, it shall clear the caller's own claim at once, and another session's only once it expired. [7WY0YT]

## Project registry
SYNC-61  When a board opens, the jaira CLI shall record it in ~/.jaira/projects.json ($JAIRA_HOME if set) under a lock, replacing the file whole. [BPR7C6]
SYNC-62  The project registry shall list a board in a linked worktree as its main checkout, once per repository, and skip the temp directory. [ZXJ72E]
SYNC-63  The project registry shall keep at most 20 boards, most recently opened first, and drop boards whose .jaira is gone. [ZXJ72E]
SYNC-64  If projects.json cannot be parsed, then the jaira CLI shall copy it to projects.json.broken (or .broken-<time>) and say so before rewriting it. [BPR7C6]
