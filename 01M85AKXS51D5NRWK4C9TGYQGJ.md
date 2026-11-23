---
id: 01M85AKXS51D5NRWK4C9TGYQGJ
title: "Das Katalog-Repo steht: Aufbau, README zum Einreichen, CI mit lint und index"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Es gibt ein oeffentliches Repo, in das jemand eine Lane oder ein Board per Pull Request legt, und CI sagt ihm in Minuten, ob der Beitrag hineindarf."
context: |-
  Heute reicht man eine Lane per Pull Request ins Code-Repo BeMuCa/jaira ein (lanes/README.md, Abschnitt 'Adding yours').
  Das kostet: die ganze Go-Testsuite laeuft (scripts/ci-code-changed.sh zaehlt lanes/*.md als Code), und die Lane ist erst nach dem naechsten Release sichtbar, weil der Client auf den Tag des Binaries zeigt.
  Loesung: ein eigenes Repo nur fuer Inhalt. Wo es liegt und wer mergt, entscheidet das Frage-Ticket in diesem Milestone - vorher nicht anfangen.
  Inhalt des Repos:
  - lanes/ und boards/ nach docs/MARKET.md; changelog-writer.md und secrets-scan.md ziehen aus BeMuCa/jaira/lanes/ um.
  - README: was eine Lane und ein Board ist, wie man einreicht (eine Datei, ein Pull Request), wie man von Hand installiert (Datei nach ~/.jaira/lanes/<id>.md, Board-Ordner nach ~/.jaira/boards/<id>/), was Reviewer pruefen.
  - Pull-Request-Vorlage mit Checkliste; CODEOWNERS, damit kein Beitrag ohne Maintainer hineinkommt; Branch-Schutz auf main.
  - Workflow 'check' auf pull_request (nicht pull_request_target - ein Fork-PR bekommt so keine Secrets): laedt ein festes jaira-Release, prueft dessen Pruefsumme, ruft 'jaira market lint .' auf.
  - Workflow 'index' auf push nach main: 'jaira market index . > index.json', veroeffentlicht ueber GitHub Pages.
  - Fremde Actions sind auf einen vollen Commit-SHA gepinnt.
  lanes/ im Code-Repo bleibt vorerst liegen: alte Binaries zeigen auf ihren Tag und lesen es dort weiter.
  Diese Arbeit findet groesstenteils ausserhalb dieses Repos statt; das Ticket haelt sie hier fest.
definition-of-done: "Das Repo existiert mit lanes/, boards/, README, CONTRIBUTING, PR-Vorlage und CODEOWNERS"
tags:
  - ci
  - docs
blocked-by:
  - 3DWPW6
related: []
commits: []
created-at: 2026-11-23T23:07:20Z
updated-at: 2026-11-23T23:07:20Z
---

# Das Katalog-Repo steht: Aufbau, README zum Einreichen, CI mit lint und index

## Definition of Done

- [ ] Das Repo existiert mit lanes/, boards/, README, CONTRIBUTING, PR-Vorlage und CODEOWNERS
- [ ] Ein Test-PR mit einer Datei ausserhalb der erlaubten Pfade wird von CI rot markiert, einer mit einer gueltigen Lane gruen
- [ ] Nach einem Merge liegt index.json unter der Pages-Adresse und nennt den neuen Eintrag
- [ ] README erklaert Einreichen und Installation von Hand, mit Zielpfaden
- [ ] docs/MARKET.md nennt die Adresse des Repos

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

