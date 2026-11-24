# One jaira name per person, kept in a people list on a shared ref
Status: Accepted
Date: 2026-11-23
## Context
One human appeared on a board under every git user.name and email they had; nothing said "BeMuCa" and "berkc" were one person.
The ownership gate compared raw strings and refused people their own tickets, which trained --force on every move.
## Decision
The list lives in the file `people` on refs/jaira/people, one line `<jaira name>: <alias>, ...`, hand-editable and diff-readable.
jaira resolves git user.name/email through it and writes the jaira name; tickets are never rewritten, readers canonicalise.
Writes are synchronous: fetch, change the newest list, push with a lease, redo on a lost race (max 4); offline writes nothing.
## Rejected
- Ask the name once per project on first start and keep it in a local, gitignored setting (berkc's approach): each clone could answer differently.
- Rewrite every ticket file on rename: a migration per rename instead of one write.
- Queue people writes in the outbox like tickets: a second store holding unsent names.
## Consequences
Rename is one write; `list --assignee`, gates, TUI and web show and match any name of a person.
People commands need the board's remote; `jaira snapshot` backs the list up.
## Sources
DVPAPE (decision notes 2026-11-23 19:41, 19:44, 19:53)
