---
slug: dod-verb
created: 2026-08-12
status: in-progress
---

# Writing checklist state

Agents need to record which step they are on and which criteria they have met.
Reading the checklists landed in `dod-checkbox-states`; this task adds the write
path.

## Why a core function rather than a sed in the CLI

Checkbox position is load-bearing. An item only counts inside a recognised
heading, so a naive append lands it under whatever section happens to be last —
that mistake was made and caught while verifying the previous task. The write
therefore has to find the section, count items within it, and replace exactly one
marker character in place.

## What changes

1. `Section` (`SectionDoD` / `SectionPlan`) selecting which checklist to address.
2. `SetItemState(body string, sec Section, i int, st State) (string, error)` —
   replaces the marker of the i-th item (0-based) of that section and returns the
   new body. Every other byte, including line endings and the item's own text, is
   untouched. Out-of-range indexes are an error, never a silent no-op.
3. Setting an item to `doing` clears any other `doing` item **in the same
   section**, so "which step is the agent on" has exactly one answer.
4. `jaira dod <id> <n> [--plan] --doing|--done|--todo`, 1-based to match what the
   TUI and `jaira show` display.

## Out of scope

- TUI rendering (next task)
- Archive, review-verdict, --signal removal (separate tasks)

## Verification

- Byte-fidelity test: exactly one character differs between input and output.
- CRLF bodies keep their line endings.
- Out-of-range index returns an error and leaves the body unchanged.
- Setting a second item to doing moves the marker rather than adding one.
- End-to-end: `jaira dod` on a real ticket, then `jaira show`, then a move to
  `done` that is refused while an item is in progress.
