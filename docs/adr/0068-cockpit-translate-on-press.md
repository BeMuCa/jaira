# The Acceptance pane translates only on a press, in one call, and never writes the translation
Status: Accepted
Date: 2026-11-23

## Context
0GMYQB rephrased notes into the person's settings.json language automatically on every opening.
That cost two haiku calls per ticket and the text jumped from original to translation while being read.

## Decision
A Translate button (key t) translates head and all notes in one haiku call; a second press returns the original.
Every opening starts in the original; answers are cached in the store per language and text.

## Rejected
- Automatic rephrasing on opening (0GMYQB) - cost and text jumping under the reader.

## Consequences
What Execute writes onto the ticket is always the original wording.

## Sources
FMWBA8, 0GMYQB
