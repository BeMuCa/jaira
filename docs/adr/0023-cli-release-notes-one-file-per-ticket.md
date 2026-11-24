# Release notes are one file per ticket, and the folder names the release
Status: Accepted
Date: 2026-10-12

## Context
Every PR edited the Unreleased section of core/release/NOTES.md, so two PRs always conflicted there.
A release also needed its own commit to rename Unreleased into the version number.

## Decision
Each client-facing change adds core/release/notes/<TICKET>.md with one '- ' line per change ('- ! ' = must act).
A file directly in notes/ is Unreleased; the release cut moves it into notes/<X.Y.Z>/, and the folder name is the version.
Notes are embedded with go:embed; NOTES.md is closed history up to 0.3.6.

## Rejected
- One NOTES.md with an Unreleased heading (conflicts on every PR pair).
- A changelog generated from git at release time (written when nobody remembers why it matters).
- A version heading inside each note file (the file would have to change at release).

## Consequences
Two PRs never touch the same note lines; .gitattributes must pin notes/**/*.md to LF (5HDAJZ).
Moving or editing a note inside a version folder by hand rewrites shipped history.

## Sources
AR97K9, 5HDAJZ
