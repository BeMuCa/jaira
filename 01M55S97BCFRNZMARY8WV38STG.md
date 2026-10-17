---
id: 01M55S97BCFRNZMARY8WV38STG
title: "Die Installation per curl richtet Rollen und Cockpit gleich mit ein und zeigt dabei, was sie tut"
status: testing
ready: true
creator: berkc
assignee: berkc
goal: "Wer jaira mit dem curl-Einzeiler (scripts/install.sh) oder install.ps1 installiert, hat danach Binary, Rollen-Skills und Cockpit eingerichtet, ohne einen zweiten Befehl, und sieht dabei Fortschritt: Download-Balken, 'Setting up the skills', 'Setting up the Cockpit'."
context: |-
  Was heute ist: scripts/install.sh (Zeile ~86-94) und scripts/install.ps1 hoeren auf, sobald das Binary liegt. Rollen und Cockpit kommen erst mit 'jaira roles install --global', und das sagt niemandem jemand.
  Ausloeser: Berk am 17.10.: hat per curl installiert, im Repo 'jaira' getippt, Board war da - 'jaira init' brauchte er nie. Rollen und Cockpit fehlten aber. Er will KEIN Onboarding-Dialog, sondern: der curl-Install installiert und aktiviert Rollen und Cockpit gleich mit.
  Bekannt: 'jaira roles install --global' schreibt 34 Dateien nach ~/.claude/skills, inklusive Cockpit (am 17.10. mit 0.5.2 laufen lassen: zweiter Lauf '0 written, 34 unchanged', exit 0). Das Cockpit laedt erst in der naechsten Claude-Code-Session.
  Gewuenscht sind sichtbare Schritte: Download mit Fortschrittsbalken (curl --progress-bar o.ae.), dann Zeilen wie 'Setting up the skills', 'Setting up the Cockpit', am Ende was jetzt aktiv ist.
  Offen: ein Opt-out (Umgebungsvariable, z.B. JAIRA_NO_ROLES=1) fuer CI und Leute ohne Claude Code; nicht-interaktiv lauffaehig bleiben (curl | sh hat kein TTY auf stdin).
definition-of-done: |-
  install.sh installiert nach dem Binary die Rollen und das Cockpit global, mit Fortschrittsbalken beim Download und je einer Statuszeile fuer Skills und Cockpit
  install.ps1 tut dasselbe
  Ein Opt-out per Umgebungsvariable laesst nur das Binary installieren
  Ein zweiter Lauf ist idempotent und meldet unveraenderte Dateien statt Fehler
  Release-Note in core/release/notes/<ID>.md
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-17T20:36:33Z
updated-at: 2026-10-17T23:51:53Z
updated-by: berkc
claimed-by: berk-ca-283567
claimed-at: 2026-10-17T20:37:26Z
outcome-what: "install scripts run roles install --global; ps1 failure branch fixed, opt-out only =1"
outcome-why: curl install left roles and Cockpit missing
outcome-resolves: "install.sh verified; install.ps1 not run"
review-summary: "install.sh and install.ps1 now run 'jaira roles install --global' after placing the binary, as ONE step 'Setting up the skills and the Cockpit…' (roles install has no flag to split them - checked with --help), print the 'N written, M unchanged' line, what is active, and that the Cockpit loads next Claude Code session. Failure is printed with the command's own output and does not undo the binary install. JAIRA_NO_ROLES=1 (only exactly 1) skips it. install.sh shows a curl --progress-bar for the archive only when stderr is a TTY; install.ps1 turns the IWR progress bar on for the archive when [Environment]::UserInteractive. ps1 switches ErrorActionPreference to Continue around the native call and back to Stop (function-local, cannot leak into the caller session under irm|iex). README and release note V38STG.md updated."
review-gaps: "1) install.ps1:112 re-enables the Invoke-WebRequest progress bar for the 5.3MB archive, which the file's own comment at install.ps1:84-85 says makes a PS 5.1 download crawl; UserInteractive is true for every irm|iex user, so all Windows PowerShell 5.1 installs take this path - speed unmeasured (no pwsh here). 2) If a user has edited an installed role, roles install exits 3 (34 unchanged, 1 skipped): install.sh:117-121 / install.ps1:160-163 then print 'could not set up' and 'retry with: jaira roles install --global', which fails the same way again; jaira's own line does say --force. 3) PS 5.1 '2>&1 | Out-String' renders stderr as ErrorRecords ('jaira : ... + CategoryInfo ... NativeCommandError') in the failure text - cosmetic, unverified. 4) Ticket body still has four stray untagged DoD lines under the 5 parsed items (leftover of round 1); parsed DoD is correct. 5) DoD 1 asked for one status line each for skills and Cockpit; one combined line shipped - justified, roles install --help shows no flag to split them. install.ps1 entirely unrun."
test-verdict: "pass: sh -n, dash runs: fresh 35 written, rerun 0 written/35 unchanged, NO_ROLES=1 skips, =0 installs, read-only failure printed. install.ps1 NOT run (no pwsh)."
review-verdict: "Approve with gaps. install.sh verified under dash against the real v0.5.3 release: fresh 35 written exit 0, rerun 0 written/35 unchanged exit 0, JAIRA_NO_ROLES=1 skips and creates no ~/.claude, JAIRA_NO_ROLES=0 installs, read-only ~/.claude prints the permission-denied error plus retry hint and still exits 0 with the binary in place. Round-1 findings are fixed in the diff (EAP Continue around the native call, =1 only). install.ps1 judged by reading only - unsure about the 5.1 progress-bar slowdown (gap 1); a person should try it on Windows before accepting."
review-round: 2
review-check: "1. cd into the jaira checkout on branch feat/v38stg  2. mkdir -p /tmp/jt/home  3. HOME=/tmp/jt/home JAIRA_INSTALL_DIR=/tmp/jt/bin sh scripts/install.sh  4. you see a #### progress bar, then 'Setting up the skills and the Cockpit…' and '✓ 35 written, 0 unchanged, ...'  5. ls /tmp/jt/home/.claude/skills shows jaira-cockpit and jaira-role-* folders  6. run step 3 again: line reads '✓ 0 written, 35 unchanged', exit 0  7. HOME=/tmp/jt/home2 JAIRA_NO_ROLES=1 JAIRA_INSTALL_DIR=/tmp/jt/bin2 sh scripts/install.sh: last line 'skipped the skills and the Cockpit ...', /tmp/jt/home2 has no .claude  8. On Windows PowerShell 5.1: make %USERPROFILE%\\.claude read-only, run scripts/install.ps1 - expected 'FAILED to set up ...' plus retry hint, not a red NativeCommandError  9. Normal Windows run: 'OK 35 written, ...' and jaira-cockpit in %USERPROFILE%\\.claude\\skills"
round: 2
---

# Die Installation per curl richtet Rollen und Cockpit gleich mit ein und zeigt dabei, was sie tut

## Definition of Done

- [x] install.sh installiert nach dem Binary die Rollen und das Cockpit global, mit Fortschrittsbalken beim Download und je einer Statuszeile fuer Skills und Cockpit
  proof: install.sh: one combined status line, not two - roles install has no flag to split skills from Cockpit
- [x] install.ps1 does the same as install.sh
  proof: scripts/install.ps1 Install-Jaira tail; NOT run, no pwsh
- [x] An opt-out env var installs only the binary
  proof: install.sh: JAIRA_NO_ROLES=1 run, no ~/.claude created; =0 installs roles
- [x] A second run is idempotent and reports unchanged files instead of an error
  proof: dash rerun: 0 written, 35 unchanged
- [x] Release note in core/release/notes/<ID>.md
  proof: core/release/notes/V38STG.md
install.ps1 tut dasselbe
Ein Opt-out per Umgebungsvariable laesst nur das Binary installieren
Ein zweiter Lauf ist idempotent und meldet unveraenderte Dateien statt Fehler
Release-Note in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-17 20:37 · berkc** — berkc took this ticket over from BeMuCa
- **2026-10-17 20:38 · berkc** — One roles install call writes skills+cockpit (no separate flag), so one honest step 'Setting up the skills and the Cockpit'. Progress bar: curl only, only when stderr is a TTY. Opt-out JAIRA_NO_ROLES. install.ps1 unverified: no pwsh on this box.
- **2026-10-17 23:29 · berkc** — review (fresh context) sends back: (1) install.ps1:150 under PS 5.1 with ErrorActionPreference=Stop, '& $exe roles install --global 2>&1' throws on jaira's stderr failure line, so the FAILED/retry branch never runs - wrap the call with ErrorActionPreference='Continue' (unverified, no pwsh); (2) DoD items 2-5 are plain body lines, not checkboxes - re-enter and tick with proof; (3) DoD wants a separate status line for skills and Cockpit, script prints one combined line - split or get the DoD amended. install.sh itself verified: dash fresh=35 written, rerun=0 written/35 unchanged, JAIRA_NO_ROLES skip, pty bar, read-only ~/.claude prints error + retry, exit 0.
- **2026-10-17 23:48 · berkc** — review round: ps1 ErrorActionPreference Continue around native call; opt-out only ==1; roles install has no mods flag so one step stays; DoD items re-added.
