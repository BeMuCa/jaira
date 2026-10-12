---
id: 01M1KAFM6C20J8H7Z3178W1R94
title: "Repo ziehen und jaira updaten ist ein einfacher, durchgespielter Prozess - und der Share-Zustand bleibt dabei immer erhalten"
status: testing
ready: true
creator: BeMuCa
goal: "Der komplette Weg 'jaira installieren/updaten + Board-Repo ziehen/aktualisieren' ist einmal Schritt fuer Schritt durchgespielt, so einfach wie moeglich gemacht und dokumentiert; bei JEDEM Update-Pfad bleibt der Share-Zustand des Boards erhalten (shared bleibt shared, privat bleibt privat)"
context: "Berk am 03.09. vor einem Context-Clear: wenn wir ein Update von jaira ziehen, muss der Share-Zustand im Repo beibehalten werden. Teilweise schon gebaut: ETR0PX (accepted) - jaira update schreibt nur noch den Agent-Block, nie .gitignore; nur jaira init gitignoriert. ZU PRUEFEN: gilt das auf ALLEN Update-Pfaden (jaira self upgrade, go build/install eines neuen Stands, jaira update auf geteiltem vs. privatem Board, erster Befehl nach Binary-Wechsel inkl. Lane-Dir-Migration und Agent-Block-Refresh)? Zweiter Teil: den Prozess einmal als Nutzer durchgehen - frisches Clone eines geteilten Boards, jaira installieren (welcher Befehl? go install trifft NICHT ~/.local/bin - bekannte Falle), Board oeffnen, Update ziehen, jaira update ausfuehren - jede Reibung notieren und glaetten (READMEs, Fehlermeldungen, evtl. ein 'jaira doctor'-artiger Check). DoD-Messlatte: ein Teammate ohne Vorwissen kommt mit README allein durch."
definition-of-done: "Ein dokumentierter Durchlauf clone->install->open->update existiert (Skript oder README-Abschnitt); ein Test oder Durchlauf belegt je Update-Pfad, dass .gitignore/Share-Zustand unangetastet bleibt; gefundene Reibungspunkte sind als Tickets eingefangen oder gefixt"
tags: []
blocked-by: []
commits: []
created-at: 2026-09-03T09:44:29Z
updated-at: 2026-10-12T20:19:02Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-62139
claimed-at: 2026-10-12T20:10:08Z
assignee: Alexander Sacharov
outcome-what: "jaira init laesst .gitignore auf einem schon in git liegenden Board unangetastet (gitrepo.Repo.Tracks + .gitattributes); Test je Update-Pfad (init/update/self upgrade/erster Befehl nach Binary-Wechsel) auf geteiltem Clone und privatem Board; README-Durchlauf clone->install->init->fetch->open->update; init-Hilfe, COMMANDS.md, Release-Note"
outcome-why: "Der natuerliche erste Befehl eines Teammates in einem frischen Clone - jaira init, den jede andere Meldung dort verlangt - haengte /.jaira/ an und machte das geteilte Board wieder privat; das README nannte jaira fetch als ersten Schritt, der in einem solchen Clone scheitert"
outcome-resolves: "Alle Update-Pfade erhalten den Share-Zustand nachweislich (TestUpdatePathsKeepTheShareState), der Beitritt ist im README Schritt fuer Schritt und durchgespielt"
review-summary: none
review-gaps: "sharestate_test.go: eigenen Env-Helper gitIdentity entfernt, Fixtures setzen die Identitaet wie boardremote_test.go/filemode_test.go per 'git config user.*'; DoD-Proof auf die verschobenen Zeilen nachgezogen. Gelassen: gitrepo.Repo.Tracks ist neu, es gab keinen Helfer fuer getrackte Dateien (nur ls-files --others in WorktreeDiff); git ls-files laeuft nur bei init. Vorbestehend und nicht angefasst: 13 rote Tests in internal/cli und TestArchiveFromTheBoard in internal/tui, identisch auf der Basis."
---

# Repo ziehen und jaira updaten ist ein einfacher, durchgespielter Prozess - und der Share-Zustand bleibt dabei immer erhalten

## Definition of Done

- [x] Ein dokumentierter Durchlauf clone->install->open->update existiert (Skript oder README-Abschnitt); ein Test oder Durchlauf belegt je Update-Pfad, dass .gitignore/Share-Zustand unangetastet bleibt; gefundene Reibungspunkte sind als Tickets eingefangen oder gefixt
  proof: README.md:128 (Beitritt clone->install->init->fetch->open) + README.md:86 (Update); TestUpdatePathsKeepTheShareState (internal/cli/sharestate_test.go:78: init/update/self upgrade/erster Befehl nach Binary-Wechsel x geteilt/privat, .gitignore byte-gleich); Reibung gefixt: internal/cli/tickets.go:76 init privatisiert geteiltes Board nicht mehr, TestInitStillMakesAnUntrackedExistingBoardPrivate (sharestate_test.go:139)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] init: nur ein in diesem Lauf angelegtes Board gitignorieren; Ausgabe/JSON melden den echten Share-Zustand
- [x] Test je Update-Pfad: init/update/erster Befehl nach Binary-Wechsel auf geteiltem und privatem Board lassen .gitignore byte-gleich
- [x] Durchlauf clone->install->open->update real durchspielen, README-Abschnitt schreiben, Reibungen fixen oder notieren
- [x] Release-Note fuer die init-Aenderung

## Progress
- **2026-10-12 20:09 · Alexander Sacharov** — Teamlead 2026-10-12: Worker gestartet (jaira-f5). Laeuft parallel zu 5A6RJZ (chat-model/roles) und 0EW5XK (Cockpit-Knoepfe) - keine gemeinsamen Dateien erwartet; Update-/Share-Pfad ist core/update + .gitignore.
- **2026-10-12 20:10 · Alexander Sacharov** — parent session: jaira-f5
- **2026-10-12 20:11 · Alexander Sacharov** — Offene Entscheidungen: keine, die nur eine Person treffen kann. Befund vor dem Plan: einziger Schreiber von /.jaira/ in .gitignore ist board.AddIgnore, aufgerufen von board.Prepare (jaira init, TUI-Taste i) und jaira share. TUI-i verweigert ein existierendes Board (isBoardDir); 'jaira init' NICHT - reproduziert: Clone eines geteilten Boards + 'jaira init' haengt '/.jaira/' an und meldet 'This board is private', d.h. der natuerliche erste Befehl eines Teammates privatisiert das Board. jaira update (ETR0PX), jaira self upgrade (core/selfupdate tauscht nur die Binary) und der erste Befehl nach Binary-Wechsel (nudgeIfStale, bindDriverIfShared, lane.Load) schreiben .gitignore nicht. Entscheidung: init gitignoriert nur ein Board, das es in diesem Lauf selbst anlegt (created==true); ein schon vorhandenes .jaira/ behaelt seinen Share-Zustand, weil Teilen eine bereits getroffene Entscheidung ist (gleiche Begruendung wie update). JSON 'private' und die Schlusszeile von init melden dann den echten Zustand statt fest 'true'.
- **2026-10-12 20:14 · Alexander Sacharov** — Fix gebaut: init entscheidet 'geteilt' an git ls-files .jaira bzw. .jaira/.gitattributes (gitrepo.Repo.Tracks), NICHT an created==false - verworfen, weil 'jaira fetch'/lane.Load/tag .jaira/ vor init anlegen koennen und so ein nie geteiltes Board committbar bliebe (TestInitStillMakesAnUntrackedExistingBoardPrivate). Gegenprobe: alter tickets.go laesst TestUpdatePathsKeepTheShareState rot werden. Zweiter Befund: frischer Clone eines geteilten Boards hat nur .jaira/.gitattributes, kein tickets/ (leere Ordner trackt git nicht, Tickets reisen auf Refs) -> 'jaira fetch' und alle Befehle sagen 'run jaira init'. README sagte 'clone, run jaira fetch' - falsch. Loesung: init ist jetzt im Clone sicher und wird der dokumentierte erste Schritt. Erwogen und verworfen: share schreibt tickets/.gitkeep - hilft nur Boards, die nach der Aenderung geteilt werden.
- **2026-10-12 20:17 · Alexander Sacharov** — Vorbestehend, nicht von hier: 13 Tests in internal/cli (TestDod*, TestSetAppend*, TestLanesAdd*, TestListShowsTheDoDCounter, TestMoveDryRunOnAnAllowedMove, TestReleaseClearsAnExpiredForeignClaim, TestTwoBranchesThatBothCreate...) und internal/tui TestArchiveFromTheBoard (nil-Panic) scheitern identisch auf Basis 43080ed7 in einem frischen Worktree - vermutlich maschinenlokale Einstellungen; nicht in diesem Ticket angefasst. Weitere Reibung gefunden und im README geglaettet: README verwies auf einen 'jaira update'-Abschnitt 'just below', den es nicht gab; zwei Installationen (~/.local/bin vs ~/go/bin) - which -a jaira dokumentiert. Bei 'jaira update' auf einem dev-Build listet es trotz frischem Stempel alle Unreleased-Notes - nur dev, nicht verfolgt.
- **2026-10-12 20:18 · Alexander Sacharov** — critique (Pass 1, ganzer Diff 80daf1bd): nichts. Erwogen und stehen gelassen: (a) nur hasAttributes (bestehendes isShared-Muster, internal/cli/share.go) statt zusaetzlich gitrepo.Repo.Tracks - reicht nicht, ein von Hand geteiltes Board (Ignore-Zeile auskommentiert, nie 'jaira share', Fall aus TestUpdateLeavesASharedBoardShared) hat keine .gitattributes, im Clone aber getrackte Tickets; ohne Tracks wuerde init es privatisieren. (b) Entscheidung im Aufrufer (init) statt in board.Prepare - richtig, Prepare ist laut Doku fuer frisch angelegte Boards, und der TUI-Pfad (browse.go 'i') verweigert existierende Boards schon selbst.
