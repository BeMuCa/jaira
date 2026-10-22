---
id: 01M5H9HS15YA6ZTGREC4T22Z2C
title: "CI und ein lokaler Hook lehnen Commits ab, die Arbeitgeber- oder Kundennamen tragen"
status: testing
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
updated-at: 2026-10-22T08:25:10Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-4193148
claimed-at: 2026-10-22T07:55:54Z
outcome-what: "Added scripts/check-mentions.sh and a CI job 'mentions' in .github/workflows/ci.yaml: a pull request fails when an added line, a file name, an author or committer email, or a commit message matches the denylist (secret MENTION_DENYLIST in CI, ~/.config/jaira/mention-denylist locally); --staged checks the index and git's ident for a pre-commit hook. A hit prints only where it is, never the text or the pattern. docs/COMMANDS.md explains the list and the pre-commit/pre-push hooks. The PR role's example host is now git.example.com."
outcome-why: "The list of forbidden strings cannot live in the repo, or it would itself be the mention, so CI reads it from a secret and the local hooks from a file outside the repo. The CI log of a public repo must not repeat the mention, so the script names only file:line or the commit."
outcome-resolves: "DoD 1: job mentions + scripts/check-mentions.sh, exit 2 with a message when the list is missing. DoD 2: the script's local list source and the hook section in docs/COMMANDS.md. DoD 3: SKILL.md:84 names git.example.com. DoD 4: release note core/release/notes/T22Z2C.md written; the full race suite is the maintainer's before the PR."
review-summary: none
round: 2
review-gaps: "Removed: the $diff temp file and its trap entry - the diff now pipes into added_lines (only awk in the pipe, so no shell state is lost; pipefail keeps a failing git diff fatal), and scan_added lost its argument. Folded the two copies of the denylist normalizer (tr -d CR, drop blank lines) into one filter after if/elif/fi. Rewrapped the ci.yaml mentions comment (first line was 102 chars). Left: the base/head/usage/bad/footer idiom shared with scripts/check-ticket-commits.sh - two standalone CI scripts, a sourced helper for six lines would be a new abstraction; the mentions job's own checkout next to ticket-commits - a separate job is its own status check, merging them changes what CI reports; the --staged and range branches repeating the two git diff calls - unifying changes the printed 'file N of' label for --staged or needs a label parameter; the w_up/tr ident loop - unrolling duplicates the sed; four git calls per commit and sed per hit (critique let stand, PRs small, hits rare); the up-front invalid-regex check - matches() ignores grep's exit 2, so it is not redundant; sub(/\\t$/) in added_lines - live, git appends a tab to a +++ path with a space."
---

# CI und ein lokaler Hook lehnen Commits ab, die Arbeitgeber- oder Kundennamen tragen

## Definition of Done

- [x] Ein CI-Job in .github/workflows/ci.yaml prueft bei jedem PR den Diff gegen master und Author-/Committer-E-Mail jedes Commits gegen die Regex-Liste aus dem Secret MENTION_DENYLIST und schlaegt bei einem Treffer mit Zeile und Commit fehl; fehlt das Secret, schlaegt der Job mit einer klaren Meldung fehl statt still zu bestehen
  proof: .github/workflows/ci.yaml job mentions; scripts/check-mentions.sh scan_names + scan_added + commit loop; shell test cases 2,4,5,6,8,15-23
- [x] scripts/check-mentions.sh macht dieselbe Pruefung lokal, liest die Liste aus $XDG_CONFIG_HOME/jaira/mention-denylist (Fallback ~/.config/jaira/mention-denylist), und docs/COMMANDS.md sagt, wie man es als pre-commit- und pre-push-Hook einhaengt
  proof: scripts/check-mentions.sh:25-30 (XDG list, fallback HOME/.config); docs/COMMANDS.md section 'Keeping a name out of the repository' (end of file)
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
- **2026-10-22 08:12 · BeMuCa** — critique: scripts/check-mentions.sh added_lines: file names come only from '+++' headers, and git writes none for an empty new file or a binary file - probe: empty acme-corp-empty.txt and binary acme-corp.bin both exit 0 (the 08:03 note says an empty one is caught). Take names from `git -c core.quotePath=false diff --name-only --no-renames --diff-filter=d <range>` (with --cached for --staged) and grep that list; drop the awk path-only entry. quotePath=false also fixes non-ASCII names: pattern 'acmé' misses the file acmé-corp.txt today because the header path is octal-quoted.
scripts/check-mentions.sh added_lines: `print p ":" substr($0, 2) > adds` matches the path again on every added line, so one matching file name prints one hit per line (probe: 3-line acme-corp.txt prints 4 hits), and a '^' anchor never reaches the content. Write only substr($0, 2) to $adds; the name list from the finding above checks each name once.
scripts/check-mentions.sh added_lines: /^\+\+\+ / also fires on an added line whose content starts with '++ ' - probe: real.txt with '++ b/fake' on line 2 and a hit on line 4 reports 'fake:3'. Treat '+++' as a header only between 'diff --git' and the first '@@' (h=1 on /^diff --git /, h=0 on /^@@ /, content rule !h && /^\+/). Also pass --no-prefix to both git diff calls and drop sub(/^b\//): with a user's diff.mnemonicPrefix=true, --staged reports 'i/staged.txt:1' today.
scripts/check-mentions.sh scan_added: for a file-name hit the printed location is the mention itself (probe prints 'acme-corp.txt: ...' and 'acme-corp.txt:1: ...'), against decision (3) of the 07:59 note. When the path itself matches, print its number in the name list instead ('file 3 of git diff --name-only <base>...<head>: its name matches') and use that number in place of the path for its content hits.
scripts/check-mentions.sh commit loop: `git log -1 --format=%B "$c" | matches` under pipefail - grep -q exits on the first match, git log dies of SIGPIPE (141), the pipeline returns 141, '&& hit' is skipped and '|| true' swallows it, so a hit passes silently. Probe: hit on line 1 of a ~150 KB message, rc=0 in 5 of 5 runs (flaky at 100 KB, caught up to 70 KB). Use a here-string, no pipe: `matches <<<"$(git log -1 --format=%B "$c")" && hit ...`; same for %ae, %ce and the git var line.
docs/COMMANDS.md: the new section sits between '## Agent plumbing' and '## Keys', inside the binary's reference, while its first line says it is not part of the binary. Move it to the end of the file, after '## Keys'.
.github/workflows/ci.yaml mentions comment: names what a pull request adds and who authored and committed it; the script also checks commit messages - add 'and its commit messages'.
Let stand: per-commit git calls (same loop shape as scripts/check-ticket-commits.sh, PRs are small); the two parallel temp files for path:line (grep stays the matcher); the rc check in scan_added (a grep I/O error must not pass); base...head and --no-renames (07:59 note item 4); --staged (item 5); hit prints only the location for emails and messages (item 3); fork PRs failing (item 6 and the 08:00 note: the driver's question for the maintainer). Probes ran against a scratch repo with MENTION_DENYLIST=acme-corp, not against this board.
- **2026-10-22 08:15 · BeMuCa** — round 2: names now come from git diff --name-only (catches empty/binary/non-ASCII files), matching names print as 'file N' since the name is the mention; commit fields via here-strings because pipefail + grep -q SIGPIPE hid hits on large messages. Self-check with pattern equal to the docs example (the sample word) hits docs/COMMANDS.md itself - inherent, use a non-matching pattern to check exit 0.
- **2026-10-22 08:24 · BeMuCa** — optimize: dropped the $diff temp file - the 08:03 dead end (pipe loses bad=1) only bites when hit() runs in the pipe; now only awk (added_lines, writes files) sits in the pipe and scan_added runs after it in the main shell, and pipefail still exits on a failing git diff (probe: bad base rc=128). The two denylist sources share one tr/grep normalizer (if/elif/fi | filter). Checked live, not dead: sub(/\\t$/) in added_lines - git appends a tab to a +++ path that contains a space (probed). Behaviour unchanged: 27/27 shell tests.
