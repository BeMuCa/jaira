---
id: 01M83W793FFMNNE370MYXSGWPX
title: "Cockpit und jaira pr erkennen eine SSH-Remote 'user@eigener-host:owner/repo.git' (GitHub Enterprise) und fragen gh mit Host und owner/repo"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Für ein Repo auf GitHub Enterprise mit SSH-Remote der Form <user>@<host>:<owner>/<repo>.git finden Cockpit-Abnahme und 'jaira pr' den PR: gh läuft mit --hostname/GH_HOST=<host> und -R <owner>/<repo>."
context: |-
  Was falsch ist (Teamlead eines anderen Repos, 23.11.2026, jaira 0.6.3): die Cockpit-Abnahme prüft PRs mit 'gh repo view' auf der rohen Remote-URL '<user>@<ghe-host>:<owner>/<repo>.git' -> 'Could not resolve to a Repository'.
  Die SSH-Form mit Benutzer vor '@' und eigenem Host wird nicht in Host + owner/repo zerlegt. github.com-Remotes gehen.
  gh braucht GH_HOST=<host> (oder --hostname bei gh api) und -R <owner>/<repo>.
  Stellen: core/pr (Forge/Ziel-Repo aus Remote, CRASZS), Cockpit register.tsx (openPullRequest/hasPullRequest), jaira-role-pr scripts.
definition-of-done: "Remote-Parser zerlegt https://, ssh://, git@host:owner/repo(.git) und <beliebiger-user>@<host>:owner/repo(.git) in host + owner/repo (Tabellentest inkl. GHE-Host, Port, ohne .git)"
tags:
  - cockpit
blocked-by: []
related: []
commits:
  - 91f5e8ab0d48bdffe68ee036127fd6698f90dff9
  - 53d1eb583980f0366bb497f06ca0901dd142cf40
  - df630383fcf6949b0c4124b655e0cb457739c42b
  - 1ed1bc3793db8579dd7f71fbc71e351a8978f6e9
created-at: 2026-11-23T09:36:31Z
updated-at: 2026-11-23T12:03:44Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-61389
claimed-at: 2026-11-23T09:37:02Z
outcome-what: "Remote-Parser (Go core/pr/remote.go, TS board.ts parseRemote) zerlegt jede git-Remote-Form in Host + owner/repo; gh bekommt auf einem Host außerhalb von github.com/*.github.com host/owner/repo und --hostname bei gh api — in jaira pr/State und in der Cockpit-PR-Abfrage."
outcome-why: "gh repo view auf 'user@ghe-host:owner/repo.git' scheitert mit 'Could not resolve to a Repository'; gh versteht nur git@-Präfix bzw. [HOST/]OWNER/REPO."
outcome-resolves: Cockpit-Abnahme und jaira pr finden PRs auf GitHub Enterprise mit SSH-Remote eines beliebigen Users.
review-summary: "Full review (first). core/pr: new private parseRemote takes apart https/ssh/git+ssh/scp remotes (any user@, port dropped, .git dropped); ghRepo turns a non-github.com origin into host/owner/repo plus a host, and connect() then passes host/owner/repo to 'gh repo view' and prefixes every 'gh api' call with --hostname <host>, so Read and State (jaira pr, blocked-by pr#N) ask the GHE instance instead of github.com. github.com and every *.github.com (ssh.github.com, www) stay exactly as before (raw origin, no --hostname). Cockpit: the same parser in board.ts (parseRemote/ghRepoOf/onHost); forgeRepo carries the host, and pullTarget, decideStep and askPullRequest name the repo as host/owner/repo to gh repo view and gh pr list --repo. Release note and shipped.sha256 added. Forge detection is unchanged: a GHE host still needs git config jaira.forge github."
review-gaps: |-
  note: limit: a GHE remote with a port (https://ghe:8443/...) loses the port in --hostname and host/owner/repo - stated in the ticket notes, no instance to check against
  note: limit: an SSH host alias as remote (git@github.com-work:me/repo.git) is asked as host github.com-work - it failed the same way before (GH_DEBUG: gh repo view on the raw URL requests https://github.com-work/api/graphql), no regression
  note: elsewhere: core/role/builtin/jaira-dispatcher/scripts/reap.sh:175/:224 and jaira-role-pr/scripts/pr-round.sh:81 still talk to github.com only on a GHE remote (pre-existing, bughunt note)
  note: elsewhere: the uncommitted AGENTS.md, CLAUDE.md and .jaira/lanes/*.md are a jaira update regeneration to 0.6.3, not this ticket - keep them out of the XSGWPX commits
test-verdict: "pass: round 2 since 53d1eb583: go test -race core/pr internal/cli RC=0, cockpit 501/501 RC=0; DoD 1-3 verified; real binary + fake gh: alice@ghe.example.com:team/app.git -> 'gh repo view ghe.example.com/team/app' + 'gh api --hostname ghe.example.com graphql'; ssh.github.com origin unchanged (raw, no --hostname); gh failing -> jaira pr exit 1 forge_failed, no success claimed; writes nothing"
hunt-verdict: "pass: bughunt done at df630383f, 0 change; round 1: 1 change fixed (df630383f), 1 pre-existing (reap.sh/pr-round.sh github.com-only) left as note for the person"
round: 2
back-from: testing
back-by: lane
back-count: 1
review-verdict: "accept: DoD 1-3 met; go test ./core/pr ./core/shipped ok, claude plugin test 501/501, shipped-history --check green; no defect in core/pr or the cockpit forge path"
review-round: 2
review-verdict-gaps: |-
  note: limit: a GHE remote with a port (https://ghe:8443/...) loses the port in --hostname and host/owner/repo - stated in the ticket notes, no instance to check against
  note: limit: an SSH host alias as remote (git@github.com-work:me/repo.git) is asked as host github.com-work - it failed the same way before (GH_DEBUG: gh repo view on the raw URL requests https://github.com-work/api/graphql), no regression
  note: elsewhere: core/role/builtin/jaira-dispatcher/scripts/reap.sh:175/:224 and jaira-role-pr/scripts/pr-round.sh:81 still talk to github.com only on a GHE remote (pre-existing, bughunt note)
  note: elsewhere: the uncommitted AGENTS.md, CLAUDE.md and .jaira/lanes/*.md are a jaira update regeneration to 0.6.3, not this ticket - keep them out of the XSGWPX commits
review-check: "1. cd core/pr && go test -run 'TestParseRemote|TestReadGitHubEnterprise' -v .  - both PASS; the GHE case shows 'gh repo view ghe.example.com/team/app' and 'gh api --hostname ghe.example.com graphql'  2. cd core/mod/builtin/jaira-cockpit && claude plugin test .  - 0 fail; look for 'GitHub Enterprise behind a user@host:owner/repo remote' (pass)  3. Optional, live with a GHE account: in a repo whose origin is <user>@<ghe-host>:<owner>/<repo>.git run git config jaira.forge github, then jaira pr  - it lists the pull requests instead of 'Could not resolve to a Repository'"
question: ""
---

# Cockpit und jaira pr erkennen eine SSH-Remote 'user@eigener-host:owner/repo.git' (GitHub Enterprise) und fragen gh mit Host und owner/repo

## Definition of Done

- [x] Remote-Parser zerlegt https://, ssh://, git@host:owner/repo(.git) und <beliebiger-user>@<host>:owner/repo(.git) in host + owner/repo (Tabellentest inkl. GHE-Host, Port, ohne .git)
  proof: core/pr/remote_test.go TestParseRemote; tests/board.test.ts describe('parseRemote')
- [x] jaira pr und die Cockpit-PR-Abfrage rufen gh für einen nicht-github.com-Host mit GH_HOST bzw. --hostname und -R owner/repo auf (Test mit Fake-gh)
  proof: core/pr/remote_test.go TestReadGitHubEnterprise (Read+State, Fake-gh: --hostname + host/owner/repo); tests/board.test.ts 'GitHub Enterprise behind a user@host:owner/repo remote'
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/XSGWPX.md; core/shipped/shipped.sha256 (scripts/shipped-history.sh --check grün)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Go: remoteOf(origin) -> host, owner/repo in core/pr + Tabellentest
- [x] Go: connect gibt gh bei Nicht-github.com-Host --hostname (api) und HOST/OWNER/REPO (repo view); Fake-gh-Test
- [x] TS: remoteOf in board.ts + Tabellentest; repoViewCommand/pullRequestsCommand nehmen host; forgeRepo trägt host
- [x] TS-Test mit Fake-gh im Cockpit (register-Ebene)
- [x] Release-Note + shipped.sha256, Gates

## Progress
- **2026-11-23 09:37 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/6958.sock
- **2026-11-23 09:38 · Alexander Sacharov** — decision: keine offenen Fragen. gh bekommt den Host als Argument statt GH_HOST-Env: 'gh api --hostname <host>' und '-R/repo view <host>/<owner>/<repo>' (gh's dokumentierte [HOST/]OWNER/REPO-Form). Grund: Cockpit-$.process.run und pr.Run kennen keine Env, und ein Argument ist im Fake-gh-Test sichtbar. github.com-Remotes werden zu owner/repo normalisiert (bisher roh durchgereicht). Forge-Erkennung bleibt: GHE-Host ohne 'github.com' braucht weiter git config jaira.forge github. pr-round.sh (jaira-role-pr) prüft 'gh auth status --hostname github.com' fest — nicht in diesem Ticket, known limit.
- **2026-11-23 09:39 · Alexander Sacharov** — Go-Teil fertig: core/pr/remote.go ParseRemote+ghRepo, connect() in pr.go setzt --hostname bei gh api und host/owner/repo bei repo view; Tests TestParseRemote, TestReadGitHubEnterprise grün. Nächstes: Cockpit board.ts/register.tsx.
- **2026-11-23 09:43 · Alexander Sacharov** — critique: none. Geprüft: Parser doppelt in Go (core/pr/remote.go) und TS (board.ts parseRemote) ist das bestehende Muster (forgeOf existiert schon in beiden, pr.go:152 verweist darauf); --hostname im Wrapper von connect() (pr.go) ist die Stelle, durch die jeder gh-Aufruf geht; github.com bleibt roh, damit bestehendes Verhalten gleich bleibt. Exportiertes ParseRemote/Remote → optimize.
- **2026-11-23 09:46 · Alexander Sacharov** — bughunt: change: blocker: core/pr/remote.go:65 - class: every host other than exactly github.com/www.github.com is treated as GitHub Enterprise, including subdomains of github.com that gh itself maps to github.com (ssh.github.com, the documented SSH-over-port-443 host); pr.go:151 then adds --hostname <that host> to every gh api call, and gh builds https://api.<host>/graphql from it; example: origin ssh://git@ssh.github.com:443/me/repo.git (forge detected as github by the 'github.com' substring) -> ghRepo gives host=ssh.github.com -> 'gh api --hostname ssh.github.com graphql' fails with 'error connecting to api.ssh.github.com', so jaira pr and State (blocked-by pr#N) break where they worked before this change (gh repo view on the raw URL went to api.github.com); same rule, harmless for now because no --hostname is passed there, in core/mod/builtin/jaira-cockpit/hooks/board.ts ghRepoOf; fix: count a host as github.com when it is github.com or ends in .github.com, as gh does; ran: go test with ghRepo("ssh://git@ssh.github.com:443/me/repo.git") -> host "ssh.github.com"; GH_DEBUG=api gh api --hostname ssh.github.com graphql -> Request to https://api.ssh.github.com/graphql, error connecting; GH_DEBUG=api gh repo view ssh://git@ssh.github.com:443/cli/cli.git -> Request to https://api.github.com/graphql (old path worked)
- **2026-11-23 09:46 · Alexander Sacharov** — bughunt: pre-existing: core/role/builtin/jaira-dispatcher/scripts/reap.sh:175 - class: the role scripts that run after the cockpit's pull-request check still talk to github.com only, so on a GitHub Enterprise remote the check now passes and the next step fails; places: reap.sh:175 and :224 pass the raw origin to 'gh repo view' (the form this ticket fixed in core/pr); jaira-role-pr/scripts/pr-round.sh:81 requires 'gh auth status --hostname github.com' and dies 'not signed in to github.com' for a GHE-only sign-in, and its gh api calls carry no --hostname; example: origin alice@ghe.example.com:team/app.git -> gh repo view asks api.github.com; blame: 305549c6fe (reap.sh:175), 423ec7e1c9 (pr-round.sh:81), both before the ticket; ran: GH_TOKEN=bogus gh repo view alice@ghe.example.com:team/app.git -> HTTP 401 from https://api.github.com/graphql, exit 1
- **2026-11-23 09:46 · Alexander Sacharov** — hunt report: 2 findings (1 change, blocker; 1 pre-existing), both notes on XSGWPX | read whole: core/pr (pr.go, remote.go, github.go, state.go, watch.go, gitlab.go target/pages), board.ts remote/forge helpers, register.tsx forgeRepo/pullTarget/decideStep/askPullRequest, the uncommitted .jaira/lanes, AGENTS.md, CLAUDE.md (jaira update text, no code); first hunt | not checked: a GHE remote with a port (https://ghe:8443/...) loses the port in --hostname and host/owner/repo - no GHE instance with a port to run against; the cockpit's gh pr list --repo ssh.github.com/... was not run, gh's -R parsing normalises it the way repo view does
- **2026-11-23 09:47 · Alexander Sacharov** — bughunt: done at 53d1eb583 - 2 findings (1 change, 1 pre-existing) in 3 modules, first hunt; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/tests/board.test.ts core/pr/pr.go core/pr/remote.go core/pr/remote_test.go core/release/notes/XSGWPX.md core/shipped/shipped.sha256
- **2026-11-23 09:47 · Alexander Sacharov** — testing: class: a host gh maps to github.com is taken for GHE; example: ssh://git@ssh.github.com:443/me/repo.git -> --hostname ssh.github.com (bughunt note); fix: github.com or any *.github.com stays as before (raw origin, no --hostname), in Go ghRepo and TS ghRepoOf, test over the class (ssh.github.com, www.github.com, GitHub.com). pre-existing reap.sh/pr-round.sh: same subject (GHE), left as note for the person per 'no satellite tickets'.
- **2026-11-23 09:49 · Alexander Sacharov** — fix round 2: *.github.com (ssh.github.com, www.github.com) bleibt github.com — Go ghRepo (remote.go) und TS ghRepoOf; Tests: TestReadGitHubEnterprise mit ssh.github.com/www.GitHub.com ohne --hostname, describe('ghRepoOf'). note: limit: GHE mit Port (https://ghe:8443/...) verliert den Port in --hostname — gh-Hosts stehen in hosts.yml ohne Port üblich, keine Instanz zum Prüfen.
- **2026-11-23 09:49 · Alexander Sacharov** — critique round 2: none — liest nur den Fix df630383f: der bughunt-Befund (*.github.com) ist in Go und TS an derselben Stelle behoben, Tests über die Klasse.
- **2026-11-23 09:50 · Alexander Sacharov** — hunt report: 0 findings, nothing written beyond this report | read as a diff since 53d1eb583: core/pr/remote.go, core/pr/remote_test.go, core/mod/builtin/jaira-cockpit/hooks/board.ts, tests/board.test.ts, core/shipped/shipped.sha256, .jaira/lanes/*.md, AGENTS.md, CLAUDE.md (regenerated board text, no code); no new files | ran: go test ./core/pr ./core/shipped (ok), ghRepoOf via node --experimental-strip-types on ssh.github.com:443, www.GitHub.com, GITHUB.COM, ghe host, notgithub.com (all as intended); could not run the cockpit TS suite itself (no runner found in the worktree), and did not check live gh against a real GHE host
- **2026-11-23 09:50 · Alexander Sacharov** — bughunt: done at df630383f - 0 findings in 3 modules, follow-up hunt since 53d1eb583; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/tests/board.test.ts core/pr/pr.go core/pr/remote.go core/pr/remote_test.go core/release/notes/XSGWPX.md core/shipped/shipped.sha256
- **2026-11-23 09:51 · Alexander Sacharov** — waiting: fresh review lane (run-lane.sh, Herdr tab) since 10:51
- **2026-11-23 09:55 · Alexander Sacharov** — review: done at df630383f
- **2026-11-23 09:56 · Alexander Sacharov** — review: move to signoff refused by the bughunt gate - the uncommitted jaira update regeneration (AGENTS.md, CLAUDE.md, .jaira/lanes/{critique,in-progress,review,testing}.md, no code) counts as unhunted change. Question: how do these files leave the worktree?
1. (recommended) teamlead/person moves the ticket with --force: the files carry no code, the review read them
2. commit them on this branch as chore(XSGWPX): jaira update regeneration - then re-hunt is needed for the gate
3. git restore them here and let master take the update separately - throws away the regeneration in this worktree
- **2026-11-23 09:56 · Alexander Sacharov** — question (worker): review accept, signoff-Gate verweigert wegen uncommitteter jaira-update-Dateien (AGENTS.md, CLAUDE.md, .jaira/lanes/*.md — 'brought this board up to 0.6.4', vom jaira-Binary beim ersten Aufruf geschrieben, nicht von mir). Optionen: 1) restore: git checkout -- die 6 Dateien, Review-Lane neu (Recommended — Board-Update gehört nicht in den Ticket-PR) 2) eigener Commit chore(XSGWPX) auf diesem Branch 3) --force am Gate. waiting: Antwort des Parents seit 10:56
- **2026-11-23 09:57 · Alexander Sacharov** — decision (teamlead 2026-11-23): option 1 — 6 jaira-update-Dateien per git checkout verworfen, Board-Update geht über eigenen PR nach master; Review-accept bleibt (Code unverändert), direkt nach signoff.
- **2026-11-23 09:57 · Alexander Sacharov** — report: GHE-Remote user@host:owner/repo -> gh mit host/owner/repo + --hostname in jaira pr/State und Cockpit-PR-Abfrage; *.github.com unverändert | signoff, context unknown | offen für die Person: 4 notes (Port-limit, SSH-Alias-limit, reap.sh/pr-round.sh nur github.com, Regen-Dateien erledigt); Branch noch nicht gepusht
- **2026-11-23 11:59 · Alexander Sacharov** — Accepted from signoff: angenommen:
- note: limit: a GHE remote with a port (https://ghe:8443/...) loses the port in --hostname and host/owner/repo - stated in the ticket notes, no instance to check against
- note: limit: an SSH host alias as remote (git@github.com-work:me/repo.git) is asked as host github.com-work - it failed the same way before (GH_DEBUG: gh repo view on the raw URL requests https://github.com-work/api/graphql), no regression
- note: elsewhere: core/role/builtin/jaira-dispatcher/scripts/reap.sh:175/:224 and jaira-role-pr/scripts/pr-round.sh:81 still talk to github.com only on a GHE remote (pre-existing, bughunt note)
- note: elsewhere: the uncommitted AGENTS.md, CLAUDE.md and .jaira/lanes/*.md are a jaira update regeneration to 0.6.3, not this ticket - keep them out of the XSGWPX commits
- **2026-11-23 11:59 · Alexander Sacharov** — cockpit: accepted by the person, signoff -> done
