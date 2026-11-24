---
id: 01M86HDDAQ3ER9GEJCYVAA2C3G
title: TestCurrentHonoursJairaUserAboveEverything liest den echten people-Ref und ist in jedem Klon mit Aliassen rot
status: in-progress
ready: true
creator: Sascha
assignee: Sascha
goal: "Die Tests in core/identity lesen nie den people-Ref des Repos, in dem sie laufen; sie sind in jedem Klon gruen, egal welche Aliasse dort stehen."
context: |-
  Was falsch ist: core/identity/identity_test.go:14 TestCurrentHonoursJairaUserAboveEverything erwartet Current()='berk' bei JAIRA_USER=berk, bekommt aber 'Berk'. Grund: identity.People liest den people-Ref dieses Klons (gitref.PeopleRef), und seit 24.11.2026 steht dort 'Berk - also BeMuCa, berkc'.
  Folge: go test ./... ist in jedem lokalen Klon mit gefetchtem refs/jaira/people rot; der Gate jedes Workers und von jaira-role-pr bleibt stehen (8NWE4P am 24.11. gestoppt). Master-CI ist gruen, weil CI-Klone den Ref nicht haben.
  Test aus DVPAPE/Q030WZ. Gefunden im Review zu WV2J2W und beim PR von 8NWE4P.
definition-of-done: "Die identity-Tests laufen gegen ein eigenes temporaeres Repo bzw. einen leeren people-Ref und sind gruen in einem Klon, dessen refs/jaira/people 'Berk - also berkc' enthaelt; Test im Repo-Checkout mit gesetztem Ref nachgestellt"
tags:
  - ci
blocked-by: []
related: []
commits:
  - 7f6c4f59b
  - 7f6c4f59b d8c662e2c
created-at: 2026-11-24T10:25:21Z
updated-at: 2026-11-24T10:44:52Z
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-74925
claimed-at: 2026-11-24T10:25:54Z
outcome-what: "identity test runs in t.TempDir(); testhome.Run unsets git's repository variables (GIT_DIR, GIT_WORK_TREE, ...) and sets GIT_CEILING_DIRECTORIES=os.TempDir(), with a test over both cases"
outcome-why: "the checkout's refs/jaira/people turned go test ./... red in every clone with aliases; under an exported GIT_DIR the same tests even overwrote that ref and .git/config"
outcome-resolves: "DoD 1-3: identity tests green in a clone with 'Berk: BeMuCa, berkc', no test reads or writes the developer's repo, also under GIT_DIR/TMPDIR"
review-summary: none
review-gaps: none
test-verdict: "pass: go build + go test -count=1 ./... green (RC=0) after d8c662e2c; DoD 1-3 verified: identity tests green in this worktree with 'Berk: BeMuCa, berkc'; GIT_DIR=<scratch clone> 4 FAIL before/ok after, TMPDIR inside the clone 1 FAIL before/ok after, clone's people ref unchanged; TestIsolateGit red with isolateGit disabled; change writes nothing outside tests, no concurrent writer to check"
hunt-verdict: "fail: bughunt done at d8c662e, 0 change, 1 pre-existing same subject (packages without testhome.Run: board, gitref, release, shipped, snapshot, scripts) folded into this ticket"
round: 3
back-from: testing
back-by: lane
back-count: 2
---

# TestCurrentHonoursJairaUserAboveEverything liest den echten people-Ref und ist in jedem Klon mit Aliassen rot

## Definition of Done

- [x] Die identity-Tests laufen gegen ein eigenes temporaeres Repo bzw. einen leeren people-Ref und sind gruen in einem Klon, dessen refs/jaira/people 'Berk - also berkc' enthaelt; Test im Repo-Checkout mit gesetztem Ref nachgestellt
  proof: core/identity/identity_test.go:13 TestCurrentHonoursJairaUserAboveEverything uses t.TempDir(); red before/green after in this worktree whose refs/jaira/people has 'Berk: BeMuCa, berkc'
- [x] Kein anderer Test in core/identity, core/people oder internal/cli liest den echten people-Ref (grep/Lauf im echten Checkout gruen)
  proof: grep: only identity_test.go:13 passed '.'; mergebranches_test.go repoRoot only feeds go build; go test ./core/identity ./core/people ./internal/cli green in this checkout with the ref set
- [x] Auch mit exportiertem GIT_DIR oder einem TMPDIR innerhalb eines Repos erreicht kein Test ueber testhome.Run das Repo des Entwicklers (lesen oder schreiben)
  proof: internal/testhome/testhome.go isolateGit; TestIsolateGitKeepsTempDirsOutOfTheDevelopersRepository (both cases red without isolateGit); GIT_DIR=<scratch clone> and TMPDIR=<inside scratch clone> go test ./core/identity: 4 resp. 1 FAIL before, ok after, clone's people ref unchanged
- [x] Jedes Paket, dessen Tests git aufrufen, hat einen TestMain mit testhome.Run; ein Test in internal/testhome schlaegt fuer jedes fehlende fehl
  proof: core/{board,gitref,release,shipped,snapshot}/main_test.go, scripts/main_test.go; internal/testhome TestEveryPackageThatRunsGitInTestsCallsRun (fails naming core/board with its main_test.go removed); GIT_DIR=<scratch repo> go test over all eight packages: scratch repo's refs and config unchanged

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-24 10:25 · Sascha** — parent session: uds:/run/user/1000/cc-socks/3417.sock
- **2026-11-24 10:26 · Sascha** — no open decisions: fix is to run the test in t.TempDir() (no repo, so no people ref) instead of '.'; reproduced red in this worktree whose refs/jaira/people has 'Berk: BeMuCa, berkc'
- **2026-11-24 10:29 · Sascha** — four categories: concurrency n/a (test only); error paths: t.TempDir is outside any repo so People() gets git's error and returns the empty list; unusual input n/a; permissions n/a. Only sibling candidate internal/cli repoRoot() feeds 'go build' only, reads no identity.
- **2026-11-24 10:29 · Sascha** — critique+optimize in one pass (3 lines): t.TempDir() is the pattern TestCurrentFallsThroughWhenJairaUserIsEmpty already uses in this file; nothing simpler, nothing to remove
- **2026-11-24 10:31 · Sascha** — bughunt: change: core/identity/identity_test.go:15 - class: the test is isolated from the checkout only by its working directory, but git resolves the repository from GIT_DIR/GIT_WORK_TREE first and otherwise walks up from -C dir, so Current(t.TempDir()) still reads the checkout's refs/jaira/people whenever GIT_DIR is exported (any git hook, a wrapper script) or TMPDIR lies inside a repository; example: with this checkout's people list 'Berk: BeMuCa, berkc' the test gets Current() = "Berk", want "berk"; not a blocker: plain go test and CI pass; fix direction: build the dir as its own git repo (git init) or set GIT_CEILING_DIRECTORIES and unset GIT_DIR in the test or in internal/testhome; ran: GIT_DIR=$(git rev-parse --git-dir) go test ./core/identity/ -run TestCurrentHonoursJairaUserAboveEverything (exit 1) and TMPDIR=<dir inside checkout> go test ... (exit 1)
- **2026-11-24 10:31 · Sascha** — bughunt: pre-existing: core/identity/identity_test.go:167 (gitInit) and :228/:269 (update-ref refs/jaira/people) - class: test helpers that write git state into a temp dir write into the real repository when GIT_DIR is exported, because git -C dir honours GIT_DIR over dir; also TestAliasesIncludeTheGitEmail, TestAliasesReadTheAliasFile, TestCurrentIsTheJairaNameThePeopleListGives fail then; example: GIT_DIR=$(git rev-parse --git-dir) go test ./core/identity/ set user.name=BeMuCa and user.email=berk@example.test in the main checkout's .git/config and replaced the shared refs/jaira/people with the test's 'grace:' list (data loss of the board's people list, no reflog for that ref); blame: e51f9ae, 5f2c28b (before the ticket); internal/testhome.Run does not clear GIT_DIR/GIT_WORK_TREE; ran: the command above, then restored refs/jaira/people to f35141c (= upstream) and unset the two config keys by hand
- **2026-11-24 10:31 · Sascha** — hunt report: 2 findings (1 change edge, 1 pre-existing), both as notes on AA2C3G | read whole: core/identity (identity.go, identity_test.go, main_test.go) plus internal/testhome; nothing as a diff (first hunt) | could not check whether any teammate's hook runs go test with GIT_DIR set; the repo itself ships no hook; my reproduction overwrote refs/jaira/people and .git/config locally and I restored both (people ref = upstream f35141c, user.* back from ~/.gitconfig)
- **2026-11-24 10:31 · Sascha** — bughunt: done at 7f6c4f59b - 2 findings (1 change, 1 pre-existing) in 1 module, first hunt; concurrency, error paths, input, permissions checked; files: core/identity/identity_test.go
- **2026-11-24 10:33 · Sascha** — testing: class: a test's git calls find the developer's repository through the environment, not through its dir - GIT_DIR/GIT_WORK_TREE exported, or TMPDIR inside a checkout; example: bughunt change note (GIT_DIR set -> Current()='Berk'), and the pre-existing one where the same env made gitInit/update-ref overwrite the real refs/jaira/people; fix: at the cause, internal/testhome.Run unsets git's repository variables and sets GIT_CEILING_DIRECTORIES to os.TempDir(), covering every package that calls it
- **2026-11-24 10:38 · Sascha** — fix in testhome, not per test: 21 packages call testhome.Run, so every git-touching test is covered; GIT_CEILING_DIRECTORIES=os.TempDir() handles TMPDIR inside a checkout. Repro only in a scratch clone - the first hunt's repro with GIT_DIR on the real checkout overwrote refs/jaira/people (restored to upstream f35141c). Dead end in my own test: GIT_DIR set to '' makes git fail outright, so the case must unset it
- **2026-11-24 10:38 · Sascha** — critique round 2 (reads only the testhome change answering the hunt): fix sits at the one place all 21 test packages pass through (testhome.Run), same shape as its existing JAIRA_* unsets; no finding
- **2026-11-24 10:38 · Sascha** — optimize round 2: nothing to remove; gitRepoVars is a list because Run's own JAIRA_* unsets are one
- **2026-11-24 10:39 · Sascha** — waiting: bughunt round 2 (subagent) since 11:39
- **2026-11-24 10:40 · Sascha** — bughunt: pre-existing: core/gitref/gitref_test.go:38, core/snapshot/snapshot_test.go:25, core/shipped/history_test.go:128 (auch core/release/notes_test.go:32, scripts/*_test.go:58,126) - Klasse: Testpakete, die git in t.TempDir() aufrufen, aber kein TestMain mit testhome.Run haben, bekommen die neue Git-Isolation (isolateGit) nicht; ein exportiertes GIT_DIR (Hook/Wrapper) lenkt ihre git-Aufrufe ins Repo des Entwicklers; Beispiel: GIT_DIR=<wegwerf-repo>/.git go test ./core/gitref ./core/snapshot ./core/shipped -> refs/heads/master des Repos von 06a655d auf 7203dbd umgebogen, user.name/email, remote upstream, jaira.remote in .git/config ueberschrieben, fremde Branches a-d/main/feature und Tags v0.1.0-v0.3.0 angelegt; Gegenprobe core/identity, core/people (mit testhome.Run): Repo unveraendert; Testdateien nicht in Commits von AA2C3G (letzte Aenderung 996c3e735); ran: siehe oben, nur gegen Wegwerf-Repos im scratchpad
- **2026-11-24 10:40 · Sascha** — hunt report: 1 Befund (0 change, 1 pre-existing), als Note auf AA2C3G | ganz gelesen: internal/testhome/testhome.go, internal/testhome/testhome_test.go (neu in der Liste); als Diff: core/identity/identity_test.go (kein Diff seit 7f6c4f59b) | nicht geprueft: Windows/macOS-Pfade von GIT_CEILING_DIRECTORIES (nur Linux; Slash am Ende von TMPDIR geprueft, harmlos)
- **2026-11-24 10:40 · Sascha** — bughunt: done at d8c662e2c - 1 finding (0 change, 1 pre-existing) in 2 modules, follow-up hunt; concurrency, error paths, input, permissions checked; files: core/identity/identity_test.go internal/testhome/testhome.go internal/testhome/testhome_test.go
- **2026-11-24 10:42 · Sascha** — testing round 2: class: test packages that run git but have no TestMain with testhome.Run miss isolateGit; example: GIT_DIR=<scratch> go test ./core/gitref ./core/snapshot ./core/shipped rewrote that repo's master, .git/config, tags; fix: main_test.go with testhome.Run in all six, plus a guard test in internal/testhome that fails for any package running git in tests without it. Folded here, not a follow-up ticket: same subject as the goal
