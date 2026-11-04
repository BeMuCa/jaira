---
id: 01M6HZQPZZFAF3MXY3AWQFWRXH
title: "Eine PR-Lane zeigt gebuendelt, welche Tickets in welchem PR liegen, mit Name, Branch und Pipeline-Farbe"
status: todo
ready: true
creator: berkc
assignee: berkc
goal: "Das Board zeigt rechts nach done eine Nur-Lese-Spalte 'PRs' mit jedem offenen PR des Repos (auch der Kollegen) als Buendel in CI-Farbe; das Buendel klappt auf, solange der Cursor darauf steht, und nennt seine Tickets mit ihrer Lane; der Stand kommt beim Refresh aus gh und liegt nur im Speicher."
context: |-
  Wunsch berkc 04.11.2026 (Punkt 11): 'PR Lane klingt eig fun, gebuendelt darstellen, wobei es sich ausbreitet wenn man drueber hovered. Es sollte den PR-Namen oder Branch-Namen halten, und Farbe nach State der Pipeline. Vielleicht Cron-Job, der stuendlich checkt?'
  Heute: PRs entstehen pro Milestone (6TCTFV): jaira-role-pr baut milestone/<name> aus origin/master und merged die Ticket-Branches mit --no-ff. Den PR-Stand kennt jaira nicht.
  Randbedingung: kein Daemon, kein Server in jaira. Kandidaten fuer den Stand: 'gh pr view --json' beim Auffrischen des Boards (nur wenn gh da ist), oder eine GitHub Action, die den Stand auf den Ticket-Ref schreibt.
  Offen fuer brainstorm: Lane oder nur Gruppierung im Board; Quelle und Haeufigkeit des PR-Stands; GitLab.
definition-of-done: Brainstorm hat Darstellung und Quelle des PR-Stands festgelegt
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-11-04T00:35:54Z
updated-at: 2026-11-04T21:42:12Z
updated-by: berkc
claimed-by: berk-ca-1311236
claimed-at: 2026-11-04T21:27:24Z
question: ""
---

# Eine PR-Lane zeigt gebuendelt, welche Tickets in welchem PR liegen, mit Name, Branch und Pipeline-Farbe

## Definition of Done

- [ ] Brainstorm hat Darstellung und Quelle des PR-Stands festgelegt
- [ ] Tickets eines offenen PRs erscheinen gebuendelt mit PR-Name/Branch und Farbe nach CI-Stand; Test
- [ ] Release-Notiz unter core/release/notes/
- [ ] Spalte 'PRs' nach done mit allen offenen PRs des Repos (gh pr list), je PR: Nummer+Titel, Branch, CI-Zustand mit Glyph+Wort und Farbe; ohne gh oder ohne offene PRs sieht das Board aus wie heute; TUI-Test
- [ ] Cursor auf einem Buendel klappt es auf: erster fehlgeschlagener Check, je Ticket eine Zeile mit Lane; enter springt zur Karte, o oeffnet den PR im Browser; TUI-Test
- [ ] PR->Ticket ueber das Handle im Branchnamen und milestone/<name> -> Tickets des Milestones; Test
- [ ] Stand wird beim Start, beim Hintergrund-Fetch und beim manuellen Refresh geholt, nur im Speicher, nie auf Ticket oder Ref; Test
- [ ] Release-Note-Zeile in core/release/notes/QFWRXH.md

## Options

- [x] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-04 21:30 · berkc** — Brainstorm: what is true in the code today
- jaira knows nothing about PRs. No gh/glab call exists in Go code. Only the skill jaira-role-pr (core/role/builtin/jaira-role-pr/SKILL.md) uses gh/glab.
- There is no "milestone grouping" on the board. A milestone is one coloured cell on the right edge of a card (milestoneColors, model.go:1832) plus the M picker/filter (renderMilestones, view.go:1673).
- Card = 3-row band (title / handle+assignee / flags). Lane box is minColWidth 22, stretched to fill. An empty lane is 4 wide (thin).
- Refresh: the background ref fetch runs every 10 min by default (refFetchEvery, internal/tui/refs.go; DefaultFetchEvery in core/refsync/schedule.go).
- gh pr list --state open --json number,title,headRefName,isDraft,statusCheckRollup runs on this machine (returned []). glab is not installed, so nothing about GitLab is verified.
Mockups (ASCII, exact widths, 120 columns): /tmp/claude-1000/-home-berkc-code-X-General-jaira/15e0b874-a533-433c-887a-ca9272cf72bc/scratchpad/qfwrxh-mockups.md (scratch, may vanish)

Question: how should a PR's tickets show up on the board?
1. (recommended) Read-only "PRs" column on the right, after done. It is not a lane. One bundle per open PR in card format (#nr title / branch / CI glyph+word + ticket count), left cell in the CI colour. It expands while the cursor is on it: first failing check + one row per ticket with its lane, and the existing link lines go to the real cards. Trade-off: the whole PR in one place even across lanes; costs a non-lane column in boardFit/cursor.
2. Bundles inside the existing lanes. Trade-off: no extra column, but a milestone PR whose tickets sit in two lanes splits into two bundles, and it breaks the "one entry = one ticket" assumption in cardsInBudget/scroll/cardIdx.
3. Only a flag on the card (✓#131 / ✗#140 / ●#138) plus a P picker like M, with a filter. Trade-off: the smallest change, but there is no bundle on the board, and the flag row is already full.
Common to all: a grouping, not a lane. PR state changes outside jaira; a lane would need status writes on every CI run, which means ref churn and conflicts, and a ticket cannot be in signoff and in "PR" at once. Source: gh pr list when the board starts, on the existing ref fetch and on manual refresh. No daemon/cron/Action. Memory only, never written to tickets/refs. Mapping: handle segment in headRefName (as branchTicket, internal/cli/tickets.go:152) + milestone/<name> -> members (6TCTFV, not on master yet). CI: fail > running > pass > no CI. GitLab via glab, same shape, after GitHub (unverified).
- **2026-11-04 21:41 · berkc** — Entscheidung berkc 04.11.2026: Option 1, PRs-Spalte rechts. Zusatz wortgetreu: 'Man koennte dann quasi PR's von kollegen sehen. Ich denk das muesste damit es sinnvoll ist auch von allen die PRs anzeigen. also muessen wir auch aktive darstellen und pullen.' Also: alle offenen PRs des Repos (auch von Kollegen), nicht nur eigene; aktive PRs zeigen und den Stand beim Refresh holen.
