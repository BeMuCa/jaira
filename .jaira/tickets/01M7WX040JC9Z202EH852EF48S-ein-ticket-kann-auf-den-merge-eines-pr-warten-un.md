---
id: 01M7WX040JC9Z202EH852EF48S
title: "Ein Ticket kann auf den Merge eines PR warten, und nach dem Merge startet sein Worker"
status: critique
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
commits: []
created-at: 2026-11-20T16:35:24Z
updated-at: 2026-11-22T19:55:45Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-81804
claimed-at: 2026-11-22T19:38:55Z
outcome-what: "blocked-by takes pr#<n>; list/next/move ask gh or glab for the pull request's state (core/pr State + Checker, cached in ~/.jaira/repo/<key>/pr-state.json: merged for good, open/closed 2 min); list shows [waits for PR #12: open], list --json waiting_on, next --json waiting_on_pr; TUI reads the cache only; teamlead skill starts a worker for tickets a merge freed"
outcome-why: "a dependency on a pull request lived only in the teamlead's notes, so a merged PR freed nothing and the follow-up waited until somebody remembered"
outcome-resolves: "2EF48S: a ticket can wait for a PR merge and becomes actionable on the next run after it"
review-summary: |-
  core/pr/state.go:19-33 repeats the forge/tool/call setup of Read (core/pr/pr.go:113-128: forgeOf, look, classify closure); extract one helper (connect) both call
  internal/tui/model.go:2005 adds a package-level sync.Map for the cache path and builds a new pr.Checker per gateEnv call; the model already caches derived state lazily in fields reset on reload (m.index, linkIndex) - keep one *pr.Checker on the model the same way, reset where m.index is reset; this also gives gateEnv its doc comment back, which the insertion detached
round: 2
back-from: critique
back-by: lane
back-count: 1
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
  proof: core/pr/cache_test.go TestCheckerCache (one call per process, shared file, merged kept, open re-asked after OpenFor, failures not cached, no temp left) + TestCheckerCacheFailures (corrupt, unwritable, offline); core/pr/cache.go store temp+rename
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
