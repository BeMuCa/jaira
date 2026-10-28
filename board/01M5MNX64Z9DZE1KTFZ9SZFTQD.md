---
id: 01M5MNX64Z9DZE1KTFZ9SZFTQD
title: "Ergebnisse fuer eine Person liegen auf einem stabilen Pfad, nie unter /tmp"
status: human
ready: true
creator: BeMuCa
goal: "Was ein Worker fuer einen Menschen erzeugt (Vergleichsseite, Bericht, Screenshot), ueberlebt einen Crash: die Rollen sagen, wohin es gehoert, und der Pfad steht am Ticket."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkt 16): die Vergleichsseite eines Workers lag im Scratchpad unter /tmp und war nach dem WSL-Crash weg. Dasselbe passierte dieser Session mit Backup-Bundle und SHA-Map.
  Das Scratchpad ist fuer Zwischenergebnisse richtig; ein Ergebnis fuer eine Person ist keins. Die Rollen (jaira-worker, jaira-role-lane, jaira-role-acceptance) nennen heute keinen Ort dafuer.
definition-of-done: "jaira-worker, jaira-dispatcher, jaira-role-lane und jaira-role-acceptance sagen: ein Ergebnis fuer eine Person wird unter einem stabilen Pfad abgelegt (Board-Setting 'deliverables' oder, wenn leer, <repo>/.jaira/deliverables/<ticket>/), nie unter /tmp, und der Pfad steht als Notiz am Ticket"
tags:
  - docs
blocked-by: []
related: []
commits:
  - a780e3fca9d45d977c7495252e356612ee0ebc1c
created-at: 2026-10-23T15:26:09Z
updated-at: 2026-10-27T20:53:36Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: ""
claimed-at: ""
outcome-what: "Vier Rollen-Skills (jaira-worker, jaira-dispatcher, jaira-role-lane, jaira-role-acceptance) bekommen je einen Absatz: ein Ergebnis fuer eine Person (Vergleichsseite, Bericht, Screenshot, Abnahmeseite) liegt unter <repo>/.jaira/deliverables/<ticket-handle>/, nie unter /tmp oder im Scratchpad, und der Pfad kommt per jaira note ans Ticket. Release-Notiz core/release/notes/SZFTQD.md."
outcome-why: "Eine Vergleichsseite eines Workers lag im Scratchpad unter /tmp und war nach dem WSL-Crash weg (Feedback Punkt 16). Die Rollen nannten keinen Ort fuer Ergebnisse, die fuer einen Menschen bestimmt sind."
outcome-resolves: "DoD 1: alle vier Skills nennen den stabilen Pfad, verbieten /tmp und verlangen die Notiz am Ticket (Fundstellen im Proof). Ein Setting 'deliverables' gibt es nicht und wurde nicht gebaut (nicht text-only), es gilt der Fallback. DoD 2: core/release/notes/SZFTQD.md:1."
review-summary: |-
  core/role/builtin/jaira-worker/SKILL.md:331, jaira-dispatcher/SKILL.md:474, jaira-role-lane/SKILL.md:342, jaira-role-acceptance/SKILL.md:52 — "the board setting `deliverables` if the board has one, otherwise" names a setting that does not exist and that nothing reads (driver note 2026-10-27 00:03 (a)); the brief allows naming it only if it exists or is trivial to add. Drop the clause and name only the fallback path.
  same four lines — "an acceptance page" in the example list contradicts jaira-role-acceptance/SKILL.md:45, which keeps that page under ~/.claude/acceptance/<repo>/<name>/; the ticket goal lists only comparison page, report, screenshot. Drop "an acceptance page" in all four; in jaira-role-acceptance:52 the closing sentence then only has to say the page stays where the listing puts it — and that listing is above, not "as below".
  core/role/builtin/jaira-role-acceptance/SKILL.md:52 — the paragraph sits between the listing (lines 45-50) and the line that defines `<repo>` for it (line 54, "the basename of the main worktree"), so that definition now also reads onto `<repo>/.jaira/deliverables/`, where `<repo>` has to be a path. Move the paragraph below the "If `url.txt` exists" paragraph (after line 60) and do not give `<repo>` a second meaning in this file.
  same four lines — each paragraph is one line of 442-581 characters; every other prose line in these files wraps at about 80 columns (only frontmatter and code-block lines exceed 100). Wrap at the file's width; in jaira-role-lane:342 indent the continuation lines two spaces like the bullets above it.
  DECISION FOR THE MAINTAINER, same four lines — the path puts results for a person inside the repository, beside an existing pattern that keeps them out: jaira-role-acceptance/SKILL.md:45 (~/.claude/acceptance/<repo>/<name>/) and jaira-role-logbook-summary/SKILL.md:103-104 ("beside nothing in the repository — in ~/.claude/summaries/<repo>/<from>_<to>/"). `<repo>` is undefined in worker/dispatcher/lane, so an agent in its ticket worktree (jaira-worker/SKILL.md:299) writes there; on a shared board .jaira/deliverables/ is not ignored (git check-ignore exits 1 in this repo), and core/gitrepo/git.go:229 lists every untracked file outside .jaira/tickets and diffs it in full into `worktree_diff`, so a comparison page reaches every later lane's payload as part of the change, `git add -A` commits it, and removing the worktree after done deletes it while the ticket note still points at it. Writing into the main worktree instead trips "Your worktree only" (jaira-role-lane/SKILL.md:336) and the stop rules (jaira-worker/SKILL.md:322, jaira-dispatcher/SKILL.md:460). Options: (a) keep the path and add /.jaira/deliverables/ to the gitignore beside /.jaira/lanes/ (core/board/gitignore.go:19) — code change, breaks text-only, still lost with the worktree; (b) move it outside the repository like acceptance and summaries, e.g. ~/.claude/deliverables/<repo>/<ticket-handle>/ — text-only, out of git and of every diff, but not the path the DoD names; (c) keep it as is and accept the payload and git-add cost.
question: "Wo liegen Ergebnisse fuer eine Person (Vergleichsseite, Bericht, Screenshot)? (a) .jaira/deliverables/<handle>/ wie im DoD, plus /.jaira/deliverables/ in die gitignore neben /.jaira/lanes/ (core/board/gitignore.go:19). Kostet: Code-Aenderung, nicht mehr text-only; liegt im Ticket-Worktree und ist mit dem Worktree nach done trotzdem weg. (b) ausserhalb des Repos wie Abnahme und Summary, z. B. ~/.claude/deliverables/<repo>/<handle>/. Kostet: weicht vom Pfad im DoD ab (DoD-Text muesste angepasst werden); dafuer text-only, nie in git, in keinem Diff, ueberlebt Worktree-Entfernung. (c) so lassen. Kostet: auf einem geteilten Board landet jede Vergleichsseite im worktree_diff aller spaeteren Lanes (core/gitrepo/git.go:229), git add -A committet sie, nach done ist sie mit dem Worktree weg. Danach zurueck nach in-progress mit den Funden 1-4 aus review-summary."
---

# Ergebnisse fuer eine Person liegen auf einem stabilen Pfad, nie unter /tmp

## Definition of Done

- [x] jaira-worker, jaira-dispatcher, jaira-role-lane und jaira-role-acceptance sagen: ein Ergebnis fuer eine Person wird unter einem stabilen Pfad abgelegt (Board-Setting 'deliverables' oder, wenn leer, <repo>/.jaira/deliverables/<ticket>/), nie unter /tmp, und der Pfad steht als Notiz am Ticket
  proof: core/role/builtin/jaira-worker/SKILL.md:331, jaira-dispatcher/SKILL.md:474, jaira-role-lane/SKILL.md:342, jaira-role-acceptance/SKILL.md:52
- [x] Release-Notiz unter core/release/notes/
  proof: core/release/notes/SZFTQD.md:1

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-27 00:03 · BeMuCa** — Driver: brainstorm und pre-process uebersprungen (Optionen nicht gesetzt, next_lane lt. lane.OnPath). in-progress an Subagent (tier cheap) gegeben; Vorgabe: kein neues deliverables-Setting, nur Fallback .jaira/deliverables/<handle>/, ausser es ist trivial.
- **2026-10-27 00:03 · BeMuCa** — (a) No 'deliverables' setting exists: grep -rni deliverables over *.go finds nothing; settings are the machine JSON in core/settings/settings.go plus single files per board. The skills say 'the board setting deliverables if the board has one' but nothing reads it, so the fallback .jaira/deliverables/<handle>/ is what actually applies. Not added (not text-only). (b) .jaira/deliverables/ is safe: .jaira subdirs are read by name only (boards/ in core/lane/boards.go:60, tickets/logbook in core/ticket/store.go:519); board discovery core/project/discover.go:45 checks only .jaira/tickets. No scan of .jaira itself. (c) 'jaira init' writes /.jaira/ to .gitignore (core/board/gitignore.go:10,44-56); 'jaira share' removes it so deliverables are then untracked-but-visible and get committed if someone git adds them; only /.jaira/lanes/ stays ignored (gitignore.go:19). Large files there would reach the repo.
- **2026-10-27 00:10 · BeMuCa** — critique (pass 1): core/role/builtin/jaira-worker/SKILL.md:331, jaira-dispatcher/SKILL.md:474, jaira-role-lane/SKILL.md:342, jaira-role-acceptance/SKILL.md:52 — "the board setting `deliverables` if the board has one, otherwise" names a setting that does not exist and that nothing reads (driver note 2026-10-27 00:03 (a)); the brief allows naming it only if it exists or is trivial to add. Drop the clause and name only the fallback path.
same four lines — "an acceptance page" in the example list contradicts jaira-role-acceptance/SKILL.md:45, which keeps that page under ~/.claude/acceptance/<repo>/<name>/; the ticket goal lists only comparison page, report, screenshot. Drop "an acceptance page" in all four; in jaira-role-acceptance:52 the closing sentence then only has to say the page stays where the listing puts it — and that listing is above, not "as below".
core/role/builtin/jaira-role-acceptance/SKILL.md:52 — the paragraph sits between the listing (lines 45-50) and the line that defines `<repo>` for it (line 54, "the basename of the main worktree"), so that definition now also reads onto `<repo>/.jaira/deliverables/`, where `<repo>` has to be a path. Move the paragraph below the "If `url.txt` exists" paragraph (after line 60) and do not give `<repo>` a second meaning in this file.
same four lines — each paragraph is one line of 442-581 characters; every other prose line in these files wraps at about 80 columns (only frontmatter and code-block lines exceed 100). Wrap at the file's width; in jaira-role-lane:342 indent the continuation lines two spaces like the bullets above it.
DECISION FOR THE MAINTAINER, same four lines — the path puts results for a person inside the repository, beside an existing pattern that keeps them out: jaira-role-acceptance/SKILL.md:45 (~/.claude/acceptance/<repo>/<name>/) and jaira-role-logbook-summary/SKILL.md:103-104 ("beside nothing in the repository — in ~/.claude/summaries/<repo>/<from>_<to>/"). `<repo>` is undefined in worker/dispatcher/lane, so an agent in its ticket worktree (jaira-worker/SKILL.md:299) writes there; on a shared board .jaira/deliverables/ is not ignored (git check-ignore exits 1 in this repo), and core/gitrepo/git.go:229 lists every untracked file outside .jaira/tickets and diffs it in full into `worktree_diff`, so a comparison page reaches every later lane's payload as part of the change, `git add -A` commits it, and removing the worktree after done deletes it while the ticket note still points at it. Writing into the main worktree instead trips "Your worktree only" (jaira-role-lane/SKILL.md:336) and the stop rules (jaira-worker/SKILL.md:322, jaira-dispatcher/SKILL.md:460). Options: (a) keep the path and add /.jaira/deliverables/ to the gitignore beside /.jaira/lanes/ (core/board/gitignore.go:19) — code change, breaks text-only, still lost with the worktree; (b) move it outside the repository like acceptance and summaries, e.g. ~/.claude/deliverables/<repo>/<ticket-handle>/ — text-only, out of git and of every diff, but not the path the DoD names; (c) keep it as is and accept the payload and git-add cost.
Left alone: four verbatim copies (the brief asks for one paragraph per skill); the release note's 'Run jaira roles install again' matches core/release/notes/0.5.3/YD4BE9.md:1; placement in worker/dispatcher (end of 'Do not swallow what the human should hear') and in role-lane (Boundaries) fits. Findings 1-4 are clear fixes; finding 5 needs the maintainer and changes the same text, so fix all after the answer.
- **2026-10-27 00:20 · BeMuCa** — Driver: an human gegeben, weil critique (pass 1) eine echte Entscheidung fand (Ablageort, siehe question). Funde 1-4 sind klar und warten auf in-progress, damit die vier Absaetze nur einmal angefasst werden. Ticket-Datei bleibt untracked im Worktree; der naechste Code-Commit nimmt sie mit. Hinweis: mehrere Schreibvorgaenge meldeten 'was not sent: BeMuCa wrote it first' - der Ref hing hinter der Datei.
- **2026-10-27 20:53 · BeMuCa** — Aus O6 (27.10.2026): nicht nur Ergebnisse fuer eine Person, auch Session-Zustand lag unter /tmp und ging mit dem WSL-Absturz verloren - Backup-Bundle, SHA-Map der History-Umschreibung, review-check-Dateien. Welcher stabile Pfad auch immer gewaehlt wird: er muss auch fuer solche Zwischenstaende eines Drivers gelten, nicht nur fuer Deliverables.
