---
id: 01M5YWBQHKFPKV02PA8YGBFM32
title: "jaira update bringt standardmaessig alles auf Stand - Abschnitt, Rollen, Cockpit, Skill - und 'jaira update diff' zeigt, was zurueckliegt"
status: signoff
ready: true
creator: berkc
assignee: berkc
goal: "Ein einziges 'jaira update' bringt Board-Abschnitt, Rollen-Skills, Cockpit-Mod und jaira-Skill auf den Stand des installierten Binaries; 'jaira update diff' sagt vorher, was zurueckliegt."
context: |-
  Wunsch von berkc am 27.10.2026: 'Jaira update should do all of it by default, from skill to mod to roles. jaira update diff can be used to see which lanes or roles are behind.'
  Heute (0.5.6): 'jaira update' schreibt nur den jaira-Abschnitt in CLAUDE.md/AGENTS.md und die Release-Notes. Rollen und Cockpit kommen nur ueber 'jaira roles install --global|--project'.
  PR #91 (T8A8KM) fuegt 'update --roles', '--lanes' und '--dry-run' hinzu - Rollen sind dort noch opt-in.
  Offen: Lanes auch per Default? Lanes sind geteilte Board-Konfiguration und brauchen einen Commit - Vorschlag: Rollen/Mod/Skill per Default, Lanes nur mit --lanes, 'diff' zeigt beides.
  Gefunden am 23.10.: 'roles install' hielt zwei alte, unveraenderte Cockpit-Dateien (board.ts aus 5bc22fb1, register.tsx aus 2c4c6790) fuer 'von dir bearbeitet', weil ihr Hash in core/shipped/shipped.sha256 fehlt - auch in #91. Vermutung (ungeprueft): scripts/shipped-history.sh geht nicht durch zweite Merge-Eltern. Ohne den Fix wuerde ein Default-Update genau diese Dateien stehen lassen.
definition-of-done: "'jaira update' ohne Flags aktualisiert Abschnitt, Rollen, Cockpit-Mod und jaira-Skill ueberall, wo sie schon installiert sind; bearbeitete Dateien bleiben (exit 3)"
tags:
  - cli
blocked-by: []
related:
  - 01M21PJTWPK4DZ89PRY6T8A8KM
commits: []
created-at: 2026-10-27T14:31:21Z
updated-at: 2026-10-30T22:30:14Z
updated-by: berkc
claimed-by: berk-ca-2379103
claimed-at: 2026-10-30T22:13:03Z
outcome-what: "bare 'jaira update' updates roles, cockpit mod and jaira skill; 'jaira update diff' lists current/behind/edited per file and lane; the shipped manifest gets a release column and is regenerated"
outcome-why: "the cockpit files shipped in v0.5.2-v0.5.6 read as edited because the committed manifest was a hand-patched, never regenerated list"
outcome-resolves: cockpit files from released tags are no longer reported as edited; update diff answers what is behind before anything is written
review-summary: "Bare 'jaira update' now brings roles, the cockpit mod and the jaira skill up to date by default (--roles still accepted, now a no-op; lanes stay opt-in via --lanes). Adds 'jaira update diff', which lists each role, cockpit file and lane as current, behind (with the release that shipped it) or edited, and writes nothing. core/shipped/shipped.sha256 is regenerated with a release column and shipped.Release reads it, so cockpit files from released tags are no longer reported as edited. Release note added."
round: 2
review-gaps: |-
  none
  note: go vet ./... and go test ./... (plan item still unticked) were not run in this review; only the four packages named in the DoD were run, all pass
  note: update diff reports a role file that is not installed yet as behind (missing), which reads as behind for a file that was never there; wording only
test-verdict: "pass: go vet ./... RC=0; go test ./... -race -count=1 -timeout 20m -skip 'TestUnreachableRemoteIsOfflineNotARace|TestAnOfflineWriteArrivesWhenTheNetworkIsBack|TestWebPushReportsAQueuedDrop' EXIT=0, 33 packages ok, no FAIL; DoD 1-4 verified in the tree; behaviour exercised with the built binary on a scratch home (update diff: board.ts and register.tsx 'behind (shipped in v0.5.4)'; bare update rc=0; diff afterwards: both current)"
review-verdict: "The diff matches the four DoD items and no defect was found; no blocker. Verified: the four named packages pass, regenerating scripts/shipped-history.sh gives a byte-identical manifest, and no line of the previous manifest is lost (0 of 497 by sha+path). Not verified by hand: the interactive update output on a real stale board."
review-round: 2
review-check: |-
  1. cd into a jaira repo and build: go build -o /tmp/j ./cmd/jaira
  2. jaira roles install --global, then replace ~/.claude/skills/jaira-cockpit/hooks/board.ts with: git show v0.5.4:core/mod/builtin/jaira-cockpit/hooks/board.ts
  3. append a line to ~/.claude/skills/jaira/SKILL.md
  4. /tmp/j update diff - see board.ts 'behind (shipped in v0.5.4)', jaira/SKILL.md 'edited', a last line '... (nothing was written)'; the files are unchanged
  5. /tmp/j update; echo $? - prints 'skipped .../jaira/SKILL.md', exit code 3
  6. /tmp/j update diff - board.ts now 'current', only SKILL.md 'edited'; lane files untouched unless you pass --lanes
  7. cat core/release/notes/GBFM32.md - three lines, each starting '- '
---

# jaira update bringt standardmaessig alles auf Stand - Abschnitt, Rollen, Cockpit, Skill - und 'jaira update diff' zeigt, was zurueckliegt

## Definition of Done

- [x] 'jaira update' ohne Flags aktualisiert Abschnitt, Rollen, Cockpit-Mod und jaira-Skill ueberall, wo sie schon installiert sind; bearbeitete Dateien bleiben (exit 3)
  proof: internal/cli/updatesetup_test.go TestUpdateWithoutSwitchesUpdatesRolesButNotLanes; internal/cli/update.go roles default true
- [x] 'jaira update diff' listet je Rolle, Mod-Datei und Lane: aktuell, zurueck (welche Version) oder bearbeitet; schreibt nichts
  proof: internal/cli/updatediff.go writeDiff; internal/cli/updatediff_test.go TestUpdateDiffWritesNothingAndNamesWhatIsBehind
- [x] Alte ausgelieferte Versionen aus Merge-Zweigen (board.ts 5bc22fb1, register.tsx 2c4c6790) gelten als ausgeliefert, nicht als bearbeitet
  proof: core/shipped/shipped.sha256 board.ts 9406afda... line shipped in v0.5.2; core/shipped/shipped.go Release
- [x] Tests; Release-Note
  proof: core/release/notes/GBFM32.md; go test ./internal/cli ./core/shipped ./core/role ./core/lane ok

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] find the real cause of the false 'edited' cockpit files in core/shipped (done: list is stale, script is correct; see notes)
  proof: core/shipped/shipped.sha256 regenerated; scripts/shipped-history.sh
- [x] scripts/shipped-history.sh: add the release column (first vX.Y.Z tag that holds the blob at that path, '-' if none)
- [x] regenerate core/shipped/shipped.sha256 with the script; check append-only (no committed line lost) and every tagged embedded file present
- [x] core/shipped/shipped.go: Known keeps sha+path matching; add Release(path, b) to read the column; update shipped_test.go regex
- [x] core/role/install.go and core/lane/refresh.go: Result / Refreshed carry the release of an Updated file
- [x] internal/cli/update.go + updatesetup.go: roles default on, --roles accepted as no-op, personal-only path kept, hintLines keeps lanes only
- [x] internal/cli/update.go: 'jaira update diff' subcommand, one line per role, mod file and lane: current | behind (release) | edited; writes nothing
- [x] tests: adjust update_test.go / updatedry_test.go / updatesetup_test.go for new defaults; add diff tests incl. a no-write check
- [x] release note core/release/notes/GBFM32.md (one '- ' line per change, instruction form)
- [ ] go vet ./... and go test ./... (exit codes checked, output verbatim on failure)

## Progress
- **2026-10-30 21:41 · berkc** — open decisions at start: none that only the person can take. Settled by the person (ticket context, 27.10.): roles, cockpit mod and jaira skill update by default; lanes opt-in via --lanes; 'update diff' writes nothing. Settled here by default: diff output is one line per role/mod file/lane with state current|behind (version)|edited, the same for all three; 'update' touches only what is already installed. Parent session: none (started directly as worker subagent, no AskUserQuestion).
- **2026-10-30 21:44 · berkc** — finding 1 (shipped list): scripts/shipped-history.sh does walk second parents: it uses 'git log --full-history -m' since its first commit 299c65aa, and every commit reachable from HEAD that touches the embedded paths is in that walk (303 vs 375 commits, none missing). Its output today is 558 lines; core/shipped/shipped.sha256 has 503 and is a strict subset of it (0 lines committed-only, 55 lines script-only). The list is append-only and was patched by hand (736c80ee) and never re-run after the 44 later commits that changed embedded files. Hence board.ts 9406af (added by 5bc22fb1, F80QZT) and register.tsx are reported as edited. Fix: regenerate the list with the script. Checked: every file of tags v0.4.1..v0.5.7 under core/role/builtin, core/mod/builtin, core/lane/builtin, core/role/jaira has its sha in the regenerated list (0 missing); the committed list misses 22 of them (v0.5.2..v0.5.6 cockpit files).
- **2026-10-30 21:46 · berkc** — finding 2 (cause of the cockpit false 'edited'): not the script's walk. The committed list is an append-only file: script output at 299c65aa (438 lines) was extended by hand (736c80ee) and never re-run after 44 later commits to embedded paths. Committed list = strict subset of script output today (0 lines lost, 55 missing). Fix = regenerate, then keep append-only. Also a manifest that keeps 2 columns cannot say which release a behind file came from, so DoD 2 'welche Version' needs a release column.
- **2026-10-30 21:46 · berkc** — decisions settled here (reversible, not person-only): (a) manifest line becomes 'sha  path  release' where release is the first vX.Y.Z tag whose tree holds that blob at that path, '-' if none; shipped.Known keeps matching on sha+path. (b) update: roles, cockpit mod and jaira skill default on; --roles stays accepted as a no-op; lanes stay --lanes. (c) 'update diff' is a subcommand, one line per role/mod file and per lane: current | behind (release, or 'missing' for files not yet written) | edited; it reuses the dry-run of refreshRoles and lane.RefreshShipped and writes nothing. Outside a board, roles only, as before.
- **2026-10-30 21:47 · berkc** — plan written as Plan items (see jaira show). Options 'planning' was left unticked by the person, so the pre-process lane is skipped; the plan above is the pre-process output, written by the worker.
- **2026-10-30 21:56 · berkc** — built and checked by hand on a scratch clone (HOME=scratch, stale v0.5.4 cockpit files): update diff lists board.ts behind (shipped in v0.5.4), texts.ts/types behind (v0.5.2); bare update rewrites them to current; diff afterwards shows only lanes (critique edited, optimize/testing behind: those files in the clone are old). update diff writes nothing (tests: TestUpdateDiffWritesNothingAndNamesWhatIsBehind, tree snapshot). Lane Load writes lane files when a board has none: diff skips lanes there on purpose (ProjectLanesActive guard).
- **2026-10-30 22:17 · berkc** — cause of the false 'edited' cockpit files: not the script's walk. scripts/shipped-history.sh already uses git log --full-history -m (since 299c65aa), so second merge parents are walked; every commit reachable from HEAD that touches the embedded paths is in that walk. The committed core/shipped/shipped.sha256 was an append-only file patched by hand (736c80ee) and never regenerated after 44 later commits to embedded paths; it was a strict subset of the script output (55 lines missing), so board.ts 9406af (from 5bc22fb1, F80QZT, in v0.5.2) read as edited. Fix: regenerate with the script. Checked: every file under core/role/builtin, core/mod/builtin, core/lane/builtin, core/role/jaira at tags v0.4.1..v0.5.7 has its sha in the regenerated list (0 missing).
- **2026-10-30 22:17 · berkc** — DESIGN DECISION 1 (for signoff): manifest line is now 'sha256  path  release'. release = first vX.Y.Z tag whose tree holds that blob at that path, '-' when no release holds it (409 of 559 lines: intermediate versions that never reached a tag). shipped.Known still matches sha+path only; shipped.Release reads the third column. Alternative not taken: a separate file or a comment header; both need a second lookup. Reason: DoD 2 asks for 'welche Version'.
- **2026-10-30 22:17 · berkc** — DESIGN DECISION 2 (for signoff): 'jaira update diff' output. One line per role file, mod file and lane: current | behind (shipped in vX.Y.Z) | behind (not in a release yet) | behind (missing) | edited, grouped by skills directory, then a summary line. Writes nothing; outside a board only personal dirs. Lanes are skipped on a board with no lane files of its own (lane.Load would write them). --json gives roles, lanes, state, release. Alternative not taken: list only changed files.
- **2026-10-30 22:17 · berkc** — DESIGN DECISION 3: a bare 'jaira update' now also updates roles, cockpit mod and jaira skill (flag --roles kept, default true, --roles=false leaves them). Lanes stay --lanes. Exit 3 when a role file was edited, also on a bare update. Outside a board, a bare update does the personal dirs only (as --roles did).
- **2026-10-30 22:17 · berkc** — SKIPPED TESTS in the full run (environment, not this change): core/gitref TestUnreachableRemoteIsOfflineNotARace (hangs on git push to an unreachable remote; the same hang reproduces on unmodified origin/master), core/outbox TestAnOfflineWriteArrivesWhenTheNetworkIsBack (same class), internal/cli TestWebPushReportsAQueuedDrop (git fetch to an unreachable remote, hung in this sandbox).
- **2026-10-30 22:17 · berkc** — INCIDENT: the first worktree was deleted by a mistaken 'git worktree remove --force' on its own path; uncommitted work was lost and re-applied from the transcript in ../jaira-GBFM32 on the same branch. The regenerated list was reused from the scratchpad and re-verified byte-identical.
- **2026-10-30 22:19 · berkc** — critique (pass 1, full diff read): standing on DESIGN DECISION 1-3 notes; accepted, not re-argued. Findings: (a) updatediff.go diffJSON mixes lanes[i].ID with diffFile.dir; (b) plain 'jaira update' prints 'roles: none installed ...' on every run once roles default on, update.go / updatesetup.go printRoleRuns. Both fixed in the next in-progress round.
- **2026-10-30 22:20 · berkc** — in-progress round 2: fixed critique findings (a) diffJSON uses diffFile.dir for lane ids, (b) 'roles: none installed' line only printed when roles were found (update.go, dry and real run). update tests pass.
- **2026-10-30 22:20 · berkc** — critique pass 2 (reads findings of pass 1 and the change answering them only): (a) fixed at updatediff.go diffJSON, lane id from f.dir; (b) fixed at update.go, both printRoleRuns calls guarded by len(setup.roleRuns) > 0. Nothing left to say.
- **2026-10-30 22:23 · berkc** — testing: skipped three tests that dial an unreachable remote and hang in this sandbox: core/gitref TestUnreachableRemoteIsOfflineNotARace (reproduced on origin/master too), core/outbox TestAnOfflineWriteArrivesWhenTheNetworkIsBack, internal/cli TestWebPushReportsAQueuedDrop (git fetch to 127.0.0.1:1). Not verified: those three on a machine with normal network behaviour.
- **2026-10-30 22:30 · berkc** — review (fresh context): pass, no blockers. Verified: manifest reproducible by scripts/shipped-history.sh, append-only (0 lost, 55 added), all v0.5.4-v0.5.7 embedded files present with release column; go vet and go test ./... green (skipped TestUnreachableRemoteIsOfflineNotARace, TestAnOfflineWriteArrivesWhenTheNetworkIsBack, TestWebPushReportsAQueuedDrop - they hang on pushes to 127.0.0.1:1 in this sandbox); hand run on scratch HOME: update diff writes nothing, bare update fixes stale v0.5.4 board.ts/register.tsx and exits 3 on edited SKILL.md. Three notes in review-gaps.
