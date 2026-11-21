---
id: 01M7ZYEXTPHSC0PNW5STG8HNHT
title: "Ohne Bughunt auf HEAD kein signoff: der Gate verweigert den Schritt, wenn bughunt.state fail ist"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Ticket, das Code geändert hat, kommt nur nach signoff, wenn der letzte Bughunt den Stand HEAD abdeckt."
context: |-
  Was heute schiefging (21.11.): 1MYFKE kam nach signoff, obwohl nach dem letzten Bughunt (10e5c341c) noch ein Commit (c3728ba70) die gejagten Dateien änderte. jaira show 1MYFKE --brief zeigte bughunt.state = fail, niemand hat es aufgehalten.
  Heute prüft das nur der Teamlead von Hand (grep auf "— bughunt: done at <sha>") und das Cockpit zeigt es im Explain an; der Gate selbst kennt bughunt nicht.
  Die Daten gibt es schon: jaira show <id> --brief liefert bughunt.at und bughunt.state (ok | fail | unknown | skipped, null = kein Hunt), eingeführt mit 0.5.15.
  Gewünscht: jaira move <id> --to signoff (und jeder Weg in einen Personen-Lane hinter testing) wird mit exit 3 verweigert, wenn das Ticket Code ausserhalb von .jaira/ geändert hat und bughunt.state fail oder null ist; skipped (stop-after-fix) und ein Ticket ohne Codeänderung gehen durch. Die Meldung nennt den Hunt-Commit und HEAD und sagt: /jaira-role-bughunt <id> laufen lassen.
  unknown (uncommittete Änderungen) ebenfalls verweigern, mit eigener Meldung "erst committen".
  Person mit --force kommt durch.
definition-of-done: "move nach signoff mit bughunt.state fail oder ohne Hunt bei Codeänderung: exit 3 mit Hunt-Commit, HEAD und Hinweis (Test in core/gate)"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-11-21T20:58:41Z
updated-at: 2026-11-21T21:13:45Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2700790
claimed-at: 2026-11-21T20:59:16Z
outcome-what: "jaira move into signoff refuses a ticket that changed code and whose newest bug hunt does not cover HEAD (needs_bughunt), or whose hunted files are uncommitted (hunt_uncommitted), on boards whose testing lane owes hunt-verdict"
outcome-why: 1MYFKE reached signoff on 21.11. with a commit after its last hunt and bughunt.state fail; only the teamlead checked that by hand
outcome-resolves: "G8HNHT: the gate reads bughunt.state; ok, skipped and tickets without code changes pass; --force overrides"
review-summary: none
review-gaps: "folded TouchesCode (one caller) into gitrepo.HuntCheck; no existing helper lists a commit set's files (Diff/Stat render text) so the git show stays; left gitrepo.HuntReport beside gate.HuntReport (import cycle, see note); proofs re-checked, none pointed at TouchesCode"
test-verdict: "pass: go vet + go test -count=1 ./... green (RC=0); DoD 1-5 verified (tests named in proofs run); scratch board with testing lane via built binary: no hunt → exit 3 needs_bughunt, hunt at HEAD passes, uncommitted → 'commit them first', commit after hunt → refused naming hunt sha and HEAD, --force lands in signoff exit 0, board without testing lane → no hunt refusal; change writes nothing (gate only), so no two-writer case"
hunt-verdict: "fail: bughunt done at 2551070, 1 change blocker (board in a repo subdirectory: ticket-file-only commit reads as code → false needs_bughunt), 1 change edge (uncommitted-only code passes; note for the person)"
round: 2
back-from: testing
back-by: lane
back-count: 1
---

# Ohne Bughunt auf HEAD kein signoff: der Gate verweigert den Schritt, wenn bughunt.state fail ist

## Definition of Done

- [x] move nach signoff mit bughunt.state fail oder ohne Hunt bei Codeänderung: exit 3 mit Hunt-Commit, HEAD und Hinweis (Test in core/gate)
  proof: core/gate/bughunt_test.go TestSignoffNeedsHuntOnHEAD (stale hunt, no hunt); core/gate/gate.go huntAgainstMove; exit 3 via refusalCode→ExitValidation
- [x] skipped, ok und Tickets ohne Codeänderung ausserhalb .jaira/ gehen durch (Tabellentest)
  proof: core/gate/bughunt_test.go TestSignoffNeedsHuntOnHEAD (covers HEAD, stop-after-fix, no code changed); core/gitrepo/huntcheck_test.go TestHuntCheck (board data only)
- [x] unknown wird mit 'erst committen' verweigert; --force der Person kommt durch (Test)
  proof: core/gate/bughunt_test.go TestSignoffNeedsHuntOnHEAD/uncommitted; core/move/move_test.go TestForceOverridesStaleHunt
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/G8HNHT.md
- [x] Ein Board ohne Lane, die hunt-verdict produziert, verhält sich wie vorher (Test)
  proof: core/gate/bughunt_test.go TestHuntRuleScope (board without hunt-verdict, Hunt never called)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-21 20:59 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock
- **2026-11-21 21:00 · Alexander Sacharov** — decision (worker, settled): the gate guards only lanes with requires-human-exit (signoff), not the human lane — human takes questions, and testing's own max-rounds escape goes into human; blocking a question on a missing hunt would trap that. 'code changed' = a commit of the ticket (CommitsForTicket) touches a path outside .jaira/ (gitrepo anyCode rule); uncommitted code with no commits yields no hunt requirement but the unknown state still refuses. Check is injected into gate.Env like DeriveCommits, so gate stays filesystem-free.
- **2026-11-21 21:00 · Alexander Sacharov** — question (open, for the person): A fresh 'jaira init' board has signoff but no testing lane and no bughunt. Should the gate apply there too? 1) Only on boards with a lane that produces hunt-verdict (Recommended) — boards without hunts unchanged, no surprise for users who never run bughunt. 2) On every board — every code ticket needs a /jaira-role-bughunt note before signoff, or --force.
- **2026-11-21 21:00 · Alexander Sacharov** — decision: Variante 1 (Teamlead 21.11.): der Gate gilt nur auf Boards, deren Lanes hunt-verdict produzieren; Boards ohne Bughunt ändern sich nicht. Gate nur vor Lanes mit requires-human-exit (signoff), nicht vor human.
- **2026-11-21 21:06 · Alexander Sacharov** — in-progress: hunt state logic moved from internal/cli/tickets.go to core/gitrepo/hunt.go (LastHunt, TouchesCode, HuntCheck) so CLI and TUI share it; gitrepo cannot import gate (core/board tests import gitrepo → import cycle), so gitrepo.HuntReport mirrors gate.HuntReport and callers convert. A TUI move (Interactive) is not refused, like the verdict rail: a person in the board decides. 'unknown' also covers git failing to answer; the message says so.
- **2026-11-21 21:06 · Alexander Sacharov** — critique (pass 1, read the diff): let stand — (a) gitrepo.HuntReport mirroring gate.HuntReport: forced by the import cycle (in-progress note), a struct conversion keeps the callers one line each; (b) state strings as literals in gate: gate stays git-free by design, like DeriveCommits; (c) blocked→signoff re-checks the hunt: the ticket returns to a person's lane, the check belongs there too. No finding.
- **2026-11-21 21:09 · Alexander Sacharov** — testing: note: refusal names HEAD as the full 40-char sha beside the hunt's short one (HEAD 04e7f55485e0…); cosmetic, readable; fix would be rev-parse --short in HuntCheck
- **2026-11-21 21:10 · Alexander Sacharov** — bughunt: change: blocker: core/gitrepo/hunt.go:97-111 (HuntCheck -> anyCode -> isBoardData core/gitrepo/git.go:409) - class: a board that sits in a subdirectory of the git repository (a supported layout: jaira init/share are tested there, internal/cli/startupgit_test.go:75, privateexclude_test.go:324); git show --name-only prints repo-root paths ('sub/.jaira/tickets/X.md'), isBoardData only recognises a '.jaira/' prefix at the root, so board data reads as code; the same root anchoring sits in codeOnly (git.go:105, ':(exclude,top).jaira/...') that HuntCoverage falls back to for a hunt without a file list, so on such a board a hunt can also read fail because of the ticket's own ticket-file commits; example: repo with board in sub/, the ticket's only commit is chore(G8HNHT) carrying sub/.jaira/tickets/<id>.md -> HuntCheck returns CodeChanged:true State:'' -> the move into signoff is refused with needs_bughunt (exit 3) although the ticket changed no code; ran: go test -overlay with a scratch test in core/gitrepo (TestScratchSubdirBoard logged {CodeChanged:true State: ...}); isBoardData/codeOnly are older (DF0YD6, a3deff81d), the new gate call makes it a refusal
- **2026-11-21 21:10 · Alexander Sacharov** — bughunt: change: core/gitrepo/hunt.go:98-104 - class: code changes that exist only in the worktree (not committed) never count as code changed, because HuntCheck derives CodeChanged from the ticket's commits alone; example: a conversational ticket whose worker hands back a commit line, so all code is uncommitted and no hunt ran -> CommitsForTicket finds no code commit -> HuntReport{} -> the move into signoff passes without any hunt, while the same ticket with one code commit plus uncommitted edits is refused (hunt_uncommitted); ran: read (HuntCheck returns before LastHunt when no listed commit touches code; HuntCoverage's worktree check only runs when a hunt exists). Edge, not a blocker: the gate lets this through, it does not refuse wrongly or lose data
- **2026-11-21 21:11 · Alexander Sacharov** — bughunt: done at 2551070be - 2 findings (2 change, 1 of them blocker, 0 pre-existing) in 5 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/gate/bughunt_test.go core/gate/gate.go core/gitrepo/hunt.go core/gitrepo/huntcheck_test.go core/move/move_test.go core/release/notes/G8HNHT.md internal/cli/root.go internal/cli/tickets.go internal/tui/model.go
- **2026-11-21 21:11 · Alexander Sacharov** — testing: class: board-data detection compares repo-root-relative git paths against a board-relative '.jaira/' prefix; example: board at sub/.jaira, commit carrying only sub/.jaira/tickets/X.md → CodeChanged true, refused; fix: strip the board's prefix (git rev-parse --show-prefix) before isBoardData, in anyCode and codeOnly alike, test with a board in a subdirectory. Edge (uncommitted-only code passes the gate) stays a note for the person: refusing it would also stop a conversational ticket whose code the person has not committed yet.
- **2026-11-21 21:13 · Alexander Sacharov** — in-progress (round 2): anyCode is now a Repo method that strips the board's prefix (git rev-parse --show-prefix) before isBoardData; a path outside the board's directory counts as code. This fixes HuntCheck and HuntCoverage's codeOnly fallback alike, since both filter git's output through anyCode. Pre-existing and left alone: codeOnly's ':(exclude,top).jaira/tickets' in Diff/WorktreeDiff has the same root anchoring for a subdirectory board (ticket files show in the review diff there) — not this ticket's subject.
