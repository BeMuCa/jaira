---
id: 01M6NHNWE5JR91RRD6Q8D44NVR
title: "panda-bobr/jaira-releases: plugins.json, CI-Pruefung und CONTRIBUTING"
status: testing
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
updated-at: 2026-11-05T12:08:29Z
updated-by: Alexander Sacharov
question: "Pushes nach GitHub sind dieser Sitzung verweigert. 1) (empfohlen) Alex fuehrt publish-1-before-merge.sh aus, merged den PR, dann publish-2-after-merge.sh und oeffnet den Test-PR. 2) Alex erlaubt dem Worker git push/gh repo create per Permission-Regel. Zwei Fixture-Repos statt einem - ok?"
outcome-what: "panda-bobr/jaira-releases: README auf main, per PR #1 plugins.json ([]), schema/ (Kopien aus R6X6FH), CONTRIBUTING.md, .github/scripts/check_catalog.py mit 23 Unit-Tests und Workflow check-catalog (pull_request_target, nur plugins.json-PRs, Kommentar mit allen Problemen). Fixture-Repos panda-bobr/jaira-plugin-ci-fixture (Release ohne checksums.txt) und -norelease. Test-PR #2 zeigt die Ablehnung."
outcome-why: "Ein Plugin-Autor soll seinen Eintrag per PR einreichen koennen und vor jedem Menschen erfahren, was fehlt; der Katalog aus R6X6FH existierte nur als leeres Repo."
outcome-resolves: "DoD 1 (main 44d293d) und DoD 2 (PR #2, Run 42024115800, Bot-Kommentar mit doppelter id, fehlendem Release, fehlender Pruefsumme)"
review-summary: none
round: 2
review-gaps: "owner_repo (ein Aufrufer, zwei Zeilen) in check_entry inlined, 310de33. Geprueft und gelassen: get_json (3 Aufrufer), keine toten Imports (pyflakes sauber), kein Duplikat - jaira selbst hat noch keinen Katalog-/Manifest-Validator (8Z3VRF kommt spaeter); die Schema-Kopien in schema/ sind eine bewusste Entscheidung (Notiz). Kosten: pro geaenderter Eintrag 3-5 API-Aufrufe, unveraenderte Eintraege werden nicht abgefragt (Test test_unchanged_entries_are_not_fetched)."
test-verdict: "pass: unit RC=0 (25 tests, fix/d44nvr-errors 310de33); DoD 1 - origin/main 44d293d hat plugins.json, README.md, CONTRIBUTING.md; DoD 2 - PR #2 Bot-Kommentar gelesen, gleicher Lauf lokal gegen echte API mit RC=1 und denselben drei Faellen; gruener Pfad live: hello-demo mit Autor=Besitzer -> 'Catalogue check passed' RC=0"
---

# panda-bobr/jaira-releases: plugins.json, CI-Pruefung und CONTRIBUTING

## Definition of Done

- [x] plugins.json, README und CONTRIBUTING.md liegen auf main von panda-bobr/jaira-releases
  proof: panda-bobr/jaira-releases main 44d293d (Merge PR #1): plugins.json, README.md, CONTRIBUTING.md; check-catalog push run success on 44d293d
- [x] Die Action lehnt einen PR mit doppelter id, fehlendem Release oder fehlender Pruefsumme ab; in einem Test-PR gezeigt
  proof: panda-bobr/jaira-releases PR #2 (draft test PR): check-catalog run 42024115800 failed; bot comment issuecomment-6657527419 lists duplicate id ci-fixture, no release v0.0.1 (ci-fixture-norelease), no checksums.txt (ci-fixture); unit tests .github/scripts/test_check_catalog.py (23)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Repo jaira-releases lokal anlegen, README-Start-Commit auf main pushen
- [x] plugins.json, schema/, CONTRIBUTING.md, check_catalog.py + Unit-Tests, Workflow auf feat/d44nvr
- [x] Fixture-Repo sashasoft90/jaira-plugin-ci-fixture mit Release v0.0.1 ohne checksums.txt
- [x] Test-Branches (doppelte id, fehlendes Release, fehlende Pruefsumme) pushen, gh-Zeilen fuer Alex

## Progress
- **2026-11-05 10:23 · Alexander Sacharov** — Teamlead (05.11.2026): der Standard (R6X6FH, abgenommen) liegt noch nicht auf master, sondern auf origin/feat/r6x6fh: docs/PLUGINS.md und docs/schema/*.json - lesen mit 'git show origin/feat/r6x6fh:docs/PLUGINS.md'. Die Arbeit selbst gehoert in das Repo panda-bobr/jaira-releases; Alex ist dort admin. Ein Agent pusht seinen Branch und oeffnet keinen PR.
- **2026-11-05 10:24 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/1770.sock
- **2026-11-05 10:24 · Alexander Sacharov** — decision (worker, settled): Pruefung als Python-Skript (stdlib + jsonschema per pip) in .github/scripts/check_catalog.py; Workflow auf pull_request_target, checkt NUR main aus und liest plugins.json des PR-Kopfs als Daten (git show), fuehrt nie PR-Code aus - sonst koennte ein PR das Pruefskript selbst aendern und sich durchwinken. PR, der andere Dateien als plugins.json aendert, wird abgelehnt.
- **2026-11-05 10:24 · Alexander Sacharov** — decision (worker, settled): Schemas aus R6X6FH (docs/schema/*.json, Branch feat/r6x6fh) werden nach jaira-releases/schema/ kopiert, nicht zur CI-Zeit aus jaira geholt - Repo bleibt eigenstaendig; bei Schema-Aenderung in jaira muss die Kopie nachgezogen werden (steht im README). Besitzer-Pruefung: PR-Autor == Besitzer von repo, oder Besitzer ist eine Organisation und der Autor ist dort oeffentliches Mitglied (/orgs/<org>/public_members/<user>, ohne Token abfragbar). jaira-Befehlsnamen als feste Liste im Skript.
- **2026-11-05 10:25 · Alexander Sacharov** — question (open, an Alex): 1) Wie kommt der Inhalt auf main des leeren Repos? (a, empfohlen) Worker pusht nur README als Start-Commit auf main, alles andere auf Branch feat/d44nvr, Alex oeffnet+merged den PR. (b) Worker pusht alles direkt auf main, ohne Review. 2) Test-PR fuer DoD 2: (a, empfohlen) Worker legt oeffentliches Fixture-Repo sashasoft90/jaira-plugin-ci-fixture an (Release v0.0.1 mit jaira-plugin.json, ohne checksums.txt), pusht Test-Branches, Alex oeffnet die Test-PRs mit gelieferten gh-Zeilen. (b) Wie a, aber Worker oeffnet die Test-PRs selbst (als Test markiert, danach geschlossen) - Ausnahme von 'Agents oeffnen keine PRs'. (c) Kein Fixture-Repo: Test-PR zeigt nur doppelte id und fehlendes Release; fehlende Pruefsumme nur im Unit-Test mit aufgezeichneten API-Antworten.
- **2026-11-05 10:27 · Alexander Sacharov** — Alex (05.11.2026, via Teamlead): 1a - nur README als Start-Commit auf main von panda-bobr/jaira-releases, alles andere auf feat/d44nvr, Alex oeffnet und merged den PR. 2a - oeffentliches Fixture-Repo sashasoft90/jaira-plugin-ci-fixture (Release v0.0.1 ohne checksums.txt), Test-Branches pushen, Alex oeffnet die Test-PRs mit fertigen gh-Zeilen.
- **2026-11-05 10:29 · Alexander Sacharov** — dead end: git push origin main nach panda-bobr/jaira-releases wurde vom Auto-Mode-Klassifikator dieser Sitzung verweigert. Nicht umgangen. Alles liegt lokal im Scratchpad-Clone (main: README-Commit eb57b38, Branch feat/d44nvr); ein neu gestarteter Worker findet den Clone nicht - Inhalt muss dann neu erzeugt oder Alex pusht. Push ist Alex' Entscheidung.
- **2026-11-05 10:33 · Alexander Sacharov** — dead end: auch 'gh repo create sashasoft90/jaira-plugin-ci-fixture --public' vom Klassifikator verweigert ([Create Public Surface]). Nicht umgangen, nicht an den Teamlead weitergereicht (waere Permission-Laundering). Alles liegt lokal fertig in /tmp/claude-1000/-home-alex-projects--worktrees-jaira-d44nvr/e9c86c25-536c-4790-a6d8-3e163fb430b1/scratchpad: jaira-releases (main eb57b38 README; feat/d44nvr e4f71cd; test/d44nvr-catalog-check fe5c04b), fixtures/ (zwei Repos), publish-1-before-merge.sh und publish-2-after-merge.sh fuer Alex.
- **2026-11-05 10:33 · Alexander Sacharov** — decision (worker): ZWEI Fixture-Repos statt einem - jaira-plugin-ci-fixture (Release v0.0.1 ohne checksums.txt) und jaira-plugin-ci-fixture-norelease (Manifest, kein Release). Grund: die Pruefung liest das Manifest vom Default-Branch, ein Repo kann nicht zugleich 'Release fehlt' und 'Pruefsumme fehlt' zeigen. So zeigt EIN Test-PR alle drei Faelle (DoD 2: doppelte id, fehlendes Release, fehlende Pruefsumme). Pruefskript gegen echte GitHub-API geprueft: nicht existierendes Repo, Org-Repo ohne oeffentliche Mitgliedschaft, fehlendes Manifest werden korrekt gemeldet; 22 Unit-Tests gruen (python -m unittest discover -s .github/scripts).
- **2026-11-05 10:33 · Alexander Sacharov** — gotcha: ein Branch, der als erster in ein leeres Repo gepusht wird, wird Default-Branch. Darum feat/d44nvr erst nach main pushen (publish-1 tut das in dieser Reihenfolge). check-catalog laeuft auf pull_request_target nur mit paths: plugins.json - der PR feat/d44nvr selbst loest ihn nicht aus (Workflow liegt noch nicht auf main); Unit-Job laeuft bei push auf main.
- **2026-11-05 10:37 · Alexander Sacharov** — Alex (05.11.2026, via Teamlead): Veroeffentlichung freigegeben; Fixture-Repos nach panda-bobr statt sashasoft90 (panda-bobr/jaira-plugin-ci-fixture, panda-bobr/jaira-plugin-ci-fixture-norelease). publish-1-before-merge.sh und test/d44nvr-catalog-check (neuer Commit 2d830c6, repo-URLs) angepasst. Teamlead fuehrt publish-1 auf Alex' Anweisung aus; der Worker pusht nichts. Erwartung Test-PR: zusaetzlich zu doppelter id, fehlendem Release, fehlender Pruefsumme meldet die Pruefung 'sashasoft90 is not a public member of panda-bobr', solange Alex' Mitgliedschaft nicht oeffentlich ist. Trockenlauf gegen panda-bobr/jaira-hello-demo (id hello): Manifest, Release v0.1.0, Binaries und checksums.txt bestehen; einzig die Mitgliedschaftspruefung schlaegt fehl - Kandidat fuer einen gruenen PR, sobald die Mitgliedschaft oeffentlich ist.
- **2026-11-05 10:39 · Alexander Sacharov** — Trockenlauf mit den echten Fixture-Repos (nach publish-1): Pruefung meldet doppelte id, fehlendes Release (norelease), fehlendes checksums.txt (ci-fixture) und die Mitgliedschaft. Gefunden: eine doppelte id wurde zweimal gegen GitHub geprueft und doppelt gemeldet -> Commit 6840b4f auf feat/d44nvr prueft je id nur einmal (+Test, 23 gruen). 6840b4f ist LOKAL und muss vor dem Merge noch nach feat/d44nvr gepusht werden; test-Branch darauf rebased.
- **2026-11-05 10:39 · Alexander Sacharov** — PR https://github.com/panda-bobr/jaira-releases/pull/1 offen (feat/d44nvr bei 6840b4f), auf Alex' Anweisung vom Teamlead geoeffnet; Merge macht Alex, danach publish-2 und Test-PR.
- **2026-11-05 11:44 · Alexander Sacharov** — publish-2 umgebaut: kein Force-Push mehr. Worker hat test/d44nvr-catalog-check lokal auf origin/main 44d293d rebased (Diff = nur plugins.json); der Branch war nie gepusht, also reicht 'git push origin test/d44nvr-catalog-check'. Die Ablehnung im Teamlead lag am -f.
- **2026-11-05 12:06 · Alexander Sacharov** — critique (Runde 1): 2 Funde - (1) Nicht-404-HTTP-Fehler ergeben leeren Bot-Kommentar: fangen und als Problem melden, Kommentar-Schritt an report.md binden; (2) event/manifest aus JAIRA_COMMANDS streichen. Bewusst stehen gelassen: Schema-Kopien statt Abruf aus jaira (Notiz 'decision (worker, settled)' zu schema/), Python statt Go (separates Repo, kein Code aus jaira wiederverwendbar). Fix geht als neuer Branch fix/d44nvr-errors in jaira-releases und braucht wieder einen PR von Alex.
- **2026-11-05 12:06 · Alexander Sacharov** — Runde-1-Fix: Branch fix/d44nvr-errors in jaira-releases, Commit 5bace82 (lokal, ungepusht, auf origin/main 44d293d). Dead end: 'if: hashFiles(runner.temp/report.md)' geht nicht - hashFiles sieht nur Dateien unter GITHUB_WORKSPACE; daher Shell-Guard '[ -f report.md ] || exit 0' im Kommentar-Schritt. 25 Unit-Tests gruen.
- **2026-11-05 12:07 · Alexander Sacharov** — critique (Runde 2): beide Funde aus Runde 1 in 5bace82 erledigt (HTTPError/URLError in run() gefangen und als Problem gemeldet, Kommentar-Schritt hinter report.md; event/manifest gestrichen). Keine neuen Funde.
