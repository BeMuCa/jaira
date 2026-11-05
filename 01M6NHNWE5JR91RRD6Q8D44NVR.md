---
id: 01M6NHNWE5JR91RRD6Q8D44NVR
title: "panda-bobr/jaira-releases: plugins.json, CI-Pruefung und CONTRIBUTING"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Autor bringt sein Plugin per PR in den Katalog, und die CI sagt ihm vor jedem Menschen, was fehlt."
context: |-
  Das Repo https://github.com/panda-bobr/jaira-releases ist seit 05.11.2026 leer.
  Inhalt: plugins.json (leeres Array), README, CONTRIBUTING.md (Fork, Eintrag anhaengen, PR), GitHub-Action, die einen PR prueft.
  CI-Pruefung: id eindeutig und ohne 'jaira'; repo existiert und ist oeffentlich; jaira-plugin.json am Default-Branch passt zum Schema (R6X6FH); Release v<version> hat Assets fuer linux/darwin/windows amd64+arm64 und checksums.txt; PR-Autor ist Besitzer des repo (wie Obsidians GitHub-Verknuepfung).
  Arbeit liegt in einem anderen Repo als jaira: Branch und PR dort, das Ticket faehrt mit einem Notiz-Commit hier.
definition-of-done: "plugins.json, README und CONTRIBUTING.md liegen auf main von panda-bobr/jaira-releases"
tags:
  - ci
blocked-by:
  - R6X6FH
related: []
commits: []
created-at: 2026-11-05T09:47:12Z
updated-at: 2026-11-05T10:28:23Z
updated-by: Alexander Sacharov
---

# panda-bobr/jaira-releases: plugins.json, CI-Pruefung und CONTRIBUTING

## Definition of Done

- [ ] plugins.json, README und CONTRIBUTING.md liegen auf main von panda-bobr/jaira-releases
- [ ] Die Action lehnt einen PR mit doppelter id, fehlendem Release oder fehlender Pruefsumme ab; in einem Test-PR gezeigt

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [ ] Repo jaira-releases lokal anlegen, README-Start-Commit auf main pushen
- [ ] plugins.json, schema/, CONTRIBUTING.md, check_catalog.py + Unit-Tests, Workflow auf feat/d44nvr
- [ ] Fixture-Repo sashasoft90/jaira-plugin-ci-fixture mit Release v0.0.1 ohne checksums.txt
- [ ] Test-Branches (doppelte id, fehlendes Release, fehlende Pruefsumme) pushen, gh-Zeilen fuer Alex

## Progress
- **2026-11-05 10:23 · Alexander Sacharov** — Teamlead (05.11.2026): der Standard (R6X6FH, abgenommen) liegt noch nicht auf master, sondern auf origin/feat/r6x6fh: docs/PLUGINS.md und docs/schema/*.json - lesen mit 'git show origin/feat/r6x6fh:docs/PLUGINS.md'. Die Arbeit selbst gehoert in das Repo panda-bobr/jaira-releases; Alex ist dort admin. Ein Agent pusht seinen Branch und oeffnet keinen PR.
- **2026-11-05 10:24 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/1770.sock
- **2026-11-05 10:24 · Alexander Sacharov** — decision (worker, settled): Pruefung als Python-Skript (stdlib + jsonschema per pip) in .github/scripts/check_catalog.py; Workflow auf pull_request_target, checkt NUR main aus und liest plugins.json des PR-Kopfs als Daten (git show), fuehrt nie PR-Code aus - sonst koennte ein PR das Pruefskript selbst aendern und sich durchwinken. PR, der andere Dateien als plugins.json aendert, wird abgelehnt.
- **2026-11-05 10:24 · Alexander Sacharov** — decision (worker, settled): Schemas aus R6X6FH (docs/schema/*.json, Branch feat/r6x6fh) werden nach jaira-releases/schema/ kopiert, nicht zur CI-Zeit aus jaira geholt - Repo bleibt eigenstaendig; bei Schema-Aenderung in jaira muss die Kopie nachgezogen werden (steht im README). Besitzer-Pruefung: PR-Autor == Besitzer von repo, oder Besitzer ist eine Organisation und der Autor ist dort oeffentliches Mitglied (/orgs/<org>/public_members/<user>, ohne Token abfragbar). jaira-Befehlsnamen als feste Liste im Skript.
- **2026-11-05 10:25 · Alexander Sacharov** — question (open, an Alex): 1) Wie kommt der Inhalt auf main des leeren Repos? (a, empfohlen) Worker pusht nur README als Start-Commit auf main, alles andere auf Branch feat/d44nvr, Alex oeffnet+merged den PR. (b) Worker pusht alles direkt auf main, ohne Review. 2) Test-PR fuer DoD 2: (a, empfohlen) Worker legt oeffentliches Fixture-Repo sashasoft90/jaira-plugin-ci-fixture an (Release v0.0.1 mit jaira-plugin.json, ohne checksums.txt), pusht Test-Branches, Alex oeffnet die Test-PRs mit gelieferten gh-Zeilen. (b) Wie a, aber Worker oeffnet die Test-PRs selbst (als Test markiert, danach geschlossen) - Ausnahme von 'Agents oeffnen keine PRs'. (c) Kein Fixture-Repo: Test-PR zeigt nur doppelte id und fehlendes Release; fehlende Pruefsumme nur im Unit-Test mit aufgezeichneten API-Antworten.
- **2026-11-05 10:27 · Alexander Sacharov** — Alex (05.11.2026, via Teamlead): 1a - nur README als Start-Commit auf main von panda-bobr/jaira-releases, alles andere auf feat/d44nvr, Alex oeffnet und merged den PR. 2a - oeffentliches Fixture-Repo sashasoft90/jaira-plugin-ci-fixture (Release v0.0.1 ohne checksums.txt), Test-Branches pushen, Alex oeffnet die Test-PRs mit fertigen gh-Zeilen.
