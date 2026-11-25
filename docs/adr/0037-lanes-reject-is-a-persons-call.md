# Rejecting a ticket is a person's decision and records why
Status: Accepted
Date: 2026-10-19

## Context
An archived ticket looked the same whether it was finished elsewhere, abandoned or refused; months later nobody knew why.

## Decision
`jaira reject <id> --why` writes rejected-kind/why/by/at and archives the ticket; without --force it refuses (exit 3)
and tells an agent to move the ticket to the question lane with "reject? reason: ...". A ref-only ticket is rejected on its ref.

## Rejected
- Plain archive for refused work: the archive could not say why.
- Letting agents reject: whether work is done is not the worker's call.
- Requiring a pull before rejecting a ref-only ticket (4539MZ).

## Consequences
`jaira restore` takes the rejection off and keeps a related link that existed before; `jaira claim` refuses a rejected ticket.

## Sources
A0C081 (decision notes 2026-10-19, 2026-10-21), 4539MZ (2026-11-22), X8KZQ7
