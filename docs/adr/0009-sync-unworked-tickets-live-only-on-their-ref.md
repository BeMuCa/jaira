# An unworked ticket lives only on its ref until someone pulls it
Status: Accepted
Date: 2026-09-11

## Context
A ticket logged in one clone while another still had it under tickets/ existed twice after the merge (rename/modify),
and add/add conflicts appeared for the same ticket created on two branches.

## Decision
On a board with a usable remote, 'jaira create' sends the ticket to its ref and deletes the file once the remote took it.
'jaira pull' claims the ticket on the ref first (assignee + CAS) and only then writes the file, so exactly one clone holds it.
'jaira release' clears the assignee on the ref first, then deletes the file. A board without a remote keeps files (one place decides: fileOnRefOnly).

## Rejected
- Keeping a file of every ticket in every checkout: duplicates through merges, hides who owns the ticket.
- Writing the file before the ref push in pull: the race loser would hold a ticket somebody else works.

## Consequences
Write commands refuse a ref-only ticket (exit 3, on_ref_only) and point to 'jaira pull'.
The untouched backlog lives only on the remote; the snapshot branch keeps a copy. Two storage modes exist and must stay named.

## Sources
RFC7GA (context: "Drei Preise"), VM0A76, internal/cli/refs.go fileOnRefOnly
