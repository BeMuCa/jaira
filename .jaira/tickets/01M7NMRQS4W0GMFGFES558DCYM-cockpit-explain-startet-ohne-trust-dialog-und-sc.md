---
id: 01M7NMRQS4W0GMFGFES558DCYM
title: Cockpit-Explain startet ohne Trust-Dialog und schreibt im richtigen Worktree
status: critique
ready: true
creator: Alexander Sacharov
goal: "Ein Druck auf Explain im Cockpit laeuft ohne Eingriff durch: der Tab startet claude ohne Rueckfrage, der Chat schreibt seine explain:-Note, und der Tab schliesst sich danach von selbst."
context: |-
  Was heute kaputt ist (live 17.11., Herdr App, Claude Code 2.1.326, ein Ticket in einem anderen Repo):
  1. spawn.sh startet Explain/Discuss mit 'claude --add-dir /tmp/jaira-chat.XXXX'. Claude Code 2.1.326 fragt dafuer 'Adding a new directory ... trust?' - der Tab bleibt im Dialog stehen, nichts wird getippt.
  2. Nach Bestaetigung von Hand lief der Chat im Herdr-Review-Worktree (Branch review/feat/<id>). Dort liegt das Ticket nicht auf der Platte; 'jaira note <id>' scheitert mit 'on its ref and not on your disk; pull it first'. Ohne explain:-Note schliesst das Cockpit den Tab nie.
  Ursache 2: worktreeFor() in core/mod/builtin/jaira-cockpit/hooks/board.ts nimmt den ersten Worktree, dessen Branch auf /<id> endet - review/feat/<id> matcht genauso wie feat/<id>, wo das Ticket liegt.
  Geprueft: 'claude -p --allowedTools=Read(//tmp/jaira-chat.X/**)' liest die Datei ohne Rueckfrage; ohne Regel verweigert. --allowedTools ist variadisch: nur in der Form mit '=' schreiben.
definition-of-done: spawn.sh startet Explain/Discuss ohne --add-dir; das Lesen der Auftragsdatei erlaubt eine Read-Regel auf genau ihren Ordner (auch auf Windows)
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-11-17T20:56:53Z
updated-at: 2026-11-17T21:04:37Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-27759
claimed-at: 2026-11-17T20:58:46Z
outcome-what: "Explain/Discuss starten claude mit einer Read-Regel auf den Auftragsordner statt --add-dir; das Cockpit waehlt den Worktree des Ticket-Branches statt Herdrs review/feat/<id>"
outcome-why: "Claude Code 2.1.326 haelt den Tab am Trust-Dialog fuer --add-dir an; im Review-Worktree scheitert jaira note, die explain:-Note fehlt und der Tab schliesst nie"
outcome-resolves: DoD 1-4
review-summary: none
review-gaps: "none - considered cygpath -u for the rule path instead of the drive regex: Cygwin would give /cygdrive/c and the test stub answers only -w, so kept; no other Windows-path conversion in core/role/builtin/jaira-dispatcher/scripts; nothing orphaned (chat_dir_win still read by the regex, --add-dir gone from spawn.sh and tests)"
test-verdict: "fail: go test ./... red - internal/wintrap TestRepositoryIsClean flags core/role/role_test.go:736 (new Read rule suffix glued with literal /)"
round: 2
---

# Cockpit-Explain startet ohne Trust-Dialog und schreibt im richtigen Worktree

## Definition of Done

- [x] spawn.sh startet Explain/Discuss ohne --add-dir; das Lesen der Auftragsdatei erlaubt eine Read-Regel auf genau ihren Ordner (auch auf Windows)
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:312 chat_rule + :321/:340 --allowedTools; TestSpawnChatsStartOnTheChatModelWithTheirBrief, TestSpawnStartsClaudeWhereItsShellRuns/git bash chat
- [x] worktreeFor/branchesFor waehlen bei mehreren passenden Worktrees den Branch des Tickets (feat/<id>), nicht review/feat/<id>
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts ticketsTree(); Fixture WORKTREES in tests/board.test.ts mit review/feat/mg677r zuerst: alter Code 14 Fehler, neuer keiner
- [x] Tests in core/role und im Cockpit decken beides ab; go test ./core/... und claude plugin test gruen bis auf vorbestehende Fehler
  proof: go test ./... RC=0 (incl. internal/wintrap, role_test.go:736 marked //wintrap:ok); claude plugin test . 176/177, der eine Fehler (borderStyle round) auch auf master
- [x] Release-Notiz core/release/notes/<ID>.md
  proof: core/release/notes/58DCYM.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] spawn.sh: --add-dir durch '--allowedTools=Read(//<ordner>/**)' ersetzen, Windows-Pfad in POSIX-Form //c/...
- [x] board.ts: unter passenden Worktrees den Branch mit den wenigsten Teilen waehlen (worktreeFor, branchesFor)
- [x] Tests: role_test.go Startzeilen, Cockpit-Fixture mit review/feat/mg677r vor feat/mg677r

## Progress
- **2026-11-17 20:58 · Alexander Sacharov** — Interaktiver Start im Herdr-Tab mit --allowedTools nicht live geprueft, nur claude -p (liest ohne Rueckfrage, ohne Regel verweigert). Windows-Regelform //c/... nach Claude-Code-Doku zu Permission-Pfaden, nicht auf echtem Windows getestet. Fehler 'borderStyle round' in board.test.ts besteht auf master unabhaengig davon.
- **2026-11-17 20:58 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/51013.sock
- **2026-11-17 20:59 · Alexander Sacharov** — critique: nothing raised. board.ts ticketsTree (fewest branch parts) is general rather than Herdr-specific - a fix/<id> or feat/<id> branch beats review/feat/<id>, a detached tree matched by path loses to any branch; spawn.sh keeps one rule variable for both start lines, same shape as the --add-dir it replaces.
- **2026-11-17 21:02 · Alexander Sacharov** — testing: go test ./... RC=1: internal/wintrap TestRepositoryIsClean: 'core/role/role_test.go:736: windows trap 5: a path is glued together with a literal "/" and handed to strings.HasSuffix'. The string is a Claude Code permission rule in a shell line, not a filesystem path - mark the site //wintrap:ok with that reason. go test ./core/... green; cockpit 176/177, the one fail (borderStyle round) pre-existing per earlier note.
