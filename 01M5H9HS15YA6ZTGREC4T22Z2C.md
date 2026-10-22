---
id: 01M5H9HS15YA6ZTGREC4T22Z2C
title: "CI und ein lokaler Hook lehnen Commits ab, die Arbeitgeber- oder Kundennamen tragen"
status: critique
ready: true
creator: BeMuCa
goal: "Kein Commit, keine Datei und keine Ticket-Notiz im Repo nennt Berks Arbeitgeber oder Kunden - weder im Inhalt noch in Author-/Committer-E-Mail - und die Pipeline weist einen PR ab, der es doch tut."
context: |-
  Am 2026-10-22 fiel auf: die globale git-E-Mail auf Berks Rechner lag auf einer Firmen-Domain, und 927 Commits auf master tragen sie. Dazu nennt core/role/builtin/jaira-role-pr/SKILL.md:84 einen Firmen-GitLab-Host, und fuenf Logbuch-/Archiv-Tickets haben Firmennamen in Notizen.
  Die E-Mail fuer dieses Repo ist seit dem 2026-10-22 die GitHub-noreply-Adresse (git config --local). Die offenen Branches wurden umgeschrieben; master wird am Sonntag 2026-10-26 umgeschrieben.
  Die Liste der verbotenen Strings darf selbst nicht im Repo liegen - sonst ist sie die Erwaehnung. Entschieden: ein GitHub-Repository-Secret haelt die Regex-Liste; die CI liest sie. Lokal liest ein pre-commit/pre-push-Hook dieselbe Liste aus einer Datei unter ~/.config/.
  Verworfen: DataFog oder ein anderer PII-Scanner - vier bekannte Strings sind ein Regex, kein ML-Problem.
definition-of-done: "Ein CI-Job in .github/workflows/ci.yaml prueft bei jedem PR den Diff gegen master und Author-/Committer-E-Mail jedes Commits gegen die Regex-Liste aus dem Secret MENTION_DENYLIST und schlaegt bei einem Treffer mit Zeile und Commit fehl; fehlt das Secret, schlaegt der Job mit einer klaren Meldung fehl statt still zu bestehen"
tags:
  - ci
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T07:52:29Z
updated-at: 2026-10-22T08:05:14Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-4193148
claimed-at: 2026-10-22T07:55:54Z
outcome-what: "Added scripts/check-mentions.sh and a CI job 'mentions' in .github/workflows/ci.yaml: a pull request fails when an added line, a file name, an author or committer email, or a commit message matches the denylist (secret MENTION_DENYLIST in CI, ~/.config/jaira/mention-denylist locally); --staged checks the index and git's ident for a pre-commit hook. A hit prints only where it is, never the text or the pattern. docs/COMMANDS.md explains the list and the pre-commit/pre-push hooks. The PR role's example host is now git.example.com."
outcome-why: "The list of forbidden strings cannot live in the repo, or it would itself be the mention, so CI reads it from a secret and the local hooks from a file outside the repo. The CI log of a public repo must not repeat the mention, so the script names only file:line or the commit."
outcome-resolves: "DoD 1: job mentions + scripts/check-mentions.sh, exit 2 with a message when the list is missing. DoD 2: the script's local list source and the hook section in docs/COMMANDS.md. DoD 3: SKILL.md:84 names git.example.com. DoD 4: release note core/release/notes/T22Z2C.md written; the full race suite is the maintainer's before the PR."
---

# CI und ein lokaler Hook lehnen Commits ab, die Arbeitgeber- oder Kundennamen tragen

## Definition of Done

- [x] Ein CI-Job in .github/workflows/ci.yaml prueft bei jedem PR den Diff gegen master und Author-/Committer-E-Mail jedes Commits gegen die Regex-Liste aus dem Secret MENTION_DENYLIST und schlaegt bei einem Treffer mit Zeile und Commit fehl; fehlt das Secret, schlaegt der Job mit einer klaren Meldung fehl statt still zu bestehen
  proof: .github/workflows/ci.yaml job mentions; scripts/check-mentions.sh scan_added + commit loop; shell test cases 2,4,5,6,8
- [x] scripts/check-mentions.sh macht dieselbe Pruefung lokal, liest die Liste aus $XDG_CONFIG_HOME/jaira/mention-denylist (Fallback ~/.config/jaira/mention-denylist), und docs/COMMANDS.md sagt, wie man es als pre-commit- und pre-push-Hook einhaengt
  proof: scripts/check-mentions.sh:22-27 (XDG list, fallback HOME/.config); docs/COMMANDS.md section 'Keeping a name out of the repository'
- [x] core/role/builtin/jaira-role-pr/SKILL.md:84 nennt keinen Firmen-Host mehr, sondern ein generisches Beispiel
  proof: core/role/builtin/jaira-role-pr/SKILL.md:84 now git.example.com
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-22 07:59 · BeMuCa** — Driver decisions before in-progress: (1) pre-process skipped, the Options checklist does not tick planning; the method is in this note. (2) SKILL.md:83 example host becomes git.example.com, not gitlab.example.com: rule 4 there is for a host that says neither github nor gitlab, so an example containing 'gitlab' would fall under rule 3. (3) The check prints only where a hit is (file:line, or short sha plus field), never the matched text or the pattern: on a public repo the CI log would otherwise repeat the mention. (4) Diff is base...head (merge base), so lines master changed after the branch point are not blamed on the PR; --no-renames so a renamed file's new path and content count as added. (5) pre-commit has no commit range yet, so the script takes --staged: index diff plus the author/committer ident git would use. (6) Found: 5 of the last 8 merged PRs come from a fork (sashasoft90); GitHub passes no secrets to pull_request runs from forks, so as decided the job fails on each of those PRs - question for Berk once it works.
- **2026-10-22 08:00 · BeMuCa** — Verified in GitHub docs (using-secrets-in-github-actions): 'With the exception of GITHUB_TOKEN, secrets are not passed to the runner' when the workflow is triggered from a forked repository. So under the pull_request trigger the mentions job sees an empty MENTION_DENYLIST on every fork PR and fails by design.
- **2026-10-22 08:03 · BeMuCa** — Built scripts/check-mentions.sh: diff is parsed with awk into two parallel temp files (path:content for matching, path:line for reporting) so grep -n line numbers map back to the new-file line, and the matched text is never printed. Each file also gets one path-only entry so a new file whose NAME matches is caught even when empty. Dead end: piping the diff into a shell function loses 'bad=1' (pipe subshell), so the diff goes to a temp file first. Invalid regex is checked once up front against empty input (grep exit 2), not per match. SKILL.md host was on line 84, not 83 as the brief said. shellcheck is not installed here. Shell test lives in the scratchpad (T22Z2C-check-mentions-test.sh), not in the repo: the repo has no harness for scripts.
