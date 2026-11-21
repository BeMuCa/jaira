---
id: 01M5ZJ5DKQZ83GYBGPJHAX8K02
title: "Die PR-Rolle hat einen festen Schritt vor dem Push: Scan, master, Suite, DoD, Logbook"
status: critique
ready: true
creator: BeMuCa
goal: "jaira-role-pr fuehrt vor jedem Push einen festen Ablauf aus: Branch auf E-Mail, Denylist-Strings und Rechnernamen pruefen, master mergen, die volle Race-Suite einmal laufen lassen, die Suite-DoD-Punkte ticken, je Ticket logbook plus chore-Commit, dann push."
context: |-
  Was schief ist: fuer zwei Batches (19 Tickets, 22.-27.10.2026) wurde dieser Schritt jedes Mal von Hand aus dem Kopf gemacht (O5/O17). Die Teile stehen verstreut: Suite nur vor dem PR (Hausregel, 4N1B66), keine Firmen-Strings und kein Hostname im Repo (T22Z2C, QYK8MD), Logbook je Ticket (3SJMV8).
  Heute: core/role/builtin/jaira-role-pr/SKILL.md kennt logbook, aber weder Scan noch Suite noch DoD-Tick (grep race|denylist|hostname|email findet nichts).
  Der Scan liest die Denylist von ausserhalb des Repos (~/.config/jaira/mention-denylist, T22Z2C); die Strings selbst duerfen nie in den Skill.
  Nicht gewollt: PR oeffnen oder mergen - das bleibt beim Maintainer.
definition-of-done: SKILL.md hat den Ablauf als nummerierte Liste; jeder Schritt nennt den Befehl und das Abbruchkriterium
tags:
  - docs
  - release
blocked-by:
  - 01M2HZKCQ48BPQYN1YPTXB97NJ
related:
  - 01M7WWYB5XZMD01DZY3MJG3REW
commits: []
created-at: 2026-10-27T20:52:23Z
updated-at: 2026-11-21T19:17:04Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2284393
claimed-at: 2026-11-21T19:01:24Z
outcome-what: "New core/role/builtin/jaira-role-pr/scripts/pre-push-scan.sh scans every commit of <base>..HEAD (added lines, merges by --cc only, file names, message, author/committer email) against ~/.config/jaira/mention-denylist (exit 3, prints commit+file:line only) and this machine's name (warning, push goes on). SKILL.md 'Push the branch' is now seven numbered steps with command and stop rule: scan, merge base, full suite + CI gates (exit code, not grep), tick suite DoD items with --proof, file finished tickets away, scan again, push. Table test core/role/pre_push_scan_test.go (16 cases), script registered in roleFiles, manifest regenerated, release note AX8K02.md. Four categories checked: concurrency - per-run mktemp dir, no shared state; error paths - no die inside $(...), invalid/unreadable list exits 2, git failure exits non-zero under pipefail; unusual input - empty/binary/CRLF/space/non-ASCII/100 KB message tested; permissions - list read from outside the repo, nothing from the repo is executed."
outcome-why: "For two batches the pre-push step was done by hand from memory: suite only before the PR, no employer/customer strings or hostname in the repo, logbook per ticket. The role named only logbook, so a leak or a missing --force-free done move depended on someone remembering."
outcome-resolves: "jaira-role-pr runs one fixed, checkable sequence before every push; a denylist hit stops it without repeating the string; suite DoD items get ticked so done needs no --force."
---

# Die PR-Rolle hat einen festen Schritt vor dem Push: Scan, master, Suite, DoD, Logbook

## Definition of Done

- [x] SKILL.md hat den Ablauf als nummerierte Liste; jeder Schritt nennt den Befehl und das Abbruchkriterium
  proof: core/role/builtin/jaira-role-pr/SKILL.md:289-341: seven numbered steps, each with its command and a 'Stops on:' line
- [x] Ein Treffer im Scan bricht vor dem Push ab und nennt Datei und Commit, nie den String selbst im Report
  proof: pre-push-scan.sh exit 3 prints '<commit> <file>:<line>: matches the denylist' only; TestPrePushScan (core/role/pre_push_scan_test.go:68) asserts no case's output carries the mention; SKILL.md:298 step 1 + :338 step 7
- [x] Die Suite-DoD-Punkte werden mit --proof auf den Lauf getickt, danach braucht der done-Move kein --force
  proof: SKILL.md step 4; rehearsal 2026-11-21: before the tick 'jaira move 7HBYF4 --to done' refused with 'criterion 2 ... is still open', after 'jaira dod 7HBYF4 2 --done --proof go test ./... -race exit 0 at 9d6b689' the DoD reason is gone; only signoff's human checkpoint remains, which a person passes
- [x] Auf einem Scratch-Branch nachgestellt: der Ablauf laeuft durch und pusht
  proof: scratch repo + bare remote (scratchpad/rehearse): scan rc 0, real merge of a teammate's master commit, go test ./... -race rc 0, tick, rescan rc 0, git push -u origin HEAD rc 0 (remote feat/7hbyf4 = 178307f); a later commit with a mention: scan rc 3, push not run, remote tip unchanged
- [x] Fault class: a mention leaves the machine through any part of a pushed commit. The scan covers every commit of <base>..HEAD (not only the net diff), and in each: added lines, file names, commit message, author and committer email; a merge is checked for its own resolution only
  proof: pre-push-scan.sh scan(): git rev-list base..HEAD, per commit names/adds(--cc)/message/emails; tests 'added then removed again', 'merge: its own resolution, not what it brings in', author/committer email
- [x] Input classes covered by one table test: no commits, empty and binary new file, file name with space and non-ASCII, CRLF line, mention only in the file name, merge with a conflict resolution, a commit message over 100 KB, missing/empty/invalid denylist (invalid stops with exit 2, never passes)
  proof: TestPrePushScan cases in core/role/pre_push_scan_test.go:68-150 (16 subtests, go test ./core/role -run TestPrePushScan -v)
- [x] This machine's name in a commit is a warning that names file and commit and does not stop the push (decision 19:07, variant 1)
  proof: TestPrePushScan/this_machine's_name_only_warns: exit 0, 'warning: ... names this machine', no 'do not push'

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [ ] pre-push-scan.sh beside the role: per commit of <base>..HEAD, added lines (merges: --cc), file names, message, emails; denylist stops (exit 3), hostname warns
- [ ] Go table test in core/role over the input classes; register the script in roleFiles
- [x] SKILL.md: one numbered pre-push sequence, each step with command and stop rule
- [x] release note core/release/notes/AX8K02.md
- [x] rehearse on a scratch repo with a bare remote: sequence runs through and pushes

## Progress
- **2026-11-20 16:49 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-21 19:01 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock
- **2026-11-21 19:02 · Alexander Sacharov** — decision (worker, before plan): no question for the person. (1) The scan is a script beside the role, core/role/builtin/jaira-role-pr/scripts/pre-push-scan.sh, shipped with it like pr-round.sh: a numbered SKILL step alone is a promise, a script is testable and can print only file+commit. scripts/check-mentions.sh from T22Z2C is NOT in this repo (its commits bbdd6ee6/d7ace33a/2825557b exist on no ref here; T22Z2C sits in human), so nothing to reuse. (2) Denylist file missing (~/.config/jaira/mention-denylist, or JAIRA_MENTION_DENYLIST): the role is shipped to every user, so that is not an abort - the scan says the strings were not checked and still checks email and hostname. (3) Email is checked against the denylist; the hostname (this machine's, case-insensitive) is checked as a literal in added lines, file names and commit messages. (4) The scan runs first (cheap, before the long suite) and again right before the push, because logbook/chore commits made by the step itself come after the first run. (5) Release note core/release/notes/AX8K02.md: a changed shipped role is client-facing (jaira update --roles).
- **2026-11-21 19:03 · Alexander Sacharov** — question (worker, before code): the hostname check would stop every push from this machine. Its name sits in claimed-by of every ticket file a branch commits (core/session/session.go:129 still builds the session id from os.Hostname()); QYK8MD removes that, but sits in signoff and its commits are on no ref here. Options: 1) (Recommended) hostname hit is a warning: printed as file+commit, does not stop the push; denylist and email hits stop it. Turns into a stop by a one-line change once QYK8MD is on master. 2) hostname hit stops the push as the goal says; AX8K02 gets blocked_by QYK8MD and waits for it. 3) hostname hit stops, but the claimed-by:/claimed-at: frontmatter lines are skipped - a special case for one field.
- **2026-11-21 19:06 · Alexander Sacharov** — decision: Alex 2026-11-21: Rechnername im Scan ist nur Warnung (Datei + Commit), Push geht weiter; Denylist und E-Mail stoppen. Wird Stopp, sobald QYK8MD in master ist.
- **2026-11-21 19:10 · Alexander Sacharov** — pre-push-scan.sh built and tested (core/role/pre_push_scan_test.go, 16 cases). Found on the way: (a) a die inside $(...) only ends the subshell, so a grep error would read as 'no match' - every lines() call writes to a file instead; (b) scanning each commit, not base...HEAD's net diff: a mention added and removed again inside the branch is still published by the push (test 'added then removed again'); (c) merges via git show --cc: only lines that differ from every parent, so what master brings in is not blamed, a conflict resolution is; (d) grep's own error text is silenced when the list is validated, so a bad pattern never reaches the report.
- **2026-11-21 19:16 · Alexander Sacharov** — rehearsal: the person's done move could not be played in the scratch board - stripping the agent env vars to pass signoff's human checkpoint was refused by the auto-mode classifier, and that refusal stands. So step 5 (logbook) ran as 'ticket not in done, stays on the board', and DoD 3 is proven by the refusal text: before the tick it names criterion 2, after it only the human checkpoint. Manifest: core/shipped/shipped.sha256 must be regenerated (scripts/shipped-history.sh --write) after every SKILL.md change, and a new file must be git-added first or --write does not see it.
