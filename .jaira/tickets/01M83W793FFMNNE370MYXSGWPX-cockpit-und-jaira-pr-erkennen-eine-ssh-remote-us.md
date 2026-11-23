---
id: 01M83W793FFMNNE370MYXSGWPX
title: "Cockpit und jaira pr erkennen eine SSH-Remote 'user@eigener-host:owner/repo.git' (GitHub Enterprise) und fragen gh mit Host und owner/repo"
status: testing
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
commits: []
created-at: 2026-11-23T09:36:31Z
updated-at: 2026-11-23T09:44:19Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-61389
claimed-at: 2026-11-23T09:37:02Z
outcome-what: "Remote-Parser (Go core/pr/remote.go, TS board.ts parseRemote) zerlegt jede git-Remote-Form in Host + owner/repo; gh bekommt auf einem Nicht-github.com-Host host/owner/repo und --hostname bei gh api — in jaira pr/State und in der Cockpit-PR-Abfrage."
outcome-why: "gh repo view auf 'user@ghe-host:owner/repo.git' scheitert mit 'Could not resolve to a Repository'; gh versteht nur git@-Präfix bzw. [HOST/]OWNER/REPO."
outcome-resolves: Cockpit-Abnahme und jaira pr finden PRs auf GitHub Enterprise mit SSH-Remote eines beliebigen Users.
review-summary: none
review-gaps: "parseRemote/remote in core/pr unexportiert (nur im Paket benutzt), onGitHubCom mit einem Aufrufer in ghRepo eingefaltet. Keine Duplikation: kein anderer Remote-URL-Parser in Go oder im Cockpit (grep '://', ssh://, .git-Suffix). Go/TS-Paar bleibt bewusst (Cockpit läuft ohne Go). onHost (3 Aufrufer) bleibt."
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
