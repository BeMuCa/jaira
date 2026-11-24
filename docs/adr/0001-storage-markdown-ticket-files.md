# A ticket is one markdown file with flat frontmatter, rewritten by splicing
Status: Accepted
Date: 2026-08-11
## Context
The file format is the API: tickets are hand-edited and read in git diffs.
Many sessions and teammates create tickets at once with no server to hand out numbers.
## Decision
One file per ticket, .jaira/tickets/<ULID>-<slug>.md; id from frontmatter, handle = last 6 chars.
Writes splice one value into the original bytes; unknown keys, order and comments survive.
Frontmatter stays flat (outcome-what/-why/-resolves, rejected-*) so every field merges as one key.
## Rejected
- Sequential ids (JAIRA-1): a counter file is the worst merge-conflict hotspot.
- Unmarshal/marshal or re-printing the YAML AST: renormalises fields nobody touched.
- Nested outcome mapping: one change rewrites a block, merge driver works per top-level key.
- Renaming the file on title change: churns history; the slug is only a breadcrumb.
## Consequences
YAML anchors/aliases are refused (unsafe_yaml) rather than risked.
A ULID head is a timestamp, so references resolve by prefix or suffix.
## Sources
core/ticket/id.go, core/ticket/frontmatter.go (initial store commit 2026-08-11); KA9CFA, BWT53T
