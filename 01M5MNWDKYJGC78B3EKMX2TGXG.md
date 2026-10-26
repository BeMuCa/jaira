---
id: 01M5MNWDKYJGC78B3EKMX2TGXG
title: "Ein Worker findet seinen Teamlead nach einem Neustart wieder, und eine Frage an die Person kommt als fertige Zeile"
status: in-progress
ready: true
creator: BeMuCa
goal: "Nach einem Crash oder Neustart des Teamleads erreichen die Nachrichten der Worker ihn weiter, ohne dass jemand den neuen Namen von Hand nachreicht; und eine Frage, die nur die Person per Kommando loesen kann, kommt mit der Zeile, die sie tippen muss."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkte 3 und 4).
  Nach einem Neustart bekommt jede Session einen neuen Namen. Worker behalten das beim Start uebergebene --parent und schicken ihre Berichte ins Leere, bis der Teamlead jedem den neuen Namen nachsendet.
  Zweiter Fall: ein Worker brauchte Zugangsdaten in einem Container; Agenten duerfen das zu Recht nicht kopieren. Die Frage ging an die Person, aber ohne die Zeile, die sie haette tippen koennen - eine Runde verloren. Claude Code kennt '! <befehl>' im Prompt der Person.
  Wo das Elternfeld gesetzt wird: spawn.sh (--parent), jaira-worker/SKILL.md und jaira-dispatcher/SKILL.md (Berichtsweg).
definition-of-done: Der Worker schlaegt die aktuelle Adresse seines Teamleads beim Senden nach (Session-Record oder Board) statt den Startwert einzufrieren; ein Test deckt den Namenswechsel ab
tags:
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:25:44Z
updated-at: 2026-10-26T19:17:51Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-252449
claimed-at: 2026-10-26T19:10:49Z
---

# Ein Worker findet seinen Teamlead nach einem Neustart wieder, und eine Frage an die Person kommt als fertige Zeile

## Definition of Done

- [x] Der Worker schlaegt die aktuelle Adresse seines Teamleads beim Senden nach (Session-Record oder Board) statt den Startwert einzufrieren; ein Test deckt den Namenswechsel ab
  proof: core/role/role_test.go TestTeamleadFollowsARestart
- [ ] jaira-worker/SKILL.md und jaira-dispatcher/SKILL.md sagen: eine Frage, die nur die Person per Kommando loesen kann, enthaelt die fertige '! <befehl>'-Zeile
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] teamlead.sh in core/role/builtin/jaira-dispatcher/scripts/: '--here [repo]' writes uds:$CLAUDE_CODE_MESSAGING_SOCKET to <git common dir>/jaira/teamlead (exit 2 without a live socket); '<parent> [repo]' prints the parent while it is a uds: address whose socket exists, else the recorded address when its socket exists, else the parent unchanged
- [x] Go test in core/role/role_test.go: record A; A gone and B recorded (also from a second worktree) -> B; A alive and B recorded -> A; both gone -> A; --here without a socket -> exit 2
- [x] jaira-teamlead/SKILL.md: run 'teamlead.sh --here' first thing in every session
- [x] jaira-worker/SKILL.md and jaira-dispatcher/SKILL.md: before every SendMessage to the parent run 'teamlead.sh <session>' and send to what it prints; a changed address is noted as the new parent; the 'no other peer' rule names the recorded teamlead as the only replacement
- [x] jaira-worker/SKILL.md and jaira-dispatcher/SKILL.md: a question only the person can resolve by running a command carries the ready-to-type '! <command>' line
- [x] core/release/notes/X2TGXG.md; gofmt/vet/build; go test ./core/role -run <touched tests>; commit code together with the ticket file

## Progress
- **2026-10-26 19:15 · BeMuCa** — design (driver, pre-process skipped by ticket options): the record is one file, <git common dir>/jaira/teamlead, beside spawn.sh's jaira/panes/ - shared by every worktree, machine-local, never in git status. Parent first, record second: a parent whose uds: socket still exists keeps the reports, so a second live teamlead on the same clone never takes another's workers (the rule in both SKILL.md files). Only a gone socket is replaced, and only by the address the teamlead wrote itself on start - not a guess among peers. Rejected: a Go command over RepoStateDir (binary change and new CLI surface for one file); record-first (would route a live teamlead's reports to a second one); spawn.sh writing the record (a dispatcher calling spawn.sh would overwrite it).
