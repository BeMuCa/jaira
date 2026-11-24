# jaira web
_Owner: Berk · EARS · each line one requirement · [tickets] hold the detail_

## Serving
WEB-01  When `jaira web` runs, it shall listen on 127.0.0.1 only (`--port`, 0 = free port), print the URL and open the browser unless `--no-open`. [RFTQKP]
WEB-02  If a request names a non-loopback Host, then the server shall answer 403; a write with a foreign Origin shall also get 403. [RFTQKP]
WEB-03  The page shall be a Svelte build embedded in the binary, so running it needs no Node and no server beyond the command. [RFTQKP]
WEB-04  While the command runs, it shall fetch every minute and push after every change; a failed push shall stay queued for the next one. [RFTQKP]
WEB-05  The page shall show when it last synchronised and any push error. [RFTQKP]
WEB-06  The command shall log every answer >=400, failed fetch or push, panic with stack and the ending signal to `<state>/web.log`. [HEY6T7]
WEB-07  When the command starts, it shall print the log's path; if the log cannot be opened, it shall warn and serve on. [HEY6T7]
WEB-08  If a handler panics, then the server shall answer 500 instead of dropping the connection. [HEY6T7]
WEB-09  If the command stops answering, then the page shall show a banner within 10 seconds. [HEY6T7]
WEB-10  The server shall answer not-found with 404, a taken name or a change that does not apply with 409, and invalid input with 400.

## Epics
WEB-20  The page shall show open tickets in no epic on the left and epics with description and done X of Y on the right. [RFTQKP]
WEB-21  When a ticket is dragged onto an epic, the page shall write the epic's ref and push it; dragging back shall take it out. [RFTQKP]
WEB-22  Each direct ticket card in an epic shall have a Remove button that works by tap, without dragging. [3XXF5B]
WEB-23  A ticket counted only through a member it follows shall name that member and offer no drag or Remove. [3XXF5B]
WEB-24  When "+ New epic" or "Edit" is used, the page shall create an epic or rewrite title, description and text as the CLI does. [CRDH6B]
WEB-25  If an epic edit is based on a stale version, then the server shall answer 409 and the form shall keep the typed text. [CRDH6B]
WEB-26  The page shall hide closed epics unless "Show closed" is on, and offer Close and Reopen on an epic. [182JT1]
WEB-27  When a ticket is dropped on a closed epic, the page shall ask for confirmation before adding it. [182JT1]
WEB-28  Where an epic has `external:` keys, the page shall show each key as a link and filter by key. [PG9P9Y]
WEB-29  Where an epic is named, the page shall show its title with its name small beside it. [PG9P9Y]

## Tickets and filters
WEB-40  When a ticket is clicked, the page shall show goal, context, DoD and every note as markdown, read-only. [RFTQKP]
WEB-41  Each card shall show its assignee cut to 20 characters with the full name as tooltip, and nothing without one. [ZDBXN1]
WEB-42  Each card shall show its created date and `untouched <age>` since its last write. [ZKPT4H]
WEB-43  If a ticket file is older than its ref, then the page shall show the ticket in the ref's lane. [ZKPT4H]
WEB-44  The filter bar shall filter by text over title, goal and id, by lane, tag, person, where and a date range. [TPCKZB, ZDBXN1]
WEB-45  The lane, tag, person and where filters shall take several chips, any of which matches; a click shall turn a chip into `not`. [ZDBXN1]
WEB-46  The page shall keep every filter in the URL query, so a link shows the same list; `?all=1` shall read as every place. [TPCKZB, DF70XX]
WEB-47  The where filter shall pick board (ref-only included), logbook, archive and rejected; filed tickets open but never drag. [DF70XX]
WEB-48  A filed ticket shall show how it left: done struck through, abandoned dimmed with its lane, rejected with red edge and reason. [DF70XX]
WEB-49  Every lane shall have one colour on badges and in the lane filter: human lanes amber, terminal green, parking red. [TPCKZB]
