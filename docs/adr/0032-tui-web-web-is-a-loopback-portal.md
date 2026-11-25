# jaira web is a loopback page inside the binary, synced only through git
Status: Accepted
Date: 2026-10-14

## Context
People without an agent or terminal board need to sort tickets into epics and read the work.
The project rules forbid a server and any runtime dependency.

## Decision
`jaira web` serves an embedded Svelte build on 127.0.0.1 for as long as the command runs.
It reads and writes the same refs as the CLI: fetch every minute, push after each change.
Requests must name a loopback host; writes with a foreign Origin are refused.

## Rejected
- A hosted or long-running server.
- Building the page at install time (would need Node); the dist is committed and checked by CI.

## Consequences
Errors, failed syncs and panics go to a log file under the state dir (HEY6T7).
The page is read-only for ticket content; it edits epics only.

## Sources
RFTQKP, HEY6T7, TPCKZB
