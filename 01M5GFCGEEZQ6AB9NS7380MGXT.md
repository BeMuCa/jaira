---
id: 01M5GFCGEEZQ6AB9NS7380MGXT
title: "Ein Ticket in done raeumt seinen Worktree weg, sobald der Branch gepusht und sauber ist"
status: in-progress
ready: true
creator: berkc
assignee: berkc
goal: "Erreicht ein Ticket done, verschwindet sein Worktree (und Herdr-Tab), ohne auf den Merge des PR zu warten; der Branch bleibt, bis der PR gemergt ist."
context: |-
  Entscheidung von berkc am 22.10.2026: Worktrees werden bei done entfernt, nicht bei review.
  Heute (geprueft): reap.sh (core/role/builtin/jaira-dispatcher/scripts/reap.sh, YD4BE9) entfernt Tab, Worktree und lokalen Branch erst, wenn das Ticket done ist UND sein PR gemergt ist. Laeuft nur, wenn der Teamlead es pro Runde startet. AYBDJ2 (Backlog) betrifft nur die Arbeitskopie nach jaira logbook.
  Neu wegen 6TCTFV: PRs entstehen pro Milestone, nicht pro Ticket. Auf den Merge zu warten hiesse, Worktrees liegen bis zum Ende des Milestones herum.
  Sicherheitsregeln: nur entfernen, wenn git status im Worktree sauber ist (oder nur die eigene Ticket-Datei geaendert) und der Branch auf dem Remote denselben Stand hat; sonst stehen lassen und sagen warum. Branch erst nach Merge loeschen.
  Wie andere Sitzungen wissen, woran gearbeitet wird: claim (claimed-by/claimed-at auf dem Ref, 30-min-Lease), jaira next ueberspringt beanspruchte Tickets; C8V6DQ zeigt es auf dem Board.
definition-of-done: "Erreicht ein Ticket done und ist sein Branch gepusht und der Worktree sauber, wird der Worktree entfernt; der Branch bleibt; Test"
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T00:15:13Z
updated-at: 2026-11-04T21:37:27Z
updated-by: berkc
claimed-by: berk-ca-1279800
claimed-at: 2026-11-04T21:24:25Z
---

# Ein Ticket in done raeumt seinen Worktree weg, sobald der Branch gepusht und sauber ist

## Definition of Done

- [x] Steht das Ticket in done, ist jeder seiner Worktrees sauber (oder nur die eigene Ticket-Datei geaendert) und steht origin/<branch> auf demselben Commit wie der lokale Branch, schliesst reap.sh seine Tabs und entfernt den Worktree, ohne auf einen PR zu warten; der Branch bleibt; Test in core/role/reap_test.go
  proof: core/role/reap_test.go TestReapRemovesAPushedWorktreeBeforeTheMerge; reap.sh:189-194 (push check) and 209-226 (worktree removal before merge)
- [x] Ist das Ticket nicht in done, der Worktree nicht sauber, der Branch nicht gepusht oder origin/<branch> auf anderem Stand, bleibt alles (Tabs, Worktree, Branch) und reap.sh druckt je Ticket eine Zeile mit dem Grund; Test je Fall
  proof: core/role/reap_test.go TestReapKeepsWhatIsNotFinished (subtests not pushed, origin elsewhere, not done, uncommitted, detached HEAD, herdr not reachable); reap.sh:143 waiting, 182 dirty, 189-191 push
- [x] Der lokale Branch wird erst geloescht, wenn sein eigener PR gemergt ist oder sein Stand im Default-Branch von origin liegt (Milestone-PR); reap.sh --dry-run zeigt 'would remove <worktree>' und 'would delete branch <branch>' und aendert nichts; Test
  proof: core/role/reap_test.go TestReapClosesAMergedTicket (--dry-run lines 'would remove', 'would delete branch', nothing changed); TestReapDeletesABranchInOriginsDefaultBranch; reap.sh:248-266
- [ ] Release-Note-Zeile in core/release/notes/80MGXT.md
- [ ] jaira-teamlead/SKILL.md beschreibt den neuen Zeitpunkt (Worktree bei done, Branch nach Merge); neue Hashes von reap.sh und SKILL.md in core/shipped/shipped.sha256 angehaengt, core/role- und core/shipped-Tests gruen

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] reap.sh lesen (core/role/builtin/jaira-dispatcher/scripts/reap.sh) und core/role/reap_test.go; Tests vorher laufen lassen: go test ./core/role -run 'Reap|Spawn' -count=1
- [x] reap.sh Schleife je Ticket neu ordnen, Phase A (Lane): jaira show zuerst; Lane != done -> 'waiting <T>: the ticket is in <lane>, not done' und continue; exit 5 (abgelegt) zaehlt als done; anderer Fehler -> 'kept <T>: jaira show <T> failed (exit N)'
- [x] Branch bestimmen: zuerst Feld branch= aus den Records des Tickets, sonst symbolic-ref eines Worktrees, sonst ${JAIRA_BRANCH_PREFIX:-feat}/<t klein>. Record-Felder lesen mit field(); Records ohne pane= (leer) im Tab-Loop ueberspringen, Records ohne worktree= im Worktree-Loop ueberspringen (passiert schon durch [ -d ])
- [x] Phase B (Worktree, nur wenn wts nicht leer ODER ein Record ein pane= hat): detached/dirty-Check wie heute, aber zusaetzlich jede Status-Zeile ignorieren, deren Pfad unter .jaira/tickets/ liegt und '<T>-' im Dateinamen enthaelt (eigene Ticket-Datei, jeder Status-Code). Dann Push-Check: tip=rev-parse refs/heads/<branch>, up=rev-parse refs/remotes/origin/<branch>; fehlt up -> 'kept <T>: <branch> is not pushed'; up != tip -> 'kept <T>: <branch> is not pushed at its head (origin at <12>, local at <12>)'; beides continue. Dann herdr-Check wie heute, Tabs schliessen und Worktrees entfernen wie heute (--dry-run: would ...)
- [x] Nach erfolgreicher Phase B (nicht --dry-run, kein failed): die Records in gone loeschen und genau einen Record '$panes/branch-<T>' schreiben mit 'ticket=<T>\nbranch=<branch>\n' (der Record des eigenen HERDR_PANE_ID bleibt wie heute). Der Dateiname kollidiert nicht mit Pane-Records? pruefen; sonst anderes Praefix
- [x] Phase C (Branch): existiert refs/heads/<branch> nicht -> Records des Tickets loeschen (nicht bei --dry-run), continue. Ist er in root oder main ausgecheckt -> stehen lassen, Records bleiben, continue (stumm wie heute). Gemergt, wenn tip Vorfahre von refs/remotes/origin/HEAD ist (falls der Ref existiert; Milestone-PR) ODER der Forge-PR gemergt ist und tip == head oder Vorfahre von head (bestehende merged_head-Logik, Forge-Fehlermeldungen wie heute). Kein Merge -> 'waiting <T>: <branch> has no merged pull request'. Tip hinter/neben head -> 'kept <T>: <branch> is not at the merged head ...'. Sonst branch -D, 'deleted branch <branch>' (--dry-run: 'would delete branch'), dann Records des Tickets loeschen
- [x] usage()-Text oben in reap.sh an die neue Reihenfolge anpassen (Worktree+Tabs bei done+sauber+gepusht, Branch nach Merge inkl. origin/HEAD), Kopfkommentar Zeile 2-3 ebenso
- [x] Tests in core/role/reap_test.go: Helfer pushed(t) = git update-ref refs/remotes/origin/feat/ab12cd <tip>. Bestehende Tests anpassen (wo Worktree weg soll: pushed setzen; 'given back' erwartet jetzt waiting...not done; 'committed after the merge' als Branch-Phasen-Test). Neue Tests: (a) done+sauber+gepusht, kein PR -> Worktree weg, Tabs zu, Branch da, Ausgabe 'waiting ... no merged pull request', ein Record mit branch= bleibt; zweiter Lauf nach mergedAt -> Branch weg, keine Records, keine Tabs doppelt geschlossen. (b) Keep-Tabelle: nicht gepusht, origin auf anderem Stand, nicht done, dirty, detached, herdr weg - je Grundzeile, alles bleibt. (c) nur eigene Ticket-Datei geaendert (.jaira/tickets/01X...AB12CD-x.md committed, dann geaendert) -> Worktree wird entfernt. (d) Milestone: tip Vorfahre von origin/HEAD (update-ref refs/remotes/origin/master <tip>; symbolic-ref refs/remotes/origin/HEAD refs/remotes/origin/master), kein PR -> Branch geloescht. (e) --dry-run zeigt 'would remove' und 'would delete branch' und aendert nichts
- [x] core/role/builtin/jaira-teamlead/SKILL.md Zeilen ~36-42 und ~253-262 umschreiben: Tabs und Worktree gehen bei done + sauber + gepusht, lokaler Branch erst nach Merge (eigener PR oder im Default-Branch von origin); Gruende werden als kept genannt
- [x] core/shipped: go test ./core/role ./core/shipped -count=1; der Shipped-Test verlangt neue Hashes von reap.sh und jaira-teamlead/SKILL.md - Zeilen an core/shipped/shipped.sha256 anhaengen wie die Fehlermeldung sagt (alte bleiben)
- [x] core/release/notes/80MGXT.md mit einer Zeile '- ...' (eine Zeile, nicht umbrechen); commit feat(80MGXT): ... mit Code + Ticket-Datei

## Progress
- **2026-10-31 01:25 · berkc** — Maintainer 31.10.2026 (Punkt 7): Integrationsbranch ok, aber das Ziel ist: wird ein Ticket angenommen und steht in done, wird sein Branch gemerged und sein Worktree geloescht - 'ich will nicht den speicher mit worktrees voll machen die nicht geloescht wurden.' Merge bleibt Sache der Person (PR); das Aufraeumen des Worktrees nach done ist dieses Ticket.
- **2026-11-03 17:18 · berkc** — Maintainer 03.11.2026: Worktree-Entfernen ist auch deshalb sicher, weil die Commits im gepushten Branch liegen. Ergaenzung: das gemeinsame Review laeuft in einem eigenen Integrations-Worktree/-Branch je Milestone (Ticket-Branches hineingemergt, PR von dort); die Ticket-Worktrees verschwinden bei done.
- **2026-11-04 21:26 · berkc** — Entscheidungen (Worker, vor dem Plan): (1) 'gepusht' = refs/remotes/origin/<branch> zeigt auf denselben Commit wie refs/heads/<branch>; lokal gelesen, kein fetch/ls-remote - offline-sicher, Tests haben nur einen Fake-Remote git@github.com. git push setzt den Tracking-Ref ohnehin. (2) Milestone-PR: Branch gilt als gemergt, wenn sein eigener PR gemergt ist (wie heute) ODER sein Tip Vorfahre von refs/remotes/origin/HEAD ist. Squash-gemergter Milestone-PR -> Branch bleibt (sichere Seite, lokal ist nichts verloren, Branch liegt auf origin). (3) Ticket nicht in done -> 'waiting <T>: the ticket is in <lane>, not done' statt 'kept'. (4) Nur Worktrees mit spawn.sh-Record (git-common-dir/jaira/panes) werden geraeumt; Worktrees ohne Record bleiben ausser Scope; spawn.sh wird nicht angefasst (DSDPZ4). (5) Nach dem Worktree-Abbau ersetzt reap.sh die Tab-Records des Tickets durch einen Record ticket=/branch= ohne pane, damit die naechste Runde den Branch kennt und keine Tabs erneut schliesst.
