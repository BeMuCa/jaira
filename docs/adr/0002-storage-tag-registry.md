# Tags are a closed, normalised vocabulary with colours in one .jaira/tags file
Status: Accepted
Date: 2026-09-04
## Context
"ui", "frontend" and "gui" on one board are three names for one subject and filter to nothing.
## Decision
A ticket carries a tags list; colours live per tag in .jaira/tags, one "name: <ansi256>" line.
Names are stored lowercase kebab; other characters are refused, never trimmed.
jaira only appends lines for new names and rewrites a recoloured line; the file merges by union.
## Rejected
- Colour on each ticket: one tag could render two ways; recolouring rewrites every ticket.
- YAML or a directory for the registry: a line is the smallest thing git merges.
- Silently trimming bad characters: invents a second name for one subject.
## Consequences
Agents must read jaira tags before tagging; past sixteen tags colours repeat.
## Sources
79GEPW (2026-09-04), AXTFG3, S1VM40
