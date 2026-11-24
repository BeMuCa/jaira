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
commits: []
created-at: 2026-11-24T10:25:21Z
updated-at: 2026-11-24T10:29:31Z
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-74925
claimed-at: 2026-11-24T10:25:54Z
---

# TestCurrentHonoursJairaUserAboveEverything liest den echten people-Ref und ist in jedem Klon mit Aliassen rot

## Definition of Done

- [x] Die identity-Tests laufen gegen ein eigenes temporaeres Repo bzw. einen leeren people-Ref und sind gruen in einem Klon, dessen refs/jaira/people 'Berk - also berkc' enthaelt; Test im Repo-Checkout mit gesetztem Ref nachgestellt
  proof: core/identity/identity_test.go:13 TestCurrentHonoursJairaUserAboveEverything uses t.TempDir(); red before/green after in this worktree whose refs/jaira/people has 'Berk: BeMuCa, berkc'
- [x] Kein anderer Test in core/identity, core/people oder internal/cli liest den echten people-Ref (grep/Lauf im echten Checkout gruen)
  proof: grep: only identity_test.go:13 passed '.'; mergebranches_test.go repoRoot only feeds go build; go test ./core/identity ./core/people ./internal/cli green in this checkout with the ref set

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-24 10:25 · Sascha** — parent session: uds:/run/user/1000/cc-socks/3417.sock
- **2026-11-24 10:26 · Sascha** — no open decisions: fix is to run the test in t.TempDir() (no repo, so no people ref) instead of '.'; reproduced red in this worktree whose refs/jaira/people has 'Berk: BeMuCa, berkc'
- **2026-11-24 10:29 · Sascha** — four categories: concurrency n/a (test only); error paths: t.TempDir is outside any repo so People() gets git's error and returns the empty list; unusual input n/a; permissions n/a. Only sibling candidate internal/cli repoRoot() feeds 'go build' only, reads no identity.
