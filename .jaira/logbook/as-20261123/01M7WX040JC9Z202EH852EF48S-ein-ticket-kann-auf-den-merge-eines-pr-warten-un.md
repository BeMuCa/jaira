---
id: 01M7WX040JC9Z202EH852EF48S
title: "Ein Ticket kann auf den Merge eines PR warten, und nach dem Merge startet sein Worker"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "blocked_by nimmt pr#<n>; jaira zeigt das Ticket als wartend, und sobald der PR gemergt ist, ist es actionable und der Teamlead startet den Worker ohne Notiz im Kopf."
context: |-
  Heute: Abhaengigkeiten wie "startet nach Merge von PR #12" stehen nur in Notizen des Teamleads.
  blocked_by kennt nur Ticket-IDs (jaira create --blocked-by).
  Alex 20.11.2026: "klingt gut". Kein Daemon: der Merge wird beim naechsten Lauf erkannt (next, list, reap, Cockpit), gh/glab als Quelle.
definition-of-done: "jaira set <id> blocked_by=pr#12 wird angenommen; list/next zeigen es als wartend auf PR #12 (Test)"
tags: []
blocked-by:
  - 01M7WWYG213BWB06A7D1CRASZS
related:
  - 01M7WWYG213BWB06A7D1CRASZS
  - 01M263RFEQMYPPT67D40Y5SW81
commits:
  - cca6f9e01ad1c55a51a83765562fe7ac6a3e3bad
  - b38b1bb09000ff606c43e08598918401610287d6
  - cee1269c9dc562da3a2e1ff1929eb406235a6bc6
  - db4be8be7a71c6f7cd9f09c5d3193018ec2f8dfe
  - 3b15b18c42add621e1b2fb2bc1e23c47a8298a4a
created-at: 2026-11-20T16:35:24Z
updated-at: 2026-11-22T23:23:35Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-81804
claimed-at: 2026-11-22T19:38:55Z
outcome-what: "blocked-by takes pr#<n>; list/next/move ask gh or glab for the pull request's state (core/pr State + Checker, cached in ~/.jaira/repo/<key>/pr-state.json: merged for good, open/closed 2 min); list shows [waits for PR #12: open], list --json waiting_on, next --json waiting_on_pr; TUI reads the cache only; teamlead skill starts a worker for tickets a merge freed"
outcome-why: "a dependency on a pull request lived only in the teamlead's notes, so a merged PR freed nothing and the follow-up waited until somebody remembered"
outcome-resolves: "2EF48S: a ticket can wait for a PR merge and becomes actionable on the next run after it"
review-summary: "blocked-by now takes pr#<n> beside ticket ids. core/ticket/prref.go parses it (any case, positive decimal up to 2^31-1, stored as pr#<n>); set and create refuse a malformed one with exit 3, and create now resolves ticket handles to full ids like set. core/pr/state.go asks gh (graphql state) or glab (MR state) for just the state, sharing connect() with Read; core/pr/cache.go remembers answers in ~/.jaira/repo/<repo>/pr-state.json (merged for good, open/closed 2 minutes, failures never cached, atomic write, corrupt/unwritable file only costs a call). gate.blockedBy treats pr#<n> through an injected PR hook: only merged clears it, every other answer (open, closed, no tool, signed out, nil hook) keeps the ticket waiting with code waits_on_pr (exit blocked on move). list shows [blocked] [waits for PR #12: open], list --json waiting_on, next prints 'Waiting on a pull request:' and next --json waiting_on_pr. validate, link index, show, TUI link lines and the TUI b key skip pr refs as tickets; the TUI answers from the cache only, the checker reset on reload and on a store root change. The teamlead skill starts a worker for every ticket a merge freed."
round: 3
back-from: testing
back-by: lane
back-count: 1
review-gaps: |-
  note: githubState calls githubTarget (gh repo view) for every PR number, so each uncached PR costs two gh calls instead of one plus one per process - memo the target on the Checker (testing measured 2.3s for 2 PRs on a cold cache)
  note: an open PR is re-asked every 2 minutes, so on a board with several open pr# blockers every list/next after the window pays ~1s per open PR - acceptable by decision (3), but it is the one place jaira list is no longer instant
  note: docs/COMMANDS.md says next lists tickets that wait "only" for a pull request, while waitsOnPR also lists a ticket blocked by a PR and a ticket together - drop "only"
  note: the TUI move gate refuses a ticket whose PR merged but whose cache holds no fresh answer ("not asked here") until list/next runs - intended by decision (6), worth one line in the TUI notice
test-verdict: "pass: round 2 (diff since cee1269c9: internal/tui only): go test -race ./internal/tui RC=0, TestBoardSwitchReadsTheNewRepositorysPRCache covers the class (store change); round 1: go test ./... RC=0, DoD 1-8 verified, real-gh function pass, failure path and two writers checked"
hunt-verdict: "pass: bughunt done at db4be8be7, 0 change, 0 pre-existing (round 1 blocker at cee1269c9 fixed in db4be8be7)"
review-verdict: "accept: all 8 DoD items are met by code and tests; go build, vet and go test ./core/... ./internal/... pass; no blocker found along concurrency, error paths, input and permissions"
review-round: 3
review-verdict-gaps: |-
  note: githubState calls githubTarget (gh repo view) for every PR number, so each uncached PR costs two gh calls instead of one plus one per process - memo the target on the Checker (testing measured 2.3s for 2 PRs on a cold cache)
  note: an open PR is re-asked every 2 minutes, so on a board with several open pr# blockers every list/next after the window pays ~1s per open PR - acceptable by decision (3), but it is the one place jaira list is no longer instant
  note: docs/COMMANDS.md says next lists tickets that wait "only" for a pull request, while waitsOnPR also lists a ticket blocked by a PR and a ticket together - drop "only"
  note: the TUI move gate refuses a ticket whose PR merged but whose cache holds no fresh answer ("not asked here") until list/next runs - intended by decision (6), worth one line in the TUI notice
review-check: |-
  1. go build -o /tmp/jaira ./cmd/jaira
  2. in a clone whose origin is on GitHub with gh signed in: /tmp/jaira create "probe" --goal g --context c --dod d --blocked-by pr#<an open PR number>
  3. /tmp/jaira list - the row shows [blocked] [waits for PR #<n>: open]
  4. /tmp/jaira next --json | jq .waiting_on_pr - the ticket is there with waiting_on "waits for PR #<n>: open"
  5. /tmp/jaira set <id> blocked-by=pr#<a merged PR number>, then /tmp/jaira next --all --json - the ticket is in tickets
  6. /tmp/jaira set <id> blocked-by=pr#abc - exit 3, nothing written
  7. PATH=/usr/bin:/bin /tmp/jaira list (gh not on PATH, cache older than 2 min for an open PR) - row says state unknown (forge tool missing ...)
  8. jaira reject or delete the probe ticket afterwards
---

# Ein Ticket kann auf den Merge eines PR warten, und nach dem Merge startet sein Worker

## Definition of Done

- [x] jaira set <id> blocked_by=pr#12 wird angenommen; list/next zeigen es als wartend auf PR #12 (Test)
  proof: internal/cli/prblock_test.go TestBlockedByTakesAPullRequest (set blocked-by=pr#12 stored pr#12), TestTicketWaitsForItsPullRequest (list '[waits for PR #12: open]', next --json waiting_on_pr, next text 'Waiting on a pull request:'); field is blocked-by
- [x] Ist PR #12 gemergt, ist das Ticket in jaira next --json actionable
  proof: internal/cli/prblock_test.go TestTicketWaitsForItsPullRequest: pr#13 MERGED -> handed out by next --all --json; core/gate/prblock_test.go TestPRBlocker merged case
- [x] Der Teamlead-Skill startet beim naechsten Durchgang einen Worker fuer so frei gewordene Tickets
  proof: core/role/builtin/jaira-teamlead/SKILL.md 'Then start what a merge freed': every ticket that left waiting_on_pr for tickets gets its worker in this pass, blocker written as pr#<n>
- [x] pr#<n> input classes: PR#12/pr#12 accepted and stored as pr#12; pr#, pr#0, pr#-1, pr#abc, pr# 12, pr#99999999999999999999 refused by set and create with exit 3 (as an unknown ticket is); a hand-edited malformed pr ref is a validate error, not a crash (table test)
  proof: core/ticket/prref_test.go TestParsePRRef (table incl. pr#, pr#0, pr#-1, pr#abc, pr# 12, overflow, full-width digits); internal/cli/prblock_test.go TestBlockedByTakesAPullRequest (set+create exit 3, nothing written); core/validate/prref_test.go TestPRRefInBlockedBy (bad_pr_ref)
- [x] forge failure classes keep the ticket waiting with the reason in list/next: gh/glab missing, not signed in, no forge origin, PR closed without merge, unreadable answer (table test with a fake runner)
  proof: core/pr/state_test.go TestState (gh missing, signed out, no origin, foreign host, no such PR, graphql error, unreadable, unknown state, closed); internal/cli/prblock_test.go TestTicketWaitsForItsPullRequest (closed, 401) + TestNoForgeToolKeepsTheTicketWaiting
- [x] cache: merged is remembered across processes and worktrees, open/closed re-asked after 2 minutes; written atomically (temp+rename), an unwritable or corrupt cache file only costs a forge call, never an error (test)
  proof: core/pr/cache_test.go TestCheckerCache (one call per process, shared file, merged kept, open re-asked after OpenFor, failures not cached, no temp left) + TestCheckerCacheFailures (corrupt, unwritable, offline); core/pr/cache.go store writes through ticket.WriteAtomic
- [x] every reader of blocked-by treats pr#<n> as a PR, not a ticket: validate (no dangling-dep error), link index, sync export, show, TUI blocker lines; grep of BlockedBy users is the checklist (test where code)
  proof: validate: core/validate/prref_test.go; link: core/link/prref_test.go TestPRRefIsNoRelation; show/TUI display: ticket.Handle keeps pr#n (core/ticket/prref_test.go TestHandleKeepsPRRef); TUI: internal/tui/prblock_test.go (cache-only gate, b key), linklines skips pr refs; sync keeps existing entries (no change needed); create resolves like set
- [x] acceptance rule: starting work before the PR merged is worse than waiting a round longer
  proof: core/gate/gate.go blockedBy: every non-merged answer and nil PR hook keep CodeWaitsOnPR; core/gate/prblock_test.go TestPRBlocker 'no hook'/'unknown'/'malformed'

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] pr ref parse/format in core/ticket + set/create accept it
- [x] pr.State light query in core/pr (gh graphql state, glab mr state) + tests
- [x] cache in git common dir + CLI env hook gate.Env.PR
- [x] gate.blockedBy: pr refs via hook; list/next text+json show waiting on PR
- [x] siblings: validate, link, sync, show, tui
- [x] teamlead skill: pass picks up tickets freed by a merge
- [x] release note + docs/COMMANDS + gates

## Progress
- **2026-11-20 16:49 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-22 19:39 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 19:39 · Alexander Sacharov** — Stacked on feat/craszs (PR #201, CRASZS done, not merged yet): core/pr reads a PR by number with state open/merged/closed through gh/glab - reused instead of a second forge client. Open this ticket's PR only after #201 merged.
- **2026-11-22 19:41 · Alexander Sacharov** — decisions (worker, settled without the person; disagree at acceptance): (1) form: pr#<n>, case-insensitive, stored as pr#<n> in blocked-by beside ticket ids; anything else after pr# is refused by set/create. (2) source: core/pr from CRASZS, origin's forge (fork -> parent), gh/glab; new light pr.State asks only the state, not the whole PR. (3) cost: one forge call per PR number per process, answers cached in <git-common-dir>/jaira/pr-state.json - merged forever, open/closed for 2 minutes - so list/next stay fast and every worktree shares the answer. (4) acceptance rule: starting work before the PR merged is worse than waiting a round longer - no gh, not signed in, offline, unknown state all keep the ticket waiting, with the reason in the blocker message. (5) a PR closed without merge keeps it waiting and says so (clear blocked-by or reject). (6) TUI reads the cache only, never the network (instant startup). (7) field name is blocked-by (the board's own spelling); the DoD's blocked_by means that field.
- **2026-11-22 19:47 · Alexander Sacharov** — cache location changed from the first decision note: ~/.jaira/repo/<key>/pr-state.json (Store.RepoStateDir, shared by every worktree of a clone, as the snapshot stamps are) instead of inside .git - jaira already keeps repo-wide state there. Blocker code: gate.CodeWaitsOnPR (exit as blocked on move), list rows '[blocked] [waits for PR #12: open]', list --json waiting_on, next --json waiting_on_pr, next text 'Waiting on a pull request:'. TUI reads the cache only (pr.Checker Offline).
- **2026-11-22 19:52 · Alexander Sacharov** — self-check before move: concurrency - two processes writing pr-state.json: rename is atomic, last writer wins, a lost entry only costs one forge call; Checker memo is per process, gate calls are sequential, TUI builds one per call. Error paths - forge errors never cached, always 'state unknown (...)' and waiting. Input - pr ref table in core/ticket. Permissions - nothing trusted beyond gh/glab's own answer; cache 0600 via CreateTemp. Every forward move of a ticket carrying pr#n also asks (gate.CheckAdvance -> blockedBy) - cached, so a merged PR costs nothing after the first ask. Side fix in scope: create --blocked-by now resolves handles to full ids like set (it stored them raw and they never cleared).
- **2026-11-22 19:53 · Alexander Sacharov** — critique (round 1): core/pr/state.go:19-33 repeats the forge/tool/call setup of Read (core/pr/pr.go:113-128: forgeOf, look, classify closure); extract one helper (connect) both call
internal/tui/model.go:2005 adds a package-level sync.Map for the cache path and builds a new pr.Checker per gateEnv call; the model already caches derived state lazily in fields reset on reload (m.index, linkIndex) - keep one *pr.Checker on the model the same way, reset where m.index is reset; this also gives gateEnv its doc comment back, which the insertion detached
- **2026-11-22 19:55 · Alexander Sacharov** — critique round 1 answered: core/pr connect() shared by Read and State; TUI keeps one pr.Checker on the model, reset on reload like m.index, cache path computed once per model; gateEnv's doc comment back on gateEnv. TUI test now reloads after writing the cache.
- **2026-11-22 19:58 · Alexander Sacharov** — testing (function pass, real gh against the fork's parent BeMuCa/jaira): pr#201 open -> list '[waits for PR #201: open]', pr#193 merged -> next hands it out; first list 2.3s (2 PRs x 2 calls), second 16ms from the cache. gh off PATH -> 'state unknown (forge tool missing ...)', next count 0, exit 0, no cache written. set blocked-by=pr#abc exit 3. Two parallel commands on an empty cache: both answers in the file, no temp left. Pitfall: a scratch repo whose origin is the fork pushes every created ticket ref there - two scratch refs (P5VV50, T96HBG) landed on sashasoft90/jaira and were deleted again at once; set the push url to nowhere before creating tickets in a fixture.
- **2026-11-22 19:59 · Alexander Sacharov** — bughunt: change: blocker: internal/tui/model.go:2043 - class: per-model state derived from m.store that switchBoard (model.go:391) does not reset; prCacheFile is set once per model and reload clears only prCheck, so after a board switch every pr#<n> is answered from the previous repository's pr-state.json; example: repo A caches PR 12 as merged (kept for good), the TUI switches to repo B, a B ticket with blocked-by pr#12 shows as actionable and the TUI move gate (applyMove -> m.gateEnv) lets it through although B's PR 12 never merged; also the reverse (A's open/closed keeps B waiting); blame: b38b1bb (names 2EF48S); ran: scratch test in internal/tui (newTestModel, cache merged in A, switchBoard to a second store, ticket blocked-by pr#12, gate.Actionable == true) - go test failed as predicted, test removed
- **2026-11-22 19:59 · Alexander Sacharov** — bughunt: done at cee1269c9 - 1 finding (1 change, 0 pre-existing) in 9 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/gate/gate.go core/gate/prblock_test.go core/link/link.go core/link/prref_test.go core/pr/cache.go core/pr/cache_test.go core/pr/pr.go core/pr/state.go core/pr/state_test.go core/release/notes/2EF48S.md core/role/builtin/jaira-teamlead/SKILL.md core/shipped/shipped.sha256 core/ticket/id.go core/ticket/prref.go core/ticket/prref_test.go core/validate/prref_test.go core/validate/validate.go docs/COMMANDS.md internal/cli/flow.go internal/cli/pr.go internal/cli/prblock_test.go internal/cli/root.go internal/cli/tickets.go internal/tui/linklines.go internal/tui/model.go internal/tui/prblock_test.go
- **2026-11-22 20:00 · Alexander Sacharov** — testing: class: TUI state derived from m.store kept across a store change; example: prCacheFile set once per model, switchBoard to another repo answers pr#12 from the old repo's cache (bughunt note); fix: tie the cached path to the store root it was computed for, test switching stores
- **2026-11-22 20:01 · Alexander Sacharov** — fix for the bughunt blocker: prWait resets prCacheRoot/prCacheFile/prCheck whenever m.store.Root differs from the root the path was computed for - covers switchBoard and any other store change, not only the one path; test internal/tui/prblock_test.go TestBoardSwitchReadsTheNewRepositorysPRCache. Other per-model state of this change: none (the CLI builds its checker per command).
- **2026-11-22 20:02 · Alexander Sacharov** — bughunt: done at db4be8be7 - 0 findings in 1 module, follow-up hunt since cee1269c9 (diff in internal/tui/model.go prWait and its test only, no new files); concurrency, error paths, input, permissions checked; files: core/gate/gate.go core/gate/prblock_test.go core/link/link.go core/link/prref_test.go core/pr/cache.go core/pr/cache_test.go core/pr/pr.go core/pr/state.go core/pr/state_test.go core/release/notes/2EF48S.md core/role/builtin/jaira-teamlead/SKILL.md core/shipped/shipped.sha256 core/ticket/id.go core/ticket/prref.go core/ticket/prref_test.go core/validate/prref_test.go core/validate/validate.go docs/COMMANDS.md internal/cli/flow.go internal/cli/pr.go internal/cli/prblock_test.go internal/cli/root.go internal/cli/tickets.go internal/tui/linklines.go internal/tui/model.go internal/tui/prblock_test.go
- **2026-11-22 20:18 · Alexander Sacharov** — Accepted from signoff: angenommen:
- note: githubState calls githubTarget (gh repo view) for every PR number, so each uncached PR costs two gh calls instead of one plus one per process - memo the target on the Checker (testing measured 2.3s for 2 PRs on a cold cache)
- note: an open PR is re-asked every 2 minutes, so on a board with several open pr# blockers every list/next after the window pays ~1s per open PR - acceptable by decision (3), but it is the one place jaira list is no longer instant
- note: docs/COMMANDS.md says next lists tickets that wait "only" for a pull request, while waitsOnPR also lists a ticket blocked by a PR and a ticket together - drop "only"
- note: the TUI move gate refuses a ticket whose PR merged but whose cache holds no fresh answer ("not asked here") until list/next runs - intended by decision (6), worth one line in the TUI notice
