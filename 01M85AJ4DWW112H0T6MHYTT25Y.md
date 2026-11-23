---
id: 01M85AJ4DWW112H0T6MHYTT25Y
title: "docs/MARKET.md haelt fest, wie der Katalog fuer Lanes und Boards aufgebaut ist"
status: backlog
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Wer am Marktplatz baut oder einen Eintrag einreicht, liest in einer Datei, wie der Katalog aufgebaut ist, was CI prueft und was bewusst nicht gebaut wird."
context: |-
  Es gibt kein Dokument, das den Marktplatz fuer Lanes und Boards beschreibt. Jedes Ticket im Milestone extension-marketplace braucht dieselben Festlegungen.
  Fuer Plugins (eigene Programme) gibt es das schon: docs/PLUGINS.md, Katalog plugins.json in github.com/panda-bobr/jaira-releases (Ticket R6X6FH). Fuer Lanes und Boards (reine Textdateien) fehlt es.
  Ausloeser: Berks Idee vom 24.11.2026 - ein Erweiterungs-Repo, Beitraege per Pull Request, automatisch geprueft, eine Website, Installation per ID.
  Was heute da ist: 'jaira lanes market' liest das Verzeichnis lanes/ im Repo BeMuCa/jaira ueber die GitHub-Contents-API, festgenagelt auf den Tag des laufenden Binaries (core/market/market.go, internal/cli/market.go). CI prueft nur, dass jede Datei parst (core/lane/shipped_test.go). Kein Pfad-Filter, keine Groessengrenze, keine Pruefsumme, keine Boards.
  Entwurf (Teamlead-Recherche 24.11.2026), das Dokument schreibt ihn aus:
  - Modell: EIN Katalog-Repo, der Inhalt liegt darin (lanes/<id>.md, boards/<id>/ mit description, order, lanes/*.md). Jede Aenderung ist ein Pull Request. Kein Autoren-Repo je Lane: bei Obsidian laufen Updates ohne Review durch, und ein Lane-Prompt ist eine Anweisung an einen Agenten.
  - index.json wird von CI aus dem Repo erzeugt, nie von Hand: je Eintrag kind, id, name, description, path, sha256, min-jaira.
  - Der Client liest nur statische Dateien per HTTPS (<basis>/index.json, <basis>/<path>). Damit ist ein Spiegel eine Kopie des Verzeichnisses.
  - ID: der Dateiname, kebab-case, je Art eindeutig. Aus einer zweiten Quelle: <quelle>/<id>.
  - Vertrauen: Review durch Maintainer, CI-Regeln, Pruefsumme, Inhalt vor dem Uebernehmen zeigen, Uebernehmen ist Sache einer Person.
  - Nicht gebaut: Portal mit Datenbank und Anmeldung, Bewertungen, Signaturen (cosign), Autoren-Repos fuer Lanes, automatische Updates, Verteilung von Skills (dafuer gibt es Claude-Code-Plugin-Marktplaetze).
  Vorbilder mit Quellen: Obsidian (github.com/obsidianmd/obsidian-releases), Claude Code marketplace.json (code.claude.com/docs/en/plugin-marketplaces), Homebrew Taps (docs.brew.sh/Taps), GitHub Actions SHA-Pinning (docs.github.com/en/actions/reference/security/secure-use), VS Code extensions.allowed (code.visualstudio.com/docs/enterprise/extensions), Terraform Registry (developer.hashicorp.com/terraform/registry/modules/publish), Backstage (backstage.io/docs/features/software-templates/adding-templates).
  Haengt an der Entscheidung, wo das Repo liegt (Frage-Ticket im selben Milestone).
definition-of-done: "docs/MARKET.md beschreibt Repo-Aufbau, index.json mit jedem Feld, die ID-Regel und die Regel fuer min-jaira"
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-11-23T23:06:22Z
updated-at: 2026-11-23T23:06:22Z
---

# docs/MARKET.md haelt fest, wie der Katalog fuer Lanes und Boards aufgebaut ist

## Definition of Done

- [ ] docs/MARKET.md beschreibt Repo-Aufbau, index.json mit jedem Feld, die ID-Regel und die Regel fuer min-jaira
- [ ] Das Dokument nennt jede CI-Regel (erlaubte Pfade, Groesse, kein script:, versteckte Zeichen) mit ihrem Grund
- [ ] Ein Abschnitt 'Nicht gebaut' nennt Portal, Bewertungen, Signaturen, Autoren-Repos, Auto-Update, Skills - je mit einem Satz warum
- [ ] docs/PLUGINS.md und lanes/README.md verweisen auf docs/MARKET.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

