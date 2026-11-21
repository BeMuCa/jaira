---
id: 01M7ZYQW91W2NET85VHEP122GE
title: "Ein Worker-Tab bleibt in signoff offen, solange sein Cache warm ist, damit ein Give back ihn weiterverwendet"
status: testing
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Give back aus signoff innerhalb des Cache-Fensters landet beim selben Worker statt bei einem neuen, der den ganzen Kontext neu liest."
context: |-
  Was heute teuer war (21.11.): seit 0.5.15 schliesst close-tab.sh den Worker-Tab sofort, wenn das Ticket signoff erreicht. 9 Give-backs an diesem Tag = 9 neue Worker, die je 100-200k Kontext neu gelesen haben; SendMessage an den alten Socket schlug mit ENOENT fehl.
  Das Cockpit-Tool workers und die Give-back-Regel kennen schon "same worker (cache warm, context under limit)" — sie finden nur keinen lebenden Tab mehr.
  Gewünscht: core/role/builtin/jaira-dispatcher/scripts/close-tab.sh schliesst in signoff erst, wenn das Cache-Fenster (60 min nach der letzten Antwort) abgelaufen ist oder das Ticket done ist; ein person-Lane mit offener Frage hält den Tab wie bisher. reap.sh bzw. der Cockpit-Reap (alle 10 min) schliesst ihn danach.
  Alex, 21.11.: Token sparen ohne Qualität zu verlieren.
definition-of-done: "close-tab.sh lässt den Tab eines Tickets in signoff offen, solange die letzte Antwort jünger als das Cache-Fenster ist (Test)"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-11-21T21:03:35Z
updated-at: 2026-11-21T21:12:44Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2835032
claimed-at: 2026-11-21T21:04:38Z
outcome-what: "close-tab.sh parks a worker's tab in signoff/human (parked=<epoch> in its spawn.sh record) instead of closing; closes at once only in done. reap.sh hands every parked record to close-tab.sh --reap --pane, which closes after JAIRA_CACHE_WINDOW (default 3600 s), on done or filed away, and keeps a working tab or a ticket given back. Skill texts of worker, dispatcher, teamlead say so; release note; shipped.sha256."
outcome-why: "Since 0.5.15 every give-back from signoff started a new worker that re-read 100-200k of context; the cockpit's same-worker rule found no live tab."
outcome-resolves: "DoD 1-3: the tab stays while the cache is warm, the next reap closes it after the window or on done, release note and manifest. Four categories: concurrent park/reap via atomic rename and reap.lock, missing record/sed failure keep the tab, odd pane ids sanitized as spawn.sh does, no new trust."
review-summary: "none"
round: 2
back-from: critique
back-by: lane
back-count: 1
review-gaps: "reap.sh: the parked pass reads its record through the existing field() (moved up beside herdr_up) instead of three inline seds; close-tab.sh: terminal is set in the same lane-flag loop instead of a second jaira json get. Left: the pane-to-record name (tr -c 'A-Za-z0-9_.-' '_') is written in spawn.sh and repeated in close-tab.sh - one line, two scripts, no shared lib between them today."
---

# Ein Worker-Tab bleibt in signoff offen, solange sein Cache warm ist, damit ein Give back ihn weiterverwendet

## Definition of Done

- [x] close-tab.sh lässt den Tab eines Tickets in signoff offen, solange die letzte Antwort jünger als das Cache-Fenster ist (Test)
  proof: TestCloseTabParksAWorkerWhoseTicketIsWithAPerson, TestReapClosesAParkedTabOnceItsCacheIsCold/cache_warm (core/role/closetab_test.go); close-tab.sh parks via park()
- [x] Nach Ablauf des Fensters oder bei done schliesst der nächste Reap den Tab (Test)
  proof: TestReapClosesAParkedTabOnceItsCacheIsCold/{window_over,ticket_done,ticket_filed_away}; reap.sh parked pass before the forge
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/P122GE.md; core/shipped/shipped.sha256 regenerated (shipped-history.sh --check green)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-21 21:04 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock
- **2026-11-21 21:06 · Alexander Sacharov** — decision: no open decisions for the person. Mechanism: close-tab.sh, in a person's lane that is not terminal (signoff, human without question), does not close but writes parked=<epoch> into the tab's spawn.sh record (git common dir jaira/panes/<pane>) and prints 'kept <T>: cache warm until HH:MM'. reap.sh, before any forge question, calls close-tab.sh --reap --pane <pane> <T> <worktree> for every record with parked=; that mode skips the parent check, closes when the ticket is done/filed away or parked is older than JAIRA_CACHE_WINDOW (seconds, default 3600 = the subscription TTL in cockpit/hooks/cache.ts), keeps when the tab's agent works or the ticket went back to an agent's lane (the same worker carries it again and writes parked anew when it finishes). One place holds the closing rules; reap only finds the parked tabs. Why not the cockpit's session records with lastAnswerAt: they live in $.store, which reap.sh (bash) cannot read. Known limit: a person typing into the parked tab refreshes the cache but not parked=, so reap may close up to one window early.
- **2026-11-21 21:11 · Alexander Sacharov** — critique: teamlead SKILL.md sentence about reap lost its verb; fix wording. Shape checked and left: parked= in the spawn.sh record + close-tab.sh --reap (see decision note); a double close of a parked tab by reap's merged-ticket loop in the same round is harmless (herdr closes once, the second finds no pane).
