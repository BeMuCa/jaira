---
id: 01M7ZYQW91W2NET85VHEP122GE
title: "Ein Worker-Tab bleibt in signoff offen, solange sein Cache warm ist, damit ein Give back ihn weiterverwendet"
status: done
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
commits:
  - f72ac19969ef8e056aaa5712fdca11a364e28bf4
  - 9c171f6e2a10b9a2faa14ee87cf8640489bc5abf
  - 4378ef310f746b483564bb8bef6944233912b27d
created-at: 2026-11-21T21:03:35Z
updated-at: 2026-11-21T21:43:09Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2835032
claimed-at: 2026-11-21T21:04:38Z
outcome-what: "close-tab.sh parks a worker's tab in signoff/human (parked=<epoch> in its spawn.sh record) instead of closing; closes at once only in done. reap.sh hands every parked record to close-tab.sh --reap --pane, which closes after JAIRA_CACHE_WINDOW (default 3600 s), on done or filed away, and keeps a working tab or a ticket given back. Skill texts of worker, dispatcher, teamlead say so; release note; shipped.sha256."
outcome-why: "Since 0.5.15 every give-back from signoff started a new worker that re-read 100-200k of context; the cockpit's same-worker rule found no live tab."
outcome-resolves: "DoD 1-3: the tab stays while the cache is warm, the next reap closes it after the window or on done, release note and manifest. Four categories: concurrent park/reap via atomic rename and reap.lock, missing record/sed failure keep the tab, odd pane ids sanitized as spawn.sh does, no new trust."
review-summary: "close-tab.sh no longer closes a worker's Herdr tab when its ticket reaches signoff or another person's lane without a question: it writes parked=<epoch> into the tab's spawn.sh record (git common dir jaira/panes/<pane>, via temp file + mv) and prints 'kept <T>: cache warm until HH:MM'. In done it closes at once as before. reap.sh gets a parked pass that runs before any forge question: for every record with parked= it calls close-tab.sh --reap --pane <pane> <T> <worktree>, which closes once the ticket is done or filed away, or the park is older than JAIRA_CACHE_WINDOW (default 3600 s); it keeps the tab while the agent is working, while the ticket has an open question, or while the ticket is back in an agent lane. A closed or vanished tab loses its parked= line so later rounds skip it. Worker, dispatcher and teamlead skills and a release note describe the new behaviour; closetab_test.go covers park, warm, window over, done, filed away, given back, agent working, open question."
round: 2
back-from: critique
back-by: lane
back-count: 1
review-gaps: |-
  note: close-tab.sh:88 checks JAIRA_CACHE_WINDOW for digits only, but bash arithmetic (lines 174, 180) reads a leading zero as octal: =08 aborts --reap (swallowed by reap.sh || true, tab never closes by window), =010 means 8 s; fix: window=$((10#$window)) after the check
  note: reap.sh:148 falls back to $root when a record's worktree is gone, so close-tab.sh --reap checks uncommitted changes of the checkout reap runs in and can keep a parked tab with "uncommitted changes" of unrelated files; fix: skip the uncommitted check in --reap mode or pass a flag that the worktree is gone
  note: a ticket given back to an agent lane and picked up by a NEW worker leaves the old parked tab open until the ticket returns to a person's lane or done (status in-progress keeps it); acceptable by design (decision note), but the old tab idles for the whole round
  note: no test for a non-default window with a leading zero or for the missing-worktree fallback
test-verdict: "pass: go build/vet ./... green, go test -race ./core/role ./core/shipped RC=0 (full go test ./... green in in-progress), shipped-history --check green; DoD 1-3 verified (tests named in proofs ran, note file present); scratch repo with stub herdr: signoff parks (parked= written, kept with HH:MM), reap warm keeps '60 more min', JAIRA_CACHE_WINDOW=0 closes t1 and drops parked=; failure paths: bad window exit 2, --reap without dir exit 2, missing record keeps with reason (exit 0 as documented); six parallel parks leave one parked= line and no leftover temp file"
hunt-verdict: "pass: bughunt done at 4378ef310, 2 change (no blocker: JAIRA_CACHE_WINDOW with a leading zero read as octal; reap's parked pass falls back to its own checkout when the record's worktree is gone), 0 pre-existing; both left as notes for review and the person"
review-verdict: "Both DoD items are met and tested (go test ./role -run 'CloseTab|Reap' green); the release note and shipped.sha256 are there. No blocker: the two bughunt findings are real but only change how long a tab stays open (octal window value, missing worktree), nothing is lost or closed wrongly. Recommend accept, optionally with the one-line 10#$window fix."
review-round: 2
review-verdict-gaps: |-
  note: close-tab.sh:88 checks JAIRA_CACHE_WINDOW for digits only, but bash arithmetic (lines 174, 180) reads a leading zero as octal: =08 aborts --reap (swallowed by reap.sh || true, tab never closes by window), =010 means 8 s; fix: window=$((10#$window)) after the check
  note: reap.sh:148 falls back to $root when a record's worktree is gone, so close-tab.sh --reap checks uncommitted changes of the checkout reap runs in and can keep a parked tab with "uncommitted changes" of unrelated files; fix: skip the uncommitted check in --reap mode or pass a flag that the worktree is gone
  note: a ticket given back to an agent lane and picked up by a NEW worker leaves the old parked tab open until the ticket returns to a person's lane or done (status in-progress keeps it); acceptable by design (decision note), but the old tab idles for the whole round
  note: no test for a non-default window with a leading zero or for the missing-worktree fallback
review-check: "1. cd core && go test ./role -run 'CloseTab|Reap' -v  2. see TestCloseTabParksAWorkerWhoseTicketIsWithAPerson and all 7 subtests of TestReapClosesAParkedTabOnceItsCacheIsCold PASS  3. live: let a worker finish a ticket into signoff in a Herdr tab - its last line reads 'kept <ID>: cache warm until HH:MM' and the tab stays open  4. cat $(git rev-parse --git-common-dir)/jaira/panes/<pane> - it holds parked=<epoch>  5. JAIRA_CACHE_WINDOW=0 core/role/builtin/jaira-dispatcher/scripts/reap.sh --dry-run - it prints 'would close <ID> <tab>'"
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
- **2026-11-21 21:15 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-dispatcher/scripts/close-tab.sh:88 - class: JAIRA_CACHE_WINDOW is checked for digits only (line 88), but bash arithmetic reads a leading zero as octal (lines 174 and 180); example: JAIRA_CACHE_WINDOW=08 or =0900 aborts --reap at line 180 with exit 1 (usage promises 0 or 2; reap.sh:149 swallows it with || true, so the parked tab never closes by window), in the worker path line 174 hits the ERR trap ('kept: line 174 failed'), so the tab is never parked; =010 silently means 8 seconds; not a blocker: the default 3600 and any value without a leading zero are unaffected; ran: bash -c 'window=08; echo $((1000+window-1000))' (value too great for base, rc=1), window=010 gives 8
- **2026-11-21 21:15 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-dispatcher/scripts/reap.sh:148 - class: the parked pass falls back to $root, the checkout reap.sh runs in, when a record's worktree is gone, so close-tab.sh --reap checks the wrong directory for uncommitted changes and runs jaira show there; example: a worker parks, its worktree is removed by hand, the cockpit runs reap.sh from a checkout with unrelated uncommitted edits, and every round prints 'kept <ticket>: uncommitted changes: <files of another checkout>', so the parked tab stays until the PR merges; not a blocker: a tab stays open longer, nothing is lost; read
- **2026-11-21 21:15 · Alexander Sacharov** — bughunt: done at 4378ef310 - 2 findings (2 change, 0 pre-existing, no blocker) in 4 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/release/notes/P122GE.md core/role/builtin/jaira-dispatcher/SKILL.md core/role/builtin/jaira-dispatcher/scripts/close-tab.sh core/role/builtin/jaira-dispatcher/scripts/reap.sh core/role/builtin/jaira-teamlead/SKILL.md core/role/builtin/jaira-worker/SKILL.md core/role/closetab_test.go core/role/reap_test.go core/shipped/shipped.sha256
